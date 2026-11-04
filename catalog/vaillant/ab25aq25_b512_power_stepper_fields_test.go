package vaillant

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func ab25aq25Call(ab25aq25T *testing.T, ab25aq25Input any) map[string]any {
	ab25aq25T.Helper()
	ab25aq25Bytes, ab25aq25Err := json.Marshal(ab25aq25Input)
	if ab25aq25Err != nil {
		ab25aq25T.Fatal(ab25aq25Err)
	}
	ab25aq25Output, ab25aq25Err := Ab25Aq25Evaluate(ab25aq25Bytes)
	if ab25aq25Err != nil {
		ab25aq25T.Fatal(ab25aq25Err)
	}
	ab25aq25Object := map[string]any{}
	if ab25aq25Err = json.Unmarshal(ab25aq25Output, &ab25aq25Object); ab25aq25Err != nil {
		ab25aq25T.Fatal(ab25aq25Err)
	}
	return ab25aq25Object
}

func ab25aq25Contains(ab25aq25T *testing.T, ab25aq25Actual, ab25aq25Expected any) {
	ab25aq25T.Helper()
	if ab25aq25Object, ab25aq25OK := ab25aq25Expected.(map[string]any); ab25aq25OK {
		ab25aq25Got, ab25aq25OK := ab25aq25Actual.(map[string]any)
		if !ab25aq25OK {
			ab25aq25T.Fatalf("expected object, got %v", ab25aq25Actual)
		}
		for ab25aq25Key, ab25aq25Value := range ab25aq25Object {
			ab25aq25Contains(ab25aq25T, ab25aq25Got[ab25aq25Key], ab25aq25Value)
		}
	} else if !reflect.DeepEqual(ab25aq25Actual, ab25aq25Expected) {
		ab25aq25T.Fatalf("got %#v want %#v", ab25aq25Actual, ab25aq25Expected)
	}
}

func ab25aq25Input() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 18, "operation": 8, "direction": "request", "view": "power-stepper", "family": "synthetic", "variant": "candidate"}, "fields": map[string]any{"inlet": []int{240, 255}, "status": []int{85}, "stepper": []int{1}}}
}

func TestAb25Aq25Fixture(ab25aq25T *testing.T) {
	ab25aq25Bytes, ab25aq25Err := os.ReadFile("testdata/ab25aq25_b512_power_stepper_fields.json")
	if ab25aq25Err != nil {
		ab25aq25T.Fatal(ab25aq25Err)
	}
	var ab25aq25Cases []map[string]any
	if ab25aq25Err = json.Unmarshal(ab25aq25Bytes, &ab25aq25Cases); ab25aq25Err != nil {
		ab25aq25T.Fatal(ab25aq25Err)
	}
	if len(ab25aq25Cases) == 0 {
		ab25aq25T.Fatal("zero fixtures")
	}
	for _, ab25aq25Case := range ab25aq25Cases {
		ab25aq25T.Run(ab25aq25Case["id"].(string), func(ab25aq25T *testing.T) {
			ab25aq25Got := ab25aq25Call(ab25aq25T, ab25aq25Case["input"])
			ab25aq25Contains(ab25aq25T, ab25aq25Got, ab25aq25Case["expected"])
			if ab25aq25Got["physical"] != "Unknown" || ab25aq25Got["qualified"] != false {
				ab25aq25T.Fatal(ab25aq25Got)
			}
			ab25aq25Results := ab25aq25Got["results"].(map[string]any)
			for _, ab25aq25Result := range ab25aq25Results {
				ab25aq25Map := ab25aq25Result.(map[string]any)
				if ab25aq25Map["state"] == "complete" && ab25aq25Map["grade"] != "Hypothesis" {
					ab25aq25T.Fatal(ab25aq25Map)
				}
				if ab25aq25Flags, ab25aq25OK := ab25aq25Map["flags"].(map[string]any); ab25aq25OK {
					if len(ab25aq25Flags) != 1 {
						ab25aq25T.Fatal("release must not imply movement", ab25aq25Flags)
					}
					if _, ab25aq25OK = ab25aq25Flags["release"]; !ab25aq25OK {
						ab25aq25T.Fatal(ab25aq25Flags)
					}
					for _, ab25aq25Key := range []string{"number", "label", "movement", "contact"} {
						if _, ab25aq25OK = ab25aq25Map[ab25aq25Key]; ab25aq25OK {
							ab25aq25T.Fatal(ab25aq25Map)
						}
					}
				}
			}
		})
	}
}

func TestAb25Aq25Selectors(ab25aq25T *testing.T) {
	for ab25aq25Key, ab25aq25Values := range map[string][]any{"primary": {180, nil}, "secondary": {19, nil}, "operation": {0, 10, 20, nil}, "direction": {"positive", "other", nil}, "view": {"other", nil}, "family": {"other", nil, 0, false}, "variant": {"other", nil, 0, false}} {
		for ab25aq25Index, ab25aq25Value := range ab25aq25Values {
			ab25aq25T.Run(fmt.Sprintf("%s_%d", ab25aq25Key, ab25aq25Index), func(ab25aq25T *testing.T) {
				ab25aq25In := ab25aq25Input()
				ab25aq25Context := ab25aq25In["context"].(map[string]any)
				if ab25aq25Value == nil {
					delete(ab25aq25Context, ab25aq25Key)
				} else {
					ab25aq25Context[ab25aq25Key] = ab25aq25Value
				}
				ab25aq25Got := ab25aq25Call(ab25aq25T, ab25aq25In)
				if ab25aq25Got["qualified"] != false {
					ab25aq25T.Fatal(ab25aq25Got)
				}
				if ab25aq25Key == "family" || ab25aq25Key == "variant" {
					if ab25aq25Got["state"] != "complete" {
						ab25aq25T.Fatal(ab25aq25Got)
					}
				} else if ab25aq25Key == "direction" && ab25aq25Value == "positive" {
					if ab25aq25Got["state"] != "partial" {
						ab25aq25T.Fatal(ab25aq25Got)
					}
					if ab25aq25Got["results"].(map[string]any)["status"].(map[string]any)["state"] != "Unknown" {
						ab25aq25T.Fatal(ab25aq25Got)
					}
				} else {
					if ab25aq25Got["state"] != "Unknown" {
						ab25aq25T.Fatal(ab25aq25Got)
					}
					for _, ab25aq25Result := range ab25aq25Got["results"].(map[string]any) {
						if len(ab25aq25Result.(map[string]any)) != 2 {
							ab25aq25T.Fatal(ab25aq25Result)
						}
					}
				}
				ab25aq25Bytes, _ := json.Marshal(ab25aq25Context)
				ab25aq25Expected := map[string]any{}
				if ab25aq25Err := json.Unmarshal(ab25aq25Bytes, &ab25aq25Expected); ab25aq25Err != nil {
					ab25aq25T.Fatal(ab25aq25Err)
				}
				if !reflect.DeepEqual(ab25aq25Got["context"], ab25aq25Expected) {
					ab25aq25T.Fatal("context retention")
				}
			})
		}
	}
}

func TestAb25Aq25Invalid(ab25aq25T *testing.T) {
	for _, ab25aq25Text := range []string{``, `[]`, `{`, `{} {}`, `{"context":{},"fields":{},"extra":1e999}`, `{"context":{},"fields":{"x":[],"x":[]}}`, `{"context":{"operation":8,"operation":8},"fields":{}}`, `{"context":{},"fields":{},"arguments":{"x":0,"x":1}}`, `{"context":{},"fields":{},"x":{"a":1,"a":2}}`, `{"context":{},"fields":{},"x":[{"a":1,"a":2}]}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{},"fields":{"a":[NaN]}}`} {
		ab25aq25Out, ab25aq25Err := Ab25Aq25Evaluate(json.RawMessage(ab25aq25Text))
		if ab25aq25Err == nil || ab25aq25Out != nil {
			ab25aq25T.Fatalf("accepted %s", ab25aq25Text)
		}
	}
	for _, ab25aq25Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, ab25aq25Value := range []any{nil, true, []int{}, map[string]any{}, 1.5, -1, 256, "8"} {
			if (ab25aq25Key == "direction" || ab25aq25Key == "view") && ab25aq25Value == "8" {
				continue
			}
			ab25aq25In := ab25aq25Input()
			ab25aq25In["context"].(map[string]any)[ab25aq25Key] = ab25aq25Value
			ab25aq25Bytes, _ := json.Marshal(ab25aq25In)
			ab25aq25Out, ab25aq25Err := Ab25Aq25Evaluate(ab25aq25Bytes)
			if ab25aq25Err == nil || ab25aq25Out != nil {
				ab25aq25T.Fatalf("accepted %s", ab25aq25Bytes)
			}
		}
	}
	for _, ab25aq25Value := range []any{nil, true, "1", 1, []any{nil}, []any{true}, []any{"1"}, []any{1.5}, []any{-1}, []any{256}} {
		ab25aq25In := ab25aq25Input()
		ab25aq25In["fields"].(map[string]any)["inlet"] = ab25aq25Value
		ab25aq25Bytes, _ := json.Marshal(ab25aq25In)
		ab25aq25Out, ab25aq25Err := Ab25Aq25Evaluate(ab25aq25Bytes)
		if ab25aq25Err == nil || ab25aq25Out != nil {
			ab25aq25T.Fatalf("accepted %s", ab25aq25Bytes)
		}
	}
	for _, ab25aq25Key := range []string{"context", "fields", "arguments"} {
		for _, ab25aq25Value := range []any{nil, 0, false, "", []any{}} {
			ab25aq25In := ab25aq25Input()
			ab25aq25In[ab25aq25Key] = ab25aq25Value
			ab25aq25Bytes, _ := json.Marshal(ab25aq25In)
			ab25aq25Out, ab25aq25Err := Ab25Aq25Evaluate(ab25aq25Bytes)
			if ab25aq25Err == nil || ab25aq25Out != nil {
				ab25aq25T.Fatalf("accepted %s", ab25aq25Bytes)
			}
		}
	}
}

func TestAb25Aq25WidthsAndRetention(ab25aq25T *testing.T) {
	for _, ab25aq25Spec := range []struct {
		ab25aq25Op        int
		ab25aq25Direction string
		ab25aq25Field     string
		ab25aq25Width     int
	}{{7, "positive", "maximum_power", 1}, {7, "positive", "minimum_power", 1}, {8, "request", "inlet", 2}, {8, "request", "status", 1}, {8, "request", "stepper", 1}, {8, "positive", "setpoint", 2}, {8, "positive", "stepper", 1}, {9, "request", "release", 1}, {9, "positive", "release", 1}} {
		for _, ab25aq25Width := range []int{-1, 0, ab25aq25Spec.ab25aq25Width - 1, ab25aq25Spec.ab25aq25Width, ab25aq25Spec.ab25aq25Width + 1} {
			ab25aq25T.Run(fmt.Sprintf("%d_%s_%s_%d", ab25aq25Spec.ab25aq25Op, ab25aq25Spec.ab25aq25Direction, ab25aq25Spec.ab25aq25Field, ab25aq25Width), func(ab25aq25T *testing.T) {
				ab25aq25In := ab25aq25Input()
				ab25aq25Context := ab25aq25In["context"].(map[string]any)
				ab25aq25Context["operation"] = ab25aq25Spec.ab25aq25Op
				ab25aq25Context["direction"] = ab25aq25Spec.ab25aq25Direction
				ab25aq25Fields := map[string]any{"maximum_power": []int{255}, "minimum_power": []int{250}, "inlet": []int{240, 255}, "status": []int{85}, "stepper": []int{1}, "setpoint": []int{16, 0}, "release": []int{1}, "opaque": []int{}}
				ab25aq25In["fields"] = ab25aq25Fields
				if ab25aq25Width < 0 {
					delete(ab25aq25Fields, ab25aq25Spec.ab25aq25Field)
				} else {
					ab25aq25Fields[ab25aq25Spec.ab25aq25Field] = make([]int, ab25aq25Width)
				}
				ab25aq25In["extra"] = map[string]any{"null": nil, "zero": 0, "false": false, "empty": []any{}}
				ab25aq25In["arguments"] = map[string]any{"opaque": nil}
				ab25aq25Got := ab25aq25Call(ab25aq25T, ab25aq25In)
				ab25aq25Results := ab25aq25Got["results"].(map[string]any)
				ab25aq25Result := ab25aq25Results[ab25aq25Spec.ab25aq25Field].(map[string]any)
				if ab25aq25Width == ab25aq25Spec.ab25aq25Width {
					if ab25aq25Spec.ab25aq25Field == "maximum_power" || ab25aq25Spec.ab25aq25Field == "minimum_power" || ab25aq25Spec.ab25aq25Field == "inlet" || ab25aq25Spec.ab25aq25Field == "setpoint" {
						if ab25aq25Result["number"] != float64(0) {
							ab25aq25T.Fatal("zero arithmetic", ab25aq25Result)
						}
					}
					if ab25aq25Got["state"] != "complete" || ab25aq25Result["state"] != "complete" {
						ab25aq25T.Fatal(ab25aq25Got)
					}
				} else {
					if ab25aq25Got["state"] != "partial" || ab25aq25Result["state"] != "incomplete" {
						ab25aq25T.Fatal(ab25aq25Got)
					}
					if _, ab25aq25OK := ab25aq25Result["number"]; ab25aq25OK {
						ab25aq25T.Fatal(ab25aq25Result)
					}
				}
				_, ab25aq25Present := ab25aq25Result["raw"]
				if ab25aq25Present != (ab25aq25Width >= 0) {
					ab25aq25T.Fatal(ab25aq25Result)
				}
				ab25aq25CanonicalBytes, _ := json.Marshal(ab25aq25In)
				ab25aq25Canonical := map[string]any{}
				if ab25aq25Err := json.Unmarshal(ab25aq25CanonicalBytes, &ab25aq25Canonical); ab25aq25Err != nil {
					ab25aq25T.Fatal(ab25aq25Err)
				}
				for _, ab25aq25Key := range []string{"context", "fields", "extra", "arguments"} {
					if !reflect.DeepEqual(ab25aq25Got[ab25aq25Key], ab25aq25Canonical[ab25aq25Key]) {
						ab25aq25T.Fatal("retention", ab25aq25Key)
					}
				}
				if ab25aq25Spec.ab25aq25Op == 7 && ab25aq25Spec.ab25aq25Field == "maximum_power" {
					ab25aq25Contains(ab25aq25T, ab25aq25Results["minimum_power"], map[string]any{"number": float64(125)})
				}
				if ab25aq25Spec.ab25aq25Op == 8 && ab25aq25Spec.ab25aq25Field != "stepper" {
					ab25aq25Contains(ab25aq25T, ab25aq25Results["stepper"], map[string]any{"state": "complete"})
				}
			})
		}
	}
}

func TestAb25Aq25Memory(ab25aq25T *testing.T) {
	ab25aq25InputBytes, _ := json.Marshal(ab25aq25Input())
	ab25aq25First, ab25aq25Err := Ab25Aq25Evaluate(ab25aq25InputBytes)
	if ab25aq25Err != nil {
		ab25aq25T.Fatal(ab25aq25Err)
	}
	ab25aq25Saved := string(ab25aq25First)
	for ab25aq25Index := range ab25aq25InputBytes {
		ab25aq25InputBytes[ab25aq25Index] = ' '
	}
	if string(ab25aq25First) != ab25aq25Saved {
		ab25aq25T.Fatal("input alias")
	}
	for ab25aq25Index := range ab25aq25First {
		ab25aq25First[ab25aq25Index] = ' '
	}
	ab25aq25NextBytes, _ := json.Marshal(ab25aq25Input())
	ab25aq25Next, ab25aq25Err := Ab25Aq25Evaluate(ab25aq25NextBytes)
	if ab25aq25Err != nil || string(ab25aq25Next) != ab25aq25Saved {
		ab25aq25T.Fatal("cross-call alias")
	}
	if strings.Contains(string(ab25aq25Next), "movement") {
		ab25aq25T.Fatal("invented movement")
	}
}

func TestAb25Aq25OpaqueIdentity(ab25aq25T *testing.T) {
	for _, ab25aq25Key := range []string{"family", "variant"} {
		for _, ab25aq25Value := range []any{nil, false, 0, "", []any{}, map[string]any{}} {
			ab25aq25In := ab25aq25Input()
			ab25aq25In["context"].(map[string]any)[ab25aq25Key] = ab25aq25Value
			ab25aq25Got := ab25aq25Call(ab25aq25T, ab25aq25In)
			ab25aq25Bytes, _ := json.Marshal(ab25aq25In["context"])
			ab25aq25Expected := map[string]any{}
			if ab25aq25Err := json.Unmarshal(ab25aq25Bytes, &ab25aq25Expected); ab25aq25Err != nil {
				ab25aq25T.Fatal(ab25aq25Err)
			}
			if !reflect.DeepEqual(ab25aq25Got["context"], ab25aq25Expected) || ab25aq25Got["qualified"] != false || ab25aq25Got["state"] != "complete" {
				ab25aq25T.Fatal(ab25aq25Got)
			}
		}
	}
}
