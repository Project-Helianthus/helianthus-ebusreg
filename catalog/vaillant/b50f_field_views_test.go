package vaillant

import (
	b50FFieldViewJSON "encoding/json"
	b50FFieldViewOS "os"
	b50FFieldViewReflect "reflect"
	b50FFieldViewStrings "strings"
	b50FFieldViewTesting "testing"
)

func b50FFieldViewInput() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 15, "operation": 0, "direction": "positive", "view": "thermal-fields", "family": map[string]any{"opaque": []any{"synthetic"}}, "variant": []any{0, "candidate"}}, "fields": map[string]any{"temperature_one": []int{240, 255}, "temperature_two": []int{0, 0}, "temperature_three": []int{16, 0}, "error_one": []int{0}, "error_two": []int{1}, "error_three": []int{2}}}
}
func b50FFieldViewCall(b50FFieldViewT *b50FFieldViewTesting.T, b50FFieldViewInputValue any) map[string]any {
	b50FFieldViewT.Helper()
	b50FFieldViewData, b50FFieldViewErr := b50FFieldViewJSON.Marshal(b50FFieldViewInputValue)
	if b50FFieldViewErr != nil {
		b50FFieldViewT.Fatal(b50FFieldViewErr)
	}
	b50FFieldViewData, b50FFieldViewErr = B50FFieldViewInterpret(b50FFieldViewData)
	if b50FFieldViewErr != nil {
		b50FFieldViewT.Fatal(b50FFieldViewErr)
	}
	var b50FFieldViewOutput map[string]any
	if b50FFieldViewErr = b50FFieldViewJSON.Unmarshal(b50FFieldViewData, &b50FFieldViewOutput); b50FFieldViewErr != nil {
		b50FFieldViewT.Fatal(b50FFieldViewErr)
	}
	if b50FFieldViewOutput["physical"] != "Unknown" || b50FFieldViewOutput["qualified"] != false {
		b50FFieldViewT.Fatal("qualification promoted")
	}
	return b50FFieldViewOutput
}
func b50FFieldViewResult(b50FFieldViewOutput map[string]any, b50FFieldViewKey string) map[string]any {
	return b50FFieldViewOutput["fields"].(map[string]any)[b50FFieldViewKey].(map[string]any)
}
func TestAb09Ao09Thermal(b50FFieldViewT *b50FFieldViewTesting.T) {
	for b50FFieldViewIndex, b50FFieldViewKey := range []string{"temperature_one", "temperature_two", "temperature_three", "error_one", "error_two", "error_three"} {
		for _, b50FFieldViewBytes := range [][]int{nil, {}, {0}, {16, 0}, {255, 255, 99}} {
			b50FFieldViewInputValue := b50FFieldViewInput()
			b50FFieldViewRaw := b50FFieldViewInputValue["fields"].(map[string]any)
			if b50FFieldViewBytes == nil {
				delete(b50FFieldViewRaw, b50FFieldViewKey)
			} else {
				b50FFieldViewRaw[b50FFieldViewKey] = b50FFieldViewBytes
			}
			b50FFieldViewOutput := b50FFieldViewCall(b50FFieldViewT, b50FFieldViewInputValue)
			b50FFieldViewField := b50FFieldViewResult(b50FFieldViewOutput, b50FFieldViewKey)
			b50FFieldViewWidth := 1
			if b50FFieldViewIndex < 3 {
				b50FFieldViewWidth = 2
			}
			if len(b50FFieldViewBytes) != b50FFieldViewWidth {
				if _, b50FFieldViewFound := b50FFieldViewField["number"]; b50FFieldViewFound {
					b50FFieldViewT.Fatal("numeric incomplete field")
				}
				if b50FFieldViewOutput["state"] != "partial" {
					b50FFieldViewT.Fatal("lost partial")
				}
			}
			if b50FFieldViewBytes == nil {
				if _, b50FFieldViewFound := b50FFieldViewField["raw"]; b50FFieldViewFound || b50FFieldViewField["state"] != "missing" {
					b50FFieldViewT.Fatal("lost absence")
				}
			} else {
				b50FFieldViewExpected, _ := b50FFieldViewJSON.Marshal(b50FFieldViewBytes)
				b50FFieldViewActual, _ := b50FFieldViewJSON.Marshal(b50FFieldViewField["raw"])
				if string(b50FFieldViewExpected) != string(b50FFieldViewActual) {
					b50FFieldViewT.Fatal("raw alias", b50FFieldViewKey)
				}
			}
			for b50FFieldViewSibling, b50FFieldViewExpected := range map[string]float64{"temperature_one": -1, "temperature_two": 0, "temperature_three": 1} {
				if b50FFieldViewSibling != b50FFieldViewKey {
					b50FFieldViewField = b50FFieldViewResult(b50FFieldViewOutput, b50FFieldViewSibling)
					if b50FFieldViewField["number"] != b50FFieldViewExpected || b50FFieldViewField["grade"] != "Hypothesis" || b50FFieldViewField["calculation"] != "Proven" {
						b50FFieldViewT.Fatal("independent arithmetic", b50FFieldViewSibling, b50FFieldViewField)
					}
				}
			}
		}
	}
	b50FFieldViewInputValue := b50FFieldViewInput()
	b50FFieldViewInputValue["context"].(map[string]any)["operation"] = 1
	b50FFieldViewRaw := b50FFieldViewInputValue["fields"].(map[string]any)
	delete(b50FFieldViewRaw, "temperature_three")
	for _, b50FFieldViewBytes := range [][]int{nil, {16, 0}, {0, 0}, {}, {1}, {1, 2, 3}} {
		if b50FFieldViewBytes == nil {
			delete(b50FFieldViewRaw, "stack_voltage")
		} else {
			b50FFieldViewRaw["stack_voltage"] = b50FFieldViewBytes
		}
		b50FFieldViewOutput := b50FFieldViewCall(b50FFieldViewT, b50FFieldViewInputValue)
		b50FFieldViewField := b50FFieldViewResult(b50FFieldViewOutput, "stack_voltage")
		if _, b50FFieldViewFound := b50FFieldViewField["number"]; b50FFieldViewFound {
			b50FFieldViewT.Fatal("invented voltage")
		}
		if len(b50FFieldViewBytes) == 2 && b50FFieldViewField["state"] != "unit-conflict" {
			b50FFieldViewT.Fatal("lost conflict")
		}
		if b50FFieldViewResult(b50FFieldViewOutput, "temperature_one")["number"] != float64(-1) {
			b50FFieldViewT.Fatal("lost sibling")
		}
	}
}
func TestAb09Ao09DemandIntegrity(b50FFieldViewT *b50FFieldViewTesting.T) {
	for _, b50FFieldViewDemand := range []int{0, 1, 2, 4, 8, 16, 31, 32, 64, 128, 224, 255} {
		for b50FFieldViewFeedback := 0; b50FFieldViewFeedback < 256; b50FFieldViewFeedback++ {
			for _, b50FFieldViewIntegrity := range [][]int{nil, {}, {18}, {52}, {18, 52}} {
				b50FFieldViewInputValue := b50FFieldViewInput()
				b50FFieldViewContext := b50FFieldViewInputValue["context"].(map[string]any)
				b50FFieldViewContext["operation"] = 3
				b50FFieldViewContext["view"] = "demand"
				b50FFieldViewRaw := map[string]any{"demand": []int{b50FFieldViewDemand}, "feedback": []int{b50FFieldViewFeedback}, "integrity_two": []int{52}}
				if b50FFieldViewIntegrity != nil {
					b50FFieldViewRaw["integrity_one"] = b50FFieldViewIntegrity
				}
				b50FFieldViewInputValue["fields"] = b50FFieldViewRaw
				b50FFieldViewOutput := b50FFieldViewCall(b50FFieldViewT, b50FFieldViewInputValue)
				b50FFieldViewField := b50FFieldViewResult(b50FFieldViewOutput, "demand")
				if b50FFieldViewField["label"] != "demand" || b50FFieldViewField["reserved"] != float64(b50FFieldViewDemand&224) || b50FFieldViewField["grade"] != "Hypothesis" {
					b50FFieldViewT.Fatal("demand promoted or reserved lost")
				}
				for b50FFieldViewIndex, b50FFieldViewName := range []string{"burner", "reformer-gas-valve", "gas-flow", "fan-power", "valve"} {
					if b50FFieldViewField["flags"].(map[string]any)[b50FFieldViewName] != (b50FFieldViewDemand&(1<<b50FFieldViewIndex) != 0) {
						b50FFieldViewT.Fatal("demand bit", b50FFieldViewName)
					}
				}
				b50FFieldViewLabel := "rejected"
				if b50FFieldViewFeedback == 1 {
					b50FFieldViewLabel = "accepted"
				}
				if b50FFieldViewResult(b50FFieldViewOutput, "feedback")["label"] != b50FFieldViewLabel {
					b50FFieldViewT.Fatal("feedback depends on integrity")
				}
				b50FFieldViewField = b50FFieldViewResult(b50FFieldViewOutput, "integrity_two")
				if b50FFieldViewField["crc"] != "Unknown" || b50FFieldViewField["word_order"] != "Unknown" || b50FFieldViewField["logical_position"] != float64(5) {
					b50FFieldViewT.Fatal("invented CRC")
				}
				if _, b50FFieldViewFound := b50FFieldViewField["number"]; b50FFieldViewFound {
					b50FFieldViewT.Fatal("numeric integrity")
				}
			}
		}
	}
	for _, b50FFieldViewOperation := range []int{0, 1, 2, 3} {
		b50FFieldViewInputValue := b50FFieldViewInput()
		b50FFieldViewContext := b50FFieldViewInputValue["context"].(map[string]any)
		b50FFieldViewContext["operation"] = b50FFieldViewOperation
		if b50FFieldViewOperation == 2 {
			b50FFieldViewContext["view"] = "integrity-only"
		}
		if b50FFieldViewOperation == 3 {
			b50FFieldViewContext["view"] = "demand"
		}
		b50FFieldViewRaw := b50FFieldViewInputValue["fields"].(map[string]any)
		b50FFieldViewRaw["integrity_one"] = []int{18}
		b50FFieldViewRaw["integrity_two"] = []int{52}
		b50FFieldViewOutput := b50FFieldViewCall(b50FFieldViewT, b50FFieldViewInputValue)
		b50FFieldViewPosition := 12
		if b50FFieldViewOperation == 2 {
			b50FFieldViewPosition = 11
		}
		if b50FFieldViewOperation == 3 {
			b50FFieldViewPosition = 4
		}
		for b50FFieldViewIndex, b50FFieldViewKey := range []string{"integrity_one", "integrity_two"} {
			b50FFieldViewField := b50FFieldViewResult(b50FFieldViewOutput, b50FFieldViewKey)
			if b50FFieldViewField["logical_position"] != float64(b50FFieldViewPosition+b50FFieldViewIndex) || b50FFieldViewField["state"] != "unknown-order" {
				b50FFieldViewT.Fatal("integrity position")
			}
		}
		if b50FFieldViewOperation == 2 && b50FFieldViewResult(b50FFieldViewOutput, "temperature_one")["state"] != "Unknown" {
			b50FFieldViewT.Fatal("borrowed thermal")
		}
	}
}
func TestAb09Ao09SelectorsCopies(b50FFieldViewT *b50FFieldViewTesting.T) {
	for b50FFieldViewKey, b50FFieldViewValue := range map[string]any{"primary": 180, "secondary": 14, "operation": 9, "direction": "negative", "view": "unknown", "family": map[string]any{"other": []int{1, 2}}, "variant": false} {
		b50FFieldViewInputValue := b50FFieldViewInput()
		b50FFieldViewInputValue["context"].(map[string]any)[b50FFieldViewKey] = b50FFieldViewValue
		b50FFieldViewOutput := b50FFieldViewCall(b50FFieldViewT, b50FFieldViewInputValue)
		b50FFieldViewExpected, _ := b50FFieldViewJSON.Marshal(b50FFieldViewInputValue["context"])
		b50FFieldViewActual, _ := b50FFieldViewJSON.Marshal(b50FFieldViewOutput["context"])
		if string(b50FFieldViewExpected) != string(b50FFieldViewActual) {
			b50FFieldViewT.Fatal("context changed")
		}
		if b50FFieldViewKey != "family" && b50FFieldViewKey != "variant" {
			if b50FFieldViewOutput["state"] != "Unknown" || b50FFieldViewResult(b50FFieldViewOutput, "temperature_one")["state"] != "Unknown" {
				b50FFieldViewT.Fatal("guessed selection", b50FFieldViewKey)
			}
		} else if b50FFieldViewResult(b50FFieldViewOutput, "temperature_one")["number"] != float64(-1) {
			b50FFieldViewT.Fatal("family selected view")
		}
	}
	for _, b50FFieldViewKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		b50FFieldViewInputValue := b50FFieldViewInput()
		delete(b50FFieldViewInputValue["context"].(map[string]any), b50FFieldViewKey)
		if b50FFieldViewCall(b50FFieldViewT, b50FFieldViewInputValue)["state"] != "Unknown" {
			b50FFieldViewT.Fatal("missing selector", b50FFieldViewKey)
		}
	}
	b50FFieldViewInputValue := b50FFieldViewInput()
	b50FFieldViewInputValue["arguments"] = map[string]any{}
	b50FFieldViewInputValue["fields"].(map[string]any)["tail"] = []int{1, 2, 3}
	b50FFieldViewData, _ := b50FFieldViewJSON.Marshal(b50FFieldViewInputValue)
	b50FFieldViewSaved := append([]byte(nil), b50FFieldViewData...)
	b50FFieldViewFirst, b50FFieldViewErr := B50FFieldViewInterpret(b50FFieldViewData)
	if b50FFieldViewErr != nil {
		b50FFieldViewT.Fatal(b50FFieldViewErr)
	}
	for b50FFieldViewIndex := range b50FFieldViewData {
		b50FFieldViewData[b50FFieldViewIndex] = ' '
	}
	b50FFieldViewSecond, b50FFieldViewErr := B50FFieldViewInterpret(b50FFieldViewSaved)
	if b50FFieldViewErr != nil || string(b50FFieldViewFirst) != string(b50FFieldViewSecond) {
		b50FFieldViewT.Fatal("input alias")
	}
	b50FFieldViewFirst[0] = '!'
	b50FFieldViewThird, _ := B50FFieldViewInterpret(b50FFieldViewSaved)
	if string(b50FFieldViewThird) != string(b50FFieldViewSecond) {
		b50FFieldViewT.Fatal("output alias")
	}
	b50FFieldViewOutput := b50FFieldViewCall(b50FFieldViewT, b50FFieldViewInputValue)
	b50FFieldViewResult(b50FFieldViewOutput, "tail")["raw"].([]any)[0] = float64(99)
	if b50FFieldViewOutput["raw"].(map[string]any)["tail"].([]any)[0] != float64(1) {
		b50FFieldViewT.Fatal("nested alias")
	}
	if _, b50FFieldViewFound := b50FFieldViewOutput["arguments"]; !b50FFieldViewFound {
		b50FFieldViewT.Fatal("empty arguments lost")
	}
	delete(b50FFieldViewInputValue, "arguments")
	if _, b50FFieldViewFound := b50FFieldViewCall(b50FFieldViewT, b50FFieldViewInputValue)["arguments"]; b50FFieldViewFound {
		b50FFieldViewT.Fatal("arguments invented")
	}
}
func TestAb09Ao09Reject(b50FFieldViewT *b50FFieldViewTesting.T) {
	for _, b50FFieldViewText := range []string{`null`, `[]`, `{}`, `{"context":{},"fields":null}`, `{"context":{},"fields":{"x":null}}`, `{"context":{},"fields":{"x":[256]}}`, `{"context":{},"fields":{"x":[-1]}}`, `{"context":{},"fields":{"x":[-0]}}`, `{"context":{},"fields":{"x":[1.0]}}`, `{"context":{},"fields":{"x":[true]}}`, `{"context":{},"fields":{"x":["1"]}}`, `{"context":{},"fields":{"x":[{}]}}`, `{"context":{},"fields":{"x":{}}}`, `{"context":{},"fields":{},"arguments":null}`, `{"context":{},"fields":{},"arguments":{"flag":true}}`, `{"context":{"primary":"181"},"fields":{}}`, `{"context":{"operation":0.1},"fields":{}}`, `{"context":{"direction":0},"fields":{}}`, `{"context":{"view":null},"fields":{}}`, `{"context":{},"fields":{},"extra":0}`, `{"context":{},"fields":{},"fields":{}}`, `{"context":{"family":{"x":1,"x":2}},"fields":{}}`, `{"context":{},"fields":{"x":[],"x":[]}}`, `{"context":{},"fields":{}} {}`, `{"context":{},"fields":{}`, `{"context":{},"fields":{"x":[1e999]}}`} {
		b50FFieldViewOutput, b50FFieldViewErr := B50FFieldViewInterpret([]byte(b50FFieldViewText))
		if b50FFieldViewErr == nil || b50FFieldViewOutput != nil {
			b50FFieldViewT.Fatal("accepted invalid input", b50FFieldViewText)
		}
	}
}
func TestAb09Ao09Fixtures(b50FFieldViewT *b50FFieldViewTesting.T) {
	b50FFieldViewData, b50FFieldViewErr := b50FFieldViewOS.ReadFile("testdata/ab09ao09_b50f_field_views.json")
	if b50FFieldViewErr != nil {
		b50FFieldViewT.Fatal(b50FFieldViewErr)
	}
	var b50FFieldViewCases []struct {
		B50FFieldViewID       string                       `json:"ab09ao09_id"`
		B50FFieldViewInput    b50FFieldViewJSON.RawMessage `json:"input"`
		B50FFieldViewExpected map[string]any               `json:"expected"`
	}
	if b50FFieldViewErr = b50FFieldViewJSON.Unmarshal(b50FFieldViewData, &b50FFieldViewCases); b50FFieldViewErr != nil || len(b50FFieldViewCases) == 0 {
		b50FFieldViewT.Fatal("invalid fixtures", b50FFieldViewErr)
	}
	for _, b50FFieldViewCase := range b50FFieldViewCases {
		b50FFieldViewT.Run(b50FFieldViewCase.B50FFieldViewID, func(b50FFieldViewT *b50FFieldViewTesting.T) {
			b50FFieldViewData, b50FFieldViewErr := B50FFieldViewInterpret(b50FFieldViewCase.B50FFieldViewInput)
			if b50FFieldViewErr != nil {
				b50FFieldViewT.Fatal(b50FFieldViewErr)
			}
			var b50FFieldViewOutput map[string]any
			if b50FFieldViewErr = b50FFieldViewJSON.Unmarshal(b50FFieldViewData, &b50FFieldViewOutput); b50FFieldViewErr != nil {
				b50FFieldViewT.Fatal(b50FFieldViewErr)
			}
			for b50FFieldViewKey, b50FFieldViewExpected := range b50FFieldViewCase.B50FFieldViewExpected {
				if !b50FFieldViewReflect.DeepEqual(b50FFieldViewOutput[b50FFieldViewKey], b50FFieldViewExpected) {
					b50FFieldViewT.Fatal(b50FFieldViewKey, b50FFieldViewOutput)
				}
			}
		})
	}
}

func TestAb09Ao09EveryFieldAxis(b50FFieldViewT *b50FFieldViewTesting.T) {
	for _, b50FFieldViewOperation := range []int{0, 1, 2, 3} {
		b50FFieldViewControl := b50FFieldViewInput()
		b50FFieldViewContext := b50FFieldViewControl["context"].(map[string]any)
		b50FFieldViewContext["operation"] = b50FFieldViewOperation
		b50FFieldViewRaw := b50FFieldViewControl["fields"].(map[string]any)
		if b50FFieldViewOperation == 1 {
			delete(b50FFieldViewRaw, "temperature_three")
			b50FFieldViewRaw["stack_voltage"] = []int{16, 0}
		}
		if b50FFieldViewOperation == 2 {
			b50FFieldViewContext["view"] = "integrity-only"
			b50FFieldViewRaw = map[string]any{}
		}
		if b50FFieldViewOperation == 3 {
			b50FFieldViewContext["view"] = "demand"
			b50FFieldViewRaw = map[string]any{"demand": []int{31}, "feedback": []int{1}}
		}
		b50FFieldViewRaw["integrity_one"] = []int{18}
		b50FFieldViewRaw["integrity_two"] = []int{52}
		b50FFieldViewControl["fields"] = b50FFieldViewRaw
		b50FFieldViewBaseline := b50FFieldViewCall(b50FFieldViewT, b50FFieldViewControl)
		for b50FFieldViewKey, b50FFieldViewOriginal := range b50FFieldViewRaw {
			for _, b50FFieldViewBytes := range [][]int{nil, {}, {0}, {255}, {0, 0}, {240, 255}, {1, 2, 3}} {
				if b50FFieldViewBytes == nil {
					delete(b50FFieldViewRaw, b50FFieldViewKey)
				} else {
					b50FFieldViewRaw[b50FFieldViewKey] = b50FFieldViewBytes
				}
				b50FFieldViewOutput := b50FFieldViewCall(b50FFieldViewT, b50FFieldViewControl)
				for b50FFieldViewSibling := range b50FFieldViewRaw {
					if b50FFieldViewSibling != b50FFieldViewKey && !b50FFieldViewReflect.DeepEqual(b50FFieldViewResult(b50FFieldViewOutput, b50FFieldViewSibling), b50FFieldViewResult(b50FFieldViewBaseline, b50FFieldViewSibling)) {
						b50FFieldViewT.Fatal("sibling coupling", b50FFieldViewOperation, b50FFieldViewKey, b50FFieldViewSibling)
					}
				}
				b50FFieldViewWidth := len(b50FFieldViewOriginal.([]int))
				b50FFieldViewFieldValue, b50FFieldViewPresent := b50FFieldViewOutput["fields"].(map[string]any)[b50FFieldViewKey]
				if !b50FFieldViewPresent {
					if b50FFieldViewBytes != nil || (b50FFieldViewKey != "integrity_one" && b50FFieldViewKey != "integrity_two") {
						b50FFieldViewT.Fatal("lost field")
					}
					continue
				}
				b50FFieldViewField := b50FFieldViewFieldValue.(map[string]any)
				if b50FFieldViewBytes == nil {
					if _, b50FFieldViewFound := b50FFieldViewField["raw"]; b50FFieldViewFound || b50FFieldViewField["state"] != "missing" {
						b50FFieldViewT.Fatal("absence collapsed")
					}
				} else {
					b50FFieldViewExpected, _ := b50FFieldViewJSON.Marshal(b50FFieldViewBytes)
					b50FFieldViewActual, _ := b50FFieldViewJSON.Marshal(b50FFieldViewField["raw"])
					if string(b50FFieldViewActual) != string(b50FFieldViewExpected) {
						b50FFieldViewT.Fatal("raw substituted", b50FFieldViewKey)
					}
					if len(b50FFieldViewBytes) != b50FFieldViewWidth {
						b50FFieldViewState := "incomplete"
						if len(b50FFieldViewBytes) > b50FFieldViewWidth {
							b50FFieldViewState = "unmapped"
						}
						if b50FFieldViewField["state"] != b50FFieldViewState {
							b50FFieldViewT.Fatal("width guessed", b50FFieldViewKey)
						}
						for _, b50FFieldViewSemantic := range []string{"number", "label", "flags", "reserved"} {
							if _, b50FFieldViewFound := b50FFieldViewField[b50FFieldViewSemantic]; b50FFieldViewFound {
								b50FFieldViewT.Fatal("incomplete semantics", b50FFieldViewKey)
							}
						}
					}
				}
			}
			b50FFieldViewRaw[b50FFieldViewKey] = b50FFieldViewOriginal
		}
	}
	for _, b50FFieldViewKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, b50FFieldViewValue := range []any{nil, true, []int{}, map[string]any{}} {
			b50FFieldViewInputValue := b50FFieldViewInput()
			b50FFieldViewInputValue["context"].(map[string]any)[b50FFieldViewKey] = b50FFieldViewValue
			b50FFieldViewData, _ := b50FFieldViewJSON.Marshal(b50FFieldViewInputValue)
			b50FFieldViewResult, b50FFieldViewErr := B50FFieldViewInterpret(b50FFieldViewData)
			if b50FFieldViewErr == nil || b50FFieldViewResult != nil {
				b50FFieldViewT.Fatal("wrong selector kind", b50FFieldViewKey)
			}
		}
	}
}

func TestAb09Ao09DepthAndRange(b50FFieldViewT *b50FFieldViewTesting.T) {
	b50FFieldViewDeep := `{"context":{"x":` + b50FFieldViewStrings.Repeat("[", 200000) + b50FFieldViewStrings.Repeat("]", 200000) + `},"fields":{}}`
	if b50FFieldViewOutput, b50FFieldViewErr := B50FFieldViewInterpret([]byte(b50FFieldViewDeep)); b50FFieldViewErr == nil || b50FFieldViewOutput != nil {
		b50FFieldViewT.Fatal("deep nesting accepted")
	}
	b50FFieldViewShallow := `{"context":{"x":` + b50FFieldViewStrings.Repeat("[", 32) + b50FFieldViewStrings.Repeat("]", 32) + `},"fields":{}}`
	if _, b50FFieldViewErr := B50FFieldViewInterpret([]byte(b50FFieldViewShallow)); b50FFieldViewErr != nil {
		b50FFieldViewT.Fatal(b50FFieldViewErr)
	}
	for _, b50FFieldViewKey := range []string{"operation", "primary", "secondary"} {
		b50FFieldViewData, b50FFieldViewErr := B50FFieldViewInterpret([]byte(`{"context":{"` + b50FFieldViewKey + `":5000000000000000000000,"direction":"positive"},"fields":{"temperature_one":[1,2]}}`))
		if b50FFieldViewErr != nil {
			b50FFieldViewT.Fatal(b50FFieldViewKey, b50FFieldViewErr)
		}
		var b50FFieldViewOutput map[string]any
		if b50FFieldViewErr = b50FFieldViewJSON.Unmarshal(b50FFieldViewData, &b50FFieldViewOutput); b50FFieldViewErr != nil || b50FFieldViewOutput["state"] != "Unknown" {
			b50FFieldViewT.Fatal(b50FFieldViewKey, "not Unknown")
		}
	}
}
