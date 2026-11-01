package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func ab15ap15Read(ab15ap15Decoder *json.Decoder) (any, error) {
	ab15ap15Token, ab15ap15Err := ab15ap15Decoder.Token()
	if ab15ap15Err != nil {
		return nil, ab15ap15Err
	}
	if ab15ap15Delimiter, ab15ap15OK := ab15ap15Token.(json.Delim); ab15ap15OK {
		switch ab15ap15Delimiter {
		case '{':
			ab15ap15Object := map[string]any{}
			for ab15ap15Decoder.More() {
				ab15ap15KeyToken, ab15ap15Err := ab15ap15Decoder.Token()
				if ab15ap15Err != nil {
					return nil, ab15ap15Err
				}
				ab15ap15Key := ab15ap15KeyToken.(string)
				if _, ab15ap15Exists := ab15ap15Object[ab15ap15Key]; ab15ap15Exists {
					return nil, fmt.Errorf("duplicate key %q", ab15ap15Key)
				}
				ab15ap15Value, ab15ap15Err := ab15ap15Read(ab15ap15Decoder)
				if ab15ap15Err != nil {
					return nil, ab15ap15Err
				}
				ab15ap15Object[ab15ap15Key] = ab15ap15Value
			}
			_, ab15ap15Err := ab15ap15Decoder.Token()
			return ab15ap15Object, ab15ap15Err
		case '[':
			ab15ap15Array := []any{}
			for ab15ap15Decoder.More() {
				ab15ap15Value, ab15ap15Err := ab15ap15Read(ab15ap15Decoder)
				if ab15ap15Err != nil {
					return nil, ab15ap15Err
				}
				ab15ap15Array = append(ab15ap15Array, ab15ap15Value)
			}
			_, ab15ap15Err := ab15ap15Decoder.Token()
			return ab15ap15Array, ab15ap15Err
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	return ab15ap15Token, nil
}

func Ab15Ap15Evaluate(ab15ap15Input json.RawMessage) (json.RawMessage, error) {
	ab15ap15Decoder := json.NewDecoder(bytes.NewReader(ab15ap15Input))
	ab15ap15Decoder.UseNumber()
	ab15ap15Value, ab15ap15Err := ab15ap15Read(ab15ap15Decoder)
	if ab15ap15Err != nil {
		return nil, ab15ap15Err
	}
	if _, ab15ap15Err = ab15ap15Decoder.Token(); ab15ap15Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab15ap15Root, ab15ap15OK := ab15ap15Value.(map[string]any)
	if !ab15ap15OK {
		return nil, fmt.Errorf("input must be object")
	}
	ab15ap15Context, ab15ap15OK := ab15ap15Root["context"].(map[string]any)
	if !ab15ap15OK {
		return nil, fmt.Errorf("context must be object")
	}
	for _, ab15ap15Key := range []string{"primary", "secondary", "operation"} {
		if ab15ap15Value, ab15ap15Exists := ab15ap15Context[ab15ap15Key]; ab15ap15Exists && ab15ap15Value != nil {
			ab15ap15Number, ab15ap15OK := ab15ap15Value.(json.Number)
			if !ab15ap15OK {
				return nil, fmt.Errorf("context number must be integer")
			}
			ab15ap15Integer, ab15ap15Err := ab15ap15Number.Int64()
			if ab15ap15Err != nil || ab15ap15Integer < 0 || ab15ap15Integer > 255 {
				return nil, fmt.Errorf("context number out of range")
			}
		}
	}
	for _, ab15ap15Key := range []string{"direction", "view"} {
		if ab15ap15Value, ab15ap15Exists := ab15ap15Context[ab15ap15Key]; ab15ap15Exists && ab15ap15Value != nil {
			if _, ab15ap15OK := ab15ap15Value.(string); !ab15ap15OK {
				return nil, fmt.Errorf("context selector must be string")
			}
		}
	}
	ab15ap15Raw, ab15ap15OK := ab15ap15Root["fields"].(map[string]any)
	if !ab15ap15OK {
		return nil, fmt.Errorf("fields must be object")
	}
	if ab15ap15Arguments, ab15ap15Exists := ab15ap15Root["arguments"]; ab15ap15Exists {
		if _, ab15ap15OK := ab15ap15Arguments.(map[string]any); !ab15ap15OK {
			return nil, fmt.Errorf("arguments must be object")
		}
	}
	for _, ab15ap15Value := range ab15ap15Raw {
		ab15ap15Array, ab15ap15OK := ab15ap15Value.([]any)
		if !ab15ap15OK {
			return nil, fmt.Errorf("field must be byte array")
		}
		for _, ab15ap15Byte := range ab15ap15Array {
			ab15ap15Number, ab15ap15OK := ab15ap15Byte.(json.Number)
			if !ab15ap15OK {
				return nil, fmt.Errorf("byte must be integer")
			}
			ab15ap15Integer, ab15ap15Err := ab15ap15Number.Int64()
			if ab15ap15Err != nil || ab15ap15Integer < 0 || ab15ap15Integer > 255 {
				return nil, fmt.Errorf("byte out of range")
			}
		}
	}
	ab15ap15Results := map[string]any{}
	ab15ap15Output := map[string]any{"context": ab15ap15Context, "raw": ab15ap15Raw, "state": "Unknown", "fields": ab15ap15Results, "physical": "Unknown", "qualified": false}
	if ab15ap15Arguments, ab15ap15Exists := ab15ap15Root["arguments"]; ab15ap15Exists {
		ab15ap15Output["arguments"] = ab15ap15Arguments
	}
	ab15ap15Operation := ab15ap15Context["operation"]
	ab15ap15Direction := ab15ap15Context["direction"]
	ab15ap15Flow := ab15ap15Operation == json.Number("17") && ab15ap15Context["view"] == "flow-io"
	ab15ap15Emergency := ab15ap15Operation == json.Number("18") && ab15ap15Context["view"] == "emergency-io"
	if ab15ap15Context["primary"] != json.Number("181") || ab15ap15Context["secondary"] != json.Number("18") || (!ab15ap15Flow && !ab15ap15Emergency) || (ab15ap15Direction != "request" && ab15ap15Direction != "positive") {
		return json.Marshal(ab15ap15Output)
	}
	ab15ap15Output["state"] = "complete"
	ab15ap15Lengths := map[string]int{"flags": 1}
	if ab15ap15Direction == "positive" {
		ab15ap15Output["candidate"] = map[string]any{"state": "Hypothesis", "declared_positive_prefix": []int{181, 18, 12}, "requested_operation": ab15ap15Operation, "literal_reply_conflict": true, "physical_association": "Unknown"}
		if ab15ap15Flow {
			ab15ap15Lengths["setpoint"] = 1
		}
	} else if ab15ap15Flow {
		for _, ab15ap15Name := range []string{"temperature1", "temperature2", "temperature3"} {
			ab15ap15Lengths[ab15ap15Name] = 2
		}
		for _, ab15ap15Name := range []string{"pressure", "flow", "pump_status"} {
			ab15ap15Lengths[ab15ap15Name] = 1
		}
	} else {
		ab15ap15Lengths["switch"] = 1
		ab15ap15Lengths["co2_1"] = 2
		ab15ap15Lengths["co2_2"] = 2
	}
	for ab15ap15Name, ab15ap15RawValue := range ab15ap15Raw {
		ab15ap15Results[ab15ap15Name] = map[string]any{"raw": ab15ap15RawValue, "state": "Unknown"}
	}
	for ab15ap15Name, ab15ap15Length := range ab15ap15Lengths {
		ab15ap15Result := map[string]any{"state": "incomplete"}
		ab15ap15Results[ab15ap15Name] = ab15ap15Result
		ab15ap15RawValue, ab15ap15Exists := ab15ap15Raw[ab15ap15Name]
		if ab15ap15Exists {
			ab15ap15Result["raw"] = ab15ap15RawValue
		}
		ab15ap15Array, _ := ab15ap15RawValue.([]any)
		if len(ab15ap15Array) != ab15ap15Length {
			ab15ap15Output["state"] = "partial"
			continue
		}
		ab15ap15Bytes := make([]int, len(ab15ap15Array))
		for ab15ap15Index, ab15ap15Value := range ab15ap15Array {
			ab15ap15Integer, _ := ab15ap15Value.(json.Number).Int64()
			ab15ap15Bytes[ab15ap15Index] = int(ab15ap15Integer)
		}
		ab15ap15Result["state"] = "Hypothesis"
		ab15ap15Number := float64(ab15ap15Bytes[0])
		switch ab15ap15Name {
		case "flags":
			ab15ap15Label := "emergency"
			if ab15ap15Flow {
				ab15ap15Label = "compressor-off"
			} else if ab15ap15Direction == "positive" {
				ab15ap15Label = "LED"
			}
			ab15ap15Flags := map[string]any{ab15ap15Label: ab15ap15Bytes[0]&1 != 0}
			ab15ap15Result["flags"] = ab15ap15Flags
			if ab15ap15Flow && ab15ap15Direction == "positive" {
				ab15ap15Flags["pump-relay"] = ab15ap15Bytes[0]&2 != 0
				ab15ap15Result["reserved"] = ab15ap15Bytes[0] & 252
			} else {
				ab15ap15Result["reserved"] = ab15ap15Bytes[0] & 126
				ab15ap15Result["unmodeled"] = ab15ap15Bytes[0] & 128
			}
			continue
		case "temperature1", "temperature2", "temperature3":
			ab15ap15Number = float64(int16(ab15ap15Bytes[0]|ab15ap15Bytes[1]<<8)) / 16
		case "pressure":
			ab15ap15Number /= 10
		case "flow":
			ab15ap15Number = float64(int8(ab15ap15Bytes[0])) * 10
		case "co2_1", "co2_2":
			ab15ap15Number = float64(ab15ap15Bytes[0] | ab15ap15Bytes[1]<<8)
		case "setpoint":
			ab15ap15Result["candidate_domain"] = ab15ap15Number <= 100
		}
		ab15ap15Result["number"] = ab15ap15Number
		ab15ap15Result["calculation"] = "Proven"
	}
	return json.Marshal(ab15ap15Output)
}
