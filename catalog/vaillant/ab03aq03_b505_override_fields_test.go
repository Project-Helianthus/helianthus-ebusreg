package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

type ab03aq03Case struct {
	Ab03Aq03ID       string          `json:"id"`
	Ab03Aq03Input    json.RawMessage `json:"input"`
	Ab03Aq03Expected map[string]any  `json:"expected"`
}

func ab03aq03Cases(ab03aq03T *testing.T) []ab03aq03Case {
	ab03aq03T.Helper()
	ab03aq03Data, ab03aq03Err := os.ReadFile("testdata/ab03aq03_b505_override_fields.json")
	if ab03aq03Err != nil {
		ab03aq03T.Fatal(ab03aq03Err)
	}
	var ab03aq03Cases []ab03aq03Case
	if ab03aq03Err := json.Unmarshal(ab03aq03Data, &ab03aq03Cases); ab03aq03Err != nil {
		ab03aq03T.Fatal(ab03aq03Err)
	}
	if len(ab03aq03Cases) != 7 {
		ab03aq03T.Fatal("fixture case count")
	}
	return ab03aq03Cases
}

func ab03aq03Subset(ab03aq03T *testing.T, ab03aq03Expected, ab03aq03Actual any) {
	ab03aq03T.Helper()
	if ab03aq03Map, ab03aq03OK := ab03aq03Expected.(map[string]any); ab03aq03OK {
		ab03aq03Got, ab03aq03OK := ab03aq03Actual.(map[string]any)
		if !ab03aq03OK {
			ab03aq03T.Fatalf("expected object, got %v", ab03aq03Actual)
		}
		for ab03aq03Key, ab03aq03Value := range ab03aq03Map {
			ab03aq03Subset(ab03aq03T, ab03aq03Value, ab03aq03Got[ab03aq03Key])
		}
		return
	}
	if ab03aq03Number, ab03aq03OK := ab03aq03Expected.(int); ab03aq03OK {
		ab03aq03Expected = float64(ab03aq03Number)
	}
	if !reflect.DeepEqual(ab03aq03Expected, ab03aq03Actual) {
		ab03aq03T.Fatalf("want %#v got %#v", ab03aq03Expected, ab03aq03Actual)
	}
}

func ab03aq03Run(ab03aq03T *testing.T, ab03aq03Input any) map[string]any {
	ab03aq03T.Helper()
	ab03aq03Raw, ab03aq03Err := json.Marshal(ab03aq03Input)
	if ab03aq03Err != nil {
		ab03aq03T.Fatal(ab03aq03Err)
	}
	ab03aq03Output, ab03aq03Err := Ab03Aq03Evaluate(ab03aq03Raw)
	if ab03aq03Err != nil {
		ab03aq03T.Fatal(ab03aq03Err)
	}
	var ab03aq03Got map[string]any
	if ab03aq03Err := json.Unmarshal(ab03aq03Output, &ab03aq03Got); ab03aq03Err != nil {
		ab03aq03T.Fatal(ab03aq03Err)
	}
	return ab03aq03Got
}

func ab03aq03Input(ab03aq03T *testing.T, ab03aq03Index int) map[string]any {
	ab03aq03T.Helper()
	var ab03aq03Input map[string]any
	if ab03aq03Err := json.Unmarshal(ab03aq03Cases(ab03aq03T)[ab03aq03Index].Ab03Aq03Input, &ab03aq03Input); ab03aq03Err != nil {
		ab03aq03T.Fatal(ab03aq03Err)
	}
	return ab03aq03Input
}

func TestAb03Aq03Fixture(ab03aq03T *testing.T) {
	for _, ab03aq03Case := range ab03aq03Cases(ab03aq03T) {
		ab03aq03T.Run(ab03aq03Case.Ab03Aq03ID, func(ab03aq03T *testing.T) {
			ab03aq03Got := ab03aq03Run(ab03aq03T, ab03aq03Case.Ab03Aq03Input)
			ab03aq03Subset(ab03aq03T, ab03aq03Case.Ab03Aq03Expected, ab03aq03Got)
			var ab03aq03Original map[string]any
			if ab03aq03Err := json.Unmarshal(ab03aq03Case.Ab03Aq03Input, &ab03aq03Original); ab03aq03Err != nil {
				ab03aq03T.Fatal(ab03aq03Err)
			}
			for ab03aq03Key, ab03aq03Value := range ab03aq03Original {
				if !reflect.DeepEqual(ab03aq03Value, ab03aq03Got[ab03aq03Key]) {
					ab03aq03T.Fatalf("input retention: %s", ab03aq03Key)
				}
			}
			if ab03aq03Original["context"].(map[string]any)["operation"] == float64(78) {
				if _, ab03aq03OK := ab03aq03Got["candidate_validity_seconds"]; ab03aq03OK {
					ab03aq03T.Fatal("priority inherited duration")
				}
			}
			if ab03aq03Original["context"].(map[string]any)["operation"] == float64(97) {
				if _, ab03aq03OK := ab03aq03Got["results"].(map[string]any)["sensor_status"].(map[string]any)["label"]; ab03aq03OK {
					ab03aq03T.Fatal("OP61 sensor label")
				}
			}
		})
	}
}

func TestAb03Aq03Identity(ab03aq03T *testing.T) {
	for _, ab03aq03Case := range []struct {
		ab03aq03Key   string
		ab03aq03Value any
	}{{"primary", 180}, {"secondary", 6}, {"operation", nil}, {"operation", 0}, {"operation", 50}, {"direction", "response"}, {"view", "other"}, {"family", "other-family"}, {"variant", "other-variant"}} {
		ab03aq03T.Run(ab03aq03Case.ab03aq03Key+"/"+strings.TrimSpace(string(ab03aq03MustJSON(ab03aq03Case.ab03aq03Value))), func(ab03aq03T *testing.T) {
			ab03aq03Input := ab03aq03Input(ab03aq03T, 0)
			ab03aq03Context := ab03aq03Input["context"].(map[string]any)
			if ab03aq03Case.ab03aq03Value == nil {
				delete(ab03aq03Context, ab03aq03Case.ab03aq03Key)
			} else {
				ab03aq03Context[ab03aq03Case.ab03aq03Key] = ab03aq03Case.ab03aq03Value
			}
			ab03aq03Got := ab03aq03Run(ab03aq03T, ab03aq03Input)
			ab03aq03Subset(ab03aq03T, map[string]any{"context": ab03aq03Context, "fields": ab03aq03Input["fields"], "qualified": false, "physical": "Unknown"}, ab03aq03Got)
			if ab03aq03Case.ab03aq03Key == "family" || ab03aq03Case.ab03aq03Key == "variant" {
				if ab03aq03Got["grade"] != "Hypothesis" {
					ab03aq03T.Fatal("opaque identity changed grade")
				}
				return
			}
			if ab03aq03Got["state"] != "Unknown" {
				ab03aq03T.Fatal("unsupported identity admitted")
			}
			for _, ab03aq03Result := range ab03aq03Got["results"].(map[string]any) {
				if len(ab03aq03Result.(map[string]any)) != 2 {
					ab03aq03T.Fatal("unsupported conversion")
				}
			}
			if _, ab03aq03OK := ab03aq03Got["candidate_validity_seconds"]; ab03aq03OK {
				ab03aq03T.Fatal("unsupported duration")
			}
		})
	}
}

func ab03aq03MustJSON(ab03aq03Value any) []byte {
	ab03aq03Raw, _ := json.Marshal(ab03aq03Value)
	return ab03aq03Raw
}

func TestAb03Aq03WidthsAndArguments(ab03aq03T *testing.T) {
	for _, ab03aq03Key := range []string{"sensor_status", "remote_control", "priority"} {
		for ab03aq03Index, ab03aq03Value := range []any{nil, []any{}, []any{0, 1}, []any{0}} {
			ab03aq03Input := ab03aq03Input(ab03aq03T, 0)
			if ab03aq03Key == "priority" {
				ab03aq03Input = ab03aq03InputFromPriority(ab03aq03T)
			}
			ab03aq03Fields := ab03aq03Input["fields"].(map[string]any)
			if ab03aq03Value == nil {
				delete(ab03aq03Fields, ab03aq03Key)
			} else {
				ab03aq03Fields[ab03aq03Key] = ab03aq03Value
			}
			ab03aq03Got := ab03aq03Run(ab03aq03T, ab03aq03Input)
			ab03aq03Result := ab03aq03Got["results"].(map[string]any)[ab03aq03Key].(map[string]any)
			if ab03aq03Index < 3 {
				if ab03aq03Got["state"] != "partial" || ab03aq03Result["state"] != "incomplete" {
					ab03aq03T.Fatal("width admitted")
				}
				if _, ab03aq03OK := ab03aq03Result["number"]; ab03aq03OK {
					ab03aq03T.Fatal("incomplete number")
				}
				if ab03aq03Value == nil {
					if _, ab03aq03OK := ab03aq03Result["raw"]; ab03aq03OK {
						ab03aq03T.Fatal("missing became present")
					}
				}
			} else if ab03aq03Got["state"] != "complete" {
				ab03aq03T.Fatal("zero incomplete")
			}
			if ab03aq03Key == "sensor_status" {
				ab03aq03Subset(ab03aq03T, map[string]any{"label": "enabled", "state": "complete"}, ab03aq03Got["results"].(map[string]any)["remote_control"])
			}
			if ab03aq03Key == "remote_control" {
				ab03aq03Subset(ab03aq03T, map[string]any{"label": "short-circuit", "state": "complete"}, ab03aq03Got["results"].(map[string]any)["sensor_status"])
			}
		}
	}
	for _, ab03aq03Value := range []int{0, 25, 50, 75, 100, 26, 255} {
		ab03aq03Input := ab03aq03Input(ab03aq03T, 0)
		ab03aq03Input["arguments"].(map[string]any)["display_percent"] = ab03aq03Value
		ab03aq03Got := ab03aq03Run(ab03aq03T, ab03aq03Input)["results"].(map[string]any)["display_percent"].(map[string]any)
		ab03aq03Want := ab03aq03Value == 0 || ab03aq03Value == 25 || ab03aq03Value == 50 || ab03aq03Value == 75 || ab03aq03Value == 100
		if ab03aq03Got["number"] != float64(ab03aq03Value) || ab03aq03Got["candidate_domain"] != ab03aq03Want {
			ab03aq03T.Fatal("display rounded or domain changed")
		}
	}
}

func ab03aq03InputFromPriority(ab03aq03T *testing.T) map[string]any {
	return ab03aq03Input(ab03aq03T, 4)
}

func TestAb03Aq03Reject(ab03aq03T *testing.T) {
	for _, ab03aq03Raw := range []string{``, `null`, `[]`, `{`, `{} {}`, `{"context":{},"context":{}}`, `{"context":{"x":0,"x":1},"fields":{}}`, `{"context":{},"fields":{},"x":[{"y":1,"y":2}]}`, `{"context":{},"fields":{},"x":NaN}`, `{"context":{},"fields":{},"x":1e999}`, `{"context":{},"fields":{},"arguments":{"x":1,"x":2}}`, `{"context":{},"fields":{"x":[0],"x":[1]}}`} {
		ab03aq03Output, ab03aq03Err := Ab03Aq03Evaluate([]byte(ab03aq03Raw))
		if ab03aq03Err == nil || ab03aq03Output != nil {
			ab03aq03T.Fatalf("accepted malformed %s", ab03aq03Raw)
		}
	}
	for _, ab03aq03Path := range []string{"primary", "secondary", "operation", "display_percent", "sensor_status", "remote_control", "priority"} {
		for _, ab03aq03Value := range []any{nil, true, "1", -1, 256, 1.5} {
			ab03aq03Input := ab03aq03Input(ab03aq03T, 0)
			if ab03aq03Path == "priority" {
				ab03aq03Input = ab03aq03InputFromPriority(ab03aq03T)
			}
			switch ab03aq03Path {
			case "display_percent":
				ab03aq03Input["arguments"].(map[string]any)[ab03aq03Path] = ab03aq03Value
			case "sensor_status", "remote_control", "priority":
				ab03aq03Input["fields"].(map[string]any)[ab03aq03Path] = []any{ab03aq03Value}
			default:
				ab03aq03Input["context"].(map[string]any)[ab03aq03Path] = ab03aq03Value
			}
			ab03aq03Output, ab03aq03Err := Ab03Aq03Evaluate(ab03aq03MustJSON(ab03aq03Input))
			if ab03aq03Err == nil || ab03aq03Output != nil {
				ab03aq03T.Fatalf("accepted %s=%v", ab03aq03Path, ab03aq03Value)
			}
		}
	}
	for _, ab03aq03Key := range []string{"context", "fields", "arguments"} {
		for _, ab03aq03Value := range []any{nil, true, 1, "x", []any{}} {
			ab03aq03Input := ab03aq03Input(ab03aq03T, 0)
			ab03aq03Input[ab03aq03Key] = ab03aq03Value
			if ab03aq03Output, ab03aq03Err := Ab03Aq03Evaluate(ab03aq03MustJSON(ab03aq03Input)); ab03aq03Err == nil || ab03aq03Output != nil {
				ab03aq03T.Fatal("object kind accepted")
			}
		}
	}
	ab03aq03Deep := `{"context":{},"fields":{},"x":` + strings.Repeat("[", 3000000) + strings.Repeat("]", 3000000) + `}`
	if ab03aq03Output, ab03aq03Err := Ab03Aq03Evaluate([]byte(ab03aq03Deep)); ab03aq03Err == nil || ab03aq03Output != nil {
		ab03aq03T.Fatal("deep nesting accepted")
	}
	if ab03aq03Output, ab03aq03Err := Ab03Aq03Evaluate([]byte(`{"context":{},"fields":{},"arguments":{"display_percent":-0}}`)); ab03aq03Err == nil || ab03aq03Output != nil {
		ab03aq03T.Fatal("negative zero accepted")
	}
	for _, ab03aq03Key := range []string{"physical", "qualified", "state", "results", "grade", "candidate_validity_seconds"} {
		ab03aq03Input := ab03aq03Input(ab03aq03T, 0)
		ab03aq03Input[ab03aq03Key] = "x"
		if ab03aq03Output, ab03aq03Err := Ab03Aq03Evaluate(ab03aq03MustJSON(ab03aq03Input)); ab03aq03Err == nil || ab03aq03Output != nil {
			ab03aq03T.Fatalf("colliding key %s accepted", ab03aq03Key)
		}
	}
}

func TestAb03Aq03OptionalDisplay(ab03aq03T *testing.T) {
	ab03aq03Input := ab03aq03Input(ab03aq03T, 0)
	delete(ab03aq03Input, "arguments")
	ab03aq03Got := ab03aq03Run(ab03aq03T, ab03aq03Input)
	ab03aq03Results := ab03aq03Got["results"].(map[string]any)
	if _, ab03aq03OK := ab03aq03Results["display_percent"]; ab03aq03OK {
		ab03aq03T.Fatal("display result invented")
	}
	if ab03aq03Got["state"] != "complete" || ab03aq03Results["sensor_status"].(map[string]any)["label"] == nil {
		ab03aq03T.Fatal("siblings unavailable without display_percent")
	}
}

func TestAb03Aq03Memory(ab03aq03T *testing.T) {
	ab03aq03Input := append(json.RawMessage(nil), ab03aq03Cases(ab03aq03T)[0].Ab03Aq03Input...)
	ab03aq03Output, ab03aq03Err := Ab03Aq03Evaluate(ab03aq03Input)
	if ab03aq03Err != nil {
		ab03aq03T.Fatal(ab03aq03Err)
	}
	ab03aq03Control := append([]byte(nil), ab03aq03Output...)
	for ab03aq03Index := range ab03aq03Input {
		ab03aq03Input[ab03aq03Index] = ' '
	}
	if !reflect.DeepEqual([]byte(ab03aq03Output), ab03aq03Control) {
		ab03aq03T.Fatal("input aliases output")
	}
	for ab03aq03Index := range ab03aq03Output {
		ab03aq03Output[ab03aq03Index] = ' '
	}
	ab03aq03Next, ab03aq03Err := Ab03Aq03Evaluate(ab03aq03Cases(ab03aq03T)[0].Ab03Aq03Input)
	if ab03aq03Err != nil || !reflect.DeepEqual([]byte(ab03aq03Next), ab03aq03Control) {
		ab03aq03T.Fatal("call shared storage")
	}
}

func TestAb03Aq03OpaqueAndKinds(ab03aq03T *testing.T) {
	for _, ab03aq03Value := range []any{nil, []any{}, []any{0}, []any{85, 170}} {
		ab03aq03Input := ab03aq03Input(ab03aq03T, 3)
		if ab03aq03Value == nil {
			delete(ab03aq03Input["fields"].(map[string]any), "dhw_setpoint")
		} else {
			ab03aq03Input["fields"].(map[string]any)["dhw_setpoint"] = ab03aq03Value
		}
		ab03aq03Got := ab03aq03Run(ab03aq03T, ab03aq03Input)
		var ab03aq03Fields any
		if ab03aq03Err := json.Unmarshal(ab03aq03MustJSON(ab03aq03Input["fields"]), &ab03aq03Fields); ab03aq03Err != nil {
			ab03aq03T.Fatal(ab03aq03Err)
		}
		if !reflect.DeepEqual(ab03aq03Fields, ab03aq03Got["fields"]) {
			ab03aq03T.Fatal("opaque raw changed")
		}
		if ab03aq03Value != nil {
			ab03aq03Result := ab03aq03Got["results"].(map[string]any)["dhw_setpoint"].(map[string]any)
			if len(ab03aq03Result) != 2 || ab03aq03Result["state"] != "Unknown" {
				ab03aq03T.Fatal("opaque conversion")
			}
		}
	}
	for _, ab03aq03Key := range []string{"direction", "view", "family", "variant"} {
		for _, ab03aq03Value := range []any{nil, false, 1, []any{}, map[string]any{}} {
			ab03aq03Input := ab03aq03Input(ab03aq03T, 0)
			ab03aq03Input["context"].(map[string]any)[ab03aq03Key] = ab03aq03Value
			if ab03aq03Output, ab03aq03Err := Ab03Aq03Evaluate(ab03aq03MustJSON(ab03aq03Input)); ab03aq03Err == nil || ab03aq03Output != nil {
				ab03aq03T.Fatal("identifier kind admitted")
			}
		}
	}
	for _, ab03aq03Key := range []string{"sensor_status", "remote_control", "priority", "dhw_setpoint"} {
		for _, ab03aq03Value := range []any{nil, false, 0, "x", map[string]any{}} {
			ab03aq03Input := ab03aq03Input(ab03aq03T, 0)
			ab03aq03Input["fields"].(map[string]any)[ab03aq03Key] = ab03aq03Value
			if ab03aq03Output, ab03aq03Err := Ab03Aq03Evaluate(ab03aq03MustJSON(ab03aq03Input)); ab03aq03Err == nil || ab03aq03Output != nil {
				ab03aq03T.Fatal("field kind admitted")
			}
		}
	}
	for _, ab03aq03Index := range []int{3, 4} {
		ab03aq03Input := ab03aq03Input(ab03aq03T, ab03aq03Index)
		ab03aq03Input["arguments"] = map[string]any{"display_percent": 25}
		if ab03aq03Output, ab03aq03Err := Ab03Aq03Evaluate(ab03aq03MustJSON(ab03aq03Input)); ab03aq03Err == nil || ab03aq03Output != nil {
			ab03aq03T.Fatal("unmodeled argument admitted")
		}
	}
}

func TestAb03Aq03PresenceAndIsolation(ab03aq03T *testing.T) {
	for _, ab03aq03Key := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
		for _, ab03aq03Present := range []bool{false, true} {
			ab03aq03Input := ab03aq03Input(ab03aq03T, 0)
			ab03aq03Context := ab03aq03Input["context"].(map[string]any)
			if !ab03aq03Present {
				delete(ab03aq03Context, ab03aq03Key)
			} else if ab03aq03Key == "primary" || ab03aq03Key == "secondary" || ab03aq03Key == "operation" {
				ab03aq03Context[ab03aq03Key] = float64(0)
			} else {
				ab03aq03Context[ab03aq03Key] = ""
			}
			ab03aq03Got := ab03aq03Run(ab03aq03T, ab03aq03Input)
			if !reflect.DeepEqual(ab03aq03Context, ab03aq03Got["context"]) {
				ab03aq03T.Fatal("context presence changed")
			}
			if ab03aq03Key != "family" && ab03aq03Key != "variant" && ab03aq03Got["state"] != "Unknown" {
				ab03aq03T.Fatal("absent or zero identity admitted")
			}
			if ab03aq03Got["qualified"] != false {
				ab03aq03T.Fatal("opaque qualification")
			}
		}
	}
	for _, ab03aq03Value := range []any{[]any{}, []any{0}, []any{255}, []any{25, 26}} {
		ab03aq03Input := ab03aq03Input(ab03aq03T, 0)
		ab03aq03Input["fields"].(map[string]any)["pump_power"] = ab03aq03Value
		ab03aq03Got := ab03aq03Run(ab03aq03T, ab03aq03Input)
		ab03aq03Result := ab03aq03Got["results"].(map[string]any)["pump_power"].(map[string]any)
		if len(ab03aq03Result) != 2 || ab03aq03Result["state"] != "Unknown" {
			ab03aq03T.Fatal("raw pump converted")
		}
		ab03aq03Subset(ab03aq03T, map[string]any{"number": 25, "candidate_domain": true}, ab03aq03Got["results"].(map[string]any)["display_percent"])
	}
}
