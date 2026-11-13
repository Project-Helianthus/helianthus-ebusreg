package vaillant

import (
	b531OrdinaryLimitBytesPackage "bytes"
	b531OrdinaryLimitJSON "encoding/json"
	b531OrdinaryLimitOs "os"
	b531OrdinaryLimitReflect "reflect"
	b531OrdinaryLimitTesting "testing"
)

type b531OrdinaryLimitCase struct {
	B531OrdinaryLimitID        string         `json:"id"`
	B531OrdinaryLimitInput     map[string]any `json:"input"`
	B531OrdinaryLimitToken     float64        `json:"token_sum"`
	B531OrdinaryLimitMaximum   float64        `json:"maximum_sum"`
	B531OrdinaryLimitEffective any            `json:"effective"`
	B531OrdinaryLimitRequested any            `json:"requested_state"`
	B531OrdinaryLimitState     string         `json:"state"`
	B531OrdinaryLimitIndices   []any          `json:"admitted_slots"`
}

func b531OrdinaryLimitCases(b531OrdinaryLimitT *b531OrdinaryLimitTesting.T) []b531OrdinaryLimitCase {
	b531OrdinaryLimitT.Helper()
	b531OrdinaryLimitData, b531OrdinaryLimitError := b531OrdinaryLimitOs.ReadFile("testdata/ab03ap03_b531_ordinary_limit.json")
	if b531OrdinaryLimitError != nil {
		b531OrdinaryLimitT.Fatal(b531OrdinaryLimitError)
	}
	var b531OrdinaryLimitCases []b531OrdinaryLimitCase
	if b531OrdinaryLimitError := b531OrdinaryLimitJSON.Unmarshal(b531OrdinaryLimitData, &b531OrdinaryLimitCases); b531OrdinaryLimitError != nil {
		b531OrdinaryLimitT.Fatal(b531OrdinaryLimitError)
	}
	if len(b531OrdinaryLimitCases) != 16 {
		b531OrdinaryLimitT.Fatalf("fixture count: %d", len(b531OrdinaryLimitCases))
	}
	return b531OrdinaryLimitCases
}

func b531OrdinaryLimitRun(b531OrdinaryLimitT *b531OrdinaryLimitTesting.T, b531OrdinaryLimitInput map[string]any) map[string]any {
	b531OrdinaryLimitT.Helper()
	b531OrdinaryLimitData, b531OrdinaryLimitError := b531OrdinaryLimitJSON.Marshal(b531OrdinaryLimitInput)
	if b531OrdinaryLimitError != nil {
		b531OrdinaryLimitT.Fatal(b531OrdinaryLimitError)
	}
	b531OrdinaryLimitResult, b531OrdinaryLimitError := B531OrdinaryLimitEvaluate(b531OrdinaryLimitData)
	if b531OrdinaryLimitError != nil {
		b531OrdinaryLimitT.Fatal(b531OrdinaryLimitError)
	}
	var b531OrdinaryLimitOutput map[string]any
	if b531OrdinaryLimitError := b531OrdinaryLimitJSON.Unmarshal(b531OrdinaryLimitResult, &b531OrdinaryLimitOutput); b531OrdinaryLimitError != nil {
		b531OrdinaryLimitT.Fatal(b531OrdinaryLimitError)
	}
	for _, b531OrdinaryLimitKey := range []string{"context", "arguments"} {
		if !b531OrdinaryLimitReflect.DeepEqual(b531OrdinaryLimitInput[b531OrdinaryLimitKey], b531OrdinaryLimitOutput[b531OrdinaryLimitKey]) {
			b531OrdinaryLimitT.Fatalf("lost %s", b531OrdinaryLimitKey)
		}
	}
	if !b531OrdinaryLimitReflect.DeepEqual(b531OrdinaryLimitInput["fields"], b531OrdinaryLimitOutput["raw"]) {
		b531OrdinaryLimitT.Fatal("raw changed")
	}
	if b531OrdinaryLimitOutput["qualified"] != false || b531OrdinaryLimitOutput["physical"] != "Unknown" {
		b531OrdinaryLimitT.Fatal("qualification")
	}
	return b531OrdinaryLimitOutput
}

func TestAb03Ap03Synthetic(b531OrdinaryLimitT *b531OrdinaryLimitTesting.T) {
	for _, b531OrdinaryLimitCase := range b531OrdinaryLimitCases(b531OrdinaryLimitT) {
		b531OrdinaryLimitT.Run(b531OrdinaryLimitCase.B531OrdinaryLimitID, func(b531OrdinaryLimitT *b531OrdinaryLimitTesting.T) {
			b531OrdinaryLimitOutput := b531OrdinaryLimitRun(b531OrdinaryLimitT, b531OrdinaryLimitCase.B531OrdinaryLimitInput)
			if b531OrdinaryLimitOutput["state"] != b531OrdinaryLimitCase.B531OrdinaryLimitState {
				b531OrdinaryLimitT.Fatal(b531OrdinaryLimitOutput)
			}
			b531OrdinaryLimitFields := b531OrdinaryLimitOutput["fields"].(map[string]any)
			for b531OrdinaryLimitKey, b531OrdinaryLimitExpected := range map[string]float64{"token_sum": b531OrdinaryLimitCase.B531OrdinaryLimitToken, "maximum_sum": b531OrdinaryLimitCase.B531OrdinaryLimitMaximum} {
				b531OrdinaryLimitField := b531OrdinaryLimitFields[b531OrdinaryLimitKey].(map[string]any)
				if b531OrdinaryLimitField["number"] != b531OrdinaryLimitExpected || b531OrdinaryLimitField["state"] != "Hypothesis" || b531OrdinaryLimitField["calculation"] != "Proven" {
					b531OrdinaryLimitT.Fatalf("%s: %v", b531OrdinaryLimitKey, b531OrdinaryLimitField)
				}
			}
			b531OrdinaryLimitEffective := b531OrdinaryLimitFields["effective"].(map[string]any)
			b531OrdinaryLimitRequested := b531OrdinaryLimitFields["requested_state"].(map[string]any)
			if b531OrdinaryLimitEffective["value"] != b531OrdinaryLimitCase.B531OrdinaryLimitEffective || b531OrdinaryLimitRequested["number"] != b531OrdinaryLimitCase.B531OrdinaryLimitRequested {
				b531OrdinaryLimitT.Fatal(b531OrdinaryLimitFields)
			}
			for _, b531OrdinaryLimitPair := range []struct {
				b531OrdinaryLimitField map[string]any
				b531OrdinaryLimitValue any
			}{{b531OrdinaryLimitEffective, b531OrdinaryLimitCase.B531OrdinaryLimitEffective}, {b531OrdinaryLimitRequested, b531OrdinaryLimitCase.B531OrdinaryLimitRequested}} {
				b531OrdinaryLimitGrade := "Hypothesis"
				if b531OrdinaryLimitPair.b531OrdinaryLimitValue == nil {
					b531OrdinaryLimitGrade = "Unknown"
				}
				if b531OrdinaryLimitPair.b531OrdinaryLimitField["state"] != b531OrdinaryLimitGrade {
					b531OrdinaryLimitT.Fatal("grade", b531OrdinaryLimitPair)
				}
			}
			if !b531OrdinaryLimitReflect.DeepEqual(b531OrdinaryLimitFields["admitted_slots"].(map[string]any)["indices"], b531OrdinaryLimitCase.B531OrdinaryLimitIndices) {
				b531OrdinaryLimitT.Fatal("indices")
			}
			b531OrdinaryLimitSlots := b531OrdinaryLimitFields["slots"].([]any)
			for b531OrdinaryLimitIndex, b531OrdinaryLimitRow := range b531OrdinaryLimitCase.B531OrdinaryLimitInput["arguments"].(map[string]any)["rows"].([]any) {
				if !b531OrdinaryLimitReflect.DeepEqual(b531OrdinaryLimitSlots[b531OrdinaryLimitIndex].(map[string]any)["raw"], b531OrdinaryLimitRow) {
					b531OrdinaryLimitT.Fatal("slot raw")
				}
			}
		})
	}
}

func TestAb03Ap03Context(b531OrdinaryLimitT *b531OrdinaryLimitTesting.T) {
	for b531OrdinaryLimitKey, b531OrdinaryLimitValues := range map[string][]any{
		"primary": {float64(180), nil, "absent", float64(0)}, "secondary": {float64(50), nil, "absent", float64(0)}, "operation": {float64(4), nil, "absent", float64(0)}, "direction": {"remote", nil, "absent", ""}, "view": {"other", nil, "absent", ""},
		"family": {"other-family", nil, "absent", false}, "variant": {"other-variant", nil, "absent", false}, "opaque": {nil, "absent", float64(0), false},
	} {
		for _, b531OrdinaryLimitValue := range b531OrdinaryLimitValues {
			b531OrdinaryLimitInput := b531OrdinaryLimitCases(b531OrdinaryLimitT)[0].B531OrdinaryLimitInput
			b531OrdinaryLimitContext := b531OrdinaryLimitInput["context"].(map[string]any)
			if b531OrdinaryLimitValue == "absent" {
				delete(b531OrdinaryLimitContext, b531OrdinaryLimitKey)
			} else {
				b531OrdinaryLimitContext[b531OrdinaryLimitKey] = b531OrdinaryLimitValue
			}
			b531OrdinaryLimitOutput := b531OrdinaryLimitRun(b531OrdinaryLimitT, b531OrdinaryLimitInput)
			b531OrdinaryLimitExpected := "Unknown"
			if b531OrdinaryLimitKey == "family" || b531OrdinaryLimitKey == "variant" || b531OrdinaryLimitKey == "opaque" {
				b531OrdinaryLimitExpected = "complete"
			}
			if b531OrdinaryLimitOutput["state"] != b531OrdinaryLimitExpected {
				b531OrdinaryLimitT.Fatal(b531OrdinaryLimitKey, b531OrdinaryLimitValue, b531OrdinaryLimitOutput)
			}
			if b531OrdinaryLimitExpected == "Unknown" && len(b531OrdinaryLimitOutput["fields"].(map[string]any)) != 0 {
				b531OrdinaryLimitT.Fatal("unsupported formula")
			}
		}
	}
}

func TestAb03Ap03IncompleteGate(b531OrdinaryLimitT *b531OrdinaryLimitTesting.T) {
	for _, b531OrdinaryLimitGate := range []any{"absent", []any{}, []any{float64(0), float64(2)}} {
		b531OrdinaryLimitInput := b531OrdinaryLimitCases(b531OrdinaryLimitT)[0].B531OrdinaryLimitInput
		b531OrdinaryLimitFields := b531OrdinaryLimitInput["fields"].(map[string]any)
		if b531OrdinaryLimitGate == "absent" {
			delete(b531OrdinaryLimitFields, "gate")
		} else {
			b531OrdinaryLimitFields["gate"] = b531OrdinaryLimitGate
		}
		b531OrdinaryLimitFields["unmapped"] = []any{float64(255)}
		b531OrdinaryLimitOutput := b531OrdinaryLimitRun(b531OrdinaryLimitT, b531OrdinaryLimitInput)
		b531OrdinaryLimitSemantic := b531OrdinaryLimitOutput["fields"].(map[string]any)
		if b531OrdinaryLimitOutput["state"] != "partial" || b531OrdinaryLimitSemantic["effective"].(map[string]any)["value"] != true || b531OrdinaryLimitSemantic["token_sum"].(map[string]any)["number"] != float64(3) {
			b531OrdinaryLimitT.Fatal(b531OrdinaryLimitOutput)
		}
		if b531OrdinaryLimitSemantic["requested_state"].(map[string]any)["state"] != "Unknown" {
			b531OrdinaryLimitT.Fatal("invented request")
		}
	}
}

func TestAb03Ap03Reject(b531OrdinaryLimitT *b531OrdinaryLimitTesting.T) {
	for _, b531OrdinaryLimitData := range []string{``, `{`, `[]`, `null`, `{} {}`, `{"x":1,"x":2}`, `{"context":{"x":1,"x":2}}`, `{"x":[{"y":0,"y":1}]}`} {
		if b531OrdinaryLimitResult, b531OrdinaryLimitError := B531OrdinaryLimitEvaluate([]byte(b531OrdinaryLimitData)); b531OrdinaryLimitError == nil || b531OrdinaryLimitResult != nil {
			b531OrdinaryLimitT.Fatal("accepted", b531OrdinaryLimitData)
		}
	}
	for _, b531OrdinaryLimitPath := range []string{"context", "fields", "arguments", "gate", "limited", "rows", "row", "registered", "capable", "failure", "token", "maximum", "primary", "secondary", "operation", "direction", "view"} {
		for _, b531OrdinaryLimitBad := range []any{nil, true, "1", float64(-1), float64(1.5), float64(4294967296), []any{}, map[string]any{}} {
			if b531OrdinaryLimitPath == "context" || b531OrdinaryLimitPath == "fields" {
				if _, b531OrdinaryLimitOK := b531OrdinaryLimitBad.(map[string]any); b531OrdinaryLimitOK {
					continue
				}
			}
			if (b531OrdinaryLimitPath == "primary" || b531OrdinaryLimitPath == "secondary" || b531OrdinaryLimitPath == "operation" || b531OrdinaryLimitPath == "direction" || b531OrdinaryLimitPath == "view") && b531OrdinaryLimitBad == nil {
				continue
			}
			if b531OrdinaryLimitPath == "direction" || b531OrdinaryLimitPath == "view" {
				if _, b531OrdinaryLimitOK := b531OrdinaryLimitBad.(string); b531OrdinaryLimitOK {
					continue
				}
			}
			if (b531OrdinaryLimitPath == "registered" || b531OrdinaryLimitPath == "capable" || b531OrdinaryLimitPath == "failure") && b531OrdinaryLimitBad == true {
				continue
			}
			b531OrdinaryLimitInput := b531OrdinaryLimitCases(b531OrdinaryLimitT)[0].B531OrdinaryLimitInput
			switch b531OrdinaryLimitPath {
			case "context", "fields", "arguments":
				b531OrdinaryLimitInput[b531OrdinaryLimitPath] = b531OrdinaryLimitBad
			case "gate":
				b531OrdinaryLimitInput["fields"].(map[string]any)["gate"] = []any{b531OrdinaryLimitBad}
			case "limited", "rows":
				b531OrdinaryLimitInput["arguments"].(map[string]any)[b531OrdinaryLimitPath] = b531OrdinaryLimitBad
			case "row":
				b531OrdinaryLimitInput["arguments"].(map[string]any)["rows"].([]any)[0] = b531OrdinaryLimitBad
			case "primary", "secondary", "operation", "direction", "view":
				b531OrdinaryLimitInput["context"].(map[string]any)[b531OrdinaryLimitPath] = b531OrdinaryLimitBad
			default:
				b531OrdinaryLimitInput["arguments"].(map[string]any)["rows"].([]any)[0].(map[string]any)[b531OrdinaryLimitPath] = b531OrdinaryLimitBad
			}
			b531OrdinaryLimitData, _ := b531OrdinaryLimitJSON.Marshal(b531OrdinaryLimitInput)
			if b531OrdinaryLimitResult, b531OrdinaryLimitError := B531OrdinaryLimitEvaluate(b531OrdinaryLimitData); b531OrdinaryLimitError == nil || b531OrdinaryLimitResult != nil {
				b531OrdinaryLimitT.Fatalf("accepted %s=%v", b531OrdinaryLimitPath, b531OrdinaryLimitBad)
			}
		}
	}
	for _, b531OrdinaryLimitKey := range []string{"limited", "rows", "registered", "capable", "failure", "token", "maximum"} {
		b531OrdinaryLimitInput := b531OrdinaryLimitCases(b531OrdinaryLimitT)[0].B531OrdinaryLimitInput
		b531OrdinaryLimitParent := b531OrdinaryLimitInput["arguments"].(map[string]any)
		if b531OrdinaryLimitKey != "limited" && b531OrdinaryLimitKey != "rows" {
			b531OrdinaryLimitParent = b531OrdinaryLimitParent["rows"].([]any)[0].(map[string]any)
		}
		delete(b531OrdinaryLimitParent, b531OrdinaryLimitKey)
		b531OrdinaryLimitData, _ := b531OrdinaryLimitJSON.Marshal(b531OrdinaryLimitInput)
		if b531OrdinaryLimitResult, b531OrdinaryLimitError := B531OrdinaryLimitEvaluate(b531OrdinaryLimitData); b531OrdinaryLimitError == nil || b531OrdinaryLimitResult != nil {
			b531OrdinaryLimitT.Fatal("missing", b531OrdinaryLimitKey)
		}
	}
}

func TestAb03Ap03Isolation(b531OrdinaryLimitT *b531OrdinaryLimitTesting.T) {
	b531OrdinaryLimitInput := b531OrdinaryLimitCases(b531OrdinaryLimitT)[0].B531OrdinaryLimitInput
	b531OrdinaryLimitData, _ := b531OrdinaryLimitJSON.Marshal(b531OrdinaryLimitInput)
	b531OrdinaryLimitOriginal := append([]byte(nil), b531OrdinaryLimitData...)
	b531OrdinaryLimitFirst, b531OrdinaryLimitError := B531OrdinaryLimitEvaluate(b531OrdinaryLimitData)
	if b531OrdinaryLimitError != nil {
		b531OrdinaryLimitT.Fatal(b531OrdinaryLimitError)
	}
	b531OrdinaryLimitSaved := append([]byte(nil), b531OrdinaryLimitFirst...)
	for b531OrdinaryLimitIndex := range b531OrdinaryLimitData {
		b531OrdinaryLimitData[b531OrdinaryLimitIndex] = ' '
	}
	if !b531OrdinaryLimitReflect.DeepEqual([]byte(b531OrdinaryLimitFirst), b531OrdinaryLimitSaved) {
		b531OrdinaryLimitT.Fatal("input alias")
	}
	b531OrdinaryLimitSecond, _ := B531OrdinaryLimitEvaluate(b531OrdinaryLimitOriginal)
	for b531OrdinaryLimitIndex := range b531OrdinaryLimitFirst {
		b531OrdinaryLimitFirst[b531OrdinaryLimitIndex] = ' '
	}
	if !b531OrdinaryLimitReflect.DeepEqual([]byte(b531OrdinaryLimitSecond), b531OrdinaryLimitSaved) {
		b531OrdinaryLimitT.Fatal("result alias")
	}
	b531OrdinaryLimitOutput := b531OrdinaryLimitRun(b531OrdinaryLimitT, b531OrdinaryLimitInput)
	b531OrdinaryLimitOutput["raw"].(map[string]any)["gate"].([]any)[0] = float64(99)
	b531OrdinaryLimitOutput["arguments"].(map[string]any)["rows"].([]any)[0].(map[string]any)["token"] = float64(99)
	if b531OrdinaryLimitOutput["fields"].(map[string]any)["requested_state"].(map[string]any)["raw"].([]any)[0] != float64(2) {
		b531OrdinaryLimitT.Fatal("nested raw alias")
	}
	if b531OrdinaryLimitOutput["fields"].(map[string]any)["slots"].([]any)[0].(map[string]any)["raw"].(map[string]any)["token"] != float64(3) {
		b531OrdinaryLimitT.Fatal("nested argument alias")
	}
	b531OrdinaryLimitThird, _ := B531OrdinaryLimitEvaluate(b531OrdinaryLimitOriginal)
	if !b531OrdinaryLimitReflect.DeepEqual([]byte(b531OrdinaryLimitThird), b531OrdinaryLimitSaved) {
		b531OrdinaryLimitT.Fatal("later call changed")
	}
}

func TestAb03Ap03ArgumentBoundaries(b531OrdinaryLimitT *b531OrdinaryLimitTesting.T) {
	for b531OrdinaryLimitIndex := 0; b531OrdinaryLimitIndex < 9; b531OrdinaryLimitIndex++ {
		b531OrdinaryLimitInput := b531OrdinaryLimitCases(b531OrdinaryLimitT)[0].B531OrdinaryLimitInput
		b531OrdinaryLimitArguments := b531OrdinaryLimitInput["arguments"].(map[string]any)
		b531OrdinaryLimitRows := b531OrdinaryLimitArguments["rows"].([]any)
		b531OrdinaryLimitRows[0].(map[string]any)["registered"] = false
		b531OrdinaryLimitRows[b531OrdinaryLimitIndex].(map[string]any)["registered"] = true
		b531OrdinaryLimitRows[b531OrdinaryLimitIndex].(map[string]any)["token"] = float64(7)
		b531OrdinaryLimitRows[b531OrdinaryLimitIndex].(map[string]any)["maximum"] = float64(11)
		b531OrdinaryLimitRows[b531OrdinaryLimitIndex].(map[string]any)["opaque"] = map[string]any{"value": false}
		b531OrdinaryLimitArguments["limited"] = float64(4294967295)
		b531OrdinaryLimitOutput := b531OrdinaryLimitRun(b531OrdinaryLimitT, b531OrdinaryLimitInput)
		b531OrdinaryLimitFields := b531OrdinaryLimitOutput["fields"].(map[string]any)
		b531OrdinaryLimitExpected := float64(7)
		if b531OrdinaryLimitIndex == 1 {
			b531OrdinaryLimitExpected = 0
		}
		if b531OrdinaryLimitFields["token_sum"].(map[string]any)["number"] != b531OrdinaryLimitExpected {
			b531OrdinaryLimitT.Fatal("slot", b531OrdinaryLimitIndex)
		}
		if !b531OrdinaryLimitReflect.DeepEqual(b531OrdinaryLimitFields["admitted_slots"].(map[string]any)["indices"], []any{float64(b531OrdinaryLimitIndex)}) {
			b531OrdinaryLimitT.Fatal("slot index", b531OrdinaryLimitIndex)
		}
	}
	for _, b531OrdinaryLimitLength := range []int{0, 8, 10} {
		b531OrdinaryLimitInput := b531OrdinaryLimitCases(b531OrdinaryLimitT)[0].B531OrdinaryLimitInput
		b531OrdinaryLimitArguments := b531OrdinaryLimitInput["arguments"].(map[string]any)
		b531OrdinaryLimitRows := b531OrdinaryLimitArguments["rows"].([]any)
		if b531OrdinaryLimitLength == 10 {
			b531OrdinaryLimitRows = append(b531OrdinaryLimitRows, b531OrdinaryLimitRows[0])
		} else {
			b531OrdinaryLimitRows = b531OrdinaryLimitRows[:b531OrdinaryLimitLength]
		}
		b531OrdinaryLimitArguments["rows"] = b531OrdinaryLimitRows
		b531OrdinaryLimitData, _ := b531OrdinaryLimitJSON.Marshal(b531OrdinaryLimitInput)
		if b531OrdinaryLimitResult, b531OrdinaryLimitError := B531OrdinaryLimitEvaluate(b531OrdinaryLimitData); b531OrdinaryLimitError == nil || b531OrdinaryLimitResult != nil {
			b531OrdinaryLimitT.Fatal("rows length", b531OrdinaryLimitLength)
		}
	}
	for _, b531OrdinaryLimitBad := range []any{nil, false, float64(0), "", map[string]any{}} {
		b531OrdinaryLimitInput := b531OrdinaryLimitCases(b531OrdinaryLimitT)[0].B531OrdinaryLimitInput
		b531OrdinaryLimitInput["fields"].(map[string]any)["gate"] = b531OrdinaryLimitBad
		b531OrdinaryLimitData, _ := b531OrdinaryLimitJSON.Marshal(b531OrdinaryLimitInput)
		if b531OrdinaryLimitResult, b531OrdinaryLimitError := B531OrdinaryLimitEvaluate(b531OrdinaryLimitData); b531OrdinaryLimitError == nil || b531OrdinaryLimitResult != nil {
			b531OrdinaryLimitT.Fatal("gate kind", b531OrdinaryLimitBad)
		}
	}
	for _, b531OrdinaryLimitKey := range []string{"context", "fields", "arguments"} {
		b531OrdinaryLimitInput := b531OrdinaryLimitCases(b531OrdinaryLimitT)[0].B531OrdinaryLimitInput
		delete(b531OrdinaryLimitInput, b531OrdinaryLimitKey)
		b531OrdinaryLimitData, _ := b531OrdinaryLimitJSON.Marshal(b531OrdinaryLimitInput)
		if b531OrdinaryLimitResult, b531OrdinaryLimitError := B531OrdinaryLimitEvaluate(b531OrdinaryLimitData); b531OrdinaryLimitError == nil || b531OrdinaryLimitResult != nil {
			b531OrdinaryLimitT.Fatal("missing", b531OrdinaryLimitKey)
		}
	}
	for _, b531OrdinaryLimitExtra := range []any{nil, false, float64(0), []any{}, map[string]any{}} {
		b531OrdinaryLimitInput := b531OrdinaryLimitCases(b531OrdinaryLimitT)[0].B531OrdinaryLimitInput
		b531OrdinaryLimitInput["arguments"].(map[string]any)["opaque"] = b531OrdinaryLimitExtra
		b531OrdinaryLimitRun(b531OrdinaryLimitT, b531OrdinaryLimitInput)
	}
}

func TestAb03Ap03RawOnlyAndDuplicates(b531OrdinaryLimitT *b531OrdinaryLimitTesting.T) {
	for _, b531OrdinaryLimitArguments := range []any{"absent", map[string]any{}} {
		b531OrdinaryLimitInput := b531OrdinaryLimitCases(b531OrdinaryLimitT)[0].B531OrdinaryLimitInput
		b531OrdinaryLimitInput["context"].(map[string]any)["view"] = "unmapped"
		if b531OrdinaryLimitArguments == "absent" {
			delete(b531OrdinaryLimitInput, "arguments")
		} else {
			b531OrdinaryLimitInput["arguments"] = b531OrdinaryLimitArguments
		}
		b531OrdinaryLimitOutput := b531OrdinaryLimitRun(b531OrdinaryLimitT, b531OrdinaryLimitInput)
		if _, b531OrdinaryLimitPresent := b531OrdinaryLimitOutput["arguments"]; b531OrdinaryLimitPresent != (b531OrdinaryLimitArguments != "absent") {
			b531OrdinaryLimitT.Fatal("arguments presence")
		}
	}
	b531OrdinaryLimitInput := b531OrdinaryLimitCases(b531OrdinaryLimitT)[0].B531OrdinaryLimitInput
	b531OrdinaryLimitData, _ := b531OrdinaryLimitJSON.Marshal(b531OrdinaryLimitInput)
	for _, b531OrdinaryLimitPair := range [][2]string{{`"primary":181`, `"primary":181,"primary":181`}, {`"gate":[2]`, `"gate":[2],"gate":[2]`}, {`"limited":1`, `"limited":1,"limited":1`}, {`"token":3`, `"token":3,"token":3`}, {`"x":false`, `"x":false,"x":false`}} {
		b531OrdinaryLimitDuplicate := b531OrdinaryLimitBytesPackage.Replace(b531OrdinaryLimitData, []byte(b531OrdinaryLimitPair[0]), []byte(b531OrdinaryLimitPair[1]), 1)
		if b531OrdinaryLimitBytesPackage.Equal(b531OrdinaryLimitData, b531OrdinaryLimitDuplicate) {
			b531OrdinaryLimitT.Fatal("duplicate control missing")
		}
		if b531OrdinaryLimitResult, b531OrdinaryLimitError := B531OrdinaryLimitEvaluate(b531OrdinaryLimitDuplicate); b531OrdinaryLimitError == nil || b531OrdinaryLimitResult != nil {
			b531OrdinaryLimitT.Fatal("duplicate accepted", b531OrdinaryLimitPair[0])
		}
	}
}

func TestAb03Ap03UnsupportedValidation(b531OrdinaryLimitT *b531OrdinaryLimitTesting.T) {
	for _, b531OrdinaryLimitKey := range []string{"limited", "rows", "token", "maximum", "registered", "capable", "failure", "gate"} {
		b531OrdinaryLimitInput := b531OrdinaryLimitCases(b531OrdinaryLimitT)[0].B531OrdinaryLimitInput
		b531OrdinaryLimitInput["context"].(map[string]any)["primary"] = float64(0)
		switch b531OrdinaryLimitKey {
		case "limited", "rows":
			b531OrdinaryLimitInput["arguments"].(map[string]any)[b531OrdinaryLimitKey] = nil
		case "gate":
			b531OrdinaryLimitInput["fields"].(map[string]any)["gate"] = []any{float64(256)}
		default:
			b531OrdinaryLimitInput["arguments"].(map[string]any)["rows"].([]any)[0].(map[string]any)[b531OrdinaryLimitKey] = nil
		}
		b531OrdinaryLimitData, _ := b531OrdinaryLimitJSON.Marshal(b531OrdinaryLimitInput)
		if b531OrdinaryLimitResult, b531OrdinaryLimitError := B531OrdinaryLimitEvaluate(b531OrdinaryLimitData); b531OrdinaryLimitError == nil || b531OrdinaryLimitResult != nil {
			b531OrdinaryLimitT.Fatal("unsupported invalid argument", b531OrdinaryLimitKey)
		}
	}
}
