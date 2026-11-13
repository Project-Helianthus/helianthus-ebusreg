package vaillant

import (
	b503ControlFieldJSON "encoding/json"
	b503ControlFieldFmt "fmt"
	b503ControlFieldOs "os"
	b503ControlFieldReflect "reflect"
	b503ControlFieldStrings "strings"
	b503ControlFieldTesting "testing"
)

func b503ControlFieldEvaluateTest(b503ControlFieldT *b503ControlFieldTesting.T, b503ControlFieldInput any) map[string]any {
	b503ControlFieldT.Helper()
	b503ControlFieldBytes, b503ControlFieldError := b503ControlFieldJSON.Marshal(b503ControlFieldInput)
	if b503ControlFieldError != nil {
		b503ControlFieldT.Fatal(b503ControlFieldError)
	}
	b503ControlFieldBytes, b503ControlFieldError = B503ControlFieldEvaluate(b503ControlFieldBytes)
	if b503ControlFieldError != nil {
		b503ControlFieldT.Fatal(b503ControlFieldError)
	}
	b503ControlFieldOutput := map[string]any{}
	if b503ControlFieldError = b503ControlFieldJSON.Unmarshal(b503ControlFieldBytes, &b503ControlFieldOutput); b503ControlFieldError != nil {
		b503ControlFieldT.Fatal(b503ControlFieldError)
	}
	return b503ControlFieldOutput
}

func b503ControlFieldSubset(b503ControlFieldT *b503ControlFieldTesting.T, b503ControlFieldWant, b503ControlFieldGot any) {
	b503ControlFieldT.Helper()
	if b503ControlFieldMap, b503ControlFieldOK := b503ControlFieldWant.(map[string]any); b503ControlFieldOK {
		b503ControlFieldActual, b503ControlFieldOK := b503ControlFieldGot.(map[string]any)
		if !b503ControlFieldOK {
			b503ControlFieldT.Fatalf("want object %v got %v", b503ControlFieldWant, b503ControlFieldGot)
		}
		for b503ControlFieldKey, b503ControlFieldValue := range b503ControlFieldMap {
			b503ControlFieldSubset(b503ControlFieldT, b503ControlFieldValue, b503ControlFieldActual[b503ControlFieldKey])
		}
	} else if !b503ControlFieldReflect.DeepEqual(b503ControlFieldWant, b503ControlFieldGot) {
		b503ControlFieldT.Fatalf("want %v got %v", b503ControlFieldWant, b503ControlFieldGot)
	}
}

func b503ControlFieldFixtures(b503ControlFieldT *b503ControlFieldTesting.T) []map[string]any {
	b503ControlFieldT.Helper()
	b503ControlFieldBytes, b503ControlFieldError := b503ControlFieldOs.ReadFile("testdata/ab07aq07_b503_control_fields.json")
	if b503ControlFieldError != nil {
		b503ControlFieldT.Fatal(b503ControlFieldError)
	}
	var b503ControlFieldCases []map[string]any
	if b503ControlFieldError = b503ControlFieldJSON.Unmarshal(b503ControlFieldBytes, &b503ControlFieldCases); b503ControlFieldError != nil {
		b503ControlFieldT.Fatal(b503ControlFieldError)
	}
	if len(b503ControlFieldCases) != 22 {
		b503ControlFieldT.Fatalf("fixture count %d", len(b503ControlFieldCases))
	}
	return b503ControlFieldCases
}

func TestAb07Aq07Fixtures(b503ControlFieldT *b503ControlFieldTesting.T) {
	for _, b503ControlFieldCase := range b503ControlFieldFixtures(b503ControlFieldT) {
		b503ControlFieldT.Run(b503ControlFieldCase["ab07aq07id"].(string), func(b503ControlFieldT *b503ControlFieldTesting.T) {
			b503ControlFieldInput := b503ControlFieldCase["input"].(map[string]any)
			b503ControlFieldOutput := b503ControlFieldEvaluateTest(b503ControlFieldT, b503ControlFieldInput)
			b503ControlFieldSubset(b503ControlFieldT, b503ControlFieldCase["expected"], b503ControlFieldOutput)
			for b503ControlFieldKey, b503ControlFieldValue := range b503ControlFieldInput {
				if !b503ControlFieldReflect.DeepEqual(b503ControlFieldValue, b503ControlFieldOutput[b503ControlFieldKey]) {
					b503ControlFieldT.Fatalf("lost input %s", b503ControlFieldKey)
				}
			}
			if len(b503ControlFieldOutput) != len(b503ControlFieldInput)+4 {
				b503ControlFieldT.Fatal("invented output claim")
			}
			b503ControlFieldResults := b503ControlFieldOutput["results"].(map[string]any)
			if b503ControlFieldValue, b503ControlFieldPresent := b503ControlFieldResults["flags"]; b503ControlFieldPresent {
				b503ControlFieldResult := b503ControlFieldValue.(map[string]any)
				if b503ControlFieldFlags, b503ControlFieldOK := b503ControlFieldResult["flags"].(map[string]any); b503ControlFieldOK && len(b503ControlFieldFlags) != 5 {
					b503ControlFieldT.Fatal("reserved bits labeled")
				}
				if _, b503ControlFieldPresent := b503ControlFieldResult["label"]; b503ControlFieldPresent {
					b503ControlFieldT.Fatal("flags given semantic label")
				}
			}
			for _, b503ControlFieldValue := range b503ControlFieldResults {
				b503ControlFieldResult := b503ControlFieldValue.(map[string]any)
				for _, b503ControlFieldKey := range []string{"reset_success", "effects", "freshness", "qualified", "executed", "delivery"} {
					if _, b503ControlFieldPresent := b503ControlFieldResult[b503ControlFieldKey]; b503ControlFieldPresent {
						b503ControlFieldT.Fatalf("invented %s", b503ControlFieldKey)
					}
				}
				if b503ControlFieldResult["state"] == "Unknown" {
					if _, b503ControlFieldPresent := b503ControlFieldResult["label"]; b503ControlFieldPresent {
						b503ControlFieldT.Fatal("invented label")
					}
				}
			}
		})
	}
}

func TestAb07Aq07Context(b503ControlFieldT *b503ControlFieldTesting.T) {
	for _, b503ControlFieldCase := range b503ControlFieldFixtures(b503ControlFieldT) {
		for b503ControlFieldKey, b503ControlFieldValues := range map[string][]any{"primary": {float64(180), nil}, "secondary": {float64(4), nil}, "operation": {float64(0), float64(1), nil}, "direction": {"other", nil}, "view": {"other", nil}, "family": {"other", nil, false, float64(0)}, "variant": {"other", nil, false, float64(0)}} {
			for _, b503ControlFieldValue := range b503ControlFieldValues {
				b503ControlFieldT.Run(b503ControlFieldFmt.Sprint(b503ControlFieldCase["ab07aq07id"], b503ControlFieldKey, b503ControlFieldValue), func(b503ControlFieldT *b503ControlFieldTesting.T) {
					b503ControlFieldBytes, _ := b503ControlFieldJSON.Marshal(b503ControlFieldCase["input"])
					b503ControlFieldInput := map[string]any{}
					_ = b503ControlFieldJSON.Unmarshal(b503ControlFieldBytes, &b503ControlFieldInput)
					b503ControlFieldContext := b503ControlFieldInput["context"].(map[string]any)
					if b503ControlFieldValue == nil {
						delete(b503ControlFieldContext, b503ControlFieldKey)
					} else {
						b503ControlFieldContext[b503ControlFieldKey] = b503ControlFieldValue
					}
					b503ControlFieldOutput := b503ControlFieldEvaluateTest(b503ControlFieldT, b503ControlFieldInput)
					if !b503ControlFieldReflect.DeepEqual(b503ControlFieldContext, b503ControlFieldOutput["context"]) || b503ControlFieldOutput["qualified"] != false {
						b503ControlFieldT.Fatal("identity lost or qualified")
					}
					if b503ControlFieldKey == "family" || b503ControlFieldKey == "variant" {
						b503ControlFieldSubset(b503ControlFieldT, b503ControlFieldCase["expected"], b503ControlFieldOutput)
					} else if b503ControlFieldOutput["state"] != "Unknown" || len(b503ControlFieldOutput["results"].(map[string]any)) != 0 {
						b503ControlFieldT.Fatal("unsupported context interpreted")
					}
				})
			}
		}
	}
	for _, b503ControlFieldView := range []string{"reset-acceptance", "category-reset"} {
		b503ControlFieldDirection := "positive"
		b503ControlFieldField := "acceptance"
		if b503ControlFieldView == "category-reset" {
			b503ControlFieldDirection = "request"
			b503ControlFieldField = "category"
		}
		b503ControlFieldOutput := b503ControlFieldEvaluateTest(b503ControlFieldT, map[string]any{"context": map[string]any{"primary": 181, "secondary": 3, "operation": 6, "direction": b503ControlFieldDirection, "view": b503ControlFieldView}, "fields": map[string]any{b503ControlFieldField: []int{1}}})
		if b503ControlFieldOutput["state"] != "complete" {
			b503ControlFieldT.Fatal("OP06 rejected")
		}
	}
}

func TestAb07Aq07WidthsAndFields(b503ControlFieldT *b503ControlFieldTesting.T) {
	for _, b503ControlFieldCase := range b503ControlFieldFixtures(b503ControlFieldT) {
		b503ControlFieldOriginal := b503ControlFieldCase["input"].(map[string]any)
		for b503ControlFieldKey := range b503ControlFieldOriginal["fields"].(map[string]any) {
			for _, b503ControlFieldReplacement := range []any{nil, []int{}, []int{0}, []int{0, 0}, []int{0, 0, 0}} {
				b503ControlFieldT.Run(b503ControlFieldFmt.Sprint(b503ControlFieldCase["ab07aq07id"], b503ControlFieldKey, b503ControlFieldReplacement), func(b503ControlFieldT *b503ControlFieldTesting.T) {
					b503ControlFieldBytes, _ := b503ControlFieldJSON.Marshal(b503ControlFieldOriginal)
					b503ControlFieldInput := map[string]any{}
					_ = b503ControlFieldJSON.Unmarshal(b503ControlFieldBytes, &b503ControlFieldInput)
					b503ControlFieldFields := b503ControlFieldInput["fields"].(map[string]any)
					if b503ControlFieldReplacement == nil {
						delete(b503ControlFieldFields, b503ControlFieldKey)
					} else {
						b503ControlFieldFields[b503ControlFieldKey] = b503ControlFieldReplacement
					}
					b503ControlFieldOutput := b503ControlFieldEvaluateTest(b503ControlFieldT, b503ControlFieldInput)
					b503ControlFieldControl := b503ControlFieldEvaluateTest(b503ControlFieldT, b503ControlFieldOriginal)
					b503ControlFieldResults := b503ControlFieldOutput["results"].(map[string]any)
					for b503ControlFieldSibling, b503ControlFieldValue := range b503ControlFieldControl["results"].(map[string]any) {
						if b503ControlFieldSibling != b503ControlFieldKey && !b503ControlFieldReflect.DeepEqual(b503ControlFieldValue, b503ControlFieldResults[b503ControlFieldSibling]) {
							b503ControlFieldT.Fatalf("changed sibling %s", b503ControlFieldSibling)
						}
					}
					b503ControlFieldWidth := 1
					if b503ControlFieldStrings.HasPrefix(b503ControlFieldKey, "error") {
						b503ControlFieldWidth = 2
					}
					b503ControlFieldArray, _ := b503ControlFieldReplacement.([]int)
					b503ControlFieldResult := b503ControlFieldResults[b503ControlFieldKey].(map[string]any)
					if len(b503ControlFieldArray) != b503ControlFieldWidth {
						if b503ControlFieldOutput["state"] != "partial" || b503ControlFieldResult["state"] != "incomplete" {
							b503ControlFieldT.Fatal("bad width complete")
						}
						if _, b503ControlFieldPresent := b503ControlFieldResult["number"]; b503ControlFieldPresent {
							b503ControlFieldT.Fatal("invented number")
						}
					} else if b503ControlFieldResult["number"] != float64(0) {
						b503ControlFieldT.Fatal("field not independently evaluated")
					}
					if b503ControlFieldResult["present"] != (b503ControlFieldReplacement != nil) {
						b503ControlFieldT.Fatal("absence conflated")
					}
				})
			}
		}
	}
}

func TestAb07Aq07Reject(b503ControlFieldT *b503ControlFieldTesting.T) {
	for _, b503ControlFieldInput := range []string{``, `{`, `[]`, `null`, `{} {}`, `{"context":{},"context":{}}`, `{"x":[{"a":0,"a":1}]}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":null}`} {
		if b503ControlFieldResult, b503ControlFieldError := B503ControlFieldEvaluate([]byte(b503ControlFieldInput)); b503ControlFieldError == nil || b503ControlFieldResult != nil {
			b503ControlFieldT.Fatalf("accepted %s", b503ControlFieldInput)
		}
	}
	for _, b503ControlFieldBad := range []string{`null`, `true`, `"1"`, `1.5`, `-1`, `256`, `1e999`, `NaN`} {
		for _, b503ControlFieldKey := range []string{"primary", "secondary", "operation", "field", "arguments", "direction", "view"} {
			b503ControlFieldInput := `{"context":{"primary":181,"secondary":3,"operation":3,"direction":"positive","view":"reset-acceptance"},"fields":{"acceptance":[1]}}`
			switch b503ControlFieldKey {
			case "field":
				b503ControlFieldInput = b503ControlFieldStrings.Replace(b503ControlFieldInput, `[1]`, `[`+b503ControlFieldBad+`]`, 1)
			case "arguments":
				b503ControlFieldInput = b503ControlFieldStrings.TrimSuffix(b503ControlFieldInput, "}") + `,"arguments":{"undefined":` + b503ControlFieldBad + `}}`
			default:
				if (b503ControlFieldKey == "direction" || b503ControlFieldKey == "view") && b503ControlFieldBad == `"1"` {
					continue
				}
				b503ControlFieldOld := map[string]string{"primary": "181", "secondary": "3", "operation": "3", "direction": `"positive"`, "view": `"reset-acceptance"`}[b503ControlFieldKey]
				b503ControlFieldInput = b503ControlFieldStrings.Replace(b503ControlFieldInput, `"`+b503ControlFieldKey+`":`+b503ControlFieldOld, `"`+b503ControlFieldKey+`":`+b503ControlFieldBad, 1)
			}
			if b503ControlFieldResult, b503ControlFieldError := B503ControlFieldEvaluate([]byte(b503ControlFieldInput)); b503ControlFieldError == nil || b503ControlFieldResult != nil {
				b503ControlFieldT.Fatalf("accepted %s", b503ControlFieldInput)
			}
		}
	}
}

func TestAb07Aq07MemoryAndExtras(b503ControlFieldT *b503ControlFieldTesting.T) {
	b503ControlFieldInput := []byte(`{"context":{"primary":181,"secondary":3,"operation":3,"direction":"positive","view":"reset-acceptance"},"fields":{"acceptance":[1],"opaque":[]},"arguments":{},"extra":[null,0,false,[]]}`)
	b503ControlFieldControl, b503ControlFieldError := B503ControlFieldEvaluate(b503ControlFieldInput)
	if b503ControlFieldError != nil {
		b503ControlFieldT.Fatal(b503ControlFieldError)
	}
	b503ControlFieldOriginal := append([]byte(nil), b503ControlFieldInput...)
	b503ControlFieldOutput, _ := B503ControlFieldEvaluate(b503ControlFieldInput)
	for b503ControlFieldIndex := range b503ControlFieldInput {
		b503ControlFieldInput[b503ControlFieldIndex] = 'x'
	}
	if string(b503ControlFieldOutput) != string(b503ControlFieldControl) {
		b503ControlFieldT.Fatal("input aliases output")
	}
	for b503ControlFieldIndex := range b503ControlFieldOutput {
		b503ControlFieldOutput[b503ControlFieldIndex] = 'x'
	}
	b503ControlFieldAgain, b503ControlFieldError := B503ControlFieldEvaluate(b503ControlFieldOriginal)
	if b503ControlFieldError != nil || string(b503ControlFieldAgain) != string(b503ControlFieldControl) {
		b503ControlFieldT.Fatal("cross-call alias")
	}
	b503ControlFieldDecoded := map[string]any{}
	_ = b503ControlFieldJSON.Unmarshal(b503ControlFieldAgain, &b503ControlFieldDecoded)
	if !b503ControlFieldReflect.DeepEqual(b503ControlFieldDecoded["extra"], []any{nil, float64(0), false, []any{}}) {
		b503ControlFieldT.Fatal("extras lost")
	}
}

func TestAb07Aq07FeedbackDomain(b503ControlFieldT *b503ControlFieldTesting.T) {
	for b503ControlFieldNumber := 0; b503ControlFieldNumber < 256; b503ControlFieldNumber++ {
		b503ControlFieldOutput := b503ControlFieldEvaluateTest(b503ControlFieldT, map[string]any{"context": map[string]any{"primary": 181, "secondary": 3, "operation": 7, "direction": "positive", "view": "current-errors"}, "fields": map[string]any{"feedback": []int{b503ControlFieldNumber}}})
		b503ControlFieldResult := b503ControlFieldOutput["results"].(map[string]any)["feedback"].(map[string]any)
		b503ControlFieldLabel := "rejected"
		if b503ControlFieldNumber == 1 {
			b503ControlFieldLabel = "accepted"
		}
		if b503ControlFieldResult["label"] != b503ControlFieldLabel || b503ControlFieldOutput["state"] != "partial" {
			b503ControlFieldT.Fatalf("feedback %d: %v", b503ControlFieldNumber, b503ControlFieldOutput)
		}
		if len(b503ControlFieldOutput) != 6 || len(b503ControlFieldResult) != 7 {
			b503ControlFieldT.Fatal("feedback inferred effects")
		}
	}
}

func TestAb07Aq07KindsAndPreservation(b503ControlFieldT *b503ControlFieldTesting.T) {
	for _, b503ControlFieldBad := range []string{`null`, `true`, `1`, `"x"`, `{}`} {
		for _, b503ControlFieldInput := range []string{`{"context":{},"fields":{"opaque":` + b503ControlFieldBad + `}}`, `{"context":{},"fields":{},"arguments":` + b503ControlFieldBad + `}`} {
			if b503ControlFieldBad == `{}` && b503ControlFieldStrings.Contains(b503ControlFieldInput, `"arguments"`) {
				b503ControlFieldResult, b503ControlFieldError := B503ControlFieldEvaluate([]byte(b503ControlFieldInput))
				if b503ControlFieldError != nil || b503ControlFieldResult == nil {
					b503ControlFieldT.Fatal("empty arguments rejected")
				}
				continue
			}
			if b503ControlFieldResult, b503ControlFieldError := B503ControlFieldEvaluate([]byte(b503ControlFieldInput)); b503ControlFieldError == nil || b503ControlFieldResult != nil {
				b503ControlFieldT.Fatalf("accepted %s", b503ControlFieldInput)
			}
		}
	}
	b503ControlFieldInput := map[string]any{"context": map[string]any{"primary": float64(181), "secondary": float64(3), "operation": float64(0), "direction": "positive", "view": "reset-acceptance", "variant": nil}, "fields": map[string]any{"acceptance": []any{}, "opaque": []any{float64(0)}}, "arguments": map[string]any{}, "extra": []any{nil, false, float64(0), []any{}}}
	b503ControlFieldOutput := b503ControlFieldEvaluateTest(b503ControlFieldT, b503ControlFieldInput)
	for b503ControlFieldKey, b503ControlFieldValue := range b503ControlFieldInput {
		if !b503ControlFieldReflect.DeepEqual(b503ControlFieldValue, b503ControlFieldOutput[b503ControlFieldKey]) {
			b503ControlFieldT.Fatalf("unsupported context lost %s", b503ControlFieldKey)
		}
	}
	if b503ControlFieldOutput["state"] != "Unknown" || len(b503ControlFieldOutput["results"].(map[string]any)) != 0 {
		b503ControlFieldT.Fatal("unsupported context evaluated")
	}
}

func TestAb07Aq07RejectReservedOutputKeys(b503ControlFieldT *b503ControlFieldTesting.T) {
	for _, b503ControlFieldKey := range []string{"results", "state", "physical", "qualified"} {
		b503ControlFieldInput := `{"context":{},"fields":{},"` + b503ControlFieldKey + `":{"opaque":1}}`
		if b503ControlFieldResult, b503ControlFieldError := B503ControlFieldEvaluate([]byte(b503ControlFieldInput)); b503ControlFieldError == nil || b503ControlFieldResult != nil {
			b503ControlFieldT.Fatalf("reserved key %s accepted", b503ControlFieldKey)
		}
	}
}
