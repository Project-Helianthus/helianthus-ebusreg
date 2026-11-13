package system

import (
	b524CalendarViewJSON "encoding/json"
	b524CalendarViewOS "os"
	b524CalendarViewReflect "reflect"
	b524CalendarViewTesting "testing"
)

type b524CalendarViewVector struct {
	B524CalendarViewLabel       string  `json:"Ab01An01Label"`
	B524CalendarViewFamily      string  `json:"Ab01An01Family"`
	B524CalendarViewGG          byte    `json:"Ab01An01GG"`
	B524CalendarViewII          byte    `json:"Ab01An01II"`
	B524CalendarViewRR          uint16  `json:"Ab01An01RR"`
	B524CalendarViewAlternative string  `json:"Ab01An01Alternative"`
	B524CalendarViewBytes       [3]byte `json:"Ab01An01Bytes"`
	B524CalendarViewExpected    [3]byte `json:"Ab01An01Expected"`
}

func b524CalendarViewBase() B524CalendarViewInput {
	b524CalendarViewII := byte(0)
	return B524CalendarViewInput{B524CalendarViewFamily: "controller-700", B524CalendarViewSoftware: "synthetic-software", B524CalendarViewHardware: "synthetic-hardware", B524CalendarViewAPIVersion: "synthetic-api", B524CalendarViewAPIRevision: "synthetic-revision", B524CalendarViewOP: 2, B524CalendarViewII: &b524CalendarViewII, B524CalendarViewRequestedRR: 0x34, B524CalendarViewReturnedRR: 0x34, B524CalendarViewAlternative: "date", B524CalendarViewResponseState: "echo-payload", B524CalendarViewRaw: []byte{31, 12, 99}, B524CalendarViewCells: [3]B524CalendarViewCell{{31, "present"}, {12, "present"}, {99, "present"}}}
}

func b524CalendarViewAssert(b524CalendarViewT *b524CalendarViewTesting.T, b524CalendarViewInput B524CalendarViewInput, b524CalendarViewExpected [3]*byte, b524CalendarViewStates [3]string, b524CalendarViewComplete bool, b524CalendarViewGrade string) {
	b524CalendarViewT.Helper()
	b524CalendarViewResult := B524CalendarViewInterpret(b524CalendarViewInput)
	if !b524CalendarViewReflect.DeepEqual(b524CalendarViewResult.B524CalendarViewContext, b524CalendarViewInput) {
		b524CalendarViewT.Fatal("ab01an01 context lost")
	}
	if b524CalendarViewResult.B524CalendarViewCompleteCandidate != b524CalendarViewComplete || b524CalendarViewResult.B524CalendarViewGrade != b524CalendarViewGrade || b524CalendarViewResult.B524CalendarViewContractGrade != "Proven" || b524CalendarViewResult.B524CalendarViewNativeQualification != "Unknown" || b524CalendarViewResult.B524CalendarViewCombinedTimestamp != "Unknown" || b524CalendarViewResult.B524CalendarViewEditValidation != "Unknown" {
		b524CalendarViewT.Fatalf("ab01an01 result: %+v", b524CalendarViewResult)
	}
	for b524CalendarViewIndex, b524CalendarViewComponent := range b524CalendarViewResult.B524CalendarViewComponents {
		if b524CalendarViewComponent.B524CalendarViewCell != b524CalendarViewInput.B524CalendarViewCells[b524CalendarViewIndex] || b524CalendarViewComponent.B524CalendarViewState != b524CalendarViewStates[b524CalendarViewIndex] || !b524CalendarViewReflect.DeepEqual(b524CalendarViewComponent.B524CalendarViewCandidate, b524CalendarViewExpected[b524CalendarViewIndex]) {
			b524CalendarViewT.Fatalf("ab01an01 component %d: %+v", b524CalendarViewIndex, b524CalendarViewComponent)
		}
	}
}

func TestAb01An01Vectors(b524CalendarViewT *b524CalendarViewTesting.T) {
	b524CalendarViewData, b524CalendarViewErr := b524CalendarViewOS.ReadFile("testdata/ab01an01_b524_calendar_views.json")
	if b524CalendarViewErr != nil {
		b524CalendarViewT.Fatal(b524CalendarViewErr)
	}
	var b524CalendarViewVectors []b524CalendarViewVector
	if b524CalendarViewErr = b524CalendarViewJSON.Unmarshal(b524CalendarViewData, &b524CalendarViewVectors); b524CalendarViewErr != nil {
		b524CalendarViewT.Fatal(b524CalendarViewErr)
	}
	if len(b524CalendarViewVectors) != 26 {
		b524CalendarViewT.Fatal("ab01an01 expected 26 rows")
	}
	b524CalendarViewSeen := map[string]map[[3]uint16]bool{}
	for _, b524CalendarViewVector := range b524CalendarViewVectors {
		b524CalendarViewT.Run(b524CalendarViewVector.B524CalendarViewLabel, func(b524CalendarViewT *b524CalendarViewTesting.T) {
			b524CalendarViewInput := b524CalendarViewBase()
			b524CalendarViewInput.B524CalendarViewFamily = b524CalendarViewVector.B524CalendarViewFamily
			b524CalendarViewInput.B524CalendarViewRequestedGG = b524CalendarViewVector.B524CalendarViewGG
			b524CalendarViewInput.B524CalendarViewReturnedGG = b524CalendarViewVector.B524CalendarViewGG
			*b524CalendarViewInput.B524CalendarViewII = b524CalendarViewVector.B524CalendarViewII
			b524CalendarViewInput.B524CalendarViewRequestedRR = b524CalendarViewVector.B524CalendarViewRR
			b524CalendarViewInput.B524CalendarViewReturnedRR = b524CalendarViewVector.B524CalendarViewRR
			b524CalendarViewInput.B524CalendarViewAlternative = b524CalendarViewVector.B524CalendarViewAlternative
			for b524CalendarViewIndex, b524CalendarViewByte := range b524CalendarViewVector.B524CalendarViewBytes {
				b524CalendarViewInput.B524CalendarViewCells[b524CalendarViewIndex].B524CalendarViewRaw = b524CalendarViewByte
			}
			b524CalendarViewAssert(b524CalendarViewT, b524CalendarViewInput, [3]*byte{&b524CalendarViewVector.B524CalendarViewExpected[0], &b524CalendarViewVector.B524CalendarViewExpected[1], &b524CalendarViewVector.B524CalendarViewExpected[2]}, [3]string{"candidate", "candidate", "candidate"}, true, "Hypothesis")
			if b524CalendarViewSeen[b524CalendarViewVector.B524CalendarViewFamily] == nil {
				b524CalendarViewSeen[b524CalendarViewVector.B524CalendarViewFamily] = map[[3]uint16]bool{}
			}
			b524CalendarViewKey := [3]uint16{uint16(b524CalendarViewVector.B524CalendarViewGG), uint16(b524CalendarViewVector.B524CalendarViewII), b524CalendarViewVector.B524CalendarViewRR}
			if b524CalendarViewSeen[b524CalendarViewVector.B524CalendarViewFamily][b524CalendarViewKey] {
				b524CalendarViewT.Fatal("ab01an01 duplicate row")
			}
			b524CalendarViewSeen[b524CalendarViewVector.B524CalendarViewFamily][b524CalendarViewKey] = true
		})
	}
}

func TestAb01An01Context(b524CalendarViewT *b524CalendarViewTesting.T) {
	for b524CalendarViewAxis := 0; b524CalendarViewAxis < 20; b524CalendarViewAxis++ {
		b524CalendarViewInput := b524CalendarViewBase()
		switch b524CalendarViewAxis {
		case 0:
			b524CalendarViewInput.B524CalendarViewFamily = ""
		case 1:
			b524CalendarViewInput.B524CalendarViewFamily = "unknown-family"
		case 2:
			b524CalendarViewInput.B524CalendarViewSoftware = ""
		case 3:
			b524CalendarViewInput.B524CalendarViewHardware = ""
		case 4:
			b524CalendarViewInput.B524CalendarViewAPIVersion = ""
		case 5:
			b524CalendarViewInput.B524CalendarViewAPIRevision = ""
		case 6:
			b524CalendarViewInput.B524CalendarViewOP = 6
		case 7:
			b524CalendarViewInput.B524CalendarViewRequestedGG = 2
			b524CalendarViewInput.B524CalendarViewReturnedGG = 2
		case 8:
			*b524CalendarViewInput.B524CalendarViewII = 1
		case 9:
			b524CalendarViewInput.B524CalendarViewII = nil
		case 10:
			b524CalendarViewInput.B524CalendarViewRequestedRR = 0xffff
			b524CalendarViewInput.B524CalendarViewReturnedRR = 0xffff
		case 11:
			b524CalendarViewInput.B524CalendarViewAlternative = "time"
		case 12:
			b524CalendarViewInput.B524CalendarViewReturnedGG = 1
		case 13:
			b524CalendarViewInput.B524CalendarViewReturnedRR = 0x35
		case 14:
			b524CalendarViewInput.B524CalendarViewResponseState = "one-byte-status"
		case 15:
			b524CalendarViewInput.B524CalendarViewResponseState = "transport-error"
		case 16:
			b524CalendarViewInput.B524CalendarViewResponseState = ""
		case 17:
			b524CalendarViewInput.B524CalendarViewRequestedRR = 0xda
			b524CalendarViewInput.B524CalendarViewReturnedRR = 0xda
		case 18:
			b524CalendarViewInput.B524CalendarViewFamily = "controller-720"
			b524CalendarViewInput.B524CalendarViewRequestedRR = 0x71
			b524CalendarViewInput.B524CalendarViewReturnedRR = 0x71
		case 19:
			b524CalendarViewInput.B524CalendarViewAlternative = "timestamp"
		}
		b524CalendarViewAssert(b524CalendarViewT, b524CalendarViewInput, [3]*byte{}, [3]string{"Unknown", "Unknown", "Unknown"}, false, "Unknown")
	}
	for b524CalendarViewAxis := 0; b524CalendarViewAxis < 5; b524CalendarViewAxis++ {
		b524CalendarViewInput := b524CalendarViewBase()
		switch b524CalendarViewAxis {
		case 0:
			b524CalendarViewInput.B524CalendarViewSoftware = "synthetic-other"
		case 1:
			b524CalendarViewInput.B524CalendarViewHardware = "synthetic-other"
		case 2:
			b524CalendarViewInput.B524CalendarViewAPIVersion = "synthetic-other"
		case 3:
			b524CalendarViewInput.B524CalendarViewAPIRevision = "synthetic-other"
		case 4:
			b524CalendarViewInput.B524CalendarViewRequestedGG = 3
			b524CalendarViewInput.B524CalendarViewReturnedGG = 3
			b524CalendarViewInput.B524CalendarViewRequestedRR = 3
			b524CalendarViewInput.B524CalendarViewReturnedRR = 3
			*b524CalendarViewInput.B524CalendarViewII = 255
		}
		b524CalendarViewExpected := [3]byte{31, 12, 99}
		b524CalendarViewAssert(b524CalendarViewT, b524CalendarViewInput, [3]*byte{&b524CalendarViewExpected[0], &b524CalendarViewExpected[1], &b524CalendarViewExpected[2]}, [3]string{"candidate", "candidate", "candidate"}, true, "Hypothesis")
	}
}

func TestAb01An01Components(b524CalendarViewT *b524CalendarViewTesting.T) {
	for b524CalendarViewIndex := 0; b524CalendarViewIndex < 3; b524CalendarViewIndex++ {
		for _, b524CalendarViewState := range []string{"missing", "failed", "invalid"} {
			b524CalendarViewInput := b524CalendarViewBase()
			b524CalendarViewInput.B524CalendarViewCells[b524CalendarViewIndex].B524CalendarViewState = b524CalendarViewState
			b524CalendarViewExpected := [3]*byte{&b524CalendarViewInput.B524CalendarViewCells[0].B524CalendarViewRaw, &b524CalendarViewInput.B524CalendarViewCells[1].B524CalendarViewRaw, &b524CalendarViewInput.B524CalendarViewCells[2].B524CalendarViewRaw}
			b524CalendarViewExpected[b524CalendarViewIndex] = nil
			b524CalendarViewStates := [3]string{"candidate", "candidate", "candidate"}
			if b524CalendarViewState == "invalid" {
				b524CalendarViewState = "Unknown"
			}
			b524CalendarViewStates[b524CalendarViewIndex] = b524CalendarViewState
			b524CalendarViewAssert(b524CalendarViewT, b524CalendarViewInput, b524CalendarViewExpected, b524CalendarViewStates, false, "Hypothesis")
		}
		for b524CalendarViewValue := 0; b524CalendarViewValue < 256; b524CalendarViewValue++ {
			b524CalendarViewInput := b524CalendarViewBase()
			b524CalendarViewInput.B524CalendarViewCells[b524CalendarViewIndex].B524CalendarViewRaw = byte(b524CalendarViewValue)
			b524CalendarViewExpected := [3]*byte{&b524CalendarViewInput.B524CalendarViewCells[0].B524CalendarViewRaw, &b524CalendarViewInput.B524CalendarViewCells[1].B524CalendarViewRaw, &b524CalendarViewInput.B524CalendarViewCells[2].B524CalendarViewRaw}
			b524CalendarViewStates := [3]string{"candidate", "candidate", "candidate"}
			b524CalendarViewValid := (b524CalendarViewIndex == 0 && b524CalendarViewValue >= 1 && b524CalendarViewValue <= 31) || (b524CalendarViewIndex == 1 && b524CalendarViewValue >= 1 && b524CalendarViewValue <= 12) || (b524CalendarViewIndex == 2 && b524CalendarViewValue <= 99)
			if !b524CalendarViewValid {
				b524CalendarViewExpected[b524CalendarViewIndex] = nil
				b524CalendarViewStates[b524CalendarViewIndex] = "out-of-domain"
			}
			b524CalendarViewAssert(b524CalendarViewT, b524CalendarViewInput, b524CalendarViewExpected, b524CalendarViewStates, b524CalendarViewValid, "Hypothesis")
			b524CalendarViewInput.B524CalendarViewRequestedRR = 0x35
			b524CalendarViewInput.B524CalendarViewReturnedRR = 0x35
			b524CalendarViewInput.B524CalendarViewAlternative = "time"
			b524CalendarViewExpected[b524CalendarViewIndex] = &b524CalendarViewInput.B524CalendarViewCells[b524CalendarViewIndex].B524CalendarViewRaw
			b524CalendarViewAssert(b524CalendarViewT, b524CalendarViewInput, b524CalendarViewExpected, [3]string{"candidate", "candidate", "candidate"}, true, "Hypothesis")
		}
	}
}

func TestAb01An01Detached(b524CalendarViewT *b524CalendarViewTesting.T) {
	b524CalendarViewInput := b524CalendarViewBase()
	b524CalendarViewResult := B524CalendarViewInterpret(b524CalendarViewInput)
	b524CalendarViewInput.B524CalendarViewRaw[0] = 0
	*b524CalendarViewInput.B524CalendarViewII = 9
	b524CalendarViewInput.B524CalendarViewCells[0].B524CalendarViewRaw = 0
	if b524CalendarViewResult.B524CalendarViewContext.B524CalendarViewRaw[0] != 31 || *b524CalendarViewResult.B524CalendarViewContext.B524CalendarViewII != 0 || *b524CalendarViewResult.B524CalendarViewComponents[0].B524CalendarViewCandidate != 31 {
		b524CalendarViewT.Fatal("ab01an01 input aliases result")
	}
	b524CalendarViewOther := B524CalendarViewInterpret(b524CalendarViewResult.B524CalendarViewContext)
	b524CalendarViewResult.B524CalendarViewContext.B524CalendarViewRaw[1] = 0
	*b524CalendarViewResult.B524CalendarViewContext.B524CalendarViewII = 8
	*b524CalendarViewResult.B524CalendarViewComponents[1].B524CalendarViewCandidate = 0
	if b524CalendarViewInput.B524CalendarViewRaw[1] != 12 || *b524CalendarViewInput.B524CalendarViewII != 9 || b524CalendarViewOther.B524CalendarViewContext.B524CalendarViewRaw[1] != 12 || *b524CalendarViewOther.B524CalendarViewContext.B524CalendarViewII != 0 || *b524CalendarViewOther.B524CalendarViewComponents[1].B524CalendarViewCandidate != 12 || b524CalendarViewResult.B524CalendarViewContext.B524CalendarViewCells[1].B524CalendarViewRaw != 12 {
		b524CalendarViewT.Fatal("ab01an01 returned storage aliases")
	}
}
