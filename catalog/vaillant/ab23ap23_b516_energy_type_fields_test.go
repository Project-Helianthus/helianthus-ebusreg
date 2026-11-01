package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func ab23ap23Input() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 22, "operation": nil, "direction": "read-positive", "view": "energy-type", "family": "synthetic-family", "variant": "synthetic-variant", "opaque": map[string]any{"nested": []any{false, nil, 0}}}, "fields": map[string]any{"value": []any{193}}, "arguments": map[string]any{"opaque": false}}
}

func ab23ap23Run(ab23ap23T *testing.T, ab23ap23InputValue any) map[string]any {
	ab23ap23T.Helper()
	ab23ap23Encoded, ab23ap23Err := json.Marshal(ab23ap23InputValue)
	if ab23ap23Err != nil {
		ab23ap23T.Fatal(ab23ap23Err)
	}
	ab23ap23Output, ab23ap23Err := Ab23Ap23Evaluate(ab23ap23Encoded)
	if ab23ap23Err != nil {
		ab23ap23T.Fatal(ab23ap23Err)
	}
	ab23ap23Result := map[string]any{}
	if ab23ap23Err := json.Unmarshal(ab23ap23Output, &ab23ap23Result); ab23ap23Err != nil {
		ab23ap23T.Fatal(ab23ap23Err)
	}
	if ab23ap23Result["physical"] != "Unknown" || ab23ap23Result["qualified"] != false {
		ab23ap23T.Fatal(ab23ap23Result)
	}
	return ab23ap23Result
}

func ab23ap23Value(ab23ap23Output map[string]any) map[string]any {
	return ab23ap23Output["fields"].(map[string]any)["value"].(map[string]any)
}

func TestAb23Ap23Directions(ab23ap23T *testing.T) {
	ab23ap23Data, ab23ap23Err := os.ReadFile("testdata/ab23ap23_b516_energy_type_fields.json")
	if ab23ap23Err != nil {
		ab23ap23T.Fatal(ab23ap23Err)
	}
	var ab23ap23Cases []map[string]any
	if ab23ap23Err := json.Unmarshal(ab23ap23Data, &ab23ap23Cases); ab23ap23Err != nil {
		ab23ap23T.Fatal(ab23ap23Err)
	}
	if len(ab23ap23Cases) != 8 {
		ab23ap23T.Fatal("missing vectors")
	}
	for _, ab23ap23Case := range ab23ap23Cases {
		for _, ab23ap23Direction := range []string{"request", "read-positive", "write-positive"} {
			ab23ap23T.Run(ab23ap23Case["id"].(string)+"/"+ab23ap23Direction, func(ab23ap23T *testing.T) {
				ab23ap23In := ab23ap23Input()
				ab23ap23In["context"].(map[string]any)["direction"] = ab23ap23Direction
				ab23ap23In["fields"].(map[string]any)["value"] = []any{ab23ap23Case["value"]}
				ab23ap23Out := ab23ap23Run(ab23ap23T, ab23ap23In)
				ab23ap23Expected := map[string]any{"raw": []any{ab23ap23Case["value"]}, "present": true, "state": "Hypothesis"}
				switch ab23ap23Direction {
				case "request":
					ab23ap23Expected["number"] = ab23ap23Case["value"]
					ab23ap23Expected["unmapped"] = true
				case "read-positive":
					ab23ap23Expected["low_six"] = ab23ap23Case["low"]
					ab23ap23Expected["flags"] = map[string]any{"flag6": ab23ap23Case["flag6"], "flag7": ab23ap23Case["flag7"]}
					ab23ap23Expected["calculation"] = "Proven"
				case "write-positive":
					ab23ap23Expected["low_six"] = ab23ap23Case["low"]
					ab23ap23Expected["unmodeled_high"] = ab23ap23Case["high"]
					ab23ap23Expected["calculation"] = "Proven"
				}
				if ab23ap23Out["state"] != "complete" || !reflect.DeepEqual(ab23ap23Value(ab23ap23Out), ab23ap23Expected) {
					ab23ap23T.Fatalf("got %#v want %#v", ab23ap23Out, ab23ap23Expected)
				}
			})
		}
	}
}

func TestAb23Ap23Context(ab23ap23T *testing.T) {
	for _, ab23ap23Key := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant", "opaque"} {
		for _, ab23ap23Replacement := range []any{nil, 0, false, "", "other", []any{}, map[string]any{}} {
			ab23ap23In := ab23ap23Input()
			ab23ap23Context := ab23ap23In["context"].(map[string]any)
			ab23ap23Context[ab23ap23Key] = ab23ap23Replacement
			ab23ap23Bytes, _ := json.Marshal(ab23ap23In)
			ab23ap23Out, ab23ap23Err := Ab23Ap23Evaluate(ab23ap23Bytes)
			ab23ap23Typed := ab23ap23Key == "primary" || ab23ap23Key == "secondary" || ab23ap23Key == "operation" || ab23ap23Key == "direction" || ab23ap23Key == "view"
			ab23ap23Valid := ab23ap23Replacement == nil
			if ab23ap23Key == "direction" || ab23ap23Key == "view" {
				_, ab23ap23Valid = ab23ap23Replacement.(string)
				ab23ap23Valid = ab23ap23Valid || ab23ap23Replacement == nil
			} else if ab23ap23Typed {
				_, ab23ap23Numeric := ab23ap23Replacement.(int)
				ab23ap23Valid = ab23ap23Valid || ab23ap23Numeric
			}
			if ab23ap23Typed && !ab23ap23Valid {
				if ab23ap23Err == nil || ab23ap23Out != nil {
					ab23ap23T.Fatalf("accepted %s=%v", ab23ap23Key, ab23ap23Replacement)
				}
				continue
			}
			if ab23ap23Err != nil {
				ab23ap23T.Fatal(ab23ap23Err)
			}
			ab23ap23Result := ab23ap23Run(ab23ap23T, ab23ap23In)
			ab23ap23RoundTrip, _ := json.Marshal(ab23ap23Result["context"])
			ab23ap23Original, _ := json.Marshal(ab23ap23Context)
			if string(ab23ap23RoundTrip) != string(ab23ap23Original) {
				ab23ap23T.Fatal("context changed")
			}
			ab23ap23Want := "Unknown"
			if !ab23ap23Typed || (ab23ap23Key == "operation" && ab23ap23Replacement == nil) {
				ab23ap23Want = "complete"
			}
			if ab23ap23Want == "Unknown" {
				ab23ap23Expected := map[string]any{"raw": []any{float64(193)}, "present": true, "state": "Unknown"}
				if !reflect.DeepEqual(ab23ap23Value(ab23ap23Result), ab23ap23Expected) {
					ab23ap23T.Fatal("unsupported context gained interpretation")
				}
			}
			if ab23ap23Result["state"] != ab23ap23Want {
				ab23ap23T.Fatalf("%s=%v: %v", ab23ap23Key, ab23ap23Replacement, ab23ap23Result)
			}
		}
		ab23ap23In := ab23ap23Input()
		delete(ab23ap23In["context"].(map[string]any), ab23ap23Key)
		ab23ap23Out := ab23ap23Run(ab23ap23T, ab23ap23In)
		ab23ap23Want := "Unknown"
		if ab23ap23Key == "family" || ab23ap23Key == "variant" || ab23ap23Key == "opaque" {
			ab23ap23Want = "complete"
		}
		if ab23ap23Out["state"] != ab23ap23Want {
			ab23ap23T.Fatal(ab23ap23Out)
		}
		if _, ab23ap23Present := ab23ap23Out["context"].(map[string]any)[ab23ap23Key]; ab23ap23Present {
			ab23ap23T.Fatal("missing key invented")
		}
	}
}

func TestAb23Ap23IncompleteAndOpaque(ab23ap23T *testing.T) {
	for _, ab23ap23Raw := range []any{nil, []any{}, []any{0, 1}} {
		ab23ap23In := ab23ap23Input()
		ab23ap23Fields := ab23ap23In["fields"].(map[string]any)
		ab23ap23Fields["sibling"] = []any{0, 255}
		if ab23ap23Raw == nil {
			delete(ab23ap23Fields, "value")
		} else {
			ab23ap23Fields["value"] = ab23ap23Raw
		}
		ab23ap23Out := ab23ap23Run(ab23ap23T, ab23ap23In)
		if ab23ap23Out["state"] != "partial" || ab23ap23Value(ab23ap23Out)["state"] != "incomplete" {
			ab23ap23T.Fatal(ab23ap23Out)
		}
		ab23ap23ExpectedValue := map[string]any{"raw": ab23ap23Raw, "present": ab23ap23Raw != nil, "state": "incomplete"}
		ab23ap23ExpectedJSON, _ := json.Marshal(ab23ap23ExpectedValue)
		ab23ap23ActualJSON, _ := json.Marshal(ab23ap23Value(ab23ap23Out))
		if string(ab23ap23ExpectedJSON) != string(ab23ap23ActualJSON) {
			ab23ap23T.Fatal("incomplete evidence changed")
		}
		ab23ap23Expected, _ := json.Marshal(ab23ap23Fields)
		ab23ap23Actual, _ := json.Marshal(ab23ap23Out["raw"])
		if string(ab23ap23Expected) != string(ab23ap23Actual) {
			ab23ap23T.Fatal("raw lost")
		}
		if ab23ap23Out["fields"].(map[string]any)["sibling"].(map[string]any)["state"] != "Unknown" {
			ab23ap23T.Fatal("invented sibling interpretation")
		}
	}
	for _, ab23ap23Argument := range []any{false, nil, 0, []any{}, map[string]any{}, ""} {
		ab23ap23In := ab23ap23Input()
		ab23ap23In["arguments"] = map[string]any{"opaque": ab23ap23Argument}
		ab23ap23Out := ab23ap23Run(ab23ap23T, ab23ap23In)
		ab23ap23Expected, _ := json.Marshal(ab23ap23In["arguments"])
		ab23ap23Actual, _ := json.Marshal(ab23ap23Out["arguments"])
		if string(ab23ap23Expected) != string(ab23ap23Actual) {
			ab23ap23T.Fatal("arguments changed")
		}
	}
	for _, ab23ap23Raw := range []any{[]any{}, []any{0}, []any{255}} {
		ab23ap23In := ab23ap23Input()
		ab23ap23In["fields"].(map[string]any)["sibling"] = ab23ap23Raw
		ab23ap23Out := ab23ap23Run(ab23ap23T, ab23ap23In)
		if ab23ap23Out["state"] != "complete" || ab23ap23Value(ab23ap23Out)["low_six"] != float64(1) {
			ab23ap23T.Fatal("sibling affected value")
		}
	}
	ab23ap23Empty := ab23ap23Input()
	ab23ap23Empty["arguments"] = map[string]any{}
	if len(ab23ap23Run(ab23ap23T, ab23ap23Empty)["arguments"].(map[string]any)) != 0 {
		ab23ap23T.Fatal("empty arguments changed")
	}
	ab23ap23In := ab23ap23Input()
	delete(ab23ap23In, "arguments")
	if _, ab23ap23Present := ab23ap23Run(ab23ap23T, ab23ap23In)["arguments"]; ab23ap23Present {
		ab23ap23T.Fatal("arguments invented")
	}
}

func TestAb23Ap23Reject(ab23ap23T *testing.T) {
	for _, ab23ap23Bad := range []string{``, `{`, `null`, `[]`, `{}`, `{"context":{},"fields":{},"fields":{}}`, `{"context":{"a":{"b":0,"b":1}},"fields":{}}`, `{"context":{},"fields":{},"arguments":{"a":[{"b":0,"b":1}]}}`, `{"context":{},"fields":{}} false`} {
		if ab23ap23Out, ab23ap23Err := Ab23Ap23Evaluate([]byte(ab23ap23Bad)); ab23ap23Err == nil || ab23ap23Out != nil {
			ab23ap23T.Fatalf("accepted %s", ab23ap23Bad)
		}
	}
	for _, ab23ap23Bad := range []any{nil, false, "0", 0, []any{true}, []any{nil}, []any{"1"}, []any{-1}, []any{256}, []any{1.5}, []any{[]any{}}} {
		ab23ap23In := ab23ap23Input()
		ab23ap23In["fields"].(map[string]any)["value"] = ab23ap23Bad
		ab23ap23Bytes, _ := json.Marshal(ab23ap23In)
		if ab23ap23Out, ab23ap23Err := Ab23Ap23Evaluate(ab23ap23Bytes); ab23ap23Err == nil || ab23ap23Out != nil {
			ab23ap23T.Fatalf("accepted %s", ab23ap23Bytes)
		}
	}
	for _, ab23ap23Key := range []string{"context", "fields", "arguments"} {
		for _, ab23ap23Bad := range []any{nil, false, 0, "", []any{}} {
			ab23ap23In := ab23ap23Input()
			ab23ap23In[ab23ap23Key] = ab23ap23Bad
			ab23ap23Bytes, _ := json.Marshal(ab23ap23In)
			if ab23ap23Out, ab23ap23Err := Ab23Ap23Evaluate(ab23ap23Bytes); ab23ap23Err == nil || ab23ap23Out != nil {
				ab23ap23T.Fatal("accepted wrong object kind")
			}
		}
	}
	for _, ab23ap23Key := range []string{"primary", "secondary", "operation"} {
		for _, ab23ap23Bad := range []any{-1, 256, 1.5, "181", false} {
			ab23ap23In := ab23ap23Input()
			ab23ap23In["context"].(map[string]any)[ab23ap23Key] = ab23ap23Bad
			ab23ap23Bytes, _ := json.Marshal(ab23ap23In)
			if ab23ap23Out, ab23ap23Err := Ab23Ap23Evaluate(ab23ap23Bytes); ab23ap23Err == nil || ab23ap23Out != nil {
				ab23ap23T.Fatal("accepted invalid numeric context")
			}
		}
	}
}

func TestAb23Ap23Isolation(ab23ap23T *testing.T) {
	ab23ap23Bytes, _ := json.Marshal(ab23ap23Input())
	ab23ap23Out, ab23ap23Err := Ab23Ap23Evaluate(ab23ap23Bytes)
	if ab23ap23Err != nil {
		ab23ap23T.Fatal(ab23ap23Err)
	}
	ab23ap23Saved := string(ab23ap23Out)
	for ab23ap23Index := range ab23ap23Bytes {
		ab23ap23Bytes[ab23ap23Index] = ' '
	}
	if string(ab23ap23Out) != ab23ap23Saved {
		ab23ap23T.Fatal("input aliases output")
	}
	ab23ap23Decoded := map[string]any{}
	if ab23ap23Err := json.Unmarshal(ab23ap23Out, &ab23ap23Decoded); ab23ap23Err != nil {
		ab23ap23T.Fatal(ab23ap23Err)
	}
	ab23ap23Decoded["raw"].(map[string]any)["value"].([]any)[0] = float64(0)
	ab23ap23Decoded["context"].(map[string]any)["opaque"].(map[string]any)["nested"].([]any)[0] = true
	ab23ap23Decoded["arguments"].(map[string]any)["opaque"] = true
	if ab23ap23Value(ab23ap23Decoded)["raw"].([]any)[0] != float64(193) {
		ab23ap23T.Fatal("raw aliases semantic raw")
	}
	for ab23ap23Index := range ab23ap23Out {
		ab23ap23Out[ab23ap23Index] = ' '
	}
	ab23ap23Fresh, _ := json.Marshal(ab23ap23Input())
	ab23ap23Next, ab23ap23Err := Ab23Ap23Evaluate(ab23ap23Fresh)
	if ab23ap23Err != nil || string(ab23ap23Next) != ab23ap23Saved {
		ab23ap23T.Fatal("result aliases later call")
	}
	if strings.Contains(string(ab23ap23Next), "meter") || strings.Contains(string(ab23ap23Next), "label") {
		ab23ap23T.Fatal("invented semantic association")
	}
}
