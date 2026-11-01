package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func ab23ap23Decode(ab23ap23Decoder *json.Decoder) (any, error) {
	ab23ap23Token, ab23ap23Err := ab23ap23Decoder.Token()
	if ab23ap23Err != nil {
		return nil, ab23ap23Err
	}
	switch ab23ap23Token {
	case json.Delim('{'):
		ab23ap23Object := map[string]any{}
		for ab23ap23Decoder.More() {
			ab23ap23KeyToken, ab23ap23Err := ab23ap23Decoder.Token()
			if ab23ap23Err != nil {
				return nil, ab23ap23Err
			}
			ab23ap23Key := ab23ap23KeyToken.(string)
			if _, ab23ap23Exists := ab23ap23Object[ab23ap23Key]; ab23ap23Exists {
				return nil, fmt.Errorf("duplicate key %q", ab23ap23Key)
			}
			ab23ap23Value, ab23ap23Err := ab23ap23Decode(ab23ap23Decoder)
			if ab23ap23Err != nil {
				return nil, ab23ap23Err
			}
			ab23ap23Object[ab23ap23Key] = ab23ap23Value
		}
		_, ab23ap23Err := ab23ap23Decoder.Token()
		return ab23ap23Object, ab23ap23Err
	case json.Delim('['):
		ab23ap23Array := []any{}
		for ab23ap23Decoder.More() {
			ab23ap23Value, ab23ap23Err := ab23ap23Decode(ab23ap23Decoder)
			if ab23ap23Err != nil {
				return nil, ab23ap23Err
			}
			ab23ap23Array = append(ab23ap23Array, ab23ap23Value)
		}
		_, ab23ap23Err := ab23ap23Decoder.Token()
		return ab23ap23Array, ab23ap23Err
	default:
		if _, ab23ap23Delimiter := ab23ap23Token.(json.Delim); ab23ap23Delimiter {
			return nil, fmt.Errorf("unexpected delimiter")
		}
		return ab23ap23Token, nil
	}
}

func ab23ap23Byte(ab23ap23Value any) (uint8, bool) {
	ab23ap23Number, ab23ap23OK := ab23ap23Value.(json.Number)
	if !ab23ap23OK {
		return 0, false
	}
	ab23ap23Integer, ab23ap23Err := strconv.ParseUint(string(ab23ap23Number), 10, 8)
	return uint8(ab23ap23Integer), ab23ap23Err == nil
}

func Ab23Ap23Evaluate(ab23ap23Input json.RawMessage) (json.RawMessage, error) {
	ab23ap23Decoder := json.NewDecoder(bytes.NewReader(ab23ap23Input))
	ab23ap23Decoder.UseNumber()
	ab23ap23Decoded, ab23ap23Err := ab23ap23Decode(ab23ap23Decoder)
	if ab23ap23Err != nil {
		return nil, ab23ap23Err
	}
	if _, ab23ap23Err = ab23ap23Decoder.Token(); ab23ap23Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab23ap23Root, ab23ap23OK := ab23ap23Decoded.(map[string]any)
	if !ab23ap23OK {
		return nil, fmt.Errorf("input must be an object")
	}
	ab23ap23Context, ab23ap23OK := ab23ap23Root["context"].(map[string]any)
	if !ab23ap23OK {
		return nil, fmt.Errorf("context must be an object")
	}
	ab23ap23Raw, ab23ap23OK := ab23ap23Root["fields"].(map[string]any)
	if !ab23ap23OK {
		return nil, fmt.Errorf("fields must be an object")
	}
	for _, ab23ap23Key := range []string{"primary", "secondary", "operation"} {
		if ab23ap23Value, ab23ap23Present := ab23ap23Context[ab23ap23Key]; ab23ap23Present && ab23ap23Value != nil {
			if _, ab23ap23OK := ab23ap23Byte(ab23ap23Value); !ab23ap23OK {
				return nil, fmt.Errorf("%s must be a byte or null", ab23ap23Key)
			}
		}
	}
	for _, ab23ap23Key := range []string{"direction", "view"} {
		if ab23ap23Value, ab23ap23Present := ab23ap23Context[ab23ap23Key]; ab23ap23Present && ab23ap23Value != nil {
			if _, ab23ap23OK := ab23ap23Value.(string); !ab23ap23OK {
				return nil, fmt.Errorf("%s must be a string or null", ab23ap23Key)
			}
		}
	}
	for ab23ap23Key, ab23ap23Value := range ab23ap23Raw {
		ab23ap23Array, ab23ap23OK := ab23ap23Value.([]any)
		if !ab23ap23OK {
			return nil, fmt.Errorf("field %s must be a byte array", ab23ap23Key)
		}
		for _, ab23ap23Element := range ab23ap23Array {
			if _, ab23ap23OK := ab23ap23Byte(ab23ap23Element); !ab23ap23OK {
				return nil, fmt.Errorf("field %s contains an invalid byte", ab23ap23Key)
			}
		}
	}
	ab23ap23Results := map[string]any{}
	for ab23ap23Key, ab23ap23Value := range ab23ap23Raw {
		ab23ap23Results[ab23ap23Key] = map[string]any{"raw": ab23ap23Value, "state": "Unknown", "present": true}
	}
	ab23ap23Output := map[string]any{"context": ab23ap23Context, "raw": ab23ap23Raw, "state": "Unknown", "fields": ab23ap23Results, "physical": "Unknown", "qualified": false}
	if ab23ap23Arguments, ab23ap23Present := ab23ap23Root["arguments"]; ab23ap23Present {
		if _, ab23ap23OK := ab23ap23Arguments.(map[string]any); !ab23ap23OK {
			return nil, fmt.Errorf("arguments must be an object")
		}
		ab23ap23Output["arguments"] = ab23ap23Arguments
	}
	ab23ap23Primary, _ := ab23ap23Byte(ab23ap23Context["primary"])
	ab23ap23Secondary, _ := ab23ap23Byte(ab23ap23Context["secondary"])
	ab23ap23Operation, ab23ap23OperationPresent := ab23ap23Context["operation"]
	ab23ap23Direction, _ := ab23ap23Context["direction"].(string)
	if ab23ap23Primary == 181 && ab23ap23Secondary == 22 && ab23ap23OperationPresent && ab23ap23Operation == nil && ab23ap23Context["view"] == "energy-type" && (ab23ap23Direction == "request" || ab23ap23Direction == "read-positive" || ab23ap23Direction == "write-positive") {
		ab23ap23Value, ab23ap23Present := ab23ap23Raw["value"]
		ab23ap23Result := map[string]any{"raw": ab23ap23Value, "present": ab23ap23Present, "state": "incomplete"}
		ab23ap23Results["value"] = ab23ap23Result
		ab23ap23Output["state"] = "partial"
		if ab23ap23Array, ab23ap23OK := ab23ap23Value.([]any); ab23ap23OK && len(ab23ap23Array) == 1 {
			ab23ap23ByteValue, _ := ab23ap23Byte(ab23ap23Array[0])
			ab23ap23Output["state"] = "complete"
			ab23ap23Result["state"] = "Hypothesis"
			if ab23ap23Direction == "request" {
				ab23ap23Result["number"] = ab23ap23ByteValue
				ab23ap23Result["unmapped"] = true
			} else {
				ab23ap23Result["low_six"] = ab23ap23ByteValue & 63
				ab23ap23Result["calculation"] = "Proven"
				if ab23ap23Direction == "read-positive" {
					ab23ap23Result["flags"] = map[string]any{"flag6": ab23ap23ByteValue&64 != 0, "flag7": ab23ap23ByteValue&128 != 0}
				} else {
					ab23ap23Result["unmodeled_high"] = ab23ap23ByteValue & 192
				}
			}
		}
	}
	return json.Marshal(ab23ap23Output)
}
