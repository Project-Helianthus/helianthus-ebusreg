package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func ab01aq01JSON(ab01aq01Decoder *json.Decoder) error {
	ab01aq01Token, ab01aq01Err := ab01aq01Decoder.Token()
	if ab01aq01Err != nil {
		return ab01aq01Err
	}
	if ab01aq01Number, ab01aq01OK := ab01aq01Token.(json.Number); ab01aq01OK {
		if _, ab01aq01Err = strconv.ParseFloat(string(ab01aq01Number), 64); ab01aq01Err != nil {
			return ab01aq01Err
		}
	}
	if ab01aq01Delimiter, ab01aq01OK := ab01aq01Token.(json.Delim); ab01aq01OK {
		switch ab01aq01Delimiter {
		case '{':
			ab01aq01Seen := map[string]bool{}
			for ab01aq01Decoder.More() {
				ab01aq01Key, ab01aq01Err := ab01aq01Decoder.Token()
				if ab01aq01Err != nil {
					return ab01aq01Err
				}
				ab01aq01Name, ab01aq01OK := ab01aq01Key.(string)
				if !ab01aq01OK || ab01aq01Seen[ab01aq01Name] {
					return fmt.Errorf("duplicate or invalid object key")
				}
				ab01aq01Seen[ab01aq01Name] = true
				if ab01aq01Err = ab01aq01JSON(ab01aq01Decoder); ab01aq01Err != nil {
					return ab01aq01Err
				}
			}
		case '[':
			for ab01aq01Decoder.More() {
				if ab01aq01Err = ab01aq01JSON(ab01aq01Decoder); ab01aq01Err != nil {
					return ab01aq01Err
				}
			}
		default:
			return fmt.Errorf("unexpected delimiter")
		}
		_, ab01aq01Err = ab01aq01Decoder.Token()
	}
	return ab01aq01Err
}

func ab01aq01Object(ab01aq01Raw json.RawMessage) (map[string]json.RawMessage, error) {
	ab01aq01Trimmed := bytes.TrimSpace(ab01aq01Raw)
	if len(ab01aq01Trimmed) == 0 || ab01aq01Trimmed[0] != '{' {
		return nil, fmt.Errorf("expected object")
	}
	var ab01aq01Value map[string]json.RawMessage
	ab01aq01Err := json.Unmarshal(ab01aq01Raw, &ab01aq01Value)
	return ab01aq01Value, ab01aq01Err
}

func ab01aq01Integer(ab01aq01Raw json.RawMessage) (int, error) {
	ab01aq01Value, ab01aq01Err := strconv.Atoi(string(ab01aq01Raw))
	if ab01aq01Err != nil || ab01aq01Value < 0 || ab01aq01Value > 255 {
		return 0, fmt.Errorf("expected integer in 0..255")
	}
	return ab01aq01Value, nil
}

func ab01aq01String(ab01aq01Raw json.RawMessage) (string, error) {
	if len(ab01aq01Raw) == 0 || ab01aq01Raw[0] != '"' {
		return "", fmt.Errorf("expected string")
	}
	var ab01aq01Value string
	ab01aq01Err := json.Unmarshal(ab01aq01Raw, &ab01aq01Value)
	return ab01aq01Value, ab01aq01Err
}

func Ab01Aq01Evaluate(ab01aq01InputRaw json.RawMessage) (json.RawMessage, error) {
	ab01aq01Decoder := json.NewDecoder(bytes.NewReader(ab01aq01InputRaw))
	ab01aq01Decoder.UseNumber()
	if ab01aq01Err := ab01aq01JSON(ab01aq01Decoder); ab01aq01Err != nil {
		return nil, ab01aq01Err
	}
	if _, ab01aq01Err := ab01aq01Decoder.Token(); ab01aq01Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab01aq01Input, ab01aq01Err := ab01aq01Object(ab01aq01InputRaw)
	if ab01aq01Err != nil {
		return nil, ab01aq01Err
	}
	ab01aq01Context, ab01aq01Err := ab01aq01Object(ab01aq01Input["context"])
	if ab01aq01Err != nil {
		return nil, ab01aq01Err
	}
	ab01aq01Fields, ab01aq01Err := ab01aq01Object(ab01aq01Input["fields"])
	if ab01aq01Err != nil {
		return nil, ab01aq01Err
	}
	ab01aq01Arguments := map[string]json.RawMessage{}
	if ab01aq01Raw, ab01aq01Present := ab01aq01Input["arguments"]; ab01aq01Present {
		ab01aq01Arguments, ab01aq01Err = ab01aq01Object(ab01aq01Raw)
		if ab01aq01Err != nil {
			return nil, ab01aq01Err
		}
	}
	ab01aq01Integers := map[string]int{}
	for _, ab01aq01Key := range []string{"primary", "secondary", "operation"} {
		if ab01aq01Raw, ab01aq01Present := ab01aq01Context[ab01aq01Key]; ab01aq01Present {
			ab01aq01Integers[ab01aq01Key], ab01aq01Err = ab01aq01Integer(ab01aq01Raw)
			if ab01aq01Err != nil {
				return nil, ab01aq01Err
			}
		}
	}
	ab01aq01Strings := map[string]string{}
	for _, ab01aq01Key := range []string{"direction", "view", "family", "variant"} {
		if ab01aq01Raw, ab01aq01Present := ab01aq01Context[ab01aq01Key]; ab01aq01Present {
			ab01aq01Strings[ab01aq01Key], ab01aq01Err = ab01aq01String(ab01aq01Raw)
			if ab01aq01Err != nil {
				return nil, ab01aq01Err
			}
		}
	}
	ab01aq01Bytes := map[string][]int{}
	ab01aq01Results := map[string]any{}
	for ab01aq01Key, ab01aq01Raw := range ab01aq01Fields {
		if len(ab01aq01Raw) == 0 || ab01aq01Raw[0] != '[' {
			return nil, fmt.Errorf("field must be byte array")
		}
		var ab01aq01Array []json.RawMessage
		if ab01aq01Err = json.Unmarshal(ab01aq01Raw, &ab01aq01Array); ab01aq01Err != nil {
			return nil, ab01aq01Err
		}
		for _, ab01aq01Element := range ab01aq01Array {
			ab01aq01Byte, ab01aq01Err := ab01aq01Integer(ab01aq01Element)
			if ab01aq01Err != nil {
				return nil, ab01aq01Err
			}
			ab01aq01Bytes[ab01aq01Key] = append(ab01aq01Bytes[ab01aq01Key], ab01aq01Byte)
		}
		ab01aq01Results[ab01aq01Key] = map[string]any{"raw": ab01aq01Raw, "state": "Unknown", "present": true}
	}
	ab01aq01Divisor := 0
	ab01aq01Role := ""
	for ab01aq01Key, ab01aq01Raw := range ab01aq01Arguments {
		switch ab01aq01Key {
		case "divisor":
			ab01aq01Divisor, ab01aq01Err = strconv.Atoi(string(ab01aq01Raw))
			if ab01aq01Divisor < 0 || ab01aq01Divisor > 256 {
				return nil, fmt.Errorf("invalid divisor")
			}
		case "role":
			ab01aq01Role, ab01aq01Err = ab01aq01String(ab01aq01Raw)
		}
		if ab01aq01Err != nil {
			return nil, ab01aq01Err
		}
	}
	ab01aq01View := ab01aq01Strings["view"]
	ab01aq01Operations := map[string]int{"corrected-outside": 22, "demand": 32, "protection": 33, "adjustment": 38, "circuit-role": 13}
	ab01aq01Operation, ab01aq01Known := ab01aq01Operations[ab01aq01View]
	_, ab01aq01HasOperation := ab01aq01Context["operation"]
	ab01aq01Supported := ab01aq01Known && ab01aq01HasOperation && ab01aq01Integers["primary"] == 181 && ab01aq01Integers["secondary"] == 4 && ab01aq01Integers["operation"] == ab01aq01Operation && ab01aq01Strings["direction"] == "positive"
	ab01aq01State := "Unknown"
	if ab01aq01Supported {
		ab01aq01Argument := ""
		if ab01aq01View == "corrected-outside" {
			ab01aq01Argument = "divisor"
			if ab01aq01Divisor != 16 && ab01aq01Divisor != 256 {
				return nil, fmt.Errorf("divisor must be explicitly 16 or 256")
			}
		}
		if ab01aq01View == "circuit-role" {
			ab01aq01Argument = "role"
			if _, ab01aq01Present := ab01aq01Arguments["role"]; !ab01aq01Present {
				return nil, fmt.Errorf("missing role")
			}
		}
		for ab01aq01Key := range ab01aq01Arguments {
			if ab01aq01Key != ab01aq01Argument {
				return nil, fmt.Errorf("unexpected argument")
			}
		}
		ab01aq01Widths := map[string]int{}
		switch ab01aq01View {
		case "corrected-outside":
			ab01aq01Widths["temperature"] = 2
		case "demand":
			ab01aq01Widths["setpoint"] = 1
		case "protection", "circuit-role":
			ab01aq01Widths["temperature"] = 1
		case "adjustment":
			ab01aq01Widths["adjustment"] = 1
			ab01aq01Widths["room"] = 1
		}
		if ab01aq01View == "circuit-role" {
			if _, ab01aq01Present := ab01aq01Fields["dhw_request"]; ab01aq01Present {
				ab01aq01Widths["dhw_request"] = 1
			}
		}
		ab01aq01State = "complete"
		for ab01aq01Key, ab01aq01Width := range ab01aq01Widths {
			ab01aq01Raw, ab01aq01Present := ab01aq01Fields[ab01aq01Key]
			ab01aq01Result := map[string]any{"raw": ab01aq01Raw, "present": ab01aq01Present, "state": "incomplete"}
			ab01aq01Results[ab01aq01Key] = ab01aq01Result
			ab01aq01Data := ab01aq01Bytes[ab01aq01Key]
			if !ab01aq01Present || len(ab01aq01Data) != ab01aq01Width {
				ab01aq01State = "partial"
				continue
			}
			ab01aq01Result["state"] = "complete"
			ab01aq01Result["grade"] = "Hypothesis"
			switch ab01aq01View {
			case "corrected-outside":
				ab01aq01Result["number"] = float64(int16(uint16(ab01aq01Data[0])|uint16(ab01aq01Data[1])<<8)) / float64(ab01aq01Divisor)
				ab01aq01Result["calculation"] = "Proven"
			case "demand":
				if ab01aq01Data[0] == 0 {
					ab01aq01Result["label"] = "no-demand"
				} else {
					ab01aq01Result["label"] = "offset-included-candidate"
				}
			case "protection":
				if ab01aq01Data[0] == 0 {
					ab01aq01Result["label"] = "off"
				} else {
					ab01aq01Result["label"] = "protection-active"
				}
			case "adjustment":
				if ab01aq01Key == "adjustment" && ab01aq01Data[0] == 1 {
					ab01aq01Result["label"] = "change-permitted"
				}
				if ab01aq01Key == "room" && ab01aq01Data[0] == 255 {
					ab01aq01Result["label"] = "unavailable-candidate"
				}
			case "circuit-role":
				ab01aq01Labels := map[string]string{"burner": "common-sensor", "mixer": "flow", "dhw": "tank"}
				if ab01aq01Label, ab01aq01OK := ab01aq01Labels[ab01aq01Role]; ab01aq01OK && ab01aq01Key == "temperature" {
					ab01aq01Result["label"] = ab01aq01Label
					ab01aq01Result["number"] = ab01aq01Data[0]
				} else if ab01aq01Key == "dhw_request" && ab01aq01Role == "dhw" {
					ab01aq01Result["label"] = "dhw"
				} else {
					ab01aq01Result["state"] = "Unknown"
					delete(ab01aq01Result, "grade")
				}
			}
		}
	}
	ab01aq01Output := map[string]any{}
	for ab01aq01Key, ab01aq01Raw := range ab01aq01Input {
		ab01aq01Output[ab01aq01Key] = ab01aq01Raw
	}
	ab01aq01Output["state"] = ab01aq01State
	ab01aq01Output["results"] = ab01aq01Results
	ab01aq01Output["physical"] = "Unknown"
	ab01aq01Output["qualified"] = false
	return json.Marshal(ab01aq01Output)
}
