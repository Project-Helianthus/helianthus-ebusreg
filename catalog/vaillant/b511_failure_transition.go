package vaillant

import (
	b511FailureTransitionBytesPackage "bytes"
	"encoding/json"
	b511FailureTransitionFmt "fmt"
	b511FailureTransitionIO "io"
	b511FailureTransitionStrconv "strconv"
)

func b511FailureTransitionRead(b511FailureTransitionDecoder *json.Decoder) (any, error) {
	b511FailureTransitionToken, b511FailureTransitionErr := b511FailureTransitionDecoder.Token()
	if b511FailureTransitionErr != nil {
		return nil, b511FailureTransitionErr
	}
	if b511FailureTransitionDelimiter, b511FailureTransitionOK := b511FailureTransitionToken.(json.Delim); b511FailureTransitionOK {
		switch b511FailureTransitionDelimiter {
		case '{':
			b511FailureTransitionObject := map[string]any{}
			for b511FailureTransitionDecoder.More() {
				b511FailureTransitionKeyToken, b511FailureTransitionErr := b511FailureTransitionDecoder.Token()
				if b511FailureTransitionErr != nil {
					return nil, b511FailureTransitionErr
				}
				b511FailureTransitionKey := b511FailureTransitionKeyToken.(string)
				if _, b511FailureTransitionExists := b511FailureTransitionObject[b511FailureTransitionKey]; b511FailureTransitionExists {
					return nil, b511FailureTransitionFmt.Errorf("duplicate key %q", b511FailureTransitionKey)
				}
				b511FailureTransitionValue, b511FailureTransitionErr := b511FailureTransitionRead(b511FailureTransitionDecoder)
				if b511FailureTransitionErr != nil {
					return nil, b511FailureTransitionErr
				}
				b511FailureTransitionObject[b511FailureTransitionKey] = b511FailureTransitionValue
			}
			_, b511FailureTransitionErr := b511FailureTransitionDecoder.Token()
			return b511FailureTransitionObject, b511FailureTransitionErr
		case '[':
			b511FailureTransitionArray := []any{}
			for b511FailureTransitionDecoder.More() {
				b511FailureTransitionValue, b511FailureTransitionErr := b511FailureTransitionRead(b511FailureTransitionDecoder)
				if b511FailureTransitionErr != nil {
					return nil, b511FailureTransitionErr
				}
				b511FailureTransitionArray = append(b511FailureTransitionArray, b511FailureTransitionValue)
			}
			_, b511FailureTransitionErr := b511FailureTransitionDecoder.Token()
			return b511FailureTransitionArray, b511FailureTransitionErr
		}
		return nil, b511FailureTransitionFmt.Errorf("unexpected delimiter")
	}
	return b511FailureTransitionToken, nil
}

func b511FailureTransitionInteger(b511FailureTransitionValue any) (int64, bool) {
	b511FailureTransitionNumber, b511FailureTransitionOK := b511FailureTransitionValue.(json.Number)
	if !b511FailureTransitionOK {
		return 0, false
	}
	b511FailureTransitionInt, b511FailureTransitionErr := b511FailureTransitionStrconv.ParseInt(string(b511FailureTransitionNumber), 10, 64)
	return b511FailureTransitionInt, b511FailureTransitionErr == nil
}

func B511FailureTransitionInterpret(b511FailureTransitionInput json.RawMessage) (json.RawMessage, error) {
	b511FailureTransitionDecoder := json.NewDecoder(b511FailureTransitionBytesPackage.NewReader(b511FailureTransitionInput))
	b511FailureTransitionDecoder.UseNumber()
	b511FailureTransitionValue, b511FailureTransitionErr := b511FailureTransitionRead(b511FailureTransitionDecoder)
	if b511FailureTransitionErr != nil {
		return nil, b511FailureTransitionErr
	}
	if _, b511FailureTransitionErr = b511FailureTransitionDecoder.Token(); b511FailureTransitionErr != b511FailureTransitionIO.EOF {
		return nil, b511FailureTransitionFmt.Errorf("trailing JSON")
	}
	b511FailureTransitionInputObject, b511FailureTransitionOK := b511FailureTransitionValue.(map[string]any)
	if !b511FailureTransitionOK {
		return nil, b511FailureTransitionFmt.Errorf("input must be object")
	}
	b511FailureTransitionContext, b511FailureTransitionOK := b511FailureTransitionInputObject["context"].(map[string]any)
	if !b511FailureTransitionOK {
		return nil, b511FailureTransitionFmt.Errorf("context must be object")
	}
	b511FailureTransitionRaw, b511FailureTransitionOK := b511FailureTransitionInputObject["fields"].(map[string]any)
	if !b511FailureTransitionOK {
		return nil, b511FailureTransitionFmt.Errorf("fields must be object")
	}
	b511FailureTransitionBytes := map[string][]int64{}
	for b511FailureTransitionKey, b511FailureTransitionValue := range b511FailureTransitionRaw {
		b511FailureTransitionArray, b511FailureTransitionOK := b511FailureTransitionValue.([]any)
		if !b511FailureTransitionOK {
			return nil, b511FailureTransitionFmt.Errorf("field %s must be byte array", b511FailureTransitionKey)
		}
		b511FailureTransitionBytes[b511FailureTransitionKey] = []int64{}
		for _, b511FailureTransitionValue := range b511FailureTransitionArray {
			b511FailureTransitionInt, b511FailureTransitionOK := b511FailureTransitionInteger(b511FailureTransitionValue)
			if !b511FailureTransitionOK || b511FailureTransitionInt < 0 || b511FailureTransitionInt > 255 {
				return nil, b511FailureTransitionFmt.Errorf("invalid byte")
			}
			b511FailureTransitionBytes[b511FailureTransitionKey] = append(b511FailureTransitionBytes[b511FailureTransitionKey], b511FailureTransitionInt)
		}
	}
	for _, b511FailureTransitionKey := range []string{"primary", "secondary", "operation"} {
		if b511FailureTransitionValue, b511FailureTransitionExists := b511FailureTransitionContext[b511FailureTransitionKey]; b511FailureTransitionExists {
			if _, b511FailureTransitionOK := b511FailureTransitionInteger(b511FailureTransitionValue); !b511FailureTransitionOK {
				return nil, b511FailureTransitionFmt.Errorf("invalid context integer")
			}
		}
	}
	for _, b511FailureTransitionKey := range []string{"direction", "view"} {
		if b511FailureTransitionValue, b511FailureTransitionExists := b511FailureTransitionContext[b511FailureTransitionKey]; b511FailureTransitionExists {
			if _, b511FailureTransitionOK := b511FailureTransitionValue.(string); !b511FailureTransitionOK {
				return nil, b511FailureTransitionFmt.Errorf("invalid context string")
			}
		}
	}
	b511FailureTransitionArguments := map[string]any{}
	b511FailureTransitionArgumentValue, b511FailureTransitionHasArguments := b511FailureTransitionInputObject["arguments"]
	if b511FailureTransitionHasArguments {
		b511FailureTransitionArguments, b511FailureTransitionOK = b511FailureTransitionArgumentValue.(map[string]any)
		if !b511FailureTransitionOK {
			return nil, b511FailureTransitionFmt.Errorf("arguments must be object")
		}
	}
	for b511FailureTransitionKey, b511FailureTransitionValue := range b511FailureTransitionArguments {
		switch b511FailureTransitionKey {
		case "parsed":
			if _, b511FailureTransitionOK := b511FailureTransitionValue.(bool); !b511FailureTransitionOK {
				return nil, b511FailureTransitionFmt.Errorf("parsed must be boolean")
			}
		case "index", "selector", "local_selector", "count":
			if b511FailureTransitionInt, b511FailureTransitionOK := b511FailureTransitionInteger(b511FailureTransitionValue); !b511FailureTransitionOK || b511FailureTransitionInt < 0 {
				return nil, b511FailureTransitionFmt.Errorf("invalid argument integer")
			}
		default:
			return nil, b511FailureTransitionFmt.Errorf("unknown argument")
		}
	}
	b511FailureTransitionFields := map[string]any{}
	b511FailureTransitionOutput := map[string]any{"context": b511FailureTransitionContext, "raw": b511FailureTransitionRaw, "state": "Unknown", "fields": b511FailureTransitionFields, "physical": "Unknown", "qualified": false}
	if b511FailureTransitionHasArguments {
		b511FailureTransitionOutput["arguments"] = b511FailureTransitionArguments
	}
	b511FailureTransitionFinish := func() (json.RawMessage, error) { return json.Marshal(b511FailureTransitionOutput) }
	b511FailureTransitionPrimary, _ := b511FailureTransitionInteger(b511FailureTransitionContext["primary"])
	b511FailureTransitionSecondary, _ := b511FailureTransitionInteger(b511FailureTransitionContext["secondary"])
	if b511FailureTransitionPrimary != 181 || b511FailureTransitionSecondary != 17 || b511FailureTransitionContext["direction"] != "local" {
		return b511FailureTransitionFinish()
	}
	b511FailureTransitionView, _ := b511FailureTransitionContext["view"].(string)
	b511FailureTransitionOperation, b511FailureTransitionHasOperation := b511FailureTransitionContext["operation"]
	if b511FailureTransitionView == "refresh" {
		b511FailureTransitionOp, _ := b511FailureTransitionInteger(b511FailureTransitionOperation)
		if !b511FailureTransitionHasOperation || b511FailureTransitionOp != 1 {
			return b511FailureTransitionFinish()
		}
	} else if b511FailureTransitionHasOperation || (b511FailureTransitionView != "initialize" && b511FailureTransitionView != "timer-consume") {
		return b511FailureTransitionFinish()
	}
	b511FailureTransitionOutput["state"] = "Hypothesis"
	if b511FailureTransitionView == "initialize" {
		b511FailureTransitionFields["failure"] = map[string]any{"raw": []any{}, "state": "Hypothesis", "flags": map[string]bool{"failure": false}}
		return b511FailureTransitionFinish()
	}
	b511FailureTransitionIndex, b511FailureTransitionIndexOK := b511FailureTransitionInteger(b511FailureTransitionArguments["index"])
	if b511FailureTransitionView == "timer-consume" && (!b511FailureTransitionIndexOK || b511FailureTransitionIndex != 0) {
		b511FailureTransitionOutput["state"] = "Unknown"
		return b511FailureTransitionFinish()
	}
	b511FailureTransitionPrior := b511FailureTransitionBytes["prior_failure"]
	b511FailureTransitionFailure := false
	b511FailureTransitionFailureField := map[string]any{}
	if len(b511FailureTransitionPrior) == 1 {
		b511FailureTransitionFailure = b511FailureTransitionPrior[0] != 0
		b511FailureTransitionFailureField = map[string]any{"raw": b511FailureTransitionRaw["prior_failure"], "state": "Hypothesis", "flags": map[string]bool{"failure": b511FailureTransitionFailure}}
		b511FailureTransitionFields["failure"] = b511FailureTransitionFailureField
	} else {
		b511FailureTransitionOutput["state"] = "partial"
		if b511FailureTransitionView == "timer-consume" {
			return b511FailureTransitionFinish()
		}
	}
	if b511FailureTransitionView == "timer-consume" {
		b511FailureTransitionPath := map[string]any{"raw": b511FailureTransitionRaw["prior_failure"], "state": "Hypothesis", "label": "capability-path"}
		if b511FailureTransitionFailure {
			b511FailureTransitionPath["label"] = "error"
			b511FailureTransitionPath["number"] = 120
		}
		b511FailureTransitionFields["path"] = b511FailureTransitionPath
		return b511FailureTransitionFinish()
	}
	b511FailureTransitionRegistration := b511FailureTransitionBytes["registration"]
	b511FailureTransitionSelector, b511FailureTransitionSelectorOK := b511FailureTransitionInteger(b511FailureTransitionArguments["selector"])
	b511FailureTransitionLocalSelector, b511FailureTransitionLocalOK := b511FailureTransitionInteger(b511FailureTransitionArguments["local_selector"])
	b511FailureTransitionCount, b511FailureTransitionCountOK := b511FailureTransitionInteger(b511FailureTransitionArguments["count"])
	b511FailureTransitionParsed, b511FailureTransitionParsedOK := b511FailureTransitionArguments["parsed"].(bool)
	if len(b511FailureTransitionRegistration) != 1 || !b511FailureTransitionIndexOK || !b511FailureTransitionSelectorOK || !b511FailureTransitionLocalOK || !b511FailureTransitionCountOK || !b511FailureTransitionParsedOK {
		b511FailureTransitionOutput["state"] = "partial"
		return b511FailureTransitionFinish()
	}
	if b511FailureTransitionRegistration[0]&1 == 0 || b511FailureTransitionIndex == 1 || b511FailureTransitionSelector == b511FailureTransitionLocalSelector || !b511FailureTransitionParsed || b511FailureTransitionCount == 0 {
		return b511FailureTransitionFinish()
	}
	b511FailureTransitionBody := b511FailureTransitionBytes["local_body"]
	if len(b511FailureTransitionBody) < 9 {
		b511FailureTransitionOutput["state"] = "partial"
		return b511FailureTransitionFinish()
	}
	b511FailureTransitionFields["copied"] = map[string]any{"raw": b511FailureTransitionBody[:9], "state": "Hypothesis"}
	b511FailureTransitionFailureField["state"] = "Hypothesis"
	b511FailureTransitionFields["failure"] = b511FailureTransitionFailureField
	b511FailureTransitionFailureField["raw"] = b511FailureTransitionBody[:9]
	b511FailureTransitionFailureField["flags"] = map[string]bool{"failure": b511FailureTransitionBody[8]&1 != 0}
	return b511FailureTransitionFinish()
}
