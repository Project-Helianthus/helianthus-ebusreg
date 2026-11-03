package vaillant

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func ab25ao25Run(ab25ao25T *testing.T, ab25ao25Text string) map[string]any {
	ab25ao25T.Helper()
	ab25ao25Out, ab25ao25Err := Ab25Ao25Interpret(json.RawMessage(ab25ao25Text))
	if ab25ao25Err != nil {
		ab25ao25T.Fatalf("unexpected error: %v", ab25ao25Err)
	}
	var ab25ao25Map map[string]any
	if ab25ao25Err = json.Unmarshal(ab25ao25Out, &ab25ao25Map); ab25ao25Err != nil {
		ab25ao25T.Fatalf("output is not an object: %v", ab25ao25Err)
	}
	return ab25ao25Map
}

func ab25ao25Input(ab25ao25Context string, ab25ao25Fields string) string {
	return `{"context":` + ab25ao25Context + `,"fields":` + ab25ao25Fields + `}`
}

func ab25ao25Context(ab25ao25Direction string) string {
	return `{"primary":181,"secondary":34,"operation":1,"direction":"` + ab25ao25Direction + `","view":"isolated-fields"}`
}

func ab25ao25FieldOf(ab25ao25T *testing.T, ab25ao25Out map[string]any, ab25ao25Name string) map[string]any {
	ab25ao25T.Helper()
	ab25ao25Fields := ab25ao25Out["fields"].(map[string]any)
	ab25ao25Field, ab25ao25OK := ab25ao25Fields[ab25ao25Name].(map[string]any)
	if !ab25ao25OK {
		ab25ao25T.Fatalf("field %q missing in %v", ab25ao25Name, ab25ao25Fields)
	}
	return ab25ao25Field
}

func ab25ao25Request(ab25ao25T *testing.T, ab25ao25Fields string) map[string]any {
	ab25ao25T.Helper()
	return ab25ao25Run(ab25ao25T, ab25ao25Input(ab25ao25Context("request"), ab25ao25Fields))
}

func TestAb25Ao25FixtureCases(ab25ao25T *testing.T) {
	ab25ao25Data, ab25ao25Err := os.ReadFile("testdata/ab25ao25_b522_demand_views.json")
	if ab25ao25Err != nil {
		ab25ao25T.Fatal(ab25ao25Err)
	}
	var ab25ao25File struct {
		Cases []struct {
			ID       string          `json:"id"`
			Input    json.RawMessage `json:"input"`
			Expected json.RawMessage `json:"expected"`
		} `json:"ab25ao25_cases"`
	}
	if ab25ao25Err = json.Unmarshal(ab25ao25Data, &ab25ao25File); ab25ao25Err != nil {
		ab25ao25T.Fatal(ab25ao25Err)
	}
	if len(ab25ao25File.Cases) != 4 {
		ab25ao25T.Fatalf("fixture cases = %d", len(ab25ao25File.Cases))
	}
	for _, ab25ao25Case := range ab25ao25File.Cases {
		if !strings.HasPrefix(ab25ao25Case.ID, "ab25ao25") {
			ab25ao25T.Fatalf("fixture id %q", ab25ao25Case.ID)
		}
		ab25ao25Got := ab25ao25Run(ab25ao25T, string(ab25ao25Case.Input))
		var ab25ao25Want map[string]any
		if ab25ao25Err = json.Unmarshal(ab25ao25Case.Expected, &ab25ao25Want); ab25ao25Err != nil {
			ab25ao25T.Fatal(ab25ao25Err)
		}
		if !reflect.DeepEqual(ab25ao25Got, ab25ao25Want) {
			ab25ao25T.Fatalf("%s:\n got %v\nwant %v", ab25ao25Case.ID, ab25ao25Got, ab25ao25Want)
		}
	}
}

func TestAb25Ao25Room(ab25ao25T *testing.T) {
	for ab25ao25Byte := 0; ab25ao25Byte <= 255; ab25ao25Byte++ {
		ab25ao25Out := ab25ao25Request(ab25ao25T, fmt.Sprintf(`{"room_temperature":[%d]}`, ab25ao25Byte))
		ab25ao25F := ab25ao25FieldOf(ab25ao25T, ab25ao25Out, "room_temperature")
		ab25ao25Want := float64(ab25ao25Byte / 2)
		if ab25ao25Byte%2 == 1 {
			ab25ao25Want += 0.5
		}
		if ab25ao25F["number"] != ab25ao25Want || ab25ao25F["state"] != "decoded" ||
			ab25ao25F["candidate_domain"] != true || ab25ao25F["calculation"] != "Proven" || ab25ao25F["grade"] != "Hypothesis" {
			ab25ao25T.Fatalf("room %d: %v", ab25ao25Byte, ab25ao25F)
		}
		if _, ab25ao25Has := ab25ao25F["label"]; ab25ao25Has {
			ab25ao25T.Fatalf("room %d carries a label (sentinel)", ab25ao25Byte)
		}
	}
	ab25ao25Cases := map[int]float64{0: 0, 1: 0.5, 255: 127.5}
	for ab25ao25Byte, ab25ao25Want := range ab25ao25Cases {
		ab25ao25F := ab25ao25FieldOf(ab25ao25T, ab25ao25Request(ab25ao25T, fmt.Sprintf(`{"room_temperature":[%d]}`, ab25ao25Byte)), "room_temperature")
		if ab25ao25F["number"] != ab25ao25Want {
			ab25ao25T.Fatalf("room %d = %v, want %v", ab25ao25Byte, ab25ao25F["number"], ab25ao25Want)
		}
	}
}

func TestAb25Ao25Fan(ab25ao25T *testing.T) {
	for ab25ao25Byte := 0; ab25ao25Byte <= 255; ab25ao25Byte++ {
		ab25ao25F := ab25ao25FieldOf(ab25ao25T, ab25ao25Request(ab25ao25T, fmt.Sprintf(`{"fan":[%d]}`, ab25ao25Byte)), "fan")
		if ab25ao25F["number"] != float64(ab25ao25Byte) || ab25ao25F["state"] != "decoded" {
			ab25ao25T.Fatalf("fan %d: %v", ab25ao25Byte, ab25ao25F)
		}
		for _, ab25ao25Key := range []string{"label", "calculation", "candidate_domain", "unit", "reserved", "flags"} {
			if _, ab25ao25Has := ab25ao25F[ab25ao25Key]; ab25ao25Has {
				ab25ao25T.Fatalf("fan %d carries %s", ab25ao25Byte, ab25ao25Key)
			}
		}
	}
}

func TestAb25Ao25Mode(ab25ao25T *testing.T) {
	ab25ao25Labels := map[int]string{
		0: "no-heating-demand",
		2: "source-heating-demand",
		3: "no-source-cooling-demand-sink-demand",
		4: "source-cooling-demand",
	}
	for ab25ao25Byte := 0; ab25ao25Byte <= 255; ab25ao25Byte++ {
		ab25ao25F := ab25ao25FieldOf(ab25ao25T, ab25ao25Request(ab25ao25T, fmt.Sprintf(`{"mode":[%d]}`, ab25ao25Byte)), "mode")
		switch {
		case ab25ao25Byte == 1:
			if ab25ao25F["state"] != "label-conflict" || ab25ao25F["label"] != "no-source-demand-sink-demand" || ab25ao25F["number"] != float64(1) {
				ab25ao25T.Fatalf("mode 1: %v", ab25ao25F)
			}
		case ab25ao25Byte <= 4:
			if ab25ao25F["state"] != "decoded" || ab25ao25F["label"] != ab25ao25Labels[ab25ao25Byte] || ab25ao25F["number"] != float64(ab25ao25Byte) {
				ab25ao25T.Fatalf("mode %d: %v", ab25ao25Byte, ab25ao25F)
			}
		default:
			if ab25ao25F["state"] != "unmapped" {
				ab25ao25T.Fatalf("mode %d: %v", ab25ao25Byte, ab25ao25F)
			}
			for _, ab25ao25Key := range []string{"number", "label"} {
				if _, ab25ao25Has := ab25ao25F[ab25ao25Key]; ab25ao25Has {
					ab25ao25T.Fatalf("mode %d unmapped carries %s", ab25ao25Byte, ab25ao25Key)
				}
			}
		}
		if ab25ao25F["grade"] != "Hypothesis" {
			ab25ao25T.Fatalf("mode %d grade %v", ab25ao25Byte, ab25ao25F["grade"])
		}
	}
}

func TestAb25Ao25AutomaticBits(ab25ao25T *testing.T) {
	for ab25ao25Byte := 0; ab25ao25Byte <= 255; ab25ao25Byte++ {
		ab25ao25F := ab25ao25FieldOf(ab25ao25T, ab25ao25Request(ab25ao25T, fmt.Sprintf(`{"automatic":[%d]}`, ab25ao25Byte)), "automatic")
		ab25ao25Label, ab25ao25Number := "disabled", 0
		if ab25ao25Byte%2 == 1 {
			ab25ao25Label, ab25ao25Number = "enabled", 1
		}
		ab25ao25Reserved := ab25ao25Byte - ab25ao25Byte%2
		if ab25ao25F["label"] != ab25ao25Label || ab25ao25F["number"] != float64(ab25ao25Number) ||
			ab25ao25F["reserved"] != float64(ab25ao25Reserved) || ab25ao25F["state"] != "decoded" {
			ab25ao25T.Fatalf("automatic %d: %v", ab25ao25Byte, ab25ao25F)
		}
	}
	ab25ao25One := ab25ao25FieldOf(ab25ao25T, ab25ao25Request(ab25ao25T, `{"automatic":[1]}`), "automatic")
	ab25ao25Hi := ab25ao25FieldOf(ab25ao25T, ab25ao25Request(ab25ao25T, `{"automatic":[128]}`), "automatic")
	if ab25ao25One["label"] != "enabled" || ab25ao25One["reserved"] != float64(0) || ab25ao25Hi["label"] != "disabled" || ab25ao25Hi["reserved"] != float64(128) {
		ab25ao25T.Fatalf("automatic 01/80: %v %v", ab25ao25One, ab25ao25Hi)
	}
}

func TestAb25Ao25PositiveBits(ab25ao25T *testing.T) {
	for ab25ao25Byte := 0; ab25ao25Byte <= 255; ab25ao25Byte++ {
		ab25ao25Out := ab25ao25Run(ab25ao25T, ab25ao25Input(ab25ao25Context("positive"), fmt.Sprintf(`{"positive_flags":[%d]}`, ab25ao25Byte)))
		ab25ao25F := ab25ao25FieldOf(ab25ao25T, ab25ao25Out, "positive_flags")
		ab25ao25Flags := ab25ao25F["flags"].([]any)
		if len(ab25ao25Flags) != 3 {
			ab25ao25T.Fatalf("positive %d has %d flags", ab25ao25Byte, len(ab25ao25Flags))
		}
		ab25ao25Reserved := 0
		for ab25ao25Bit := 0; ab25ao25Bit < 8; ab25ao25Bit++ {
			ab25ao25Set := (ab25ao25Byte>>ab25ao25Bit)&1 == 1
			if ab25ao25Bit < 3 {
				if ab25ao25Flags[ab25ao25Bit] != ab25ao25Set {
					ab25ao25T.Fatalf("positive %d flag %d", ab25ao25Byte, ab25ao25Bit)
				}
			} else if ab25ao25Set {
				ab25ao25Reserved += 1 << ab25ao25Bit
			}
		}
		if ab25ao25F["reserved"] != float64(ab25ao25Reserved) {
			ab25ao25T.Fatalf("positive %d reserved %v want %d", ab25ao25Byte, ab25ao25F["reserved"], ab25ao25Reserved)
		}
		if _, ab25ao25Has := ab25ao25F["label"]; ab25ao25Has {
			ab25ao25T.Fatalf("positive %d invented a role label", ab25ao25Byte)
		}
		if _, ab25ao25Has := ab25ao25F["number"]; ab25ao25Has {
			ab25ao25T.Fatalf("positive %d carries number", ab25ao25Byte)
		}
	}
}

func TestAb25Ao25WidthAndAbsence(ab25ao25T *testing.T) {
	ab25ao25Bad := map[string]string{
		"empty":    `[]`,
		"overlong": `[2,2]`,
	}
	ab25ao25Want := map[string]string{"empty": "incomplete", "overlong": "overlong"}
	for ab25ao25Name, ab25ao25Value := range ab25ao25Bad {
		for _, ab25ao25Target := range []string{"mode", "fan", "room_temperature", "automatic"} {
			ab25ao25Fields := map[string]string{"mode": `[2]`, "fan": `[7]`, "room_temperature": `[4]`, "automatic": `[1]`}
			ab25ao25Fields[ab25ao25Target] = ab25ao25Value
			ab25ao25Parts := []string{}
			for ab25ao25Key, ab25ao25Val := range ab25ao25Fields {
				ab25ao25Parts = append(ab25ao25Parts, `"`+ab25ao25Key+`":`+ab25ao25Val)
			}
			ab25ao25Out := ab25ao25Request(ab25ao25T, `{`+strings.Join(ab25ao25Parts, ",")+`}`)
			if ab25ao25Out["state"] != "partial" {
				ab25ao25T.Fatalf("%s %s: state %v", ab25ao25Name, ab25ao25Target, ab25ao25Out["state"])
			}
			for ab25ao25Key := range ab25ao25Fields {
				ab25ao25F := ab25ao25FieldOf(ab25ao25T, ab25ao25Out, ab25ao25Key)
				if ab25ao25Key == ab25ao25Target {
					if ab25ao25F["state"] != ab25ao25Want[ab25ao25Name] {
						ab25ao25T.Fatalf("%s %s: %v", ab25ao25Name, ab25ao25Key, ab25ao25F)
					}
					if _, ab25ao25Has := ab25ao25F["number"]; ab25ao25Has {
						ab25ao25T.Fatalf("%s %s carries number", ab25ao25Name, ab25ao25Key)
					}
				} else if ab25ao25F["state"] != "decoded" {
					ab25ao25T.Fatalf("sibling %s erased by %s %s", ab25ao25Key, ab25ao25Name, ab25ao25Target)
				}
			}
		}
	}
	for _, ab25ao25Target := range []string{"mode", "fan", "room_temperature", "automatic"} {
		ab25ao25Fields := map[string]string{"mode": `"mode":[2]`, "fan": `"fan":[7]`, "room_temperature": `"room_temperature":[4]`, "automatic": `"automatic":[1]`}
		delete(ab25ao25Fields, ab25ao25Target)
		ab25ao25Parts := []string{}
		for _, ab25ao25Val := range ab25ao25Fields {
			ab25ao25Parts = append(ab25ao25Parts, ab25ao25Val)
		}
		ab25ao25Out := ab25ao25Request(ab25ao25T, `{`+strings.Join(ab25ao25Parts, ",")+`}`)
		ab25ao25F := ab25ao25FieldOf(ab25ao25T, ab25ao25Out, ab25ao25Target)
		if ab25ao25F["state"] != "absent" || ab25ao25F["raw"] != nil {
			ab25ao25T.Fatalf("omitted %s: %v", ab25ao25Target, ab25ao25F)
		}
		if _, ab25ao25Has := ab25ao25F["number"]; ab25ao25Has {
			ab25ao25T.Fatalf("omitted %s was filled with a value", ab25ao25Target)
		}
		if ab25ao25Out["state"] != "partial" {
			ab25ao25T.Fatalf("omitted %s: state %v", ab25ao25Target, ab25ao25Out["state"])
		}
	}
	ab25ao25None := ab25ao25Request(ab25ao25T, `{}`)
	if ab25ao25None["state"] != "unresolved" {
		ab25ao25T.Fatalf("no fields state %v", ab25ao25None["state"])
	}
	ab25ao25Empty := ab25ao25Run(ab25ao25T, ab25ao25Input(ab25ao25Context("positive"), `{"positive_flags":[]}`))
	ab25ao25Missing := ab25ao25Run(ab25ao25T, ab25ao25Input(ab25ao25Context("positive"), `{}`))
	ab25ao25EF := ab25ao25FieldOf(ab25ao25T, ab25ao25Empty, "positive_flags")
	ab25ao25MF := ab25ao25FieldOf(ab25ao25T, ab25ao25Missing, "positive_flags")
	if ab25ao25EF["state"] != "incomplete" || ab25ao25MF["state"] != "absent" || reflect.DeepEqual(ab25ao25EF["raw"], ab25ao25MF["raw"]) {
		ab25ao25T.Fatalf("empty and missing not distinct: %v %v", ab25ao25EF, ab25ao25MF)
	}
	if _, ab25ao25Has := ab25ao25EF["flags"]; ab25ao25Has {
		ab25ao25T.Fatal("empty positive_flags produced flags")
	}
}

func TestAb25Ao25CrossDirection(ab25ao25T *testing.T) {
	ab25ao25Req := ab25ao25Run(ab25ao25T, ab25ao25Input(ab25ao25Context("request"), `{"positive_flags":[7],"mode":[2]}`))
	ab25ao25Fields := ab25ao25Req["fields"].(map[string]any)
	if _, ab25ao25Has := ab25ao25Fields["positive_flags"]; ab25ao25Has || len(ab25ao25Fields) != 4 {
		ab25ao25T.Fatalf("request interpreted positive_flags: %v", ab25ao25Fields)
	}
	if !reflect.DeepEqual(ab25ao25Req["raw"].(map[string]any)["positive_flags"], []any{float64(7)}) {
		ab25ao25T.Fatal("positive_flags raw lost in request")
	}
	ab25ao25Pos := ab25ao25Run(ab25ao25T, ab25ao25Input(ab25ao25Context("positive"), `{"positive_flags":[7],"mode":[2],"fan":[1],"room_temperature":[1],"automatic":[1]}`))
	ab25ao25PF := ab25ao25Pos["fields"].(map[string]any)
	if len(ab25ao25PF) != 1 || ab25ao25Pos["state"] != "complete" {
		ab25ao25T.Fatalf("positive reused request fields: %v", ab25ao25PF)
	}
	if len(ab25ao25Pos["raw"].(map[string]any)) != 5 {
		ab25ao25T.Fatal("positive raw lost sibling bytes")
	}
}

func TestAb25Ao25ContextAxes(ab25ao25T *testing.T) {
	ab25ao25Make := func(ab25ao25P, ab25ao25S, ab25ao25O, ab25ao25D, ab25ao25V string) string {
		return `{"primary":` + ab25ao25P + `,"secondary":` + ab25ao25S + `,"operation":` + ab25ao25O + `,"direction":` + ab25ao25D + `,"view":` + ab25ao25V + `,"family":"x","variant":{"k":[1]}}`
	}
	ab25ao25Fields := `{"mode":[2]}`
	ab25ao25Control := ab25ao25Run(ab25ao25T, ab25ao25Input(ab25ao25Make("181", "34", "1", `"request"`, `"isolated-fields"`), ab25ao25Fields))
	if ab25ao25Control["state"] != "partial" || ab25ao25FieldOf(ab25ao25T, ab25ao25Control, "mode")["label"] != "source-heating-demand" {
		ab25ao25T.Fatalf("control: %v", ab25ao25Control)
	}
	ab25ao25Axes := []struct {
		name, ctx, reason string
	}{
		{"primary", ab25ao25Make("180", "34", "1", `"request"`, `"isolated-fields"`), "primary"},
		{"secondary", ab25ao25Make("181", "35", "1", `"request"`, `"isolated-fields"`), "secondary"},
		{"operation-high", ab25ao25Make("181", "34", "256", `"request"`, `"isolated-fields"`), "operation"},
		{"operation-negative", ab25ao25Make("181", "34", "-1", `"request"`, `"isolated-fields"`), "operation"},
		{"direction", ab25ao25Make("181", "34", "1", `"response"`, `"isolated-fields"`), "direction"},
		{"direction-case", ab25ao25Make("181", "34", "1", `"Request"`, `"isolated-fields"`), "direction"},
		{"view", ab25ao25Make("181", "34", "1", `"request"`, `"payload-length"`), "view"},
		{"view-empty", ab25ao25Make("181", "34", "1", `"request"`, `""`), "view"},
	}
	for _, ab25ao25Axis := range ab25ao25Axes {
		ab25ao25Out := ab25ao25Run(ab25ao25T, ab25ao25Input(ab25ao25Axis.ctx, ab25ao25Fields))
		if ab25ao25Out["state"] != "unknown" || !reflect.DeepEqual(ab25ao25Out["unknown_reasons"], []any{ab25ao25Axis.reason}) {
			ab25ao25T.Fatalf("%s: %v", ab25ao25Axis.name, ab25ao25Out)
		}
		if len(ab25ao25Out["fields"].(map[string]any)) != 0 || ab25ao25Out["qualified"] != false || ab25ao25Out["physical"] != "Unknown" {
			ab25ao25T.Fatalf("%s: guessed selection %v", ab25ao25Axis.name, ab25ao25Out)
		}
		if !reflect.DeepEqual(ab25ao25Out["raw"], ab25ao25Control["raw"]) {
			ab25ao25T.Fatalf("%s: raw lost", ab25ao25Axis.name)
		}
	}
	for _, ab25ao25Op := range []string{"0", "1", "127", "255"} {
		ab25ao25Out := ab25ao25Run(ab25ao25T, ab25ao25Input(ab25ao25Make("181", "34", ab25ao25Op, `"request"`, `"isolated-fields"`), ab25ao25Fields))
		if ab25ao25Out["state"] != "partial" {
			ab25ao25T.Fatalf("operation %s: %v", ab25ao25Op, ab25ao25Out["state"])
		}
	}
	ab25ao25Multi := ab25ao25Run(ab25ao25T, ab25ao25Input(ab25ao25Make("1", "2", "300", `"x"`, `"y"`), ab25ao25Fields))
	if !reflect.DeepEqual(ab25ao25Multi["unknown_reasons"], []any{"primary", "secondary", "operation", "direction", "view"}) {
		ab25ao25T.Fatalf("reasons: %v", ab25ao25Multi["unknown_reasons"])
	}
	ab25ao25Fam := []string{`"a"`, `"b"`, `null`, `{"deep":[1,2,{"z":true}]}`, `1.50`}
	for _, ab25ao25Value := range ab25ao25Fam {
		ab25ao25Ctx := `{"primary":181,"secondary":34,"operation":1,"direction":"request","view":"isolated-fields","family":` + ab25ao25Value + `,"variant":` + ab25ao25Value + `}`
		ab25ao25Out := ab25ao25Run(ab25ao25T, ab25ao25Input(ab25ao25Ctx, ab25ao25Fields))
		var ab25ao25Want map[string]any
		if ab25ao25Err := json.Unmarshal([]byte(ab25ao25Ctx), &ab25ao25Want); ab25ao25Err != nil {
			ab25ao25T.Fatal(ab25ao25Err)
		}
		if !reflect.DeepEqual(ab25ao25Out["context"], ab25ao25Want) || ab25ao25Out["qualified"] != false || ab25ao25Out["physical"] != "Unknown" {
			ab25ao25T.Fatalf("context %s: %v", ab25ao25Value, ab25ao25Out)
		}
	}
}

func TestAb25Ao25InvalidInput(ab25ao25T *testing.T) {
	ab25ao25Ctx := ab25ao25Context("request")
	ab25ao25Cases := map[string]string{
		"empty":             ``,
		"malformed":         `{"context":`,
		"not-object":        `[]`,
		"trailing":          ab25ao25Input(ab25ao25Ctx, `{}`) + ` {}`,
		"dup-top":           `{"context":` + ab25ao25Ctx + `,"fields":{},"fields":{}}`,
		"dup-context":       `{"context":{"primary":181,"primary":181,"secondary":34,"operation":1,"direction":"request","view":"isolated-fields"},"fields":{}}`,
		"dup-field":         ab25ao25Input(ab25ao25Ctx, `{"mode":[1],"mode":[2]}`),
		"dup-nested-opaque": `{"context":{"primary":181,"secondary":34,"operation":1,"direction":"request","view":"isolated-fields","f":{"a":1,"a":2}},"fields":{}}`,
		"dup-arguments":     `{"context":` + ab25ao25Ctx + `,"fields":{},"arguments":{"a":1,"a":2}}`,
		"element-256":       ab25ao25Input(ab25ao25Ctx, `{"mode":[256]}`),
		"element-negative":  ab25ao25Input(ab25ao25Ctx, `{"mode":[-1]}`),
		"element-fraction":  ab25ao25Input(ab25ao25Ctx, `{"mode":[1.5]}`),
		"element-exponent":  ab25ao25Input(ab25ao25Ctx, `{"mode":[1e1]}`),
		"element-string":    ab25ao25Input(ab25ao25Ctx, `{"mode":["1"]}`),
		"element-bool":      ab25ao25Input(ab25ao25Ctx, `{"mode":[true]}`),
		"element-null":      ab25ao25Input(ab25ao25Ctx, `{"mode":[null]}`),
		"element-nested":    ab25ao25Input(ab25ao25Ctx, `{"mode":[[1]]}`),
		"field-scalar":      ab25ao25Input(ab25ao25Ctx, `{"mode":2}`),
		"field-null":        ab25ao25Input(ab25ao25Ctx, `{"mode":null}`),
		"field-string":      ab25ao25Input(ab25ao25Ctx, `{"mode":"02"}`),
		"unmodelled-bad":    ab25ao25Input(ab25ao25Ctx, `{"tail":[300]}`),
		"fields-array":      ab25ao25Input(ab25ao25Ctx, `[]`),
		"fields-missing":    `{"context":` + ab25ao25Ctx + `}`,
		"context-missing":   `{"fields":{}}`,
		"context-array":     `{"context":[],"fields":{}}`,
		"primary-string":    ab25ao25Input(`{"primary":"181","secondary":34,"operation":1,"direction":"request","view":"isolated-fields"}`, `{}`),
		"primary-float":     ab25ao25Input(`{"primary":181.0,"secondary":34,"operation":1,"direction":"request","view":"isolated-fields"}`, `{}`),
		"primary-missing":   ab25ao25Input(`{"secondary":34,"operation":1,"direction":"request","view":"isolated-fields"}`, `{}`),
		"secondary-missing": ab25ao25Input(`{"primary":181,"operation":1,"direction":"request","view":"isolated-fields"}`, `{}`),
		"operation-string":  ab25ao25Input(`{"primary":181,"secondary":34,"operation":"1","direction":"request","view":"isolated-fields"}`, `{}`),
		"operation-missing": ab25ao25Input(`{"primary":181,"secondary":34,"direction":"request","view":"isolated-fields"}`, `{}`),
		"direction-number":  ab25ao25Input(`{"primary":181,"secondary":34,"operation":1,"direction":1,"view":"isolated-fields"}`, `{}`),
		"direction-missing": ab25ao25Input(`{"primary":181,"secondary":34,"operation":1,"view":"isolated-fields"}`, `{}`),
		"view-missing":      ab25ao25Input(`{"primary":181,"secondary":34,"operation":1,"direction":"request"}`, `{}`),
		"view-number":       ab25ao25Input(`{"primary":181,"secondary":34,"operation":1,"direction":"request","view":2}`, `{}`),
		"arguments-null":    `{"context":` + ab25ao25Ctx + `,"fields":{},"arguments":null}`,
		"arguments-array":   `{"context":` + ab25ao25Ctx + `,"fields":{},"arguments":[]}`,
		"arguments-string":  `{"context":` + ab25ao25Ctx + `,"fields":{},"arguments":{"a":"x"}}`,
		"arguments-float":   `{"context":` + ab25ao25Ctx + `,"fields":{},"arguments":{"a":1.5}}`,
		"arguments-nested":  `{"context":` + ab25ao25Ctx + `,"fields":{},"arguments":{"a":{}}}`,
	}
	for ab25ao25Name, ab25ao25Text := range ab25ao25Cases {
		ab25ao25Out, ab25ao25Err := Ab25Ao25Interpret(json.RawMessage(ab25ao25Text))
		if ab25ao25Err == nil || ab25ao25Out != nil {
			ab25ao25T.Fatalf("%s accepted: %s", ab25ao25Name, ab25ao25Out)
		}
	}
	ab25ao25Run(ab25ao25T, ab25ao25Input(ab25ao25Ctx, `{"mode":[0],"tail":[255,0,1]}`))
}

func TestAb25Ao25UnrecognisedKeysAndArguments(ab25ao25T *testing.T) {
	ab25ao25Out := ab25ao25Run(ab25ao25T, `{"context":`+ab25ao25Context("request")+`,"fields":{"mode":[2,9,9],"tail":[0,255],"extra_empty":[]},"arguments":{"flag":false,"zero":0,"neg":-3},"note":{"a":[1]}}`)
	ab25ao25Raw := ab25ao25Out["raw"].(map[string]any)
	if !reflect.DeepEqual(ab25ao25Raw["tail"], []any{float64(0), float64(255)}) ||
		!reflect.DeepEqual(ab25ao25Raw["extra_empty"], []any{}) ||
		!reflect.DeepEqual(ab25ao25Raw["mode"], []any{float64(2), float64(9), float64(9)}) {
		ab25ao25T.Fatalf("raw: %v", ab25ao25Raw)
	}
	if !reflect.DeepEqual(ab25ao25Out["arguments"], map[string]any{"flag": false, "zero": float64(0), "neg": float64(-3)}) {
		ab25ao25T.Fatalf("arguments: %v", ab25ao25Out["arguments"])
	}
	if !reflect.DeepEqual(ab25ao25Out["unrecognised"], map[string]any{"note": map[string]any{"a": []any{float64(1)}}}) {
		ab25ao25T.Fatalf("unrecognised: %v", ab25ao25Out["unrecognised"])
	}
	if _, ab25ao25Has := ab25ao25FieldOf(ab25ao25T, ab25ao25Out, "mode")["number"]; ab25ao25Has {
		ab25ao25T.Fatal("overlong mode carried a number")
	}
	ab25ao25NoArgs := ab25ao25Request(ab25ao25T, `{}`)
	if _, ab25ao25Has := ab25ao25NoArgs["arguments"]; ab25ao25Has {
		ab25ao25T.Fatal("arguments invented")
	}
	if _, ab25ao25Has := ab25ao25NoArgs["unrecognised"]; ab25ao25Has {
		ab25ao25T.Fatal("unrecognised invented")
	}
	ab25ao25EmptyArgs := ab25ao25Run(ab25ao25T, `{"context":`+ab25ao25Context("request")+`,"fields":{},"arguments":{}}`)
	if !reflect.DeepEqual(ab25ao25EmptyArgs["arguments"], map[string]any{}) {
		ab25ao25T.Fatalf("empty arguments not preserved: %v", ab25ao25EmptyArgs["arguments"])
	}
}

func TestAb25Ao25IndependenceAndRoundTrip(ab25ao25T *testing.T) {
	ab25ao25Text := `{"context":` + ab25ao25Context("request") + `,"fields":{"mode":[2],"fan":[],"tail":[1]},"arguments":{"a":1}}`
	ab25ao25Input := json.RawMessage(ab25ao25Text)
	ab25ao25First, ab25ao25Err := Ab25Ao25Interpret(ab25ao25Input)
	if ab25ao25Err != nil {
		ab25ao25T.Fatal(ab25ao25Err)
	}
	ab25ao25Snapshot := string(ab25ao25First)
	for ab25ao25Index := range ab25ao25Input {
		ab25ao25Input[ab25ao25Index] = ' '
	}
	for ab25ao25Index := range ab25ao25First {
		ab25ao25First[ab25ao25Index] = 'x'
	}
	ab25ao25Second, ab25ao25Err := Ab25Ao25Interpret(json.RawMessage(ab25ao25Text))
	if ab25ao25Err != nil || string(ab25ao25Second) != ab25ao25Snapshot {
		ab25ao25T.Fatalf("later call changed: %v\n%s\n%s", ab25ao25Err, ab25ao25Second, ab25ao25Snapshot)
	}
	var ab25ao25Decoded map[string]any
	if ab25ao25Err = json.Unmarshal(ab25ao25Second, &ab25ao25Decoded); ab25ao25Err != nil {
		ab25ao25T.Fatal(ab25ao25Err)
	}
	ab25ao25Decoded["raw"].(map[string]any)["mode"].([]any)[0] = float64(99)
	ab25ao25Decoded["context"].(map[string]any)["primary"] = float64(1)
	ab25ao25Decoded["fields"].(map[string]any)["mode"].(map[string]any)["raw"].([]any)[0] = float64(98)
	ab25ao25Third, _ := Ab25Ao25Interpret(json.RawMessage(ab25ao25Text))
	if string(ab25ao25Third) != ab25ao25Snapshot {
		ab25ao25T.Fatal("mutating a decoded result altered a later call")
	}
	var ab25ao25Fresh map[string]any
	_ = json.Unmarshal(ab25ao25Third, &ab25ao25Fresh)
	if ab25ao25Fresh["raw"].(map[string]any)["mode"].([]any)[0] != float64(2) ||
		ab25ao25Fresh["fields"].(map[string]any)["mode"].(map[string]any)["raw"].([]any)[0] != float64(2) {
		ab25ao25T.Fatal("raw and field raw are not independent copies")
	}
	ab25ao25Round, _ := json.Marshal(ab25ao25Fresh)
	var ab25ao25Again map[string]any
	_ = json.Unmarshal(ab25ao25Round, &ab25ao25Again)
	if !reflect.DeepEqual(ab25ao25Fresh, ab25ao25Again) {
		ab25ao25T.Fatal("round trip changed result")
	}
	ab25ao25F := ab25ao25Fresh["fields"].(map[string]any)
	if ab25ao25F["fan"].(map[string]any)["raw"] == nil || ab25ao25F["room_temperature"].(map[string]any)["raw"] != nil {
		ab25ao25T.Fatal("empty and absent raw not distinct after round trip")
	}
	if ab25ao25Fresh["state"] != "partial" || ab25ao25Fresh["qualified"] != false || ab25ao25Fresh["physical"] != "Unknown" {
		ab25ao25T.Fatalf("top-level: %v", ab25ao25Fresh)
	}
}
