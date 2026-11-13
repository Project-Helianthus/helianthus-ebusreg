package vaillant

import (
	b512FanFieldJSON "encoding/json"
	b512FanFieldFmt "fmt"
	b512FanFieldOs "os"
	b512FanFieldReflect "reflect"
	b512FanFieldStrings "strings"
	b512FanFieldTesting "testing"
)

func b512FanFieldSubset(b512FanFieldT *b512FanFieldTesting.T, b512FanFieldWant, b512FanFieldGot any) {
	b512FanFieldT.Helper()
	if b512FanFieldMap, b512FanFieldOK := b512FanFieldWant.(map[string]any); b512FanFieldOK {
		b512FanFieldActual, b512FanFieldOK := b512FanFieldGot.(map[string]any)
		if !b512FanFieldOK {
			b512FanFieldT.Fatalf("expected object, got %v", b512FanFieldGot)
		}
		for b512FanFieldKey, b512FanFieldValue := range b512FanFieldMap {
			b512FanFieldSubset(b512FanFieldT, b512FanFieldValue, b512FanFieldActual[b512FanFieldKey])
		}
		return
	}
	if !b512FanFieldReflect.DeepEqual(b512FanFieldWant, b512FanFieldGot) {
		b512FanFieldT.Fatalf("want %#v, got %#v", b512FanFieldWant, b512FanFieldGot)
	}
}

func b512FanFieldCall(b512FanFieldT *b512FanFieldTesting.T, b512FanFieldInput any) map[string]any {
	b512FanFieldT.Helper()
	b512FanFieldBytes, b512FanFieldErr := b512FanFieldJSON.Marshal(b512FanFieldInput)
	if b512FanFieldErr != nil {
		b512FanFieldT.Fatal(b512FanFieldErr)
	}
	b512FanFieldResult, b512FanFieldErr := B512FanFieldEvaluate(b512FanFieldBytes)
	if b512FanFieldErr != nil {
		b512FanFieldT.Fatal(b512FanFieldErr)
	}
	b512FanFieldOutput := map[string]any{}
	if b512FanFieldErr := b512FanFieldJSON.Unmarshal(b512FanFieldResult, &b512FanFieldOutput); b512FanFieldErr != nil {
		b512FanFieldT.Fatal(b512FanFieldErr)
	}
	return b512FanFieldOutput
}

func TestAb11Ap11Synthetic(b512FanFieldT *b512FanFieldTesting.T) {
	b512FanFieldData, b512FanFieldErr := b512FanFieldOs.ReadFile("testdata/ab11ap11_b512_fan_fields.json")
	if b512FanFieldErr != nil {
		b512FanFieldT.Fatal(b512FanFieldErr)
	}
	var b512FanFieldCases []struct {
		B512FanFieldID       string         `json:"id"`
		B512FanFieldInput    map[string]any `json:"input"`
		B512FanFieldExpected map[string]any `json:"expected"`
	}
	if b512FanFieldErr := b512FanFieldJSON.Unmarshal(b512FanFieldData, &b512FanFieldCases); b512FanFieldErr != nil {
		b512FanFieldT.Fatal(b512FanFieldErr)
	}
	if len(b512FanFieldCases) < 100 {
		b512FanFieldT.Fatal("missing independent cases")
	}
	for _, b512FanFieldCase := range b512FanFieldCases {
		b512FanFieldT.Run(b512FanFieldCase.B512FanFieldID, func(b512FanFieldT *b512FanFieldTesting.T) {
			b512FanFieldOutput := b512FanFieldCall(b512FanFieldT, b512FanFieldCase.B512FanFieldInput)
			b512FanFieldSubset(b512FanFieldT, b512FanFieldCase.B512FanFieldExpected, b512FanFieldOutput)
			if !b512FanFieldReflect.DeepEqual(b512FanFieldOutput["context"], b512FanFieldCase.B512FanFieldInput["context"]) || !b512FanFieldReflect.DeepEqual(b512FanFieldOutput["raw"], b512FanFieldCase.B512FanFieldInput["fields"]) {
				b512FanFieldT.Fatal("raw/context changed")
			}
			for b512FanFieldKey, b512FanFieldValue := range b512FanFieldOutput["fields"].(map[string]any) {
				b512FanFieldField := b512FanFieldValue.(map[string]any)
				if b512FanFieldField["state"] == "incomplete" || b512FanFieldField["state"] == "Unknown" {
					if _, b512FanFieldOK := b512FanFieldField["number"]; b512FanFieldOK {
						b512FanFieldT.Fatal("invented number")
					}
					if _, b512FanFieldOK := b512FanFieldField["label"]; b512FanFieldOK {
						b512FanFieldT.Fatal("invented label")
					}
				}
				b512FanFieldRaw, b512FanFieldExists := b512FanFieldCase.B512FanFieldInput["fields"].(map[string]any)[b512FanFieldKey]
				b512FanFieldActual, b512FanFieldHasRaw := b512FanFieldField["raw"]
				if b512FanFieldExists != b512FanFieldHasRaw || !b512FanFieldReflect.DeepEqual(b512FanFieldRaw, b512FanFieldActual) {
					b512FanFieldT.Fatal("field raw distinction lost")
				}
			}
			b512FanFieldControlInput := b512FanFieldCase.B512FanFieldInput
			if b512FanFieldOutput["state"] == "partial" {
				for b512FanFieldKey, b512FanFieldValue := range b512FanFieldOutput["fields"].(map[string]any) {
					if b512FanFieldValue.(map[string]any)["state"] == "incomplete" {
						continue
					}
					if b512FanFieldControlInput["fields"].(map[string]any)[b512FanFieldKey] != nil && b512FanFieldValue.(map[string]any)["state"] == nil {
						b512FanFieldT.Fatal("sibling lost")
					}
				}
				if len(b512FanFieldOutput["fields"].(map[string]any)) < len(b512FanFieldControlInput["fields"].(map[string]any)) {
					b512FanFieldT.Fatal("siblings dropped")
				}
			}
		})
	}
}

func b512FanFieldBase() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 18, "operation": 14, "direction": "positive", "view": "air-quality", "family": "synthetic", "variant": "a", "opaque": map[string]any{"false": false, "null": nil, "array": []any{0, "x"}}}, "fields": map[string]any{"temperature": []int{16, 0}, "humidity": []int{80, 6}, "air_quality": []int{255, 255}}}
}

func TestAb11Ap11ContextAndArguments(b512FanFieldT *b512FanFieldTesting.T) {
	for _, b512FanFieldKey := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
		for _, b512FanFieldValue := range []any{nil, 0, "other", "", false} {
			if (b512FanFieldKey == "primary" || b512FanFieldKey == "secondary" || b512FanFieldKey == "operation") && (b512FanFieldValue == "other" || b512FanFieldValue == "" || b512FanFieldValue == false) {
				continue
			}
			if (b512FanFieldKey == "direction" || b512FanFieldKey == "view") && (b512FanFieldValue == 0 || b512FanFieldValue == false) {
				continue
			}
			b512FanFieldT.Run(b512FanFieldKey+b512FanFieldFmt.Sprint(b512FanFieldValue), func(b512FanFieldT *b512FanFieldTesting.T) {
				b512FanFieldInput := b512FanFieldBase()
				b512FanFieldInput["context"].(map[string]any)[b512FanFieldKey] = b512FanFieldValue
				b512FanFieldOutput := b512FanFieldCall(b512FanFieldT, b512FanFieldInput)
				b512FanFieldBytes, _ := b512FanFieldJSON.Marshal(b512FanFieldInput["context"])
				var b512FanFieldContext any
				_ = b512FanFieldJSON.Unmarshal(b512FanFieldBytes, &b512FanFieldContext)
				if !b512FanFieldReflect.DeepEqual(b512FanFieldContext, b512FanFieldOutput["context"]) {
					b512FanFieldT.Fatal("context changed")
				}
				if b512FanFieldKey != "family" && b512FanFieldKey != "variant" {
					if b512FanFieldOutput["state"] != "Unknown" {
						b512FanFieldT.Fatal("unsupported admitted")
					}
					for _, b512FanFieldField := range b512FanFieldOutput["fields"].(map[string]any) {
						if len(b512FanFieldField.(map[string]any)) != 2 {
							b512FanFieldT.Fatal("unsupported formula")
						}
					}
				} else if b512FanFieldOutput["state"] != "complete" {
					b512FanFieldT.Fatal("opaque selector used")
				}
				if b512FanFieldOutput["qualified"] != false || b512FanFieldOutput["physical"] != "Unknown" {
					b512FanFieldT.Fatal("qualified")
				}
			})
		}
		b512FanFieldInput := b512FanFieldBase()
		delete(b512FanFieldInput["context"].(map[string]any), b512FanFieldKey)
		b512FanFieldOutput := b512FanFieldCall(b512FanFieldT, b512FanFieldInput)
		if _, b512FanFieldExists := b512FanFieldOutput["context"].(map[string]any)[b512FanFieldKey]; b512FanFieldExists {
			b512FanFieldT.Fatal("missing context filled")
		}
	}
	for _, b512FanFieldArguments := range []any{map[string]any{}, map[string]any{"opaque": false}, map[string]any{"opaque": 0}, map[string]any{"opaque": nil}, map[string]any{"opaque": []any{false, 0, nil, "x"}}} {
		b512FanFieldInput := b512FanFieldBase()
		b512FanFieldInput["arguments"] = b512FanFieldArguments
		b512FanFieldOutput := b512FanFieldCall(b512FanFieldT, b512FanFieldInput)
		b512FanFieldBytes, _ := b512FanFieldJSON.Marshal(b512FanFieldArguments)
		var b512FanFieldWant any
		_ = b512FanFieldJSON.Unmarshal(b512FanFieldBytes, &b512FanFieldWant)
		if !b512FanFieldReflect.DeepEqual(b512FanFieldWant, b512FanFieldOutput["arguments"]) {
			b512FanFieldT.Fatal("arguments changed")
		}
	}
	if _, b512FanFieldExists := b512FanFieldCall(b512FanFieldT, b512FanFieldBase())["arguments"]; b512FanFieldExists {
		b512FanFieldT.Fatal("absent arguments filled")
	}
}

func TestAb11Ap11Reject(b512FanFieldT *b512FanFieldTesting.T) {
	for b512FanFieldIndex, b512FanFieldInput := range []string{``, `{`, `[]`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"x":{"a":0,"a":1}},"fields":{}}`, `{"context":{},"fields":{"x":[{"a":0,"a":1}]}}`, `{"context":{},"fields":{},"arguments":{"x":[{"a":0,"a":1}]}}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":null}`, `{"context":{},"fields":{},"arguments":false}`} {
		b512FanFieldT.Run(b512FanFieldFmt.Sprint(b512FanFieldIndex), func(b512FanFieldT *b512FanFieldTesting.T) {
			b512FanFieldOutput, b512FanFieldErr := B512FanFieldEvaluate([]byte(b512FanFieldInput))
			if b512FanFieldErr == nil || b512FanFieldOutput != nil {
				b512FanFieldT.Fatal("invalid accepted")
			}
		})
	}
	for _, b512FanFieldBad := range []string{`true`, `false`, `null`, `"0"`, `-1`, `256`, `1.5`, `1e0`, `{}`, `[]`} {
		for _, b512FanFieldTemplate := range []string{`{"context":{},"fields":{"x":[VALUE]}}`, `{"context":{"primary":VALUE},"fields":{}}`, `{"context":{"secondary":VALUE},"fields":{}}`, `{"context":{"operation":VALUE},"fields":{}}`} {
			if b512FanFieldBad == `null` && !b512FanFieldStrings.Contains(b512FanFieldTemplate, `[VALUE]`) {
				continue
			}
			b512FanFieldOutput, b512FanFieldErr := B512FanFieldEvaluate([]byte(b512FanFieldStrings.ReplaceAll(b512FanFieldTemplate, "VALUE", b512FanFieldBad)))
			if b512FanFieldErr == nil || b512FanFieldOutput != nil {
				b512FanFieldT.Fatalf("accepted %s %s", b512FanFieldTemplate, b512FanFieldBad)
			}
		}
	}
	for _, b512FanFieldBad := range []any{nil, false, 0, "x", map[string]any{}} {
		b512FanFieldInput := b512FanFieldBase()
		b512FanFieldInput["fields"].(map[string]any)["humidity"] = b512FanFieldBad
		b512FanFieldBytes, _ := b512FanFieldJSON.Marshal(b512FanFieldInput)
		if b512FanFieldOutput, b512FanFieldErr := B512FanFieldEvaluate(b512FanFieldBytes); b512FanFieldErr == nil || b512FanFieldOutput != nil {
			b512FanFieldT.Fatal("non-array accepted")
		}
	}
}

func TestAb11Ap11Isolation(b512FanFieldT *b512FanFieldTesting.T) {
	b512FanFieldInput, _ := b512FanFieldJSON.Marshal(b512FanFieldBase())
	b512FanFieldOriginal := append([]byte(nil), b512FanFieldInput...)
	b512FanFieldFirst, b512FanFieldErr := B512FanFieldEvaluate(b512FanFieldInput)
	if b512FanFieldErr != nil {
		b512FanFieldT.Fatal(b512FanFieldErr)
	}
	b512FanFieldSaved := append([]byte(nil), b512FanFieldFirst...)
	for b512FanFieldIndex := range b512FanFieldInput {
		b512FanFieldInput[b512FanFieldIndex] = ' '
	}
	if !b512FanFieldReflect.DeepEqual(b512FanFieldFirst, b512FanFieldJSON.RawMessage(b512FanFieldSaved)) {
		b512FanFieldT.Fatal("input alias")
	}
	var b512FanFieldDecoded map[string]any
	_ = b512FanFieldJSON.Unmarshal(b512FanFieldFirst, &b512FanFieldDecoded)
	b512FanFieldDecoded["raw"].(map[string]any)["humidity"].([]any)[0] = 0.0
	if b512FanFieldDecoded["fields"].(map[string]any)["humidity"].(map[string]any)["raw"].([]any)[0] != float64(80) {
		b512FanFieldT.Fatal("nested alias")
	}
	for b512FanFieldIndex := range b512FanFieldFirst {
		b512FanFieldFirst[b512FanFieldIndex] = ' '
	}
	b512FanFieldSecond, b512FanFieldErr := B512FanFieldEvaluate(b512FanFieldOriginal)
	if b512FanFieldErr != nil || !b512FanFieldReflect.DeepEqual([]byte(b512FanFieldSecond), b512FanFieldSaved) {
		b512FanFieldT.Fatal("later result affected")
	}
}

func TestAb11Ap11RequestIsolation(b512FanFieldT *b512FanFieldTesting.T) {
	for _, b512FanFieldView := range []string{"digital-air", "air-quality"} {
		b512FanFieldInput := b512FanFieldBase()
		b512FanFieldContext := b512FanFieldInput["context"].(map[string]any)
		b512FanFieldContext["view"] = b512FanFieldView
		b512FanFieldContext["direction"] = "request"
		if b512FanFieldView == "digital-air" {
			b512FanFieldContext["operation"] = 13
		}
		b512FanFieldOutput := b512FanFieldCall(b512FanFieldT, b512FanFieldInput)
		if b512FanFieldOutput["state"] != "complete" {
			b512FanFieldT.Fatal("request incomplete")
		}
		for _, b512FanFieldField := range b512FanFieldOutput["fields"].(map[string]any) {
			if len(b512FanFieldField.(map[string]any)) != 2 {
				b512FanFieldT.Fatal("positive semantics on request")
			}
		}
	}
	for _, b512FanFieldKey := range []string{"direction", "view"} {
		for _, b512FanFieldBad := range []any{false, 0, []any{}, map[string]any{}} {
			b512FanFieldInput := b512FanFieldBase()
			b512FanFieldInput["context"].(map[string]any)[b512FanFieldKey] = b512FanFieldBad
			b512FanFieldBytes, _ := b512FanFieldJSON.Marshal(b512FanFieldInput)
			if b512FanFieldResult, b512FanFieldErr := B512FanFieldEvaluate(b512FanFieldBytes); b512FanFieldErr == nil || b512FanFieldResult != nil {
				b512FanFieldT.Fatal("wrong context kind admitted")
			}
		}
	}
}
