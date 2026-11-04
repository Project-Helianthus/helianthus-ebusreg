package vaillant

import (
	ab13aq13bytes "bytes"
	ab13aq13json "encoding/json"
	ab13aq13fmt "fmt"
	ab13aq13io "io"
	ab13aq13math "math"
	ab13aq13strconv "strconv"
)

func ab13aq13Decode(ab13aq13Decoder *ab13aq13json.Decoder) (any, error) {
	ab13aq13Token, ab13aq13Error := ab13aq13Decoder.Token()
	if ab13aq13Error != nil {
		return nil, ab13aq13Error
	}
	if ab13aq13Number, ab13aq13OK := ab13aq13Token.(ab13aq13json.Number); ab13aq13OK {
		ab13aq13Float, ab13aq13Error := ab13aq13Number.Float64()
		if ab13aq13Error != nil || ab13aq13math.IsInf(ab13aq13Float, 0) || ab13aq13math.IsNaN(ab13aq13Float) {
			return nil, ab13aq13fmt.Errorf("nonfinite number")
		}
	}
	if ab13aq13Delimiter, ab13aq13OK := ab13aq13Token.(ab13aq13json.Delim); ab13aq13OK {
		switch ab13aq13Delimiter {
		case '{':
			ab13aq13Object := map[string]any{}
			for ab13aq13Decoder.More() {
				ab13aq13KeyToken, ab13aq13Error := ab13aq13Decoder.Token()
				if ab13aq13Error != nil {
					return nil, ab13aq13Error
				}
				ab13aq13Key, ab13aq13OK := ab13aq13KeyToken.(string)
				if !ab13aq13OK {
					return nil, ab13aq13fmt.Errorf("object key required")
				}
				if _, ab13aq13Exists := ab13aq13Object[ab13aq13Key]; ab13aq13Exists {
					return nil, ab13aq13fmt.Errorf("duplicate key %q", ab13aq13Key)
				}
				ab13aq13Value, ab13aq13Error := ab13aq13Decode(ab13aq13Decoder)
				if ab13aq13Error != nil {
					return nil, ab13aq13Error
				}
				ab13aq13Object[ab13aq13Key] = ab13aq13Value
			}
			_, ab13aq13Error := ab13aq13Decoder.Token()
			return ab13aq13Object, ab13aq13Error
		case '[':
			ab13aq13Array := []any{}
			for ab13aq13Decoder.More() {
				ab13aq13Value, ab13aq13Error := ab13aq13Decode(ab13aq13Decoder)
				if ab13aq13Error != nil {
					return nil, ab13aq13Error
				}
				ab13aq13Array = append(ab13aq13Array, ab13aq13Value)
			}
			_, ab13aq13Error := ab13aq13Decoder.Token()
			return ab13aq13Array, ab13aq13Error
		default:
			return nil, ab13aq13fmt.Errorf("unexpected delimiter")
		}
	}
	return ab13aq13Token, nil
}

func ab13aq13Byte(ab13aq13Value any) (uint64, error) {
	ab13aq13Number, ab13aq13OK := ab13aq13Value.(ab13aq13json.Number)
	if !ab13aq13OK {
		return 0, ab13aq13fmt.Errorf("byte integer required")
	}
	return ab13aq13strconv.ParseUint(string(ab13aq13Number), 10, 8)
}

func Ab13Aq13Evaluate(ab13aq13Input ab13aq13json.RawMessage) (ab13aq13json.RawMessage, error) {
	ab13aq13Decoder := ab13aq13json.NewDecoder(ab13aq13bytes.NewReader(ab13aq13Input))
	ab13aq13Decoder.UseNumber()
	ab13aq13Value, ab13aq13Error := ab13aq13Decode(ab13aq13Decoder)
	if ab13aq13Error != nil {
		return nil, ab13aq13Error
	}
	if _, ab13aq13Error = ab13aq13Decoder.Token(); ab13aq13Error != ab13aq13io.EOF {
		return nil, ab13aq13fmt.Errorf("trailing JSON")
	}
	ab13aq13Output, ab13aq13OK := ab13aq13Value.(map[string]any)
	if !ab13aq13OK {
		return nil, ab13aq13fmt.Errorf("input object required")
	}
	ab13aq13Context, ab13aq13OK := ab13aq13Output["context"].(map[string]any)
	if !ab13aq13OK {
		return nil, ab13aq13fmt.Errorf("context object required")
	}
	ab13aq13Fields, ab13aq13OK := ab13aq13Output["fields"].(map[string]any)
	if !ab13aq13OK {
		return nil, ab13aq13fmt.Errorf("fields object required")
	}
	for ab13aq13Name, ab13aq13Value := range ab13aq13Fields {
		ab13aq13Array, ab13aq13OK := ab13aq13Value.([]any)
		if !ab13aq13OK {
			return nil, ab13aq13fmt.Errorf("field %s requires byte array", ab13aq13Name)
		}
		for _, ab13aq13Value := range ab13aq13Array {
			if _, ab13aq13Error := ab13aq13Byte(ab13aq13Value); ab13aq13Error != nil {
				return nil, ab13aq13Error
			}
		}
	}
	ab13aq13Arguments := map[string]any{}
	if ab13aq13Value, ab13aq13Present := ab13aq13Output["arguments"]; ab13aq13Present {
		ab13aq13Arguments, ab13aq13OK = ab13aq13Value.(map[string]any)
		if !ab13aq13OK {
			return nil, ab13aq13fmt.Errorf("arguments object required")
		}
	}
	ab13aq13Supported := true
	for ab13aq13Name, ab13aq13Expected := range map[string]uint64{"primary": 181, "secondary": 33} {
		ab13aq13Value, ab13aq13Present := ab13aq13Context[ab13aq13Name]
		if !ab13aq13Present || ab13aq13Value == nil {
			ab13aq13Supported = false
			continue
		}
		ab13aq13Number, ab13aq13Error := ab13aq13Byte(ab13aq13Value)
		if ab13aq13Error != nil {
			return nil, ab13aq13Error
		}
		if ab13aq13Number != ab13aq13Expected {
			ab13aq13Supported = false
		}
	}
	if ab13aq13Value, ab13aq13Present := ab13aq13Context["operation"]; ab13aq13Present {
		ab13aq13Supported = false
		if ab13aq13Value != nil {
			if _, ab13aq13Error := ab13aq13Byte(ab13aq13Value); ab13aq13Error != nil {
				return nil, ab13aq13Error
			}
		}
	}
	for _, ab13aq13Name := range []string{"direction", "view"} {
		ab13aq13Value, ab13aq13Present := ab13aq13Context[ab13aq13Name]
		if !ab13aq13Present || ab13aq13Value == nil {
			ab13aq13Supported = false
			continue
		}
		ab13aq13String, ab13aq13OK := ab13aq13Value.(string)
		if !ab13aq13OK {
			return nil, ab13aq13fmt.Errorf("context %s requires string", ab13aq13Name)
		}
		if (ab13aq13Name == "view" && ab13aq13String != "module-fields") || (ab13aq13Name == "direction" && ab13aq13String != "request" && ab13aq13String != "positive") {
			ab13aq13Supported = false
		}
	}
	ab13aq13Results := map[string]any{}
	for ab13aq13Name, ab13aq13Raw := range ab13aq13Fields {
		ab13aq13Results[ab13aq13Name] = map[string]any{"raw": ab13aq13Raw, "state": "Unknown", "unmodeled": true}
	}
	ab13aq13Output["results"] = ab13aq13Results
	ab13aq13Output["state"] = "Unknown"
	ab13aq13Output["physical"] = "Unknown"
	ab13aq13Output["qualified"] = false
	if !ab13aq13Supported {
		return ab13aq13json.Marshal(ab13aq13Output)
	}
	if len(ab13aq13Arguments) != 0 {
		return nil, ab13aq13fmt.Errorf("module-fields declares no arguments")
	}
	ab13aq13Widths := map[string]int{"flags": 1, "appliance_bits": 1, "noise_percent": 1, "frost_temperature": 1}
	ab13aq13Labels := map[string][]string{"flags": {"bit0", "bit1", "bit2", "bit3", "bit4", "bit5"}, "appliance_bits": {"bit0", "bit1", "bit2", "bit3", "bit4"}}
	if ab13aq13Context["direction"] == "positive" {
		ab13aq13Widths = map[string]int{"status0": 1, "status1": 1, "brine_percent": 1, "fan_percent": 1, "fan_speed": 1, "temperature0": 2, "temperature1": 2}
		ab13aq13Labels = map[string][]string{"status0": {"ready", "deicing-request", "deicing-active", "deicing-timeout-error", "deicing-with-heater", "", "", "cleaning-request"}, "status1": {"heater", "relay-output", "fan-running", "brine-flow-reverse", "fan-error", "safety-limiter-error", "hot-summer", "bad-weather"}}
	}
	ab13aq13Output["state"] = "complete"
	for ab13aq13Name, ab13aq13Width := range ab13aq13Widths {
		ab13aq13Raw, ab13aq13Present := ab13aq13Fields[ab13aq13Name]
		ab13aq13Result := map[string]any{"raw": ab13aq13Raw, "present": ab13aq13Present, "state": "incomplete", "grade": "Hypothesis"}
		ab13aq13Results[ab13aq13Name] = ab13aq13Result
		ab13aq13Array, _ := ab13aq13Raw.([]any)
		if len(ab13aq13Array) != ab13aq13Width {
			ab13aq13Output["state"] = "partial"
			continue
		}
		ab13aq13Result["state"] = "complete"
		ab13aq13Result["calculation"] = "Proven"
		ab13aq13Number, _ := ab13aq13Byte(ab13aq13Array[0])
		if ab13aq13Names, ab13aq13IsFlags := ab13aq13Labels[ab13aq13Name]; ab13aq13IsFlags {
			ab13aq13Flags := map[string]bool{}
			ab13aq13Reserved := ab13aq13Number
			for ab13aq13Bit, ab13aq13Label := range ab13aq13Names {
				if ab13aq13Label == "" {
					continue
				}
				ab13aq13Flags[ab13aq13Label] = ab13aq13Number&(1<<ab13aq13Bit) != 0
				ab13aq13Reserved &^= 1 << ab13aq13Bit
			}
			ab13aq13Result["flags"] = ab13aq13Flags
			ab13aq13Result["reserved"] = ab13aq13Reserved
			continue
		}
		ab13aq13Scalar := float64(ab13aq13Number)
		switch ab13aq13Name {
		case "fan_speed":
			ab13aq13Scalar *= 10
			ab13aq13Result["unit"] = "rpm"
		case "temperature0", "temperature1":
			ab13aq13High, _ := ab13aq13Byte(ab13aq13Array[1])
			ab13aq13Scalar = float64(int16(uint16(ab13aq13Number)|uint16(ab13aq13High)<<8)) / 16
			ab13aq13Result["unit"] = "Celsius"
		case "frost_temperature":
			ab13aq13Result["unit"] = "Celsius"
		default:
			ab13aq13Result["unit"] = "percent"
		}
		ab13aq13Result["number"] = ab13aq13Scalar
	}
	return ab13aq13json.Marshal(ab13aq13Output)
}
