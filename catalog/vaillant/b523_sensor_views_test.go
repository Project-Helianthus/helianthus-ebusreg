package vaillant

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func b523SensorViewInput(b523SensorViewOperation int, b523SensorViewBytes []int) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 35, "operation": b523SensorViewOperation, "direction": "positive", "view": "status-zero-valid", "family": map[string]any{"opaque": []any{"synthetic-family", 0}}, "variant": []any{"synthetic-variant", nil}}, "fields": map[string]any{"logical_object": b523SensorViewBytes, "unmodelled": []int{0, 255, 7}}}
}

func b523SensorViewCall(b523SensorViewT *testing.T, b523SensorViewInputValue any) map[string]any {
	b523SensorViewT.Helper()
	b523SensorViewInputJSON, b523SensorViewErr := json.Marshal(b523SensorViewInputValue)
	if b523SensorViewErr != nil {
		b523SensorViewT.Fatal(b523SensorViewErr)
	}
	b523SensorViewResult, b523SensorViewErr := B523SensorViewInterpret(b523SensorViewInputJSON)
	if b523SensorViewErr != nil {
		b523SensorViewT.Fatal(b523SensorViewErr)
	}
	var b523SensorViewOutput map[string]any
	if b523SensorViewErr := json.Unmarshal(b523SensorViewResult, &b523SensorViewOutput); b523SensorViewErr != nil {
		b523SensorViewT.Fatal(b523SensorViewErr)
	}
	if b523SensorViewOutput["physical"] != "Unknown" || b523SensorViewOutput["qualified"] != false {
		b523SensorViewT.Fatal("qualification inferred")
	}
	var b523SensorViewRoundTrip map[string]any
	if b523SensorViewErr := json.Unmarshal(b523SensorViewInputJSON, &b523SensorViewRoundTrip); b523SensorViewErr != nil {
		b523SensorViewT.Fatal(b523SensorViewErr)
	}
	if !reflect.DeepEqual(b523SensorViewOutput["context"], b523SensorViewRoundTrip["context"]) || !reflect.DeepEqual(b523SensorViewOutput["raw"], b523SensorViewRoundTrip["fields"]) {
		b523SensorViewT.Fatal("native evidence lost")
	}
	if b523SensorViewExpected, b523SensorViewPresent := b523SensorViewRoundTrip["arguments"]; b523SensorViewPresent {
		if !reflect.DeepEqual(b523SensorViewOutput["arguments"], b523SensorViewExpected) {
			b523SensorViewT.Fatal("arguments lost")
		}
	} else if _, b523SensorViewPresent := b523SensorViewOutput["arguments"]; b523SensorViewPresent {
		b523SensorViewT.Fatal("arguments defaulted")
	}
	return b523SensorViewOutput
}

func b523SensorViewField(b523SensorViewOutput map[string]any, b523SensorViewName string) map[string]any {
	return b523SensorViewOutput["fields"].(map[string]any)[b523SensorViewName].(map[string]any)
}

func TestAb27Ao27Fixture(b523SensorViewT *testing.T) {
	b523SensorViewData, b523SensorViewErr := os.ReadFile("testdata/ab27ao27_b523_sensor_views.json")
	if b523SensorViewErr != nil {
		b523SensorViewT.Fatal(b523SensorViewErr)
	}
	var b523SensorViewCases []struct {
		B523SensorViewID       string                    `json:"id"`
		B523SensorViewInput    map[string]any            `json:"input"`
		B523SensorViewExpected map[string]map[string]any `json:"expected"`
	}
	if b523SensorViewErr := json.Unmarshal(b523SensorViewData, &b523SensorViewCases); b523SensorViewErr != nil {
		b523SensorViewT.Fatal(b523SensorViewErr)
	}
	if len(b523SensorViewCases) != 3 {
		b523SensorViewT.Fatal("fixture count")
	}
	for _, b523SensorViewCase := range b523SensorViewCases {
		b523SensorViewT.Run(b523SensorViewCase.B523SensorViewID, func(b523SensorViewT *testing.T) {
			b523SensorViewOutput := b523SensorViewCall(b523SensorViewT, b523SensorViewCase.B523SensorViewInput)
			for b523SensorViewName, b523SensorViewExpected := range b523SensorViewCase.B523SensorViewExpected {
				for b523SensorViewKey, b523SensorViewValue := range b523SensorViewExpected {
					if !reflect.DeepEqual(b523SensorViewField(b523SensorViewOutput, b523SensorViewName)[b523SensorViewKey], b523SensorViewValue) {
						b523SensorViewT.Fatalf("%s %s", b523SensorViewName, b523SensorViewKey)
					}
				}
			}
		})
	}
}

func TestAb27Ao27Sensors(b523SensorViewT *testing.T) {
	for _, b523SensorViewOperation := range []int{3, 6, 7} {
		b523SensorViewPositions := []int{3, 5, 7, 9, 11, 13}
		b523SensorViewStart := 1
		b523SensorViewExtent := 18
		if b523SensorViewOperation == 6 {
			b523SensorViewPositions = []int{3, 5, 7, 9, 11, 13, 15}
			b523SensorViewExtent = 19
		}
		if b523SensorViewOperation == 7 {
			b523SensorViewPositions = []int{3, 5, 7, 9, 12, 14}
			b523SensorViewStart = 8
		}
		for b523SensorViewIndex, b523SensorViewPosition := range b523SensorViewPositions {
			for _, b523SensorViewVector := range []struct {
				b523SensorViewLow, b523SensorViewHigh int
				b523SensorViewNumber                  float64
				b523SensorViewReplacement             bool
			}{{0, 128, 0, true}, {16, 0, 1, false}, {240, 255, -1, false}, {255, 127, 2047.9375, false}, {1, 128, -2047.9375, false}} {
				b523SensorViewT.Run(fmt.Sprintf("op%d_S%d_%d_%d", b523SensorViewOperation, b523SensorViewStart+b523SensorViewIndex, b523SensorViewVector.b523SensorViewLow, b523SensorViewVector.b523SensorViewHigh), func(b523SensorViewT *testing.T) {
					b523SensorViewBytes := make([]int, b523SensorViewExtent+2)
					b523SensorViewBytes[0], b523SensorViewBytes[1], b523SensorViewBytes[2] = 255, 99, 0
					b523SensorViewBytes[b523SensorViewPosition], b523SensorViewBytes[b523SensorViewPosition+1] = b523SensorViewVector.b523SensorViewLow, b523SensorViewVector.b523SensorViewHigh
					if b523SensorViewOperation == 7 {
						b523SensorViewBytes[11] = 77
					}
					b523SensorViewBytes[b523SensorViewExtent], b523SensorViewBytes[b523SensorViewExtent+1] = 222, 111
					b523SensorViewOutput := b523SensorViewCall(b523SensorViewT, b523SensorViewInput(b523SensorViewOperation, b523SensorViewBytes))
					if b523SensorViewOutput["state"] != "Hypothesis" {
						b523SensorViewT.Fatal("complete state")
					}
					for b523SensorViewSibling := range b523SensorViewPositions {
						b523SensorViewFieldValue := b523SensorViewField(b523SensorViewOutput, fmt.Sprintf("S%d", b523SensorViewStart+b523SensorViewSibling))
						if b523SensorViewFieldValue["state"] != "Hypothesis" || b523SensorViewFieldValue["grade"] != "Hypothesis" {
							b523SensorViewT.Fatal("grade")
						}
						if b523SensorViewSibling == b523SensorViewIndex && b523SensorViewVector.b523SensorViewReplacement {
							if b523SensorViewFieldValue["raw"] != float64(32768) || b523SensorViewFieldValue["label"] != "replacement" {
								b523SensorViewT.Fatal("replacement selector")
							}
							if _, b523SensorViewPresent := b523SensorViewFieldValue["number"]; b523SensorViewPresent {
								b523SensorViewT.Fatal("replacement converted")
							}
							if _, b523SensorViewPresent := b523SensorViewFieldValue["calculation"]; b523SensorViewPresent {
								b523SensorViewT.Fatal("replacement calculation")
							}
						} else {
							b523SensorViewWant := float64(0)
							if b523SensorViewSibling == b523SensorViewIndex {
								b523SensorViewWant = b523SensorViewVector.b523SensorViewNumber
							}
							if b523SensorViewFieldValue["number"] != b523SensorViewWant || b523SensorViewFieldValue["calculation"] != "Proven" {
								b523SensorViewT.Fatalf("sensor %d: %#v want %v", b523SensorViewSibling, b523SensorViewFieldValue, b523SensorViewWant)
							}
						}
					}
					if b523SensorViewOperation == 7 {
						b523SensorViewOpaque := b523SensorViewField(b523SensorViewOutput, "opaque11")
						if b523SensorViewOpaque["raw"] != float64(77) || b523SensorViewOpaque["state"] != "Unknown" {
							b523SensorViewT.Fatal("opaque position")
						}
						if _, b523SensorViewPresent := b523SensorViewOpaque["number"]; b523SensorViewPresent {
							b523SensorViewT.Fatal("opaque numeric guess")
						}
					}
				})
			}
		}
	}
}

func TestAb27Ao27StatusPairs(b523SensorViewT *testing.T) {
	for _, b523SensorViewOperation := range []int{3, 6, 7} {
		b523SensorViewCount, b523SensorViewStart, b523SensorViewStatus, b523SensorViewExtent := 6, 1, 16, 18
		if b523SensorViewOperation == 6 {
			b523SensorViewCount, b523SensorViewStatus, b523SensorViewExtent = 7, 17, 19
		}
		if b523SensorViewOperation == 7 {
			b523SensorViewStart = 8
		}
		for _, b523SensorViewView := range []string{"status-zero-valid", "status-zero-valid-or-low"} {
			for b523SensorViewIndex := 0; b523SensorViewIndex < b523SensorViewCount; b523SensorViewIndex++ {
				for b523SensorViewPair := 0; b523SensorViewPair < 4; b523SensorViewPair++ {
					b523SensorViewT.Run(fmt.Sprintf("op%d_%s_pair%d_value%d", b523SensorViewOperation, b523SensorViewView, b523SensorViewIndex, b523SensorViewPair), func(b523SensorViewT *testing.T) {
						b523SensorViewBytes := make([]int, b523SensorViewExtent)
						b523SensorViewBytes[b523SensorViewStatus+b523SensorViewIndex/4] = b523SensorViewPair << (2 * (b523SensorViewIndex % 4))
						b523SensorViewInputValue := b523SensorViewInput(b523SensorViewOperation, b523SensorViewBytes)
						b523SensorViewInputValue["context"].(map[string]any)["view"] = b523SensorViewView
						b523SensorViewOutput := b523SensorViewCall(b523SensorViewT, b523SensorViewInputValue)
						for b523SensorViewSibling := 0; b523SensorViewSibling < b523SensorViewCount; b523SensorViewSibling++ {
							b523SensorViewExpectedPair := 0
							if b523SensorViewSibling == b523SensorViewIndex {
								b523SensorViewExpectedPair = b523SensorViewPair
							}
							b523SensorViewLabels := []string{"valid-value", "open-circuit", "closed-circuit", "high-level"}
							if b523SensorViewView == "status-zero-valid-or-low" {
								b523SensorViewLabels[0] = "valid-value-or-low-level"
							}
							b523SensorViewFieldValue := b523SensorViewField(b523SensorViewOutput, fmt.Sprintf("statusS%d", b523SensorViewStart+b523SensorViewSibling))
							if b523SensorViewFieldValue["raw"] != float64(b523SensorViewExpectedPair) || b523SensorViewFieldValue["label"] != b523SensorViewLabels[b523SensorViewExpectedPair] || b523SensorViewFieldValue["grade"] != "Hypothesis" {
								b523SensorViewT.Fatal(b523SensorViewFieldValue)
							}
						}
					})
				}
			}
			b523SensorViewFirstReserved := 4
			if b523SensorViewOperation == 6 {
				b523SensorViewFirstReserved = 6
			}
			for b523SensorViewBit := b523SensorViewFirstReserved; b523SensorViewBit < 8; b523SensorViewBit++ {
				b523SensorViewT.Run(fmt.Sprintf("op%d_%s_reserved%d", b523SensorViewOperation, b523SensorViewView, b523SensorViewBit), func(b523SensorViewT *testing.T) {
					b523SensorViewBytes := make([]int, b523SensorViewExtent)
					b523SensorViewBytes[b523SensorViewStatus+1] = 1 << b523SensorViewBit
					b523SensorViewInputValue := b523SensorViewInput(b523SensorViewOperation, b523SensorViewBytes)
					b523SensorViewInputValue["context"].(map[string]any)["view"] = b523SensorViewView
					b523SensorViewOutput := b523SensorViewCall(b523SensorViewT, b523SensorViewInputValue)
					if b523SensorViewField(b523SensorViewOutput, "status_reserved")["reserved"] != float64(int(1)<<(b523SensorViewBit-b523SensorViewFirstReserved)) {
						b523SensorViewT.Fatal("reserved bit lost")
					}
					for b523SensorViewIndex := 0; b523SensorViewIndex < b523SensorViewCount; b523SensorViewIndex++ {
						if b523SensorViewField(b523SensorViewOutput, fmt.Sprintf("statusS%d", b523SensorViewStart+b523SensorViewIndex))["raw"] != float64(0) {
							b523SensorViewT.Fatal("reserved interpreted as status")
						}
					}
				})
			}
		}
	}
}

func TestAb27Ao27TruncationPWM(b523SensorViewT *testing.T) {
	for _, b523SensorViewOperation := range []int{3, 6, 7} {
		b523SensorViewExtent := 18
		b523SensorViewPositions := []int{3, 5, 7, 9, 11, 13}
		b523SensorViewStart := 1
		if b523SensorViewOperation == 6 {
			b523SensorViewExtent = 19
			b523SensorViewPositions = []int{3, 5, 7, 9, 11, 13, 15}
		}
		if b523SensorViewOperation == 7 {
			b523SensorViewPositions = []int{3, 5, 7, 9, 12, 14}
			b523SensorViewStart = 8
		}
		for b523SensorViewWidth := 0; b523SensorViewWidth <= b523SensorViewExtent+2; b523SensorViewWidth++ {
			b523SensorViewT.Run(fmt.Sprintf("op%d_width%d", b523SensorViewOperation, b523SensorViewWidth), func(b523SensorViewT *testing.T) {
				b523SensorViewBytes := make([]int, b523SensorViewWidth)
				for b523SensorViewIndex := range b523SensorViewBytes {
					b523SensorViewBytes[b523SensorViewIndex] = 255
				}
				b523SensorViewOutput := b523SensorViewCall(b523SensorViewT, b523SensorViewInput(b523SensorViewOperation, b523SensorViewBytes))
				b523SensorViewState := "partial"
				if b523SensorViewWidth >= b523SensorViewExtent {
					b523SensorViewState = "Hypothesis"
				}
				if b523SensorViewOutput["state"] != b523SensorViewState {
					b523SensorViewT.Fatal("extent state")
				}
				for b523SensorViewIndex, b523SensorViewPosition := range b523SensorViewPositions {
					b523SensorViewFieldValue := b523SensorViewField(b523SensorViewOutput, fmt.Sprintf("S%d", b523SensorViewStart+b523SensorViewIndex))
					if b523SensorViewPosition+2 <= b523SensorViewWidth {
						if b523SensorViewFieldValue["number"] != -0.0625 {
							b523SensorViewT.Fatal("complete sibling discarded")
						}
					} else {
						if b523SensorViewFieldValue["state"] != "incomplete" {
							b523SensorViewT.Fatal("incomplete state")
						}
						if _, b523SensorViewPresent := b523SensorViewFieldValue["number"]; b523SensorViewPresent {
							b523SensorViewT.Fatal("missing word defaulted")
						}
					}
				}
				if b523SensorViewOperation == 3 && b523SensorViewWidth >= 16 {
					if b523SensorViewField(b523SensorViewOutput, "PWM")["number"] != 127.5 {
						b523SensorViewT.Fatal("PWM clipped")
					}
				}
			})
		}
	}
}

func TestAb27Ao27SelectionPresence(b523SensorViewT *testing.T) {
	for _, b523SensorViewKey := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
		b523SensorViewVariants := []any{nil, 0, 182}
		if b523SensorViewKey == "secondary" {
			b523SensorViewVariants = []any{nil, 0, 36}
		}
		if b523SensorViewKey == "operation" {
			b523SensorViewVariants = []any{nil, 0, 8}
		}
		if b523SensorViewKey == "direction" || b523SensorViewKey == "view" {
			b523SensorViewVariants = []any{nil, "", "unknown"}
		}
		if b523SensorViewKey == "family" || b523SensorViewKey == "variant" {
			b523SensorViewVariants = []any{nil, "", map[string]any{"different": true}}
		}
		for b523SensorViewIndex, b523SensorViewValue := range b523SensorViewVariants {
			b523SensorViewT.Run(fmt.Sprintf("%s_%d", b523SensorViewKey, b523SensorViewIndex), func(b523SensorViewT *testing.T) {
				b523SensorViewInputValue := b523SensorViewInput(3, make([]int, 18))
				b523SensorViewContext := b523SensorViewInputValue["context"].(map[string]any)
				if b523SensorViewValue == nil {
					delete(b523SensorViewContext, b523SensorViewKey)
				} else {
					b523SensorViewContext[b523SensorViewKey] = b523SensorViewValue
				}
				b523SensorViewOutput := b523SensorViewCall(b523SensorViewT, b523SensorViewInputValue)
				if b523SensorViewKey == "family" || b523SensorViewKey == "variant" {
					if b523SensorViewOutput["state"] != "Hypothesis" {
						b523SensorViewT.Fatal("opaque context selects")
					}
				} else {
					if b523SensorViewOutput["state"] != "Unknown" || len(b523SensorViewOutput["fields"].(map[string]any)) != 0 {
						b523SensorViewT.Fatal("guessed selection")
					}
				}
			})
		}
	}
	for _, b523SensorViewKey := range []string{"logical_object", "unmodelled"} {
		for _, b523SensorViewMode := range []string{"missing", "empty", "zero", "changed"} {
			b523SensorViewT.Run(b523SensorViewKey+"_"+b523SensorViewMode, func(b523SensorViewT *testing.T) {
				b523SensorViewInputValue := b523SensorViewInput(3, make([]int, 18))
				b523SensorViewRaw := b523SensorViewInputValue["fields"].(map[string]any)
				switch b523SensorViewMode {
				case "missing":
					delete(b523SensorViewRaw, b523SensorViewKey)
				case "empty":
					b523SensorViewRaw[b523SensorViewKey] = []int{}
				case "zero":
					b523SensorViewRaw[b523SensorViewKey] = []int{0}
				case "changed":
					b523SensorViewRaw[b523SensorViewKey] = []int{99, 1, 7, 16, 0}
				}
				b523SensorViewOutput := b523SensorViewCall(b523SensorViewT, b523SensorViewInputValue)
				if b523SensorViewKey == "unmodelled" && b523SensorViewOutput["state"] != "Hypothesis" {
					b523SensorViewT.Fatal("unmodelled field selected")
				}
			})
		}
	}
	b523SensorViewInputValue := b523SensorViewInput(3, make([]int, 18))
	b523SensorViewInputValue["arguments"] = map[string]any{}
	b523SensorViewCall(b523SensorViewT, b523SensorViewInputValue)
}

func TestAb27Ao27InvalidJSON(b523SensorViewT *testing.T) {
	for b523SensorViewIndex, b523SensorViewInputJSON := range []string{
		``, `null`, `[]`, `{`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"family":{"x":1,"x":2}},"fields":{}}`, `{"context":{},"fields":{"x":[],"x":[]}}`, `{"context":{},"fields":{},"arguments":{"x":1,"x":2}}`, `{"context":{},"fields":{},"fields":{}}`, `{"context":{},"fields":{}} true`, `{"context":[],"fields":{}}`, `{"context":{},"fields":null}`, `{"context":{},"fields":[]}`, `{"context":{},"fields":{"x":null}}`, `{"context":{},"fields":{"x":"AA=="}}`, `{"context":{},"fields":{"x":[-1]}}`, `{"context":{},"fields":{"x":[256]}}`, `{"context":{},"fields":{"x":[1.5]}}`, `{"context":{},"fields":{"x":[1e0]}}`, `{"context":{},"fields":{"x":[true]}}`, `{"context":{},"fields":{"x":[null]}}`, `{"context":{},"fields":{"x":["1"]}}`, `{"context":{},"fields":{"x":[{}]}}`, `{"context":{},"fields":{"x":[[]]}}`, `{"context":{},"fields":{},"arguments":null}`, `{"context":{},"fields":{},"arguments":[]}`, `{"context":{},"fields":{},"arguments":{"flag":true}}`, `{"context":{},"fields":{},"extra":0}`,
	} {
		b523SensorViewT.Run(fmt.Sprint(b523SensorViewIndex), func(b523SensorViewT *testing.T) {
			b523SensorViewResult, b523SensorViewErr := B523SensorViewInterpret(json.RawMessage(b523SensorViewInputJSON))
			if b523SensorViewErr == nil || b523SensorViewResult != nil {
				b523SensorViewT.Fatal("invalid JSON admitted")
			}
		})
	}
	for _, b523SensorViewKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, b523SensorViewValue := range []any{nil, true, []any{}, map[string]any{}, 1.5} {
			b523SensorViewInputValue := b523SensorViewInput(3, make([]int, 18))
			b523SensorViewInputValue["context"].(map[string]any)[b523SensorViewKey] = b523SensorViewValue
			b523SensorViewInputJSON, b523SensorViewErr := json.Marshal(b523SensorViewInputValue)
			if b523SensorViewErr != nil {
				b523SensorViewT.Fatal(b523SensorViewErr)
			}
			b523SensorViewResult, b523SensorViewErr := B523SensorViewInterpret(b523SensorViewInputJSON)
			if b523SensorViewErr == nil || b523SensorViewResult != nil {
				b523SensorViewT.Fatal("wrong context kind admitted", b523SensorViewKey)
			}
		}
	}
}

func TestAb27Ao27Ownership(b523SensorViewT *testing.T) {
	b523SensorViewInputValue := b523SensorViewInput(7, make([]int, 18))
	b523SensorViewInputValue["arguments"] = map[string]any{}
	b523SensorViewInputJSON, b523SensorViewErr := json.Marshal(b523SensorViewInputValue)
	if b523SensorViewErr != nil {
		b523SensorViewT.Fatal(b523SensorViewErr)
	}
	b523SensorViewControl := append(json.RawMessage(nil), b523SensorViewInputJSON...)
	b523SensorViewFirst, b523SensorViewErr := B523SensorViewInterpret(b523SensorViewInputJSON)
	if b523SensorViewErr != nil {
		b523SensorViewT.Fatal(b523SensorViewErr)
	}
	b523SensorViewSaved := append(json.RawMessage(nil), b523SensorViewFirst...)
	for b523SensorViewIndex := range b523SensorViewInputJSON {
		b523SensorViewInputJSON[b523SensorViewIndex] = ' '
	}
	if !reflect.DeepEqual(b523SensorViewFirst, b523SensorViewSaved) {
		b523SensorViewT.Fatal("input aliases output")
	}
	b523SensorViewSecond, b523SensorViewErr := B523SensorViewInterpret(b523SensorViewControl)
	if b523SensorViewErr != nil {
		b523SensorViewT.Fatal(b523SensorViewErr)
	}
	for b523SensorViewIndex := range b523SensorViewFirst {
		b523SensorViewFirst[b523SensorViewIndex] = ' '
	}
	if !reflect.DeepEqual(b523SensorViewSecond, b523SensorViewSaved) {
		b523SensorViewT.Fatal("output aliases another result")
	}
	var b523SensorViewNested map[string]any
	if b523SensorViewErr := json.Unmarshal(b523SensorViewSecond, &b523SensorViewNested); b523SensorViewErr != nil {
		b523SensorViewT.Fatal(b523SensorViewErr)
	}
	b523SensorViewNested["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "changed"
	b523SensorViewNested["raw"].(map[string]any)["logical_object"].([]any)[0] = float64(99)
	b523SensorViewNested["arguments"].(map[string]any)["changed"] = true
	b523SensorViewField(b523SensorViewNested, "S8")["number"] = float64(999)
	b523SensorViewThird, b523SensorViewErr := B523SensorViewInterpret(b523SensorViewControl)
	if b523SensorViewErr != nil || !reflect.DeepEqual(b523SensorViewThird, b523SensorViewSaved) {
		b523SensorViewT.Fatal("nested mutation affects later call")
	}
}

func TestAb27Ao27OpaquePrefixAndPWM(b523SensorViewT *testing.T) {
	for b523SensorViewPosition := 0; b523SensorViewPosition < 3; b523SensorViewPosition++ {
		for _, b523SensorViewByte := range []int{0, 3, 6, 7, 181, 255} {
			b523SensorViewT.Run(fmt.Sprintf("prefix%d_%d", b523SensorViewPosition, b523SensorViewByte), func(b523SensorViewT *testing.T) {
				b523SensorViewBytes := make([]int, 18)
				b523SensorViewBytes[b523SensorViewPosition] = b523SensorViewByte
				b523SensorViewOutput := b523SensorViewCall(b523SensorViewT, b523SensorViewInput(7, b523SensorViewBytes))
				if b523SensorViewOutput["state"] != "Hypothesis" || b523SensorViewField(b523SensorViewOutput, "S12")["number"] != float64(0) {
					b523SensorViewT.Fatal("prefix selected operation")
				}
			})
		}
	}
	for _, b523SensorViewCase := range []struct {
		b523SensorViewByte int
		b523SensorViewWant float64
	}{{0, 0}, {1, 0.5}, {200, 100}, {255, 127.5}} {
		b523SensorViewBytes := make([]int, 18)
		b523SensorViewBytes[15] = b523SensorViewCase.b523SensorViewByte
		b523SensorViewOutput := b523SensorViewCall(b523SensorViewT, b523SensorViewInput(3, b523SensorViewBytes))
		if b523SensorViewField(b523SensorViewOutput, "PWM")["number"] != b523SensorViewCase.b523SensorViewWant {
			b523SensorViewT.Fatal("PWM carrier")
		}
	}
	for _, b523SensorViewKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		b523SensorViewInputValue := b523SensorViewInput(3, make([]int, 18))
		if b523SensorViewKey == "direction" || b523SensorViewKey == "view" {
			b523SensorViewInputValue["context"].(map[string]any)[b523SensorViewKey] = 0
		} else {
			b523SensorViewInputValue["context"].(map[string]any)[b523SensorViewKey] = "3"
		}
		b523SensorViewJSON, b523SensorViewErr := json.Marshal(b523SensorViewInputValue)
		if b523SensorViewErr != nil {
			b523SensorViewT.Fatal(b523SensorViewErr)
		}
		b523SensorViewResult, b523SensorViewErr := B523SensorViewInterpret(b523SensorViewJSON)
		if b523SensorViewErr == nil || b523SensorViewResult != nil {
			b523SensorViewT.Fatal("context kind admitted")
		}
	}
}

func TestAb27Ao27PartialIndependentFields(b523SensorViewT *testing.T) {
	for _, b523SensorViewOperation := range []int{3, 6, 7} {
		b523SensorViewExtent, b523SensorViewStatus, b523SensorViewCount, b523SensorViewStart := 18, 16, 6, 1
		if b523SensorViewOperation == 6 {
			b523SensorViewExtent, b523SensorViewStatus, b523SensorViewCount = 19, 17, 7
		}
		if b523SensorViewOperation == 7 {
			b523SensorViewStart = 8
		}
		for b523SensorViewWidth := 0; b523SensorViewWidth <= b523SensorViewExtent; b523SensorViewWidth++ {
			b523SensorViewBytes := make([]int, b523SensorViewWidth)
			for b523SensorViewIndex := range b523SensorViewBytes {
				b523SensorViewBytes[b523SensorViewIndex] = 255
			}
			b523SensorViewOutput := b523SensorViewCall(b523SensorViewT, b523SensorViewInput(b523SensorViewOperation, b523SensorViewBytes))
			for b523SensorViewIndex := 0; b523SensorViewIndex < b523SensorViewCount; b523SensorViewIndex++ {
				b523SensorViewFieldValue := b523SensorViewField(b523SensorViewOutput, fmt.Sprintf("statusS%d", b523SensorViewStart+b523SensorViewIndex))
				if b523SensorViewWidth > b523SensorViewStatus+b523SensorViewIndex/4 {
					if b523SensorViewFieldValue["label"] != "high-level" || b523SensorViewFieldValue["raw"] != float64(3) {
						b523SensorViewT.Fatal("complete status lost")
					}
				} else {
					if b523SensorViewFieldValue["state"] != "incomplete" {
						b523SensorViewT.Fatal("status defaulted")
					}
					if _, b523SensorViewPresent := b523SensorViewFieldValue["label"]; b523SensorViewPresent {
						b523SensorViewT.Fatal("missing status labelled")
					}
				}
			}
			b523SensorViewReserved := b523SensorViewField(b523SensorViewOutput, "status_reserved")
			if b523SensorViewWidth < b523SensorViewExtent {
				if b523SensorViewReserved["state"] != "incomplete" {
					b523SensorViewT.Fatal("reserved defaulted")
				}
			} else {
				b523SensorViewWant := float64(15)
				if b523SensorViewOperation == 6 {
					b523SensorViewWant = 3
				}
				if b523SensorViewReserved["reserved"] != b523SensorViewWant {
					b523SensorViewT.Fatal("reserved carrier")
				}
			}
			if b523SensorViewOperation == 3 || b523SensorViewOperation == 7 {
				b523SensorViewName, b523SensorViewPosition := "PWM", 15
				if b523SensorViewOperation == 7 {
					b523SensorViewName, b523SensorViewPosition = "opaque11", 11
				}
				b523SensorViewFieldValue := b523SensorViewField(b523SensorViewOutput, b523SensorViewName)
				if b523SensorViewWidth <= b523SensorViewPosition {
					if b523SensorViewFieldValue["state"] != "incomplete" {
						b523SensorViewT.Fatal("byte field defaulted")
					}
					if _, b523SensorViewPresent := b523SensorViewFieldValue["number"]; b523SensorViewPresent {
						b523SensorViewT.Fatal("incomplete byte converted")
					}
				} else if b523SensorViewFieldValue["raw"] != float64(255) {
					b523SensorViewT.Fatal("complete byte lost")
				}
			}
		}
	}
	b523SensorViewResult, b523SensorViewErr := B523SensorViewInterpret(json.RawMessage{'{', '"', 255, '"', ':', '0', '}'})
	if b523SensorViewErr == nil || b523SensorViewResult != nil {
		b523SensorViewT.Fatal("invalid UTF-8 admitted")
	}
}
