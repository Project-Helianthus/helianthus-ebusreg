package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func ab17ao17Read(ab17ao17Decoder *json.Decoder) (any, error) {
	ab17ao17Token, ab17ao17Error := ab17ao17Decoder.Token()
	if ab17ao17Error != nil {
		return nil, ab17ao17Error
	}
	switch ab17ao17Token {
	case json.Delim('{'):
		ab17ao17Object := map[string]any{}
		for ab17ao17Decoder.More() {
			ab17ao17KeyToken, ab17ao17Error := ab17ao17Decoder.Token()
			if ab17ao17Error != nil {
				return nil, ab17ao17Error
			}
			ab17ao17Key := ab17ao17KeyToken.(string)
			if _, ab17ao17Exists := ab17ao17Object[ab17ao17Key]; ab17ao17Exists {
				return nil, fmt.Errorf("duplicate key %q", ab17ao17Key)
			}
			ab17ao17Value, ab17ao17Error := ab17ao17Read(ab17ao17Decoder)
			if ab17ao17Error != nil {
				return nil, ab17ao17Error
			}
			ab17ao17Object[ab17ao17Key] = ab17ao17Value
		}
		_, ab17ao17Error = ab17ao17Decoder.Token()
		return ab17ao17Object, ab17ao17Error
	case json.Delim('['):
		ab17ao17Array := []any{}
		for ab17ao17Decoder.More() {
			ab17ao17Value, ab17ao17Error := ab17ao17Read(ab17ao17Decoder)
			if ab17ao17Error != nil {
				return nil, ab17ao17Error
			}
			ab17ao17Array = append(ab17ao17Array, ab17ao17Value)
		}
		_, ab17ao17Error = ab17ao17Decoder.Token()
		return ab17ao17Array, ab17ao17Error
	}
	return ab17ao17Token, nil
}

func ab17ao17Integer(ab17ao17Value any) (int64, bool) {
	ab17ao17Number, ab17ao17OK := ab17ao17Value.(json.Number)
	if !ab17ao17OK {
		return 0, false
	}
	ab17ao17Int, ab17ao17Error := strconv.ParseInt(string(ab17ao17Number), 10, 64)
	return ab17ao17Int, ab17ao17Error == nil
}

func Ab17Ao17Interpret(ab17ao17Input json.RawMessage) (json.RawMessage, error) {
	ab17ao17Decoder := json.NewDecoder(bytes.NewReader(ab17ao17Input))
	ab17ao17Decoder.UseNumber()
	ab17ao17Value, ab17ao17Error := ab17ao17Read(ab17ao17Decoder)
	if ab17ao17Error != nil {
		return nil, ab17ao17Error
	}
	if _, ab17ao17Error = ab17ao17Decoder.Token(); ab17ao17Error != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab17ao17Root, ab17ao17OK := ab17ao17Value.(map[string]any)
	if !ab17ao17OK {
		return nil, fmt.Errorf("input must be object")
	}
	for ab17ao17Key := range ab17ao17Root {
		if ab17ao17Key != "context" && ab17ao17Key != "fields" && ab17ao17Key != "arguments" {
			return nil, fmt.Errorf("unknown input key %q", ab17ao17Key)
		}
	}
	ab17ao17Context, ab17ao17OK := ab17ao17Root["context"].(map[string]any)
	if !ab17ao17OK {
		return nil, fmt.Errorf("context must be object")
	}
	ab17ao17Raw, ab17ao17OK := ab17ao17Root["fields"].(map[string]any)
	if !ab17ao17OK {
		return nil, fmt.Errorf("fields must be object")
	}
	for _, ab17ao17Key := range []string{"primary", "secondary", "operation"} {
		if ab17ao17Value, ab17ao17Exists := ab17ao17Context[ab17ao17Key]; ab17ao17Exists {
			if _, ab17ao17OK := ab17ao17Integer(ab17ao17Value); !ab17ao17OK {
				return nil, fmt.Errorf("%s must be integer", ab17ao17Key)
			}
		}
	}
	for _, ab17ao17Key := range []string{"direction", "view"} {
		if ab17ao17Value, ab17ao17Exists := ab17ao17Context[ab17ao17Key]; ab17ao17Exists {
			if _, ab17ao17OK := ab17ao17Value.(string); !ab17ao17OK {
				return nil, fmt.Errorf("%s must be string", ab17ao17Key)
			}
		}
	}
	for ab17ao17Key, ab17ao17Value := range ab17ao17Raw {
		ab17ao17Array, ab17ao17OK := ab17ao17Value.([]any)
		if !ab17ao17OK {
			return nil, fmt.Errorf("%s must be byte array", ab17ao17Key)
		}
		for _, ab17ao17Byte := range ab17ao17Array {
			ab17ao17Int, ab17ao17OK := ab17ao17Integer(ab17ao17Byte)
			if !ab17ao17OK || ab17ao17Int < 0 || ab17ao17Int > 255 {
				return nil, fmt.Errorf("invalid byte in %s", ab17ao17Key)
			}
		}
	}
	ab17ao17Fields := map[string]any{}
	ab17ao17Output := map[string]any{"context": ab17ao17Context, "raw": ab17ao17Raw, "state": "Unknown", "fields": ab17ao17Fields, "physical": "Unknown", "qualified": false}
	if ab17ao17Arguments, ab17ao17Exists := ab17ao17Root["arguments"]; ab17ao17Exists {
		ab17ao17Object, ab17ao17OK := ab17ao17Arguments.(map[string]any)
		if !ab17ao17OK || len(ab17ao17Object) != 0 {
			return nil, fmt.Errorf("no arguments defined for these views")
		}
		ab17ao17Output["arguments"] = ab17ao17Object
	}
	ab17ao17Primary, _ := ab17ao17Integer(ab17ao17Context["primary"])
	ab17ao17Secondary, _ := ab17ao17Integer(ab17ao17Context["secondary"])
	ab17ao17Operation, _ := ab17ao17Integer(ab17ao17Context["operation"])
	if ab17ao17Primary != 181 || ab17ao17Secondary != 19 || ab17ao17Context["direction"] != "positive" {
		return json.Marshal(ab17ao17Output)
	}
	if ab17ao17Operation == 7 && ab17ao17Context["view"] == "information" {
		ab17ao17Object, _ := ab17ao17Raw["logical_object"].([]any)
		ab17ao17Output["state"] = "Hypothesis"
		if len(ab17ao17Object) < 12 {
			ab17ao17Output["state"] = "partial"
		}
		for ab17ao17Index, ab17ao17Name := range []string{"type", "capabilities", "fuel", "storage", "features", "pump_mode"} {
			ab17ao17Position := ab17ao17Index + 3
			ab17ao17Field := map[string]any{"raw": []any{}, "state": "Unknown"}
			ab17ao17Fields[ab17ao17Name] = ab17ao17Field
			if len(ab17ao17Object) <= ab17ao17Position {
				continue
			}
			ab17ao17Field["raw"] = ab17ao17Object[ab17ao17Position : ab17ao17Position+1]
			ab17ao17Number, _ := ab17ao17Integer(ab17ao17Object[ab17ao17Position])
			if ab17ao17Index == 0 || ab17ao17Index == 5 {
				ab17ao17Labels := []string{"modulating", "single-step", "two-step"}
				if ab17ao17Index == 5 {
					ab17ao17Labels = []string{"pump-mode0", "pump-mode1"}
				}
				if ab17ao17Number >= int64(len(ab17ao17Labels)) {
					ab17ao17Output["state"] = "partial"
					continue
				}
				ab17ao17Field["label"] = ab17ao17Labels[ab17ao17Number]
			} else {
				ab17ao17Labels := [][]string{{"floor-standing", "condensing", "fan-assisted", "open-vented", "high-modulation-range", "water-heater", "step-modulation", "ECS"}, {"gas", "oil", "ABC"}, {"layer-tank", "warm-start", "external-tank", "combination"}, {"solar-compact", "basic-solar", "AAPS", "advanced-power", "partial-load", "parameter-master"}}[ab17ao17Index-1]
				ab17ao17Flags := map[string]bool{}
				for ab17ao17Bit, ab17ao17Label := range ab17ao17Labels {
					ab17ao17Flags[ab17ao17Label] = ab17ao17Number&(1<<ab17ao17Bit) != 0
				}
				ab17ao17Field["flags"] = ab17ao17Flags
				ab17ao17Field["reserved"] = ab17ao17Number & (255 ^ ((1 << len(ab17ao17Labels)) - 1))
			}
			ab17ao17Field["number"] = ab17ao17Number
			ab17ao17Field["state"] = "Hypothesis"
		}
		ab17ao17Opaque := []any{}
		if len(ab17ao17Object) > 9 {
			ab17ao17End := len(ab17ao17Object)
			if ab17ao17End > 12 {
				ab17ao17End = 12
			}
			ab17ao17Opaque = ab17ao17Object[9:ab17ao17End]
		}
		ab17ao17Fields["opaque"] = map[string]any{"raw": ab17ao17Opaque, "state": "Unknown"}
	} else if (ab17ao17Operation == 5 || ab17ao17Operation == 6) && ab17ao17Context["view"] == "feedback" {
		ab17ao17Feedback, _ := ab17ao17Raw["feedback"].([]any)
		ab17ao17Field := map[string]any{"raw": ab17ao17Feedback, "state": "Unknown"}
		if ab17ao17Feedback == nil {
			ab17ao17Field["raw"] = []any{}
		}
		ab17ao17Fields["feedback"] = ab17ao17Field
		ab17ao17Output["state"] = "partial"
		if len(ab17ao17Feedback) == 1 {
			ab17ao17Number, _ := ab17ao17Integer(ab17ao17Feedback[0])
			ab17ao17Label := "rejected"
			if ab17ao17Operation == 5 {
				ab17ao17Label = "not-accepted-or-not-supported"
			}
			if ab17ao17Number == 1 {
				ab17ao17Label = "accepted"
			}
			ab17ao17Field["number"] = ab17ao17Number
			ab17ao17Field["label"] = ab17ao17Label
			ab17ao17Field["state"] = "Hypothesis"
			ab17ao17Output["state"] = "Hypothesis"
		}
	}
	return json.Marshal(ab17ao17Output)
}
