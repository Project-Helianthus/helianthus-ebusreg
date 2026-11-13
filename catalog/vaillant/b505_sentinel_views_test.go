package vaillant

import (
	b505SentinelViewBytes "bytes"
	b505SentinelViewJSON "encoding/json"
	b505SentinelViewOS "os"
	b505SentinelViewReflect "reflect"
	b505SentinelViewTesting "testing"
)

func b505SentinelViewControl() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 5, "operation": 50, "direction": "request", "view": "replacement-pair", "family": "synthetic-family", "variant": map[string]any{"opaque": []any{0, "variant"}}}, "fields": map[string]any{"temperature_one": []int{0, 128}, "temperature_two": []int{0, 0}, "unmodelled": []int{11, 255, 0}}}
}

func b505SentinelViewCall(b505SentinelViewT *b505SentinelViewTesting.T, b505SentinelViewInput any) map[string]any {
	b505SentinelViewT.Helper()
	b505SentinelViewEncoded, b505SentinelViewError := b505SentinelViewJSON.Marshal(b505SentinelViewInput)
	if b505SentinelViewError != nil {
		b505SentinelViewT.Fatal(b505SentinelViewError)
	}
	b505SentinelViewOutput, b505SentinelViewError := B505SentinelViewInterpret(b505SentinelViewEncoded)
	if b505SentinelViewError != nil {
		b505SentinelViewT.Fatal(b505SentinelViewError)
	}
	var b505SentinelViewResult map[string]any
	if b505SentinelViewError := b505SentinelViewJSON.Unmarshal(b505SentinelViewOutput, &b505SentinelViewResult); b505SentinelViewError != nil {
		b505SentinelViewT.Fatal(b505SentinelViewError)
	}
	var b505SentinelViewOriginal map[string]any
	if b505SentinelViewError := b505SentinelViewJSON.Unmarshal(b505SentinelViewEncoded, &b505SentinelViewOriginal); b505SentinelViewError != nil {
		b505SentinelViewT.Fatal(b505SentinelViewError)
	}
	if !b505SentinelViewReflect.DeepEqual(b505SentinelViewOriginal["context"], b505SentinelViewResult["context"]) || !b505SentinelViewReflect.DeepEqual(b505SentinelViewOriginal["fields"], b505SentinelViewResult["raw"]) {
		b505SentinelViewT.Fatal("context or raw evidence changed")
	}
	if b505SentinelViewResult["physical"] != "Unknown" || b505SentinelViewResult["qualified"] != false {
		b505SentinelViewT.Fatal("native qualification inferred")
	}
	b505SentinelViewArgument, b505SentinelViewPresent := b505SentinelViewOriginal["arguments"]
	b505SentinelViewReturned, b505SentinelViewReturnedPresent := b505SentinelViewResult["arguments"]
	if b505SentinelViewPresent != b505SentinelViewReturnedPresent || !b505SentinelViewReflect.DeepEqual(b505SentinelViewArgument, b505SentinelViewReturned) {
		b505SentinelViewT.Fatal("arguments presence changed")
	}
	for _, b505SentinelViewValue := range b505SentinelViewResult["fields"].(map[string]any) {
		b505SentinelViewField := b505SentinelViewValue.(map[string]any)
		for _, b505SentinelViewKey := range []string{"number", "calculation", "candidate_domain", "flags", "reserved"} {
			if _, b505SentinelViewPresent := b505SentinelViewField[b505SentinelViewKey]; b505SentinelViewPresent {
				b505SentinelViewT.Fatalf("invented %s", b505SentinelViewKey)
			}
		}
	}
	return b505SentinelViewResult
}

func TestAb01Ao01IndependentWords(b505SentinelViewT *b505SentinelViewTesting.T) {
	b505SentinelViewData, b505SentinelViewError := b505SentinelViewOS.ReadFile("testdata/ab01ao01_b505_sentinel_views.json")
	if b505SentinelViewError != nil {
		b505SentinelViewT.Fatal(b505SentinelViewError)
	}
	var b505SentinelViewCases []struct {
		B505SentinelViewID    string         `json:"id"`
		B505SentinelViewInput map[string]any `json:"input"`
		B505SentinelViewField string         `json:"field"`
		B505SentinelViewLabel any            `json:"label"`
		B505SentinelViewState string         `json:"state"`
	}
	if b505SentinelViewError := b505SentinelViewJSON.Unmarshal(b505SentinelViewData, &b505SentinelViewCases); b505SentinelViewError != nil {
		b505SentinelViewT.Fatal(b505SentinelViewError)
	}
	if len(b505SentinelViewCases) != 16 {
		b505SentinelViewT.Fatal("fixture cases missing")
	}
	for _, b505SentinelViewCase := range b505SentinelViewCases {
		b505SentinelViewT.Run(b505SentinelViewCase.B505SentinelViewID, func(b505SentinelViewT *b505SentinelViewTesting.T) {
			b505SentinelViewResult := b505SentinelViewCall(b505SentinelViewT, b505SentinelViewCase.B505SentinelViewInput)
			b505SentinelViewFields := b505SentinelViewResult["fields"].(map[string]any)
			if len(b505SentinelViewFields) != 2 {
				b505SentinelViewT.Fatal("independent result fields lost")
			}
			b505SentinelViewField := b505SentinelViewFields[b505SentinelViewCase.B505SentinelViewField].(map[string]any)
			if b505SentinelViewField["state"] != b505SentinelViewCase.B505SentinelViewState || b505SentinelViewField["label"] != b505SentinelViewCase.B505SentinelViewLabel {
				b505SentinelViewT.Fatalf("wrong branch: %v", b505SentinelViewField)
			}
			b505SentinelViewRaw := b505SentinelViewCase.B505SentinelViewInput["fields"].(map[string]any)[b505SentinelViewCase.B505SentinelViewField]
			if !b505SentinelViewReflect.DeepEqual(b505SentinelViewRaw, b505SentinelViewField["raw"]) {
				b505SentinelViewT.Fatal("field raw changed")
			}
			b505SentinelViewSibling := "temperature_two"
			b505SentinelViewLabel := "circuit-default"
			if b505SentinelViewCase.B505SentinelViewField == "temperature_two" {
				b505SentinelViewSibling = "temperature_one"
				b505SentinelViewLabel = "preserve-existing"
			}
			b505SentinelViewOther := b505SentinelViewFields[b505SentinelViewSibling].(map[string]any)
			if b505SentinelViewOther["state"] != "Hypothesis" || b505SentinelViewOther["label"] != b505SentinelViewLabel {
				b505SentinelViewT.Fatal("independent valid sibling lost")
			}
			b505SentinelViewExpectedState := "partial"
			if b505SentinelViewCase.B505SentinelViewState == "Hypothesis" {
				b505SentinelViewExpectedState = "Hypothesis"
			}
			if b505SentinelViewResult["state"] != b505SentinelViewExpectedState {
				b505SentinelViewT.Fatal("wrong aggregate state", b505SentinelViewResult["state"])
			}
			if b505SentinelViewCase.B505SentinelViewLabel != nil && b505SentinelViewField["grade"] != "Hypothesis" {
				b505SentinelViewT.Fatal("semantic grade changed")
			}
		})
	}
	b505SentinelViewInput := b505SentinelViewControl()
	b505SentinelViewInput["fields"].(map[string]any)["temperature_one"] = []int{0, 0}
	b505SentinelViewInput["fields"].(map[string]any)["temperature_two"] = []int{0, 128}
	b505SentinelViewResult := b505SentinelViewCall(b505SentinelViewT, b505SentinelViewInput)["fields"].(map[string]any)
	if b505SentinelViewResult["temperature_one"].(map[string]any)["label"] != "circuit-default" || b505SentinelViewResult["temperature_two"].(map[string]any)["label"] != "preserve-existing" {
		b505SentinelViewT.Fatal("reverse pair failed")
	}
}

func TestAb01Ao01SelectionAxes(b505SentinelViewT *b505SentinelViewTesting.T) {
	for b505SentinelViewKey, b505SentinelViewValue := range map[string]any{"primary": 180, "secondary": 6, "operation": 49, "direction": "response", "view": "other"} {
		b505SentinelViewT.Run(b505SentinelViewKey, func(b505SentinelViewT *b505SentinelViewTesting.T) {
			for _, b505SentinelViewMode := range []string{"different", "missing", "empty-or-zero"} {
				b505SentinelViewInput := b505SentinelViewControl()
				b505SentinelViewContext := b505SentinelViewInput["context"].(map[string]any)
				switch b505SentinelViewMode {
				case "different":
					b505SentinelViewContext[b505SentinelViewKey] = b505SentinelViewValue
				case "missing":
					delete(b505SentinelViewContext, b505SentinelViewKey)
				default:
					if _, b505SentinelViewOK := b505SentinelViewValue.(string); b505SentinelViewOK {
						b505SentinelViewContext[b505SentinelViewKey] = ""
					} else {
						b505SentinelViewContext[b505SentinelViewKey] = 0
					}
				}
				b505SentinelViewResult := b505SentinelViewCall(b505SentinelViewT, b505SentinelViewInput)
				if b505SentinelViewResult["state"] != "Unknown" || len(b505SentinelViewResult["fields"].(map[string]any)) != 0 {
					b505SentinelViewT.Fatal("selection guessed", b505SentinelViewMode)
				}
			}
		})
	}
	for _, b505SentinelViewKey := range []string{"family", "variant"} {
		for _, b505SentinelViewValue := range []any{nil, "other-family", 0, []any{}, map[string]any{"opaque": true}} {
			b505SentinelViewInput := b505SentinelViewControl()
			b505SentinelViewInput["context"].(map[string]any)[b505SentinelViewKey] = b505SentinelViewValue
			if b505SentinelViewCall(b505SentinelViewT, b505SentinelViewInput)["state"] != "Hypothesis" {
				b505SentinelViewT.Fatal("opaque context selected view")
			}
		}
		b505SentinelViewInput := b505SentinelViewControl()
		delete(b505SentinelViewInput["context"].(map[string]any), b505SentinelViewKey)
		if b505SentinelViewCall(b505SentinelViewT, b505SentinelViewInput)["state"] != "Hypothesis" {
			b505SentinelViewT.Fatal("opaque context required")
		}
	}
	b505SentinelViewInput := b505SentinelViewControl()
	b505SentinelViewInput["arguments"] = map[string]any{}
	b505SentinelViewCall(b505SentinelViewT, b505SentinelViewInput)
}

func TestAb01Ao01UnknownAndPresence(b505SentinelViewT *b505SentinelViewTesting.T) {
	for _, b505SentinelViewWords := range []map[string]any{{}, {"temperature_one": []int{}, "temperature_two": []int{}}, {"temperature_one": []int{1, 0}, "temperature_two": []int{255, 255}}} {
		b505SentinelViewInput := b505SentinelViewControl()
		b505SentinelViewInput["fields"] = b505SentinelViewWords
		b505SentinelViewResult := b505SentinelViewCall(b505SentinelViewT, b505SentinelViewInput)
		b505SentinelViewExpected := "partial"
		if len(b505SentinelViewWords) == 2 && len(b505SentinelViewWords["temperature_one"].([]int)) == 2 {
			b505SentinelViewExpected = "Unknown"
		}
		if b505SentinelViewResult["state"] != b505SentinelViewExpected {
			b505SentinelViewT.Fatal("unknown or incomplete state lost")
		}
	}
	for _, b505SentinelViewValue := range [][]int{{128, 0}, {0, 1}, {255, 255}, {0, 0, 0}, {}} {
		b505SentinelViewInput := b505SentinelViewControl()
		b505SentinelViewInput["fields"].(map[string]any)["unmodelled"] = b505SentinelViewValue
		if b505SentinelViewCall(b505SentinelViewT, b505SentinelViewInput)["state"] != "Hypothesis" {
			b505SentinelViewT.Fatal("unmodelled field selected view")
		}
	}
}

func TestAb01Ao01RejectInvalid(b505SentinelViewT *b505SentinelViewTesting.T) {
	for _, b505SentinelViewInput := range []string{``, `{`, `[]`, `null`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"family":{"x":0,"x":1}},"fields":{}}`, `{"context":{},"fields":{"x":[],"x":[]}}`, `{"context":{},"fields":{},"arguments":{"x":0,"x":1}}`, `{"context":{},"fields":{},"extra":0}`} {
		b505SentinelViewOutput, b505SentinelViewError := B505SentinelViewInterpret([]byte(b505SentinelViewInput))
		if b505SentinelViewError == nil || b505SentinelViewOutput != nil {
			b505SentinelViewT.Fatalf("accepted invalid JSON: %s", b505SentinelViewInput)
		}
	}
	for _, b505SentinelViewKey := range []string{"context", "fields", "arguments"} {
		for _, b505SentinelViewValue := range []any{nil, 0, false, "", []any{}} {
			b505SentinelViewInput := b505SentinelViewControl()
			b505SentinelViewInput[b505SentinelViewKey] = b505SentinelViewValue
			b505SentinelViewEncoded, _ := b505SentinelViewJSON.Marshal(b505SentinelViewInput)
			b505SentinelViewOutput, b505SentinelViewError := B505SentinelViewInterpret(b505SentinelViewEncoded)
			if b505SentinelViewError == nil || b505SentinelViewOutput != nil {
				b505SentinelViewT.Fatal("accepted wrong object kind", b505SentinelViewKey, b505SentinelViewValue)
			}
		}
	}
	for _, b505SentinelViewKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, b505SentinelViewValue := range []any{nil, true, []any{}, map[string]any{}} {
			b505SentinelViewInput := b505SentinelViewControl()
			b505SentinelViewInput["context"].(map[string]any)[b505SentinelViewKey] = b505SentinelViewValue
			b505SentinelViewEncoded, _ := b505SentinelViewJSON.Marshal(b505SentinelViewInput)
			b505SentinelViewOutput, b505SentinelViewError := B505SentinelViewInterpret(b505SentinelViewEncoded)
			if b505SentinelViewError == nil || b505SentinelViewOutput != nil {
				b505SentinelViewT.Fatal("accepted invalid selector", b505SentinelViewKey)
			}
		}
	}
	for _, b505SentinelViewKey := range []string{"temperature_one", "temperature_two", "unmodelled"} {
		for _, b505SentinelViewValue := range []any{nil, 0, "AA==", map[string]any{}, []any{-1}, []any{256}, []any{1.5}, []any{true}, []any{nil}, []any{"0"}, []any{[]any{}}, []any{b505SentinelViewJSON.Number("1e0")}} {
			b505SentinelViewInput := b505SentinelViewControl()
			b505SentinelViewInput["fields"].(map[string]any)[b505SentinelViewKey] = b505SentinelViewValue
			b505SentinelViewEncoded, _ := b505SentinelViewJSON.Marshal(b505SentinelViewInput)
			b505SentinelViewOutput, b505SentinelViewError := B505SentinelViewInterpret(b505SentinelViewEncoded)
			if b505SentinelViewError == nil || b505SentinelViewOutput != nil {
				b505SentinelViewT.Fatal("accepted invalid byte array", b505SentinelViewKey, b505SentinelViewValue)
			}
		}
	}
	b505SentinelViewInput := b505SentinelViewControl()
	b505SentinelViewInput["arguments"] = map[string]any{"flag": true}
	b505SentinelViewEncoded, _ := b505SentinelViewJSON.Marshal(b505SentinelViewInput)
	if b505SentinelViewOutput, b505SentinelViewError := B505SentinelViewInterpret(b505SentinelViewEncoded); b505SentinelViewError == nil || b505SentinelViewOutput != nil {
		b505SentinelViewT.Fatal("accepted undefined argument")
	}
}

func TestAb01Ao01DetachedMemory(b505SentinelViewT *b505SentinelViewTesting.T) {
	b505SentinelViewInput, _ := b505SentinelViewJSON.Marshal(b505SentinelViewControl())
	b505SentinelViewFirst, b505SentinelViewError := B505SentinelViewInterpret(b505SentinelViewInput)
	if b505SentinelViewError != nil {
		b505SentinelViewT.Fatal(b505SentinelViewError)
	}
	b505SentinelViewSaved := append([]byte(nil), b505SentinelViewFirst...)
	for b505SentinelViewIndex := range b505SentinelViewInput {
		b505SentinelViewInput[b505SentinelViewIndex] = ' '
	}
	if !b505SentinelViewBytes.Equal(b505SentinelViewFirst, b505SentinelViewSaved) {
		b505SentinelViewT.Fatal("input aliases output")
	}
	var b505SentinelViewNested map[string]any
	if b505SentinelViewError := b505SentinelViewJSON.Unmarshal(b505SentinelViewFirst, &b505SentinelViewNested); b505SentinelViewError != nil {
		b505SentinelViewT.Fatal(b505SentinelViewError)
	}
	b505SentinelViewNested["raw"].(map[string]any)["temperature_one"].([]any)[0] = float64(99)
	if b505SentinelViewNested["fields"].(map[string]any)["temperature_one"].(map[string]any)["raw"].([]any)[0] != float64(0) {
		b505SentinelViewT.Fatal("raw aliases semantic field")
	}
	b505SentinelViewNested["context"].(map[string]any)["variant"].(map[string]any)["opaque"].([]any)[0] = "changed"
	for b505SentinelViewIndex := range b505SentinelViewFirst {
		b505SentinelViewFirst[b505SentinelViewIndex] = ' '
	}
	b505SentinelViewFresh, _ := b505SentinelViewJSON.Marshal(b505SentinelViewControl())
	b505SentinelViewSecond, b505SentinelViewError := B505SentinelViewInterpret(b505SentinelViewFresh)
	if b505SentinelViewError != nil || !b505SentinelViewBytes.Equal(b505SentinelViewSecond, b505SentinelViewSaved) {
		b505SentinelViewT.Fatal("later call changed")
	}
	if !b505SentinelViewJSON.Valid(b505SentinelViewFresh) {
		b505SentinelViewT.Fatal("output mutation changed later input")
	}
}

func TestAb01Ao01IntegerAndKinds(b505SentinelViewT *b505SentinelViewTesting.T) {
	for _, b505SentinelViewKey := range []string{"primary", "secondary", "operation"} {
		b505SentinelViewInput := b505SentinelViewControl()
		b505SentinelViewInput["context"].(map[string]any)[b505SentinelViewKey] = b505SentinelViewJSON.Number("99999999999999999999999999999999999")
		if b505SentinelViewCall(b505SentinelViewT, b505SentinelViewInput)["state"] != "Unknown" {
			b505SentinelViewT.Fatal("unknown integer selected view")
		}
		for _, b505SentinelViewValue := range []any{"50", 1.5, b505SentinelViewJSON.Number("1e0")} {
			b505SentinelViewInput["context"].(map[string]any)[b505SentinelViewKey] = b505SentinelViewValue
			b505SentinelViewEncoded, _ := b505SentinelViewJSON.Marshal(b505SentinelViewInput)
			if b505SentinelViewOutput, b505SentinelViewError := B505SentinelViewInterpret(b505SentinelViewEncoded); b505SentinelViewError == nil || b505SentinelViewOutput != nil {
				b505SentinelViewT.Fatal("invalid integer accepted", b505SentinelViewKey)
			}
		}
	}
	for _, b505SentinelViewKey := range []string{"direction", "view"} {
		b505SentinelViewInput := b505SentinelViewControl()
		b505SentinelViewInput["context"].(map[string]any)[b505SentinelViewKey] = 0
		b505SentinelViewEncoded, _ := b505SentinelViewJSON.Marshal(b505SentinelViewInput)
		if b505SentinelViewOutput, b505SentinelViewError := B505SentinelViewInterpret(b505SentinelViewEncoded); b505SentinelViewError == nil || b505SentinelViewOutput != nil {
			b505SentinelViewT.Fatal("numeric string selector accepted")
		}
	}
	for _, b505SentinelViewKey := range []string{"context", "fields"} {
		b505SentinelViewInput := b505SentinelViewControl()
		delete(b505SentinelViewInput, b505SentinelViewKey)
		b505SentinelViewEncoded, _ := b505SentinelViewJSON.Marshal(b505SentinelViewInput)
		if b505SentinelViewOutput, b505SentinelViewError := B505SentinelViewInterpret(b505SentinelViewEncoded); b505SentinelViewError == nil || b505SentinelViewOutput != nil {
			b505SentinelViewT.Fatal("missing required object accepted")
		}
	}
	b505SentinelViewInput := b505SentinelViewControl()
	b505SentinelViewInput["fields"].(map[string]any)["temperature_one"] = []any{b505SentinelViewJSON.Number("-0"), 128}
	if b505SentinelViewCall(b505SentinelViewT, b505SentinelViewInput)["state"] != "Hypothesis" {
		b505SentinelViewT.Fatal("integer negative zero misinterpreted")
	}
}
