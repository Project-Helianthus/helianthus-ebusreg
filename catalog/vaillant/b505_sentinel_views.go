package vaillant

import (
	b505SentinelViewBytes "bytes"
	b505SentinelViewJSON "encoding/json"
	b505SentinelViewFmt "fmt"
	b505SentinelViewIO "io"
	b505SentinelViewStrconv "strconv"
	b505SentinelViewStrings "strings"
)

func b505SentinelViewRead(b505SentinelViewDecoder *b505SentinelViewJSON.Decoder) (any, error) {
	b505SentinelViewToken, b505SentinelViewError := b505SentinelViewDecoder.Token()
	if b505SentinelViewError != nil {
		return nil, b505SentinelViewError
	}
	if b505SentinelViewDelimiter, b505SentinelViewOK := b505SentinelViewToken.(b505SentinelViewJSON.Delim); b505SentinelViewOK {
		switch b505SentinelViewDelimiter {
		case '{':
			b505SentinelViewObject := map[string]any{}
			for b505SentinelViewDecoder.More() {
				b505SentinelViewKeyToken, b505SentinelViewError := b505SentinelViewDecoder.Token()
				if b505SentinelViewError != nil {
					return nil, b505SentinelViewError
				}
				b505SentinelViewKey, b505SentinelViewOK := b505SentinelViewKeyToken.(string)
				if !b505SentinelViewOK {
					return nil, b505SentinelViewFmt.Errorf("ab01ao01: invalid object key")
				}
				if _, b505SentinelViewExists := b505SentinelViewObject[b505SentinelViewKey]; b505SentinelViewExists {
					return nil, b505SentinelViewFmt.Errorf("ab01ao01: duplicate key %q", b505SentinelViewKey)
				}
				b505SentinelViewValue, b505SentinelViewError := b505SentinelViewRead(b505SentinelViewDecoder)
				if b505SentinelViewError != nil {
					return nil, b505SentinelViewError
				}
				b505SentinelViewObject[b505SentinelViewKey] = b505SentinelViewValue
			}
			_, b505SentinelViewError := b505SentinelViewDecoder.Token()
			return b505SentinelViewObject, b505SentinelViewError
		case '[':
			b505SentinelViewArray := []any{}
			for b505SentinelViewDecoder.More() {
				b505SentinelViewValue, b505SentinelViewError := b505SentinelViewRead(b505SentinelViewDecoder)
				if b505SentinelViewError != nil {
					return nil, b505SentinelViewError
				}
				b505SentinelViewArray = append(b505SentinelViewArray, b505SentinelViewValue)
			}
			_, b505SentinelViewError := b505SentinelViewDecoder.Token()
			return b505SentinelViewArray, b505SentinelViewError
		default:
			return nil, b505SentinelViewFmt.Errorf("ab01ao01: unexpected delimiter")
		}
	}
	return b505SentinelViewToken, nil
}

func B505SentinelViewInterpret(b505SentinelViewInput b505SentinelViewJSON.RawMessage) (b505SentinelViewJSON.RawMessage, error) {
	b505SentinelViewDecoder := b505SentinelViewJSON.NewDecoder(b505SentinelViewBytes.NewReader(b505SentinelViewInput))
	b505SentinelViewDecoder.UseNumber()
	b505SentinelViewValue, b505SentinelViewError := b505SentinelViewRead(b505SentinelViewDecoder)
	if b505SentinelViewError != nil {
		return nil, b505SentinelViewError
	}
	if _, b505SentinelViewError := b505SentinelViewDecoder.Token(); b505SentinelViewError != b505SentinelViewIO.EOF {
		return nil, b505SentinelViewFmt.Errorf("ab01ao01: trailing JSON")
	}
	b505SentinelViewRoot, b505SentinelViewOK := b505SentinelViewValue.(map[string]any)
	if !b505SentinelViewOK {
		return nil, b505SentinelViewFmt.Errorf("ab01ao01: input must be an object")
	}
	for b505SentinelViewKey := range b505SentinelViewRoot {
		if b505SentinelViewKey != "context" && b505SentinelViewKey != "fields" && b505SentinelViewKey != "arguments" {
			return nil, b505SentinelViewFmt.Errorf("ab01ao01: unknown input key %q", b505SentinelViewKey)
		}
	}
	b505SentinelViewContext, b505SentinelViewOK := b505SentinelViewRoot["context"].(map[string]any)
	if !b505SentinelViewOK {
		return nil, b505SentinelViewFmt.Errorf("ab01ao01: context must be an object")
	}
	for _, b505SentinelViewKey := range []string{"primary", "secondary", "operation"} {
		if b505SentinelViewValue, b505SentinelViewPresent := b505SentinelViewContext[b505SentinelViewKey]; b505SentinelViewPresent {
			b505SentinelViewNumber, b505SentinelViewOK := b505SentinelViewValue.(b505SentinelViewJSON.Number)
			if !b505SentinelViewOK {
				return nil, b505SentinelViewFmt.Errorf("ab01ao01: %s must be an integer", b505SentinelViewKey)
			}
			if b505SentinelViewStrings.ContainsAny(string(b505SentinelViewNumber), ".eE") {
				return nil, b505SentinelViewFmt.Errorf("ab01ao01: %s must be an integer", b505SentinelViewKey)
			}
		}
	}
	for _, b505SentinelViewKey := range []string{"direction", "view"} {
		if b505SentinelViewValue, b505SentinelViewPresent := b505SentinelViewContext[b505SentinelViewKey]; b505SentinelViewPresent {
			if _, b505SentinelViewOK := b505SentinelViewValue.(string); !b505SentinelViewOK {
				return nil, b505SentinelViewFmt.Errorf("ab01ao01: %s must be a string", b505SentinelViewKey)
			}
		}
	}
	b505SentinelViewRaw, b505SentinelViewOK := b505SentinelViewRoot["fields"].(map[string]any)
	if !b505SentinelViewOK {
		return nil, b505SentinelViewFmt.Errorf("ab01ao01: fields must be an object")
	}
	for b505SentinelViewKey, b505SentinelViewValue := range b505SentinelViewRaw {
		b505SentinelViewArray, b505SentinelViewOK := b505SentinelViewValue.([]any)
		if !b505SentinelViewOK {
			return nil, b505SentinelViewFmt.Errorf("ab01ao01: %s must be a byte array", b505SentinelViewKey)
		}
		for _, b505SentinelViewElement := range b505SentinelViewArray {
			b505SentinelViewNumber, b505SentinelViewOK := b505SentinelViewElement.(b505SentinelViewJSON.Number)
			if !b505SentinelViewOK {
				return nil, b505SentinelViewFmt.Errorf("ab01ao01: invalid byte")
			}
			b505SentinelViewByte, b505SentinelViewError := b505SentinelViewStrconv.ParseInt(string(b505SentinelViewNumber), 10, 64)
			if b505SentinelViewError != nil || b505SentinelViewByte < 0 || b505SentinelViewByte > 255 {
				return nil, b505SentinelViewFmt.Errorf("ab01ao01: invalid byte")
			}
		}
	}
	b505SentinelViewFields := map[string]any{}
	b505SentinelViewOutput := map[string]any{"context": b505SentinelViewContext, "raw": b505SentinelViewRaw, "fields": b505SentinelViewFields, "state": "Unknown", "physical": "Unknown", "qualified": false}
	if b505SentinelViewArguments, b505SentinelViewPresent := b505SentinelViewRoot["arguments"]; b505SentinelViewPresent {
		b505SentinelViewObject, b505SentinelViewOK := b505SentinelViewArguments.(map[string]any)
		if !b505SentinelViewOK || len(b505SentinelViewObject) != 0 {
			return nil, b505SentinelViewFmt.Errorf("ab01ao01: no arguments are defined for this view")
		}
		b505SentinelViewOutput["arguments"] = b505SentinelViewObject
	}
	if b505SentinelViewEquals(b505SentinelViewContext["primary"], 181) && b505SentinelViewEquals(b505SentinelViewContext["secondary"], 5) && b505SentinelViewEquals(b505SentinelViewContext["operation"], 50) && b505SentinelViewContext["direction"] == "request" && b505SentinelViewContext["view"] == "replacement-pair" {
		b505SentinelViewKnown := 0
		b505SentinelViewIncomplete := false
		for _, b505SentinelViewKey := range []string{"temperature_one", "temperature_two"} {
			b505SentinelViewValue, b505SentinelViewPresent := b505SentinelViewRaw[b505SentinelViewKey]
			b505SentinelViewField := map[string]any{"raw": b505SentinelViewValue, "state": "missing"}
			b505SentinelViewFields[b505SentinelViewKey] = b505SentinelViewField
			if !b505SentinelViewPresent {
				b505SentinelViewIncomplete = true
				continue
			}
			b505SentinelViewArray := b505SentinelViewValue.([]any)
			if len(b505SentinelViewArray) != 2 {
				b505SentinelViewField["state"] = "incomplete"
				b505SentinelViewIncomplete = true
				continue
			}
			b505SentinelViewField["grade"] = "Hypothesis"
			b505SentinelViewField["state"] = "Unknown"
			b505SentinelViewField["label"] = "numeric-unqualified"
			if b505SentinelViewEquals(b505SentinelViewArray[0], 0) && (b505SentinelViewEquals(b505SentinelViewArray[1], 128) || b505SentinelViewEquals(b505SentinelViewArray[1], 0)) {
				b505SentinelViewField["state"] = "Hypothesis"
				b505SentinelViewField["label"] = "circuit-default"
				if b505SentinelViewEquals(b505SentinelViewArray[1], 128) {
					b505SentinelViewField["label"] = "preserve-existing"
				}
				b505SentinelViewKnown++
			}
		}
		if b505SentinelViewKnown == 2 {
			b505SentinelViewOutput["state"] = "Hypothesis"
		} else if b505SentinelViewKnown > 0 || b505SentinelViewIncomplete {
			b505SentinelViewOutput["state"] = "partial"
		}
	}
	return b505SentinelViewJSON.Marshal(b505SentinelViewOutput)
}

func b505SentinelViewEquals(b505SentinelViewValue any, b505SentinelViewExpected int64) bool {
	b505SentinelViewNumber, b505SentinelViewOK := b505SentinelViewValue.(b505SentinelViewJSON.Number)
	if !b505SentinelViewOK {
		return false
	}
	b505SentinelViewInteger, b505SentinelViewError := b505SentinelViewNumber.Int64()
	return b505SentinelViewError == nil && b505SentinelViewInteger == b505SentinelViewExpected
}
