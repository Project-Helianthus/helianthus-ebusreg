package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func b516EnergyTypeFieldDecode(b516EnergyTypeFieldDecoder *json.Decoder) (any, error) {
	b516EnergyTypeFieldToken, b516EnergyTypeFieldErr := b516EnergyTypeFieldDecoder.Token()
	if b516EnergyTypeFieldErr != nil {
		return nil, b516EnergyTypeFieldErr
	}
	switch b516EnergyTypeFieldToken {
	case json.Delim('{'):
		b516EnergyTypeFieldObject := map[string]any{}
		for b516EnergyTypeFieldDecoder.More() {
			b516EnergyTypeFieldKeyToken, b516EnergyTypeFieldErr := b516EnergyTypeFieldDecoder.Token()
			if b516EnergyTypeFieldErr != nil {
				return nil, b516EnergyTypeFieldErr
			}
			b516EnergyTypeFieldKey := b516EnergyTypeFieldKeyToken.(string)
			if _, b516EnergyTypeFieldExists := b516EnergyTypeFieldObject[b516EnergyTypeFieldKey]; b516EnergyTypeFieldExists {
				return nil, fmt.Errorf("duplicate key %q", b516EnergyTypeFieldKey)
			}
			b516EnergyTypeFieldValue, b516EnergyTypeFieldErr := b516EnergyTypeFieldDecode(b516EnergyTypeFieldDecoder)
			if b516EnergyTypeFieldErr != nil {
				return nil, b516EnergyTypeFieldErr
			}
			b516EnergyTypeFieldObject[b516EnergyTypeFieldKey] = b516EnergyTypeFieldValue
		}
		_, b516EnergyTypeFieldErr := b516EnergyTypeFieldDecoder.Token()
		return b516EnergyTypeFieldObject, b516EnergyTypeFieldErr
	case json.Delim('['):
		b516EnergyTypeFieldArray := []any{}
		for b516EnergyTypeFieldDecoder.More() {
			b516EnergyTypeFieldValue, b516EnergyTypeFieldErr := b516EnergyTypeFieldDecode(b516EnergyTypeFieldDecoder)
			if b516EnergyTypeFieldErr != nil {
				return nil, b516EnergyTypeFieldErr
			}
			b516EnergyTypeFieldArray = append(b516EnergyTypeFieldArray, b516EnergyTypeFieldValue)
		}
		_, b516EnergyTypeFieldErr := b516EnergyTypeFieldDecoder.Token()
		return b516EnergyTypeFieldArray, b516EnergyTypeFieldErr
	default:
		if _, b516EnergyTypeFieldDelimiter := b516EnergyTypeFieldToken.(json.Delim); b516EnergyTypeFieldDelimiter {
			return nil, fmt.Errorf("unexpected delimiter")
		}
		return b516EnergyTypeFieldToken, nil
	}
}

func b516EnergyTypeFieldByte(b516EnergyTypeFieldValue any) (uint8, bool) {
	b516EnergyTypeFieldNumber, b516EnergyTypeFieldOK := b516EnergyTypeFieldValue.(json.Number)
	if !b516EnergyTypeFieldOK {
		return 0, false
	}
	b516EnergyTypeFieldInteger, b516EnergyTypeFieldErr := strconv.ParseUint(string(b516EnergyTypeFieldNumber), 10, 8)
	return uint8(b516EnergyTypeFieldInteger), b516EnergyTypeFieldErr == nil
}

func B516EnergyTypeFieldEvaluate(b516EnergyTypeFieldInput json.RawMessage) (json.RawMessage, error) {
	b516EnergyTypeFieldDecoder := json.NewDecoder(bytes.NewReader(b516EnergyTypeFieldInput))
	b516EnergyTypeFieldDecoder.UseNumber()
	b516EnergyTypeFieldDecoded, b516EnergyTypeFieldErr := b516EnergyTypeFieldDecode(b516EnergyTypeFieldDecoder)
	if b516EnergyTypeFieldErr != nil {
		return nil, b516EnergyTypeFieldErr
	}
	if _, b516EnergyTypeFieldErr = b516EnergyTypeFieldDecoder.Token(); b516EnergyTypeFieldErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b516EnergyTypeFieldRoot, b516EnergyTypeFieldOK := b516EnergyTypeFieldDecoded.(map[string]any)
	if !b516EnergyTypeFieldOK {
		return nil, fmt.Errorf("input must be an object")
	}
	b516EnergyTypeFieldContext, b516EnergyTypeFieldOK := b516EnergyTypeFieldRoot["context"].(map[string]any)
	if !b516EnergyTypeFieldOK {
		return nil, fmt.Errorf("context must be an object")
	}
	b516EnergyTypeFieldRaw, b516EnergyTypeFieldOK := b516EnergyTypeFieldRoot["fields"].(map[string]any)
	if !b516EnergyTypeFieldOK {
		return nil, fmt.Errorf("fields must be an object")
	}
	for _, b516EnergyTypeFieldKey := range []string{"primary", "secondary", "operation"} {
		if b516EnergyTypeFieldValue, b516EnergyTypeFieldPresent := b516EnergyTypeFieldContext[b516EnergyTypeFieldKey]; b516EnergyTypeFieldPresent && b516EnergyTypeFieldValue != nil {
			if _, b516EnergyTypeFieldOK := b516EnergyTypeFieldByte(b516EnergyTypeFieldValue); !b516EnergyTypeFieldOK {
				return nil, fmt.Errorf("%s must be a byte or null", b516EnergyTypeFieldKey)
			}
		}
	}
	for _, b516EnergyTypeFieldKey := range []string{"direction", "view"} {
		if b516EnergyTypeFieldValue, b516EnergyTypeFieldPresent := b516EnergyTypeFieldContext[b516EnergyTypeFieldKey]; b516EnergyTypeFieldPresent && b516EnergyTypeFieldValue != nil {
			if _, b516EnergyTypeFieldOK := b516EnergyTypeFieldValue.(string); !b516EnergyTypeFieldOK {
				return nil, fmt.Errorf("%s must be a string or null", b516EnergyTypeFieldKey)
			}
		}
	}
	for b516EnergyTypeFieldKey, b516EnergyTypeFieldValue := range b516EnergyTypeFieldRaw {
		b516EnergyTypeFieldArray, b516EnergyTypeFieldOK := b516EnergyTypeFieldValue.([]any)
		if !b516EnergyTypeFieldOK {
			return nil, fmt.Errorf("field %s must be a byte array", b516EnergyTypeFieldKey)
		}
		for _, b516EnergyTypeFieldElement := range b516EnergyTypeFieldArray {
			if _, b516EnergyTypeFieldOK := b516EnergyTypeFieldByte(b516EnergyTypeFieldElement); !b516EnergyTypeFieldOK {
				return nil, fmt.Errorf("field %s contains an invalid byte", b516EnergyTypeFieldKey)
			}
		}
	}
	b516EnergyTypeFieldResults := map[string]any{}
	for b516EnergyTypeFieldKey, b516EnergyTypeFieldValue := range b516EnergyTypeFieldRaw {
		b516EnergyTypeFieldResults[b516EnergyTypeFieldKey] = map[string]any{"raw": b516EnergyTypeFieldValue, "state": "Unknown", "present": true}
	}
	b516EnergyTypeFieldOutput := map[string]any{"context": b516EnergyTypeFieldContext, "raw": b516EnergyTypeFieldRaw, "state": "Unknown", "fields": b516EnergyTypeFieldResults, "physical": "Unknown", "qualified": false}
	if b516EnergyTypeFieldArguments, b516EnergyTypeFieldPresent := b516EnergyTypeFieldRoot["arguments"]; b516EnergyTypeFieldPresent {
		if _, b516EnergyTypeFieldOK := b516EnergyTypeFieldArguments.(map[string]any); !b516EnergyTypeFieldOK {
			return nil, fmt.Errorf("arguments must be an object")
		}
		b516EnergyTypeFieldOutput["arguments"] = b516EnergyTypeFieldArguments
	}
	b516EnergyTypeFieldPrimary, _ := b516EnergyTypeFieldByte(b516EnergyTypeFieldContext["primary"])
	b516EnergyTypeFieldSecondary, _ := b516EnergyTypeFieldByte(b516EnergyTypeFieldContext["secondary"])
	b516EnergyTypeFieldOperation, b516EnergyTypeFieldOperationPresent := b516EnergyTypeFieldContext["operation"]
	b516EnergyTypeFieldDirection, _ := b516EnergyTypeFieldContext["direction"].(string)
	if b516EnergyTypeFieldPrimary == 181 && b516EnergyTypeFieldSecondary == 22 && b516EnergyTypeFieldOperationPresent && b516EnergyTypeFieldOperation == nil && b516EnergyTypeFieldContext["view"] == "energy-type" && (b516EnergyTypeFieldDirection == "request" || b516EnergyTypeFieldDirection == "read-positive" || b516EnergyTypeFieldDirection == "write-positive") {
		b516EnergyTypeFieldValue, b516EnergyTypeFieldPresent := b516EnergyTypeFieldRaw["value"]
		b516EnergyTypeFieldResult := map[string]any{"raw": b516EnergyTypeFieldValue, "present": b516EnergyTypeFieldPresent, "state": "incomplete"}
		b516EnergyTypeFieldResults["value"] = b516EnergyTypeFieldResult
		b516EnergyTypeFieldOutput["state"] = "partial"
		if b516EnergyTypeFieldArray, b516EnergyTypeFieldOK := b516EnergyTypeFieldValue.([]any); b516EnergyTypeFieldOK && len(b516EnergyTypeFieldArray) == 1 {
			b516EnergyTypeFieldByteValue, _ := b516EnergyTypeFieldByte(b516EnergyTypeFieldArray[0])
			b516EnergyTypeFieldOutput["state"] = "complete"
			b516EnergyTypeFieldResult["state"] = "Hypothesis"
			if b516EnergyTypeFieldDirection == "request" {
				b516EnergyTypeFieldResult["number"] = b516EnergyTypeFieldByteValue
				b516EnergyTypeFieldResult["unmapped"] = true
			} else {
				b516EnergyTypeFieldResult["low_six"] = b516EnergyTypeFieldByteValue & 63
				b516EnergyTypeFieldResult["calculation"] = "Proven"
				if b516EnergyTypeFieldDirection == "read-positive" {
					b516EnergyTypeFieldResult["flags"] = map[string]any{"flag6": b516EnergyTypeFieldByteValue&64 != 0, "flag7": b516EnergyTypeFieldByteValue&128 != 0}
				} else {
					b516EnergyTypeFieldResult["unmodeled_high"] = b516EnergyTypeFieldByteValue & 192
				}
			}
		}
	}
	return json.Marshal(b516EnergyTypeFieldOutput)
}
