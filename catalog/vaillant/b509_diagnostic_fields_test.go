package vaillant

import (
	b509DiagnosticFieldBytes "bytes"
	b509DiagnosticFieldJSON "encoding/json"
	b509DiagnosticFieldOS "os"
	b509DiagnosticFieldReflect "reflect"
	b509DiagnosticFieldStrings "strings"
	b509DiagnosticFieldTesting "testing"
)

func b509DiagnosticFieldDecode(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T, b509DiagnosticFieldRaw []byte) map[string]any {
	b509DiagnosticFieldT.Helper()
	var b509DiagnosticFieldObject map[string]any
	b509DiagnosticFieldDecoder := b509DiagnosticFieldJSON.NewDecoder(b509DiagnosticFieldBytes.NewReader(b509DiagnosticFieldRaw))
	b509DiagnosticFieldDecoder.UseNumber()
	if b509DiagnosticFieldErr := b509DiagnosticFieldDecoder.Decode(&b509DiagnosticFieldObject); b509DiagnosticFieldErr != nil {
		b509DiagnosticFieldT.Fatal(b509DiagnosticFieldErr)
	}
	return b509DiagnosticFieldObject
}

func b509DiagnosticFieldRun(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T, b509DiagnosticFieldInput any) map[string]any {
	b509DiagnosticFieldT.Helper()
	b509DiagnosticFieldRaw, b509DiagnosticFieldErr := b509DiagnosticFieldJSON.Marshal(b509DiagnosticFieldInput)
	if b509DiagnosticFieldErr != nil {
		b509DiagnosticFieldT.Fatal(b509DiagnosticFieldErr)
	}
	b509DiagnosticFieldOutput, b509DiagnosticFieldErr := B509DiagnosticFieldEvaluate(b509DiagnosticFieldRaw)
	if b509DiagnosticFieldErr != nil {
		b509DiagnosticFieldT.Fatal(b509DiagnosticFieldErr)
	}
	return b509DiagnosticFieldDecode(b509DiagnosticFieldT, b509DiagnosticFieldOutput)
}

func b509DiagnosticFieldBase(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T) map[string]any {
	return b509DiagnosticFieldDecode(b509DiagnosticFieldT, []byte(`{"context":{"primary":181,"secondary":9,"operation":52,"direction":"request","view":"diagnostic-get","family":"synthetic-family","variant":"synthetic-variant","opaque":{"enabled":false,"items":[null,0] }},"fields":{"address":[120,86,52,18],"count":[255]},"arguments":{"requested_address":0,"extra":false}}`))
}

func TestAb17Ap17SyntheticCases(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T) {
	b509DiagnosticFieldRaw, b509DiagnosticFieldErr := b509DiagnosticFieldOS.ReadFile("testdata/ab17ap17_b509_diagnostic_fields.json")
	if b509DiagnosticFieldErr != nil {
		b509DiagnosticFieldT.Fatal(b509DiagnosticFieldErr)
	}
	var b509DiagnosticFieldCases []struct {
		B509DiagnosticFieldID       string                             `json:"id"`
		B509DiagnosticFieldInput    b509DiagnosticFieldJSON.RawMessage `json:"input"`
		B509DiagnosticFieldExpected b509DiagnosticFieldJSON.RawMessage `json:"expected"`
	}
	if b509DiagnosticFieldErr := b509DiagnosticFieldJSON.Unmarshal(b509DiagnosticFieldRaw, &b509DiagnosticFieldCases); b509DiagnosticFieldErr != nil {
		b509DiagnosticFieldT.Fatal(b509DiagnosticFieldErr)
	}
	if len(b509DiagnosticFieldCases) != 21 {
		b509DiagnosticFieldT.Fatalf("fixture count %d", len(b509DiagnosticFieldCases))
	}
	for _, b509DiagnosticFieldCase := range b509DiagnosticFieldCases {
		b509DiagnosticFieldT.Run(b509DiagnosticFieldCase.B509DiagnosticFieldID, func(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T) {
			b509DiagnosticFieldOutput, b509DiagnosticFieldErr := B509DiagnosticFieldEvaluate(b509DiagnosticFieldCase.B509DiagnosticFieldInput)
			if b509DiagnosticFieldErr != nil {
				b509DiagnosticFieldT.Fatal(b509DiagnosticFieldErr)
			}
			b509DiagnosticFieldActual := b509DiagnosticFieldDecode(b509DiagnosticFieldT, b509DiagnosticFieldOutput)
			b509DiagnosticFieldExpected := b509DiagnosticFieldDecode(b509DiagnosticFieldT, b509DiagnosticFieldCase.B509DiagnosticFieldExpected)
			if !b509DiagnosticFieldReflect.DeepEqual(b509DiagnosticFieldActual, b509DiagnosticFieldExpected) {
				b509DiagnosticFieldT.Fatalf("got %s; want %s", b509DiagnosticFieldOutput, b509DiagnosticFieldCase.B509DiagnosticFieldExpected)
			}
		})
	}
}

func TestAb17Ap17ContextAxes(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T) {
	for _, b509DiagnosticFieldKey := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant", "opaque"} {
		for _, b509DiagnosticFieldMode := range []string{"absent", "null", "changed", "zero"} {
			b509DiagnosticFieldT.Run("ab17ap17_"+b509DiagnosticFieldKey+"_"+b509DiagnosticFieldMode, func(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T) {
				b509DiagnosticFieldInput := b509DiagnosticFieldBase(b509DiagnosticFieldT)
				b509DiagnosticFieldContext := b509DiagnosticFieldInput["context"].(map[string]any)
				switch b509DiagnosticFieldMode {
				case "absent":
					delete(b509DiagnosticFieldContext, b509DiagnosticFieldKey)
				case "null":
					b509DiagnosticFieldContext[b509DiagnosticFieldKey] = nil
				default:
					switch b509DiagnosticFieldKey {
					case "primary", "secondary", "operation":
						if b509DiagnosticFieldMode == "zero" {
							b509DiagnosticFieldContext[b509DiagnosticFieldKey] = 0
						} else {
							b509DiagnosticFieldContext[b509DiagnosticFieldKey] = 53
						}
					case "direction", "view":
						b509DiagnosticFieldContext[b509DiagnosticFieldKey] = "unsupported"
					default:
						b509DiagnosticFieldContext[b509DiagnosticFieldKey] = false
					}
				}
				b509DiagnosticFieldOutput := b509DiagnosticFieldRun(b509DiagnosticFieldT, b509DiagnosticFieldInput)
				b509DiagnosticFieldWant := "Unknown"
				if b509DiagnosticFieldKey == "family" || b509DiagnosticFieldKey == "variant" || b509DiagnosticFieldKey == "opaque" {
					b509DiagnosticFieldWant = "complete"
				}
				if b509DiagnosticFieldOutput["state"] != b509DiagnosticFieldWant || b509DiagnosticFieldOutput["qualified"] != false || b509DiagnosticFieldOutput["physical"] != "Unknown" {
					b509DiagnosticFieldT.Fatalf("selector admitted or qualified: %v", b509DiagnosticFieldOutput)
				}
				b509DiagnosticFieldNormalized := b509DiagnosticFieldDecode(b509DiagnosticFieldT, b509DiagnosticFieldMustMarshal(b509DiagnosticFieldT, b509DiagnosticFieldInput))
				if !b509DiagnosticFieldReflect.DeepEqual(b509DiagnosticFieldOutput["context"], b509DiagnosticFieldNormalized["context"]) || !b509DiagnosticFieldReflect.DeepEqual(b509DiagnosticFieldOutput["raw"], b509DiagnosticFieldNormalized["fields"]) {
					b509DiagnosticFieldT.Fatal("context/raw changed")
				}
				if b509DiagnosticFieldWant == "Unknown" && len(b509DiagnosticFieldOutput["fields"].(map[string]any)) != 0 {
					b509DiagnosticFieldT.Fatal("unsupported formula applied")
				}
			})
		}
	}
}

func b509DiagnosticFieldMustMarshal(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T, b509DiagnosticFieldValue any) []byte {
	b509DiagnosticFieldT.Helper()
	b509DiagnosticFieldRaw, b509DiagnosticFieldErr := b509DiagnosticFieldJSON.Marshal(b509DiagnosticFieldValue)
	if b509DiagnosticFieldErr != nil {
		b509DiagnosticFieldT.Fatal(b509DiagnosticFieldErr)
	}
	return b509DiagnosticFieldRaw
}

func TestAb17Ap17PartialFields(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T) {
	for _, b509DiagnosticFieldKey := range []string{"address", "count", "status", "data"} {
		for _, b509DiagnosticFieldMode := range []string{"absent", "empty", "short", "overlong"} {
			b509DiagnosticFieldT.Run("ab17ap17_"+b509DiagnosticFieldKey+"_"+b509DiagnosticFieldMode, func(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T) {
				b509DiagnosticFieldInput := b509DiagnosticFieldBase(b509DiagnosticFieldT)
				b509DiagnosticFieldContext := b509DiagnosticFieldInput["context"].(map[string]any)
				b509DiagnosticFieldRaw := b509DiagnosticFieldInput["fields"].(map[string]any)
				b509DiagnosticFieldSibling := "address"
				if b509DiagnosticFieldKey == "address" {
					b509DiagnosticFieldSibling = "count"
				}
				if b509DiagnosticFieldKey == "data" {
					b509DiagnosticFieldContext["direction"] = "response"
					delete(b509DiagnosticFieldRaw, "count")
				}
				if b509DiagnosticFieldKey == "status" {
					b509DiagnosticFieldContext["operation"] = 53
					b509DiagnosticFieldContext["view"] = "diagnostic-put"
					b509DiagnosticFieldContext["direction"] = "response"
					b509DiagnosticFieldRaw["status"] = []int{6}
				}
				switch b509DiagnosticFieldMode {
				case "absent":
					delete(b509DiagnosticFieldRaw, b509DiagnosticFieldKey)
				case "empty":
					b509DiagnosticFieldRaw[b509DiagnosticFieldKey] = []int{}
				case "short":
					b509DiagnosticFieldRaw[b509DiagnosticFieldKey] = []int{1, 2, 3}
				case "overlong":
					b509DiagnosticFieldRaw[b509DiagnosticFieldKey] = []int{1, 2, 3, 4, 5}
				}
				b509DiagnosticFieldOutput := b509DiagnosticFieldRun(b509DiagnosticFieldT, b509DiagnosticFieldInput)
				b509DiagnosticFieldFields := b509DiagnosticFieldOutput["fields"].(map[string]any)
				b509DiagnosticFieldResult := b509DiagnosticFieldFields[b509DiagnosticFieldKey].(map[string]any)
				b509DiagnosticFieldWant := "incomplete"
				if b509DiagnosticFieldKey == "data" && b509DiagnosticFieldMode != "absent" {
					b509DiagnosticFieldWant = "complete"
				}
				if b509DiagnosticFieldResult["state"] != b509DiagnosticFieldWant {
					b509DiagnosticFieldT.Fatalf("wrong completeness: %v", b509DiagnosticFieldResult)
				}
				b509DiagnosticFieldTopWant := "partial"
				if b509DiagnosticFieldWant == "complete" {
					b509DiagnosticFieldTopWant = "complete"
				}
				if b509DiagnosticFieldOutput["state"] != b509DiagnosticFieldTopWant {
					b509DiagnosticFieldT.Fatal(b509DiagnosticFieldOutput)
				}
				if b509DiagnosticFieldWant == "incomplete" {
					if _, b509DiagnosticFieldExists := b509DiagnosticFieldResult["number"]; b509DiagnosticFieldExists {
						b509DiagnosticFieldT.Fatal("invented number")
					}
				}
				if b509DiagnosticFieldKey != "status" && b509DiagnosticFieldFields[b509DiagnosticFieldSibling].(map[string]any)["state"] != "complete" {
					b509DiagnosticFieldT.Fatal("independent sibling lost")
				}
				_, b509DiagnosticFieldHasRaw := b509DiagnosticFieldResult["raw"]
				if b509DiagnosticFieldHasRaw != (b509DiagnosticFieldMode != "absent") {
					b509DiagnosticFieldT.Fatal("missing raw fabricated")
				}
			})
		}
	}
}

func TestAb17Ap17RejectInvalid(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T) {
	b509DiagnosticFieldValid := string(b509DiagnosticFieldMustMarshal(b509DiagnosticFieldT, b509DiagnosticFieldBase(b509DiagnosticFieldT)))
	b509DiagnosticFieldCases := []string{"", "null", "[]", "{}", `{"context":{}}`, `{"fields":{}}`, "{", b509DiagnosticFieldValid + " {}", `{"context":{},"fields":{},"fields":{}}`, `{"context":{"nested":[{"x":1,"x":2}]},"fields":{}}`, `{"context":{},"fields":{},"arguments":{"x":false,"x":true}}`, `{"context":{},"fields":{"data":[],"data":[]}}`, `{"context":{},"fields":{},"extra":{"x":1,"x":2}}`, `{"context":{},"fields":{},"arguments":{"requested_address":0,"requested_\u0061ddress":1}}`}
	for _, b509DiagnosticFieldKey := range []string{"address", "count", "data", "status", "extra"} {
		for _, b509DiagnosticFieldBad := range []string{`null`, `false`, `0`, `"x"`, `{}`, `[-1]`, `[256]`, `[1.5]`, `[true]`, `[null]`, `["1"]`, `[[1]]`, `[1e0]`} {
			b509DiagnosticFieldCases = append(b509DiagnosticFieldCases, `{"context":{},"fields":{"`+b509DiagnosticFieldKey+`":`+b509DiagnosticFieldBad+`}}`)
		}
	}
	for _, b509DiagnosticFieldKey := range []string{"primary", "secondary", "operation"} {
		for _, b509DiagnosticFieldBad := range []string{`false`, `"181"`, `[]`, `{}`, `-1`, `256`, `181.5`} {
			b509DiagnosticFieldCases = append(b509DiagnosticFieldCases, `{"context":{"`+b509DiagnosticFieldKey+`":`+b509DiagnosticFieldBad+`},"fields":{}}`)
		}
	}
	for _, b509DiagnosticFieldKey := range []string{"direction", "view"} {
		for _, b509DiagnosticFieldBad := range []string{`false`, `0`, `[]`, `{}`} {
			b509DiagnosticFieldCases = append(b509DiagnosticFieldCases, `{"context":{"`+b509DiagnosticFieldKey+`":`+b509DiagnosticFieldBad+`},"fields":{}}`)
		}
	}
	for _, b509DiagnosticFieldBad := range []string{`null`, `false`, `[]`, `"1"`, `-1`, `4294967296`, `1.2`, `1e0`} {
		b509DiagnosticFieldCases = append(b509DiagnosticFieldCases, `{"context":{},"fields":{},"arguments":{"requested_address":`+b509DiagnosticFieldBad+`}}`)
	}
	for _, b509DiagnosticFieldKey := range []string{"context", "fields", "arguments"} {
		for _, b509DiagnosticFieldBad := range []string{`null`, `false`, `[]`, `0`, `"x"`} {
			b509DiagnosticFieldInput := b509DiagnosticFieldBase(b509DiagnosticFieldT)
			b509DiagnosticFieldInput[b509DiagnosticFieldKey] = b509DiagnosticFieldJSON.RawMessage(b509DiagnosticFieldBad)
			b509DiagnosticFieldCases = append(b509DiagnosticFieldCases, string(b509DiagnosticFieldMustMarshal(b509DiagnosticFieldT, b509DiagnosticFieldInput)))
		}
	}
	for b509DiagnosticFieldIndex, b509DiagnosticFieldCase := range b509DiagnosticFieldCases {
		b509DiagnosticFieldOutput, b509DiagnosticFieldErr := B509DiagnosticFieldEvaluate([]byte(b509DiagnosticFieldCase))
		if b509DiagnosticFieldErr == nil || b509DiagnosticFieldOutput != nil {
			b509DiagnosticFieldT.Fatalf("invalid case %d accepted: %s", b509DiagnosticFieldIndex, b509DiagnosticFieldCase)
		}
	}
}

func TestAb17Ap17ArgumentsAndOwnership(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T) {
	for _, b509DiagnosticFieldAddress := range []any{nil, 0, 1, 305419896, 4294967295} {
		b509DiagnosticFieldInput := b509DiagnosticFieldBase(b509DiagnosticFieldT)
		b509DiagnosticFieldInput["context"].(map[string]any)["direction"] = "response"
		b509DiagnosticFieldInput["fields"].(map[string]any)["data"] = []int{}
		b509DiagnosticFieldArgs := b509DiagnosticFieldInput["arguments"].(map[string]any)
		if b509DiagnosticFieldAddress == nil {
			delete(b509DiagnosticFieldArgs, "requested_address")
		} else {
			b509DiagnosticFieldArgs["requested_address"] = b509DiagnosticFieldAddress
		}
		b509DiagnosticFieldOutput := b509DiagnosticFieldRun(b509DiagnosticFieldT, b509DiagnosticFieldInput)
		b509DiagnosticFieldNormalized := b509DiagnosticFieldDecode(b509DiagnosticFieldT, b509DiagnosticFieldMustMarshal(b509DiagnosticFieldT, b509DiagnosticFieldInput))
		if !b509DiagnosticFieldReflect.DeepEqual(b509DiagnosticFieldOutput["arguments"], b509DiagnosticFieldNormalized["arguments"]) || b509DiagnosticFieldOutput["correlation"] != "Unknown" {
			b509DiagnosticFieldT.Fatal("argument used as echo constraint")
		}
	}
	for _, b509DiagnosticFieldMode := range []string{"absent", "empty", "false", "nested"} {
		b509DiagnosticFieldInput := b509DiagnosticFieldBase(b509DiagnosticFieldT)
		switch b509DiagnosticFieldMode {
		case "absent":
			delete(b509DiagnosticFieldInput, "arguments")
		case "empty":
			b509DiagnosticFieldInput["arguments"] = map[string]any{}
		case "false":
			b509DiagnosticFieldInput["arguments"].(map[string]any)["extra"] = false
		case "nested":
			b509DiagnosticFieldInput["arguments"].(map[string]any)["extra"] = map[string]any{"list": []any{nil, false, 0}}
		}
		b509DiagnosticFieldOutput := b509DiagnosticFieldRun(b509DiagnosticFieldT, b509DiagnosticFieldInput)
		b509DiagnosticFieldNormalized := b509DiagnosticFieldDecode(b509DiagnosticFieldT, b509DiagnosticFieldMustMarshal(b509DiagnosticFieldT, b509DiagnosticFieldInput))
		if !b509DiagnosticFieldReflect.DeepEqual(b509DiagnosticFieldOutput["arguments"], b509DiagnosticFieldNormalized["arguments"]) {
			b509DiagnosticFieldT.Fatal("arguments changed")
		}
		_, b509DiagnosticFieldExists := b509DiagnosticFieldOutput["arguments"]
		if b509DiagnosticFieldExists != (b509DiagnosticFieldMode != "absent") {
			b509DiagnosticFieldT.Fatal("arguments presence lost")
		}
	}
	b509DiagnosticFieldInput := b509DiagnosticFieldMustMarshal(b509DiagnosticFieldT, b509DiagnosticFieldBase(b509DiagnosticFieldT))
	b509DiagnosticFieldOriginal := append([]byte(nil), b509DiagnosticFieldInput...)
	b509DiagnosticFieldFirst, b509DiagnosticFieldErr := B509DiagnosticFieldEvaluate(b509DiagnosticFieldInput)
	if b509DiagnosticFieldErr != nil {
		b509DiagnosticFieldT.Fatal(b509DiagnosticFieldErr)
	}
	b509DiagnosticFieldSaved := append([]byte(nil), b509DiagnosticFieldFirst...)
	for b509DiagnosticFieldIndex := range b509DiagnosticFieldInput {
		b509DiagnosticFieldInput[b509DiagnosticFieldIndex] = ' '
	}
	if !b509DiagnosticFieldBytes.Equal(b509DiagnosticFieldSaved, b509DiagnosticFieldFirst) {
		b509DiagnosticFieldT.Fatal("input aliases result")
	}
	b509DiagnosticFieldDecoded := b509DiagnosticFieldDecode(b509DiagnosticFieldT, b509DiagnosticFieldFirst)
	b509DiagnosticFieldDecoded["raw"].(map[string]any)["address"].([]any)[0] = b509DiagnosticFieldJSON.Number("0")
	b509DiagnosticFieldDecoded["context"].(map[string]any)["opaque"].(map[string]any)["enabled"] = true
	b509DiagnosticFieldDecoded["arguments"].(map[string]any)["extra"] = true
	if b509DiagnosticFieldDecoded["fields"].(map[string]any)["address"].(map[string]any)["raw"].([]any)[0] != b509DiagnosticFieldJSON.Number("120") {
		b509DiagnosticFieldT.Fatal("nested raw aliases semantic field")
	}
	for b509DiagnosticFieldIndex := range b509DiagnosticFieldFirst {
		b509DiagnosticFieldFirst[b509DiagnosticFieldIndex] = ' '
	}
	b509DiagnosticFieldSecond, b509DiagnosticFieldErr := B509DiagnosticFieldEvaluate(b509DiagnosticFieldOriginal)
	if b509DiagnosticFieldErr != nil || !b509DiagnosticFieldBytes.Equal(b509DiagnosticFieldSaved, b509DiagnosticFieldSecond) {
		b509DiagnosticFieldT.Fatal("later call changed")
	}
	if b509DiagnosticFieldStrings.Contains(string(b509DiagnosticFieldSecond), `"label"`) {
		b509DiagnosticFieldT.Fatal("invented label")
	}
}

func TestAb17Ap17UnmappedFields(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T) {
	for _, b509DiagnosticFieldRaw := range [][]int{{}, {0}, {255}, {1, 2, 3}} {
		b509DiagnosticFieldInput := b509DiagnosticFieldBase(b509DiagnosticFieldT)
		b509DiagnosticFieldInput["fields"].(map[string]any)["extra"] = b509DiagnosticFieldRaw
		b509DiagnosticFieldOutput := b509DiagnosticFieldRun(b509DiagnosticFieldT, b509DiagnosticFieldInput)
		b509DiagnosticFieldExpected := b509DiagnosticFieldDecode(b509DiagnosticFieldT, b509DiagnosticFieldMustMarshal(b509DiagnosticFieldT, map[string]any{"raw": b509DiagnosticFieldRaw, "state": "Unknown"}))
		if !b509DiagnosticFieldReflect.DeepEqual(b509DiagnosticFieldOutput["fields"].(map[string]any)["extra"], b509DiagnosticFieldExpected) || b509DiagnosticFieldOutput["state"] != "complete" {
			b509DiagnosticFieldT.Fatal("unmapped field gained meaning or hid siblings")
		}
	}
}

func TestAb17Ap17RegistryFieldNames(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T) {
	b509DiagnosticFieldReq := b509DiagnosticFieldBase(b509DiagnosticFieldT)
	b509DiagnosticFieldFields := b509DiagnosticFieldReq["fields"].(map[string]any)
	b509DiagnosticFieldFields["requested_count"] = b509DiagnosticFieldFields["count"]
	delete(b509DiagnosticFieldFields, "count")
	b509DiagnosticFieldOut := b509DiagnosticFieldRun(b509DiagnosticFieldT, b509DiagnosticFieldReq)
	if b509DiagnosticFieldOut["state"] != "complete" {
		b509DiagnosticFieldT.Fatalf("requested_count not consumed: %v", b509DiagnosticFieldOut)
	}
	if b509DiagnosticFieldOut["fields"].(map[string]any)["requested_count"].(map[string]any)["state"] != "complete" {
		b509DiagnosticFieldT.Fatalf("requested_count unresolved: %v", b509DiagnosticFieldOut)
	}
	b509DiagnosticFieldPut := b509DiagnosticFieldBase(b509DiagnosticFieldT)
	b509DiagnosticFieldCtx := b509DiagnosticFieldPut["context"].(map[string]any)
	b509DiagnosticFieldCtx["operation"] = 53
	b509DiagnosticFieldCtx["direction"] = "response"
	b509DiagnosticFieldCtx["view"] = "diagnostic-put"
	b509DiagnosticFieldPut["fields"] = map[string]any{"write_status": []any{3}}
	b509DiagnosticFieldPutOut := b509DiagnosticFieldRun(b509DiagnosticFieldT, b509DiagnosticFieldPut)
	if b509DiagnosticFieldPutOut["state"] != "complete" || b509DiagnosticFieldPutOut["fields"].(map[string]any)["write_status"].(map[string]any)["candidate_domain"] != true {
		b509DiagnosticFieldT.Fatalf("write_status not consumed: %v", b509DiagnosticFieldPutOut)
	}
}

func TestAb17Ap17RegistryScalarBoundaries(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T) {
	for _, b509DiagnosticFieldKey := range []string{"requested_count", "write_status"} {
		for _, b509DiagnosticFieldValue := range []int{0, 6, 7, 255} {
			b509DiagnosticFieldT.Run("ab17ap17_"+b509DiagnosticFieldKey+"_"+string(b509DiagnosticFieldMustMarshal(b509DiagnosticFieldT, b509DiagnosticFieldValue)), func(b509DiagnosticFieldT *b509DiagnosticFieldTesting.T) {
				b509DiagnosticFieldInput := b509DiagnosticFieldBase(b509DiagnosticFieldT)
				b509DiagnosticFieldAlias := "count"
				b509DiagnosticFieldRaw := map[string]any{"address": []int{120, 86, 52, 18}, b509DiagnosticFieldKey: []int{b509DiagnosticFieldValue}}
				if b509DiagnosticFieldKey == "write_status" {
					b509DiagnosticFieldAlias = "status"
					b509DiagnosticFieldContext := b509DiagnosticFieldInput["context"].(map[string]any)
					b509DiagnosticFieldContext["operation"] = 53
					b509DiagnosticFieldContext["direction"] = "response"
					b509DiagnosticFieldContext["view"] = "diagnostic-put"
					delete(b509DiagnosticFieldRaw, "address")
				}
				b509DiagnosticFieldInput["fields"] = b509DiagnosticFieldRaw
				b509DiagnosticFieldOutput := b509DiagnosticFieldRun(b509DiagnosticFieldT, b509DiagnosticFieldInput)
				b509DiagnosticFieldFields := b509DiagnosticFieldOutput["fields"].(map[string]any)
				b509DiagnosticFieldExpected := map[string]any{
					"raw": []int{b509DiagnosticFieldValue}, "state": "complete", "grade": "Hypothesis",
					"number": b509DiagnosticFieldValue, "calculation": "Proven",
				}
				if b509DiagnosticFieldKey == "write_status" {
					b509DiagnosticFieldExpected["candidate_domain"] = b509DiagnosticFieldValue <= 6
				}
				b509DiagnosticFieldWant := b509DiagnosticFieldDecode(b509DiagnosticFieldT, b509DiagnosticFieldMustMarshal(b509DiagnosticFieldT, b509DiagnosticFieldExpected))
				if !b509DiagnosticFieldReflect.DeepEqual(b509DiagnosticFieldFields[b509DiagnosticFieldKey], b509DiagnosticFieldWant) {
					b509DiagnosticFieldT.Fatalf("canonical scalar differs: got %v; want %v", b509DiagnosticFieldFields[b509DiagnosticFieldKey], b509DiagnosticFieldWant)
				}
				if _, b509DiagnosticFieldExists := b509DiagnosticFieldFields[b509DiagnosticFieldAlias]; b509DiagnosticFieldExists {
					b509DiagnosticFieldT.Fatal("fabricated short alias")
				}
				if b509DiagnosticFieldOutput["state"] != "complete" || b509DiagnosticFieldOutput["physical"] != "Unknown" || b509DiagnosticFieldOutput["qualified"] != false || b509DiagnosticFieldOutput["correlation"] != "Unknown" {
					b509DiagnosticFieldT.Fatalf("canonical scalar changed evidence boundaries: %v", b509DiagnosticFieldOutput)
				}
				b509DiagnosticFieldNormalized := b509DiagnosticFieldDecode(b509DiagnosticFieldT, b509DiagnosticFieldMustMarshal(b509DiagnosticFieldT, b509DiagnosticFieldInput))
				if !b509DiagnosticFieldReflect.DeepEqual(b509DiagnosticFieldOutput["raw"], b509DiagnosticFieldNormalized["fields"]) {
					b509DiagnosticFieldT.Fatal("canonical raw fields changed")
				}
				if b509DiagnosticFieldKey == "requested_count" && b509DiagnosticFieldFields["address"].(map[string]any)["number"] != b509DiagnosticFieldJSON.Number("305419896") {
					b509DiagnosticFieldT.Fatal("independent address lost")
				}
			})
		}
	}
}
