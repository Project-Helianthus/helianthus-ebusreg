package vaillant

import (
	b516CycleReducerJSON "encoding/json"
	b516CycleReducerOs "os"
	b516CycleReducerReflect "reflect"
	b516CycleReducerStrings "strings"
	b516CycleReducerTesting "testing"
)

func b516CycleReducerInput() map[string]any {
	b516CycleReducerSlots := []any{}
	for b516CycleReducerIndex := 0; b516CycleReducerIndex < 9; b516CycleReducerIndex++ {
		b516CycleReducerSlots = append(b516CycleReducerSlots, map[string]any{"registered": true, "capable": true, "prior": 10, "prior_previous": 7, "prior_changed": false, "outcome": "success", "value": 0})
	}
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 22, "operation": 14, "direction": "local", "view": "cycle", "family": "synthetic-family", "variant": "candidate"}, "fields": map[string]any{"opaque": []any{0, 255}}, "arguments": map[string]any{"slots": b516CycleReducerSlots, "failure_counter": 19}}
}
func b516CycleReducerRun(b516CycleReducerT *b516CycleReducerTesting.T, b516CycleReducerInputValue map[string]any) map[string]any {
	b516CycleReducerT.Helper()
	b516CycleReducerEncoded, b516CycleReducerErr := b516CycleReducerJSON.Marshal(b516CycleReducerInputValue)
	if b516CycleReducerErr != nil {
		b516CycleReducerT.Fatal(b516CycleReducerErr)
	}
	b516CycleReducerOutput, b516CycleReducerErr := B516CycleReducerEvaluate(b516CycleReducerEncoded)
	if b516CycleReducerErr != nil {
		b516CycleReducerT.Fatal(b516CycleReducerErr)
	}
	b516CycleReducerDecoded := map[string]any{}
	if b516CycleReducerErr = b516CycleReducerJSON.Unmarshal(b516CycleReducerOutput, &b516CycleReducerDecoded); b516CycleReducerErr != nil {
		b516CycleReducerT.Fatal(b516CycleReducerErr)
	}
	return b516CycleReducerDecoded
}
func b516CycleReducerArgs(b516CycleReducerInputValue map[string]any) map[string]any {
	return b516CycleReducerInputValue["arguments"].(map[string]any)
}
func b516CycleReducerSlot(b516CycleReducerInputValue map[string]any, b516CycleReducerIndex int) map[string]any {
	return b516CycleReducerArgs(b516CycleReducerInputValue)["slots"].([]any)[b516CycleReducerIndex].(map[string]any)
}
func b516CycleReducerFields(b516CycleReducerOutput map[string]any) map[string]any {
	return b516CycleReducerOutput["fields"].(map[string]any)
}
func b516CycleReducerOutcome(b516CycleReducerOutput map[string]any, b516CycleReducerIndex int) map[string]any {
	return b516CycleReducerFields(b516CycleReducerOutput)["slots"].(map[string]any)["outcomes"].([]any)[b516CycleReducerIndex].(map[string]any)
}
func b516CycleReducerEqual(b516CycleReducerT *b516CycleReducerTesting.T, b516CycleReducerActual, b516CycleReducerExpected any) {
	b516CycleReducerT.Helper()
	if !b516CycleReducerReflect.DeepEqual(b516CycleReducerActual, b516CycleReducerExpected) {
		b516CycleReducerT.Fatalf("got %#v want %#v", b516CycleReducerActual, b516CycleReducerExpected)
	}
}
func b516CycleReducerNumber(b516CycleReducerValue any) any {
	return b516CycleReducerValue.(map[string]any)["number"]
}
func b516CycleReducerReject(b516CycleReducerT *b516CycleReducerTesting.T, b516CycleReducerInputValue map[string]any) {
	b516CycleReducerT.Helper()
	b516CycleReducerEncoded, _ := b516CycleReducerJSON.Marshal(b516CycleReducerInputValue)
	b516CycleReducerResult, b516CycleReducerErr := B516CycleReducerEvaluate(b516CycleReducerEncoded)
	if b516CycleReducerErr == nil || b516CycleReducerResult != nil {
		b516CycleReducerT.Fatalf("accepted invalid input: %s", b516CycleReducerEncoded)
	}
}

func TestAb05Ap05Transitions(b516CycleReducerT *b516CycleReducerTesting.T) {
	b516CycleReducerData, b516CycleReducerErr := b516CycleReducerOs.ReadFile("testdata/ab05ap05_b516_cycle_reducer.json")
	if b516CycleReducerErr != nil {
		b516CycleReducerT.Fatal(b516CycleReducerErr)
	}
	var b516CycleReducerCases []map[string]any
	if b516CycleReducerErr = b516CycleReducerJSON.Unmarshal(b516CycleReducerData, &b516CycleReducerCases); b516CycleReducerErr != nil {
		b516CycleReducerT.Fatal(b516CycleReducerErr)
	}
	for _, b516CycleReducerCase := range b516CycleReducerCases {
		for b516CycleReducerIndex := 0; b516CycleReducerIndex < 9; b516CycleReducerIndex++ {
			b516CycleReducerInputValue := b516CycleReducerInput()
			b516CycleReducerSlotValue := b516CycleReducerSlot(b516CycleReducerInputValue, b516CycleReducerIndex)
			b516CycleReducerSlotValue["prior"] = b516CycleReducerCase["prior"]
			b516CycleReducerSlotValue["value"] = b516CycleReducerCase["value"]
			b516CycleReducerOutput := b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
			b516CycleReducerResult := b516CycleReducerOutcome(b516CycleReducerOutput, b516CycleReducerIndex)
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerResult["changed"].(map[string]any)["value"], b516CycleReducerCase["changed"])
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerResult["current"]), b516CycleReducerCase["value"])
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerResult["previous"]), b516CycleReducerCase["prior"])
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerFields(b516CycleReducerOutput)["sum"]), b516CycleReducerCase["value"])
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["qualified"], false)
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerResult["state"], "Hypothesis")
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerResult["current"].(map[string]any)["state"], "Hypothesis")
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerResult["current"].(map[string]any)["calculation"], "Proven")
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["physical"], "Unknown")
		}
	}
}
func TestAb05Ap05EligibilityOutcomes(b516CycleReducerT *b516CycleReducerTesting.T) {
	for b516CycleReducerIndex := 0; b516CycleReducerIndex < 9; b516CycleReducerIndex++ {
		for _, b516CycleReducerRegistered := range []bool{false, true} {
			for _, b516CycleReducerCapable := range []bool{false, true} {
				for _, b516CycleReducerOutcomeValue := range []string{"success", "read-failure", "parse-failure"} {
					for _, b516CycleReducerChanged := range []bool{false, true} {
						b516CycleReducerInputValue := b516CycleReducerInput()
						b516CycleReducerSlotValue := b516CycleReducerSlot(b516CycleReducerInputValue, b516CycleReducerIndex)
						b516CycleReducerSlotValue["registered"] = b516CycleReducerRegistered
						b516CycleReducerSlotValue["capable"] = b516CycleReducerCapable
						b516CycleReducerSlotValue["prior_changed"] = b516CycleReducerChanged
						b516CycleReducerSlotValue["prior_previous"] = 123
						b516CycleReducerSlotValue["outcome"] = b516CycleReducerOutcomeValue
						b516CycleReducerSlotValue["value"] = 3
						if b516CycleReducerOutcomeValue != "success" {
							delete(b516CycleReducerSlotValue, "value")
						}
						b516CycleReducerOutput := b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
						b516CycleReducerResult := b516CycleReducerOutcome(b516CycleReducerOutput, b516CycleReducerIndex)
						b516CycleReducerCounts := b516CycleReducerFields(b516CycleReducerOutput)["counts"].(map[string]any)
						b516CycleReducerAdmitted := b516CycleReducerRegistered && b516CycleReducerCapable
						b516CycleReducerExpectedOutcome := "skipped"
						if b516CycleReducerAdmitted {
							b516CycleReducerExpectedOutcome = b516CycleReducerOutcomeValue
						}
						b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerResult["outcome"], b516CycleReducerExpectedOutcome)
						b516CycleReducerExpectedSuccessful := float64(8)
						if b516CycleReducerAdmitted && b516CycleReducerOutcomeValue == "success" {
							b516CycleReducerExpectedSuccessful = 9
						}
						b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerCounts["successful"], b516CycleReducerExpectedSuccessful)
						b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerResult["attempted"], b516CycleReducerAdmitted)
						b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerResult["prior_changed"], b516CycleReducerChanged)
						if !b516CycleReducerAdmitted {
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerResult["previous"]), float64(123))
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerCounts["skipped"], float64(1))
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerCounts["failed"], float64(0))
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerCounts["attempted"], float64(8))
						} else {
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerResult["previous"]), float64(10))
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerCounts["attempted"], float64(9))
						}
						if !b516CycleReducerAdmitted || b516CycleReducerOutcomeValue != "success" {
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerResult["current"]), float64(4294967295))
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerResult["changed"].(map[string]any)["state"], "Unknown")
							if _, b516CycleReducerPresent := b516CycleReducerResult["changed"].(map[string]any)["value"]; b516CycleReducerPresent {
								b516CycleReducerT.Fatal("invented changed value")
							}
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["partial_cycle"], true)
						} else {
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerResult["current"]), float64(3))
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["partial_cycle"], false)
						}
						if b516CycleReducerAdmitted && b516CycleReducerOutcomeValue != "success" {
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerCounts["failed"], float64(1))
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["state"], "partial")
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerFields(b516CycleReducerOutput)["failure_counter"].(map[string]any)["prior"], float64(19))
						} else {
							b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerFields(b516CycleReducerOutput)["failure_counter"]), float64(0))
						}
						b516CycleReducerExpectedSum := float64(0)
						if b516CycleReducerAdmitted && b516CycleReducerOutcomeValue == "success" {
							b516CycleReducerExpectedSum = 3
						}
						b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerFields(b516CycleReducerOutput)["sum"]), b516CycleReducerExpectedSum)
					}
				}
			}
		}
	}
}
func TestAb05Ap05WrapAndPartial(b516CycleReducerT *b516CycleReducerTesting.T) {
	b516CycleReducerInputValue := b516CycleReducerInput()
	b516CycleReducerSlot(b516CycleReducerInputValue, 0)["value"] = uint32(4294967294)
	b516CycleReducerSlot(b516CycleReducerInputValue, 1)["value"] = 3
	b516CycleReducerSlot(b516CycleReducerInputValue, 2)["outcome"] = "read-failure"
	delete(b516CycleReducerSlot(b516CycleReducerInputValue, 2), "value")
	b516CycleReducerArgs(b516CycleReducerInputValue)["failure_counter"] = uint32(4294967295)
	b516CycleReducerArgs(b516CycleReducerInputValue)["counter_increment"] = 1
	b516CycleReducerOutput := b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
	b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerFields(b516CycleReducerOutput)["sum"]), float64(1))
	b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerFields(b516CycleReducerOutput)["failure_counter"]), float64(0))
	delete(b516CycleReducerArgs(b516CycleReducerInputValue), "counter_increment")
	b516CycleReducerOutput = b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
	b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerFields(b516CycleReducerOutput)["failure_counter"].(map[string]any)["prior"], float64(4294967295))
	b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["state"], "partial")
	b516CycleReducerSlot(b516CycleReducerInputValue, 2)["capable"] = false
	b516CycleReducerArgs(b516CycleReducerInputValue)["counter_increment"] = 7
	b516CycleReducerOutput = b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
	b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerFields(b516CycleReducerOutput)["failure_counter"]), float64(0))
}

func TestAb05Ap05ContextAndOwnership(b516CycleReducerT *b516CycleReducerTesting.T) {
	for b516CycleReducerKey, b516CycleReducerValue := range map[string]any{"primary": 180, "secondary": 21, "operation": 13, "direction": "other", "view": "other"} {
		for _, b516CycleReducerMode := range []string{"different", "absent", "null"} {
			b516CycleReducerInputValue := b516CycleReducerInput()
			b516CycleReducerContext := b516CycleReducerInputValue["context"].(map[string]any)
			b516CycleReducerContext[b516CycleReducerKey] = b516CycleReducerValue
			if b516CycleReducerMode == "absent" {
				delete(b516CycleReducerContext, b516CycleReducerKey)
			}
			if b516CycleReducerMode == "null" {
				b516CycleReducerContext[b516CycleReducerKey] = nil
			}
			b516CycleReducerOutput := b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["state"], "Unknown")
			if _, b516CycleReducerPresent := b516CycleReducerFields(b516CycleReducerOutput)["sum"]; b516CycleReducerPresent {
				b516CycleReducerT.Fatal("formula on unsupported context")
			}
			b516CycleReducerEncoded, _ := b516CycleReducerJSON.Marshal(b516CycleReducerContext)
			var b516CycleReducerExpected any
			_ = b516CycleReducerJSON.Unmarshal(b516CycleReducerEncoded, &b516CycleReducerExpected)
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["context"], b516CycleReducerExpected)
		}
	}
	for _, b516CycleReducerKey := range []string{"family", "variant", "extra"} {
		for _, b516CycleReducerValue := range []any{nil, "", 0, false, []any{map[string]any{"nested": 1}}} {
			b516CycleReducerInputValue := b516CycleReducerInput()
			b516CycleReducerInputValue["context"].(map[string]any)[b516CycleReducerKey] = b516CycleReducerValue
			b516CycleReducerOutput := b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["state"], "complete")
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["qualified"], false)
		}
	}
	b516CycleReducerInputValue := b516CycleReducerInput()
	b516CycleReducerArgs(b516CycleReducerInputValue)["opaque"] = map[string]any{"nested": []any{false, nil, 0}}
	b516CycleReducerEncoded, _ := b516CycleReducerJSON.Marshal(b516CycleReducerInputValue)
	b516CycleReducerOutput, b516CycleReducerErr := B516CycleReducerEvaluate(b516CycleReducerEncoded)
	if b516CycleReducerErr != nil {
		b516CycleReducerT.Fatal(b516CycleReducerErr)
	}
	b516CycleReducerSaved := append([]byte(nil), b516CycleReducerOutput...)
	for b516CycleReducerIndex := range b516CycleReducerEncoded {
		b516CycleReducerEncoded[b516CycleReducerIndex] = ' '
	}
	b516CycleReducerEqual(b516CycleReducerT, []byte(b516CycleReducerOutput), b516CycleReducerSaved)
	for b516CycleReducerIndex := range b516CycleReducerOutput {
		b516CycleReducerOutput[b516CycleReducerIndex] = ' '
	}
	b516CycleReducerAgain := b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
	b516CycleReducerAgain["raw"].(map[string]any)["opaque"].([]any)[0] = 99
	b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerFields(b516CycleReducerAgain)["opaque"].(map[string]any)["raw"].([]any)[0], float64(0))
	b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)["raw"].(map[string]any)["opaque"].([]any)[0], float64(0))
}

func TestAb05Ap05Validation(b516CycleReducerT *b516CycleReducerTesting.T) {
	for _, b516CycleReducerText := range []string{``, `{`, `[]`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"x":[{"a":1,"a":2}]},"fields":{}}`} {
		b516CycleReducerResult, b516CycleReducerErr := B516CycleReducerEvaluate([]byte(b516CycleReducerText))
		if b516CycleReducerErr == nil || b516CycleReducerResult != nil {
			b516CycleReducerT.Fatalf("accepted %s", b516CycleReducerText)
		}
	}
	for _, b516CycleReducerKey := range []string{"registered", "capable", "prior_changed", "prior", "prior_previous", "outcome", "value"} {
		for _, b516CycleReducerBad := range []any{nil, []any{}, map[string]any{}, "bad", -1, 1.5, uint64(4294967296)} {
			b516CycleReducerInputValue := b516CycleReducerInput()
			b516CycleReducerSlot(b516CycleReducerInputValue, 0)[b516CycleReducerKey] = b516CycleReducerBad
			b516CycleReducerReject(b516CycleReducerT, b516CycleReducerInputValue)
		}
		b516CycleReducerInputValue := b516CycleReducerInput()
		delete(b516CycleReducerSlot(b516CycleReducerInputValue, 0), b516CycleReducerKey)
		b516CycleReducerReject(b516CycleReducerT, b516CycleReducerInputValue)
	}
	for _, b516CycleReducerKey := range []string{"failure_counter", "counter_increment"} {
		for _, b516CycleReducerBad := range []any{nil, false, "0", []any{}, -1, 0.5, uint64(4294967296)} {
			b516CycleReducerInputValue := b516CycleReducerInput()
			b516CycleReducerArgs(b516CycleReducerInputValue)[b516CycleReducerKey] = b516CycleReducerBad
			b516CycleReducerReject(b516CycleReducerT, b516CycleReducerInputValue)
		}
	}
	for _, b516CycleReducerCount := range []int{0, 8, 10} {
		b516CycleReducerInputValue := b516CycleReducerInput()
		b516CycleReducerSlots := []any{}
		for b516CycleReducerIndex := 0; b516CycleReducerIndex < b516CycleReducerCount; b516CycleReducerIndex++ {
			b516CycleReducerSlots = append(b516CycleReducerSlots, b516CycleReducerSlot(b516CycleReducerInputValue, 0))
		}
		b516CycleReducerArgs(b516CycleReducerInputValue)["slots"] = b516CycleReducerSlots
		b516CycleReducerReject(b516CycleReducerT, b516CycleReducerInputValue)
	}
	for _, b516CycleReducerKey := range []string{"slots", "failure_counter"} {
		b516CycleReducerInputValue := b516CycleReducerInput()
		delete(b516CycleReducerArgs(b516CycleReducerInputValue), b516CycleReducerKey)
		b516CycleReducerReject(b516CycleReducerT, b516CycleReducerInputValue)
	}
	for _, b516CycleReducerBad := range []any{nil, false, 0, "bytes", []any{256}, []any{-1}, []any{true}, []any{0.5}, []any{"0"}} {
		b516CycleReducerInputValue := b516CycleReducerInput()
		b516CycleReducerInputValue["fields"].(map[string]any)["opaque"] = b516CycleReducerBad
		b516CycleReducerReject(b516CycleReducerT, b516CycleReducerInputValue)
	}
	for _, b516CycleReducerValue := range [][]any{{}, {0}, {255, 0, 1}} {
		b516CycleReducerInputValue := b516CycleReducerInput()
		b516CycleReducerInputValue["fields"].(map[string]any)["opaque"] = b516CycleReducerValue
		b516CycleReducerOutput := b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
		b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerFields(b516CycleReducerOutput)["sum"]), float64(0))
		b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerFields(b516CycleReducerOutput)["opaque"].(map[string]any)["state"], "Unknown")
	}
	b516CycleReducerInputValue := b516CycleReducerInput()
	b516CycleReducerEncoded, _ := b516CycleReducerJSON.Marshal(b516CycleReducerInputValue)
	for _, b516CycleReducerReplacement := range []string{`"prior":true`, `"prior":10.0`, `"prior":1e1`} {
		b516CycleReducerText := b516CycleReducerStrings.Replace(string(b516CycleReducerEncoded), `"prior":10`, b516CycleReducerReplacement, 1)
		if _, b516CycleReducerErr := B516CycleReducerEvaluate([]byte(b516CycleReducerText)); b516CycleReducerErr == nil {
			b516CycleReducerT.Fatal("accepted wrong numeric kind")
		}
	}
}

func TestAb05Ap05Boundaries(b516CycleReducerT *b516CycleReducerTesting.T) {
	for _, b516CycleReducerKey := range []string{"primary", "secondary", "operation"} {
		for _, b516CycleReducerBad := range []any{false, "181", -1, 1.5, 256, []any{}} {
			b516CycleReducerInputValue := b516CycleReducerInput()
			b516CycleReducerInputValue["context"].(map[string]any)[b516CycleReducerKey] = b516CycleReducerBad
			b516CycleReducerReject(b516CycleReducerT, b516CycleReducerInputValue)
		}
	}
	for _, b516CycleReducerKey := range []string{"direction", "view"} {
		for _, b516CycleReducerBad := range []any{false, 0, []any{}, map[string]any{}} {
			b516CycleReducerInputValue := b516CycleReducerInput()
			b516CycleReducerInputValue["context"].(map[string]any)[b516CycleReducerKey] = b516CycleReducerBad
			b516CycleReducerReject(b516CycleReducerT, b516CycleReducerInputValue)
		}
	}
	for _, b516CycleReducerKey := range []string{"prior", "prior_previous", "value"} {
		for _, b516CycleReducerBad := range []any{false, true} {
			b516CycleReducerInputValue := b516CycleReducerInput()
			b516CycleReducerSlot(b516CycleReducerInputValue, 0)[b516CycleReducerKey] = b516CycleReducerBad
			b516CycleReducerReject(b516CycleReducerT, b516CycleReducerInputValue)
		}
	}
	for _, b516CycleReducerKey := range []string{"arguments", "context", "fields"} {
		b516CycleReducerInputValue := b516CycleReducerInput()
		delete(b516CycleReducerInputValue, b516CycleReducerKey)
		b516CycleReducerReject(b516CycleReducerT, b516CycleReducerInputValue)
		for _, b516CycleReducerBad := range []any{nil, false, []any{}, 0, ""} {
			b516CycleReducerInputValue := b516CycleReducerInput()
			b516CycleReducerInputValue[b516CycleReducerKey] = b516CycleReducerBad
			b516CycleReducerReject(b516CycleReducerT, b516CycleReducerInputValue)
		}
	}
	for _, b516CycleReducerBad := range []any{nil, false, map[string]any{}, ""} {
		b516CycleReducerInputValue := b516CycleReducerInput()
		b516CycleReducerArgs(b516CycleReducerInputValue)["slots"] = b516CycleReducerBad
		b516CycleReducerReject(b516CycleReducerT, b516CycleReducerInputValue)
	}
	b516CycleReducerInputValue := b516CycleReducerInput()
	b516CycleReducerSlot(b516CycleReducerInputValue, 0)["outcome"] = "parse-failure"
	b516CycleReducerReject(b516CycleReducerT, b516CycleReducerInputValue)
	for _, b516CycleReducerCounter := range []uint32{0, 1, 4294967295} {
		for _, b516CycleReducerIncrement := range []uint32{0, 1, 4294967295} {
			b516CycleReducerInputValue := b516CycleReducerInput()
			b516CycleReducerArgs(b516CycleReducerInputValue)["failure_counter"] = b516CycleReducerCounter
			b516CycleReducerArgs(b516CycleReducerInputValue)["counter_increment"] = b516CycleReducerIncrement
			b516CycleReducerOutput := b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerFields(b516CycleReducerOutput)["failure_counter"]), float64(0))
		}
	}
	for _, b516CycleReducerKey := range []string{"family", "variant", "opaque"} {
		for _, b516CycleReducerValue := range []any{nil, false, "", "other", 0, []any{map[string]any{"nested": false}}} {
			b516CycleReducerInputValue := b516CycleReducerInput()
			b516CycleReducerInputValue["context"].(map[string]any)[b516CycleReducerKey] = b516CycleReducerValue
			b516CycleReducerArgs(b516CycleReducerInputValue)["extra"] = b516CycleReducerValue
			b516CycleReducerOutput := b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
			b516CycleReducerEncoded, _ := b516CycleReducerJSON.Marshal(b516CycleReducerInputValue)
			var b516CycleReducerExpected map[string]any
			_ = b516CycleReducerJSON.Unmarshal(b516CycleReducerEncoded, &b516CycleReducerExpected)
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["context"], b516CycleReducerExpected["context"])
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["arguments"], b516CycleReducerExpected["arguments"])
		}
	}
	b516CycleReducerInputValue = b516CycleReducerInput()
	delete(b516CycleReducerInputValue["fields"].(map[string]any), "opaque")
	b516CycleReducerOutput := b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
	b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["raw"], map[string]any{})
	b516CycleReducerInputValue["fields"].(map[string]any)["other"] = []any{0}
	b516CycleReducerOutput = b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
	b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["raw"].(map[string]any)["other"], []any{float64(0)})
}

func TestAb05Ap05CounterAndPreviousBoundaries(b516CycleReducerT *b516CycleReducerTesting.T) {
	for _, b516CycleReducerCase := range []struct {
		b516CycleReducerPrior     uint32
		b516CycleReducerIncrement uint32
		b516CycleReducerExpected  float64
	}{
		{0, 0, 0}, {0, 4294967295, 4294967295}, {4294967295, 0, 4294967295}, {4294967295, 1, 0}, {4294967295, 4294967295, 4294967294},
	} {
		b516CycleReducerInputValue := b516CycleReducerInput()
		b516CycleReducerArgs(b516CycleReducerInputValue)["failure_counter"] = b516CycleReducerCase.b516CycleReducerPrior
		b516CycleReducerArgs(b516CycleReducerInputValue)["counter_increment"] = b516CycleReducerCase.b516CycleReducerIncrement
		b516CycleReducerSlot(b516CycleReducerInputValue, 0)["outcome"] = "parse-failure"
		delete(b516CycleReducerSlot(b516CycleReducerInputValue, 0), "value")
		b516CycleReducerOutput := b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
		b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerFields(b516CycleReducerOutput)["failure_counter"]), b516CycleReducerCase.b516CycleReducerExpected)
		b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerFields(b516CycleReducerOutput)["failure_counter"].(map[string]any)["native_increments"], "Unknown")
	}
	for _, b516CycleReducerPrevious := range []uint32{0, 4294967295} {
		for _, b516CycleReducerPrior := range []uint32{0, 4294967295} {
			b516CycleReducerInputValue := b516CycleReducerInput()
			b516CycleReducerSlot(b516CycleReducerInputValue, 0)["prior_previous"] = b516CycleReducerPrevious
			b516CycleReducerSlot(b516CycleReducerInputValue, 0)["prior"] = b516CycleReducerPrior
			b516CycleReducerSlot(b516CycleReducerInputValue, 0)["registered"] = false
			b516CycleReducerOutput := b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerOutcome(b516CycleReducerOutput, 0)["previous"]), float64(b516CycleReducerPrevious))
			b516CycleReducerSlot(b516CycleReducerInputValue, 0)["registered"] = true
			b516CycleReducerSlot(b516CycleReducerInputValue, 0)["outcome"] = "read-failure"
			delete(b516CycleReducerSlot(b516CycleReducerInputValue, 0), "value")
			b516CycleReducerOutput = b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
			b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerNumber(b516CycleReducerOutcome(b516CycleReducerOutput, 0)["previous"]), float64(b516CycleReducerPrior))
		}
	}
	for _, b516CycleReducerKey := range []string{"primary", "secondary", "operation"} {
		b516CycleReducerInputValue := b516CycleReducerInput()
		b516CycleReducerInputValue["context"].(map[string]any)[b516CycleReducerKey] = 0
		b516CycleReducerOutput := b516CycleReducerRun(b516CycleReducerT, b516CycleReducerInputValue)
		b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["state"], "Unknown")
		b516CycleReducerEqual(b516CycleReducerT, b516CycleReducerOutput["context"].(map[string]any)[b516CycleReducerKey], float64(0))
	}
}
