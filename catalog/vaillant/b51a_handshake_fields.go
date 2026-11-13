package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func b51AHandshakeFieldRead(b51AHandshakeFieldDecoder *json.Decoder) (any, error) {
	b51AHandshakeFieldToken, b51AHandshakeFieldErr := b51AHandshakeFieldDecoder.Token()
	if b51AHandshakeFieldErr != nil {
		return nil, b51AHandshakeFieldErr
	}
	if b51AHandshakeFieldDelimiter, b51AHandshakeFieldOK := b51AHandshakeFieldToken.(json.Delim); b51AHandshakeFieldOK {
		switch b51AHandshakeFieldDelimiter {
		case '{':
			b51AHandshakeFieldObject := map[string]any{}
			for b51AHandshakeFieldDecoder.More() {
				b51AHandshakeFieldKeyToken, b51AHandshakeFieldErr := b51AHandshakeFieldDecoder.Token()
				if b51AHandshakeFieldErr != nil {
					return nil, b51AHandshakeFieldErr
				}
				b51AHandshakeFieldKey := b51AHandshakeFieldKeyToken.(string)
				if _, b51AHandshakeFieldExists := b51AHandshakeFieldObject[b51AHandshakeFieldKey]; b51AHandshakeFieldExists {
					return nil, fmt.Errorf("duplicate key %q", b51AHandshakeFieldKey)
				}
				b51AHandshakeFieldValue, b51AHandshakeFieldErr := b51AHandshakeFieldRead(b51AHandshakeFieldDecoder)
				if b51AHandshakeFieldErr != nil {
					return nil, b51AHandshakeFieldErr
				}
				b51AHandshakeFieldObject[b51AHandshakeFieldKey] = b51AHandshakeFieldValue
			}
			_, b51AHandshakeFieldErr := b51AHandshakeFieldDecoder.Token()
			return b51AHandshakeFieldObject, b51AHandshakeFieldErr
		case '[':
			b51AHandshakeFieldArray := []any{}
			for b51AHandshakeFieldDecoder.More() {
				b51AHandshakeFieldValue, b51AHandshakeFieldErr := b51AHandshakeFieldRead(b51AHandshakeFieldDecoder)
				if b51AHandshakeFieldErr != nil {
					return nil, b51AHandshakeFieldErr
				}
				b51AHandshakeFieldArray = append(b51AHandshakeFieldArray, b51AHandshakeFieldValue)
			}
			_, b51AHandshakeFieldErr := b51AHandshakeFieldDecoder.Token()
			return b51AHandshakeFieldArray, b51AHandshakeFieldErr
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	if b51AHandshakeFieldNumber, b51AHandshakeFieldOK := b51AHandshakeFieldToken.(json.Number); b51AHandshakeFieldOK {
		if _, b51AHandshakeFieldErr := b51AHandshakeFieldNumber.Float64(); b51AHandshakeFieldErr != nil {
			return nil, fmt.Errorf("nonfinite number")
		}
	}
	return b51AHandshakeFieldToken, nil
}

func b51AHandshakeFieldByte(b51AHandshakeFieldValue any) (int, bool) {
	b51AHandshakeFieldNumber, b51AHandshakeFieldOK := b51AHandshakeFieldValue.(json.Number)
	if !b51AHandshakeFieldOK {
		return 0, false
	}
	b51AHandshakeFieldInt, b51AHandshakeFieldErr := b51AHandshakeFieldNumber.Int64()
	return int(b51AHandshakeFieldInt), b51AHandshakeFieldErr == nil && b51AHandshakeFieldInt >= 0 && b51AHandshakeFieldInt <= 255
}

func B51AHandshakeFieldEvaluate(b51AHandshakeFieldInput json.RawMessage) (json.RawMessage, error) {
	b51AHandshakeFieldDecoder := json.NewDecoder(bytes.NewReader(b51AHandshakeFieldInput))
	b51AHandshakeFieldDecoder.UseNumber()
	b51AHandshakeFieldValue, b51AHandshakeFieldErr := b51AHandshakeFieldRead(b51AHandshakeFieldDecoder)
	if b51AHandshakeFieldErr != nil {
		return nil, b51AHandshakeFieldErr
	}
	if _, b51AHandshakeFieldErr = b51AHandshakeFieldDecoder.Token(); b51AHandshakeFieldErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b51AHandshakeFieldRoot, b51AHandshakeFieldOK := b51AHandshakeFieldValue.(map[string]any)
	if !b51AHandshakeFieldOK {
		return nil, fmt.Errorf("input must be object")
	}
	b51AHandshakeFieldContext, b51AHandshakeFieldOK := b51AHandshakeFieldRoot["context"].(map[string]any)
	if !b51AHandshakeFieldOK {
		return nil, fmt.Errorf("context must be object")
	}
	b51AHandshakeFieldFields, b51AHandshakeFieldOK := b51AHandshakeFieldRoot["fields"].(map[string]any)
	if !b51AHandshakeFieldOK {
		return nil, fmt.Errorf("fields must be object")
	}
	for _, b51AHandshakeFieldKey := range []string{"primary", "secondary", "operation"} {
		if b51AHandshakeFieldValue, b51AHandshakeFieldPresent := b51AHandshakeFieldContext[b51AHandshakeFieldKey]; b51AHandshakeFieldPresent {
			if _, b51AHandshakeFieldOK := b51AHandshakeFieldByte(b51AHandshakeFieldValue); !b51AHandshakeFieldOK {
				return nil, fmt.Errorf("invalid %s", b51AHandshakeFieldKey)
			}
		}
	}
	for _, b51AHandshakeFieldKey := range []string{"direction", "view"} {
		if b51AHandshakeFieldValue, b51AHandshakeFieldPresent := b51AHandshakeFieldContext[b51AHandshakeFieldKey]; b51AHandshakeFieldPresent {
			if _, b51AHandshakeFieldOK := b51AHandshakeFieldValue.(string); !b51AHandshakeFieldOK {
				return nil, fmt.Errorf("invalid %s", b51AHandshakeFieldKey)
			}
		}
	}
	if b51AHandshakeFieldValue, b51AHandshakeFieldPresent := b51AHandshakeFieldRoot["arguments"]; b51AHandshakeFieldPresent {
		b51AHandshakeFieldArguments, b51AHandshakeFieldOK := b51AHandshakeFieldValue.(map[string]any)
		if !b51AHandshakeFieldOK || len(b51AHandshakeFieldArguments) != 0 {
			return nil, fmt.Errorf("no arguments supported")
		}
	}
	b51AHandshakeFieldBytes := map[string][]int{}
	for b51AHandshakeFieldKey, b51AHandshakeFieldValue := range b51AHandshakeFieldFields {
		b51AHandshakeFieldArray, b51AHandshakeFieldOK := b51AHandshakeFieldValue.([]any)
		if !b51AHandshakeFieldOK {
			return nil, fmt.Errorf("field %s must be byte array", b51AHandshakeFieldKey)
		}
		b51AHandshakeFieldBytes[b51AHandshakeFieldKey] = []int{}
		for _, b51AHandshakeFieldValue := range b51AHandshakeFieldArray {
			b51AHandshakeFieldInt, b51AHandshakeFieldOK := b51AHandshakeFieldByte(b51AHandshakeFieldValue)
			if !b51AHandshakeFieldOK {
				return nil, fmt.Errorf("invalid byte in %s", b51AHandshakeFieldKey)
			}
			b51AHandshakeFieldBytes[b51AHandshakeFieldKey] = append(b51AHandshakeFieldBytes[b51AHandshakeFieldKey], b51AHandshakeFieldInt)
		}
	}
	b51AHandshakeFieldPrimary, _ := b51AHandshakeFieldByte(b51AHandshakeFieldContext["primary"])
	b51AHandshakeFieldSecondary, _ := b51AHandshakeFieldByte(b51AHandshakeFieldContext["secondary"])
	b51AHandshakeFieldOperation, b51AHandshakeFieldHasOperation := b51AHandshakeFieldContext["operation"]
	b51AHandshakeFieldOp, _ := b51AHandshakeFieldByte(b51AHandshakeFieldOperation)
	b51AHandshakeFieldView, _ := b51AHandshakeFieldContext["view"].(string)
	b51AHandshakeFieldDirection, _ := b51AHandshakeFieldContext["direction"].(string)
	b51AHandshakeFieldReset := !b51AHandshakeFieldHasOperation
	b51AHandshakeFieldHandshake := b51AHandshakeFieldHasOperation && b51AHandshakeFieldOp == 1
	b51AHandshakeFieldSupported := b51AHandshakeFieldPrimary == 181 && b51AHandshakeFieldSecondary == 26 && ((b51AHandshakeFieldView == "reset" && b51AHandshakeFieldReset || b51AHandshakeFieldView == "handshake" && b51AHandshakeFieldHandshake) && (b51AHandshakeFieldDirection == "request" || b51AHandshakeFieldDirection == "positive") || b51AHandshakeFieldView == "sequence-pair" && b51AHandshakeFieldDirection == "pair" && (b51AHandshakeFieldReset || b51AHandshakeFieldHandshake))
	b51AHandshakeFieldResults := map[string]any{}
	for b51AHandshakeFieldKey, b51AHandshakeFieldRaw := range b51AHandshakeFieldFields {
		b51AHandshakeFieldResults[b51AHandshakeFieldKey] = map[string]any{"raw": b51AHandshakeFieldRaw, "state": "Unknown"}
	}
	b51AHandshakeFieldRoot["results"] = b51AHandshakeFieldResults
	b51AHandshakeFieldRoot["physical"] = "Unknown"
	b51AHandshakeFieldRoot["qualified"] = false
	b51AHandshakeFieldRoot["state"] = "Unknown"
	if !b51AHandshakeFieldSupported {
		return json.Marshal(b51AHandshakeFieldRoot)
	}
	b51AHandshakeFieldRoot["state"] = "complete"
	b51AHandshakeFieldFixed := func(b51AHandshakeFieldKey string, b51AHandshakeFieldWidth int) (map[string]any, []int, bool) {
		b51AHandshakeFieldRaw, b51AHandshakeFieldPresent := b51AHandshakeFieldFields[b51AHandshakeFieldKey]
		b51AHandshakeFieldResult := map[string]any{"state": "incomplete", "grade": "Hypothesis"}
		if b51AHandshakeFieldPresent {
			b51AHandshakeFieldResult["raw"] = b51AHandshakeFieldRaw
		}
		b51AHandshakeFieldResults[b51AHandshakeFieldKey] = b51AHandshakeFieldResult
		b51AHandshakeFieldData := b51AHandshakeFieldBytes[b51AHandshakeFieldKey]
		if !b51AHandshakeFieldPresent || len(b51AHandshakeFieldData) != b51AHandshakeFieldWidth {
			b51AHandshakeFieldRoot["state"] = "partial"
			return b51AHandshakeFieldResult, nil, false
		}
		b51AHandshakeFieldResult["state"] = "complete"
		return b51AHandshakeFieldResult, b51AHandshakeFieldData, true
	}
	if b51AHandshakeFieldView == "sequence-pair" {
		b51AHandshakeFieldRequest, b51AHandshakeFieldRequestData, b51AHandshakeFieldRequestOK := b51AHandshakeFieldFixed("request_sequence", 1)
		_, b51AHandshakeFieldResponseData, b51AHandshakeFieldResponseOK := b51AHandshakeFieldFixed("response_sequence", 1)
		if b51AHandshakeFieldRequestOK && b51AHandshakeFieldResponseOK {
			b51AHandshakeFieldRequest["flags"] = map[string]bool{"correlated": b51AHandshakeFieldRequestData[0] == b51AHandshakeFieldResponseData[0]}
			b51AHandshakeFieldRequest["calculation"] = "Proven"
		}
		return json.Marshal(b51AHandshakeFieldRoot)
	}
	if b51AHandshakeFieldResult, b51AHandshakeFieldData, b51AHandshakeFieldOK := b51AHandshakeFieldFixed("sequence", 1); b51AHandshakeFieldOK {
		b51AHandshakeFieldResult["number"] = b51AHandshakeFieldData[0]
		b51AHandshakeFieldResult["calculation"] = "Proven"
	}
	if b51AHandshakeFieldReset {
		if b51AHandshakeFieldDirection == "positive" {
			if b51AHandshakeFieldResult, b51AHandshakeFieldData, b51AHandshakeFieldOK := b51AHandshakeFieldFixed("ack", 1); b51AHandshakeFieldOK {
				b51AHandshakeFieldResult["flags"] = map[string]bool{"denied": b51AHandshakeFieldData[0]&1 != 0}
				b51AHandshakeFieldResult["unmodeled"] = b51AHandshakeFieldData[0] & 254
			}
		}
		return json.Marshal(b51AHandshakeFieldRoot)
	}
	b51AHandshakeFieldFlags, b51AHandshakeFieldData, b51AHandshakeFieldFlagsOK := b51AHandshakeFieldFixed("flags", 1)
	b51AHandshakeFieldChange := false
	if b51AHandshakeFieldFlagsOK {
		b51AHandshakeFieldLabels := []string{"power-on/start", "change", "CRC-support"}
		b51AHandshakeFieldMask := 248
		if b51AHandshakeFieldDirection == "positive" {
			b51AHandshakeFieldLabels = []string{"valid", "invalid", "blocked", "change", "denied", "invalid-range", "not-user-adjustable"}
			b51AHandshakeFieldMask = 128
			b51AHandshakeFieldChange = b51AHandshakeFieldData[0]&8 != 0
		}
		b51AHandshakeFieldDecoded := map[string]bool{}
		for b51AHandshakeFieldIndex, b51AHandshakeFieldLabel := range b51AHandshakeFieldLabels {
			b51AHandshakeFieldDecoded[b51AHandshakeFieldLabel] = b51AHandshakeFieldData[0]&(1<<b51AHandshakeFieldIndex) != 0
		}
		b51AHandshakeFieldFlags["flags"] = b51AHandshakeFieldDecoded
		b51AHandshakeFieldFlags["reserved"] = b51AHandshakeFieldData[0] & b51AHandshakeFieldMask
	}
	for _, b51AHandshakeFieldKey := range []string{"dsn", "start", "offset"} {
		b51AHandshakeFieldResult := map[string]any{"state": "Unknown", "grade": "Hypothesis"}
		if b51AHandshakeFieldRaw, b51AHandshakeFieldPresent := b51AHandshakeFieldFields[b51AHandshakeFieldKey]; b51AHandshakeFieldPresent {
			b51AHandshakeFieldResult["raw"] = b51AHandshakeFieldRaw
			b51AHandshakeFieldResult["state"] = "complete"
		}
		if b51AHandshakeFieldDirection == "positive" && b51AHandshakeFieldFlagsOK && !b51AHandshakeFieldChange {
			b51AHandshakeFieldResult["state"] = "absent-unless-change"
		}
		b51AHandshakeFieldResults[b51AHandshakeFieldKey] = b51AHandshakeFieldResult
	}
	if b51AHandshakeFieldDirection == "request" {
		if b51AHandshakeFieldResult, b51AHandshakeFieldData, b51AHandshakeFieldOK := b51AHandshakeFieldFixed("crc", 2); b51AHandshakeFieldOK {
			b51AHandshakeFieldResult["number"] = b51AHandshakeFieldData[0] + 256*b51AHandshakeFieldData[1]
			b51AHandshakeFieldResult["calculation"] = "Proven"
		}
	} else {
		if b51AHandshakeFieldResult, b51AHandshakeFieldData, b51AHandshakeFieldOK := b51AHandshakeFieldFixed("family_code", 1); b51AHandshakeFieldOK {
			b51AHandshakeFieldResult["number"] = b51AHandshakeFieldData[0]
			b51AHandshakeFieldResult["calculation"] = "Proven"
		}
	}
	return json.Marshal(b51AHandshakeFieldRoot)
}
