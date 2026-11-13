package vaillant

import (
	b531OrdinaryLimitBytesPackage "bytes"
	b531OrdinaryLimitJSON "encoding/json"
	b531OrdinaryLimitFmt "fmt"
	b531OrdinaryLimitIO "io"
	b531OrdinaryLimitStrconv "strconv"
)

func b531OrdinaryLimitDecode(b531OrdinaryLimitDecoder *b531OrdinaryLimitJSON.Decoder) (any, error) {
	b531OrdinaryLimitToken, b531OrdinaryLimitError := b531OrdinaryLimitDecoder.Token()
	if b531OrdinaryLimitError != nil {
		return nil, b531OrdinaryLimitError
	}
	if b531OrdinaryLimitDelimiter, b531OrdinaryLimitOK := b531OrdinaryLimitToken.(b531OrdinaryLimitJSON.Delim); b531OrdinaryLimitOK {
		switch b531OrdinaryLimitDelimiter {
		case '{':
			b531OrdinaryLimitObject := map[string]any{}
			for b531OrdinaryLimitDecoder.More() {
				b531OrdinaryLimitKeyToken, b531OrdinaryLimitError := b531OrdinaryLimitDecoder.Token()
				if b531OrdinaryLimitError != nil {
					return nil, b531OrdinaryLimitError
				}
				b531OrdinaryLimitKey, b531OrdinaryLimitOK := b531OrdinaryLimitKeyToken.(string)
				if !b531OrdinaryLimitOK {
					return nil, b531OrdinaryLimitFmt.Errorf("invalid object key")
				}
				if _, b531OrdinaryLimitExists := b531OrdinaryLimitObject[b531OrdinaryLimitKey]; b531OrdinaryLimitExists {
					return nil, b531OrdinaryLimitFmt.Errorf("duplicate key %q", b531OrdinaryLimitKey)
				}
				b531OrdinaryLimitValue, b531OrdinaryLimitError := b531OrdinaryLimitDecode(b531OrdinaryLimitDecoder)
				if b531OrdinaryLimitError != nil {
					return nil, b531OrdinaryLimitError
				}
				b531OrdinaryLimitObject[b531OrdinaryLimitKey] = b531OrdinaryLimitValue
			}
			_, b531OrdinaryLimitError := b531OrdinaryLimitDecoder.Token()
			return b531OrdinaryLimitObject, b531OrdinaryLimitError
		case '[':
			b531OrdinaryLimitArray := []any{}
			for b531OrdinaryLimitDecoder.More() {
				b531OrdinaryLimitValue, b531OrdinaryLimitError := b531OrdinaryLimitDecode(b531OrdinaryLimitDecoder)
				if b531OrdinaryLimitError != nil {
					return nil, b531OrdinaryLimitError
				}
				b531OrdinaryLimitArray = append(b531OrdinaryLimitArray, b531OrdinaryLimitValue)
			}
			_, b531OrdinaryLimitError := b531OrdinaryLimitDecoder.Token()
			return b531OrdinaryLimitArray, b531OrdinaryLimitError
		}
		return nil, b531OrdinaryLimitFmt.Errorf("unexpected delimiter")
	}
	return b531OrdinaryLimitToken, nil
}

func b531OrdinaryLimitUint(b531OrdinaryLimitValue any, b531OrdinaryLimitBits int) (uint64, error) {
	b531OrdinaryLimitNumber, b531OrdinaryLimitOK := b531OrdinaryLimitValue.(b531OrdinaryLimitJSON.Number)
	if !b531OrdinaryLimitOK {
		return 0, b531OrdinaryLimitFmt.Errorf("unsigned integer required")
	}
	return b531OrdinaryLimitStrconv.ParseUint(string(b531OrdinaryLimitNumber), 10, b531OrdinaryLimitBits)
}

func B531OrdinaryLimitEvaluate(b531OrdinaryLimitInput b531OrdinaryLimitJSON.RawMessage) (b531OrdinaryLimitJSON.RawMessage, error) {
	b531OrdinaryLimitDecoder := b531OrdinaryLimitJSON.NewDecoder(b531OrdinaryLimitBytesPackage.NewReader(b531OrdinaryLimitInput))
	b531OrdinaryLimitDecoder.UseNumber()
	b531OrdinaryLimitValue, b531OrdinaryLimitError := b531OrdinaryLimitDecode(b531OrdinaryLimitDecoder)
	if b531OrdinaryLimitError != nil {
		return nil, b531OrdinaryLimitError
	}
	if _, b531OrdinaryLimitError = b531OrdinaryLimitDecoder.Token(); b531OrdinaryLimitError != b531OrdinaryLimitIO.EOF {
		return nil, b531OrdinaryLimitFmt.Errorf("trailing JSON")
	}
	b531OrdinaryLimitInputObject, b531OrdinaryLimitOK := b531OrdinaryLimitValue.(map[string]any)
	if !b531OrdinaryLimitOK {
		return nil, b531OrdinaryLimitFmt.Errorf("input object required")
	}
	b531OrdinaryLimitContext, b531OrdinaryLimitOK := b531OrdinaryLimitInputObject["context"].(map[string]any)
	if !b531OrdinaryLimitOK {
		return nil, b531OrdinaryLimitFmt.Errorf("context object required")
	}
	b531OrdinaryLimitRaw, b531OrdinaryLimitOK := b531OrdinaryLimitInputObject["fields"].(map[string]any)
	if !b531OrdinaryLimitOK {
		return nil, b531OrdinaryLimitFmt.Errorf("fields object required")
	}
	for b531OrdinaryLimitKey, b531OrdinaryLimitValue := range b531OrdinaryLimitRaw {
		b531OrdinaryLimitBytes, b531OrdinaryLimitOK := b531OrdinaryLimitValue.([]any)
		if !b531OrdinaryLimitOK {
			return nil, b531OrdinaryLimitFmt.Errorf("field %s must be a byte array", b531OrdinaryLimitKey)
		}
		for _, b531OrdinaryLimitByte := range b531OrdinaryLimitBytes {
			if _, b531OrdinaryLimitError := b531OrdinaryLimitUint(b531OrdinaryLimitByte, 8); b531OrdinaryLimitError != nil {
				return nil, b531OrdinaryLimitError
			}
		}
	}
	b531OrdinaryLimitSemantic := map[string]any{}
	b531OrdinaryLimitOutput := map[string]any{"context": b531OrdinaryLimitContext, "raw": b531OrdinaryLimitRaw, "state": "Unknown", "fields": b531OrdinaryLimitSemantic, "physical": "Unknown", "qualified": false}
	b531OrdinaryLimitArguments := map[string]any{}
	if b531OrdinaryLimitArgumentValue, b531OrdinaryLimitPresent := b531OrdinaryLimitInputObject["arguments"]; b531OrdinaryLimitPresent {
		b531OrdinaryLimitArguments, b531OrdinaryLimitOK = b531OrdinaryLimitArgumentValue.(map[string]any)
		if !b531OrdinaryLimitOK {
			return nil, b531OrdinaryLimitFmt.Errorf("arguments object required")
		}
		b531OrdinaryLimitOutput["arguments"] = b531OrdinaryLimitArguments
	}
	b531OrdinaryLimitSupported := true
	for b531OrdinaryLimitKey, b531OrdinaryLimitExpected := range map[string]uint64{"primary": 181, "secondary": 49, "operation": 3} {
		b531OrdinaryLimitValue, b531OrdinaryLimitPresent := b531OrdinaryLimitContext[b531OrdinaryLimitKey]
		if !b531OrdinaryLimitPresent || b531OrdinaryLimitValue == nil {
			b531OrdinaryLimitSupported = false
			continue
		}
		b531OrdinaryLimitNumber, b531OrdinaryLimitError := b531OrdinaryLimitUint(b531OrdinaryLimitValue, 8)
		if b531OrdinaryLimitError != nil {
			return nil, b531OrdinaryLimitError
		}
		if b531OrdinaryLimitNumber != b531OrdinaryLimitExpected {
			b531OrdinaryLimitSupported = false
		}
	}
	for b531OrdinaryLimitKey, b531OrdinaryLimitExpected := range map[string]string{"direction": "local", "view": "ordinary-limit"} {
		b531OrdinaryLimitValue, b531OrdinaryLimitPresent := b531OrdinaryLimitContext[b531OrdinaryLimitKey]
		if !b531OrdinaryLimitPresent || b531OrdinaryLimitValue == nil {
			b531OrdinaryLimitSupported = false
			continue
		}
		b531OrdinaryLimitString, b531OrdinaryLimitOK := b531OrdinaryLimitValue.(string)
		if !b531OrdinaryLimitOK {
			return nil, b531OrdinaryLimitFmt.Errorf("context %s must be a string", b531OrdinaryLimitKey)
		}
		if b531OrdinaryLimitString != b531OrdinaryLimitExpected {
			b531OrdinaryLimitSupported = false
		}
	}
	if !b531OrdinaryLimitSupported {
		if b531OrdinaryLimitError := b531OrdinaryLimitValidateOptional(b531OrdinaryLimitArguments); b531OrdinaryLimitError != nil {
			return nil, b531OrdinaryLimitError
		}
		return b531OrdinaryLimitJSON.Marshal(b531OrdinaryLimitOutput)
	}
	b531OrdinaryLimitLimited, b531OrdinaryLimitError := b531OrdinaryLimitUint(b531OrdinaryLimitArguments["limited"], 32)
	if b531OrdinaryLimitError != nil {
		return nil, b531OrdinaryLimitError
	}
	b531OrdinaryLimitRows, b531OrdinaryLimitOK := b531OrdinaryLimitArguments["rows"].([]any)
	if !b531OrdinaryLimitOK || len(b531OrdinaryLimitRows) != 9 {
		return nil, b531OrdinaryLimitFmt.Errorf("rows requires nine slots")
	}
	b531OrdinaryLimitIndices := []int{}
	b531OrdinaryLimitSlots := []any{}
	var b531OrdinaryLimitTokenSum, b531OrdinaryLimitMaximumSum uint64
	b531OrdinaryLimitSpecial := false
	for b531OrdinaryLimitIndex, b531OrdinaryLimitValue := range b531OrdinaryLimitRows {
		b531OrdinaryLimitRow, b531OrdinaryLimitOK := b531OrdinaryLimitValue.(map[string]any)
		if !b531OrdinaryLimitOK {
			return nil, b531OrdinaryLimitFmt.Errorf("slot object required")
		}
		b531OrdinaryLimitPredicates := map[string]bool{}
		for _, b531OrdinaryLimitKey := range []string{"registered", "capable", "failure"} {
			b531OrdinaryLimitPredicate, b531OrdinaryLimitOK := b531OrdinaryLimitRow[b531OrdinaryLimitKey].(bool)
			if !b531OrdinaryLimitOK {
				return nil, b531OrdinaryLimitFmt.Errorf("boolean %s required", b531OrdinaryLimitKey)
			}
			b531OrdinaryLimitPredicates[b531OrdinaryLimitKey] = b531OrdinaryLimitPredicate
		}
		b531OrdinaryLimitToken, b531OrdinaryLimitError := b531OrdinaryLimitUint(b531OrdinaryLimitRow["token"], 32)
		if b531OrdinaryLimitError != nil {
			return nil, b531OrdinaryLimitError
		}
		b531OrdinaryLimitMaximum, b531OrdinaryLimitError := b531OrdinaryLimitUint(b531OrdinaryLimitRow["maximum"], 32)
		if b531OrdinaryLimitError != nil {
			return nil, b531OrdinaryLimitError
		}
		b531OrdinaryLimitAdmitted := b531OrdinaryLimitPredicates["registered"] && b531OrdinaryLimitPredicates["capable"] && !b531OrdinaryLimitPredicates["failure"]
		b531OrdinaryLimitSlots = append(b531OrdinaryLimitSlots, map[string]any{"slot": b531OrdinaryLimitIndex, "raw": b531OrdinaryLimitRow, "admitted": b531OrdinaryLimitAdmitted, "state": "Hypothesis"})
		if !b531OrdinaryLimitAdmitted {
			continue
		}
		b531OrdinaryLimitIndices = append(b531OrdinaryLimitIndices, b531OrdinaryLimitIndex)
		if b531OrdinaryLimitIndex == 1 {
			b531OrdinaryLimitSpecial = true
			continue
		}
		b531OrdinaryLimitTokenSum += b531OrdinaryLimitToken
		b531OrdinaryLimitMaximumSum += b531OrdinaryLimitMaximum
	}
	b531OrdinaryLimitSemantic["slots"] = b531OrdinaryLimitSlots
	b531OrdinaryLimitSemantic["admitted_slots"] = map[string]any{"raw": b531OrdinaryLimitRows, "state": "Hypothesis", "indices": b531OrdinaryLimitIndices}
	b531OrdinaryLimitSemantic["token_sum"] = map[string]any{"raw": b531OrdinaryLimitRows, "state": "Hypothesis", "calculation": "Proven", "number": b531OrdinaryLimitTokenSum}
	b531OrdinaryLimitSemantic["maximum_sum"] = map[string]any{"raw": b531OrdinaryLimitRows, "state": "Hypothesis", "calculation": "Proven", "number": b531OrdinaryLimitMaximumSum}
	b531OrdinaryLimitEffective := b531OrdinaryLimitTokenSum < b531OrdinaryLimitMaximumSum && b531OrdinaryLimitLimited != 0
	b531OrdinaryLimitEffectiveResult := map[string]any{"raw": b531OrdinaryLimitRows, "state": "Unknown"}
	b531OrdinaryLimitRequested := map[string]any{"state": "Unknown"}
	if b531OrdinaryLimitGate, b531OrdinaryLimitPresent := b531OrdinaryLimitRaw["gate"]; b531OrdinaryLimitPresent {
		b531OrdinaryLimitRequested["raw"] = b531OrdinaryLimitGate
	}
	b531OrdinaryLimitOutput["state"] = "complete"
	if b531OrdinaryLimitSpecial {
		b531OrdinaryLimitOutput["state"] = "partial"
	} else {
		b531OrdinaryLimitEffectiveResult["state"] = "Hypothesis"
		b531OrdinaryLimitEffectiveResult["calculation"] = "Proven"
		b531OrdinaryLimitEffectiveResult["value"] = b531OrdinaryLimitEffective
	}
	b531OrdinaryLimitGate, b531OrdinaryLimitOK := b531OrdinaryLimitRaw["gate"].([]any)
	if !b531OrdinaryLimitOK || len(b531OrdinaryLimitGate) != 1 {
		b531OrdinaryLimitOutput["state"] = "partial"
	} else {
		b531OrdinaryLimitGateNumber, _ := b531OrdinaryLimitUint(b531OrdinaryLimitGate[0], 8)
		if b531OrdinaryLimitGateNumber == 0 {
			b531OrdinaryLimitRequested["state"] = "Hypothesis"
			b531OrdinaryLimitRequested["number"] = 0
			b531OrdinaryLimitRequested["calculation"] = "Proven"
		} else if !b531OrdinaryLimitSpecial {
			b531OrdinaryLimitNumber := 1
			if b531OrdinaryLimitEffective {
				b531OrdinaryLimitNumber = 2
			}
			b531OrdinaryLimitRequested["state"] = "Hypothesis"
			b531OrdinaryLimitRequested["number"] = b531OrdinaryLimitNumber
			b531OrdinaryLimitRequested["calculation"] = "Proven"
		}
	}
	b531OrdinaryLimitSemantic["effective"] = b531OrdinaryLimitEffectiveResult
	b531OrdinaryLimitSemantic["requested_state"] = b531OrdinaryLimitRequested
	return b531OrdinaryLimitJSON.Marshal(b531OrdinaryLimitOutput)
}

func b531OrdinaryLimitValidateOptional(b531OrdinaryLimitArguments map[string]any) error {
	if b531OrdinaryLimitValue, b531OrdinaryLimitPresent := b531OrdinaryLimitArguments["limited"]; b531OrdinaryLimitPresent {
		if _, b531OrdinaryLimitError := b531OrdinaryLimitUint(b531OrdinaryLimitValue, 32); b531OrdinaryLimitError != nil {
			return b531OrdinaryLimitError
		}
	}
	if b531OrdinaryLimitValue, b531OrdinaryLimitPresent := b531OrdinaryLimitArguments["rows"]; b531OrdinaryLimitPresent {
		b531OrdinaryLimitRows, b531OrdinaryLimitOK := b531OrdinaryLimitValue.([]any)
		if !b531OrdinaryLimitOK || len(b531OrdinaryLimitRows) != 9 {
			return b531OrdinaryLimitFmt.Errorf("rows requires nine slots")
		}
		for _, b531OrdinaryLimitValue := range b531OrdinaryLimitRows {
			b531OrdinaryLimitRow, b531OrdinaryLimitOK := b531OrdinaryLimitValue.(map[string]any)
			if !b531OrdinaryLimitOK {
				return b531OrdinaryLimitFmt.Errorf("slot object required")
			}
			for _, b531OrdinaryLimitKey := range []string{"registered", "capable", "failure"} {
				if _, b531OrdinaryLimitOK := b531OrdinaryLimitRow[b531OrdinaryLimitKey].(bool); !b531OrdinaryLimitOK {
					return b531OrdinaryLimitFmt.Errorf("boolean %s required", b531OrdinaryLimitKey)
				}
			}
			for _, b531OrdinaryLimitKey := range []string{"token", "maximum"} {
				if _, b531OrdinaryLimitError := b531OrdinaryLimitUint(b531OrdinaryLimitRow[b531OrdinaryLimitKey], 32); b531OrdinaryLimitError != nil {
					return b531OrdinaryLimitError
				}
			}
		}
	}
	return nil
}
