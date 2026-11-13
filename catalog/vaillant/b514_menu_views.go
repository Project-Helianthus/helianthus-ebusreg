package vaillant

import (
	b514MenuViewBytes "bytes"
	b514MenuViewJSON "encoding/json"
	b514MenuViewFmt "fmt"
	b514MenuViewIO "io"
	b514MenuViewStrconv "strconv"
	b514MenuViewUTF8 "unicode/utf8"
)

func b514MenuViewRead(b514MenuViewDecoder *b514MenuViewJSON.Decoder, b514MenuViewDepth int) (any, error) {
	if b514MenuViewDepth > 64 {
		return nil, b514MenuViewFmt.Errorf("JSON nesting too deep")
	}
	b514MenuViewToken, b514MenuViewErr := b514MenuViewDecoder.Token()
	if b514MenuViewErr != nil {
		return nil, b514MenuViewErr
	}
	if b514MenuViewDelimiter, b514MenuViewOK := b514MenuViewToken.(b514MenuViewJSON.Delim); b514MenuViewOK {
		switch b514MenuViewDelimiter {
		case '{':
			b514MenuViewObject := map[string]any{}
			for b514MenuViewDecoder.More() {
				b514MenuViewKeyToken, b514MenuViewErr := b514MenuViewDecoder.Token()
				if b514MenuViewErr != nil {
					return nil, b514MenuViewErr
				}
				b514MenuViewKey, b514MenuViewOK := b514MenuViewKeyToken.(string)
				if !b514MenuViewOK {
					return nil, b514MenuViewFmt.Errorf("invalid object key")
				}
				if _, b514MenuViewExists := b514MenuViewObject[b514MenuViewKey]; b514MenuViewExists {
					return nil, b514MenuViewFmt.Errorf("duplicate key %q", b514MenuViewKey)
				}
				b514MenuViewValue, b514MenuViewErr := b514MenuViewRead(b514MenuViewDecoder, b514MenuViewDepth+1)
				if b514MenuViewErr != nil {
					return nil, b514MenuViewErr
				}
				b514MenuViewObject[b514MenuViewKey] = b514MenuViewValue
			}
			if _, b514MenuViewErr := b514MenuViewDecoder.Token(); b514MenuViewErr != nil {
				return nil, b514MenuViewErr
			}
			return b514MenuViewObject, nil
		case '[':
			b514MenuViewArray := []any{}
			for b514MenuViewDecoder.More() {
				b514MenuViewValue, b514MenuViewErr := b514MenuViewRead(b514MenuViewDecoder, b514MenuViewDepth+1)
				if b514MenuViewErr != nil {
					return nil, b514MenuViewErr
				}
				b514MenuViewArray = append(b514MenuViewArray, b514MenuViewValue)
			}
			if _, b514MenuViewErr := b514MenuViewDecoder.Token(); b514MenuViewErr != nil {
				return nil, b514MenuViewErr
			}
			return b514MenuViewArray, nil
		default:
			return nil, b514MenuViewFmt.Errorf("unexpected delimiter")
		}
	}
	return b514MenuViewToken, nil
}

func b514MenuViewInteger(b514MenuViewValue any) (int, bool) {
	b514MenuViewNumber, b514MenuViewOK := b514MenuViewValue.(b514MenuViewJSON.Number)
	if !b514MenuViewOK {
		return 0, false
	}
	b514MenuViewInt, b514MenuViewErr := b514MenuViewStrconv.Atoi(string(b514MenuViewNumber))
	return b514MenuViewInt, b514MenuViewErr == nil
}

func B514MenuViewInterpret(b514MenuViewInput b514MenuViewJSON.RawMessage) (b514MenuViewJSON.RawMessage, error) {
	if !b514MenuViewUTF8.Valid(b514MenuViewInput) {
		return nil, b514MenuViewFmt.Errorf("invalid UTF-8 JSON")
	}
	b514MenuViewDecoder := b514MenuViewJSON.NewDecoder(b514MenuViewBytes.NewReader(b514MenuViewInput))
	b514MenuViewDecoder.UseNumber()
	b514MenuViewValue, b514MenuViewErr := b514MenuViewRead(b514MenuViewDecoder, 0)
	if b514MenuViewErr != nil {
		return nil, b514MenuViewErr
	}
	if _, b514MenuViewErr = b514MenuViewDecoder.Token(); b514MenuViewErr != b514MenuViewIO.EOF {
		return nil, b514MenuViewFmt.Errorf("trailing JSON")
	}
	b514MenuViewInputObject, b514MenuViewOK := b514MenuViewValue.(map[string]any)
	if !b514MenuViewOK {
		return nil, b514MenuViewFmt.Errorf("input must be an object")
	}
	for b514MenuViewKey := range b514MenuViewInputObject {
		if b514MenuViewKey != "context" && b514MenuViewKey != "fields" && b514MenuViewKey != "arguments" {
			return nil, b514MenuViewFmt.Errorf("unknown input key %q", b514MenuViewKey)
		}
	}
	b514MenuViewContext, b514MenuViewOK := b514MenuViewInputObject["context"].(map[string]any)
	if !b514MenuViewOK {
		return nil, b514MenuViewFmt.Errorf("context must be an object")
	}
	b514MenuViewRaw, b514MenuViewOK := b514MenuViewInputObject["fields"].(map[string]any)
	if !b514MenuViewOK {
		return nil, b514MenuViewFmt.Errorf("fields must be an object")
	}
	for b514MenuViewKey, b514MenuViewValue := range b514MenuViewContext {
		switch b514MenuViewKey {
		case "primary", "secondary", "operation":
			if _, b514MenuViewOK := b514MenuViewInteger(b514MenuViewValue); !b514MenuViewOK {
				return nil, b514MenuViewFmt.Errorf("%s must be an integer", b514MenuViewKey)
			}
		case "direction", "view":
			if _, b514MenuViewOK := b514MenuViewValue.(string); !b514MenuViewOK {
				return nil, b514MenuViewFmt.Errorf("%s must be a string", b514MenuViewKey)
			}
		}
	}
	for b514MenuViewKey, b514MenuViewValue := range b514MenuViewRaw {
		b514MenuViewArray, b514MenuViewOK := b514MenuViewValue.([]any)
		if !b514MenuViewOK {
			return nil, b514MenuViewFmt.Errorf("%s must be a byte array", b514MenuViewKey)
		}
		for _, b514MenuViewElement := range b514MenuViewArray {
			b514MenuViewByte, b514MenuViewOK := b514MenuViewInteger(b514MenuViewElement)
			if !b514MenuViewOK || b514MenuViewByte < 0 || b514MenuViewByte > 255 {
				return nil, b514MenuViewFmt.Errorf("invalid byte in %s", b514MenuViewKey)
			}
		}
	}
	b514MenuViewFields := map[string]any{}
	b514MenuViewOutput := map[string]any{"context": b514MenuViewContext, "raw": b514MenuViewRaw, "fields": b514MenuViewFields, "state": "Unknown", "physical": "Unknown", "qualified": false}
	if b514MenuViewArguments, b514MenuViewPresent := b514MenuViewInputObject["arguments"]; b514MenuViewPresent {
		b514MenuViewObject, b514MenuViewOK := b514MenuViewArguments.(map[string]any)
		if !b514MenuViewOK || len(b514MenuViewObject) != 0 {
			return nil, b514MenuViewFmt.Errorf("no arguments are defined for these views")
		}
		b514MenuViewOutput["arguments"] = b514MenuViewObject
	}
	b514MenuViewPrimary, _ := b514MenuViewInteger(b514MenuViewContext["primary"])
	b514MenuViewSecondary, _ := b514MenuViewInteger(b514MenuViewContext["secondary"])
	b514MenuViewOperation, _ := b514MenuViewInteger(b514MenuViewContext["operation"])
	b514MenuViewView, _ := b514MenuViewContext["view"].(string)
	b514MenuViewSelected := b514MenuViewPrimary == 245 && b514MenuViewSecondary == 20 && b514MenuViewContext["direction"] == "positive"
	if !b514MenuViewSelected {
		return b514MenuViewJSON.Marshal(b514MenuViewOutput)
	}
	b514MenuViewAdd := func(b514MenuViewName string, b514MenuViewData []any, b514MenuViewStart, b514MenuViewWidth int) map[string]any {
		b514MenuViewEnd := b514MenuViewStart + b514MenuViewWidth
		if b514MenuViewStart > len(b514MenuViewData) {
			b514MenuViewStart = len(b514MenuViewData)
		}
		if b514MenuViewEnd > len(b514MenuViewData) {
			b514MenuViewEnd = len(b514MenuViewData)
		}
		b514MenuViewField := map[string]any{"raw": b514MenuViewData[b514MenuViewStart:b514MenuViewEnd], "state": "incomplete", "grade": "Hypothesis"}
		b514MenuViewFields[b514MenuViewName] = b514MenuViewField
		return b514MenuViewField
	}
	if (b514MenuViewOperation == 6 && (b514MenuViewView == "generic-result" || b514MenuViewView == "specific-a-result" || b514MenuViewView == "specific-b-result")) || (b514MenuViewOperation == 7 && b514MenuViewView == "progress") {
		b514MenuViewData, b514MenuViewPresent := b514MenuViewRaw["feedback"].([]any)
		if !b514MenuViewPresent {
			b514MenuViewData = []any{}
		}
		b514MenuViewField := b514MenuViewAdd("feedback", b514MenuViewData, 0, 1)
		b514MenuViewOutput["state"] = "partial"
		if len(b514MenuViewData) > 0 {
			b514MenuViewOutput["state"] = "Hypothesis"
			b514MenuViewField["state"] = "unmapped"
			b514MenuViewByte, _ := b514MenuViewInteger(b514MenuViewData[0])
			b514MenuViewLabel := ""
			if b514MenuViewOperation == 6 {
				if b514MenuViewByte == 1 {
					if b514MenuViewView == "generic-result" {
						b514MenuViewLabel = "in-progress"
					} else {
						b514MenuViewLabel = "not-successful"
					}
				}
				if b514MenuViewByte == 255 && b514MenuViewView == "generic-result" {
					b514MenuViewLabel = "unsupported"
				}
				if b514MenuViewByte >= 3 && b514MenuViewView != "generic-result" {
					b514MenuViewLabel = "invalid"
				}
			} else {
				b514MenuViewLabel = map[int]string{0: "idle", 1: "progress", 2: "passed", 3: "warning", 4: "combustion-error", 5: "gas-type-error", 6: "aborted", 255: "unsupported"}[b514MenuViewByte]
			}
			if b514MenuViewLabel != "" {
				b514MenuViewField["label"] = b514MenuViewLabel
				b514MenuViewField["state"] = "Hypothesis"
			}
		}
	} else if b514MenuViewOperation == 8 && b514MenuViewView == "pressure" {
		b514MenuViewData, _ := b514MenuViewRaw["logical_object"].([]any)
		if b514MenuViewData == nil {
			b514MenuViewData = []any{}
		}
		b514MenuViewOutput["state"] = "Hypothesis"
		if len(b514MenuViewData) < 13 {
			b514MenuViewOutput["state"] = "partial"
		}
		for b514MenuViewIndex, b514MenuViewName := range []string{"feedback", "pressure", "minimum", "maximum", "step", "flow"} {
			b514MenuViewStart := []int{3, 4, 6, 8, 10, 11}[b514MenuViewIndex]
			b514MenuViewWidth := []int{1, 2, 2, 2, 1, 2}[b514MenuViewIndex]
			b514MenuViewField := b514MenuViewAdd(b514MenuViewName, b514MenuViewData, b514MenuViewStart, b514MenuViewWidth)
			if len(b514MenuViewData) < b514MenuViewStart+b514MenuViewWidth {
				continue
			}
			b514MenuViewNumber, _ := b514MenuViewInteger(b514MenuViewData[b514MenuViewStart])
			if b514MenuViewWidth == 2 {
				b514MenuViewHigh, _ := b514MenuViewInteger(b514MenuViewData[b514MenuViewStart+1])
				b514MenuViewNumber += b514MenuViewHigh << 8
			}
			b514MenuViewField["state"] = "Hypothesis"
			switch b514MenuViewName {
			case "feedback":
				b514MenuViewField["state"] = "unmapped"
			case "step":
				b514MenuViewField["state"] = "unit-conflict"
				b514MenuViewField["flags"] = []string{"constant-1-mbar-versus-default-0"}
			default:
				b514MenuViewField["calculation"] = "Proven"
				if b514MenuViewName == "flow" {
					b514MenuViewField["number"] = float64(b514MenuViewNumber) / 100
				} else {
					b514MenuViewField["number"] = b514MenuViewNumber
				}
			}
		}
	}
	return b514MenuViewJSON.Marshal(b514MenuViewOutput)
}
