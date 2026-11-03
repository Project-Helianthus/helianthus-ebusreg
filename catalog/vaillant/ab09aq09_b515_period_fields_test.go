package vaillant

import (
	ab09aq09bytes "bytes"
	ab09aq09json "encoding/json"
	ab09aq09fmt "fmt"
	ab09aq09os "os"
	ab09aq09reflect "reflect"
	ab09aq09strings "strings"
	ab09aq09testing "testing"
)

func ab09aq09MakeInput(ab09aq09View string) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 21, "direction": "request", "view": ab09aq09View}, "fields": map[string]any{"period0": []int{0}, "period1": []int{54}, "period2": []int{144}}}
}

func ab09aq09Run(ab09aq09T *ab09aq09testing.T, ab09aq09Input any) map[string]any {
	ab09aq09T.Helper()
	ab09aq09Encoded, ab09aq09Error := ab09aq09json.Marshal(ab09aq09Input)
	if ab09aq09Error != nil {
		ab09aq09T.Fatal(ab09aq09Error)
	}
	ab09aq09Output, ab09aq09Error := Ab09Aq09Evaluate(ab09aq09Encoded)
	if ab09aq09Error != nil {
		ab09aq09T.Fatal(ab09aq09Error)
	}
	ab09aq09Result := map[string]any{}
	if ab09aq09Error := ab09aq09json.Unmarshal(ab09aq09Output, &ab09aq09Result); ab09aq09Error != nil {
		ab09aq09T.Fatal(ab09aq09Error)
	}
	if ab09aq09Result["physical"] != "Unknown" || ab09aq09Result["qualified"] != false {
		ab09aq09T.Fatalf("qualification: %s", ab09aq09Output)
	}
	ab09aq09Retained := map[string]any{}
	if ab09aq09Error := ab09aq09json.Unmarshal(ab09aq09Encoded, &ab09aq09Retained); ab09aq09Error != nil {
		ab09aq09T.Fatal(ab09aq09Error)
	}
	for _, ab09aq09Key := range []string{"context", "fields", "arguments"} {
		ab09aq09Expected, ab09aq09ExpectedPresent := ab09aq09Retained[ab09aq09Key]
		ab09aq09Actual, ab09aq09ActualPresent := ab09aq09Result[ab09aq09Key]
		if ab09aq09ExpectedPresent != ab09aq09ActualPresent || !ab09aq09reflect.DeepEqual(ab09aq09Actual, ab09aq09Expected) {
			ab09aq09T.Fatalf("retention of %s", ab09aq09Key)
		}
		delete(ab09aq09Retained, ab09aq09Key)
	}
	if !ab09aq09reflect.DeepEqual(ab09aq09Result["extra"], ab09aq09Retained) {
		ab09aq09T.Fatal("extra keys lost")
	}
	return ab09aq09Result
}

func ab09aq09Subset(ab09aq09T *ab09aq09testing.T, ab09aq09Actual, ab09aq09Expected any) {
	ab09aq09T.Helper()
	switch ab09aq09ExpectedValue := ab09aq09Expected.(type) {
	case map[string]any:
		ab09aq09ActualValue, ab09aq09OK := ab09aq09Actual.(map[string]any)
		if !ab09aq09OK {
			ab09aq09T.Fatalf("expected object, got %v", ab09aq09Actual)
		}
		for ab09aq09Key, ab09aq09Value := range ab09aq09ExpectedValue {
			ab09aq09ActualChild, ab09aq09Present := ab09aq09ActualValue[ab09aq09Key]
			if !ab09aq09Present {
				ab09aq09T.Fatalf("missing %s", ab09aq09Key)
			}
			ab09aq09Subset(ab09aq09T, ab09aq09ActualChild, ab09aq09Value)
		}
	case []any:
		ab09aq09ActualValue, ab09aq09OK := ab09aq09Actual.([]any)
		if !ab09aq09OK || len(ab09aq09ActualValue) != len(ab09aq09ExpectedValue) {
			ab09aq09T.Fatalf("array mismatch %v", ab09aq09Actual)
		}
		for ab09aq09Index := range ab09aq09ExpectedValue {
			ab09aq09Subset(ab09aq09T, ab09aq09ActualValue[ab09aq09Index], ab09aq09ExpectedValue[ab09aq09Index])
		}
	default:
		if !ab09aq09reflect.DeepEqual(ab09aq09Actual, ab09aq09Expected) {
			ab09aq09T.Fatalf("got %v, want %v", ab09aq09Actual, ab09aq09Expected)
		}
	}
}

func ab09aq09GetResult(ab09aq09Output map[string]any, ab09aq09Key string) map[string]any {
	return ab09aq09Output["results"].(map[string]any)[ab09aq09Key].(map[string]any)
}

func ab09aq09NoNumber(ab09aq09T *ab09aq09testing.T, ab09aq09Result map[string]any) {
	ab09aq09T.Helper()
	for _, ab09aq09Key := range []string{"number", "hour", "minute", "label", "value", "calculation"} {
		if _, ab09aq09Present := ab09aq09Result[ab09aq09Key]; ab09aq09Present {
			ab09aq09T.Fatalf("unexpected %s in %v", ab09aq09Key, ab09aq09Result)
		}
	}
}

func TestAb09Aq09Fixture(ab09aq09T *ab09aq09testing.T) {
	ab09aq09Data, ab09aq09Error := ab09aq09os.ReadFile("testdata/ab09aq09_b515_period_fields.json")
	if ab09aq09Error != nil {
		ab09aq09T.Fatal(ab09aq09Error)
	}
	var ab09aq09Cases []map[string]any
	if ab09aq09Error := ab09aq09json.Unmarshal(ab09aq09Data, &ab09aq09Cases); ab09aq09Error != nil {
		ab09aq09T.Fatal(ab09aq09Error)
	}
	if len(ab09aq09Cases) != 5 {
		ab09aq09T.Fatal("fixture cases missing")
	}
	for _, ab09aq09Case := range ab09aq09Cases {
		ab09aq09T.Run(ab09aq09Case["id"].(string), func(ab09aq09T *ab09aq09testing.T) {
			ab09aq09Output := ab09aq09Run(ab09aq09T, ab09aq09Case["input"])
			ab09aq09Subset(ab09aq09T, ab09aq09Output, ab09aq09Case["expected"])
			if ab09aq09Case["id"] == "ab09aq09_unmapped" {
				for _, ab09aq09Key := range []string{"period0", "period1", "minutes"} {
					ab09aq09NoNumber(ab09aq09T, ab09aq09GetResult(ab09aq09Output, ab09aq09Key))
				}
			}
			if ab09aq09Case["id"] == "ab09aq09_unmapped_equal" {
				ab09aq09NoNumber(ab09aq09T, ab09aq09GetResult(ab09aq09Output, "start_order"))
			}
		})
	}
}

func TestAb09Aq09ExhaustPeriods(ab09aq09T *ab09aq09testing.T) {
	ab09aq09ExpectedMinutes := 0
	ab09aq09ExpectedHour, ab09aq09ExpectedMinute := 0, 0
	for ab09aq09Raw := 0; ab09aq09Raw <= 255; ab09aq09Raw++ {
		ab09aq09Input := ab09aq09MakeInput("periods")
		for _, ab09aq09Key := range []string{"period0", "period1", "period2"} {
			ab09aq09Input["fields"].(map[string]any)[ab09aq09Key] = []int{ab09aq09Raw}
		}
		ab09aq09Output := ab09aq09Run(ab09aq09T, ab09aq09Input)
		if ab09aq09Output["state"] != "complete" {
			ab09aq09T.Fatal(ab09aq09Output)
		}
		for _, ab09aq09Key := range []string{"period0", "period1", "period2"} {
			ab09aq09Result := ab09aq09GetResult(ab09aq09Output, ab09aq09Key)
			if ab09aq09Raw <= 144 {
				if ab09aq09Result["number"] != float64(ab09aq09ExpectedMinutes) || ab09aq09Result["hour"] != float64(ab09aq09ExpectedHour) || ab09aq09Result["minute"] != float64(ab09aq09ExpectedMinute) || ab09aq09Result["label"] != ab09aq09fmt.Sprintf("%02d:%02d", ab09aq09ExpectedHour, ab09aq09ExpectedMinute) || ab09aq09Result["state"] != "Hypothesis" || ab09aq09Result["calculation"] != "Proven" {
					ab09aq09T.Fatalf("raw %d: %v", ab09aq09Raw, ab09aq09Result)
				}
			} else {
				if ab09aq09Result["state"] != "Unknown" {
					ab09aq09T.Fatal(ab09aq09Result)
				}
				ab09aq09NoNumber(ab09aq09T, ab09aq09Result)
			}
		}
		ab09aq09ExpectedMinutes += 10
		ab09aq09ExpectedMinute += 10
		if ab09aq09ExpectedMinute == 60 {
			ab09aq09ExpectedHour++
			ab09aq09ExpectedMinute = 0
		}
	}
}

func TestAb09Aq09Inverse(ab09aq09T *ab09aq09testing.T) {
	for ab09aq09Minutes := 0; ab09aq09Minutes <= 1440; ab09aq09Minutes++ {
		ab09aq09Input := ab09aq09MakeInput("periods")
		ab09aq09Input["arguments"] = map[string]any{"minutes": ab09aq09Minutes}
		ab09aq09Output := ab09aq09Run(ab09aq09T, ab09aq09Input)
		ab09aq09Result := ab09aq09GetResult(ab09aq09Output, "minutes")
		if ab09aq09Minutes%10 != 0 {
			if ab09aq09Result["state"] != "Unknown" {
				ab09aq09T.Fatal(ab09aq09Result)
			}
			ab09aq09NoNumber(ab09aq09T, ab09aq09Result)
		} else {
			if ab09aq09Result["number"] != float64(ab09aq09Minutes/10) || ab09aq09Result["state"] != "Hypothesis" || ab09aq09Result["calculation"] != "Proven" {
				ab09aq09T.Fatal(ab09aq09Result)
			}
			ab09aq09Input["fields"].(map[string]any)["period0"] = []int{ab09aq09Minutes / 10}
			if ab09aq09GetResult(ab09aq09Run(ab09aq09T, ab09aq09Input), "period0")["number"] != float64(ab09aq09Minutes) {
				ab09aq09T.Fatal("round trip")
			}
		}
	}
}

func TestAb09Aq09ContextAxes(ab09aq09T *ab09aq09testing.T) {
	for ab09aq09Key, ab09aq09Values := range map[string][]any{"primary": {180, 0, 255}, "secondary": {20, 0, 255}, "operation": {0, 1, 255}, "direction": {"negative", "local", ""}, "view": {"other", ""}} {
		for _, ab09aq09Value := range ab09aq09Values {
			ab09aq09Input := ab09aq09MakeInput("periods")
			ab09aq09Input["context"].(map[string]any)[ab09aq09Key] = ab09aq09Value
			ab09aq09Output := ab09aq09Run(ab09aq09T, ab09aq09Input)
			if ab09aq09Output["state"] != "Unknown" || len(ab09aq09Output["results"].(map[string]any)) != 0 {
				ab09aq09T.Fatalf("%s=%v: %v", ab09aq09Key, ab09aq09Value, ab09aq09Output)
			}
		}
	}
	for _, ab09aq09Key := range []string{"primary", "secondary", "direction", "view"} {
		ab09aq09Input := ab09aq09MakeInput("periods")
		delete(ab09aq09Input["context"].(map[string]any), ab09aq09Key)
		if ab09aq09Run(ab09aq09T, ab09aq09Input)["state"] != "Unknown" {
			ab09aq09T.Fatal(ab09aq09Key)
		}
	}
	for _, ab09aq09Key := range []string{"family", "variant"} {
		for _, ab09aq09Value := range []any{nil, false, 0, "", []any{}, map[string]any{}, "synthetic-family"} {
			ab09aq09Input := ab09aq09MakeInput("periods")
			ab09aq09Input["context"].(map[string]any)[ab09aq09Key] = ab09aq09Value
			ab09aq09Output := ab09aq09Run(ab09aq09T, ab09aq09Input)
			ab09aq09Baseline := ab09aq09Run(ab09aq09T, ab09aq09MakeInput("periods"))
			if !ab09aq09reflect.DeepEqual(ab09aq09Output["results"], ab09aq09Baseline["results"]) {
				ab09aq09T.Fatal("opaque identity changed arithmetic")
			}
			ab09aq09Encoded, _ := ab09aq09json.Marshal(ab09aq09Value)
			var ab09aq09Expected any
			if ab09aq09Error := ab09aq09json.Unmarshal(ab09aq09Encoded, &ab09aq09Expected); ab09aq09Error != nil {
				ab09aq09T.Fatal(ab09aq09Error)
			}
			if !ab09aq09reflect.DeepEqual(ab09aq09Output["context"].(map[string]any)[ab09aq09Key], ab09aq09Expected) {
				ab09aq09T.Fatal("opaque identity lost")
			}
		}
	}
}

func TestAb09Aq09PartialAndRetention(ab09aq09T *ab09aq09testing.T) {
	for _, ab09aq09Key := range []string{"period0", "period1", "period2"} {
		for _, ab09aq09Width := range []int{-1, 0, 2} {
			ab09aq09Input := ab09aq09MakeInput("periods")
			ab09aq09Fields := ab09aq09Input["fields"].(map[string]any)
			if ab09aq09Width == -1 {
				delete(ab09aq09Fields, ab09aq09Key)
			} else {
				ab09aq09Fields[ab09aq09Key] = make([]int, ab09aq09Width)
			}
			ab09aq09Fields["opaque"] = []int{}
			ab09aq09Input["extra"] = map[string]any{"null": nil, "zero": 0, "false": false, "empty": []any{}}
			ab09aq09Input["state"] = "retained"
			ab09aq09Output := ab09aq09Run(ab09aq09T, ab09aq09Input)
			if ab09aq09Output["state"] != "partial" {
				ab09aq09T.Fatal(ab09aq09Output)
			}
			ab09aq09Result := ab09aq09GetResult(ab09aq09Output, ab09aq09Key)
			if ab09aq09Result["state"] != "incomplete" {
				ab09aq09T.Fatal(ab09aq09Result)
			}
			ab09aq09NoNumber(ab09aq09T, ab09aq09Result)
			if _, ab09aq09Present := ab09aq09Result["raw"]; ab09aq09Present != (ab09aq09Width != -1) {
				ab09aq09T.Fatal("absence lost")
			}
			for _, ab09aq09Sibling := range []string{"period0", "period1", "period2"} {
				if ab09aq09Sibling != ab09aq09Key && ab09aq09GetResult(ab09aq09Output, ab09aq09Sibling)["state"] != "Hypothesis" {
					ab09aq09T.Fatal("sibling lost")
				}
			}
			if ab09aq09Output["extra"].(map[string]any)["state"] != "retained" {
				ab09aq09T.Fatal("extra collision")
			}
			if len(ab09aq09GetResult(ab09aq09Output, "opaque")["raw"].([]any)) != 0 {
				ab09aq09T.Fatal("empty opaque lost")
			}
			ab09aq09Subset(ab09aq09T, ab09aq09Output["extra"], map[string]any{"extra": map[string]any{"null": nil, "zero": float64(0), "false": false, "empty": []any{}}})
		}
	}
}

func TestAb09Aq09IntervalsAxes(ab09aq09T *ab09aq09testing.T) {
	for ab09aq09Index := 0; ab09aq09Index < 3; ab09aq09Index++ {
		for ab09aq09Endpoint := 0; ab09aq09Endpoint < 2; ab09aq09Endpoint++ {
			for _, ab09aq09Value := range []int{0, 10, 20, 40, 144, 145, 255} {
				ab09aq09Pairs := [3][2]int{{10, 10}, {20, 30}, {40, 50}}
				ab09aq09Pairs[ab09aq09Index][ab09aq09Endpoint] = ab09aq09Value
				ab09aq09Input := ab09aq09MakeInput("intervals")
				ab09aq09Input["arguments"] = map[string]any{"pairs": ab09aq09Pairs}
				ab09aq09Output := ab09aq09Run(ab09aq09T, ab09aq09Input)
				ab09aq09Order := ab09aq09GetResult(ab09aq09Output, "start_order")
				if ab09aq09Value > 144 {
					if ab09aq09Order["state"] != "Unknown" {
						ab09aq09T.Fatal(ab09aq09Order)
					}
					ab09aq09NoNumber(ab09aq09T, ab09aq09Order)
				} else {
					ab09aq09Expected := ab09aq09Pairs[0][0] < ab09aq09Pairs[1][0] && ab09aq09Pairs[1][0] < ab09aq09Pairs[2][0]
					if ab09aq09Order["value"] != ab09aq09Expected {
						ab09aq09T.Fatal(ab09aq09Order)
					}
				}
				for ab09aq09PairIndex, ab09aq09PairResult := range ab09aq09GetResult(ab09aq09Output, "pairs")["alternatives"].([]any) {
					ab09aq09Expected := "enabled-candidate"
					if ab09aq09Pairs[ab09aq09PairIndex][0] == ab09aq09Pairs[ab09aq09PairIndex][1] {
						ab09aq09Expected = "disabled"
					}
					ab09aq09Result := ab09aq09PairResult.(map[string]any)
					if ab09aq09Result["label"] != ab09aq09Expected || ab09aq09Result["start"] != float64(ab09aq09Pairs[ab09aq09PairIndex][0]) || ab09aq09Result["end"] != float64(ab09aq09Pairs[ab09aq09PairIndex][1]) {
						ab09aq09T.Fatal(ab09aq09Result)
					}
					if _, ab09aq09Present := ab09aq09Result["deletion"]; ab09aq09Present {
						ab09aq09T.Fatal("invented deletion")
					}
				}
			}
		}
	}
	ab09aq09Input := ab09aq09MakeInput("intervals")
	ab09aq09Input["arguments"] = map[string]any{"pairs": [3][2]int{{140, 1}, {30, 10}, {20, 40}}}
	if ab09aq09Run(ab09aq09T, ab09aq09Input)["state"] != "complete" {
		ab09aq09T.Fatal("reversed intervals rejected")
	}
	ab09aq09Input["context"].(map[string]any)["direction"] = "positive"
	if ab09aq09Run(ab09aq09T, ab09aq09Input)["state"] != "Unknown" {
		ab09aq09T.Fatal("interval direction admitted")
	}
}

func TestAb09Aq09RejectJSON(ab09aq09T *ab09aq09testing.T) {
	for _, ab09aq09Input := range []string{
		``, `{`, `null`, `[]`, `{} {}`, `{} true`, `{"context":{},"context":{}}`,
		`{"context":{"family":{"x":0,"x":1}},"fields":{}}`,
		`{"context":{},"fields":{},"extra":[{"x":0,"x":1}]}`,
		`{"context":{},"fields":{"x":[NaN]}}`, `{"context":{},"fields":{},"extra":1e400}`,
		`{"context":{},"fields":{"x":[0],"x":[]}}`,
		`{"context":{},"fields":{},"arguments":{"x":0,"x":1}}`,
	} {
		ab09aq09Output, ab09aq09Error := Ab09Aq09Evaluate([]byte(ab09aq09Input))
		if ab09aq09Error == nil || ab09aq09Output != nil {
			ab09aq09T.Fatalf("accepted %q", ab09aq09Input)
		}
	}
}

func TestAb09Aq09RejectKinds(ab09aq09T *ab09aq09testing.T) {
	for _, ab09aq09Value := range []any{nil, true, "1", -1, 256, 1.5, []any{}, map[string]any{}} {
		for _, ab09aq09Key := range []string{"primary", "secondary", "operation", "byte", "minutes", "endpoint"} {
			ab09aq09Input := ab09aq09MakeInput("periods")
			switch ab09aq09Key {
			case "byte":
				ab09aq09Input["fields"].(map[string]any)["period0"] = []any{ab09aq09Value}
			case "minutes":
				if ab09aq09Value == 256 {
					continue
				}
				ab09aq09Input["arguments"] = map[string]any{"minutes": ab09aq09Value}
			case "endpoint":
				ab09aq09Input["arguments"] = map[string]any{"pairs": []any{[]any{ab09aq09Value, 10}, []int{20, 30}, []int{40, 50}}}
			default:
				ab09aq09Input["context"].(map[string]any)[ab09aq09Key] = ab09aq09Value
			}
			ab09aq09Encoded, _ := ab09aq09json.Marshal(ab09aq09Input)
			ab09aq09Output, ab09aq09Error := Ab09Aq09Evaluate(ab09aq09Encoded)
			if ab09aq09Error == nil || ab09aq09Output != nil {
				ab09aq09T.Fatalf("accepted %s", ab09aq09Encoded)
			}
		}
	}
	for _, ab09aq09Patch := range []string{
		`"arguments":{"minutes":1441}`, `"arguments":null`, `"fields":null`, `"context":null`,
		`"arguments":{"pairs":[]}`, `"arguments":{"pairs":[[1,2],[3,4],[5]]}`,
		`"arguments":{"pairs":[[1,2],[3,4],[5,6,7]]}`, `"fields":{"period0":null}`,
		`"context":{"direction":false}`, `"context":{"view":0}`,
	} {
		ab09aq09Input := ab09aq09MakeInput("periods")
		ab09aq09PatchObject := map[string]any{}
		if ab09aq09Error := ab09aq09json.Unmarshal([]byte("{"+ab09aq09Patch+"}"), &ab09aq09PatchObject); ab09aq09Error != nil {
			ab09aq09T.Fatal(ab09aq09Error)
		}
		for ab09aq09Key, ab09aq09Value := range ab09aq09PatchObject {
			ab09aq09Input[ab09aq09Key] = ab09aq09Value
		}
		ab09aq09Encoded, _ := ab09aq09json.Marshal(ab09aq09Input)
		if ab09aq09Output, ab09aq09Error := Ab09Aq09Evaluate(ab09aq09Encoded); ab09aq09Error == nil || ab09aq09Output != nil {
			ab09aq09T.Fatalf("accepted %s", ab09aq09Encoded)
		}
	}
	ab09aq09Encoded, _ := ab09aq09json.Marshal(ab09aq09MakeInput("intervals"))
	if ab09aq09Output, ab09aq09Error := Ab09Aq09Evaluate(ab09aq09Encoded); ab09aq09Error == nil || ab09aq09Output != nil {
		ab09aq09T.Fatal("missing pairs accepted")
	}
}

func TestAb09Aq09IndependentMemory(ab09aq09T *ab09aq09testing.T) {
	ab09aq09Encoded, _ := ab09aq09json.Marshal(ab09aq09MakeInput("periods"))
	ab09aq09Original := append([]byte{}, ab09aq09Encoded...)
	ab09aq09First, ab09aq09Error := Ab09Aq09Evaluate(ab09aq09Encoded)
	if ab09aq09Error != nil {
		ab09aq09T.Fatal(ab09aq09Error)
	}
	ab09aq09Saved := append([]byte{}, ab09aq09First...)
	for ab09aq09Index := range ab09aq09Encoded {
		ab09aq09Encoded[ab09aq09Index] = ' '
	}
	if !ab09aq09bytes.Equal(ab09aq09Saved, ab09aq09First) {
		ab09aq09T.Fatal("input aliases output")
	}
	for ab09aq09Index := range ab09aq09First {
		ab09aq09First[ab09aq09Index] = ' '
	}
	ab09aq09Second, ab09aq09Error := Ab09Aq09Evaluate(ab09aq09Original)
	if ab09aq09Error != nil || !ab09aq09bytes.Equal(ab09aq09Saved, ab09aq09Second) {
		ab09aq09T.Fatal("calls share storage")
	}
	if ab09aq09strings.Contains(string(ab09aq09Second), "qualified\":true") {
		ab09aq09T.Fatal("qualified")
	}
}

func TestAb09Aq09FieldAxes(ab09aq09T *ab09aq09testing.T) {
	for _, ab09aq09Key := range []string{"period0", "period1", "period2"} {
		for _, ab09aq09Value := range []int{0, 1, 54, 144, 145, 255} {
			ab09aq09Input := ab09aq09MakeInput("periods")
			ab09aq09Input["fields"].(map[string]any)[ab09aq09Key] = []int{ab09aq09Value}
			ab09aq09Output := ab09aq09Run(ab09aq09T, ab09aq09Input)
			ab09aq09Baseline := ab09aq09Run(ab09aq09T, ab09aq09MakeInput("periods"))
			for _, ab09aq09Sibling := range []string{"period0", "period1", "period2"} {
				ab09aq09Result := ab09aq09GetResult(ab09aq09Output, ab09aq09Sibling)
				if ab09aq09Sibling != ab09aq09Key {
					if !ab09aq09reflect.DeepEqual(ab09aq09Result, ab09aq09GetResult(ab09aq09Baseline, ab09aq09Sibling)) {
						ab09aq09T.Fatal("field cross-talk")
					}
				} else {
					if !ab09aq09reflect.DeepEqual(ab09aq09Result["raw"], []any{float64(ab09aq09Value)}) {
						ab09aq09T.Fatal("raw field lost")
					}
					if ab09aq09Value <= 144 {
						if ab09aq09Result["number"] != float64(ab09aq09Value*10) {
							ab09aq09T.Fatal("field ignored")
						}
					} else {
						ab09aq09NoNumber(ab09aq09T, ab09aq09Result)
					}
				}
			}
		}
	}
	for _, ab09aq09Key := range []string{"weekday", "channel", "mode"} {
		for _, ab09aq09Value := range [][]int{{}, {0}, {1}, {255}, {0, 255}} {
			ab09aq09Input := ab09aq09MakeInput("periods")
			ab09aq09Input["fields"].(map[string]any)[ab09aq09Key] = ab09aq09Value
			ab09aq09Output := ab09aq09Run(ab09aq09T, ab09aq09Input)
			ab09aq09Result := ab09aq09GetResult(ab09aq09Output, ab09aq09Key)
			if ab09aq09Result["state"] != "Unknown" {
				ab09aq09T.Fatal("opaque carrier decoded")
			}
			ab09aq09NoNumber(ab09aq09T, ab09aq09Result)
		}
	}
}
