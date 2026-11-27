package owner

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"strings"
	"testing"
	"time"
)

// TestCircuitSpecialFunctionModePublishesOneByteReply reads the one-byte
// special-function selector and publishes it unchanged on the circuit record
// and its MCP JSON.
func TestVaillantOwnerCircuitSpecialFunctionModePublishesOneByteReply(t *testing.T) {
	for _, value := range []byte{0x00, 0x03, 0x7E} {
		t.Run(fmt.Sprintf("value_%d", value), func(t *testing.T) {
			now := time.Unix(1000, 0)
			poller, provider := sfmCircuitPoller(&now, map[uint16][]byte{sfmTestRegister: {value}})
			before := sfmCounter(t, "u8")
			poller.refreshCircuits(context.Background())

			config, mcpConfig := sfmPublishedConfig(t, provider)
			if config.SpecialFunctionMode == nil || *config.SpecialFunctionMode != int(value) {
				t.Fatalf("circuit special function mode = %v; want %d", config.SpecialFunctionMode, value)
			}
			got, ok := mcpConfig["circuit_special_function_mode"]
			if !ok {
				t.Fatalf("MCP circuit config lacks circuit_special_function_mode: %v", mcpConfig)
			}
			if number, ok := got.(float64); !ok || number != float64(value) {
				t.Fatalf("MCP circuit_special_function_mode = %#v; want %d", got, value)
			}
			if delta := sfmCounter(t, "u8") - before; delta != 0 {
				t.Fatalf("one-byte reply counted as mismatch %d times", delta)
			}
		})
	}
}

// TestCircuitSpecialFunctionModeHasNoAliasAndNoWrite checks that the retired
// field is absent from the published record and its MCP JSON, and that a
// circuit refresh does not write the bus.
func TestVaillantOwnerCircuitSpecialFunctionModeHasNoAliasAndNoWrite(t *testing.T) {
	now := time.Unix(1000, 0)
	poller, provider := sfmCircuitPoller(&now, map[uint16][]byte{sfmTestRegister: {0x02}})
	writes := 0
	inner := poller.sendFrameFn
	poller.sendFrameFn = func(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
		if len(frame.Data) >= 2 && frame.Data[1] == vaillantB524OpWrite {
			writes++
		}
		return inner(ctx, frame)
	}
	poller.refreshCircuits(context.Background())
	if writes != 0 {
		t.Fatalf("circuit refresh writes = %d; want none", writes)
	}
	config, mcpConfig := sfmPublishedConfig(t, provider)
	retired := "frost" + "_prot_c"
	if _, ok := mcpConfig[retired]; ok {
		t.Errorf("MCP circuit config still publishes the retired field")
	}
	raw, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	var encoded map[string]any
	if err := json.Unmarshal(raw, &encoded); err != nil {
		t.Fatal(err)
	}
	if _, ok := encoded[retired]; ok {
		t.Errorf("published circuit config still has the retired field")
	}
}

// TestB524ReplyWidthMismatchOnCircuitRefresh covers both directions on the
// circuit refresh: a four-byte RR001D reply to the one-byte reader and a
// one-byte reply to a four-byte reader. Each is counted once under its reader
// kind, publishes no value and logs one line with selector and widths.
func TestVaillantOwnerB524ReplyWidthMismatchOnCircuitRefresh(t *testing.T) {
	cases := []struct {
		name      string
		addr      uint16
		payload   []byte
		kind      string
		want, got int
		absent    func(output.CircuitConfig) bool
	}{
		{
			name: "four_bytes_to_one_byte_reader", addr: sfmTestRegister, payload: []byte{0x00, 0x00, 0xA0, 0x41},
			kind: "u8", want: 1, got: 4,
		},
		{
			name: "one_byte_to_four_byte_reader", addr: sfmTestFloatRegister, payload: []byte{0x2A},
			kind: "float32le", want: 4, got: 1,
			absent: func(config output.CircuitConfig) bool { return config.FlowTempMaxC == nil },
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lines := sfmCaptureLog(t)
			now := time.Unix(1000, 0)
			poller, provider := sfmCircuitPoller(&now, map[uint16][]byte{tc.addr: tc.payload})
			before := sfmCounter(t, tc.kind)
			poller.refreshCircuits(context.Background())

			if delta := sfmCounter(t, tc.kind) - before; delta != 1 {
				t.Fatalf("%s mismatch counter delta = %d; want 1", tc.kind, delta)
			}
			wantLine := sfmMismatchLine(0x02, sfmTestCircuitGroup, tc.addr, tc.want, tc.got)
			if got := lines(); len(got) != 1 || got[0] != wantLine {
				t.Fatalf("mismatch log = %q; want exactly %q", got, wantLine)
			}
			config, mcpConfig := sfmPublishedConfig(t, provider)
			if tc.absent == nil {
				if config.SpecialFunctionMode != nil {
					t.Fatalf("published special function mode = %v for a wrong-width reply", *config.SpecialFunctionMode)
				}
				if value, ok := mcpConfig["circuit_special_function_mode"]; ok && value != nil {
					t.Fatalf("MCP publishes %v for a wrong-width reply", value)
				}
			} else if !tc.absent(config) {
				t.Fatalf("wrong-width reply published a value: %+v", config)
			}
		})
	}
}

// TestB524ReplyWidthMismatchEveryFixedWidthReader sends one reply of the
// wrong width to each fixed-width reader on a selector of its own, and one
// reply of the right width, and checks value, counter and log for both.
func TestVaillantOwnerB524ReplyWidthMismatchEveryFixedWidthReader(t *testing.T) {
	type reader func(p *vaillantSemanticPoller, sel sfmSelector) bool
	cases := []struct {
		name     string
		kind     string
		sel      sfmSelector
		wrong    []byte
		right    []byte
		want     int
		hasValue reader
	}{
		{"u8", "u8", sfmSelector{0x06, 0x0A, 0x01, 0x0019}, []byte{0x05, 0x00}, []byte{0x05}, 1,
			func(p *vaillantSemanticPoller, s sfmSelector) bool {
				return p.readB524Uint8(context.Background(), s.opcode, s.group, s.instance, s.addr) != nil
			}},
		{"bool", "u8", sfmSelector{0x06, 0x0A, 0x01, 0x001E}, []byte{0x01, 0x00, 0x00, 0x00}, []byte{0x01}, 1,
			func(p *vaillantSemanticPoller, s sfmSelector) bool {
				return p.readB524Bool(context.Background(), s.opcode, s.group, s.instance, s.addr) != nil
			}},
		{"uint16", "uint16", sfmSelector{0x02, 0x02, 0x00, 0x0015}, []byte{0x01, 0x00, 0x00}, []byte{0x01, 0x00}, 2,
			func(p *vaillantSemanticPoller, s sfmSelector) bool {
				v, ok := p.readB524Uint16(context.Background(), s.opcode, s.group, s.instance, s.addr)
				return ok && v != nil
			}},
		{"u16", "uint16", sfmSelector{0x06, 0x0C, 0x02, 0x0023}, []byte{0x01, 0x02, 0x03, 0x04}, []byte{0x34, 0x12}, 2,
			func(p *vaillantSemanticPoller, s sfmSelector) bool {
				return p.readB524Uint16Pointer(context.Background(), s.opcode, s.group, s.instance, s.addr) != nil
			}},
		{"float32le", "float32le", sfmSelector{0x02, 0x03, 0x01, 0x0014}, []byte{0x00, 0x00, 0xA0}, testB524Float32Payload(20), 4,
			func(p *vaillantSemanticPoller, s sfmSelector) bool {
				_, ok := p.readB524Float32LE(context.Background(), s.opcode, s.group, s.instance, s.addr)
				return ok
			}},
		{"f32", "float32le", sfmSelector{0x06, 0x09, 0x03, 0x000F}, []byte{0x00, 0x00, 0xA0, 0x41, 0x00}, testB524Float32Payload(21), 4,
			func(p *vaillantSemanticPoller, s sfmSelector) bool {
				return p.readB524Float32(context.Background(), s.opcode, s.group, s.instance, s.addr) != nil
			}},
		{"uint32le", "uint32le", sfmSelector{0x02, 0x02, 0x01, 0x0024}, []byte{0x10, 0x00}, []byte{0x10, 0x00, 0x00, 0x00}, 4,
			func(p *vaillantSemanticPoller, s sfmSelector) bool {
				_, ok := p.readB524Uint32LE(context.Background(), s.opcode, s.group, s.instance, s.addr)
				return ok
			}},
		{"energy_register", "uint32le", sfmSelector{vaillantB524OpcodeLocal, localRegulator.group, regulatorInstance, b524EnergyQueries[0].addr}, []byte{0x10, 0x00, 0x00, 0x00, 0x00}, []byte{0x10, 0x00, 0x00, 0x00}, 4,
			func(p *vaillantSemanticPoller, s sfmSelector) bool {
				_, ok := p.readEnergyRegister(context.Background(), b524EnergyQueries[0])
				return ok
			}},
		{"date_hda3", "date_hda3", sfmSelector{0x02, 0x00, 0x00, 0x002C}, []byte{0x02, 0x01}, []byte{0x02, 0x01, 0x18}, 3,
			func(p *vaillantSemanticPoller, s sfmSelector) bool {
				_, ok := p.readB524DateThreeByte(context.Background(), s.opcode, s.group, s.instance, s.addr)
				return ok
			}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lines := sfmCaptureLog(t)
			now := time.Unix(1000, 0)

			right := sfmReaderPoller(&now, map[sfmSelector][]byte{tc.sel: tc.right})
			before := sfmCounter(t, tc.kind)
			if !tc.hasValue(right, tc.sel) {
				t.Fatalf("right-width reply of %d bytes returned no value", len(tc.right))
			}
			if delta := sfmCounter(t, tc.kind) - before; delta != 0 || len(lines()) != 0 {
				t.Fatalf("right-width reply reported: counter delta %d, log %q", delta, lines())
			}

			wrong := sfmReaderPoller(&now, map[sfmSelector][]byte{tc.sel: tc.wrong})
			before = sfmCounter(t, tc.kind)
			if tc.hasValue(wrong, tc.sel) {
				t.Fatalf("wrong-width reply of %d bytes returned a value", len(tc.wrong))
			}
			if delta := sfmCounter(t, tc.kind) - before; delta != 1 {
				t.Fatalf("%s counter delta = %d; want 1", tc.kind, delta)
			}
			wantLine := sfmMismatchLine(tc.sel.opcode, tc.sel.group, tc.sel.addr, tc.want, len(tc.wrong))
			if got := lines(); len(got) != 1 || got[0] != wantLine {
				t.Fatalf("mismatch log = %q; want exactly %q", got, wantLine)
			}
		})
	}
}

// TestB524ReplyWidthMismatchCircuitRR0020StaysSilent keeps the deliberate
// rejection of the circuit RR0020 selector out of the mismatch report.
func TestVaillantOwnerB524ReplyWidthMismatchCircuitRawFlowRegisterStaysSilent(t *testing.T) {
	lines := sfmCaptureLog(t)
	now := time.Unix(1000, 0)
	for _, payload := range [][]byte{{0x01}, testB524Float32Payload(46.25)} {
		sel := sfmSelector{0x02, 0x02, 0x00, 0x0020}
		poller := sfmReaderPoller(&now, map[sfmSelector][]byte{sel: payload})
		before := sfmCounter(t, "float32le")
		if _, ok := poller.readB524Float32LE(context.Background(), sel.opcode, sel.group, sel.instance, sel.addr); ok {
			t.Fatal("circuit RR0020 returned a value")
		}
		if delta := sfmCounter(t, "float32le") - before; delta != 0 || len(lines()) != 0 {
			t.Fatalf("circuit RR0020 rejection reported: counter delta %d, log %q", delta, lines())
		}
	}
}

// TestB524ReplyWidthMismatchLogIsRateLimitedPerSelector checks the ten-minute
// window per selector. Every field of the logged selector is varied on its
// own (opcode, group, register): a change of any one of them is a different
// selector and logs at once. The instance is not part of the logged selector.
func TestVaillantOwnerB524ReplyWidthMismatchLogIsRateLimitedPerSelector(t *testing.T) {
	lines := sfmCaptureLog(t)
	now := time.Unix(1000, 0)
	base := sfmSelector{0x02, 0x02, 0x00, 0x001D}
	otherOpcode := sfmSelector{0x06, 0x02, 0x00, 0x001D}
	otherGroup := sfmSelector{0x02, 0x03, 0x00, 0x001D}
	otherRegister := sfmSelector{0x02, 0x02, 0x00, 0x001E}
	otherInstance := sfmSelector{0x02, 0x02, 0x01, 0x001D}
	wrong := []byte{0x01, 0x02, 0x03, 0x04}
	poller := sfmReaderPoller(&now, map[sfmSelector][]byte{
		base: wrong, otherOpcode: wrong, otherGroup: wrong, otherRegister: wrong, otherInstance: wrong,
	})
	read := func(sel sfmSelector) {
		t.Helper()
		if poller.readB524Uint8(context.Background(), sel.opcode, sel.group, sel.instance, sel.addr) != nil {
			t.Fatal("wrong-width reply returned a value")
		}
	}
	line := func(sel sfmSelector) string { return sfmMismatchLine(sel.opcode, sel.group, sel.addr, 1, 4) }

	before := sfmCounter(t, "u8")
	read(base)
	now = now.Add(time.Minute)
	read(base)
	if delta := sfmCounter(t, "u8") - before; delta != 2 {
		t.Fatalf("counter delta after two mismatches = %d; want 2", delta)
	}
	if got := lines(); len(got) != 1 || got[0] != line(base) {
		t.Fatalf("log after two mismatches of one selector = %q; want one line", got)
	}

	read(otherInstance)
	if got := lines(); len(got) != 1 {
		t.Fatalf("another instance of the same logged selector logged again: %q", got)
	}
	want := []string{line(base)}
	for _, sel := range []sfmSelector{otherOpcode, otherGroup, otherRegister} {
		read(sel)
		want = append(want, line(sel))
		if got := lines(); strings.Join(got, "\n") != strings.Join(want, "\n") {
			t.Fatalf("log = %q; want %q", got, want)
		}
	}

	now = time.Unix(1000, 0).Add(10*time.Minute - time.Second)
	read(base)
	if got := lines(); len(got) != len(want) {
		t.Fatalf("selector logged again inside the window: %q", got)
	}
	now = time.Unix(1000, 0).Add(10 * time.Minute)
	read(base)
	want = append(want, line(base))
	if got := lines(); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("selector not logged again after the window: %q", got)
	}
	if delta := sfmCounter(t, "u8") - before; delta != 8 {
		t.Fatalf("counter delta = %d; want 8 (every mismatch counted)", delta)
	}
}
