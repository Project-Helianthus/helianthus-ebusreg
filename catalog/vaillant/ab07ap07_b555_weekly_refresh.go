package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func ab07ap07Read(ab07ap07Decoder *json.Decoder) (any, error) {
	ab07ap07Token, ab07ap07Err := ab07ap07Decoder.Token()
	if ab07ap07Err != nil {
		return nil, ab07ap07Err
	}
	if ab07ap07Delimiter, ab07ap07OK := ab07ap07Token.(json.Delim); ab07ap07OK {
		switch ab07ap07Delimiter {
		case '{':
			ab07ap07Object := map[string]any{}
			for ab07ap07Decoder.More() {
				ab07ap07KeyToken, ab07ap07Err := ab07ap07Decoder.Token()
				if ab07ap07Err != nil {
					return nil, ab07ap07Err
				}
				ab07ap07Key := ab07ap07KeyToken.(string)
				if _, ab07ap07Found := ab07ap07Object[ab07ap07Key]; ab07ap07Found {
					return nil, fmt.Errorf("duplicate key %q", ab07ap07Key)
				}
				ab07ap07Value, ab07ap07Err := ab07ap07Read(ab07ap07Decoder)
				if ab07ap07Err != nil {
					return nil, ab07ap07Err
				}
				ab07ap07Object[ab07ap07Key] = ab07ap07Value
			}
			_, ab07ap07Err := ab07ap07Decoder.Token()
			return ab07ap07Object, ab07ap07Err
		case '[':
			ab07ap07Array := []any{}
			for ab07ap07Decoder.More() {
				ab07ap07Value, ab07ap07Err := ab07ap07Read(ab07ap07Decoder)
				if ab07ap07Err != nil {
					return nil, ab07ap07Err
				}
				ab07ap07Array = append(ab07ap07Array, ab07ap07Value)
			}
			_, ab07ap07Err := ab07ap07Decoder.Token()
			return ab07ap07Array, ab07ap07Err
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	return ab07ap07Token, nil
}

func ab07ap07Uint(ab07ap07Value any, ab07ap07Max uint64) (uint64, error) {
	ab07ap07Number, ab07ap07OK := ab07ap07Value.(json.Number)
	if !ab07ap07OK {
		return 0, fmt.Errorf("expected unsigned integer")
	}
	ab07ap07Result, ab07ap07Err := strconv.ParseUint(string(ab07ap07Number), 10, 64)
	if ab07ap07Err != nil || ab07ap07Result > ab07ap07Max {
		return 0, fmt.Errorf("unsigned integer out of range or nonintegral")
	}
	return ab07ap07Result, nil
}

func ab07ap07Bytes(ab07ap07Value any) error {
	ab07ap07Array, ab07ap07OK := ab07ap07Value.([]any)
	if !ab07ap07OK {
		return fmt.Errorf("expected byte array")
	}
	for _, ab07ap07Byte := range ab07ap07Array {
		if _, ab07ap07Err := ab07ap07Uint(ab07ap07Byte, 255); ab07ap07Err != nil {
			return ab07ap07Err
		}
	}
	return nil
}

func Ab07Ap07Evaluate(ab07ap07Input json.RawMessage) (json.RawMessage, error) {
	ab07ap07Decoder := json.NewDecoder(bytes.NewReader(ab07ap07Input))
	ab07ap07Decoder.UseNumber()
	ab07ap07Value, ab07ap07Err := ab07ap07Read(ab07ap07Decoder)
	if ab07ap07Err != nil {
		return nil, ab07ap07Err
	}
	if _, ab07ap07Err = ab07ap07Decoder.Token(); ab07ap07Err != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab07ap07Object, ab07ap07OK := ab07ap07Value.(map[string]any)
	if !ab07ap07OK {
		return nil, fmt.Errorf("expected object")
	}
	ab07ap07Context, ab07ap07OK := ab07ap07Object["context"].(map[string]any)
	if !ab07ap07OK {
		return nil, fmt.Errorf("expected context object")
	}
	ab07ap07Raw, ab07ap07OK := ab07ap07Object["fields"].(map[string]any)
	if !ab07ap07OK {
		return nil, fmt.Errorf("expected fields object")
	}
	ab07ap07Fields := map[string]any{}
	for ab07ap07Key, ab07ap07RawValue := range ab07ap07Raw {
		if ab07ap07Err := ab07ap07Bytes(ab07ap07RawValue); ab07ap07Err != nil {
			return nil, ab07ap07Err
		}
		ab07ap07Fields[ab07ap07Key] = map[string]any{"raw": ab07ap07RawValue, "state": "Unknown"}
	}
	ab07ap07Output := map[string]any{"context": ab07ap07Context, "raw": ab07ap07Raw, "state": "Unknown", "fields": ab07ap07Fields, "physical": "Unknown", "qualified": false}
	ab07ap07Arguments := map[string]any(nil)
	if ab07ap07ArgumentValue, ab07ap07Present := ab07ap07Object["arguments"]; ab07ap07Present {
		ab07ap07Arguments, ab07ap07OK = ab07ap07ArgumentValue.(map[string]any)
		if !ab07ap07OK {
			return nil, fmt.Errorf("expected arguments object")
		}
		ab07ap07Output["arguments"] = ab07ap07Arguments
	}
	for _, ab07ap07Key := range []string{"primary", "secondary", "operation"} {
		if ab07ap07Value, ab07ap07Present := ab07ap07Context[ab07ap07Key]; ab07ap07Present && ab07ap07Value != nil {
			if _, ab07ap07Err := ab07ap07Uint(ab07ap07Value, 255); ab07ap07Err != nil {
				return nil, ab07ap07Err
			}
		}
	}
	for _, ab07ap07Key := range []string{"direction", "view"} {
		if ab07ap07Value, ab07ap07Present := ab07ap07Context[ab07ap07Key]; ab07ap07Present && ab07ap07Value != nil {
			if _, ab07ap07OK := ab07ap07Value.(string); !ab07ap07OK {
				return nil, fmt.Errorf("expected context string")
			}
		}
	}
	if ab07ap07Context["primary"] != json.Number("181") || ab07ap07Context["secondary"] != json.Number("85") || ab07ap07Context["direction"] != "local" {
		return json.Marshal(ab07ap07Output)
	}
	ab07ap07Weekly := ab07ap07Context["operation"] == json.Number("164") && ab07ap07Context["view"] == "weekly"
	ab07ap07Shape := ab07ap07Context["operation"] == json.Number("166") && ab07ap07Context["view"] == "nonheating-shape"
	if !ab07ap07Weekly && !ab07ap07Shape {
		return json.Marshal(ab07ap07Output)
	}
	if ab07ap07Arguments == nil {
		return nil, fmt.Errorf("missing arguments")
	}
	ab07ap07Result := map[string]any{"raw": ab07ap07Arguments, "state": "Hypothesis", "complete": false}
	if ab07ap07Weekly {
		ab07ap07Status, ab07ap07OK := ab07ap07Arguments["status_ok"].(bool)
		if !ab07ap07OK {
			return nil, fmt.Errorf("expected status_ok Boolean")
		}
		ab07ap07Limit, ab07ap07Err := ab07ap07Uint(ab07ap07Arguments["slot_limit"], 255)
		if ab07ap07Err != nil {
			return nil, ab07ap07Err
		}
		ab07ap07Days, ab07ap07OK := ab07ap07Arguments["days"].([]any)
		if !ab07ap07OK || len(ab07ap07Days) != 7 {
			return nil, fmt.Errorf("expected seven days")
		}
		for _, ab07ap07DayValue := range ab07ap07Days {
			ab07ap07Day, ab07ap07OK := ab07ap07DayValue.(map[string]any)
			if !ab07ap07OK {
				return nil, fmt.Errorf("expected day object")
			}
			if _, ab07ap07Err := ab07ap07Uint(ab07ap07Day["count"], 255); ab07ap07Err != nil {
				return nil, ab07ap07Err
			}
			ab07ap07Slots, ab07ap07OK := ab07ap07Day["slots"].([]any)
			if !ab07ap07OK {
				return nil, fmt.Errorf("expected slots array")
			}
			for _, ab07ap07SlotValue := range ab07ap07Slots {
				ab07ap07Slot, ab07ap07OK := ab07ap07SlotValue.(map[string]any)
				if !ab07ap07OK {
					return nil, fmt.Errorf("expected slot object")
				}
				if _, ab07ap07OK := ab07ap07Slot["ok"].(bool); !ab07ap07OK {
					return nil, fmt.Errorf("expected slot Boolean")
				}
				if ab07ap07Raw, ab07ap07Present := ab07ap07Slot["raw"]; ab07ap07Present {
					if ab07ap07Err := ab07ap07Bytes(ab07ap07Raw); ab07ap07Err != nil {
						return nil, ab07ap07Err
					}
				}
			}
		}
		ab07ap07Attempts, ab07ap07Completed := 0, 0
		ab07ap07Progress := []any{}
		ab07ap07Reason := "status-failure"
		if ab07ap07Status {
			ab07ap07Reason = "complete"
			for ab07ap07Index, ab07ap07DayValue := range ab07ap07Days {
				ab07ap07Day := ab07ap07DayValue.(map[string]any)
				ab07ap07Count, _ := ab07ap07Uint(ab07ap07Day["count"], 255)
				if ab07ap07Count > ab07ap07Limit {
					ab07ap07Reason = "count-overflow"
					break
				}
				ab07ap07Slots := ab07ap07Day["slots"].([]any)
				ab07ap07Success := []any{}
				for ab07ap07SlotIndex := 0; ab07ap07SlotIndex < int(ab07ap07Count); ab07ap07SlotIndex++ {
					if ab07ap07SlotIndex >= len(ab07ap07Slots) {
						ab07ap07Reason = "missing-slot"
						break
					}
					ab07ap07Attempts++
					ab07ap07Slot := ab07ap07Slots[ab07ap07SlotIndex].(map[string]any)
					if !ab07ap07Slot["ok"].(bool) {
						ab07ap07Reason = "slot-failure"
						break
					}
					if _, ab07ap07Present := ab07ap07Slot["raw"]; !ab07ap07Present {
						ab07ap07Reason = "missing-slot-data"
						break
					}
					ab07ap07Success = append(ab07ap07Success, ab07ap07Slot)
				}
				ab07ap07DayComplete := len(ab07ap07Success) == int(ab07ap07Count)
				ab07ap07Progress = append(ab07ap07Progress, map[string]any{"day": ab07ap07Index, "slots": ab07ap07Success, "complete": ab07ap07DayComplete})
				if !ab07ap07DayComplete {
					break
				}
				ab07ap07Completed++
			}
		}
		ab07ap07Result["attempts"] = ab07ap07Attempts
		ab07ap07Result["completed_days"] = ab07ap07Completed
		ab07ap07Result["days"] = ab07ap07Progress
		ab07ap07Result["reason"] = ab07ap07Reason
		ab07ap07Result["complete"] = ab07ap07Completed == 7
		if ab07ap07RawValue, ab07ap07Present := ab07ap07Raw["weekly"]; ab07ap07Present {
			ab07ap07Result["raw"] = ab07ap07RawValue
		}
		ab07ap07Fields["weekly"] = ab07ap07Result
	} else {
		ab07ap07Initial, ab07ap07OK := ab07ap07Arguments["initial"].([]any)
		if !ab07ap07OK || len(ab07ap07Initial) != 5 {
			return nil, fmt.Errorf("expected five initial scalars")
		}
		for _, ab07ap07Scalar := range ab07ap07Initial {
			switch ab07ap07Scalar.(type) {
			case []any, map[string]any:
				return nil, fmt.Errorf("expected scalar")
			}
		}
		ab07ap07Count, ab07ap07Err := ab07ap07Uint(ab07ap07Arguments["count"], 255)
		if ab07ap07Err != nil {
			return nil, ab07ap07Err
		}
		ab07ap07Setpoints, ab07ap07OK := ab07ap07Arguments["slot_setpoints"].([]any)
		if !ab07ap07OK {
			return nil, fmt.Errorf("expected setpoints array")
		}
		for _, ab07ap07Setpoint := range ab07ap07Setpoints {
			if _, ab07ap07Err := ab07ap07Uint(ab07ap07Setpoint, 65535); ab07ap07Err != nil {
				return nil, ab07ap07Err
			}
		}
		if len(ab07ap07Setpoints) < int(ab07ap07Count) {
			return nil, fmt.Errorf("insufficient setpoints")
		}
		ab07ap07Processed := []any{}
		for ab07ap07Index := 0; ab07ap07Index < int(ab07ap07Count); ab07ap07Index++ {
			ab07ap07Processed = append(ab07ap07Processed, 65535)
		}
		ab07ap07Result["initial"] = ab07ap07Initial
		ab07ap07Result["processed_setpoints"] = ab07ap07Processed
		ab07ap07Result["processed_slots"] = ab07ap07Count
		ab07ap07Result["unused_setpoints"] = ab07ap07Setpoints[int(ab07ap07Count):]
		ab07ap07Result["complete"] = true
		if ab07ap07RawValue, ab07ap07Present := ab07ap07Raw["nonheating-shape"]; ab07ap07Present {
			ab07ap07Result["raw"] = ab07ap07RawValue
		}
		ab07ap07Fields["nonheating-shape"] = ab07ap07Result
	}
	ab07ap07Output["state"] = "partial"
	if ab07ap07Result["complete"] == true {
		ab07ap07Output["state"] = "complete"
	}
	return json.Marshal(ab07ap07Output)
}
