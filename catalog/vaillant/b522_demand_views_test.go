package vaillant

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func b522DemandViewRun(b522DemandViewT *testing.T, b522DemandViewText string) map[string]any {
	b522DemandViewT.Helper()
	b522DemandViewOut, b522DemandViewErr := B522DemandViewInterpret(json.RawMessage(b522DemandViewText))
	if b522DemandViewErr != nil {
		b522DemandViewT.Fatalf("unexpected error: %v", b522DemandViewErr)
	}
	var b522DemandViewMap map[string]any
	if b522DemandViewErr = json.Unmarshal(b522DemandViewOut, &b522DemandViewMap); b522DemandViewErr != nil {
		b522DemandViewT.Fatalf("output is not an object: %v", b522DemandViewErr)
	}
	return b522DemandViewMap
}

func b522DemandViewInput(b522DemandViewContext string, b522DemandViewFields string) string {
	return `{"context":` + b522DemandViewContext + `,"fields":` + b522DemandViewFields + `}`
}

func b522DemandViewContext(b522DemandViewDirection string) string {
	return `{"primary":181,"secondary":34,"operation":1,"direction":"` + b522DemandViewDirection + `","view":"isolated-fields"}`
}

func b522DemandViewFieldOf(b522DemandViewT *testing.T, b522DemandViewOut map[string]any, b522DemandViewName string) map[string]any {
	b522DemandViewT.Helper()
	b522DemandViewFields := b522DemandViewOut["fields"].(map[string]any)
	b522DemandViewField, b522DemandViewOK := b522DemandViewFields[b522DemandViewName].(map[string]any)
	if !b522DemandViewOK {
		b522DemandViewT.Fatalf("field %q missing in %v", b522DemandViewName, b522DemandViewFields)
	}
	return b522DemandViewField
}

func b522DemandViewRequest(b522DemandViewT *testing.T, b522DemandViewFields string) map[string]any {
	b522DemandViewT.Helper()
	return b522DemandViewRun(b522DemandViewT, b522DemandViewInput(b522DemandViewContext("request"), b522DemandViewFields))
}

func TestAb25Ao25FixtureCases(b522DemandViewT *testing.T) {
	b522DemandViewData, b522DemandViewErr := os.ReadFile("testdata/ab25ao25_b522_demand_views.json")
	if b522DemandViewErr != nil {
		b522DemandViewT.Fatal(b522DemandViewErr)
	}
	var b522DemandViewFile struct {
		Cases []struct {
			ID       string          `json:"id"`
			Input    json.RawMessage `json:"input"`
			Expected json.RawMessage `json:"expected"`
		} `json:"ab25ao25_cases"`
	}
	if b522DemandViewErr = json.Unmarshal(b522DemandViewData, &b522DemandViewFile); b522DemandViewErr != nil {
		b522DemandViewT.Fatal(b522DemandViewErr)
	}
	if len(b522DemandViewFile.Cases) != 4 {
		b522DemandViewT.Fatalf("fixture cases = %d", len(b522DemandViewFile.Cases))
	}
	for _, b522DemandViewCase := range b522DemandViewFile.Cases {
		if !strings.HasPrefix(b522DemandViewCase.ID, "ab25ao25") {
			b522DemandViewT.Fatalf("fixture id %q", b522DemandViewCase.ID)
		}
		b522DemandViewGot := b522DemandViewRun(b522DemandViewT, string(b522DemandViewCase.Input))
		var b522DemandViewWant map[string]any
		if b522DemandViewErr = json.Unmarshal(b522DemandViewCase.Expected, &b522DemandViewWant); b522DemandViewErr != nil {
			b522DemandViewT.Fatal(b522DemandViewErr)
		}
		if !reflect.DeepEqual(b522DemandViewGot, b522DemandViewWant) {
			b522DemandViewT.Fatalf("%s:\n got %v\nwant %v", b522DemandViewCase.ID, b522DemandViewGot, b522DemandViewWant)
		}
	}
}

func TestAb25Ao25Room(b522DemandViewT *testing.T) {
	for b522DemandViewByte := 0; b522DemandViewByte <= 255; b522DemandViewByte++ {
		b522DemandViewOut := b522DemandViewRequest(b522DemandViewT, fmt.Sprintf(`{"room_temperature":[%d]}`, b522DemandViewByte))
		b522DemandViewF := b522DemandViewFieldOf(b522DemandViewT, b522DemandViewOut, "room_temperature")
		b522DemandViewWant := float64(b522DemandViewByte / 2)
		if b522DemandViewByte%2 == 1 {
			b522DemandViewWant += 0.5
		}
		if b522DemandViewF["number"] != b522DemandViewWant || b522DemandViewF["state"] != "decoded" ||
			b522DemandViewF["candidate_domain"] != true || b522DemandViewF["calculation"] != "Proven" || b522DemandViewF["grade"] != "Hypothesis" {
			b522DemandViewT.Fatalf("room %d: %v", b522DemandViewByte, b522DemandViewF)
		}
		if _, b522DemandViewHas := b522DemandViewF["label"]; b522DemandViewHas {
			b522DemandViewT.Fatalf("room %d carries a label (sentinel)", b522DemandViewByte)
		}
	}
	b522DemandViewCases := map[int]float64{0: 0, 1: 0.5, 255: 127.5}
	for b522DemandViewByte, b522DemandViewWant := range b522DemandViewCases {
		b522DemandViewF := b522DemandViewFieldOf(b522DemandViewT, b522DemandViewRequest(b522DemandViewT, fmt.Sprintf(`{"room_temperature":[%d]}`, b522DemandViewByte)), "room_temperature")
		if b522DemandViewF["number"] != b522DemandViewWant {
			b522DemandViewT.Fatalf("room %d = %v, want %v", b522DemandViewByte, b522DemandViewF["number"], b522DemandViewWant)
		}
	}
}

func TestAb25Ao25Fan(b522DemandViewT *testing.T) {
	for b522DemandViewByte := 0; b522DemandViewByte <= 255; b522DemandViewByte++ {
		b522DemandViewF := b522DemandViewFieldOf(b522DemandViewT, b522DemandViewRequest(b522DemandViewT, fmt.Sprintf(`{"fan":[%d]}`, b522DemandViewByte)), "fan")
		if b522DemandViewF["number"] != float64(b522DemandViewByte) || b522DemandViewF["state"] != "decoded" {
			b522DemandViewT.Fatalf("fan %d: %v", b522DemandViewByte, b522DemandViewF)
		}
		for _, b522DemandViewKey := range []string{"label", "calculation", "candidate_domain", "unit", "reserved", "flags"} {
			if _, b522DemandViewHas := b522DemandViewF[b522DemandViewKey]; b522DemandViewHas {
				b522DemandViewT.Fatalf("fan %d carries %s", b522DemandViewByte, b522DemandViewKey)
			}
		}
	}
}

func TestAb25Ao25Mode(b522DemandViewT *testing.T) {
	b522DemandViewLabels := map[int]string{
		0: "no-heating-demand",
		2: "source-heating-demand",
		3: "no-source-cooling-demand-sink-demand",
		4: "source-cooling-demand",
	}
	for b522DemandViewByte := 0; b522DemandViewByte <= 255; b522DemandViewByte++ {
		b522DemandViewF := b522DemandViewFieldOf(b522DemandViewT, b522DemandViewRequest(b522DemandViewT, fmt.Sprintf(`{"mode":[%d]}`, b522DemandViewByte)), "mode")
		switch {
		case b522DemandViewByte == 1:
			if b522DemandViewF["state"] != "label-conflict" || b522DemandViewF["label"] != "no-source-demand-sink-demand" || b522DemandViewF["number"] != float64(1) {
				b522DemandViewT.Fatalf("mode 1: %v", b522DemandViewF)
			}
		case b522DemandViewByte <= 4:
			if b522DemandViewF["state"] != "decoded" || b522DemandViewF["label"] != b522DemandViewLabels[b522DemandViewByte] || b522DemandViewF["number"] != float64(b522DemandViewByte) {
				b522DemandViewT.Fatalf("mode %d: %v", b522DemandViewByte, b522DemandViewF)
			}
		default:
			if b522DemandViewF["state"] != "unmapped" {
				b522DemandViewT.Fatalf("mode %d: %v", b522DemandViewByte, b522DemandViewF)
			}
			for _, b522DemandViewKey := range []string{"number", "label"} {
				if _, b522DemandViewHas := b522DemandViewF[b522DemandViewKey]; b522DemandViewHas {
					b522DemandViewT.Fatalf("mode %d unmapped carries %s", b522DemandViewByte, b522DemandViewKey)
				}
			}
		}
		if b522DemandViewF["grade"] != "Hypothesis" {
			b522DemandViewT.Fatalf("mode %d grade %v", b522DemandViewByte, b522DemandViewF["grade"])
		}
	}
}

func TestAb25Ao25AutomaticBits(b522DemandViewT *testing.T) {
	for b522DemandViewByte := 0; b522DemandViewByte <= 255; b522DemandViewByte++ {
		b522DemandViewF := b522DemandViewFieldOf(b522DemandViewT, b522DemandViewRequest(b522DemandViewT, fmt.Sprintf(`{"automatic":[%d]}`, b522DemandViewByte)), "automatic")
		b522DemandViewLabel, b522DemandViewNumber := "disabled", 0
		if b522DemandViewByte%2 == 1 {
			b522DemandViewLabel, b522DemandViewNumber = "enabled", 1
		}
		b522DemandViewReserved := b522DemandViewByte - b522DemandViewByte%2
		if b522DemandViewF["label"] != b522DemandViewLabel || b522DemandViewF["number"] != float64(b522DemandViewNumber) ||
			b522DemandViewF["reserved"] != float64(b522DemandViewReserved) || b522DemandViewF["state"] != "decoded" {
			b522DemandViewT.Fatalf("automatic %d: %v", b522DemandViewByte, b522DemandViewF)
		}
	}
	b522DemandViewOne := b522DemandViewFieldOf(b522DemandViewT, b522DemandViewRequest(b522DemandViewT, `{"automatic":[1]}`), "automatic")
	b522DemandViewHi := b522DemandViewFieldOf(b522DemandViewT, b522DemandViewRequest(b522DemandViewT, `{"automatic":[128]}`), "automatic")
	if b522DemandViewOne["label"] != "enabled" || b522DemandViewOne["reserved"] != float64(0) || b522DemandViewHi["label"] != "disabled" || b522DemandViewHi["reserved"] != float64(128) {
		b522DemandViewT.Fatalf("automatic 01/80: %v %v", b522DemandViewOne, b522DemandViewHi)
	}
}

func TestAb25Ao25PositiveBits(b522DemandViewT *testing.T) {
	for b522DemandViewByte := 0; b522DemandViewByte <= 255; b522DemandViewByte++ {
		b522DemandViewOut := b522DemandViewRun(b522DemandViewT, b522DemandViewInput(b522DemandViewContext("positive"), fmt.Sprintf(`{"positive_flags":[%d]}`, b522DemandViewByte)))
		b522DemandViewF := b522DemandViewFieldOf(b522DemandViewT, b522DemandViewOut, "positive_flags")
		b522DemandViewFlags := b522DemandViewF["flags"].([]any)
		if len(b522DemandViewFlags) != 3 {
			b522DemandViewT.Fatalf("positive %d has %d flags", b522DemandViewByte, len(b522DemandViewFlags))
		}
		b522DemandViewReserved := 0
		for b522DemandViewBit := 0; b522DemandViewBit < 8; b522DemandViewBit++ {
			b522DemandViewSet := (b522DemandViewByte>>b522DemandViewBit)&1 == 1
			if b522DemandViewBit < 3 {
				if b522DemandViewFlags[b522DemandViewBit] != b522DemandViewSet {
					b522DemandViewT.Fatalf("positive %d flag %d", b522DemandViewByte, b522DemandViewBit)
				}
			} else if b522DemandViewSet {
				b522DemandViewReserved += 1 << b522DemandViewBit
			}
		}
		if b522DemandViewF["reserved"] != float64(b522DemandViewReserved) {
			b522DemandViewT.Fatalf("positive %d reserved %v want %d", b522DemandViewByte, b522DemandViewF["reserved"], b522DemandViewReserved)
		}
		if _, b522DemandViewHas := b522DemandViewF["label"]; b522DemandViewHas {
			b522DemandViewT.Fatalf("positive %d invented a role label", b522DemandViewByte)
		}
		if _, b522DemandViewHas := b522DemandViewF["number"]; b522DemandViewHas {
			b522DemandViewT.Fatalf("positive %d carries number", b522DemandViewByte)
		}
	}
}

func TestAb25Ao25WidthAndAbsence(b522DemandViewT *testing.T) {
	b522DemandViewBad := map[string]string{
		"empty":    `[]`,
		"overlong": `[2,2]`,
	}
	b522DemandViewWant := map[string]string{"empty": "incomplete", "overlong": "overlong"}
	for b522DemandViewName, b522DemandViewValue := range b522DemandViewBad {
		for _, b522DemandViewTarget := range []string{"mode", "fan", "room_temperature", "automatic"} {
			b522DemandViewFields := map[string]string{"mode": `[2]`, "fan": `[7]`, "room_temperature": `[4]`, "automatic": `[1]`}
			b522DemandViewFields[b522DemandViewTarget] = b522DemandViewValue
			b522DemandViewParts := []string{}
			for b522DemandViewKey, b522DemandViewVal := range b522DemandViewFields {
				b522DemandViewParts = append(b522DemandViewParts, `"`+b522DemandViewKey+`":`+b522DemandViewVal)
			}
			b522DemandViewOut := b522DemandViewRequest(b522DemandViewT, `{`+strings.Join(b522DemandViewParts, ",")+`}`)
			if b522DemandViewOut["state"] != "partial" {
				b522DemandViewT.Fatalf("%s %s: state %v", b522DemandViewName, b522DemandViewTarget, b522DemandViewOut["state"])
			}
			for b522DemandViewKey := range b522DemandViewFields {
				b522DemandViewF := b522DemandViewFieldOf(b522DemandViewT, b522DemandViewOut, b522DemandViewKey)
				if b522DemandViewKey == b522DemandViewTarget {
					if b522DemandViewF["state"] != b522DemandViewWant[b522DemandViewName] {
						b522DemandViewT.Fatalf("%s %s: %v", b522DemandViewName, b522DemandViewKey, b522DemandViewF)
					}
					if _, b522DemandViewHas := b522DemandViewF["number"]; b522DemandViewHas {
						b522DemandViewT.Fatalf("%s %s carries number", b522DemandViewName, b522DemandViewKey)
					}
				} else if b522DemandViewF["state"] != "decoded" {
					b522DemandViewT.Fatalf("sibling %s erased by %s %s", b522DemandViewKey, b522DemandViewName, b522DemandViewTarget)
				}
			}
		}
	}
	for _, b522DemandViewTarget := range []string{"mode", "fan", "room_temperature", "automatic"} {
		b522DemandViewFields := map[string]string{"mode": `"mode":[2]`, "fan": `"fan":[7]`, "room_temperature": `"room_temperature":[4]`, "automatic": `"automatic":[1]`}
		delete(b522DemandViewFields, b522DemandViewTarget)
		b522DemandViewParts := []string{}
		for _, b522DemandViewVal := range b522DemandViewFields {
			b522DemandViewParts = append(b522DemandViewParts, b522DemandViewVal)
		}
		b522DemandViewOut := b522DemandViewRequest(b522DemandViewT, `{`+strings.Join(b522DemandViewParts, ",")+`}`)
		b522DemandViewF := b522DemandViewFieldOf(b522DemandViewT, b522DemandViewOut, b522DemandViewTarget)
		if b522DemandViewF["state"] != "absent" || b522DemandViewF["raw"] != nil {
			b522DemandViewT.Fatalf("omitted %s: %v", b522DemandViewTarget, b522DemandViewF)
		}
		if _, b522DemandViewHas := b522DemandViewF["number"]; b522DemandViewHas {
			b522DemandViewT.Fatalf("omitted %s was filled with a value", b522DemandViewTarget)
		}
		if b522DemandViewOut["state"] != "partial" {
			b522DemandViewT.Fatalf("omitted %s: state %v", b522DemandViewTarget, b522DemandViewOut["state"])
		}
	}
	b522DemandViewNone := b522DemandViewRequest(b522DemandViewT, `{}`)
	if b522DemandViewNone["state"] != "unresolved" {
		b522DemandViewT.Fatalf("no fields state %v", b522DemandViewNone["state"])
	}
	b522DemandViewEmpty := b522DemandViewRun(b522DemandViewT, b522DemandViewInput(b522DemandViewContext("positive"), `{"positive_flags":[]}`))
	b522DemandViewMissing := b522DemandViewRun(b522DemandViewT, b522DemandViewInput(b522DemandViewContext("positive"), `{}`))
	b522DemandViewEF := b522DemandViewFieldOf(b522DemandViewT, b522DemandViewEmpty, "positive_flags")
	b522DemandViewMF := b522DemandViewFieldOf(b522DemandViewT, b522DemandViewMissing, "positive_flags")
	if b522DemandViewEF["state"] != "incomplete" || b522DemandViewMF["state"] != "absent" || reflect.DeepEqual(b522DemandViewEF["raw"], b522DemandViewMF["raw"]) {
		b522DemandViewT.Fatalf("empty and missing not distinct: %v %v", b522DemandViewEF, b522DemandViewMF)
	}
	if _, b522DemandViewHas := b522DemandViewEF["flags"]; b522DemandViewHas {
		b522DemandViewT.Fatal("empty positive_flags produced flags")
	}
}

func TestAb25Ao25CrossDirection(b522DemandViewT *testing.T) {
	b522DemandViewReq := b522DemandViewRun(b522DemandViewT, b522DemandViewInput(b522DemandViewContext("request"), `{"positive_flags":[7],"mode":[2]}`))
	b522DemandViewFields := b522DemandViewReq["fields"].(map[string]any)
	if _, b522DemandViewHas := b522DemandViewFields["positive_flags"]; b522DemandViewHas || len(b522DemandViewFields) != 4 {
		b522DemandViewT.Fatalf("request interpreted positive_flags: %v", b522DemandViewFields)
	}
	if !reflect.DeepEqual(b522DemandViewReq["raw"].(map[string]any)["positive_flags"], []any{float64(7)}) {
		b522DemandViewT.Fatal("positive_flags raw lost in request")
	}
	b522DemandViewPos := b522DemandViewRun(b522DemandViewT, b522DemandViewInput(b522DemandViewContext("positive"), `{"positive_flags":[7],"mode":[2],"fan":[1],"room_temperature":[1],"automatic":[1]}`))
	b522DemandViewPF := b522DemandViewPos["fields"].(map[string]any)
	if len(b522DemandViewPF) != 1 || b522DemandViewPos["state"] != "complete" {
		b522DemandViewT.Fatalf("positive reused request fields: %v", b522DemandViewPF)
	}
	if len(b522DemandViewPos["raw"].(map[string]any)) != 5 {
		b522DemandViewT.Fatal("positive raw lost sibling bytes")
	}
}

func TestAb25Ao25ContextAxes(b522DemandViewT *testing.T) {
	b522DemandViewMake := func(b522DemandViewP, b522DemandViewS, b522DemandViewO, b522DemandViewD, b522DemandViewV string) string {
		return `{"primary":` + b522DemandViewP + `,"secondary":` + b522DemandViewS + `,"operation":` + b522DemandViewO + `,"direction":` + b522DemandViewD + `,"view":` + b522DemandViewV + `,"family":"x","variant":{"k":[1]}}`
	}
	b522DemandViewFields := `{"mode":[2]}`
	b522DemandViewControl := b522DemandViewRun(b522DemandViewT, b522DemandViewInput(b522DemandViewMake("181", "34", "1", `"request"`, `"isolated-fields"`), b522DemandViewFields))
	if b522DemandViewControl["state"] != "partial" || b522DemandViewFieldOf(b522DemandViewT, b522DemandViewControl, "mode")["label"] != "source-heating-demand" {
		b522DemandViewT.Fatalf("control: %v", b522DemandViewControl)
	}
	b522DemandViewAxes := []struct {
		name, ctx, reason string
	}{
		{"primary", b522DemandViewMake("180", "34", "1", `"request"`, `"isolated-fields"`), "primary"},
		{"secondary", b522DemandViewMake("181", "35", "1", `"request"`, `"isolated-fields"`), "secondary"},
		{"operation-high", b522DemandViewMake("181", "34", "256", `"request"`, `"isolated-fields"`), "operation"},
		{"operation-negative", b522DemandViewMake("181", "34", "-1", `"request"`, `"isolated-fields"`), "operation"},
		{"direction", b522DemandViewMake("181", "34", "1", `"response"`, `"isolated-fields"`), "direction"},
		{"direction-case", b522DemandViewMake("181", "34", "1", `"Request"`, `"isolated-fields"`), "direction"},
		{"view", b522DemandViewMake("181", "34", "1", `"request"`, `"payload-length"`), "view"},
		{"view-empty", b522DemandViewMake("181", "34", "1", `"request"`, `""`), "view"},
	}
	for _, b522DemandViewAxis := range b522DemandViewAxes {
		b522DemandViewOut := b522DemandViewRun(b522DemandViewT, b522DemandViewInput(b522DemandViewAxis.ctx, b522DemandViewFields))
		if b522DemandViewOut["state"] != "unknown" || !reflect.DeepEqual(b522DemandViewOut["unknown_reasons"], []any{b522DemandViewAxis.reason}) {
			b522DemandViewT.Fatalf("%s: %v", b522DemandViewAxis.name, b522DemandViewOut)
		}
		if len(b522DemandViewOut["fields"].(map[string]any)) != 0 || b522DemandViewOut["qualified"] != false || b522DemandViewOut["physical"] != "Unknown" {
			b522DemandViewT.Fatalf("%s: guessed selection %v", b522DemandViewAxis.name, b522DemandViewOut)
		}
		if !reflect.DeepEqual(b522DemandViewOut["raw"], b522DemandViewControl["raw"]) {
			b522DemandViewT.Fatalf("%s: raw lost", b522DemandViewAxis.name)
		}
	}
	for _, b522DemandViewOp := range []string{"0", "1", "127", "255"} {
		b522DemandViewOut := b522DemandViewRun(b522DemandViewT, b522DemandViewInput(b522DemandViewMake("181", "34", b522DemandViewOp, `"request"`, `"isolated-fields"`), b522DemandViewFields))
		if b522DemandViewOut["state"] != "partial" {
			b522DemandViewT.Fatalf("operation %s: %v", b522DemandViewOp, b522DemandViewOut["state"])
		}
	}
	b522DemandViewMulti := b522DemandViewRun(b522DemandViewT, b522DemandViewInput(b522DemandViewMake("1", "2", "300", `"x"`, `"y"`), b522DemandViewFields))
	if !reflect.DeepEqual(b522DemandViewMulti["unknown_reasons"], []any{"primary", "secondary", "operation", "direction", "view"}) {
		b522DemandViewT.Fatalf("reasons: %v", b522DemandViewMulti["unknown_reasons"])
	}
	b522DemandViewFam := []string{`"a"`, `"b"`, `null`, `{"deep":[1,2,{"z":true}]}`, `1.50`}
	for _, b522DemandViewValue := range b522DemandViewFam {
		b522DemandViewCtx := `{"primary":181,"secondary":34,"operation":1,"direction":"request","view":"isolated-fields","family":` + b522DemandViewValue + `,"variant":` + b522DemandViewValue + `}`
		b522DemandViewOut := b522DemandViewRun(b522DemandViewT, b522DemandViewInput(b522DemandViewCtx, b522DemandViewFields))
		var b522DemandViewWant map[string]any
		if b522DemandViewErr := json.Unmarshal([]byte(b522DemandViewCtx), &b522DemandViewWant); b522DemandViewErr != nil {
			b522DemandViewT.Fatal(b522DemandViewErr)
		}
		if !reflect.DeepEqual(b522DemandViewOut["context"], b522DemandViewWant) || b522DemandViewOut["qualified"] != false || b522DemandViewOut["physical"] != "Unknown" {
			b522DemandViewT.Fatalf("context %s: %v", b522DemandViewValue, b522DemandViewOut)
		}
	}
}

func TestAb25Ao25InvalidInput(b522DemandViewT *testing.T) {
	b522DemandViewCtx := b522DemandViewContext("request")
	b522DemandViewCases := map[string]string{
		"empty":             ``,
		"malformed":         `{"context":`,
		"not-object":        `[]`,
		"trailing":          b522DemandViewInput(b522DemandViewCtx, `{}`) + ` {}`,
		"dup-top":           `{"context":` + b522DemandViewCtx + `,"fields":{},"fields":{}}`,
		"dup-context":       `{"context":{"primary":181,"primary":181,"secondary":34,"operation":1,"direction":"request","view":"isolated-fields"},"fields":{}}`,
		"dup-field":         b522DemandViewInput(b522DemandViewCtx, `{"mode":[1],"mode":[2]}`),
		"dup-nested-opaque": `{"context":{"primary":181,"secondary":34,"operation":1,"direction":"request","view":"isolated-fields","f":{"a":1,"a":2}},"fields":{}}`,
		"dup-arguments":     `{"context":` + b522DemandViewCtx + `,"fields":{},"arguments":{"a":1,"a":2}}`,
		"element-256":       b522DemandViewInput(b522DemandViewCtx, `{"mode":[256]}`),
		"element-negative":  b522DemandViewInput(b522DemandViewCtx, `{"mode":[-1]}`),
		"element-fraction":  b522DemandViewInput(b522DemandViewCtx, `{"mode":[1.5]}`),
		"element-exponent":  b522DemandViewInput(b522DemandViewCtx, `{"mode":[1e1]}`),
		"element-string":    b522DemandViewInput(b522DemandViewCtx, `{"mode":["1"]}`),
		"element-bool":      b522DemandViewInput(b522DemandViewCtx, `{"mode":[true]}`),
		"element-null":      b522DemandViewInput(b522DemandViewCtx, `{"mode":[null]}`),
		"element-nested":    b522DemandViewInput(b522DemandViewCtx, `{"mode":[[1]]}`),
		"field-scalar":      b522DemandViewInput(b522DemandViewCtx, `{"mode":2}`),
		"field-null":        b522DemandViewInput(b522DemandViewCtx, `{"mode":null}`),
		"field-string":      b522DemandViewInput(b522DemandViewCtx, `{"mode":"02"}`),
		"unmodelled-bad":    b522DemandViewInput(b522DemandViewCtx, `{"tail":[300]}`),
		"fields-array":      b522DemandViewInput(b522DemandViewCtx, `[]`),
		"fields-missing":    `{"context":` + b522DemandViewCtx + `}`,
		"context-missing":   `{"fields":{}}`,
		"context-array":     `{"context":[],"fields":{}}`,
		"primary-string":    b522DemandViewInput(`{"primary":"181","secondary":34,"operation":1,"direction":"request","view":"isolated-fields"}`, `{}`),
		"primary-float":     b522DemandViewInput(`{"primary":181.0,"secondary":34,"operation":1,"direction":"request","view":"isolated-fields"}`, `{}`),
		"primary-missing":   b522DemandViewInput(`{"secondary":34,"operation":1,"direction":"request","view":"isolated-fields"}`, `{}`),
		"secondary-missing": b522DemandViewInput(`{"primary":181,"operation":1,"direction":"request","view":"isolated-fields"}`, `{}`),
		"operation-string":  b522DemandViewInput(`{"primary":181,"secondary":34,"operation":"1","direction":"request","view":"isolated-fields"}`, `{}`),
		"operation-missing": b522DemandViewInput(`{"primary":181,"secondary":34,"direction":"request","view":"isolated-fields"}`, `{}`),
		"direction-number":  b522DemandViewInput(`{"primary":181,"secondary":34,"operation":1,"direction":1,"view":"isolated-fields"}`, `{}`),
		"direction-missing": b522DemandViewInput(`{"primary":181,"secondary":34,"operation":1,"view":"isolated-fields"}`, `{}`),
		"view-missing":      b522DemandViewInput(`{"primary":181,"secondary":34,"operation":1,"direction":"request"}`, `{}`),
		"view-number":       b522DemandViewInput(`{"primary":181,"secondary":34,"operation":1,"direction":"request","view":2}`, `{}`),
		"arguments-null":    `{"context":` + b522DemandViewCtx + `,"fields":{},"arguments":null}`,
		"arguments-array":   `{"context":` + b522DemandViewCtx + `,"fields":{},"arguments":[]}`,
		"arguments-string":  `{"context":` + b522DemandViewCtx + `,"fields":{},"arguments":{"a":"x"}}`,
		"arguments-float":   `{"context":` + b522DemandViewCtx + `,"fields":{},"arguments":{"a":1.5}}`,
		"arguments-nested":  `{"context":` + b522DemandViewCtx + `,"fields":{},"arguments":{"a":{}}}`,
	}
	for b522DemandViewName, b522DemandViewText := range b522DemandViewCases {
		b522DemandViewOut, b522DemandViewErr := B522DemandViewInterpret(json.RawMessage(b522DemandViewText))
		if b522DemandViewErr == nil || b522DemandViewOut != nil {
			b522DemandViewT.Fatalf("%s accepted: %s", b522DemandViewName, b522DemandViewOut)
		}
	}
	b522DemandViewRun(b522DemandViewT, b522DemandViewInput(b522DemandViewCtx, `{"mode":[0],"tail":[255,0,1]}`))
}

func TestAb25Ao25UnrecognisedKeysAndArguments(b522DemandViewT *testing.T) {
	b522DemandViewOut := b522DemandViewRun(b522DemandViewT, `{"context":`+b522DemandViewContext("request")+`,"fields":{"mode":[2,9,9],"tail":[0,255],"extra_empty":[]},"arguments":{"flag":false,"zero":0,"neg":-3},"note":{"a":[1]}}`)
	b522DemandViewRaw := b522DemandViewOut["raw"].(map[string]any)
	if !reflect.DeepEqual(b522DemandViewRaw["tail"], []any{float64(0), float64(255)}) ||
		!reflect.DeepEqual(b522DemandViewRaw["extra_empty"], []any{}) ||
		!reflect.DeepEqual(b522DemandViewRaw["mode"], []any{float64(2), float64(9), float64(9)}) {
		b522DemandViewT.Fatalf("raw: %v", b522DemandViewRaw)
	}
	if !reflect.DeepEqual(b522DemandViewOut["arguments"], map[string]any{"flag": false, "zero": float64(0), "neg": float64(-3)}) {
		b522DemandViewT.Fatalf("arguments: %v", b522DemandViewOut["arguments"])
	}
	if !reflect.DeepEqual(b522DemandViewOut["unrecognised"], map[string]any{"note": map[string]any{"a": []any{float64(1)}}}) {
		b522DemandViewT.Fatalf("unrecognised: %v", b522DemandViewOut["unrecognised"])
	}
	if _, b522DemandViewHas := b522DemandViewFieldOf(b522DemandViewT, b522DemandViewOut, "mode")["number"]; b522DemandViewHas {
		b522DemandViewT.Fatal("overlong mode carried a number")
	}
	b522DemandViewNoArgs := b522DemandViewRequest(b522DemandViewT, `{}`)
	if _, b522DemandViewHas := b522DemandViewNoArgs["arguments"]; b522DemandViewHas {
		b522DemandViewT.Fatal("arguments invented")
	}
	if _, b522DemandViewHas := b522DemandViewNoArgs["unrecognised"]; b522DemandViewHas {
		b522DemandViewT.Fatal("unrecognised invented")
	}
	b522DemandViewEmptyArgs := b522DemandViewRun(b522DemandViewT, `{"context":`+b522DemandViewContext("request")+`,"fields":{},"arguments":{}}`)
	if !reflect.DeepEqual(b522DemandViewEmptyArgs["arguments"], map[string]any{}) {
		b522DemandViewT.Fatalf("empty arguments not preserved: %v", b522DemandViewEmptyArgs["arguments"])
	}
}

func TestAb25Ao25IndependenceAndRoundTrip(b522DemandViewT *testing.T) {
	b522DemandViewText := `{"context":` + b522DemandViewContext("request") + `,"fields":{"mode":[2],"fan":[],"tail":[1]},"arguments":{"a":1}}`
	b522DemandViewInput := json.RawMessage(b522DemandViewText)
	b522DemandViewFirst, b522DemandViewErr := B522DemandViewInterpret(b522DemandViewInput)
	if b522DemandViewErr != nil {
		b522DemandViewT.Fatal(b522DemandViewErr)
	}
	b522DemandViewSnapshot := string(b522DemandViewFirst)
	for b522DemandViewIndex := range b522DemandViewInput {
		b522DemandViewInput[b522DemandViewIndex] = ' '
	}
	for b522DemandViewIndex := range b522DemandViewFirst {
		b522DemandViewFirst[b522DemandViewIndex] = 'x'
	}
	b522DemandViewSecond, b522DemandViewErr := B522DemandViewInterpret(json.RawMessage(b522DemandViewText))
	if b522DemandViewErr != nil || string(b522DemandViewSecond) != b522DemandViewSnapshot {
		b522DemandViewT.Fatalf("later call changed: %v\n%s\n%s", b522DemandViewErr, b522DemandViewSecond, b522DemandViewSnapshot)
	}
	var b522DemandViewDecoded map[string]any
	if b522DemandViewErr = json.Unmarshal(b522DemandViewSecond, &b522DemandViewDecoded); b522DemandViewErr != nil {
		b522DemandViewT.Fatal(b522DemandViewErr)
	}
	b522DemandViewDecoded["raw"].(map[string]any)["mode"].([]any)[0] = float64(99)
	b522DemandViewDecoded["context"].(map[string]any)["primary"] = float64(1)
	b522DemandViewDecoded["fields"].(map[string]any)["mode"].(map[string]any)["raw"].([]any)[0] = float64(98)
	b522DemandViewThird, _ := B522DemandViewInterpret(json.RawMessage(b522DemandViewText))
	if string(b522DemandViewThird) != b522DemandViewSnapshot {
		b522DemandViewT.Fatal("mutating a decoded result altered a later call")
	}
	var b522DemandViewFresh map[string]any
	_ = json.Unmarshal(b522DemandViewThird, &b522DemandViewFresh)
	if b522DemandViewFresh["raw"].(map[string]any)["mode"].([]any)[0] != float64(2) ||
		b522DemandViewFresh["fields"].(map[string]any)["mode"].(map[string]any)["raw"].([]any)[0] != float64(2) {
		b522DemandViewT.Fatal("raw and field raw are not independent copies")
	}
	b522DemandViewRound, _ := json.Marshal(b522DemandViewFresh)
	var b522DemandViewAgain map[string]any
	_ = json.Unmarshal(b522DemandViewRound, &b522DemandViewAgain)
	if !reflect.DeepEqual(b522DemandViewFresh, b522DemandViewAgain) {
		b522DemandViewT.Fatal("round trip changed result")
	}
	b522DemandViewF := b522DemandViewFresh["fields"].(map[string]any)
	if b522DemandViewF["fan"].(map[string]any)["raw"] == nil || b522DemandViewF["room_temperature"].(map[string]any)["raw"] != nil {
		b522DemandViewT.Fatal("empty and absent raw not distinct after round trip")
	}
	if b522DemandViewFresh["state"] != "partial" || b522DemandViewFresh["qualified"] != false || b522DemandViewFresh["physical"] != "Unknown" {
		b522DemandViewT.Fatalf("top-level: %v", b522DemandViewFresh)
	}
}
