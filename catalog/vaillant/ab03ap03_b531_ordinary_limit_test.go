package vaillant

import (
	ab03ap03bytes "bytes"
	ab03ap03json "encoding/json"
	ab03ap03os "os"
	ab03ap03reflect "reflect"
	ab03ap03testing "testing"
)

type ab03ap03Case struct {
	Ab03Ap03ID        string         `json:"id"`
	Ab03Ap03Input     map[string]any `json:"input"`
	Ab03Ap03Token     float64        `json:"token_sum"`
	Ab03Ap03Maximum   float64        `json:"maximum_sum"`
	Ab03Ap03Effective any            `json:"effective"`
	Ab03Ap03Requested any            `json:"requested_state"`
	Ab03Ap03State     string         `json:"state"`
	Ab03Ap03Indices   []any          `json:"admitted_slots"`
}

func ab03ap03Cases(ab03ap03T *ab03ap03testing.T) []ab03ap03Case {
	ab03ap03T.Helper()
	ab03ap03Data, ab03ap03Error := ab03ap03os.ReadFile("testdata/ab03ap03_b531_ordinary_limit.json")
	if ab03ap03Error != nil {
		ab03ap03T.Fatal(ab03ap03Error)
	}
	var ab03ap03Cases []ab03ap03Case
	if ab03ap03Error := ab03ap03json.Unmarshal(ab03ap03Data, &ab03ap03Cases); ab03ap03Error != nil {
		ab03ap03T.Fatal(ab03ap03Error)
	}
	if len(ab03ap03Cases) != 16 {
		ab03ap03T.Fatalf("fixture count: %d", len(ab03ap03Cases))
	}
	return ab03ap03Cases
}

func ab03ap03Run(ab03ap03T *ab03ap03testing.T, ab03ap03Input map[string]any) map[string]any {
	ab03ap03T.Helper()
	ab03ap03Data, ab03ap03Error := ab03ap03json.Marshal(ab03ap03Input)
	if ab03ap03Error != nil {
		ab03ap03T.Fatal(ab03ap03Error)
	}
	ab03ap03Result, ab03ap03Error := Ab03Ap03Evaluate(ab03ap03Data)
	if ab03ap03Error != nil {
		ab03ap03T.Fatal(ab03ap03Error)
	}
	var ab03ap03Output map[string]any
	if ab03ap03Error := ab03ap03json.Unmarshal(ab03ap03Result, &ab03ap03Output); ab03ap03Error != nil {
		ab03ap03T.Fatal(ab03ap03Error)
	}
	for _, ab03ap03Key := range []string{"context", "arguments"} {
		if !ab03ap03reflect.DeepEqual(ab03ap03Input[ab03ap03Key], ab03ap03Output[ab03ap03Key]) {
			ab03ap03T.Fatalf("lost %s", ab03ap03Key)
		}
	}
	if !ab03ap03reflect.DeepEqual(ab03ap03Input["fields"], ab03ap03Output["raw"]) {
		ab03ap03T.Fatal("raw changed")
	}
	if ab03ap03Output["qualified"] != false || ab03ap03Output["physical"] != "Unknown" {
		ab03ap03T.Fatal("qualification")
	}
	return ab03ap03Output
}

func TestAb03Ap03Synthetic(ab03ap03T *ab03ap03testing.T) {
	for _, ab03ap03Case := range ab03ap03Cases(ab03ap03T) {
		ab03ap03T.Run(ab03ap03Case.Ab03Ap03ID, func(ab03ap03T *ab03ap03testing.T) {
			ab03ap03Output := ab03ap03Run(ab03ap03T, ab03ap03Case.Ab03Ap03Input)
			if ab03ap03Output["state"] != ab03ap03Case.Ab03Ap03State {
				ab03ap03T.Fatal(ab03ap03Output)
			}
			ab03ap03Fields := ab03ap03Output["fields"].(map[string]any)
			for ab03ap03Key, ab03ap03Expected := range map[string]float64{"token_sum": ab03ap03Case.Ab03Ap03Token, "maximum_sum": ab03ap03Case.Ab03Ap03Maximum} {
				ab03ap03Field := ab03ap03Fields[ab03ap03Key].(map[string]any)
				if ab03ap03Field["number"] != ab03ap03Expected || ab03ap03Field["state"] != "Hypothesis" || ab03ap03Field["calculation"] != "Proven" {
					ab03ap03T.Fatalf("%s: %v", ab03ap03Key, ab03ap03Field)
				}
			}
			ab03ap03Effective := ab03ap03Fields["effective"].(map[string]any)
			ab03ap03Requested := ab03ap03Fields["requested_state"].(map[string]any)
			if ab03ap03Effective["value"] != ab03ap03Case.Ab03Ap03Effective || ab03ap03Requested["number"] != ab03ap03Case.Ab03Ap03Requested {
				ab03ap03T.Fatal(ab03ap03Fields)
			}
			for _, ab03ap03Pair := range []struct {
				ab03ap03Field map[string]any
				ab03ap03Value any
			}{{ab03ap03Effective, ab03ap03Case.Ab03Ap03Effective}, {ab03ap03Requested, ab03ap03Case.Ab03Ap03Requested}} {
				ab03ap03Grade := "Hypothesis"
				if ab03ap03Pair.ab03ap03Value == nil {
					ab03ap03Grade = "Unknown"
				}
				if ab03ap03Pair.ab03ap03Field["state"] != ab03ap03Grade {
					ab03ap03T.Fatal("grade", ab03ap03Pair)
				}
			}
			if !ab03ap03reflect.DeepEqual(ab03ap03Fields["admitted_slots"].(map[string]any)["indices"], ab03ap03Case.Ab03Ap03Indices) {
				ab03ap03T.Fatal("indices")
			}
			ab03ap03Slots := ab03ap03Fields["slots"].([]any)
			for ab03ap03Index, ab03ap03Row := range ab03ap03Case.Ab03Ap03Input["arguments"].(map[string]any)["rows"].([]any) {
				if !ab03ap03reflect.DeepEqual(ab03ap03Slots[ab03ap03Index].(map[string]any)["raw"], ab03ap03Row) {
					ab03ap03T.Fatal("slot raw")
				}
			}
		})
	}
}

func TestAb03Ap03Context(ab03ap03T *ab03ap03testing.T) {
	for ab03ap03Key, ab03ap03Values := range map[string][]any{
		"primary": {float64(180), nil, "absent", float64(0)}, "secondary": {float64(50), nil, "absent", float64(0)}, "operation": {float64(4), nil, "absent", float64(0)}, "direction": {"remote", nil, "absent", ""}, "view": {"other", nil, "absent", ""},
		"family": {"other-family", nil, "absent", false}, "variant": {"other-variant", nil, "absent", false}, "opaque": {nil, "absent", float64(0), false},
	} {
		for _, ab03ap03Value := range ab03ap03Values {
			ab03ap03Input := ab03ap03Cases(ab03ap03T)[0].Ab03Ap03Input
			ab03ap03Context := ab03ap03Input["context"].(map[string]any)
			if ab03ap03Value == "absent" {
				delete(ab03ap03Context, ab03ap03Key)
			} else {
				ab03ap03Context[ab03ap03Key] = ab03ap03Value
			}
			ab03ap03Output := ab03ap03Run(ab03ap03T, ab03ap03Input)
			ab03ap03Expected := "Unknown"
			if ab03ap03Key == "family" || ab03ap03Key == "variant" || ab03ap03Key == "opaque" {
				ab03ap03Expected = "complete"
			}
			if ab03ap03Output["state"] != ab03ap03Expected {
				ab03ap03T.Fatal(ab03ap03Key, ab03ap03Value, ab03ap03Output)
			}
			if ab03ap03Expected == "Unknown" && len(ab03ap03Output["fields"].(map[string]any)) != 0 {
				ab03ap03T.Fatal("unsupported formula")
			}
		}
	}
}

func TestAb03Ap03IncompleteGate(ab03ap03T *ab03ap03testing.T) {
	for _, ab03ap03Gate := range []any{"absent", []any{}, []any{float64(0), float64(2)}} {
		ab03ap03Input := ab03ap03Cases(ab03ap03T)[0].Ab03Ap03Input
		ab03ap03Fields := ab03ap03Input["fields"].(map[string]any)
		if ab03ap03Gate == "absent" {
			delete(ab03ap03Fields, "gate")
		} else {
			ab03ap03Fields["gate"] = ab03ap03Gate
		}
		ab03ap03Fields["unmapped"] = []any{float64(255)}
		ab03ap03Output := ab03ap03Run(ab03ap03T, ab03ap03Input)
		ab03ap03Semantic := ab03ap03Output["fields"].(map[string]any)
		if ab03ap03Output["state"] != "partial" || ab03ap03Semantic["effective"].(map[string]any)["value"] != true || ab03ap03Semantic["token_sum"].(map[string]any)["number"] != float64(3) {
			ab03ap03T.Fatal(ab03ap03Output)
		}
		if ab03ap03Semantic["requested_state"].(map[string]any)["state"] != "Unknown" {
			ab03ap03T.Fatal("invented request")
		}
	}
}

func TestAb03Ap03Reject(ab03ap03T *ab03ap03testing.T) {
	for _, ab03ap03Data := range []string{``, `{`, `[]`, `null`, `{} {}`, `{"x":1,"x":2}`, `{"context":{"x":1,"x":2}}`, `{"x":[{"y":0,"y":1}]}`} {
		if ab03ap03Result, ab03ap03Error := Ab03Ap03Evaluate([]byte(ab03ap03Data)); ab03ap03Error == nil || ab03ap03Result != nil {
			ab03ap03T.Fatal("accepted", ab03ap03Data)
		}
	}
	for _, ab03ap03Path := range []string{"context", "fields", "arguments", "gate", "limited", "rows", "row", "registered", "capable", "failure", "token", "maximum", "primary", "secondary", "operation", "direction", "view"} {
		for _, ab03ap03Bad := range []any{nil, true, "1", float64(-1), float64(1.5), float64(4294967296), []any{}, map[string]any{}} {
			if ab03ap03Path == "context" || ab03ap03Path == "fields" {
				if _, ab03ap03OK := ab03ap03Bad.(map[string]any); ab03ap03OK {
					continue
				}
			}
			if (ab03ap03Path == "primary" || ab03ap03Path == "secondary" || ab03ap03Path == "operation" || ab03ap03Path == "direction" || ab03ap03Path == "view") && ab03ap03Bad == nil {
				continue
			}
			if ab03ap03Path == "direction" || ab03ap03Path == "view" {
				if _, ab03ap03OK := ab03ap03Bad.(string); ab03ap03OK {
					continue
				}
			}
			if (ab03ap03Path == "registered" || ab03ap03Path == "capable" || ab03ap03Path == "failure") && ab03ap03Bad == true {
				continue
			}
			ab03ap03Input := ab03ap03Cases(ab03ap03T)[0].Ab03Ap03Input
			switch ab03ap03Path {
			case "context", "fields", "arguments":
				ab03ap03Input[ab03ap03Path] = ab03ap03Bad
			case "gate":
				ab03ap03Input["fields"].(map[string]any)["gate"] = []any{ab03ap03Bad}
			case "limited", "rows":
				ab03ap03Input["arguments"].(map[string]any)[ab03ap03Path] = ab03ap03Bad
			case "row":
				ab03ap03Input["arguments"].(map[string]any)["rows"].([]any)[0] = ab03ap03Bad
			case "primary", "secondary", "operation", "direction", "view":
				ab03ap03Input["context"].(map[string]any)[ab03ap03Path] = ab03ap03Bad
			default:
				ab03ap03Input["arguments"].(map[string]any)["rows"].([]any)[0].(map[string]any)[ab03ap03Path] = ab03ap03Bad
			}
			ab03ap03Data, _ := ab03ap03json.Marshal(ab03ap03Input)
			if ab03ap03Result, ab03ap03Error := Ab03Ap03Evaluate(ab03ap03Data); ab03ap03Error == nil || ab03ap03Result != nil {
				ab03ap03T.Fatalf("accepted %s=%v", ab03ap03Path, ab03ap03Bad)
			}
		}
	}
	for _, ab03ap03Key := range []string{"limited", "rows", "registered", "capable", "failure", "token", "maximum"} {
		ab03ap03Input := ab03ap03Cases(ab03ap03T)[0].Ab03Ap03Input
		ab03ap03Parent := ab03ap03Input["arguments"].(map[string]any)
		if ab03ap03Key != "limited" && ab03ap03Key != "rows" {
			ab03ap03Parent = ab03ap03Parent["rows"].([]any)[0].(map[string]any)
		}
		delete(ab03ap03Parent, ab03ap03Key)
		ab03ap03Data, _ := ab03ap03json.Marshal(ab03ap03Input)
		if ab03ap03Result, ab03ap03Error := Ab03Ap03Evaluate(ab03ap03Data); ab03ap03Error == nil || ab03ap03Result != nil {
			ab03ap03T.Fatal("missing", ab03ap03Key)
		}
	}
}

func TestAb03Ap03Isolation(ab03ap03T *ab03ap03testing.T) {
	ab03ap03Input := ab03ap03Cases(ab03ap03T)[0].Ab03Ap03Input
	ab03ap03Data, _ := ab03ap03json.Marshal(ab03ap03Input)
	ab03ap03Original := append([]byte(nil), ab03ap03Data...)
	ab03ap03First, ab03ap03Error := Ab03Ap03Evaluate(ab03ap03Data)
	if ab03ap03Error != nil {
		ab03ap03T.Fatal(ab03ap03Error)
	}
	ab03ap03Saved := append([]byte(nil), ab03ap03First...)
	for ab03ap03Index := range ab03ap03Data {
		ab03ap03Data[ab03ap03Index] = ' '
	}
	if !ab03ap03reflect.DeepEqual([]byte(ab03ap03First), ab03ap03Saved) {
		ab03ap03T.Fatal("input alias")
	}
	ab03ap03Second, _ := Ab03Ap03Evaluate(ab03ap03Original)
	for ab03ap03Index := range ab03ap03First {
		ab03ap03First[ab03ap03Index] = ' '
	}
	if !ab03ap03reflect.DeepEqual([]byte(ab03ap03Second), ab03ap03Saved) {
		ab03ap03T.Fatal("result alias")
	}
	ab03ap03Output := ab03ap03Run(ab03ap03T, ab03ap03Input)
	ab03ap03Output["raw"].(map[string]any)["gate"].([]any)[0] = float64(99)
	ab03ap03Output["arguments"].(map[string]any)["rows"].([]any)[0].(map[string]any)["token"] = float64(99)
	if ab03ap03Output["fields"].(map[string]any)["requested_state"].(map[string]any)["raw"].([]any)[0] != float64(2) {
		ab03ap03T.Fatal("nested raw alias")
	}
	if ab03ap03Output["fields"].(map[string]any)["slots"].([]any)[0].(map[string]any)["raw"].(map[string]any)["token"] != float64(3) {
		ab03ap03T.Fatal("nested argument alias")
	}
	ab03ap03Third, _ := Ab03Ap03Evaluate(ab03ap03Original)
	if !ab03ap03reflect.DeepEqual([]byte(ab03ap03Third), ab03ap03Saved) {
		ab03ap03T.Fatal("later call changed")
	}
}

func TestAb03Ap03ArgumentBoundaries(ab03ap03T *ab03ap03testing.T) {
	for ab03ap03Index := 0; ab03ap03Index < 9; ab03ap03Index++ {
		ab03ap03Input := ab03ap03Cases(ab03ap03T)[0].Ab03Ap03Input
		ab03ap03Arguments := ab03ap03Input["arguments"].(map[string]any)
		ab03ap03Rows := ab03ap03Arguments["rows"].([]any)
		ab03ap03Rows[0].(map[string]any)["registered"] = false
		ab03ap03Rows[ab03ap03Index].(map[string]any)["registered"] = true
		ab03ap03Rows[ab03ap03Index].(map[string]any)["token"] = float64(7)
		ab03ap03Rows[ab03ap03Index].(map[string]any)["maximum"] = float64(11)
		ab03ap03Rows[ab03ap03Index].(map[string]any)["opaque"] = map[string]any{"value": false}
		ab03ap03Arguments["limited"] = float64(4294967295)
		ab03ap03Output := ab03ap03Run(ab03ap03T, ab03ap03Input)
		ab03ap03Fields := ab03ap03Output["fields"].(map[string]any)
		ab03ap03Expected := float64(7)
		if ab03ap03Index == 1 {
			ab03ap03Expected = 0
		}
		if ab03ap03Fields["token_sum"].(map[string]any)["number"] != ab03ap03Expected {
			ab03ap03T.Fatal("slot", ab03ap03Index)
		}
		if !ab03ap03reflect.DeepEqual(ab03ap03Fields["admitted_slots"].(map[string]any)["indices"], []any{float64(ab03ap03Index)}) {
			ab03ap03T.Fatal("slot index", ab03ap03Index)
		}
	}
	for _, ab03ap03Length := range []int{0, 8, 10} {
		ab03ap03Input := ab03ap03Cases(ab03ap03T)[0].Ab03Ap03Input
		ab03ap03Arguments := ab03ap03Input["arguments"].(map[string]any)
		ab03ap03Rows := ab03ap03Arguments["rows"].([]any)
		if ab03ap03Length == 10 {
			ab03ap03Rows = append(ab03ap03Rows, ab03ap03Rows[0])
		} else {
			ab03ap03Rows = ab03ap03Rows[:ab03ap03Length]
		}
		ab03ap03Arguments["rows"] = ab03ap03Rows
		ab03ap03Data, _ := ab03ap03json.Marshal(ab03ap03Input)
		if ab03ap03Result, ab03ap03Error := Ab03Ap03Evaluate(ab03ap03Data); ab03ap03Error == nil || ab03ap03Result != nil {
			ab03ap03T.Fatal("rows length", ab03ap03Length)
		}
	}
	for _, ab03ap03Bad := range []any{nil, false, float64(0), "", map[string]any{}} {
		ab03ap03Input := ab03ap03Cases(ab03ap03T)[0].Ab03Ap03Input
		ab03ap03Input["fields"].(map[string]any)["gate"] = ab03ap03Bad
		ab03ap03Data, _ := ab03ap03json.Marshal(ab03ap03Input)
		if ab03ap03Result, ab03ap03Error := Ab03Ap03Evaluate(ab03ap03Data); ab03ap03Error == nil || ab03ap03Result != nil {
			ab03ap03T.Fatal("gate kind", ab03ap03Bad)
		}
	}
	for _, ab03ap03Key := range []string{"context", "fields", "arguments"} {
		ab03ap03Input := ab03ap03Cases(ab03ap03T)[0].Ab03Ap03Input
		delete(ab03ap03Input, ab03ap03Key)
		ab03ap03Data, _ := ab03ap03json.Marshal(ab03ap03Input)
		if ab03ap03Result, ab03ap03Error := Ab03Ap03Evaluate(ab03ap03Data); ab03ap03Error == nil || ab03ap03Result != nil {
			ab03ap03T.Fatal("missing", ab03ap03Key)
		}
	}
	for _, ab03ap03Extra := range []any{nil, false, float64(0), []any{}, map[string]any{}} {
		ab03ap03Input := ab03ap03Cases(ab03ap03T)[0].Ab03Ap03Input
		ab03ap03Input["arguments"].(map[string]any)["opaque"] = ab03ap03Extra
		ab03ap03Run(ab03ap03T, ab03ap03Input)
	}
}

func TestAb03Ap03RawOnlyAndDuplicates(ab03ap03T *ab03ap03testing.T) {
	for _, ab03ap03Arguments := range []any{"absent", map[string]any{}} {
		ab03ap03Input := ab03ap03Cases(ab03ap03T)[0].Ab03Ap03Input
		ab03ap03Input["context"].(map[string]any)["view"] = "unmapped"
		if ab03ap03Arguments == "absent" {
			delete(ab03ap03Input, "arguments")
		} else {
			ab03ap03Input["arguments"] = ab03ap03Arguments
		}
		ab03ap03Output := ab03ap03Run(ab03ap03T, ab03ap03Input)
		if _, ab03ap03Present := ab03ap03Output["arguments"]; ab03ap03Present != (ab03ap03Arguments != "absent") {
			ab03ap03T.Fatal("arguments presence")
		}
	}
	ab03ap03Input := ab03ap03Cases(ab03ap03T)[0].Ab03Ap03Input
	ab03ap03Data, _ := ab03ap03json.Marshal(ab03ap03Input)
	for _, ab03ap03Pair := range [][2]string{{`"primary":181`, `"primary":181,"primary":181`}, {`"gate":[2]`, `"gate":[2],"gate":[2]`}, {`"limited":1`, `"limited":1,"limited":1`}, {`"token":3`, `"token":3,"token":3`}, {`"x":false`, `"x":false,"x":false`}} {
		ab03ap03Duplicate := ab03ap03bytes.Replace(ab03ap03Data, []byte(ab03ap03Pair[0]), []byte(ab03ap03Pair[1]), 1)
		if ab03ap03bytes.Equal(ab03ap03Data, ab03ap03Duplicate) {
			ab03ap03T.Fatal("duplicate control missing")
		}
		if ab03ap03Result, ab03ap03Error := Ab03Ap03Evaluate(ab03ap03Duplicate); ab03ap03Error == nil || ab03ap03Result != nil {
			ab03ap03T.Fatal("duplicate accepted", ab03ap03Pair[0])
		}
	}
}

func TestAb03Ap03UnsupportedValidation(ab03ap03T *ab03ap03testing.T) {
	for _, ab03ap03Key := range []string{"limited", "rows", "token", "maximum", "registered", "capable", "failure", "gate"} {
		ab03ap03Input := ab03ap03Cases(ab03ap03T)[0].Ab03Ap03Input
		ab03ap03Input["context"].(map[string]any)["primary"] = float64(0)
		switch ab03ap03Key {
		case "limited", "rows":
			ab03ap03Input["arguments"].(map[string]any)[ab03ap03Key] = nil
		case "gate":
			ab03ap03Input["fields"].(map[string]any)["gate"] = []any{float64(256)}
		default:
			ab03ap03Input["arguments"].(map[string]any)["rows"].([]any)[0].(map[string]any)[ab03ap03Key] = nil
		}
		ab03ap03Data, _ := ab03ap03json.Marshal(ab03ap03Input)
		if ab03ap03Result, ab03ap03Error := Ab03Ap03Evaluate(ab03ap03Data); ab03ap03Error == nil || ab03ap03Result != nil {
			ab03ap03T.Fatal("unsupported invalid argument", ab03ap03Key)
		}
	}
}
