package vaillant

import (
	b531LimitNumericViewHex "encoding/hex"
	b531LimitNumericViewJSON "encoding/json"
	b531LimitNumericViewOS "os"
	b531LimitNumericViewReflect "reflect"
	b531LimitNumericViewTesting "testing"
)

func b531LimitNumericViewRaw(b531LimitNumericViewT *b531LimitNumericViewTesting.T, b531LimitNumericViewText string) []byte {
	b531LimitNumericViewT.Helper()
	b531LimitNumericViewBytes, b531LimitNumericViewError := b531LimitNumericViewHex.DecodeString(b531LimitNumericViewText)
	if b531LimitNumericViewError != nil {
		b531LimitNumericViewT.Fatal(b531LimitNumericViewError)
	}
	return b531LimitNumericViewBytes
}

func TestAb19An19Vectors(b531LimitNumericViewT *b531LimitNumericViewTesting.T) {
	var b531LimitNumericViewFixture struct {
		B531LimitNumericViewLimits []struct {
			B531LimitNumericViewLabel        string `json:"Ab19An19Label"`
			B531LimitNumericViewMinimumRaw   string `json:"Ab19An19MinimumRaw"`
			B531LimitNumericViewMaximumRaw   string `json:"Ab19An19MaximumRaw"`
			B531LimitNumericViewMinimumValue uint32 `json:"Ab19An19MinimumValue"`
			B531LimitNumericViewMaximumValue uint32 `json:"Ab19An19MaximumValue"`
		} `json:"Ab19An19Limits"`
		B531LimitNumericViewGetters []struct {
			B531LimitNumericViewLabel    string `json:"Ab19An19Label"`
			B531LimitNumericViewRaw      uint32 `json:"Ab19An19Raw"`
			B531LimitNumericViewExpected uint16 `json:"Ab19An19Expected"`
		} `json:"Ab19An19Getters"`
	}
	b531LimitNumericViewBytes, b531LimitNumericViewError := b531LimitNumericViewOS.ReadFile("testdata/ab19an19_b531_limit_numeric_views.json")
	if b531LimitNumericViewError != nil {
		b531LimitNumericViewT.Fatal(b531LimitNumericViewError)
	}
	if b531LimitNumericViewError = b531LimitNumericViewJSON.Unmarshal(b531LimitNumericViewBytes, &b531LimitNumericViewFixture); b531LimitNumericViewError != nil {
		b531LimitNumericViewT.Fatal(b531LimitNumericViewError)
	}
	if len(b531LimitNumericViewFixture.B531LimitNumericViewLimits) != 6 || len(b531LimitNumericViewFixture.B531LimitNumericViewGetters) != 5 {
		b531LimitNumericViewT.Fatal("ab19an19 incomplete vectors")
	}
	for _, b531LimitNumericViewVector := range b531LimitNumericViewFixture.B531LimitNumericViewLimits {
		b531LimitNumericViewT.Run(b531LimitNumericViewVector.B531LimitNumericViewLabel, func(b531LimitNumericViewT *b531LimitNumericViewTesting.T) {
			b531LimitNumericViewMinimum := b531LimitNumericViewRaw(b531LimitNumericViewT, b531LimitNumericViewVector.B531LimitNumericViewMinimumRaw)
			b531LimitNumericViewMaximum := b531LimitNumericViewRaw(b531LimitNumericViewT, b531LimitNumericViewVector.B531LimitNumericViewMaximumRaw)
			b531LimitNumericViewResult := B531LimitNumericViewInterpret(B531LimitNumericViewFamily, B531LimitNumericViewVariant, 2, B531LimitNumericViewDirection, b531LimitNumericViewMinimum, b531LimitNumericViewMaximum, true, true, B531LimitNumericViewPair{})
			b531LimitNumericViewExpected := B531LimitNumericViewPair{
				B531LimitNumericViewMinimum: B531LimitNumericViewField{B531LimitNumericViewRaw: b531LimitNumericViewMinimum, B531LimitNumericViewState: "success", B531LimitNumericViewAvailable: true, B531LimitNumericViewValue: b531LimitNumericViewVector.B531LimitNumericViewMinimumValue, B531LimitNumericViewFresh: true},
				B531LimitNumericViewMaximum: B531LimitNumericViewField{B531LimitNumericViewRaw: b531LimitNumericViewMaximum, B531LimitNumericViewState: "success", B531LimitNumericViewAvailable: true, B531LimitNumericViewValue: b531LimitNumericViewVector.B531LimitNumericViewMaximumValue, B531LimitNumericViewFresh: true}, B531LimitNumericViewFresh: true,
			}
			if !b531LimitNumericViewResult.B531LimitNumericViewAdmitted || !b531LimitNumericViewReflect.DeepEqual(b531LimitNumericViewResult.B531LimitNumericViewAttempt, b531LimitNumericViewExpected) || !b531LimitNumericViewReflect.DeepEqual(b531LimitNumericViewResult.B531LimitNumericViewLastKnown, b531LimitNumericViewExpected) {
				b531LimitNumericViewT.Fatalf("ab19an19 result: %+v", b531LimitNumericViewResult)
			}
			if b531LimitNumericViewResult.B531LimitNumericViewContractGrade != "Proven" || b531LimitNumericViewResult.B531LimitNumericViewInterpretationGrade != "Hypothesis" || b531LimitNumericViewResult.B531LimitNumericViewNativeUnit != "Unknown" || b531LimitNumericViewResult.B531LimitNumericViewQualification != "Unknown" || b531LimitNumericViewResult.B531LimitNumericViewPresence != "Unknown" || b531LimitNumericViewResult.B531LimitNumericViewExecution != "Unknown" {
				b531LimitNumericViewT.Fatal("ab19an19 grade conflation")
			}
		})
	}
	for _, b531LimitNumericViewVector := range b531LimitNumericViewFixture.B531LimitNumericViewGetters {
		b531LimitNumericViewT.Run(b531LimitNumericViewVector.B531LimitNumericViewLabel, func(b531LimitNumericViewT *b531LimitNumericViewTesting.T) {
			if b531LimitNumericViewActual := B531LimitNumericViewGetter(b531LimitNumericViewVector.B531LimitNumericViewRaw); b531LimitNumericViewActual != b531LimitNumericViewVector.B531LimitNumericViewExpected {
				b531LimitNumericViewT.Fatalf("ab19an19 got %d want %d", b531LimitNumericViewActual, b531LimitNumericViewVector.B531LimitNumericViewExpected)
			}
		})
	}
}

func TestAb19An19Admission(b531LimitNumericViewT *b531LimitNumericViewTesting.T) {
	for _, b531LimitNumericViewCase := range []struct {
		B531LimitNumericViewLabel, B531LimitNumericViewFamily, B531LimitNumericViewVariant, B531LimitNumericViewDirection string
		B531LimitNumericViewOperation                                                                                     byte
	}{
		{"ab19an19-family", "other-family", B531LimitNumericViewVariant, B531LimitNumericViewDirection, 2},
		{"ab19an19-variant", B531LimitNumericViewFamily, "other-variant", B531LimitNumericViewDirection, 2},
		{"ab19an19-direction", B531LimitNumericViewFamily, B531LimitNumericViewVariant, "request", 2},
		{"ab19an19-operation", B531LimitNumericViewFamily, B531LimitNumericViewVariant, B531LimitNumericViewDirection, 3},
	} {
		b531LimitNumericViewT.Run(b531LimitNumericViewCase.B531LimitNumericViewLabel, func(b531LimitNumericViewT *b531LimitNumericViewTesting.T) {
			b531LimitNumericViewResult := B531LimitNumericViewInterpret(b531LimitNumericViewCase.B531LimitNumericViewFamily, b531LimitNumericViewCase.B531LimitNumericViewVariant, b531LimitNumericViewCase.B531LimitNumericViewOperation, b531LimitNumericViewCase.B531LimitNumericViewDirection, []byte{1, 0}, []byte{2, 0}, true, true, B531LimitNumericViewPair{})
			if b531LimitNumericViewResult.B531LimitNumericViewAdmitted || b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewFresh || b531LimitNumericViewResult.B531LimitNumericViewLastKnown.B531LimitNumericViewFresh || b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewMinimum.B531LimitNumericViewAvailable || b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewMaximum.B531LimitNumericViewAvailable || b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewMinimum.B531LimitNumericViewState != "not-admitted" || b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewMaximum.B531LimitNumericViewState != "not-admitted" {
				b531LimitNumericViewT.Fatalf("ab19an19 admission: %+v", b531LimitNumericViewResult)
			}
		})
	}
}

func TestAb19An19Retention(b531LimitNumericViewT *b531LimitNumericViewTesting.T) {
	for _, b531LimitNumericViewSide := range []string{"ab19an19-minimum", "ab19an19-maximum"} {
		for _, b531LimitNumericViewCase := range []struct {
			B531LimitNumericViewLabel, B531LimitNumericViewRaw, B531LimitNumericViewState string
			B531LimitNumericViewReadOK                                                    bool
		}{
			{"ab19an19-failed", "0100", "read-failed", false},
			{"ab19an19-missing", "", "missing", true},
			{"ab19an19-short", "01", "malformed", true},
			{"ab19an19-long", "010000", "malformed", true},
		} {
			for _, b531LimitNumericViewPrior := range []B531LimitNumericViewField{
				{},
				{B531LimitNumericViewRaw: []byte{3, 0}, B531LimitNumericViewState: "success", B531LimitNumericViewAvailable: true, B531LimitNumericViewValue: 30, B531LimitNumericViewFresh: true},
				{B531LimitNumericViewRaw: []byte{4, 0}, B531LimitNumericViewState: "success", B531LimitNumericViewAvailable: true, B531LimitNumericViewValue: 40},
				{B531LimitNumericViewRaw: []byte{9}, B531LimitNumericViewState: "malformed", B531LimitNumericViewValue: 77},
			} {
				b531LimitNumericViewT.Run(b531LimitNumericViewSide+"/"+b531LimitNumericViewCase.B531LimitNumericViewLabel, func(b531LimitNumericViewT *b531LimitNumericViewTesting.T) {
					b531LimitNumericViewPair := B531LimitNumericViewPair{B531LimitNumericViewMinimum: b531LimitNumericViewPrior, B531LimitNumericViewMaximum: b531LimitNumericViewPrior, B531LimitNumericViewFresh: true}
					b531LimitNumericViewBadRaw := b531LimitNumericViewRaw(b531LimitNumericViewT, b531LimitNumericViewCase.B531LimitNumericViewRaw)
					b531LimitNumericViewMin, b531LimitNumericViewMax := []byte{5, 0}, []byte{5, 0}
					b531LimitNumericViewMinOK, b531LimitNumericViewMaxOK := true, true
					if b531LimitNumericViewSide == "ab19an19-minimum" {
						b531LimitNumericViewMin, b531LimitNumericViewMinOK = b531LimitNumericViewBadRaw, b531LimitNumericViewCase.B531LimitNumericViewReadOK
					} else {
						b531LimitNumericViewMax, b531LimitNumericViewMaxOK = b531LimitNumericViewBadRaw, b531LimitNumericViewCase.B531LimitNumericViewReadOK
					}
					b531LimitNumericViewResult := B531LimitNumericViewInterpret(B531LimitNumericViewFamily, B531LimitNumericViewVariant, 2, B531LimitNumericViewDirection, b531LimitNumericViewMin, b531LimitNumericViewMax, b531LimitNumericViewMinOK, b531LimitNumericViewMaxOK, b531LimitNumericViewPair)
					b531LimitNumericViewFailed, b531LimitNumericViewRetained, b531LimitNumericViewGood := b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewMinimum, b531LimitNumericViewResult.B531LimitNumericViewLastKnown.B531LimitNumericViewMinimum, b531LimitNumericViewResult.B531LimitNumericViewLastKnown.B531LimitNumericViewMaximum
					if b531LimitNumericViewSide == "ab19an19-maximum" {
						b531LimitNumericViewFailed, b531LimitNumericViewRetained, b531LimitNumericViewGood = b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewMaximum, b531LimitNumericViewResult.B531LimitNumericViewLastKnown.B531LimitNumericViewMaximum, b531LimitNumericViewResult.B531LimitNumericViewLastKnown.B531LimitNumericViewMinimum
					}
					b531LimitNumericViewExpectedFailed := B531LimitNumericViewField{B531LimitNumericViewRaw: append([]byte(nil), b531LimitNumericViewBadRaw...), B531LimitNumericViewState: b531LimitNumericViewCase.B531LimitNumericViewState}
					b531LimitNumericViewExpectedRetained := b531LimitNumericViewPrior
					b531LimitNumericViewExpectedRetained.B531LimitNumericViewFresh = false
					b531LimitNumericViewExpectedGood := B531LimitNumericViewField{B531LimitNumericViewRaw: []byte{5, 0}, B531LimitNumericViewState: "success", B531LimitNumericViewAvailable: true, B531LimitNumericViewValue: 50, B531LimitNumericViewFresh: true}
					if !b531LimitNumericViewReflect.DeepEqual(b531LimitNumericViewFailed, b531LimitNumericViewExpectedFailed) || !b531LimitNumericViewReflect.DeepEqual(b531LimitNumericViewRetained, b531LimitNumericViewExpectedRetained) || !b531LimitNumericViewReflect.DeepEqual(b531LimitNumericViewGood, b531LimitNumericViewExpectedGood) || !b531LimitNumericViewReflect.DeepEqual(b531LimitNumericViewResult.B531LimitNumericViewPrior, b531LimitNumericViewPair) || b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewFresh || b531LimitNumericViewResult.B531LimitNumericViewLastKnown.B531LimitNumericViewFresh {
						b531LimitNumericViewT.Fatalf("ab19an19 partial: %+v", b531LimitNumericViewResult)
					}
				})
			}
		}
	}
}

func TestAb19An19PriorDimensions(b531LimitNumericViewT *b531LimitNumericViewTesting.T) {
	for _, b531LimitNumericViewCase := range []struct {
		B531LimitNumericViewLabel string
		B531LimitNumericViewField B531LimitNumericViewField
	}{
		{"ab19an19-baseline", B531LimitNumericViewField{[]byte{3, 0}, "success", true, 30, true}},
		{"ab19an19-raw", B531LimitNumericViewField{[]byte{4, 0}, "success", true, 30, true}},
		{"ab19an19-state", B531LimitNumericViewField{[]byte{3, 0}, "malformed", true, 30, true}},
		{"ab19an19-available", B531LimitNumericViewField{[]byte{3, 0}, "success", false, 30, true}},
		{"ab19an19-value", B531LimitNumericViewField{[]byte{3, 0}, "success", true, 40, true}},
		{"ab19an19-fresh", B531LimitNumericViewField{[]byte{3, 0}, "success", true, 30, false}},
	} {
		for _, b531LimitNumericViewSide := range []string{"ab19an19-minimum", "ab19an19-maximum"} {
			for _, b531LimitNumericViewFresh := range []bool{false, true} {
				b531LimitNumericViewT.Run(b531LimitNumericViewCase.B531LimitNumericViewLabel+"/"+b531LimitNumericViewSide, func(b531LimitNumericViewT *b531LimitNumericViewTesting.T) {
					b531LimitNumericViewPrior := B531LimitNumericViewPair{B531LimitNumericViewMinimum: B531LimitNumericViewField{[]byte{3, 0}, "success", true, 30, true}, B531LimitNumericViewMaximum: B531LimitNumericViewField{[]byte{3, 0}, "success", true, 30, true}, B531LimitNumericViewFresh: b531LimitNumericViewFresh}
					if b531LimitNumericViewSide == "ab19an19-minimum" {
						b531LimitNumericViewPrior.B531LimitNumericViewMinimum = b531LimitNumericViewCase.B531LimitNumericViewField
					} else {
						b531LimitNumericViewPrior.B531LimitNumericViewMaximum = b531LimitNumericViewCase.B531LimitNumericViewField
					}
					b531LimitNumericViewExpected := b531LimitNumericViewPrior
					b531LimitNumericViewExpected.B531LimitNumericViewFresh = false
					b531LimitNumericViewExpected.B531LimitNumericViewMinimum.B531LimitNumericViewFresh = false
					b531LimitNumericViewExpected.B531LimitNumericViewMaximum.B531LimitNumericViewFresh = false
					b531LimitNumericViewResult := B531LimitNumericViewInterpret(B531LimitNumericViewFamily, B531LimitNumericViewVariant, 2, B531LimitNumericViewDirection, nil, nil, false, false, b531LimitNumericViewPrior)
					if !b531LimitNumericViewReflect.DeepEqual(b531LimitNumericViewResult.B531LimitNumericViewPrior, b531LimitNumericViewPrior) || !b531LimitNumericViewReflect.DeepEqual(b531LimitNumericViewResult.B531LimitNumericViewLastKnown, b531LimitNumericViewExpected) {
						b531LimitNumericViewT.Fatalf("ab19an19 prior dimensions: %+v", b531LimitNumericViewResult)
					}
				})
			}
		}
	}
}

func TestAb19An19TotalFailureAndOwnership(b531LimitNumericViewT *b531LimitNumericViewTesting.T) {
	b531LimitNumericViewMin, b531LimitNumericViewMax := []byte{1, 0}, []byte{2, 0}
	b531LimitNumericViewPrior := B531LimitNumericViewPair{B531LimitNumericViewMinimum: B531LimitNumericViewField{B531LimitNumericViewRaw: []byte{3, 0}, B531LimitNumericViewState: "success", B531LimitNumericViewAvailable: true, B531LimitNumericViewValue: 30, B531LimitNumericViewFresh: true}, B531LimitNumericViewMaximum: B531LimitNumericViewField{B531LimitNumericViewRaw: []byte{4, 0}, B531LimitNumericViewState: "success", B531LimitNumericViewAvailable: true, B531LimitNumericViewValue: 40, B531LimitNumericViewFresh: true}, B531LimitNumericViewFresh: true}
	b531LimitNumericViewFirst := B531LimitNumericViewInterpret(B531LimitNumericViewFamily, B531LimitNumericViewVariant, 2, B531LimitNumericViewDirection, b531LimitNumericViewMin, b531LimitNumericViewMax, false, false, b531LimitNumericViewPrior)
	b531LimitNumericViewSecond := B531LimitNumericViewInterpret(B531LimitNumericViewFamily, B531LimitNumericViewVariant, 2, B531LimitNumericViewDirection, b531LimitNumericViewMin, b531LimitNumericViewMax, false, false, b531LimitNumericViewPrior)
	if b531LimitNumericViewFirst.B531LimitNumericViewLastKnown.B531LimitNumericViewMinimum.B531LimitNumericViewValue != 30 || b531LimitNumericViewFirst.B531LimitNumericViewLastKnown.B531LimitNumericViewMaximum.B531LimitNumericViewValue != 40 || b531LimitNumericViewFirst.B531LimitNumericViewLastKnown.B531LimitNumericViewFresh || b531LimitNumericViewFirst.B531LimitNumericViewLastKnown.B531LimitNumericViewMinimum.B531LimitNumericViewFresh || b531LimitNumericViewFirst.B531LimitNumericViewLastKnown.B531LimitNumericViewMaximum.B531LimitNumericViewFresh || b531LimitNumericViewFirst.B531LimitNumericViewAttempt.B531LimitNumericViewMinimum.B531LimitNumericViewState != "read-failed" || b531LimitNumericViewFirst.B531LimitNumericViewAttempt.B531LimitNumericViewMaximum.B531LimitNumericViewState != "read-failed" {
		b531LimitNumericViewT.Fatalf("ab19an19 total failure: %+v", b531LimitNumericViewFirst)
	}
	b531LimitNumericViewMin[0], b531LimitNumericViewMax[0] = 99, 99
	b531LimitNumericViewPrior.B531LimitNumericViewMinimum.B531LimitNumericViewRaw[0], b531LimitNumericViewPrior.B531LimitNumericViewMaximum.B531LimitNumericViewRaw[0] = 99, 99
	b531LimitNumericViewFirst.B531LimitNumericViewAttempt.B531LimitNumericViewMinimum.B531LimitNumericViewRaw[0] = 88
	b531LimitNumericViewFirst.B531LimitNumericViewAttempt.B531LimitNumericViewMaximum.B531LimitNumericViewRaw[0] = 88
	b531LimitNumericViewFirst.B531LimitNumericViewPrior.B531LimitNumericViewMinimum.B531LimitNumericViewRaw[0] = 88
	b531LimitNumericViewFirst.B531LimitNumericViewPrior.B531LimitNumericViewMaximum.B531LimitNumericViewRaw[0] = 88
	if b531LimitNumericViewSecond.B531LimitNumericViewAttempt.B531LimitNumericViewMinimum.B531LimitNumericViewRaw[0] != 1 || b531LimitNumericViewSecond.B531LimitNumericViewAttempt.B531LimitNumericViewMaximum.B531LimitNumericViewRaw[0] != 2 || b531LimitNumericViewSecond.B531LimitNumericViewPrior.B531LimitNumericViewMinimum.B531LimitNumericViewRaw[0] != 3 || b531LimitNumericViewSecond.B531LimitNumericViewPrior.B531LimitNumericViewMaximum.B531LimitNumericViewRaw[0] != 4 || b531LimitNumericViewFirst.B531LimitNumericViewLastKnown.B531LimitNumericViewMinimum.B531LimitNumericViewRaw[0] != 3 || b531LimitNumericViewFirst.B531LimitNumericViewLastKnown.B531LimitNumericViewMaximum.B531LimitNumericViewRaw[0] != 4 {
		b531LimitNumericViewT.Fatal("ab19an19 raw buffers alias")
	}
	b531LimitNumericViewFirst = B531LimitNumericViewInterpret(B531LimitNumericViewFamily, B531LimitNumericViewVariant, 2, B531LimitNumericViewDirection, []byte{1, 0}, []byte{2, 0}, true, true, B531LimitNumericViewPair{})
	b531LimitNumericViewFirst.B531LimitNumericViewAttempt.B531LimitNumericViewMinimum.B531LimitNumericViewRaw[0], b531LimitNumericViewFirst.B531LimitNumericViewAttempt.B531LimitNumericViewMaximum.B531LimitNumericViewRaw[0] = 88, 88
	if b531LimitNumericViewFirst.B531LimitNumericViewLastKnown.B531LimitNumericViewMinimum.B531LimitNumericViewRaw[0] != 1 || b531LimitNumericViewFirst.B531LimitNumericViewLastKnown.B531LimitNumericViewMaximum.B531LimitNumericViewRaw[0] != 2 {
		b531LimitNumericViewT.Fatal("ab19an19 successful buffers alias")
	}
}
