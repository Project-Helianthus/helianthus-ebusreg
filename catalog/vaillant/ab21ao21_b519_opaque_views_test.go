package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func ab21ao21Call(ab21ao21T *testing.T, ab21ao21Input any) map[string]any {
	ab21ao21T.Helper()
	ab21ao21JSON, ab21ao21Err := json.Marshal(ab21ao21Input)
	if ab21ao21Err != nil {
		ab21ao21T.Fatal(ab21ao21Err)
	}
	ab21ao21Result, ab21ao21Err := Ab21Ao21Interpret(ab21ao21JSON)
	if ab21ao21Err != nil {
		ab21ao21T.Fatal(ab21ao21Err)
	}
	ab21ao21Output := map[string]any{}
	if ab21ao21Err = json.Unmarshal(ab21ao21Result, &ab21ao21Output); ab21ao21Err != nil {
		ab21ao21T.Fatal(ab21ao21Err)
	}
	if ab21ao21Output["physical"] != "Unknown" || ab21ao21Output["qualified"] != false {
		ab21ao21T.Fatal(ab21ao21Output)
	}
	return ab21ao21Output
}

func ab21ao21Input(ab21ao21Operation int, ab21ao21Direction string, ab21ao21Bytes []int) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 25, "operation": ab21ao21Operation, "direction": ab21ao21Direction, "view": "logical-body", "family": map[string]any{"opaque": []any{"synthetic", 0}}, "variant": "synthetic"}, "fields": map[string]any{"logical_object": ab21ao21Bytes, "extra": []int{0, 255}}}
}

func TestAb21Ao21SelectorsWidths(ab21ao21T *testing.T) {
	for _, ab21ao21Operation := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 170} {
		for _, ab21ao21Direction := range []string{"request", "positive"} {
			for ab21ao21Width := 0; ab21ao21Width <= 5; ab21ao21Width++ {
				ab21ao21T.Run(fmt.Sprintf("op%d_%s_width%d", ab21ao21Operation, ab21ao21Direction, ab21ao21Width), func(ab21ao21T *testing.T) {
					ab21ao21Bytes := []int{19, 38, 57, 0, 255}[:ab21ao21Width]
					ab21ao21InputValue := ab21ao21Input(ab21ao21Operation, ab21ao21Direction, ab21ao21Bytes)
					ab21ao21Output := ab21ao21Call(ab21ao21T, ab21ao21InputValue)
					ab21ao21Raw, _ := json.Marshal(ab21ao21Output["raw"])
					ab21ao21Expected, _ := json.Marshal(ab21ao21InputValue["fields"])
					if !bytes.Equal(ab21ao21Raw, ab21ao21Expected) {
						ab21ao21T.Fatal("raw changed")
					}
					ab21ao21Fields := ab21ao21Output["fields"].(map[string]any)
					ab21ao21Unknown := ab21ao21Direction == "positive" && ab21ao21Operation != 3 && ab21ao21Operation != 6 && ab21ao21Operation != 170
					if ab21ao21Unknown {
						if ab21ao21Output["state"] != "Unknown" || len(ab21ao21Fields) != 0 {
							ab21ao21T.Fatal(ab21ao21Output)
						}
						return
					}
					if ab21ao21Width < 3 {
						if ab21ao21Output["state"] != "partial" {
							ab21ao21T.Fatal(ab21ao21Output)
						}
						return
					}
					ab21ao21Form := ab21ao21Direction == "positive" && ab21ao21Operation == 6 || ab21ao21Direction == "request" && (ab21ao21Operation == 0 || ab21ao21Operation == 2 || ab21ao21Operation == 170)
					if ab21ao21Form {
						ab21ao21Label := "no-declared-argument"
						if ab21ao21Direction == "positive" {
							ab21ao21Label = "prefix-only"
						}
						if ab21ao21Fields["form"].(map[string]any)["label"] != ab21ao21Label {
							ab21ao21T.Fatal(ab21ao21Output)
						}
						if ab21ao21Width > 3 {
							if ab21ao21Fields["unmodelled"] == nil || ab21ao21Output["state"] != "partial" {
								ab21ao21T.Fatal(ab21ao21Output)
							}
						}
						return
					}
					if ab21ao21Operation == 8 && ab21ao21Direction == "request" && ab21ao21Width == 3 {
						if ab21ao21Output["state"] != "partial" {
							ab21ao21T.Fatal(ab21ao21Output)
						}
						return
					}
					ab21ao21Offset := 3
					if ab21ao21Operation == 8 && ab21ao21Direction == "request" {
						ab21ao21Offset = 4
					}
					ab21ao21Want, _ := json.Marshal(ab21ao21Bytes[ab21ao21Offset:])
					ab21ao21Got, _ := json.Marshal(ab21ao21Fields["opaque"].(map[string]any)["raw"])
					if !bytes.Equal(ab21ao21Want, ab21ao21Got) {
						ab21ao21T.Fatal(string(ab21ao21Got), string(ab21ao21Want))
					}
				})
			}
		}
	}
}

func TestAb21Ao21EveryByte(ab21ao21T *testing.T) {
	for _, ab21ao21Direction := range []string{"request", "positive"} {
		for _, ab21ao21Operation := range []int{3, 170, 8} {
			for ab21ao21Position := 0; ab21ao21Position < 7; ab21ao21Position++ {
				for ab21ao21Byte := 0; ab21ao21Byte < 256; ab21ao21Byte++ {
					ab21ao21Bytes := []int{0, 0, 0, 127, 0, 255, 65}
					ab21ao21Bytes[ab21ao21Position] = ab21ao21Byte
					ab21ao21Output := ab21ao21Call(ab21ao21T, ab21ao21Input(ab21ao21Operation, ab21ao21Direction, ab21ao21Bytes))
					ab21ao21Fields := ab21ao21Output["fields"].(map[string]any)
					if ab21ao21Direction == "positive" && ab21ao21Operation == 8 {
						continue
					}
					if ab21ao21Direction == "request" && ab21ao21Operation == 170 {
						continue
					}
					ab21ao21Offset := 3
					if ab21ao21Operation == 8 {
						ab21ao21Offset = 4
						ab21ao21Address := ab21ao21Fields["slave_address"].(map[string]any)
						if !reflect.DeepEqual(ab21ao21Address["raw"], []any{float64(ab21ao21Bytes[3])}) {
							ab21ao21T.Fatal(ab21ao21Address)
						}
					}
					ab21ao21Opaque := ab21ao21Fields["opaque"].(map[string]any)
					ab21ao21Want := []any{}
					for _, ab21ao21Value := range ab21ao21Bytes[ab21ao21Offset:] {
						ab21ao21Want = append(ab21ao21Want, float64(ab21ao21Value))
					}
					if !reflect.DeepEqual(ab21ao21Opaque["raw"], ab21ao21Want) || ab21ao21Opaque["number"] != nil || ab21ao21Opaque["state"] != "Hypothesis" {
						ab21ao21T.Fatal(ab21ao21Opaque)
					}
					if !reflect.DeepEqual(ab21ao21Opaque["flags"], map[string]any{"accepted_empty": "Unknown", "unlimited_carrier": "Unknown"}) {
						ab21ao21T.Fatal(ab21ao21Opaque)
					}
				}
			}
		}
	}
}

func TestAb21Ao21ContextPresence(ab21ao21T *testing.T) {
	for ab21ao21Key, ab21ao21Value := range map[string]any{"primary": 180, "secondary": 24, "operation": 9, "direction": "response", "view": "guess"} {
		ab21ao21InputValue := ab21ao21Input(3, "request", []int{0, 0, 0, 0, 255, 65})
		ab21ao21InputValue["context"].(map[string]any)[ab21ao21Key] = ab21ao21Value
		if ab21ao21Output := ab21ao21Call(ab21ao21T, ab21ao21InputValue); ab21ao21Output["state"] != "Unknown" || len(ab21ao21Output["fields"].(map[string]any)) != 0 {
			ab21ao21T.Fatal(ab21ao21Output)
		}
		delete(ab21ao21InputValue["context"].(map[string]any), ab21ao21Key)
		if ab21ao21Call(ab21ao21T, ab21ao21InputValue)["state"] != "Unknown" {
			ab21ao21T.Fatal(ab21ao21Key)
		}
	}
	for _, ab21ao21Key := range []string{"family", "variant", "unmodelled_context"} {
		ab21ao21InputValue := ab21ao21Input(3, "request", []int{0, 0, 0, 0})
		ab21ao21InputValue["context"].(map[string]any)[ab21ao21Key] = map[string]any{"nested": []any{nil, false, "opaque", 0}}
		ab21ao21Output := ab21ao21Call(ab21ao21T, ab21ao21InputValue)
		ab21ao21Want, _ := json.Marshal(ab21ao21InputValue["context"])
		ab21ao21Got, _ := json.Marshal(ab21ao21Output["context"])
		if !bytes.Equal(ab21ao21Want, ab21ao21Got) || ab21ao21Output["state"] != "Hypothesis" {
			ab21ao21T.Fatal(ab21ao21Output)
		}
	}
	for _, ab21ao21Bytes := range [][]int{nil, {}, {0}} {
		ab21ao21InputValue := ab21ao21Input(3, "request", []int{})
		if ab21ao21Bytes == nil {
			delete(ab21ao21InputValue["fields"].(map[string]any), "logical_object")
		} else {
			ab21ao21InputValue["fields"].(map[string]any)["logical_object"] = ab21ao21Bytes
		}
		ab21ao21Output := ab21ao21Call(ab21ao21T, ab21ao21InputValue)
		if ab21ao21Output["state"] != "partial" {
			ab21ao21T.Fatal(ab21ao21Output)
		}
	}
	ab21ao21InputValue := ab21ao21Input(0, "request", []int{0, 0, 0})
	ab21ao21Without := ab21ao21Call(ab21ao21T, ab21ao21InputValue)
	ab21ao21InputValue["arguments"] = map[string]any{}
	ab21ao21With := ab21ao21Call(ab21ao21T, ab21ao21InputValue)
	if _, ab21ao21Present := ab21ao21Without["arguments"]; ab21ao21Present || ab21ao21With["arguments"] == nil {
		ab21ao21T.Fatal("argument presence lost")
	}
}

func TestAb21Ao21Invalid(ab21ao21T *testing.T) {
	for _, ab21ao21JSON := range []string{``, `{`, `[]`, `null`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"x":{"a":1,"a":2}},"fields":{}}`, `{"context":{},"fields":{"a":[],"a":[]}}`, `{"context":{},"fields":null}`, `{"context":[],"fields":{}}`, `{"context":{"operation":1.0},"fields":{}}`, `{"context":{"primary":true},"fields":{}}`, `{"context":{"secondary":null},"fields":{}}`, `{"context":{"direction":0},"fields":{}}`, `{"context":{"view":[]},"fields":{}}`, `{"context":{},"fields":{},"arguments":null}`, `{"context":{},"fields":{},"arguments":{"flag":false}}`} {
		if ab21ao21Result, ab21ao21Err := Ab21Ao21Interpret([]byte(ab21ao21JSON)); ab21ao21Err == nil || ab21ao21Result != nil {
			ab21ao21T.Fatal(ab21ao21JSON, string(ab21ao21Result), ab21ao21Err)
		}
	}
	for _, ab21ao21Bad := range []string{`null`, `"AA"`, `[null]`, `[true]`, `["0"]`, `[-1]`, `[256]`, `[1.5]`, `[1e0]`, `[[]]`, `[{}]`} {
		ab21ao21JSON := `{"context":{},"fields":{"other":` + ab21ao21Bad + `}}`
		if ab21ao21Result, ab21ao21Err := Ab21Ao21Interpret([]byte(ab21ao21JSON)); ab21ao21Err == nil || ab21ao21Result != nil {
			ab21ao21T.Fatal(ab21ao21JSON)
		}
	}
}

func TestAb21Ao21IsolationFixtures(ab21ao21T *testing.T) {
	ab21ao21Data, ab21ao21Err := os.ReadFile("testdata/ab21ao21_b519_opaque_views.json")
	if ab21ao21Err != nil {
		ab21ao21T.Fatal(ab21ao21Err)
	}
	var ab21ao21Cases []struct {
		Ab21Ao21ID     string          `json:"id"`
		Ab21Ao21Input  json.RawMessage `json:"input"`
		Ab21Ao21Fields map[string]any  `json:"fields"`
	}
	if ab21ao21Err = json.Unmarshal(ab21ao21Data, &ab21ao21Cases); ab21ao21Err != nil || len(ab21ao21Cases) == 0 {
		ab21ao21T.Fatal(ab21ao21Err)
	}
	for _, ab21ao21Case := range ab21ao21Cases {
		ab21ao21T.Run(ab21ao21Case.Ab21Ao21ID, func(ab21ao21T *testing.T) {
			ab21ao21InputCopy := append(json.RawMessage{}, ab21ao21Case.Ab21Ao21Input...)
			ab21ao21Result, ab21ao21Err := Ab21Ao21Interpret(ab21ao21InputCopy)
			if ab21ao21Err != nil {
				ab21ao21T.Fatal(ab21ao21Err)
			}
			ab21ao21Snapshot := append([]byte{}, ab21ao21Result...)
			for ab21ao21Index := range ab21ao21InputCopy {
				ab21ao21InputCopy[ab21ao21Index] = ' '
			}
			if !bytes.Equal(ab21ao21Snapshot, ab21ao21Result) {
				ab21ao21T.Fatal("input aliases result")
			}
			ab21ao21Output := map[string]any{}
			if ab21ao21Err = json.Unmarshal(ab21ao21Result, &ab21ao21Output); ab21ao21Err != nil {
				ab21ao21T.Fatal(ab21ao21Err)
			}
			if !reflect.DeepEqual(ab21ao21Output["fields"], ab21ao21Case.Ab21Ao21Fields) {
				ab21ao21T.Fatal(ab21ao21Output)
			}
			ab21ao21Output["raw"].(map[string]any)["logical_object"].([]any)[0] = 999
			for ab21ao21Index := range ab21ao21Result {
				ab21ao21Result[ab21ao21Index] = ' '
			}
			ab21ao21Again, ab21ao21Err := Ab21Ao21Interpret(ab21ao21Case.Ab21Ao21Input)
			if ab21ao21Err != nil || !bytes.Equal(ab21ao21Snapshot, ab21ao21Again) {
				ab21ao21T.Fatal("result aliases later call")
			}
		})
	}
}

func TestAb21Ao21PresenceAndSiblings(ab21ao21T *testing.T) {
	for _, ab21ao21Key := range []string{"logical_object", "extra"} {
		for _, ab21ao21Bytes := range [][]int{{}, {0}, {255}, {0, 255, 0, 255, 65, 1}} {
			ab21ao21Value := ab21ao21Input(8, "request", []int{91, 92, 93, 127, 0, 255})
			ab21ao21Value["fields"].(map[string]any)[ab21ao21Key] = ab21ao21Bytes
			ab21ao21Output := ab21ao21Call(ab21ao21T, ab21ao21Value)
			ab21ao21Want, _ := json.Marshal(ab21ao21Bytes)
			ab21ao21Got, _ := json.Marshal(ab21ao21Output["raw"].(map[string]any)[ab21ao21Key])
			if !bytes.Equal(ab21ao21Want, ab21ao21Got) {
				ab21ao21T.Fatal("raw width/value lost")
			}
			if ab21ao21Key == "extra" {
				if ab21ao21Output["fields"].(map[string]any)["opaque"] == nil || ab21ao21Output["fields"].(map[string]any)["slave_address"] == nil {
					ab21ao21T.Fatal("sibling erased")
				}
			}
		}
	}
	ab21ao21Missing := ab21ao21Input(3, "request", []int{})
	delete(ab21ao21Missing["fields"].(map[string]any), "logical_object")
	ab21ao21Absent := ab21ao21Call(ab21ao21T, ab21ao21Missing)
	ab21ao21Empty := ab21ao21Call(ab21ao21T, ab21ao21Input(3, "request", []int{}))
	ab21ao21Zero := ab21ao21Call(ab21ao21T, ab21ao21Input(3, "request", []int{0}))
	if reflect.DeepEqual(ab21ao21Absent, ab21ao21Empty) || reflect.DeepEqual(ab21ao21Empty, ab21ao21Zero) {
		ab21ao21T.Fatal("presence collapsed")
	}
	ab21ao21Complete := ab21ao21Call(ab21ao21T, ab21ao21Input(8, "request", []int{0, 0, 0, 0}))
	if !reflect.DeepEqual(ab21ao21Complete["fields"].(map[string]any)["slave_address"].(map[string]any)["raw"], []any{float64(0)}) || !reflect.DeepEqual(ab21ao21Complete["fields"].(map[string]any)["opaque"].(map[string]any)["raw"], []any{}) {
		ab21ao21T.Fatal(ab21ao21Complete)
	}
	for _, ab21ao21Operation := range []string{"99999999999999999999999999999999999999", "-99999999999999999999999999999999999999"} {
		ab21ao21JSON := []byte(`{"context":{"primary":181,"secondary":25,"operation":` + ab21ao21Operation + `,"direction":"request","view":"logical-body"},"fields":{"logical_object":[0,0,0]}}`)
		ab21ao21Result, ab21ao21Err := Ab21Ao21Interpret(ab21ao21JSON)
		if ab21ao21Err != nil {
			ab21ao21T.Fatal(ab21ao21Err)
		}
		ab21ao21Output := map[string]any{}
		if ab21ao21Err = json.Unmarshal(ab21ao21Result, &ab21ao21Output); ab21ao21Err != nil {
			ab21ao21T.Fatal(ab21ao21Err)
		}
		if ab21ao21Output["state"] != "Unknown" {
			ab21ao21T.Fatal(ab21ao21Output)
		}
	}
}
