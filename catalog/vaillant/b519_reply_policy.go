package vaillant

import (
	b519ReplyPolicyEmbed "embed"
	b519ReplyPolicyJSON "encoding/json"
	b519ReplyPolicyFmt "fmt"
)

//go:embed data/ab21_b519_reply_policy.json
var b519ReplyPolicyFiles b519ReplyPolicyEmbed.FS

type B519ReplyPolicyKey struct {
	B519ReplyPolicyFamily        string  `json:"family"`
	B519ReplyPolicyVariant       string  `json:"variant"`
	B519ReplyPolicyPrimary       uint8   `json:"primary"`
	B519ReplyPolicySecondary     *uint8  `json:"secondary"`
	B519ReplyPolicyOperation     *uint8  `json:"operation"`
	B519ReplyPolicyDirection     string  `json:"direction"`
	B519ReplyPolicyField         string  `json:"field"`
	B519ReplyPolicyAlternative   string  `json:"alternative"`
	B519ReplyPolicyGroup         *uint8  `json:"group"`
	B519ReplyPolicyInstanceScope string  `json:"instance_scope"`
	B519ReplyPolicyRegister      *uint16 `json:"register"`
}

func b519ReplyPolicyValidate(b519ReplyPolicyData []byte) ([]b519ReplyPolicyJSON.RawMessage, error) {
	var b519ReplyPolicyRows []b519ReplyPolicyJSON.RawMessage
	if b519ReplyPolicyError := b519ReplyPolicyJSON.Unmarshal(b519ReplyPolicyData, &b519ReplyPolicyRows); b519ReplyPolicyError != nil {
		return nil, b519ReplyPolicyError
	}
	if len(b519ReplyPolicyRows) == 0 {
		return nil, b519ReplyPolicyFmt.Errorf("ab21: empty catalog")
	}
	b519ReplyPolicySeen := make(map[string]bool)
	for b519ReplyPolicyIndex, b519ReplyPolicyRow := range b519ReplyPolicyRows {
		var b519ReplyPolicyObject map[string]b519ReplyPolicyJSON.RawMessage
		if b519ReplyPolicyError := b519ReplyPolicyJSON.Unmarshal(b519ReplyPolicyRow, &b519ReplyPolicyObject); b519ReplyPolicyError != nil {
			return nil, b519ReplyPolicyError
		}
		var b519ReplyPolicyAxes map[string]b519ReplyPolicyJSON.RawMessage
		if b519ReplyPolicyError := b519ReplyPolicyJSON.Unmarshal(b519ReplyPolicyObject["key"], &b519ReplyPolicyAxes); b519ReplyPolicyError != nil {
			return nil, b519ReplyPolicyError
		}
		if len(b519ReplyPolicyAxes) != 11 {
			return nil, b519ReplyPolicyFmt.Errorf("ab21: incomplete identity at %d", b519ReplyPolicyIndex)
		}
		for _, b519ReplyPolicyAxis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			if _, b519ReplyPolicyPresent := b519ReplyPolicyAxes[b519ReplyPolicyAxis]; !b519ReplyPolicyPresent {
				return nil, b519ReplyPolicyFmt.Errorf("ab21: missing axis %s", b519ReplyPolicyAxis)
			}
		}
		var b519ReplyPolicyKey B519ReplyPolicyKey
		if b519ReplyPolicyError := b519ReplyPolicyJSON.Unmarshal(b519ReplyPolicyObject["key"], &b519ReplyPolicyKey); b519ReplyPolicyError != nil {
			return nil, b519ReplyPolicyError
		}
		if b519ReplyPolicyKey.B519ReplyPolicyFamily == "" || b519ReplyPolicyKey.B519ReplyPolicyVariant == "" || b519ReplyPolicyKey.B519ReplyPolicyDirection == "" || b519ReplyPolicyKey.B519ReplyPolicyField == "" || b519ReplyPolicyKey.B519ReplyPolicyAlternative == "" || b519ReplyPolicyKey.B519ReplyPolicyInstanceScope == "" || string(b519ReplyPolicyAxes["primary"]) == "null" {
			return nil, b519ReplyPolicyFmt.Errorf("ab21: invalid identity at %d", b519ReplyPolicyIndex)
		}
		b519ReplyPolicyIdentity, b519ReplyPolicyError := b519ReplyPolicyJSON.Marshal(b519ReplyPolicyKey)
		if b519ReplyPolicyError != nil {
			return nil, b519ReplyPolicyError
		}
		if b519ReplyPolicySeen[string(b519ReplyPolicyIdentity)] {
			return nil, b519ReplyPolicyFmt.Errorf("ab21: duplicate identity at %d", b519ReplyPolicyIndex)
		}
		b519ReplyPolicySeen[string(b519ReplyPolicyIdentity)] = true
	}
	return b519ReplyPolicyRows, nil
}

func B519ReplyPolicyEntries() (b519ReplyPolicyJSON.RawMessage, error) {
	b519ReplyPolicyData, b519ReplyPolicyError := b519ReplyPolicyFiles.ReadFile("data/ab21_b519_reply_policy.json")
	if b519ReplyPolicyError != nil {
		return nil, b519ReplyPolicyError
	}
	if _, b519ReplyPolicyError := b519ReplyPolicyValidate(b519ReplyPolicyData); b519ReplyPolicyError != nil {
		return nil, b519ReplyPolicyError
	}
	return append(b519ReplyPolicyJSON.RawMessage(nil), b519ReplyPolicyData...), nil
}

func b519ReplyPolicyLookup(b519ReplyPolicyData []byte, b519ReplyPolicyKey B519ReplyPolicyKey) (b519ReplyPolicyJSON.RawMessage, bool, error) {
	b519ReplyPolicyRows, b519ReplyPolicyError := b519ReplyPolicyValidate(b519ReplyPolicyData)
	if b519ReplyPolicyError != nil {
		return nil, false, b519ReplyPolicyError
	}
	b519ReplyPolicyWanted, b519ReplyPolicyError := b519ReplyPolicyJSON.Marshal(b519ReplyPolicyKey)
	if b519ReplyPolicyError != nil {
		return nil, false, b519ReplyPolicyError
	}
	for _, b519ReplyPolicyRow := range b519ReplyPolicyRows {
		var b519ReplyPolicyEntry struct {
			B519ReplyPolicyKey B519ReplyPolicyKey `json:"key"`
		}
		if b519ReplyPolicyError := b519ReplyPolicyJSON.Unmarshal(b519ReplyPolicyRow, &b519ReplyPolicyEntry); b519ReplyPolicyError != nil {
			return nil, false, b519ReplyPolicyError
		}
		b519ReplyPolicyIdentity, b519ReplyPolicyError := b519ReplyPolicyJSON.Marshal(b519ReplyPolicyEntry.B519ReplyPolicyKey)
		if b519ReplyPolicyError != nil {
			return nil, false, b519ReplyPolicyError
		}
		if string(b519ReplyPolicyWanted) == string(b519ReplyPolicyIdentity) {
			return append(b519ReplyPolicyJSON.RawMessage(nil), b519ReplyPolicyRow...), true, nil
		}
	}
	return nil, false, nil
}

func B519ReplyPolicyLookup(b519ReplyPolicyKey B519ReplyPolicyKey) (b519ReplyPolicyJSON.RawMessage, bool, error) {
	b519ReplyPolicyData, b519ReplyPolicyError := b519ReplyPolicyFiles.ReadFile("data/ab21_b519_reply_policy.json")
	if b519ReplyPolicyError != nil {
		return nil, false, b519ReplyPolicyError
	}
	return b519ReplyPolicyLookup(b519ReplyPolicyData, b519ReplyPolicyKey)
}
