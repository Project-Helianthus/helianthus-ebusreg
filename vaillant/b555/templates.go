package b555

import (
	"fmt"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
)

const (
	opcodeConfigRead = byte(0xA3)
	opcodeSlotsRead  = byte(0xA4)
	opcodeTimerRead  = byte(0xA5)
	opcodeTimerWrite = byte(0xA6)
)

// configReadTemplate builds a B555 CONFIG_READ (A3) request.
//
// Wire format: [A3, ZONE, HC]
type configReadTemplate struct {
	primary   byte
	secondary byte
}

func (t configReadTemplate) Primary() byte   { return t.primary }
func (t configReadTemplate) Secondary() byte { return t.secondary }

func (t configReadTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("b555 config read missing params: %w", ebuserrors.ErrInvalidPayload)
	}
	zone, ok := uint8Param(params, "zone")
	if !ok {
		return nil, fmt.Errorf("b555 config read zone: %w", ebuserrors.ErrInvalidPayload)
	}
	hc, ok := uint8Param(params, "hc")
	if !ok {
		return nil, fmt.Errorf("b555 config read hc: %w", ebuserrors.ErrInvalidPayload)
	}
	if hc > 0x04 {
		return nil, fmt.Errorf("b555 config read hc %d > 4: %w", hc, ebuserrors.ErrInvalidPayload)
	}
	return []byte{opcodeConfigRead, zone, hc}, nil
}

// slotsReadTemplate builds a B555 SLOTS_READ (A4) request.
//
// Wire format: [A4, ZONE, HC]
type slotsReadTemplate struct {
	primary   byte
	secondary byte
}

func (t slotsReadTemplate) Primary() byte   { return t.primary }
func (t slotsReadTemplate) Secondary() byte { return t.secondary }

func (t slotsReadTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("b555 slots read missing params: %w", ebuserrors.ErrInvalidPayload)
	}
	zone, ok := uint8Param(params, "zone")
	if !ok {
		return nil, fmt.Errorf("b555 slots read zone: %w", ebuserrors.ErrInvalidPayload)
	}
	hc, ok := uint8Param(params, "hc")
	if !ok {
		return nil, fmt.Errorf("b555 slots read hc: %w", ebuserrors.ErrInvalidPayload)
	}
	if hc > 0x04 {
		return nil, fmt.Errorf("b555 slots read hc %d > 4: %w", hc, ebuserrors.ErrInvalidPayload)
	}
	return []byte{opcodeSlotsRead, zone, hc}, nil
}

// timerReadTemplate builds a B555 TIMER_READ (A5) request.
//
// Wire format: [A5, ZONE, HC, DD, SS]
type timerReadTemplate struct {
	primary   byte
	secondary byte
}

func (t timerReadTemplate) Primary() byte   { return t.primary }
func (t timerReadTemplate) Secondary() byte { return t.secondary }

func (t timerReadTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("b555 timer read missing params: %w", ebuserrors.ErrInvalidPayload)
	}
	zone, ok := uint8Param(params, "zone")
	if !ok {
		return nil, fmt.Errorf("b555 timer read zone: %w", ebuserrors.ErrInvalidPayload)
	}
	hc, ok := uint8Param(params, "hc")
	if !ok {
		return nil, fmt.Errorf("b555 timer read hc: %w", ebuserrors.ErrInvalidPayload)
	}
	if hc > 0x04 {
		return nil, fmt.Errorf("b555 timer read hc %d > 4: %w", hc, ebuserrors.ErrInvalidPayload)
	}
	dd, ok := uint8Param(params, "weekday")
	if !ok {
		return nil, fmt.Errorf("b555 timer read weekday: %w", ebuserrors.ErrInvalidPayload)
	}
	if dd > 0x06 {
		return nil, fmt.Errorf("b555 timer read weekday %d > 6: %w", dd, ebuserrors.ErrInvalidPayload)
	}
	ss, ok := uint8Param(params, "slot")
	if !ok {
		return nil, fmt.Errorf("b555 timer read slot: %w", ebuserrors.ErrInvalidPayload)
	}
	return []byte{opcodeTimerRead, zone, hc, dd, ss}, nil
}

// timerWriteTemplate builds a B555 TIMER_WRITE (A6) request.
//
// Wire format: [A6, ZONE, HC, DD, SI, SC, Sh, Sm, Eh, Em, Tlo, Thi]
type timerWriteTemplate struct {
	primary   byte
	secondary byte
}

func (t timerWriteTemplate) Primary() byte   { return t.primary }
func (t timerWriteTemplate) Secondary() byte { return t.secondary }

func (t timerWriteTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("b555 timer write missing params: %w", ebuserrors.ErrInvalidPayload)
	}
	zone, ok := uint8Param(params, "zone")
	if !ok {
		return nil, fmt.Errorf("b555 timer write zone: %w", ebuserrors.ErrInvalidPayload)
	}
	hc, ok := uint8Param(params, "hc")
	if !ok {
		return nil, fmt.Errorf("b555 timer write hc: %w", ebuserrors.ErrInvalidPayload)
	}
	if hc > 0x04 {
		return nil, fmt.Errorf("b555 timer write hc %d > 4: %w", hc, ebuserrors.ErrInvalidPayload)
	}
	dd, ok := uint8Param(params, "weekday")
	if !ok {
		return nil, fmt.Errorf("b555 timer write weekday: %w", ebuserrors.ErrInvalidPayload)
	}
	if dd > 0x06 {
		return nil, fmt.Errorf("b555 timer write weekday %d > 6: %w", dd, ebuserrors.ErrInvalidPayload)
	}
	si, ok := uint8Param(params, "slot")
	if !ok {
		return nil, fmt.Errorf("b555 timer write slot: %w", ebuserrors.ErrInvalidPayload)
	}
	sc, ok := uint8Param(params, "slot_count")
	if !ok {
		return nil, fmt.Errorf("b555 timer write slot_count: %w", ebuserrors.ErrInvalidPayload)
	}
	if sc == 0 {
		return nil, fmt.Errorf("b555 timer write slot_count must be > 0: %w", ebuserrors.ErrInvalidPayload)
	}
	startHour, ok := uint8Param(params, "start_hour")
	if !ok {
		return nil, fmt.Errorf("b555 timer write start_hour: %w", ebuserrors.ErrInvalidPayload)
	}
	if startHour > 24 {
		return nil, fmt.Errorf("b555 timer write start_hour %d > 24: %w", startHour, ebuserrors.ErrInvalidPayload)
	}
	startMinute, ok := uint8Param(params, "start_minute")
	if !ok {
		return nil, fmt.Errorf("b555 timer write start_minute: %w", ebuserrors.ErrInvalidPayload)
	}
	if startMinute > 59 {
		return nil, fmt.Errorf("b555 timer write start_minute %d > 59: %w", startMinute, ebuserrors.ErrInvalidPayload)
	}
	endHour, ok := uint8Param(params, "end_hour")
	if !ok {
		return nil, fmt.Errorf("b555 timer write end_hour: %w", ebuserrors.ErrInvalidPayload)
	}
	if endHour > 24 {
		return nil, fmt.Errorf("b555 timer write end_hour %d > 24: %w", endHour, ebuserrors.ErrInvalidPayload)
	}
	endMinute, ok := uint8Param(params, "end_minute")
	if !ok {
		return nil, fmt.Errorf("b555 timer write end_minute: %w", ebuserrors.ErrInvalidPayload)
	}
	if endMinute > 59 {
		return nil, fmt.Errorf("b555 timer write end_minute %d > 59: %w", endMinute, ebuserrors.ErrInvalidPayload)
	}
	temp, ok := uint16Param(params, "temperature")
	if !ok {
		return nil, fmt.Errorf("b555 timer write temperature: %w", ebuserrors.ErrInvalidPayload)
	}

	return []byte{
		opcodeTimerWrite,
		zone, hc, dd, si, sc,
		startHour, startMinute,
		endHour, endMinute,
		byte(temp), byte(temp >> 8),
	}, nil
}
