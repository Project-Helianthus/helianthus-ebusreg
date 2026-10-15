package system

import "testing"

// Byte-order regression. Register numbers on the wire are little-endian
// (RRlo, RRhi) and decoded as uint16(lo)|uint16(hi)<<8 everywhere else in this
// package. The BASV0 description fixture in helianthus-docs-ebus
// (protocols/vaillant/fixtures/b524-basv0-description-profile-v1.json) records
// OP02/GG00/RR0001 as 0x0001, OP02/GG00/RR0003 as UIN 0..12, OP02/GG02/RR0004
// as EXP 15..80 and OP02/GG03/RR0002 as EXP 15..30 step 0.5.
func TestB524ConstraintCatalogUsesLittleEndianRegisterNumbers(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		group    byte
		register uint16
		min, max float64
		step     float64
	}{
		{"frost_protection_delay GG00 RR0003", 0x00, 0x0003, 0, 12, 1},
		{"circuit_backflow_temperature_setpoint GG02 RR0004", 0x02, 0x0004, 15, 80, 1},
		{"zone_room_temperature_setpoint_cooling GG03 RR0002", 0x03, 0x0002, 15, 30, 0.5},
	}
	for _, tc := range cases {
		got, ok := lookupB524Constraint(B524OpcodeLocal, tc.group, tc.register)
		if !ok {
			t.Errorf("%s: no constraint at little-endian register 0x%04X", tc.name, tc.register)
			continue
		}
		if got.Min != tc.min || got.Max != tc.max || got.Step != tc.step {
			t.Errorf("%s: got %v..%v step %v; want %v..%v step %v", tc.name, got.Min, got.Max, got.Step, tc.min, tc.max, tc.step)
		}
		swapped := tc.register<<8 | tc.register>>8
		if _, found := lookupB524Constraint(B524OpcodeLocal, tc.group, swapped); found {
			t.Errorf("%s: byte-swapped key 0x%04X must not resolve", tc.name, swapped)
		}
	}
}
