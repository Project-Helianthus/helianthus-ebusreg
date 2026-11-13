package system

import "fmt"

type EventIdentity struct {
	EventPlanFamily     string `json:"Ab24Family"`
	EventPlanSystemType int    `json:"Ab24SystemType"`
	EventPlanInstance   int    `json:"Ab24Instance"`
	EventPlanAddress    int    `json:"Ab24Address"`
	EventPlanCode       int    `json:"Ab24Code"`
}

type EventBaseline struct {
	EventPlanIdentity EventIdentity `json:"Ab24Identity"`
	EventPlanOpcode   int           `json:"Ab24Opcode"`
	EventPlanSuccess  bool          `json:"Ab24Success"`
	EventPlanRaw      []byte        `json:"Ab24Raw"`
}

type EventEdit struct {
	EventPlanIdentity EventIdentity  `json:"Ab24Identity"`
	EventPlanOpcode   int            `json:"Ab24Opcode"`
	EventPlanValues   []int          `json:"Ab24Values"`
	EventPlanOP09     *EventBaseline `json:"Ab24OP09"`
	EventPlanOP0B     *EventBaseline `json:"Ab24OP0B"`
}

type EventPlan struct {
	EventPlanIdentity EventIdentity `json:"Ab24Identity"`
	EventPlanPayload  [12]byte      `json:"Ab24Payload"`
	EventPlanValues   [7]byte       `json:"Ab24Values"`
	EventPlanOP09     [8]byte       `json:"Ab24OP09"`
	EventPlanOP0B     [8]byte       `json:"Ab24OP0B"`
}

func BuildEventPlan(eventPlanEdit EventEdit) (*EventPlan, error) {
	eventPlanIdentity := eventPlanEdit.EventPlanIdentity
	if eventPlanIdentity.EventPlanFamily == "" {
		return nil, fmt.Errorf("ab24: missing family")
	}
	for _, eventPlanField := range []int{eventPlanIdentity.EventPlanSystemType, eventPlanIdentity.EventPlanInstance, eventPlanIdentity.EventPlanAddress, eventPlanIdentity.EventPlanCode, eventPlanEdit.EventPlanOpcode} {
		if eventPlanField < 0 || eventPlanField > 255 {
			return nil, fmt.Errorf("ab24: field outside raw u8 range")
		}
	}
	eventPlanMaxAddress := 2
	switch eventPlanIdentity.EventPlanSystemType {
	case 0:
		eventPlanMaxAddress = 3
	case 1, 3:
	default:
		return nil, fmt.Errorf("ab24: unsupported system type")
	}
	if eventPlanIdentity.EventPlanAddress < 1 || eventPlanIdentity.EventPlanAddress > eventPlanMaxAddress {
		return nil, fmt.Errorf("ab24: invalid address")
	}
	if eventPlanEdit.EventPlanOpcode != 0x0a && eventPlanEdit.EventPlanOpcode != 0x0c {
		return nil, fmt.Errorf("ab24: invalid edit opcode")
	}
	if len(eventPlanEdit.EventPlanValues) != 7 {
		return nil, fmt.Errorf("ab24: exactly seven raw values required")
	}
	for _, eventPlanValue := range eventPlanEdit.EventPlanValues {
		if eventPlanValue < 0 || eventPlanValue > 255 {
			return nil, fmt.Errorf("ab24: value outside raw u8 range")
		}
	}
	for eventPlanIndex, eventPlanBaseline := range []*EventBaseline{eventPlanEdit.EventPlanOP09, eventPlanEdit.EventPlanOP0B} {
		eventPlanOpcode := 0x09 + eventPlanIndex*2
		if eventPlanBaseline == nil || !eventPlanBaseline.EventPlanSuccess || eventPlanBaseline.EventPlanOpcode != eventPlanOpcode || eventPlanBaseline.EventPlanIdentity != eventPlanIdentity || len(eventPlanBaseline.EventPlanRaw) != 8 {
			return nil, fmt.Errorf("ab24: missing, failed or uncorrelated baseline %02x", eventPlanOpcode)
		}
	}
	eventPlanPlan := &EventPlan{EventPlanIdentity: eventPlanIdentity}
	eventPlanPlan.EventPlanPayload = [12]byte{byte(eventPlanEdit.EventPlanOpcode), byte(eventPlanIdentity.EventPlanSystemType), byte(eventPlanIdentity.EventPlanInstance), byte(eventPlanIdentity.EventPlanAddress), byte(eventPlanIdentity.EventPlanCode)}
	for eventPlanIndex, eventPlanValue := range eventPlanEdit.EventPlanValues {
		eventPlanPlan.EventPlanValues[eventPlanIndex] = byte(eventPlanValue)
		eventPlanPlan.EventPlanPayload[5+eventPlanIndex] = byte(eventPlanValue)
	}
	copy(eventPlanPlan.EventPlanOP09[:], eventPlanEdit.EventPlanOP09.EventPlanRaw)
	copy(eventPlanPlan.EventPlanOP0B[:], eventPlanEdit.EventPlanOP0B.EventPlanRaw)
	return eventPlanPlan, nil
}
