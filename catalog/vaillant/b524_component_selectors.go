package vaillant

import (
	b524ComponentSelectorEmbed "embed"
	b524ComponentSelectorJSON "encoding/json"
	b524ComponentSelectorFmt "fmt"
	b524ComponentSelectorReflect "reflect"
	b524ComponentSelectorStrings "strings"
)

//go:embed data/ab32_b524_component_selectors.json
var b524ComponentSelectorFiles b524ComponentSelectorEmbed.FS

type B524ComponentSelectorKey struct {
	B524ComponentSelectorFamily        string  `json:"family"`
	B524ComponentSelectorVariant       string  `json:"variant"`
	B524ComponentSelectorPrimary       uint8   `json:"primary"`
	B524ComponentSelectorSecondary     *uint8  `json:"secondary"`
	B524ComponentSelectorOperation     *uint8  `json:"operation"`
	B524ComponentSelectorDirection     string  `json:"direction"`
	B524ComponentSelectorField         string  `json:"field"`
	B524ComponentSelectorAlternative   string  `json:"alternative"`
	B524ComponentSelectorGroup         *uint8  `json:"group"`
	B524ComponentSelectorInstanceScope string  `json:"instance_scope"`
	B524ComponentSelectorRegister      *uint16 `json:"register"`
}

func b524ComponentSelectorValidate(b524ComponentSelectorData []byte) ([]b524ComponentSelectorJSON.RawMessage, []B524ComponentSelectorKey, error) {
	var b524ComponentSelectorRows []b524ComponentSelectorJSON.RawMessage
	if b524ComponentSelectorError := b524ComponentSelectorJSON.Unmarshal(b524ComponentSelectorData, &b524ComponentSelectorRows); b524ComponentSelectorError != nil {
		return nil, nil, b524ComponentSelectorError
	}
	if b524ComponentSelectorRows == nil {
		return nil, nil, b524ComponentSelectorFmt.Errorf("ab32: expected array")
	}
	b524ComponentSelectorKeys := make([]B524ComponentSelectorKey, 0, len(b524ComponentSelectorRows))
	b524ComponentSelectorIDs := make(map[string]bool)
	for _, b524ComponentSelectorRow := range b524ComponentSelectorRows {
		var b524ComponentSelectorEnvelope struct {
			B524ComponentSelectorID  string                               `json:"id"`
			B524ComponentSelectorKey b524ComponentSelectorJSON.RawMessage `json:"key"`
		}
		if b524ComponentSelectorError := b524ComponentSelectorJSON.Unmarshal(b524ComponentSelectorRow, &b524ComponentSelectorEnvelope); b524ComponentSelectorError != nil {
			return nil, nil, b524ComponentSelectorError
		}
		if !b524ComponentSelectorStrings.HasPrefix(b524ComponentSelectorEnvelope.B524ComponentSelectorID, "ab32-") || b524ComponentSelectorIDs[b524ComponentSelectorEnvelope.B524ComponentSelectorID] {
			return nil, nil, b524ComponentSelectorFmt.Errorf("ab32: invalid or duplicate id")
		}
		b524ComponentSelectorIDs[b524ComponentSelectorEnvelope.B524ComponentSelectorID] = true
		var b524ComponentSelectorFields map[string]b524ComponentSelectorJSON.RawMessage
		if b524ComponentSelectorError := b524ComponentSelectorJSON.Unmarshal(b524ComponentSelectorEnvelope.B524ComponentSelectorKey, &b524ComponentSelectorFields); b524ComponentSelectorError != nil {
			return nil, nil, b524ComponentSelectorError
		}
		b524ComponentSelectorAxes := []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"}
		if len(b524ComponentSelectorFields) != len(b524ComponentSelectorAxes) {
			return nil, nil, b524ComponentSelectorFmt.Errorf("ab32: incomplete identity")
		}
		for _, b524ComponentSelectorAxis := range b524ComponentSelectorAxes {
			if _, b524ComponentSelectorPresent := b524ComponentSelectorFields[b524ComponentSelectorAxis]; !b524ComponentSelectorPresent {
				return nil, nil, b524ComponentSelectorFmt.Errorf("ab32: missing %s", b524ComponentSelectorAxis)
			}
		}
		if string(b524ComponentSelectorFields["primary"]) == "null" {
			return nil, nil, b524ComponentSelectorFmt.Errorf("ab32: null primary")
		}
		var b524ComponentSelectorKey B524ComponentSelectorKey
		if b524ComponentSelectorError := b524ComponentSelectorJSON.Unmarshal(b524ComponentSelectorEnvelope.B524ComponentSelectorKey, &b524ComponentSelectorKey); b524ComponentSelectorError != nil {
			return nil, nil, b524ComponentSelectorError
		}
		for _, b524ComponentSelectorText := range []string{b524ComponentSelectorKey.B524ComponentSelectorFamily, b524ComponentSelectorKey.B524ComponentSelectorVariant, b524ComponentSelectorKey.B524ComponentSelectorDirection, b524ComponentSelectorKey.B524ComponentSelectorField, b524ComponentSelectorKey.B524ComponentSelectorAlternative, b524ComponentSelectorKey.B524ComponentSelectorInstanceScope} {
			if b524ComponentSelectorText == "" {
				return nil, nil, b524ComponentSelectorFmt.Errorf("ab32: empty identity string")
			}
		}
		for _, b524ComponentSelectorPrevious := range b524ComponentSelectorKeys {
			if b524ComponentSelectorReflect.DeepEqual(b524ComponentSelectorPrevious, b524ComponentSelectorKey) {
				return nil, nil, b524ComponentSelectorFmt.Errorf("ab32: duplicate key")
			}
		}
		b524ComponentSelectorKeys = append(b524ComponentSelectorKeys, b524ComponentSelectorKey)
	}
	return b524ComponentSelectorRows, b524ComponentSelectorKeys, nil
}

func b524ComponentSelectorRead() ([]byte, []b524ComponentSelectorJSON.RawMessage, []B524ComponentSelectorKey, error) {
	b524ComponentSelectorData, b524ComponentSelectorError := b524ComponentSelectorFiles.ReadFile("data/ab32_b524_component_selectors.json")
	if b524ComponentSelectorError != nil {
		return nil, nil, nil, b524ComponentSelectorError
	}
	b524ComponentSelectorRows, b524ComponentSelectorKeys, b524ComponentSelectorError := b524ComponentSelectorValidate(b524ComponentSelectorData)
	return b524ComponentSelectorData, b524ComponentSelectorRows, b524ComponentSelectorKeys, b524ComponentSelectorError
}

func B524ComponentSelectorEntries() (b524ComponentSelectorJSON.RawMessage, error) {
	b524ComponentSelectorData, _, _, b524ComponentSelectorError := b524ComponentSelectorRead()
	if b524ComponentSelectorError != nil {
		return nil, b524ComponentSelectorError
	}
	return append(b524ComponentSelectorJSON.RawMessage(nil), b524ComponentSelectorData...), nil
}

func B524ComponentSelectorLookup(b524ComponentSelectorKey B524ComponentSelectorKey) (b524ComponentSelectorJSON.RawMessage, bool, error) {
	b524ComponentSelectorData, _, _, b524ComponentSelectorError := b524ComponentSelectorRead()
	if b524ComponentSelectorError != nil {
		return nil, false, b524ComponentSelectorError
	}
	return b524ComponentSelectorLookupData(b524ComponentSelectorData, b524ComponentSelectorKey)
}

func b524ComponentSelectorLookupData(b524ComponentSelectorData []byte, b524ComponentSelectorKey B524ComponentSelectorKey) (b524ComponentSelectorJSON.RawMessage, bool, error) {
	b524ComponentSelectorRows, b524ComponentSelectorKeys, b524ComponentSelectorError := b524ComponentSelectorValidate(b524ComponentSelectorData)
	if b524ComponentSelectorError != nil {
		return nil, false, b524ComponentSelectorError
	}
	for b524ComponentSelectorIndex, b524ComponentSelectorCandidate := range b524ComponentSelectorKeys {
		if b524ComponentSelectorReflect.DeepEqual(b524ComponentSelectorCandidate, b524ComponentSelectorKey) {
			return append(b524ComponentSelectorJSON.RawMessage(nil), b524ComponentSelectorRows[b524ComponentSelectorIndex]...), true, nil
		}
	}
	return nil, false, nil
}
