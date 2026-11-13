package vaillant

import (
	b51ADiagnosticFieldJSON "encoding/json"
	b51ADiagnosticFieldOs "os"
	b51ADiagnosticFieldReflect "reflect"
	b51ADiagnosticFieldStrings "strings"
	b51ADiagnosticFieldTesting "testing"
)

type b51ADiagnosticFieldCase struct {
	B51ADiagnosticFieldID       string         `json:"id"`
	B51ADiagnosticFieldInput    map[string]any `json:"input"`
	B51ADiagnosticFieldExpected map[string]any `json:"expected"`
}

func b51ADiagnosticFieldCases(b51ADiagnosticFieldT *b51ADiagnosticFieldTesting.T) []b51ADiagnosticFieldCase {
	b51ADiagnosticFieldT.Helper()
	b51ADiagnosticFieldData, b51ADiagnosticFieldErr := b51ADiagnosticFieldOs.ReadFile("testdata/ab21aq21_b51a_diagnostic_fields.json")
	if b51ADiagnosticFieldErr != nil {
		b51ADiagnosticFieldT.Fatal(b51ADiagnosticFieldErr)
	}
	var b51ADiagnosticFieldCases []b51ADiagnosticFieldCase
	if b51ADiagnosticFieldErr := b51ADiagnosticFieldJSON.Unmarshal(b51ADiagnosticFieldData, &b51ADiagnosticFieldCases); b51ADiagnosticFieldErr != nil {
		b51ADiagnosticFieldT.Fatal(b51ADiagnosticFieldErr)
	}
	if len(b51ADiagnosticFieldCases) != 9 {
		b51ADiagnosticFieldT.Fatal("fixture count", len(b51ADiagnosticFieldCases))
	}
	return b51ADiagnosticFieldCases
}

func b51ADiagnosticFieldEval(b51ADiagnosticFieldT *b51ADiagnosticFieldTesting.T, b51ADiagnosticFieldInput map[string]any) map[string]any {
	b51ADiagnosticFieldT.Helper()
	b51ADiagnosticFieldData, b51ADiagnosticFieldErr := b51ADiagnosticFieldJSON.Marshal(b51ADiagnosticFieldInput)
	if b51ADiagnosticFieldErr != nil {
		b51ADiagnosticFieldT.Fatal(b51ADiagnosticFieldErr)
	}
	b51ADiagnosticFieldOutput, b51ADiagnosticFieldErr := B51ADiagnosticFieldEvaluate(b51ADiagnosticFieldData)
	if b51ADiagnosticFieldErr != nil {
		b51ADiagnosticFieldT.Fatal(b51ADiagnosticFieldErr)
	}
	var b51ADiagnosticFieldResult map[string]any
	if b51ADiagnosticFieldErr := b51ADiagnosticFieldJSON.Unmarshal(b51ADiagnosticFieldOutput, &b51ADiagnosticFieldResult); b51ADiagnosticFieldErr != nil {
		b51ADiagnosticFieldT.Fatal(b51ADiagnosticFieldErr)
	}
	if b51ADiagnosticFieldResult["qualified"] != false || b51ADiagnosticFieldResult["physical"] != "Unknown" {
		b51ADiagnosticFieldT.Fatal("qualification", b51ADiagnosticFieldResult)
	}
	if !b51ADiagnosticFieldReflect.DeepEqual(b51ADiagnosticFieldInput, b51ADiagnosticFieldResult["input"]) {
		b51ADiagnosticFieldT.Fatal("input preservation")
	}
	if !b51ADiagnosticFieldReflect.DeepEqual(b51ADiagnosticFieldInput["fields"], b51ADiagnosticFieldResult["raw"]) {
		b51ADiagnosticFieldT.Fatal("raw preservation")
	}
	return b51ADiagnosticFieldResult
}

func b51ADiagnosticFieldSubset(b51ADiagnosticFieldT *b51ADiagnosticFieldTesting.T, b51ADiagnosticFieldGot, b51ADiagnosticFieldWant any) {
	b51ADiagnosticFieldT.Helper()
	if b51ADiagnosticFieldMap, b51ADiagnosticFieldOK := b51ADiagnosticFieldWant.(map[string]any); b51ADiagnosticFieldOK {
		b51ADiagnosticFieldActual, b51ADiagnosticFieldOK := b51ADiagnosticFieldGot.(map[string]any)
		if !b51ADiagnosticFieldOK {
			b51ADiagnosticFieldT.Fatalf("want object, got %v", b51ADiagnosticFieldGot)
		}
		for b51ADiagnosticFieldKey, b51ADiagnosticFieldValue := range b51ADiagnosticFieldMap {
			b51ADiagnosticFieldSubset(b51ADiagnosticFieldT, b51ADiagnosticFieldActual[b51ADiagnosticFieldKey], b51ADiagnosticFieldValue)
		}
		return
	}
	if !b51ADiagnosticFieldReflect.DeepEqual(b51ADiagnosticFieldGot, b51ADiagnosticFieldWant) {
		b51ADiagnosticFieldT.Fatalf("got %v want %v", b51ADiagnosticFieldGot, b51ADiagnosticFieldWant)
	}
}

func TestAb21Aq21Fixtures(b51ADiagnosticFieldT *b51ADiagnosticFieldTesting.T) {
	for _, b51ADiagnosticFieldCase := range b51ADiagnosticFieldCases(b51ADiagnosticFieldT) {
		b51ADiagnosticFieldT.Run(b51ADiagnosticFieldCase.B51ADiagnosticFieldID, func(b51ADiagnosticFieldT *b51ADiagnosticFieldTesting.T) {
			b51ADiagnosticFieldOutput := b51ADiagnosticFieldEval(b51ADiagnosticFieldT, b51ADiagnosticFieldCase.B51ADiagnosticFieldInput)
			if b51ADiagnosticFieldOutput["state"] != "complete" {
				b51ADiagnosticFieldT.Fatal(b51ADiagnosticFieldOutput)
			}
			b51ADiagnosticFieldSubset(b51ADiagnosticFieldT, b51ADiagnosticFieldOutput["results"], b51ADiagnosticFieldCase.B51ADiagnosticFieldExpected)
			for b51ADiagnosticFieldKey, b51ADiagnosticFieldValue := range b51ADiagnosticFieldOutput["results"].(map[string]any) {
				b51ADiagnosticFieldResult := b51ADiagnosticFieldValue.(map[string]any)
				if !b51ADiagnosticFieldReflect.DeepEqual(b51ADiagnosticFieldResult["raw"], b51ADiagnosticFieldCase.B51ADiagnosticFieldInput["fields"].(map[string]any)[b51ADiagnosticFieldKey]) {
					b51ADiagnosticFieldT.Fatal("result raw")
				}
				if b51ADiagnosticFieldKey == "alternatives" {
					if _, b51ADiagnosticFieldPresent := b51ADiagnosticFieldResult["number"]; b51ADiagnosticFieldPresent {
						b51ADiagnosticFieldT.Fatal("enum current number")
					}
				}
				if b51ADiagnosticFieldResult["state"] == "Unknown" {
					if _, b51ADiagnosticFieldPresent := b51ADiagnosticFieldResult["number"]; b51ADiagnosticFieldPresent {
						b51ADiagnosticFieldT.Fatal("disabled interpreted", b51ADiagnosticFieldKey)
					}
					continue
				}
				if b51ADiagnosticFieldResult["state"] != "Hypothesis" {
					b51ADiagnosticFieldT.Fatal("candidate grade", b51ADiagnosticFieldResult)
				}
				if _, b51ADiagnosticFieldPresent := b51ADiagnosticFieldResult["number"]; b51ADiagnosticFieldPresent && b51ADiagnosticFieldResult["calculation"] != "Proven" {
					b51ADiagnosticFieldT.Fatal("calculation grade")
				}
				for _, b51ADiagnosticFieldForbidden := range []string{"current", "current_value", "success", "effective_authorization"} {
					if _, b51ADiagnosticFieldPresent := b51ADiagnosticFieldResult[b51ADiagnosticFieldForbidden]; b51ADiagnosticFieldPresent {
						b51ADiagnosticFieldT.Fatal("invented result", b51ADiagnosticFieldForbidden)
					}
				}
			}
		})
	}
}

func TestAb21Aq21EnablePositions(b51ADiagnosticFieldT *b51ADiagnosticFieldTesting.T) {
	for b51ADiagnosticFieldPosition := 0; b51ADiagnosticFieldPosition < 104; b51ADiagnosticFieldPosition++ {
		b51ADiagnosticFieldInput := b51ADiagnosticFieldCases(b51ADiagnosticFieldT)[0].B51ADiagnosticFieldInput
		b51ADiagnosticFieldInput["fields"].(map[string]any)["enable"].([]any)[b51ADiagnosticFieldPosition/8] = float64(uint64(1) << uint(b51ADiagnosticFieldPosition%8))
		b51ADiagnosticFieldBits := b51ADiagnosticFieldEval(b51ADiagnosticFieldT, b51ADiagnosticFieldInput)["results"].(map[string]any)["enable"].(map[string]any)["flags"].([]any)
		if len(b51ADiagnosticFieldBits) != 104 {
			b51ADiagnosticFieldT.Fatal("bit count")
		}
		for b51ADiagnosticFieldIndex, b51ADiagnosticFieldBit := range b51ADiagnosticFieldBits {
			if b51ADiagnosticFieldBit != (b51ADiagnosticFieldIndex == b51ADiagnosticFieldPosition) {
				b51ADiagnosticFieldT.Fatal("bit position", b51ADiagnosticFieldPosition, b51ADiagnosticFieldIndex)
			}
		}
	}
}

func TestAb21Aq21ContextAxes(b51ADiagnosticFieldT *b51ADiagnosticFieldTesting.T) {
	for _, b51ADiagnosticFieldKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, b51ADiagnosticFieldAbsent := range []bool{false, true} {
			b51ADiagnosticFieldInput := b51ADiagnosticFieldCases(b51ADiagnosticFieldT)[1].B51ADiagnosticFieldInput
			b51ADiagnosticFieldContext := b51ADiagnosticFieldInput["context"].(map[string]any)
			if b51ADiagnosticFieldKey == "direction" || b51ADiagnosticFieldKey == "view" {
				b51ADiagnosticFieldContext[b51ADiagnosticFieldKey] = "unsupported"
			} else {
				b51ADiagnosticFieldContext[b51ADiagnosticFieldKey] = float64(0)
			}
			if b51ADiagnosticFieldAbsent {
				delete(b51ADiagnosticFieldContext, b51ADiagnosticFieldKey)
			}
			b51ADiagnosticFieldOutput := b51ADiagnosticFieldEval(b51ADiagnosticFieldT, b51ADiagnosticFieldInput)
			if b51ADiagnosticFieldOutput["state"] != "Unknown" || len(b51ADiagnosticFieldOutput["results"].(map[string]any)) != 0 {
				b51ADiagnosticFieldT.Fatal("context admitted", b51ADiagnosticFieldKey, b51ADiagnosticFieldAbsent)
			}
		}
	}
	for _, b51ADiagnosticFieldKey := range []string{"family", "variant"} {
		for _, b51ADiagnosticFieldValue := range []any{nil, false, float64(0), "", "another-family", []any{}} {
			b51ADiagnosticFieldInput := b51ADiagnosticFieldCases(b51ADiagnosticFieldT)[1].B51ADiagnosticFieldInput
			b51ADiagnosticFieldInput["context"].(map[string]any)[b51ADiagnosticFieldKey] = b51ADiagnosticFieldValue
			if b51ADiagnosticFieldEval(b51ADiagnosticFieldT, b51ADiagnosticFieldInput)["state"] != "complete" {
				b51ADiagnosticFieldT.Fatal("opaque identity selects layout")
			}
		}
	}
}

func TestAb21Aq21FieldWidthsAndSiblings(b51ADiagnosticFieldT *b51ADiagnosticFieldTesting.T) {
	for _, b51ADiagnosticFieldCase := range b51ADiagnosticFieldCases(b51ADiagnosticFieldT) {
		for b51ADiagnosticFieldKey, b51ADiagnosticFieldValue := range b51ADiagnosticFieldCase.B51ADiagnosticFieldInput["fields"].(map[string]any) {
			for _, b51ADiagnosticFieldMode := range []string{"absent", "empty", "short", "long"} {
				if b51ADiagnosticFieldMode == "absent" && b51ADiagnosticFieldCase.B51ADiagnosticFieldID == "ab21aq21_disabled" && b51ADiagnosticFieldKey != "flags" {
					continue
				}
				b51ADiagnosticFieldInput := b51ADiagnosticFieldCases(b51ADiagnosticFieldT)[0].B51ADiagnosticFieldInput
				for _, b51ADiagnosticFieldCandidate := range b51ADiagnosticFieldCases(b51ADiagnosticFieldT) {
					if b51ADiagnosticFieldCandidate.B51ADiagnosticFieldID == b51ADiagnosticFieldCase.B51ADiagnosticFieldID {
						b51ADiagnosticFieldInput = b51ADiagnosticFieldCandidate.B51ADiagnosticFieldInput
					}
				}
				b51ADiagnosticFieldFields := b51ADiagnosticFieldInput["fields"].(map[string]any)
				b51ADiagnosticFieldArray := b51ADiagnosticFieldValue.([]any)
				switch b51ADiagnosticFieldMode {
				case "absent":
					delete(b51ADiagnosticFieldFields, b51ADiagnosticFieldKey)
				case "empty":
					b51ADiagnosticFieldFields[b51ADiagnosticFieldKey] = []any{}
				case "short":
					b51ADiagnosticFieldFields[b51ADiagnosticFieldKey] = b51ADiagnosticFieldArray[:len(b51ADiagnosticFieldArray)-1]
				case "long":
					b51ADiagnosticFieldFields[b51ADiagnosticFieldKey] = append(append([]any{}, b51ADiagnosticFieldArray...), float64(0))
				}
				b51ADiagnosticFieldOutput := b51ADiagnosticFieldEval(b51ADiagnosticFieldT, b51ADiagnosticFieldInput)
				b51ADiagnosticFieldResult := b51ADiagnosticFieldOutput["results"].(map[string]any)[b51ADiagnosticFieldKey].(map[string]any)
				if b51ADiagnosticFieldOutput["state"] != "partial" || b51ADiagnosticFieldResult["state"] != "incomplete" {
					b51ADiagnosticFieldT.Fatal("width", b51ADiagnosticFieldKey, b51ADiagnosticFieldMode)
				}
				if _, b51ADiagnosticFieldPresent := b51ADiagnosticFieldResult["number"]; b51ADiagnosticFieldPresent {
					b51ADiagnosticFieldT.Fatal("incomplete number")
				}
				if _, b51ADiagnosticFieldPresent := b51ADiagnosticFieldResult["raw"]; b51ADiagnosticFieldPresent != (b51ADiagnosticFieldMode != "absent") {
					b51ADiagnosticFieldT.Fatal("absence conflated")
				}
				if b51ADiagnosticFieldCase.B51ADiagnosticFieldID == "ab21aq21_range" && b51ADiagnosticFieldKey != "flags" && b51ADiagnosticFieldKey != "value" {
					b51ADiagnosticFieldSubset(b51ADiagnosticFieldT, b51ADiagnosticFieldOutput["results"].(map[string]any)["value"], map[string]any{"number": float64(4660), "state": "Hypothesis"})
				}
			}
		}
	}
}

func TestAb21Aq21RejectJSONAndKinds(b51ADiagnosticFieldT *b51ADiagnosticFieldTesting.T) {
	for _, b51ADiagnosticFieldInput := range []string{``, `{`, `[]`, `{} {}`, `{"context":{},"fields":{},"x":{"y":1,"y":2}}`, `{"context":{"primary":181,"primary":181},"fields":{}}`, `{"context":{},"fields":{"x":[{"a":0,"a":1}]}}`, `{"context":{},"fields":{},"x":1e999}`, `{"context":{},"fields":{},"x":NaN}`} {
		b51ADiagnosticFieldOutput, b51ADiagnosticFieldErr := B51ADiagnosticFieldEvaluate([]byte(b51ADiagnosticFieldInput))
		if b51ADiagnosticFieldErr == nil || b51ADiagnosticFieldOutput != nil {
			b51ADiagnosticFieldT.Fatal("accepted", b51ADiagnosticFieldInput)
		}
	}
	for _, b51ADiagnosticFieldKey := range []string{"primary", "secondary", "operation", "byte", "value_width", "context", "fields", "arguments", "direction", "view"} {
		for _, b51ADiagnosticFieldBad := range []any{nil, true, "1", float64(-1), float64(256), 1.5, []any{}} {
			if (b51ADiagnosticFieldKey == "direction" || b51ADiagnosticFieldKey == "view") && b51ADiagnosticFieldBad == "1" {
				continue
			}
			b51ADiagnosticFieldInput := b51ADiagnosticFieldCases(b51ADiagnosticFieldT)[5].B51ADiagnosticFieldInput
			switch b51ADiagnosticFieldKey {
			case "byte":
				b51ADiagnosticFieldInput["fields"].(map[string]any)["value"] = []any{b51ADiagnosticFieldBad, 0}
			case "value_width":
				b51ADiagnosticFieldInput["arguments"].(map[string]any)[b51ADiagnosticFieldKey] = b51ADiagnosticFieldBad
			case "context", "fields", "arguments":
				b51ADiagnosticFieldInput[b51ADiagnosticFieldKey] = b51ADiagnosticFieldBad
			default:
				b51ADiagnosticFieldInput["context"].(map[string]any)[b51ADiagnosticFieldKey] = b51ADiagnosticFieldBad
			}
			b51ADiagnosticFieldData, _ := b51ADiagnosticFieldJSON.Marshal(b51ADiagnosticFieldInput)
			b51ADiagnosticFieldOutput, b51ADiagnosticFieldErr := B51ADiagnosticFieldEvaluate(b51ADiagnosticFieldData)
			if b51ADiagnosticFieldErr == nil || b51ADiagnosticFieldOutput != nil {
				b51ADiagnosticFieldT.Fatal("accepted kind", b51ADiagnosticFieldKey, b51ADiagnosticFieldBad)
			}
		}
	}
	for _, b51ADiagnosticFieldWidth := range []any{nil, float64(0), float64(3)} {
		b51ADiagnosticFieldInput := b51ADiagnosticFieldCases(b51ADiagnosticFieldT)[5].B51ADiagnosticFieldInput
		if b51ADiagnosticFieldWidth == nil {
			delete(b51ADiagnosticFieldInput, "arguments")
		} else {
			b51ADiagnosticFieldInput["arguments"].(map[string]any)["value_width"] = b51ADiagnosticFieldWidth
		}
		b51ADiagnosticFieldData, _ := b51ADiagnosticFieldJSON.Marshal(b51ADiagnosticFieldInput)
		if b51ADiagnosticFieldOutput, b51ADiagnosticFieldErr := B51ADiagnosticFieldEvaluate(b51ADiagnosticFieldData); b51ADiagnosticFieldErr == nil || b51ADiagnosticFieldOutput != nil {
			b51ADiagnosticFieldT.Fatal("inferred width")
		}
	}
}

func TestAb21Aq21PreservationAndMemory(b51ADiagnosticFieldT *b51ADiagnosticFieldTesting.T) {
	b51ADiagnosticFieldInput := b51ADiagnosticFieldCases(b51ADiagnosticFieldT)[1].B51ADiagnosticFieldInput
	for _, b51ADiagnosticFieldKey := range []string{"state", "results", "physical", "qualified", "input", "raw", "extra"} {
		b51ADiagnosticFieldInput[b51ADiagnosticFieldKey] = map[string]any{"null": nil, "zero": float64(0), "false": false, "empty": []any{}}
	}
	b51ADiagnosticFieldInput["fields"].(map[string]any)["opaque"] = []any{}
	b51ADiagnosticFieldOutput := b51ADiagnosticFieldEval(b51ADiagnosticFieldT, b51ADiagnosticFieldInput)
	if !b51ADiagnosticFieldReflect.DeepEqual(b51ADiagnosticFieldOutput["extra"], b51ADiagnosticFieldInput["extra"]) {
		b51ADiagnosticFieldT.Fatal("extra lost")
	}
	b51ADiagnosticFieldData, _ := b51ADiagnosticFieldJSON.Marshal(b51ADiagnosticFieldInput)
	b51ADiagnosticFieldFirst, b51ADiagnosticFieldErr := B51ADiagnosticFieldEvaluate(b51ADiagnosticFieldData)
	if b51ADiagnosticFieldErr != nil {
		b51ADiagnosticFieldT.Fatal(b51ADiagnosticFieldErr)
	}
	b51ADiagnosticFieldSaved := string(b51ADiagnosticFieldFirst)
	for b51ADiagnosticFieldIndex := range b51ADiagnosticFieldData {
		b51ADiagnosticFieldData[b51ADiagnosticFieldIndex] = ' '
	}
	if string(b51ADiagnosticFieldFirst) != b51ADiagnosticFieldSaved {
		b51ADiagnosticFieldT.Fatal("input aliases output")
	}
	for b51ADiagnosticFieldIndex := range b51ADiagnosticFieldFirst {
		b51ADiagnosticFieldFirst[b51ADiagnosticFieldIndex] = ' '
	}
	b51ADiagnosticFieldNext, _ := b51ADiagnosticFieldJSON.Marshal(b51ADiagnosticFieldEval(b51ADiagnosticFieldT, b51ADiagnosticFieldInput))
	var b51ADiagnosticFieldExpected any
	if b51ADiagnosticFieldErr := b51ADiagnosticFieldJSON.Unmarshal([]byte(b51ADiagnosticFieldSaved), &b51ADiagnosticFieldExpected); b51ADiagnosticFieldErr != nil {
		b51ADiagnosticFieldT.Fatal(b51ADiagnosticFieldErr)
	}
	var b51ADiagnosticFieldActual any
	if b51ADiagnosticFieldErr := b51ADiagnosticFieldJSON.Unmarshal(b51ADiagnosticFieldNext, &b51ADiagnosticFieldActual); b51ADiagnosticFieldErr != nil {
		b51ADiagnosticFieldT.Fatal(b51ADiagnosticFieldErr)
	}
	if !b51ADiagnosticFieldReflect.DeepEqual(b51ADiagnosticFieldActual, b51ADiagnosticFieldExpected) {
		b51ADiagnosticFieldT.Fatal("cross-call storage")
	}
	if b51ADiagnosticFieldStrings.Contains(string(b51ADiagnosticFieldNext), `"success"`) {
		b51ADiagnosticFieldT.Fatal("success inferred")
	}
}

func TestAb21Aq21IndependentFieldValues(b51ADiagnosticFieldT *b51ADiagnosticFieldTesting.T) {
	for _, b51ADiagnosticFieldCase := range b51ADiagnosticFieldCases(b51ADiagnosticFieldT) {
		if b51ADiagnosticFieldCase.B51ADiagnosticFieldID == "ab21aq21_disabled" {
			continue
		}
		for b51ADiagnosticFieldKey, b51ADiagnosticFieldValue := range b51ADiagnosticFieldCase.B51ADiagnosticFieldInput["fields"].(map[string]any) {
			if b51ADiagnosticFieldKey == "enable" || b51ADiagnosticFieldKey == "flags" {
				continue
			}
			for b51ADiagnosticFieldPosition := range b51ADiagnosticFieldValue.([]any) {
				var b51ADiagnosticFieldInput map[string]any
				for _, b51ADiagnosticFieldCandidate := range b51ADiagnosticFieldCases(b51ADiagnosticFieldT) {
					if b51ADiagnosticFieldCandidate.B51ADiagnosticFieldID == b51ADiagnosticFieldCase.B51ADiagnosticFieldID {
						b51ADiagnosticFieldInput = b51ADiagnosticFieldCandidate.B51ADiagnosticFieldInput
					}
				}
				b51ADiagnosticFieldArray := make([]any, len(b51ADiagnosticFieldValue.([]any)))
				for b51ADiagnosticFieldIndex := range b51ADiagnosticFieldArray {
					b51ADiagnosticFieldArray[b51ADiagnosticFieldIndex] = float64(0)
				}
				b51ADiagnosticFieldArray[b51ADiagnosticFieldPosition] = float64(1)
				b51ADiagnosticFieldInput["fields"].(map[string]any)[b51ADiagnosticFieldKey] = b51ADiagnosticFieldArray
				b51ADiagnosticFieldResult := b51ADiagnosticFieldEval(b51ADiagnosticFieldT, b51ADiagnosticFieldInput)["results"].(map[string]any)[b51ADiagnosticFieldKey].(map[string]any)
				switch b51ADiagnosticFieldKey {
				case "alternatives":
					b51ADiagnosticFieldSubset(b51ADiagnosticFieldT, b51ADiagnosticFieldResult["alternatives"], b51ADiagnosticFieldArray)
				case "reserved":
					b51ADiagnosticFieldSubset(b51ADiagnosticFieldT, b51ADiagnosticFieldResult["reserved"], b51ADiagnosticFieldArray)
				default:
					b51ADiagnosticFieldSubset(b51ADiagnosticFieldT, b51ADiagnosticFieldResult["number"], []float64{1, 256, 65536, 16777216}[b51ADiagnosticFieldPosition])
				}
			}
		}
	}
	for _, b51ADiagnosticFieldIndex := range []int{1, 3, 4} {
		b51ADiagnosticFieldInput := b51ADiagnosticFieldCases(b51ADiagnosticFieldT)[b51ADiagnosticFieldIndex].B51ADiagnosticFieldInput
		b51ADiagnosticFieldInput["fields"].(map[string]any)["flags"] = []any{float64(1)}
		for b51ADiagnosticFieldKey, b51ADiagnosticFieldValue := range b51ADiagnosticFieldEval(b51ADiagnosticFieldT, b51ADiagnosticFieldInput)["results"].(map[string]any) {
			if b51ADiagnosticFieldKey == "flags" {
				continue
			}
			b51ADiagnosticFieldResult := b51ADiagnosticFieldValue.(map[string]any)
			if b51ADiagnosticFieldResult["state"] != "Unknown" || len(b51ADiagnosticFieldResult) != 2 {
				b51ADiagnosticFieldT.Fatal("disabled interpretation", b51ADiagnosticFieldKey)
			}
		}
	}
	for _, b51ADiagnosticFieldIndex := range []int{1, 3, 4} {
		b51ADiagnosticFieldInput := b51ADiagnosticFieldCases(b51ADiagnosticFieldT)[b51ADiagnosticFieldIndex].B51ADiagnosticFieldInput
		b51ADiagnosticFieldInput["fields"] = map[string]any{"flags": []any{float64(1)}}
		b51ADiagnosticFieldOutput := b51ADiagnosticFieldEval(b51ADiagnosticFieldT, b51ADiagnosticFieldInput)
		if b51ADiagnosticFieldOutput["state"] != "complete" {
			b51ADiagnosticFieldT.Fatal("disabled flags-only state", b51ADiagnosticFieldOutput["state"])
		}
		for b51ADiagnosticFieldKey, b51ADiagnosticFieldValue := range b51ADiagnosticFieldOutput["results"].(map[string]any) {
			if b51ADiagnosticFieldKey != "flags" && b51ADiagnosticFieldValue.(map[string]any)["state"] != "Unknown" {
				b51ADiagnosticFieldT.Fatal("disabled flags-only field", b51ADiagnosticFieldKey)
			}
		}
		b51ADiagnosticFieldInput["fields"] = map[string]any{"flags": []any{float64(0)}}
		if b51ADiagnosticFieldEval(b51ADiagnosticFieldT, b51ADiagnosticFieldInput)["state"] != "partial" {
			b51ADiagnosticFieldT.Fatal("enabled short reply must stay partial")
		}
	}
	for b51ADiagnosticFieldStatus, b51ADiagnosticFieldLabel := range map[int]string{3: "range-violation", 4: "misconfiguration", 255: "unmapped"} {
		b51ADiagnosticFieldInput := b51ADiagnosticFieldCases(b51ADiagnosticFieldT)[5].B51ADiagnosticFieldInput
		b51ADiagnosticFieldInput["fields"].(map[string]any)["status"] = []any{float64(b51ADiagnosticFieldStatus)}
		b51ADiagnosticFieldSubset(b51ADiagnosticFieldT, b51ADiagnosticFieldEval(b51ADiagnosticFieldT, b51ADiagnosticFieldInput)["results"].(map[string]any)["status"], map[string]any{"label": b51ADiagnosticFieldLabel})
	}
	for b51ADiagnosticFieldBit, b51ADiagnosticFieldName := range []string{"disabled", "writable", "enum-versus-range", "additional-four-byte/no-range", "reserved", "reserved", "text", "text"} {
		b51ADiagnosticFieldInput := b51ADiagnosticFieldCases(b51ADiagnosticFieldT)[1].B51ADiagnosticFieldInput
		b51ADiagnosticFieldInput["fields"].(map[string]any)["flags"] = []any{float64(uint64(1) << b51ADiagnosticFieldBit)}
		b51ADiagnosticFieldResult := b51ADiagnosticFieldEval(b51ADiagnosticFieldT, b51ADiagnosticFieldInput)["results"].(map[string]any)["flags"].(map[string]any)
		if b51ADiagnosticFieldBit < 4 {
			for b51ADiagnosticFieldKey, b51ADiagnosticFieldValue := range b51ADiagnosticFieldResult["flags"].(map[string]any) {
				if b51ADiagnosticFieldValue != (b51ADiagnosticFieldKey == b51ADiagnosticFieldName) {
					b51ADiagnosticFieldT.Fatal("flag position", b51ADiagnosticFieldBit)
				}
			}
		} else if b51ADiagnosticFieldBit < 6 {
			b51ADiagnosticFieldSubset(b51ADiagnosticFieldT, b51ADiagnosticFieldResult["reserved"], []float64{1, 2}[b51ADiagnosticFieldBit-4])
		} else {
			b51ADiagnosticFieldSubset(b51ADiagnosticFieldT, b51ADiagnosticFieldResult["unmodeled"], []float64{1, 2}[b51ADiagnosticFieldBit-6])
		}
	}
	b51ADiagnosticFieldInput := b51ADiagnosticFieldCases(b51ADiagnosticFieldT)[5].B51ADiagnosticFieldInput
	b51ADiagnosticFieldInput["arguments"].(map[string]any)["value_width"] = float64(1)
	b51ADiagnosticFieldResult := b51ADiagnosticFieldEval(b51ADiagnosticFieldT, b51ADiagnosticFieldInput)["results"].(map[string]any)["value"].(map[string]any)
	if b51ADiagnosticFieldResult["state"] != "incomplete" {
		b51ADiagnosticFieldT.Fatal("width inferred from bytes")
	}
}
