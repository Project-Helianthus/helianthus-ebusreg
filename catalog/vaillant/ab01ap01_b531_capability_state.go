package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func ab01ap01Read(ab01ap01Decoder *json.Decoder) (any, error) {
	ab01ap01Token, ab01ap01Err := ab01ap01Decoder.Token()
	if ab01ap01Err != nil {
		return nil, ab01ap01Err
	}
	if ab01ap01Delimiter, ab01ap01OK := ab01ap01Token.(json.Delim); ab01ap01OK {
		switch ab01ap01Delimiter {
		case '{':
			ab01ap01Object := map[string]any{}
			for ab01ap01Decoder.More() {
				ab01ap01KeyToken, ab01ap01Err := ab01ap01Decoder.Token()
				if ab01ap01Err != nil {
					return nil, ab01ap01Err
				}
				ab01ap01Key := ab01ap01KeyToken.(string)
				if _, ab01ap01Exists := ab01ap01Object[ab01ap01Key]; ab01ap01Exists {
					return nil, fmt.Errorf("duplicate key %q", ab01ap01Key)
				}
				ab01ap01Value, ab01ap01Err := ab01ap01Read(ab01ap01Decoder)
				if ab01ap01Err != nil {
					return nil, ab01ap01Err
				}
				ab01ap01Object[ab01ap01Key] = ab01ap01Value
			}
			_, ab01ap01Err := ab01ap01Decoder.Token()
			return ab01ap01Object, ab01ap01Err
		case '[':
			ab01ap01Array := []any{}
			for ab01ap01Decoder.More() {
				ab01ap01Value, ab01ap01Err := ab01ap01Read(ab01ap01Decoder)
				if ab01ap01Err != nil {
					return nil, ab01ap01Err
				}
				ab01ap01Array = append(ab01ap01Array, ab01ap01Value)
			}
			_, ab01ap01Err := ab01ap01Decoder.Token()
			return ab01ap01Array, ab01ap01Err
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	return ab01ap01Token, nil
}

func ab01ap01Uint(ab01ap01Value any, ab01ap01Bits int) (uint64, error) {
	ab01ap01Number, ab01ap01OK := ab01ap01Value.(json.Number)
	if !ab01ap01OK {
		return 0, fmt.Errorf("expected unsigned integer")
	}
	return strconv.ParseUint(string(ab01ap01Number), 10, ab01ap01Bits)
}

func Ab01Ap01Evaluate(ab01ap01Input json.RawMessage) (json.RawMessage, error) {
	ab01ap01Decoder := json.NewDecoder(bytes.NewReader(ab01ap01Input))
	ab01ap01Decoder.UseNumber()
	ab01ap01Value, ab01ap01Err := ab01ap01Read(ab01ap01Decoder)
	if ab01ap01Err != nil {
		return nil, ab01ap01Err
	}
	if _, ab01ap01Err = ab01ap01Decoder.Token(); ab01ap01Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab01ap01Root, ab01ap01OK := ab01ap01Value.(map[string]any)
	if !ab01ap01OK {
		return nil, fmt.Errorf("expected object")
	}
	ab01ap01Context, ab01ap01OK := ab01ap01Root["context"].(map[string]any)
	if !ab01ap01OK {
		return nil, fmt.Errorf("expected context object")
	}
	ab01ap01Raw, ab01ap01OK := ab01ap01Root["fields"].(map[string]any)
	if !ab01ap01OK {
		return nil, fmt.Errorf("expected fields object")
	}
	ab01ap01Supported := true
	for ab01ap01Key, ab01ap01Expected := range map[string]uint64{"primary": 181, "secondary": 49, "operation": 1} {
		ab01ap01Item, ab01ap01Exists := ab01ap01Context[ab01ap01Key]
		if !ab01ap01Exists || ab01ap01Item == nil {
			ab01ap01Supported = false
			continue
		}
		ab01ap01Number, ab01ap01Err := ab01ap01Uint(ab01ap01Item, 8)
		if ab01ap01Err != nil {
			return nil, ab01ap01Err
		}
		if ab01ap01Number != ab01ap01Expected {
			ab01ap01Supported = false
		}
	}
	for _, ab01ap01Key := range []string{"direction", "view"} {
		if ab01ap01Item := ab01ap01Context[ab01ap01Key]; ab01ap01Item != nil {
			if _, ab01ap01OK := ab01ap01Item.(string); !ab01ap01OK {
				return nil, fmt.Errorf("expected string selector")
			}
		}
	}
	ab01ap01View, _ := ab01ap01Context["view"].(string)
	ab01ap01Supported = ab01ap01Supported && ab01ap01Context["direction"] == "local" && (ab01ap01View == "capability" || ab01ap01View == "select")
	ab01ap01Fields := map[string]any{}
	ab01ap01State := "complete"
	for ab01ap01Key, ab01ap01Item := range ab01ap01Raw {
		ab01ap01Bytes, ab01ap01OK := ab01ap01Item.([]any)
		if !ab01ap01OK {
			return nil, fmt.Errorf("expected byte array")
		}
		for _, ab01ap01Byte := range ab01ap01Bytes {
			if _, ab01ap01Err := ab01ap01Uint(ab01ap01Byte, 8); ab01ap01Err != nil {
				return nil, ab01ap01Err
			}
		}
		ab01ap01Fields[ab01ap01Key] = map[string]any{"raw": ab01ap01Item, "state": "Unknown", "unmodeled": true}
	}
	if ab01ap01Supported {
		for ab01ap01Key, ab01ap01Item := range ab01ap01Raw {
			if len(ab01ap01Item.([]any)) == 0 {
				ab01ap01Fields[ab01ap01Key] = map[string]any{"raw": ab01ap01Item, "state": "incomplete", "unmodeled": true}
				ab01ap01State = "partial"
			}
		}
	}
	ab01ap01Output := map[string]any{"context": ab01ap01Context, "raw": ab01ap01Raw, "fields": ab01ap01Fields, "physical": "Unknown", "qualified": false}
	ab01ap01Arguments := map[string]any{}
	if ab01ap01Item, ab01ap01Exists := ab01ap01Root["arguments"]; ab01ap01Exists {
		ab01ap01Arguments, ab01ap01OK = ab01ap01Item.(map[string]any)
		if !ab01ap01OK {
			return nil, fmt.Errorf("expected arguments object")
		}
		ab01ap01Output["arguments"] = ab01ap01Arguments
	}
	for _, ab01ap01Key := range []string{"registered", "prior", "result"} {
		if ab01ap01Item, ab01ap01Exists := ab01ap01Arguments[ab01ap01Key]; ab01ap01Exists {
			if _, ab01ap01Err := ab01ap01Uint(ab01ap01Item, 32); ab01ap01Err != nil {
				return nil, ab01ap01Err
			}
		}
	}
	ab01ap01Rows := []any(nil)
	if ab01ap01Item, ab01ap01Exists := ab01ap01Arguments["rows"]; ab01ap01Exists {
		ab01ap01Rows, ab01ap01OK = ab01ap01Item.([]any)
		if !ab01ap01OK || len(ab01ap01Rows) != 9 {
			return nil, fmt.Errorf("expected nine rows")
		}
		for _, ab01ap01Row := range ab01ap01Rows {
			ab01ap01Object, ab01ap01OK := ab01ap01Row.(map[string]any)
			if !ab01ap01OK {
				return nil, fmt.Errorf("expected row object")
			}
			for _, ab01ap01Key := range []string{"a", "b", "c"} {
				if _, ab01ap01Err := ab01ap01Uint(ab01ap01Object[ab01ap01Key], 32); ab01ap01Err != nil {
					return nil, ab01ap01Err
				}
			}
		}
	}
	if !ab01ap01Supported {
		ab01ap01Output["state"] = "Unknown"
		return json.Marshal(ab01ap01Output)
	}
	if ab01ap01View == "capability" {
		ab01ap01Numbers := map[string]uint64{}
		for _, ab01ap01Key := range []string{"registered", "prior", "result"} {
			ab01ap01Number, ab01ap01Err := ab01ap01Uint(ab01ap01Arguments[ab01ap01Key], 32)
			if ab01ap01Err != nil {
				return nil, ab01ap01Err
			}
			ab01ap01Numbers[ab01ap01Key] = ab01ap01Number
		}
		ab01ap01Updated := ab01ap01Numbers["prior"]
		var ab01ap01Outcome string
		if ab01ap01Numbers["registered"] == 0 {
			ab01ap01Outcome = "skipped"
		} else {
			switch ab01ap01Numbers["result"] {
			case 0:
				ab01ap01Updated = 0
				ab01ap01Outcome = "cleared"
			case 1:
				if ab01ap01Updated == 0 {
					ab01ap01Updated = 1
				}
				ab01ap01Outcome = "admitted"
			default:
				ab01ap01Outcome = "Unknown"
			}
		}
		ab01ap01Output["capability"] = map[string]any{"state": "Hypothesis", "calculation": "Proven", "marker": 65535, "updated": ab01ap01Updated, "cached_bit0": ab01ap01Updated&1 != 0, "whole_nonzero": ab01ap01Updated != 0, "outcome": ab01ap01Outcome}
	} else {
		if ab01ap01Rows == nil {
			return nil, fmt.Errorf("missing rows")
		}
		ab01ap01Index := 9
		for ab01ap01Position, ab01ap01Row := range ab01ap01Rows {
			ab01ap01Object := ab01ap01Row.(map[string]any)
			ab01ap01Match := true
			for _, ab01ap01Key := range []string{"a", "b", "c"} {
				ab01ap01Number, _ := ab01ap01Uint(ab01ap01Object[ab01ap01Key], 32)
				ab01ap01Match = ab01ap01Match && ab01ap01Number != 0
			}
			if ab01ap01Match {
				ab01ap01Index = ab01ap01Position
				break
			}
		}
		ab01ap01Output["select"] = map[string]any{"state": "Hypothesis", "calculation": "Proven", "index": ab01ap01Index}
	}
	ab01ap01Output["state"] = ab01ap01State
	return json.Marshal(ab01ap01Output)
}
