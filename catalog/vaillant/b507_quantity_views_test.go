package vaillant

import (
	b507QuantityViewBytesPackage "bytes"
	b507QuantityViewJSON "encoding/json"
	b507QuantityViewFmt "fmt"
	b507QuantityViewOs "os"
	b507QuantityViewReflect "reflect"
	b507QuantityViewTesting "testing"
)

func b507QuantityViewInput(b507QuantityViewOperation int, b507QuantityViewView string, b507QuantityViewFields map[string]any) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 7, "operation": b507QuantityViewOperation, "direction": "positive", "view": b507QuantityViewView, "family": map[string]any{"opaque": []any{"synthetic-family"}}, "variant": []any{"synthetic-variant"}}, "fields": b507QuantityViewFields}
}

func b507QuantityViewRun(b507QuantityViewT *b507QuantityViewTesting.T, b507QuantityViewObject map[string]any) map[string]any {
	b507QuantityViewT.Helper()
	b507QuantityViewInputBytes, b507QuantityViewError := b507QuantityViewJSON.Marshal(b507QuantityViewObject)
	if b507QuantityViewError != nil {
		b507QuantityViewT.Fatal(b507QuantityViewError)
	}
	b507QuantityViewOutput, b507QuantityViewError := B507QuantityViewInterpret(b507QuantityViewInputBytes)
	if b507QuantityViewError != nil {
		b507QuantityViewT.Fatal(b507QuantityViewError)
	}
	b507QuantityViewResult := map[string]any{}
	if b507QuantityViewError := b507QuantityViewJSON.Unmarshal(b507QuantityViewOutput, &b507QuantityViewResult); b507QuantityViewError != nil {
		b507QuantityViewT.Fatal(b507QuantityViewError)
	}
	if b507QuantityViewResult["physical"] != "Unknown" || b507QuantityViewResult["qualified"] != false {
		b507QuantityViewT.Fatalf("qualification: %s", b507QuantityViewOutput)
	}
	b507QuantityViewCanonical := map[string]any{}
	if b507QuantityViewError := b507QuantityViewJSON.Unmarshal(b507QuantityViewInputBytes, &b507QuantityViewCanonical); b507QuantityViewError != nil {
		b507QuantityViewT.Fatal(b507QuantityViewError)
	}
	if !b507QuantityViewReflect.DeepEqual(b507QuantityViewResult["context"], b507QuantityViewCanonical["context"]) || !b507QuantityViewReflect.DeepEqual(b507QuantityViewResult["raw"], b507QuantityViewCanonical["fields"]) {
		b507QuantityViewT.Fatalf("lost evidence: %s", b507QuantityViewOutput)
	}
	for b507QuantityViewKey := range b507QuantityViewResult {
		switch b507QuantityViewKey {
		case "context", "raw", "arguments", "state", "fields", "physical", "qualified":
		default:
			b507QuantityViewT.Fatalf("unexpected output %s", b507QuantityViewKey)
		}
	}
	for _, b507QuantityViewItem := range b507QuantityViewResult["fields"].(map[string]any) {
		for b507QuantityViewKey := range b507QuantityViewItem.(map[string]any) {
			switch b507QuantityViewKey {
			case "raw", "state", "grade", "number", "label", "calculation", "candidate_domain":
			default:
				b507QuantityViewT.Fatalf("invented meaning or permission %s", b507QuantityViewKey)
			}
		}
	}
	return b507QuantityViewResult
}

func TestAb05Ao05Capacity(b507QuantityViewT *b507QuantityViewTesting.T) {
	for b507QuantityViewIndex, b507QuantityViewByte := range []int{0, 1, 250, 251, 255} {
		b507QuantityViewT.Run(b507QuantityViewFmt.Sprint(b507QuantityViewByte), func(b507QuantityViewT *b507QuantityViewTesting.T) {
			b507QuantityViewResult := b507QuantityViewRun(b507QuantityViewT, b507QuantityViewInput(6, "capacity", map[string]any{"capacity": []int{b507QuantityViewByte}, "cost": []int{1, 2}, "tail": []int{99}}))
			b507QuantityViewField := b507QuantityViewResult["fields"].(map[string]any)["capacity"].(map[string]any)
			b507QuantityViewNumber, b507QuantityViewHasNumber := b507QuantityViewField["number"]
			if !b507QuantityViewHasNumber || b507QuantityViewNumber != []float64{0, 0.1, 25, 25.1, 25.5}[b507QuantityViewIndex] || b507QuantityViewField["calculation"] != "Proven" {
				b507QuantityViewT.Fatalf("capacity: %#v", b507QuantityViewField)
			}
			if b507QuantityViewField["grade"] != "Hypothesis" || b507QuantityViewField["candidate_domain"] != (b507QuantityViewByte <= 250) {
				b507QuantityViewT.Fatalf("capacity: %#v", b507QuantityViewField)
			}
			b507QuantityViewState, b507QuantityViewAggregate := "Hypothesis", "Hypothesis"
			if b507QuantityViewByte > 250 {
				b507QuantityViewState, b507QuantityViewAggregate = "outside-candidate-domain", "partial"
			}
			if b507QuantityViewField["state"] != b507QuantityViewState || b507QuantityViewResult["state"] != b507QuantityViewAggregate {
				b507QuantityViewT.Fatal(b507QuantityViewResult)
			}
			for _, b507QuantityViewKey := range []string{"cost", "tail"} {
				b507QuantityViewUnknown := b507QuantityViewResult["fields"].(map[string]any)[b507QuantityViewKey].(map[string]any)
				if len(b507QuantityViewUnknown) != 2 || b507QuantityViewUnknown["state"] != "Unknown" {
					b507QuantityViewT.Fatal(b507QuantityViewUnknown)
				}
			}
		})
	}
}

func TestAb05Ao05Preference(b507QuantityViewT *b507QuantityViewTesting.T) {
	for b507QuantityViewIndex, b507QuantityViewWord := range [][]int{{0, 0}, {255, 255}, {0, 1}, {1, 0}} {
		b507QuantityViewT.Run(b507QuantityViewFmt.Sprint(b507QuantityViewIndex), func(b507QuantityViewT *b507QuantityViewTesting.T) {
			b507QuantityViewResult := b507QuantityViewRun(b507QuantityViewT, b507QuantityViewInput(8, "preference", map[string]any{"preference": b507QuantityViewWord}))
			b507QuantityViewField := b507QuantityViewResult["fields"].(map[string]any)["preference"].(map[string]any)
			if _, b507QuantityViewPresent := b507QuantityViewField["number"]; b507QuantityViewPresent {
				b507QuantityViewT.Fatal(b507QuantityViewField)
			}
			if b507QuantityViewField["grade"] != "Hypothesis" {
				b507QuantityViewT.Fatal(b507QuantityViewField)
			}
			if b507QuantityViewIndex < 2 {
				if b507QuantityViewField["label"] != []string{"appliance-unavailable candidate", "maximum-efficiency-preference"}[b507QuantityViewIndex] || b507QuantityViewField["state"] != "Hypothesis" {
					b507QuantityViewT.Fatal(b507QuantityViewField)
				}
			} else if b507QuantityViewField["label"] != nil || b507QuantityViewField["state"] != "unmapped" || b507QuantityViewResult["state"] != "partial" {
				b507QuantityViewT.Fatal(b507QuantityViewResult)
			}
		})
	}
}

func TestAb05Ao05Selectors(b507QuantityViewT *b507QuantityViewTesting.T) {
	for _, b507QuantityViewView := range []string{"capacity", "preference"} {
		b507QuantityViewOperation := 6
		b507QuantityViewBytes := []int{1}
		if b507QuantityViewView == "preference" {
			b507QuantityViewOperation = 8
			b507QuantityViewBytes = []int{255, 255}
		}
		for b507QuantityViewKey, b507QuantityViewValues := range map[string][]any{"primary": {180, 0, b507QuantityViewJSON.Number("999999999999999999999999999999"), nil}, "secondary": {8, 0, b507QuantityViewJSON.Number("999999999999999999999999999999"), nil}, "operation": {9, 11, 0, b507QuantityViewJSON.Number("999999999999999999999999999999"), nil}, "direction": {"negative", "", nil}, "view": {"unknown", "", nil}} {
			for b507QuantityViewIndex, b507QuantityViewValue := range b507QuantityViewValues {
				b507QuantityViewT.Run(b507QuantityViewFmt.Sprintf("%s/%s/%d", b507QuantityViewView, b507QuantityViewKey, b507QuantityViewIndex), func(b507QuantityViewT *b507QuantityViewTesting.T) {
					b507QuantityViewObject := b507QuantityViewInput(b507QuantityViewOperation, b507QuantityViewView, map[string]any{b507QuantityViewView: b507QuantityViewBytes})
					if b507QuantityViewValue == nil {
						delete(b507QuantityViewObject["context"].(map[string]any), b507QuantityViewKey)
					} else {
						b507QuantityViewObject["context"].(map[string]any)[b507QuantityViewKey] = b507QuantityViewValue
					}
					b507QuantityViewResult := b507QuantityViewRun(b507QuantityViewT, b507QuantityViewObject)
					b507QuantityViewField := b507QuantityViewResult["fields"].(map[string]any)[b507QuantityViewView].(map[string]any)
					if b507QuantityViewResult["state"] != "Unknown" || b507QuantityViewField["state"] != "Unknown" || len(b507QuantityViewField) != 2 {
						b507QuantityViewT.Fatal(b507QuantityViewResult)
					}
				})
			}
		}
		for _, b507QuantityViewKey := range []string{"family", "variant"} {
			b507QuantityViewObject := b507QuantityViewInput(b507QuantityViewOperation, b507QuantityViewView, map[string]any{b507QuantityViewView: b507QuantityViewBytes})
			b507QuantityViewObject["context"].(map[string]any)[b507QuantityViewKey] = map[string]any{"opaque": []any{nil, 0, "other-family"}}
			if b507QuantityViewRun(b507QuantityViewT, b507QuantityViewObject)["state"] != "Hypothesis" {
				b507QuantityViewT.Fatal("family selected view")
			}
		}
	}
	for _, b507QuantityViewOperation := range []int{9, 11} {
		for _, b507QuantityViewWord := range [][]int{{0, 0}, {255, 255}, {0, 1}} {
			b507QuantityViewResult := b507QuantityViewRun(b507QuantityViewT, b507QuantityViewInput(b507QuantityViewOperation, "preference", map[string]any{"preference": b507QuantityViewWord}))
			if b507QuantityViewResult["state"] != "Unknown" {
				b507QuantityViewT.Fatal(b507QuantityViewResult)
			}
		}
	}
	for _, b507QuantityViewOperation := range []int{6, 8} {
		b507QuantityViewView := "preference"
		if b507QuantityViewOperation == 8 {
			b507QuantityViewView = "capacity"
		}
		if b507QuantityViewRun(b507QuantityViewT, b507QuantityViewInput(b507QuantityViewOperation, b507QuantityViewView, map[string]any{"capacity": []int{1}, "preference": []int{0, 0}}))["state"] != "Unknown" {
			b507QuantityViewT.Fatal("crossed view")
		}
	}
}

func TestAb05Ao05WidthsPresenceAndSiblings(b507QuantityViewT *b507QuantityViewTesting.T) {
	for _, b507QuantityViewView := range []string{"capacity", "preference"} {
		b507QuantityViewOperation := 6
		b507QuantityViewWidths := [][]int{nil, {}, {0, 1}, {0, 1, 2}}
		if b507QuantityViewView == "preference" {
			b507QuantityViewOperation = 8
			b507QuantityViewWidths = [][]int{nil, {}, {0}, {0, 1, 2}}
		}
		for b507QuantityViewIndex, b507QuantityViewBytes := range b507QuantityViewWidths {
			b507QuantityViewT.Run(b507QuantityViewFmt.Sprintf("%s/%d", b507QuantityViewView, b507QuantityViewIndex), func(b507QuantityViewT *b507QuantityViewTesting.T) {
				b507QuantityViewFields := map[string]any{"tail": []int{0, 255}}
				if b507QuantityViewIndex != 0 {
					b507QuantityViewFields[b507QuantityViewView] = b507QuantityViewBytes
				}
				b507QuantityViewResult := b507QuantityViewRun(b507QuantityViewT, b507QuantityViewInput(b507QuantityViewOperation, b507QuantityViewView, b507QuantityViewFields))
				b507QuantityViewField := b507QuantityViewResult["fields"].(map[string]any)[b507QuantityViewView].(map[string]any)
				_, b507QuantityViewPresent := b507QuantityViewField["raw"]
				if b507QuantityViewPresent != (b507QuantityViewIndex != 0) || b507QuantityViewField["state"] != "incomplete" || b507QuantityViewField["number"] != nil || b507QuantityViewResult["state"] != "partial" {
					b507QuantityViewT.Fatal(b507QuantityViewResult)
				}
			})
		}
	}
	for _, b507QuantityViewKey := range []string{"preference", "cost", "tail"} {
		for _, b507QuantityViewBytes := range [][]int{{}, {0}, {255, 0, 1}} {
			b507QuantityViewObject := b507QuantityViewInput(6, "capacity", map[string]any{"capacity": []int{1}, b507QuantityViewKey: b507QuantityViewBytes})
			if b507QuantityViewRun(b507QuantityViewT, b507QuantityViewObject)["fields"].(map[string]any)["capacity"].(map[string]any)["number"] != 0.1 {
				b507QuantityViewT.Fatal("sibling erased capacity")
			}
		}
	}
}

func TestAb05Ao05RejectInvalid(b507QuantityViewT *b507QuantityViewTesting.T) {
	for b507QuantityViewIndex, b507QuantityViewInput := range []string{
		``, `{`, `null`, `[]`, `{} {}`, `{"context":{},"fields":{},"context":{}}`,
		`{"context":{"family":{"x":1,"x":2}},"fields":{}}`,
		`{"context":{},"fields":{"x":[],"x":[]}}`,
		`{"context":{},"fields":{},"arguments":{"x":1,"x":2}}`,
		`{"context":null,"fields":{}}`, `{"context":{},"fields":null}`,
		`{"context":{},"fields":{"x":null}}`, `{"context":{},"fields":{"x":"AA=="}}`,
		`{"context":{},"fields":{"x":[-1]}}`, `{"context":{},"fields":{"x":[256]}}`,
		`{"context":{},"fields":{"x":[1.0]}}`, `{"context":{},"fields":{"x":[1e0]}}`,
		`{"context":{},"fields":{"x":[true]}}`, `{"context":{},"fields":{"x":["0"]}}`,
		`{"context":{},"fields":{"x":[null]}}`, `{"context":{},"fields":{"x":[{}]}}`,
		`{"context":{},"fields":{"x":[[]]}}`, `{"context":{},"fields":{},"arguments":null}`,
		`{"context":{},"fields":{},"arguments":{"execute":true}}`,
		`{"context":{},"fields":{},"extra":0}`, `{"context":{},"fields":{}} true`,
	} {
		b507QuantityViewT.Run(b507QuantityViewFmt.Sprint(b507QuantityViewIndex), func(b507QuantityViewT *b507QuantityViewTesting.T) {
			b507QuantityViewOutput, b507QuantityViewError := B507QuantityViewInterpret([]byte(b507QuantityViewInput))
			if b507QuantityViewError == nil || b507QuantityViewOutput != nil {
				b507QuantityViewT.Fatalf("accepted %s: %s, %v", b507QuantityViewInput, b507QuantityViewOutput, b507QuantityViewError)
			}
		})
	}
	for _, b507QuantityViewKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, b507QuantityViewValue := range []any{nil, true, []int{}, map[string]any{}, 1.5} {
			b507QuantityViewObject := b507QuantityViewInput(6, "capacity", map[string]any{"capacity": []int{1}})
			b507QuantityViewObject["context"].(map[string]any)[b507QuantityViewKey] = b507QuantityViewValue
			b507QuantityViewBytes, _ := b507QuantityViewJSON.Marshal(b507QuantityViewObject)
			if b507QuantityViewOutput, b507QuantityViewError := B507QuantityViewInterpret(b507QuantityViewBytes); b507QuantityViewError == nil || b507QuantityViewOutput != nil {
				b507QuantityViewT.Fatal("invalid context accepted", b507QuantityViewKey, b507QuantityViewValue)
			}
		}
	}
}

func TestAb05Ao05ContextKinds(b507QuantityViewT *b507QuantityViewTesting.T) {
	for _, b507QuantityViewKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		b507QuantityViewObject := b507QuantityViewInput(6, "capacity", map[string]any{"capacity": []int{1}})
		if b507QuantityViewKey == "direction" || b507QuantityViewKey == "view" {
			b507QuantityViewObject["context"].(map[string]any)[b507QuantityViewKey] = 0
		} else {
			b507QuantityViewObject["context"].(map[string]any)[b507QuantityViewKey] = "6"
		}
		b507QuantityViewBytes, _ := b507QuantityViewJSON.Marshal(b507QuantityViewObject)
		if b507QuantityViewOutput, b507QuantityViewError := B507QuantityViewInterpret(b507QuantityViewBytes); b507QuantityViewError == nil || b507QuantityViewOutput != nil {
			b507QuantityViewT.Fatal("wrong context kind", b507QuantityViewKey)
		}
	}
}

func TestAb05Ao05OwnershipAndArguments(b507QuantityViewT *b507QuantityViewTesting.T) {
	b507QuantityViewObject := b507QuantityViewInput(6, "capacity", map[string]any{"capacity": []int{1}})
	b507QuantityViewAbsent := b507QuantityViewRun(b507QuantityViewT, b507QuantityViewObject)
	if _, b507QuantityViewPresent := b507QuantityViewAbsent["arguments"]; b507QuantityViewPresent {
		b507QuantityViewT.Fatal("invented arguments")
	}
	b507QuantityViewObject["arguments"] = map[string]any{}
	b507QuantityViewResult := b507QuantityViewRun(b507QuantityViewT, b507QuantityViewObject)
	if !b507QuantityViewReflect.DeepEqual(b507QuantityViewResult["arguments"], map[string]any{}) {
		b507QuantityViewT.Fatal("lost empty arguments")
	}
	b507QuantityViewInputBytes, _ := b507QuantityViewJSON.Marshal(b507QuantityViewObject)
	b507QuantityViewOutput, b507QuantityViewError := B507QuantityViewInterpret(b507QuantityViewInputBytes)
	if b507QuantityViewError != nil {
		b507QuantityViewT.Fatal(b507QuantityViewError)
	}
	b507QuantityViewSaved := append([]byte(nil), b507QuantityViewOutput...)
	for b507QuantityViewIndex := range b507QuantityViewInputBytes {
		b507QuantityViewInputBytes[b507QuantityViewIndex] = ' '
	}
	if !b507QuantityViewBytesPackage.Equal(b507QuantityViewSaved, b507QuantityViewOutput) {
		b507QuantityViewT.Fatal("input alias")
	}
	b507QuantityViewResult["raw"].(map[string]any)["capacity"].([]any)[0] = 255.0
	b507QuantityViewResult["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "mutated"
	b507QuantityViewResult["arguments"].(map[string]any)["execute"] = true
	if b507QuantityViewResult["fields"].(map[string]any)["capacity"].(map[string]any)["raw"].([]any)[0] != 1.0 {
		b507QuantityViewT.Fatal("nested alias")
	}
	b507QuantityViewResult["fields"].(map[string]any)["capacity"].(map[string]any)["raw"].([]any)[0] = 7.0
	b507QuantityViewAgain := b507QuantityViewRun(b507QuantityViewT, b507QuantityViewObject)
	if b507QuantityViewAgain["fields"].(map[string]any)["capacity"].(map[string]any)["number"] != 0.1 {
		b507QuantityViewT.Fatal("cross-call alias")
	}
	for b507QuantityViewIndex := range b507QuantityViewOutput {
		b507QuantityViewOutput[b507QuantityViewIndex] = ' '
	}
	b507QuantityViewFreshBytes, _ := b507QuantityViewJSON.Marshal(b507QuantityViewObject)
	b507QuantityViewFresh, b507QuantityViewError := B507QuantityViewInterpret(b507QuantityViewFreshBytes)
	if b507QuantityViewError != nil || !b507QuantityViewBytesPackage.Equal(b507QuantityViewFresh, b507QuantityViewSaved) {
		b507QuantityViewT.Fatal("output alias")
	}
}

func TestAb05Ao05SyntheticFixture(b507QuantityViewT *b507QuantityViewTesting.T) {
	b507QuantityViewBytes, b507QuantityViewError := b507QuantityViewOs.ReadFile("testdata/ab05ao05_b507_quantity_views.json")
	if b507QuantityViewError != nil {
		b507QuantityViewT.Fatal(b507QuantityViewError)
	}
	b507QuantityViewCases := map[string]map[string]any{}
	if b507QuantityViewError := b507QuantityViewJSON.Unmarshal(b507QuantityViewBytes, &b507QuantityViewCases); b507QuantityViewError != nil {
		b507QuantityViewT.Fatal(b507QuantityViewError)
	}
	if len(b507QuantityViewCases) != 8 {
		b507QuantityViewT.Fatal("fixture count")
	}
	for b507QuantityViewID, b507QuantityViewObject := range b507QuantityViewCases {
		b507QuantityViewT.Run(b507QuantityViewID, func(b507QuantityViewT *b507QuantityViewTesting.T) {
			b507QuantityViewRun(b507QuantityViewT, b507QuantityViewObject)
		})
	}
}
