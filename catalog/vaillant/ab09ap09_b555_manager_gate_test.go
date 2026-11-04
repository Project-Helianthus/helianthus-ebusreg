package vaillant

import (
	ab09ap09json "encoding/json"
	ab09ap09os "os"
	ab09ap09reflect "reflect"
	ab09ap09strings "strings"
	ab09ap09testing "testing"
)

func ab09ap09Input() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 85, "operation": 0, "direction": "local", "view": "manager-gate", "family": "synthetic-family", "variant": "synthetic-variant", "opaque": map[string]any{"zero": 0, "null": nil, "false": false}}, "fields": map[string]any{"reply": []any{3}, "independent": []any{0, 255}}, "arguments": map[string]any{"controller_type": 2, "spn": 385, "family_match": true, "opaque": []any{false, nil, 0}}}
}

func ab09ap09Evaluate(ab09ap09T *ab09ap09testing.T, ab09ap09Input map[string]any) map[string]any {
	ab09ap09T.Helper()
	ab09ap09Bytes, ab09ap09Err := ab09ap09json.Marshal(ab09ap09Input)
	if ab09ap09Err != nil {
		ab09ap09T.Fatal(ab09ap09Err)
	}
	ab09ap09Result, ab09ap09Err := Ab09Ap09Evaluate(ab09ap09Bytes)
	if ab09ap09Err != nil {
		ab09ap09T.Fatal(ab09ap09Err)
	}
	ab09ap09Output := map[string]any{}
	if ab09ap09Err := ab09ap09json.Unmarshal(ab09ap09Result, &ab09ap09Output); ab09ap09Err != nil {
		ab09ap09T.Fatal(ab09ap09Err)
	}
	for ab09ap09OutputKey, ab09ap09InputKey := range map[string]string{"context": "context", "raw": "fields", "arguments": "arguments"} {
		ab09ap09ExpectedBytes, _ := ab09ap09json.Marshal(ab09ap09Input[ab09ap09InputKey])
		ab09ap09Expected := any(nil)
		_ = ab09ap09json.Unmarshal(ab09ap09ExpectedBytes, &ab09ap09Expected)
		if !ab09ap09reflect.DeepEqual(ab09ap09Output[ab09ap09OutputKey], ab09ap09Expected) {
			ab09ap09T.Fatalf("%s not preserved", ab09ap09OutputKey)
		}
	}
	if ab09ap09Output["physical"] != "Unknown" || ab09ap09Output["qualified"] != false {
		ab09ap09T.Fatal("qualification boundary", ab09ap09Output)
	}
	return ab09ap09Output
}

func ab09ap09Reply(ab09ap09Output map[string]any) map[string]any {
	return ab09ap09Output["fields"].(map[string]any)["reply"].(map[string]any)
}

func TestAb09Ap09SyntheticCases(ab09ap09T *ab09ap09testing.T) {
	ab09ap09Bytes, ab09ap09Err := ab09ap09os.ReadFile("testdata/ab09ap09_b555_manager_gate.json")
	if ab09ap09Err != nil {
		ab09ap09T.Fatal(ab09ap09Err)
	}
	ab09ap09Cases := []map[string]any{}
	if ab09ap09Err := ab09ap09json.Unmarshal(ab09ap09Bytes, &ab09ap09Cases); ab09ap09Err != nil {
		ab09ap09T.Fatal(ab09ap09Err)
	}
	if len(ab09ap09Cases) != 12 {
		ab09ap09T.Fatal("fixture case count")
	}
	for _, ab09ap09Case := range ab09ap09Cases {
		ab09ap09T.Run(ab09ap09Case["id"].(string), func(ab09ap09T *ab09ap09testing.T) {
			ab09ap09Request := ab09ap09Input()
			ab09ap09Request["fields"].(map[string]any)["reply"] = ab09ap09Case["reply"]
			for _, ab09ap09Key := range []string{"controller_type", "spn", "family_match"} {
				ab09ap09Request["arguments"].(map[string]any)[ab09ap09Key] = ab09ap09Case[ab09ap09Key]
			}
			ab09ap09Output := ab09ap09Evaluate(ab09ap09T, ab09ap09Request)
			ab09ap09Result := ab09ap09Reply(ab09ap09Output)
			if ab09ap09Output["state"] != "complete" || ab09ap09Result["state"] != "Hypothesis" || ab09ap09Result["calculation"] != "Proven" {
				ab09ap09T.Fatal("grades", ab09ap09Output)
			}
			ab09ap09Flags := ab09ap09Result["flags"].(map[string]any)
			for _, ab09ap09Key := range []string{"bit0", "bit1"} {
				if ab09ap09Flags[ab09ap09Key] != ab09ap09Case[ab09ap09Key] {
					ab09ap09T.Fatal(ab09ap09Key, ab09ap09Result)
				}
			}
			if ab09ap09Result["eligible"] != ab09ap09Case["eligible"] || ab09ap09Result["first_bit_policy"] != ab09ap09Case["bit0"] || ab09ap09Result["unmodeled"] != ab09ap09Case["high"] {
				ab09ap09T.Fatal("predicates/high bits", ab09ap09Result)
			}
			if ab09ap09Result["family_role"] != float64(21) || ab09ap09Result["route_role"] != float64(8) || ab09ap09Result["first_bit_policy_label"] != "caller predicate" {
				ab09ap09T.Fatal("roles/policy", ab09ap09Result)
			}
			ab09ap09Expected := "Hypothesis"
			if ab09ap09Case["family_match"] == nil {
				ab09ap09Expected = "Unknown"
			}
			if ab09ap09Result["family_applicability"] != ab09ap09Expected {
				ab09ap09T.Fatal("family applicability", ab09ap09Result)
			}
		})
	}
}

func TestAb09Ap09Context(ab09ap09T *ab09ap09testing.T) {
	for ab09ap09Key, ab09ap09Other := range map[string]any{"primary": 180, "secondary": 84, "operation": 1, "direction": "reply", "view": "other", "family": "other-family", "variant": "other-variant", "opaque": []any{nil, 0, false}} {
		for _, ab09ap09Mode := range []string{"other", "missing", "null", "zero"} {
			ab09ap09T.Run(ab09ap09Key+"_"+ab09ap09Mode, func(ab09ap09T *ab09ap09testing.T) {
				ab09ap09Request := ab09ap09Input()
				ab09ap09Context := ab09ap09Request["context"].(map[string]any)
				switch ab09ap09Mode {
				case "other":
					ab09ap09Context[ab09ap09Key] = ab09ap09Other
				case "missing":
					delete(ab09ap09Context, ab09ap09Key)
				case "null":
					ab09ap09Context[ab09ap09Key] = nil
				case "zero":
					if ab09ap09Key == "direction" || ab09ap09Key == "view" {
						ab09ap09Context[ab09ap09Key] = ""
					} else {
						ab09ap09Context[ab09ap09Key] = 0
					}
				}
				ab09ap09Output := ab09ap09Evaluate(ab09ap09T, ab09ap09Request)
				ab09ap09Known := ab09ap09Key == "family" || ab09ap09Key == "variant" || ab09ap09Key == "opaque" || (ab09ap09Key == "operation" && ab09ap09Mode == "zero")
				if ab09ap09Known {
					if ab09ap09Output["state"] != "complete" {
						ab09ap09T.Fatal(ab09ap09Output)
					}
				} else if ab09ap09Output["state"] != "Unknown" || len(ab09ap09Output["fields"].(map[string]any)) != 0 {
					ab09ap09T.Fatal("unsupported context interpreted", ab09ap09Output)
				}
			})
		}
	}
}

func TestAb09Ap09IncompleteAndIndependentRaw(ab09ap09T *ab09ap09testing.T) {
	for _, ab09ap09Mode := range []string{"missing", "empty", "long"} {
		ab09ap09T.Run(ab09ap09Mode, func(ab09ap09T *ab09ap09testing.T) {
			ab09ap09Request := ab09ap09Input()
			ab09ap09Fields := ab09ap09Request["fields"].(map[string]any)
			switch ab09ap09Mode {
			case "missing":
				delete(ab09ap09Fields, "reply")
			case "empty":
				ab09ap09Fields["reply"] = []any{}
			case "long":
				ab09ap09Fields["reply"] = []any{3, 0}
			}
			ab09ap09Output := ab09ap09Evaluate(ab09ap09T, ab09ap09Request)
			ab09ap09Result := ab09ap09Reply(ab09ap09Output)
			if ab09ap09Output["state"] != "partial" || ab09ap09Result["state"] != "incomplete" {
				ab09ap09T.Fatal(ab09ap09Output)
			}
			for _, ab09ap09Key := range []string{"eligible", "number", "label", "flags"} {
				if _, ab09ap09Present := ab09ap09Result[ab09ap09Key]; ab09ap09Present {
					ab09ap09T.Fatal("invented incomplete result", ab09ap09Result)
				}
			}
		})
	}
	for _, ab09ap09Value := range [][]any{{}, {0}, {255, 0}} {
		ab09ap09Request := ab09ap09Input()
		ab09ap09Request["fields"].(map[string]any)["independent"] = ab09ap09Value
		if ab09ap09Evaluate(ab09ap09T, ab09ap09Request)["state"] != "complete" {
			ab09ap09T.Fatal("independent raw field affected reply")
		}
	}
	ab09ap09Request := ab09ap09Input()
	delete(ab09ap09Request["fields"].(map[string]any), "independent")
	ab09ap09Evaluate(ab09ap09T, ab09ap09Request)
}

func TestAb09Ap09RejectInvalid(ab09ap09T *ab09ap09testing.T) {
	ab09ap09Base, _ := ab09ap09json.Marshal(ab09ap09Input())
	ab09ap09Invalid := []string{"", "{", "[]", "null", string(ab09ap09Base) + " {}", `{"context":{},"context":{},"fields":{}}`, `{"context":{"x":1,"x":2},"fields":{}}`, `{"context":{},"fields":{"reply":[],"reply":[]}}`, `{"context":{},"fields":{},"arguments":{"x":[{"a":0,"a":1}]}}`}
	for _, ab09ap09Key := range []string{"context", "fields", "arguments"} {
		for _, ab09ap09Value := range []any{nil, false, 0, "", []any{}} {
			ab09ap09Request := ab09ap09Input()
			ab09ap09Request[ab09ap09Key] = ab09ap09Value
			ab09ap09Bytes, _ := ab09ap09json.Marshal(ab09ap09Request)
			ab09ap09Invalid = append(ab09ap09Invalid, string(ab09ap09Bytes))
		}
	}
	for _, ab09ap09Key := range []string{"controller_type", "spn", "family_match"} {
		ab09ap09Request := ab09ap09Input()
		delete(ab09ap09Request["arguments"].(map[string]any), ab09ap09Key)
		ab09ap09Bytes, _ := ab09ap09json.Marshal(ab09ap09Request)
		ab09ap09Invalid = append(ab09ap09Invalid, string(ab09ap09Bytes))
		ab09ap09Values := []any{nil, false, "2", []any{}, map[string]any{}, -1, 1.5, 256}
		if ab09ap09Key == "spn" {
			ab09ap09Values = []any{nil, false, "385", []any{}, map[string]any{}, -1, 1.5, 65536}
		}
		if ab09ap09Key == "family_match" {
			ab09ap09Values = []any{0, 1, "true", []any{}, map[string]any{}}
		}
		for _, ab09ap09Value := range ab09ap09Values {
			ab09ap09Request := ab09ap09Input()
			ab09ap09Request["arguments"].(map[string]any)[ab09ap09Key] = ab09ap09Value
			ab09ap09Bytes, _ := ab09ap09json.Marshal(ab09ap09Request)
			ab09ap09Invalid = append(ab09ap09Invalid, string(ab09ap09Bytes))
		}
	}
	for _, ab09ap09Value := range []any{nil, false, 0, "", map[string]any{}, []any{nil}, []any{false}, []any{"3"}, []any{-1}, []any{256}, []any{1.5}} {
		for _, ab09ap09Key := range []string{"reply", "independent"} {
			ab09ap09Request := ab09ap09Input()
			ab09ap09Request["fields"].(map[string]any)[ab09ap09Key] = ab09ap09Value
			ab09ap09Bytes, _ := ab09ap09json.Marshal(ab09ap09Request)
			ab09ap09Invalid = append(ab09ap09Invalid, string(ab09ap09Bytes))
		}
	}
	for _, ab09ap09Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		ab09ap09Values := []any{false, "181", -1, 256, 1.5, []any{}}
		if ab09ap09Key == "direction" || ab09ap09Key == "view" {
			ab09ap09Values = []any{false, 0, []any{}}
		}
		for _, ab09ap09Value := range ab09ap09Values {
			ab09ap09Request := ab09ap09Input()
			ab09ap09Request["context"].(map[string]any)[ab09ap09Key] = ab09ap09Value
			ab09ap09Bytes, _ := ab09ap09json.Marshal(ab09ap09Request)
			ab09ap09Invalid = append(ab09ap09Invalid, string(ab09ap09Bytes))
		}
	}
	for ab09ap09Index, ab09ap09Input := range ab09ap09Invalid {
		ab09ap09Result, ab09ap09Err := Ab09Ap09Evaluate([]byte(ab09ap09Input))
		if ab09ap09Err == nil || ab09ap09Result != nil {
			ab09ap09T.Fatalf("invalid case %d accepted: %s", ab09ap09Index, ab09ap09Input)
		}
	}
}

func TestAb09Ap09IsolationAndBounds(ab09ap09T *ab09ap09testing.T) {
	ab09ap09Request := ab09ap09Input()
	ab09ap09InputBytes, _ := ab09ap09json.Marshal(ab09ap09Request)
	ab09ap09First, ab09ap09Err := Ab09Ap09Evaluate(ab09ap09InputBytes)
	if ab09ap09Err != nil {
		ab09ap09T.Fatal(ab09ap09Err)
	}
	ab09ap09Saved := string(ab09ap09First)
	for ab09ap09Index := range ab09ap09InputBytes {
		ab09ap09InputBytes[ab09ap09Index] = ' '
	}
	if string(ab09ap09First) != ab09ap09Saved {
		ab09ap09T.Fatal("input aliases result")
	}
	ab09ap09Output := map[string]any{}
	_ = ab09ap09json.Unmarshal(ab09ap09First, &ab09ap09Output)
	ab09ap09Output["raw"].(map[string]any)["reply"].([]any)[0] = 99
	ab09ap09Output["context"].(map[string]any)["opaque"].(map[string]any)["zero"] = 99
	ab09ap09Output["arguments"].(map[string]any)["opaque"].([]any)[0] = true
	if ab09ap09Reply(ab09ap09Output)["raw"].([]any)[0] != float64(3) {
		ab09ap09T.Fatal("raw aliases field result")
	}
	for ab09ap09Index := range ab09ap09First {
		ab09ap09First[ab09ap09Index] = ' '
	}
	ab09ap09Fresh, _ := ab09ap09json.Marshal(ab09ap09Input())
	ab09ap09Second, ab09ap09Err := Ab09Ap09Evaluate(ab09ap09Fresh)
	if ab09ap09Err != nil || string(ab09ap09Second) != ab09ap09Saved {
		ab09ap09T.Fatal("call storage reused")
	}
	for _, ab09ap09Key := range []string{"controller_type", "spn"} {
		ab09ap09Request := ab09ap09Input()
		ab09ap09Max := 255
		if ab09ap09Key == "spn" {
			ab09ap09Max = 65535
		}
		ab09ap09Request["arguments"].(map[string]any)[ab09ap09Key] = ab09ap09Max
		if ab09ap09Reply(ab09ap09Evaluate(ab09ap09T, ab09ap09Request))["eligible"] != false {
			ab09ap09T.Fatal("maximum accepted as match")
		}
	}
	ab09ap09Request = ab09ap09Input()
	ab09ap09Request["arguments"].(map[string]any)["opaque"] = map[string]any{"array": []any{false, nil, 0}, "number": 4294967295}
	ab09ap09Evaluate(ab09ap09T, ab09ap09Request)
	delete(ab09ap09Request["arguments"].(map[string]any), "opaque")
	ab09ap09Evaluate(ab09ap09T, ab09ap09Request)
	ab09ap09Request["context"].(map[string]any)["view"] = "unknown"
	delete(ab09ap09Request, "arguments")
	ab09ap09Bytes, _ := ab09ap09json.Marshal(ab09ap09Request)
	ab09ap09Result, ab09ap09Err := Ab09Ap09Evaluate(ab09ap09Bytes)
	if ab09ap09Err != nil || ab09ap09strings.Contains(string(ab09ap09Result), `"arguments"`) {
		ab09ap09T.Fatal("absent arguments changed", ab09ap09Err)
	}
}

func TestAb09Ap09EveryReplyBit(ab09ap09T *ab09ap09testing.T) {
	for _, ab09ap09High := range []int{4, 8, 16, 32, 64, 128, 252} {
		ab09ap09Request := ab09ap09Input()
		ab09ap09Request["fields"].(map[string]any)["reply"] = []any{ab09ap09High + 3}
		ab09ap09Result := ab09ap09Reply(ab09ap09Evaluate(ab09ap09T, ab09ap09Request))
		if ab09ap09Result["eligible"] != true || ab09ap09Result["unmodeled"] != float64(ab09ap09High) {
			ab09ap09T.Fatal("high bits changed eligibility or were discarded", ab09ap09Result)
		}
	}
}
