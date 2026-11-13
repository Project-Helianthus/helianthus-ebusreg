package vaillant

import (
	b503ControlFieldBytesPackage "bytes"
	b503ControlFieldJSON "encoding/json"
	b503ControlFieldFmt "fmt"
	b503ControlFieldIO "io"
	b503ControlFieldStrconv "strconv"
)

func b503ControlFieldDecode(b503ControlFieldDecoder *b503ControlFieldJSON.Decoder) (any, error) {
	b503ControlFieldToken, b503ControlFieldError := b503ControlFieldDecoder.Token()
	if b503ControlFieldError != nil {
		return nil, b503ControlFieldError
	}
	if b503ControlFieldDelimiter, b503ControlFieldOK := b503ControlFieldToken.(b503ControlFieldJSON.Delim); b503ControlFieldOK {
		switch b503ControlFieldDelimiter {
		case '{':
			b503ControlFieldObject := map[string]any{}
			for b503ControlFieldDecoder.More() {
				b503ControlFieldKeyToken, b503ControlFieldError := b503ControlFieldDecoder.Token()
				if b503ControlFieldError != nil {
					return nil, b503ControlFieldError
				}
				b503ControlFieldKey, b503ControlFieldOK := b503ControlFieldKeyToken.(string)
				if !b503ControlFieldOK {
					return nil, b503ControlFieldFmt.Errorf("invalid object key")
				}
				if _, b503ControlFieldExists := b503ControlFieldObject[b503ControlFieldKey]; b503ControlFieldExists {
					return nil, b503ControlFieldFmt.Errorf("duplicate key %q", b503ControlFieldKey)
				}
				b503ControlFieldValue, b503ControlFieldError := b503ControlFieldDecode(b503ControlFieldDecoder)
				if b503ControlFieldError != nil {
					return nil, b503ControlFieldError
				}
				b503ControlFieldObject[b503ControlFieldKey] = b503ControlFieldValue
			}
			_, b503ControlFieldError := b503ControlFieldDecoder.Token()
			return b503ControlFieldObject, b503ControlFieldError
		case '[':
			b503ControlFieldArray := []any{}
			for b503ControlFieldDecoder.More() {
				b503ControlFieldValue, b503ControlFieldError := b503ControlFieldDecode(b503ControlFieldDecoder)
				if b503ControlFieldError != nil {
					return nil, b503ControlFieldError
				}
				b503ControlFieldArray = append(b503ControlFieldArray, b503ControlFieldValue)
			}
			_, b503ControlFieldError := b503ControlFieldDecoder.Token()
			return b503ControlFieldArray, b503ControlFieldError
		}
		return nil, b503ControlFieldFmt.Errorf("unexpected delimiter")
	}
	return b503ControlFieldToken, nil
}

func b503ControlFieldUint(b503ControlFieldValue any, b503ControlFieldBits int) (uint64, error) {
	b503ControlFieldNumber, b503ControlFieldOK := b503ControlFieldValue.(b503ControlFieldJSON.Number)
	if !b503ControlFieldOK {
		return 0, b503ControlFieldFmt.Errorf("unsigned integer required")
	}
	return b503ControlFieldStrconv.ParseUint(string(b503ControlFieldNumber), 10, b503ControlFieldBits)
}

func B503ControlFieldEvaluate(b503ControlFieldInput b503ControlFieldJSON.RawMessage) (b503ControlFieldJSON.RawMessage, error) {
	b503ControlFieldDecoder := b503ControlFieldJSON.NewDecoder(b503ControlFieldBytesPackage.NewReader(b503ControlFieldInput))
	b503ControlFieldDecoder.UseNumber()
	b503ControlFieldValue, b503ControlFieldError := b503ControlFieldDecode(b503ControlFieldDecoder)
	if b503ControlFieldError != nil {
		return nil, b503ControlFieldError
	}
	if _, b503ControlFieldError = b503ControlFieldDecoder.Token(); b503ControlFieldError != b503ControlFieldIO.EOF {
		return nil, b503ControlFieldFmt.Errorf("trailing JSON")
	}
	b503ControlFieldOutput, b503ControlFieldOK := b503ControlFieldValue.(map[string]any)
	if !b503ControlFieldOK {
		return nil, b503ControlFieldFmt.Errorf("input object required")
	}
	for _, b503ControlFieldReserved := range []string{"results", "state", "physical", "qualified"} {
		if _, b503ControlFieldPresent := b503ControlFieldOutput[b503ControlFieldReserved]; b503ControlFieldPresent {
			return nil, b503ControlFieldFmt.Errorf("input key %s is reserved for output", b503ControlFieldReserved)
		}
	}
	b503ControlFieldContext, b503ControlFieldOK := b503ControlFieldOutput["context"].(map[string]any)
	if !b503ControlFieldOK {
		return nil, b503ControlFieldFmt.Errorf("context object required")
	}
	b503ControlFieldFields, b503ControlFieldOK := b503ControlFieldOutput["fields"].(map[string]any)
	if !b503ControlFieldOK {
		return nil, b503ControlFieldFmt.Errorf("fields object required")
	}
	for b503ControlFieldKey, b503ControlFieldValue := range b503ControlFieldFields {
		b503ControlFieldArray, b503ControlFieldOK := b503ControlFieldValue.([]any)
		if !b503ControlFieldOK {
			return nil, b503ControlFieldFmt.Errorf("field %s requires byte array", b503ControlFieldKey)
		}
		for _, b503ControlFieldByte := range b503ControlFieldArray {
			if _, b503ControlFieldError := b503ControlFieldUint(b503ControlFieldByte, 8); b503ControlFieldError != nil {
				return nil, b503ControlFieldError
			}
		}
	}
	if b503ControlFieldValue, b503ControlFieldPresent := b503ControlFieldOutput["arguments"]; b503ControlFieldPresent {
		b503ControlFieldArguments, b503ControlFieldOK := b503ControlFieldValue.(map[string]any)
		if !b503ControlFieldOK || len(b503ControlFieldArguments) != 0 {
			return nil, b503ControlFieldFmt.Errorf("no arguments defined for these views")
		}
	}
	b503ControlFieldSupported := true
	b503ControlFieldNumbers := map[string]uint64{}
	for _, b503ControlFieldKey := range []string{"primary", "secondary", "operation"} {
		b503ControlFieldValue, b503ControlFieldPresent := b503ControlFieldContext[b503ControlFieldKey]
		if !b503ControlFieldPresent {
			b503ControlFieldSupported = false
			continue
		}
		b503ControlFieldNumber, b503ControlFieldError := b503ControlFieldUint(b503ControlFieldValue, 8)
		if b503ControlFieldError != nil {
			return nil, b503ControlFieldError
		}
		b503ControlFieldNumbers[b503ControlFieldKey] = b503ControlFieldNumber
	}
	for _, b503ControlFieldKey := range []string{"direction", "view"} {
		if b503ControlFieldValue, b503ControlFieldPresent := b503ControlFieldContext[b503ControlFieldKey]; b503ControlFieldPresent {
			if _, b503ControlFieldOK := b503ControlFieldValue.(string); !b503ControlFieldOK {
				return nil, b503ControlFieldFmt.Errorf("context %s requires string", b503ControlFieldKey)
			}
		} else {
			b503ControlFieldSupported = false
		}
	}
	b503ControlFieldView, _ := b503ControlFieldContext["view"].(string)
	b503ControlFieldDirection, _ := b503ControlFieldContext["direction"].(string)
	b503ControlFieldOperation := b503ControlFieldNumbers["operation"]
	b503ControlFieldMatch := false
	switch b503ControlFieldView {
	case "reset-acceptance":
		b503ControlFieldMatch = (b503ControlFieldOperation == 3 || b503ControlFieldOperation == 6) && b503ControlFieldDirection == "positive"
	case "current-errors":
		b503ControlFieldMatch = b503ControlFieldOperation == 7 && b503ControlFieldDirection == "positive"
	case "category-common":
		b503ControlFieldMatch = b503ControlFieldOperation == 4 && b503ControlFieldDirection == "request"
	case "category-reset":
		b503ControlFieldMatch = (b503ControlFieldOperation == 3 || b503ControlFieldOperation == 6) && b503ControlFieldDirection == "request"
	}
	b503ControlFieldResults := map[string]any{}
	b503ControlFieldOutput["results"] = b503ControlFieldResults
	b503ControlFieldOutput["state"] = "Unknown"
	b503ControlFieldOutput["physical"] = "Unknown"
	b503ControlFieldOutput["qualified"] = false
	if !b503ControlFieldSupported || !b503ControlFieldMatch || b503ControlFieldNumbers["primary"] != 181 || b503ControlFieldNumbers["secondary"] != 3 {
		return b503ControlFieldJSON.Marshal(b503ControlFieldOutput)
	}
	b503ControlFieldOutput["state"] = "complete"
	b503ControlFieldWidths := map[string]int{}
	switch b503ControlFieldView {
	case "reset-acceptance":
		b503ControlFieldWidths["acceptance"] = 1
	case "current-errors":
		b503ControlFieldWidths["flags"] = 1
		b503ControlFieldWidths["feedback"] = 1
		for b503ControlFieldIndex := 0; b503ControlFieldIndex < 5; b503ControlFieldIndex++ {
			b503ControlFieldWidths[b503ControlFieldFmt.Sprintf("error%d", b503ControlFieldIndex)] = 2
		}
	default:
		b503ControlFieldWidths["category"] = 1
	}
	for b503ControlFieldKey, b503ControlFieldWidth := range b503ControlFieldWidths {
		b503ControlFieldRaw, b503ControlFieldPresent := b503ControlFieldFields[b503ControlFieldKey]
		b503ControlFieldResult := map[string]any{"raw": b503ControlFieldRaw, "present": b503ControlFieldPresent, "state": "incomplete", "grade": "Hypothesis"}
		b503ControlFieldResults[b503ControlFieldKey] = b503ControlFieldResult
		b503ControlFieldArray, _ := b503ControlFieldRaw.([]any)
		if !b503ControlFieldPresent || len(b503ControlFieldArray) != b503ControlFieldWidth {
			b503ControlFieldOutput["state"] = "partial"
			continue
		}
		b503ControlFieldNumber, _ := b503ControlFieldUint(b503ControlFieldArray[0], 8)
		if b503ControlFieldWidth == 2 {
			b503ControlFieldHigh, _ := b503ControlFieldUint(b503ControlFieldArray[1], 8)
			b503ControlFieldNumber += b503ControlFieldHigh << 8
		}
		b503ControlFieldResult["state"] = "Hypothesis"
		b503ControlFieldResult["number"] = b503ControlFieldNumber
		b503ControlFieldResult["calculation"] = "Proven"
		switch b503ControlFieldKey {
		case "acceptance":
			b503ControlFieldResult["default"] = 0
			if b503ControlFieldNumber <= 1 {
				b503ControlFieldResult["value"] = b503ControlFieldNumber == 1
			} else {
				b503ControlFieldResult["state"] = "Unknown"
				b503ControlFieldResult["unmodeled"] = b503ControlFieldNumber
			}
		case "feedback":
			b503ControlFieldResult["label"] = "rejected"
			if b503ControlFieldNumber == 1 {
				b503ControlFieldResult["label"] = "accepted"
			}
		case "flags":
			b503ControlFieldFlags := map[string]bool{}
			for b503ControlFieldIndex := 0; b503ControlFieldIndex < 5; b503ControlFieldIndex++ {
				b503ControlFieldFlags[b503ControlFieldFmt.Sprintf("bit%d", b503ControlFieldIndex)] = b503ControlFieldNumber&(1<<b503ControlFieldIndex) != 0
			}
			b503ControlFieldResult["flags"] = b503ControlFieldFlags
			b503ControlFieldResult["reserved"] = b503ControlFieldNumber & 0xe0
		case "category":
			b503ControlFieldResult["state"] = "Unknown"
			if b503ControlFieldView == "category-common" {
				b503ControlFieldResult["candidate_domain"] = []int{1, 2, 3, 4, 5}
				if b503ControlFieldNumber == 2 {
					b503ControlFieldResult["label"] = "predictive-maintenance"
					b503ControlFieldResult["state"] = "Hypothesis"
				}
			} else {
				if b503ControlFieldNumber == 1 {
					b503ControlFieldResult["label"] = "error"
					b503ControlFieldResult["state"] = "Hypothesis"
				}
				if b503ControlFieldNumber == 3 {
					b503ControlFieldResult["label"] = "status"
					b503ControlFieldResult["state"] = "Hypothesis"
				}
			}
		}
	}
	return b503ControlFieldJSON.Marshal(b503ControlFieldOutput)
}
