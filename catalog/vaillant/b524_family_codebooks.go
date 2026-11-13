package vaillant

import (
	_ "embed"
	b524FamilyCodebookJSON "encoding/json"
	b524FamilyCodebookFmt "fmt"
	b524FamilyCodebookReflect "reflect"
)

//go:embed data/ab30_b524_family_codebooks.json
var b524FamilyCodebookEmbedded string

type B524FamilyCodebookKey struct {
	B524Family                      string  `json:"family"`
	B524FamilyCodebookVariant       string  `json:"variant"`
	B524FamilyCodebookPrimary       uint8   `json:"primary"`
	B524FamilyCodebookSecondary     *uint8  `json:"secondary"`
	B524FamilyCodebookOperation     *uint8  `json:"operation"`
	B524FamilyCodebookDirection     string  `json:"direction"`
	B524FamilyCodebookField         string  `json:"field"`
	B524FamilyCodebookAlternative   string  `json:"alternative"`
	B524FamilyCodebookGroup         *uint8  `json:"group"`
	B524FamilyCodebookInstanceScope string  `json:"instance_scope"`
	B524FamilyCodebookRegister      *uint16 `json:"register"`
}

func b524FamilyCodebookValidate(b524FamilyCodebookData []byte) ([]b524FamilyCodebookJSON.RawMessage, []B524FamilyCodebookKey, error) {
	var b524FamilyCodebookRows []b524FamilyCodebookJSON.RawMessage
	if b524FamilyCodebookError := b524FamilyCodebookJSON.Unmarshal(b524FamilyCodebookData, &b524FamilyCodebookRows); b524FamilyCodebookError != nil {
		return nil, nil, b524FamilyCodebookError
	}
	if b524FamilyCodebookRows == nil {
		return nil, nil, b524FamilyCodebookFmt.Errorf("ab30 entries must be an array")
	}
	b524FamilyCodebookKeys := make([]B524FamilyCodebookKey, 0, len(b524FamilyCodebookRows))
	for _, b524FamilyCodebookRow := range b524FamilyCodebookRows {
		var b524FamilyCodebookEntry struct {
			B524FamilyCodebookKey B524FamilyCodebookKey `json:"key"`
		}
		if b524FamilyCodebookError := b524FamilyCodebookJSON.Unmarshal(b524FamilyCodebookRow, &b524FamilyCodebookEntry); b524FamilyCodebookError != nil {
			return nil, nil, b524FamilyCodebookError
		}
		var b524FamilyCodebookObject map[string]b524FamilyCodebookJSON.RawMessage
		if b524FamilyCodebookError := b524FamilyCodebookJSON.Unmarshal(b524FamilyCodebookRow, &b524FamilyCodebookObject); b524FamilyCodebookError != nil {
			return nil, nil, b524FamilyCodebookError
		}
		var b524FamilyCodebookIdentity map[string]b524FamilyCodebookJSON.RawMessage
		if b524FamilyCodebookError := b524FamilyCodebookJSON.Unmarshal(b524FamilyCodebookObject["key"], &b524FamilyCodebookIdentity); b524FamilyCodebookError != nil {
			return nil, nil, b524FamilyCodebookError
		}
		for _, b524FamilyCodebookAxis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			b524FamilyCodebookValue, b524FamilyCodebookPresent := b524FamilyCodebookIdentity[b524FamilyCodebookAxis]
			if !b524FamilyCodebookPresent || (string(b524FamilyCodebookValue) == "null" && b524FamilyCodebookAxis != "secondary" && b524FamilyCodebookAxis != "operation" && b524FamilyCodebookAxis != "group" && b524FamilyCodebookAxis != "register") {
				return nil, nil, b524FamilyCodebookFmt.Errorf("ab30 missing identity axis %s", b524FamilyCodebookAxis)
			}
		}
		if b524FamilyCodebookEntry.B524FamilyCodebookKey.B524Family == "" || b524FamilyCodebookEntry.B524FamilyCodebookKey.B524FamilyCodebookVariant == "" || b524FamilyCodebookEntry.B524FamilyCodebookKey.B524FamilyCodebookDirection == "" || b524FamilyCodebookEntry.B524FamilyCodebookKey.B524FamilyCodebookField == "" || b524FamilyCodebookEntry.B524FamilyCodebookKey.B524FamilyCodebookAlternative == "" || b524FamilyCodebookEntry.B524FamilyCodebookKey.B524FamilyCodebookInstanceScope == "" {
			return nil, nil, b524FamilyCodebookFmt.Errorf("ab30 string axes require explicit tokens")
		}
		for _, b524FamilyCodebookPrior := range b524FamilyCodebookKeys {
			if b524FamilyCodebookReflect.DeepEqual(b524FamilyCodebookPrior, b524FamilyCodebookEntry.B524FamilyCodebookKey) {
				return nil, nil, b524FamilyCodebookFmt.Errorf("ab30 duplicate identity")
			}
		}
		b524FamilyCodebookKeys = append(b524FamilyCodebookKeys, b524FamilyCodebookEntry.B524FamilyCodebookKey)
	}
	return b524FamilyCodebookRows, b524FamilyCodebookKeys, nil
}

func B524FamilyCodebookEntries() (b524FamilyCodebookJSON.RawMessage, error) {
	if _, _, b524FamilyCodebookError := b524FamilyCodebookValidate([]byte(b524FamilyCodebookEmbedded)); b524FamilyCodebookError != nil {
		return nil, b524FamilyCodebookError
	}
	return b524FamilyCodebookJSON.RawMessage([]byte(b524FamilyCodebookEmbedded)), nil
}

func b524FamilyCodebookLookup(b524FamilyCodebookData []byte, b524FamilyCodebookKey B524FamilyCodebookKey) (b524FamilyCodebookJSON.RawMessage, bool, error) {
	b524FamilyCodebookRows, b524FamilyCodebookKeys, b524FamilyCodebookError := b524FamilyCodebookValidate(b524FamilyCodebookData)
	if b524FamilyCodebookError != nil {
		return nil, false, b524FamilyCodebookError
	}
	for b524FamilyCodebookIndex, b524FamilyCodebookCandidate := range b524FamilyCodebookKeys {
		if b524FamilyCodebookReflect.DeepEqual(b524FamilyCodebookCandidate, b524FamilyCodebookKey) {
			return append(b524FamilyCodebookJSON.RawMessage(nil), b524FamilyCodebookRows[b524FamilyCodebookIndex]...), true, nil
		}
	}
	return nil, false, nil
}

func B524FamilyCodebookLookup(b524FamilyCodebookKey B524FamilyCodebookKey) (b524FamilyCodebookJSON.RawMessage, bool, error) {
	return b524FamilyCodebookLookup([]byte(b524FamilyCodebookEmbedded), b524FamilyCodebookKey)
}
