package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

type b505OverrideFieldCase struct {
	B505OverrideFieldID       string          `json:"id"`
	B505OverrideFieldInput    json.RawMessage `json:"input"`
	B505OverrideFieldExpected map[string]any  `json:"expected"`
}

func b505OverrideFieldCases(b505OverrideFieldT *testing.T) []b505OverrideFieldCase {
	b505OverrideFieldT.Helper()
	b505OverrideFieldData, b505OverrideFieldErr := os.ReadFile("testdata/ab03aq03_b505_override_fields.json")
	if b505OverrideFieldErr != nil {
		b505OverrideFieldT.Fatal(b505OverrideFieldErr)
	}
	var b505OverrideFieldCases []b505OverrideFieldCase
	if b505OverrideFieldErr := json.Unmarshal(b505OverrideFieldData, &b505OverrideFieldCases); b505OverrideFieldErr != nil {
		b505OverrideFieldT.Fatal(b505OverrideFieldErr)
	}
	if len(b505OverrideFieldCases) != 7 {
		b505OverrideFieldT.Fatal("fixture case count")
	}
	return b505OverrideFieldCases
}

func b505OverrideFieldSubset(b505OverrideFieldT *testing.T, b505OverrideFieldExpected, b505OverrideFieldActual any) {
	b505OverrideFieldT.Helper()
	if b505OverrideFieldMap, b505OverrideFieldOK := b505OverrideFieldExpected.(map[string]any); b505OverrideFieldOK {
		b505OverrideFieldGot, b505OverrideFieldOK := b505OverrideFieldActual.(map[string]any)
		if !b505OverrideFieldOK {
			b505OverrideFieldT.Fatalf("expected object, got %v", b505OverrideFieldActual)
		}
		for b505OverrideFieldKey, b505OverrideFieldValue := range b505OverrideFieldMap {
			b505OverrideFieldSubset(b505OverrideFieldT, b505OverrideFieldValue, b505OverrideFieldGot[b505OverrideFieldKey])
		}
		return
	}
	if b505OverrideFieldNumber, b505OverrideFieldOK := b505OverrideFieldExpected.(int); b505OverrideFieldOK {
		b505OverrideFieldExpected = float64(b505OverrideFieldNumber)
	}
	if !reflect.DeepEqual(b505OverrideFieldExpected, b505OverrideFieldActual) {
		b505OverrideFieldT.Fatalf("want %#v got %#v", b505OverrideFieldExpected, b505OverrideFieldActual)
	}
}

func b505OverrideFieldRun(b505OverrideFieldT *testing.T, b505OverrideFieldInput any) map[string]any {
	b505OverrideFieldT.Helper()
	b505OverrideFieldRaw, b505OverrideFieldErr := json.Marshal(b505OverrideFieldInput)
	if b505OverrideFieldErr != nil {
		b505OverrideFieldT.Fatal(b505OverrideFieldErr)
	}
	b505OverrideFieldOutput, b505OverrideFieldErr := B505OverrideFieldEvaluate(b505OverrideFieldRaw)
	if b505OverrideFieldErr != nil {
		b505OverrideFieldT.Fatal(b505OverrideFieldErr)
	}
	var b505OverrideFieldGot map[string]any
	if b505OverrideFieldErr := json.Unmarshal(b505OverrideFieldOutput, &b505OverrideFieldGot); b505OverrideFieldErr != nil {
		b505OverrideFieldT.Fatal(b505OverrideFieldErr)
	}
	return b505OverrideFieldGot
}

func b505OverrideFieldInput(b505OverrideFieldT *testing.T, b505OverrideFieldIndex int) map[string]any {
	b505OverrideFieldT.Helper()
	var b505OverrideFieldInput map[string]any
	if b505OverrideFieldErr := json.Unmarshal(b505OverrideFieldCases(b505OverrideFieldT)[b505OverrideFieldIndex].B505OverrideFieldInput, &b505OverrideFieldInput); b505OverrideFieldErr != nil {
		b505OverrideFieldT.Fatal(b505OverrideFieldErr)
	}
	return b505OverrideFieldInput
}

func TestAb03Aq03Fixture(b505OverrideFieldT *testing.T) {
	for _, b505OverrideFieldCase := range b505OverrideFieldCases(b505OverrideFieldT) {
		b505OverrideFieldT.Run(b505OverrideFieldCase.B505OverrideFieldID, func(b505OverrideFieldT *testing.T) {
			b505OverrideFieldGot := b505OverrideFieldRun(b505OverrideFieldT, b505OverrideFieldCase.B505OverrideFieldInput)
			b505OverrideFieldSubset(b505OverrideFieldT, b505OverrideFieldCase.B505OverrideFieldExpected, b505OverrideFieldGot)
			var b505OverrideFieldOriginal map[string]any
			if b505OverrideFieldErr := json.Unmarshal(b505OverrideFieldCase.B505OverrideFieldInput, &b505OverrideFieldOriginal); b505OverrideFieldErr != nil {
				b505OverrideFieldT.Fatal(b505OverrideFieldErr)
			}
			for b505OverrideFieldKey, b505OverrideFieldValue := range b505OverrideFieldOriginal {
				if !reflect.DeepEqual(b505OverrideFieldValue, b505OverrideFieldGot[b505OverrideFieldKey]) {
					b505OverrideFieldT.Fatalf("input retention: %s", b505OverrideFieldKey)
				}
			}
			if b505OverrideFieldOriginal["context"].(map[string]any)["operation"] == float64(78) {
				if _, b505OverrideFieldOK := b505OverrideFieldGot["candidate_validity_seconds"]; b505OverrideFieldOK {
					b505OverrideFieldT.Fatal("priority inherited duration")
				}
			}
			if b505OverrideFieldOriginal["context"].(map[string]any)["operation"] == float64(97) {
				if _, b505OverrideFieldOK := b505OverrideFieldGot["results"].(map[string]any)["sensor_status"].(map[string]any)["label"]; b505OverrideFieldOK {
					b505OverrideFieldT.Fatal("OP61 sensor label")
				}
			}
		})
	}
}

func TestAb03Aq03Identity(b505OverrideFieldT *testing.T) {
	for _, b505OverrideFieldCase := range []struct {
		b505OverrideFieldKey   string
		b505OverrideFieldValue any
	}{{"primary", 180}, {"secondary", 6}, {"operation", nil}, {"operation", 0}, {"operation", 50}, {"direction", "response"}, {"view", "other"}, {"family", "other-family"}, {"variant", "other-variant"}} {
		b505OverrideFieldT.Run(b505OverrideFieldCase.b505OverrideFieldKey+"/"+strings.TrimSpace(string(b505OverrideFieldMustJSON(b505OverrideFieldCase.b505OverrideFieldValue))), func(b505OverrideFieldT *testing.T) {
			b505OverrideFieldInput := b505OverrideFieldInput(b505OverrideFieldT, 0)
			b505OverrideFieldContext := b505OverrideFieldInput["context"].(map[string]any)
			if b505OverrideFieldCase.b505OverrideFieldValue == nil {
				delete(b505OverrideFieldContext, b505OverrideFieldCase.b505OverrideFieldKey)
			} else {
				b505OverrideFieldContext[b505OverrideFieldCase.b505OverrideFieldKey] = b505OverrideFieldCase.b505OverrideFieldValue
			}
			b505OverrideFieldGot := b505OverrideFieldRun(b505OverrideFieldT, b505OverrideFieldInput)
			b505OverrideFieldSubset(b505OverrideFieldT, map[string]any{"context": b505OverrideFieldContext, "fields": b505OverrideFieldInput["fields"], "qualified": false, "physical": "Unknown"}, b505OverrideFieldGot)
			if b505OverrideFieldCase.b505OverrideFieldKey == "family" || b505OverrideFieldCase.b505OverrideFieldKey == "variant" {
				if b505OverrideFieldGot["grade"] != "Hypothesis" {
					b505OverrideFieldT.Fatal("opaque identity changed grade")
				}
				return
			}
			if b505OverrideFieldGot["state"] != "Unknown" {
				b505OverrideFieldT.Fatal("unsupported identity admitted")
			}
			for _, b505OverrideFieldResult := range b505OverrideFieldGot["results"].(map[string]any) {
				if len(b505OverrideFieldResult.(map[string]any)) != 2 {
					b505OverrideFieldT.Fatal("unsupported conversion")
				}
			}
			if _, b505OverrideFieldOK := b505OverrideFieldGot["candidate_validity_seconds"]; b505OverrideFieldOK {
				b505OverrideFieldT.Fatal("unsupported duration")
			}
		})
	}
}

func b505OverrideFieldMustJSON(b505OverrideFieldValue any) []byte {
	b505OverrideFieldRaw, _ := json.Marshal(b505OverrideFieldValue)
	return b505OverrideFieldRaw
}

func TestAb03Aq03WidthsAndArguments(b505OverrideFieldT *testing.T) {
	for _, b505OverrideFieldKey := range []string{"sensor_status", "remote_control", "priority"} {
		for b505OverrideFieldIndex, b505OverrideFieldValue := range []any{nil, []any{}, []any{0, 1}, []any{0}} {
			b505OverrideFieldInput := b505OverrideFieldInput(b505OverrideFieldT, 0)
			if b505OverrideFieldKey == "priority" {
				b505OverrideFieldInput = b505OverrideFieldInputFromPriority(b505OverrideFieldT)
			}
			b505OverrideFieldFields := b505OverrideFieldInput["fields"].(map[string]any)
			if b505OverrideFieldValue == nil {
				delete(b505OverrideFieldFields, b505OverrideFieldKey)
			} else {
				b505OverrideFieldFields[b505OverrideFieldKey] = b505OverrideFieldValue
			}
			b505OverrideFieldGot := b505OverrideFieldRun(b505OverrideFieldT, b505OverrideFieldInput)
			b505OverrideFieldResult := b505OverrideFieldGot["results"].(map[string]any)[b505OverrideFieldKey].(map[string]any)
			if b505OverrideFieldIndex < 3 {
				if b505OverrideFieldGot["state"] != "partial" || b505OverrideFieldResult["state"] != "incomplete" {
					b505OverrideFieldT.Fatal("width admitted")
				}
				if _, b505OverrideFieldOK := b505OverrideFieldResult["number"]; b505OverrideFieldOK {
					b505OverrideFieldT.Fatal("incomplete number")
				}
				if b505OverrideFieldValue == nil {
					if _, b505OverrideFieldOK := b505OverrideFieldResult["raw"]; b505OverrideFieldOK {
						b505OverrideFieldT.Fatal("missing became present")
					}
				}
			} else if b505OverrideFieldGot["state"] != "complete" {
				b505OverrideFieldT.Fatal("zero incomplete")
			}
			if b505OverrideFieldKey == "sensor_status" {
				b505OverrideFieldSubset(b505OverrideFieldT, map[string]any{"label": "enabled", "state": "complete"}, b505OverrideFieldGot["results"].(map[string]any)["remote_control"])
			}
			if b505OverrideFieldKey == "remote_control" {
				b505OverrideFieldSubset(b505OverrideFieldT, map[string]any{"label": "short-circuit", "state": "complete"}, b505OverrideFieldGot["results"].(map[string]any)["sensor_status"])
			}
		}
	}
	for _, b505OverrideFieldValue := range []int{0, 25, 50, 75, 100, 26, 255} {
		b505OverrideFieldInput := b505OverrideFieldInput(b505OverrideFieldT, 0)
		b505OverrideFieldInput["arguments"].(map[string]any)["display_percent"] = b505OverrideFieldValue
		b505OverrideFieldGot := b505OverrideFieldRun(b505OverrideFieldT, b505OverrideFieldInput)["results"].(map[string]any)["display_percent"].(map[string]any)
		b505OverrideFieldWant := b505OverrideFieldValue == 0 || b505OverrideFieldValue == 25 || b505OverrideFieldValue == 50 || b505OverrideFieldValue == 75 || b505OverrideFieldValue == 100
		if b505OverrideFieldGot["number"] != float64(b505OverrideFieldValue) || b505OverrideFieldGot["candidate_domain"] != b505OverrideFieldWant {
			b505OverrideFieldT.Fatal("display rounded or domain changed")
		}
	}
}

func b505OverrideFieldInputFromPriority(b505OverrideFieldT *testing.T) map[string]any {
	return b505OverrideFieldInput(b505OverrideFieldT, 4)
}

func TestAb03Aq03Reject(b505OverrideFieldT *testing.T) {
	for _, b505OverrideFieldRaw := range []string{``, `null`, `[]`, `{`, `{} {}`, `{"context":{},"context":{}}`, `{"context":{"x":0,"x":1},"fields":{}}`, `{"context":{},"fields":{},"x":[{"y":1,"y":2}]}`, `{"context":{},"fields":{},"x":NaN}`, `{"context":{},"fields":{},"x":1e999}`, `{"context":{},"fields":{},"arguments":{"x":1,"x":2}}`, `{"context":{},"fields":{"x":[0],"x":[1]}}`} {
		b505OverrideFieldOutput, b505OverrideFieldErr := B505OverrideFieldEvaluate([]byte(b505OverrideFieldRaw))
		if b505OverrideFieldErr == nil || b505OverrideFieldOutput != nil {
			b505OverrideFieldT.Fatalf("accepted malformed %s", b505OverrideFieldRaw)
		}
	}
	for _, b505OverrideFieldPath := range []string{"primary", "secondary", "operation", "display_percent", "sensor_status", "remote_control", "priority"} {
		for _, b505OverrideFieldValue := range []any{nil, true, "1", -1, 256, 1.5} {
			b505OverrideFieldInput := b505OverrideFieldInput(b505OverrideFieldT, 0)
			if b505OverrideFieldPath == "priority" {
				b505OverrideFieldInput = b505OverrideFieldInputFromPriority(b505OverrideFieldT)
			}
			switch b505OverrideFieldPath {
			case "display_percent":
				b505OverrideFieldInput["arguments"].(map[string]any)[b505OverrideFieldPath] = b505OverrideFieldValue
			case "sensor_status", "remote_control", "priority":
				b505OverrideFieldInput["fields"].(map[string]any)[b505OverrideFieldPath] = []any{b505OverrideFieldValue}
			default:
				b505OverrideFieldInput["context"].(map[string]any)[b505OverrideFieldPath] = b505OverrideFieldValue
			}
			b505OverrideFieldOutput, b505OverrideFieldErr := B505OverrideFieldEvaluate(b505OverrideFieldMustJSON(b505OverrideFieldInput))
			if b505OverrideFieldErr == nil || b505OverrideFieldOutput != nil {
				b505OverrideFieldT.Fatalf("accepted %s=%v", b505OverrideFieldPath, b505OverrideFieldValue)
			}
		}
	}
	for _, b505OverrideFieldKey := range []string{"context", "fields", "arguments"} {
		for _, b505OverrideFieldValue := range []any{nil, true, 1, "x", []any{}} {
			b505OverrideFieldInput := b505OverrideFieldInput(b505OverrideFieldT, 0)
			b505OverrideFieldInput[b505OverrideFieldKey] = b505OverrideFieldValue
			if b505OverrideFieldOutput, b505OverrideFieldErr := B505OverrideFieldEvaluate(b505OverrideFieldMustJSON(b505OverrideFieldInput)); b505OverrideFieldErr == nil || b505OverrideFieldOutput != nil {
				b505OverrideFieldT.Fatal("object kind accepted")
			}
		}
	}
	b505OverrideFieldDeep := `{"context":{},"fields":{},"x":` + strings.Repeat("[", 3000000) + strings.Repeat("]", 3000000) + `}`
	if b505OverrideFieldOutput, b505OverrideFieldErr := B505OverrideFieldEvaluate([]byte(b505OverrideFieldDeep)); b505OverrideFieldErr == nil || b505OverrideFieldOutput != nil {
		b505OverrideFieldT.Fatal("deep nesting accepted")
	}
	if b505OverrideFieldOutput, b505OverrideFieldErr := B505OverrideFieldEvaluate([]byte(`{"context":{},"fields":{},"arguments":{"display_percent":-0}}`)); b505OverrideFieldErr == nil || b505OverrideFieldOutput != nil {
		b505OverrideFieldT.Fatal("negative zero accepted")
	}
	for _, b505OverrideFieldKey := range []string{"physical", "qualified", "state", "results", "grade", "candidate_validity_seconds"} {
		b505OverrideFieldInput := b505OverrideFieldInput(b505OverrideFieldT, 0)
		b505OverrideFieldInput[b505OverrideFieldKey] = "x"
		if b505OverrideFieldOutput, b505OverrideFieldErr := B505OverrideFieldEvaluate(b505OverrideFieldMustJSON(b505OverrideFieldInput)); b505OverrideFieldErr == nil || b505OverrideFieldOutput != nil {
			b505OverrideFieldT.Fatalf("colliding key %s accepted", b505OverrideFieldKey)
		}
	}
}

func TestAb03Aq03OptionalDisplay(b505OverrideFieldT *testing.T) {
	b505OverrideFieldInput := b505OverrideFieldInput(b505OverrideFieldT, 0)
	delete(b505OverrideFieldInput, "arguments")
	b505OverrideFieldGot := b505OverrideFieldRun(b505OverrideFieldT, b505OverrideFieldInput)
	b505OverrideFieldResults := b505OverrideFieldGot["results"].(map[string]any)
	if _, b505OverrideFieldOK := b505OverrideFieldResults["display_percent"]; b505OverrideFieldOK {
		b505OverrideFieldT.Fatal("display result invented")
	}
	if b505OverrideFieldGot["state"] != "complete" || b505OverrideFieldResults["sensor_status"].(map[string]any)["label"] == nil {
		b505OverrideFieldT.Fatal("siblings unavailable without display_percent")
	}
}

func TestAb03Aq03Memory(b505OverrideFieldT *testing.T) {
	b505OverrideFieldInput := append(json.RawMessage(nil), b505OverrideFieldCases(b505OverrideFieldT)[0].B505OverrideFieldInput...)
	b505OverrideFieldOutput, b505OverrideFieldErr := B505OverrideFieldEvaluate(b505OverrideFieldInput)
	if b505OverrideFieldErr != nil {
		b505OverrideFieldT.Fatal(b505OverrideFieldErr)
	}
	b505OverrideFieldControl := append([]byte(nil), b505OverrideFieldOutput...)
	for b505OverrideFieldIndex := range b505OverrideFieldInput {
		b505OverrideFieldInput[b505OverrideFieldIndex] = ' '
	}
	if !reflect.DeepEqual([]byte(b505OverrideFieldOutput), b505OverrideFieldControl) {
		b505OverrideFieldT.Fatal("input aliases output")
	}
	for b505OverrideFieldIndex := range b505OverrideFieldOutput {
		b505OverrideFieldOutput[b505OverrideFieldIndex] = ' '
	}
	b505OverrideFieldNext, b505OverrideFieldErr := B505OverrideFieldEvaluate(b505OverrideFieldCases(b505OverrideFieldT)[0].B505OverrideFieldInput)
	if b505OverrideFieldErr != nil || !reflect.DeepEqual([]byte(b505OverrideFieldNext), b505OverrideFieldControl) {
		b505OverrideFieldT.Fatal("call shared storage")
	}
}

func TestAb03Aq03OpaqueAndKinds(b505OverrideFieldT *testing.T) {
	for _, b505OverrideFieldValue := range []any{nil, []any{}, []any{0}, []any{85, 170}} {
		b505OverrideFieldInput := b505OverrideFieldInput(b505OverrideFieldT, 3)
		if b505OverrideFieldValue == nil {
			delete(b505OverrideFieldInput["fields"].(map[string]any), "dhw_setpoint")
		} else {
			b505OverrideFieldInput["fields"].(map[string]any)["dhw_setpoint"] = b505OverrideFieldValue
		}
		b505OverrideFieldGot := b505OverrideFieldRun(b505OverrideFieldT, b505OverrideFieldInput)
		var b505OverrideFieldFields any
		if b505OverrideFieldErr := json.Unmarshal(b505OverrideFieldMustJSON(b505OverrideFieldInput["fields"]), &b505OverrideFieldFields); b505OverrideFieldErr != nil {
			b505OverrideFieldT.Fatal(b505OverrideFieldErr)
		}
		if !reflect.DeepEqual(b505OverrideFieldFields, b505OverrideFieldGot["fields"]) {
			b505OverrideFieldT.Fatal("opaque raw changed")
		}
		if b505OverrideFieldValue != nil {
			b505OverrideFieldResult := b505OverrideFieldGot["results"].(map[string]any)["dhw_setpoint"].(map[string]any)
			if len(b505OverrideFieldResult) != 2 || b505OverrideFieldResult["state"] != "Unknown" {
				b505OverrideFieldT.Fatal("opaque conversion")
			}
		}
	}
	for _, b505OverrideFieldKey := range []string{"direction", "view", "family", "variant"} {
		for _, b505OverrideFieldValue := range []any{nil, false, 1, []any{}, map[string]any{}} {
			b505OverrideFieldInput := b505OverrideFieldInput(b505OverrideFieldT, 0)
			b505OverrideFieldInput["context"].(map[string]any)[b505OverrideFieldKey] = b505OverrideFieldValue
			if b505OverrideFieldOutput, b505OverrideFieldErr := B505OverrideFieldEvaluate(b505OverrideFieldMustJSON(b505OverrideFieldInput)); b505OverrideFieldErr == nil || b505OverrideFieldOutput != nil {
				b505OverrideFieldT.Fatal("identifier kind admitted")
			}
		}
	}
	for _, b505OverrideFieldKey := range []string{"sensor_status", "remote_control", "priority", "dhw_setpoint"} {
		for _, b505OverrideFieldValue := range []any{nil, false, 0, "x", map[string]any{}} {
			b505OverrideFieldInput := b505OverrideFieldInput(b505OverrideFieldT, 0)
			b505OverrideFieldInput["fields"].(map[string]any)[b505OverrideFieldKey] = b505OverrideFieldValue
			if b505OverrideFieldOutput, b505OverrideFieldErr := B505OverrideFieldEvaluate(b505OverrideFieldMustJSON(b505OverrideFieldInput)); b505OverrideFieldErr == nil || b505OverrideFieldOutput != nil {
				b505OverrideFieldT.Fatal("field kind admitted")
			}
		}
	}
	for _, b505OverrideFieldIndex := range []int{3, 4} {
		b505OverrideFieldInput := b505OverrideFieldInput(b505OverrideFieldT, b505OverrideFieldIndex)
		b505OverrideFieldInput["arguments"] = map[string]any{"display_percent": 25}
		if b505OverrideFieldOutput, b505OverrideFieldErr := B505OverrideFieldEvaluate(b505OverrideFieldMustJSON(b505OverrideFieldInput)); b505OverrideFieldErr == nil || b505OverrideFieldOutput != nil {
			b505OverrideFieldT.Fatal("unmodeled argument admitted")
		}
	}
}

func TestAb03Aq03PresenceAndIsolation(b505OverrideFieldT *testing.T) {
	for _, b505OverrideFieldKey := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
		for _, b505OverrideFieldPresent := range []bool{false, true} {
			b505OverrideFieldInput := b505OverrideFieldInput(b505OverrideFieldT, 0)
			b505OverrideFieldContext := b505OverrideFieldInput["context"].(map[string]any)
			if !b505OverrideFieldPresent {
				delete(b505OverrideFieldContext, b505OverrideFieldKey)
			} else if b505OverrideFieldKey == "primary" || b505OverrideFieldKey == "secondary" || b505OverrideFieldKey == "operation" {
				b505OverrideFieldContext[b505OverrideFieldKey] = float64(0)
			} else {
				b505OverrideFieldContext[b505OverrideFieldKey] = ""
			}
			b505OverrideFieldGot := b505OverrideFieldRun(b505OverrideFieldT, b505OverrideFieldInput)
			if !reflect.DeepEqual(b505OverrideFieldContext, b505OverrideFieldGot["context"]) {
				b505OverrideFieldT.Fatal("context presence changed")
			}
			if b505OverrideFieldKey != "family" && b505OverrideFieldKey != "variant" && b505OverrideFieldGot["state"] != "Unknown" {
				b505OverrideFieldT.Fatal("absent or zero identity admitted")
			}
			if b505OverrideFieldGot["qualified"] != false {
				b505OverrideFieldT.Fatal("opaque qualification")
			}
		}
	}
	for _, b505OverrideFieldValue := range []any{[]any{}, []any{0}, []any{255}, []any{25, 26}} {
		b505OverrideFieldInput := b505OverrideFieldInput(b505OverrideFieldT, 0)
		b505OverrideFieldInput["fields"].(map[string]any)["pump_power"] = b505OverrideFieldValue
		b505OverrideFieldGot := b505OverrideFieldRun(b505OverrideFieldT, b505OverrideFieldInput)
		b505OverrideFieldResult := b505OverrideFieldGot["results"].(map[string]any)["pump_power"].(map[string]any)
		if len(b505OverrideFieldResult) != 2 || b505OverrideFieldResult["state"] != "Unknown" {
			b505OverrideFieldT.Fatal("raw pump converted")
		}
		b505OverrideFieldSubset(b505OverrideFieldT, map[string]any{"number": 25, "candidate_domain": true}, b505OverrideFieldGot["results"].(map[string]any)["display_percent"])
	}
}
