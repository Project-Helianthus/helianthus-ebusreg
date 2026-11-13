package vaillant

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func b512PowerStepperFieldCall(b512PowerStepperFieldT *testing.T, b512PowerStepperFieldInput any) map[string]any {
	b512PowerStepperFieldT.Helper()
	b512PowerStepperFieldBytes, b512PowerStepperFieldErr := json.Marshal(b512PowerStepperFieldInput)
	if b512PowerStepperFieldErr != nil {
		b512PowerStepperFieldT.Fatal(b512PowerStepperFieldErr)
	}
	b512PowerStepperFieldOutput, b512PowerStepperFieldErr := B512PowerStepperFieldEvaluate(b512PowerStepperFieldBytes)
	if b512PowerStepperFieldErr != nil {
		b512PowerStepperFieldT.Fatal(b512PowerStepperFieldErr)
	}
	b512PowerStepperFieldObject := map[string]any{}
	if b512PowerStepperFieldErr = json.Unmarshal(b512PowerStepperFieldOutput, &b512PowerStepperFieldObject); b512PowerStepperFieldErr != nil {
		b512PowerStepperFieldT.Fatal(b512PowerStepperFieldErr)
	}
	return b512PowerStepperFieldObject
}

func b512PowerStepperFieldContains(b512PowerStepperFieldT *testing.T, b512PowerStepperFieldActual, b512PowerStepperFieldExpected any) {
	b512PowerStepperFieldT.Helper()
	if b512PowerStepperFieldObject, b512PowerStepperFieldOK := b512PowerStepperFieldExpected.(map[string]any); b512PowerStepperFieldOK {
		b512PowerStepperFieldGot, b512PowerStepperFieldOK := b512PowerStepperFieldActual.(map[string]any)
		if !b512PowerStepperFieldOK {
			b512PowerStepperFieldT.Fatalf("expected object, got %v", b512PowerStepperFieldActual)
		}
		for b512PowerStepperFieldKey, b512PowerStepperFieldValue := range b512PowerStepperFieldObject {
			b512PowerStepperFieldContains(b512PowerStepperFieldT, b512PowerStepperFieldGot[b512PowerStepperFieldKey], b512PowerStepperFieldValue)
		}
	} else if !reflect.DeepEqual(b512PowerStepperFieldActual, b512PowerStepperFieldExpected) {
		b512PowerStepperFieldT.Fatalf("got %#v want %#v", b512PowerStepperFieldActual, b512PowerStepperFieldExpected)
	}
}

func b512PowerStepperFieldInput() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 18, "operation": 8, "direction": "request", "view": "power-stepper", "family": "synthetic", "variant": "candidate"}, "fields": map[string]any{"inlet": []int{240, 255}, "status": []int{85}, "stepper": []int{1}}}
}

func TestAb25Aq25Fixture(b512PowerStepperFieldT *testing.T) {
	b512PowerStepperFieldBytes, b512PowerStepperFieldErr := os.ReadFile("testdata/ab25aq25_b512_power_stepper_fields.json")
	if b512PowerStepperFieldErr != nil {
		b512PowerStepperFieldT.Fatal(b512PowerStepperFieldErr)
	}
	var b512PowerStepperFieldCases []map[string]any
	if b512PowerStepperFieldErr = json.Unmarshal(b512PowerStepperFieldBytes, &b512PowerStepperFieldCases); b512PowerStepperFieldErr != nil {
		b512PowerStepperFieldT.Fatal(b512PowerStepperFieldErr)
	}
	if len(b512PowerStepperFieldCases) == 0 {
		b512PowerStepperFieldT.Fatal("zero fixtures")
	}
	for _, b512PowerStepperFieldCase := range b512PowerStepperFieldCases {
		b512PowerStepperFieldT.Run(b512PowerStepperFieldCase["id"].(string), func(b512PowerStepperFieldT *testing.T) {
			b512PowerStepperFieldGot := b512PowerStepperFieldCall(b512PowerStepperFieldT, b512PowerStepperFieldCase["input"])
			b512PowerStepperFieldContains(b512PowerStepperFieldT, b512PowerStepperFieldGot, b512PowerStepperFieldCase["expected"])
			if b512PowerStepperFieldGot["physical"] != "Unknown" || b512PowerStepperFieldGot["qualified"] != false {
				b512PowerStepperFieldT.Fatal(b512PowerStepperFieldGot)
			}
			b512PowerStepperFieldResults := b512PowerStepperFieldGot["results"].(map[string]any)
			for _, b512PowerStepperFieldResult := range b512PowerStepperFieldResults {
				b512PowerStepperFieldMap := b512PowerStepperFieldResult.(map[string]any)
				if b512PowerStepperFieldMap["state"] == "complete" && b512PowerStepperFieldMap["grade"] != "Hypothesis" {
					b512PowerStepperFieldT.Fatal(b512PowerStepperFieldMap)
				}
				if b512PowerStepperFieldFlags, b512PowerStepperFieldOK := b512PowerStepperFieldMap["flags"].(map[string]any); b512PowerStepperFieldOK {
					if len(b512PowerStepperFieldFlags) != 1 {
						b512PowerStepperFieldT.Fatal("release must not imply movement", b512PowerStepperFieldFlags)
					}
					if _, b512PowerStepperFieldOK = b512PowerStepperFieldFlags["release"]; !b512PowerStepperFieldOK {
						b512PowerStepperFieldT.Fatal(b512PowerStepperFieldFlags)
					}
					for _, b512PowerStepperFieldKey := range []string{"number", "label", "movement", "contact"} {
						if _, b512PowerStepperFieldOK = b512PowerStepperFieldMap[b512PowerStepperFieldKey]; b512PowerStepperFieldOK {
							b512PowerStepperFieldT.Fatal(b512PowerStepperFieldMap)
						}
					}
				}
			}
		})
	}
}

func TestAb25Aq25Selectors(b512PowerStepperFieldT *testing.T) {
	for b512PowerStepperFieldKey, b512PowerStepperFieldValues := range map[string][]any{"primary": {180, nil}, "secondary": {19, nil}, "operation": {0, 10, 20, nil}, "direction": {"positive", "other", nil}, "view": {"other", nil}, "family": {"other", nil, 0, false}, "variant": {"other", nil, 0, false}} {
		for b512PowerStepperFieldIndex, b512PowerStepperFieldValue := range b512PowerStepperFieldValues {
			b512PowerStepperFieldT.Run(fmt.Sprintf("%s_%d", b512PowerStepperFieldKey, b512PowerStepperFieldIndex), func(b512PowerStepperFieldT *testing.T) {
				b512PowerStepperFieldIn := b512PowerStepperFieldInput()
				b512PowerStepperFieldContext := b512PowerStepperFieldIn["context"].(map[string]any)
				if b512PowerStepperFieldValue == nil {
					delete(b512PowerStepperFieldContext, b512PowerStepperFieldKey)
				} else {
					b512PowerStepperFieldContext[b512PowerStepperFieldKey] = b512PowerStepperFieldValue
				}
				b512PowerStepperFieldGot := b512PowerStepperFieldCall(b512PowerStepperFieldT, b512PowerStepperFieldIn)
				if b512PowerStepperFieldGot["qualified"] != false {
					b512PowerStepperFieldT.Fatal(b512PowerStepperFieldGot)
				}
				if b512PowerStepperFieldKey == "family" || b512PowerStepperFieldKey == "variant" {
					if b512PowerStepperFieldGot["state"] != "complete" {
						b512PowerStepperFieldT.Fatal(b512PowerStepperFieldGot)
					}
				} else if b512PowerStepperFieldKey == "direction" && b512PowerStepperFieldValue == "positive" {
					if b512PowerStepperFieldGot["state"] != "partial" {
						b512PowerStepperFieldT.Fatal(b512PowerStepperFieldGot)
					}
					if b512PowerStepperFieldGot["results"].(map[string]any)["status"].(map[string]any)["state"] != "Unknown" {
						b512PowerStepperFieldT.Fatal(b512PowerStepperFieldGot)
					}
				} else {
					if b512PowerStepperFieldGot["state"] != "Unknown" {
						b512PowerStepperFieldT.Fatal(b512PowerStepperFieldGot)
					}
					for _, b512PowerStepperFieldResult := range b512PowerStepperFieldGot["results"].(map[string]any) {
						if len(b512PowerStepperFieldResult.(map[string]any)) != 2 {
							b512PowerStepperFieldT.Fatal(b512PowerStepperFieldResult)
						}
					}
				}
				b512PowerStepperFieldBytes, _ := json.Marshal(b512PowerStepperFieldContext)
				b512PowerStepperFieldExpected := map[string]any{}
				if b512PowerStepperFieldErr := json.Unmarshal(b512PowerStepperFieldBytes, &b512PowerStepperFieldExpected); b512PowerStepperFieldErr != nil {
					b512PowerStepperFieldT.Fatal(b512PowerStepperFieldErr)
				}
				if !reflect.DeepEqual(b512PowerStepperFieldGot["context"], b512PowerStepperFieldExpected) {
					b512PowerStepperFieldT.Fatal("context retention")
				}
			})
		}
	}
}

func TestAb25Aq25Invalid(b512PowerStepperFieldT *testing.T) {
	for _, b512PowerStepperFieldText := range []string{``, `[]`, `{`, `{} {}`, `{"context":{},"fields":{},"extra":1e999}`, `{"context":{},"fields":{"x":[],"x":[]}}`, `{"context":{"operation":8,"operation":8},"fields":{}}`, `{"context":{},"fields":{},"arguments":{"x":0,"x":1}}`, `{"context":{},"fields":{},"x":{"a":1,"a":2}}`, `{"context":{},"fields":{},"x":[{"a":1,"a":2}]}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{},"fields":{"a":[NaN]}}`} {
		b512PowerStepperFieldOut, b512PowerStepperFieldErr := B512PowerStepperFieldEvaluate(json.RawMessage(b512PowerStepperFieldText))
		if b512PowerStepperFieldErr == nil || b512PowerStepperFieldOut != nil {
			b512PowerStepperFieldT.Fatalf("accepted %s", b512PowerStepperFieldText)
		}
	}
	for _, b512PowerStepperFieldKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, b512PowerStepperFieldValue := range []any{nil, true, []int{}, map[string]any{}, 1.5, -1, 256, "8"} {
			if (b512PowerStepperFieldKey == "direction" || b512PowerStepperFieldKey == "view") && b512PowerStepperFieldValue == "8" {
				continue
			}
			b512PowerStepperFieldIn := b512PowerStepperFieldInput()
			b512PowerStepperFieldIn["context"].(map[string]any)[b512PowerStepperFieldKey] = b512PowerStepperFieldValue
			b512PowerStepperFieldBytes, _ := json.Marshal(b512PowerStepperFieldIn)
			b512PowerStepperFieldOut, b512PowerStepperFieldErr := B512PowerStepperFieldEvaluate(b512PowerStepperFieldBytes)
			if b512PowerStepperFieldErr == nil || b512PowerStepperFieldOut != nil {
				b512PowerStepperFieldT.Fatalf("accepted %s", b512PowerStepperFieldBytes)
			}
		}
	}
	for _, b512PowerStepperFieldValue := range []any{nil, true, "1", 1, []any{nil}, []any{true}, []any{"1"}, []any{1.5}, []any{-1}, []any{256}} {
		b512PowerStepperFieldIn := b512PowerStepperFieldInput()
		b512PowerStepperFieldIn["fields"].(map[string]any)["inlet"] = b512PowerStepperFieldValue
		b512PowerStepperFieldBytes, _ := json.Marshal(b512PowerStepperFieldIn)
		b512PowerStepperFieldOut, b512PowerStepperFieldErr := B512PowerStepperFieldEvaluate(b512PowerStepperFieldBytes)
		if b512PowerStepperFieldErr == nil || b512PowerStepperFieldOut != nil {
			b512PowerStepperFieldT.Fatalf("accepted %s", b512PowerStepperFieldBytes)
		}
	}
	for _, b512PowerStepperFieldKey := range []string{"context", "fields", "arguments"} {
		for _, b512PowerStepperFieldValue := range []any{nil, 0, false, "", []any{}} {
			b512PowerStepperFieldIn := b512PowerStepperFieldInput()
			b512PowerStepperFieldIn[b512PowerStepperFieldKey] = b512PowerStepperFieldValue
			b512PowerStepperFieldBytes, _ := json.Marshal(b512PowerStepperFieldIn)
			b512PowerStepperFieldOut, b512PowerStepperFieldErr := B512PowerStepperFieldEvaluate(b512PowerStepperFieldBytes)
			if b512PowerStepperFieldErr == nil || b512PowerStepperFieldOut != nil {
				b512PowerStepperFieldT.Fatalf("accepted %s", b512PowerStepperFieldBytes)
			}
		}
	}
}

func TestAb25Aq25WidthsAndRetention(b512PowerStepperFieldT *testing.T) {
	for _, b512PowerStepperFieldSpec := range []struct {
		b512PowerStepperFieldOp        int
		b512PowerStepperFieldDirection string
		b512PowerStepperFieldField     string
		b512PowerStepperFieldWidth     int
	}{{7, "positive", "maximum_power", 1}, {7, "positive", "minimum_power", 1}, {8, "request", "inlet", 2}, {8, "request", "status", 1}, {8, "request", "stepper", 1}, {8, "positive", "setpoint", 2}, {8, "positive", "stepper", 1}, {9, "request", "release", 1}, {9, "positive", "release", 1}} {
		for _, b512PowerStepperFieldWidth := range []int{-1, 0, b512PowerStepperFieldSpec.b512PowerStepperFieldWidth - 1, b512PowerStepperFieldSpec.b512PowerStepperFieldWidth, b512PowerStepperFieldSpec.b512PowerStepperFieldWidth + 1} {
			b512PowerStepperFieldT.Run(fmt.Sprintf("%d_%s_%s_%d", b512PowerStepperFieldSpec.b512PowerStepperFieldOp, b512PowerStepperFieldSpec.b512PowerStepperFieldDirection, b512PowerStepperFieldSpec.b512PowerStepperFieldField, b512PowerStepperFieldWidth), func(b512PowerStepperFieldT *testing.T) {
				b512PowerStepperFieldIn := b512PowerStepperFieldInput()
				b512PowerStepperFieldContext := b512PowerStepperFieldIn["context"].(map[string]any)
				b512PowerStepperFieldContext["operation"] = b512PowerStepperFieldSpec.b512PowerStepperFieldOp
				b512PowerStepperFieldContext["direction"] = b512PowerStepperFieldSpec.b512PowerStepperFieldDirection
				b512PowerStepperFieldFields := map[string]any{"maximum_power": []int{255}, "minimum_power": []int{250}, "inlet": []int{240, 255}, "status": []int{85}, "stepper": []int{1}, "setpoint": []int{16, 0}, "release": []int{1}, "opaque": []int{}}
				b512PowerStepperFieldIn["fields"] = b512PowerStepperFieldFields
				if b512PowerStepperFieldWidth < 0 {
					delete(b512PowerStepperFieldFields, b512PowerStepperFieldSpec.b512PowerStepperFieldField)
				} else {
					b512PowerStepperFieldFields[b512PowerStepperFieldSpec.b512PowerStepperFieldField] = make([]int, b512PowerStepperFieldWidth)
				}
				b512PowerStepperFieldIn["extra"] = map[string]any{"null": nil, "zero": 0, "false": false, "empty": []any{}}
				b512PowerStepperFieldIn["arguments"] = map[string]any{"opaque": nil}
				b512PowerStepperFieldGot := b512PowerStepperFieldCall(b512PowerStepperFieldT, b512PowerStepperFieldIn)
				b512PowerStepperFieldResults := b512PowerStepperFieldGot["results"].(map[string]any)
				b512PowerStepperFieldResult := b512PowerStepperFieldResults[b512PowerStepperFieldSpec.b512PowerStepperFieldField].(map[string]any)
				if b512PowerStepperFieldWidth == b512PowerStepperFieldSpec.b512PowerStepperFieldWidth {
					if b512PowerStepperFieldSpec.b512PowerStepperFieldField == "maximum_power" || b512PowerStepperFieldSpec.b512PowerStepperFieldField == "minimum_power" || b512PowerStepperFieldSpec.b512PowerStepperFieldField == "inlet" || b512PowerStepperFieldSpec.b512PowerStepperFieldField == "setpoint" {
						if b512PowerStepperFieldResult["number"] != float64(0) {
							b512PowerStepperFieldT.Fatal("zero arithmetic", b512PowerStepperFieldResult)
						}
					}
					if b512PowerStepperFieldGot["state"] != "complete" || b512PowerStepperFieldResult["state"] != "complete" {
						b512PowerStepperFieldT.Fatal(b512PowerStepperFieldGot)
					}
				} else {
					if b512PowerStepperFieldGot["state"] != "partial" || b512PowerStepperFieldResult["state"] != "incomplete" {
						b512PowerStepperFieldT.Fatal(b512PowerStepperFieldGot)
					}
					if _, b512PowerStepperFieldOK := b512PowerStepperFieldResult["number"]; b512PowerStepperFieldOK {
						b512PowerStepperFieldT.Fatal(b512PowerStepperFieldResult)
					}
				}
				_, b512PowerStepperFieldPresent := b512PowerStepperFieldResult["raw"]
				if b512PowerStepperFieldPresent != (b512PowerStepperFieldWidth >= 0) {
					b512PowerStepperFieldT.Fatal(b512PowerStepperFieldResult)
				}
				b512PowerStepperFieldCanonicalBytes, _ := json.Marshal(b512PowerStepperFieldIn)
				b512PowerStepperFieldCanonical := map[string]any{}
				if b512PowerStepperFieldErr := json.Unmarshal(b512PowerStepperFieldCanonicalBytes, &b512PowerStepperFieldCanonical); b512PowerStepperFieldErr != nil {
					b512PowerStepperFieldT.Fatal(b512PowerStepperFieldErr)
				}
				for _, b512PowerStepperFieldKey := range []string{"context", "fields", "extra", "arguments"} {
					if !reflect.DeepEqual(b512PowerStepperFieldGot[b512PowerStepperFieldKey], b512PowerStepperFieldCanonical[b512PowerStepperFieldKey]) {
						b512PowerStepperFieldT.Fatal("retention", b512PowerStepperFieldKey)
					}
				}
				if b512PowerStepperFieldSpec.b512PowerStepperFieldOp == 7 && b512PowerStepperFieldSpec.b512PowerStepperFieldField == "maximum_power" {
					b512PowerStepperFieldContains(b512PowerStepperFieldT, b512PowerStepperFieldResults["minimum_power"], map[string]any{"number": float64(125)})
				}
				if b512PowerStepperFieldSpec.b512PowerStepperFieldOp == 8 && b512PowerStepperFieldSpec.b512PowerStepperFieldField != "stepper" {
					b512PowerStepperFieldContains(b512PowerStepperFieldT, b512PowerStepperFieldResults["stepper"], map[string]any{"state": "complete"})
				}
			})
		}
	}
}

func TestAb25Aq25Memory(b512PowerStepperFieldT *testing.T) {
	b512PowerStepperFieldInputBytes, _ := json.Marshal(b512PowerStepperFieldInput())
	b512PowerStepperFieldFirst, b512PowerStepperFieldErr := B512PowerStepperFieldEvaluate(b512PowerStepperFieldInputBytes)
	if b512PowerStepperFieldErr != nil {
		b512PowerStepperFieldT.Fatal(b512PowerStepperFieldErr)
	}
	b512PowerStepperFieldSaved := string(b512PowerStepperFieldFirst)
	for b512PowerStepperFieldIndex := range b512PowerStepperFieldInputBytes {
		b512PowerStepperFieldInputBytes[b512PowerStepperFieldIndex] = ' '
	}
	if string(b512PowerStepperFieldFirst) != b512PowerStepperFieldSaved {
		b512PowerStepperFieldT.Fatal("input alias")
	}
	for b512PowerStepperFieldIndex := range b512PowerStepperFieldFirst {
		b512PowerStepperFieldFirst[b512PowerStepperFieldIndex] = ' '
	}
	b512PowerStepperFieldNextBytes, _ := json.Marshal(b512PowerStepperFieldInput())
	b512PowerStepperFieldNext, b512PowerStepperFieldErr := B512PowerStepperFieldEvaluate(b512PowerStepperFieldNextBytes)
	if b512PowerStepperFieldErr != nil || string(b512PowerStepperFieldNext) != b512PowerStepperFieldSaved {
		b512PowerStepperFieldT.Fatal("cross-call alias")
	}
	if strings.Contains(string(b512PowerStepperFieldNext), "movement") {
		b512PowerStepperFieldT.Fatal("invented movement")
	}
}

func TestAb25Aq25OpaqueIdentity(b512PowerStepperFieldT *testing.T) {
	for _, b512PowerStepperFieldKey := range []string{"family", "variant"} {
		for _, b512PowerStepperFieldValue := range []any{nil, false, 0, "", []any{}, map[string]any{}} {
			b512PowerStepperFieldIn := b512PowerStepperFieldInput()
			b512PowerStepperFieldIn["context"].(map[string]any)[b512PowerStepperFieldKey] = b512PowerStepperFieldValue
			b512PowerStepperFieldGot := b512PowerStepperFieldCall(b512PowerStepperFieldT, b512PowerStepperFieldIn)
			b512PowerStepperFieldBytes, _ := json.Marshal(b512PowerStepperFieldIn["context"])
			b512PowerStepperFieldExpected := map[string]any{}
			if b512PowerStepperFieldErr := json.Unmarshal(b512PowerStepperFieldBytes, &b512PowerStepperFieldExpected); b512PowerStepperFieldErr != nil {
				b512PowerStepperFieldT.Fatal(b512PowerStepperFieldErr)
			}
			if !reflect.DeepEqual(b512PowerStepperFieldGot["context"], b512PowerStepperFieldExpected) || b512PowerStepperFieldGot["qualified"] != false || b512PowerStepperFieldGot["state"] != "complete" {
				b512PowerStepperFieldT.Fatal(b512PowerStepperFieldGot)
			}
		}
	}
}
