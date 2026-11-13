package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func b505ConditionalRoleDecode(b505ConditionalRoleT *testing.T, b505ConditionalRoleInput []byte) map[string]any {
	b505ConditionalRoleT.Helper()
	b505ConditionalRoleOutput, b505ConditionalRoleError := B505ConditionalRoleEvaluate(b505ConditionalRoleInput)
	if b505ConditionalRoleError != nil {
		b505ConditionalRoleT.Fatal(b505ConditionalRoleError)
	}
	var b505ConditionalRoleObject map[string]any
	if b505ConditionalRoleError = json.Unmarshal(b505ConditionalRoleOutput, &b505ConditionalRoleObject); b505ConditionalRoleError != nil {
		b505ConditionalRoleT.Fatal(b505ConditionalRoleError)
	}
	return b505ConditionalRoleObject
}

func b505ConditionalRoleInput(b505ConditionalRoleView string, b505ConditionalRoleOperation int, b505ConditionalRoleFields, b505ConditionalRoleArguments map[string]any) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 5, "operation": b505ConditionalRoleOperation, "direction": "request", "view": b505ConditionalRoleView, "family": "synthetic-family", "variant": "synthetic-variant"}, "fields": b505ConditionalRoleFields, "arguments": b505ConditionalRoleArguments}
}

func b505ConditionalRoleRun(b505ConditionalRoleT *testing.T, b505ConditionalRoleObject map[string]any) map[string]any {
	b505ConditionalRoleT.Helper()
	b505ConditionalRoleBytes, b505ConditionalRoleError := json.Marshal(b505ConditionalRoleObject)
	if b505ConditionalRoleError != nil {
		b505ConditionalRoleT.Fatal(b505ConditionalRoleError)
	}
	return b505ConditionalRoleDecode(b505ConditionalRoleT, b505ConditionalRoleBytes)
}

func TestAb05Aq05Fixture(b505ConditionalRoleT *testing.T) {
	b505ConditionalRoleData, b505ConditionalRoleError := os.ReadFile("testdata/ab05aq05_b505_conditional_roles.json")
	if b505ConditionalRoleError != nil {
		b505ConditionalRoleT.Fatal(b505ConditionalRoleError)
	}
	var b505ConditionalRoleCases []struct {
		B505ConditionalRoleID       string          `json:"id"`
		B505ConditionalRoleInput    json.RawMessage `json:"input"`
		B505ConditionalRoleKey      string          `json:"key"`
		B505ConditionalRoleExpected map[string]any  `json:"expected"`
	}
	if b505ConditionalRoleError = json.Unmarshal(b505ConditionalRoleData, &b505ConditionalRoleCases); b505ConditionalRoleError != nil {
		b505ConditionalRoleT.Fatal(b505ConditionalRoleError)
	}
	if len(b505ConditionalRoleCases) != 15 {
		b505ConditionalRoleT.Fatalf("fixture case count %d", len(b505ConditionalRoleCases))
	}
	for _, b505ConditionalRoleCase := range b505ConditionalRoleCases {
		b505ConditionalRoleT.Run(b505ConditionalRoleCase.B505ConditionalRoleID, func(b505ConditionalRoleT *testing.T) {
			b505ConditionalRoleOutput := b505ConditionalRoleDecode(b505ConditionalRoleT, b505ConditionalRoleCase.B505ConditionalRoleInput)
			if b505ConditionalRoleOutput["state"] != "complete" || b505ConditionalRoleOutput["physical"] != "Unknown" || b505ConditionalRoleOutput["qualified"] != false {
				b505ConditionalRoleT.Fatal(b505ConditionalRoleOutput)
			}
			b505ConditionalRoleResult := b505ConditionalRoleOutput["results"].(map[string]any)[b505ConditionalRoleCase.B505ConditionalRoleKey].(map[string]any)
			if b505ConditionalRoleResult["grade"] != "Hypothesis" || b505ConditionalRoleResult["state"] != "complete" {
				b505ConditionalRoleT.Fatal(b505ConditionalRoleResult)
			}
			for b505ConditionalRoleKey, b505ConditionalRoleExpected := range b505ConditionalRoleCase.B505ConditionalRoleExpected {
				if !reflect.DeepEqual(b505ConditionalRoleResult[b505ConditionalRoleKey], b505ConditionalRoleExpected) {
					b505ConditionalRoleT.Fatalf("%s: got %v want %v", b505ConditionalRoleKey, b505ConditionalRoleResult[b505ConditionalRoleKey], b505ConditionalRoleExpected)
				}
			}
		})
	}
}

func TestAb05Aq05Identity(b505ConditionalRoleT *testing.T) {
	for _, b505ConditionalRoleAxis := range []string{"primary", "secondary", "operation-absent", "operation-zero", "operation-other", "direction", "view", "family", "variant"} {
		b505ConditionalRoleT.Run(b505ConditionalRoleAxis, func(b505ConditionalRoleT *testing.T) {
			b505ConditionalRoleObject := b505ConditionalRoleInput("dhw-enable", 49, map[string]any{"enabled": []int{1}}, map[string]any{})
			b505ConditionalRoleContext := b505ConditionalRoleObject["context"].(map[string]any)
			switch b505ConditionalRoleAxis {
			case "primary":
				b505ConditionalRoleContext["primary"] = 180
			case "secondary":
				b505ConditionalRoleContext["secondary"] = 6
			case "operation-absent":
				delete(b505ConditionalRoleContext, "operation")
			case "operation-zero":
				b505ConditionalRoleContext["operation"] = 0
			case "operation-other":
				b505ConditionalRoleContext["operation"] = 9
			default:
				b505ConditionalRoleContext[b505ConditionalRoleAxis] = "other"
			}
			b505ConditionalRoleObject["extra"] = map[string]any{"null": nil, "zero": 0, "false": false, "empty": []any{}}
			b505ConditionalRoleOutput := b505ConditionalRoleRun(b505ConditionalRoleT, b505ConditionalRoleObject)
			b505ConditionalRoleState := "Unknown"
			if b505ConditionalRoleAxis == "family" || b505ConditionalRoleAxis == "variant" {
				b505ConditionalRoleState = "complete"
			}
			if b505ConditionalRoleOutput["state"] != b505ConditionalRoleState || b505ConditionalRoleOutput["qualified"] != false {
				b505ConditionalRoleT.Fatal(b505ConditionalRoleOutput)
			}
			b505ConditionalRoleResult := b505ConditionalRoleOutput["results"].(map[string]any)["enabled"].(map[string]any)
			if b505ConditionalRoleState == "Unknown" {
				if len(b505ConditionalRoleResult) != 2 || b505ConditionalRoleResult["state"] != "Unknown" {
					b505ConditionalRoleT.Fatal(b505ConditionalRoleResult)
				}
			}
			b505ConditionalRoleBytes, _ := json.Marshal(b505ConditionalRoleObject)
			var b505ConditionalRoleOriginal map[string]any
			_ = json.Unmarshal(b505ConditionalRoleBytes, &b505ConditionalRoleOriginal)
			for _, b505ConditionalRoleKey := range []string{"context", "fields", "arguments", "extra"} {
				if !reflect.DeepEqual(b505ConditionalRoleOutput[b505ConditionalRoleKey], b505ConditionalRoleOriginal[b505ConditionalRoleKey]) {
					b505ConditionalRoleT.Fatal("lost raw", b505ConditionalRoleKey)
				}
			}
		})
	}
	for _, b505ConditionalRoleOperation := range []int{61, 90} {
		b505ConditionalRoleObject := b505ConditionalRoleInput("unchanged", b505ConditionalRoleOperation, map[string]any{}, map[string]any{"replacement": true, "previous": 17, "requested": 23})
		if b505ConditionalRoleRun(b505ConditionalRoleT, b505ConditionalRoleObject)["state"] != "complete" {
			b505ConditionalRoleT.Fatal("operation alias")
		}
	}
}

func TestAb05Aq05PartialAndOpaque(b505ConditionalRoleT *testing.T) {
	for _, b505ConditionalRoleField := range []string{"enabled", "disable", "start"} {
		for _, b505ConditionalRoleWidth := range []int{-1, 0, 1, 2} {
			b505ConditionalRoleView := "dhw-enable"
			b505ConditionalRoleOperation := 49
			b505ConditionalRoleArgs := map[string]any{}
			if b505ConditionalRoleField == "disable" {
				b505ConditionalRoleView = "parallel-charge"
				b505ConditionalRoleOperation = 24
				b505ConditionalRoleArgs = map[string]any{"role": "mixer", "parallel": true}
			}
			if b505ConditionalRoleField == "start" {
				b505ConditionalRoleView = "programmed-start"
				b505ConditionalRoleOperation = 9
				b505ConditionalRoleArgs = map[string]any{"program_mode": 0}
			}
			b505ConditionalRoleFields := map[string]any{"independent": []int{17}, "empty": []int{}}
			if b505ConditionalRoleWidth >= 0 {
				b505ConditionalRoleFields[b505ConditionalRoleField] = make([]int, b505ConditionalRoleWidth)
			}
			b505ConditionalRoleObject := b505ConditionalRoleInput(b505ConditionalRoleView, b505ConditionalRoleOperation, b505ConditionalRoleFields, b505ConditionalRoleArgs)
			b505ConditionalRoleOutput := b505ConditionalRoleRun(b505ConditionalRoleT, b505ConditionalRoleObject)
			b505ConditionalRoleComplete := b505ConditionalRoleWidth == 1 || (b505ConditionalRoleField == "start" && b505ConditionalRoleWidth >= 0)
			b505ConditionalRoleState := "partial"
			if b505ConditionalRoleComplete {
				b505ConditionalRoleState = "complete"
			}
			if b505ConditionalRoleOutput["state"] != b505ConditionalRoleState {
				b505ConditionalRoleT.Fatal(b505ConditionalRoleField, b505ConditionalRoleWidth, b505ConditionalRoleOutput)
			}
			b505ConditionalRoleResults := b505ConditionalRoleOutput["results"].(map[string]any)
			if !reflect.DeepEqual(b505ConditionalRoleResults["independent"].(map[string]any)["raw"], []any{float64(17)}) || !reflect.DeepEqual(b505ConditionalRoleResults["empty"].(map[string]any)["raw"], []any{}) {
				b505ConditionalRoleT.Fatal("siblings lost")
			}
			b505ConditionalRoleResult := b505ConditionalRoleResults[b505ConditionalRoleField].(map[string]any)
			_, b505ConditionalRoleHasRaw := b505ConditionalRoleResult["raw"]
			if b505ConditionalRoleHasRaw != (b505ConditionalRoleWidth >= 0) {
				b505ConditionalRoleT.Fatal("missing/empty collapsed")
			}
			if !b505ConditionalRoleComplete {
				if _, b505ConditionalRoleExists := b505ConditionalRoleResult["number"]; b505ConditionalRoleExists {
					b505ConditionalRoleT.Fatal("invented number")
				}
			}
		}
	}
}

func TestAb05Aq05Reject(b505ConditionalRoleT *testing.T) {
	for _, b505ConditionalRoleBad := range []string{``, `{`, `[]`, `{} {}`, `{"context":{},"fields":{},"x":{"a":1,"a":2}}`, `{"context":{},"fields":{},"x":[{"a":1,"a":2}]}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{},"fields":{},"x":NaN}`, `{"context":{},"fields":{},"x":1e999}`, `{"context":{},"fields":null}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":{},"arguments":null}`} {
		if b505ConditionalRoleOutput, b505ConditionalRoleError := B505ConditionalRoleEvaluate([]byte(b505ConditionalRoleBad)); b505ConditionalRoleError == nil || b505ConditionalRoleOutput != nil {
			b505ConditionalRoleT.Fatalf("accepted %s", b505ConditionalRoleBad)
		}
	}
	for _, b505ConditionalRoleKey := range []string{"primary", "secondary", "operation", "enabled", "program_mode", "replacement", "previous", "requested", "role", "parallel", "demand", "excess"} {
		b505ConditionalRoleBadValues := []any{nil, "1", true, 1.5, -1, 256}
		if b505ConditionalRoleKey == "replacement" || b505ConditionalRoleKey == "parallel" {
			b505ConditionalRoleBadValues = []any{nil, "true", 0, []any{}}
		}
		if b505ConditionalRoleKey == "previous" || b505ConditionalRoleKey == "requested" {
			b505ConditionalRoleBadValues = []any{[]any{}, map[string]any{}}
		}
		if b505ConditionalRoleKey == "role" {
			b505ConditionalRoleBadValues = []any{nil, 0, false, "other"}
		}
		if b505ConditionalRoleKey == "demand" || b505ConditionalRoleKey == "excess" {
			b505ConditionalRoleBadValues = []any{nil, "1", true, []any{}}
		}
		for _, b505ConditionalRoleBad := range b505ConditionalRoleBadValues {
			b505ConditionalRoleObject := b505ConditionalRoleInput("dhw-enable", 49, map[string]any{"enabled": []int{1}}, map[string]any{})
			b505ConditionalRoleContext := b505ConditionalRoleObject["context"].(map[string]any)
			b505ConditionalRoleArgs := b505ConditionalRoleObject["arguments"].(map[string]any)
			switch b505ConditionalRoleKey {
			case "primary", "secondary", "operation":
				b505ConditionalRoleContext[b505ConditionalRoleKey] = b505ConditionalRoleBad
			case "enabled":
				b505ConditionalRoleObject["fields"] = map[string]any{"enabled": []any{b505ConditionalRoleBad}}
			case "program_mode":
				b505ConditionalRoleContext["view"] = "programmed-start"
				b505ConditionalRoleContext["operation"] = 9
				b505ConditionalRoleArgs[b505ConditionalRoleKey] = b505ConditionalRoleBad
			case "replacement", "previous", "requested":
				b505ConditionalRoleContext["view"] = "unchanged"
				b505ConditionalRoleContext["operation"] = 90
				b505ConditionalRoleArgs["replacement"] = true
				b505ConditionalRoleArgs["previous"] = 17
				b505ConditionalRoleArgs["requested"] = 23
				b505ConditionalRoleArgs[b505ConditionalRoleKey] = b505ConditionalRoleBad
			case "role", "parallel":
				b505ConditionalRoleContext["view"] = "parallel-charge"
				b505ConditionalRoleContext["operation"] = 24
				b505ConditionalRoleArgs["role"] = "mixer"
				b505ConditionalRoleArgs["parallel"] = true
				b505ConditionalRoleArgs[b505ConditionalRoleKey] = b505ConditionalRoleBad
			case "demand", "excess":
				b505ConditionalRoleContext["view"] = "flow-excess"
				b505ConditionalRoleContext["operation"] = 23
				b505ConditionalRoleArgs["demand"] = 40
				b505ConditionalRoleArgs["excess"] = 5
				b505ConditionalRoleArgs[b505ConditionalRoleKey] = b505ConditionalRoleBad
			}
			b505ConditionalRoleBytes, _ := json.Marshal(b505ConditionalRoleObject)
			if b505ConditionalRoleOutput, b505ConditionalRoleError := B505ConditionalRoleEvaluate(b505ConditionalRoleBytes); b505ConditionalRoleError == nil || b505ConditionalRoleOutput != nil {
				b505ConditionalRoleT.Fatalf("accepted %s=%v", b505ConditionalRoleKey, b505ConditionalRoleBad)
			}
		}
	}
	for _, b505ConditionalRoleCase := range []struct {
		B505ConditionalRoleView      string
		B505ConditionalRoleOperation int
		B505ConditionalRoleArgs      map[string]any
	}{{"unchanged", 90, map[string]any{"replacement": true, "previous": 17, "requested": 23}}, {"parallel-charge", 24, map[string]any{"role": "mixer", "parallel": true}}, {"flow-excess", 23, map[string]any{"demand": 40, "excess": 5}}, {"programmed-start", 9, map[string]any{"program_mode": 0}}} {
		for b505ConditionalRoleKey := range b505ConditionalRoleCase.B505ConditionalRoleArgs {
			b505ConditionalRoleArgs := map[string]any{}
			for b505ConditionalRoleName, b505ConditionalRoleValue := range b505ConditionalRoleCase.B505ConditionalRoleArgs {
				if b505ConditionalRoleName != b505ConditionalRoleKey {
					b505ConditionalRoleArgs[b505ConditionalRoleName] = b505ConditionalRoleValue
				}
			}
			b505ConditionalRoleBytes, _ := json.Marshal(b505ConditionalRoleInput(b505ConditionalRoleCase.B505ConditionalRoleView, b505ConditionalRoleCase.B505ConditionalRoleOperation, map[string]any{}, b505ConditionalRoleArgs))
			if b505ConditionalRoleOutput, b505ConditionalRoleError := B505ConditionalRoleEvaluate(b505ConditionalRoleBytes); b505ConditionalRoleError == nil || b505ConditionalRoleOutput != nil {
				b505ConditionalRoleT.Fatal("accepted missing", b505ConditionalRoleKey)
			}
		}
	}
}

func TestAb05Aq05ScalarsAndMemory(b505ConditionalRoleT *testing.T) {
	for _, b505ConditionalRoleValue := range []any{nil, false, 0, "", 17.5} {
		for _, b505ConditionalRoleReplacement := range []bool{false, true} {
			b505ConditionalRoleObject := b505ConditionalRoleInput("unchanged", 90, map[string]any{}, map[string]any{"replacement": b505ConditionalRoleReplacement, "previous": b505ConditionalRoleValue, "requested": b505ConditionalRoleValue})
			b505ConditionalRoleOutput := b505ConditionalRoleRun(b505ConditionalRoleT, b505ConditionalRoleObject)
			b505ConditionalRoleBytes, _ := json.Marshal(b505ConditionalRoleValue)
			var b505ConditionalRoleExpected any
			_ = json.Unmarshal(b505ConditionalRoleBytes, &b505ConditionalRoleExpected)
			if !reflect.DeepEqual(b505ConditionalRoleOutput["results"].(map[string]any)["unchanged"].(map[string]any)["value"], b505ConditionalRoleExpected) {
				b505ConditionalRoleT.Fatal("scalar collapsed")
			}
		}
	}
	b505ConditionalRoleObject := b505ConditionalRoleInput("flow-excess", 23, map[string]any{}, map[string]any{"demand": 0.1, "excess": 0.2})
	if b505ConditionalRoleRun(b505ConditionalRoleT, b505ConditionalRoleObject)["results"].(map[string]any)["flow-excess"].(map[string]any)["number"] != 0.3 {
		b505ConditionalRoleT.Fatal("decimal sum")
	}
	b505ConditionalRoleInputBytes, _ := json.Marshal(b505ConditionalRoleObject)
	b505ConditionalRoleOriginal := append([]byte(nil), b505ConditionalRoleInputBytes...)
	b505ConditionalRoleFirst, b505ConditionalRoleError := B505ConditionalRoleEvaluate(b505ConditionalRoleInputBytes)
	if b505ConditionalRoleError != nil {
		b505ConditionalRoleT.Fatal(b505ConditionalRoleError)
	}
	b505ConditionalRoleSaved := append([]byte(nil), b505ConditionalRoleFirst...)
	for b505ConditionalRoleIndex := range b505ConditionalRoleInputBytes {
		b505ConditionalRoleInputBytes[b505ConditionalRoleIndex] = ' '
	}
	if string(b505ConditionalRoleFirst) != string(b505ConditionalRoleSaved) {
		b505ConditionalRoleT.Fatal("input aliases output")
	}
	for b505ConditionalRoleIndex := range b505ConditionalRoleFirst {
		b505ConditionalRoleFirst[b505ConditionalRoleIndex] = ' '
	}
	b505ConditionalRoleSecond, b505ConditionalRoleError := B505ConditionalRoleEvaluate(b505ConditionalRoleOriginal)
	if b505ConditionalRoleError != nil || string(b505ConditionalRoleSecond) != string(b505ConditionalRoleSaved) {
		b505ConditionalRoleT.Fatal("cross-call storage")
	}
	if strings.Contains(string(b505ConditionalRoleSecond), "replacement_code") {
		b505ConditionalRoleT.Fatal("cross-view contamination")
	}
}

func TestAb05Aq05IndependentInputs(b505ConditionalRoleT *testing.T) {
	for _, b505ConditionalRoleCase := range []struct {
		B505ConditionalRoleDemand float64
		B505ConditionalRoleExcess float64
		B505ConditionalRoleWant   float64
	}{{41, 5, 46}, {40, 6, 46}, {0, 5, 5}, {40, 0, 40}, {-40, 5, -35}} {
		b505ConditionalRoleResult := b505ConditionalRoleRun(b505ConditionalRoleT, b505ConditionalRoleInput("flow-excess", 23, map[string]any{}, map[string]any{"demand": b505ConditionalRoleCase.B505ConditionalRoleDemand, "excess": b505ConditionalRoleCase.B505ConditionalRoleExcess}))["results"].(map[string]any)["flow-excess"].(map[string]any)
		if b505ConditionalRoleResult["number"] != b505ConditionalRoleCase.B505ConditionalRoleWant {
			b505ConditionalRoleT.Fatal(b505ConditionalRoleResult)
		}
	}
	for _, b505ConditionalRoleKey := range []string{"previous", "requested"} {
		b505ConditionalRoleArgs := map[string]any{"replacement": b505ConditionalRoleKey == "previous", "previous": 17, "requested": 23}
		b505ConditionalRoleArgs[b505ConditionalRoleKey] = 31
		b505ConditionalRoleResult := b505ConditionalRoleRun(b505ConditionalRoleT, b505ConditionalRoleInput("unchanged", 61, map[string]any{"replacement": []int{0}}, b505ConditionalRoleArgs))["results"].(map[string]any)["unchanged"].(map[string]any)
		if b505ConditionalRoleResult["number"] != float64(31) || b505ConditionalRoleResult["replacement_code"] != "Unknown" {
			b505ConditionalRoleT.Fatal(b505ConditionalRoleResult)
		}
	}
	for _, b505ConditionalRoleMode := range []int{0, 1, 255} {
		for _, b505ConditionalRoleStart := range [][]int{{}, {0}, {255, 1, 2, 3}} {
			b505ConditionalRoleResult := b505ConditionalRoleRun(b505ConditionalRoleT, b505ConditionalRoleInput("programmed-start", 9, map[string]any{"start": b505ConditionalRoleStart}, map[string]any{"program_mode": b505ConditionalRoleMode}))["results"].(map[string]any)["start"].(map[string]any)
			b505ConditionalRoleLabel := "ignored"
			if b505ConditionalRoleMode == 0 {
				b505ConditionalRoleLabel = "applicable"
			}
			if b505ConditionalRoleResult["label"] != b505ConditionalRoleLabel {
				b505ConditionalRoleT.Fatal(b505ConditionalRoleResult)
			}
			b505ConditionalRoleBytes, _ := json.Marshal(b505ConditionalRoleStart)
			var b505ConditionalRoleExpected any
			_ = json.Unmarshal(b505ConditionalRoleBytes, &b505ConditionalRoleExpected)
			if !reflect.DeepEqual(b505ConditionalRoleResult["raw"], b505ConditionalRoleExpected) {
				b505ConditionalRoleT.Fatal("start changed")
			}
		}
	}
	for _, b505ConditionalRoleField := range []string{"enabled", "disable", "start"} {
		for _, b505ConditionalRoleBad := range []any{nil, 1, "1", false, []any{nil}, []any{true}, []any{"1"}, []any{256}, []any{1.5}} {
			b505ConditionalRoleBytes, _ := json.Marshal(b505ConditionalRoleInput("unknown", 0, map[string]any{b505ConditionalRoleField: b505ConditionalRoleBad}, map[string]any{}))
			if b505ConditionalRoleOutput, b505ConditionalRoleError := B505ConditionalRoleEvaluate(b505ConditionalRoleBytes); b505ConditionalRoleError == nil || b505ConditionalRoleOutput != nil {
				b505ConditionalRoleT.Fatal("invalid raw field accepted", b505ConditionalRoleField, b505ConditionalRoleBad)
			}
		}
	}
	for _, b505ConditionalRoleBad := range []string{`{"context":{"operation":1.0000000000000000001},"fields":{}}`, `{"context":{},"fields":{"enabled":[0.99999999999999999999]}}`, `{"context":{},"fields":{}} false`} {
		if b505ConditionalRoleOutput, b505ConditionalRoleError := B505ConditionalRoleEvaluate([]byte(b505ConditionalRoleBad)); b505ConditionalRoleError == nil || b505ConditionalRoleOutput != nil {
			b505ConditionalRoleT.Fatal("precision/trailing accepted")
		}
	}
}

func TestAb05Aq05BoundedNumbers(b505ConditionalRoleT *testing.T) {
	b505ConditionalRolePrefix := `{"context":{"primary":181,"secondary":5,"operation":23,"direction":"request","view":"flow-excess"},"fields":{},"arguments":{`
	for _, b505ConditionalRoleBad := range []string{`"demand":1e-100000,"excess":5`, `"demand":1,"excess":1e1000000`, `"demand":1e-401,"excess":5`, `"demand":0.` + strings.Repeat("0", 70) + `1,"excess":5`} {
		if b505ConditionalRoleOutput, b505ConditionalRoleError := B505ConditionalRoleEvaluate([]byte(b505ConditionalRolePrefix + b505ConditionalRoleBad + `}}`)); b505ConditionalRoleError == nil || b505ConditionalRoleOutput != nil {
			b505ConditionalRoleT.Fatalf("accepted %s", b505ConditionalRoleBad)
		}
	}
	b505ConditionalRoleOutput, b505ConditionalRoleError := B505ConditionalRoleEvaluate([]byte(b505ConditionalRolePrefix + `"demand":40,"excess":5e0}}`))
	if b505ConditionalRoleError != nil || !strings.Contains(string(b505ConditionalRoleOutput), `"number":45`) {
		b505ConditionalRoleT.Fatalf("bounded sum failed: %s %v", b505ConditionalRoleOutput, b505ConditionalRoleError)
	}
}

func TestAb05Aq05ReservedKeysAndCollisions(b505ConditionalRoleT *testing.T) {
	for _, b505ConditionalRoleKey := range []string{"state", "results", "physical", "qualified"} {
		b505ConditionalRoleObject := b505ConditionalRoleInput("dhw-enable", 49, map[string]any{"enabled": []int{1}}, map[string]any{})
		b505ConditionalRoleObject[b505ConditionalRoleKey] = "Proven"
		b505ConditionalRoleBytes, _ := json.Marshal(b505ConditionalRoleObject)
		if b505ConditionalRoleOutput, b505ConditionalRoleError := B505ConditionalRoleEvaluate(b505ConditionalRoleBytes); b505ConditionalRoleError == nil || b505ConditionalRoleOutput != nil {
			b505ConditionalRoleT.Fatalf("accepted reserved key %s", b505ConditionalRoleKey)
		}
	}
	b505ConditionalRoleObject := b505ConditionalRoleInput("flow-excess", 23, map[string]any{"flow-excess": []int{7, 8}}, map[string]any{"demand": json.Number("40"), "excess": json.Number("5")})
	b505ConditionalRoleResults := b505ConditionalRoleRun(b505ConditionalRoleT, b505ConditionalRoleObject)["results"].(map[string]any)
	b505ConditionalRoleView := b505ConditionalRoleResults["flow-excess"].(map[string]any)
	if !reflect.DeepEqual(b505ConditionalRoleView["field_raw"], []any{float64(7), float64(8)}) || b505ConditionalRoleView["number"] != float64(45) {
		b505ConditionalRoleT.Fatalf("collision lost field raw: %v", b505ConditionalRoleView)
	}
}

func TestAb05Aq05CanonicalSum(b505ConditionalRoleT *testing.T) {
	b505ConditionalRoleObject := b505ConditionalRoleInput("flow-excess", 23, map[string]any{}, map[string]any{"demand": json.Number("0.25"), "excess": json.Number("0.1")})
	b505ConditionalRoleBytes, _ := json.Marshal(b505ConditionalRoleObject)
	b505ConditionalRoleOutput, b505ConditionalRoleError := B505ConditionalRoleEvaluate(b505ConditionalRoleBytes)
	if b505ConditionalRoleError != nil {
		b505ConditionalRoleT.Fatal(b505ConditionalRoleError)
	}
	if !strings.Contains(string(b505ConditionalRoleOutput), `"number":0.35,`) && !strings.Contains(string(b505ConditionalRoleOutput), `"number":0.35}`) {
		b505ConditionalRoleT.Fatalf("non-canonical sum: %s", b505ConditionalRoleOutput)
	}
}

func TestAb05Aq05RejectDeepNesting(b505ConditionalRoleT *testing.T) {
	b505ConditionalRoleDeep := `{"context":{},"fields":{},"x":` + strings.Repeat("[", 100000) + strings.Repeat("]", 100000) + `}`
	if b505ConditionalRoleOutput, b505ConditionalRoleError := B505ConditionalRoleEvaluate([]byte(b505ConditionalRoleDeep)); b505ConditionalRoleError == nil || b505ConditionalRoleOutput != nil {
		b505ConditionalRoleT.Fatalf("accepted over-deep input")
	}
	b505ConditionalRoleShallow := `{"context":{},"fields":{},"x":` + strings.Repeat("[", 8) + strings.Repeat("]", 8) + `}`
	if _, b505ConditionalRoleError := B505ConditionalRoleEvaluate([]byte(b505ConditionalRoleShallow)); b505ConditionalRoleError != nil {
		b505ConditionalRoleT.Fatalf("rejected shallow input: %v", b505ConditionalRoleError)
	}
}
