package vaillant

import (
	ab05ap05bytes "bytes"
	ab05ap05json "encoding/json"
	ab05ap05fmt "fmt"
	ab05ap05io "io"
	ab05ap05strconv "strconv"
)

func ab05ap05Read(ab05ap05Decoder *ab05ap05json.Decoder) (any, error) {
	ab05ap05Token, ab05ap05Err := ab05ap05Decoder.Token()
	if ab05ap05Err != nil {
		return nil, ab05ap05Err
	}
	switch ab05ap05Token {
	case ab05ap05json.Delim('{'):
		ab05ap05Object := map[string]any{}
		for ab05ap05Decoder.More() {
			ab05ap05KeyToken, ab05ap05Err := ab05ap05Decoder.Token()
			if ab05ap05Err != nil {
				return nil, ab05ap05Err
			}
			ab05ap05Key := ab05ap05KeyToken.(string)
			if _, ab05ap05Exists := ab05ap05Object[ab05ap05Key]; ab05ap05Exists {
				return nil, ab05ap05fmt.Errorf("duplicate key %q", ab05ap05Key)
			}
			ab05ap05Value, ab05ap05Err := ab05ap05Read(ab05ap05Decoder)
			if ab05ap05Err != nil {
				return nil, ab05ap05Err
			}
			ab05ap05Object[ab05ap05Key] = ab05ap05Value
		}
		_, ab05ap05Err = ab05ap05Decoder.Token()
		return ab05ap05Object, ab05ap05Err
	case ab05ap05json.Delim('['):
		ab05ap05Array := []any{}
		for ab05ap05Decoder.More() {
			ab05ap05Value, ab05ap05Err := ab05ap05Read(ab05ap05Decoder)
			if ab05ap05Err != nil {
				return nil, ab05ap05Err
			}
			ab05ap05Array = append(ab05ap05Array, ab05ap05Value)
		}
		_, ab05ap05Err = ab05ap05Decoder.Token()
		return ab05ap05Array, ab05ap05Err
	}
	if _, ab05ap05Delimiter := ab05ap05Token.(ab05ap05json.Delim); ab05ap05Delimiter {
		return nil, ab05ap05fmt.Errorf("unexpected delimiter")
	}
	return ab05ap05Token, nil
}

func ab05ap05Uint(ab05ap05Value any, ab05ap05Bits int) (uint64, error) {
	ab05ap05Number, ab05ap05OK := ab05ap05Value.(ab05ap05json.Number)
	if !ab05ap05OK {
		return 0, ab05ap05fmt.Errorf("expected unsigned integer")
	}
	return ab05ap05strconv.ParseUint(string(ab05ap05Number), 10, ab05ap05Bits)
}

func ab05ap05Quantity(ab05ap05Number any) map[string]any {
	return map[string]any{"state": "Hypothesis", "calculation": "Proven", "number": ab05ap05Number}
}

func Ab05Ap05Evaluate(ab05ap05Input ab05ap05json.RawMessage) (ab05ap05json.RawMessage, error) {
	ab05ap05Decoder := ab05ap05json.NewDecoder(ab05ap05bytes.NewReader(ab05ap05Input))
	ab05ap05Decoder.UseNumber()
	ab05ap05Value, ab05ap05Err := ab05ap05Read(ab05ap05Decoder)
	if ab05ap05Err != nil {
		return nil, ab05ap05Err
	}
	if _, ab05ap05Err = ab05ap05Decoder.Token(); ab05ap05Err != ab05ap05io.EOF {
		return nil, ab05ap05fmt.Errorf("trailing JSON")
	}
	ab05ap05Root, ab05ap05OK := ab05ap05Value.(map[string]any)
	if !ab05ap05OK {
		return nil, ab05ap05fmt.Errorf("expected object")
	}
	ab05ap05Context, ab05ap05OK := ab05ap05Root["context"].(map[string]any)
	if !ab05ap05OK {
		return nil, ab05ap05fmt.Errorf("expected context object")
	}
	ab05ap05Raw, ab05ap05OK := ab05ap05Root["fields"].(map[string]any)
	if !ab05ap05OK {
		return nil, ab05ap05fmt.Errorf("expected fields object")
	}
	ab05ap05Fields := map[string]any{}
	for ab05ap05Key, ab05ap05Value := range ab05ap05Raw {
		ab05ap05Array, ab05ap05OK := ab05ap05Value.([]any)
		if !ab05ap05OK {
			return nil, ab05ap05fmt.Errorf("field %s: expected byte array", ab05ap05Key)
		}
		for _, ab05ap05Byte := range ab05ap05Array {
			if _, ab05ap05Err := ab05ap05Uint(ab05ap05Byte, 8); ab05ap05Err != nil {
				return nil, ab05ap05Err
			}
		}
		ab05ap05Fields[ab05ap05Key] = map[string]any{"raw": ab05ap05Value, "state": "Unknown"}
	}
	ab05ap05Output := map[string]any{"context": ab05ap05Context, "raw": ab05ap05Raw, "state": "Unknown", "fields": ab05ap05Fields, "physical": "Unknown", "qualified": false}
	ab05ap05Arguments, ab05ap05Supplied := ab05ap05Root["arguments"]
	if ab05ap05Supplied {
		if _, ab05ap05OK := ab05ap05Arguments.(map[string]any); !ab05ap05OK {
			return nil, ab05ap05fmt.Errorf("expected arguments object")
		}
		ab05ap05Output["arguments"] = ab05ap05Arguments
	}
	ab05ap05Supported := true
	for ab05ap05Key, ab05ap05Expected := range map[string]uint64{"primary": 181, "secondary": 22, "operation": 14} {
		ab05ap05Value, ab05ap05Present := ab05ap05Context[ab05ap05Key]
		if !ab05ap05Present || ab05ap05Value == nil {
			ab05ap05Supported = false
			continue
		}
		ab05ap05Number, ab05ap05Err := ab05ap05Uint(ab05ap05Value, 8)
		if ab05ap05Err != nil {
			return nil, ab05ap05Err
		}
		if ab05ap05Number != ab05ap05Expected {
			ab05ap05Supported = false
		}
	}
	for ab05ap05Key, ab05ap05Expected := range map[string]string{"direction": "local", "view": "cycle"} {
		ab05ap05Value := ab05ap05Context[ab05ap05Key]
		if ab05ap05Value == nil {
			ab05ap05Supported = false
			continue
		}
		ab05ap05String, ab05ap05OK := ab05ap05Value.(string)
		if !ab05ap05OK {
			return nil, ab05ap05fmt.Errorf("expected context string")
		}
		if ab05ap05String != ab05ap05Expected {
			ab05ap05Supported = false
		}
	}
	if !ab05ap05Supported {
		return ab05ap05json.Marshal(ab05ap05Output)
	}
	ab05ap05Args, ab05ap05OK := ab05ap05Arguments.(map[string]any)
	if !ab05ap05OK {
		return nil, ab05ap05fmt.Errorf("missing arguments")
	}
	ab05ap05Counter, ab05ap05Err := ab05ap05Uint(ab05ap05Args["failure_counter"], 32)
	if ab05ap05Err != nil {
		return nil, ab05ap05Err
	}
	var ab05ap05Increment uint64
	ab05ap05IncrementValue, ab05ap05HasIncrement := ab05ap05Args["counter_increment"]
	if ab05ap05HasIncrement {
		ab05ap05Increment, ab05ap05Err = ab05ap05Uint(ab05ap05IncrementValue, 32)
		if ab05ap05Err != nil {
			return nil, ab05ap05Err
		}
	}
	ab05ap05Slots, ab05ap05OK := ab05ap05Args["slots"].([]any)
	if !ab05ap05OK || len(ab05ap05Slots) != 9 {
		return nil, ab05ap05fmt.Errorf("expected nine slots")
	}
	ab05ap05Results := []any{}
	var ab05ap05Sum uint32
	ab05ap05Successful, ab05ap05Failed, ab05ap05Skipped := 0, 0, 0
	for _, ab05ap05Value := range ab05ap05Slots {
		ab05ap05Slot, ab05ap05OK := ab05ap05Value.(map[string]any)
		if !ab05ap05OK {
			return nil, ab05ap05fmt.Errorf("expected slot object")
		}
		ab05ap05Registered, ab05ap05OK := ab05ap05Slot["registered"].(bool)
		if !ab05ap05OK {
			return nil, ab05ap05fmt.Errorf("expected registered Boolean")
		}
		ab05ap05Capable, ab05ap05OK := ab05ap05Slot["capable"].(bool)
		if !ab05ap05OK {
			return nil, ab05ap05fmt.Errorf("expected capable Boolean")
		}
		ab05ap05PriorChanged, ab05ap05OK := ab05ap05Slot["prior_changed"].(bool)
		if !ab05ap05OK {
			return nil, ab05ap05fmt.Errorf("expected prior_changed Boolean")
		}
		ab05ap05Prior, ab05ap05Err := ab05ap05Uint(ab05ap05Slot["prior"], 32)
		if ab05ap05Err != nil {
			return nil, ab05ap05Err
		}
		ab05ap05Previous, ab05ap05Err := ab05ap05Uint(ab05ap05Slot["prior_previous"], 32)
		if ab05ap05Err != nil {
			return nil, ab05ap05Err
		}
		ab05ap05Outcome, ab05ap05OK := ab05ap05Slot["outcome"].(string)
		if !ab05ap05OK || (ab05ap05Outcome != "success" && ab05ap05Outcome != "read-failure" && ab05ap05Outcome != "parse-failure") {
			return nil, ab05ap05fmt.Errorf("invalid outcome")
		}
		var ab05ap05Current uint64
		if ab05ap05Outcome == "success" {
			ab05ap05Current, ab05ap05Err = ab05ap05Uint(ab05ap05Slot["value"], 32)
			if ab05ap05Err != nil {
				return nil, ab05ap05Err
			}
		} else if _, ab05ap05Present := ab05ap05Slot["value"]; ab05ap05Present {
			return nil, ab05ap05fmt.Errorf("value only allowed for success")
		}
		ab05ap05Result := map[string]any{"raw": ab05ap05Slot, "state": "Hypothesis", "prior_changed": ab05ap05PriorChanged, "changed": map[string]any{"state": "Unknown"}, "current": ab05ap05Quantity(uint32(0xffffffff)), "attempted": false, "outcome": "skipped"}
		if !ab05ap05Registered || !ab05ap05Capable {
			ab05ap05Skipped++
		} else {
			ab05ap05Previous = ab05ap05Prior
			ab05ap05Result["attempted"] = true
			ab05ap05Result["outcome"] = ab05ap05Outcome
			if ab05ap05Outcome != "success" {
				ab05ap05Failed++
			} else {
				ab05ap05Successful++
				ab05ap05Sum += uint32(ab05ap05Current)
				ab05ap05Result["current"] = ab05ap05Quantity(ab05ap05Current)
				ab05ap05Result["changed"] = map[string]any{"state": "Hypothesis", "calculation": "Proven", "value": ab05ap05Prior == 0xffffffff || uint32(ab05ap05Current)-uint32(ab05ap05Prior) >= 2}
			}
		}
		ab05ap05Result["current"].(map[string]any)["raw"] = ab05ap05Slot
		ab05ap05Result["changed"].(map[string]any)["raw"] = ab05ap05Slot
		ab05ap05Result["previous"] = ab05ap05Quantity(ab05ap05Previous)
		ab05ap05Result["previous"].(map[string]any)["raw"] = ab05ap05Slot
		ab05ap05Results = append(ab05ap05Results, ab05ap05Result)
	}
	ab05ap05Output["state"] = "complete"
	ab05ap05Fields["slots"] = map[string]any{"raw": ab05ap05Slots, "state": "Hypothesis", "outcomes": ab05ap05Results}
	ab05ap05Fields["sum"] = ab05ap05Quantity(ab05ap05Sum)
	ab05ap05Fields["counts"] = map[string]any{"raw": ab05ap05Slots, "state": "Hypothesis", "successful": ab05ap05Successful, "failed": ab05ap05Failed, "skipped": ab05ap05Skipped, "attempted": ab05ap05Successful + ab05ap05Failed}
	ab05ap05Output["partial_cycle"] = ab05ap05Failed+ab05ap05Skipped > 0
	ab05ap05Fields["failure_counter"] = ab05ap05Quantity(uint32(0))
	if ab05ap05Failed > 0 {
		if ab05ap05HasIncrement {
			ab05ap05Fields["failure_counter"] = ab05ap05Quantity(uint32(ab05ap05Counter + ab05ap05Increment))
		} else {
			ab05ap05Output["state"] = "partial"
			ab05ap05Fields["failure_counter"] = map[string]any{"state": "Unknown", "prior": ab05ap05Counter}
		}
	}
	ab05ap05Fields["sum"].(map[string]any)["raw"] = ab05ap05Slots
	ab05ap05Fields["failure_counter"].(map[string]any)["raw"] = ab05ap05Args["failure_counter"]
	if ab05ap05Failed > 0 {
		ab05ap05Fields["failure_counter"].(map[string]any)["native_increments"] = "Unknown"
	}
	return ab05ap05json.Marshal(ab05ap05Output)
}
