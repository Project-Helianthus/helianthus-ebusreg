package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func b504CandidateFieldJSON(b504CandidateFieldDecoder *json.Decoder) error {
	b504CandidateFieldToken, b504CandidateFieldErr := b504CandidateFieldDecoder.Token()
	if b504CandidateFieldErr != nil {
		return b504CandidateFieldErr
	}
	if b504CandidateFieldNumber, b504CandidateFieldOK := b504CandidateFieldToken.(json.Number); b504CandidateFieldOK {
		if _, b504CandidateFieldErr = strconv.ParseFloat(string(b504CandidateFieldNumber), 64); b504CandidateFieldErr != nil {
			return b504CandidateFieldErr
		}
	}
	if b504CandidateFieldDelimiter, b504CandidateFieldOK := b504CandidateFieldToken.(json.Delim); b504CandidateFieldOK {
		switch b504CandidateFieldDelimiter {
		case '{':
			b504CandidateFieldSeen := map[string]bool{}
			for b504CandidateFieldDecoder.More() {
				b504CandidateFieldKey, b504CandidateFieldErr := b504CandidateFieldDecoder.Token()
				if b504CandidateFieldErr != nil {
					return b504CandidateFieldErr
				}
				b504CandidateFieldName, b504CandidateFieldOK := b504CandidateFieldKey.(string)
				if !b504CandidateFieldOK || b504CandidateFieldSeen[b504CandidateFieldName] {
					return fmt.Errorf("duplicate or invalid object key")
				}
				b504CandidateFieldSeen[b504CandidateFieldName] = true
				if b504CandidateFieldErr = b504CandidateFieldJSON(b504CandidateFieldDecoder); b504CandidateFieldErr != nil {
					return b504CandidateFieldErr
				}
			}
		case '[':
			for b504CandidateFieldDecoder.More() {
				if b504CandidateFieldErr = b504CandidateFieldJSON(b504CandidateFieldDecoder); b504CandidateFieldErr != nil {
					return b504CandidateFieldErr
				}
			}
		default:
			return fmt.Errorf("unexpected delimiter")
		}
		_, b504CandidateFieldErr = b504CandidateFieldDecoder.Token()
	}
	return b504CandidateFieldErr
}

func b504CandidateFieldObject(b504CandidateFieldRaw json.RawMessage) (map[string]json.RawMessage, error) {
	b504CandidateFieldTrimmed := bytes.TrimSpace(b504CandidateFieldRaw)
	if len(b504CandidateFieldTrimmed) == 0 || b504CandidateFieldTrimmed[0] != '{' {
		return nil, fmt.Errorf("expected object")
	}
	var b504CandidateFieldValue map[string]json.RawMessage
	b504CandidateFieldErr := json.Unmarshal(b504CandidateFieldRaw, &b504CandidateFieldValue)
	return b504CandidateFieldValue, b504CandidateFieldErr
}

func b504CandidateFieldInteger(b504CandidateFieldRaw json.RawMessage) (int, error) {
	b504CandidateFieldValue, b504CandidateFieldErr := strconv.Atoi(string(b504CandidateFieldRaw))
	if b504CandidateFieldErr != nil || b504CandidateFieldValue < 0 || b504CandidateFieldValue > 255 {
		return 0, fmt.Errorf("expected integer in 0..255")
	}
	return b504CandidateFieldValue, nil
}

func b504CandidateFieldString(b504CandidateFieldRaw json.RawMessage) (string, error) {
	if len(b504CandidateFieldRaw) == 0 || b504CandidateFieldRaw[0] != '"' {
		return "", fmt.Errorf("expected string")
	}
	var b504CandidateFieldValue string
	b504CandidateFieldErr := json.Unmarshal(b504CandidateFieldRaw, &b504CandidateFieldValue)
	return b504CandidateFieldValue, b504CandidateFieldErr
}

func B504CandidateFieldEvaluate(b504CandidateFieldInputRaw json.RawMessage) (json.RawMessage, error) {
	b504CandidateFieldDecoder := json.NewDecoder(bytes.NewReader(b504CandidateFieldInputRaw))
	b504CandidateFieldDecoder.UseNumber()
	if b504CandidateFieldErr := b504CandidateFieldJSON(b504CandidateFieldDecoder); b504CandidateFieldErr != nil {
		return nil, b504CandidateFieldErr
	}
	if _, b504CandidateFieldErr := b504CandidateFieldDecoder.Token(); b504CandidateFieldErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b504CandidateFieldInput, b504CandidateFieldErr := b504CandidateFieldObject(b504CandidateFieldInputRaw)
	if b504CandidateFieldErr != nil {
		return nil, b504CandidateFieldErr
	}
	b504CandidateFieldContext, b504CandidateFieldErr := b504CandidateFieldObject(b504CandidateFieldInput["context"])
	if b504CandidateFieldErr != nil {
		return nil, b504CandidateFieldErr
	}
	b504CandidateFieldFields, b504CandidateFieldErr := b504CandidateFieldObject(b504CandidateFieldInput["fields"])
	if b504CandidateFieldErr != nil {
		return nil, b504CandidateFieldErr
	}
	b504CandidateFieldArguments := map[string]json.RawMessage{}
	if b504CandidateFieldRaw, b504CandidateFieldPresent := b504CandidateFieldInput["arguments"]; b504CandidateFieldPresent {
		b504CandidateFieldArguments, b504CandidateFieldErr = b504CandidateFieldObject(b504CandidateFieldRaw)
		if b504CandidateFieldErr != nil {
			return nil, b504CandidateFieldErr
		}
	}
	b504CandidateFieldIntegers := map[string]int{}
	for _, b504CandidateFieldKey := range []string{"primary", "secondary", "operation"} {
		if b504CandidateFieldRaw, b504CandidateFieldPresent := b504CandidateFieldContext[b504CandidateFieldKey]; b504CandidateFieldPresent {
			b504CandidateFieldIntegers[b504CandidateFieldKey], b504CandidateFieldErr = b504CandidateFieldInteger(b504CandidateFieldRaw)
			if b504CandidateFieldErr != nil {
				return nil, b504CandidateFieldErr
			}
		}
	}
	b504CandidateFieldStrings := map[string]string{}
	for _, b504CandidateFieldKey := range []string{"direction", "view", "family", "variant"} {
		if b504CandidateFieldRaw, b504CandidateFieldPresent := b504CandidateFieldContext[b504CandidateFieldKey]; b504CandidateFieldPresent {
			b504CandidateFieldStrings[b504CandidateFieldKey], b504CandidateFieldErr = b504CandidateFieldString(b504CandidateFieldRaw)
			if b504CandidateFieldErr != nil {
				return nil, b504CandidateFieldErr
			}
		}
	}
	b504CandidateFieldBytes := map[string][]int{}
	b504CandidateFieldResults := map[string]any{}
	for b504CandidateFieldKey, b504CandidateFieldRaw := range b504CandidateFieldFields {
		if len(b504CandidateFieldRaw) == 0 || b504CandidateFieldRaw[0] != '[' {
			return nil, fmt.Errorf("field must be byte array")
		}
		var b504CandidateFieldArray []json.RawMessage
		if b504CandidateFieldErr = json.Unmarshal(b504CandidateFieldRaw, &b504CandidateFieldArray); b504CandidateFieldErr != nil {
			return nil, b504CandidateFieldErr
		}
		for _, b504CandidateFieldElement := range b504CandidateFieldArray {
			b504CandidateFieldByte, b504CandidateFieldErr := b504CandidateFieldInteger(b504CandidateFieldElement)
			if b504CandidateFieldErr != nil {
				return nil, b504CandidateFieldErr
			}
			b504CandidateFieldBytes[b504CandidateFieldKey] = append(b504CandidateFieldBytes[b504CandidateFieldKey], b504CandidateFieldByte)
		}
		b504CandidateFieldResults[b504CandidateFieldKey] = map[string]any{"raw": b504CandidateFieldRaw, "state": "Unknown", "present": true}
	}
	b504CandidateFieldDivisor := 0
	b504CandidateFieldRole := ""
	for b504CandidateFieldKey, b504CandidateFieldRaw := range b504CandidateFieldArguments {
		switch b504CandidateFieldKey {
		case "divisor":
			b504CandidateFieldDivisor, b504CandidateFieldErr = strconv.Atoi(string(b504CandidateFieldRaw))
			if b504CandidateFieldDivisor < 0 || b504CandidateFieldDivisor > 256 {
				return nil, fmt.Errorf("invalid divisor")
			}
		case "role":
			b504CandidateFieldRole, b504CandidateFieldErr = b504CandidateFieldString(b504CandidateFieldRaw)
		}
		if b504CandidateFieldErr != nil {
			return nil, b504CandidateFieldErr
		}
	}
	b504CandidateFieldView := b504CandidateFieldStrings["view"]
	b504CandidateFieldOperations := map[string]int{"corrected-outside": 22, "demand": 32, "protection": 33, "adjustment": 38, "circuit-role": 13}
	b504CandidateFieldOperation, b504CandidateFieldKnown := b504CandidateFieldOperations[b504CandidateFieldView]
	_, b504CandidateFieldHasOperation := b504CandidateFieldContext["operation"]
	b504CandidateFieldSupported := b504CandidateFieldKnown && b504CandidateFieldHasOperation && b504CandidateFieldIntegers["primary"] == 181 && b504CandidateFieldIntegers["secondary"] == 4 && b504CandidateFieldIntegers["operation"] == b504CandidateFieldOperation && b504CandidateFieldStrings["direction"] == "positive"
	b504CandidateFieldState := "Unknown"
	if b504CandidateFieldSupported {
		b504CandidateFieldArgument := ""
		if b504CandidateFieldView == "corrected-outside" {
			b504CandidateFieldArgument = "divisor"
			if b504CandidateFieldDivisor != 16 && b504CandidateFieldDivisor != 256 {
				return nil, fmt.Errorf("divisor must be explicitly 16 or 256")
			}
		}
		if b504CandidateFieldView == "circuit-role" {
			b504CandidateFieldArgument = "role"
			if _, b504CandidateFieldPresent := b504CandidateFieldArguments["role"]; !b504CandidateFieldPresent {
				return nil, fmt.Errorf("missing role")
			}
		}
		for b504CandidateFieldKey := range b504CandidateFieldArguments {
			if b504CandidateFieldKey != b504CandidateFieldArgument {
				return nil, fmt.Errorf("unexpected argument")
			}
		}
		b504CandidateFieldWidths := map[string]int{}
		switch b504CandidateFieldView {
		case "corrected-outside":
			b504CandidateFieldWidths["temperature"] = 2
		case "demand":
			b504CandidateFieldWidths["setpoint"] = 1
		case "protection", "circuit-role":
			b504CandidateFieldWidths["temperature"] = 1
		case "adjustment":
			b504CandidateFieldWidths["adjustment"] = 1
			b504CandidateFieldWidths["room"] = 1
		}
		if b504CandidateFieldView == "circuit-role" {
			if _, b504CandidateFieldPresent := b504CandidateFieldFields["dhw_request"]; b504CandidateFieldPresent {
				b504CandidateFieldWidths["dhw_request"] = 1
			}
		}
		b504CandidateFieldState = "complete"
		for b504CandidateFieldKey, b504CandidateFieldWidth := range b504CandidateFieldWidths {
			b504CandidateFieldRaw, b504CandidateFieldPresent := b504CandidateFieldFields[b504CandidateFieldKey]
			b504CandidateFieldResult := map[string]any{"raw": b504CandidateFieldRaw, "present": b504CandidateFieldPresent, "state": "incomplete"}
			b504CandidateFieldResults[b504CandidateFieldKey] = b504CandidateFieldResult
			b504CandidateFieldData := b504CandidateFieldBytes[b504CandidateFieldKey]
			if !b504CandidateFieldPresent || len(b504CandidateFieldData) != b504CandidateFieldWidth {
				b504CandidateFieldState = "partial"
				continue
			}
			b504CandidateFieldResult["state"] = "complete"
			b504CandidateFieldResult["grade"] = "Hypothesis"
			switch b504CandidateFieldView {
			case "corrected-outside":
				b504CandidateFieldResult["number"] = float64(int16(uint16(b504CandidateFieldData[0])|uint16(b504CandidateFieldData[1])<<8)) / float64(b504CandidateFieldDivisor)
				b504CandidateFieldResult["calculation"] = "Proven"
			case "demand":
				if b504CandidateFieldData[0] == 0 {
					b504CandidateFieldResult["label"] = "no-demand"
				} else {
					b504CandidateFieldResult["label"] = "offset-included-candidate"
				}
			case "protection":
				if b504CandidateFieldData[0] == 0 {
					b504CandidateFieldResult["label"] = "off"
				} else {
					b504CandidateFieldResult["label"] = "protection-active"
				}
			case "adjustment":
				if b504CandidateFieldKey == "adjustment" && b504CandidateFieldData[0] == 1 {
					b504CandidateFieldResult["label"] = "change-permitted"
				}
				if b504CandidateFieldKey == "room" && b504CandidateFieldData[0] == 255 {
					b504CandidateFieldResult["label"] = "unavailable-candidate"
				}
			case "circuit-role":
				b504CandidateFieldLabels := map[string]string{"burner": "common-sensor", "mixer": "flow", "dhw": "tank"}
				if b504CandidateFieldLabel, b504CandidateFieldOK := b504CandidateFieldLabels[b504CandidateFieldRole]; b504CandidateFieldOK && b504CandidateFieldKey == "temperature" {
					b504CandidateFieldResult["label"] = b504CandidateFieldLabel
					b504CandidateFieldResult["number"] = b504CandidateFieldData[0]
				} else if b504CandidateFieldKey == "dhw_request" && b504CandidateFieldRole == "dhw" {
					b504CandidateFieldResult["label"] = "dhw"
				} else {
					b504CandidateFieldResult["state"] = "Unknown"
					delete(b504CandidateFieldResult, "grade")
				}
			}
		}
	}
	b504CandidateFieldOutput := map[string]any{}
	for b504CandidateFieldKey, b504CandidateFieldRaw := range b504CandidateFieldInput {
		b504CandidateFieldOutput[b504CandidateFieldKey] = b504CandidateFieldRaw
	}
	b504CandidateFieldOutput["state"] = b504CandidateFieldState
	b504CandidateFieldOutput["results"] = b504CandidateFieldResults
	b504CandidateFieldOutput["physical"] = "Unknown"
	b504CandidateFieldOutput["qualified"] = false
	return json.Marshal(b504CandidateFieldOutput)
}
