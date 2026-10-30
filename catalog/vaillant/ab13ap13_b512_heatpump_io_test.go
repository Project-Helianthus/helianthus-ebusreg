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

type ab13ap13Case struct {
	Ab13Ap13ID       string          `json:"id"`
	Ab13Ap13Input    json.RawMessage `json:"input"`
	Ab13Ap13Expected map[string]any  `json:"expected"`
}

func ab13ap13Cases(ab13ap13T *testing.T) []ab13ap13Case {
	ab13ap13T.Helper()
	ab13ap13Bytes, ab13ap13Err := os.ReadFile("testdata/ab13ap13_b512_heatpump_io.json")
	if ab13ap13Err != nil {
		ab13ap13T.Fatal(ab13ap13Err)
	}
	var ab13ap13Cases []ab13ap13Case
	if ab13ap13Err := json.Unmarshal(ab13ap13Bytes, &ab13ap13Cases); ab13ap13Err != nil {
		ab13ap13T.Fatal(ab13ap13Err)
	}
	if len(ab13ap13Cases) != 4 {
		ab13ap13T.Fatal("four synthetic views required")
	}
	return ab13ap13Cases
}

func ab13ap13Object(ab13ap13T *testing.T, ab13ap13Bytes []byte) map[string]any {
	ab13ap13T.Helper()
	var ab13ap13Object map[string]any
	if ab13ap13Err := json.Unmarshal(ab13ap13Bytes, &ab13ap13Object); ab13ap13Err != nil {
		ab13ap13T.Fatal(ab13ap13Err)
	}
	return ab13ap13Object
}

func ab13ap13Encode(ab13ap13T *testing.T, ab13ap13Value any) []byte {
	ab13ap13T.Helper()
	ab13ap13Bytes, ab13ap13Err := json.Marshal(ab13ap13Value)
	if ab13ap13Err != nil {
		ab13ap13T.Fatal(ab13ap13Err)
	}
	return ab13ap13Bytes
}

func ab13ap13Evaluate(ab13ap13T *testing.T, ab13ap13Input any) map[string]any {
	ab13ap13T.Helper()
	ab13ap13Bytes, ab13ap13Err := Ab13Ap13Evaluate(ab13ap13Encode(ab13ap13T, ab13ap13Input))
	if ab13ap13Err != nil {
		ab13ap13T.Fatal(ab13ap13Err)
	}
	ab13ap13Output := ab13ap13Object(ab13ap13T, ab13ap13Bytes)
	if ab13ap13Output["physical"] != "Unknown" || ab13ap13Output["qualified"] != false {
		ab13ap13T.Fatal("local interpretation qualified physical support")
	}
	return ab13ap13Output
}

func ab13ap13Field(ab13ap13Output map[string]any, ab13ap13Name string) map[string]any {
	return ab13ap13Output["fields"].(map[string]any)[ab13ap13Name].(map[string]any)
}

func ab13ap13Equal(ab13ap13T *testing.T, ab13ap13Actual, ab13ap13Expected any) {
	ab13ap13T.Helper()
	if !reflect.DeepEqual(ab13ap13Actual, ab13ap13Expected) {
		ab13ap13T.Fatalf("got %#v, want %#v", ab13ap13Actual, ab13ap13Expected)
	}
}

func ab13ap13Absent(ab13ap13T *testing.T, ab13ap13Object map[string]any, ab13ap13Keys ...string) {
	ab13ap13T.Helper()
	for _, ab13ap13Key := range ab13ap13Keys {
		if _, ab13ap13Exists := ab13ap13Object[ab13ap13Key]; ab13ap13Exists {
			ab13ap13T.Fatalf("unexpected %s: %#v", ab13ap13Key, ab13ap13Object)
		}
	}
}

func TestAb13Ap13Fixtures(ab13ap13T *testing.T) {
	for _, ab13ap13Case := range ab13ap13Cases(ab13ap13T) {
		ab13ap13T.Run(ab13ap13Case.Ab13Ap13ID, func(ab13ap13T *testing.T) {
			ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Case.Ab13Ap13Input)
			ab13ap13Output := ab13ap13Evaluate(ab13ap13T, ab13ap13Input)
			ab13ap13Equal(ab13ap13T, ab13ap13Output["context"], ab13ap13Input["context"])
			ab13ap13Equal(ab13ap13T, ab13ap13Output["raw"], ab13ap13Input["fields"])
			ab13ap13Equal(ab13ap13T, ab13ap13Output["state"], "complete")
			for ab13ap13Name, ab13ap13Expected := range ab13ap13Case.Ab13Ap13Expected {
				for ab13ap13Key, ab13ap13Value := range ab13ap13Expected.(map[string]any) {
					ab13ap13Equal(ab13ap13T, ab13ap13Field(ab13ap13Output, ab13ap13Name)[ab13ap13Key], ab13ap13Value)
				}
			}
		})
	}
}

func TestAb13Ap13Temperatures(ab13ap13T *testing.T) {
	for _, ab13ap13Name := range []string{"temperature1", "temperature2"} {
		for ab13ap13Index, ab13ap13Bytes := range [][]int{{0, 128}, {240, 255}, {0, 0}, {16, 0}, {255, 127}, {1, 128}} {
			ab13ap13T.Run(fmt.Sprintf("ab13ap13_%s_%d", ab13ap13Name, ab13ap13Index), func(ab13ap13T *testing.T) {
				ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Cases(ab13ap13T)[1].Ab13Ap13Input)
				ab13ap13Input["fields"].(map[string]any)[ab13ap13Name] = ab13ap13Bytes
				ab13ap13Result := ab13ap13Field(ab13ap13Evaluate(ab13ap13T, ab13ap13Input), ab13ap13Name)
				ab13ap13Equal(ab13ap13T, ab13ap13Result["state"], "Hypothesis")
				if ab13ap13Index == 0 {
					ab13ap13Equal(ab13ap13T, ab13ap13Result["label"], "replacement")
					ab13ap13Absent(ab13ap13T, ab13ap13Result, "number", "calculation")
				} else {
					ab13ap13Equal(ab13ap13T, ab13ap13Result["number"], []float64{0, -1, 0, 1, 2047.9375, -2047.9375}[ab13ap13Index])
					ab13ap13Equal(ab13ap13T, ab13ap13Result["calculation"], "Proven")
					ab13ap13Absent(ab13ap13T, ab13ap13Result, "label")
				}
			})
		}
	}
}

func TestAb13Ap13CodesAndPressure(ab13ap13T *testing.T) {
	for _, ab13ap13Row := range []struct {
		ab13ap13View  int
		ab13ap13Name  string
		ab13ap13Byte  int
		ab13ap13Label string
	}{
		{0, "valve", 0, "off"}, {0, "valve", 1, "CH"}, {0, "valve", 2, "DHW"}, {0, "valve", 3, ""}, {0, "valve", 255, ""},
		{1, "input", 0, "multipurpose"}, {1, "input", 1, "utility-contact"}, {1, "input", 2, ""},
		{1, "backup", 0, "off"}, {1, "backup", 1, "step"}, {1, "backup", 2, "step"}, {1, "backup", 3, "step"}, {1, "backup", 4, ""}, {1, "backup", 254, ""}, {1, "backup", 255, "replacement"},
		{3, "board", 0, "opaque-code-0"}, {3, "board", 1, "opaque-code-1"}, {3, "board", 2, "opaque-code-2"}, {3, "board", 3, "opaque-code-3"}, {3, "board", 11, "opaque-code-11"}, {3, "board", 12, "opaque-code-12"}, {3, "board", 13, ""}, {3, "board", 255, ""},
	} {
		ab13ap13T.Run(fmt.Sprintf("ab13ap13_%s_%d", ab13ap13Row.ab13ap13Name, ab13ap13Row.ab13ap13Byte), func(ab13ap13T *testing.T) {
			ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Cases(ab13ap13T)[ab13ap13Row.ab13ap13View].Ab13Ap13Input)
			ab13ap13Input["fields"].(map[string]any)[ab13ap13Row.ab13ap13Name] = []int{ab13ap13Row.ab13ap13Byte}
			ab13ap13Result := ab13ap13Field(ab13ap13Evaluate(ab13ap13T, ab13ap13Input), ab13ap13Row.ab13ap13Name)
			if ab13ap13Row.ab13ap13Label == "" {
				ab13ap13Equal(ab13ap13T, ab13ap13Result["state"], "Unknown")
				ab13ap13Absent(ab13ap13T, ab13ap13Result, "label", "number")
			} else {
				ab13ap13Equal(ab13ap13T, ab13ap13Result["label"], ab13ap13Row.ab13ap13Label)
				if ab13ap13Row.ab13ap13Label == "step" {
					ab13ap13Equal(ab13ap13T, ab13ap13Result["number"], float64(ab13ap13Row.ab13ap13Byte))
				} else {
					ab13ap13Absent(ab13ap13T, ab13ap13Result, "number")
				}
			}
		})
	}
	for _, ab13ap13Row := range []struct {
		ab13ap13Byte   int
		ab13ap13Number float64
	}{{0, 0}, {1, 0.1}, {10, 1}, {254, 25.4}, {255, 25.5}} {
		ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Cases(ab13ap13T)[1].Ab13Ap13Input)
		ab13ap13Input["fields"].(map[string]any)["pressure"] = []int{ab13ap13Row.ab13ap13Byte}
		ab13ap13Result := ab13ap13Field(ab13ap13Evaluate(ab13ap13T, ab13ap13Input), "pressure")
		ab13ap13Equal(ab13ap13T, ab13ap13Result["number"], ab13ap13Row.ab13ap13Number)
		ab13ap13Default := ab13ap13Result["alternatives"].([]any)[0].(map[string]any)
		ab13ap13Equal(ab13ap13T, ab13ap13Default["number"], float64(255))
		ab13ap13Equal(ab13ap13T, ab13ap13Default["candidate_domain"], false)
		ab13ap13Equal(ab13ap13T, ab13ap13Default["conflict"], true)
	}
}

func TestAb13Ap13EveryBit(ab13ap13T *testing.T) {
	for _, ab13ap13Row := range []struct {
		ab13ap13View   int
		ab13ap13Name   string
		ab13ap13Labels []string
	}{
		{0, "output", []string{"output1", "output2"}},
		{2, "actuator1", []string{"cooling-valve", "passive-mixer1", "passive-mixer2", "cooling-active", "well-pump", "backup-safety", "backup1", "backup2"}},
		{2, "actuator2", []string{"backup3", "external-valve", "external-backup"}},
		{3, "status", []string{"anode1", "anode2", "smart-grid", "well", "safety-limiter", "safety"}},
	} {
		for ab13ap13Bit := -1; ab13ap13Bit < 8; ab13ap13Bit++ {
			ab13ap13T.Run(fmt.Sprintf("ab13ap13_%s_bit%d", ab13ap13Row.ab13ap13Name, ab13ap13Bit), func(ab13ap13T *testing.T) {
				ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Cases(ab13ap13T)[ab13ap13Row.ab13ap13View].Ab13Ap13Input)
				ab13ap13Byte := 0
				if ab13ap13Bit >= 0 {
					ab13ap13Byte = 1 << ab13ap13Bit
				}
				ab13ap13Input["fields"].(map[string]any)[ab13ap13Row.ab13ap13Name] = []int{ab13ap13Byte}
				ab13ap13Result := ab13ap13Field(ab13ap13Evaluate(ab13ap13T, ab13ap13Input), ab13ap13Row.ab13ap13Name)
				ab13ap13Flags := ab13ap13Result["flags"].(map[string]any)
				ab13ap13Equal(ab13ap13T, len(ab13ap13Flags), len(ab13ap13Row.ab13ap13Labels))
				for ab13ap13Index, ab13ap13Label := range ab13ap13Row.ab13ap13Labels {
					ab13ap13Equal(ab13ap13T, ab13ap13Flags[ab13ap13Label], ab13ap13Index == ab13ap13Bit)
				}
				ab13ap13Reserved := 0
				if ab13ap13Bit >= len(ab13ap13Row.ab13ap13Labels) {
					ab13ap13Reserved = ab13ap13Byte
				}
				ab13ap13Equal(ab13ap13T, ab13ap13Result["reserved"], float64(ab13ap13Reserved))
			})
		}
	}
}

func TestAb13Ap13ADCChannels(ab13ap13T *testing.T) {
	for _, ab13ap13Name := range []string{"adc0", "adc1", "adc2", "adc3", "adc4", "adc5"} {
		for _, ab13ap13Row := range []struct {
			ab13ap13Bytes  []int
			ab13ap13Number float64
		}{{[]int{0, 0}, 0}, {[]int{1, 0}, 1}, {[]int{0, 1}, 256}, {[]int{18, 52}, 13330}, {[]int{255, 255}, 65535}} {
			ab13ap13T.Run(fmt.Sprintf("ab13ap13_%s_%v", ab13ap13Name, ab13ap13Row.ab13ap13Bytes), func(ab13ap13T *testing.T) {
				ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Cases(ab13ap13T)[3].Ab13Ap13Input)
				ab13ap13Control := ab13ap13Evaluate(ab13ap13T, ab13ap13Input)
				ab13ap13Input["fields"].(map[string]any)[ab13ap13Name] = ab13ap13Row.ab13ap13Bytes
				ab13ap13Output := ab13ap13Evaluate(ab13ap13T, ab13ap13Input)
				ab13ap13Equal(ab13ap13T, ab13ap13Field(ab13ap13Output, ab13ap13Name)["number"], ab13ap13Row.ab13ap13Number)
				for ab13ap13Sibling, ab13ap13Result := range ab13ap13Control["fields"].(map[string]any) {
					if ab13ap13Sibling != ab13ap13Name {
						ab13ap13Equal(ab13ap13T, ab13ap13Output["fields"].(map[string]any)[ab13ap13Sibling], ab13ap13Result)
					}
				}
			})
		}
	}
}

func TestAb13Ap13OpaqueInverseAndPhase(ab13ap13T *testing.T) {
	for _, ab13ap13Bytes := range [][]int{{18, 52}, {52, 18}, {0, 0}, {255, 255}} {
		ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Cases(ab13ap13T)[2].Ab13Ap13Input)
		ab13ap13Input["fields"].(map[string]any)["inverse"] = ab13ap13Bytes
		ab13ap13Result := ab13ap13Field(ab13ap13Evaluate(ab13ap13T, ab13ap13Input), "inverse")
		ab13ap13Equal(ab13ap13T, ab13ap13Result["raw"], ab13ap13Object(ab13ap13T, ab13ap13Encode(ab13ap13T, map[string]any{"bytes": ab13ap13Bytes}))["bytes"])
		ab13ap13Equal(ab13ap13T, ab13ap13Result["state"], "Unknown")
		ab13ap13Absent(ab13ap13T, ab13ap13Result, "number", "label", "flags")
		ab13ap13Equal(ab13ap13T, ab13ap13Result["unmodeled"], map[string]any{"order": "Unknown", "algorithm": "Unknown"})
	}
	for _, ab13ap13Name := range []string{"phase1", "phase2"} {
		for _, ab13ap13Byte := range []int{0, 33, 50, 66, 100, 255} {
			ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Cases(ab13ap13T)[3].Ab13Ap13Input)
			ab13ap13Input["fields"].(map[string]any)[ab13ap13Name] = []int{ab13ap13Byte}
			ab13ap13Result := ab13ap13Field(ab13ap13Evaluate(ab13ap13T, ab13ap13Input), ab13ap13Name)
			ab13ap13Equal(ab13ap13T, ab13ap13Result["number"], float64(ab13ap13Byte))
			ab13ap13Absent(ab13ap13T, ab13ap13Result, "label", "phase_kind", "selected", "tolerance")
			ab13ap13Equal(ab13ap13T, ab13ap13Result["alternatives"], []any{map[string]any{"state": "Hypothesis", "label": "single-phase", "approximate": []any{float64(50)}}, map[string]any{"state": "Hypothesis", "label": "three-phase", "approximate": []any{float64(0), float64(33), float64(66), float64(100)}}})
		}
	}
}

func TestAb13Ap13ContextAxes(ab13ap13T *testing.T) {
	for _, ab13ap13Case := range ab13ap13Cases(ab13ap13T) {
		for _, ab13ap13Axis := range []string{"primary", "secondary", "operation", "direction", "view", "family", "variant"} {
			for ab13ap13Index, ab13ap13Value := range []any{nil, float64(0), "other", false, map[string]any{"opaque": []any{nil, false, float64(0)}}} {
				if ab13ap13Axis == "primary" || ab13ap13Axis == "secondary" || ab13ap13Axis == "operation" {
					if ab13ap13Index > 1 {
						continue
					}
				}
				if ab13ap13Axis == "direction" || ab13ap13Axis == "view" {
					if ab13ap13Index == 1 || ab13ap13Index > 2 {
						continue
					}
				}
				ab13ap13T.Run(fmt.Sprintf("%s_%s_%d", ab13ap13Case.Ab13Ap13ID, ab13ap13Axis, ab13ap13Index), func(ab13ap13T *testing.T) {
					ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Case.Ab13Ap13Input)
					ab13ap13Context := ab13ap13Input["context"].(map[string]any)
					ab13ap13Context[ab13ap13Axis] = ab13ap13Value
					ab13ap13Output := ab13ap13Evaluate(ab13ap13T, ab13ap13Input)
					ab13ap13Equal(ab13ap13T, ab13ap13Output["context"], ab13ap13Context)
					ab13ap13Equal(ab13ap13T, ab13ap13Output["raw"], ab13ap13Input["fields"])
					if ab13ap13Axis == "family" || ab13ap13Axis == "variant" {
						ab13ap13Equal(ab13ap13T, ab13ap13Output["state"], "complete")
					} else {
						ab13ap13Equal(ab13ap13T, ab13ap13Output["state"], "Unknown")
						ab13ap13Equal(ab13ap13T, len(ab13ap13Output["fields"].(map[string]any)), 0)
					}
				})
			}
			ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Case.Ab13Ap13Input)
			delete(ab13ap13Input["context"].(map[string]any), ab13ap13Axis)
			ab13ap13Output := ab13ap13Evaluate(ab13ap13T, ab13ap13Input)
			ab13ap13Absent(ab13ap13T, ab13ap13Output["context"].(map[string]any), ab13ap13Axis)
			if ab13ap13Axis != "family" && ab13ap13Axis != "variant" {
				ab13ap13Equal(ab13ap13T, ab13ap13Output["state"], "Unknown")
			}
		}
	}
	for _, ab13ap13Pair := range []struct {
		ab13ap13Operation int
		ab13ap13View      string
	}{{15, "board-io"}, {16, "valve-io"}, {17, "valve-io"}} {
		ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Cases(ab13ap13T)[0].Ab13Ap13Input)
		ab13ap13Input["context"].(map[string]any)["operation"] = ab13ap13Pair.ab13ap13Operation
		ab13ap13Input["context"].(map[string]any)["view"] = ab13ap13Pair.ab13ap13View
		ab13ap13Equal(ab13ap13T, ab13ap13Evaluate(ab13ap13T, ab13ap13Input)["state"], "Unknown")
	}
}

func TestAb13Ap13EveryFieldPartial(ab13ap13T *testing.T) {
	for _, ab13ap13Case := range ab13ap13Cases(ab13ap13T) {
		ab13ap13Original := ab13ap13Object(ab13ap13T, ab13ap13Case.Ab13Ap13Input)
		for ab13ap13Name, ab13ap13OriginalBytes := range ab13ap13Original["fields"].(map[string]any) {
			ab13ap13Width := len(ab13ap13OriginalBytes.([]any))
			for ab13ap13Length := -1; ab13ap13Length <= ab13ap13Width+1; ab13ap13Length++ {
				ab13ap13T.Run(fmt.Sprintf("%s_%s_length%d", ab13ap13Case.Ab13Ap13ID, ab13ap13Name, ab13ap13Length), func(ab13ap13T *testing.T) {
					ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Case.Ab13Ap13Input)
					ab13ap13Control := ab13ap13Evaluate(ab13ap13T, ab13ap13Input)
					ab13ap13Raw := ab13ap13Input["fields"].(map[string]any)
					if ab13ap13Length < 0 {
						delete(ab13ap13Raw, ab13ap13Name)
					} else {
						ab13ap13Raw[ab13ap13Name] = make([]int, ab13ap13Length)
					}
					ab13ap13Output := ab13ap13Evaluate(ab13ap13T, ab13ap13Input)
					ab13ap13Result := ab13ap13Field(ab13ap13Output, ab13ap13Name)
					if ab13ap13Length != ab13ap13Width {
						ab13ap13Equal(ab13ap13T, ab13ap13Output["state"], "partial")
						ab13ap13Equal(ab13ap13T, ab13ap13Result["state"], "incomplete")
						ab13ap13Absent(ab13ap13T, ab13ap13Result, "number", "label", "flags", "calculation")
					} else {
						ab13ap13Equal(ab13ap13T, ab13ap13Output["state"], "complete")
					}
					if ab13ap13Length < 0 {
						ab13ap13Absent(ab13ap13T, ab13ap13Result, "raw")
					} else {
						ab13ap13Equal(ab13ap13T, len(ab13ap13Result["raw"].([]any)), ab13ap13Length)
					}
					ab13ap13Equal(ab13ap13T, ab13ap13Output["raw"], ab13ap13Object(ab13ap13T, ab13ap13Encode(ab13ap13T, ab13ap13Raw)))
					for ab13ap13Sibling, ab13ap13Result := range ab13ap13Control["fields"].(map[string]any) {
						if ab13ap13Sibling != ab13ap13Name {
							ab13ap13Equal(ab13ap13T, ab13ap13Output["fields"].(map[string]any)[ab13ap13Sibling], ab13ap13Result)
						}
					}
				})
			}
		}
	}
}

func TestAb13Ap13InvalidJSONAndKinds(ab13ap13T *testing.T) {
	for ab13ap13Index, ab13ap13Input := range []string{
		``, `{`, `[]`, `null`, `{} {}`, `{"context":{},"fields":{},}`, `{"context":{},"context":{},"fields":{}}`,
		`{"context":{"family":{"x":1,"x":2}},"fields":{}}`, `{"context":{},"fields":{"x":[1],"x":[2]}}`,
		`{"context":{},"fields":{},"arguments":{"extra":[{"x":false,"x":true}]}}`,
		`{"context":null,"fields":{}}`, `{"context":{},"fields":null}`, `{"context":{},"fields":{},"arguments":null}`,
	} {
		ab13ap13T.Run(fmt.Sprintf("ab13ap13_invalid_%d", ab13ap13Index), func(ab13ap13T *testing.T) {
			ab13ap13Result, ab13ap13Err := Ab13Ap13Evaluate(json.RawMessage(ab13ap13Input))
			if ab13ap13Err == nil || ab13ap13Result != nil {
				ab13ap13T.Fatalf("invalid input succeeded: %s", ab13ap13Input)
			}
		})
	}
	for _, ab13ap13Case := range ab13ap13Cases(ab13ap13T) {
		ab13ap13Original := ab13ap13Object(ab13ap13T, ab13ap13Case.Ab13Ap13Input)
		for ab13ap13Name := range ab13ap13Original["fields"].(map[string]any) {
			for ab13ap13Index, ab13ap13Value := range []any{nil, false, "0", 0, []any{false}, []any{nil}, []any{"1"}, []any{-1}, []any{256}, []any{1.5}, []any{4294967296.0}, []any{[]any{0}}} {
				ab13ap13T.Run(fmt.Sprintf("%s_%s_invalid%d", ab13ap13Case.Ab13Ap13ID, ab13ap13Name, ab13ap13Index), func(ab13ap13T *testing.T) {
					ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Case.Ab13Ap13Input)
					ab13ap13Input["fields"].(map[string]any)[ab13ap13Name] = ab13ap13Value
					ab13ap13Result, ab13ap13Err := Ab13Ap13Evaluate(ab13ap13Encode(ab13ap13T, ab13ap13Input))
					if ab13ap13Err == nil || ab13ap13Result != nil {
						ab13ap13T.Fatalf("invalid byte kind succeeded: %#v", ab13ap13Value)
					}
				})
			}
		}
	}
	for _, ab13ap13Name := range []string{"primary", "secondary", "operation", "direction", "view"} {
		ab13ap13Values := []any{false, []any{}, map[string]any{}, 1.5, -1, 256, "181"}
		if ab13ap13Name == "direction" || ab13ap13Name == "view" {
			ab13ap13Values = []any{false, 0, []any{}, map[string]any{}}
		}
		for _, ab13ap13Value := range ab13ap13Values {
			ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Cases(ab13ap13T)[0].Ab13Ap13Input)
			ab13ap13Input["context"].(map[string]any)[ab13ap13Name] = ab13ap13Value
			ab13ap13Result, ab13ap13Err := Ab13Ap13Evaluate(ab13ap13Encode(ab13ap13T, ab13ap13Input))
			if ab13ap13Err == nil || ab13ap13Result != nil {
				ab13ap13T.Fatalf("invalid context %s succeeded", ab13ap13Name)
			}
		}
	}
	for _, ab13ap13Value := range []any{nil, false, "", 0, []any{}} {
		ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Cases(ab13ap13T)[0].Ab13Ap13Input)
		ab13ap13Input["arguments"] = ab13ap13Value
		ab13ap13Result, ab13ap13Err := Ab13Ap13Evaluate(ab13ap13Encode(ab13ap13T, ab13ap13Input))
		if ab13ap13Err == nil || ab13ap13Result != nil {
			ab13ap13T.Fatal("invalid arguments object succeeded")
		}
	}
	for _, ab13ap13Number := range []string{"1.0", "1e0", "-0", "256", "1e999"} {
		ab13ap13Result, ab13ap13Err := Ab13Ap13Evaluate(json.RawMessage(`{"context":{},"fields":{"extra":[` + ab13ap13Number + `]}}`))
		if ab13ap13Err == nil || ab13ap13Result != nil {
			ab13ap13T.Fatal("invalid byte representation succeeded")
		}
	}
}

func TestAb13Ap13ArgumentNumbers(ab13ap13T *testing.T) {
	for _, ab13ap13Number := range []string{"-1", "0.5", "1.0", "1e0", "4294967296", "1e999"} {
		ab13ap13Input := json.RawMessage(`{"context":{},"fields":{},"arguments":{"extra":[{"nested":` + ab13ap13Number + `}]}}`)
		ab13ap13Result, ab13ap13Err := Ab13Ap13Evaluate(ab13ap13Input)
		if ab13ap13Err == nil || ab13ap13Result != nil {
			ab13ap13T.Fatalf("invalid argument number succeeded: %s", ab13ap13Number)
		}
	}
}

func TestAb13Ap13ArgumentsExtrasAndStorage(ab13ap13T *testing.T) {
	for _, ab13ap13Case := range ab13ap13Cases(ab13ap13T) {
		for _, ab13ap13Value := range []any{nil, false, true, "", float64(0), float64(4294967295), []any{}, map[string]any{}, []any{nil, false, map[string]any{"nested": []any{float64(0), "opaque"}}}} {
			ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Case.Ab13Ap13Input)
			ab13ap13Input["arguments"] = map[string]any{"extra": ab13ap13Value}
			ab13ap13Input["context"].(map[string]any)["opaque"] = ab13ap13Value
			ab13ap13Input["fields"].(map[string]any)["extra"] = []int{0, 255}
			ab13ap13Output := ab13ap13Evaluate(ab13ap13T, ab13ap13Input)
			ab13ap13Equal(ab13ap13T, ab13ap13Output["arguments"], ab13ap13Input["arguments"])
			ab13ap13Equal(ab13ap13T, ab13ap13Output["context"], ab13ap13Input["context"])
			ab13ap13Equal(ab13ap13T, ab13ap13Output["state"], "complete")
			ab13ap13Absent(ab13ap13T, ab13ap13Output["fields"].(map[string]any), "extra")
		}
		ab13ap13Input := ab13ap13Object(ab13ap13T, ab13ap13Case.Ab13Ap13Input)
		ab13ap13Absent(ab13ap13T, ab13ap13Evaluate(ab13ap13T, ab13ap13Input), "arguments")
		ab13ap13Input["arguments"] = map[string]any{}
		ab13ap13Equal(ab13ap13T, ab13ap13Evaluate(ab13ap13T, ab13ap13Input)["arguments"], map[string]any{})
	}
	ab13ap13Input := []byte(`{"context":{"primary":181,"secondary":18,"operation":16,"direction":"request","view":"board-io","family":{"nested":[false,0]},"variant":null},"fields":{"actuator1":[0],"actuator2":[0],"inverse":[18,52]},"arguments":{"extra":{"nested":[false,0]}}}`)
	ab13ap13Original := append([]byte(nil), ab13ap13Input...)
	ab13ap13First, ab13ap13Err := Ab13Ap13Evaluate(ab13ap13Input)
	if ab13ap13Err != nil {
		ab13ap13T.Fatal(ab13ap13Err)
	}
	ab13ap13Equal(ab13ap13T, ab13ap13Input, ab13ap13Original)
	ab13ap13Second, ab13ap13Err := Ab13Ap13Evaluate(ab13ap13Input)
	if ab13ap13Err != nil {
		ab13ap13T.Fatal(ab13ap13Err)
	}
	ab13ap13Saved := append([]byte(nil), ab13ap13First...)
	for ab13ap13Index := range ab13ap13Input {
		ab13ap13Input[ab13ap13Index] = ' '
	}
	ab13ap13Equal(ab13ap13T, []byte(ab13ap13First), ab13ap13Saved)
	for ab13ap13Index := range ab13ap13First {
		ab13ap13First[ab13ap13Index] = ' '
	}
	ab13ap13Equal(ab13ap13T, []byte(ab13ap13Second), ab13ap13Saved)
	ab13ap13Nested := ab13ap13Object(ab13ap13T, ab13ap13Second)
	ab13ap13Field(ab13ap13Nested, "inverse")["raw"].([]any)[0] = float64(99)
	ab13ap13Equal(ab13ap13T, ab13ap13Nested["raw"].(map[string]any)["inverse"].([]any)[0], float64(18))
	ab13ap13Nested["context"].(map[string]any)["family"].(map[string]any)["nested"].([]any)[0] = true
	ab13ap13Nested["arguments"].(map[string]any)["extra"].(map[string]any)["nested"].([]any)[0] = true
	ab13ap13Third, ab13ap13Err := Ab13Ap13Evaluate(ab13ap13Original)
	if ab13ap13Err != nil {
		ab13ap13T.Fatal(ab13ap13Err)
	}
	ab13ap13Equal(ab13ap13T, []byte(ab13ap13Third), ab13ap13Saved)
	if !bytes.Equal(ab13ap13Second, ab13ap13Third) || strings.Contains(string(ab13ap13Third), "phase_kind") {
		ab13ap13T.Fatal("cross-call storage or unexpected selection")
	}
}
