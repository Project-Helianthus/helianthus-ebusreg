package vaillant

import (
	ab11aq11bytes "bytes"
	ab11aq11json "encoding/json"
	ab11aq11fmt "fmt"
	ab11aq11io "io"
	ab11aq11big "math/big"
	ab11aq11strconv "strconv"
	ab11aq11utf8 "unicode/utf8"
)

func ab11aq11Read(ab11aq11Decoder *ab11aq11json.Decoder) (any, error) {
	ab11aq11Token, ab11aq11Err := ab11aq11Decoder.Token()
	if ab11aq11Err != nil {
		return nil, ab11aq11Err
	}
	if ab11aq11Delimiter, ab11aq11OK := ab11aq11Token.(ab11aq11json.Delim); ab11aq11OK {
		switch ab11aq11Delimiter {
		case '{':
			ab11aq11Object := map[string]any{}
			for ab11aq11Decoder.More() {
				ab11aq11Key, ab11aq11Err := ab11aq11Decoder.Token()
				if ab11aq11Err != nil {
					return nil, ab11aq11Err
				}
				ab11aq11Name, ab11aq11OK := ab11aq11Key.(string)
				if !ab11aq11OK {
					return nil, ab11aq11fmt.Errorf("expected object key")
				}
				if _, ab11aq11Exists := ab11aq11Object[ab11aq11Name]; ab11aq11Exists {
					return nil, ab11aq11fmt.Errorf("duplicate key")
				}
				ab11aq11Value, ab11aq11Err := ab11aq11Read(ab11aq11Decoder)
				if ab11aq11Err != nil {
					return nil, ab11aq11Err
				}
				ab11aq11Object[ab11aq11Name] = ab11aq11Value
			}
			_, ab11aq11Err = ab11aq11Decoder.Token()
			return ab11aq11Object, ab11aq11Err
		case '[':
			ab11aq11Array := []any{}
			for ab11aq11Decoder.More() {
				ab11aq11Value, ab11aq11Err := ab11aq11Read(ab11aq11Decoder)
				if ab11aq11Err != nil {
					return nil, ab11aq11Err
				}
				ab11aq11Array = append(ab11aq11Array, ab11aq11Value)
			}
			_, ab11aq11Err = ab11aq11Decoder.Token()
			return ab11aq11Array, ab11aq11Err
		default:
			return nil, ab11aq11fmt.Errorf("unexpected delimiter")
		}
	}
	if ab11aq11Number, ab11aq11OK := ab11aq11Token.(ab11aq11json.Number); ab11aq11OK {
		if _, ab11aq11Err := ab11aq11strconv.ParseFloat(string(ab11aq11Number), 64); ab11aq11Err != nil {
			return nil, ab11aq11fmt.Errorf("nonfinite number")
		}
	}
	return ab11aq11Token, nil
}

func ab11aq11Integer(ab11aq11Value any, ab11aq11Max int64) (int64, bool) {
	ab11aq11Number, ab11aq11OK := ab11aq11Value.(ab11aq11json.Number)
	if !ab11aq11OK {
		return 0, false
	}
	ab11aq11Rational, ab11aq11OK := new(ab11aq11big.Rat).SetString(string(ab11aq11Number))
	if !ab11aq11OK || !ab11aq11Rational.IsInt() || !ab11aq11Rational.Num().IsInt64() {
		return 0, false
	}
	ab11aq11Integer := ab11aq11Rational.Num().Int64()
	return ab11aq11Integer, ab11aq11Integer >= 0 && ab11aq11Integer <= ab11aq11Max
}

func ab11aq11Encoding(ab11aq11Data []byte) bool {
	if !ab11aq11utf8.Valid(ab11aq11Data) {
		return false
	}
	ab11aq11Inside := false
	ab11aq11Unit := func(ab11aq11At int) (int, bool) {
		if ab11aq11At+6 > len(ab11aq11Data) || ab11aq11Data[ab11aq11At] != '\\' || ab11aq11Data[ab11aq11At+1] != 'u' {
			return 0, false
		}
		ab11aq11Parsed, ab11aq11Err := ab11aq11strconv.ParseUint(string(ab11aq11Data[ab11aq11At+2:ab11aq11At+6]), 16, 32)
		return int(ab11aq11Parsed), ab11aq11Err == nil
	}
	for ab11aq11Index := 0; ab11aq11Index < len(ab11aq11Data); ab11aq11Index++ {
		switch {
		case !ab11aq11Inside:
			ab11aq11Inside = ab11aq11Data[ab11aq11Index] == '"'
		case ab11aq11Data[ab11aq11Index] == '"':
			ab11aq11Inside = false
		case ab11aq11Data[ab11aq11Index] == '\\':
			if ab11aq11Index+1 < len(ab11aq11Data) && ab11aq11Data[ab11aq11Index+1] == 'u' {
				ab11aq11Code, ab11aq11OK := ab11aq11Unit(ab11aq11Index)
				if !ab11aq11OK {
					return false
				}
				if ab11aq11Code >= 0xDC00 && ab11aq11Code <= 0xDFFF {
					return false
				}
				if ab11aq11Code >= 0xD800 && ab11aq11Code <= 0xDBFF {
					ab11aq11Low, ab11aq11OK := ab11aq11Unit(ab11aq11Index + 6)
					if !ab11aq11OK || ab11aq11Low < 0xDC00 || ab11aq11Low > 0xDFFF {
						return false
					}
					ab11aq11Index += 6
				}
				ab11aq11Index += 5
			} else {
				ab11aq11Index++
			}
		}
	}
	return true
}

func Ab11Aq11Evaluate(ab11aq11InputJSON ab11aq11json.RawMessage) (ab11aq11json.RawMessage, error) {
	if !ab11aq11Encoding(ab11aq11InputJSON) {
		return nil, ab11aq11fmt.Errorf("invalid string encoding")
	}
	ab11aq11Decoder := ab11aq11json.NewDecoder(ab11aq11bytes.NewReader(ab11aq11InputJSON))
	ab11aq11Decoder.UseNumber()
	ab11aq11Value, ab11aq11Err := ab11aq11Read(ab11aq11Decoder)
	if ab11aq11Err != nil {
		return nil, ab11aq11Err
	}
	if _, ab11aq11Err = ab11aq11Decoder.Token(); ab11aq11Err != ab11aq11io.EOF {
		return nil, ab11aq11fmt.Errorf("trailing JSON")
	}
	ab11aq11Object, ab11aq11OK := ab11aq11Value.(map[string]any)
	if !ab11aq11OK {
		return nil, ab11aq11fmt.Errorf("expected object")
	}
	ab11aq11Context, ab11aq11OK := ab11aq11Object["context"].(map[string]any)
	if !ab11aq11OK {
		return nil, ab11aq11fmt.Errorf("expected context object")
	}
	ab11aq11Fields, ab11aq11OK := ab11aq11Object["fields"].(map[string]any)
	if !ab11aq11OK {
		return nil, ab11aq11fmt.Errorf("expected fields object")
	}
	ab11aq11Arguments := map[string]any{}
	if ab11aq11Value, ab11aq11Present := ab11aq11Object["arguments"]; ab11aq11Present {
		ab11aq11Arguments, ab11aq11OK = ab11aq11Value.(map[string]any)
		if !ab11aq11OK {
			return nil, ab11aq11fmt.Errorf("expected arguments object")
		}
	}
	for _, ab11aq11Name := range []string{"primary", "secondary", "operation"} {
		if ab11aq11Value, ab11aq11Present := ab11aq11Context[ab11aq11Name]; ab11aq11Present {
			if _, ab11aq11OK := ab11aq11Integer(ab11aq11Value, 255); !ab11aq11OK {
				return nil, ab11aq11fmt.Errorf("invalid context %s", ab11aq11Name)
			}
		} else if ab11aq11Name != "operation" {
			return nil, ab11aq11fmt.Errorf("missing context %s", ab11aq11Name)
		}
	}
	for _, ab11aq11Name := range []string{"direction", "view"} {
		if _, ab11aq11OK := ab11aq11Context[ab11aq11Name].(string); !ab11aq11OK {
			return nil, ab11aq11fmt.Errorf("invalid context %s", ab11aq11Name)
		}
	}
	for ab11aq11Name, ab11aq11Value := range ab11aq11Fields {
		ab11aq11Array, ab11aq11OK := ab11aq11Value.([]any)
		if !ab11aq11OK {
			return nil, ab11aq11fmt.Errorf("invalid field %s", ab11aq11Name)
		}
		for _, ab11aq11Byte := range ab11aq11Array {
			if _, ab11aq11OK := ab11aq11Integer(ab11aq11Byte, 255); !ab11aq11OK {
				return nil, ab11aq11fmt.Errorf("invalid byte in %s", ab11aq11Name)
			}
		}
	}
	ab11aq11Numbers := map[string]int64{}
	for _, ab11aq11Name := range []string{"ready_state", "final_option", "received_length"} {
		if ab11aq11Value, ab11aq11Present := ab11aq11Arguments[ab11aq11Name]; ab11aq11Present {
			ab11aq11Max := int64(255)
			if ab11aq11Name == "received_length" {
				ab11aq11Max = 65535
			}
			ab11aq11Number, ab11aq11OK := ab11aq11Integer(ab11aq11Value, ab11aq11Max)
			if !ab11aq11OK {
				return nil, ab11aq11fmt.Errorf("invalid argument %s", ab11aq11Name)
			}
			ab11aq11Numbers[ab11aq11Name] = ab11aq11Number
		}
	}
	ab11aq11Primary, _ := ab11aq11Integer(ab11aq11Context["primary"], 255)
	ab11aq11Secondary, _ := ab11aq11Integer(ab11aq11Context["secondary"], 255)
	_, ab11aq11Operation := ab11aq11Context["operation"]
	ab11aq11Raw := map[string]any{}
	for ab11aq11Name, ab11aq11Value := range ab11aq11Object {
		ab11aq11Raw[ab11aq11Name] = ab11aq11Value
	}
	ab11aq11Object["raw"] = ab11aq11Raw
	ab11aq11Results := map[string]any{}
	ab11aq11Object["state"] = "Unknown"
	ab11aq11Object["physical"] = "Unknown"
	ab11aq11Object["qualified"] = false
	ab11aq11Object["results"] = ab11aq11Results
	if ab11aq11Primary != 181 || ab11aq11Secondary != 23 || ab11aq11Operation || ab11aq11Context["direction"] != "local" || ab11aq11Context["view"] != "adaptation-guard" {
		return ab11aq11json.Marshal(ab11aq11Object)
	}
	for _, ab11aq11Name := range []string{"ready_state", "final_option", "received_length"} {
		if _, ab11aq11Present := ab11aq11Numbers[ab11aq11Name]; !ab11aq11Present {
			return nil, ab11aq11fmt.Errorf("missing argument %s", ab11aq11Name)
		}
	}
	ab11aq11Request, _ := ab11aq11Fields["request"].([]any)
	ab11aq11Input, _ := ab11aq11Fields["input"].([]any)
	_, ab11aq11InputPresent := ab11aq11Fields["input"]
	ab11aq11RequestComplete := len(ab11aq11Request) >= 2
	ab11aq11InputComplete := len(ab11aq11Input) >= 5
	ab11aq11Object["state"] = "complete"
	if !ab11aq11RequestComplete || !ab11aq11InputComplete {
		ab11aq11Object["state"] = "partial"
	}
	ab11aq11Eligible := true
	ab11aq11Add := func(ab11aq11Name string, ab11aq11Raw any, ab11aq11Complete bool, ab11aq11Flag bool) {
		ab11aq11State := "complete"
		if !ab11aq11Complete {
			ab11aq11State = "incomplete"
		}
		ab11aq11Results[ab11aq11Name] = map[string]any{"raw": ab11aq11Raw, "state": ab11aq11State, "grade": "Hypothesis", "flags": map[string]bool{"value": ab11aq11Complete && ab11aq11Flag}}
		if ab11aq11Complete {
			ab11aq11Results[ab11aq11Name].(map[string]any)["calculation"] = "Proven"
		}
		ab11aq11Eligible = ab11aq11Eligible && ab11aq11Complete && ab11aq11Flag
	}
	ab11aq11Add("ready_state", ab11aq11Arguments["ready_state"], true, ab11aq11Numbers["ready_state"] == 3)
	ab11aq11Add("final_option", ab11aq11Arguments["final_option"], true, ab11aq11Numbers["final_option"] != 0)
	ab11aq11Add("initialized_input", ab11aq11Fields, ab11aq11InputComplete, ab11aq11InputComplete)
	ab11aq11PayloadLength := false
	if ab11aq11InputComplete {
		ab11aq11Byte, _ := ab11aq11Integer(ab11aq11Input[4], 255)
		ab11aq11PayloadLength = ab11aq11Byte+6 == ab11aq11Numbers["received_length"]
	}
	ab11aq11Add("payload_length", map[string]any{"fields": ab11aq11Fields, "received_length": ab11aq11Arguments["received_length"]}, ab11aq11InputComplete, ab11aq11PayloadLength)
	ab11aq11Add("initialized_length", map[string]any{"fields": ab11aq11Fields, "received_length": ab11aq11Arguments["received_length"]}, ab11aq11InputPresent, ab11aq11Numbers["received_length"] == int64(len(ab11aq11Input)))
	ab11aq11RequestIdentity := false
	if ab11aq11RequestComplete {
		ab11aq11First, _ := ab11aq11Integer(ab11aq11Request[0], 255)
		ab11aq11Second, _ := ab11aq11Integer(ab11aq11Request[1], 255)
		ab11aq11RequestIdentity = ab11aq11First == 181 && ab11aq11Second == 23
	}
	ab11aq11Add("request_identity", ab11aq11Fields, ab11aq11RequestComplete, ab11aq11RequestIdentity)
	ab11aq11InputIdentity := false
	if len(ab11aq11Input) >= 2 {
		ab11aq11First, _ := ab11aq11Integer(ab11aq11Input[0], 255)
		ab11aq11Second, _ := ab11aq11Integer(ab11aq11Input[1], 255)
		ab11aq11InputIdentity = ab11aq11First == 181 && ab11aq11Second == 24
	}
	ab11aq11Add("input_identity", ab11aq11Fields, len(ab11aq11Input) >= 2, ab11aq11InputIdentity)
	ab11aq11State := "complete"
	if ab11aq11Object["state"] == "partial" {
		ab11aq11State = "incomplete"
	}
	ab11aq11Results["eligible"] = map[string]any{"raw": map[string]any{"fields": ab11aq11Fields, "arguments": ab11aq11Arguments}, "state": ab11aq11State, "grade": "Hypothesis", "flags": map[string]bool{"value": ab11aq11Eligible}}
	if ab11aq11State == "complete" {
		ab11aq11Results["eligible"].(map[string]any)["calculation"] = "Proven"
	}
	return ab11aq11json.Marshal(ab11aq11Object)
}
