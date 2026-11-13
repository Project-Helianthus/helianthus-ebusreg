package vaillant

import (
	b521FieldViewBytes "bytes"
	b521FieldViewJSON "encoding/json"
	b521FieldViewOs "os"
	b521FieldViewReflect "reflect"
	b521FieldViewStrconv "strconv"
	b521FieldViewStrings "strings"
	b521FieldViewTesting "testing"
)

func b521FieldViewInput(b521FieldViewT *b521FieldViewTesting.T, b521FieldViewDirection string) map[string]any {
	b521FieldViewT.Helper()
	b521FieldViewFields := map[string]any{"status0": []any{0}, "status1": []any{0}, "brine_percent": []any{0}, "fan_percent": []any{0}, "fan_speed": []any{0}, "temperature0": []any{0, 0}, "temperature1": []any{0, 0}}
	if b521FieldViewDirection == "request" {
		b521FieldViewFields = map[string]any{"flags": []any{0}, "appliance_bits": []any{0}, "noise_percent": []any{0}, "frost_temperature": []any{0}}
	}
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 33, "view": "module-fields", "direction": b521FieldViewDirection}, "fields": b521FieldViewFields}
}

func b521FieldViewEvaluate(b521FieldViewT *b521FieldViewTesting.T, b521FieldViewInput any) map[string]any {
	b521FieldViewT.Helper()
	b521FieldViewEncoded, b521FieldViewError := b521FieldViewJSON.Marshal(b521FieldViewInput)
	if b521FieldViewError != nil {
		b521FieldViewT.Fatal(b521FieldViewError)
	}
	b521FieldViewRaw, b521FieldViewError := B521FieldViewEvaluate(b521FieldViewEncoded)
	if b521FieldViewError != nil {
		b521FieldViewT.Fatal(b521FieldViewError)
	}
	b521FieldViewOutput := map[string]any{}
	if b521FieldViewError := b521FieldViewJSON.Unmarshal(b521FieldViewRaw, &b521FieldViewOutput); b521FieldViewError != nil {
		b521FieldViewT.Fatal(b521FieldViewError)
	}
	if b521FieldViewOutput["physical"] != "Unknown" || b521FieldViewOutput["qualified"] != false {
		b521FieldViewT.Fatal("qualification changed", b521FieldViewOutput)
	}
	return b521FieldViewOutput
}

func b521FieldViewResults(b521FieldViewOutput map[string]any) map[string]any {
	return b521FieldViewOutput["results"].(map[string]any)
}

func TestAb13Aq13Fixture(b521FieldViewT *b521FieldViewTesting.T) {
	b521FieldViewRaw, b521FieldViewError := b521FieldViewOs.ReadFile("testdata/ab13aq13_b521_field_views.json")
	if b521FieldViewError != nil {
		b521FieldViewT.Fatal(b521FieldViewError)
	}
	var b521FieldViewCases []struct {
		B521FieldViewID       string                    `json:"id"`
		B521FieldViewInput    map[string]any            `json:"input"`
		B521FieldViewState    string                    `json:"state"`
		B521FieldViewExpected map[string]map[string]any `json:"expected"`
	}
	if b521FieldViewError := b521FieldViewJSON.Unmarshal(b521FieldViewRaw, &b521FieldViewCases); b521FieldViewError != nil {
		b521FieldViewT.Fatal(b521FieldViewError)
	}
	if len(b521FieldViewCases) != 41 {
		b521FieldViewT.Fatalf("fixture count %d", len(b521FieldViewCases))
	}
	for _, b521FieldViewCase := range b521FieldViewCases {
		b521FieldViewT.Run(b521FieldViewCase.B521FieldViewID, func(b521FieldViewT *b521FieldViewTesting.T) {
			b521FieldViewOutput := b521FieldViewEvaluate(b521FieldViewT, b521FieldViewCase.B521FieldViewInput)
			if b521FieldViewOutput["state"] != b521FieldViewCase.B521FieldViewState {
				b521FieldViewT.Fatal(b521FieldViewOutput)
			}
			if !b521FieldViewReflect.DeepEqual(b521FieldViewOutput["fields"], b521FieldViewCase.B521FieldViewInput["fields"]) {
				b521FieldViewT.Fatal("raw fields changed")
			}
			for b521FieldViewName, b521FieldViewExpected := range b521FieldViewCase.B521FieldViewExpected {
				b521FieldViewResult := b521FieldViewResults(b521FieldViewOutput)[b521FieldViewName].(map[string]any)
				for b521FieldViewKey, b521FieldViewValue := range b521FieldViewExpected {
					if !b521FieldViewReflect.DeepEqual(b521FieldViewResult[b521FieldViewKey], b521FieldViewValue) {
						b521FieldViewT.Fatalf("%s.%s got %v want %v", b521FieldViewName, b521FieldViewKey, b521FieldViewResult[b521FieldViewKey], b521FieldViewValue)
					}
				}
			}
			for _, b521FieldViewValue := range b521FieldViewResults(b521FieldViewOutput) {
				b521FieldViewResult := b521FieldViewValue.(map[string]any)
				if b521FieldViewResult["grade"] != "Hypothesis" {
					b521FieldViewT.Fatal("candidate grade lost")
				}
				if b521FieldViewResult["state"] == "complete" && b521FieldViewResult["calculation"] != "Proven" {
					b521FieldViewT.Fatal("calculation grade lost")
				}
				for b521FieldViewKey := range b521FieldViewResult {
					switch b521FieldViewKey {
					case "raw", "present", "state", "grade", "calculation", "flags", "reserved", "number", "unit":
					default:
						b521FieldViewT.Fatalf("invented result property %s", b521FieldViewKey)
					}
				}
			}
		})
	}
}

func TestAb13Aq13Context(b521FieldViewT *b521FieldViewTesting.T) {
	for b521FieldViewKey, b521FieldViewValues := range map[string][]any{"primary": {0, 180, 182, 255, nil}, "secondary": {0, 32, 34, 255, nil}, "operation": {0, 1, 255, nil}, "direction": {"negative", "local", "", nil}, "view": {"other", "", nil}} {
		for b521FieldViewIndex, b521FieldViewValue := range b521FieldViewValues {
			b521FieldViewT.Run(b521FieldViewKey+b521FieldViewStrconv.Itoa(b521FieldViewIndex), func(b521FieldViewT *b521FieldViewTesting.T) {
				b521FieldViewInput := b521FieldViewInput(b521FieldViewT, "positive")
				b521FieldViewInput["context"].(map[string]any)[b521FieldViewKey] = b521FieldViewValue
				b521FieldViewOutput := b521FieldViewEvaluate(b521FieldViewT, b521FieldViewInput)
				b521FieldViewContextBefore, _ := b521FieldViewJSON.Marshal(b521FieldViewInput["context"])
				b521FieldViewContextAfter, _ := b521FieldViewJSON.Marshal(b521FieldViewOutput["context"])
				if !b521FieldViewBytes.Equal(b521FieldViewContextBefore, b521FieldViewContextAfter) {
					b521FieldViewT.Fatal("unsupported context lost")
				}
				if b521FieldViewOutput["state"] != "Unknown" {
					b521FieldViewT.Fatal("unsupported context evaluated")
				}
				for _, b521FieldViewValue := range b521FieldViewResults(b521FieldViewOutput) {
					b521FieldViewResult := b521FieldViewValue.(map[string]any)
					if len(b521FieldViewResult) != 3 || b521FieldViewResult["state"] != "Unknown" {
						b521FieldViewT.Fatal("unsupported result", b521FieldViewResult)
					}
				}
			})
		}
	}
	for _, b521FieldViewKey := range []string{"primary", "secondary", "direction", "view"} {
		b521FieldViewInput := b521FieldViewInput(b521FieldViewT, "positive")
		delete(b521FieldViewInput["context"].(map[string]any), b521FieldViewKey)
		if b521FieldViewEvaluate(b521FieldViewT, b521FieldViewInput)["state"] != "Unknown" {
			b521FieldViewT.Fatal("missing context evaluated", b521FieldViewKey)
		}
	}
	for _, b521FieldViewKey := range []string{"family", "variant"} {
		for _, b521FieldViewValue := range []any{nil, false, 0, "", "synthetic-family", map[string]any{"opaque": true}} {
			b521FieldViewInput := b521FieldViewInput(b521FieldViewT, "positive")
			b521FieldViewInput["context"].(map[string]any)[b521FieldViewKey] = b521FieldViewValue
			if b521FieldViewEvaluate(b521FieldViewT, b521FieldViewInput)["state"] != "complete" {
				b521FieldViewT.Fatal("opaque identity selected view")
			}
		}
	}
}

func TestAb13Aq13WidthsAndSiblings(b521FieldViewT *b521FieldViewTesting.T) {
	for _, b521FieldViewDirection := range []string{"request", "positive"} {
		for b521FieldViewName, b521FieldViewOriginal := range b521FieldViewInput(b521FieldViewT, b521FieldViewDirection)["fields"].(map[string]any) {
			for _, b521FieldViewLength := range []int{-1, 0, 1, 2, 3} {
				b521FieldViewT.Run(b521FieldViewDirection+b521FieldViewName+b521FieldViewStrconv.Itoa(b521FieldViewLength), func(b521FieldViewT *b521FieldViewTesting.T) {
					b521FieldViewInput := b521FieldViewInput(b521FieldViewT, b521FieldViewDirection)
					b521FieldViewFields := b521FieldViewInput["fields"].(map[string]any)
					delete(b521FieldViewFields, b521FieldViewName)
					if b521FieldViewLength >= 0 {
						b521FieldViewArray := []any{}
						for b521FieldViewIndex := 0; b521FieldViewIndex < b521FieldViewLength; b521FieldViewIndex++ {
							b521FieldViewArray = append(b521FieldViewArray, 0)
						}
						b521FieldViewFields[b521FieldViewName] = b521FieldViewArray
					}
					b521FieldViewFields["opaque"] = []any{}
					b521FieldViewBaseline := b521FieldViewResults(b521FieldViewEvaluate(b521FieldViewT, b521FieldViewInputWithoutDirectionChange(b521FieldViewT, b521FieldViewDirection)))
					b521FieldViewOutput := b521FieldViewEvaluate(b521FieldViewT, b521FieldViewInput)
					b521FieldViewComplete := b521FieldViewLength == len(b521FieldViewOriginal.([]any))
					b521FieldViewExpected := "partial"
					if b521FieldViewComplete {
						b521FieldViewExpected = "complete"
					}
					if b521FieldViewOutput["state"] != b521FieldViewExpected {
						b521FieldViewT.Fatal(b521FieldViewOutput)
					}
					for b521FieldViewOther, b521FieldViewValue := range b521FieldViewResults(b521FieldViewOutput) {
						b521FieldViewResult := b521FieldViewValue.(map[string]any)
						if b521FieldViewOther == "opaque" {
							if b521FieldViewResult["state"] != "Unknown" || !b521FieldViewReflect.DeepEqual(b521FieldViewResult["raw"], []any{}) {
								b521FieldViewT.Fatal("empty opaque lost")
							}
							continue
						}
						if b521FieldViewOther == b521FieldViewName && !b521FieldViewComplete {
							if b521FieldViewResult["state"] != "incomplete" || b521FieldViewResult["present"] != (b521FieldViewLength >= 0) {
								b521FieldViewT.Fatal("missing/width confused")
							}
							if _, b521FieldViewExists := b521FieldViewResult["number"]; b521FieldViewExists {
								b521FieldViewT.Fatal("invented number")
							}
						} else if !b521FieldViewReflect.DeepEqual(b521FieldViewResult, b521FieldViewBaseline[b521FieldViewOther]) {
							b521FieldViewT.Fatal("sibling changed", b521FieldViewOther)
						}
					}
				})
			}
		}
	}
}

func TestAb13Aq13Reject(b521FieldViewT *b521FieldViewTesting.T) {
	for b521FieldViewIndex, b521FieldViewRaw := range []string{``, `{`, `[]`, `null`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"primary":181,"primary":181},"fields":{}}`, `{"context":{},"fields":{"x":[0],"x":[1]}}`, `{"context":{},"fields":{},"extra":[{"x":0,"x":1}]}`, `{"context":{},"fields":{},"extra":1e999}`, `{"context":{},"fields":{},"arguments":{"x":1,"x":2}}`} {
		b521FieldViewT.Run("syntax"+b521FieldViewStrconv.Itoa(b521FieldViewIndex), func(b521FieldViewT *b521FieldViewTesting.T) {
			b521FieldViewOutput, b521FieldViewError := B521FieldViewEvaluate([]byte(b521FieldViewRaw))
			if b521FieldViewError == nil || b521FieldViewOutput != nil {
				b521FieldViewT.Fatal("invalid JSON accepted")
			}
		})
	}
	for _, b521FieldViewKey := range []string{"primary", "secondary", "operation", "direction", "view", "fields", "context", "arguments", "field"} {
		for b521FieldViewIndex, b521FieldViewBad := range []any{nil, true, "1", -1, 256, 1.5, []any{}, map[string]any{}} {
			if b521FieldViewBad == nil && (b521FieldViewKey == "primary" || b521FieldViewKey == "secondary" || b521FieldViewKey == "operation" || b521FieldViewKey == "direction" || b521FieldViewKey == "view") {
				continue
			}
			if (b521FieldViewKey == "direction" || b521FieldViewKey == "view") && b521FieldViewIndex == 2 {
				continue
			}
			if (b521FieldViewKey == "fields" || b521FieldViewKey == "context" || b521FieldViewKey == "arguments") && b521FieldViewIndex == 7 {
				continue
			}
			b521FieldViewT.Run(b521FieldViewKey+b521FieldViewStrconv.Itoa(b521FieldViewIndex), func(b521FieldViewT *b521FieldViewTesting.T) {
				b521FieldViewInput := b521FieldViewInput(b521FieldViewT, "positive")
				switch b521FieldViewKey {
				case "fields", "context", "arguments":
					b521FieldViewInput[b521FieldViewKey] = b521FieldViewBad
				case "field":
					b521FieldViewInput["fields"].(map[string]any)["status0"] = []any{b521FieldViewBad}
				default:
					b521FieldViewInput["context"].(map[string]any)[b521FieldViewKey] = b521FieldViewBad
				}
				b521FieldViewRaw, _ := b521FieldViewJSON.Marshal(b521FieldViewInput)
				b521FieldViewOutput, b521FieldViewError := B521FieldViewEvaluate(b521FieldViewRaw)
				if b521FieldViewError == nil || b521FieldViewOutput != nil {
					b521FieldViewT.Fatal("wrong kind accepted", string(b521FieldViewRaw))
				}
			})
		}
	}
	for _, b521FieldViewBad := range []any{nil, true, "0", 0, 1.5, []any{}} {
		b521FieldViewInput := b521FieldViewInput(b521FieldViewT, "positive")
		b521FieldViewInput["fields"].(map[string]any)["status0"] = b521FieldViewBad
		if b521FieldViewArray, b521FieldViewOK := b521FieldViewBad.([]any); b521FieldViewOK && len(b521FieldViewArray) == 0 {
			continue
		}
		b521FieldViewRaw, _ := b521FieldViewJSON.Marshal(b521FieldViewInput)
		if b521FieldViewOutput, b521FieldViewError := B521FieldViewEvaluate(b521FieldViewRaw); b521FieldViewError == nil || b521FieldViewOutput != nil {
			b521FieldViewT.Fatal("nonarray accepted")
		}
	}
	for _, b521FieldViewValue := range []any{nil, false, 0, "", 1.5} {
		b521FieldViewInput := b521FieldViewInput(b521FieldViewT, "positive")
		b521FieldViewInput["arguments"] = map[string]any{"undeclared": b521FieldViewValue}
		b521FieldViewRaw, _ := b521FieldViewJSON.Marshal(b521FieldViewInput)
		if b521FieldViewOutput, b521FieldViewError := B521FieldViewEvaluate(b521FieldViewRaw); b521FieldViewError == nil || b521FieldViewOutput != nil {
			b521FieldViewT.Fatal("undeclared argument accepted")
		}
	}
}

func TestAb13Aq13RetentionAndMemory(b521FieldViewT *b521FieldViewTesting.T) {
	b521FieldViewInput := b521FieldViewInput(b521FieldViewT, "positive")
	b521FieldViewInput["extra"] = map[string]any{"absent-check": nil, "null": nil, "zero": 0, "false": false, "empty": []any{}, "nested": map[string]any{"value": "synthetic"}}
	b521FieldViewInput["arguments"] = map[string]any{}
	b521FieldViewRaw, _ := b521FieldViewJSON.Marshal(b521FieldViewInput)
	b521FieldViewBefore := append([]byte(nil), b521FieldViewRaw...)
	b521FieldViewFirst, b521FieldViewError := B521FieldViewEvaluate(b521FieldViewRaw)
	if b521FieldViewError != nil {
		b521FieldViewT.Fatal(b521FieldViewError)
	}
	if !b521FieldViewBytes.Equal(b521FieldViewRaw, b521FieldViewBefore) {
		b521FieldViewT.Fatal("input mutated")
	}
	b521FieldViewOutput := b521FieldViewEvaluate(b521FieldViewT, b521FieldViewInput)
	b521FieldViewNormalized := map[string]any{}
	_ = b521FieldViewJSON.Unmarshal(b521FieldViewRaw, &b521FieldViewNormalized)
	for _, b521FieldViewKey := range []string{"context", "fields", "arguments", "extra"} {
		if !b521FieldViewReflect.DeepEqual(b521FieldViewOutput[b521FieldViewKey], b521FieldViewNormalized[b521FieldViewKey]) {
			b521FieldViewT.Fatal("retention", b521FieldViewKey)
		}
	}
	b521FieldViewSecond, _ := B521FieldViewEvaluate(b521FieldViewBefore)
	b521FieldViewFirst[0] = '!'
	b521FieldViewRaw[0] = '!'
	b521FieldViewThird, b521FieldViewError := B521FieldViewEvaluate(b521FieldViewBefore)
	if b521FieldViewError != nil || !b521FieldViewBytes.Equal(b521FieldViewSecond, b521FieldViewThird) {
		b521FieldViewT.Fatal("shared storage")
	}
	b521FieldViewNoArgs := b521FieldViewEvaluate(b521FieldViewT, b521FieldViewInputWithoutArguments(b521FieldViewT))
	if _, b521FieldViewPresent := b521FieldViewNoArgs["arguments"]; b521FieldViewPresent {
		b521FieldViewT.Fatal("invented empty arguments")
	}
	if b521FieldViewStrings.Contains(string(b521FieldViewThird), "offset") || b521FieldViewStrings.Contains(string(b521FieldViewThird), "position") {
		b521FieldViewT.Fatal("invented coordinate")
	}
}

func b521FieldViewInputWithoutArguments(b521FieldViewT *b521FieldViewTesting.T) map[string]any {
	return b521FieldViewInput(b521FieldViewT, "positive")
}

func TestAb13Aq13SignedDomain(b521FieldViewT *b521FieldViewTesting.T) {
	for b521FieldViewCarrier := 0; b521FieldViewCarrier < 65536; b521FieldViewCarrier++ {
		b521FieldViewInput := b521FieldViewInput(b521FieldViewT, "positive")
		b521FieldViewInput["fields"].(map[string]any)["temperature0"] = []any{b521FieldViewCarrier % 256, b521FieldViewCarrier / 256}
		b521FieldViewInput["fields"].(map[string]any)["temperature1"] = []any{255 - b521FieldViewCarrier%256, 255 - b521FieldViewCarrier/256}
		b521FieldViewExpected := float64(b521FieldViewCarrier) / 16
		if b521FieldViewCarrier >= 32768 {
			b521FieldViewExpected -= 4096
		}
		b521FieldViewAllResults := b521FieldViewResults(b521FieldViewEvaluate(b521FieldViewT, b521FieldViewInput))
		b521FieldViewResult := b521FieldViewAllResults["temperature0"].(map[string]any)
		b521FieldViewComplementExpected := float64(65535-b521FieldViewCarrier) / 16
		if b521FieldViewCarrier < 32768 {
			b521FieldViewComplementExpected -= 4096
		}
		if b521FieldViewAllResults["temperature1"].(map[string]any)["number"] != b521FieldViewComplementExpected {
			b521FieldViewT.Fatal("temperature1 carrier", b521FieldViewCarrier)
		}
		if b521FieldViewResult["number"] != b521FieldViewExpected {
			b521FieldViewT.Fatalf("carrier %d: %v want %v", b521FieldViewCarrier, b521FieldViewResult["number"], b521FieldViewExpected)
		}
	}
}

func b521FieldViewInputWithoutDirectionChange(b521FieldViewT *b521FieldViewTesting.T, b521FieldViewDirection string) map[string]any {
	return b521FieldViewInput(b521FieldViewT, b521FieldViewDirection)
}

func TestAb13Aq13OpaqueExtent(b521FieldViewT *b521FieldViewTesting.T) {
	for b521FieldViewDirection, b521FieldViewLength := range map[string]int{"request": 7, "positive": 12} {
		b521FieldViewInput := b521FieldViewInput(b521FieldViewT, b521FieldViewDirection)
		b521FieldViewPayload := []any{}
		for b521FieldViewIndex := 0; b521FieldViewIndex < b521FieldViewLength; b521FieldViewIndex++ {
			b521FieldViewPayload = append(b521FieldViewPayload, 255)
		}
		b521FieldViewInput["fields"] = map[string]any{"opaque": b521FieldViewPayload}
		b521FieldViewOutput := b521FieldViewEvaluate(b521FieldViewT, b521FieldViewInput)
		if b521FieldViewOutput["state"] != "partial" {
			b521FieldViewT.Fatal("abstract extent interpreted")
		}
		for b521FieldViewName, b521FieldViewValue := range b521FieldViewResults(b521FieldViewOutput) {
			b521FieldViewResult := b521FieldViewValue.(map[string]any)
			if b521FieldViewName == "opaque" {
				if b521FieldViewResult["state"] != "Unknown" {
					b521FieldViewT.Fatal("opaque field interpreted")
				}
				continue
			}
			if b521FieldViewResult["state"] != "incomplete" || b521FieldViewResult["present"] != false {
				b521FieldViewT.Fatal("unlocated field invented")
			}
			if _, b521FieldViewExists := b521FieldViewResult["number"]; b521FieldViewExists {
				b521FieldViewT.Fatal("unlocated scalar invented")
			}
		}
	}
}

func TestAb13Aq13UnsignedDomain(b521FieldViewT *b521FieldViewTesting.T) {
	for _, b521FieldViewDirection := range []string{"request", "positive"} {
		for b521FieldViewCarrier := 0; b521FieldViewCarrier < 256; b521FieldViewCarrier++ {
			b521FieldViewInput := b521FieldViewInput(b521FieldViewT, b521FieldViewDirection)
			b521FieldViewExpected := map[string]float64{"brine_percent": float64(b521FieldViewCarrier), "fan_percent": float64(b521FieldViewCarrier), "fan_speed": float64(b521FieldViewCarrier * 10)}
			if b521FieldViewDirection == "request" {
				b521FieldViewExpected = map[string]float64{"noise_percent": float64(b521FieldViewCarrier), "frost_temperature": float64(b521FieldViewCarrier)}
			}
			for b521FieldViewName := range b521FieldViewExpected {
				b521FieldViewInput["fields"].(map[string]any)[b521FieldViewName] = []any{b521FieldViewCarrier}
			}
			b521FieldViewResults := b521FieldViewResults(b521FieldViewEvaluate(b521FieldViewT, b521FieldViewInput))
			for b521FieldViewName, b521FieldViewNumber := range b521FieldViewExpected {
				if b521FieldViewResults[b521FieldViewName].(map[string]any)["number"] != b521FieldViewNumber {
					b521FieldViewT.Fatal("unsigned carrier", b521FieldViewName, b521FieldViewCarrier)
				}
			}
		}
	}
}
