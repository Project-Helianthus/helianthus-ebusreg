package system

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusgo/types"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

// Synthetic vectors use the request/response layouts in docs-ebus
// protocols/vaillant/ebus-vaillant-B516-energy.md, sections 3, 5 and 6.
func TestK09DecodedEvidenceOwnsBytes(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name        string
		payload     []byte
		selectorLen int
	}{
		{"request", []byte{0x10, 0x00, 0xFF, 0xFF, 0x04, 0x03, 0x00, 0x30}, 8},
		{"response", []byte{0x00, 0x00, 0xFF, 0x04, 0x03, 0x00, 0x30, 0x00, 0x00, 0x80, 0x3F}, 7},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Each returned slice is tested as the mutator, with a fresh decode pair.
			for _, mutatedField := range []string{"payload", "selector"} {
				t.Run(mutatedField, func(t *testing.T) {
					caller := append([]byte(nil), tc.payload...)
					frame := protocol.Frame{Primary: 0xB5, Secondary: 0x16, Data: caller}
					decoder := newSystemPlane(registry.DeviceInfo{})
					first, handled, err := decoder.DecodeBroadcast(frame)
					if err != nil || !handled {
						t.Fatalf("first decode: handled %v, error %v", handled, err)
					}
					second, handled, err := decoder.DecodeBroadcast(frame)
					if err != nil || !handled {
						t.Fatalf("second decode: handled %v, error %v", handled, err)
					}
					assertEvidence := func(values map[string]types.Value) {
						t.Helper()
						k09AssertBytes(t, values, "payload", tc.payload)
						k09AssertBytes(t, values, "selector", tc.payload[:tc.selectorLen])
						if got := values["selector_is_rsp"]; !got.Valid || got.Value != (tc.selectorLen == 7) {
							t.Errorf("selector_is_rsp = %+v", got)
						}
						if tc.selectorLen == 7 {
							if got := values["wh"]; !got.Valid || got.Value != float32(1) {
								t.Errorf("wh = %+v; want 1 valid", got)
							}
						} else if _, exists := values["wh"]; exists {
							t.Error("request invented energy")
						}
					}
					for i := range caller {
						caller[i] ^= 0xFF
					}
					assertEvidence(first)
					assertEvidence(second)

					mutatedWant := tc.payload
					if mutatedField == "selector" {
						mutatedWant = tc.payload[:tc.selectorLen]
					}
					raw := k09AssertBytes(t, first, mutatedField, mutatedWant)
					for i := range raw {
						raw[i] ^= 0xFF
					}
					otherField, otherWant := "selector", tc.payload[:tc.selectorLen]
					if mutatedField == "selector" {
						otherField, otherWant = "payload", tc.payload
					}
					k09AssertBytes(t, first, otherField, otherWant)
					assertEvidence(second)
					for i, b := range caller {
						if b != tc.payload[i]^0xFF {
							t.Fatalf("returned %s aliases caller at byte %d", mutatedField, i)
						}
					}
				})
			}
		})
	}
}

func TestK09UnknownSelectorsPreserveEvidence(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name           string
		requestOffset  int
		responseOffset int
		field          string
	}{
		{"period", 1, 0, "period"},
		{"source", 4, 3, "source"},
		{"usage", 5, 4, "usage"},
	} {
		for _, response := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/response=%v", tc.name, response), func(t *testing.T) {
				payload := []byte{0x10, 0x00, 0xFF, 0xFF, 0x04, 0x03, 0x00, 0x30}
				offset, selectorLen := tc.requestOffset, 8
				if response {
					payload = []byte{0x00, 0x00, 0xFF, 0x04, 0x03, 0x00, 0x30, 0x00, 0x00, 0x80, 0x3F}
					offset, selectorLen = tc.responseOffset, 7
				}
				payload[offset] = 0x0E
				values, handled, err := newSystemPlane(registry.DeviceInfo{}).DecodeBroadcast(protocol.Frame{Primary: 0xB5, Secondary: 0x16, Data: payload})
				if response && tc.field == "period" {
					if err != nil || handled || values != nil {
						t.Fatalf("unknown response period: values %v, handled %v, error %v; want nil, false, nil", values, handled, err)
					}
					return
				}
				if err != nil || !handled {
					t.Fatalf("decode: handled %v, error %v", handled, err)
				}
				if got, exists := values[tc.field]; !exists || got.Valid || got.Value != nil {
					t.Errorf("unknown %s invented a label: %+v", tc.field, got)
				}
				if got := values[tc.field+"_code"]; !got.Valid || got.Value != uint8(0x0E) {
					t.Errorf("raw code = %+v; want 14 valid", got)
				}
				for _, field := range []string{"year", "month", "day", "year_kind"} {
					if got := values[field]; got.Valid || got.Value != nil {
						t.Errorf("%s invented calendar identity: %+v", field, got)
					}
				}
				if response {
					// Unknown labels do not erase independently decoded native energy.
					if got := values["wh"]; !got.Valid || got.Value != float32(1) {
						t.Errorf("wh = %+v; want wire value 1 valid", got)
					}
				} else if _, exists := values["wh"]; exists {
					t.Error("request invented energy")
				}
				want := append([]byte(nil), payload...)
				for i := range payload {
					payload[i] = 0
				}
				k09AssertBytes(t, values, "payload", want)
				k09AssertBytes(t, values, "selector", want[:selectorLen])
			})
		}
	}
}

func TestK09IncompleteFramesPublishNothing(t *testing.T) {
	t.Parallel()
	for _, full := range [][]byte{
		{0x10, 0x00, 0xFF, 0xFF, 0x04, 0x03, 0x00, 0x30},
		{0x00, 0x00, 0xFF, 0x04, 0x03, 0x00, 0x30, 0x00, 0x00, 0x80, 0x3F},
	} {
		for n := 0; n < len(full); n++ {
			t.Run(fmt.Sprintf("length_%d/truncated_%d", len(full), n), func(t *testing.T) {
				values, handled, err := newSystemPlane(registry.DeviceInfo{}).DecodeBroadcast(protocol.Frame{Primary: 0xB5, Secondary: 0x16, Data: append([]byte(nil), full[:n]...)})
				if err != nil || handled || values != nil {
					t.Fatalf("incomplete frame: values %v, handled %v, error %v; want nil, false, nil", values, handled, err)
				}
			})
		}
	}
}

func k09AssertBytes(t *testing.T, values map[string]types.Value, field string, want []byte) []byte {
	t.Helper()
	value := values[field]
	raw, ok := value.Value.([]byte)
	if !ok || !value.Valid || !bytes.Equal(raw, want) {
		t.Fatalf("%s = %+v; want %x valid", field, value, want)
	}
	return raw
}
