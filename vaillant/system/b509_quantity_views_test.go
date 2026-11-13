package system

import (
	b509QuantityViewHexPackage "encoding/hex"
	b509QuantityViewJSON "encoding/json"
	b509QuantityViewOs "os"
	b509QuantityViewReflect "reflect"
	b509QuantityViewTesting "testing"
)

type b509QuantityViewVector struct {
	B509QuantityViewName      string                    `json:"Ab09An09Name"`
	B509QuantityViewIdentity  B509QuantityViewSelection `json:"Ab09An09Identity"`
	B509QuantityViewSelection B509QuantityViewSelection `json:"Ab09An09Selection"`
	B509QuantityViewSelected  bool                      `json:"Ab09An09Selected"`
	B509QuantityViewRaw       []*string                 `json:"Ab09An09Raw"`
	B509QuantityViewStates    []string                  `json:"Ab09An09States"`
	B509QuantityViewValues    []*float64                `json:"Ab09An09Values"`
}

func b509QuantityViewBytes(b509QuantityViewT *b509QuantityViewTesting.T, b509QuantityViewHex *string) []byte {
	b509QuantityViewT.Helper()
	if b509QuantityViewHex == nil {
		return nil
	}
	b509QuantityViewBytes, b509QuantityViewErr := b509QuantityViewHexPackage.DecodeString(*b509QuantityViewHex)
	if b509QuantityViewErr != nil {
		b509QuantityViewT.Fatal(b509QuantityViewErr)
	}
	return b509QuantityViewBytes
}

func b509QuantityViewCandidates(b509QuantityViewView B509QuantityView) []B509QuantityViewCandidate {
	return []B509QuantityViewCandidate{b509QuantityViewView.B509QuantityViewTemperature, b509QuantityViewView.B509QuantityViewFlow, b509QuantityViewView.B509QuantityViewVolume, b509QuantityViewView.B509QuantityViewMinimumVolume, b509QuantityViewView.B509QuantityViewMaximumVolume}
}

func TestAb09An09IndependentVectors(b509QuantityViewT *b509QuantityViewTesting.T) {
	b509QuantityViewData, b509QuantityViewErr := b509QuantityViewOs.ReadFile("testdata/ab09an09_b509_quantity_views.json")
	if b509QuantityViewErr != nil {
		b509QuantityViewT.Fatal(b509QuantityViewErr)
	}
	var b509QuantityViewVectors []b509QuantityViewVector
	if b509QuantityViewErr = b509QuantityViewJSON.Unmarshal(b509QuantityViewData, &b509QuantityViewVectors); b509QuantityViewErr != nil {
		b509QuantityViewT.Fatal(b509QuantityViewErr)
	}
	if len(b509QuantityViewVectors) != 41 {
		b509QuantityViewT.Fatalf("ab09an09 vector count = %d", len(b509QuantityViewVectors))
	}
	for _, b509QuantityViewCase := range b509QuantityViewVectors {
		b509QuantityViewT.Run(b509QuantityViewCase.B509QuantityViewName, func(b509QuantityViewT *b509QuantityViewTesting.T) {
			if len(b509QuantityViewCase.B509QuantityViewRaw) != 5 || len(b509QuantityViewCase.B509QuantityViewStates) != 5 || len(b509QuantityViewCase.B509QuantityViewValues) != 5 {
				b509QuantityViewT.Fatal("ab09an09 incomplete oracle")
			}
			b509QuantityViewInput := B509QuantityViewFields{
				b509QuantityViewBytes(b509QuantityViewT, b509QuantityViewCase.B509QuantityViewRaw[0]),
				b509QuantityViewBytes(b509QuantityViewT, b509QuantityViewCase.B509QuantityViewRaw[1]),
				b509QuantityViewBytes(b509QuantityViewT, b509QuantityViewCase.B509QuantityViewRaw[2]),
				b509QuantityViewBytes(b509QuantityViewT, b509QuantityViewCase.B509QuantityViewRaw[3]),
				b509QuantityViewBytes(b509QuantityViewT, b509QuantityViewCase.B509QuantityViewRaw[4]),
			}
			b509QuantityViewBefore := B509QuantityViewFields{b509QuantityViewBytes(b509QuantityViewT, b509QuantityViewCase.B509QuantityViewRaw[0]), b509QuantityViewBytes(b509QuantityViewT, b509QuantityViewCase.B509QuantityViewRaw[1]), b509QuantityViewBytes(b509QuantityViewT, b509QuantityViewCase.B509QuantityViewRaw[2]), b509QuantityViewBytes(b509QuantityViewT, b509QuantityViewCase.B509QuantityViewRaw[3]), b509QuantityViewBytes(b509QuantityViewT, b509QuantityViewCase.B509QuantityViewRaw[4])}
			b509QuantityViewView := B509QuantityViews(b509QuantityViewCase.B509QuantityViewIdentity, b509QuantityViewCase.B509QuantityViewSelection, b509QuantityViewInput)
			if !b509QuantityViewReflect.DeepEqual(b509QuantityViewInput, b509QuantityViewBefore) {
				b509QuantityViewT.Fatal("ab09an09 input mutated")
			}
			if b509QuantityViewView.B509QuantityViewSelected != b509QuantityViewCase.B509QuantityViewSelected || b509QuantityViewView.B509QuantityViewIdentity != b509QuantityViewCase.B509QuantityViewIdentity || b509QuantityViewView.B509QuantityViewSelection != b509QuantityViewCase.B509QuantityViewSelection || b509QuantityViewView.B509QuantityViewContractGrade != "Proven" {
				b509QuantityViewT.Fatalf("ab09an09 selection/contract = %+v", b509QuantityViewView)
			}
			for b509QuantityViewIndex, b509QuantityViewGot := range b509QuantityViewCandidates(b509QuantityViewView) {
				b509QuantityViewGrade := "Hypothesis"
				if b509QuantityViewIndex >= 3 {
					b509QuantityViewGrade = "Unknown"
				}
				if !b509QuantityViewReflect.DeepEqual(b509QuantityViewGot.B509QuantityViewRaw, b509QuantityViewBytes(b509QuantityViewT, b509QuantityViewCase.B509QuantityViewRaw[b509QuantityViewIndex])) || b509QuantityViewGot.B509QuantityViewState != b509QuantityViewCase.B509QuantityViewStates[b509QuantityViewIndex] || !b509QuantityViewReflect.DeepEqual(b509QuantityViewGot.B509QuantityViewValue, b509QuantityViewCase.B509QuantityViewValues[b509QuantityViewIndex]) || b509QuantityViewGot.B509QuantityViewSemanticGrade != b509QuantityViewGrade || b509QuantityViewGot.B509QuantityViewNativeQualification != "Unknown" || b509QuantityViewGot.B509QuantityViewUnit != []string{"degrees Celsius", "litres per minute", "litres", "", ""}[b509QuantityViewIndex] {
					b509QuantityViewT.Errorf("ab09an09 field %d = %+v; expected state %s value %v", b509QuantityViewIndex, b509QuantityViewGot, b509QuantityViewCase.B509QuantityViewStates[b509QuantityViewIndex], b509QuantityViewCase.B509QuantityViewValues[b509QuantityViewIndex])
				}
			}
		})
	}
}

func TestAb09An09Ownership(b509QuantityViewT *b509QuantityViewTesting.T) {
	b509QuantityViewIdentity := B509QuantityViewSelection{"ab09an09_synthetic_family", "ab09an09_synthetic_variant", 0x59, "ab09an09_response"}
	b509QuantityViewShared := []byte{0x10, 0}
	b509QuantityViewInput := B509QuantityViewFields{b509QuantityViewShared, b509QuantityViewShared, b509QuantityViewShared, b509QuantityViewShared, b509QuantityViewShared}
	b509QuantityViewFirst := B509QuantityViews(b509QuantityViewIdentity, b509QuantityViewIdentity, b509QuantityViewInput)
	b509QuantityViewSecond := B509QuantityViews(b509QuantityViewIdentity, b509QuantityViewIdentity, b509QuantityViewInput)
	b509QuantityViewShared[0] = 0xff
	for _, b509QuantityViewGot := range b509QuantityViewCandidates(b509QuantityViewFirst) {
		if !b509QuantityViewReflect.DeepEqual(b509QuantityViewGot.B509QuantityViewRaw, []byte{0x10, 0}) {
			b509QuantityViewT.Fatal("ab09an09 input aliases result")
		}
	}
	for b509QuantityViewIndex, b509QuantityViewGot := range b509QuantityViewCandidates(b509QuantityViewFirst) {
		b509QuantityViewGot.B509QuantityViewRaw[0] = byte(b509QuantityViewIndex)
		if b509QuantityViewGot.B509QuantityViewValue != nil {
			*b509QuantityViewGot.B509QuantityViewValue = -999
		}
		for b509QuantityViewOtherIndex, b509QuantityViewOther := range b509QuantityViewCandidates(b509QuantityViewFirst) {
			if b509QuantityViewOtherIndex > b509QuantityViewIndex && b509QuantityViewOther.B509QuantityViewRaw[0] != 0x10 {
				b509QuantityViewT.Fatal("ab09an09 fields alias")
			}
		}
	}
	for b509QuantityViewIndex, b509QuantityViewGot := range b509QuantityViewCandidates(b509QuantityViewSecond) {
		if !b509QuantityViewReflect.DeepEqual(b509QuantityViewGot.B509QuantityViewRaw, []byte{0x10, 0}) {
			b509QuantityViewT.Fatal("ab09an09 results alias")
		}
		if b509QuantityViewIndex < 3 && (b509QuantityViewGot.B509QuantityViewValue == nil || *b509QuantityViewGot.B509QuantityViewValue != []float64{1, 0.16, 1.6}[b509QuantityViewIndex]) {
			b509QuantityViewT.Fatal("ab09an09 scalar results alias")
		}
	}
	if b509QuantityViewShared[0] != 0xff {
		b509QuantityViewT.Fatal("ab09an09 returned raw aliases input")
	}
}
