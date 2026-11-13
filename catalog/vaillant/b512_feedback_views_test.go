package vaillant

import (
	b512FeedbackViewJSON "encoding/json"
	b512FeedbackViewOs "os"
	b512FeedbackViewReflect "reflect"
	b512FeedbackViewStrings "strings"
	b512FeedbackViewTesting "testing"
)

func b512FeedbackViewControl() map[string]any {
	return map[string]any{
		"context": map[string]any{"primary": 181, "secondary": 18, "operation": 1, "direction": "request", "view": "directional-fields", "family": map[string]any{"opaque": []any{"synthetic-family", 0}}, "variant": []any{"synthetic-variant", false}, "setpoint": []int{17}, "feedback": []int{17}},
		"fields":  map[string]any{"secondary_heating_pump": []int{11}, "circulation_pump": []int{22}, "solar_pump": []int{33}, "error_relay": []int{15}, "loading_pump": []int{44}, "fume_flap": []int{240}, "unmodelled": []int{0, 255, 19}},
	}
}

func b512FeedbackViewEncode(b512FeedbackViewT *b512FeedbackViewTesting.T, b512FeedbackViewValue any) b512FeedbackViewJSON.RawMessage {
	b512FeedbackViewT.Helper()
	b512FeedbackViewRaw, b512FeedbackViewErr := b512FeedbackViewJSON.Marshal(b512FeedbackViewValue)
	if b512FeedbackViewErr != nil {
		b512FeedbackViewT.Fatal(b512FeedbackViewErr)
	}
	return b512FeedbackViewRaw
}

func b512FeedbackViewDecode(b512FeedbackViewT *b512FeedbackViewTesting.T, b512FeedbackViewRaw b512FeedbackViewJSON.RawMessage) map[string]any {
	b512FeedbackViewT.Helper()
	var b512FeedbackViewValue map[string]any
	if b512FeedbackViewErr := b512FeedbackViewJSON.Unmarshal(b512FeedbackViewRaw, &b512FeedbackViewValue); b512FeedbackViewErr != nil {
		b512FeedbackViewT.Fatal(b512FeedbackViewErr)
	}
	return b512FeedbackViewValue
}

func b512FeedbackViewRun(b512FeedbackViewT *b512FeedbackViewTesting.T, b512FeedbackViewInput map[string]any) map[string]any {
	b512FeedbackViewT.Helper()
	b512FeedbackViewRaw := b512FeedbackViewEncode(b512FeedbackViewT, b512FeedbackViewInput)
	b512FeedbackViewResult, b512FeedbackViewErr := B512FeedbackViewInterpret(b512FeedbackViewRaw)
	if b512FeedbackViewErr != nil {
		b512FeedbackViewT.Fatal(b512FeedbackViewErr)
	}
	b512FeedbackViewOutput := b512FeedbackViewDecode(b512FeedbackViewT, b512FeedbackViewResult)
	b512FeedbackViewNormalized := b512FeedbackViewDecode(b512FeedbackViewT, b512FeedbackViewRaw)
	b512FeedbackViewWantKeys := 6
	if _, b512FeedbackViewPresent := b512FeedbackViewInput["arguments"]; b512FeedbackViewPresent {
		b512FeedbackViewWantKeys = 7
	}
	if len(b512FeedbackViewOutput) != b512FeedbackViewWantKeys {
		b512FeedbackViewT.Fatalf("unexpected output keys: %#v", b512FeedbackViewOutput)
	}
	_, b512FeedbackViewInputArguments := b512FeedbackViewNormalized["arguments"]
	_, b512FeedbackViewOutputArguments := b512FeedbackViewOutput["arguments"]
	if b512FeedbackViewInputArguments != b512FeedbackViewOutputArguments {
		b512FeedbackViewT.Fatal("argument presence changed")
	}

	for _, b512FeedbackViewKey := range []string{"context", "arguments"} {
		if !b512FeedbackViewReflect.DeepEqual(b512FeedbackViewOutput[b512FeedbackViewKey], b512FeedbackViewNormalized[b512FeedbackViewKey]) {
			b512FeedbackViewT.Fatalf("preservation %s: %#v", b512FeedbackViewKey, b512FeedbackViewOutput)
		}
	}
	if !b512FeedbackViewReflect.DeepEqual(b512FeedbackViewOutput["raw"], b512FeedbackViewNormalized["fields"]) || b512FeedbackViewOutput["physical"] != "Unknown" || b512FeedbackViewOutput["qualified"] != false {
		b512FeedbackViewT.Fatalf("raw/qualification: %#v", b512FeedbackViewOutput)
	}
	if !b512FeedbackViewReflect.DeepEqual(b512FeedbackViewOutput, b512FeedbackViewDecode(b512FeedbackViewT, b512FeedbackViewEncode(b512FeedbackViewT, b512FeedbackViewOutput))) {
		b512FeedbackViewT.Fatal("round trip")
	}
	return b512FeedbackViewOutput
}

func b512FeedbackViewExpected() map[string]any {
	return map[string]any{
		"secondary_heating_pump": map[string]any{"raw": []any{float64(11)}, "state": "Hypothesis", "flags": map[string]any{"origin": "origin-Unknown"}},
		"circulation_pump":       map[string]any{"raw": []any{float64(22)}, "state": "Hypothesis", "flags": map[string]any{"origin": "origin-Unknown"}},
		"solar_pump":             map[string]any{"raw": []any{float64(33)}, "state": "Hypothesis", "flags": map[string]any{"origin": "origin-Unknown"}},
		"error_relay":            map[string]any{"raw": []any{float64(15)}, "state": "Hypothesis", "label": "state-0F", "flags": map[string]any{"origin": "origin-Unknown"}},
		"loading_pump":           map[string]any{"raw": []any{float64(44)}, "state": "Hypothesis", "flags": map[string]any{"origin": "origin-Unknown"}},
		"fume_flap":              map[string]any{"raw": []any{float64(240)}, "state": "Hypothesis", "label": "state-F0", "flags": map[string]any{"origin": "origin-Unknown"}},
	}
}

func TestAb15Ao15IndependentFields(b512FeedbackViewT *b512FeedbackViewTesting.T) {
	for _, b512FeedbackViewName := range []string{"secondary_heating_pump", "circulation_pump", "solar_pump", "error_relay", "loading_pump", "fume_flap"} {
		for _, b512FeedbackViewDirection := range []string{"request", "positive"} {
			for _, b512FeedbackViewCase := range []struct {
				b512FeedbackViewByte                                                             int
				b512FeedbackViewPumpRequest, b512FeedbackViewPumpPositive, b512FeedbackViewRelay string
			}{
				{0, "", "", "no-action-candidate"}, {1, "", "", "no-action-candidate"}, {15, "", "", "state-0F"}, {240, "", "", "state-F0"}, {255, "leave-unchanged", "function-unavailable", "no-action-candidate"}, {42, "", "", "no-action-candidate"},
			} {
				b512FeedbackViewT.Run(b512FeedbackViewName+"/"+b512FeedbackViewDirection+"/"+string(b512FeedbackViewEncode(b512FeedbackViewT, b512FeedbackViewCase.b512FeedbackViewByte)), func(b512FeedbackViewT *b512FeedbackViewTesting.T) {
					b512FeedbackViewInput := b512FeedbackViewControl()
					b512FeedbackViewInput["context"].(map[string]any)["direction"] = b512FeedbackViewDirection
					b512FeedbackViewInput["fields"].(map[string]any)[b512FeedbackViewName] = []int{b512FeedbackViewCase.b512FeedbackViewByte}
					b512FeedbackViewWant := b512FeedbackViewExpected()
					b512FeedbackViewField := b512FeedbackViewWant[b512FeedbackViewName].(map[string]any)
					b512FeedbackViewField["raw"] = []any{float64(b512FeedbackViewCase.b512FeedbackViewByte)}
					b512FeedbackViewLabel := b512FeedbackViewCase.b512FeedbackViewPumpRequest
					if b512FeedbackViewDirection == "positive" {
						b512FeedbackViewLabel = b512FeedbackViewCase.b512FeedbackViewPumpPositive
					}
					if b512FeedbackViewName == "error_relay" || b512FeedbackViewName == "fume_flap" {
						b512FeedbackViewLabel = b512FeedbackViewCase.b512FeedbackViewRelay
					}
					delete(b512FeedbackViewField, "label")
					if b512FeedbackViewLabel != "" {
						b512FeedbackViewField["label"] = b512FeedbackViewLabel
					}
					b512FeedbackViewOutput := b512FeedbackViewRun(b512FeedbackViewT, b512FeedbackViewInput)
					if b512FeedbackViewOutput["state"] != "Hypothesis" || !b512FeedbackViewReflect.DeepEqual(b512FeedbackViewOutput["fields"], b512FeedbackViewWant) {
						b512FeedbackViewT.Fatalf("exact independent fields: got %#v want %#v", b512FeedbackViewOutput, b512FeedbackViewWant)
					}
				})
			}
		}
	}
}

func TestAb15Ao15Origin(b512FeedbackViewT *b512FeedbackViewTesting.T) {
	for _, b512FeedbackViewCase := range []struct {
		b512FeedbackViewID, b512FeedbackViewOrigin string
		b512FeedbackViewArguments                  any
	}{
		{"ab15ao15Absent", "origin-Unknown", nil}, {"ab15ao15Empty", "origin-Unknown", map[string]any{}}, {"ab15ao15False", "fallback-candidate", map[string]any{"actual_signal": false}}, {"ab15ao15True", "actual-input-candidate", map[string]any{"actual_signal": true}},
	} {
		for _, b512FeedbackViewFeedback := range []int{17, 93} {
			for _, b512FeedbackViewDirection := range []string{"request", "positive"} {
				b512FeedbackViewT.Run(b512FeedbackViewCase.b512FeedbackViewID+"/"+b512FeedbackViewDirection+"/"+string(b512FeedbackViewEncode(b512FeedbackViewT, b512FeedbackViewFeedback)), func(b512FeedbackViewT *b512FeedbackViewTesting.T) {
					b512FeedbackViewInput := b512FeedbackViewControl()
					b512FeedbackViewInput["context"].(map[string]any)["feedback"] = []int{b512FeedbackViewFeedback}
					b512FeedbackViewInput["context"].(map[string]any)["direction"] = b512FeedbackViewDirection
					if b512FeedbackViewCase.b512FeedbackViewArguments != nil {
						b512FeedbackViewInput["arguments"] = b512FeedbackViewCase.b512FeedbackViewArguments
					}
					b512FeedbackViewWant := b512FeedbackViewExpected()
					for _, b512FeedbackViewField := range b512FeedbackViewWant {
						b512FeedbackViewField.(map[string]any)["flags"] = map[string]any{"origin": b512FeedbackViewCase.b512FeedbackViewOrigin}
					}
					b512FeedbackViewOutput := b512FeedbackViewRun(b512FeedbackViewT, b512FeedbackViewInput)
					if !b512FeedbackViewReflect.DeepEqual(b512FeedbackViewOutput["fields"], b512FeedbackViewWant) {
						b512FeedbackViewT.Fatalf("origin/default/polarity: %#v", b512FeedbackViewOutput)
					}
				})
			}
		}
	}
}

func TestAb15Ao15Partial(b512FeedbackViewT *b512FeedbackViewTesting.T) {
	for b512FeedbackViewName := range b512FeedbackViewExpected() {
		for _, b512FeedbackViewWidth := range []string{"missing", "empty", "overlong", "zero"} {
			b512FeedbackViewT.Run(b512FeedbackViewName+"/"+b512FeedbackViewWidth, func(b512FeedbackViewT *b512FeedbackViewTesting.T) {
				b512FeedbackViewInput := b512FeedbackViewControl()
				b512FeedbackViewRaw := b512FeedbackViewInput["fields"].(map[string]any)
				b512FeedbackViewWant := b512FeedbackViewExpected()
				b512FeedbackViewState := "partial"
				b512FeedbackViewField := map[string]any{"raw": nil, "state": "Unknown"}
				switch b512FeedbackViewWidth {
				case "missing":
					delete(b512FeedbackViewRaw, b512FeedbackViewName)
				case "empty":
					b512FeedbackViewRaw[b512FeedbackViewName] = []int{}
					b512FeedbackViewField["raw"] = []any{}
				case "overlong":
					b512FeedbackViewRaw[b512FeedbackViewName] = []int{255, 15}
					b512FeedbackViewField["raw"] = []any{float64(255), float64(15)}
				case "zero":
					b512FeedbackViewRaw[b512FeedbackViewName] = []int{0}
					b512FeedbackViewState = "Hypothesis"
					b512FeedbackViewField = map[string]any{"raw": []any{float64(0)}, "state": "Hypothesis", "flags": map[string]any{"origin": "origin-Unknown"}}
					if b512FeedbackViewName == "error_relay" || b512FeedbackViewName == "fume_flap" {
						b512FeedbackViewField["label"] = "no-action-candidate"
					}
				}
				b512FeedbackViewWant[b512FeedbackViewName] = b512FeedbackViewField
				b512FeedbackViewOutput := b512FeedbackViewRun(b512FeedbackViewT, b512FeedbackViewInput)
				if b512FeedbackViewOutput["state"] != b512FeedbackViewState || !b512FeedbackViewReflect.DeepEqual(b512FeedbackViewOutput["fields"], b512FeedbackViewWant) {
					b512FeedbackViewT.Fatalf("partial siblings: %#v", b512FeedbackViewOutput)
				}
			})
		}
	}
}

func TestAb15Ao15SelectorsAndOpaqueContext(b512FeedbackViewT *b512FeedbackViewTesting.T) {
	for b512FeedbackViewKey, b512FeedbackViewValues := range map[string][]any{
		"primary": {0, 182}, "secondary": {0, 19}, "operation": {0, 2}, "direction": {"", "unknown"}, "view": {"", "unknown"},
	} {
		for _, b512FeedbackViewValue := range append(b512FeedbackViewValues, nil) {
			b512FeedbackViewT.Run(b512FeedbackViewKey+"/"+string(b512FeedbackViewEncode(b512FeedbackViewT, b512FeedbackViewValue)), func(b512FeedbackViewT *b512FeedbackViewTesting.T) {
				b512FeedbackViewInput := b512FeedbackViewControl()
				if b512FeedbackViewValue == nil {
					delete(b512FeedbackViewInput["context"].(map[string]any), b512FeedbackViewKey)
				} else {
					b512FeedbackViewInput["context"].(map[string]any)[b512FeedbackViewKey] = b512FeedbackViewValue
				}
				b512FeedbackViewOutput := b512FeedbackViewRun(b512FeedbackViewT, b512FeedbackViewInput)
				if b512FeedbackViewOutput["state"] != "Unknown" || len(b512FeedbackViewOutput["fields"].(map[string]any)) != 0 {
					b512FeedbackViewT.Fatalf("guessed selection: %#v", b512FeedbackViewOutput)
				}
			})
		}
	}
	for _, b512FeedbackViewKey := range []string{"family", "variant", "setpoint", "feedback", "opaque"} {
		for _, b512FeedbackViewValue := range []any{nil, 0, "directional-fields", []any{}, map[string]any{"nested": []any{false, 1}}} {
			b512FeedbackViewInput := b512FeedbackViewControl()
			b512FeedbackViewInput["context"].(map[string]any)[b512FeedbackViewKey] = b512FeedbackViewValue
			b512FeedbackViewOutput := b512FeedbackViewRun(b512FeedbackViewT, b512FeedbackViewInput)
			if !b512FeedbackViewReflect.DeepEqual(b512FeedbackViewOutput["fields"], b512FeedbackViewExpected()) {
				b512FeedbackViewT.Fatalf("opaque context affected selection: %#v", b512FeedbackViewOutput)
			}
		}
	}
}

func TestAb15Ao15Reject(b512FeedbackViewT *b512FeedbackViewTesting.T) {
	for _, b512FeedbackViewBad := range []string{"", "null", "[]", "{}", "{", `{"context":{},"context":{},"fields":{}}`, `{"context":{"family":{"x":1,"x":2}},"fields":{}}`, `{"context":{"family":[{"x":1,"\u0078":2}]},"fields":{}}`, `{"context":{},"fields":{"x":[],"x":[]}}`, `{"context":{},"fields":{},"arguments":{"actual_signal":true,"actual_signal":false}}`} {
		b512FeedbackViewResult, b512FeedbackViewErr := B512FeedbackViewInterpret([]byte(b512FeedbackViewBad))
		if b512FeedbackViewErr == nil || b512FeedbackViewResult != nil {
			b512FeedbackViewT.Fatalf("accepted %q", b512FeedbackViewBad)
		}
	}
	for b512FeedbackViewKey, b512FeedbackViewValues := range map[string][]any{
		"primary": {nil, "181", 1.5, true, []int{}, map[string]any{}}, "secondary": {nil, "18", false}, "operation": {nil, "1", 1.5}, "direction": {nil, 1, []int{}}, "view": {nil, true, map[string]any{}},
	} {
		for _, b512FeedbackViewValue := range b512FeedbackViewValues {
			b512FeedbackViewInput := b512FeedbackViewControl()
			b512FeedbackViewInput["context"].(map[string]any)[b512FeedbackViewKey] = b512FeedbackViewValue
			b512FeedbackViewResult, b512FeedbackViewErr := B512FeedbackViewInterpret(b512FeedbackViewEncode(b512FeedbackViewT, b512FeedbackViewInput))
			if b512FeedbackViewErr == nil || b512FeedbackViewResult != nil {
				b512FeedbackViewT.Fatalf("accepted identity %s %#v", b512FeedbackViewKey, b512FeedbackViewValue)
			}
		}
	}
	for _, b512FeedbackViewValue := range []any{nil, true, 0, "", map[string]any{}, []any{-1}, []any{256}, []any{1.5}, []any{nil}, []any{true}, []any{"1"}, []any{[]int{1}}, []any{map[string]any{}}} {
		for _, b512FeedbackViewName := range []string{"secondary_heating_pump", "circulation_pump", "solar_pump", "error_relay", "loading_pump", "fume_flap", "unmodelled"} {
			b512FeedbackViewInput := b512FeedbackViewControl()
			b512FeedbackViewInput["fields"].(map[string]any)[b512FeedbackViewName] = b512FeedbackViewValue
			b512FeedbackViewResult, b512FeedbackViewErr := B512FeedbackViewInterpret(b512FeedbackViewEncode(b512FeedbackViewT, b512FeedbackViewInput))
			if b512FeedbackViewErr == nil || b512FeedbackViewResult != nil {
				b512FeedbackViewT.Fatalf("accepted bytes %#v", b512FeedbackViewValue)
			}
		}
	}
	for _, b512FeedbackViewKey := range []string{"context", "fields", "arguments"} {
		for _, b512FeedbackViewValue := range []any{nil, []int{}, 0, "", true} {
			b512FeedbackViewInput := b512FeedbackViewControl()
			b512FeedbackViewInput[b512FeedbackViewKey] = b512FeedbackViewValue
			b512FeedbackViewResult, b512FeedbackViewErr := B512FeedbackViewInterpret(b512FeedbackViewEncode(b512FeedbackViewT, b512FeedbackViewInput))
			if b512FeedbackViewErr == nil || b512FeedbackViewResult != nil {
				b512FeedbackViewT.Fatalf("accepted kind %s", b512FeedbackViewKey)
			}
		}
	}
	for _, b512FeedbackViewValue := range []any{nil, 0, 1, "true", []int{}, map[string]any{}} {
		b512FeedbackViewInput := b512FeedbackViewControl()
		b512FeedbackViewInput["arguments"] = map[string]any{"actual_signal": b512FeedbackViewValue}
		b512FeedbackViewResult, b512FeedbackViewErr := B512FeedbackViewInterpret(b512FeedbackViewEncode(b512FeedbackViewT, b512FeedbackViewInput))
		if b512FeedbackViewErr == nil || b512FeedbackViewResult != nil {
			b512FeedbackViewT.Fatalf("accepted flag %#v", b512FeedbackViewValue)
		}
	}
	for _, b512FeedbackViewBad := range []string{
		string(b512FeedbackViewEncode(b512FeedbackViewT, b512FeedbackViewControl())) + ` {}`, b512FeedbackViewStrings.Replace(string(b512FeedbackViewEncode(b512FeedbackViewT, b512FeedbackViewControl())), `[11]`, `[1e0]`, 1),
	} {
		b512FeedbackViewResult, b512FeedbackViewErr := B512FeedbackViewInterpret([]byte(b512FeedbackViewBad))
		if b512FeedbackViewErr == nil || b512FeedbackViewResult != nil {
			b512FeedbackViewT.Fatalf("accepted malformed/integer syntax %s", b512FeedbackViewBad)
		}
	}
	b512FeedbackViewInput := b512FeedbackViewControl()
	b512FeedbackViewInput["arguments"] = map[string]any{"default": true}
	b512FeedbackViewResult, b512FeedbackViewErr := B512FeedbackViewInterpret(b512FeedbackViewEncode(b512FeedbackViewT, b512FeedbackViewInput))
	if b512FeedbackViewErr == nil || b512FeedbackViewResult != nil {
		b512FeedbackViewT.Fatal("accepted invented default")
	}
}

func TestAb15Ao15Isolation(b512FeedbackViewT *b512FeedbackViewTesting.T) {
	b512FeedbackViewInput := b512FeedbackViewControl()
	b512FeedbackViewInput["arguments"] = map[string]any{"actual_signal": false}
	b512FeedbackViewRaw := b512FeedbackViewEncode(b512FeedbackViewT, b512FeedbackViewInput)
	b512FeedbackViewFirst, b512FeedbackViewErr := B512FeedbackViewInterpret(b512FeedbackViewRaw)
	if b512FeedbackViewErr != nil {
		b512FeedbackViewT.Fatal(b512FeedbackViewErr)
	}
	b512FeedbackViewSaved := append([]byte(nil), b512FeedbackViewFirst...)
	for b512FeedbackViewIndex := range b512FeedbackViewRaw {
		b512FeedbackViewRaw[b512FeedbackViewIndex] = ' '
	}
	if string(b512FeedbackViewFirst) != string(b512FeedbackViewSaved) {
		b512FeedbackViewT.Fatal("input aliases output")
	}
	b512FeedbackViewNested := b512FeedbackViewDecode(b512FeedbackViewT, b512FeedbackViewFirst)
	b512FeedbackViewNested["raw"].(map[string]any)["solar_pump"].([]any)[0] = float64(99)
	if b512FeedbackViewNested["fields"].(map[string]any)["solar_pump"].(map[string]any)["raw"].([]any)[0] != float64(33) {
		b512FeedbackViewT.Fatal("raw aliases semantic raw")
	}
	b512FeedbackViewNested["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "changed"
	b512FeedbackViewNested["arguments"].(map[string]any)["actual_signal"] = true
	b512FeedbackViewNested["fields"].(map[string]any)["circulation_pump"].(map[string]any)["flags"].(map[string]any)["origin"] = "changed"
	b512FeedbackViewSecond, b512FeedbackViewErr := B512FeedbackViewInterpret(b512FeedbackViewEncode(b512FeedbackViewT, b512FeedbackViewInput))
	if b512FeedbackViewErr != nil {
		b512FeedbackViewT.Fatal(b512FeedbackViewErr)
	}
	if string(b512FeedbackViewSecond) != string(b512FeedbackViewSaved) {
		b512FeedbackViewT.Fatal("nested mutation affected later call")
	}
	for b512FeedbackViewIndex := range b512FeedbackViewFirst {
		b512FeedbackViewFirst[b512FeedbackViewIndex] = ' '
	}
	if string(b512FeedbackViewSecond) != string(b512FeedbackViewSaved) {
		b512FeedbackViewT.Fatal("output aliases another call")
	}
}

func TestAb15Ao15Fixtures(b512FeedbackViewT *b512FeedbackViewTesting.T) {
	b512FeedbackViewData, b512FeedbackViewErr := b512FeedbackViewOs.ReadFile("testdata/ab15ao15_b512_feedback_views.json")
	if b512FeedbackViewErr != nil {
		b512FeedbackViewT.Fatal(b512FeedbackViewErr)
	}
	var b512FeedbackViewCases []struct {
		B512FeedbackViewID     string            `json:"id"`
		B512FeedbackViewInput  map[string]any    `json:"input"`
		B512FeedbackViewState  string            `json:"state"`
		B512FeedbackViewLabels map[string]string `json:"labels"`
	}
	if b512FeedbackViewErr := b512FeedbackViewJSON.Unmarshal(b512FeedbackViewData, &b512FeedbackViewCases); b512FeedbackViewErr != nil {
		b512FeedbackViewT.Fatal(b512FeedbackViewErr)
	}
	if len(b512FeedbackViewCases) != 3 {
		b512FeedbackViewT.Fatal("fixture count")
	}
	for _, b512FeedbackViewCase := range b512FeedbackViewCases {
		b512FeedbackViewT.Run(b512FeedbackViewCase.B512FeedbackViewID, func(b512FeedbackViewT *b512FeedbackViewTesting.T) {
			if !b512FeedbackViewStrings.HasPrefix(b512FeedbackViewCase.B512FeedbackViewID, "ab15ao15") {
				b512FeedbackViewT.Fatal("fixture ID")
			}
			b512FeedbackViewOutput := b512FeedbackViewRun(b512FeedbackViewT, b512FeedbackViewCase.B512FeedbackViewInput)
			if b512FeedbackViewOutput["state"] != b512FeedbackViewCase.B512FeedbackViewState {
				b512FeedbackViewT.Fatalf("fixture state %#v", b512FeedbackViewOutput)
			}
			for b512FeedbackViewName, b512FeedbackViewLabel := range b512FeedbackViewCase.B512FeedbackViewLabels {
				if b512FeedbackViewOutput["fields"].(map[string]any)[b512FeedbackViewName].(map[string]any)["label"] != b512FeedbackViewLabel {
					b512FeedbackViewT.Fatalf("fixture label %s", b512FeedbackViewName)
				}
			}
		})
	}
}
