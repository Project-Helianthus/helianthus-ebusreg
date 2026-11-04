package system

import (
	"fmt"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
)

type Ab08TimerIdentity struct {
	Ab08Family     string
	Ab08Channel    string
	Ab08SystemType int
	Ab08Instance   int
	Ab08Address    int
	Ab08Day        int
}

type Ab08TimerBaseline struct {
	Ab08Identity Ab08TimerIdentity
	Ab08Raw      []byte
}

type Ab08TimerSlot struct {
	Ab08Start int
	Ab08Stop  int
}

type Ab08TimerPlan struct {
	ab08Identity Ab08TimerIdentity
	ab08Baseline [7]byte
	ab08Payload  [11]byte
}

func (ab08Plan Ab08TimerPlan) Ab08Identity() Ab08TimerIdentity {
	return ab08Plan.ab08Identity
}

func (ab08Plan Ab08TimerPlan) Ab08Baseline() [7]byte {
	return ab08Plan.ab08Baseline
}

func (ab08Plan Ab08TimerPlan) Ab08Payload() [11]byte {
	return ab08Plan.ab08Payload
}

func (ab08Plan Ab08TimerPlan) Ab08NativeWriteAdmission() string {
	return "Unknown"
}

func Ab08BuildTimerPlan(ab08Identity Ab08TimerIdentity, ab08Baseline Ab08TimerBaseline, ab08Slots []Ab08TimerSlot) (Ab08TimerPlan, error) {
	if ab08Err := ab08ValidateIdentity(ab08Identity); ab08Err != nil {
		return Ab08TimerPlan{}, ab08Err
	}
	if ab08Identity.Ab08Family != ab08Baseline.Ab08Identity.Ab08Family ||
		ab08Identity.Ab08Channel != ab08Baseline.Ab08Identity.Ab08Channel ||
		ab08Identity.Ab08SystemType != ab08Baseline.Ab08Identity.Ab08SystemType ||
		ab08Identity.Ab08Instance != ab08Baseline.Ab08Identity.Ab08Instance ||
		ab08Identity.Ab08Address != ab08Baseline.Ab08Identity.Ab08Address ||
		ab08Identity.Ab08Day != ab08Baseline.Ab08Identity.Ab08Day {
		return Ab08TimerPlan{}, fmt.Errorf("ab08 baseline identity mismatch: %w", ebuserrors.ErrInvalidPayload)
	}
	if len(ab08Baseline.Ab08Raw) != 7 || len(ab08Slots) != 3 {
		return Ab08TimerPlan{}, fmt.Errorf("ab08 requires seven baseline bytes and three slots: %w", ebuserrors.ErrInvalidPayload)
	}
	ab08Plan := Ab08TimerPlan{ab08Identity: ab08Identity}
	copy(ab08Plan.ab08Baseline[:], ab08Baseline.Ab08Raw)
	ab08Plan.ab08Payload = [11]byte{0x04, byte(ab08Identity.Ab08SystemType), byte(ab08Identity.Ab08Instance), byte(ab08Identity.Ab08Address), byte(ab08Identity.Ab08Day)}
	for ab08Index, ab08Slot := range ab08Slots {
		if !ab08Byte(ab08Slot.Ab08Start) || !ab08Byte(ab08Slot.Ab08Stop) {
			return Ab08TimerPlan{}, fmt.Errorf("ab08 slot %d outside byte range: %w", ab08Index, ebuserrors.ErrInvalidPayload)
		}
		if (ab08Slot.Ab08Start != 0x90 || ab08Slot.Ab08Stop != 0x90) &&
			(ab08Slot.Ab08Start > 143 || ab08Slot.Ab08Start >= ab08Slot.Ab08Stop || ab08Slot.Ab08Stop > 144) {
			return Ab08TimerPlan{}, fmt.Errorf("ab08 slot %d invalid raw ten-minute pair: %w", ab08Index, ebuserrors.ErrInvalidPayload)
		}
		ab08Plan.ab08Payload[5+2*ab08Index] = byte(ab08Slot.Ab08Start)
		ab08Plan.ab08Payload[6+2*ab08Index] = byte(ab08Slot.Ab08Stop)
	}
	return ab08Plan, nil
}

func ab08Byte(ab08Value int) bool {
	return ab08Value >= 0 && ab08Value <= 255
}

func ab08ValidateIdentity(ab08Identity Ab08TimerIdentity) error {
	if ab08Identity.Ab08Family == "" || ab08Identity.Ab08Channel == "" ||
		!ab08Byte(ab08Identity.Ab08SystemType) || !ab08Byte(ab08Identity.Ab08Instance) ||
		!ab08Byte(ab08Identity.Ab08Address) || !ab08Byte(ab08Identity.Ab08Day) || ab08Identity.Ab08Day > 6 {
		return fmt.Errorf("ab08 invalid explicit identity: %w", ebuserrors.ErrInvalidPayload)
	}
	ab08Known := false
	switch ab08Identity.Ab08SystemType {
	case 0x00:
		ab08Known = ab08Identity.Ab08Address >= 1 && ab08Identity.Ab08Address <= 3
	case 0x01, 0x03:
		ab08Known = ab08Identity.Ab08Address == 1 || ab08Identity.Ab08Address == 2
	}
	if !ab08Known || (ab08Identity.Ab08SystemType != 0x03 && ab08Identity.Ab08Instance != 0) {
		return fmt.Errorf("ab08 unsupported raw selector or non-zone instance: %w", ebuserrors.ErrInvalidPayload)
	}
	return nil
}
