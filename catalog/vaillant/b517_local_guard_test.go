package vaillant

import (
	b517LocalGuardBytes "bytes"
	b517LocalGuardJSON "encoding/json"
	b517LocalGuardOs "os"
	b517LocalGuardTesting "testing"
)

type b517LocalGuardCase struct {
	B517LocalGuardIncomplete []string                      `json:"incomplete"`
	B517LocalGuardID         string                        `json:"id"`
	B517LocalGuardInput      b517LocalGuardJSON.RawMessage `json:"input"`
	B517LocalGuardState      string                        `json:"state"`
	B517LocalGuardEligible   bool                          `json:"eligible"`
	B517LocalGuardFalse      []string                      `json:"false"`
}

func b517LocalGuardCases(b517LocalGuardT *b517LocalGuardTesting.T) []b517LocalGuardCase {
	b517LocalGuardT.Helper()
	b517LocalGuardData, b517LocalGuardErr := b517LocalGuardOs.ReadFile("testdata/ab11aq11_b517_local_guard.json")
	if b517LocalGuardErr != nil {
		b517LocalGuardT.Fatal(b517LocalGuardErr)
	}
	var b517LocalGuardCases []b517LocalGuardCase
	if b517LocalGuardErr := b517LocalGuardJSON.Unmarshal(b517LocalGuardData, &b517LocalGuardCases); b517LocalGuardErr != nil {
		b517LocalGuardT.Fatal(b517LocalGuardErr)
	}
	if len(b517LocalGuardCases) != 22 {
		b517LocalGuardT.Fatalf("fixture count: %d", len(b517LocalGuardCases))
	}
	return b517LocalGuardCases
}

func b517LocalGuardObject(b517LocalGuardT *b517LocalGuardTesting.T, b517LocalGuardData []byte) map[string]any {
	b517LocalGuardT.Helper()
	var b517LocalGuardObject map[string]any
	if b517LocalGuardErr := b517LocalGuardJSON.Unmarshal(b517LocalGuardData, &b517LocalGuardObject); b517LocalGuardErr != nil {
		b517LocalGuardT.Fatal(b517LocalGuardErr)
	}
	return b517LocalGuardObject
}

func b517LocalGuardEvaluate(b517LocalGuardT *b517LocalGuardTesting.T, b517LocalGuardInput any) map[string]any {
	b517LocalGuardT.Helper()
	b517LocalGuardData, b517LocalGuardErr := b517LocalGuardJSON.Marshal(b517LocalGuardInput)
	if b517LocalGuardErr != nil {
		b517LocalGuardT.Fatal(b517LocalGuardErr)
	}
	b517LocalGuardOutput, b517LocalGuardErr := B517LocalGuardEvaluate(b517LocalGuardData)
	if b517LocalGuardErr != nil {
		b517LocalGuardT.Fatal(b517LocalGuardErr)
	}
	return b517LocalGuardObject(b517LocalGuardT, b517LocalGuardOutput)
}

func TestAb11Aq11Fixture(b517LocalGuardT *b517LocalGuardTesting.T) {
	for _, b517LocalGuardCase := range b517LocalGuardCases(b517LocalGuardT) {
		b517LocalGuardT.Run(b517LocalGuardCase.B517LocalGuardID, func(b517LocalGuardT *b517LocalGuardTesting.T) {
			b517LocalGuardInput := b517LocalGuardObject(b517LocalGuardT, b517LocalGuardCase.B517LocalGuardInput)
			b517LocalGuardOutput := b517LocalGuardEvaluate(b517LocalGuardT, b517LocalGuardInput)
			if b517LocalGuardOutput["state"] != b517LocalGuardCase.B517LocalGuardState || b517LocalGuardOutput["physical"] != "Unknown" || b517LocalGuardOutput["qualified"] != false {
				b517LocalGuardT.Fatalf("envelope: %v", b517LocalGuardOutput)
			}
			for b517LocalGuardKey, b517LocalGuardValue := range b517LocalGuardInput {
				b517LocalGuardBefore, _ := b517LocalGuardJSON.Marshal(b517LocalGuardValue)
				b517LocalGuardAfter, _ := b517LocalGuardJSON.Marshal(b517LocalGuardOutput[b517LocalGuardKey])
				if !b517LocalGuardBytes.Equal(b517LocalGuardBefore, b517LocalGuardAfter) {
					b517LocalGuardT.Fatalf("raw retention %s", b517LocalGuardKey)
				}
			}
			b517LocalGuardResults := b517LocalGuardOutput["results"].(map[string]any)
			if len(b517LocalGuardResults) != 8 {
				b517LocalGuardT.Fatalf("results: %v", b517LocalGuardResults)
			}
			for b517LocalGuardName, b517LocalGuardValue := range b517LocalGuardResults {
				b517LocalGuardResult := b517LocalGuardValue.(map[string]any)
				b517LocalGuardExpectedState := "complete"
				for _, b517LocalGuardIncomplete := range b517LocalGuardCase.B517LocalGuardIncomplete {
					if b517LocalGuardName == b517LocalGuardIncomplete {
						b517LocalGuardExpectedState = "incomplete"
					}
				}
				if b517LocalGuardResult["state"] != b517LocalGuardExpectedState {
					b517LocalGuardT.Errorf("state %s: %v", b517LocalGuardName, b517LocalGuardResult)
				}
				b517LocalGuardExpected := true
				if b517LocalGuardName == "eligible" {
					b517LocalGuardExpected = b517LocalGuardCase.B517LocalGuardEligible
				}
				for _, b517LocalGuardFalse := range b517LocalGuardCase.B517LocalGuardFalse {
					if b517LocalGuardName == b517LocalGuardFalse {
						b517LocalGuardExpected = false
					}
				}
				if b517LocalGuardResult["flags"].(map[string]any)["value"] != b517LocalGuardExpected {
					b517LocalGuardT.Errorf("predicate %s: %v", b517LocalGuardName, b517LocalGuardResult)
				}
				if _, b517LocalGuardOK := b517LocalGuardResult["raw"]; !b517LocalGuardOK {
					b517LocalGuardT.Error("missing raw")
				}
				if b517LocalGuardResult["grade"] != "Hypothesis" || (b517LocalGuardResult["state"] == "complete" && b517LocalGuardResult["calculation"] != "Proven") || (b517LocalGuardResult["state"] == "incomplete" && b517LocalGuardResult["calculation"] != nil) {
					b517LocalGuardT.Error("grades")
				}
			}
			if b517LocalGuardCase.B517LocalGuardState == "partial" {
				if b517LocalGuardResults["eligible"].(map[string]any)["state"] != "incomplete" || b517LocalGuardResults["ready_state"].(map[string]any)["state"] != "complete" {
					b517LocalGuardT.Error("partial sibling state")
				}
			}
		})
	}
}

func TestAb11Aq11Context(b517LocalGuardT *b517LocalGuardTesting.T) {
	for b517LocalGuardName, b517LocalGuardValue := range map[string]any{"primary": 180, "secondary": 24, "operation": 0, "direction": "remote", "view": "B518"} {
		b517LocalGuardT.Run(b517LocalGuardName, func(b517LocalGuardT *b517LocalGuardTesting.T) {
			b517LocalGuardInput := b517LocalGuardObject(b517LocalGuardT, b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput)
			b517LocalGuardInput["context"].(map[string]any)[b517LocalGuardName] = b517LocalGuardValue
			b517LocalGuardOutput := b517LocalGuardEvaluate(b517LocalGuardT, b517LocalGuardInput)
			if b517LocalGuardOutput["state"] != "Unknown" || len(b517LocalGuardOutput["results"].(map[string]any)) != 0 || b517LocalGuardOutput["qualified"] != false {
				b517LocalGuardT.Fatalf("unsupported context: %v", b517LocalGuardOutput)
			}
			b517LocalGuardBefore, _ := b517LocalGuardJSON.Marshal(b517LocalGuardInput["context"])
			b517LocalGuardAfter, _ := b517LocalGuardJSON.Marshal(b517LocalGuardOutput["context"])
			if !b517LocalGuardBytes.Equal(b517LocalGuardBefore, b517LocalGuardAfter) {
				b517LocalGuardT.Error("context retention")
			}
		})
	}
	for _, b517LocalGuardName := range []string{"family", "variant"} {
		b517LocalGuardInput := b517LocalGuardObject(b517LocalGuardT, b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput)
		delete(b517LocalGuardInput["context"].(map[string]any), b517LocalGuardName)
		b517LocalGuardOutput := b517LocalGuardEvaluate(b517LocalGuardT, b517LocalGuardInput)
		if _, b517LocalGuardPresent := b517LocalGuardOutput["context"].(map[string]any)[b517LocalGuardName]; b517LocalGuardPresent {
			b517LocalGuardT.Error("opaque absence lost")
		}
		if b517LocalGuardOutput["state"] != "complete" || b517LocalGuardOutput["qualified"] != false {
			b517LocalGuardT.Error("opaque absence qualified")
		}
		for _, b517LocalGuardValue := range []any{nil, false, 0, "B518", []any{}} {
			b517LocalGuardInput := b517LocalGuardObject(b517LocalGuardT, b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput)
			b517LocalGuardInput["context"].(map[string]any)[b517LocalGuardName] = b517LocalGuardValue
			b517LocalGuardOutput := b517LocalGuardEvaluate(b517LocalGuardT, b517LocalGuardInput)
			if b517LocalGuardOutput["qualified"] != false || b517LocalGuardOutput["physical"] != "Unknown" || b517LocalGuardOutput["state"] != "complete" {
				b517LocalGuardT.Error("opaque identity")
			}
		}
	}
}

func TestAb11Aq11Reject(b517LocalGuardT *b517LocalGuardTesting.T) {
	for _, b517LocalGuardData := range []string{`{`, `[]`, `null`, `{} {}`, `{"x":1,"x":2}`, `{"x":{"a":1,"\u0061":2}}`, `{"x":[{"a":1,"a":2}]}`, `{"x":1e999}`, `{"x":NaN}`} {
		b517LocalGuardOutput, b517LocalGuardErr := B517LocalGuardEvaluate([]byte(b517LocalGuardData))
		if b517LocalGuardErr == nil || b517LocalGuardOutput != nil {
			b517LocalGuardT.Errorf("accepted %s", b517LocalGuardData)
		}
	}
	b517LocalGuardValid := b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput
	b517LocalGuardEmbedded := []string{`"x":1,"x":2`, `"x":{"a":1,"\u0061":2}`, `"x":[{"a":1,"a":2}]`, `"x":"\ud800"`, `"x":"\udc00"`, `"x":"\ud800\u0041"`, "\"x\":\"\xff\"", "\"x\":\"\xc3\""}
	if _, b517LocalGuardErr := B517LocalGuardEvaluate(b517LocalGuardValid); b517LocalGuardErr != nil {
		b517LocalGuardT.Fatal(b517LocalGuardErr)
	}
	for _, b517LocalGuardMember := range b517LocalGuardEmbedded {
		b517LocalGuardData := append([]byte(`{`+b517LocalGuardMember+`,`), b517LocalGuardBytes.TrimLeft(b517LocalGuardValid, " \t\r\n")[1:]...)
		if b517LocalGuardOutput, b517LocalGuardErr := B517LocalGuardEvaluate(b517LocalGuardData); b517LocalGuardErr == nil || b517LocalGuardOutput != nil {
			b517LocalGuardT.Errorf("accepted embedded %q", b517LocalGuardMember)
		}
	}
	for _, b517LocalGuardGroup := range []string{"context", "fields", "arguments"} {
		b517LocalGuardNames := map[string][]string{"context": {"primary", "secondary", "operation"}, "fields": {"request", "input"}, "arguments": {"ready_state", "final_option", "received_length"}}[b517LocalGuardGroup]
		for _, b517LocalGuardName := range b517LocalGuardNames {
			for _, b517LocalGuardBad := range []any{nil, true, "3", -1, 1.5, 65536, 1.0000000001} {
				b517LocalGuardInput := b517LocalGuardObject(b517LocalGuardT, b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput)
				if b517LocalGuardGroup == "fields" {
					b517LocalGuardInput[b517LocalGuardGroup].(map[string]any)[b517LocalGuardName] = []any{b517LocalGuardBad}
				} else {
					b517LocalGuardInput[b517LocalGuardGroup].(map[string]any)[b517LocalGuardName] = b517LocalGuardBad
				}
				b517LocalGuardData, _ := b517LocalGuardJSON.Marshal(b517LocalGuardInput)
				b517LocalGuardOutput, b517LocalGuardErr := B517LocalGuardEvaluate(b517LocalGuardData)
				if b517LocalGuardErr == nil || b517LocalGuardOutput != nil {
					b517LocalGuardT.Errorf("accepted %s.%s=%v", b517LocalGuardGroup, b517LocalGuardName, b517LocalGuardBad)
				}
			}
		}
	}
	for _, b517LocalGuardName := range []string{"ready_state", "final_option", "received_length"} {
		b517LocalGuardInput := b517LocalGuardObject(b517LocalGuardT, b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput)
		delete(b517LocalGuardInput["arguments"].(map[string]any), b517LocalGuardName)
		b517LocalGuardData, _ := b517LocalGuardJSON.Marshal(b517LocalGuardInput)
		if b517LocalGuardOutput, b517LocalGuardErr := B517LocalGuardEvaluate(b517LocalGuardData); b517LocalGuardErr == nil || b517LocalGuardOutput != nil {
			b517LocalGuardT.Errorf("missing argument %s", b517LocalGuardName)
		}
	}
}

func TestAb11Aq11Memory(b517LocalGuardT *b517LocalGuardTesting.T) {
	b517LocalGuardInput := append([]byte(nil), b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput...)
	b517LocalGuardOriginal := append([]byte(nil), b517LocalGuardInput...)
	b517LocalGuardFirst, b517LocalGuardErr := B517LocalGuardEvaluate(b517LocalGuardInput)
	if b517LocalGuardErr != nil {
		b517LocalGuardT.Fatal(b517LocalGuardErr)
	}
	b517LocalGuardSaved := append([]byte(nil), b517LocalGuardFirst...)
	for b517LocalGuardIndex := range b517LocalGuardInput {
		b517LocalGuardInput[b517LocalGuardIndex] = 'x'
	}
	if !b517LocalGuardBytes.Equal(b517LocalGuardFirst, b517LocalGuardSaved) {
		b517LocalGuardT.Fatal("input aliases output")
	}
	for b517LocalGuardIndex := range b517LocalGuardFirst {
		b517LocalGuardFirst[b517LocalGuardIndex] = 'x'
	}
	b517LocalGuardSecond, b517LocalGuardErr := B517LocalGuardEvaluate(b517LocalGuardOriginal)
	if b517LocalGuardErr != nil || !b517LocalGuardBytes.Equal(b517LocalGuardSecond, b517LocalGuardSaved) {
		b517LocalGuardT.Fatal("cross-call storage")
	}
}

func TestAb11Aq11KindsAndBounds(b517LocalGuardT *b517LocalGuardTesting.T) {
	for _, b517LocalGuardGroup := range []string{"context", "fields", "arguments"} {
		for _, b517LocalGuardBad := range []any{nil, false, 0, "", []any{}} {
			b517LocalGuardInput := b517LocalGuardObject(b517LocalGuardT, b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput)
			b517LocalGuardInput[b517LocalGuardGroup] = b517LocalGuardBad
			b517LocalGuardData, _ := b517LocalGuardJSON.Marshal(b517LocalGuardInput)
			if b517LocalGuardOutput, b517LocalGuardErr := B517LocalGuardEvaluate(b517LocalGuardData); b517LocalGuardErr == nil || b517LocalGuardOutput != nil {
				b517LocalGuardT.Errorf("accepted group %s: %v", b517LocalGuardGroup, b517LocalGuardBad)
			}
		}
	}
	for _, b517LocalGuardName := range []string{"direction", "view"} {
		b517LocalGuardInput := b517LocalGuardObject(b517LocalGuardT, b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput)
		b517LocalGuardInput["context"].(map[string]any)[b517LocalGuardName] = nil
		b517LocalGuardData, _ := b517LocalGuardJSON.Marshal(b517LocalGuardInput)
		if b517LocalGuardOutput, b517LocalGuardErr := B517LocalGuardEvaluate(b517LocalGuardData); b517LocalGuardErr == nil || b517LocalGuardOutput != nil {
			b517LocalGuardT.Errorf("accepted context kind %s", b517LocalGuardName)
		}
	}
	for _, b517LocalGuardName := range []string{"primary", "secondary", "operation"} {
		b517LocalGuardInput := b517LocalGuardObject(b517LocalGuardT, b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput)
		b517LocalGuardInput["context"].(map[string]any)[b517LocalGuardName] = 256
		b517LocalGuardData, _ := b517LocalGuardJSON.Marshal(b517LocalGuardInput)
		if b517LocalGuardOutput, b517LocalGuardErr := B517LocalGuardEvaluate(b517LocalGuardData); b517LocalGuardErr == nil || b517LocalGuardOutput != nil {
			b517LocalGuardT.Errorf("accepted context overflow %s", b517LocalGuardName)
		}
	}
	for _, b517LocalGuardName := range []string{"ready_state", "final_option"} {
		b517LocalGuardInput := b517LocalGuardObject(b517LocalGuardT, b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput)
		b517LocalGuardInput["arguments"].(map[string]any)[b517LocalGuardName] = 256
		b517LocalGuardData, _ := b517LocalGuardJSON.Marshal(b517LocalGuardInput)
		if b517LocalGuardOutput, b517LocalGuardErr := B517LocalGuardEvaluate(b517LocalGuardData); b517LocalGuardErr == nil || b517LocalGuardOutput != nil {
			b517LocalGuardT.Errorf("accepted u8 overflow %s", b517LocalGuardName)
		}
	}
	for _, b517LocalGuardName := range []string{"request", "input"} {
		for _, b517LocalGuardBad := range []any{nil, false, "", 0, []any{256}} {
			b517LocalGuardInput := b517LocalGuardObject(b517LocalGuardT, b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput)
			b517LocalGuardInput["fields"].(map[string]any)[b517LocalGuardName] = b517LocalGuardBad
			b517LocalGuardData, _ := b517LocalGuardJSON.Marshal(b517LocalGuardInput)
			if b517LocalGuardOutput, b517LocalGuardErr := B517LocalGuardEvaluate(b517LocalGuardData); b517LocalGuardErr == nil || b517LocalGuardOutput != nil {
				b517LocalGuardT.Errorf("accepted field kind/width %s", b517LocalGuardName)
			}
		}
	}
	for _, b517LocalGuardValue := range []int{0, 65535} {
		b517LocalGuardInput := b517LocalGuardObject(b517LocalGuardT, b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput)
		b517LocalGuardInput["arguments"].(map[string]any)["received_length"] = b517LocalGuardValue
		b517LocalGuardOutput := b517LocalGuardEvaluate(b517LocalGuardT, b517LocalGuardInput)
		if b517LocalGuardOutput["results"].(map[string]any)["eligible"].(map[string]any)["flags"].(map[string]any)["value"] != false {
			b517LocalGuardT.Error("bounded length accepted as eligible")
		}
	}
	b517LocalGuardInput := b517LocalGuardObject(b517LocalGuardT, b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput)
	b517LocalGuardInput["arguments"].(map[string]any)["opaque"] = map[string]any{"null": nil, "zero": 0, "false": false, "empty": []any{}}
	b517LocalGuardOutput := b517LocalGuardEvaluate(b517LocalGuardT, b517LocalGuardInput)
	b517LocalGuardBefore, _ := b517LocalGuardJSON.Marshal(b517LocalGuardInput["arguments"])
	b517LocalGuardAfter, _ := b517LocalGuardJSON.Marshal(b517LocalGuardOutput["arguments"])
	if !b517LocalGuardBytes.Equal(b517LocalGuardBefore, b517LocalGuardAfter) {
		b517LocalGuardT.Error("argument extras retention")
	}
	for _, b517LocalGuardKey := range []string{"state", "results", "physical", "qualified", "raw"} {
		b517LocalGuardInput := b517LocalGuardObject(b517LocalGuardT, b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput)
		b517LocalGuardInput[b517LocalGuardKey] = map[string]any{"opaque": nil}
		b517LocalGuardOutput := b517LocalGuardEvaluate(b517LocalGuardT, b517LocalGuardInput)
		b517LocalGuardBefore, _ := b517LocalGuardJSON.Marshal(b517LocalGuardInput)
		b517LocalGuardAfter, _ := b517LocalGuardJSON.Marshal(b517LocalGuardOutput["raw"])
		if !b517LocalGuardBytes.Equal(b517LocalGuardBefore, b517LocalGuardAfter) {
			b517LocalGuardT.Errorf("reserved output key loses raw input: %s", b517LocalGuardKey)
		}
	}
	for _, b517LocalGuardSuffix := range []string{" {}", " true", " null", " []"} {
		b517LocalGuardData := append(append([]byte(nil), b517LocalGuardCases(b517LocalGuardT)[0].B517LocalGuardInput...), []byte(b517LocalGuardSuffix)...)
		if b517LocalGuardOutput, b517LocalGuardErr := B517LocalGuardEvaluate(b517LocalGuardData); b517LocalGuardErr == nil || b517LocalGuardOutput != nil {
			b517LocalGuardT.Error("accepted trailing JSON")
		}
	}
}
