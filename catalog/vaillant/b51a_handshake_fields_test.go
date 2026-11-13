package vaillant

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func b51AHandshakeFieldEvaluateTest(b51AHandshakeFieldT *testing.T, b51AHandshakeFieldInput any) map[string]any {
	b51AHandshakeFieldT.Helper()
	b51AHandshakeFieldJSON, b51AHandshakeFieldErr := json.Marshal(b51AHandshakeFieldInput)
	if b51AHandshakeFieldErr != nil {
		b51AHandshakeFieldT.Fatal(b51AHandshakeFieldErr)
	}
	b51AHandshakeFieldOutput, b51AHandshakeFieldErr := B51AHandshakeFieldEvaluate(b51AHandshakeFieldJSON)
	if b51AHandshakeFieldErr != nil {
		b51AHandshakeFieldT.Fatal(b51AHandshakeFieldErr)
	}
	b51AHandshakeFieldResult := map[string]any{}
	if b51AHandshakeFieldErr := json.Unmarshal(b51AHandshakeFieldOutput, &b51AHandshakeFieldResult); b51AHandshakeFieldErr != nil {
		b51AHandshakeFieldT.Fatal(b51AHandshakeFieldErr)
	}
	return b51AHandshakeFieldResult
}

func b51AHandshakeFieldCases(b51AHandshakeFieldT *testing.T) []map[string]any {
	b51AHandshakeFieldT.Helper()
	b51AHandshakeFieldData, b51AHandshakeFieldErr := os.ReadFile("testdata/ab19aq19_b51a_handshake_fields.json")
	if b51AHandshakeFieldErr != nil {
		b51AHandshakeFieldT.Fatal(b51AHandshakeFieldErr)
	}
	var b51AHandshakeFieldCases []map[string]any
	if b51AHandshakeFieldErr := json.Unmarshal(b51AHandshakeFieldData, &b51AHandshakeFieldCases); b51AHandshakeFieldErr != nil {
		b51AHandshakeFieldT.Fatal(b51AHandshakeFieldErr)
	}
	return b51AHandshakeFieldCases
}

func b51AHandshakeFieldPath(b51AHandshakeFieldValue any, b51AHandshakeFieldPath string) any {
	for _, b51AHandshakeFieldKey := range strings.Split(b51AHandshakeFieldPath, ".") {
		b51AHandshakeFieldMap, b51AHandshakeFieldOK := b51AHandshakeFieldValue.(map[string]any)
		if !b51AHandshakeFieldOK {
			return nil
		}
		b51AHandshakeFieldValue = b51AHandshakeFieldMap[b51AHandshakeFieldKey]
	}
	return b51AHandshakeFieldValue
}

func TestAb19Aq19Fixtures(b51AHandshakeFieldT *testing.T) {
	for _, b51AHandshakeFieldCase := range b51AHandshakeFieldCases(b51AHandshakeFieldT) {
		b51AHandshakeFieldT.Run(b51AHandshakeFieldCase["id"].(string), func(b51AHandshakeFieldT *testing.T) {
			b51AHandshakeFieldInput := b51AHandshakeFieldCase["input"].(map[string]any)
			b51AHandshakeFieldOutput := b51AHandshakeFieldEvaluateTest(b51AHandshakeFieldT, b51AHandshakeFieldInput)
			for b51AHandshakeFieldPathKey, b51AHandshakeFieldExpected := range b51AHandshakeFieldCase["expected"].(map[string]any) {
				if b51AHandshakeFieldGot := b51AHandshakeFieldPath(b51AHandshakeFieldOutput, b51AHandshakeFieldPathKey); !reflect.DeepEqual(b51AHandshakeFieldGot, b51AHandshakeFieldExpected) {
					b51AHandshakeFieldT.Errorf("%s = %v want %v", b51AHandshakeFieldPathKey, b51AHandshakeFieldGot, b51AHandshakeFieldExpected)
				}
			}
			for _, b51AHandshakeFieldKey := range []string{"context", "fields", "arguments", "extra"} {
				if !reflect.DeepEqual(b51AHandshakeFieldOutput[b51AHandshakeFieldKey], b51AHandshakeFieldInput[b51AHandshakeFieldKey]) {
					b51AHandshakeFieldT.Fatal("raw retention", b51AHandshakeFieldKey)
				}
			}
			for b51AHandshakeFieldKey, b51AHandshakeFieldRaw := range b51AHandshakeFieldInput["fields"].(map[string]any) {
				if !reflect.DeepEqual(b51AHandshakeFieldPath(b51AHandshakeFieldOutput, "results."+b51AHandshakeFieldKey+".raw"), b51AHandshakeFieldRaw) {
					b51AHandshakeFieldT.Fatal("result raw", b51AHandshakeFieldKey)
				}
			}
			if b51AHandshakeFieldOutput["physical"] != "Unknown" || b51AHandshakeFieldOutput["qualified"] != false {
				b51AHandshakeFieldT.Fatal("qualification")
			}
			b51AHandshakeFieldJSON, _ := json.Marshal(b51AHandshakeFieldOutput)
			for _, b51AHandshakeFieldForbidden := range []string{"success", "crc-valid", "validated"} {
				if bytes.Contains(b51AHandshakeFieldJSON, []byte(b51AHandshakeFieldForbidden)) {
					b51AHandshakeFieldT.Fatal("invented claim", b51AHandshakeFieldForbidden)
				}
			}
			if b51AHandshakeFieldPath(b51AHandshakeFieldOutput, "results.crc.grade") != nil && b51AHandshakeFieldPath(b51AHandshakeFieldOutput, "results.crc.grade") != "Hypothesis" {
				b51AHandshakeFieldT.Fatal("calculation promoted meaning")
			}
		})
	}
}

func TestAb19Aq19Context(b51AHandshakeFieldT *testing.T) {
	for _, b51AHandshakeFieldCase := range b51AHandshakeFieldCases(b51AHandshakeFieldT) {
		for _, b51AHandshakeFieldKey := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
			b51AHandshakeFieldT.Run(b51AHandshakeFieldCase["id"].(string)+"/"+b51AHandshakeFieldKey, func(b51AHandshakeFieldT *testing.T) {
				b51AHandshakeFieldJSON, _ := json.Marshal(b51AHandshakeFieldCase["input"])
				var b51AHandshakeFieldInput map[string]any
				if b51AHandshakeFieldErr := json.Unmarshal(b51AHandshakeFieldJSON, &b51AHandshakeFieldInput); b51AHandshakeFieldErr != nil {
					b51AHandshakeFieldT.Fatal(b51AHandshakeFieldErr)
				}
				b51AHandshakeFieldContext := b51AHandshakeFieldInput["context"].(map[string]any)
				switch b51AHandshakeFieldKey {
				case "primary", "secondary":
					b51AHandshakeFieldContext[b51AHandshakeFieldKey] = 0
				case "operation":
					b51AHandshakeFieldContext[b51AHandshakeFieldKey] = 0
				default:
					b51AHandshakeFieldContext[b51AHandshakeFieldKey] = "other"
				}
				b51AHandshakeFieldOutput := b51AHandshakeFieldEvaluateTest(b51AHandshakeFieldT, b51AHandshakeFieldInput)
				if b51AHandshakeFieldKey == "family" || b51AHandshakeFieldKey == "variant" {
					if b51AHandshakeFieldOutput["qualified"] != false || b51AHandshakeFieldOutput["state"] != "complete" {
						b51AHandshakeFieldT.Fatal("opaque identity")
					}
					return
				}
				if b51AHandshakeFieldOutput["state"] != "Unknown" {
					b51AHandshakeFieldT.Fatal("unsupported context")
				}
				for _, b51AHandshakeFieldResult := range b51AHandshakeFieldOutput["results"].(map[string]any) {
					if len(b51AHandshakeFieldResult.(map[string]any)) != 2 {
						b51AHandshakeFieldT.Fatal("formula leaked")
					}
				}
			})
		}
	}
}

func TestAb19Aq19Invalid(b51AHandshakeFieldT *testing.T) {
	for _, b51AHandshakeFieldInput := range []string{
		`{`, `[]`, `{} {}`, `{"context":{},"context":{},"fields":{}}`,
		`{"context":{},"fields":{},"extra":[{"x":0,"x":1}]}`,
		`{"context":{"operation":null},"fields":{}}`,
		`{"context":{"primary":true},"fields":{}}`,
		`{"context":{"secondary":1.5},"fields":{}}`,
		`{"context":{"operation":256},"fields":{}}`,
		`{"context":{"direction":null},"fields":{}}`,
		`{"context":{"view":false},"fields":{}}`,
		`{"context":{},"fields":{"sequence":null}}`,
		`{"context":{},"fields":{"sequence":["1"]}}`,
		`{"context":{},"fields":{"sequence":[false]}}`,
		`{"context":{},"fields":{"sequence":[null]}}`,
		`{"context":{},"fields":{"sequence":[-1]}}`,
		`{"context":{},"fields":{"sequence":[256]}}`,
		`{"context":{},"fields":{"sequence":[1.5]}}`,
		`{"context":{},"fields":{},"arguments":null}`,
		`{"context":{},"fields":{},"arguments":{"value":1}}`,
		`{"context":{},"fields":{},"extra":1e999}`,
		`{"context":{},"fields":{}} true`,
	} {
		b51AHandshakeFieldOutput, b51AHandshakeFieldErr := B51AHandshakeFieldEvaluate(json.RawMessage(b51AHandshakeFieldInput))
		if b51AHandshakeFieldErr == nil || b51AHandshakeFieldOutput != nil {
			b51AHandshakeFieldT.Errorf("accepted %s", b51AHandshakeFieldInput)
		}
	}
}

func TestAb19Aq19WidthsAndFields(b51AHandshakeFieldT *testing.T) {
	for _, b51AHandshakeFieldCase := range b51AHandshakeFieldCases(b51AHandshakeFieldT) {
		for b51AHandshakeFieldKey := range b51AHandshakeFieldCase["input"].(map[string]any)["fields"].(map[string]any) {
			for _, b51AHandshakeFieldWidth := range []int{-1, 0, 1, 2, 4} {
				b51AHandshakeFieldJSON, _ := json.Marshal(b51AHandshakeFieldCase["input"])
				var b51AHandshakeFieldInput map[string]any
				if b51AHandshakeFieldErr := json.Unmarshal(b51AHandshakeFieldJSON, &b51AHandshakeFieldInput); b51AHandshakeFieldErr != nil {
					b51AHandshakeFieldT.Fatal(b51AHandshakeFieldErr)
				}
				b51AHandshakeFieldFields := b51AHandshakeFieldInput["fields"].(map[string]any)
				b51AHandshakeFieldData := []any{}
				for b51AHandshakeFieldIndex := 0; b51AHandshakeFieldIndex < b51AHandshakeFieldWidth; b51AHandshakeFieldIndex++ {
					b51AHandshakeFieldData = append(b51AHandshakeFieldData, float64(7+b51AHandshakeFieldIndex))
				}
				if b51AHandshakeFieldWidth < 0 {
					delete(b51AHandshakeFieldFields, b51AHandshakeFieldKey)
				} else {
					b51AHandshakeFieldFields[b51AHandshakeFieldKey] = b51AHandshakeFieldData
				}
				b51AHandshakeFieldOutput := b51AHandshakeFieldEvaluateTest(b51AHandshakeFieldT, b51AHandshakeFieldInput)
				b51AHandshakeFieldResult := b51AHandshakeFieldPath(b51AHandshakeFieldOutput, "results."+b51AHandshakeFieldKey).(map[string]any)
				b51AHandshakeFieldOpaque := b51AHandshakeFieldKey == "dsn" || b51AHandshakeFieldKey == "start" || b51AHandshakeFieldKey == "offset"
				if b51AHandshakeFieldOpaque {
					if b51AHandshakeFieldWidth >= 0 && !reflect.DeepEqual(b51AHandshakeFieldResult["raw"], b51AHandshakeFieldData) {
						b51AHandshakeFieldT.Fatal("opaque lost")
					}
					if b51AHandshakeFieldResult["state"] == "incomplete" || b51AHandshakeFieldOutput["state"] != "complete" {
						b51AHandshakeFieldT.Fatal("opaque width imposed")
					}
				} else {
					b51AHandshakeFieldRequired := 1
					if b51AHandshakeFieldKey == "crc" {
						b51AHandshakeFieldRequired = 2
					}
					if b51AHandshakeFieldWidth != b51AHandshakeFieldRequired {
						if b51AHandshakeFieldOutput["state"] != "partial" || b51AHandshakeFieldResult["state"] != "incomplete" || b51AHandshakeFieldResult["number"] != nil || b51AHandshakeFieldResult["flags"] != nil {
							b51AHandshakeFieldT.Fatal("incomplete conversion", b51AHandshakeFieldKey, b51AHandshakeFieldWidth)
						}
					} else if b51AHandshakeFieldResult["state"] != "complete" {
						b51AHandshakeFieldT.Fatal("valid width")
					}
					if b51AHandshakeFieldKey != "sequence" && b51AHandshakeFieldFields["sequence"] != nil && b51AHandshakeFieldPath(b51AHandshakeFieldOutput, "results.sequence.state") != "complete" {
						b51AHandshakeFieldT.Fatal("sibling lost")
					}
				}
				if b51AHandshakeFieldWidth < 0 {
					if _, b51AHandshakeFieldPresent := b51AHandshakeFieldResult["raw"]; b51AHandshakeFieldPresent {
						b51AHandshakeFieldT.Fatal("missing became null")
					}
				}
			}
		}
	}
}

func TestAb19Aq19Memory(b51AHandshakeFieldT *testing.T) {
	b51AHandshakeFieldJSON, _ := json.Marshal(b51AHandshakeFieldCases(b51AHandshakeFieldT)[3]["input"])
	b51AHandshakeFieldOriginal := append([]byte(nil), b51AHandshakeFieldJSON...)
	b51AHandshakeFieldFirst, b51AHandshakeFieldErr := B51AHandshakeFieldEvaluate(b51AHandshakeFieldJSON)
	if b51AHandshakeFieldErr != nil {
		b51AHandshakeFieldT.Fatal(b51AHandshakeFieldErr)
	}
	b51AHandshakeFieldSaved := append([]byte(nil), b51AHandshakeFieldFirst...)
	for b51AHandshakeFieldIndex := range b51AHandshakeFieldJSON {
		b51AHandshakeFieldJSON[b51AHandshakeFieldIndex] = ' '
	}
	if !bytes.Equal(b51AHandshakeFieldFirst, b51AHandshakeFieldSaved) {
		b51AHandshakeFieldT.Fatal("input alias")
	}
	for b51AHandshakeFieldIndex := range b51AHandshakeFieldFirst {
		b51AHandshakeFieldFirst[b51AHandshakeFieldIndex] = ' '
	}
	b51AHandshakeFieldSecond, b51AHandshakeFieldErr := B51AHandshakeFieldEvaluate(b51AHandshakeFieldOriginal)
	if b51AHandshakeFieldErr != nil || !bytes.Equal(b51AHandshakeFieldSecond, b51AHandshakeFieldSaved) {
		b51AHandshakeFieldT.Fatal("cross call alias")
	}
}

func TestAb19Aq19IndependentBitsAndNumbers(b51AHandshakeFieldT *testing.T) {
	for _, b51AHandshakeFieldDirection := range []string{"request", "positive"} {
		b51AHandshakeFieldLabels := []string{"power-on/start", "change", "CRC-support"}
		if b51AHandshakeFieldDirection == "positive" {
			b51AHandshakeFieldLabels = []string{"valid", "invalid", "blocked", "change", "denied", "invalid-range", "not-user-adjustable"}
		}
		for b51AHandshakeFieldBit := 0; b51AHandshakeFieldBit < 8; b51AHandshakeFieldBit++ {
			b51AHandshakeFieldInput := map[string]any{"context": map[string]any{"primary": 181, "secondary": 26, "operation": 1, "view": "handshake", "direction": b51AHandshakeFieldDirection}, "fields": map[string]any{"sequence": []int{13}, "flags": []int{1 << b51AHandshakeFieldBit}, "crc": []int{1, 2}, "family_code": []int{27}}}
			b51AHandshakeFieldOutput := b51AHandshakeFieldEvaluateTest(b51AHandshakeFieldT, b51AHandshakeFieldInput)
			for b51AHandshakeFieldIndex, b51AHandshakeFieldLabel := range b51AHandshakeFieldLabels {
				if b51AHandshakeFieldPath(b51AHandshakeFieldOutput, "results.flags.flags."+b51AHandshakeFieldLabel) != (b51AHandshakeFieldIndex == b51AHandshakeFieldBit) {
					b51AHandshakeFieldT.Fatal("wrong bit", b51AHandshakeFieldDirection, b51AHandshakeFieldLabel, b51AHandshakeFieldBit)
				}
			}
			b51AHandshakeFieldReserved := 0
			if b51AHandshakeFieldBit >= len(b51AHandshakeFieldLabels) {
				b51AHandshakeFieldReserved = 1 << b51AHandshakeFieldBit
			}
			if b51AHandshakeFieldPath(b51AHandshakeFieldOutput, "results.flags.reserved") != float64(b51AHandshakeFieldReserved) {
				b51AHandshakeFieldT.Fatal("reserved")
			}
			if b51AHandshakeFieldPath(b51AHandshakeFieldOutput, "results.sequence.number") != float64(13) {
				b51AHandshakeFieldT.Fatal("sequence")
			}
			if b51AHandshakeFieldDirection == "request" && b51AHandshakeFieldPath(b51AHandshakeFieldOutput, "results.crc.number") != float64(513) {
				b51AHandshakeFieldT.Fatal("CRC word")
			}
			if b51AHandshakeFieldDirection == "positive" && b51AHandshakeFieldPath(b51AHandshakeFieldOutput, "results.family_code.number") != float64(27) {
				b51AHandshakeFieldT.Fatal("family code")
			}
		}
	}
	for b51AHandshakeFieldBit := 0; b51AHandshakeFieldBit < 8; b51AHandshakeFieldBit++ {
		b51AHandshakeFieldInput := map[string]any{"context": map[string]any{"primary": 181, "secondary": 26, "view": "reset", "direction": "positive"}, "fields": map[string]any{"sequence": []int{0}, "ack": []int{1 << b51AHandshakeFieldBit}}}
		b51AHandshakeFieldOutput := b51AHandshakeFieldEvaluateTest(b51AHandshakeFieldT, b51AHandshakeFieldInput)
		if b51AHandshakeFieldPath(b51AHandshakeFieldOutput, "results.ack.flags.denied") != (b51AHandshakeFieldBit == 0) {
			b51AHandshakeFieldT.Fatal("ack bit")
		}
		b51AHandshakeFieldExpected := 0
		if b51AHandshakeFieldBit > 0 {
			b51AHandshakeFieldExpected = 1 << b51AHandshakeFieldBit
		}
		if b51AHandshakeFieldPath(b51AHandshakeFieldOutput, "results.ack.unmodeled") != float64(b51AHandshakeFieldExpected) {
			b51AHandshakeFieldT.Fatal("ack unmodeled")
		}
	}
}

func TestAb19Aq19AbsentIdentityAndOpaqueValues(b51AHandshakeFieldT *testing.T) {
	for _, b51AHandshakeFieldKey := range []string{"primary", "secondary", "operation", "view", "direction"} {
		b51AHandshakeFieldInput := b51AHandshakeFieldCases(b51AHandshakeFieldT)[3]["input"].(map[string]any)
		delete(b51AHandshakeFieldInput["context"].(map[string]any), b51AHandshakeFieldKey)
		if b51AHandshakeFieldEvaluateTest(b51AHandshakeFieldT, b51AHandshakeFieldInput)["state"] != "Unknown" {
			b51AHandshakeFieldT.Fatal("missing context admitted", b51AHandshakeFieldKey)
		}
	}
	for _, b51AHandshakeFieldKey := range []string{"family", "variant"} {
		for _, b51AHandshakeFieldValue := range []any{nil, false, float64(0), "", []any{}} {
			b51AHandshakeFieldInput := b51AHandshakeFieldCases(b51AHandshakeFieldT)[0]["input"].(map[string]any)
			b51AHandshakeFieldInput["context"].(map[string]any)[b51AHandshakeFieldKey] = b51AHandshakeFieldValue
			b51AHandshakeFieldOutput := b51AHandshakeFieldEvaluateTest(b51AHandshakeFieldT, b51AHandshakeFieldInput)
			if !reflect.DeepEqual(b51AHandshakeFieldOutput["context"], b51AHandshakeFieldInput["context"]) || b51AHandshakeFieldOutput["qualified"] != false {
				b51AHandshakeFieldT.Fatal("opaque value")
			}
		}
	}
}
