package vaillant

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func ab23aq23EvaluateTest(ab23aq23T *testing.T, ab23aq23Input any) map[string]any {
	ab23aq23T.Helper()
	ab23aq23Bytes, ab23aq23Err := json.Marshal(ab23aq23Input)
	if ab23aq23Err != nil {
		ab23aq23T.Fatal(ab23aq23Err)
	}
	ab23aq23Output, ab23aq23Err := Ab23Aq23Evaluate(ab23aq23Bytes)
	if ab23aq23Err != nil {
		ab23aq23T.Fatal(ab23aq23Err)
	}
	ab23aq23Result := map[string]any{}
	if ab23aq23Err := json.Unmarshal(ab23aq23Output, &ab23aq23Result); ab23aq23Err != nil {
		ab23aq23T.Fatal(ab23aq23Err)
	}
	return ab23aq23Result
}

func ab23aq23InputTest(ab23aq23Op int, ab23aq23Direction string, ab23aq23Fields map[string]any) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 18, "operation": ab23aq23Op, "direction": ab23aq23Direction, "view": "actuator-fields", "family": "synthetic-family", "variant": "synthetic-variant"}, "fields": ab23aq23Fields}
}

func TestAb23Aq23Fixture(ab23aq23T *testing.T) {
	ab23aq23Data, ab23aq23Err := os.ReadFile("testdata/ab23aq23_b512_valve_pump_fields.json")
	if ab23aq23Err != nil {
		ab23aq23T.Fatal(ab23aq23Err)
	}
	var ab23aq23Cases []struct {
		Ab23Aq23ID       string                    `json:"id"`
		Ab23Aq23Input    map[string]any            `json:"input"`
		Ab23Aq23Expected map[string]map[string]any `json:"expected"`
	}
	if ab23aq23Err := json.Unmarshal(ab23aq23Data, &ab23aq23Cases); ab23aq23Err != nil {
		ab23aq23T.Fatal(ab23aq23Err)
	}
	if len(ab23aq23Cases) == 0 {
		ab23aq23T.Fatal("zero cases")
	}
	for _, ab23aq23Case := range ab23aq23Cases {
		ab23aq23T.Run(ab23aq23Case.Ab23Aq23ID, func(ab23aq23T *testing.T) {
			ab23aq23Output := ab23aq23EvaluateTest(ab23aq23T, ab23aq23Case.Ab23Aq23Input)
			if ab23aq23Output["state"] != "complete" || ab23aq23Output["physical"] != "Unknown" || ab23aq23Output["qualified"] != false {
				ab23aq23T.Fatal(ab23aq23Output)
			}
			if !reflect.DeepEqual(ab23aq23Output["context"], ab23aq23Case.Ab23Aq23Input["context"]) || !reflect.DeepEqual(ab23aq23Output["arguments"], ab23aq23Case.Ab23Aq23Input["arguments"]) {
				ab23aq23T.Fatal("context or arguments changed")
			}
			if !reflect.DeepEqual(ab23aq23Output["fields"], ab23aq23Case.Ab23Aq23Input["fields"]) {
				ab23aq23T.Fatal("raw changed")
			}
			for ab23aq23Field, ab23aq23Expected := range ab23aq23Case.Ab23Aq23Expected {
				ab23aq23Result := ab23aq23Output["results"].(map[string]any)[ab23aq23Field].(map[string]any)
				if ab23aq23Result["state"] != "complete" || ab23aq23Result["meaning"] != "Hypothesis" {
					ab23aq23T.Fatal(ab23aq23Result)
				}
				for ab23aq23Key, ab23aq23Want := range ab23aq23Expected {
					if !reflect.DeepEqual(ab23aq23Result[ab23aq23Key], ab23aq23Want) {
						ab23aq23T.Fatalf("%s.%s got %v want %v", ab23aq23Field, ab23aq23Key, ab23aq23Result[ab23aq23Key], ab23aq23Want)
					}
				}
				if ab23aq23Result["label"] == "unresolved" {
					if _, ab23aq23Exists := ab23aq23Result["number"]; ab23aq23Exists {
						ab23aq23T.Fatal("inferred actual feedback")
					}
				} else if ab23aq23Result["calculation"] != "Proven" {
					ab23aq23T.Fatal("calculation grade")
				}
			}
			if ab23aq23Case.Ab23Aq23Input["context"].(map[string]any)["operation"] == float64(5) {
				ab23aq23Bytes, _ := json.Marshal(ab23aq23Output)
				if bytes.Contains(ab23aq23Bytes, []byte("900")) || bytes.Contains(ab23aq23Bytes, []byte("timeout")) {
					ab23aq23T.Fatal("inherited timeout")
				}
			}
		})
	}
}

func TestAb23Aq23Context(ab23aq23T *testing.T) {
	for _, ab23aq23Operation := range []int{0, 5} {
		ab23aq23Input := ab23aq23InputTest(ab23aq23Operation, "positive", map[string]any{"setpoint": []int{100}})
		ab23aq23Output := ab23aq23EvaluateTest(ab23aq23T, ab23aq23Input)
		if ab23aq23Output["state"] != "Unknown" || len(ab23aq23Output["results"].(map[string]any)["setpoint"].(map[string]any)) != 2 {
			ab23aq23T.Fatal("undeclared positive view")
		}
	}
	for ab23aq23Key, ab23aq23Value := range map[string]any{"primary": 180, "secondary": 19, "operation": 1, "direction": "response", "view": "other"} {
		ab23aq23T.Run(ab23aq23Key, func(ab23aq23T *testing.T) {
			ab23aq23Input := ab23aq23InputTest(0, "request", map[string]any{"setpoint": []int{100}})
			ab23aq23Input["context"].(map[string]any)[ab23aq23Key] = ab23aq23Value
			ab23aq23Output := ab23aq23EvaluateTest(ab23aq23T, ab23aq23Input)
			if ab23aq23Output["state"] != "Unknown" || len(ab23aq23Output["results"].(map[string]any)["setpoint"].(map[string]any)) != 2 {
				ab23aq23T.Fatal(ab23aq23Output)
			}
		})
	}
	for _, ab23aq23Key := range []string{"primary", "secondary", "operation"} {
		ab23aq23Input := ab23aq23InputTest(0, "request", map[string]any{"setpoint": []int{100}})
		delete(ab23aq23Input["context"].(map[string]any), ab23aq23Key)
		if ab23aq23EvaluateTest(ab23aq23T, ab23aq23Input)["state"] != "Unknown" {
			ab23aq23T.Fatal("absent selector")
		}
	}
	for _, ab23aq23Key := range []string{"family", "variant"} {
		ab23aq23Absent := ab23aq23InputTest(0, "request", map[string]any{"setpoint": []int{100}})
		delete(ab23aq23Absent["context"].(map[string]any), ab23aq23Key)
		ab23aq23AbsentOutput := ab23aq23EvaluateTest(ab23aq23T, ab23aq23Absent)
		if _, ab23aq23Exists := ab23aq23AbsentOutput["context"].(map[string]any)[ab23aq23Key]; ab23aq23Exists || ab23aq23AbsentOutput["qualified"] != false {
			ab23aq23T.Fatal("absent opaque identity changed")
		}
		for _, ab23aq23Value := range []any{nil, false, 0, "", "another-family", []any{}} {
			ab23aq23Input := ab23aq23InputTest(0, "request", map[string]any{"setpoint": []int{100}})
			ab23aq23Input["context"].(map[string]any)[ab23aq23Key] = ab23aq23Value
			ab23aq23Output := ab23aq23EvaluateTest(ab23aq23T, ab23aq23Input)
			if ab23aq23Output["qualified"] != false || ab23aq23Output["state"] != "complete" || !ab23aq23EqualTest(ab23aq23Output["context"].(map[string]any)[ab23aq23Key], ab23aq23Value) {
				ab23aq23T.Fatal(ab23aq23Output)
			}
		}
	}
}

func TestAb23Aq23WidthsAndIndependence(ab23aq23T *testing.T) {
	for _, ab23aq23Op := range []int{0, 2, 3, 4, 5} {
		for _, ab23aq23Direction := range []string{"request", "positive"} {
			ab23aq23Fields := map[string]any{}
			ab23aq23Widths := map[string]int{}
			if ab23aq23Direction == "request" {
				switch ab23aq23Op {
				case 0, 4:
					ab23aq23Widths["setpoint"] = 1
				case 2:
					ab23aq23Widths["valve"] = 1
				case 3:
					ab23aq23Widths["selector"] = 1
				case 5:
					ab23aq23Widths = map[string]int{"channel0": 1, "channel1": 1, "channel2": 1}
				}
			} else {
				switch ab23aq23Op {
				case 2:
					ab23aq23Widths = map[string]int{"valve": 1, "contact": 1, "tank": 2, "virtual_error": 2}
				case 3:
					ab23aq23Widths = map[string]int{"status": 1, "temperature_raw": 2}
				case 4:
					ab23aq23Widths["feedback"] = 1
				}
			}
			for ab23aq23Key, ab23aq23Width := range ab23aq23Widths {
				ab23aq23Fields[ab23aq23Key] = make([]int, ab23aq23Width)
			}
			for ab23aq23Key, ab23aq23Width := range ab23aq23Widths {
				for _, ab23aq23Size := range []int{-1, 0, ab23aq23Width - 1, ab23aq23Width + 1} {
					ab23aq23Input := ab23aq23InputTest(ab23aq23Op, ab23aq23Direction, map[string]any{})
					for ab23aq23Sibling, ab23aq23Data := range ab23aq23Fields {
						ab23aq23Input["fields"].(map[string]any)[ab23aq23Sibling] = ab23aq23Data
					}
					if ab23aq23Size < 0 {
						delete(ab23aq23Input["fields"].(map[string]any), ab23aq23Key)
					} else {
						ab23aq23Input["fields"].(map[string]any)[ab23aq23Key] = make([]int, ab23aq23Size)
					}
					if ab23aq23Op == 5 {
						ab23aq23Input["arguments"] = map[string]any{"external_modes": []int{8, 0, 8}}
					}
					ab23aq23Output := ab23aq23EvaluateTest(ab23aq23T, ab23aq23Input)
					if ab23aq23Output["state"] != "partial" {
						ab23aq23T.Fatal(ab23aq23Output)
					}
					for ab23aq23Sibling := range ab23aq23Widths {
						ab23aq23Result := ab23aq23Output["results"].(map[string]any)[ab23aq23Sibling].(map[string]any)
						if ab23aq23Sibling == ab23aq23Key {
							if ab23aq23Result["state"] != "incomplete" || ab23aq23Result["number"] != nil {
								ab23aq23T.Fatal(ab23aq23Result)
							}
							_, ab23aq23HasRaw := ab23aq23Result["raw"]
							if ab23aq23HasRaw != (ab23aq23Size >= 0) {
								ab23aq23T.Fatal("missing vs empty")
							}
						} else if ab23aq23Result["state"] != "complete" {
							ab23aq23T.Fatal("sibling lost")
						}
					}
				}
			}
		}
	}
	for ab23aq23Index := 0; ab23aq23Index < 3; ab23aq23Index++ {
		ab23aq23Modes := []int{0, 0, 0}
		ab23aq23Modes[ab23aq23Index] = 8
		ab23aq23Input := ab23aq23InputTest(5, "request", map[string]any{"channel0": []int{0}, "channel1": []int{1}, "channel2": []int{255}})
		ab23aq23Input["arguments"] = map[string]any{"external_modes": ab23aq23Modes}
		ab23aq23Output := ab23aq23EvaluateTest(ab23aq23T, ab23aq23Input)
		for ab23aq23Channel, ab23aq23Key := range []string{"channel0", "channel1", "channel2"} {
			if ab23aq23Output["results"].(map[string]any)[ab23aq23Key].(map[string]any)["prerequisite"] != (ab23aq23Channel == ab23aq23Index) {
				ab23aq23T.Fatal("shared mode")
			}
		}
	}
}

func TestAb23Aq23Reject(ab23aq23T *testing.T) {
	for _, ab23aq23JSON := range []string{``, `[]`, `{`, `{} {}`, `{"x":1,"x":2}`, `{"context":{"x":1,"x":2},"fields":{}}`, `{"context":{},"fields":{"x":[{"a":1,"a":2}]}}`, `{"context":{},"fields":{},"extra":[{"x":1,"x":2}]}`, `{"context":{},"fields":{},"extra":NaN}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":null}`, `{"context":{},"fields":{},"arguments":null}`} {
		if ab23aq23Output, ab23aq23Err := Ab23Aq23Evaluate([]byte(ab23aq23JSON)); ab23aq23Err == nil || ab23aq23Output != nil {
			ab23aq23T.Fatalf("accepted %s", ab23aq23JSON)
		}
	}
	for _, ab23aq23Bad := range []any{nil, true, "0", -1, 256, 1.5, []int{}} {
		for _, ab23aq23Key := range []string{"primary", "secondary", "operation", "field", "mode"} {
			ab23aq23Input := ab23aq23InputTest(5, "request", map[string]any{"channel0": []int{0}})
			ab23aq23Input["arguments"] = map[string]any{"external_modes": []any{8, 0, 8}}
			switch ab23aq23Key {
			case "field":
				ab23aq23Input["fields"].(map[string]any)["channel0"] = []any{ab23aq23Bad}
			case "mode":
				ab23aq23Input["arguments"].(map[string]any)["external_modes"] = []any{8, ab23aq23Bad, 8}
			default:
				ab23aq23Input["context"].(map[string]any)[ab23aq23Key] = ab23aq23Bad
			}
			ab23aq23Bytes, _ := json.Marshal(ab23aq23Input)
			if ab23aq23Output, ab23aq23Err := Ab23Aq23Evaluate(ab23aq23Bytes); ab23aq23Err == nil || ab23aq23Output != nil {
				ab23aq23T.Fatalf("accepted %s", ab23aq23Bytes)
			}
		}
	}
	for _, ab23aq23Args := range []any{nil, map[string]any{}, map[string]any{"external_modes": []int{8, 0}}, map[string]any{"external_modes": []int{8, 0, 8, 8}}, map[string]any{"external_modes": []int{8, 0, 8}, "other": 0}} {
		ab23aq23Input := ab23aq23InputTest(5, "request", map[string]any{})
		if ab23aq23Args != nil {
			ab23aq23Input["arguments"] = ab23aq23Args
		}
		ab23aq23Bytes, _ := json.Marshal(ab23aq23Input)
		if ab23aq23Output, ab23aq23Err := Ab23Aq23Evaluate(ab23aq23Bytes); ab23aq23Err == nil || ab23aq23Output != nil {
			ab23aq23T.Fatal("accepted invalid modes")
		}
	}
}

func TestAb23Aq23Memory(ab23aq23T *testing.T) {
	ab23aq23Input := []byte(`{"context":{"primary":181,"secondary":18,"operation":0,"direction":"request","view":"actuator-fields"},"fields":{"setpoint":[0],"opaque":[]},"extra":{"nil":null,"zero":0,"false":false,"empty":[]}}`)
	ab23aq23Original := append([]byte(nil), ab23aq23Input...)
	ab23aq23First, ab23aq23Err := Ab23Aq23Evaluate(ab23aq23Input)
	if ab23aq23Err != nil {
		ab23aq23T.Fatal(ab23aq23Err)
	}
	ab23aq23Saved := append([]byte(nil), ab23aq23First...)
	for ab23aq23Index := range ab23aq23Input {
		ab23aq23Input[ab23aq23Index] = ' '
	}
	if !bytes.Equal(ab23aq23First, ab23aq23Saved) {
		ab23aq23T.Fatal("input alias")
	}
	for ab23aq23Index := range ab23aq23First {
		ab23aq23First[ab23aq23Index] = ' '
	}
	ab23aq23Later, ab23aq23Err := Ab23Aq23Evaluate(ab23aq23Original)
	if ab23aq23Err != nil || !bytes.Equal(ab23aq23Later, ab23aq23Saved) {
		ab23aq23T.Fatal("output alias")
	}
	if !strings.Contains(string(ab23aq23Later), `"opaque":[]`) || !strings.Contains(string(ab23aq23Later), `"nil":null`) || !strings.Contains(string(ab23aq23Later), `"false":false`) {
		ab23aq23T.Fatal("explicit values lost")
	}
}

func ab23aq23EqualTest(ab23aq23Left any, ab23aq23Right any) bool {
	ab23aq23LeftBytes, _ := json.Marshal(ab23aq23Left)
	ab23aq23RightBytes, _ := json.Marshal(ab23aq23Right)
	return bytes.Equal(ab23aq23LeftBytes, ab23aq23RightBytes)
}

func TestAb23Aq23RetentionAndKinds(ab23aq23T *testing.T) {
	ab23aq23Input := ab23aq23InputTest(0, "request", map[string]any{"setpoint": []int{0}, "opaque": []int{}})
	for _, ab23aq23Key := range []string{"state", "results", "qualified", "physical", "input", "unknown"} {
		ab23aq23Input[ab23aq23Key] = map[string]any{"null": nil, "false": false, "zero": 0, "empty": []int{}}
	}
	ab23aq23Output := ab23aq23EvaluateTest(ab23aq23T, ab23aq23Input)
	if !ab23aq23EqualTest(ab23aq23Output["input"], ab23aq23Input) {
		ab23aq23T.Fatal("reserved input keys lost")
	}
	for _, ab23aq23Bad := range []any{nil, false, "0", 0, map[string]any{}} {
		ab23aq23Input := ab23aq23InputTest(0, "request", map[string]any{"setpoint": ab23aq23Bad})
		ab23aq23Bytes, _ := json.Marshal(ab23aq23Input)
		if ab23aq23Result, ab23aq23Err := Ab23Aq23Evaluate(ab23aq23Bytes); ab23aq23Err == nil || ab23aq23Result != nil {
			ab23aq23T.Fatal("invalid field kind")
		}
	}
	for _, ab23aq23Key := range []string{"direction", "view"} {
		for _, ab23aq23Bad := range []any{nil, false, 0, []any{}} {
			ab23aq23Input := ab23aq23InputTest(0, "request", map[string]any{})
			ab23aq23Input["context"].(map[string]any)[ab23aq23Key] = ab23aq23Bad
			ab23aq23Bytes, _ := json.Marshal(ab23aq23Input)
			if ab23aq23Result, ab23aq23Err := Ab23Aq23Evaluate(ab23aq23Bytes); ab23aq23Err == nil || ab23aq23Result != nil {
				ab23aq23T.Fatal("invalid context kind")
			}
		}
	}
	for _, ab23aq23JSON := range []string{`{"context":{},"fields":{},"extra":1e999}`, `{"context":{},"fields":{},"arguments":{"x":1,"x":2}}`, `{"context":{},"fields":{},"fields":{}}`, `{"context":{},"fields":{"x":[],"x":[]}}`} {
		if ab23aq23Result, ab23aq23Err := Ab23Aq23Evaluate([]byte(ab23aq23JSON)); ab23aq23Err == nil || ab23aq23Result != nil {
			ab23aq23T.Fatal("invalid JSON accepted")
		}
	}
}

func TestAb23Aq23RegisteredTemperatureRaw(ab23aq23T *testing.T) {
	ab23aq23Input := ab23aq23InputTest(3, "positive", map[string]any{
		"status": []int{0}, "temperature_raw": []int{0x34, 0x12},
	})
	ab23aq23Output := ab23aq23EvaluateTest(ab23aq23T, ab23aq23Input)
	ab23aq23Results := ab23aq23Output["results"].(map[string]any)
	ab23aq23Result := ab23aq23Results["temperature_raw"].(map[string]any)
	if ab23aq23Output["state"] != "complete" || ab23aq23Result["state"] != "complete" || ab23aq23Result["number"] != float64(4660) || ab23aq23Result["label"] != "raw" || ab23aq23Result["meaning"] != "Hypothesis" || ab23aq23Result["calculation"] != "Proven" {
		ab23aq23T.Fatal(ab23aq23Output)
	}
	if !ab23aq23EqualTest(ab23aq23Result["raw"], []int{0x34, 0x12}) {
		ab23aq23T.Fatal("registered raw bytes changed")
	}
	if _, ab23aq23Exists := ab23aq23Results["value"]; ab23aq23Exists {
		ab23aq23T.Fatal("synthesized unregistered value field")
	}
	ab23aq23Input["fields"].(map[string]any)["value"] = []int{0xff, 0xff}
	ab23aq23Output = ab23aq23EvaluateTest(ab23aq23T, ab23aq23Input)
	ab23aq23Results = ab23aq23Output["results"].(map[string]any)
	if ab23aq23Output["state"] != "complete" || ab23aq23Results["temperature_raw"].(map[string]any)["number"] != float64(4660) || !ab23aq23EqualTest(ab23aq23Results["value"], map[string]any{"raw": []int{255, 255}, "state": "Unknown"}) {
		ab23aq23T.Fatal("extra field changed registered interpretation", ab23aq23Output)
	}
}
