package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func ab03ao03Input(ab03ao03Key any, ab03ao03Value any) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 6, "operation": 17, "direction": "request", "view": "parameter-union", "family": map[string]any{"opaque": []any{"synthetic"}}, "variant": []any{0, "opaque"}}, "fields": map[string]any{"key": ab03ao03Key, "value": ab03ao03Value, "tail": []int{0, 255}}, "arguments": map[string]any{}}
}

func ab03ao03Run(ab03ao03T *testing.T, ab03ao03Input any) map[string]any {
	ab03ao03T.Helper()
	ab03ao03JSON, ab03ao03Err := json.Marshal(ab03ao03Input)
	if ab03ao03Err != nil {
		ab03ao03T.Fatal(ab03ao03Err)
	}
	ab03ao03Result, ab03ao03Err := Ab03Ao03Interpret(ab03ao03JSON)
	if ab03ao03Err != nil {
		ab03ao03T.Fatal(ab03ao03Err)
	}
	ab03ao03Output := map[string]any{}
	if ab03ao03Err := json.Unmarshal(ab03ao03Result, &ab03ao03Output); ab03ao03Err != nil {
		ab03ao03T.Fatal(ab03ao03Err)
	}
	ab03ao03Canonical := map[string]any{}
	if ab03ao03Err := json.Unmarshal(ab03ao03JSON, &ab03ao03Canonical); ab03ao03Err != nil {
		ab03ao03T.Fatal(ab03ao03Err)
	}
	if !reflect.DeepEqual(ab03ao03Output["context"], ab03ao03Canonical["context"]) || !reflect.DeepEqual(ab03ao03Output["raw"], ab03ao03Canonical["fields"]) {
		ab03ao03T.Fatalf("evidence changed: %s", ab03ao03Result)
	}
	if ab03ao03Arguments, ab03ao03Present := ab03ao03Canonical["arguments"]; ab03ao03Present {
		if !reflect.DeepEqual(ab03ao03Arguments, ab03ao03Output["arguments"]) {
			ab03ao03T.Fatal("arguments changed")
		}
	} else if _, ab03ao03Present := ab03ao03Output["arguments"]; ab03ao03Present {
		ab03ao03T.Fatal("absent arguments fabricated")
	}
	if ab03ao03Output["physical"] != "Unknown" || ab03ao03Output["qualified"] != false {
		ab03ao03T.Fatal("physical qualification fabricated")
	}
	for ab03ao03Name, ab03ao03Semantic := range ab03ao03Output["fields"].(map[string]any) {
		ab03ao03Raw, ab03ao03HasRaw := ab03ao03Semantic.(map[string]any)["raw"]
		ab03ao03Supplied, ab03ao03Present := ab03ao03Canonical["fields"].(map[string]any)[ab03ao03Name]
		if ab03ao03HasRaw != ab03ao03Present || (ab03ao03HasRaw && !reflect.DeepEqual(ab03ao03Raw, ab03ao03Supplied)) {
			ab03ao03T.Fatalf("semantic raw changed: %s", ab03ao03Name)
		}
	}
	return ab03ao03Output
}

func ab03ao03ValueField(ab03ao03Output map[string]any) map[string]any {
	return ab03ao03Output["fields"].(map[string]any)["value"].(map[string]any)
}

func TestAb03Ao03Alternatives(ab03ao03T *testing.T) {
	ab03ao03Data, ab03ao03Err := os.ReadFile("testdata/ab03ao03_b506_parameter_views.json")
	if ab03ao03Err != nil {
		ab03ao03T.Fatal(ab03ao03Err)
	}
	var ab03ao03Cases []struct {
		Ab03Ao03ID    string `json:"ab03ao03ID"`
		Ab03Ao03Key   int    `json:"ab03ao03Key"`
		Ab03Ao03Value []int  `json:"ab03ao03Value"`
		Ab03Ao03State string `json:"ab03ao03State"`
	}
	if ab03ao03Err := json.Unmarshal(ab03ao03Data, &ab03ao03Cases); ab03ao03Err != nil {
		ab03ao03T.Fatal(ab03ao03Err)
	}
	if len(ab03ao03Cases) != 22 {
		ab03ao03T.Fatal("missing cases")
	}
	for ab03ao03Pass := 0; ab03ao03Pass < 2; ab03ao03Pass++ {
		for ab03ao03Index := range ab03ao03Cases {
			ab03ao03Case := ab03ao03Cases[(ab03ao03Index*7+ab03ao03Pass)%len(ab03ao03Cases)]
			ab03ao03T.Run(fmt.Sprintf("%s-permutation-%d", ab03ao03Case.Ab03Ao03ID, ab03ao03Pass), func(ab03ao03T *testing.T) {
				ab03ao03Output := ab03ao03Run(ab03ao03T, ab03ao03Input([]int{ab03ao03Case.Ab03Ao03Key}, ab03ao03Case.Ab03Ao03Value))
				if ab03ao03Output["state"] != ab03ao03Case.Ab03Ao03State {
					ab03ao03T.Fatalf("state: %v", ab03ao03Output)
				}
				if ab03ao03Case.Ab03Ao03State == "Unknown" {
					if len(ab03ao03Output["fields"].(map[string]any)) != 0 {
						ab03ao03T.Fatal("unknown key selected")
					}
					return
				}
				ab03ao03KeyField := ab03ao03Output["fields"].(map[string]any)["key"].(map[string]any)
				if ab03ao03KeyField["state"] != "Hypothesis" || ab03ao03KeyField["number"] != float64(ab03ao03Case.Ab03Ao03Key) {
					ab03ao03T.Fatal("key lost")
				}
				ab03ao03Field := ab03ao03ValueField(ab03ao03Output)
				switch ab03ao03Case.Ab03Ao03Key {
				case 0, 1, 2, 4, 5:
					if len(ab03ao03Field) != 2 || ab03ao03Field["state"] != "Hypothesis" {
						ab03ao03T.Fatal("opaque bytes interpreted")
					}
					if _, ab03ao03OK := ab03ao03Field["raw"].([]any); !ab03ao03OK {
						ab03ao03T.Fatal("opaque bytes converted to text")
					}
				case 10:
					if ab03ao03Case.Ab03Ao03Value[0] == 2 {
						if ab03ao03Field["state"] != "unmapped" || len(ab03ao03Field) != 2 {
							ab03ao03T.Fatal("Boolean 2 mapped")
						}
					} else {
						ab03ao03Label := "no"
						if ab03ao03Case.Ab03Ao03Value[0] == 1 {
							ab03ao03Label = "yes"
						}
						if ab03ao03Field["label"] != ab03ao03Label || ab03ao03Field["number"] != float64(ab03ao03Case.Ab03Ao03Value[0]) || ab03ao03Field["state"] != "Hypothesis" {
							ab03ao03T.Fatal("Boolean mismatch")
						}
					}
				default:
					if ab03ao03Field["number"] != float64(127) || ab03ao03Field["state"] != "Hypothesis" {
						ab03ao03T.Fatal("scalar mismatch")
					}
				}
			})
		}
	}
}

func TestAb03Ao03WidthsAndPresence(ab03ao03T *testing.T) {
	for _, ab03ao03Key := range []int{0, 1, 2, 4, 5, 3, 6, 7, 8, 9, 10, 32, 33, 34, 35, 36, 37} {
		for _, ab03ao03Width := range []int{-1, 0, 1, 2, 3} {
			ab03ao03T.Run(fmt.Sprintf("key-%d-width-%d", ab03ao03Key, ab03ao03Width), func(ab03ao03T *testing.T) {
				ab03ao03Bytes := []int{}
				for ab03ao03Index := 0; ab03ao03Index < ab03ao03Width; ab03ao03Index++ {
					ab03ao03Bytes = append(ab03ao03Bytes, 0)
				}
				ab03ao03In := ab03ao03Input([]int{ab03ao03Key}, ab03ao03Bytes)
				if ab03ao03Width < 0 {
					delete(ab03ao03In["fields"].(map[string]any), "value")
				}
				ab03ao03Output := ab03ao03Run(ab03ao03T, ab03ao03In)
				ab03ao03Field := ab03ao03ValueField(ab03ao03Output)
				ab03ao03Opaque := ab03ao03Key == 0 || ab03ao03Key == 1 || ab03ao03Key == 2 || ab03ao03Key == 4 || ab03ao03Key == 5
				ab03ao03Expected := "Hypothesis"
				if ab03ao03Width < 0 || (!ab03ao03Opaque && ab03ao03Width == 0) {
					ab03ao03Expected = "incomplete"
				}
				if !ab03ao03Opaque && ab03ao03Width > 1 {
					ab03ao03Expected = "width-conflict"
				}
				if ab03ao03Field["state"] != ab03ao03Expected {
					ab03ao03T.Fatalf("width state: %v", ab03ao03Field)
				}
				_, ab03ao03HasRaw := ab03ao03Field["raw"]
				if ab03ao03HasRaw != (ab03ao03Width >= 0) {
					ab03ao03T.Fatal("missing/empty conflated")
				}
				_, ab03ao03HasNumber := ab03ao03Field["number"]
				if ab03ao03HasNumber != (!ab03ao03Opaque && ab03ao03Width == 1) {
					ab03ao03T.Fatal("invalid numeric presence")
				}
				if ab03ao03Expected != "Hypothesis" {
					if ab03ao03Output["state"] != "partial" || ab03ao03Output["fields"].(map[string]any)["key"].(map[string]any)["number"] != float64(ab03ao03Key) {
						ab03ao03T.Fatal("complete independent key erased")
					}
				} else if ab03ao03Output["state"] != "Hypothesis" {
					ab03ao03T.Fatal("complete state")
				}
			})
		}
	}
	for _, ab03ao03Key := range []any{nil, []int{}, []int{2, 3}} {
		ab03ao03In := ab03ao03Input(ab03ao03Key, []int{0})
		if ab03ao03Key == nil {
			delete(ab03ao03In["fields"].(map[string]any), "key")
		}
		ab03ao03Output := ab03ao03Run(ab03ao03T, ab03ao03In)
		if ab03ao03Output["state"] != "partial" {
			ab03ao03T.Fatal("invalid key width selected")
		}
	}
}

func TestAb03Ao03ContextAxes(ab03ao03T *testing.T) {
	for ab03ao03Axis, ab03ao03Value := range map[string]any{"primary": 180, "secondary": 7, "operation": 18, "direction": "response", "view": "other"} {
		for _, ab03ao03Missing := range []bool{false, true} {
			ab03ao03T.Run(fmt.Sprintf("%s-missing-%v", ab03ao03Axis, ab03ao03Missing), func(ab03ao03T *testing.T) {
				ab03ao03In := ab03ao03Input([]int{3}, []int{127})
				ab03ao03Context := ab03ao03In["context"].(map[string]any)
				if ab03ao03Missing {
					delete(ab03ao03Context, ab03ao03Axis)
				} else {
					ab03ao03Context[ab03ao03Axis] = ab03ao03Value
				}
				ab03ao03Output := ab03ao03Run(ab03ao03T, ab03ao03In)
				if ab03ao03Output["state"] != "Unknown" || len(ab03ao03Output["fields"].(map[string]any)) != 0 {
					ab03ao03T.Fatal("selector guessed")
				}
			})
		}
	}
	for _, ab03ao03Axis := range []string{"family", "variant"} {
		for _, ab03ao03Value := range []any{nil, 0, "parameter-union", []any{}, map[string]any{"different": true}} {
			ab03ao03In := ab03ao03Input([]int{3}, []int{127})
			ab03ao03In["context"].(map[string]any)[ab03ao03Axis] = ab03ao03Value
			if ab03ao03Run(ab03ao03T, ab03ao03In)["state"] != "Hypothesis" {
				ab03ao03T.Fatal("opaque context selects")
			}
		}
	}
	for _, ab03ao03Tail := range [][]int{{}, {1}, {255, 0, 127}} {
		ab03ao03In := ab03ao03Input([]int{3}, []int{0})
		ab03ao03In["fields"].(map[string]any)["tail"] = ab03ao03Tail
		delete(ab03ao03In, "arguments")
		if ab03ao03ValueField(ab03ao03Run(ab03ao03T, ab03ao03In))["number"] != float64(0) {
			ab03ao03T.Fatal("tail changed scalar")
		}
	}
}

func TestAb03Ao03Reject(ab03ao03T *testing.T) {
	for _, ab03ao03JSON := range []string{``, `{`, `{} {}`, `[]`, `null`, `{"context":{},"fields":{},"fields":{}}`, `{"context":{"x":1,"x":2},"fields":{}}`, `{"context":{"x":[{"a":0,"a":1}]},"fields":{}}`, `{"context":{},"fields":null}`} {
		if ab03ao03Result, ab03ao03Err := Ab03Ao03Interpret(json.RawMessage(ab03ao03JSON)); ab03ao03Err == nil || ab03ao03Result != nil {
			ab03ao03T.Fatalf("accepted %q", ab03ao03JSON)
		}
	}
	for _, ab03ao03Axis := range []string{"key", "value", "tail"} {
		for _, ab03ao03Bad := range []any{nil, "bytes", 0, map[string]any{}, []any{-1}, []any{256}, []any{1.5}, []any{true}, []any{"0"}, []any{nil}, []any{[]int{0}}} {
			ab03ao03In := ab03ao03Input([]int{3}, []int{0})
			ab03ao03In["fields"].(map[string]any)[ab03ao03Axis] = ab03ao03Bad
			ab03ao03JSON, _ := json.Marshal(ab03ao03In)
			if ab03ao03Result, ab03ao03Err := Ab03Ao03Interpret(ab03ao03JSON); ab03ao03Err == nil || ab03ao03Result != nil {
				ab03ao03T.Fatalf("accepted %s", ab03ao03JSON)
			}
		}
	}
	for _, ab03ao03Axis := range []string{"primary", "secondary", "operation", "direction", "view"} {
		ab03ao03In := ab03ao03Input([]int{3}, []int{0})
		ab03ao03In["context"].(map[string]any)[ab03ao03Axis] = []any{}
		ab03ao03JSON, _ := json.Marshal(ab03ao03In)
		if _, ab03ao03Err := Ab03Ao03Interpret(ab03ao03JSON); ab03ao03Err == nil {
			ab03ao03T.Fatal("wrong selector kind")
		}
	}
	for _, ab03ao03Axis := range []string{"context", "fields", "arguments"} {
		ab03ao03In := ab03ao03Input([]int{3}, []int{0})
		ab03ao03In[ab03ao03Axis] = []any{}
		ab03ao03JSON, _ := json.Marshal(ab03ao03In)
		if _, ab03ao03Err := Ab03Ao03Interpret(ab03ao03JSON); ab03ao03Err == nil {
			ab03ao03T.Fatal("wrong object kind")
		}
	}
	for _, ab03ao03Arguments := range []any{nil, map[string]any{"flag": true}, map[string]any{"count": 0}} {
		ab03ao03In := ab03ao03Input([]int{3}, []int{0})
		ab03ao03In["arguments"] = ab03ao03Arguments
		ab03ao03JSON, _ := json.Marshal(ab03ao03In)
		if _, ab03ao03Err := Ab03Ao03Interpret(ab03ao03JSON); ab03ao03Err == nil {
			ab03ao03T.Fatal("undefined arguments")
		}
	}
}

func TestAb03Ao03Memory(ab03ao03T *testing.T) {
	ab03ao03InputBytes, _ := json.Marshal(ab03ao03Input([]int{2}, []int{0, 170, 255}))
	ab03ao03Original := append([]byte(nil), ab03ao03InputBytes...)
	ab03ao03First, ab03ao03Err := Ab03Ao03Interpret(ab03ao03InputBytes)
	if ab03ao03Err != nil {
		ab03ao03T.Fatal(ab03ao03Err)
	}
	ab03ao03Saved := append([]byte(nil), ab03ao03First...)
	for ab03ao03Index := range ab03ao03InputBytes {
		ab03ao03InputBytes[ab03ao03Index] = ' '
	}
	if !bytes.Equal(ab03ao03First, ab03ao03Saved) {
		ab03ao03T.Fatal("input aliases output")
	}
	ab03ao03Second, ab03ao03Err := Ab03Ao03Interpret(ab03ao03Original)
	if ab03ao03Err != nil {
		ab03ao03T.Fatal(ab03ao03Err)
	}
	ab03ao03Nested := map[string]any{}
	if ab03ao03Err := json.Unmarshal(ab03ao03First, &ab03ao03Nested); ab03ao03Err != nil {
		ab03ao03T.Fatal(ab03ao03Err)
	}
	ab03ao03Nested["raw"].(map[string]any)["value"].([]any)[0] = float64(99)
	ab03ao03Nested["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "changed"
	ab03ao03ValueField(ab03ao03Nested)["raw"].([]any)[1] = float64(99)
	for ab03ao03Index := range ab03ao03First {
		ab03ao03First[ab03ao03Index] = ' '
	}
	ab03ao03Third, ab03ao03Err := Ab03Ao03Interpret(ab03ao03Original)
	if ab03ao03Err != nil || !bytes.Equal(ab03ao03Second, ab03ao03Saved) || !bytes.Equal(ab03ao03Third, ab03ao03Saved) {
		ab03ao03T.Fatal("result aliases another result or later call")
	}
}

func TestAb03Ao03KeyVersusValue(ab03ao03T *testing.T) {
	for _, ab03ao03Key := range []int{2, 3} {
		ab03ao03Output := ab03ao03Run(ab03ao03T, ab03ao03Input([]int{ab03ao03Key}, []int{0, 170, 255}))
		ab03ao03Expected := "Hypothesis"
		if ab03ao03Key == 3 {
			ab03ao03Expected = "width-conflict"
		}
		if ab03ao03ValueField(ab03ao03Output)["state"] != ab03ao03Expected {
			ab03ao03T.Fatal("key did not select independent alternative")
		}
	}
	for _, ab03ao03Value := range []int{0, 1, 127, 255} {
		if ab03ao03ValueField(ab03ao03Run(ab03ao03T, ab03ao03Input([]int{3}, []int{ab03ao03Value})))["number"] != float64(ab03ao03Value) {
			ab03ao03T.Fatal("scalar ignored byte")
		}
	}
	for _, ab03ao03Axis := range []string{"primary", "secondary", "operation"} {
		for _, ab03ao03Value := range []any{0, -1, json.Number("999999999999999999999999999999999999")} {
			ab03ao03In := ab03ao03Input([]int{3}, []int{127})
			ab03ao03In["context"].(map[string]any)[ab03ao03Axis] = ab03ao03Value
			if ab03ao03Run(ab03ao03T, ab03ao03In)["state"] != "Unknown" {
				ab03ao03T.Fatal("unknown integer selector admitted")
			}
		}
	}
}
