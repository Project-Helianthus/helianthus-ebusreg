package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
)

func b523ActuatorFieldRead(b523ActuatorFieldDecoder *json.Decoder) (any, error) {
	b523ActuatorFieldToken, b523ActuatorFieldError := b523ActuatorFieldDecoder.Token()
	if b523ActuatorFieldError != nil {
		return nil, b523ActuatorFieldError
	}
	if b523ActuatorFieldDelimiter, b523ActuatorFieldOK := b523ActuatorFieldToken.(json.Delim); b523ActuatorFieldOK {
		switch b523ActuatorFieldDelimiter {
		case '{':
			b523ActuatorFieldObject := map[string]any{}
			for b523ActuatorFieldDecoder.More() {
				b523ActuatorFieldKeyToken, b523ActuatorFieldError := b523ActuatorFieldDecoder.Token()
				if b523ActuatorFieldError != nil {
					return nil, b523ActuatorFieldError
				}
				b523ActuatorFieldKey := b523ActuatorFieldKeyToken.(string)
				if _, b523ActuatorFieldExists := b523ActuatorFieldObject[b523ActuatorFieldKey]; b523ActuatorFieldExists {
					return nil, fmt.Errorf("duplicate key %q", b523ActuatorFieldKey)
				}
				b523ActuatorFieldValue, b523ActuatorFieldError := b523ActuatorFieldRead(b523ActuatorFieldDecoder)
				if b523ActuatorFieldError != nil {
					return nil, b523ActuatorFieldError
				}
				b523ActuatorFieldObject[b523ActuatorFieldKey] = b523ActuatorFieldValue
			}
			_, b523ActuatorFieldError = b523ActuatorFieldDecoder.Token()
			return b523ActuatorFieldObject, b523ActuatorFieldError
		case '[':
			b523ActuatorFieldArray := []any{}
			for b523ActuatorFieldDecoder.More() {
				b523ActuatorFieldValue, b523ActuatorFieldError := b523ActuatorFieldRead(b523ActuatorFieldDecoder)
				if b523ActuatorFieldError != nil {
					return nil, b523ActuatorFieldError
				}
				b523ActuatorFieldArray = append(b523ActuatorFieldArray, b523ActuatorFieldValue)
			}
			_, b523ActuatorFieldError = b523ActuatorFieldDecoder.Token()
			return b523ActuatorFieldArray, b523ActuatorFieldError
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	return b523ActuatorFieldToken, nil
}

func b523ActuatorFieldByte(b523ActuatorFieldValue any) (int, bool) {
	b523ActuatorFieldNumber, b523ActuatorFieldOK := b523ActuatorFieldValue.(json.Number)
	if !b523ActuatorFieldOK {
		return 0, false
	}
	b523ActuatorFieldRational, b523ActuatorFieldOK := new(big.Rat).SetString(string(b523ActuatorFieldNumber))
	if !b523ActuatorFieldOK || !b523ActuatorFieldRational.IsInt() || b523ActuatorFieldRational.Sign() < 0 || b523ActuatorFieldRational.Cmp(big.NewRat(255, 1)) > 0 {
		return 0, false
	}
	return int(b523ActuatorFieldRational.Num().Int64()), true
}

func b523ActuatorFieldResult(b523ActuatorFieldRaw any) map[string]any {
	return map[string]any{"raw": b523ActuatorFieldRaw, "state": "complete", "grade": "Hypothesis"}
}

func B523ActuatorFieldEvaluate(b523ActuatorFieldInput json.RawMessage) (json.RawMessage, error) {
	b523ActuatorFieldDecoder := json.NewDecoder(bytes.NewReader(b523ActuatorFieldInput))
	b523ActuatorFieldDecoder.UseNumber()
	b523ActuatorFieldValue, b523ActuatorFieldError := b523ActuatorFieldRead(b523ActuatorFieldDecoder)
	if b523ActuatorFieldError != nil {
		return nil, b523ActuatorFieldError
	}
	if _, b523ActuatorFieldError = b523ActuatorFieldDecoder.Token(); b523ActuatorFieldError != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b523ActuatorFieldObject, b523ActuatorFieldOK := b523ActuatorFieldValue.(map[string]any)
	if !b523ActuatorFieldOK {
		return nil, fmt.Errorf("input must be an object")
	}
	b523ActuatorFieldContext, b523ActuatorFieldOK := b523ActuatorFieldObject["context"].(map[string]any)
	if !b523ActuatorFieldOK {
		return nil, fmt.Errorf("context must be an object")
	}
	b523ActuatorFieldRaw, b523ActuatorFieldOK := b523ActuatorFieldObject["fields"].(map[string]any)
	if !b523ActuatorFieldOK {
		return nil, fmt.Errorf("fields must be an object")
	}
	for _, b523ActuatorFieldKey := range []string{"primary", "secondary", "operation"} {
		if b523ActuatorFieldValue := b523ActuatorFieldContext[b523ActuatorFieldKey]; b523ActuatorFieldValue != nil {
			if _, b523ActuatorFieldOK := b523ActuatorFieldByte(b523ActuatorFieldValue); !b523ActuatorFieldOK {
				return nil, fmt.Errorf("invalid context %s", b523ActuatorFieldKey)
			}
		}
	}
	for _, b523ActuatorFieldKey := range []string{"direction", "view"} {
		if b523ActuatorFieldValue := b523ActuatorFieldContext[b523ActuatorFieldKey]; b523ActuatorFieldValue != nil {
			if _, b523ActuatorFieldOK := b523ActuatorFieldValue.(string); !b523ActuatorFieldOK {
				return nil, fmt.Errorf("invalid context %s", b523ActuatorFieldKey)
			}
		}
	}
	for b523ActuatorFieldKey, b523ActuatorFieldValue := range b523ActuatorFieldRaw {
		b523ActuatorFieldArray, b523ActuatorFieldOK := b523ActuatorFieldValue.([]any)
		if !b523ActuatorFieldOK {
			return nil, fmt.Errorf("field %s must be a byte array", b523ActuatorFieldKey)
		}
		for _, b523ActuatorFieldValue := range b523ActuatorFieldArray {
			if _, b523ActuatorFieldOK := b523ActuatorFieldByte(b523ActuatorFieldValue); !b523ActuatorFieldOK {
				return nil, fmt.Errorf("invalid byte in %s", b523ActuatorFieldKey)
			}
		}
	}
	b523ActuatorFieldArguments := map[string]any{}
	if b523ActuatorFieldValue, b523ActuatorFieldPresent := b523ActuatorFieldObject["arguments"]; b523ActuatorFieldPresent {
		b523ActuatorFieldArguments, b523ActuatorFieldOK = b523ActuatorFieldValue.(map[string]any)
		if !b523ActuatorFieldOK {
			return nil, fmt.Errorf("arguments must be an object")
		}
	}
	b523ActuatorFieldFields := map[string]any{}
	b523ActuatorFieldOutput := map[string]any{"context": b523ActuatorFieldContext, "raw": b523ActuatorFieldRaw, "fields": b523ActuatorFieldFields, "state": "Unknown", "physical": "Unknown", "qualified": false}
	if _, b523ActuatorFieldPresent := b523ActuatorFieldObject["arguments"]; b523ActuatorFieldPresent {
		b523ActuatorFieldOutput["arguments"] = b523ActuatorFieldArguments
	}
	b523ActuatorFieldPrimary, _ := b523ActuatorFieldByte(b523ActuatorFieldContext["primary"])
	b523ActuatorFieldSecondary, _ := b523ActuatorFieldByte(b523ActuatorFieldContext["secondary"])
	b523ActuatorFieldOperation, _ := b523ActuatorFieldByte(b523ActuatorFieldContext["operation"])
	b523ActuatorFieldView, _ := b523ActuatorFieldContext["view"].(string)
	b523ActuatorFieldDirection, _ := b523ActuatorFieldContext["direction"].(string)
	b523ActuatorFieldSupported := b523ActuatorFieldPrimary == 181 && b523ActuatorFieldSecondary == 35 && ((b523ActuatorFieldOperation == 1 && b523ActuatorFieldView == "relay-six") || (b523ActuatorFieldOperation == 5 && b523ActuatorFieldView == "relay-twelve") || (b523ActuatorFieldOperation == 2 && b523ActuatorFieldView == "mixer")) && (b523ActuatorFieldDirection == "request" || b523ActuatorFieldDirection == "positive-feedback")
	if !b523ActuatorFieldSupported {
		return json.Marshal(b523ActuatorFieldOutput)
	}
	b523ActuatorFieldVariant := ""
	if b523ActuatorFieldOperation == 2 && b523ActuatorFieldDirection == "request" {
		b523ActuatorFieldVariant, _ = b523ActuatorFieldArguments["variant"].(string)
		if b523ActuatorFieldVariant != "two-circuit" && b523ActuatorFieldVariant != "three-circuit" {
			return nil, fmt.Errorf("arguments.variant must be two-circuit or three-circuit")
		}
	}
	b523ActuatorFieldOutput["state"] = "complete"
	b523ActuatorFieldGet := func(b523ActuatorFieldKey string, b523ActuatorFieldLength int) ([]any, map[string]any, bool) {
		b523ActuatorFieldValue, b523ActuatorFieldPresent := b523ActuatorFieldRaw[b523ActuatorFieldKey]
		b523ActuatorFieldArray, _ := b523ActuatorFieldValue.([]any)
		b523ActuatorFieldField := b523ActuatorFieldResult(b523ActuatorFieldValue)
		b523ActuatorFieldField["present"] = b523ActuatorFieldPresent
		b523ActuatorFieldFields[b523ActuatorFieldKey] = b523ActuatorFieldField
		if !b523ActuatorFieldPresent || len(b523ActuatorFieldArray) != b523ActuatorFieldLength {
			b523ActuatorFieldField["state"] = "incomplete"
			b523ActuatorFieldOutput["state"] = "partial"
			return b523ActuatorFieldArray, b523ActuatorFieldField, false
		}
		return b523ActuatorFieldArray, b523ActuatorFieldField, true
	}
	if b523ActuatorFieldDirection == "positive-feedback" {
		if b523ActuatorFieldArray, b523ActuatorFieldField, b523ActuatorFieldOK := b523ActuatorFieldGet("feedback", 1); b523ActuatorFieldOK {
			b523ActuatorFieldByteValue, _ := b523ActuatorFieldByte(b523ActuatorFieldArray[0])
			b523ActuatorFieldField["label"] = "rejected"
			if b523ActuatorFieldOperation == 5 {
				b523ActuatorFieldField["label"] = "rejected-or-unsupported"
			}
			if b523ActuatorFieldByteValue == 1 {
				b523ActuatorFieldField["label"] = "accepted"
			}
		}
		if b523ActuatorFieldOperation == 2 {
			if b523ActuatorFieldArray, b523ActuatorFieldField, b523ActuatorFieldOK := b523ActuatorFieldGet("movement", 1); b523ActuatorFieldOK {
				b523ActuatorFieldByteValue, _ := b523ActuatorFieldByte(b523ActuatorFieldArray[0])
				b523ActuatorFieldSigned := int(int8(b523ActuatorFieldByteValue))
				if b523ActuatorFieldSigned >= -100 && b523ActuatorFieldSigned <= 100 {
					b523ActuatorFieldField["number"] = b523ActuatorFieldSigned
					b523ActuatorFieldField["unit"] = "percent"
					b523ActuatorFieldField["calculation"] = "Proven"
					b523ActuatorFieldField["maximum"] = b523ActuatorFieldSigned == -100 || b523ActuatorFieldSigned == 100
					b523ActuatorFieldField["label"] = "stopped"
					if b523ActuatorFieldSigned < 0 {
						b523ActuatorFieldField["label"] = "closing"
					}
					if b523ActuatorFieldSigned > 0 {
						b523ActuatorFieldField["label"] = "opening"
					}
				} else {
					b523ActuatorFieldField["state"] = "Unknown"
				}
			}
		}
	} else if b523ActuatorFieldOperation == 1 || b523ActuatorFieldOperation == 5 {
		b523ActuatorFieldLength := 6
		if b523ActuatorFieldOperation == 5 {
			b523ActuatorFieldLength = 12
		}
		if b523ActuatorFieldArray, b523ActuatorFieldField, b523ActuatorFieldOK := b523ActuatorFieldGet("relays", b523ActuatorFieldLength); b523ActuatorFieldOK {
			b523ActuatorFieldRelays := []any{}
			for _, b523ActuatorFieldValue := range b523ActuatorFieldArray {
				b523ActuatorFieldByteValue, _ := b523ActuatorFieldByte(b523ActuatorFieldValue)
				b523ActuatorFieldRelay := b523ActuatorFieldResult(b523ActuatorFieldValue)
				switch b523ActuatorFieldByteValue {
				case 0:
					b523ActuatorFieldRelay["label"] = "off"
				case 254:
					b523ActuatorFieldRelay["label"] = "on"
				case 255:
					b523ActuatorFieldRelay["label"] = "unchanged"
				default:
					b523ActuatorFieldRelay["label"] = "on-for-seconds-then-off"
					b523ActuatorFieldRelay["number"] = b523ActuatorFieldByteValue
					b523ActuatorFieldRelay["unit"] = "seconds"
				}
				b523ActuatorFieldRelays = append(b523ActuatorFieldRelays, b523ActuatorFieldRelay)
			}
			b523ActuatorFieldField["relays"] = b523ActuatorFieldRelays
		}
		b523ActuatorFieldPWMKeys := []string{"pwm1"}
		if b523ActuatorFieldOperation == 5 {
			b523ActuatorFieldPWMKeys = append(b523ActuatorFieldPWMKeys, "pwm2")
		}
		for _, b523ActuatorFieldKey := range b523ActuatorFieldPWMKeys {
			if b523ActuatorFieldArray, b523ActuatorFieldField, b523ActuatorFieldOK := b523ActuatorFieldGet(b523ActuatorFieldKey, 1); b523ActuatorFieldOK {
				b523ActuatorFieldByteValue, _ := b523ActuatorFieldByte(b523ActuatorFieldArray[0])
				b523ActuatorFieldNumeric := b523ActuatorFieldResult(b523ActuatorFieldArray)
				b523ActuatorFieldNumeric["candidate"] = "numeric"
				b523ActuatorFieldNumeric["number"] = float64(b523ActuatorFieldByteValue) / 2
				b523ActuatorFieldNumeric["unit"] = "percent"
				b523ActuatorFieldNumeric["candidate_domain"] = b523ActuatorFieldByteValue <= 250
				b523ActuatorFieldNumeric["calculation"] = "Proven"
				b523ActuatorFieldBehavior := b523ActuatorFieldResult(b523ActuatorFieldArray)
				b523ActuatorFieldBehavior["candidate"] = "behavior"
				if b523ActuatorFieldByteValue <= 200 {
					b523ActuatorFieldBehavior["number"] = float64(b523ActuatorFieldByteValue) / 2
					b523ActuatorFieldBehavior["unit"] = "percent"
					b523ActuatorFieldBehavior["calculation"] = "Proven"
				} else {
					b523ActuatorFieldBehavior["label"] = "unchanged"
				}
				b523ActuatorFieldField["alternatives"] = []any{b523ActuatorFieldNumeric, b523ActuatorFieldBehavior}
			}
		}
	} else {
		if b523ActuatorFieldArray, b523ActuatorFieldField, b523ActuatorFieldOK := b523ActuatorFieldGet("circuit", 1); b523ActuatorFieldOK {
			b523ActuatorFieldByteValue, _ := b523ActuatorFieldByte(b523ActuatorFieldArray[0])
			if b523ActuatorFieldByteValue < 2 || (b523ActuatorFieldByteValue == 2 && b523ActuatorFieldVariant == "three-circuit") {
				b523ActuatorFieldField["label"] = fmt.Sprintf("circuit-%d", b523ActuatorFieldByteValue)
			} else {
				b523ActuatorFieldField["label"] = "ignored"
			}
		}
		b523ActuatorFieldMode := -1
		if b523ActuatorFieldArray, b523ActuatorFieldField, b523ActuatorFieldOK := b523ActuatorFieldGet("mode", 1); b523ActuatorFieldOK {
			b523ActuatorFieldMode, _ = b523ActuatorFieldByte(b523ActuatorFieldArray[0])
			b523ActuatorFieldLabels := []string{"stop/pump-off/close", "heating-regulation", "cooling-regulation", "pump-off/open", "pump-on/close", "pump-on/open"}
			b523ActuatorFieldField["label"] = "ignored"
			if b523ActuatorFieldMode < len(b523ActuatorFieldLabels) {
				b523ActuatorFieldField["label"] = b523ActuatorFieldLabels[b523ActuatorFieldMode]
			}
		}
		if b523ActuatorFieldArray, b523ActuatorFieldField, b523ActuatorFieldOK := b523ActuatorFieldGet("setpoint", 1); b523ActuatorFieldOK {
			if b523ActuatorFieldMode >= 0 {
				b523ActuatorFieldField["relevant"] = b523ActuatorFieldMode == 1 || b523ActuatorFieldMode == 2
			} else {
				b523ActuatorFieldField["relevance"] = "Unknown"
			}
			b523ActuatorFieldByteValue, _ := b523ActuatorFieldByte(b523ActuatorFieldArray[0])
			b523ActuatorFieldSigned := b523ActuatorFieldResult(b523ActuatorFieldArray)
			b523ActuatorFieldSigned["candidate"] = "signed-half"
			b523ActuatorFieldSigned["number"] = float64(int8(b523ActuatorFieldByteValue)) / 2
			b523ActuatorFieldSigned["unit"] = "Celsius"
			b523ActuatorFieldSigned["calculation"] = "Proven"
			b523ActuatorFieldExample := b523ActuatorFieldResult(b523ActuatorFieldArray)
			b523ActuatorFieldExample["candidate"] = "conflicting-example"
			b523ActuatorFieldExample["clamp"] = 95
			b523ActuatorFieldExample["state"] = "Unknown"
			if b523ActuatorFieldByteValue == 190 {
				b523ActuatorFieldExample["number"] = 95
				b523ActuatorFieldExample["unit"] = "Celsius"
				b523ActuatorFieldExample["state"] = "complete"
			}
			b523ActuatorFieldField["alternatives"] = []any{b523ActuatorFieldSigned, b523ActuatorFieldExample}
		}
	}
	return json.Marshal(b523ActuatorFieldOutput)
}
