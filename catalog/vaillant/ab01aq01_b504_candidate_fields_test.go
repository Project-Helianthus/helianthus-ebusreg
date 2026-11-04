package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type ab01aq01Case struct {
	Ab01Aq01ID       string          `json:"id"`
	Ab01Aq01Input    json.RawMessage `json:"input"`
	Ab01Aq01Expected map[string]any  `json:"expected"`
	Ab01Aq01State    string          `json:"state"`
}

func ab01aq01Cases(ab01aq01T *testing.T) []ab01aq01Case {
	ab01aq01T.Helper()
	ab01aq01Raw, ab01aq01Err := os.ReadFile("testdata/ab01aq01_b504_candidate_fields.json")
	if ab01aq01Err != nil {
		ab01aq01T.Fatal(ab01aq01Err)
	}
	var ab01aq01Cases []ab01aq01Case
	if ab01aq01Err = json.Unmarshal(ab01aq01Raw, &ab01aq01Cases); ab01aq01Err != nil {
		ab01aq01T.Fatal(ab01aq01Err)
	}
	if len(ab01aq01Cases) == 0 {
		ab01aq01T.Fatal("zero fixture cases")
	}
	return ab01aq01Cases
}

func ab01aq01Evaluate(ab01aq01T *testing.T, ab01aq01Input any) map[string]any {
	ab01aq01T.Helper()
	ab01aq01Raw, ab01aq01Err := json.Marshal(ab01aq01Input)
	if ab01aq01Err != nil {
		ab01aq01T.Fatal(ab01aq01Err)
	}
	ab01aq01Output, ab01aq01Err := Ab01Aq01Evaluate(ab01aq01Raw)
	if ab01aq01Err != nil {
		ab01aq01T.Fatal(ab01aq01Err)
	}
	var ab01aq01Value map[string]any
	if ab01aq01Err = json.Unmarshal(ab01aq01Output, &ab01aq01Value); ab01aq01Err != nil {
		ab01aq01T.Fatal(ab01aq01Err)
	}
	if ab01aq01Value["physical"] != "Unknown" || ab01aq01Value["qualified"] != false {
		ab01aq01T.Fatal(ab01aq01Value)
	}
	return ab01aq01Value
}

func ab01aq01Decode(ab01aq01T *testing.T, ab01aq01Raw json.RawMessage) map[string]any {
	ab01aq01T.Helper()
	var ab01aq01Input map[string]any
	if ab01aq01Err := json.Unmarshal(ab01aq01Raw, &ab01aq01Input); ab01aq01Err != nil {
		ab01aq01T.Fatal(ab01aq01Err)
	}
	return ab01aq01Input
}

func TestAb01Aq01Fixture(ab01aq01T *testing.T) {
	for _, ab01aq01Case := range ab01aq01Cases(ab01aq01T) {
		ab01aq01T.Run(ab01aq01Case.Ab01Aq01ID, func(ab01aq01T *testing.T) {
			ab01aq01Input := ab01aq01Decode(ab01aq01T, ab01aq01Case.Ab01Aq01Input)
			ab01aq01Output := ab01aq01Evaluate(ab01aq01T, ab01aq01Input)
			if ab01aq01Output["state"] != ab01aq01Case.Ab01Aq01State || !reflect.DeepEqual(ab01aq01Output["results"], ab01aq01Case.Ab01Aq01Expected) {
				ab01aq01T.Fatalf("got %v; expected %v (%s)", ab01aq01Output, ab01aq01Case.Ab01Aq01Expected, ab01aq01Case.Ab01Aq01State)
			}
			for ab01aq01Key, ab01aq01Value := range ab01aq01Input {
				if !reflect.DeepEqual(ab01aq01Output[ab01aq01Key], ab01aq01Value) {
					ab01aq01T.Fatalf("lost %s", ab01aq01Key)
				}
			}
		})
	}
}

func TestAb01Aq01ContextAxes(ab01aq01T *testing.T) {
	for _, ab01aq01Case := range ab01aq01Cases(ab01aq01T) {
		for _, ab01aq01Key := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
			for _, ab01aq01Mode := range []string{"changed", "absent", "zero"} {
				ab01aq01T.Run(ab01aq01Case.Ab01Aq01ID+"/"+ab01aq01Key+"/"+ab01aq01Mode, func(ab01aq01T *testing.T) {
					ab01aq01Input := ab01aq01Decode(ab01aq01T, ab01aq01Case.Ab01Aq01Input)
					ab01aq01Context := ab01aq01Input["context"].(map[string]any)
					ab01aq01Opaque := ab01aq01Key == "family" || ab01aq01Key == "variant"
					switch ab01aq01Key {
					case "primary", "secondary", "operation":
						ab01aq01Context[ab01aq01Key] = float64(0)
						if ab01aq01Mode == "changed" {
							ab01aq01Context[ab01aq01Key] = float64(7)
						}
					default:
						ab01aq01Context[ab01aq01Key] = ""
						if ab01aq01Mode == "changed" {
							ab01aq01Context[ab01aq01Key] = "other"
						}
					}
					if ab01aq01Mode == "absent" {
						delete(ab01aq01Context, ab01aq01Key)
					}
					ab01aq01Output := ab01aq01Evaluate(ab01aq01T, ab01aq01Input)
					if !reflect.DeepEqual(ab01aq01Output["context"], ab01aq01Context) {
						ab01aq01T.Fatal("context lost")
					}
					if ab01aq01Opaque {
						if !reflect.DeepEqual(ab01aq01Output["results"], ab01aq01Case.Ab01Aq01Expected) {
							ab01aq01T.Fatal("opaque identity selected semantics")
						}
						return
					}
					if ab01aq01Output["state"] != "Unknown" {
						ab01aq01T.Fatal(ab01aq01Output)
					}
					for _, ab01aq01Result := range ab01aq01Output["results"].(map[string]any) {
						ab01aq01Field := ab01aq01Result.(map[string]any)
						if ab01aq01Field["state"] != "Unknown" || len(ab01aq01Field) != 3 {
							ab01aq01T.Fatal("unsupported conversion", ab01aq01Field)
						}
					}
				})
			}
		}
	}
}

func TestAb01Aq01WidthsAndPresence(ab01aq01T *testing.T) {
	for _, ab01aq01Case := range ab01aq01Cases(ab01aq01T) {
		if ab01aq01Case.Ab01Aq01State != "complete" {
			continue
		}
		ab01aq01Original := ab01aq01Decode(ab01aq01T, ab01aq01Case.Ab01Aq01Input)
		for ab01aq01Key := range ab01aq01Original["fields"].(map[string]any) {
			for _, ab01aq01Mode := range []string{"missing", "empty", "long"} {
				ab01aq01T.Run(ab01aq01Case.Ab01Aq01ID+"/"+ab01aq01Key+"/"+ab01aq01Mode, func(ab01aq01T *testing.T) {
					ab01aq01Input := ab01aq01Decode(ab01aq01T, ab01aq01Case.Ab01Aq01Input)
					ab01aq01Fields := ab01aq01Input["fields"].(map[string]any)
					switch ab01aq01Mode {
					case "missing":
						delete(ab01aq01Fields, ab01aq01Key)
					case "empty":
						ab01aq01Fields[ab01aq01Key] = []any{}
					case "long":
						ab01aq01Fields[ab01aq01Key] = []any{float64(0), float64(0), float64(0)}
					}
					ab01aq01Output := ab01aq01Evaluate(ab01aq01T, ab01aq01Input)
					ab01aq01Results := ab01aq01Output["results"].(map[string]any)
					if ab01aq01Key == "dhw_request" && ab01aq01Mode == "missing" {
						if _, ab01aq01Present := ab01aq01Results[ab01aq01Key]; ab01aq01Present {
							ab01aq01T.Fatal("optional absent became present")
						}
						return
					}
					ab01aq01Result := ab01aq01Results[ab01aq01Key].(map[string]any)
					if ab01aq01Output["state"] != "partial" || ab01aq01Result["state"] != "incomplete" || ab01aq01Result["number"] != nil {
						ab01aq01T.Fatal(ab01aq01Output)
					}
					if ab01aq01Result["present"] != (ab01aq01Mode != "missing") || !reflect.DeepEqual(ab01aq01Result["raw"], ab01aq01Fields[ab01aq01Key]) {
						ab01aq01T.Fatal("presence/raw lost")
					}
					for ab01aq01Sibling, ab01aq01Expected := range ab01aq01Case.Ab01Aq01Expected {
						if ab01aq01Sibling != ab01aq01Key && !reflect.DeepEqual(ab01aq01Results[ab01aq01Sibling], ab01aq01Expected) {
							ab01aq01T.Fatal("sibling erased")
						}
					}
				})
			}
		}
	}
}

func TestAb01Aq01RejectInvalid(ab01aq01T *testing.T) {
	for ab01aq01Index, ab01aq01Raw := range []string{
		``, `null`, `[]`, `{`, `{} {}`, `{"context":{},"fields":{},"x":{"a":1,"a":2}}`,
		`{"context":{},"fields":{},"x":[{"a":1,"\u0061":2}]}`, `{"context":{},"fields":{},"x":1e999}`, `{"context":{},"fields":{},"x":NaN}`,
	} {
		ab01aq01T.Run(string(rune('a'+ab01aq01Index)), func(ab01aq01T *testing.T) {
			ab01aq01Output, ab01aq01Err := Ab01Aq01Evaluate([]byte(ab01aq01Raw))
			if ab01aq01Err == nil || ab01aq01Output != nil {
				ab01aq01T.Fatal("accepted invalid input", ab01aq01Raw)
			}
		})
	}
	for _, ab01aq01Key := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant", "fields", "context", "arguments", "temperature", "divisor", "role"} {
		ab01aq01Invalid := []any{nil, true, "1", 1.5, -1, 256, []any{}, map[string]any{}}
		switch ab01aq01Key {
		case "divisor":
			ab01aq01Invalid = []any{nil, true, "1", 1.5, -1, 0, 15, 255, 257, []any{}, map[string]any{}}
		case "role", "direction", "view", "family", "variant":
			ab01aq01Invalid = []any{nil, true, 1.5, -1, 256, []any{}, map[string]any{}}
		case "fields", "context", "arguments":
			ab01aq01Invalid = []any{nil, true, "1", 1.5, -1, 256, []any{}}
		}
		for _, ab01aq01Value := range ab01aq01Invalid {
			ab01aq01T.Run(ab01aq01Key+"/"+ab01aq01Describe(ab01aq01Value), func(ab01aq01T *testing.T) {
				ab01aq01Input := ab01aq01Decode(ab01aq01T, ab01aq01Cases(ab01aq01T)[0].Ab01Aq01Input)
				switch ab01aq01Key {
				case "fields", "context", "arguments":
					ab01aq01Input[ab01aq01Key] = ab01aq01Value
				case "temperature":
					ab01aq01Input["fields"].(map[string]any)[ab01aq01Key] = []any{ab01aq01Value, 0}
				case "divisor":
					ab01aq01Input["arguments"].(map[string]any)[ab01aq01Key] = ab01aq01Value
				case "role":
					ab01aq01Input["context"].(map[string]any)["view"] = "circuit-role"
					ab01aq01Input["context"].(map[string]any)["operation"] = 13
					ab01aq01Input["arguments"] = map[string]any{"role": ab01aq01Value}
				default:
					ab01aq01Input["context"].(map[string]any)[ab01aq01Key] = ab01aq01Value
				}
				ab01aq01Raw, _ := json.Marshal(ab01aq01Input)
				ab01aq01Output, ab01aq01Err := Ab01Aq01Evaluate(ab01aq01Raw)
				if ab01aq01Err == nil || ab01aq01Output != nil {
					ab01aq01T.Fatalf("accepted %s", ab01aq01Raw)
				}
			})
		}
	}
}

func ab01aq01Describe(ab01aq01Value any) string {
	ab01aq01Raw, _ := json.Marshal(ab01aq01Value)
	return string(ab01aq01Raw)
}

func TestAb01Aq01ArgumentsAndMemory(ab01aq01T *testing.T) {
	for _, ab01aq01Index := range []int{0, 14} {
		ab01aq01Input := ab01aq01Decode(ab01aq01T, ab01aq01Cases(ab01aq01T)[ab01aq01Index].Ab01Aq01Input)
		delete(ab01aq01Input, "arguments")
		ab01aq01Raw, _ := json.Marshal(ab01aq01Input)
		if ab01aq01Output, ab01aq01Err := Ab01Aq01Evaluate(ab01aq01Raw); ab01aq01Err == nil || ab01aq01Output != nil {
			ab01aq01T.Fatal("missing required argument")
		}
	}
	ab01aq01Input := ab01aq01Decode(ab01aq01T, ab01aq01Cases(ab01aq01T)[0].Ab01Aq01Input)
	ab01aq01Input["extra"] = map[string]any{"absent": nil, "zero": float64(0), "false": false, "empty": []any{}, "object": map[string]any{}}
	ab01aq01Input["fields"].(map[string]any)["opaque"] = []any{}
	ab01aq01Expected := ab01aq01Evaluate(ab01aq01T, ab01aq01Input)
	if !reflect.DeepEqual(ab01aq01Expected["extra"], ab01aq01Input["extra"]) || !reflect.DeepEqual(ab01aq01Expected["fields"], ab01aq01Input["fields"]) {
		ab01aq01T.Fatal("extra values lost")
	}
	ab01aq01Raw, _ := json.Marshal(ab01aq01Input)
	ab01aq01Original := append([]byte(nil), ab01aq01Raw...)
	ab01aq01Output, ab01aq01Err := Ab01Aq01Evaluate(ab01aq01Raw)
	if ab01aq01Err != nil {
		ab01aq01T.Fatal(ab01aq01Err)
	}
	for ab01aq01Index := range ab01aq01Raw {
		ab01aq01Raw[ab01aq01Index] = ' '
	}
	if !reflect.DeepEqual(ab01aq01Decode(ab01aq01T, ab01aq01Output), ab01aq01Expected) {
		ab01aq01T.Fatal("input aliases output")
	}
	for ab01aq01Index := range ab01aq01Output {
		ab01aq01Output[ab01aq01Index] = ' '
	}
	if !reflect.DeepEqual(ab01aq01Evaluate(ab01aq01T, json.RawMessage(ab01aq01Original)), ab01aq01Expected) {
		ab01aq01T.Fatal("call state retained")
	}
}

func TestAb01Aq01InputBoundaries(ab01aq01T *testing.T) {
	for _, ab01aq01Raw := range []string{
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
		ab01aq01Output, ab01aq01Err := Ab01Aq01Evaluate([]byte(ab01aq01Raw))
		if ab01aq01Err == nil || ab01aq01Output != nil {
			ab01aq01T.Fatalf("accepted %s", ab01aq01Raw)
		}
	}
	for _, ab01aq01Index := range []int{0, 4, 7, 10, 14} {
		ab01aq01Input := ab01aq01Decode(ab01aq01T, ab01aq01Cases(ab01aq01T)[ab01aq01Index].Ab01Aq01Input)
		ab01aq01Input["arguments"] = map[string]any{"unexpected": false}
		ab01aq01Raw, _ := json.Marshal(ab01aq01Input)
		if ab01aq01Output, ab01aq01Err := Ab01Aq01Evaluate(ab01aq01Raw); ab01aq01Err == nil || ab01aq01Output != nil {
			ab01aq01T.Fatal("unexpected argument accepted")
		}
	}
	ab01aq01Input := ab01aq01Decode(ab01aq01T, ab01aq01Cases(ab01aq01T)[0].Ab01Aq01Input)
	ab01aq01Input["context"].(map[string]any)["opaque"] = map[string]any{"null": nil, "zero": float64(0), "false": false, "empty": ""}
	ab01aq01Input["fields"].(map[string]any)["tail"] = []any{float64(255), float64(0)}
	ab01aq01Output := ab01aq01Evaluate(ab01aq01T, ab01aq01Input)
	if !reflect.DeepEqual(ab01aq01Output["context"], ab01aq01Input["context"]) {
		ab01aq01T.Fatal("opaque context lost")
	}
	ab01aq01Tail := ab01aq01Output["results"].(map[string]any)["tail"].(map[string]any)
	if ab01aq01Tail["state"] != "Unknown" || len(ab01aq01Tail) != 3 {
		ab01aq01T.Fatal("tail interpreted")
	}
	ab01aq01Input["context"].(map[string]any)["operation"] = float64(0)
	ab01aq01Input["arguments"] = map[string]any{"extra": nil, "zero": float64(0), "false": false, "empty": []any{}}
	ab01aq01Output = ab01aq01Evaluate(ab01aq01T, ab01aq01Input)
	if ab01aq01Output["state"] != "Unknown" || !reflect.DeepEqual(ab01aq01Output["arguments"], ab01aq01Input["arguments"]) {
		ab01aq01T.Fatal("unknown arguments lost")
	}
}
