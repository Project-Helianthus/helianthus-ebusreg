package vaillant

import (
	b514MenuViewJSON "encoding/json"
	b514MenuViewFmt "fmt"
	b514MenuViewOs "os"
	b514MenuViewReflect "reflect"
	b514MenuViewStrings "strings"
	b514MenuViewTesting "testing"
)

func b514MenuViewInput(b514MenuViewOperation int, b514MenuViewView string, b514MenuViewFields map[string]any) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 245, "secondary": 20, "operation": b514MenuViewOperation, "direction": "positive", "view": b514MenuViewView, "family": map[string]any{"opaque": []any{"synthetic-family", 0}}, "variant": []any{"synthetic-variant", false}}, "fields": b514MenuViewFields}
}
func b514MenuViewCall(b514MenuViewT *b514MenuViewTesting.T, b514MenuViewInputObject map[string]any) map[string]any {
	b514MenuViewT.Helper()
	b514MenuViewInputJSON, b514MenuViewErr := b514MenuViewJSON.Marshal(b514MenuViewInputObject)
	if b514MenuViewErr != nil {
		b514MenuViewT.Fatal(b514MenuViewErr)
	}
	b514MenuViewOutputJSON, b514MenuViewErr := B514MenuViewInterpret(b514MenuViewInputJSON)
	if b514MenuViewErr != nil {
		b514MenuViewT.Fatal(b514MenuViewErr)
	}
	b514MenuViewOutput := map[string]any{}
	if b514MenuViewErr := b514MenuViewJSON.Unmarshal(b514MenuViewOutputJSON, &b514MenuViewOutput); b514MenuViewErr != nil {
		b514MenuViewT.Fatal(b514MenuViewErr)
	}
	if b514MenuViewOutput["physical"] != "Unknown" || b514MenuViewOutput["qualified"] != false {
		b514MenuViewT.Fatal("qualification promoted", b514MenuViewOutput)
	}
	b514MenuViewExpected := map[string]any{}
	if b514MenuViewErr := b514MenuViewJSON.Unmarshal(b514MenuViewInputJSON, &b514MenuViewExpected); b514MenuViewErr != nil {
		b514MenuViewT.Fatal(b514MenuViewErr)
	}
	for b514MenuViewKey, b514MenuViewValue := range map[string]any{"context": b514MenuViewExpected["context"], "raw": b514MenuViewExpected["fields"]} {
		if !b514MenuViewReflect.DeepEqual(b514MenuViewOutput[b514MenuViewKey], b514MenuViewValue) {
			b514MenuViewT.Fatalf("lost %s: %#v", b514MenuViewKey, b514MenuViewOutput)
		}
	}
	b514MenuViewActualArguments, b514MenuViewActualPresent := b514MenuViewOutput["arguments"]
	b514MenuViewExpectedArguments, b514MenuViewExpectedPresent := b514MenuViewExpected["arguments"]
	if b514MenuViewActualPresent != b514MenuViewExpectedPresent || !b514MenuViewReflect.DeepEqual(b514MenuViewActualArguments, b514MenuViewExpectedArguments) {
		b514MenuViewT.Fatal("argument presence changed")
	}
	for _, b514MenuViewField := range b514MenuViewOutput["fields"].(map[string]any) {
		if b514MenuViewField.(map[string]any)["grade"] != "Hypothesis" {
			b514MenuViewT.Fatal("family grade changed")
		}
	}
	return b514MenuViewOutput
}
func b514MenuViewField(b514MenuViewOutput map[string]any, b514MenuViewName string) map[string]any {
	return b514MenuViewOutput["fields"].(map[string]any)[b514MenuViewName].(map[string]any)
}
func b514MenuViewNoNumber(b514MenuViewT *b514MenuViewTesting.T, b514MenuViewField map[string]any) {
	b514MenuViewT.Helper()
	if _, b514MenuViewPresent := b514MenuViewField["number"]; b514MenuViewPresent {
		b514MenuViewT.Fatal("unexpected number", b514MenuViewField)
	}
}
func TestAb19Ao19Feedback(b514MenuViewT *b514MenuViewTesting.T) {
	for _, b514MenuViewView := range []string{"generic-result", "specific-a-result", "specific-b-result"} {
		for b514MenuViewByte := 0; b514MenuViewByte <= 255; b514MenuViewByte++ {
			b514MenuViewT.Run(b514MenuViewFmt.Sprintf("%s/%d", b514MenuViewView, b514MenuViewByte), func(b514MenuViewT *b514MenuViewTesting.T) {
				b514MenuViewOutput := b514MenuViewCall(b514MenuViewT, b514MenuViewInput(6, b514MenuViewView, map[string]any{"feedback": []int{b514MenuViewByte, 99}, "unmodelled": []int{4, 5}}))
				b514MenuViewExpected := ""
				if b514MenuViewByte == 1 {
					b514MenuViewExpected = "not-successful"
					if b514MenuViewView == "generic-result" {
						b514MenuViewExpected = "in-progress"
					}
				}
				if b514MenuViewByte == 255 && b514MenuViewView == "generic-result" {
					b514MenuViewExpected = "unsupported"
				}
				if b514MenuViewByte >= 3 && b514MenuViewView != "generic-result" {
					b514MenuViewExpected = "invalid"
				}
				b514MenuViewFieldValue := b514MenuViewField(b514MenuViewOutput, "feedback")
				b514MenuViewLabel, _ := b514MenuViewFieldValue["label"].(string)
				if b514MenuViewLabel != b514MenuViewExpected || b514MenuViewOutput["state"] != "Hypothesis" {
					b514MenuViewT.Fatal(b514MenuViewOutput)
				}
				b514MenuViewState := "Hypothesis"
				if b514MenuViewExpected == "" {
					b514MenuViewState = "unmapped"
				}
				if b514MenuViewFieldValue["state"] != b514MenuViewState || !b514MenuViewReflect.DeepEqual(b514MenuViewFieldValue["raw"], []any{float64(b514MenuViewByte)}) {
					b514MenuViewT.Fatal(b514MenuViewFieldValue)
				}
				b514MenuViewNoNumber(b514MenuViewT, b514MenuViewFieldValue)
			})
		}
	}
	for b514MenuViewByte := 0; b514MenuViewByte <= 255; b514MenuViewByte++ {
		b514MenuViewT.Run(b514MenuViewFmt.Sprintf("progress/%d", b514MenuViewByte), func(b514MenuViewT *b514MenuViewTesting.T) {
			b514MenuViewOutput := b514MenuViewCall(b514MenuViewT, b514MenuViewInput(7, "progress", map[string]any{"feedback": []int{b514MenuViewByte}}))
			b514MenuViewExpected := ""
			if b514MenuViewByte < 7 {
				b514MenuViewExpected = []string{"idle", "progress", "passed", "warning", "combustion-error", "gas-type-error", "aborted"}[b514MenuViewByte]
			}
			if b514MenuViewByte == 255 {
				b514MenuViewExpected = "unsupported"
			}
			b514MenuViewFieldValue := b514MenuViewField(b514MenuViewOutput, "feedback")
			b514MenuViewLabel, _ := b514MenuViewFieldValue["label"].(string)
			b514MenuViewState := "Hypothesis"
			if b514MenuViewExpected == "" {
				b514MenuViewState = "unmapped"
			}
			if b514MenuViewLabel != b514MenuViewExpected || b514MenuViewFieldValue["state"] != b514MenuViewState {
				b514MenuViewT.Fatal(b514MenuViewOutput)
			}
			b514MenuViewNoNumber(b514MenuViewT, b514MenuViewFieldValue)
		})
	}
}
func b514MenuViewPressure() []int { return []int{19, 23, 29, 1, 1, 0, 0, 1, 52, 18, 0, 100, 0} }
func b514MenuViewPressureCheck(b514MenuViewT *b514MenuViewTesting.T, b514MenuViewData []int, b514MenuViewExpected map[string]float64) map[string]any {
	b514MenuViewT.Helper()
	b514MenuViewOutput := b514MenuViewCall(b514MenuViewT, b514MenuViewInput(8, "pressure", map[string]any{"logical_object": b514MenuViewData, "other": []int{7}}))
	for b514MenuViewName, b514MenuViewNumber := range b514MenuViewExpected {
		b514MenuViewValue := b514MenuViewField(b514MenuViewOutput, b514MenuViewName)
		if b514MenuViewValue["number"] != b514MenuViewNumber || b514MenuViewValue["calculation"] != "Proven" || b514MenuViewValue["state"] != "Hypothesis" {
			b514MenuViewT.Fatalf("%s: %#v", b514MenuViewName, b514MenuViewValue)
		}
	}
	b514MenuViewStep := b514MenuViewField(b514MenuViewOutput, "step")
	if len(b514MenuViewData) > 10 {
		if b514MenuViewStep["state"] != "unit-conflict" || !b514MenuViewReflect.DeepEqual(b514MenuViewStep["flags"], []any{"constant-1-mbar-versus-default-0"}) {
			b514MenuViewT.Fatal(b514MenuViewStep)
		}
	}
	b514MenuViewNoNumber(b514MenuViewT, b514MenuViewStep)
	b514MenuViewNoNumber(b514MenuViewT, b514MenuViewField(b514MenuViewOutput, "feedback"))
	return b514MenuViewOutput
}
func TestAb19Ao19Pressure(b514MenuViewT *b514MenuViewTesting.T) {
	for _, b514MenuViewCase := range []struct {
		b514MenuViewIndex, b514MenuViewByte int
		b514MenuViewName                    string
		b514MenuViewNumber                  float64
	}{
		{4, 2, "pressure", 2}, {5, 1, "pressure", 257}, {6, 3, "minimum", 259}, {7, 2, "minimum", 512}, {8, 53, "maximum", 4661}, {9, 19, "maximum", 4916}, {11, 200, "flow", 2}, {12, 1, "flow", 3.56},
	} {
		b514MenuViewT.Run(b514MenuViewFmt.Sprintf("word/%d", b514MenuViewCase.b514MenuViewIndex), func(b514MenuViewT *b514MenuViewTesting.T) {
			b514MenuViewData := b514MenuViewPressure()
			b514MenuViewData[b514MenuViewCase.b514MenuViewIndex] = b514MenuViewCase.b514MenuViewByte
			b514MenuViewExpected := map[string]float64{"pressure": 1, "minimum": 256, "maximum": 4660, "flow": 1}
			b514MenuViewExpected[b514MenuViewCase.b514MenuViewName] = b514MenuViewCase.b514MenuViewNumber
			b514MenuViewPressureCheck(b514MenuViewT, b514MenuViewData, b514MenuViewExpected)
		})
	}
	for _, b514MenuViewIndex := range []int{0, 1, 2, 3, 10, 13} {
		b514MenuViewT.Run(b514MenuViewFmt.Sprintf("opaque-feedback-step-tail/%d", b514MenuViewIndex), func(b514MenuViewT *b514MenuViewTesting.T) {
			b514MenuViewData := append(b514MenuViewPressure(), 81)
			b514MenuViewData[b514MenuViewIndex] = 255
			b514MenuViewOutput := b514MenuViewPressureCheck(b514MenuViewT, b514MenuViewData, map[string]float64{"pressure": 1, "minimum": 256, "maximum": 4660, "flow": 1})
			if b514MenuViewOutput["state"] != "Hypothesis" || b514MenuViewField(b514MenuViewOutput, "feedback")["state"] != "unmapped" {
				b514MenuViewT.Fatal(b514MenuViewOutput)
			}
		})
	}
	for _, b514MenuViewStep := range []int{0, 1, 2, 255} {
		b514MenuViewData := b514MenuViewPressure()
		b514MenuViewData[10] = b514MenuViewStep
		b514MenuViewPressureCheck(b514MenuViewT, b514MenuViewData, map[string]float64{"pressure": 1, "minimum": 256, "maximum": 4660, "flow": 1})
	}
	for b514MenuViewLength := 0; b514MenuViewLength <= 13; b514MenuViewLength++ {
		b514MenuViewT.Run(b514MenuViewFmt.Sprintf("length/%d", b514MenuViewLength), func(b514MenuViewT *b514MenuViewTesting.T) {
			b514MenuViewData := b514MenuViewPressure()[:b514MenuViewLength]
			b514MenuViewExpected := map[string]float64{}
			for b514MenuViewName, b514MenuViewEnd := range map[string]int{"pressure": 6, "minimum": 8, "maximum": 10, "flow": 13} {
				if b514MenuViewLength >= b514MenuViewEnd {
					b514MenuViewExpected[b514MenuViewName] = map[string]float64{"pressure": 1, "minimum": 256, "maximum": 4660, "flow": 1}[b514MenuViewName]
				}
			}
			b514MenuViewOutput := b514MenuViewPressureCheck(b514MenuViewT, b514MenuViewData, b514MenuViewExpected)
			b514MenuViewState := "partial"
			if b514MenuViewLength == 13 {
				b514MenuViewState = "Hypothesis"
			}
			if b514MenuViewOutput["state"] != b514MenuViewState {
				b514MenuViewT.Fatal(b514MenuViewOutput)
			}
			for b514MenuViewName, b514MenuViewEnd := range map[string]int{"feedback": 4, "pressure": 6, "minimum": 8, "maximum": 10, "step": 11, "flow": 13} {
				if b514MenuViewLength < b514MenuViewEnd {
					b514MenuViewValue := b514MenuViewField(b514MenuViewOutput, b514MenuViewName)
					if b514MenuViewValue["state"] != "incomplete" {
						b514MenuViewT.Fatal(b514MenuViewValue)
					}
					b514MenuViewNoNumber(b514MenuViewT, b514MenuViewValue)
				}
			}
		})
	}
	b514MenuViewData := make([]int, 13)
	b514MenuViewPressureCheck(b514MenuViewT, b514MenuViewData, map[string]float64{"pressure": 0, "minimum": 0, "maximum": 0, "flow": 0})
}
func TestAb19Ao19Selection(b514MenuViewT *b514MenuViewTesting.T) {
	for _, b514MenuViewSelection := range []struct {
		b514MenuViewOperation             int
		b514MenuViewView, b514MenuViewRaw string
	}{{6, "generic-result", "feedback"}, {6, "specific-a-result", "feedback"}, {6, "specific-b-result", "feedback"}, {7, "progress", "feedback"}, {8, "pressure", "logical_object"}} {
		for _, b514MenuViewKey := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
			for _, b514MenuViewMissing := range []bool{false, true} {
				b514MenuViewT.Run(b514MenuViewFmt.Sprintf("%d/%s/%s/%t", b514MenuViewSelection.b514MenuViewOperation, b514MenuViewSelection.b514MenuViewView, b514MenuViewKey, b514MenuViewMissing), func(b514MenuViewT *b514MenuViewTesting.T) {
					b514MenuViewData := []int{1}
					if b514MenuViewSelection.b514MenuViewOperation == 8 {
						b514MenuViewData = b514MenuViewPressure()
					}
					b514MenuViewInputObject := b514MenuViewInput(b514MenuViewSelection.b514MenuViewOperation, b514MenuViewSelection.b514MenuViewView, map[string]any{b514MenuViewSelection.b514MenuViewRaw: b514MenuViewData})
					b514MenuViewControl := b514MenuViewCall(b514MenuViewT, b514MenuViewInputObject)
					b514MenuViewContext := b514MenuViewInputObject["context"].(map[string]any)
					if b514MenuViewMissing {
						delete(b514MenuViewContext, b514MenuViewKey)
					} else {
						b514MenuViewContext[b514MenuViewKey] = map[string]any{"primary": 182, "secondary": 21, "operation": 2, "direction": "request", "view": "other", "family": map[string]any{"opaque": []any{0, true}}, "variant": []any{17, "other"}}[b514MenuViewKey]
					}
					b514MenuViewOutput := b514MenuViewCall(b514MenuViewT, b514MenuViewInputObject)
					if b514MenuViewKey == "family" || b514MenuViewKey == "variant" {
						if !b514MenuViewReflect.DeepEqual(b514MenuViewOutput["fields"], b514MenuViewControl["fields"]) {
							b514MenuViewT.Fatal("opaque label selected view")
						}
					} else if b514MenuViewOutput["state"] != "Unknown" || len(b514MenuViewOutput["fields"].(map[string]any)) != 0 {
						b514MenuViewT.Fatal(b514MenuViewOutput)
					}
				})
			}
		}
	}
	for _, b514MenuViewOperation := range []int{0, 2, 5, 9} {
		b514MenuViewOutput := b514MenuViewCall(b514MenuViewT, b514MenuViewInput(b514MenuViewOperation, "pressure", map[string]any{"logical_object": b514MenuViewPressure()}))
		if b514MenuViewOutput["state"] != "Unknown" {
			b514MenuViewT.Fatal(b514MenuViewOutput)
		}
	}
}
func TestAb19Ao19Presence(b514MenuViewT *b514MenuViewTesting.T) {
	for _, b514MenuViewOperation := range []int{6, 7, 8} {
		for _, b514MenuViewMode := range []string{"missing", "empty", "zero"} {
			b514MenuViewT.Run(b514MenuViewFmt.Sprintf("%d/%s", b514MenuViewOperation, b514MenuViewMode), func(b514MenuViewT *b514MenuViewTesting.T) {
				b514MenuViewView := map[int]string{6: "generic-result", 7: "progress", 8: "pressure"}[b514MenuViewOperation]
				b514MenuViewKey := "feedback"
				if b514MenuViewOperation == 8 {
					b514MenuViewKey = "logical_object"
				}
				b514MenuViewRaw := map[string]any{"other": []int{0}}
				if b514MenuViewMode == "empty" {
					b514MenuViewRaw[b514MenuViewKey] = []int{}
				}
				if b514MenuViewMode == "zero" {
					b514MenuViewRaw[b514MenuViewKey] = []int{0}
				}
				b514MenuViewInputObject := b514MenuViewInput(b514MenuViewOperation, b514MenuViewView, b514MenuViewRaw)
				b514MenuViewOutput := b514MenuViewCall(b514MenuViewT, b514MenuViewInputObject)
				b514MenuViewState := "partial"
				if b514MenuViewMode == "zero" && b514MenuViewOperation != 8 {
					b514MenuViewState = "Hypothesis"
				}
				if b514MenuViewOutput["state"] != b514MenuViewState {
					b514MenuViewT.Fatal(b514MenuViewOutput)
				}
				b514MenuViewInputObject["arguments"] = map[string]any{}
				b514MenuViewCall(b514MenuViewT, b514MenuViewInputObject)
			})
		}
	}
}
func TestAb19Ao19Reject(b514MenuViewT *b514MenuViewTesting.T) {
	for b514MenuViewIndex, b514MenuViewInputText := range []string{
		``, `{`, `null`, `[]`, `{} {}`, `{"context":{},"fields":{},"fields":{}}`,
		`{"context":{"family":{"x":0,"x":1}},"fields":{}}`, `{"context":{"variant":[{"x":0,"x":1}]},"fields":{}}`,
		`{"context":{},"fields":{"feedback":[],"feedback":[0]}}`, `{"context":{},"fields":{},"arguments":{"x":1,"x":2}}`,
		`{"context":null,"fields":{}}`, `{"context":{},"fields":null}`, `{"context":[],"fields":{}}`, `{"context":{},"fields":[]}`, `{"context":{},"fields":{},"extra":0}`,
		`{"context":{},"fields":{},"arguments":null}`, `{"context":{},"fields":{},"arguments":[]}`, `{"context":{},"fields":{},"arguments":{"flag":true}}`,
	} {
		b514MenuViewT.Run(b514MenuViewFmt.Sprintf("structure/%d", b514MenuViewIndex), func(b514MenuViewT *b514MenuViewTesting.T) {
			b514MenuViewOutput, b514MenuViewErr := B514MenuViewInterpret([]byte(b514MenuViewInputText))
			if b514MenuViewErr == nil || b514MenuViewOutput != nil {
				b514MenuViewT.Fatal("accepted invalid structure", b514MenuViewInputText)
			}
		})
	}
	for _, b514MenuViewKey := range []string{"feedback", "logical_object", "unmodelled"} {
		for _, b514MenuViewBad := range []string{`null`, `0`, `"AA=="`, `{}`, `[null]`, `[true]`, `["1"]`, `[-1]`, `[256]`, `[1.5]`, `[1.0]`, `[1e0]`, `[999999999999999999999999999999]`, `[[0]]`, `[{}]`} {
			b514MenuViewInputText := `{"context":{},"fields":{"` + b514MenuViewKey + `":` + b514MenuViewBad + `}}`
			b514MenuViewOutput, b514MenuViewErr := B514MenuViewInterpret([]byte(b514MenuViewInputText))
			if b514MenuViewErr == nil || b514MenuViewOutput != nil {
				b514MenuViewT.Fatal("accepted invalid bytes", b514MenuViewInputText)
			}
		}
	}
	for _, b514MenuViewKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, b514MenuViewBad := range []any{nil, true, []int{}, map[string]any{}} {
			b514MenuViewInputObject := b514MenuViewInput(8, "pressure", map[string]any{})
			b514MenuViewInputObject["context"].(map[string]any)[b514MenuViewKey] = b514MenuViewBad
			b514MenuViewInputJSON, _ := b514MenuViewJSON.Marshal(b514MenuViewInputObject)
			b514MenuViewOutput, b514MenuViewErr := B514MenuViewInterpret(b514MenuViewInputJSON)
			if b514MenuViewErr == nil || b514MenuViewOutput != nil {
				b514MenuViewT.Fatal("accepted invalid context", string(b514MenuViewInputJSON))
			}
		}
	}
	for _, b514MenuViewInputText := range []string{`{"context":{"primary":"245"},"fields":{}}`, `{"context":{"secondary":20.5},"fields":{}}`, `{"context":{"operation":6.0},"fields":{}}`, `{"context":{"direction":0},"fields":{}}`, `{"context":{"view":1},"fields":{}}`} {
		b514MenuViewOutput, b514MenuViewErr := B514MenuViewInterpret([]byte(b514MenuViewInputText))
		if b514MenuViewErr == nil || b514MenuViewOutput != nil {
			b514MenuViewT.Fatal("accepted wrong kind")
		}
	}
}
func TestAb19Ao19Isolation(b514MenuViewT *b514MenuViewTesting.T) {
	b514MenuViewInputJSON, _ := b514MenuViewJSON.Marshal(b514MenuViewInput(8, "pressure", map[string]any{"logical_object": b514MenuViewPressure()}))
	b514MenuViewOriginal := append([]byte(nil), b514MenuViewInputJSON...)
	b514MenuViewFirst, b514MenuViewErr := B514MenuViewInterpret(b514MenuViewInputJSON)
	if b514MenuViewErr != nil {
		b514MenuViewT.Fatal(b514MenuViewErr)
	}
	b514MenuViewSnapshot := append([]byte(nil), b514MenuViewFirst...)
	for b514MenuViewIndex := range b514MenuViewInputJSON {
		b514MenuViewInputJSON[b514MenuViewIndex] = ' '
	}
	if string(b514MenuViewFirst) != string(b514MenuViewSnapshot) {
		b514MenuViewT.Fatal("input aliases output")
	}
	b514MenuViewSecond, b514MenuViewErr := B514MenuViewInterpret(b514MenuViewOriginal)
	if b514MenuViewErr != nil {
		b514MenuViewT.Fatal(b514MenuViewErr)
	}
	b514MenuViewDecoded := map[string]any{}
	if b514MenuViewErr := b514MenuViewJSON.Unmarshal(b514MenuViewFirst, &b514MenuViewDecoded); b514MenuViewErr != nil {
		b514MenuViewT.Fatal(b514MenuViewErr)
	}
	b514MenuViewDecoded["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "changed"
	b514MenuViewDecoded["raw"].(map[string]any)["logical_object"].([]any)[4] = float64(200)
	if b514MenuViewField(b514MenuViewDecoded, "pressure")["raw"].([]any)[0] != float64(1) {
		b514MenuViewT.Fatal("raw evidence aliases semantic raw")
	}
	b514MenuViewField(b514MenuViewDecoded, "pressure")["raw"].([]any)[0] = float64(201)
	if b514MenuViewDecoded["raw"].(map[string]any)["logical_object"].([]any)[4] != float64(200) {
		b514MenuViewT.Fatal("semantic raw aliases raw evidence")
	}
	if b514MenuViewField(b514MenuViewDecoded, "pressure")["number"] != float64(1) {
		b514MenuViewT.Fatal("raw mutation changed calculation")
	}
	for b514MenuViewIndex := range b514MenuViewFirst {
		b514MenuViewFirst[b514MenuViewIndex] = ' '
	}
	b514MenuViewThird, b514MenuViewErr := B514MenuViewInterpret(b514MenuViewOriginal)
	if b514MenuViewErr != nil || string(b514MenuViewSecond) != string(b514MenuViewSnapshot) || string(b514MenuViewThird) != string(b514MenuViewSnapshot) {
		b514MenuViewT.Fatal("output aliases another call")
	}
}
func TestAb19Ao19Fixture(b514MenuViewT *b514MenuViewTesting.T) {
	b514MenuViewData, b514MenuViewErr := b514MenuViewOs.ReadFile("testdata/ab19ao19_b514_menu_views.json")
	if b514MenuViewErr != nil {
		b514MenuViewT.Fatal(b514MenuViewErr)
	}
	b514MenuViewCases := []struct {
		B514MenuViewID      string             `json:"id"`
		B514MenuViewInput   map[string]any     `json:"input"`
		B514MenuViewState   string             `json:"state"`
		B514MenuViewNumbers map[string]float64 `json:"numbers"`
	}{}
	if b514MenuViewErr := b514MenuViewJSON.Unmarshal(b514MenuViewData, &b514MenuViewCases); b514MenuViewErr != nil {
		b514MenuViewT.Fatal(b514MenuViewErr)
	}
	if len(b514MenuViewCases) == 0 {
		b514MenuViewT.Fatal("empty fixture")
	}
	for _, b514MenuViewCase := range b514MenuViewCases {
		if !b514MenuViewStrings.HasPrefix(b514MenuViewCase.B514MenuViewID, "ab19ao19") {
			b514MenuViewT.Fatal("fixture ID")
		}
		b514MenuViewOutput := b514MenuViewCall(b514MenuViewT, b514MenuViewCase.B514MenuViewInput)
		if b514MenuViewOutput["state"] != b514MenuViewCase.B514MenuViewState {
			b514MenuViewT.Fatal(b514MenuViewOutput)
		}
		for b514MenuViewName, b514MenuViewNumber := range b514MenuViewCase.B514MenuViewNumbers {
			if b514MenuViewField(b514MenuViewOutput, b514MenuViewName)["number"] != b514MenuViewNumber {
				b514MenuViewT.Fatal(b514MenuViewOutput)
			}
		}
	}
}

func TestAb19Ao19UnknownSiblings(b514MenuViewT *b514MenuViewTesting.T) {
	for _, b514MenuViewName := range []string{"other", "feedback", "logical_object"} {
		for _, b514MenuViewData := range [][]int{{}, {0}, {255, 1, 2, 3}} {
			b514MenuViewRaw := map[string]any{"logical_object": b514MenuViewPressure(), "feedback": []int{1}, "other": []int{0}}
			b514MenuViewRaw[b514MenuViewName] = b514MenuViewData
			b514MenuViewOutput := b514MenuViewCall(b514MenuViewT, b514MenuViewInput(6, "generic-result", b514MenuViewRaw))
			if b514MenuViewName != "feedback" && b514MenuViewField(b514MenuViewOutput, "feedback")["label"] != "in-progress" {
				b514MenuViewT.Fatal("unknown sibling erased feedback")
			}
			b514MenuViewOutput = b514MenuViewCall(b514MenuViewT, b514MenuViewInput(8, "pressure", b514MenuViewRaw))
			if b514MenuViewName != "logical_object" && b514MenuViewField(b514MenuViewOutput, "pressure")["number"] != float64(1) {
				b514MenuViewT.Fatal("unknown sibling erased pressure")
			}
		}
	}
	b514MenuViewBad := append([]byte(`{"context":{"opaque":"`), 255)
	b514MenuViewBad = append(b514MenuViewBad, []byte(`"},"fields":{}}`)...)
	b514MenuViewOutput, b514MenuViewErr := B514MenuViewInterpret(b514MenuViewBad)
	if b514MenuViewErr == nil || b514MenuViewOutput != nil {
		b514MenuViewT.Fatal("accepted malformed UTF-8")
	}
}
