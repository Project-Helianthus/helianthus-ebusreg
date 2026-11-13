package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func b506ParameterViewRead(b506ParameterViewDecoder *json.Decoder) (any, error) {
	b506ParameterViewToken, b506ParameterViewErr := b506ParameterViewDecoder.Token()
	if b506ParameterViewErr != nil {
		return nil, b506ParameterViewErr
	}
	if b506ParameterViewDelimiter, b506ParameterViewOK := b506ParameterViewToken.(json.Delim); b506ParameterViewOK {
		switch b506ParameterViewDelimiter {
		case '{':
			b506ParameterViewObject := map[string]any{}
			for b506ParameterViewDecoder.More() {
				b506ParameterViewKeyToken, b506ParameterViewErr := b506ParameterViewDecoder.Token()
				if b506ParameterViewErr != nil {
					return nil, b506ParameterViewErr
				}
				b506ParameterViewKey, b506ParameterViewOK := b506ParameterViewKeyToken.(string)
				if !b506ParameterViewOK {
					return nil, fmt.Errorf("invalid object key")
				}
				if _, b506ParameterViewExists := b506ParameterViewObject[b506ParameterViewKey]; b506ParameterViewExists {
					return nil, fmt.Errorf("duplicate key %q", b506ParameterViewKey)
				}
				b506ParameterViewValue, b506ParameterViewErr := b506ParameterViewRead(b506ParameterViewDecoder)
				if b506ParameterViewErr != nil {
					return nil, b506ParameterViewErr
				}
				b506ParameterViewObject[b506ParameterViewKey] = b506ParameterViewValue
			}
			_, b506ParameterViewErr := b506ParameterViewDecoder.Token()
			return b506ParameterViewObject, b506ParameterViewErr
		case '[':
			b506ParameterViewArray := []any{}
			for b506ParameterViewDecoder.More() {
				b506ParameterViewValue, b506ParameterViewErr := b506ParameterViewRead(b506ParameterViewDecoder)
				if b506ParameterViewErr != nil {
					return nil, b506ParameterViewErr
				}
				b506ParameterViewArray = append(b506ParameterViewArray, b506ParameterViewValue)
			}
			_, b506ParameterViewErr := b506ParameterViewDecoder.Token()
			return b506ParameterViewArray, b506ParameterViewErr
		default:
			return nil, fmt.Errorf("unexpected delimiter")
		}
	}
	return b506ParameterViewToken, nil
}

func b506ParameterViewInteger(b506ParameterViewValue any) (int64, bool) {
	b506ParameterViewNumber, b506ParameterViewOK := b506ParameterViewValue.(json.Number)
	if !b506ParameterViewOK {
		return 0, false
	}
	b506ParameterViewInt, b506ParameterViewErr := strconv.ParseInt(string(b506ParameterViewNumber), 10, 64)
	return b506ParameterViewInt, b506ParameterViewErr == nil
}

func B506ParameterViewInterpret(b506ParameterViewInput json.RawMessage) (json.RawMessage, error) {
	b506ParameterViewDecoder := json.NewDecoder(bytes.NewReader(b506ParameterViewInput))
	b506ParameterViewDecoder.UseNumber()
	b506ParameterViewValue, b506ParameterViewErr := b506ParameterViewRead(b506ParameterViewDecoder)
	if b506ParameterViewErr != nil {
		return nil, b506ParameterViewErr
	}
	if _, b506ParameterViewErr = b506ParameterViewDecoder.Token(); b506ParameterViewErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b506ParameterViewObject, b506ParameterViewOK := b506ParameterViewValue.(map[string]any)
	if !b506ParameterViewOK {
		return nil, fmt.Errorf("input must be an object")
	}
	for b506ParameterViewKey := range b506ParameterViewObject {
		if b506ParameterViewKey != "context" && b506ParameterViewKey != "fields" && b506ParameterViewKey != "arguments" {
			return nil, fmt.Errorf("unknown input member")
		}
	}
	b506ParameterViewContext, b506ParameterViewOK := b506ParameterViewObject["context"].(map[string]any)
	if !b506ParameterViewOK {
		return nil, fmt.Errorf("context must be an object")
	}
	for _, b506ParameterViewKey := range []string{"primary", "secondary", "operation"} {
		if b506ParameterViewValue, b506ParameterViewPresent := b506ParameterViewContext[b506ParameterViewKey]; b506ParameterViewPresent {
			b506ParameterViewNumber, b506ParameterViewOK := b506ParameterViewValue.(json.Number)
			if !b506ParameterViewOK || strings.ContainsAny(string(b506ParameterViewNumber), ".eE") {
				return nil, fmt.Errorf("%s must be an integer", b506ParameterViewKey)
			}
		}
	}
	for _, b506ParameterViewKey := range []string{"direction", "view"} {
		if b506ParameterViewValue, b506ParameterViewPresent := b506ParameterViewContext[b506ParameterViewKey]; b506ParameterViewPresent {
			if _, b506ParameterViewOK := b506ParameterViewValue.(string); !b506ParameterViewOK {
				return nil, fmt.Errorf("%s must be a string", b506ParameterViewKey)
			}
		}
	}
	b506ParameterViewRaw, b506ParameterViewOK := b506ParameterViewObject["fields"].(map[string]any)
	if !b506ParameterViewOK {
		return nil, fmt.Errorf("fields must be an object")
	}
	for _, b506ParameterViewValue := range b506ParameterViewRaw {
		b506ParameterViewArray, b506ParameterViewOK := b506ParameterViewValue.([]any)
		if !b506ParameterViewOK {
			return nil, fmt.Errorf("field must be a byte array")
		}
		for _, b506ParameterViewByte := range b506ParameterViewArray {
			b506ParameterViewInt, b506ParameterViewOK := b506ParameterViewInteger(b506ParameterViewByte)
			if !b506ParameterViewOK || b506ParameterViewInt < 0 || b506ParameterViewInt > 255 {
				return nil, fmt.Errorf("invalid byte")
			}
		}
	}
	b506ParameterViewFields := map[string]any{}
	b506ParameterViewOutput := map[string]any{"context": b506ParameterViewContext, "raw": b506ParameterViewRaw, "state": "Unknown", "fields": b506ParameterViewFields, "physical": "Unknown", "qualified": false}
	if b506ParameterViewArguments, b506ParameterViewPresent := b506ParameterViewObject["arguments"]; b506ParameterViewPresent {
		b506ParameterViewMap, b506ParameterViewOK := b506ParameterViewArguments.(map[string]any)
		if !b506ParameterViewOK || len(b506ParameterViewMap) != 0 {
			return nil, fmt.Errorf("no arguments defined for this view")
		}
		b506ParameterViewOutput["arguments"] = b506ParameterViewMap
	}
	b506ParameterViewPrimary, _ := b506ParameterViewInteger(b506ParameterViewContext["primary"])
	b506ParameterViewSecondary, _ := b506ParameterViewInteger(b506ParameterViewContext["secondary"])
	b506ParameterViewOperation, _ := b506ParameterViewInteger(b506ParameterViewContext["operation"])
	if b506ParameterViewPrimary != 181 || b506ParameterViewSecondary != 6 || b506ParameterViewOperation != 17 || b506ParameterViewContext["direction"] != "request" || b506ParameterViewContext["view"] != "parameter-union" {
		return json.Marshal(b506ParameterViewOutput)
	}
	b506ParameterViewKey, b506ParameterViewPresent := b506ParameterViewRaw["key"].([]any)
	if !b506ParameterViewPresent || len(b506ParameterViewKey) != 1 {
		b506ParameterViewOutput["state"] = "partial"
		b506ParameterViewState := "incomplete"
		if len(b506ParameterViewKey) > 1 {
			b506ParameterViewState = "width-conflict"
		}
		b506ParameterViewField := map[string]any{"state": b506ParameterViewState}
		if b506ParameterViewPresent {
			b506ParameterViewField["raw"] = b506ParameterViewKey
		}
		b506ParameterViewFields["key"] = b506ParameterViewField
		return json.Marshal(b506ParameterViewOutput)
	}
	b506ParameterViewKeyNumber, _ := b506ParameterViewInteger(b506ParameterViewKey[0])
	b506ParameterViewOpaque := false
	switch b506ParameterViewKeyNumber {
	case 0, 1, 2, 4, 5:
		b506ParameterViewOpaque = true
	case 3, 6, 7, 8, 9, 10, 32, 33, 34, 35, 36, 37:
	default:
		return json.Marshal(b506ParameterViewOutput)
	}
	b506ParameterViewFields["key"] = map[string]any{"raw": b506ParameterViewKey, "state": "Hypothesis", "number": b506ParameterViewKeyNumber}
	b506ParameterViewField := map[string]any{"state": "incomplete"}
	b506ParameterViewFields["value"] = b506ParameterViewField
	b506ParameterViewOutput["state"] = "partial"
	b506ParameterViewBytes, b506ParameterViewPresent := b506ParameterViewRaw["value"].([]any)
	if !b506ParameterViewPresent {
		return json.Marshal(b506ParameterViewOutput)
	}
	b506ParameterViewField["raw"] = b506ParameterViewBytes
	if b506ParameterViewOpaque {
		b506ParameterViewField["state"] = "Hypothesis"
		b506ParameterViewOutput["state"] = "Hypothesis"
		return json.Marshal(b506ParameterViewOutput)
	}
	if len(b506ParameterViewBytes) == 0 {
		return json.Marshal(b506ParameterViewOutput)
	}
	if len(b506ParameterViewBytes) != 1 {
		b506ParameterViewField["state"] = "width-conflict"
		return json.Marshal(b506ParameterViewOutput)
	}
	b506ParameterViewNumber, _ := b506ParameterViewInteger(b506ParameterViewBytes[0])
	if b506ParameterViewKeyNumber == 10 {
		if b506ParameterViewNumber > 1 {
			b506ParameterViewField["state"] = "unmapped"
			return json.Marshal(b506ParameterViewOutput)
		}
		if b506ParameterViewNumber == 0 {
			b506ParameterViewField["label"] = "no"
		} else {
			b506ParameterViewField["label"] = "yes"
		}
	}
	b506ParameterViewField["state"] = "Hypothesis"
	b506ParameterViewField["number"] = b506ParameterViewNumber
	b506ParameterViewOutput["state"] = "Hypothesis"
	return json.Marshal(b506ParameterViewOutput)
}
