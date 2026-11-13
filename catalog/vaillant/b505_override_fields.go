package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

const b505OverrideFieldMaxDepth = 64

func b505OverrideFieldValidate(b505OverrideFieldDecoder *json.Decoder, b505OverrideFieldDepth int) error {
	if b505OverrideFieldDepth > b505OverrideFieldMaxDepth {
		return fmt.Errorf("nesting too deep")
	}
	b505OverrideFieldToken, b505OverrideFieldErr := b505OverrideFieldDecoder.Token()
	if b505OverrideFieldErr != nil {
		return b505OverrideFieldErr
	}
	if b505OverrideFieldDelimiter, b505OverrideFieldOK := b505OverrideFieldToken.(json.Delim); b505OverrideFieldOK {
		switch b505OverrideFieldDelimiter {
		case '{':
			b505OverrideFieldSeen := map[string]bool{}
			for b505OverrideFieldDecoder.More() {
				b505OverrideFieldKey, b505OverrideFieldErr := b505OverrideFieldDecoder.Token()
				if b505OverrideFieldErr != nil {
					return b505OverrideFieldErr
				}
				b505OverrideFieldName, b505OverrideFieldOK := b505OverrideFieldKey.(string)
				if !b505OverrideFieldOK || b505OverrideFieldSeen[b505OverrideFieldName] {
					return fmt.Errorf("duplicate or invalid object key")
				}
				b505OverrideFieldSeen[b505OverrideFieldName] = true
				if b505OverrideFieldErr := b505OverrideFieldValidate(b505OverrideFieldDecoder, b505OverrideFieldDepth+1); b505OverrideFieldErr != nil {
					return b505OverrideFieldErr
				}
			}
		case '[':
			for b505OverrideFieldDecoder.More() {
				if b505OverrideFieldErr := b505OverrideFieldValidate(b505OverrideFieldDecoder, b505OverrideFieldDepth+1); b505OverrideFieldErr != nil {
					return b505OverrideFieldErr
				}
			}
		default:
			return fmt.Errorf("unexpected delimiter")
		}
		_, b505OverrideFieldErr = b505OverrideFieldDecoder.Token()
		return b505OverrideFieldErr
	}
	if b505OverrideFieldNumber, b505OverrideFieldOK := b505OverrideFieldToken.(json.Number); b505OverrideFieldOK {
		if _, b505OverrideFieldErr := strconv.ParseFloat(string(b505OverrideFieldNumber), 64); b505OverrideFieldErr != nil {
			return fmt.Errorf("nonfinite number")
		}
	}
	return nil
}

func b505OverrideFieldObject(b505OverrideFieldRaw json.RawMessage) (map[string]json.RawMessage, error) {
	if len(bytes.TrimSpace(b505OverrideFieldRaw)) == 0 || bytes.TrimSpace(b505OverrideFieldRaw)[0] != '{' {
		return nil, fmt.Errorf("expected object")
	}
	var b505OverrideFieldObject map[string]json.RawMessage
	b505OverrideFieldErr := json.Unmarshal(b505OverrideFieldRaw, &b505OverrideFieldObject)
	return b505OverrideFieldObject, b505OverrideFieldErr
}

func b505OverrideFieldByte(b505OverrideFieldRaw json.RawMessage) (int, error) {
	b505OverrideFieldNumber, b505OverrideFieldErr := strconv.Atoi(string(b505OverrideFieldRaw))
	if b505OverrideFieldErr != nil || string(b505OverrideFieldRaw) == "-0" || b505OverrideFieldNumber < 0 || b505OverrideFieldNumber > 255 {
		return 0, fmt.Errorf("expected integer in 0..255")
	}
	return b505OverrideFieldNumber, nil
}

func B505OverrideFieldEvaluate(b505OverrideFieldInputRaw json.RawMessage) (json.RawMessage, error) {
	b505OverrideFieldDecoder := json.NewDecoder(bytes.NewReader(b505OverrideFieldInputRaw))
	b505OverrideFieldDecoder.UseNumber()
	if b505OverrideFieldErr := b505OverrideFieldValidate(b505OverrideFieldDecoder, 0); b505OverrideFieldErr != nil {
		return nil, b505OverrideFieldErr
	}
	if _, b505OverrideFieldErr := b505OverrideFieldDecoder.Token(); b505OverrideFieldErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b505OverrideFieldInput, b505OverrideFieldErr := b505OverrideFieldObject(b505OverrideFieldInputRaw)
	if b505OverrideFieldErr != nil {
		return nil, b505OverrideFieldErr
	}
	b505OverrideFieldContext, b505OverrideFieldErr := b505OverrideFieldObject(b505OverrideFieldInput["context"])
	if b505OverrideFieldErr != nil {
		return nil, b505OverrideFieldErr
	}
	b505OverrideFieldFields, b505OverrideFieldErr := b505OverrideFieldObject(b505OverrideFieldInput["fields"])
	if b505OverrideFieldErr != nil {
		return nil, b505OverrideFieldErr
	}
	b505OverrideFieldNumbers := map[string]int{}
	for _, b505OverrideFieldKey := range []string{"primary", "secondary", "operation"} {
		if b505OverrideFieldRaw, b505OverrideFieldOK := b505OverrideFieldContext[b505OverrideFieldKey]; b505OverrideFieldOK {
			b505OverrideFieldNumber, b505OverrideFieldErr := b505OverrideFieldByte(bytes.TrimSpace(b505OverrideFieldRaw))
			if b505OverrideFieldErr != nil {
				return nil, b505OverrideFieldErr
			}
			b505OverrideFieldNumbers[b505OverrideFieldKey] = b505OverrideFieldNumber
		}
	}
	b505OverrideFieldStrings := map[string]string{}
	for _, b505OverrideFieldKey := range []string{"direction", "view", "family", "variant"} {
		if b505OverrideFieldRaw, b505OverrideFieldOK := b505OverrideFieldContext[b505OverrideFieldKey]; b505OverrideFieldOK {
			if bytes.Equal(bytes.TrimSpace(b505OverrideFieldRaw), []byte("null")) {
				return nil, fmt.Errorf("expected identifier string")
			}
			var b505OverrideFieldValue string
			if b505OverrideFieldErr := json.Unmarshal(b505OverrideFieldRaw, &b505OverrideFieldValue); b505OverrideFieldErr != nil {
				return nil, b505OverrideFieldErr
			}
			b505OverrideFieldStrings[b505OverrideFieldKey] = b505OverrideFieldValue
		}
	}
	b505OverrideFieldBytes := map[string][]int{}
	for b505OverrideFieldKey, b505OverrideFieldRaw := range b505OverrideFieldFields {
		if len(bytes.TrimSpace(b505OverrideFieldRaw)) == 0 || bytes.TrimSpace(b505OverrideFieldRaw)[0] != '[' {
			return nil, fmt.Errorf("expected byte array")
		}
		var b505OverrideFieldArray []json.RawMessage
		if b505OverrideFieldErr := json.Unmarshal(b505OverrideFieldRaw, &b505OverrideFieldArray); b505OverrideFieldErr != nil {
			return nil, b505OverrideFieldErr
		}
		for _, b505OverrideFieldElement := range b505OverrideFieldArray {
			b505OverrideFieldNumber, b505OverrideFieldErr := b505OverrideFieldByte(b505OverrideFieldElement)
			if b505OverrideFieldErr != nil {
				return nil, b505OverrideFieldErr
			}
			b505OverrideFieldBytes[b505OverrideFieldKey] = append(b505OverrideFieldBytes[b505OverrideFieldKey], b505OverrideFieldNumber)
		}
	}
	b505OverrideFieldArguments := map[string]json.RawMessage{}
	if b505OverrideFieldRaw, b505OverrideFieldOK := b505OverrideFieldInput["arguments"]; b505OverrideFieldOK {
		b505OverrideFieldArguments, b505OverrideFieldErr = b505OverrideFieldObject(b505OverrideFieldRaw)
		if b505OverrideFieldErr != nil {
			return nil, b505OverrideFieldErr
		}
	}
	b505OverrideFieldDisplay := 0
	for b505OverrideFieldKey, b505OverrideFieldRaw := range b505OverrideFieldArguments {
		if b505OverrideFieldKey != "display_percent" {
			return nil, fmt.Errorf("unknown argument")
		}
		b505OverrideFieldDisplay, b505OverrideFieldErr = b505OverrideFieldByte(bytes.TrimSpace(b505OverrideFieldRaw))
		if b505OverrideFieldErr != nil {
			return nil, b505OverrideFieldErr
		}
	}
	b505OverrideFieldOperation := b505OverrideFieldNumbers["operation"]
	b505OverrideFieldOverride := b505OverrideFieldStrings["view"] == "override" && (b505OverrideFieldOperation == 77 || b505OverrideFieldOperation == 97)
	b505OverrideFieldPriority := b505OverrideFieldStrings["view"] == "priority" && b505OverrideFieldOperation == 78
	b505OverrideFieldSupported := b505OverrideFieldNumbers["primary"] == 181 && b505OverrideFieldNumbers["secondary"] == 5 && b505OverrideFieldStrings["direction"] == "request" && (b505OverrideFieldOverride || b505OverrideFieldPriority)
	b505OverrideFieldOutput := map[string]any{}
	for _, b505OverrideFieldKey := range []string{"physical", "qualified", "state", "results", "grade", "candidate_validity_seconds"} {
		if _, b505OverrideFieldOK := b505OverrideFieldInput[b505OverrideFieldKey]; b505OverrideFieldOK {
			return nil, fmt.Errorf("input key %s collides with output key", b505OverrideFieldKey)
		}
	}
	for b505OverrideFieldKey, b505OverrideFieldRaw := range b505OverrideFieldInput {
		b505OverrideFieldOutput[b505OverrideFieldKey] = b505OverrideFieldRaw
	}
	b505OverrideFieldOutput["physical"] = "Unknown"
	b505OverrideFieldOutput["qualified"] = false
	b505OverrideFieldOutput["state"] = "Unknown"
	b505OverrideFieldResults := map[string]any{}
	for b505OverrideFieldKey, b505OverrideFieldRaw := range b505OverrideFieldFields {
		b505OverrideFieldResults[b505OverrideFieldKey] = map[string]any{"raw": b505OverrideFieldRaw, "state": "Unknown"}
	}
	b505OverrideFieldOutput["results"] = b505OverrideFieldResults
	if !b505OverrideFieldSupported {
		return json.Marshal(b505OverrideFieldOutput)
	}
	b505OverrideFieldOutput["state"] = "complete"
	b505OverrideFieldOutput["grade"] = "Hypothesis"
	if b505OverrideFieldOverride {
		b505OverrideFieldOutput["candidate_validity_seconds"] = 900
	}
	b505OverrideFieldRequired := []string{}
	if b505OverrideFieldPriority {
		b505OverrideFieldRequired = []string{"priority"}
	}
	if b505OverrideFieldOperation == 77 {
		if b505OverrideFieldRaw, b505OverrideFieldOK := b505OverrideFieldArguments["display_percent"]; b505OverrideFieldOK {
			b505OverrideFieldResults["display_percent"] = map[string]any{"raw": b505OverrideFieldRaw, "state": "complete", "grade": "Hypothesis", "number": b505OverrideFieldDisplay, "candidate_domain": b505OverrideFieldDisplay == 0 || b505OverrideFieldDisplay == 25 || b505OverrideFieldDisplay == 50 || b505OverrideFieldDisplay == 75 || b505OverrideFieldDisplay == 100}
		}
		b505OverrideFieldRequired = []string{"sensor_status", "remote_control"}
	}
	if b505OverrideFieldOperation != 77 && len(b505OverrideFieldArguments) != 0 {
		return nil, fmt.Errorf("arguments not supported for view")
	}
	for _, b505OverrideFieldKey := range b505OverrideFieldRequired {
		b505OverrideFieldResult := map[string]any{"state": "incomplete", "grade": "Hypothesis"}
		if b505OverrideFieldRaw, b505OverrideFieldOK := b505OverrideFieldFields[b505OverrideFieldKey]; b505OverrideFieldOK {
			b505OverrideFieldResult["raw"] = b505OverrideFieldRaw
		}
		b505OverrideFieldResults[b505OverrideFieldKey] = b505OverrideFieldResult
		if len(b505OverrideFieldBytes[b505OverrideFieldKey]) != 1 {
			b505OverrideFieldOutput["state"] = "partial"
			continue
		}
		b505OverrideFieldNumber := b505OverrideFieldBytes[b505OverrideFieldKey][0]
		b505OverrideFieldResult["state"] = "complete"
		if b505OverrideFieldKey == "priority" {
			b505OverrideFieldResult["number"] = b505OverrideFieldNumber
			b505OverrideFieldResult["calculation"] = "Proven"
			continue
		}
		b505OverrideFieldResult["label"] = "unmapped"
		if b505OverrideFieldKey == "sensor_status" {
			if b505OverrideFieldNumber == 85 {
				b505OverrideFieldResult["label"] = "short-circuit"
			}
			if b505OverrideFieldNumber == 170 {
				b505OverrideFieldResult["label"] = "open"
			}
		}
		if b505OverrideFieldKey == "remote_control" && b505OverrideFieldNumber == 1 {
			b505OverrideFieldResult["label"] = "enabled"
		}
	}
	return json.Marshal(b505OverrideFieldOutput)
}
