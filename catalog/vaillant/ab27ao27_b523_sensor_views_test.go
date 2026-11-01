package vaillant

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func ab27ao27Input(ab27ao27Operation int, ab27ao27Bytes []int) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 35, "operation": ab27ao27Operation, "direction": "positive", "view": "status-zero-valid", "family": map[string]any{"opaque": []any{"synthetic-family", 0}}, "variant": []any{"synthetic-variant", nil}}, "fields": map[string]any{"logical_object": ab27ao27Bytes, "unmodelled": []int{0, 255, 7}}}
}

func ab27ao27Call(ab27ao27T *testing.T, ab27ao27InputValue any) map[string]any {
	ab27ao27T.Helper()
	ab27ao27InputJSON, ab27ao27Err := json.Marshal(ab27ao27InputValue)
	if ab27ao27Err != nil {
		ab27ao27T.Fatal(ab27ao27Err)
	}
	ab27ao27Result, ab27ao27Err := Ab27Ao27Interpret(ab27ao27InputJSON)
	if ab27ao27Err != nil {
		ab27ao27T.Fatal(ab27ao27Err)
	}
	var ab27ao27Output map[string]any
	if ab27ao27Err := json.Unmarshal(ab27ao27Result, &ab27ao27Output); ab27ao27Err != nil {
		ab27ao27T.Fatal(ab27ao27Err)
	}
	if ab27ao27Output["physical"] != "Unknown" || ab27ao27Output["qualified"] != false {
		ab27ao27T.Fatal("qualification inferred")
	}
	var ab27ao27RoundTrip map[string]any
	if ab27ao27Err := json.Unmarshal(ab27ao27InputJSON, &ab27ao27RoundTrip); ab27ao27Err != nil {
		ab27ao27T.Fatal(ab27ao27Err)
	}
	if !reflect.DeepEqual(ab27ao27Output["context"], ab27ao27RoundTrip["context"]) || !reflect.DeepEqual(ab27ao27Output["raw"], ab27ao27RoundTrip["fields"]) {
		ab27ao27T.Fatal("native evidence lost")
	}
	if ab27ao27Expected, ab27ao27Present := ab27ao27RoundTrip["arguments"]; ab27ao27Present {
		if !reflect.DeepEqual(ab27ao27Output["arguments"], ab27ao27Expected) {
			ab27ao27T.Fatal("arguments lost")
		}
	} else if _, ab27ao27Present := ab27ao27Output["arguments"]; ab27ao27Present {
		ab27ao27T.Fatal("arguments defaulted")
	}
	return ab27ao27Output
}

func ab27ao27Field(ab27ao27Output map[string]any, ab27ao27Name string) map[string]any {
	return ab27ao27Output["fields"].(map[string]any)[ab27ao27Name].(map[string]any)
}

func TestAb27Ao27Fixture(ab27ao27T *testing.T) {
	ab27ao27Data, ab27ao27Err := os.ReadFile("testdata/ab27ao27_b523_sensor_views.json")
	if ab27ao27Err != nil {
		ab27ao27T.Fatal(ab27ao27Err)
	}
	var ab27ao27Cases []struct {
		Ab27Ao27ID       string                    `json:"id"`
		Ab27Ao27Input    map[string]any            `json:"input"`
		Ab27Ao27Expected map[string]map[string]any `json:"expected"`
	}
	if ab27ao27Err := json.Unmarshal(ab27ao27Data, &ab27ao27Cases); ab27ao27Err != nil {
		ab27ao27T.Fatal(ab27ao27Err)
	}
	if len(ab27ao27Cases) != 3 {
		ab27ao27T.Fatal("fixture count")
	}
	for _, ab27ao27Case := range ab27ao27Cases {
		ab27ao27T.Run(ab27ao27Case.Ab27Ao27ID, func(ab27ao27T *testing.T) {
			ab27ao27Output := ab27ao27Call(ab27ao27T, ab27ao27Case.Ab27Ao27Input)
			for ab27ao27Name, ab27ao27Expected := range ab27ao27Case.Ab27Ao27Expected {
				for ab27ao27Key, ab27ao27Value := range ab27ao27Expected {
					if !reflect.DeepEqual(ab27ao27Field(ab27ao27Output, ab27ao27Name)[ab27ao27Key], ab27ao27Value) {
						ab27ao27T.Fatalf("%s %s", ab27ao27Name, ab27ao27Key)
					}
				}
			}
		})
	}
}

func TestAb27Ao27Sensors(ab27ao27T *testing.T) {
	for _, ab27ao27Operation := range []int{3, 6, 7} {
		ab27ao27Positions := []int{3, 5, 7, 9, 11, 13}
		ab27ao27Start := 1
		ab27ao27Extent := 18
		if ab27ao27Operation == 6 {
			ab27ao27Positions = []int{3, 5, 7, 9, 11, 13, 15}
			ab27ao27Extent = 19
		}
		if ab27ao27Operation == 7 {
			ab27ao27Positions = []int{3, 5, 7, 9, 12, 14}
			ab27ao27Start = 8
		}
		for ab27ao27Index, ab27ao27Position := range ab27ao27Positions {
			for _, ab27ao27Vector := range []struct {
				ab27ao27Low, ab27ao27High int
				ab27ao27Number            float64
				ab27ao27Replacement       bool
			}{{0, 128, 0, true}, {16, 0, 1, false}, {240, 255, -1, false}, {255, 127, 2047.9375, false}, {1, 128, -2047.9375, false}} {
				ab27ao27T.Run(fmt.Sprintf("op%d_S%d_%d_%d", ab27ao27Operation, ab27ao27Start+ab27ao27Index, ab27ao27Vector.ab27ao27Low, ab27ao27Vector.ab27ao27High), func(ab27ao27T *testing.T) {
					ab27ao27Bytes := make([]int, ab27ao27Extent+2)
					ab27ao27Bytes[0], ab27ao27Bytes[1], ab27ao27Bytes[2] = 255, 99, 0
					ab27ao27Bytes[ab27ao27Position], ab27ao27Bytes[ab27ao27Position+1] = ab27ao27Vector.ab27ao27Low, ab27ao27Vector.ab27ao27High
					if ab27ao27Operation == 7 {
						ab27ao27Bytes[11] = 77
					}
					ab27ao27Bytes[ab27ao27Extent], ab27ao27Bytes[ab27ao27Extent+1] = 222, 111
					ab27ao27Output := ab27ao27Call(ab27ao27T, ab27ao27Input(ab27ao27Operation, ab27ao27Bytes))
					if ab27ao27Output["state"] != "Hypothesis" {
						ab27ao27T.Fatal("complete state")
					}
					for ab27ao27Sibling := range ab27ao27Positions {
						ab27ao27FieldValue := ab27ao27Field(ab27ao27Output, fmt.Sprintf("S%d", ab27ao27Start+ab27ao27Sibling))
						if ab27ao27FieldValue["state"] != "Hypothesis" || ab27ao27FieldValue["grade"] != "Hypothesis" {
							ab27ao27T.Fatal("grade")
						}
						if ab27ao27Sibling == ab27ao27Index && ab27ao27Vector.ab27ao27Replacement {
							if ab27ao27FieldValue["raw"] != float64(32768) || ab27ao27FieldValue["label"] != "replacement" {
								ab27ao27T.Fatal("replacement selector")
							}
							if _, ab27ao27Present := ab27ao27FieldValue["number"]; ab27ao27Present {
								ab27ao27T.Fatal("replacement converted")
							}
							if _, ab27ao27Present := ab27ao27FieldValue["calculation"]; ab27ao27Present {
								ab27ao27T.Fatal("replacement calculation")
							}
						} else {
							ab27ao27Want := float64(0)
							if ab27ao27Sibling == ab27ao27Index {
								ab27ao27Want = ab27ao27Vector.ab27ao27Number
							}
							if ab27ao27FieldValue["number"] != ab27ao27Want || ab27ao27FieldValue["calculation"] != "Proven" {
								ab27ao27T.Fatalf("sensor %d: %#v want %v", ab27ao27Sibling, ab27ao27FieldValue, ab27ao27Want)
							}
						}
					}
					if ab27ao27Operation == 7 {
						ab27ao27Opaque := ab27ao27Field(ab27ao27Output, "opaque11")
						if ab27ao27Opaque["raw"] != float64(77) || ab27ao27Opaque["state"] != "Unknown" {
							ab27ao27T.Fatal("opaque position")
						}
						if _, ab27ao27Present := ab27ao27Opaque["number"]; ab27ao27Present {
							ab27ao27T.Fatal("opaque numeric guess")
						}
					}
				})
			}
		}
	}
}

func TestAb27Ao27StatusPairs(ab27ao27T *testing.T) {
	for _, ab27ao27Operation := range []int{3, 6, 7} {
		ab27ao27Count, ab27ao27Start, ab27ao27Status, ab27ao27Extent := 6, 1, 16, 18
		if ab27ao27Operation == 6 {
			ab27ao27Count, ab27ao27Status, ab27ao27Extent = 7, 17, 19
		}
		if ab27ao27Operation == 7 {
			ab27ao27Start = 8
		}
		for _, ab27ao27View := range []string{"status-zero-valid", "status-zero-valid-or-low"} {
			for ab27ao27Index := 0; ab27ao27Index < ab27ao27Count; ab27ao27Index++ {
				for ab27ao27Pair := 0; ab27ao27Pair < 4; ab27ao27Pair++ {
					ab27ao27T.Run(fmt.Sprintf("op%d_%s_pair%d_value%d", ab27ao27Operation, ab27ao27View, ab27ao27Index, ab27ao27Pair), func(ab27ao27T *testing.T) {
						ab27ao27Bytes := make([]int, ab27ao27Extent)
						ab27ao27Bytes[ab27ao27Status+ab27ao27Index/4] = ab27ao27Pair << (2 * (ab27ao27Index % 4))
						ab27ao27InputValue := ab27ao27Input(ab27ao27Operation, ab27ao27Bytes)
						ab27ao27InputValue["context"].(map[string]any)["view"] = ab27ao27View
						ab27ao27Output := ab27ao27Call(ab27ao27T, ab27ao27InputValue)
						for ab27ao27Sibling := 0; ab27ao27Sibling < ab27ao27Count; ab27ao27Sibling++ {
							ab27ao27ExpectedPair := 0
							if ab27ao27Sibling == ab27ao27Index {
								ab27ao27ExpectedPair = ab27ao27Pair
							}
							ab27ao27Labels := []string{"valid-value", "open-circuit", "closed-circuit", "high-level"}
							if ab27ao27View == "status-zero-valid-or-low" {
								ab27ao27Labels[0] = "valid-value-or-low-level"
							}
							ab27ao27FieldValue := ab27ao27Field(ab27ao27Output, fmt.Sprintf("statusS%d", ab27ao27Start+ab27ao27Sibling))
							if ab27ao27FieldValue["raw"] != float64(ab27ao27ExpectedPair) || ab27ao27FieldValue["label"] != ab27ao27Labels[ab27ao27ExpectedPair] || ab27ao27FieldValue["grade"] != "Hypothesis" {
								ab27ao27T.Fatal(ab27ao27FieldValue)
							}
						}
					})
				}
			}
			ab27ao27FirstReserved := 4
			if ab27ao27Operation == 6 {
				ab27ao27FirstReserved = 6
			}
			for ab27ao27Bit := ab27ao27FirstReserved; ab27ao27Bit < 8; ab27ao27Bit++ {
				ab27ao27T.Run(fmt.Sprintf("op%d_%s_reserved%d", ab27ao27Operation, ab27ao27View, ab27ao27Bit), func(ab27ao27T *testing.T) {
					ab27ao27Bytes := make([]int, ab27ao27Extent)
					ab27ao27Bytes[ab27ao27Status+1] = 1 << ab27ao27Bit
					ab27ao27InputValue := ab27ao27Input(ab27ao27Operation, ab27ao27Bytes)
					ab27ao27InputValue["context"].(map[string]any)["view"] = ab27ao27View
					ab27ao27Output := ab27ao27Call(ab27ao27T, ab27ao27InputValue)
					if ab27ao27Field(ab27ao27Output, "status_reserved")["reserved"] != float64(int(1)<<ab27ao27Bit) {
						ab27ao27T.Fatal("reserved bit lost")
					}
					for ab27ao27Index := 0; ab27ao27Index < ab27ao27Count; ab27ao27Index++ {
						if ab27ao27Field(ab27ao27Output, fmt.Sprintf("statusS%d", ab27ao27Start+ab27ao27Index))["raw"] != float64(0) {
							ab27ao27T.Fatal("reserved interpreted as status")
						}
					}
				})
			}
		}
	}
}

func TestAb27Ao27TruncationPWM(ab27ao27T *testing.T) {
	for _, ab27ao27Operation := range []int{3, 6, 7} {
		ab27ao27Extent := 18
		ab27ao27Positions := []int{3, 5, 7, 9, 11, 13}
		ab27ao27Start := 1
		if ab27ao27Operation == 6 {
			ab27ao27Extent = 19
			ab27ao27Positions = []int{3, 5, 7, 9, 11, 13, 15}
		}
		if ab27ao27Operation == 7 {
			ab27ao27Positions = []int{3, 5, 7, 9, 12, 14}
			ab27ao27Start = 8
		}
		for ab27ao27Width := 0; ab27ao27Width <= ab27ao27Extent+2; ab27ao27Width++ {
			ab27ao27T.Run(fmt.Sprintf("op%d_width%d", ab27ao27Operation, ab27ao27Width), func(ab27ao27T *testing.T) {
				ab27ao27Bytes := make([]int, ab27ao27Width)
				for ab27ao27Index := range ab27ao27Bytes {
					ab27ao27Bytes[ab27ao27Index] = 255
				}
				ab27ao27Output := ab27ao27Call(ab27ao27T, ab27ao27Input(ab27ao27Operation, ab27ao27Bytes))
				ab27ao27State := "partial"
				if ab27ao27Width >= ab27ao27Extent {
					ab27ao27State = "Hypothesis"
				}
				if ab27ao27Output["state"] != ab27ao27State {
					ab27ao27T.Fatal("extent state")
				}
				for ab27ao27Index, ab27ao27Position := range ab27ao27Positions {
					ab27ao27FieldValue := ab27ao27Field(ab27ao27Output, fmt.Sprintf("S%d", ab27ao27Start+ab27ao27Index))
					if ab27ao27Position+2 <= ab27ao27Width {
						if ab27ao27FieldValue["number"] != -0.0625 {
							ab27ao27T.Fatal("complete sibling discarded")
						}
					} else {
						if ab27ao27FieldValue["state"] != "incomplete" {
							ab27ao27T.Fatal("incomplete state")
						}
						if _, ab27ao27Present := ab27ao27FieldValue["number"]; ab27ao27Present {
							ab27ao27T.Fatal("missing word defaulted")
						}
					}
				}
				if ab27ao27Operation == 3 && ab27ao27Width >= 16 {
					if ab27ao27Field(ab27ao27Output, "PWM")["number"] != 127.5 {
						ab27ao27T.Fatal("PWM clipped")
					}
				}
			})
		}
	}
}

func TestAb27Ao27SelectionPresence(ab27ao27T *testing.T) {
	for _, ab27ao27Key := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
		ab27ao27Variants := []any{nil, 0, 182}
		if ab27ao27Key == "secondary" {
			ab27ao27Variants = []any{nil, 0, 36}
		}
		if ab27ao27Key == "operation" {
			ab27ao27Variants = []any{nil, 0, 8}
		}
		if ab27ao27Key == "direction" || ab27ao27Key == "view" {
			ab27ao27Variants = []any{nil, "", "unknown"}
		}
		if ab27ao27Key == "family" || ab27ao27Key == "variant" {
			ab27ao27Variants = []any{nil, "", map[string]any{"different": true}}
		}
		for ab27ao27Index, ab27ao27Value := range ab27ao27Variants {
			ab27ao27T.Run(fmt.Sprintf("%s_%d", ab27ao27Key, ab27ao27Index), func(ab27ao27T *testing.T) {
				ab27ao27InputValue := ab27ao27Input(3, make([]int, 18))
				ab27ao27Context := ab27ao27InputValue["context"].(map[string]any)
				if ab27ao27Value == nil {
					delete(ab27ao27Context, ab27ao27Key)
				} else {
					ab27ao27Context[ab27ao27Key] = ab27ao27Value
				}
				ab27ao27Output := ab27ao27Call(ab27ao27T, ab27ao27InputValue)
				if ab27ao27Key == "family" || ab27ao27Key == "variant" {
					if ab27ao27Output["state"] != "Hypothesis" {
						ab27ao27T.Fatal("opaque context selects")
					}
				} else {
					if ab27ao27Output["state"] != "Unknown" || len(ab27ao27Output["fields"].(map[string]any)) != 0 {
						ab27ao27T.Fatal("guessed selection")
					}
				}
			})
		}
	}
	for _, ab27ao27Key := range []string{"logical_object", "unmodelled"} {
		for _, ab27ao27Mode := range []string{"missing", "empty", "zero", "changed"} {
			ab27ao27T.Run(ab27ao27Key+"_"+ab27ao27Mode, func(ab27ao27T *testing.T) {
				ab27ao27InputValue := ab27ao27Input(3, make([]int, 18))
				ab27ao27Raw := ab27ao27InputValue["fields"].(map[string]any)
				switch ab27ao27Mode {
				case "missing":
					delete(ab27ao27Raw, ab27ao27Key)
				case "empty":
					ab27ao27Raw[ab27ao27Key] = []int{}
				case "zero":
					ab27ao27Raw[ab27ao27Key] = []int{0}
				case "changed":
					ab27ao27Raw[ab27ao27Key] = []int{99, 1, 7, 16, 0}
				}
				ab27ao27Output := ab27ao27Call(ab27ao27T, ab27ao27InputValue)
				if ab27ao27Key == "unmodelled" && ab27ao27Output["state"] != "Hypothesis" {
					ab27ao27T.Fatal("unmodelled field selected")
				}
			})
		}
	}
	ab27ao27InputValue := ab27ao27Input(3, make([]int, 18))
	ab27ao27InputValue["arguments"] = map[string]any{}
	ab27ao27Call(ab27ao27T, ab27ao27InputValue)
}

func TestAb27Ao27InvalidJSON(ab27ao27T *testing.T) {
	for ab27ao27Index, ab27ao27InputJSON := range []string{
		``, `null`, `[]`, `{`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"family":{"x":1,"x":2}},"fields":{}}`, `{"context":{},"fields":{"x":[],"x":[]}}`, `{"context":{},"fields":{},"arguments":{"x":1,"x":2}}`, `{"context":{},"fields":{},"fields":{}}`, `{"context":{},"fields":{}} true`, `{"context":[],"fields":{}}`, `{"context":{},"fields":null}`, `{"context":{},"fields":[]}`, `{"context":{},"fields":{"x":null}}`, `{"context":{},"fields":{"x":"AA=="}}`, `{"context":{},"fields":{"x":[-1]}}`, `{"context":{},"fields":{"x":[256]}}`, `{"context":{},"fields":{"x":[1.5]}}`, `{"context":{},"fields":{"x":[1e0]}}`, `{"context":{},"fields":{"x":[true]}}`, `{"context":{},"fields":{"x":[null]}}`, `{"context":{},"fields":{"x":["1"]}}`, `{"context":{},"fields":{"x":[{}]}}`, `{"context":{},"fields":{"x":[[]]}}`, `{"context":{},"fields":{},"arguments":null}`, `{"context":{},"fields":{},"arguments":[]}`, `{"context":{},"fields":{},"arguments":{"flag":true}}`, `{"context":{},"fields":{},"extra":0}`,
	} {
		ab27ao27T.Run(fmt.Sprint(ab27ao27Index), func(ab27ao27T *testing.T) {
			ab27ao27Result, ab27ao27Err := Ab27Ao27Interpret(json.RawMessage(ab27ao27InputJSON))
			if ab27ao27Err == nil || ab27ao27Result != nil {
				ab27ao27T.Fatal("invalid JSON admitted")
			}
		})
	}
	for _, ab27ao27Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, ab27ao27Value := range []any{nil, true, []any{}, map[string]any{}, 1.5} {
			ab27ao27InputValue := ab27ao27Input(3, make([]int, 18))
			ab27ao27InputValue["context"].(map[string]any)[ab27ao27Key] = ab27ao27Value
			ab27ao27InputJSON, ab27ao27Err := json.Marshal(ab27ao27InputValue)
			if ab27ao27Err != nil {
				ab27ao27T.Fatal(ab27ao27Err)
			}
			ab27ao27Result, ab27ao27Err := Ab27Ao27Interpret(ab27ao27InputJSON)
			if ab27ao27Err == nil || ab27ao27Result != nil {
				ab27ao27T.Fatal("wrong context kind admitted", ab27ao27Key)
			}
		}
	}
}

func TestAb27Ao27Ownership(ab27ao27T *testing.T) {
	ab27ao27InputValue := ab27ao27Input(7, make([]int, 18))
	ab27ao27InputValue["arguments"] = map[string]any{}
	ab27ao27InputJSON, ab27ao27Err := json.Marshal(ab27ao27InputValue)
	if ab27ao27Err != nil {
		ab27ao27T.Fatal(ab27ao27Err)
	}
	ab27ao27Control := append(json.RawMessage(nil), ab27ao27InputJSON...)
	ab27ao27First, ab27ao27Err := Ab27Ao27Interpret(ab27ao27InputJSON)
	if ab27ao27Err != nil {
		ab27ao27T.Fatal(ab27ao27Err)
	}
	ab27ao27Saved := append(json.RawMessage(nil), ab27ao27First...)
	for ab27ao27Index := range ab27ao27InputJSON {
		ab27ao27InputJSON[ab27ao27Index] = ' '
	}
	if !reflect.DeepEqual(ab27ao27First, ab27ao27Saved) {
		ab27ao27T.Fatal("input aliases output")
	}
	ab27ao27Second, ab27ao27Err := Ab27Ao27Interpret(ab27ao27Control)
	if ab27ao27Err != nil {
		ab27ao27T.Fatal(ab27ao27Err)
	}
	for ab27ao27Index := range ab27ao27First {
		ab27ao27First[ab27ao27Index] = ' '
	}
	if !reflect.DeepEqual(ab27ao27Second, ab27ao27Saved) {
		ab27ao27T.Fatal("output aliases another result")
	}
	var ab27ao27Nested map[string]any
	if ab27ao27Err := json.Unmarshal(ab27ao27Second, &ab27ao27Nested); ab27ao27Err != nil {
		ab27ao27T.Fatal(ab27ao27Err)
	}
	ab27ao27Nested["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "changed"
	ab27ao27Nested["raw"].(map[string]any)["logical_object"].([]any)[0] = float64(99)
	ab27ao27Nested["arguments"].(map[string]any)["changed"] = true
	ab27ao27Field(ab27ao27Nested, "S8")["number"] = float64(999)
	ab27ao27Third, ab27ao27Err := Ab27Ao27Interpret(ab27ao27Control)
	if ab27ao27Err != nil || !reflect.DeepEqual(ab27ao27Third, ab27ao27Saved) {
		ab27ao27T.Fatal("nested mutation affects later call")
	}
}

func TestAb27Ao27OpaquePrefixAndPWM(ab27ao27T *testing.T) {
	for ab27ao27Position := 0; ab27ao27Position < 3; ab27ao27Position++ {
		for _, ab27ao27Byte := range []int{0, 3, 6, 7, 181, 255} {
			ab27ao27T.Run(fmt.Sprintf("prefix%d_%d", ab27ao27Position, ab27ao27Byte), func(ab27ao27T *testing.T) {
				ab27ao27Bytes := make([]int, 18)
				ab27ao27Bytes[ab27ao27Position] = ab27ao27Byte
				ab27ao27Output := ab27ao27Call(ab27ao27T, ab27ao27Input(7, ab27ao27Bytes))
				if ab27ao27Output["state"] != "Hypothesis" || ab27ao27Field(ab27ao27Output, "S12")["number"] != float64(0) {
					ab27ao27T.Fatal("prefix selected operation")
				}
			})
		}
	}
	for _, ab27ao27Case := range []struct {
		ab27ao27Byte int
		ab27ao27Want float64
	}{{0, 0}, {1, 0.5}, {200, 100}, {255, 127.5}} {
		ab27ao27Bytes := make([]int, 18)
		ab27ao27Bytes[15] = ab27ao27Case.ab27ao27Byte
		ab27ao27Output := ab27ao27Call(ab27ao27T, ab27ao27Input(3, ab27ao27Bytes))
		if ab27ao27Field(ab27ao27Output, "PWM")["number"] != ab27ao27Case.ab27ao27Want {
			ab27ao27T.Fatal("PWM carrier")
		}
	}
	for _, ab27ao27Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		ab27ao27InputValue := ab27ao27Input(3, make([]int, 18))
		if ab27ao27Key == "direction" || ab27ao27Key == "view" {
			ab27ao27InputValue["context"].(map[string]any)[ab27ao27Key] = 0
		} else {
			ab27ao27InputValue["context"].(map[string]any)[ab27ao27Key] = "3"
		}
		ab27ao27JSON, ab27ao27Err := json.Marshal(ab27ao27InputValue)
		if ab27ao27Err != nil {
			ab27ao27T.Fatal(ab27ao27Err)
		}
		ab27ao27Result, ab27ao27Err := Ab27Ao27Interpret(ab27ao27JSON)
		if ab27ao27Err == nil || ab27ao27Result != nil {
			ab27ao27T.Fatal("context kind admitted")
		}
	}
}

func TestAb27Ao27PartialIndependentFields(ab27ao27T *testing.T) {
	for _, ab27ao27Operation := range []int{3, 6, 7} {
		ab27ao27Extent, ab27ao27Status, ab27ao27Count, ab27ao27Start := 18, 16, 6, 1
		if ab27ao27Operation == 6 {
			ab27ao27Extent, ab27ao27Status, ab27ao27Count = 19, 17, 7
		}
		if ab27ao27Operation == 7 {
			ab27ao27Start = 8
		}
		for ab27ao27Width := 0; ab27ao27Width <= ab27ao27Extent; ab27ao27Width++ {
			ab27ao27Bytes := make([]int, ab27ao27Width)
			for ab27ao27Index := range ab27ao27Bytes {
				ab27ao27Bytes[ab27ao27Index] = 255
			}
			ab27ao27Output := ab27ao27Call(ab27ao27T, ab27ao27Input(ab27ao27Operation, ab27ao27Bytes))
			for ab27ao27Index := 0; ab27ao27Index < ab27ao27Count; ab27ao27Index++ {
				ab27ao27FieldValue := ab27ao27Field(ab27ao27Output, fmt.Sprintf("statusS%d", ab27ao27Start+ab27ao27Index))
				if ab27ao27Width > ab27ao27Status+ab27ao27Index/4 {
					if ab27ao27FieldValue["label"] != "high-level" || ab27ao27FieldValue["raw"] != float64(3) {
						ab27ao27T.Fatal("complete status lost")
					}
				} else {
					if ab27ao27FieldValue["state"] != "incomplete" {
						ab27ao27T.Fatal("status defaulted")
					}
					if _, ab27ao27Present := ab27ao27FieldValue["label"]; ab27ao27Present {
						ab27ao27T.Fatal("missing status labelled")
					}
				}
			}
			ab27ao27Reserved := ab27ao27Field(ab27ao27Output, "status_reserved")
			if ab27ao27Width < ab27ao27Extent {
				if ab27ao27Reserved["state"] != "incomplete" {
					ab27ao27T.Fatal("reserved defaulted")
				}
			} else {
				ab27ao27Want := float64(240)
				if ab27ao27Operation == 6 {
					ab27ao27Want = 192
				}
				if ab27ao27Reserved["reserved"] != ab27ao27Want {
					ab27ao27T.Fatal("reserved carrier")
				}
			}
			if ab27ao27Operation == 3 || ab27ao27Operation == 7 {
				ab27ao27Name, ab27ao27Position := "PWM", 15
				if ab27ao27Operation == 7 {
					ab27ao27Name, ab27ao27Position = "opaque11", 11
				}
				ab27ao27FieldValue := ab27ao27Field(ab27ao27Output, ab27ao27Name)
				if ab27ao27Width <= ab27ao27Position {
					if ab27ao27FieldValue["state"] != "incomplete" {
						ab27ao27T.Fatal("byte field defaulted")
					}
					if _, ab27ao27Present := ab27ao27FieldValue["number"]; ab27ao27Present {
						ab27ao27T.Fatal("incomplete byte converted")
					}
				} else if ab27ao27FieldValue["raw"] != float64(255) {
					ab27ao27T.Fatal("complete byte lost")
				}
			}
		}
	}
	ab27ao27Result, ab27ao27Err := Ab27Ao27Interpret(json.RawMessage{'{', '"', 255, '"', ':', '0', '}'})
	if ab27ao27Err == nil || ab27ao27Result != nil {
		ab27ao27T.Fatal("invalid UTF-8 admitted")
	}
}
