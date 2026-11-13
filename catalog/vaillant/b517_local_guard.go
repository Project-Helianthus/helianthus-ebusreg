package vaillant

import (
	b517LocalGuardBytes "bytes"
	b517LocalGuardJSON "encoding/json"
	b517LocalGuardFmt "fmt"
	b517LocalGuardIO "io"
	b517LocalGuardBig "math/big"
	b517LocalGuardStrconv "strconv"
	b517LocalGuardUTF8 "unicode/utf8"
)

func b517LocalGuardRead(b517LocalGuardDecoder *b517LocalGuardJSON.Decoder) (any, error) {
	b517LocalGuardToken, b517LocalGuardErr := b517LocalGuardDecoder.Token()
	if b517LocalGuardErr != nil {
		return nil, b517LocalGuardErr
	}
	if b517LocalGuardDelimiter, b517LocalGuardOK := b517LocalGuardToken.(b517LocalGuardJSON.Delim); b517LocalGuardOK {
		switch b517LocalGuardDelimiter {
		case '{':
			b517LocalGuardObject := map[string]any{}
			for b517LocalGuardDecoder.More() {
				b517LocalGuardKey, b517LocalGuardErr := b517LocalGuardDecoder.Token()
				if b517LocalGuardErr != nil {
					return nil, b517LocalGuardErr
				}
				b517LocalGuardName, b517LocalGuardOK := b517LocalGuardKey.(string)
				if !b517LocalGuardOK {
					return nil, b517LocalGuardFmt.Errorf("expected object key")
				}
				if _, b517LocalGuardExists := b517LocalGuardObject[b517LocalGuardName]; b517LocalGuardExists {
					return nil, b517LocalGuardFmt.Errorf("duplicate key")
				}
				b517LocalGuardValue, b517LocalGuardErr := b517LocalGuardRead(b517LocalGuardDecoder)
				if b517LocalGuardErr != nil {
					return nil, b517LocalGuardErr
				}
				b517LocalGuardObject[b517LocalGuardName] = b517LocalGuardValue
			}
			_, b517LocalGuardErr = b517LocalGuardDecoder.Token()
			return b517LocalGuardObject, b517LocalGuardErr
		case '[':
			b517LocalGuardArray := []any{}
			for b517LocalGuardDecoder.More() {
				b517LocalGuardValue, b517LocalGuardErr := b517LocalGuardRead(b517LocalGuardDecoder)
				if b517LocalGuardErr != nil {
					return nil, b517LocalGuardErr
				}
				b517LocalGuardArray = append(b517LocalGuardArray, b517LocalGuardValue)
			}
			_, b517LocalGuardErr = b517LocalGuardDecoder.Token()
			return b517LocalGuardArray, b517LocalGuardErr
		default:
			return nil, b517LocalGuardFmt.Errorf("unexpected delimiter")
		}
	}
	if b517LocalGuardNumber, b517LocalGuardOK := b517LocalGuardToken.(b517LocalGuardJSON.Number); b517LocalGuardOK {
		if _, b517LocalGuardErr := b517LocalGuardStrconv.ParseFloat(string(b517LocalGuardNumber), 64); b517LocalGuardErr != nil {
			return nil, b517LocalGuardFmt.Errorf("nonfinite number")
		}
	}
	return b517LocalGuardToken, nil
}

func b517LocalGuardInteger(b517LocalGuardValue any, b517LocalGuardMax int64) (int64, bool) {
	b517LocalGuardNumber, b517LocalGuardOK := b517LocalGuardValue.(b517LocalGuardJSON.Number)
	if !b517LocalGuardOK {
		return 0, false
	}
	b517LocalGuardRational, b517LocalGuardOK := new(b517LocalGuardBig.Rat).SetString(string(b517LocalGuardNumber))
	if !b517LocalGuardOK || !b517LocalGuardRational.IsInt() || !b517LocalGuardRational.Num().IsInt64() {
		return 0, false
	}
	b517LocalGuardInteger := b517LocalGuardRational.Num().Int64()
	return b517LocalGuardInteger, b517LocalGuardInteger >= 0 && b517LocalGuardInteger <= b517LocalGuardMax
}

func b517LocalGuardEncoding(b517LocalGuardData []byte) bool {
	if !b517LocalGuardUTF8.Valid(b517LocalGuardData) {
		return false
	}
	b517LocalGuardInside := false
	b517LocalGuardUnit := func(b517LocalGuardAt int) (int, bool) {
		if b517LocalGuardAt+6 > len(b517LocalGuardData) || b517LocalGuardData[b517LocalGuardAt] != '\\' || b517LocalGuardData[b517LocalGuardAt+1] != 'u' {
			return 0, false
		}
		b517LocalGuardParsed, b517LocalGuardErr := b517LocalGuardStrconv.ParseUint(string(b517LocalGuardData[b517LocalGuardAt+2:b517LocalGuardAt+6]), 16, 32)
		return int(b517LocalGuardParsed), b517LocalGuardErr == nil
	}
	for b517LocalGuardIndex := 0; b517LocalGuardIndex < len(b517LocalGuardData); b517LocalGuardIndex++ {
		switch {
		case !b517LocalGuardInside:
			b517LocalGuardInside = b517LocalGuardData[b517LocalGuardIndex] == '"'
		case b517LocalGuardData[b517LocalGuardIndex] == '"':
			b517LocalGuardInside = false
		case b517LocalGuardData[b517LocalGuardIndex] == '\\':
			if b517LocalGuardIndex+1 < len(b517LocalGuardData) && b517LocalGuardData[b517LocalGuardIndex+1] == 'u' {
				b517LocalGuardCode, b517LocalGuardOK := b517LocalGuardUnit(b517LocalGuardIndex)
				if !b517LocalGuardOK {
					return false
				}
				if b517LocalGuardCode >= 0xDC00 && b517LocalGuardCode <= 0xDFFF {
					return false
				}
				if b517LocalGuardCode >= 0xD800 && b517LocalGuardCode <= 0xDBFF {
					b517LocalGuardLow, b517LocalGuardOK := b517LocalGuardUnit(b517LocalGuardIndex + 6)
					if !b517LocalGuardOK || b517LocalGuardLow < 0xDC00 || b517LocalGuardLow > 0xDFFF {
						return false
					}
					b517LocalGuardIndex += 6
				}
				b517LocalGuardIndex += 5
			} else {
				b517LocalGuardIndex++
			}
		}
	}
	return true
}

func B517LocalGuardEvaluate(b517LocalGuardInputJSON b517LocalGuardJSON.RawMessage) (b517LocalGuardJSON.RawMessage, error) {
	if !b517LocalGuardEncoding(b517LocalGuardInputJSON) {
		return nil, b517LocalGuardFmt.Errorf("invalid string encoding")
	}
	b517LocalGuardDecoder := b517LocalGuardJSON.NewDecoder(b517LocalGuardBytes.NewReader(b517LocalGuardInputJSON))
	b517LocalGuardDecoder.UseNumber()
	b517LocalGuardValue, b517LocalGuardErr := b517LocalGuardRead(b517LocalGuardDecoder)
	if b517LocalGuardErr != nil {
		return nil, b517LocalGuardErr
	}
	if _, b517LocalGuardErr = b517LocalGuardDecoder.Token(); b517LocalGuardErr != b517LocalGuardIO.EOF {
		return nil, b517LocalGuardFmt.Errorf("trailing JSON")
	}
	b517LocalGuardObject, b517LocalGuardOK := b517LocalGuardValue.(map[string]any)
	if !b517LocalGuardOK {
		return nil, b517LocalGuardFmt.Errorf("expected object")
	}
	b517LocalGuardContext, b517LocalGuardOK := b517LocalGuardObject["context"].(map[string]any)
	if !b517LocalGuardOK {
		return nil, b517LocalGuardFmt.Errorf("expected context object")
	}
	b517LocalGuardFields, b517LocalGuardOK := b517LocalGuardObject["fields"].(map[string]any)
	if !b517LocalGuardOK {
		return nil, b517LocalGuardFmt.Errorf("expected fields object")
	}
	b517LocalGuardArguments := map[string]any{}
	if b517LocalGuardValue, b517LocalGuardPresent := b517LocalGuardObject["arguments"]; b517LocalGuardPresent {
		b517LocalGuardArguments, b517LocalGuardOK = b517LocalGuardValue.(map[string]any)
		if !b517LocalGuardOK {
			return nil, b517LocalGuardFmt.Errorf("expected arguments object")
		}
	}
	for _, b517LocalGuardName := range []string{"primary", "secondary", "operation"} {
		if b517LocalGuardValue, b517LocalGuardPresent := b517LocalGuardContext[b517LocalGuardName]; b517LocalGuardPresent {
			if _, b517LocalGuardOK := b517LocalGuardInteger(b517LocalGuardValue, 255); !b517LocalGuardOK {
				return nil, b517LocalGuardFmt.Errorf("invalid context %s", b517LocalGuardName)
			}
		} else if b517LocalGuardName != "operation" {
			return nil, b517LocalGuardFmt.Errorf("missing context %s", b517LocalGuardName)
		}
	}
	for _, b517LocalGuardName := range []string{"direction", "view"} {
		if _, b517LocalGuardOK := b517LocalGuardContext[b517LocalGuardName].(string); !b517LocalGuardOK {
			return nil, b517LocalGuardFmt.Errorf("invalid context %s", b517LocalGuardName)
		}
	}
	for b517LocalGuardName, b517LocalGuardValue := range b517LocalGuardFields {
		b517LocalGuardArray, b517LocalGuardOK := b517LocalGuardValue.([]any)
		if !b517LocalGuardOK {
			return nil, b517LocalGuardFmt.Errorf("invalid field %s", b517LocalGuardName)
		}
		for _, b517LocalGuardByte := range b517LocalGuardArray {
			if _, b517LocalGuardOK := b517LocalGuardInteger(b517LocalGuardByte, 255); !b517LocalGuardOK {
				return nil, b517LocalGuardFmt.Errorf("invalid byte in %s", b517LocalGuardName)
			}
		}
	}
	b517LocalGuardNumbers := map[string]int64{}
	for _, b517LocalGuardName := range []string{"ready_state", "final_option", "received_length"} {
		if b517LocalGuardValue, b517LocalGuardPresent := b517LocalGuardArguments[b517LocalGuardName]; b517LocalGuardPresent {
			b517LocalGuardMax := int64(255)
			if b517LocalGuardName == "received_length" {
				b517LocalGuardMax = 65535
			}
			b517LocalGuardNumber, b517LocalGuardOK := b517LocalGuardInteger(b517LocalGuardValue, b517LocalGuardMax)
			if !b517LocalGuardOK {
				return nil, b517LocalGuardFmt.Errorf("invalid argument %s", b517LocalGuardName)
			}
			b517LocalGuardNumbers[b517LocalGuardName] = b517LocalGuardNumber
		}
	}
	b517LocalGuardPrimary, _ := b517LocalGuardInteger(b517LocalGuardContext["primary"], 255)
	b517LocalGuardSecondary, _ := b517LocalGuardInteger(b517LocalGuardContext["secondary"], 255)
	_, b517LocalGuardOperation := b517LocalGuardContext["operation"]
	b517LocalGuardRaw := map[string]any{}
	for b517LocalGuardName, b517LocalGuardValue := range b517LocalGuardObject {
		b517LocalGuardRaw[b517LocalGuardName] = b517LocalGuardValue
	}
	b517LocalGuardObject["raw"] = b517LocalGuardRaw
	b517LocalGuardResults := map[string]any{}
	b517LocalGuardObject["state"] = "Unknown"
	b517LocalGuardObject["physical"] = "Unknown"
	b517LocalGuardObject["qualified"] = false
	b517LocalGuardObject["results"] = b517LocalGuardResults
	if b517LocalGuardPrimary != 181 || b517LocalGuardSecondary != 23 || b517LocalGuardOperation || b517LocalGuardContext["direction"] != "local" || b517LocalGuardContext["view"] != "adaptation-guard" {
		return b517LocalGuardJSON.Marshal(b517LocalGuardObject)
	}
	for _, b517LocalGuardName := range []string{"ready_state", "final_option", "received_length"} {
		if _, b517LocalGuardPresent := b517LocalGuardNumbers[b517LocalGuardName]; !b517LocalGuardPresent {
			return nil, b517LocalGuardFmt.Errorf("missing argument %s", b517LocalGuardName)
		}
	}
	b517LocalGuardRequest, _ := b517LocalGuardFields["request"].([]any)
	b517LocalGuardInput, _ := b517LocalGuardFields["input"].([]any)
	_, b517LocalGuardInputPresent := b517LocalGuardFields["input"]
	b517LocalGuardRequestComplete := len(b517LocalGuardRequest) >= 2
	b517LocalGuardInputComplete := len(b517LocalGuardInput) >= 5
	b517LocalGuardObject["state"] = "complete"
	if !b517LocalGuardRequestComplete || !b517LocalGuardInputComplete {
		b517LocalGuardObject["state"] = "partial"
	}
	b517LocalGuardEligible := true
	b517LocalGuardAdd := func(b517LocalGuardName string, b517LocalGuardRaw any, b517LocalGuardComplete bool, b517LocalGuardFlag bool) {
		b517LocalGuardState := "complete"
		if !b517LocalGuardComplete {
			b517LocalGuardState = "incomplete"
		}
		b517LocalGuardResults[b517LocalGuardName] = map[string]any{"raw": b517LocalGuardRaw, "state": b517LocalGuardState, "grade": "Hypothesis", "flags": map[string]bool{"value": b517LocalGuardComplete && b517LocalGuardFlag}}
		if b517LocalGuardComplete {
			b517LocalGuardResults[b517LocalGuardName].(map[string]any)["calculation"] = "Proven"
		}
		b517LocalGuardEligible = b517LocalGuardEligible && b517LocalGuardComplete && b517LocalGuardFlag
	}
	b517LocalGuardAdd("ready_state", b517LocalGuardArguments["ready_state"], true, b517LocalGuardNumbers["ready_state"] == 3)
	b517LocalGuardAdd("final_option", b517LocalGuardArguments["final_option"], true, b517LocalGuardNumbers["final_option"] != 0)
	b517LocalGuardAdd("initialized_input", b517LocalGuardFields, b517LocalGuardInputComplete, b517LocalGuardInputComplete)
	b517LocalGuardPayloadLength := false
	if b517LocalGuardInputComplete {
		b517LocalGuardByte, _ := b517LocalGuardInteger(b517LocalGuardInput[4], 255)
		b517LocalGuardPayloadLength = b517LocalGuardByte+6 == b517LocalGuardNumbers["received_length"]
	}
	b517LocalGuardAdd("payload_length", map[string]any{"fields": b517LocalGuardFields, "received_length": b517LocalGuardArguments["received_length"]}, b517LocalGuardInputComplete, b517LocalGuardPayloadLength)
	b517LocalGuardAdd("initialized_length", map[string]any{"fields": b517LocalGuardFields, "received_length": b517LocalGuardArguments["received_length"]}, b517LocalGuardInputPresent, b517LocalGuardNumbers["received_length"] == int64(len(b517LocalGuardInput)))
	b517LocalGuardRequestIdentity := false
	if b517LocalGuardRequestComplete {
		b517LocalGuardFirst, _ := b517LocalGuardInteger(b517LocalGuardRequest[0], 255)
		b517LocalGuardSecond, _ := b517LocalGuardInteger(b517LocalGuardRequest[1], 255)
		b517LocalGuardRequestIdentity = b517LocalGuardFirst == 181 && b517LocalGuardSecond == 23
	}
	b517LocalGuardAdd("request_identity", b517LocalGuardFields, b517LocalGuardRequestComplete, b517LocalGuardRequestIdentity)
	b517LocalGuardInputIdentity := false
	if len(b517LocalGuardInput) >= 2 {
		b517LocalGuardFirst, _ := b517LocalGuardInteger(b517LocalGuardInput[0], 255)
		b517LocalGuardSecond, _ := b517LocalGuardInteger(b517LocalGuardInput[1], 255)
		b517LocalGuardInputIdentity = b517LocalGuardFirst == 181 && b517LocalGuardSecond == 24
	}
	b517LocalGuardAdd("input_identity", b517LocalGuardFields, len(b517LocalGuardInput) >= 2, b517LocalGuardInputIdentity)
	b517LocalGuardState := "complete"
	if b517LocalGuardObject["state"] == "partial" {
		b517LocalGuardState = "incomplete"
	}
	b517LocalGuardResults["eligible"] = map[string]any{"raw": map[string]any{"fields": b517LocalGuardFields, "arguments": b517LocalGuardArguments}, "state": b517LocalGuardState, "grade": "Hypothesis", "flags": map[string]bool{"value": b517LocalGuardEligible}}
	if b517LocalGuardState == "complete" {
		b517LocalGuardResults["eligible"].(map[string]any)["calculation"] = "Proven"
	}
	return b517LocalGuardJSON.Marshal(b517LocalGuardObject)
}
