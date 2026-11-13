package vaillant

import (
	b50FFieldViewBytes "bytes"
	b50FFieldViewJSON "encoding/json"
	b50FFieldViewFmt "fmt"
	b50FFieldViewIO "io"
	b50FFieldViewStrconv "strconv"
)

const b50FFieldViewMaxDepth = 64

func b50FFieldViewRead(b50FFieldViewDecoder *b50FFieldViewJSON.Decoder, b50FFieldViewDepth int) (any, error) {
	if b50FFieldViewDepth > b50FFieldViewMaxDepth {
		return nil, b50FFieldViewFmt.Errorf("JSON nesting too deep")
	}
	b50FFieldViewToken, b50FFieldViewErr := b50FFieldViewDecoder.Token()
	if b50FFieldViewErr != nil {
		return nil, b50FFieldViewErr
	}
	if b50FFieldViewDelimiter, b50FFieldViewOK := b50FFieldViewToken.(b50FFieldViewJSON.Delim); b50FFieldViewOK {
		switch b50FFieldViewDelimiter {
		case '{':
			b50FFieldViewObject := map[string]any{}
			for b50FFieldViewDecoder.More() {
				b50FFieldViewKeyToken, b50FFieldViewErr := b50FFieldViewDecoder.Token()
				if b50FFieldViewErr != nil {
					return nil, b50FFieldViewErr
				}
				b50FFieldViewKey, b50FFieldViewOK := b50FFieldViewKeyToken.(string)
				if !b50FFieldViewOK {
					return nil, b50FFieldViewFmt.Errorf("invalid object key")
				}
				if _, b50FFieldViewFound := b50FFieldViewObject[b50FFieldViewKey]; b50FFieldViewFound {
					return nil, b50FFieldViewFmt.Errorf("duplicate object key")
				}
				b50FFieldViewValue, b50FFieldViewErr := b50FFieldViewRead(b50FFieldViewDecoder, b50FFieldViewDepth+1)
				if b50FFieldViewErr != nil {
					return nil, b50FFieldViewErr
				}
				b50FFieldViewObject[b50FFieldViewKey] = b50FFieldViewValue
			}
			_, b50FFieldViewErr := b50FFieldViewDecoder.Token()
			return b50FFieldViewObject, b50FFieldViewErr
		case '[':
			b50FFieldViewArray := []any{}
			for b50FFieldViewDecoder.More() {
				b50FFieldViewValue, b50FFieldViewErr := b50FFieldViewRead(b50FFieldViewDecoder, b50FFieldViewDepth+1)
				if b50FFieldViewErr != nil {
					return nil, b50FFieldViewErr
				}
				b50FFieldViewArray = append(b50FFieldViewArray, b50FFieldViewValue)
			}
			_, b50FFieldViewErr := b50FFieldViewDecoder.Token()
			return b50FFieldViewArray, b50FFieldViewErr
		}
		return nil, b50FFieldViewFmt.Errorf("unexpected delimiter")
	}
	return b50FFieldViewToken, nil
}

func B50FFieldViewInterpret(b50FFieldViewInput b50FFieldViewJSON.RawMessage) (b50FFieldViewJSON.RawMessage, error) {
	b50FFieldViewDecoder := b50FFieldViewJSON.NewDecoder(b50FFieldViewBytes.NewReader(b50FFieldViewInput))
	b50FFieldViewDecoder.UseNumber()
	b50FFieldViewValue, b50FFieldViewErr := b50FFieldViewRead(b50FFieldViewDecoder, 0)
	if b50FFieldViewErr != nil {
		return nil, b50FFieldViewErr
	}
	if _, b50FFieldViewErr = b50FFieldViewDecoder.Token(); b50FFieldViewErr != b50FFieldViewIO.EOF {
		return nil, b50FFieldViewFmt.Errorf("trailing JSON")
	}
	b50FFieldViewInputObject, b50FFieldViewOK := b50FFieldViewValue.(map[string]any)
	if !b50FFieldViewOK {
		return nil, b50FFieldViewFmt.Errorf("input must be an object")
	}
	for b50FFieldViewKey := range b50FFieldViewInputObject {
		if b50FFieldViewKey != "context" && b50FFieldViewKey != "fields" && b50FFieldViewKey != "arguments" {
			return nil, b50FFieldViewFmt.Errorf("unknown input key")
		}
	}
	b50FFieldViewContext, b50FFieldViewOK := b50FFieldViewInputObject["context"].(map[string]any)
	if !b50FFieldViewOK {
		return nil, b50FFieldViewFmt.Errorf("context must be an object")
	}
	b50FFieldViewRaw, b50FFieldViewOK := b50FFieldViewInputObject["fields"].(map[string]any)
	if !b50FFieldViewOK {
		return nil, b50FFieldViewFmt.Errorf("fields must be an object")
	}
	for _, b50FFieldViewKey := range []string{"primary", "secondary", "operation"} {
		if b50FFieldViewValue, b50FFieldViewPresent := b50FFieldViewContext[b50FFieldViewKey]; b50FFieldViewPresent {
			b50FFieldViewNumber, b50FFieldViewOK := b50FFieldViewValue.(b50FFieldViewJSON.Number)
			if !b50FFieldViewOK {
				return nil, b50FFieldViewFmt.Errorf("integer selector required")
			}
			if _, b50FFieldViewErr := b50FFieldViewStrconv.ParseInt(string(b50FFieldViewNumber), 10, 64); b50FFieldViewErr != nil && !b50FFieldViewIntegerOutOfRange(string(b50FFieldViewNumber)) {
				return nil, b50FFieldViewFmt.Errorf("integer selector required")
			}
		}
	}
	for _, b50FFieldViewKey := range []string{"direction", "view"} {
		if b50FFieldViewValue, b50FFieldViewPresent := b50FFieldViewContext[b50FFieldViewKey]; b50FFieldViewPresent {
			if _, b50FFieldViewOK := b50FFieldViewValue.(string); !b50FFieldViewOK {
				return nil, b50FFieldViewFmt.Errorf("string selector required")
			}
		}
	}
	b50FFieldViewArrays := map[string][]int{}
	for b50FFieldViewKey, b50FFieldViewValue := range b50FFieldViewRaw {
		b50FFieldViewArray, b50FFieldViewOK := b50FFieldViewValue.([]any)
		if !b50FFieldViewOK {
			return nil, b50FFieldViewFmt.Errorf("field must be a byte array")
		}
		b50FFieldViewData := []int{}
		for _, b50FFieldViewElement := range b50FFieldViewArray {
			b50FFieldViewNumber, b50FFieldViewOK := b50FFieldViewElement.(b50FFieldViewJSON.Number)
			if !b50FFieldViewOK {
				return nil, b50FFieldViewFmt.Errorf("invalid byte")
			}
			b50FFieldViewByte, b50FFieldViewErr := b50FFieldViewStrconv.Atoi(string(b50FFieldViewNumber))
			if b50FFieldViewErr != nil || b50FFieldViewByte < 0 || b50FFieldViewByte > 255 || string(b50FFieldViewNumber) == "-0" {
				return nil, b50FFieldViewFmt.Errorf("invalid byte")
			}
			b50FFieldViewData = append(b50FFieldViewData, b50FFieldViewByte)
		}
		b50FFieldViewArrays[b50FFieldViewKey] = b50FFieldViewData
	}
	b50FFieldViewFields := map[string]any{}
	for b50FFieldViewKey, b50FFieldViewValue := range b50FFieldViewRaw {
		b50FFieldViewFields[b50FFieldViewKey] = map[string]any{"raw": b50FFieldViewValue, "state": "Unknown"}
	}
	b50FFieldViewOutput := map[string]any{"context": b50FFieldViewContext, "raw": b50FFieldViewRaw, "fields": b50FFieldViewFields, "state": "Unknown", "physical": "Unknown", "qualified": false}
	if b50FFieldViewValue, b50FFieldViewPresent := b50FFieldViewInputObject["arguments"]; b50FFieldViewPresent {
		b50FFieldViewArguments, b50FFieldViewOK := b50FFieldViewValue.(map[string]any)
		if !b50FFieldViewOK || len(b50FFieldViewArguments) != 0 {
			return nil, b50FFieldViewFmt.Errorf("no arguments defined for this view")
		}
		b50FFieldViewOutput["arguments"] = b50FFieldViewArguments
	}
	b50FFieldViewOperation := b50FFieldViewContext["operation"]
	b50FFieldViewView := b50FFieldViewContext["view"]
	if b50FFieldViewContext["primary"] != b50FFieldViewJSON.Number("181") || b50FFieldViewContext["secondary"] != b50FFieldViewJSON.Number("15") || b50FFieldViewContext["direction"] != "positive" {
		return b50FFieldViewJSON.Marshal(b50FFieldViewOutput)
	}
	b50FFieldViewSelected := (b50FFieldViewOperation == b50FFieldViewJSON.Number("0") || b50FFieldViewOperation == b50FFieldViewJSON.Number("1")) && b50FFieldViewView == "thermal-fields" || b50FFieldViewOperation == b50FFieldViewJSON.Number("3") && b50FFieldViewView == "demand" || b50FFieldViewOperation == b50FFieldViewJSON.Number("2") && b50FFieldViewView == "integrity-only"
	if !b50FFieldViewSelected {
		return b50FFieldViewJSON.Marshal(b50FFieldViewOutput)
	}
	b50FFieldViewOutput["state"] = "complete"
	b50FFieldViewField := func(b50FFieldViewKey string, b50FFieldViewWidth int) (map[string]any, []int, bool) {
		b50FFieldViewData, b50FFieldViewPresent := b50FFieldViewArrays[b50FFieldViewKey]
		b50FFieldViewResult := map[string]any{"state": "missing", "grade": "Hypothesis"}
		b50FFieldViewFields[b50FFieldViewKey] = b50FFieldViewResult
		if b50FFieldViewPresent {
			b50FFieldViewResult["raw"] = b50FFieldViewRaw[b50FFieldViewKey]
			b50FFieldViewResult["state"] = "incomplete"
		}
		if !b50FFieldViewPresent || len(b50FFieldViewData) != b50FFieldViewWidth {
			b50FFieldViewOutput["state"] = "partial"
			if len(b50FFieldViewData) > b50FFieldViewWidth {
				b50FFieldViewResult["state"] = "unmapped"
			}
			return b50FFieldViewResult, b50FFieldViewData, false
		}
		b50FFieldViewResult["state"] = "complete"
		return b50FFieldViewResult, b50FFieldViewData, true
	}
	if b50FFieldViewView == "thermal-fields" {
		for b50FFieldViewIndex, b50FFieldViewSuffix := range []string{"one", "two", "three"} {
			b50FFieldViewKey := "temperature_" + b50FFieldViewSuffix
			if b50FFieldViewIndex == 2 && b50FFieldViewOperation == b50FFieldViewJSON.Number("1") {
				b50FFieldViewKey = "stack_voltage"
			}
			b50FFieldViewResult, b50FFieldViewData, b50FFieldViewComplete := b50FFieldViewField(b50FFieldViewKey, 2)
			if b50FFieldViewComplete {
				if b50FFieldViewKey == "stack_voltage" {
					b50FFieldViewResult["state"] = "unit-conflict"
					b50FFieldViewOutput["state"] = "partial"
				} else {
					b50FFieldViewResult["number"] = float64(int16(uint16(b50FFieldViewData[0])|uint16(b50FFieldViewData[1])<<8)) / 16
					b50FFieldViewResult["label"] = "Celsius"
					b50FFieldViewResult["calculation"] = "Proven"
				}
			}
			b50FFieldViewField("error_"+b50FFieldViewSuffix, 1)
		}
	}
	if b50FFieldViewView == "demand" {
		b50FFieldViewResult, b50FFieldViewData, b50FFieldViewComplete := b50FFieldViewField("demand", 1)
		if b50FFieldViewComplete {
			b50FFieldViewFlags := map[string]bool{}
			for b50FFieldViewIndex, b50FFieldViewName := range []string{"burner", "reformer-gas-valve", "gas-flow", "fan-power", "valve"} {
				b50FFieldViewFlags[b50FFieldViewName] = b50FFieldViewData[0]&(1<<b50FFieldViewIndex) != 0
			}
			b50FFieldViewResult["flags"] = b50FFieldViewFlags
			b50FFieldViewResult["reserved"] = b50FFieldViewData[0] & 224
			b50FFieldViewResult["label"] = "demand"
		}
		b50FFieldViewResult, b50FFieldViewData, b50FFieldViewComplete = b50FFieldViewField("feedback", 1)
		if b50FFieldViewComplete {
			b50FFieldViewResult["label"] = "rejected"
			if b50FFieldViewData[0] == 1 {
				b50FFieldViewResult["label"] = "accepted"
			}
		}
	}
	for b50FFieldViewIndex, b50FFieldViewKey := range []string{"integrity_one", "integrity_two"} {
		if _, b50FFieldViewPresent := b50FFieldViewArrays[b50FFieldViewKey]; !b50FFieldViewPresent && b50FFieldViewView != "integrity-only" {
			continue
		}
		b50FFieldViewResult, _, b50FFieldViewComplete := b50FFieldViewField(b50FFieldViewKey, 1)
		if b50FFieldViewComplete {
			b50FFieldViewPosition := 12
			if b50FFieldViewOperation == b50FFieldViewJSON.Number("2") {
				b50FFieldViewPosition = 11
			}
			if b50FFieldViewOperation == b50FFieldViewJSON.Number("3") {
				b50FFieldViewPosition = 4
			}
			b50FFieldViewResult["logical_position"] = b50FFieldViewPosition + b50FFieldViewIndex
			b50FFieldViewResult["state"] = "unknown-order"
			b50FFieldViewResult["word_order"] = "Unknown"
			b50FFieldViewResult["crc"] = "Unknown"
			b50FFieldViewOutput["state"] = "partial"
		}
	}
	for _, b50FFieldViewValue := range b50FFieldViewFields {
		if b50FFieldViewValue.(map[string]any)["state"] == "Unknown" {
			b50FFieldViewOutput["state"] = "partial"
		}
	}
	return b50FFieldViewJSON.Marshal(b50FFieldViewOutput)
}

// b50FFieldViewIntegerOutOfRange reports a plain integer literal too large for int64;
// such a selector is an unknown value, not a wrong kind.
func b50FFieldViewIntegerOutOfRange(b50FFieldViewText string) bool {
	b50FFieldViewDigits := b50FFieldViewText
	if len(b50FFieldViewDigits) > 0 && b50FFieldViewDigits[0] == '-' {
		b50FFieldViewDigits = b50FFieldViewDigits[1:]
	}
	if b50FFieldViewDigits == "" {
		return false
	}
	for _, b50FFieldViewChar := range b50FFieldViewDigits {
		if b50FFieldViewChar < '0' || b50FFieldViewChar > '9' {
			return false
		}
	}
	return true
}
