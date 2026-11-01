package vaillant

import (
	ab11ap11bytes "bytes"
	ab11ap11json "encoding/json"
	ab11ap11fmt "fmt"
	ab11ap11io "io"
	ab11ap11strconv "strconv"
)

func ab11ap11Read(ab11ap11Decoder *ab11ap11json.Decoder) (any, error) {
	ab11ap11Token, ab11ap11Err := ab11ap11Decoder.Token()
	if ab11ap11Err != nil {
		return nil, ab11ap11Err
	}
	switch ab11ap11Token {
	case ab11ap11json.Delim('{'):
		ab11ap11Object := map[string]any{}
		for ab11ap11Decoder.More() {
			ab11ap11KeyToken, ab11ap11Err := ab11ap11Decoder.Token()
			if ab11ap11Err != nil {
				return nil, ab11ap11Err
			}
			ab11ap11Key := ab11ap11KeyToken.(string)
			if _, ab11ap11Exists := ab11ap11Object[ab11ap11Key]; ab11ap11Exists {
				return nil, ab11ap11fmt.Errorf("duplicate key %q", ab11ap11Key)
			}
			ab11ap11Value, ab11ap11Err := ab11ap11Read(ab11ap11Decoder)
			if ab11ap11Err != nil {
				return nil, ab11ap11Err
			}
			ab11ap11Object[ab11ap11Key] = ab11ap11Value
		}
		_, ab11ap11Err := ab11ap11Decoder.Token()
		return ab11ap11Object, ab11ap11Err
	case ab11ap11json.Delim('['):
		ab11ap11Array := []any{}
		for ab11ap11Decoder.More() {
			ab11ap11Value, ab11ap11Err := ab11ap11Read(ab11ap11Decoder)
			if ab11ap11Err != nil {
				return nil, ab11ap11Err
			}
			ab11ap11Array = append(ab11ap11Array, ab11ap11Value)
		}
		_, ab11ap11Err := ab11ap11Decoder.Token()
		return ab11ap11Array, ab11ap11Err
	default:
		return ab11ap11Token, nil
	}
}

func ab11ap11Uint(ab11ap11Value any, ab11ap11Bits int) (uint64, bool) {
	ab11ap11Number, ab11ap11OK := ab11ap11Value.(ab11ap11json.Number)
	if !ab11ap11OK {
		return 0, false
	}
	ab11ap11Result, ab11ap11Err := ab11ap11strconv.ParseUint(string(ab11ap11Number), 10, ab11ap11Bits)
	return ab11ap11Result, ab11ap11Err == nil
}

func Ab11Ap11Evaluate(ab11ap11Input ab11ap11json.RawMessage) (ab11ap11json.RawMessage, error) {
	ab11ap11Decoder := ab11ap11json.NewDecoder(ab11ap11bytes.NewReader(ab11ap11Input))
	ab11ap11Decoder.UseNumber()
	ab11ap11Value, ab11ap11Err := ab11ap11Read(ab11ap11Decoder)
	if ab11ap11Err != nil {
		return nil, ab11ap11Err
	}
	if _, ab11ap11Err = ab11ap11Decoder.Token(); ab11ap11Err != ab11ap11io.EOF {
		return nil, ab11ap11fmt.Errorf("trailing JSON")
	}
	ab11ap11InputObject, ab11ap11OK := ab11ap11Value.(map[string]any)
	if !ab11ap11OK {
		return nil, ab11ap11fmt.Errorf("input must be object")
	}
	ab11ap11Context, ab11ap11OK := ab11ap11InputObject["context"].(map[string]any)
	if !ab11ap11OK {
		return nil, ab11ap11fmt.Errorf("context must be object")
	}
	ab11ap11Raw, ab11ap11OK := ab11ap11InputObject["fields"].(map[string]any)
	if !ab11ap11OK {
		return nil, ab11ap11fmt.Errorf("fields must be object")
	}
	for _, ab11ap11Key := range []string{"primary", "secondary", "operation"} {
		if ab11ap11Value, ab11ap11Exists := ab11ap11Context[ab11ap11Key]; ab11ap11Exists && ab11ap11Value != nil {
			if _, ab11ap11OK := ab11ap11Uint(ab11ap11Value, 8); !ab11ap11OK {
				return nil, ab11ap11fmt.Errorf("invalid %s", ab11ap11Key)
			}
		}
	}
	for _, ab11ap11Key := range []string{"view", "direction"} {
		if ab11ap11Value, ab11ap11Exists := ab11ap11Context[ab11ap11Key]; ab11ap11Exists && ab11ap11Value != nil {
			if _, ab11ap11OK := ab11ap11Value.(string); !ab11ap11OK {
				return nil, ab11ap11fmt.Errorf("invalid %s", ab11ap11Key)
			}
		}
	}
	for ab11ap11Key, ab11ap11Value := range ab11ap11Raw {
		ab11ap11Array, ab11ap11OK := ab11ap11Value.([]any)
		if !ab11ap11OK {
			return nil, ab11ap11fmt.Errorf("field %s must be byte array", ab11ap11Key)
		}
		for _, ab11ap11Byte := range ab11ap11Array {
			if _, ab11ap11OK := ab11ap11Uint(ab11ap11Byte, 8); !ab11ap11OK {
				return nil, ab11ap11fmt.Errorf("invalid byte in %s", ab11ap11Key)
			}
		}
	}
	ab11ap11Fields := map[string]any{}
	ab11ap11Output := map[string]any{"context": ab11ap11Context, "raw": ab11ap11Raw, "fields": ab11ap11Fields, "state": "Unknown", "physical": "Unknown", "qualified": false}
	if ab11ap11Arguments, ab11ap11Exists := ab11ap11InputObject["arguments"]; ab11ap11Exists {
		if _, ab11ap11OK := ab11ap11Arguments.(map[string]any); !ab11ap11OK {
			return nil, ab11ap11fmt.Errorf("arguments must be object")
		}
		ab11ap11Output["arguments"] = ab11ap11Arguments
	}
	for ab11ap11Key, ab11ap11Value := range ab11ap11Raw {
		ab11ap11Fields[ab11ap11Key] = map[string]any{"raw": ab11ap11Value, "state": "Unknown"}
	}
	ab11ap11Primary, _ := ab11ap11Uint(ab11ap11Context["primary"], 8)
	ab11ap11Secondary, _ := ab11ap11Uint(ab11ap11Context["secondary"], 8)
	ab11ap11Operation, _ := ab11ap11Uint(ab11ap11Context["operation"], 8)
	ab11ap11View, _ := ab11ap11Context["view"].(string)
	ab11ap11Direction, _ := ab11ap11Context["direction"].(string)
	ab11ap11Views := map[string]uint64{"fan": 11, "fan-control": 12, "digital-air": 13, "air-quality": 14}
	if ab11ap11Primary != 181 || ab11ap11Secondary != 18 || ab11ap11Views[ab11ap11View] == 0 || ab11ap11Views[ab11ap11View] != ab11ap11Operation || (ab11ap11Direction != "request" && ab11ap11Direction != "positive") {
		return ab11ap11json.Marshal(ab11ap11Output)
	}
	ab11ap11Specs := map[string]string{}
	switch ab11ap11View {
	case "fan":
		ab11ap11Specs = map[string]string{"fan1": "rpm", "fan2": "rpm"}
		if ab11ap11Direction == "request" {
			ab11ap11Specs["relay"] = "relay"
		} else {
			ab11ap11Specs["status"] = "status"
			for _, ab11ap11Key := range []string{"temperature1", "temperature2", "temperature3"} {
				ab11ap11Specs[ab11ap11Key] = "Celsius"
			}
		}
	case "fan-control":
		if ab11ap11Direction == "request" {
			ab11ap11Specs = map[string]string{"pulse": "u8", "p": "u8", "i": "u8", "d": "u8", "fan_count": "u8", "minimum_pwm": "PWM", "maximum_pwm": "PWM"}
		} else {
			ab11ap11Specs["feedback"] = "feedback"
		}
	case "digital-air":
		if ab11ap11Direction == "positive" {
			ab11ap11Specs = map[string]string{"digital": "digital", "temperature1": "Celsius", "temperature2": "Celsius"}
		}
	case "air-quality":
		if ab11ap11Direction == "positive" {
			ab11ap11Specs = map[string]string{"temperature": "Celsius", "humidity": "percent", "air_quality": "Unknown"}
		}
	}
	ab11ap11Output["state"] = "complete"
	for ab11ap11Key, ab11ap11Kind := range ab11ap11Specs {
		ab11ap11Width := 1
		if ab11ap11Kind == "rpm" || ab11ap11Kind == "Celsius" || ab11ap11Kind == "PWM" || ab11ap11Kind == "percent" || ab11ap11Kind == "Unknown" {
			ab11ap11Width = 2
		}
		ab11ap11Field := map[string]any{"state": "incomplete"}
		ab11ap11Fields[ab11ap11Key] = ab11ap11Field
		ab11ap11Value, ab11ap11Exists := ab11ap11Raw[ab11ap11Key]
		if ab11ap11Exists {
			ab11ap11Field["raw"] = ab11ap11Value
		}
		ab11ap11Array, _ := ab11ap11Value.([]any)
		if !ab11ap11Exists || len(ab11ap11Array) != ab11ap11Width {
			ab11ap11Output["state"] = "partial"
			continue
		}
		ab11ap11Low, _ := ab11ap11Uint(ab11ap11Array[0], 8)
		ab11ap11Word := ab11ap11Low
		if ab11ap11Width == 2 {
			ab11ap11High, _ := ab11ap11Uint(ab11ap11Array[1], 8)
			ab11ap11Word += ab11ap11High * 256
		}
		ab11ap11Field["state"] = "Hypothesis"
		switch ab11ap11Kind {
		case "relay":
			ab11ap11Field["flags"] = map[string]any{"relay1": map[bool]string{false: "15V", true: "tray-heater-relay"}[ab11ap11Word&1 != 0], "relay2": map[bool]string{false: "15V", true: "tray-heater-relay"}[ab11ap11Word&2 != 0]}
			ab11ap11Field["reserved"] = ab11ap11Word & 252
		case "status":
			if ab11ap11Label, ab11ap11OK := map[uint64]string{0: "idle", 1: "running", 255: "fault"}[ab11ap11Word]; ab11ap11OK {
				ab11ap11Field["label"] = ab11ap11Label
			} else {
				ab11ap11Field["state"] = "Unknown"
				ab11ap11Field["unmodeled"] = ab11ap11Word
			}
		case "feedback":
			if ab11ap11Word == 1 {
				ab11ap11Field["label"] = "accepted"
			} else {
				ab11ap11Field["label"] = "rejected"
			}
		case "digital":
			if ab11ap11Word <= 1 {
				ab11ap11Field["label"] = ab11ap11Word == 1
			} else {
				ab11ap11Field["state"] = "Unknown"
				ab11ap11Field["unmodeled"] = ab11ap11Word
			}
		default:
			ab11ap11Number := float64(ab11ap11Word)
			switch ab11ap11Kind {
			case "Celsius":
				ab11ap11Number = float64(int16(ab11ap11Word)) / 16
			case "PWM":
				ab11ap11Number = float64(int16(ab11ap11Word)) / 256
			case "percent":
				ab11ap11Number /= 16
			}
			ab11ap11Field["number"] = ab11ap11Number
			ab11ap11Field["calculation"] = "Proven"
			if ab11ap11Kind != "u8" {
				ab11ap11Field["unit"] = ab11ap11Kind
			}
			if ab11ap11Kind == "PWM" {
				ab11ap11Field["unit"] = "percent"
				ab11ap11Field["candidate_domain"] = true
				ab11ap11Field["candidate_default"] = 1
			}
		}
	}
	return ab11ap11json.Marshal(ab11ap11Output)
}
