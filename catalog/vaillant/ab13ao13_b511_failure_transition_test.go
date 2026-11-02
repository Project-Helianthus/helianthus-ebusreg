package vaillant

import (
	ab13ao13json "encoding/json"
	ab13ao13os "os"
	ab13ao13reflect "reflect"
	ab13ao13strings "strings"
	ab13ao13testing "testing"
)

func ab13ao13Base() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 17, "operation": 1, "direction": "local", "view": "refresh", "family": map[string]any{"opaque": "synthetic-family"}, "variant": []any{"synthetic-variant"}}, "fields": map[string]any{"registration": []int{1}, "prior_failure": []int{1}, "local_body": []int{10, 20, 30, 40, 50, 60, 70, 80, 0}, "unmodelled": []int{255, 0}}, "arguments": map[string]any{"index": 0, "selector": 2, "local_selector": 3, "count": 1, "parsed": true}}
}

func ab13ao13Call(ab13ao13T *ab13ao13testing.T, ab13ao13Input map[string]any) map[string]any {
	ab13ao13T.Helper()
	ab13ao13Encoded, ab13ao13Err := ab13ao13json.Marshal(ab13ao13Input)
	if ab13ao13Err != nil {
		ab13ao13T.Fatal(ab13ao13Err)
	}
	ab13ao13Result, ab13ao13Err := Ab13Ao13Interpret(ab13ao13Encoded)
	if ab13ao13Err != nil {
		ab13ao13T.Fatal(ab13ao13Err)
	}
	ab13ao13Output := map[string]any{}
	if ab13ao13Err := ab13ao13json.Unmarshal(ab13ao13Result, &ab13ao13Output); ab13ao13Err != nil {
		ab13ao13T.Fatal(ab13ao13Err)
	}
	if ab13ao13Output["physical"] != "Unknown" || ab13ao13Output["qualified"] != false {
		ab13ao13T.Fatal("native qualification inferred")
	}
	ab13ao13Roundtrip, _ := ab13ao13json.Marshal(ab13ao13Output)
	ab13ao13Again := map[string]any{}
	_ = ab13ao13json.Unmarshal(ab13ao13Roundtrip, &ab13ao13Again)
	if !ab13ao13reflect.DeepEqual(ab13ao13Output, ab13ao13Again) {
		ab13ao13T.Fatal("roundtrip")
	}
	ab13ao13InputRoundtrip := map[string]any{}
	_ = ab13ao13json.Unmarshal(ab13ao13Encoded, &ab13ao13InputRoundtrip)
	for _, ab13ao13Pair := range [][2]string{{"context", "context"}, {"fields", "raw"}, {"arguments", "arguments"}} {
		if !ab13ao13reflect.DeepEqual(ab13ao13InputRoundtrip[ab13ao13Pair[0]], ab13ao13Output[ab13ao13Pair[1]]) {
			ab13ao13T.Fatalf("preservation %s", ab13ao13Pair[0])
		}
	}
	return ab13ao13Output
}

func ab13ao13CheckFailure(ab13ao13T *ab13ao13testing.T, ab13ao13Output map[string]any, ab13ao13Want bool, ab13ao13State string, ab13ao13Copied bool) {
	ab13ao13T.Helper()
	if ab13ao13Output["state"] != ab13ao13State {
		ab13ao13T.Fatalf("state: %v", ab13ao13Output)
	}
	ab13ao13Fields := ab13ao13Output["fields"].(map[string]any)
	ab13ao13Failure := ab13ao13Fields["failure"].(map[string]any)
	if ab13ao13Failure["state"] != "Hypothesis" || ab13ao13Failure["flags"].(map[string]any)["failure"] != ab13ao13Want {
		ab13ao13T.Fatalf("failure: %v", ab13ao13Failure)
	}
	if _, ab13ao13Exists := ab13ao13Failure["number"]; ab13ao13Exists {
		ab13ao13T.Fatal("invented failure number")
	}
	ab13ao13Copy, ab13ao13Exists := ab13ao13Fields["copied"]
	if ab13ao13Exists != ab13ao13Copied {
		ab13ao13T.Fatalf("copy presence: %v", ab13ao13Fields)
	}
	if ab13ao13Copied {
		ab13ao13Raw := ab13ao13Output["raw"].(map[string]any)["local_body"].([]any)
		if !ab13ao13reflect.DeepEqual(ab13ao13Copy.(map[string]any)["raw"], ab13ao13Raw[:9]) || len(ab13ao13Copy.(map[string]any)["raw"].([]any)) != 9 {
			ab13ao13T.Fatal("copy must be exactly first nine")
		}
	}
}

func TestAb13Ao13Refresh(ab13ao13T *ab13ao13testing.T) {
	for _, ab13ao13Case := range []struct {
		ab13ao13Name  string
		ab13ao13Group string
		ab13ao13Key   string
		ab13ao13Value any
		ab13ao13Want  bool
		ab13ao13Copy  bool
	}{
		{"ab13ao13_control", "arguments", "index", 0, false, true},
		{"ab13ao13_registration_bit", "fields", "registration", []int{2}, true, false},
		{"ab13ao13_registration_zero", "fields", "registration", []int{0}, true, false},
		{"ab13ao13_registration_high", "fields", "registration", []int{3}, false, true},
		{"ab13ao13_index_exclusion", "arguments", "index", 1, true, false},
		{"ab13ao13_index_two", "arguments", "index", 2, false, true},
		{"ab13ao13_selector_exclusion", "arguments", "selector", 3, true, false},
		{"ab13ao13_local_selector_exclusion", "arguments", "local_selector", 2, true, false},
		{"ab13ao13_selector_other", "arguments", "selector", 4, false, true},
		{"ab13ao13_local_selector_other", "arguments", "local_selector", 4, false, true},
		{"ab13ao13_count_zero", "arguments", "count", 0, true, false},
		{"ab13ao13_count_two", "arguments", "count", 2, false, true},
		{"ab13ao13_parse_false", "arguments", "parsed", false, true, false},
		{"ab13ao13_prior_zero", "fields", "prior_failure", []int{0}, false, true},
	} {
		ab13ao13T.Run(ab13ao13Case.ab13ao13Name, func(ab13ao13T *ab13ao13testing.T) {
			ab13ao13Input := ab13ao13Base()
			ab13ao13Input[ab13ao13Case.ab13ao13Group].(map[string]any)[ab13ao13Case.ab13ao13Key] = ab13ao13Case.ab13ao13Value
			ab13ao13CheckFailure(ab13ao13T, ab13ao13Call(ab13ao13T, ab13ao13Input), ab13ao13Case.ab13ao13Want, "Hypothesis", ab13ao13Case.ab13ao13Copy)
		})
	}
	for ab13ao13Last := 0; ab13ao13Last < 3; ab13ao13Last++ {
		ab13ao13Input := ab13ao13Base()
		ab13ao13Input["fields"].(map[string]any)["local_body"].([]int)[8] = ab13ao13Last
		ab13ao13CheckFailure(ab13ao13T, ab13ao13Call(ab13ao13T, ab13ao13Input), ab13ao13Last == 1, "Hypothesis", true)
	}
	for ab13ao13Length := 0; ab13ao13Length <= 10; ab13ao13Length++ {
		ab13ao13Input := ab13ao13Base()
		ab13ao13Body := make([]int, ab13ao13Length)
		for ab13ao13Index := range ab13ao13Body {
			ab13ao13Body[ab13ao13Index] = ab13ao13Index + 10
		}
		ab13ao13Input["fields"].(map[string]any)["local_body"] = ab13ao13Body
		ab13ao13State := "Hypothesis"
		if ab13ao13Length < 9 {
			ab13ao13State = "partial"
		}
		ab13ao13CheckFailure(ab13ao13T, ab13ao13Call(ab13ao13T, ab13ao13Input), ab13ao13Length < 9, ab13ao13State, ab13ao13Length >= 9)
	}
}

func TestAb13Ao13LocalViews(ab13ao13T *ab13ao13testing.T) {
	for _, ab13ao13Prior := range []int{0, 1, 2, 255} {
		ab13ao13Input := ab13ao13Base()
		ab13ao13Context := ab13ao13Input["context"].(map[string]any)
		delete(ab13ao13Context, "operation")
		ab13ao13Context["view"] = "timer-consume"
		ab13ao13Input["fields"] = map[string]any{"prior_failure": []int{ab13ao13Prior}}
		ab13ao13Input["arguments"] = map[string]any{"index": 0}
		ab13ao13Output := ab13ao13Call(ab13ao13T, ab13ao13Input)
		ab13ao13CheckFailure(ab13ao13T, ab13ao13Output, ab13ao13Prior != 0, "Hypothesis", false)
		ab13ao13Path := ab13ao13Output["fields"].(map[string]any)["path"].(map[string]any)
		if ab13ao13Prior == 0 {
			if ab13ao13Path["label"] != "capability-path" || ab13ao13Path["number"] != nil {
				ab13ao13T.Fatal(ab13ao13Path)
			}
		} else if ab13ao13Path["label"] != "error" || ab13ao13Path["number"] != float64(120) {
			ab13ao13T.Fatal(ab13ao13Path)
		}
		ab13ao13Input["arguments"].(map[string]any)["index"] = 1
		ab13ao13Output = ab13ao13Call(ab13ao13T, ab13ao13Input)
		if ab13ao13Output["state"] != "Unknown" || len(ab13ao13Output["fields"].(map[string]any)) != 0 {
			ab13ao13T.Fatal(ab13ao13Output)
		}
	}
	ab13ao13Input := ab13ao13Base()
	delete(ab13ao13Input["context"].(map[string]any), "operation")
	ab13ao13Input["context"].(map[string]any)["view"] = "initialize"
	ab13ao13Input["fields"] = map[string]any{}
	delete(ab13ao13Input, "arguments")
	ab13ao13CheckFailure(ab13ao13T, ab13ao13Call(ab13ao13T, ab13ao13Input), false, "Hypothesis", false)
}

func TestAb13Ao13IdentityAndPresence(ab13ao13T *ab13ao13testing.T) {
	for _, ab13ao13Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, ab13ao13Missing := range []bool{false, true} {
			ab13ao13Input := ab13ao13Base()
			ab13ao13Context := ab13ao13Input["context"].(map[string]any)
			if ab13ao13Missing {
				delete(ab13ao13Context, ab13ao13Key)
			} else if ab13ao13Key == "direction" || ab13ao13Key == "view" {
				ab13ao13Context[ab13ao13Key] = "unrecognised"
			} else {
				ab13ao13Context[ab13ao13Key] = 0
			}
			ab13ao13Output := ab13ao13Call(ab13ao13T, ab13ao13Input)
			if ab13ao13Output["state"] != "Unknown" || len(ab13ao13Output["fields"].(map[string]any)) != 0 {
				ab13ao13T.Fatalf("%s: %v", ab13ao13Key, ab13ao13Output)
			}
		}
	}
	for _, ab13ao13Key := range []string{"family", "variant"} {
		ab13ao13Input := ab13ao13Base()
		ab13ao13Input["context"].(map[string]any)[ab13ao13Key] = map[string]any{"opaque": []any{0, false, "other-family"}}
		ab13ao13CheckFailure(ab13ao13T, ab13ao13Call(ab13ao13T, ab13ao13Input), false, "Hypothesis", true)
	}
	for _, ab13ao13Key := range []string{"registration", "prior_failure", "local_body"} {
		for _, ab13ao13Width := range []int{-1, 0, 2} {
			ab13ao13Input := ab13ao13Base()
			if ab13ao13Width < 0 {
				delete(ab13ao13Input["fields"].(map[string]any), ab13ao13Key)
			} else {
				ab13ao13Input["fields"].(map[string]any)[ab13ao13Key] = make([]int, ab13ao13Width)
			}
			ab13ao13Output := ab13ao13Call(ab13ao13T, ab13ao13Input)
			if ab13ao13Output["state"] != "partial" {
				ab13ao13T.Fatal(ab13ao13Output)
			}
			if ab13ao13Key != "prior_failure" {
				ab13ao13CheckFailure(ab13ao13T, ab13ao13Output, true, "partial", false)
			} else {
				ab13ao13CheckFailure(ab13ao13T, ab13ao13Output, false, "partial", true)
			}
		}
	}
	for _, ab13ao13Key := range []string{"index", "selector", "local_selector", "count", "parsed"} {
		ab13ao13Input := ab13ao13Base()
		delete(ab13ao13Input["arguments"].(map[string]any), ab13ao13Key)
		ab13ao13CheckFailure(ab13ao13T, ab13ao13Call(ab13ao13T, ab13ao13Input), true, "partial", false)
	}
	for _, ab13ao13Value := range []any{[]int{}, []int{0}, []int{4, 5, 6}} {
		ab13ao13Input := ab13ao13Base()
		ab13ao13Input["fields"].(map[string]any)["unmodelled"] = ab13ao13Value
		ab13ao13CheckFailure(ab13ao13T, ab13ao13Call(ab13ao13T, ab13ao13Input), false, "Hypothesis", true)
	}
}

func TestAb13Ao13Malformed(ab13ao13T *ab13ao13testing.T) {
	for _, ab13ao13Text := range []string{``, `{`, `[]`, `null`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"opaque":{"x":0,"x":1}},"fields":{}}`, `{"context":{},"fields":{"x":[256]}}`, `{"context":{},"fields":{"x":[-1]}}`, `{"context":{},"fields":{"x":[1.0]}}`, `{"context":{},"fields":{"x":[true]}}`, `{"context":{},"fields":{"x":[null]}}`, `{"context":{},"fields":{"x":"AA=="}}`, `{"context":{},"fields":{"x":null}}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":[]}`, `{"context":{},"fields":{},"arguments":null}`, `{"context":{"primary":"181"},"fields":{}}`, `{"context":{"view":1},"fields":{}}`, `{"context":{},"fields":{},"arguments":{"parsed":0}}`, `{"context":{},"fields":{},"arguments":{"index":false}}`, `{"context":{},"fields":{},"arguments":{"count":-1}}`, `{"context":{},"fields":{},"arguments":{"extra":0}}`, `{"context":{},"fields":{},"arguments":{"count":0,"count":1}}`} {
		ab13ao13Result, ab13ao13Err := Ab13Ao13Interpret([]byte(ab13ao13Text))
		if ab13ao13Err == nil || ab13ao13Result != nil {
			ab13ao13T.Fatalf("accepted %s: %s", ab13ao13Text, ab13ao13Result)
		}
	}
}

func TestAb13Ao13Isolation(ab13ao13T *ab13ao13testing.T) {
	ab13ao13Input, _ := ab13ao13json.Marshal(ab13ao13Base())
	ab13ao13Original := append([]byte(nil), ab13ao13Input...)
	ab13ao13First, ab13ao13Err := Ab13Ao13Interpret(ab13ao13Input)
	if ab13ao13Err != nil {
		ab13ao13T.Fatal(ab13ao13Err)
	}
	ab13ao13Expected := append([]byte(nil), ab13ao13First...)
	for ab13ao13Index := range ab13ao13Input {
		ab13ao13Input[ab13ao13Index] = ' '
	}
	if string(ab13ao13First) != string(ab13ao13Expected) {
		ab13ao13T.Fatal("input alias")
	}
	ab13ao13Decoded := map[string]any{}
	_ = ab13ao13json.Unmarshal(ab13ao13First, &ab13ao13Decoded)
	ab13ao13Decoded["context"].(map[string]any)["family"].(map[string]any)["opaque"] = "changed"
	ab13ao13Decoded["raw"].(map[string]any)["local_body"].([]any)[0] = float64(255)
	ab13ao13Decoded["arguments"].(map[string]any)["index"] = 99
	ab13ao13Decoded["fields"].(map[string]any)["copied"].(map[string]any)["raw"].([]any)[0] = 0
	if ab13ao13Decoded["fields"].(map[string]any)["failure"].(map[string]any)["raw"].([]any)[0] != float64(10) {
		ab13ao13T.Fatal("nested result alias")
	}
	for ab13ao13Index := range ab13ao13First {
		ab13ao13First[ab13ao13Index] = ' '
	}
	ab13ao13Second, ab13ao13Err := Ab13Ao13Interpret(ab13ao13Original)
	if ab13ao13Err != nil || string(ab13ao13Second) != string(ab13ao13Expected) {
		ab13ao13T.Fatal("call alias", ab13ao13Err)
	}
}

func TestAb13Ao13Fixture(ab13ao13T *ab13ao13testing.T) {
	ab13ao13Data, ab13ao13Err := ab13ao13os.ReadFile("testdata/ab13ao13_b511_failure_transition.json")
	if ab13ao13Err != nil {
		ab13ao13T.Fatal(ab13ao13Err)
	}
	var ab13ao13Cases []struct {
		Ab13Ao13ID      string         `json:"id"`
		Ab13Ao13Input   map[string]any `json:"input"`
		Ab13Ao13Failure bool           `json:"failure"`
	}
	if ab13ao13Err := ab13ao13json.Unmarshal(ab13ao13Data, &ab13ao13Cases); ab13ao13Err != nil || len(ab13ao13Cases) == 0 {
		ab13ao13T.Fatal("fixture", ab13ao13Err)
	}
	for _, ab13ao13Case := range ab13ao13Cases {
		if !ab13ao13strings.HasPrefix(ab13ao13Case.Ab13Ao13ID, "ab13ao13") {
			ab13ao13T.Fatal("fixture id")
		}
		ab13ao13CheckFailure(ab13ao13T, ab13ao13Call(ab13ao13T, ab13ao13Case.Ab13Ao13Input), ab13ao13Case.Ab13Ao13Failure, "Hypothesis", true)
	}
}

func TestAb13Ao13RetainedPrior(ab13ao13T *ab13ao13testing.T) {
	for _, ab13ao13Prior := range []int{0, 1, 2, 255} {
		for _, ab13ao13Skip := range []string{"parsed", "count", "registration", "index", "selector", "short"} {
			ab13ao13Input := ab13ao13Base()
			ab13ao13Input["fields"].(map[string]any)["prior_failure"] = []int{ab13ao13Prior}
			switch ab13ao13Skip {
			case "parsed":
				ab13ao13Input["arguments"].(map[string]any)["parsed"] = false
			case "count":
				ab13ao13Input["arguments"].(map[string]any)["count"] = 0
			case "registration":
				ab13ao13Input["fields"].(map[string]any)["registration"] = []int{2}
			case "index":
				ab13ao13Input["arguments"].(map[string]any)["index"] = 1
			case "selector":
				ab13ao13Input["arguments"].(map[string]any)["selector"] = 3
			case "short":
				ab13ao13Input["fields"].(map[string]any)["local_body"] = []int{}
			}
			ab13ao13State := "Hypothesis"
			if ab13ao13Skip == "short" {
				ab13ao13State = "partial"
			}
			ab13ao13CheckFailure(ab13ao13T, ab13ao13Call(ab13ao13T, ab13ao13Input), ab13ao13Prior != 0, ab13ao13State, false)
		}
	}
	for ab13ao13Index := 0; ab13ao13Index < 9; ab13ao13Index++ {
		ab13ao13Input := ab13ao13Base()
		ab13ao13Input["fields"].(map[string]any)["local_body"].([]int)[ab13ao13Index] = 255
		ab13ao13CheckFailure(ab13ao13T, ab13ao13Call(ab13ao13T, ab13ao13Input), ab13ao13Index == 8, "Hypothesis", true)
	}
	for _, ab13ao13View := range []string{"initialize", "timer-consume"} {
		for _, ab13ao13Operation := range []int{0, 1} {
			ab13ao13Input := ab13ao13Base()
			ab13ao13Input["context"].(map[string]any)["view"] = ab13ao13View
			ab13ao13Input["context"].(map[string]any)["operation"] = ab13ao13Operation
			ab13ao13Output := ab13ao13Call(ab13ao13T, ab13ao13Input)
			if ab13ao13Output["state"] != "Unknown" {
				ab13ao13T.Fatal(ab13ao13Output)
			}
		}
	}
}
