package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

type b512HeatpumpIOEvaluatorCase struct {
	B512HeatpumpIOEvaluatorID       string          `json:"id"`
	B512HeatpumpIOEvaluatorInput    json.RawMessage `json:"input"`
	B512HeatpumpIOEvaluatorExpected map[string]any  `json:"expected"`
}

func b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT *testing.T) []b512HeatpumpIOEvaluatorCase {
	b512HeatpumpIOEvaluatorT.Helper()
	b512HeatpumpIOEvaluatorBytes, b512HeatpumpIOEvaluatorErr := os.ReadFile("testdata/ab13ap13_b512_heatpump_io.json")
	if b512HeatpumpIOEvaluatorErr != nil {
		b512HeatpumpIOEvaluatorT.Fatal(b512HeatpumpIOEvaluatorErr)
	}
	var b512HeatpumpIOEvaluatorCases []b512HeatpumpIOEvaluatorCase
	if b512HeatpumpIOEvaluatorErr := json.Unmarshal(b512HeatpumpIOEvaluatorBytes, &b512HeatpumpIOEvaluatorCases); b512HeatpumpIOEvaluatorErr != nil {
		b512HeatpumpIOEvaluatorT.Fatal(b512HeatpumpIOEvaluatorErr)
	}
	if len(b512HeatpumpIOEvaluatorCases) != 4 {
		b512HeatpumpIOEvaluatorT.Fatal("four synthetic views required")
	}
	return b512HeatpumpIOEvaluatorCases
}

func b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT *testing.T, b512HeatpumpIOEvaluatorBytes []byte) map[string]any {
	b512HeatpumpIOEvaluatorT.Helper()
	var b512HeatpumpIOEvaluatorObject map[string]any
	if b512HeatpumpIOEvaluatorErr := json.Unmarshal(b512HeatpumpIOEvaluatorBytes, &b512HeatpumpIOEvaluatorObject); b512HeatpumpIOEvaluatorErr != nil {
		b512HeatpumpIOEvaluatorT.Fatal(b512HeatpumpIOEvaluatorErr)
	}
	return b512HeatpumpIOEvaluatorObject
}

func b512HeatpumpIOEvaluatorEncode(b512HeatpumpIOEvaluatorT *testing.T, b512HeatpumpIOEvaluatorValue any) []byte {
	b512HeatpumpIOEvaluatorT.Helper()
	b512HeatpumpIOEvaluatorBytes, b512HeatpumpIOEvaluatorErr := json.Marshal(b512HeatpumpIOEvaluatorValue)
	if b512HeatpumpIOEvaluatorErr != nil {
		b512HeatpumpIOEvaluatorT.Fatal(b512HeatpumpIOEvaluatorErr)
	}
	return b512HeatpumpIOEvaluatorBytes
}

func b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT *testing.T, b512HeatpumpIOEvaluatorInput any) map[string]any {
	b512HeatpumpIOEvaluatorT.Helper()
	b512HeatpumpIOEvaluatorBytes, b512HeatpumpIOEvaluatorErr := B512HeatpumpIOEvaluate(b512HeatpumpIOEvaluatorEncode(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput))
	if b512HeatpumpIOEvaluatorErr != nil {
		b512HeatpumpIOEvaluatorT.Fatal(b512HeatpumpIOEvaluatorErr)
	}
	b512HeatpumpIOEvaluatorOutput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorBytes)
	if b512HeatpumpIOEvaluatorOutput["physical"] != "Unknown" || b512HeatpumpIOEvaluatorOutput["qualified"] != false {
		b512HeatpumpIOEvaluatorT.Fatal("local interpretation qualified physical support")
	}
	return b512HeatpumpIOEvaluatorOutput
}

func b512HeatpumpIOEvaluatorField(b512HeatpumpIOEvaluatorOutput map[string]any, b512HeatpumpIOEvaluatorName string) map[string]any {
	return b512HeatpumpIOEvaluatorOutput["fields"].(map[string]any)[b512HeatpumpIOEvaluatorName].(map[string]any)
}

func b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT *testing.T, b512HeatpumpIOEvaluatorActual, b512HeatpumpIOEvaluatorExpected any) {
	b512HeatpumpIOEvaluatorT.Helper()
	if !reflect.DeepEqual(b512HeatpumpIOEvaluatorActual, b512HeatpumpIOEvaluatorExpected) {
		b512HeatpumpIOEvaluatorT.Fatalf("got %#v, want %#v", b512HeatpumpIOEvaluatorActual, b512HeatpumpIOEvaluatorExpected)
	}
}

func b512HeatpumpIOEvaluatorAbsent(b512HeatpumpIOEvaluatorT *testing.T, b512HeatpumpIOEvaluatorObject map[string]any, b512HeatpumpIOEvaluatorKeys ...string) {
	b512HeatpumpIOEvaluatorT.Helper()
	for _, b512HeatpumpIOEvaluatorKey := range b512HeatpumpIOEvaluatorKeys {
		if _, b512HeatpumpIOEvaluatorExists := b512HeatpumpIOEvaluatorObject[b512HeatpumpIOEvaluatorKey]; b512HeatpumpIOEvaluatorExists {
			b512HeatpumpIOEvaluatorT.Fatalf("unexpected %s: %#v", b512HeatpumpIOEvaluatorKey, b512HeatpumpIOEvaluatorObject)
		}
	}
}

func TestAb13Ap13Fixtures(b512HeatpumpIOEvaluatorT *testing.T) {
	for _, b512HeatpumpIOEvaluatorCase := range b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT) {
		b512HeatpumpIOEvaluatorT.Run(b512HeatpumpIOEvaluatorCase.B512HeatpumpIOEvaluatorID, func(b512HeatpumpIOEvaluatorT *testing.T) {
			b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCase.B512HeatpumpIOEvaluatorInput)
			b512HeatpumpIOEvaluatorOutput := b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput)
			b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["context"], b512HeatpumpIOEvaluatorInput["context"])
			b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["raw"], b512HeatpumpIOEvaluatorInput["fields"])
			b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["state"], "complete")
			for b512HeatpumpIOEvaluatorName, b512HeatpumpIOEvaluatorExpected := range b512HeatpumpIOEvaluatorCase.B512HeatpumpIOEvaluatorExpected {
				for b512HeatpumpIOEvaluatorKey, b512HeatpumpIOEvaluatorValue := range b512HeatpumpIOEvaluatorExpected.(map[string]any) {
					b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorField(b512HeatpumpIOEvaluatorOutput, b512HeatpumpIOEvaluatorName)[b512HeatpumpIOEvaluatorKey], b512HeatpumpIOEvaluatorValue)
				}
			}
		})
	}
}

func TestAb13Ap13Temperatures(b512HeatpumpIOEvaluatorT *testing.T) {
	for _, b512HeatpumpIOEvaluatorName := range []string{"temperature1", "temperature2"} {
		for b512HeatpumpIOEvaluatorIndex, b512HeatpumpIOEvaluatorBytes := range [][]int{{0, 128}, {240, 255}, {0, 0}, {16, 0}, {255, 127}, {1, 128}} {
			b512HeatpumpIOEvaluatorT.Run(fmt.Sprintf("ab13ap13_%s_%d", b512HeatpumpIOEvaluatorName, b512HeatpumpIOEvaluatorIndex), func(b512HeatpumpIOEvaluatorT *testing.T) {
				b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT)[1].B512HeatpumpIOEvaluatorInput)
				b512HeatpumpIOEvaluatorInput["fields"].(map[string]any)[b512HeatpumpIOEvaluatorName] = b512HeatpumpIOEvaluatorBytes
				b512HeatpumpIOEvaluatorResult := b512HeatpumpIOEvaluatorField(b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput), b512HeatpumpIOEvaluatorName)
				b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["state"], "Hypothesis")
				if b512HeatpumpIOEvaluatorIndex == 0 {
					b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["label"], "replacement")
					b512HeatpumpIOEvaluatorAbsent(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult, "number", "calculation")
				} else {
					b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["number"], []float64{0, -1, 0, 1, 2047.9375, -2047.9375}[b512HeatpumpIOEvaluatorIndex])
					b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["calculation"], "Proven")
					b512HeatpumpIOEvaluatorAbsent(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult, "label")
				}
			})
		}
	}
}

func TestAb13Ap13CodesAndPressure(b512HeatpumpIOEvaluatorT *testing.T) {
	for _, b512HeatpumpIOEvaluatorRow := range []struct {
		b512HeatpumpIOEvaluatorView  int
		b512HeatpumpIOEvaluatorName  string
		b512HeatpumpIOEvaluatorByte  int
		b512HeatpumpIOEvaluatorLabel string
	}{
		{0, "valve", 0, "off"}, {0, "valve", 1, "CH"}, {0, "valve", 2, "DHW"}, {0, "valve", 3, ""}, {0, "valve", 255, ""},
		{1, "input", 0, "multipurpose"}, {1, "input", 1, "utility-contact"}, {1, "input", 2, ""},
		{1, "backup", 0, "off"}, {1, "backup", 1, "step"}, {1, "backup", 2, "step"}, {1, "backup", 3, "step"}, {1, "backup", 4, ""}, {1, "backup", 254, ""}, {1, "backup", 255, "replacement"},
		{3, "board", 0, "opaque-code-0"}, {3, "board", 1, "opaque-code-1"}, {3, "board", 2, "opaque-code-2"}, {3, "board", 3, "opaque-code-3"}, {3, "board", 11, "opaque-code-11"}, {3, "board", 12, "opaque-code-12"}, {3, "board", 13, ""}, {3, "board", 255, ""},
	} {
		b512HeatpumpIOEvaluatorT.Run(fmt.Sprintf("ab13ap13_%s_%d", b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorName, b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorByte), func(b512HeatpumpIOEvaluatorT *testing.T) {
			b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT)[b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorView].B512HeatpumpIOEvaluatorInput)
			b512HeatpumpIOEvaluatorInput["fields"].(map[string]any)[b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorName] = []int{b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorByte}
			b512HeatpumpIOEvaluatorResult := b512HeatpumpIOEvaluatorField(b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput), b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorName)
			if b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorLabel == "" {
				b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["state"], "Unknown")
				b512HeatpumpIOEvaluatorAbsent(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult, "label", "number")
			} else {
				b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["label"], b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorLabel)
				if b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorLabel == "step" {
					b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["number"], float64(b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorByte))
				} else {
					b512HeatpumpIOEvaluatorAbsent(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult, "number")
				}
			}
		})
	}
	for _, b512HeatpumpIOEvaluatorRow := range []struct {
		b512HeatpumpIOEvaluatorByte   int
		b512HeatpumpIOEvaluatorNumber float64
	}{{0, 0}, {1, 0.1}, {10, 1}, {254, 25.4}, {255, 25.5}} {
		b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT)[1].B512HeatpumpIOEvaluatorInput)
		b512HeatpumpIOEvaluatorInput["fields"].(map[string]any)["pressure"] = []int{b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorByte}
		b512HeatpumpIOEvaluatorResult := b512HeatpumpIOEvaluatorField(b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput), "pressure")
		b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["number"], b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorNumber)
		b512HeatpumpIOEvaluatorDefault := b512HeatpumpIOEvaluatorResult["alternatives"].([]any)[0].(map[string]any)
		b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorDefault["number"], float64(255))
		b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorDefault["candidate_domain"], false)
		b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorDefault["conflict"], true)
	}
}

func TestAb13Ap13EveryBit(b512HeatpumpIOEvaluatorT *testing.T) {
	for _, b512HeatpumpIOEvaluatorRow := range []struct {
		b512HeatpumpIOEvaluatorView   int
		b512HeatpumpIOEvaluatorName   string
		b512HeatpumpIOEvaluatorLabels []string
	}{
		{0, "output", []string{"output1", "output2"}},
		{2, "actuator1", []string{"cooling-valve", "passive-mixer1", "passive-mixer2", "cooling-active", "well-pump", "backup-safety", "backup1", "backup2"}},
		{2, "actuator2", []string{"backup3", "external-valve", "external-backup"}},
		{3, "status", []string{"anode1", "anode2", "smart-grid", "well", "safety-limiter", "safety"}},
	} {
		for b512HeatpumpIOEvaluatorBit := -1; b512HeatpumpIOEvaluatorBit < 8; b512HeatpumpIOEvaluatorBit++ {
			b512HeatpumpIOEvaluatorT.Run(fmt.Sprintf("ab13ap13_%s_bit%d", b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorName, b512HeatpumpIOEvaluatorBit), func(b512HeatpumpIOEvaluatorT *testing.T) {
				b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT)[b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorView].B512HeatpumpIOEvaluatorInput)
				b512HeatpumpIOEvaluatorByte := 0
				if b512HeatpumpIOEvaluatorBit >= 0 {
					b512HeatpumpIOEvaluatorByte = 1 << b512HeatpumpIOEvaluatorBit
				}
				b512HeatpumpIOEvaluatorInput["fields"].(map[string]any)[b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorName] = []int{b512HeatpumpIOEvaluatorByte}
				b512HeatpumpIOEvaluatorResult := b512HeatpumpIOEvaluatorField(b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput), b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorName)
				b512HeatpumpIOEvaluatorFlags := b512HeatpumpIOEvaluatorResult["flags"].(map[string]any)
				b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, len(b512HeatpumpIOEvaluatorFlags), len(b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorLabels))
				for b512HeatpumpIOEvaluatorIndex, b512HeatpumpIOEvaluatorLabel := range b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorLabels {
					b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorFlags[b512HeatpumpIOEvaluatorLabel], b512HeatpumpIOEvaluatorIndex == b512HeatpumpIOEvaluatorBit)
				}
				b512HeatpumpIOEvaluatorReserved := 0
				if b512HeatpumpIOEvaluatorBit >= len(b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorLabels) {
					b512HeatpumpIOEvaluatorReserved = b512HeatpumpIOEvaluatorByte
				}
				b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["reserved"], float64(b512HeatpumpIOEvaluatorReserved))
			})
		}
	}
}

func TestAb13Ap13ADCChannels(b512HeatpumpIOEvaluatorT *testing.T) {
	for _, b512HeatpumpIOEvaluatorName := range []string{"adc0", "adc1", "adc2", "adc3", "adc4", "adc5"} {
		for _, b512HeatpumpIOEvaluatorRow := range []struct {
			b512HeatpumpIOEvaluatorBytes  []int
			b512HeatpumpIOEvaluatorNumber float64
		}{{[]int{0, 0}, 0}, {[]int{1, 0}, 1}, {[]int{0, 1}, 256}, {[]int{18, 52}, 13330}, {[]int{255, 255}, 65535}} {
			b512HeatpumpIOEvaluatorT.Run(fmt.Sprintf("ab13ap13_%s_%v", b512HeatpumpIOEvaluatorName, b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorBytes), func(b512HeatpumpIOEvaluatorT *testing.T) {
				b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT)[3].B512HeatpumpIOEvaluatorInput)
				b512HeatpumpIOEvaluatorControl := b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput)
				b512HeatpumpIOEvaluatorInput["fields"].(map[string]any)[b512HeatpumpIOEvaluatorName] = b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorBytes
				b512HeatpumpIOEvaluatorOutput := b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput)
				b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorField(b512HeatpumpIOEvaluatorOutput, b512HeatpumpIOEvaluatorName)["number"], b512HeatpumpIOEvaluatorRow.b512HeatpumpIOEvaluatorNumber)
				for b512HeatpumpIOEvaluatorSibling, b512HeatpumpIOEvaluatorResult := range b512HeatpumpIOEvaluatorControl["fields"].(map[string]any) {
					if b512HeatpumpIOEvaluatorSibling != b512HeatpumpIOEvaluatorName {
						b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["fields"].(map[string]any)[b512HeatpumpIOEvaluatorSibling], b512HeatpumpIOEvaluatorResult)
					}
				}
			})
		}
	}
}

func TestAb13Ap13OpaqueInverseAndPhase(b512HeatpumpIOEvaluatorT *testing.T) {
	for _, b512HeatpumpIOEvaluatorBytes := range [][]int{{18, 52}, {52, 18}, {0, 0}, {255, 255}} {
		b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT)[2].B512HeatpumpIOEvaluatorInput)
		b512HeatpumpIOEvaluatorInput["fields"].(map[string]any)["inverse"] = b512HeatpumpIOEvaluatorBytes
		b512HeatpumpIOEvaluatorResult := b512HeatpumpIOEvaluatorField(b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput), "inverse")
		b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["raw"], b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorEncode(b512HeatpumpIOEvaluatorT, map[string]any{"bytes": b512HeatpumpIOEvaluatorBytes}))["bytes"])
		b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["state"], "Unknown")
		b512HeatpumpIOEvaluatorAbsent(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult, "number", "label", "flags")
		b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["unmodeled"], map[string]any{"order": "Unknown", "algorithm": "Unknown"})
	}
	for _, b512HeatpumpIOEvaluatorName := range []string{"phase1", "phase2"} {
		for _, b512HeatpumpIOEvaluatorByte := range []int{0, 33, 50, 66, 100, 255} {
			b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT)[3].B512HeatpumpIOEvaluatorInput)
			b512HeatpumpIOEvaluatorInput["fields"].(map[string]any)[b512HeatpumpIOEvaluatorName] = []int{b512HeatpumpIOEvaluatorByte}
			b512HeatpumpIOEvaluatorResult := b512HeatpumpIOEvaluatorField(b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput), b512HeatpumpIOEvaluatorName)
			b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["number"], float64(b512HeatpumpIOEvaluatorByte))
			b512HeatpumpIOEvaluatorAbsent(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult, "label", "phase_kind", "selected", "tolerance")
			b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["alternatives"], []any{map[string]any{"state": "Hypothesis", "label": "single-phase", "approximate": []any{float64(50)}}, map[string]any{"state": "Hypothesis", "label": "three-phase", "approximate": []any{float64(0), float64(33), float64(66), float64(100)}}})
		}
	}
}

func TestAb13Ap13ContextAxes(b512HeatpumpIOEvaluatorT *testing.T) {
	for _, b512HeatpumpIOEvaluatorCase := range b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT) {
		for _, b512HeatpumpIOEvaluatorAxis := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
			for b512HeatpumpIOEvaluatorIndex, b512HeatpumpIOEvaluatorValue := range []any{nil, float64(0), "other", false, map[string]any{"opaque": []any{nil, false, float64(0)}}} {
				if b512HeatpumpIOEvaluatorAxis == "primary" || b512HeatpumpIOEvaluatorAxis == "secondary" || b512HeatpumpIOEvaluatorAxis == "operation" {
					if b512HeatpumpIOEvaluatorIndex > 1 {
						continue
					}
				}
				if b512HeatpumpIOEvaluatorAxis == "direction" || b512HeatpumpIOEvaluatorAxis == "view" {
					if b512HeatpumpIOEvaluatorIndex == 1 || b512HeatpumpIOEvaluatorIndex > 2 {
						continue
					}
				}
				b512HeatpumpIOEvaluatorT.Run(fmt.Sprintf("%s_%s_%d", b512HeatpumpIOEvaluatorCase.B512HeatpumpIOEvaluatorID, b512HeatpumpIOEvaluatorAxis, b512HeatpumpIOEvaluatorIndex), func(b512HeatpumpIOEvaluatorT *testing.T) {
					b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCase.B512HeatpumpIOEvaluatorInput)
					b512HeatpumpIOEvaluatorContext := b512HeatpumpIOEvaluatorInput["context"].(map[string]any)
					b512HeatpumpIOEvaluatorContext[b512HeatpumpIOEvaluatorAxis] = b512HeatpumpIOEvaluatorValue
					b512HeatpumpIOEvaluatorOutput := b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput)
					b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["context"], b512HeatpumpIOEvaluatorContext)
					b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["raw"], b512HeatpumpIOEvaluatorInput["fields"])
					if b512HeatpumpIOEvaluatorAxis == "family" || b512HeatpumpIOEvaluatorAxis == "variant" {
						b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["state"], "complete")
					} else {
						b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["state"], "Unknown")
						b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, len(b512HeatpumpIOEvaluatorOutput["fields"].(map[string]any)), 0)
					}
				})
			}
			b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCase.B512HeatpumpIOEvaluatorInput)
			delete(b512HeatpumpIOEvaluatorInput["context"].(map[string]any), b512HeatpumpIOEvaluatorAxis)
			b512HeatpumpIOEvaluatorOutput := b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput)
			b512HeatpumpIOEvaluatorAbsent(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["context"].(map[string]any), b512HeatpumpIOEvaluatorAxis)
			if b512HeatpumpIOEvaluatorAxis != "family" && b512HeatpumpIOEvaluatorAxis != "variant" {
				b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["state"], "Unknown")
			}
		}
	}
	for _, b512HeatpumpIOEvaluatorPair := range []struct {
		b512HeatpumpIOEvaluatorOperation int
		b512HeatpumpIOEvaluatorView      string
	}{{15, "board-io"}, {16, "valve-io"}, {17, "valve-io"}} {
		b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT)[0].B512HeatpumpIOEvaluatorInput)
		b512HeatpumpIOEvaluatorInput["context"].(map[string]any)["operation"] = b512HeatpumpIOEvaluatorPair.b512HeatpumpIOEvaluatorOperation
		b512HeatpumpIOEvaluatorInput["context"].(map[string]any)["view"] = b512HeatpumpIOEvaluatorPair.b512HeatpumpIOEvaluatorView
		b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput)["state"], "Unknown")
	}
}

func TestAb13Ap13EveryFieldPartial(b512HeatpumpIOEvaluatorT *testing.T) {
	for _, b512HeatpumpIOEvaluatorCase := range b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT) {
		b512HeatpumpIOEvaluatorOriginal := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCase.B512HeatpumpIOEvaluatorInput)
		for b512HeatpumpIOEvaluatorName, b512HeatpumpIOEvaluatorOriginalBytes := range b512HeatpumpIOEvaluatorOriginal["fields"].(map[string]any) {
			b512HeatpumpIOEvaluatorWidth := len(b512HeatpumpIOEvaluatorOriginalBytes.([]any))
			for b512HeatpumpIOEvaluatorLength := -1; b512HeatpumpIOEvaluatorLength <= b512HeatpumpIOEvaluatorWidth+1; b512HeatpumpIOEvaluatorLength++ {
				b512HeatpumpIOEvaluatorT.Run(fmt.Sprintf("%s_%s_length%d", b512HeatpumpIOEvaluatorCase.B512HeatpumpIOEvaluatorID, b512HeatpumpIOEvaluatorName, b512HeatpumpIOEvaluatorLength), func(b512HeatpumpIOEvaluatorT *testing.T) {
					b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCase.B512HeatpumpIOEvaluatorInput)
					b512HeatpumpIOEvaluatorControl := b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput)
					b512HeatpumpIOEvaluatorRaw := b512HeatpumpIOEvaluatorInput["fields"].(map[string]any)
					if b512HeatpumpIOEvaluatorLength < 0 {
						delete(b512HeatpumpIOEvaluatorRaw, b512HeatpumpIOEvaluatorName)
					} else {
						b512HeatpumpIOEvaluatorRaw[b512HeatpumpIOEvaluatorName] = make([]int, b512HeatpumpIOEvaluatorLength)
					}
					b512HeatpumpIOEvaluatorOutput := b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput)
					b512HeatpumpIOEvaluatorResult := b512HeatpumpIOEvaluatorField(b512HeatpumpIOEvaluatorOutput, b512HeatpumpIOEvaluatorName)
					if b512HeatpumpIOEvaluatorLength != b512HeatpumpIOEvaluatorWidth {
						b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["state"], "partial")
						b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult["state"], "incomplete")
						b512HeatpumpIOEvaluatorAbsent(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult, "number", "label", "flags", "calculation")
					} else {
						b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["state"], "complete")
					}
					if b512HeatpumpIOEvaluatorLength < 0 {
						b512HeatpumpIOEvaluatorAbsent(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorResult, "raw")
					} else {
						b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, len(b512HeatpumpIOEvaluatorResult["raw"].([]any)), b512HeatpumpIOEvaluatorLength)
					}
					b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["raw"], b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorEncode(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorRaw)))
					for b512HeatpumpIOEvaluatorSibling, b512HeatpumpIOEvaluatorResult := range b512HeatpumpIOEvaluatorControl["fields"].(map[string]any) {
						if b512HeatpumpIOEvaluatorSibling != b512HeatpumpIOEvaluatorName {
							b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["fields"].(map[string]any)[b512HeatpumpIOEvaluatorSibling], b512HeatpumpIOEvaluatorResult)
						}
					}
				})
			}
		}
	}
}

func TestAb13Ap13InvalidJSONAndKinds(b512HeatpumpIOEvaluatorT *testing.T) {
	for b512HeatpumpIOEvaluatorIndex, b512HeatpumpIOEvaluatorInput := range []string{
		``, `{`, `[]`, `null`, `{} {}`, `{"context":{},"fields":{},}`, `{"context":{},"context":{},"fields":{}}`,
		`{"context":{"family":{"x":1,"x":2}},"fields":{}}`, `{"context":{},"fields":{"x":[1],"x":[2]}}`,
		`{"context":{},"fields":{},"arguments":{"extra":[{"x":false,"x":true}]}}`,
		`{"context":null,"fields":{}}`, `{"context":{},"fields":null}`, `{"context":{},"fields":{},"arguments":null}`,
	} {
		b512HeatpumpIOEvaluatorT.Run(fmt.Sprintf("ab13ap13_invalid_%d", b512HeatpumpIOEvaluatorIndex), func(b512HeatpumpIOEvaluatorT *testing.T) {
			b512HeatpumpIOEvaluatorResult, b512HeatpumpIOEvaluatorErr := B512HeatpumpIOEvaluate(json.RawMessage(b512HeatpumpIOEvaluatorInput))
			if b512HeatpumpIOEvaluatorErr == nil || b512HeatpumpIOEvaluatorResult != nil {
				b512HeatpumpIOEvaluatorT.Fatalf("invalid input succeeded: %s", b512HeatpumpIOEvaluatorInput)
			}
		})
	}
	for _, b512HeatpumpIOEvaluatorCase := range b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT) {
		b512HeatpumpIOEvaluatorOriginal := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCase.B512HeatpumpIOEvaluatorInput)
		for b512HeatpumpIOEvaluatorName := range b512HeatpumpIOEvaluatorOriginal["fields"].(map[string]any) {
			for b512HeatpumpIOEvaluatorIndex, b512HeatpumpIOEvaluatorValue := range []any{nil, false, "0", 0, []any{false}, []any{nil}, []any{"1"}, []any{-1}, []any{256}, []any{1.5}, []any{4294967296.0}, []any{[]any{0}}} {
				b512HeatpumpIOEvaluatorT.Run(fmt.Sprintf("%s_%s_invalid%d", b512HeatpumpIOEvaluatorCase.B512HeatpumpIOEvaluatorID, b512HeatpumpIOEvaluatorName, b512HeatpumpIOEvaluatorIndex), func(b512HeatpumpIOEvaluatorT *testing.T) {
					b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCase.B512HeatpumpIOEvaluatorInput)
					b512HeatpumpIOEvaluatorInput["fields"].(map[string]any)[b512HeatpumpIOEvaluatorName] = b512HeatpumpIOEvaluatorValue
					b512HeatpumpIOEvaluatorResult, b512HeatpumpIOEvaluatorErr := B512HeatpumpIOEvaluate(b512HeatpumpIOEvaluatorEncode(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput))
					if b512HeatpumpIOEvaluatorErr == nil || b512HeatpumpIOEvaluatorResult != nil {
						b512HeatpumpIOEvaluatorT.Fatalf("invalid byte kind succeeded: %#v", b512HeatpumpIOEvaluatorValue)
					}
				})
			}
		}
	}
	for _, b512HeatpumpIOEvaluatorName := range []string{"primary", "secondary", "operation", "direction", "view"} {
		b512HeatpumpIOEvaluatorValues := []any{false, []any{}, map[string]any{}, 1.5, -1, 256, "181"}
		if b512HeatpumpIOEvaluatorName == "direction" || b512HeatpumpIOEvaluatorName == "view" {
			b512HeatpumpIOEvaluatorValues = []any{false, 0, []any{}, map[string]any{}}
		}
		for _, b512HeatpumpIOEvaluatorValue := range b512HeatpumpIOEvaluatorValues {
			b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT)[0].B512HeatpumpIOEvaluatorInput)
			b512HeatpumpIOEvaluatorInput["context"].(map[string]any)[b512HeatpumpIOEvaluatorName] = b512HeatpumpIOEvaluatorValue
			b512HeatpumpIOEvaluatorResult, b512HeatpumpIOEvaluatorErr := B512HeatpumpIOEvaluate(b512HeatpumpIOEvaluatorEncode(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput))
			if b512HeatpumpIOEvaluatorErr == nil || b512HeatpumpIOEvaluatorResult != nil {
				b512HeatpumpIOEvaluatorT.Fatalf("invalid context %s succeeded", b512HeatpumpIOEvaluatorName)
			}
		}
	}
	for _, b512HeatpumpIOEvaluatorValue := range []any{nil, false, "", 0, []any{}} {
		b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT)[0].B512HeatpumpIOEvaluatorInput)
		b512HeatpumpIOEvaluatorInput["arguments"] = b512HeatpumpIOEvaluatorValue
		b512HeatpumpIOEvaluatorResult, b512HeatpumpIOEvaluatorErr := B512HeatpumpIOEvaluate(b512HeatpumpIOEvaluatorEncode(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput))
		if b512HeatpumpIOEvaluatorErr == nil || b512HeatpumpIOEvaluatorResult != nil {
			b512HeatpumpIOEvaluatorT.Fatal("invalid arguments object succeeded")
		}
	}
	for _, b512HeatpumpIOEvaluatorNumber := range []string{"1.0", "1e0", "-0", "256", "1e999"} {
		b512HeatpumpIOEvaluatorResult, b512HeatpumpIOEvaluatorErr := B512HeatpumpIOEvaluate(json.RawMessage(`{"context":{},"fields":{"extra":[` + b512HeatpumpIOEvaluatorNumber + `]}}`))
		if b512HeatpumpIOEvaluatorErr == nil || b512HeatpumpIOEvaluatorResult != nil {
			b512HeatpumpIOEvaluatorT.Fatal("invalid byte representation succeeded")
		}
	}
}

func TestAb13Ap13ArgumentNumbers(b512HeatpumpIOEvaluatorT *testing.T) {
	for _, b512HeatpumpIOEvaluatorNumber := range []string{"-1", "0.5", "1.0", "1e0", "4294967296", "1e999"} {
		b512HeatpumpIOEvaluatorInput := json.RawMessage(`{"context":{},"fields":{},"arguments":{"extra":[{"nested":` + b512HeatpumpIOEvaluatorNumber + `}]}}`)
		b512HeatpumpIOEvaluatorResult, b512HeatpumpIOEvaluatorErr := B512HeatpumpIOEvaluate(b512HeatpumpIOEvaluatorInput)
		if b512HeatpumpIOEvaluatorErr == nil || b512HeatpumpIOEvaluatorResult != nil {
			b512HeatpumpIOEvaluatorT.Fatalf("invalid argument number succeeded: %s", b512HeatpumpIOEvaluatorNumber)
		}
	}
}

func TestAb13Ap13ArgumentsExtrasAndStorage(b512HeatpumpIOEvaluatorT *testing.T) {
	for _, b512HeatpumpIOEvaluatorCase := range b512HeatpumpIOEvaluatorCases(b512HeatpumpIOEvaluatorT) {
		for _, b512HeatpumpIOEvaluatorValue := range []any{nil, false, true, "", float64(0), float64(4294967295), []any{}, map[string]any{}, []any{nil, false, map[string]any{"nested": []any{float64(0), "opaque"}}}} {
			b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCase.B512HeatpumpIOEvaluatorInput)
			b512HeatpumpIOEvaluatorInput["arguments"] = map[string]any{"extra": b512HeatpumpIOEvaluatorValue}
			b512HeatpumpIOEvaluatorInput["context"].(map[string]any)["opaque"] = b512HeatpumpIOEvaluatorValue
			b512HeatpumpIOEvaluatorInput["fields"].(map[string]any)["extra"] = []int{0, 255}
			b512HeatpumpIOEvaluatorOutput := b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput)
			b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["arguments"], b512HeatpumpIOEvaluatorInput["arguments"])
			b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["context"], b512HeatpumpIOEvaluatorInput["context"])
			b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["state"], "complete")
			b512HeatpumpIOEvaluatorAbsent(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorOutput["fields"].(map[string]any), "extra")
		}
		b512HeatpumpIOEvaluatorInput := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorCase.B512HeatpumpIOEvaluatorInput)
		b512HeatpumpIOEvaluatorAbsent(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput), "arguments")
		b512HeatpumpIOEvaluatorInput["arguments"] = map[string]any{}
		b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorEvaluate(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput)["arguments"], map[string]any{})
	}
	b512HeatpumpIOEvaluatorInput := []byte(`{"context":{"primary":181,"secondary":18,"operation":16,"direction":"request","view":"board-io","family":{"nested":[false,0]},"variant":null},"fields":{"actuator1":[0],"actuator2":[0],"inverse":[18,52]},"arguments":{"extra":{"nested":[false,0]}}}`)
	b512HeatpumpIOEvaluatorOriginal := append([]byte(nil), b512HeatpumpIOEvaluatorInput...)
	b512HeatpumpIOEvaluatorFirst, b512HeatpumpIOEvaluatorErr := B512HeatpumpIOEvaluate(b512HeatpumpIOEvaluatorInput)
	if b512HeatpumpIOEvaluatorErr != nil {
		b512HeatpumpIOEvaluatorT.Fatal(b512HeatpumpIOEvaluatorErr)
	}
	b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorInput, b512HeatpumpIOEvaluatorOriginal)
	b512HeatpumpIOEvaluatorSecond, b512HeatpumpIOEvaluatorErr := B512HeatpumpIOEvaluate(b512HeatpumpIOEvaluatorInput)
	if b512HeatpumpIOEvaluatorErr != nil {
		b512HeatpumpIOEvaluatorT.Fatal(b512HeatpumpIOEvaluatorErr)
	}
	b512HeatpumpIOEvaluatorSaved := append([]byte(nil), b512HeatpumpIOEvaluatorFirst...)
	for b512HeatpumpIOEvaluatorIndex := range b512HeatpumpIOEvaluatorInput {
		b512HeatpumpIOEvaluatorInput[b512HeatpumpIOEvaluatorIndex] = ' '
	}
	b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, []byte(b512HeatpumpIOEvaluatorFirst), b512HeatpumpIOEvaluatorSaved)
	for b512HeatpumpIOEvaluatorIndex := range b512HeatpumpIOEvaluatorFirst {
		b512HeatpumpIOEvaluatorFirst[b512HeatpumpIOEvaluatorIndex] = ' '
	}
	b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, []byte(b512HeatpumpIOEvaluatorSecond), b512HeatpumpIOEvaluatorSaved)
	b512HeatpumpIOEvaluatorNested := b512HeatpumpIOEvaluatorObject(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorSecond)
	b512HeatpumpIOEvaluatorField(b512HeatpumpIOEvaluatorNested, "inverse")["raw"].([]any)[0] = float64(99)
	b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, b512HeatpumpIOEvaluatorNested["raw"].(map[string]any)["inverse"].([]any)[0], float64(18))
	b512HeatpumpIOEvaluatorNested["context"].(map[string]any)["family"].(map[string]any)["nested"].([]any)[0] = true
	b512HeatpumpIOEvaluatorNested["arguments"].(map[string]any)["extra"].(map[string]any)["nested"].([]any)[0] = true
	b512HeatpumpIOEvaluatorThird, b512HeatpumpIOEvaluatorErr := B512HeatpumpIOEvaluate(b512HeatpumpIOEvaluatorOriginal)
	if b512HeatpumpIOEvaluatorErr != nil {
		b512HeatpumpIOEvaluatorT.Fatal(b512HeatpumpIOEvaluatorErr)
	}
	b512HeatpumpIOEvaluatorEqual(b512HeatpumpIOEvaluatorT, []byte(b512HeatpumpIOEvaluatorThird), b512HeatpumpIOEvaluatorSaved)
	if !bytes.Equal(b512HeatpumpIOEvaluatorSecond, b512HeatpumpIOEvaluatorThird) || strings.Contains(string(b512HeatpumpIOEvaluatorThird), "phase_kind") {
		b512HeatpumpIOEvaluatorT.Fatal("cross-call storage or unexpected selection")
	}
}
