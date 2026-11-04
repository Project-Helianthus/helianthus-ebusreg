package vaillant

import (
	ab13ao13bytes "bytes"
	ab13ao13json "encoding/json"
	ab13ao13fmt "fmt"
	ab13ao13io "io"
	ab13ao13strconv "strconv"
)

func ab13ao13Read(ab13ao13Decoder *ab13ao13json.Decoder) (any, error) {
	ab13ao13Token, ab13ao13Err := ab13ao13Decoder.Token()
	if ab13ao13Err != nil {
		return nil, ab13ao13Err
	}
	if ab13ao13Delimiter, ab13ao13OK := ab13ao13Token.(ab13ao13json.Delim); ab13ao13OK {
		switch ab13ao13Delimiter {
		case '{':
			ab13ao13Object := map[string]any{}
			for ab13ao13Decoder.More() {
				ab13ao13KeyToken, ab13ao13Err := ab13ao13Decoder.Token()
				if ab13ao13Err != nil {
					return nil, ab13ao13Err
				}
				ab13ao13Key := ab13ao13KeyToken.(string)
				if _, ab13ao13Exists := ab13ao13Object[ab13ao13Key]; ab13ao13Exists {
					return nil, ab13ao13fmt.Errorf("duplicate key %q", ab13ao13Key)
				}
				ab13ao13Value, ab13ao13Err := ab13ao13Read(ab13ao13Decoder)
				if ab13ao13Err != nil {
					return nil, ab13ao13Err
				}
				ab13ao13Object[ab13ao13Key] = ab13ao13Value
			}
			_, ab13ao13Err := ab13ao13Decoder.Token()
			return ab13ao13Object, ab13ao13Err
		case '[':
			ab13ao13Array := []any{}
			for ab13ao13Decoder.More() {
				ab13ao13Value, ab13ao13Err := ab13ao13Read(ab13ao13Decoder)
				if ab13ao13Err != nil {
					return nil, ab13ao13Err
				}
				ab13ao13Array = append(ab13ao13Array, ab13ao13Value)
			}
			_, ab13ao13Err := ab13ao13Decoder.Token()
			return ab13ao13Array, ab13ao13Err
		}
		return nil, ab13ao13fmt.Errorf("unexpected delimiter")
	}
	return ab13ao13Token, nil
}

func ab13ao13Integer(ab13ao13Value any) (int64, bool) {
	ab13ao13Number, ab13ao13OK := ab13ao13Value.(ab13ao13json.Number)
	if !ab13ao13OK {
		return 0, false
	}
	ab13ao13Int, ab13ao13Err := ab13ao13strconv.ParseInt(string(ab13ao13Number), 10, 64)
	return ab13ao13Int, ab13ao13Err == nil
}

func Ab13Ao13Interpret(ab13ao13Input ab13ao13json.RawMessage) (ab13ao13json.RawMessage, error) {
	ab13ao13Decoder := ab13ao13json.NewDecoder(ab13ao13bytes.NewReader(ab13ao13Input))
	ab13ao13Decoder.UseNumber()
	ab13ao13Value, ab13ao13Err := ab13ao13Read(ab13ao13Decoder)
	if ab13ao13Err != nil {
		return nil, ab13ao13Err
	}
	if _, ab13ao13Err = ab13ao13Decoder.Token(); ab13ao13Err != ab13ao13io.EOF {
		return nil, ab13ao13fmt.Errorf("trailing JSON")
	}
	ab13ao13InputObject, ab13ao13OK := ab13ao13Value.(map[string]any)
	if !ab13ao13OK {
		return nil, ab13ao13fmt.Errorf("input must be object")
	}
	ab13ao13Context, ab13ao13OK := ab13ao13InputObject["context"].(map[string]any)
	if !ab13ao13OK {
		return nil, ab13ao13fmt.Errorf("context must be object")
	}
	ab13ao13Raw, ab13ao13OK := ab13ao13InputObject["fields"].(map[string]any)
	if !ab13ao13OK {
		return nil, ab13ao13fmt.Errorf("fields must be object")
	}
	ab13ao13Bytes := map[string][]int64{}
	for ab13ao13Key, ab13ao13Value := range ab13ao13Raw {
		ab13ao13Array, ab13ao13OK := ab13ao13Value.([]any)
		if !ab13ao13OK {
			return nil, ab13ao13fmt.Errorf("field %s must be byte array", ab13ao13Key)
		}
		ab13ao13Bytes[ab13ao13Key] = []int64{}
		for _, ab13ao13Value := range ab13ao13Array {
			ab13ao13Int, ab13ao13OK := ab13ao13Integer(ab13ao13Value)
			if !ab13ao13OK || ab13ao13Int < 0 || ab13ao13Int > 255 {
				return nil, ab13ao13fmt.Errorf("invalid byte")
			}
			ab13ao13Bytes[ab13ao13Key] = append(ab13ao13Bytes[ab13ao13Key], ab13ao13Int)
		}
	}
	for _, ab13ao13Key := range []string{"primary", "secondary", "operation"} {
		if ab13ao13Value, ab13ao13Exists := ab13ao13Context[ab13ao13Key]; ab13ao13Exists {
			if _, ab13ao13OK := ab13ao13Integer(ab13ao13Value); !ab13ao13OK {
				return nil, ab13ao13fmt.Errorf("invalid context integer")
			}
		}
	}
	for _, ab13ao13Key := range []string{"direction", "view"} {
		if ab13ao13Value, ab13ao13Exists := ab13ao13Context[ab13ao13Key]; ab13ao13Exists {
			if _, ab13ao13OK := ab13ao13Value.(string); !ab13ao13OK {
				return nil, ab13ao13fmt.Errorf("invalid context string")
			}
		}
	}
	ab13ao13Arguments := map[string]any{}
	ab13ao13ArgumentValue, ab13ao13HasArguments := ab13ao13InputObject["arguments"]
	if ab13ao13HasArguments {
		ab13ao13Arguments, ab13ao13OK = ab13ao13ArgumentValue.(map[string]any)
		if !ab13ao13OK {
			return nil, ab13ao13fmt.Errorf("arguments must be object")
		}
	}
	for ab13ao13Key, ab13ao13Value := range ab13ao13Arguments {
		switch ab13ao13Key {
		case "parsed":
			if _, ab13ao13OK := ab13ao13Value.(bool); !ab13ao13OK {
				return nil, ab13ao13fmt.Errorf("parsed must be boolean")
			}
		case "index", "selector", "local_selector", "count":
			if ab13ao13Int, ab13ao13OK := ab13ao13Integer(ab13ao13Value); !ab13ao13OK || ab13ao13Int < 0 {
				return nil, ab13ao13fmt.Errorf("invalid argument integer")
			}
		default:
			return nil, ab13ao13fmt.Errorf("unknown argument")
		}
	}
	ab13ao13Fields := map[string]any{}
	ab13ao13Output := map[string]any{"context": ab13ao13Context, "raw": ab13ao13Raw, "state": "Unknown", "fields": ab13ao13Fields, "physical": "Unknown", "qualified": false}
	if ab13ao13HasArguments {
		ab13ao13Output["arguments"] = ab13ao13Arguments
	}
	ab13ao13Finish := func() (ab13ao13json.RawMessage, error) { return ab13ao13json.Marshal(ab13ao13Output) }
	ab13ao13Primary, _ := ab13ao13Integer(ab13ao13Context["primary"])
	ab13ao13Secondary, _ := ab13ao13Integer(ab13ao13Context["secondary"])
	if ab13ao13Primary != 181 || ab13ao13Secondary != 17 || ab13ao13Context["direction"] != "local" {
		return ab13ao13Finish()
	}
	ab13ao13View, _ := ab13ao13Context["view"].(string)
	ab13ao13Operation, ab13ao13HasOperation := ab13ao13Context["operation"]
	if ab13ao13View == "refresh" {
		ab13ao13Op, _ := ab13ao13Integer(ab13ao13Operation)
		if !ab13ao13HasOperation || ab13ao13Op != 1 {
			return ab13ao13Finish()
		}
	} else if ab13ao13HasOperation || (ab13ao13View != "initialize" && ab13ao13View != "timer-consume") {
		return ab13ao13Finish()
	}
	ab13ao13Output["state"] = "Hypothesis"
	if ab13ao13View == "initialize" {
		ab13ao13Fields["failure"] = map[string]any{"raw": []any{}, "state": "Hypothesis", "flags": map[string]bool{"failure": false}}
		return ab13ao13Finish()
	}
	ab13ao13Index, ab13ao13IndexOK := ab13ao13Integer(ab13ao13Arguments["index"])
	if ab13ao13View == "timer-consume" && (!ab13ao13IndexOK || ab13ao13Index != 0) {
		ab13ao13Output["state"] = "Unknown"
		return ab13ao13Finish()
	}
	ab13ao13Prior := ab13ao13Bytes["prior_failure"]
	ab13ao13Failure := false
	ab13ao13FailureField := map[string]any{}
	if len(ab13ao13Prior) == 1 {
		ab13ao13Failure = ab13ao13Prior[0] != 0
		ab13ao13FailureField = map[string]any{"raw": ab13ao13Raw["prior_failure"], "state": "Hypothesis", "flags": map[string]bool{"failure": ab13ao13Failure}}
		ab13ao13Fields["failure"] = ab13ao13FailureField
	} else {
		ab13ao13Output["state"] = "partial"
		if ab13ao13View == "timer-consume" {
			return ab13ao13Finish()
		}
	}
	if ab13ao13View == "timer-consume" {
		ab13ao13Path := map[string]any{"raw": ab13ao13Raw["prior_failure"], "state": "Hypothesis", "label": "capability-path"}
		if ab13ao13Failure {
			ab13ao13Path["label"] = "error"
			ab13ao13Path["number"] = 120
		}
		ab13ao13Fields["path"] = ab13ao13Path
		return ab13ao13Finish()
	}
	ab13ao13Registration := ab13ao13Bytes["registration"]
	ab13ao13Selector, ab13ao13SelectorOK := ab13ao13Integer(ab13ao13Arguments["selector"])
	ab13ao13LocalSelector, ab13ao13LocalOK := ab13ao13Integer(ab13ao13Arguments["local_selector"])
	ab13ao13Count, ab13ao13CountOK := ab13ao13Integer(ab13ao13Arguments["count"])
	ab13ao13Parsed, ab13ao13ParsedOK := ab13ao13Arguments["parsed"].(bool)
	if len(ab13ao13Registration) != 1 || !ab13ao13IndexOK || !ab13ao13SelectorOK || !ab13ao13LocalOK || !ab13ao13CountOK || !ab13ao13ParsedOK {
		ab13ao13Output["state"] = "partial"
		return ab13ao13Finish()
	}
	if ab13ao13Registration[0]&1 == 0 || ab13ao13Index == 1 || ab13ao13Selector == ab13ao13LocalSelector || !ab13ao13Parsed || ab13ao13Count == 0 {
		return ab13ao13Finish()
	}
	ab13ao13Body := ab13ao13Bytes["local_body"]
	if len(ab13ao13Body) < 9 {
		ab13ao13Output["state"] = "partial"
		return ab13ao13Finish()
	}
	ab13ao13Fields["copied"] = map[string]any{"raw": ab13ao13Body[:9], "state": "Hypothesis"}
	ab13ao13FailureField["state"] = "Hypothesis"
	ab13ao13Fields["failure"] = ab13ao13FailureField
	ab13ao13FailureField["raw"] = ab13ao13Body[:9]
	ab13ao13FailureField["flags"] = map[string]bool{"failure": ab13ao13Body[8]&1 != 0}
	return ab13ao13Finish()
}
