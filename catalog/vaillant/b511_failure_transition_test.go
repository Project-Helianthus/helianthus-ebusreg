package vaillant

import (
	"encoding/json"
	b511FailureTransitionOs "os"
	b511FailureTransitionReflect "reflect"
	b511FailureTransitionStrings "strings"
	b511FailureTransitionTesting "testing"
)

func b511FailureTransitionBase() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 17, "operation": 1, "direction": "local", "view": "refresh", "family": map[string]any{"opaque": "synthetic-family"}, "variant": []any{"synthetic-variant"}}, "fields": map[string]any{"registration": []int{1}, "prior_failure": []int{1}, "local_body": []int{10, 20, 30, 40, 50, 60, 70, 80, 0}, "unmodelled": []int{255, 0}}, "arguments": map[string]any{"index": 0, "selector": 2, "local_selector": 3, "count": 1, "parsed": true}}
}

func b511FailureTransitionCall(b511FailureTransitionT *b511FailureTransitionTesting.T, b511FailureTransitionInput map[string]any) map[string]any {
	b511FailureTransitionT.Helper()
	b511FailureTransitionEncoded, b511FailureTransitionErr := json.Marshal(b511FailureTransitionInput)
	if b511FailureTransitionErr != nil {
		b511FailureTransitionT.Fatal(b511FailureTransitionErr)
	}
	b511FailureTransitionResult, b511FailureTransitionErr := B511FailureTransitionInterpret(b511FailureTransitionEncoded)
	if b511FailureTransitionErr != nil {
		b511FailureTransitionT.Fatal(b511FailureTransitionErr)
	}
	b511FailureTransitionOutput := map[string]any{}
	if b511FailureTransitionErr := json.Unmarshal(b511FailureTransitionResult, &b511FailureTransitionOutput); b511FailureTransitionErr != nil {
		b511FailureTransitionT.Fatal(b511FailureTransitionErr)
	}
	if b511FailureTransitionOutput["physical"] != "Unknown" || b511FailureTransitionOutput["qualified"] != false {
		b511FailureTransitionT.Fatal("native qualification inferred")
	}
	b511FailureTransitionRoundtrip, _ := json.Marshal(b511FailureTransitionOutput)
	b511FailureTransitionAgain := map[string]any{}
	_ = json.Unmarshal(b511FailureTransitionRoundtrip, &b511FailureTransitionAgain)
	if !b511FailureTransitionReflect.DeepEqual(b511FailureTransitionOutput, b511FailureTransitionAgain) {
		b511FailureTransitionT.Fatal("roundtrip")
	}
	b511FailureTransitionInputRoundtrip := map[string]any{}
	_ = json.Unmarshal(b511FailureTransitionEncoded, &b511FailureTransitionInputRoundtrip)
	for _, b511FailureTransitionPair := range [][2]string{{"context", "context"}, {"fields", "raw"}, {"arguments", "arguments"}} {
		if !b511FailureTransitionReflect.DeepEqual(b511FailureTransitionInputRoundtrip[b511FailureTransitionPair[0]], b511FailureTransitionOutput[b511FailureTransitionPair[1]]) {
			b511FailureTransitionT.Fatalf("preservation %s", b511FailureTransitionPair[0])
		}
	}
	return b511FailureTransitionOutput
}

func b511FailureTransitionCheckFailure(b511FailureTransitionT *b511FailureTransitionTesting.T, b511FailureTransitionOutput map[string]any, b511FailureTransitionWant bool, b511FailureTransitionState string, b511FailureTransitionCopied bool) {
	b511FailureTransitionT.Helper()
	if b511FailureTransitionOutput["state"] != b511FailureTransitionState {
		b511FailureTransitionT.Fatalf("state: %v", b511FailureTransitionOutput)
	}
	b511FailureTransitionFields := b511FailureTransitionOutput["fields"].(map[string]any)
	b511FailureTransitionFailure := b511FailureTransitionFields["failure"].(map[string]any)
	if b511FailureTransitionFailure["state"] != "Hypothesis" || b511FailureTransitionFailure["flags"].(map[string]any)["failure"] != b511FailureTransitionWant {
		b511FailureTransitionT.Fatalf("failure: %v", b511FailureTransitionFailure)
	}
	if _, b511FailureTransitionExists := b511FailureTransitionFailure["number"]; b511FailureTransitionExists {
		b511FailureTransitionT.Fatal("invented failure number")
	}
	b511FailureTransitionCopy, b511FailureTransitionExists := b511FailureTransitionFields["copied"]
	if b511FailureTransitionExists != b511FailureTransitionCopied {
		b511FailureTransitionT.Fatalf("copy presence: %v", b511FailureTransitionFields)
	}
	if b511FailureTransitionCopied {
		b511FailureTransitionRaw := b511FailureTransitionOutput["raw"].(map[string]any)["local_body"].([]any)
		if !b511FailureTransitionReflect.DeepEqual(b511FailureTransitionCopy.(map[string]any)["raw"], b511FailureTransitionRaw[:9]) || len(b511FailureTransitionCopy.(map[string]any)["raw"].([]any)) != 9 {
			b511FailureTransitionT.Fatal("copy must be exactly first nine")
		}
	}
}

func TestAb13Ao13Refresh(b511FailureTransitionT *b511FailureTransitionTesting.T) {
	for _, b511FailureTransitionCase := range []struct {
		b511FailureTransitionName  string
		b511FailureTransitionGroup string
		b511FailureTransitionKey   string
		b511FailureTransitionValue any
		b511FailureTransitionWant  bool
		b511FailureTransitionCopy  bool
	}{
		{"ab13ao13_control", "arguments", "index", 0, false, true},
		{"ab13ao13_registration_bit", "fields", "registration", []int{2}, true, false},
		{"ab13ao13_registration_zero", "fields", "registration", []int{0}, true, false},
		{"ab13ao13_registration_high", "fields", "registration", []int{3}, false, true},
		{"ab13ao13_index_exclusion", "arguments", "index", 1, true, false},
		{"ab13ao13_index_two", "arguments", "index", 2, false, true},
		{"ab13ao13_selector_exclusion", "arguments", "selector", 3, true, false},
		{"ab13ao13_local_selector_exclusion", "arguments", "local_selector", 2, true, false},
		{"ab13ao13_selector_other", "arguments", "selector", 4, false, true},
		{"ab13ao13_local_selector_other", "arguments", "local_selector", 4, false, true},
		{"ab13ao13_count_zero", "arguments", "count", 0, true, false},
		{"ab13ao13_count_two", "arguments", "count", 2, false, true},
		{"ab13ao13_parse_false", "arguments", "parsed", false, true, false},
		{"ab13ao13_prior_zero", "fields", "prior_failure", []int{0}, false, true},
	} {
		b511FailureTransitionT.Run(b511FailureTransitionCase.b511FailureTransitionName, func(b511FailureTransitionT *b511FailureTransitionTesting.T) {
			b511FailureTransitionInput := b511FailureTransitionBase()
			b511FailureTransitionInput[b511FailureTransitionCase.b511FailureTransitionGroup].(map[string]any)[b511FailureTransitionCase.b511FailureTransitionKey] = b511FailureTransitionCase.b511FailureTransitionValue
			b511FailureTransitionCheckFailure(b511FailureTransitionT, b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionInput), b511FailureTransitionCase.b511FailureTransitionWant, "Hypothesis", b511FailureTransitionCase.b511FailureTransitionCopy)
		})
	}
	for b511FailureTransitionLast := 0; b511FailureTransitionLast < 3; b511FailureTransitionLast++ {
		b511FailureTransitionInput := b511FailureTransitionBase()
		b511FailureTransitionInput["fields"].(map[string]any)["local_body"].([]int)[8] = b511FailureTransitionLast
		b511FailureTransitionCheckFailure(b511FailureTransitionT, b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionInput), b511FailureTransitionLast == 1, "Hypothesis", true)
	}
	for b511FailureTransitionLength := 0; b511FailureTransitionLength <= 10; b511FailureTransitionLength++ {
		b511FailureTransitionInput := b511FailureTransitionBase()
		b511FailureTransitionBody := make([]int, b511FailureTransitionLength)
		for b511FailureTransitionIndex := range b511FailureTransitionBody {
			b511FailureTransitionBody[b511FailureTransitionIndex] = b511FailureTransitionIndex + 10
		}
		b511FailureTransitionInput["fields"].(map[string]any)["local_body"] = b511FailureTransitionBody
		b511FailureTransitionState := "Hypothesis"
		if b511FailureTransitionLength < 9 {
			b511FailureTransitionState = "partial"
		}
		b511FailureTransitionCheckFailure(b511FailureTransitionT, b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionInput), b511FailureTransitionLength < 9, b511FailureTransitionState, b511FailureTransitionLength >= 9)
	}
}

func TestAb13Ao13LocalViews(b511FailureTransitionT *b511FailureTransitionTesting.T) {
	for _, b511FailureTransitionPrior := range []int{0, 1, 2, 255} {
		b511FailureTransitionInput := b511FailureTransitionBase()
		b511FailureTransitionContext := b511FailureTransitionInput["context"].(map[string]any)
		delete(b511FailureTransitionContext, "operation")
		b511FailureTransitionContext["view"] = "timer-consume"
		b511FailureTransitionInput["fields"] = map[string]any{"prior_failure": []int{b511FailureTransitionPrior}}
		b511FailureTransitionInput["arguments"] = map[string]any{"index": 0}
		b511FailureTransitionOutput := b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionInput)
		b511FailureTransitionCheckFailure(b511FailureTransitionT, b511FailureTransitionOutput, b511FailureTransitionPrior != 0, "Hypothesis", false)
		b511FailureTransitionPath := b511FailureTransitionOutput["fields"].(map[string]any)["path"].(map[string]any)
		if b511FailureTransitionPrior == 0 {
			if b511FailureTransitionPath["label"] != "capability-path" || b511FailureTransitionPath["number"] != nil {
				b511FailureTransitionT.Fatal(b511FailureTransitionPath)
			}
		} else if b511FailureTransitionPath["label"] != "error" || b511FailureTransitionPath["number"] != float64(120) {
			b511FailureTransitionT.Fatal(b511FailureTransitionPath)
		}
		b511FailureTransitionInput["arguments"].(map[string]any)["index"] = 1
		b511FailureTransitionOutput = b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionInput)
		if b511FailureTransitionOutput["state"] != "Unknown" || len(b511FailureTransitionOutput["fields"].(map[string]any)) != 0 {
			b511FailureTransitionT.Fatal(b511FailureTransitionOutput)
		}
	}
	b511FailureTransitionInput := b511FailureTransitionBase()
	delete(b511FailureTransitionInput["context"].(map[string]any), "operation")
	b511FailureTransitionInput["context"].(map[string]any)["view"] = "initialize"
	b511FailureTransitionInput["fields"] = map[string]any{}
	delete(b511FailureTransitionInput, "arguments")
	b511FailureTransitionCheckFailure(b511FailureTransitionT, b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionInput), false, "Hypothesis", false)
}

func TestAb13Ao13IdentityAndPresence(b511FailureTransitionT *b511FailureTransitionTesting.T) {
	for _, b511FailureTransitionKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, b511FailureTransitionMissing := range []bool{false, true} {
			b511FailureTransitionInput := b511FailureTransitionBase()
			b511FailureTransitionContext := b511FailureTransitionInput["context"].(map[string]any)
			if b511FailureTransitionMissing {
				delete(b511FailureTransitionContext, b511FailureTransitionKey)
			} else if b511FailureTransitionKey == "direction" || b511FailureTransitionKey == "view" {
				b511FailureTransitionContext[b511FailureTransitionKey] = "unrecognised"
			} else {
				b511FailureTransitionContext[b511FailureTransitionKey] = 0
			}
			b511FailureTransitionOutput := b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionInput)
			if b511FailureTransitionOutput["state"] != "Unknown" || len(b511FailureTransitionOutput["fields"].(map[string]any)) != 0 {
				b511FailureTransitionT.Fatalf("%s: %v", b511FailureTransitionKey, b511FailureTransitionOutput)
			}
		}
	}
	for _, b511FailureTransitionKey := range []string{"family", "variant"} {
		b511FailureTransitionInput := b511FailureTransitionBase()
		b511FailureTransitionInput["context"].(map[string]any)[b511FailureTransitionKey] = map[string]any{"opaque": []any{0, false, "other-family"}}
		b511FailureTransitionCheckFailure(b511FailureTransitionT, b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionInput), false, "Hypothesis", true)
	}
	for _, b511FailureTransitionKey := range []string{"registration", "prior_failure", "local_body"} {
		for _, b511FailureTransitionWidth := range []int{-1, 0, 2} {
			b511FailureTransitionInput := b511FailureTransitionBase()
			if b511FailureTransitionWidth < 0 {
				delete(b511FailureTransitionInput["fields"].(map[string]any), b511FailureTransitionKey)
			} else {
				b511FailureTransitionInput["fields"].(map[string]any)[b511FailureTransitionKey] = make([]int, b511FailureTransitionWidth)
			}
			b511FailureTransitionOutput := b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionInput)
			if b511FailureTransitionOutput["state"] != "partial" {
				b511FailureTransitionT.Fatal(b511FailureTransitionOutput)
			}
			if b511FailureTransitionKey != "prior_failure" {
				b511FailureTransitionCheckFailure(b511FailureTransitionT, b511FailureTransitionOutput, true, "partial", false)
			} else {
				b511FailureTransitionCheckFailure(b511FailureTransitionT, b511FailureTransitionOutput, false, "partial", true)
			}
		}
	}
	for _, b511FailureTransitionKey := range []string{"index", "selector", "local_selector", "count", "parsed"} {
		b511FailureTransitionInput := b511FailureTransitionBase()
		delete(b511FailureTransitionInput["arguments"].(map[string]any), b511FailureTransitionKey)
		b511FailureTransitionCheckFailure(b511FailureTransitionT, b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionInput), true, "partial", false)
	}
	for _, b511FailureTransitionValue := range []any{[]int{}, []int{0}, []int{4, 5, 6}} {
		b511FailureTransitionInput := b511FailureTransitionBase()
		b511FailureTransitionInput["fields"].(map[string]any)["unmodelled"] = b511FailureTransitionValue
		b511FailureTransitionCheckFailure(b511FailureTransitionT, b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionInput), false, "Hypothesis", true)
	}
}

func TestAb13Ao13Malformed(b511FailureTransitionT *b511FailureTransitionTesting.T) {
	for _, b511FailureTransitionText := range []string{``, `{`, `[]`, `null`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"opaque":{"x":0,"x":1}},"fields":{}}`, `{"context":{},"fields":{"x":[256]}}`, `{"context":{},"fields":{"x":[-1]}}`, `{"context":{},"fields":{"x":[1.0]}}`, `{"context":{},"fields":{"x":[true]}}`, `{"context":{},"fields":{"x":[null]}}`, `{"context":{},"fields":{"x":"AA=="}}`, `{"context":{},"fields":{"x":null}}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":[]}`, `{"context":{},"fields":{},"arguments":null}`, `{"context":{"primary":"181"},"fields":{}}`, `{"context":{"view":1},"fields":{}}`, `{"context":{},"fields":{},"arguments":{"parsed":0}}`, `{"context":{},"fields":{},"arguments":{"index":false}}`, `{"context":{},"fields":{},"arguments":{"count":-1}}`, `{"context":{},"fields":{},"arguments":{"extra":0}}`, `{"context":{},"fields":{},"arguments":{"count":0,"count":1}}`} {
		b511FailureTransitionResult, b511FailureTransitionErr := B511FailureTransitionInterpret([]byte(b511FailureTransitionText))
		if b511FailureTransitionErr == nil || b511FailureTransitionResult != nil {
			b511FailureTransitionT.Fatalf("accepted %s: %s", b511FailureTransitionText, b511FailureTransitionResult)
		}
	}
}

func TestAb13Ao13Isolation(b511FailureTransitionT *b511FailureTransitionTesting.T) {
	b511FailureTransitionInput, _ := json.Marshal(b511FailureTransitionBase())
	b511FailureTransitionOriginal := append([]byte(nil), b511FailureTransitionInput...)
	b511FailureTransitionFirst, b511FailureTransitionErr := B511FailureTransitionInterpret(b511FailureTransitionInput)
	if b511FailureTransitionErr != nil {
		b511FailureTransitionT.Fatal(b511FailureTransitionErr)
	}
	b511FailureTransitionExpected := append([]byte(nil), b511FailureTransitionFirst...)
	for b511FailureTransitionIndex := range b511FailureTransitionInput {
		b511FailureTransitionInput[b511FailureTransitionIndex] = ' '
	}
	if string(b511FailureTransitionFirst) != string(b511FailureTransitionExpected) {
		b511FailureTransitionT.Fatal("input alias")
	}
	b511FailureTransitionDecoded := map[string]any{}
	_ = json.Unmarshal(b511FailureTransitionFirst, &b511FailureTransitionDecoded)
	b511FailureTransitionDecoded["context"].(map[string]any)["family"].(map[string]any)["opaque"] = "changed"
	b511FailureTransitionDecoded["raw"].(map[string]any)["local_body"].([]any)[0] = float64(255)
	b511FailureTransitionDecoded["arguments"].(map[string]any)["index"] = 99
	b511FailureTransitionDecoded["fields"].(map[string]any)["copied"].(map[string]any)["raw"].([]any)[0] = 0
	if b511FailureTransitionDecoded["fields"].(map[string]any)["failure"].(map[string]any)["raw"].([]any)[0] != float64(10) {
		b511FailureTransitionT.Fatal("nested result alias")
	}
	for b511FailureTransitionIndex := range b511FailureTransitionFirst {
		b511FailureTransitionFirst[b511FailureTransitionIndex] = ' '
	}
	b511FailureTransitionSecond, b511FailureTransitionErr := B511FailureTransitionInterpret(b511FailureTransitionOriginal)
	if b511FailureTransitionErr != nil || string(b511FailureTransitionSecond) != string(b511FailureTransitionExpected) {
		b511FailureTransitionT.Fatal("call alias", b511FailureTransitionErr)
	}
}

func TestAb13Ao13Fixture(b511FailureTransitionT *b511FailureTransitionTesting.T) {
	b511FailureTransitionData, b511FailureTransitionErr := b511FailureTransitionOs.ReadFile("testdata/ab13ao13_b511_failure_transition.json")
	if b511FailureTransitionErr != nil {
		b511FailureTransitionT.Fatal(b511FailureTransitionErr)
	}
	var b511FailureTransitionCases []struct {
		B511FailureTransitionID      string         `json:"id"`
		B511FailureTransitionInput   map[string]any `json:"input"`
		B511FailureTransitionFailure bool           `json:"failure"`
	}
	if b511FailureTransitionErr := json.Unmarshal(b511FailureTransitionData, &b511FailureTransitionCases); b511FailureTransitionErr != nil || len(b511FailureTransitionCases) == 0 {
		b511FailureTransitionT.Fatal("fixture", b511FailureTransitionErr)
	}
	for _, b511FailureTransitionCase := range b511FailureTransitionCases {
		if !b511FailureTransitionStrings.HasPrefix(b511FailureTransitionCase.B511FailureTransitionID, "ab13ao13") {
			b511FailureTransitionT.Fatal("fixture id")
		}
		b511FailureTransitionCheckFailure(b511FailureTransitionT, b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionCase.B511FailureTransitionInput), b511FailureTransitionCase.B511FailureTransitionFailure, "Hypothesis", true)
	}
}

func TestAb13Ao13RetainedPrior(b511FailureTransitionT *b511FailureTransitionTesting.T) {
	for _, b511FailureTransitionPrior := range []int{0, 1, 2, 255} {
		for _, b511FailureTransitionSkip := range []string{"parsed", "count", "registration", "index", "selector", "short"} {
			b511FailureTransitionInput := b511FailureTransitionBase()
			b511FailureTransitionInput["fields"].(map[string]any)["prior_failure"] = []int{b511FailureTransitionPrior}
			switch b511FailureTransitionSkip {
			case "parsed":
				b511FailureTransitionInput["arguments"].(map[string]any)["parsed"] = false
			case "count":
				b511FailureTransitionInput["arguments"].(map[string]any)["count"] = 0
			case "registration":
				b511FailureTransitionInput["fields"].(map[string]any)["registration"] = []int{2}
			case "index":
				b511FailureTransitionInput["arguments"].(map[string]any)["index"] = 1
			case "selector":
				b511FailureTransitionInput["arguments"].(map[string]any)["selector"] = 3
			case "short":
				b511FailureTransitionInput["fields"].(map[string]any)["local_body"] = []int{}
			}
			b511FailureTransitionState := "Hypothesis"
			if b511FailureTransitionSkip == "short" {
				b511FailureTransitionState = "partial"
			}
			b511FailureTransitionCheckFailure(b511FailureTransitionT, b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionInput), b511FailureTransitionPrior != 0, b511FailureTransitionState, false)
		}
	}
	for b511FailureTransitionIndex := 0; b511FailureTransitionIndex < 9; b511FailureTransitionIndex++ {
		b511FailureTransitionInput := b511FailureTransitionBase()
		b511FailureTransitionInput["fields"].(map[string]any)["local_body"].([]int)[b511FailureTransitionIndex] = 255
		b511FailureTransitionCheckFailure(b511FailureTransitionT, b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionInput), b511FailureTransitionIndex == 8, "Hypothesis", true)
	}
	for _, b511FailureTransitionView := range []string{"initialize", "timer-consume"} {
		for _, b511FailureTransitionOperation := range []int{0, 1} {
			b511FailureTransitionInput := b511FailureTransitionBase()
			b511FailureTransitionInput["context"].(map[string]any)["view"] = b511FailureTransitionView
			b511FailureTransitionInput["context"].(map[string]any)["operation"] = b511FailureTransitionOperation
			b511FailureTransitionOutput := b511FailureTransitionCall(b511FailureTransitionT, b511FailureTransitionInput)
			if b511FailureTransitionOutput["state"] != "Unknown" {
				b511FailureTransitionT.Fatal(b511FailureTransitionOutput)
			}
		}
	}
}
