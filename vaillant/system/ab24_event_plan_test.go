package system

import (
	"fmt"
	"reflect"
	"testing"
)

func ab24Fixture() Ab24EventEdit {
	ab24Identity := Ab24EventIdentity{"ab24-synthetic-family", 3, 2, 1, 255}
	return Ab24EventEdit{
		Ab24Identity: ab24Identity, Ab24Opcode: 0x0a,
		Ab24Values: []int{0, 1, 2, 3, 4, 5, 6},
		Ab24OP09:   &Ab24EventBaseline{ab24Identity, 0x09, true, []byte{0, 1, 2, 3, 4, 5, 6, 7}},
		Ab24OP0B:   &Ab24EventBaseline{ab24Identity, 0x0b, true, []byte{8, 9, 10, 11, 12, 13, 14, 15}},
	}
}

func TestAb24ExactVectorsAndOwnership(ab24T *testing.T) {
	for _, ab24Opcode := range []int{0x0a, 0x0c} {
		ab24T.Run(fmt.Sprintf("ab24_opcode_%02x", ab24Opcode), func(ab24T *testing.T) {
			ab24Edit := ab24Fixture()
			ab24Edit.Ab24Opcode = ab24Opcode
			ab24Plan, ab24Err := Ab24BuildEventPlan(ab24Edit)
			ab24Want := [12]byte{byte(ab24Opcode), 3, 2, 1, 255, 0, 1, 2, 3, 4, 5, 6}
			if ab24Err != nil || ab24Plan == nil {
				ab24T.Fatalf("plan: %v", ab24Err)
			}
			if ab24Plan.Ab24Payload != ab24Want || ab24Plan.Ab24Identity != ab24Edit.Ab24Identity || ab24Plan.Ab24Values != [7]byte{0, 1, 2, 3, 4, 5, 6} || ab24Plan.Ab24OP09 != [8]byte{0, 1, 2, 3, 4, 5, 6, 7} || ab24Plan.Ab24OP0B != [8]byte{8, 9, 10, 11, 12, 13, 14, 15} {
				ab24T.Fatalf("unexpected plan: %+v", ab24Plan)
			}
			ab24Before := *ab24Plan
			ab24Edit.Ab24Values[0] = 99
			ab24Edit.Ab24OP09.Ab24Raw[0] = 99
			ab24Edit.Ab24OP0B.Ab24Raw[0] = 99
			ab24Edit.Ab24Identity.Ab24Family = "ab24-other-family"
			if *ab24Plan != ab24Before {
				ab24T.Fatal("input mutation changed owned plan")
			}
		})
	}
}

func TestAb24IdentityAndRawValues(ab24T *testing.T) {
	for ab24Field, ab24Value := range map[string]any{"Ab24Family": "ab24-other-family", "Ab24SystemType": 0, "Ab24Instance": 255, "Ab24Address": 2, "Ab24Code": 0} {
		ab24T.Run("ab24_"+ab24Field, func(ab24T *testing.T) {
			ab24Edit := ab24Fixture()
			reflect.ValueOf(&ab24Edit.Ab24Identity).Elem().FieldByName(ab24Field).Set(reflect.ValueOf(ab24Value))
			ab24Edit.Ab24OP09.Ab24Identity = ab24Edit.Ab24Identity
			ab24Edit.Ab24OP0B.Ab24Identity = ab24Edit.Ab24Identity
			ab24Plan, ab24Err := Ab24BuildEventPlan(ab24Edit)
			if ab24Err != nil || ab24Plan == nil {
				ab24T.Fatalf("valid identity: %v", ab24Err)
			}
			ab24Want := [12]byte{0x0a, byte(ab24Edit.Ab24Identity.Ab24SystemType), byte(ab24Edit.Ab24Identity.Ab24Instance), byte(ab24Edit.Ab24Identity.Ab24Address), byte(ab24Edit.Ab24Identity.Ab24Code), 0, 1, 2, 3, 4, 5, 6}
			if ab24Plan.Ab24Identity != ab24Edit.Ab24Identity || ab24Plan.Ab24Payload != ab24Want {
				ab24T.Fatalf("identity lost: %+v", ab24Plan)
			}
		})
	}
	for ab24Index := 0; ab24Index < 7; ab24Index++ {
		ab24T.Run(fmt.Sprintf("ab24_value_%d", ab24Index), func(ab24T *testing.T) {
			ab24Edit := ab24Fixture()
			ab24Edit.Ab24Values[ab24Index] = 255
			ab24Plan, ab24Err := Ab24BuildEventPlan(ab24Edit)
			if ab24Err != nil || ab24Plan == nil {
				ab24T.Fatalf("raw value: %v", ab24Err)
			}
			ab24Want := [12]byte{0x0a, 3, 2, 1, 255, 0, 1, 2, 3, 4, 5, 6}
			ab24Want[5+ab24Index] = 255
			ab24WantValues := [7]byte{0, 1, 2, 3, 4, 5, 6}
			ab24WantValues[ab24Index] = 255
			if ab24Plan.Ab24Payload != ab24Want || ab24Plan.Ab24Values != ab24WantValues {
				ab24T.Fatalf("value order: %+v", ab24Plan)
			}
		})
	}
}

func TestAb24RejectBaseline(ab24T *testing.T) {
	for _, ab24Slot := range []string{"ab24_op09", "ab24_op0b"} {
		for _, ab24Case := range []string{"missing", "failed", "opcode", "short", "long", "family", "system", "instance", "address", "code"} {
			ab24T.Run(ab24Slot+"_"+ab24Case, func(ab24T *testing.T) {
				ab24Edit := ab24Fixture()
				ab24Baseline := &ab24Edit.Ab24OP09
				if ab24Slot == "ab24_op0b" {
					ab24Baseline = &ab24Edit.Ab24OP0B
				}
				switch ab24Case {
				case "missing":
					*ab24Baseline = nil
				case "failed":
					(*ab24Baseline).Ab24Success = false
				case "opcode":
					(*ab24Baseline).Ab24Opcode = 5
				case "short":
					(*ab24Baseline).Ab24Raw = (*ab24Baseline).Ab24Raw[:7]
				case "long":
					(*ab24Baseline).Ab24Raw = append((*ab24Baseline).Ab24Raw, 9)
				case "family":
					(*ab24Baseline).Ab24Identity.Ab24Family = "ab24-other-family"
				case "system":
					(*ab24Baseline).Ab24Identity.Ab24SystemType = 1
				case "instance":
					(*ab24Baseline).Ab24Identity.Ab24Instance = 3
				case "address":
					(*ab24Baseline).Ab24Identity.Ab24Address = 2
				case "code":
					(*ab24Baseline).Ab24Identity.Ab24Code = 254
				}
				ab24Plan, ab24Err := Ab24BuildEventPlan(ab24Edit)
				if ab24Err == nil || ab24Plan != nil {
					ab24T.Fatalf("accepted invalid baseline: %+v %v", ab24Plan, ab24Err)
				}
			})
		}
	}
}

func TestAb24RejectFields(ab24T *testing.T) {
	for _, ab24Field := range []string{"Ab24SystemType", "Ab24Instance", "Ab24Address", "Ab24Code", "Ab24Opcode"} {
		for _, ab24Value := range []int{-1, 256} {
			ab24T.Run(fmt.Sprintf("ab24_%s_%d", ab24Field, ab24Value), func(ab24T *testing.T) {
				ab24Edit := ab24Fixture()
				if ab24Field == "Ab24Opcode" {
					ab24Edit.Ab24Opcode = ab24Value
				} else {
					reflect.ValueOf(&ab24Edit.Ab24Identity).Elem().FieldByName(ab24Field).SetInt(int64(ab24Value))
					ab24Edit.Ab24OP09.Ab24Identity = ab24Edit.Ab24Identity
					ab24Edit.Ab24OP0B.Ab24Identity = ab24Edit.Ab24Identity
				}
				ab24Plan, ab24Err := Ab24BuildEventPlan(ab24Edit)
				if ab24Err == nil || ab24Plan != nil {
					ab24T.Fatal("accepted narrowed field")
				}
			})
		}
	}
	for ab24Index := 0; ab24Index < 7; ab24Index++ {
		for _, ab24Value := range []int{-1, 256} {
			ab24T.Run(fmt.Sprintf("ab24_value_%d_%d", ab24Index, ab24Value), func(ab24T *testing.T) {
				ab24Edit := ab24Fixture()
				ab24Edit.Ab24Values[ab24Index] = ab24Value
				ab24Plan, ab24Err := Ab24BuildEventPlan(ab24Edit)
				if ab24Err == nil || ab24Plan != nil {
					ab24T.Fatal("accepted narrowed value")
				}
			})
		}
	}
}

func TestAb24ShapeAndSelectors(ab24T *testing.T) {
	for _, ab24Case := range []string{"six", "eight", "op05", "empty_family", "system02", "address00", "system00_address04", "system01_address03", "system03_address03"} {
		ab24T.Run("ab24_"+ab24Case, func(ab24T *testing.T) {
			ab24Edit := ab24Fixture()
			switch ab24Case {
			case "six":
				ab24Edit.Ab24Values = ab24Edit.Ab24Values[:6]
			case "eight":
				ab24Edit.Ab24Values = append(ab24Edit.Ab24Values, 7)
			case "op05":
				ab24Edit.Ab24Opcode = 5
			case "empty_family":
				ab24Edit.Ab24Identity.Ab24Family = ""
			case "system02":
				ab24Edit.Ab24Identity.Ab24SystemType = 2
			case "address00":
				ab24Edit.Ab24Identity.Ab24Address = 0
			case "system00_address04":
				ab24Edit.Ab24Identity.Ab24SystemType = 0
				ab24Edit.Ab24Identity.Ab24Address = 4
			case "system01_address03":
				ab24Edit.Ab24Identity.Ab24SystemType = 1
				ab24Edit.Ab24Identity.Ab24Address = 3
			case "system03_address03":
				ab24Edit.Ab24Identity.Ab24Address = 3
			}
			ab24Edit.Ab24OP09.Ab24Identity = ab24Edit.Ab24Identity
			ab24Edit.Ab24OP0B.Ab24Identity = ab24Edit.Ab24Identity
			ab24Plan, ab24Err := Ab24BuildEventPlan(ab24Edit)
			if ab24Err == nil || ab24Plan != nil {
				ab24T.Fatal("accepted invalid shape or selector")
			}
		})
	}
	for _, ab24System := range []int{0, 1, 3} {
		ab24Max := 2
		if ab24System == 0 {
			ab24Max = 3
		}
		for ab24Address := 1; ab24Address <= ab24Max; ab24Address++ {
			ab24T.Run(fmt.Sprintf("ab24_valid_%d_%d", ab24System, ab24Address), func(ab24T *testing.T) {
				ab24Edit := ab24Fixture()
				ab24Edit.Ab24Identity.Ab24SystemType = ab24System
				ab24Edit.Ab24Identity.Ab24Address = ab24Address
				ab24Edit.Ab24OP09.Ab24Identity = ab24Edit.Ab24Identity
				ab24Edit.Ab24OP0B.Ab24Identity = ab24Edit.Ab24Identity
				ab24Plan, ab24Err := Ab24BuildEventPlan(ab24Edit)
				if ab24Err != nil || ab24Plan == nil || ab24Plan.Ab24Payload[1] != byte(ab24System) || ab24Plan.Ab24Payload[3] != byte(ab24Address) {
					ab24T.Fatalf("valid selector: %v", ab24Err)
				}
			})
		}
	}
}
