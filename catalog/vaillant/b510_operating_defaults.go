package vaillant

import (
	b510OperatingDefaultBytes "bytes"
	_ "embed"
	b510OperatingDefaultJSON "encoding/json"
	b510OperatingDefaultFmt "fmt"
	b510OperatingDefaultReflect "reflect"
)

//go:embed data/ab48_b510_operating_defaults.json
var b510OperatingDefaultData string

type B510OperatingDefaultKey struct {
	B510OperatingDefaultFamily        string  `json:"family"`
	B510OperatingDefaultVariant       string  `json:"variant"`
	B510OperatingDefaultPrimary       uint8   `json:"primary"`
	B510OperatingDefaultSecondary     *uint8  `json:"secondary"`
	B510OperatingDefaultOperation     *uint8  `json:"operation"`
	B510OperatingDefaultDirection     string  `json:"direction"`
	B510OperatingDefaultField         string  `json:"field"`
	B510OperatingDefaultAlternative   string  `json:"alternative"`
	B510OperatingDefaultGroup         *uint8  `json:"group"`
	B510OperatingDefaultInstanceScope string  `json:"instance_scope"`
	B510OperatingDefaultRegister      *uint16 `json:"register"`
}

func b510OperatingDefaultValidate(b510OperatingDefaultSource []byte) ([]b510OperatingDefaultJSON.RawMessage, []B510OperatingDefaultKey, error) {
	var b510OperatingDefaultRows []b510OperatingDefaultJSON.RawMessage
	if b510OperatingDefaultErr := b510OperatingDefaultJSON.Unmarshal(b510OperatingDefaultSource, &b510OperatingDefaultRows); b510OperatingDefaultErr != nil {
		return nil, nil, b510OperatingDefaultErr
	}
	if b510OperatingDefaultRows == nil {
		return nil, nil, b510OperatingDefaultFmt.Errorf("ab48 catalog: expected array")
	}
	b510OperatingDefaultKeys := make([]B510OperatingDefaultKey, len(b510OperatingDefaultRows))
	for b510OperatingDefaultIndex, b510OperatingDefaultRaw := range b510OperatingDefaultRows {
		var b510OperatingDefaultObject map[string]b510OperatingDefaultJSON.RawMessage
		if b510OperatingDefaultErr := b510OperatingDefaultJSON.Unmarshal(b510OperatingDefaultRaw, &b510OperatingDefaultObject); b510OperatingDefaultErr != nil {
			return nil, nil, b510OperatingDefaultErr
		}
		var b510OperatingDefaultMembers map[string]b510OperatingDefaultJSON.RawMessage
		if b510OperatingDefaultErr := b510OperatingDefaultJSON.Unmarshal(b510OperatingDefaultObject["key"], &b510OperatingDefaultMembers); b510OperatingDefaultErr != nil {
			return nil, nil, b510OperatingDefaultErr
		}
		for _, b510OperatingDefaultName := range []string{"family", "variant", "primary", "direction", "field", "alternative", "instance_scope"} {
			if b510OperatingDefaultValue, b510OperatingDefaultOK := b510OperatingDefaultMembers[b510OperatingDefaultName]; !b510OperatingDefaultOK || b510OperatingDefaultBytes.Equal(b510OperatingDefaultBytes.TrimSpace(b510OperatingDefaultValue), []byte("null")) {
				return nil, nil, b510OperatingDefaultFmt.Errorf("ab48 catalog: missing key axis %s", b510OperatingDefaultName)
			}
		}
		b510OperatingDefaultDecoder := b510OperatingDefaultJSON.NewDecoder(b510OperatingDefaultBytes.NewReader(b510OperatingDefaultObject["key"]))
		b510OperatingDefaultDecoder.DisallowUnknownFields()
		if b510OperatingDefaultErr := b510OperatingDefaultDecoder.Decode(&b510OperatingDefaultKeys[b510OperatingDefaultIndex]); b510OperatingDefaultErr != nil {
			return nil, nil, b510OperatingDefaultErr
		}
		b510OperatingDefaultKey := b510OperatingDefaultKeys[b510OperatingDefaultIndex]
		if b510OperatingDefaultKey.B510OperatingDefaultFamily == "" || b510OperatingDefaultKey.B510OperatingDefaultVariant == "" || b510OperatingDefaultKey.B510OperatingDefaultDirection == "" || b510OperatingDefaultKey.B510OperatingDefaultField == "" || b510OperatingDefaultKey.B510OperatingDefaultAlternative == "" || b510OperatingDefaultKey.B510OperatingDefaultInstanceScope == "" {
			return nil, nil, b510OperatingDefaultFmt.Errorf("ab48 catalog: empty key axis; use none for non-applicable strings")
		}
		for b510OperatingDefaultPrevious := 0; b510OperatingDefaultPrevious < b510OperatingDefaultIndex; b510OperatingDefaultPrevious++ {
			if b510OperatingDefaultReflect.DeepEqual(b510OperatingDefaultKeys[b510OperatingDefaultPrevious], b510OperatingDefaultKey) {
				return nil, nil, b510OperatingDefaultFmt.Errorf("ab48 catalog: duplicate key")
			}
		}
	}
	return b510OperatingDefaultRows, b510OperatingDefaultKeys, nil
}

func B510OperatingDefaultEntries() (b510OperatingDefaultJSON.RawMessage, error) {
	b510OperatingDefaultSource := []byte(b510OperatingDefaultData)
	if _, _, b510OperatingDefaultErr := b510OperatingDefaultValidate(b510OperatingDefaultSource); b510OperatingDefaultErr != nil {
		return nil, b510OperatingDefaultErr
	}
	return b510OperatingDefaultSource, nil
}

func b510OperatingDefaultLookupSource(b510OperatingDefaultSource []byte, b510OperatingDefaultKey B510OperatingDefaultKey) (b510OperatingDefaultJSON.RawMessage, bool, error) {
	b510OperatingDefaultRows, b510OperatingDefaultKeys, b510OperatingDefaultErr := b510OperatingDefaultValidate(b510OperatingDefaultSource)
	if b510OperatingDefaultErr != nil {
		return nil, false, b510OperatingDefaultErr
	}
	for b510OperatingDefaultIndex, b510OperatingDefaultCandidate := range b510OperatingDefaultKeys {
		if b510OperatingDefaultReflect.DeepEqual(b510OperatingDefaultCandidate, b510OperatingDefaultKey) {
			return b510OperatingDefaultBytes.Clone(b510OperatingDefaultRows[b510OperatingDefaultIndex]), true, nil
		}
	}
	return nil, false, nil
}

func B510OperatingDefaultLookup(b510OperatingDefaultKey B510OperatingDefaultKey) (b510OperatingDefaultJSON.RawMessage, bool, error) {
	return b510OperatingDefaultLookupSource([]byte(b510OperatingDefaultData), b510OperatingDefaultKey)
}
