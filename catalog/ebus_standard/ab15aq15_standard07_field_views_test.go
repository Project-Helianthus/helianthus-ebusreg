package ebus_standard_catalog

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

type ab15aq15Case struct {
	Ab15Aq15ID       string                    `json:"id"`
	Ab15Aq15Input    json.RawMessage           `json:"input"`
	Ab15Aq15Expected map[string]map[string]any `json:"expected"`
}

func ab15aq15Cases(ab15aq15T *testing.T) []ab15aq15Case {
	ab15aq15T.Helper()
	ab15aq15Data, ab15aq15Err := os.ReadFile("testdata/ab15aq15_standard07_field_views.json")
	if ab15aq15Err != nil {
		ab15aq15T.Fatal(ab15aq15Err)
	}
	var ab15aq15Cases []ab15aq15Case
	if ab15aq15Err := json.Unmarshal(ab15aq15Data, &ab15aq15Cases); ab15aq15Err != nil {
		ab15aq15T.Fatal(ab15aq15Err)
	}
	if len(ab15aq15Cases) != 12 {
		ab15aq15T.Fatal("fixture count", len(ab15aq15Cases))
	}
	return ab15aq15Cases
}

func ab15aq15Evaluate(ab15aq15T *testing.T, ab15aq15Input any) map[string]any {
	ab15aq15T.Helper()
	ab15aq15Data, ab15aq15Err := json.Marshal(ab15aq15Input)
	if ab15aq15Err != nil {
		ab15aq15T.Fatal(ab15aq15Err)
	}
	return ab15aq15Decode(ab15aq15T, ab15aq15Data)
}

func ab15aq15Decode(ab15aq15T *testing.T, ab15aq15Input []byte) map[string]any {
	ab15aq15T.Helper()
	ab15aq15Data, ab15aq15Err := Ab15Aq15Evaluate(ab15aq15Input)
	if ab15aq15Err != nil {
		ab15aq15T.Fatal(ab15aq15Err)
	}
	var ab15aq15Output map[string]any
	if ab15aq15Err := json.Unmarshal(ab15aq15Data, &ab15aq15Output); ab15aq15Err != nil {
		ab15aq15T.Fatal(ab15aq15Err)
	}
	return ab15aq15Output
}

func ab15aq15Input(ab15aq15T *testing.T, ab15aq15Raw []byte) map[string]any {
	ab15aq15T.Helper()
	var ab15aq15Input map[string]any
	if ab15aq15Err := json.Unmarshal(ab15aq15Raw, &ab15aq15Input); ab15aq15Err != nil {
		ab15aq15T.Fatal(ab15aq15Err)
	}
	return ab15aq15Input
}

func TestAb15Aq15Fixture(ab15aq15T *testing.T) {
	for _, ab15aq15Case := range ab15aq15Cases(ab15aq15T) {
		ab15aq15T.Run(ab15aq15Case.Ab15Aq15ID, func(ab15aq15T *testing.T) {
			ab15aq15Input := ab15aq15Input(ab15aq15T, ab15aq15Case.Ab15Aq15Input)
			ab15aq15Output := ab15aq15Decode(ab15aq15T, ab15aq15Case.Ab15Aq15Input)
			if ab15aq15Output["state"] != "complete" || ab15aq15Output["physical"] != "Unknown" || ab15aq15Output["qualified"] != false {
				ab15aq15T.Fatal(ab15aq15Output)
			}
			if !reflect.DeepEqual(ab15aq15Input["context"], ab15aq15Output["context"]) || !reflect.DeepEqual(ab15aq15Input["fields"], ab15aq15Output["fields"]) {
				ab15aq15T.Fatal("raw input lost")
			}
			for ab15aq15Key := range ab15aq15Output {
				if _, ab15aq15Present := ab15aq15Input[ab15aq15Key]; !ab15aq15Present && ab15aq15Key != "state" && ab15aq15Key != "results" && ab15aq15Key != "physical" && ab15aq15Key != "qualified" {
					ab15aq15T.Fatal("invented output", ab15aq15Key)
				}
			}
			ab15aq15Results := ab15aq15Output["results"].(map[string]any)
			for ab15aq15Key, ab15aq15Expected := range ab15aq15Case.Ab15Aq15Expected {
				ab15aq15Result := ab15aq15Results[ab15aq15Key].(map[string]any)
				if ab15aq15Result["state"] != "complete" || ab15aq15Result["grade"] != "Hypothesis" || !reflect.DeepEqual(ab15aq15Result["raw"], ab15aq15Input["fields"].(map[string]any)[ab15aq15Key]) {
					ab15aq15T.Fatal(ab15aq15Result)
				}
				for ab15aq15Property, ab15aq15Value := range ab15aq15Expected {
					if !reflect.DeepEqual(ab15aq15Result[ab15aq15Property], ab15aq15Value) {
						ab15aq15T.Fatalf("%s.%s got %v want %v", ab15aq15Key, ab15aq15Property, ab15aq15Result[ab15aq15Property], ab15aq15Value)
					}
				}
				ab15aq15Allowed := map[string]bool{"raw": true, "state": true, "grade": true}
				switch ab15aq15Key {
				case "temperature":
					for _, ab15aq15Property := range []string{"number", "unit", "calculation", "candidate_domain"} {
						ab15aq15Allowed[ab15aq15Property] = true
					}
				case "duration":
					for _, ab15aq15Property := range []string{"number", "unit", "calculation"} {
						ab15aq15Allowed[ab15aq15Property] = true
					}
				case "weekday":
					ab15aq15Allowed["label"] = true
				case "second", "minute", "hour", "day", "month", "year":
					ab15aq15Allowed["candidate_domain"] = true
					ab15aq15Allowed["conversion"] = true
				}
				for ab15aq15Property := range ab15aq15Result {
					if !ab15aq15Allowed[ab15aq15Property] {
						ab15aq15T.Fatalf("invented property %s", ab15aq15Property)
					}
				}
			}
		})
	}
}

func TestAb15Aq15Axes(ab15aq15T *testing.T) {
	for _, ab15aq15Case := range ab15aq15Cases(ab15aq15T) {
		for ab15aq15Key, ab15aq15Value := range map[string]any{"primary": 8, "secondary": 1, "operation": 0, "direction": "negative", "view": "other"} {
			ab15aq15T.Run(ab15aq15Case.Ab15Aq15ID+"_"+ab15aq15Key, func(ab15aq15T *testing.T) {
				ab15aq15Input := ab15aq15Input(ab15aq15T, ab15aq15Case.Ab15Aq15Input)
				ab15aq15Input["context"].(map[string]any)[ab15aq15Key] = ab15aq15Value
				ab15aq15Output := ab15aq15Evaluate(ab15aq15T, ab15aq15Input)
				if ab15aq15Output["state"] != "Unknown" {
					ab15aq15T.Fatal(ab15aq15Output)
				}
				for _, ab15aq15Value := range ab15aq15Output["results"].(map[string]any) {
					ab15aq15Result := ab15aq15Value.(map[string]any)
					if len(ab15aq15Result) != 2 || ab15aq15Result["state"] != "Unknown" {
						ab15aq15T.Fatal(ab15aq15Result)
					}
				}
			})
		}
		for _, ab15aq15Key := range []string{"family", "variant"} {
			for _, ab15aq15Value := range []any{nil, false, float64(0), "", "other-family"} {
				ab15aq15Input := ab15aq15Input(ab15aq15T, ab15aq15Case.Ab15Aq15Input)
				ab15aq15Input["context"].(map[string]any)[ab15aq15Key] = ab15aq15Value
				ab15aq15Output := ab15aq15Evaluate(ab15aq15T, ab15aq15Input)
				if ab15aq15Output["qualified"] != false || ab15aq15Output["physical"] != "Unknown" || !reflect.DeepEqual(ab15aq15Output["results"], ab15aq15Decode(ab15aq15T, ab15aq15Case.Ab15Aq15Input)["results"]) || !reflect.DeepEqual(ab15aq15Output["context"], ab15aq15Input["context"]) {
					ab15aq15T.Fatal(ab15aq15Output)
				}
			}
		}
	}
}

func TestAb15Aq15WidthsAndIndependence(ab15aq15T *testing.T) {
	for _, ab15aq15Case := range ab15aq15Cases(ab15aq15T) {
		ab15aq15Original := ab15aq15Input(ab15aq15T, ab15aq15Case.Ab15Aq15Input)
		for ab15aq15Key, ab15aq15Raw := range ab15aq15Original["fields"].(map[string]any) {
			for ab15aq15Index, ab15aq15Replacement := range []any{nil, []any{}, append(ab15aq15Raw.([]any), float64(0)), ab15aq15Raw.([]any)[:len(ab15aq15Raw.([]any))-1]} {
				ab15aq15Input := ab15aq15Input(ab15aq15T, ab15aq15Case.Ab15Aq15Input)
				ab15aq15Fields := ab15aq15Input["fields"].(map[string]any)
				if ab15aq15Index == 0 {
					delete(ab15aq15Fields, ab15aq15Key)
				} else {
					ab15aq15Fields[ab15aq15Key] = ab15aq15Replacement
				}
				ab15aq15Output := ab15aq15Evaluate(ab15aq15T, ab15aq15Input)
				ab15aq15Results := ab15aq15Output["results"].(map[string]any)
				ab15aq15Result := ab15aq15Results[ab15aq15Key].(map[string]any)
				if ab15aq15Output["state"] != "partial" || ab15aq15Result["state"] != "incomplete" || ab15aq15Result["number"] != nil {
					ab15aq15T.Fatal(ab15aq15Output)
				}
				if ab15aq15Index == 0 {
					if ab15aq15Result["absent"] != true {
						ab15aq15T.Fatal(ab15aq15Result)
					}
					if _, ab15aq15Found := ab15aq15Result["raw"]; ab15aq15Found {
						ab15aq15T.Fatal("invented raw")
					}
				} else if !reflect.DeepEqual(ab15aq15Result["raw"], ab15aq15Replacement) {
					ab15aq15T.Fatal("raw changed")
				}
				ab15aq15Control := ab15aq15Decode(ab15aq15T, ab15aq15Case.Ab15Aq15Input)["results"].(map[string]any)
				for ab15aq15Sibling, ab15aq15Value := range ab15aq15Control {
					if ab15aq15Sibling != ab15aq15Key && !reflect.DeepEqual(ab15aq15Results[ab15aq15Sibling], ab15aq15Value) {
						ab15aq15T.Fatal("sibling lost")
					}
				}
			}
			ab15aq15Input := ab15aq15Input(ab15aq15T, ab15aq15Case.Ab15Aq15Input)
			ab15aq15Array := ab15aq15Input["fields"].(map[string]any)[ab15aq15Key].([]any)
			ab15aq15Array[0] = float64(129)
			ab15aq15Output := ab15aq15Evaluate(ab15aq15T, ab15aq15Input)
			ab15aq15Control := ab15aq15Decode(ab15aq15T, ab15aq15Case.Ab15Aq15Input)["results"].(map[string]any)
			for ab15aq15Sibling, ab15aq15Value := range ab15aq15Control {
				if ab15aq15Sibling != ab15aq15Key && !reflect.DeepEqual(ab15aq15Output["results"].(map[string]any)[ab15aq15Sibling], ab15aq15Value) {
					ab15aq15T.Fatal("field coupling")
				}
			}
		}
	}
}

func TestAb15Aq15Reject(ab15aq15T *testing.T) {
	for _, ab15aq15Raw := range []string{``, `null`, `[]`, `{} {} `, `{"context":{},"context":{},"fields":{}}`, `{"context":{"x":{"a":1,"a":2}},"fields":{}}`, `{"context":{},"fields":{},"x":[{"a":1,"a":2}]}`, `{"context":{},"fields":{},"x":1e999}`, `{"context":{},"fields":{}} true`, `{"context":{},"fields":{},}`} {
		ab15aq15Result, ab15aq15Err := Ab15Aq15Evaluate([]byte(ab15aq15Raw))
		if ab15aq15Err == nil || ab15aq15Result != nil {
			ab15aq15T.Fatalf("accepted %s", ab15aq15Raw)
		}
	}
	for _, ab15aq15Value := range []any{nil, true, "1", 1.5, -1, 256, []any{}, map[string]any{}} {
		for _, ab15aq15Key := range []string{"primary", "secondary", "operation", "field"} {
			ab15aq15Input := ab15aq15Input(ab15aq15T, ab15aq15Cases(ab15aq15T)[0].Ab15Aq15Input)
			if ab15aq15Key == "field" {
				ab15aq15Input["fields"].(map[string]any)["temperature"] = []any{ab15aq15Value, 0}
			} else {
				ab15aq15Input["context"].(map[string]any)[ab15aq15Key] = ab15aq15Value
			}
			ab15aq15Data, _ := json.Marshal(ab15aq15Input)
			ab15aq15Result, ab15aq15Err := Ab15Aq15Evaluate(ab15aq15Data)
			if ab15aq15Err == nil || ab15aq15Result != nil {
				ab15aq15T.Fatalf("accepted %s", ab15aq15Data)
			}
		}
	}
	for _, ab15aq15Key := range []string{"context", "fields", "arguments"} {
		for _, ab15aq15Value := range []any{nil, false, float64(0), "", []any{}} {
			ab15aq15Input := ab15aq15Input(ab15aq15T, ab15aq15Cases(ab15aq15T)[0].Ab15Aq15Input)
			ab15aq15Input[ab15aq15Key] = ab15aq15Value
			ab15aq15Data, _ := json.Marshal(ab15aq15Input)
			if ab15aq15Result, ab15aq15Err := Ab15Aq15Evaluate(ab15aq15Data); ab15aq15Err == nil || ab15aq15Result != nil {
				ab15aq15T.Fatal("wrong kind accepted")
			}
		}
	}
	ab15aq15Input := ab15aq15Input(ab15aq15T, ab15aq15Cases(ab15aq15T)[0].Ab15Aq15Input)
	ab15aq15Input["arguments"] = map[string]any{"unexpected": 0}
	ab15aq15Data, _ := json.Marshal(ab15aq15Input)
	if ab15aq15Result, ab15aq15Err := Ab15Aq15Evaluate(ab15aq15Data); ab15aq15Err == nil || ab15aq15Result != nil {
		ab15aq15T.Fatal("argument accepted")
	}
}

func TestAb15Aq15Storage(ab15aq15T *testing.T) {
	ab15aq15Input := ab15aq15Input(ab15aq15T, ab15aq15Cases(ab15aq15T)[0].Ab15Aq15Input)
	ab15aq15Input["arguments"] = map[string]any{}
	ab15aq15Input["extra"] = []any{nil, false, float64(0), "", []any{}}
	ab15aq15Input["fields"].(map[string]any)["opaque"] = []any{}
	ab15aq15Output := ab15aq15Evaluate(ab15aq15T, ab15aq15Input)
	for _, ab15aq15Key := range []string{"arguments", "extra", "fields"} {
		if !reflect.DeepEqual(ab15aq15Output[ab15aq15Key], ab15aq15Input[ab15aq15Key]) {
			ab15aq15T.Fatal("retention", ab15aq15Key)
		}
	}
	ab15aq15Data, _ := json.Marshal(ab15aq15Input)
	ab15aq15First, ab15aq15Err := Ab15Aq15Evaluate(ab15aq15Data)
	if ab15aq15Err != nil {
		ab15aq15T.Fatal(ab15aq15Err)
	}
	ab15aq15Saved := append([]byte(nil), ab15aq15First...)
	for ab15aq15Index := range ab15aq15Data {
		ab15aq15Data[ab15aq15Index] = ' '
	}
	if !reflect.DeepEqual([]byte(ab15aq15First), ab15aq15Saved) {
		ab15aq15T.Fatal("input alias")
	}
	for ab15aq15Index := range ab15aq15First {
		ab15aq15First[ab15aq15Index] = ' '
	}
	if !reflect.DeepEqual(ab15aq15Evaluate(ab15aq15T, ab15aq15Input), ab15aq15Output) {
		ab15aq15T.Fatal("output alias")
	}
}

func TestAb15Aq15ExactNumbersAndMissingSelectors(ab15aq15T *testing.T) {
	for _, ab15aq15Number := range []string{"1e-999", "255.000000000000000001", "0.000000000000000001"} {
		for _, ab15aq15Raw := range []string{`{"context":{"primary":` + ab15aq15Number + `},"fields":{}}`, `{"context":{},"fields":{"x":[` + ab15aq15Number + `]}}`} {
			if ab15aq15Result, ab15aq15Err := Ab15Aq15Evaluate([]byte(ab15aq15Raw)); ab15aq15Err == nil || ab15aq15Result != nil {
				ab15aq15T.Fatal("fraction accepted", ab15aq15Raw)
			}
		}
	}
	for _, ab15aq15Key := range []string{"primary", "secondary", "direction", "view"} {
		ab15aq15Input := ab15aq15Input(ab15aq15T, ab15aq15Cases(ab15aq15T)[0].Ab15Aq15Input)
		delete(ab15aq15Input["context"].(map[string]any), ab15aq15Key)
		if ab15aq15Output := ab15aq15Evaluate(ab15aq15T, ab15aq15Input); ab15aq15Output["state"] != "Unknown" {
			ab15aq15T.Fatal(ab15aq15Output)
		}
	}
	for _, ab15aq15Secondary := range []int{1, 4, 254, 255} {
		ab15aq15Input := ab15aq15Input(ab15aq15T, ab15aq15Cases(ab15aq15T)[0].Ab15Aq15Input)
		ab15aq15Input["context"].(map[string]any)["secondary"] = ab15aq15Secondary
		if ab15aq15Output := ab15aq15Evaluate(ab15aq15T, ab15aq15Input); ab15aq15Output["state"] != "Unknown" {
			ab15aq15T.Fatal(ab15aq15Output)
		}
	}
	for _, ab15aq15Key := range []string{"direction", "view"} {
		for _, ab15aq15Value := range []any{nil, 0, false, []any{}, map[string]any{}} {
			ab15aq15Input := ab15aq15Input(ab15aq15T, ab15aq15Cases(ab15aq15T)[0].Ab15Aq15Input)
			ab15aq15Input["context"].(map[string]any)[ab15aq15Key] = ab15aq15Value
			ab15aq15Data, _ := json.Marshal(ab15aq15Input)
			if ab15aq15Result, ab15aq15Err := Ab15Aq15Evaluate(ab15aq15Data); ab15aq15Err == nil || ab15aq15Result != nil {
				ab15aq15T.Fatal("identifier kind accepted")
			}
		}
	}
}

func TestAb15Aq15ReservedKeysAndUnsupportedArguments(ab15aq15T *testing.T) {
	for _, ab15aq15Key := range []string{"results", "state", "physical", "qualified"} {
		ab15aq15Input := ab15aq15Input(ab15aq15T, ab15aq15Cases(ab15aq15T)[0].Ab15Aq15Input)
		ab15aq15Input[ab15aq15Key] = "x"
		ab15aq15Data, _ := json.Marshal(ab15aq15Input)
		if ab15aq15Result, ab15aq15Err := Ab15Aq15Evaluate(ab15aq15Data); ab15aq15Err == nil || ab15aq15Result != nil {
			ab15aq15T.Fatal("reserved key accepted", ab15aq15Key)
		}
	}
	ab15aq15Input := ab15aq15Input(ab15aq15T, ab15aq15Cases(ab15aq15T)[0].Ab15Aq15Input)
	ab15aq15Input["context"].(map[string]any)["secondary"] = float64(1)
	ab15aq15Input["arguments"] = map[string]any{"unexpected": float64(0)}
	ab15aq15Output := ab15aq15Evaluate(ab15aq15T, ab15aq15Input)
	if ab15aq15Output["state"] != "Unknown" || !reflect.DeepEqual(ab15aq15Output["arguments"], ab15aq15Input["arguments"]) {
		ab15aq15T.Fatal("unsupported context must stay raw-only Unknown")
	}
}

func TestAb15Aq15DepthLimit(ab15aq15T *testing.T) {
	ab15aq15Deep := `{"context":{},"fields":{},"x":` + strings.Repeat("[", 100) + strings.Repeat("]", 100) + `}`
	if ab15aq15Result, ab15aq15Err := Ab15Aq15Evaluate([]byte(ab15aq15Deep)); ab15aq15Err == nil || ab15aq15Result != nil {
		ab15aq15T.Fatal("deep nesting accepted")
	}
}

func TestAb15Aq15RejectArgumentsOnUnsupportedContext(ab15aq15T *testing.T) {
	for _, ab15aq15Raw := range []string{`{"context":{"primary":8},"fields":{},"arguments":5}`, `{"context":{},"fields":{},"arguments":[]}`} {
		if ab15aq15Result, ab15aq15Err := Ab15Aq15Evaluate([]byte(ab15aq15Raw)); ab15aq15Err == nil || ab15aq15Result != nil {
			ab15aq15T.Fatal("arguments accepted", ab15aq15Raw)
		}
	}
}
