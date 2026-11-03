package ebus_standard_catalog

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func ab17aq17Run(ab17aq17T *testing.T, ab17aq17Input any) map[string]any {
	ab17aq17T.Helper()
	ab17aq17Raw, ab17aq17Err := json.Marshal(ab17aq17Input)
	if ab17aq17Err != nil {
		ab17aq17T.Fatal(ab17aq17Err)
	}
	ab17aq17Output, ab17aq17Err := Ab17Aq17Evaluate(ab17aq17Raw)
	if ab17aq17Err != nil {
		ab17aq17T.Fatal(ab17aq17Err)
	}
	var ab17aq17Result map[string]any
	if ab17aq17Err := json.Unmarshal(ab17aq17Output, &ab17aq17Result); ab17aq17Err != nil {
		ab17aq17T.Fatal(ab17aq17Err)
	}
	if ab17aq17Result["physical"] != "Unknown" || ab17aq17Result["qualified"] != false {
		ab17aq17T.Fatal("arithmetic qualified physical support", ab17aq17Result)
	}
	return ab17aq17Result
}

func ab17aq17Base() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 255, "secondary": 4, "direction": "positive", "view": "body-extent", "family": "synthetic-family", "variant": ""}, "fields": map[string]any{"body": []int{0, 1}, "opaque": []int{}}, "arguments": map[string]any{"declared_length": 2}}
}

func TestAb17Aq17Fixture(ab17aq17T *testing.T) {
	ab17aq17Raw, ab17aq17Err := os.ReadFile("testdata/ab17aq17_standardff_extent_view.json")
	if ab17aq17Err != nil {
		ab17aq17T.Fatal(ab17aq17Err)
	}
	var ab17aq17Cases []struct {
		Ab17Aq17ID       string         `json:"id"`
		Ab17Aq17Input    map[string]any `json:"input"`
		Ab17Aq17Expected map[string]any `json:"expected"`
	}
	if ab17aq17Err := json.Unmarshal(ab17aq17Raw, &ab17aq17Cases); ab17aq17Err != nil {
		ab17aq17T.Fatal(ab17aq17Err)
	}
	if len(ab17aq17Cases) != 11 {
		ab17aq17T.Fatal("fixture count")
	}
	for _, ab17aq17Case := range ab17aq17Cases {
		ab17aq17T.Run(ab17aq17Case.Ab17Aq17ID, func(ab17aq17T *testing.T) {
			ab17aq17Out := ab17aq17Run(ab17aq17T, ab17aq17Case.Ab17Aq17Input)
			for _, ab17aq17Key := range []string{"context", "fields", "arguments", "extra"} {
				if !reflect.DeepEqual(ab17aq17Out[ab17aq17Key], ab17aq17Case.Ab17Aq17Input[ab17aq17Key]) {
					ab17aq17T.Fatal("raw input changed", ab17aq17Key)
				}
			}
			ab17aq17Results := ab17aq17Out["results"].(map[string]any)
			if len(ab17aq17Results) != 2 {
				ab17aq17T.Fatal("decoded unprovided fields")
			}
			ab17aq17Body := ab17aq17Results["body"].(map[string]any)
			if ab17aq17Body["state"] != "Unknown" || ab17aq17Body["calculation"] != "Proven" {
				ab17aq17T.Fatal("semantic grade", ab17aq17Body)
			}
			if !reflect.DeepEqual(ab17aq17Body["raw"], ab17aq17Case.Ab17Aq17Input["fields"].(map[string]any)["body"]) {
				ab17aq17T.Fatal("body raw")
			}
			for _, ab17aq17Key := range []string{"flags", "declared_pair_count", "available_pair_count"} {
				if !reflect.DeepEqual(ab17aq17Body[ab17aq17Key], ab17aq17Case.Ab17Aq17Expected[ab17aq17Key]) {
					ab17aq17T.Fatalf("%s: got %v want %v", ab17aq17Key, ab17aq17Body[ab17aq17Key], ab17aq17Case.Ab17Aq17Expected[ab17aq17Key])
				}
			}
			ab17aq17Flags := ab17aq17Body["flags"].(map[string]any)
			ab17aq17Complete := ab17aq17Flags["bounded_minimum"] == true || ab17aq17Flags["whole_declared_span"] == true
			if (ab17aq17Out["state"] == "complete") != ab17aq17Complete {
				ab17aq17T.Fatal("extent completeness", ab17aq17Out)
			}
		})
	}
}

func TestAb17Aq17Identity(ab17aq17T *testing.T) {
	for ab17aq17Key, ab17aq17Values := range map[string][]any{"primary": {254, nil}, "secondary": {5, nil}, "operation": {0, 1}, "direction": {"negative", ""}, "view": {"other", ""}} {
		for _, ab17aq17Value := range ab17aq17Values {
			ab17aq17Input := ab17aq17Base()
			ab17aq17Context := ab17aq17Input["context"].(map[string]any)
			if ab17aq17Value == nil {
				delete(ab17aq17Context, ab17aq17Key)
			} else {
				ab17aq17Context[ab17aq17Key] = ab17aq17Value
			}
			ab17aq17Out := ab17aq17Run(ab17aq17T, ab17aq17Input)
			ab17aq17Body := ab17aq17Out["results"].(map[string]any)["body"].(map[string]any)
			if ab17aq17Out["state"] != "Unknown" || len(ab17aq17Body) != 2 || ab17aq17Body["state"] != "Unknown" {
				ab17aq17T.Fatal("unsupported identity calculated", ab17aq17Key, ab17aq17Out)
			}
		}
	}
	for _, ab17aq17Key := range []string{"family", "variant"} {
		for _, ab17aq17Value := range []any{nil, false, 0, "another-family", "", []any{}} {
			ab17aq17Input := ab17aq17Base()
			ab17aq17Input["context"].(map[string]any)[ab17aq17Key] = ab17aq17Value
			ab17aq17Out := ab17aq17Run(ab17aq17T, ab17aq17Input)
			if ab17aq17Out["state"] != "complete" || !reflect.DeepEqual(ab17aq17Out["context"].(map[string]any)[ab17aq17Key], func() any {
				ab17aq17Raw, _ := json.Marshal(ab17aq17Value)
				var ab17aq17Normalized any
				_ = json.Unmarshal(ab17aq17Raw, &ab17aq17Normalized)
				return ab17aq17Normalized
			}()) {
				ab17aq17T.Fatal("opaque identity", ab17aq17Key)
			}
		}
	}
}

func TestAb17Aq17Invalid(ab17aq17T *testing.T) {
	for _, ab17aq17Raw := range []string{``, `{`, `[]`, `null`, `{} {}`, `{"x":{"a":1,"a":2}}`, `{"x":[{"a":1,"\u0061":2}]}`, `{"x":NaN}`} {
		if ab17aq17Out, ab17aq17Err := Ab17Aq17Evaluate(json.RawMessage(ab17aq17Raw)); ab17aq17Err == nil || ab17aq17Out != nil {
			ab17aq17T.Fatal("accepted malformed", ab17aq17Raw)
		}
	}
	for _, ab17aq17Key := range []string{"primary", "secondary", "operation", "declared_length", "body", "opaque", "direction", "view", "arguments", "fields", "context"} {
		for _, ab17aq17Bad := range []any{nil, true, "2", -1, 256, 1.5} {
			ab17aq17Input := ab17aq17Base()
			switch ab17aq17Key {
			case "primary", "secondary", "operation":
				ab17aq17Input["context"].(map[string]any)[ab17aq17Key] = ab17aq17Bad
			case "direction", "view":
				if _, ab17aq17OK := ab17aq17Bad.(string); ab17aq17OK {
					continue
				}
				ab17aq17Input["context"].(map[string]any)[ab17aq17Key] = ab17aq17Bad
			case "declared_length":
				ab17aq17Input["arguments"].(map[string]any)[ab17aq17Key] = ab17aq17Bad
			case "body", "opaque":
				ab17aq17Input["fields"].(map[string]any)[ab17aq17Key] = []any{ab17aq17Bad}
			default:
				ab17aq17Input[ab17aq17Key] = ab17aq17Bad
			}
			ab17aq17Raw, _ := json.Marshal(ab17aq17Input)
			if ab17aq17Out, ab17aq17Err := Ab17Aq17Evaluate(ab17aq17Raw); ab17aq17Err == nil || ab17aq17Out != nil {
				ab17aq17T.Fatal("accepted wrong kind", ab17aq17Key, ab17aq17Bad)
			}
		}
	}
	for _, ab17aq17Raw := range []string{`null`, `true`, `{}`, `""`} {
		ab17aq17Raw = []string{`{"context":{"primary":255,"secondary":4,"direction":"positive","view":"body-extent"},"fields":{"body":` + ab17aq17Raw + `},"arguments":{"declared_length":2}}`}[0]
		if ab17aq17Out, ab17aq17Err := Ab17Aq17Evaluate([]byte(ab17aq17Raw)); ab17aq17Err == nil || ab17aq17Out != nil {
			ab17aq17T.Fatal("nonarray body")
		}
	}
	for _, ab17aq17Key := range []string{"arguments", "declared_length"} {
		ab17aq17Input := ab17aq17Base()
		if ab17aq17Key == "arguments" {
			delete(ab17aq17Input, ab17aq17Key)
		} else {
			delete(ab17aq17Input["arguments"].(map[string]any), ab17aq17Key)
		}
		ab17aq17Raw, _ := json.Marshal(ab17aq17Input)
		if _, ab17aq17Err := Ab17Aq17Evaluate(ab17aq17Raw); ab17aq17Err == nil {
			ab17aq17T.Fatal("missing argument")
		}
	}
}

func TestAb17Aq17PresenceAndMemory(ab17aq17T *testing.T) {
	for _, ab17aq17Body := range [][]int{nil, {}, {0}, {0, 1}, {0, 1, 2}} {
		ab17aq17Input := ab17aq17Base()
		ab17aq17Fields := ab17aq17Input["fields"].(map[string]any)
		if ab17aq17Body == nil {
			delete(ab17aq17Fields, "body")
		} else {
			ab17aq17Fields["body"] = ab17aq17Body
		}
		ab17aq17Out := ab17aq17Run(ab17aq17T, ab17aq17Input)
		ab17aq17Results := ab17aq17Out["results"].(map[string]any)
		if !reflect.DeepEqual(ab17aq17Results["opaque"], map[string]any{"raw": []any{}, "state": "Unknown"}) {
			ab17aq17T.Fatal("independent sibling lost")
		}
		ab17aq17Result := ab17aq17Results["body"].(map[string]any)
		if ab17aq17Body == nil {
			if _, ab17aq17OK := ab17aq17Result["raw"]; ab17aq17OK {
				ab17aq17T.Fatal("invented raw")
			}
			if _, ab17aq17OK := ab17aq17Result["available_length"]; ab17aq17OK {
				ab17aq17T.Fatal("invented length")
			}
		} else if ab17aq17Result["available_length"] != float64(len(ab17aq17Body)) {
			ab17aq17T.Fatal("available length")
		}
	}
	ab17aq17Raw, _ := json.Marshal(ab17aq17Base())
	ab17aq17Original := append([]byte(nil), ab17aq17Raw...)
	ab17aq17First, ab17aq17Err := Ab17Aq17Evaluate(ab17aq17Raw)
	if ab17aq17Err != nil {
		ab17aq17T.Fatal(ab17aq17Err)
	}
	ab17aq17Saved := string(ab17aq17First)
	for ab17aq17Index := range ab17aq17Raw {
		ab17aq17Raw[ab17aq17Index] = ' '
	}
	if string(ab17aq17First) != ab17aq17Saved {
		ab17aq17T.Fatal("input aliases output")
	}
	for ab17aq17Index := range ab17aq17First {
		ab17aq17First[ab17aq17Index] = ' '
	}
	ab17aq17Next, ab17aq17Err := Ab17Aq17Evaluate(ab17aq17Original)
	if ab17aq17Err != nil || string(ab17aq17Next) != ab17aq17Saved {
		ab17aq17T.Fatal("cross call mutation")
	}
	if strings.Contains(ab17aq17Saved, "address") || strings.Contains(ab17aq17Saved, "status") {
		ab17aq17T.Fatal("invented decoder")
	}
}

func TestAb17Aq17AdditionalBoundaries(ab17aq17T *testing.T) {
	for _, ab17aq17Key := range []string{"family", "variant"} {
		ab17aq17Input := ab17aq17Base()
		delete(ab17aq17Input["context"].(map[string]any), ab17aq17Key)
		if ab17aq17Run(ab17aq17T, ab17aq17Input)["state"] != "complete" {
			ab17aq17T.Fatal("opaque absence changes arithmetic")
		}
	}
	for _, ab17aq17Length := range []int{0, 1, 11, 255} {
		ab17aq17Input := ab17aq17Base()
		ab17aq17Input["arguments"].(map[string]any)["declared_length"] = ab17aq17Length
		if ab17aq17Run(ab17aq17T, ab17aq17Input)["state"] != "partial" {
			ab17aq17T.Fatal("invalid minimum admitted", ab17aq17Length)
		}
	}
	ab17aq17Input := ab17aq17Base()
	ab17aq17Input["fields"].(map[string]any)["body"] = []int{0, 255, 254}
	ab17aq17Out := ab17aq17Run(ab17aq17T, ab17aq17Input)
	if ab17aq17Out["state"] != "complete" {
		ab17aq17T.Fatal("extra byte repairs or breaks minimum")
	}
	delete(ab17aq17Input["fields"].(map[string]any), "opaque")
	ab17aq17Out = ab17aq17Run(ab17aq17T, ab17aq17Input)
	if len(ab17aq17Out["results"].(map[string]any)) != 1 {
		ab17aq17T.Fatal("invented absent sibling")
	}
	ab17aq17Input["arguments"].(map[string]any)["extra"] = 0
	ab17aq17Raw, _ := json.Marshal(ab17aq17Input)
	if _, ab17aq17Err := Ab17Aq17Evaluate(ab17aq17Raw); ab17aq17Err == nil {
		ab17aq17T.Fatal("extra typed argument admitted")
	}
	for _, ab17aq17Body := range [][]int{nil, {}} {
		ab17aq17Input = ab17aq17Base()
		ab17aq17Input["context"].(map[string]any)["secondary"] = 6
		if ab17aq17Body == nil {
			delete(ab17aq17Input["fields"].(map[string]any), "body")
		} else {
			ab17aq17Input["fields"].(map[string]any)["body"] = ab17aq17Body
		}
		ab17aq17Result := ab17aq17Run(ab17aq17T, ab17aq17Input)["results"].(map[string]any)["body"].(map[string]any)
		if _, ab17aq17OK := ab17aq17Result["available_pair_count"]; ab17aq17OK {
			ab17aq17T.Fatal("empty or absent data invented pair count")
		}
	}
}

// TestAb17Aq17EvaluateExportedContract drives the exported entry point with
// literal JSON and compares the whole decoded result (or the error) with
// independently written expectations, so a change to the extent arithmetic,
// the identity gate, raw preservation or input rejection fails here.
func TestAb17Aq17EvaluateExportedContract(ab17aq17T *testing.T) {
	ab17aq17Cases := []struct {
		ab17aq17Name  string
		ab17aq17Input string
		ab17aq17Want  string
		ab17aq17Err   string
	}{
		{
			ab17aq17Name:  "ff06-declared5-available3",
			ab17aq17Input: `{"context":{"primary":255,"secondary":6,"direction":"positive","view":"body-extent"},"fields":{"body":[1,2,3]},"arguments":{"declared_length":5}}`,
			ab17aq17Want:  `{"context":{"primary":255,"secondary":6,"direction":"positive","view":"body-extent"},"fields":{"body":[1,2,3]},"arguments":{"declared_length":5},"state":"partial","physical":"Unknown","qualified":false,"results":{"body":{"state":"Unknown","calculation":"Proven","declared_length":5,"raw":[1,2,3],"available_length":3,"declared_pair_count":2,"available_pair_count":1,"flags":{"declared_incomplete_tail":false,"whole_declared_span":false}}}}`,
		},
		{
			ab17aq17Name:  "ff06-declared4-available6",
			ab17aq17Input: `{"context":{"primary":255,"secondary":6,"direction":"positive","view":"body-extent"},"fields":{"body":[1,2,3,4,5,6]},"arguments":{"declared_length":4}}`,
			ab17aq17Want:  `{"context":{"primary":255,"secondary":6,"direction":"positive","view":"body-extent"},"fields":{"body":[1,2,3,4,5,6]},"arguments":{"declared_length":4},"state":"partial","physical":"Unknown","qualified":false,"results":{"body":{"state":"Unknown","calculation":"Proven","declared_length":4,"raw":[1,2,3,4,5,6],"available_length":6,"declared_pair_count":1,"available_pair_count":1,"flags":{"declared_incomplete_tail":true,"whole_declared_span":false}}}}`,
		},
		{
			ab17aq17Name:  "ff04-declared10-available10",
			ab17aq17Input: `{"context":{"primary":255,"secondary":4,"direction":"positive","view":"body-extent"},"fields":{"body":[0,1,2,3,4,5,6,7,8,9]},"arguments":{"declared_length":10}}`,
			ab17aq17Want:  `{"context":{"primary":255,"secondary":4,"direction":"positive","view":"body-extent"},"fields":{"body":[0,1,2,3,4,5,6,7,8,9]},"arguments":{"declared_length":10},"state":"complete","physical":"Unknown","qualified":false,"results":{"body":{"state":"Unknown","calculation":"Proven","declared_length":10,"raw":[0,1,2,3,4,5,6,7,8,9],"available_length":10,"flags":{"declared_minimum_satisfied":true,"initialized_minimum_satisfied":true,"bounded_minimum":true}}}}`,
		},
		{
			ab17aq17Name:  "ff04-declared1-available1",
			ab17aq17Input: `{"context":{"primary":255,"secondary":4,"direction":"positive","view":"body-extent"},"fields":{"body":[7]},"arguments":{"declared_length":1}}`,
			ab17aq17Want:  `{"context":{"primary":255,"secondary":4,"direction":"positive","view":"body-extent"},"fields":{"body":[7]},"arguments":{"declared_length":1},"state":"partial","physical":"Unknown","qualified":false,"results":{"body":{"state":"Unknown","calculation":"Proven","declared_length":1,"raw":[7],"available_length":1,"flags":{"declared_minimum_satisfied":false,"initialized_minimum_satisfied":false,"bounded_minimum":false}}}}`,
		},
		{
			ab17aq17Name:  "zero-operation-stays-raw-only",
			ab17aq17Input: `{"context":{"primary":255,"secondary":4,"operation":0,"direction":"positive","view":"body-extent"},"fields":{"body":[0,1]},"arguments":{"declared_length":2}}`,
			ab17aq17Want:  `{"context":{"primary":255,"secondary":4,"operation":0,"direction":"positive","view":"body-extent"},"fields":{"body":[0,1]},"arguments":{"declared_length":2},"state":"Unknown","physical":"Unknown","qualified":false,"results":{"body":{"raw":[0,1],"state":"Unknown"}}}`,
		},
		{
			ab17aq17Name:  "trailing-json",
			ab17aq17Input: `{"context":{"primary":255,"secondary":6,"direction":"positive","view":"body-extent"},"fields":{"body":[1,2,3]},"arguments":{"declared_length":5}} {}`,
			ab17aq17Err:   "trailing JSON",
		},
		{
			ab17aq17Name:  "declared-length-outside-byte",
			ab17aq17Input: `{"context":{"primary":255,"secondary":6,"direction":"positive","view":"body-extent"},"fields":{"body":[1,2,3]},"arguments":{"declared_length":256}}`,
			ab17aq17Err:   "integer outside byte range",
		},
	}
	for _, ab17aq17Case := range ab17aq17Cases {
		ab17aq17T.Run(ab17aq17Case.ab17aq17Name, func(ab17aq17T *testing.T) {
			ab17aq17Out, ab17aq17Err := Ab17Aq17Evaluate(json.RawMessage(ab17aq17Case.ab17aq17Input))
			if ab17aq17Case.ab17aq17Err != "" {
				if ab17aq17Err == nil || !strings.Contains(ab17aq17Err.Error(), ab17aq17Case.ab17aq17Err) || ab17aq17Out != nil {
					ab17aq17T.Fatalf("got output %s error %v, want nil output and error containing %q", ab17aq17Out, ab17aq17Err, ab17aq17Case.ab17aq17Err)
				}
				return
			}
			if ab17aq17Err != nil {
				ab17aq17T.Fatalf("unexpected error: %v", ab17aq17Err)
			}
			var ab17aq17Got, ab17aq17Want any
			if ab17aq17Err := json.Unmarshal(ab17aq17Out, &ab17aq17Got); ab17aq17Err != nil {
				ab17aq17T.Fatal(ab17aq17Err)
			}
			if ab17aq17Err := json.Unmarshal([]byte(ab17aq17Case.ab17aq17Want), &ab17aq17Want); ab17aq17Err != nil {
				ab17aq17T.Fatal(ab17aq17Err)
			}
			if !reflect.DeepEqual(ab17aq17Got, ab17aq17Want) {
				ab17aq17T.Fatalf("got %s\nwant %s", ab17aq17Out, ab17aq17Case.ab17aq17Want)
			}
		})
	}
}

func TestAb17Aq17OutputKeyCollision(ab17aq17T *testing.T) {
	for _, ab17aq17Key := range []string{"state", "results", "physical", "qualified"} {
		ab17aq17Input := ab17aq17Base()
		ab17aq17Input[ab17aq17Key] = 7
		ab17aq17Raw, ab17aq17Err := json.Marshal(ab17aq17Input)
		if ab17aq17Err != nil {
			ab17aq17T.Fatal(ab17aq17Err)
		}
		if ab17aq17Out, ab17aq17Err := Ab17Aq17Evaluate(ab17aq17Raw); ab17aq17Err == nil || ab17aq17Out != nil {
			ab17aq17T.Fatal("accepted colliding key", ab17aq17Key)
		}
	}
}
