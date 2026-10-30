package vaillant

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func ab17ao17Input(ab17ao17Bytes []int) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 19, "operation": 7, "direction": "positive", "view": "information", "family": map[string]any{"opaque": []any{"synthetic-family", 0}}, "variant": "synthetic-variant"}, "fields": map[string]any{"logical_object": ab17ao17Bytes, "unmodelled": []int{0, 255}}}
}
func ab17ao17Run(ab17ao17T *testing.T, ab17ao17InputValue any) map[string]any {
	ab17ao17T.Helper()
	ab17ao17Encoded, ab17ao17Error := json.Marshal(ab17ao17InputValue)
	if ab17ao17Error != nil {
		ab17ao17T.Fatal(ab17ao17Error)
	}
	ab17ao17Result, ab17ao17Error := Ab17Ao17Interpret(ab17ao17Encoded)
	if ab17ao17Error != nil {
		ab17ao17T.Fatal(ab17ao17Error)
	}
	ab17ao17Output := map[string]any{}
	if ab17ao17Error = json.Unmarshal(ab17ao17Result, &ab17ao17Output); ab17ao17Error != nil {
		ab17ao17T.Fatal(ab17ao17Error)
	}
	if ab17ao17Output["physical"] != "Unknown" || ab17ao17Output["qualified"] != false {
		ab17ao17T.Fatal("qualification inferred", ab17ao17Output)
	}
	ab17ao17Normalized := map[string]any{}
	_ = json.Unmarshal(ab17ao17Encoded, &ab17ao17Normalized)
	if !reflect.DeepEqual(ab17ao17Output["context"], ab17ao17Normalized["context"]) || !reflect.DeepEqual(ab17ao17Output["raw"], ab17ao17Normalized["fields"]) {
		ab17ao17T.Fatal("context or raw changed", ab17ao17Output)
	}
	if ab17ao17Arguments, ab17ao17Present := ab17ao17Normalized["arguments"]; ab17ao17Present {
		if !reflect.DeepEqual(ab17ao17Arguments, ab17ao17Output["arguments"]) {
			ab17ao17T.Fatal("arguments changed")
		}
	} else if _, ab17ao17Present := ab17ao17Output["arguments"]; ab17ao17Present {
		ab17ao17T.Fatal("arguments invented")
	}
	return ab17ao17Output
}
func ab17ao17Field(ab17ao17Output map[string]any, ab17ao17Name string) map[string]any {
	return ab17ao17Output["fields"].(map[string]any)[ab17ao17Name].(map[string]any)
}
func TestAb17Ao17Flags(ab17ao17T *testing.T) {
	ab17ao17Names := []string{"capabilities", "fuel", "storage", "features"}
	ab17ao17Labels := [][]string{{"floor-standing", "condensing", "fan-assisted", "open-vented", "high-modulation-range", "water-heater", "step-modulation", "ECS"}, {"gas", "oil", "ABC"}, {"layer-tank", "warm-start", "external-tank", "combination"}, {"solar-compact", "basic-solar", "AAPS", "advanced-power", "partial-load", "parameter-master"}}
	for ab17ao17Index, ab17ao17Name := range ab17ao17Names {
		for ab17ao17Bit := 0; ab17ao17Bit < 8; ab17ao17Bit++ {
			ab17ao17T.Run(fmt.Sprintf("%s_bit%d", ab17ao17Name, ab17ao17Bit), func(ab17ao17T *testing.T) {
				ab17ao17Bytes := []int{91, 82, 73, 0, 0, 0, 0, 0, 0, 7, 8, 9, 201}
				ab17ao17Bytes[4+ab17ao17Index] = 1 << ab17ao17Bit
				ab17ao17Output := ab17ao17Run(ab17ao17T, ab17ao17Input(ab17ao17Bytes))
				if ab17ao17Output["state"] != "Hypothesis" {
					ab17ao17T.Fatal(ab17ao17Output)
				}
				for ab17ao17Other, ab17ao17OtherName := range ab17ao17Names {
					ab17ao17Expected := map[string]any{}
					for ab17ao17LabelIndex, ab17ao17Label := range ab17ao17Labels[ab17ao17Other] {
						ab17ao17Expected[ab17ao17Label] = ab17ao17Other == ab17ao17Index && ab17ao17LabelIndex == ab17ao17Bit
					}
					ab17ao17FieldValue := ab17ao17Field(ab17ao17Output, ab17ao17OtherName)
					ab17ao17Reserved := 0
					if ab17ao17Other == ab17ao17Index && ab17ao17Bit >= len(ab17ao17Labels[ab17ao17Other]) {
						ab17ao17Reserved = 1 << ab17ao17Bit
					}
					if !reflect.DeepEqual(ab17ao17FieldValue["flags"], ab17ao17Expected) || ab17ao17FieldValue["reserved"] != float64(ab17ao17Reserved) || ab17ao17FieldValue["state"] != "Hypothesis" || ab17ao17FieldValue["number"] != float64(ab17ao17Bytes[4+ab17ao17Other]) {
						ab17ao17T.Fatal(ab17ao17OtherName, ab17ao17FieldValue)
					}
				}
				ab17ao17Opaque := ab17ao17Field(ab17ao17Output, "opaque")
				if !reflect.DeepEqual(ab17ao17Opaque, map[string]any{"raw": []any{float64(7), float64(8), float64(9)}, "state": "Unknown"}) {
					ab17ao17T.Fatal(ab17ao17Opaque)
				}
			})
		}
	}
}
func TestAb17Ao17EnumsAndPartial(ab17ao17T *testing.T) {
	for _, ab17ao17Case := range []struct {
		ab17ao17Name     string
		ab17ao17Position int
		ab17ao17Values   []int
		ab17ao17Labels   []string
	}{{"type", 3, []int{0, 1, 2, 3, 255}, []string{"modulating", "single-step", "two-step"}}, {"pump_mode", 8, []int{0, 1, 2}, []string{"pump-mode0", "pump-mode1"}}} {
		for _, ab17ao17Value := range ab17ao17Case.ab17ao17Values {
			ab17ao17T.Run(fmt.Sprintf("%s_%d", ab17ao17Case.ab17ao17Name, ab17ao17Value), func(ab17ao17T *testing.T) {
				ab17ao17Bytes := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 9}
				ab17ao17Bytes[ab17ao17Case.ab17ao17Position] = ab17ao17Value
				ab17ao17Output := ab17ao17Run(ab17ao17T, ab17ao17Input(ab17ao17Bytes))
				ab17ao17Actual := ab17ao17Field(ab17ao17Output, ab17ao17Case.ab17ao17Name)
				if ab17ao17Value < len(ab17ao17Case.ab17ao17Labels) {
					if ab17ao17Actual["number"] != float64(ab17ao17Value) || ab17ao17Actual["label"] != ab17ao17Case.ab17ao17Labels[ab17ao17Value] || ab17ao17Actual["state"] != "Hypothesis" {
						ab17ao17T.Fatal(ab17ao17Actual)
					}
				} else {
					if ab17ao17Actual["state"] != "Unknown" || ab17ao17Actual["label"] != nil || ab17ao17Actual["number"] != nil || ab17ao17Output["state"] != "partial" {
						ab17ao17T.Fatal(ab17ao17Actual)
					}
				}
				if ab17ao17Field(ab17ao17Output, "fuel")["state"] != "Hypothesis" {
					ab17ao17T.Fatal("unknown sibling erased fuel")
				}
			})
		}
	}
	for ab17ao17Length := 0; ab17ao17Length < 12; ab17ao17Length++ {
		ab17ao17T.Run(fmt.Sprintf("length_%d", ab17ao17Length), func(ab17ao17T *testing.T) {
			ab17ao17Bytes := []int{201, 202, 203, 2, 255, 7, 15, 63, 1, 7, 8, 9}
			ab17ao17Output := ab17ao17Run(ab17ao17T, ab17ao17Input(ab17ao17Bytes[:ab17ao17Length]))
			if ab17ao17Output["state"] != "partial" {
				ab17ao17T.Fatal(ab17ao17Output)
			}
			for ab17ao17Index, ab17ao17Name := range []string{"type", "capabilities", "fuel", "storage", "features", "pump_mode"} {
				ab17ao17Actual := ab17ao17Field(ab17ao17Output, ab17ao17Name)
				if ab17ao17Length > ab17ao17Index+3 {
					if ab17ao17Actual["number"] != float64(ab17ao17Bytes[ab17ao17Index+3]) || ab17ao17Actual["state"] != "Hypothesis" {
						ab17ao17T.Fatal(ab17ao17Name, ab17ao17Actual)
					}
				} else if ab17ao17Actual["number"] != nil || ab17ao17Actual["state"] != "Unknown" {
					ab17ao17T.Fatal(ab17ao17Actual)
				}
			}
		})
	}
}
func TestAb17Ao17Feedback(ab17ao17T *testing.T) {
	for _, ab17ao17Operation := range []int{5, 6} {
		for _, ab17ao17Value := range []int{0, 1, 2, 255} {
			ab17ao17T.Run(fmt.Sprintf("op%d_%d", ab17ao17Operation, ab17ao17Value), func(ab17ao17T *testing.T) {
				ab17ao17InputValue := ab17ao17Input([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 9})
				ab17ao17Context := ab17ao17InputValue["context"].(map[string]any)
				ab17ao17Context["operation"] = ab17ao17Operation
				ab17ao17Context["view"] = "feedback"
				ab17ao17InputValue["fields"].(map[string]any)["feedback"] = []int{ab17ao17Value}
				ab17ao17Output := ab17ao17Run(ab17ao17T, ab17ao17InputValue)
				ab17ao17Actual := ab17ao17Field(ab17ao17Output, "feedback")
				ab17ao17Expected := "rejected"
				if ab17ao17Operation == 5 {
					ab17ao17Expected = "not-accepted-or-not-supported"
				}
				if ab17ao17Value == 1 {
					ab17ao17Expected = "accepted"
				}
				if ab17ao17Actual["label"] != ab17ao17Expected || ab17ao17Actual["number"] != float64(ab17ao17Value) || ab17ao17Actual["state"] != "Hypothesis" || len(ab17ao17Output["fields"].(map[string]any)) != 1 {
					ab17ao17T.Fatal(ab17ao17Output)
				}
				for _, ab17ao17Width := range []int{0, 2, 3} {
					ab17ao17InputValue["fields"].(map[string]any)["feedback"] = make([]int, ab17ao17Width)
					ab17ao17Output = ab17ao17Run(ab17ao17T, ab17ao17InputValue)
					if ab17ao17Output["state"] != "partial" || ab17ao17Field(ab17ao17Output, "feedback")["number"] != nil {
						ab17ao17T.Fatal(ab17ao17Output)
					}
				}
				delete(ab17ao17InputValue["fields"].(map[string]any), "feedback")
				if ab17ao17Run(ab17ao17T, ab17ao17InputValue)["state"] != "partial" {
					ab17ao17T.Fatal("missing feedback")
				}
			})
		}
	}
}
func TestAb17Ao17Selection(ab17ao17T *testing.T) {
	for ab17ao17Key, ab17ao17Values := range map[string][]any{"primary": {0, 180, 182}, "secondary": {0, 18, 20}, "operation": {0, 4, 5, 6, 8}, "direction": {"request", "unknown", ""}, "view": {"feedback", "unknown", ""}} {
		for _, ab17ao17Value := range ab17ao17Values {
			ab17ao17T.Run(fmt.Sprint(ab17ao17Key, ab17ao17Value), func(ab17ao17T *testing.T) {
				ab17ao17InputValue := ab17ao17Input([]int{181, 19, 7, 0, 255, 7, 15, 63, 1, 7, 8, 9})
				ab17ao17InputValue["context"].(map[string]any)[ab17ao17Key] = ab17ao17Value
				ab17ao17Output := ab17ao17Run(ab17ao17T, ab17ao17InputValue)
				if ab17ao17Output["state"] != "Unknown" || len(ab17ao17Output["fields"].(map[string]any)) != 0 {
					ab17ao17T.Fatal(ab17ao17Output)
				}
			})
		}
	}
	for _, ab17ao17Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		ab17ao17InputValue := ab17ao17Input([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 9})
		delete(ab17ao17InputValue["context"].(map[string]any), ab17ao17Key)
		if ab17ao17Run(ab17ao17T, ab17ao17InputValue)["state"] != "Unknown" {
			ab17ao17T.Fatal(ab17ao17Key)
		}
	}
	for _, ab17ao17Key := range []string{"family", "variant"} {
		for _, ab17ao17Value := range []any{nil, 0, "other", map[string]any{"opaque": []int{0, 255}}} {
			ab17ao17InputValue := ab17ao17Input([]int{255, 254, 253, 0, 0, 0, 0, 0, 0, 7, 8, 9})
			ab17ao17InputValue["context"].(map[string]any)[ab17ao17Key] = ab17ao17Value
			if ab17ao17Run(ab17ao17T, ab17ao17InputValue)["state"] != "Hypothesis" {
				ab17ao17T.Fatal("opaque family selected interpretation")
			}
		}
	}
	for _, ab17ao17Key := range []string{"logical_object", "unmodelled"} {
		ab17ao17InputValue := ab17ao17Input([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 9})
		ab17ao17InputValue["fields"].(map[string]any)[ab17ao17Key] = []int{}
		ab17ao17Run(ab17ao17T, ab17ao17InputValue)
		delete(ab17ao17InputValue["fields"].(map[string]any), ab17ao17Key)
		ab17ao17Run(ab17ao17T, ab17ao17InputValue)
	}
	for _, ab17ao17Operation := range []int{5, 6} {
		ab17ao17InputValue := ab17ao17Input([]int{})
		ab17ao17Context := ab17ao17InputValue["context"].(map[string]any)
		ab17ao17Context["operation"] = ab17ao17Operation
		ab17ao17Context["view"] = "feedback"
		ab17ao17Context["direction"] = "request"
		ab17ao17InputValue["fields"].(map[string]any)["feedback"] = []int{1}
		if ab17ao17Run(ab17ao17T, ab17ao17InputValue)["state"] != "Unknown" {
			ab17ao17T.Fatal("request accepted")
		}
	}
}
func TestAb17Ao17Invalid(ab17ao17T *testing.T) {
	for ab17ao17Index, ab17ao17InputValue := range []string{``, `null`, `[]`, `{}`, `{"context":{},"fields":{}}{}`, `{"context":{},"fields":{},}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"family":{"x":0,"x":1}},"fields":{}}`, `{"context":{},"fields":{"x":[],"x":[]}}`, `{"context":{},"fields":{},"arguments":{"x":0,"x":1}}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":null}`, `{"context":{},"fields":{"x":null}}`, `{"context":{},"fields":{"x":"AA=="}}`, `{"context":{},"fields":{"x":[-1]}}`, `{"context":{},"fields":{"x":[256]}}`, `{"context":{},"fields":{"x":[1.0]}}`, `{"context":{},"fields":{"x":[1e0]}}`, `{"context":{},"fields":{"x":[true]}}`, `{"context":{},"fields":{"x":["0"]}}`, `{"context":{},"fields":{"x":[{}]}}`, `{"context":{},"fields":{"x":[[]]}}`, `{"context":{},"fields":{},"arguments":null}`, `{"context":{},"fields":{},"arguments":{"flag":1}}`, `{"context":{"primary":"181"},"fields":{}}`, `{"context":{"secondary":19.0},"fields":{}}`, `{"context":{"operation":true},"fields":{}}`, `{"context":{"direction":0},"fields":{}}`, `{"context":{"view":null},"fields":{}}`} {
		ab17ao17T.Run(fmt.Sprint(ab17ao17Index), func(ab17ao17T *testing.T) {
			ab17ao17Output, ab17ao17Error := Ab17Ao17Interpret(json.RawMessage(ab17ao17InputValue))
			if ab17ao17Error == nil || ab17ao17Output != nil {
				ab17ao17T.Fatal("invalid accepted", ab17ao17InputValue, string(ab17ao17Output))
			}
		})
	}
}
func TestAb17Ao17Isolation(ab17ao17T *testing.T) {
	ab17ao17InputValue := ab17ao17Input([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 9, 255})
	ab17ao17InputValue["arguments"] = map[string]any{}
	ab17ao17Encoded, _ := json.Marshal(ab17ao17InputValue)
	ab17ao17Original := append(json.RawMessage(nil), ab17ao17Encoded...)
	ab17ao17First, ab17ao17Error := Ab17Ao17Interpret(ab17ao17Encoded)
	if ab17ao17Error != nil {
		ab17ao17T.Fatal(ab17ao17Error)
	}
	ab17ao17Saved := append(json.RawMessage(nil), ab17ao17First...)
	for ab17ao17Index := range ab17ao17Encoded {
		ab17ao17Encoded[ab17ao17Index] = ' '
	}
	if !reflect.DeepEqual(ab17ao17Saved, ab17ao17First) {
		ab17ao17T.Fatal("input aliases output")
	}
	ab17ao17Decoded := ab17ao17Run(ab17ao17T, ab17ao17InputValue)
	ab17ao17Decoded["raw"].(map[string]any)["logical_object"].([]any)[3] = float64(255)
	if ab17ao17Field(ab17ao17Decoded, "type")["raw"].([]any)[0] != float64(0) {
		ab17ao17T.Fatal("raw aliases semantic field")
	}
	ab17ao17Decoded["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "changed"
	ab17ao17Field(ab17ao17Decoded, "fuel")["flags"].(map[string]any)["gas"] = true
	for ab17ao17Index := range ab17ao17First {
		ab17ao17First[ab17ao17Index] = ' '
	}
	ab17ao17Second, ab17ao17Error := Ab17Ao17Interpret(ab17ao17Original)
	if ab17ao17Error != nil || !reflect.DeepEqual(ab17ao17Second, ab17ao17Saved) {
		ab17ao17T.Fatal("later call changed", ab17ao17Error)
	}
}
func TestAb17Ao17Fixtures(ab17ao17T *testing.T) {
	ab17ao17Data, ab17ao17Error := os.ReadFile("testdata/ab17ao17_b513_information_views.json")
	if ab17ao17Error != nil {
		ab17ao17T.Fatal(ab17ao17Error)
	}
	var ab17ao17Cases []struct {
		Ab17Ao17ID    string          `json:"id"`
		Ab17Ao17Input json.RawMessage `json:"input"`
		Ab17Ao17State string          `json:"state"`
	}
	if ab17ao17Error = json.Unmarshal(ab17ao17Data, &ab17ao17Cases); ab17ao17Error != nil || len(ab17ao17Cases) == 0 {
		ab17ao17T.Fatal("invalid fixtures", ab17ao17Error)
	}
	for _, ab17ao17Case := range ab17ao17Cases {
		if !strings.HasPrefix(ab17ao17Case.Ab17Ao17ID, "ab17ao17") {
			ab17ao17T.Fatal("fixture ID")
		}
		ab17ao17T.Run(ab17ao17Case.Ab17Ao17ID, func(ab17ao17T *testing.T) {
			var ab17ao17InputValue any
			_ = json.Unmarshal(ab17ao17Case.Ab17Ao17Input, &ab17ao17InputValue)
			if ab17ao17Run(ab17ao17T, ab17ao17InputValue)["state"] != ab17ao17Case.Ab17Ao17State {
				ab17ao17T.Fatal("fixture state")
			}
		})
	}
}

func TestAb17Ao17RawAxes(ab17ao17T *testing.T) {
	for _, ab17ao17Position := range []int{0, 1, 2, 9, 10, 11, 12} {
		for _, ab17ao17Value := range []int{0, 255} {
			ab17ao17T.Run(fmt.Sprintf("position%d_%d", ab17ao17Position, ab17ao17Value), func(ab17ao17T *testing.T) {
				ab17ao17Bytes := []int{91, 82, 73, 2, 255, 7, 15, 63, 1, 7, 8, 9, 201}
				ab17ao17Bytes[ab17ao17Position] = ab17ao17Value
				ab17ao17Output := ab17ao17Run(ab17ao17T, ab17ao17Input(ab17ao17Bytes))
				if ab17ao17Output["state"] != "Hypothesis" {
					ab17ao17T.Fatal("opaque position influenced admission")
				}
				ab17ao17Opaque := ab17ao17Field(ab17ao17Output, "opaque")
				if !reflect.DeepEqual(ab17ao17Opaque["raw"], []any{float64(ab17ao17Bytes[9]), float64(ab17ao17Bytes[10]), float64(ab17ao17Bytes[11])}) || ab17ao17Opaque["number"] != nil || ab17ao17Opaque["label"] != nil || ab17ao17Opaque["state"] != "Unknown" {
					ab17ao17T.Fatal(ab17ao17Opaque)
				}
				for ab17ao17Index, ab17ao17Name := range []string{"type", "capabilities", "fuel", "storage", "features", "pump_mode"} {
					if !reflect.DeepEqual(ab17ao17Field(ab17ao17Output, ab17ao17Name)["raw"], []any{float64(ab17ao17Bytes[3+ab17ao17Index])}) {
						ab17ao17T.Fatal("field raw differs", ab17ao17Name)
					}
				}
			})
		}
	}
	for _, ab17ao17Value := range [][]int{{}, {0}, {255}, {0, 255, 1}} {
		ab17ao17InputValue := ab17ao17Input([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 9})
		ab17ao17InputValue["fields"].(map[string]any)["unmodelled"] = ab17ao17Value
		if ab17ao17Run(ab17ao17T, ab17ao17InputValue)["state"] != "Hypothesis" {
			ab17ao17T.Fatal("unmodelled bytes selected view")
		}
	}
	ab17ao17InputValue := ab17ao17Input([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 8, 9})
	ab17ao17InputValue["context"].(map[string]any)["additional"] = map[string]any{"nested": []any{nil, true, "opaque", 0}}
	ab17ao17Run(ab17ao17T, ab17ao17InputValue)
}
