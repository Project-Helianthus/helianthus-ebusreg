package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func b531CapabilityStateRead(b531CapabilityStateDecoder *json.Decoder) (any, error) {
	b531CapabilityStateToken, b531CapabilityStateErr := b531CapabilityStateDecoder.Token()
	if b531CapabilityStateErr != nil {
		return nil, b531CapabilityStateErr
	}
	if b531CapabilityStateDelimiter, b531CapabilityStateOK := b531CapabilityStateToken.(json.Delim); b531CapabilityStateOK {
		switch b531CapabilityStateDelimiter {
		case '{':
			b531CapabilityStateObject := map[string]any{}
			for b531CapabilityStateDecoder.More() {
				b531CapabilityStateKeyToken, b531CapabilityStateErr := b531CapabilityStateDecoder.Token()
				if b531CapabilityStateErr != nil {
					return nil, b531CapabilityStateErr
				}
				b531CapabilityStateKey := b531CapabilityStateKeyToken.(string)
				if _, b531CapabilityStateExists := b531CapabilityStateObject[b531CapabilityStateKey]; b531CapabilityStateExists {
					return nil, fmt.Errorf("duplicate key %q", b531CapabilityStateKey)
				}
				b531CapabilityStateValue, b531CapabilityStateErr := b531CapabilityStateRead(b531CapabilityStateDecoder)
				if b531CapabilityStateErr != nil {
					return nil, b531CapabilityStateErr
				}
				b531CapabilityStateObject[b531CapabilityStateKey] = b531CapabilityStateValue
			}
			_, b531CapabilityStateErr := b531CapabilityStateDecoder.Token()
			return b531CapabilityStateObject, b531CapabilityStateErr
		case '[':
			b531CapabilityStateArray := []any{}
			for b531CapabilityStateDecoder.More() {
				b531CapabilityStateValue, b531CapabilityStateErr := b531CapabilityStateRead(b531CapabilityStateDecoder)
				if b531CapabilityStateErr != nil {
					return nil, b531CapabilityStateErr
				}
				b531CapabilityStateArray = append(b531CapabilityStateArray, b531CapabilityStateValue)
			}
			_, b531CapabilityStateErr := b531CapabilityStateDecoder.Token()
			return b531CapabilityStateArray, b531CapabilityStateErr
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	return b531CapabilityStateToken, nil
}

func b531CapabilityStateUint(b531CapabilityStateValue any, b531CapabilityStateBits int) (uint64, error) {
	b531CapabilityStateNumber, b531CapabilityStateOK := b531CapabilityStateValue.(json.Number)
	if !b531CapabilityStateOK {
		return 0, fmt.Errorf("expected unsigned integer")
	}
	return strconv.ParseUint(string(b531CapabilityStateNumber), 10, b531CapabilityStateBits)
}

func B531CapabilityStateEvaluate(b531CapabilityStateInput json.RawMessage) (json.RawMessage, error) {
	b531CapabilityStateDecoder := json.NewDecoder(bytes.NewReader(b531CapabilityStateInput))
	b531CapabilityStateDecoder.UseNumber()
	b531CapabilityStateValue, b531CapabilityStateErr := b531CapabilityStateRead(b531CapabilityStateDecoder)
	if b531CapabilityStateErr != nil {
		return nil, b531CapabilityStateErr
	}
	if _, b531CapabilityStateErr = b531CapabilityStateDecoder.Token(); b531CapabilityStateErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b531CapabilityStateRoot, b531CapabilityStateOK := b531CapabilityStateValue.(map[string]any)
	if !b531CapabilityStateOK {
		return nil, fmt.Errorf("expected object")
	}
	b531CapabilityStateContext, b531CapabilityStateOK := b531CapabilityStateRoot["context"].(map[string]any)
	if !b531CapabilityStateOK {
		return nil, fmt.Errorf("expected context object")
	}
	b531CapabilityStateRaw, b531CapabilityStateOK := b531CapabilityStateRoot["fields"].(map[string]any)
	if !b531CapabilityStateOK {
		return nil, fmt.Errorf("expected fields object")
	}
	b531CapabilityStateSupported := true
	for b531CapabilityStateKey, b531CapabilityStateExpected := range map[string]uint64{"primary": 181, "secondary": 49, "operation": 1} {
		b531CapabilityStateItem, b531CapabilityStateExists := b531CapabilityStateContext[b531CapabilityStateKey]
		if !b531CapabilityStateExists || b531CapabilityStateItem == nil {
			b531CapabilityStateSupported = false
			continue
		}
		b531CapabilityStateNumber, b531CapabilityStateErr := b531CapabilityStateUint(b531CapabilityStateItem, 8)
		if b531CapabilityStateErr != nil {
			return nil, b531CapabilityStateErr
		}
		if b531CapabilityStateNumber != b531CapabilityStateExpected {
			b531CapabilityStateSupported = false
		}
	}
	for _, b531CapabilityStateKey := range []string{"direction", "view"} {
		if b531CapabilityStateItem := b531CapabilityStateContext[b531CapabilityStateKey]; b531CapabilityStateItem != nil {
			if _, b531CapabilityStateOK := b531CapabilityStateItem.(string); !b531CapabilityStateOK {
				return nil, fmt.Errorf("expected string selector")
			}
		}
	}
	b531CapabilityStateView, _ := b531CapabilityStateContext["view"].(string)
	b531CapabilityStateSupported = b531CapabilityStateSupported && b531CapabilityStateContext["direction"] == "local" && (b531CapabilityStateView == "capability" || b531CapabilityStateView == "select")
	b531CapabilityStateFields := map[string]any{}
	b531CapabilityStateState := "complete"
	for b531CapabilityStateKey, b531CapabilityStateItem := range b531CapabilityStateRaw {
		b531CapabilityStateBytes, b531CapabilityStateOK := b531CapabilityStateItem.([]any)
		if !b531CapabilityStateOK {
			return nil, fmt.Errorf("expected byte array")
		}
		for _, b531CapabilityStateByte := range b531CapabilityStateBytes {
			if _, b531CapabilityStateErr := b531CapabilityStateUint(b531CapabilityStateByte, 8); b531CapabilityStateErr != nil {
				return nil, b531CapabilityStateErr
			}
		}
		b531CapabilityStateFields[b531CapabilityStateKey] = map[string]any{"raw": b531CapabilityStateItem, "state": "Unknown", "unmodeled": true}
	}
	if b531CapabilityStateSupported {
		for b531CapabilityStateKey, b531CapabilityStateItem := range b531CapabilityStateRaw {
			if len(b531CapabilityStateItem.([]any)) == 0 {
				b531CapabilityStateFields[b531CapabilityStateKey] = map[string]any{"raw": b531CapabilityStateItem, "state": "incomplete", "unmodeled": true}
				b531CapabilityStateState = "partial"
			}
		}
	}
	b531CapabilityStateOutput := map[string]any{"context": b531CapabilityStateContext, "raw": b531CapabilityStateRaw, "fields": b531CapabilityStateFields, "physical": "Unknown", "qualified": false}
	b531CapabilityStateArguments := map[string]any{}
	if b531CapabilityStateItem, b531CapabilityStateExists := b531CapabilityStateRoot["arguments"]; b531CapabilityStateExists {
		b531CapabilityStateArguments, b531CapabilityStateOK = b531CapabilityStateItem.(map[string]any)
		if !b531CapabilityStateOK {
			return nil, fmt.Errorf("expected arguments object")
		}
		b531CapabilityStateOutput["arguments"] = b531CapabilityStateArguments
	}
	for _, b531CapabilityStateKey := range []string{"registered", "prior", "result"} {
		if b531CapabilityStateItem, b531CapabilityStateExists := b531CapabilityStateArguments[b531CapabilityStateKey]; b531CapabilityStateExists {
			if _, b531CapabilityStateErr := b531CapabilityStateUint(b531CapabilityStateItem, 32); b531CapabilityStateErr != nil {
				return nil, b531CapabilityStateErr
			}
		}
	}
	b531CapabilityStateRows := []any(nil)
	if b531CapabilityStateItem, b531CapabilityStateExists := b531CapabilityStateArguments["rows"]; b531CapabilityStateExists {
		b531CapabilityStateRows, b531CapabilityStateOK = b531CapabilityStateItem.([]any)
		if !b531CapabilityStateOK || len(b531CapabilityStateRows) != 9 {
			return nil, fmt.Errorf("expected nine rows")
		}
		for _, b531CapabilityStateRow := range b531CapabilityStateRows {
			b531CapabilityStateObject, b531CapabilityStateOK := b531CapabilityStateRow.(map[string]any)
			if !b531CapabilityStateOK {
				return nil, fmt.Errorf("expected row object")
			}
			for _, b531CapabilityStateKey := range []string{"a", "b", "c"} {
				if _, b531CapabilityStateErr := b531CapabilityStateUint(b531CapabilityStateObject[b531CapabilityStateKey], 32); b531CapabilityStateErr != nil {
					return nil, b531CapabilityStateErr
				}
			}
		}
	}
	if !b531CapabilityStateSupported {
		b531CapabilityStateOutput["state"] = "Unknown"
		return json.Marshal(b531CapabilityStateOutput)
	}
	if b531CapabilityStateView == "capability" {
		b531CapabilityStateNumbers := map[string]uint64{}
		for _, b531CapabilityStateKey := range []string{"registered", "prior", "result"} {
			b531CapabilityStateNumber, b531CapabilityStateErr := b531CapabilityStateUint(b531CapabilityStateArguments[b531CapabilityStateKey], 32)
			if b531CapabilityStateErr != nil {
				return nil, b531CapabilityStateErr
			}
			b531CapabilityStateNumbers[b531CapabilityStateKey] = b531CapabilityStateNumber
		}
		b531CapabilityStateUpdated := b531CapabilityStateNumbers["prior"]
		var b531CapabilityStateOutcome string
		if b531CapabilityStateNumbers["registered"] == 0 {
			b531CapabilityStateOutcome = "skipped"
		} else {
			switch b531CapabilityStateNumbers["result"] {
			case 0:
				b531CapabilityStateUpdated = 0
				b531CapabilityStateOutcome = "cleared"
			case 1:
				if b531CapabilityStateUpdated == 0 {
					b531CapabilityStateUpdated = 1
				}
				b531CapabilityStateOutcome = "admitted"
			default:
				b531CapabilityStateOutcome = "Unknown"
			}
		}
		b531CapabilityStateOutput["capability"] = map[string]any{"state": "Hypothesis", "calculation": "Proven", "marker": 65535, "updated": b531CapabilityStateUpdated, "cached_bit0": b531CapabilityStateUpdated&1 != 0, "whole_nonzero": b531CapabilityStateUpdated != 0, "outcome": b531CapabilityStateOutcome}
	} else {
		if b531CapabilityStateRows == nil {
			return nil, fmt.Errorf("missing rows")
		}
		b531CapabilityStateIndex := 9
		for b531CapabilityStatePosition, b531CapabilityStateRow := range b531CapabilityStateRows {
			b531CapabilityStateObject := b531CapabilityStateRow.(map[string]any)
			b531CapabilityStateMatch := true
			for _, b531CapabilityStateKey := range []string{"a", "b", "c"} {
				b531CapabilityStateNumber, _ := b531CapabilityStateUint(b531CapabilityStateObject[b531CapabilityStateKey], 32)
				b531CapabilityStateMatch = b531CapabilityStateMatch && b531CapabilityStateNumber != 0
			}
			if b531CapabilityStateMatch {
				b531CapabilityStateIndex = b531CapabilityStatePosition
				break
			}
		}
		b531CapabilityStateOutput["select"] = map[string]any{"state": "Hypothesis", "calculation": "Proven", "index": b531CapabilityStateIndex}
	}
	b531CapabilityStateOutput["state"] = b531CapabilityStateState
	return json.Marshal(b531CapabilityStateOutput)
}
