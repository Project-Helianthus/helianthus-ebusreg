package ebusstandard

import (
	standard03CandidateViewJSON "encoding/json"
	standard03CandidateViewFmt "fmt"
	standard03CandidateViewOs "os"
	standard03CandidateViewReflect "reflect"
	standard03CandidateViewTesting "testing"
)

func standard03CandidateViewByte(standard03CandidateViewValue uint8) *uint8 {
	return &standard03CandidateViewValue
}

func standard03CandidateViewSelection(standard03CandidateViewOperation uint8, standard03CandidateViewSelector, standard03CandidateViewDirection string) Standard03CandidateViewSelection {
	standard03CandidateViewEncoding := Standard03CandidateViewUnsignedU8
	if standard03CandidateViewOperation == 4 {
		standard03CandidateViewEncoding = Standard03CandidateViewDecimalChunks
	}
	return Standard03CandidateViewSelection{Standard03CandidateViewFamily: "standard-counter", Standard03CandidateViewVariant: standard03CandidateViewFmt.Sprintf("03/%02d", standard03CandidateViewOperation), Standard03CandidateViewFunction: standard03CandidateViewByte(3), Standard03CandidateViewOperation: standard03CandidateViewByte(standard03CandidateViewOperation), Standard03CandidateViewSelector: standard03CandidateViewSelector, Standard03CandidateViewDirection: standard03CandidateViewDirection, Standard03CandidateViewEncoding: standard03CandidateViewEncoding}
}

func standard03CandidateViewNoProjection(standard03CandidateViewT *standard03CandidateViewTesting.T, standard03CandidateViewView Standard03CandidateView) {
	standard03CandidateViewT.Helper()
	if standard03CandidateViewView.Standard03CandidateViewLocalCounter != nil || standard03CandidateViewView.Standard03CandidateViewLocalHours != nil || standard03CandidateViewView.Standard03CandidateViewLocalMinute != nil {
		standard03CandidateViewT.Fatalf("unexpected projection: %+v", standard03CandidateViewView)
	}
}

func TestAb76Fixtures(standard03CandidateViewT *standard03CandidateViewTesting.T) {
	standard03CandidateViewData, standard03CandidateViewErr := standard03CandidateViewOs.ReadFile("testdata/ab76_standard03_candidate_views.json")
	if standard03CandidateViewErr != nil {
		standard03CandidateViewT.Fatal(standard03CandidateViewErr)
	}
	var standard03CandidateViewCases []struct {
		Standard03CandidateViewLabel     string   `json:"Ab76Label"`
		Standard03CandidateViewOperation uint8    `json:"Ab76Operation"`
		Standard03CandidateViewBody      []uint8  `json:"Ab76Body"`
		Standard03CandidateViewNames     []string `json:"Ab76Names"`
		Standard03CandidateViewValues    []uint8  `json:"Ab76Values"`
		Standard03CandidateViewCounter   *int     `json:"Ab76Counter"`
		Standard03CandidateViewHours     *int     `json:"Ab76Hours"`
		Standard03CandidateViewMinute    *int     `json:"Ab76Minute"`
	}
	if standard03CandidateViewErr = standard03CandidateViewJSON.Unmarshal(standard03CandidateViewData, &standard03CandidateViewCases); standard03CandidateViewErr != nil {
		standard03CandidateViewT.Fatal(standard03CandidateViewErr)
	}
	if len(standard03CandidateViewCases) != 3 {
		standard03CandidateViewT.Fatal("ab76 fixture inventory")
	}
	for _, standard03CandidateViewCase := range standard03CandidateViewCases {
		standard03CandidateViewT.Run(standard03CandidateViewCase.Standard03CandidateViewLabel, func(standard03CandidateViewT *standard03CandidateViewTesting.T) {
			for _, standard03CandidateViewMode := range []struct{ Standard03CandidateViewSelector, Standard03CandidateViewDirection string }{{"candidate", "request"}, {"standard", "positive"}} {
				if standard03CandidateViewCase.Standard03CandidateViewOperation == 5 && standard03CandidateViewMode.Standard03CandidateViewSelector == "standard" {
					continue
				}
				standard03CandidateViewView := Standard03CandidateViewInterpret(standard03CandidateViewSelection(standard03CandidateViewCase.Standard03CandidateViewOperation, standard03CandidateViewMode.Standard03CandidateViewSelector, standard03CandidateViewMode.Standard03CandidateViewDirection), standard03CandidateViewCase.Standard03CandidateViewBody)
				if standard03CandidateViewView.Standard03CandidateViewState != Standard03CandidateViewComplete || !standard03CandidateViewReflect.DeepEqual(standard03CandidateViewView.Standard03CandidateViewRaw, standard03CandidateViewCase.Standard03CandidateViewBody) || len(standard03CandidateViewView.Standard03CandidateViewComponents) != len(standard03CandidateViewCase.Standard03CandidateViewNames) || len(standard03CandidateViewView.Standard03CandidateViewRemainder) != 0 {
					standard03CandidateViewT.Fatalf("ab76 shape: %+v", standard03CandidateViewView)
				}
				for standard03CandidateViewIndex, standard03CandidateViewComponent := range standard03CandidateViewView.Standard03CandidateViewComponents {
					if standard03CandidateViewComponent.Standard03CandidateViewName != standard03CandidateViewCase.Standard03CandidateViewNames[standard03CandidateViewIndex] || standard03CandidateViewComponent.Standard03CandidateViewState != Standard03CandidateViewValid || standard03CandidateViewComponent.Standard03CandidateViewValue == nil || *standard03CandidateViewComponent.Standard03CandidateViewValue != standard03CandidateViewCase.Standard03CandidateViewValues[standard03CandidateViewIndex] || standard03CandidateViewComponent.Standard03CandidateViewRaw == nil || *standard03CandidateViewComponent.Standard03CandidateViewRaw != standard03CandidateViewCase.Standard03CandidateViewBody[standard03CandidateViewIndex] || standard03CandidateViewComponent.Standard03CandidateViewWidth != 1 || standard03CandidateViewComponent.Standard03CandidateViewEncoding != standard03CandidateViewView.Standard03CandidateViewSelection.Standard03CandidateViewEncoding {
						standard03CandidateViewT.Fatalf("ab76 component: %+v", standard03CandidateViewComponent)
					}
					if standard03CandidateViewCase.Standard03CandidateViewOperation == 4 {
						if standard03CandidateViewComponent.Standard03CandidateViewDefault == nil || *standard03CandidateViewComponent.Standard03CandidateViewDefault != 0 {
							standard03CandidateViewT.Fatal("ab76 decimal default")
						}
					} else if standard03CandidateViewComponent.Standard03CandidateViewDefault != nil {
						standard03CandidateViewT.Fatal("ab76 invented default")
					}
				}
				if !standard03CandidateViewReflect.DeepEqual(standard03CandidateViewView.Standard03CandidateViewLocalCounter, standard03CandidateViewCase.Standard03CandidateViewCounter) || !standard03CandidateViewReflect.DeepEqual(standard03CandidateViewView.Standard03CandidateViewLocalHours, standard03CandidateViewCase.Standard03CandidateViewHours) || !standard03CandidateViewReflect.DeepEqual(standard03CandidateViewView.Standard03CandidateViewLocalMinute, standard03CandidateViewCase.Standard03CandidateViewMinute) {
					standard03CandidateViewT.Fatalf("ab76 independent projections: %+v", standard03CandidateViewView)
				}
				if standard03CandidateViewView.Standard03CandidateViewImplementationGrade != "Proven" || standard03CandidateViewView.Standard03CandidateViewMeaningGrade != "Hypothesis" || len(standard03CandidateViewView.Standard03CandidateViewNativeProperties) != 10 {
					standard03CandidateViewT.Fatal("ab76 grades")
				}
				for _, standard03CandidateViewGrade := range standard03CandidateViewView.Standard03CandidateViewNativeProperties {
					if standard03CandidateViewGrade != "Unknown" {
						standard03CandidateViewT.Fatal("ab76 native qualification")
					}
				}
			}
		})
	}
}

func TestAb76SelectionAxes(standard03CandidateViewT *standard03CandidateViewTesting.T) {
	for _, standard03CandidateViewOperation := range []uint8{4, 5, 6} {
		for _, standard03CandidateViewMode := range []struct{ Standard03CandidateViewSelector, Standard03CandidateViewDirection string }{{"candidate", "request"}, {"standard", "positive"}} {
			if standard03CandidateViewOperation == 5 && standard03CandidateViewMode.Standard03CandidateViewSelector == "standard" {
				continue
			}
			standard03CandidateViewControl := standard03CandidateViewSelection(standard03CandidateViewOperation, standard03CandidateViewMode.Standard03CandidateViewSelector, standard03CandidateViewMode.Standard03CandidateViewDirection)
			standard03CandidateViewBody := []byte{1, 2, 3}
			if standard03CandidateViewOperation == 6 {
				standard03CandidateViewBody = []byte{59, 1, 2, 3}
			}
			if Standard03CandidateViewInterpret(standard03CandidateViewControl, standard03CandidateViewBody).Standard03CandidateViewState != Standard03CandidateViewComplete {
				standard03CandidateViewT.Fatal("ab76 control")
			}
			for _, standard03CandidateViewCase := range []struct {
				Standard03CandidateViewLabel  string
				Standard03CandidateViewMutate func(*Standard03CandidateViewSelection)
			}{
				{"ab76-family", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewFamily = "ab76-unknown"
				}},
				{"ab76-family-absent", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewFamily = ""
				}},
				{"ab76-variant", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewVariant = "03/07"
				}},
				{"ab76-variant08", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewVariant = "03/08"
				}},
				{"ab76-variant10", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewVariant = "03/10"
				}},
				{"ab76-variant-absent", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewVariant = ""
				}},
				{"ab76-variant05-06", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewVariant = "03/05"
					if standard03CandidateViewOperation == 5 {
						standard03CandidateViewS.Standard03CandidateViewVariant = "03/06"
					}
				}},
				{"ab76-function", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewFunction = standard03CandidateViewByte(9)
				}},
				{"ab76-function-absent", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewFunction = nil
				}},
				{"ab76-function-zero", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewFunction = standard03CandidateViewByte(0)
				}},
				{"ab76-operation", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewOperation = standard03CandidateViewByte(7)
				}},
				{"ab76-operation-absent", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewOperation = nil
				}},
				{"ab76-operation-zero", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewOperation = standard03CandidateViewByte(0)
				}},
				{"ab76-direction", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewDirection = "ab76-unknown"
				}},
				{"ab76-direction-absent", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewDirection = ""
				}},
				{"ab76-direction-opposite", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewDirection = "positive"
					if standard03CandidateViewMode.Standard03CandidateViewDirection == "positive" {
						standard03CandidateViewS.Standard03CandidateViewDirection = "request"
					}
				}},
				{"ab76-selector", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewSelector = "ab76-unknown"
				}},
				{"ab76-selector-absent", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewSelector = ""
				}},
				{"ab76-selector-opposite", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewSelector = "standard"
					if standard03CandidateViewMode.Standard03CandidateViewSelector == "standard" {
						standard03CandidateViewS.Standard03CandidateViewSelector = "candidate"
					}
				}},
				{"ab76-encoding-unspecified", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewEncoding = Standard03CandidateViewEncodingUnspecified
				}},
				{"ab76-encoding-unknown", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewEncoding = "ab76-unknown"
				}},
				{"ab76-encoding-opposite", func(standard03CandidateViewS *Standard03CandidateViewSelection) {
					standard03CandidateViewS.Standard03CandidateViewEncoding = Standard03CandidateViewDecimalChunks
					if standard03CandidateViewOperation == 4 {
						standard03CandidateViewS.Standard03CandidateViewEncoding = Standard03CandidateViewUnsignedU8
					}
				}},
			} {
				standard03CandidateViewT.Run(standard03CandidateViewFmt.Sprintf("ab76-%d-%s-%s", standard03CandidateViewOperation, standard03CandidateViewMode.Standard03CandidateViewSelector, standard03CandidateViewCase.Standard03CandidateViewLabel), func(standard03CandidateViewT *standard03CandidateViewTesting.T) {
					standard03CandidateViewChanged := standard03CandidateViewControl
					standard03CandidateViewCase.Standard03CandidateViewMutate(&standard03CandidateViewChanged)
					standard03CandidateViewView := Standard03CandidateViewInterpret(standard03CandidateViewChanged, standard03CandidateViewBody)
					if standard03CandidateViewView.Standard03CandidateViewState != Standard03CandidateViewRawOnly || len(standard03CandidateViewView.Standard03CandidateViewComponents) != 0 || !standard03CandidateViewReflect.DeepEqual(standard03CandidateViewView.Standard03CandidateViewRaw, standard03CandidateViewBody) || !standard03CandidateViewReflect.DeepEqual(standard03CandidateViewView.Standard03CandidateViewSelection, standard03CandidateViewChanged) {
						standard03CandidateViewT.Fatalf("ab76 axis admitted: %+v", standard03CandidateViewView)
					}
					standard03CandidateViewNoProjection(standard03CandidateViewT, standard03CandidateViewView)
				})
			}
		}
	}
	standard03CandidateViewUnsupported := standard03CandidateViewSelection(5, "standard", "positive")
	if Standard03CandidateViewInterpret(standard03CandidateViewUnsupported, []byte{1, 2, 3}).Standard03CandidateViewState != Standard03CandidateViewRawOnly {
		standard03CandidateViewT.Fatal("ab76 05 standard admitted")
	}
}

func TestAb76LengthsAndFields(standard03CandidateViewT *standard03CandidateViewTesting.T) {
	for _, standard03CandidateViewOperation := range []uint8{4, 5, 6} {
		standard03CandidateViewS := standard03CandidateViewSelection(standard03CandidateViewOperation, "candidate", "request")
		standard03CandidateViewMaxima := []byte{99, 99, 99}
		if standard03CandidateViewOperation == 5 {
			standard03CandidateViewMaxima = []byte{59, 59, 59}
		}
		if standard03CandidateViewOperation == 6 {
			standard03CandidateViewMaxima = []byte{59, 99, 99, 99}
		}
		for standard03CandidateViewLength := 0; standard03CandidateViewLength <= len(standard03CandidateViewMaxima)+1; standard03CandidateViewLength++ {
			standard03CandidateViewT.Run(standard03CandidateViewFmt.Sprintf("ab76-%d-length-%d", standard03CandidateViewOperation, standard03CandidateViewLength), func(standard03CandidateViewT *standard03CandidateViewTesting.T) {
				standard03CandidateViewBody := append(append([]byte(nil), standard03CandidateViewMaxima...), 42)[:standard03CandidateViewLength]
				standard03CandidateViewView := Standard03CandidateViewInterpret(standard03CandidateViewS, standard03CandidateViewBody)
				standard03CandidateViewWantState := Standard03CandidateViewComplete
				if standard03CandidateViewLength < len(standard03CandidateViewMaxima) {
					standard03CandidateViewWantState = Standard03CandidateViewIncomplete
				}
				if standard03CandidateViewLength > len(standard03CandidateViewMaxima) {
					standard03CandidateViewWantState = Standard03CandidateViewExtra
				}
				if standard03CandidateViewView.Standard03CandidateViewState != standard03CandidateViewWantState || !standard03CandidateViewReflect.DeepEqual(standard03CandidateViewView.Standard03CandidateViewRaw, standard03CandidateViewBody) {
					standard03CandidateViewT.Fatalf("ab76 length state: %+v", standard03CandidateViewView)
				}
				if standard03CandidateViewWantState != Standard03CandidateViewComplete {
					standard03CandidateViewNoProjection(standard03CandidateViewT, standard03CandidateViewView)
				}
				if standard03CandidateViewWantState == Standard03CandidateViewExtra && !standard03CandidateViewReflect.DeepEqual(standard03CandidateViewView.Standard03CandidateViewRemainder, []byte{42}) {
					standard03CandidateViewT.Fatal("ab76 remainder")
				}
				for standard03CandidateViewIndex, standard03CandidateViewComponent := range standard03CandidateViewView.Standard03CandidateViewComponents {
					if standard03CandidateViewIndex < standard03CandidateViewLength {
						if standard03CandidateViewComponent.Standard03CandidateViewState != Standard03CandidateViewValid || standard03CandidateViewComponent.Standard03CandidateViewValue == nil || *standard03CandidateViewComponent.Standard03CandidateViewValue != standard03CandidateViewMaxima[standard03CandidateViewIndex] {
							standard03CandidateViewT.Fatal("ab76 dropped partial component")
						}
					} else if standard03CandidateViewComponent.Standard03CandidateViewState != Standard03CandidateViewAbsent || standard03CandidateViewComponent.Standard03CandidateViewRaw != nil || standard03CandidateViewComponent.Standard03CandidateViewValue != nil {
						standard03CandidateViewT.Fatal("ab76 absent became zero")
					}
				}
			})
		}
		for standard03CandidateViewIndex, standard03CandidateViewMaximum := range standard03CandidateViewMaxima {
			for _, standard03CandidateViewValue := range []byte{0, standard03CandidateViewMaximum, standard03CandidateViewMaximum + 1, 255} {
				standard03CandidateViewT.Run(standard03CandidateViewFmt.Sprintf("ab76-%d-field-%d-value-%d", standard03CandidateViewOperation, standard03CandidateViewIndex, standard03CandidateViewValue), func(standard03CandidateViewT *standard03CandidateViewTesting.T) {
					standard03CandidateViewBody := append([]byte(nil), standard03CandidateViewMaxima...)
					standard03CandidateViewBody[standard03CandidateViewIndex] = standard03CandidateViewValue
					standard03CandidateViewView := Standard03CandidateViewInterpret(standard03CandidateViewS, standard03CandidateViewBody)
					standard03CandidateViewValid := standard03CandidateViewValue <= standard03CandidateViewMaximum
					if (standard03CandidateViewView.Standard03CandidateViewState == Standard03CandidateViewComplete) != standard03CandidateViewValid {
						standard03CandidateViewT.Fatal("ab76 boundary")
					}
					if !standard03CandidateViewValid {
						if standard03CandidateViewView.Standard03CandidateViewState != Standard03CandidateViewInvalid {
							standard03CandidateViewT.Fatal("ab76 invalid state")
						}
						standard03CandidateViewNoProjection(standard03CandidateViewT, standard03CandidateViewView)
					}
					for standard03CandidateViewOther, standard03CandidateViewComponent := range standard03CandidateViewView.Standard03CandidateViewComponents {
						if standard03CandidateViewComponent.Standard03CandidateViewMaximum != standard03CandidateViewMaxima[standard03CandidateViewOther] || standard03CandidateViewComponent.Standard03CandidateViewRaw == nil || *standard03CandidateViewComponent.Standard03CandidateViewRaw != standard03CandidateViewBody[standard03CandidateViewOther] {
							standard03CandidateViewT.Fatal("ab76 raw or domain lost")
						}
						if standard03CandidateViewOther == standard03CandidateViewIndex && !standard03CandidateViewValid {
							if standard03CandidateViewComponent.Standard03CandidateViewState != Standard03CandidateViewInvalid || standard03CandidateViewComponent.Standard03CandidateViewValue != nil {
								standard03CandidateViewT.Fatal("ab76 invalid accepted")
							}
						} else if standard03CandidateViewComponent.Standard03CandidateViewState != Standard03CandidateViewValid || standard03CandidateViewComponent.Standard03CandidateViewValue == nil || *standard03CandidateViewComponent.Standard03CandidateViewValue != standard03CandidateViewBody[standard03CandidateViewOther] {
							standard03CandidateViewT.Fatal("ab76 unaffected field lost")
						}
					}
				})
			}
		}
	}
}

func TestAb76OwnedMemory(standard03CandidateViewT *standard03CandidateViewTesting.T) {
	standard03CandidateViewS := standard03CandidateViewSelection(4, "candidate", "request")
	standard03CandidateViewBody := []byte{99, 1, 2, 42}
	standard03CandidateViewView := Standard03CandidateViewInterpret(standard03CandidateViewS, standard03CandidateViewBody)
	standard03CandidateViewBaseline := Standard03CandidateViewInterpret(standard03CandidateViewS, standard03CandidateViewBody)
	*standard03CandidateViewS.Standard03CandidateViewFunction, *standard03CandidateViewS.Standard03CandidateViewOperation = 0, 0
	standard03CandidateViewBody[0], standard03CandidateViewBody[3] = 0, 0
	if !standard03CandidateViewReflect.DeepEqual(standard03CandidateViewView, standard03CandidateViewBaseline) {
		standard03CandidateViewT.Fatal("ab76 input alias")
	}
	standard03CandidateViewView.Standard03CandidateViewRaw[0], standard03CandidateViewView.Standard03CandidateViewRaw[3] = 8, 8
	standard03CandidateViewView.Standard03CandidateViewRemainder[0] = 9
	*standard03CandidateViewView.Standard03CandidateViewComponents[0].Standard03CandidateViewRaw = 7
	*standard03CandidateViewView.Standard03CandidateViewComponents[0].Standard03CandidateViewValue = 6
	*standard03CandidateViewView.Standard03CandidateViewComponents[0].Standard03CandidateViewDefault = 5
	*standard03CandidateViewView.Standard03CandidateViewSelection.Standard03CandidateViewFunction = 0
	*standard03CandidateViewView.Standard03CandidateViewSelection.Standard03CandidateViewOperation = 0
	standard03CandidateViewView.Standard03CandidateViewNativeProperties["physical-support"] = "changed"
	if *standard03CandidateViewBaseline.Standard03CandidateViewComponents[0].Standard03CandidateViewValue != 99 || standard03CandidateViewView.Standard03CandidateViewRaw[0] != 8 || standard03CandidateViewView.Standard03CandidateViewRaw[3] != 8 || *standard03CandidateViewView.Standard03CandidateViewComponents[0].Standard03CandidateViewRaw != 7 || *standard03CandidateViewView.Standard03CandidateViewComponents[1].Standard03CandidateViewDefault != 0 {
		standard03CandidateViewT.Fatal("ab76 return internal alias")
	}
	if !standard03CandidateViewReflect.DeepEqual(standard03CandidateViewBaseline, Standard03CandidateViewInterpret(standard03CandidateViewSelection(4, "candidate", "request"), []byte{99, 1, 2, 42})) {
		standard03CandidateViewT.Fatal("ab76 shared return storage")
	}
	if standard03CandidateViewBody[0] != 0 || standard03CandidateViewBody[3] != 0 {
		standard03CandidateViewT.Fatal("ab76 output aliases input")
	}
	standard03CandidateViewComplete := Standard03CandidateViewInterpret(standard03CandidateViewSelection(6, "candidate", "request"), []byte{59, 1, 2, 3})
	*standard03CandidateViewComplete.Standard03CandidateViewLocalMinute, *standard03CandidateViewComplete.Standard03CandidateViewLocalHours = 0, 0
	standard03CandidateViewFresh := Standard03CandidateViewInterpret(standard03CandidateViewSelection(6, "candidate", "request"), []byte{59, 1, 2, 3})
	if *standard03CandidateViewFresh.Standard03CandidateViewLocalMinute != 59 || *standard03CandidateViewFresh.Standard03CandidateViewLocalHours != 30201 || *standard03CandidateViewComplete.Standard03CandidateViewComponents[0].Standard03CandidateViewValue != 59 {
		standard03CandidateViewT.Fatal("ab76 projection alias")
	}
}

func TestAb76ProjectionAxes(standard03CandidateViewT *standard03CandidateViewTesting.T) {
	for _, standard03CandidateViewCase := range []struct {
		Standard03CandidateViewLabel     string
		Standard03CandidateViewOperation uint8
		Standard03CandidateViewBody      []byte
		Standard03CandidateViewCounter   int
		Standard03CandidateViewHours     int
		Standard03CandidateViewMinute    int
	}{
		{"ab76-c0-zero", 4, []byte{0, 1, 2}, 20100, 0, 0},
		{"ab76-c1-zero", 4, []byte{99, 0, 2}, 20099, 0, 0},
		{"ab76-c2-zero", 4, []byte{99, 1, 0}, 199, 0, 0},
		{"ab76-c-all-max", 4, []byte{99, 99, 99}, 999999, 0, 0},
		{"ab76-minute-zero", 6, []byte{0, 1, 2, 3}, 0, 30201, 0},
		{"ab76-h0-zero", 6, []byte{59, 0, 2, 3}, 0, 30200, 59},
		{"ab76-h1-zero", 6, []byte{59, 1, 0, 3}, 0, 30001, 59},
		{"ab76-h2-zero", 6, []byte{59, 1, 2, 0}, 0, 201, 59},
		{"ab76-h-all-max", 6, []byte{59, 99, 99, 99}, 0, 999999, 59},
	} {
		standard03CandidateViewT.Run(standard03CandidateViewCase.Standard03CandidateViewLabel, func(standard03CandidateViewT *standard03CandidateViewTesting.T) {
			standard03CandidateViewView := Standard03CandidateViewInterpret(standard03CandidateViewSelection(standard03CandidateViewCase.Standard03CandidateViewOperation, "candidate", "request"), standard03CandidateViewCase.Standard03CandidateViewBody)
			if standard03CandidateViewView.Standard03CandidateViewState != Standard03CandidateViewComplete {
				standard03CandidateViewT.Fatal("ab76 zero component is present")
			}
			if standard03CandidateViewCase.Standard03CandidateViewOperation == 4 {
				if standard03CandidateViewView.Standard03CandidateViewLocalCounter == nil || *standard03CandidateViewView.Standard03CandidateViewLocalCounter != standard03CandidateViewCase.Standard03CandidateViewCounter || standard03CandidateViewView.Standard03CandidateViewLocalHours != nil || standard03CandidateViewView.Standard03CandidateViewLocalMinute != nil {
					standard03CandidateViewT.Fatal("ab76 counter weight")
				}
			} else if standard03CandidateViewView.Standard03CandidateViewLocalHours == nil || *standard03CandidateViewView.Standard03CandidateViewLocalHours != standard03CandidateViewCase.Standard03CandidateViewHours || standard03CandidateViewView.Standard03CandidateViewLocalMinute == nil || *standard03CandidateViewView.Standard03CandidateViewLocalMinute != standard03CandidateViewCase.Standard03CandidateViewMinute || standard03CandidateViewView.Standard03CandidateViewLocalCounter != nil {
				standard03CandidateViewT.Fatal("ab76 hour or minute weight")
			}
		})
	}
}

func TestAb76RawOnlyAndUnselectedFunctions(standard03CandidateViewT *standard03CandidateViewTesting.T) {
	for _, standard03CandidateViewOperation := range []uint8{0, 7, 8, 10, 255} {
		standard03CandidateViewS := standard03CandidateViewSelection(standard03CandidateViewOperation, "candidate", "request")
		for standard03CandidateViewLength := 0; standard03CandidateViewLength <= 5; standard03CandidateViewLength++ {
			standard03CandidateViewT.Run(standard03CandidateViewFmt.Sprintf("ab76-raw-%d-length-%d", standard03CandidateViewOperation, standard03CandidateViewLength), func(standard03CandidateViewT *standard03CandidateViewTesting.T) {
				standard03CandidateViewBody := make([]byte, standard03CandidateViewLength)
				standard03CandidateViewView := Standard03CandidateViewInterpret(standard03CandidateViewS, standard03CandidateViewBody)
				if standard03CandidateViewView.Standard03CandidateViewState != Standard03CandidateViewRawOnly || len(standard03CandidateViewView.Standard03CandidateViewRaw) != standard03CandidateViewLength || len(standard03CandidateViewView.Standard03CandidateViewComponents) != 0 {
					standard03CandidateViewT.Fatal("ab76 unselected body interpreted")
				}
				standard03CandidateViewNoProjection(standard03CandidateViewT, standard03CandidateViewView)
				if standard03CandidateViewLength > 0 {
					standard03CandidateViewBody[0] = 42
					if standard03CandidateViewView.Standard03CandidateViewRaw[0] != 0 {
						standard03CandidateViewT.Fatal("ab76 raw-only input alias")
					}
					standard03CandidateViewView.Standard03CandidateViewRaw[0] = 99
					if standard03CandidateViewBody[0] != 42 {
						standard03CandidateViewT.Fatal("ab76 raw-only output alias")
					}
				}
			})
		}
	}
}
