package system

import "fmt"

// EventIdentity names one event record: the device family and the raw
// system type, instance, address and code bytes. The values are ints so that
// out-of-range input can be rejected rather than truncated.
type EventIdentity struct {
	EventPlanFamily     string `json:"family"`
	EventPlanSystemType int    `json:"system_type"`
	EventPlanInstance   int    `json:"instance"`
	EventPlanAddress    int    `json:"address"`
	EventPlanCode       int    `json:"code"`
}

// EventBaseline is a prior read of an event record: its identity, the read
// opcode, whether the read succeeded and the raw reply.
type EventBaseline struct {
	EventPlanIdentity EventIdentity `json:"identity"`
	EventPlanOpcode   int           `json:"opcode"`
	EventPlanSuccess  bool          `json:"success"`
	EventPlanRaw      []byte        `json:"raw"`
}

// EventEdit is a requested change to an event record: its identity, the
// edit opcode, the seven raw values and the baselines read with opcodes
// 0x09 and 0x0B.
type EventEdit struct {
	EventPlanIdentity EventIdentity  `json:"identity"`
	EventPlanOpcode   int            `json:"opcode"`
	EventPlanValues   []int          `json:"values"`
	EventPlanOP09     *EventBaseline `json:"op09"`
	EventPlanOP0B     *EventBaseline `json:"op0b"`
}

// EventPlan is a validated edit: the twelve-byte request payload (opcode,
// system type, instance, address, code and the seven values), the values,
// and both eight-byte baselines it was planned against.
type EventPlan struct {
	EventPlanIdentity EventIdentity `json:"identity"`
	EventPlanPayload  [12]byte      `json:"payload"`
	EventPlanValues   [7]byte       `json:"values"`
	EventPlanOP09     [8]byte       `json:"op09"`
	EventPlanOP0B     [8]byte       `json:"op0b"`
}

// BuildEventPlan validates an edit and builds its plan; it sends nothing.
// It fails when the family is missing, a field or value is outside 0-255,
// the system type is unsupported, the address is outside the system type's
// range, the opcode is not 0x0A or 0x0C, there are not exactly seven
// values, or either baseline is missing, failed, read with the wrong opcode,
// for another identity or not eight bytes long.
func BuildEventPlan(eventPlanEdit EventEdit) (*EventPlan, error) {
	eventPlanIdentity := eventPlanEdit.EventPlanIdentity
	if eventPlanIdentity.EventPlanFamily == "" {
		return nil, fmt.Errorf("event plan: missing family")
	}
	for _, eventPlanField := range []int{eventPlanIdentity.EventPlanSystemType, eventPlanIdentity.EventPlanInstance, eventPlanIdentity.EventPlanAddress, eventPlanIdentity.EventPlanCode, eventPlanEdit.EventPlanOpcode} {
		if eventPlanField < 0 || eventPlanField > 255 {
			return nil, fmt.Errorf("event plan: field outside raw u8 range")
		}
	}
	eventPlanMaxAddress := 2
	switch eventPlanIdentity.EventPlanSystemType {
	case 0:
		eventPlanMaxAddress = 3
	case 1, 3:
	default:
		return nil, fmt.Errorf("event plan: unsupported system type")
	}
	if eventPlanIdentity.EventPlanAddress < 1 || eventPlanIdentity.EventPlanAddress > eventPlanMaxAddress {
		return nil, fmt.Errorf("event plan: invalid address")
	}
	if eventPlanEdit.EventPlanOpcode != 0x0a && eventPlanEdit.EventPlanOpcode != 0x0c {
		return nil, fmt.Errorf("event plan: invalid edit opcode")
	}
	if len(eventPlanEdit.EventPlanValues) != 7 {
		return nil, fmt.Errorf("event plan: exactly seven raw values required")
	}
	for _, eventPlanValue := range eventPlanEdit.EventPlanValues {
		if eventPlanValue < 0 || eventPlanValue > 255 {
			return nil, fmt.Errorf("event plan: value outside raw u8 range")
		}
	}
	for eventPlanIndex, eventPlanBaseline := range []*EventBaseline{eventPlanEdit.EventPlanOP09, eventPlanEdit.EventPlanOP0B} {
		eventPlanOpcode := 0x09 + eventPlanIndex*2
		if eventPlanBaseline == nil || !eventPlanBaseline.EventPlanSuccess || eventPlanBaseline.EventPlanOpcode != eventPlanOpcode || eventPlanBaseline.EventPlanIdentity != eventPlanIdentity || len(eventPlanBaseline.EventPlanRaw) != 8 {
			return nil, fmt.Errorf("event plan: missing, failed or uncorrelated baseline %02x", eventPlanOpcode)
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
