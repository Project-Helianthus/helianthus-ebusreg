package vaillant

import (
	b51AMirrorFieldEvaluatorBytesPackage "bytes"
	b51AMirrorFieldEvaluatorJSON "encoding/json"
	b51AMirrorFieldEvaluatorFmt "fmt"
	b51AMirrorFieldEvaluatorIO "io"
	b51AMirrorFieldEvaluatorBig "math/big"
)

func b51AMirrorFieldEvaluatorRead(b51AMirrorFieldEvaluatorDecoder *b51AMirrorFieldEvaluatorJSON.Decoder) (any, error) {
	b51AMirrorFieldEvaluatorToken, b51AMirrorFieldEvaluatorErr := b51AMirrorFieldEvaluatorDecoder.Token()
	if b51AMirrorFieldEvaluatorErr != nil {
		return nil, b51AMirrorFieldEvaluatorErr
	}
	switch b51AMirrorFieldEvaluatorToken {
	case b51AMirrorFieldEvaluatorJSON.Delim('{'):
		b51AMirrorFieldEvaluatorObject := map[string]any{}
		for b51AMirrorFieldEvaluatorDecoder.More() {
			b51AMirrorFieldEvaluatorKeyToken, b51AMirrorFieldEvaluatorErr := b51AMirrorFieldEvaluatorDecoder.Token()
			if b51AMirrorFieldEvaluatorErr != nil {
				return nil, b51AMirrorFieldEvaluatorErr
			}
			b51AMirrorFieldEvaluatorKey := b51AMirrorFieldEvaluatorKeyToken.(string)
			if _, b51AMirrorFieldEvaluatorExists := b51AMirrorFieldEvaluatorObject[b51AMirrorFieldEvaluatorKey]; b51AMirrorFieldEvaluatorExists {
				return nil, b51AMirrorFieldEvaluatorFmt.Errorf("duplicate key %q", b51AMirrorFieldEvaluatorKey)
			}
			b51AMirrorFieldEvaluatorValue, b51AMirrorFieldEvaluatorErr := b51AMirrorFieldEvaluatorRead(b51AMirrorFieldEvaluatorDecoder)
			if b51AMirrorFieldEvaluatorErr != nil {
				return nil, b51AMirrorFieldEvaluatorErr
			}
			b51AMirrorFieldEvaluatorObject[b51AMirrorFieldEvaluatorKey] = b51AMirrorFieldEvaluatorValue
		}
		_, b51AMirrorFieldEvaluatorErr = b51AMirrorFieldEvaluatorDecoder.Token()
		return b51AMirrorFieldEvaluatorObject, b51AMirrorFieldEvaluatorErr
	case b51AMirrorFieldEvaluatorJSON.Delim('['):
		b51AMirrorFieldEvaluatorArray := []any{}
		for b51AMirrorFieldEvaluatorDecoder.More() {
			b51AMirrorFieldEvaluatorValue, b51AMirrorFieldEvaluatorErr := b51AMirrorFieldEvaluatorRead(b51AMirrorFieldEvaluatorDecoder)
			if b51AMirrorFieldEvaluatorErr != nil {
				return nil, b51AMirrorFieldEvaluatorErr
			}
			b51AMirrorFieldEvaluatorArray = append(b51AMirrorFieldEvaluatorArray, b51AMirrorFieldEvaluatorValue)
		}
		_, b51AMirrorFieldEvaluatorErr = b51AMirrorFieldEvaluatorDecoder.Token()
		return b51AMirrorFieldEvaluatorArray, b51AMirrorFieldEvaluatorErr
	}
	return b51AMirrorFieldEvaluatorToken, nil
}

func b51AMirrorFieldEvaluatorUint(b51AMirrorFieldEvaluatorValue any, b51AMirrorFieldEvaluatorMax float64) (uint64, bool) {
	b51AMirrorFieldEvaluatorNumber, b51AMirrorFieldEvaluatorOK := b51AMirrorFieldEvaluatorValue.(b51AMirrorFieldEvaluatorJSON.Number)
	if !b51AMirrorFieldEvaluatorOK {
		return 0, false
	}
	b51AMirrorFieldEvaluatorRational, b51AMirrorFieldEvaluatorOK := new(b51AMirrorFieldEvaluatorBig.Rat).SetString(string(b51AMirrorFieldEvaluatorNumber))
	if !b51AMirrorFieldEvaluatorOK || !b51AMirrorFieldEvaluatorRational.IsInt() || b51AMirrorFieldEvaluatorRational.Sign() < 0 || !b51AMirrorFieldEvaluatorRational.Num().IsUint64() || b51AMirrorFieldEvaluatorRational.Num().Uint64() > uint64(b51AMirrorFieldEvaluatorMax) {
		return 0, false
	}
	return b51AMirrorFieldEvaluatorRational.Num().Uint64(), true
}

func B51AMirrorFieldEvaluate(b51AMirrorFieldEvaluatorInput b51AMirrorFieldEvaluatorJSON.RawMessage) (b51AMirrorFieldEvaluatorJSON.RawMessage, error) {
	b51AMirrorFieldEvaluatorDecoder := b51AMirrorFieldEvaluatorJSON.NewDecoder(b51AMirrorFieldEvaluatorBytesPackage.NewReader(b51AMirrorFieldEvaluatorInput))
	b51AMirrorFieldEvaluatorDecoder.UseNumber()
	b51AMirrorFieldEvaluatorValue, b51AMirrorFieldEvaluatorErr := b51AMirrorFieldEvaluatorRead(b51AMirrorFieldEvaluatorDecoder)
	if b51AMirrorFieldEvaluatorErr != nil {
		return nil, b51AMirrorFieldEvaluatorErr
	}
	if _, b51AMirrorFieldEvaluatorErr = b51AMirrorFieldEvaluatorDecoder.Token(); b51AMirrorFieldEvaluatorErr != b51AMirrorFieldEvaluatorIO.EOF {
		return nil, b51AMirrorFieldEvaluatorFmt.Errorf("trailing JSON")
	}
	b51AMirrorFieldEvaluatorObject, b51AMirrorFieldEvaluatorOK := b51AMirrorFieldEvaluatorValue.(map[string]any)
	if !b51AMirrorFieldEvaluatorOK {
		return nil, b51AMirrorFieldEvaluatorFmt.Errorf("input must be object")
	}
	b51AMirrorFieldEvaluatorContext, b51AMirrorFieldEvaluatorOK := b51AMirrorFieldEvaluatorObject["context"].(map[string]any)
	if !b51AMirrorFieldEvaluatorOK {
		return nil, b51AMirrorFieldEvaluatorFmt.Errorf("context must be object")
	}
	b51AMirrorFieldEvaluatorRaw, b51AMirrorFieldEvaluatorOK := b51AMirrorFieldEvaluatorObject["fields"].(map[string]any)
	if !b51AMirrorFieldEvaluatorOK {
		return nil, b51AMirrorFieldEvaluatorFmt.Errorf("fields must be object")
	}
	for _, b51AMirrorFieldEvaluatorKey := range []string{"primary", "secondary", "operation"} {
		if b51AMirrorFieldEvaluatorValue, b51AMirrorFieldEvaluatorPresent := b51AMirrorFieldEvaluatorContext[b51AMirrorFieldEvaluatorKey]; b51AMirrorFieldEvaluatorPresent && b51AMirrorFieldEvaluatorValue != nil {
			if _, b51AMirrorFieldEvaluatorOK := b51AMirrorFieldEvaluatorUint(b51AMirrorFieldEvaluatorValue, 255); !b51AMirrorFieldEvaluatorOK {
				return nil, b51AMirrorFieldEvaluatorFmt.Errorf("invalid context %s", b51AMirrorFieldEvaluatorKey)
			}
		}
	}
	for _, b51AMirrorFieldEvaluatorKey := range []string{"direction", "view"} {
		if b51AMirrorFieldEvaluatorValue, b51AMirrorFieldEvaluatorPresent := b51AMirrorFieldEvaluatorContext[b51AMirrorFieldEvaluatorKey]; b51AMirrorFieldEvaluatorPresent && b51AMirrorFieldEvaluatorValue != nil {
			if _, b51AMirrorFieldEvaluatorOK := b51AMirrorFieldEvaluatorValue.(string); !b51AMirrorFieldEvaluatorOK {
				return nil, b51AMirrorFieldEvaluatorFmt.Errorf("invalid context %s", b51AMirrorFieldEvaluatorKey)
			}
		}
	}
	for b51AMirrorFieldEvaluatorKey, b51AMirrorFieldEvaluatorValue := range b51AMirrorFieldEvaluatorRaw {
		b51AMirrorFieldEvaluatorArray, b51AMirrorFieldEvaluatorOK := b51AMirrorFieldEvaluatorValue.([]any)
		if !b51AMirrorFieldEvaluatorOK {
			return nil, b51AMirrorFieldEvaluatorFmt.Errorf("field %s must be byte array", b51AMirrorFieldEvaluatorKey)
		}
		for _, b51AMirrorFieldEvaluatorByte := range b51AMirrorFieldEvaluatorArray {
			if _, b51AMirrorFieldEvaluatorOK := b51AMirrorFieldEvaluatorUint(b51AMirrorFieldEvaluatorByte, 255); !b51AMirrorFieldEvaluatorOK {
				return nil, b51AMirrorFieldEvaluatorFmt.Errorf("invalid byte in %s", b51AMirrorFieldEvaluatorKey)
			}
		}
	}
	b51AMirrorFieldEvaluatorOutput := map[string]any{"context": b51AMirrorFieldEvaluatorContext, "raw": b51AMirrorFieldEvaluatorRaw, "fields": map[string]any{}, "state": "Unknown", "physical": "Unknown", "qualified": false}
	if b51AMirrorFieldEvaluatorArguments, b51AMirrorFieldEvaluatorPresent := b51AMirrorFieldEvaluatorObject["arguments"]; b51AMirrorFieldEvaluatorPresent {
		if _, b51AMirrorFieldEvaluatorOK := b51AMirrorFieldEvaluatorArguments.(map[string]any); !b51AMirrorFieldEvaluatorOK {
			return nil, b51AMirrorFieldEvaluatorFmt.Errorf("arguments must be object")
		}
		b51AMirrorFieldEvaluatorOutput["arguments"] = b51AMirrorFieldEvaluatorArguments
	}
	b51AMirrorFieldEvaluatorPrimary, _ := b51AMirrorFieldEvaluatorUint(b51AMirrorFieldEvaluatorContext["primary"], 255)
	b51AMirrorFieldEvaluatorSecondary, _ := b51AMirrorFieldEvaluatorUint(b51AMirrorFieldEvaluatorContext["secondary"], 255)
	b51AMirrorFieldEvaluatorOperation, _ := b51AMirrorFieldEvaluatorUint(b51AMirrorFieldEvaluatorContext["operation"], 255)
	b51AMirrorFieldEvaluatorDirection := b51AMirrorFieldEvaluatorContext["direction"]
	b51AMirrorFieldEvaluatorView := b51AMirrorFieldEvaluatorContext["view"]
	if b51AMirrorFieldEvaluatorPrimary != 181 || b51AMirrorFieldEvaluatorSecondary != 26 || ((b51AMirrorFieldEvaluatorOperation != 2 || b51AMirrorFieldEvaluatorView != "mirror-get") && (b51AMirrorFieldEvaluatorOperation != 3 || b51AMirrorFieldEvaluatorView != "mirror-put")) || (b51AMirrorFieldEvaluatorDirection != "request" && b51AMirrorFieldEvaluatorDirection != "positive") {
		return b51AMirrorFieldEvaluatorJSON.Marshal(b51AMirrorFieldEvaluatorOutput)
	}
	b51AMirrorFieldEvaluatorWidths := map[string]int{"sequence": 1}
	if b51AMirrorFieldEvaluatorDirection == "request" {
		b51AMirrorFieldEvaluatorWidths["flags"] = 1
		if b51AMirrorFieldEvaluatorOperation == 2 {
			b51AMirrorFieldEvaluatorWidths["length"] = 1
			b51AMirrorFieldEvaluatorWidths["address"] = 2
		} else {
			b51AMirrorFieldEvaluatorWidths["mirror_address"] = 2
			b51AMirrorFieldEvaluatorWidths["symbol_address"] = 2
			b51AMirrorFieldEvaluatorWidths["data"] = -1
		}
	} else if b51AMirrorFieldEvaluatorOperation == 2 {
		b51AMirrorFieldEvaluatorWidths["status"] = 1
		b51AMirrorFieldEvaluatorWidths["data"] = -1
	} else {
		b51AMirrorFieldEvaluatorWidths["flags"] = 1
	}
	b51AMirrorFieldEvaluatorResults := b51AMirrorFieldEvaluatorOutput["fields"].(map[string]any)
	b51AMirrorFieldEvaluatorOutput["state"] = "complete"
	for b51AMirrorFieldEvaluatorKey, b51AMirrorFieldEvaluatorWidth := range b51AMirrorFieldEvaluatorWidths {
		b51AMirrorFieldEvaluatorResult := map[string]any{"state": "incomplete"}
		b51AMirrorFieldEvaluatorResults[b51AMirrorFieldEvaluatorKey] = b51AMirrorFieldEvaluatorResult
		b51AMirrorFieldEvaluatorValue, b51AMirrorFieldEvaluatorPresent := b51AMirrorFieldEvaluatorRaw[b51AMirrorFieldEvaluatorKey]
		if b51AMirrorFieldEvaluatorPresent {
			b51AMirrorFieldEvaluatorResult["raw"] = b51AMirrorFieldEvaluatorValue
		}
		b51AMirrorFieldEvaluatorArray, _ := b51AMirrorFieldEvaluatorValue.([]any)
		if !b51AMirrorFieldEvaluatorPresent || (b51AMirrorFieldEvaluatorWidth >= 0 && len(b51AMirrorFieldEvaluatorArray) != b51AMirrorFieldEvaluatorWidth) {
			b51AMirrorFieldEvaluatorOutput["state"] = "partial"
			continue
		}
		b51AMirrorFieldEvaluatorResult["state"] = "Hypothesis"
		if b51AMirrorFieldEvaluatorKey == "data" {
			b51AMirrorFieldEvaluatorResult["label"] = "opaque"
			continue
		}
		b51AMirrorFieldEvaluatorNumber, _ := b51AMirrorFieldEvaluatorUint(b51AMirrorFieldEvaluatorArray[0], 255)
		if b51AMirrorFieldEvaluatorWidth == 2 {
			b51AMirrorFieldEvaluatorHigh, _ := b51AMirrorFieldEvaluatorUint(b51AMirrorFieldEvaluatorArray[1], 255)
			b51AMirrorFieldEvaluatorNumber += 256 * b51AMirrorFieldEvaluatorHigh
			b51AMirrorFieldEvaluatorResult["logical_validity"] = "Unknown"
		}
		b51AMirrorFieldEvaluatorResult["number"] = b51AMirrorFieldEvaluatorNumber
		b51AMirrorFieldEvaluatorResult["calculation"] = "Proven"
		if b51AMirrorFieldEvaluatorKey == "status" {
			b51AMirrorFieldEvaluatorBoolean := map[string]any{"state": "Unknown"}
			if b51AMirrorFieldEvaluatorNumber <= 1 {
				b51AMirrorFieldEvaluatorBoolean["state"] = "Hypothesis"
				b51AMirrorFieldEvaluatorBoolean["value"] = b51AMirrorFieldEvaluatorNumber == 1
			}
			b51AMirrorFieldEvaluatorResult["alternatives"] = map[string]any{"whole_boolean": b51AMirrorFieldEvaluatorBoolean, "bit0_out_of_area": map[string]any{"state": "Hypothesis", "value": b51AMirrorFieldEvaluatorNumber&1 != 0}}
		}
		if b51AMirrorFieldEvaluatorKey == "flags" {
			b51AMirrorFieldEvaluatorNames := []string{"triple-storage"}
			if b51AMirrorFieldEvaluatorOperation == 3 {
				if b51AMirrorFieldEvaluatorDirection == "request" {
					b51AMirrorFieldEvaluatorNames = append(b51AMirrorFieldEvaluatorNames, "direct-update")
				} else {
					b51AMirrorFieldEvaluatorNames = []string{"not-written", "invalid-symbol", "invalid-parameter", "storage-error"}
				}
			}
			b51AMirrorFieldEvaluatorFlags := map[string]bool{}
			for b51AMirrorFieldEvaluatorBit, b51AMirrorFieldEvaluatorName := range b51AMirrorFieldEvaluatorNames {
				b51AMirrorFieldEvaluatorFlags[b51AMirrorFieldEvaluatorName] = b51AMirrorFieldEvaluatorNumber&(1<<b51AMirrorFieldEvaluatorBit) != 0
			}
			b51AMirrorFieldEvaluatorResult["flags"] = b51AMirrorFieldEvaluatorFlags
			b51AMirrorFieldEvaluatorRest := b51AMirrorFieldEvaluatorNumber & (255 << len(b51AMirrorFieldEvaluatorNames))
			if b51AMirrorFieldEvaluatorOperation == 3 {
				b51AMirrorFieldEvaluatorResult["reserved"] = b51AMirrorFieldEvaluatorRest
			} else {
				b51AMirrorFieldEvaluatorResult["unresolved"] = b51AMirrorFieldEvaluatorRest
			}
		}
	}
	return b51AMirrorFieldEvaluatorJSON.Marshal(b51AMirrorFieldEvaluatorOutput)
}
