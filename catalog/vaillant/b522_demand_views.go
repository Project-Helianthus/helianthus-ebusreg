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

var b522DemandViewIntegerPattern = regexp.MustCompile(`^-?(0|[1-9][0-9]*)$`)

const b522DemandViewMaxDepth = 64

var b522DemandViewModeLabels = map[int]string{
	0: "no-heating-demand",
	1: "no-source-demand-sink-demand",
	2: "source-heating-demand",
	3: "no-source-cooling-demand-sink-demand",
	4: "source-cooling-demand",
}

func b522DemandViewParse(b522DemandViewDec *json.Decoder, b522DemandViewDepth int) (any, error) {
	if b522DemandViewDepth > b522DemandViewMaxDepth {
		return nil, errors.New("ab25ao25: nesting too deep")
	}
	b522DemandViewTok, b522DemandViewErr := b522DemandViewDec.Token()
	if b522DemandViewErr != nil {
		return nil, b522DemandViewErr
	}
	b522DemandViewDelim, b522DemandViewIsDelim := b522DemandViewTok.(json.Delim)
	if !b522DemandViewIsDelim {
		return b522DemandViewTok, nil
	}
	if b522DemandViewDelim == '{' {
		b522DemandViewMap := map[string]any{}
		for b522DemandViewDec.More() {
			b522DemandViewKeyTok, b522DemandViewKeyErr := b522DemandViewDec.Token()
			if b522DemandViewKeyErr != nil {
				return nil, b522DemandViewKeyErr
			}
			b522DemandViewKey, b522DemandViewKeyOK := b522DemandViewKeyTok.(string)
			if !b522DemandViewKeyOK {
				return nil, errors.New("ab25ao25: object key is not a string")
			}
			if _, b522DemandViewDup := b522DemandViewMap[b522DemandViewKey]; b522DemandViewDup {
				return nil, fmt.Errorf("ab25ao25: duplicate key %q", b522DemandViewKey)
			}
			b522DemandViewVal, b522DemandViewValErr := b522DemandViewParse(b522DemandViewDec, b522DemandViewDepth+1)
			if b522DemandViewValErr != nil {
				return nil, b522DemandViewValErr
			}
			b522DemandViewMap[b522DemandViewKey] = b522DemandViewVal
		}
		if _, b522DemandViewEndErr := b522DemandViewDec.Token(); b522DemandViewEndErr != nil {
			return nil, b522DemandViewEndErr
		}
		return b522DemandViewMap, nil
	}
	b522DemandViewList := []any{}
	for b522DemandViewDec.More() {
		b522DemandViewVal, b522DemandViewValErr := b522DemandViewParse(b522DemandViewDec, b522DemandViewDepth+1)
		if b522DemandViewValErr != nil {
			return nil, b522DemandViewValErr
		}
		b522DemandViewList = append(b522DemandViewList, b522DemandViewVal)
	}
	if _, b522DemandViewEndErr := b522DemandViewDec.Token(); b522DemandViewEndErr != nil {
		return nil, b522DemandViewEndErr
	}
	return b522DemandViewList, nil
}

func b522DemandViewCopy(b522DemandViewValue any) any {
	switch b522DemandViewTyped := b522DemandViewValue.(type) {
	case map[string]any:
		b522DemandViewOut := make(map[string]any, len(b522DemandViewTyped))
		for b522DemandViewKey, b522DemandViewItem := range b522DemandViewTyped {
			b522DemandViewOut[b522DemandViewKey] = b522DemandViewCopy(b522DemandViewItem)
		}
		return b522DemandViewOut
	case []any:
		b522DemandViewOut := make([]any, len(b522DemandViewTyped))
		for b522DemandViewIndex, b522DemandViewItem := range b522DemandViewTyped {
			b522DemandViewOut[b522DemandViewIndex] = b522DemandViewCopy(b522DemandViewItem)
		}
		return b522DemandViewOut
	}
	return b522DemandViewValue
}

func b522DemandViewInteger(b522DemandViewValue any) (string, bool) {
	b522DemandViewNum, b522DemandViewOK := b522DemandViewValue.(json.Number)
	if !b522DemandViewOK || !b522DemandViewIntegerPattern.MatchString(b522DemandViewNum.String()) {
		return "", false
	}
	return b522DemandViewNum.String(), true
}

func b522DemandViewBytes(b522DemandViewName string, b522DemandViewValue any) ([]int, error) {
	b522DemandViewList, b522DemandViewOK := b522DemandViewValue.([]any)
	if !b522DemandViewOK {
		return nil, fmt.Errorf("ab25ao25: field %q is not an array", b522DemandViewName)
	}
	b522DemandViewOut := make([]int, 0, len(b522DemandViewList))
	for _, b522DemandViewItem := range b522DemandViewList {
		b522DemandViewText, b522DemandViewIntOK := b522DemandViewInteger(b522DemandViewItem)
		if !b522DemandViewIntOK {
			return nil, fmt.Errorf("ab25ao25: field %q has a non-integer element", b522DemandViewName)
		}
		b522DemandViewN, b522DemandViewErr := strconv.Atoi(b522DemandViewText)
		if b522DemandViewErr != nil || b522DemandViewN < 0 || b522DemandViewN > 255 {
			return nil, fmt.Errorf("ab25ao25: field %q has an element outside 0..255", b522DemandViewName)
		}
		b522DemandViewOut = append(b522DemandViewOut, b522DemandViewN)
	}
	return b522DemandViewOut, nil
}

func b522DemandViewField(b522DemandViewRaw []int, b522DemandViewPresent bool, b522DemandViewDecode func(int) map[string]any) map[string]any {
	b522DemandViewOut := map[string]any{"grade": "Hypothesis"}
	if !b522DemandViewPresent {
		b522DemandViewOut["raw"] = nil
		b522DemandViewOut["state"] = "absent"
		return b522DemandViewOut
	}
	b522DemandViewOut["raw"] = append([]int{}, b522DemandViewRaw...)
	switch {
	case len(b522DemandViewRaw) == 0:
		b522DemandViewOut["state"] = "incomplete"
	case len(b522DemandViewRaw) > 1:
		b522DemandViewOut["state"] = "overlong"
	default:
		for b522DemandViewKey, b522DemandViewVal := range b522DemandViewDecode(b522DemandViewRaw[0]) {
			b522DemandViewOut[b522DemandViewKey] = b522DemandViewVal
		}
	}
	return b522DemandViewOut
}

func b522DemandViewDecodeMode(b522DemandViewByte int) map[string]any {
	b522DemandViewLabel, b522DemandViewKnown := b522DemandViewModeLabels[b522DemandViewByte]
	if !b522DemandViewKnown {
		return map[string]any{"state": "unmapped"}
	}
	if b522DemandViewByte == 1 {
		return map[string]any{"state": "label-conflict", "number": b522DemandViewByte, "label": b522DemandViewLabel}
	}
	return map[string]any{"state": "decoded", "number": b522DemandViewByte, "label": b522DemandViewLabel}
}

func b522DemandViewDecodeFan(b522DemandViewByte int) map[string]any {
	return map[string]any{"state": "decoded", "number": b522DemandViewByte}
}

func b522DemandViewDecodeRoom(b522DemandViewByte int) map[string]any {
	return map[string]any{
		"state":            "decoded",
		"number":           float64(b522DemandViewByte) / 2,
		"unit":             "celsius",
		"calculation":      "Proven",
		"candidate_domain": true,
	}
}

func b522DemandViewDecodeAutomatic(b522DemandViewByte int) map[string]any {
	b522DemandViewLabel := "disabled"
	if b522DemandViewByte&1 == 1 {
		b522DemandViewLabel = "enabled"
	}
	return map[string]any{
		"state":    "decoded",
		"number":   b522DemandViewByte & 1,
		"label":    b522DemandViewLabel,
		"reserved": b522DemandViewByte & 0xFE,
	}
}

func b522DemandViewDecodePositive(b522DemandViewByte int) map[string]any {
	return map[string]any{
		"state":    "decoded",
		"flags":    []bool{b522DemandViewByte&1 != 0, b522DemandViewByte&2 != 0, b522DemandViewByte&4 != 0},
		"reserved": b522DemandViewByte & 0xF8,
	}
}

// B522DemandViewInterpret interprets B522 demand and raw scalar views selected
// explicitly by the input context. The result is a software interpretation;
// physical support stays Unknown and nothing is qualified.
func B522DemandViewInterpret(b522DemandViewInput json.RawMessage) (json.RawMessage, error) {
	b522DemandViewDec := json.NewDecoder(bytes.NewReader(b522DemandViewInput))
	b522DemandViewDec.UseNumber()
	b522DemandViewRoot, b522DemandViewErr := b522DemandViewParse(b522DemandViewDec, 0)
	if b522DemandViewErr != nil {
		return nil, fmt.Errorf("ab25ao25: invalid input: %w", b522DemandViewErr)
	}
	if _, b522DemandViewTailErr := b522DemandViewDec.Token(); !errors.Is(b522DemandViewTailErr, io.EOF) {
		return nil, errors.New("ab25ao25: trailing data after input")
	}
	b522DemandViewObj, b522DemandViewObjOK := b522DemandViewRoot.(map[string]any)
	if !b522DemandViewObjOK {
		return nil, errors.New("ab25ao25: input is not an object")
	}
	b522DemandViewCtx, b522DemandViewCtxOK := b522DemandViewObj["context"].(map[string]any)
	if !b522DemandViewCtxOK {
		return nil, errors.New("ab25ao25: context missing or not an object")
	}
	b522DemandViewPrimary, b522DemandViewPrimaryOK := b522DemandViewInteger(b522DemandViewCtx["primary"])
	b522DemandViewSecondary, b522DemandViewSecondaryOK := b522DemandViewInteger(b522DemandViewCtx["secondary"])
	b522DemandViewOperation, b522DemandViewOperationOK := b522DemandViewInteger(b522DemandViewCtx["operation"])
	b522DemandViewDirection, b522DemandViewDirectionOK := b522DemandViewCtx["direction"].(string)
	b522DemandViewView, b522DemandViewViewOK := b522DemandViewCtx["view"].(string)
	if !b522DemandViewPrimaryOK || !b522DemandViewSecondaryOK || !b522DemandViewOperationOK || !b522DemandViewDirectionOK || !b522DemandViewViewOK {
		return nil, errors.New("ab25ao25: context requires integer primary, secondary, operation and string direction, view")
	}
	b522DemandViewFieldsIn, b522DemandViewFieldsOK := b522DemandViewObj["fields"].(map[string]any)
	if !b522DemandViewFieldsOK {
		return nil, errors.New("ab25ao25: fields missing or not an object")
	}
	b522DemandViewRaw := map[string]any{}
	b522DemandViewByteMap := map[string][]int{}
	for b522DemandViewName, b522DemandViewValue := range b522DemandViewFieldsIn {
		b522DemandViewParsed, b522DemandViewByteErr := b522DemandViewBytes(b522DemandViewName, b522DemandViewValue)
		if b522DemandViewByteErr != nil {
			return nil, b522DemandViewByteErr
		}
		b522DemandViewByteMap[b522DemandViewName] = b522DemandViewParsed
		b522DemandViewRaw[b522DemandViewName] = b522DemandViewCopy(b522DemandViewValue)
	}
	b522DemandViewOut := map[string]any{
		"context":   b522DemandViewCopy(b522DemandViewCtx),
		"raw":       b522DemandViewRaw,
		"physical":  "Unknown",
		"qualified": false,
	}
	if b522DemandViewArgsIn, b522DemandViewArgsPresent := b522DemandViewObj["arguments"]; b522DemandViewArgsPresent {
		b522DemandViewArgs, b522DemandViewArgsOK := b522DemandViewArgsIn.(map[string]any)
		if !b522DemandViewArgsOK {
			return nil, errors.New("ab25ao25: arguments is not an object")
		}
		for b522DemandViewName, b522DemandViewValue := range b522DemandViewArgs {
			if _, b522DemandViewIsInt := b522DemandViewInteger(b522DemandViewValue); !b522DemandViewIsInt {
				if _, b522DemandViewIsBool := b522DemandViewValue.(bool); !b522DemandViewIsBool {
					return nil, fmt.Errorf("ab25ao25: argument %q is not a boolean or integer", b522DemandViewName)
				}
			}
		}
		b522DemandViewOut["arguments"] = b522DemandViewCopy(b522DemandViewArgs)
	}
	b522DemandViewExtra := map[string]any{}
	for b522DemandViewName, b522DemandViewValue := range b522DemandViewObj {
		if b522DemandViewName != "context" && b522DemandViewName != "fields" && b522DemandViewName != "arguments" {
			b522DemandViewExtra[b522DemandViewName] = b522DemandViewCopy(b522DemandViewValue)
		}
	}
	if len(b522DemandViewExtra) > 0 {
		b522DemandViewOut["unrecognised"] = b522DemandViewExtra
	}

	b522DemandViewReasons := []string{}
	if b522DemandViewPrimary != "181" {
		b522DemandViewReasons = append(b522DemandViewReasons, "primary")
	}
	if b522DemandViewSecondary != "34" {
		b522DemandViewReasons = append(b522DemandViewReasons, "secondary")
	}
	if b522DemandViewOpN, b522DemandViewOpErr := strconv.Atoi(b522DemandViewOperation); b522DemandViewOpErr != nil || b522DemandViewOpN < 0 || b522DemandViewOpN > 255 {
		b522DemandViewReasons = append(b522DemandViewReasons, "operation")
	}
	if b522DemandViewDirection != "request" && b522DemandViewDirection != "positive" {
		b522DemandViewReasons = append(b522DemandViewReasons, "direction")
	}
	if b522DemandViewView != "isolated-fields" {
		b522DemandViewReasons = append(b522DemandViewReasons, "view")
	}
	if len(b522DemandViewReasons) > 0 {
		b522DemandViewOut["state"] = "unknown"
		b522DemandViewOut["unknown_reasons"] = b522DemandViewReasons
		b522DemandViewOut["fields"] = map[string]any{}
		return json.Marshal(b522DemandViewOut)
	}

	type b522DemandViewSpec struct {
		name   string
		decode func(int) map[string]any
	}
	b522DemandViewSpecs := []b522DemandViewSpec{{"positive_flags", b522DemandViewDecodePositive}}
	if b522DemandViewDirection == "request" {
		b522DemandViewSpecs = []b522DemandViewSpec{
			{"mode", b522DemandViewDecodeMode},
			{"fan", b522DemandViewDecodeFan},
			{"room_temperature", b522DemandViewDecodeRoom},
			{"automatic", b522DemandViewDecodeAutomatic},
		}
	}
	b522DemandViewResults := map[string]any{}
	b522DemandViewDecoded := 0
	for _, b522DemandViewItem := range b522DemandViewSpecs {
		b522DemandViewData, b522DemandViewPresent := b522DemandViewByteMap[b522DemandViewItem.name]
		b522DemandViewResult := b522DemandViewField(b522DemandViewData, b522DemandViewPresent, b522DemandViewItem.decode)
		if b522DemandViewResult["state"] == "decoded" {
			b522DemandViewDecoded++
		}
		b522DemandViewResults[b522DemandViewItem.name] = b522DemandViewResult
	}
	switch {
	case b522DemandViewDecoded == len(b522DemandViewSpecs):
		b522DemandViewOut["state"] = "complete"
	case b522DemandViewDecoded > 0:
		b522DemandViewOut["state"] = "partial"
	default:
		b522DemandViewOut["state"] = "unresolved"
	}
	b522DemandViewOut["fields"] = b522DemandViewResults
	return json.Marshal(b522DemandViewOut)
}
