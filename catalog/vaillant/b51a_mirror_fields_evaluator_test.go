package vaillant

import (
	b51AMirrorFieldEvaluatorJSON "encoding/json"
	b51AMirrorFieldEvaluatorOs "os"
	b51AMirrorFieldEvaluatorReflect "reflect"
	b51AMirrorFieldEvaluatorTesting "testing"
)

func b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT *b51AMirrorFieldEvaluatorTesting.T) []map[string]any {
	b51AMirrorFieldEvaluatorT.Helper()
	b51AMirrorFieldEvaluatorBytes, b51AMirrorFieldEvaluatorErr := b51AMirrorFieldEvaluatorOs.ReadFile("testdata/ab19ap19_b51a_mirror_fields.json")
	if b51AMirrorFieldEvaluatorErr != nil {
		b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorErr)
	}
	var b51AMirrorFieldEvaluatorCases []map[string]any
	if b51AMirrorFieldEvaluatorErr := b51AMirrorFieldEvaluatorJSON.Unmarshal(b51AMirrorFieldEvaluatorBytes, &b51AMirrorFieldEvaluatorCases); b51AMirrorFieldEvaluatorErr != nil {
		b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorErr)
	}
	return b51AMirrorFieldEvaluatorCases
}
func b51AMirrorFieldEvaluatorRun(b51AMirrorFieldEvaluatorT *b51AMirrorFieldEvaluatorTesting.T, b51AMirrorFieldEvaluatorInput any) map[string]any {
	b51AMirrorFieldEvaluatorT.Helper()
	b51AMirrorFieldEvaluatorBytes, _ := b51AMirrorFieldEvaluatorJSON.Marshal(b51AMirrorFieldEvaluatorInput)
	b51AMirrorFieldEvaluatorOutput, b51AMirrorFieldEvaluatorErr := B51AMirrorFieldEvaluate(b51AMirrorFieldEvaluatorBytes)
	if b51AMirrorFieldEvaluatorErr != nil {
		b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorErr)
	}
	var b51AMirrorFieldEvaluatorResult map[string]any
	if b51AMirrorFieldEvaluatorErr := b51AMirrorFieldEvaluatorJSON.Unmarshal(b51AMirrorFieldEvaluatorOutput, &b51AMirrorFieldEvaluatorResult); b51AMirrorFieldEvaluatorErr != nil {
		b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorErr)
	}
	return b51AMirrorFieldEvaluatorResult
}
func b51AMirrorFieldEvaluatorField(b51AMirrorFieldEvaluatorOutput map[string]any, b51AMirrorFieldEvaluatorName string) map[string]any {
	return b51AMirrorFieldEvaluatorOutput["fields"].(map[string]any)[b51AMirrorFieldEvaluatorName].(map[string]any)
}
func TestAb19Ap19Vectors(b51AMirrorFieldEvaluatorT *b51AMirrorFieldEvaluatorTesting.T) {
	for _, b51AMirrorFieldEvaluatorCase := range b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT) {
		b51AMirrorFieldEvaluatorT.Run(b51AMirrorFieldEvaluatorCase["id"].(string), func(b51AMirrorFieldEvaluatorT *b51AMirrorFieldEvaluatorTesting.T) {
			b51AMirrorFieldEvaluatorInput := b51AMirrorFieldEvaluatorCase["input"].(map[string]any)
			b51AMirrorFieldEvaluatorOutput := b51AMirrorFieldEvaluatorRun(b51AMirrorFieldEvaluatorT, b51AMirrorFieldEvaluatorInput)
			if b51AMirrorFieldEvaluatorOutput["state"] != "complete" || b51AMirrorFieldEvaluatorOutput["physical"] != "Unknown" || b51AMirrorFieldEvaluatorOutput["qualified"] != false || !b51AMirrorFieldEvaluatorReflect.DeepEqual(b51AMirrorFieldEvaluatorOutput["context"], b51AMirrorFieldEvaluatorInput["context"]) || !b51AMirrorFieldEvaluatorReflect.DeepEqual(b51AMirrorFieldEvaluatorOutput["raw"], b51AMirrorFieldEvaluatorInput["fields"]) {
				b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorOutput)
			}
			for b51AMirrorFieldEvaluatorName, b51AMirrorFieldEvaluatorWant := range b51AMirrorFieldEvaluatorCase["numbers"].(map[string]any) {
				b51AMirrorFieldEvaluatorField := b51AMirrorFieldEvaluatorField(b51AMirrorFieldEvaluatorOutput, b51AMirrorFieldEvaluatorName)
				if b51AMirrorFieldEvaluatorField["number"] != b51AMirrorFieldEvaluatorWant || b51AMirrorFieldEvaluatorField["state"] != "Hypothesis" || b51AMirrorFieldEvaluatorField["calculation"] != "Proven" {
					b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorField)
				}
			}
			if b51AMirrorFieldEvaluatorAddress, b51AMirrorFieldEvaluatorOK := b51AMirrorFieldEvaluatorOutput["fields"].(map[string]any)["address"]; b51AMirrorFieldEvaluatorOK && b51AMirrorFieldEvaluatorAddress.(map[string]any)["logical_validity"] != "Unknown" {
				b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorAddress)
			}
		})
	}
}
func TestAb19Ap19IndependentFields(b51AMirrorFieldEvaluatorT *b51AMirrorFieldEvaluatorTesting.T) {
	for b51AMirrorFieldEvaluatorIndex := 0; b51AMirrorFieldEvaluatorIndex < 4; b51AMirrorFieldEvaluatorIndex++ {
		for b51AMirrorFieldEvaluatorName := range b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT)[b51AMirrorFieldEvaluatorIndex]["input"].(map[string]any)["fields"].(map[string]any) {
			for _, b51AMirrorFieldEvaluatorReplacement := range []any{nil, []any{}, []any{float64(0)}, []any{float64(255)}, []any{float64(0), float64(0)}, []any{float64(255), float64(255)}, []any{float64(1), float64(2), float64(3)}} {
				b51AMirrorFieldEvaluatorInput := b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT)[b51AMirrorFieldEvaluatorIndex]["input"].(map[string]any)
				b51AMirrorFieldEvaluatorRaw := b51AMirrorFieldEvaluatorInput["fields"].(map[string]any)
				b51AMirrorFieldEvaluatorOriginal := b51AMirrorFieldEvaluatorRun(b51AMirrorFieldEvaluatorT, b51AMirrorFieldEvaluatorInput)
				if b51AMirrorFieldEvaluatorReplacement == nil {
					delete(b51AMirrorFieldEvaluatorRaw, b51AMirrorFieldEvaluatorName)
				} else {
					b51AMirrorFieldEvaluatorRaw[b51AMirrorFieldEvaluatorName] = b51AMirrorFieldEvaluatorReplacement
				}
				b51AMirrorFieldEvaluatorOutput := b51AMirrorFieldEvaluatorRun(b51AMirrorFieldEvaluatorT, b51AMirrorFieldEvaluatorInput)
				for b51AMirrorFieldEvaluatorSibling, b51AMirrorFieldEvaluatorExpected := range b51AMirrorFieldEvaluatorOriginal["fields"].(map[string]any) {
					if b51AMirrorFieldEvaluatorSibling != b51AMirrorFieldEvaluatorName && !b51AMirrorFieldEvaluatorReflect.DeepEqual(b51AMirrorFieldEvaluatorExpected, b51AMirrorFieldEvaluatorOutput["fields"].(map[string]any)[b51AMirrorFieldEvaluatorSibling]) {
						b51AMirrorFieldEvaluatorT.Fatal("sibling changed", b51AMirrorFieldEvaluatorSibling)
					}
				}
				b51AMirrorFieldEvaluatorResult := b51AMirrorFieldEvaluatorField(b51AMirrorFieldEvaluatorOutput, b51AMirrorFieldEvaluatorName)
				if b51AMirrorFieldEvaluatorReplacement == nil {
					if _, b51AMirrorFieldEvaluatorPresent := b51AMirrorFieldEvaluatorResult["raw"]; b51AMirrorFieldEvaluatorPresent {
						b51AMirrorFieldEvaluatorT.Fatal("invented missing raw")
					}
				} else if !b51AMirrorFieldEvaluatorReflect.DeepEqual(b51AMirrorFieldEvaluatorResult["raw"], b51AMirrorFieldEvaluatorReplacement) {
					b51AMirrorFieldEvaluatorT.Fatal("changed field raw lost", b51AMirrorFieldEvaluatorName)
				}
				b51AMirrorFieldEvaluatorWidth := 1
				if b51AMirrorFieldEvaluatorName == "address" || b51AMirrorFieldEvaluatorName == "mirror_address" || b51AMirrorFieldEvaluatorName == "symbol_address" {
					b51AMirrorFieldEvaluatorWidth = 2
				}
				b51AMirrorFieldEvaluatorIncomplete := b51AMirrorFieldEvaluatorReplacement == nil
				if b51AMirrorFieldEvaluatorReplacement != nil && b51AMirrorFieldEvaluatorName != "data" {
					b51AMirrorFieldEvaluatorIncomplete = len(b51AMirrorFieldEvaluatorReplacement.([]any)) != b51AMirrorFieldEvaluatorWidth
				}
				if b51AMirrorFieldEvaluatorIncomplete {
					if b51AMirrorFieldEvaluatorResult["state"] != "incomplete" || b51AMirrorFieldEvaluatorResult["number"] != nil || b51AMirrorFieldEvaluatorOutput["state"] != "partial" {
						b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorOutput)
					}
				} else {
					if b51AMirrorFieldEvaluatorResult["state"] != "Hypothesis" {
						b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorResult)
					}
					if b51AMirrorFieldEvaluatorName == "data" {
						if b51AMirrorFieldEvaluatorResult["label"] != "opaque" || b51AMirrorFieldEvaluatorResult["number"] != nil {
							b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorResult)
						}
					} else if b51AMirrorFieldEvaluatorWidth == 1 && b51AMirrorFieldEvaluatorResult["number"] != b51AMirrorFieldEvaluatorReplacement.([]any)[0] {
						b51AMirrorFieldEvaluatorT.Fatal("scalar value changed", b51AMirrorFieldEvaluatorName, b51AMirrorFieldEvaluatorResult)
					}
				}
			}
		}
	}
}
func TestAb19Ap19StatusAndFlags(b51AMirrorFieldEvaluatorT *b51AMirrorFieldEvaluatorTesting.T) {
	for b51AMirrorFieldEvaluatorStatus := 0; b51AMirrorFieldEvaluatorStatus < 256; b51AMirrorFieldEvaluatorStatus++ {
		b51AMirrorFieldEvaluatorInput := b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT)[1]["input"].(map[string]any)
		b51AMirrorFieldEvaluatorInput["fields"].(map[string]any)["status"] = []int{b51AMirrorFieldEvaluatorStatus}
		b51AMirrorFieldEvaluatorOutput := b51AMirrorFieldEvaluatorRun(b51AMirrorFieldEvaluatorT, b51AMirrorFieldEvaluatorInput)
		b51AMirrorFieldEvaluatorAlternatives := b51AMirrorFieldEvaluatorField(b51AMirrorFieldEvaluatorOutput, "status")["alternatives"].(map[string]any)
		if len(b51AMirrorFieldEvaluatorAlternatives) != 2 {
			b51AMirrorFieldEvaluatorT.Fatal("both status interpretations required", b51AMirrorFieldEvaluatorAlternatives)
		}
		b51AMirrorFieldEvaluatorBoolean := b51AMirrorFieldEvaluatorAlternatives["whole_boolean"].(map[string]any)
		b51AMirrorFieldEvaluatorBit := b51AMirrorFieldEvaluatorAlternatives["bit0_out_of_area"].(map[string]any)
		if len(b51AMirrorFieldEvaluatorAlternatives) != 2 || b51AMirrorFieldEvaluatorBit["value"] != (b51AMirrorFieldEvaluatorStatus%2 == 1) {
			b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorAlternatives)
		}
		if b51AMirrorFieldEvaluatorStatus < 2 {
			if b51AMirrorFieldEvaluatorBoolean["value"] != (b51AMirrorFieldEvaluatorStatus == 1) {
				b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorBoolean)
			}
		} else if b51AMirrorFieldEvaluatorBoolean["state"] != "Unknown" || b51AMirrorFieldEvaluatorBoolean["value"] != nil {
			b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorBoolean)
		}
		if len(b51AMirrorFieldEvaluatorField(b51AMirrorFieldEvaluatorOutput, "data")["raw"].([]any)) != 2 {
			b51AMirrorFieldEvaluatorT.Fatal("returned data resized")
		}
	}
	for _, b51AMirrorFieldEvaluatorIndex := range []int{0, 2, 3} {
		b51AMirrorFieldEvaluatorNames := map[int][]string{0: {"triple-storage"}, 2: {"triple-storage", "direct-update"}, 3: {"not-written", "invalid-symbol", "invalid-parameter", "storage-error"}}[b51AMirrorFieldEvaluatorIndex]
		for b51AMirrorFieldEvaluatorBit := -1; b51AMirrorFieldEvaluatorBit < 8; b51AMirrorFieldEvaluatorBit++ {
			b51AMirrorFieldEvaluatorNumber := 0
			if b51AMirrorFieldEvaluatorBit >= 0 {
				b51AMirrorFieldEvaluatorNumber = 1 << b51AMirrorFieldEvaluatorBit
			}
			b51AMirrorFieldEvaluatorInput := b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT)[b51AMirrorFieldEvaluatorIndex]["input"].(map[string]any)
			b51AMirrorFieldEvaluatorInput["fields"].(map[string]any)["flags"] = []int{b51AMirrorFieldEvaluatorNumber}
			b51AMirrorFieldEvaluatorResult := b51AMirrorFieldEvaluatorField(b51AMirrorFieldEvaluatorRun(b51AMirrorFieldEvaluatorT, b51AMirrorFieldEvaluatorInput), "flags")
			for b51AMirrorFieldEvaluatorIndex, b51AMirrorFieldEvaluatorName := range b51AMirrorFieldEvaluatorNames {
				if b51AMirrorFieldEvaluatorResult["flags"].(map[string]any)[b51AMirrorFieldEvaluatorName] != (b51AMirrorFieldEvaluatorIndex == b51AMirrorFieldEvaluatorBit) {
					b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorResult)
				}
			}
			b51AMirrorFieldEvaluatorReserved := 0
			if b51AMirrorFieldEvaluatorBit >= len(b51AMirrorFieldEvaluatorNames) {
				b51AMirrorFieldEvaluatorReserved = b51AMirrorFieldEvaluatorNumber
			}
			b51AMirrorFieldEvaluatorKey := "reserved"
			if b51AMirrorFieldEvaluatorIndex == 0 {
				b51AMirrorFieldEvaluatorKey = "unresolved"
			}
			if b51AMirrorFieldEvaluatorResult[b51AMirrorFieldEvaluatorKey] != float64(b51AMirrorFieldEvaluatorReserved) {
				b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorResult)
			}
			if b51AMirrorFieldEvaluatorIndex == 0 {
				if _, b51AMirrorFieldEvaluatorHas := b51AMirrorFieldEvaluatorResult["reserved"]; b51AMirrorFieldEvaluatorHas {
					b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorResult)
				}
			}
		}
	}
}
func TestAb19Ap19ContextAndOwnership(b51AMirrorFieldEvaluatorT *b51AMirrorFieldEvaluatorTesting.T) {
	for _, b51AMirrorFieldEvaluatorKey := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant", "opaque"} {
		for _, b51AMirrorFieldEvaluatorValue := range []any{nil, false, float64(0), "other", map[string]any{"nested": []any{false, nil, float64(0)}}} {
			b51AMirrorFieldEvaluatorInput := b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT)[0]["input"].(map[string]any)
			b51AMirrorFieldEvaluatorContext := b51AMirrorFieldEvaluatorInput["context"].(map[string]any)
			if b51AMirrorFieldEvaluatorKey == "primary" || b51AMirrorFieldEvaluatorKey == "secondary" || b51AMirrorFieldEvaluatorKey == "operation" {
				if b51AMirrorFieldEvaluatorValue != nil && b51AMirrorFieldEvaluatorValue != float64(0) {
					continue
				}
			}
			if b51AMirrorFieldEvaluatorKey == "direction" || b51AMirrorFieldEvaluatorKey == "view" {
				if b51AMirrorFieldEvaluatorValue != nil && b51AMirrorFieldEvaluatorValue != "other" {
					continue
				}
			}
			b51AMirrorFieldEvaluatorContext[b51AMirrorFieldEvaluatorKey] = b51AMirrorFieldEvaluatorValue
			b51AMirrorFieldEvaluatorInput["arguments"] = map[string]any{"opaque": []any{false, nil, float64(0), "", []any{}}}
			b51AMirrorFieldEvaluatorOutput := b51AMirrorFieldEvaluatorRun(b51AMirrorFieldEvaluatorT, b51AMirrorFieldEvaluatorInput)
			if !b51AMirrorFieldEvaluatorReflect.DeepEqual(b51AMirrorFieldEvaluatorOutput["context"], b51AMirrorFieldEvaluatorContext) || !b51AMirrorFieldEvaluatorReflect.DeepEqual(b51AMirrorFieldEvaluatorOutput["arguments"], b51AMirrorFieldEvaluatorInput["arguments"]) || b51AMirrorFieldEvaluatorOutput["qualified"] != false {
				b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorOutput)
			}
			b51AMirrorFieldEvaluatorUnsupported := b51AMirrorFieldEvaluatorKey != "family" && b51AMirrorFieldEvaluatorKey != "variant" && b51AMirrorFieldEvaluatorKey != "opaque"
			if b51AMirrorFieldEvaluatorUnsupported {
				if b51AMirrorFieldEvaluatorOutput["state"] != "Unknown" || len(b51AMirrorFieldEvaluatorOutput["fields"].(map[string]any)) != 0 {
					b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorOutput)
				}
			} else if b51AMirrorFieldEvaluatorOutput["state"] != "complete" {
				b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorOutput)
			}
			delete(b51AMirrorFieldEvaluatorContext, b51AMirrorFieldEvaluatorKey)
			b51AMirrorFieldEvaluatorOutput = b51AMirrorFieldEvaluatorRun(b51AMirrorFieldEvaluatorT, b51AMirrorFieldEvaluatorInput)
			if _, b51AMirrorFieldEvaluatorPresent := b51AMirrorFieldEvaluatorOutput["context"].(map[string]any)[b51AMirrorFieldEvaluatorKey]; b51AMirrorFieldEvaluatorPresent {
				b51AMirrorFieldEvaluatorT.Fatal("missing context invented")
			}
		}
	}
	b51AMirrorFieldEvaluatorInput := b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT)[0]["input"].(map[string]any)
	b51AMirrorFieldEvaluatorInput["context"].(map[string]any)["register"] = float64(6)
	b51AMirrorFieldEvaluatorInput["arguments"] = map[string]any{"register": float64(6), "sequence": float64(255), "length": float64(3)}
	b51AMirrorFieldEvaluatorOutput := b51AMirrorFieldEvaluatorRun(b51AMirrorFieldEvaluatorT, b51AMirrorFieldEvaluatorInput)
	if b51AMirrorFieldEvaluatorOutput["state"] != "complete" || len(b51AMirrorFieldEvaluatorOutput["fields"].(map[string]any)) != 4 || b51AMirrorFieldEvaluatorField(b51AMirrorFieldEvaluatorOutput, "sequence")["number"] != float64(0) {
		b51AMirrorFieldEvaluatorT.Fatal("selector alias or sequence correlation", b51AMirrorFieldEvaluatorOutput)
	}
	b51AMirrorFieldEvaluatorBytes, _ := b51AMirrorFieldEvaluatorJSON.Marshal(b51AMirrorFieldEvaluatorInput)
	b51AMirrorFieldEvaluatorResult, b51AMirrorFieldEvaluatorErr := B51AMirrorFieldEvaluate(b51AMirrorFieldEvaluatorBytes)
	if b51AMirrorFieldEvaluatorErr != nil {
		b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorErr)
	}
	b51AMirrorFieldEvaluatorSaved := append([]byte(nil), b51AMirrorFieldEvaluatorResult...)
	for b51AMirrorFieldEvaluatorIndex := range b51AMirrorFieldEvaluatorBytes {
		b51AMirrorFieldEvaluatorBytes[b51AMirrorFieldEvaluatorIndex] = ' '
	}
	if string(b51AMirrorFieldEvaluatorSaved) != string(b51AMirrorFieldEvaluatorResult) {
		b51AMirrorFieldEvaluatorT.Fatal("input aliases result")
	}
	for b51AMirrorFieldEvaluatorIndex := range b51AMirrorFieldEvaluatorResult {
		b51AMirrorFieldEvaluatorResult[b51AMirrorFieldEvaluatorIndex] = ' '
	}
	b51AMirrorFieldEvaluatorAgain, _ := b51AMirrorFieldEvaluatorJSON.Marshal(b51AMirrorFieldEvaluatorInput)
	b51AMirrorFieldEvaluatorResult, b51AMirrorFieldEvaluatorErr = B51AMirrorFieldEvaluate(b51AMirrorFieldEvaluatorAgain)
	if b51AMirrorFieldEvaluatorErr != nil || string(b51AMirrorFieldEvaluatorResult) != string(b51AMirrorFieldEvaluatorSaved) {
		b51AMirrorFieldEvaluatorT.Fatal("result aliases later call")
	}
	b51AMirrorFieldEvaluatorOutput["raw"].(map[string]any)["sequence"].([]any)[0] = float64(99)
	if b51AMirrorFieldEvaluatorField(b51AMirrorFieldEvaluatorOutput, "sequence")["raw"].([]any)[0] != float64(0) || b51AMirrorFieldEvaluatorRun(b51AMirrorFieldEvaluatorT, b51AMirrorFieldEvaluatorInput)["raw"].(map[string]any)["sequence"].([]any)[0] != float64(0) {
		b51AMirrorFieldEvaluatorT.Fatal("nested storage aliases")
	}
}
func TestAb19Ap19Reject(b51AMirrorFieldEvaluatorT *b51AMirrorFieldEvaluatorTesting.T) {
	for _, b51AMirrorFieldEvaluatorInput := range []string{``, `{`, `[]`, `{} {}`, `{"context":{},"fields":{},"fields":{}}`, `{"context":{"x":{"a":1,"a":2}},"fields":{}}`, `{"context":{},"fields":{},"arguments":{"x":[{"a":0,"a":1}]}}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":null}`, `{"context":{},"fields":{},"arguments":false}`} {
		b51AMirrorFieldEvaluatorOutput, b51AMirrorFieldEvaluatorErr := B51AMirrorFieldEvaluate([]byte(b51AMirrorFieldEvaluatorInput))
		if b51AMirrorFieldEvaluatorErr == nil || b51AMirrorFieldEvaluatorOutput != nil {
			b51AMirrorFieldEvaluatorT.Fatal("accepted invalid JSON", b51AMirrorFieldEvaluatorInput)
		}
	}
	for _, b51AMirrorFieldEvaluatorInvalid := range []any{nil, false, "0", float64(-1), float64(256), float64(0.5), b51AMirrorFieldEvaluatorJSON.Number("0.00000000000000000001"), b51AMirrorFieldEvaluatorJSON.Number("255.00000000000000000001"), map[string]any{}} {
		for _, b51AMirrorFieldEvaluatorKey := range []string{"sequence", "flags", "length", "address", "mirror_address", "symbol_address", "status", "data"} {
			b51AMirrorFieldEvaluatorInput := b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT)[0]["input"].(map[string]any)
			b51AMirrorFieldEvaluatorInput["fields"].(map[string]any)[b51AMirrorFieldEvaluatorKey] = []any{b51AMirrorFieldEvaluatorInvalid}
			b51AMirrorFieldEvaluatorBytes, _ := b51AMirrorFieldEvaluatorJSON.Marshal(b51AMirrorFieldEvaluatorInput)
			b51AMirrorFieldEvaluatorOutput, b51AMirrorFieldEvaluatorErr := B51AMirrorFieldEvaluate(b51AMirrorFieldEvaluatorBytes)
			if b51AMirrorFieldEvaluatorErr == nil || b51AMirrorFieldEvaluatorOutput != nil {
				b51AMirrorFieldEvaluatorT.Fatal("accepted invalid byte", b51AMirrorFieldEvaluatorKey, b51AMirrorFieldEvaluatorInvalid)
			}
		}
	}
	for _, b51AMirrorFieldEvaluatorKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		b51AMirrorFieldEvaluatorInput := b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT)[0]["input"].(map[string]any)
		b51AMirrorFieldEvaluatorInput["context"].(map[string]any)[b51AMirrorFieldEvaluatorKey] = []int{}
		b51AMirrorFieldEvaluatorBytes, _ := b51AMirrorFieldEvaluatorJSON.Marshal(b51AMirrorFieldEvaluatorInput)
		if _, b51AMirrorFieldEvaluatorErr := B51AMirrorFieldEvaluate(b51AMirrorFieldEvaluatorBytes); b51AMirrorFieldEvaluatorErr == nil {
			b51AMirrorFieldEvaluatorT.Fatal("accepted wrong context kind")
		}
	}
}

func TestAb19Ap19IndependentCoordinates(b51AMirrorFieldEvaluatorT *b51AMirrorFieldEvaluatorTesting.T) {
	for b51AMirrorFieldEvaluatorIndex := 0; b51AMirrorFieldEvaluatorIndex < 4; b51AMirrorFieldEvaluatorIndex++ {
		for _, b51AMirrorFieldEvaluatorSequence := range []float64{0, 255} {
			b51AMirrorFieldEvaluatorInput := b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT)[b51AMirrorFieldEvaluatorIndex]["input"].(map[string]any)
			b51AMirrorFieldEvaluatorInput["fields"].(map[string]any)["sequence"] = []float64{b51AMirrorFieldEvaluatorSequence}
			b51AMirrorFieldEvaluatorInput["arguments"] = map[string]any{"sequence": 255 - b51AMirrorFieldEvaluatorSequence, "length": 3, "register": 6}
			b51AMirrorFieldEvaluatorOutput := b51AMirrorFieldEvaluatorRun(b51AMirrorFieldEvaluatorT, b51AMirrorFieldEvaluatorInput)
			if b51AMirrorFieldEvaluatorOutput["state"] != "complete" || b51AMirrorFieldEvaluatorField(b51AMirrorFieldEvaluatorOutput, "sequence")["number"] != b51AMirrorFieldEvaluatorSequence || !b51AMirrorFieldEvaluatorReflect.DeepEqual(b51AMirrorFieldEvaluatorOutput["context"], b51AMirrorFieldEvaluatorInput["context"]) {
				b51AMirrorFieldEvaluatorT.Fatal("sequence correlation", b51AMirrorFieldEvaluatorOutput)
			}
			if b51AMirrorFieldEvaluatorIndex == 1 && len(b51AMirrorFieldEvaluatorField(b51AMirrorFieldEvaluatorOutput, "data")["raw"].([]any)) != 2 {
				b51AMirrorFieldEvaluatorT.Fatal("requested length resized data")
			}
		}
	}
	for _, b51AMirrorFieldEvaluatorName := range []string{"address", "mirror_address", "symbol_address"} {
		for _, b51AMirrorFieldEvaluatorPair := range []struct {
			b51AMirrorFieldEvaluatorBytes []int
			b51AMirrorFieldEvaluatorWant  float64
		}{{[]int{0, 0}, 0}, {[]int{255, 255}, 65535}, {[]int{52, 18}, 4660}, {[]int{120, 86}, 22136}} {
			b51AMirrorFieldEvaluatorIndex := 2
			if b51AMirrorFieldEvaluatorName == "address" {
				b51AMirrorFieldEvaluatorIndex = 0
			}
			b51AMirrorFieldEvaluatorInput := b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT)[b51AMirrorFieldEvaluatorIndex]["input"].(map[string]any)
			b51AMirrorFieldEvaluatorInput["fields"].(map[string]any)[b51AMirrorFieldEvaluatorName] = b51AMirrorFieldEvaluatorPair.b51AMirrorFieldEvaluatorBytes
			b51AMirrorFieldEvaluatorResult := b51AMirrorFieldEvaluatorField(b51AMirrorFieldEvaluatorRun(b51AMirrorFieldEvaluatorT, b51AMirrorFieldEvaluatorInput), b51AMirrorFieldEvaluatorName)
			if b51AMirrorFieldEvaluatorResult["number"] != b51AMirrorFieldEvaluatorPair.b51AMirrorFieldEvaluatorWant || b51AMirrorFieldEvaluatorResult["logical_validity"] != "Unknown" || b51AMirrorFieldEvaluatorResult["label"] != nil {
				b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorResult)
			}
		}
	}
	for _, b51AMirrorFieldEvaluatorArgument := range []any{map[string]any{}, map[string]any{"opaque": false}, map[string]any{"opaque": 0}, map[string]any{"opaque": []any{}}, map[string]any{"opaque": nil}} {
		b51AMirrorFieldEvaluatorInput := b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT)[0]["input"].(map[string]any)
		b51AMirrorFieldEvaluatorInput["arguments"] = b51AMirrorFieldEvaluatorArgument
		b51AMirrorFieldEvaluatorOutput := b51AMirrorFieldEvaluatorRun(b51AMirrorFieldEvaluatorT, b51AMirrorFieldEvaluatorInput)
		b51AMirrorFieldEvaluatorBytes, _ := b51AMirrorFieldEvaluatorJSON.Marshal(b51AMirrorFieldEvaluatorArgument)
		var b51AMirrorFieldEvaluatorWant any
		_ = b51AMirrorFieldEvaluatorJSON.Unmarshal(b51AMirrorFieldEvaluatorBytes, &b51AMirrorFieldEvaluatorWant)
		if !b51AMirrorFieldEvaluatorReflect.DeepEqual(b51AMirrorFieldEvaluatorWant, b51AMirrorFieldEvaluatorOutput["arguments"]) || b51AMirrorFieldEvaluatorOutput["state"] != "complete" {
			b51AMirrorFieldEvaluatorT.Fatal(b51AMirrorFieldEvaluatorOutput)
		}
	}
}

func TestAb19Ap19RawBoundaries(b51AMirrorFieldEvaluatorT *b51AMirrorFieldEvaluatorTesting.T) {
	for _, b51AMirrorFieldEvaluatorName := range []string{"sequence", "flags", "length", "address", "mirror_address", "symbol_address", "status", "data"} {
		for _, b51AMirrorFieldEvaluatorValue := range []any{nil, false, 0, "", map[string]any{}} {
			b51AMirrorFieldEvaluatorInput := b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT)[0]["input"].(map[string]any)
			b51AMirrorFieldEvaluatorInput["fields"].(map[string]any)[b51AMirrorFieldEvaluatorName] = b51AMirrorFieldEvaluatorValue
			b51AMirrorFieldEvaluatorBytes, _ := b51AMirrorFieldEvaluatorJSON.Marshal(b51AMirrorFieldEvaluatorInput)
			if b51AMirrorFieldEvaluatorOutput, b51AMirrorFieldEvaluatorErr := B51AMirrorFieldEvaluate(b51AMirrorFieldEvaluatorBytes); b51AMirrorFieldEvaluatorErr == nil || b51AMirrorFieldEvaluatorOutput != nil {
				b51AMirrorFieldEvaluatorT.Fatal("field must be array", b51AMirrorFieldEvaluatorName)
			}
		}
	}
	for _, b51AMirrorFieldEvaluatorValue := range []any{[]any{}, []int{0}, []int{255}} {
		b51AMirrorFieldEvaluatorInput := b51AMirrorFieldEvaluatorCases(b51AMirrorFieldEvaluatorT)[0]["input"].(map[string]any)
		b51AMirrorFieldEvaluatorInput["fields"].(map[string]any)["configured_register"] = b51AMirrorFieldEvaluatorValue
		b51AMirrorFieldEvaluatorOutput := b51AMirrorFieldEvaluatorRun(b51AMirrorFieldEvaluatorT, b51AMirrorFieldEvaluatorInput)
		if len(b51AMirrorFieldEvaluatorOutput["fields"].(map[string]any)) != 4 || b51AMirrorFieldEvaluatorOutput["state"] != "complete" {
			b51AMirrorFieldEvaluatorT.Fatal("unknown field interpreted", b51AMirrorFieldEvaluatorOutput)
		}
		b51AMirrorFieldEvaluatorBytes, _ := b51AMirrorFieldEvaluatorJSON.Marshal(b51AMirrorFieldEvaluatorValue)
		var b51AMirrorFieldEvaluatorWant any
		_ = b51AMirrorFieldEvaluatorJSON.Unmarshal(b51AMirrorFieldEvaluatorBytes, &b51AMirrorFieldEvaluatorWant)
		if !b51AMirrorFieldEvaluatorReflect.DeepEqual(b51AMirrorFieldEvaluatorOutput["raw"].(map[string]any)["configured_register"], b51AMirrorFieldEvaluatorWant) {
			b51AMirrorFieldEvaluatorT.Fatal("unknown raw lost")
		}
	}
}
