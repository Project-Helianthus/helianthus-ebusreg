package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func b519OpaqueViewCall(b519OpaqueViewT *testing.T, b519OpaqueViewInput any) map[string]any {
	b519OpaqueViewT.Helper()
	b519OpaqueViewJSON, b519OpaqueViewErr := json.Marshal(b519OpaqueViewInput)
	if b519OpaqueViewErr != nil {
		b519OpaqueViewT.Fatal(b519OpaqueViewErr)
	}
	b519OpaqueViewResult, b519OpaqueViewErr := B519OpaqueViewInterpret(b519OpaqueViewJSON)
	if b519OpaqueViewErr != nil {
		b519OpaqueViewT.Fatal(b519OpaqueViewErr)
	}
	b519OpaqueViewOutput := map[string]any{}
	if b519OpaqueViewErr = json.Unmarshal(b519OpaqueViewResult, &b519OpaqueViewOutput); b519OpaqueViewErr != nil {
		b519OpaqueViewT.Fatal(b519OpaqueViewErr)
	}
	if b519OpaqueViewOutput["physical"] != "Unknown" || b519OpaqueViewOutput["qualified"] != false {
		b519OpaqueViewT.Fatal(b519OpaqueViewOutput)
	}
	return b519OpaqueViewOutput
}

func b519OpaqueViewInput(b519OpaqueViewOperation int, b519OpaqueViewDirection string, b519OpaqueViewBytes []int) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 25, "operation": b519OpaqueViewOperation, "direction": b519OpaqueViewDirection, "view": "logical-body", "family": map[string]any{"opaque": []any{"synthetic", 0}}, "variant": "synthetic"}, "fields": map[string]any{"logical_object": b519OpaqueViewBytes, "extra": []int{0, 255}}}
}

func TestAb21Ao21SelectorsWidths(b519OpaqueViewT *testing.T) {
	for _, b519OpaqueViewOperation := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 170} {
		for _, b519OpaqueViewDirection := range []string{"request", "positive"} {
			for b519OpaqueViewWidth := 0; b519OpaqueViewWidth <= 5; b519OpaqueViewWidth++ {
				b519OpaqueViewT.Run(fmt.Sprintf("op%d_%s_width%d", b519OpaqueViewOperation, b519OpaqueViewDirection, b519OpaqueViewWidth), func(b519OpaqueViewT *testing.T) {
					b519OpaqueViewBytes := []int{19, 38, 57, 0, 255}[:b519OpaqueViewWidth]
					b519OpaqueViewInputValue := b519OpaqueViewInput(b519OpaqueViewOperation, b519OpaqueViewDirection, b519OpaqueViewBytes)
					b519OpaqueViewOutput := b519OpaqueViewCall(b519OpaqueViewT, b519OpaqueViewInputValue)
					b519OpaqueViewRaw, _ := json.Marshal(b519OpaqueViewOutput["raw"])
					b519OpaqueViewExpected, _ := json.Marshal(b519OpaqueViewInputValue["fields"])
					if !bytes.Equal(b519OpaqueViewRaw, b519OpaqueViewExpected) {
						b519OpaqueViewT.Fatal("raw changed")
					}
					b519OpaqueViewFields := b519OpaqueViewOutput["fields"].(map[string]any)
					b519OpaqueViewUnknown := b519OpaqueViewDirection == "positive" && b519OpaqueViewOperation != 3 && b519OpaqueViewOperation != 6 && b519OpaqueViewOperation != 170
					if b519OpaqueViewUnknown {
						if b519OpaqueViewOutput["state"] != "Unknown" || len(b519OpaqueViewFields) != 0 {
							b519OpaqueViewT.Fatal(b519OpaqueViewOutput)
						}
						return
					}
					if b519OpaqueViewWidth < 3 {
						if b519OpaqueViewOutput["state"] != "partial" {
							b519OpaqueViewT.Fatal(b519OpaqueViewOutput)
						}
						return
					}
					b519OpaqueViewForm := b519OpaqueViewDirection == "positive" && b519OpaqueViewOperation == 6 || b519OpaqueViewDirection == "request" && (b519OpaqueViewOperation == 0 || b519OpaqueViewOperation == 2 || b519OpaqueViewOperation == 170)
					if b519OpaqueViewForm {
						b519OpaqueViewLabel := "no-declared-argument"
						if b519OpaqueViewDirection == "positive" {
							b519OpaqueViewLabel = "prefix-only"
						}
						if b519OpaqueViewFields["form"].(map[string]any)["label"] != b519OpaqueViewLabel {
							b519OpaqueViewT.Fatal(b519OpaqueViewOutput)
						}
						if b519OpaqueViewWidth > 3 {
							if b519OpaqueViewFields["unmodelled"] == nil || b519OpaqueViewOutput["state"] != "partial" {
								b519OpaqueViewT.Fatal(b519OpaqueViewOutput)
							}
						}
						return
					}
					if b519OpaqueViewOperation == 8 && b519OpaqueViewDirection == "request" && b519OpaqueViewWidth == 3 {
						if b519OpaqueViewOutput["state"] != "partial" {
							b519OpaqueViewT.Fatal(b519OpaqueViewOutput)
						}
						return
					}
					b519OpaqueViewOffset := 3
					if b519OpaqueViewOperation == 8 && b519OpaqueViewDirection == "request" {
						b519OpaqueViewOffset = 4
					}
					b519OpaqueViewWant, _ := json.Marshal(b519OpaqueViewBytes[b519OpaqueViewOffset:])
					b519OpaqueViewGot, _ := json.Marshal(b519OpaqueViewFields["opaque"].(map[string]any)["raw"])
					if !bytes.Equal(b519OpaqueViewWant, b519OpaqueViewGot) {
						b519OpaqueViewT.Fatal(string(b519OpaqueViewGot), string(b519OpaqueViewWant))
					}
				})
			}
		}
	}
}

func TestAb21Ao21EveryByte(b519OpaqueViewT *testing.T) {
	for _, b519OpaqueViewDirection := range []string{"request", "positive"} {
		for _, b519OpaqueViewOperation := range []int{3, 170, 8} {
			for b519OpaqueViewPosition := 0; b519OpaqueViewPosition < 7; b519OpaqueViewPosition++ {
				for b519OpaqueViewByte := 0; b519OpaqueViewByte < 256; b519OpaqueViewByte++ {
					b519OpaqueViewBytes := []int{0, 0, 0, 127, 0, 255, 65}
					b519OpaqueViewBytes[b519OpaqueViewPosition] = b519OpaqueViewByte
					b519OpaqueViewOutput := b519OpaqueViewCall(b519OpaqueViewT, b519OpaqueViewInput(b519OpaqueViewOperation, b519OpaqueViewDirection, b519OpaqueViewBytes))
					b519OpaqueViewFields := b519OpaqueViewOutput["fields"].(map[string]any)
					if b519OpaqueViewDirection == "positive" && b519OpaqueViewOperation == 8 {
						continue
					}
					if b519OpaqueViewDirection == "request" && b519OpaqueViewOperation == 170 {
						continue
					}
					b519OpaqueViewOffset := 3
					if b519OpaqueViewOperation == 8 {
						b519OpaqueViewOffset = 4
						b519OpaqueViewAddress := b519OpaqueViewFields["slave_address"].(map[string]any)
						if !reflect.DeepEqual(b519OpaqueViewAddress["raw"], []any{float64(b519OpaqueViewBytes[3])}) {
							b519OpaqueViewT.Fatal(b519OpaqueViewAddress)
						}
					}
					b519OpaqueViewOpaque := b519OpaqueViewFields["opaque"].(map[string]any)
					b519OpaqueViewWant := []any{}
					for _, b519OpaqueViewValue := range b519OpaqueViewBytes[b519OpaqueViewOffset:] {
						b519OpaqueViewWant = append(b519OpaqueViewWant, float64(b519OpaqueViewValue))
					}
					if !reflect.DeepEqual(b519OpaqueViewOpaque["raw"], b519OpaqueViewWant) || b519OpaqueViewOpaque["number"] != nil || b519OpaqueViewOpaque["state"] != "Hypothesis" {
						b519OpaqueViewT.Fatal(b519OpaqueViewOpaque)
					}
					if !reflect.DeepEqual(b519OpaqueViewOpaque["flags"], map[string]any{"accepted_empty": "Unknown", "unlimited_carrier": "Unknown"}) {
						b519OpaqueViewT.Fatal(b519OpaqueViewOpaque)
					}
				}
			}
		}
	}
}

func TestAb21Ao21ContextPresence(b519OpaqueViewT *testing.T) {
	for b519OpaqueViewKey, b519OpaqueViewValue := range map[string]any{"primary": 180, "secondary": 24, "operation": 9, "direction": "response", "view": "guess"} {
		b519OpaqueViewInputValue := b519OpaqueViewInput(3, "request", []int{0, 0, 0, 0, 255, 65})
		b519OpaqueViewInputValue["context"].(map[string]any)[b519OpaqueViewKey] = b519OpaqueViewValue
		if b519OpaqueViewOutput := b519OpaqueViewCall(b519OpaqueViewT, b519OpaqueViewInputValue); b519OpaqueViewOutput["state"] != "Unknown" || len(b519OpaqueViewOutput["fields"].(map[string]any)) != 0 {
			b519OpaqueViewT.Fatal(b519OpaqueViewOutput)
		}
		delete(b519OpaqueViewInputValue["context"].(map[string]any), b519OpaqueViewKey)
		if b519OpaqueViewCall(b519OpaqueViewT, b519OpaqueViewInputValue)["state"] != "Unknown" {
			b519OpaqueViewT.Fatal(b519OpaqueViewKey)
		}
	}
	for _, b519OpaqueViewKey := range []string{"family", "variant", "unmodelled_context"} {
		b519OpaqueViewInputValue := b519OpaqueViewInput(3, "request", []int{0, 0, 0, 0})
		b519OpaqueViewInputValue["context"].(map[string]any)[b519OpaqueViewKey] = map[string]any{"nested": []any{nil, false, "opaque", 0}}
		b519OpaqueViewOutput := b519OpaqueViewCall(b519OpaqueViewT, b519OpaqueViewInputValue)
		b519OpaqueViewWant, _ := json.Marshal(b519OpaqueViewInputValue["context"])
		b519OpaqueViewGot, _ := json.Marshal(b519OpaqueViewOutput["context"])
		if !bytes.Equal(b519OpaqueViewWant, b519OpaqueViewGot) || b519OpaqueViewOutput["state"] != "Hypothesis" {
			b519OpaqueViewT.Fatal(b519OpaqueViewOutput)
		}
	}
	for _, b519OpaqueViewBytes := range [][]int{nil, {}, {0}} {
		b519OpaqueViewInputValue := b519OpaqueViewInput(3, "request", []int{})
		if b519OpaqueViewBytes == nil {
			delete(b519OpaqueViewInputValue["fields"].(map[string]any), "logical_object")
		} else {
			b519OpaqueViewInputValue["fields"].(map[string]any)["logical_object"] = b519OpaqueViewBytes
		}
		b519OpaqueViewOutput := b519OpaqueViewCall(b519OpaqueViewT, b519OpaqueViewInputValue)
		if b519OpaqueViewOutput["state"] != "partial" {
			b519OpaqueViewT.Fatal(b519OpaqueViewOutput)
		}
	}
	b519OpaqueViewInputValue := b519OpaqueViewInput(0, "request", []int{0, 0, 0})
	b519OpaqueViewWithout := b519OpaqueViewCall(b519OpaqueViewT, b519OpaqueViewInputValue)
	b519OpaqueViewInputValue["arguments"] = map[string]any{}
	b519OpaqueViewWith := b519OpaqueViewCall(b519OpaqueViewT, b519OpaqueViewInputValue)
	if _, b519OpaqueViewPresent := b519OpaqueViewWithout["arguments"]; b519OpaqueViewPresent || b519OpaqueViewWith["arguments"] == nil {
		b519OpaqueViewT.Fatal("argument presence lost")
	}
}

func TestAb21Ao21Invalid(b519OpaqueViewT *testing.T) {
	for _, b519OpaqueViewJSON := range []string{``, `{`, `[]`, `null`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"x":{"a":1,"a":2}},"fields":{}}`, `{"context":{},"fields":{"a":[],"a":[]}}`, `{"context":{},"fields":null}`, `{"context":[],"fields":{}}`, `{"context":{"operation":1.0},"fields":{}}`, `{"context":{"primary":true},"fields":{}}`, `{"context":{"secondary":null},"fields":{}}`, `{"context":{"direction":0},"fields":{}}`, `{"context":{"view":[]},"fields":{}}`, `{"context":{},"fields":{},"arguments":null}`, `{"context":{},"fields":{},"arguments":{"flag":false}}`} {
		if b519OpaqueViewResult, b519OpaqueViewErr := B519OpaqueViewInterpret([]byte(b519OpaqueViewJSON)); b519OpaqueViewErr == nil || b519OpaqueViewResult != nil {
			b519OpaqueViewT.Fatal(b519OpaqueViewJSON, string(b519OpaqueViewResult), b519OpaqueViewErr)
		}
	}
	for _, b519OpaqueViewBad := range []string{`null`, `"AA"`, `[null]`, `[true]`, `["0"]`, `[-1]`, `[256]`, `[1.5]`, `[1e0]`, `[[]]`, `[{}]`} {
		b519OpaqueViewJSON := `{"context":{},"fields":{"other":` + b519OpaqueViewBad + `}}`
		if b519OpaqueViewResult, b519OpaqueViewErr := B519OpaqueViewInterpret([]byte(b519OpaqueViewJSON)); b519OpaqueViewErr == nil || b519OpaqueViewResult != nil {
			b519OpaqueViewT.Fatal(b519OpaqueViewJSON)
		}
	}
}

func TestAb21Ao21IsolationFixtures(b519OpaqueViewT *testing.T) {
	b519OpaqueViewData, b519OpaqueViewErr := os.ReadFile("testdata/ab21ao21_b519_opaque_views.json")
	if b519OpaqueViewErr != nil {
		b519OpaqueViewT.Fatal(b519OpaqueViewErr)
	}
	var b519OpaqueViewCases []struct {
		B519OpaqueViewID     string          `json:"id"`
		B519OpaqueViewInput  json.RawMessage `json:"input"`
		B519OpaqueViewFields map[string]any  `json:"fields"`
	}
	if b519OpaqueViewErr = json.Unmarshal(b519OpaqueViewData, &b519OpaqueViewCases); b519OpaqueViewErr != nil || len(b519OpaqueViewCases) == 0 {
		b519OpaqueViewT.Fatal(b519OpaqueViewErr)
	}
	for _, b519OpaqueViewCase := range b519OpaqueViewCases {
		b519OpaqueViewT.Run(b519OpaqueViewCase.B519OpaqueViewID, func(b519OpaqueViewT *testing.T) {
			b519OpaqueViewInputCopy := append(json.RawMessage{}, b519OpaqueViewCase.B519OpaqueViewInput...)
			b519OpaqueViewResult, b519OpaqueViewErr := B519OpaqueViewInterpret(b519OpaqueViewInputCopy)
			if b519OpaqueViewErr != nil {
				b519OpaqueViewT.Fatal(b519OpaqueViewErr)
			}
			b519OpaqueViewSnapshot := append([]byte{}, b519OpaqueViewResult...)
			for b519OpaqueViewIndex := range b519OpaqueViewInputCopy {
				b519OpaqueViewInputCopy[b519OpaqueViewIndex] = ' '
			}
			if !bytes.Equal(b519OpaqueViewSnapshot, b519OpaqueViewResult) {
				b519OpaqueViewT.Fatal("input aliases result")
			}
			b519OpaqueViewOutput := map[string]any{}
			if b519OpaqueViewErr = json.Unmarshal(b519OpaqueViewResult, &b519OpaqueViewOutput); b519OpaqueViewErr != nil {
				b519OpaqueViewT.Fatal(b519OpaqueViewErr)
			}
			if !reflect.DeepEqual(b519OpaqueViewOutput["fields"], b519OpaqueViewCase.B519OpaqueViewFields) {
				b519OpaqueViewT.Fatal(b519OpaqueViewOutput)
			}
			b519OpaqueViewOutput["raw"].(map[string]any)["logical_object"].([]any)[0] = 999
			for b519OpaqueViewIndex := range b519OpaqueViewResult {
				b519OpaqueViewResult[b519OpaqueViewIndex] = ' '
			}
			b519OpaqueViewAgain, b519OpaqueViewErr := B519OpaqueViewInterpret(b519OpaqueViewCase.B519OpaqueViewInput)
			if b519OpaqueViewErr != nil || !bytes.Equal(b519OpaqueViewSnapshot, b519OpaqueViewAgain) {
				b519OpaqueViewT.Fatal("result aliases later call")
			}
		})
	}
}

func TestAb21Ao21PresenceAndSiblings(b519OpaqueViewT *testing.T) {
	for _, b519OpaqueViewKey := range []string{"logical_object", "extra"} {
		for _, b519OpaqueViewBytes := range [][]int{{}, {0}, {255}, {0, 255, 0, 255, 65, 1}} {
			b519OpaqueViewValue := b519OpaqueViewInput(8, "request", []int{91, 92, 93, 127, 0, 255})
			b519OpaqueViewValue["fields"].(map[string]any)[b519OpaqueViewKey] = b519OpaqueViewBytes
			b519OpaqueViewOutput := b519OpaqueViewCall(b519OpaqueViewT, b519OpaqueViewValue)
			b519OpaqueViewWant, _ := json.Marshal(b519OpaqueViewBytes)
			b519OpaqueViewGot, _ := json.Marshal(b519OpaqueViewOutput["raw"].(map[string]any)[b519OpaqueViewKey])
			if !bytes.Equal(b519OpaqueViewWant, b519OpaqueViewGot) {
				b519OpaqueViewT.Fatal("raw width/value lost")
			}
			if b519OpaqueViewKey == "extra" {
				if b519OpaqueViewOutput["fields"].(map[string]any)["opaque"] == nil || b519OpaqueViewOutput["fields"].(map[string]any)["slave_address"] == nil {
					b519OpaqueViewT.Fatal("sibling erased")
				}
			}
		}
	}
	b519OpaqueViewMissing := b519OpaqueViewInput(3, "request", []int{})
	delete(b519OpaqueViewMissing["fields"].(map[string]any), "logical_object")
	b519OpaqueViewAbsent := b519OpaqueViewCall(b519OpaqueViewT, b519OpaqueViewMissing)
	b519OpaqueViewEmpty := b519OpaqueViewCall(b519OpaqueViewT, b519OpaqueViewInput(3, "request", []int{}))
	b519OpaqueViewZero := b519OpaqueViewCall(b519OpaqueViewT, b519OpaqueViewInput(3, "request", []int{0}))
	if reflect.DeepEqual(b519OpaqueViewAbsent, b519OpaqueViewEmpty) || reflect.DeepEqual(b519OpaqueViewEmpty, b519OpaqueViewZero) {
		b519OpaqueViewT.Fatal("presence collapsed")
	}
	b519OpaqueViewComplete := b519OpaqueViewCall(b519OpaqueViewT, b519OpaqueViewInput(8, "request", []int{0, 0, 0, 0}))
	if !reflect.DeepEqual(b519OpaqueViewComplete["fields"].(map[string]any)["slave_address"].(map[string]any)["raw"], []any{float64(0)}) || !reflect.DeepEqual(b519OpaqueViewComplete["fields"].(map[string]any)["opaque"].(map[string]any)["raw"], []any{}) {
		b519OpaqueViewT.Fatal(b519OpaqueViewComplete)
	}
	for _, b519OpaqueViewOperation := range []string{"99999999999999999999999999999999999999", "-99999999999999999999999999999999999999"} {
		b519OpaqueViewJSON := []byte(`{"context":{"primary":181,"secondary":25,"operation":` + b519OpaqueViewOperation + `,"direction":"request","view":"logical-body"},"fields":{"logical_object":[0,0,0]}}`)
		b519OpaqueViewResult, b519OpaqueViewErr := B519OpaqueViewInterpret(b519OpaqueViewJSON)
		if b519OpaqueViewErr != nil {
			b519OpaqueViewT.Fatal(b519OpaqueViewErr)
		}
		b519OpaqueViewOutput := map[string]any{}
		if b519OpaqueViewErr = json.Unmarshal(b519OpaqueViewResult, &b519OpaqueViewOutput); b519OpaqueViewErr != nil {
			b519OpaqueViewT.Fatal(b519OpaqueViewErr)
		}
		if b519OpaqueViewOutput["state"] != "Unknown" {
			b519OpaqueViewT.Fatal(b519OpaqueViewOutput)
		}
	}
}
