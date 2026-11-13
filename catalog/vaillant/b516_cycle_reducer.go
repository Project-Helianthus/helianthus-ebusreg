package vaillant

import (
	b516CycleReducerBytes "bytes"
	b516CycleReducerJSON "encoding/json"
	b516CycleReducerFmt "fmt"
	b516CycleReducerIO "io"
	b516CycleReducerStrconv "strconv"
)

func b516CycleReducerRead(b516CycleReducerDecoder *b516CycleReducerJSON.Decoder) (any, error) {
	b516CycleReducerToken, b516CycleReducerErr := b516CycleReducerDecoder.Token()
	if b516CycleReducerErr != nil {
		return nil, b516CycleReducerErr
	}
	switch b516CycleReducerToken {
	case b516CycleReducerJSON.Delim('{'):
		b516CycleReducerObject := map[string]any{}
		for b516CycleReducerDecoder.More() {
			b516CycleReducerKeyToken, b516CycleReducerErr := b516CycleReducerDecoder.Token()
			if b516CycleReducerErr != nil {
				return nil, b516CycleReducerErr
			}
			b516CycleReducerKey := b516CycleReducerKeyToken.(string)
			if _, b516CycleReducerExists := b516CycleReducerObject[b516CycleReducerKey]; b516CycleReducerExists {
				return nil, b516CycleReducerFmt.Errorf("duplicate key %q", b516CycleReducerKey)
			}
			b516CycleReducerValue, b516CycleReducerErr := b516CycleReducerRead(b516CycleReducerDecoder)
			if b516CycleReducerErr != nil {
				return nil, b516CycleReducerErr
			}
			b516CycleReducerObject[b516CycleReducerKey] = b516CycleReducerValue
		}
		_, b516CycleReducerErr = b516CycleReducerDecoder.Token()
		return b516CycleReducerObject, b516CycleReducerErr
	case b516CycleReducerJSON.Delim('['):
		b516CycleReducerArray := []any{}
		for b516CycleReducerDecoder.More() {
			b516CycleReducerValue, b516CycleReducerErr := b516CycleReducerRead(b516CycleReducerDecoder)
			if b516CycleReducerErr != nil {
				return nil, b516CycleReducerErr
			}
			b516CycleReducerArray = append(b516CycleReducerArray, b516CycleReducerValue)
		}
		_, b516CycleReducerErr = b516CycleReducerDecoder.Token()
		return b516CycleReducerArray, b516CycleReducerErr
	}
	if _, b516CycleReducerDelimiter := b516CycleReducerToken.(b516CycleReducerJSON.Delim); b516CycleReducerDelimiter {
		return nil, b516CycleReducerFmt.Errorf("unexpected delimiter")
	}
	return b516CycleReducerToken, nil
}

func b516CycleReducerUint(b516CycleReducerValue any, b516CycleReducerBits int) (uint64, error) {
	b516CycleReducerNumber, b516CycleReducerOK := b516CycleReducerValue.(b516CycleReducerJSON.Number)
	if !b516CycleReducerOK {
		return 0, b516CycleReducerFmt.Errorf("expected unsigned integer")
	}
	return b516CycleReducerStrconv.ParseUint(string(b516CycleReducerNumber), 10, b516CycleReducerBits)
}

func b516CycleReducerQuantity(b516CycleReducerNumber any) map[string]any {
	return map[string]any{"state": "Hypothesis", "calculation": "Proven", "number": b516CycleReducerNumber}
}

func B516CycleReducerEvaluate(b516CycleReducerInput b516CycleReducerJSON.RawMessage) (b516CycleReducerJSON.RawMessage, error) {
	b516CycleReducerDecoder := b516CycleReducerJSON.NewDecoder(b516CycleReducerBytes.NewReader(b516CycleReducerInput))
	b516CycleReducerDecoder.UseNumber()
	b516CycleReducerValue, b516CycleReducerErr := b516CycleReducerRead(b516CycleReducerDecoder)
	if b516CycleReducerErr != nil {
		return nil, b516CycleReducerErr
	}
	if _, b516CycleReducerErr = b516CycleReducerDecoder.Token(); b516CycleReducerErr != b516CycleReducerIO.EOF {
		return nil, b516CycleReducerFmt.Errorf("trailing JSON")
	}
	b516CycleReducerRoot, b516CycleReducerOK := b516CycleReducerValue.(map[string]any)
	if !b516CycleReducerOK {
		return nil, b516CycleReducerFmt.Errorf("expected object")
	}
	b516CycleReducerContext, b516CycleReducerOK := b516CycleReducerRoot["context"].(map[string]any)
	if !b516CycleReducerOK {
		return nil, b516CycleReducerFmt.Errorf("expected context object")
	}
	b516CycleReducerRaw, b516CycleReducerOK := b516CycleReducerRoot["fields"].(map[string]any)
	if !b516CycleReducerOK {
		return nil, b516CycleReducerFmt.Errorf("expected fields object")
	}
	b516CycleReducerFields := map[string]any{}
	for b516CycleReducerKey, b516CycleReducerValue := range b516CycleReducerRaw {
		b516CycleReducerArray, b516CycleReducerOK := b516CycleReducerValue.([]any)
		if !b516CycleReducerOK {
			return nil, b516CycleReducerFmt.Errorf("field %s: expected byte array", b516CycleReducerKey)
		}
		for _, b516CycleReducerByte := range b516CycleReducerArray {
			if _, b516CycleReducerErr := b516CycleReducerUint(b516CycleReducerByte, 8); b516CycleReducerErr != nil {
				return nil, b516CycleReducerErr
			}
		}
		b516CycleReducerFields[b516CycleReducerKey] = map[string]any{"raw": b516CycleReducerValue, "state": "Unknown"}
	}
	b516CycleReducerOutput := map[string]any{"context": b516CycleReducerContext, "raw": b516CycleReducerRaw, "state": "Unknown", "fields": b516CycleReducerFields, "physical": "Unknown", "qualified": false}
	b516CycleReducerArguments, b516CycleReducerSupplied := b516CycleReducerRoot["arguments"]
	if b516CycleReducerSupplied {
		if _, b516CycleReducerOK := b516CycleReducerArguments.(map[string]any); !b516CycleReducerOK {
			return nil, b516CycleReducerFmt.Errorf("expected arguments object")
		}
		b516CycleReducerOutput["arguments"] = b516CycleReducerArguments
	}
	b516CycleReducerSupported := true
	for b516CycleReducerKey, b516CycleReducerExpected := range map[string]uint64{"primary": 181, "secondary": 22, "operation": 14} {
		b516CycleReducerValue, b516CycleReducerPresent := b516CycleReducerContext[b516CycleReducerKey]
		if !b516CycleReducerPresent || b516CycleReducerValue == nil {
			b516CycleReducerSupported = false
			continue
		}
		b516CycleReducerNumber, b516CycleReducerErr := b516CycleReducerUint(b516CycleReducerValue, 8)
		if b516CycleReducerErr != nil {
			return nil, b516CycleReducerErr
		}
		if b516CycleReducerNumber != b516CycleReducerExpected {
			b516CycleReducerSupported = false
		}
	}
	for b516CycleReducerKey, b516CycleReducerExpected := range map[string]string{"direction": "local", "view": "cycle"} {
		b516CycleReducerValue := b516CycleReducerContext[b516CycleReducerKey]
		if b516CycleReducerValue == nil {
			b516CycleReducerSupported = false
			continue
		}
		b516CycleReducerString, b516CycleReducerOK := b516CycleReducerValue.(string)
		if !b516CycleReducerOK {
			return nil, b516CycleReducerFmt.Errorf("expected context string")
		}
		if b516CycleReducerString != b516CycleReducerExpected {
			b516CycleReducerSupported = false
		}
	}
	if !b516CycleReducerSupported {
		return b516CycleReducerJSON.Marshal(b516CycleReducerOutput)
	}
	b516CycleReducerArgs, b516CycleReducerOK := b516CycleReducerArguments.(map[string]any)
	if !b516CycleReducerOK {
		return nil, b516CycleReducerFmt.Errorf("missing arguments")
	}
	b516CycleReducerCounter, b516CycleReducerErr := b516CycleReducerUint(b516CycleReducerArgs["failure_counter"], 32)
	if b516CycleReducerErr != nil {
		return nil, b516CycleReducerErr
	}
	var b516CycleReducerIncrement uint64
	b516CycleReducerIncrementValue, b516CycleReducerHasIncrement := b516CycleReducerArgs["counter_increment"]
	if b516CycleReducerHasIncrement {
		b516CycleReducerIncrement, b516CycleReducerErr = b516CycleReducerUint(b516CycleReducerIncrementValue, 32)
		if b516CycleReducerErr != nil {
			return nil, b516CycleReducerErr
		}
	}
	b516CycleReducerSlots, b516CycleReducerOK := b516CycleReducerArgs["slots"].([]any)
	if !b516CycleReducerOK || len(b516CycleReducerSlots) != 9 {
		return nil, b516CycleReducerFmt.Errorf("expected nine slots")
	}
	b516CycleReducerResults := []any{}
	var b516CycleReducerSum uint32
	b516CycleReducerSuccessful, b516CycleReducerFailed, b516CycleReducerSkipped := 0, 0, 0
	for _, b516CycleReducerValue := range b516CycleReducerSlots {
		b516CycleReducerSlot, b516CycleReducerOK := b516CycleReducerValue.(map[string]any)
		if !b516CycleReducerOK {
			return nil, b516CycleReducerFmt.Errorf("expected slot object")
		}
		b516CycleReducerRegistered, b516CycleReducerOK := b516CycleReducerSlot["registered"].(bool)
		if !b516CycleReducerOK {
			return nil, b516CycleReducerFmt.Errorf("expected registered Boolean")
		}
		b516CycleReducerCapable, b516CycleReducerOK := b516CycleReducerSlot["capable"].(bool)
		if !b516CycleReducerOK {
			return nil, b516CycleReducerFmt.Errorf("expected capable Boolean")
		}
		b516CycleReducerPriorChanged, b516CycleReducerOK := b516CycleReducerSlot["prior_changed"].(bool)
		if !b516CycleReducerOK {
			return nil, b516CycleReducerFmt.Errorf("expected prior_changed Boolean")
		}
		b516CycleReducerPrior, b516CycleReducerErr := b516CycleReducerUint(b516CycleReducerSlot["prior"], 32)
		if b516CycleReducerErr != nil {
			return nil, b516CycleReducerErr
		}
		b516CycleReducerPrevious, b516CycleReducerErr := b516CycleReducerUint(b516CycleReducerSlot["prior_previous"], 32)
		if b516CycleReducerErr != nil {
			return nil, b516CycleReducerErr
		}
		b516CycleReducerOutcome, b516CycleReducerOK := b516CycleReducerSlot["outcome"].(string)
		if !b516CycleReducerOK || (b516CycleReducerOutcome != "success" && b516CycleReducerOutcome != "read-failure" && b516CycleReducerOutcome != "parse-failure") {
			return nil, b516CycleReducerFmt.Errorf("invalid outcome")
		}
		var b516CycleReducerCurrent uint64
		if b516CycleReducerOutcome == "success" {
			b516CycleReducerCurrent, b516CycleReducerErr = b516CycleReducerUint(b516CycleReducerSlot["value"], 32)
			if b516CycleReducerErr != nil {
				return nil, b516CycleReducerErr
			}
		} else if _, b516CycleReducerPresent := b516CycleReducerSlot["value"]; b516CycleReducerPresent {
			return nil, b516CycleReducerFmt.Errorf("value only allowed for success")
		}
		b516CycleReducerResult := map[string]any{"raw": b516CycleReducerSlot, "state": "Hypothesis", "prior_changed": b516CycleReducerPriorChanged, "changed": map[string]any{"state": "Unknown"}, "current": b516CycleReducerQuantity(uint32(0xffffffff)), "attempted": false, "outcome": "skipped"}
		if !b516CycleReducerRegistered || !b516CycleReducerCapable {
			b516CycleReducerSkipped++
		} else {
			b516CycleReducerPrevious = b516CycleReducerPrior
			b516CycleReducerResult["attempted"] = true
			b516CycleReducerResult["outcome"] = b516CycleReducerOutcome
			if b516CycleReducerOutcome != "success" {
				b516CycleReducerFailed++
			} else {
				b516CycleReducerSuccessful++
				b516CycleReducerSum += uint32(b516CycleReducerCurrent)
				b516CycleReducerResult["current"] = b516CycleReducerQuantity(b516CycleReducerCurrent)
				b516CycleReducerResult["changed"] = map[string]any{"state": "Hypothesis", "calculation": "Proven", "value": b516CycleReducerPrior == 0xffffffff || uint32(b516CycleReducerCurrent)-uint32(b516CycleReducerPrior) >= 2}
			}
		}
		b516CycleReducerResult["current"].(map[string]any)["raw"] = b516CycleReducerSlot
		b516CycleReducerResult["changed"].(map[string]any)["raw"] = b516CycleReducerSlot
		b516CycleReducerResult["previous"] = b516CycleReducerQuantity(b516CycleReducerPrevious)
		b516CycleReducerResult["previous"].(map[string]any)["raw"] = b516CycleReducerSlot
		b516CycleReducerResults = append(b516CycleReducerResults, b516CycleReducerResult)
	}
	b516CycleReducerOutput["state"] = "complete"
	b516CycleReducerFields["slots"] = map[string]any{"raw": b516CycleReducerSlots, "state": "Hypothesis", "outcomes": b516CycleReducerResults}
	b516CycleReducerFields["sum"] = b516CycleReducerQuantity(b516CycleReducerSum)
	b516CycleReducerFields["counts"] = map[string]any{"raw": b516CycleReducerSlots, "state": "Hypothesis", "successful": b516CycleReducerSuccessful, "failed": b516CycleReducerFailed, "skipped": b516CycleReducerSkipped, "attempted": b516CycleReducerSuccessful + b516CycleReducerFailed}
	b516CycleReducerOutput["partial_cycle"] = b516CycleReducerFailed+b516CycleReducerSkipped > 0
	b516CycleReducerFields["failure_counter"] = b516CycleReducerQuantity(uint32(0))
	if b516CycleReducerFailed > 0 {
		if b516CycleReducerHasIncrement {
			b516CycleReducerFields["failure_counter"] = b516CycleReducerQuantity(uint32(b516CycleReducerCounter + b516CycleReducerIncrement))
		} else {
			b516CycleReducerOutput["state"] = "partial"
			b516CycleReducerFields["failure_counter"] = map[string]any{"state": "Unknown", "prior": b516CycleReducerCounter}
		}
	}
	b516CycleReducerFields["sum"].(map[string]any)["raw"] = b516CycleReducerSlots
	b516CycleReducerFields["failure_counter"].(map[string]any)["raw"] = b516CycleReducerArgs["failure_counter"]
	if b516CycleReducerFailed > 0 {
		b516CycleReducerFields["failure_counter"].(map[string]any)["native_increments"] = "Unknown"
	}
	return b516CycleReducerJSON.Marshal(b516CycleReducerOutput)
}
