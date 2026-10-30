package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
)

func ab25ap25Read(ab25ap25Decoder *json.Decoder) (any, error) {
	ab25ap25Token, ab25ap25Error := ab25ap25Decoder.Token()
	if ab25ap25Error != nil {
		return nil, ab25ap25Error
	}
	if ab25ap25Delimiter, ab25ap25OK := ab25ap25Token.(json.Delim); ab25ap25OK {
		switch ab25ap25Delimiter {
		case '{':
			ab25ap25Object := map[string]any{}
			for ab25ap25Decoder.More() {
				ab25ap25KeyToken, ab25ap25Error := ab25ap25Decoder.Token()
				if ab25ap25Error != nil {
					return nil, ab25ap25Error
				}
				ab25ap25Key := ab25ap25KeyToken.(string)
				if _, ab25ap25Exists := ab25ap25Object[ab25ap25Key]; ab25ap25Exists {
					return nil, fmt.Errorf("duplicate key %q", ab25ap25Key)
				}
				ab25ap25Value, ab25ap25Error := ab25ap25Read(ab25ap25Decoder)
				if ab25ap25Error != nil {
					return nil, ab25ap25Error
				}
				ab25ap25Object[ab25ap25Key] = ab25ap25Value
			}
			if _, ab25ap25Error := ab25ap25Decoder.Token(); ab25ap25Error != nil {
				return nil, ab25ap25Error
			}
			return ab25ap25Object, nil
		case '[':
			ab25ap25Array := []any{}
			for ab25ap25Decoder.More() {
				ab25ap25Value, ab25ap25Error := ab25ap25Read(ab25ap25Decoder)
				if ab25ap25Error != nil {
					return nil, ab25ap25Error
				}
				ab25ap25Array = append(ab25ap25Array, ab25ap25Value)
			}
			if _, ab25ap25Error := ab25ap25Decoder.Token(); ab25ap25Error != nil {
				return nil, ab25ap25Error
			}
			return ab25ap25Array, nil
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	return ab25ap25Token, nil
}

func ab25ap25Byte(ab25ap25Value any) (int, bool) {
	ab25ap25Number, ab25ap25OK := ab25ap25Value.(json.Number)
	if !ab25ap25OK {
		return 0, false
	}
	ab25ap25Rational, ab25ap25OK := new(big.Rat).SetString(string(ab25ap25Number))
	if !ab25ap25OK || !ab25ap25Rational.IsInt() || !ab25ap25Rational.Num().IsInt64() {
		return 0, false
	}
	ab25ap25Integer := ab25ap25Rational.Num().Int64()
	if ab25ap25Integer < 0 || ab25ap25Integer > 255 {
		return 0, false
	}
	return int(ab25ap25Integer), true
}

func Ab25Ap25Evaluate(ab25ap25Input json.RawMessage) (json.RawMessage, error) {
	ab25ap25Decoder := json.NewDecoder(bytes.NewReader(ab25ap25Input))
	ab25ap25Decoder.UseNumber()
	ab25ap25Value, ab25ap25Error := ab25ap25Read(ab25ap25Decoder)
	if ab25ap25Error != nil {
		return nil, ab25ap25Error
	}
	if _, ab25ap25Error = ab25ap25Decoder.Token(); ab25ap25Error != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab25ap25InputObject, ab25ap25OK := ab25ap25Value.(map[string]any)
	if !ab25ap25OK {
		return nil, fmt.Errorf("input must be object")
	}
	ab25ap25Context, ab25ap25OK := ab25ap25InputObject["context"].(map[string]any)
	if !ab25ap25OK {
		return nil, fmt.Errorf("context must be object")
	}
	ab25ap25Raw, ab25ap25OK := ab25ap25InputObject["fields"].(map[string]any)
	if !ab25ap25OK {
		return nil, fmt.Errorf("fields must be object")
	}
	for _, ab25ap25Name := range []string{"primary", "secondary", "operation"} {
		if ab25ap25Value, ab25ap25Exists := ab25ap25Context[ab25ap25Name]; ab25ap25Exists && ab25ap25Value != nil {
			if _, ab25ap25OK := ab25ap25Byte(ab25ap25Value); !ab25ap25OK {
				return nil, fmt.Errorf("invalid context %s", ab25ap25Name)
			}
		}
	}
	for _, ab25ap25Name := range []string{"direction", "view"} {
		if ab25ap25Value, ab25ap25Exists := ab25ap25Context[ab25ap25Name]; ab25ap25Exists && ab25ap25Value != nil {
			if _, ab25ap25OK := ab25ap25Value.(string); !ab25ap25OK {
				return nil, fmt.Errorf("invalid context %s", ab25ap25Name)
			}
		}
	}
	for ab25ap25Name, ab25ap25Value := range ab25ap25Raw {
		ab25ap25Array, ab25ap25OK := ab25ap25Value.([]any)
		if !ab25ap25OK {
			return nil, fmt.Errorf("field %s must be byte array", ab25ap25Name)
		}
		for _, ab25ap25Item := range ab25ap25Array {
			if _, ab25ap25OK := ab25ap25Byte(ab25ap25Item); !ab25ap25OK {
				return nil, fmt.Errorf("invalid byte in %s", ab25ap25Name)
			}
		}
	}
	ab25ap25Fields := map[string]any{}
	ab25ap25Output := map[string]any{"context": ab25ap25Context, "raw": ab25ap25Raw, "state": "Unknown", "fields": ab25ap25Fields, "physical": "Unknown", "qualified": false}
	if ab25ap25Arguments, ab25ap25Exists := ab25ap25InputObject["arguments"]; ab25ap25Exists {
		if _, ab25ap25OK := ab25ap25Arguments.(map[string]any); !ab25ap25OK {
			return nil, fmt.Errorf("arguments must be object")
		}
		ab25ap25Output["arguments"] = ab25ap25Arguments
	}
	ab25ap25Primary, _ := ab25ap25Byte(ab25ap25Context["primary"])
	ab25ap25Secondary, _ := ab25ap25Byte(ab25ap25Context["secondary"])
	ab25ap25Operation, _ := ab25ap25Byte(ab25ap25Context["operation"])
	ab25ap25Quality := ab25ap25Operation == 3 && ab25ap25Context["view"] == "quality" && ab25ap25Context["direction"] == "positive"
	ab25ap25Configuration := ab25ap25Operation >= 24 && ab25ap25Operation <= 27 && ab25ap25Context["view"] == "configuration" && (ab25ap25Context["direction"] == "request" || ab25ap25Context["direction"] == "positive")
	if ab25ap25Primary != 181 || ab25ap25Secondary != 37 || (!ab25ap25Quality && !ab25ap25Configuration) {
		return json.Marshal(ab25ap25Output)
	}
	ab25ap25Output["state"] = "complete"
	ab25ap25Names := []string{"rssi", "quality"}
	if ab25ap25Configuration {
		ab25ap25Names = []string{"configuration", "data"}
	}
	for _, ab25ap25Name := range ab25ap25Names {
		ab25ap25Result := map[string]any{"state": "Unknown"}
		ab25ap25Fields[ab25ap25Name] = ab25ap25Result
		ab25ap25Array, ab25ap25Exists := ab25ap25Raw[ab25ap25Name].([]any)
		if ab25ap25Exists {
			ab25ap25Result["raw"] = ab25ap25Array
		}
		if !ab25ap25Exists || (ab25ap25Name != "data" && len(ab25ap25Array) != 1) {
			ab25ap25Output["state"] = "partial"
			continue
		}
		ab25ap25Result["state"] = "Hypothesis"
		if ab25ap25Name == "data" {
			continue
		}
		ab25ap25Number, _ := ab25ap25Byte(ab25ap25Array[0])
		switch ab25ap25Name {
		case "rssi":
			if ab25ap25Number >= 128 {
				ab25ap25Number -= 256
			}
			ab25ap25Result["number"] = ab25ap25Number
			ab25ap25Result["unit"] = "Unknown"
			ab25ap25Result["calculation"] = "Proven"
		case "quality":
			ab25ap25Flags := map[string]bool{}
			for ab25ap25Bit := 0; ab25ap25Bit < 7; ab25ap25Bit++ {
				ab25ap25Flags[fmt.Sprintf("bit%d", ab25ap25Bit)] = ab25ap25Number&(1<<ab25ap25Bit) != 0
			}
			ab25ap25Flags["crc_okay"] = ab25ap25Number&128 != 0
			ab25ap25Result["flags"] = ab25ap25Flags
			ab25ap25Result["checksum_verified"] = false
		case "configuration":
			ab25ap25Result["flags"] = map[string]bool{"bit0": ab25ap25Number&1 != 0}
			ab25ap25Result["unmodeled"] = ab25ap25Number & 254
			ab25ap25Result["candidate_default"] = map[string]any{"number": 0, "label": "read", "state": "Hypothesis"}
		}
	}
	return json.Marshal(ab25ap25Output)
}
