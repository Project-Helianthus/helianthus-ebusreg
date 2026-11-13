package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func b512PowerStepperFieldRead(b512PowerStepperFieldDecoder *json.Decoder) (any, error) {
	b512PowerStepperFieldToken, b512PowerStepperFieldErr := b512PowerStepperFieldDecoder.Token()
	if b512PowerStepperFieldErr != nil {
		return nil, b512PowerStepperFieldErr
	}
	if b512PowerStepperFieldDelimiter, b512PowerStepperFieldOK := b512PowerStepperFieldToken.(json.Delim); b512PowerStepperFieldOK {
		switch b512PowerStepperFieldDelimiter {
		case '{':
			b512PowerStepperFieldObject := map[string]any{}
			for b512PowerStepperFieldDecoder.More() {
				b512PowerStepperFieldKeyToken, b512PowerStepperFieldErr := b512PowerStepperFieldDecoder.Token()
				if b512PowerStepperFieldErr != nil {
					return nil, b512PowerStepperFieldErr
				}
				b512PowerStepperFieldKey := b512PowerStepperFieldKeyToken.(string)
				if _, b512PowerStepperFieldExists := b512PowerStepperFieldObject[b512PowerStepperFieldKey]; b512PowerStepperFieldExists {
					return nil, fmt.Errorf("duplicate key %q", b512PowerStepperFieldKey)
				}
				b512PowerStepperFieldValue, b512PowerStepperFieldErr := b512PowerStepperFieldRead(b512PowerStepperFieldDecoder)
				if b512PowerStepperFieldErr != nil {
					return nil, b512PowerStepperFieldErr
				}
				b512PowerStepperFieldObject[b512PowerStepperFieldKey] = b512PowerStepperFieldValue
			}
			_, b512PowerStepperFieldErr := b512PowerStepperFieldDecoder.Token()
			return b512PowerStepperFieldObject, b512PowerStepperFieldErr
		case '[':
			b512PowerStepperFieldArray := []any{}
			for b512PowerStepperFieldDecoder.More() {
				b512PowerStepperFieldValue, b512PowerStepperFieldErr := b512PowerStepperFieldRead(b512PowerStepperFieldDecoder)
				if b512PowerStepperFieldErr != nil {
					return nil, b512PowerStepperFieldErr
				}
				b512PowerStepperFieldArray = append(b512PowerStepperFieldArray, b512PowerStepperFieldValue)
			}
			_, b512PowerStepperFieldErr := b512PowerStepperFieldDecoder.Token()
			return b512PowerStepperFieldArray, b512PowerStepperFieldErr
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	if b512PowerStepperFieldNumber, b512PowerStepperFieldOK := b512PowerStepperFieldToken.(json.Number); b512PowerStepperFieldOK {
		if _, b512PowerStepperFieldErr := b512PowerStepperFieldNumber.Float64(); b512PowerStepperFieldErr != nil {
			return nil, fmt.Errorf("nonfinite number")
		}
	}
	return b512PowerStepperFieldToken, nil
}

func B512PowerStepperFieldEvaluate(b512PowerStepperFieldInput json.RawMessage) (json.RawMessage, error) {
	b512PowerStepperFieldDecoder := json.NewDecoder(bytes.NewReader(b512PowerStepperFieldInput))
	b512PowerStepperFieldDecoder.UseNumber()
	b512PowerStepperFieldValue, b512PowerStepperFieldErr := b512PowerStepperFieldRead(b512PowerStepperFieldDecoder)
	if b512PowerStepperFieldErr != nil {
		return nil, b512PowerStepperFieldErr
	}
	if _, b512PowerStepperFieldErr = b512PowerStepperFieldDecoder.Token(); b512PowerStepperFieldErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b512PowerStepperFieldRoot, b512PowerStepperFieldOK := b512PowerStepperFieldValue.(map[string]any)
	if !b512PowerStepperFieldOK {
		return nil, fmt.Errorf("input must be object")
	}
	b512PowerStepperFieldContext, b512PowerStepperFieldOK := b512PowerStepperFieldRoot["context"].(map[string]any)
	if !b512PowerStepperFieldOK {
		return nil, fmt.Errorf("context must be object")
	}
	for _, b512PowerStepperFieldKey := range []string{"primary", "secondary", "operation"} {
		if b512PowerStepperFieldValue, b512PowerStepperFieldExists := b512PowerStepperFieldContext[b512PowerStepperFieldKey]; b512PowerStepperFieldExists {
			b512PowerStepperFieldNumber, b512PowerStepperFieldOK := b512PowerStepperFieldValue.(json.Number)
			if !b512PowerStepperFieldOK {
				return nil, fmt.Errorf("context number must be integer")
			}
			b512PowerStepperFieldInteger, b512PowerStepperFieldErr := b512PowerStepperFieldNumber.Int64()
			if b512PowerStepperFieldErr != nil || b512PowerStepperFieldInteger < 0 || b512PowerStepperFieldInteger > 255 {
				return nil, fmt.Errorf("context number out of range")
			}
		}
	}
	for _, b512PowerStepperFieldKey := range []string{"direction", "view"} {
		if b512PowerStepperFieldValue, b512PowerStepperFieldExists := b512PowerStepperFieldContext[b512PowerStepperFieldKey]; b512PowerStepperFieldExists {
			if _, b512PowerStepperFieldOK := b512PowerStepperFieldValue.(string); !b512PowerStepperFieldOK {
				return nil, fmt.Errorf("context selector must be string")
			}
		}
	}
	b512PowerStepperFieldRaw, b512PowerStepperFieldOK := b512PowerStepperFieldRoot["fields"].(map[string]any)
	if !b512PowerStepperFieldOK {
		return nil, fmt.Errorf("fields must be object")
	}
	if b512PowerStepperFieldArguments, b512PowerStepperFieldExists := b512PowerStepperFieldRoot["arguments"]; b512PowerStepperFieldExists {
		if _, b512PowerStepperFieldOK := b512PowerStepperFieldArguments.(map[string]any); !b512PowerStepperFieldOK {
			return nil, fmt.Errorf("arguments must be object")
		}
	}
	for _, b512PowerStepperFieldValue := range b512PowerStepperFieldRaw {
		b512PowerStepperFieldArray, b512PowerStepperFieldOK := b512PowerStepperFieldValue.([]any)
		if !b512PowerStepperFieldOK {
			return nil, fmt.Errorf("field must be byte array")
		}
		for _, b512PowerStepperFieldByte := range b512PowerStepperFieldArray {
			b512PowerStepperFieldNumber, b512PowerStepperFieldOK := b512PowerStepperFieldByte.(json.Number)
			if !b512PowerStepperFieldOK {
				return nil, fmt.Errorf("byte must be integer")
			}
			b512PowerStepperFieldInteger, b512PowerStepperFieldErr := b512PowerStepperFieldNumber.Int64()
			if b512PowerStepperFieldErr != nil || b512PowerStepperFieldInteger < 0 || b512PowerStepperFieldInteger > 255 {
				return nil, fmt.Errorf("byte out of range")
			}
		}
	}
	b512PowerStepperFieldResults := map[string]any{}
	for b512PowerStepperFieldName, b512PowerStepperFieldValue := range b512PowerStepperFieldRaw {
		b512PowerStepperFieldResults[b512PowerStepperFieldName] = map[string]any{"raw": b512PowerStepperFieldValue, "state": "Unknown"}
	}
	b512PowerStepperFieldRoot["results"] = b512PowerStepperFieldResults
	b512PowerStepperFieldRoot["state"] = "Unknown"
	b512PowerStepperFieldRoot["physical"] = "Unknown"
	b512PowerStepperFieldRoot["qualified"] = false
	b512PowerStepperFieldOperation, _ := b512PowerStepperFieldContext["operation"].(json.Number)
	b512PowerStepperFieldOp, _ := b512PowerStepperFieldOperation.Int64()
	b512PowerStepperFieldPrimary, _ := b512PowerStepperFieldContext["primary"].(json.Number)
	b512PowerStepperFieldSecondary, _ := b512PowerStepperFieldContext["secondary"].(json.Number)
	b512PowerStepperFieldP, _ := b512PowerStepperFieldPrimary.Int64()
	b512PowerStepperFieldS, _ := b512PowerStepperFieldSecondary.Int64()
	b512PowerStepperFieldDirection := b512PowerStepperFieldContext["direction"]
	if b512PowerStepperFieldP != 181 || b512PowerStepperFieldS != 18 || b512PowerStepperFieldOp < 7 || b512PowerStepperFieldOp > 9 || b512PowerStepperFieldContext["view"] != "power-stepper" || (b512PowerStepperFieldDirection != "request" && b512PowerStepperFieldDirection != "positive") {
		return json.Marshal(b512PowerStepperFieldRoot)
	}
	b512PowerStepperFieldLengths := map[string]int{}
	switch b512PowerStepperFieldOp {
	case 7:
		if b512PowerStepperFieldDirection == "positive" {
			b512PowerStepperFieldLengths["maximum_power"] = 1
			b512PowerStepperFieldLengths["minimum_power"] = 1
		}
	case 8:
		b512PowerStepperFieldLengths["stepper"] = 1
		if b512PowerStepperFieldDirection == "request" {
			b512PowerStepperFieldLengths["inlet"] = 2
			b512PowerStepperFieldLengths["status"] = 1
		} else {
			b512PowerStepperFieldLengths["setpoint"] = 2
		}
	case 9:
		b512PowerStepperFieldLengths["release"] = 1
	}
	b512PowerStepperFieldRoot["state"] = "complete"
	b512PowerStepperFieldRoot["grade"] = "Hypothesis"
	for b512PowerStepperFieldName, b512PowerStepperFieldWidth := range b512PowerStepperFieldLengths {
		b512PowerStepperFieldResult := map[string]any{"state": "incomplete"}
		b512PowerStepperFieldResults[b512PowerStepperFieldName] = b512PowerStepperFieldResult
		b512PowerStepperFieldValue, b512PowerStepperFieldPresent := b512PowerStepperFieldRaw[b512PowerStepperFieldName]
		if b512PowerStepperFieldPresent {
			b512PowerStepperFieldResult["raw"] = b512PowerStepperFieldValue
		}
		b512PowerStepperFieldArray, _ := b512PowerStepperFieldValue.([]any)
		if len(b512PowerStepperFieldArray) != b512PowerStepperFieldWidth {
			b512PowerStepperFieldRoot["state"] = "partial"
			continue
		}
		b512PowerStepperFieldBytes := make([]int, b512PowerStepperFieldWidth)
		for b512PowerStepperFieldIndex, b512PowerStepperFieldByte := range b512PowerStepperFieldArray {
			b512PowerStepperFieldInteger, _ := b512PowerStepperFieldByte.(json.Number).Int64()
			b512PowerStepperFieldBytes[b512PowerStepperFieldIndex] = int(b512PowerStepperFieldInteger)
		}
		b512PowerStepperFieldResult["state"] = "complete"
		b512PowerStepperFieldResult["grade"] = "Hypothesis"
		switch b512PowerStepperFieldName {
		case "maximum_power":
			b512PowerStepperFieldResult["number"] = b512PowerStepperFieldBytes[0]
			b512PowerStepperFieldResult["unit"] = "kW"
			b512PowerStepperFieldResult["default"] = 0
			b512PowerStepperFieldResult["calculation"] = "Proven"
		case "minimum_power":
			b512PowerStepperFieldResult["number"] = float64(b512PowerStepperFieldBytes[0]) / 2
			b512PowerStepperFieldResult["unit"] = "%"
			b512PowerStepperFieldResult["default"] = 0
			b512PowerStepperFieldResult["candidate_domain"] = b512PowerStepperFieldBytes[0] <= 250
			b512PowerStepperFieldResult["calculation"] = "Proven"
		case "inlet", "setpoint":
			b512PowerStepperFieldResult["number"] = float64(int16(b512PowerStepperFieldBytes[0]|b512PowerStepperFieldBytes[1]<<8)) / 16
			b512PowerStepperFieldResult["unit"] = "Celsius"
			b512PowerStepperFieldResult["calculation"] = "Proven"
		case "status":
			b512PowerStepperFieldLabel := "unmapped"
			switch b512PowerStepperFieldBytes[0] {
			case 0:
				b512PowerStepperFieldLabel = "valid"
			case 0x55:
				b512PowerStepperFieldLabel = "short"
			case 0xaa:
				b512PowerStepperFieldLabel = "interruption-or-disconnected"
			}
			b512PowerStepperFieldResult["label"] = b512PowerStepperFieldLabel
		case "stepper", "release":
			b512PowerStepperFieldResult["flags"] = map[string]any{"release": b512PowerStepperFieldBytes[0]&1 != 0}
			b512PowerStepperFieldResult["reserved"] = b512PowerStepperFieldBytes[0] & 254
		}
	}
	return json.Marshal(b512PowerStepperFieldRoot)
}
