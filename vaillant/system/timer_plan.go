package system

import (
	"fmt"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
)

type TimerIdentity struct {
	TimerPlanFamily     string `json:"family"`
	TimerPlanChannel    string `json:"channel"`
	TimerPlanSystemType int    `json:"system_type"`
	TimerPlanInstance   int    `json:"instance"`
	TimerPlanAddress    int    `json:"address"`
	TimerPlanDay        int    `json:"day"`
}

type TimerBaseline struct {
	TimerPlanIdentity TimerIdentity `json:"identity"`
	TimerPlanRaw      []byte        `json:"raw"`
}

type TimerSlot struct {
	TimerPlanStart int `json:"start"`
	TimerPlanStop  int `json:"stop"`
}

type TimerPlan struct {
	timerPlanIdentity TimerIdentity
	timerPlanBaseline [7]byte
	timerPlanPayload  [11]byte
}

func (timerPlanPlan TimerPlan) TimerPlanIdentity() TimerIdentity {
	return timerPlanPlan.timerPlanIdentity
}

func (timerPlanPlan TimerPlan) TimerPlanBaseline() [7]byte {
	return timerPlanPlan.timerPlanBaseline
}

func (timerPlanPlan TimerPlan) TimerPlanPayload() [11]byte {
	return timerPlanPlan.timerPlanPayload
}

func (timerPlanPlan TimerPlan) TimerPlanNativeWriteAdmission() string {
	return "Unknown"
}

func BuildTimerPlan(timerPlanIdentity TimerIdentity, timerPlanBaseline TimerBaseline, timerPlanSlots []TimerSlot) (TimerPlan, error) {
	if timerPlanErr := timerPlanValidateIdentity(timerPlanIdentity); timerPlanErr != nil {
		return TimerPlan{}, timerPlanErr
	}
	if timerPlanIdentity.TimerPlanFamily != timerPlanBaseline.TimerPlanIdentity.TimerPlanFamily ||
		timerPlanIdentity.TimerPlanChannel != timerPlanBaseline.TimerPlanIdentity.TimerPlanChannel ||
		timerPlanIdentity.TimerPlanSystemType != timerPlanBaseline.TimerPlanIdentity.TimerPlanSystemType ||
		timerPlanIdentity.TimerPlanInstance != timerPlanBaseline.TimerPlanIdentity.TimerPlanInstance ||
		timerPlanIdentity.TimerPlanAddress != timerPlanBaseline.TimerPlanIdentity.TimerPlanAddress ||
		timerPlanIdentity.TimerPlanDay != timerPlanBaseline.TimerPlanIdentity.TimerPlanDay {
		return TimerPlan{}, fmt.Errorf("timer plan baseline identity mismatch: %w", ebuserrors.ErrInvalidPayload)
	}
	if len(timerPlanBaseline.TimerPlanRaw) != 7 || len(timerPlanSlots) != 3 {
		return TimerPlan{}, fmt.Errorf("timer plan requires seven baseline bytes and three slots: %w", ebuserrors.ErrInvalidPayload)
	}
	timerPlanPlan := TimerPlan{timerPlanIdentity: timerPlanIdentity}
	copy(timerPlanPlan.timerPlanBaseline[:], timerPlanBaseline.TimerPlanRaw)
	timerPlanPlan.timerPlanPayload = [11]byte{0x04, byte(timerPlanIdentity.TimerPlanSystemType), byte(timerPlanIdentity.TimerPlanInstance), byte(timerPlanIdentity.TimerPlanAddress), byte(timerPlanIdentity.TimerPlanDay)}
	for timerPlanIndex, timerPlanSlot := range timerPlanSlots {
		if !timerPlanByte(timerPlanSlot.TimerPlanStart) || !timerPlanByte(timerPlanSlot.TimerPlanStop) {
			return TimerPlan{}, fmt.Errorf("timer plan slot %d outside byte range: %w", timerPlanIndex, ebuserrors.ErrInvalidPayload)
		}
		if (timerPlanSlot.TimerPlanStart != 0x90 || timerPlanSlot.TimerPlanStop != 0x90) &&
			(timerPlanSlot.TimerPlanStart > 143 || timerPlanSlot.TimerPlanStart >= timerPlanSlot.TimerPlanStop || timerPlanSlot.TimerPlanStop > 144) {
			return TimerPlan{}, fmt.Errorf("timer plan slot %d invalid raw ten-minute pair: %w", timerPlanIndex, ebuserrors.ErrInvalidPayload)
		}
		timerPlanPlan.timerPlanPayload[5+2*timerPlanIndex] = byte(timerPlanSlot.TimerPlanStart)
		timerPlanPlan.timerPlanPayload[6+2*timerPlanIndex] = byte(timerPlanSlot.TimerPlanStop)
	}
	return timerPlanPlan, nil
}

func timerPlanByte(timerPlanValue int) bool {
	return timerPlanValue >= 0 && timerPlanValue <= 255
}

func timerPlanValidateIdentity(timerPlanIdentity TimerIdentity) error {
	if timerPlanIdentity.TimerPlanFamily == "" || timerPlanIdentity.TimerPlanChannel == "" ||
		!timerPlanByte(timerPlanIdentity.TimerPlanSystemType) || !timerPlanByte(timerPlanIdentity.TimerPlanInstance) ||
		!timerPlanByte(timerPlanIdentity.TimerPlanAddress) || !timerPlanByte(timerPlanIdentity.TimerPlanDay) || timerPlanIdentity.TimerPlanDay > 6 {
		return fmt.Errorf("timer plan invalid explicit identity: %w", ebuserrors.ErrInvalidPayload)
	}
	timerPlanKnown := false
	switch timerPlanIdentity.TimerPlanSystemType {
	case 0x00:
		timerPlanKnown = timerPlanIdentity.TimerPlanAddress >= 1 && timerPlanIdentity.TimerPlanAddress <= 3
	case 0x01, 0x03:
		timerPlanKnown = timerPlanIdentity.TimerPlanAddress == 1 || timerPlanIdentity.TimerPlanAddress == 2
	}
	if !timerPlanKnown || (timerPlanIdentity.TimerPlanSystemType != 0x03 && timerPlanIdentity.TimerPlanInstance != 0) {
		return fmt.Errorf("timer plan unsupported raw selector or non-zone instance: %w", ebuserrors.ErrInvalidPayload)
	}
	return nil
}
