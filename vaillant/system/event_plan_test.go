package system

import (
	"fmt"
	"reflect"
	"testing"
)

func eventPlanFixture() EventEdit {
	eventPlanIdentity := EventIdentity{"ab24-synthetic-family", 3, 2, 1, 255}
	return EventEdit{
		EventPlanIdentity: eventPlanIdentity, EventPlanOpcode: 0x0a,
		EventPlanValues: []int{0, 1, 2, 3, 4, 5, 6},
		EventPlanOP09:   &EventBaseline{eventPlanIdentity, 0x09, true, []byte{0, 1, 2, 3, 4, 5, 6, 7}},
		EventPlanOP0B:   &EventBaseline{eventPlanIdentity, 0x0b, true, []byte{8, 9, 10, 11, 12, 13, 14, 15}},
	}
}

func TestAb24ExactVectorsAndOwnership(eventPlanT *testing.T) {
	for _, eventPlanOpcode := range []int{0x0a, 0x0c} {
		eventPlanT.Run(fmt.Sprintf("ab24_opcode_%02x", eventPlanOpcode), func(eventPlanT *testing.T) {
			eventPlanEdit := eventPlanFixture()
			eventPlanEdit.EventPlanOpcode = eventPlanOpcode
			eventPlanPlan, eventPlanErr := BuildEventPlan(eventPlanEdit)
			eventPlanWant := [12]byte{byte(eventPlanOpcode), 3, 2, 1, 255, 0, 1, 2, 3, 4, 5, 6}
			if eventPlanErr != nil || eventPlanPlan == nil {
				eventPlanT.Fatalf("plan: %v", eventPlanErr)
			}
			if eventPlanPlan.EventPlanPayload != eventPlanWant || eventPlanPlan.EventPlanIdentity != eventPlanEdit.EventPlanIdentity || eventPlanPlan.EventPlanValues != [7]byte{0, 1, 2, 3, 4, 5, 6} || eventPlanPlan.EventPlanOP09 != [8]byte{0, 1, 2, 3, 4, 5, 6, 7} || eventPlanPlan.EventPlanOP0B != [8]byte{8, 9, 10, 11, 12, 13, 14, 15} {
				eventPlanT.Fatalf("unexpected plan: %+v", eventPlanPlan)
			}
			eventPlanBefore := *eventPlanPlan
			eventPlanEdit.EventPlanValues[0] = 99
			eventPlanEdit.EventPlanOP09.EventPlanRaw[0] = 99
			eventPlanEdit.EventPlanOP0B.EventPlanRaw[0] = 99
			eventPlanEdit.EventPlanIdentity.EventPlanFamily = "ab24-other-family"
			if *eventPlanPlan != eventPlanBefore {
				eventPlanT.Fatal("input mutation changed owned plan")
			}
		})
	}
}

func TestAb24IdentityAndRawValues(eventPlanT *testing.T) {
	for eventPlanField, eventPlanValue := range map[string]any{"EventPlanFamily": "ab24-other-family", "EventPlanSystemType": 0, "EventPlanInstance": 255, "EventPlanAddress": 2, "EventPlanCode": 0} {
		eventPlanT.Run("ab24_"+eventPlanField, func(eventPlanT *testing.T) {
			eventPlanEdit := eventPlanFixture()
			reflect.ValueOf(&eventPlanEdit.EventPlanIdentity).Elem().FieldByName(eventPlanField).Set(reflect.ValueOf(eventPlanValue))
			eventPlanEdit.EventPlanOP09.EventPlanIdentity = eventPlanEdit.EventPlanIdentity
			eventPlanEdit.EventPlanOP0B.EventPlanIdentity = eventPlanEdit.EventPlanIdentity
			eventPlanPlan, eventPlanErr := BuildEventPlan(eventPlanEdit)
			if eventPlanErr != nil || eventPlanPlan == nil {
				eventPlanT.Fatalf("valid identity: %v", eventPlanErr)
			}
			eventPlanWant := [12]byte{0x0a, byte(eventPlanEdit.EventPlanIdentity.EventPlanSystemType), byte(eventPlanEdit.EventPlanIdentity.EventPlanInstance), byte(eventPlanEdit.EventPlanIdentity.EventPlanAddress), byte(eventPlanEdit.EventPlanIdentity.EventPlanCode), 0, 1, 2, 3, 4, 5, 6}
			if eventPlanPlan.EventPlanIdentity != eventPlanEdit.EventPlanIdentity || eventPlanPlan.EventPlanPayload != eventPlanWant {
				eventPlanT.Fatalf("identity lost: %+v", eventPlanPlan)
			}
		})
	}
	for eventPlanIndex := 0; eventPlanIndex < 7; eventPlanIndex++ {
		eventPlanT.Run(fmt.Sprintf("ab24_value_%d", eventPlanIndex), func(eventPlanT *testing.T) {
			eventPlanEdit := eventPlanFixture()
			eventPlanEdit.EventPlanValues[eventPlanIndex] = 255
			eventPlanPlan, eventPlanErr := BuildEventPlan(eventPlanEdit)
			if eventPlanErr != nil || eventPlanPlan == nil {
				eventPlanT.Fatalf("raw value: %v", eventPlanErr)
			}
			eventPlanWant := [12]byte{0x0a, 3, 2, 1, 255, 0, 1, 2, 3, 4, 5, 6}
			eventPlanWant[5+eventPlanIndex] = 255
			eventPlanWantValues := [7]byte{0, 1, 2, 3, 4, 5, 6}
			eventPlanWantValues[eventPlanIndex] = 255
			if eventPlanPlan.EventPlanPayload != eventPlanWant || eventPlanPlan.EventPlanValues != eventPlanWantValues {
				eventPlanT.Fatalf("value order: %+v", eventPlanPlan)
			}
		})
	}
}

func TestAb24RejectBaseline(eventPlanT *testing.T) {
	for _, eventPlanSlot := range []string{"ab24_op09", "ab24_op0b"} {
		for _, eventPlanCase := range []string{"missing", "failed", "opcode", "short", "long", "family", "system", "instance", "address", "code"} {
			eventPlanT.Run(eventPlanSlot+"_"+eventPlanCase, func(eventPlanT *testing.T) {
				eventPlanEdit := eventPlanFixture()
				eventPlanBaseline := &eventPlanEdit.EventPlanOP09
				if eventPlanSlot == "ab24_op0b" {
					eventPlanBaseline = &eventPlanEdit.EventPlanOP0B
				}
				switch eventPlanCase {
				case "missing":
					*eventPlanBaseline = nil
				case "failed":
					(*eventPlanBaseline).EventPlanSuccess = false
				case "opcode":
					(*eventPlanBaseline).EventPlanOpcode = 5
				case "short":
					(*eventPlanBaseline).EventPlanRaw = (*eventPlanBaseline).EventPlanRaw[:7]
				case "long":
					(*eventPlanBaseline).EventPlanRaw = append((*eventPlanBaseline).EventPlanRaw, 9)
				case "family":
					(*eventPlanBaseline).EventPlanIdentity.EventPlanFamily = "ab24-other-family"
				case "system":
					(*eventPlanBaseline).EventPlanIdentity.EventPlanSystemType = 1
				case "instance":
					(*eventPlanBaseline).EventPlanIdentity.EventPlanInstance = 3
				case "address":
					(*eventPlanBaseline).EventPlanIdentity.EventPlanAddress = 2
				case "code":
					(*eventPlanBaseline).EventPlanIdentity.EventPlanCode = 254
				}
				eventPlanPlan, eventPlanErr := BuildEventPlan(eventPlanEdit)
				if eventPlanErr == nil || eventPlanPlan != nil {
					eventPlanT.Fatalf("accepted invalid baseline: %+v %v", eventPlanPlan, eventPlanErr)
				}
			})
		}
	}
}

func TestAb24RejectFields(eventPlanT *testing.T) {
	for _, eventPlanField := range []string{"EventPlanSystemType", "EventPlanInstance", "EventPlanAddress", "EventPlanCode", "EventPlanOpcode"} {
		for _, eventPlanValue := range []int{-1, 256} {
			eventPlanT.Run(fmt.Sprintf("ab24_%s_%d", eventPlanField, eventPlanValue), func(eventPlanT *testing.T) {
				eventPlanEdit := eventPlanFixture()
				if eventPlanField == "EventPlanOpcode" {
					eventPlanEdit.EventPlanOpcode = eventPlanValue
				} else {
					reflect.ValueOf(&eventPlanEdit.EventPlanIdentity).Elem().FieldByName(eventPlanField).SetInt(int64(eventPlanValue))
					eventPlanEdit.EventPlanOP09.EventPlanIdentity = eventPlanEdit.EventPlanIdentity
					eventPlanEdit.EventPlanOP0B.EventPlanIdentity = eventPlanEdit.EventPlanIdentity
				}
				eventPlanPlan, eventPlanErr := BuildEventPlan(eventPlanEdit)
				if eventPlanErr == nil || eventPlanPlan != nil {
					eventPlanT.Fatal("accepted narrowed field")
				}
			})
		}
	}
	for eventPlanIndex := 0; eventPlanIndex < 7; eventPlanIndex++ {
		for _, eventPlanValue := range []int{-1, 256} {
			eventPlanT.Run(fmt.Sprintf("ab24_value_%d_%d", eventPlanIndex, eventPlanValue), func(eventPlanT *testing.T) {
				eventPlanEdit := eventPlanFixture()
				eventPlanEdit.EventPlanValues[eventPlanIndex] = eventPlanValue
				eventPlanPlan, eventPlanErr := BuildEventPlan(eventPlanEdit)
				if eventPlanErr == nil || eventPlanPlan != nil {
					eventPlanT.Fatal("accepted narrowed value")
				}
			})
		}
	}
}

func TestAb24ShapeAndSelectors(eventPlanT *testing.T) {
	for _, eventPlanCase := range []string{"six", "eight", "op05", "empty_family", "system02", "address00", "system00_address04", "system01_address03", "system03_address03"} {
		eventPlanT.Run("ab24_"+eventPlanCase, func(eventPlanT *testing.T) {
			eventPlanEdit := eventPlanFixture()
			switch eventPlanCase {
			case "six":
				eventPlanEdit.EventPlanValues = eventPlanEdit.EventPlanValues[:6]
			case "eight":
				eventPlanEdit.EventPlanValues = append(eventPlanEdit.EventPlanValues, 7)
			case "op05":
				eventPlanEdit.EventPlanOpcode = 5
			case "empty_family":
				eventPlanEdit.EventPlanIdentity.EventPlanFamily = ""
			case "system02":
				eventPlanEdit.EventPlanIdentity.EventPlanSystemType = 2
			case "address00":
				eventPlanEdit.EventPlanIdentity.EventPlanAddress = 0
			case "system00_address04":
				eventPlanEdit.EventPlanIdentity.EventPlanSystemType = 0
				eventPlanEdit.EventPlanIdentity.EventPlanAddress = 4
			case "system01_address03":
				eventPlanEdit.EventPlanIdentity.EventPlanSystemType = 1
				eventPlanEdit.EventPlanIdentity.EventPlanAddress = 3
			case "system03_address03":
				eventPlanEdit.EventPlanIdentity.EventPlanAddress = 3
			}
			eventPlanEdit.EventPlanOP09.EventPlanIdentity = eventPlanEdit.EventPlanIdentity
			eventPlanEdit.EventPlanOP0B.EventPlanIdentity = eventPlanEdit.EventPlanIdentity
			eventPlanPlan, eventPlanErr := BuildEventPlan(eventPlanEdit)
			if eventPlanErr == nil || eventPlanPlan != nil {
				eventPlanT.Fatal("accepted invalid shape or selector")
			}
		})
	}
	for _, eventPlanSystem := range []int{0, 1, 3} {
		eventPlanMax := 2
		if eventPlanSystem == 0 {
			eventPlanMax = 3
		}
		for eventPlanAddress := 1; eventPlanAddress <= eventPlanMax; eventPlanAddress++ {
			eventPlanT.Run(fmt.Sprintf("ab24_valid_%d_%d", eventPlanSystem, eventPlanAddress), func(eventPlanT *testing.T) {
				eventPlanEdit := eventPlanFixture()
				eventPlanEdit.EventPlanIdentity.EventPlanSystemType = eventPlanSystem
				eventPlanEdit.EventPlanIdentity.EventPlanAddress = eventPlanAddress
				eventPlanEdit.EventPlanOP09.EventPlanIdentity = eventPlanEdit.EventPlanIdentity
				eventPlanEdit.EventPlanOP0B.EventPlanIdentity = eventPlanEdit.EventPlanIdentity
				eventPlanPlan, eventPlanErr := BuildEventPlan(eventPlanEdit)
				if eventPlanErr != nil || eventPlanPlan == nil || eventPlanPlan.EventPlanPayload[1] != byte(eventPlanSystem) || eventPlanPlan.EventPlanPayload[3] != byte(eventPlanAddress) {
					eventPlanT.Fatalf("valid selector: %v", eventPlanErr)
				}
			})
		}
	}
}
