package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ab03ao03Read(ab03ao03Decoder *json.Decoder) (any, error) {
	ab03ao03Token, ab03ao03Err := ab03ao03Decoder.Token()
	if ab03ao03Err != nil {
		return nil, ab03ao03Err
	}
	if ab03ao03Delimiter, ab03ao03OK := ab03ao03Token.(json.Delim); ab03ao03OK {
		switch ab03ao03Delimiter {
		case '{':
			ab03ao03Object := map[string]any{}
			for ab03ao03Decoder.More() {
				ab03ao03KeyToken, ab03ao03Err := ab03ao03Decoder.Token()
				if ab03ao03Err != nil {
					return nil, ab03ao03Err
				}
				ab03ao03Key, ab03ao03OK := ab03ao03KeyToken.(string)
				if !ab03ao03OK {
					return nil, fmt.Errorf("invalid object key")
				}
				if _, ab03ao03Exists := ab03ao03Object[ab03ao03Key]; ab03ao03Exists {
					return nil, fmt.Errorf("duplicate key %q", ab03ao03Key)
				}
				ab03ao03Value, ab03ao03Err := ab03ao03Read(ab03ao03Decoder)
				if ab03ao03Err != nil {
					return nil, ab03ao03Err
				}
				ab03ao03Object[ab03ao03Key] = ab03ao03Value
			}
			_, ab03ao03Err := ab03ao03Decoder.Token()
			return ab03ao03Object, ab03ao03Err
		case '[':
			ab03ao03Array := []any{}
			for ab03ao03Decoder.More() {
				ab03ao03Value, ab03ao03Err := ab03ao03Read(ab03ao03Decoder)
				if ab03ao03Err != nil {
					return nil, ab03ao03Err
				}
				ab03ao03Array = append(ab03ao03Array, ab03ao03Value)
			}
			_, ab03ao03Err := ab03ao03Decoder.Token()
			return ab03ao03Array, ab03ao03Err
		default:
			return nil, fmt.Errorf("unexpected delimiter")
		}
	}
	return ab03ao03Token, nil
}

func ab03ao03Integer(ab03ao03Value any) (int64, bool) {
	ab03ao03Number, ab03ao03OK := ab03ao03Value.(json.Number)
	if !ab03ao03OK {
		return 0, false
	}
	ab03ao03Int, ab03ao03Err := strconv.ParseInt(string(ab03ao03Number), 10, 64)
	return ab03ao03Int, ab03ao03Err == nil
}

func Ab03Ao03Interpret(ab03ao03Input json.RawMessage) (json.RawMessage, error) {
	ab03ao03Decoder := json.NewDecoder(bytes.NewReader(ab03ao03Input))
	ab03ao03Decoder.UseNumber()
	ab03ao03Value, ab03ao03Err := ab03ao03Read(ab03ao03Decoder)
	if ab03ao03Err != nil {
		return nil, ab03ao03Err
	}
	if _, ab03ao03Err = ab03ao03Decoder.Token(); ab03ao03Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab03ao03Object, ab03ao03OK := ab03ao03Value.(map[string]any)
	if !ab03ao03OK {
		return nil, fmt.Errorf("input must be an object")
	}
	for ab03ao03Key := range ab03ao03Object {
		if ab03ao03Key != "context" && ab03ao03Key != "fields" && ab03ao03Key != "arguments" {
			return nil, fmt.Errorf("unknown input member")
		}
	}
	ab03ao03Context, ab03ao03OK := ab03ao03Object["context"].(map[string]any)
	if !ab03ao03OK {
		return nil, fmt.Errorf("context must be an object")
	}
	for _, ab03ao03Key := range []string{"primary", "secondary", "operation"} {
		if ab03ao03Value, ab03ao03Present := ab03ao03Context[ab03ao03Key]; ab03ao03Present {
			ab03ao03Number, ab03ao03OK := ab03ao03Value.(json.Number)
			if !ab03ao03OK || strings.ContainsAny(string(ab03ao03Number), ".eE") {
				return nil, fmt.Errorf("%s must be an integer", ab03ao03Key)
			}
		}
	}
	for _, ab03ao03Key := range []string{"direction", "view"} {
		if ab03ao03Value, ab03ao03Present := ab03ao03Context[ab03ao03Key]; ab03ao03Present {
			if _, ab03ao03OK := ab03ao03Value.(string); !ab03ao03OK {
				return nil, fmt.Errorf("%s must be a string", ab03ao03Key)
			}
		}
	}
	ab03ao03Raw, ab03ao03OK := ab03ao03Object["fields"].(map[string]any)
	if !ab03ao03OK {
		return nil, fmt.Errorf("fields must be an object")
	}
	for _, ab03ao03Value := range ab03ao03Raw {
		ab03ao03Array, ab03ao03OK := ab03ao03Value.([]any)
		if !ab03ao03OK {
			return nil, fmt.Errorf("field must be a byte array")
		}
		for _, ab03ao03Byte := range ab03ao03Array {
			ab03ao03Int, ab03ao03OK := ab03ao03Integer(ab03ao03Byte)
			if !ab03ao03OK || ab03ao03Int < 0 || ab03ao03Int > 255 {
				return nil, fmt.Errorf("invalid byte")
			}
		}
	}
	ab03ao03Fields := map[string]any{}
	ab03ao03Output := map[string]any{"context": ab03ao03Context, "raw": ab03ao03Raw, "state": "Unknown", "fields": ab03ao03Fields, "physical": "Unknown", "qualified": false}
	if ab03ao03Arguments, ab03ao03Present := ab03ao03Object["arguments"]; ab03ao03Present {
		ab03ao03Map, ab03ao03OK := ab03ao03Arguments.(map[string]any)
		if !ab03ao03OK || len(ab03ao03Map) != 0 {
			return nil, fmt.Errorf("no arguments defined for this view")
		}
		ab03ao03Output["arguments"] = ab03ao03Map
	}
	ab03ao03Primary, _ := ab03ao03Integer(ab03ao03Context["primary"])
	ab03ao03Secondary, _ := ab03ao03Integer(ab03ao03Context["secondary"])
	ab03ao03Operation, _ := ab03ao03Integer(ab03ao03Context["operation"])
	if ab03ao03Primary != 181 || ab03ao03Secondary != 6 || ab03ao03Operation != 17 || ab03ao03Context["direction"] != "request" || ab03ao03Context["view"] != "parameter-union" {
		return json.Marshal(ab03ao03Output)
	}
	ab03ao03Key, ab03ao03Present := ab03ao03Raw["key"].([]any)
	if !ab03ao03Present || len(ab03ao03Key) != 1 {
		ab03ao03Output["state"] = "partial"
		ab03ao03State := "incomplete"
		if len(ab03ao03Key) > 1 {
			ab03ao03State = "width-conflict"
		}
		ab03ao03Field := map[string]any{"state": ab03ao03State}
		if ab03ao03Present {
			ab03ao03Field["raw"] = ab03ao03Key
		}
		ab03ao03Fields["key"] = ab03ao03Field
		return json.Marshal(ab03ao03Output)
	}
	ab03ao03KeyNumber, _ := ab03ao03Integer(ab03ao03Key[0])
	ab03ao03Opaque := false
	switch ab03ao03KeyNumber {
	case 0, 1, 2, 4, 5:
		ab03ao03Opaque = true
	case 3, 6, 7, 8, 9, 10, 32, 33, 34, 35, 36, 37:
	default:
		return json.Marshal(ab03ao03Output)
	}
	ab03ao03Fields["key"] = map[string]any{"raw": ab03ao03Key, "state": "Hypothesis", "number": ab03ao03KeyNumber}
	ab03ao03Field := map[string]any{"state": "incomplete"}
	ab03ao03Fields["value"] = ab03ao03Field
	ab03ao03Output["state"] = "partial"
	ab03ao03Bytes, ab03ao03Present := ab03ao03Raw["value"].([]any)
	if !ab03ao03Present {
		return json.Marshal(ab03ao03Output)
	}
	ab03ao03Field["raw"] = ab03ao03Bytes
	if ab03ao03Opaque {
		ab03ao03Field["state"] = "Hypothesis"
		ab03ao03Output["state"] = "Hypothesis"
		return json.Marshal(ab03ao03Output)
	}
	if len(ab03ao03Bytes) == 0 {
		return json.Marshal(ab03ao03Output)
	}
	if len(ab03ao03Bytes) != 1 {
		ab03ao03Field["state"] = "width-conflict"
		return json.Marshal(ab03ao03Output)
	}
	ab03ao03Number, _ := ab03ao03Integer(ab03ao03Bytes[0])
	if ab03ao03KeyNumber == 10 {
		if ab03ao03Number > 1 {
			ab03ao03Field["state"] = "unmapped"
			return json.Marshal(ab03ao03Output)
		}
		if ab03ao03Number == 0 {
			ab03ao03Field["label"] = "no"
		} else {
			ab03ao03Field["label"] = "yes"
		}
	}
	ab03ao03Field["state"] = "Hypothesis"
	ab03ao03Field["number"] = ab03ao03Number
	ab03ao03Output["state"] = "Hypothesis"
	return json.Marshal(ab03ao03Output)
}
