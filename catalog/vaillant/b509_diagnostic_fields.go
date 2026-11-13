package vaillant

import (
	b509DiagnosticFieldJSON "encoding/json"
	b509DiagnosticFieldFmt "fmt"
	b509DiagnosticFieldIO "io"
	b509DiagnosticFieldStrconv "strconv"
	b509DiagnosticFieldStrings "strings"
)

func b509DiagnosticFieldRead(b509DiagnosticFieldDecoder *b509DiagnosticFieldJSON.Decoder) (any, error) {
	b509DiagnosticFieldToken, b509DiagnosticFieldErr := b509DiagnosticFieldDecoder.Token()
	if b509DiagnosticFieldErr != nil {
		return nil, b509DiagnosticFieldErr
	}
	switch b509DiagnosticFieldToken {
	case b509DiagnosticFieldJSON.Delim('{'):
		b509DiagnosticFieldObject := map[string]any{}
		for b509DiagnosticFieldDecoder.More() {
			b509DiagnosticFieldKeyToken, b509DiagnosticFieldErr := b509DiagnosticFieldDecoder.Token()
			if b509DiagnosticFieldErr != nil {
				return nil, b509DiagnosticFieldErr
			}
			b509DiagnosticFieldKey, b509DiagnosticFieldOK := b509DiagnosticFieldKeyToken.(string)
			if !b509DiagnosticFieldOK {
				return nil, b509DiagnosticFieldFmt.Errorf("object key required")
			}
			if _, b509DiagnosticFieldExists := b509DiagnosticFieldObject[b509DiagnosticFieldKey]; b509DiagnosticFieldExists {
				return nil, b509DiagnosticFieldFmt.Errorf("duplicate key %q", b509DiagnosticFieldKey)
			}
			b509DiagnosticFieldValue, b509DiagnosticFieldErr := b509DiagnosticFieldRead(b509DiagnosticFieldDecoder)
			if b509DiagnosticFieldErr != nil {
				return nil, b509DiagnosticFieldErr
			}
			b509DiagnosticFieldObject[b509DiagnosticFieldKey] = b509DiagnosticFieldValue
		}
		b509DiagnosticFieldEnd, b509DiagnosticFieldErr := b509DiagnosticFieldDecoder.Token()
		if b509DiagnosticFieldErr != nil || b509DiagnosticFieldEnd != b509DiagnosticFieldJSON.Delim('}') {
			return nil, b509DiagnosticFieldFmt.Errorf("unclosed object")
		}
		return b509DiagnosticFieldObject, nil
	case b509DiagnosticFieldJSON.Delim('['):
		b509DiagnosticFieldArray := []any{}
		for b509DiagnosticFieldDecoder.More() {
			b509DiagnosticFieldValue, b509DiagnosticFieldErr := b509DiagnosticFieldRead(b509DiagnosticFieldDecoder)
			if b509DiagnosticFieldErr != nil {
				return nil, b509DiagnosticFieldErr
			}
			b509DiagnosticFieldArray = append(b509DiagnosticFieldArray, b509DiagnosticFieldValue)
		}
		b509DiagnosticFieldEnd, b509DiagnosticFieldErr := b509DiagnosticFieldDecoder.Token()
		if b509DiagnosticFieldErr != nil || b509DiagnosticFieldEnd != b509DiagnosticFieldJSON.Delim(']') {
			return nil, b509DiagnosticFieldFmt.Errorf("unclosed array")
		}
		return b509DiagnosticFieldArray, nil
	case b509DiagnosticFieldJSON.Delim('}'), b509DiagnosticFieldJSON.Delim(']'):
		return nil, b509DiagnosticFieldFmt.Errorf("unexpected delimiter")
	default:
		return b509DiagnosticFieldToken, nil
	}
}

func b509DiagnosticFieldUnsigned(b509DiagnosticFieldValue any, b509DiagnosticFieldMax uint64) (uint64, bool) {
	b509DiagnosticFieldNumber, b509DiagnosticFieldOK := b509DiagnosticFieldValue.(b509DiagnosticFieldJSON.Number)
	if !b509DiagnosticFieldOK {
		return 0, false
	}
	b509DiagnosticFieldParsed, b509DiagnosticFieldErr := b509DiagnosticFieldStrconv.ParseUint(string(b509DiagnosticFieldNumber), 10, 64)
	return b509DiagnosticFieldParsed, b509DiagnosticFieldErr == nil && b509DiagnosticFieldParsed <= b509DiagnosticFieldMax
}

func B509DiagnosticFieldEvaluate(b509DiagnosticFieldInput b509DiagnosticFieldJSON.RawMessage) (b509DiagnosticFieldJSON.RawMessage, error) {
	b509DiagnosticFieldDecoder := b509DiagnosticFieldJSON.NewDecoder(b509DiagnosticFieldStrings.NewReader(string(b509DiagnosticFieldInput)))
	b509DiagnosticFieldDecoder.UseNumber()
	b509DiagnosticFieldValue, b509DiagnosticFieldErr := b509DiagnosticFieldRead(b509DiagnosticFieldDecoder)
	if b509DiagnosticFieldErr != nil {
		return nil, b509DiagnosticFieldErr
	}
	if _, b509DiagnosticFieldErr = b509DiagnosticFieldDecoder.Token(); b509DiagnosticFieldErr != b509DiagnosticFieldIO.EOF {
		return nil, b509DiagnosticFieldFmt.Errorf("trailing JSON")
	}
	b509DiagnosticFieldInputObject, b509DiagnosticFieldOK := b509DiagnosticFieldValue.(map[string]any)
	if !b509DiagnosticFieldOK {
		return nil, b509DiagnosticFieldFmt.Errorf("input object required")
	}
	b509DiagnosticFieldContext, b509DiagnosticFieldOK := b509DiagnosticFieldInputObject["context"].(map[string]any)
	if !b509DiagnosticFieldOK {
		return nil, b509DiagnosticFieldFmt.Errorf("context object required")
	}
	b509DiagnosticFieldRaw, b509DiagnosticFieldOK := b509DiagnosticFieldInputObject["fields"].(map[string]any)
	if !b509DiagnosticFieldOK {
		return nil, b509DiagnosticFieldFmt.Errorf("fields object required")
	}
	for _, b509DiagnosticFieldKey := range []string{"primary", "secondary", "operation"} {
		if b509DiagnosticFieldValue, b509DiagnosticFieldExists := b509DiagnosticFieldContext[b509DiagnosticFieldKey]; b509DiagnosticFieldExists && b509DiagnosticFieldValue != nil {
			if _, b509DiagnosticFieldOK := b509DiagnosticFieldUnsigned(b509DiagnosticFieldValue, 255); !b509DiagnosticFieldOK {
				return nil, b509DiagnosticFieldFmt.Errorf("%s must be u8 or null", b509DiagnosticFieldKey)
			}
		}
	}
	for _, b509DiagnosticFieldKey := range []string{"direction", "view"} {
		if b509DiagnosticFieldValue, b509DiagnosticFieldExists := b509DiagnosticFieldContext[b509DiagnosticFieldKey]; b509DiagnosticFieldExists && b509DiagnosticFieldValue != nil {
			if _, b509DiagnosticFieldOK := b509DiagnosticFieldValue.(string); !b509DiagnosticFieldOK {
				return nil, b509DiagnosticFieldFmt.Errorf("%s must be string or null", b509DiagnosticFieldKey)
			}
		}
	}
	for b509DiagnosticFieldKey, b509DiagnosticFieldValue := range b509DiagnosticFieldRaw {
		b509DiagnosticFieldArray, b509DiagnosticFieldOK := b509DiagnosticFieldValue.([]any)
		if !b509DiagnosticFieldOK {
			return nil, b509DiagnosticFieldFmt.Errorf("%s must be a byte array", b509DiagnosticFieldKey)
		}
		for _, b509DiagnosticFieldByte := range b509DiagnosticFieldArray {
			if _, b509DiagnosticFieldOK := b509DiagnosticFieldUnsigned(b509DiagnosticFieldByte, 255); !b509DiagnosticFieldOK {
				return nil, b509DiagnosticFieldFmt.Errorf("%s contains a non-byte", b509DiagnosticFieldKey)
			}
		}
	}
	b509DiagnosticFieldResults := map[string]any{}
	b509DiagnosticFieldOutput := map[string]any{"context": b509DiagnosticFieldContext, "raw": b509DiagnosticFieldRaw, "state": "Unknown", "fields": b509DiagnosticFieldResults, "physical": "Unknown", "qualified": false}
	if b509DiagnosticFieldArgumentsValue, b509DiagnosticFieldExists := b509DiagnosticFieldInputObject["arguments"]; b509DiagnosticFieldExists {
		b509DiagnosticFieldArguments, b509DiagnosticFieldOK := b509DiagnosticFieldArgumentsValue.(map[string]any)
		if !b509DiagnosticFieldOK {
			return nil, b509DiagnosticFieldFmt.Errorf("arguments object required")
		}
		if b509DiagnosticFieldRequested, b509DiagnosticFieldExists := b509DiagnosticFieldArguments["requested_address"]; b509DiagnosticFieldExists {
			if _, b509DiagnosticFieldOK := b509DiagnosticFieldUnsigned(b509DiagnosticFieldRequested, 4294967295); !b509DiagnosticFieldOK {
				return nil, b509DiagnosticFieldFmt.Errorf("requested_address must be u32")
			}
		}
		b509DiagnosticFieldOutput["arguments"] = b509DiagnosticFieldArguments
	}
	b509DiagnosticFieldPrimary, b509DiagnosticFieldPrimaryOK := b509DiagnosticFieldUnsigned(b509DiagnosticFieldContext["primary"], 255)
	b509DiagnosticFieldSecondary, b509DiagnosticFieldSecondaryOK := b509DiagnosticFieldUnsigned(b509DiagnosticFieldContext["secondary"], 255)
	b509DiagnosticFieldOperation, b509DiagnosticFieldOperationOK := b509DiagnosticFieldUnsigned(b509DiagnosticFieldContext["operation"], 255)
	b509DiagnosticFieldDirection, _ := b509DiagnosticFieldContext["direction"].(string)
	b509DiagnosticFieldView, _ := b509DiagnosticFieldContext["view"].(string)
	if !b509DiagnosticFieldPrimaryOK || !b509DiagnosticFieldSecondaryOK || !b509DiagnosticFieldOperationOK || b509DiagnosticFieldPrimary != 181 || b509DiagnosticFieldSecondary != 9 ||
		(b509DiagnosticFieldDirection != "request" && b509DiagnosticFieldDirection != "response") ||
		((b509DiagnosticFieldOperation != 52 || b509DiagnosticFieldView != "diagnostic-get") && (b509DiagnosticFieldOperation != 53 || b509DiagnosticFieldView != "diagnostic-put")) {
		return b509DiagnosticFieldJSON.Marshal(b509DiagnosticFieldOutput)
	}
	b509DiagnosticFieldOutput["state"] = "complete"
	b509DiagnosticFieldOutput["grade"] = "Hypothesis"
	b509DiagnosticFieldOutput["correlation"] = "Unknown"
	b509DiagnosticFieldAliases := map[string]string{"count": "requested_count", "status": "write_status"}
	for b509DiagnosticFieldShort, b509DiagnosticFieldLong := range b509DiagnosticFieldAliases {
		_, b509DiagnosticFieldShortExists := b509DiagnosticFieldRaw[b509DiagnosticFieldShort]
		_, b509DiagnosticFieldLongExists := b509DiagnosticFieldRaw[b509DiagnosticFieldLong]
		if b509DiagnosticFieldShortExists && b509DiagnosticFieldLongExists {
			return nil, b509DiagnosticFieldFmt.Errorf("%s and %s are aliases and must not both be present", b509DiagnosticFieldShort, b509DiagnosticFieldLong)
		}
	}
	b509DiagnosticFieldExpected := map[string]int{}
	if b509DiagnosticFieldDirection == "request" || b509DiagnosticFieldOperation == 52 {
		b509DiagnosticFieldExpected["address"] = 4
	}
	if b509DiagnosticFieldOperation == 52 && b509DiagnosticFieldDirection == "request" {
		b509DiagnosticFieldExpected[b509DiagnosticFieldAliasKey(b509DiagnosticFieldRaw, "count", b509DiagnosticFieldAliases)] = 1
	} else if b509DiagnosticFieldOperation == 53 && b509DiagnosticFieldDirection == "response" {
		b509DiagnosticFieldExpected[b509DiagnosticFieldAliasKey(b509DiagnosticFieldRaw, "status", b509DiagnosticFieldAliases)] = 1
	} else {
		b509DiagnosticFieldExpected["data"] = -1
	}
	for b509DiagnosticFieldKey, b509DiagnosticFieldRawValue := range b509DiagnosticFieldRaw {
		b509DiagnosticFieldResults[b509DiagnosticFieldKey] = map[string]any{"raw": b509DiagnosticFieldRawValue, "state": "Unknown"}
	}
	for b509DiagnosticFieldKey, b509DiagnosticFieldWidth := range b509DiagnosticFieldExpected {
		b509DiagnosticFieldRawValue, b509DiagnosticFieldExists := b509DiagnosticFieldRaw[b509DiagnosticFieldKey]
		b509DiagnosticFieldResult := map[string]any{"state": "incomplete", "grade": "Hypothesis"}
		b509DiagnosticFieldResults[b509DiagnosticFieldKey] = b509DiagnosticFieldResult
		if b509DiagnosticFieldExists {
			b509DiagnosticFieldResult["raw"] = b509DiagnosticFieldRawValue
		}
		b509DiagnosticFieldArray, _ := b509DiagnosticFieldRawValue.([]any)
		if !b509DiagnosticFieldExists || (b509DiagnosticFieldWidth >= 0 && len(b509DiagnosticFieldArray) != b509DiagnosticFieldWidth) {
			b509DiagnosticFieldOutput["state"] = "partial"
			continue
		}
		b509DiagnosticFieldResult["state"] = "complete"
		if b509DiagnosticFieldKey == "data" {
			b509DiagnosticFieldResult["length"] = len(b509DiagnosticFieldArray)
			b509DiagnosticFieldResult["unmodeled"] = true
			continue
		}
		var b509DiagnosticFieldNumber uint64
		for b509DiagnosticFieldIndex, b509DiagnosticFieldByte := range b509DiagnosticFieldArray {
			b509DiagnosticFieldPart, _ := b509DiagnosticFieldUnsigned(b509DiagnosticFieldByte, 255)
			b509DiagnosticFieldNumber |= b509DiagnosticFieldPart << (8 * b509DiagnosticFieldIndex)
		}
		b509DiagnosticFieldResult["number"] = b509DiagnosticFieldNumber
		b509DiagnosticFieldResult["calculation"] = "Proven"
		if b509DiagnosticFieldKey == "status" || b509DiagnosticFieldKey == "write_status" {
			b509DiagnosticFieldResult["candidate_domain"] = b509DiagnosticFieldNumber <= 6
		}
	}
	return b509DiagnosticFieldJSON.Marshal(b509DiagnosticFieldOutput)
}

func b509DiagnosticFieldAliasKey(b509DiagnosticFieldRaw map[string]any, b509DiagnosticFieldShort string, b509DiagnosticFieldAliases map[string]string) string {
	if _, b509DiagnosticFieldExists := b509DiagnosticFieldRaw[b509DiagnosticFieldAliases[b509DiagnosticFieldShort]]; b509DiagnosticFieldExists {
		return b509DiagnosticFieldAliases[b509DiagnosticFieldShort]
	}
	return b509DiagnosticFieldShort
}
