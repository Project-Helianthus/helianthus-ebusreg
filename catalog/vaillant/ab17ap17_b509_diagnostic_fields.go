package vaillant

import (
	ab17ap17JSON "encoding/json"
	ab17ap17Fmt "fmt"
	ab17ap17IO "io"
	ab17ap17Strconv "strconv"
	ab17ap17Strings "strings"
)

func ab17ap17Read(ab17ap17Decoder *ab17ap17JSON.Decoder) (any, error) {
	ab17ap17Token, ab17ap17Err := ab17ap17Decoder.Token()
	if ab17ap17Err != nil {
		return nil, ab17ap17Err
	}
	switch ab17ap17Token {
	case ab17ap17JSON.Delim('{'):
		ab17ap17Object := map[string]any{}
		for ab17ap17Decoder.More() {
			ab17ap17KeyToken, ab17ap17Err := ab17ap17Decoder.Token()
			if ab17ap17Err != nil {
				return nil, ab17ap17Err
			}
			ab17ap17Key, ab17ap17OK := ab17ap17KeyToken.(string)
			if !ab17ap17OK {
				return nil, ab17ap17Fmt.Errorf("object key required")
			}
			if _, ab17ap17Exists := ab17ap17Object[ab17ap17Key]; ab17ap17Exists {
				return nil, ab17ap17Fmt.Errorf("duplicate key %q", ab17ap17Key)
			}
			ab17ap17Value, ab17ap17Err := ab17ap17Read(ab17ap17Decoder)
			if ab17ap17Err != nil {
				return nil, ab17ap17Err
			}
			ab17ap17Object[ab17ap17Key] = ab17ap17Value
		}
		ab17ap17End, ab17ap17Err := ab17ap17Decoder.Token()
		if ab17ap17Err != nil || ab17ap17End != ab17ap17JSON.Delim('}') {
			return nil, ab17ap17Fmt.Errorf("unclosed object")
		}
		return ab17ap17Object, nil
	case ab17ap17JSON.Delim('['):
		ab17ap17Array := []any{}
		for ab17ap17Decoder.More() {
			ab17ap17Value, ab17ap17Err := ab17ap17Read(ab17ap17Decoder)
			if ab17ap17Err != nil {
				return nil, ab17ap17Err
			}
			ab17ap17Array = append(ab17ap17Array, ab17ap17Value)
		}
		ab17ap17End, ab17ap17Err := ab17ap17Decoder.Token()
		if ab17ap17Err != nil || ab17ap17End != ab17ap17JSON.Delim(']') {
			return nil, ab17ap17Fmt.Errorf("unclosed array")
		}
		return ab17ap17Array, nil
	case ab17ap17JSON.Delim('}'), ab17ap17JSON.Delim(']'):
		return nil, ab17ap17Fmt.Errorf("unexpected delimiter")
	default:
		return ab17ap17Token, nil
	}
}

func ab17ap17Unsigned(ab17ap17Value any, ab17ap17Max uint64) (uint64, bool) {
	ab17ap17Number, ab17ap17OK := ab17ap17Value.(ab17ap17JSON.Number)
	if !ab17ap17OK {
		return 0, false
	}
	ab17ap17Parsed, ab17ap17Err := ab17ap17Strconv.ParseUint(string(ab17ap17Number), 10, 64)
	return ab17ap17Parsed, ab17ap17Err == nil && ab17ap17Parsed <= ab17ap17Max
}

func Ab17Ap17Evaluate(ab17ap17Input ab17ap17JSON.RawMessage) (ab17ap17JSON.RawMessage, error) {
	ab17ap17Decoder := ab17ap17JSON.NewDecoder(ab17ap17Strings.NewReader(string(ab17ap17Input)))
	ab17ap17Decoder.UseNumber()
	ab17ap17Value, ab17ap17Err := ab17ap17Read(ab17ap17Decoder)
	if ab17ap17Err != nil {
		return nil, ab17ap17Err
	}
	if _, ab17ap17Err = ab17ap17Decoder.Token(); ab17ap17Err != ab17ap17IO.EOF {
		return nil, ab17ap17Fmt.Errorf("trailing JSON")
	}
	ab17ap17InputObject, ab17ap17OK := ab17ap17Value.(map[string]any)
	if !ab17ap17OK {
		return nil, ab17ap17Fmt.Errorf("input object required")
	}
	ab17ap17Context, ab17ap17OK := ab17ap17InputObject["context"].(map[string]any)
	if !ab17ap17OK {
		return nil, ab17ap17Fmt.Errorf("context object required")
	}
	ab17ap17Raw, ab17ap17OK := ab17ap17InputObject["fields"].(map[string]any)
	if !ab17ap17OK {
		return nil, ab17ap17Fmt.Errorf("fields object required")
	}
	for _, ab17ap17Key := range []string{"primary", "secondary", "operation"} {
		if ab17ap17Value, ab17ap17Exists := ab17ap17Context[ab17ap17Key]; ab17ap17Exists && ab17ap17Value != nil {
			if _, ab17ap17OK := ab17ap17Unsigned(ab17ap17Value, 255); !ab17ap17OK {
				return nil, ab17ap17Fmt.Errorf("%s must be u8 or null", ab17ap17Key)
			}
		}
	}
	for _, ab17ap17Key := range []string{"direction", "view"} {
		if ab17ap17Value, ab17ap17Exists := ab17ap17Context[ab17ap17Key]; ab17ap17Exists && ab17ap17Value != nil {
			if _, ab17ap17OK := ab17ap17Value.(string); !ab17ap17OK {
				return nil, ab17ap17Fmt.Errorf("%s must be string or null", ab17ap17Key)
			}
		}
	}
	for ab17ap17Key, ab17ap17Value := range ab17ap17Raw {
		ab17ap17Array, ab17ap17OK := ab17ap17Value.([]any)
		if !ab17ap17OK {
			return nil, ab17ap17Fmt.Errorf("%s must be a byte array", ab17ap17Key)
		}
		for _, ab17ap17Byte := range ab17ap17Array {
			if _, ab17ap17OK := ab17ap17Unsigned(ab17ap17Byte, 255); !ab17ap17OK {
				return nil, ab17ap17Fmt.Errorf("%s contains a non-byte", ab17ap17Key)
			}
		}
	}
	ab17ap17Results := map[string]any{}
	ab17ap17Output := map[string]any{"context": ab17ap17Context, "raw": ab17ap17Raw, "state": "Unknown", "fields": ab17ap17Results, "physical": "Unknown", "qualified": false}
	if ab17ap17ArgumentsValue, ab17ap17Exists := ab17ap17InputObject["arguments"]; ab17ap17Exists {
		ab17ap17Arguments, ab17ap17OK := ab17ap17ArgumentsValue.(map[string]any)
		if !ab17ap17OK {
			return nil, ab17ap17Fmt.Errorf("arguments object required")
		}
		if ab17ap17Requested, ab17ap17Exists := ab17ap17Arguments["requested_address"]; ab17ap17Exists {
			if _, ab17ap17OK := ab17ap17Unsigned(ab17ap17Requested, 4294967295); !ab17ap17OK {
				return nil, ab17ap17Fmt.Errorf("requested_address must be u32")
			}
		}
		ab17ap17Output["arguments"] = ab17ap17Arguments
	}
	ab17ap17Primary, ab17ap17PrimaryOK := ab17ap17Unsigned(ab17ap17Context["primary"], 255)
	ab17ap17Secondary, ab17ap17SecondaryOK := ab17ap17Unsigned(ab17ap17Context["secondary"], 255)
	ab17ap17Operation, ab17ap17OperationOK := ab17ap17Unsigned(ab17ap17Context["operation"], 255)
	ab17ap17Direction, _ := ab17ap17Context["direction"].(string)
	ab17ap17View, _ := ab17ap17Context["view"].(string)
	if !ab17ap17PrimaryOK || !ab17ap17SecondaryOK || !ab17ap17OperationOK || ab17ap17Primary != 181 || ab17ap17Secondary != 9 ||
		(ab17ap17Direction != "request" && ab17ap17Direction != "response") ||
		((ab17ap17Operation != 52 || ab17ap17View != "diagnostic-get") && (ab17ap17Operation != 53 || ab17ap17View != "diagnostic-put")) {
		return ab17ap17JSON.Marshal(ab17ap17Output)
	}
	ab17ap17Output["state"] = "complete"
	ab17ap17Output["grade"] = "Hypothesis"
	ab17ap17Output["correlation"] = "Unknown"
	ab17ap17Expected := map[string]int{}
	if ab17ap17Direction == "request" || ab17ap17Operation == 52 {
		ab17ap17Expected["address"] = 4
	}
	if ab17ap17Operation == 52 && ab17ap17Direction == "request" {
		ab17ap17Expected["count"] = 1
	} else if ab17ap17Operation == 53 && ab17ap17Direction == "response" {
		ab17ap17Expected["status"] = 1
	} else {
		ab17ap17Expected["data"] = -1
	}
	for ab17ap17Key, ab17ap17RawValue := range ab17ap17Raw {
		ab17ap17Results[ab17ap17Key] = map[string]any{"raw": ab17ap17RawValue, "state": "Unknown"}
	}
	for ab17ap17Key, ab17ap17Width := range ab17ap17Expected {
		ab17ap17RawValue, ab17ap17Exists := ab17ap17Raw[ab17ap17Key]
		ab17ap17Result := map[string]any{"state": "incomplete", "grade": "Hypothesis"}
		ab17ap17Results[ab17ap17Key] = ab17ap17Result
		if ab17ap17Exists {
			ab17ap17Result["raw"] = ab17ap17RawValue
		}
		ab17ap17Array, _ := ab17ap17RawValue.([]any)
		if !ab17ap17Exists || (ab17ap17Width >= 0 && len(ab17ap17Array) != ab17ap17Width) {
			ab17ap17Output["state"] = "partial"
			continue
		}
		ab17ap17Result["state"] = "complete"
		if ab17ap17Key == "data" {
			ab17ap17Result["length"] = len(ab17ap17Array)
			ab17ap17Result["unmodeled"] = true
			continue
		}
		var ab17ap17Number uint64
		for ab17ap17Index, ab17ap17Byte := range ab17ap17Array {
			ab17ap17Part, _ := ab17ap17Unsigned(ab17ap17Byte, 255)
			ab17ap17Number |= ab17ap17Part << (8 * ab17ap17Index)
		}
		ab17ap17Result["number"] = ab17ap17Number
		ab17ap17Result["calculation"] = "Proven"
		if ab17ap17Key == "status" {
			ab17ap17Result["candidate_domain"] = ab17ap17Number <= 6
		}
	}
	return ab17ap17JSON.Marshal(ab17ap17Output)
}
