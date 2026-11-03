package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func ab25aq25Read(ab25aq25Decoder *json.Decoder) (any, error) {
	ab25aq25Token, ab25aq25Err := ab25aq25Decoder.Token()
	if ab25aq25Err != nil {
		return nil, ab25aq25Err
	}
	if ab25aq25Delimiter, ab25aq25OK := ab25aq25Token.(json.Delim); ab25aq25OK {
		switch ab25aq25Delimiter {
		case '{':
			ab25aq25Object := map[string]any{}
			for ab25aq25Decoder.More() {
				ab25aq25KeyToken, ab25aq25Err := ab25aq25Decoder.Token()
				if ab25aq25Err != nil {
					return nil, ab25aq25Err
				}
				ab25aq25Key := ab25aq25KeyToken.(string)
				if _, ab25aq25Exists := ab25aq25Object[ab25aq25Key]; ab25aq25Exists {
					return nil, fmt.Errorf("duplicate key %q", ab25aq25Key)
				}
				ab25aq25Value, ab25aq25Err := ab25aq25Read(ab25aq25Decoder)
				if ab25aq25Err != nil {
					return nil, ab25aq25Err
				}
				ab25aq25Object[ab25aq25Key] = ab25aq25Value
			}
			_, ab25aq25Err := ab25aq25Decoder.Token()
			return ab25aq25Object, ab25aq25Err
		case '[':
			ab25aq25Array := []any{}
			for ab25aq25Decoder.More() {
				ab25aq25Value, ab25aq25Err := ab25aq25Read(ab25aq25Decoder)
				if ab25aq25Err != nil {
					return nil, ab25aq25Err
				}
				ab25aq25Array = append(ab25aq25Array, ab25aq25Value)
			}
			_, ab25aq25Err := ab25aq25Decoder.Token()
			return ab25aq25Array, ab25aq25Err
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	if ab25aq25Number, ab25aq25OK := ab25aq25Token.(json.Number); ab25aq25OK {
		if _, ab25aq25Err := ab25aq25Number.Float64(); ab25aq25Err != nil {
			return nil, fmt.Errorf("nonfinite number")
		}
	}
	return ab25aq25Token, nil
}

func Ab25Aq25Evaluate(ab25aq25Input json.RawMessage) (json.RawMessage, error) {
	ab25aq25Decoder := json.NewDecoder(bytes.NewReader(ab25aq25Input))
	ab25aq25Decoder.UseNumber()
	ab25aq25Value, ab25aq25Err := ab25aq25Read(ab25aq25Decoder)
	if ab25aq25Err != nil {
		return nil, ab25aq25Err
	}
	if _, ab25aq25Err = ab25aq25Decoder.Token(); ab25aq25Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab25aq25Root, ab25aq25OK := ab25aq25Value.(map[string]any)
	if !ab25aq25OK {
		return nil, fmt.Errorf("input must be object")
	}
	ab25aq25Context, ab25aq25OK := ab25aq25Root["context"].(map[string]any)
	if !ab25aq25OK {
		return nil, fmt.Errorf("context must be object")
	}
	for _, ab25aq25Key := range []string{"primary", "secondary", "operation"} {
		if ab25aq25Value, ab25aq25Exists := ab25aq25Context[ab25aq25Key]; ab25aq25Exists {
			ab25aq25Number, ab25aq25OK := ab25aq25Value.(json.Number)
			if !ab25aq25OK {
				return nil, fmt.Errorf("context number must be integer")
			}
			ab25aq25Integer, ab25aq25Err := ab25aq25Number.Int64()
			if ab25aq25Err != nil || ab25aq25Integer < 0 || ab25aq25Integer > 255 {
				return nil, fmt.Errorf("context number out of range")
			}
		}
	}
	for _, ab25aq25Key := range []string{"direction", "view"} {
		if ab25aq25Value, ab25aq25Exists := ab25aq25Context[ab25aq25Key]; ab25aq25Exists {
			if _, ab25aq25OK := ab25aq25Value.(string); !ab25aq25OK {
				return nil, fmt.Errorf("context selector must be string")
			}
		}
	}
	ab25aq25Raw, ab25aq25OK := ab25aq25Root["fields"].(map[string]any)
	if !ab25aq25OK {
		return nil, fmt.Errorf("fields must be object")
	}
	if ab25aq25Arguments, ab25aq25Exists := ab25aq25Root["arguments"]; ab25aq25Exists {
		if _, ab25aq25OK := ab25aq25Arguments.(map[string]any); !ab25aq25OK {
			return nil, fmt.Errorf("arguments must be object")
		}
	}
	for _, ab25aq25Value := range ab25aq25Raw {
		ab25aq25Array, ab25aq25OK := ab25aq25Value.([]any)
		if !ab25aq25OK {
			return nil, fmt.Errorf("field must be byte array")
		}
		for _, ab25aq25Byte := range ab25aq25Array {
			ab25aq25Number, ab25aq25OK := ab25aq25Byte.(json.Number)
			if !ab25aq25OK {
				return nil, fmt.Errorf("byte must be integer")
			}
			ab25aq25Integer, ab25aq25Err := ab25aq25Number.Int64()
			if ab25aq25Err != nil || ab25aq25Integer < 0 || ab25aq25Integer > 255 {
				return nil, fmt.Errorf("byte out of range")
			}
		}
	}
	ab25aq25Results := map[string]any{}
	for ab25aq25Name, ab25aq25Value := range ab25aq25Raw {
		ab25aq25Results[ab25aq25Name] = map[string]any{"raw": ab25aq25Value, "state": "Unknown"}
	}
	ab25aq25Root["results"] = ab25aq25Results
	ab25aq25Root["state"] = "Unknown"
	ab25aq25Root["physical"] = "Unknown"
	ab25aq25Root["qualified"] = false
	ab25aq25Operation, _ := ab25aq25Context["operation"].(json.Number)
	ab25aq25Op, _ := ab25aq25Operation.Int64()
	ab25aq25Primary, _ := ab25aq25Context["primary"].(json.Number)
	ab25aq25Secondary, _ := ab25aq25Context["secondary"].(json.Number)
	ab25aq25P, _ := ab25aq25Primary.Int64()
	ab25aq25S, _ := ab25aq25Secondary.Int64()
	ab25aq25Direction := ab25aq25Context["direction"]
	if ab25aq25P != 181 || ab25aq25S != 18 || ab25aq25Op < 7 || ab25aq25Op > 9 || ab25aq25Context["view"] != "power-stepper" || (ab25aq25Direction != "request" && ab25aq25Direction != "positive") {
		return json.Marshal(ab25aq25Root)
	}
	ab25aq25Lengths := map[string]int{}
	switch ab25aq25Op {
	case 7:
		if ab25aq25Direction == "positive" {
			ab25aq25Lengths["maximum_power"] = 1
			ab25aq25Lengths["minimum_power"] = 1
		}
	case 8:
		ab25aq25Lengths["stepper"] = 1
		if ab25aq25Direction == "request" {
			ab25aq25Lengths["inlet"] = 2
			ab25aq25Lengths["status"] = 1
		} else {
			ab25aq25Lengths["setpoint"] = 2
		}
	case 9:
		ab25aq25Lengths["release"] = 1
	}
	ab25aq25Root["state"] = "complete"
	ab25aq25Root["grade"] = "Hypothesis"
	for ab25aq25Name, ab25aq25Width := range ab25aq25Lengths {
		ab25aq25Result := map[string]any{"state": "incomplete"}
		ab25aq25Results[ab25aq25Name] = ab25aq25Result
		ab25aq25Value, ab25aq25Present := ab25aq25Raw[ab25aq25Name]
		if ab25aq25Present {
			ab25aq25Result["raw"] = ab25aq25Value
		}
		ab25aq25Array, _ := ab25aq25Value.([]any)
		if len(ab25aq25Array) != ab25aq25Width {
			ab25aq25Root["state"] = "partial"
			continue
		}
		ab25aq25Bytes := make([]int, ab25aq25Width)
		for ab25aq25Index, ab25aq25Byte := range ab25aq25Array {
			ab25aq25Integer, _ := ab25aq25Byte.(json.Number).Int64()
			ab25aq25Bytes[ab25aq25Index] = int(ab25aq25Integer)
		}
		ab25aq25Result["state"] = "complete"
		ab25aq25Result["grade"] = "Hypothesis"
		switch ab25aq25Name {
		case "maximum_power":
			ab25aq25Result["number"] = ab25aq25Bytes[0]
			ab25aq25Result["unit"] = "kW"
			ab25aq25Result["default"] = 0
			ab25aq25Result["calculation"] = "Proven"
		case "minimum_power":
			ab25aq25Result["number"] = float64(ab25aq25Bytes[0]) / 2
			ab25aq25Result["unit"] = "%"
			ab25aq25Result["default"] = 0
			ab25aq25Result["candidate_domain"] = ab25aq25Bytes[0] <= 250
			ab25aq25Result["calculation"] = "Proven"
		case "inlet", "setpoint":
			ab25aq25Result["number"] = float64(int16(ab25aq25Bytes[0]|ab25aq25Bytes[1]<<8)) / 16
			ab25aq25Result["unit"] = "Celsius"
			ab25aq25Result["calculation"] = "Proven"
		case "status":
			ab25aq25Label := "unmapped"
			switch ab25aq25Bytes[0] {
			case 0:
				ab25aq25Label = "valid"
			case 0x55:
				ab25aq25Label = "short"
			case 0xaa:
				ab25aq25Label = "interruption-or-disconnected"
			}
			ab25aq25Result["label"] = ab25aq25Label
		case "stepper", "release":
			ab25aq25Result["flags"] = map[string]any{"release": ab25aq25Bytes[0]&1 != 0}
			ab25aq25Result["reserved"] = ab25aq25Bytes[0] & 254
		}
	}
	return json.Marshal(ab25aq25Root)
}
