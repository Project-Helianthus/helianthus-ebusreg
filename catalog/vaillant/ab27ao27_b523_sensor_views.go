package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

func ab27ao27Read(ab27ao27Decoder *json.Decoder) (any, error) {
	ab27ao27Token, ab27ao27Err := ab27ao27Decoder.Token()
	if ab27ao27Err != nil {
		return nil, ab27ao27Err
	}
	if ab27ao27Delimiter, ab27ao27OK := ab27ao27Token.(json.Delim); ab27ao27OK {
		switch ab27ao27Delimiter {
		case '{':
			ab27ao27Object := map[string]any{}
			for ab27ao27Decoder.More() {
				ab27ao27KeyToken, ab27ao27Err := ab27ao27Decoder.Token()
				if ab27ao27Err != nil {
					return nil, ab27ao27Err
				}
				ab27ao27Key, ab27ao27OK := ab27ao27KeyToken.(string)
				if !ab27ao27OK {
					return nil, fmt.Errorf("expected object key")
				}
				if _, ab27ao27Exists := ab27ao27Object[ab27ao27Key]; ab27ao27Exists {
					return nil, fmt.Errorf("duplicate object key")
				}
				ab27ao27Value, ab27ao27Err := ab27ao27Read(ab27ao27Decoder)
				if ab27ao27Err != nil {
					return nil, ab27ao27Err
				}
				ab27ao27Object[ab27ao27Key] = ab27ao27Value
			}
			ab27ao27End, ab27ao27Err := ab27ao27Decoder.Token()
			if ab27ao27Err != nil || ab27ao27End != json.Delim('}') {
				return nil, fmt.Errorf("invalid object end")
			}
			return ab27ao27Object, nil
		case '[':
			ab27ao27Array := []any{}
			for ab27ao27Decoder.More() {
				ab27ao27Value, ab27ao27Err := ab27ao27Read(ab27ao27Decoder)
				if ab27ao27Err != nil {
					return nil, ab27ao27Err
				}
				ab27ao27Array = append(ab27ao27Array, ab27ao27Value)
			}
			ab27ao27End, ab27ao27Err := ab27ao27Decoder.Token()
			if ab27ao27Err != nil || ab27ao27End != json.Delim(']') {
				return nil, fmt.Errorf("invalid array end")
			}
			return ab27ao27Array, nil
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	return ab27ao27Token, nil
}

func ab27ao27Integer(ab27ao27Value any) (int, bool) {
	ab27ao27Number, ab27ao27OK := ab27ao27Value.(json.Number)
	if !ab27ao27OK {
		return 0, false
	}
	ab27ao27Int, ab27ao27Err := strconv.Atoi(string(ab27ao27Number))
	return ab27ao27Int, ab27ao27Err == nil
}

func Ab27Ao27Interpret(ab27ao27Input json.RawMessage) (json.RawMessage, error) {
	if !utf8.Valid(ab27ao27Input) {
		return nil, fmt.Errorf("invalid UTF-8 JSON")
	}
	ab27ao27Decoder := json.NewDecoder(bytes.NewReader(ab27ao27Input))
	ab27ao27Decoder.UseNumber()
	ab27ao27Value, ab27ao27Err := ab27ao27Read(ab27ao27Decoder)
	if ab27ao27Err != nil {
		return nil, ab27ao27Err
	}
	if _, ab27ao27Err = ab27ao27Decoder.Token(); ab27ao27Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab27ao27Root, ab27ao27OK := ab27ao27Value.(map[string]any)
	if !ab27ao27OK {
		return nil, fmt.Errorf("input must be an object")
	}
	for ab27ao27Key := range ab27ao27Root {
		if ab27ao27Key != "context" && ab27ao27Key != "fields" && ab27ao27Key != "arguments" {
			return nil, fmt.Errorf("unknown input member")
		}
	}
	ab27ao27Context, ab27ao27OK := ab27ao27Root["context"].(map[string]any)
	if !ab27ao27OK {
		return nil, fmt.Errorf("context must be an object")
	}
	ab27ao27Raw, ab27ao27OK := ab27ao27Root["fields"].(map[string]any)
	if !ab27ao27OK {
		return nil, fmt.Errorf("fields must be an object")
	}
	for _, ab27ao27Key := range []string{"primary", "secondary", "operation"} {
		if ab27ao27Value, ab27ao27Present := ab27ao27Context[ab27ao27Key]; ab27ao27Present {
			if _, ab27ao27OK := ab27ao27Integer(ab27ao27Value); !ab27ao27OK {
				return nil, fmt.Errorf("context %s must be an integer", ab27ao27Key)
			}
		}
	}
	for _, ab27ao27Key := range []string{"direction", "view"} {
		if ab27ao27Value, ab27ao27Present := ab27ao27Context[ab27ao27Key]; ab27ao27Present {
			if _, ab27ao27OK := ab27ao27Value.(string); !ab27ao27OK {
				return nil, fmt.Errorf("context %s must be a string", ab27ao27Key)
			}
		}
	}
	for _, ab27ao27Value := range ab27ao27Raw {
		ab27ao27Array, ab27ao27OK := ab27ao27Value.([]any)
		if !ab27ao27OK {
			return nil, fmt.Errorf("raw field must be a byte array")
		}
		for _, ab27ao27Value := range ab27ao27Array {
			ab27ao27Byte, ab27ao27OK := ab27ao27Integer(ab27ao27Value)
			if !ab27ao27OK || ab27ao27Byte < 0 || ab27ao27Byte > 255 {
				return nil, fmt.Errorf("invalid byte")
			}
		}
	}
	ab27ao27Fields := map[string]any{}
	ab27ao27Output := map[string]any{"context": ab27ao27Context, "raw": ab27ao27Raw, "state": "Unknown", "fields": ab27ao27Fields, "physical": "Unknown", "qualified": false}
	if ab27ao27Arguments, ab27ao27Present := ab27ao27Root["arguments"]; ab27ao27Present {
		ab27ao27Object, ab27ao27OK := ab27ao27Arguments.(map[string]any)
		if !ab27ao27OK || len(ab27ao27Object) != 0 {
			return nil, fmt.Errorf("no arguments are defined for this view")
		}
		ab27ao27Output["arguments"] = ab27ao27Object
	}
	ab27ao27Primary, _ := ab27ao27Integer(ab27ao27Context["primary"])
	ab27ao27Secondary, _ := ab27ao27Integer(ab27ao27Context["secondary"])
	ab27ao27Operation, _ := ab27ao27Integer(ab27ao27Context["operation"])
	ab27ao27View, _ := ab27ao27Context["view"].(string)
	if ab27ao27Primary != 181 || ab27ao27Secondary != 35 || ab27ao27Context["direction"] != "positive" || (ab27ao27View != "status-zero-valid" && ab27ao27View != "status-zero-valid-or-low") || (ab27ao27Operation != 3 && ab27ao27Operation != 6 && ab27ao27Operation != 7) {
		return json.Marshal(ab27ao27Output)
	}
	ab27ao27Object, ab27ao27Present := ab27ao27Raw["logical_object"].([]any)
	ab27ao27Output["state"] = "Hypothesis"
	ab27ao27Slice := func(ab27ao27Position, ab27ao27Width int) ([]any, bool) {
		if ab27ao27Position >= len(ab27ao27Object) {
			ab27ao27Output["state"] = "partial"
			return []any{}, false
		}
		ab27ao27End := ab27ao27Position + ab27ao27Width
		if ab27ao27End > len(ab27ao27Object) {
			ab27ao27Output["state"] = "partial"
			return ab27ao27Object[ab27ao27Position:], false
		}
		return ab27ao27Object[ab27ao27Position:ab27ao27End], true
	}
	ab27ao27Start, ab27ao27Count, ab27ao27StatusPosition := 1, 6, 16
	if ab27ao27Operation == 6 {
		ab27ao27Count, ab27ao27StatusPosition = 7, 17
	}
	if ab27ao27Operation == 7 {
		ab27ao27Start = 8
	}
	for ab27ao27Index := 0; ab27ao27Index < ab27ao27Count; ab27ao27Index++ {
		ab27ao27Position := 3 + 2*ab27ao27Index
		if ab27ao27Operation == 7 && ab27ao27Index >= 4 {
			ab27ao27Position++
		}
		ab27ao27Bytes, ab27ao27Complete := ab27ao27Slice(ab27ao27Position, 2)
		ab27ao27Field := map[string]any{"raw": ab27ao27Bytes, "state": "incomplete", "grade": "Hypothesis"}
		if ab27ao27Complete {
			ab27ao27Low, _ := ab27ao27Integer(ab27ao27Bytes[0])
			ab27ao27High, _ := ab27ao27Integer(ab27ao27Bytes[1])
			ab27ao27Word := ab27ao27Low | ab27ao27High<<8
			ab27ao27Field["raw"], ab27ao27Field["state"] = ab27ao27Word, "Hypothesis"
			if ab27ao27Word == 32768 {
				ab27ao27Field["label"] = "replacement"
			} else {
				ab27ao27Field["number"] = float64(int16(ab27ao27Word)) / 16
				ab27ao27Field["calculation"] = "Proven"
			}
		}
		ab27ao27Fields[fmt.Sprintf("S%d", ab27ao27Start+ab27ao27Index)] = ab27ao27Field
		ab27ao27Bytes, ab27ao27Complete = ab27ao27Slice(ab27ao27StatusPosition+ab27ao27Index/4, 1)
		ab27ao27Field = map[string]any{"raw": ab27ao27Bytes, "state": "incomplete", "grade": "Hypothesis"}
		if ab27ao27Complete {
			ab27ao27Byte, _ := ab27ao27Integer(ab27ao27Bytes[0])
			ab27ao27Pair := (ab27ao27Byte >> (2 * (ab27ao27Index % 4))) & 3
			ab27ao27Labels := []string{"valid-value", "open-circuit", "closed-circuit", "high-level"}
			if ab27ao27View == "status-zero-valid-or-low" {
				ab27ao27Labels[0] = "valid-value-or-low-level"
			}
			ab27ao27Field["raw"], ab27ao27Field["state"], ab27ao27Field["label"] = ab27ao27Pair, "Hypothesis", ab27ao27Labels[ab27ao27Pair]
		}
		ab27ao27Fields[fmt.Sprintf("statusS%d", ab27ao27Start+ab27ao27Index)] = ab27ao27Field
	}
	ab27ao27Bytes, ab27ao27Complete := ab27ao27Slice(ab27ao27StatusPosition+1, 1)
	ab27ao27Reserved := map[string]any{"raw": ab27ao27Bytes, "state": "incomplete", "grade": "Hypothesis"}
	if ab27ao27Complete {
		ab27ao27Byte, _ := ab27ao27Integer(ab27ao27Bytes[0])
		ab27ao27Mask := 240
		if ab27ao27Operation == 6 {
			ab27ao27Mask = 192
		}
		ab27ao27Reserved["raw"], ab27ao27Reserved["reserved"], ab27ao27Reserved["state"] = ab27ao27Byte, ab27ao27Byte&ab27ao27Mask, "Hypothesis"
	}
	ab27ao27Fields["status_reserved"] = ab27ao27Reserved
	if ab27ao27Operation == 3 || ab27ao27Operation == 7 {
		ab27ao27Position, ab27ao27Name := 15, "PWM"
		if ab27ao27Operation == 7 {
			ab27ao27Position, ab27ao27Name = 11, "opaque11"
		}
		ab27ao27Bytes, ab27ao27Complete := ab27ao27Slice(ab27ao27Position, 1)
		ab27ao27Field := map[string]any{"raw": ab27ao27Bytes, "state": "incomplete", "grade": "Hypothesis"}
		if ab27ao27Complete {
			ab27ao27Byte, _ := ab27ao27Integer(ab27ao27Bytes[0])
			ab27ao27Field["raw"] = ab27ao27Byte
			if ab27ao27Operation == 3 {
				ab27ao27Field["state"], ab27ao27Field["number"], ab27ao27Field["calculation"] = "Hypothesis", float64(ab27ao27Byte)/2, "Proven"
			} else {
				ab27ao27Field["state"] = "Unknown"
			}
		}
		ab27ao27Fields[ab27ao27Name] = ab27ao27Field
	}
	if !ab27ao27Present {
		ab27ao27Output["state"] = "partial"
	}
	return json.Marshal(ab27ao27Output)
}
