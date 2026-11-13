package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func b555WeeklyRefreshRead(b555WeeklyRefreshDecoder *json.Decoder) (any, error) {
	b555WeeklyRefreshToken, b555WeeklyRefreshErr := b555WeeklyRefreshDecoder.Token()
	if b555WeeklyRefreshErr != nil {
		return nil, b555WeeklyRefreshErr
	}
	if b555WeeklyRefreshDelimiter, b555WeeklyRefreshOK := b555WeeklyRefreshToken.(json.Delim); b555WeeklyRefreshOK {
		switch b555WeeklyRefreshDelimiter {
		case '{':
			b555WeeklyRefreshObject := map[string]any{}
			for b555WeeklyRefreshDecoder.More() {
				b555WeeklyRefreshKeyToken, b555WeeklyRefreshErr := b555WeeklyRefreshDecoder.Token()
				if b555WeeklyRefreshErr != nil {
					return nil, b555WeeklyRefreshErr
				}
				b555WeeklyRefreshKey := b555WeeklyRefreshKeyToken.(string)
				if _, b555WeeklyRefreshFound := b555WeeklyRefreshObject[b555WeeklyRefreshKey]; b555WeeklyRefreshFound {
					return nil, fmt.Errorf("duplicate key %q", b555WeeklyRefreshKey)
				}
				b555WeeklyRefreshValue, b555WeeklyRefreshErr := b555WeeklyRefreshRead(b555WeeklyRefreshDecoder)
				if b555WeeklyRefreshErr != nil {
					return nil, b555WeeklyRefreshErr
				}
				b555WeeklyRefreshObject[b555WeeklyRefreshKey] = b555WeeklyRefreshValue
			}
			_, b555WeeklyRefreshErr := b555WeeklyRefreshDecoder.Token()
			return b555WeeklyRefreshObject, b555WeeklyRefreshErr
		case '[':
			b555WeeklyRefreshArray := []any{}
			for b555WeeklyRefreshDecoder.More() {
				b555WeeklyRefreshValue, b555WeeklyRefreshErr := b555WeeklyRefreshRead(b555WeeklyRefreshDecoder)
				if b555WeeklyRefreshErr != nil {
					return nil, b555WeeklyRefreshErr
				}
				b555WeeklyRefreshArray = append(b555WeeklyRefreshArray, b555WeeklyRefreshValue)
			}
			_, b555WeeklyRefreshErr := b555WeeklyRefreshDecoder.Token()
			return b555WeeklyRefreshArray, b555WeeklyRefreshErr
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	return b555WeeklyRefreshToken, nil
}

func b555WeeklyRefreshUint(b555WeeklyRefreshValue any, b555WeeklyRefreshMax uint64) (uint64, error) {
	b555WeeklyRefreshNumber, b555WeeklyRefreshOK := b555WeeklyRefreshValue.(json.Number)
	if !b555WeeklyRefreshOK {
		return 0, fmt.Errorf("expected unsigned integer")
	}
	b555WeeklyRefreshResult, b555WeeklyRefreshErr := strconv.ParseUint(string(b555WeeklyRefreshNumber), 10, 64)
	if b555WeeklyRefreshErr != nil || b555WeeklyRefreshResult > b555WeeklyRefreshMax {
		return 0, fmt.Errorf("unsigned integer out of range or nonintegral")
	}
	return b555WeeklyRefreshResult, nil
}

func b555WeeklyRefreshBytes(b555WeeklyRefreshValue any) error {
	b555WeeklyRefreshArray, b555WeeklyRefreshOK := b555WeeklyRefreshValue.([]any)
	if !b555WeeklyRefreshOK {
		return fmt.Errorf("expected byte array")
	}
	for _, b555WeeklyRefreshByte := range b555WeeklyRefreshArray {
		if _, b555WeeklyRefreshErr := b555WeeklyRefreshUint(b555WeeklyRefreshByte, 255); b555WeeklyRefreshErr != nil {
			return b555WeeklyRefreshErr
		}
	}
	return nil
}

func B555WeeklyRefreshEvaluate(b555WeeklyRefreshInput json.RawMessage) (json.RawMessage, error) {
	b555WeeklyRefreshDecoder := json.NewDecoder(bytes.NewReader(b555WeeklyRefreshInput))
	b555WeeklyRefreshDecoder.UseNumber()
	b555WeeklyRefreshValue, b555WeeklyRefreshErr := b555WeeklyRefreshRead(b555WeeklyRefreshDecoder)
	if b555WeeklyRefreshErr != nil {
		return nil, b555WeeklyRefreshErr
	}
	if _, b555WeeklyRefreshErr = b555WeeklyRefreshDecoder.Token(); b555WeeklyRefreshErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b555WeeklyRefreshObject, b555WeeklyRefreshOK := b555WeeklyRefreshValue.(map[string]any)
	if !b555WeeklyRefreshOK {
		return nil, fmt.Errorf("expected object")
	}
	b555WeeklyRefreshContext, b555WeeklyRefreshOK := b555WeeklyRefreshObject["context"].(map[string]any)
	if !b555WeeklyRefreshOK {
		return nil, fmt.Errorf("expected context object")
	}
	b555WeeklyRefreshRaw, b555WeeklyRefreshOK := b555WeeklyRefreshObject["fields"].(map[string]any)
	if !b555WeeklyRefreshOK {
		return nil, fmt.Errorf("expected fields object")
	}
	b555WeeklyRefreshFields := map[string]any{}
	for b555WeeklyRefreshKey, b555WeeklyRefreshRawValue := range b555WeeklyRefreshRaw {
		if b555WeeklyRefreshErr := b555WeeklyRefreshBytes(b555WeeklyRefreshRawValue); b555WeeklyRefreshErr != nil {
			return nil, b555WeeklyRefreshErr
		}
		b555WeeklyRefreshFields[b555WeeklyRefreshKey] = map[string]any{"raw": b555WeeklyRefreshRawValue, "state": "Unknown"}
	}
	b555WeeklyRefreshOutput := map[string]any{"context": b555WeeklyRefreshContext, "raw": b555WeeklyRefreshRaw, "state": "Unknown", "fields": b555WeeklyRefreshFields, "physical": "Unknown", "qualified": false}
	b555WeeklyRefreshArguments := map[string]any(nil)
	if b555WeeklyRefreshArgumentValue, b555WeeklyRefreshPresent := b555WeeklyRefreshObject["arguments"]; b555WeeklyRefreshPresent {
		b555WeeklyRefreshArguments, b555WeeklyRefreshOK = b555WeeklyRefreshArgumentValue.(map[string]any)
		if !b555WeeklyRefreshOK {
			return nil, fmt.Errorf("expected arguments object")
		}
		b555WeeklyRefreshOutput["arguments"] = b555WeeklyRefreshArguments
	}
	for _, b555WeeklyRefreshKey := range []string{"primary", "secondary", "operation"} {
		if b555WeeklyRefreshValue, b555WeeklyRefreshPresent := b555WeeklyRefreshContext[b555WeeklyRefreshKey]; b555WeeklyRefreshPresent && b555WeeklyRefreshValue != nil {
			if _, b555WeeklyRefreshErr := b555WeeklyRefreshUint(b555WeeklyRefreshValue, 255); b555WeeklyRefreshErr != nil {
				return nil, b555WeeklyRefreshErr
			}
		}
	}
	for _, b555WeeklyRefreshKey := range []string{"direction", "view"} {
		if b555WeeklyRefreshValue, b555WeeklyRefreshPresent := b555WeeklyRefreshContext[b555WeeklyRefreshKey]; b555WeeklyRefreshPresent && b555WeeklyRefreshValue != nil {
			if _, b555WeeklyRefreshOK := b555WeeklyRefreshValue.(string); !b555WeeklyRefreshOK {
				return nil, fmt.Errorf("expected context string")
			}
		}
	}
	if b555WeeklyRefreshContext["primary"] != json.Number("181") || b555WeeklyRefreshContext["secondary"] != json.Number("85") || b555WeeklyRefreshContext["direction"] != "local" {
		return json.Marshal(b555WeeklyRefreshOutput)
	}
	b555WeeklyRefreshWeekly := b555WeeklyRefreshContext["operation"] == json.Number("164") && b555WeeklyRefreshContext["view"] == "weekly"
	b555WeeklyRefreshShape := b555WeeklyRefreshContext["operation"] == json.Number("166") && b555WeeklyRefreshContext["view"] == "nonheating-shape"
	if !b555WeeklyRefreshWeekly && !b555WeeklyRefreshShape {
		return json.Marshal(b555WeeklyRefreshOutput)
	}
	if b555WeeklyRefreshArguments == nil {
		return nil, fmt.Errorf("missing arguments")
	}
	b555WeeklyRefreshResult := map[string]any{"raw": b555WeeklyRefreshArguments, "state": "Hypothesis", "complete": false}
	if b555WeeklyRefreshWeekly {
		b555WeeklyRefreshStatus, b555WeeklyRefreshOK := b555WeeklyRefreshArguments["status_ok"].(bool)
		if !b555WeeklyRefreshOK {
			return nil, fmt.Errorf("expected status_ok Boolean")
		}
		b555WeeklyRefreshLimit, b555WeeklyRefreshErr := b555WeeklyRefreshUint(b555WeeklyRefreshArguments["slot_limit"], 255)
		if b555WeeklyRefreshErr != nil {
			return nil, b555WeeklyRefreshErr
		}
		b555WeeklyRefreshDays, b555WeeklyRefreshOK := b555WeeklyRefreshArguments["days"].([]any)
		if !b555WeeklyRefreshOK || len(b555WeeklyRefreshDays) != 7 {
			return nil, fmt.Errorf("expected seven days")
		}
		for _, b555WeeklyRefreshDayValue := range b555WeeklyRefreshDays {
			b555WeeklyRefreshDay, b555WeeklyRefreshOK := b555WeeklyRefreshDayValue.(map[string]any)
			if !b555WeeklyRefreshOK {
				return nil, fmt.Errorf("expected day object")
			}
			if _, b555WeeklyRefreshErr := b555WeeklyRefreshUint(b555WeeklyRefreshDay["count"], 255); b555WeeklyRefreshErr != nil {
				return nil, b555WeeklyRefreshErr
			}
			b555WeeklyRefreshSlots, b555WeeklyRefreshOK := b555WeeklyRefreshDay["slots"].([]any)
			if !b555WeeklyRefreshOK {
				return nil, fmt.Errorf("expected slots array")
			}
			for _, b555WeeklyRefreshSlotValue := range b555WeeklyRefreshSlots {
				b555WeeklyRefreshSlot, b555WeeklyRefreshOK := b555WeeklyRefreshSlotValue.(map[string]any)
				if !b555WeeklyRefreshOK {
					return nil, fmt.Errorf("expected slot object")
				}
				if _, b555WeeklyRefreshOK := b555WeeklyRefreshSlot["ok"].(bool); !b555WeeklyRefreshOK {
					return nil, fmt.Errorf("expected slot Boolean")
				}
				if b555WeeklyRefreshRaw, b555WeeklyRefreshPresent := b555WeeklyRefreshSlot["raw"]; b555WeeklyRefreshPresent {
					if b555WeeklyRefreshErr := b555WeeklyRefreshBytes(b555WeeklyRefreshRaw); b555WeeklyRefreshErr != nil {
						return nil, b555WeeklyRefreshErr
					}
				}
			}
		}
		b555WeeklyRefreshAttempts, b555WeeklyRefreshCompleted := 0, 0
		b555WeeklyRefreshProgress := []any{}
		b555WeeklyRefreshReason := "status-failure"
		if b555WeeklyRefreshStatus {
			b555WeeklyRefreshReason = "complete"
			for b555WeeklyRefreshIndex, b555WeeklyRefreshDayValue := range b555WeeklyRefreshDays {
				b555WeeklyRefreshDay := b555WeeklyRefreshDayValue.(map[string]any)
				b555WeeklyRefreshCount, _ := b555WeeklyRefreshUint(b555WeeklyRefreshDay["count"], 255)
				if b555WeeklyRefreshCount > b555WeeklyRefreshLimit {
					b555WeeklyRefreshReason = "count-overflow"
					break
				}
				b555WeeklyRefreshSlots := b555WeeklyRefreshDay["slots"].([]any)
				b555WeeklyRefreshSuccess := []any{}
				for b555WeeklyRefreshSlotIndex := 0; b555WeeklyRefreshSlotIndex < int(b555WeeklyRefreshCount); b555WeeklyRefreshSlotIndex++ {
					if b555WeeklyRefreshSlotIndex >= len(b555WeeklyRefreshSlots) {
						b555WeeklyRefreshReason = "missing-slot"
						break
					}
					b555WeeklyRefreshAttempts++
					b555WeeklyRefreshSlot := b555WeeklyRefreshSlots[b555WeeklyRefreshSlotIndex].(map[string]any)
					if !b555WeeklyRefreshSlot["ok"].(bool) {
						b555WeeklyRefreshReason = "slot-failure"
						break
					}
					if _, b555WeeklyRefreshPresent := b555WeeklyRefreshSlot["raw"]; !b555WeeklyRefreshPresent {
						b555WeeklyRefreshReason = "missing-slot-data"
						break
					}
					b555WeeklyRefreshSuccess = append(b555WeeklyRefreshSuccess, b555WeeklyRefreshSlot)
				}
				b555WeeklyRefreshDayComplete := len(b555WeeklyRefreshSuccess) == int(b555WeeklyRefreshCount)
				b555WeeklyRefreshProgress = append(b555WeeklyRefreshProgress, map[string]any{"day": b555WeeklyRefreshIndex, "slots": b555WeeklyRefreshSuccess, "complete": b555WeeklyRefreshDayComplete})
				if !b555WeeklyRefreshDayComplete {
					break
				}
				b555WeeklyRefreshCompleted++
			}
		}
		b555WeeklyRefreshResult["attempts"] = b555WeeklyRefreshAttempts
		b555WeeklyRefreshResult["completed_days"] = b555WeeklyRefreshCompleted
		b555WeeklyRefreshResult["days"] = b555WeeklyRefreshProgress
		b555WeeklyRefreshResult["reason"] = b555WeeklyRefreshReason
		b555WeeklyRefreshResult["complete"] = b555WeeklyRefreshCompleted == 7
		if b555WeeklyRefreshRawValue, b555WeeklyRefreshPresent := b555WeeklyRefreshRaw["weekly"]; b555WeeklyRefreshPresent {
			b555WeeklyRefreshResult["raw"] = b555WeeklyRefreshRawValue
		}
		b555WeeklyRefreshFields["weekly"] = b555WeeklyRefreshResult
	} else {
		b555WeeklyRefreshInitial, b555WeeklyRefreshOK := b555WeeklyRefreshArguments["initial"].([]any)
		if !b555WeeklyRefreshOK || len(b555WeeklyRefreshInitial) != 5 {
			return nil, fmt.Errorf("expected five initial scalars")
		}
		for _, b555WeeklyRefreshScalar := range b555WeeklyRefreshInitial {
			switch b555WeeklyRefreshScalar.(type) {
			case []any, map[string]any:
				return nil, fmt.Errorf("expected scalar")
			}
		}
		b555WeeklyRefreshCount, b555WeeklyRefreshErr := b555WeeklyRefreshUint(b555WeeklyRefreshArguments["count"], 255)
		if b555WeeklyRefreshErr != nil {
			return nil, b555WeeklyRefreshErr
		}
		b555WeeklyRefreshSetpoints, b555WeeklyRefreshOK := b555WeeklyRefreshArguments["slot_setpoints"].([]any)
		if !b555WeeklyRefreshOK {
			return nil, fmt.Errorf("expected setpoints array")
		}
		for _, b555WeeklyRefreshSetpoint := range b555WeeklyRefreshSetpoints {
			if _, b555WeeklyRefreshErr := b555WeeklyRefreshUint(b555WeeklyRefreshSetpoint, 65535); b555WeeklyRefreshErr != nil {
				return nil, b555WeeklyRefreshErr
			}
		}
		if len(b555WeeklyRefreshSetpoints) < int(b555WeeklyRefreshCount) {
			return nil, fmt.Errorf("insufficient setpoints")
		}
		b555WeeklyRefreshProcessed := []any{}
		for b555WeeklyRefreshIndex := 0; b555WeeklyRefreshIndex < int(b555WeeklyRefreshCount); b555WeeklyRefreshIndex++ {
			b555WeeklyRefreshProcessed = append(b555WeeklyRefreshProcessed, 65535)
		}
		b555WeeklyRefreshResult["initial"] = b555WeeklyRefreshInitial
		b555WeeklyRefreshResult["processed_setpoints"] = b555WeeklyRefreshProcessed
		b555WeeklyRefreshResult["processed_slots"] = b555WeeklyRefreshCount
		b555WeeklyRefreshResult["unused_setpoints"] = b555WeeklyRefreshSetpoints[int(b555WeeklyRefreshCount):]
		b555WeeklyRefreshResult["complete"] = true
		if b555WeeklyRefreshRawValue, b555WeeklyRefreshPresent := b555WeeklyRefreshRaw["nonheating-shape"]; b555WeeklyRefreshPresent {
			b555WeeklyRefreshResult["raw"] = b555WeeklyRefreshRawValue
		}
		b555WeeklyRefreshFields["nonheating-shape"] = b555WeeklyRefreshResult
	}
	b555WeeklyRefreshOutput["state"] = "partial"
	if b555WeeklyRefreshResult["complete"] == true {
		b555WeeklyRefreshOutput["state"] = "complete"
	}
	return json.Marshal(b555WeeklyRefreshOutput)
}
