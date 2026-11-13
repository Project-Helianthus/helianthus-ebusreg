package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func b523ActuatorFieldInputFor(b523ActuatorFieldOperation int, b523ActuatorFieldDirection string) map[string]any {
	b523ActuatorFieldView := "relay-six"
	b523ActuatorFieldFields := map[string]any{"relays": []int{0, 0, 0, 0, 0, 0}, "pwm1": []int{0}}
	if b523ActuatorFieldOperation == 5 {
		b523ActuatorFieldView = "relay-twelve"
		b523ActuatorFieldFields["relays"] = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		b523ActuatorFieldFields["pwm2"] = []int{0}
	}
	if b523ActuatorFieldOperation == 2 {
		b523ActuatorFieldView = "mixer"
		b523ActuatorFieldFields = map[string]any{"circuit": []int{0}, "mode": []int{1}, "setpoint": []int{190}}
	}
	if b523ActuatorFieldDirection == "positive-feedback" {
		b523ActuatorFieldFields = map[string]any{"feedback": []int{1}}
		if b523ActuatorFieldOperation == 2 {
			b523ActuatorFieldFields["movement"] = []int{0}
		}
	}
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 35, "operation": b523ActuatorFieldOperation, "direction": b523ActuatorFieldDirection, "view": b523ActuatorFieldView, "family": "synthetic-family", "variant": "opaque"}, "fields": b523ActuatorFieldFields, "arguments": map[string]any{"variant": "two-circuit", "opaque": map[string]any{"false": false, "zero": 0, "null": nil, "array": []any{0, false, nil}}}}
}

func b523ActuatorFieldEvaluateFor(b523ActuatorFieldT *testing.T, b523ActuatorFieldInput any) map[string]any {
	b523ActuatorFieldT.Helper()
	b523ActuatorFieldBytes, b523ActuatorFieldError := json.Marshal(b523ActuatorFieldInput)
	if b523ActuatorFieldError != nil {
		b523ActuatorFieldT.Fatal(b523ActuatorFieldError)
	}
	b523ActuatorFieldOutput, b523ActuatorFieldError := B523ActuatorFieldEvaluate(b523ActuatorFieldBytes)
	if b523ActuatorFieldError != nil {
		b523ActuatorFieldT.Fatal(b523ActuatorFieldError)
	}
	b523ActuatorFieldDecoded := map[string]any{}
	if b523ActuatorFieldError = json.Unmarshal(b523ActuatorFieldOutput, &b523ActuatorFieldDecoded); b523ActuatorFieldError != nil {
		b523ActuatorFieldT.Fatal(b523ActuatorFieldError)
	}
	if b523ActuatorFieldDecoded["physical"] != "Unknown" || b523ActuatorFieldDecoded["qualified"] != false {
		b523ActuatorFieldT.Fatal("native qualification")
	}
	b523ActuatorFieldOriginal := map[string]any{}
	_ = json.Unmarshal(b523ActuatorFieldBytes, &b523ActuatorFieldOriginal)
	for _, b523ActuatorFieldPair := range [][2]string{{"context", "context"}, {"fields", "raw"}, {"arguments", "arguments"}} {
		if !reflect.DeepEqual(b523ActuatorFieldOriginal[b523ActuatorFieldPair[0]], b523ActuatorFieldDecoded[b523ActuatorFieldPair[1]]) {
			b523ActuatorFieldT.Fatalf("lost %s", b523ActuatorFieldPair[0])
		}
	}
	return b523ActuatorFieldDecoded
}

func b523ActuatorFieldFieldFor(b523ActuatorFieldOutput map[string]any, b523ActuatorFieldKey string) map[string]any {
	return b523ActuatorFieldOutput["fields"].(map[string]any)[b523ActuatorFieldKey].(map[string]any)
}

func TestAb21Ap21Relays(b523ActuatorFieldT *testing.T) {
	for _, b523ActuatorFieldOperation := range []int{1, 5} {
		b523ActuatorFieldLength := 6
		if b523ActuatorFieldOperation == 5 {
			b523ActuatorFieldLength = 12
		}
		for b523ActuatorFieldSlot := 0; b523ActuatorFieldSlot < b523ActuatorFieldLength; b523ActuatorFieldSlot++ {
			for _, b523ActuatorFieldCase := range []struct {
				b523ActuatorFieldRaw   int
				b523ActuatorFieldLabel string
			}{{0, "off"}, {1, "on-for-seconds-then-off"}, {253, "on-for-seconds-then-off"}, {254, "on"}, {255, "unchanged"}} {
				b523ActuatorFieldT.Run(fmt.Sprintf("ab21ap21-%d-%d-%d", b523ActuatorFieldOperation, b523ActuatorFieldSlot, b523ActuatorFieldCase.b523ActuatorFieldRaw), func(b523ActuatorFieldT *testing.T) {
					b523ActuatorFieldInput := b523ActuatorFieldInputFor(b523ActuatorFieldOperation, "request")
					b523ActuatorFieldInput["fields"].(map[string]any)["relays"].([]int)[b523ActuatorFieldSlot] = b523ActuatorFieldCase.b523ActuatorFieldRaw
					b523ActuatorFieldOutput := b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput)
					b523ActuatorFieldRelays := b523ActuatorFieldFieldFor(b523ActuatorFieldOutput, "relays")["relays"].([]any)
					if len(b523ActuatorFieldRelays) != b523ActuatorFieldLength {
						b523ActuatorFieldT.Fatal("relay lost")
					}
					for b523ActuatorFieldIndex, b523ActuatorFieldValue := range b523ActuatorFieldRelays {
						b523ActuatorFieldRelay := b523ActuatorFieldValue.(map[string]any)
						b523ActuatorFieldWant := "off"
						if b523ActuatorFieldIndex == b523ActuatorFieldSlot {
							b523ActuatorFieldWant = b523ActuatorFieldCase.b523ActuatorFieldLabel
						}
						if b523ActuatorFieldRelay["label"] != b523ActuatorFieldWant || b523ActuatorFieldRelay["grade"] != "Hypothesis" {
							b523ActuatorFieldT.Fatal(b523ActuatorFieldRelay)
						}
						if b523ActuatorFieldIndex == b523ActuatorFieldSlot && (b523ActuatorFieldCase.b523ActuatorFieldRaw == 1 || b523ActuatorFieldCase.b523ActuatorFieldRaw == 253) {
							if b523ActuatorFieldRelay["number"] != float64(b523ActuatorFieldCase.b523ActuatorFieldRaw) {
								b523ActuatorFieldT.Fatal("duration")
							}
						} else if _, b523ActuatorFieldExists := b523ActuatorFieldRelay["number"]; b523ActuatorFieldExists {
							b523ActuatorFieldT.Fatal("invented duration")
						}
					}
				})
			}
		}
	}
}

func TestAb21Ap21PWM(b523ActuatorFieldT *testing.T) {
	for _, b523ActuatorFieldOperation := range []int{1, 5} {
		b523ActuatorFieldKeys := []string{"pwm1"}
		if b523ActuatorFieldOperation == 5 {
			b523ActuatorFieldKeys = append(b523ActuatorFieldKeys, "pwm2")
		}
		for _, b523ActuatorFieldKey := range b523ActuatorFieldKeys {
			for _, b523ActuatorFieldCase := range []struct {
				b523ActuatorFieldRaw      int
				b523ActuatorFieldNumber   float64
				b523ActuatorFieldDomain   bool
				b523ActuatorFieldBehavior any
			}{{0, 0, true, float64(0)}, {200, 100, true, float64(100)}, {201, 100.5, true, "unchanged"}, {250, 125, true, "unchanged"}, {255, 127.5, false, "unchanged"}} {
				b523ActuatorFieldT.Run(fmt.Sprintf("ab21ap21-%d-%s-%d", b523ActuatorFieldOperation, b523ActuatorFieldKey, b523ActuatorFieldCase.b523ActuatorFieldRaw), func(b523ActuatorFieldT *testing.T) {
					b523ActuatorFieldInput := b523ActuatorFieldInputFor(b523ActuatorFieldOperation, "request")
					b523ActuatorFieldInput["fields"].(map[string]any)[b523ActuatorFieldKey] = []int{b523ActuatorFieldCase.b523ActuatorFieldRaw}
					b523ActuatorFieldOutput := b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput)
					b523ActuatorFieldAlternatives := b523ActuatorFieldFieldFor(b523ActuatorFieldOutput, b523ActuatorFieldKey)["alternatives"].([]any)
					if len(b523ActuatorFieldAlternatives) != 2 {
						b523ActuatorFieldT.Fatal("lost PWM alternative")
					}
					b523ActuatorFieldNumeric := b523ActuatorFieldAlternatives[0].(map[string]any)
					b523ActuatorFieldBehavior := b523ActuatorFieldAlternatives[1].(map[string]any)
					if b523ActuatorFieldNumeric["number"] != b523ActuatorFieldCase.b523ActuatorFieldNumber || b523ActuatorFieldNumeric["candidate_domain"] != b523ActuatorFieldCase.b523ActuatorFieldDomain || b523ActuatorFieldNumeric["calculation"] != "Proven" || b523ActuatorFieldNumeric["grade"] != "Hypothesis" {
						b523ActuatorFieldT.Fatal(b523ActuatorFieldNumeric)
					}
					if b523ActuatorFieldCase.b523ActuatorFieldRaw <= 200 {
						if b523ActuatorFieldBehavior["number"] != b523ActuatorFieldCase.b523ActuatorFieldBehavior {
							b523ActuatorFieldT.Fatal(b523ActuatorFieldBehavior)
						}
					} else if b523ActuatorFieldBehavior["label"] != "unchanged" || b523ActuatorFieldBehavior["number"] != nil {
						b523ActuatorFieldT.Fatal(b523ActuatorFieldBehavior)
					}
					for _, b523ActuatorFieldSibling := range b523ActuatorFieldKeys {
						if b523ActuatorFieldSibling != b523ActuatorFieldKey {
							if b523ActuatorFieldFieldFor(b523ActuatorFieldOutput, b523ActuatorFieldSibling)["alternatives"].([]any)[0].(map[string]any)["number"] != float64(0) {
								b523ActuatorFieldT.Fatal("PWM sibling changed")
							}
						}
					}
				})
			}
		}
	}
}

func TestAb21Ap21Mixer(b523ActuatorFieldT *testing.T) {
	for _, b523ActuatorFieldVariant := range []string{"two-circuit", "three-circuit"} {
		for _, b523ActuatorFieldCircuit := range []int{0, 1, 2, 3, 255} {
			b523ActuatorFieldInput := b523ActuatorFieldInputFor(2, "request")
			b523ActuatorFieldInput["arguments"].(map[string]any)["variant"] = b523ActuatorFieldVariant
			b523ActuatorFieldInput["fields"].(map[string]any)["circuit"] = []int{b523ActuatorFieldCircuit}
			b523ActuatorFieldWant := "ignored"
			if b523ActuatorFieldCircuit == 0 {
				b523ActuatorFieldWant = "circuit-0"
			}
			if b523ActuatorFieldCircuit == 1 {
				b523ActuatorFieldWant = "circuit-1"
			}
			if b523ActuatorFieldCircuit == 2 && b523ActuatorFieldVariant == "three-circuit" {
				b523ActuatorFieldWant = "circuit-2"
			}
			if b523ActuatorFieldFieldFor(b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput), "circuit")["label"] != b523ActuatorFieldWant {
				b523ActuatorFieldT.Fatal("circuit variant")
			}
		}
	}
	for b523ActuatorFieldMode, b523ActuatorFieldLabel := range []string{"stop/pump-off/close", "heating-regulation", "cooling-regulation", "pump-off/open", "pump-on/close", "pump-on/open", "ignored"} {
		b523ActuatorFieldInput := b523ActuatorFieldInputFor(2, "request")
		b523ActuatorFieldInput["fields"].(map[string]any)["mode"] = []int{b523ActuatorFieldMode}
		b523ActuatorFieldOutput := b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput)
		if b523ActuatorFieldFieldFor(b523ActuatorFieldOutput, "mode")["label"] != b523ActuatorFieldLabel {
			b523ActuatorFieldT.Fatal("mode")
		}
		b523ActuatorFieldField := b523ActuatorFieldFieldFor(b523ActuatorFieldOutput, "setpoint")
		if b523ActuatorFieldField["relevant"] != (b523ActuatorFieldMode == 1 || b523ActuatorFieldMode == 2) {
			b523ActuatorFieldT.Fatal("relevance")
		}
		b523ActuatorFieldAlternatives := b523ActuatorFieldField["alternatives"].([]any)
		if len(b523ActuatorFieldAlternatives) != 2 || b523ActuatorFieldAlternatives[0].(map[string]any)["number"] != float64(-33) || b523ActuatorFieldAlternatives[1].(map[string]any)["number"] != float64(95) || b523ActuatorFieldAlternatives[1].(map[string]any)["clamp"] != float64(95) {
			b523ActuatorFieldT.Fatal(b523ActuatorFieldField)
		}
	}
	for _, b523ActuatorFieldCase := range []struct {
		b523ActuatorFieldRaw    int
		b523ActuatorFieldNumber float64
	}{{0, 0}, {127, 63.5}, {128, -64}, {255, -0.5}} {
		b523ActuatorFieldInput := b523ActuatorFieldInputFor(2, "request")
		b523ActuatorFieldInput["fields"].(map[string]any)["setpoint"] = []int{b523ActuatorFieldCase.b523ActuatorFieldRaw}
		b523ActuatorFieldAlternatives := b523ActuatorFieldFieldFor(b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput), "setpoint")["alternatives"].([]any)
		if b523ActuatorFieldAlternatives[0].(map[string]any)["number"] != b523ActuatorFieldCase.b523ActuatorFieldNumber || b523ActuatorFieldAlternatives[1].(map[string]any)["state"] != "Unknown" || b523ActuatorFieldAlternatives[1].(map[string]any)["number"] != nil {
			b523ActuatorFieldT.Fatal(b523ActuatorFieldAlternatives)
		}
	}
}

func TestAb21Ap21Feedback(b523ActuatorFieldT *testing.T) {
	for _, b523ActuatorFieldOperation := range []int{1, 2, 5} {
		for _, b523ActuatorFieldFeedback := range []int{0, 1, 2, 255} {
			b523ActuatorFieldInput := b523ActuatorFieldInputFor(b523ActuatorFieldOperation, "positive-feedback")
			b523ActuatorFieldInput["fields"].(map[string]any)["feedback"] = []int{b523ActuatorFieldFeedback}
			b523ActuatorFieldWant := "rejected"
			if b523ActuatorFieldOperation == 5 {
				b523ActuatorFieldWant = "rejected-or-unsupported"
			}
			if b523ActuatorFieldFeedback == 1 {
				b523ActuatorFieldWant = "accepted"
			}
			b523ActuatorFieldLabel := b523ActuatorFieldFieldFor(b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput), "feedback")["label"]
			if b523ActuatorFieldLabel != b523ActuatorFieldWant {
				b523ActuatorFieldT.Fatal("feedback operation")
			}
		}
	}
	for _, b523ActuatorFieldCase := range []struct {
		b523ActuatorFieldRaw    int
		b523ActuatorFieldNumber any
		b523ActuatorFieldLabel  any
	}{{0, float64(0), "stopped"}, {156, float64(-100), "closing"}, {100, float64(100), "opening"}, {155, nil, nil}, {101, nil, nil}, {255, float64(-1), "closing"}, {1, float64(1), "opening"}} {
		b523ActuatorFieldInput := b523ActuatorFieldInputFor(2, "positive-feedback")
		b523ActuatorFieldInput["fields"].(map[string]any)["movement"] = []int{b523ActuatorFieldCase.b523ActuatorFieldRaw}
		b523ActuatorFieldField := b523ActuatorFieldFieldFor(b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput), "movement")
		if b523ActuatorFieldField["number"] != b523ActuatorFieldCase.b523ActuatorFieldNumber || b523ActuatorFieldField["label"] != b523ActuatorFieldCase.b523ActuatorFieldLabel {
			b523ActuatorFieldT.Fatal(b523ActuatorFieldField)
		}
		if b523ActuatorFieldCase.b523ActuatorFieldNumber == nil && b523ActuatorFieldField["state"] != "Unknown" {
			b523ActuatorFieldT.Fatal(b523ActuatorFieldField)
		}
	}
}

func TestAb21Ap21Context(b523ActuatorFieldT *testing.T) {
	for _, b523ActuatorFieldKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, b523ActuatorFieldValue := range []any{nil, "ab21ap21-absent", 0} {
			b523ActuatorFieldInput := b523ActuatorFieldInputFor(1, "request")
			b523ActuatorFieldContext := b523ActuatorFieldInput["context"].(map[string]any)
			if b523ActuatorFieldValue == "ab21ap21-absent" {
				delete(b523ActuatorFieldContext, b523ActuatorFieldKey)
			} else if b523ActuatorFieldKey == "direction" || b523ActuatorFieldKey == "view" {
				if b523ActuatorFieldValue == nil {
					b523ActuatorFieldContext[b523ActuatorFieldKey] = nil
				} else {
					b523ActuatorFieldContext[b523ActuatorFieldKey] = "unsupported"
				}
			} else {
				b523ActuatorFieldContext[b523ActuatorFieldKey] = b523ActuatorFieldValue
			}
			b523ActuatorFieldOutput := b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput)
			if b523ActuatorFieldOutput["state"] != "Unknown" || len(b523ActuatorFieldOutput["fields"].(map[string]any)) != 0 {
				b523ActuatorFieldT.Fatal("unsupported interpreted")
			}
		}
	}
	for _, b523ActuatorFieldKey := range []string{"family", "variant", "opaque"} {
		for _, b523ActuatorFieldValue := range []any{nil, 0, false, "synthetic-other", []any{}, map[string]any{"nested": false}} {
			b523ActuatorFieldInput := b523ActuatorFieldInputFor(1, "request")
			b523ActuatorFieldInput["context"].(map[string]any)[b523ActuatorFieldKey] = b523ActuatorFieldValue
			if b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput)["state"] != "complete" {
				b523ActuatorFieldT.Fatal("opaque selector")
			}
		}
		b523ActuatorFieldInput := b523ActuatorFieldInputFor(1, "request")
		delete(b523ActuatorFieldInput["context"].(map[string]any), b523ActuatorFieldKey)
		b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput)
	}
}

func TestAb21Ap21Partial(b523ActuatorFieldT *testing.T) {
	for _, b523ActuatorFieldOperation := range []int{1, 2, 5} {
		for _, b523ActuatorFieldDirection := range []string{"request", "positive-feedback"} {
			b523ActuatorFieldControl := b523ActuatorFieldInputFor(b523ActuatorFieldOperation, b523ActuatorFieldDirection)
			for b523ActuatorFieldKey := range b523ActuatorFieldControl["fields"].(map[string]any) {
				for _, b523ActuatorFieldLength := range []int{-1, 0, 2, 13} {
					b523ActuatorFieldInput := b523ActuatorFieldInputFor(b523ActuatorFieldOperation, b523ActuatorFieldDirection)
					b523ActuatorFieldRaw := b523ActuatorFieldInput["fields"].(map[string]any)
					if b523ActuatorFieldLength == -1 {
						delete(b523ActuatorFieldRaw, b523ActuatorFieldKey)
					} else {
						b523ActuatorFieldRaw[b523ActuatorFieldKey] = make([]int, b523ActuatorFieldLength)
					}
					b523ActuatorFieldOutput := b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput)
					if b523ActuatorFieldOutput["state"] != "partial" || b523ActuatorFieldFieldFor(b523ActuatorFieldOutput, b523ActuatorFieldKey)["state"] != "incomplete" {
						b523ActuatorFieldT.Fatal("incomplete lost")
					}
					for b523ActuatorFieldSibling := range b523ActuatorFieldControl["fields"].(map[string]any) {
						if b523ActuatorFieldSibling != b523ActuatorFieldKey && b523ActuatorFieldFieldFor(b523ActuatorFieldOutput, b523ActuatorFieldSibling)["state"] != "complete" {
							b523ActuatorFieldT.Fatal("independent sibling lost")
						}
					}
				}
			}
		}
	}
}

func TestAb21Ap21Reject(b523ActuatorFieldT *testing.T) {
	for _, b523ActuatorFieldInput := range []string{`{`, `[]`, `null`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"x":{"a":0,"a":1}},"fields":{}}`, `{"context":{},"fields":{},"arguments":{"x":[{"a":0,"a":1}]}}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":null}`, `{"context":{},"fields":{},"arguments":null}`} {
		if b523ActuatorFieldOutput, b523ActuatorFieldError := B523ActuatorFieldEvaluate([]byte(b523ActuatorFieldInput)); b523ActuatorFieldError == nil || b523ActuatorFieldOutput != nil {
			b523ActuatorFieldT.Fatalf("accepted %s", b523ActuatorFieldInput)
		}
	}
	for _, b523ActuatorFieldOperation := range []int{1, 2, 5} {
		for _, b523ActuatorFieldDirection := range []string{"request", "positive-feedback"} {
			b523ActuatorFieldControl := b523ActuatorFieldInputFor(b523ActuatorFieldOperation, b523ActuatorFieldDirection)
			for b523ActuatorFieldKey := range b523ActuatorFieldControl["fields"].(map[string]any) {
				for _, b523ActuatorFieldBad := range []any{nil, false, 0, "0", []any{false}, []any{nil}, []any{"1"}, []any{-1}, []any{256}, []any{1.5}, []any{json.Number("1.00000000000000000001")}} {
					b523ActuatorFieldInput := b523ActuatorFieldInputFor(b523ActuatorFieldOperation, b523ActuatorFieldDirection)
					b523ActuatorFieldInput["fields"].(map[string]any)[b523ActuatorFieldKey] = b523ActuatorFieldBad
					b523ActuatorFieldBytes, _ := json.Marshal(b523ActuatorFieldInput)
					if b523ActuatorFieldOutput, b523ActuatorFieldError := B523ActuatorFieldEvaluate(b523ActuatorFieldBytes); b523ActuatorFieldError == nil || b523ActuatorFieldOutput != nil {
						b523ActuatorFieldT.Fatalf("accepted %s", b523ActuatorFieldBytes)
					}
				}
			}
		}
	}
	for _, b523ActuatorFieldBad := range []any{nil, false, 0, []any{}, "unknown", "missing", "no-arguments"} {
		b523ActuatorFieldInput := b523ActuatorFieldInputFor(2, "request")
		switch b523ActuatorFieldBad {
		case "missing":
			delete(b523ActuatorFieldInput["arguments"].(map[string]any), "variant")
		case "no-arguments":
			delete(b523ActuatorFieldInput, "arguments")
		default:
			b523ActuatorFieldInput["arguments"].(map[string]any)["variant"] = b523ActuatorFieldBad
		}
		b523ActuatorFieldBytes, _ := json.Marshal(b523ActuatorFieldInput)
		if b523ActuatorFieldOutput, b523ActuatorFieldError := B523ActuatorFieldEvaluate(b523ActuatorFieldBytes); b523ActuatorFieldError == nil || b523ActuatorFieldOutput != nil {
			b523ActuatorFieldT.Fatalf("accepted variant %v", b523ActuatorFieldBad)
		}
	}
}

func TestAb21Ap21Isolation(b523ActuatorFieldT *testing.T) {
	b523ActuatorFieldInput, _ := json.Marshal(b523ActuatorFieldInputFor(5, "request"))
	b523ActuatorFieldOriginal := append([]byte(nil), b523ActuatorFieldInput...)
	b523ActuatorFieldFirst, b523ActuatorFieldError := B523ActuatorFieldEvaluate(b523ActuatorFieldInput)
	if b523ActuatorFieldError != nil {
		b523ActuatorFieldT.Fatal(b523ActuatorFieldError)
	}
	b523ActuatorFieldExpected := append([]byte(nil), b523ActuatorFieldFirst...)
	for b523ActuatorFieldIndex := range b523ActuatorFieldInput {
		b523ActuatorFieldInput[b523ActuatorFieldIndex] = ' '
	}
	if !bytes.Equal(b523ActuatorFieldFirst, b523ActuatorFieldExpected) {
		b523ActuatorFieldT.Fatal("input alias")
	}
	b523ActuatorFieldDecoded := map[string]any{}
	_ = json.Unmarshal(b523ActuatorFieldFirst, &b523ActuatorFieldDecoded)
	b523ActuatorFieldDecoded["raw"].(map[string]any)["pwm1"].([]any)[0] = float64(201)
	b523ActuatorFieldDecoded["arguments"].(map[string]any)["opaque"].(map[string]any)["false"] = true
	if b523ActuatorFieldFieldFor(b523ActuatorFieldDecoded, "pwm1")["raw"].([]any)[0] != float64(0) {
		b523ActuatorFieldT.Fatal("nested raw alias")
	}
	for b523ActuatorFieldIndex := range b523ActuatorFieldFirst {
		b523ActuatorFieldFirst[b523ActuatorFieldIndex] = ' '
	}
	b523ActuatorFieldSecond, b523ActuatorFieldError := B523ActuatorFieldEvaluate(b523ActuatorFieldOriginal)
	if b523ActuatorFieldError != nil || !bytes.Equal(b523ActuatorFieldSecond, b523ActuatorFieldExpected) {
		b523ActuatorFieldT.Fatal("result alias")
	}
}

func TestAb21Ap21SyntheticFixture(b523ActuatorFieldT *testing.T) {
	b523ActuatorFieldBytes, b523ActuatorFieldError := os.ReadFile("testdata/ab21ap21_b523_actuator_fields.json")
	if b523ActuatorFieldError != nil {
		b523ActuatorFieldT.Fatal(b523ActuatorFieldError)
	}
	var b523ActuatorFieldCases []struct {
		B523ActuatorFieldID    string          `json:"id"`
		B523ActuatorFieldInput json.RawMessage `json:"input"`
		B523ActuatorFieldState string          `json:"state"`
	}
	if b523ActuatorFieldError = json.Unmarshal(b523ActuatorFieldBytes, &b523ActuatorFieldCases); b523ActuatorFieldError != nil || len(b523ActuatorFieldCases) == 0 {
		b523ActuatorFieldT.Fatal("empty fixture", b523ActuatorFieldError)
	}
	for _, b523ActuatorFieldCase := range b523ActuatorFieldCases {
		b523ActuatorFieldT.Run(b523ActuatorFieldCase.B523ActuatorFieldID, func(b523ActuatorFieldT *testing.T) {
			b523ActuatorFieldOutput := b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldCase.B523ActuatorFieldInput)
			if b523ActuatorFieldOutput["state"] != b523ActuatorFieldCase.B523ActuatorFieldState {
				b523ActuatorFieldT.Fatal(b523ActuatorFieldOutput)
			}
		})
	}
}

func TestAb21Ap21ValidationBoundaries(b523ActuatorFieldT *testing.T) {
	for _, b523ActuatorFieldKey := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, b523ActuatorFieldBad := range []any{false, []any{}, map[string]any{}, 1.5, -1, 256} {
			b523ActuatorFieldInput := b523ActuatorFieldInputFor(1, "request")
			b523ActuatorFieldInput["context"].(map[string]any)[b523ActuatorFieldKey] = b523ActuatorFieldBad
			b523ActuatorFieldBytes, _ := json.Marshal(b523ActuatorFieldInput)
			if b523ActuatorFieldOutput, b523ActuatorFieldError := B523ActuatorFieldEvaluate(b523ActuatorFieldBytes); b523ActuatorFieldError == nil || b523ActuatorFieldOutput != nil {
				b523ActuatorFieldT.Fatalf("bad context accepted %s", b523ActuatorFieldBytes)
			}
		}
	}
	for _, b523ActuatorFieldKey := range []string{"primary", "secondary", "operation"} {
		b523ActuatorFieldInput := b523ActuatorFieldInputFor(1, "request")
		b523ActuatorFieldInput["context"].(map[string]any)[b523ActuatorFieldKey] = "1"
		b523ActuatorFieldBytes, _ := json.Marshal(b523ActuatorFieldInput)
		if _, b523ActuatorFieldError := B523ActuatorFieldEvaluate(b523ActuatorFieldBytes); b523ActuatorFieldError == nil {
			b523ActuatorFieldT.Fatal("coerced context")
		}
	}
	for _, b523ActuatorFieldInput := range []string{`{"context":{},"fields":{"x":[],"x":[]}}`, `{"context":{},"fields":{},"arguments":{"variant":"two-circuit","variant":"three-circuit"}}`, `{"context":{"x":[{"y":1,"y":2}]},"fields":{}}`, `{"context":{},"fields":{"x":[1e999]}}`} {
		if b523ActuatorFieldOutput, b523ActuatorFieldError := B523ActuatorFieldEvaluate([]byte(b523ActuatorFieldInput)); b523ActuatorFieldError == nil || b523ActuatorFieldOutput != nil {
			b523ActuatorFieldT.Fatal("invalid JSON accepted")
		}
	}
	for _, b523ActuatorFieldArguments := range []any{map[string]any{}, map[string]any{"extra": false}, map[string]any{"extra": nil}, map[string]any{"extra": 0}, map[string]any{"extra": []any{}}} {
		b523ActuatorFieldInput := b523ActuatorFieldInputFor(1, "request")
		b523ActuatorFieldInput["arguments"] = b523ActuatorFieldArguments
		b523ActuatorFieldInput["fields"].(map[string]any)["extra"] = []int{0, 255}
		b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput)
	}
	b523ActuatorFieldInput := b523ActuatorFieldInputFor(1, "request")
	delete(b523ActuatorFieldInput, "arguments")
	b523ActuatorFieldOutput := b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput)
	if _, b523ActuatorFieldPresent := b523ActuatorFieldOutput["arguments"]; b523ActuatorFieldPresent {
		b523ActuatorFieldT.Fatal("invented arguments")
	}
	b523ActuatorFieldInput = b523ActuatorFieldInputFor(2, "request")
	delete(b523ActuatorFieldInput["fields"].(map[string]any), "mode")
	b523ActuatorFieldField := b523ActuatorFieldFieldFor(b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput), "setpoint")
	if b523ActuatorFieldField["relevance"] != "Unknown" || b523ActuatorFieldField["relevant"] != nil {
		b523ActuatorFieldT.Fatal("invented relevance")
	}
	for _, b523ActuatorFieldMovement := range []int{156, 100} {
		b523ActuatorFieldInput := b523ActuatorFieldInputFor(2, "positive-feedback")
		b523ActuatorFieldInput["fields"].(map[string]any)["movement"] = []int{b523ActuatorFieldMovement}
		b523ActuatorFieldField := b523ActuatorFieldFieldFor(b523ActuatorFieldEvaluateFor(b523ActuatorFieldT, b523ActuatorFieldInput), "movement")
		if b523ActuatorFieldField["maximum"] != true || b523ActuatorFieldField["grade"] != "Hypothesis" || b523ActuatorFieldField["calculation"] != "Proven" {
			b523ActuatorFieldT.Fatal("movement grade or maximum")
		}
	}
}
