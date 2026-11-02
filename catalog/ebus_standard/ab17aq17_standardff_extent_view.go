package ebus_standard_catalog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func ab17aq17Value(ab17aq17Decoder *json.Decoder) error {
	ab17aq17Token, ab17aq17Err := ab17aq17Decoder.Token()
	if ab17aq17Err != nil {
		return ab17aq17Err
	}
	if ab17aq17Delimiter, ab17aq17OK := ab17aq17Token.(json.Delim); ab17aq17OK {
		switch ab17aq17Delimiter {
		case '{':
			ab17aq17Seen := map[string]bool{}
			for ab17aq17Decoder.More() {
				ab17aq17Key, ab17aq17Err := ab17aq17Decoder.Token()
				if ab17aq17Err != nil {
					return ab17aq17Err
				}
				ab17aq17Name, ab17aq17OK := ab17aq17Key.(string)
				if !ab17aq17OK || ab17aq17Seen[ab17aq17Name] {
					return fmt.Errorf("duplicate or invalid object key")
				}
				ab17aq17Seen[ab17aq17Name] = true
				if ab17aq17Err := ab17aq17Value(ab17aq17Decoder); ab17aq17Err != nil {
					return ab17aq17Err
				}
			}
		case '[':
			for ab17aq17Decoder.More() {
				if ab17aq17Err := ab17aq17Value(ab17aq17Decoder); ab17aq17Err != nil {
					return ab17aq17Err
				}
			}
		default:
			return fmt.Errorf("unexpected delimiter")
		}
		_, ab17aq17Err = ab17aq17Decoder.Token()
	}
	return ab17aq17Err
}

func ab17aq17Object(ab17aq17Raw json.RawMessage) (map[string]json.RawMessage, error) {
	ab17aq17Map := map[string]json.RawMessage{}
	if len(bytes.TrimSpace(ab17aq17Raw)) == 0 || bytes.TrimSpace(ab17aq17Raw)[0] != '{' {
		return nil, fmt.Errorf("expected object")
	}
	ab17aq17Err := json.Unmarshal(ab17aq17Raw, &ab17aq17Map)
	return ab17aq17Map, ab17aq17Err
}

func ab17aq17Byte(ab17aq17Raw json.RawMessage) (int, error) {
	var ab17aq17Number int
	if len(ab17aq17Raw) == 0 || bytes.Equal(bytes.TrimSpace(ab17aq17Raw), []byte("null")) {
		return 0, fmt.Errorf("expected integer byte")
	}
	if ab17aq17Err := json.Unmarshal(ab17aq17Raw, &ab17aq17Number); ab17aq17Err != nil {
		return 0, ab17aq17Err
	}
	if ab17aq17Number < 0 || ab17aq17Number > 255 {
		return 0, fmt.Errorf("integer outside byte range")
	}
	return ab17aq17Number, nil
}

func Ab17Aq17Evaluate(ab17aq17InputRaw json.RawMessage) (json.RawMessage, error) {
	ab17aq17Decoder := json.NewDecoder(bytes.NewReader(ab17aq17InputRaw))
	ab17aq17Decoder.UseNumber()
	if ab17aq17Err := ab17aq17Value(ab17aq17Decoder); ab17aq17Err != nil {
		return nil, ab17aq17Err
	}
	if _, ab17aq17Err := ab17aq17Decoder.Token(); ab17aq17Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab17aq17Input, ab17aq17Err := ab17aq17Object(ab17aq17InputRaw)
	if ab17aq17Err != nil {
		return nil, ab17aq17Err
	}
	ab17aq17Context, ab17aq17Err := ab17aq17Object(ab17aq17Input["context"])
	if ab17aq17Err != nil {
		return nil, ab17aq17Err
	}
	ab17aq17Fields, ab17aq17Err := ab17aq17Object(ab17aq17Input["fields"])
	if ab17aq17Err != nil {
		return nil, ab17aq17Err
	}
	ab17aq17Numbers := map[string]int{}
	for _, ab17aq17Key := range []string{"primary", "secondary", "operation"} {
		if ab17aq17Raw, ab17aq17OK := ab17aq17Context[ab17aq17Key]; ab17aq17OK {
			ab17aq17Number, ab17aq17Err := ab17aq17Byte(ab17aq17Raw)
			if ab17aq17Err != nil {
				return nil, ab17aq17Err
			}
			ab17aq17Numbers[ab17aq17Key] = ab17aq17Number
		}
	}
	for _, ab17aq17Key := range []string{"direction", "view"} {
		if ab17aq17Raw, ab17aq17OK := ab17aq17Context[ab17aq17Key]; ab17aq17OK {
			var ab17aq17String string
			if bytes.Equal(bytes.TrimSpace(ab17aq17Raw), []byte("null")) {
				return nil, fmt.Errorf("expected context string")
			}
			if ab17aq17Err := json.Unmarshal(ab17aq17Raw, &ab17aq17String); ab17aq17Err != nil {
				return nil, ab17aq17Err
			}
		}
	}
	ab17aq17Results := map[string]any{}
	ab17aq17Lengths := map[string]int{}
	for ab17aq17Key, ab17aq17Raw := range ab17aq17Fields {
		if len(bytes.TrimSpace(ab17aq17Raw)) == 0 || bytes.TrimSpace(ab17aq17Raw)[0] != '[' {
			return nil, fmt.Errorf("expected byte array")
		}
		var ab17aq17Bytes []json.RawMessage
		if ab17aq17Err := json.Unmarshal(ab17aq17Raw, &ab17aq17Bytes); ab17aq17Err != nil {
			return nil, ab17aq17Err
		}
		for _, ab17aq17ByteRaw := range ab17aq17Bytes {
			if _, ab17aq17Err := ab17aq17Byte(ab17aq17ByteRaw); ab17aq17Err != nil {
				return nil, ab17aq17Err
			}
		}
		ab17aq17Lengths[ab17aq17Key] = len(ab17aq17Bytes)
		ab17aq17Results[ab17aq17Key] = map[string]any{"raw": ab17aq17Raw, "state": "Unknown"}
	}
	ab17aq17Arguments := map[string]json.RawMessage{}
	if ab17aq17Raw, ab17aq17OK := ab17aq17Input["arguments"]; ab17aq17OK {
		ab17aq17Arguments, ab17aq17Err = ab17aq17Object(ab17aq17Raw)
		if ab17aq17Err != nil {
			return nil, ab17aq17Err
		}
	}
	ab17aq17Declared := 0
	if ab17aq17Raw, ab17aq17OK := ab17aq17Arguments["declared_length"]; ab17aq17OK {
		ab17aq17Declared, ab17aq17Err = ab17aq17Byte(ab17aq17Raw)
		if ab17aq17Err != nil {
			return nil, ab17aq17Err
		}
	}
	ab17aq17Output := map[string]any{}
	for ab17aq17Key, ab17aq17Raw := range ab17aq17Input {
		ab17aq17Output[ab17aq17Key] = ab17aq17Raw
	}
	ab17aq17Output["state"] = "Unknown"
	ab17aq17Output["results"] = ab17aq17Results
	ab17aq17Output["physical"] = "Unknown"
	ab17aq17Output["qualified"] = false
	_, ab17aq17Operation := ab17aq17Context["operation"]
	var ab17aq17Direction, ab17aq17View string
	_ = json.Unmarshal(ab17aq17Context["direction"], &ab17aq17Direction)
	_ = json.Unmarshal(ab17aq17Context["view"], &ab17aq17View)
	if ab17aq17Numbers["primary"] != 255 || (ab17aq17Numbers["secondary"] != 4 && ab17aq17Numbers["secondary"] != 6) || ab17aq17Operation || ab17aq17Direction != "positive" || ab17aq17View != "body-extent" {
		return json.Marshal(ab17aq17Output)
	}
	if _, ab17aq17OK := ab17aq17Arguments["declared_length"]; !ab17aq17OK || len(ab17aq17Arguments) != 1 {
		return nil, fmt.Errorf("expected declared_length argument only")
	}
	ab17aq17Available, ab17aq17Present := ab17aq17Lengths["body"]
	ab17aq17Body := map[string]any{"state": "Unknown", "calculation": "Proven", "declared_length": ab17aq17Declared}
	if ab17aq17Present {
		ab17aq17Body["raw"] = ab17aq17Fields["body"]
		ab17aq17Body["available_length"] = ab17aq17Available
	} else {
		ab17aq17Body["state"] = "partial"
	}
	ab17aq17Flags := map[string]bool{}
	ab17aq17Complete := false
	if ab17aq17Numbers["secondary"] == 4 {
		ab17aq17Flags["declared_minimum_satisfied"] = ab17aq17Declared >= 2
		ab17aq17Flags["initialized_minimum_satisfied"] = ab17aq17Present && ab17aq17Available >= 2
		ab17aq17Complete = ab17aq17Declared >= 2 && ab17aq17Declared <= 10 && ab17aq17Present && ab17aq17Available >= ab17aq17Declared
		ab17aq17Flags["bounded_minimum"] = ab17aq17Complete
	} else {
		ab17aq17Tail := ab17aq17Declared >= 1 && (ab17aq17Declared-1)%2 == 1
		ab17aq17Flags["declared_incomplete_tail"] = ab17aq17Tail
		if ab17aq17Declared >= 1 {
			ab17aq17Body["declared_pair_count"] = (ab17aq17Declared - 1) / 2
		}
		ab17aq17Bounded := ab17aq17Available
		if ab17aq17Bounded > ab17aq17Declared {
			ab17aq17Bounded = ab17aq17Declared
		}
		if ab17aq17Present && ab17aq17Bounded >= 1 {
			ab17aq17Body["available_pair_count"] = (ab17aq17Bounded - 1) / 2
		}
		ab17aq17Complete = ab17aq17Declared >= 1 && ab17aq17Present && ab17aq17Available >= ab17aq17Declared && !ab17aq17Tail
		ab17aq17Flags["whole_declared_span"] = ab17aq17Complete
	}
	ab17aq17Body["flags"] = ab17aq17Flags
	ab17aq17Results["body"] = ab17aq17Body
	ab17aq17Output["state"] = "partial"
	if ab17aq17Complete {
		ab17aq17Output["state"] = "complete"
	}
	return json.Marshal(ab17aq17Output)
}
