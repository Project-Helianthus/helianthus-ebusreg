package ebusstandard

import (
	standard05CandidateViewJSON "encoding/json"
	standard05CandidateViewOs "os"
	standard05CandidateViewReflect "reflect"
	standard05CandidateViewTesting "testing"
)

func standard05CandidateViewByte(standard05CandidateViewValue uint8) *uint8 {
	return &standard05CandidateViewValue
}

func standard05CandidateViewControl() Standard05CandidateViewKey {
	return Standard05CandidateViewKey{Standard05CandidateViewFamily: "burner", Standard05CandidateViewVariant: "0501", Standard05CandidateViewPrimary: standard05CandidateViewByte(5), Standard05CandidateViewSecondary: standard05CandidateViewByte(1), Standard05CandidateViewDirection: "abstract-request", Standard05CandidateViewSelector: "burner-five-field"}
}

func TestAb78Examples(standard05CandidateViewT *standard05CandidateViewTesting.T) {
	standard05CandidateViewData, standard05CandidateViewErr := standard05CandidateViewOs.ReadFile("testdata/ab78_standard05_candidate_views.json")
	if standard05CandidateViewErr != nil {
		standard05CandidateViewT.Fatal(standard05CandidateViewErr)
	}
	var standard05CandidateViewCases []struct {
		Standard05CandidateViewLabel   string                         `json:"ab78Label"`
		Standard05CandidateViewBody    []int                          `json:"ab78Body"`
		Standard05CandidateViewValues  []*int                         `json:"ab78Values"`
		Standard05CandidateViewName    string                         `json:"ab78Name"`
		Standard05CandidateViewDomains []Standard05CandidateViewState `json:"ab78Domains"`
	}
	if standard05CandidateViewErr := standard05CandidateViewJSON.Unmarshal(standard05CandidateViewData, &standard05CandidateViewCases); standard05CandidateViewErr != nil {
		standard05CandidateViewT.Fatal(standard05CandidateViewErr)
	}
	if len(standard05CandidateViewCases) != 7 {
		standard05CandidateViewT.Fatal("ab78 fixture inventory")
	}
	for _, standard05CandidateViewCase := range standard05CandidateViewCases {
		standard05CandidateViewT.Run(standard05CandidateViewCase.Standard05CandidateViewLabel, func(standard05CandidateViewT *standard05CandidateViewTesting.T) {
			standard05CandidateViewBody := make([]byte, len(standard05CandidateViewCase.Standard05CandidateViewBody))
			for standard05CandidateViewIndex, standard05CandidateViewValue := range standard05CandidateViewCase.Standard05CandidateViewBody {
				standard05CandidateViewBody[standard05CandidateViewIndex] = byte(standard05CandidateViewValue)
			}
			standard05CandidateViewView := Standard05CandidateViewInterpret(standard05CandidateViewControl(), standard05CandidateViewBody)
			if !standard05CandidateViewView.Standard05CandidateViewComplete || standard05CandidateViewView.Standard05CandidateViewState != "candidate-body" || len(standard05CandidateViewView.Standard05CandidateViewRemainder) != 0 || !standard05CandidateViewReflect.DeepEqual(standard05CandidateViewView.Standard05CandidateViewRaw, standard05CandidateViewBody) {
				standard05CandidateViewT.Fatal("ab78 body boundary")
			}
			for standard05CandidateViewIndex, standard05CandidateViewField := range standard05CandidateViewView.Standard05CandidateViewFields {
				standard05CandidateViewWantState := Standard05CandidateViewState("interpreted")
				if standard05CandidateViewCase.Standard05CandidateViewValues[standard05CandidateViewIndex] == nil {
					standard05CandidateViewWantState = "unmapped"
				}
				if standard05CandidateViewField.Standard05CandidateViewRaw == nil || *standard05CandidateViewField.Standard05CandidateViewRaw != standard05CandidateViewBody[standard05CandidateViewIndex] || !standard05CandidateViewReflect.DeepEqual(standard05CandidateViewField.Standard05CandidateViewValue, standard05CandidateViewCase.Standard05CandidateViewValues[standard05CandidateViewIndex]) || standard05CandidateViewField.Standard05CandidateViewDomain != standard05CandidateViewCase.Standard05CandidateViewDomains[standard05CandidateViewIndex] || standard05CandidateViewField.Standard05CandidateViewInterpretation != standard05CandidateViewWantState {
					standard05CandidateViewT.Fatalf("ab78 field %d: %+v", standard05CandidateViewIndex, standard05CandidateViewField)
				}
			}
			if standard05CandidateViewView.Standard05CandidateViewFields[0].Standard05CandidateViewName != standard05CandidateViewCase.Standard05CandidateViewName {
				standard05CandidateViewT.Fatal("ab78 heat name")
			}
		})
	}
}

func TestAb78SelectionAxes(standard05CandidateViewT *standard05CandidateViewTesting.T) {
	standard05CandidateViewCases := []struct {
		Standard05CandidateViewLabel  string
		Standard05CandidateViewMutate func(*Standard05CandidateViewKey)
	}{
		{"ab78-family", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewFamily = "unknown"
		}},
		{"ab78-variant-unknown", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewVariant = "unknown"
		}},
		{"ab78-variant", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewVariant = "0504"
		}},
		{"ab78-primary", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewPrimary = standard05CandidateViewByte(4)
		}},
		{"ab78-primary-absent", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewPrimary = nil
		}},
		{"ab78-primary-zero", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewPrimary = standard05CandidateViewByte(0)
		}},
		{"ab78-secondary", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewSecondary = standard05CandidateViewByte(4)
		}},
		{"ab78-secondary-absent", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewSecondary = nil
		}},
		{"ab78-secondary-zero", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewSecondary = standard05CandidateViewByte(0)
		}},
		{"ab78-operation-zero", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewOperation = standard05CandidateViewByte(0)
		}},
		{"ab78-operation", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewOperation = standard05CandidateViewByte(1)
		}},
		{"ab78-positive", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewDirection = "abstract-positive"
		}},
		{"ab78-direction", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewDirection = "unknown"
		}},
		{"ab78-selector", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewSelector = "unknown"
		}},
		{"ab78-selector-absent", func(standard05CandidateViewKey *Standard05CandidateViewKey) {
			standard05CandidateViewKey.Standard05CandidateViewSelector = ""
		}},
	}
	for _, standard05CandidateViewCase := range standard05CandidateViewCases {
		standard05CandidateViewT.Run(standard05CandidateViewCase.Standard05CandidateViewLabel, func(standard05CandidateViewT *standard05CandidateViewTesting.T) {
			standard05CandidateViewKey := standard05CandidateViewControl()
			if Standard05CandidateViewInterpret(standard05CandidateViewKey, []byte{0, 0, 0, 0, 0}).Standard05CandidateViewState != "candidate-body" {
				standard05CandidateViewT.Fatal("ab78 valid control")
			}
			standard05CandidateViewCase.Standard05CandidateViewMutate(&standard05CandidateViewKey)
			for _, standard05CandidateViewBody := range [][]byte{nil, {5, 1}, {170, 90, 35, 255, 255}, {5, 1, 170, 90, 35, 255, 255}} {
				standard05CandidateViewView := Standard05CandidateViewInterpret(standard05CandidateViewKey, standard05CandidateViewBody)
				if standard05CandidateViewView.Standard05CandidateViewState != "raw-only" || standard05CandidateViewView.Standard05CandidateViewComplete || standard05CandidateViewView.Standard05CandidateViewFields != [5]Standard05CandidateViewField{} || !standard05CandidateViewReflect.DeepEqual(standard05CandidateViewView.Standard05CandidateViewRaw, standard05CandidateViewBody) || !standard05CandidateViewReflect.DeepEqual(standard05CandidateViewView.Standard05CandidateViewRemainder, standard05CandidateViewBody) {
					standard05CandidateViewT.Fatal("ab78 changed axis interpreted")
				}
			}
		})
	}
	standard05CandidateViewKey := standard05CandidateViewControl()
	standard05CandidateViewKey.Standard05CandidateViewVariant = "0504"
	standard05CandidateViewKey.Standard05CandidateViewSecondary = standard05CandidateViewByte(4)
	standard05CandidateViewView := Standard05CandidateViewInterpret(standard05CandidateViewKey, []byte{0, 0, 0, 0, 0})
	if standard05CandidateViewView.Standard05CandidateViewState != "raw-only" || standard05CandidateViewView.Standard05CandidateViewKey.Standard05CandidateViewVariant != "0504" || *standard05CandidateViewView.Standard05CandidateViewKey.Standard05CandidateViewSecondary != 4 {
		standard05CandidateViewT.Fatal("ab78 separate identity lost")
	}
}

func TestAb78Lengths(standard05CandidateViewT *standard05CandidateViewTesting.T) {
	standard05CandidateViewBody := []byte{170, 90, 35, 255, 255, 17, 34}
	for standard05CandidateViewLength := 0; standard05CandidateViewLength <= len(standard05CandidateViewBody); standard05CandidateViewLength++ {
		standard05CandidateViewView := Standard05CandidateViewInterpret(standard05CandidateViewControl(), standard05CandidateViewBody[:standard05CandidateViewLength])
		if standard05CandidateViewView.Standard05CandidateViewComplete != (standard05CandidateViewLength == 5) {
			standard05CandidateViewT.Fatalf("ab78 completeness %d", standard05CandidateViewLength)
		}
		for standard05CandidateViewIndex, standard05CandidateViewField := range standard05CandidateViewView.Standard05CandidateViewFields {
			if standard05CandidateViewIndex >= standard05CandidateViewLength {
				if standard05CandidateViewField.Standard05CandidateViewRaw != nil || standard05CandidateViewField.Standard05CandidateViewValue != nil || standard05CandidateViewField.Standard05CandidateViewName != "" || standard05CandidateViewField.Standard05CandidateViewInterpretation != "absent" || standard05CandidateViewField.Standard05CandidateViewDomain != "unknown" {
					standard05CandidateViewT.Fatalf("ab78 filled absent %d", standard05CandidateViewIndex)
				}
			} else if standard05CandidateViewField.Standard05CandidateViewRaw == nil || *standard05CandidateViewField.Standard05CandidateViewRaw != standard05CandidateViewBody[standard05CandidateViewIndex] || standard05CandidateViewField.Standard05CandidateViewValue == nil || *standard05CandidateViewField.Standard05CandidateViewValue != []int{170, 90, 35, -1, 255}[standard05CandidateViewIndex] {
				standard05CandidateViewT.Fatal("ab78 partial value lost")
			}
		}
		if standard05CandidateViewLength > 5 && !standard05CandidateViewReflect.DeepEqual(standard05CandidateViewView.Standard05CandidateViewRemainder, standard05CandidateViewBody[5:standard05CandidateViewLength]) {
			standard05CandidateViewT.Fatal("ab78 remainder lost")
		}
	}
	standard05CandidateViewZero := Standard05CandidateViewInterpret(standard05CandidateViewControl(), []byte{0, 0, 0, 0, 0})
	for _, standard05CandidateViewField := range standard05CandidateViewZero.Standard05CandidateViewFields {
		if standard05CandidateViewField.Standard05CandidateViewRaw == nil || standard05CandidateViewField.Standard05CandidateViewValue == nil || *standard05CandidateViewField.Standard05CandidateViewValue != 0 || standard05CandidateViewField.Standard05CandidateViewInterpretation != "interpreted" {
			standard05CandidateViewT.Fatal("ab78 explicit zero absent")
		}
	}
}

func TestAb78FieldAxes(standard05CandidateViewT *standard05CandidateViewTesting.T) {
	for standard05CandidateViewIndex := 0; standard05CandidateViewIndex < 5; standard05CandidateViewIndex++ {
		for standard05CandidateViewValue := 0; standard05CandidateViewValue <= 255; standard05CandidateViewValue++ {
			standard05CandidateViewBody := []byte{0, 0, 0, 0, 0}
			standard05CandidateViewBody[standard05CandidateViewIndex] = byte(standard05CandidateViewValue)
			standard05CandidateViewView := Standard05CandidateViewInterpret(standard05CandidateViewControl(), standard05CandidateViewBody)
			for standard05CandidateViewOther, standard05CandidateViewField := range standard05CandidateViewView.Standard05CandidateViewFields {
				if standard05CandidateViewOther != standard05CandidateViewIndex && (standard05CandidateViewField.Standard05CandidateViewValue == nil || *standard05CandidateViewField.Standard05CandidateViewValue != 0) {
					standard05CandidateViewT.Fatal("ab78 field interference")
				}
			}
			standard05CandidateViewField := standard05CandidateViewView.Standard05CandidateViewFields[standard05CandidateViewIndex]
			standard05CandidateViewWant := standard05CandidateViewValue
			if standard05CandidateViewIndex == 3 && standard05CandidateViewValue >= 128 {
				standard05CandidateViewWant = standard05CandidateViewValue - 256
			}
			if standard05CandidateViewIndex == 0 {
				standard05CandidateViewNames := map[int]string{0: "shutdown", 85: "DHW", 170: "heating", 204: "emission-check", 238: "controller-stop"}
				standard05CandidateViewName, standard05CandidateViewKnown := standard05CandidateViewNames[standard05CandidateViewValue]
				if !standard05CandidateViewKnown {
					if standard05CandidateViewField.Standard05CandidateViewInterpretation != "unmapped" || standard05CandidateViewField.Standard05CandidateViewValue != nil || standard05CandidateViewField.Standard05CandidateViewName != "" || standard05CandidateViewField.Standard05CandidateViewDomain != "unknown" {
						standard05CandidateViewT.Fatal("ab78 unmapped heat named")
					}
					continue
				}
				if standard05CandidateViewField.Standard05CandidateViewName != standard05CandidateViewName {
					standard05CandidateViewT.Fatal("ab78 heat mapping")
				}
			}
			standard05CandidateViewDomain := Standard05CandidateViewState("within-candidate-domain")
			if (standard05CandidateViewIndex == 1 || standard05CandidateViewIndex == 2) && standard05CandidateViewValue > 90 {
				standard05CandidateViewDomain = "out-of-candidate-domain"
			}
			if standard05CandidateViewField.Standard05CandidateViewRaw == nil || int(*standard05CandidateViewField.Standard05CandidateViewRaw) != standard05CandidateViewValue || standard05CandidateViewField.Standard05CandidateViewValue == nil || *standard05CandidateViewField.Standard05CandidateViewValue != standard05CandidateViewWant || standard05CandidateViewField.Standard05CandidateViewDomain != standard05CandidateViewDomain || standard05CandidateViewField.Standard05CandidateViewInterpretation != "interpreted" {
				standard05CandidateViewT.Fatalf("ab78 axis %d value %d", standard05CandidateViewIndex, standard05CandidateViewValue)
			}
		}
	}
}

func TestAb78Metadata(standard05CandidateViewT *standard05CandidateViewTesting.T) {
	standard05CandidateViewView := Standard05CandidateViewInterpret(standard05CandidateViewControl(), nil)
	if standard05CandidateViewView.Standard05CandidateViewAPIGrade != "Proven" || standard05CandidateViewView.Standard05CandidateViewMeaningGrade != "Hypothesis" || standard05CandidateViewView.Standard05CandidateViewNativeGrade != "Unknown" || !standard05CandidateViewReflect.DeepEqual(standard05CandidateViewView.Standard05CandidateViewUnknowns, []string{"physical support", "presence", "freshness", "physical length", "delivery", "effects", "persistence", "execution", "ACK", "command acceptance", "empty wire response"}) {
		standard05CandidateViewT.Fatal("ab78 grades or unknowns")
	}
	for standard05CandidateViewIndex, standard05CandidateViewField := range standard05CandidateViewView.Standard05CandidateViewFields {
		if standard05CandidateViewField.Standard05CandidateViewLabel != []string{"heat-state", "boiler-target", "DHW-target", "outside-temperature", "setting"}[standard05CandidateViewIndex] || standard05CandidateViewField.Standard05CandidateViewDefault != 0 || standard05CandidateViewField.Standard05CandidateViewMinimum != []int{0, 0, 0, -128, 0}[standard05CandidateViewIndex] || standard05CandidateViewField.Standard05CandidateViewMaximum != []int{255, 90, 90, 127, 255}[standard05CandidateViewIndex] {
			standard05CandidateViewT.Fatal("ab78 metadata")
		}
	}
}

func TestAb78OwnedMemory(standard05CandidateViewT *standard05CandidateViewTesting.T) {
	standard05CandidateViewKey := standard05CandidateViewControl()
	standard05CandidateViewBody := []byte{170, 90, 35, 255, 255, 17}
	standard05CandidateViewBaseline := Standard05CandidateViewInterpret(standard05CandidateViewControl(), append([]byte(nil), standard05CandidateViewBody...))
	standard05CandidateViewView := Standard05CandidateViewInterpret(standard05CandidateViewKey, standard05CandidateViewBody)
	*standard05CandidateViewKey.Standard05CandidateViewPrimary = 0
	*standard05CandidateViewKey.Standard05CandidateViewSecondary = 0
	standard05CandidateViewBody[0] = 0
	standard05CandidateViewBody[5] = 0
	if !standard05CandidateViewReflect.DeepEqual(standard05CandidateViewView, standard05CandidateViewBaseline) {
		standard05CandidateViewT.Fatal("ab78 input alias")
	}
	standard05CandidateViewView.Standard05CandidateViewRaw[0] = 1
	standard05CandidateViewView.Standard05CandidateViewRaw[5] = 2
	if *standard05CandidateViewView.Standard05CandidateViewFields[0].Standard05CandidateViewRaw != 170 || standard05CandidateViewView.Standard05CandidateViewRemainder[0] != 17 {
		standard05CandidateViewT.Fatal("ab78 internal alias")
	}
	for standard05CandidateViewIndex := range standard05CandidateViewView.Standard05CandidateViewFields {
		*standard05CandidateViewView.Standard05CandidateViewFields[standard05CandidateViewIndex].Standard05CandidateViewRaw = 0
		*standard05CandidateViewView.Standard05CandidateViewFields[standard05CandidateViewIndex].Standard05CandidateViewValue = 0
	}
	standard05CandidateViewView.Standard05CandidateViewUnknowns[0] = "changed"
	*standard05CandidateViewView.Standard05CandidateViewKey.Standard05CandidateViewPrimary = 0
	*standard05CandidateViewView.Standard05CandidateViewKey.Standard05CandidateViewSecondary = 0
	if !standard05CandidateViewReflect.DeepEqual(Standard05CandidateViewInterpret(standard05CandidateViewControl(), []byte{170, 90, 35, 255, 255, 17}), standard05CandidateViewBaseline) {
		standard05CandidateViewT.Fatal("ab78 return alias")
	}
	standard05CandidateViewKey = standard05CandidateViewControl()
	standard05CandidateViewKey.Standard05CandidateViewOperation = standard05CandidateViewByte(0)
	standard05CandidateViewRaw := Standard05CandidateViewInterpret(standard05CandidateViewKey, standard05CandidateViewBody)
	*standard05CandidateViewKey.Standard05CandidateViewOperation = 1
	standard05CandidateViewRaw.Standard05CandidateViewRaw[0] = 9
	if *standard05CandidateViewRaw.Standard05CandidateViewKey.Standard05CandidateViewOperation != 0 || standard05CandidateViewRaw.Standard05CandidateViewRemainder[0] != standard05CandidateViewBody[0] {
		standard05CandidateViewT.Fatal("ab78 raw-only alias")
	}
}
