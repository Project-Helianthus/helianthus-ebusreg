package vaillant

import (
	b521FieldViewBytes "bytes"
	b521FieldViewJSON "encoding/json"
	b521FieldViewFmt "fmt"
	b521FieldViewIO "io"
	b521FieldViewMath "math"
	b521FieldViewStrconv "strconv"
)

func b521FieldViewDecode(b521FieldViewDecoder *b521FieldViewJSON.Decoder) (any, error) {
	b521FieldViewToken, b521FieldViewError := b521FieldViewDecoder.Token()
	if b521FieldViewError != nil {
		return nil, b521FieldViewError
	}
	if b521FieldViewNumber, b521FieldViewOK := b521FieldViewToken.(b521FieldViewJSON.Number); b521FieldViewOK {
		b521FieldViewFloat, b521FieldViewError := b521FieldViewNumber.Float64()
		if b521FieldViewError != nil || b521FieldViewMath.IsInf(b521FieldViewFloat, 0) || b521FieldViewMath.IsNaN(b521FieldViewFloat) {
			return nil, b521FieldViewFmt.Errorf("nonfinite number")
		}
	}
	if b521FieldViewDelimiter, b521FieldViewOK := b521FieldViewToken.(b521FieldViewJSON.Delim); b521FieldViewOK {
		switch b521FieldViewDelimiter {
		case '{':
			b521FieldViewObject := map[string]any{}
			for b521FieldViewDecoder.More() {
				b521FieldViewKeyToken, b521FieldViewError := b521FieldViewDecoder.Token()
				if b521FieldViewError != nil {
					return nil, b521FieldViewError
				}
				b521FieldViewKey, b521FieldViewOK := b521FieldViewKeyToken.(string)
				if !b521FieldViewOK {
					return nil, b521FieldViewFmt.Errorf("object key required")
				}
				if _, b521FieldViewExists := b521FieldViewObject[b521FieldViewKey]; b521FieldViewExists {
					return nil, b521FieldViewFmt.Errorf("duplicate key %q", b521FieldViewKey)
				}
				b521FieldViewValue, b521FieldViewError := b521FieldViewDecode(b521FieldViewDecoder)
				if b521FieldViewError != nil {
					return nil, b521FieldViewError
				}
				b521FieldViewObject[b521FieldViewKey] = b521FieldViewValue
			}
			_, b521FieldViewError := b521FieldViewDecoder.Token()
			return b521FieldViewObject, b521FieldViewError
		case '[':
			b521FieldViewArray := []any{}
			for b521FieldViewDecoder.More() {
				b521FieldViewValue, b521FieldViewError := b521FieldViewDecode(b521FieldViewDecoder)
				if b521FieldViewError != nil {
					return nil, b521FieldViewError
				}
				b521FieldViewArray = append(b521FieldViewArray, b521FieldViewValue)
			}
			_, b521FieldViewError := b521FieldViewDecoder.Token()
			return b521FieldViewArray, b521FieldViewError
		default:
			return nil, b521FieldViewFmt.Errorf("unexpected delimiter")
		}
	}
	return b521FieldViewToken, nil
}

func b521FieldViewByte(b521FieldViewValue any) (uint64, error) {
	b521FieldViewNumber, b521FieldViewOK := b521FieldViewValue.(b521FieldViewJSON.Number)
	if !b521FieldViewOK {
		return 0, b521FieldViewFmt.Errorf("byte integer required")
	}
	return b521FieldViewStrconv.ParseUint(string(b521FieldViewNumber), 10, 8)
}

func B521FieldViewEvaluate(b521FieldViewInput b521FieldViewJSON.RawMessage) (b521FieldViewJSON.RawMessage, error) {
	b521FieldViewDecoder := b521FieldViewJSON.NewDecoder(b521FieldViewBytes.NewReader(b521FieldViewInput))
	b521FieldViewDecoder.UseNumber()
	b521FieldViewValue, b521FieldViewError := b521FieldViewDecode(b521FieldViewDecoder)
	if b521FieldViewError != nil {
		return nil, b521FieldViewError
	}
	if _, b521FieldViewError = b521FieldViewDecoder.Token(); b521FieldViewError != b521FieldViewIO.EOF {
		return nil, b521FieldViewFmt.Errorf("trailing JSON")
	}
	b521FieldViewOutput, b521FieldViewOK := b521FieldViewValue.(map[string]any)
	if !b521FieldViewOK {
		return nil, b521FieldViewFmt.Errorf("input object required")
	}
	b521FieldViewContext, b521FieldViewOK := b521FieldViewOutput["context"].(map[string]any)
	if !b521FieldViewOK {
		return nil, b521FieldViewFmt.Errorf("context object required")
	}
	b521FieldViewFields, b521FieldViewOK := b521FieldViewOutput["fields"].(map[string]any)
	if !b521FieldViewOK {
		return nil, b521FieldViewFmt.Errorf("fields object required")
	}
	for b521FieldViewName, b521FieldViewValue := range b521FieldViewFields {
		b521FieldViewArray, b521FieldViewOK := b521FieldViewValue.([]any)
		if !b521FieldViewOK {
			return nil, b521FieldViewFmt.Errorf("field %s requires byte array", b521FieldViewName)
		}
		for _, b521FieldViewValue := range b521FieldViewArray {
			if _, b521FieldViewError := b521FieldViewByte(b521FieldViewValue); b521FieldViewError != nil {
				return nil, b521FieldViewError
			}
		}
	}
	b521FieldViewArguments := map[string]any{}
	if b521FieldViewValue, b521FieldViewPresent := b521FieldViewOutput["arguments"]; b521FieldViewPresent {
		b521FieldViewArguments, b521FieldViewOK = b521FieldViewValue.(map[string]any)
		if !b521FieldViewOK {
			return nil, b521FieldViewFmt.Errorf("arguments object required")
		}
	}
	b521FieldViewSupported := true
	for b521FieldViewName, b521FieldViewExpected := range map[string]uint64{"primary": 181, "secondary": 33} {
		b521FieldViewValue, b521FieldViewPresent := b521FieldViewContext[b521FieldViewName]
		if !b521FieldViewPresent || b521FieldViewValue == nil {
			b521FieldViewSupported = false
			continue
		}
		b521FieldViewNumber, b521FieldViewError := b521FieldViewByte(b521FieldViewValue)
		if b521FieldViewError != nil {
			return nil, b521FieldViewError
		}
		if b521FieldViewNumber != b521FieldViewExpected {
			b521FieldViewSupported = false
		}
	}
	if b521FieldViewValue, b521FieldViewPresent := b521FieldViewContext["operation"]; b521FieldViewPresent {
		b521FieldViewSupported = false
		if b521FieldViewValue != nil {
			if _, b521FieldViewError := b521FieldViewByte(b521FieldViewValue); b521FieldViewError != nil {
				return nil, b521FieldViewError
			}
		}
	}
	for _, b521FieldViewName := range []string{"direction", "view"} {
		b521FieldViewValue, b521FieldViewPresent := b521FieldViewContext[b521FieldViewName]
		if !b521FieldViewPresent || b521FieldViewValue == nil {
			b521FieldViewSupported = false
			continue
		}
		b521FieldViewString, b521FieldViewOK := b521FieldViewValue.(string)
		if !b521FieldViewOK {
			return nil, b521FieldViewFmt.Errorf("context %s requires string", b521FieldViewName)
		}
		if (b521FieldViewName == "view" && b521FieldViewString != "module-fields") || (b521FieldViewName == "direction" && b521FieldViewString != "request" && b521FieldViewString != "positive") {
			b521FieldViewSupported = false
		}
	}
	b521FieldViewResults := map[string]any{}
	for b521FieldViewName, b521FieldViewRaw := range b521FieldViewFields {
		b521FieldViewResults[b521FieldViewName] = map[string]any{"raw": b521FieldViewRaw, "state": "Unknown", "unmodeled": true}
	}
	b521FieldViewOutput["results"] = b521FieldViewResults
	b521FieldViewOutput["state"] = "Unknown"
	b521FieldViewOutput["physical"] = "Unknown"
	b521FieldViewOutput["qualified"] = false
	if !b521FieldViewSupported {
		return b521FieldViewJSON.Marshal(b521FieldViewOutput)
	}
	if len(b521FieldViewArguments) != 0 {
		return nil, b521FieldViewFmt.Errorf("module-fields declares no arguments")
	}
	b521FieldViewWidths := map[string]int{"flags": 1, "appliance_bits": 1, "noise_percent": 1, "frost_temperature": 1}
	b521FieldViewLabels := map[string][]string{"flags": {"bit0", "bit1", "bit2", "bit3", "bit4", "bit5"}, "appliance_bits": {"bit0", "bit1", "bit2", "bit3", "bit4"}}
	if b521FieldViewContext["direction"] == "positive" {
		b521FieldViewWidths = map[string]int{"status0": 1, "status1": 1, "brine_percent": 1, "fan_percent": 1, "fan_speed": 1, "temperature0": 2, "temperature1": 2}
		b521FieldViewLabels = map[string][]string{"status0": {"ready", "deicing-request", "deicing-active", "deicing-timeout-error", "deicing-with-heater", "", "", "cleaning-request"}, "status1": {"heater", "relay-output", "fan-running", "brine-flow-reverse", "fan-error", "safety-limiter-error", "hot-summer", "bad-weather"}}
	}
	b521FieldViewOutput["state"] = "complete"
	for b521FieldViewName, b521FieldViewWidth := range b521FieldViewWidths {
		b521FieldViewRaw, b521FieldViewPresent := b521FieldViewFields[b521FieldViewName]
		b521FieldViewResult := map[string]any{"raw": b521FieldViewRaw, "present": b521FieldViewPresent, "state": "incomplete", "grade": "Hypothesis"}
		b521FieldViewResults[b521FieldViewName] = b521FieldViewResult
		b521FieldViewArray, _ := b521FieldViewRaw.([]any)
		if len(b521FieldViewArray) != b521FieldViewWidth {
			b521FieldViewOutput["state"] = "partial"
			continue
		}
		b521FieldViewResult["state"] = "complete"
		b521FieldViewResult["calculation"] = "Proven"
		b521FieldViewNumber, _ := b521FieldViewByte(b521FieldViewArray[0])
		if b521FieldViewNames, b521FieldViewIsFlags := b521FieldViewLabels[b521FieldViewName]; b521FieldViewIsFlags {
			b521FieldViewFlags := map[string]bool{}
			b521FieldViewReserved := b521FieldViewNumber
			for b521FieldViewBit, b521FieldViewLabel := range b521FieldViewNames {
				if b521FieldViewLabel == "" {
					continue
				}
				b521FieldViewFlags[b521FieldViewLabel] = b521FieldViewNumber&(1<<b521FieldViewBit) != 0
				b521FieldViewReserved &^= 1 << b521FieldViewBit
			}
			b521FieldViewResult["flags"] = b521FieldViewFlags
			b521FieldViewResult["reserved"] = b521FieldViewReserved
			continue
		}
		b521FieldViewScalar := float64(b521FieldViewNumber)
		switch b521FieldViewName {
		case "fan_speed":
			b521FieldViewScalar *= 10
			b521FieldViewResult["unit"] = "rpm"
		case "temperature0", "temperature1":
			b521FieldViewHigh, _ := b521FieldViewByte(b521FieldViewArray[1])
			b521FieldViewScalar = float64(int16(uint16(b521FieldViewNumber)|uint16(b521FieldViewHigh)<<8)) / 16
			b521FieldViewResult["unit"] = "Celsius"
		case "frost_temperature":
			b521FieldViewResult["unit"] = "Celsius"
		default:
			b521FieldViewResult["unit"] = "percent"
		}
		b521FieldViewResult["number"] = b521FieldViewScalar
	}
	return b521FieldViewJSON.Marshal(b521FieldViewOutput)
}
