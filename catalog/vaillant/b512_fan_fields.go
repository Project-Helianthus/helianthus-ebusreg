package vaillant

import (
	b512FanFieldBytesPackage "bytes"
	b512FanFieldJSON "encoding/json"
	b512FanFieldFmt "fmt"
	b512FanFieldIO "io"
	b512FanFieldStrconv "strconv"
)

func b512FanFieldRead(b512FanFieldDecoder *b512FanFieldJSON.Decoder) (any, error) {
	b512FanFieldToken, b512FanFieldErr := b512FanFieldDecoder.Token()
	if b512FanFieldErr != nil {
		return nil, b512FanFieldErr
	}
	switch b512FanFieldToken {
	case b512FanFieldJSON.Delim('{'):
		b512FanFieldObject := map[string]any{}
		for b512FanFieldDecoder.More() {
			b512FanFieldKeyToken, b512FanFieldErr := b512FanFieldDecoder.Token()
			if b512FanFieldErr != nil {
				return nil, b512FanFieldErr
			}
			b512FanFieldKey := b512FanFieldKeyToken.(string)
			if _, b512FanFieldExists := b512FanFieldObject[b512FanFieldKey]; b512FanFieldExists {
				return nil, b512FanFieldFmt.Errorf("duplicate key %q", b512FanFieldKey)
			}
			b512FanFieldValue, b512FanFieldErr := b512FanFieldRead(b512FanFieldDecoder)
			if b512FanFieldErr != nil {
				return nil, b512FanFieldErr
			}
			b512FanFieldObject[b512FanFieldKey] = b512FanFieldValue
		}
		_, b512FanFieldErr := b512FanFieldDecoder.Token()
		return b512FanFieldObject, b512FanFieldErr
	case b512FanFieldJSON.Delim('['):
		b512FanFieldArray := []any{}
		for b512FanFieldDecoder.More() {
			b512FanFieldValue, b512FanFieldErr := b512FanFieldRead(b512FanFieldDecoder)
			if b512FanFieldErr != nil {
				return nil, b512FanFieldErr
			}
			b512FanFieldArray = append(b512FanFieldArray, b512FanFieldValue)
		}
		_, b512FanFieldErr := b512FanFieldDecoder.Token()
		return b512FanFieldArray, b512FanFieldErr
	default:
		return b512FanFieldToken, nil
	}
}

func b512FanFieldUint(b512FanFieldValue any, b512FanFieldBits int) (uint64, bool) {
	b512FanFieldNumber, b512FanFieldOK := b512FanFieldValue.(b512FanFieldJSON.Number)
	if !b512FanFieldOK {
		return 0, false
	}
	b512FanFieldResult, b512FanFieldErr := b512FanFieldStrconv.ParseUint(string(b512FanFieldNumber), 10, b512FanFieldBits)
	return b512FanFieldResult, b512FanFieldErr == nil
}

func B512FanFieldEvaluate(b512FanFieldInput b512FanFieldJSON.RawMessage) (b512FanFieldJSON.RawMessage, error) {
	b512FanFieldDecoder := b512FanFieldJSON.NewDecoder(b512FanFieldBytesPackage.NewReader(b512FanFieldInput))
	b512FanFieldDecoder.UseNumber()
	b512FanFieldValue, b512FanFieldErr := b512FanFieldRead(b512FanFieldDecoder)
	if b512FanFieldErr != nil {
		return nil, b512FanFieldErr
	}
	if _, b512FanFieldErr = b512FanFieldDecoder.Token(); b512FanFieldErr != b512FanFieldIO.EOF {
		return nil, b512FanFieldFmt.Errorf("trailing JSON")
	}
	b512FanFieldInputObject, b512FanFieldOK := b512FanFieldValue.(map[string]any)
	if !b512FanFieldOK {
		return nil, b512FanFieldFmt.Errorf("input must be object")
	}
	b512FanFieldContext, b512FanFieldOK := b512FanFieldInputObject["context"].(map[string]any)
	if !b512FanFieldOK {
		return nil, b512FanFieldFmt.Errorf("context must be object")
	}
	b512FanFieldRaw, b512FanFieldOK := b512FanFieldInputObject["fields"].(map[string]any)
	if !b512FanFieldOK {
		return nil, b512FanFieldFmt.Errorf("fields must be object")
	}
	for _, b512FanFieldKey := range []string{"primary", "secondary", "operation"} {
		if b512FanFieldValue, b512FanFieldExists := b512FanFieldContext[b512FanFieldKey]; b512FanFieldExists && b512FanFieldValue != nil {
			if _, b512FanFieldOK := b512FanFieldUint(b512FanFieldValue, 8); !b512FanFieldOK {
				return nil, b512FanFieldFmt.Errorf("invalid %s", b512FanFieldKey)
			}
		}
	}
	for _, b512FanFieldKey := range []string{"view", "direction"} {
		if b512FanFieldValue, b512FanFieldExists := b512FanFieldContext[b512FanFieldKey]; b512FanFieldExists && b512FanFieldValue != nil {
			if _, b512FanFieldOK := b512FanFieldValue.(string); !b512FanFieldOK {
				return nil, b512FanFieldFmt.Errorf("invalid %s", b512FanFieldKey)
			}
		}
	}
	for b512FanFieldKey, b512FanFieldValue := range b512FanFieldRaw {
		b512FanFieldArray, b512FanFieldOK := b512FanFieldValue.([]any)
		if !b512FanFieldOK {
			return nil, b512FanFieldFmt.Errorf("field %s must be byte array", b512FanFieldKey)
		}
		for _, b512FanFieldByte := range b512FanFieldArray {
			if _, b512FanFieldOK := b512FanFieldUint(b512FanFieldByte, 8); !b512FanFieldOK {
				return nil, b512FanFieldFmt.Errorf("invalid byte in %s", b512FanFieldKey)
			}
		}
	}
	b512FanFieldFields := map[string]any{}
	b512FanFieldOutput := map[string]any{"context": b512FanFieldContext, "raw": b512FanFieldRaw, "fields": b512FanFieldFields, "state": "Unknown", "physical": "Unknown", "qualified": false}
	if b512FanFieldArguments, b512FanFieldExists := b512FanFieldInputObject["arguments"]; b512FanFieldExists {
		if _, b512FanFieldOK := b512FanFieldArguments.(map[string]any); !b512FanFieldOK {
			return nil, b512FanFieldFmt.Errorf("arguments must be object")
		}
		b512FanFieldOutput["arguments"] = b512FanFieldArguments
	}
	for b512FanFieldKey, b512FanFieldValue := range b512FanFieldRaw {
		b512FanFieldFields[b512FanFieldKey] = map[string]any{"raw": b512FanFieldValue, "state": "Unknown"}
	}
	b512FanFieldPrimary, _ := b512FanFieldUint(b512FanFieldContext["primary"], 8)
	b512FanFieldSecondary, _ := b512FanFieldUint(b512FanFieldContext["secondary"], 8)
	b512FanFieldOperation, _ := b512FanFieldUint(b512FanFieldContext["operation"], 8)
	b512FanFieldView, _ := b512FanFieldContext["view"].(string)
	b512FanFieldDirection, _ := b512FanFieldContext["direction"].(string)
	b512FanFieldViews := map[string]uint64{"fan": 11, "fan-control": 12, "digital-air": 13, "air-quality": 14}
	if b512FanFieldPrimary != 181 || b512FanFieldSecondary != 18 || b512FanFieldViews[b512FanFieldView] == 0 || b512FanFieldViews[b512FanFieldView] != b512FanFieldOperation || (b512FanFieldDirection != "request" && b512FanFieldDirection != "positive") {
		return b512FanFieldJSON.Marshal(b512FanFieldOutput)
	}
	b512FanFieldSpecs := map[string]string{}
	switch b512FanFieldView {
	case "fan":
		b512FanFieldSpecs = map[string]string{"fan1": "rpm", "fan2": "rpm"}
		if b512FanFieldDirection == "request" {
			b512FanFieldSpecs["relay"] = "relay"
		} else {
			b512FanFieldSpecs["status"] = "status"
			for _, b512FanFieldKey := range []string{"temperature1", "temperature2", "temperature3"} {
				b512FanFieldSpecs[b512FanFieldKey] = "Celsius"
			}
		}
	case "fan-control":
		if b512FanFieldDirection == "request" {
			b512FanFieldSpecs = map[string]string{"pulse": "u8", "p": "u8", "i": "u8", "d": "u8", "fan_count": "u8", "minimum_pwm": "PWM", "maximum_pwm": "PWM"}
		} else {
			b512FanFieldSpecs["feedback"] = "feedback"
		}
	case "digital-air":
		if b512FanFieldDirection == "positive" {
			b512FanFieldSpecs = map[string]string{"digital": "digital", "temperature1": "Celsius", "temperature2": "Celsius"}
		}
	case "air-quality":
		if b512FanFieldDirection == "positive" {
			b512FanFieldSpecs = map[string]string{"temperature": "Celsius", "humidity": "percent", "air_quality": "Unknown"}
		}
	}
	b512FanFieldOutput["state"] = "complete"
	for b512FanFieldKey, b512FanFieldKind := range b512FanFieldSpecs {
		b512FanFieldWidth := 1
		if b512FanFieldKind == "rpm" || b512FanFieldKind == "Celsius" || b512FanFieldKind == "PWM" || b512FanFieldKind == "percent" || b512FanFieldKind == "Unknown" {
			b512FanFieldWidth = 2
		}
		b512FanFieldField := map[string]any{"state": "incomplete"}
		b512FanFieldFields[b512FanFieldKey] = b512FanFieldField
		b512FanFieldValue, b512FanFieldExists := b512FanFieldRaw[b512FanFieldKey]
		if b512FanFieldExists {
			b512FanFieldField["raw"] = b512FanFieldValue
		}
		b512FanFieldArray, _ := b512FanFieldValue.([]any)
		if !b512FanFieldExists || len(b512FanFieldArray) != b512FanFieldWidth {
			b512FanFieldOutput["state"] = "partial"
			continue
		}
		b512FanFieldLow, _ := b512FanFieldUint(b512FanFieldArray[0], 8)
		b512FanFieldWord := b512FanFieldLow
		if b512FanFieldWidth == 2 {
			b512FanFieldHigh, _ := b512FanFieldUint(b512FanFieldArray[1], 8)
			b512FanFieldWord += b512FanFieldHigh * 256
		}
		b512FanFieldField["state"] = "Hypothesis"
		switch b512FanFieldKind {
		case "relay":
			b512FanFieldField["flags"] = map[string]any{"relay1": map[bool]string{false: "15V", true: "tray-heater-relay"}[b512FanFieldWord&1 != 0], "relay2": map[bool]string{false: "15V", true: "tray-heater-relay"}[b512FanFieldWord&2 != 0]}
			b512FanFieldField["reserved"] = b512FanFieldWord & 252
		case "status":
			if b512FanFieldLabel, b512FanFieldOK := map[uint64]string{0: "idle", 1: "running", 255: "fault"}[b512FanFieldWord]; b512FanFieldOK {
				b512FanFieldField["label"] = b512FanFieldLabel
			} else {
				b512FanFieldField["state"] = "Unknown"
				b512FanFieldField["unmodeled"] = b512FanFieldWord
			}
		case "feedback":
			if b512FanFieldWord == 1 {
				b512FanFieldField["label"] = "accepted"
			} else {
				b512FanFieldField["label"] = "rejected"
			}
		case "digital":
			if b512FanFieldWord <= 1 {
				b512FanFieldField["label"] = b512FanFieldWord == 1
			} else {
				b512FanFieldField["state"] = "Unknown"
				b512FanFieldField["unmodeled"] = b512FanFieldWord
			}
		default:
			b512FanFieldNumber := float64(b512FanFieldWord)
			switch b512FanFieldKind {
			case "Celsius":
				b512FanFieldNumber = float64(int16(b512FanFieldWord)) / 16
			case "PWM":
				b512FanFieldNumber = float64(int16(b512FanFieldWord)) / 256
			case "percent":
				b512FanFieldNumber /= 16
			}
			b512FanFieldField["number"] = b512FanFieldNumber
			b512FanFieldField["calculation"] = "Proven"
			if b512FanFieldKind != "u8" {
				b512FanFieldField["unit"] = b512FanFieldKind
			}
			if b512FanFieldKind == "PWM" {
				b512FanFieldField["unit"] = "percent"
				b512FanFieldField["candidate_domain"] = true
				b512FanFieldField["candidate_default"] = 1
			}
		}
	}
	return b512FanFieldJSON.Marshal(b512FanFieldOutput)
}
