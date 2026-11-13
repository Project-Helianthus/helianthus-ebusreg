package vaillant

import (
	b515PeriodFieldBytes "bytes"
	b515PeriodFieldJSON "encoding/json"
	b515PeriodFieldFmt "fmt"
	b515PeriodFieldOs "os"
	b515PeriodFieldReflect "reflect"
	b515PeriodFieldStrings "strings"
	b515PeriodFieldTesting "testing"
)

func b515PeriodFieldMakeInput(b515PeriodFieldView string) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 21, "direction": "request", "view": b515PeriodFieldView}, "fields": map[string]any{"period0": []int{0}, "period1": []int{54}, "period2": []int{144}}}
}

func b515PeriodFieldRun(b515PeriodFieldT *b515PeriodFieldTesting.T, b515PeriodFieldInput any) map[string]any {
	b515PeriodFieldT.Helper()
	b515PeriodFieldEncoded, b515PeriodFieldError := b515PeriodFieldJSON.Marshal(b515PeriodFieldInput)
	if b515PeriodFieldError != nil {
		b515PeriodFieldT.Fatal(b515PeriodFieldError)
	}
	b515PeriodFieldOutput, b515PeriodFieldError := B515PeriodFieldEvaluate(b515PeriodFieldEncoded)
	if b515PeriodFieldError != nil {
		b515PeriodFieldT.Fatal(b515PeriodFieldError)
	}
	b515PeriodFieldResult := map[string]any{}
	if b515PeriodFieldError := b515PeriodFieldJSON.Unmarshal(b515PeriodFieldOutput, &b515PeriodFieldResult); b515PeriodFieldError != nil {
		b515PeriodFieldT.Fatal(b515PeriodFieldError)
	}
	if b515PeriodFieldResult["physical"] != "Unknown" || b515PeriodFieldResult["qualified"] != false {
		b515PeriodFieldT.Fatalf("qualification: %s", b515PeriodFieldOutput)
	}
	b515PeriodFieldRetained := map[string]any{}
	if b515PeriodFieldError := b515PeriodFieldJSON.Unmarshal(b515PeriodFieldEncoded, &b515PeriodFieldRetained); b515PeriodFieldError != nil {
		b515PeriodFieldT.Fatal(b515PeriodFieldError)
	}
	for _, b515PeriodFieldKey := range []string{"context", "fields", "arguments"} {
		b515PeriodFieldExpected, b515PeriodFieldExpectedPresent := b515PeriodFieldRetained[b515PeriodFieldKey]
		b515PeriodFieldActual, b515PeriodFieldActualPresent := b515PeriodFieldResult[b515PeriodFieldKey]
		if b515PeriodFieldExpectedPresent != b515PeriodFieldActualPresent || !b515PeriodFieldReflect.DeepEqual(b515PeriodFieldActual, b515PeriodFieldExpected) {
			b515PeriodFieldT.Fatalf("retention of %s", b515PeriodFieldKey)
		}
		delete(b515PeriodFieldRetained, b515PeriodFieldKey)
	}
	if !b515PeriodFieldReflect.DeepEqual(b515PeriodFieldResult["extra"], b515PeriodFieldRetained) {
		b515PeriodFieldT.Fatal("extra keys lost")
	}
	return b515PeriodFieldResult
}

func b515PeriodFieldSubset(b515PeriodFieldT *b515PeriodFieldTesting.T, b515PeriodFieldActual, b515PeriodFieldExpected any) {
	b515PeriodFieldT.Helper()
	switch b515PeriodFieldExpectedValue := b515PeriodFieldExpected.(type) {
	case map[string]any:
		b515PeriodFieldActualValue, b515PeriodFieldOK := b515PeriodFieldActual.(map[string]any)
		if !b515PeriodFieldOK {
			b515PeriodFieldT.Fatalf("expected object, got %v", b515PeriodFieldActual)
		}
		for b515PeriodFieldKey, b515PeriodFieldValue := range b515PeriodFieldExpectedValue {
			b515PeriodFieldActualChild, b515PeriodFieldPresent := b515PeriodFieldActualValue[b515PeriodFieldKey]
			if !b515PeriodFieldPresent {
				b515PeriodFieldT.Fatalf("missing %s", b515PeriodFieldKey)
			}
			b515PeriodFieldSubset(b515PeriodFieldT, b515PeriodFieldActualChild, b515PeriodFieldValue)
		}
	case []any:
		b515PeriodFieldActualValue, b515PeriodFieldOK := b515PeriodFieldActual.([]any)
		if !b515PeriodFieldOK || len(b515PeriodFieldActualValue) != len(b515PeriodFieldExpectedValue) {
			b515PeriodFieldT.Fatalf("array mismatch %v", b515PeriodFieldActual)
		}
		for b515PeriodFieldIndex := range b515PeriodFieldExpectedValue {
			b515PeriodFieldSubset(b515PeriodFieldT, b515PeriodFieldActualValue[b515PeriodFieldIndex], b515PeriodFieldExpectedValue[b515PeriodFieldIndex])
		}
	default:
		if !b515PeriodFieldReflect.DeepEqual(b515PeriodFieldActual, b515PeriodFieldExpected) {
			b515PeriodFieldT.Fatalf("got %v, want %v", b515PeriodFieldActual, b515PeriodFieldExpected)
		}
	}
}

func b515PeriodFieldGetResult(b515PeriodFieldOutput map[string]any, b515PeriodFieldKey string) map[string]any {
	return b515PeriodFieldOutput["results"].(map[string]any)[b515PeriodFieldKey].(map[string]any)
}

func b515PeriodFieldNoNumber(b515PeriodFieldT *b515PeriodFieldTesting.T, b515PeriodFieldResult map[string]any) {
	b515PeriodFieldT.Helper()
	for _, b515PeriodFieldKey := range []string{"number", "hour", "minute", "label", "value", "calculation"} {
		if _, b515PeriodFieldPresent := b515PeriodFieldResult[b515PeriodFieldKey]; b515PeriodFieldPresent {
			b515PeriodFieldT.Fatalf("unexpected %s in %v", b515PeriodFieldKey, b515PeriodFieldResult)
		}
	}
}

func TestAb09Aq09Fixture(b515PeriodFieldT *b515PeriodFieldTesting.T) {
	b515PeriodFieldData, b515PeriodFieldError := b515PeriodFieldOs.ReadFile("testdata/ab09aq09_b515_period_fields.json")
	if b515PeriodFieldError != nil {
		b515PeriodFieldT.Fatal(b515PeriodFieldError)
	}
	var b515PeriodFieldCases []map[string]any
	if b515PeriodFieldError := b515PeriodFieldJSON.Unmarshal(b515PeriodFieldData, &b515PeriodFieldCases); b515PeriodFieldError != nil {
		b515PeriodFieldT.Fatal(b515PeriodFieldError)
	}
	if len(b515PeriodFieldCases) != 5 {
		b515PeriodFieldT.Fatal("fixture cases missing")
	}
	for _, b515PeriodFieldCase := range b515PeriodFieldCases {
		b515PeriodFieldT.Run(b515PeriodFieldCase["id"].(string), func(b515PeriodFieldT *b515PeriodFieldTesting.T) {
			b515PeriodFieldOutput := b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldCase["input"])
			b515PeriodFieldSubset(b515PeriodFieldT, b515PeriodFieldOutput, b515PeriodFieldCase["expected"])
			if b515PeriodFieldCase["id"] == "ab09aq09_unmapped" {
				for _, b515PeriodFieldKey := range []string{"period0", "period1", "minutes"} {
					b515PeriodFieldNoNumber(b515PeriodFieldT, b515PeriodFieldGetResult(b515PeriodFieldOutput, b515PeriodFieldKey))
				}
			}
			if b515PeriodFieldCase["id"] == "ab09aq09_unmapped_equal" {
				b515PeriodFieldNoNumber(b515PeriodFieldT, b515PeriodFieldGetResult(b515PeriodFieldOutput, "start_order"))
			}
		})
	}
}

func TestAb09Aq09ExhaustPeriods(b515PeriodFieldT *b515PeriodFieldTesting.T) {
	b515PeriodFieldExpectedMinutes := 0
	b515PeriodFieldExpectedHour, b515PeriodFieldExpectedMinute := 0, 0
	for b515PeriodFieldRaw := 0; b515PeriodFieldRaw <= 255; b515PeriodFieldRaw++ {
		b515PeriodFieldInput := b515PeriodFieldMakeInput("periods")
		for _, b515PeriodFieldKey := range []string{"period0", "period1", "period2"} {
			b515PeriodFieldInput["fields"].(map[string]any)[b515PeriodFieldKey] = []int{b515PeriodFieldRaw}
		}
		b515PeriodFieldOutput := b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldInput)
		if b515PeriodFieldOutput["state"] != "complete" {
			b515PeriodFieldT.Fatal(b515PeriodFieldOutput)
		}
		for _, b515PeriodFieldKey := range []string{"period0", "period1", "period2"} {
			b515PeriodFieldResult := b515PeriodFieldGetResult(b515PeriodFieldOutput, b515PeriodFieldKey)
			if b515PeriodFieldRaw <= 144 {
				if b515PeriodFieldResult["number"] != float64(b515PeriodFieldExpectedMinutes) || b515PeriodFieldResult["hour"] != float64(b515PeriodFieldExpectedHour) || b515PeriodFieldResult["minute"] != float64(b515PeriodFieldExpectedMinute) || b515PeriodFieldResult["label"] != b515PeriodFieldFmt.Sprintf("%02d:%02d", b515PeriodFieldExpectedHour, b515PeriodFieldExpectedMinute) || b515PeriodFieldResult["state"] != "Hypothesis" || b515PeriodFieldResult["calculation"] != "Proven" {
					b515PeriodFieldT.Fatalf("raw %d: %v", b515PeriodFieldRaw, b515PeriodFieldResult)
				}
			} else {
				if b515PeriodFieldResult["state"] != "Unknown" {
					b515PeriodFieldT.Fatal(b515PeriodFieldResult)
				}
				b515PeriodFieldNoNumber(b515PeriodFieldT, b515PeriodFieldResult)
			}
		}
		b515PeriodFieldExpectedMinutes += 10
		b515PeriodFieldExpectedMinute += 10
		if b515PeriodFieldExpectedMinute == 60 {
			b515PeriodFieldExpectedHour++
			b515PeriodFieldExpectedMinute = 0
		}
	}
}

func TestAb09Aq09Inverse(b515PeriodFieldT *b515PeriodFieldTesting.T) {
	for b515PeriodFieldMinutes := 0; b515PeriodFieldMinutes <= 1440; b515PeriodFieldMinutes++ {
		b515PeriodFieldInput := b515PeriodFieldMakeInput("periods")
		b515PeriodFieldInput["arguments"] = map[string]any{"minutes": b515PeriodFieldMinutes}
		b515PeriodFieldOutput := b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldInput)
		b515PeriodFieldResult := b515PeriodFieldGetResult(b515PeriodFieldOutput, "minutes")
		if b515PeriodFieldMinutes%10 != 0 {
			if b515PeriodFieldResult["state"] != "Unknown" {
				b515PeriodFieldT.Fatal(b515PeriodFieldResult)
			}
			b515PeriodFieldNoNumber(b515PeriodFieldT, b515PeriodFieldResult)
		} else {
			if b515PeriodFieldResult["number"] != float64(b515PeriodFieldMinutes/10) || b515PeriodFieldResult["state"] != "Hypothesis" || b515PeriodFieldResult["calculation"] != "Proven" {
				b515PeriodFieldT.Fatal(b515PeriodFieldResult)
			}
			b515PeriodFieldInput["fields"].(map[string]any)["period0"] = []int{b515PeriodFieldMinutes / 10}
			if b515PeriodFieldGetResult(b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldInput), "period0")["number"] != float64(b515PeriodFieldMinutes) {
				b515PeriodFieldT.Fatal("round trip")
			}
		}
	}
}

func TestAb09Aq09ContextAxes(b515PeriodFieldT *b515PeriodFieldTesting.T) {
	for b515PeriodFieldKey, b515PeriodFieldValues := range map[string][]any{"primary": {180, 0, 255}, "secondary": {20, 0, 255}, "operation": {0, 1, 255}, "direction": {"negative", "local", ""}, "view": {"other", ""}} {
		for _, b515PeriodFieldValue := range b515PeriodFieldValues {
			b515PeriodFieldInput := b515PeriodFieldMakeInput("periods")
			b515PeriodFieldInput["context"].(map[string]any)[b515PeriodFieldKey] = b515PeriodFieldValue
			b515PeriodFieldOutput := b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldInput)
			if b515PeriodFieldOutput["state"] != "Unknown" || len(b515PeriodFieldOutput["results"].(map[string]any)) != 0 {
				b515PeriodFieldT.Fatalf("%s=%v: %v", b515PeriodFieldKey, b515PeriodFieldValue, b515PeriodFieldOutput)
			}
		}
	}
	for _, b515PeriodFieldKey := range []string{"primary", "secondary", "direction", "view"} {
		b515PeriodFieldInput := b515PeriodFieldMakeInput("periods")
		delete(b515PeriodFieldInput["context"].(map[string]any), b515PeriodFieldKey)
		if b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldInput)["state"] != "Unknown" {
			b515PeriodFieldT.Fatal(b515PeriodFieldKey)
		}
	}
	for _, b515PeriodFieldKey := range []string{"family", "variant"} {
		for _, b515PeriodFieldValue := range []any{nil, false, 0, "", []any{}, map[string]any{}, "synthetic-family"} {
			b515PeriodFieldInput := b515PeriodFieldMakeInput("periods")
			b515PeriodFieldInput["context"].(map[string]any)[b515PeriodFieldKey] = b515PeriodFieldValue
			b515PeriodFieldOutput := b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldInput)
			b515PeriodFieldBaseline := b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldMakeInput("periods"))
			if !b515PeriodFieldReflect.DeepEqual(b515PeriodFieldOutput["results"], b515PeriodFieldBaseline["results"]) {
				b515PeriodFieldT.Fatal("opaque identity changed arithmetic")
			}
			b515PeriodFieldEncoded, _ := b515PeriodFieldJSON.Marshal(b515PeriodFieldValue)
			var b515PeriodFieldExpected any
			if b515PeriodFieldError := b515PeriodFieldJSON.Unmarshal(b515PeriodFieldEncoded, &b515PeriodFieldExpected); b515PeriodFieldError != nil {
				b515PeriodFieldT.Fatal(b515PeriodFieldError)
			}
			if !b515PeriodFieldReflect.DeepEqual(b515PeriodFieldOutput["context"].(map[string]any)[b515PeriodFieldKey], b515PeriodFieldExpected) {
				b515PeriodFieldT.Fatal("opaque identity lost")
			}
		}
	}
}

func TestAb09Aq09PartialAndRetention(b515PeriodFieldT *b515PeriodFieldTesting.T) {
	for _, b515PeriodFieldKey := range []string{"period0", "period1", "period2"} {
		for _, b515PeriodFieldWidth := range []int{-1, 0, 2} {
			b515PeriodFieldInput := b515PeriodFieldMakeInput("periods")
			b515PeriodFieldFields := b515PeriodFieldInput["fields"].(map[string]any)
			if b515PeriodFieldWidth == -1 {
				delete(b515PeriodFieldFields, b515PeriodFieldKey)
			} else {
				b515PeriodFieldFields[b515PeriodFieldKey] = make([]int, b515PeriodFieldWidth)
			}
			b515PeriodFieldFields["opaque"] = []int{}
			b515PeriodFieldInput["extra"] = map[string]any{"null": nil, "zero": 0, "false": false, "empty": []any{}}
			b515PeriodFieldInput["state"] = "retained"
			b515PeriodFieldOutput := b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldInput)
			if b515PeriodFieldOutput["state"] != "partial" {
				b515PeriodFieldT.Fatal(b515PeriodFieldOutput)
			}
			b515PeriodFieldResult := b515PeriodFieldGetResult(b515PeriodFieldOutput, b515PeriodFieldKey)
			if b515PeriodFieldResult["state"] != "incomplete" {
				b515PeriodFieldT.Fatal(b515PeriodFieldResult)
			}
			b515PeriodFieldNoNumber(b515PeriodFieldT, b515PeriodFieldResult)
			if _, b515PeriodFieldPresent := b515PeriodFieldResult["raw"]; b515PeriodFieldPresent != (b515PeriodFieldWidth != -1) {
				b515PeriodFieldT.Fatal("absence lost")
			}
			for _, b515PeriodFieldSibling := range []string{"period0", "period1", "period2"} {
				if b515PeriodFieldSibling != b515PeriodFieldKey && b515PeriodFieldGetResult(b515PeriodFieldOutput, b515PeriodFieldSibling)["state"] != "Hypothesis" {
					b515PeriodFieldT.Fatal("sibling lost")
				}
			}
			if b515PeriodFieldOutput["extra"].(map[string]any)["state"] != "retained" {
				b515PeriodFieldT.Fatal("extra collision")
			}
			if len(b515PeriodFieldGetResult(b515PeriodFieldOutput, "opaque")["raw"].([]any)) != 0 {
				b515PeriodFieldT.Fatal("empty opaque lost")
			}
			b515PeriodFieldSubset(b515PeriodFieldT, b515PeriodFieldOutput["extra"], map[string]any{"extra": map[string]any{"null": nil, "zero": float64(0), "false": false, "empty": []any{}}})
		}
	}
}

func TestAb09Aq09IntervalsAxes(b515PeriodFieldT *b515PeriodFieldTesting.T) {
	for b515PeriodFieldIndex := 0; b515PeriodFieldIndex < 3; b515PeriodFieldIndex++ {
		for b515PeriodFieldEndpoint := 0; b515PeriodFieldEndpoint < 2; b515PeriodFieldEndpoint++ {
			for _, b515PeriodFieldValue := range []int{0, 10, 20, 40, 144, 145, 255} {
				b515PeriodFieldPairs := [3][2]int{{10, 10}, {20, 30}, {40, 50}}
				b515PeriodFieldPairs[b515PeriodFieldIndex][b515PeriodFieldEndpoint] = b515PeriodFieldValue
				b515PeriodFieldInput := b515PeriodFieldMakeInput("intervals")
				b515PeriodFieldInput["arguments"] = map[string]any{"pairs": b515PeriodFieldPairs}
				b515PeriodFieldOutput := b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldInput)
				b515PeriodFieldOrder := b515PeriodFieldGetResult(b515PeriodFieldOutput, "start_order")
				if b515PeriodFieldValue > 144 {
					if b515PeriodFieldOrder["state"] != "Unknown" {
						b515PeriodFieldT.Fatal(b515PeriodFieldOrder)
					}
					b515PeriodFieldNoNumber(b515PeriodFieldT, b515PeriodFieldOrder)
				} else {
					b515PeriodFieldExpected := b515PeriodFieldPairs[0][0] < b515PeriodFieldPairs[1][0] && b515PeriodFieldPairs[1][0] < b515PeriodFieldPairs[2][0]
					if b515PeriodFieldOrder["value"] != b515PeriodFieldExpected {
						b515PeriodFieldT.Fatal(b515PeriodFieldOrder)
					}
				}
				for b515PeriodFieldPairIndex, b515PeriodFieldPairResult := range b515PeriodFieldGetResult(b515PeriodFieldOutput, "pairs")["alternatives"].([]any) {
					b515PeriodFieldExpected := "enabled-candidate"
					if b515PeriodFieldPairs[b515PeriodFieldPairIndex][0] == b515PeriodFieldPairs[b515PeriodFieldPairIndex][1] {
						b515PeriodFieldExpected = "disabled"
					}
					b515PeriodFieldResult := b515PeriodFieldPairResult.(map[string]any)
					if b515PeriodFieldResult["label"] != b515PeriodFieldExpected || b515PeriodFieldResult["start"] != float64(b515PeriodFieldPairs[b515PeriodFieldPairIndex][0]) || b515PeriodFieldResult["end"] != float64(b515PeriodFieldPairs[b515PeriodFieldPairIndex][1]) {
						b515PeriodFieldT.Fatal(b515PeriodFieldResult)
					}
					if _, b515PeriodFieldPresent := b515PeriodFieldResult["deletion"]; b515PeriodFieldPresent {
						b515PeriodFieldT.Fatal("invented deletion")
					}
				}
			}
		}
	}
	b515PeriodFieldInput := b515PeriodFieldMakeInput("intervals")
	b515PeriodFieldInput["arguments"] = map[string]any{"pairs": [3][2]int{{140, 1}, {30, 10}, {20, 40}}}
	if b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldInput)["state"] != "complete" {
		b515PeriodFieldT.Fatal("reversed intervals rejected")
	}
	b515PeriodFieldInput["context"].(map[string]any)["direction"] = "positive"
	if b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldInput)["state"] != "Unknown" {
		b515PeriodFieldT.Fatal("interval direction admitted")
	}
}

func TestAb09Aq09RejectJSON(b515PeriodFieldT *b515PeriodFieldTesting.T) {
	for _, b515PeriodFieldInput := range []string{
		``, `{`, `null`, `[]`, `{} {}`, `{} true`, `{"context":{},"context":{}}`,
		`{"context":{"family":{"x":0,"x":1}},"fields":{}}`,
		`{"context":{},"fields":{},"extra":[{"x":0,"x":1}]}`,
		`{"context":{},"fields":{"x":[NaN]}}`, `{"context":{},"fields":{},"extra":1e400}`,
		`{"context":{},"fields":{"x":[0],"x":[]}}`,
		`{"context":{},"fields":{},"arguments":{"x":0,"x":1}}`,
	} {
		b515PeriodFieldOutput, b515PeriodFieldError := B515PeriodFieldEvaluate([]byte(b515PeriodFieldInput))
		if b515PeriodFieldError == nil || b515PeriodFieldOutput != nil {
			b515PeriodFieldT.Fatalf("accepted %q", b515PeriodFieldInput)
		}
	}
}

func TestAb09Aq09RejectKinds(b515PeriodFieldT *b515PeriodFieldTesting.T) {
	for _, b515PeriodFieldValue := range []any{nil, true, "1", -1, 256, 1.5, []any{}, map[string]any{}} {
		for _, b515PeriodFieldKey := range []string{"primary", "secondary", "operation", "byte", "minutes", "endpoint"} {
			b515PeriodFieldInput := b515PeriodFieldMakeInput("periods")
			switch b515PeriodFieldKey {
			case "byte":
				b515PeriodFieldInput["fields"].(map[string]any)["period0"] = []any{b515PeriodFieldValue}
			case "minutes":
				if b515PeriodFieldValue == 256 {
					continue
				}
				b515PeriodFieldInput["arguments"] = map[string]any{"minutes": b515PeriodFieldValue}
			case "endpoint":
				b515PeriodFieldInput["arguments"] = map[string]any{"pairs": []any{[]any{b515PeriodFieldValue, 10}, []int{20, 30}, []int{40, 50}}}
			default:
				b515PeriodFieldInput["context"].(map[string]any)[b515PeriodFieldKey] = b515PeriodFieldValue
			}
			b515PeriodFieldEncoded, _ := b515PeriodFieldJSON.Marshal(b515PeriodFieldInput)
			b515PeriodFieldOutput, b515PeriodFieldError := B515PeriodFieldEvaluate(b515PeriodFieldEncoded)
			if b515PeriodFieldError == nil || b515PeriodFieldOutput != nil {
				b515PeriodFieldT.Fatalf("accepted %s", b515PeriodFieldEncoded)
			}
		}
	}
	for _, b515PeriodFieldPatch := range []string{
		`"arguments":{"minutes":1441}`, `"arguments":null`, `"fields":null`, `"context":null`,
		`"arguments":{"pairs":[]}`, `"arguments":{"pairs":[[1,2],[3,4],[5]]}`,
		`"arguments":{"pairs":[[1,2],[3,4],[5,6,7]]}`, `"fields":{"period0":null}`,
		`"context":{"direction":false}`, `"context":{"view":0}`,
	} {
		b515PeriodFieldInput := b515PeriodFieldMakeInput("periods")
		b515PeriodFieldPatchObject := map[string]any{}
		if b515PeriodFieldError := b515PeriodFieldJSON.Unmarshal([]byte("{"+b515PeriodFieldPatch+"}"), &b515PeriodFieldPatchObject); b515PeriodFieldError != nil {
			b515PeriodFieldT.Fatal(b515PeriodFieldError)
		}
		for b515PeriodFieldKey, b515PeriodFieldValue := range b515PeriodFieldPatchObject {
			b515PeriodFieldInput[b515PeriodFieldKey] = b515PeriodFieldValue
		}
		b515PeriodFieldEncoded, _ := b515PeriodFieldJSON.Marshal(b515PeriodFieldInput)
		if b515PeriodFieldOutput, b515PeriodFieldError := B515PeriodFieldEvaluate(b515PeriodFieldEncoded); b515PeriodFieldError == nil || b515PeriodFieldOutput != nil {
			b515PeriodFieldT.Fatalf("accepted %s", b515PeriodFieldEncoded)
		}
	}
	b515PeriodFieldEncoded, _ := b515PeriodFieldJSON.Marshal(b515PeriodFieldMakeInput("intervals"))
	if b515PeriodFieldOutput, b515PeriodFieldError := B515PeriodFieldEvaluate(b515PeriodFieldEncoded); b515PeriodFieldError == nil || b515PeriodFieldOutput != nil {
		b515PeriodFieldT.Fatal("missing pairs accepted")
	}
}

func TestAb09Aq09IndependentMemory(b515PeriodFieldT *b515PeriodFieldTesting.T) {
	b515PeriodFieldEncoded, _ := b515PeriodFieldJSON.Marshal(b515PeriodFieldMakeInput("periods"))
	b515PeriodFieldOriginal := append([]byte{}, b515PeriodFieldEncoded...)
	b515PeriodFieldFirst, b515PeriodFieldError := B515PeriodFieldEvaluate(b515PeriodFieldEncoded)
	if b515PeriodFieldError != nil {
		b515PeriodFieldT.Fatal(b515PeriodFieldError)
	}
	b515PeriodFieldSaved := append([]byte{}, b515PeriodFieldFirst...)
	for b515PeriodFieldIndex := range b515PeriodFieldEncoded {
		b515PeriodFieldEncoded[b515PeriodFieldIndex] = ' '
	}
	if !b515PeriodFieldBytes.Equal(b515PeriodFieldSaved, b515PeriodFieldFirst) {
		b515PeriodFieldT.Fatal("input aliases output")
	}
	for b515PeriodFieldIndex := range b515PeriodFieldFirst {
		b515PeriodFieldFirst[b515PeriodFieldIndex] = ' '
	}
	b515PeriodFieldSecond, b515PeriodFieldError := B515PeriodFieldEvaluate(b515PeriodFieldOriginal)
	if b515PeriodFieldError != nil || !b515PeriodFieldBytes.Equal(b515PeriodFieldSaved, b515PeriodFieldSecond) {
		b515PeriodFieldT.Fatal("calls share storage")
	}
	if b515PeriodFieldStrings.Contains(string(b515PeriodFieldSecond), "qualified\":true") {
		b515PeriodFieldT.Fatal("qualified")
	}
}

func TestAb09Aq09FieldAxes(b515PeriodFieldT *b515PeriodFieldTesting.T) {
	for _, b515PeriodFieldKey := range []string{"period0", "period1", "period2"} {
		for _, b515PeriodFieldValue := range []int{0, 1, 54, 144, 145, 255} {
			b515PeriodFieldInput := b515PeriodFieldMakeInput("periods")
			b515PeriodFieldInput["fields"].(map[string]any)[b515PeriodFieldKey] = []int{b515PeriodFieldValue}
			b515PeriodFieldOutput := b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldInput)
			b515PeriodFieldBaseline := b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldMakeInput("periods"))
			for _, b515PeriodFieldSibling := range []string{"period0", "period1", "period2"} {
				b515PeriodFieldResult := b515PeriodFieldGetResult(b515PeriodFieldOutput, b515PeriodFieldSibling)
				if b515PeriodFieldSibling != b515PeriodFieldKey {
					if !b515PeriodFieldReflect.DeepEqual(b515PeriodFieldResult, b515PeriodFieldGetResult(b515PeriodFieldBaseline, b515PeriodFieldSibling)) {
						b515PeriodFieldT.Fatal("field cross-talk")
					}
				} else {
					if !b515PeriodFieldReflect.DeepEqual(b515PeriodFieldResult["raw"], []any{float64(b515PeriodFieldValue)}) {
						b515PeriodFieldT.Fatal("raw field lost")
					}
					if b515PeriodFieldValue <= 144 {
						if b515PeriodFieldResult["number"] != float64(b515PeriodFieldValue*10) {
							b515PeriodFieldT.Fatal("field ignored")
						}
					} else {
						b515PeriodFieldNoNumber(b515PeriodFieldT, b515PeriodFieldResult)
					}
				}
			}
		}
	}
	for _, b515PeriodFieldKey := range []string{"weekday", "channel", "mode"} {
		for _, b515PeriodFieldValue := range [][]int{{}, {0}, {1}, {255}, {0, 255}} {
			b515PeriodFieldInput := b515PeriodFieldMakeInput("periods")
			b515PeriodFieldInput["fields"].(map[string]any)[b515PeriodFieldKey] = b515PeriodFieldValue
			b515PeriodFieldOutput := b515PeriodFieldRun(b515PeriodFieldT, b515PeriodFieldInput)
			b515PeriodFieldResult := b515PeriodFieldGetResult(b515PeriodFieldOutput, b515PeriodFieldKey)
			if b515PeriodFieldResult["state"] != "Unknown" {
				b515PeriodFieldT.Fatal("opaque carrier decoded")
			}
			b515PeriodFieldNoNumber(b515PeriodFieldT, b515PeriodFieldResult)
		}
	}
}
