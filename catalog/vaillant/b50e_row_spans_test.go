package vaillant

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func b50ERowSpanInput(b50ERowSpanKey int, b50ERowSpanValue []int) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 14, "operation": 2, "direction": "positive", "view": "row-spans", "family": map[string]any{"opaque": []any{"synthetic", 0}}, "variant": "candidate"}, "fields": map[string]any{"key": []int{b50ERowSpanKey}, "value": b50ERowSpanValue, "unmodelled": []int{99, 0, 255}}}
}

func b50ERowSpanCall(b50ERowSpanT *testing.T, b50ERowSpanObject map[string]any) map[string]any {
	b50ERowSpanT.Helper()
	b50ERowSpanInputBytes, b50ERowSpanErr := json.Marshal(b50ERowSpanObject)
	if b50ERowSpanErr != nil {
		b50ERowSpanT.Fatal(b50ERowSpanErr)
	}
	b50ERowSpanOutputBytes, b50ERowSpanErr := B50ERowSpanInterpret(b50ERowSpanInputBytes)
	if b50ERowSpanErr != nil {
		b50ERowSpanT.Fatal(b50ERowSpanErr)
	}
	b50ERowSpanOutput := map[string]any{}
	if b50ERowSpanErr = json.Unmarshal(b50ERowSpanOutputBytes, &b50ERowSpanOutput); b50ERowSpanErr != nil {
		b50ERowSpanT.Fatal(b50ERowSpanErr)
	}
	b50ERowSpanCanonical := map[string]any{}
	_ = json.Unmarshal(b50ERowSpanInputBytes, &b50ERowSpanCanonical)
	if !reflect.DeepEqual(b50ERowSpanOutput["context"], b50ERowSpanCanonical["context"]) || !reflect.DeepEqual(b50ERowSpanOutput["raw"], b50ERowSpanCanonical["fields"]) {
		b50ERowSpanT.Fatal("context/raw changed", b50ERowSpanOutput)
	}
	if b50ERowSpanOutput["qualified"] != false || b50ERowSpanOutput["physical"] != "Unknown" {
		b50ERowSpanT.Fatal("qualification inferred")
	}
	if b50ERowSpanArguments, b50ERowSpanExists := b50ERowSpanCanonical["arguments"]; b50ERowSpanExists {
		if !reflect.DeepEqual(b50ERowSpanOutput["arguments"], b50ERowSpanArguments) {
			b50ERowSpanT.Fatal("arguments changed")
		}
	} else if _, b50ERowSpanExists := b50ERowSpanOutput["arguments"]; b50ERowSpanExists {
		b50ERowSpanT.Fatal("absent arguments invented")
	}
	return b50ERowSpanOutput
}

func TestAb07Ao07SyntheticSpans(b50ERowSpanT *testing.T) {
	b50ERowSpanBytes, b50ERowSpanErr := os.ReadFile("testdata/ab07ao07_b50e_row_spans.json")
	if b50ERowSpanErr != nil {
		b50ERowSpanT.Fatal(b50ERowSpanErr)
	}
	var b50ERowSpanRows []struct {
		B50ERowSpanID     string         `json:"Ab07Ao07ID"`
		B50ERowSpanKey    int            `json:"Ab07Ao07Key"`
		B50ERowSpanValue  []int          `json:"Ab07Ao07Value"`
		B50ERowSpanFields map[string]any `json:"Ab07Ao07Fields"`
		B50ERowSpanState  string         `json:"Ab07Ao07State"`
	}
	if b50ERowSpanErr = json.Unmarshal(b50ERowSpanBytes, &b50ERowSpanRows); b50ERowSpanErr != nil {
		b50ERowSpanT.Fatal(b50ERowSpanErr)
	}
	if len(b50ERowSpanRows) != 57 {
		b50ERowSpanT.Fatal("missing synthetic controls")
	}
	for _, b50ERowSpanRow := range b50ERowSpanRows {
		b50ERowSpanT.Run(b50ERowSpanRow.B50ERowSpanID, func(b50ERowSpanT *testing.T) {
			b50ERowSpanOutput := b50ERowSpanCall(b50ERowSpanT, b50ERowSpanInput(b50ERowSpanRow.B50ERowSpanKey, b50ERowSpanRow.B50ERowSpanValue))
			if b50ERowSpanOutput["state"] != b50ERowSpanRow.B50ERowSpanState || !reflect.DeepEqual(b50ERowSpanOutput["fields"], b50ERowSpanRow.B50ERowSpanFields) {
				b50ERowSpanT.Fatalf("got %#v; want state %s fields %#v", b50ERowSpanOutput, b50ERowSpanRow.B50ERowSpanState, b50ERowSpanRow.B50ERowSpanFields)
			}
		})
	}
}

func TestAb07Ao07LabelsAndIndependentSpans(b50ERowSpanT *testing.T) {
	for _, b50ERowSpanCase := range []struct {
		b50ERowSpanKey                                      int
		b50ERowSpanValue                                    []int
		b50ERowSpanName, b50ERowSpanLabel, b50ERowSpanState string
	}{
		{8, []int{15}, "value", "true", "Hypothesis"}, {8, []int{240}, "value", "false", "Hypothesis"}, {8, []int{0}, "value", "", "Unknown"}, {8, []int{2}, "value", "", "Unknown"},
		{17, []int{16, 0, 0}, "status", "okay", "Hypothesis"}, {17, []int{16, 0, 85}, "status", "short-circuit", "Hypothesis"}, {17, []int{16, 0, 170}, "status", "interrupted", "Hypothesis"}, {17, []int{16, 0, 1}, "status", "", "Unknown"},
	} {
		b50ERowSpanOutput := b50ERowSpanCall(b50ERowSpanT, b50ERowSpanInput(b50ERowSpanCase.b50ERowSpanKey, b50ERowSpanCase.b50ERowSpanValue))
		b50ERowSpanField := b50ERowSpanOutput["fields"].(map[string]any)[b50ERowSpanCase.b50ERowSpanName].(map[string]any)
		if b50ERowSpanField["state"] != b50ERowSpanCase.b50ERowSpanState {
			b50ERowSpanT.Fatal(b50ERowSpanField)
		}
		if b50ERowSpanCase.b50ERowSpanLabel == "" {
			if _, b50ERowSpanExists := b50ERowSpanField["label"]; b50ERowSpanExists {
				b50ERowSpanT.Fatal("borrowed table")
			}
		} else if b50ERowSpanField["label"] != b50ERowSpanCase.b50ERowSpanLabel {
			b50ERowSpanT.Fatal(b50ERowSpanField)
		}
		if _, b50ERowSpanExists := b50ERowSpanField["number"]; b50ERowSpanExists {
			b50ERowSpanT.Fatal("invented numeric value")
		}
	}
	for _, b50ERowSpanValue := range [][]int{{16, 0, 18, 52}, {16, 0, 0, 0}, {16, 0, 255, 255}} {
		b50ERowSpanOutput := b50ERowSpanCall(b50ERowSpanT, b50ERowSpanInput(18, b50ERowSpanValue))
		if b50ERowSpanOutput["state"] != "Hypothesis" {
			b50ERowSpanT.Fatal("complement validation")
		}
		b50ERowSpanInverse := b50ERowSpanOutput["fields"].(map[string]any)["inverse"].(map[string]any)
		if !reflect.DeepEqual(b50ERowSpanInverse, map[string]any{"raw": []any{float64(b50ERowSpanValue[2]), float64(b50ERowSpanValue[3])}, "state": "Hypothesis"}) {
			b50ERowSpanT.Fatal("inverse bytes changed", b50ERowSpanInverse)
		}
	}
	for _, b50ERowSpanCase := range []struct {
		b50ERowSpanKey    int
		b50ERowSpanValue  []int
		b50ERowSpanNumber float64
	}{{2, []int{0}, 0}, {2, []int{255}, 255}, {3, []int{127}, 127}, {3, []int{128}, -128}, {3, []int{255}, -1}, {9, []int{42}, 42}, {10, []int{255}, -1}} {
		b50ERowSpanOutput := b50ERowSpanCall(b50ERowSpanT, b50ERowSpanInput(b50ERowSpanCase.b50ERowSpanKey, b50ERowSpanCase.b50ERowSpanValue))
		b50ERowSpanField := b50ERowSpanOutput["fields"].(map[string]any)["value"].(map[string]any)
		if b50ERowSpanField["number"] != b50ERowSpanCase.b50ERowSpanNumber || b50ERowSpanField["calculation"] != "Proven" || b50ERowSpanField["state"] != "Hypothesis" {
			b50ERowSpanT.Fatal(b50ERowSpanField)
		}
	}
}

func TestAb07Ao07ContextAxes(b50ERowSpanT *testing.T) {
	for _, b50ERowSpanCase := range []struct {
		b50ERowSpanName    string
		b50ERowSpanValue   any
		b50ERowSpanUnknown bool
	}{
		{"primary", 180, true}, {"secondary", 15, true}, {"operation", 3, true}, {"operation", 0, true}, {"direction", "negative", true}, {"view", "other", true}, {"family", map[string]any{"opaque": "other"}, false}, {"variant", []any{0, "other"}, false}, {"extra", map[string]any{"nested": []any{false, nil}}, false},
	} {
		b50ERowSpanT.Run(b50ERowSpanCase.b50ERowSpanName+"-"+b50ERowSpanCaseName(b50ERowSpanCase.b50ERowSpanValue), func(b50ERowSpanT *testing.T) {
			b50ERowSpanObject := b50ERowSpanInput(2, []int{42})
			b50ERowSpanObject["context"].(map[string]any)[b50ERowSpanCase.b50ERowSpanName] = b50ERowSpanCase.b50ERowSpanValue
			b50ERowSpanOutput := b50ERowSpanCall(b50ERowSpanT, b50ERowSpanObject)
			if b50ERowSpanCase.b50ERowSpanUnknown {
				if b50ERowSpanOutput["state"] != "Unknown" || len(b50ERowSpanOutput["fields"].(map[string]any)) != 0 {
					b50ERowSpanT.Fatal("selection guessed")
				}
			} else if b50ERowSpanOutput["state"] != "Hypothesis" {
				b50ERowSpanT.Fatal("opaque context selected view")
			}
		})
	}
	for _, b50ERowSpanName := range []string{"primary", "secondary", "operation", "direction", "view"} {
		b50ERowSpanObject := b50ERowSpanInput(2, []int{42})
		delete(b50ERowSpanObject["context"].(map[string]any), b50ERowSpanName)
		if b50ERowSpanCall(b50ERowSpanT, b50ERowSpanObject)["state"] != "Unknown" {
			b50ERowSpanT.Fatal("missing context selected")
		}
	}
}

func b50ERowSpanCaseName(b50ERowSpanValue any) string {
	b50ERowSpanBytes, _ := json.Marshal(b50ERowSpanValue)
	return string(b50ERowSpanBytes)
}

func TestAb07Ao07PresenceAndWidths(b50ERowSpanT *testing.T) {
	for _, b50ERowSpanName := range []string{"key", "value", "unmodelled"} {
		for _, b50ERowSpanMode := range []string{"missing", "empty", "zero", "extra"} {
			b50ERowSpanT.Run(b50ERowSpanName+b50ERowSpanMode, func(b50ERowSpanT *testing.T) {
				b50ERowSpanObject := b50ERowSpanInput(2, []int{42})
				b50ERowSpanRaw := b50ERowSpanObject["fields"].(map[string]any)
				switch b50ERowSpanMode {
				case "missing":
					delete(b50ERowSpanRaw, b50ERowSpanName)
				case "empty":
					b50ERowSpanRaw[b50ERowSpanName] = []int{}
				case "zero":
					b50ERowSpanRaw[b50ERowSpanName] = []int{0}
				case "extra":
					b50ERowSpanRaw[b50ERowSpanName] = []int{42, 43}
				}
				b50ERowSpanOutput := b50ERowSpanCall(b50ERowSpanT, b50ERowSpanObject)
				if b50ERowSpanName == "unmodelled" && (b50ERowSpanOutput["state"] != "Hypothesis" || b50ERowSpanOutput["fields"].(map[string]any)["value"].(map[string]any)["number"] != float64(42)) {
					b50ERowSpanT.Fatal("independent value lost")
				}
				switch b50ERowSpanName {
				case "key":
					if b50ERowSpanOutput["state"] != "Unknown" {
						b50ERowSpanT.Fatal("key guessed")
					}
				case "value":
					b50ERowSpanField := b50ERowSpanOutput["fields"].(map[string]any)["value"].(map[string]any)
					if b50ERowSpanMode == "missing" || b50ERowSpanMode == "empty" {
						if b50ERowSpanOutput["state"] != "partial" {
							b50ERowSpanT.Fatal("incomplete lost")
						}
						_, b50ERowSpanPresent := b50ERowSpanField["raw"]
						if b50ERowSpanPresent != (b50ERowSpanMode == "empty") {
							b50ERowSpanT.Fatal("missing/empty conflated")
						}
					} else if b50ERowSpanField["number"] != float64(b50ERowSpanRaw["value"].([]int)[0]) {
						b50ERowSpanT.Fatal("zero/value lost")
					}
				}
			})
		}
	}
	b50ERowSpanObject := b50ERowSpanInput(19, []int{42})
	if len(b50ERowSpanCall(b50ERowSpanT, b50ERowSpanObject)["fields"].(map[string]any)) != 0 {
		b50ERowSpanT.Fatal("unknown key interpreted")
	}
	b50ERowSpanObject = b50ERowSpanInput(2, []int{42})
	b50ERowSpanObject["arguments"] = map[string]any{}
	b50ERowSpanCall(b50ERowSpanT, b50ERowSpanObject)
}

func TestAb07Ao07RejectJSON(b50ERowSpanT *testing.T) {
	for _, b50ERowSpanInput := range []string{
		``, `null`, `[]`, `{}`, "{\"context\":{\"family\":\"" + string([]byte{255}) + "\"},\"fields\":{}}", `{"context":{},"fields":{}} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"x":{"a":1,"a":2}},"fields":{}}`, `{"context":{},"fields":{"x":[0],"x":[]}}`, `{"context":{},"fields":{"x":[1.0]}}`, `{"context":{},"fields":{"x":[1e0]}}`, `{"context":{},"fields":{"x":[-1]}}`, `{"context":{},"fields":{"x":[256]}}`, `{"context":{},"fields":{"x":[true]}}`, `{"context":{},"fields":{"x":[null]}}`, `{"context":{},"fields":{"x":["0"]}}`, `{"context":{},"fields":{"x":null}}`, `{"context":{},"fields":{"x":"AA=="}}`, `{"context":{},"fields":[]}`, `{"context":null,"fields":{}}`, `{"context":{"operation":"2"},"fields":{}}`, `{"context":{"primary":null},"fields":{}}`, `{"context":{"secondary":1.5},"fields":{}}`, `{"context":{"direction":2},"fields":{}}`, `{"context":{"view":null},"fields":{}}`, `{"context":{},"fields":{},"arguments":null}`, `{"context":{},"fields":{},"arguments":{"flag":true}}`, `{"context":{},"fields":{},"arguments":{"x":0,"x":0}}`, `{"context":{},"fields":{},"other":0}`, `{"context":{},"fields":{"x":[0,]}}`,
	} {
		b50ERowSpanOutput, b50ERowSpanErr := B50ERowSpanInterpret(json.RawMessage(b50ERowSpanInput))
		if b50ERowSpanErr == nil || b50ERowSpanOutput != nil {
			b50ERowSpanT.Fatalf("accepted %s", b50ERowSpanInput)
		}
	}
}

func TestAb07Ao07Ownership(b50ERowSpanT *testing.T) {
	b50ERowSpanObject := b50ERowSpanInput(17, []int{16, 0, 85})
	b50ERowSpanBytes, _ := json.Marshal(b50ERowSpanObject)
	b50ERowSpanOriginal := append([]byte(nil), b50ERowSpanBytes...)
	b50ERowSpanFirst, b50ERowSpanErr := B50ERowSpanInterpret(b50ERowSpanBytes)
	if b50ERowSpanErr != nil {
		b50ERowSpanT.Fatal(b50ERowSpanErr)
	}
	b50ERowSpanSecond, b50ERowSpanErr := B50ERowSpanInterpret(b50ERowSpanBytes)
	if b50ERowSpanErr != nil {
		b50ERowSpanT.Fatal(b50ERowSpanErr)
	}
	for b50ERowSpanIndex := range b50ERowSpanBytes {
		b50ERowSpanBytes[b50ERowSpanIndex] = ' '
	}
	if !json.Valid(b50ERowSpanFirst) || string(b50ERowSpanFirst) != string(b50ERowSpanSecond) {
		b50ERowSpanT.Fatal("input alias")
	}
	b50ERowSpanNested := map[string]any{}
	_ = json.Unmarshal(b50ERowSpanFirst, &b50ERowSpanNested)
	b50ERowSpanNested["raw"].(map[string]any)["value"].([]any)[0] = float64(99)
	if b50ERowSpanNested["fields"].(map[string]any)["value"].(map[string]any)["raw"].([]any)[0] != float64(16) {
		b50ERowSpanT.Fatal("nested alias")
	}
	b50ERowSpanNested["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "changed"
	for b50ERowSpanIndex := range b50ERowSpanFirst {
		b50ERowSpanFirst[b50ERowSpanIndex] = ' '
	}
	b50ERowSpanThird, b50ERowSpanErr := B50ERowSpanInterpret(b50ERowSpanOriginal)
	if b50ERowSpanErr != nil || string(b50ERowSpanThird) != string(b50ERowSpanSecond) || strings.Contains(string(b50ERowSpanThird), "changed") {
		b50ERowSpanT.Fatal("cross-call alias")
	}
}

func TestAb07Ao07CarrierDomains(b50ERowSpanT *testing.T) {
	for b50ERowSpanByte := 0; b50ERowSpanByte < 256; b50ERowSpanByte++ {
		for _, b50ERowSpanKey := range []int{1, 8, 17} {
			b50ERowSpanValue := []int{b50ERowSpanByte}
			b50ERowSpanName := "value"
			b50ERowSpanLabel := ""
			if b50ERowSpanKey == 8 {
				switch b50ERowSpanByte {
				case 15:
					b50ERowSpanLabel = "true"
				case 240:
					b50ERowSpanLabel = "false"
				}
			}
			if b50ERowSpanKey == 17 {
				b50ERowSpanValue = []int{16, 0, b50ERowSpanByte}
				b50ERowSpanName = "status"
				switch b50ERowSpanByte {
				case 0:
					b50ERowSpanLabel = "okay"
				case 85:
					b50ERowSpanLabel = "short-circuit"
				case 170:
					b50ERowSpanLabel = "interrupted"
				}
			}
			b50ERowSpanOutput := b50ERowSpanCall(b50ERowSpanT, b50ERowSpanInput(b50ERowSpanKey, b50ERowSpanValue))
			b50ERowSpanField := b50ERowSpanOutput["fields"].(map[string]any)[b50ERowSpanName].(map[string]any)
			if _, b50ERowSpanExists := b50ERowSpanField["number"]; b50ERowSpanExists {
				b50ERowSpanT.Fatal("numeric truth/status invented")
			}
			if b50ERowSpanLabel == "" {
				if b50ERowSpanField["state"] != "Unknown" || b50ERowSpanField["label"] != nil {
					b50ERowSpanT.Fatal("unstated mapping", b50ERowSpanKey, b50ERowSpanByte, b50ERowSpanField)
				}
			} else if b50ERowSpanField["state"] != "Hypothesis" || b50ERowSpanField["label"] != b50ERowSpanLabel {
				b50ERowSpanT.Fatal("mapping lost", b50ERowSpanField)
			}
		}
	}
}
