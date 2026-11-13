package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

func b520ActuatorViewValue(b520ActuatorViewDecoder *json.Decoder) (any, error) {
	b520ActuatorViewToken, b520ActuatorViewErr := b520ActuatorViewDecoder.Token()
	if b520ActuatorViewErr != nil {
		return nil, b520ActuatorViewErr
	}
	if b520ActuatorViewDelimiter, b520ActuatorViewOK := b520ActuatorViewToken.(json.Delim); b520ActuatorViewOK {
		switch b520ActuatorViewDelimiter {
		case '{':
			b520ActuatorViewObject := map[string]any{}
			for b520ActuatorViewDecoder.More() {
				b520ActuatorViewKeyToken, b520ActuatorViewErr := b520ActuatorViewDecoder.Token()
				if b520ActuatorViewErr != nil {
					return nil, b520ActuatorViewErr
				}
				b520ActuatorViewKey, b520ActuatorViewOK := b520ActuatorViewKeyToken.(string)
				if !b520ActuatorViewOK {
					return nil, fmt.Errorf("object key must be a string")
				}
				if _, b520ActuatorViewExists := b520ActuatorViewObject[b520ActuatorViewKey]; b520ActuatorViewExists {
					return nil, fmt.Errorf("duplicate object key %q", b520ActuatorViewKey)
				}
				b520ActuatorViewChild, b520ActuatorViewErr := b520ActuatorViewValue(b520ActuatorViewDecoder)
				if b520ActuatorViewErr != nil {
					return nil, b520ActuatorViewErr
				}
				b520ActuatorViewObject[b520ActuatorViewKey] = b520ActuatorViewChild
			}
			b520ActuatorViewEnd, b520ActuatorViewErr := b520ActuatorViewDecoder.Token()
			if b520ActuatorViewErr != nil || b520ActuatorViewEnd != json.Delim('}') {
				return nil, fmt.Errorf("unterminated object")
			}
			return b520ActuatorViewObject, nil
		case '[':
			b520ActuatorViewArray := []any{}
			for b520ActuatorViewDecoder.More() {
				b520ActuatorViewChild, b520ActuatorViewErr := b520ActuatorViewValue(b520ActuatorViewDecoder)
				if b520ActuatorViewErr != nil {
					return nil, b520ActuatorViewErr
				}
				b520ActuatorViewArray = append(b520ActuatorViewArray, b520ActuatorViewChild)
			}
			b520ActuatorViewEnd, b520ActuatorViewErr := b520ActuatorViewDecoder.Token()
			if b520ActuatorViewErr != nil || b520ActuatorViewEnd != json.Delim(']') {
				return nil, fmt.Errorf("unterminated array")
			}
			return b520ActuatorViewArray, nil
		default:
			return nil, fmt.Errorf("unexpected delimiter")
		}
	}
	return b520ActuatorViewToken, nil
}

func b520ActuatorViewInteger(b520ActuatorViewValue any) (int64, bool) {
	b520ActuatorViewNumber, b520ActuatorViewOK := b520ActuatorViewValue.(json.Number)
	if !b520ActuatorViewOK {
		return 0, false
	}
	b520ActuatorViewInt, b520ActuatorViewErr := strconv.ParseInt(string(b520ActuatorViewNumber), 10, 64)
	return b520ActuatorViewInt, b520ActuatorViewErr == nil
}

// Check escapes before decoding, which otherwise replaces unpaired surrogates.
// The JSON decoder remains responsible for the rest of the syntax validation.
func b520ActuatorViewValidSurrogates(input []byte) bool {
	for i := 0; i < len(input); i++ {
		if input[i] != '\\' {
			continue
		}
		i++
		if i >= len(input) || input[i] != 'u' {
			continue // Skip escaped backslashes and other non-Unicode escapes.
		}
		if i+4 >= len(input) {
			return false
		}
		code, err := strconv.ParseUint(string(input[i+1:i+5]), 16, 16)
		if err != nil {
			return false
		}
		i += 4
		if code >= 0xdc00 && code <= 0xdfff {
			return false
		}
		if code >= 0xd800 && code <= 0xdbff {
			if i+6 >= len(input) || input[i+1] != '\\' || input[i+2] != 'u' {
				return false
			}
			low, err := strconv.ParseUint(string(input[i+3:i+7]), 16, 16)
			if err != nil || low < 0xdc00 || low > 0xdfff {
				return false
			}
			i += 6
		}
	}
	return true
}

func B520ActuatorViewInterpret(b520ActuatorViewInput json.RawMessage) (json.RawMessage, error) {
	if !utf8.Valid(b520ActuatorViewInput) {
		return nil, fmt.Errorf("input must be UTF-8 JSON")
	}
	if !b520ActuatorViewValidSurrogates(b520ActuatorViewInput) {
		return nil, fmt.Errorf("input must contain valid Unicode escapes with paired surrogates")
	}
	b520ActuatorViewDecoder := json.NewDecoder(bytes.NewReader(b520ActuatorViewInput))
	b520ActuatorViewDecoder.UseNumber()
	b520ActuatorViewValue, b520ActuatorViewErr := b520ActuatorViewValue(b520ActuatorViewDecoder)
	if b520ActuatorViewErr != nil {
		return nil, b520ActuatorViewErr
	}
	if _, b520ActuatorViewErr = b520ActuatorViewDecoder.Token(); b520ActuatorViewErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON data")
	}
	b520ActuatorViewInputObject, b520ActuatorViewOK := b520ActuatorViewValue.(map[string]any)
	if !b520ActuatorViewOK {
		return nil, fmt.Errorf("input must be an object")
	}
	for b520ActuatorViewKey := range b520ActuatorViewInputObject {
		if b520ActuatorViewKey != "context" && b520ActuatorViewKey != "fields" && b520ActuatorViewKey != "arguments" {
			return nil, fmt.Errorf("unknown input key %q", b520ActuatorViewKey)
		}
	}
	b520ActuatorViewContext, b520ActuatorViewOK := b520ActuatorViewInputObject["context"].(map[string]any)
	if !b520ActuatorViewOK {
		return nil, fmt.Errorf("context must be an object")
	}
	b520ActuatorViewRaw, b520ActuatorViewOK := b520ActuatorViewInputObject["fields"].(map[string]any)
	if !b520ActuatorViewOK {
		return nil, fmt.Errorf("fields must be an object")
	}
	for _, b520ActuatorViewKey := range []string{"primary", "secondary", "operation"} {
		if b520ActuatorViewSelector, b520ActuatorViewExists := b520ActuatorViewContext[b520ActuatorViewKey]; b520ActuatorViewExists {
			b520ActuatorViewNumber, b520ActuatorViewOK := b520ActuatorViewSelector.(json.Number)
			if !b520ActuatorViewOK || strings.ContainsAny(string(b520ActuatorViewNumber), ".eE") {
				return nil, fmt.Errorf("%s must be an integer", b520ActuatorViewKey)
			}
		}
	}
	for _, b520ActuatorViewKey := range []string{"direction", "view"} {
		if b520ActuatorViewSelector, b520ActuatorViewExists := b520ActuatorViewContext[b520ActuatorViewKey]; b520ActuatorViewExists {
			if _, b520ActuatorViewOK := b520ActuatorViewSelector.(string); !b520ActuatorViewOK {
				return nil, fmt.Errorf("%s must be a string", b520ActuatorViewKey)
			}
		}
	}
	for b520ActuatorViewKey, b520ActuatorViewField := range b520ActuatorViewRaw {
		b520ActuatorViewBytes, b520ActuatorViewOK := b520ActuatorViewField.([]any)
		if !b520ActuatorViewOK {
			return nil, fmt.Errorf("field %q must be a byte array", b520ActuatorViewKey)
		}
		for _, b520ActuatorViewByte := range b520ActuatorViewBytes {
			b520ActuatorViewInt, b520ActuatorViewOK := b520ActuatorViewInteger(b520ActuatorViewByte)
			if !b520ActuatorViewOK || b520ActuatorViewInt < 0 || b520ActuatorViewInt > 255 {
				return nil, fmt.Errorf("field %q has an invalid byte", b520ActuatorViewKey)
			}
		}
	}
	b520ActuatorViewFields := map[string]any{}
	b520ActuatorViewOutput := map[string]any{"context": b520ActuatorViewContext, "raw": b520ActuatorViewRaw, "state": "Unknown", "fields": b520ActuatorViewFields, "physical": "Unknown", "qualified": false}
	if b520ActuatorViewArguments, b520ActuatorViewExists := b520ActuatorViewInputObject["arguments"]; b520ActuatorViewExists {
		b520ActuatorViewObject, b520ActuatorViewOK := b520ActuatorViewArguments.(map[string]any)
		if !b520ActuatorViewOK || len(b520ActuatorViewObject) != 0 {
			return nil, fmt.Errorf("arguments must be an empty object; no scalar arguments are defined")
		}
		b520ActuatorViewOutput["arguments"] = b520ActuatorViewObject
	}
	b520ActuatorViewPrimary, _ := b520ActuatorViewInteger(b520ActuatorViewContext["primary"])
	b520ActuatorViewSecondary, _ := b520ActuatorViewInteger(b520ActuatorViewContext["secondary"])
	b520ActuatorViewOperation, b520ActuatorViewHasOperation := b520ActuatorViewInteger(b520ActuatorViewContext["operation"])
	b520ActuatorViewDirection, _ := b520ActuatorViewContext["direction"].(string)
	if b520ActuatorViewPrimary != 181 || b520ActuatorViewSecondary != 32 || !b520ActuatorViewHasOperation || b520ActuatorViewOperation != 0 || b520ActuatorViewContext["view"] != "logical-object" || (b520ActuatorViewDirection != "request" && b520ActuatorViewDirection != "positive") {
		return json.Marshal(b520ActuatorViewOutput)
	}
	b520ActuatorViewOutput["state"] = "Hypothesis"
	b520ActuatorViewObject, b520ActuatorViewPresent := b520ActuatorViewRaw["logical_object"].([]any)
	b520ActuatorViewNames := []string{"relay", "PWM1", "PWM2", "stepper", "LED", "ED", "temperature-table"}
	b520ActuatorViewOffsets := []int{3, 4, 5, 6, 7, 8, 9}
	b520ActuatorViewWidths := []int{1, 1, 1, 1, 1, 1, 1}
	if b520ActuatorViewDirection == "positive" {
		b520ActuatorViewNames = []string{"ADC", "temperature1", "temperature2", "temperature3", "temperature4", "reserved", "impeller", "vortex"}
		b520ActuatorViewOffsets = []int{3, 5, 7, 9, 11, 13, 14, 16}
		b520ActuatorViewWidths = []int{2, 2, 2, 2, 2, 1, 2, 2}
	}
	for b520ActuatorViewIndex, b520ActuatorViewName := range b520ActuatorViewNames {
		b520ActuatorViewOffset, b520ActuatorViewWidth := b520ActuatorViewOffsets[b520ActuatorViewIndex], b520ActuatorViewWidths[b520ActuatorViewIndex]
		b520ActuatorViewField := map[string]any{"state": "incomplete", "grade": "Hypothesis"}
		b520ActuatorViewFields[b520ActuatorViewName] = b520ActuatorViewField
		if !b520ActuatorViewPresent {
			b520ActuatorViewField["raw"] = nil
		} else {
			b520ActuatorViewStart := min(b520ActuatorViewOffset, len(b520ActuatorViewObject))
			b520ActuatorViewEnd := min(b520ActuatorViewOffset+b520ActuatorViewWidth, len(b520ActuatorViewObject))
			b520ActuatorViewField["raw"] = b520ActuatorViewObject[b520ActuatorViewStart:b520ActuatorViewEnd]
		}
		if len(b520ActuatorViewObject) < b520ActuatorViewOffset+b520ActuatorViewWidth {
			b520ActuatorViewOutput["state"] = "partial"
			continue
		}
		b520ActuatorViewField["state"] = "Hypothesis"
		b520ActuatorViewLow, _ := b520ActuatorViewInteger(b520ActuatorViewObject[b520ActuatorViewOffset])
		b520ActuatorViewCarrier := b520ActuatorViewLow
		if b520ActuatorViewWidth == 2 {
			b520ActuatorViewHigh, _ := b520ActuatorViewInteger(b520ActuatorViewObject[b520ActuatorViewOffset+1])
			b520ActuatorViewCarrier += b520ActuatorViewHigh * 256
		}
		switch b520ActuatorViewName {
		case "PWM1", "PWM2", "ED":
			b520ActuatorViewField["number"] = float64(b520ActuatorViewCarrier) / 2
			b520ActuatorViewField["calculation"] = "Proven"
		case "stepper":
			if b520ActuatorViewCarrier <= 240 {
				b520ActuatorViewField["number"] = b520ActuatorViewCarrier
			} else if b520ActuatorViewCarrier == 255 {
				b520ActuatorViewField["label"] = "calibrate/go-to-zero"
			} else {
				b520ActuatorViewField["state"] = "unmapped"
			}
		case "LED":
			if b520ActuatorViewCarrier <= 1 {
				b520ActuatorViewField["label"] = []string{"disabled", "enabled"}[b520ActuatorViewCarrier]
			} else {
				b520ActuatorViewField["state"] = "unmapped"
			}
		case "temperature-table":
			if b520ActuatorViewCarrier <= 3 {
				b520ActuatorViewField["label"] = []string{"standard3K3", "extended1K", "low-temperature15K", "standard1K"}[b520ActuatorViewCarrier]
				b520ActuatorViewField["flags"] = map[string]any{"candidate_celsius_range": [][2]float64{{-13.5, 116}, {-20, 155}, {-39, 80}, {-13.9, 130}}[b520ActuatorViewCarrier]}
			} else {
				b520ActuatorViewField["state"] = "unmapped"
			}
		case "temperature1", "temperature2", "temperature3", "temperature4":
			b520ActuatorViewField["number"] = float64(int16(b520ActuatorViewCarrier)) / 16
			b520ActuatorViewField["calculation"] = "Proven"
			b520ActuatorViewField["candidate_domain"] = b520ActuatorViewCarrier <= 3200
		case "ADC":
			b520ActuatorViewField["number"] = b520ActuatorViewCarrier
			b520ActuatorViewField["candidate_domain"] = b520ActuatorViewCarrier <= 1023
		case "impeller":
			b520ActuatorViewField["number"] = b520ActuatorViewCarrier
			b520ActuatorViewField["candidate_domain"] = b520ActuatorViewCarrier <= 1000
		case "vortex":
			b520ActuatorViewField["number"] = b520ActuatorViewCarrier
			b520ActuatorViewField["candidate_domain"] = b520ActuatorViewCarrier <= 350
			b520ActuatorViewField["flags"] = map[string]any{"candidate_operational_domain": b520ActuatorViewCarrier >= 5 && b520ActuatorViewCarrier <= 350, "candidate_hysteresis_hz": map[string]int{"on": 5, "off": 4}}
		case "reserved":
			b520ActuatorViewField["reserved"] = true
		}
		if b520ActuatorViewField["state"] == "unmapped" {
			b520ActuatorViewOutput["state"] = "partial"
		}
	}
	return json.Marshal(b520ActuatorViewOutput)
}
