package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

func b523SensorViewRead(b523SensorViewDecoder *json.Decoder) (any, error) {
	b523SensorViewToken, b523SensorViewErr := b523SensorViewDecoder.Token()
	if b523SensorViewErr != nil {
		return nil, b523SensorViewErr
	}
	if b523SensorViewDelimiter, b523SensorViewOK := b523SensorViewToken.(json.Delim); b523SensorViewOK {
		switch b523SensorViewDelimiter {
		case '{':
			b523SensorViewObject := map[string]any{}
			for b523SensorViewDecoder.More() {
				b523SensorViewKeyToken, b523SensorViewErr := b523SensorViewDecoder.Token()
				if b523SensorViewErr != nil {
					return nil, b523SensorViewErr
				}
				b523SensorViewKey, b523SensorViewOK := b523SensorViewKeyToken.(string)
				if !b523SensorViewOK {
					return nil, fmt.Errorf("expected object key")
				}
				if _, b523SensorViewExists := b523SensorViewObject[b523SensorViewKey]; b523SensorViewExists {
					return nil, fmt.Errorf("duplicate object key")
				}
				b523SensorViewValue, b523SensorViewErr := b523SensorViewRead(b523SensorViewDecoder)
				if b523SensorViewErr != nil {
					return nil, b523SensorViewErr
				}
				b523SensorViewObject[b523SensorViewKey] = b523SensorViewValue
			}
			b523SensorViewEnd, b523SensorViewErr := b523SensorViewDecoder.Token()
			if b523SensorViewErr != nil || b523SensorViewEnd != json.Delim('}') {
				return nil, fmt.Errorf("invalid object end")
			}
			return b523SensorViewObject, nil
		case '[':
			b523SensorViewArray := []any{}
			for b523SensorViewDecoder.More() {
				b523SensorViewValue, b523SensorViewErr := b523SensorViewRead(b523SensorViewDecoder)
				if b523SensorViewErr != nil {
					return nil, b523SensorViewErr
				}
				b523SensorViewArray = append(b523SensorViewArray, b523SensorViewValue)
			}
			b523SensorViewEnd, b523SensorViewErr := b523SensorViewDecoder.Token()
			if b523SensorViewErr != nil || b523SensorViewEnd != json.Delim(']') {
				return nil, fmt.Errorf("invalid array end")
			}
			return b523SensorViewArray, nil
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	return b523SensorViewToken, nil
}

func b523SensorViewInteger(b523SensorViewValue any) (int, bool) {
	b523SensorViewNumber, b523SensorViewOK := b523SensorViewValue.(json.Number)
	if !b523SensorViewOK {
		return 0, false
	}
	b523SensorViewInt, b523SensorViewErr := strconv.Atoi(string(b523SensorViewNumber))
	return b523SensorViewInt, b523SensorViewErr == nil
}

func B523SensorViewInterpret(b523SensorViewInput json.RawMessage) (json.RawMessage, error) {
	if !utf8.Valid(b523SensorViewInput) {
		return nil, fmt.Errorf("invalid UTF-8 JSON")
	}
	b523SensorViewDecoder := json.NewDecoder(bytes.NewReader(b523SensorViewInput))
	b523SensorViewDecoder.UseNumber()
	b523SensorViewValue, b523SensorViewErr := b523SensorViewRead(b523SensorViewDecoder)
	if b523SensorViewErr != nil {
		return nil, b523SensorViewErr
	}
	if _, b523SensorViewErr = b523SensorViewDecoder.Token(); b523SensorViewErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b523SensorViewRoot, b523SensorViewOK := b523SensorViewValue.(map[string]any)
	if !b523SensorViewOK {
		return nil, fmt.Errorf("input must be an object")
	}
	for b523SensorViewKey := range b523SensorViewRoot {
		if b523SensorViewKey != "context" && b523SensorViewKey != "fields" && b523SensorViewKey != "arguments" {
			return nil, fmt.Errorf("unknown input member")
		}
	}
	b523SensorViewContext, b523SensorViewOK := b523SensorViewRoot["context"].(map[string]any)
	if !b523SensorViewOK {
		return nil, fmt.Errorf("context must be an object")
	}
	b523SensorViewRaw, b523SensorViewOK := b523SensorViewRoot["fields"].(map[string]any)
	if !b523SensorViewOK {
		return nil, fmt.Errorf("fields must be an object")
	}
	for _, b523SensorViewKey := range []string{"primary", "secondary", "operation"} {
		if b523SensorViewValue, b523SensorViewPresent := b523SensorViewContext[b523SensorViewKey]; b523SensorViewPresent {
			if _, b523SensorViewOK := b523SensorViewInteger(b523SensorViewValue); !b523SensorViewOK {
				return nil, fmt.Errorf("context %s must be an integer", b523SensorViewKey)
			}
		}
	}
	for _, b523SensorViewKey := range []string{"direction", "view"} {
		if b523SensorViewValue, b523SensorViewPresent := b523SensorViewContext[b523SensorViewKey]; b523SensorViewPresent {
			if _, b523SensorViewOK := b523SensorViewValue.(string); !b523SensorViewOK {
				return nil, fmt.Errorf("context %s must be a string", b523SensorViewKey)
			}
		}
	}
	for _, b523SensorViewValue := range b523SensorViewRaw {
		b523SensorViewArray, b523SensorViewOK := b523SensorViewValue.([]any)
		if !b523SensorViewOK {
			return nil, fmt.Errorf("raw field must be a byte array")
		}
		for _, b523SensorViewValue := range b523SensorViewArray {
			b523SensorViewByte, b523SensorViewOK := b523SensorViewInteger(b523SensorViewValue)
			if !b523SensorViewOK || b523SensorViewByte < 0 || b523SensorViewByte > 255 {
				return nil, fmt.Errorf("invalid byte")
			}
		}
	}
	b523SensorViewFields := map[string]any{}
	b523SensorViewOutput := map[string]any{"context": b523SensorViewContext, "raw": b523SensorViewRaw, "state": "Unknown", "fields": b523SensorViewFields, "physical": "Unknown", "qualified": false}
	if b523SensorViewArguments, b523SensorViewPresent := b523SensorViewRoot["arguments"]; b523SensorViewPresent {
		b523SensorViewObject, b523SensorViewOK := b523SensorViewArguments.(map[string]any)
		if !b523SensorViewOK || len(b523SensorViewObject) != 0 {
			return nil, fmt.Errorf("no arguments are defined for this view")
		}
		b523SensorViewOutput["arguments"] = b523SensorViewObject
	}
	b523SensorViewPrimary, _ := b523SensorViewInteger(b523SensorViewContext["primary"])
	b523SensorViewSecondary, _ := b523SensorViewInteger(b523SensorViewContext["secondary"])
	b523SensorViewOperation, _ := b523SensorViewInteger(b523SensorViewContext["operation"])
	b523SensorViewView, _ := b523SensorViewContext["view"].(string)
	if b523SensorViewPrimary != 181 || b523SensorViewSecondary != 35 || b523SensorViewContext["direction"] != "positive" || (b523SensorViewView != "status-zero-valid" && b523SensorViewView != "status-zero-valid-or-low") || (b523SensorViewOperation != 3 && b523SensorViewOperation != 6 && b523SensorViewOperation != 7) {
		return json.Marshal(b523SensorViewOutput)
	}
	b523SensorViewObject, b523SensorViewPresent := b523SensorViewRaw["logical_object"].([]any)
	b523SensorViewOutput["state"] = "Hypothesis"
	b523SensorViewSlice := func(b523SensorViewPosition, b523SensorViewWidth int) ([]any, bool) {
		if b523SensorViewPosition >= len(b523SensorViewObject) {
			b523SensorViewOutput["state"] = "partial"
			return []any{}, false
		}
		b523SensorViewEnd := b523SensorViewPosition + b523SensorViewWidth
		if b523SensorViewEnd > len(b523SensorViewObject) {
			b523SensorViewOutput["state"] = "partial"
			return b523SensorViewObject[b523SensorViewPosition:], false
		}
		return b523SensorViewObject[b523SensorViewPosition:b523SensorViewEnd], true
	}
	b523SensorViewStart, b523SensorViewCount, b523SensorViewStatusPosition := 1, 6, 16
	if b523SensorViewOperation == 6 {
		b523SensorViewCount, b523SensorViewStatusPosition = 7, 17
	}
	if b523SensorViewOperation == 7 {
		b523SensorViewStart = 8
	}
	for b523SensorViewIndex := 0; b523SensorViewIndex < b523SensorViewCount; b523SensorViewIndex++ {
		b523SensorViewPosition := 3 + 2*b523SensorViewIndex
		if b523SensorViewOperation == 7 && b523SensorViewIndex >= 4 {
			b523SensorViewPosition++
		}
		b523SensorViewBytes, b523SensorViewComplete := b523SensorViewSlice(b523SensorViewPosition, 2)
		b523SensorViewField := map[string]any{"raw": b523SensorViewBytes, "state": "incomplete", "grade": "Hypothesis"}
		if b523SensorViewComplete {
			b523SensorViewLow, _ := b523SensorViewInteger(b523SensorViewBytes[0])
			b523SensorViewHigh, _ := b523SensorViewInteger(b523SensorViewBytes[1])
			b523SensorViewWord := b523SensorViewLow | b523SensorViewHigh<<8
			b523SensorViewField["raw"], b523SensorViewField["state"] = b523SensorViewWord, "Hypothesis"
			if b523SensorViewWord == 32768 {
				b523SensorViewField["label"] = "replacement"
			} else {
				b523SensorViewField["number"] = float64(int16(b523SensorViewWord)) / 16
				b523SensorViewField["calculation"] = "Proven"
			}
		}
		b523SensorViewFields[fmt.Sprintf("S%d", b523SensorViewStart+b523SensorViewIndex)] = b523SensorViewField
		b523SensorViewBytes, b523SensorViewComplete = b523SensorViewSlice(b523SensorViewStatusPosition+b523SensorViewIndex/4, 1)
		b523SensorViewField = map[string]any{"raw": b523SensorViewBytes, "state": "incomplete", "grade": "Hypothesis"}
		if b523SensorViewComplete {
			b523SensorViewByte, _ := b523SensorViewInteger(b523SensorViewBytes[0])
			b523SensorViewPair := (b523SensorViewByte >> (2 * (b523SensorViewIndex % 4))) & 3
			b523SensorViewLabels := []string{"valid-value", "open-circuit", "closed-circuit", "high-level"}
			if b523SensorViewView == "status-zero-valid-or-low" {
				b523SensorViewLabels[0] = "valid-value-or-low-level"
			}
			b523SensorViewField["raw"], b523SensorViewField["state"], b523SensorViewField["label"] = b523SensorViewPair, "Hypothesis", b523SensorViewLabels[b523SensorViewPair]
		}
		b523SensorViewFields[fmt.Sprintf("statusS%d", b523SensorViewStart+b523SensorViewIndex)] = b523SensorViewField
	}
	b523SensorViewBytes, b523SensorViewComplete := b523SensorViewSlice(b523SensorViewStatusPosition+1, 1)
	b523SensorViewReserved := map[string]any{"raw": b523SensorViewBytes, "state": "incomplete", "grade": "Hypothesis"}
	if b523SensorViewComplete {
		b523SensorViewByte, _ := b523SensorViewInteger(b523SensorViewBytes[0])
		b523SensorViewMask, b523SensorViewShift := 240, 4
		if b523SensorViewOperation == 6 {
			b523SensorViewMask, b523SensorViewShift = 192, 6
		}
		b523SensorViewReserved["raw"], b523SensorViewReserved["reserved"], b523SensorViewReserved["state"] = b523SensorViewByte, (b523SensorViewByte&b523SensorViewMask)>>b523SensorViewShift, "Hypothesis"
	}
	b523SensorViewFields["status_reserved"] = b523SensorViewReserved
	if b523SensorViewOperation == 3 || b523SensorViewOperation == 7 {
		b523SensorViewPosition, b523SensorViewName := 15, "PWM"
		if b523SensorViewOperation == 7 {
			b523SensorViewPosition, b523SensorViewName = 11, "opaque11"
		}
		b523SensorViewBytes, b523SensorViewComplete := b523SensorViewSlice(b523SensorViewPosition, 1)
		b523SensorViewField := map[string]any{"raw": b523SensorViewBytes, "state": "incomplete", "grade": "Hypothesis"}
		if b523SensorViewComplete {
			b523SensorViewByte, _ := b523SensorViewInteger(b523SensorViewBytes[0])
			b523SensorViewField["raw"] = b523SensorViewByte
			if b523SensorViewOperation == 3 {
				b523SensorViewField["state"], b523SensorViewField["number"], b523SensorViewField["calculation"] = "Hypothesis", float64(b523SensorViewByte)/2, "Proven"
			} else {
				b523SensorViewField["state"] = "Unknown"
			}
		}
		b523SensorViewFields[b523SensorViewName] = b523SensorViewField
	}
	if !b523SensorViewPresent {
		b523SensorViewOutput["state"] = "partial"
	}
	return json.Marshal(b523SensorViewOutput)
}
