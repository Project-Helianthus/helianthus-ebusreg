package vaillant

import (
	ab17ap17Bytes "bytes"
	ab17ap17JSON "encoding/json"
	ab17ap17OS "os"
	ab17ap17Reflect "reflect"
	ab17ap17Strings "strings"
	ab17ap17Testing "testing"
)

func ab17ap17Decode(ab17ap17T *ab17ap17Testing.T, ab17ap17Raw []byte) map[string]any {
	ab17ap17T.Helper()
	var ab17ap17Object map[string]any
	ab17ap17Decoder := ab17ap17JSON.NewDecoder(ab17ap17Bytes.NewReader(ab17ap17Raw))
	ab17ap17Decoder.UseNumber()
	if ab17ap17Err := ab17ap17Decoder.Decode(&ab17ap17Object); ab17ap17Err != nil {
		ab17ap17T.Fatal(ab17ap17Err)
	}
	return ab17ap17Object
}

func ab17ap17Run(ab17ap17T *ab17ap17Testing.T, ab17ap17Input any) map[string]any {
	ab17ap17T.Helper()
	ab17ap17Raw, ab17ap17Err := ab17ap17JSON.Marshal(ab17ap17Input)
	if ab17ap17Err != nil {
		ab17ap17T.Fatal(ab17ap17Err)
	}
	ab17ap17Output, ab17ap17Err := Ab17Ap17Evaluate(ab17ap17Raw)
	if ab17ap17Err != nil {
		ab17ap17T.Fatal(ab17ap17Err)
	}
	return ab17ap17Decode(ab17ap17T, ab17ap17Output)
}

func ab17ap17Base(ab17ap17T *ab17ap17Testing.T) map[string]any {
	return ab17ap17Decode(ab17ap17T, []byte(`{"context":{"primary":181,"secondary":9,"operation":52,"direction":"request","view":"diagnostic-get","family":"synthetic-family","variant":"synthetic-variant","opaque":{"enabled":false,"items":[null,0] }},"fields":{"address":[120,86,52,18],"count":[255]},"arguments":{"requested_address":0,"extra":false}}`))
}

func TestAb17Ap17SyntheticCases(ab17ap17T *ab17ap17Testing.T) {
	ab17ap17Raw, ab17ap17Err := ab17ap17OS.ReadFile("testdata/ab17ap17_b509_diagnostic_fields.json")
	if ab17ap17Err != nil {
		ab17ap17T.Fatal(ab17ap17Err)
	}
	var ab17ap17Cases []struct {
		Ab17Ap17ID       string                  `json:"id"`
		Ab17Ap17Input    ab17ap17JSON.RawMessage `json:"input"`
		Ab17Ap17Expected ab17ap17JSON.RawMessage `json:"expected"`
	}
	if ab17ap17Err := ab17ap17JSON.Unmarshal(ab17ap17Raw, &ab17ap17Cases); ab17ap17Err != nil {
		ab17ap17T.Fatal(ab17ap17Err)
	}
	if len(ab17ap17Cases) != 21 {
		ab17ap17T.Fatalf("fixture count %d", len(ab17ap17Cases))
	}
	for _, ab17ap17Case := range ab17ap17Cases {
		ab17ap17T.Run(ab17ap17Case.Ab17Ap17ID, func(ab17ap17T *ab17ap17Testing.T) {
			ab17ap17Output, ab17ap17Err := Ab17Ap17Evaluate(ab17ap17Case.Ab17Ap17Input)
			if ab17ap17Err != nil {
				ab17ap17T.Fatal(ab17ap17Err)
			}
			ab17ap17Actual := ab17ap17Decode(ab17ap17T, ab17ap17Output)
			ab17ap17Expected := ab17ap17Decode(ab17ap17T, ab17ap17Case.Ab17Ap17Expected)
			if !ab17ap17Reflect.DeepEqual(ab17ap17Actual, ab17ap17Expected) {
				ab17ap17T.Fatalf("got %s; want %s", ab17ap17Output, ab17ap17Case.Ab17Ap17Expected)
			}
		})
	}
}

func TestAb17Ap17ContextAxes(ab17ap17T *ab17ap17Testing.T) {
	for _, ab17ap17Key := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant", "opaque"} {
		for _, ab17ap17Mode := range []string{"absent", "null", "changed", "zero"} {
			ab17ap17T.Run("ab17ap17_"+ab17ap17Key+"_"+ab17ap17Mode, func(ab17ap17T *ab17ap17Testing.T) {
				ab17ap17Input := ab17ap17Base(ab17ap17T)
				ab17ap17Context := ab17ap17Input["context"].(map[string]any)
				switch ab17ap17Mode {
				case "absent":
					delete(ab17ap17Context, ab17ap17Key)
				case "null":
					ab17ap17Context[ab17ap17Key] = nil
				default:
					switch ab17ap17Key {
					case "primary", "secondary", "operation":
						if ab17ap17Mode == "zero" {
							ab17ap17Context[ab17ap17Key] = 0
						} else {
							ab17ap17Context[ab17ap17Key] = 53
						}
					case "direction", "view":
						ab17ap17Context[ab17ap17Key] = "unsupported"
					default:
						ab17ap17Context[ab17ap17Key] = false
					}
				}
				ab17ap17Output := ab17ap17Run(ab17ap17T, ab17ap17Input)
				ab17ap17Want := "Unknown"
				if ab17ap17Key == "family" || ab17ap17Key == "variant" || ab17ap17Key == "opaque" {
					ab17ap17Want = "complete"
				}
				if ab17ap17Output["state"] != ab17ap17Want || ab17ap17Output["qualified"] != false || ab17ap17Output["physical"] != "Unknown" {
					ab17ap17T.Fatalf("selector admitted or qualified: %v", ab17ap17Output)
				}
				ab17ap17Normalized := ab17ap17Decode(ab17ap17T, ab17ap17MustMarshal(ab17ap17T, ab17ap17Input))
				if !ab17ap17Reflect.DeepEqual(ab17ap17Output["context"], ab17ap17Normalized["context"]) || !ab17ap17Reflect.DeepEqual(ab17ap17Output["raw"], ab17ap17Normalized["fields"]) {
					ab17ap17T.Fatal("context/raw changed")
				}
				if ab17ap17Want == "Unknown" && len(ab17ap17Output["fields"].(map[string]any)) != 0 {
					ab17ap17T.Fatal("unsupported formula applied")
				}
			})
		}
	}
}

func ab17ap17MustMarshal(ab17ap17T *ab17ap17Testing.T, ab17ap17Value any) []byte {
	ab17ap17T.Helper()
	ab17ap17Raw, ab17ap17Err := ab17ap17JSON.Marshal(ab17ap17Value)
	if ab17ap17Err != nil {
		ab17ap17T.Fatal(ab17ap17Err)
	}
	return ab17ap17Raw
}

func TestAb17Ap17PartialFields(ab17ap17T *ab17ap17Testing.T) {
	for _, ab17ap17Key := range []string{"address", "count", "status", "data"} {
		for _, ab17ap17Mode := range []string{"absent", "empty", "short", "overlong"} {
			ab17ap17T.Run("ab17ap17_"+ab17ap17Key+"_"+ab17ap17Mode, func(ab17ap17T *ab17ap17Testing.T) {
				ab17ap17Input := ab17ap17Base(ab17ap17T)
				ab17ap17Context := ab17ap17Input["context"].(map[string]any)
				ab17ap17Raw := ab17ap17Input["fields"].(map[string]any)
				ab17ap17Sibling := "address"
				if ab17ap17Key == "address" {
					ab17ap17Sibling = "count"
				}
				if ab17ap17Key == "data" {
					ab17ap17Context["direction"] = "response"
					delete(ab17ap17Raw, "count")
				}
				if ab17ap17Key == "status" {
					ab17ap17Context["operation"] = 53
					ab17ap17Context["view"] = "diagnostic-put"
					ab17ap17Context["direction"] = "response"
					ab17ap17Raw["status"] = []int{6}
				}
				switch ab17ap17Mode {
				case "absent":
					delete(ab17ap17Raw, ab17ap17Key)
				case "empty":
					ab17ap17Raw[ab17ap17Key] = []int{}
				case "short":
					ab17ap17Raw[ab17ap17Key] = []int{1, 2, 3}
				case "overlong":
					ab17ap17Raw[ab17ap17Key] = []int{1, 2, 3, 4, 5}
				}
				ab17ap17Output := ab17ap17Run(ab17ap17T, ab17ap17Input)
				ab17ap17Fields := ab17ap17Output["fields"].(map[string]any)
				ab17ap17Result := ab17ap17Fields[ab17ap17Key].(map[string]any)
				ab17ap17Want := "incomplete"
				if ab17ap17Key == "data" && ab17ap17Mode != "absent" {
					ab17ap17Want = "complete"
				}
				if ab17ap17Result["state"] != ab17ap17Want {
					ab17ap17T.Fatalf("wrong completeness: %v", ab17ap17Result)
				}
				ab17ap17TopWant := "partial"
				if ab17ap17Want == "complete" {
					ab17ap17TopWant = "complete"
				}
				if ab17ap17Output["state"] != ab17ap17TopWant {
					ab17ap17T.Fatal(ab17ap17Output)
				}
				if ab17ap17Want == "incomplete" {
					if _, ab17ap17Exists := ab17ap17Result["number"]; ab17ap17Exists {
						ab17ap17T.Fatal("invented number")
					}
				}
				if ab17ap17Key != "status" && ab17ap17Fields[ab17ap17Sibling].(map[string]any)["state"] != "complete" {
					ab17ap17T.Fatal("independent sibling lost")
				}
				_, ab17ap17HasRaw := ab17ap17Result["raw"]
				if ab17ap17HasRaw != (ab17ap17Mode != "absent") {
					ab17ap17T.Fatal("missing raw fabricated")
				}
			})
		}
	}
}

func TestAb17Ap17RejectInvalid(ab17ap17T *ab17ap17Testing.T) {
	ab17ap17Valid := string(ab17ap17MustMarshal(ab17ap17T, ab17ap17Base(ab17ap17T)))
	ab17ap17Cases := []string{"", "null", "[]", "{}", `{"context":{}}`, `{"fields":{}}`, "{", ab17ap17Valid + " {}", `{"context":{},"fields":{},"fields":{}}`, `{"context":{"nested":[{"x":1,"x":2}]},"fields":{}}`, `{"context":{},"fields":{},"arguments":{"x":false,"x":true}}`, `{"context":{},"fields":{"data":[],"data":[]}}`, `{"context":{},"fields":{},"extra":{"x":1,"x":2}}`, `{"context":{},"fields":{},"arguments":{"requested_address":0,"requested_\u0061ddress":1}}`}
	for _, ab17ap17Key := range []string{"address", "count", "data", "status", "extra"} {
		for _, ab17ap17Bad := range []string{`null`, `false`, `0`, `"x"`, `{}`, `[-1]`, `[256]`, `[1.5]`, `[true]`, `[null]`, `["1"]`, `[[1]]`, `[1e0]`} {
			ab17ap17Cases = append(ab17ap17Cases, `{"context":{},"fields":{"`+ab17ap17Key+`":`+ab17ap17Bad+`}}`)
		}
	}
	for _, ab17ap17Key := range []string{"primary", "secondary", "operation"} {
		for _, ab17ap17Bad := range []string{`false`, `"181"`, `[]`, `{}`, `-1`, `256`, `181.5`} {
			ab17ap17Cases = append(ab17ap17Cases, `{"context":{"`+ab17ap17Key+`":`+ab17ap17Bad+`},"fields":{}}`)
		}
	}
	for _, ab17ap17Key := range []string{"direction", "view"} {
		for _, ab17ap17Bad := range []string{`false`, `0`, `[]`, `{}`} {
			ab17ap17Cases = append(ab17ap17Cases, `{"context":{"`+ab17ap17Key+`":`+ab17ap17Bad+`},"fields":{}}`)
		}
	}
	for _, ab17ap17Bad := range []string{`null`, `false`, `[]`, `"1"`, `-1`, `4294967296`, `1.2`, `1e0`} {
		ab17ap17Cases = append(ab17ap17Cases, `{"context":{},"fields":{},"arguments":{"requested_address":`+ab17ap17Bad+`}}`)
	}
	for _, ab17ap17Key := range []string{"context", "fields", "arguments"} {
		for _, ab17ap17Bad := range []string{`null`, `false`, `[]`, `0`, `"x"`} {
			ab17ap17Input := ab17ap17Base(ab17ap17T)
			ab17ap17Input[ab17ap17Key] = ab17ap17JSON.RawMessage(ab17ap17Bad)
			ab17ap17Cases = append(ab17ap17Cases, string(ab17ap17MustMarshal(ab17ap17T, ab17ap17Input)))
		}
	}
	for ab17ap17Index, ab17ap17Case := range ab17ap17Cases {
		ab17ap17Output, ab17ap17Err := Ab17Ap17Evaluate([]byte(ab17ap17Case))
		if ab17ap17Err == nil || ab17ap17Output != nil {
			ab17ap17T.Fatalf("invalid case %d accepted: %s", ab17ap17Index, ab17ap17Case)
		}
	}
}

func TestAb17Ap17ArgumentsAndOwnership(ab17ap17T *ab17ap17Testing.T) {
	for _, ab17ap17Address := range []any{nil, 0, 1, 305419896, 4294967295} {
		ab17ap17Input := ab17ap17Base(ab17ap17T)
		ab17ap17Input["context"].(map[string]any)["direction"] = "response"
		ab17ap17Input["fields"].(map[string]any)["data"] = []int{}
		ab17ap17Args := ab17ap17Input["arguments"].(map[string]any)
		if ab17ap17Address == nil {
			delete(ab17ap17Args, "requested_address")
		} else {
			ab17ap17Args["requested_address"] = ab17ap17Address
		}
		ab17ap17Output := ab17ap17Run(ab17ap17T, ab17ap17Input)
		ab17ap17Normalized := ab17ap17Decode(ab17ap17T, ab17ap17MustMarshal(ab17ap17T, ab17ap17Input))
		if !ab17ap17Reflect.DeepEqual(ab17ap17Output["arguments"], ab17ap17Normalized["arguments"]) || ab17ap17Output["correlation"] != "Unknown" {
			ab17ap17T.Fatal("argument used as echo constraint")
		}
	}
	for _, ab17ap17Mode := range []string{"absent", "empty", "false", "nested"} {
		ab17ap17Input := ab17ap17Base(ab17ap17T)
		switch ab17ap17Mode {
		case "absent":
			delete(ab17ap17Input, "arguments")
		case "empty":
			ab17ap17Input["arguments"] = map[string]any{}
		case "false":
			ab17ap17Input["arguments"].(map[string]any)["extra"] = false
		case "nested":
			ab17ap17Input["arguments"].(map[string]any)["extra"] = map[string]any{"list": []any{nil, false, 0}}
		}
		ab17ap17Output := ab17ap17Run(ab17ap17T, ab17ap17Input)
		ab17ap17Normalized := ab17ap17Decode(ab17ap17T, ab17ap17MustMarshal(ab17ap17T, ab17ap17Input))
		if !ab17ap17Reflect.DeepEqual(ab17ap17Output["arguments"], ab17ap17Normalized["arguments"]) {
			ab17ap17T.Fatal("arguments changed")
		}
		_, ab17ap17Exists := ab17ap17Output["arguments"]
		if ab17ap17Exists != (ab17ap17Mode != "absent") {
			ab17ap17T.Fatal("arguments presence lost")
		}
	}
	ab17ap17Input := ab17ap17MustMarshal(ab17ap17T, ab17ap17Base(ab17ap17T))
	ab17ap17Original := append([]byte(nil), ab17ap17Input...)
	ab17ap17First, ab17ap17Err := Ab17Ap17Evaluate(ab17ap17Input)
	if ab17ap17Err != nil {
		ab17ap17T.Fatal(ab17ap17Err)
	}
	ab17ap17Saved := append([]byte(nil), ab17ap17First...)
	for ab17ap17Index := range ab17ap17Input {
		ab17ap17Input[ab17ap17Index] = ' '
	}
	if !ab17ap17Bytes.Equal(ab17ap17Saved, ab17ap17First) {
		ab17ap17T.Fatal("input aliases result")
	}
	ab17ap17Decoded := ab17ap17Decode(ab17ap17T, ab17ap17First)
	ab17ap17Decoded["raw"].(map[string]any)["address"].([]any)[0] = ab17ap17JSON.Number("0")
	ab17ap17Decoded["context"].(map[string]any)["opaque"].(map[string]any)["enabled"] = true
	ab17ap17Decoded["arguments"].(map[string]any)["extra"] = true
	if ab17ap17Decoded["fields"].(map[string]any)["address"].(map[string]any)["raw"].([]any)[0] != ab17ap17JSON.Number("120") {
		ab17ap17T.Fatal("nested raw aliases semantic field")
	}
	for ab17ap17Index := range ab17ap17First {
		ab17ap17First[ab17ap17Index] = ' '
	}
	ab17ap17Second, ab17ap17Err := Ab17Ap17Evaluate(ab17ap17Original)
	if ab17ap17Err != nil || !ab17ap17Bytes.Equal(ab17ap17Saved, ab17ap17Second) {
		ab17ap17T.Fatal("later call changed")
	}
	if ab17ap17Strings.Contains(string(ab17ap17Second), `"label"`) {
		ab17ap17T.Fatal("invented label")
	}
}

func TestAb17Ap17UnmappedFields(ab17ap17T *ab17ap17Testing.T) {
	for _, ab17ap17Raw := range [][]int{{}, {0}, {255}, {1, 2, 3}} {
		ab17ap17Input := ab17ap17Base(ab17ap17T)
		ab17ap17Input["fields"].(map[string]any)["extra"] = ab17ap17Raw
		ab17ap17Output := ab17ap17Run(ab17ap17T, ab17ap17Input)
		ab17ap17Expected := ab17ap17Decode(ab17ap17T, ab17ap17MustMarshal(ab17ap17T, map[string]any{"raw": ab17ap17Raw, "state": "Unknown"}))
		if !ab17ap17Reflect.DeepEqual(ab17ap17Output["fields"].(map[string]any)["extra"], ab17ap17Expected) || ab17ap17Output["state"] != "complete" {
			ab17ap17T.Fatal("unmapped field gained meaning or hid siblings")
		}
	}
}
