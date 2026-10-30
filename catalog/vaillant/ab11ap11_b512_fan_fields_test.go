package vaillant

import (
	ab11ap11json "encoding/json"
	ab11ap11fmt "fmt"
	ab11ap11os "os"
	ab11ap11reflect "reflect"
	ab11ap11strings "strings"
	ab11ap11testing "testing"
)

func ab11ap11Subset(ab11ap11T *ab11ap11testing.T, ab11ap11Want, ab11ap11Got any) {
	ab11ap11T.Helper()
	if ab11ap11Map, ab11ap11OK := ab11ap11Want.(map[string]any); ab11ap11OK {
		ab11ap11Actual, ab11ap11OK := ab11ap11Got.(map[string]any)
		if !ab11ap11OK {
			ab11ap11T.Fatalf("expected object, got %v", ab11ap11Got)
		}
		for ab11ap11Key, ab11ap11Value := range ab11ap11Map {
			ab11ap11Subset(ab11ap11T, ab11ap11Value, ab11ap11Actual[ab11ap11Key])
		}
		return
	}
	if !ab11ap11reflect.DeepEqual(ab11ap11Want, ab11ap11Got) {
		ab11ap11T.Fatalf("want %#v, got %#v", ab11ap11Want, ab11ap11Got)
	}
}

func ab11ap11Call(ab11ap11T *ab11ap11testing.T, ab11ap11Input any) map[string]any {
	ab11ap11T.Helper()
	ab11ap11Bytes, ab11ap11Err := ab11ap11json.Marshal(ab11ap11Input)
	if ab11ap11Err != nil {
		ab11ap11T.Fatal(ab11ap11Err)
	}
	ab11ap11Result, ab11ap11Err := Ab11Ap11Evaluate(ab11ap11Bytes)
	if ab11ap11Err != nil {
		ab11ap11T.Fatal(ab11ap11Err)
	}
	ab11ap11Output := map[string]any{}
	if ab11ap11Err := ab11ap11json.Unmarshal(ab11ap11Result, &ab11ap11Output); ab11ap11Err != nil {
		ab11ap11T.Fatal(ab11ap11Err)
	}
	return ab11ap11Output
}

func TestAb11Ap11Synthetic(ab11ap11T *ab11ap11testing.T) {
	ab11ap11Data, ab11ap11Err := ab11ap11os.ReadFile("testdata/ab11ap11_b512_fan_fields.json")
	if ab11ap11Err != nil {
		ab11ap11T.Fatal(ab11ap11Err)
	}
	var ab11ap11Cases []struct {
		Ab11Ap11ID       string         `json:"id"`
		Ab11Ap11Input    map[string]any `json:"input"`
		Ab11Ap11Expected map[string]any `json:"expected"`
	}
	if ab11ap11Err := ab11ap11json.Unmarshal(ab11ap11Data, &ab11ap11Cases); ab11ap11Err != nil {
		ab11ap11T.Fatal(ab11ap11Err)
	}
	if len(ab11ap11Cases) < 100 {
		ab11ap11T.Fatal("missing independent cases")
	}
	for _, ab11ap11Case := range ab11ap11Cases {
		ab11ap11T.Run(ab11ap11Case.Ab11Ap11ID, func(ab11ap11T *ab11ap11testing.T) {
			ab11ap11Output := ab11ap11Call(ab11ap11T, ab11ap11Case.Ab11Ap11Input)
			ab11ap11Subset(ab11ap11T, ab11ap11Case.Ab11Ap11Expected, ab11ap11Output)
			if !ab11ap11reflect.DeepEqual(ab11ap11Output["context"], ab11ap11Case.Ab11Ap11Input["context"]) || !ab11ap11reflect.DeepEqual(ab11ap11Output["raw"], ab11ap11Case.Ab11Ap11Input["fields"]) {
				ab11ap11T.Fatal("raw/context changed")
			}
			for ab11ap11Key, ab11ap11Value := range ab11ap11Output["fields"].(map[string]any) {
				ab11ap11Field := ab11ap11Value.(map[string]any)
				if ab11ap11Field["state"] == "incomplete" || ab11ap11Field["state"] == "Unknown" {
					if _, ab11ap11OK := ab11ap11Field["number"]; ab11ap11OK {
						ab11ap11T.Fatal("invented number")
					}
					if _, ab11ap11OK := ab11ap11Field["label"]; ab11ap11OK {
						ab11ap11T.Fatal("invented label")
					}
				}
				ab11ap11Raw, ab11ap11Exists := ab11ap11Case.Ab11Ap11Input["fields"].(map[string]any)[ab11ap11Key]
				ab11ap11Actual, ab11ap11HasRaw := ab11ap11Field["raw"]
				if ab11ap11Exists != ab11ap11HasRaw || !ab11ap11reflect.DeepEqual(ab11ap11Raw, ab11ap11Actual) {
					ab11ap11T.Fatal("field raw distinction lost")
				}
			}
			ab11ap11ControlInput := ab11ap11Case.Ab11Ap11Input
			if ab11ap11Output["state"] == "partial" {
				for ab11ap11Key, ab11ap11Value := range ab11ap11Output["fields"].(map[string]any) {
					if ab11ap11Value.(map[string]any)["state"] == "incomplete" {
						continue
					}
					if ab11ap11ControlInput["fields"].(map[string]any)[ab11ap11Key] != nil && ab11ap11Value.(map[string]any)["state"] == nil {
						ab11ap11T.Fatal("sibling lost")
					}
				}
				if len(ab11ap11Output["fields"].(map[string]any)) < len(ab11ap11ControlInput["fields"].(map[string]any)) {
					ab11ap11T.Fatal("siblings dropped")
				}
			}
		})
	}
}

func ab11ap11Base() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 18, "operation": 14, "direction": "positive", "view": "air-quality", "family": "synthetic", "variant": "a", "opaque": map[string]any{"false": false, "null": nil, "array": []any{0, "x"}}}, "fields": map[string]any{"temperature": []int{16, 0}, "humidity": []int{80, 6}, "air_quality": []int{255, 255}}}
}

func TestAb11Ap11ContextAndArguments(ab11ap11T *ab11ap11testing.T) {
	for _, ab11ap11Key := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
		for _, ab11ap11Value := range []any{nil, 0, "other", "", false} {
			if (ab11ap11Key == "primary" || ab11ap11Key == "secondary" || ab11ap11Key == "operation") && (ab11ap11Value == "other" || ab11ap11Value == "" || ab11ap11Value == false) {
				continue
			}
			if (ab11ap11Key == "direction" || ab11ap11Key == "view") && (ab11ap11Value == 0 || ab11ap11Value == false) {
				continue
			}
			ab11ap11T.Run(ab11ap11Key+ab11ap11fmt.Sprint(ab11ap11Value), func(ab11ap11T *ab11ap11testing.T) {
				ab11ap11Input := ab11ap11Base()
				ab11ap11Input["context"].(map[string]any)[ab11ap11Key] = ab11ap11Value
				ab11ap11Output := ab11ap11Call(ab11ap11T, ab11ap11Input)
				ab11ap11Bytes, _ := ab11ap11json.Marshal(ab11ap11Input["context"])
				var ab11ap11Context any
				_ = ab11ap11json.Unmarshal(ab11ap11Bytes, &ab11ap11Context)
				if !ab11ap11reflect.DeepEqual(ab11ap11Context, ab11ap11Output["context"]) {
					ab11ap11T.Fatal("context changed")
				}
				if ab11ap11Key != "family" && ab11ap11Key != "variant" {
					if ab11ap11Output["state"] != "Unknown" {
						ab11ap11T.Fatal("unsupported admitted")
					}
					for _, ab11ap11Field := range ab11ap11Output["fields"].(map[string]any) {
						if len(ab11ap11Field.(map[string]any)) != 2 {
							ab11ap11T.Fatal("unsupported formula")
						}
					}
				} else if ab11ap11Output["state"] != "complete" {
					ab11ap11T.Fatal("opaque selector used")
				}
				if ab11ap11Output["qualified"] != false || ab11ap11Output["physical"] != "Unknown" {
					ab11ap11T.Fatal("qualified")
				}
			})
		}
		ab11ap11Input := ab11ap11Base()
		delete(ab11ap11Input["context"].(map[string]any), ab11ap11Key)
		ab11ap11Output := ab11ap11Call(ab11ap11T, ab11ap11Input)
		if _, ab11ap11Exists := ab11ap11Output["context"].(map[string]any)[ab11ap11Key]; ab11ap11Exists {
			ab11ap11T.Fatal("missing context filled")
		}
	}
	for _, ab11ap11Arguments := range []any{map[string]any{}, map[string]any{"opaque": false}, map[string]any{"opaque": 0}, map[string]any{"opaque": nil}, map[string]any{"opaque": []any{false, 0, nil, "x"}}} {
		ab11ap11Input := ab11ap11Base()
		ab11ap11Input["arguments"] = ab11ap11Arguments
		ab11ap11Output := ab11ap11Call(ab11ap11T, ab11ap11Input)
		ab11ap11Bytes, _ := ab11ap11json.Marshal(ab11ap11Arguments)
		var ab11ap11Want any
		_ = ab11ap11json.Unmarshal(ab11ap11Bytes, &ab11ap11Want)
		if !ab11ap11reflect.DeepEqual(ab11ap11Want, ab11ap11Output["arguments"]) {
			ab11ap11T.Fatal("arguments changed")
		}
	}
	if _, ab11ap11Exists := ab11ap11Call(ab11ap11T, ab11ap11Base())["arguments"]; ab11ap11Exists {
		ab11ap11T.Fatal("absent arguments filled")
	}
}

func TestAb11Ap11Reject(ab11ap11T *ab11ap11testing.T) {
	for ab11ap11Index, ab11ap11Input := range []string{``, `{`, `[]`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"x":{"a":0,"a":1}},"fields":{}}`, `{"context":{},"fields":{"x":[{"a":0,"a":1}]}}`, `{"context":{},"fields":{},"arguments":{"x":[{"a":0,"a":1}]}}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":null}`, `{"context":{},"fields":{},"arguments":false}`} {
		ab11ap11T.Run(ab11ap11fmt.Sprint(ab11ap11Index), func(ab11ap11T *ab11ap11testing.T) {
			ab11ap11Output, ab11ap11Err := Ab11Ap11Evaluate([]byte(ab11ap11Input))
			if ab11ap11Err == nil || ab11ap11Output != nil {
				ab11ap11T.Fatal("invalid accepted")
			}
		})
	}
	for _, ab11ap11Bad := range []string{`true`, `false`, `null`, `"0"`, `-1`, `256`, `1.5`, `1e0`, `{}`, `[]`} {
		for _, ab11ap11Template := range []string{`{"context":{},"fields":{"x":[VALUE]}}`, `{"context":{"primary":VALUE},"fields":{}}`, `{"context":{"secondary":VALUE},"fields":{}}`, `{"context":{"operation":VALUE},"fields":{}}`} {
			if ab11ap11Bad == `null` && !ab11ap11strings.Contains(ab11ap11Template, `[VALUE]`) {
				continue
			}
			ab11ap11Output, ab11ap11Err := Ab11Ap11Evaluate([]byte(ab11ap11strings.ReplaceAll(ab11ap11Template, "VALUE", ab11ap11Bad)))
			if ab11ap11Err == nil || ab11ap11Output != nil {
				ab11ap11T.Fatalf("accepted %s %s", ab11ap11Template, ab11ap11Bad)
			}
		}
	}
	for _, ab11ap11Bad := range []any{nil, false, 0, "x", map[string]any{}} {
		ab11ap11Input := ab11ap11Base()
		ab11ap11Input["fields"].(map[string]any)["humidity"] = ab11ap11Bad
		ab11ap11Bytes, _ := ab11ap11json.Marshal(ab11ap11Input)
		if ab11ap11Output, ab11ap11Err := Ab11Ap11Evaluate(ab11ap11Bytes); ab11ap11Err == nil || ab11ap11Output != nil {
			ab11ap11T.Fatal("non-array accepted")
		}
	}
}

func TestAb11Ap11Isolation(ab11ap11T *ab11ap11testing.T) {
	ab11ap11Input, _ := ab11ap11json.Marshal(ab11ap11Base())
	ab11ap11Original := append([]byte(nil), ab11ap11Input...)
	ab11ap11First, ab11ap11Err := Ab11Ap11Evaluate(ab11ap11Input)
	if ab11ap11Err != nil {
		ab11ap11T.Fatal(ab11ap11Err)
	}
	ab11ap11Saved := append([]byte(nil), ab11ap11First...)
	for ab11ap11Index := range ab11ap11Input {
		ab11ap11Input[ab11ap11Index] = ' '
	}
	if !ab11ap11reflect.DeepEqual(ab11ap11First, ab11ap11json.RawMessage(ab11ap11Saved)) {
		ab11ap11T.Fatal("input alias")
	}
	var ab11ap11Decoded map[string]any
	_ = ab11ap11json.Unmarshal(ab11ap11First, &ab11ap11Decoded)
	ab11ap11Decoded["raw"].(map[string]any)["humidity"].([]any)[0] = 0.0
	if ab11ap11Decoded["fields"].(map[string]any)["humidity"].(map[string]any)["raw"].([]any)[0] != float64(80) {
		ab11ap11T.Fatal("nested alias")
	}
	for ab11ap11Index := range ab11ap11First {
		ab11ap11First[ab11ap11Index] = ' '
	}
	ab11ap11Second, ab11ap11Err := Ab11Ap11Evaluate(ab11ap11Original)
	if ab11ap11Err != nil || !ab11ap11reflect.DeepEqual([]byte(ab11ap11Second), ab11ap11Saved) {
		ab11ap11T.Fatal("later result affected")
	}
}

func TestAb11Ap11RequestIsolation(ab11ap11T *ab11ap11testing.T) {
	for _, ab11ap11View := range []string{"digital-air", "air-quality"} {
		ab11ap11Input := ab11ap11Base()
		ab11ap11Context := ab11ap11Input["context"].(map[string]any)
		ab11ap11Context["view"] = ab11ap11View
		ab11ap11Context["direction"] = "request"
		if ab11ap11View == "digital-air" {
			ab11ap11Context["operation"] = 13
		}
		ab11ap11Output := ab11ap11Call(ab11ap11T, ab11ap11Input)
		if ab11ap11Output["state"] != "complete" {
			ab11ap11T.Fatal("request incomplete")
		}
		for _, ab11ap11Field := range ab11ap11Output["fields"].(map[string]any) {
			if len(ab11ap11Field.(map[string]any)) != 2 {
				ab11ap11T.Fatal("positive semantics on request")
			}
		}
	}
	for _, ab11ap11Key := range []string{"direction", "view"} {
		for _, ab11ap11Bad := range []any{false, 0, []any{}, map[string]any{}} {
			ab11ap11Input := ab11ap11Base()
			ab11ap11Input["context"].(map[string]any)[ab11ap11Key] = ab11ap11Bad
			ab11ap11Bytes, _ := ab11ap11json.Marshal(ab11ap11Input)
			if ab11ap11Result, ab11ap11Err := Ab11Ap11Evaluate(ab11ap11Bytes); ab11ap11Err == nil || ab11ap11Result != nil {
				ab11ap11T.Fatal("wrong context kind admitted")
			}
		}
	}
}
