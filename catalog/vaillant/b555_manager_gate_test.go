package vaillant

import (
	b555ManagerGateJSON "encoding/json"
	b555ManagerGateOs "os"
	b555ManagerGateReflect "reflect"
	b555ManagerGateStrings "strings"
	b555ManagerGateTesting "testing"
)

func b555ManagerGateInput() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 85, "operation": 0, "direction": "local", "view": "manager-gate", "family": "synthetic-family", "variant": "synthetic-variant", "opaque": map[string]any{"zero": 0, "null": nil, "false": false}}, "fields": map[string]any{"reply": []any{3}, "independent": []any{0, 255}}, "arguments": map[string]any{"controller_type": 2, "spn": 385, "family_match": true, "opaque": []any{false, nil, 0}}}
}

func b555ManagerGateEvaluate(b555ManagerGateT *b555ManagerGateTesting.T, b555ManagerGateInput map[string]any) map[string]any {
	b555ManagerGateT.Helper()
	b555ManagerGateBytes, b555ManagerGateErr := b555ManagerGateJSON.Marshal(b555ManagerGateInput)
	if b555ManagerGateErr != nil {
		b555ManagerGateT.Fatal(b555ManagerGateErr)
	}
	b555ManagerGateResult, b555ManagerGateErr := B555ManagerGateEvaluate(b555ManagerGateBytes)
	if b555ManagerGateErr != nil {
		b555ManagerGateT.Fatal(b555ManagerGateErr)
	}
	b555ManagerGateOutput := map[string]any{}
	if b555ManagerGateErr := b555ManagerGateJSON.Unmarshal(b555ManagerGateResult, &b555ManagerGateOutput); b555ManagerGateErr != nil {
		b555ManagerGateT.Fatal(b555ManagerGateErr)
	}
	for b555ManagerGateOutputKey, b555ManagerGateInputKey := range map[string]string{"context": "context", "raw": "fields", "arguments": "arguments"} {
		b555ManagerGateExpectedBytes, _ := b555ManagerGateJSON.Marshal(b555ManagerGateInput[b555ManagerGateInputKey])
		b555ManagerGateExpected := any(nil)
		_ = b555ManagerGateJSON.Unmarshal(b555ManagerGateExpectedBytes, &b555ManagerGateExpected)
		if !b555ManagerGateReflect.DeepEqual(b555ManagerGateOutput[b555ManagerGateOutputKey], b555ManagerGateExpected) {
			b555ManagerGateT.Fatalf("%s not preserved", b555ManagerGateOutputKey)
		}
	}
	if b555ManagerGateOutput["physical"] != "Unknown" || b555ManagerGateOutput["qualified"] != false {
		b555ManagerGateT.Fatal("qualification boundary", b555ManagerGateOutput)
	}
	return b555ManagerGateOutput
}

func b555ManagerGateReply(b555ManagerGateOutput map[string]any) map[string]any {
	return b555ManagerGateOutput["fields"].(map[string]any)["reply"].(map[string]any)
}

func TestAb09Ap09SyntheticCases(b555ManagerGateT *b555ManagerGateTesting.T) {
	b555ManagerGateBytes, b555ManagerGateErr := b555ManagerGateOs.ReadFile("testdata/ab09ap09_b555_manager_gate.json")
	if b555ManagerGateErr != nil {
		b555ManagerGateT.Fatal(b555ManagerGateErr)
	}
	b555ManagerGateCases := []map[string]any{}
	if b555ManagerGateErr := b555ManagerGateJSON.Unmarshal(b555ManagerGateBytes, &b555ManagerGateCases); b555ManagerGateErr != nil {
		b555ManagerGateT.Fatal(b555ManagerGateErr)
	}
	if len(b555ManagerGateCases) != 12 {
		b555ManagerGateT.Fatal("fixture case count")
	}
	for _, b555ManagerGateCase := range b555ManagerGateCases {
		b555ManagerGateT.Run(b555ManagerGateCase["id"].(string), func(b555ManagerGateT *b555ManagerGateTesting.T) {
			b555ManagerGateRequest := b555ManagerGateInput()
			b555ManagerGateRequest["fields"].(map[string]any)["reply"] = b555ManagerGateCase["reply"]
			for _, b555ManagerGateKey := range []string{"controller_type", "spn", "family_match"} {
				b555ManagerGateRequest["arguments"].(map[string]any)[b555ManagerGateKey] = b555ManagerGateCase[b555ManagerGateKey]
			}
			b555ManagerGateOutput := b555ManagerGateEvaluate(b555ManagerGateT, b555ManagerGateRequest)
			b555ManagerGateResult := b555ManagerGateReply(b555ManagerGateOutput)
			if b555ManagerGateOutput["state"] != "complete" || b555ManagerGateResult["state"] != "Hypothesis" || b555ManagerGateResult["calculation"] != "Proven" {
				b555ManagerGateT.Fatal("grades", b555ManagerGateOutput)
			}
			b555ManagerGateFlags := b555ManagerGateResult["flags"].(map[string]any)
			for _, b555ManagerGateKey := range []string{"bit0", "bit1"} {
				if b555ManagerGateFlags[b555ManagerGateKey] != b555ManagerGateCase[b555ManagerGateKey] {
					b555ManagerGateT.Fatal(b555ManagerGateKey, b555ManagerGateResult)
				}
			}
			if b555ManagerGateResult["eligible"] != b555ManagerGateCase["eligible"] || b555ManagerGateResult["first_bit_policy"] != b555ManagerGateCase["bit0"] || b555ManagerGateResult["unmodeled"] != b555ManagerGateCase["high"] {
				b555ManagerGateT.Fatal("predicates/high bits", b555ManagerGateResult)
			}
			if b555ManagerGateResult["family_role"] != float64(21) || b555ManagerGateResult["route_role"] != float64(8) || b555ManagerGateResult["first_bit_policy_label"] != "caller predicate" {
				b555ManagerGateT.Fatal("roles/policy", b555ManagerGateResult)
			}
			b555ManagerGateExpected := "Hypothesis"
			if b555ManagerGateCase["family_match"] == nil {
				b555ManagerGateExpected = "Unknown"
			}
			if b555ManagerGateResult["family_applicability"] != b555ManagerGateExpected {
				b555ManagerGateT.Fatal("family applicability", b555ManagerGateResult)
			}
		})
	}
}

func TestAb09Ap09Context(b555ManagerGateT *b555ManagerGateTesting.T) {
	for b555ManagerGateKey, b555ManagerGateOther := range map[string]any{"primary": 180, "secondary": 84, "operation": 1, "direction": "reply", "view": "other", "family": "other-family", "variant": "other-variant", "opaque": []any{nil, 0, false}} {
		for _, b555ManagerGateMode := range []string{"other", "missing", "null", "zero"} {
			b555ManagerGateT.Run(b555ManagerGateKey+"_"+b555ManagerGateMode, func(b555ManagerGateT *b555ManagerGateTesting.T) {
				b555ManagerGateRequest := b555ManagerGateInput()
				b555ManagerGateContext := b555ManagerGateRequest["context"].(map[string]any)
				switch b555ManagerGateMode {
				case "other":
					b555ManagerGateContext[b555ManagerGateKey] = b555ManagerGateOther
				case "missing":
					delete(b555ManagerGateContext, b555ManagerGateKey)
				case "null":
					b555ManagerGateContext[b555ManagerGateKey] = nil
				case "zero":
					if b555ManagerGateKey == "direction" || b555ManagerGateKey == "view" {
						b555ManagerGateContext[b555ManagerGateKey] = ""
					} else {
						b555ManagerGateContext[b555ManagerGateKey] = 0
					}
				}
				b555ManagerGateOutput := b555ManagerGateEvaluate(b555ManagerGateT, b555ManagerGateRequest)
				b555ManagerGateKnown := b555ManagerGateKey == "family" || b555ManagerGateKey == "variant" || b555ManagerGateKey == "opaque" || (b555ManagerGateKey == "operation" && b555ManagerGateMode == "zero")
				if b555ManagerGateKnown {
					if b555ManagerGateOutput["state"] != "complete" {
						b555ManagerGateT.Fatal(b555ManagerGateOutput)
					}
				} else if b555ManagerGateOutput["state"] != "Unknown" || len(b555ManagerGateOutput["fields"].(map[string]any)) != 0 {
					b555ManagerGateT.Fatal("unsupported context interpreted", b555ManagerGateOutput)
				}
			})
		}
	}
}

func TestAb09Ap09IncompleteAndIndependentRaw(b555ManagerGateT *b555ManagerGateTesting.T) {
	for _, b555ManagerGateMode := range []string{"missing", "empty", "long"} {
		b555ManagerGateT.Run(b555ManagerGateMode, func(b555ManagerGateT *b555ManagerGateTesting.T) {
			b555ManagerGateRequest := b555ManagerGateInput()
			b555ManagerGateFields := b555ManagerGateRequest["fields"].(map[string]any)
			switch b555ManagerGateMode {
			case "missing":
				delete(b555ManagerGateFields, "reply")
			case "empty":
				b555ManagerGateFields["reply"] = []any{}
			case "long":
				b555ManagerGateFields["reply"] = []any{3, 0}
			}
			b555ManagerGateOutput := b555ManagerGateEvaluate(b555ManagerGateT, b555ManagerGateRequest)
			b555ManagerGateResult := b555ManagerGateReply(b555ManagerGateOutput)
			if b555ManagerGateOutput["state"] != "partial" || b555ManagerGateResult["state"] != "incomplete" {
				b555ManagerGateT.Fatal(b555ManagerGateOutput)
			}
			for _, b555ManagerGateKey := range []string{"eligible", "number", "label", "flags"} {
				if _, b555ManagerGatePresent := b555ManagerGateResult[b555ManagerGateKey]; b555ManagerGatePresent {
					b555ManagerGateT.Fatal("invented incomplete result", b555ManagerGateResult)
				}
			}
		})
	}
	for _, b555ManagerGateValue := range [][]any{{}, {0}, {255, 0}} {
		b555ManagerGateRequest := b555ManagerGateInput()
		b555ManagerGateRequest["fields"].(map[string]any)["independent"] = b555ManagerGateValue
		if b555ManagerGateEvaluate(b555ManagerGateT, b555ManagerGateRequest)["state"] != "complete" {
			b555ManagerGateT.Fatal("independent raw field affected reply")
		}
	}
	b555ManagerGateRequest := b555ManagerGateInput()
	delete(b555ManagerGateRequest["fields"].(map[string]any), "independent")
	b555ManagerGateEvaluate(b555ManagerGateT, b555ManagerGateRequest)
}

func TestAb09Ap09RejectInvalid(b555ManagerGateT *b555ManagerGateTesting.T) {
	b555ManagerGateBase, _ := b555ManagerGateJSON.Marshal(b555ManagerGateInput())
	b555ManagerGateInvalid := []string{"", "{", "[]", "null", string(b555ManagerGateBase) + " {}", `{"context":{},"context":{},"fields":{}}`, `{"context":{"x":1,"x":2},"fields":{}}`, `{"context":{},"fields":{"reply":[],"reply":[]}}`, `{"context":{},"fields":{},"arguments":{"x":[{"a":0,"a":1}]}}`}
	for _, b555ManagerGateKey := range []string{"context", "fields", "arguments"} {
		for _, b555ManagerGateValue := range []any{nil, false, 0, "", []any{}} {
			b555ManagerGateRequest := b555ManagerGateInput()
			b555ManagerGateRequest[b555ManagerGateKey] = b555ManagerGateValue
			b555ManagerGateBytes, _ := b555ManagerGateJSON.Marshal(b555ManagerGateRequest)
			b555ManagerGateInvalid = append(b555ManagerGateInvalid, string(b555ManagerGateBytes))
		}
	}
	for _, b555ManagerGateKey := range []string{"controller_type", "spn", "family_match"} {
		b555ManagerGateRequest := b555ManagerGateInput()
		delete(b555ManagerGateRequest["arguments"].(map[string]any), b555ManagerGateKey)
		b555ManagerGateBytes, _ := b555ManagerGateJSON.Marshal(b555ManagerGateRequest)
		b555ManagerGateInvalid = append(b555ManagerGateInvalid, string(b555ManagerGateBytes))
		b555ManagerGateValues := []any{nil, false, "2", []any{}, map[string]any{}, -1, 1.5, 256}
		if b555ManagerGateKey == "spn" {
			b555ManagerGateValues = []any{nil, false, "385", []any{}, map[string]any{}, -1, 1.5, 65536}
		}
		if b555ManagerGateKey == "family_match" {
			b555ManagerGateValues = []any{0, 1, "true", []any{}, map[string]any{}}
		}
		for _, b555ManagerGateValue := range b555ManagerGateValues {
			b555ManagerGateRequest := b555ManagerGateInput()
			b555ManagerGateRequest["arguments"].(map[string]any)[b555ManagerGateKey] = b555ManagerGateValue
			b555ManagerGateBytes, _ := b555ManagerGateJSON.Marshal(b555ManagerGateRequest)
			b555ManagerGateInvalid = append(b555ManagerGateInvalid, string(b555ManagerGateBytes))
		}
	}
	for _, b555ManagerGateValue := range []any{nil, false, 0, "", map[string]any{}, []any{nil}, []any{false}, []any{"3"}, []any{-1}, []any{256}, []any{1.5}} {
		for _, b555ManagerGateKey := range []string{"reply", "independent"} {
			b555ManagerGateRequest := b555ManagerGateInput()
			b555ManagerGateRequest["fields"].(map[string]any)[b555ManagerGateKey] = b555ManagerGateValue
			b555ManagerGateBytes, _ := b555ManagerGateJSON.Marshal(b555ManagerGateRequest)
			b555ManagerGateInvalid = append(b555ManagerGateInvalid, string(b555ManagerGateBytes))
		}
	}
	for _, b555ManagerGateKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		b555ManagerGateValues := []any{false, "181", -1, 256, 1.5, []any{}}
		if b555ManagerGateKey == "direction" || b555ManagerGateKey == "view" {
			b555ManagerGateValues = []any{false, 0, []any{}}
		}
		for _, b555ManagerGateValue := range b555ManagerGateValues {
			b555ManagerGateRequest := b555ManagerGateInput()
			b555ManagerGateRequest["context"].(map[string]any)[b555ManagerGateKey] = b555ManagerGateValue
			b555ManagerGateBytes, _ := b555ManagerGateJSON.Marshal(b555ManagerGateRequest)
			b555ManagerGateInvalid = append(b555ManagerGateInvalid, string(b555ManagerGateBytes))
		}
	}
	for b555ManagerGateIndex, b555ManagerGateInput := range b555ManagerGateInvalid {
		b555ManagerGateResult, b555ManagerGateErr := B555ManagerGateEvaluate([]byte(b555ManagerGateInput))
		if b555ManagerGateErr == nil || b555ManagerGateResult != nil {
			b555ManagerGateT.Fatalf("invalid case %d accepted: %s", b555ManagerGateIndex, b555ManagerGateInput)
		}
	}
}

func TestAb09Ap09IsolationAndBounds(b555ManagerGateT *b555ManagerGateTesting.T) {
	b555ManagerGateRequest := b555ManagerGateInput()
	b555ManagerGateInputBytes, _ := b555ManagerGateJSON.Marshal(b555ManagerGateRequest)
	b555ManagerGateFirst, b555ManagerGateErr := B555ManagerGateEvaluate(b555ManagerGateInputBytes)
	if b555ManagerGateErr != nil {
		b555ManagerGateT.Fatal(b555ManagerGateErr)
	}
	b555ManagerGateSaved := string(b555ManagerGateFirst)
	for b555ManagerGateIndex := range b555ManagerGateInputBytes {
		b555ManagerGateInputBytes[b555ManagerGateIndex] = ' '
	}
	if string(b555ManagerGateFirst) != b555ManagerGateSaved {
		b555ManagerGateT.Fatal("input aliases result")
	}
	b555ManagerGateOutput := map[string]any{}
	_ = b555ManagerGateJSON.Unmarshal(b555ManagerGateFirst, &b555ManagerGateOutput)
	b555ManagerGateOutput["raw"].(map[string]any)["reply"].([]any)[0] = 99
	b555ManagerGateOutput["context"].(map[string]any)["opaque"].(map[string]any)["zero"] = 99
	b555ManagerGateOutput["arguments"].(map[string]any)["opaque"].([]any)[0] = true
	if b555ManagerGateReply(b555ManagerGateOutput)["raw"].([]any)[0] != float64(3) {
		b555ManagerGateT.Fatal("raw aliases field result")
	}
	for b555ManagerGateIndex := range b555ManagerGateFirst {
		b555ManagerGateFirst[b555ManagerGateIndex] = ' '
	}
	b555ManagerGateFresh, _ := b555ManagerGateJSON.Marshal(b555ManagerGateInput())
	b555ManagerGateSecond, b555ManagerGateErr := B555ManagerGateEvaluate(b555ManagerGateFresh)
	if b555ManagerGateErr != nil || string(b555ManagerGateSecond) != b555ManagerGateSaved {
		b555ManagerGateT.Fatal("call storage reused")
	}
	for _, b555ManagerGateKey := range []string{"controller_type", "spn"} {
		b555ManagerGateRequest := b555ManagerGateInput()
		b555ManagerGateMax := 255
		if b555ManagerGateKey == "spn" {
			b555ManagerGateMax = 65535
		}
		b555ManagerGateRequest["arguments"].(map[string]any)[b555ManagerGateKey] = b555ManagerGateMax
		if b555ManagerGateReply(b555ManagerGateEvaluate(b555ManagerGateT, b555ManagerGateRequest))["eligible"] != false {
			b555ManagerGateT.Fatal("maximum accepted as match")
		}
	}
	b555ManagerGateRequest = b555ManagerGateInput()
	b555ManagerGateRequest["arguments"].(map[string]any)["opaque"] = map[string]any{"array": []any{false, nil, 0}, "number": 4294967295}
	b555ManagerGateEvaluate(b555ManagerGateT, b555ManagerGateRequest)
	delete(b555ManagerGateRequest["arguments"].(map[string]any), "opaque")
	b555ManagerGateEvaluate(b555ManagerGateT, b555ManagerGateRequest)
	b555ManagerGateRequest["context"].(map[string]any)["view"] = "unknown"
	delete(b555ManagerGateRequest, "arguments")
	b555ManagerGateBytes, _ := b555ManagerGateJSON.Marshal(b555ManagerGateRequest)
	b555ManagerGateResult, b555ManagerGateErr := B555ManagerGateEvaluate(b555ManagerGateBytes)
	if b555ManagerGateErr != nil || b555ManagerGateStrings.Contains(string(b555ManagerGateResult), `"arguments"`) {
		b555ManagerGateT.Fatal("absent arguments changed", b555ManagerGateErr)
	}
}

func TestAb09Ap09EveryReplyBit(b555ManagerGateT *b555ManagerGateTesting.T) {
	for _, b555ManagerGateHigh := range []int{4, 8, 16, 32, 64, 128, 252} {
		b555ManagerGateRequest := b555ManagerGateInput()
		b555ManagerGateRequest["fields"].(map[string]any)["reply"] = []any{b555ManagerGateHigh + 3}
		b555ManagerGateResult := b555ManagerGateReply(b555ManagerGateEvaluate(b555ManagerGateT, b555ManagerGateRequest))
		if b555ManagerGateResult["eligible"] != true || b555ManagerGateResult["unmodeled"] != float64(b555ManagerGateHigh) {
			b555ManagerGateT.Fatal("high bits changed eligibility or were discarded", b555ManagerGateResult)
		}
	}
}
