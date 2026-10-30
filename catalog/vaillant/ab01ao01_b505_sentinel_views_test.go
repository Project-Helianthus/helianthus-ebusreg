package vaillant

import (
	ab01ao01Bytes "bytes"
	ab01ao01JSON "encoding/json"
	ab01ao01OS "os"
	ab01ao01Reflect "reflect"
	ab01ao01Testing "testing"
)

func ab01ao01Control() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 5, "operation": 50, "direction": "request", "view": "replacement-pair", "family": "synthetic-family", "variant": map[string]any{"opaque": []any{0, "variant"}}}, "fields": map[string]any{"temperature_one": []int{0, 128}, "temperature_two": []int{0, 0}, "unmodelled": []int{11, 255, 0}}}
}

func ab01ao01Call(ab01ao01T *ab01ao01Testing.T, ab01ao01Input any) map[string]any {
	ab01ao01T.Helper()
	ab01ao01Encoded, ab01ao01Error := ab01ao01JSON.Marshal(ab01ao01Input)
	if ab01ao01Error != nil {
		ab01ao01T.Fatal(ab01ao01Error)
	}
	ab01ao01Output, ab01ao01Error := Ab01Ao01Interpret(ab01ao01Encoded)
	if ab01ao01Error != nil {
		ab01ao01T.Fatal(ab01ao01Error)
	}
	var ab01ao01Result map[string]any
	if ab01ao01Error := ab01ao01JSON.Unmarshal(ab01ao01Output, &ab01ao01Result); ab01ao01Error != nil {
		ab01ao01T.Fatal(ab01ao01Error)
	}
	var ab01ao01Original map[string]any
	if ab01ao01Error := ab01ao01JSON.Unmarshal(ab01ao01Encoded, &ab01ao01Original); ab01ao01Error != nil {
		ab01ao01T.Fatal(ab01ao01Error)
	}
	if !ab01ao01Reflect.DeepEqual(ab01ao01Original["context"], ab01ao01Result["context"]) || !ab01ao01Reflect.DeepEqual(ab01ao01Original["fields"], ab01ao01Result["raw"]) {
		ab01ao01T.Fatal("context or raw evidence changed")
	}
	if ab01ao01Result["physical"] != "Unknown" || ab01ao01Result["qualified"] != false {
		ab01ao01T.Fatal("native qualification inferred")
	}
	ab01ao01Argument, ab01ao01Present := ab01ao01Original["arguments"]
	ab01ao01Returned, ab01ao01ReturnedPresent := ab01ao01Result["arguments"]
	if ab01ao01Present != ab01ao01ReturnedPresent || !ab01ao01Reflect.DeepEqual(ab01ao01Argument, ab01ao01Returned) {
		ab01ao01T.Fatal("arguments presence changed")
	}
	for _, ab01ao01Value := range ab01ao01Result["fields"].(map[string]any) {
		ab01ao01Field := ab01ao01Value.(map[string]any)
		for _, ab01ao01Key := range []string{"number", "calculation", "candidate_domain", "flags", "reserved"} {
			if _, ab01ao01Present := ab01ao01Field[ab01ao01Key]; ab01ao01Present {
				ab01ao01T.Fatalf("invented %s", ab01ao01Key)
			}
		}
	}
	return ab01ao01Result
}

func TestAb01Ao01IndependentWords(ab01ao01T *ab01ao01Testing.T) {
	ab01ao01Data, ab01ao01Error := ab01ao01OS.ReadFile("testdata/ab01ao01_b505_sentinel_views.json")
	if ab01ao01Error != nil {
		ab01ao01T.Fatal(ab01ao01Error)
	}
	var ab01ao01Cases []struct {
		Ab01Ao01ID    string         `json:"id"`
		Ab01Ao01Input map[string]any `json:"input"`
		Ab01Ao01Field string         `json:"field"`
		Ab01Ao01Label any            `json:"label"`
		Ab01Ao01State string         `json:"state"`
	}
	if ab01ao01Error := ab01ao01JSON.Unmarshal(ab01ao01Data, &ab01ao01Cases); ab01ao01Error != nil {
		ab01ao01T.Fatal(ab01ao01Error)
	}
	if len(ab01ao01Cases) != 16 {
		ab01ao01T.Fatal("fixture cases missing")
	}
	for _, ab01ao01Case := range ab01ao01Cases {
		ab01ao01T.Run(ab01ao01Case.Ab01Ao01ID, func(ab01ao01T *ab01ao01Testing.T) {
			ab01ao01Result := ab01ao01Call(ab01ao01T, ab01ao01Case.Ab01Ao01Input)
			ab01ao01Fields := ab01ao01Result["fields"].(map[string]any)
			if len(ab01ao01Fields) != 2 {
				ab01ao01T.Fatal("independent result fields lost")
			}
			ab01ao01Field := ab01ao01Fields[ab01ao01Case.Ab01Ao01Field].(map[string]any)
			if ab01ao01Field["state"] != ab01ao01Case.Ab01Ao01State || ab01ao01Field["label"] != ab01ao01Case.Ab01Ao01Label {
				ab01ao01T.Fatalf("wrong branch: %v", ab01ao01Field)
			}
			ab01ao01Raw := ab01ao01Case.Ab01Ao01Input["fields"].(map[string]any)[ab01ao01Case.Ab01Ao01Field]
			if !ab01ao01Reflect.DeepEqual(ab01ao01Raw, ab01ao01Field["raw"]) {
				ab01ao01T.Fatal("field raw changed")
			}
			ab01ao01Sibling := "temperature_two"
			ab01ao01Label := "circuit-default"
			if ab01ao01Case.Ab01Ao01Field == "temperature_two" {
				ab01ao01Sibling = "temperature_one"
				ab01ao01Label = "preserve-existing"
			}
			ab01ao01Other := ab01ao01Fields[ab01ao01Sibling].(map[string]any)
			if ab01ao01Other["state"] != "Hypothesis" || ab01ao01Other["label"] != ab01ao01Label {
				ab01ao01T.Fatal("independent valid sibling lost")
			}
			ab01ao01ExpectedState := "partial"
			if ab01ao01Case.Ab01Ao01State == "Hypothesis" {
				ab01ao01ExpectedState = "Hypothesis"
			}
			if ab01ao01Result["state"] != ab01ao01ExpectedState {
				ab01ao01T.Fatal("wrong aggregate state", ab01ao01Result["state"])
			}
			if ab01ao01Case.Ab01Ao01Label != nil && ab01ao01Field["grade"] != "Hypothesis" {
				ab01ao01T.Fatal("semantic grade changed")
			}
		})
	}
	ab01ao01Input := ab01ao01Control()
	ab01ao01Input["fields"].(map[string]any)["temperature_one"] = []int{0, 0}
	ab01ao01Input["fields"].(map[string]any)["temperature_two"] = []int{0, 128}
	ab01ao01Result := ab01ao01Call(ab01ao01T, ab01ao01Input)["fields"].(map[string]any)
	if ab01ao01Result["temperature_one"].(map[string]any)["label"] != "circuit-default" || ab01ao01Result["temperature_two"].(map[string]any)["label"] != "preserve-existing" {
		ab01ao01T.Fatal("reverse pair failed")
	}
}

func TestAb01Ao01SelectionAxes(ab01ao01T *ab01ao01Testing.T) {
	for ab01ao01Key, ab01ao01Value := range map[string]any{"primary": 180, "secondary": 6, "operation": 49, "direction": "response", "view": "other"} {
		ab01ao01T.Run(ab01ao01Key, func(ab01ao01T *ab01ao01Testing.T) {
			for _, ab01ao01Mode := range []string{"different", "missing", "empty-or-zero"} {
				ab01ao01Input := ab01ao01Control()
				ab01ao01Context := ab01ao01Input["context"].(map[string]any)
				switch ab01ao01Mode {
				case "different":
					ab01ao01Context[ab01ao01Key] = ab01ao01Value
				case "missing":
					delete(ab01ao01Context, ab01ao01Key)
				default:
					if _, ab01ao01OK := ab01ao01Value.(string); ab01ao01OK {
						ab01ao01Context[ab01ao01Key] = ""
					} else {
						ab01ao01Context[ab01ao01Key] = 0
					}
				}
				ab01ao01Result := ab01ao01Call(ab01ao01T, ab01ao01Input)
				if ab01ao01Result["state"] != "Unknown" || len(ab01ao01Result["fields"].(map[string]any)) != 0 {
					ab01ao01T.Fatal("selection guessed", ab01ao01Mode)
				}
			}
		})
	}
	for _, ab01ao01Key := range []string{"family", "variant"} {
		for _, ab01ao01Value := range []any{nil, "other-family", 0, []any{}, map[string]any{"opaque": true}} {
			ab01ao01Input := ab01ao01Control()
			ab01ao01Input["context"].(map[string]any)[ab01ao01Key] = ab01ao01Value
			if ab01ao01Call(ab01ao01T, ab01ao01Input)["state"] != "Hypothesis" {
				ab01ao01T.Fatal("opaque context selected view")
			}
		}
		ab01ao01Input := ab01ao01Control()
		delete(ab01ao01Input["context"].(map[string]any), ab01ao01Key)
		if ab01ao01Call(ab01ao01T, ab01ao01Input)["state"] != "Hypothesis" {
			ab01ao01T.Fatal("opaque context required")
		}
	}
	ab01ao01Input := ab01ao01Control()
	ab01ao01Input["arguments"] = map[string]any{}
	ab01ao01Call(ab01ao01T, ab01ao01Input)
}

func TestAb01Ao01UnknownAndPresence(ab01ao01T *ab01ao01Testing.T) {
	for _, ab01ao01Words := range []map[string]any{{}, {"temperature_one": []int{}, "temperature_two": []int{}}, {"temperature_one": []int{1, 0}, "temperature_two": []int{255, 255}}} {
		ab01ao01Input := ab01ao01Control()
		ab01ao01Input["fields"] = ab01ao01Words
		ab01ao01Result := ab01ao01Call(ab01ao01T, ab01ao01Input)
		ab01ao01Expected := "partial"
		if len(ab01ao01Words) == 2 && len(ab01ao01Words["temperature_one"].([]int)) == 2 {
			ab01ao01Expected = "Unknown"
		}
		if ab01ao01Result["state"] != ab01ao01Expected {
			ab01ao01T.Fatal("unknown or incomplete state lost")
		}
	}
	for _, ab01ao01Value := range [][]int{{128, 0}, {0, 1}, {255, 255}, {0, 0, 0}, {}} {
		ab01ao01Input := ab01ao01Control()
		ab01ao01Input["fields"].(map[string]any)["unmodelled"] = ab01ao01Value
		if ab01ao01Call(ab01ao01T, ab01ao01Input)["state"] != "Hypothesis" {
			ab01ao01T.Fatal("unmodelled field selected view")
		}
	}
}

func TestAb01Ao01RejectInvalid(ab01ao01T *ab01ao01Testing.T) {
	for _, ab01ao01Input := range []string{``, `{`, `[]`, `null`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"family":{"x":0,"x":1}},"fields":{}}`, `{"context":{},"fields":{"x":[],"x":[]}}`, `{"context":{},"fields":{},"arguments":{"x":0,"x":1}}`, `{"context":{},"fields":{},"extra":0}`} {
		ab01ao01Output, ab01ao01Error := Ab01Ao01Interpret([]byte(ab01ao01Input))
		if ab01ao01Error == nil || ab01ao01Output != nil {
			ab01ao01T.Fatalf("accepted invalid JSON: %s", ab01ao01Input)
		}
	}
	for _, ab01ao01Key := range []string{"context", "fields", "arguments"} {
		for _, ab01ao01Value := range []any{nil, 0, false, "", []any{}} {
			ab01ao01Input := ab01ao01Control()
			ab01ao01Input[ab01ao01Key] = ab01ao01Value
			ab01ao01Encoded, _ := ab01ao01JSON.Marshal(ab01ao01Input)
			ab01ao01Output, ab01ao01Error := Ab01Ao01Interpret(ab01ao01Encoded)
			if ab01ao01Error == nil || ab01ao01Output != nil {
				ab01ao01T.Fatal("accepted wrong object kind", ab01ao01Key, ab01ao01Value)
			}
		}
	}
	for _, ab01ao01Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, ab01ao01Value := range []any{nil, true, []any{}, map[string]any{}} {
			ab01ao01Input := ab01ao01Control()
			ab01ao01Input["context"].(map[string]any)[ab01ao01Key] = ab01ao01Value
			ab01ao01Encoded, _ := ab01ao01JSON.Marshal(ab01ao01Input)
			ab01ao01Output, ab01ao01Error := Ab01Ao01Interpret(ab01ao01Encoded)
			if ab01ao01Error == nil || ab01ao01Output != nil {
				ab01ao01T.Fatal("accepted invalid selector", ab01ao01Key)
			}
		}
	}
	for _, ab01ao01Key := range []string{"temperature_one", "temperature_two", "unmodelled"} {
		for _, ab01ao01Value := range []any{nil, 0, "AA==", map[string]any{}, []any{-1}, []any{256}, []any{1.5}, []any{true}, []any{nil}, []any{"0"}, []any{[]any{}}, []any{ab01ao01JSON.Number("1e0")}} {
			ab01ao01Input := ab01ao01Control()
			ab01ao01Input["fields"].(map[string]any)[ab01ao01Key] = ab01ao01Value
			ab01ao01Encoded, _ := ab01ao01JSON.Marshal(ab01ao01Input)
			ab01ao01Output, ab01ao01Error := Ab01Ao01Interpret(ab01ao01Encoded)
			if ab01ao01Error == nil || ab01ao01Output != nil {
				ab01ao01T.Fatal("accepted invalid byte array", ab01ao01Key, ab01ao01Value)
			}
		}
	}
	ab01ao01Input := ab01ao01Control()
	ab01ao01Input["arguments"] = map[string]any{"flag": true}
	ab01ao01Encoded, _ := ab01ao01JSON.Marshal(ab01ao01Input)
	if ab01ao01Output, ab01ao01Error := Ab01Ao01Interpret(ab01ao01Encoded); ab01ao01Error == nil || ab01ao01Output != nil {
		ab01ao01T.Fatal("accepted undefined argument")
	}
}

func TestAb01Ao01DetachedMemory(ab01ao01T *ab01ao01Testing.T) {
	ab01ao01Input, _ := ab01ao01JSON.Marshal(ab01ao01Control())
	ab01ao01First, ab01ao01Error := Ab01Ao01Interpret(ab01ao01Input)
	if ab01ao01Error != nil {
		ab01ao01T.Fatal(ab01ao01Error)
	}
	ab01ao01Saved := append([]byte(nil), ab01ao01First...)
	for ab01ao01Index := range ab01ao01Input {
		ab01ao01Input[ab01ao01Index] = ' '
	}
	if !ab01ao01Bytes.Equal(ab01ao01First, ab01ao01Saved) {
		ab01ao01T.Fatal("input aliases output")
	}
	var ab01ao01Nested map[string]any
	if ab01ao01Error := ab01ao01JSON.Unmarshal(ab01ao01First, &ab01ao01Nested); ab01ao01Error != nil {
		ab01ao01T.Fatal(ab01ao01Error)
	}
	ab01ao01Nested["raw"].(map[string]any)["temperature_one"].([]any)[0] = float64(99)
	if ab01ao01Nested["fields"].(map[string]any)["temperature_one"].(map[string]any)["raw"].([]any)[0] != float64(0) {
		ab01ao01T.Fatal("raw aliases semantic field")
	}
	ab01ao01Nested["context"].(map[string]any)["variant"].(map[string]any)["opaque"].([]any)[0] = "changed"
	for ab01ao01Index := range ab01ao01First {
		ab01ao01First[ab01ao01Index] = ' '
	}
	ab01ao01Fresh, _ := ab01ao01JSON.Marshal(ab01ao01Control())
	ab01ao01Second, ab01ao01Error := Ab01Ao01Interpret(ab01ao01Fresh)
	if ab01ao01Error != nil || !ab01ao01Bytes.Equal(ab01ao01Second, ab01ao01Saved) {
		ab01ao01T.Fatal("later call changed")
	}
	if !ab01ao01JSON.Valid(ab01ao01Fresh) {
		ab01ao01T.Fatal("output mutation changed later input")
	}
}

func TestAb01Ao01IntegerAndKinds(ab01ao01T *ab01ao01Testing.T) {
	for _, ab01ao01Key := range []string{"primary", "secondary", "operation"} {
		ab01ao01Input := ab01ao01Control()
		ab01ao01Input["context"].(map[string]any)[ab01ao01Key] = ab01ao01JSON.Number("99999999999999999999999999999999999")
		if ab01ao01Call(ab01ao01T, ab01ao01Input)["state"] != "Unknown" {
			ab01ao01T.Fatal("unknown integer selected view")
		}
		for _, ab01ao01Value := range []any{"50", 1.5, ab01ao01JSON.Number("1e0")} {
			ab01ao01Input["context"].(map[string]any)[ab01ao01Key] = ab01ao01Value
			ab01ao01Encoded, _ := ab01ao01JSON.Marshal(ab01ao01Input)
			if ab01ao01Output, ab01ao01Error := Ab01Ao01Interpret(ab01ao01Encoded); ab01ao01Error == nil || ab01ao01Output != nil {
				ab01ao01T.Fatal("invalid integer accepted", ab01ao01Key)
			}
		}
	}
	for _, ab01ao01Key := range []string{"direction", "view"} {
		ab01ao01Input := ab01ao01Control()
		ab01ao01Input["context"].(map[string]any)[ab01ao01Key] = 0
		ab01ao01Encoded, _ := ab01ao01JSON.Marshal(ab01ao01Input)
		if ab01ao01Output, ab01ao01Error := Ab01Ao01Interpret(ab01ao01Encoded); ab01ao01Error == nil || ab01ao01Output != nil {
			ab01ao01T.Fatal("numeric string selector accepted")
		}
	}
	for _, ab01ao01Key := range []string{"context", "fields"} {
		ab01ao01Input := ab01ao01Control()
		delete(ab01ao01Input, ab01ao01Key)
		ab01ao01Encoded, _ := ab01ao01JSON.Marshal(ab01ao01Input)
		if ab01ao01Output, ab01ao01Error := Ab01Ao01Interpret(ab01ao01Encoded); ab01ao01Error == nil || ab01ao01Output != nil {
			ab01ao01T.Fatal("missing required object accepted")
		}
	}
	ab01ao01Input := ab01ao01Control()
	ab01ao01Input["fields"].(map[string]any)["temperature_one"] = []any{ab01ao01JSON.Number("-0"), 128}
	if ab01ao01Call(ab01ao01T, ab01ao01Input)["state"] != "Hypothesis" {
		ab01ao01T.Fatal("integer negative zero misinterpreted")
	}
}
