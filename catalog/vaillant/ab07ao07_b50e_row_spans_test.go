package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func ab07ao07Input(ab07ao07Key int, ab07ao07Value []int) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 14, "operation": 2, "direction": "positive", "view": "row-spans", "family": map[string]any{"opaque": []any{"synthetic", 0}}, "variant": "candidate"}, "fields": map[string]any{"key": []int{ab07ao07Key}, "value": ab07ao07Value, "unmodelled": []int{99, 0, 255}}}
}

func ab07ao07Call(ab07ao07T *testing.T, ab07ao07Object map[string]any) map[string]any {
	ab07ao07T.Helper()
	ab07ao07InputBytes, ab07ao07Err := json.Marshal(ab07ao07Object)
	if ab07ao07Err != nil {
		ab07ao07T.Fatal(ab07ao07Err)
	}
	ab07ao07OutputBytes, ab07ao07Err := Ab07Ao07Interpret(ab07ao07InputBytes)
	if ab07ao07Err != nil {
		ab07ao07T.Fatal(ab07ao07Err)
	}
	ab07ao07Output := map[string]any{}
	if ab07ao07Err = json.Unmarshal(ab07ao07OutputBytes, &ab07ao07Output); ab07ao07Err != nil {
		ab07ao07T.Fatal(ab07ao07Err)
	}
	ab07ao07Canonical := map[string]any{}
	_ = json.Unmarshal(ab07ao07InputBytes, &ab07ao07Canonical)
	if !reflect.DeepEqual(ab07ao07Output["context"], ab07ao07Canonical["context"]) || !reflect.DeepEqual(ab07ao07Output["raw"], ab07ao07Canonical["fields"]) {
		ab07ao07T.Fatal("context/raw changed", ab07ao07Output)
	}
	if ab07ao07Output["qualified"] != false || ab07ao07Output["physical"] != "Unknown" {
		ab07ao07T.Fatal("qualification inferred")
	}
	if ab07ao07Arguments, ab07ao07Exists := ab07ao07Canonical["arguments"]; ab07ao07Exists {
		if !reflect.DeepEqual(ab07ao07Output["arguments"], ab07ao07Arguments) {
			ab07ao07T.Fatal("arguments changed")
		}
	} else if _, ab07ao07Exists := ab07ao07Output["arguments"]; ab07ao07Exists {
		ab07ao07T.Fatal("absent arguments invented")
	}
	return ab07ao07Output
}

func TestAb07Ao07SyntheticSpans(ab07ao07T *testing.T) {
	ab07ao07Bytes, ab07ao07Err := os.ReadFile("testdata/ab07ao07_b50e_row_spans.json")
	if ab07ao07Err != nil {
		ab07ao07T.Fatal(ab07ao07Err)
	}
	var ab07ao07Rows []struct {
		Ab07Ao07ID     string
		Ab07Ao07Key    int
		Ab07Ao07Value  []int
		Ab07Ao07Fields map[string]any
		Ab07Ao07State  string
	}
	if ab07ao07Err = json.Unmarshal(ab07ao07Bytes, &ab07ao07Rows); ab07ao07Err != nil {
		ab07ao07T.Fatal(ab07ao07Err)
	}
	if len(ab07ao07Rows) != 57 {
		ab07ao07T.Fatal("missing synthetic controls")
	}
	for _, ab07ao07Row := range ab07ao07Rows {
		ab07ao07T.Run(ab07ao07Row.Ab07Ao07ID, func(ab07ao07T *testing.T) {
			ab07ao07Output := ab07ao07Call(ab07ao07T, ab07ao07Input(ab07ao07Row.Ab07Ao07Key, ab07ao07Row.Ab07Ao07Value))
			if ab07ao07Output["state"] != ab07ao07Row.Ab07Ao07State || !reflect.DeepEqual(ab07ao07Output["fields"], ab07ao07Row.Ab07Ao07Fields) {
				ab07ao07T.Fatalf("got %#v; want state %s fields %#v", ab07ao07Output, ab07ao07Row.Ab07Ao07State, ab07ao07Row.Ab07Ao07Fields)
			}
		})
	}
}

func TestAb07Ao07LabelsAndIndependentSpans(ab07ao07T *testing.T) {
	for _, ab07ao07Case := range []struct {
		ab07ao07Key                                int
		ab07ao07Value                              []int
		ab07ao07Name, ab07ao07Label, ab07ao07State string
	}{
		{8, []int{15}, "value", "true", "Hypothesis"}, {8, []int{240}, "value", "false", "Hypothesis"}, {8, []int{0}, "value", "", "Unknown"}, {8, []int{2}, "value", "", "Unknown"},
		{17, []int{16, 0, 0}, "status", "okay", "Hypothesis"}, {17, []int{16, 0, 85}, "status", "short-circuit", "Hypothesis"}, {17, []int{16, 0, 170}, "status", "interrupted", "Hypothesis"}, {17, []int{16, 0, 1}, "status", "", "Unknown"},
	} {
		ab07ao07Output := ab07ao07Call(ab07ao07T, ab07ao07Input(ab07ao07Case.ab07ao07Key, ab07ao07Case.ab07ao07Value))
		ab07ao07Field := ab07ao07Output["fields"].(map[string]any)[ab07ao07Case.ab07ao07Name].(map[string]any)
		if ab07ao07Field["state"] != ab07ao07Case.ab07ao07State {
			ab07ao07T.Fatal(ab07ao07Field)
		}
		if ab07ao07Case.ab07ao07Label == "" {
			if _, ab07ao07Exists := ab07ao07Field["label"]; ab07ao07Exists {
				ab07ao07T.Fatal("borrowed table")
			}
		} else if ab07ao07Field["label"] != ab07ao07Case.ab07ao07Label {
			ab07ao07T.Fatal(ab07ao07Field)
		}
		if _, ab07ao07Exists := ab07ao07Field["number"]; ab07ao07Exists {
			ab07ao07T.Fatal("invented numeric value")
		}
	}
	for _, ab07ao07Value := range [][]int{{16, 0, 18, 52}, {16, 0, 0, 0}, {16, 0, 255, 255}} {
		ab07ao07Output := ab07ao07Call(ab07ao07T, ab07ao07Input(18, ab07ao07Value))
		if ab07ao07Output["state"] != "Hypothesis" {
			ab07ao07T.Fatal("complement validation")
		}
		ab07ao07Inverse := ab07ao07Output["fields"].(map[string]any)["inverse"].(map[string]any)
		if !reflect.DeepEqual(ab07ao07Inverse, map[string]any{"raw": []any{float64(ab07ao07Value[2]), float64(ab07ao07Value[3])}, "state": "Hypothesis"}) {
			ab07ao07T.Fatal("inverse bytes changed", ab07ao07Inverse)
		}
	}
	for _, ab07ao07Case := range []struct {
		ab07ao07Key    int
		ab07ao07Value  []int
		ab07ao07Number float64
	}{{2, []int{0}, 0}, {2, []int{255}, 255}, {3, []int{127}, 127}, {3, []int{128}, -128}, {3, []int{255}, -1}, {9, []int{42}, 42}, {10, []int{255}, -1}} {
		ab07ao07Output := ab07ao07Call(ab07ao07T, ab07ao07Input(ab07ao07Case.ab07ao07Key, ab07ao07Case.ab07ao07Value))
		ab07ao07Field := ab07ao07Output["fields"].(map[string]any)["value"].(map[string]any)
		if ab07ao07Field["number"] != ab07ao07Case.ab07ao07Number || ab07ao07Field["calculation"] != "Proven" || ab07ao07Field["state"] != "Hypothesis" {
			ab07ao07T.Fatal(ab07ao07Field)
		}
	}
}

func TestAb07Ao07ContextAxes(ab07ao07T *testing.T) {
	for _, ab07ao07Case := range []struct {
		ab07ao07Name    string
		ab07ao07Value   any
		ab07ao07Unknown bool
	}{
		{"primary", 180, true}, {"secondary", 15, true}, {"operation", 3, true}, {"operation", 0, true}, {"direction", "negative", true}, {"view", "other", true}, {"family", map[string]any{"opaque": "other"}, false}, {"variant", []any{0, "other"}, false}, {"extra", map[string]any{"nested": []any{false, nil}}, false},
	} {
		ab07ao07T.Run(ab07ao07Case.ab07ao07Name+"-"+ab07ao07CaseName(ab07ao07Case.ab07ao07Value), func(ab07ao07T *testing.T) {
			ab07ao07Object := ab07ao07Input(2, []int{42})
			ab07ao07Object["context"].(map[string]any)[ab07ao07Case.ab07ao07Name] = ab07ao07Case.ab07ao07Value
			ab07ao07Output := ab07ao07Call(ab07ao07T, ab07ao07Object)
			if ab07ao07Case.ab07ao07Unknown {
				if ab07ao07Output["state"] != "Unknown" || len(ab07ao07Output["fields"].(map[string]any)) != 0 {
					ab07ao07T.Fatal("selection guessed")
				}
			} else if ab07ao07Output["state"] != "Hypothesis" {
				ab07ao07T.Fatal("opaque context selected view")
			}
		})
	}
	for _, ab07ao07Name := range []string{"primary", "secondary", "operation", "direction", "view"} {
		ab07ao07Object := ab07ao07Input(2, []int{42})
		delete(ab07ao07Object["context"].(map[string]any), ab07ao07Name)
		if ab07ao07Call(ab07ao07T, ab07ao07Object)["state"] != "Unknown" {
			ab07ao07T.Fatal("missing context selected")
		}
	}
}

func ab07ao07CaseName(ab07ao07Value any) string {
	ab07ao07Bytes, _ := json.Marshal(ab07ao07Value)
	return string(ab07ao07Bytes)
}

func TestAb07Ao07PresenceAndWidths(ab07ao07T *testing.T) {
	for _, ab07ao07Name := range []string{"key", "value", "unmodelled"} {
		for _, ab07ao07Mode := range []string{"missing", "empty", "zero", "extra"} {
			ab07ao07T.Run(ab07ao07Name+ab07ao07Mode, func(ab07ao07T *testing.T) {
				ab07ao07Object := ab07ao07Input(2, []int{42})
				ab07ao07Raw := ab07ao07Object["fields"].(map[string]any)
				switch ab07ao07Mode {
				case "missing":
					delete(ab07ao07Raw, ab07ao07Name)
				case "empty":
					ab07ao07Raw[ab07ao07Name] = []int{}
				case "zero":
					ab07ao07Raw[ab07ao07Name] = []int{0}
				case "extra":
					ab07ao07Raw[ab07ao07Name] = []int{42, 43}
				}
				ab07ao07Output := ab07ao07Call(ab07ao07T, ab07ao07Object)
				if ab07ao07Name == "unmodelled" && (ab07ao07Output["state"] != "Hypothesis" || ab07ao07Output["fields"].(map[string]any)["value"].(map[string]any)["number"] != float64(42)) {
					ab07ao07T.Fatal("independent value lost")
				}
				switch ab07ao07Name {
				case "key":
					if ab07ao07Output["state"] != "Unknown" {
						ab07ao07T.Fatal("key guessed")
					}
				case "value":
					ab07ao07Field := ab07ao07Output["fields"].(map[string]any)["value"].(map[string]any)
					if ab07ao07Mode == "missing" || ab07ao07Mode == "empty" {
						if ab07ao07Output["state"] != "partial" {
							ab07ao07T.Fatal("incomplete lost")
						}
						_, ab07ao07Present := ab07ao07Field["raw"]
						if ab07ao07Present != (ab07ao07Mode == "empty") {
							ab07ao07T.Fatal("missing/empty conflated")
						}
					} else if ab07ao07Field["number"] != float64(ab07ao07Raw["value"].([]int)[0]) {
						ab07ao07T.Fatal("zero/value lost")
					}
				}
			})
		}
	}
	ab07ao07Object := ab07ao07Input(19, []int{42})
	if len(ab07ao07Call(ab07ao07T, ab07ao07Object)["fields"].(map[string]any)) != 0 {
		ab07ao07T.Fatal("unknown key interpreted")
	}
	ab07ao07Object = ab07ao07Input(2, []int{42})
	ab07ao07Object["arguments"] = map[string]any{}
	ab07ao07Call(ab07ao07T, ab07ao07Object)
}

func TestAb07Ao07RejectJSON(ab07ao07T *testing.T) {
	for _, ab07ao07Input := range []string{
		``, `null`, `[]`, `{}`, "{\"context\":{\"family\":\"" + string([]byte{255}) + "\"},\"fields\":{}}", `{"context":{},"fields":{}} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"x":{"a":1,"a":2}},"fields":{}}`, `{"context":{},"fields":{"x":[0],"x":[]}}`, `{"context":{},"fields":{"x":[1.0]}}`, `{"context":{},"fields":{"x":[1e0]}}`, `{"context":{},"fields":{"x":[-1]}}`, `{"context":{},"fields":{"x":[256]}}`, `{"context":{},"fields":{"x":[true]}}`, `{"context":{},"fields":{"x":[null]}}`, `{"context":{},"fields":{"x":["0"]}}`, `{"context":{},"fields":{"x":null}}`, `{"context":{},"fields":{"x":"AA=="}}`, `{"context":{},"fields":[]}`, `{"context":null,"fields":{}}`, `{"context":{"operation":"2"},"fields":{}}`, `{"context":{"primary":null},"fields":{}}`, `{"context":{"secondary":1.5},"fields":{}}`, `{"context":{"direction":2},"fields":{}}`, `{"context":{"view":null},"fields":{}}`, `{"context":{},"fields":{},"arguments":null}`, `{"context":{},"fields":{},"arguments":{"flag":true}}`, `{"context":{},"fields":{},"arguments":{"x":0,"x":0}}`, `{"context":{},"fields":{},"other":0}`, `{"context":{},"fields":{"x":[0,]}}`,
	} {
		ab07ao07Output, ab07ao07Err := Ab07Ao07Interpret(json.RawMessage(ab07ao07Input))
		if ab07ao07Err == nil || ab07ao07Output != nil {
			ab07ao07T.Fatalf("accepted %s", ab07ao07Input)
		}
	}
}

func TestAb07Ao07Ownership(ab07ao07T *testing.T) {
	ab07ao07Object := ab07ao07Input(17, []int{16, 0, 85})
	ab07ao07Bytes, _ := json.Marshal(ab07ao07Object)
	ab07ao07Original := append([]byte(nil), ab07ao07Bytes...)
	ab07ao07First, ab07ao07Err := Ab07Ao07Interpret(ab07ao07Bytes)
	if ab07ao07Err != nil {
		ab07ao07T.Fatal(ab07ao07Err)
	}
	ab07ao07Second, ab07ao07Err := Ab07Ao07Interpret(ab07ao07Bytes)
	if ab07ao07Err != nil {
		ab07ao07T.Fatal(ab07ao07Err)
	}
	for ab07ao07Index := range ab07ao07Bytes {
		ab07ao07Bytes[ab07ao07Index] = ' '
	}
	if !json.Valid(ab07ao07First) || string(ab07ao07First) != string(ab07ao07Second) {
		ab07ao07T.Fatal("input alias")
	}
	ab07ao07Nested := map[string]any{}
	_ = json.Unmarshal(ab07ao07First, &ab07ao07Nested)
	ab07ao07Nested["raw"].(map[string]any)["value"].([]any)[0] = float64(99)
	if ab07ao07Nested["fields"].(map[string]any)["value"].(map[string]any)["raw"].([]any)[0] != float64(16) {
		ab07ao07T.Fatal("nested alias")
	}
	ab07ao07Nested["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "changed"
	for ab07ao07Index := range ab07ao07First {
		ab07ao07First[ab07ao07Index] = ' '
	}
	ab07ao07Third, ab07ao07Err := Ab07Ao07Interpret(ab07ao07Original)
	if ab07ao07Err != nil || string(ab07ao07Third) != string(ab07ao07Second) || strings.Contains(string(ab07ao07Third), "changed") {
		ab07ao07T.Fatal("cross-call alias")
	}
}

func TestAb07Ao07CarrierDomains(ab07ao07T *testing.T) {
	for ab07ao07Byte := 0; ab07ao07Byte < 256; ab07ao07Byte++ {
		for _, ab07ao07Key := range []int{1, 8, 17} {
			ab07ao07Value := []int{ab07ao07Byte}
			ab07ao07Name := "value"
			ab07ao07Label := ""
			if ab07ao07Key == 8 {
				switch ab07ao07Byte {
				case 15:
					ab07ao07Label = "true"
				case 240:
					ab07ao07Label = "false"
				}
			}
			if ab07ao07Key == 17 {
				ab07ao07Value = []int{16, 0, ab07ao07Byte}
				ab07ao07Name = "status"
				switch ab07ao07Byte {
				case 0:
					ab07ao07Label = "okay"
				case 85:
					ab07ao07Label = "short-circuit"
				case 170:
					ab07ao07Label = "interrupted"
				}
			}
			ab07ao07Output := ab07ao07Call(ab07ao07T, ab07ao07Input(ab07ao07Key, ab07ao07Value))
			ab07ao07Field := ab07ao07Output["fields"].(map[string]any)[ab07ao07Name].(map[string]any)
			if _, ab07ao07Exists := ab07ao07Field["number"]; ab07ao07Exists {
				ab07ao07T.Fatal("numeric truth/status invented")
			}
			if ab07ao07Label == "" {
				if ab07ao07Field["state"] != "Unknown" || ab07ao07Field["label"] != nil {
					ab07ao07T.Fatal("unstated mapping", ab07ao07Key, ab07ao07Byte, ab07ao07Field)
				}
			} else if ab07ao07Field["state"] != "Hypothesis" || ab07ao07Field["label"] != ab07ao07Label {
				ab07ao07T.Fatal("mapping lost", ab07ao07Field)
			}
		}
	}
}
