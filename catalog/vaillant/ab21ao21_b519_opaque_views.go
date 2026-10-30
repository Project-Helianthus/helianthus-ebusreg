package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func ab21ao21Value(ab21ao21Decoder *json.Decoder) (any, error) {
	ab21ao21Token, ab21ao21Err := ab21ao21Decoder.Token()
	if ab21ao21Err != nil {
		return nil, ab21ao21Err
	}
	switch ab21ao21Token {
	case json.Delim('{'):
		ab21ao21Object := map[string]any{}
		for ab21ao21Decoder.More() {
			ab21ao21KeyToken, ab21ao21Err := ab21ao21Decoder.Token()
			if ab21ao21Err != nil {
				return nil, ab21ao21Err
			}
			ab21ao21Key := ab21ao21KeyToken.(string)
			if _, ab21ao21Found := ab21ao21Object[ab21ao21Key]; ab21ao21Found {
				return nil, fmt.Errorf("duplicate key %q", ab21ao21Key)
			}
			ab21ao21Object[ab21ao21Key], ab21ao21Err = ab21ao21Value(ab21ao21Decoder)
			if ab21ao21Err != nil {
				return nil, ab21ao21Err
			}
		}
		_, ab21ao21Err := ab21ao21Decoder.Token()
		return ab21ao21Object, ab21ao21Err
	case json.Delim('['):
		ab21ao21Array := []any{}
		for ab21ao21Decoder.More() {
			ab21ao21Item, ab21ao21Err := ab21ao21Value(ab21ao21Decoder)
			if ab21ao21Err != nil {
				return nil, ab21ao21Err
			}
			ab21ao21Array = append(ab21ao21Array, ab21ao21Item)
		}
		_, ab21ao21Err := ab21ao21Decoder.Token()
		return ab21ao21Array, ab21ao21Err
	}
	return ab21ao21Token, nil
}

func ab21ao21Integer(ab21ao21Value any) (int, bool) {
	ab21ao21Number, ab21ao21OK := ab21ao21Value.(json.Number)
	if !ab21ao21OK {
		return 0, false
	}
	for ab21ao21Index, ab21ao21Digit := range string(ab21ao21Number) {
		if ab21ao21Index == 0 && ab21ao21Digit == '-' {
			continue
		}
		if ab21ao21Digit < '0' || ab21ao21Digit > '9' {
			return 0, false
		}
	}
	ab21ao21Int, ab21ao21Err := strconv.Atoi(string(ab21ao21Number))
	if ab21ao21Err != nil {
		return -1, true
	}
	return ab21ao21Int, true
}

func Ab21Ao21Interpret(ab21ao21Input json.RawMessage) (json.RawMessage, error) {
	ab21ao21Decoder := json.NewDecoder(bytes.NewReader(ab21ao21Input))
	ab21ao21Decoder.UseNumber()
	ab21ao21Parsed, ab21ao21Err := ab21ao21Value(ab21ao21Decoder)
	if ab21ao21Err != nil {
		return nil, ab21ao21Err
	}
	if _, ab21ao21Err = ab21ao21Decoder.Token(); ab21ao21Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab21ao21Root, ab21ao21OK := ab21ao21Parsed.(map[string]any)
	if !ab21ao21OK {
		return nil, fmt.Errorf("input must be object")
	}
	for ab21ao21Key := range ab21ao21Root {
		if ab21ao21Key != "context" && ab21ao21Key != "fields" && ab21ao21Key != "arguments" {
			return nil, fmt.Errorf("unknown input key %q", ab21ao21Key)
		}
	}
	ab21ao21Context, ab21ao21OK := ab21ao21Root["context"].(map[string]any)
	if !ab21ao21OK {
		return nil, fmt.Errorf("context must be object")
	}
	ab21ao21Raw, ab21ao21OK := ab21ao21Root["fields"].(map[string]any)
	if !ab21ao21OK {
		return nil, fmt.Errorf("fields must be object")
	}
	for ab21ao21Key, ab21ao21Value := range ab21ao21Raw {
		ab21ao21Array, ab21ao21OK := ab21ao21Value.([]any)
		if !ab21ao21OK {
			return nil, fmt.Errorf("field %q must be byte array", ab21ao21Key)
		}
		for _, ab21ao21Byte := range ab21ao21Array {
			ab21ao21Int, ab21ao21OK := ab21ao21Integer(ab21ao21Byte)
			if !ab21ao21OK || ab21ao21Int < 0 || ab21ao21Int > 255 {
				return nil, fmt.Errorf("invalid byte in %q", ab21ao21Key)
			}
		}
	}
	for _, ab21ao21Key := range []string{"primary", "secondary", "operation"} {
		if ab21ao21Value, ab21ao21Present := ab21ao21Context[ab21ao21Key]; ab21ao21Present {
			if _, ab21ao21OK := ab21ao21Integer(ab21ao21Value); !ab21ao21OK {
				return nil, fmt.Errorf("%s must be integer", ab21ao21Key)
			}
		}
	}
	for _, ab21ao21Key := range []string{"direction", "view"} {
		if ab21ao21Value, ab21ao21Present := ab21ao21Context[ab21ao21Key]; ab21ao21Present {
			if _, ab21ao21OK := ab21ao21Value.(string); !ab21ao21OK {
				return nil, fmt.Errorf("%s must be string", ab21ao21Key)
			}
		}
	}
	ab21ao21Fields := map[string]any{}
	ab21ao21Output := map[string]any{"context": ab21ao21Context, "raw": ab21ao21Raw, "state": "Unknown", "fields": ab21ao21Fields, "physical": "Unknown", "qualified": false}
	if ab21ao21Arguments, ab21ao21Present := ab21ao21Root["arguments"]; ab21ao21Present {
		ab21ao21Object, ab21ao21OK := ab21ao21Arguments.(map[string]any)
		if !ab21ao21OK || len(ab21ao21Object) != 0 {
			return nil, fmt.Errorf("no arguments declared")
		}
		ab21ao21Output["arguments"] = ab21ao21Object
	}
	ab21ao21Primary, _ := ab21ao21Integer(ab21ao21Context["primary"])
	ab21ao21Secondary, _ := ab21ao21Integer(ab21ao21Context["secondary"])
	ab21ao21Operation, ab21ao21HasOperation := ab21ao21Integer(ab21ao21Context["operation"])
	ab21ao21Direction := ab21ao21Context["direction"]
	if ab21ao21Primary != 181 || ab21ao21Secondary != 25 || !ab21ao21HasOperation || ((ab21ao21Operation < 0 || ab21ao21Operation > 8) && ab21ao21Operation != 170) || ab21ao21Context["view"] != "logical-body" || (ab21ao21Direction != "request" && ab21ao21Direction != "positive") {
		return json.Marshal(ab21ao21Output)
	}
	if ab21ao21Direction == "positive" && ab21ao21Operation != 3 && ab21ao21Operation != 6 && ab21ao21Operation != 170 {
		return json.Marshal(ab21ao21Output)
	}
	ab21ao21Logical, ab21ao21Present := ab21ao21Raw["logical_object"].([]any)
	if !ab21ao21Present || len(ab21ao21Logical) < 3 {
		ab21ao21Output["state"] = "partial"
		if ab21ao21Present {
			ab21ao21Fields["logical_prefix"] = map[string]any{"raw": ab21ao21Logical, "state": "Unknown"}
		}
		return json.Marshal(ab21ao21Output)
	}
	ab21ao21Output["state"] = "Hypothesis"
	if ab21ao21Direction == "positive" && ab21ao21Operation == 6 || ab21ao21Direction == "request" && (ab21ao21Operation == 0 || ab21ao21Operation == 2 || ab21ao21Operation == 170) {
		ab21ao21Label := "no-declared-argument"
		if ab21ao21Direction == "positive" {
			ab21ao21Label = "prefix-only"
		}
		ab21ao21Fields["form"] = map[string]any{"raw": ab21ao21Logical[:3], "state": "Hypothesis", "label": ab21ao21Label}
		if len(ab21ao21Logical) > 3 {
			ab21ao21Fields["unmodelled"] = map[string]any{"raw": ab21ao21Logical[3:], "state": "Unknown"}
			ab21ao21Output["state"] = "partial"
		}
		return json.Marshal(ab21ao21Output)
	}
	ab21ao21Offset := 3
	if ab21ao21Direction == "request" && ab21ao21Operation == 8 {
		if len(ab21ao21Logical) == 3 {
			ab21ao21Fields["slave_address"] = map[string]any{"raw": []any{}, "state": "Unknown"}
			ab21ao21Output["state"] = "partial"
			return json.Marshal(ab21ao21Output)
		}
		ab21ao21Fields["slave_address"] = map[string]any{"raw": ab21ao21Logical[3:4], "state": "Hypothesis"}
		ab21ao21Offset = 4
	}
	ab21ao21Fields["opaque"] = map[string]any{"raw": ab21ao21Logical[ab21ao21Offset:], "state": "Hypothesis", "flags": map[string]any{"accepted_empty": "Unknown", "unlimited_carrier": "Unknown"}}
	return json.Marshal(ab21ao21Output)
}
