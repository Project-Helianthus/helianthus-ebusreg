package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func b555WeeklyRefreshWeek() map[string]any {
	b555WeeklyRefreshDays := []any{}
	for b555WeeklyRefreshIndex := 0; b555WeeklyRefreshIndex < 7; b555WeeklyRefreshIndex++ {
		b555WeeklyRefreshDays = append(b555WeeklyRefreshDays, map[string]any{"count": 0, "slots": []any{}})
	}
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 85, "operation": 164, "direction": "local", "view": "weekly", "family": "synthetic-family", "variant": "synthetic-variant"}, "fields": map[string]any{"evidence": []any{0, 255}}, "arguments": map[string]any{"status_ok": true, "slot_limit": 2, "days": b555WeeklyRefreshDays}}
}
func b555WeeklyRefreshDay(b555WeeklyRefreshInput map[string]any, b555WeeklyRefreshIndex int, b555WeeklyRefreshCount int, b555WeeklyRefreshSlots ...any) {
	if b555WeeklyRefreshSlots == nil {
		b555WeeklyRefreshSlots = []any{}
	}
	b555WeeklyRefreshInput["arguments"].(map[string]any)["days"].([]any)[b555WeeklyRefreshIndex] = map[string]any{"count": b555WeeklyRefreshCount, "slots": b555WeeklyRefreshSlots}
}
func b555WeeklyRefreshSlot(b555WeeklyRefreshOK bool, b555WeeklyRefreshByte int) any {
	return map[string]any{"ok": b555WeeklyRefreshOK, "raw": []any{b555WeeklyRefreshByte}}
}
func b555WeeklyRefreshRun(b555WeeklyRefreshT *testing.T, b555WeeklyRefreshInput map[string]any) map[string]any {
	b555WeeklyRefreshT.Helper()
	b555WeeklyRefreshJSON, b555WeeklyRefreshErr := json.Marshal(b555WeeklyRefreshInput)
	if b555WeeklyRefreshErr != nil {
		b555WeeklyRefreshT.Fatal(b555WeeklyRefreshErr)
	}
	b555WeeklyRefreshOutput, b555WeeklyRefreshErr := B555WeeklyRefreshEvaluate(b555WeeklyRefreshJSON)
	if b555WeeklyRefreshErr != nil {
		b555WeeklyRefreshT.Fatal(b555WeeklyRefreshErr)
	}
	b555WeeklyRefreshResult := map[string]any{}
	if b555WeeklyRefreshErr := json.Unmarshal(b555WeeklyRefreshOutput, &b555WeeklyRefreshResult); b555WeeklyRefreshErr != nil {
		b555WeeklyRefreshT.Fatal(b555WeeklyRefreshErr)
	}
	if b555WeeklyRefreshResult["physical"] != "Unknown" || b555WeeklyRefreshResult["qualified"] != false {
		b555WeeklyRefreshT.Fatalf("qualification: %s", b555WeeklyRefreshOutput)
	}
	return b555WeeklyRefreshResult
}
func b555WeeklyRefreshWeeklyResult(b555WeeklyRefreshT *testing.T, b555WeeklyRefreshInput map[string]any, b555WeeklyRefreshAttempts, b555WeeklyRefreshDays int, b555WeeklyRefreshComplete bool, b555WeeklyRefreshReason string) map[string]any {
	b555WeeklyRefreshT.Helper()
	b555WeeklyRefreshOutput := b555WeeklyRefreshRun(b555WeeklyRefreshT, b555WeeklyRefreshInput)
	b555WeeklyRefreshResult := b555WeeklyRefreshOutput["fields"].(map[string]any)["weekly"].(map[string]any)
	if b555WeeklyRefreshResult["attempts"] != float64(b555WeeklyRefreshAttempts) || b555WeeklyRefreshResult["completed_days"] != float64(b555WeeklyRefreshDays) || b555WeeklyRefreshResult["complete"] != b555WeeklyRefreshComplete || b555WeeklyRefreshResult["reason"] != b555WeeklyRefreshReason || b555WeeklyRefreshResult["state"] != "Hypothesis" {
		b555WeeklyRefreshT.Fatalf("wrong reduction: %#v", b555WeeklyRefreshResult)
	}
	b555WeeklyRefreshState := "partial"
	if b555WeeklyRefreshComplete {
		b555WeeklyRefreshState = "complete"
	}
	if b555WeeklyRefreshOutput["state"] != b555WeeklyRefreshState {
		b555WeeklyRefreshT.Fatalf("wrong state: %#v", b555WeeklyRefreshOutput)
	}
	if b555WeeklyRefreshOutput["fields"].(map[string]any)["evidence"].(map[string]any)["state"] != "Unknown" {
		b555WeeklyRefreshT.Fatal("invented raw interpretation")
	}
	return b555WeeklyRefreshResult
}
func TestAb07Ap07WeeklyTransitions(b555WeeklyRefreshT *testing.T) {
	b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
	b555WeeklyRefreshWeeklyResult(b555WeeklyRefreshT, b555WeeklyRefreshInput, 0, 7, true, "complete")
	b555WeeklyRefreshDay(b555WeeklyRefreshInput, 0, 1, b555WeeklyRefreshSlot(true, 11))
	b555WeeklyRefreshInput["arguments"].(map[string]any)["status_ok"] = false
	b555WeeklyRefreshWeeklyResult(b555WeeklyRefreshT, b555WeeklyRefreshInput, 0, 0, false, "status-failure")
	b555WeeklyRefreshInput["arguments"].(map[string]any)["status_ok"] = true
	b555WeeklyRefreshWeeklyResult(b555WeeklyRefreshT, b555WeeklyRefreshInput, 1, 7, true, "complete")
	b555WeeklyRefreshDay(b555WeeklyRefreshInput, 1, 1, b555WeeklyRefreshSlot(true, 22))
	b555WeeklyRefreshDay(b555WeeklyRefreshInput, 2, 3)
	b555WeeklyRefreshResult := b555WeeklyRefreshWeeklyResult(b555WeeklyRefreshT, b555WeeklyRefreshInput, 2, 2, false, "count-overflow")
	if len(b555WeeklyRefreshResult["days"].([]any)) != 2 {
		b555WeeklyRefreshT.Fatal("lost completed days")
	}
	b555WeeklyRefreshInput = b555WeeklyRefreshWeek()
	b555WeeklyRefreshDay(b555WeeklyRefreshInput, 0, 1, b555WeeklyRefreshSlot(false, 33))
	b555WeeklyRefreshResult = b555WeeklyRefreshWeeklyResult(b555WeeklyRefreshT, b555WeeklyRefreshInput, 1, 0, false, "slot-failure")
	if len(b555WeeklyRefreshResult["days"].([]any)[0].(map[string]any)["slots"].([]any)) != 0 {
		b555WeeklyRefreshT.Fatal("failed slot succeeded")
	}
	b555WeeklyRefreshDay(b555WeeklyRefreshInput, 0, 1)
	b555WeeklyRefreshWeeklyResult(b555WeeklyRefreshT, b555WeeklyRefreshInput, 0, 0, false, "missing-slot")
	b555WeeklyRefreshDay(b555WeeklyRefreshInput, 0, 1, map[string]any{"ok": true})
	b555WeeklyRefreshWeeklyResult(b555WeeklyRefreshT, b555WeeklyRefreshInput, 1, 0, false, "missing-slot-data")
}
func TestAb07Ap07EveryDayAndSlot(b555WeeklyRefreshT *testing.T) {
	for b555WeeklyRefreshDayIndex := 0; b555WeeklyRefreshDayIndex < 7; b555WeeklyRefreshDayIndex++ {
		for b555WeeklyRefreshSlotIndex := 0; b555WeeklyRefreshSlotIndex < 2; b555WeeklyRefreshSlotIndex++ {
			b555WeeklyRefreshT.Run(string(rune('A'+b555WeeklyRefreshDayIndex))+string(rune('0'+b555WeeklyRefreshSlotIndex)), func(b555WeeklyRefreshT *testing.T) {
				b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
				for b555WeeklyRefreshIndex := 0; b555WeeklyRefreshIndex < 7; b555WeeklyRefreshIndex++ {
					b555WeeklyRefreshDay(b555WeeklyRefreshInput, b555WeeklyRefreshIndex, 2, b555WeeklyRefreshSlot(true, 10), b555WeeklyRefreshSlot(true, 20))
				}
				b555WeeklyRefreshWeeklyResult(b555WeeklyRefreshT, b555WeeklyRefreshInput, 14, 7, true, "complete")
				b555WeeklyRefreshInput["arguments"].(map[string]any)["days"].([]any)[b555WeeklyRefreshDayIndex].(map[string]any)["slots"].([]any)[b555WeeklyRefreshSlotIndex].(map[string]any)["ok"] = false
				b555WeeklyRefreshResult := b555WeeklyRefreshWeeklyResult(b555WeeklyRefreshT, b555WeeklyRefreshInput, b555WeeklyRefreshDayIndex*2+b555WeeklyRefreshSlotIndex+1, b555WeeklyRefreshDayIndex, false, "slot-failure")
				b555WeeklyRefreshSlots := b555WeeklyRefreshResult["days"].([]any)[b555WeeklyRefreshDayIndex].(map[string]any)["slots"].([]any)
				if len(b555WeeklyRefreshSlots) != b555WeeklyRefreshSlotIndex {
					b555WeeklyRefreshT.Fatal("lost successful prefix")
				}
				if b555WeeklyRefreshSlotIndex == 1 && b555WeeklyRefreshSlots[0].(map[string]any)["raw"].([]any)[0] != float64(10) {
					b555WeeklyRefreshT.Fatal("raw slot altered")
				}
			})
		}
	}
}
func TestAb07Ap07Nonheating(b555WeeklyRefreshT *testing.T) {
	for b555WeeklyRefreshCount := 0; b555WeeklyRefreshCount <= 2; b555WeeklyRefreshCount++ {
		b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
		b555WeeklyRefreshInput["context"].(map[string]any)["operation"] = 166
		b555WeeklyRefreshInput["context"].(map[string]any)["view"] = "nonheating-shape"
		b555WeeklyRefreshInput["arguments"] = map[string]any{"initial": []any{1, 2, 3, 4, 5}, "count": b555WeeklyRefreshCount, "slot_setpoints": []any{100, 200, 300}}
		b555WeeklyRefreshResult := b555WeeklyRefreshRun(b555WeeklyRefreshT, b555WeeklyRefreshInput)["fields"].(map[string]any)["nonheating-shape"].(map[string]any)
		if !reflect.DeepEqual(b555WeeklyRefreshResult["initial"], []any{float64(1), float64(2), float64(3), float64(4), float64(5)}) || b555WeeklyRefreshResult["processed_slots"] != float64(b555WeeklyRefreshCount) || b555WeeklyRefreshResult["complete"] != true || b555WeeklyRefreshResult["state"] != "Hypothesis" {
			b555WeeklyRefreshT.Fatalf("shape: %#v", b555WeeklyRefreshResult)
		}
		b555WeeklyRefreshProcessed := b555WeeklyRefreshResult["processed_setpoints"].([]any)
		if len(b555WeeklyRefreshProcessed) != b555WeeklyRefreshCount {
			b555WeeklyRefreshT.Fatal("invented processed slots")
		}
		for _, b555WeeklyRefreshValue := range b555WeeklyRefreshProcessed {
			if b555WeeklyRefreshValue != float64(65535) {
				b555WeeklyRefreshT.Fatal("missing FFFF")
			}
		}
		if !reflect.DeepEqual(b555WeeklyRefreshResult["unused_setpoints"], []any{float64(100), float64(200), float64(300)}[b555WeeklyRefreshCount:]) {
			b555WeeklyRefreshT.Fatal("prefilled unused values")
		}
	}
}
func TestAb07Ap07ContextAndRaw(b555WeeklyRefreshT *testing.T) {
	for _, b555WeeklyRefreshKey := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
		b555WeeklyRefreshT.Run(b555WeeklyRefreshKey, func(b555WeeklyRefreshT *testing.T) {
			for _, b555WeeklyRefreshValue := range []any{nil, "different"} {
				b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
				b555WeeklyRefreshContext := b555WeeklyRefreshInput["context"].(map[string]any)
				if b555WeeklyRefreshKey == "primary" || b555WeeklyRefreshKey == "secondary" || b555WeeklyRefreshKey == "operation" {
					b555WeeklyRefreshValue = 0
				}
				b555WeeklyRefreshContext[b555WeeklyRefreshKey] = b555WeeklyRefreshValue
				b555WeeklyRefreshOutput := b555WeeklyRefreshRun(b555WeeklyRefreshT, b555WeeklyRefreshInput)
				if (b555WeeklyRefreshKey == "family" || b555WeeklyRefreshKey == "variant") && b555WeeklyRefreshOutput["state"] != "complete" {
					b555WeeklyRefreshT.Fatal("opaque context gated local interpretation")
				}
				if b555WeeklyRefreshKey != "family" && b555WeeklyRefreshKey != "variant" && b555WeeklyRefreshOutput["state"] != "Unknown" {
					b555WeeklyRefreshT.Fatal("selector gained interpretation")
				}
				b555WeeklyRefreshEncoded, _ := json.Marshal(b555WeeklyRefreshContext)
				b555WeeklyRefreshReturned, _ := json.Marshal(b555WeeklyRefreshOutput["context"])
				if string(b555WeeklyRefreshEncoded) != string(b555WeeklyRefreshReturned) {
					b555WeeklyRefreshT.Fatal("context changed")
				}
			}
			b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
			delete(b555WeeklyRefreshInput["context"].(map[string]any), b555WeeklyRefreshKey)
			b555WeeklyRefreshOutput := b555WeeklyRefreshRun(b555WeeklyRefreshT, b555WeeklyRefreshInput)
			if _, b555WeeklyRefreshPresent := b555WeeklyRefreshOutput["context"].(map[string]any)[b555WeeklyRefreshKey]; b555WeeklyRefreshPresent {
				b555WeeklyRefreshT.Fatal("missing became present")
			}
		})
	}
	for _, b555WeeklyRefreshRaw := range []map[string]any{{}, {"evidence": []any{}}, {"evidence": []any{0}}, {"evidence": []any{255}, "sibling": []any{1, 2, 3}}} {
		b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
		b555WeeklyRefreshInput["fields"] = b555WeeklyRefreshRaw
		b555WeeklyRefreshInput["arguments"].(map[string]any)["extra"] = map[string]any{"nested": []any{nil, false, 0}}
		b555WeeklyRefreshOutput := b555WeeklyRefreshRun(b555WeeklyRefreshT, b555WeeklyRefreshInput)
		b555WeeklyRefreshExpected, _ := json.Marshal(b555WeeklyRefreshRaw)
		b555WeeklyRefreshActual, _ := json.Marshal(b555WeeklyRefreshOutput["raw"])
		if string(b555WeeklyRefreshExpected) != string(b555WeeklyRefreshActual) {
			b555WeeklyRefreshT.Fatal("raw changed")
		}
		b555WeeklyRefreshExpected, _ = json.Marshal(b555WeeklyRefreshInput["arguments"])
		b555WeeklyRefreshActual, _ = json.Marshal(b555WeeklyRefreshOutput["arguments"])
		if string(b555WeeklyRefreshExpected) != string(b555WeeklyRefreshActual) {
			b555WeeklyRefreshT.Fatal("arguments changed")
		}
	}
}
func TestAb07Ap07Rejections(b555WeeklyRefreshT *testing.T) {
	b555WeeklyRefreshJSON, _ := json.Marshal(b555WeeklyRefreshWeek())
	b555WeeklyRefreshBase := string(b555WeeklyRefreshJSON)
	for _, b555WeeklyRefreshInvalid := range []string{
		"", "null", "[]", "{", b555WeeklyRefreshBase + " {}",
		strings.Replace(b555WeeklyRefreshBase, `"primary":181`, `"primary":181,"primary":181`, 1),
		strings.Replace(b555WeeklyRefreshBase, `"count":0`, `"count":0,"count":0`, 1),
		strings.Replace(b555WeeklyRefreshBase, `"evidence":[0,255]`, `"evidence":[256]`, 1),
		strings.Replace(b555WeeklyRefreshBase, `"evidence":[0,255]`, `"evidence":[true]`, 1),
		strings.Replace(b555WeeklyRefreshBase, `"evidence":[0,255]`, `"evidence":[1.5]`, 1),
		strings.Replace(b555WeeklyRefreshBase, `"slot_limit":2`, `"slot_limit":null`, 1),
		strings.Replace(b555WeeklyRefreshBase, `"status_ok":true`, `"status_ok":0`, 1),
		strings.Replace(b555WeeklyRefreshBase, `"primary":181`, `"primary":"181"`, 1),
	} {
		b555WeeklyRefreshOutput, b555WeeklyRefreshErr := B555WeeklyRefreshEvaluate(json.RawMessage(b555WeeklyRefreshInvalid))
		if b555WeeklyRefreshErr == nil || b555WeeklyRefreshOutput != nil {
			b555WeeklyRefreshT.Fatalf("accepted invalid: %s", b555WeeklyRefreshInvalid)
		}
	}
	for _, b555WeeklyRefreshKey := range []string{"status_ok", "slot_limit", "days"} {
		b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
		delete(b555WeeklyRefreshInput["arguments"].(map[string]any), b555WeeklyRefreshKey)
		b555WeeklyRefreshJSON, _ := json.Marshal(b555WeeklyRefreshInput)
		if _, b555WeeklyRefreshErr := B555WeeklyRefreshEvaluate(b555WeeklyRefreshJSON); b555WeeklyRefreshErr == nil {
			b555WeeklyRefreshT.Fatalf("missing %s accepted", b555WeeklyRefreshKey)
		}
	}
	for _, b555WeeklyRefreshInvalid := range []any{nil, true, "1", -1, 256, 1.5, []any{}} {
		for _, b555WeeklyRefreshKey := range []string{"slot_limit", "count"} {
			b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
			if b555WeeklyRefreshKey == "count" {
				b555WeeklyRefreshInput["arguments"].(map[string]any)["days"].([]any)[0].(map[string]any)["count"] = b555WeeklyRefreshInvalid
			} else {
				b555WeeklyRefreshInput["arguments"].(map[string]any)[b555WeeklyRefreshKey] = b555WeeklyRefreshInvalid
			}
			b555WeeklyRefreshJSON, _ := json.Marshal(b555WeeklyRefreshInput)
			if _, b555WeeklyRefreshErr := B555WeeklyRefreshEvaluate(b555WeeklyRefreshJSON); b555WeeklyRefreshErr == nil {
				b555WeeklyRefreshT.Fatalf("accepted %s=%v", b555WeeklyRefreshKey, b555WeeklyRefreshInvalid)
			}
		}
	}
}
func TestAb07Ap07Isolation(b555WeeklyRefreshT *testing.T) {
	b555WeeklyRefreshInput, _ := json.Marshal(b555WeeklyRefreshWeek())
	b555WeeklyRefreshSaved := append([]byte(nil), b555WeeklyRefreshInput...)
	b555WeeklyRefreshFirst, b555WeeklyRefreshErr := B555WeeklyRefreshEvaluate(b555WeeklyRefreshInput)
	if b555WeeklyRefreshErr != nil {
		b555WeeklyRefreshT.Fatal(b555WeeklyRefreshErr)
	}
	for b555WeeklyRefreshIndex := range b555WeeklyRefreshInput {
		b555WeeklyRefreshInput[b555WeeklyRefreshIndex] = ' '
	}
	if !json.Valid(b555WeeklyRefreshFirst) {
		b555WeeklyRefreshT.Fatal("input aliases output")
	}
	b555WeeklyRefreshSecond, _ := B555WeeklyRefreshEvaluate(b555WeeklyRefreshSaved)
	if string(b555WeeklyRefreshFirst) != string(b555WeeklyRefreshSecond) {
		b555WeeklyRefreshT.Fatal("input mutation changed later result")
	}
	b555WeeklyRefreshDecoded := map[string]any{}
	_ = json.Unmarshal(b555WeeklyRefreshFirst, &b555WeeklyRefreshDecoded)
	b555WeeklyRefreshDecoded["raw"].(map[string]any)["evidence"].([]any)[0] = 99
	if b555WeeklyRefreshDecoded["fields"].(map[string]any)["evidence"].(map[string]any)["raw"].([]any)[0] != float64(0) {
		b555WeeklyRefreshT.Fatal("nested raw aliases semantic raw")
	}
	for b555WeeklyRefreshIndex := range b555WeeklyRefreshFirst {
		b555WeeklyRefreshFirst[b555WeeklyRefreshIndex] = ' '
	}
	b555WeeklyRefreshThird, _ := B555WeeklyRefreshEvaluate(b555WeeklyRefreshSaved)
	if string(b555WeeklyRefreshSecond) != string(b555WeeklyRefreshThird) {
		b555WeeklyRefreshT.Fatal("output aliases another result")
	}
}
func TestAb07Ap07Fixtures(b555WeeklyRefreshT *testing.T) {
	b555WeeklyRefreshData, b555WeeklyRefreshErr := os.ReadFile("testdata/ab07ap07_b555_weekly_refresh.json")
	if b555WeeklyRefreshErr != nil {
		b555WeeklyRefreshT.Fatal(b555WeeklyRefreshErr)
	}
	var b555WeeklyRefreshCases []struct {
		B555WeeklyRefreshID       string          `json:"id"`
		B555WeeklyRefreshInput    json.RawMessage `json:"input"`
		B555WeeklyRefreshExpected json.RawMessage `json:"expected"`
	}
	if b555WeeklyRefreshErr := json.Unmarshal(b555WeeklyRefreshData, &b555WeeklyRefreshCases); b555WeeklyRefreshErr != nil {
		b555WeeklyRefreshT.Fatal(b555WeeklyRefreshErr)
	}
	if len(b555WeeklyRefreshCases) == 0 {
		b555WeeklyRefreshT.Fatal("no fixtures")
	}
	for _, b555WeeklyRefreshCase := range b555WeeklyRefreshCases {
		b555WeeklyRefreshT.Run(b555WeeklyRefreshCase.B555WeeklyRefreshID, func(b555WeeklyRefreshT *testing.T) {
			b555WeeklyRefreshOutput, b555WeeklyRefreshErr := B555WeeklyRefreshEvaluate(b555WeeklyRefreshCase.B555WeeklyRefreshInput)
			if b555WeeklyRefreshErr != nil {
				b555WeeklyRefreshT.Fatal(b555WeeklyRefreshErr)
			}
			b555WeeklyRefreshActual := map[string]any{}
			b555WeeklyRefreshExpected := map[string]any{}
			_ = json.Unmarshal(b555WeeklyRefreshOutput, &b555WeeklyRefreshActual)
			_ = json.Unmarshal(b555WeeklyRefreshCase.B555WeeklyRefreshExpected, &b555WeeklyRefreshExpected)
			for b555WeeklyRefreshKey, b555WeeklyRefreshValue := range b555WeeklyRefreshExpected {
				if !reflect.DeepEqual(b555WeeklyRefreshActual[b555WeeklyRefreshKey], b555WeeklyRefreshValue) {
					b555WeeklyRefreshT.Fatalf("%s: got %#v want %#v", b555WeeklyRefreshKey, b555WeeklyRefreshActual[b555WeeklyRefreshKey], b555WeeklyRefreshValue)
				}
			}
		})
	}
}

func TestAb07Ap07ArgumentBoundaries(b555WeeklyRefreshT *testing.T) {
	b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
	b555WeeklyRefreshInput["arguments"].(map[string]any)["slot_limit"] = 0
	b555WeeklyRefreshDay(b555WeeklyRefreshInput, 0, 0, b555WeeklyRefreshSlot(false, 255))
	b555WeeklyRefreshResult := b555WeeklyRefreshWeeklyResult(b555WeeklyRefreshT, b555WeeklyRefreshInput, 0, 7, true, "complete")
	if len(b555WeeklyRefreshResult["days"].([]any)[0].(map[string]any)["slots"].([]any)) != 0 {
		b555WeeklyRefreshT.Fatal("count zero used extra slot")
	}
	b555WeeklyRefreshDay(b555WeeklyRefreshInput, 0, 1, b555WeeklyRefreshSlot(true, 0))
	b555WeeklyRefreshWeeklyResult(b555WeeklyRefreshT, b555WeeklyRefreshInput, 0, 0, false, "count-overflow")
	b555WeeklyRefreshInput["arguments"].(map[string]any)["slot_limit"] = 255
	b555WeeklyRefreshDay(b555WeeklyRefreshInput, 0, 1, b555WeeklyRefreshSlot(true, 0), b555WeeklyRefreshSlot(false, 255))
	b555WeeklyRefreshResult = b555WeeklyRefreshWeeklyResult(b555WeeklyRefreshT, b555WeeklyRefreshInput, 1, 7, true, "complete")
	if len(b555WeeklyRefreshResult["days"].([]any)[0].(map[string]any)["slots"].([]any)) != 1 {
		b555WeeklyRefreshT.Fatal("extra slot attempted")
	}
	for _, b555WeeklyRefreshRaw := range []any{[]any{}, []any{255}, []any{0, 1, 2}} {
		b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
		b555WeeklyRefreshDay(b555WeeklyRefreshInput, 0, 1, map[string]any{"ok": true, "raw": b555WeeklyRefreshRaw})
		b555WeeklyRefreshResult := b555WeeklyRefreshWeeklyResult(b555WeeklyRefreshT, b555WeeklyRefreshInput, 1, 7, true, "complete")
		b555WeeklyRefreshActual, _ := json.Marshal(b555WeeklyRefreshResult["days"].([]any)[0].(map[string]any)["slots"].([]any)[0].(map[string]any)["raw"])
		b555WeeklyRefreshExpected, _ := json.Marshal(b555WeeklyRefreshRaw)
		if string(b555WeeklyRefreshActual) != string(b555WeeklyRefreshExpected) {
			b555WeeklyRefreshT.Fatal("slot raw altered")
		}
	}
	for _, b555WeeklyRefreshKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, b555WeeklyRefreshValue := range []any{nil, 0} {
			b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
			if b555WeeklyRefreshKey == "direction" || b555WeeklyRefreshKey == "view" {
				if b555WeeklyRefreshValue != nil {
					b555WeeklyRefreshValue = ""
				}
			}
			b555WeeklyRefreshInput["context"].(map[string]any)[b555WeeklyRefreshKey] = b555WeeklyRefreshValue
			delete(b555WeeklyRefreshInput, "arguments")
			if b555WeeklyRefreshRun(b555WeeklyRefreshT, b555WeeklyRefreshInput)["state"] != "Unknown" {
				b555WeeklyRefreshT.Fatal("unknown selector evaluated")
			}
		}
	}
}

func TestAb07Ap07ShapeScalarsAndValidation(b555WeeklyRefreshT *testing.T) {
	for b555WeeklyRefreshIndex := 0; b555WeeklyRefreshIndex < 5; b555WeeklyRefreshIndex++ {
		for _, b555WeeklyRefreshScalar := range []any{nil, false, "opaque", 0, 1.5} {
			b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
			b555WeeklyRefreshInput["context"].(map[string]any)["operation"] = 166
			b555WeeklyRefreshInput["context"].(map[string]any)["view"] = "nonheating-shape"
			b555WeeklyRefreshInitial := []any{1, 2, 3, 4, 5}
			b555WeeklyRefreshInitial[b555WeeklyRefreshIndex] = b555WeeklyRefreshScalar
			b555WeeklyRefreshInput["arguments"] = map[string]any{"initial": b555WeeklyRefreshInitial, "count": 1, "slot_setpoints": []any{0, 65535}}
			b555WeeklyRefreshResult := b555WeeklyRefreshRun(b555WeeklyRefreshT, b555WeeklyRefreshInput)["fields"].(map[string]any)["nonheating-shape"].(map[string]any)
			b555WeeklyRefreshExpected, _ := json.Marshal(b555WeeklyRefreshInitial)
			b555WeeklyRefreshActual, _ := json.Marshal(b555WeeklyRefreshResult["initial"])
			if string(b555WeeklyRefreshActual) != string(b555WeeklyRefreshExpected) {
				b555WeeklyRefreshT.Fatal("scalar altered")
			}
		}
	}
	for _, b555WeeklyRefreshKey := range []string{"initial", "count", "slot_setpoints"} {
		for _, b555WeeklyRefreshBad := range []any{nil, true, "1", -1, 1.5, []any{}, map[string]any{}} {
			b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
			b555WeeklyRefreshInput["context"].(map[string]any)["operation"] = 166
			b555WeeklyRefreshInput["context"].(map[string]any)["view"] = "nonheating-shape"
			b555WeeklyRefreshInput["arguments"] = map[string]any{"initial": []any{1, 2, 3, 4, 5}, "count": 1, "slot_setpoints": []any{10, 20}}
			b555WeeklyRefreshInput["arguments"].(map[string]any)[b555WeeklyRefreshKey] = b555WeeklyRefreshBad
			b555WeeklyRefreshJSON, _ := json.Marshal(b555WeeklyRefreshInput)
			b555WeeklyRefreshOutput, b555WeeklyRefreshErr := B555WeeklyRefreshEvaluate(b555WeeklyRefreshJSON)
			if b555WeeklyRefreshErr == nil || b555WeeklyRefreshOutput != nil {
				b555WeeklyRefreshT.Fatalf("accepted invalid %s=%v", b555WeeklyRefreshKey, b555WeeklyRefreshBad)
			}
		}
		b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
		b555WeeklyRefreshInput["context"].(map[string]any)["operation"] = 166
		b555WeeklyRefreshInput["context"].(map[string]any)["view"] = "nonheating-shape"
		b555WeeklyRefreshInput["arguments"] = map[string]any{"initial": []any{1, 2, 3, 4, 5}, "count": 1, "slot_setpoints": []any{10}}
		delete(b555WeeklyRefreshInput["arguments"].(map[string]any), b555WeeklyRefreshKey)
		b555WeeklyRefreshJSON, _ := json.Marshal(b555WeeklyRefreshInput)
		if _, b555WeeklyRefreshErr := B555WeeklyRefreshEvaluate(b555WeeklyRefreshJSON); b555WeeklyRefreshErr == nil {
			b555WeeklyRefreshT.Fatalf("missing %s", b555WeeklyRefreshKey)
		}
	}
	for _, b555WeeklyRefreshBad := range []any{nil, true, "1", -1, 65536, 1.5} {
		b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
		b555WeeklyRefreshInput["context"].(map[string]any)["operation"] = 166
		b555WeeklyRefreshInput["context"].(map[string]any)["view"] = "nonheating-shape"
		b555WeeklyRefreshInput["arguments"] = map[string]any{"initial": []any{1, 2, 3, 4, 5}, "count": 1, "slot_setpoints": []any{b555WeeklyRefreshBad}}
		b555WeeklyRefreshJSON, _ := json.Marshal(b555WeeklyRefreshInput)
		if _, b555WeeklyRefreshErr := B555WeeklyRefreshEvaluate(b555WeeklyRefreshJSON); b555WeeklyRefreshErr == nil {
			b555WeeklyRefreshT.Fatal("invalid u16 accepted")
		}
	}
}

func TestAb07Ap07NestedRejections(b555WeeklyRefreshT *testing.T) {
	for _, b555WeeklyRefreshKey := range []string{"ok", "raw"} {
		for _, b555WeeklyRefreshBad := range []any{nil, 0, "false", []any{-1}, []any{256}, []any{nil}} {
			b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
			b555WeeklyRefreshSlotValue := map[string]any{"ok": true, "raw": []any{1}}
			b555WeeklyRefreshSlotValue[b555WeeklyRefreshKey] = b555WeeklyRefreshBad
			b555WeeklyRefreshDay(b555WeeklyRefreshInput, 0, 1, b555WeeklyRefreshSlotValue)
			b555WeeklyRefreshJSON, _ := json.Marshal(b555WeeklyRefreshInput)
			b555WeeklyRefreshOutput, b555WeeklyRefreshErr := B555WeeklyRefreshEvaluate(b555WeeklyRefreshJSON)
			if b555WeeklyRefreshErr == nil || b555WeeklyRefreshOutput != nil {
				b555WeeklyRefreshT.Fatalf("accepted slot %s=%v", b555WeeklyRefreshKey, b555WeeklyRefreshBad)
			}
		}
	}
	for _, b555WeeklyRefreshBad := range []string{
		`{"context":{},"fields":{},"arguments":{"extra":{"x":1,"x":2}}}`,
		`{"context":{"family":{"x":1,"x":2}},"fields":{}}`,
		`{"context":{},"fields":{"x":[],"x":[]}}`,
		`{"context":{},"fields":null}`,
		`{"context":null,"fields":{}}`,
		`{"context":{},"fields":{},"arguments":null}`,
	} {
		b555WeeklyRefreshOutput, b555WeeklyRefreshErr := B555WeeklyRefreshEvaluate(json.RawMessage(b555WeeklyRefreshBad))
		if b555WeeklyRefreshErr == nil || b555WeeklyRefreshOutput != nil {
			b555WeeklyRefreshT.Fatalf("accepted %s", b555WeeklyRefreshBad)
		}
	}
	for _, b555WeeklyRefreshSize := range []int{0, 6, 8} {
		b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
		b555WeeklyRefreshDays := b555WeeklyRefreshInput["arguments"].(map[string]any)["days"].([]any)
		if b555WeeklyRefreshSize == 8 {
			b555WeeklyRefreshDays = append(b555WeeklyRefreshDays, b555WeeklyRefreshDays[0])
		} else {
			b555WeeklyRefreshDays = b555WeeklyRefreshDays[:b555WeeklyRefreshSize]
		}
		b555WeeklyRefreshInput["arguments"].(map[string]any)["days"] = b555WeeklyRefreshDays
		b555WeeklyRefreshJSON, _ := json.Marshal(b555WeeklyRefreshInput)
		if _, b555WeeklyRefreshErr := B555WeeklyRefreshEvaluate(b555WeeklyRefreshJSON); b555WeeklyRefreshErr == nil {
			b555WeeklyRefreshT.Fatal("wrong day length accepted")
		}
	}
}

func TestAb07Ap07NamedRawAndShapeBounds(b555WeeklyRefreshT *testing.T) {
	for _, b555WeeklyRefreshView := range []string{"weekly", "nonheating-shape"} {
		b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
		b555WeeklyRefreshInput["fields"].(map[string]any)[b555WeeklyRefreshView] = []any{17, 0, 255}
		if b555WeeklyRefreshView == "nonheating-shape" {
			b555WeeklyRefreshInput["context"].(map[string]any)["operation"] = 166
			b555WeeklyRefreshInput["context"].(map[string]any)["view"] = b555WeeklyRefreshView
			b555WeeklyRefreshInput["arguments"] = map[string]any{"initial": []any{1, 2, 3, 4, 5}, "count": 0, "slot_setpoints": []any{}}
		}
		b555WeeklyRefreshOutput := b555WeeklyRefreshRun(b555WeeklyRefreshT, b555WeeklyRefreshInput)
		if !reflect.DeepEqual(b555WeeklyRefreshOutput["fields"].(map[string]any)[b555WeeklyRefreshView].(map[string]any)["raw"], []any{float64(17), float64(0), float64(255)}) {
			b555WeeklyRefreshT.Fatal("named raw overwritten")
		}
	}
	for _, b555WeeklyRefreshCount := range []int{0, 255, 256} {
		b555WeeklyRefreshInput := b555WeeklyRefreshWeek()
		b555WeeklyRefreshInput["context"].(map[string]any)["operation"] = 166
		b555WeeklyRefreshInput["context"].(map[string]any)["view"] = "nonheating-shape"
		b555WeeklyRefreshSetpoints := []any{}
		for b555WeeklyRefreshIndex := 0; b555WeeklyRefreshIndex < b555WeeklyRefreshCount; b555WeeklyRefreshIndex++ {
			b555WeeklyRefreshSetpoints = append(b555WeeklyRefreshSetpoints, b555WeeklyRefreshIndex)
		}
		b555WeeklyRefreshInput["arguments"] = map[string]any{"initial": []any{1, 2, 3, 4, 5}, "count": b555WeeklyRefreshCount, "slot_setpoints": b555WeeklyRefreshSetpoints}
		if b555WeeklyRefreshCount == 256 {
			b555WeeklyRefreshJSON, _ := json.Marshal(b555WeeklyRefreshInput)
			if _, b555WeeklyRefreshErr := B555WeeklyRefreshEvaluate(b555WeeklyRefreshJSON); b555WeeklyRefreshErr == nil {
				b555WeeklyRefreshT.Fatal("count overflow accepted")
			}
			continue
		}
		b555WeeklyRefreshResult := b555WeeklyRefreshRun(b555WeeklyRefreshT, b555WeeklyRefreshInput)["fields"].(map[string]any)["nonheating-shape"].(map[string]any)
		if len(b555WeeklyRefreshResult["processed_setpoints"].([]any)) != b555WeeklyRefreshCount {
			b555WeeklyRefreshT.Fatal("count boundary changed")
		}
	}
}
