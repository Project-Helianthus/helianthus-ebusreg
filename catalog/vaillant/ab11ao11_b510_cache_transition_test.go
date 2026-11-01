package vaillant

import (
	ab11ao11json "encoding/json"
	ab11ao11os "os"
	ab11ao11reflect "reflect"
	ab11ao11strings "strings"
	ab11ao11testing "testing"
)

func ab11ao11Input() map[string]any {
	return map[string]any{
		"context":   map[string]any{"primary": float64(181), "secondary": float64(16), "direction": "local", "view": "capability-query", "family": "synthetic-family", "variant": map[string]any{"opaque": []any{"synthetic-variant"}}},
		"fields":    map[string]any{"prior": []any{float64(165)}, "feedback": []any{float64(1)}, "bit_six": []any{float64(64)}},
		"arguments": map[string]any{"parsed": true, "count": float64(3)},
	}
}

func ab11ao11Call(ab11ao11T *ab11ao11testing.T, ab11ao11InputValue map[string]any) map[string]any {
	ab11ao11T.Helper()
	ab11ao11Encoded, ab11ao11Err := ab11ao11json.Marshal(ab11ao11InputValue)
	if ab11ao11Err != nil {
		ab11ao11T.Fatal(ab11ao11Err)
	}
	ab11ao11Result, ab11ao11Err := Ab11Ao11Interpret(ab11ao11Encoded)
	if ab11ao11Err != nil {
		ab11ao11T.Fatal(ab11ao11Err)
	}
	var ab11ao11Output map[string]any
	if ab11ao11Err := ab11ao11json.Unmarshal(ab11ao11Result, &ab11ao11Output); ab11ao11Err != nil {
		ab11ao11T.Fatal(ab11ao11Err)
	}
	if ab11ao11Output["physical"] != "Unknown" || ab11ao11Output["qualified"] != false {
		ab11ao11T.Fatalf("native qualification: %s", ab11ao11Result)
	}
	for _, ab11ao11Name := range []string{"context", "arguments"} {
		if !ab11ao11reflect.DeepEqual(ab11ao11InputValue[ab11ao11Name], ab11ao11Output[ab11ao11Name]) && ab11ao11InputValue[ab11ao11Name] != nil {
			ab11ao11T.Fatalf("%s changed", ab11ao11Name)
		}
	}
	if ab11ao11InputValue["fields"] != nil && !ab11ao11reflect.DeepEqual(ab11ao11InputValue["fields"], ab11ao11Output["raw"]) {
		ab11ao11T.Fatal("raw changed")
	}
	return ab11ao11Output
}

func ab11ao11Field(ab11ao11Output map[string]any, ab11ao11Name string) map[string]any {
	return ab11ao11Output["fields"].(map[string]any)[ab11ao11Name].(map[string]any)
}

func TestAb11Ao11Transitions(ab11ao11T *ab11ao11testing.T) {
	ab11ao11Data, ab11ao11Err := ab11ao11os.ReadFile("testdata/ab11ao11_b510_cache_transition.json")
	if ab11ao11Err != nil {
		ab11ao11T.Fatal(ab11ao11Err)
	}
	var ab11ao11Cases []map[string]any
	if ab11ao11Err := ab11ao11json.Unmarshal(ab11ao11Data, &ab11ao11Cases); ab11ao11Err != nil {
		ab11ao11T.Fatal(ab11ao11Err)
	}
	if len(ab11ao11Cases) != 13 {
		ab11ao11T.Fatal("missing independent vectors")
	}
	for _, ab11ao11Case := range ab11ao11Cases {
		ab11ao11T.Run(ab11ao11Case["id"].(string), func(ab11ao11T *ab11ao11testing.T) {
			ab11ao11Value := ab11ao11Input()
			for _, ab11ao11Name := range []string{"prior", "feedback", "bit_six"} {
				ab11ao11Value["fields"].(map[string]any)[ab11ao11Name] = []any{ab11ao11Case[ab11ao11Name]}
			}
			for _, ab11ao11Name := range []string{"parsed", "count"} {
				ab11ao11Value["arguments"].(map[string]any)[ab11ao11Name] = ab11ao11Case[ab11ao11Name]
			}
			ab11ao11Output := ab11ao11Call(ab11ao11T, ab11ao11Value)
			if ab11ao11Output["state"] != ab11ao11Case["state"] || ab11ao11Field(ab11ao11Output, "cache")["number"] != ab11ao11Case["number"] || ab11ao11Field(ab11ao11Output, "freshness")["label"] != ab11ao11Case["freshness"] || ab11ao11Field(ab11ao11Output, "rule")["label"] != ab11ao11Case["rule"] {
				ab11ao11T.Fatalf("transition: %#v", ab11ao11Output)
			}
			ab11ao11RuleState := "Hypothesis"
			if ab11ao11Case["rule"] == "Unknown" {
				ab11ao11RuleState = "Unknown"
			}
			if ab11ao11Field(ab11ao11Output, "rule")["state"] != ab11ao11RuleState {
				ab11ao11T.Fatal("rule grade changed")
			}
			for _, ab11ao11Name := range []string{"freshness", "rule"} {
				if _, ab11ao11Present := ab11ao11Field(ab11ao11Output, ab11ao11Name)["number"]; ab11ao11Present {
					ab11ao11T.Fatal("invented numeric result")
				}
			}
			for _, ab11ao11Name := range []string{"cache", "freshness"} {
				if ab11ao11Field(ab11ao11Output, ab11ao11Name)["state"] != "Hypothesis" {
					ab11ao11T.Fatal("family grade promoted")
				}
			}
			if !ab11ao11reflect.DeepEqual(ab11ao11Field(ab11ao11Output, "cache")["raw"], ab11ao11Value["fields"].(map[string]any)["prior"]) {
				ab11ao11T.Fatal("cache evidence lost")
			}
		})
	}
}

func TestAb11Ao11Selection(ab11ao11T *ab11ao11testing.T) {
	for _, ab11ao11Name := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, ab11ao11Mode := range []string{"different", "missing", "zero"} {
			ab11ao11T.Run("ab11ao11_"+ab11ao11Name+"_"+ab11ao11Mode, func(ab11ao11T *ab11ao11testing.T) {
				ab11ao11Value := ab11ao11Input()
				ab11ao11Context := ab11ao11Value["context"].(map[string]any)
				switch ab11ao11Name {
				case "direction", "view":
					ab11ao11Context[ab11ao11Name] = "unknown"
				default:
					ab11ao11Context[ab11ao11Name] = float64(255)
				}
				if ab11ao11Mode == "missing" {
					delete(ab11ao11Context, ab11ao11Name)
				}
				if ab11ao11Mode == "zero" {
					if ab11ao11Name == "direction" || ab11ao11Name == "view" {
						ab11ao11Context[ab11ao11Name] = ""
					} else {
						ab11ao11Context[ab11ao11Name] = float64(0)
					}
				}
				ab11ao11Output := ab11ao11Call(ab11ao11T, ab11ao11Value)
				if ab11ao11Name == "operation" && ab11ao11Mode == "missing" {
					if ab11ao11Field(ab11ao11Output, "cache")["number"] != float64(64) {
						ab11ao11T.Fatal("absent operation control")
					}
				} else if ab11ao11Output["state"] != "Unknown" || len(ab11ao11Output["fields"].(map[string]any)) != 0 {
					ab11ao11T.Fatal("guessed selection")
				}
			})
		}
	}
	for _, ab11ao11Name := range []string{"family", "variant", "extra"} {
		ab11ao11T.Run("ab11ao11_opaque_"+ab11ao11Name, func(ab11ao11T *ab11ao11testing.T) {
			ab11ao11Value := ab11ao11Input()
			ab11ao11Value["context"].(map[string]any)[ab11ao11Name] = map[string]any{"unmodelled": []any{false, float64(0), "other-family"}}
			ab11ao11Output := ab11ao11Call(ab11ao11T, ab11ao11Value)
			if ab11ao11Field(ab11ao11Output, "cache")["number"] != float64(64) {
				ab11ao11T.Fatal("opaque context selected a view")
			}
		})
	}
}

func TestAb11Ao11Incomplete(ab11ao11T *ab11ao11testing.T) {
	for _, ab11ao11Name := range []string{"prior", "feedback", "bit_six"} {
		for _, ab11ao11Mode := range []string{"missing", "empty", "tail"} {
			ab11ao11T.Run("ab11ao11_"+ab11ao11Name+"_"+ab11ao11Mode, func(ab11ao11T *ab11ao11testing.T) {
				ab11ao11Value := ab11ao11Input()
				ab11ao11Raw := ab11ao11Value["fields"].(map[string]any)
				switch ab11ao11Mode {
				case "missing":
					delete(ab11ao11Raw, ab11ao11Name)
				case "empty":
					ab11ao11Raw[ab11ao11Name] = []any{}
				case "tail":
					ab11ao11Raw[ab11ao11Name] = append(ab11ao11Raw[ab11ao11Name].([]any), float64(0))
				}
				ab11ao11Output := ab11ao11Call(ab11ao11T, ab11ao11Value)
				if ab11ao11Output["state"] != "Unknown" || len(ab11ao11Output["fields"].(map[string]any)) != 0 {
					ab11ao11T.Fatal("incomplete transition admitted")
				}
			})
		}
	}
	for _, ab11ao11Name := range []string{"parsed", "count", "arguments", "context", "fields"} {
		ab11ao11T.Run("ab11ao11_missing_"+ab11ao11Name, func(ab11ao11T *ab11ao11testing.T) {
			ab11ao11Value := ab11ao11Input()
			if ab11ao11Name == "parsed" || ab11ao11Name == "count" {
				delete(ab11ao11Value["arguments"].(map[string]any), ab11ao11Name)
			} else {
				delete(ab11ao11Value, ab11ao11Name)
			}
			ab11ao11Output := ab11ao11Call(ab11ao11T, ab11ao11Value)
			if ab11ao11Output["state"] != "Unknown" || len(ab11ao11Output["fields"].(map[string]any)) != 0 {
				ab11ao11T.Fatal("missing admitted")
			}
			if ab11ao11Name == "arguments" {
				if _, ab11ao11Present := ab11ao11Output["arguments"]; ab11ao11Present {
					ab11ao11T.Fatal("invented arguments")
				}
			}
		})
	}
	for _, ab11ao11Name := range []string{"negative-count", "invalid-bit", "empty-arguments", "unmodelled"} {
		ab11ao11T.Run("ab11ao11_"+ab11ao11Name, func(ab11ao11T *ab11ao11testing.T) {
			ab11ao11Value := ab11ao11Input()
			switch ab11ao11Name {
			case "negative-count":
				ab11ao11Value["arguments"].(map[string]any)["count"] = float64(-1)
			case "invalid-bit":
				ab11ao11Value["fields"].(map[string]any)["bit_six"] = []any{float64(1)}
			case "empty-arguments":
				ab11ao11Value["arguments"] = map[string]any{}
			case "unmodelled":
				ab11ao11Value["fields"].(map[string]any)["unknown"] = []any{float64(255), float64(0)}
			}
			ab11ao11Output := ab11ao11Call(ab11ao11T, ab11ao11Value)
			if ab11ao11Name == "unmodelled" {
				if ab11ao11Field(ab11ao11Output, "cache")["number"] != float64(64) {
					ab11ao11T.Fatal("sibling erased valid transition")
				}
			} else if len(ab11ao11Output["fields"].(map[string]any)) != 0 {
				ab11ao11T.Fatal("invalid admitted")
			}
		})
	}
}

func TestAb11Ao11Malformed(ab11ao11T *ab11ao11testing.T) {
	for ab11ao11Index, ab11ao11Text := range []string{
		``, `{`, `{} {}`, `null`, `[]`, `{"context":null}`, `{"fields":null}`, `{"arguments":null}`,
		`{"context":{},"context":{}}`, `{"context":{"family":{"x":1,"x":2}}}`, `{"fields":{"prior":[1],"prior":[2]}}`, `{"arguments":{"parsed":true,"parsed":false}}`,
		`{"context":{"primary":"181"}}`, `{"context":{"secondary":1.5}}`, `{"context":{"operation":null}}`, `{"context":{"direction":0}}`, `{"context":{"view":[]}}`,
		`{"arguments":{"parsed":1}}`, `{"arguments":{"count":true}}`, `{"arguments":{"count":1.0}}`, `{"arguments":{"extra":0}}`,
		`{"fields":{"prior":null}}`, `{"fields":{"prior":"AA=="}}`, `{"fields":{"prior":[-1]}}`, `{"fields":{"prior":[256]}}`, `{"fields":{"prior":[1.5]}}`, `{"fields":{"prior":[1e0]}}`, `{"fields":{"prior":[true]}}`, `{"fields":{"prior":[null]}}`, `{"fields":{"extra":["0"]}}`,
	} {
		ab11ao11Result, ab11ao11Err := Ab11Ao11Interpret([]byte(ab11ao11Text))
		if ab11ao11Err == nil || ab11ao11Result != nil {
			ab11ao11T.Fatalf("case %d accepted: %s", ab11ao11Index, ab11ao11Text)
		}
	}
}

func TestAb11Ao11IsolationAndLargeCount(ab11ao11T *ab11ao11testing.T) {
	ab11ao11Value := ab11ao11Input()
	ab11ao11Encoded, _ := ab11ao11json.Marshal(ab11ao11Value)
	ab11ao11First, ab11ao11Err := Ab11Ao11Interpret(ab11ao11Encoded)
	if ab11ao11Err != nil {
		ab11ao11T.Fatal(ab11ao11Err)
	}
	ab11ao11Saved := string(ab11ao11First)
	for ab11ao11Index := range ab11ao11Encoded {
		ab11ao11Encoded[ab11ao11Index] = ' '
	}
	if string(ab11ao11First) != ab11ao11Saved {
		ab11ao11T.Fatal("input alias")
	}
	var ab11ao11Decoded map[string]any
	if ab11ao11Err := ab11ao11json.Unmarshal(ab11ao11First, &ab11ao11Decoded); ab11ao11Err != nil {
		ab11ao11T.Fatal(ab11ao11Err)
	}
	ab11ao11Field(ab11ao11Decoded, "cache")["raw"].([]any)[0] = float64(9)
	if ab11ao11Decoded["raw"].(map[string]any)["prior"].([]any)[0] != float64(165) {
		ab11ao11T.Fatal("nested raw alias")
	}
	ab11ao11Decoded["context"].(map[string]any)["variant"].(map[string]any)["opaque"].([]any)[0] = "changed"
	ab11ao11Decoded["arguments"].(map[string]any)["count"] = float64(0)
	for ab11ao11Index := range ab11ao11First {
		ab11ao11First[ab11ao11Index] = ' '
	}
	ab11ao11SecondInput, _ := ab11ao11json.Marshal(ab11ao11Value)
	ab11ao11Second, ab11ao11Err := Ab11Ao11Interpret(ab11ao11SecondInput)
	if ab11ao11Err != nil || string(ab11ao11Second) != ab11ao11Saved {
		ab11ao11T.Fatal("cross-call mutation")
	}
	ab11ao11Huge := ab11ao11strings.Replace(string(ab11ao11SecondInput), `"count":3`, `"count":184467440737095516160000`, 1)
	ab11ao11Result, ab11ao11Err := Ab11Ao11Interpret([]byte(ab11ao11Huge))
	if ab11ao11Err != nil {
		ab11ao11T.Fatal(ab11ao11Err)
	}
	if ab11ao11Err := ab11ao11json.Unmarshal(ab11ao11Result, &ab11ao11Decoded); ab11ao11Err != nil {
		ab11ao11T.Fatal(ab11ao11Err)
	}
	if ab11ao11Field(ab11ao11Decoded, "cache")["number"] != float64(64) {
		ab11ao11T.Fatal("count narrowed")
	}
}

func TestAb11Ao11KindsAndOpaquePresence(ab11ao11T *ab11ao11testing.T) {
	for _, ab11ao11Name := range []string{"prior", "feedback", "bit_six", "unmodelled"} {
		for _, ab11ao11Bad := range []any{nil, "AA==", []any{nil}, []any{true}, []any{float64(-1)}, []any{float64(256)}, []any{float64(0.5)}} {
			ab11ao11Value := ab11ao11Input()
			ab11ao11Value["fields"].(map[string]any)[ab11ao11Name] = ab11ao11Bad
			ab11ao11Encoded, _ := ab11ao11json.Marshal(ab11ao11Value)
			ab11ao11Result, ab11ao11Err := Ab11Ao11Interpret(ab11ao11Encoded)
			if ab11ao11Err == nil || ab11ao11Result != nil {
				ab11ao11T.Fatalf("bad field kind accepted: %s", ab11ao11Encoded)
			}
		}
	}
	for _, ab11ao11Name := range []string{"family", "variant"} {
		for _, ab11ao11Mode := range []string{"missing", "null", "empty", "zero"} {
			ab11ao11T.Run("ab11ao11_"+ab11ao11Name+"_"+ab11ao11Mode, func(ab11ao11T *ab11ao11testing.T) {
				ab11ao11Value := ab11ao11Input()
				ab11ao11Context := ab11ao11Value["context"].(map[string]any)
				switch ab11ao11Mode {
				case "missing":
					delete(ab11ao11Context, ab11ao11Name)
				case "null":
					ab11ao11Context[ab11ao11Name] = nil
				case "empty":
					ab11ao11Context[ab11ao11Name] = ""
				case "zero":
					ab11ao11Context[ab11ao11Name] = float64(0)
				}
				ab11ao11Output := ab11ao11Call(ab11ao11T, ab11ao11Value)
				if ab11ao11Field(ab11ao11Output, "cache")["number"] != float64(64) {
					ab11ao11T.Fatal("opaque presence affects selection")
				}
			})
		}
	}
}
