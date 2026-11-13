package vaillant

import (
	b510CacheTransitionJSON "encoding/json"
	b510CacheTransitionOs "os"
	b510CacheTransitionReflect "reflect"
	b510CacheTransitionStrings "strings"
	b510CacheTransitionTesting "testing"
)

func b510CacheTransitionInput() map[string]any {
	return map[string]any{
		"context":   map[string]any{"primary": float64(181), "secondary": float64(16), "direction": "local", "view": "capability-query", "family": "synthetic-family", "variant": map[string]any{"opaque": []any{"synthetic-variant"}}},
		"fields":    map[string]any{"prior": []any{float64(165)}, "feedback": []any{float64(1)}, "bit_six": []any{float64(64)}},
		"arguments": map[string]any{"parsed": true, "count": float64(3)},
	}
}

func b510CacheTransitionCall(b510CacheTransitionT *b510CacheTransitionTesting.T, b510CacheTransitionInputValue map[string]any) map[string]any {
	b510CacheTransitionT.Helper()
	b510CacheTransitionEncoded, b510CacheTransitionErr := b510CacheTransitionJSON.Marshal(b510CacheTransitionInputValue)
	if b510CacheTransitionErr != nil {
		b510CacheTransitionT.Fatal(b510CacheTransitionErr)
	}
	b510CacheTransitionResult, b510CacheTransitionErr := B510CacheTransitionInterpret(b510CacheTransitionEncoded)
	if b510CacheTransitionErr != nil {
		b510CacheTransitionT.Fatal(b510CacheTransitionErr)
	}
	var b510CacheTransitionOutput map[string]any
	if b510CacheTransitionErr := b510CacheTransitionJSON.Unmarshal(b510CacheTransitionResult, &b510CacheTransitionOutput); b510CacheTransitionErr != nil {
		b510CacheTransitionT.Fatal(b510CacheTransitionErr)
	}
	if b510CacheTransitionOutput["physical"] != "Unknown" || b510CacheTransitionOutput["qualified"] != false {
		b510CacheTransitionT.Fatalf("native qualification: %s", b510CacheTransitionResult)
	}
	for _, b510CacheTransitionName := range []string{"context", "arguments"} {
		if !b510CacheTransitionReflect.DeepEqual(b510CacheTransitionInputValue[b510CacheTransitionName], b510CacheTransitionOutput[b510CacheTransitionName]) && b510CacheTransitionInputValue[b510CacheTransitionName] != nil {
			b510CacheTransitionT.Fatalf("%s changed", b510CacheTransitionName)
		}
	}
	if b510CacheTransitionInputValue["fields"] != nil && !b510CacheTransitionReflect.DeepEqual(b510CacheTransitionInputValue["fields"], b510CacheTransitionOutput["raw"]) {
		b510CacheTransitionT.Fatal("raw changed")
	}
	return b510CacheTransitionOutput
}

func b510CacheTransitionField(b510CacheTransitionOutput map[string]any, b510CacheTransitionName string) map[string]any {
	return b510CacheTransitionOutput["fields"].(map[string]any)[b510CacheTransitionName].(map[string]any)
}

func TestAb11Ao11Transitions(b510CacheTransitionT *b510CacheTransitionTesting.T) {
	b510CacheTransitionData, b510CacheTransitionErr := b510CacheTransitionOs.ReadFile("testdata/ab11ao11_b510_cache_transition.json")
	if b510CacheTransitionErr != nil {
		b510CacheTransitionT.Fatal(b510CacheTransitionErr)
	}
	var b510CacheTransitionCases []map[string]any
	if b510CacheTransitionErr := b510CacheTransitionJSON.Unmarshal(b510CacheTransitionData, &b510CacheTransitionCases); b510CacheTransitionErr != nil {
		b510CacheTransitionT.Fatal(b510CacheTransitionErr)
	}
	if len(b510CacheTransitionCases) != 13 {
		b510CacheTransitionT.Fatal("missing independent vectors")
	}
	for _, b510CacheTransitionCase := range b510CacheTransitionCases {
		b510CacheTransitionT.Run(b510CacheTransitionCase["id"].(string), func(b510CacheTransitionT *b510CacheTransitionTesting.T) {
			b510CacheTransitionValue := b510CacheTransitionInput()
			for _, b510CacheTransitionName := range []string{"prior", "feedback", "bit_six"} {
				b510CacheTransitionValue["fields"].(map[string]any)[b510CacheTransitionName] = []any{b510CacheTransitionCase[b510CacheTransitionName]}
			}
			for _, b510CacheTransitionName := range []string{"parsed", "count"} {
				b510CacheTransitionValue["arguments"].(map[string]any)[b510CacheTransitionName] = b510CacheTransitionCase[b510CacheTransitionName]
			}
			b510CacheTransitionOutput := b510CacheTransitionCall(b510CacheTransitionT, b510CacheTransitionValue)
			if b510CacheTransitionOutput["state"] != b510CacheTransitionCase["state"] || b510CacheTransitionField(b510CacheTransitionOutput, "cache")["number"] != b510CacheTransitionCase["number"] || b510CacheTransitionField(b510CacheTransitionOutput, "freshness")["label"] != b510CacheTransitionCase["freshness"] || b510CacheTransitionField(b510CacheTransitionOutput, "rule")["label"] != b510CacheTransitionCase["rule"] {
				b510CacheTransitionT.Fatalf("transition: %#v", b510CacheTransitionOutput)
			}
			b510CacheTransitionRuleState := "Hypothesis"
			if b510CacheTransitionCase["rule"] == "Unknown" {
				b510CacheTransitionRuleState = "Unknown"
			}
			if b510CacheTransitionField(b510CacheTransitionOutput, "rule")["state"] != b510CacheTransitionRuleState {
				b510CacheTransitionT.Fatal("rule grade changed")
			}
			for _, b510CacheTransitionName := range []string{"freshness", "rule"} {
				if _, b510CacheTransitionPresent := b510CacheTransitionField(b510CacheTransitionOutput, b510CacheTransitionName)["number"]; b510CacheTransitionPresent {
					b510CacheTransitionT.Fatal("invented numeric result")
				}
			}
			for _, b510CacheTransitionName := range []string{"cache", "freshness"} {
				if b510CacheTransitionField(b510CacheTransitionOutput, b510CacheTransitionName)["state"] != "Hypothesis" {
					b510CacheTransitionT.Fatal("family grade promoted")
				}
			}
			if !b510CacheTransitionReflect.DeepEqual(b510CacheTransitionField(b510CacheTransitionOutput, "cache")["raw"], b510CacheTransitionValue["fields"].(map[string]any)["prior"]) {
				b510CacheTransitionT.Fatal("cache evidence lost")
			}
		})
	}
}

func TestAb11Ao11Selection(b510CacheTransitionT *b510CacheTransitionTesting.T) {
	for _, b510CacheTransitionName := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, b510CacheTransitionMode := range []string{"different", "missing", "zero"} {
			b510CacheTransitionT.Run("ab11ao11_"+b510CacheTransitionName+"_"+b510CacheTransitionMode, func(b510CacheTransitionT *b510CacheTransitionTesting.T) {
				b510CacheTransitionValue := b510CacheTransitionInput()
				b510CacheTransitionContext := b510CacheTransitionValue["context"].(map[string]any)
				switch b510CacheTransitionName {
				case "direction", "view":
					b510CacheTransitionContext[b510CacheTransitionName] = "unknown"
				default:
					b510CacheTransitionContext[b510CacheTransitionName] = float64(255)
				}
				if b510CacheTransitionMode == "missing" {
					delete(b510CacheTransitionContext, b510CacheTransitionName)
				}
				if b510CacheTransitionMode == "zero" {
					if b510CacheTransitionName == "direction" || b510CacheTransitionName == "view" {
						b510CacheTransitionContext[b510CacheTransitionName] = ""
					} else {
						b510CacheTransitionContext[b510CacheTransitionName] = float64(0)
					}
				}
				b510CacheTransitionOutput := b510CacheTransitionCall(b510CacheTransitionT, b510CacheTransitionValue)
				if b510CacheTransitionName == "operation" && b510CacheTransitionMode == "missing" {
					if b510CacheTransitionField(b510CacheTransitionOutput, "cache")["number"] != float64(64) {
						b510CacheTransitionT.Fatal("absent operation control")
					}
				} else if b510CacheTransitionOutput["state"] != "Unknown" || len(b510CacheTransitionOutput["fields"].(map[string]any)) != 0 {
					b510CacheTransitionT.Fatal("guessed selection")
				}
			})
		}
	}
	for _, b510CacheTransitionName := range []string{"family", "variant", "extra"} {
		b510CacheTransitionT.Run("ab11ao11_opaque_"+b510CacheTransitionName, func(b510CacheTransitionT *b510CacheTransitionTesting.T) {
			b510CacheTransitionValue := b510CacheTransitionInput()
			b510CacheTransitionValue["context"].(map[string]any)[b510CacheTransitionName] = map[string]any{"unmodelled": []any{false, float64(0), "other-family"}}
			b510CacheTransitionOutput := b510CacheTransitionCall(b510CacheTransitionT, b510CacheTransitionValue)
			if b510CacheTransitionField(b510CacheTransitionOutput, "cache")["number"] != float64(64) {
				b510CacheTransitionT.Fatal("opaque context selected a view")
			}
		})
	}
}

func TestAb11Ao11Incomplete(b510CacheTransitionT *b510CacheTransitionTesting.T) {
	for _, b510CacheTransitionName := range []string{"prior", "feedback", "bit_six"} {
		for _, b510CacheTransitionMode := range []string{"missing", "empty", "tail"} {
			b510CacheTransitionT.Run("ab11ao11_"+b510CacheTransitionName+"_"+b510CacheTransitionMode, func(b510CacheTransitionT *b510CacheTransitionTesting.T) {
				b510CacheTransitionValue := b510CacheTransitionInput()
				b510CacheTransitionRaw := b510CacheTransitionValue["fields"].(map[string]any)
				switch b510CacheTransitionMode {
				case "missing":
					delete(b510CacheTransitionRaw, b510CacheTransitionName)
				case "empty":
					b510CacheTransitionRaw[b510CacheTransitionName] = []any{}
				case "tail":
					b510CacheTransitionRaw[b510CacheTransitionName] = append(b510CacheTransitionRaw[b510CacheTransitionName].([]any), float64(0))
				}
				b510CacheTransitionOutput := b510CacheTransitionCall(b510CacheTransitionT, b510CacheTransitionValue)
				if b510CacheTransitionOutput["state"] != "Unknown" || len(b510CacheTransitionOutput["fields"].(map[string]any)) != 0 {
					b510CacheTransitionT.Fatal("incomplete transition admitted")
				}
			})
		}
	}
	for _, b510CacheTransitionName := range []string{"parsed", "count", "arguments", "context", "fields"} {
		b510CacheTransitionT.Run("ab11ao11_missing_"+b510CacheTransitionName, func(b510CacheTransitionT *b510CacheTransitionTesting.T) {
			b510CacheTransitionValue := b510CacheTransitionInput()
			if b510CacheTransitionName == "parsed" || b510CacheTransitionName == "count" {
				delete(b510CacheTransitionValue["arguments"].(map[string]any), b510CacheTransitionName)
			} else {
				delete(b510CacheTransitionValue, b510CacheTransitionName)
			}
			b510CacheTransitionOutput := b510CacheTransitionCall(b510CacheTransitionT, b510CacheTransitionValue)
			if b510CacheTransitionOutput["state"] != "Unknown" || len(b510CacheTransitionOutput["fields"].(map[string]any)) != 0 {
				b510CacheTransitionT.Fatal("missing admitted")
			}
			if b510CacheTransitionName == "arguments" {
				if _, b510CacheTransitionPresent := b510CacheTransitionOutput["arguments"]; b510CacheTransitionPresent {
					b510CacheTransitionT.Fatal("invented arguments")
				}
			}
		})
	}
	for _, b510CacheTransitionName := range []string{"negative-count", "invalid-bit", "empty-arguments", "unmodelled"} {
		b510CacheTransitionT.Run("ab11ao11_"+b510CacheTransitionName, func(b510CacheTransitionT *b510CacheTransitionTesting.T) {
			b510CacheTransitionValue := b510CacheTransitionInput()
			switch b510CacheTransitionName {
			case "negative-count":
				b510CacheTransitionValue["arguments"].(map[string]any)["count"] = float64(-1)
			case "invalid-bit":
				b510CacheTransitionValue["fields"].(map[string]any)["bit_six"] = []any{float64(1)}
			case "empty-arguments":
				b510CacheTransitionValue["arguments"] = map[string]any{}
			case "unmodelled":
				b510CacheTransitionValue["fields"].(map[string]any)["unknown"] = []any{float64(255), float64(0)}
			}
			b510CacheTransitionOutput := b510CacheTransitionCall(b510CacheTransitionT, b510CacheTransitionValue)
			if b510CacheTransitionName == "unmodelled" {
				if b510CacheTransitionField(b510CacheTransitionOutput, "cache")["number"] != float64(64) {
					b510CacheTransitionT.Fatal("sibling erased valid transition")
				}
			} else if len(b510CacheTransitionOutput["fields"].(map[string]any)) != 0 {
				b510CacheTransitionT.Fatal("invalid admitted")
			}
		})
	}
}

func TestAb11Ao11Malformed(b510CacheTransitionT *b510CacheTransitionTesting.T) {
	for b510CacheTransitionIndex, b510CacheTransitionText := range []string{
		``, `{`, `{} {}`, `null`, `[]`, `{"context":null}`, `{"fields":null}`, `{"arguments":null}`,
		`{"context":{},"context":{}}`, `{"context":{"family":{"x":1,"x":2}}}`, `{"fields":{"prior":[1],"prior":[2]}}`, `{"arguments":{"parsed":true,"parsed":false}}`,
		`{"context":{"primary":"181"}}`, `{"context":{"secondary":1.5}}`, `{"context":{"operation":null}}`, `{"context":{"direction":0}}`, `{"context":{"view":[]}}`,
		`{"arguments":{"parsed":1}}`, `{"arguments":{"count":true}}`, `{"arguments":{"count":1.0}}`, `{"arguments":{"extra":0}}`,
		`{"fields":{"prior":null}}`, `{"fields":{"prior":"AA=="}}`, `{"fields":{"prior":[-1]}}`, `{"fields":{"prior":[256]}}`, `{"fields":{"prior":[1.5]}}`, `{"fields":{"prior":[1e0]}}`, `{"fields":{"prior":[true]}}`, `{"fields":{"prior":[null]}}`, `{"fields":{"extra":["0"]}}`,
	} {
		b510CacheTransitionResult, b510CacheTransitionErr := B510CacheTransitionInterpret([]byte(b510CacheTransitionText))
		if b510CacheTransitionErr == nil || b510CacheTransitionResult != nil {
			b510CacheTransitionT.Fatalf("case %d accepted: %s", b510CacheTransitionIndex, b510CacheTransitionText)
		}
	}
}

func TestAb11Ao11IsolationAndLargeCount(b510CacheTransitionT *b510CacheTransitionTesting.T) {
	b510CacheTransitionValue := b510CacheTransitionInput()
	b510CacheTransitionEncoded, _ := b510CacheTransitionJSON.Marshal(b510CacheTransitionValue)
	b510CacheTransitionFirst, b510CacheTransitionErr := B510CacheTransitionInterpret(b510CacheTransitionEncoded)
	if b510CacheTransitionErr != nil {
		b510CacheTransitionT.Fatal(b510CacheTransitionErr)
	}
	b510CacheTransitionSaved := string(b510CacheTransitionFirst)
	for b510CacheTransitionIndex := range b510CacheTransitionEncoded {
		b510CacheTransitionEncoded[b510CacheTransitionIndex] = ' '
	}
	if string(b510CacheTransitionFirst) != b510CacheTransitionSaved {
		b510CacheTransitionT.Fatal("input alias")
	}
	var b510CacheTransitionDecoded map[string]any
	if b510CacheTransitionErr := b510CacheTransitionJSON.Unmarshal(b510CacheTransitionFirst, &b510CacheTransitionDecoded); b510CacheTransitionErr != nil {
		b510CacheTransitionT.Fatal(b510CacheTransitionErr)
	}
	b510CacheTransitionField(b510CacheTransitionDecoded, "cache")["raw"].([]any)[0] = float64(9)
	if b510CacheTransitionDecoded["raw"].(map[string]any)["prior"].([]any)[0] != float64(165) {
		b510CacheTransitionT.Fatal("nested raw alias")
	}
	b510CacheTransitionDecoded["context"].(map[string]any)["variant"].(map[string]any)["opaque"].([]any)[0] = "changed"
	b510CacheTransitionDecoded["arguments"].(map[string]any)["count"] = float64(0)
	for b510CacheTransitionIndex := range b510CacheTransitionFirst {
		b510CacheTransitionFirst[b510CacheTransitionIndex] = ' '
	}
	b510CacheTransitionSecondInput, _ := b510CacheTransitionJSON.Marshal(b510CacheTransitionValue)
	b510CacheTransitionSecond, b510CacheTransitionErr := B510CacheTransitionInterpret(b510CacheTransitionSecondInput)
	if b510CacheTransitionErr != nil || string(b510CacheTransitionSecond) != b510CacheTransitionSaved {
		b510CacheTransitionT.Fatal("cross-call mutation")
	}
	b510CacheTransitionHuge := b510CacheTransitionStrings.Replace(string(b510CacheTransitionSecondInput), `"count":3`, `"count":184467440737095516160000`, 1)
	b510CacheTransitionResult, b510CacheTransitionErr := B510CacheTransitionInterpret([]byte(b510CacheTransitionHuge))
	if b510CacheTransitionErr != nil {
		b510CacheTransitionT.Fatal(b510CacheTransitionErr)
	}
	if b510CacheTransitionErr := b510CacheTransitionJSON.Unmarshal(b510CacheTransitionResult, &b510CacheTransitionDecoded); b510CacheTransitionErr != nil {
		b510CacheTransitionT.Fatal(b510CacheTransitionErr)
	}
	if b510CacheTransitionField(b510CacheTransitionDecoded, "cache")["number"] != float64(64) {
		b510CacheTransitionT.Fatal("count narrowed")
	}
}

func TestAb11Ao11KindsAndOpaquePresence(b510CacheTransitionT *b510CacheTransitionTesting.T) {
	for _, b510CacheTransitionName := range []string{"prior", "feedback", "bit_six", "unmodelled"} {
		for _, b510CacheTransitionBad := range []any{nil, "AA==", []any{nil}, []any{true}, []any{float64(-1)}, []any{float64(256)}, []any{float64(0.5)}} {
			b510CacheTransitionValue := b510CacheTransitionInput()
			b510CacheTransitionValue["fields"].(map[string]any)[b510CacheTransitionName] = b510CacheTransitionBad
			b510CacheTransitionEncoded, _ := b510CacheTransitionJSON.Marshal(b510CacheTransitionValue)
			b510CacheTransitionResult, b510CacheTransitionErr := B510CacheTransitionInterpret(b510CacheTransitionEncoded)
			if b510CacheTransitionErr == nil || b510CacheTransitionResult != nil {
				b510CacheTransitionT.Fatalf("bad field kind accepted: %s", b510CacheTransitionEncoded)
			}
		}
	}
	for _, b510CacheTransitionName := range []string{"family", "variant"} {
		for _, b510CacheTransitionMode := range []string{"missing", "null", "empty", "zero"} {
			b510CacheTransitionT.Run("ab11ao11_"+b510CacheTransitionName+"_"+b510CacheTransitionMode, func(b510CacheTransitionT *b510CacheTransitionTesting.T) {
				b510CacheTransitionValue := b510CacheTransitionInput()
				b510CacheTransitionContext := b510CacheTransitionValue["context"].(map[string]any)
				switch b510CacheTransitionMode {
				case "missing":
					delete(b510CacheTransitionContext, b510CacheTransitionName)
				case "null":
					b510CacheTransitionContext[b510CacheTransitionName] = nil
				case "empty":
					b510CacheTransitionContext[b510CacheTransitionName] = ""
				case "zero":
					b510CacheTransitionContext[b510CacheTransitionName] = float64(0)
				}
				b510CacheTransitionOutput := b510CacheTransitionCall(b510CacheTransitionT, b510CacheTransitionValue)
				if b510CacheTransitionField(b510CacheTransitionOutput, "cache")["number"] != float64(64) {
					b510CacheTransitionT.Fatal("opaque presence affects selection")
				}
			})
		}
	}
}
