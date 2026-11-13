package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type b504CandidateFieldCase struct {
	B504CandidateFieldID       string          `json:"id"`
	B504CandidateFieldInput    json.RawMessage `json:"input"`
	B504CandidateFieldExpected map[string]any  `json:"expected"`
	B504CandidateFieldState    string          `json:"state"`
}

func b504CandidateFieldCases(b504CandidateFieldT *testing.T) []b504CandidateFieldCase {
	b504CandidateFieldT.Helper()
	b504CandidateFieldRaw, b504CandidateFieldErr := os.ReadFile("testdata/ab01aq01_b504_candidate_fields.json")
	if b504CandidateFieldErr != nil {
		b504CandidateFieldT.Fatal(b504CandidateFieldErr)
	}
	var b504CandidateFieldCases []b504CandidateFieldCase
	if b504CandidateFieldErr = json.Unmarshal(b504CandidateFieldRaw, &b504CandidateFieldCases); b504CandidateFieldErr != nil {
		b504CandidateFieldT.Fatal(b504CandidateFieldErr)
	}
	if len(b504CandidateFieldCases) == 0 {
		b504CandidateFieldT.Fatal("zero fixture cases")
	}
	return b504CandidateFieldCases
}

func b504CandidateFieldEvaluate(b504CandidateFieldT *testing.T, b504CandidateFieldInput any) map[string]any {
	b504CandidateFieldT.Helper()
	b504CandidateFieldRaw, b504CandidateFieldErr := json.Marshal(b504CandidateFieldInput)
	if b504CandidateFieldErr != nil {
		b504CandidateFieldT.Fatal(b504CandidateFieldErr)
	}
	b504CandidateFieldOutput, b504CandidateFieldErr := B504CandidateFieldEvaluate(b504CandidateFieldRaw)
	if b504CandidateFieldErr != nil {
		b504CandidateFieldT.Fatal(b504CandidateFieldErr)
	}
	var b504CandidateFieldValue map[string]any
	if b504CandidateFieldErr = json.Unmarshal(b504CandidateFieldOutput, &b504CandidateFieldValue); b504CandidateFieldErr != nil {
		b504CandidateFieldT.Fatal(b504CandidateFieldErr)
	}
	if b504CandidateFieldValue["physical"] != "Unknown" || b504CandidateFieldValue["qualified"] != false {
		b504CandidateFieldT.Fatal(b504CandidateFieldValue)
	}
	return b504CandidateFieldValue
}

func b504CandidateFieldDecode(b504CandidateFieldT *testing.T, b504CandidateFieldRaw json.RawMessage) map[string]any {
	b504CandidateFieldT.Helper()
	var b504CandidateFieldInput map[string]any
	if b504CandidateFieldErr := json.Unmarshal(b504CandidateFieldRaw, &b504CandidateFieldInput); b504CandidateFieldErr != nil {
		b504CandidateFieldT.Fatal(b504CandidateFieldErr)
	}
	return b504CandidateFieldInput
}

func TestAb01Aq01Fixture(b504CandidateFieldT *testing.T) {
	for _, b504CandidateFieldCase := range b504CandidateFieldCases(b504CandidateFieldT) {
		b504CandidateFieldT.Run(b504CandidateFieldCase.B504CandidateFieldID, func(b504CandidateFieldT *testing.T) {
			b504CandidateFieldInput := b504CandidateFieldDecode(b504CandidateFieldT, b504CandidateFieldCase.B504CandidateFieldInput)
			b504CandidateFieldOutput := b504CandidateFieldEvaluate(b504CandidateFieldT, b504CandidateFieldInput)
			if b504CandidateFieldOutput["state"] != b504CandidateFieldCase.B504CandidateFieldState || !reflect.DeepEqual(b504CandidateFieldOutput["results"], b504CandidateFieldCase.B504CandidateFieldExpected) {
				b504CandidateFieldT.Fatalf("got %v; expected %v (%s)", b504CandidateFieldOutput, b504CandidateFieldCase.B504CandidateFieldExpected, b504CandidateFieldCase.B504CandidateFieldState)
			}
			for b504CandidateFieldKey, b504CandidateFieldValue := range b504CandidateFieldInput {
				if !reflect.DeepEqual(b504CandidateFieldOutput[b504CandidateFieldKey], b504CandidateFieldValue) {
					b504CandidateFieldT.Fatalf("lost %s", b504CandidateFieldKey)
				}
			}
		})
	}
}

func TestAb01Aq01ContextAxes(b504CandidateFieldT *testing.T) {
	for _, b504CandidateFieldCase := range b504CandidateFieldCases(b504CandidateFieldT) {
		for _, b504CandidateFieldKey := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
			for _, b504CandidateFieldMode := range []string{"changed", "absent", "zero"} {
				b504CandidateFieldT.Run(b504CandidateFieldCase.B504CandidateFieldID+"/"+b504CandidateFieldKey+"/"+b504CandidateFieldMode, func(b504CandidateFieldT *testing.T) {
					b504CandidateFieldInput := b504CandidateFieldDecode(b504CandidateFieldT, b504CandidateFieldCase.B504CandidateFieldInput)
					b504CandidateFieldContext := b504CandidateFieldInput["context"].(map[string]any)
					b504CandidateFieldOpaque := b504CandidateFieldKey == "family" || b504CandidateFieldKey == "variant"
					switch b504CandidateFieldKey {
					case "primary", "secondary", "operation":
						b504CandidateFieldContext[b504CandidateFieldKey] = float64(0)
						if b504CandidateFieldMode == "changed" {
							b504CandidateFieldContext[b504CandidateFieldKey] = float64(7)
						}
					default:
						b504CandidateFieldContext[b504CandidateFieldKey] = ""
						if b504CandidateFieldMode == "changed" {
							b504CandidateFieldContext[b504CandidateFieldKey] = "other"
						}
					}
					if b504CandidateFieldMode == "absent" {
						delete(b504CandidateFieldContext, b504CandidateFieldKey)
					}
					b504CandidateFieldOutput := b504CandidateFieldEvaluate(b504CandidateFieldT, b504CandidateFieldInput)
					if !reflect.DeepEqual(b504CandidateFieldOutput["context"], b504CandidateFieldContext) {
						b504CandidateFieldT.Fatal("context lost")
					}
					if b504CandidateFieldOpaque {
						if !reflect.DeepEqual(b504CandidateFieldOutput["results"], b504CandidateFieldCase.B504CandidateFieldExpected) {
							b504CandidateFieldT.Fatal("opaque identity selected semantics")
						}
						return
					}
					if b504CandidateFieldOutput["state"] != "Unknown" {
						b504CandidateFieldT.Fatal(b504CandidateFieldOutput)
					}
					for _, b504CandidateFieldResult := range b504CandidateFieldOutput["results"].(map[string]any) {
						b504CandidateFieldField := b504CandidateFieldResult.(map[string]any)
						if b504CandidateFieldField["state"] != "Unknown" || len(b504CandidateFieldField) != 3 {
							b504CandidateFieldT.Fatal("unsupported conversion", b504CandidateFieldField)
						}
					}
				})
			}
		}
	}
}

func TestAb01Aq01WidthsAndPresence(b504CandidateFieldT *testing.T) {
	for _, b504CandidateFieldCase := range b504CandidateFieldCases(b504CandidateFieldT) {
		if b504CandidateFieldCase.B504CandidateFieldState != "complete" {
			continue
		}
		b504CandidateFieldOriginal := b504CandidateFieldDecode(b504CandidateFieldT, b504CandidateFieldCase.B504CandidateFieldInput)
		for b504CandidateFieldKey := range b504CandidateFieldOriginal["fields"].(map[string]any) {
			for _, b504CandidateFieldMode := range []string{"missing", "empty", "long"} {
				b504CandidateFieldT.Run(b504CandidateFieldCase.B504CandidateFieldID+"/"+b504CandidateFieldKey+"/"+b504CandidateFieldMode, func(b504CandidateFieldT *testing.T) {
					b504CandidateFieldInput := b504CandidateFieldDecode(b504CandidateFieldT, b504CandidateFieldCase.B504CandidateFieldInput)
					b504CandidateFieldFields := b504CandidateFieldInput["fields"].(map[string]any)
					switch b504CandidateFieldMode {
					case "missing":
						delete(b504CandidateFieldFields, b504CandidateFieldKey)
					case "empty":
						b504CandidateFieldFields[b504CandidateFieldKey] = []any{}
					case "long":
						b504CandidateFieldFields[b504CandidateFieldKey] = []any{float64(0), float64(0), float64(0)}
					}
					b504CandidateFieldOutput := b504CandidateFieldEvaluate(b504CandidateFieldT, b504CandidateFieldInput)
					b504CandidateFieldResults := b504CandidateFieldOutput["results"].(map[string]any)
					if b504CandidateFieldKey == "dhw_request" && b504CandidateFieldMode == "missing" {
						if _, b504CandidateFieldPresent := b504CandidateFieldResults[b504CandidateFieldKey]; b504CandidateFieldPresent {
							b504CandidateFieldT.Fatal("optional absent became present")
						}
						return
					}
					b504CandidateFieldResult := b504CandidateFieldResults[b504CandidateFieldKey].(map[string]any)
					if b504CandidateFieldOutput["state"] != "partial" || b504CandidateFieldResult["state"] != "incomplete" || b504CandidateFieldResult["number"] != nil {
						b504CandidateFieldT.Fatal(b504CandidateFieldOutput)
					}
					if b504CandidateFieldResult["present"] != (b504CandidateFieldMode != "missing") || !reflect.DeepEqual(b504CandidateFieldResult["raw"], b504CandidateFieldFields[b504CandidateFieldKey]) {
						b504CandidateFieldT.Fatal("presence/raw lost")
					}
					for b504CandidateFieldSibling, b504CandidateFieldExpected := range b504CandidateFieldCase.B504CandidateFieldExpected {
						if b504CandidateFieldSibling != b504CandidateFieldKey && !reflect.DeepEqual(b504CandidateFieldResults[b504CandidateFieldSibling], b504CandidateFieldExpected) {
							b504CandidateFieldT.Fatal("sibling erased")
						}
					}
				})
			}
		}
	}
}

func TestAb01Aq01RejectInvalid(b504CandidateFieldT *testing.T) {
	for b504CandidateFieldIndex, b504CandidateFieldRaw := range []string{
		``, `null`, `[]`, `{`, `{} {}`, `{"context":{},"fields":{},"x":{"a":1,"a":2}}`,
		`{"context":{},"fields":{},"x":[{"a":1,"\u0061":2}]}`, `{"context":{},"fields":{},"x":1e999}`, `{"context":{},"fields":{},"x":NaN}`,
	} {
		b504CandidateFieldT.Run(string(rune('a'+b504CandidateFieldIndex)), func(b504CandidateFieldT *testing.T) {
			b504CandidateFieldOutput, b504CandidateFieldErr := B504CandidateFieldEvaluate([]byte(b504CandidateFieldRaw))
			if b504CandidateFieldErr == nil || b504CandidateFieldOutput != nil {
				b504CandidateFieldT.Fatal("accepted invalid input", b504CandidateFieldRaw)
			}
		})
	}
	for _, b504CandidateFieldKey := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant", "fields", "context", "arguments", "temperature", "divisor", "role"} {
		b504CandidateFieldInvalid := []any{nil, true, "1", 1.5, -1, 256, []any{}, map[string]any{}}
		switch b504CandidateFieldKey {
		case "divisor":
			b504CandidateFieldInvalid = []any{nil, true, "1", 1.5, -1, 0, 15, 255, 257, []any{}, map[string]any{}}
		case "role", "direction", "view", "family", "variant":
			b504CandidateFieldInvalid = []any{nil, true, 1.5, -1, 256, []any{}, map[string]any{}}
		case "fields", "context", "arguments":
			b504CandidateFieldInvalid = []any{nil, true, "1", 1.5, -1, 256, []any{}}
		}
		for _, b504CandidateFieldValue := range b504CandidateFieldInvalid {
			b504CandidateFieldT.Run(b504CandidateFieldKey+"/"+b504CandidateFieldDescribe(b504CandidateFieldValue), func(b504CandidateFieldT *testing.T) {
				b504CandidateFieldInput := b504CandidateFieldDecode(b504CandidateFieldT, b504CandidateFieldCases(b504CandidateFieldT)[0].B504CandidateFieldInput)
				switch b504CandidateFieldKey {
				case "fields", "context", "arguments":
					b504CandidateFieldInput[b504CandidateFieldKey] = b504CandidateFieldValue
				case "temperature":
					b504CandidateFieldInput["fields"].(map[string]any)[b504CandidateFieldKey] = []any{b504CandidateFieldValue, 0}
				case "divisor":
					b504CandidateFieldInput["arguments"].(map[string]any)[b504CandidateFieldKey] = b504CandidateFieldValue
				case "role":
					b504CandidateFieldInput["context"].(map[string]any)["view"] = "circuit-role"
					b504CandidateFieldInput["context"].(map[string]any)["operation"] = 13
					b504CandidateFieldInput["arguments"] = map[string]any{"role": b504CandidateFieldValue}
				default:
					b504CandidateFieldInput["context"].(map[string]any)[b504CandidateFieldKey] = b504CandidateFieldValue
				}
				b504CandidateFieldRaw, _ := json.Marshal(b504CandidateFieldInput)
				b504CandidateFieldOutput, b504CandidateFieldErr := B504CandidateFieldEvaluate(b504CandidateFieldRaw)
				if b504CandidateFieldErr == nil || b504CandidateFieldOutput != nil {
					b504CandidateFieldT.Fatalf("accepted %s", b504CandidateFieldRaw)
				}
			})
		}
	}
}

func b504CandidateFieldDescribe(b504CandidateFieldValue any) string {
	b504CandidateFieldRaw, _ := json.Marshal(b504CandidateFieldValue)
	return string(b504CandidateFieldRaw)
}

func TestAb01Aq01ArgumentsAndMemory(b504CandidateFieldT *testing.T) {
	for _, b504CandidateFieldIndex := range []int{0, 14} {
		b504CandidateFieldInput := b504CandidateFieldDecode(b504CandidateFieldT, b504CandidateFieldCases(b504CandidateFieldT)[b504CandidateFieldIndex].B504CandidateFieldInput)
		delete(b504CandidateFieldInput, "arguments")
		b504CandidateFieldRaw, _ := json.Marshal(b504CandidateFieldInput)
		if b504CandidateFieldOutput, b504CandidateFieldErr := B504CandidateFieldEvaluate(b504CandidateFieldRaw); b504CandidateFieldErr == nil || b504CandidateFieldOutput != nil {
			b504CandidateFieldT.Fatal("missing required argument")
		}
	}
	b504CandidateFieldInput := b504CandidateFieldDecode(b504CandidateFieldT, b504CandidateFieldCases(b504CandidateFieldT)[0].B504CandidateFieldInput)
	b504CandidateFieldInput["extra"] = map[string]any{"absent": nil, "zero": float64(0), "false": false, "empty": []any{}, "object": map[string]any{}}
	b504CandidateFieldInput["fields"].(map[string]any)["opaque"] = []any{}
	b504CandidateFieldExpected := b504CandidateFieldEvaluate(b504CandidateFieldT, b504CandidateFieldInput)
	if !reflect.DeepEqual(b504CandidateFieldExpected["extra"], b504CandidateFieldInput["extra"]) || !reflect.DeepEqual(b504CandidateFieldExpected["fields"], b504CandidateFieldInput["fields"]) {
		b504CandidateFieldT.Fatal("extra values lost")
	}
	b504CandidateFieldRaw, _ := json.Marshal(b504CandidateFieldInput)
	b504CandidateFieldOriginal := append([]byte(nil), b504CandidateFieldRaw...)
	b504CandidateFieldOutput, b504CandidateFieldErr := B504CandidateFieldEvaluate(b504CandidateFieldRaw)
	if b504CandidateFieldErr != nil {
		b504CandidateFieldT.Fatal(b504CandidateFieldErr)
	}
	for b504CandidateFieldIndex := range b504CandidateFieldRaw {
		b504CandidateFieldRaw[b504CandidateFieldIndex] = ' '
	}
	if !reflect.DeepEqual(b504CandidateFieldDecode(b504CandidateFieldT, b504CandidateFieldOutput), b504CandidateFieldExpected) {
		b504CandidateFieldT.Fatal("input aliases output")
	}
	for b504CandidateFieldIndex := range b504CandidateFieldOutput {
		b504CandidateFieldOutput[b504CandidateFieldIndex] = ' '
	}
	if !reflect.DeepEqual(b504CandidateFieldEvaluate(b504CandidateFieldT, json.RawMessage(b504CandidateFieldOriginal)), b504CandidateFieldExpected) {
		b504CandidateFieldT.Fatal("call state retained")
	}
}

func TestAb01Aq01InputBoundaries(b504CandidateFieldT *testing.T) {
	for _, b504CandidateFieldRaw := range []string{
		`{"context":{},"context":{},"fields":{}}`,
		`{"context":{"primary":181,"primary":181},"fields":{}}`,
		`{"context":{},"fields":{"x":[],"x":[]}}`,
		`{"context":{},"fields":{},"arguments":{"role":"dhw","role":"dhw"}}`,
		`{"context":{},"fields":{}} true`,
		`{"context":{},"fields":{"x":null}}`,
		`{"context":{},"fields":{"x":false}}`,
		`{"context":{},"fields":{"x":""}}`,
		`{"context":{},"fields":{"x":{}}}`,
	} {
		b504CandidateFieldOutput, b504CandidateFieldErr := B504CandidateFieldEvaluate([]byte(b504CandidateFieldRaw))
		if b504CandidateFieldErr == nil || b504CandidateFieldOutput != nil {
			b504CandidateFieldT.Fatalf("accepted %s", b504CandidateFieldRaw)
		}
	}
	for _, b504CandidateFieldIndex := range []int{0, 4, 7, 10, 14} {
		b504CandidateFieldInput := b504CandidateFieldDecode(b504CandidateFieldT, b504CandidateFieldCases(b504CandidateFieldT)[b504CandidateFieldIndex].B504CandidateFieldInput)
		b504CandidateFieldInput["arguments"] = map[string]any{"unexpected": false}
		b504CandidateFieldRaw, _ := json.Marshal(b504CandidateFieldInput)
		if b504CandidateFieldOutput, b504CandidateFieldErr := B504CandidateFieldEvaluate(b504CandidateFieldRaw); b504CandidateFieldErr == nil || b504CandidateFieldOutput != nil {
			b504CandidateFieldT.Fatal("unexpected argument accepted")
		}
	}
	b504CandidateFieldInput := b504CandidateFieldDecode(b504CandidateFieldT, b504CandidateFieldCases(b504CandidateFieldT)[0].B504CandidateFieldInput)
	b504CandidateFieldInput["context"].(map[string]any)["opaque"] = map[string]any{"null": nil, "zero": float64(0), "false": false, "empty": ""}
	b504CandidateFieldInput["fields"].(map[string]any)["tail"] = []any{float64(255), float64(0)}
	b504CandidateFieldOutput := b504CandidateFieldEvaluate(b504CandidateFieldT, b504CandidateFieldInput)
	if !reflect.DeepEqual(b504CandidateFieldOutput["context"], b504CandidateFieldInput["context"]) {
		b504CandidateFieldT.Fatal("opaque context lost")
	}
	b504CandidateFieldTail := b504CandidateFieldOutput["results"].(map[string]any)["tail"].(map[string]any)
	if b504CandidateFieldTail["state"] != "Unknown" || len(b504CandidateFieldTail) != 3 {
		b504CandidateFieldT.Fatal("tail interpreted")
	}
	b504CandidateFieldInput["context"].(map[string]any)["operation"] = float64(0)
	b504CandidateFieldInput["arguments"] = map[string]any{"extra": nil, "zero": float64(0), "false": false, "empty": []any{}}
	b504CandidateFieldOutput = b504CandidateFieldEvaluate(b504CandidateFieldT, b504CandidateFieldInput)
	if b504CandidateFieldOutput["state"] != "Unknown" || !reflect.DeepEqual(b504CandidateFieldOutput["arguments"], b504CandidateFieldInput["arguments"]) {
		b504CandidateFieldT.Fatal("unknown arguments lost")
	}
}
