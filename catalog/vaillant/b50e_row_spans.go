package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

func b50ERowSpanRead(b50ERowSpanDecoder *json.Decoder) (any, error) {
	b50ERowSpanToken, b50ERowSpanErr := b50ERowSpanDecoder.Token()
	if b50ERowSpanErr != nil {
		return nil, b50ERowSpanErr
	}
	switch b50ERowSpanToken {
	case json.Delim('{'):
		b50ERowSpanObject := map[string]any{}
		for b50ERowSpanDecoder.More() {
			b50ERowSpanKey, b50ERowSpanErr := b50ERowSpanDecoder.Token()
			if b50ERowSpanErr != nil {
				return nil, b50ERowSpanErr
			}
			b50ERowSpanName := b50ERowSpanKey.(string)
			if _, b50ERowSpanExists := b50ERowSpanObject[b50ERowSpanName]; b50ERowSpanExists {
				return nil, fmt.Errorf("duplicate key %q", b50ERowSpanName)
			}
			b50ERowSpanValue, b50ERowSpanErr := b50ERowSpanRead(b50ERowSpanDecoder)
			if b50ERowSpanErr != nil {
				return nil, b50ERowSpanErr
			}
			b50ERowSpanObject[b50ERowSpanName] = b50ERowSpanValue
		}
		_, b50ERowSpanErr := b50ERowSpanDecoder.Token()
		return b50ERowSpanObject, b50ERowSpanErr
	case json.Delim('['):
		b50ERowSpanArray := []any{}
		for b50ERowSpanDecoder.More() {
			b50ERowSpanValue, b50ERowSpanErr := b50ERowSpanRead(b50ERowSpanDecoder)
			if b50ERowSpanErr != nil {
				return nil, b50ERowSpanErr
			}
			b50ERowSpanArray = append(b50ERowSpanArray, b50ERowSpanValue)
		}
		_, b50ERowSpanErr := b50ERowSpanDecoder.Token()
		return b50ERowSpanArray, b50ERowSpanErr
	}
	return b50ERowSpanToken, nil
}

func b50ERowSpanInteger(b50ERowSpanValue any) (int64, bool) {
	b50ERowSpanNumber, b50ERowSpanOK := b50ERowSpanValue.(json.Number)
	if !b50ERowSpanOK {
		return 0, false
	}
	b50ERowSpanInt, b50ERowSpanErr := strconv.ParseInt(string(b50ERowSpanNumber), 10, 64)
	return b50ERowSpanInt, b50ERowSpanErr == nil
}

func B50ERowSpanInterpret(b50ERowSpanInput json.RawMessage) (json.RawMessage, error) {
	if !utf8.Valid(b50ERowSpanInput) {
		return nil, fmt.Errorf("invalid UTF-8 JSON")
	}
	b50ERowSpanDecoder := json.NewDecoder(bytes.NewReader(b50ERowSpanInput))
	b50ERowSpanDecoder.UseNumber()
	b50ERowSpanParsed, b50ERowSpanErr := b50ERowSpanRead(b50ERowSpanDecoder)
	if b50ERowSpanErr != nil {
		return nil, b50ERowSpanErr
	}
	if _, b50ERowSpanErr = b50ERowSpanDecoder.Token(); b50ERowSpanErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b50ERowSpanObject, b50ERowSpanOK := b50ERowSpanParsed.(map[string]any)
	if !b50ERowSpanOK {
		return nil, fmt.Errorf("input must be object")
	}
	for b50ERowSpanName := range b50ERowSpanObject {
		if b50ERowSpanName != "context" && b50ERowSpanName != "fields" && b50ERowSpanName != "arguments" {
			return nil, fmt.Errorf("unknown input member")
		}
	}
	b50ERowSpanContext, b50ERowSpanOK := b50ERowSpanObject["context"].(map[string]any)
	if !b50ERowSpanOK {
		return nil, fmt.Errorf("context must be object")
	}
	b50ERowSpanRaw, b50ERowSpanOK := b50ERowSpanObject["fields"].(map[string]any)
	if !b50ERowSpanOK {
		return nil, fmt.Errorf("fields must be object")
	}
	for _, b50ERowSpanValue := range b50ERowSpanRaw {
		b50ERowSpanArray, b50ERowSpanOK := b50ERowSpanValue.([]any)
		if !b50ERowSpanOK {
			return nil, fmt.Errorf("field must be byte array")
		}
		for _, b50ERowSpanByte := range b50ERowSpanArray {
			b50ERowSpanInt, b50ERowSpanOK := b50ERowSpanInteger(b50ERowSpanByte)
			if !b50ERowSpanOK || b50ERowSpanInt < 0 || b50ERowSpanInt > 255 {
				return nil, fmt.Errorf("invalid byte")
			}
		}
	}
	for _, b50ERowSpanName := range []string{"primary", "secondary", "operation"} {
		if b50ERowSpanValue, b50ERowSpanExists := b50ERowSpanContext[b50ERowSpanName]; b50ERowSpanExists {
			if _, b50ERowSpanOK := b50ERowSpanInteger(b50ERowSpanValue); !b50ERowSpanOK {
				return nil, fmt.Errorf("invalid context integer")
			}
		}
	}
	for _, b50ERowSpanName := range []string{"direction", "view"} {
		if b50ERowSpanValue, b50ERowSpanExists := b50ERowSpanContext[b50ERowSpanName]; b50ERowSpanExists {
			if _, b50ERowSpanOK := b50ERowSpanValue.(string); !b50ERowSpanOK {
				return nil, fmt.Errorf("invalid context token")
			}
		}
	}
	b50ERowSpanFields := map[string]any{}
	b50ERowSpanOutput := map[string]any{"context": b50ERowSpanContext, "raw": b50ERowSpanRaw, "state": "Unknown", "fields": b50ERowSpanFields, "physical": "Unknown", "qualified": false}
	if b50ERowSpanArguments, b50ERowSpanExists := b50ERowSpanObject["arguments"]; b50ERowSpanExists {
		b50ERowSpanMap, b50ERowSpanOK := b50ERowSpanArguments.(map[string]any)
		if !b50ERowSpanOK || len(b50ERowSpanMap) != 0 {
			return nil, fmt.Errorf("no arguments defined for row-spans")
		}
		b50ERowSpanOutput["arguments"] = b50ERowSpanMap
	}
	b50ERowSpanPrimary, _ := b50ERowSpanInteger(b50ERowSpanContext["primary"])
	b50ERowSpanSecondary, _ := b50ERowSpanInteger(b50ERowSpanContext["secondary"])
	b50ERowSpanOperation, _ := b50ERowSpanInteger(b50ERowSpanContext["operation"])
	if b50ERowSpanPrimary != 181 || b50ERowSpanSecondary != 14 || b50ERowSpanOperation != 2 || b50ERowSpanContext["direction"] != "positive" || b50ERowSpanContext["view"] != "row-spans" {
		return json.Marshal(b50ERowSpanOutput)
	}
	b50ERowSpanKeys, _ := b50ERowSpanRaw["key"].([]any)
	if len(b50ERowSpanKeys) != 1 {
		return json.Marshal(b50ERowSpanOutput)
	}
	b50ERowSpanKey, _ := b50ERowSpanInteger(b50ERowSpanKeys[0])
	if b50ERowSpanKey > 18 {
		return json.Marshal(b50ERowSpanOutput)
	}
	b50ERowSpanValue, b50ERowSpanPresent := b50ERowSpanRaw["value"].([]any)
	if b50ERowSpanKey == 0 || b50ERowSpanKey == 15 || b50ERowSpanKey == 16 {
		b50ERowSpanField := map[string]any{"state": "Unknown"}
		if b50ERowSpanPresent {
			b50ERowSpanField["raw"] = b50ERowSpanValue
		}
		b50ERowSpanFields["value"] = b50ERowSpanField
		return json.Marshal(b50ERowSpanOutput)
	}
	b50ERowSpanOutput["state"] = "Hypothesis"
	b50ERowSpanWidth := 1
	switch b50ERowSpanKey {
	case 4, 5, 11, 12, 17, 18:
		b50ERowSpanWidth = 2
	case 6, 7, 13, 14:
		b50ERowSpanWidth = 4
	}
	b50ERowSpanSpan := func(b50ERowSpanName string, b50ERowSpanOffset, b50ERowSpanSize int) map[string]any {
		b50ERowSpanEnd := b50ERowSpanOffset + b50ERowSpanSize
		if b50ERowSpanEnd > len(b50ERowSpanValue) {
			b50ERowSpanEnd = len(b50ERowSpanValue)
		}
		if b50ERowSpanOffset > b50ERowSpanEnd {
			b50ERowSpanOffset = b50ERowSpanEnd
		}
		b50ERowSpanField := map[string]any{"state": "Hypothesis"}
		if b50ERowSpanPresent {
			b50ERowSpanField["raw"] = b50ERowSpanValue[b50ERowSpanOffset:b50ERowSpanEnd]
		}
		if b50ERowSpanEnd-b50ERowSpanOffset < b50ERowSpanSize {
			b50ERowSpanField["state"] = "incomplete"
			b50ERowSpanOutput["state"] = "partial"
		}
		b50ERowSpanFields[b50ERowSpanName] = b50ERowSpanField
		return b50ERowSpanField
	}
	b50ERowSpanField := b50ERowSpanSpan("value", 0, b50ERowSpanWidth)
	if b50ERowSpanField["state"] == "Hypothesis" && b50ERowSpanWidth == 1 {
		b50ERowSpanNumber, _ := b50ERowSpanInteger(b50ERowSpanValue[0])
		switch b50ERowSpanKey {
		case 1:
			b50ERowSpanField["state"] = "Unknown"
		case 8:
			switch b50ERowSpanNumber {
			case 15:
				b50ERowSpanField["label"] = "true"
			case 240:
				b50ERowSpanField["label"] = "false"
			default:
				b50ERowSpanField["state"] = "Unknown"
			}
		default:
			if (b50ERowSpanKey == 3 || b50ERowSpanKey == 10) && b50ERowSpanNumber >= 128 {
				b50ERowSpanNumber -= 256
			}
			b50ERowSpanField["number"] = b50ERowSpanNumber
			b50ERowSpanField["calculation"] = "Proven"
		}
	}
	if b50ERowSpanKey >= 9 && b50ERowSpanKey <= 14 || b50ERowSpanKey == 18 {
		b50ERowSpanSpan("inverse", b50ERowSpanWidth, b50ERowSpanWidth)
	}
	if b50ERowSpanKey == 17 {
		b50ERowSpanStatus := b50ERowSpanSpan("status", 2, 1)
		if b50ERowSpanStatus["state"] == "Hypothesis" {
			b50ERowSpanNumber, _ := b50ERowSpanInteger(b50ERowSpanValue[2])
			switch b50ERowSpanNumber {
			case 0:
				b50ERowSpanStatus["label"] = "okay"
			case 85:
				b50ERowSpanStatus["label"] = "short-circuit"
			case 170:
				b50ERowSpanStatus["label"] = "interrupted"
			default:
				b50ERowSpanStatus["state"] = "Unknown"
			}
		}
	}
	return json.Marshal(b50ERowSpanOutput)
}
