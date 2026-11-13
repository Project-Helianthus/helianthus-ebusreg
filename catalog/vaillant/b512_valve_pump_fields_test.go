package vaillant

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func b512ValvePumpFieldEvaluateTest(b512ValvePumpFieldT *testing.T, b512ValvePumpFieldInput any) map[string]any {
	b512ValvePumpFieldT.Helper()
	b512ValvePumpFieldBytes, b512ValvePumpFieldErr := json.Marshal(b512ValvePumpFieldInput)
	if b512ValvePumpFieldErr != nil {
		b512ValvePumpFieldT.Fatal(b512ValvePumpFieldErr)
	}
	b512ValvePumpFieldOutput, b512ValvePumpFieldErr := B512ValvePumpFieldEvaluate(b512ValvePumpFieldBytes)
	if b512ValvePumpFieldErr != nil {
		b512ValvePumpFieldT.Fatal(b512ValvePumpFieldErr)
	}
	b512ValvePumpFieldResult := map[string]any{}
	if b512ValvePumpFieldErr := json.Unmarshal(b512ValvePumpFieldOutput, &b512ValvePumpFieldResult); b512ValvePumpFieldErr != nil {
		b512ValvePumpFieldT.Fatal(b512ValvePumpFieldErr)
	}
	return b512ValvePumpFieldResult
}

func b512ValvePumpFieldInputTest(b512ValvePumpFieldOp int, b512ValvePumpFieldDirection string, b512ValvePumpFieldFields map[string]any) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 18, "operation": b512ValvePumpFieldOp, "direction": b512ValvePumpFieldDirection, "view": "actuator-fields", "family": "synthetic-family", "variant": "synthetic-variant"}, "fields": b512ValvePumpFieldFields}
}

func TestAb23Aq23Fixture(b512ValvePumpFieldT *testing.T) {
	b512ValvePumpFieldData, b512ValvePumpFieldErr := os.ReadFile("testdata/ab23aq23_b512_valve_pump_fields.json")
	if b512ValvePumpFieldErr != nil {
		b512ValvePumpFieldT.Fatal(b512ValvePumpFieldErr)
	}
	var b512ValvePumpFieldCases []struct {
		B512ValvePumpFieldID       string                    `json:"id"`
		B512ValvePumpFieldInput    map[string]any            `json:"input"`
		B512ValvePumpFieldExpected map[string]map[string]any `json:"expected"`
	}
	if b512ValvePumpFieldErr := json.Unmarshal(b512ValvePumpFieldData, &b512ValvePumpFieldCases); b512ValvePumpFieldErr != nil {
		b512ValvePumpFieldT.Fatal(b512ValvePumpFieldErr)
	}
	if len(b512ValvePumpFieldCases) == 0 {
		b512ValvePumpFieldT.Fatal("zero cases")
	}
	for _, b512ValvePumpFieldCase := range b512ValvePumpFieldCases {
		b512ValvePumpFieldT.Run(b512ValvePumpFieldCase.B512ValvePumpFieldID, func(b512ValvePumpFieldT *testing.T) {
			b512ValvePumpFieldOutput := b512ValvePumpFieldEvaluateTest(b512ValvePumpFieldT, b512ValvePumpFieldCase.B512ValvePumpFieldInput)
			if b512ValvePumpFieldOutput["state"] != "complete" || b512ValvePumpFieldOutput["physical"] != "Unknown" || b512ValvePumpFieldOutput["qualified"] != false {
				b512ValvePumpFieldT.Fatal(b512ValvePumpFieldOutput)
			}
			if !reflect.DeepEqual(b512ValvePumpFieldOutput["context"], b512ValvePumpFieldCase.B512ValvePumpFieldInput["context"]) || !reflect.DeepEqual(b512ValvePumpFieldOutput["arguments"], b512ValvePumpFieldCase.B512ValvePumpFieldInput["arguments"]) {
				b512ValvePumpFieldT.Fatal("context or arguments changed")
			}
			if !reflect.DeepEqual(b512ValvePumpFieldOutput["fields"], b512ValvePumpFieldCase.B512ValvePumpFieldInput["fields"]) {
				b512ValvePumpFieldT.Fatal("raw changed")
			}
			for b512ValvePumpFieldField, b512ValvePumpFieldExpected := range b512ValvePumpFieldCase.B512ValvePumpFieldExpected {
				b512ValvePumpFieldResult := b512ValvePumpFieldOutput["results"].(map[string]any)[b512ValvePumpFieldField].(map[string]any)
				if b512ValvePumpFieldResult["state"] != "complete" || b512ValvePumpFieldResult["meaning"] != "Hypothesis" {
					b512ValvePumpFieldT.Fatal(b512ValvePumpFieldResult)
				}
				for b512ValvePumpFieldKey, b512ValvePumpFieldWant := range b512ValvePumpFieldExpected {
					if !reflect.DeepEqual(b512ValvePumpFieldResult[b512ValvePumpFieldKey], b512ValvePumpFieldWant) {
						b512ValvePumpFieldT.Fatalf("%s.%s got %v want %v", b512ValvePumpFieldField, b512ValvePumpFieldKey, b512ValvePumpFieldResult[b512ValvePumpFieldKey], b512ValvePumpFieldWant)
					}
				}
				if b512ValvePumpFieldResult["label"] == "unresolved" {
					if _, b512ValvePumpFieldExists := b512ValvePumpFieldResult["number"]; b512ValvePumpFieldExists {
						b512ValvePumpFieldT.Fatal("inferred actual feedback")
					}
				} else if b512ValvePumpFieldResult["calculation"] != "Proven" {
					b512ValvePumpFieldT.Fatal("calculation grade")
				}
			}
			if b512ValvePumpFieldCase.B512ValvePumpFieldInput["context"].(map[string]any)["operation"] == float64(5) {
				b512ValvePumpFieldBytes, _ := json.Marshal(b512ValvePumpFieldOutput)
				if bytes.Contains(b512ValvePumpFieldBytes, []byte("900")) || bytes.Contains(b512ValvePumpFieldBytes, []byte("timeout")) {
					b512ValvePumpFieldT.Fatal("inherited timeout")
				}
			}
		})
	}
}

func TestAb23Aq23Context(b512ValvePumpFieldT *testing.T) {
	for _, b512ValvePumpFieldOperation := range []int{0, 5} {
		b512ValvePumpFieldInput := b512ValvePumpFieldInputTest(b512ValvePumpFieldOperation, "positive", map[string]any{"setpoint": []int{100}})
		b512ValvePumpFieldOutput := b512ValvePumpFieldEvaluateTest(b512ValvePumpFieldT, b512ValvePumpFieldInput)
		if b512ValvePumpFieldOutput["state"] != "Unknown" || len(b512ValvePumpFieldOutput["results"].(map[string]any)["setpoint"].(map[string]any)) != 2 {
			b512ValvePumpFieldT.Fatal("undeclared positive view")
		}
	}
	for b512ValvePumpFieldKey, b512ValvePumpFieldValue := range map[string]any{"primary": 180, "secondary": 19, "operation": 1, "direction": "response", "view": "other"} {
		b512ValvePumpFieldT.Run(b512ValvePumpFieldKey, func(b512ValvePumpFieldT *testing.T) {
			b512ValvePumpFieldInput := b512ValvePumpFieldInputTest(0, "request", map[string]any{"setpoint": []int{100}})
			b512ValvePumpFieldInput["context"].(map[string]any)[b512ValvePumpFieldKey] = b512ValvePumpFieldValue
			b512ValvePumpFieldOutput := b512ValvePumpFieldEvaluateTest(b512ValvePumpFieldT, b512ValvePumpFieldInput)
			if b512ValvePumpFieldOutput["state"] != "Unknown" || len(b512ValvePumpFieldOutput["results"].(map[string]any)["setpoint"].(map[string]any)) != 2 {
				b512ValvePumpFieldT.Fatal(b512ValvePumpFieldOutput)
			}
		})
	}
	for _, b512ValvePumpFieldKey := range []string{"primary", "secondary", "operation"} {
		b512ValvePumpFieldInput := b512ValvePumpFieldInputTest(0, "request", map[string]any{"setpoint": []int{100}})
		delete(b512ValvePumpFieldInput["context"].(map[string]any), b512ValvePumpFieldKey)
		if b512ValvePumpFieldEvaluateTest(b512ValvePumpFieldT, b512ValvePumpFieldInput)["state"] != "Unknown" {
			b512ValvePumpFieldT.Fatal("absent selector")
		}
	}
	for _, b512ValvePumpFieldKey := range []string{"family", "variant"} {
		b512ValvePumpFieldAbsent := b512ValvePumpFieldInputTest(0, "request", map[string]any{"setpoint": []int{100}})
		delete(b512ValvePumpFieldAbsent["context"].(map[string]any), b512ValvePumpFieldKey)
		b512ValvePumpFieldAbsentOutput := b512ValvePumpFieldEvaluateTest(b512ValvePumpFieldT, b512ValvePumpFieldAbsent)
		if _, b512ValvePumpFieldExists := b512ValvePumpFieldAbsentOutput["context"].(map[string]any)[b512ValvePumpFieldKey]; b512ValvePumpFieldExists || b512ValvePumpFieldAbsentOutput["qualified"] != false {
			b512ValvePumpFieldT.Fatal("absent opaque identity changed")
		}
		for _, b512ValvePumpFieldValue := range []any{nil, false, 0, "", "another-family", []any{}} {
			b512ValvePumpFieldInput := b512ValvePumpFieldInputTest(0, "request", map[string]any{"setpoint": []int{100}})
			b512ValvePumpFieldInput["context"].(map[string]any)[b512ValvePumpFieldKey] = b512ValvePumpFieldValue
			b512ValvePumpFieldOutput := b512ValvePumpFieldEvaluateTest(b512ValvePumpFieldT, b512ValvePumpFieldInput)
			if b512ValvePumpFieldOutput["qualified"] != false || b512ValvePumpFieldOutput["state"] != "complete" || !b512ValvePumpFieldEqualTest(b512ValvePumpFieldOutput["context"].(map[string]any)[b512ValvePumpFieldKey], b512ValvePumpFieldValue) {
				b512ValvePumpFieldT.Fatal(b512ValvePumpFieldOutput)
			}
		}
	}
}

func TestAb23Aq23WidthsAndIndependence(b512ValvePumpFieldT *testing.T) {
	for _, b512ValvePumpFieldOp := range []int{0, 2, 3, 4, 5} {
		for _, b512ValvePumpFieldDirection := range []string{"request", "positive"} {
			b512ValvePumpFieldFields := map[string]any{}
			b512ValvePumpFieldWidths := map[string]int{}
			if b512ValvePumpFieldDirection == "request" {
				switch b512ValvePumpFieldOp {
				case 0, 4:
					b512ValvePumpFieldWidths["setpoint"] = 1
				case 2:
					b512ValvePumpFieldWidths["valve"] = 1
				case 3:
					b512ValvePumpFieldWidths["selector"] = 1
				case 5:
					b512ValvePumpFieldWidths = map[string]int{"channel0": 1, "channel1": 1, "channel2": 1}
				}
			} else {
				switch b512ValvePumpFieldOp {
				case 2:
					b512ValvePumpFieldWidths = map[string]int{"valve": 1, "contact": 1, "tank": 2, "virtual_error": 2}
				case 3:
					b512ValvePumpFieldWidths = map[string]int{"status": 1, "temperature_raw": 2}
				case 4:
					b512ValvePumpFieldWidths["feedback"] = 1
				}
			}
			for b512ValvePumpFieldKey, b512ValvePumpFieldWidth := range b512ValvePumpFieldWidths {
				b512ValvePumpFieldFields[b512ValvePumpFieldKey] = make([]int, b512ValvePumpFieldWidth)
			}
			for b512ValvePumpFieldKey, b512ValvePumpFieldWidth := range b512ValvePumpFieldWidths {
				for _, b512ValvePumpFieldSize := range []int{-1, 0, b512ValvePumpFieldWidth - 1, b512ValvePumpFieldWidth + 1} {
					b512ValvePumpFieldInput := b512ValvePumpFieldInputTest(b512ValvePumpFieldOp, b512ValvePumpFieldDirection, map[string]any{})
					for b512ValvePumpFieldSibling, b512ValvePumpFieldData := range b512ValvePumpFieldFields {
						b512ValvePumpFieldInput["fields"].(map[string]any)[b512ValvePumpFieldSibling] = b512ValvePumpFieldData
					}
					if b512ValvePumpFieldSize < 0 {
						delete(b512ValvePumpFieldInput["fields"].(map[string]any), b512ValvePumpFieldKey)
					} else {
						b512ValvePumpFieldInput["fields"].(map[string]any)[b512ValvePumpFieldKey] = make([]int, b512ValvePumpFieldSize)
					}
					if b512ValvePumpFieldOp == 5 {
						b512ValvePumpFieldInput["arguments"] = map[string]any{"external_modes": []int{8, 0, 8}}
					}
					b512ValvePumpFieldOutput := b512ValvePumpFieldEvaluateTest(b512ValvePumpFieldT, b512ValvePumpFieldInput)
					if b512ValvePumpFieldOutput["state"] != "partial" {
						b512ValvePumpFieldT.Fatal(b512ValvePumpFieldOutput)
					}
					for b512ValvePumpFieldSibling := range b512ValvePumpFieldWidths {
						b512ValvePumpFieldResult := b512ValvePumpFieldOutput["results"].(map[string]any)[b512ValvePumpFieldSibling].(map[string]any)
						if b512ValvePumpFieldSibling == b512ValvePumpFieldKey {
							if b512ValvePumpFieldResult["state"] != "incomplete" || b512ValvePumpFieldResult["number"] != nil {
								b512ValvePumpFieldT.Fatal(b512ValvePumpFieldResult)
							}
							_, b512ValvePumpFieldHasRaw := b512ValvePumpFieldResult["raw"]
							if b512ValvePumpFieldHasRaw != (b512ValvePumpFieldSize >= 0) {
								b512ValvePumpFieldT.Fatal("missing vs empty")
							}
						} else if b512ValvePumpFieldResult["state"] != "complete" {
							b512ValvePumpFieldT.Fatal("sibling lost")
						}
					}
				}
			}
		}
	}
	for b512ValvePumpFieldIndex := 0; b512ValvePumpFieldIndex < 3; b512ValvePumpFieldIndex++ {
		b512ValvePumpFieldModes := []int{0, 0, 0}
		b512ValvePumpFieldModes[b512ValvePumpFieldIndex] = 8
		b512ValvePumpFieldInput := b512ValvePumpFieldInputTest(5, "request", map[string]any{"channel0": []int{0}, "channel1": []int{1}, "channel2": []int{255}})
		b512ValvePumpFieldInput["arguments"] = map[string]any{"external_modes": b512ValvePumpFieldModes}
		b512ValvePumpFieldOutput := b512ValvePumpFieldEvaluateTest(b512ValvePumpFieldT, b512ValvePumpFieldInput)
		for b512ValvePumpFieldChannel, b512ValvePumpFieldKey := range []string{"channel0", "channel1", "channel2"} {
			if b512ValvePumpFieldOutput["results"].(map[string]any)[b512ValvePumpFieldKey].(map[string]any)["prerequisite"] != (b512ValvePumpFieldChannel == b512ValvePumpFieldIndex) {
				b512ValvePumpFieldT.Fatal("shared mode")
			}
		}
	}
}

func TestAb23Aq23Reject(b512ValvePumpFieldT *testing.T) {
	for _, b512ValvePumpFieldJSON := range []string{``, `[]`, `{`, `{} {}`, `{"x":1,"x":2}`, `{"context":{"x":1,"x":2},"fields":{}}`, `{"context":{},"fields":{"x":[{"a":1,"a":2}]}}`, `{"context":{},"fields":{},"extra":[{"x":1,"x":2}]}`, `{"context":{},"fields":{},"extra":NaN}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":null}`, `{"context":{},"fields":{},"arguments":null}`} {
		if b512ValvePumpFieldOutput, b512ValvePumpFieldErr := B512ValvePumpFieldEvaluate([]byte(b512ValvePumpFieldJSON)); b512ValvePumpFieldErr == nil || b512ValvePumpFieldOutput != nil {
			b512ValvePumpFieldT.Fatalf("accepted %s", b512ValvePumpFieldJSON)
		}
	}
	for _, b512ValvePumpFieldBad := range []any{nil, true, "0", -1, 256, 1.5, []int{}} {
		for _, b512ValvePumpFieldKey := range []string{"primary", "secondary", "operation", "field", "mode"} {
			b512ValvePumpFieldInput := b512ValvePumpFieldInputTest(5, "request", map[string]any{"channel0": []int{0}})
			b512ValvePumpFieldInput["arguments"] = map[string]any{"external_modes": []any{8, 0, 8}}
			switch b512ValvePumpFieldKey {
			case "field":
				b512ValvePumpFieldInput["fields"].(map[string]any)["channel0"] = []any{b512ValvePumpFieldBad}
			case "mode":
				b512ValvePumpFieldInput["arguments"].(map[string]any)["external_modes"] = []any{8, b512ValvePumpFieldBad, 8}
			default:
				b512ValvePumpFieldInput["context"].(map[string]any)[b512ValvePumpFieldKey] = b512ValvePumpFieldBad
			}
			b512ValvePumpFieldBytes, _ := json.Marshal(b512ValvePumpFieldInput)
			if b512ValvePumpFieldOutput, b512ValvePumpFieldErr := B512ValvePumpFieldEvaluate(b512ValvePumpFieldBytes); b512ValvePumpFieldErr == nil || b512ValvePumpFieldOutput != nil {
				b512ValvePumpFieldT.Fatalf("accepted %s", b512ValvePumpFieldBytes)
			}
		}
	}
	for _, b512ValvePumpFieldArgs := range []any{nil, map[string]any{}, map[string]any{"external_modes": []int{8, 0}}, map[string]any{"external_modes": []int{8, 0, 8, 8}}, map[string]any{"external_modes": []int{8, 0, 8}, "other": 0}} {
		b512ValvePumpFieldInput := b512ValvePumpFieldInputTest(5, "request", map[string]any{})
		if b512ValvePumpFieldArgs != nil {
			b512ValvePumpFieldInput["arguments"] = b512ValvePumpFieldArgs
		}
		b512ValvePumpFieldBytes, _ := json.Marshal(b512ValvePumpFieldInput)
		if b512ValvePumpFieldOutput, b512ValvePumpFieldErr := B512ValvePumpFieldEvaluate(b512ValvePumpFieldBytes); b512ValvePumpFieldErr == nil || b512ValvePumpFieldOutput != nil {
			b512ValvePumpFieldT.Fatal("accepted invalid modes")
		}
	}
}

func TestAb23Aq23Memory(b512ValvePumpFieldT *testing.T) {
	b512ValvePumpFieldInput := []byte(`{"context":{"primary":181,"secondary":18,"operation":0,"direction":"request","view":"actuator-fields"},"fields":{"setpoint":[0],"opaque":[]},"extra":{"nil":null,"zero":0,"false":false,"empty":[]}}`)
	b512ValvePumpFieldOriginal := append([]byte(nil), b512ValvePumpFieldInput...)
	b512ValvePumpFieldFirst, b512ValvePumpFieldErr := B512ValvePumpFieldEvaluate(b512ValvePumpFieldInput)
	if b512ValvePumpFieldErr != nil {
		b512ValvePumpFieldT.Fatal(b512ValvePumpFieldErr)
	}
	b512ValvePumpFieldSaved := append([]byte(nil), b512ValvePumpFieldFirst...)
	for b512ValvePumpFieldIndex := range b512ValvePumpFieldInput {
		b512ValvePumpFieldInput[b512ValvePumpFieldIndex] = ' '
	}
	if !bytes.Equal(b512ValvePumpFieldFirst, b512ValvePumpFieldSaved) {
		b512ValvePumpFieldT.Fatal("input alias")
	}
	for b512ValvePumpFieldIndex := range b512ValvePumpFieldFirst {
		b512ValvePumpFieldFirst[b512ValvePumpFieldIndex] = ' '
	}
	b512ValvePumpFieldLater, b512ValvePumpFieldErr := B512ValvePumpFieldEvaluate(b512ValvePumpFieldOriginal)
	if b512ValvePumpFieldErr != nil || !bytes.Equal(b512ValvePumpFieldLater, b512ValvePumpFieldSaved) {
		b512ValvePumpFieldT.Fatal("output alias")
	}
	if !strings.Contains(string(b512ValvePumpFieldLater), `"opaque":[]`) || !strings.Contains(string(b512ValvePumpFieldLater), `"nil":null`) || !strings.Contains(string(b512ValvePumpFieldLater), `"false":false`) {
		b512ValvePumpFieldT.Fatal("explicit values lost")
	}
}

func b512ValvePumpFieldEqualTest(b512ValvePumpFieldLeft any, b512ValvePumpFieldRight any) bool {
	b512ValvePumpFieldLeftBytes, _ := json.Marshal(b512ValvePumpFieldLeft)
	b512ValvePumpFieldRightBytes, _ := json.Marshal(b512ValvePumpFieldRight)
	return bytes.Equal(b512ValvePumpFieldLeftBytes, b512ValvePumpFieldRightBytes)
}

func TestAb23Aq23RetentionAndKinds(b512ValvePumpFieldT *testing.T) {
	b512ValvePumpFieldInput := b512ValvePumpFieldInputTest(0, "request", map[string]any{"setpoint": []int{0}, "opaque": []int{}})
	for _, b512ValvePumpFieldKey := range []string{"state", "results", "qualified", "physical", "input", "unknown"} {
		b512ValvePumpFieldInput[b512ValvePumpFieldKey] = map[string]any{"null": nil, "false": false, "zero": 0, "empty": []int{}}
	}
	b512ValvePumpFieldOutput := b512ValvePumpFieldEvaluateTest(b512ValvePumpFieldT, b512ValvePumpFieldInput)
	if !b512ValvePumpFieldEqualTest(b512ValvePumpFieldOutput["input"], b512ValvePumpFieldInput) {
		b512ValvePumpFieldT.Fatal("reserved input keys lost")
	}
	for _, b512ValvePumpFieldBad := range []any{nil, false, "0", 0, map[string]any{}} {
		b512ValvePumpFieldInput := b512ValvePumpFieldInputTest(0, "request", map[string]any{"setpoint": b512ValvePumpFieldBad})
		b512ValvePumpFieldBytes, _ := json.Marshal(b512ValvePumpFieldInput)
		if b512ValvePumpFieldResult, b512ValvePumpFieldErr := B512ValvePumpFieldEvaluate(b512ValvePumpFieldBytes); b512ValvePumpFieldErr == nil || b512ValvePumpFieldResult != nil {
			b512ValvePumpFieldT.Fatal("invalid field kind")
		}
	}
	for _, b512ValvePumpFieldKey := range []string{"direction", "view"} {
		for _, b512ValvePumpFieldBad := range []any{nil, false, 0, []any{}} {
			b512ValvePumpFieldInput := b512ValvePumpFieldInputTest(0, "request", map[string]any{})
			b512ValvePumpFieldInput["context"].(map[string]any)[b512ValvePumpFieldKey] = b512ValvePumpFieldBad
			b512ValvePumpFieldBytes, _ := json.Marshal(b512ValvePumpFieldInput)
			if b512ValvePumpFieldResult, b512ValvePumpFieldErr := B512ValvePumpFieldEvaluate(b512ValvePumpFieldBytes); b512ValvePumpFieldErr == nil || b512ValvePumpFieldResult != nil {
				b512ValvePumpFieldT.Fatal("invalid context kind")
			}
		}
	}
	for _, b512ValvePumpFieldJSON := range []string{`{"context":{},"fields":{},"extra":1e999}`, `{"context":{},"fields":{},"arguments":{"x":1,"x":2}}`, `{"context":{},"fields":{},"fields":{}}`, `{"context":{},"fields":{"x":[],"x":[]}}`} {
		if b512ValvePumpFieldResult, b512ValvePumpFieldErr := B512ValvePumpFieldEvaluate([]byte(b512ValvePumpFieldJSON)); b512ValvePumpFieldErr == nil || b512ValvePumpFieldResult != nil {
			b512ValvePumpFieldT.Fatal("invalid JSON accepted")
		}
	}
}

func TestAb23Aq23RegisteredTemperatureRaw(b512ValvePumpFieldT *testing.T) {
	b512ValvePumpFieldInput := b512ValvePumpFieldInputTest(3, "positive", map[string]any{
		"status": []int{0}, "temperature_raw": []int{0x34, 0x12},
	})
	b512ValvePumpFieldOutput := b512ValvePumpFieldEvaluateTest(b512ValvePumpFieldT, b512ValvePumpFieldInput)
	b512ValvePumpFieldResults := b512ValvePumpFieldOutput["results"].(map[string]any)
	b512ValvePumpFieldResult := b512ValvePumpFieldResults["temperature_raw"].(map[string]any)
	if b512ValvePumpFieldOutput["state"] != "complete" || b512ValvePumpFieldResult["state"] != "complete" || b512ValvePumpFieldResult["number"] != float64(4660) || b512ValvePumpFieldResult["label"] != "raw" || b512ValvePumpFieldResult["meaning"] != "Hypothesis" || b512ValvePumpFieldResult["calculation"] != "Proven" {
		b512ValvePumpFieldT.Fatal(b512ValvePumpFieldOutput)
	}
	if !b512ValvePumpFieldEqualTest(b512ValvePumpFieldResult["raw"], []int{0x34, 0x12}) {
		b512ValvePumpFieldT.Fatal("registered raw bytes changed")
	}
	if _, b512ValvePumpFieldExists := b512ValvePumpFieldResults["value"]; b512ValvePumpFieldExists {
		b512ValvePumpFieldT.Fatal("synthesized unregistered value field")
	}
	b512ValvePumpFieldInput["fields"].(map[string]any)["value"] = []int{0xff, 0xff}
	b512ValvePumpFieldOutput = b512ValvePumpFieldEvaluateTest(b512ValvePumpFieldT, b512ValvePumpFieldInput)
	b512ValvePumpFieldResults = b512ValvePumpFieldOutput["results"].(map[string]any)
	if b512ValvePumpFieldOutput["state"] != "complete" || b512ValvePumpFieldResults["temperature_raw"].(map[string]any)["number"] != float64(4660) || !b512ValvePumpFieldEqualTest(b512ValvePumpFieldResults["value"], map[string]any{"raw": []int{255, 255}, "state": "Unknown"}) {
		b512ValvePumpFieldT.Fatal("extra field changed registered interpretation", b512ValvePumpFieldOutput)
	}
}
