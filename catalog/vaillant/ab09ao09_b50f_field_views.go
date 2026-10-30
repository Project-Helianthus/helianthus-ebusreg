package vaillant

import (
	ab09ao09Bytes "bytes"
	ab09ao09JSON "encoding/json"
	ab09ao09Fmt "fmt"
	ab09ao09IO "io"
	ab09ao09Strconv "strconv"
)

func ab09ao09Read(ab09ao09Decoder *ab09ao09JSON.Decoder) (any, error) {
	ab09ao09Token, ab09ao09Err := ab09ao09Decoder.Token()
	if ab09ao09Err != nil {
		return nil, ab09ao09Err
	}
	if ab09ao09Delimiter, ab09ao09OK := ab09ao09Token.(ab09ao09JSON.Delim); ab09ao09OK {
		switch ab09ao09Delimiter {
		case '{':
			ab09ao09Object := map[string]any{}
			for ab09ao09Decoder.More() {
				ab09ao09KeyToken, ab09ao09Err := ab09ao09Decoder.Token()
				if ab09ao09Err != nil {
					return nil, ab09ao09Err
				}
				ab09ao09Key, ab09ao09OK := ab09ao09KeyToken.(string)
				if !ab09ao09OK {
					return nil, ab09ao09Fmt.Errorf("invalid object key")
				}
				if _, ab09ao09Found := ab09ao09Object[ab09ao09Key]; ab09ao09Found {
					return nil, ab09ao09Fmt.Errorf("duplicate object key")
				}
				ab09ao09Value, ab09ao09Err := ab09ao09Read(ab09ao09Decoder)
				if ab09ao09Err != nil {
					return nil, ab09ao09Err
				}
				ab09ao09Object[ab09ao09Key] = ab09ao09Value
			}
			_, ab09ao09Err := ab09ao09Decoder.Token()
			return ab09ao09Object, ab09ao09Err
		case '[':
			ab09ao09Array := []any{}
			for ab09ao09Decoder.More() {
				ab09ao09Value, ab09ao09Err := ab09ao09Read(ab09ao09Decoder)
				if ab09ao09Err != nil {
					return nil, ab09ao09Err
				}
				ab09ao09Array = append(ab09ao09Array, ab09ao09Value)
			}
			_, ab09ao09Err := ab09ao09Decoder.Token()
			return ab09ao09Array, ab09ao09Err
		}
		return nil, ab09ao09Fmt.Errorf("unexpected delimiter")
	}
	return ab09ao09Token, nil
}

func Ab09Ao09Interpret(ab09ao09Input ab09ao09JSON.RawMessage) (ab09ao09JSON.RawMessage, error) {
	ab09ao09Decoder := ab09ao09JSON.NewDecoder(ab09ao09Bytes.NewReader(ab09ao09Input))
	ab09ao09Decoder.UseNumber()
	ab09ao09Value, ab09ao09Err := ab09ao09Read(ab09ao09Decoder)
	if ab09ao09Err != nil {
		return nil, ab09ao09Err
	}
	if _, ab09ao09Err = ab09ao09Decoder.Token(); ab09ao09Err != ab09ao09IO.EOF {
		return nil, ab09ao09Fmt.Errorf("trailing JSON")
	}
	ab09ao09InputObject, ab09ao09OK := ab09ao09Value.(map[string]any)
	if !ab09ao09OK {
		return nil, ab09ao09Fmt.Errorf("input must be an object")
	}
	for ab09ao09Key := range ab09ao09InputObject {
		if ab09ao09Key != "context" && ab09ao09Key != "fields" && ab09ao09Key != "arguments" {
			return nil, ab09ao09Fmt.Errorf("unknown input key")
		}
	}
	ab09ao09Context, ab09ao09OK := ab09ao09InputObject["context"].(map[string]any)
	if !ab09ao09OK {
		return nil, ab09ao09Fmt.Errorf("context must be an object")
	}
	ab09ao09Raw, ab09ao09OK := ab09ao09InputObject["fields"].(map[string]any)
	if !ab09ao09OK {
		return nil, ab09ao09Fmt.Errorf("fields must be an object")
	}
	for _, ab09ao09Key := range []string{"primary", "secondary", "operation"} {
		if ab09ao09Value, ab09ao09Present := ab09ao09Context[ab09ao09Key]; ab09ao09Present {
			ab09ao09Number, ab09ao09OK := ab09ao09Value.(ab09ao09JSON.Number)
			if !ab09ao09OK {
				return nil, ab09ao09Fmt.Errorf("integer selector required")
			}
			if _, ab09ao09Err := ab09ao09Strconv.ParseInt(string(ab09ao09Number), 10, 64); ab09ao09Err != nil {
				return nil, ab09ao09Fmt.Errorf("integer selector required")
			}
		}
	}
	for _, ab09ao09Key := range []string{"direction", "view"} {
		if ab09ao09Value, ab09ao09Present := ab09ao09Context[ab09ao09Key]; ab09ao09Present {
			if _, ab09ao09OK := ab09ao09Value.(string); !ab09ao09OK {
				return nil, ab09ao09Fmt.Errorf("string selector required")
			}
		}
	}
	ab09ao09Arrays := map[string][]int{}
	for ab09ao09Key, ab09ao09Value := range ab09ao09Raw {
		ab09ao09Array, ab09ao09OK := ab09ao09Value.([]any)
		if !ab09ao09OK {
			return nil, ab09ao09Fmt.Errorf("field must be a byte array")
		}
		ab09ao09Bytes := []int{}
		for _, ab09ao09Element := range ab09ao09Array {
			ab09ao09Number, ab09ao09OK := ab09ao09Element.(ab09ao09JSON.Number)
			if !ab09ao09OK {
				return nil, ab09ao09Fmt.Errorf("invalid byte")
			}
			ab09ao09Byte, ab09ao09Err := ab09ao09Strconv.Atoi(string(ab09ao09Number))
			if ab09ao09Err != nil || ab09ao09Byte < 0 || ab09ao09Byte > 255 {
				return nil, ab09ao09Fmt.Errorf("invalid byte")
			}
			ab09ao09Bytes = append(ab09ao09Bytes, ab09ao09Byte)
		}
		ab09ao09Arrays[ab09ao09Key] = ab09ao09Bytes
	}
	ab09ao09Fields := map[string]any{}
	for ab09ao09Key, ab09ao09Value := range ab09ao09Raw {
		ab09ao09Fields[ab09ao09Key] = map[string]any{"raw": ab09ao09Value, "state": "Unknown"}
	}
	ab09ao09Output := map[string]any{"context": ab09ao09Context, "raw": ab09ao09Raw, "fields": ab09ao09Fields, "state": "Unknown", "physical": "Unknown", "qualified": false}
	if ab09ao09Value, ab09ao09Present := ab09ao09InputObject["arguments"]; ab09ao09Present {
		ab09ao09Arguments, ab09ao09OK := ab09ao09Value.(map[string]any)
		if !ab09ao09OK || len(ab09ao09Arguments) != 0 {
			return nil, ab09ao09Fmt.Errorf("no arguments defined for this view")
		}
		ab09ao09Output["arguments"] = ab09ao09Arguments
	}
	ab09ao09Operation := ab09ao09Context["operation"]
	ab09ao09View := ab09ao09Context["view"]
	if ab09ao09Context["primary"] != ab09ao09JSON.Number("181") || ab09ao09Context["secondary"] != ab09ao09JSON.Number("15") || ab09ao09Context["direction"] != "positive" {
		return ab09ao09JSON.Marshal(ab09ao09Output)
	}
	ab09ao09Selected := (ab09ao09Operation == ab09ao09JSON.Number("0") || ab09ao09Operation == ab09ao09JSON.Number("1")) && ab09ao09View == "thermal-fields" || ab09ao09Operation == ab09ao09JSON.Number("3") && ab09ao09View == "demand" || ab09ao09Operation == ab09ao09JSON.Number("2") && ab09ao09View == "integrity-only"
	if !ab09ao09Selected {
		return ab09ao09JSON.Marshal(ab09ao09Output)
	}
	ab09ao09Output["state"] = "complete"
	ab09ao09Field := func(ab09ao09Key string, ab09ao09Width int) (map[string]any, []int, bool) {
		ab09ao09Bytes, ab09ao09Present := ab09ao09Arrays[ab09ao09Key]
		ab09ao09Result := map[string]any{"state": "missing", "grade": "Hypothesis"}
		ab09ao09Fields[ab09ao09Key] = ab09ao09Result
		if ab09ao09Present {
			ab09ao09Result["raw"] = ab09ao09Raw[ab09ao09Key]
			ab09ao09Result["state"] = "incomplete"
		}
		if !ab09ao09Present || len(ab09ao09Bytes) != ab09ao09Width {
			ab09ao09Output["state"] = "partial"
			if len(ab09ao09Bytes) > ab09ao09Width {
				ab09ao09Result["state"] = "unmapped"
			}
			return ab09ao09Result, ab09ao09Bytes, false
		}
		ab09ao09Result["state"] = "complete"
		return ab09ao09Result, ab09ao09Bytes, true
	}
	if ab09ao09View == "thermal-fields" {
		for ab09ao09Index, ab09ao09Suffix := range []string{"one", "two", "three"} {
			ab09ao09Key := "temperature_" + ab09ao09Suffix
			if ab09ao09Index == 2 && ab09ao09Operation == ab09ao09JSON.Number("1") {
				ab09ao09Key = "stack_voltage"
			}
			ab09ao09Result, ab09ao09Bytes, ab09ao09Complete := ab09ao09Field(ab09ao09Key, 2)
			if ab09ao09Complete {
				if ab09ao09Key == "stack_voltage" {
					ab09ao09Result["state"] = "unit-conflict"
					ab09ao09Output["state"] = "partial"
				} else {
					ab09ao09Result["number"] = float64(int16(uint16(ab09ao09Bytes[0])|uint16(ab09ao09Bytes[1])<<8)) / 16
					ab09ao09Result["label"] = "Celsius"
					ab09ao09Result["calculation"] = "Proven"
				}
			}
			ab09ao09Field("error_"+ab09ao09Suffix, 1)
		}
	}
	if ab09ao09View == "demand" {
		ab09ao09Result, ab09ao09Bytes, ab09ao09Complete := ab09ao09Field("demand", 1)
		if ab09ao09Complete {
			ab09ao09Flags := map[string]bool{}
			for ab09ao09Index, ab09ao09Name := range []string{"burner", "reformer-gas-valve", "gas-flow", "fan-power", "valve"} {
				ab09ao09Flags[ab09ao09Name] = ab09ao09Bytes[0]&(1<<ab09ao09Index) != 0
			}
			ab09ao09Result["flags"] = ab09ao09Flags
			ab09ao09Result["reserved"] = ab09ao09Bytes[0] & 224
			ab09ao09Result["label"] = "demand"
		}
		ab09ao09Result, ab09ao09Bytes, ab09ao09Complete = ab09ao09Field("feedback", 1)
		if ab09ao09Complete {
			ab09ao09Result["label"] = "rejected"
			if ab09ao09Bytes[0] == 1 {
				ab09ao09Result["label"] = "accepted"
			}
		}
	}
	for ab09ao09Index, ab09ao09Key := range []string{"integrity_one", "integrity_two"} {
		if _, ab09ao09Present := ab09ao09Arrays[ab09ao09Key]; !ab09ao09Present && ab09ao09View != "integrity-only" {
			continue
		}
		ab09ao09Result, _, ab09ao09Complete := ab09ao09Field(ab09ao09Key, 1)
		if ab09ao09Complete {
			ab09ao09Position := 12
			if ab09ao09Operation == ab09ao09JSON.Number("2") {
				ab09ao09Position = 11
			}
			if ab09ao09Operation == ab09ao09JSON.Number("3") {
				ab09ao09Position = 4
			}
			ab09ao09Result["logical_position"] = ab09ao09Position + ab09ao09Index
			ab09ao09Result["state"] = "unknown-order"
			ab09ao09Result["word_order"] = "Unknown"
			ab09ao09Result["crc"] = "Unknown"
			ab09ao09Output["state"] = "partial"
		}
	}
	for _, ab09ao09Value := range ab09ao09Fields {
		if ab09ao09Value.(map[string]any)["state"] == "Unknown" {
			ab09ao09Output["state"] = "partial"
		}
	}
	return ab09ao09JSON.Marshal(ab09ao09Output)
}
