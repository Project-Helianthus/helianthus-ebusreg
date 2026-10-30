package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func ab23ao23Input(ab23ao23Direction string, ab23ao23Bytes []int) map[string]any {
	return map[string]any{
		"context": map[string]any{"primary": 181, "secondary": 32, "operation": 0, "direction": ab23ao23Direction, "view": "logical-object", "family": map[string]any{"opaque": []any{"synthetic", 0}}, "variant": []any{nil, "candidate"}},
		"fields":  map[string]any{"logical_object": ab23ao23Bytes, "unmodelled": []int{0, 255}},
	}
}

func ab23ao23Call(ab23ao23T *testing.T, ab23ao23Input any) map[string]any {
	ab23ao23T.Helper()
	ab23ao23JSON, ab23ao23Err := json.Marshal(ab23ao23Input)
	if ab23ao23Err != nil {
		ab23ao23T.Fatal(ab23ao23Err)
	}
	ab23ao23Result, ab23ao23Err := Ab23Ao23Interpret(ab23ao23JSON)
	if ab23ao23Err != nil {
		ab23ao23T.Fatal(ab23ao23Err)
	}
	ab23ao23Output := map[string]any{}
	if ab23ao23Err := json.Unmarshal(ab23ao23Result, &ab23ao23Output); ab23ao23Err != nil {
		ab23ao23T.Fatal(ab23ao23Err)
	}
	if ab23ao23Output["physical"] != "Unknown" || ab23ao23Output["qualified"] != false {
		ab23ao23T.Fatalf("native qualification inferred: %s", ab23ao23Result)
	}
	for ab23ao23Name, ab23ao23Value := range ab23ao23Output["fields"].(map[string]any) {
		ab23ao23ResultField := ab23ao23Value.(map[string]any)
		ab23ao23Equal(ab23ao23T, ab23ao23ResultField["grade"], "Hypothesis")
		if ab23ao23Name == "relay" || ab23ao23Name == "reserved" || ab23ao23ResultField["state"] == "incomplete" || ab23ao23ResultField["state"] == "unmapped" || ab23ao23ResultField["label"] != nil {
			if _, ab23ao23Exists := ab23ao23ResultField["number"]; ab23ao23Exists {
				ab23ao23T.Fatalf("raw-only, labelled or incomplete field %s gained a number", ab23ao23Name)
			}
		}
	}
	ab23ao23RoundTripJSON, ab23ao23Err := json.Marshal(ab23ao23Output)
	if ab23ao23Err != nil {
		ab23ao23T.Fatal(ab23ao23Err)
	}
	var ab23ao23RoundTrip map[string]any
	if ab23ao23Err := json.Unmarshal(ab23ao23RoundTripJSON, &ab23ao23RoundTrip); ab23ao23Err != nil {
		ab23ao23T.Fatal(ab23ao23Err)
	}
	if !reflect.DeepEqual(ab23ao23Output, ab23ao23RoundTrip) {
		ab23ao23T.Fatal("JSON round trip changed state or presence")
	}
	return ab23ao23Output
}

func ab23ao23Field(ab23ao23Output map[string]any, ab23ao23Name string) map[string]any {
	return ab23ao23Output["fields"].(map[string]any)[ab23ao23Name].(map[string]any)
}

func ab23ao23Equal(ab23ao23T *testing.T, ab23ao23Got, ab23ao23Want any) {
	ab23ao23T.Helper()
	ab23ao23GotJSON, _ := json.Marshal(ab23ao23Got)
	ab23ao23WantJSON, _ := json.Marshal(ab23ao23Want)
	if !bytes.Equal(ab23ao23GotJSON, ab23ao23WantJSON) {
		ab23ao23T.Fatalf("got %s; want %s", ab23ao23GotJSON, ab23ao23WantJSON)
	}
}

func TestAb23Ao23Fixtures(ab23ao23T *testing.T) {
	ab23ao23JSON, ab23ao23Err := os.ReadFile("testdata/ab23ao23_b520_actuator_views.json")
	if ab23ao23Err != nil {
		ab23ao23T.Fatal(ab23ao23Err)
	}
	var ab23ao23Fixtures []map[string]json.RawMessage
	if ab23ao23Err := json.Unmarshal(ab23ao23JSON, &ab23ao23Fixtures); ab23ao23Err != nil {
		ab23ao23T.Fatal(ab23ao23Err)
	}
	if len(ab23ao23Fixtures) != 2 {
		ab23ao23T.Fatal("expected two independent vectors")
	}
	for _, ab23ao23Fixture := range ab23ao23Fixtures {
		var ab23ao23ID, ab23ao23Direction string
		var ab23ao23Bytes []int
		var ab23ao23Numbers map[string]float64
		for ab23ao23Key, ab23ao23Target := range map[string]any{"id": &ab23ao23ID, "direction": &ab23ao23Direction, "bytes": &ab23ao23Bytes, "numbers": &ab23ao23Numbers} {
			if ab23ao23Err := json.Unmarshal(ab23ao23Fixture[ab23ao23Key], ab23ao23Target); ab23ao23Err != nil {
				ab23ao23T.Fatal(ab23ao23Err)
			}
		}
		ab23ao23T.Run(ab23ao23ID, func(ab23ao23T *testing.T) {
			ab23ao23Input := ab23ao23Input(ab23ao23Direction, ab23ao23Bytes)
			ab23ao23Output := ab23ao23Call(ab23ao23T, ab23ao23Input)
			ab23ao23Equal(ab23ao23T, ab23ao23Output["state"], "Hypothesis")
			ab23ao23Equal(ab23ao23T, ab23ao23Output["raw"], ab23ao23Input["fields"])
			ab23ao23Equal(ab23ao23T, ab23ao23Output["context"], ab23ao23Input["context"])
			for ab23ao23Name, ab23ao23Number := range ab23ao23Numbers {
				ab23ao23Field := ab23ao23Field(ab23ao23Output, ab23ao23Name)
				ab23ao23Equal(ab23ao23T, ab23ao23Field["number"], ab23ao23Number)
				ab23ao23Equal(ab23ao23T, ab23ao23Field["grade"], "Hypothesis")
				ab23ao23Equal(ab23ao23T, ab23ao23Field["state"], "Hypothesis")
			}
		})
	}
}

func TestAb23Ao23RequestMappings(ab23ao23T *testing.T) {
	for _, ab23ao23Stepper := range []int{0, 240, 241, 254, 255} {
		ab23ao23T.Run(fmt.Sprintf("ab23ao23_stepper_%d", ab23ao23Stepper), func(ab23ao23T *testing.T) {
			ab23ao23Output := ab23ao23Call(ab23ao23T, ab23ao23Input("request", []int{83, 12, 251, 37, 255, 255, ab23ao23Stepper, 0, 255, 0}))
			ab23ao23Field := ab23ao23Field(ab23ao23Output, "stepper")
			if ab23ao23Stepper <= 240 {
				ab23ao23Equal(ab23ao23T, ab23ao23Field["number"], ab23ao23Stepper)
			} else {
				if _, ab23ao23Exists := ab23ao23Field["number"]; ab23ao23Exists {
					ab23ao23T.Fatal("unmapped/calibration carrier gained a number")
				}
				if ab23ao23Stepper == 255 {
					ab23ao23Equal(ab23ao23T, ab23ao23Field["label"], "calibrate/go-to-zero")
				} else {
					ab23ao23Equal(ab23ao23T, ab23ao23Field["state"], "unmapped")
					ab23ao23Equal(ab23ao23T, ab23ao23Output["state"], "partial")
				}
			}
			for _, ab23ao23Name := range []string{"PWM1", "PWM2", "ED"} {
				ab23ao23Equal(ab23ao23T, ab23ao23FieldFrom(ab23ao23Output, ab23ao23Name, "number"), 127.5)
				ab23ao23Equal(ab23ao23T, ab23ao23FieldFrom(ab23ao23Output, ab23ao23Name, "calculation"), "Proven")
			}
		})
	}
	for ab23ao23Index, ab23ao23Labels := range map[int][]string{7: {"disabled", "enabled"}, 9: {"standard3K3", "extended1K", "low-temperature15K", "standard1K"}} {
		for ab23ao23Value := 0; ab23ao23Value <= len(ab23ao23Labels); ab23ao23Value++ {
			ab23ao23T.Run(fmt.Sprintf("ab23ao23_map_%d_%d", ab23ao23Index, ab23ao23Value), func(ab23ao23T *testing.T) {
				ab23ao23Bytes := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
				ab23ao23Bytes[ab23ao23Index] = ab23ao23Value
				ab23ao23Name := "LED"
				if ab23ao23Index == 9 {
					ab23ao23Name = "temperature-table"
				}
				ab23ao23Output := ab23ao23Call(ab23ao23T, ab23ao23Input("request", ab23ao23Bytes))
				ab23ao23Field := ab23ao23Field(ab23ao23Output, ab23ao23Name)
				if ab23ao23Value == len(ab23ao23Labels) {
					ab23ao23Equal(ab23ao23T, ab23ao23Field["state"], "unmapped")
					if _, ab23ao23Exists := ab23ao23Field["number"]; ab23ao23Exists {
						ab23ao23T.Fatal("unmapped number")
					}
				} else {
					ab23ao23Equal(ab23ao23T, ab23ao23Field["label"], ab23ao23Labels[ab23ao23Value])
					if ab23ao23Index == 9 {
						ab23ao23Equal(ab23ao23T, ab23ao23Field["flags"].(map[string]any)["candidate_celsius_range"], [][2]float64{{-13.5, 116}, {-20, 155}, {-39, 80}, {-13.9, 130}}[ab23ao23Value])
					}
				}
			})
		}
	}
}

func ab23ao23FieldFrom(ab23ao23Output map[string]any, ab23ao23Name, ab23ao23Key string) any {
	return ab23ao23Field(ab23ao23Output, ab23ao23Name)[ab23ao23Key]
}

func TestAb23Ao23Domains(ab23ao23T *testing.T) {
	for ab23ao23Index, ab23ao23Name := range map[int]string{3: "ADC", 5: "temperature1", 7: "temperature2", 9: "temperature3", 11: "temperature4", 14: "impeller", 16: "vortex"} {
		ab23ao23Values := []int{0, 1023, 1024}
		ab23ao23Limit := 1023
		if ab23ao23Index >= 5 && ab23ao23Index <= 11 {
			ab23ao23Values = []int{0, 16, 3072, 3200, 3201, 65535}
			ab23ao23Limit = 3200
		} else if ab23ao23Index == 14 {
			ab23ao23Values = []int{0, 4, 1000, 1001}
			ab23ao23Limit = 1000
		} else if ab23ao23Index == 16 {
			ab23ao23Values = []int{0, 4, 5, 350, 351}
			ab23ao23Limit = 350
		}
		for _, ab23ao23Value := range ab23ao23Values {
			ab23ao23T.Run(fmt.Sprintf("ab23ao23_%s_%d", ab23ao23Name, ab23ao23Value), func(ab23ao23T *testing.T) {
				ab23ao23Bytes := make([]int, 18)
				ab23ao23Bytes[ab23ao23Index], ab23ao23Bytes[ab23ao23Index+1] = ab23ao23Value%256, ab23ao23Value/256
				ab23ao23Output := ab23ao23Call(ab23ao23T, ab23ao23Input("positive", ab23ao23Bytes))
				ab23ao23Field := ab23ao23Field(ab23ao23Output, ab23ao23Name)
				ab23ao23Want := float64(ab23ao23Value)
				if ab23ao23Index >= 5 && ab23ao23Index <= 11 {
					ab23ao23Want = map[int]float64{0: 0, 16: 1, 3072: 192, 3200: 200, 3201: 200.0625, 65535: -0.0625}[ab23ao23Value]
					ab23ao23Equal(ab23ao23T, ab23ao23Field["calculation"], "Proven")
				}
				ab23ao23Equal(ab23ao23T, ab23ao23Field["number"], ab23ao23Want)
				ab23ao23Equal(ab23ao23T, ab23ao23Field["candidate_domain"], ab23ao23Value <= ab23ao23Limit)
				ab23ao23Equal(ab23ao23T, ab23ao23Output["state"], "Hypothesis")
				if ab23ao23Index == 16 {
					ab23ao23Equal(ab23ao23T, ab23ao23Field["flags"], map[string]any{"candidate_operational_domain": ab23ao23Value >= 5 && ab23ao23Value <= 350, "candidate_hysteresis_hz": map[string]int{"on": 5, "off": 4}})
				}
			})
		}
	}
}

func TestAb23Ao23IndependentBytes(ab23ao23T *testing.T) {
	for _, ab23ao23Direction := range []string{"request", "positive"} {
		ab23ao23Names := []string{"relay", "PWM1", "PWM2", "stepper", "LED", "ED", "temperature-table"}
		ab23ao23Offsets, ab23ao23Widths, ab23ao23Length := []int{3, 4, 5, 6, 7, 8, 9}, []int{1, 1, 1, 1, 1, 1, 1}, 10
		if ab23ao23Direction == "positive" {
			ab23ao23Names = []string{"ADC", "temperature1", "temperature2", "temperature3", "temperature4", "reserved", "impeller", "vortex"}
			ab23ao23Offsets, ab23ao23Widths, ab23ao23Length = []int{3, 5, 7, 9, 11, 13, 14, 16}, []int{2, 2, 2, 2, 2, 1, 2, 2}, 18
		}
		ab23ao23Control := ab23ao23Call(ab23ao23T, ab23ao23Input(ab23ao23Direction, make([]int, ab23ao23Length)))
		for ab23ao23Position := 0; ab23ao23Position < ab23ao23Length+2; ab23ao23Position++ {
			ab23ao23T.Run(fmt.Sprintf("ab23ao23_%s_byte_%d", ab23ao23Direction, ab23ao23Position), func(ab23ao23T *testing.T) {
				ab23ao23Bytes := make([]int, ab23ao23Length+2)
				ab23ao23Bytes[ab23ao23Position] = 1
				ab23ao23Input := ab23ao23Input(ab23ao23Direction, ab23ao23Bytes)
				ab23ao23Output := ab23ao23Call(ab23ao23T, ab23ao23Input)
				ab23ao23Equal(ab23ao23T, ab23ao23Output["raw"], ab23ao23Input["fields"])
				for ab23ao23Index, ab23ao23Name := range ab23ao23Names {
					ab23ao23Field := ab23ao23Field(ab23ao23Output, ab23ao23Name)
					if ab23ao23Position < ab23ao23Offsets[ab23ao23Index] || ab23ao23Position >= ab23ao23Offsets[ab23ao23Index]+ab23ao23Widths[ab23ao23Index] {
						ab23ao23Equal(ab23ao23T, ab23ao23Field, ab23ao23FieldFromControl(ab23ao23Control, ab23ao23Name))
					} else {
						ab23ao23Equal(ab23ao23T, ab23ao23Field["raw"], ab23ao23Bytes[ab23ao23Offsets[ab23ao23Index]:ab23ao23Offsets[ab23ao23Index]+ab23ao23Widths[ab23ao23Index]])
						if ab23ao23Name == "reserved" {
							ab23ao23Equal(ab23ao23T, ab23ao23Field["reserved"], true)
						} else if ab23ao23Direction == "positive" {
							ab23ao23Number := float64(1)
							if ab23ao23Position != ab23ao23Offsets[ab23ao23Index] {
								ab23ao23Number = 256
							}
							if ab23ao23Index >= 1 && ab23ao23Index <= 4 {
								ab23ao23Number /= 16
							}
							ab23ao23Equal(ab23ao23T, ab23ao23Field["number"], ab23ao23Number)
						} else if ab23ao23Name == "PWM1" || ab23ao23Name == "PWM2" || ab23ao23Name == "ED" {
							ab23ao23Equal(ab23ao23T, ab23ao23Field["number"], 0.5)
						}
					}
				}
			})
		}
	}
}

func ab23ao23FieldFromControl(ab23ao23Output map[string]any, ab23ao23Name string) map[string]any {
	return ab23ao23Field(ab23ao23Output, ab23ao23Name)
}

func TestAb23Ao23Truncation(ab23ao23T *testing.T) {
	for _, ab23ao23Direction := range []string{"request", "positive"} {
		ab23ao23Names, ab23ao23Ends := []string{"relay", "PWM1", "PWM2", "stepper", "LED", "ED", "temperature-table"}, []int{4, 5, 6, 7, 8, 9, 10}
		if ab23ao23Direction == "positive" {
			ab23ao23Names, ab23ao23Ends = []string{"ADC", "temperature1", "temperature2", "temperature3", "temperature4", "reserved", "impeller", "vortex"}, []int{5, 7, 9, 11, 13, 14, 16, 18}
		}
		ab23ao23Full := make([]int, ab23ao23Ends[len(ab23ao23Ends)-1])
		ab23ao23Control := ab23ao23Call(ab23ao23T, ab23ao23Input(ab23ao23Direction, ab23ao23Full))
		for ab23ao23Length := 0; ab23ao23Length <= len(ab23ao23Full); ab23ao23Length++ {
			ab23ao23T.Run(fmt.Sprintf("ab23ao23_%s_length_%d", ab23ao23Direction, ab23ao23Length), func(ab23ao23T *testing.T) {
				ab23ao23Output := ab23ao23Call(ab23ao23T, ab23ao23Input(ab23ao23Direction, ab23ao23Full[:ab23ao23Length]))
				ab23ao23WantState := "partial"
				if ab23ao23Length == len(ab23ao23Full) {
					ab23ao23WantState = "Hypothesis"
				}
				ab23ao23Equal(ab23ao23T, ab23ao23Output["state"], ab23ao23WantState)
				for ab23ao23Index, ab23ao23Name := range ab23ao23Names {
					ab23ao23Field := ab23ao23Field(ab23ao23Output, ab23ao23Name)
					if ab23ao23Length >= ab23ao23Ends[ab23ao23Index] {
						ab23ao23Equal(ab23ao23T, ab23ao23Field, ab23ao23FieldFromControl(ab23ao23Control, ab23ao23Name))
					} else {
						ab23ao23Equal(ab23ao23T, ab23ao23Field["state"], "incomplete")
						if _, ab23ao23Exists := ab23ao23Field["number"]; ab23ao23Exists {
							ab23ao23T.Fatal("default filled incomplete field")
						}
					}
				}
			})
		}
	}
}

func TestAb23Ao23Selectors(ab23ao23T *testing.T) {
	for _, ab23ao23Direction := range []string{"request", "positive"} {
		for ab23ao23Key, ab23ao23Values := range map[string][]any{"primary": {180, nil}, "secondary": {33, nil}, "operation": {1, nil}, "direction": {"response", "negative", "unknown", nil}, "view": {"unknown", "", nil}} {
			for ab23ao23Index, ab23ao23Value := range ab23ao23Values {
				ab23ao23T.Run(fmt.Sprintf("ab23ao23_%s_%s_%d", ab23ao23Direction, ab23ao23Key, ab23ao23Index), func(ab23ao23T *testing.T) {
					ab23ao23Input := ab23ao23Input(ab23ao23Direction, make([]int, 20))
					ab23ao23Control := ab23ao23Call(ab23ao23T, ab23ao23Input)
					ab23ao23Equal(ab23ao23T, ab23ao23Control["state"], "Hypothesis")
					ab23ao23Context := ab23ao23Input["context"].(map[string]any)
					if ab23ao23Value == nil {
						delete(ab23ao23Context, ab23ao23Key)
					} else {
						ab23ao23Context[ab23ao23Key] = ab23ao23Value
					}
					ab23ao23Output := ab23ao23Call(ab23ao23T, ab23ao23Input)
					ab23ao23Equal(ab23ao23T, ab23ao23Output["state"], "Unknown")
					ab23ao23Equal(ab23ao23T, ab23ao23Output["fields"], map[string]any{})
					ab23ao23Equal(ab23ao23T, ab23ao23Output["raw"], ab23ao23Input["fields"])
					ab23ao23Equal(ab23ao23T, ab23ao23Output["context"], ab23ao23Context)
				})
			}
		}
	}
	for _, ab23ao23Key := range []string{"family", "variant", "extra"} {
		for _, ab23ao23Value := range []any{nil, 0, "unknown", map[string]any{"opaque": []any{true, 3}}} {
			ab23ao23Input := ab23ao23Input("positive", make([]int, 18))
			ab23ao23Input["context"].(map[string]any)[ab23ao23Key] = ab23ao23Value
			ab23ao23Output := ab23ao23Call(ab23ao23T, ab23ao23Input)
			ab23ao23Equal(ab23ao23T, ab23ao23Output["context"], ab23ao23Input["context"])
			ab23ao23Equal(ab23ao23T, ab23ao23Output["state"], "Hypothesis")
		}
	}
}

func TestAb23Ao23Presence(ab23ao23T *testing.T) {
	for _, ab23ao23Direction := range []string{"request", "positive"} {
		for _, ab23ao23Kind := range []string{"missing", "empty", "zero"} {
			ab23ao23T.Run("ab23ao23_"+ab23ao23Direction+"_"+ab23ao23Kind, func(ab23ao23T *testing.T) {
				ab23ao23Input := ab23ao23Input(ab23ao23Direction, []int{})
				ab23ao23Fields := ab23ao23Input["fields"].(map[string]any)
				switch ab23ao23Kind {
				case "missing":
					delete(ab23ao23Fields, "logical_object")
				case "zero":
					ab23ao23Fields["logical_object"] = []int{0, 0, 0, 0, 0}
				}
				ab23ao23Input["arguments"] = map[string]any{}
				ab23ao23Output := ab23ao23Call(ab23ao23T, ab23ao23Input)
				ab23ao23Equal(ab23ao23T, ab23ao23Output["raw"], ab23ao23Fields)
				ab23ao23Equal(ab23ao23T, ab23ao23Output["arguments"], map[string]any{})
				ab23ao23Name := "relay"
				if ab23ao23Direction == "positive" {
					ab23ao23Name = "ADC"
				}
				ab23ao23Field := ab23ao23Field(ab23ao23Output, ab23ao23Name)
				if ab23ao23Kind == "missing" {
					ab23ao23Equal(ab23ao23T, ab23ao23Field["raw"], nil)
				} else if ab23ao23Kind == "empty" {
					ab23ao23Equal(ab23ao23T, ab23ao23Field["raw"], []int{})
				} else if ab23ao23Direction == "positive" {
					ab23ao23Equal(ab23ao23T, ab23ao23Field["number"], 0)
				}
				delete(ab23ao23Input, "arguments")
				if _, ab23ao23Exists := ab23ao23Call(ab23ao23T, ab23ao23Input)["arguments"]; ab23ao23Exists {
					ab23ao23T.Fatal("absent arguments gained presence")
				}
			})
		}
	}
}

func TestAb23Ao23OpaqueFields(ab23ao23T *testing.T) {
	for _, ab23ao23Direction := range []string{"request", "positive"} {
		ab23ao23Control := ab23ao23Call(ab23ao23T, ab23ao23Input(ab23ao23Direction, make([]int, 18)))
		for ab23ao23Index, ab23ao23Value := range []any{nil, []int{}, []int{0}, []int{255}, []int{1, 2, 3}} {
			ab23ao23T.Run(fmt.Sprintf("ab23ao23_%s_opaque_%d", ab23ao23Direction, ab23ao23Index), func(ab23ao23T *testing.T) {
				ab23ao23Input := ab23ao23Input(ab23ao23Direction, make([]int, 18))
				ab23ao23Fields := ab23ao23Input["fields"].(map[string]any)
				if ab23ao23Value == nil {
					delete(ab23ao23Fields, "unmodelled")
				} else {
					ab23ao23Fields["unmodelled"] = ab23ao23Value
				}
				ab23ao23Output := ab23ao23Call(ab23ao23T, ab23ao23Input)
				ab23ao23Equal(ab23ao23T, ab23ao23Output["raw"], ab23ao23Fields)
				ab23ao23Equal(ab23ao23T, ab23ao23Output["fields"], ab23ao23Control["fields"])
			})
		}
	}
	for _, ab23ao23Key := range []string{"family", "variant"} {
		ab23ao23Input := ab23ao23Input("request", make([]int, 10))
		delete(ab23ao23Input["context"].(map[string]any), ab23ao23Key)
		ab23ao23Output := ab23ao23Call(ab23ao23T, ab23ao23Input)
		ab23ao23Equal(ab23ao23T, ab23ao23Output["context"], ab23ao23Input["context"])
		ab23ao23Equal(ab23ao23T, ab23ao23Output["state"], "Hypothesis")
	}
	ab23ao23Output := ab23ao23Call(ab23ao23T, map[string]any{"context": map[string]any{"primary": 181, "secondary": 32, "operation": json.Number("18446744073709551616"), "direction": "request", "view": "logical-object"}, "fields": map[string]any{}})
	ab23ao23Equal(ab23ao23T, ab23ao23Output["state"], "Unknown")
}

func TestAb23Ao23InvalidJSON(ab23ao23T *testing.T) {
	for ab23ao23Index, ab23ao23JSON := range []string{
		``, `{`, `null`, `[]`, `{} {}`, `{"context":{},"fields":{},"context":{}}`,
		`{"context":{"family":{"x":1,"x":2}},"fields":{}}`,
		`{"context":{},"fields":{"x":[],"x":[]}}`,
		`{"context":{},"fields":{},"arguments":{"x":0,"x":1}}`,
		`{"context":null,"fields":{}}`, `{"context":{},"fields":null}`,
		`{"context":{},"fields":{"x":null}}`, `{"context":{},"fields":{"x":"AA=="}}`,
		`{"context":{},"fields":{"x":[-1]}}`, `{"context":{},"fields":{"x":[256]}}`,
		`{"context":{},"fields":{"x":[1.5]}}`, `{"context":{},"fields":{"x":[1.0]}}`,
		`{"context":{},"fields":{"x":[true]}}`, `{"context":{},"fields":{"x":[null]}}`,
		`{"context":{},"fields":{"x":["1"]}}`, `{"context":{},"fields":{"x":[{}]}}`,
		`{"context":{},"fields":{"x":[[]]}}`, `{"context":{},"fields":{"x":[1e0]}}`,
		`{"context":{},"fields":{},"arguments":null}`, `{"context":{},"fields":{},"arguments":[]}`,
		`{"context":{},"fields":{},"arguments":{"positive":true}}`,
		`{"context":{"operation":"0"},"fields":{}}`, `{"context":{"primary":181.1},"fields":{}}`,
		`{"context":{"secondary":null},"fields":{}}`, `{"context":{"direction":0},"fields":{}}`,
		`{"context":{"view":false},"fields":{}}`, `{"context":{},"fields":{},"extra":0}`,
		"{\"context\":{\"family\":\"" + string([]byte{255}) + "\"},\"fields\":{}}",
	} {
		ab23ao23T.Run(fmt.Sprintf("ab23ao23_invalid_%d", ab23ao23Index), func(ab23ao23T *testing.T) {
			ab23ao23Result, ab23ao23Err := Ab23Ao23Interpret(json.RawMessage(ab23ao23JSON))
			if ab23ao23Err == nil || ab23ao23Result != nil {
				ab23ao23T.Fatalf("invalid input succeeded: %s => %s, %v", ab23ao23JSON, ab23ao23Result, ab23ao23Err)
			}
		})
	}
}

func TestAb23Ao23UnicodeEscapes(t *testing.T) {
	for _, input := range []string{
		`{"context":{"family":"\ud800"},"fields":{}}`,
		`{"context":{"family":"\udfff"},"fields":{}}`,
		`{"context":{"family":"\ud800\u0041"},"fields":{}}`,
		`{"context":{"family":"\ud800\ud800"},"fields":{}}`,
		`{"context":{"family":"\udc00\ud800"},"fields":{}}`,
		`{"context":{"family":"\ud800x\udc00"},"fields":{}}`,
		`{"context":{"family":"\ud800\\udc00"},"fields":{}}`,
		`{"context":{"family":{"\ud800":[]}},"fields":{}}`,
		`{"context":{"family":["\udc00"]},"fields":{}}`,
		`{"context":{},"fields":{"\ud800":[]}}`,
	} {
		t.Run(input, func(t *testing.T) {
			output, err := Ab23Ao23Interpret(json.RawMessage(input))
			if err == nil || output != nil {
				t.Fatalf("unpaired surrogate accepted: %s, %v", output, err)
			}
		})
	}
	for _, tc := range []struct{ encoded, want string }{
		{`"\ud800\udc00"`, "\U00010000"},
		{`"\uDBFF\uDFFF"`, "\U0010ffff"},
		{`"\ud83d\ude00"`, "\U0001f600"},
		{`"\\ud800"`, `\ud800`},
		{`"\\\ud800\udc00"`, "\\\U00010000"},
		{`"\ufffd"`, "\ufffd"},
		{`"�"`, "\ufffd"},
		{`"\"\u0041\\"`, "\"A\\"},
	} {
		t.Run(tc.encoded, func(t *testing.T) {
			input := json.RawMessage(`{"context":{"family":` + tc.encoded + `},"fields":{}}`)
			output, err := Ab23Ao23Interpret(input)
			if err != nil {
				t.Fatal(err)
			}
			var decoded struct{ Context struct{ Family string } }
			if err := json.Unmarshal(output, &decoded); err != nil {
				t.Fatal(err)
			}
			if decoded.Context.Family != tc.want {
				t.Fatalf("got %q; want %q", decoded.Context.Family, tc.want)
			}
		})
	}
}

func TestAb23Ao23Isolation(ab23ao23T *testing.T) {
	ab23ao23Input := ab23ao23Input("positive", []int{0, 0, 0, 255, 3, 0, 12, 255, 255, 0, 0, 16, 0, 171, 232, 3, 94, 1})
	ab23ao23Input["arguments"] = map[string]any{}
	ab23ao23JSON, _ := json.Marshal(ab23ao23Input)
	ab23ao23Original := append(json.RawMessage(nil), ab23ao23JSON...)
	ab23ao23First, ab23ao23Err := Ab23Ao23Interpret(ab23ao23JSON)
	if ab23ao23Err != nil {
		ab23ao23T.Fatal(ab23ao23Err)
	}
	ab23ao23Second, ab23ao23Err := Ab23Ao23Interpret(ab23ao23JSON)
	if ab23ao23Err != nil {
		ab23ao23T.Fatal(ab23ao23Err)
	}
	ab23ao23Snapshot := append(json.RawMessage(nil), ab23ao23First...)
	for ab23ao23Index := range ab23ao23JSON {
		ab23ao23JSON[ab23ao23Index] = 'x'
	}
	if !bytes.Equal(ab23ao23First, ab23ao23Snapshot) {
		ab23ao23T.Fatal("input mutation affected output")
	}
	var ab23ao23Nested map[string]any
	if ab23ao23Err := json.Unmarshal(ab23ao23First, &ab23ao23Nested); ab23ao23Err != nil {
		ab23ao23T.Fatal(ab23ao23Err)
	}
	ab23ao23Nested["raw"].(map[string]any)["logical_object"].([]any)[5] = 99
	ab23ao23Equal(ab23ao23T, ab23ao23FieldFrom(ab23ao23Nested, "temperature1", "raw"), []int{0, 12})
	ab23ao23Field(ab23ao23Nested, "temperature1")["raw"].([]any)[0] = 88
	ab23ao23Equal(ab23ao23T, ab23ao23FieldFrom(ab23ao23Nested, "temperature2", "raw"), []int{255, 255})
	ab23ao23Nested["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "changed"
	ab23ao23Nested["arguments"].(map[string]any)["changed"] = true
	for ab23ao23Index := range ab23ao23First {
		ab23ao23First[ab23ao23Index] = 'y'
	}
	if !bytes.Equal(ab23ao23Second, ab23ao23Snapshot) {
		ab23ao23T.Fatal("output mutation affected another result")
	}
	ab23ao23Third, ab23ao23Err := Ab23Ao23Interpret(ab23ao23Original)
	if ab23ao23Err != nil || !reflect.DeepEqual(ab23ao23Third, ab23ao23Second) {
		ab23ao23T.Fatal("mutation affected later call", ab23ao23Err)
	}
}
