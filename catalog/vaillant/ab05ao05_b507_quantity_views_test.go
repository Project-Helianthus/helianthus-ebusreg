package vaillant

import (
	ab05ao05bytes "bytes"
	ab05ao05json "encoding/json"
	ab05ao05fmt "fmt"
	ab05ao05os "os"
	ab05ao05reflect "reflect"
	ab05ao05testing "testing"
)

func ab05ao05Input(ab05ao05Operation int, ab05ao05View string, ab05ao05Fields map[string]any) map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 7, "operation": ab05ao05Operation, "direction": "positive", "view": ab05ao05View, "family": map[string]any{"opaque": []any{"synthetic-family"}}, "variant": []any{"synthetic-variant"}}, "fields": ab05ao05Fields}
}

func ab05ao05Run(ab05ao05T *ab05ao05testing.T, ab05ao05Object map[string]any) map[string]any {
	ab05ao05T.Helper()
	ab05ao05InputBytes, ab05ao05Error := ab05ao05json.Marshal(ab05ao05Object)
	if ab05ao05Error != nil {
		ab05ao05T.Fatal(ab05ao05Error)
	}
	ab05ao05Output, ab05ao05Error := Ab05Ao05Interpret(ab05ao05InputBytes)
	if ab05ao05Error != nil {
		ab05ao05T.Fatal(ab05ao05Error)
	}
	ab05ao05Result := map[string]any{}
	if ab05ao05Error := ab05ao05json.Unmarshal(ab05ao05Output, &ab05ao05Result); ab05ao05Error != nil {
		ab05ao05T.Fatal(ab05ao05Error)
	}
	if ab05ao05Result["physical"] != "Unknown" || ab05ao05Result["qualified"] != false {
		ab05ao05T.Fatalf("qualification: %s", ab05ao05Output)
	}
	ab05ao05Canonical := map[string]any{}
	if ab05ao05Error := ab05ao05json.Unmarshal(ab05ao05InputBytes, &ab05ao05Canonical); ab05ao05Error != nil {
		ab05ao05T.Fatal(ab05ao05Error)
	}
	if !ab05ao05reflect.DeepEqual(ab05ao05Result["context"], ab05ao05Canonical["context"]) || !ab05ao05reflect.DeepEqual(ab05ao05Result["raw"], ab05ao05Canonical["fields"]) {
		ab05ao05T.Fatalf("lost evidence: %s", ab05ao05Output)
	}
	for ab05ao05Key := range ab05ao05Result {
		switch ab05ao05Key {
		case "context", "raw", "arguments", "state", "fields", "physical", "qualified":
		default:
			ab05ao05T.Fatalf("unexpected output %s", ab05ao05Key)
		}
	}
	for _, ab05ao05Item := range ab05ao05Result["fields"].(map[string]any) {
		for ab05ao05Key := range ab05ao05Item.(map[string]any) {
			switch ab05ao05Key {
			case "raw", "state", "grade", "number", "label", "calculation", "candidate_domain":
			default:
				ab05ao05T.Fatalf("invented meaning or permission %s", ab05ao05Key)
			}
		}
	}
	return ab05ao05Result
}

func TestAb05Ao05Capacity(ab05ao05T *ab05ao05testing.T) {
	for ab05ao05Index, ab05ao05Byte := range []int{0, 1, 250, 251, 255} {
		ab05ao05T.Run(ab05ao05fmt.Sprint(ab05ao05Byte), func(ab05ao05T *ab05ao05testing.T) {
			ab05ao05Result := ab05ao05Run(ab05ao05T, ab05ao05Input(6, "capacity", map[string]any{"capacity": []int{ab05ao05Byte}, "cost": []int{1, 2}, "tail": []int{99}}))
			ab05ao05Field := ab05ao05Result["fields"].(map[string]any)["capacity"].(map[string]any)
			if ab05ao05Field["number"] != []float64{0, 0.1, 25, 25.1, 25.5}[ab05ao05Index] || ab05ao05Field["calculation"] != "Proven" || ab05ao05Field["grade"] != "Hypothesis" || ab05ao05Field["candidate_domain"] != (ab05ao05Byte <= 250) {
				ab05ao05T.Fatalf("capacity: %#v", ab05ao05Field)
			}
			ab05ao05State := "Hypothesis"
			if ab05ao05Byte > 250 {
				ab05ao05State = "outside-candidate-domain"
			}
			if ab05ao05Field["state"] != ab05ao05State || ab05ao05Result["state"] != "Hypothesis" {
				ab05ao05T.Fatal(ab05ao05Result)
			}
			for _, ab05ao05Key := range []string{"cost", "tail"} {
				ab05ao05Unknown := ab05ao05Result["fields"].(map[string]any)[ab05ao05Key].(map[string]any)
				if len(ab05ao05Unknown) != 2 || ab05ao05Unknown["state"] != "Unknown" {
					ab05ao05T.Fatal(ab05ao05Unknown)
				}
			}
		})
	}
}

func TestAb05Ao05Preference(ab05ao05T *ab05ao05testing.T) {
	for ab05ao05Index, ab05ao05Word := range [][]int{{0, 0}, {255, 255}, {0, 1}, {1, 0}} {
		ab05ao05T.Run(ab05ao05fmt.Sprint(ab05ao05Index), func(ab05ao05T *ab05ao05testing.T) {
			ab05ao05Result := ab05ao05Run(ab05ao05T, ab05ao05Input(8, "preference", map[string]any{"preference": ab05ao05Word}))
			ab05ao05Field := ab05ao05Result["fields"].(map[string]any)["preference"].(map[string]any)
			if _, ab05ao05Present := ab05ao05Field["number"]; ab05ao05Present {
				ab05ao05T.Fatal(ab05ao05Field)
			}
			if ab05ao05Field["grade"] != "Hypothesis" {
				ab05ao05T.Fatal(ab05ao05Field)
			}
			if ab05ao05Index < 2 {
				if ab05ao05Field["label"] != []string{"appliance-unavailable", "maximum-efficiency-preference"}[ab05ao05Index] || ab05ao05Field["state"] != "Hypothesis" {
					ab05ao05T.Fatal(ab05ao05Field)
				}
			} else if ab05ao05Field["label"] != nil || ab05ao05Field["state"] != "unmapped" || ab05ao05Result["state"] != "partial" {
				ab05ao05T.Fatal(ab05ao05Result)
			}
		})
	}
}

func TestAb05Ao05Selectors(ab05ao05T *ab05ao05testing.T) {
	for _, ab05ao05View := range []string{"capacity", "preference"} {
		ab05ao05Operation := 6
		ab05ao05Bytes := []int{1}
		if ab05ao05View == "preference" {
			ab05ao05Operation = 8
			ab05ao05Bytes = []int{255, 255}
		}
		for ab05ao05Key, ab05ao05Values := range map[string][]any{"primary": {180, 0, ab05ao05json.Number("999999999999999999999999999999"), nil}, "secondary": {8, 0, ab05ao05json.Number("999999999999999999999999999999"), nil}, "operation": {9, 11, 0, ab05ao05json.Number("999999999999999999999999999999"), nil}, "direction": {"negative", "", nil}, "view": {"unknown", "", nil}} {
			for ab05ao05Index, ab05ao05Value := range ab05ao05Values {
				ab05ao05T.Run(ab05ao05fmt.Sprintf("%s/%s/%d", ab05ao05View, ab05ao05Key, ab05ao05Index), func(ab05ao05T *ab05ao05testing.T) {
					ab05ao05Object := ab05ao05Input(ab05ao05Operation, ab05ao05View, map[string]any{ab05ao05View: ab05ao05Bytes})
					if ab05ao05Value == nil {
						delete(ab05ao05Object["context"].(map[string]any), ab05ao05Key)
					} else {
						ab05ao05Object["context"].(map[string]any)[ab05ao05Key] = ab05ao05Value
					}
					ab05ao05Result := ab05ao05Run(ab05ao05T, ab05ao05Object)
					ab05ao05Field := ab05ao05Result["fields"].(map[string]any)[ab05ao05View].(map[string]any)
					if ab05ao05Result["state"] != "Unknown" || ab05ao05Field["state"] != "Unknown" || len(ab05ao05Field) != 2 {
						ab05ao05T.Fatal(ab05ao05Result)
					}
				})
			}
		}
		for _, ab05ao05Key := range []string{"family", "variant"} {
			ab05ao05Object := ab05ao05Input(ab05ao05Operation, ab05ao05View, map[string]any{ab05ao05View: ab05ao05Bytes})
			ab05ao05Object["context"].(map[string]any)[ab05ao05Key] = map[string]any{"opaque": []any{nil, 0, "other-family"}}
			if ab05ao05Run(ab05ao05T, ab05ao05Object)["state"] != "Hypothesis" {
				ab05ao05T.Fatal("family selected view")
			}
		}
	}
	for _, ab05ao05Operation := range []int{9, 11} {
		for _, ab05ao05Word := range [][]int{{0, 0}, {255, 255}, {0, 1}} {
			ab05ao05Result := ab05ao05Run(ab05ao05T, ab05ao05Input(ab05ao05Operation, "preference", map[string]any{"preference": ab05ao05Word}))
			if ab05ao05Result["state"] != "Unknown" {
				ab05ao05T.Fatal(ab05ao05Result)
			}
		}
	}
	for _, ab05ao05Operation := range []int{6, 8} {
		ab05ao05View := "preference"
		if ab05ao05Operation == 8 {
			ab05ao05View = "capacity"
		}
		if ab05ao05Run(ab05ao05T, ab05ao05Input(ab05ao05Operation, ab05ao05View, map[string]any{"capacity": []int{1}, "preference": []int{0, 0}}))["state"] != "Unknown" {
			ab05ao05T.Fatal("crossed view")
		}
	}
}

func TestAb05Ao05WidthsPresenceAndSiblings(ab05ao05T *ab05ao05testing.T) {
	for _, ab05ao05View := range []string{"capacity", "preference"} {
		ab05ao05Operation := 6
		ab05ao05Widths := [][]int{nil, {}, {0, 1}, {0, 1, 2}}
		if ab05ao05View == "preference" {
			ab05ao05Operation = 8
			ab05ao05Widths = [][]int{nil, {}, {0}, {0, 1, 2}}
		}
		for ab05ao05Index, ab05ao05Bytes := range ab05ao05Widths {
			ab05ao05T.Run(ab05ao05fmt.Sprintf("%s/%d", ab05ao05View, ab05ao05Index), func(ab05ao05T *ab05ao05testing.T) {
				ab05ao05Fields := map[string]any{"tail": []int{0, 255}}
				if ab05ao05Index != 0 {
					ab05ao05Fields[ab05ao05View] = ab05ao05Bytes
				}
				ab05ao05Result := ab05ao05Run(ab05ao05T, ab05ao05Input(ab05ao05Operation, ab05ao05View, ab05ao05Fields))
				ab05ao05Field := ab05ao05Result["fields"].(map[string]any)[ab05ao05View].(map[string]any)
				_, ab05ao05Present := ab05ao05Field["raw"]
				if ab05ao05Present != (ab05ao05Index != 0) || ab05ao05Field["state"] != "incomplete" || ab05ao05Field["number"] != nil || ab05ao05Result["state"] != "partial" {
					ab05ao05T.Fatal(ab05ao05Result)
				}
			})
		}
	}
	for _, ab05ao05Key := range []string{"preference", "cost", "tail"} {
		for _, ab05ao05Bytes := range [][]int{{}, {0}, {255, 0, 1}} {
			ab05ao05Object := ab05ao05Input(6, "capacity", map[string]any{"capacity": []int{1}, ab05ao05Key: ab05ao05Bytes})
			if ab05ao05Run(ab05ao05T, ab05ao05Object)["fields"].(map[string]any)["capacity"].(map[string]any)["number"] != 0.1 {
				ab05ao05T.Fatal("sibling erased capacity")
			}
		}
	}
}

func TestAb05Ao05RejectInvalid(ab05ao05T *ab05ao05testing.T) {
	for ab05ao05Index, ab05ao05Input := range []string{
		``, `{`, `null`, `[]`, `{} {}`, `{"context":{},"fields":{},"context":{}}`,
		`{"context":{"family":{"x":1,"x":2}},"fields":{}}`,
		`{"context":{},"fields":{"x":[],"x":[]}}`,
		`{"context":{},"fields":{},"arguments":{"x":1,"x":2}}`,
		`{"context":null,"fields":{}}`, `{"context":{},"fields":null}`,
		`{"context":{},"fields":{"x":null}}`, `{"context":{},"fields":{"x":"AA=="}}`,
		`{"context":{},"fields":{"x":[-1]}}`, `{"context":{},"fields":{"x":[256]}}`,
		`{"context":{},"fields":{"x":[1.0]}}`, `{"context":{},"fields":{"x":[1e0]}}`,
		`{"context":{},"fields":{"x":[true]}}`, `{"context":{},"fields":{"x":["0"]}}`,
		`{"context":{},"fields":{"x":[null]}}`, `{"context":{},"fields":{"x":[{}]}}`,
		`{"context":{},"fields":{"x":[[]]}}`, `{"context":{},"fields":{},"arguments":null}`,
		`{"context":{},"fields":{},"arguments":{"execute":true}}`,
		`{"context":{},"fields":{},"extra":0}`, `{"context":{},"fields":{}} true`,
	} {
		ab05ao05T.Run(ab05ao05fmt.Sprint(ab05ao05Index), func(ab05ao05T *ab05ao05testing.T) {
			ab05ao05Output, ab05ao05Error := Ab05Ao05Interpret([]byte(ab05ao05Input))
			if ab05ao05Error == nil || ab05ao05Output != nil {
				ab05ao05T.Fatalf("accepted %s: %s, %v", ab05ao05Input, ab05ao05Output, ab05ao05Error)
			}
		})
	}
	for _, ab05ao05Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, ab05ao05Value := range []any{nil, true, []int{}, map[string]any{}, 1.5} {
			ab05ao05Object := ab05ao05Input(6, "capacity", map[string]any{"capacity": []int{1}})
			ab05ao05Object["context"].(map[string]any)[ab05ao05Key] = ab05ao05Value
			ab05ao05Bytes, _ := ab05ao05json.Marshal(ab05ao05Object)
			if ab05ao05Output, ab05ao05Error := Ab05Ao05Interpret(ab05ao05Bytes); ab05ao05Error == nil || ab05ao05Output != nil {
				ab05ao05T.Fatal("invalid context accepted", ab05ao05Key, ab05ao05Value)
			}
		}
	}
}

func TestAb05Ao05ContextKinds(ab05ao05T *ab05ao05testing.T) {
	for _, ab05ao05Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		ab05ao05Object := ab05ao05Input(6, "capacity", map[string]any{"capacity": []int{1}})
		if ab05ao05Key == "direction" || ab05ao05Key == "view" {
			ab05ao05Object["context"].(map[string]any)[ab05ao05Key] = 0
		} else {
			ab05ao05Object["context"].(map[string]any)[ab05ao05Key] = "6"
		}
		ab05ao05Bytes, _ := ab05ao05json.Marshal(ab05ao05Object)
		if ab05ao05Output, ab05ao05Error := Ab05Ao05Interpret(ab05ao05Bytes); ab05ao05Error == nil || ab05ao05Output != nil {
			ab05ao05T.Fatal("wrong context kind", ab05ao05Key)
		}
	}
}

func TestAb05Ao05OwnershipAndArguments(ab05ao05T *ab05ao05testing.T) {
	ab05ao05Object := ab05ao05Input(6, "capacity", map[string]any{"capacity": []int{1}})
	ab05ao05Absent := ab05ao05Run(ab05ao05T, ab05ao05Object)
	if _, ab05ao05Present := ab05ao05Absent["arguments"]; ab05ao05Present {
		ab05ao05T.Fatal("invented arguments")
	}
	ab05ao05Object["arguments"] = map[string]any{}
	ab05ao05Result := ab05ao05Run(ab05ao05T, ab05ao05Object)
	if !ab05ao05reflect.DeepEqual(ab05ao05Result["arguments"], map[string]any{}) {
		ab05ao05T.Fatal("lost empty arguments")
	}
	ab05ao05InputBytes, _ := ab05ao05json.Marshal(ab05ao05Object)
	ab05ao05Output, ab05ao05Error := Ab05Ao05Interpret(ab05ao05InputBytes)
	if ab05ao05Error != nil {
		ab05ao05T.Fatal(ab05ao05Error)
	}
	ab05ao05Saved := append([]byte(nil), ab05ao05Output...)
	for ab05ao05Index := range ab05ao05InputBytes {
		ab05ao05InputBytes[ab05ao05Index] = ' '
	}
	if !ab05ao05bytes.Equal(ab05ao05Saved, ab05ao05Output) {
		ab05ao05T.Fatal("input alias")
	}
	ab05ao05Result["raw"].(map[string]any)["capacity"].([]any)[0] = 255.0
	ab05ao05Result["context"].(map[string]any)["family"].(map[string]any)["opaque"].([]any)[0] = "mutated"
	ab05ao05Result["arguments"].(map[string]any)["execute"] = true
	if ab05ao05Result["fields"].(map[string]any)["capacity"].(map[string]any)["raw"].([]any)[0] != 1.0 {
		ab05ao05T.Fatal("nested alias")
	}
	ab05ao05Result["fields"].(map[string]any)["capacity"].(map[string]any)["raw"].([]any)[0] = 7.0
	ab05ao05Again := ab05ao05Run(ab05ao05T, ab05ao05Object)
	if ab05ao05Again["fields"].(map[string]any)["capacity"].(map[string]any)["number"] != 0.1 {
		ab05ao05T.Fatal("cross-call alias")
	}
	for ab05ao05Index := range ab05ao05Output {
		ab05ao05Output[ab05ao05Index] = ' '
	}
	ab05ao05FreshBytes, _ := ab05ao05json.Marshal(ab05ao05Object)
	ab05ao05Fresh, ab05ao05Error := Ab05Ao05Interpret(ab05ao05FreshBytes)
	if ab05ao05Error != nil || !ab05ao05bytes.Equal(ab05ao05Fresh, ab05ao05Saved) {
		ab05ao05T.Fatal("output alias")
	}
}

func TestAb05Ao05SyntheticFixture(ab05ao05T *ab05ao05testing.T) {
	ab05ao05Bytes, ab05ao05Error := ab05ao05os.ReadFile("testdata/ab05ao05_b507_quantity_views.json")
	if ab05ao05Error != nil {
		ab05ao05T.Fatal(ab05ao05Error)
	}
	ab05ao05Cases := map[string]map[string]any{}
	if ab05ao05Error := ab05ao05json.Unmarshal(ab05ao05Bytes, &ab05ao05Cases); ab05ao05Error != nil {
		ab05ao05T.Fatal(ab05ao05Error)
	}
	if len(ab05ao05Cases) != 8 {
		ab05ao05T.Fatal("fixture count")
	}
	for ab05ao05ID, ab05ao05Object := range ab05ao05Cases {
		ab05ao05T.Run(ab05ao05ID, func(ab05ao05T *ab05ao05testing.T) { ab05ao05Run(ab05ao05T, ab05ao05Object) })
	}
}
