package system

import (
	"errors"
	"testing"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
)

func timerPlanFixture() (TimerIdentity, TimerBaseline, []TimerSlot) {
	timerPlanIdentity := TimerIdentity{"synthetic-family", "raw-channel", 3, 2, 2, 6}
	return timerPlanIdentity, TimerBaseline{timerPlanIdentity, []byte{6, 0x90, 0x90, 0x90, 0x90, 0x90, 0x90}},
		[]TimerSlot{{6, 12}, {0x90, 0x90}, {0x8f, 0x90}}
}

func TestAb08ExactPayloadAndOwnership(timerPlanT *testing.T) {
	timerPlanIdentity, timerPlanBaseline, timerPlanSlots := timerPlanFixture()
	timerPlanPlan, timerPlanErr := BuildTimerPlan(timerPlanIdentity, timerPlanBaseline, timerPlanSlots)
	if timerPlanErr != nil {
		timerPlanT.Fatal(timerPlanErr)
	}
	timerPlanWant := [11]byte{0x04, 0x03, 0x02, 0x02, 0x06, 0x06, 0x0c, 0x90, 0x90, 0x8f, 0x90}
	if timerPlanPlan.TimerPlanPayload() != timerPlanWant {
		timerPlanT.Fatalf("payload = % X; want % X", timerPlanPlan.TimerPlanPayload(), timerPlanWant)
	}
	timerPlanRaw := timerPlanPlan.TimerPlanBaseline()
	if timerPlanRaw != [7]byte{6, 0x90, 0x90, 0x90, 0x90, 0x90, 0x90} || timerPlanPlan.TimerPlanIdentity() != timerPlanIdentity {
		timerPlanT.Fatal("baseline or identity not preserved")
	}
	timerPlanBaseline.TimerPlanRaw[0] = 99
	timerPlanBaseline.TimerPlanIdentity.TimerPlanFamily = "changed-family"
	timerPlanSlots[0] = TimerSlot{0, 1}
	timerPlanIdentity.TimerPlanDay = 0
	timerPlanReturnedIdentity := timerPlanPlan.TimerPlanIdentity()
	timerPlanReturnedIdentity.TimerPlanDay = 1
	timerPlanRaw[0] = 100
	timerPlanReturnedPayload := timerPlanPlan.TimerPlanPayload()
	timerPlanReturnedPayload[0] = 0
	if timerPlanPlan.TimerPlanBaseline()[0] != 6 || timerPlanPlan.TimerPlanPayload() != timerPlanWant || timerPlanPlan.TimerPlanIdentity().TimerPlanDay != 6 || timerPlanPlan.TimerPlanIdentity().TimerPlanFamily != "synthetic-family" {
		timerPlanT.Fatal("plan aliases caller or accessor memory")
	}
	if timerPlanPlan.TimerPlanNativeWriteAdmission() != "Unknown" {
		timerPlanT.Fatal("offline planning must not promote native write admission")
	}
}

func TestAb08IdentityMismatch(timerPlanT *testing.T) {
	timerPlanCases := []struct {
		timerPlanName   string
		timerPlanChange func(*TimerIdentity)
	}{
		{"ab08Family", func(timerPlanID *TimerIdentity) { timerPlanID.TimerPlanFamily = "other-family" }},
		{"ab08Channel", func(timerPlanID *TimerIdentity) { timerPlanID.TimerPlanChannel = "other-channel" }},
		{"ab08SystemType", func(timerPlanID *TimerIdentity) { timerPlanID.TimerPlanSystemType = 1 }},
		{"ab08Instance", func(timerPlanID *TimerIdentity) { timerPlanID.TimerPlanInstance = 255 }},
		{"ab08Address", func(timerPlanID *TimerIdentity) { timerPlanID.TimerPlanAddress = 1 }},
		{"ab08Day", func(timerPlanID *TimerIdentity) { timerPlanID.TimerPlanDay = 5 }},
	}
	for _, timerPlanCase := range timerPlanCases {
		timerPlanT.Run(timerPlanCase.timerPlanName, func(timerPlanT *testing.T) {
			timerPlanID, timerPlanBaseline, timerPlanSlots := timerPlanFixture()
			timerPlanID.TimerPlanInstance = 0
			timerPlanBaseline.TimerPlanIdentity = timerPlanID
			if _, timerPlanErr := BuildTimerPlan(timerPlanID, timerPlanBaseline, timerPlanSlots); timerPlanErr != nil {
				timerPlanT.Fatal("valid control:", timerPlanErr)
			}
			timerPlanCase.timerPlanChange(&timerPlanID)
			timerPlanExpectInvalid(timerPlanT, timerPlanID, timerPlanBaseline, timerPlanSlots)
			timerPlanBaseline.TimerPlanIdentity = timerPlanID
			if _, timerPlanErr := BuildTimerPlan(timerPlanID, timerPlanBaseline, timerPlanSlots); timerPlanErr != nil {
				timerPlanT.Fatal("independently changed valid control:", timerPlanErr)
			}
		})
	}
}

func timerPlanExpectInvalid(timerPlanT *testing.T, timerPlanID TimerIdentity, timerPlanBaseline TimerBaseline, timerPlanSlots []TimerSlot) {
	timerPlanT.Helper()
	timerPlanPlan, timerPlanErr := BuildTimerPlan(timerPlanID, timerPlanBaseline, timerPlanSlots)
	if !errors.Is(timerPlanErr, ebuserrors.ErrInvalidPayload) || timerPlanPlan != (TimerPlan{}) {
		timerPlanT.Fatalf("expected invalid payload and empty plan; got %v, %+v", timerPlanErr, timerPlanPlan)
	}
}

func TestAb08InvalidWidthsAndSlots(timerPlanT *testing.T) {
	for _, timerPlanWidth := range []int{0, 6, 8} {
		timerPlanID, timerPlanBaseline, timerPlanSlots := timerPlanFixture()
		timerPlanBaseline.TimerPlanRaw = make([]byte, timerPlanWidth)
		timerPlanExpectInvalid(timerPlanT, timerPlanID, timerPlanBaseline, timerPlanSlots)
	}
	for _, timerPlanCount := range []int{0, 2, 4} {
		timerPlanID, timerPlanBaseline, _ := timerPlanFixture()
		timerPlanSlots := make([]TimerSlot, timerPlanCount)
		for timerPlanIndex := range timerPlanSlots {
			timerPlanSlots[timerPlanIndex] = TimerSlot{0, 1}
		}
		timerPlanExpectInvalid(timerPlanT, timerPlanID, timerPlanBaseline, timerPlanSlots)
	}
	for _, timerPlanBad := range []TimerSlot{{144, 145}, {143, 145}, {12, 6}, {6, 6}, {144, 143}, {0, 0}, {-1, 12}, {256, 12}, {6, -1}, {6, 256}, {0x90, 0x8f}} {
		for timerPlanIndex := 0; timerPlanIndex < 3; timerPlanIndex++ {
			timerPlanID, timerPlanBaseline, timerPlanSlots := timerPlanFixture()
			timerPlanSlots[timerPlanIndex] = timerPlanBad
			timerPlanExpectInvalid(timerPlanT, timerPlanID, timerPlanBaseline, timerPlanSlots)
		}
	}
}

func TestAb08IdentityBounds(timerPlanT *testing.T) {
	timerPlanChanges := []func(*TimerIdentity, int){
		func(timerPlanID *TimerIdentity, timerPlanValue int) { timerPlanID.TimerPlanSystemType = timerPlanValue },
		func(timerPlanID *TimerIdentity, timerPlanValue int) { timerPlanID.TimerPlanInstance = timerPlanValue },
		func(timerPlanID *TimerIdentity, timerPlanValue int) { timerPlanID.TimerPlanAddress = timerPlanValue },
		func(timerPlanID *TimerIdentity, timerPlanValue int) { timerPlanID.TimerPlanDay = timerPlanValue },
	}
	for _, timerPlanChange := range timerPlanChanges {
		for _, timerPlanValue := range []int{-1, 256} {
			timerPlanID, timerPlanBaseline, timerPlanSlots := timerPlanFixture()
			timerPlanChange(&timerPlanID, timerPlanValue)
			timerPlanBaseline.TimerPlanIdentity = timerPlanID
			timerPlanExpectInvalid(timerPlanT, timerPlanID, timerPlanBaseline, timerPlanSlots)
		}
	}
	for _, timerPlanChange := range []func(*TimerIdentity){
		func(timerPlanID *TimerIdentity) { timerPlanID.TimerPlanDay = 7 },
		func(timerPlanID *TimerIdentity) { timerPlanID.TimerPlanFamily = "" },
		func(timerPlanID *TimerIdentity) { timerPlanID.TimerPlanChannel = "" },
		func(timerPlanID *TimerIdentity) { timerPlanID.TimerPlanSystemType = 2 },
		func(timerPlanID *TimerIdentity) { timerPlanID.TimerPlanAddress = 3 },
		func(timerPlanID *TimerIdentity) { timerPlanID.TimerPlanSystemType = 0 },
	} {
		timerPlanID, timerPlanBaseline, timerPlanSlots := timerPlanFixture()
		timerPlanChange(&timerPlanID)
		timerPlanBaseline.TimerPlanIdentity = timerPlanID
		timerPlanExpectInvalid(timerPlanT, timerPlanID, timerPlanBaseline, timerPlanSlots)
	}
}

func TestAb08RawSelectorsAndBoundaries(timerPlanT *testing.T) {
	for _, timerPlanPair := range [][2]int{{0, 1}, {0, 2}, {0, 3}, {1, 1}, {1, 2}, {3, 1}, {3, 2}} {
		for _, timerPlanInstance := range []int{0, 255} {
			timerPlanID, timerPlanBaseline, _ := timerPlanFixture()
			timerPlanID.TimerPlanSystemType, timerPlanID.TimerPlanAddress = timerPlanPair[0], timerPlanPair[1]
			timerPlanID.TimerPlanInstance = timerPlanInstance
			timerPlanID.TimerPlanDay = 0
			timerPlanBaseline.TimerPlanIdentity = timerPlanID
			timerPlanSlots := []TimerSlot{{0, 1}, {1, 144}, {144, 144}}
			if timerPlanPair[0] != 3 && timerPlanInstance != 0 {
				timerPlanExpectInvalid(timerPlanT, timerPlanID, timerPlanBaseline, timerPlanSlots)
				continue
			}
			timerPlanPlan, timerPlanErr := BuildTimerPlan(timerPlanID, timerPlanBaseline, timerPlanSlots)
			if timerPlanErr != nil {
				timerPlanT.Fatal(timerPlanErr)
			}
			timerPlanWant := [11]byte{4, byte(timerPlanPair[0]), byte(timerPlanInstance), byte(timerPlanPair[1]), 0, 0, 1, 1, 144, 144, 144}
			if timerPlanPlan.TimerPlanPayload() != timerPlanWant {
				timerPlanT.Fatalf("raw selectors changed: % X; want % X", timerPlanPlan.TimerPlanPayload(), timerPlanWant)
			}
		}
	}
}
