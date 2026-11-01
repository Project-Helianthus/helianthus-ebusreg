package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func ab21ap21InputFor(ab21ap21Operation int, ab21ap21Direction string) map[string]any {
	ab21ap21View := "relay-six"
	ab21ap21Fields := map[string]any{"relays": []int{0, 0, 0, 0, 0, 0}, "pwm1": []int{0}}
	if ab21ap21Operation == 5 {
		ab21ap21View = "relay-twelve"
		ab21ap21Fields["relays"] = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		ab21ap21Fields["pwm2"] = []int{0}
	}
	if ab21ap21Operation == 2 {
		ab21ap21View = "mixer"
		ab21ap21Fields = map[string]any{"circuit": []int{0}, "mode": []int{1}, "setpoint": []int{190}}
	}
	if ab21ap21Direction == "positive-feedback" {
		ab21ap21Fields = map[string]any{"feedback": []int{1}}
		if ab21ap21Operation == 2 {
			ab21ap21Fields["movement"] = []int{0}
		}
	}
	return map[string]any{"context": map[string]any{"primary": 181, "secondary": 35, "operation": ab21ap21Operation, "direction": ab21ap21Direction, "view": ab21ap21View, "family": "synthetic-family", "variant": "opaque"}, "fields": ab21ap21Fields, "arguments": map[string]any{"variant": "two-circuit", "opaque": map[string]any{"false": false, "zero": 0, "null": nil, "array": []any{0, false, nil}}}}
}

func ab21ap21EvaluateFor(ab21ap21T *testing.T, ab21ap21Input any) map[string]any {
	ab21ap21T.Helper()
	ab21ap21Bytes, ab21ap21Error := json.Marshal(ab21ap21Input)
	if ab21ap21Error != nil {
		ab21ap21T.Fatal(ab21ap21Error)
	}
	ab21ap21Output, ab21ap21Error := Ab21Ap21Evaluate(ab21ap21Bytes)
	if ab21ap21Error != nil {
		ab21ap21T.Fatal(ab21ap21Error)
	}
	ab21ap21Decoded := map[string]any{}
	if ab21ap21Error = json.Unmarshal(ab21ap21Output, &ab21ap21Decoded); ab21ap21Error != nil {
		ab21ap21T.Fatal(ab21ap21Error)
	}
	if ab21ap21Decoded["physical"] != "Unknown" || ab21ap21Decoded["qualified"] != false {
		ab21ap21T.Fatal("native qualification")
	}
	ab21ap21Original := map[string]any{}
	_ = json.Unmarshal(ab21ap21Bytes, &ab21ap21Original)
	for _, ab21ap21Pair := range [][2]string{{"context", "context"}, {"fields", "raw"}, {"arguments", "arguments"}} {
		if !reflect.DeepEqual(ab21ap21Original[ab21ap21Pair[0]], ab21ap21Decoded[ab21ap21Pair[1]]) {
			ab21ap21T.Fatalf("lost %s", ab21ap21Pair[0])
		}
	}
	return ab21ap21Decoded
}

func ab21ap21FieldFor(ab21ap21Output map[string]any, ab21ap21Key string) map[string]any {
	return ab21ap21Output["fields"].(map[string]any)[ab21ap21Key].(map[string]any)
}

func TestAb21Ap21Relays(ab21ap21T *testing.T) {
	for _, ab21ap21Operation := range []int{1, 5} {
		ab21ap21Length := 6
		if ab21ap21Operation == 5 {
			ab21ap21Length = 12
		}
		for ab21ap21Slot := 0; ab21ap21Slot < ab21ap21Length; ab21ap21Slot++ {
			for _, ab21ap21Case := range []struct {
				ab21ap21Raw   int
				ab21ap21Label string
			}{{0, "off"}, {1, "on-for-seconds-then-off"}, {253, "on-for-seconds-then-off"}, {254, "on"}, {255, "unchanged"}} {
				ab21ap21T.Run(fmt.Sprintf("ab21ap21-%d-%d-%d", ab21ap21Operation, ab21ap21Slot, ab21ap21Case.ab21ap21Raw), func(ab21ap21T *testing.T) {
					ab21ap21Input := ab21ap21InputFor(ab21ap21Operation, "request")
					ab21ap21Input["fields"].(map[string]any)["relays"].([]int)[ab21ap21Slot] = ab21ap21Case.ab21ap21Raw
					ab21ap21Output := ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input)
					ab21ap21Relays := ab21ap21FieldFor(ab21ap21Output, "relays")["relays"].([]any)
					if len(ab21ap21Relays) != ab21ap21Length {
						ab21ap21T.Fatal("relay lost")
					}
					for ab21ap21Index, ab21ap21Value := range ab21ap21Relays {
						ab21ap21Relay := ab21ap21Value.(map[string]any)
						ab21ap21Want := "off"
						if ab21ap21Index == ab21ap21Slot {
							ab21ap21Want = ab21ap21Case.ab21ap21Label
						}
						if ab21ap21Relay["label"] != ab21ap21Want || ab21ap21Relay["grade"] != "Hypothesis" {
							ab21ap21T.Fatal(ab21ap21Relay)
						}
						if ab21ap21Index == ab21ap21Slot && (ab21ap21Case.ab21ap21Raw == 1 || ab21ap21Case.ab21ap21Raw == 253) {
							if ab21ap21Relay["number"] != float64(ab21ap21Case.ab21ap21Raw) {
								ab21ap21T.Fatal("duration")
							}
						} else if _, ab21ap21Exists := ab21ap21Relay["number"]; ab21ap21Exists {
							ab21ap21T.Fatal("invented duration")
						}
					}
				})
			}
		}
	}
}

func TestAb21Ap21PWM(ab21ap21T *testing.T) {
	for _, ab21ap21Operation := range []int{1, 5} {
		ab21ap21Keys := []string{"pwm1"}
		if ab21ap21Operation == 5 {
			ab21ap21Keys = append(ab21ap21Keys, "pwm2")
		}
		for _, ab21ap21Key := range ab21ap21Keys {
			for _, ab21ap21Case := range []struct {
				ab21ap21Raw      int
				ab21ap21Number   float64
				ab21ap21Domain   bool
				ab21ap21Behavior any
			}{{0, 0, true, float64(0)}, {200, 100, true, float64(100)}, {201, 100.5, true, "unchanged"}, {250, 125, true, "unchanged"}, {255, 127.5, false, "unchanged"}} {
				ab21ap21T.Run(fmt.Sprintf("ab21ap21-%d-%s-%d", ab21ap21Operation, ab21ap21Key, ab21ap21Case.ab21ap21Raw), func(ab21ap21T *testing.T) {
					ab21ap21Input := ab21ap21InputFor(ab21ap21Operation, "request")
					ab21ap21Input["fields"].(map[string]any)[ab21ap21Key] = []int{ab21ap21Case.ab21ap21Raw}
					ab21ap21Output := ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input)
					ab21ap21Alternatives := ab21ap21FieldFor(ab21ap21Output, ab21ap21Key)["alternatives"].([]any)
					if len(ab21ap21Alternatives) != 2 {
						ab21ap21T.Fatal("lost PWM alternative")
					}
					ab21ap21Numeric := ab21ap21Alternatives[0].(map[string]any)
					ab21ap21Behavior := ab21ap21Alternatives[1].(map[string]any)
					if ab21ap21Numeric["number"] != ab21ap21Case.ab21ap21Number || ab21ap21Numeric["candidate_domain"] != ab21ap21Case.ab21ap21Domain || ab21ap21Numeric["calculation"] != "Proven" || ab21ap21Numeric["grade"] != "Hypothesis" {
						ab21ap21T.Fatal(ab21ap21Numeric)
					}
					if ab21ap21Case.ab21ap21Raw <= 200 {
						if ab21ap21Behavior["number"] != ab21ap21Case.ab21ap21Behavior {
							ab21ap21T.Fatal(ab21ap21Behavior)
						}
					} else if ab21ap21Behavior["label"] != "unchanged" || ab21ap21Behavior["number"] != nil {
						ab21ap21T.Fatal(ab21ap21Behavior)
					}
					for _, ab21ap21Sibling := range ab21ap21Keys {
						if ab21ap21Sibling != ab21ap21Key {
							if ab21ap21FieldFor(ab21ap21Output, ab21ap21Sibling)["alternatives"].([]any)[0].(map[string]any)["number"] != float64(0) {
								ab21ap21T.Fatal("PWM sibling changed")
							}
						}
					}
				})
			}
		}
	}
}

func TestAb21Ap21Mixer(ab21ap21T *testing.T) {
	for _, ab21ap21Variant := range []string{"two-circuit", "three-circuit"} {
		for _, ab21ap21Circuit := range []int{0, 1, 2, 3, 255} {
			ab21ap21Input := ab21ap21InputFor(2, "request")
			ab21ap21Input["arguments"].(map[string]any)["variant"] = ab21ap21Variant
			ab21ap21Input["fields"].(map[string]any)["circuit"] = []int{ab21ap21Circuit}
			ab21ap21Want := "ignored"
			if ab21ap21Circuit == 0 {
				ab21ap21Want = "circuit-0"
			}
			if ab21ap21Circuit == 1 {
				ab21ap21Want = "circuit-1"
			}
			if ab21ap21Circuit == 2 && ab21ap21Variant == "three-circuit" {
				ab21ap21Want = "circuit-2"
			}
			if ab21ap21FieldFor(ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input), "circuit")["label"] != ab21ap21Want {
				ab21ap21T.Fatal("circuit variant")
			}
		}
	}
	for ab21ap21Mode, ab21ap21Label := range []string{"stop/pump-off/close", "heating-regulation", "cooling-regulation", "pump-off/open", "pump-on/close", "pump-on/open", "ignored"} {
		ab21ap21Input := ab21ap21InputFor(2, "request")
		ab21ap21Input["fields"].(map[string]any)["mode"] = []int{ab21ap21Mode}
		ab21ap21Output := ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input)
		if ab21ap21FieldFor(ab21ap21Output, "mode")["label"] != ab21ap21Label {
			ab21ap21T.Fatal("mode")
		}
		ab21ap21Field := ab21ap21FieldFor(ab21ap21Output, "setpoint")
		if ab21ap21Field["relevant"] != (ab21ap21Mode == 1 || ab21ap21Mode == 2) {
			ab21ap21T.Fatal("relevance")
		}
		ab21ap21Alternatives := ab21ap21Field["alternatives"].([]any)
		if len(ab21ap21Alternatives) != 2 || ab21ap21Alternatives[0].(map[string]any)["number"] != float64(-33) || ab21ap21Alternatives[1].(map[string]any)["number"] != float64(95) || ab21ap21Alternatives[1].(map[string]any)["clamp"] != float64(95) {
			ab21ap21T.Fatal(ab21ap21Field)
		}
	}
	for _, ab21ap21Case := range []struct {
		ab21ap21Raw    int
		ab21ap21Number float64
	}{{0, 0}, {127, 63.5}, {128, -64}, {255, -0.5}} {
		ab21ap21Input := ab21ap21InputFor(2, "request")
		ab21ap21Input["fields"].(map[string]any)["setpoint"] = []int{ab21ap21Case.ab21ap21Raw}
		ab21ap21Alternatives := ab21ap21FieldFor(ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input), "setpoint")["alternatives"].([]any)
		if ab21ap21Alternatives[0].(map[string]any)["number"] != ab21ap21Case.ab21ap21Number || ab21ap21Alternatives[1].(map[string]any)["state"] != "Unknown" || ab21ap21Alternatives[1].(map[string]any)["number"] != nil {
			ab21ap21T.Fatal(ab21ap21Alternatives)
		}
	}
}

func TestAb21Ap21Feedback(ab21ap21T *testing.T) {
	for _, ab21ap21Operation := range []int{1, 2, 5} {
		for _, ab21ap21Feedback := range []int{0, 1, 2, 255} {
			ab21ap21Input := ab21ap21InputFor(ab21ap21Operation, "positive-feedback")
			ab21ap21Input["fields"].(map[string]any)["feedback"] = []int{ab21ap21Feedback}
			ab21ap21Want := "rejected"
			if ab21ap21Operation == 5 {
				ab21ap21Want = "rejected-or-unsupported"
			}
			if ab21ap21Feedback == 1 {
				ab21ap21Want = "accepted"
			}
			if ab21ap21FieldFor(ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input), "feedback")["label"] != ab21ap21Want {
				ab21ap21T.Fatal("feedback operation")
			}
		}
	}
	for _, ab21ap21Case := range []struct {
		ab21ap21Raw    int
		ab21ap21Number any
		ab21ap21Label  any
	}{{0, float64(0), "stopped"}, {156, float64(-100), "closing"}, {100, float64(100), "opening"}, {155, nil, nil}, {101, nil, nil}, {255, float64(-1), "closing"}, {1, float64(1), "opening"}} {
		ab21ap21Input := ab21ap21InputFor(2, "positive-feedback")
		ab21ap21Input["fields"].(map[string]any)["movement"] = []int{ab21ap21Case.ab21ap21Raw}
		ab21ap21Field := ab21ap21FieldFor(ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input), "movement")
		if ab21ap21Field["number"] != ab21ap21Case.ab21ap21Number || ab21ap21Field["label"] != ab21ap21Case.ab21ap21Label {
			ab21ap21T.Fatal(ab21ap21Field)
		}
		if ab21ap21Case.ab21ap21Number == nil && ab21ap21Field["state"] != "Unknown" {
			ab21ap21T.Fatal(ab21ap21Field)
		}
	}
}

func TestAb21Ap21Context(ab21ap21T *testing.T) {
	for _, ab21ap21Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, ab21ap21Value := range []any{nil, "ab21ap21-absent", 0} {
			ab21ap21Input := ab21ap21InputFor(1, "request")
			ab21ap21Context := ab21ap21Input["context"].(map[string]any)
			if ab21ap21Value == "ab21ap21-absent" {
				delete(ab21ap21Context, ab21ap21Key)
			} else if ab21ap21Key == "direction" || ab21ap21Key == "view" {
				if ab21ap21Value == nil {
					ab21ap21Context[ab21ap21Key] = nil
				} else {
					ab21ap21Context[ab21ap21Key] = "unsupported"
				}
			} else {
				ab21ap21Context[ab21ap21Key] = ab21ap21Value
			}
			ab21ap21Output := ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input)
			if ab21ap21Output["state"] != "Unknown" || len(ab21ap21Output["fields"].(map[string]any)) != 0 {
				ab21ap21T.Fatal("unsupported interpreted")
			}
		}
	}
	for _, ab21ap21Key := range []string{"family", "variant", "opaque"} {
		for _, ab21ap21Value := range []any{nil, 0, false, "synthetic-other", []any{}, map[string]any{"nested": false}} {
			ab21ap21Input := ab21ap21InputFor(1, "request")
			ab21ap21Input["context"].(map[string]any)[ab21ap21Key] = ab21ap21Value
			if ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input)["state"] != "complete" {
				ab21ap21T.Fatal("opaque selector")
			}
		}
		ab21ap21Input := ab21ap21InputFor(1, "request")
		delete(ab21ap21Input["context"].(map[string]any), ab21ap21Key)
		ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input)
	}
}

func TestAb21Ap21Partial(ab21ap21T *testing.T) {
	for _, ab21ap21Operation := range []int{1, 2, 5} {
		for _, ab21ap21Direction := range []string{"request", "positive-feedback"} {
			ab21ap21Control := ab21ap21InputFor(ab21ap21Operation, ab21ap21Direction)
			for ab21ap21Key := range ab21ap21Control["fields"].(map[string]any) {
				for _, ab21ap21Length := range []int{-1, 0, 2, 13} {
					ab21ap21Input := ab21ap21InputFor(ab21ap21Operation, ab21ap21Direction)
					ab21ap21Raw := ab21ap21Input["fields"].(map[string]any)
					if ab21ap21Length == -1 {
						delete(ab21ap21Raw, ab21ap21Key)
					} else {
						ab21ap21Raw[ab21ap21Key] = make([]int, ab21ap21Length)
					}
					ab21ap21Output := ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input)
					if ab21ap21Output["state"] != "partial" || ab21ap21FieldFor(ab21ap21Output, ab21ap21Key)["state"] != "incomplete" {
						ab21ap21T.Fatal("incomplete lost")
					}
					for ab21ap21Sibling := range ab21ap21Control["fields"].(map[string]any) {
						if ab21ap21Sibling != ab21ap21Key && ab21ap21FieldFor(ab21ap21Output, ab21ap21Sibling)["state"] != "complete" {
							ab21ap21T.Fatal("independent sibling lost")
						}
					}
				}
			}
		}
	}
}

func TestAb21Ap21Reject(ab21ap21T *testing.T) {
	for _, ab21ap21Input := range []string{`{`, `[]`, `null`, `{} {}`, `{"context":{},"context":{},"fields":{}}`, `{"context":{"x":{"a":0,"a":1}},"fields":{}}`, `{"context":{},"fields":{},"arguments":{"x":[{"a":0,"a":1}]}}`, `{"context":null,"fields":{}}`, `{"context":{},"fields":null}`, `{"context":{},"fields":{},"arguments":null}`} {
		if ab21ap21Output, ab21ap21Error := Ab21Ap21Evaluate([]byte(ab21ap21Input)); ab21ap21Error == nil || ab21ap21Output != nil {
			ab21ap21T.Fatalf("accepted %s", ab21ap21Input)
		}
	}
	for _, ab21ap21Operation := range []int{1, 2, 5} {
		for _, ab21ap21Direction := range []string{"request", "positive-feedback"} {
			ab21ap21Control := ab21ap21InputFor(ab21ap21Operation, ab21ap21Direction)
			for ab21ap21Key := range ab21ap21Control["fields"].(map[string]any) {
				for _, ab21ap21Bad := range []any{nil, false, 0, "0", []any{false}, []any{nil}, []any{"1"}, []any{-1}, []any{256}, []any{1.5}, []any{json.Number("1.00000000000000000001")}} {
					ab21ap21Input := ab21ap21InputFor(ab21ap21Operation, ab21ap21Direction)
					ab21ap21Input["fields"].(map[string]any)[ab21ap21Key] = ab21ap21Bad
					ab21ap21Bytes, _ := json.Marshal(ab21ap21Input)
					if ab21ap21Output, ab21ap21Error := Ab21Ap21Evaluate(ab21ap21Bytes); ab21ap21Error == nil || ab21ap21Output != nil {
						ab21ap21T.Fatalf("accepted %s", ab21ap21Bytes)
					}
				}
			}
		}
	}
	for _, ab21ap21Bad := range []any{nil, false, 0, []any{}, "unknown"} {
		ab21ap21Input := ab21ap21InputFor(2, "request")
		ab21ap21Input["arguments"].(map[string]any)["variant"] = ab21ap21Bad
		ab21ap21Bytes, _ := json.Marshal(ab21ap21Input)
		if _, ab21ap21Error := Ab21Ap21Evaluate(ab21ap21Bytes); ab21ap21Error == nil {
			ab21ap21T.Fatal("bad variant")
		}
	}
	for _, ab21ap21RemoveArguments := range []bool{true, false} {
		ab21ap21Input := ab21ap21InputFor(2, "request")
		if ab21ap21RemoveArguments {
			delete(ab21ap21Input, "arguments")
		} else {
			delete(ab21ap21Input["arguments"].(map[string]any), "variant")
		}
		ab21ap21Bytes, _ := json.Marshal(ab21ap21Input)
		if _, ab21ap21Error := Ab21Ap21Evaluate(ab21ap21Bytes); ab21ap21Error == nil {
			ab21ap21T.Fatal("missing variant")
		}
	}
}

func TestAb21Ap21Isolation(ab21ap21T *testing.T) {
	ab21ap21Input, _ := json.Marshal(ab21ap21InputFor(5, "request"))
	ab21ap21Original := append([]byte(nil), ab21ap21Input...)
	ab21ap21First, ab21ap21Error := Ab21Ap21Evaluate(ab21ap21Input)
	if ab21ap21Error != nil {
		ab21ap21T.Fatal(ab21ap21Error)
	}
	ab21ap21Expected := append([]byte(nil), ab21ap21First...)
	for ab21ap21Index := range ab21ap21Input {
		ab21ap21Input[ab21ap21Index] = ' '
	}
	if !bytes.Equal(ab21ap21First, ab21ap21Expected) {
		ab21ap21T.Fatal("input alias")
	}
	ab21ap21Decoded := map[string]any{}
	_ = json.Unmarshal(ab21ap21First, &ab21ap21Decoded)
	ab21ap21Decoded["raw"].(map[string]any)["pwm1"].([]any)[0] = float64(201)
	ab21ap21Decoded["arguments"].(map[string]any)["opaque"].(map[string]any)["false"] = true
	if ab21ap21FieldFor(ab21ap21Decoded, "pwm1")["raw"].([]any)[0] != float64(0) {
		ab21ap21T.Fatal("nested raw alias")
	}
	for ab21ap21Index := range ab21ap21First {
		ab21ap21First[ab21ap21Index] = ' '
	}
	ab21ap21Second, ab21ap21Error := Ab21Ap21Evaluate(ab21ap21Original)
	if ab21ap21Error != nil || !bytes.Equal(ab21ap21Second, ab21ap21Expected) {
		ab21ap21T.Fatal("result alias")
	}
}

func TestAb21Ap21SyntheticFixture(ab21ap21T *testing.T) {
	ab21ap21Bytes, ab21ap21Error := os.ReadFile("testdata/ab21ap21_b523_actuator_fields.json")
	if ab21ap21Error != nil {
		ab21ap21T.Fatal(ab21ap21Error)
	}
	var ab21ap21Cases []struct {
		Ab21Ap21ID    string          `json:"id"`
		Ab21Ap21Input json.RawMessage `json:"input"`
		Ab21Ap21State string          `json:"state"`
	}
	if ab21ap21Error = json.Unmarshal(ab21ap21Bytes, &ab21ap21Cases); ab21ap21Error != nil || len(ab21ap21Cases) == 0 {
		ab21ap21T.Fatal("empty fixture", ab21ap21Error)
	}
	for _, ab21ap21Case := range ab21ap21Cases {
		ab21ap21T.Run(ab21ap21Case.Ab21Ap21ID, func(ab21ap21T *testing.T) {
			ab21ap21Output := ab21ap21EvaluateFor(ab21ap21T, ab21ap21Case.Ab21Ap21Input)
			if ab21ap21Output["state"] != ab21ap21Case.Ab21Ap21State {
				ab21ap21T.Fatal(ab21ap21Output)
			}
		})
	}
}

func TestAb21Ap21ValidationBoundaries(ab21ap21T *testing.T) {
	for _, ab21ap21Key := range []string{"primary", "secondary", "operation", "direction", "view"} {
		for _, ab21ap21Bad := range []any{false, []any{}, map[string]any{}, 1.5, -1, 256} {
			ab21ap21Input := ab21ap21InputFor(1, "request")
			ab21ap21Input["context"].(map[string]any)[ab21ap21Key] = ab21ap21Bad
			ab21ap21Bytes, _ := json.Marshal(ab21ap21Input)
			if ab21ap21Output, ab21ap21Error := Ab21Ap21Evaluate(ab21ap21Bytes); ab21ap21Error == nil || ab21ap21Output != nil {
				ab21ap21T.Fatalf("bad context accepted %s", ab21ap21Bytes)
			}
		}
	}
	for _, ab21ap21Key := range []string{"primary", "secondary", "operation"} {
		ab21ap21Input := ab21ap21InputFor(1, "request")
		ab21ap21Input["context"].(map[string]any)[ab21ap21Key] = "1"
		ab21ap21Bytes, _ := json.Marshal(ab21ap21Input)
		if _, ab21ap21Error := Ab21Ap21Evaluate(ab21ap21Bytes); ab21ap21Error == nil {
			ab21ap21T.Fatal("coerced context")
		}
	}
	for _, ab21ap21Input := range []string{`{"context":{},"fields":{"x":[],"x":[]}}`, `{"context":{},"fields":{},"arguments":{"variant":"two-circuit","variant":"three-circuit"}}`, `{"context":{"x":[{"y":1,"y":2}]},"fields":{}}`, `{"context":{},"fields":{"x":[1e999]}}`} {
		if ab21ap21Output, ab21ap21Error := Ab21Ap21Evaluate([]byte(ab21ap21Input)); ab21ap21Error == nil || ab21ap21Output != nil {
			ab21ap21T.Fatal("invalid JSON accepted")
		}
	}
	for _, ab21ap21Arguments := range []any{map[string]any{}, map[string]any{"extra": false}, map[string]any{"extra": nil}, map[string]any{"extra": 0}, map[string]any{"extra": []any{}}} {
		ab21ap21Input := ab21ap21InputFor(1, "request")
		ab21ap21Input["arguments"] = ab21ap21Arguments
		ab21ap21Input["fields"].(map[string]any)["extra"] = []int{0, 255}
		ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input)
	}
	ab21ap21Input := ab21ap21InputFor(1, "request")
	delete(ab21ap21Input, "arguments")
	ab21ap21Output := ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input)
	if _, ab21ap21Present := ab21ap21Output["arguments"]; ab21ap21Present {
		ab21ap21T.Fatal("invented arguments")
	}
	ab21ap21Input = ab21ap21InputFor(2, "request")
	delete(ab21ap21Input["fields"].(map[string]any), "mode")
	ab21ap21Field := ab21ap21FieldFor(ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input), "setpoint")
	if ab21ap21Field["relevance"] != "Unknown" || ab21ap21Field["relevant"] != nil {
		ab21ap21T.Fatal("invented relevance")
	}
	for _, ab21ap21Movement := range []int{156, 100} {
		ab21ap21Input := ab21ap21InputFor(2, "positive-feedback")
		ab21ap21Input["fields"].(map[string]any)["movement"] = []int{ab21ap21Movement}
		ab21ap21Field := ab21ap21FieldFor(ab21ap21EvaluateFor(ab21ap21T, ab21ap21Input), "movement")
		if ab21ap21Field["maximum"] != true || ab21ap21Field["grade"] != "Hypothesis" || ab21ap21Field["calculation"] != "Proven" {
			ab21ap21T.Fatal("movement grade or maximum")
		}
	}
}
