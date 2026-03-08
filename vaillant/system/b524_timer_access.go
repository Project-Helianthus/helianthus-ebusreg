package system

import (
	"fmt"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/types"
)

const (
	methodReadTimer = "read_timer"
	methodReadRaw   = "read_raw"

	timerOpcodeRead  = byte(0x03)
	timerOpcodeWrite = byte(0x04)
)

// timerReadTemplate builds a B524 timer read request (opcode 0x03).
//
// Wire format: [0x03, SEL1, SEL2, SEL3, WD]
//
// SEL1/SEL2/SEL3 are timer-specific selectors whose meaning depends on the
// controller model. WD is the weekday index (0x00=Mon .. 0x06=Sun).
type timerReadTemplate struct {
	primary   byte
	secondary byte
}

func (template timerReadTemplate) Primary() byte {
	return template.primary
}

func (template timerReadTemplate) Secondary() byte {
	return template.secondary
}

func (template timerReadTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("timer read template missing params: %w", ebuserrors.ErrInvalidPayload)
	}

	sel1, ok := uint8Param(params, "sel1")
	if !ok {
		return nil, fmt.Errorf("timer read template sel1: %w", ebuserrors.ErrInvalidPayload)
	}
	sel2, ok := uint8Param(params, "sel2")
	if !ok {
		return nil, fmt.Errorf("timer read template sel2: %w", ebuserrors.ErrInvalidPayload)
	}
	sel3, ok := uint8Param(params, "sel3")
	if !ok {
		return nil, fmt.Errorf("timer read template sel3: %w", ebuserrors.ErrInvalidPayload)
	}
	weekday, ok := uint8Param(params, "weekday")
	if !ok {
		return nil, fmt.Errorf("timer read template weekday: %w", ebuserrors.ErrInvalidPayload)
	}
	if weekday > 6 {
		return nil, fmt.Errorf("timer read template weekday %d > 6: %w", weekday, ebuserrors.ErrInvalidPayload)
	}

	return []byte{timerOpcodeRead, sel1, sel2, sel3, weekday}, nil
}

// rawReadTemplate builds a B524 raw opcode request for investigation.
//
// Wire format: caller-provided payload bytes sent verbatim.
// Intended for probing uncharted opcodes (e.g. 0x0B array/table transport)
// without committing to a typed parameter schema.
type rawReadTemplate struct {
	primary   byte
	secondary byte
}

func (template rawReadTemplate) Primary() byte {
	return template.primary
}

func (template rawReadTemplate) Secondary() byte {
	return template.secondary
}

func (template rawReadTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("raw read template missing params: %w", ebuserrors.ErrInvalidPayload)
	}

	payload, hasPayload, err := bytesParam(params, "payload")
	if err != nil {
		return nil, fmt.Errorf("raw read template payload: %w", err)
	}
	if !hasPayload || len(payload) == 0 {
		return nil, fmt.Errorf("raw read template payload empty: %w", ebuserrors.ErrInvalidPayload)
	}

	if len(payload) > 16 {
		return nil, fmt.Errorf("raw read template payload too long (%d > 16): %w", len(payload), ebuserrors.ErrInvalidPayload)
	}

	return append([]byte(nil), payload...), nil
}

func decodeTimerResponse(sel1, sel2, sel3, weekday byte, payload []byte) map[string]types.Value {
	values := map[string]types.Value{
		"opcode":  {Value: timerOpcodeRead, Valid: true},
		"sel1":    {Value: sel1, Valid: true},
		"sel2":    {Value: sel2, Valid: true},
		"sel3":    {Value: sel3, Valid: true},
		"weekday": {Value: weekday, Valid: true},
		"payload": {Value: append([]byte(nil), payload...), Valid: true},
	}

	if len(payload) == 0 {
		values["value"] = types.Value{Valid: false}
		return values
	}

	values["value"] = types.Value{Value: append([]byte(nil), payload...), Valid: true}
	values["slot_count"] = types.Value{Value: len(payload), Valid: true}
	return values
}

func decodeRawResponse(requestPayload []byte, responsePayload []byte) map[string]types.Value {
	values := map[string]types.Value{
		"request_payload":  {Value: append([]byte(nil), requestPayload...), Valid: true},
		"response_payload": {Value: append([]byte(nil), responsePayload...), Valid: true},
	}

	if len(responsePayload) == 0 {
		values["value"] = types.Value{Valid: false}
		return values
	}

	values["value"] = types.Value{Value: append([]byte(nil), responsePayload...), Valid: true}
	return values
}
