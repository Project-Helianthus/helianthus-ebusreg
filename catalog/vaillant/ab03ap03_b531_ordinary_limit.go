package vaillant

import (
	ab03ap03bytes "bytes"
	ab03ap03json "encoding/json"
	ab03ap03fmt "fmt"
	ab03ap03io "io"
	ab03ap03strconv "strconv"
)

func ab03ap03Decode(ab03ap03Decoder *ab03ap03json.Decoder) (any, error) {
	ab03ap03Token, ab03ap03Error := ab03ap03Decoder.Token()
	if ab03ap03Error != nil {
		return nil, ab03ap03Error
	}
	if ab03ap03Delimiter, ab03ap03OK := ab03ap03Token.(ab03ap03json.Delim); ab03ap03OK {
		switch ab03ap03Delimiter {
		case '{':
			ab03ap03Object := map[string]any{}
			for ab03ap03Decoder.More() {
				ab03ap03KeyToken, ab03ap03Error := ab03ap03Decoder.Token()
				if ab03ap03Error != nil {
					return nil, ab03ap03Error
				}
				ab03ap03Key, ab03ap03OK := ab03ap03KeyToken.(string)
				if !ab03ap03OK {
					return nil, ab03ap03fmt.Errorf("invalid object key")
				}
				if _, ab03ap03Exists := ab03ap03Object[ab03ap03Key]; ab03ap03Exists {
					return nil, ab03ap03fmt.Errorf("duplicate key %q", ab03ap03Key)
				}
				ab03ap03Value, ab03ap03Error := ab03ap03Decode(ab03ap03Decoder)
				if ab03ap03Error != nil {
					return nil, ab03ap03Error
				}
				ab03ap03Object[ab03ap03Key] = ab03ap03Value
			}
			_, ab03ap03Error := ab03ap03Decoder.Token()
			return ab03ap03Object, ab03ap03Error
		case '[':
			ab03ap03Array := []any{}
			for ab03ap03Decoder.More() {
				ab03ap03Value, ab03ap03Error := ab03ap03Decode(ab03ap03Decoder)
				if ab03ap03Error != nil {
					return nil, ab03ap03Error
				}
				ab03ap03Array = append(ab03ap03Array, ab03ap03Value)
			}
			_, ab03ap03Error := ab03ap03Decoder.Token()
			return ab03ap03Array, ab03ap03Error
		}
		return nil, ab03ap03fmt.Errorf("unexpected delimiter")
	}
	return ab03ap03Token, nil
}

func ab03ap03Uint(ab03ap03Value any, ab03ap03Bits int) (uint64, error) {
	ab03ap03Number, ab03ap03OK := ab03ap03Value.(ab03ap03json.Number)
	if !ab03ap03OK {
		return 0, ab03ap03fmt.Errorf("unsigned integer required")
	}
	return ab03ap03strconv.ParseUint(string(ab03ap03Number), 10, ab03ap03Bits)
}

func Ab03Ap03Evaluate(ab03ap03Input ab03ap03json.RawMessage) (ab03ap03json.RawMessage, error) {
	ab03ap03Decoder := ab03ap03json.NewDecoder(ab03ap03bytes.NewReader(ab03ap03Input))
	ab03ap03Decoder.UseNumber()
	ab03ap03Value, ab03ap03Error := ab03ap03Decode(ab03ap03Decoder)
	if ab03ap03Error != nil {
		return nil, ab03ap03Error
	}
	if _, ab03ap03Error = ab03ap03Decoder.Token(); ab03ap03Error != ab03ap03io.EOF {
		return nil, ab03ap03fmt.Errorf("trailing JSON")
	}
	ab03ap03InputObject, ab03ap03OK := ab03ap03Value.(map[string]any)
	if !ab03ap03OK {
		return nil, ab03ap03fmt.Errorf("input object required")
	}
	ab03ap03Context, ab03ap03OK := ab03ap03InputObject["context"].(map[string]any)
	if !ab03ap03OK {
		return nil, ab03ap03fmt.Errorf("context object required")
	}
	ab03ap03Raw, ab03ap03OK := ab03ap03InputObject["fields"].(map[string]any)
	if !ab03ap03OK {
		return nil, ab03ap03fmt.Errorf("fields object required")
	}
	for ab03ap03Key, ab03ap03Value := range ab03ap03Raw {
		ab03ap03Bytes, ab03ap03OK := ab03ap03Value.([]any)
		if !ab03ap03OK {
			return nil, ab03ap03fmt.Errorf("field %s must be a byte array", ab03ap03Key)
		}
		for _, ab03ap03Byte := range ab03ap03Bytes {
			if _, ab03ap03Error := ab03ap03Uint(ab03ap03Byte, 8); ab03ap03Error != nil {
				return nil, ab03ap03Error
			}
		}
	}
	ab03ap03Semantic := map[string]any{}
	ab03ap03Output := map[string]any{"context": ab03ap03Context, "raw": ab03ap03Raw, "state": "Unknown", "fields": ab03ap03Semantic, "physical": "Unknown", "qualified": false}
	ab03ap03Arguments := map[string]any{}
	if ab03ap03ArgumentValue, ab03ap03Present := ab03ap03InputObject["arguments"]; ab03ap03Present {
		ab03ap03Arguments, ab03ap03OK = ab03ap03ArgumentValue.(map[string]any)
		if !ab03ap03OK {
			return nil, ab03ap03fmt.Errorf("arguments object required")
		}
		ab03ap03Output["arguments"] = ab03ap03Arguments
	}
	ab03ap03Supported := true
	for ab03ap03Key, ab03ap03Expected := range map[string]uint64{"primary": 181, "secondary": 49, "operation": 3} {
		ab03ap03Value, ab03ap03Present := ab03ap03Context[ab03ap03Key]
		if !ab03ap03Present || ab03ap03Value == nil {
			ab03ap03Supported = false
			continue
		}
		ab03ap03Number, ab03ap03Error := ab03ap03Uint(ab03ap03Value, 8)
		if ab03ap03Error != nil {
			return nil, ab03ap03Error
		}
		if ab03ap03Number != ab03ap03Expected {
			ab03ap03Supported = false
		}
	}
	for ab03ap03Key, ab03ap03Expected := range map[string]string{"direction": "local", "view": "ordinary-limit"} {
		ab03ap03Value, ab03ap03Present := ab03ap03Context[ab03ap03Key]
		if !ab03ap03Present || ab03ap03Value == nil {
			ab03ap03Supported = false
			continue
		}
		ab03ap03String, ab03ap03OK := ab03ap03Value.(string)
		if !ab03ap03OK {
			return nil, ab03ap03fmt.Errorf("context %s must be a string", ab03ap03Key)
		}
		if ab03ap03String != ab03ap03Expected {
			ab03ap03Supported = false
		}
	}
	if !ab03ap03Supported {
		if ab03ap03Error := ab03ap03ValidateOptional(ab03ap03Arguments); ab03ap03Error != nil {
			return nil, ab03ap03Error
		}
		return ab03ap03json.Marshal(ab03ap03Output)
	}
	ab03ap03Limited, ab03ap03Error := ab03ap03Uint(ab03ap03Arguments["limited"], 32)
	if ab03ap03Error != nil {
		return nil, ab03ap03Error
	}
	ab03ap03Rows, ab03ap03OK := ab03ap03Arguments["rows"].([]any)
	if !ab03ap03OK || len(ab03ap03Rows) != 9 {
		return nil, ab03ap03fmt.Errorf("rows requires nine slots")
	}
	ab03ap03Indices := []int{}
	ab03ap03Slots := []any{}
	var ab03ap03TokenSum, ab03ap03MaximumSum uint64
	ab03ap03Special := false
	for ab03ap03Index, ab03ap03Value := range ab03ap03Rows {
		ab03ap03Row, ab03ap03OK := ab03ap03Value.(map[string]any)
		if !ab03ap03OK {
			return nil, ab03ap03fmt.Errorf("slot object required")
		}
		ab03ap03Predicates := map[string]bool{}
		for _, ab03ap03Key := range []string{"registered", "capable", "failure"} {
			ab03ap03Predicate, ab03ap03OK := ab03ap03Row[ab03ap03Key].(bool)
			if !ab03ap03OK {
				return nil, ab03ap03fmt.Errorf("boolean %s required", ab03ap03Key)
			}
			ab03ap03Predicates[ab03ap03Key] = ab03ap03Predicate
		}
		ab03ap03Token, ab03ap03Error := ab03ap03Uint(ab03ap03Row["token"], 32)
		if ab03ap03Error != nil {
			return nil, ab03ap03Error
		}
		ab03ap03Maximum, ab03ap03Error := ab03ap03Uint(ab03ap03Row["maximum"], 32)
		if ab03ap03Error != nil {
			return nil, ab03ap03Error
		}
		ab03ap03Admitted := ab03ap03Predicates["registered"] && ab03ap03Predicates["capable"] && !ab03ap03Predicates["failure"]
		ab03ap03Slots = append(ab03ap03Slots, map[string]any{"slot": ab03ap03Index, "raw": ab03ap03Row, "admitted": ab03ap03Admitted, "state": "Hypothesis"})
		if !ab03ap03Admitted {
			continue
		}
		ab03ap03Indices = append(ab03ap03Indices, ab03ap03Index)
		if ab03ap03Index == 1 {
			ab03ap03Special = true
			continue
		}
		ab03ap03TokenSum += ab03ap03Token
		ab03ap03MaximumSum += ab03ap03Maximum
	}
	ab03ap03Semantic["slots"] = ab03ap03Slots
	ab03ap03Semantic["admitted_slots"] = map[string]any{"raw": ab03ap03Rows, "state": "Hypothesis", "indices": ab03ap03Indices}
	ab03ap03Semantic["token_sum"] = map[string]any{"raw": ab03ap03Rows, "state": "Hypothesis", "calculation": "Proven", "number": ab03ap03TokenSum}
	ab03ap03Semantic["maximum_sum"] = map[string]any{"raw": ab03ap03Rows, "state": "Hypothesis", "calculation": "Proven", "number": ab03ap03MaximumSum}
	ab03ap03Effective := ab03ap03TokenSum < ab03ap03MaximumSum && ab03ap03Limited != 0
	ab03ap03EffectiveResult := map[string]any{"raw": ab03ap03Rows, "state": "Unknown"}
	ab03ap03Requested := map[string]any{"state": "Unknown"}
	if ab03ap03Gate, ab03ap03Present := ab03ap03Raw["gate"]; ab03ap03Present {
		ab03ap03Requested["raw"] = ab03ap03Gate
	}
	ab03ap03Output["state"] = "complete"
	if ab03ap03Special {
		ab03ap03Output["state"] = "partial"
	} else {
		ab03ap03EffectiveResult["state"] = "Hypothesis"
		ab03ap03EffectiveResult["calculation"] = "Proven"
		ab03ap03EffectiveResult["value"] = ab03ap03Effective
	}
	ab03ap03Gate, ab03ap03OK := ab03ap03Raw["gate"].([]any)
	if !ab03ap03OK || len(ab03ap03Gate) != 1 {
		ab03ap03Output["state"] = "partial"
	} else {
		ab03ap03GateNumber, _ := ab03ap03Uint(ab03ap03Gate[0], 8)
		if ab03ap03GateNumber == 0 {
			ab03ap03Requested["state"] = "Hypothesis"
			ab03ap03Requested["number"] = 0
			ab03ap03Requested["calculation"] = "Proven"
		} else if !ab03ap03Special {
			ab03ap03Number := 1
			if ab03ap03Effective {
				ab03ap03Number = 2
			}
			ab03ap03Requested["state"] = "Hypothesis"
			ab03ap03Requested["number"] = ab03ap03Number
			ab03ap03Requested["calculation"] = "Proven"
		}
	}
	ab03ap03Semantic["effective"] = ab03ap03EffectiveResult
	ab03ap03Semantic["requested_state"] = ab03ap03Requested
	return ab03ap03json.Marshal(ab03ap03Output)
}

func ab03ap03ValidateOptional(ab03ap03Arguments map[string]any) error {
	if ab03ap03Value, ab03ap03Present := ab03ap03Arguments["limited"]; ab03ap03Present {
		if _, ab03ap03Error := ab03ap03Uint(ab03ap03Value, 32); ab03ap03Error != nil {
			return ab03ap03Error
		}
	}
	if ab03ap03Value, ab03ap03Present := ab03ap03Arguments["rows"]; ab03ap03Present {
		ab03ap03Rows, ab03ap03OK := ab03ap03Value.([]any)
		if !ab03ap03OK || len(ab03ap03Rows) != 9 {
			return ab03ap03fmt.Errorf("rows requires nine slots")
		}
		for _, ab03ap03Value := range ab03ap03Rows {
			ab03ap03Row, ab03ap03OK := ab03ap03Value.(map[string]any)
			if !ab03ap03OK {
				return ab03ap03fmt.Errorf("slot object required")
			}
			for _, ab03ap03Key := range []string{"registered", "capable", "failure"} {
				if _, ab03ap03OK := ab03ap03Row[ab03ap03Key].(bool); !ab03ap03OK {
					return ab03ap03fmt.Errorf("boolean %s required", ab03ap03Key)
				}
			}
			for _, ab03ap03Key := range []string{"token", "maximum"} {
				if _, ab03ap03Error := ab03ap03Uint(ab03ap03Row[ab03ap03Key], 32); ab03ap03Error != nil {
					return ab03ap03Error
				}
			}
		}
	}
	return nil
}
