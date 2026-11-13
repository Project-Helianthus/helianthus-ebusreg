package vaillant

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func b525QualityConfigInput(b525QualityConfigOperation int, b525QualityConfigDirection, b525QualityConfigView string, b525QualityConfigFields map[string]any) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 37, "operation": b525QualityConfigOperation, "direction": b525QualityConfigDirection, "view": b525QualityConfigView, "family": "synthetic-family", "variant": "synthetic-variant", "opaque": map[string]any{"flag": false, "zero": 0, "null": nil}}, "fields": b525QualityConfigFields}
}

func b525QualityConfigEvaluate(b525QualityConfigT *testing.T, b525QualityConfigInputObject map[string]any) map[string]any {
	b525QualityConfigT.Helper()
	b525QualityConfigBytes, b525QualityConfigError := json.Marshal(b525QualityConfigInputObject)
	if b525QualityConfigError != nil {
		b525QualityConfigT.Fatal(b525QualityConfigError)
	}
	b525QualityConfigOutput, b525QualityConfigError := B525QualityConfigEvaluate(b525QualityConfigBytes)
	if b525QualityConfigError != nil {
		b525QualityConfigT.Fatal(b525QualityConfigError)
	}
	var b525QualityConfigResult map[string]any
	if b525QualityConfigError := json.Unmarshal(b525QualityConfigOutput, &b525QualityConfigResult); b525QualityConfigError != nil {
		b525QualityConfigT.Fatal(b525QualityConfigError)
	}
	var b525QualityConfigRoundTrip map[string]any
	_ = json.Unmarshal(b525QualityConfigBytes, &b525QualityConfigRoundTrip)
	if !reflect.DeepEqual(b525QualityConfigRoundTrip["context"], b525QualityConfigResult["context"]) || !reflect.DeepEqual(b525QualityConfigRoundTrip["fields"], b525QualityConfigResult["raw"]) {
		b525QualityConfigT.Fatalf("raw/context changed: %s", b525QualityConfigOutput)
	}
	if b525QualityConfigArguments, b525QualityConfigExists := b525QualityConfigRoundTrip["arguments"]; b525QualityConfigExists && !reflect.DeepEqual(b525QualityConfigArguments, b525QualityConfigResult["arguments"]) {
		b525QualityConfigT.Fatal("arguments changed")
	}
	if b525QualityConfigResult["physical"] != "Unknown" || b525QualityConfigResult["qualified"] != false {
		b525QualityConfigT.Fatal("qualification gained")
	}
	return b525QualityConfigResult
}

func b525QualityConfigField(b525QualityConfigOutput map[string]any, b525QualityConfigName string) map[string]any {
	return b525QualityConfigOutput["fields"].(map[string]any)[b525QualityConfigName].(map[string]any)
}

func TestAb25Ap25RSSI(b525QualityConfigT *testing.T) {
	b525QualityConfigBytes, b525QualityConfigError := os.ReadFile("testdata/ab25ap25_b525_quality_config.json")
	if b525QualityConfigError != nil {
		b525QualityConfigT.Fatal(b525QualityConfigError)
	}
	var b525QualityConfigCases []map[string]any
	if b525QualityConfigError := json.Unmarshal(b525QualityConfigBytes, &b525QualityConfigCases); b525QualityConfigError != nil {
		b525QualityConfigT.Fatal(b525QualityConfigError)
	}
	if len(b525QualityConfigCases) != 4 {
		b525QualityConfigT.Fatal("missing cases")
	}
	for _, b525QualityConfigCase := range b525QualityConfigCases {
		b525QualityConfigT.Run(b525QualityConfigCase["id"].(string), func(b525QualityConfigT *testing.T) {
			b525QualityConfigResult := b525QualityConfigEvaluate(b525QualityConfigT, b525QualityConfigInput(3, "positive", "quality", map[string]any{"rssi": []any{b525QualityConfigCase["byte"]}, "quality": []int{0}}))
			b525QualityConfigRSSI := b525QualityConfigField(b525QualityConfigResult, "rssi")
			if b525QualityConfigResult["state"] != "complete" || b525QualityConfigRSSI["number"] != b525QualityConfigCase["number"] || b525QualityConfigRSSI["unit"] != "Unknown" || b525QualityConfigRSSI["state"] != "Hypothesis" || b525QualityConfigRSSI["calculation"] != "Proven" {
				b525QualityConfigT.Fatalf("RSSI: %v", b525QualityConfigRSSI)
			}
		})
	}
}

func TestAb25Ap25Quality(b525QualityConfigT *testing.T) {
	for b525QualityConfigBit := -1; b525QualityConfigBit < 8; b525QualityConfigBit++ {
		b525QualityConfigByteValue := 0
		if b525QualityConfigBit >= 0 {
			b525QualityConfigByteValue = 1 << b525QualityConfigBit
		}
		b525QualityConfigT.Run(fmt.Sprintf("ab25ap25_bit_%d", b525QualityConfigBit), func(b525QualityConfigT *testing.T) {
			b525QualityConfigResult := b525QualityConfigEvaluate(b525QualityConfigT, b525QualityConfigInput(3, "positive", "quality", map[string]any{"rssi": []int{255}, "quality": []int{b525QualityConfigByteValue}}))
			b525QualityConfigQuality := b525QualityConfigField(b525QualityConfigResult, "quality")
			b525QualityConfigFlags, b525QualityConfigOK := b525QualityConfigQuality["flags"].(map[string]any)
			if !b525QualityConfigOK {
				b525QualityConfigT.Fatal("missing independent quality flags")
			}
			if len(b525QualityConfigFlags) != 8 || b525QualityConfigQuality["checksum_verified"] != false || b525QualityConfigQuality["state"] != "Hypothesis" {
				b525QualityConfigT.Fatal(b525QualityConfigQuality)
			}
			for b525QualityConfigIndex := 0; b525QualityConfigIndex < 8; b525QualityConfigIndex++ {
				b525QualityConfigName := fmt.Sprintf("bit%d", b525QualityConfigIndex)
				if b525QualityConfigIndex == 7 {
					b525QualityConfigName = "crc_okay"
				}
				if b525QualityConfigFlags[b525QualityConfigName] != (b525QualityConfigIndex == b525QualityConfigBit) {
					b525QualityConfigT.Fatal(b525QualityConfigFlags)
				}
			}
			for _, b525QualityConfigName := range []string{"number", "score", "label"} {
				if _, b525QualityConfigExists := b525QualityConfigQuality[b525QualityConfigName]; b525QualityConfigExists {
					b525QualityConfigT.Fatal("invented quality quantity")
				}
			}
		})
	}
}

func TestAb25Ap25Configuration(b525QualityConfigT *testing.T) {
	for _, b525QualityConfigOperation := range []int{24, 25, 26, 27} {
		for _, b525QualityConfigDirection := range []string{"request", "positive"} {
			for _, b525QualityConfigConfig := range []int{0, 129, 254, 1} {
				for b525QualityConfigDataCase := 0; b525QualityConfigDataCase < 3; b525QualityConfigDataCase++ {
					b525QualityConfigT.Run(fmt.Sprintf("ab25ap25_%d_%s_%d_%d", b525QualityConfigOperation, b525QualityConfigDirection, b525QualityConfigConfig, b525QualityConfigDataCase), func(b525QualityConfigT *testing.T) {
						b525QualityConfigFields := map[string]any{"configuration": []int{b525QualityConfigConfig}}
						if b525QualityConfigDataCase == 0 {
							b525QualityConfigFields["data"] = []int{0, 170, 255}
						}
						if b525QualityConfigDataCase == 1 {
							b525QualityConfigFields["data"] = []int{}
						}
						b525QualityConfigResult := b525QualityConfigEvaluate(b525QualityConfigT, b525QualityConfigInput(b525QualityConfigOperation, b525QualityConfigDirection, "configuration", b525QualityConfigFields))
						b525QualityConfigConfigResult := b525QualityConfigField(b525QualityConfigResult, "configuration")
						if b525QualityConfigConfigResult["state"] != "Hypothesis" || b525QualityConfigConfigResult["unmodeled"] != float64(b525QualityConfigConfig&254) || b525QualityConfigConfigResult["flags"].(map[string]any)["bit0"] != (b525QualityConfigConfig&1 != 0) {
							b525QualityConfigT.Fatal(b525QualityConfigConfigResult)
						}
						if !reflect.DeepEqual(b525QualityConfigConfigResult["candidate_default"], map[string]any{"number": float64(0), "label": "read", "state": "Hypothesis"}) {
							b525QualityConfigT.Fatal("default metadata")
						}
						b525QualityConfigData := b525QualityConfigField(b525QualityConfigResult, "data")
						if b525QualityConfigDataCase == 2 {
							if b525QualityConfigResult["state"] != "partial" || b525QualityConfigData["state"] != "Unknown" {
								b525QualityConfigT.Fatal(b525QualityConfigResult)
							}
							if _, b525QualityConfigExists := b525QualityConfigData["raw"]; b525QualityConfigExists {
								b525QualityConfigT.Fatal("missing became empty")
							}
						} else if b525QualityConfigResult["state"] != "complete" || b525QualityConfigData["state"] != "Hypothesis" {
							b525QualityConfigT.Fatal(b525QualityConfigResult)
						}
					})
				}
			}
		}
	}
}

func TestAb25Ap25Incomplete(b525QualityConfigT *testing.T) {
	for _, b525QualityConfigView := range []string{"quality", "configuration"} {
		b525QualityConfigNames := []string{"rssi", "quality"}
		b525QualityConfigOperation := 3
		if b525QualityConfigView == "configuration" {
			b525QualityConfigNames = []string{"configuration", "data"}
			b525QualityConfigOperation = 24
		}
		for _, b525QualityConfigName := range b525QualityConfigNames {
			for b525QualityConfigCase := 0; b525QualityConfigCase < 3; b525QualityConfigCase++ {
				if b525QualityConfigName == "data" && b525QualityConfigCase > 0 {
					continue
				}
				b525QualityConfigFields := map[string]any{b525QualityConfigNames[0]: []int{0}, b525QualityConfigNames[1]: []int{0}}
				delete(b525QualityConfigFields, b525QualityConfigName)
				if b525QualityConfigCase == 1 {
					b525QualityConfigFields[b525QualityConfigName] = []int{}
				}
				if b525QualityConfigCase == 2 {
					b525QualityConfigFields[b525QualityConfigName] = []int{0, 255}
				}
				b525QualityConfigResult := b525QualityConfigEvaluate(b525QualityConfigT, b525QualityConfigInput(b525QualityConfigOperation, "positive", b525QualityConfigView, b525QualityConfigFields))
				if b525QualityConfigResult["state"] != "partial" || b525QualityConfigField(b525QualityConfigResult, b525QualityConfigName)["state"] != "Unknown" {
					b525QualityConfigT.Fatal(b525QualityConfigResult)
				}
				for _, b525QualityConfigSibling := range b525QualityConfigNames {
					if b525QualityConfigSibling != b525QualityConfigName && b525QualityConfigField(b525QualityConfigResult, b525QualityConfigSibling)["state"] != "Hypothesis" {
						b525QualityConfigT.Fatal("lost sibling")
					}
				}
				for _, b525QualityConfigKey := range []string{"number", "label", "flags"} {
					if _, b525QualityConfigExists := b525QualityConfigField(b525QualityConfigResult, b525QualityConfigName)[b525QualityConfigKey]; b525QualityConfigExists {
						b525QualityConfigT.Fatal("invented incomplete value")
					}
				}
			}
		}
	}
}

func TestAb25Ap25ContextArguments(b525QualityConfigT *testing.T) {
	for _, b525QualityConfigKey := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant", "opaque"} {
		for _, b525QualityConfigValue := range []any{nil, 0, false, "other", map[string]any{"nested": []any{0, false, nil}}} {
			if (b525QualityConfigKey == "primary" || b525QualityConfigKey == "secondary" || b525QualityConfigKey == "operation") && b525QualityConfigValue != nil && b525QualityConfigValue != 0 {
				continue
			}
			if (b525QualityConfigKey == "direction" || b525QualityConfigKey == "view") && b525QualityConfigValue != nil && b525QualityConfigValue != "other" {
				continue
			}
			b525QualityConfigInputObject := b525QualityConfigInput(3, "positive", "quality", map[string]any{"rssi": []int{0}, "quality": []int{0}})
			b525QualityConfigInputObject["context"].(map[string]any)[b525QualityConfigKey] = b525QualityConfigValue
			b525QualityConfigResult := b525QualityConfigEvaluate(b525QualityConfigT, b525QualityConfigInputObject)
			if (b525QualityConfigKey == "family" || b525QualityConfigKey == "variant" || b525QualityConfigKey == "opaque") && b525QualityConfigResult["state"] != "complete" {
				b525QualityConfigT.Fatal("opaque key affected admission")
			}
			if b525QualityConfigKey != "family" && b525QualityConfigKey != "variant" && b525QualityConfigKey != "opaque" && (b525QualityConfigResult["state"] != "Unknown" || len(b525QualityConfigResult["fields"].(map[string]any)) != 0) {
				b525QualityConfigT.Fatal("selector admitted")
			}
		}
		b525QualityConfigInputObject := b525QualityConfigInput(3, "positive", "quality", map[string]any{"rssi": []int{0}, "quality": []int{0}})
		delete(b525QualityConfigInputObject["context"].(map[string]any), b525QualityConfigKey)
		b525QualityConfigResult := b525QualityConfigEvaluate(b525QualityConfigT, b525QualityConfigInputObject)
		if b525QualityConfigKey != "family" && b525QualityConfigKey != "variant" && b525QualityConfigKey != "opaque" && b525QualityConfigResult["state"] != "Unknown" {
			b525QualityConfigT.Fatal("absent selector admitted")
		}
	}
	for _, b525QualityConfigValue := range []any{nil, false, 0, []any{}, map[string]any{}, "opaque"} {
		b525QualityConfigInputObject := b525QualityConfigInput(3, "positive", "quality", map[string]any{"rssi": []int{0}, "quality": []int{0}, "extra": []int{255}})
		b525QualityConfigInputObject["arguments"] = map[string]any{"extra": b525QualityConfigValue}
		b525QualityConfigEvaluate(b525QualityConfigT, b525QualityConfigInputObject)
	}
}

func TestAb25Ap25Invalid(b525QualityConfigT *testing.T) {
	b525QualityConfigValid := `{"context":{"primary":181,"secondary":37,"operation":3,"direction":"positive","view":"quality"},"fields":{"rssi":[0],"quality":[0]}}`
	b525QualityConfigCases := []string{``, `{`, `[]`, b525QualityConfigValid + ` {}`, strings.Replace(b525QualityConfigValid, `"primary":181`, `"primary":181,"primary":181`, 1), strings.Replace(b525QualityConfigValid, `"rssi":[0]`, `"rssi":[0],"rssi":[0]`, 1), strings.Replace(b525QualityConfigValid, `"rssi":[0]`, `"rssi":null`, 1), strings.Replace(b525QualityConfigValid, `"context":{`, `"arguments":null,"context":{`, 1), strings.Replace(b525QualityConfigValid, `"context":{`, `"arguments":{"x":{"a":0,"a":1}},"context":{`, 1)}
	for _, b525QualityConfigValue := range []string{`true`, `false`, `null`, `"0"`, `-1`, `256`, `0.5`, `1.00000000000000000001`, `[]`, `{}`} {
		b525QualityConfigCases = append(b525QualityConfigCases, strings.Replace(b525QualityConfigValid, `"rssi":[0]`, `"rssi":[`+b525QualityConfigValue+`]`, 1))
	}
	for _, b525QualityConfigInputBytes := range b525QualityConfigCases {
		b525QualityConfigOutput, b525QualityConfigError := B525QualityConfigEvaluate([]byte(b525QualityConfigInputBytes))
		if b525QualityConfigError == nil || b525QualityConfigOutput != nil {
			b525QualityConfigT.Fatalf("invalid admitted %s", b525QualityConfigInputBytes)
		}
	}
}

func TestAb25Ap25UTF8(b525QualityConfigT *testing.T) {
	for _, b525QualityConfigLocation := range []string{"context", "arguments"} {
		b525QualityConfigT.Run(b525QualityConfigLocation, func(b525QualityConfigT *testing.T) {
			b525QualityConfigTemplate := `{"context":{},"fields":{},"arguments":{"opaque":"%s"}}`
			if b525QualityConfigLocation == "context" {
				b525QualityConfigTemplate = `{"context":{"opaque":"%s"},"fields":{}}`
			}
			for _, b525QualityConfigText := range []string{"\xff", "\xc3", "\xc0\xaf", "\xed\xa0\x80"} {
				b525QualityConfigInputBytes := fmt.Sprintf(b525QualityConfigTemplate, b525QualityConfigText)
				b525QualityConfigOutput, b525QualityConfigError := B525QualityConfigEvaluate([]byte(b525QualityConfigInputBytes))
				if b525QualityConfigError == nil || b525QualityConfigOutput != nil {
					b525QualityConfigT.Fatalf("invalid UTF-8 admitted: %q", b525QualityConfigInputBytes)
				}
			}
			b525QualityConfigInputObject := b525QualityConfigInput(3, "positive", "quality", map[string]any{"rssi": []int{0}, "quality": []int{0}})
			b525QualityConfigInputObject["arguments"] = map[string]any{}
			b525QualityConfigInputObject[b525QualityConfigLocation].(map[string]any)["opaque"] = "é 世界 \ufffd"
			b525QualityConfigEvaluate(b525QualityConfigT, b525QualityConfigInputObject)
		})
	}
}

func TestAb25Ap25SurrogateEscapes(b525QualityConfigT *testing.T) {
	for _, b525QualityConfigTemplate := range []string{
		`{"context":{},"fields":{},"arguments":{"opaque":"%s"}}`,
		`{"context":{"opaque":"%s"},"fields":{}}`,
		`{"context":{},"fields":{},"arguments":{"%s":null}}`,
		`{"context":{},"fields":{},"arguments":{"nested":[{"opaque":"%s"}]}}`,
	} {
		for _, b525QualityConfigText := range []string{
			`\ud800`, `\udbff`, `\udc00`, `\udfff`,
			`\ud800x`, `\ud800\u0041`, `\ud800\ud800`,
			`\udc00\ud800`, `\ud800\\udc00`, `\ud800\udc00\udfff`,
		} {
			b525QualityConfigInputBytes := fmt.Sprintf(b525QualityConfigTemplate, b525QualityConfigText)
			b525QualityConfigOutput, b525QualityConfigError := B525QualityConfigEvaluate([]byte(b525QualityConfigInputBytes))
			if b525QualityConfigError == nil || b525QualityConfigOutput != nil {
				b525QualityConfigT.Errorf("invalid surrogate escape admitted: %s", b525QualityConfigInputBytes)
			}
		}
		for _, b525QualityConfigText := range []string{
			`\ud800\udc00`, `\uDBFF\uDFFF`, `\uD83D\uDE00`,
			`\ufffd`, `�`, `\\ud800`, `\\\ud800\udc00`, `\"\ud800\udc00`,
		} {
			b525QualityConfigInputBytes := fmt.Sprintf(b525QualityConfigTemplate, b525QualityConfigText)
			b525QualityConfigOutput, b525QualityConfigError := B525QualityConfigEvaluate([]byte(b525QualityConfigInputBytes))
			if b525QualityConfigError != nil {
				b525QualityConfigT.Errorf("valid Unicode rejected: %s: %v", b525QualityConfigInputBytes, b525QualityConfigError)
				continue
			}
			var b525QualityConfigBefore, b525QualityConfigAfter map[string]any
			if b525QualityConfigError := json.Unmarshal([]byte(b525QualityConfigInputBytes), &b525QualityConfigBefore); b525QualityConfigError != nil {
				b525QualityConfigT.Fatal(b525QualityConfigError)
			}
			if b525QualityConfigError := json.Unmarshal(b525QualityConfigOutput, &b525QualityConfigAfter); b525QualityConfigError != nil {
				b525QualityConfigT.Fatal(b525QualityConfigError)
			}
			for _, b525QualityConfigKey := range []string{"context", "arguments"} {
				if !reflect.DeepEqual(b525QualityConfigBefore[b525QualityConfigKey], b525QualityConfigAfter[b525QualityConfigKey]) {
					b525QualityConfigT.Errorf("%s changed for %s", b525QualityConfigKey, b525QualityConfigInputBytes)
				}
			}
		}
	}
}

func TestAb25Ap25Ownership(b525QualityConfigT *testing.T) {
	b525QualityConfigInputBytes := json.RawMessage(`{"context":{"primary":181,"secondary":37,"operation":24,"direction":"request","view":"configuration","opaque":{"a":[0]}},"fields":{"configuration":[0],"data":[0,170,255]},"arguments":{"opaque":{"a":[false]}}}`)
	b525QualityConfigSaved := append(json.RawMessage(nil), b525QualityConfigInputBytes...)
	b525QualityConfigFirst, b525QualityConfigError := B525QualityConfigEvaluate(b525QualityConfigInputBytes)
	if b525QualityConfigError != nil {
		b525QualityConfigT.Fatal(b525QualityConfigError)
	}
	b525QualityConfigExpected := append(json.RawMessage(nil), b525QualityConfigFirst...)
	for b525QualityConfigIndex := range b525QualityConfigInputBytes {
		b525QualityConfigInputBytes[b525QualityConfigIndex] = ' '
	}
	if string(b525QualityConfigFirst) != string(b525QualityConfigExpected) {
		b525QualityConfigT.Fatal("input aliases output")
	}
	var b525QualityConfigDecoded map[string]any
	_ = json.Unmarshal(b525QualityConfigFirst, &b525QualityConfigDecoded)
	b525QualityConfigDecoded["raw"].(map[string]any)["data"].([]any)[0] = 99
	b525QualityConfigDecoded["context"].(map[string]any)["opaque"].(map[string]any)["a"].([]any)[0] = 99
	b525QualityConfigDecoded["arguments"].(map[string]any)["opaque"].(map[string]any)["a"].([]any)[0] = true
	if b525QualityConfigField(b525QualityConfigDecoded, "data")["raw"].([]any)[0] != float64(0) {
		b525QualityConfigT.Fatal("raw aliases semantic data")
	}
	for b525QualityConfigIndex := range b525QualityConfigFirst {
		b525QualityConfigFirst[b525QualityConfigIndex] = ' '
	}
	b525QualityConfigSecond, b525QualityConfigError := B525QualityConfigEvaluate(b525QualityConfigSaved)
	if b525QualityConfigError != nil || string(b525QualityConfigSecond) != string(b525QualityConfigExpected) {
		b525QualityConfigT.Fatal("cross-call alias")
	}
}

func TestAb25Ap25Admission(b525QualityConfigT *testing.T) {
	for _, b525QualityConfigView := range []string{"quality", "configuration"} {
		b525QualityConfigOperation := 3
		b525QualityConfigNames := []string{"rssi", "quality"}
		if b525QualityConfigView == "configuration" {
			b525QualityConfigOperation = 24
			b525QualityConfigNames = []string{"configuration", "data"}
		}
		for _, b525QualityConfigName := range b525QualityConfigNames {
			for _, b525QualityConfigBad := range []any{nil, false, 0, "", []any{nil}, []any{false}, []any{"0"}, []any{-1}, []any{256}, []any{0.5}} {
				b525QualityConfigInputObject := b525QualityConfigInput(b525QualityConfigOperation, "positive", b525QualityConfigView, map[string]any{b525QualityConfigNames[0]: []int{0}, b525QualityConfigNames[1]: []int{0}})
				b525QualityConfigInputObject["fields"].(map[string]any)[b525QualityConfigName] = b525QualityConfigBad
				b525QualityConfigBytes, _ := json.Marshal(b525QualityConfigInputObject)
				b525QualityConfigOutput, b525QualityConfigError := B525QualityConfigEvaluate(b525QualityConfigBytes)
				if b525QualityConfigError == nil || b525QualityConfigOutput != nil {
					b525QualityConfigT.Fatalf("invalid %s admitted: %s", b525QualityConfigName, b525QualityConfigBytes)
				}
			}
		}
		for _, b525QualityConfigKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
			b525QualityConfigInputObject := b525QualityConfigInput(b525QualityConfigOperation, "positive", b525QualityConfigView, map[string]any{b525QualityConfigNames[0]: []int{0}, b525QualityConfigNames[1]: []int{0}})
			b525QualityConfigInputObject["context"].(map[string]any)[b525QualityConfigKey] = nil
			b525QualityConfigResult := b525QualityConfigEvaluate(b525QualityConfigT, b525QualityConfigInputObject)
			if b525QualityConfigResult["state"] != "Unknown" || len(b525QualityConfigResult["fields"].(map[string]any)) != 0 {
				b525QualityConfigT.Fatal("null selector admitted")
			}
		}
	}
	for _, b525QualityConfigDirection := range []string{"request", "negative", ""} {
		b525QualityConfigResult := b525QualityConfigEvaluate(b525QualityConfigT, b525QualityConfigInput(3, b525QualityConfigDirection, "quality", map[string]any{"rssi": []int{0}, "quality": []int{0}}))
		if b525QualityConfigResult["state"] != "Unknown" {
			b525QualityConfigT.Fatal("quality direction admitted")
		}
	}
	for _, b525QualityConfigOperation := range []int{3, 23, 28} {
		b525QualityConfigResult := b525QualityConfigEvaluate(b525QualityConfigT, b525QualityConfigInput(b525QualityConfigOperation, "positive", "configuration", map[string]any{"configuration": []int{0}, "data": []int{}}))
		if b525QualityConfigResult["state"] != "Unknown" {
			b525QualityConfigT.Fatal("configuration operation admitted")
		}
	}
	b525QualityConfigValid := json.RawMessage(`{"context":{"primary":181,"secondary":37,"operation":3,"direction":"positive","view":"quality"},"fields":{"rssi":[0e0],"quality":[1.0]}}`)
	if _, b525QualityConfigError := B525QualityConfigEvaluate(b525QualityConfigValid); b525QualityConfigError != nil {
		b525QualityConfigT.Fatal("integral JSON numbers rejected", b525QualityConfigError)
	}
}

func TestAb25Ap25OversizedExponentRejected(b525QualityConfigT *testing.T) {
	for _, b525QualityConfigNumber := range []string{"1e999999999", "0e999999999", "1e-999999999", "1e65", "1E65", "1e99999999999999999999"} {
		for _, b525QualityConfigTemplate := range []string{
			`{"context":{"primary":%s},"fields":{}}`,
			`{"context":{},"fields":{"rssi":[%s]}}`,
		} {
			b525QualityConfigInputBytes := fmt.Sprintf(b525QualityConfigTemplate, b525QualityConfigNumber)
			b525QualityConfigOutput, b525QualityConfigError := B525QualityConfigEvaluate([]byte(b525QualityConfigInputBytes))
			if b525QualityConfigError == nil || b525QualityConfigOutput != nil {
				b525QualityConfigT.Errorf("oversized exponent admitted: %s", b525QualityConfigInputBytes)
			}
		}
	}
	b525QualityConfigOutput, b525QualityConfigError := B525QualityConfigEvaluate([]byte(`{"context":{"primary":1.81e2},"fields":{"rssi":[2.5e1]}}`))
	if b525QualityConfigError != nil || b525QualityConfigOutput == nil {
		b525QualityConfigT.Errorf("small exponent rejected: %v", b525QualityConfigError)
	}
}

func TestAb25Ap25NestingDepthRejected(b525QualityConfigT *testing.T) {
	b525QualityConfigDeep := strings.Repeat("[", 1000000)
	for _, b525QualityConfigInput := range []string{
		`{"context":{"x":` + b525QualityConfigDeep + `},"fields":{}}`,
		`{"context":{},"fields":{"rssi":` + b525QualityConfigDeep + `}}`,
		`{"context":{"x":` + strings.Repeat(`{"a":`, 1000000) + `1},"fields":{}}`,
	} {
		b525QualityConfigOutput, b525QualityConfigError := B525QualityConfigEvaluate([]byte(b525QualityConfigInput))
		if b525QualityConfigError == nil || b525QualityConfigOutput != nil {
			b525QualityConfigT.Errorf("excessive nesting admitted")
		}
	}
	b525QualityConfigOutput, b525QualityConfigError := B525QualityConfigEvaluate([]byte(`{"context":{"x":` + strings.Repeat("[", 50) + strings.Repeat("]", 50) + `},"fields":{}}`))
	if b525QualityConfigError != nil || b525QualityConfigOutput == nil {
		b525QualityConfigT.Errorf("moderate nesting rejected: %v", b525QualityConfigError)
	}
}
