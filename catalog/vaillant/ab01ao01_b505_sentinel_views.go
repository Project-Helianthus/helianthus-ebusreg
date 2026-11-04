package vaillant

import (
	ab01ao01Bytes "bytes"
	ab01ao01JSON "encoding/json"
	ab01ao01Fmt "fmt"
	ab01ao01IO "io"
	ab01ao01Strconv "strconv"
	ab01ao01Strings "strings"
)

func ab01ao01Read(ab01ao01Decoder *ab01ao01JSON.Decoder) (any, error) {
	ab01ao01Token, ab01ao01Error := ab01ao01Decoder.Token()
	if ab01ao01Error != nil {
		return nil, ab01ao01Error
	}
	if ab01ao01Delimiter, ab01ao01OK := ab01ao01Token.(ab01ao01JSON.Delim); ab01ao01OK {
		switch ab01ao01Delimiter {
		case '{':
			ab01ao01Object := map[string]any{}
			for ab01ao01Decoder.More() {
				ab01ao01KeyToken, ab01ao01Error := ab01ao01Decoder.Token()
				if ab01ao01Error != nil {
					return nil, ab01ao01Error
				}
				ab01ao01Key, ab01ao01OK := ab01ao01KeyToken.(string)
				if !ab01ao01OK {
					return nil, ab01ao01Fmt.Errorf("ab01ao01: invalid object key")
				}
				if _, ab01ao01Exists := ab01ao01Object[ab01ao01Key]; ab01ao01Exists {
					return nil, ab01ao01Fmt.Errorf("ab01ao01: duplicate key %q", ab01ao01Key)
				}
				ab01ao01Value, ab01ao01Error := ab01ao01Read(ab01ao01Decoder)
				if ab01ao01Error != nil {
					return nil, ab01ao01Error
				}
				ab01ao01Object[ab01ao01Key] = ab01ao01Value
			}
			_, ab01ao01Error := ab01ao01Decoder.Token()
			return ab01ao01Object, ab01ao01Error
		case '[':
			ab01ao01Array := []any{}
			for ab01ao01Decoder.More() {
				ab01ao01Value, ab01ao01Error := ab01ao01Read(ab01ao01Decoder)
				if ab01ao01Error != nil {
					return nil, ab01ao01Error
				}
				ab01ao01Array = append(ab01ao01Array, ab01ao01Value)
			}
			_, ab01ao01Error := ab01ao01Decoder.Token()
			return ab01ao01Array, ab01ao01Error
		default:
			return nil, ab01ao01Fmt.Errorf("ab01ao01: unexpected delimiter")
		}
	}
	return ab01ao01Token, nil
}

func Ab01Ao01Interpret(ab01ao01Input ab01ao01JSON.RawMessage) (ab01ao01JSON.RawMessage, error) {
	ab01ao01Decoder := ab01ao01JSON.NewDecoder(ab01ao01Bytes.NewReader(ab01ao01Input))
	ab01ao01Decoder.UseNumber()
	ab01ao01Value, ab01ao01Error := ab01ao01Read(ab01ao01Decoder)
	if ab01ao01Error != nil {
		return nil, ab01ao01Error
	}
	if _, ab01ao01Error := ab01ao01Decoder.Token(); ab01ao01Error != ab01ao01IO.EOF {
		return nil, ab01ao01Fmt.Errorf("ab01ao01: trailing JSON")
	}
	ab01ao01Root, ab01ao01OK := ab01ao01Value.(map[string]any)
	if !ab01ao01OK {
		return nil, ab01ao01Fmt.Errorf("ab01ao01: input must be an object")
	}
	for ab01ao01Key := range ab01ao01Root {
		if ab01ao01Key != "context" && ab01ao01Key != "fields" && ab01ao01Key != "arguments" {
			return nil, ab01ao01Fmt.Errorf("ab01ao01: unknown input key %q", ab01ao01Key)
		}
	}
	ab01ao01Context, ab01ao01OK := ab01ao01Root["context"].(map[string]any)
	if !ab01ao01OK {
		return nil, ab01ao01Fmt.Errorf("ab01ao01: context must be an object")
	}
	for _, ab01ao01Key := range []string{"primary", "secondary", "operation"} {
		if ab01ao01Value, ab01ao01Present := ab01ao01Context[ab01ao01Key]; ab01ao01Present {
			ab01ao01Number, ab01ao01OK := ab01ao01Value.(ab01ao01JSON.Number)
			if !ab01ao01OK {
				return nil, ab01ao01Fmt.Errorf("ab01ao01: %s must be an integer", ab01ao01Key)
			}
			if ab01ao01Strings.ContainsAny(string(ab01ao01Number), ".eE") {
				return nil, ab01ao01Fmt.Errorf("ab01ao01: %s must be an integer", ab01ao01Key)
			}
		}
	}
	for _, ab01ao01Key := range []string{"direction", "view"} {
		if ab01ao01Value, ab01ao01Present := ab01ao01Context[ab01ao01Key]; ab01ao01Present {
			if _, ab01ao01OK := ab01ao01Value.(string); !ab01ao01OK {
				return nil, ab01ao01Fmt.Errorf("ab01ao01: %s must be a string", ab01ao01Key)
			}
		}
	}
	ab01ao01Raw, ab01ao01OK := ab01ao01Root["fields"].(map[string]any)
	if !ab01ao01OK {
		return nil, ab01ao01Fmt.Errorf("ab01ao01: fields must be an object")
	}
	for ab01ao01Key, ab01ao01Value := range ab01ao01Raw {
		ab01ao01Array, ab01ao01OK := ab01ao01Value.([]any)
		if !ab01ao01OK {
			return nil, ab01ao01Fmt.Errorf("ab01ao01: %s must be a byte array", ab01ao01Key)
		}
		for _, ab01ao01Element := range ab01ao01Array {
			ab01ao01Number, ab01ao01OK := ab01ao01Element.(ab01ao01JSON.Number)
			if !ab01ao01OK {
				return nil, ab01ao01Fmt.Errorf("ab01ao01: invalid byte")
			}
			ab01ao01Byte, ab01ao01Error := ab01ao01Strconv.ParseInt(string(ab01ao01Number), 10, 64)
			if ab01ao01Error != nil || ab01ao01Byte < 0 || ab01ao01Byte > 255 {
				return nil, ab01ao01Fmt.Errorf("ab01ao01: invalid byte")
			}
		}
	}
	ab01ao01Fields := map[string]any{}
	ab01ao01Output := map[string]any{"context": ab01ao01Context, "raw": ab01ao01Raw, "fields": ab01ao01Fields, "state": "Unknown", "physical": "Unknown", "qualified": false}
	if ab01ao01Arguments, ab01ao01Present := ab01ao01Root["arguments"]; ab01ao01Present {
		ab01ao01Object, ab01ao01OK := ab01ao01Arguments.(map[string]any)
		if !ab01ao01OK || len(ab01ao01Object) != 0 {
			return nil, ab01ao01Fmt.Errorf("ab01ao01: no arguments are defined for this view")
		}
		ab01ao01Output["arguments"] = ab01ao01Object
	}
	if ab01ao01Equals(ab01ao01Context["primary"], 181) && ab01ao01Equals(ab01ao01Context["secondary"], 5) && ab01ao01Equals(ab01ao01Context["operation"], 50) && ab01ao01Context["direction"] == "request" && ab01ao01Context["view"] == "replacement-pair" {
		ab01ao01Known := 0
		ab01ao01Incomplete := false
		for _, ab01ao01Key := range []string{"temperature_one", "temperature_two"} {
			ab01ao01Value, ab01ao01Present := ab01ao01Raw[ab01ao01Key]
			ab01ao01Field := map[string]any{"raw": ab01ao01Value, "state": "missing"}
			ab01ao01Fields[ab01ao01Key] = ab01ao01Field
			if !ab01ao01Present {
				ab01ao01Incomplete = true
				continue
			}
			ab01ao01Array := ab01ao01Value.([]any)
			if len(ab01ao01Array) != 2 {
				ab01ao01Field["state"] = "incomplete"
				ab01ao01Incomplete = true
				continue
			}
			ab01ao01Field["grade"] = "Hypothesis"
			ab01ao01Field["state"] = "Unknown"
			ab01ao01Field["label"] = "numeric-unqualified"
			if ab01ao01Equals(ab01ao01Array[0], 0) && (ab01ao01Equals(ab01ao01Array[1], 128) || ab01ao01Equals(ab01ao01Array[1], 0)) {
				ab01ao01Field["state"] = "Hypothesis"
				ab01ao01Field["label"] = "circuit-default"
				if ab01ao01Equals(ab01ao01Array[1], 128) {
					ab01ao01Field["label"] = "preserve-existing"
				}
				ab01ao01Known++
			}
		}
		if ab01ao01Known == 2 {
			ab01ao01Output["state"] = "Hypothesis"
		} else if ab01ao01Known > 0 || ab01ao01Incomplete {
			ab01ao01Output["state"] = "partial"
		}
	}
	return ab01ao01JSON.Marshal(ab01ao01Output)
}

func ab01ao01Equals(ab01ao01Value any, ab01ao01Expected int64) bool {
	ab01ao01Number, ab01ao01OK := ab01ao01Value.(ab01ao01JSON.Number)
	if !ab01ao01OK {
		return false
	}
	ab01ao01Integer, ab01ao01Error := ab01ao01Number.Int64()
	return ab01ao01Error == nil && ab01ao01Integer == ab01ao01Expected
}
