package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func ab07ap07Week() map[string]any {
	ab07ap07Days := []any{}
	for ab07ap07Index := 0; ab07ap07Index < 7; ab07ap07Index++ {
		ab07ap07Days = append(ab07ap07Days, map[string]any{"count": 0, "slots": []any{}})
	}
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 85, "operation": 164, "direction": "local", "view": "weekly", "family": "synthetic-family", "variant": "synthetic-variant"}, "fields": map[string]any{"evidence": []any{0, 255}}, "arguments": map[string]any{"status_ok": true, "slot_limit": 2, "days": ab07ap07Days}}
}
func ab07ap07Day(ab07ap07Input map[string]any, ab07ap07Index int, ab07ap07Count int, ab07ap07Slots ...any) {
	if ab07ap07Slots == nil {
		ab07ap07Slots = []any{}
	}
	ab07ap07Input["arguments"].(map[string]any)["days"].([]any)[ab07ap07Index] = map[string]any{"count": ab07ap07Count, "slots": ab07ap07Slots}
}
func ab07ap07Slot(ab07ap07OK bool, ab07ap07Byte int) any {
	return map[string]any{"ok": ab07ap07OK, "raw": []any{ab07ap07Byte}}
}
func ab07ap07Run(ab07ap07T *testing.T, ab07ap07Input map[string]any) map[string]any {
	ab07ap07T.Helper()
	ab07ap07JSON, ab07ap07Err := json.Marshal(ab07ap07Input)
	if ab07ap07Err != nil {
		ab07ap07T.Fatal(ab07ap07Err)
	}
	ab07ap07Output, ab07ap07Err := Ab07Ap07Evaluate(ab07ap07JSON)
	if ab07ap07Err != nil {
		ab07ap07T.Fatal(ab07ap07Err)
	}
	ab07ap07Result := map[string]any{}
	if ab07ap07Err := json.Unmarshal(ab07ap07Output, &ab07ap07Result); ab07ap07Err != nil {
		ab07ap07T.Fatal(ab07ap07Err)
	}
	if ab07ap07Result["physical"] != "Unknown" || ab07ap07Result["qualified"] != false {
		ab07ap07T.Fatalf("qualification: %s", ab07ap07Output)
	}
	return ab07ap07Result
}
func ab07ap07WeeklyResult(ab07ap07T *testing.T, ab07ap07Input map[string]any, ab07ap07Attempts, ab07ap07Days int, ab07ap07Complete bool, ab07ap07Reason string) map[string]any {
	ab07ap07T.Helper()
	ab07ap07Output := ab07ap07Run(ab07ap07T, ab07ap07Input)
	ab07ap07Result := ab07ap07Output["fields"].(map[string]any)["weekly"].(map[string]any)
	if ab07ap07Result["attempts"] != float64(ab07ap07Attempts) || ab07ap07Result["completed_days"] != float64(ab07ap07Days) || ab07ap07Result["complete"] != ab07ap07Complete || ab07ap07Result["reason"] != ab07ap07Reason || ab07ap07Result["state"] != "Hypothesis" {
		ab07ap07T.Fatalf("wrong reduction: %#v", ab07ap07Result)
	}
	ab07ap07State := "partial"
	if ab07ap07Complete {
		ab07ap07State = "complete"
	}
	if ab07ap07Output["state"] != ab07ap07State {
		ab07ap07T.Fatalf("wrong state: %#v", ab07ap07Output)
	}
	if ab07ap07Output["fields"].(map[string]any)["evidence"].(map[string]any)["state"] != "Unknown" {
		ab07ap07T.Fatal("invented raw interpretation")
	}
	return ab07ap07Result
}
func TestAb07Ap07WeeklyTransitions(ab07ap07T *testing.T) {
	ab07ap07Input := ab07ap07Week()
	ab07ap07WeeklyResult(ab07ap07T, ab07ap07Input, 0, 7, true, "complete")
	ab07ap07Day(ab07ap07Input, 0, 1, ab07ap07Slot(true, 11))
	ab07ap07Input["arguments"].(map[string]any)["status_ok"] = false
	ab07ap07WeeklyResult(ab07ap07T, ab07ap07Input, 0, 0, false, "status-failure")
	ab07ap07Input["arguments"].(map[string]any)["status_ok"] = true
	ab07ap07WeeklyResult(ab07ap07T, ab07ap07Input, 1, 7, true, "complete")
	ab07ap07Day(ab07ap07Input, 1, 1, ab07ap07Slot(true, 22))
	ab07ap07Day(ab07ap07Input, 2, 3)
	ab07ap07Result := ab07ap07WeeklyResult(ab07ap07T, ab07ap07Input, 2, 2, false, "count-overflow")
	if len(ab07ap07Result["days"].([]any)) != 2 {
		ab07ap07T.Fatal("lost completed days")
	}
	ab07ap07Input = ab07ap07Week()
	ab07ap07Day(ab07ap07Input, 0, 1, ab07ap07Slot(false, 33))
	ab07ap07Result = ab07ap07WeeklyResult(ab07ap07T, ab07ap07Input, 1, 0, false, "slot-failure")
	if len(ab07ap07Result["days"].([]any)[0].(map[string]any)["slots"].([]any)) != 0 {
		ab07ap07T.Fatal("failed slot succeeded")
	}
	ab07ap07Day(ab07ap07Input, 0, 1)
	ab07ap07WeeklyResult(ab07ap07T, ab07ap07Input, 0, 0, false, "missing-slot")
	ab07ap07Day(ab07ap07Input, 0, 1, map[string]any{"ok": true})
	ab07ap07WeeklyResult(ab07ap07T, ab07ap07Input, 1, 0, false, "missing-slot-data")
}
func TestAb07Ap07EveryDayAndSlot(ab07ap07T *testing.T) {
	for ab07ap07DayIndex := 0; ab07ap07DayIndex < 7; ab07ap07DayIndex++ {
		for ab07ap07SlotIndex := 0; ab07ap07SlotIndex < 2; ab07ap07SlotIndex++ {
			ab07ap07T.Run(string(rune('A'+ab07ap07DayIndex))+string(rune('0'+ab07ap07SlotIndex)), func(ab07ap07T *testing.T) {
				ab07ap07Input := ab07ap07Week()
				for ab07ap07Index := 0; ab07ap07Index < 7; ab07ap07Index++ {
					ab07ap07Day(ab07ap07Input, ab07ap07Index, 2, ab07ap07Slot(true, 10), ab07ap07Slot(true, 20))
				}
				ab07ap07WeeklyResult(ab07ap07T, ab07ap07Input, 14, 7, true, "complete")
				ab07ap07Input["arguments"].(map[string]any)["days"].([]any)[ab07ap07DayIndex].(map[string]any)["slots"].([]any)[ab07ap07SlotIndex].(map[string]any)["ok"] = false
				ab07ap07Result := ab07ap07WeeklyResult(ab07ap07T, ab07ap07Input, ab07ap07DayIndex*2+ab07ap07SlotIndex+1, ab07ap07DayIndex, false, "slot-failure")
				ab07ap07Slots := ab07ap07Result["days"].([]any)[ab07ap07DayIndex].(map[string]any)["slots"].([]any)
				if len(ab07ap07Slots) != ab07ap07SlotIndex {
					ab07ap07T.Fatal("lost successful prefix")
				}
				if ab07ap07SlotIndex == 1 && ab07ap07Slots[0].(map[string]any)["raw"].([]any)[0] != float64(10) {
					ab07ap07T.Fatal("raw slot altered")
				}
			})
		}
	}
}
func TestAb07Ap07Nonheating(ab07ap07T *testing.T) {
	for ab07ap07Count := 0; ab07ap07Count <= 2; ab07ap07Count++ {
		ab07ap07Input := ab07ap07Week()
		ab07ap07Input["context"].(map[string]any)["operation"] = 166
		ab07ap07Input["context"].(map[string]any)["view"] = "nonheating-shape"
		ab07ap07Input["arguments"] = map[string]any{"initial": []any{1, 2, 3, 4, 5}, "count": ab07ap07Count, "slot_setpoints": []any{100, 200, 300}}
		ab07ap07Result := ab07ap07Run(ab07ap07T, ab07ap07Input)["fields"].(map[string]any)["nonheating-shape"].(map[string]any)
		if !reflect.DeepEqual(ab07ap07Result["initial"], []any{float64(1), float64(2), float64(3), float64(4), float64(5)}) || ab07ap07Result["processed_slots"] != float64(ab07ap07Count) || ab07ap07Result["complete"] != true || ab07ap07Result["state"] != "Hypothesis" {
			ab07ap07T.Fatalf("shape: %#v", ab07ap07Result)
		}
		ab07ap07Processed := ab07ap07Result["processed_setpoints"].([]any)
		if len(ab07ap07Processed) != ab07ap07Count {
			ab07ap07T.Fatal("invented processed slots")
		}
		for _, ab07ap07Value := range ab07ap07Processed {
			if ab07ap07Value != float64(65535) {
				ab07ap07T.Fatal("missing FFFF")
			}
		}
		if !reflect.DeepEqual(ab07ap07Result["unused_setpoints"], []any{float64(100), float64(200), float64(300)}[ab07ap07Count:]) {
			ab07ap07T.Fatal("prefilled unused values")
		}
	}
}
func TestAb07Ap07ContextAndRaw(ab07ap07T *testing.T) {
	for _, ab07ap07Key := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
		ab07ap07T.Run(ab07ap07Key, func(ab07ap07T *testing.T) {
			for _, ab07ap07Value := range []any{nil, "different"} {
				ab07ap07Input := ab07ap07Week()
				ab07ap07Context := ab07ap07Input["context"].(map[string]any)
				if ab07ap07Key == "primary" || ab07ap07Key == "secondary" || ab07ap07Key == "operation" {
					ab07ap07Value = 0
				}
				ab07ap07Context[ab07ap07Key] = ab07ap07Value
				ab07ap07Output := ab07ap07Run(ab07ap07T, ab07ap07Input)
				if (ab07ap07Key == "family" || ab07ap07Key == "variant") && ab07ap07Output["state"] != "complete" {
					ab07ap07T.Fatal("opaque context gated local interpretation")
				}
				if ab07ap07Key != "family" && ab07ap07Key != "variant" && ab07ap07Output["state"] != "Unknown" {
					ab07ap07T.Fatal("selector gained interpretation")
				}
				ab07ap07Encoded, _ := json.Marshal(ab07ap07Context)
				ab07ap07Returned, _ := json.Marshal(ab07ap07Output["context"])
				if string(ab07ap07Encoded) != string(ab07ap07Returned) {
					ab07ap07T.Fatal("context changed")
				}
			}
			ab07ap07Input := ab07ap07Week()
			delete(ab07ap07Input["context"].(map[string]any), ab07ap07Key)
			ab07ap07Output := ab07ap07Run(ab07ap07T, ab07ap07Input)
			if _, ab07ap07Present := ab07ap07Output["context"].(map[string]any)[ab07ap07Key]; ab07ap07Present {
				ab07ap07T.Fatal("missing became present")
			}
		})
	}
	for _, ab07ap07Raw := range []map[string]any{{}, {"evidence": []any{}}, {"evidence": []any{0}}, {"evidence": []any{255}, "sibling": []any{1, 2, 3}}} {
		ab07ap07Input := ab07ap07Week()
		ab07ap07Input["fields"] = ab07ap07Raw
		ab07ap07Input["arguments"].(map[string]any)["extra"] = map[string]any{"nested": []any{nil, false, 0}}
		ab07ap07Output := ab07ap07Run(ab07ap07T, ab07ap07Input)
		ab07ap07Expected, _ := json.Marshal(ab07ap07Raw)
		ab07ap07Actual, _ := json.Marshal(ab07ap07Output["raw"])
		if string(ab07ap07Expected) != string(ab07ap07Actual) {
			ab07ap07T.Fatal("raw changed")
		}
		ab07ap07Expected, _ = json.Marshal(ab07ap07Input["arguments"])
		ab07ap07Actual, _ = json.Marshal(ab07ap07Output["arguments"])
		if string(ab07ap07Expected) != string(ab07ap07Actual) {
			ab07ap07T.Fatal("arguments changed")
		}
	}
}
func TestAb07Ap07Rejections(ab07ap07T *testing.T) {
	ab07ap07JSON, _ := json.Marshal(ab07ap07Week())
	ab07ap07Base := string(ab07ap07JSON)
	for _, ab07ap07Invalid := range []string{
		"", "null", "[]", "{", ab07ap07Base + " {}",
		strings.Replace(ab07ap07Base, `"primary":181`, `"primary":181,"primary":181`, 1),
		strings.Replace(ab07ap07Base, `"count":0`, `"count":0,"count":0`, 1),
		strings.Replace(ab07ap07Base, `"evidence":[0,255]`, `"evidence":[256]`, 1),
		strings.Replace(ab07ap07Base, `"evidence":[0,255]`, `"evidence":[true]`, 1),
		strings.Replace(ab07ap07Base, `"evidence":[0,255]`, `"evidence":[1.5]`, 1),
		strings.Replace(ab07ap07Base, `"slot_limit":2`, `"slot_limit":null`, 1),
		strings.Replace(ab07ap07Base, `"status_ok":true`, `"status_ok":0`, 1),
		strings.Replace(ab07ap07Base, `"primary":181`, `"primary":"181"`, 1),
	} {
		ab07ap07Output, ab07ap07Err := Ab07Ap07Evaluate(json.RawMessage(ab07ap07Invalid))
		if ab07ap07Err == nil || ab07ap07Output != nil {
			ab07ap07T.Fatalf("accepted invalid: %s", ab07ap07Invalid)
		}
	}
	for _, ab07ap07Key := range []string{"status_ok", "slot_limit", "days"} {
		ab07ap07Input := ab07ap07Week()
		delete(ab07ap07Input["arguments"].(map[string]any), ab07ap07Key)
		ab07ap07JSON, _ := json.Marshal(ab07ap07Input)
		if _, ab07ap07Err := Ab07Ap07Evaluate(ab07ap07JSON); ab07ap07Err == nil {
			ab07ap07T.Fatalf("missing %s accepted", ab07ap07Key)
		}
	}
	for _, ab07ap07Invalid := range []any{nil, true, "1", -1, 256, 1.5, []any{}} {
		for _, ab07ap07Key := range []string{"slot_limit", "count"} {
			ab07ap07Input := ab07ap07Week()
			if ab07ap07Key == "count" {
				ab07ap07Input["arguments"].(map[string]any)["days"].([]any)[0].(map[string]any)["count"] = ab07ap07Invalid
			} else {
				ab07ap07Input["arguments"].(map[string]any)[ab07ap07Key] = ab07ap07Invalid
			}
			ab07ap07JSON, _ := json.Marshal(ab07ap07Input)
			if _, ab07ap07Err := Ab07Ap07Evaluate(ab07ap07JSON); ab07ap07Err == nil {
				ab07ap07T.Fatalf("accepted %s=%v", ab07ap07Key, ab07ap07Invalid)
			}
		}
	}
}
func TestAb07Ap07Isolation(ab07ap07T *testing.T) {
	ab07ap07Input, _ := json.Marshal(ab07ap07Week())
	ab07ap07Saved := append([]byte(nil), ab07ap07Input...)
	ab07ap07First, ab07ap07Err := Ab07Ap07Evaluate(ab07ap07Input)
	if ab07ap07Err != nil {
		ab07ap07T.Fatal(ab07ap07Err)
	}
	for ab07ap07Index := range ab07ap07Input {
		ab07ap07Input[ab07ap07Index] = ' '
	}
	if !json.Valid(ab07ap07First) {
		ab07ap07T.Fatal("input aliases output")
	}
	ab07ap07Second, _ := Ab07Ap07Evaluate(ab07ap07Saved)
	if string(ab07ap07First) != string(ab07ap07Second) {
		ab07ap07T.Fatal("input mutation changed later result")
	}
	ab07ap07Decoded := map[string]any{}
	_ = json.Unmarshal(ab07ap07First, &ab07ap07Decoded)
	ab07ap07Decoded["raw"].(map[string]any)["evidence"].([]any)[0] = 99
	if ab07ap07Decoded["fields"].(map[string]any)["evidence"].(map[string]any)["raw"].([]any)[0] != float64(0) {
		ab07ap07T.Fatal("nested raw aliases semantic raw")
	}
	for ab07ap07Index := range ab07ap07First {
		ab07ap07First[ab07ap07Index] = ' '
	}
	ab07ap07Third, _ := Ab07Ap07Evaluate(ab07ap07Saved)
	if string(ab07ap07Second) != string(ab07ap07Third) {
		ab07ap07T.Fatal("output aliases another result")
	}
}
func TestAb07Ap07Fixtures(ab07ap07T *testing.T) {
	ab07ap07Data, ab07ap07Err := os.ReadFile("testdata/ab07ap07_b555_weekly_refresh.json")
	if ab07ap07Err != nil {
		ab07ap07T.Fatal(ab07ap07Err)
	}
	var ab07ap07Cases []struct {
		Ab07Ap07ID       string          `json:"id"`
		Ab07Ap07Input    json.RawMessage `json:"input"`
		Ab07Ap07Expected json.RawMessage `json:"expected"`
	}
	if ab07ap07Err := json.Unmarshal(ab07ap07Data, &ab07ap07Cases); ab07ap07Err != nil {
		ab07ap07T.Fatal(ab07ap07Err)
	}
	if len(ab07ap07Cases) == 0 {
		ab07ap07T.Fatal("no fixtures")
	}
	for _, ab07ap07Case := range ab07ap07Cases {
		ab07ap07T.Run(ab07ap07Case.Ab07Ap07ID, func(ab07ap07T *testing.T) {
			ab07ap07Output, ab07ap07Err := Ab07Ap07Evaluate(ab07ap07Case.Ab07Ap07Input)
			if ab07ap07Err != nil {
				ab07ap07T.Fatal(ab07ap07Err)
			}
			ab07ap07Actual := map[string]any{}
			ab07ap07Expected := map[string]any{}
			_ = json.Unmarshal(ab07ap07Output, &ab07ap07Actual)
			_ = json.Unmarshal(ab07ap07Case.Ab07Ap07Expected, &ab07ap07Expected)
			for ab07ap07Key, ab07ap07Value := range ab07ap07Expected {
				if !reflect.DeepEqual(ab07ap07Actual[ab07ap07Key], ab07ap07Value) {
					ab07ap07T.Fatalf("%s: got %#v want %#v", ab07ap07Key, ab07ap07Actual[ab07ap07Key], ab07ap07Value)
				}
			}
		})
	}
}

func TestAb07Ap07ArgumentBoundaries(ab07ap07T *testing.T) {
	ab07ap07Input := ab07ap07Week()
	ab07ap07Input["arguments"].(map[string]any)["slot_limit"] = 0
	ab07ap07Day(ab07ap07Input, 0, 0, ab07ap07Slot(false, 255))
	ab07ap07Result := ab07ap07WeeklyResult(ab07ap07T, ab07ap07Input, 0, 7, true, "complete")
	if len(ab07ap07Result["days"].([]any)[0].(map[string]any)["slots"].([]any)) != 0 {
		ab07ap07T.Fatal("count zero used extra slot")
	}
	ab07ap07Day(ab07ap07Input, 0, 1, ab07ap07Slot(true, 0))
	ab07ap07WeeklyResult(ab07ap07T, ab07ap07Input, 0, 0, false, "count-overflow")
	ab07ap07Input["arguments"].(map[string]any)["slot_limit"] = 255
	ab07ap07Day(ab07ap07Input, 0, 1, ab07ap07Slot(true, 0), ab07ap07Slot(false, 255))
	ab07ap07Result = ab07ap07WeeklyResult(ab07ap07T, ab07ap07Input, 1, 7, true, "complete")
	if len(ab07ap07Result["days"].([]any)[0].(map[string]any)["slots"].([]any)) != 1 {
		ab07ap07T.Fatal("extra slot attempted")
	}
	for _, ab07ap07Raw := range []any{[]any{}, []any{255}, []any{0, 1, 2}} {
		ab07ap07Input := ab07ap07Week()
		ab07ap07Day(ab07ap07Input, 0, 1, map[string]any{"ok": true, "raw": ab07ap07Raw})
		ab07ap07Result := ab07ap07WeeklyResult(ab07ap07T, ab07ap07Input, 1, 7, true, "complete")
		ab07ap07Actual, _ := json.Marshal(ab07ap07Result["days"].([]any)[0].(map[string]any)["slots"].([]any)[0].(map[string]any)["raw"])
		ab07ap07Expected, _ := json.Marshal(ab07ap07Raw)
		if string(ab07ap07Actual) != string(ab07ap07Expected) {
			ab07ap07T.Fatal("slot raw altered")
		}
	}
	for _, ab07ap07Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, ab07ap07Value := range []any{nil, 0} {
			ab07ap07Input := ab07ap07Week()
			if ab07ap07Key == "direction" || ab07ap07Key == "view" {
				if ab07ap07Value != nil {
					ab07ap07Value = ""
				}
			}
			ab07ap07Input["context"].(map[string]any)[ab07ap07Key] = ab07ap07Value
			delete(ab07ap07Input, "arguments")
			if ab07ap07Run(ab07ap07T, ab07ap07Input)["state"] != "Unknown" {
				ab07ap07T.Fatal("unknown selector evaluated")
			}
		}
	}
}

func TestAb07Ap07ShapeScalarsAndValidation(ab07ap07T *testing.T) {
	for ab07ap07Index := 0; ab07ap07Index < 5; ab07ap07Index++ {
		for _, ab07ap07Scalar := range []any{nil, false, "opaque", 0, 1.5} {
			ab07ap07Input := ab07ap07Week()
			ab07ap07Input["context"].(map[string]any)["operation"] = 166
			ab07ap07Input["context"].(map[string]any)["view"] = "nonheating-shape"
			ab07ap07Initial := []any{1, 2, 3, 4, 5}
			ab07ap07Initial[ab07ap07Index] = ab07ap07Scalar
			ab07ap07Input["arguments"] = map[string]any{"initial": ab07ap07Initial, "count": 1, "slot_setpoints": []any{0, 65535}}
			ab07ap07Result := ab07ap07Run(ab07ap07T, ab07ap07Input)["fields"].(map[string]any)["nonheating-shape"].(map[string]any)
			ab07ap07Expected, _ := json.Marshal(ab07ap07Initial)
			ab07ap07Actual, _ := json.Marshal(ab07ap07Result["initial"])
			if string(ab07ap07Actual) != string(ab07ap07Expected) {
				ab07ap07T.Fatal("scalar altered")
			}
		}
	}
	for _, ab07ap07Key := range []string{"initial", "count", "slot_setpoints"} {
		for _, ab07ap07Bad := range []any{nil, true, "1", -1, 1.5, []any{}, map[string]any{}} {
			ab07ap07Input := ab07ap07Week()
			ab07ap07Input["context"].(map[string]any)["operation"] = 166
			ab07ap07Input["context"].(map[string]any)["view"] = "nonheating-shape"
			ab07ap07Input["arguments"] = map[string]any{"initial": []any{1, 2, 3, 4, 5}, "count": 1, "slot_setpoints": []any{10, 20}}
			ab07ap07Input["arguments"].(map[string]any)[ab07ap07Key] = ab07ap07Bad
			ab07ap07JSON, _ := json.Marshal(ab07ap07Input)
			ab07ap07Output, ab07ap07Err := Ab07Ap07Evaluate(ab07ap07JSON)
			if ab07ap07Err == nil || ab07ap07Output != nil {
				ab07ap07T.Fatalf("accepted invalid %s=%v", ab07ap07Key, ab07ap07Bad)
			}
		}
		ab07ap07Input := ab07ap07Week()
		ab07ap07Input["context"].(map[string]any)["operation"] = 166
		ab07ap07Input["context"].(map[string]any)["view"] = "nonheating-shape"
		ab07ap07Input["arguments"] = map[string]any{"initial": []any{1, 2, 3, 4, 5}, "count": 1, "slot_setpoints": []any{10}}
		delete(ab07ap07Input["arguments"].(map[string]any), ab07ap07Key)
		ab07ap07JSON, _ := json.Marshal(ab07ap07Input)
		if _, ab07ap07Err := Ab07Ap07Evaluate(ab07ap07JSON); ab07ap07Err == nil {
			ab07ap07T.Fatalf("missing %s", ab07ap07Key)
		}
	}
	for _, ab07ap07Bad := range []any{nil, true, "1", -1, 65536, 1.5} {
		ab07ap07Input := ab07ap07Week()
		ab07ap07Input["context"].(map[string]any)["operation"] = 166
		ab07ap07Input["context"].(map[string]any)["view"] = "nonheating-shape"
		ab07ap07Input["arguments"] = map[string]any{"initial": []any{1, 2, 3, 4, 5}, "count": 1, "slot_setpoints": []any{ab07ap07Bad}}
		ab07ap07JSON, _ := json.Marshal(ab07ap07Input)
		if _, ab07ap07Err := Ab07Ap07Evaluate(ab07ap07JSON); ab07ap07Err == nil {
			ab07ap07T.Fatal("invalid u16 accepted")
		}
	}
}

func TestAb07Ap07NestedRejections(ab07ap07T *testing.T) {
	for _, ab07ap07Key := range []string{"ok", "raw"} {
		for _, ab07ap07Bad := range []any{nil, 0, "false", []any{-1}, []any{256}, []any{nil}} {
			ab07ap07Input := ab07ap07Week()
			ab07ap07SlotValue := map[string]any{"ok": true, "raw": []any{1}}
			ab07ap07SlotValue[ab07ap07Key] = ab07ap07Bad
			ab07ap07Day(ab07ap07Input, 0, 1, ab07ap07SlotValue)
			ab07ap07JSON, _ := json.Marshal(ab07ap07Input)
			ab07ap07Output, ab07ap07Err := Ab07Ap07Evaluate(ab07ap07JSON)
			if ab07ap07Err == nil || ab07ap07Output != nil {
				ab07ap07T.Fatalf("accepted slot %s=%v", ab07ap07Key, ab07ap07Bad)
			}
		}
	}
	for _, ab07ap07Bad := range []string{
		`{"context":{},"fields":{},"arguments":{"extra":{"x":1,"x":2}}}`,
		`{"context":{"family":{"x":1,"x":2}},"fields":{}}`,
		`{"context":{},"fields":{"x":[],"x":[]}}`,
		`{"context":{},"fields":null}`,
		`{"context":null,"fields":{}}`,
		`{"context":{},"fields":{},"arguments":null}`,
	} {
		ab07ap07Output, ab07ap07Err := Ab07Ap07Evaluate(json.RawMessage(ab07ap07Bad))
		if ab07ap07Err == nil || ab07ap07Output != nil {
			ab07ap07T.Fatalf("accepted %s", ab07ap07Bad)
		}
	}
	for _, ab07ap07Size := range []int{0, 6, 8} {
		ab07ap07Input := ab07ap07Week()
		ab07ap07Days := ab07ap07Input["arguments"].(map[string]any)["days"].([]any)
		if ab07ap07Size == 8 {
			ab07ap07Days = append(ab07ap07Days, ab07ap07Days[0])
		} else {
			ab07ap07Days = ab07ap07Days[:ab07ap07Size]
		}
		ab07ap07Input["arguments"].(map[string]any)["days"] = ab07ap07Days
		ab07ap07JSON, _ := json.Marshal(ab07ap07Input)
		if _, ab07ap07Err := Ab07Ap07Evaluate(ab07ap07JSON); ab07ap07Err == nil {
			ab07ap07T.Fatal("wrong day length accepted")
		}
	}
}

func TestAb07Ap07NamedRawAndShapeBounds(ab07ap07T *testing.T) {
	for _, ab07ap07View := range []string{"weekly", "nonheating-shape"} {
		ab07ap07Input := ab07ap07Week()
		ab07ap07Input["fields"].(map[string]any)[ab07ap07View] = []any{17, 0, 255}
		if ab07ap07View == "nonheating-shape" {
			ab07ap07Input["context"].(map[string]any)["operation"] = 166
			ab07ap07Input["context"].(map[string]any)["view"] = ab07ap07View
			ab07ap07Input["arguments"] = map[string]any{"initial": []any{1, 2, 3, 4, 5}, "count": 0, "slot_setpoints": []any{}}
		}
		ab07ap07Output := ab07ap07Run(ab07ap07T, ab07ap07Input)
		if !reflect.DeepEqual(ab07ap07Output["fields"].(map[string]any)[ab07ap07View].(map[string]any)["raw"], []any{float64(17), float64(0), float64(255)}) {
			ab07ap07T.Fatal("named raw overwritten")
		}
	}
	for _, ab07ap07Count := range []int{0, 255, 256} {
		ab07ap07Input := ab07ap07Week()
		ab07ap07Input["context"].(map[string]any)["operation"] = 166
		ab07ap07Input["context"].(map[string]any)["view"] = "nonheating-shape"
		ab07ap07Setpoints := []any{}
		for ab07ap07Index := 0; ab07ap07Index < ab07ap07Count; ab07ap07Index++ {
			ab07ap07Setpoints = append(ab07ap07Setpoints, ab07ap07Index)
		}
		ab07ap07Input["arguments"] = map[string]any{"initial": []any{1, 2, 3, 4, 5}, "count": ab07ap07Count, "slot_setpoints": ab07ap07Setpoints}
		if ab07ap07Count == 256 {
			ab07ap07JSON, _ := json.Marshal(ab07ap07Input)
			if _, ab07ap07Err := Ab07Ap07Evaluate(ab07ap07JSON); ab07ap07Err == nil {
				ab07ap07T.Fatal("count overflow accepted")
			}
			continue
		}
		ab07ap07Result := ab07ap07Run(ab07ap07T, ab07ap07Input)["fields"].(map[string]any)["nonheating-shape"].(map[string]any)
		if len(ab07ap07Result["processed_setpoints"].([]any)) != ab07ap07Count {
			ab07ap07T.Fatal("count boundary changed")
		}
	}
}
