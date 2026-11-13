package vaillant

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func b513InformationViewInput(b513InformationViewBytes []int) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 19, "operation": 7, "direction": "positive", "view": "information", "family": map[string]any{"opaque": []any{"synthetic-family", 0}}, "variant": "synthetic-variant"}, "fields": map[string]any{"logical_object": b513InformationViewBytes, "unmodelled": []int{0, 255}}}
}
func b513InformationViewRun(b513InformationViewT *testing.T, b513InformationViewInputValue any) map[string]any {
	b513InformationViewT.Helper()
	b513InformationViewEncoded, b513InformationViewError := json.Marshal(b513InformationViewInputValue)
	if b513InformationViewError != nil {
		b513InformationViewT.Fatal(b513InformationViewError)
	}
	b513InformationViewResult, b513InformationViewError := B513InformationViewInterpret(b513InformationViewEncoded)
	if b513InformationViewError != nil {
		b513InformationViewT.Fatal(b513InformationViewError)
	}
	b513InformationViewOutput := map[string]any{}
	if b513InformationViewError = json.Unmarshal(b513InformationViewResult, &b513InformationViewOutput); b513InformationViewError != nil {
		b513InformationViewT.Fatal(b513InformationViewError)
	}
	if b513InformationViewOutput["physical"] != "Unknown" || b513InformationViewOutput["qualified"] != false {
		b513InformationViewT.Fatal("qualification inferred", b513InformationViewOutput)
	}
	b513InformationViewNormalized := map[string]any{}
	_ = json.Unmarshal(b513InformationViewEncoded, &b513InformationViewNormalized)
	if !reflect.DeepEqual(b513InformationViewOutput["context"], b513InformationViewNormalized["context"]) || !reflect.DeepEqual(b513InformationViewOutput["raw"], b513InformationViewNormalized["fields"]) {
		b513InformationViewT.Fatal("context or raw changed", b513InformationViewOutput)
	}
	if b513InformationViewArguments, b513InformationViewPresent := b513InformationViewNormalized["arguments"]; b513InformationViewPresent {
		if !reflect.DeepEqual(b513InformationViewArguments, b513InformationViewOutput["arguments"]) {
			b513InformationViewT.Fatal("arguments changed")
		}
	} else if _, b513InformationViewPresent := b513InformationViewOutput["arguments"]; b513InformationViewPresent {
		b513InformationViewT.Fatal("arguments invented")
	}
	return b513InformationViewOutput
}
func b513InformationViewField(b513InformationViewOutput map[string]any, b513InformationViewName string) map[string]any {
	return b513InformationViewOutput["fields"].(map[string]any)[b513InformationViewName].(map[string]any)
}
func TestAb17Ao17Flags(b513InformationViewT *testing.T) {
	b513InformationViewNames := []string{"capabilities", "fuel", "storage", "features"}
	b513InformationViewLabels := [][]string{{"floor-standing", "condensing", "fan-assisted", "open-vented", "high-modulation-range", "water-heater", "step-modulation", "ECS"}, {"gas", "oil", "ABC"}, {"layer-tank", "warm-start", "external-tank", "combination"}, {"solar-compact", "basic-solar", "AAPS", "advanced-power", "partial-load", "parameter-primary"}}
	for b513InformationViewIndex, b513InformationViewName := range b513InformationViewNames {
		for b513InformationViewBit := 0; b513InformationViewBit < 8; b513InformationViewBit++ {
			b513InformationViewT.Run(fmt.Sprintf("%s_bit%d", b513InformationViewName, b513InformationViewBit), func(b513InformationViewT *testing.T) {
				b513InformationViewBytes := []int{91, 82, 73, 0, 0, 0, 0, 0, 0, 7, 8, 9, 201}
				b513InformationViewBytes[4+b513InformationViewIndex] = 1 << b513InformationViewBit
				b513InformationViewOutput := b513InformationViewRun(b513InformationViewT, b513InformationViewInput(b513InformationViewBytes))
				if b513InformationViewOutput["state"] != "Hypothesis" {
					b513InformationViewT.Fatal(b513InformationViewOutput)
				}
				for b513InformationViewOther, b513InformationViewOtherName := range b513InformationViewNames {
					b513InformationViewExpected := map[string]any{}
					for b513InformationViewLabelIndex, b513InformationViewLabel := range b513InformationViewLabels[b513InformationViewOther] {
						b513InformationViewExpected[b513InformationViewLabel] = b513InformationViewOther == b513InformationViewIndex && b513InformationViewLabelIndex == b513InformationViewBit
					}
					b513InformationViewFieldValue := b513InformationViewField(b513InformationViewOutput, b513InformationViewOtherName)
					b513InformationViewReserved := 0
					if b513InformationViewOther == b513InformationViewIndex && b513InformationViewBit >= len(b513InformationViewLabels[b513InformationViewOther]) {
						b513InformationViewReserved = 1 << b513InformationViewBit
					}
					if !reflect.DeepEqual(b513InformationViewFieldValue["flags"], b513InformationViewExpected) || b513InformationViewFieldValue["reserved"] != float64(b513InformationViewReserved) || b513InformationViewFieldValue["state"] != "Hypothesis" || b513InformationViewFieldValue["number"] != float64(b513InformationViewBytes[4+b513InformationViewOther]) {
						b513InformationViewT.Fatal(b513InformationViewOtherName, b513InformationViewFieldValue)
					}
				}
				b513InformationViewOpaque := b513InformationViewField(b513InformationViewOutput, "opaque")
				if !reflect.DeepEqual(b513InformationViewOpaque, map[string]any{"raw": []any{float64(7), float64(8), float64(9)}, "state": "Unknown"}) {
					b513InformationViewT.Fatal(b513InformationViewOpaque)
				}
			})
		}
	}
}
func TestAb17Ao17EnumsAndPartial(b513InformationViewT *testing.T) {
	for _, b513InformationViewCase := range []struct {
		b513InformationViewName     string
		b513InformationViewPosition int
		b513InformationViewValues   []int
		b513InformationViewLabels   []string
	}{{"type", 3, []int{0, 1, 2, 3, 255}, []string{"modulating", "single-step", "two-step"}}, {"pump_mode", 8, []int{0, 1, 2}, []string{"pump-mode0", "pump-mode1"}}} {
		for _, b513InformationViewValue := range b513InformationViewCase.b513InformationViewValues {
			b513InformationViewT.Run(fmt.Sprintf("%s_%d", b513InformationViewCase.b513InformationViewName, b513InformationViewValue), func(b513InformationViewT *testing.T) {
				b513InformationViewBytes := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 9}
				b513InformationViewBytes[b513InformationViewCase.b513InformationViewPosition] = b513InformationViewValue
				b513InformationViewOutput := b513InformationViewRun(b513InformationViewT, b513InformationViewInput(b513InformationViewBytes))
				b513InformationViewActual := b513InformationViewField(b513InformationViewOutput, b513InformationViewCase.b513InformationViewName)
				if b513InformationViewValue < len(b513InformationViewCase.b513InformationViewLabels) {
					if b513InformationViewActual["number"] != float64(b513InformationViewValue) || b513InformationViewActual["label"] != b513InformationViewCase.b513InformationViewLabels[b513InformationViewValue] || b513InformationViewActual["state"] != "Hypothesis" {
						b513InformationViewT.Fatal(b513InformationViewActual)
					}
				} else {
					if b513InformationViewActual["state"] != "Unknown" || b513InformationViewActual["label"] != nil || b513InformationViewActual["number"] != nil || b513InformationViewOutput["state"] != "partial" {
						b513InformationViewT.Fatal(b513InformationViewActual)
					}
				}
				if b513InformationViewField(b513InformationViewOutput, "fuel")["state"] != "Hypothesis" {
					b513InformationViewT.Fatal("unknown sibling erased fuel")
				}
			})
		}
	}
	for b513InformationViewLength := 0; b513InformationViewLength < 12; b513InformationViewLength++ {
		b513InformationViewT.Run(fmt.Sprintf("length_%d", b513InformationViewLength), func(b513InformationViewT *testing.T) {
			b513InformationViewBytes := []int{201, 202, 203, 2, 255, 7, 15, 63, 1, 7, 8, 9}
			b513InformationViewOutput := b513InformationViewRun(b513InformationViewT, b513InformationViewInput(b513InformationViewBytes[:b513InformationViewLength]))
			if b513InformationViewOutput["state"] != "partial" {
				b513InformationViewT.Fatal(b513InformationViewOutput)
			}
			for b513InformationViewIndex, b513InformationViewName := range []string{"type", "capabilities", "fuel", "storage", "features", "pump_mode"} {
				b513InformationViewActual := b513InformationViewField(b513InformationViewOutput, b513InformationViewName)
				if b513InformationViewLength > b513InformationViewIndex+3 {
					if b513InformationViewActual["number"] != float64(b513InformationViewBytes[b513InformationViewIndex+3]) || b513InformationViewActual["state"] != "Hypothesis" {
						b513InformationViewT.Fatal(b513InformationViewName, b513InformationViewActual)
					}
				} else if b513InformationViewActual["number"] != nil || b513InformationViewActual["state"] != "Unknown" {
					b513InformationViewT.Fatal(b513InformationViewActual)
				}
			}
		})
	}
}
func TestAb17Ao17Feedback(b513InformationViewT *testing.T) {
	for _, b513InformationViewOperation := range []int{5, 6} {
		for _, b513InformationViewValue := range []int{0, 1, 2, 255} {
			b513InformationViewT.Run(fmt.Sprintf("op%d_%d", b513InformationViewOperation, b513InformationViewValue), func(b513InformationViewT *testing.T) {
				b513InformationViewInputValue := b513InformationViewInput([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 9})
				b513InformationViewContext := b513InformationViewInputValue["context"].(map[string]any)
				b513InformationViewContext["operation"] = b513InformationViewOperation
				b513InformationViewContext["view"] = "feedback"
				b513InformationViewInputValue["fields"].(map[string]any)["feedback"] = []int{b513InformationViewValue}
				b513InformationViewOutput := b513InformationViewRun(b513InformationViewT, b513InformationViewInputValue)
				b513InformationViewActual := b513InformationViewField(b513InformationViewOutput, "feedback")
				b513InformationViewExpected := "rejected"
				if b513InformationViewOperation == 5 {
					b513InformationViewExpected = "not-accepted-or-not-supported"
				}
				if b513InformationViewValue == 1 {
					b513InformationViewExpected = "accepted"
				}
				if b513InformationViewActual["label"] != b513InformationViewExpected || b513InformationViewActual["number"] != float64(b513InformationViewValue) || b513InformationViewActual["state"] != "Hypothesis" || len(b513InformationViewOutput["fields"].(map[string]any)) != 1 {
					b513InformationViewT.Fatal(b513InformationViewOutput)
				}
				for _, b513InformationViewWidth := range []int{0, 2, 3} {
					b513InformationViewInputValue["fields"].(map[string]any)["feedback"] = make([]int, b513InformationViewWidth)
					b513InformationViewOutput = b513InformationViewRun(b513InformationViewT, b513InformationViewInputValue)
					if b513InformationViewOutput["state"] != "partial" || b513InformationViewField(b513InformationViewOutput, "feedback")["number"] != nil {
						b513InformationViewT.Fatal(b513InformationViewOutput)
					}
				}
				delete(b513InformationViewInputValue["fields"].(map[string]any), "feedback")
				if b513InformationViewRun(b513InformationViewT, b513InformationViewInputValue)["state"] != "partial" {
					b513InformationViewT.Fatal("missing feedback")
				}
			})
		}
	}
}
func TestAb17Ao17Selection(b513InformationViewT *testing.T) {
	for b513InformationViewKey, b513InformationViewValues := range map[string][]any{"primary": {0, 180, 182}, "secondary": {0, 18, 20}, "operation": {0, 4, 5, 6, 8}, "direction": {"request", "unknown", ""}, "view": {"feedback", "unknown", ""}} {
		for _, b513InformationViewValue := range b513InformationViewValues {
			b513InformationViewT.Run(fmt.Sprint(b513InformationViewKey, b513InformationViewValue), func(b513InformationViewT *testing.T) {
				b513InformationViewInputValue := b513InformationViewInput([]int{181, 19, 7, 0, 255, 7, 15, 63, 1, 7, 8, 9})
				b513InformationViewInputValue["context"].(map[string]any)[b513InformationViewKey] = b513InformationViewValue
				b513InformationViewOutput := b513InformationViewRun(b513InformationViewT, b513InformationViewInputValue)
				if b513InformationViewOutput["state"] != "Unknown" || len(b513InformationViewOutput["fields"].(map[string]any)) != 0 {
					b513InformationViewT.Fatal(b513InformationViewOutput)
				}
			})
		}
	}
	for _, b513InformationViewKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		b513InformationViewInputValue := b513InformationViewInput([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 9})
		delete(b513InformationViewInputValue["context"].(map[string]any), b513InformationViewKey)
		if b513InformationViewRun(b513InformationViewT, b513InformationViewInputValue)["state"] != "Unknown" {
			b513InformationViewT.Fatal(b513InformationViewKey)
		}
	}
	for _, b513InformationViewKey := range []string{"family", "variant"} {
		for _, b513InformationViewValue := range []any{nil, 0, "other", map[string]any{"opaque": []int{0, 255}}} {
			b513InformationViewInputValue := b513InformationViewInput([]int{255, 254, 253, 0, 0, 0, 0, 0, 0, 7, 8, 9})
			b513InformationViewInputValue["context"].(map[string]any)[b513InformationViewKey] = b513InformationViewValue
			if b513InformationViewRun(b513InformationViewT, b513InformationViewInputValue)["state"] != "Hypothesis" {
				b513InformationViewT.Fatal("opaque family selected interpretation")
			}
		}
	}
	for _, b513InformationViewKey := range []string{"logical_object", "unmodelled"} {
		b513InformationViewInputValue := b513InformationViewInput([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 9})
		b513InformationViewInputValue["fields"].(map[string]any)[b513InformationViewKey] = []int{}
		b513InformationViewRun(b513InformationViewT, b513InformationViewInputValue)
		delete(b513InformationViewInputValue["fields"].(map[string]any), b513InformationViewKey)
		b513InformationViewRun(b513InformationViewT, b513InformationViewInputValue)
	}
	for _, b513InformationViewOperation := range []int{5, 6} {
		b513InformationViewInputValue := b513InformationViewInput([]int{})
		b513InformationViewContext := b513InformationViewInputValue["context"].(map[string]any)
		b513InformationViewContext["operation"] = b513InformationViewOperation
		b513InformationViewContext["view"] = "feedback"
		b513InformationViewContext["direction"] = "request"
		b513InformationViewInputValue["fields"].(map[string]any)["feedback"] = []int{1}
		if b513InformationViewRun(b513InformationViewT, b513InformationViewInputValue)["state"] != "Unknown" {
			b513InformationViewT.Fatal("request accepted")
		}
	}
}
func TestAb17Ao17Invalid(b513InformationViewT *testing.T) {
	for b513InformationViewIndex, b513InformationViewInputValue := range []string{``, `null`, `[]`, `{}`, `{"context":{},"fields":{}}{}`, `{"context":{},"fields":{},}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"family":{"x":0,"x":1}},"fields":{}}`, `{"context":{},"fields":{"x":[],"x":[]}}`, `{"context":{},"fields":{},"arguments":{"x":0,"x":1}}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":null}`, `{"context":{},"fields":{"x":null}}`, `{"context":{},"fields":{"x":"AA=="}}`, `{"context":{},"fields":{"x":[-1]}}`, `{"context":{},"fields":{"x":[256]}}`, `{"context":{},"fields":{"x":[1.0]}}`, `{"context":{},"fields":{"x":[1e0]}}`, `{"context":{},"fields":{"x":[true]}}`, `{"context":{},"fields":{"x":["0"]}}`, `{"context":{},"fields":{"x":[{}]}}`, `{"context":{},"fields":{"x":[[]]}}`, `{"context":{},"fields":{},"arguments":null}`, `{"context":{},"fields":{},"arguments":{"flag":1}}`, `{"context":{"primary":"181"},"fields":{}}`, `{"context":{"secondary":19.0},"fields":{}}`, `{"context":{"operation":true},"fields":{}}`, `{"context":{"direction":0},"fields":{}}`, `{"context":{"view":null},"fields":{}}`} {
		b513InformationViewT.Run(fmt.Sprint(b513InformationViewIndex), func(b513InformationViewT *testing.T) {
			b513InformationViewOutput, b513InformationViewError := B513InformationViewInterpret(json.RawMessage(b513InformationViewInputValue))
			if b513InformationViewError == nil || b513InformationViewOutput != nil {
				b513InformationViewT.Fatal("invalid accepted", b513InformationViewInputValue, string(b513InformationViewOutput))
			}
		})
	}
}
func TestAb17Ao17Isolation(b513InformationViewT *testing.T) {
	b513InformationViewInputValue := b513InformationViewInput([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 9, 255})
	b513InformationViewInputValue["arguments"] = map[string]any{}
	b513InformationViewEncoded, _ := json.Marshal(b513InformationViewInputValue)
	b513InformationViewOriginal := append(json.RawMessage(nil), b513InformationViewEncoded...)
	b513InformationViewFirst, b513InformationViewError := B513InformationViewInterpret(b513InformationViewEncoded)
	if b513InformationViewError != nil {
		b513InformationViewT.Fatal(b513InformationViewError)
	}
	b513InformationViewSaved := append(json.RawMessage(nil), b513InformationViewFirst...)
	for b513InformationViewIndex := range b513InformationViewEncoded {
		b513InformationViewEncoded[b513InformationViewIndex] = ' '
	}
	if !reflect.DeepEqual(b513InformationViewSaved, b513InformationViewFirst) {
		b513InformationViewT.Fatal("input aliases output")
	}
	b513InformationViewDecoded := b513InformationViewRun(b513InformationViewT, b513InformationViewInputValue)
	b513InformationViewDecoded["raw"].(map[string]any)["logical_object"].([]any)[3] = float64(255)
	if b513InformationViewField(b513InformationViewDecoded, "type")["raw"].([]any)[0] != float64(0) {
		b513InformationViewT.Fatal("raw aliases semantic field")
	}
	b513InformationViewDecoded["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "changed"
	b513InformationViewField(b513InformationViewDecoded, "fuel")["flags"].(map[string]any)["gas"] = true
	for b513InformationViewIndex := range b513InformationViewFirst {
		b513InformationViewFirst[b513InformationViewIndex] = ' '
	}
	b513InformationViewSecond, b513InformationViewError := B513InformationViewInterpret(b513InformationViewOriginal)
	if b513InformationViewError != nil || !reflect.DeepEqual(b513InformationViewSecond, b513InformationViewSaved) {
		b513InformationViewT.Fatal("later call changed", b513InformationViewError)
	}
}
func TestAb17Ao17Fixtures(b513InformationViewT *testing.T) {
	b513InformationViewData, b513InformationViewError := os.ReadFile("testdata/ab17ao17_b513_information_views.json")
	if b513InformationViewError != nil {
		b513InformationViewT.Fatal(b513InformationViewError)
	}
	var b513InformationViewCases []struct {
		B513InformationViewID    string          `json:"id"`
		B513InformationViewInput json.RawMessage `json:"input"`
		B513InformationViewState string          `json:"state"`
	}
	if b513InformationViewError = json.Unmarshal(b513InformationViewData, &b513InformationViewCases); b513InformationViewError != nil || len(b513InformationViewCases) == 0 {
		b513InformationViewT.Fatal("invalid fixtures", b513InformationViewError)
	}
	for _, b513InformationViewCase := range b513InformationViewCases {
		if !strings.HasPrefix(b513InformationViewCase.B513InformationViewID, "ab17ao17") {
			b513InformationViewT.Fatal("fixture ID")
		}
		b513InformationViewT.Run(b513InformationViewCase.B513InformationViewID, func(b513InformationViewT *testing.T) {
			var b513InformationViewInputValue any
			_ = json.Unmarshal(b513InformationViewCase.B513InformationViewInput, &b513InformationViewInputValue)
			if b513InformationViewRun(b513InformationViewT, b513InformationViewInputValue)["state"] != b513InformationViewCase.B513InformationViewState {
				b513InformationViewT.Fatal("fixture state")
			}
		})
	}
}

func TestAb17Ao17RawAxes(b513InformationViewT *testing.T) {
	for _, b513InformationViewPosition := range []int{0, 1, 2, 9, 10, 11, 12} {
		for _, b513InformationViewValue := range []int{0, 255} {
			b513InformationViewT.Run(fmt.Sprintf("position%d_%d", b513InformationViewPosition, b513InformationViewValue), func(b513InformationViewT *testing.T) {
				b513InformationViewBytes := []int{91, 82, 73, 2, 255, 7, 15, 63, 1, 7, 8, 9, 201}
				b513InformationViewBytes[b513InformationViewPosition] = b513InformationViewValue
				b513InformationViewOutput := b513InformationViewRun(b513InformationViewT, b513InformationViewInput(b513InformationViewBytes))
				if b513InformationViewOutput["state"] != "Hypothesis" {
					b513InformationViewT.Fatal("opaque position influenced admission")
				}
				b513InformationViewOpaque := b513InformationViewField(b513InformationViewOutput, "opaque")
				if !reflect.DeepEqual(b513InformationViewOpaque["raw"], []any{float64(b513InformationViewBytes[9]), float64(b513InformationViewBytes[10]), float64(b513InformationViewBytes[11])}) || b513InformationViewOpaque["number"] != nil || b513InformationViewOpaque["label"] != nil || b513InformationViewOpaque["state"] != "Unknown" {
					b513InformationViewT.Fatal(b513InformationViewOpaque)
				}
				for b513InformationViewIndex, b513InformationViewName := range []string{"type", "capabilities", "fuel", "storage", "features", "pump_mode"} {
					if !reflect.DeepEqual(b513InformationViewField(b513InformationViewOutput, b513InformationViewName)["raw"], []any{float64(b513InformationViewBytes[3+b513InformationViewIndex])}) {
						b513InformationViewT.Fatal("field raw differs", b513InformationViewName)
					}
				}
			})
		}
	}
	for _, b513InformationViewValue := range [][]int{{}, {0}, {255}, {0, 255, 1}} {
		b513InformationViewInputValue := b513InformationViewInput([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 9})
		b513InformationViewInputValue["fields"].(map[string]any)["unmodelled"] = b513InformationViewValue
		if b513InformationViewRun(b513InformationViewT, b513InformationViewInputValue)["state"] != "Hypothesis" {
			b513InformationViewT.Fatal("unmodelled bytes selected view")
		}
	}
	b513InformationViewInputValue := b513InformationViewInput([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 9})
	b513InformationViewInputValue["context"].(map[string]any)["additional"] = map[string]any{"nested": []any{nil, true, "opaque", 0}}
	b513InformationViewRun(b513InformationViewT, b513InformationViewInputValue)
}
