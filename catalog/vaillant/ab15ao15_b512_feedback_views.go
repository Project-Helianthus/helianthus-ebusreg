package vaillant

import (
	ab15ao15bytes "bytes"
	ab15ao15json "encoding/json"
	ab15ao15fmt "fmt"
	ab15ao15io "io"
)

func Ab15Ao15Interpret(ab15ao15Input ab15ao15json.RawMessage) (ab15ao15json.RawMessage, error) {
	ab15ao15Decoder := ab15ao15json.NewDecoder(ab15ao15bytes.NewReader(ab15ao15Input))
	ab15ao15Decoder.UseNumber()
	if ab15ao15Err := ab15ao15Validate(ab15ao15Decoder); ab15ao15Err != nil {
		return nil, ab15ao15Err
	}
	if _, ab15ao15Err := ab15ao15Decoder.Token(); ab15ao15Err != ab15ao15io.EOF {
		return nil, ab15ao15fmt.Errorf("ab15ao15 trailing JSON")
	}
	ab15ao15Root, ab15ao15Err := ab15ao15Object(ab15ao15Input)
	if ab15ao15Err != nil {
		return nil, ab15ao15Err
	}
	for ab15ao15Key := range ab15ao15Root {
		if ab15ao15Key != "context" && ab15ao15Key != "fields" && ab15ao15Key != "arguments" {
			return nil, ab15ao15fmt.Errorf("ab15ao15 unknown input key")
		}
	}
	ab15ao15Context, ab15ao15Err := ab15ao15Object(ab15ao15Root["context"])
	if ab15ao15Err != nil {
		return nil, ab15ao15Err
	}
	ab15ao15Raw, ab15ao15Err := ab15ao15Object(ab15ao15Root["fields"])
	if ab15ao15Err != nil {
		return nil, ab15ao15Err
	}
	ab15ao15Bytes := make(map[string][]int)
	for ab15ao15Key, ab15ao15Value := range ab15ao15Raw {
		if len(ab15ao15bytes.TrimSpace(ab15ao15Value)) == 0 || ab15ao15bytes.TrimSpace(ab15ao15Value)[0] != '[' {
			return nil, ab15ao15fmt.Errorf("ab15ao15 field must be byte array")
		}
		var ab15ao15Elements []ab15ao15json.RawMessage
		if ab15ao15Err := ab15ao15json.Unmarshal(ab15ao15Value, &ab15ao15Elements); ab15ao15Err != nil {
			return nil, ab15ao15Err
		}
		ab15ao15Values := make([]int, len(ab15ao15Elements))
		for ab15ao15Index, ab15ao15Element := range ab15ao15Elements {
			if ab15ao15bytes.Equal(ab15ao15Element, []byte("null")) {
				return nil, ab15ao15fmt.Errorf("ab15ao15 null byte")
			}
			if ab15ao15Err := ab15ao15json.Unmarshal(ab15ao15Element, &ab15ao15Values[ab15ao15Index]); ab15ao15Err != nil || ab15ao15Values[ab15ao15Index] < 0 || ab15ao15Values[ab15ao15Index] > 255 {
				return nil, ab15ao15fmt.Errorf("ab15ao15 invalid byte")
			}
		}
		ab15ao15Bytes[ab15ao15Key] = ab15ao15Values
	}
	ab15ao15Selected := true
	for ab15ao15Key, ab15ao15Expected := range map[string]int{"primary": 181, "secondary": 18, "operation": 1} {
		ab15ao15Value, ab15ao15Present := ab15ao15Context[ab15ao15Key]
		if !ab15ao15Present {
			ab15ao15Selected = false
			continue
		}
		var ab15ao15Number int
		if ab15ao15bytes.Equal(ab15ao15Value, []byte("null")) {
			return nil, ab15ao15fmt.Errorf("ab15ao15 null identity")
		}
		if ab15ao15Err := ab15ao15json.Unmarshal(ab15ao15Value, &ab15ao15Number); ab15ao15Err != nil {
			return nil, ab15ao15Err
		}
		if ab15ao15Number != ab15ao15Expected {
			ab15ao15Selected = false
		}
	}
	ab15ao15Direction := ""
	for ab15ao15Key := range map[string]bool{"direction": true, "view": true} {
		ab15ao15Value, ab15ao15Present := ab15ao15Context[ab15ao15Key]
		if !ab15ao15Present {
			ab15ao15Selected = false
			continue
		}
		var ab15ao15Text string
		if ab15ao15bytes.Equal(ab15ao15Value, []byte("null")) {
			return nil, ab15ao15fmt.Errorf("ab15ao15 null selector")
		}
		if ab15ao15Err := ab15ao15json.Unmarshal(ab15ao15Value, &ab15ao15Text); ab15ao15Err != nil {
			return nil, ab15ao15Err
		}
		if ab15ao15Key == "direction" {
			ab15ao15Direction = ab15ao15Text
			if ab15ao15Text != "request" && ab15ao15Text != "positive" {
				ab15ao15Selected = false
			}
		} else if ab15ao15Text != "directional-fields" {
			ab15ao15Selected = false
		}
	}
	ab15ao15Origin := "origin-Unknown"
	if ab15ao15Arguments, ab15ao15Present := ab15ao15Root["arguments"]; ab15ao15Present {
		ab15ao15Flags, ab15ao15Err := ab15ao15Object(ab15ao15Arguments)
		if ab15ao15Err != nil {
			return nil, ab15ao15Err
		}
		for ab15ao15Key, ab15ao15Value := range ab15ao15Flags {
			if ab15ao15Key != "actual_signal" || (string(ab15ao15Value) != "true" && string(ab15ao15Value) != "false") {
				return nil, ab15ao15fmt.Errorf("ab15ao15 invalid argument")
			}
			if string(ab15ao15Value) == "true" {
				ab15ao15Origin = "actual-input-candidate"
			} else {
				ab15ao15Origin = "fallback-candidate"
			}
		}
	}
	ab15ao15Fields := make(map[string]any)
	ab15ao15State := "Unknown"
	if ab15ao15Selected {
		ab15ao15State = "Hypothesis"
		for _, ab15ao15Name := range []string{"secondary_heating_pump", "circulation_pump", "solar_pump", "error_relay", "loading_pump", "fume_flap"} {
			ab15ao15Field := map[string]any{"raw": nil, "state": "Unknown"}
			if ab15ao15Value, ab15ao15Present := ab15ao15Raw[ab15ao15Name]; ab15ao15Present {
				ab15ao15Field["raw"] = ab15ao15Value
			}
			if len(ab15ao15Bytes[ab15ao15Name]) != 1 {
				ab15ao15State = "partial"
				ab15ao15Fields[ab15ao15Name] = ab15ao15Field
				continue
			}
			ab15ao15Field["state"] = "Hypothesis"
			ab15ao15Field["flags"] = map[string]string{"origin": ab15ao15Origin}
			ab15ao15Byte := ab15ao15Bytes[ab15ao15Name][0]
			if ab15ao15Name == "error_relay" || ab15ao15Name == "fume_flap" {
				switch ab15ao15Byte {
				case 15:
					ab15ao15Field["label"] = "state-0F"
				case 240:
					ab15ao15Field["label"] = "state-F0"
				default:
					ab15ao15Field["label"] = "no-action-candidate"
				}
			} else if ab15ao15Byte == 255 {
				if ab15ao15Direction == "request" {
					ab15ao15Field["label"] = "leave-unchanged"
				} else {
					ab15ao15Field["label"] = "function-unavailable"
				}
			}
			ab15ao15Fields[ab15ao15Name] = ab15ao15Field
		}
	}
	ab15ao15Output := map[string]any{"context": ab15ao15Root["context"], "raw": ab15ao15Root["fields"], "state": ab15ao15State, "fields": ab15ao15Fields, "physical": "Unknown", "qualified": false}
	if ab15ao15Arguments, ab15ao15Present := ab15ao15Root["arguments"]; ab15ao15Present {
		ab15ao15Output["arguments"] = ab15ao15Arguments
	}
	return ab15ao15json.Marshal(ab15ao15Output)
}

func ab15ao15Object(ab15ao15Raw ab15ao15json.RawMessage) (map[string]ab15ao15json.RawMessage, error) {
	if len(ab15ao15bytes.TrimSpace(ab15ao15Raw)) == 0 || ab15ao15bytes.TrimSpace(ab15ao15Raw)[0] != '{' {
		return nil, ab15ao15fmt.Errorf("ab15ao15 expected object")
	}
	var ab15ao15Object map[string]ab15ao15json.RawMessage
	ab15ao15Err := ab15ao15json.Unmarshal(ab15ao15Raw, &ab15ao15Object)
	return ab15ao15Object, ab15ao15Err
}

func ab15ao15Validate(ab15ao15Decoder *ab15ao15json.Decoder) error {
	ab15ao15Token, ab15ao15Err := ab15ao15Decoder.Token()
	if ab15ao15Err != nil {
		return ab15ao15Err
	}
	ab15ao15Delimiter, ab15ao15Container := ab15ao15Token.(ab15ao15json.Delim)
	if !ab15ao15Container {
		return nil
	}
	switch ab15ao15Delimiter {
	case '{':
		ab15ao15Seen := make(map[string]bool)
		for ab15ao15Decoder.More() {
			ab15ao15KeyToken, ab15ao15Err := ab15ao15Decoder.Token()
			if ab15ao15Err != nil {
				return ab15ao15Err
			}
			ab15ao15Key, ab15ao15OK := ab15ao15KeyToken.(string)
			if !ab15ao15OK || ab15ao15Seen[ab15ao15Key] {
				return ab15ao15fmt.Errorf("ab15ao15 duplicate or invalid key")
			}
			ab15ao15Seen[ab15ao15Key] = true
			if ab15ao15Err := ab15ao15Validate(ab15ao15Decoder); ab15ao15Err != nil {
				return ab15ao15Err
			}
		}
	case '[':
		for ab15ao15Decoder.More() {
			if ab15ao15Err := ab15ao15Validate(ab15ao15Decoder); ab15ao15Err != nil {
				return ab15ao15Err
			}
		}
	default:
		return ab15ao15fmt.Errorf("ab15ao15 unexpected delimiter")
	}
	_, ab15ao15Err = ab15ao15Decoder.Token()
	return ab15ao15Err
}
