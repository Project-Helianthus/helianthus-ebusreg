package vaillant

import (
	b512FeedbackViewBytesPackage "bytes"
	b512FeedbackViewJSON "encoding/json"
	b512FeedbackViewFmt "fmt"
	b512FeedbackViewIO "io"
)

func B512FeedbackViewInterpret(b512FeedbackViewInput b512FeedbackViewJSON.RawMessage) (b512FeedbackViewJSON.RawMessage, error) {
	b512FeedbackViewDecoder := b512FeedbackViewJSON.NewDecoder(b512FeedbackViewBytesPackage.NewReader(b512FeedbackViewInput))
	b512FeedbackViewDecoder.UseNumber()
	if b512FeedbackViewErr := b512FeedbackViewValidate(b512FeedbackViewDecoder); b512FeedbackViewErr != nil {
		return nil, b512FeedbackViewErr
	}
	if _, b512FeedbackViewErr := b512FeedbackViewDecoder.Token(); b512FeedbackViewErr != b512FeedbackViewIO.EOF {
		return nil, b512FeedbackViewFmt.Errorf("ab15ao15 trailing JSON")
	}
	b512FeedbackViewRoot, b512FeedbackViewErr := b512FeedbackViewObject(b512FeedbackViewInput)
	if b512FeedbackViewErr != nil {
		return nil, b512FeedbackViewErr
	}
	for b512FeedbackViewKey := range b512FeedbackViewRoot {
		if b512FeedbackViewKey != "context" && b512FeedbackViewKey != "fields" && b512FeedbackViewKey != "arguments" {
			return nil, b512FeedbackViewFmt.Errorf("ab15ao15 unknown input key")
		}
	}
	b512FeedbackViewContext, b512FeedbackViewErr := b512FeedbackViewObject(b512FeedbackViewRoot["context"])
	if b512FeedbackViewErr != nil {
		return nil, b512FeedbackViewErr
	}
	b512FeedbackViewRaw, b512FeedbackViewErr := b512FeedbackViewObject(b512FeedbackViewRoot["fields"])
	if b512FeedbackViewErr != nil {
		return nil, b512FeedbackViewErr
	}
	b512FeedbackViewBytes := make(map[string][]int)
	for b512FeedbackViewKey, b512FeedbackViewValue := range b512FeedbackViewRaw {
		if len(b512FeedbackViewBytesPackage.TrimSpace(b512FeedbackViewValue)) == 0 || b512FeedbackViewBytesPackage.TrimSpace(b512FeedbackViewValue)[0] != '[' {
			return nil, b512FeedbackViewFmt.Errorf("ab15ao15 field must be byte array")
		}
		var b512FeedbackViewElements []b512FeedbackViewJSON.RawMessage
		if b512FeedbackViewErr := b512FeedbackViewJSON.Unmarshal(b512FeedbackViewValue, &b512FeedbackViewElements); b512FeedbackViewErr != nil {
			return nil, b512FeedbackViewErr
		}
		b512FeedbackViewValues := make([]int, len(b512FeedbackViewElements))
		for b512FeedbackViewIndex, b512FeedbackViewElement := range b512FeedbackViewElements {
			if b512FeedbackViewBytesPackage.Equal(b512FeedbackViewElement, []byte("null")) {
				return nil, b512FeedbackViewFmt.Errorf("ab15ao15 null byte")
			}
			if b512FeedbackViewErr := b512FeedbackViewJSON.Unmarshal(b512FeedbackViewElement, &b512FeedbackViewValues[b512FeedbackViewIndex]); b512FeedbackViewErr != nil || b512FeedbackViewValues[b512FeedbackViewIndex] < 0 || b512FeedbackViewValues[b512FeedbackViewIndex] > 255 {
				return nil, b512FeedbackViewFmt.Errorf("ab15ao15 invalid byte")
			}
		}
		b512FeedbackViewBytes[b512FeedbackViewKey] = b512FeedbackViewValues
	}
	b512FeedbackViewSelected := true
	for b512FeedbackViewKey, b512FeedbackViewExpected := range map[string]int{"primary": 181, "secondary": 18, "operation": 1} {
		b512FeedbackViewValue, b512FeedbackViewPresent := b512FeedbackViewContext[b512FeedbackViewKey]
		if !b512FeedbackViewPresent {
			b512FeedbackViewSelected = false
			continue
		}
		var b512FeedbackViewNumber int
		if b512FeedbackViewBytesPackage.Equal(b512FeedbackViewValue, []byte("null")) {
			return nil, b512FeedbackViewFmt.Errorf("ab15ao15 null identity")
		}
		if b512FeedbackViewErr := b512FeedbackViewJSON.Unmarshal(b512FeedbackViewValue, &b512FeedbackViewNumber); b512FeedbackViewErr != nil {
			return nil, b512FeedbackViewErr
		}
		if b512FeedbackViewNumber != b512FeedbackViewExpected {
			b512FeedbackViewSelected = false
		}
	}
	b512FeedbackViewDirection := ""
	for b512FeedbackViewKey := range map[string]bool{"direction": true, "view": true} {
		b512FeedbackViewValue, b512FeedbackViewPresent := b512FeedbackViewContext[b512FeedbackViewKey]
		if !b512FeedbackViewPresent {
			b512FeedbackViewSelected = false
			continue
		}
		var b512FeedbackViewText string
		if b512FeedbackViewBytesPackage.Equal(b512FeedbackViewValue, []byte("null")) {
			return nil, b512FeedbackViewFmt.Errorf("ab15ao15 null selector")
		}
		if b512FeedbackViewErr := b512FeedbackViewJSON.Unmarshal(b512FeedbackViewValue, &b512FeedbackViewText); b512FeedbackViewErr != nil {
			return nil, b512FeedbackViewErr
		}
		if b512FeedbackViewKey == "direction" {
			b512FeedbackViewDirection = b512FeedbackViewText
			if b512FeedbackViewText != "request" && b512FeedbackViewText != "positive" {
				b512FeedbackViewSelected = false
			}
		} else if b512FeedbackViewText != "directional-fields" {
			b512FeedbackViewSelected = false
		}
	}
	b512FeedbackViewOrigin := "origin-Unknown"
	if b512FeedbackViewArguments, b512FeedbackViewPresent := b512FeedbackViewRoot["arguments"]; b512FeedbackViewPresent {
		b512FeedbackViewFlags, b512FeedbackViewErr := b512FeedbackViewObject(b512FeedbackViewArguments)
		if b512FeedbackViewErr != nil {
			return nil, b512FeedbackViewErr
		}
		for b512FeedbackViewKey, b512FeedbackViewValue := range b512FeedbackViewFlags {
			if b512FeedbackViewKey != "actual_signal" || (string(b512FeedbackViewValue) != "true" && string(b512FeedbackViewValue) != "false") {
				return nil, b512FeedbackViewFmt.Errorf("ab15ao15 invalid argument")
			}
			if string(b512FeedbackViewValue) == "true" {
				b512FeedbackViewOrigin = "actual-input-candidate"
			} else {
				b512FeedbackViewOrigin = "fallback-candidate"
			}
		}
	}
	b512FeedbackViewFields := make(map[string]any)
	b512FeedbackViewState := "Unknown"
	if b512FeedbackViewSelected {
		b512FeedbackViewState = "Hypothesis"
		for _, b512FeedbackViewName := range []string{"secondary_heating_pump", "circulation_pump", "solar_pump", "error_relay", "loading_pump", "fume_flap"} {
			b512FeedbackViewField := map[string]any{"raw": nil, "state": "Unknown"}
			if b512FeedbackViewValue, b512FeedbackViewPresent := b512FeedbackViewRaw[b512FeedbackViewName]; b512FeedbackViewPresent {
				b512FeedbackViewField["raw"] = b512FeedbackViewValue
			}
			if len(b512FeedbackViewBytes[b512FeedbackViewName]) != 1 {
				b512FeedbackViewState = "partial"
				b512FeedbackViewFields[b512FeedbackViewName] = b512FeedbackViewField
				continue
			}
			b512FeedbackViewField["state"] = "Hypothesis"
			b512FeedbackViewField["flags"] = map[string]string{"origin": b512FeedbackViewOrigin}
			b512FeedbackViewByte := b512FeedbackViewBytes[b512FeedbackViewName][0]
			if b512FeedbackViewName == "error_relay" || b512FeedbackViewName == "fume_flap" {
				switch b512FeedbackViewByte {
				case 15:
					b512FeedbackViewField["label"] = "state-0F"
				case 240:
					b512FeedbackViewField["label"] = "state-F0"
				default:
					b512FeedbackViewField["label"] = "no-action-candidate"
				}
			} else if b512FeedbackViewByte == 255 {
				if b512FeedbackViewDirection == "request" {
					b512FeedbackViewField["label"] = "leave-unchanged"
				} else {
					b512FeedbackViewField["label"] = "function-unavailable"
				}
			}
			b512FeedbackViewFields[b512FeedbackViewName] = b512FeedbackViewField
		}
	}
	b512FeedbackViewOutput := map[string]any{"context": b512FeedbackViewRoot["context"], "raw": b512FeedbackViewRoot["fields"], "state": b512FeedbackViewState, "fields": b512FeedbackViewFields, "physical": "Unknown", "qualified": false}
	if b512FeedbackViewArguments, b512FeedbackViewPresent := b512FeedbackViewRoot["arguments"]; b512FeedbackViewPresent {
		b512FeedbackViewOutput["arguments"] = b512FeedbackViewArguments
	}
	return b512FeedbackViewJSON.Marshal(b512FeedbackViewOutput)
}

func b512FeedbackViewObject(b512FeedbackViewRaw b512FeedbackViewJSON.RawMessage) (map[string]b512FeedbackViewJSON.RawMessage, error) {
	if len(b512FeedbackViewBytesPackage.TrimSpace(b512FeedbackViewRaw)) == 0 || b512FeedbackViewBytesPackage.TrimSpace(b512FeedbackViewRaw)[0] != '{' {
		return nil, b512FeedbackViewFmt.Errorf("ab15ao15 expected object")
	}
	var b512FeedbackViewObject map[string]b512FeedbackViewJSON.RawMessage
	b512FeedbackViewErr := b512FeedbackViewJSON.Unmarshal(b512FeedbackViewRaw, &b512FeedbackViewObject)
	return b512FeedbackViewObject, b512FeedbackViewErr
}

func b512FeedbackViewValidate(b512FeedbackViewDecoder *b512FeedbackViewJSON.Decoder) error {
	b512FeedbackViewToken, b512FeedbackViewErr := b512FeedbackViewDecoder.Token()
	if b512FeedbackViewErr != nil {
		return b512FeedbackViewErr
	}
	b512FeedbackViewDelimiter, b512FeedbackViewContainer := b512FeedbackViewToken.(b512FeedbackViewJSON.Delim)
	if !b512FeedbackViewContainer {
		return nil
	}
	switch b512FeedbackViewDelimiter {
	case '{':
		b512FeedbackViewSeen := make(map[string]bool)
		for b512FeedbackViewDecoder.More() {
			b512FeedbackViewKeyToken, b512FeedbackViewErr := b512FeedbackViewDecoder.Token()
			if b512FeedbackViewErr != nil {
				return b512FeedbackViewErr
			}
			b512FeedbackViewKey, b512FeedbackViewOK := b512FeedbackViewKeyToken.(string)
			if !b512FeedbackViewOK || b512FeedbackViewSeen[b512FeedbackViewKey] {
				return b512FeedbackViewFmt.Errorf("ab15ao15 duplicate or invalid key")
			}
			b512FeedbackViewSeen[b512FeedbackViewKey] = true
			if b512FeedbackViewErr := b512FeedbackViewValidate(b512FeedbackViewDecoder); b512FeedbackViewErr != nil {
				return b512FeedbackViewErr
			}
		}
	case '[':
		for b512FeedbackViewDecoder.More() {
			if b512FeedbackViewErr := b512FeedbackViewValidate(b512FeedbackViewDecoder); b512FeedbackViewErr != nil {
				return b512FeedbackViewErr
			}
		}
	default:
		return b512FeedbackViewFmt.Errorf("ab15ao15 unexpected delimiter")
	}
	_, b512FeedbackViewErr = b512FeedbackViewDecoder.Token()
	return b512FeedbackViewErr
}
