package owner

import (
	"context"
	"errors"
	"fmt"
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/owner/systemconfig"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"reflect"
	"sort"
	"testing"
	"time"
)

func newCfgKeysFakeBus() *cfgKeysFakeBus {
	return &cfgKeysFakeBus{b524: map[uint16][]byte{}, b509: map[uint16][]byte{}}
}

func (bus *cfgKeysFakeBus) send(_ context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	bus.requests++
	data := frame.Data
	if frame.Secondary == vaillantB509Secondary {
		if len(data) < 3 {
			return nil, errors.New("short B509 frame")
		}
		addr := uint16(data[1])<<8 | uint16(data[2])
		switch data[0] {
		case vaillantB509OpcodeRead:
			bus.ops = append(bus.ops, fmt.Sprintf("RB509:%04X", addr))
			value, ok := bus.b509[addr]
			if !ok {
				return nil, errors.New("unseeded B509 register")
			}
			return &protocol.Frame{Data: append([]byte{vaillantB509OpcodeRead, data[1], data[2]}, value...)}, nil
		case vaillantB509OpcodeWrite:
			bus.ops = append(bus.ops, fmt.Sprintf("WB509:%04X", addr))
			payload := append([]byte(nil), data[3:]...)
			bus.b509[addr] = payload
			bus.writes = append(bus.writes, cfgKeysWrite{service: "B509", addr: addr, payload: payload})
			return &protocol.Frame{Data: []byte{vaillantB509OpcodeWrite, data[1], data[2]}}, nil
		}
		return nil, errors.New("unexpected B509 operation")
	}
	if len(data) < 6 || data[0] != systemconfig.OpcodeLocal {
		return nil, errors.New("unexpected B524 frame")
	}
	addr := uint16(data[4]) | uint16(data[5])<<8
	switch data[1] {
	case vaillantB524OpWrite:
		bus.ops = append(bus.ops, fmt.Sprintf("WB524:%04X", addr))
		payload := append([]byte(nil), data[6:]...)
		bus.b524[addr] = payload
		bus.writes = append(bus.writes, cfgKeysWrite{service: "B524", group: data[2], instance: data[3], addr: addr, payload: payload})
		return &protocol.Frame{Data: []byte{0x00}}, nil
	case vaillantB524OpRead:
		bus.ops = append(bus.ops, fmt.Sprintf("RB524:%04X", addr))
		value, ok := bus.b524[addr]
		if !ok {
			return nil, errors.New("unwritten B524 register")
		}
		return testB524ResponseForSelectorPayload(data, value...), nil
	}
	return nil, errors.New("unexpected B524 operation")
}

func (bus *cfgKeysFakeBus) snapshot() (int, []cfgKeysWrite, []string) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	return bus.requests, append([]cfgKeysWrite(nil), bus.writes...), append([]string(nil), bus.ops...)
}

func newCfgKeysPoller(t *testing.T, bus *cfgKeysFakeBus) (*vaillantSemanticPoller, *retainedSemanticState) {
	t.Helper()
	reg := registry.NewDeviceRegistry(nil)
	// Synthetic controller identity without API qualifiers: the DHW bivalence
	// key resolves the conservative unmatched policy bound [-20, 50].
	reg.Register(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV0", SoftwareVersion: "0217", HardwareVersion: "8803"})
	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		scheduler:      engine.NewSemanticReadScheduler(),
		reg:            reg,
		publication:    newVaillantSemanticPublication(provider, nil, nil),
		source:         0x71,
		controller:     0x15,
		boilerAddress:  cfgKeysBoilerAddress,
		requestTimeout: 50 * time.Millisecond,
		system:         &vaillantSystemSnapshot{Controller: 0x15},
		sendFrameFn:    bus.send,
	}
	return poller, provider
}

// cfgKeysSet calls the owner setter named by mutation and returns its result.
func cfgKeysSet(t *testing.T, mutation string, poller *vaillantSemanticPoller, field, value string) (bool, string) {
	t.Helper()
	switch mutation {
	case "set_system_config":
		result := poller.SetSystemConfig(context.Background(), field, value)
		return result.Success, result.Error
	case "set_boiler_config":
		result := poller.SetBoilerConfig(context.Background(), field, value)
		return result.Success, result.Error
	default:
		t.Fatalf("unknown setter %s", mutation)
		return false, ""
	}
}

type cfgKeysCase struct {
	key   string
	value string
	// seed pre-loads B509 registers so the selector resolution reads a value.
	seed map[uint16][]byte
	// writes are the expected acknowledged writes, in order.
	writes []cfgKeysWrite
	// readback is true when the setter confirms its write with a live read of
	// the written register.
	readback bool
	// published reads the value the setter published; nil when the setter
	// publishes no value for this key.
	published func(*retainedSemanticState) any
	want      any
}

func cfgKeysB524(addr uint16, payload ...byte) cfgKeysWrite {
	return cfgKeysWrite{service: "B524", group: 0x00, instance: 0x00, addr: addr, payload: payload}
}

func cfgKeysB509(addr uint16, payload ...byte) cfgKeysWrite {
	return cfgKeysWrite{service: "B509", addr: addr, payload: payload}
}

func cfgKeysPadded(text string) []byte {
	out := make([]byte, 7)
	copy(out, text)
	return out
}

func cfgKeysSystemString(read func(output.SystemConfig) *string) func(*retainedSemanticState) any {
	return func(provider *retainedSemanticState) any {
		status := provider.System()
		if status == nil || read(status.Config) == nil {
			return nil
		}
		return *read(status.Config)
	}
}

func cfgKeysBoilerFloat(read func(output.BoilerConfig) *float64) func(*retainedSemanticState) any {
	return func(provider *retainedSemanticState) any {
		status := provider.BoilerStatus()
		if status == nil || read(status.Config) == nil {
			return nil
		}
		return *read(status.Config)
	}
}

// cfgKeysSystemCases covers every set_system_config key. Per key the case
// catches: the key not accepted (unknown field), a moved register (address),
// a changed value type or encoding (payload), a pair key leaving the pair
// path (one 13-byte write instead of two 7-byte writes), a lost read-back and,
// where the setter publishes a value, a snapshot switch that no longer
// matches the key (published value nil).
func cfgKeysSystemCases() []cfgKeysCase {
	return []cfgKeysCase{
		{key: "dhw_outdoor_bivalence_temperature_limit", value: "7.5", writes: []cfgKeysWrite{cfgKeysB524(0x0001, 0x00, 0x00, 0xF0, 0x40)}, readback: true},
		{key: "room_humidity_limit_maximum", value: "55", writes: []cfgKeysWrite{cfgKeysB524(0x000E, 0x37, 0x00)}, readback: true},
		{key: "adaptive_heating_curve_enabled", value: "true", writes: []cfgKeysWrite{cfgKeysB524(0x0014, 0x01)}, readback: true},
		{key: "heating_outdoor_bivalence_temperature_limit", value: "-3.5", writes: []cfgKeysWrite{cfgKeysB524(0x0023, 0x00, 0x00, 0x60, 0xC0)}, readback: true},
		{key: "heating_flow_temperature_setpoint_maximum_at_heat_pump_fault", value: "45", writes: []cfgKeysWrite{cfgKeysB524(0x0026, 0x00, 0x00, 0x34, 0x42)}, readback: true},
		{key: "dhw_flow_temperature_setpoint_maximum", value: "62", writes: []cfgKeysWrite{cfgKeysB524(0x0046, 0x00, 0x00, 0x78, 0x42)}, readback: true},
		{
			key: "maintenance_date", value: "2027-03-14",
			writes: []cfgKeysWrite{cfgKeysB524(0x002C, 0x0E, 0x03, 0x1B)}, readback: true,
			published: cfgKeysSystemString(func(c output.SystemConfig) *string { return c.MaintenanceDate }), want: "2027-03-14",
		},
		{
			key: "installer_name", value: "Sunflower Co",
			writes:    []cfgKeysWrite{cfgKeysB524(0x006C, cfgKeysPadded("Sunflo")...), cfgKeysB524(0x006D, cfgKeysPadded("wer Co")...)},
			published: cfgKeysSystemString(func(c output.SystemConfig) *string { return c.InstallerName }), want: "Sunflower Co",
		},
		{
			key: "installer_phone_number", value: "+99 123 456",
			writes:    []cfgKeysWrite{cfgKeysB524(0x006F, cfgKeysPadded("+99 12")...), cfgKeysB524(0x0070, cfgKeysPadded("3 456")...)},
			published: cfgKeysSystemString(func(c output.SystemConfig) *string { return c.InstallerPhone }), want: "+99 123 456",
		},
		{
			key: "installer_access_code", value: "417",
			writes: []cfgKeysWrite{cfgKeysB524(0x0076, 0xA1, 0x01)}, readback: true,
			published: func(provider *retainedSemanticState) any {
				status := provider.System()
				if status == nil || status.Config.InstallerMenuCode == nil {
					return nil
				}
				return *status.Config.InstallerMenuCode
			},
			want: 417,
		},
	}
}

// cfgKeysBoilerCases covers every set_boiler_config key, including the
// fallback selector of the heating flow maximum and the 0x4904 key under its
// documentation name, written as one unsigned byte as before.
func cfgKeysBoilerCases() []cfgKeysCase {
	return []cfgKeysCase{
		{
			key: "maximum_heating_flow_temperature_setpoint", value: "55",
			seed:   map[uint16][]byte{0x0E04: {0x20, 0x03}},
			writes: []cfgKeysWrite{cfgKeysB509(0x0E04, 0x70, 0x03)}, readback: true,
			published: cfgKeysBoilerFloat(func(c output.BoilerConfig) *float64 { return c.FlowsetHcMaxC }), want: 55.0,
		},
		{
			key: "maximum_heating_flow_temperature_setpoint", value: "56",
			seed:   map[uint16][]byte{0xA500: {0x20, 0x03}},
			writes: []cfgKeysWrite{cfgKeysB509(0xA500, 0x80, 0x03)}, readback: true,
			published: cfgKeysBoilerFloat(func(c output.BoilerConfig) *float64 { return c.FlowsetHcMaxC }), want: 56.0,
		},
		{
			key: "maximum_dhw_flow_temperature_setpoint", value: "59",
			seed:   map[uint16][]byte{0x0F04: {0x00, 0x03}},
			writes: []cfgKeysWrite{cfgKeysB509(0x0F04, 0xB0, 0x03)}, readback: true,
			published: cfgKeysBoilerFloat(func(c output.BoilerConfig) *float64 { return c.FlowsetHwcMaxC }), want: 59.0,
		},
		{
			key: "heating_partial_load_power", value: "18",
			seed:   map[uint16][]byte{0x0704: {0x0A}},
			writes: []cfgKeysWrite{cfgKeysB509(0x0704, 0x12)}, readback: true,
			published: cfgKeysBoilerFloat(func(c output.BoilerConfig) *float64 { return c.PartloadHcKW }), want: 18.0,
		},
		{
			key: "dhw_partial_load_power", value: "24",
			seed:   map[uint16][]byte{0x0804: {0x0A}},
			writes: []cfgKeysWrite{cfgKeysB509(0x0804, 0x18)}, readback: true,
			published: cfgKeysBoilerFloat(func(c output.BoilerConfig) *float64 { return c.PartloadHwcKW }), want: 24.0,
		},
		{
			key: "installer_phone_number", value: "0123 456 789",
			writes: []cfgKeysWrite{cfgKeysB509(0x8104, 0x01, 0x23, 0x45, 0x67, 0x89, 0xFF, 0xFF, 0xFF)},
			published: func(provider *retainedSemanticState) any {
				status := provider.BoilerStatus()
				if status == nil || status.Config.PhoneNumber == nil {
					return nil
				}
				return *status.Config.PhoneNumber
			},
			want: "0123456789",
		},
		{
			key: "central_heating_pump_mode", value: "17",
			writes: []cfgKeysWrite{cfgKeysB509(0x4904, 0x11)},
			published: func(provider *retainedSemanticState) any {
				status := provider.BoilerStatus()
				if status == nil || status.Config.InstallerMenuCode == nil {
					return nil
				}
				return *status.Config.InstallerMenuCode
			},
			want: 17,
		},
	}
}

func runCfgKeysCase(t *testing.T, mutation string, tc cfgKeysCase) {
	bus := newCfgKeysFakeBus()
	for addr, value := range tc.seed {
		bus.b509[addr] = append([]byte(nil), value...)
	}
	poller, provider := newCfgKeysPoller(t, bus)
	success, message := cfgKeysSet(t, mutation, poller, tc.key, tc.value)
	if !success {
		t.Fatalf("%s(%s=%q) failed: %q", mutation, tc.key, tc.value, message)
	}
	_, writes, ops := bus.snapshot()
	if !reflect.DeepEqual(writes, tc.writes) {
		t.Fatalf("writes = %+v; want %+v", writes, tc.writes)
	}
	if tc.readback {
		last := tc.writes[len(tc.writes)-1]
		write := fmt.Sprintf("W%s:%04X", last.service, last.addr)
		read := fmt.Sprintf("R%s:%04X", last.service, last.addr)
		at := -1
		for i, op := range ops {
			if op == write {
				at = i
			}
		}
		confirmed := false
		for _, op := range ops[at+1:] {
			confirmed = confirmed || op == read
		}
		if at < 0 || !confirmed {
			t.Fatalf("operations = %v; want a read-back of %s after the write", ops, read)
		}
	}
	if tc.published != nil {
		if got := tc.published(provider); got != tc.want {
			t.Fatalf("published value = %#v; want %#v", got, tc.want)
		}
	}
}

type cfgKeysRowDoc struct {
	Key      string `json:"key"`
	Selector struct {
		Opcode    string   `json:"opcode,omitempty"`
		Group     string   `json:"group,omitempty"`
		Instance  string   `json:"instance,omitempty"`
		Registers []string `json:"registers"`
	} `json:"selector"`
	ValueType    string            `json:"value_type"`
	Enum         map[string]uint16 `json:"enum,omitempty"`
	Provisional  bool              `json:"provisional"`
	OpenQuestion string            `json:"open_question,omitempty"`
}

func cfgKeysSystemValueType(valueType systemconfig.ValueType) string {
	switch valueType {
	case systemconfig.ValueFloat32:
		return "float32"
	case systemconfig.ValueUint16:
		return "uint16"
	case systemconfig.ValueBoolUint8:
		return "bool_u8"
	case systemconfig.ValueEnumUint16:
		return "enum_u16"
	case systemconfig.ValueCString:
		return "cstring"
	case systemconfig.ValueDateThreeByte:
		return "date_hda3"
	default:
		return fmt.Sprintf("unknown(%d)", valueType)
	}
}

// cfgKeysSystemRows renders the system setter specs as fixture rows; a key of
// the text-pair writer lists its prefix and suffix registers.
func cfgKeysSystemRows(t *testing.T) []cfgKeysRowDoc {
	specs := systemconfig.FieldSpecs()
	for key, pair := range systemCStringPairSpecs {
		spec, ok := specs[key]
		if !ok || spec.RegisterAddr != pair.addr1 {
			t.Fatalf("pair key %s is not the system spec of its prefix register (spec=%+v)", key, spec)
		}
	}
	keys := make([]string, 0, len(specs))
	for key := range specs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	rows := make([]cfgKeysRowDoc, 0, len(keys))
	for _, key := range keys {
		spec := specs[key]
		row := cfgKeysRowDoc{Key: key, ValueType: cfgKeysSystemValueType(spec.ValueType), Enum: spec.Enum}
		row.Selector.Opcode = fmt.Sprintf("0x%02X", systemconfig.OpcodeLocal)
		row.Selector.Group = fmt.Sprintf("0x%02X", spec.RegisterGroup)
		row.Selector.Instance = fmt.Sprintf("0x%02X", systemconfig.Instance)
		row.Selector.Registers = []string{fmt.Sprintf("0x%04X", spec.RegisterAddr)}
		if pair, ok := systemCStringPairSpecs[key]; ok {
			row.ValueType = "cstring_pair"
			row.Selector.Registers = []string{fmt.Sprintf("0x%04X", pair.addr1), fmt.Sprintf("0x%04X", pair.addr2)}
		}
		rows = append(rows, row)
	}
	return rows
}

// cfgKeysBoilerRows renders the boiler setter: the float specs and the two
// writers with their own codec.
func cfgKeysBoilerRows() []cfgKeysRowDoc {
	rows := make([]cfgKeysRowDoc, 0, len(boilerConfigFieldSpecs)+2)
	add := func(key, valueType string, addrs ...uint16) {
		row := cfgKeysRowDoc{Key: key, ValueType: valueType}
		for _, addr := range addrs {
			row.Selector.Registers = append(row.Selector.Registers, fmt.Sprintf("0x%04X", addr))
		}
		rows = append(rows, row)
	}
	for key, spec := range boilerConfigFieldSpecs {
		valueType := fmt.Sprintf("unknown(%d)", spec.codec)
		switch spec.codec {
		case boilerConfigCodecTempData2c:
			valueType = "data2c"
		case boilerConfigCodecUCH:
			valueType = "uch"
		}
		add(key, valueType, spec.addrs...)
	}
	add("installer_phone_number", "bcd", boilerB509PhoneNumber)
	add("central_heating_pump_mode", "uch", boilerB509InstallerMenuCode)
	sort.Slice(rows, func(i, j int) bool { return rows[i].Key < rows[j].Key })
	return rows
}

func cfgKeysDiff(name string, want, got []cfgKeysRowDoc) error {
	if len(want) != len(got) {
		return fmt.Errorf("%s has %d keys; setter has %d", name, len(got), len(want))
	}
	for i := range want {
		row := got[i]
		// The documentation status is checked against the documentation in the
		// graphql package; the setter fixes everything else of the row.
		row.Provisional = false
		row.OpenQuestion = ""
		if !reflect.DeepEqual(want[i], row) {
			return fmt.Errorf("%s key %d = %+v; setter gives %+v", name, i, row, want[i])
		}
	}
	return nil
}
