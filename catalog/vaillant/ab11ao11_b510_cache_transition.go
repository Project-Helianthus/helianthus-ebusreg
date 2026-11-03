package vaillant

import (
	ab11ao11bytes "bytes"
	ab11ao11json "encoding/json"
	ab11ao11fmt "fmt"
	ab11ao11io "io"
	ab11ao11big "math/big"
)

func ab11ao11Read(ab11ao11Decoder *ab11ao11json.Decoder) (any, error) {
	ab11ao11Token, ab11ao11Err := ab11ao11Decoder.Token()
	if ab11ao11Err != nil {
		return nil, ab11ao11Err
	}
	if ab11ao11Delimiter, ab11ao11OK := ab11ao11Token.(ab11ao11json.Delim); ab11ao11OK {
		switch ab11ao11Delimiter {
		case '{':
			ab11ao11Object := map[string]any{}
			for ab11ao11Decoder.More() {
				ab11ao11Key, ab11ao11Err := ab11ao11Decoder.Token()
				if ab11ao11Err != nil {
					return nil, ab11ao11Err
				}
				ab11ao11Name, ab11ao11OK := ab11ao11Key.(string)
				if !ab11ao11OK {
					return nil, ab11ao11fmt.Errorf("expected object key")
				}
				if _, ab11ao11Exists := ab11ao11Object[ab11ao11Name]; ab11ao11Exists {
					return nil, ab11ao11fmt.Errorf("duplicate object key")
				}
				ab11ao11Value, ab11ao11Err := ab11ao11Read(ab11ao11Decoder)
				if ab11ao11Err != nil {
					return nil, ab11ao11Err
				}
				ab11ao11Object[ab11ao11Name] = ab11ao11Value
			}
			_, ab11ao11Err = ab11ao11Decoder.Token()
			return ab11ao11Object, ab11ao11Err
		case '[':
			ab11ao11Array := []any{}
			for ab11ao11Decoder.More() {
				ab11ao11Value, ab11ao11Err := ab11ao11Read(ab11ao11Decoder)
				if ab11ao11Err != nil {
					return nil, ab11ao11Err
				}
				ab11ao11Array = append(ab11ao11Array, ab11ao11Value)
			}
			_, ab11ao11Err = ab11ao11Decoder.Token()
			return ab11ao11Array, ab11ao11Err
		default:
			return nil, ab11ao11fmt.Errorf("unexpected delimiter")
		}
	}
	return ab11ao11Token, nil
}

func ab11ao11Integer(ab11ao11Value any) (*ab11ao11big.Int, bool) {
	ab11ao11Number, ab11ao11OK := ab11ao11Value.(ab11ao11json.Number)
	if !ab11ao11OK {
		return nil, false
	}
	ab11ao11Result, ab11ao11OK := new(ab11ao11big.Int).SetString(string(ab11ao11Number), 10)
	return ab11ao11Result, ab11ao11OK
}

func Ab11Ao11Interpret(ab11ao11Input ab11ao11json.RawMessage) (ab11ao11json.RawMessage, error) {
	ab11ao11Decoder := ab11ao11json.NewDecoder(ab11ao11bytes.NewReader(ab11ao11Input))
	ab11ao11Decoder.UseNumber()
	ab11ao11Value, ab11ao11Err := ab11ao11Read(ab11ao11Decoder)
	if ab11ao11Err != nil {
		return nil, ab11ao11Err
	}
	if _, ab11ao11Err = ab11ao11Decoder.Token(); ab11ao11Err != ab11ao11io.EOF {
		return nil, ab11ao11fmt.Errorf("trailing JSON")
	}
	ab11ao11Object, ab11ao11OK := ab11ao11Value.(map[string]any)
	if !ab11ao11OK {
		return nil, ab11ao11fmt.Errorf("expected input object")
	}
	ab11ao11Context := map[string]any{}
	if ab11ao11Value, ab11ao11Present := ab11ao11Object["context"]; ab11ao11Present {
		ab11ao11Context, ab11ao11OK = ab11ao11Value.(map[string]any)
		if !ab11ao11OK {
			return nil, ab11ao11fmt.Errorf("expected context object")
		}
	}
	ab11ao11Raw := map[string]any{}
	if ab11ao11Value, ab11ao11Present := ab11ao11Object["fields"]; ab11ao11Present {
		ab11ao11Raw, ab11ao11OK = ab11ao11Value.(map[string]any)
		if !ab11ao11OK {
			return nil, ab11ao11fmt.Errorf("expected fields object")
		}
	}
	for _, ab11ao11Value := range ab11ao11Raw {
		ab11ao11Array, ab11ao11OK := ab11ao11Value.([]any)
		if !ab11ao11OK {
			return nil, ab11ao11fmt.Errorf("expected byte array")
		}
		for _, ab11ao11Element := range ab11ao11Array {
			ab11ao11Byte, ab11ao11OK := ab11ao11Integer(ab11ao11Element)
			if !ab11ao11OK || ab11ao11Byte.Sign() < 0 || ab11ao11Byte.Cmp(ab11ao11big.NewInt(255)) > 0 {
				return nil, ab11ao11fmt.Errorf("invalid byte element")
			}
		}
	}
	ab11ao11Selected := true
	for _, ab11ao11Axis := range []struct {
		ab11ao11Name     string
		ab11ao11Expected int64
	}{{"primary", 181}, {"secondary", 16}} {
		ab11ao11Value, ab11ao11Present := ab11ao11Context[ab11ao11Axis.ab11ao11Name]
		if !ab11ao11Present {
			ab11ao11Selected = false
			continue
		}
		ab11ao11Number, ab11ao11OK := ab11ao11Integer(ab11ao11Value)
		if !ab11ao11OK {
			return nil, ab11ao11fmt.Errorf("expected context integer")
		}
		if ab11ao11Number.Cmp(ab11ao11big.NewInt(ab11ao11Axis.ab11ao11Expected)) != 0 {
			ab11ao11Selected = false
		}
	}
	if ab11ao11Value, ab11ao11Present := ab11ao11Context["operation"]; ab11ao11Present {
		if _, ab11ao11OK := ab11ao11Integer(ab11ao11Value); !ab11ao11OK {
			return nil, ab11ao11fmt.Errorf("expected operation integer")
		}
		ab11ao11Selected = false
	}
	for _, ab11ao11Axis := range []struct{ ab11ao11Name, ab11ao11Expected string }{{"direction", "local"}, {"view", "capability-query"}} {
		ab11ao11Value, ab11ao11Present := ab11ao11Context[ab11ao11Axis.ab11ao11Name]
		if !ab11ao11Present {
			ab11ao11Selected = false
			continue
		}
		ab11ao11Text, ab11ao11OK := ab11ao11Value.(string)
		if !ab11ao11OK {
			return nil, ab11ao11fmt.Errorf("expected context string")
		}
		if ab11ao11Text != ab11ao11Axis.ab11ao11Expected {
			ab11ao11Selected = false
		}
	}
	ab11ao11Arguments := map[string]any{}
	ab11ao11ArgumentValue, ab11ao11ArgumentsPresent := ab11ao11Object["arguments"]
	if ab11ao11ArgumentsPresent {
		ab11ao11Arguments, ab11ao11OK = ab11ao11ArgumentValue.(map[string]any)
		if !ab11ao11OK {
			return nil, ab11ao11fmt.Errorf("expected arguments object")
		}
	}
	for ab11ao11Name, ab11ao11Value := range ab11ao11Arguments {
		switch ab11ao11Name {
		case "parsed":
			if _, ab11ao11OK := ab11ao11Value.(bool); !ab11ao11OK {
				return nil, ab11ao11fmt.Errorf("expected parsed boolean")
			}
		case "count":
			if _, ab11ao11OK := ab11ao11Integer(ab11ao11Value); !ab11ao11OK {
				return nil, ab11ao11fmt.Errorf("expected count integer")
			}
		default:
			return nil, ab11ao11fmt.Errorf("unknown argument")
		}
	}
	ab11ao11Fields := map[string]any{}
	ab11ao11Output := map[string]any{"context": ab11ao11Context, "raw": ab11ao11Raw, "state": "Unknown", "fields": ab11ao11Fields, "physical": "Unknown", "qualified": false}
	if ab11ao11ArgumentsPresent {
		ab11ao11Output["arguments"] = ab11ao11Arguments
	}
	ab11ao11Bytes := map[string]int64{}
	for _, ab11ao11Name := range []string{"prior", "feedback", "bit_six"} {
		ab11ao11Array, ab11ao11OK := ab11ao11Raw[ab11ao11Name].([]any)
		if !ab11ao11OK || len(ab11ao11Array) != 1 {
			ab11ao11Selected = false
			continue
		}
		ab11ao11Number, _ := ab11ao11Integer(ab11ao11Array[0])
		ab11ao11Bytes[ab11ao11Name] = ab11ao11Number.Int64()
	}
	if ab11ao11Bytes["bit_six"] != 0 && ab11ao11Bytes["bit_six"] != 64 {
		ab11ao11Selected = false
	}
	ab11ao11Parsed, ab11ao11HasParsed := ab11ao11Arguments["parsed"].(bool)
	ab11ao11Count, ab11ao11HasCount := ab11ao11Integer(ab11ao11Arguments["count"])
	if !ab11ao11HasParsed || !ab11ao11HasCount || ab11ao11Count.Sign() < 0 {
		ab11ao11Selected = false
	}
	if !ab11ao11Selected {
		return ab11ao11json.Marshal(ab11ao11Output)
	}
	ab11ao11Cell := ab11ao11Bytes["prior"]
	ab11ao11Freshness, ab11ao11Rule, ab11ao11RuleState := "Unknown", "retain", "Hypothesis"
	if ab11ao11Parsed && ab11ao11Bytes["feedback"] == 1 {
		if ab11ao11Count.Sign() == 0 {
			ab11ao11Rule, ab11ao11RuleState = "Unknown", "Unknown"
		} else {
			ab11ao11Freshness = "fresh"
			if ab11ao11Count.Cmp(ab11ao11big.NewInt(2)) <= 0 {
				ab11ao11Cell, ab11ao11Rule = 0, "replace-zero"
			} else {
				ab11ao11Cell, ab11ao11Rule = ab11ao11Bytes["bit_six"], "replace-bit-six"
			}
		}
	}
	ab11ao11Fields["cache"] = map[string]any{"raw": ab11ao11Raw["prior"], "state": "Hypothesis", "number": ab11ao11Cell}
	ab11ao11Fields["freshness"] = map[string]any{"raw": ab11ao11Raw["feedback"], "state": "Hypothesis", "label": ab11ao11Freshness}
	ab11ao11Fields["rule"] = map[string]any{"raw": ab11ao11Raw["bit_six"], "state": ab11ao11RuleState, "label": ab11ao11Rule}
	if ab11ao11RuleState != "Unknown" {
		ab11ao11Output["state"] = "Hypothesis"
	}
	return ab11ao11json.Marshal(ab11ao11Output)
}
