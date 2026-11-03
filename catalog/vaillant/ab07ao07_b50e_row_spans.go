package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

func ab07ao07Read(ab07ao07Decoder *json.Decoder) (any, error) {
	ab07ao07Token, ab07ao07Err := ab07ao07Decoder.Token()
	if ab07ao07Err != nil {
		return nil, ab07ao07Err
	}
	switch ab07ao07Token {
	case json.Delim('{'):
		ab07ao07Object := map[string]any{}
		for ab07ao07Decoder.More() {
			ab07ao07Key, ab07ao07Err := ab07ao07Decoder.Token()
			if ab07ao07Err != nil {
				return nil, ab07ao07Err
			}
			ab07ao07Name := ab07ao07Key.(string)
			if _, ab07ao07Exists := ab07ao07Object[ab07ao07Name]; ab07ao07Exists {
				return nil, fmt.Errorf("duplicate key %q", ab07ao07Name)
			}
			ab07ao07Value, ab07ao07Err := ab07ao07Read(ab07ao07Decoder)
			if ab07ao07Err != nil {
				return nil, ab07ao07Err
			}
			ab07ao07Object[ab07ao07Name] = ab07ao07Value
		}
		_, ab07ao07Err := ab07ao07Decoder.Token()
		return ab07ao07Object, ab07ao07Err
	case json.Delim('['):
		ab07ao07Array := []any{}
		for ab07ao07Decoder.More() {
			ab07ao07Value, ab07ao07Err := ab07ao07Read(ab07ao07Decoder)
			if ab07ao07Err != nil {
				return nil, ab07ao07Err
			}
			ab07ao07Array = append(ab07ao07Array, ab07ao07Value)
		}
		_, ab07ao07Err := ab07ao07Decoder.Token()
		return ab07ao07Array, ab07ao07Err
	}
	return ab07ao07Token, nil
}

func ab07ao07Integer(ab07ao07Value any) (int64, bool) {
	ab07ao07Number, ab07ao07OK := ab07ao07Value.(json.Number)
	if !ab07ao07OK {
		return 0, false
	}
	ab07ao07Int, ab07ao07Err := strconv.ParseInt(string(ab07ao07Number), 10, 64)
	return ab07ao07Int, ab07ao07Err == nil
}

func Ab07Ao07Interpret(ab07ao07Input json.RawMessage) (json.RawMessage, error) {
	if !utf8.Valid(ab07ao07Input) {
		return nil, fmt.Errorf("invalid UTF-8 JSON")
	}
	ab07ao07Decoder := json.NewDecoder(bytes.NewReader(ab07ao07Input))
	ab07ao07Decoder.UseNumber()
	ab07ao07Parsed, ab07ao07Err := ab07ao07Read(ab07ao07Decoder)
	if ab07ao07Err != nil {
		return nil, ab07ao07Err
	}
	if _, ab07ao07Err = ab07ao07Decoder.Token(); ab07ao07Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab07ao07Object, ab07ao07OK := ab07ao07Parsed.(map[string]any)
	if !ab07ao07OK {
		return nil, fmt.Errorf("input must be object")
	}
	for ab07ao07Name := range ab07ao07Object {
		if ab07ao07Name != "context" && ab07ao07Name != "fields" && ab07ao07Name != "arguments" {
			return nil, fmt.Errorf("unknown input member")
		}
	}
	ab07ao07Context, ab07ao07OK := ab07ao07Object["context"].(map[string]any)
	if !ab07ao07OK {
		return nil, fmt.Errorf("context must be object")
	}
	ab07ao07Raw, ab07ao07OK := ab07ao07Object["fields"].(map[string]any)
	if !ab07ao07OK {
		return nil, fmt.Errorf("fields must be object")
	}
	for _, ab07ao07Value := range ab07ao07Raw {
		ab07ao07Array, ab07ao07OK := ab07ao07Value.([]any)
		if !ab07ao07OK {
			return nil, fmt.Errorf("field must be byte array")
		}
		for _, ab07ao07Byte := range ab07ao07Array {
			ab07ao07Int, ab07ao07OK := ab07ao07Integer(ab07ao07Byte)
			if !ab07ao07OK || ab07ao07Int < 0 || ab07ao07Int > 255 {
				return nil, fmt.Errorf("invalid byte")
			}
		}
	}
	for _, ab07ao07Name := range []string{"primary", "secondary", "operation"} {
		if ab07ao07Value, ab07ao07Exists := ab07ao07Context[ab07ao07Name]; ab07ao07Exists {
			if _, ab07ao07OK := ab07ao07Integer(ab07ao07Value); !ab07ao07OK {
				return nil, fmt.Errorf("invalid context integer")
			}
		}
	}
	for _, ab07ao07Name := range []string{"direction", "view"} {
		if ab07ao07Value, ab07ao07Exists := ab07ao07Context[ab07ao07Name]; ab07ao07Exists {
			if _, ab07ao07OK := ab07ao07Value.(string); !ab07ao07OK {
				return nil, fmt.Errorf("invalid context token")
			}
		}
	}
	ab07ao07Fields := map[string]any{}
	ab07ao07Output := map[string]any{"context": ab07ao07Context, "raw": ab07ao07Raw, "state": "Unknown", "fields": ab07ao07Fields, "physical": "Unknown", "qualified": false}
	if ab07ao07Arguments, ab07ao07Exists := ab07ao07Object["arguments"]; ab07ao07Exists {
		ab07ao07Map, ab07ao07OK := ab07ao07Arguments.(map[string]any)
		if !ab07ao07OK || len(ab07ao07Map) != 0 {
			return nil, fmt.Errorf("no arguments defined for row-spans")
		}
		ab07ao07Output["arguments"] = ab07ao07Map
	}
	ab07ao07Primary, _ := ab07ao07Integer(ab07ao07Context["primary"])
	ab07ao07Secondary, _ := ab07ao07Integer(ab07ao07Context["secondary"])
	ab07ao07Operation, _ := ab07ao07Integer(ab07ao07Context["operation"])
	if ab07ao07Primary != 181 || ab07ao07Secondary != 14 || ab07ao07Operation != 2 || ab07ao07Context["direction"] != "positive" || ab07ao07Context["view"] != "row-spans" {
		return json.Marshal(ab07ao07Output)
	}
	ab07ao07Keys, _ := ab07ao07Raw["key"].([]any)
	if len(ab07ao07Keys) != 1 {
		return json.Marshal(ab07ao07Output)
	}
	ab07ao07Key, _ := ab07ao07Integer(ab07ao07Keys[0])
	if ab07ao07Key > 18 {
		return json.Marshal(ab07ao07Output)
	}
	ab07ao07Value, ab07ao07Present := ab07ao07Raw["value"].([]any)
	if ab07ao07Key == 0 || ab07ao07Key == 15 || ab07ao07Key == 16 {
		ab07ao07Field := map[string]any{"state": "Unknown"}
		if ab07ao07Present {
			ab07ao07Field["raw"] = ab07ao07Value
		}
		ab07ao07Fields["value"] = ab07ao07Field
		return json.Marshal(ab07ao07Output)
	}
	ab07ao07Output["state"] = "Hypothesis"
	ab07ao07Width := 1
	switch ab07ao07Key {
	case 4, 5, 11, 12, 17, 18:
		ab07ao07Width = 2
	case 6, 7, 13, 14:
		ab07ao07Width = 4
	}
	ab07ao07Span := func(ab07ao07Name string, ab07ao07Offset, ab07ao07Size int) map[string]any {
		ab07ao07End := ab07ao07Offset + ab07ao07Size
		if ab07ao07End > len(ab07ao07Value) {
			ab07ao07End = len(ab07ao07Value)
		}
		if ab07ao07Offset > ab07ao07End {
			ab07ao07Offset = ab07ao07End
		}
		ab07ao07Field := map[string]any{"state": "Hypothesis"}
		if ab07ao07Present {
			ab07ao07Field["raw"] = ab07ao07Value[ab07ao07Offset:ab07ao07End]
		}
		if ab07ao07End-ab07ao07Offset < ab07ao07Size {
			ab07ao07Field["state"] = "incomplete"
			ab07ao07Output["state"] = "partial"
		}
		ab07ao07Fields[ab07ao07Name] = ab07ao07Field
		return ab07ao07Field
	}
	ab07ao07Field := ab07ao07Span("value", 0, ab07ao07Width)
	if ab07ao07Field["state"] == "Hypothesis" && ab07ao07Width == 1 {
		ab07ao07Number, _ := ab07ao07Integer(ab07ao07Value[0])
		switch ab07ao07Key {
		case 1:
			ab07ao07Field["state"] = "Unknown"
		case 8:
			switch ab07ao07Number {
			case 15:
				ab07ao07Field["label"] = "true"
			case 240:
				ab07ao07Field["label"] = "false"
			default:
				ab07ao07Field["state"] = "Unknown"
			}
		default:
			if (ab07ao07Key == 3 || ab07ao07Key == 10) && ab07ao07Number >= 128 {
				ab07ao07Number -= 256
			}
			ab07ao07Field["number"] = ab07ao07Number
			ab07ao07Field["calculation"] = "Proven"
		}
	}
	if ab07ao07Key >= 9 && ab07ao07Key <= 14 || ab07ao07Key == 18 {
		ab07ao07Span("inverse", ab07ao07Width, ab07ao07Width)
	}
	if ab07ao07Key == 17 {
		ab07ao07Status := ab07ao07Span("status", 2, 1)
		if ab07ao07Status["state"] == "Hypothesis" {
			ab07ao07Number, _ := ab07ao07Integer(ab07ao07Value[2])
			switch ab07ao07Number {
			case 0:
				ab07ao07Status["label"] = "okay"
			case 85:
				ab07ao07Status["label"] = "short-circuit"
			case 170:
				ab07ao07Status["label"] = "interrupted"
			default:
				ab07ao07Status["state"] = "Unknown"
			}
		}
	}
	return json.Marshal(ab07ao07Output)
}
