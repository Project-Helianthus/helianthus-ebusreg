package vaillant

import (
	b515PeriodFieldBytes "bytes"
	b515PeriodFieldJSON "encoding/json"
	b515PeriodFieldFmt "fmt"
	b515PeriodFieldIO "io"
	b515PeriodFieldMath "math"
	b515PeriodFieldStrconv "strconv"
)

func b515PeriodFieldDecode(b515PeriodFieldDecoder *b515PeriodFieldJSON.Decoder) (any, error) {
	b515PeriodFieldToken, b515PeriodFieldError := b515PeriodFieldDecoder.Token()
	if b515PeriodFieldError != nil {
		return nil, b515PeriodFieldError
	}
	if b515PeriodFieldDelimiter, b515PeriodFieldOK := b515PeriodFieldToken.(b515PeriodFieldJSON.Delim); b515PeriodFieldOK {
		switch b515PeriodFieldDelimiter {
		case '{':
			b515PeriodFieldObject := map[string]any{}
			for b515PeriodFieldDecoder.More() {
				b515PeriodFieldKeyToken, b515PeriodFieldError := b515PeriodFieldDecoder.Token()
				if b515PeriodFieldError != nil {
					return nil, b515PeriodFieldError
				}
				b515PeriodFieldKey, b515PeriodFieldOK := b515PeriodFieldKeyToken.(string)
				if !b515PeriodFieldOK {
					return nil, b515PeriodFieldFmt.Errorf("object key required")
				}
				if _, b515PeriodFieldPresent := b515PeriodFieldObject[b515PeriodFieldKey]; b515PeriodFieldPresent {
					return nil, b515PeriodFieldFmt.Errorf("duplicate key %q", b515PeriodFieldKey)
				}
				b515PeriodFieldValue, b515PeriodFieldError := b515PeriodFieldDecode(b515PeriodFieldDecoder)
				if b515PeriodFieldError != nil {
					return nil, b515PeriodFieldError
				}
				b515PeriodFieldObject[b515PeriodFieldKey] = b515PeriodFieldValue
			}
			_, b515PeriodFieldError = b515PeriodFieldDecoder.Token()
			return b515PeriodFieldObject, b515PeriodFieldError
		case '[':
			b515PeriodFieldArray := []any{}
			for b515PeriodFieldDecoder.More() {
				b515PeriodFieldValue, b515PeriodFieldError := b515PeriodFieldDecode(b515PeriodFieldDecoder)
				if b515PeriodFieldError != nil {
					return nil, b515PeriodFieldError
				}
				b515PeriodFieldArray = append(b515PeriodFieldArray, b515PeriodFieldValue)
			}
			_, b515PeriodFieldError = b515PeriodFieldDecoder.Token()
			return b515PeriodFieldArray, b515PeriodFieldError
		}
		return nil, b515PeriodFieldFmt.Errorf("unexpected delimiter")
	}
	if b515PeriodFieldNumber, b515PeriodFieldOK := b515PeriodFieldToken.(b515PeriodFieldJSON.Number); b515PeriodFieldOK {
		b515PeriodFieldFloat, b515PeriodFieldError := b515PeriodFieldStrconv.ParseFloat(string(b515PeriodFieldNumber), 64)
		if b515PeriodFieldError != nil || b515PeriodFieldMath.IsInf(b515PeriodFieldFloat, 0) || b515PeriodFieldMath.IsNaN(b515PeriodFieldFloat) {
			return nil, b515PeriodFieldFmt.Errorf("nonfinite number")
		}
	}
	return b515PeriodFieldToken, nil
}

func b515PeriodFieldInteger(b515PeriodFieldValue any, b515PeriodFieldMaximum int) (int, error) {
	b515PeriodFieldNumber, b515PeriodFieldOK := b515PeriodFieldValue.(b515PeriodFieldJSON.Number)
	if !b515PeriodFieldOK {
		return 0, b515PeriodFieldFmt.Errorf("integer required")
	}
	b515PeriodFieldInteger, b515PeriodFieldError := b515PeriodFieldStrconv.ParseUint(string(b515PeriodFieldNumber), 10, 16)
	if b515PeriodFieldError != nil || b515PeriodFieldInteger > uint64(b515PeriodFieldMaximum) {
		return 0, b515PeriodFieldFmt.Errorf("integer out of range or nonintegral")
	}
	return int(b515PeriodFieldInteger), nil
}

func b515PeriodFieldPeriod(b515PeriodFieldRaw any) map[string]any {
	b515PeriodFieldResult := map[string]any{"raw": b515PeriodFieldRaw, "state": "Unknown"}
	b515PeriodFieldNumber, _ := b515PeriodFieldInteger(b515PeriodFieldRaw, 255)
	if b515PeriodFieldNumber <= 144 {
		b515PeriodFieldMinutes := b515PeriodFieldNumber * 10
		b515PeriodFieldResult["state"] = "Hypothesis"
		b515PeriodFieldResult["calculation"] = "Proven"
		b515PeriodFieldResult["number"] = b515PeriodFieldMinutes
		b515PeriodFieldResult["hour"] = b515PeriodFieldMinutes / 60
		b515PeriodFieldResult["minute"] = b515PeriodFieldMinutes % 60
		b515PeriodFieldResult["label"] = b515PeriodFieldFmt.Sprintf("%02d:%02d", b515PeriodFieldMinutes/60, b515PeriodFieldMinutes%60)
	}
	return b515PeriodFieldResult
}

func B515PeriodFieldEvaluate(b515PeriodFieldInput b515PeriodFieldJSON.RawMessage) (b515PeriodFieldJSON.RawMessage, error) {
	b515PeriodFieldDecoder := b515PeriodFieldJSON.NewDecoder(b515PeriodFieldBytes.NewReader(b515PeriodFieldInput))
	b515PeriodFieldDecoder.UseNumber()
	b515PeriodFieldValue, b515PeriodFieldError := b515PeriodFieldDecode(b515PeriodFieldDecoder)
	if b515PeriodFieldError != nil {
		return nil, b515PeriodFieldError
	}
	if _, b515PeriodFieldError = b515PeriodFieldDecoder.Token(); b515PeriodFieldError != b515PeriodFieldIO.EOF {
		return nil, b515PeriodFieldFmt.Errorf("trailing JSON")
	}
	b515PeriodFieldObject, b515PeriodFieldOK := b515PeriodFieldValue.(map[string]any)
	if !b515PeriodFieldOK {
		return nil, b515PeriodFieldFmt.Errorf("input object required")
	}
	b515PeriodFieldContext, b515PeriodFieldOK := b515PeriodFieldObject["context"].(map[string]any)
	if !b515PeriodFieldOK {
		return nil, b515PeriodFieldFmt.Errorf("context object required")
	}
	b515PeriodFieldFields, b515PeriodFieldOK := b515PeriodFieldObject["fields"].(map[string]any)
	if !b515PeriodFieldOK {
		return nil, b515PeriodFieldFmt.Errorf("fields object required")
	}
	for b515PeriodFieldKey, b515PeriodFieldValue := range b515PeriodFieldFields {
		b515PeriodFieldArray, b515PeriodFieldOK := b515PeriodFieldValue.([]any)
		if !b515PeriodFieldOK {
			return nil, b515PeriodFieldFmt.Errorf("field %s requires byte array", b515PeriodFieldKey)
		}
		for _, b515PeriodFieldByte := range b515PeriodFieldArray {
			if _, b515PeriodFieldError := b515PeriodFieldInteger(b515PeriodFieldByte, 255); b515PeriodFieldError != nil {
				return nil, b515PeriodFieldError
			}
		}
	}
	b515PeriodFieldSupported := true
	for b515PeriodFieldKey, b515PeriodFieldExpected := range map[string]int{"primary": 181, "secondary": 21} {
		b515PeriodFieldValue, b515PeriodFieldPresent := b515PeriodFieldContext[b515PeriodFieldKey]
		if !b515PeriodFieldPresent {
			b515PeriodFieldSupported = false
			continue
		}
		b515PeriodFieldNumber, b515PeriodFieldError := b515PeriodFieldInteger(b515PeriodFieldValue, 255)
		if b515PeriodFieldError != nil {
			return nil, b515PeriodFieldError
		}
		if b515PeriodFieldNumber != b515PeriodFieldExpected {
			b515PeriodFieldSupported = false
		}
	}
	if b515PeriodFieldOperation, b515PeriodFieldPresent := b515PeriodFieldContext["operation"]; b515PeriodFieldPresent {
		if _, b515PeriodFieldError := b515PeriodFieldInteger(b515PeriodFieldOperation, 255); b515PeriodFieldError != nil {
			return nil, b515PeriodFieldError
		}
		b515PeriodFieldSupported = false
	}
	for _, b515PeriodFieldKey := range []string{"direction", "view"} {
		if b515PeriodFieldValue, b515PeriodFieldPresent := b515PeriodFieldContext[b515PeriodFieldKey]; b515PeriodFieldPresent {
			if _, b515PeriodFieldOK := b515PeriodFieldValue.(string); !b515PeriodFieldOK {
				return nil, b515PeriodFieldFmt.Errorf("context %s requires string", b515PeriodFieldKey)
			}
		}
	}
	b515PeriodFieldView, _ := b515PeriodFieldContext["view"].(string)
	b515PeriodFieldDirection, _ := b515PeriodFieldContext["direction"].(string)
	b515PeriodFieldSupported = b515PeriodFieldSupported && ((b515PeriodFieldView == "periods" && (b515PeriodFieldDirection == "request" || b515PeriodFieldDirection == "positive")) || (b515PeriodFieldView == "intervals" && b515PeriodFieldDirection == "request"))
	b515PeriodFieldArguments := map[string]any{}
	if b515PeriodFieldValue, b515PeriodFieldPresent := b515PeriodFieldObject["arguments"]; b515PeriodFieldPresent {
		b515PeriodFieldArguments, b515PeriodFieldOK = b515PeriodFieldValue.(map[string]any)
		if !b515PeriodFieldOK {
			return nil, b515PeriodFieldFmt.Errorf("arguments object required")
		}
	}
	b515PeriodFieldMinutes := 0
	if b515PeriodFieldValue, b515PeriodFieldPresent := b515PeriodFieldArguments["minutes"]; b515PeriodFieldPresent {
		b515PeriodFieldMinutes, b515PeriodFieldError = b515PeriodFieldInteger(b515PeriodFieldValue, 1440)
		if b515PeriodFieldError != nil {
			return nil, b515PeriodFieldError
		}
	}
	b515PeriodFieldPairs := [3][2]int{}
	if b515PeriodFieldValue, b515PeriodFieldPresent := b515PeriodFieldArguments["pairs"]; b515PeriodFieldPresent {
		b515PeriodFieldArray, b515PeriodFieldOK := b515PeriodFieldValue.([]any)
		if !b515PeriodFieldOK || len(b515PeriodFieldArray) != 3 {
			return nil, b515PeriodFieldFmt.Errorf("three pairs required")
		}
		for b515PeriodFieldIndex, b515PeriodFieldValue := range b515PeriodFieldArray {
			b515PeriodFieldPair, b515PeriodFieldOK := b515PeriodFieldValue.([]any)
			if !b515PeriodFieldOK || len(b515PeriodFieldPair) != 2 {
				return nil, b515PeriodFieldFmt.Errorf("two endpoints required")
			}
			for b515PeriodFieldEndpoint, b515PeriodFieldValue := range b515PeriodFieldPair {
				b515PeriodFieldPairs[b515PeriodFieldIndex][b515PeriodFieldEndpoint], b515PeriodFieldError = b515PeriodFieldInteger(b515PeriodFieldValue, 255)
				if b515PeriodFieldError != nil {
					return nil, b515PeriodFieldError
				}
			}
		}
	} else if b515PeriodFieldSupported && b515PeriodFieldView == "intervals" {
		return nil, b515PeriodFieldFmt.Errorf("pairs required")
	}
	b515PeriodFieldResults := map[string]any{}
	b515PeriodFieldExtra := map[string]any{}
	for b515PeriodFieldKey, b515PeriodFieldValue := range b515PeriodFieldObject {
		if b515PeriodFieldKey != "context" && b515PeriodFieldKey != "fields" && b515PeriodFieldKey != "arguments" {
			b515PeriodFieldExtra[b515PeriodFieldKey] = b515PeriodFieldValue
		}
	}
	b515PeriodFieldOutput := map[string]any{"context": b515PeriodFieldContext, "fields": b515PeriodFieldFields, "extra": b515PeriodFieldExtra, "results": b515PeriodFieldResults, "state": "Unknown", "physical": "Unknown", "qualified": false}
	if _, b515PeriodFieldPresent := b515PeriodFieldObject["arguments"]; b515PeriodFieldPresent {
		b515PeriodFieldOutput["arguments"] = b515PeriodFieldArguments
	}
	if !b515PeriodFieldSupported {
		return b515PeriodFieldJSON.Marshal(b515PeriodFieldOutput)
	}
	b515PeriodFieldOutput["state"] = "complete"
	for b515PeriodFieldKey, b515PeriodFieldValue := range b515PeriodFieldFields {
		b515PeriodFieldResults[b515PeriodFieldKey] = map[string]any{"raw": b515PeriodFieldValue, "state": "Unknown"}
	}
	if b515PeriodFieldView == "periods" {
		for _, b515PeriodFieldKey := range []string{"period0", "period1", "period2"} {
			b515PeriodFieldResult := map[string]any{"state": "incomplete"}
			b515PeriodFieldValue, b515PeriodFieldPresent := b515PeriodFieldFields[b515PeriodFieldKey]
			if b515PeriodFieldPresent {
				b515PeriodFieldResult["raw"] = b515PeriodFieldValue
			}
			b515PeriodFieldArray, _ := b515PeriodFieldValue.([]any)
			if len(b515PeriodFieldArray) == 1 {
				b515PeriodFieldResult = b515PeriodFieldPeriod(b515PeriodFieldArray[0])
				b515PeriodFieldResult["raw"] = b515PeriodFieldValue
			} else {
				b515PeriodFieldOutput["state"] = "partial"
			}
			b515PeriodFieldResults[b515PeriodFieldKey] = b515PeriodFieldResult
		}
		if b515PeriodFieldValue, b515PeriodFieldPresent := b515PeriodFieldArguments["minutes"]; b515PeriodFieldPresent {
			b515PeriodFieldResult := map[string]any{"raw": b515PeriodFieldValue, "state": "Unknown"}
			if b515PeriodFieldMinutes%10 == 0 {
				b515PeriodFieldResult["state"] = "Hypothesis"
				b515PeriodFieldResult["calculation"] = "Proven"
				b515PeriodFieldResult["number"] = b515PeriodFieldMinutes / 10
			}
			b515PeriodFieldResults["minutes"] = b515PeriodFieldResult
		}
	} else {
		b515PeriodFieldPairResults := []any{}
		b515PeriodFieldMapped := true
		for b515PeriodFieldIndex, b515PeriodFieldPair := range b515PeriodFieldPairs {
			b515PeriodFieldLabel := "enabled-candidate"
			if b515PeriodFieldPair[0] == b515PeriodFieldPair[1] {
				b515PeriodFieldLabel = "disabled"
			}
			b515PeriodFieldPairResults = append(b515PeriodFieldPairResults, map[string]any{"raw": b515PeriodFieldArguments["pairs"].([]any)[b515PeriodFieldIndex], "state": "Hypothesis", "calculation": "Proven", "label": b515PeriodFieldLabel, "start": b515PeriodFieldPair[0], "end": b515PeriodFieldPair[1]})
			if b515PeriodFieldPair[0] > 144 || b515PeriodFieldPair[1] > 144 {
				b515PeriodFieldMapped = false
			}
		}
		b515PeriodFieldResults["pairs"] = map[string]any{"raw": b515PeriodFieldArguments["pairs"], "state": "Hypothesis", "alternatives": b515PeriodFieldPairResults}
		b515PeriodFieldOrder := map[string]any{"raw": b515PeriodFieldArguments["pairs"], "state": "Unknown"}
		if b515PeriodFieldMapped {
			b515PeriodFieldOrder["state"] = "Hypothesis"
			b515PeriodFieldOrder["calculation"] = "Proven"
			b515PeriodFieldOrder["value"] = b515PeriodFieldPairs[0][0] < b515PeriodFieldPairs[1][0] && b515PeriodFieldPairs[1][0] < b515PeriodFieldPairs[2][0]
		}
		b515PeriodFieldResults["start_order"] = b515PeriodFieldOrder
	}
	return b515PeriodFieldJSON.Marshal(b515PeriodFieldOutput)
}
