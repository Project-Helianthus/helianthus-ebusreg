package vaillant

import (
	ab19ap19bytes "bytes"
	ab19ap19json "encoding/json"
	ab19ap19fmt "fmt"
	ab19ap19io "io"
	ab19ap19big "math/big"
)

func ab19ap19Read(ab19ap19Decoder *ab19ap19json.Decoder) (any, error) {
	ab19ap19Token, ab19ap19Err := ab19ap19Decoder.Token()
	if ab19ap19Err != nil {
		return nil, ab19ap19Err
	}
	switch ab19ap19Token {
	case ab19ap19json.Delim('{'):
		ab19ap19Object := map[string]any{}
		for ab19ap19Decoder.More() {
			ab19ap19KeyToken, ab19ap19Err := ab19ap19Decoder.Token()
			if ab19ap19Err != nil {
				return nil, ab19ap19Err
			}
			ab19ap19Key := ab19ap19KeyToken.(string)
			if _, ab19ap19Exists := ab19ap19Object[ab19ap19Key]; ab19ap19Exists {
				return nil, ab19ap19fmt.Errorf("duplicate key %q", ab19ap19Key)
			}
			ab19ap19Value, ab19ap19Err := ab19ap19Read(ab19ap19Decoder)
			if ab19ap19Err != nil {
				return nil, ab19ap19Err
			}
			ab19ap19Object[ab19ap19Key] = ab19ap19Value
		}
		_, ab19ap19Err = ab19ap19Decoder.Token()
		return ab19ap19Object, ab19ap19Err
	case ab19ap19json.Delim('['):
		ab19ap19Array := []any{}
		for ab19ap19Decoder.More() {
			ab19ap19Value, ab19ap19Err := ab19ap19Read(ab19ap19Decoder)
			if ab19ap19Err != nil {
				return nil, ab19ap19Err
			}
			ab19ap19Array = append(ab19ap19Array, ab19ap19Value)
		}
		_, ab19ap19Err = ab19ap19Decoder.Token()
		return ab19ap19Array, ab19ap19Err
	}
	return ab19ap19Token, nil
}

func ab19ap19Uint(ab19ap19Value any, ab19ap19Max float64) (uint64, bool) {
	ab19ap19Number, ab19ap19OK := ab19ap19Value.(ab19ap19json.Number)
	if !ab19ap19OK {
		return 0, false
	}
	ab19ap19Rational, ab19ap19OK := new(ab19ap19big.Rat).SetString(string(ab19ap19Number))
	if !ab19ap19OK || !ab19ap19Rational.IsInt() || ab19ap19Rational.Sign() < 0 || !ab19ap19Rational.Num().IsUint64() || ab19ap19Rational.Num().Uint64() > uint64(ab19ap19Max) {
		return 0, false
	}
	return ab19ap19Rational.Num().Uint64(), true
}

func Ab19Ap19Evaluate(ab19ap19Input ab19ap19json.RawMessage) (ab19ap19json.RawMessage, error) {
	ab19ap19Decoder := ab19ap19json.NewDecoder(ab19ap19bytes.NewReader(ab19ap19Input))
	ab19ap19Decoder.UseNumber()
	ab19ap19Value, ab19ap19Err := ab19ap19Read(ab19ap19Decoder)
	if ab19ap19Err != nil {
		return nil, ab19ap19Err
	}
	if _, ab19ap19Err = ab19ap19Decoder.Token(); ab19ap19Err != ab19ap19io.EOF {
		return nil, ab19ap19fmt.Errorf("trailing JSON")
	}
	ab19ap19Object, ab19ap19OK := ab19ap19Value.(map[string]any)
	if !ab19ap19OK {
		return nil, ab19ap19fmt.Errorf("input must be object")
	}
	ab19ap19Context, ab19ap19OK := ab19ap19Object["context"].(map[string]any)
	if !ab19ap19OK {
		return nil, ab19ap19fmt.Errorf("context must be object")
	}
	ab19ap19Raw, ab19ap19OK := ab19ap19Object["fields"].(map[string]any)
	if !ab19ap19OK {
		return nil, ab19ap19fmt.Errorf("fields must be object")
	}
	for _, ab19ap19Key := range []string{"primary", "secondary", "operation"} {
		if ab19ap19Value, ab19ap19Present := ab19ap19Context[ab19ap19Key]; ab19ap19Present && ab19ap19Value != nil {
			if _, ab19ap19OK := ab19ap19Uint(ab19ap19Value, 255); !ab19ap19OK {
				return nil, ab19ap19fmt.Errorf("invalid context %s", ab19ap19Key)
			}
		}
	}
	for _, ab19ap19Key := range []string{"direction", "view"} {
		if ab19ap19Value, ab19ap19Present := ab19ap19Context[ab19ap19Key]; ab19ap19Present && ab19ap19Value != nil {
			if _, ab19ap19OK := ab19ap19Value.(string); !ab19ap19OK {
				return nil, ab19ap19fmt.Errorf("invalid context %s", ab19ap19Key)
			}
		}
	}
	for ab19ap19Key, ab19ap19Value := range ab19ap19Raw {
		ab19ap19Array, ab19ap19OK := ab19ap19Value.([]any)
		if !ab19ap19OK {
			return nil, ab19ap19fmt.Errorf("field %s must be byte array", ab19ap19Key)
		}
		for _, ab19ap19Byte := range ab19ap19Array {
			if _, ab19ap19OK := ab19ap19Uint(ab19ap19Byte, 255); !ab19ap19OK {
				return nil, ab19ap19fmt.Errorf("invalid byte in %s", ab19ap19Key)
			}
		}
	}
	ab19ap19Output := map[string]any{"context": ab19ap19Context, "raw": ab19ap19Raw, "fields": map[string]any{}, "state": "Unknown", "physical": "Unknown", "qualified": false}
	if ab19ap19Arguments, ab19ap19Present := ab19ap19Object["arguments"]; ab19ap19Present {
		if _, ab19ap19OK := ab19ap19Arguments.(map[string]any); !ab19ap19OK {
			return nil, ab19ap19fmt.Errorf("arguments must be object")
		}
		ab19ap19Output["arguments"] = ab19ap19Arguments
	}
	ab19ap19Primary, _ := ab19ap19Uint(ab19ap19Context["primary"], 255)
	ab19ap19Secondary, _ := ab19ap19Uint(ab19ap19Context["secondary"], 255)
	ab19ap19Operation, _ := ab19ap19Uint(ab19ap19Context["operation"], 255)
	ab19ap19Direction := ab19ap19Context["direction"]
	ab19ap19View := ab19ap19Context["view"]
	if ab19ap19Primary != 181 || ab19ap19Secondary != 26 || ((ab19ap19Operation != 2 || ab19ap19View != "mirror-get") && (ab19ap19Operation != 3 || ab19ap19View != "mirror-put")) || (ab19ap19Direction != "request" && ab19ap19Direction != "positive") {
		return ab19ap19json.Marshal(ab19ap19Output)
	}
	ab19ap19Widths := map[string]int{"sequence": 1}
	if ab19ap19Direction == "request" {
		ab19ap19Widths["flags"] = 1
		if ab19ap19Operation == 2 {
			ab19ap19Widths["length"] = 1
			ab19ap19Widths["address"] = 2
		} else {
			ab19ap19Widths["mirror_address"] = 2
			ab19ap19Widths["symbol_address"] = 2
			ab19ap19Widths["data"] = -1
		}
	} else if ab19ap19Operation == 2 {
		ab19ap19Widths["status"] = 1
		ab19ap19Widths["data"] = -1
	} else {
		ab19ap19Widths["flags"] = 1
	}
	ab19ap19Results := ab19ap19Output["fields"].(map[string]any)
	ab19ap19Output["state"] = "complete"
	for ab19ap19Key, ab19ap19Width := range ab19ap19Widths {
		ab19ap19Result := map[string]any{"state": "incomplete"}
		ab19ap19Results[ab19ap19Key] = ab19ap19Result
		ab19ap19Value, ab19ap19Present := ab19ap19Raw[ab19ap19Key]
		if ab19ap19Present {
			ab19ap19Result["raw"] = ab19ap19Value
		}
		ab19ap19Array, _ := ab19ap19Value.([]any)
		if !ab19ap19Present || (ab19ap19Width >= 0 && len(ab19ap19Array) != ab19ap19Width) {
			ab19ap19Output["state"] = "partial"
			continue
		}
		ab19ap19Result["state"] = "Hypothesis"
		if ab19ap19Key == "data" {
			ab19ap19Result["label"] = "opaque"
			continue
		}
		ab19ap19Number, _ := ab19ap19Uint(ab19ap19Array[0], 255)
		if ab19ap19Width == 2 {
			ab19ap19High, _ := ab19ap19Uint(ab19ap19Array[1], 255)
			ab19ap19Number += 256 * ab19ap19High
			ab19ap19Result["logical_validity"] = "Unknown"
		}
		ab19ap19Result["number"] = ab19ap19Number
		ab19ap19Result["calculation"] = "Proven"
		if ab19ap19Key == "status" {
			ab19ap19Boolean := map[string]any{"state": "Unknown"}
			if ab19ap19Number <= 1 {
				ab19ap19Boolean["state"] = "Hypothesis"
				ab19ap19Boolean["value"] = ab19ap19Number == 1
			}
			ab19ap19Result["alternatives"] = map[string]any{"whole_boolean": ab19ap19Boolean, "bit0_out_of_area": map[string]any{"state": "Hypothesis", "value": ab19ap19Number&1 != 0}}
		}
		if ab19ap19Key == "flags" {
			ab19ap19Names := []string{"triple-storage"}
			if ab19ap19Operation == 3 {
				if ab19ap19Direction == "request" {
					ab19ap19Names = append(ab19ap19Names, "direct-update")
				} else {
					ab19ap19Names = []string{"not-written", "invalid-symbol", "invalid-parameter", "storage-error"}
				}
			}
			ab19ap19Flags := map[string]bool{}
			for ab19ap19Bit, ab19ap19Name := range ab19ap19Names {
				ab19ap19Flags[ab19ap19Name] = ab19ap19Number&(1<<ab19ap19Bit) != 0
			}
			ab19ap19Result["flags"] = ab19ap19Flags
			ab19ap19Rest := ab19ap19Number & (255 << len(ab19ap19Names))
			if ab19ap19Operation == 3 {
				ab19ap19Result["reserved"] = ab19ap19Rest
			} else {
				ab19ap19Result["unresolved"] = ab19ap19Rest
			}
		}
	}
	return ab19ap19json.Marshal(ab19ap19Output)
}
