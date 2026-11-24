package brink

import (
	"bytes"
	"errors"
	"math"
	"testing"
)

func TestBrinkSensorReplyDecodesNamedValue(t *testing.T) {
	// The page's example: wire 0x00C8 on ID 0x07 is 20.0 degrees C (SIR / 10).
	r, err := DecodeSensorReply(0x07, []byte{0xC8, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if !r.Known || r.Name != "inside_temperature" || r.Unit != "deg C" || math.Abs(r.Value-20.0) > 1e-9 {
		t.Fatalf("inside temperature = %#v", r)
	}

	// Unsigned UIR with no divisor: ID 0x01 wire 0x0002 is fan mode 2.
	r, err = DecodeSensorReply(0x01, []byte{0x02, 0x00})
	if err != nil || r.Value != 2 || r.Name != "fan_mode" {
		t.Fatalf("fan mode = %#v, err %v", r, err)
	}

	// STR:13 software version is ASCII text.
	text := []byte("1.2.3-synthet")
	r, err = DecodeSensorReply(0x00, text)
	if err != nil || r.Text != string(text) {
		t.Fatalf("software version = %#v, err %v", r, err)
	}
}

func TestBrinkSensorReplyWidthAndUnknownID(t *testing.T) {
	if _, err := DecodeSensorReply(0x07, []byte{0xC8}); !errors.Is(err, ErrWidth) {
		t.Fatalf("short SIR reply err = %v; want ErrWidth", err)
	}
	// An ID the page does not define keeps the raw payload and is not decoded.
	r, err := DecodeSensorReply(0x7F, []byte{0xAB, 0xCD})
	if err != nil || r.Known || !bytes.Equal(r.Raw, []byte{0xAB, 0xCD}) {
		t.Fatalf("unknown ID = %#v, err %v", r, err)
	}
}

func TestBrinkParameterReplyFiveFields(t *testing.T) {
	// ID 0x01 (UIR, no divisor): current 100, min 50, max 300, step 1, default 100.
	p, err := DecodeParameterReply(0x01, []byte{0x64, 0x00, 0x32, 0x00, 0x2C, 0x01, 0x01, 0x00, 0x64, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if p.Current != 100 || p.Min != 50 || p.Max != 300 || p.Step != 1 || p.Default != 100 || p.Unit != "m3/h" {
		t.Fatalf("flow mode 1 = %#v", p)
	}

	// The page's example for ID 0x04 (SIR, divisor 10): 240 -> 24.0, 150 -> 15.0,
	// 350 -> 35.0, step 5 -> 0.5, default 240 -> 24.0.
	p, err = DecodeParameterReply(0x04, []byte{0xF0, 0x00, 0x96, 0x00, 0x5E, 0x01, 0x05, 0x00, 0xF0, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if p.Current != 24.0 || p.Min != 15.0 || p.Max != 35.0 || p.Step != 0.5 || p.Default != 24.0 {
		t.Fatalf("bypass open temperature = %#v", p)
	}

	if _, err := DecodeParameterReply(0x01, []byte{0x64, 0x00}); !errors.Is(err, ErrWidth) {
		t.Fatalf("short 4050 reply err = %v; want ErrWidth", err)
	}
	if _, err := DecodeParameterReply(0x80, make([]byte, 10)); !errors.Is(err, ErrUnknownID) {
		t.Fatalf("gap ID err = %v; want ErrUnknownID", err)
	}
}

func TestBrinkParameterWriteEncodesDocumentedLayout(t *testing.T) {
	// The page's example: 4080 ID 0x0A with value 4320 is the UIR bytes E0 10.
	got, err := EncodeParameterWrite(0x0A, 4320)
	if err != nil || !bytes.Equal(got, []byte{0x0A, 0xE0, 0x10}) {
		t.Fatalf("filter write = % X, err %v", got, err)
	}
	// ebus address 0xD0 (UIR) set to 0x3C.
	got, err = EncodeParameterWrite(0xD0, 0x3C)
	if err != nil || !bytes.Equal(got, []byte{0xD0, 0x3C, 0x00}) {
		t.Fatalf("ebus address write = % X, err %v", got, err)
	}
	// UCH parameter 0x0F: one value byte.
	got, err = EncodeParameterWrite(0x0F, 1)
	if err != nil || !bytes.Equal(got, []byte{0x0F, 0x01}) {
		t.Fatalf("UCH write = % X, err %v", got, err)
	}
	// SIR parameter 0x04 takes a negative raw value in two's complement.
	got, err = EncodeParameterWrite(0x04, -10)
	if err != nil || !bytes.Equal(got, []byte{0x04, 0xF6, 0xFF}) {
		t.Fatalf("SIR write = % X, err %v", got, err)
	}
	if _, err := EncodeParameterWrite(0x0F, 256); !errors.Is(err, ErrValueWidth) {
		t.Fatalf("UCH overflow err = %v; want ErrValueWidth", err)
	}
	if _, err := EncodeParameterWrite(0x7F, 1); !errors.Is(err, ErrUnknownID) {
		t.Fatalf("gap parameter write err = %v; want ErrUnknownID", err)
	}
}

// Each documented write sub-protocol builds the request the page states,
// with the NN the page states.
func TestBrinkWriteSubProtocolRequests(t *testing.T) {
	cases := []struct {
		name    string
		payload []byte
		nn      int
		sb      byte
		want    []byte
	}{
		{"factory reset 0x40FF", EncodeFactoryReset(), 0, SBFactoryReset, []byte{}},
		{"reset notifications 0x4091", EncodeResetNotifications(), 1, SBResetNotifications, []byte{0x00}},
		{"fan mode temporary 0x40A1", EncodeFanModeTemporary(0x03), 2, SBFanModeTemporary, []byte{0x00, 0x03}},
		{"fan mode persistent 0x40CB", EncodeFanModePersistent(0x01), 2, SBFanModePersistent, []byte{0x00, 0x01}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if len(c.payload) != c.nn || !bytes.Equal(c.payload, c.want) {
				t.Fatalf("payload = % X; want % X (NN %d)", c.payload, c.want, c.nn)
			}
		})
	}
}

// The error history request is ID 0x00; its reply layout is a page gap and is
// kept raw by the catalogue.
func TestBrinkErrorHistoryRequestAndRawReply(t *testing.T) {
	if got := EncodeErrorHistoryRequest(); !bytes.Equal(got, []byte{0x00}) {
		t.Fatalf("error history request = % X", got)
	}
}

func TestBrinkTablesMatchPageCounts(t *testing.T) {
	if len(sensorRegisterTable) != 89 {
		t.Fatalf("4022 registers = %d; the page lists 89", len(sensorRegisterTable))
	}
	if len(parameterTable) != 83 {
		t.Fatalf("4050 parameters = %d; the page lists 83 with a stated ID", len(parameterTable))
	}
	unknown := 0
	for _, p := range parameterTable {
		if !p.RangeKnown {
			unknown++
			if p.Name != "ebus_speed" {
				t.Fatalf("parameter %s has no numeric range; only ebus_speed states varies", p.Name)
			}
		}
	}
	if unknown != 1 {
		t.Fatalf("parameters without a numeric range = %d; want 1", unknown)
	}
	// The page leaves the Valve T01 (0x80-0x92) and MRC (0xA0-0xB9) sub-IDs
	// unresolved; no entry is guessed for them.
	for _, p := range parameterTable {
		if (p.ID >= 0x80 && p.ID <= 0x92) || (p.ID >= 0xA0 && p.ID <= 0xB9) || (p.ID >= 0x50 && p.ID <= 0x7F) {
			t.Fatalf("parameter 0x%02X is in a range the page leaves unassigned", p.ID)
		}
	}
}
