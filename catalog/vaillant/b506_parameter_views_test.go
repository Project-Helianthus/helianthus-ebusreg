package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func b506ParameterViewInput(b506ParameterViewKey any, b506ParameterViewValue any) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 6, "operation": 17, "direction": "request", "view": "parameter-union", "family": map[string]any{"opaque": []any{"synthetic"}}, "variant": []any{0, "opaque"}}, "fields": map[string]any{"key": b506ParameterViewKey, "value": b506ParameterViewValue, "tail": []int{0, 255}}, "arguments": map[string]any{}}
}

func b506ParameterViewRun(b506ParameterViewT *testing.T, b506ParameterViewInput any) map[string]any {
	b506ParameterViewT.Helper()
	b506ParameterViewJSON, b506ParameterViewErr := json.Marshal(b506ParameterViewInput)
	if b506ParameterViewErr != nil {
		b506ParameterViewT.Fatal(b506ParameterViewErr)
	}
	b506ParameterViewResult, b506ParameterViewErr := B506ParameterViewInterpret(b506ParameterViewJSON)
	if b506ParameterViewErr != nil {
		b506ParameterViewT.Fatal(b506ParameterViewErr)
	}
	b506ParameterViewOutput := map[string]any{}
	if b506ParameterViewErr := json.Unmarshal(b506ParameterViewResult, &b506ParameterViewOutput); b506ParameterViewErr != nil {
		b506ParameterViewT.Fatal(b506ParameterViewErr)
	}
	b506ParameterViewCanonical := map[string]any{}
	if b506ParameterViewErr := json.Unmarshal(b506ParameterViewJSON, &b506ParameterViewCanonical); b506ParameterViewErr != nil {
		b506ParameterViewT.Fatal(b506ParameterViewErr)
	}
	if !reflect.DeepEqual(b506ParameterViewOutput["context"], b506ParameterViewCanonical["context"]) || !reflect.DeepEqual(b506ParameterViewOutput["raw"], b506ParameterViewCanonical["fields"]) {
		b506ParameterViewT.Fatalf("evidence changed: %s", b506ParameterViewResult)
	}
	if b506ParameterViewArguments, b506ParameterViewPresent := b506ParameterViewCanonical["arguments"]; b506ParameterViewPresent {
		if !reflect.DeepEqual(b506ParameterViewArguments, b506ParameterViewOutput["arguments"]) {
			b506ParameterViewT.Fatal("arguments changed")
		}
	} else if _, b506ParameterViewPresent := b506ParameterViewOutput["arguments"]; b506ParameterViewPresent {
		b506ParameterViewT.Fatal("absent arguments fabricated")
	}
	if b506ParameterViewOutput["physical"] != "Unknown" || b506ParameterViewOutput["qualified"] != false {
		b506ParameterViewT.Fatal("physical qualification fabricated")
	}
	for b506ParameterViewName, b506ParameterViewSemantic := range b506ParameterViewOutput["fields"].(map[string]any) {
		b506ParameterViewRaw, b506ParameterViewHasRaw := b506ParameterViewSemantic.(map[string]any)["raw"]
		b506ParameterViewSupplied, b506ParameterViewPresent := b506ParameterViewCanonical["fields"].(map[string]any)[b506ParameterViewName]
		if b506ParameterViewHasRaw != b506ParameterViewPresent || (b506ParameterViewHasRaw && !reflect.DeepEqual(b506ParameterViewRaw, b506ParameterViewSupplied)) {
			b506ParameterViewT.Fatalf("semantic raw changed: %s", b506ParameterViewName)
		}
	}
	return b506ParameterViewOutput
}

func b506ParameterViewValueField(b506ParameterViewOutput map[string]any) map[string]any {
	return b506ParameterViewOutput["fields"].(map[string]any)["value"].(map[string]any)
}

func TestAb03Ao03Alternatives(b506ParameterViewT *testing.T) {
	b506ParameterViewData, b506ParameterViewErr := os.ReadFile("testdata/ab03ao03_b506_parameter_views.json")
	if b506ParameterViewErr != nil {
		b506ParameterViewT.Fatal(b506ParameterViewErr)
	}
	var b506ParameterViewCases []struct {
		B506ParameterViewID    string `json:"ab03ao03ID"`
		B506ParameterViewKey   int    `json:"ab03ao03Key"`
		B506ParameterViewValue []int  `json:"ab03ao03Value"`
		B506ParameterViewState string `json:"ab03ao03State"`
	}
	if b506ParameterViewErr := json.Unmarshal(b506ParameterViewData, &b506ParameterViewCases); b506ParameterViewErr != nil {
		b506ParameterViewT.Fatal(b506ParameterViewErr)
	}
	if len(b506ParameterViewCases) != 22 {
		b506ParameterViewT.Fatal("missing cases")
	}
	for b506ParameterViewPass := 0; b506ParameterViewPass < 2; b506ParameterViewPass++ {
		for b506ParameterViewIndex := range b506ParameterViewCases {
			b506ParameterViewCase := b506ParameterViewCases[(b506ParameterViewIndex*7+b506ParameterViewPass)%len(b506ParameterViewCases)]
			b506ParameterViewT.Run(fmt.Sprintf("%s-permutation-%d", b506ParameterViewCase.B506ParameterViewID, b506ParameterViewPass), func(b506ParameterViewT *testing.T) {
				b506ParameterViewOutput := b506ParameterViewRun(b506ParameterViewT, b506ParameterViewInput([]int{b506ParameterViewCase.B506ParameterViewKey}, b506ParameterViewCase.B506ParameterViewValue))
				if b506ParameterViewOutput["state"] != b506ParameterViewCase.B506ParameterViewState {
					b506ParameterViewT.Fatalf("state: %v", b506ParameterViewOutput)
				}
				if b506ParameterViewCase.B506ParameterViewState == "Unknown" {
					if len(b506ParameterViewOutput["fields"].(map[string]any)) != 0 {
						b506ParameterViewT.Fatal("unknown key selected")
					}
					return
				}
				b506ParameterViewKeyField := b506ParameterViewOutput["fields"].(map[string]any)["key"].(map[string]any)
				if b506ParameterViewKeyField["state"] != "Hypothesis" || b506ParameterViewKeyField["number"] != float64(b506ParameterViewCase.B506ParameterViewKey) {
					b506ParameterViewT.Fatal("key lost")
				}
				b506ParameterViewField := b506ParameterViewValueField(b506ParameterViewOutput)
				switch b506ParameterViewCase.B506ParameterViewKey {
				case 0, 1, 2, 4, 5:
					if len(b506ParameterViewField) != 2 || b506ParameterViewField["state"] != "Hypothesis" {
						b506ParameterViewT.Fatal("opaque bytes interpreted")
					}
					if _, b506ParameterViewOK := b506ParameterViewField["raw"].([]any); !b506ParameterViewOK {
						b506ParameterViewT.Fatal("opaque bytes converted to text")
					}
				case 10:
					if b506ParameterViewCase.B506ParameterViewValue[0] == 2 {
						if b506ParameterViewField["state"] != "unmapped" || len(b506ParameterViewField) != 2 {
							b506ParameterViewT.Fatal("Boolean 2 mapped")
						}
					} else {
						b506ParameterViewLabel := "no"
						if b506ParameterViewCase.B506ParameterViewValue[0] == 1 {
							b506ParameterViewLabel = "yes"
						}
						if b506ParameterViewField["label"] != b506ParameterViewLabel || b506ParameterViewField["number"] != float64(b506ParameterViewCase.B506ParameterViewValue[0]) || b506ParameterViewField["state"] != "Hypothesis" {
							b506ParameterViewT.Fatal("Boolean mismatch")
						}
					}
				default:
					if b506ParameterViewField["number"] != float64(127) || b506ParameterViewField["state"] != "Hypothesis" {
						b506ParameterViewT.Fatal("scalar mismatch")
					}
				}
			})
		}
	}
}

func TestAb03Ao03WidthsAndPresence(b506ParameterViewT *testing.T) {
	for _, b506ParameterViewKey := range []int{0, 1, 2, 4, 5, 3, 6, 7, 8, 9, 10, 32, 33, 34, 35, 36, 37} {
		for _, b506ParameterViewWidth := range []int{-1, 0, 1, 2, 3} {
			b506ParameterViewT.Run(fmt.Sprintf("key-%d-width-%d", b506ParameterViewKey, b506ParameterViewWidth), func(b506ParameterViewT *testing.T) {
				b506ParameterViewBytes := []int{}
				for b506ParameterViewIndex := 0; b506ParameterViewIndex < b506ParameterViewWidth; b506ParameterViewIndex++ {
					b506ParameterViewBytes = append(b506ParameterViewBytes, 0)
				}
				b506ParameterViewIn := b506ParameterViewInput([]int{b506ParameterViewKey}, b506ParameterViewBytes)
				if b506ParameterViewWidth < 0 {
					delete(b506ParameterViewIn["fields"].(map[string]any), "value")
				}
				b506ParameterViewOutput := b506ParameterViewRun(b506ParameterViewT, b506ParameterViewIn)
				b506ParameterViewField := b506ParameterViewValueField(b506ParameterViewOutput)
				b506ParameterViewOpaque := b506ParameterViewKey == 0 || b506ParameterViewKey == 1 || b506ParameterViewKey == 2 || b506ParameterViewKey == 4 || b506ParameterViewKey == 5
				b506ParameterViewExpected := "Hypothesis"
				if b506ParameterViewWidth < 0 || (!b506ParameterViewOpaque && b506ParameterViewWidth == 0) {
					b506ParameterViewExpected = "incomplete"
				}
				if !b506ParameterViewOpaque && b506ParameterViewWidth > 1 {
					b506ParameterViewExpected = "width-conflict"
				}
				if b506ParameterViewField["state"] != b506ParameterViewExpected {
					b506ParameterViewT.Fatalf("width state: %v", b506ParameterViewField)
				}
				_, b506ParameterViewHasRaw := b506ParameterViewField["raw"]
				if b506ParameterViewHasRaw != (b506ParameterViewWidth >= 0) {
					b506ParameterViewT.Fatal("missing/empty conflated")
				}
				_, b506ParameterViewHasNumber := b506ParameterViewField["number"]
				if b506ParameterViewHasNumber != (!b506ParameterViewOpaque && b506ParameterViewWidth == 1) {
					b506ParameterViewT.Fatal("invalid numeric presence")
				}
				if b506ParameterViewExpected != "Hypothesis" {
					if b506ParameterViewOutput["state"] != "partial" || b506ParameterViewOutput["fields"].(map[string]any)["key"].(map[string]any)["number"] != float64(b506ParameterViewKey) {
						b506ParameterViewT.Fatal("complete independent key erased")
					}
				} else if b506ParameterViewOutput["state"] != "Hypothesis" {
					b506ParameterViewT.Fatal("complete state")
				}
			})
		}
	}
	for _, b506ParameterViewKey := range []any{nil, []int{}, []int{2, 3}} {
		b506ParameterViewIn := b506ParameterViewInput(b506ParameterViewKey, []int{0})
		if b506ParameterViewKey == nil {
			delete(b506ParameterViewIn["fields"].(map[string]any), "key")
		}
		b506ParameterViewOutput := b506ParameterViewRun(b506ParameterViewT, b506ParameterViewIn)
		if b506ParameterViewOutput["state"] != "partial" {
			b506ParameterViewT.Fatal("invalid key width selected")
		}
	}
}

func TestAb03Ao03ContextAxes(b506ParameterViewT *testing.T) {
	for b506ParameterViewAxis, b506ParameterViewValue := range map[string]any{"primary": 180, "secondary": 7, "operation": 18, "direction": "response", "view": "other"} {
		for _, b506ParameterViewMissing := range []bool{false, true} {
			b506ParameterViewT.Run(fmt.Sprintf("%s-missing-%v", b506ParameterViewAxis, b506ParameterViewMissing), func(b506ParameterViewT *testing.T) {
				b506ParameterViewIn := b506ParameterViewInput([]int{3}, []int{127})
				b506ParameterViewContext := b506ParameterViewIn["context"].(map[string]any)
				if b506ParameterViewMissing {
					delete(b506ParameterViewContext, b506ParameterViewAxis)
				} else {
					b506ParameterViewContext[b506ParameterViewAxis] = b506ParameterViewValue
				}
				b506ParameterViewOutput := b506ParameterViewRun(b506ParameterViewT, b506ParameterViewIn)
				if b506ParameterViewOutput["state"] != "Unknown" || len(b506ParameterViewOutput["fields"].(map[string]any)) != 0 {
					b506ParameterViewT.Fatal("selector guessed")
				}
			})
		}
	}
	for _, b506ParameterViewAxis := range []string{"family", "variant"} {
		for _, b506ParameterViewValue := range []any{nil, 0, "parameter-union", []any{}, map[string]any{"different": true}} {
			b506ParameterViewIn := b506ParameterViewInput([]int{3}, []int{127})
			b506ParameterViewIn["context"].(map[string]any)[b506ParameterViewAxis] = b506ParameterViewValue
			if b506ParameterViewRun(b506ParameterViewT, b506ParameterViewIn)["state"] != "Hypothesis" {
				b506ParameterViewT.Fatal("opaque context selects")
			}
		}
	}
	for _, b506ParameterViewTail := range [][]int{{}, {1}, {255, 0, 127}} {
		b506ParameterViewIn := b506ParameterViewInput([]int{3}, []int{0})
		b506ParameterViewIn["fields"].(map[string]any)["tail"] = b506ParameterViewTail
		delete(b506ParameterViewIn, "arguments")
		if b506ParameterViewValueField(b506ParameterViewRun(b506ParameterViewT, b506ParameterViewIn))["number"] != float64(0) {
			b506ParameterViewT.Fatal("tail changed scalar")
		}
	}
}

func TestAb03Ao03Reject(b506ParameterViewT *testing.T) {
	for _, b506ParameterViewJSON := range []string{``, `{`, `{} {}`, `[]`, `null`, `{"context":{},"fields":{},"fields":{}}`, `{"context":{"x":1,"x":2},"fields":{}}`, `{"context":{"x":[{"a":0,"a":1}]},"fields":{}}`, `{"context":{},"fields":null}`} {
		if b506ParameterViewResult, b506ParameterViewErr := B506ParameterViewInterpret(json.RawMessage(b506ParameterViewJSON)); b506ParameterViewErr == nil || b506ParameterViewResult != nil {
			b506ParameterViewT.Fatalf("accepted %q", b506ParameterViewJSON)
		}
	}
	for _, b506ParameterViewAxis := range []string{"key", "value", "tail"} {
		for _, b506ParameterViewBad := range []any{nil, "bytes", 0, map[string]any{}, []any{-1}, []any{256}, []any{1.5}, []any{true}, []any{"0"}, []any{nil}, []any{[]int{0}}} {
			b506ParameterViewIn := b506ParameterViewInput([]int{3}, []int{0})
			b506ParameterViewIn["fields"].(map[string]any)[b506ParameterViewAxis] = b506ParameterViewBad
			b506ParameterViewJSON, _ := json.Marshal(b506ParameterViewIn)
			if b506ParameterViewResult, b506ParameterViewErr := B506ParameterViewInterpret(b506ParameterViewJSON); b506ParameterViewErr == nil || b506ParameterViewResult != nil {
				b506ParameterViewT.Fatalf("accepted %s", b506ParameterViewJSON)
			}
		}
	}
	for _, b506ParameterViewAxis := range []string{"primary", "secondary", "operation", "direction", "view"} {
		b506ParameterViewIn := b506ParameterViewInput([]int{3}, []int{0})
		b506ParameterViewIn["context"].(map[string]any)[b506ParameterViewAxis] = []any{}
		b506ParameterViewJSON, _ := json.Marshal(b506ParameterViewIn)
		if _, b506ParameterViewErr := B506ParameterViewInterpret(b506ParameterViewJSON); b506ParameterViewErr == nil {
			b506ParameterViewT.Fatal("wrong selector kind")
		}
	}
	for _, b506ParameterViewAxis := range []string{"context", "fields", "arguments"} {
		b506ParameterViewIn := b506ParameterViewInput([]int{3}, []int{0})
		b506ParameterViewIn[b506ParameterViewAxis] = []any{}
		b506ParameterViewJSON, _ := json.Marshal(b506ParameterViewIn)
		if _, b506ParameterViewErr := B506ParameterViewInterpret(b506ParameterViewJSON); b506ParameterViewErr == nil {
			b506ParameterViewT.Fatal("wrong object kind")
		}
	}
	for _, b506ParameterViewArguments := range []any{nil, map[string]any{"flag": true}, map[string]any{"count": 0}} {
		b506ParameterViewIn := b506ParameterViewInput([]int{3}, []int{0})
		b506ParameterViewIn["arguments"] = b506ParameterViewArguments
		b506ParameterViewJSON, _ := json.Marshal(b506ParameterViewIn)
		if _, b506ParameterViewErr := B506ParameterViewInterpret(b506ParameterViewJSON); b506ParameterViewErr == nil {
			b506ParameterViewT.Fatal("undefined arguments")
		}
	}
}

func TestAb03Ao03Memory(b506ParameterViewT *testing.T) {
	b506ParameterViewInputBytes, _ := json.Marshal(b506ParameterViewInput([]int{2}, []int{0, 170, 255}))
	b506ParameterViewOriginal := append([]byte(nil), b506ParameterViewInputBytes...)
	b506ParameterViewFirst, b506ParameterViewErr := B506ParameterViewInterpret(b506ParameterViewInputBytes)
	if b506ParameterViewErr != nil {
		b506ParameterViewT.Fatal(b506ParameterViewErr)
	}
	b506ParameterViewSaved := append([]byte(nil), b506ParameterViewFirst...)
	for b506ParameterViewIndex := range b506ParameterViewInputBytes {
		b506ParameterViewInputBytes[b506ParameterViewIndex] = ' '
	}
	if !bytes.Equal(b506ParameterViewFirst, b506ParameterViewSaved) {
		b506ParameterViewT.Fatal("input aliases output")
	}
	b506ParameterViewSecond, b506ParameterViewErr := B506ParameterViewInterpret(b506ParameterViewOriginal)
	if b506ParameterViewErr != nil {
		b506ParameterViewT.Fatal(b506ParameterViewErr)
	}
	b506ParameterViewNested := map[string]any{}
	if b506ParameterViewErr := json.Unmarshal(b506ParameterViewFirst, &b506ParameterViewNested); b506ParameterViewErr != nil {
		b506ParameterViewT.Fatal(b506ParameterViewErr)
	}
	b506ParameterViewNested["raw"].(map[string]any)["value"].([]any)[0] = float64(99)
	b506ParameterViewNested["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "changed"
	b506ParameterViewValueField(b506ParameterViewNested)["raw"].([]any)[1] = float64(99)
	for b506ParameterViewIndex := range b506ParameterViewFirst {
		b506ParameterViewFirst[b506ParameterViewIndex] = ' '
	}
	b506ParameterViewThird, b506ParameterViewErr := B506ParameterViewInterpret(b506ParameterViewOriginal)
	if b506ParameterViewErr != nil || !bytes.Equal(b506ParameterViewSecond, b506ParameterViewSaved) || !bytes.Equal(b506ParameterViewThird, b506ParameterViewSaved) {
		b506ParameterViewT.Fatal("result aliases another result or later call")
	}
}

func TestAb03Ao03KeyVersusValue(b506ParameterViewT *testing.T) {
	for _, b506ParameterViewKey := range []int{2, 3} {
		b506ParameterViewOutput := b506ParameterViewRun(b506ParameterViewT, b506ParameterViewInput([]int{b506ParameterViewKey}, []int{0, 170, 255}))
		b506ParameterViewExpected := "Hypothesis"
		if b506ParameterViewKey == 3 {
			b506ParameterViewExpected = "width-conflict"
		}
		if b506ParameterViewValueField(b506ParameterViewOutput)["state"] != b506ParameterViewExpected {
			b506ParameterViewT.Fatal("key did not select independent alternative")
		}
	}
	for _, b506ParameterViewValue := range []int{0, 1, 127, 255} {
		if b506ParameterViewValueField(b506ParameterViewRun(b506ParameterViewT, b506ParameterViewInput([]int{3}, []int{b506ParameterViewValue})))["number"] != float64(b506ParameterViewValue) {
			b506ParameterViewT.Fatal("scalar ignored byte")
		}
	}
	for _, b506ParameterViewAxis := range []string{"primary", "secondary", "operation"} {
		for _, b506ParameterViewValue := range []any{0, -1, json.Number("999999999999999999999999999999999999")} {
			b506ParameterViewIn := b506ParameterViewInput([]int{3}, []int{127})
			b506ParameterViewIn["context"].(map[string]any)[b506ParameterViewAxis] = b506ParameterViewValue
			if b506ParameterViewRun(b506ParameterViewT, b506ParameterViewIn)["state"] != "Unknown" {
				b506ParameterViewT.Fatal("unknown integer selector admitted")
			}
		}
	}
}
