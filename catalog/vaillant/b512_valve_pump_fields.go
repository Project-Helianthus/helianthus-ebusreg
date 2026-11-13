package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

const b512ValvePumpFieldMaxDepth = 64

func b512ValvePumpFieldRead(b512ValvePumpFieldDecoder *json.Decoder) (any, error) {
	return b512ValvePumpFieldReadDepth(b512ValvePumpFieldDecoder, 0)
}

func b512ValvePumpFieldReadDepth(b512ValvePumpFieldDecoder *json.Decoder, b512ValvePumpFieldDepth int) (any, error) {
	if b512ValvePumpFieldDepth > b512ValvePumpFieldMaxDepth {
		return nil, fmt.Errorf("nesting too deep")
	}
	b512ValvePumpFieldToken, b512ValvePumpFieldErr := b512ValvePumpFieldDecoder.Token()
	if b512ValvePumpFieldErr != nil {
		return nil, b512ValvePumpFieldErr
	}
	if b512ValvePumpFieldDelimiter, b512ValvePumpFieldOK := b512ValvePumpFieldToken.(json.Delim); b512ValvePumpFieldOK {
		switch b512ValvePumpFieldDelimiter {
		case '{':
			b512ValvePumpFieldObject := map[string]any{}
			for b512ValvePumpFieldDecoder.More() {
				b512ValvePumpFieldKey, b512ValvePumpFieldErr := b512ValvePumpFieldDecoder.Token()
				if b512ValvePumpFieldErr != nil {
					return nil, b512ValvePumpFieldErr
				}
				if _, b512ValvePumpFieldExists := b512ValvePumpFieldObject[b512ValvePumpFieldKey.(string)]; b512ValvePumpFieldExists {
					return nil, fmt.Errorf("duplicate key")
				}
				b512ValvePumpFieldValue, b512ValvePumpFieldErr := b512ValvePumpFieldReadDepth(b512ValvePumpFieldDecoder, b512ValvePumpFieldDepth+1)
				if b512ValvePumpFieldErr != nil {
					return nil, b512ValvePumpFieldErr
				}
				b512ValvePumpFieldObject[b512ValvePumpFieldKey.(string)] = b512ValvePumpFieldValue
			}
			_, b512ValvePumpFieldErr = b512ValvePumpFieldDecoder.Token()
			return b512ValvePumpFieldObject, b512ValvePumpFieldErr
		case '[':
			b512ValvePumpFieldArray := []any{}
			for b512ValvePumpFieldDecoder.More() {
				b512ValvePumpFieldValue, b512ValvePumpFieldErr := b512ValvePumpFieldReadDepth(b512ValvePumpFieldDecoder, b512ValvePumpFieldDepth+1)
				if b512ValvePumpFieldErr != nil {
					return nil, b512ValvePumpFieldErr
				}
				b512ValvePumpFieldArray = append(b512ValvePumpFieldArray, b512ValvePumpFieldValue)
			}
			_, b512ValvePumpFieldErr = b512ValvePumpFieldDecoder.Token()
			return b512ValvePumpFieldArray, b512ValvePumpFieldErr
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	if b512ValvePumpFieldNumber, b512ValvePumpFieldOK := b512ValvePumpFieldToken.(json.Number); b512ValvePumpFieldOK {
		if _, b512ValvePumpFieldErr := b512ValvePumpFieldNumber.Float64(); b512ValvePumpFieldErr != nil {
			return nil, fmt.Errorf("nonfinite number")
		}
	}
	return b512ValvePumpFieldToken, nil
}

func b512ValvePumpFieldByte(b512ValvePumpFieldValue any) (int, error) {
	b512ValvePumpFieldNumber, b512ValvePumpFieldOK := b512ValvePumpFieldValue.(json.Number)
	if !b512ValvePumpFieldOK {
		return 0, fmt.Errorf("expected byte integer")
	}
	b512ValvePumpFieldInteger, b512ValvePumpFieldErr := strconv.Atoi(string(b512ValvePumpFieldNumber))
	if b512ValvePumpFieldErr != nil || b512ValvePumpFieldInteger < 0 || b512ValvePumpFieldInteger > 255 {
		return 0, fmt.Errorf("invalid byte integer")
	}
	return b512ValvePumpFieldInteger, nil
}

func b512ValvePumpFieldBytes(b512ValvePumpFieldValue any) ([]int, error) {
	b512ValvePumpFieldArray, b512ValvePumpFieldOK := b512ValvePumpFieldValue.([]any)
	if !b512ValvePumpFieldOK {
		return nil, fmt.Errorf("expected byte array")
	}
	b512ValvePumpFieldBytes := make([]int, len(b512ValvePumpFieldArray))
	for b512ValvePumpFieldIndex, b512ValvePumpFieldItem := range b512ValvePumpFieldArray {
		b512ValvePumpFieldByteValue, b512ValvePumpFieldErr := b512ValvePumpFieldByte(b512ValvePumpFieldItem)
		if b512ValvePumpFieldErr != nil {
			return nil, b512ValvePumpFieldErr
		}
		b512ValvePumpFieldBytes[b512ValvePumpFieldIndex] = b512ValvePumpFieldByteValue
	}
	return b512ValvePumpFieldBytes, nil
}

func B512ValvePumpFieldEvaluate(b512ValvePumpFieldInput json.RawMessage) (json.RawMessage, error) {
	b512ValvePumpFieldDecoder := json.NewDecoder(bytes.NewReader(b512ValvePumpFieldInput))
	b512ValvePumpFieldDecoder.UseNumber()
	b512ValvePumpFieldValue, b512ValvePumpFieldErr := b512ValvePumpFieldRead(b512ValvePumpFieldDecoder)
	if b512ValvePumpFieldErr != nil {
		return nil, b512ValvePumpFieldErr
	}
	if _, b512ValvePumpFieldErr = b512ValvePumpFieldDecoder.Token(); b512ValvePumpFieldErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b512ValvePumpFieldRoot, b512ValvePumpFieldOK := b512ValvePumpFieldValue.(map[string]any)
	if !b512ValvePumpFieldOK {
		return nil, fmt.Errorf("expected object")
	}
	b512ValvePumpFieldContext, b512ValvePumpFieldOK := b512ValvePumpFieldRoot["context"].(map[string]any)
	if !b512ValvePumpFieldOK {
		return nil, fmt.Errorf("expected context object")
	}
	b512ValvePumpFieldFields, b512ValvePumpFieldOK := b512ValvePumpFieldRoot["fields"].(map[string]any)
	if !b512ValvePumpFieldOK {
		return nil, fmt.Errorf("expected fields object")
	}
	b512ValvePumpFieldNumbers := map[string]int{}
	for _, b512ValvePumpFieldKey := range []string{"primary", "secondary", "operation"} {
		if b512ValvePumpFieldItem, b512ValvePumpFieldExists := b512ValvePumpFieldContext[b512ValvePumpFieldKey]; b512ValvePumpFieldExists {
			b512ValvePumpFieldNumber, b512ValvePumpFieldErr := b512ValvePumpFieldByte(b512ValvePumpFieldItem)
			if b512ValvePumpFieldErr != nil {
				return nil, b512ValvePumpFieldErr
			}
			b512ValvePumpFieldNumbers[b512ValvePumpFieldKey] = b512ValvePumpFieldNumber
		}
	}
	for _, b512ValvePumpFieldKey := range []string{"direction", "view"} {
		if b512ValvePumpFieldItem, b512ValvePumpFieldExists := b512ValvePumpFieldContext[b512ValvePumpFieldKey]; b512ValvePumpFieldExists {
			if _, b512ValvePumpFieldOK := b512ValvePumpFieldItem.(string); !b512ValvePumpFieldOK {
				return nil, fmt.Errorf("expected context string")
			}
		}
	}
	b512ValvePumpFieldRaw := map[string][]int{}
	for b512ValvePumpFieldKey, b512ValvePumpFieldItem := range b512ValvePumpFieldFields {
		b512ValvePumpFieldBytesValue, b512ValvePumpFieldErr := b512ValvePumpFieldBytes(b512ValvePumpFieldItem)
		if b512ValvePumpFieldErr != nil {
			return nil, b512ValvePumpFieldErr
		}
		b512ValvePumpFieldRaw[b512ValvePumpFieldKey] = b512ValvePumpFieldBytesValue
	}
	b512ValvePumpFieldArguments := map[string]any{}
	if b512ValvePumpFieldItem, b512ValvePumpFieldExists := b512ValvePumpFieldRoot["arguments"]; b512ValvePumpFieldExists {
		b512ValvePumpFieldArguments, b512ValvePumpFieldOK = b512ValvePumpFieldItem.(map[string]any)
		if !b512ValvePumpFieldOK {
			return nil, fmt.Errorf("expected arguments object")
		}
	}
	b512ValvePumpFieldOperation, b512ValvePumpFieldHasOperation := b512ValvePumpFieldNumbers["operation"]
	b512ValvePumpFieldDirection := b512ValvePumpFieldContext["direction"]
	b512ValvePumpFieldSupported := b512ValvePumpFieldHasOperation && b512ValvePumpFieldNumbers["primary"] == 181 && b512ValvePumpFieldNumbers["secondary"] == 18 && b512ValvePumpFieldContext["view"] == "actuator-fields" && (b512ValvePumpFieldDirection == "request" || b512ValvePumpFieldDirection == "positive") && (b512ValvePumpFieldOperation == 0 || b512ValvePumpFieldOperation == 2 || b512ValvePumpFieldOperation == 3 || b512ValvePumpFieldOperation == 4 || b512ValvePumpFieldOperation == 5)
	if b512ValvePumpFieldDirection == "positive" && (b512ValvePumpFieldOperation == 0 || b512ValvePumpFieldOperation == 5) {
		b512ValvePumpFieldSupported = false
	}
	b512ValvePumpFieldModes := []int{}
	if b512ValvePumpFieldSupported {
		if b512ValvePumpFieldOperation == 5 && b512ValvePumpFieldDirection == "request" {
			b512ValvePumpFieldModes, b512ValvePumpFieldErr = b512ValvePumpFieldBytes(b512ValvePumpFieldArguments["external_modes"])
			if b512ValvePumpFieldErr != nil || len(b512ValvePumpFieldModes) != 3 || len(b512ValvePumpFieldArguments) != 1 {
				return nil, fmt.Errorf("expected exactly three external modes")
			}
		} else if len(b512ValvePumpFieldArguments) != 0 {
			return nil, fmt.Errorf("unexpected arguments")
		}
	}
	b512ValvePumpFieldOriginal := b512ValvePumpFieldRoot
	b512ValvePumpFieldRoot = map[string]any{}
	for b512ValvePumpFieldKey, b512ValvePumpFieldItem := range b512ValvePumpFieldOriginal {
		b512ValvePumpFieldRoot[b512ValvePumpFieldKey] = b512ValvePumpFieldItem
	}
	b512ValvePumpFieldRoot["input"] = b512ValvePumpFieldOriginal
	b512ValvePumpFieldResults := map[string]any{}
	for b512ValvePumpFieldKey, b512ValvePumpFieldItem := range b512ValvePumpFieldFields {
		b512ValvePumpFieldResults[b512ValvePumpFieldKey] = map[string]any{"raw": b512ValvePumpFieldItem, "state": "Unknown"}
	}
	b512ValvePumpFieldRoot["physical"] = "Unknown"
	b512ValvePumpFieldRoot["qualified"] = false
	b512ValvePumpFieldRoot["state"] = "Unknown"
	b512ValvePumpFieldRoot["results"] = b512ValvePumpFieldResults
	if !b512ValvePumpFieldSupported {
		return json.Marshal(b512ValvePumpFieldRoot)
	}
	b512ValvePumpFieldWidths := map[string]int{}
	if b512ValvePumpFieldDirection == "request" {
		switch b512ValvePumpFieldOperation {
		case 0, 4:
			b512ValvePumpFieldWidths["setpoint"] = 1
		case 2:
			b512ValvePumpFieldWidths["valve"] = 1
		case 3:
			b512ValvePumpFieldWidths["selector"] = 1
		case 5:
			b512ValvePumpFieldWidths = map[string]int{"channel0": 1, "channel1": 1, "channel2": 1}
		}
	} else {
		switch b512ValvePumpFieldOperation {
		case 2:
			b512ValvePumpFieldWidths = map[string]int{"valve": 1, "contact": 1, "tank": 2, "virtual_error": 2}
		case 3:
			b512ValvePumpFieldWidths = map[string]int{"status": 1, "temperature_raw": 2}
		case 4:
			b512ValvePumpFieldWidths["feedback"] = 1
		}
	}
	b512ValvePumpFieldRoot["state"] = "complete"
	for b512ValvePumpFieldKey, b512ValvePumpFieldWidth := range b512ValvePumpFieldWidths {
		b512ValvePumpFieldResult := map[string]any{"state": "incomplete", "meaning": "Hypothesis"}
		if b512ValvePumpFieldItem, b512ValvePumpFieldExists := b512ValvePumpFieldFields[b512ValvePumpFieldKey]; b512ValvePumpFieldExists {
			b512ValvePumpFieldResult["raw"] = b512ValvePumpFieldItem
		}
		b512ValvePumpFieldResults[b512ValvePumpFieldKey] = b512ValvePumpFieldResult
		b512ValvePumpFieldData := b512ValvePumpFieldRaw[b512ValvePumpFieldKey]
		if len(b512ValvePumpFieldData) != b512ValvePumpFieldWidth {
			b512ValvePumpFieldRoot["state"] = "partial"
			continue
		}
		b512ValvePumpFieldNumber := b512ValvePumpFieldData[0]
		if b512ValvePumpFieldWidth == 2 {
			b512ValvePumpFieldNumber += 256 * b512ValvePumpFieldData[1]
		}
		b512ValvePumpFieldResult["state"] = "complete"
		b512ValvePumpFieldResult["number"] = b512ValvePumpFieldNumber
		b512ValvePumpFieldResult["calculation"] = "Proven"
		b512ValvePumpFieldLabel := "unmapped"
		switch b512ValvePumpFieldKey {
		case "setpoint", "channel0", "channel1", "channel2":
			if b512ValvePumpFieldNumber == 0 {
				b512ValvePumpFieldLabel = "off"
			} else if b512ValvePumpFieldNumber <= 100 {
				b512ValvePumpFieldLabel = "on-or-speed"
			}
			if b512ValvePumpFieldOperation == 4 && b512ValvePumpFieldNumber == 255 {
				b512ValvePumpFieldLabel = "no-remote-control"
				b512ValvePumpFieldResult["alternatives"] = []string{"replacement", "no-remote-control"}
			}
			if b512ValvePumpFieldOperation == 5 {
				if b512ValvePumpFieldNumber == 255 {
					b512ValvePumpFieldLabel = "unchanged"
				}
				b512ValvePumpFieldResult["prerequisite"] = b512ValvePumpFieldModes[int(b512ValvePumpFieldKey[len(b512ValvePumpFieldKey)-1]-'0')] == 8
			}
			if b512ValvePumpFieldOperation == 4 {
				b512ValvePumpFieldResult["candidate_domain"] = map[string]any{"timeout_seconds": 900, "reset": "every-service"}
			}
		case "valve":
			if b512ValvePumpFieldDirection == "positive" {
				delete(b512ValvePumpFieldResult, "number")
				delete(b512ValvePumpFieldResult, "calculation")
				b512ValvePumpFieldResult["alternatives"] = []string{"actual", "setpoint-fallback"}
				b512ValvePumpFieldLabel = "unresolved"
			} else {
				switch b512ValvePumpFieldNumber {
				case 0:
					b512ValvePumpFieldLabel = "CH"
				case 50:
					b512ValvePumpFieldLabel = "midpoint"
				case 100:
					b512ValvePumpFieldLabel = "DHW"
				case 254:
					b512ValvePumpFieldLabel = "standby"
				case 255:
					b512ValvePumpFieldLabel = "unchanged"
				}
			}
		case "contact":
			if b512ValvePumpFieldNumber == 0 {
				b512ValvePumpFieldLabel = "open"
			}
			if b512ValvePumpFieldNumber == 1 {
				b512ValvePumpFieldLabel = "closed"
			}
		case "tank":
			b512ValvePumpFieldResult["number"] = float64(int16(b512ValvePumpFieldNumber)) / 16
			b512ValvePumpFieldLabel = "Celsius"
		case "virtual_error":
			if b512ValvePumpFieldNumber == 65535 {
				b512ValvePumpFieldLabel = "no-error"
			}
		case "selector":
			if b512ValvePumpFieldNumber == 0 {
				b512ValvePumpFieldLabel = "valve-temperature-input"
			}
		case "status":
			switch b512ValvePumpFieldNumber {
			case 0:
				b512ValvePumpFieldLabel = "okay"
			case 1:
				b512ValvePumpFieldLabel = "disconnected"
			case 2:
				b512ValvePumpFieldLabel = "short-circuit"
			}
		case "temperature_raw":
			b512ValvePumpFieldLabel = "raw"
		case "feedback":
			b512ValvePumpFieldLabel = "rejected"
			if b512ValvePumpFieldNumber == 1 {
				b512ValvePumpFieldLabel = "accepted"
			}
		}
		b512ValvePumpFieldResult["label"] = b512ValvePumpFieldLabel
	}
	return json.Marshal(b512ValvePumpFieldRoot)
}
