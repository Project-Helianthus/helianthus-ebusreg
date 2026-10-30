package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
)

func ab21ap21Read(ab21ap21Decoder *json.Decoder) (any, error) {
	ab21ap21Token, ab21ap21Error := ab21ap21Decoder.Token()
	if ab21ap21Error != nil {
		return nil, ab21ap21Error
	}
	if ab21ap21Delimiter, ab21ap21OK := ab21ap21Token.(json.Delim); ab21ap21OK {
		switch ab21ap21Delimiter {
		case '{':
			ab21ap21Object := map[string]any{}
			for ab21ap21Decoder.More() {
				ab21ap21KeyToken, ab21ap21Error := ab21ap21Decoder.Token()
				if ab21ap21Error != nil {
					return nil, ab21ap21Error
				}
				ab21ap21Key := ab21ap21KeyToken.(string)
				if _, ab21ap21Exists := ab21ap21Object[ab21ap21Key]; ab21ap21Exists {
					return nil, fmt.Errorf("duplicate key %q", ab21ap21Key)
				}
				ab21ap21Value, ab21ap21Error := ab21ap21Read(ab21ap21Decoder)
				if ab21ap21Error != nil {
					return nil, ab21ap21Error
				}
				ab21ap21Object[ab21ap21Key] = ab21ap21Value
			}
			_, ab21ap21Error = ab21ap21Decoder.Token()
			return ab21ap21Object, ab21ap21Error
		case '[':
			ab21ap21Array := []any{}
			for ab21ap21Decoder.More() {
				ab21ap21Value, ab21ap21Error := ab21ap21Read(ab21ap21Decoder)
				if ab21ap21Error != nil {
					return nil, ab21ap21Error
				}
				ab21ap21Array = append(ab21ap21Array, ab21ap21Value)
			}
			_, ab21ap21Error = ab21ap21Decoder.Token()
			return ab21ap21Array, ab21ap21Error
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	return ab21ap21Token, nil
}

func ab21ap21Byte(ab21ap21Value any) (int, bool) {
	ab21ap21Number, ab21ap21OK := ab21ap21Value.(json.Number)
	if !ab21ap21OK {
		return 0, false
	}
	ab21ap21Rational, ab21ap21OK := new(big.Rat).SetString(string(ab21ap21Number))
	if !ab21ap21OK || !ab21ap21Rational.IsInt() || ab21ap21Rational.Sign() < 0 || ab21ap21Rational.Cmp(big.NewRat(255, 1)) > 0 {
		return 0, false
	}
	return int(ab21ap21Rational.Num().Int64()), true
}

func ab21ap21Result(ab21ap21Raw any) map[string]any {
	return map[string]any{"raw": ab21ap21Raw, "state": "complete", "grade": "Hypothesis"}
}

func Ab21Ap21Evaluate(ab21ap21Input json.RawMessage) (json.RawMessage, error) {
	ab21ap21Decoder := json.NewDecoder(bytes.NewReader(ab21ap21Input))
	ab21ap21Decoder.UseNumber()
	ab21ap21Value, ab21ap21Error := ab21ap21Read(ab21ap21Decoder)
	if ab21ap21Error != nil {
		return nil, ab21ap21Error
	}
	if _, ab21ap21Error = ab21ap21Decoder.Token(); ab21ap21Error != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	ab21ap21Object, ab21ap21OK := ab21ap21Value.(map[string]any)
	if !ab21ap21OK {
		return nil, fmt.Errorf("input must be an object")
	}
	ab21ap21Context, ab21ap21OK := ab21ap21Object["context"].(map[string]any)
	if !ab21ap21OK {
		return nil, fmt.Errorf("context must be an object")
	}
	ab21ap21Raw, ab21ap21OK := ab21ap21Object["fields"].(map[string]any)
	if !ab21ap21OK {
		return nil, fmt.Errorf("fields must be an object")
	}
	for _, ab21ap21Key := range []string{"primary", "secondary", "operation"} {
		if ab21ap21Value := ab21ap21Context[ab21ap21Key]; ab21ap21Value != nil {
			if _, ab21ap21OK := ab21ap21Byte(ab21ap21Value); !ab21ap21OK {
				return nil, fmt.Errorf("invalid context %s", ab21ap21Key)
			}
		}
	}
	for _, ab21ap21Key := range []string{"direction", "view"} {
		if ab21ap21Value := ab21ap21Context[ab21ap21Key]; ab21ap21Value != nil {
			if _, ab21ap21OK := ab21ap21Value.(string); !ab21ap21OK {
				return nil, fmt.Errorf("invalid context %s", ab21ap21Key)
			}
		}
	}
	for ab21ap21Key, ab21ap21Value := range ab21ap21Raw {
		ab21ap21Array, ab21ap21OK := ab21ap21Value.([]any)
		if !ab21ap21OK {
			return nil, fmt.Errorf("field %s must be a byte array", ab21ap21Key)
		}
		for _, ab21ap21Value := range ab21ap21Array {
			if _, ab21ap21OK := ab21ap21Byte(ab21ap21Value); !ab21ap21OK {
				return nil, fmt.Errorf("invalid byte in %s", ab21ap21Key)
			}
		}
	}
	ab21ap21Arguments := map[string]any{}
	if ab21ap21Value, ab21ap21Present := ab21ap21Object["arguments"]; ab21ap21Present {
		ab21ap21Arguments, ab21ap21OK = ab21ap21Value.(map[string]any)
		if !ab21ap21OK {
			return nil, fmt.Errorf("arguments must be an object")
		}
	}
	ab21ap21Fields := map[string]any{}
	ab21ap21Output := map[string]any{"context": ab21ap21Context, "raw": ab21ap21Raw, "fields": ab21ap21Fields, "state": "Unknown", "physical": "Unknown", "qualified": false}
	if _, ab21ap21Present := ab21ap21Object["arguments"]; ab21ap21Present {
		ab21ap21Output["arguments"] = ab21ap21Arguments
	}
	ab21ap21Primary, _ := ab21ap21Byte(ab21ap21Context["primary"])
	ab21ap21Secondary, _ := ab21ap21Byte(ab21ap21Context["secondary"])
	ab21ap21Operation, _ := ab21ap21Byte(ab21ap21Context["operation"])
	ab21ap21View, _ := ab21ap21Context["view"].(string)
	ab21ap21Direction, _ := ab21ap21Context["direction"].(string)
	ab21ap21Supported := ab21ap21Primary == 181 && ab21ap21Secondary == 35 && ((ab21ap21Operation == 1 && ab21ap21View == "relay-six") || (ab21ap21Operation == 5 && ab21ap21View == "relay-twelve") || (ab21ap21Operation == 2 && ab21ap21View == "mixer")) && (ab21ap21Direction == "request" || ab21ap21Direction == "positive-feedback")
	if !ab21ap21Supported {
		return json.Marshal(ab21ap21Output)
	}
	ab21ap21Variant := ""
	if ab21ap21Operation == 2 && ab21ap21Direction == "request" {
		ab21ap21Variant, ab21ap21OK = ab21ap21Arguments["variant"].(string)
		if !ab21ap21OK || (ab21ap21Variant != "two-circuit" && ab21ap21Variant != "three-circuit") {
			return nil, fmt.Errorf("variant must be two-circuit or three-circuit")
		}
	}
	ab21ap21Output["state"] = "complete"
	ab21ap21Get := func(ab21ap21Key string, ab21ap21Length int) ([]any, map[string]any, bool) {
		ab21ap21Value, ab21ap21Present := ab21ap21Raw[ab21ap21Key]
		ab21ap21Array, _ := ab21ap21Value.([]any)
		ab21ap21Field := ab21ap21Result(ab21ap21Value)
		ab21ap21Field["present"] = ab21ap21Present
		ab21ap21Fields[ab21ap21Key] = ab21ap21Field
		if !ab21ap21Present || len(ab21ap21Array) != ab21ap21Length {
			ab21ap21Field["state"] = "incomplete"
			ab21ap21Output["state"] = "partial"
			return ab21ap21Array, ab21ap21Field, false
		}
		return ab21ap21Array, ab21ap21Field, true
	}
	if ab21ap21Direction == "positive-feedback" {
		if ab21ap21Array, ab21ap21Field, ab21ap21OK := ab21ap21Get("feedback", 1); ab21ap21OK {
			ab21ap21ByteValue, _ := ab21ap21Byte(ab21ap21Array[0])
			ab21ap21Field["label"] = "rejected"
			if ab21ap21Operation == 5 {
				ab21ap21Field["label"] = "rejected-or-unsupported"
			}
			if ab21ap21ByteValue == 1 {
				ab21ap21Field["label"] = "accepted"
			}
		}
		if ab21ap21Operation == 2 {
			if ab21ap21Array, ab21ap21Field, ab21ap21OK := ab21ap21Get("movement", 1); ab21ap21OK {
				ab21ap21ByteValue, _ := ab21ap21Byte(ab21ap21Array[0])
				ab21ap21Signed := int(int8(ab21ap21ByteValue))
				if ab21ap21Signed >= -100 && ab21ap21Signed <= 100 {
					ab21ap21Field["number"] = ab21ap21Signed
					ab21ap21Field["unit"] = "percent"
					ab21ap21Field["calculation"] = "Proven"
					ab21ap21Field["maximum"] = ab21ap21Signed == -100 || ab21ap21Signed == 100
					ab21ap21Field["label"] = "stopped"
					if ab21ap21Signed < 0 {
						ab21ap21Field["label"] = "closing"
					}
					if ab21ap21Signed > 0 {
						ab21ap21Field["label"] = "opening"
					}
				} else {
					ab21ap21Field["state"] = "Unknown"
				}
			}
		}
	} else if ab21ap21Operation == 1 || ab21ap21Operation == 5 {
		ab21ap21Length := 6
		if ab21ap21Operation == 5 {
			ab21ap21Length = 12
		}
		if ab21ap21Array, ab21ap21Field, ab21ap21OK := ab21ap21Get("relays", ab21ap21Length); ab21ap21OK {
			ab21ap21Relays := []any{}
			for _, ab21ap21Value := range ab21ap21Array {
				ab21ap21ByteValue, _ := ab21ap21Byte(ab21ap21Value)
				ab21ap21Relay := ab21ap21Result(ab21ap21Value)
				switch ab21ap21ByteValue {
				case 0:
					ab21ap21Relay["label"] = "off"
				case 254:
					ab21ap21Relay["label"] = "on"
				case 255:
					ab21ap21Relay["label"] = "unchanged"
				default:
					ab21ap21Relay["label"] = "on-for-seconds-then-off"
					ab21ap21Relay["number"] = ab21ap21ByteValue
					ab21ap21Relay["unit"] = "seconds"
				}
				ab21ap21Relays = append(ab21ap21Relays, ab21ap21Relay)
			}
			ab21ap21Field["relays"] = ab21ap21Relays
		}
		ab21ap21PWMKeys := []string{"pwm1"}
		if ab21ap21Operation == 5 {
			ab21ap21PWMKeys = append(ab21ap21PWMKeys, "pwm2")
		}
		for _, ab21ap21Key := range ab21ap21PWMKeys {
			if ab21ap21Array, ab21ap21Field, ab21ap21OK := ab21ap21Get(ab21ap21Key, 1); ab21ap21OK {
				ab21ap21ByteValue, _ := ab21ap21Byte(ab21ap21Array[0])
				ab21ap21Numeric := ab21ap21Result(ab21ap21Array)
				ab21ap21Numeric["candidate"] = "numeric"
				ab21ap21Numeric["number"] = float64(ab21ap21ByteValue) / 2
				ab21ap21Numeric["unit"] = "percent"
				ab21ap21Numeric["candidate_domain"] = ab21ap21ByteValue <= 250
				ab21ap21Numeric["calculation"] = "Proven"
				ab21ap21Behavior := ab21ap21Result(ab21ap21Array)
				ab21ap21Behavior["candidate"] = "behavior"
				if ab21ap21ByteValue <= 200 {
					ab21ap21Behavior["number"] = float64(ab21ap21ByteValue) / 2
					ab21ap21Behavior["unit"] = "percent"
					ab21ap21Behavior["calculation"] = "Proven"
				} else {
					ab21ap21Behavior["label"] = "unchanged"
				}
				ab21ap21Field["alternatives"] = []any{ab21ap21Numeric, ab21ap21Behavior}
			}
		}
	} else {
		if ab21ap21Array, ab21ap21Field, ab21ap21OK := ab21ap21Get("circuit", 1); ab21ap21OK {
			ab21ap21ByteValue, _ := ab21ap21Byte(ab21ap21Array[0])
			if ab21ap21ByteValue < 2 || (ab21ap21ByteValue == 2 && ab21ap21Variant == "three-circuit") {
				ab21ap21Field["label"] = fmt.Sprintf("circuit-%d", ab21ap21ByteValue)
			} else {
				ab21ap21Field["label"] = "ignored"
			}
		}
		ab21ap21Mode := -1
		if ab21ap21Array, ab21ap21Field, ab21ap21OK := ab21ap21Get("mode", 1); ab21ap21OK {
			ab21ap21Mode, _ = ab21ap21Byte(ab21ap21Array[0])
			ab21ap21Labels := []string{"stop/pump-off/close", "heating-regulation", "cooling-regulation", "pump-off/open", "pump-on/close", "pump-on/open"}
			ab21ap21Field["label"] = "ignored"
			if ab21ap21Mode < len(ab21ap21Labels) {
				ab21ap21Field["label"] = ab21ap21Labels[ab21ap21Mode]
			}
		}
		if ab21ap21Array, ab21ap21Field, ab21ap21OK := ab21ap21Get("setpoint", 1); ab21ap21OK {
			if ab21ap21Mode >= 0 {
				ab21ap21Field["relevant"] = ab21ap21Mode == 1 || ab21ap21Mode == 2
			} else {
				ab21ap21Field["relevance"] = "Unknown"
			}
			ab21ap21ByteValue, _ := ab21ap21Byte(ab21ap21Array[0])
			ab21ap21Signed := ab21ap21Result(ab21ap21Array)
			ab21ap21Signed["candidate"] = "signed-half"
			ab21ap21Signed["number"] = float64(int8(ab21ap21ByteValue)) / 2
			ab21ap21Signed["unit"] = "Celsius"
			ab21ap21Signed["calculation"] = "Proven"
			ab21ap21Example := ab21ap21Result(ab21ap21Array)
			ab21ap21Example["candidate"] = "conflicting-example"
			ab21ap21Example["clamp"] = 95
			ab21ap21Example["state"] = "Unknown"
			if ab21ap21ByteValue == 190 {
				ab21ap21Example["number"] = 95
				ab21ap21Example["unit"] = "Celsius"
				ab21ap21Example["state"] = "complete"
			}
			ab21ap21Field["alternatives"] = []any{ab21ap21Signed, ab21ap21Example}
		}
	}
	return json.Marshal(ab21ap21Output)
}
