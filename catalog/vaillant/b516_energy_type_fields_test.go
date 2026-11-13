package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func b516EnergyTypeFieldInput() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 22, "operation": nil, "direction": "read-positive", "view": "energy-type", "family": "synthetic-family", "variant": "synthetic-variant", "opaque": map[string]any{"nested": []any{false, nil, 0}}}, "fields": map[string]any{"value": []any{193}}, "arguments": map[string]any{"opaque": false}}
}

func b516EnergyTypeFieldRun(b516EnergyTypeFieldT *testing.T, b516EnergyTypeFieldInputValue any) map[string]any {
	b516EnergyTypeFieldT.Helper()
	b516EnergyTypeFieldEncoded, b516EnergyTypeFieldErr := json.Marshal(b516EnergyTypeFieldInputValue)
	if b516EnergyTypeFieldErr != nil {
		b516EnergyTypeFieldT.Fatal(b516EnergyTypeFieldErr)
	}
	b516EnergyTypeFieldOutput, b516EnergyTypeFieldErr := B516EnergyTypeFieldEvaluate(b516EnergyTypeFieldEncoded)
	if b516EnergyTypeFieldErr != nil {
		b516EnergyTypeFieldT.Fatal(b516EnergyTypeFieldErr)
	}
	b516EnergyTypeFieldResult := map[string]any{}
	if b516EnergyTypeFieldErr := json.Unmarshal(b516EnergyTypeFieldOutput, &b516EnergyTypeFieldResult); b516EnergyTypeFieldErr != nil {
		b516EnergyTypeFieldT.Fatal(b516EnergyTypeFieldErr)
	}
	if b516EnergyTypeFieldResult["physical"] != "Unknown" || b516EnergyTypeFieldResult["qualified"] != false {
		b516EnergyTypeFieldT.Fatal(b516EnergyTypeFieldResult)
	}
	return b516EnergyTypeFieldResult
}

func b516EnergyTypeFieldValue(b516EnergyTypeFieldOutput map[string]any) map[string]any {
	return b516EnergyTypeFieldOutput["fields"].(map[string]any)["value"].(map[string]any)
}

func TestAb23Ap23Directions(b516EnergyTypeFieldT *testing.T) {
	b516EnergyTypeFieldData, b516EnergyTypeFieldErr := os.ReadFile("testdata/ab23ap23_b516_energy_type_fields.json")
	if b516EnergyTypeFieldErr != nil {
		b516EnergyTypeFieldT.Fatal(b516EnergyTypeFieldErr)
	}
	var b516EnergyTypeFieldCases []map[string]any
	if b516EnergyTypeFieldErr := json.Unmarshal(b516EnergyTypeFieldData, &b516EnergyTypeFieldCases); b516EnergyTypeFieldErr != nil {
		b516EnergyTypeFieldT.Fatal(b516EnergyTypeFieldErr)
	}
	if len(b516EnergyTypeFieldCases) != 8 {
		b516EnergyTypeFieldT.Fatal("missing vectors")
	}
	for _, b516EnergyTypeFieldCase := range b516EnergyTypeFieldCases {
		for _, b516EnergyTypeFieldDirection := range []string{"request", "read-positive", "write-positive"} {
			b516EnergyTypeFieldT.Run(b516EnergyTypeFieldCase["id"].(string)+"/"+b516EnergyTypeFieldDirection, func(b516EnergyTypeFieldT *testing.T) {
				b516EnergyTypeFieldIn := b516EnergyTypeFieldInput()
				b516EnergyTypeFieldIn["context"].(map[string]any)["direction"] = b516EnergyTypeFieldDirection
				b516EnergyTypeFieldIn["fields"].(map[string]any)["value"] = []any{b516EnergyTypeFieldCase["value"]}
				b516EnergyTypeFieldOut := b516EnergyTypeFieldRun(b516EnergyTypeFieldT, b516EnergyTypeFieldIn)
				b516EnergyTypeFieldExpected := map[string]any{"raw": []any{b516EnergyTypeFieldCase["value"]}, "present": true, "state": "Hypothesis"}
				switch b516EnergyTypeFieldDirection {
				case "request":
					b516EnergyTypeFieldExpected["number"] = b516EnergyTypeFieldCase["value"]
					b516EnergyTypeFieldExpected["unmapped"] = true
				case "read-positive":
					b516EnergyTypeFieldExpected["low_six"] = b516EnergyTypeFieldCase["low"]
					b516EnergyTypeFieldExpected["flags"] = map[string]any{"flag6": b516EnergyTypeFieldCase["flag6"], "flag7": b516EnergyTypeFieldCase["flag7"]}
					b516EnergyTypeFieldExpected["calculation"] = "Proven"
				case "write-positive":
					b516EnergyTypeFieldExpected["low_six"] = b516EnergyTypeFieldCase["low"]
					b516EnergyTypeFieldExpected["unmodeled_high"] = b516EnergyTypeFieldCase["high"]
					b516EnergyTypeFieldExpected["calculation"] = "Proven"
				}
				if b516EnergyTypeFieldOut["state"] != "complete" || !reflect.DeepEqual(b516EnergyTypeFieldValue(b516EnergyTypeFieldOut), b516EnergyTypeFieldExpected) {
					b516EnergyTypeFieldT.Fatalf("got %#v want %#v", b516EnergyTypeFieldOut, b516EnergyTypeFieldExpected)
				}
			})
		}
	}
}

func TestAb23Ap23Context(b516EnergyTypeFieldT *testing.T) {
	for _, b516EnergyTypeFieldKey := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant", "opaque"} {
		for _, b516EnergyTypeFieldReplacement := range []any{nil, 0, false, "", "other", []any{}, map[string]any{}} {
			b516EnergyTypeFieldIn := b516EnergyTypeFieldInput()
			b516EnergyTypeFieldContext := b516EnergyTypeFieldIn["context"].(map[string]any)
			b516EnergyTypeFieldContext[b516EnergyTypeFieldKey] = b516EnergyTypeFieldReplacement
			b516EnergyTypeFieldBytes, _ := json.Marshal(b516EnergyTypeFieldIn)
			b516EnergyTypeFieldOut, b516EnergyTypeFieldErr := B516EnergyTypeFieldEvaluate(b516EnergyTypeFieldBytes)
			b516EnergyTypeFieldTyped := b516EnergyTypeFieldKey == "primary" || b516EnergyTypeFieldKey == "secondary" || b516EnergyTypeFieldKey == "operation" || b516EnergyTypeFieldKey == "direction" || b516EnergyTypeFieldKey == "view"
			b516EnergyTypeFieldValid := b516EnergyTypeFieldReplacement == nil
			if b516EnergyTypeFieldKey == "direction" || b516EnergyTypeFieldKey == "view" {
				_, b516EnergyTypeFieldValid = b516EnergyTypeFieldReplacement.(string)
				b516EnergyTypeFieldValid = b516EnergyTypeFieldValid || b516EnergyTypeFieldReplacement == nil
			} else if b516EnergyTypeFieldTyped {
				_, b516EnergyTypeFieldNumeric := b516EnergyTypeFieldReplacement.(int)
				b516EnergyTypeFieldValid = b516EnergyTypeFieldValid || b516EnergyTypeFieldNumeric
			}
			if b516EnergyTypeFieldTyped && !b516EnergyTypeFieldValid {
				if b516EnergyTypeFieldErr == nil || b516EnergyTypeFieldOut != nil {
					b516EnergyTypeFieldT.Fatalf("accepted %s=%v", b516EnergyTypeFieldKey, b516EnergyTypeFieldReplacement)
				}
				continue
			}
			if b516EnergyTypeFieldErr != nil {
				b516EnergyTypeFieldT.Fatal(b516EnergyTypeFieldErr)
			}
			b516EnergyTypeFieldResult := b516EnergyTypeFieldRun(b516EnergyTypeFieldT, b516EnergyTypeFieldIn)
			b516EnergyTypeFieldRoundTrip, _ := json.Marshal(b516EnergyTypeFieldResult["context"])
			b516EnergyTypeFieldOriginal, _ := json.Marshal(b516EnergyTypeFieldContext)
			if string(b516EnergyTypeFieldRoundTrip) != string(b516EnergyTypeFieldOriginal) {
				b516EnergyTypeFieldT.Fatal("context changed")
			}
			b516EnergyTypeFieldWant := "Unknown"
			if !b516EnergyTypeFieldTyped || (b516EnergyTypeFieldKey == "operation" && b516EnergyTypeFieldReplacement == nil) {
				b516EnergyTypeFieldWant = "complete"
			}
			if b516EnergyTypeFieldWant == "Unknown" {
				b516EnergyTypeFieldExpected := map[string]any{"raw": []any{float64(193)}, "present": true, "state": "Unknown"}
				if !reflect.DeepEqual(b516EnergyTypeFieldValue(b516EnergyTypeFieldResult), b516EnergyTypeFieldExpected) {
					b516EnergyTypeFieldT.Fatal("unsupported context gained interpretation")
				}
			}
			if b516EnergyTypeFieldResult["state"] != b516EnergyTypeFieldWant {
				b516EnergyTypeFieldT.Fatalf("%s=%v: %v", b516EnergyTypeFieldKey, b516EnergyTypeFieldReplacement, b516EnergyTypeFieldResult)
			}
		}
		b516EnergyTypeFieldIn := b516EnergyTypeFieldInput()
		delete(b516EnergyTypeFieldIn["context"].(map[string]any), b516EnergyTypeFieldKey)
		b516EnergyTypeFieldOut := b516EnergyTypeFieldRun(b516EnergyTypeFieldT, b516EnergyTypeFieldIn)
		b516EnergyTypeFieldWant := "Unknown"
		if b516EnergyTypeFieldKey == "family" || b516EnergyTypeFieldKey == "variant" || b516EnergyTypeFieldKey == "opaque" {
			b516EnergyTypeFieldWant = "complete"
		}
		if b516EnergyTypeFieldOut["state"] != b516EnergyTypeFieldWant {
			b516EnergyTypeFieldT.Fatal(b516EnergyTypeFieldOut)
		}
		if _, b516EnergyTypeFieldPresent := b516EnergyTypeFieldOut["context"].(map[string]any)[b516EnergyTypeFieldKey]; b516EnergyTypeFieldPresent {
			b516EnergyTypeFieldT.Fatal("missing key invented")
		}
	}
}

func TestAb23Ap23IncompleteAndOpaque(b516EnergyTypeFieldT *testing.T) {
	for _, b516EnergyTypeFieldRaw := range []any{nil, []any{}, []any{0, 1}} {
		b516EnergyTypeFieldIn := b516EnergyTypeFieldInput()
		b516EnergyTypeFieldFields := b516EnergyTypeFieldIn["fields"].(map[string]any)
		b516EnergyTypeFieldFields["sibling"] = []any{0, 255}
		if b516EnergyTypeFieldRaw == nil {
			delete(b516EnergyTypeFieldFields, "value")
		} else {
			b516EnergyTypeFieldFields["value"] = b516EnergyTypeFieldRaw
		}
		b516EnergyTypeFieldOut := b516EnergyTypeFieldRun(b516EnergyTypeFieldT, b516EnergyTypeFieldIn)
		if b516EnergyTypeFieldOut["state"] != "partial" || b516EnergyTypeFieldValue(b516EnergyTypeFieldOut)["state"] != "incomplete" {
			b516EnergyTypeFieldT.Fatal(b516EnergyTypeFieldOut)
		}
		b516EnergyTypeFieldExpectedValue := map[string]any{"raw": b516EnergyTypeFieldRaw, "present": b516EnergyTypeFieldRaw != nil, "state": "incomplete"}
		b516EnergyTypeFieldExpectedJSON, _ := json.Marshal(b516EnergyTypeFieldExpectedValue)
		b516EnergyTypeFieldActualJSON, _ := json.Marshal(b516EnergyTypeFieldValue(b516EnergyTypeFieldOut))
		if string(b516EnergyTypeFieldExpectedJSON) != string(b516EnergyTypeFieldActualJSON) {
			b516EnergyTypeFieldT.Fatal("incomplete evidence changed")
		}
		b516EnergyTypeFieldExpected, _ := json.Marshal(b516EnergyTypeFieldFields)
		b516EnergyTypeFieldActual, _ := json.Marshal(b516EnergyTypeFieldOut["raw"])
		if string(b516EnergyTypeFieldExpected) != string(b516EnergyTypeFieldActual) {
			b516EnergyTypeFieldT.Fatal("raw lost")
		}
		if b516EnergyTypeFieldOut["fields"].(map[string]any)["sibling"].(map[string]any)["state"] != "Unknown" {
			b516EnergyTypeFieldT.Fatal("invented sibling interpretation")
		}
	}
	for _, b516EnergyTypeFieldArgument := range []any{false, nil, 0, []any{}, map[string]any{}, ""} {
		b516EnergyTypeFieldIn := b516EnergyTypeFieldInput()
		b516EnergyTypeFieldIn["arguments"] = map[string]any{"opaque": b516EnergyTypeFieldArgument}
		b516EnergyTypeFieldOut := b516EnergyTypeFieldRun(b516EnergyTypeFieldT, b516EnergyTypeFieldIn)
		b516EnergyTypeFieldExpected, _ := json.Marshal(b516EnergyTypeFieldIn["arguments"])
		b516EnergyTypeFieldActual, _ := json.Marshal(b516EnergyTypeFieldOut["arguments"])
		if string(b516EnergyTypeFieldExpected) != string(b516EnergyTypeFieldActual) {
			b516EnergyTypeFieldT.Fatal("arguments changed")
		}
	}
	for _, b516EnergyTypeFieldRaw := range []any{[]any{}, []any{0}, []any{255}} {
		b516EnergyTypeFieldIn := b516EnergyTypeFieldInput()
		b516EnergyTypeFieldIn["fields"].(map[string]any)["sibling"] = b516EnergyTypeFieldRaw
		b516EnergyTypeFieldOut := b516EnergyTypeFieldRun(b516EnergyTypeFieldT, b516EnergyTypeFieldIn)
		if b516EnergyTypeFieldOut["state"] != "complete" || b516EnergyTypeFieldValue(b516EnergyTypeFieldOut)["low_six"] != float64(1) {
			b516EnergyTypeFieldT.Fatal("sibling affected value")
		}
	}
	b516EnergyTypeFieldEmpty := b516EnergyTypeFieldInput()
	b516EnergyTypeFieldEmpty["arguments"] = map[string]any{}
	if len(b516EnergyTypeFieldRun(b516EnergyTypeFieldT, b516EnergyTypeFieldEmpty)["arguments"].(map[string]any)) != 0 {
		b516EnergyTypeFieldT.Fatal("empty arguments changed")
	}
	b516EnergyTypeFieldIn := b516EnergyTypeFieldInput()
	delete(b516EnergyTypeFieldIn, "arguments")
	if _, b516EnergyTypeFieldPresent := b516EnergyTypeFieldRun(b516EnergyTypeFieldT, b516EnergyTypeFieldIn)["arguments"]; b516EnergyTypeFieldPresent {
		b516EnergyTypeFieldT.Fatal("arguments invented")
	}
}

func TestAb23Ap23Reject(b516EnergyTypeFieldT *testing.T) {
	for _, b516EnergyTypeFieldBad := range []string{``, `{`, `null`, `[]`, `{}`, `{"context":{},"fields":{},"fields":{}}`, `{"context":{"a":{"b":0,"b":1}},"fields":{}}`, `{"context":{},"fields":{},"arguments":{"a":[{"b":0,"b":1}]}}`, `{"context":{},"fields":{}} false`} {
		if b516EnergyTypeFieldOut, b516EnergyTypeFieldErr := B516EnergyTypeFieldEvaluate([]byte(b516EnergyTypeFieldBad)); b516EnergyTypeFieldErr == nil || b516EnergyTypeFieldOut != nil {
			b516EnergyTypeFieldT.Fatalf("accepted %s", b516EnergyTypeFieldBad)
		}
	}
	for _, b516EnergyTypeFieldBad := range []any{nil, false, "0", 0, []any{true}, []any{nil}, []any{"1"}, []any{-1}, []any{256}, []any{1.5}, []any{[]any{}}} {
		b516EnergyTypeFieldIn := b516EnergyTypeFieldInput()
		b516EnergyTypeFieldIn["fields"].(map[string]any)["value"] = b516EnergyTypeFieldBad
		b516EnergyTypeFieldBytes, _ := json.Marshal(b516EnergyTypeFieldIn)
		if b516EnergyTypeFieldOut, b516EnergyTypeFieldErr := B516EnergyTypeFieldEvaluate(b516EnergyTypeFieldBytes); b516EnergyTypeFieldErr == nil || b516EnergyTypeFieldOut != nil {
			b516EnergyTypeFieldT.Fatalf("accepted %s", b516EnergyTypeFieldBytes)
		}
	}
	for _, b516EnergyTypeFieldKey := range []string{"context", "fields", "arguments"} {
		for _, b516EnergyTypeFieldBad := range []any{nil, false, 0, "", []any{}} {
			b516EnergyTypeFieldIn := b516EnergyTypeFieldInput()
			b516EnergyTypeFieldIn[b516EnergyTypeFieldKey] = b516EnergyTypeFieldBad
			b516EnergyTypeFieldBytes, _ := json.Marshal(b516EnergyTypeFieldIn)
			if b516EnergyTypeFieldOut, b516EnergyTypeFieldErr := B516EnergyTypeFieldEvaluate(b516EnergyTypeFieldBytes); b516EnergyTypeFieldErr == nil || b516EnergyTypeFieldOut != nil {
				b516EnergyTypeFieldT.Fatal("accepted wrong object kind")
			}
		}
	}
	for _, b516EnergyTypeFieldKey := range []string{"primary", "secondary", "operation"} {
		for _, b516EnergyTypeFieldBad := range []any{-1, 256, 1.5, "181", false} {
			b516EnergyTypeFieldIn := b516EnergyTypeFieldInput()
			b516EnergyTypeFieldIn["context"].(map[string]any)[b516EnergyTypeFieldKey] = b516EnergyTypeFieldBad
			b516EnergyTypeFieldBytes, _ := json.Marshal(b516EnergyTypeFieldIn)
			if b516EnergyTypeFieldOut, b516EnergyTypeFieldErr := B516EnergyTypeFieldEvaluate(b516EnergyTypeFieldBytes); b516EnergyTypeFieldErr == nil || b516EnergyTypeFieldOut != nil {
				b516EnergyTypeFieldT.Fatal("accepted invalid numeric context")
			}
		}
	}
}

func TestAb23Ap23Isolation(b516EnergyTypeFieldT *testing.T) {
	b516EnergyTypeFieldBytes, _ := json.Marshal(b516EnergyTypeFieldInput())
	b516EnergyTypeFieldOut, b516EnergyTypeFieldErr := B516EnergyTypeFieldEvaluate(b516EnergyTypeFieldBytes)
	if b516EnergyTypeFieldErr != nil {
		b516EnergyTypeFieldT.Fatal(b516EnergyTypeFieldErr)
	}
	b516EnergyTypeFieldSaved := string(b516EnergyTypeFieldOut)
	for b516EnergyTypeFieldIndex := range b516EnergyTypeFieldBytes {
		b516EnergyTypeFieldBytes[b516EnergyTypeFieldIndex] = ' '
	}
	if string(b516EnergyTypeFieldOut) != b516EnergyTypeFieldSaved {
		b516EnergyTypeFieldT.Fatal("input aliases output")
	}
	b516EnergyTypeFieldDecoded := map[string]any{}
	if b516EnergyTypeFieldErr := json.Unmarshal(b516EnergyTypeFieldOut, &b516EnergyTypeFieldDecoded); b516EnergyTypeFieldErr != nil {
		b516EnergyTypeFieldT.Fatal(b516EnergyTypeFieldErr)
	}
	b516EnergyTypeFieldDecoded["raw"].(map[string]any)["value"].([]any)[0] = float64(0)
	b516EnergyTypeFieldDecoded["context"].(map[string]any)["opaque"].(map[string]any)["nested"].([]any)[0] = true
	b516EnergyTypeFieldDecoded["arguments"].(map[string]any)["opaque"] = true
	if b516EnergyTypeFieldValue(b516EnergyTypeFieldDecoded)["raw"].([]any)[0] != float64(193) {
		b516EnergyTypeFieldT.Fatal("raw aliases semantic raw")
	}
	for b516EnergyTypeFieldIndex := range b516EnergyTypeFieldOut {
		b516EnergyTypeFieldOut[b516EnergyTypeFieldIndex] = ' '
	}
	b516EnergyTypeFieldFresh, _ := json.Marshal(b516EnergyTypeFieldInput())
	b516EnergyTypeFieldNext, b516EnergyTypeFieldErr := B516EnergyTypeFieldEvaluate(b516EnergyTypeFieldFresh)
	if b516EnergyTypeFieldErr != nil || string(b516EnergyTypeFieldNext) != b516EnergyTypeFieldSaved {
		b516EnergyTypeFieldT.Fatal("result aliases later call")
	}
	if strings.Contains(string(b516EnergyTypeFieldNext), "meter") || strings.Contains(string(b516EnergyTypeFieldNext), "label") {
		b516EnergyTypeFieldT.Fatal("invented semantic association")
	}
}
