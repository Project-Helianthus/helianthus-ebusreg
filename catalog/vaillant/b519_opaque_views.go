package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func b519OpaqueViewValue(b519OpaqueViewDecoder *json.Decoder) (any, error) {
	b519OpaqueViewToken, b519OpaqueViewErr := b519OpaqueViewDecoder.Token()
	if b519OpaqueViewErr != nil {
		return nil, b519OpaqueViewErr
	}
	switch b519OpaqueViewToken {
	case json.Delim('{'):
		b519OpaqueViewObject := map[string]any{}
		for b519OpaqueViewDecoder.More() {
			b519OpaqueViewKeyToken, b519OpaqueViewErr := b519OpaqueViewDecoder.Token()
			if b519OpaqueViewErr != nil {
				return nil, b519OpaqueViewErr
			}
			b519OpaqueViewKey := b519OpaqueViewKeyToken.(string)
			if _, b519OpaqueViewFound := b519OpaqueViewObject[b519OpaqueViewKey]; b519OpaqueViewFound {
				return nil, fmt.Errorf("duplicate key %q", b519OpaqueViewKey)
			}
			b519OpaqueViewObject[b519OpaqueViewKey], b519OpaqueViewErr = b519OpaqueViewValue(b519OpaqueViewDecoder)
			if b519OpaqueViewErr != nil {
				return nil, b519OpaqueViewErr
			}
		}
		_, b519OpaqueViewErr := b519OpaqueViewDecoder.Token()
		return b519OpaqueViewObject, b519OpaqueViewErr
	case json.Delim('['):
		b519OpaqueViewArray := []any{}
		for b519OpaqueViewDecoder.More() {
			b519OpaqueViewItem, b519OpaqueViewErr := b519OpaqueViewValue(b519OpaqueViewDecoder)
			if b519OpaqueViewErr != nil {
				return nil, b519OpaqueViewErr
			}
			b519OpaqueViewArray = append(b519OpaqueViewArray, b519OpaqueViewItem)
		}
		_, b519OpaqueViewErr := b519OpaqueViewDecoder.Token()
		return b519OpaqueViewArray, b519OpaqueViewErr
	}
	return b519OpaqueViewToken, nil
}

func b519OpaqueViewInteger(b519OpaqueViewValue any) (int, bool) {
	b519OpaqueViewNumber, b519OpaqueViewOK := b519OpaqueViewValue.(json.Number)
	if !b519OpaqueViewOK {
		return 0, false
	}
	for b519OpaqueViewIndex, b519OpaqueViewDigit := range string(b519OpaqueViewNumber) {
		if b519OpaqueViewIndex == 0 && b519OpaqueViewDigit == '-' {
			continue
		}
		if b519OpaqueViewDigit < '0' || b519OpaqueViewDigit > '9' {
			return 0, false
		}
	}
	b519OpaqueViewInt, b519OpaqueViewErr := strconv.Atoi(string(b519OpaqueViewNumber))
	if b519OpaqueViewErr != nil {
		return -1, true
	}
	return b519OpaqueViewInt, true
}

func B519OpaqueViewInterpret(b519OpaqueViewInput json.RawMessage) (json.RawMessage, error) {
	b519OpaqueViewDecoder := json.NewDecoder(bytes.NewReader(b519OpaqueViewInput))
	b519OpaqueViewDecoder.UseNumber()
	b519OpaqueViewParsed, b519OpaqueViewErr := b519OpaqueViewValue(b519OpaqueViewDecoder)
	if b519OpaqueViewErr != nil {
		return nil, b519OpaqueViewErr
	}
	if _, b519OpaqueViewErr = b519OpaqueViewDecoder.Token(); b519OpaqueViewErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b519OpaqueViewRoot, b519OpaqueViewOK := b519OpaqueViewParsed.(map[string]any)
	if !b519OpaqueViewOK {
		return nil, fmt.Errorf("input must be object")
	}
	for b519OpaqueViewKey := range b519OpaqueViewRoot {
		if b519OpaqueViewKey != "context" && b519OpaqueViewKey != "fields" && b519OpaqueViewKey != "arguments" {
			return nil, fmt.Errorf("unknown input key %q", b519OpaqueViewKey)
		}
	}
	b519OpaqueViewContext, b519OpaqueViewOK := b519OpaqueViewRoot["context"].(map[string]any)
	if !b519OpaqueViewOK {
		return nil, fmt.Errorf("context must be object")
	}
	b519OpaqueViewRaw, b519OpaqueViewOK := b519OpaqueViewRoot["fields"].(map[string]any)
	if !b519OpaqueViewOK {
		return nil, fmt.Errorf("fields must be object")
	}
	for b519OpaqueViewKey, b519OpaqueViewValue := range b519OpaqueViewRaw {
		b519OpaqueViewArray, b519OpaqueViewOK := b519OpaqueViewValue.([]any)
		if !b519OpaqueViewOK {
			return nil, fmt.Errorf("field %q must be byte array", b519OpaqueViewKey)
		}
		for _, b519OpaqueViewByte := range b519OpaqueViewArray {
			b519OpaqueViewInt, b519OpaqueViewOK := b519OpaqueViewInteger(b519OpaqueViewByte)
			if !b519OpaqueViewOK || b519OpaqueViewInt < 0 || b519OpaqueViewInt > 255 {
				return nil, fmt.Errorf("invalid byte in %q", b519OpaqueViewKey)
			}
		}
	}
	for _, b519OpaqueViewKey := range []string{"primary", "secondary", "operation"} {
		if b519OpaqueViewValue, b519OpaqueViewPresent := b519OpaqueViewContext[b519OpaqueViewKey]; b519OpaqueViewPresent {
			if _, b519OpaqueViewOK := b519OpaqueViewInteger(b519OpaqueViewValue); !b519OpaqueViewOK {
				return nil, fmt.Errorf("%s must be integer", b519OpaqueViewKey)
			}
		}
	}
	for _, b519OpaqueViewKey := range []string{"direction", "view"} {
		if b519OpaqueViewValue, b519OpaqueViewPresent := b519OpaqueViewContext[b519OpaqueViewKey]; b519OpaqueViewPresent {
			if _, b519OpaqueViewOK := b519OpaqueViewValue.(string); !b519OpaqueViewOK {
				return nil, fmt.Errorf("%s must be string", b519OpaqueViewKey)
			}
		}
	}
	b519OpaqueViewFields := map[string]any{}
	b519OpaqueViewOutput := map[string]any{"context": b519OpaqueViewContext, "raw": b519OpaqueViewRaw, "state": "Unknown", "fields": b519OpaqueViewFields, "physical": "Unknown", "qualified": false}
	if b519OpaqueViewArguments, b519OpaqueViewPresent := b519OpaqueViewRoot["arguments"]; b519OpaqueViewPresent {
		b519OpaqueViewObject, b519OpaqueViewOK := b519OpaqueViewArguments.(map[string]any)
		if !b519OpaqueViewOK || len(b519OpaqueViewObject) != 0 {
			return nil, fmt.Errorf("no arguments declared")
		}
		b519OpaqueViewOutput["arguments"] = b519OpaqueViewObject
	}
	b519OpaqueViewPrimary, _ := b519OpaqueViewInteger(b519OpaqueViewContext["primary"])
	b519OpaqueViewSecondary, _ := b519OpaqueViewInteger(b519OpaqueViewContext["secondary"])
	b519OpaqueViewOperation, b519OpaqueViewHasOperation := b519OpaqueViewInteger(b519OpaqueViewContext["operation"])
	b519OpaqueViewDirection := b519OpaqueViewContext["direction"]
	if b519OpaqueViewPrimary != 181 || b519OpaqueViewSecondary != 25 || !b519OpaqueViewHasOperation || ((b519OpaqueViewOperation < 0 || b519OpaqueViewOperation > 8) && b519OpaqueViewOperation != 170) || b519OpaqueViewContext["view"] != "logical-body" || (b519OpaqueViewDirection != "request" && b519OpaqueViewDirection != "positive") {
		return json.Marshal(b519OpaqueViewOutput)
	}
	if b519OpaqueViewDirection == "positive" && b519OpaqueViewOperation != 3 && b519OpaqueViewOperation != 6 && b519OpaqueViewOperation != 170 {
		return json.Marshal(b519OpaqueViewOutput)
	}
	b519OpaqueViewLogical, b519OpaqueViewPresent := b519OpaqueViewRaw["logical_object"].([]any)
	if !b519OpaqueViewPresent || len(b519OpaqueViewLogical) < 3 {
		b519OpaqueViewOutput["state"] = "partial"
		if b519OpaqueViewPresent {
			b519OpaqueViewFields["logical_prefix"] = map[string]any{"raw": b519OpaqueViewLogical, "state": "Unknown"}
		}
		return json.Marshal(b519OpaqueViewOutput)
	}
	b519OpaqueViewOutput["state"] = "Hypothesis"
	if b519OpaqueViewDirection == "positive" && b519OpaqueViewOperation == 6 || b519OpaqueViewDirection == "request" && (b519OpaqueViewOperation == 0 || b519OpaqueViewOperation == 2 || b519OpaqueViewOperation == 170) {
		b519OpaqueViewLabel := "no-declared-argument"
		if b519OpaqueViewDirection == "positive" {
			b519OpaqueViewLabel = "prefix-only"
		}
		b519OpaqueViewFields["form"] = map[string]any{"raw": b519OpaqueViewLogical[:3], "state": "Hypothesis", "label": b519OpaqueViewLabel}
		if len(b519OpaqueViewLogical) > 3 {
			b519OpaqueViewFields["unmodelled"] = map[string]any{"raw": b519OpaqueViewLogical[3:], "state": "Unknown"}
			b519OpaqueViewOutput["state"] = "partial"
		}
		return json.Marshal(b519OpaqueViewOutput)
	}
	b519OpaqueViewOffset := 3
	if b519OpaqueViewDirection == "request" && b519OpaqueViewOperation == 8 {
		if len(b519OpaqueViewLogical) == 3 {
			b519OpaqueViewFields["slave_address"] = map[string]any{"raw": []any{}, "state": "Unknown"}
			b519OpaqueViewOutput["state"] = "partial"
			return json.Marshal(b519OpaqueViewOutput)
		}
		b519OpaqueViewFields["slave_address"] = map[string]any{"raw": b519OpaqueViewLogical[3:4], "state": "Hypothesis"}
		b519OpaqueViewOffset = 4
	}
	b519OpaqueViewFields["opaque"] = map[string]any{"raw": b519OpaqueViewLogical[b519OpaqueViewOffset:], "state": "Hypothesis", "flags": map[string]any{"accepted_empty": "Unknown", "unlimited_carrier": "Unknown"}}
	return json.Marshal(b519OpaqueViewOutput)
}
