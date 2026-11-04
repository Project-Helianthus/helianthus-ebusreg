package vaillant

import (
	ab15ao15json "encoding/json"
	ab15ao15os "os"
	ab15ao15reflect "reflect"
	ab15ao15strings "strings"
	ab15ao15testing "testing"
)

func ab15ao15Control() map[string]any {
	return map[string]any{
		"context": map[string]any{"primary": 181, "secondary": 18, "operation": 1, "direction": "request", "view": "directional-fields", "family": map[string]any{"opaque": []any{"synthetic-family", 0}}, "variant": []any{"synthetic-variant", false}, "setpoint": []int{17}, "feedback": []int{17}},
		"fields":  map[string]any{"secondary_heating_pump": []int{11}, "circulation_pump": []int{22}, "solar_pump": []int{33}, "error_relay": []int{15}, "loading_pump": []int{44}, "fume_flap": []int{240}, "unmodelled": []int{0, 255, 19}},
	}
}

func ab15ao15Encode(ab15ao15T *ab15ao15testing.T, ab15ao15Value any) ab15ao15json.RawMessage {
	ab15ao15T.Helper()
	ab15ao15Raw, ab15ao15Err := ab15ao15json.Marshal(ab15ao15Value)
	if ab15ao15Err != nil {
		ab15ao15T.Fatal(ab15ao15Err)
	}
	return ab15ao15Raw
}

func ab15ao15Decode(ab15ao15T *ab15ao15testing.T, ab15ao15Raw ab15ao15json.RawMessage) map[string]any {
	ab15ao15T.Helper()
	var ab15ao15Value map[string]any
	if ab15ao15Err := ab15ao15json.Unmarshal(ab15ao15Raw, &ab15ao15Value); ab15ao15Err != nil {
		ab15ao15T.Fatal(ab15ao15Err)
	}
	return ab15ao15Value
}

func ab15ao15Run(ab15ao15T *ab15ao15testing.T, ab15ao15Input map[string]any) map[string]any {
	ab15ao15T.Helper()
	ab15ao15Raw := ab15ao15Encode(ab15ao15T, ab15ao15Input)
	ab15ao15Result, ab15ao15Err := Ab15Ao15Interpret(ab15ao15Raw)
	if ab15ao15Err != nil {
		ab15ao15T.Fatal(ab15ao15Err)
	}
	ab15ao15Output := ab15ao15Decode(ab15ao15T, ab15ao15Result)
	ab15ao15Normalized := ab15ao15Decode(ab15ao15T, ab15ao15Raw)
	ab15ao15WantKeys := 6
	if _, ab15ao15Present := ab15ao15Input["arguments"]; ab15ao15Present {
		ab15ao15WantKeys = 7
	}
	if len(ab15ao15Output) != ab15ao15WantKeys {
		ab15ao15T.Fatalf("unexpected output keys: %#v", ab15ao15Output)
	}
	_, ab15ao15InputArguments := ab15ao15Normalized["arguments"]
	_, ab15ao15OutputArguments := ab15ao15Output["arguments"]
	if ab15ao15InputArguments != ab15ao15OutputArguments {
		ab15ao15T.Fatal("argument presence changed")
	}

	for _, ab15ao15Key := range []string{"context", "arguments"} {
		if !ab15ao15reflect.DeepEqual(ab15ao15Output[ab15ao15Key], ab15ao15Normalized[ab15ao15Key]) {
			ab15ao15T.Fatalf("preservation %s: %#v", ab15ao15Key, ab15ao15Output)
		}
	}
	if !ab15ao15reflect.DeepEqual(ab15ao15Output["raw"], ab15ao15Normalized["fields"]) || ab15ao15Output["physical"] != "Unknown" || ab15ao15Output["qualified"] != false {
		ab15ao15T.Fatalf("raw/qualification: %#v", ab15ao15Output)
	}
	if !ab15ao15reflect.DeepEqual(ab15ao15Output, ab15ao15Decode(ab15ao15T, ab15ao15Encode(ab15ao15T, ab15ao15Output))) {
		ab15ao15T.Fatal("round trip")
	}
	return ab15ao15Output
}

func ab15ao15Expected() map[string]any {
	return map[string]any{
		"secondary_heating_pump": map[string]any{"raw": []any{float64(11)}, "state": "Hypothesis", "flags": map[string]any{"origin": "origin-Unknown"}},
		"circulation_pump":       map[string]any{"raw": []any{float64(22)}, "state": "Hypothesis", "flags": map[string]any{"origin": "origin-Unknown"}},
		"solar_pump":             map[string]any{"raw": []any{float64(33)}, "state": "Hypothesis", "flags": map[string]any{"origin": "origin-Unknown"}},
		"error_relay":            map[string]any{"raw": []any{float64(15)}, "state": "Hypothesis", "label": "state-0F", "flags": map[string]any{"origin": "origin-Unknown"}},
		"loading_pump":           map[string]any{"raw": []any{float64(44)}, "state": "Hypothesis", "flags": map[string]any{"origin": "origin-Unknown"}},
		"fume_flap":              map[string]any{"raw": []any{float64(240)}, "state": "Hypothesis", "label": "state-F0", "flags": map[string]any{"origin": "origin-Unknown"}},
	}
}

func TestAb15Ao15IndependentFields(ab15ao15T *ab15ao15testing.T) {
	for _, ab15ao15Name := range []string{"secondary_heating_pump", "circulation_pump", "solar_pump", "error_relay", "loading_pump", "fume_flap"} {
		for _, ab15ao15Direction := range []string{"request", "positive"} {
			for _, ab15ao15Case := range []struct {
				ab15ao15Byte                                             int
				ab15ao15PumpRequest, ab15ao15PumpPositive, ab15ao15Relay string
			}{
				{0, "", "", "no-action-candidate"}, {1, "", "", "no-action-candidate"}, {15, "", "", "state-0F"}, {240, "", "", "state-F0"}, {255, "leave-unchanged", "function-unavailable", "no-action-candidate"}, {42, "", "", "no-action-candidate"},
			} {
				ab15ao15T.Run(ab15ao15Name+"/"+ab15ao15Direction+"/"+string(ab15ao15Encode(ab15ao15T, ab15ao15Case.ab15ao15Byte)), func(ab15ao15T *ab15ao15testing.T) {
					ab15ao15Input := ab15ao15Control()
					ab15ao15Input["context"].(map[string]any)["direction"] = ab15ao15Direction
					ab15ao15Input["fields"].(map[string]any)[ab15ao15Name] = []int{ab15ao15Case.ab15ao15Byte}
					ab15ao15Want := ab15ao15Expected()
					ab15ao15Field := ab15ao15Want[ab15ao15Name].(map[string]any)
					ab15ao15Field["raw"] = []any{float64(ab15ao15Case.ab15ao15Byte)}
					ab15ao15Label := ab15ao15Case.ab15ao15PumpRequest
					if ab15ao15Direction == "positive" {
						ab15ao15Label = ab15ao15Case.ab15ao15PumpPositive
					}
					if ab15ao15Name == "error_relay" || ab15ao15Name == "fume_flap" {
						ab15ao15Label = ab15ao15Case.ab15ao15Relay
					}
					delete(ab15ao15Field, "label")
					if ab15ao15Label != "" {
						ab15ao15Field["label"] = ab15ao15Label
					}
					ab15ao15Output := ab15ao15Run(ab15ao15T, ab15ao15Input)
					if ab15ao15Output["state"] != "Hypothesis" || !ab15ao15reflect.DeepEqual(ab15ao15Output["fields"], ab15ao15Want) {
						ab15ao15T.Fatalf("exact independent fields: got %#v want %#v", ab15ao15Output, ab15ao15Want)
					}
				})
			}
		}
	}
}

func TestAb15Ao15Origin(ab15ao15T *ab15ao15testing.T) {
	for _, ab15ao15Case := range []struct {
		ab15ao15ID, ab15ao15Origin string
		ab15ao15Arguments          any
	}{
		{"ab15ao15Absent", "origin-Unknown", nil}, {"ab15ao15Empty", "origin-Unknown", map[string]any{}}, {"ab15ao15False", "fallback-candidate", map[string]any{"actual_signal": false}}, {"ab15ao15True", "actual-input-candidate", map[string]any{"actual_signal": true}},
	} {
		for _, ab15ao15Feedback := range []int{17, 93} {
			for _, ab15ao15Direction := range []string{"request", "positive"} {
				ab15ao15T.Run(ab15ao15Case.ab15ao15ID+"/"+ab15ao15Direction+"/"+string(ab15ao15Encode(ab15ao15T, ab15ao15Feedback)), func(ab15ao15T *ab15ao15testing.T) {
					ab15ao15Input := ab15ao15Control()
					ab15ao15Input["context"].(map[string]any)["feedback"] = []int{ab15ao15Feedback}
					ab15ao15Input["context"].(map[string]any)["direction"] = ab15ao15Direction
					if ab15ao15Case.ab15ao15Arguments != nil {
						ab15ao15Input["arguments"] = ab15ao15Case.ab15ao15Arguments
					}
					ab15ao15Want := ab15ao15Expected()
					for _, ab15ao15Field := range ab15ao15Want {
						ab15ao15Field.(map[string]any)["flags"] = map[string]any{"origin": ab15ao15Case.ab15ao15Origin}
					}
					ab15ao15Output := ab15ao15Run(ab15ao15T, ab15ao15Input)
					if !ab15ao15reflect.DeepEqual(ab15ao15Output["fields"], ab15ao15Want) {
						ab15ao15T.Fatalf("origin/default/polarity: %#v", ab15ao15Output)
					}
				})
			}
		}
	}
}

func TestAb15Ao15Partial(ab15ao15T *ab15ao15testing.T) {
	for ab15ao15Name := range ab15ao15Expected() {
		for _, ab15ao15Width := range []string{"missing", "empty", "overlong", "zero"} {
			ab15ao15T.Run(ab15ao15Name+"/"+ab15ao15Width, func(ab15ao15T *ab15ao15testing.T) {
				ab15ao15Input := ab15ao15Control()
				ab15ao15Raw := ab15ao15Input["fields"].(map[string]any)
				ab15ao15Want := ab15ao15Expected()
				ab15ao15State := "partial"
				ab15ao15Field := map[string]any{"raw": nil, "state": "Unknown"}
				switch ab15ao15Width {
				case "missing":
					delete(ab15ao15Raw, ab15ao15Name)
				case "empty":
					ab15ao15Raw[ab15ao15Name] = []int{}
					ab15ao15Field["raw"] = []any{}
				case "overlong":
					ab15ao15Raw[ab15ao15Name] = []int{255, 15}
					ab15ao15Field["raw"] = []any{float64(255), float64(15)}
				case "zero":
					ab15ao15Raw[ab15ao15Name] = []int{0}
					ab15ao15State = "Hypothesis"
					ab15ao15Field = map[string]any{"raw": []any{float64(0)}, "state": "Hypothesis", "flags": map[string]any{"origin": "origin-Unknown"}}
					if ab15ao15Name == "error_relay" || ab15ao15Name == "fume_flap" {
						ab15ao15Field["label"] = "no-action-candidate"
					}
				}
				ab15ao15Want[ab15ao15Name] = ab15ao15Field
				ab15ao15Output := ab15ao15Run(ab15ao15T, ab15ao15Input)
				if ab15ao15Output["state"] != ab15ao15State || !ab15ao15reflect.DeepEqual(ab15ao15Output["fields"], ab15ao15Want) {
					ab15ao15T.Fatalf("partial siblings: %#v", ab15ao15Output)
				}
			})
		}
	}
}

func TestAb15Ao15SelectorsAndOpaqueContext(ab15ao15T *ab15ao15testing.T) {
	for ab15ao15Key, ab15ao15Values := range map[string][]any{
		"primary": {0, 182}, "secondary": {0, 19}, "operation": {0, 2}, "direction": {"", "unknown"}, "view": {"", "unknown"},
	} {
		for _, ab15ao15Value := range append(ab15ao15Values, nil) {
			ab15ao15T.Run(ab15ao15Key+"/"+string(ab15ao15Encode(ab15ao15T, ab15ao15Value)), func(ab15ao15T *ab15ao15testing.T) {
				ab15ao15Input := ab15ao15Control()
				if ab15ao15Value == nil {
					delete(ab15ao15Input["context"].(map[string]any), ab15ao15Key)
				} else {
					ab15ao15Input["context"].(map[string]any)[ab15ao15Key] = ab15ao15Value
				}
				ab15ao15Output := ab15ao15Run(ab15ao15T, ab15ao15Input)
				if ab15ao15Output["state"] != "Unknown" || len(ab15ao15Output["fields"].(map[string]any)) != 0 {
					ab15ao15T.Fatalf("guessed selection: %#v", ab15ao15Output)
				}
			})
		}
	}
	for _, ab15ao15Key := range []string{"family", "variant", "setpoint", "feedback", "opaque"} {
		for _, ab15ao15Value := range []any{nil, 0, "directional-fields", []any{}, map[string]any{"nested": []any{false, 1}}} {
			ab15ao15Input := ab15ao15Control()
			ab15ao15Input["context"].(map[string]any)[ab15ao15Key] = ab15ao15Value
			ab15ao15Output := ab15ao15Run(ab15ao15T, ab15ao15Input)
			if !ab15ao15reflect.DeepEqual(ab15ao15Output["fields"], ab15ao15Expected()) {
				ab15ao15T.Fatalf("opaque context affected selection: %#v", ab15ao15Output)
			}
		}
	}
}

func TestAb15Ao15Reject(ab15ao15T *ab15ao15testing.T) {
	for _, ab15ao15Bad := range []string{"", "null", "[]", "{}", "{", `{"context":{},"context":{},"fields":{}}`, `{"context":{"family":{"x":1,"x":2}},"fields":{}}`, `{"context":{"family":[{"x":1,"\u0078":2}]},"fields":{}}`, `{"context":{},"fields":{"x":[],"x":[]}}`, `{"context":{},"fields":{},"arguments":{"actual_signal":true,"actual_signal":false}}`} {
		ab15ao15Result, ab15ao15Err := Ab15Ao15Interpret([]byte(ab15ao15Bad))
		if ab15ao15Err == nil || ab15ao15Result != nil {
			ab15ao15T.Fatalf("accepted %q", ab15ao15Bad)
		}
	}
	for ab15ao15Key, ab15ao15Values := range map[string][]any{
		"primary": {nil, "181", 1.5, true, []int{}, map[string]any{}}, "secondary": {nil, "18", false}, "operation": {nil, "1", 1.5}, "direction": {nil, 1, []int{}}, "view": {nil, true, map[string]any{}},
	} {
		for _, ab15ao15Value := range ab15ao15Values {
			ab15ao15Input := ab15ao15Control()
			ab15ao15Input["context"].(map[string]any)[ab15ao15Key] = ab15ao15Value
			ab15ao15Result, ab15ao15Err := Ab15Ao15Interpret(ab15ao15Encode(ab15ao15T, ab15ao15Input))
			if ab15ao15Err == nil || ab15ao15Result != nil {
				ab15ao15T.Fatalf("accepted identity %s %#v", ab15ao15Key, ab15ao15Value)
			}
		}
	}
	for _, ab15ao15Value := range []any{nil, true, 0, "", map[string]any{}, []any{-1}, []any{256}, []any{1.5}, []any{nil}, []any{true}, []any{"1"}, []any{[]int{1}}, []any{map[string]any{}}} {
		for _, ab15ao15Name := range []string{"secondary_heating_pump", "circulation_pump", "solar_pump", "error_relay", "loading_pump", "fume_flap", "unmodelled"} {
			ab15ao15Input := ab15ao15Control()
			ab15ao15Input["fields"].(map[string]any)[ab15ao15Name] = ab15ao15Value
			ab15ao15Result, ab15ao15Err := Ab15Ao15Interpret(ab15ao15Encode(ab15ao15T, ab15ao15Input))
			if ab15ao15Err == nil || ab15ao15Result != nil {
				ab15ao15T.Fatalf("accepted bytes %#v", ab15ao15Value)
			}
		}
	}
	for _, ab15ao15Key := range []string{"context", "fields", "arguments"} {
		for _, ab15ao15Value := range []any{nil, []int{}, 0, "", true} {
			ab15ao15Input := ab15ao15Control()
			ab15ao15Input[ab15ao15Key] = ab15ao15Value
			ab15ao15Result, ab15ao15Err := Ab15Ao15Interpret(ab15ao15Encode(ab15ao15T, ab15ao15Input))
			if ab15ao15Err == nil || ab15ao15Result != nil {
				ab15ao15T.Fatalf("accepted kind %s", ab15ao15Key)
			}
		}
	}
	for _, ab15ao15Value := range []any{nil, 0, 1, "true", []int{}, map[string]any{}} {
		ab15ao15Input := ab15ao15Control()
		ab15ao15Input["arguments"] = map[string]any{"actual_signal": ab15ao15Value}
		ab15ao15Result, ab15ao15Err := Ab15Ao15Interpret(ab15ao15Encode(ab15ao15T, ab15ao15Input))
		if ab15ao15Err == nil || ab15ao15Result != nil {
			ab15ao15T.Fatalf("accepted flag %#v", ab15ao15Value)
		}
	}
	for _, ab15ao15Bad := range []string{
		string(ab15ao15Encode(ab15ao15T, ab15ao15Control())) + ` {}`, ab15ao15strings.Replace(string(ab15ao15Encode(ab15ao15T, ab15ao15Control())), `[11]`, `[1e0]`, 1),
	} {
		ab15ao15Result, ab15ao15Err := Ab15Ao15Interpret([]byte(ab15ao15Bad))
		if ab15ao15Err == nil || ab15ao15Result != nil {
			ab15ao15T.Fatalf("accepted malformed/integer syntax %s", ab15ao15Bad)
		}
	}
	ab15ao15Input := ab15ao15Control()
	ab15ao15Input["arguments"] = map[string]any{"default": true}
	ab15ao15Result, ab15ao15Err := Ab15Ao15Interpret(ab15ao15Encode(ab15ao15T, ab15ao15Input))
	if ab15ao15Err == nil || ab15ao15Result != nil {
		ab15ao15T.Fatal("accepted invented default")
	}
}

func TestAb15Ao15Isolation(ab15ao15T *ab15ao15testing.T) {
	ab15ao15Input := ab15ao15Control()
	ab15ao15Input["arguments"] = map[string]any{"actual_signal": false}
	ab15ao15Raw := ab15ao15Encode(ab15ao15T, ab15ao15Input)
	ab15ao15First, ab15ao15Err := Ab15Ao15Interpret(ab15ao15Raw)
	if ab15ao15Err != nil {
		ab15ao15T.Fatal(ab15ao15Err)
	}
	ab15ao15Saved := append([]byte(nil), ab15ao15First...)
	for ab15ao15Index := range ab15ao15Raw {
		ab15ao15Raw[ab15ao15Index] = ' '
	}
	if string(ab15ao15First) != string(ab15ao15Saved) {
		ab15ao15T.Fatal("input aliases output")
	}
	ab15ao15Nested := ab15ao15Decode(ab15ao15T, ab15ao15First)
	ab15ao15Nested["raw"].(map[string]any)["solar_pump"].([]any)[0] = float64(99)
	if ab15ao15Nested["fields"].(map[string]any)["solar_pump"].(map[string]any)["raw"].([]any)[0] != float64(33) {
		ab15ao15T.Fatal("raw aliases semantic raw")
	}
	ab15ao15Nested["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "changed"
	ab15ao15Nested["arguments"].(map[string]any)["actual_signal"] = true
	ab15ao15Nested["fields"].(map[string]any)["circulation_pump"].(map[string]any)["flags"].(map[string]any)["origin"] = "changed"
	ab15ao15Second, ab15ao15Err := Ab15Ao15Interpret(ab15ao15Encode(ab15ao15T, ab15ao15Input))
	if ab15ao15Err != nil {
		ab15ao15T.Fatal(ab15ao15Err)
	}
	if string(ab15ao15Second) != string(ab15ao15Saved) {
		ab15ao15T.Fatal("nested mutation affected later call")
	}
	for ab15ao15Index := range ab15ao15First {
		ab15ao15First[ab15ao15Index] = ' '
	}
	if string(ab15ao15Second) != string(ab15ao15Saved) {
		ab15ao15T.Fatal("output aliases another call")
	}
}

func TestAb15Ao15Fixtures(ab15ao15T *ab15ao15testing.T) {
	ab15ao15Data, ab15ao15Err := ab15ao15os.ReadFile("testdata/ab15ao15_b512_feedback_views.json")
	if ab15ao15Err != nil {
		ab15ao15T.Fatal(ab15ao15Err)
	}
	var ab15ao15Cases []struct {
		Ab15Ao15ID     string            `json:"id"`
		Ab15Ao15Input  map[string]any    `json:"input"`
		Ab15Ao15State  string            `json:"state"`
		Ab15Ao15Labels map[string]string `json:"labels"`
	}
	if ab15ao15Err := ab15ao15json.Unmarshal(ab15ao15Data, &ab15ao15Cases); ab15ao15Err != nil {
		ab15ao15T.Fatal(ab15ao15Err)
	}
	if len(ab15ao15Cases) != 3 {
		ab15ao15T.Fatal("fixture count")
	}
	for _, ab15ao15Case := range ab15ao15Cases {
		ab15ao15T.Run(ab15ao15Case.Ab15Ao15ID, func(ab15ao15T *ab15ao15testing.T) {
			if !ab15ao15strings.HasPrefix(ab15ao15Case.Ab15Ao15ID, "ab15ao15") {
				ab15ao15T.Fatal("fixture ID")
			}
			ab15ao15Output := ab15ao15Run(ab15ao15T, ab15ao15Case.Ab15Ao15Input)
			if ab15ao15Output["state"] != ab15ao15Case.Ab15Ao15State {
				ab15ao15T.Fatalf("fixture state %#v", ab15ao15Output)
			}
			for ab15ao15Name, ab15ao15Label := range ab15ao15Case.Ab15Ao15Labels {
				if ab15ao15Output["fields"].(map[string]any)[ab15ao15Name].(map[string]any)["label"] != ab15ao15Label {
					ab15ao15T.Fatalf("fixture label %s", ab15ao15Name)
				}
			}
		})
	}
}
