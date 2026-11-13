package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func b512EmergencyFlowRead(b512EmergencyFlowDecoder *json.Decoder) (any, error) {
	b512EmergencyFlowToken, b512EmergencyFlowErr := b512EmergencyFlowDecoder.Token()
	if b512EmergencyFlowErr != nil {
		return nil, b512EmergencyFlowErr
	}
	if b512EmergencyFlowDelimiter, b512EmergencyFlowOK := b512EmergencyFlowToken.(json.Delim); b512EmergencyFlowOK {
		switch b512EmergencyFlowDelimiter {
		case '{':
			b512EmergencyFlowObject := map[string]any{}
			for b512EmergencyFlowDecoder.More() {
				b512EmergencyFlowKeyToken, b512EmergencyFlowErr := b512EmergencyFlowDecoder.Token()
				if b512EmergencyFlowErr != nil {
					return nil, b512EmergencyFlowErr
				}
				b512EmergencyFlowKey := b512EmergencyFlowKeyToken.(string)
				if _, b512EmergencyFlowExists := b512EmergencyFlowObject[b512EmergencyFlowKey]; b512EmergencyFlowExists {
					return nil, fmt.Errorf("duplicate key %q", b512EmergencyFlowKey)
				}
				b512EmergencyFlowValue, b512EmergencyFlowErr := b512EmergencyFlowRead(b512EmergencyFlowDecoder)
				if b512EmergencyFlowErr != nil {
					return nil, b512EmergencyFlowErr
				}
				b512EmergencyFlowObject[b512EmergencyFlowKey] = b512EmergencyFlowValue
			}
			_, b512EmergencyFlowErr := b512EmergencyFlowDecoder.Token()
			return b512EmergencyFlowObject, b512EmergencyFlowErr
		case '[':
			b512EmergencyFlowArray := []any{}
			for b512EmergencyFlowDecoder.More() {
				b512EmergencyFlowValue, b512EmergencyFlowErr := b512EmergencyFlowRead(b512EmergencyFlowDecoder)
				if b512EmergencyFlowErr != nil {
					return nil, b512EmergencyFlowErr
				}
				b512EmergencyFlowArray = append(b512EmergencyFlowArray, b512EmergencyFlowValue)
			}
			_, b512EmergencyFlowErr := b512EmergencyFlowDecoder.Token()
			return b512EmergencyFlowArray, b512EmergencyFlowErr
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	return b512EmergencyFlowToken, nil
}

func B512EmergencyFlowEvaluate(b512EmergencyFlowInput json.RawMessage) (json.RawMessage, error) {
	b512EmergencyFlowDecoder := json.NewDecoder(bytes.NewReader(b512EmergencyFlowInput))
	b512EmergencyFlowDecoder.UseNumber()
	b512EmergencyFlowValue, b512EmergencyFlowErr := b512EmergencyFlowRead(b512EmergencyFlowDecoder)
	if b512EmergencyFlowErr != nil {
		return nil, b512EmergencyFlowErr
	}
	if _, b512EmergencyFlowErr = b512EmergencyFlowDecoder.Token(); b512EmergencyFlowErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b512EmergencyFlowRoot, b512EmergencyFlowOK := b512EmergencyFlowValue.(map[string]any)
	if !b512EmergencyFlowOK {
		return nil, fmt.Errorf("input must be object")
	}
	b512EmergencyFlowContext, b512EmergencyFlowOK := b512EmergencyFlowRoot["context"].(map[string]any)
	if !b512EmergencyFlowOK {
		return nil, fmt.Errorf("context must be object")
	}
	for _, b512EmergencyFlowKey := range []string{"primary", "secondary", "operation"} {
		if b512EmergencyFlowValue, b512EmergencyFlowExists := b512EmergencyFlowContext[b512EmergencyFlowKey]; b512EmergencyFlowExists && b512EmergencyFlowValue != nil {
			b512EmergencyFlowNumber, b512EmergencyFlowOK := b512EmergencyFlowValue.(json.Number)
			if !b512EmergencyFlowOK {
				return nil, fmt.Errorf("context number must be integer")
			}
			b512EmergencyFlowInteger, b512EmergencyFlowErr := b512EmergencyFlowNumber.Int64()
			if b512EmergencyFlowErr != nil || b512EmergencyFlowInteger < 0 || b512EmergencyFlowInteger > 255 {
				return nil, fmt.Errorf("context number out of range")
			}
		}
	}
	for _, b512EmergencyFlowKey := range []string{"direction", "view"} {
		if b512EmergencyFlowValue, b512EmergencyFlowExists := b512EmergencyFlowContext[b512EmergencyFlowKey]; b512EmergencyFlowExists && b512EmergencyFlowValue != nil {
			if _, b512EmergencyFlowOK := b512EmergencyFlowValue.(string); !b512EmergencyFlowOK {
				return nil, fmt.Errorf("context selector must be string")
			}
		}
	}
	b512EmergencyFlowRaw, b512EmergencyFlowOK := b512EmergencyFlowRoot["fields"].(map[string]any)
	if !b512EmergencyFlowOK {
		return nil, fmt.Errorf("fields must be object")
	}
	if b512EmergencyFlowArguments, b512EmergencyFlowExists := b512EmergencyFlowRoot["arguments"]; b512EmergencyFlowExists {
		if _, b512EmergencyFlowOK := b512EmergencyFlowArguments.(map[string]any); !b512EmergencyFlowOK {
			return nil, fmt.Errorf("arguments must be object")
		}
	}
	for _, b512EmergencyFlowValue := range b512EmergencyFlowRaw {
		b512EmergencyFlowArray, b512EmergencyFlowOK := b512EmergencyFlowValue.([]any)
		if !b512EmergencyFlowOK {
			return nil, fmt.Errorf("field must be byte array")
		}
		for _, b512EmergencyFlowByte := range b512EmergencyFlowArray {
			b512EmergencyFlowNumber, b512EmergencyFlowOK := b512EmergencyFlowByte.(json.Number)
			if !b512EmergencyFlowOK {
				return nil, fmt.Errorf("byte must be integer")
			}
			b512EmergencyFlowInteger, b512EmergencyFlowErr := b512EmergencyFlowNumber.Int64()
			if b512EmergencyFlowErr != nil || b512EmergencyFlowInteger < 0 || b512EmergencyFlowInteger > 255 {
				return nil, fmt.Errorf("byte out of range")
			}
		}
	}
	b512EmergencyFlowResults := map[string]any{}
	b512EmergencyFlowOutput := map[string]any{"context": b512EmergencyFlowContext, "raw": b512EmergencyFlowRaw, "state": "Unknown", "fields": b512EmergencyFlowResults, "physical": "Unknown", "qualified": false}
	if b512EmergencyFlowArguments, b512EmergencyFlowExists := b512EmergencyFlowRoot["arguments"]; b512EmergencyFlowExists {
		b512EmergencyFlowOutput["arguments"] = b512EmergencyFlowArguments
	}
	b512EmergencyFlowOperation := b512EmergencyFlowContext["operation"]
	b512EmergencyFlowDirection := b512EmergencyFlowContext["direction"]
	b512EmergencyFlowFlow := b512EmergencyFlowOperation == json.Number("17") && b512EmergencyFlowContext["view"] == "flow-io"
	b512EmergencyFlowEmergency := b512EmergencyFlowOperation == json.Number("18") && b512EmergencyFlowContext["view"] == "emergency-io"
	if b512EmergencyFlowContext["primary"] != json.Number("181") || b512EmergencyFlowContext["secondary"] != json.Number("18") || (!b512EmergencyFlowFlow && !b512EmergencyFlowEmergency) || (b512EmergencyFlowDirection != "request" && b512EmergencyFlowDirection != "positive") {
		return json.Marshal(b512EmergencyFlowOutput)
	}
	b512EmergencyFlowOutput["state"] = "complete"
	b512EmergencyFlowLengths := map[string]int{"flags": 1}
	if b512EmergencyFlowDirection == "positive" {
		b512EmergencyFlowOutput["candidate"] = map[string]any{"state": "Hypothesis", "declared_positive_prefix": []int{181, 18, 12}, "requested_operation": b512EmergencyFlowOperation, "literal_reply_conflict": true, "physical_association": "Unknown"}
		if b512EmergencyFlowFlow {
			b512EmergencyFlowLengths["setpoint"] = 1
		}
	} else if b512EmergencyFlowFlow {
		for _, b512EmergencyFlowName := range []string{"temperature1", "temperature2", "temperature3"} {
			b512EmergencyFlowLengths[b512EmergencyFlowName] = 2
		}
		for _, b512EmergencyFlowName := range []string{"pressure", "flow", "pump_status"} {
			b512EmergencyFlowLengths[b512EmergencyFlowName] = 1
		}
	} else {
		b512EmergencyFlowLengths["switch"] = 1
		b512EmergencyFlowLengths["co2_1"] = 2
		b512EmergencyFlowLengths["co2_2"] = 2
	}
	for b512EmergencyFlowName, b512EmergencyFlowRawValue := range b512EmergencyFlowRaw {
		b512EmergencyFlowResults[b512EmergencyFlowName] = map[string]any{"raw": b512EmergencyFlowRawValue, "state": "Unknown"}
	}
	for b512EmergencyFlowName, b512EmergencyFlowLength := range b512EmergencyFlowLengths {
		b512EmergencyFlowResult := map[string]any{"state": "incomplete"}
		b512EmergencyFlowResults[b512EmergencyFlowName] = b512EmergencyFlowResult
		b512EmergencyFlowRawValue, b512EmergencyFlowExists := b512EmergencyFlowRaw[b512EmergencyFlowName]
		if b512EmergencyFlowExists {
			b512EmergencyFlowResult["raw"] = b512EmergencyFlowRawValue
		}
		b512EmergencyFlowArray, _ := b512EmergencyFlowRawValue.([]any)
		if len(b512EmergencyFlowArray) != b512EmergencyFlowLength {
			b512EmergencyFlowOutput["state"] = "partial"
			continue
		}
		b512EmergencyFlowBytes := make([]int, len(b512EmergencyFlowArray))
		for b512EmergencyFlowIndex, b512EmergencyFlowValue := range b512EmergencyFlowArray {
			b512EmergencyFlowInteger, _ := b512EmergencyFlowValue.(json.Number).Int64()
			b512EmergencyFlowBytes[b512EmergencyFlowIndex] = int(b512EmergencyFlowInteger)
		}
		b512EmergencyFlowResult["state"] = "Hypothesis"
		b512EmergencyFlowNumber := float64(b512EmergencyFlowBytes[0])
		switch b512EmergencyFlowName {
		case "flags":
			b512EmergencyFlowLabel := "emergency"
			if b512EmergencyFlowFlow {
				b512EmergencyFlowLabel = "compressor-off"
			} else if b512EmergencyFlowDirection == "positive" {
				b512EmergencyFlowLabel = "LED"
			}
			b512EmergencyFlowFlags := map[string]any{b512EmergencyFlowLabel: b512EmergencyFlowBytes[0]&1 != 0}
			b512EmergencyFlowResult["flags"] = b512EmergencyFlowFlags
			if b512EmergencyFlowFlow && b512EmergencyFlowDirection == "positive" {
				b512EmergencyFlowFlags["pump-relay"] = b512EmergencyFlowBytes[0]&2 != 0
				b512EmergencyFlowResult["reserved"] = b512EmergencyFlowBytes[0] & 252
			} else {
				b512EmergencyFlowResult["reserved"] = b512EmergencyFlowBytes[0] & 126
				b512EmergencyFlowResult["unmodeled"] = b512EmergencyFlowBytes[0] & 128
			}
			continue
		case "temperature1", "temperature2", "temperature3":
			b512EmergencyFlowNumber = float64(int16(b512EmergencyFlowBytes[0]|b512EmergencyFlowBytes[1]<<8)) / 16
		case "pressure":
			b512EmergencyFlowNumber /= 10
		case "flow":
			b512EmergencyFlowNumber = float64(int8(b512EmergencyFlowBytes[0])) * 10
		case "co2_1", "co2_2":
			b512EmergencyFlowNumber = float64(b512EmergencyFlowBytes[0] | b512EmergencyFlowBytes[1]<<8)
		case "setpoint":
			b512EmergencyFlowResult["candidate_domain"] = b512EmergencyFlowNumber <= 100
		}
		b512EmergencyFlowResult["number"] = b512EmergencyFlowNumber
		b512EmergencyFlowResult["calculation"] = "Proven"
	}
	return json.Marshal(b512EmergencyFlowOutput)
}
