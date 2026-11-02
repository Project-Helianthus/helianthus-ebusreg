package vaillant

import (
	ab07aq07json "encoding/json"
	ab07aq07fmt "fmt"
	ab07aq07os "os"
	ab07aq07reflect "reflect"
	ab07aq07strings "strings"
	ab07aq07testing "testing"
)

func ab07aq07EvaluateTest(ab07aq07T *ab07aq07testing.T, ab07aq07Input any) map[string]any {
	ab07aq07T.Helper()
	ab07aq07Bytes, ab07aq07Error := ab07aq07json.Marshal(ab07aq07Input)
	if ab07aq07Error != nil {
		ab07aq07T.Fatal(ab07aq07Error)
	}
	ab07aq07Bytes, ab07aq07Error = Ab07Aq07Evaluate(ab07aq07Bytes)
	if ab07aq07Error != nil {
		ab07aq07T.Fatal(ab07aq07Error)
	}
	ab07aq07Output := map[string]any{}
	if ab07aq07Error = ab07aq07json.Unmarshal(ab07aq07Bytes, &ab07aq07Output); ab07aq07Error != nil {
		ab07aq07T.Fatal(ab07aq07Error)
	}
	return ab07aq07Output
}

func ab07aq07Subset(ab07aq07T *ab07aq07testing.T, ab07aq07Want, ab07aq07Got any) {
	ab07aq07T.Helper()
	if ab07aq07Map, ab07aq07OK := ab07aq07Want.(map[string]any); ab07aq07OK {
		ab07aq07Actual, ab07aq07OK := ab07aq07Got.(map[string]any)
		if !ab07aq07OK {
			ab07aq07T.Fatalf("want object %v got %v", ab07aq07Want, ab07aq07Got)
		}
		for ab07aq07Key, ab07aq07Value := range ab07aq07Map {
			ab07aq07Subset(ab07aq07T, ab07aq07Value, ab07aq07Actual[ab07aq07Key])
		}
	} else if !ab07aq07reflect.DeepEqual(ab07aq07Want, ab07aq07Got) {
		ab07aq07T.Fatalf("want %v got %v", ab07aq07Want, ab07aq07Got)
	}
}

func ab07aq07Fixtures(ab07aq07T *ab07aq07testing.T) []map[string]any {
	ab07aq07T.Helper()
	ab07aq07Bytes, ab07aq07Error := ab07aq07os.ReadFile("testdata/ab07aq07_b503_control_fields.json")
	if ab07aq07Error != nil {
		ab07aq07T.Fatal(ab07aq07Error)
	}
	var ab07aq07Cases []map[string]any
	if ab07aq07Error = ab07aq07json.Unmarshal(ab07aq07Bytes, &ab07aq07Cases); ab07aq07Error != nil {
		ab07aq07T.Fatal(ab07aq07Error)
	}
	if len(ab07aq07Cases) != 22 {
		ab07aq07T.Fatalf("fixture count %d", len(ab07aq07Cases))
	}
	return ab07aq07Cases
}

func TestAb07Aq07Fixtures(ab07aq07T *ab07aq07testing.T) {
	for _, ab07aq07Case := range ab07aq07Fixtures(ab07aq07T) {
		ab07aq07T.Run(ab07aq07Case["ab07aq07id"].(string), func(ab07aq07T *ab07aq07testing.T) {
			ab07aq07Input := ab07aq07Case["input"].(map[string]any)
			ab07aq07Output := ab07aq07EvaluateTest(ab07aq07T, ab07aq07Input)
			ab07aq07Subset(ab07aq07T, ab07aq07Case["expected"], ab07aq07Output)
			for ab07aq07Key, ab07aq07Value := range ab07aq07Input {
				if !ab07aq07reflect.DeepEqual(ab07aq07Value, ab07aq07Output[ab07aq07Key]) {
					ab07aq07T.Fatalf("lost input %s", ab07aq07Key)
				}
			}
			if len(ab07aq07Output) != len(ab07aq07Input)+4 {
				ab07aq07T.Fatal("invented output claim")
			}
			ab07aq07Results := ab07aq07Output["results"].(map[string]any)
			if ab07aq07Value, ab07aq07Present := ab07aq07Results["flags"]; ab07aq07Present {
				ab07aq07Result := ab07aq07Value.(map[string]any)
				if ab07aq07Flags, ab07aq07OK := ab07aq07Result["flags"].(map[string]any); ab07aq07OK && len(ab07aq07Flags) != 5 {
					ab07aq07T.Fatal("reserved bits labeled")
				}
				if _, ab07aq07Present := ab07aq07Result["label"]; ab07aq07Present {
					ab07aq07T.Fatal("flags given semantic label")
				}
			}
			for _, ab07aq07Value := range ab07aq07Results {
				ab07aq07Result := ab07aq07Value.(map[string]any)
				for _, ab07aq07Key := range []string{"reset_success", "effects", "freshness", "qualified", "executed", "delivery"} {
					if _, ab07aq07Present := ab07aq07Result[ab07aq07Key]; ab07aq07Present {
						ab07aq07T.Fatalf("invented %s", ab07aq07Key)
					}
				}
				if ab07aq07Result["state"] == "Unknown" {
					if _, ab07aq07Present := ab07aq07Result["label"]; ab07aq07Present {
						ab07aq07T.Fatal("invented label")
					}
				}
			}
		})
	}
}

func TestAb07Aq07Context(ab07aq07T *ab07aq07testing.T) {
	for _, ab07aq07Case := range ab07aq07Fixtures(ab07aq07T) {
		for ab07aq07Key, ab07aq07Values := range map[string][]any{"primary": {float64(180), nil}, "secondary": {float64(4), nil}, "operation": {float64(0), float64(1), nil}, "direction": {"other", nil}, "view": {"other", nil}, "family": {"other", nil, false, float64(0)}, "variant": {"other", nil, false, float64(0)}} {
			for _, ab07aq07Value := range ab07aq07Values {
				ab07aq07T.Run(ab07aq07fmt.Sprint(ab07aq07Case["ab07aq07id"], ab07aq07Key, ab07aq07Value), func(ab07aq07T *ab07aq07testing.T) {
					ab07aq07Bytes, _ := ab07aq07json.Marshal(ab07aq07Case["input"])
					ab07aq07Input := map[string]any{}
					_ = ab07aq07json.Unmarshal(ab07aq07Bytes, &ab07aq07Input)
					ab07aq07Context := ab07aq07Input["context"].(map[string]any)
					if ab07aq07Value == nil {
						delete(ab07aq07Context, ab07aq07Key)
					} else {
						ab07aq07Context[ab07aq07Key] = ab07aq07Value
					}
					ab07aq07Output := ab07aq07EvaluateTest(ab07aq07T, ab07aq07Input)
					if !ab07aq07reflect.DeepEqual(ab07aq07Context, ab07aq07Output["context"]) || ab07aq07Output["qualified"] != false {
						ab07aq07T.Fatal("identity lost or qualified")
					}
					if ab07aq07Key == "family" || ab07aq07Key == "variant" {
						ab07aq07Subset(ab07aq07T, ab07aq07Case["expected"], ab07aq07Output)
					} else if ab07aq07Output["state"] != "Unknown" || len(ab07aq07Output["results"].(map[string]any)) != 0 {
						ab07aq07T.Fatal("unsupported context interpreted")
					}
				})
			}
		}
	}
	for _, ab07aq07View := range []string{"reset-acceptance", "category-reset"} {
		ab07aq07Direction := "positive"
		ab07aq07Field := "acceptance"
		if ab07aq07View == "category-reset" {
			ab07aq07Direction = "request"
			ab07aq07Field = "category"
		}
		ab07aq07Output := ab07aq07EvaluateTest(ab07aq07T, map[string]any{"context": map[string]any{"primary": 181, "secondary": 3, "operation": 6, "direction": ab07aq07Direction, "view": ab07aq07View}, "fields": map[string]any{ab07aq07Field: []int{1}}})
		if ab07aq07Output["state"] != "complete" {
			ab07aq07T.Fatal("OP06 rejected")
		}
	}
}

func TestAb07Aq07WidthsAndFields(ab07aq07T *ab07aq07testing.T) {
	for _, ab07aq07Case := range ab07aq07Fixtures(ab07aq07T) {
		ab07aq07Original := ab07aq07Case["input"].(map[string]any)
		for ab07aq07Key := range ab07aq07Original["fields"].(map[string]any) {
			for _, ab07aq07Replacement := range []any{nil, []int{}, []int{0}, []int{0, 0}, []int{0, 0, 0}} {
				ab07aq07T.Run(ab07aq07fmt.Sprint(ab07aq07Case["ab07aq07id"], ab07aq07Key, ab07aq07Replacement), func(ab07aq07T *ab07aq07testing.T) {
					ab07aq07Bytes, _ := ab07aq07json.Marshal(ab07aq07Original)
					ab07aq07Input := map[string]any{}
					_ = ab07aq07json.Unmarshal(ab07aq07Bytes, &ab07aq07Input)
					ab07aq07Fields := ab07aq07Input["fields"].(map[string]any)
					if ab07aq07Replacement == nil {
						delete(ab07aq07Fields, ab07aq07Key)
					} else {
						ab07aq07Fields[ab07aq07Key] = ab07aq07Replacement
					}
					ab07aq07Output := ab07aq07EvaluateTest(ab07aq07T, ab07aq07Input)
					ab07aq07Control := ab07aq07EvaluateTest(ab07aq07T, ab07aq07Original)
					ab07aq07Results := ab07aq07Output["results"].(map[string]any)
					for ab07aq07Sibling, ab07aq07Value := range ab07aq07Control["results"].(map[string]any) {
						if ab07aq07Sibling != ab07aq07Key && !ab07aq07reflect.DeepEqual(ab07aq07Value, ab07aq07Results[ab07aq07Sibling]) {
							ab07aq07T.Fatalf("changed sibling %s", ab07aq07Sibling)
						}
					}
					ab07aq07Width := 1
					if ab07aq07strings.HasPrefix(ab07aq07Key, "error") {
						ab07aq07Width = 2
					}
					ab07aq07Array, _ := ab07aq07Replacement.([]int)
					ab07aq07Result := ab07aq07Results[ab07aq07Key].(map[string]any)
					if len(ab07aq07Array) != ab07aq07Width {
						if ab07aq07Output["state"] != "partial" || ab07aq07Result["state"] != "incomplete" {
							ab07aq07T.Fatal("bad width complete")
						}
						if _, ab07aq07Present := ab07aq07Result["number"]; ab07aq07Present {
							ab07aq07T.Fatal("invented number")
						}
					} else if ab07aq07Result["number"] != float64(0) {
						ab07aq07T.Fatal("field not independently evaluated")
					}
					if ab07aq07Result["present"] != (ab07aq07Replacement != nil) {
						ab07aq07T.Fatal("absence conflated")
					}
				})
			}
		}
	}
}

func TestAb07Aq07Reject(ab07aq07T *ab07aq07testing.T) {
	for _, ab07aq07Input := range []string{``, `{`, `[]`, `null`, `{} {}`, `{"context":{},"context":{}}`, `{"x":[{"a":0,"a":1}]}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":null}`} {
		if ab07aq07Result, ab07aq07Error := Ab07Aq07Evaluate([]byte(ab07aq07Input)); ab07aq07Error == nil || ab07aq07Result != nil {
			ab07aq07T.Fatalf("accepted %s", ab07aq07Input)
		}
	}
	for _, ab07aq07Bad := range []string{`null`, `true`, `"1"`, `1.5`, `-1`, `256`, `1e999`, `NaN`} {
		for _, ab07aq07Key := range []string{"primary", "secondary", "operation", "field", "arguments", "direction", "view"} {
			ab07aq07Input := `{"context":{"primary":181,"secondary":3,"operation":3,"direction":"positive","view":"reset-acceptance"},"fields":{"acceptance":[1]}}`
			switch ab07aq07Key {
			case "field":
				ab07aq07Input = ab07aq07strings.Replace(ab07aq07Input, `[1]`, `[`+ab07aq07Bad+`]`, 1)
			case "arguments":
				ab07aq07Input = ab07aq07strings.TrimSuffix(ab07aq07Input, "}") + `,"arguments":{"undefined":` + ab07aq07Bad + `}}`
			default:
				if (ab07aq07Key == "direction" || ab07aq07Key == "view") && ab07aq07Bad == `"1"` {
					continue
				}
				ab07aq07Old := map[string]string{"primary": "181", "secondary": "3", "operation": "3", "direction": `"positive"`, "view": `"reset-acceptance"`}[ab07aq07Key]
				ab07aq07Input = ab07aq07strings.Replace(ab07aq07Input, `"`+ab07aq07Key+`":`+ab07aq07Old, `"`+ab07aq07Key+`":`+ab07aq07Bad, 1)
			}
			if ab07aq07Result, ab07aq07Error := Ab07Aq07Evaluate([]byte(ab07aq07Input)); ab07aq07Error == nil || ab07aq07Result != nil {
				ab07aq07T.Fatalf("accepted %s", ab07aq07Input)
			}
		}
	}
}

func TestAb07Aq07MemoryAndExtras(ab07aq07T *ab07aq07testing.T) {
	ab07aq07Input := []byte(`{"context":{"primary":181,"secondary":3,"operation":3,"direction":"positive","view":"reset-acceptance"},"fields":{"acceptance":[1],"opaque":[]},"arguments":{},"extra":[null,0,false,[]]}`)
	ab07aq07Control, ab07aq07Error := Ab07Aq07Evaluate(ab07aq07Input)
	if ab07aq07Error != nil {
		ab07aq07T.Fatal(ab07aq07Error)
	}
	ab07aq07Original := append([]byte(nil), ab07aq07Input...)
	ab07aq07Output, _ := Ab07Aq07Evaluate(ab07aq07Input)
	for ab07aq07Index := range ab07aq07Input {
		ab07aq07Input[ab07aq07Index] = 'x'
	}
	if string(ab07aq07Output) != string(ab07aq07Control) {
		ab07aq07T.Fatal("input aliases output")
	}
	for ab07aq07Index := range ab07aq07Output {
		ab07aq07Output[ab07aq07Index] = 'x'
	}
	ab07aq07Again, ab07aq07Error := Ab07Aq07Evaluate(ab07aq07Original)
	if ab07aq07Error != nil || string(ab07aq07Again) != string(ab07aq07Control) {
		ab07aq07T.Fatal("cross-call alias")
	}
	ab07aq07Decoded := map[string]any{}
	_ = ab07aq07json.Unmarshal(ab07aq07Again, &ab07aq07Decoded)
	if !ab07aq07reflect.DeepEqual(ab07aq07Decoded["extra"], []any{nil, float64(0), false, []any{}}) {
		ab07aq07T.Fatal("extras lost")
	}
}

func TestAb07Aq07FeedbackDomain(ab07aq07T *ab07aq07testing.T) {
	for ab07aq07Number := 0; ab07aq07Number < 256; ab07aq07Number++ {
		ab07aq07Output := ab07aq07EvaluateTest(ab07aq07T, map[string]any{"context": map[string]any{"primary": 181, "secondary": 3, "operation": 7, "direction": "positive", "view": "current-errors"}, "fields": map[string]any{"feedback": []int{ab07aq07Number}}})
		ab07aq07Result := ab07aq07Output["results"].(map[string]any)["feedback"].(map[string]any)
		ab07aq07Label := "rejected"
		if ab07aq07Number == 1 {
			ab07aq07Label = "accepted"
		}
		if ab07aq07Result["label"] != ab07aq07Label || ab07aq07Output["state"] != "partial" {
			ab07aq07T.Fatalf("feedback %d: %v", ab07aq07Number, ab07aq07Output)
		}
		if len(ab07aq07Output) != 6 || len(ab07aq07Result) != 7 {
			ab07aq07T.Fatal("feedback inferred effects")
		}
	}
}

func TestAb07Aq07KindsAndPreservation(ab07aq07T *ab07aq07testing.T) {
	for _, ab07aq07Bad := range []string{`null`, `true`, `1`, `"x"`, `{}`} {
		for _, ab07aq07Input := range []string{`{"context":{},"fields":{"opaque":` + ab07aq07Bad + `}}`, `{"context":{},"fields":{},"arguments":` + ab07aq07Bad + `}`} {
			if ab07aq07Bad == `{}` && ab07aq07strings.Contains(ab07aq07Input, `"arguments"`) {
				ab07aq07Result, ab07aq07Error := Ab07Aq07Evaluate([]byte(ab07aq07Input))
				if ab07aq07Error != nil || ab07aq07Result == nil {
					ab07aq07T.Fatal("empty arguments rejected")
				}
				continue
			}
			if ab07aq07Result, ab07aq07Error := Ab07Aq07Evaluate([]byte(ab07aq07Input)); ab07aq07Error == nil || ab07aq07Result != nil {
				ab07aq07T.Fatalf("accepted %s", ab07aq07Input)
			}
		}
	}
	ab07aq07Input := map[string]any{"context": map[string]any{"primary": float64(181), "secondary": float64(3), "operation": float64(0), "direction": "positive", "view": "reset-acceptance", "variant": nil}, "fields": map[string]any{"acceptance": []any{}, "opaque": []any{float64(0)}}, "arguments": map[string]any{}, "extra": []any{nil, false, float64(0), []any{}}}
	ab07aq07Output := ab07aq07EvaluateTest(ab07aq07T, ab07aq07Input)
	for ab07aq07Key, ab07aq07Value := range ab07aq07Input {
		if !ab07aq07reflect.DeepEqual(ab07aq07Value, ab07aq07Output[ab07aq07Key]) {
			ab07aq07T.Fatalf("unsupported context lost %s", ab07aq07Key)
		}
	}
	if ab07aq07Output["state"] != "Unknown" || len(ab07aq07Output["results"].(map[string]any)) != 0 {
		ab07aq07T.Fatal("unsupported context evaluated")
	}
}
