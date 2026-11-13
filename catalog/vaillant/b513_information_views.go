package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func b513InformationViewRead(b513InformationViewDecoder *json.Decoder) (any, error) {
	b513InformationViewToken, b513InformationViewError := b513InformationViewDecoder.Token()
	if b513InformationViewError != nil {
		return nil, b513InformationViewError
	}
	switch b513InformationViewToken {
	case json.Delim('{'):
		b513InformationViewObject := map[string]any{}
		for b513InformationViewDecoder.More() {
			b513InformationViewKeyToken, b513InformationViewError := b513InformationViewDecoder.Token()
			if b513InformationViewError != nil {
				return nil, b513InformationViewError
			}
			b513InformationViewKey := b513InformationViewKeyToken.(string)
			if _, b513InformationViewExists := b513InformationViewObject[b513InformationViewKey]; b513InformationViewExists {
				return nil, fmt.Errorf("duplicate key %q", b513InformationViewKey)
			}
			b513InformationViewValue, b513InformationViewError := b513InformationViewRead(b513InformationViewDecoder)
			if b513InformationViewError != nil {
				return nil, b513InformationViewError
			}
			b513InformationViewObject[b513InformationViewKey] = b513InformationViewValue
		}
		_, b513InformationViewError = b513InformationViewDecoder.Token()
		return b513InformationViewObject, b513InformationViewError
	case json.Delim('['):
		b513InformationViewArray := []any{}
		for b513InformationViewDecoder.More() {
			b513InformationViewValue, b513InformationViewError := b513InformationViewRead(b513InformationViewDecoder)
			if b513InformationViewError != nil {
				return nil, b513InformationViewError
			}
			b513InformationViewArray = append(b513InformationViewArray, b513InformationViewValue)
		}
		_, b513InformationViewError = b513InformationViewDecoder.Token()
		return b513InformationViewArray, b513InformationViewError
	}
	return b513InformationViewToken, nil
}

func b513InformationViewInteger(b513InformationViewValue any) (int64, bool) {
	b513InformationViewNumber, b513InformationViewOK := b513InformationViewValue.(json.Number)
	if !b513InformationViewOK {
		return 0, false
	}
	b513InformationViewInt, b513InformationViewError := strconv.ParseInt(string(b513InformationViewNumber), 10, 64)
	return b513InformationViewInt, b513InformationViewError == nil
}

func B513InformationViewInterpret(b513InformationViewInput json.RawMessage) (json.RawMessage, error) {
	b513InformationViewDecoder := json.NewDecoder(bytes.NewReader(b513InformationViewInput))
	b513InformationViewDecoder.UseNumber()
	b513InformationViewValue, b513InformationViewError := b513InformationViewRead(b513InformationViewDecoder)
	if b513InformationViewError != nil {
		return nil, b513InformationViewError
	}
	if _, b513InformationViewError = b513InformationViewDecoder.Token(); b513InformationViewError != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b513InformationViewRoot, b513InformationViewOK := b513InformationViewValue.(map[string]any)
	if !b513InformationViewOK {
		return nil, fmt.Errorf("input must be object")
	}
	for b513InformationViewKey := range b513InformationViewRoot {
		if b513InformationViewKey != "context" && b513InformationViewKey != "fields" && b513InformationViewKey != "arguments" {
			return nil, fmt.Errorf("unknown input key %q", b513InformationViewKey)
		}
	}
	b513InformationViewContext, b513InformationViewOK := b513InformationViewRoot["context"].(map[string]any)
	if !b513InformationViewOK {
		return nil, fmt.Errorf("context must be object")
	}
	b513InformationViewRaw, b513InformationViewOK := b513InformationViewRoot["fields"].(map[string]any)
	if !b513InformationViewOK {
		return nil, fmt.Errorf("fields must be object")
	}
	for _, b513InformationViewKey := range []string{"primary", "secondary", "operation"} {
		if b513InformationViewValue, b513InformationViewExists := b513InformationViewContext[b513InformationViewKey]; b513InformationViewExists {
			if _, b513InformationViewOK := b513InformationViewInteger(b513InformationViewValue); !b513InformationViewOK {
				return nil, fmt.Errorf("%s must be integer", b513InformationViewKey)
			}
		}
	}
	for _, b513InformationViewKey := range []string{"direction", "view"} {
		if b513InformationViewValue, b513InformationViewExists := b513InformationViewContext[b513InformationViewKey]; b513InformationViewExists {
			if _, b513InformationViewOK := b513InformationViewValue.(string); !b513InformationViewOK {
				return nil, fmt.Errorf("%s must be string", b513InformationViewKey)
			}
		}
	}
	for b513InformationViewKey, b513InformationViewValue := range b513InformationViewRaw {
		b513InformationViewArray, b513InformationViewOK := b513InformationViewValue.([]any)
		if !b513InformationViewOK {
			return nil, fmt.Errorf("%s must be byte array", b513InformationViewKey)
		}
		for _, b513InformationViewByte := range b513InformationViewArray {
			b513InformationViewInt, b513InformationViewOK := b513InformationViewInteger(b513InformationViewByte)
			if !b513InformationViewOK || b513InformationViewInt < 0 || b513InformationViewInt > 255 {
				return nil, fmt.Errorf("invalid byte in %s", b513InformationViewKey)
			}
		}
	}
	b513InformationViewFields := map[string]any{}
	b513InformationViewOutput := map[string]any{"context": b513InformationViewContext, "raw": b513InformationViewRaw, "state": "Unknown", "fields": b513InformationViewFields, "physical": "Unknown", "qualified": false}
	if b513InformationViewArguments, b513InformationViewExists := b513InformationViewRoot["arguments"]; b513InformationViewExists {
		b513InformationViewObject, b513InformationViewOK := b513InformationViewArguments.(map[string]any)
		if !b513InformationViewOK || len(b513InformationViewObject) != 0 {
			return nil, fmt.Errorf("no arguments defined for these views")
		}
		b513InformationViewOutput["arguments"] = b513InformationViewObject
	}
	b513InformationViewPrimary, _ := b513InformationViewInteger(b513InformationViewContext["primary"])
	b513InformationViewSecondary, _ := b513InformationViewInteger(b513InformationViewContext["secondary"])
	b513InformationViewOperation, _ := b513InformationViewInteger(b513InformationViewContext["operation"])
	if b513InformationViewPrimary != 181 || b513InformationViewSecondary != 19 || b513InformationViewContext["direction"] != "positive" {
		return json.Marshal(b513InformationViewOutput)
	}
	if b513InformationViewOperation == 7 && b513InformationViewContext["view"] == "information" {
		b513InformationViewObject, _ := b513InformationViewRaw["logical_object"].([]any)
		b513InformationViewOutput["state"] = "Hypothesis"
		if len(b513InformationViewObject) < 12 {
			b513InformationViewOutput["state"] = "partial"
		}
		for b513InformationViewIndex, b513InformationViewName := range []string{"type", "capabilities", "fuel", "storage", "features", "pump_mode"} {
			b513InformationViewPosition := b513InformationViewIndex + 3
			b513InformationViewField := map[string]any{"raw": []any{}, "state": "Unknown"}
			b513InformationViewFields[b513InformationViewName] = b513InformationViewField
			if len(b513InformationViewObject) <= b513InformationViewPosition {
				continue
			}
			b513InformationViewField["raw"] = b513InformationViewObject[b513InformationViewPosition : b513InformationViewPosition+1]
			b513InformationViewNumber, _ := b513InformationViewInteger(b513InformationViewObject[b513InformationViewPosition])
			if b513InformationViewIndex == 0 || b513InformationViewIndex == 5 {
				b513InformationViewLabels := []string{"modulating", "single-step", "two-step"}
				if b513InformationViewIndex == 5 {
					b513InformationViewLabels = []string{"pump-mode0", "pump-mode1"}
				}
				if b513InformationViewNumber >= int64(len(b513InformationViewLabels)) {
					b513InformationViewOutput["state"] = "partial"
					continue
				}
				b513InformationViewField["label"] = b513InformationViewLabels[b513InformationViewNumber]
			} else {
				b513InformationViewLabels := [][]string{{"floor-standing", "condensing", "fan-assisted", "open-vented", "high-modulation-range", "water-heater", "step-modulation", "ECS"}, {"gas", "oil", "ABC"}, {"layer-tank", "warm-start", "external-tank", "combination"}, {"solar-compact", "basic-solar", "AAPS", "advanced-power", "partial-load", "parameter-primary"}}[b513InformationViewIndex-1]
				b513InformationViewFlags := map[string]bool{}
				for b513InformationViewBit, b513InformationViewLabel := range b513InformationViewLabels {
					b513InformationViewFlags[b513InformationViewLabel] = b513InformationViewNumber&(1<<b513InformationViewBit) != 0
				}
				b513InformationViewField["flags"] = b513InformationViewFlags
				b513InformationViewField["reserved"] = b513InformationViewNumber & (255 ^ ((1 << len(b513InformationViewLabels)) - 1))
			}
			b513InformationViewField["number"] = b513InformationViewNumber
			b513InformationViewField["state"] = "Hypothesis"
		}
		b513InformationViewOpaque := []any{}
		if len(b513InformationViewObject) > 9 {
			b513InformationViewEnd := len(b513InformationViewObject)
			if b513InformationViewEnd > 12 {
				b513InformationViewEnd = 12
			}
			b513InformationViewOpaque = b513InformationViewObject[9:b513InformationViewEnd]
		}
		b513InformationViewFields["opaque"] = map[string]any{"raw": b513InformationViewOpaque, "state": "Unknown"}
	} else if (b513InformationViewOperation == 5 || b513InformationViewOperation == 6) && b513InformationViewContext["view"] == "feedback" {
		b513InformationViewFeedback, _ := b513InformationViewRaw["feedback"].([]any)
		b513InformationViewField := map[string]any{"raw": b513InformationViewFeedback, "state": "Unknown"}
		if b513InformationViewFeedback == nil {
			b513InformationViewField["raw"] = []any{}
		}
		b513InformationViewFields["feedback"] = b513InformationViewField
		b513InformationViewOutput["state"] = "partial"
		if len(b513InformationViewFeedback) == 1 {
			b513InformationViewNumber, _ := b513InformationViewInteger(b513InformationViewFeedback[0])
			b513InformationViewLabel := "rejected"
			if b513InformationViewOperation == 5 {
				b513InformationViewLabel = "not-accepted-or-not-supported"
			}
			if b513InformationViewNumber == 1 {
				b513InformationViewLabel = "accepted"
			}
			b513InformationViewField["number"] = b513InformationViewNumber
			b513InformationViewField["label"] = b513InformationViewLabel
			b513InformationViewField["state"] = "Hypothesis"
			b513InformationViewOutput["state"] = "Hypothesis"
		}
	}
	return json.Marshal(b513InformationViewOutput)
}
