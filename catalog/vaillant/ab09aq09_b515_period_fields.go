package vaillant

import (
	ab09aq09bytes "bytes"
	ab09aq09json "encoding/json"
	ab09aq09fmt "fmt"
	ab09aq09io "io"
	ab09aq09math "math"
	ab09aq09strconv "strconv"
)

func ab09aq09Decode(ab09aq09Decoder *ab09aq09json.Decoder) (any, error) {
	ab09aq09Token, ab09aq09Error := ab09aq09Decoder.Token()
	if ab09aq09Error != nil {
		return nil, ab09aq09Error
	}
	if ab09aq09Delimiter, ab09aq09OK := ab09aq09Token.(ab09aq09json.Delim); ab09aq09OK {
		switch ab09aq09Delimiter {
		case '{':
			ab09aq09Object := map[string]any{}
			for ab09aq09Decoder.More() {
				ab09aq09KeyToken, ab09aq09Error := ab09aq09Decoder.Token()
				if ab09aq09Error != nil {
					return nil, ab09aq09Error
				}
				ab09aq09Key, ab09aq09OK := ab09aq09KeyToken.(string)
				if !ab09aq09OK {
					return nil, ab09aq09fmt.Errorf("object key required")
				}
				if _, ab09aq09Present := ab09aq09Object[ab09aq09Key]; ab09aq09Present {
					return nil, ab09aq09fmt.Errorf("duplicate key %q", ab09aq09Key)
				}
				ab09aq09Value, ab09aq09Error := ab09aq09Decode(ab09aq09Decoder)
				if ab09aq09Error != nil {
					return nil, ab09aq09Error
				}
				ab09aq09Object[ab09aq09Key] = ab09aq09Value
			}
			_, ab09aq09Error = ab09aq09Decoder.Token()
			return ab09aq09Object, ab09aq09Error
		case '[':
			ab09aq09Array := []any{}
			for ab09aq09Decoder.More() {
				ab09aq09Value, ab09aq09Error := ab09aq09Decode(ab09aq09Decoder)
				if ab09aq09Error != nil {
					return nil, ab09aq09Error
				}
				ab09aq09Array = append(ab09aq09Array, ab09aq09Value)
			}
			_, ab09aq09Error = ab09aq09Decoder.Token()
			return ab09aq09Array, ab09aq09Error
		}
		return nil, ab09aq09fmt.Errorf("unexpected delimiter")
	}
	if ab09aq09Number, ab09aq09OK := ab09aq09Token.(ab09aq09json.Number); ab09aq09OK {
		ab09aq09Float, ab09aq09Error := ab09aq09strconv.ParseFloat(string(ab09aq09Number), 64)
		if ab09aq09Error != nil || ab09aq09math.IsInf(ab09aq09Float, 0) || ab09aq09math.IsNaN(ab09aq09Float) {
			return nil, ab09aq09fmt.Errorf("nonfinite number")
		}
	}
	return ab09aq09Token, nil
}

func ab09aq09Integer(ab09aq09Value any, ab09aq09Maximum int) (int, error) {
	ab09aq09Number, ab09aq09OK := ab09aq09Value.(ab09aq09json.Number)
	if !ab09aq09OK {
		return 0, ab09aq09fmt.Errorf("integer required")
	}
	ab09aq09Integer, ab09aq09Error := ab09aq09strconv.ParseUint(string(ab09aq09Number), 10, 16)
	if ab09aq09Error != nil || ab09aq09Integer > uint64(ab09aq09Maximum) {
		return 0, ab09aq09fmt.Errorf("integer out of range or nonintegral")
	}
	return int(ab09aq09Integer), nil
}

func ab09aq09Period(ab09aq09Raw any) map[string]any {
	ab09aq09Result := map[string]any{"raw": ab09aq09Raw, "state": "Unknown"}
	ab09aq09Number, _ := ab09aq09Integer(ab09aq09Raw, 255)
	if ab09aq09Number <= 144 {
		ab09aq09Minutes := ab09aq09Number * 10
		ab09aq09Result["state"] = "Hypothesis"
		ab09aq09Result["calculation"] = "Proven"
		ab09aq09Result["number"] = ab09aq09Minutes
		ab09aq09Result["hour"] = ab09aq09Minutes / 60
		ab09aq09Result["minute"] = ab09aq09Minutes % 60
		ab09aq09Result["label"] = ab09aq09fmt.Sprintf("%02d:%02d", ab09aq09Minutes/60, ab09aq09Minutes%60)
	}
	return ab09aq09Result
}

func Ab09Aq09Evaluate(ab09aq09Input ab09aq09json.RawMessage) (ab09aq09json.RawMessage, error) {
	ab09aq09Decoder := ab09aq09json.NewDecoder(ab09aq09bytes.NewReader(ab09aq09Input))
	ab09aq09Decoder.UseNumber()
	ab09aq09Value, ab09aq09Error := ab09aq09Decode(ab09aq09Decoder)
	if ab09aq09Error != nil {
		return nil, ab09aq09Error
	}
	if _, ab09aq09Error = ab09aq09Decoder.Token(); ab09aq09Error != ab09aq09io.EOF {
		return nil, ab09aq09fmt.Errorf("trailing JSON")
	}
	ab09aq09Object, ab09aq09OK := ab09aq09Value.(map[string]any)
	if !ab09aq09OK {
		return nil, ab09aq09fmt.Errorf("input object required")
	}
	ab09aq09Context, ab09aq09OK := ab09aq09Object["context"].(map[string]any)
	if !ab09aq09OK {
		return nil, ab09aq09fmt.Errorf("context object required")
	}
	ab09aq09Fields, ab09aq09OK := ab09aq09Object["fields"].(map[string]any)
	if !ab09aq09OK {
		return nil, ab09aq09fmt.Errorf("fields object required")
	}
	for ab09aq09Key, ab09aq09Value := range ab09aq09Fields {
		ab09aq09Array, ab09aq09OK := ab09aq09Value.([]any)
		if !ab09aq09OK {
			return nil, ab09aq09fmt.Errorf("field %s requires byte array", ab09aq09Key)
		}
		for _, ab09aq09Byte := range ab09aq09Array {
			if _, ab09aq09Error := ab09aq09Integer(ab09aq09Byte, 255); ab09aq09Error != nil {
				return nil, ab09aq09Error
			}
		}
	}
	ab09aq09Supported := true
	for ab09aq09Key, ab09aq09Expected := range map[string]int{"primary": 181, "secondary": 21} {
		ab09aq09Value, ab09aq09Present := ab09aq09Context[ab09aq09Key]
		if !ab09aq09Present {
			ab09aq09Supported = false
			continue
		}
		ab09aq09Number, ab09aq09Error := ab09aq09Integer(ab09aq09Value, 255)
		if ab09aq09Error != nil {
			return nil, ab09aq09Error
		}
		if ab09aq09Number != ab09aq09Expected {
			ab09aq09Supported = false
		}
	}
	if ab09aq09Operation, ab09aq09Present := ab09aq09Context["operation"]; ab09aq09Present {
		if _, ab09aq09Error := ab09aq09Integer(ab09aq09Operation, 255); ab09aq09Error != nil {
			return nil, ab09aq09Error
		}
		ab09aq09Supported = false
	}
	for _, ab09aq09Key := range []string{"direction", "view"} {
		if ab09aq09Value, ab09aq09Present := ab09aq09Context[ab09aq09Key]; ab09aq09Present {
			if _, ab09aq09OK := ab09aq09Value.(string); !ab09aq09OK {
				return nil, ab09aq09fmt.Errorf("context %s requires string", ab09aq09Key)
			}
		}
	}
	ab09aq09View, _ := ab09aq09Context["view"].(string)
	ab09aq09Direction, _ := ab09aq09Context["direction"].(string)
	ab09aq09Supported = ab09aq09Supported && ((ab09aq09View == "periods" && (ab09aq09Direction == "request" || ab09aq09Direction == "positive")) || (ab09aq09View == "intervals" && ab09aq09Direction == "request"))
	ab09aq09Arguments := map[string]any{}
	if ab09aq09Value, ab09aq09Present := ab09aq09Object["arguments"]; ab09aq09Present {
		ab09aq09Arguments, ab09aq09OK = ab09aq09Value.(map[string]any)
		if !ab09aq09OK {
			return nil, ab09aq09fmt.Errorf("arguments object required")
		}
	}
	ab09aq09Minutes := 0
	if ab09aq09Value, ab09aq09Present := ab09aq09Arguments["minutes"]; ab09aq09Present {
		ab09aq09Minutes, ab09aq09Error = ab09aq09Integer(ab09aq09Value, 1440)
		if ab09aq09Error != nil {
			return nil, ab09aq09Error
		}
	}
	ab09aq09Pairs := [3][2]int{}
	if ab09aq09Value, ab09aq09Present := ab09aq09Arguments["pairs"]; ab09aq09Present {
		ab09aq09Array, ab09aq09OK := ab09aq09Value.([]any)
		if !ab09aq09OK || len(ab09aq09Array) != 3 {
			return nil, ab09aq09fmt.Errorf("three pairs required")
		}
		for ab09aq09Index, ab09aq09Value := range ab09aq09Array {
			ab09aq09Pair, ab09aq09OK := ab09aq09Value.([]any)
			if !ab09aq09OK || len(ab09aq09Pair) != 2 {
				return nil, ab09aq09fmt.Errorf("two endpoints required")
			}
			for ab09aq09Endpoint, ab09aq09Value := range ab09aq09Pair {
				ab09aq09Pairs[ab09aq09Index][ab09aq09Endpoint], ab09aq09Error = ab09aq09Integer(ab09aq09Value, 255)
				if ab09aq09Error != nil {
					return nil, ab09aq09Error
				}
			}
		}
	} else if ab09aq09Supported && ab09aq09View == "intervals" {
		return nil, ab09aq09fmt.Errorf("pairs required")
	}
	ab09aq09Results := map[string]any{}
	ab09aq09Extra := map[string]any{}
	for ab09aq09Key, ab09aq09Value := range ab09aq09Object {
		if ab09aq09Key != "context" && ab09aq09Key != "fields" && ab09aq09Key != "arguments" {
			ab09aq09Extra[ab09aq09Key] = ab09aq09Value
		}
	}
	ab09aq09Output := map[string]any{"context": ab09aq09Context, "fields": ab09aq09Fields, "extra": ab09aq09Extra, "results": ab09aq09Results, "state": "Unknown", "physical": "Unknown", "qualified": false}
	if _, ab09aq09Present := ab09aq09Object["arguments"]; ab09aq09Present {
		ab09aq09Output["arguments"] = ab09aq09Arguments
	}
	if !ab09aq09Supported {
		return ab09aq09json.Marshal(ab09aq09Output)
	}
	ab09aq09Output["state"] = "complete"
	for ab09aq09Key, ab09aq09Value := range ab09aq09Fields {
		ab09aq09Results[ab09aq09Key] = map[string]any{"raw": ab09aq09Value, "state": "Unknown"}
	}
	if ab09aq09View == "periods" {
		for _, ab09aq09Key := range []string{"period0", "period1", "period2"} {
			ab09aq09Result := map[string]any{"state": "incomplete"}
			ab09aq09Value, ab09aq09Present := ab09aq09Fields[ab09aq09Key]
			if ab09aq09Present {
				ab09aq09Result["raw"] = ab09aq09Value
			}
			ab09aq09Array, _ := ab09aq09Value.([]any)
			if len(ab09aq09Array) == 1 {
				ab09aq09Result = ab09aq09Period(ab09aq09Array[0])
				ab09aq09Result["raw"] = ab09aq09Value
			} else {
				ab09aq09Output["state"] = "partial"
			}
			ab09aq09Results[ab09aq09Key] = ab09aq09Result
		}
		if ab09aq09Value, ab09aq09Present := ab09aq09Arguments["minutes"]; ab09aq09Present {
			ab09aq09Result := map[string]any{"raw": ab09aq09Value, "state": "Unknown"}
			if ab09aq09Minutes%10 == 0 {
				ab09aq09Result["state"] = "Hypothesis"
				ab09aq09Result["calculation"] = "Proven"
				ab09aq09Result["number"] = ab09aq09Minutes / 10
			}
			ab09aq09Results["minutes"] = ab09aq09Result
		}
	} else {
		ab09aq09PairResults := []any{}
		ab09aq09Mapped := true
		for ab09aq09Index, ab09aq09Pair := range ab09aq09Pairs {
			ab09aq09Label := "enabled-candidate"
			if ab09aq09Pair[0] == ab09aq09Pair[1] {
				ab09aq09Label = "disabled"
			}
			ab09aq09PairResults = append(ab09aq09PairResults, map[string]any{"raw": ab09aq09Arguments["pairs"].([]any)[ab09aq09Index], "state": "Hypothesis", "calculation": "Proven", "label": ab09aq09Label, "start": ab09aq09Pair[0], "end": ab09aq09Pair[1]})
			if ab09aq09Pair[0] > 144 || ab09aq09Pair[1] > 144 {
				ab09aq09Mapped = false
			}
		}
		ab09aq09Results["pairs"] = map[string]any{"raw": ab09aq09Arguments["pairs"], "state": "Hypothesis", "alternatives": ab09aq09PairResults}
		ab09aq09Order := map[string]any{"raw": ab09aq09Arguments["pairs"], "state": "Unknown"}
		if ab09aq09Mapped {
			ab09aq09Order["state"] = "Hypothesis"
			ab09aq09Order["calculation"] = "Proven"
			ab09aq09Order["value"] = ab09aq09Pairs[0][0] < ab09aq09Pairs[1][0] && ab09aq09Pairs[1][0] < ab09aq09Pairs[2][0]
		}
		ab09aq09Results["start_order"] = ab09aq09Order
	}
	return ab09aq09json.Marshal(ab09aq09Output)
}
