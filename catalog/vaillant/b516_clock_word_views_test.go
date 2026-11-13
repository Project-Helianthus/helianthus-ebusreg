package vaillant

import (
	b516ClockWordViewJSON "encoding/json"
	b516ClockWordViewOS "os"
	b516ClockWordViewReflect "reflect"
	b516ClockWordViewTesting "testing"
)

func b516ClockWordViewSelector(b516ClockWordViewOperation byte) B516ClockWordViewSelector {
	return B516ClockWordViewSelector{B516ClockWordViewFamily: B516ClockWordViewFamily, B516ClockWordViewVariant: B516ClockWordViewVariant, B516ClockWordViewOperation: b516ClockWordViewOperation, B516ClockWordViewDirection: "positive"}
}

func b516ClockWordViewAssertField(b516ClockWordViewT *b516ClockWordViewTesting.T, b516ClockWordViewActual B516ClockWordViewField, b516ClockWordViewState string, b516ClockWordViewValue *uint16) {
	b516ClockWordViewT.Helper()
	if !b516ClockWordViewReflect.DeepEqual(b516ClockWordViewActual, B516ClockWordViewField{B516ClockWordViewState: b516ClockWordViewState, B516ClockWordViewValue: b516ClockWordViewValue}) {
		b516ClockWordViewT.Fatalf("ab15an15 field: %+v expected %s %v", b516ClockWordViewActual, b516ClockWordViewState, b516ClockWordViewValue)
	}
}

func TestAb15An15ClockVectors(b516ClockWordViewT *b516ClockWordViewTesting.T) {
	b516ClockWordViewData, b516ClockWordViewError := b516ClockWordViewOS.ReadFile("testdata/ab15an15_b516_clock_word_views.json")
	if b516ClockWordViewError != nil {
		b516ClockWordViewT.Fatal(b516ClockWordViewError)
	}
	var b516ClockWordViewVectors []struct {
		B516ClockWordViewLabel  string     `json:"Ab15An15Label"`
		B516ClockWordViewClock  [5][]byte  `json:"Ab15An15Clock"`
		B516ClockWordViewFailed [5]bool    `json:"Ab15An15Failed"`
		B516ClockWordViewStates [5]string  `json:"Ab15An15States"`
		B516ClockWordViewValues [5]*uint16 `json:"Ab15An15Values"`
	}
	if b516ClockWordViewError = b516ClockWordViewJSON.Unmarshal(b516ClockWordViewData, &b516ClockWordViewVectors); b516ClockWordViewError != nil {
		b516ClockWordViewT.Fatal(b516ClockWordViewError)
	}
	if len(b516ClockWordViewVectors) != 50 {
		b516ClockWordViewT.Fatalf("ab15an15 count %d", len(b516ClockWordViewVectors))
	}
	for _, b516ClockWordViewVector := range b516ClockWordViewVectors {
		b516ClockWordViewT.Run(b516ClockWordViewVector.B516ClockWordViewLabel, func(b516ClockWordViewT *b516ClockWordViewTesting.T) {
			b516ClockWordViewInput := B516ClockWordViewClockInput{}
			for b516ClockWordViewIndex := range b516ClockWordViewInput.B516ClockWordViewClock {
				b516ClockWordViewInput.B516ClockWordViewClock[b516ClockWordViewIndex] = B516ClockWordViewCell{B516ClockWordViewRaw: b516ClockWordViewVector.B516ClockWordViewClock[b516ClockWordViewIndex], B516ClockWordViewFailed: b516ClockWordViewVector.B516ClockWordViewFailed[b516ClockWordViewIndex]}
			}
			b516ClockWordViewResult := B516ClockWordViewClockView(b516ClockWordViewSelector(8), b516ClockWordViewInput)
			if !b516ClockWordViewResult.B516ClockWordViewCandidate.B516ClockWordViewAdmitted || b516ClockWordViewResult.B516ClockWordViewClock == nil || b516ClockWordViewResult.B516ClockWordViewSetpointCelsius != nil || !b516ClockWordViewReflect.DeepEqual(b516ClockWordViewResult.B516ClockWordViewRaw, b516ClockWordViewInput) {
				b516ClockWordViewT.Fatal("ab15an15 clock admission/raw/missing setpoint")
			}
			for b516ClockWordViewIndex := range b516ClockWordViewInput.B516ClockWordViewClock {
				b516ClockWordViewAssertField(b516ClockWordViewT, b516ClockWordViewResult.B516ClockWordViewClock[b516ClockWordViewIndex], b516ClockWordViewVector.B516ClockWordViewStates[b516ClockWordViewIndex], b516ClockWordViewVector.B516ClockWordViewValues[b516ClockWordViewIndex])
			}
		})
	}
}

func TestAb15An15Words(b516ClockWordViewT *b516ClockWordViewTesting.T) {
	for b516ClockWordViewSlot := 0; b516ClockWordViewSlot < 4; b516ClockWordViewSlot++ {
		for _, b516ClockWordViewCase := range []struct {
			B516ClockWordViewRaw    []byte
			B516ClockWordViewFailed bool
			B516ClockWordViewState  string
			B516ClockWordViewValue  uint16
		}{
			{[]byte{0, 0}, false, "numeric", 0}, {[]byte{1, 0}, false, "numeric", 1}, {[]byte{0, 1}, false, "numeric", 256}, {[]byte{255, 255}, false, "numeric", 65535}, {[]byte{52, 18}, false, "numeric", 4660}, {nil, false, "missing", 0}, {[]byte{1}, false, "failed", 0}, {[]byte{1, 2, 3}, false, "failed", 0}, {[]byte{1, 0}, true, "failed", 0},
		} {
			b516ClockWordViewWords := [4]B516ClockWordViewCell{{B516ClockWordViewRaw: []byte{0, 0}}, {B516ClockWordViewRaw: []byte{1, 0}}, {B516ClockWordViewRaw: []byte{0, 1}}, {B516ClockWordViewRaw: []byte{255, 255}}}
			b516ClockWordViewExpected := [4]uint16{0, 1, 256, 65535}
			b516ClockWordViewWords[b516ClockWordViewSlot] = B516ClockWordViewCell{B516ClockWordViewRaw: b516ClockWordViewCase.B516ClockWordViewRaw, B516ClockWordViewFailed: b516ClockWordViewCase.B516ClockWordViewFailed}
			b516ClockWordViewResult := B516ClockWordView(b516ClockWordViewSelector(13), b516ClockWordViewWords)
			if !b516ClockWordViewResult.B516ClockWordViewCandidate.B516ClockWordViewAdmitted || b516ClockWordViewResult.B516ClockWordViewWords == nil || !b516ClockWordViewReflect.DeepEqual(b516ClockWordViewResult.B516ClockWordViewRaw, b516ClockWordViewWords) || b516ClockWordViewResult.B516ClockWordViewPositionsGrade != "Unknown" {
				b516ClockWordViewT.Fatal("ab15an15 word admission/raw/positions")
			}
			for b516ClockWordViewIndex := range b516ClockWordViewWords {
				b516ClockWordViewState := "numeric"
				b516ClockWordViewValue := &b516ClockWordViewExpected[b516ClockWordViewIndex]
				if b516ClockWordViewIndex == b516ClockWordViewSlot {
					b516ClockWordViewState = b516ClockWordViewCase.B516ClockWordViewState
					b516ClockWordViewValue = &b516ClockWordViewCase.B516ClockWordViewValue
					if b516ClockWordViewState != "numeric" {
						b516ClockWordViewValue = nil
					}
				}
				b516ClockWordViewAssertField(b516ClockWordViewT, b516ClockWordViewResult.B516ClockWordViewWords[b516ClockWordViewIndex], b516ClockWordViewState, b516ClockWordViewValue)
			}
		}
	}
}

func TestAb15An15AdmissionAndDetachedRaw(b516ClockWordViewT *b516ClockWordViewTesting.T) {
	for b516ClockWordViewAxis := 0; b516ClockWordViewAxis < 6; b516ClockWordViewAxis++ {
		b516ClockWordViewClockSelector := b516ClockWordViewSelector(8)
		b516ClockWordViewWordSelector := b516ClockWordViewSelector(13)
		for _, b516ClockWordViewSelection := range []*B516ClockWordViewSelector{&b516ClockWordViewClockSelector, &b516ClockWordViewWordSelector} {
			switch b516ClockWordViewAxis {
			case 1:
				b516ClockWordViewSelection.B516ClockWordViewFamily = "ab15an15_other_family"
			case 2:
				b516ClockWordViewSelection.B516ClockWordViewVariant = "ab15an15_other_variant"
			case 3:
				b516ClockWordViewSelection.B516ClockWordViewOperation = 0
			case 4:
				b516ClockWordViewSelection.B516ClockWordViewDirection = "request"
			case 5:
				if b516ClockWordViewSelection.B516ClockWordViewOperation == 8 {
					b516ClockWordViewSelection.B516ClockWordViewOperation = 13
				} else {
					b516ClockWordViewSelection.B516ClockWordViewOperation = 8
				}
			}
		}
		b516ClockWordViewInput := B516ClockWordViewClockInput{}
		b516ClockWordViewWords := [4]B516ClockWordViewCell{}
		for b516ClockWordViewIndex := range b516ClockWordViewInput.B516ClockWordViewClock {
			b516ClockWordViewInput.B516ClockWordViewClock[b516ClockWordViewIndex].B516ClockWordViewRaw = []byte{18}
		}
		for b516ClockWordViewIndex := range b516ClockWordViewWords {
			b516ClockWordViewWords[b516ClockWordViewIndex].B516ClockWordViewRaw = []byte{1, 2}
		}
		b516ClockWordViewBytes := [8]byte{1, 255, 254, 0, 1, 2, 3, 255}
		b516ClockWordViewInput.B516ClockWordViewSetpoint = &b516ClockWordViewBytes[0]
		b516ClockWordViewInput.B516ClockWordViewMode = &b516ClockWordViewBytes[1]
		b516ClockWordViewInput.B516ClockWordViewSpecialFunction = &b516ClockWordViewBytes[2]
		for b516ClockWordViewIndex := range b516ClockWordViewInput.B516ClockWordViewMessages {
			b516ClockWordViewInput.B516ClockWordViewMessages[b516ClockWordViewIndex] = &b516ClockWordViewBytes[b516ClockWordViewIndex+3]
		}
		b516ClockWordViewClock := B516ClockWordViewClockView(b516ClockWordViewClockSelector, b516ClockWordViewInput)
		b516ClockWordViewWord := B516ClockWordView(b516ClockWordViewWordSelector, b516ClockWordViewWords)
		if !b516ClockWordViewReflect.DeepEqual(b516ClockWordViewClock.B516ClockWordViewRaw, b516ClockWordViewInput) || !b516ClockWordViewReflect.DeepEqual(b516ClockWordViewWord.B516ClockWordViewRaw, b516ClockWordViewWords) {
			b516ClockWordViewT.Fatal("ab15an15 raw retention")
		}
		if b516ClockWordViewAxis == 0 {
			if b516ClockWordViewClock.B516ClockWordViewSetpointCelsius == nil || *b516ClockWordViewClock.B516ClockWordViewSetpointCelsius != 0.5 {
				b516ClockWordViewT.Fatal("ab15an15 half Celsius")
			}
		} else if b516ClockWordViewClock.B516ClockWordViewCandidate.B516ClockWordViewAdmitted || b516ClockWordViewWord.B516ClockWordViewCandidate.B516ClockWordViewAdmitted || b516ClockWordViewClock.B516ClockWordViewClock != nil || b516ClockWordViewWord.B516ClockWordViewWords != nil || b516ClockWordViewClock.B516ClockWordViewSetpointCelsius != nil {
			b516ClockWordViewT.Fatal("ab15an15 raw-only admission")
		}
		b516ClockWordViewExpectedCandidate := B516ClockWordViewCandidate{B516ClockWordViewAdmitted: b516ClockWordViewAxis == 0, B516ClockWordViewContractGrade: "Proven", B516ClockWordViewFieldGrade: "Hypothesis", B516ClockWordViewNativeSupport: "Unknown", B516ClockWordViewPhysicalAdaptation: "Unknown", B516ClockWordViewPresence: "Unknown", B516ClockWordViewFreshness: "Unknown", B516ClockWordViewExecution: "Unknown"}
		if b516ClockWordViewClock.B516ClockWordViewCandidate != b516ClockWordViewExpectedCandidate || b516ClockWordViewWord.B516ClockWordViewCandidate != b516ClockWordViewExpectedCandidate || b516ClockWordViewClock.B516ClockWordViewTablesGrade != "Unknown" {
			b516ClockWordViewT.Fatal("ab15an15 evidence grades")
		}
		for b516ClockWordViewIndex := range b516ClockWordViewBytes {
			b516ClockWordViewBytes[b516ClockWordViewIndex] ^= 128
		}
		if *b516ClockWordViewClock.B516ClockWordViewRaw.B516ClockWordViewSetpoint != 1 || *b516ClockWordViewClock.B516ClockWordViewRaw.B516ClockWordViewMode != 255 || *b516ClockWordViewClock.B516ClockWordViewRaw.B516ClockWordViewSpecialFunction != 254 {
			b516ClockWordViewT.Fatal("ab15an15 optional aliases")
		}
		for b516ClockWordViewIndex, b516ClockWordViewExpected := range [5]byte{0, 1, 2, 3, 255} {
			if *b516ClockWordViewClock.B516ClockWordViewRaw.B516ClockWordViewMessages[b516ClockWordViewIndex] != b516ClockWordViewExpected {
				b516ClockWordViewT.Fatal("ab15an15 message aliases/index")
			}
		}
		for b516ClockWordViewIndex := range b516ClockWordViewInput.B516ClockWordViewClock {
			b516ClockWordViewInput.B516ClockWordViewClock[b516ClockWordViewIndex].B516ClockWordViewRaw[0] = 99
			if b516ClockWordViewClock.B516ClockWordViewRaw.B516ClockWordViewClock[b516ClockWordViewIndex].B516ClockWordViewRaw[0] != 18 {
				b516ClockWordViewT.Fatal("ab15an15 clock aliases")
			}
		}
		for b516ClockWordViewIndex := range b516ClockWordViewWords {
			b516ClockWordViewWords[b516ClockWordViewIndex].B516ClockWordViewRaw[0] = 99
			b516ClockWordViewWords[b516ClockWordViewIndex].B516ClockWordViewRaw[1] = 99
			if !b516ClockWordViewReflect.DeepEqual(b516ClockWordViewWord.B516ClockWordViewRaw[b516ClockWordViewIndex].B516ClockWordViewRaw, []byte{1, 2}) {
				b516ClockWordViewT.Fatal("ab15an15 word aliases")
			}
		}
	}
}

func TestAb15An15OptionalFields(b516ClockWordViewT *b516ClockWordViewTesting.T) {
	for b516ClockWordViewIndex := 0; b516ClockWordViewIndex < 8; b516ClockWordViewIndex++ {
		for _, b516ClockWordViewCase := range []struct {
			B516ClockWordViewRaw     byte
			B516ClockWordViewCelsius float64
		}{{0, 0}, {1, 0.5}, {127, 63.5}, {255, 127.5}} {
			b516ClockWordViewValue := b516ClockWordViewCase.B516ClockWordViewRaw
			b516ClockWordViewInput := B516ClockWordViewClockInput{}
			switch b516ClockWordViewIndex {
			case 0:
				b516ClockWordViewInput.B516ClockWordViewSetpoint = &b516ClockWordViewValue
			case 1:
				b516ClockWordViewInput.B516ClockWordViewMode = &b516ClockWordViewValue
			case 2:
				b516ClockWordViewInput.B516ClockWordViewSpecialFunction = &b516ClockWordViewValue
			default:
				b516ClockWordViewInput.B516ClockWordViewMessages[b516ClockWordViewIndex-3] = &b516ClockWordViewValue
			}
			b516ClockWordViewResult := B516ClockWordViewClockView(b516ClockWordViewSelector(8), b516ClockWordViewInput)
			if !b516ClockWordViewReflect.DeepEqual(b516ClockWordViewResult.B516ClockWordViewRaw, b516ClockWordViewInput) {
				b516ClockWordViewT.Fatal("ab15an15 optional cross-contamination")
			}
			if b516ClockWordViewIndex == 0 {
				if b516ClockWordViewResult.B516ClockWordViewSetpointCelsius == nil || *b516ClockWordViewResult.B516ClockWordViewSetpointCelsius != b516ClockWordViewCase.B516ClockWordViewCelsius {
					b516ClockWordViewT.Fatal("ab15an15 optional setpoint")
				}
			} else if b516ClockWordViewResult.B516ClockWordViewSetpointCelsius != nil {
				b516ClockWordViewT.Fatal("ab15an15 missing setpoint")
			}
		}
	}
}
