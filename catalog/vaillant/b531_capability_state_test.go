package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func b531CapabilityStateInput() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 49, "operation": 1, "direction": "local", "view": "capability", "family": "synthetic-family", "variant": map[string]any{"opaque": []any{nil, 0, false}}}, "fields": map[string]any{"alpha": []any{0, 255}, "beta": []any{2}}, "arguments": map[string]any{"registered": 2, "prior": 2, "result": 1, "extra": map[string]any{"opaque": []any{false, nil}}}}
}

func b531CapabilityStateEvaluate(b531CapabilityStateT *testing.T, b531CapabilityStateValue any) map[string]any {
	b531CapabilityStateT.Helper()
	b531CapabilityStateBytes, b531CapabilityStateErr := json.Marshal(b531CapabilityStateValue)
	if b531CapabilityStateErr != nil {
		b531CapabilityStateT.Fatal(b531CapabilityStateErr)
	}
	b531CapabilityStateResult, b531CapabilityStateErr := B531CapabilityStateEvaluate(b531CapabilityStateBytes)
	if b531CapabilityStateErr != nil {
		b531CapabilityStateT.Fatal(b531CapabilityStateErr)
	}
	b531CapabilityStateOutput := map[string]any{}
	if b531CapabilityStateErr = json.Unmarshal(b531CapabilityStateResult, &b531CapabilityStateOutput); b531CapabilityStateErr != nil {
		b531CapabilityStateT.Fatal(b531CapabilityStateErr)
	}
	if b531CapabilityStateOutput["physical"] != "Unknown" || b531CapabilityStateOutput["qualified"] != false {
		b531CapabilityStateT.Fatal(b531CapabilityStateOutput)
	}
	return b531CapabilityStateOutput
}

func TestAb01Ap01Transitions(b531CapabilityStateT *testing.T) {
	b531CapabilityStateBytes, b531CapabilityStateErr := os.ReadFile("testdata/ab01ap01_b531_capability_state.json")
	if b531CapabilityStateErr != nil {
		b531CapabilityStateT.Fatal(b531CapabilityStateErr)
	}
	var b531CapabilityStateCases []map[string]any
	if b531CapabilityStateErr = json.Unmarshal(b531CapabilityStateBytes, &b531CapabilityStateCases); b531CapabilityStateErr != nil {
		b531CapabilityStateT.Fatal(b531CapabilityStateErr)
	}
	if len(b531CapabilityStateCases) != 6 {
		b531CapabilityStateT.Fatal("missing cases")
	}
	for _, b531CapabilityStateCase := range b531CapabilityStateCases {
		b531CapabilityStateT.Run(b531CapabilityStateCase["id"].(string), func(b531CapabilityStateT *testing.T) {
			b531CapabilityStateValue := b531CapabilityStateInput()
			b531CapabilityStateArgs := b531CapabilityStateValue["arguments"].(map[string]any)
			for _, b531CapabilityStateKey := range []string{"registered", "prior", "result"} {
				b531CapabilityStateArgs[b531CapabilityStateKey] = b531CapabilityStateCase[b531CapabilityStateKey]
			}
			b531CapabilityStateOutput := b531CapabilityStateEvaluate(b531CapabilityStateT, b531CapabilityStateValue)
			if b531CapabilityStateOutput["state"] != "complete" {
				b531CapabilityStateT.Fatal(b531CapabilityStateOutput)
			}
			b531CapabilityStateState := b531CapabilityStateOutput["capability"].(map[string]any)
			for _, b531CapabilityStateKey := range []string{"updated", "cached_bit0", "whole_nonzero", "outcome"} {
				if b531CapabilityStateState[b531CapabilityStateKey] != b531CapabilityStateCase[b531CapabilityStateKey] {
					b531CapabilityStateT.Fatalf("%s: %v", b531CapabilityStateKey, b531CapabilityStateState)
				}
			}
			if b531CapabilityStateState["marker"] != float64(65535) || b531CapabilityStateState["state"] != "Hypothesis" || b531CapabilityStateState["calculation"] != "Proven" {
				b531CapabilityStateT.Fatal(b531CapabilityStateState)
			}
		})
	}
}

func b531CapabilityStateRows() []any {
	b531CapabilityStateRows := []any{}
	for b531CapabilityStateIndex := 0; b531CapabilityStateIndex < 9; b531CapabilityStateIndex++ {
		b531CapabilityStateRows = append(b531CapabilityStateRows, map[string]any{"a": 0, "b": 0, "c": 0})
	}
	return b531CapabilityStateRows
}

func TestAb01Ap01Selection(b531CapabilityStateT *testing.T) {
	for b531CapabilityStateIndex := 0; b531CapabilityStateIndex < 9; b531CapabilityStateIndex++ {
		for _, b531CapabilityStateZero := range []string{"", "a", "b", "c"} {
			b531CapabilityStateValue := b531CapabilityStateInput()
			b531CapabilityStateValue["context"].(map[string]any)["view"] = "select"
			b531CapabilityStateRows := b531CapabilityStateRows()
			b531CapabilityStateRow := map[string]any{"a": 2, "b": 2, "c": 2}
			b531CapabilityStateRows[b531CapabilityStateIndex] = b531CapabilityStateRow
			b531CapabilityStateExpected := b531CapabilityStateIndex
			if b531CapabilityStateZero != "" {
				b531CapabilityStateRow[b531CapabilityStateZero] = 0
				b531CapabilityStateExpected = 9
			}
			b531CapabilityStateValue["arguments"] = map[string]any{"rows": b531CapabilityStateRows}
			b531CapabilityStateOutput := b531CapabilityStateEvaluate(b531CapabilityStateT, b531CapabilityStateValue)
			b531CapabilityStateSelection := b531CapabilityStateOutput["select"].(map[string]any)
			if b531CapabilityStateOutput["state"] != "complete" || b531CapabilityStateSelection["state"] != "Hypothesis" || b531CapabilityStateSelection["calculation"] != "Proven" {
				b531CapabilityStateT.Fatal("selection grades", b531CapabilityStateOutput)
			}
			if b531CapabilityStateOutput["select"].(map[string]any)["index"] != float64(b531CapabilityStateExpected) {
				b531CapabilityStateT.Fatalf("index %d zero %s: %v", b531CapabilityStateIndex, b531CapabilityStateZero, b531CapabilityStateOutput)
			}
		}
	}
	for _, b531CapabilityStatePair := range [][2]int{{0, 8}, {1, 8}, {7, 8}} {
		b531CapabilityStateValue := b531CapabilityStateInput()
		b531CapabilityStateValue["context"].(map[string]any)["view"] = "select"
		b531CapabilityStateRows := b531CapabilityStateRows()
		for _, b531CapabilityStateIndex := range b531CapabilityStatePair {
			b531CapabilityStateRows[b531CapabilityStateIndex] = map[string]any{"a": 2, "b": 2, "c": 2}
		}
		b531CapabilityStateValue["arguments"] = map[string]any{"rows": b531CapabilityStateRows}
		if b531CapabilityStateEvaluate(b531CapabilityStateT, b531CapabilityStateValue)["select"].(map[string]any)["index"] != float64(b531CapabilityStatePair[0]) {
			b531CapabilityStateT.Fatal("first match ordering")
		}
	}
}

func TestAb01Ap01IdentityAndRaw(b531CapabilityStateT *testing.T) {
	for _, b531CapabilityStateKey := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
		for _, b531CapabilityStateItem := range []any{nil, 0, "opaque", false} {
			b531CapabilityStateValue := b531CapabilityStateInput()
			b531CapabilityStateContext := b531CapabilityStateValue["context"].(map[string]any)
			b531CapabilityStateContext[b531CapabilityStateKey] = b531CapabilityStateItem
			b531CapabilityStateBytes, _ := json.Marshal(b531CapabilityStateValue)
			b531CapabilityStateResult, b531CapabilityStateErr := B531CapabilityStateEvaluate(b531CapabilityStateBytes)
			b531CapabilityStateOpaque := b531CapabilityStateKey == "family" || b531CapabilityStateKey == "variant"
			b531CapabilityStateWrong := !b531CapabilityStateOpaque && b531CapabilityStateItem != nil && ((b531CapabilityStateKey == "direction" || b531CapabilityStateKey == "view") && b531CapabilityStateItem != "opaque" || (b531CapabilityStateKey != "direction" && b531CapabilityStateKey != "view") && b531CapabilityStateItem != 0)
			if b531CapabilityStateWrong {
				if b531CapabilityStateErr == nil || b531CapabilityStateResult != nil {
					b531CapabilityStateT.Fatalf("accepted wrong selector kind %s", b531CapabilityStateKey)
				}
				continue
			}
			if b531CapabilityStateErr != nil {
				b531CapabilityStateT.Fatal(b531CapabilityStateErr)
			}
			b531CapabilityStateOutput := b531CapabilityStateEvaluate(b531CapabilityStateT, b531CapabilityStateValue)
			b531CapabilityStateExpected := "Unknown"
			if b531CapabilityStateOpaque {
				b531CapabilityStateExpected = "complete"
			}
			if b531CapabilityStateOutput["state"] != b531CapabilityStateExpected {
				b531CapabilityStateT.Fatal(b531CapabilityStateOutput)
			}
			b531CapabilityStateRound, _ := json.Marshal(b531CapabilityStateOutput["context"])
			b531CapabilityStateOriginal, _ := json.Marshal(b531CapabilityStateContext)
			if string(b531CapabilityStateRound) != string(b531CapabilityStateOriginal) {
				b531CapabilityStateT.Fatal("context changed")
			}
			if !b531CapabilityStateOpaque {
				if _, b531CapabilityStateExists := b531CapabilityStateOutput["capability"]; b531CapabilityStateExists {
					b531CapabilityStateT.Fatal("unsupported formula")
				}
			}
		}
		b531CapabilityStateValue := b531CapabilityStateInput()
		delete(b531CapabilityStateValue["context"].(map[string]any), b531CapabilityStateKey)
		b531CapabilityStateOutput := b531CapabilityStateEvaluate(b531CapabilityStateT, b531CapabilityStateValue)
		if _, b531CapabilityStateExists := b531CapabilityStateOutput["context"].(map[string]any)[b531CapabilityStateKey]; b531CapabilityStateExists {
			b531CapabilityStateT.Fatal("missing became present")
		}
	}
	for _, b531CapabilityStateKey := range []string{"alpha", "beta"} {
		for _, b531CapabilityStateBytes := range [][]any{nil, {}, {0}, {0, 1, 2, 3, 4}} {
			b531CapabilityStateValue := b531CapabilityStateInput()
			b531CapabilityStateRaw := b531CapabilityStateValue["fields"].(map[string]any)
			if b531CapabilityStateBytes == nil {
				delete(b531CapabilityStateRaw, b531CapabilityStateKey)
			} else {
				b531CapabilityStateRaw[b531CapabilityStateKey] = b531CapabilityStateBytes
			}
			b531CapabilityStateOutput := b531CapabilityStateEvaluate(b531CapabilityStateT, b531CapabilityStateValue)
			b531CapabilityStateWant, _ := json.Marshal(b531CapabilityStateRaw)
			b531CapabilityStateGot, _ := json.Marshal(b531CapabilityStateOutput["raw"])
			if string(b531CapabilityStateWant) != string(b531CapabilityStateGot) {
				b531CapabilityStateT.Fatal("raw changed")
			}
			b531CapabilityStateExpected := "complete"
			if b531CapabilityStateBytes != nil && len(b531CapabilityStateBytes) == 0 {
				b531CapabilityStateExpected = "partial"
			}
			if b531CapabilityStateOutput["state"] != b531CapabilityStateExpected || b531CapabilityStateOutput["capability"].(map[string]any)["updated"] != float64(2) {
				b531CapabilityStateT.Fatal("independent result lost")
			}
		}
	}
}

func TestAb01Ap01Rejections(b531CapabilityStateT *testing.T) {
	for _, b531CapabilityStateText := range []string{``, `{`, `[]`, `{} {}`, `{"x":1,"x":2}`, `{"x":[{"a":0,"a":1}]}`, `{"context":{},"fields":{},"arguments":{"extra":{"x":1,"x":2}}}`} {
		b531CapabilityStateResult, b531CapabilityStateErr := B531CapabilityStateEvaluate(json.RawMessage(b531CapabilityStateText))
		if b531CapabilityStateErr == nil || b531CapabilityStateResult != nil {
			b531CapabilityStateT.Fatalf("accepted %s", b531CapabilityStateText)
		}
	}
	for _, b531CapabilityStateKey := range []string{"registered", "prior", "result"} {
		for _, b531CapabilityStateBad := range []any{nil, false, "1", 1.5, -1, 4294967296.0, []any{}, map[string]any{}} {
			b531CapabilityStateValue := b531CapabilityStateInput()
			b531CapabilityStateValue["arguments"].(map[string]any)[b531CapabilityStateKey] = b531CapabilityStateBad
			b531CapabilityStateReject(b531CapabilityStateT, b531CapabilityStateValue)
		}
		b531CapabilityStateValue := b531CapabilityStateInput()
		delete(b531CapabilityStateValue["arguments"].(map[string]any), b531CapabilityStateKey)
		b531CapabilityStateReject(b531CapabilityStateT, b531CapabilityStateValue)
	}
	for _, b531CapabilityStateKey := range []string{"alpha", "beta"} {
		for _, b531CapabilityStateBad := range []any{nil, false, 0, "", []any{false}, []any{-1}, []any{256}, []any{1.5}, []any{nil}} {
			b531CapabilityStateValue := b531CapabilityStateInput()
			b531CapabilityStateValue["fields"].(map[string]any)[b531CapabilityStateKey] = b531CapabilityStateBad
			b531CapabilityStateReject(b531CapabilityStateT, b531CapabilityStateValue)
		}
	}
	for _, b531CapabilityStateSize := range []int{8, 10} {
		b531CapabilityStateValue := b531CapabilityStateInput()
		b531CapabilityStateValue["context"].(map[string]any)["view"] = "select"
		b531CapabilityStateRows := b531CapabilityStateRows()
		if b531CapabilityStateSize == 8 {
			b531CapabilityStateRows = b531CapabilityStateRows[:8]
		} else {
			b531CapabilityStateRows = append(b531CapabilityStateRows, map[string]any{"a": 0, "b": 0, "c": 0})
		}
		b531CapabilityStateValue["arguments"] = map[string]any{"rows": b531CapabilityStateRows}
		b531CapabilityStateReject(b531CapabilityStateT, b531CapabilityStateValue)
	}
	for b531CapabilityStateIndex := 0; b531CapabilityStateIndex < 9; b531CapabilityStateIndex++ {
		for _, b531CapabilityStateKey := range []string{"a", "b", "c"} {
			for _, b531CapabilityStateBad := range []any{nil, false, "2", -1, 1.5, 4294967296.0} {
				b531CapabilityStateValue := b531CapabilityStateInput()
				b531CapabilityStateValue["context"].(map[string]any)["view"] = "select"
				b531CapabilityStateRows := b531CapabilityStateRows()
				b531CapabilityStateRows[b531CapabilityStateIndex].(map[string]any)[b531CapabilityStateKey] = b531CapabilityStateBad
				b531CapabilityStateValue["arguments"] = map[string]any{"rows": b531CapabilityStateRows}
				b531CapabilityStateReject(b531CapabilityStateT, b531CapabilityStateValue)
			}
		}
	}
	for _, b531CapabilityStateKey := range []string{"context", "fields", "arguments"} {
		for _, b531CapabilityStateBad := range []any{nil, false, []any{}, 0, ""} {
			b531CapabilityStateValue := b531CapabilityStateInput()
			b531CapabilityStateValue[b531CapabilityStateKey] = b531CapabilityStateBad
			b531CapabilityStateReject(b531CapabilityStateT, b531CapabilityStateValue)
		}
	}
	b531CapabilityStateValue := b531CapabilityStateInput()
	delete(b531CapabilityStateValue, "arguments")
	b531CapabilityStateReject(b531CapabilityStateT, b531CapabilityStateValue)
}

func b531CapabilityStateReject(b531CapabilityStateT *testing.T, b531CapabilityStateValue any) {
	b531CapabilityStateT.Helper()
	b531CapabilityStateBytes, _ := json.Marshal(b531CapabilityStateValue)
	b531CapabilityStateResult, b531CapabilityStateErr := B531CapabilityStateEvaluate(b531CapabilityStateBytes)
	if b531CapabilityStateErr == nil || b531CapabilityStateResult != nil {
		b531CapabilityStateT.Fatalf("accepted %s", b531CapabilityStateBytes)
	}
}

func TestAb01Ap01Isolation(b531CapabilityStateT *testing.T) {
	b531CapabilityStateBytes, _ := json.Marshal(b531CapabilityStateInput())
	b531CapabilityStateOriginal := append([]byte(nil), b531CapabilityStateBytes...)
	b531CapabilityStateResult, b531CapabilityStateErr := B531CapabilityStateEvaluate(b531CapabilityStateBytes)
	if b531CapabilityStateErr != nil {
		b531CapabilityStateT.Fatal(b531CapabilityStateErr)
	}
	b531CapabilityStateExpected := append([]byte(nil), b531CapabilityStateResult...)
	for b531CapabilityStateIndex := range b531CapabilityStateBytes {
		b531CapabilityStateBytes[b531CapabilityStateIndex] = ' '
	}
	if !reflect.DeepEqual([]byte(b531CapabilityStateResult), b531CapabilityStateExpected) {
		b531CapabilityStateT.Fatal("input aliases output")
	}
	b531CapabilityStateDecoded := map[string]any{}
	if b531CapabilityStateErr := json.Unmarshal(b531CapabilityStateResult, &b531CapabilityStateDecoded); b531CapabilityStateErr != nil {
		b531CapabilityStateT.Fatal(b531CapabilityStateErr)
	}
	b531CapabilityStateDecoded["raw"].(map[string]any)["alpha"].([]any)[0] = 99
	b531CapabilityStateDecoded["context"].(map[string]any)["variant"].(map[string]any)["opaque"].([]any)[0] = true
	b531CapabilityStateDecoded["arguments"].(map[string]any)["extra"].(map[string]any)["opaque"].([]any)[0] = true
	if b531CapabilityStateDecoded["fields"].(map[string]any)["alpha"].(map[string]any)["raw"].([]any)[0] != float64(0) {
		b531CapabilityStateT.Fatal("nested output aliases")
	}
	for b531CapabilityStateIndex := range b531CapabilityStateResult {
		b531CapabilityStateResult[b531CapabilityStateIndex] = ' '
	}
	b531CapabilityStateAgain, b531CapabilityStateErr := B531CapabilityStateEvaluate(b531CapabilityStateOriginal)
	if b531CapabilityStateErr != nil || string(b531CapabilityStateAgain) != string(b531CapabilityStateExpected) {
		b531CapabilityStateT.Fatal("later call changed")
	}
}

func TestAb01Ap01AdditionalBoundaries(b531CapabilityStateT *testing.T) {
	for _, b531CapabilityStateKey := range []string{"context", "fields"} {
		b531CapabilityStateValue := b531CapabilityStateInput()
		delete(b531CapabilityStateValue, b531CapabilityStateKey)
		b531CapabilityStateReject(b531CapabilityStateT, b531CapabilityStateValue)
	}
	for _, b531CapabilityStateKey := range []string{"primary", "secondary", "operation"} {
		for _, b531CapabilityStateBad := range []any{-1, 256, 1.5, []any{}, map[string]any{}} {
			b531CapabilityStateValue := b531CapabilityStateInput()
			b531CapabilityStateValue["context"].(map[string]any)[b531CapabilityStateKey] = b531CapabilityStateBad
			b531CapabilityStateReject(b531CapabilityStateT, b531CapabilityStateValue)
		}
	}
	for _, b531CapabilityStateKey := range []string{"rows", "a", "b", "c"} {
		b531CapabilityStateValue := b531CapabilityStateInput()
		b531CapabilityStateValue["context"].(map[string]any)["view"] = "select"
		b531CapabilityStateRows := b531CapabilityStateRows()
		b531CapabilityStateValue["arguments"] = map[string]any{"rows": b531CapabilityStateRows}
		if b531CapabilityStateKey == "rows" {
			delete(b531CapabilityStateValue["arguments"].(map[string]any), "rows")
		} else {
			delete(b531CapabilityStateRows[8].(map[string]any), b531CapabilityStateKey)
		}
		b531CapabilityStateReject(b531CapabilityStateT, b531CapabilityStateValue)
	}
	for _, b531CapabilityStateBad := range []any{nil, false, 0, "", []any{nil, nil, nil, nil, nil, nil, nil, nil, nil}} {
		b531CapabilityStateValue := b531CapabilityStateInput()
		b531CapabilityStateValue["context"].(map[string]any)["view"] = "select"
		b531CapabilityStateValue["arguments"] = map[string]any{"rows": b531CapabilityStateBad}
		b531CapabilityStateReject(b531CapabilityStateT, b531CapabilityStateValue)
	}
	for b531CapabilityStateIndex := 0; b531CapabilityStateIndex < 9; b531CapabilityStateIndex++ {
		for _, b531CapabilityStateKey := range []string{"a", "b", "c"} {
			b531CapabilityStateValue := b531CapabilityStateInput()
			b531CapabilityStateValue["context"].(map[string]any)["view"] = "select"
			b531CapabilityStateRows := b531CapabilityStateRows()
			b531CapabilityStateRows[b531CapabilityStateIndex] = map[string]any{"a": 2, "b": 2, "c": 2}
			b531CapabilityStateRows[b531CapabilityStateIndex].(map[string]any)[b531CapabilityStateKey] = uint64(4294967295)
			b531CapabilityStateValue["arguments"] = map[string]any{"rows": b531CapabilityStateRows}
			if b531CapabilityStateEvaluate(b531CapabilityStateT, b531CapabilityStateValue)["select"].(map[string]any)["index"] != float64(b531CapabilityStateIndex) {
				b531CapabilityStateT.Fatal("u32 boundary")
			}
		}
	}
	b531CapabilityStateValue := b531CapabilityStateInput()
	b531CapabilityStateValue["context"].(map[string]any)["view"] = "unsupported"
	delete(b531CapabilityStateValue, "arguments")
	b531CapabilityStateOutput := b531CapabilityStateEvaluate(b531CapabilityStateT, b531CapabilityStateValue)
	if _, b531CapabilityStatePresent := b531CapabilityStateOutput["arguments"]; b531CapabilityStatePresent || b531CapabilityStateOutput["state"] != "Unknown" {
		b531CapabilityStateT.Fatal("absent arguments changed")
	}
	for _, b531CapabilityStateKey := range []string{"alpha", "beta"} {
		b531CapabilityStateField := b531CapabilityStateOutput["fields"].(map[string]any)[b531CapabilityStateKey].(map[string]any)
		if b531CapabilityStateField["state"] != "Unknown" || b531CapabilityStateField["unmodeled"] != true {
			b531CapabilityStateT.Fatal("raw-only boundary")
		}
		for _, b531CapabilityStateInvented := range []string{"number", "label"} {
			if _, b531CapabilityStatePresent := b531CapabilityStateField[b531CapabilityStateInvented]; b531CapabilityStatePresent {
				b531CapabilityStateT.Fatal("invented semantics")
			}
		}
	}
	for _, b531CapabilityStateText := range []string{
		`{"context":{"primary":181,"primary":0},"fields":{}}`,
		`{"context":{},"fields":{"alpha":[0],"alpha":[1]}}`,
		`{"context":{},"fields":{},"arguments":{"rows":[{"a":0,"a":1}]}}`,
		`{"context":{},"fields":{}} trailing`,
		`{"context":{},"fields":{"alpha":[1e0]}}`,
	} {
		b531CapabilityStateResult, b531CapabilityStateErr := B531CapabilityStateEvaluate(json.RawMessage(b531CapabilityStateText))
		if b531CapabilityStateErr == nil || b531CapabilityStateResult != nil {
			b531CapabilityStateT.Fatal("invalid JSON accepted")
		}
	}
}
