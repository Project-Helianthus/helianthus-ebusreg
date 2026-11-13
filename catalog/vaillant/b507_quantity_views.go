package vaillant

import (
	b507QuantityViewBytesPackage "bytes"
	b507QuantityViewJSON "encoding/json"
	b507QuantityViewFmt "fmt"
	b507QuantityViewIO "io"
	b507QuantityViewStrconv "strconv"
	b507QuantityViewStrings "strings"
)

func b507QuantityViewValue(b507QuantityViewDecoder *b507QuantityViewJSON.Decoder) (any, error) {
	b507QuantityViewToken, b507QuantityViewError := b507QuantityViewDecoder.Token()
	if b507QuantityViewError != nil {
		return nil, b507QuantityViewError
	}
	switch b507QuantityViewToken {
	case b507QuantityViewJSON.Delim('{'):
		b507QuantityViewObject := map[string]any{}
		for b507QuantityViewDecoder.More() {
			b507QuantityViewKeyToken, b507QuantityViewError := b507QuantityViewDecoder.Token()
			if b507QuantityViewError != nil {
				return nil, b507QuantityViewError
			}
			b507QuantityViewKey := b507QuantityViewKeyToken.(string)
			if _, b507QuantityViewExists := b507QuantityViewObject[b507QuantityViewKey]; b507QuantityViewExists {
				return nil, b507QuantityViewFmt.Errorf("duplicate key %q", b507QuantityViewKey)
			}
			b507QuantityViewChild, b507QuantityViewError := b507QuantityViewValue(b507QuantityViewDecoder)
			if b507QuantityViewError != nil {
				return nil, b507QuantityViewError
			}
			b507QuantityViewObject[b507QuantityViewKey] = b507QuantityViewChild
		}
		_, b507QuantityViewError = b507QuantityViewDecoder.Token()
		return b507QuantityViewObject, b507QuantityViewError
	case b507QuantityViewJSON.Delim('['):
		b507QuantityViewArray := []any{}
		for b507QuantityViewDecoder.More() {
			b507QuantityViewChild, b507QuantityViewError := b507QuantityViewValue(b507QuantityViewDecoder)
			if b507QuantityViewError != nil {
				return nil, b507QuantityViewError
			}
			b507QuantityViewArray = append(b507QuantityViewArray, b507QuantityViewChild)
		}
		_, b507QuantityViewError = b507QuantityViewDecoder.Token()
		return b507QuantityViewArray, b507QuantityViewError
	default:
		return b507QuantityViewToken, nil
	}
}

func B507QuantityViewInterpret(b507QuantityViewInput b507QuantityViewJSON.RawMessage) (b507QuantityViewJSON.RawMessage, error) {
	b507QuantityViewDecoder := b507QuantityViewJSON.NewDecoder(b507QuantityViewBytesPackage.NewReader(b507QuantityViewInput))
	b507QuantityViewDecoder.UseNumber()
	b507QuantityViewValueRoot, b507QuantityViewError := b507QuantityViewValue(b507QuantityViewDecoder)
	if b507QuantityViewError != nil {
		return nil, b507QuantityViewError
	}
	if _, b507QuantityViewError = b507QuantityViewDecoder.Token(); b507QuantityViewError != b507QuantityViewIO.EOF {
		return nil, b507QuantityViewFmt.Errorf("trailing JSON")
	}
	b507QuantityViewRoot, b507QuantityViewOK := b507QuantityViewValueRoot.(map[string]any)
	if !b507QuantityViewOK {
		return nil, b507QuantityViewFmt.Errorf("input must be an object")
	}
	for b507QuantityViewKey := range b507QuantityViewRoot {
		if b507QuantityViewKey != "context" && b507QuantityViewKey != "fields" && b507QuantityViewKey != "arguments" {
			return nil, b507QuantityViewFmt.Errorf("unknown input key %q", b507QuantityViewKey)
		}
	}
	b507QuantityViewContext, b507QuantityViewOK := b507QuantityViewRoot["context"].(map[string]any)
	if !b507QuantityViewOK {
		return nil, b507QuantityViewFmt.Errorf("context must be an object")
	}
	b507QuantityViewRaw, b507QuantityViewOK := b507QuantityViewRoot["fields"].(map[string]any)
	if !b507QuantityViewOK {
		return nil, b507QuantityViewFmt.Errorf("fields must be an object")
	}
	for _, b507QuantityViewKey := range []string{"primary", "secondary", "operation"} {
		if b507QuantityViewItem, b507QuantityViewPresent := b507QuantityViewContext[b507QuantityViewKey]; b507QuantityViewPresent {
			b507QuantityViewNumber, b507QuantityViewOK := b507QuantityViewItem.(b507QuantityViewJSON.Number)
			if !b507QuantityViewOK {
				return nil, b507QuantityViewFmt.Errorf("%s must be an integer", b507QuantityViewKey)
			}
			if b507QuantityViewStrings.ContainsAny(string(b507QuantityViewNumber), ".eE") {
				return nil, b507QuantityViewFmt.Errorf("%s must be an integer", b507QuantityViewKey)
			}
		}
	}
	for _, b507QuantityViewKey := range []string{"direction", "view"} {
		if b507QuantityViewItem, b507QuantityViewPresent := b507QuantityViewContext[b507QuantityViewKey]; b507QuantityViewPresent {
			if _, b507QuantityViewOK := b507QuantityViewItem.(string); !b507QuantityViewOK {
				return nil, b507QuantityViewFmt.Errorf("%s must be a string", b507QuantityViewKey)
			}
		}
	}
	b507QuantityViewResultFields := map[string]any{}
	for b507QuantityViewKey, b507QuantityViewItem := range b507QuantityViewRaw {
		b507QuantityViewArray, b507QuantityViewOK := b507QuantityViewItem.([]any)
		if !b507QuantityViewOK {
			return nil, b507QuantityViewFmt.Errorf("%s must be a byte array", b507QuantityViewKey)
		}
		for _, b507QuantityViewItem := range b507QuantityViewArray {
			b507QuantityViewNumber, b507QuantityViewOK := b507QuantityViewItem.(b507QuantityViewJSON.Number)
			if !b507QuantityViewOK {
				return nil, b507QuantityViewFmt.Errorf("invalid byte in %s", b507QuantityViewKey)
			}
			if _, b507QuantityViewError := b507QuantityViewStrconv.ParseUint(string(b507QuantityViewNumber), 10, 8); b507QuantityViewError != nil {
				return nil, b507QuantityViewFmt.Errorf("invalid byte in %s", b507QuantityViewKey)
			}
		}
		b507QuantityViewResultFields[b507QuantityViewKey] = map[string]any{"raw": b507QuantityViewArray, "state": "Unknown"}
	}
	b507QuantityViewOutput := map[string]any{"context": b507QuantityViewContext, "raw": b507QuantityViewRaw, "state": "Unknown", "fields": b507QuantityViewResultFields, "physical": "Unknown", "qualified": false}
	if b507QuantityViewArguments, b507QuantityViewPresent := b507QuantityViewRoot["arguments"]; b507QuantityViewPresent {
		b507QuantityViewObject, b507QuantityViewOK := b507QuantityViewArguments.(map[string]any)
		if !b507QuantityViewOK || len(b507QuantityViewObject) != 0 {
			return nil, b507QuantityViewFmt.Errorf("arguments must be an empty object: no argument flags are defined")
		}
		b507QuantityViewOutput["arguments"] = b507QuantityViewObject
	}
	b507QuantityViewSelected := ""
	if b507QuantityViewContext["primary"] == b507QuantityViewJSON.Number("181") && b507QuantityViewContext["secondary"] == b507QuantityViewJSON.Number("7") && b507QuantityViewContext["direction"] == "positive" {
		if b507QuantityViewContext["operation"] == b507QuantityViewJSON.Number("6") && b507QuantityViewContext["view"] == "capacity" {
			b507QuantityViewSelected = "capacity"
		}
		if b507QuantityViewContext["operation"] == b507QuantityViewJSON.Number("8") && b507QuantityViewContext["view"] == "preference" {
			b507QuantityViewSelected = "preference"
		}
	}
	if b507QuantityViewSelected != "" {
		b507QuantityViewField := map[string]any{"state": "incomplete", "grade": "Hypothesis"}
		b507QuantityViewResultFields[b507QuantityViewSelected] = b507QuantityViewField
		b507QuantityViewOutput["state"] = "partial"
		b507QuantityViewArray, b507QuantityViewPresent := b507QuantityViewRaw[b507QuantityViewSelected].([]any)
		if b507QuantityViewPresent {
			b507QuantityViewField["raw"] = b507QuantityViewArray
		}
		b507QuantityViewWidth := 1
		if b507QuantityViewSelected == "preference" {
			b507QuantityViewWidth = 2
		}
		if b507QuantityViewPresent && len(b507QuantityViewArray) == b507QuantityViewWidth {
			b507QuantityViewField["state"] = "Hypothesis"
			b507QuantityViewOutput["state"] = "Hypothesis"
			b507QuantityViewFirst, _ := b507QuantityViewStrconv.Atoi(string(b507QuantityViewArray[0].(b507QuantityViewJSON.Number)))
			if b507QuantityViewSelected == "capacity" {
				b507QuantityViewField["candidate_domain"] = b507QuantityViewFirst <= 250
				b507QuantityViewField["number"] = float64(b507QuantityViewFirst) / 10
				b507QuantityViewField["calculation"] = "Proven"
				if b507QuantityViewFirst > 250 {
					b507QuantityViewField["state"] = "outside-candidate-domain"
					b507QuantityViewOutput["state"] = "partial"
				}
			} else {
				b507QuantityViewSecond, _ := b507QuantityViewStrconv.Atoi(string(b507QuantityViewArray[1].(b507QuantityViewJSON.Number)))
				switch b507QuantityViewFirst + 256*b507QuantityViewSecond {
				case 0:
					b507QuantityViewField["label"] = "appliance-unavailable candidate"
				case 65535:
					b507QuantityViewField["label"] = "maximum-efficiency-preference"
				default:
					b507QuantityViewField["state"] = "unmapped"
					b507QuantityViewOutput["state"] = "partial"
				}
			}
		}
	}
	return b507QuantityViewJSON.Marshal(b507QuantityViewOutput)
}
