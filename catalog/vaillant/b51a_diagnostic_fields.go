package vaillant

import (
	b51ADiagnosticFieldBytes "bytes"
	b51ADiagnosticFieldJSON "encoding/json"
	b51ADiagnosticFieldFmt "fmt"
	b51ADiagnosticFieldIO "io"
	b51ADiagnosticFieldMath "math"
	b51ADiagnosticFieldBig "math/big"
	b51ADiagnosticFieldStrconv "strconv"
)

func b51ADiagnosticFieldRead(b51ADiagnosticFieldDecoder *b51ADiagnosticFieldJSON.Decoder) (any, error) {
	b51ADiagnosticFieldToken, b51ADiagnosticFieldErr := b51ADiagnosticFieldDecoder.Token()
	if b51ADiagnosticFieldErr != nil {
		return nil, b51ADiagnosticFieldErr
	}
	switch b51ADiagnosticFieldToken {
	case b51ADiagnosticFieldJSON.Delim('{'):
		b51ADiagnosticFieldObject := map[string]any{}
		for b51ADiagnosticFieldDecoder.More() {
			b51ADiagnosticFieldKeyToken, b51ADiagnosticFieldErr := b51ADiagnosticFieldDecoder.Token()
			if b51ADiagnosticFieldErr != nil {
				return nil, b51ADiagnosticFieldErr
			}
			b51ADiagnosticFieldKey := b51ADiagnosticFieldKeyToken.(string)
			if _, b51ADiagnosticFieldExists := b51ADiagnosticFieldObject[b51ADiagnosticFieldKey]; b51ADiagnosticFieldExists {
				return nil, b51ADiagnosticFieldFmt.Errorf("duplicate key %q", b51ADiagnosticFieldKey)
			}
			b51ADiagnosticFieldValue, b51ADiagnosticFieldErr := b51ADiagnosticFieldRead(b51ADiagnosticFieldDecoder)
			if b51ADiagnosticFieldErr != nil {
				return nil, b51ADiagnosticFieldErr
			}
			b51ADiagnosticFieldObject[b51ADiagnosticFieldKey] = b51ADiagnosticFieldValue
		}
		_, b51ADiagnosticFieldErr = b51ADiagnosticFieldDecoder.Token()
		return b51ADiagnosticFieldObject, b51ADiagnosticFieldErr
	case b51ADiagnosticFieldJSON.Delim('['):
		b51ADiagnosticFieldArray := []any{}
		for b51ADiagnosticFieldDecoder.More() {
			b51ADiagnosticFieldValue, b51ADiagnosticFieldErr := b51ADiagnosticFieldRead(b51ADiagnosticFieldDecoder)
			if b51ADiagnosticFieldErr != nil {
				return nil, b51ADiagnosticFieldErr
			}
			b51ADiagnosticFieldArray = append(b51ADiagnosticFieldArray, b51ADiagnosticFieldValue)
		}
		_, b51ADiagnosticFieldErr = b51ADiagnosticFieldDecoder.Token()
		return b51ADiagnosticFieldArray, b51ADiagnosticFieldErr
	}
	if b51ADiagnosticFieldNumber, b51ADiagnosticFieldOK := b51ADiagnosticFieldToken.(b51ADiagnosticFieldJSON.Number); b51ADiagnosticFieldOK {
		b51ADiagnosticFieldFloat, b51ADiagnosticFieldErr := b51ADiagnosticFieldStrconv.ParseFloat(string(b51ADiagnosticFieldNumber), 64)
		if b51ADiagnosticFieldErr != nil || b51ADiagnosticFieldMath.IsInf(b51ADiagnosticFieldFloat, 0) || b51ADiagnosticFieldMath.IsNaN(b51ADiagnosticFieldFloat) {
			return nil, b51ADiagnosticFieldFmt.Errorf("nonfinite number")
		}
	}
	return b51ADiagnosticFieldToken, nil
}

func b51ADiagnosticFieldUint(b51ADiagnosticFieldValue any) (uint64, bool) {
	b51ADiagnosticFieldNumber, b51ADiagnosticFieldOK := b51ADiagnosticFieldValue.(b51ADiagnosticFieldJSON.Number)
	if !b51ADiagnosticFieldOK {
		return 0, false
	}
	b51ADiagnosticFieldRat, b51ADiagnosticFieldOK := new(b51ADiagnosticFieldBig.Rat).SetString(string(b51ADiagnosticFieldNumber))
	if !b51ADiagnosticFieldOK || !b51ADiagnosticFieldRat.IsInt() || b51ADiagnosticFieldRat.Sign() < 0 || !b51ADiagnosticFieldRat.Num().IsUint64() || b51ADiagnosticFieldRat.Num().Uint64() > 255 {
		return 0, false
	}
	return b51ADiagnosticFieldRat.Num().Uint64(), true
}

func B51ADiagnosticFieldEvaluate(b51ADiagnosticFieldInput b51ADiagnosticFieldJSON.RawMessage) (b51ADiagnosticFieldJSON.RawMessage, error) {
	b51ADiagnosticFieldDecoder := b51ADiagnosticFieldJSON.NewDecoder(b51ADiagnosticFieldBytes.NewReader(b51ADiagnosticFieldInput))
	b51ADiagnosticFieldDecoder.UseNumber()
	b51ADiagnosticFieldValue, b51ADiagnosticFieldErr := b51ADiagnosticFieldRead(b51ADiagnosticFieldDecoder)
	if b51ADiagnosticFieldErr != nil {
		return nil, b51ADiagnosticFieldErr
	}
	if _, b51ADiagnosticFieldErr = b51ADiagnosticFieldDecoder.Token(); b51ADiagnosticFieldErr != b51ADiagnosticFieldIO.EOF {
		return nil, b51ADiagnosticFieldFmt.Errorf("trailing JSON")
	}
	b51ADiagnosticFieldObject, b51ADiagnosticFieldOK := b51ADiagnosticFieldValue.(map[string]any)
	if !b51ADiagnosticFieldOK {
		return nil, b51ADiagnosticFieldFmt.Errorf("input must be object")
	}
	b51ADiagnosticFieldContext, b51ADiagnosticFieldOK := b51ADiagnosticFieldObject["context"].(map[string]any)
	if !b51ADiagnosticFieldOK {
		return nil, b51ADiagnosticFieldFmt.Errorf("context must be object")
	}
	b51ADiagnosticFieldFields, b51ADiagnosticFieldOK := b51ADiagnosticFieldObject["fields"].(map[string]any)
	if !b51ADiagnosticFieldOK {
		return nil, b51ADiagnosticFieldFmt.Errorf("fields must be object")
	}
	for _, b51ADiagnosticFieldKey := range []string{"primary", "secondary", "operation"} {
		if b51ADiagnosticFieldValue, b51ADiagnosticFieldPresent := b51ADiagnosticFieldContext[b51ADiagnosticFieldKey]; b51ADiagnosticFieldPresent {
			if _, b51ADiagnosticFieldOK := b51ADiagnosticFieldUint(b51ADiagnosticFieldValue); !b51ADiagnosticFieldOK {
				return nil, b51ADiagnosticFieldFmt.Errorf("invalid context %s", b51ADiagnosticFieldKey)
			}
		}
	}
	for _, b51ADiagnosticFieldKey := range []string{"direction", "view"} {
		if b51ADiagnosticFieldValue, b51ADiagnosticFieldPresent := b51ADiagnosticFieldContext[b51ADiagnosticFieldKey]; b51ADiagnosticFieldPresent {
			if _, b51ADiagnosticFieldOK := b51ADiagnosticFieldValue.(string); !b51ADiagnosticFieldOK {
				return nil, b51ADiagnosticFieldFmt.Errorf("invalid context %s", b51ADiagnosticFieldKey)
			}
		}
	}
	for b51ADiagnosticFieldKey, b51ADiagnosticFieldValue := range b51ADiagnosticFieldFields {
		b51ADiagnosticFieldArray, b51ADiagnosticFieldOK := b51ADiagnosticFieldValue.([]any)
		if !b51ADiagnosticFieldOK {
			return nil, b51ADiagnosticFieldFmt.Errorf("field %s must be byte array", b51ADiagnosticFieldKey)
		}
		for _, b51ADiagnosticFieldByte := range b51ADiagnosticFieldArray {
			if _, b51ADiagnosticFieldOK := b51ADiagnosticFieldUint(b51ADiagnosticFieldByte); !b51ADiagnosticFieldOK {
				return nil, b51ADiagnosticFieldFmt.Errorf("invalid byte in %s", b51ADiagnosticFieldKey)
			}
		}
	}
	b51ADiagnosticFieldArguments := map[string]any{}
	if b51ADiagnosticFieldValue, b51ADiagnosticFieldPresent := b51ADiagnosticFieldObject["arguments"]; b51ADiagnosticFieldPresent {
		b51ADiagnosticFieldArguments, b51ADiagnosticFieldOK = b51ADiagnosticFieldValue.(map[string]any)
		if !b51ADiagnosticFieldOK {
			return nil, b51ADiagnosticFieldFmt.Errorf("arguments must be object")
		}
	}
	if b51ADiagnosticFieldValue, b51ADiagnosticFieldPresent := b51ADiagnosticFieldArguments["value_width"]; b51ADiagnosticFieldPresent {
		b51ADiagnosticFieldWidth, b51ADiagnosticFieldOK := b51ADiagnosticFieldUint(b51ADiagnosticFieldValue)
		if !b51ADiagnosticFieldOK || (b51ADiagnosticFieldWidth != 1 && b51ADiagnosticFieldWidth != 2) {
			return nil, b51ADiagnosticFieldFmt.Errorf("value_width must be 1 or 2")
		}
	}
	b51ADiagnosticFieldResults := map[string]any{}
	b51ADiagnosticFieldOutput := map[string]any{}
	for b51ADiagnosticFieldKey, b51ADiagnosticFieldValue := range b51ADiagnosticFieldObject {
		b51ADiagnosticFieldOutput[b51ADiagnosticFieldKey] = b51ADiagnosticFieldValue
	}
	b51ADiagnosticFieldOutput["input"] = b51ADiagnosticFieldObject
	b51ADiagnosticFieldOutput["raw"] = b51ADiagnosticFieldFields
	b51ADiagnosticFieldOutput["results"] = b51ADiagnosticFieldResults
	b51ADiagnosticFieldOutput["state"] = "Unknown"
	b51ADiagnosticFieldOutput["physical"] = "Unknown"
	b51ADiagnosticFieldOutput["qualified"] = false
	b51ADiagnosticFieldPrimary, _ := b51ADiagnosticFieldUint(b51ADiagnosticFieldContext["primary"])
	b51ADiagnosticFieldSecondary, _ := b51ADiagnosticFieldUint(b51ADiagnosticFieldContext["secondary"])
	b51ADiagnosticFieldOperation, _ := b51ADiagnosticFieldUint(b51ADiagnosticFieldContext["operation"])
	b51ADiagnosticFieldView, _ := b51ADiagnosticFieldContext["view"].(string)
	b51ADiagnosticFieldSupported := (b51ADiagnosticFieldOperation == 4 && b51ADiagnosticFieldView == "enable-bits") || (b51ADiagnosticFieldOperation == 5 && (b51ADiagnosticFieldView == "range" || b51ADiagnosticFieldView == "enum" || b51ADiagnosticFieldView == "four-byte")) || (b51ADiagnosticFieldOperation == 6 && b51ADiagnosticFieldView == "save")
	if b51ADiagnosticFieldPrimary != 181 || b51ADiagnosticFieldSecondary != 26 || b51ADiagnosticFieldContext["direction"] != "positive" || !b51ADiagnosticFieldSupported {
		return b51ADiagnosticFieldJSON.Marshal(b51ADiagnosticFieldOutput)
	}
	b51ADiagnosticFieldWidths := map[string]int{}
	switch b51ADiagnosticFieldView {
	case "enable-bits":
		b51ADiagnosticFieldWidths["enable"] = 13
	case "range":
		b51ADiagnosticFieldWidths = map[string]int{"flags": 1, "display_type": 1, "value": 2, "minimum": 2, "maximum": 2, "step": 1}
	case "enum":
		b51ADiagnosticFieldWidths = map[string]int{"flags": 1, "alternatives": 6}
	case "four-byte":
		b51ADiagnosticFieldWidths = map[string]int{"flags": 1, "value": 4, "reserved": 3}
	case "save":
		b51ADiagnosticFieldWidth, b51ADiagnosticFieldOK := b51ADiagnosticFieldUint(b51ADiagnosticFieldArguments["value_width"])
		if !b51ADiagnosticFieldOK {
			return nil, b51ADiagnosticFieldFmt.Errorf("save requires value_width")
		}
		b51ADiagnosticFieldWidths = map[string]int{"status": 1, "mirror": 2, "symbol": 2, "value": int(b51ADiagnosticFieldWidth)}
	}
	b51ADiagnosticFieldOutput["state"] = "complete"
	b51ADiagnosticFieldFlagArray, _ := b51ADiagnosticFieldFields["flags"].([]any)
	b51ADiagnosticFieldFlagNumber := uint64(0)
	b51ADiagnosticFieldFlagsKnown := len(b51ADiagnosticFieldFlagArray) == 1
	if b51ADiagnosticFieldFlagsKnown {
		b51ADiagnosticFieldFlagNumber, _ = b51ADiagnosticFieldUint(b51ADiagnosticFieldFlagArray[0])
	}
	for b51ADiagnosticFieldKey, b51ADiagnosticFieldWidth := range b51ADiagnosticFieldWidths {
		b51ADiagnosticFieldResult := map[string]any{"state": "incomplete"}
		b51ADiagnosticFieldResults[b51ADiagnosticFieldKey] = b51ADiagnosticFieldResult
		b51ADiagnosticFieldRaw, b51ADiagnosticFieldPresent := b51ADiagnosticFieldFields[b51ADiagnosticFieldKey]
		if b51ADiagnosticFieldPresent {
			b51ADiagnosticFieldResult["raw"] = b51ADiagnosticFieldRaw
		}
		b51ADiagnosticFieldArray, _ := b51ADiagnosticFieldRaw.([]any)
		if !b51ADiagnosticFieldPresent && b51ADiagnosticFieldOperation == 5 && b51ADiagnosticFieldKey != "flags" && b51ADiagnosticFieldFlagsKnown && b51ADiagnosticFieldFlagNumber&1 != 0 {
			b51ADiagnosticFieldResult["state"] = "Unknown"
			continue
		}
		if !b51ADiagnosticFieldPresent || len(b51ADiagnosticFieldArray) != b51ADiagnosticFieldWidth {
			b51ADiagnosticFieldOutput["state"] = "partial"
			continue
		}
		if b51ADiagnosticFieldOperation == 5 && b51ADiagnosticFieldKey != "flags" && (!b51ADiagnosticFieldFlagsKnown || b51ADiagnosticFieldFlagNumber&1 != 0) {
			b51ADiagnosticFieldResult["state"] = "Unknown"
			continue
		}
		b51ADiagnosticFieldResult["state"] = "Hypothesis"
		if b51ADiagnosticFieldKey == "enable" {
			b51ADiagnosticFieldBits := make([]bool, 104)
			for b51ADiagnosticFieldIndex, b51ADiagnosticFieldByte := range b51ADiagnosticFieldArray {
				b51ADiagnosticFieldNumber, _ := b51ADiagnosticFieldUint(b51ADiagnosticFieldByte)
				for b51ADiagnosticFieldBit := 0; b51ADiagnosticFieldBit < 8; b51ADiagnosticFieldBit++ {
					b51ADiagnosticFieldBits[b51ADiagnosticFieldIndex*8+b51ADiagnosticFieldBit] = b51ADiagnosticFieldNumber&(1<<b51ADiagnosticFieldBit) != 0
				}
			}
			b51ADiagnosticFieldResult["flags"] = b51ADiagnosticFieldBits
			b51ADiagnosticFieldResult["calculation"] = "Proven"
			continue
		}
		if b51ADiagnosticFieldKey == "alternatives" {
			b51ADiagnosticFieldResult["alternatives"] = b51ADiagnosticFieldArray
			b51ADiagnosticFieldResult["calculation"] = "Proven"
			continue
		}
		if b51ADiagnosticFieldKey == "reserved" {
			b51ADiagnosticFieldResult["reserved"] = b51ADiagnosticFieldArray
			continue
		}
		b51ADiagnosticFieldNumber := uint64(0)
		for b51ADiagnosticFieldIndex, b51ADiagnosticFieldByte := range b51ADiagnosticFieldArray {
			b51ADiagnosticFieldPart, _ := b51ADiagnosticFieldUint(b51ADiagnosticFieldByte)
			b51ADiagnosticFieldNumber |= b51ADiagnosticFieldPart << (8 * b51ADiagnosticFieldIndex)
		}
		b51ADiagnosticFieldResult["number"] = b51ADiagnosticFieldNumber
		b51ADiagnosticFieldResult["calculation"] = "Proven"
		if b51ADiagnosticFieldKey == "flags" {
			b51ADiagnosticFieldResult["flags"] = map[string]bool{"disabled": b51ADiagnosticFieldNumber&1 != 0, "writable": b51ADiagnosticFieldNumber&2 != 0, "enum-versus-range": b51ADiagnosticFieldNumber&4 != 0, "additional-four-byte/no-range": b51ADiagnosticFieldNumber&8 != 0}
			b51ADiagnosticFieldResult["reserved"] = (b51ADiagnosticFieldNumber >> 4) & 3
			b51ADiagnosticFieldResult["unmodeled"] = (b51ADiagnosticFieldNumber >> 6) & 3
			b51ADiagnosticFieldResult["candidate_domain"] = "text alternatives"
			if b51ADiagnosticFieldNumber&8 != 0 {
				b51ADiagnosticFieldResult["alternatives"] = []any{map[string]any{"label": "additional-four-byte/no-range", "state": "Hypothesis"}, map[string]any{"label": "read-only-overrides-writable", "state": "Hypothesis"}}
			}
		}
		if b51ADiagnosticFieldKey == "status" {
			b51ADiagnosticFieldLabels := []string{"value-not-saved", "point-not-found", "write-protected", "range-violation", "misconfiguration"}
			b51ADiagnosticFieldResult["label"] = "unmapped"
			if b51ADiagnosticFieldNumber < uint64(len(b51ADiagnosticFieldLabels)) {
				b51ADiagnosticFieldResult["label"] = b51ADiagnosticFieldLabels[b51ADiagnosticFieldNumber]
			}
		}
		if b51ADiagnosticFieldView == "save" && b51ADiagnosticFieldKey == "value" {
			b51ADiagnosticFieldResult["candidate_domain"] = "caller-selected width"
			b51ADiagnosticFieldResult["alternatives"] = []int{1, 2}
		}
	}
	return b51ADiagnosticFieldJSON.Marshal(b51ADiagnosticFieldOutput)
}
