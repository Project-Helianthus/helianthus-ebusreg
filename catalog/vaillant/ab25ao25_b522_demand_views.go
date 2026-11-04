package vaillant

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

var ab25ao25IntegerPattern = regexp.MustCompile(`^-?(0|[1-9][0-9]*)$`)

const ab25ao25MaxDepth = 64

var ab25ao25ModeLabels = map[int]string{
	0: "no-heating-demand",
	1: "no-source-demand-sink-demand",
	2: "source-heating-demand",
	3: "no-source-cooling-demand-sink-demand",
	4: "source-cooling-demand",
}

func ab25ao25Parse(ab25ao25Dec *json.Decoder, ab25ao25Depth int) (any, error) {
	if ab25ao25Depth > ab25ao25MaxDepth {
		return nil, errors.New("ab25ao25: nesting too deep")
	}
	ab25ao25Tok, ab25ao25Err := ab25ao25Dec.Token()
	if ab25ao25Err != nil {
		return nil, ab25ao25Err
	}
	ab25ao25Delim, ab25ao25IsDelim := ab25ao25Tok.(json.Delim)
	if !ab25ao25IsDelim {
		return ab25ao25Tok, nil
	}
	if ab25ao25Delim == '{' {
		ab25ao25Map := map[string]any{}
		for ab25ao25Dec.More() {
			ab25ao25KeyTok, ab25ao25KeyErr := ab25ao25Dec.Token()
			if ab25ao25KeyErr != nil {
				return nil, ab25ao25KeyErr
			}
			ab25ao25Key, ab25ao25KeyOK := ab25ao25KeyTok.(string)
			if !ab25ao25KeyOK {
				return nil, errors.New("ab25ao25: object key is not a string")
			}
			if _, ab25ao25Dup := ab25ao25Map[ab25ao25Key]; ab25ao25Dup {
				return nil, fmt.Errorf("ab25ao25: duplicate key %q", ab25ao25Key)
			}
			ab25ao25Val, ab25ao25ValErr := ab25ao25Parse(ab25ao25Dec, ab25ao25Depth+1)
			if ab25ao25ValErr != nil {
				return nil, ab25ao25ValErr
			}
			ab25ao25Map[ab25ao25Key] = ab25ao25Val
		}
		if _, ab25ao25EndErr := ab25ao25Dec.Token(); ab25ao25EndErr != nil {
			return nil, ab25ao25EndErr
		}
		return ab25ao25Map, nil
	}
	ab25ao25List := []any{}
	for ab25ao25Dec.More() {
		ab25ao25Val, ab25ao25ValErr := ab25ao25Parse(ab25ao25Dec, ab25ao25Depth+1)
		if ab25ao25ValErr != nil {
			return nil, ab25ao25ValErr
		}
		ab25ao25List = append(ab25ao25List, ab25ao25Val)
	}
	if _, ab25ao25EndErr := ab25ao25Dec.Token(); ab25ao25EndErr != nil {
		return nil, ab25ao25EndErr
	}
	return ab25ao25List, nil
}

func ab25ao25Copy(ab25ao25Value any) any {
	switch ab25ao25Typed := ab25ao25Value.(type) {
	case map[string]any:
		ab25ao25Out := make(map[string]any, len(ab25ao25Typed))
		for ab25ao25Key, ab25ao25Item := range ab25ao25Typed {
			ab25ao25Out[ab25ao25Key] = ab25ao25Copy(ab25ao25Item)
		}
		return ab25ao25Out
	case []any:
		ab25ao25Out := make([]any, len(ab25ao25Typed))
		for ab25ao25Index, ab25ao25Item := range ab25ao25Typed {
			ab25ao25Out[ab25ao25Index] = ab25ao25Copy(ab25ao25Item)
		}
		return ab25ao25Out
	}
	return ab25ao25Value
}

func ab25ao25Integer(ab25ao25Value any) (string, bool) {
	ab25ao25Num, ab25ao25OK := ab25ao25Value.(json.Number)
	if !ab25ao25OK || !ab25ao25IntegerPattern.MatchString(ab25ao25Num.String()) {
		return "", false
	}
	return ab25ao25Num.String(), true
}

func ab25ao25Bytes(ab25ao25Name string, ab25ao25Value any) ([]int, error) {
	ab25ao25List, ab25ao25OK := ab25ao25Value.([]any)
	if !ab25ao25OK {
		return nil, fmt.Errorf("ab25ao25: field %q is not an array", ab25ao25Name)
	}
	ab25ao25Out := make([]int, 0, len(ab25ao25List))
	for _, ab25ao25Item := range ab25ao25List {
		ab25ao25Text, ab25ao25IntOK := ab25ao25Integer(ab25ao25Item)
		if !ab25ao25IntOK {
			return nil, fmt.Errorf("ab25ao25: field %q has a non-integer element", ab25ao25Name)
		}
		ab25ao25N, ab25ao25Err := strconv.Atoi(ab25ao25Text)
		if ab25ao25Err != nil || ab25ao25N < 0 || ab25ao25N > 255 {
			return nil, fmt.Errorf("ab25ao25: field %q has an element outside 0..255", ab25ao25Name)
		}
		ab25ao25Out = append(ab25ao25Out, ab25ao25N)
	}
	return ab25ao25Out, nil
}

func ab25ao25Field(ab25ao25Raw []int, ab25ao25Present bool, ab25ao25Decode func(int) map[string]any) map[string]any {
	ab25ao25Out := map[string]any{"grade": "Hypothesis"}
	if !ab25ao25Present {
		ab25ao25Out["raw"] = nil
		ab25ao25Out["state"] = "absent"
		return ab25ao25Out
	}
	ab25ao25Out["raw"] = append([]int{}, ab25ao25Raw...)
	switch {
	case len(ab25ao25Raw) == 0:
		ab25ao25Out["state"] = "incomplete"
	case len(ab25ao25Raw) > 1:
		ab25ao25Out["state"] = "overlong"
	default:
		for ab25ao25Key, ab25ao25Val := range ab25ao25Decode(ab25ao25Raw[0]) {
			ab25ao25Out[ab25ao25Key] = ab25ao25Val
		}
	}
	return ab25ao25Out
}

func ab25ao25DecodeMode(ab25ao25Byte int) map[string]any {
	ab25ao25Label, ab25ao25Known := ab25ao25ModeLabels[ab25ao25Byte]
	if !ab25ao25Known {
		return map[string]any{"state": "unmapped"}
	}
	if ab25ao25Byte == 1 {
		return map[string]any{"state": "label-conflict", "number": ab25ao25Byte, "label": ab25ao25Label}
	}
	return map[string]any{"state": "decoded", "number": ab25ao25Byte, "label": ab25ao25Label}
}

func ab25ao25DecodeFan(ab25ao25Byte int) map[string]any {
	return map[string]any{"state": "decoded", "number": ab25ao25Byte}
}

func ab25ao25DecodeRoom(ab25ao25Byte int) map[string]any {
	return map[string]any{
		"state":            "decoded",
		"number":           float64(ab25ao25Byte) / 2,
		"unit":             "celsius",
		"calculation":      "Proven",
		"candidate_domain": true,
	}
}

func ab25ao25DecodeAutomatic(ab25ao25Byte int) map[string]any {
	ab25ao25Label := "disabled"
	if ab25ao25Byte&1 == 1 {
		ab25ao25Label = "enabled"
	}
	return map[string]any{
		"state":    "decoded",
		"number":   ab25ao25Byte & 1,
		"label":    ab25ao25Label,
		"reserved": ab25ao25Byte & 0xFE,
	}
}

func ab25ao25DecodePositive(ab25ao25Byte int) map[string]any {
	return map[string]any{
		"state":    "decoded",
		"flags":    []bool{ab25ao25Byte&1 != 0, ab25ao25Byte&2 != 0, ab25ao25Byte&4 != 0},
		"reserved": ab25ao25Byte & 0xF8,
	}
}

// Ab25Ao25Interpret interprets B522 demand and raw scalar views selected
// explicitly by the input context. The result is a software interpretation;
// physical support stays Unknown and nothing is qualified.
func Ab25Ao25Interpret(ab25ao25Input json.RawMessage) (json.RawMessage, error) {
	ab25ao25Dec := json.NewDecoder(bytes.NewReader(ab25ao25Input))
	ab25ao25Dec.UseNumber()
	ab25ao25Root, ab25ao25Err := ab25ao25Parse(ab25ao25Dec, 0)
	if ab25ao25Err != nil {
		return nil, fmt.Errorf("ab25ao25: invalid input: %w", ab25ao25Err)
	}
	if _, ab25ao25TailErr := ab25ao25Dec.Token(); !errors.Is(ab25ao25TailErr, io.EOF) {
		return nil, errors.New("ab25ao25: trailing data after input")
	}
	ab25ao25Obj, ab25ao25ObjOK := ab25ao25Root.(map[string]any)
	if !ab25ao25ObjOK {
		return nil, errors.New("ab25ao25: input is not an object")
	}
	ab25ao25Ctx, ab25ao25CtxOK := ab25ao25Obj["context"].(map[string]any)
	if !ab25ao25CtxOK {
		return nil, errors.New("ab25ao25: context missing or not an object")
	}
	ab25ao25Primary, ab25ao25PrimaryOK := ab25ao25Integer(ab25ao25Ctx["primary"])
	ab25ao25Secondary, ab25ao25SecondaryOK := ab25ao25Integer(ab25ao25Ctx["secondary"])
	ab25ao25Operation, ab25ao25OperationOK := ab25ao25Integer(ab25ao25Ctx["operation"])
	ab25ao25Direction, ab25ao25DirectionOK := ab25ao25Ctx["direction"].(string)
	ab25ao25View, ab25ao25ViewOK := ab25ao25Ctx["view"].(string)
	if !ab25ao25PrimaryOK || !ab25ao25SecondaryOK || !ab25ao25OperationOK || !ab25ao25DirectionOK || !ab25ao25ViewOK {
		return nil, errors.New("ab25ao25: context requires integer primary, secondary, operation and string direction, view")
	}
	ab25ao25FieldsIn, ab25ao25FieldsOK := ab25ao25Obj["fields"].(map[string]any)
	if !ab25ao25FieldsOK {
		return nil, errors.New("ab25ao25: fields missing or not an object")
	}
	ab25ao25Raw := map[string]any{}
	ab25ao25ByteMap := map[string][]int{}
	for ab25ao25Name, ab25ao25Value := range ab25ao25FieldsIn {
		ab25ao25Parsed, ab25ao25ByteErr := ab25ao25Bytes(ab25ao25Name, ab25ao25Value)
		if ab25ao25ByteErr != nil {
			return nil, ab25ao25ByteErr
		}
		ab25ao25ByteMap[ab25ao25Name] = ab25ao25Parsed
		ab25ao25Raw[ab25ao25Name] = ab25ao25Copy(ab25ao25Value)
	}
	ab25ao25Out := map[string]any{
		"context":   ab25ao25Copy(ab25ao25Ctx),
		"raw":       ab25ao25Raw,
		"physical":  "Unknown",
		"qualified": false,
	}
	if ab25ao25ArgsIn, ab25ao25ArgsPresent := ab25ao25Obj["arguments"]; ab25ao25ArgsPresent {
		ab25ao25Args, ab25ao25ArgsOK := ab25ao25ArgsIn.(map[string]any)
		if !ab25ao25ArgsOK {
			return nil, errors.New("ab25ao25: arguments is not an object")
		}
		for ab25ao25Name, ab25ao25Value := range ab25ao25Args {
			if _, ab25ao25IsInt := ab25ao25Integer(ab25ao25Value); !ab25ao25IsInt {
				if _, ab25ao25IsBool := ab25ao25Value.(bool); !ab25ao25IsBool {
					return nil, fmt.Errorf("ab25ao25: argument %q is not a boolean or integer", ab25ao25Name)
				}
			}
		}
		ab25ao25Out["arguments"] = ab25ao25Copy(ab25ao25Args)
	}
	ab25ao25Extra := map[string]any{}
	for ab25ao25Name, ab25ao25Value := range ab25ao25Obj {
		if ab25ao25Name != "context" && ab25ao25Name != "fields" && ab25ao25Name != "arguments" {
			ab25ao25Extra[ab25ao25Name] = ab25ao25Copy(ab25ao25Value)
		}
	}
	if len(ab25ao25Extra) > 0 {
		ab25ao25Out["unrecognised"] = ab25ao25Extra
	}

	ab25ao25Reasons := []string{}
	if ab25ao25Primary != "181" {
		ab25ao25Reasons = append(ab25ao25Reasons, "primary")
	}
	if ab25ao25Secondary != "34" {
		ab25ao25Reasons = append(ab25ao25Reasons, "secondary")
	}
	if ab25ao25OpN, ab25ao25OpErr := strconv.Atoi(ab25ao25Operation); ab25ao25OpErr != nil || ab25ao25OpN < 0 || ab25ao25OpN > 255 {
		ab25ao25Reasons = append(ab25ao25Reasons, "operation")
	}
	if ab25ao25Direction != "request" && ab25ao25Direction != "positive" {
		ab25ao25Reasons = append(ab25ao25Reasons, "direction")
	}
	if ab25ao25View != "isolated-fields" {
		ab25ao25Reasons = append(ab25ao25Reasons, "view")
	}
	if len(ab25ao25Reasons) > 0 {
		ab25ao25Out["state"] = "unknown"
		ab25ao25Out["unknown_reasons"] = ab25ao25Reasons
		ab25ao25Out["fields"] = map[string]any{}
		return json.Marshal(ab25ao25Out)
	}

	type ab25ao25Spec struct {
		name   string
		decode func(int) map[string]any
	}
	ab25ao25Specs := []ab25ao25Spec{{"positive_flags", ab25ao25DecodePositive}}
	if ab25ao25Direction == "request" {
		ab25ao25Specs = []ab25ao25Spec{
			{"mode", ab25ao25DecodeMode},
			{"fan", ab25ao25DecodeFan},
			{"room_temperature", ab25ao25DecodeRoom},
			{"automatic", ab25ao25DecodeAutomatic},
		}
	}
	ab25ao25Results := map[string]any{}
	ab25ao25Decoded := 0
	for _, ab25ao25Item := range ab25ao25Specs {
		ab25ao25Data, ab25ao25Present := ab25ao25ByteMap[ab25ao25Item.name]
		ab25ao25Result := ab25ao25Field(ab25ao25Data, ab25ao25Present, ab25ao25Item.decode)
		if ab25ao25Result["state"] == "decoded" {
			ab25ao25Decoded++
		}
		ab25ao25Results[ab25ao25Item.name] = ab25ao25Result
	}
	switch {
	case ab25ao25Decoded == len(ab25ao25Specs):
		ab25ao25Out["state"] = "complete"
	case ab25ao25Decoded > 0:
		ab25ao25Out["state"] = "partial"
	default:
		ab25ao25Out["state"] = "unresolved"
	}
	ab25ao25Out["fields"] = ab25ao25Results
	return json.Marshal(ab25ao25Out)
}
