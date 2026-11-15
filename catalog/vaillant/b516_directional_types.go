package vaillant

import (
	b516DirectionalTypeBytes "bytes"
	_ "embed"
	b516DirectionalTypeJSON "encoding/json"
	b516DirectionalTypeFmt "fmt"
	b516DirectionalTypeReflect "reflect"
)

//go:embed data/b516_directional_types.json
var b516DirectionalTypeData string

// B516DirectionalTypeKey identifies one B516 directional type row; nil optional axes are absent and
// lookups compare every axis exactly.
type B516DirectionalTypeKey struct {
	B516DirectionalTypeFamily        string  `json:"family"`
	B516DirectionalTypeVariant       string  `json:"variant"`
	B516DirectionalTypePrimary       uint8   `json:"primary"`
	B516DirectionalTypeSecondary     *uint8  `json:"secondary,omitempty"`
	B516DirectionalTypeOperation     *uint8  `json:"operation,omitempty"`
	B516DirectionalTypeDirection     string  `json:"direction"`
	B516DirectionalTypeField         string  `json:"field"`
	B516DirectionalTypeAlternative   string  `json:"alternative"`
	B516DirectionalTypeGroup         *uint8  `json:"group,omitempty"`
	B516DirectionalTypeInstanceScope string  `json:"instance_scope"`
	B516DirectionalTypeRegister      *uint16 `json:"register,omitempty"`
}

func b516DirectionalTypeValidate(b516DirectionalTypeData []byte) ([]b516DirectionalTypeJSON.RawMessage, error) {
	var b516DirectionalTypeRows []b516DirectionalTypeJSON.RawMessage
	if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeData, &b516DirectionalTypeRows); b516DirectionalTypeErr != nil {
		return nil, b516DirectionalTypeFmt.Errorf("b516 directional types catalog: %w", b516DirectionalTypeErr)
	}
	if len(b516DirectionalTypeRows) == 0 {
		return nil, b516DirectionalTypeFmt.Errorf("b516 directional types catalog: empty array")
	}
	b516DirectionalTypeKeys := make([]B516DirectionalTypeKey, 0, len(b516DirectionalTypeRows))
	b516DirectionalTypeIDs := make(map[string]bool)
	for _, b516DirectionalTypeRow := range b516DirectionalTypeRows {
		var b516DirectionalTypeEntry struct {
			B516DirectionalTypeID  string                             `json:"id"`
			B516DirectionalTypeKey b516DirectionalTypeJSON.RawMessage `json:"key"`
		}
		if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeRow, &b516DirectionalTypeEntry); b516DirectionalTypeErr != nil {
			return nil, b516DirectionalTypeFmt.Errorf("b516 directional types catalog: %w", b516DirectionalTypeErr)
		}
		var b516DirectionalTypeObject map[string]b516DirectionalTypeJSON.RawMessage
		if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeRow, &b516DirectionalTypeObject); b516DirectionalTypeErr != nil {
			return nil, b516DirectionalTypeErr
		}
		var b516DirectionalTypeIdentity map[string]b516DirectionalTypeJSON.RawMessage
		if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeObject["key"], &b516DirectionalTypeIdentity); b516DirectionalTypeErr != nil {
			return nil, b516DirectionalTypeErr
		}
		if b516DirectionalTypePrimary, b516DirectionalTypePresent := b516DirectionalTypeIdentity["primary"]; !b516DirectionalTypePresent || string(b516DirectionalTypePrimary) == "null" {
			return nil, b516DirectionalTypeFmt.Errorf("b516 directional types catalog: missing primary")
		}
		var b516DirectionalTypeKey B516DirectionalTypeKey
		b516DirectionalTypeDecoder := b516DirectionalTypeJSON.NewDecoder(b516DirectionalTypeBytes.NewReader(b516DirectionalTypeEntry.B516DirectionalTypeKey))
		b516DirectionalTypeDecoder.DisallowUnknownFields()
		if b516DirectionalTypeErr := b516DirectionalTypeDecoder.Decode(&b516DirectionalTypeKey); b516DirectionalTypeErr != nil {
			return nil, b516DirectionalTypeFmt.Errorf("b516 directional types catalog: %w", b516DirectionalTypeErr)
		}
		if b516DirectionalTypeEntry.B516DirectionalTypeID == "" || b516DirectionalTypeIDs[b516DirectionalTypeEntry.B516DirectionalTypeID] || b516DirectionalTypeKey.B516DirectionalTypeFamily == "" || b516DirectionalTypeKey.B516DirectionalTypeVariant == "" || b516DirectionalTypeKey.B516DirectionalTypeDirection == "" || b516DirectionalTypeKey.B516DirectionalTypeField == "" || b516DirectionalTypeKey.B516DirectionalTypeAlternative == "" || b516DirectionalTypeKey.B516DirectionalTypeInstanceScope == "" {
			return nil, b516DirectionalTypeFmt.Errorf("b516 directional types catalog: incomplete identity or duplicate id")
		}
		for _, b516DirectionalTypePrevious := range b516DirectionalTypeKeys {
			if b516DirectionalTypeReflect.DeepEqual(b516DirectionalTypePrevious, b516DirectionalTypeKey) {
				return nil, b516DirectionalTypeFmt.Errorf("b516 directional types catalog: duplicate key")
			}
		}
		b516DirectionalTypeIDs[b516DirectionalTypeEntry.B516DirectionalTypeID] = true
		b516DirectionalTypeKeys = append(b516DirectionalTypeKeys, b516DirectionalTypeKey)
	}
	return b516DirectionalTypeRows, nil
}

// B516DirectionalTypeEntries validates the embedded B516 directional type catalog (unique ids, complete
// keys, no duplicate key) and returns a caller-owned copy of its JSON, or the validation error.
func B516DirectionalTypeEntries() (b516DirectionalTypeJSON.RawMessage, error) {
	if _, b516DirectionalTypeErr := b516DirectionalTypeValidate([]byte(b516DirectionalTypeData)); b516DirectionalTypeErr != nil {
		return nil, b516DirectionalTypeErr
	}
	return b516DirectionalTypeJSON.RawMessage([]byte(b516DirectionalTypeData)), nil
}

// B516DirectionalTypeLookup returns a caller-owned copy of the row whose key equals the given key, false
// when none matches, or the catalog validation error.
func B516DirectionalTypeLookup(b516DirectionalTypeKey B516DirectionalTypeKey) (b516DirectionalTypeJSON.RawMessage, bool, error) {
	return b516DirectionalTypeLookup([]byte(b516DirectionalTypeData), b516DirectionalTypeKey)
}

func b516DirectionalTypeLookup(b516DirectionalTypeData []byte, b516DirectionalTypeKey B516DirectionalTypeKey) (b516DirectionalTypeJSON.RawMessage, bool, error) {
	b516DirectionalTypeRows, b516DirectionalTypeErr := b516DirectionalTypeValidate(b516DirectionalTypeData)
	if b516DirectionalTypeErr != nil {
		return nil, false, b516DirectionalTypeErr
	}
	for _, b516DirectionalTypeRow := range b516DirectionalTypeRows {
		var b516DirectionalTypeEntry struct {
			B516DirectionalTypeKey B516DirectionalTypeKey `json:"key"`
		}
		if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeRow, &b516DirectionalTypeEntry); b516DirectionalTypeErr != nil {
			return nil, false, b516DirectionalTypeErr
		}
		if b516DirectionalTypeReflect.DeepEqual(b516DirectionalTypeKey, b516DirectionalTypeEntry.B516DirectionalTypeKey) {
			return append(b516DirectionalTypeJSON.RawMessage(nil), b516DirectionalTypeRow...), true, nil
		}
	}
	return nil, false, nil
}
