package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func ab01ap01Input() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 49, "operation": 1, "direction": "local", "view": "capability", "family": "synthetic-family", "variant": map[string]any{"opaque": []any{nil, 0, false}}}, "fields": map[string]any{"alpha": []any{0, 255}, "beta": []any{2}}, "arguments": map[string]any{"registered": 2, "prior": 2, "result": 1, "extra": map[string]any{"opaque": []any{false, nil}}}}
}

func ab01ap01Evaluate(ab01ap01T *testing.T, ab01ap01Value any) map[string]any {
	ab01ap01T.Helper()
	ab01ap01Bytes, ab01ap01Err := json.Marshal(ab01ap01Value)
	if ab01ap01Err != nil {
		ab01ap01T.Fatal(ab01ap01Err)
	}
	ab01ap01Result, ab01ap01Err := Ab01Ap01Evaluate(ab01ap01Bytes)
	if ab01ap01Err != nil {
		ab01ap01T.Fatal(ab01ap01Err)
	}
	ab01ap01Output := map[string]any{}
	if ab01ap01Err = json.Unmarshal(ab01ap01Result, &ab01ap01Output); ab01ap01Err != nil {
		ab01ap01T.Fatal(ab01ap01Err)
	}
	if ab01ap01Output["physical"] != "Unknown" || ab01ap01Output["qualified"] != false {
		ab01ap01T.Fatal(ab01ap01Output)
	}
	return ab01ap01Output
}

func TestAb01Ap01Transitions(ab01ap01T *testing.T) {
	ab01ap01Bytes, ab01ap01Err := os.ReadFile("testdata/ab01ap01_b531_capability_state.json")
	if ab01ap01Err != nil {
		ab01ap01T.Fatal(ab01ap01Err)
	}
	var ab01ap01Cases []map[string]any
	if ab01ap01Err = json.Unmarshal(ab01ap01Bytes, &ab01ap01Cases); ab01ap01Err != nil {
		ab01ap01T.Fatal(ab01ap01Err)
	}
	if len(ab01ap01Cases) != 6 {
		ab01ap01T.Fatal("missing cases")
	}
	for _, ab01ap01Case := range ab01ap01Cases {
		ab01ap01T.Run(ab01ap01Case["id"].(string), func(ab01ap01T *testing.T) {
			ab01ap01Value := ab01ap01Input()
			ab01ap01Args := ab01ap01Value["arguments"].(map[string]any)
			for _, ab01ap01Key := range []string{"registered", "prior", "result"} {
				ab01ap01Args[ab01ap01Key] = ab01ap01Case[ab01ap01Key]
			}
			ab01ap01Output := ab01ap01Evaluate(ab01ap01T, ab01ap01Value)
			if ab01ap01Output["state"] != "complete" {
				ab01ap01T.Fatal(ab01ap01Output)
			}
			ab01ap01State := ab01ap01Output["capability"].(map[string]any)
			for _, ab01ap01Key := range []string{"updated", "cached_bit0", "whole_nonzero", "outcome"} {
				if ab01ap01State[ab01ap01Key] != ab01ap01Case[ab01ap01Key] {
					ab01ap01T.Fatalf("%s: %v", ab01ap01Key, ab01ap01State)
				}
			}
			if ab01ap01State["marker"] != float64(65535) || ab01ap01State["state"] != "Hypothesis" || ab01ap01State["calculation"] != "Proven" {
				ab01ap01T.Fatal(ab01ap01State)
			}
		})
	}
}

func ab01ap01Rows() []any {
	ab01ap01Rows := []any{}
	for ab01ap01Index := 0; ab01ap01Index < 9; ab01ap01Index++ {
		ab01ap01Rows = append(ab01ap01Rows, map[string]any{"a": 0, "b": 0, "c": 0})
	}
	return ab01ap01Rows
}

func TestAb01Ap01Selection(ab01ap01T *testing.T) {
	for ab01ap01Index := 0; ab01ap01Index < 9; ab01ap01Index++ {
		for _, ab01ap01Zero := range []string{"", "a", "b", "c"} {
			ab01ap01Value := ab01ap01Input()
			ab01ap01Value["context"].(map[string]any)["view"] = "select"
			ab01ap01Rows := ab01ap01Rows()
			ab01ap01Row := map[string]any{"a": 2, "b": 2, "c": 2}
			ab01ap01Rows[ab01ap01Index] = ab01ap01Row
			ab01ap01Expected := ab01ap01Index
			if ab01ap01Zero != "" {
				ab01ap01Row[ab01ap01Zero] = 0
				ab01ap01Expected = 9
			}
			ab01ap01Value["arguments"] = map[string]any{"rows": ab01ap01Rows}
			ab01ap01Output := ab01ap01Evaluate(ab01ap01T, ab01ap01Value)
			ab01ap01Selection := ab01ap01Output["select"].(map[string]any)
			if ab01ap01Output["state"] != "complete" || ab01ap01Selection["state"] != "Hypothesis" || ab01ap01Selection["calculation"] != "Proven" {
				ab01ap01T.Fatal("selection grades", ab01ap01Output)
			}
			if ab01ap01Output["select"].(map[string]any)["index"] != float64(ab01ap01Expected) {
				ab01ap01T.Fatalf("index %d zero %s: %v", ab01ap01Index, ab01ap01Zero, ab01ap01Output)
			}
		}
	}
	for _, ab01ap01Pair := range [][2]int{{0, 8}, {1, 8}, {7, 8}} {
		ab01ap01Value := ab01ap01Input()
		ab01ap01Value["context"].(map[string]any)["view"] = "select"
		ab01ap01Rows := ab01ap01Rows()
		for _, ab01ap01Index := range ab01ap01Pair {
			ab01ap01Rows[ab01ap01Index] = map[string]any{"a": 2, "b": 2, "c": 2}
		}
		ab01ap01Value["arguments"] = map[string]any{"rows": ab01ap01Rows}
		if ab01ap01Evaluate(ab01ap01T, ab01ap01Value)["select"].(map[string]any)["index"] != float64(ab01ap01Pair[0]) {
			ab01ap01T.Fatal("first match ordering")
		}
	}
}

func TestAb01Ap01IdentityAndRaw(ab01ap01T *testing.T) {
	for _, ab01ap01Key := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
		for _, ab01ap01Item := range []any{nil, 0, "opaque", false} {
			ab01ap01Value := ab01ap01Input()
			ab01ap01Context := ab01ap01Value["context"].(map[string]any)
			ab01ap01Context[ab01ap01Key] = ab01ap01Item
			ab01ap01Bytes, _ := json.Marshal(ab01ap01Value)
			ab01ap01Result, ab01ap01Err := Ab01Ap01Evaluate(ab01ap01Bytes)
			ab01ap01Opaque := ab01ap01Key == "family" || ab01ap01Key == "variant"
			ab01ap01Wrong := !ab01ap01Opaque && ab01ap01Item != nil && ((ab01ap01Key == "direction" || ab01ap01Key == "view") && ab01ap01Item != "opaque" || (ab01ap01Key != "direction" && ab01ap01Key != "view") && ab01ap01Item != 0)
			if ab01ap01Wrong {
				if ab01ap01Err == nil || ab01ap01Result != nil {
					ab01ap01T.Fatalf("accepted wrong selector kind %s", ab01ap01Key)
				}
				continue
			}
			if ab01ap01Err != nil {
				ab01ap01T.Fatal(ab01ap01Err)
			}
			ab01ap01Output := ab01ap01Evaluate(ab01ap01T, ab01ap01Value)
			ab01ap01Expected := "Unknown"
			if ab01ap01Opaque {
				ab01ap01Expected = "complete"
			}
			if ab01ap01Output["state"] != ab01ap01Expected {
				ab01ap01T.Fatal(ab01ap01Output)
			}
			ab01ap01Round, _ := json.Marshal(ab01ap01Output["context"])
			ab01ap01Original, _ := json.Marshal(ab01ap01Context)
			if string(ab01ap01Round) != string(ab01ap01Original) {
				ab01ap01T.Fatal("context changed")
			}
			if !ab01ap01Opaque {
				if _, ab01ap01Exists := ab01ap01Output["capability"]; ab01ap01Exists {
					ab01ap01T.Fatal("unsupported formula")
				}
			}
		}
		ab01ap01Value := ab01ap01Input()
		delete(ab01ap01Value["context"].(map[string]any), ab01ap01Key)
		ab01ap01Output := ab01ap01Evaluate(ab01ap01T, ab01ap01Value)
		if _, ab01ap01Exists := ab01ap01Output["context"].(map[string]any)[ab01ap01Key]; ab01ap01Exists {
			ab01ap01T.Fatal("missing became present")
		}
	}
	for _, ab01ap01Key := range []string{"alpha", "beta"} {
		for _, ab01ap01Bytes := range [][]any{nil, {}, {0}, {0, 1, 2, 3, 4}} {
			ab01ap01Value := ab01ap01Input()
			ab01ap01Raw := ab01ap01Value["fields"].(map[string]any)
			if ab01ap01Bytes == nil {
				delete(ab01ap01Raw, ab01ap01Key)
			} else {
				ab01ap01Raw[ab01ap01Key] = ab01ap01Bytes
			}
			ab01ap01Output := ab01ap01Evaluate(ab01ap01T, ab01ap01Value)
			ab01ap01Want, _ := json.Marshal(ab01ap01Raw)
			ab01ap01Got, _ := json.Marshal(ab01ap01Output["raw"])
			if string(ab01ap01Want) != string(ab01ap01Got) {
				ab01ap01T.Fatal("raw changed")
			}
			ab01ap01Expected := "complete"
			if ab01ap01Bytes != nil && len(ab01ap01Bytes) == 0 {
				ab01ap01Expected = "partial"
			}
			if ab01ap01Output["state"] != ab01ap01Expected || ab01ap01Output["capability"].(map[string]any)["updated"] != float64(2) {
				ab01ap01T.Fatal("independent result lost")
			}
		}
	}
}

func TestAb01Ap01Rejections(ab01ap01T *testing.T) {
	for _, ab01ap01Text := range []string{``, `{`, `[]`, `{} {}`, `{"x":1,"x":2}`, `{"x":[{"a":0,"a":1}]}`, `{"context":{},"fields":{},"arguments":{"extra":{"x":1,"x":2}}}`} {
		ab01ap01Result, ab01ap01Err := Ab01Ap01Evaluate(json.RawMessage(ab01ap01Text))
		if ab01ap01Err == nil || ab01ap01Result != nil {
			ab01ap01T.Fatalf("accepted %s", ab01ap01Text)
		}
	}
	for _, ab01ap01Key := range []string{"registered", "prior", "result"} {
		for _, ab01ap01Bad := range []any{nil, false, "1", 1.5, -1, 4294967296.0, []any{}, map[string]any{}} {
			ab01ap01Value := ab01ap01Input()
			ab01ap01Value["arguments"].(map[string]any)[ab01ap01Key] = ab01ap01Bad
			ab01ap01Reject(ab01ap01T, ab01ap01Value)
		}
		ab01ap01Value := ab01ap01Input()
		delete(ab01ap01Value["arguments"].(map[string]any), ab01ap01Key)
		ab01ap01Reject(ab01ap01T, ab01ap01Value)
	}
	for _, ab01ap01Key := range []string{"alpha", "beta"} {
		for _, ab01ap01Bad := range []any{nil, false, 0, "", []any{false}, []any{-1}, []any{256}, []any{1.5}, []any{nil}} {
			ab01ap01Value := ab01ap01Input()
			ab01ap01Value["fields"].(map[string]any)[ab01ap01Key] = ab01ap01Bad
			ab01ap01Reject(ab01ap01T, ab01ap01Value)
		}
	}
	for _, ab01ap01Size := range []int{8, 10} {
		ab01ap01Value := ab01ap01Input()
		ab01ap01Value["context"].(map[string]any)["view"] = "select"
		ab01ap01Rows := ab01ap01Rows()
		if ab01ap01Size == 8 {
			ab01ap01Rows = ab01ap01Rows[:8]
		} else {
			ab01ap01Rows = append(ab01ap01Rows, map[string]any{"a": 0, "b": 0, "c": 0})
		}
		ab01ap01Value["arguments"] = map[string]any{"rows": ab01ap01Rows}
		ab01ap01Reject(ab01ap01T, ab01ap01Value)
	}
	for ab01ap01Index := 0; ab01ap01Index < 9; ab01ap01Index++ {
		for _, ab01ap01Key := range []string{"a", "b", "c"} {
			for _, ab01ap01Bad := range []any{nil, false, "2", -1, 1.5, 4294967296.0} {
				ab01ap01Value := ab01ap01Input()
				ab01ap01Value["context"].(map[string]any)["view"] = "select"
				ab01ap01Rows := ab01ap01Rows()
				ab01ap01Rows[ab01ap01Index].(map[string]any)[ab01ap01Key] = ab01ap01Bad
				ab01ap01Value["arguments"] = map[string]any{"rows": ab01ap01Rows}
				ab01ap01Reject(ab01ap01T, ab01ap01Value)
			}
		}
	}
	for _, ab01ap01Key := range []string{"context", "fields", "arguments"} {
		for _, ab01ap01Bad := range []any{nil, false, []any{}, 0, ""} {
			ab01ap01Value := ab01ap01Input()
			ab01ap01Value[ab01ap01Key] = ab01ap01Bad
			ab01ap01Reject(ab01ap01T, ab01ap01Value)
		}
	}
	ab01ap01Value := ab01ap01Input()
	delete(ab01ap01Value, "arguments")
	ab01ap01Reject(ab01ap01T, ab01ap01Value)
}

func ab01ap01Reject(ab01ap01T *testing.T, ab01ap01Value any) {
	ab01ap01T.Helper()
	ab01ap01Bytes, _ := json.Marshal(ab01ap01Value)
	ab01ap01Result, ab01ap01Err := Ab01Ap01Evaluate(ab01ap01Bytes)
	if ab01ap01Err == nil || ab01ap01Result != nil {
		ab01ap01T.Fatalf("accepted %s", ab01ap01Bytes)
	}
}

func TestAb01Ap01Isolation(ab01ap01T *testing.T) {
	ab01ap01Bytes, _ := json.Marshal(ab01ap01Input())
	ab01ap01Original := append([]byte(nil), ab01ap01Bytes...)
	ab01ap01Result, ab01ap01Err := Ab01Ap01Evaluate(ab01ap01Bytes)
	if ab01ap01Err != nil {
		ab01ap01T.Fatal(ab01ap01Err)
	}
	ab01ap01Expected := append([]byte(nil), ab01ap01Result...)
	for ab01ap01Index := range ab01ap01Bytes {
		ab01ap01Bytes[ab01ap01Index] = ' '
	}
	if !reflect.DeepEqual([]byte(ab01ap01Result), ab01ap01Expected) {
		ab01ap01T.Fatal("input aliases output")
	}
	ab01ap01Decoded := map[string]any{}
	if ab01ap01Err := json.Unmarshal(ab01ap01Result, &ab01ap01Decoded); ab01ap01Err != nil {
		ab01ap01T.Fatal(ab01ap01Err)
	}
	ab01ap01Decoded["raw"].(map[string]any)["alpha"].([]any)[0] = 99
	ab01ap01Decoded["context"].(map[string]any)["variant"].(map[string]any)["opaque"].([]any)[0] = true
	ab01ap01Decoded["arguments"].(map[string]any)["extra"].(map[string]any)["opaque"].([]any)[0] = true
	if ab01ap01Decoded["fields"].(map[string]any)["alpha"].(map[string]any)["raw"].([]any)[0] != float64(0) {
		ab01ap01T.Fatal("nested output aliases")
	}
	for ab01ap01Index := range ab01ap01Result {
		ab01ap01Result[ab01ap01Index] = ' '
	}
	ab01ap01Again, ab01ap01Err := Ab01Ap01Evaluate(ab01ap01Original)
	if ab01ap01Err != nil || string(ab01ap01Again) != string(ab01ap01Expected) {
		ab01ap01T.Fatal("later call changed")
	}
}

func TestAb01Ap01AdditionalBoundaries(ab01ap01T *testing.T) {
	for _, ab01ap01Key := range []string{"context", "fields"} {
		ab01ap01Value := ab01ap01Input()
		delete(ab01ap01Value, ab01ap01Key)
		ab01ap01Reject(ab01ap01T, ab01ap01Value)
	}
	for _, ab01ap01Key := range []string{"primary", "secondary", "operation"} {
		for _, ab01ap01Bad := range []any{-1, 256, 1.5, []any{}, map[string]any{}} {
			ab01ap01Value := ab01ap01Input()
			ab01ap01Value["context"].(map[string]any)[ab01ap01Key] = ab01ap01Bad
			ab01ap01Reject(ab01ap01T, ab01ap01Value)
		}
	}
	for _, ab01ap01Key := range []string{"rows", "a", "b", "c"} {
		ab01ap01Value := ab01ap01Input()
		ab01ap01Value["context"].(map[string]any)["view"] = "select"
		ab01ap01Rows := ab01ap01Rows()
		ab01ap01Value["arguments"] = map[string]any{"rows": ab01ap01Rows}
		if ab01ap01Key == "rows" {
			delete(ab01ap01Value["arguments"].(map[string]any), "rows")
		} else {
			delete(ab01ap01Rows[8].(map[string]any), ab01ap01Key)
		}
		ab01ap01Reject(ab01ap01T, ab01ap01Value)
	}
	for _, ab01ap01Bad := range []any{nil, false, 0, "", []any{nil, nil, nil, nil, nil, nil, nil, nil, nil}} {
		ab01ap01Value := ab01ap01Input()
		ab01ap01Value["context"].(map[string]any)["view"] = "select"
		ab01ap01Value["arguments"] = map[string]any{"rows": ab01ap01Bad}
		ab01ap01Reject(ab01ap01T, ab01ap01Value)
	}
	for ab01ap01Index := 0; ab01ap01Index < 9; ab01ap01Index++ {
		for _, ab01ap01Key := range []string{"a", "b", "c"} {
			ab01ap01Value := ab01ap01Input()
			ab01ap01Value["context"].(map[string]any)["view"] = "select"
			ab01ap01Rows := ab01ap01Rows()
			ab01ap01Rows[ab01ap01Index] = map[string]any{"a": 2, "b": 2, "c": 2}
			ab01ap01Rows[ab01ap01Index].(map[string]any)[ab01ap01Key] = uint64(4294967295)
			ab01ap01Value["arguments"] = map[string]any{"rows": ab01ap01Rows}
			if ab01ap01Evaluate(ab01ap01T, ab01ap01Value)["select"].(map[string]any)["index"] != float64(ab01ap01Index) {
				ab01ap01T.Fatal("u32 boundary")
			}
		}
	}
	ab01ap01Value := ab01ap01Input()
	ab01ap01Value["context"].(map[string]any)["view"] = "unsupported"
	delete(ab01ap01Value, "arguments")
	ab01ap01Output := ab01ap01Evaluate(ab01ap01T, ab01ap01Value)
	if _, ab01ap01Present := ab01ap01Output["arguments"]; ab01ap01Present || ab01ap01Output["state"] != "Unknown" {
		ab01ap01T.Fatal("absent arguments changed")
	}
	for _, ab01ap01Key := range []string{"alpha", "beta"} {
		ab01ap01Field := ab01ap01Output["fields"].(map[string]any)[ab01ap01Key].(map[string]any)
		if ab01ap01Field["state"] != "Unknown" || ab01ap01Field["unmodeled"] != true {
			ab01ap01T.Fatal("raw-only boundary")
		}
		for _, ab01ap01Invented := range []string{"number", "label"} {
			if _, ab01ap01Present := ab01ap01Field[ab01ap01Invented]; ab01ap01Present {
				ab01ap01T.Fatal("invented semantics")
			}
		}
	}
	for _, ab01ap01Text := range []string{
		`{"context":{"primary":181,"primary":0},"fields":{}}`,
		`{"context":{},"fields":{"alpha":[0],"alpha":[1]}}`,
		`{"context":{},"fields":{},"arguments":{"rows":[{"a":0,"a":1}]}}`,
		`{"context":{},"fields":{}} trailing`,
		`{"context":{},"fields":{"alpha":[1e0]}}`,
	} {
		ab01ap01Result, ab01ap01Err := Ab01Ap01Evaluate(json.RawMessage(ab01ap01Text))
		if ab01ap01Err == nil || ab01ap01Result != nil {
			ab01ap01T.Fatal("invalid JSON accepted")
		}
	}
}
