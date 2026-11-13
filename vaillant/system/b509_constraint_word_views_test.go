package system

import (
	b509ConstraintWordViewBytes "bytes"
	b509ConstraintWordViewJSON "encoding/json"
	b509ConstraintWordViewOs "os"
	b509ConstraintWordViewReflect "reflect"
	b509ConstraintWordViewTesting "testing"
)

type b509ConstraintWordViewVector struct {
	B509ConstraintWordViewLabel    string                      `json:"Ab07An07Label"`
	B509ConstraintWordViewInput    B509ConstraintWordViewInput `json:"Ab07An07Input"`
	B509ConstraintWordViewValues   [5]int16                    `json:"Ab07An07Values"`
	B509ConstraintWordViewStates   [5]string                   `json:"Ab07An07States"`
	B509ConstraintWordViewSelected bool                        `json:"Ab07An07Selected"`
	B509ConstraintWordViewEqual    bool                        `json:"Ab07An07Equal"`
}

func TestAb07An07Vectors(b509ConstraintWordViewT *b509ConstraintWordViewTesting.T) {
	b509ConstraintWordViewData, b509ConstraintWordViewErr := b509ConstraintWordViewOs.ReadFile("testdata/ab07an07_b509_constraint_word_views.json")
	if b509ConstraintWordViewErr != nil {
		b509ConstraintWordViewT.Fatal(b509ConstraintWordViewErr)
	}
	var b509ConstraintWordViewVectors []b509ConstraintWordViewVector
	if b509ConstraintWordViewErr = b509ConstraintWordViewJSON.Unmarshal(b509ConstraintWordViewData, &b509ConstraintWordViewVectors); b509ConstraintWordViewErr != nil {
		b509ConstraintWordViewT.Fatal(b509ConstraintWordViewErr)
	}
	if len(b509ConstraintWordViewVectors) == 0 {
		b509ConstraintWordViewT.Fatal("ab07an07_no_vectors")
	}
	for _, b509ConstraintWordViewVector := range b509ConstraintWordViewVectors {
		b509ConstraintWordViewT.Run(b509ConstraintWordViewVector.B509ConstraintWordViewLabel, func(b509ConstraintWordViewT *b509ConstraintWordViewTesting.T) {
			b509ConstraintWordViewView := B509ConstraintWordViewInterpret(b509ConstraintWordViewVector.B509ConstraintWordViewInput)
			if !b509ConstraintWordViewReflect.DeepEqual(b509ConstraintWordViewView.B509ConstraintWordViewContext, b509ConstraintWordViewVector.B509ConstraintWordViewInput) {
				b509ConstraintWordViewT.Fatal("ab07an07_context_lost")
			}
			if b509ConstraintWordViewView.B509ConstraintWordViewSelected != b509ConstraintWordViewVector.B509ConstraintWordViewSelected || b509ConstraintWordViewView.B509ConstraintWordViewRawAddressEqual != b509ConstraintWordViewVector.B509ConstraintWordViewEqual {
				b509ConstraintWordViewT.Fatal("ab07an07_selection_or_local_equality")
			}
			if b509ConstraintWordViewView.B509ConstraintWordViewRequestedEditBoundsAdmitted {
				b509ConstraintWordViewT.Fatal("ab07an07_candidate_cannot_admit_bounds")
			}
			b509ConstraintWordViewWords := []B509ConstraintWord{b509ConstraintWordViewView.B509ConstraintWordViewMinimum, b509ConstraintWordViewView.B509ConstraintWordViewMaximum, b509ConstraintWordViewView.B509ConstraintWordViewDefault, b509ConstraintWordViewView.B509ConstraintWordViewStep, b509ConstraintWordViewView.B509ConstraintWordViewFactor}
			b509ConstraintWordViewRaw := [][]byte{b509ConstraintWordViewVector.B509ConstraintWordViewInput.B509ConstraintWordViewMinimum, b509ConstraintWordViewVector.B509ConstraintWordViewInput.B509ConstraintWordViewMaximum, b509ConstraintWordViewVector.B509ConstraintWordViewInput.B509ConstraintWordViewDefault, b509ConstraintWordViewVector.B509ConstraintWordViewInput.B509ConstraintWordViewStep, b509ConstraintWordViewVector.B509ConstraintWordViewInput.B509ConstraintWordViewFactor}
			for b509ConstraintWordViewIndex, b509ConstraintWordViewWord := range b509ConstraintWordViewWords {
				if b509ConstraintWordViewWord.B509ConstraintWordViewValue != int(b509ConstraintWordViewVector.B509ConstraintWordViewValues[b509ConstraintWordViewIndex]) || b509ConstraintWordViewWord.B509ConstraintWordViewState != b509ConstraintWordViewVector.B509ConstraintWordViewStates[b509ConstraintWordViewIndex] || !b509ConstraintWordViewReflect.DeepEqual(b509ConstraintWordViewWord.B509ConstraintWordViewRaw, b509ConstraintWordViewRaw[b509ConstraintWordViewIndex]) {
					b509ConstraintWordViewT.Fatalf("ab07an07_word_%d: %+v", b509ConstraintWordViewIndex, b509ConstraintWordViewWord)
				}
			}
			b509ConstraintWordViewGrades := map[string]string{
				"B509ConstraintWordViewSoftwareContract": "Proven", "B509ConstraintWordViewInterpretation": "Hypothesis",
				"B509ConstraintWordViewNativeEchoRules": "Unknown", "B509ConstraintWordViewPhysicalQualification": "Unknown",
				"B509ConstraintWordViewPacketAssociation": "Unknown", "B509ConstraintWordViewFamilyApplicability": "Unknown",
				"B509ConstraintWordViewFactorAlgorithm": "Unknown", "B509ConstraintWordViewWritability": "Unknown",
				"B509ConstraintWordViewPeerEnforcement": "Unknown", "B509ConstraintWordViewFreshness": "Unknown", "B509ConstraintWordViewExecution": "Unknown",
			}
			for b509ConstraintWordViewName, b509ConstraintWordViewExpected := range b509ConstraintWordViewGrades {
				if b509ConstraintWordViewReflect.ValueOf(b509ConstraintWordViewView).FieldByName(b509ConstraintWordViewName).String() != b509ConstraintWordViewExpected {
					b509ConstraintWordViewT.Fatal(b509ConstraintWordViewName)
				}
			}
		})
	}
}

func TestAb07An07DetachedMemory(b509ConstraintWordViewT *b509ConstraintWordViewTesting.T) {
	b509ConstraintWordViewNames := []string{"B509ConstraintWordViewRequestedAddress", "B509ConstraintWordViewReturnedAddress", "B509ConstraintWordViewMinimum", "B509ConstraintWordViewMaximum", "B509ConstraintWordViewDefault", "B509ConstraintWordViewStep", "B509ConstraintWordViewFactor"}
	for _, b509ConstraintWordViewName := range b509ConstraintWordViewNames {
		b509ConstraintWordViewT.Run(b509ConstraintWordViewName, func(b509ConstraintWordViewT *b509ConstraintWordViewTesting.T) {
			b509ConstraintWordViewInput := B509ConstraintWordViewInput{B509ConstraintWordViewOperation: 0x33, B509ConstraintWordViewDirection: B509ConstraintWordViewPositiveResponse}
			b509ConstraintWordViewSource := []byte{1, 2}
			b509ConstraintWordViewReflect.ValueOf(&b509ConstraintWordViewInput).Elem().FieldByName(b509ConstraintWordViewName).SetBytes(b509ConstraintWordViewSource)
			b509ConstraintWordViewView := B509ConstraintWordViewInterpret(b509ConstraintWordViewInput)
			b509ConstraintWordViewOther := B509ConstraintWordViewInterpret(b509ConstraintWordViewInput)
			b509ConstraintWordViewContextRaw := b509ConstraintWordViewReflect.ValueOf(b509ConstraintWordViewView.B509ConstraintWordViewContext).FieldByName(b509ConstraintWordViewName).Bytes()
			b509ConstraintWordViewSource[0] = 9
			if !b509ConstraintWordViewBytes.Equal(b509ConstraintWordViewContextRaw, []byte{1, 2}) {
				b509ConstraintWordViewT.Fatal("ab07an07_input_alias")
			}
			b509ConstraintWordViewContextRaw[1] = 8
			if b509ConstraintWordViewSource[1] != 2 || b509ConstraintWordViewReflect.ValueOf(b509ConstraintWordViewOther.B509ConstraintWordViewContext).FieldByName(b509ConstraintWordViewName).Bytes()[1] != 2 {
				b509ConstraintWordViewT.Fatal("ab07an07_result_alias")
			}
			if b509ConstraintWordViewName != "B509ConstraintWordViewRequestedAddress" && b509ConstraintWordViewName != "B509ConstraintWordViewReturnedAddress" {
				b509ConstraintWordViewWord := b509ConstraintWordViewReflect.ValueOf(b509ConstraintWordViewView).FieldByName(b509ConstraintWordViewName).Interface().(B509ConstraintWord)
				if !b509ConstraintWordViewBytes.Equal(b509ConstraintWordViewWord.B509ConstraintWordViewRaw, []byte{1, 2}) {
					b509ConstraintWordViewT.Fatal("ab07an07_word_alias")
				}
				b509ConstraintWordViewWord.B509ConstraintWordViewRaw[0] = 7
				if b509ConstraintWordViewSource[0] != 9 || b509ConstraintWordViewContextRaw[0] != 1 || b509ConstraintWordViewReflect.ValueOf(b509ConstraintWordViewOther).FieldByName(b509ConstraintWordViewName).Interface().(B509ConstraintWord).B509ConstraintWordViewRaw[0] != 1 {
					b509ConstraintWordViewT.Fatal("ab07an07_word_output_alias")
				}
			}
		})
	}
}
