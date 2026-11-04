package vaillant

import (
	ab09ap09bytes "bytes"
	ab09ap09json "encoding/json"
	ab09ap09fmt "fmt"
	ab09ap09io "io"
	ab09ap09strconv "strconv"
)

func ab09ap09Value(ab09ap09Decoder *ab09ap09json.Decoder) (any, error) {
	ab09ap09Token, ab09ap09Err := ab09ap09Decoder.Token()
	if ab09ap09Err != nil {
		return nil, ab09ap09Err
	}
	switch ab09ap09Token {
	case ab09ap09json.Delim('{'):
		ab09ap09Object := map[string]any{}
		for ab09ap09Decoder.More() {
			ab09ap09KeyToken, ab09ap09Err := ab09ap09Decoder.Token()
			if ab09ap09Err != nil {
				return nil, ab09ap09Err
			}
			ab09ap09Key := ab09ap09KeyToken.(string)
			if _, ab09ap09Exists := ab09ap09Object[ab09ap09Key]; ab09ap09Exists {
				return nil, ab09ap09fmt.Errorf("duplicate key %q", ab09ap09Key)
			}
			ab09ap09Child, ab09ap09Err := ab09ap09Value(ab09ap09Decoder)
			if ab09ap09Err != nil {
				return nil, ab09ap09Err
			}
			ab09ap09Object[ab09ap09Key] = ab09ap09Child
		}
		_, ab09ap09Err := ab09ap09Decoder.Token()
		return ab09ap09Object, ab09ap09Err
	case ab09ap09json.Delim('['):
		ab09ap09Array := []any{}
		for ab09ap09Decoder.More() {
			ab09ap09Child, ab09ap09Err := ab09ap09Value(ab09ap09Decoder)
			if ab09ap09Err != nil {
				return nil, ab09ap09Err
			}
			ab09ap09Array = append(ab09ap09Array, ab09ap09Child)
		}
		_, ab09ap09Err := ab09ap09Decoder.Token()
		return ab09ap09Array, ab09ap09Err
	}
	return ab09ap09Token, nil
}

func ab09ap09Unsigned(ab09ap09Value any, ab09ap09Bits int) (uint64, error) {
	ab09ap09Number, ab09ap09OK := ab09ap09Value.(ab09ap09json.Number)
	if !ab09ap09OK {
		return 0, ab09ap09fmt.Errorf("expected unsigned integer")
	}
	return ab09ap09strconv.ParseUint(string(ab09ap09Number), 10, ab09ap09Bits)
}

func Ab09Ap09Evaluate(ab09ap09Input ab09ap09json.RawMessage) (ab09ap09json.RawMessage, error) {
	ab09ap09Decoder := ab09ap09json.NewDecoder(ab09ap09bytes.NewReader(ab09ap09Input))
	ab09ap09Decoder.UseNumber()
	ab09ap09Root, ab09ap09Err := ab09ap09Value(ab09ap09Decoder)
	if ab09ap09Err != nil {
		return nil, ab09ap09Err
	}
	if _, ab09ap09Err = ab09ap09Decoder.Token(); ab09ap09Err != ab09ap09io.EOF {
		return nil, ab09ap09fmt.Errorf("trailing JSON")
	}
	ab09ap09Object, ab09ap09OK := ab09ap09Root.(map[string]any)
	if !ab09ap09OK {
		return nil, ab09ap09fmt.Errorf("expected input object")
	}
	ab09ap09Context, ab09ap09OK := ab09ap09Object["context"].(map[string]any)
	if !ab09ap09OK {
		return nil, ab09ap09fmt.Errorf("expected context object")
	}
	ab09ap09Raw, ab09ap09OK := ab09ap09Object["fields"].(map[string]any)
	if !ab09ap09OK {
		return nil, ab09ap09fmt.Errorf("expected fields object")
	}
	for _, ab09ap09Field := range ab09ap09Raw {
		ab09ap09Array, ab09ap09OK := ab09ap09Field.([]any)
		if !ab09ap09OK {
			return nil, ab09ap09fmt.Errorf("expected byte array")
		}
		for _, ab09ap09Byte := range ab09ap09Array {
			if _, ab09ap09Err := ab09ap09Unsigned(ab09ap09Byte, 8); ab09ap09Err != nil {
				return nil, ab09ap09Err
			}
		}
	}
	ab09ap09Arguments := map[string]any{}
	if ab09ap09Value, ab09ap09Present := ab09ap09Object["arguments"]; ab09ap09Present {
		ab09ap09Arguments, ab09ap09OK = ab09ap09Value.(map[string]any)
		if !ab09ap09OK {
			return nil, ab09ap09fmt.Errorf("expected arguments object")
		}
	}
	for ab09ap09Key, ab09ap09Bits := range map[string]int{"controller_type": 8, "spn": 16} {
		if ab09ap09Value, ab09ap09Present := ab09ap09Arguments[ab09ap09Key]; ab09ap09Present {
			if _, ab09ap09Err := ab09ap09Unsigned(ab09ap09Value, ab09ap09Bits); ab09ap09Err != nil {
				return nil, ab09ap09Err
			}
		}
	}
	if ab09ap09Family, ab09ap09Present := ab09ap09Arguments["family_match"]; ab09ap09Present && ab09ap09Family != nil {
		if _, ab09ap09OK := ab09ap09Family.(bool); !ab09ap09OK {
			return nil, ab09ap09fmt.Errorf("expected Boolean or null family_match")
		}
	}
	ab09ap09Supported := true
	for ab09ap09Key, ab09ap09Expected := range map[string]uint64{"primary": 181, "secondary": 85, "operation": 0} {
		ab09ap09Value := ab09ap09Context[ab09ap09Key]
		if ab09ap09Value == nil {
			ab09ap09Supported = false
			continue
		}
		ab09ap09Number, ab09ap09Err := ab09ap09Unsigned(ab09ap09Value, 8)
		if ab09ap09Err != nil {
			return nil, ab09ap09Err
		}
		ab09ap09Supported = ab09ap09Supported && ab09ap09Number == ab09ap09Expected
	}
	for ab09ap09Key, ab09ap09Expected := range map[string]string{"direction": "local", "view": "manager-gate"} {
		ab09ap09Value := ab09ap09Context[ab09ap09Key]
		if ab09ap09Value == nil {
			ab09ap09Supported = false
			continue
		}
		ab09ap09String, ab09ap09OK := ab09ap09Value.(string)
		if !ab09ap09OK {
			return nil, ab09ap09fmt.Errorf("expected string selector")
		}
		ab09ap09Supported = ab09ap09Supported && ab09ap09String == ab09ap09Expected
	}
	ab09ap09Fields := map[string]any{}
	ab09ap09Output := map[string]any{"context": ab09ap09Context, "raw": ab09ap09Raw, "state": "Unknown", "fields": ab09ap09Fields, "physical": "Unknown", "qualified": false}
	if ab09ap09Value, ab09ap09Present := ab09ap09Object["arguments"]; ab09ap09Present {
		ab09ap09Output["arguments"] = ab09ap09Value
	}
	if !ab09ap09Supported {
		return ab09ap09json.Marshal(ab09ap09Output)
	}
	for _, ab09ap09Key := range []string{"controller_type", "spn", "family_match"} {
		if _, ab09ap09Present := ab09ap09Arguments[ab09ap09Key]; !ab09ap09Present {
			return nil, ab09ap09fmt.Errorf("missing argument %s", ab09ap09Key)
		}
	}
	ab09ap09Reply := map[string]any{"state": "incomplete"}
	if ab09ap09Value, ab09ap09Present := ab09ap09Raw["reply"]; ab09ap09Present {
		ab09ap09Reply["raw"] = ab09ap09Value
	}
	ab09ap09Fields["reply"] = ab09ap09Reply
	ab09ap09Output["state"] = "partial"
	ab09ap09Array, _ := ab09ap09Raw["reply"].([]any)
	if len(ab09ap09Array) != 1 {
		return ab09ap09json.Marshal(ab09ap09Output)
	}
	ab09ap09Byte, _ := ab09ap09Unsigned(ab09ap09Array[0], 8)
	ab09ap09Type, _ := ab09ap09Unsigned(ab09ap09Arguments["controller_type"], 8)
	ab09ap09SPN, _ := ab09ap09Unsigned(ab09ap09Arguments["spn"], 16)
	ab09ap09Bit0 := ab09ap09Byte&1 != 0
	ab09ap09Bit1 := ab09ap09Byte&2 != 0
	ab09ap09Family, _ := ab09ap09Arguments["family_match"].(bool)
	ab09ap09Reply["state"] = "Hypothesis"
	ab09ap09Reply["calculation"] = "Proven"
	ab09ap09Reply["flags"] = map[string]any{"bit0": ab09ap09Bit0, "bit1": ab09ap09Bit1}
	ab09ap09Reply["unmodeled"] = ab09ap09Byte & 252
	ab09ap09Reply["eligible"] = ab09ap09Type == 2 && ab09ap09SPN == 385 && ab09ap09Family && ab09ap09Bit0 && ab09ap09Bit1
	ab09ap09Reply["first_bit_policy"] = ab09ap09Bit0
	ab09ap09Reply["first_bit_policy_label"] = "caller predicate"
	ab09ap09Reply["family_role"] = 0x15
	ab09ap09Reply["route_role"] = 0x08
	ab09ap09Reply["family_applicability"] = "Hypothesis"
	if ab09ap09Arguments["family_match"] == nil {
		ab09ap09Reply["family_applicability"] = "Unknown"
	}
	ab09ap09Reply["unknown"] = []string{"exact family prefix", "reply adaptation", "actual manager state", "physical support"}
	ab09ap09Output["state"] = "complete"
	return ab09ap09json.Marshal(ab09ap09Output)
}
