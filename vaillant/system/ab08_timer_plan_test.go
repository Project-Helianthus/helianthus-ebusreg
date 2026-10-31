package system

import (
	"errors"
	"testing"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
)

func ab08Fixture() (Ab08TimerIdentity, Ab08TimerBaseline, []Ab08TimerSlot) {
	ab08Identity := Ab08TimerIdentity{"synthetic-family", "raw-channel", 3, 2, 2, 6}
	return ab08Identity, Ab08TimerBaseline{ab08Identity, []byte{6, 0x90, 0x90, 0x90, 0x90, 0x90, 0x90}},
		[]Ab08TimerSlot{{6, 12}, {0x90, 0x90}, {0x8f, 0x90}}
}

func TestAb08ExactPayloadAndOwnership(ab08T *testing.T) {
	ab08Identity, ab08Baseline, ab08Slots := ab08Fixture()
	ab08Plan, ab08Err := Ab08BuildTimerPlan(ab08Identity, ab08Baseline, ab08Slots)
	if ab08Err != nil {
		ab08T.Fatal(ab08Err)
	}
	ab08Want := [11]byte{0x04, 0x03, 0x02, 0x02, 0x06, 0x06, 0x0c, 0x90, 0x90, 0x8f, 0x90}
	if ab08Plan.Ab08Payload() != ab08Want {
		ab08T.Fatalf("payload = % X; want % X", ab08Plan.Ab08Payload(), ab08Want)
	}
	ab08Raw := ab08Plan.Ab08Baseline()
	if ab08Raw != [7]byte{6, 0x90, 0x90, 0x90, 0x90, 0x90, 0x90} || ab08Plan.Ab08Identity() != ab08Identity {
		ab08T.Fatal("baseline or identity not preserved")
	}
	ab08Baseline.Ab08Raw[0] = 99
	ab08Baseline.Ab08Identity.Ab08Family = "changed-family"
	ab08Slots[0] = Ab08TimerSlot{0, 1}
	ab08Identity.Ab08Day = 0
	ab08ReturnedIdentity := ab08Plan.Ab08Identity()
	ab08ReturnedIdentity.Ab08Day = 1
	ab08Raw[0] = 100
	ab08ReturnedPayload := ab08Plan.Ab08Payload()
	ab08ReturnedPayload[0] = 0
	if ab08Plan.Ab08Baseline()[0] != 6 || ab08Plan.Ab08Payload() != ab08Want || ab08Plan.Ab08Identity().Ab08Day != 6 || ab08Plan.Ab08Identity().Ab08Family != "synthetic-family" {
		ab08T.Fatal("plan aliases caller or accessor memory")
	}
	if ab08Plan.Ab08NativeWriteAdmission() != "Unknown" {
		ab08T.Fatal("offline planning must not promote native write admission")
	}
}

func TestAb08IdentityMismatch(ab08T *testing.T) {
	ab08Cases := []struct {
		ab08Name   string
		ab08Change func(*Ab08TimerIdentity)
	}{
		{"ab08Family", func(ab08ID *Ab08TimerIdentity) { ab08ID.Ab08Family = "other-family" }},
		{"ab08Channel", func(ab08ID *Ab08TimerIdentity) { ab08ID.Ab08Channel = "other-channel" }},
		{"ab08SystemType", func(ab08ID *Ab08TimerIdentity) { ab08ID.Ab08SystemType = 1 }},
		{"ab08Instance", func(ab08ID *Ab08TimerIdentity) { ab08ID.Ab08Instance = 255 }},
		{"ab08Address", func(ab08ID *Ab08TimerIdentity) { ab08ID.Ab08Address = 1 }},
		{"ab08Day", func(ab08ID *Ab08TimerIdentity) { ab08ID.Ab08Day = 5 }},
	}
	for _, ab08Case := range ab08Cases {
		ab08T.Run(ab08Case.ab08Name, func(ab08T *testing.T) {
			ab08ID, ab08Baseline, ab08Slots := ab08Fixture()
			ab08ID.Ab08Instance = 0
			ab08Baseline.Ab08Identity = ab08ID
			if _, ab08Err := Ab08BuildTimerPlan(ab08ID, ab08Baseline, ab08Slots); ab08Err != nil {
				ab08T.Fatal("valid control:", ab08Err)
			}
			ab08Case.ab08Change(&ab08ID)
			ab08ExpectInvalid(ab08T, ab08ID, ab08Baseline, ab08Slots)
			ab08Baseline.Ab08Identity = ab08ID
			if _, ab08Err := Ab08BuildTimerPlan(ab08ID, ab08Baseline, ab08Slots); ab08Err != nil {
				ab08T.Fatal("independently changed valid control:", ab08Err)
			}
		})
	}
}

func ab08ExpectInvalid(ab08T *testing.T, ab08ID Ab08TimerIdentity, ab08Baseline Ab08TimerBaseline, ab08Slots []Ab08TimerSlot) {
	ab08T.Helper()
	ab08Plan, ab08Err := Ab08BuildTimerPlan(ab08ID, ab08Baseline, ab08Slots)
	if !errors.Is(ab08Err, ebuserrors.ErrInvalidPayload) || ab08Plan != (Ab08TimerPlan{}) {
		ab08T.Fatalf("expected invalid payload and empty plan; got %v, %+v", ab08Err, ab08Plan)
	}
}

func TestAb08InvalidWidthsAndSlots(ab08T *testing.T) {
	for _, ab08Width := range []int{0, 6, 8} {
		ab08ID, ab08Baseline, ab08Slots := ab08Fixture()
		ab08Baseline.Ab08Raw = make([]byte, ab08Width)
		ab08ExpectInvalid(ab08T, ab08ID, ab08Baseline, ab08Slots)
	}
	for _, ab08Count := range []int{0, 2, 4} {
		ab08ID, ab08Baseline, _ := ab08Fixture()
		ab08Slots := make([]Ab08TimerSlot, ab08Count)
		for ab08Index := range ab08Slots {
			ab08Slots[ab08Index] = Ab08TimerSlot{0, 1}
		}
		ab08ExpectInvalid(ab08T, ab08ID, ab08Baseline, ab08Slots)
	}
	for _, ab08Bad := range []Ab08TimerSlot{{144, 145}, {143, 145}, {12, 6}, {6, 6}, {144, 143}, {0, 0}, {-1, 12}, {256, 12}, {6, -1}, {6, 256}, {0x90, 0x8f}} {
		for ab08Index := 0; ab08Index < 3; ab08Index++ {
			ab08ID, ab08Baseline, ab08Slots := ab08Fixture()
			ab08Slots[ab08Index] = ab08Bad
			ab08ExpectInvalid(ab08T, ab08ID, ab08Baseline, ab08Slots)
		}
	}
}

func TestAb08IdentityBounds(ab08T *testing.T) {
	ab08Changes := []func(*Ab08TimerIdentity, int){
		func(ab08ID *Ab08TimerIdentity, ab08Value int) { ab08ID.Ab08SystemType = ab08Value },
		func(ab08ID *Ab08TimerIdentity, ab08Value int) { ab08ID.Ab08Instance = ab08Value },
		func(ab08ID *Ab08TimerIdentity, ab08Value int) { ab08ID.Ab08Address = ab08Value },
		func(ab08ID *Ab08TimerIdentity, ab08Value int) { ab08ID.Ab08Day = ab08Value },
	}
	for _, ab08Change := range ab08Changes {
		for _, ab08Value := range []int{-1, 256} {
			ab08ID, ab08Baseline, ab08Slots := ab08Fixture()
			ab08Change(&ab08ID, ab08Value)
			ab08Baseline.Ab08Identity = ab08ID
			ab08ExpectInvalid(ab08T, ab08ID, ab08Baseline, ab08Slots)
		}
	}
	for _, ab08Change := range []func(*Ab08TimerIdentity){
		func(ab08ID *Ab08TimerIdentity) { ab08ID.Ab08Day = 7 },
		func(ab08ID *Ab08TimerIdentity) { ab08ID.Ab08Family = "" },
		func(ab08ID *Ab08TimerIdentity) { ab08ID.Ab08Channel = "" },
		func(ab08ID *Ab08TimerIdentity) { ab08ID.Ab08SystemType = 2 },
		func(ab08ID *Ab08TimerIdentity) { ab08ID.Ab08Address = 3 },
		func(ab08ID *Ab08TimerIdentity) { ab08ID.Ab08SystemType = 0 },
	} {
		ab08ID, ab08Baseline, ab08Slots := ab08Fixture()
		ab08Change(&ab08ID)
		ab08Baseline.Ab08Identity = ab08ID
		ab08ExpectInvalid(ab08T, ab08ID, ab08Baseline, ab08Slots)
	}
}

func TestAb08RawSelectorsAndBoundaries(ab08T *testing.T) {
	for _, ab08Pair := range [][2]int{{0, 1}, {0, 2}, {0, 3}, {1, 1}, {1, 2}, {3, 1}, {3, 2}} {
		for _, ab08Instance := range []int{0, 255} {
			ab08ID, ab08Baseline, _ := ab08Fixture()
			ab08ID.Ab08SystemType, ab08ID.Ab08Address = ab08Pair[0], ab08Pair[1]
			ab08ID.Ab08Instance = ab08Instance
			ab08ID.Ab08Day = 0
			ab08Baseline.Ab08Identity = ab08ID
			ab08Slots := []Ab08TimerSlot{{0, 1}, {1, 144}, {144, 144}}
			if ab08Pair[0] != 3 && ab08Instance != 0 {
				ab08ExpectInvalid(ab08T, ab08ID, ab08Baseline, ab08Slots)
				continue
			}
			ab08Plan, ab08Err := Ab08BuildTimerPlan(ab08ID, ab08Baseline, ab08Slots)
			if ab08Err != nil {
				ab08T.Fatal(ab08Err)
			}
			ab08Want := [11]byte{4, byte(ab08Pair[0]), byte(ab08Instance), byte(ab08Pair[1]), 0, 0, 1, 1, 144, 144, 144}
			if ab08Plan.Ab08Payload() != ab08Want {
				ab08T.Fatalf("raw selectors changed: % X; want % X", ab08Plan.Ab08Payload(), ab08Want)
			}
		}
	}
}
