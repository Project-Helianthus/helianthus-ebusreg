package vaillant

import (
	b510CacheTransitionBytesPackage "bytes"
	b510CacheTransitionJSON "encoding/json"
	b510CacheTransitionFmt "fmt"
	b510CacheTransitionIO "io"
	b510CacheTransitionBig "math/big"
)

func b510CacheTransitionRead(b510CacheTransitionDecoder *b510CacheTransitionJSON.Decoder) (any, error) {
	b510CacheTransitionToken, b510CacheTransitionErr := b510CacheTransitionDecoder.Token()
	if b510CacheTransitionErr != nil {
		return nil, b510CacheTransitionErr
	}
	if b510CacheTransitionDelimiter, b510CacheTransitionOK := b510CacheTransitionToken.(b510CacheTransitionJSON.Delim); b510CacheTransitionOK {
		switch b510CacheTransitionDelimiter {
		case '{':
			b510CacheTransitionObject := map[string]any{}
			for b510CacheTransitionDecoder.More() {
				b510CacheTransitionKey, b510CacheTransitionErr := b510CacheTransitionDecoder.Token()
				if b510CacheTransitionErr != nil {
					return nil, b510CacheTransitionErr
				}
				b510CacheTransitionName, b510CacheTransitionOK := b510CacheTransitionKey.(string)
				if !b510CacheTransitionOK {
					return nil, b510CacheTransitionFmt.Errorf("expected object key")
				}
				if _, b510CacheTransitionExists := b510CacheTransitionObject[b510CacheTransitionName]; b510CacheTransitionExists {
					return nil, b510CacheTransitionFmt.Errorf("duplicate object key")
				}
				b510CacheTransitionValue, b510CacheTransitionErr := b510CacheTransitionRead(b510CacheTransitionDecoder)
				if b510CacheTransitionErr != nil {
					return nil, b510CacheTransitionErr
				}
				b510CacheTransitionObject[b510CacheTransitionName] = b510CacheTransitionValue
			}
			_, b510CacheTransitionErr = b510CacheTransitionDecoder.Token()
			return b510CacheTransitionObject, b510CacheTransitionErr
		case '[':
			b510CacheTransitionArray := []any{}
			for b510CacheTransitionDecoder.More() {
				b510CacheTransitionValue, b510CacheTransitionErr := b510CacheTransitionRead(b510CacheTransitionDecoder)
				if b510CacheTransitionErr != nil {
					return nil, b510CacheTransitionErr
				}
				b510CacheTransitionArray = append(b510CacheTransitionArray, b510CacheTransitionValue)
			}
			_, b510CacheTransitionErr = b510CacheTransitionDecoder.Token()
			return b510CacheTransitionArray, b510CacheTransitionErr
		default:
			return nil, b510CacheTransitionFmt.Errorf("unexpected delimiter")
		}
	}
	return b510CacheTransitionToken, nil
}

func b510CacheTransitionInteger(b510CacheTransitionValue any) (*b510CacheTransitionBig.Int, bool) {
	b510CacheTransitionNumber, b510CacheTransitionOK := b510CacheTransitionValue.(b510CacheTransitionJSON.Number)
	if !b510CacheTransitionOK {
		return nil, false
	}
	b510CacheTransitionResult, b510CacheTransitionOK := new(b510CacheTransitionBig.Int).SetString(string(b510CacheTransitionNumber), 10)
	return b510CacheTransitionResult, b510CacheTransitionOK
}

func B510CacheTransitionInterpret(b510CacheTransitionInput b510CacheTransitionJSON.RawMessage) (b510CacheTransitionJSON.RawMessage, error) {
	b510CacheTransitionDecoder := b510CacheTransitionJSON.NewDecoder(b510CacheTransitionBytesPackage.NewReader(b510CacheTransitionInput))
	b510CacheTransitionDecoder.UseNumber()
	b510CacheTransitionValue, b510CacheTransitionErr := b510CacheTransitionRead(b510CacheTransitionDecoder)
	if b510CacheTransitionErr != nil {
		return nil, b510CacheTransitionErr
	}
	if _, b510CacheTransitionErr = b510CacheTransitionDecoder.Token(); b510CacheTransitionErr != b510CacheTransitionIO.EOF {
		return nil, b510CacheTransitionFmt.Errorf("trailing JSON")
	}
	b510CacheTransitionObject, b510CacheTransitionOK := b510CacheTransitionValue.(map[string]any)
	if !b510CacheTransitionOK {
		return nil, b510CacheTransitionFmt.Errorf("expected input object")
	}
	b510CacheTransitionContext := map[string]any{}
	if b510CacheTransitionValue, b510CacheTransitionPresent := b510CacheTransitionObject["context"]; b510CacheTransitionPresent {
		b510CacheTransitionContext, b510CacheTransitionOK = b510CacheTransitionValue.(map[string]any)
		if !b510CacheTransitionOK {
			return nil, b510CacheTransitionFmt.Errorf("expected context object")
		}
	}
	b510CacheTransitionRaw := map[string]any{}
	if b510CacheTransitionValue, b510CacheTransitionPresent := b510CacheTransitionObject["fields"]; b510CacheTransitionPresent {
		b510CacheTransitionRaw, b510CacheTransitionOK = b510CacheTransitionValue.(map[string]any)
		if !b510CacheTransitionOK {
			return nil, b510CacheTransitionFmt.Errorf("expected fields object")
		}
	}
	for _, b510CacheTransitionValue := range b510CacheTransitionRaw {
		b510CacheTransitionArray, b510CacheTransitionOK := b510CacheTransitionValue.([]any)
		if !b510CacheTransitionOK {
			return nil, b510CacheTransitionFmt.Errorf("expected byte array")
		}
		for _, b510CacheTransitionElement := range b510CacheTransitionArray {
			b510CacheTransitionByte, b510CacheTransitionOK := b510CacheTransitionInteger(b510CacheTransitionElement)
			if !b510CacheTransitionOK || b510CacheTransitionByte.Sign() < 0 || b510CacheTransitionByte.Cmp(b510CacheTransitionBig.NewInt(255)) > 0 {
				return nil, b510CacheTransitionFmt.Errorf("invalid byte element")
			}
		}
	}
	b510CacheTransitionSelected := true
	for _, b510CacheTransitionAxis := range []struct {
		b510CacheTransitionName     string
		b510CacheTransitionExpected int64
	}{{"primary", 181}, {"secondary", 16}} {
		b510CacheTransitionValue, b510CacheTransitionPresent := b510CacheTransitionContext[b510CacheTransitionAxis.b510CacheTransitionName]
		if !b510CacheTransitionPresent {
			b510CacheTransitionSelected = false
			continue
		}
		b510CacheTransitionNumber, b510CacheTransitionOK := b510CacheTransitionInteger(b510CacheTransitionValue)
		if !b510CacheTransitionOK {
			return nil, b510CacheTransitionFmt.Errorf("expected context integer")
		}
		if b510CacheTransitionNumber.Cmp(b510CacheTransitionBig.NewInt(b510CacheTransitionAxis.b510CacheTransitionExpected)) != 0 {
			b510CacheTransitionSelected = false
		}
	}
	if b510CacheTransitionValue, b510CacheTransitionPresent := b510CacheTransitionContext["operation"]; b510CacheTransitionPresent {
		if _, b510CacheTransitionOK := b510CacheTransitionInteger(b510CacheTransitionValue); !b510CacheTransitionOK {
			return nil, b510CacheTransitionFmt.Errorf("expected operation integer")
		}
		b510CacheTransitionSelected = false
	}
	for _, b510CacheTransitionAxis := range []struct{ b510CacheTransitionName, b510CacheTransitionExpected string }{{"direction", "local"}, {"view", "capability-query"}} {
		b510CacheTransitionValue, b510CacheTransitionPresent := b510CacheTransitionContext[b510CacheTransitionAxis.b510CacheTransitionName]
		if !b510CacheTransitionPresent {
			b510CacheTransitionSelected = false
			continue
		}
		b510CacheTransitionText, b510CacheTransitionOK := b510CacheTransitionValue.(string)
		if !b510CacheTransitionOK {
			return nil, b510CacheTransitionFmt.Errorf("expected context string")
		}
		if b510CacheTransitionText != b510CacheTransitionAxis.b510CacheTransitionExpected {
			b510CacheTransitionSelected = false
		}
	}
	b510CacheTransitionArguments := map[string]any{}
	b510CacheTransitionArgumentValue, b510CacheTransitionArgumentsPresent := b510CacheTransitionObject["arguments"]
	if b510CacheTransitionArgumentsPresent {
		b510CacheTransitionArguments, b510CacheTransitionOK = b510CacheTransitionArgumentValue.(map[string]any)
		if !b510CacheTransitionOK {
			return nil, b510CacheTransitionFmt.Errorf("expected arguments object")
		}
	}
	for b510CacheTransitionName, b510CacheTransitionValue := range b510CacheTransitionArguments {
		switch b510CacheTransitionName {
		case "parsed":
			if _, b510CacheTransitionOK := b510CacheTransitionValue.(bool); !b510CacheTransitionOK {
				return nil, b510CacheTransitionFmt.Errorf("expected parsed boolean")
			}
		case "count":
			if _, b510CacheTransitionOK := b510CacheTransitionInteger(b510CacheTransitionValue); !b510CacheTransitionOK {
				return nil, b510CacheTransitionFmt.Errorf("expected count integer")
			}
		default:
			return nil, b510CacheTransitionFmt.Errorf("unknown argument")
		}
	}
	b510CacheTransitionFields := map[string]any{}
	b510CacheTransitionOutput := map[string]any{"context": b510CacheTransitionContext, "raw": b510CacheTransitionRaw, "state": "Unknown", "fields": b510CacheTransitionFields, "physical": "Unknown", "qualified": false}
	if b510CacheTransitionArgumentsPresent {
		b510CacheTransitionOutput["arguments"] = b510CacheTransitionArguments
	}
	b510CacheTransitionBytes := map[string]int64{}
	for _, b510CacheTransitionName := range []string{"prior", "feedback", "bit_six"} {
		b510CacheTransitionArray, b510CacheTransitionOK := b510CacheTransitionRaw[b510CacheTransitionName].([]any)
		if !b510CacheTransitionOK || len(b510CacheTransitionArray) != 1 {
			b510CacheTransitionSelected = false
			continue
		}
		b510CacheTransitionNumber, _ := b510CacheTransitionInteger(b510CacheTransitionArray[0])
		b510CacheTransitionBytes[b510CacheTransitionName] = b510CacheTransitionNumber.Int64()
	}
	if b510CacheTransitionBytes["bit_six"] != 0 && b510CacheTransitionBytes["bit_six"] != 64 {
		b510CacheTransitionSelected = false
	}
	b510CacheTransitionParsed, b510CacheTransitionHasParsed := b510CacheTransitionArguments["parsed"].(bool)
	b510CacheTransitionCount, b510CacheTransitionHasCount := b510CacheTransitionInteger(b510CacheTransitionArguments["count"])
	if !b510CacheTransitionHasParsed || !b510CacheTransitionHasCount || b510CacheTransitionCount.Sign() < 0 {
		b510CacheTransitionSelected = false
	}
	if !b510CacheTransitionSelected {
		return b510CacheTransitionJSON.Marshal(b510CacheTransitionOutput)
	}
	b510CacheTransitionCell := b510CacheTransitionBytes["prior"]
	b510CacheTransitionFreshness, b510CacheTransitionRule, b510CacheTransitionRuleState := "Unknown", "retain", "Hypothesis"
	if b510CacheTransitionParsed && b510CacheTransitionBytes["feedback"] == 1 {
		if b510CacheTransitionCount.Sign() == 0 {
			b510CacheTransitionRule, b510CacheTransitionRuleState = "Unknown", "Unknown"
		} else {
			b510CacheTransitionFreshness = "fresh"
			if b510CacheTransitionCount.Cmp(b510CacheTransitionBig.NewInt(2)) <= 0 {
				b510CacheTransitionCell, b510CacheTransitionRule = 0, "replace-zero"
			} else {
				b510CacheTransitionCell, b510CacheTransitionRule = b510CacheTransitionBytes["bit_six"], "replace-bit-six"
			}
		}
	}
	b510CacheTransitionFields["cache"] = map[string]any{"raw": b510CacheTransitionRaw["prior"], "state": "Hypothesis", "number": b510CacheTransitionCell}
	b510CacheTransitionFields["freshness"] = map[string]any{"raw": b510CacheTransitionRaw["feedback"], "state": "Hypothesis", "label": b510CacheTransitionFreshness}
	b510CacheTransitionFields["rule"] = map[string]any{"raw": b510CacheTransitionRaw["bit_six"], "state": b510CacheTransitionRuleState, "label": b510CacheTransitionRule}
	if b510CacheTransitionRuleState != "Unknown" {
		b510CacheTransitionOutput["state"] = "Hypothesis"
	}
	return b510CacheTransitionJSON.Marshal(b510CacheTransitionOutput)
}
