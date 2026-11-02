package vaillant

import (
	ab21aq21bytes "bytes"
	ab21aq21json "encoding/json"
	ab21aq21fmt "fmt"
	ab21aq21io "io"
	ab21aq21math "math"
	ab21aq21big "math/big"
	ab21aq21strconv "strconv"
)

func ab21aq21Read(ab21aq21Decoder *ab21aq21json.Decoder) (any, error) {
	ab21aq21Token, ab21aq21Err := ab21aq21Decoder.Token()
	if ab21aq21Err != nil {
		return nil, ab21aq21Err
	}
	switch ab21aq21Token {
	case ab21aq21json.Delim('{'):
		ab21aq21Object := map[string]any{}
		for ab21aq21Decoder.More() {
			ab21aq21KeyToken, ab21aq21Err := ab21aq21Decoder.Token()
			if ab21aq21Err != nil {
				return nil, ab21aq21Err
			}
			ab21aq21Key := ab21aq21KeyToken.(string)
			if _, ab21aq21Exists := ab21aq21Object[ab21aq21Key]; ab21aq21Exists {
				return nil, ab21aq21fmt.Errorf("duplicate key %q", ab21aq21Key)
			}
			ab21aq21Value, ab21aq21Err := ab21aq21Read(ab21aq21Decoder)
			if ab21aq21Err != nil {
				return nil, ab21aq21Err
			}
			ab21aq21Object[ab21aq21Key] = ab21aq21Value
		}
		_, ab21aq21Err = ab21aq21Decoder.Token()
		return ab21aq21Object, ab21aq21Err
	case ab21aq21json.Delim('['):
		ab21aq21Array := []any{}
		for ab21aq21Decoder.More() {
			ab21aq21Value, ab21aq21Err := ab21aq21Read(ab21aq21Decoder)
			if ab21aq21Err != nil {
				return nil, ab21aq21Err
			}
			ab21aq21Array = append(ab21aq21Array, ab21aq21Value)
		}
		_, ab21aq21Err = ab21aq21Decoder.Token()
		return ab21aq21Array, ab21aq21Err
	}
	if ab21aq21Number, ab21aq21OK := ab21aq21Token.(ab21aq21json.Number); ab21aq21OK {
		ab21aq21Float, ab21aq21Err := ab21aq21strconv.ParseFloat(string(ab21aq21Number), 64)
		if ab21aq21Err != nil || ab21aq21math.IsInf(ab21aq21Float, 0) || ab21aq21math.IsNaN(ab21aq21Float) {
			return nil, ab21aq21fmt.Errorf("nonfinite number")
		}
	}
	return ab21aq21Token, nil
}

func ab21aq21Uint(ab21aq21Value any) (uint64, bool) {
	ab21aq21Number, ab21aq21OK := ab21aq21Value.(ab21aq21json.Number)
	if !ab21aq21OK {
		return 0, false
	}
	ab21aq21Rat, ab21aq21OK := new(ab21aq21big.Rat).SetString(string(ab21aq21Number))
	if !ab21aq21OK || !ab21aq21Rat.IsInt() || ab21aq21Rat.Sign() < 0 || !ab21aq21Rat.Num().IsUint64() || ab21aq21Rat.Num().Uint64() > 255 {
		return 0, false
	}
	return ab21aq21Rat.Num().Uint64(), true
}

func Ab21Aq21Evaluate(ab21aq21Input ab21aq21json.RawMessage) (ab21aq21json.RawMessage, error) {
	ab21aq21Decoder := ab21aq21json.NewDecoder(ab21aq21bytes.NewReader(ab21aq21Input))
	ab21aq21Decoder.UseNumber()
	ab21aq21Value, ab21aq21Err := ab21aq21Read(ab21aq21Decoder)
	if ab21aq21Err != nil {
		return nil, ab21aq21Err
	}
	if _, ab21aq21Err = ab21aq21Decoder.Token(); ab21aq21Err != ab21aq21io.EOF {
		return nil, ab21aq21fmt.Errorf("trailing JSON")
	}
	ab21aq21Object, ab21aq21OK := ab21aq21Value.(map[string]any)
	if !ab21aq21OK {
		return nil, ab21aq21fmt.Errorf("input must be object")
	}
	ab21aq21Context, ab21aq21OK := ab21aq21Object["context"].(map[string]any)
	if !ab21aq21OK {
		return nil, ab21aq21fmt.Errorf("context must be object")
	}
	ab21aq21Fields, ab21aq21OK := ab21aq21Object["fields"].(map[string]any)
	if !ab21aq21OK {
		return nil, ab21aq21fmt.Errorf("fields must be object")
	}
	for _, ab21aq21Key := range []string{"primary", "secondary", "operation"} {
		if ab21aq21Value, ab21aq21Present := ab21aq21Context[ab21aq21Key]; ab21aq21Present {
			if _, ab21aq21OK := ab21aq21Uint(ab21aq21Value); !ab21aq21OK {
				return nil, ab21aq21fmt.Errorf("invalid context %s", ab21aq21Key)
			}
		}
	}
	for _, ab21aq21Key := range []string{"direction", "view"} {
		if ab21aq21Value, ab21aq21Present := ab21aq21Context[ab21aq21Key]; ab21aq21Present {
			if _, ab21aq21OK := ab21aq21Value.(string); !ab21aq21OK {
				return nil, ab21aq21fmt.Errorf("invalid context %s", ab21aq21Key)
			}
		}
	}
	for ab21aq21Key, ab21aq21Value := range ab21aq21Fields {
		ab21aq21Array, ab21aq21OK := ab21aq21Value.([]any)
		if !ab21aq21OK {
			return nil, ab21aq21fmt.Errorf("field %s must be byte array", ab21aq21Key)
		}
		for _, ab21aq21Byte := range ab21aq21Array {
			if _, ab21aq21OK := ab21aq21Uint(ab21aq21Byte); !ab21aq21OK {
				return nil, ab21aq21fmt.Errorf("invalid byte in %s", ab21aq21Key)
			}
		}
	}
	ab21aq21Arguments := map[string]any{}
	if ab21aq21Value, ab21aq21Present := ab21aq21Object["arguments"]; ab21aq21Present {
		ab21aq21Arguments, ab21aq21OK = ab21aq21Value.(map[string]any)
		if !ab21aq21OK {
			return nil, ab21aq21fmt.Errorf("arguments must be object")
		}
	}
	if ab21aq21Value, ab21aq21Present := ab21aq21Arguments["value_width"]; ab21aq21Present {
		ab21aq21Width, ab21aq21OK := ab21aq21Uint(ab21aq21Value)
		if !ab21aq21OK || (ab21aq21Width != 1 && ab21aq21Width != 2) {
			return nil, ab21aq21fmt.Errorf("value_width must be 1 or 2")
		}
	}
	ab21aq21Results := map[string]any{}
	ab21aq21Output := map[string]any{}
	for ab21aq21Key, ab21aq21Value := range ab21aq21Object {
		ab21aq21Output[ab21aq21Key] = ab21aq21Value
	}
	ab21aq21Output["input"] = ab21aq21Object
	ab21aq21Output["raw"] = ab21aq21Fields
	ab21aq21Output["results"] = ab21aq21Results
	ab21aq21Output["state"] = "Unknown"
	ab21aq21Output["physical"] = "Unknown"
	ab21aq21Output["qualified"] = false
	ab21aq21Primary, _ := ab21aq21Uint(ab21aq21Context["primary"])
	ab21aq21Secondary, _ := ab21aq21Uint(ab21aq21Context["secondary"])
	ab21aq21Operation, _ := ab21aq21Uint(ab21aq21Context["operation"])
	ab21aq21View, _ := ab21aq21Context["view"].(string)
	ab21aq21Supported := (ab21aq21Operation == 4 && ab21aq21View == "enable-bits") || (ab21aq21Operation == 5 && (ab21aq21View == "range" || ab21aq21View == "enum" || ab21aq21View == "four-byte")) || (ab21aq21Operation == 6 && ab21aq21View == "save")
	if ab21aq21Primary != 181 || ab21aq21Secondary != 26 || ab21aq21Context["direction"] != "positive" || !ab21aq21Supported {
		return ab21aq21json.Marshal(ab21aq21Output)
	}
	ab21aq21Widths := map[string]int{}
	switch ab21aq21View {
	case "enable-bits":
		ab21aq21Widths["enable"] = 13
	case "range":
		ab21aq21Widths = map[string]int{"flags": 1, "display_type": 1, "value": 2, "minimum": 2, "maximum": 2, "step": 1}
	case "enum":
		ab21aq21Widths = map[string]int{"flags": 1, "alternatives": 6}
	case "four-byte":
		ab21aq21Widths = map[string]int{"flags": 1, "value": 4, "reserved": 3}
	case "save":
		ab21aq21Width, ab21aq21OK := ab21aq21Uint(ab21aq21Arguments["value_width"])
		if !ab21aq21OK {
			return nil, ab21aq21fmt.Errorf("save requires value_width")
		}
		ab21aq21Widths = map[string]int{"status": 1, "mirror": 2, "symbol": 2, "value": int(ab21aq21Width)}
	}
	ab21aq21Output["state"] = "complete"
	ab21aq21FlagArray, _ := ab21aq21Fields["flags"].([]any)
	ab21aq21FlagNumber := uint64(0)
	ab21aq21FlagsKnown := len(ab21aq21FlagArray) == 1
	if ab21aq21FlagsKnown {
		ab21aq21FlagNumber, _ = ab21aq21Uint(ab21aq21FlagArray[0])
	}
	for ab21aq21Key, ab21aq21Width := range ab21aq21Widths {
		ab21aq21Result := map[string]any{"state": "incomplete"}
		ab21aq21Results[ab21aq21Key] = ab21aq21Result
		ab21aq21Raw, ab21aq21Present := ab21aq21Fields[ab21aq21Key]
		if ab21aq21Present {
			ab21aq21Result["raw"] = ab21aq21Raw
		}
		ab21aq21Array, _ := ab21aq21Raw.([]any)
		if !ab21aq21Present || len(ab21aq21Array) != ab21aq21Width {
			ab21aq21Output["state"] = "partial"
			continue
		}
		if ab21aq21Operation == 5 && ab21aq21Key != "flags" && (!ab21aq21FlagsKnown || ab21aq21FlagNumber&1 != 0) {
			ab21aq21Result["state"] = "Unknown"
			continue
		}
		ab21aq21Result["state"] = "Hypothesis"
		if ab21aq21Key == "enable" {
			ab21aq21Bits := make([]bool, 104)
			for ab21aq21Index, ab21aq21Byte := range ab21aq21Array {
				ab21aq21Number, _ := ab21aq21Uint(ab21aq21Byte)
				for ab21aq21Bit := 0; ab21aq21Bit < 8; ab21aq21Bit++ {
					ab21aq21Bits[ab21aq21Index*8+ab21aq21Bit] = ab21aq21Number&(1<<ab21aq21Bit) != 0
				}
			}
			ab21aq21Result["flags"] = ab21aq21Bits
			ab21aq21Result["calculation"] = "Proven"
			continue
		}
		if ab21aq21Key == "alternatives" {
			ab21aq21Result["alternatives"] = ab21aq21Array
			ab21aq21Result["calculation"] = "Proven"
			continue
		}
		if ab21aq21Key == "reserved" {
			ab21aq21Result["reserved"] = ab21aq21Array
			continue
		}
		ab21aq21Number := uint64(0)
		for ab21aq21Index, ab21aq21Byte := range ab21aq21Array {
			ab21aq21Part, _ := ab21aq21Uint(ab21aq21Byte)
			ab21aq21Number |= ab21aq21Part << (8 * ab21aq21Index)
		}
		ab21aq21Result["number"] = ab21aq21Number
		ab21aq21Result["calculation"] = "Proven"
		if ab21aq21Key == "flags" {
			ab21aq21Result["flags"] = map[string]bool{"disabled": ab21aq21Number&1 != 0, "writable": ab21aq21Number&2 != 0, "enum-versus-range": ab21aq21Number&4 != 0, "additional-four-byte/no-range": ab21aq21Number&8 != 0}
			ab21aq21Result["reserved"] = (ab21aq21Number >> 4) & 3
			ab21aq21Result["unmodeled"] = (ab21aq21Number >> 6) & 3
			ab21aq21Result["candidate_domain"] = "text alternatives"
			if ab21aq21Number&8 != 0 {
				ab21aq21Result["alternatives"] = []any{map[string]any{"label": "additional-four-byte/no-range", "state": "Hypothesis"}, map[string]any{"label": "read-only-overrides-writable", "state": "Hypothesis"}}
			}
		}
		if ab21aq21Key == "status" {
			ab21aq21Labels := []string{"value-not-saved", "point-not-found", "write-protected", "range-violation", "misconfiguration"}
			ab21aq21Result["label"] = "unmapped"
			if ab21aq21Number < uint64(len(ab21aq21Labels)) {
				ab21aq21Result["label"] = ab21aq21Labels[ab21aq21Number]
			}
		}
		if ab21aq21View == "save" && ab21aq21Key == "value" {
			ab21aq21Result["candidate_domain"] = "caller-selected width"
			ab21aq21Result["alternatives"] = []int{1, 2}
		}
	}
	return ab21aq21json.Marshal(ab21aq21Output)
}
