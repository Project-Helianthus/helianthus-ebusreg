package system

import (
	ab10testing "testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

// Synthetic layouts follow the native contract referenced by
// b516_evidence_ownership_test.go: docs-ebus
// protocols/vaillant/ebus-vaillant-B516-energy.md, sections 3, 5 and 6.
func TestAb10PassiveShapesAndIndependentFields(ab10t *ab10testing.T) {
	for _, ab10case := range []struct {
		ab10name                                          string
		ab10period, ab10source, ab10usage, ab10lo, ab10hi byte
		ab10year                                          uint16
		ab10day                                           uint8
		ab10dateValid                                     bool
	}{
		{"ab10_control", 1, 4, 3, 0x21, 0x30, 2024, 1, true},
		{"ab10_period_only", 0, 4, 3, 0x21, 0x30, 0, 0, false},
		{"ab10_source_only", 1, 9, 3, 0x21, 0x30, 2024, 1, true},
		{"ab10_usage_only", 1, 4, 4, 0x21, 0x30, 2024, 1, true},
		{"ab10_date_low_only", 1, 4, 3, 0x22, 0x30, 2024, 2, true},
		{"ab10_date_high_only", 1, 4, 3, 0x21, 0x32, 2025, 1, true},
		{"ab10_malformed_date_preserves_energy", 1, 4, 3, 0x20, 0x30, 0, 0, false},
	} {
		ab10t.Run(ab10case.ab10name, func(ab10t *ab10testing.T) {
			for _, ab10response := range []bool{false, true} {
				ab10name := "ab10_request_eight_bytes"
				if ab10response {
					ab10name = "ab10_reply_eleven_bytes_trailing_float"
				}
				ab10t.Run(ab10name, func(ab10t *ab10testing.T) {
					ab10payload := []byte{0x10, ab10case.ab10period, 0xFF, 0xFF, ab10case.ab10source, ab10case.ab10usage, ab10case.ab10lo, ab10case.ab10hi}
					if ab10response {
						ab10payload = []byte{ab10case.ab10period, 0xFF, 0xFF, ab10case.ab10source, ab10case.ab10usage, ab10case.ab10lo, ab10case.ab10hi, 0x00, 0x50, 0xC3, 0x47}
					}
					ab10values, ab10handled, ab10err := newSystemPlane(registry.DeviceInfo{}).DecodeBroadcast(protocol.Frame{Primary: 0xB5, Secondary: 0x16, Data: ab10payload})
					if ab10err != nil || !ab10handled {
						ab10t.Fatalf("DecodeBroadcast: handled %v, error %v; want handled", ab10handled, ab10err)
					}
					for _, ab10field := range []struct {
						ab10key  string
						ab10want any
					}{
						{"period_code", ab10case.ab10period},
						{"source_code", ab10case.ab10source},
						{"usage_code", ab10case.ab10usage},
						{"selector_w", ab10case.ab10lo >> 4},
						{"selector_v", ab10case.ab10lo & 0x0F},
						{"selector_q", ab10case.ab10hi & 0x0F},
						{"selector_is_rsp", ab10response},
					} {
						if ab10got, ab10exists := ab10values[ab10field.ab10key]; !ab10exists || !ab10got.Valid || ab10got.Value != ab10field.ab10want {
							ab10t.Errorf("%s = %+v; want %v valid", ab10field.ab10key, ab10got, ab10field.ab10want)
						}
					}
					for _, ab10field := range []struct {
						ab10key  string
						ab10want any
					}{
						{"year", ab10case.ab10year},
						{"month", uint8(1)},
						{"day", ab10case.ab10day},
					} {
						ab10got, ab10exists := ab10values[ab10field.ab10key]
						if !ab10exists || ab10got.Valid != ab10case.ab10dateValid || (ab10case.ab10dateValid && ab10got.Value != ab10field.ab10want) {
							ab10t.Errorf("%s = %+v; want %v, valid %v", ab10field.ab10key, ab10got, ab10field.ab10want, ab10case.ab10dateValid)
						}
					}
					if ab10response {
						if ab10got := ab10values["wh"]; !ab10got.Valid || ab10got.Value != float32(100000) {
							ab10t.Errorf("trailing float wh = %+v; want 100000 valid", ab10got)
						}
					} else if _, ab10exists := ab10values["wh"]; ab10exists {
						ab10t.Error("request must not emit energy")
					}
				})
			}
		})
	}
}

func TestAb10ActiveThreeByteBodyIsNotPassive(ab10t *ab10testing.T) {
	for _, ab10case := range []struct {
		ab10name    string
		ab10payload []byte
	}{
		{"ab10_control", []byte{1, 4, 3}},
		{"ab10_period_only", []byte{2, 4, 3}},
		{"ab10_source_only", []byte{1, 9, 3}},
		{"ab10_usage_only", []byte{1, 4, 4}},
	} {
		ab10t.Run(ab10case.ab10name, func(ab10t *ab10testing.T) {
			ab10values, ab10handled, ab10err := newSystemPlane(registry.DeviceInfo{}).DecodeBroadcast(protocol.Frame{Primary: 0xB5, Secondary: 0x16, Data: ab10case.ab10payload})
			if ab10err != nil || ab10handled || ab10values != nil {
				ab10t.Fatalf("active three-byte body decoded as passive: values %v, handled %v, error %v", ab10values, ab10handled, ab10err)
			}
		})
	}
}
