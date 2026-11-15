package vaillant

import (
	_ "embed"
	b512FeedbackFallbackJSON "encoding/json"
	b512FeedbackFallbackFmt "fmt"
)

// B512FeedbackFallbackKey identifies one B512 feedback fallback row; absent optional axes are nil and
// every axis takes part in the lookup.
type B512FeedbackFallbackKey struct {
	B512FeedbackFallbackFamily        string  `json:"family"`
	B512FeedbackFallbackVariant       string  `json:"variant"`
	B512FeedbackFallbackPrimary       uint8   `json:"primary"`
	B512FeedbackFallbackSecondary     *uint8  `json:"secondary"`
	B512FeedbackFallbackOperation     *uint8  `json:"operation"`
	B512FeedbackFallbackDirection     string  `json:"direction"`
	B512FeedbackFallbackField         string  `json:"field"`
	B512FeedbackFallbackAlternative   string  `json:"alternative"`
	B512FeedbackFallbackGroup         *uint8  `json:"group"`
	B512FeedbackFallbackInstanceScope string  `json:"instance_scope"`
	B512FeedbackFallbackRegister      *uint16 `json:"register"`
}

//go:embed data/b512_feedback_fallback.json
var b512FeedbackFallbackData string

func b512FeedbackFallbackValidate(b512FeedbackFallbackInput []byte) (map[string]b512FeedbackFallbackJSON.RawMessage, error) {
	var b512FeedbackFallbackRows []b512FeedbackFallbackJSON.RawMessage
	if b512FeedbackFallbackError := b512FeedbackFallbackJSON.Unmarshal(b512FeedbackFallbackInput, &b512FeedbackFallbackRows); b512FeedbackFallbackError != nil {
		return nil, b512FeedbackFallbackError
	}
	if b512FeedbackFallbackRows == nil {
		return nil, b512FeedbackFallbackFmt.Errorf("b512 feedback fallback: expected JSON array")
	}
	b512FeedbackFallbackIndex := make(map[string]b512FeedbackFallbackJSON.RawMessage, len(b512FeedbackFallbackRows))
	for _, b512FeedbackFallbackRow := range b512FeedbackFallbackRows {
		var b512FeedbackFallbackEnvelope struct {
			B512FeedbackFallbackKey b512FeedbackFallbackJSON.RawMessage `json:"key"`
		}
		if b512FeedbackFallbackError := b512FeedbackFallbackJSON.Unmarshal(b512FeedbackFallbackRow, &b512FeedbackFallbackEnvelope); b512FeedbackFallbackError != nil {
			return nil, b512FeedbackFallbackError
		}
		var b512FeedbackFallbackAxes map[string]b512FeedbackFallbackJSON.RawMessage
		if b512FeedbackFallbackError := b512FeedbackFallbackJSON.Unmarshal(b512FeedbackFallbackEnvelope.B512FeedbackFallbackKey, &b512FeedbackFallbackAxes); b512FeedbackFallbackError != nil {
			return nil, b512FeedbackFallbackError
		}
		for _, b512FeedbackFallbackAxis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			if _, b512FeedbackFallbackPresent := b512FeedbackFallbackAxes[b512FeedbackFallbackAxis]; !b512FeedbackFallbackPresent {
				return nil, b512FeedbackFallbackFmt.Errorf("b512 feedback fallback: missing identity axis %s", b512FeedbackFallbackAxis)
			}
		}
		if len(b512FeedbackFallbackAxes) != 11 {
			return nil, b512FeedbackFallbackFmt.Errorf("b512 feedback fallback: unexpected identity axes")
		}
		var b512FeedbackFallbackKey B512FeedbackFallbackKey
		if b512FeedbackFallbackError := b512FeedbackFallbackJSON.Unmarshal(b512FeedbackFallbackEnvelope.B512FeedbackFallbackKey, &b512FeedbackFallbackKey); b512FeedbackFallbackError != nil {
			return nil, b512FeedbackFallbackError
		}
		if b512FeedbackFallbackKey.B512FeedbackFallbackFamily == "" || b512FeedbackFallbackKey.B512FeedbackFallbackVariant == "" || b512FeedbackFallbackKey.B512FeedbackFallbackDirection == "" || b512FeedbackFallbackKey.B512FeedbackFallbackField == "" || b512FeedbackFallbackKey.B512FeedbackFallbackAlternative == "" || b512FeedbackFallbackKey.B512FeedbackFallbackInstanceScope == "" || string(b512FeedbackFallbackAxes["primary"]) == "null" {
			return nil, b512FeedbackFallbackFmt.Errorf("b512 feedback fallback: incomplete identity")
		}
		b512FeedbackFallbackIdentity, b512FeedbackFallbackError := b512FeedbackFallbackJSON.Marshal(b512FeedbackFallbackKey)
		if b512FeedbackFallbackError != nil {
			return nil, b512FeedbackFallbackError
		}
		if _, b512FeedbackFallbackDuplicate := b512FeedbackFallbackIndex[string(b512FeedbackFallbackIdentity)]; b512FeedbackFallbackDuplicate {
			return nil, b512FeedbackFallbackFmt.Errorf("b512 feedback fallback: duplicate identity")
		}
		b512FeedbackFallbackIndex[string(b512FeedbackFallbackIdentity)] = append(b512FeedbackFallbackJSON.RawMessage(nil), b512FeedbackFallbackRow...)
	}
	return b512FeedbackFallbackIndex, nil
}

// B512FeedbackFallbackEntries validates the embedded B512 feedback fallback rows and returns them as
// caller-owned JSON. Invalid embedded data returns an error.
func B512FeedbackFallbackEntries() (b512FeedbackFallbackJSON.RawMessage, error) {
	if _, b512FeedbackFallbackError := b512FeedbackFallbackValidate([]byte(b512FeedbackFallbackData)); b512FeedbackFallbackError != nil {
		return nil, b512FeedbackFallbackError
	}
	return append(b512FeedbackFallbackJSON.RawMessage(nil), b512FeedbackFallbackData...), nil
}

func b512FeedbackFallbackLookup(b512FeedbackFallbackInput []byte, b512FeedbackFallbackKey B512FeedbackFallbackKey) (b512FeedbackFallbackJSON.RawMessage, bool, error) {
	b512FeedbackFallbackIndex, b512FeedbackFallbackError := b512FeedbackFallbackValidate(b512FeedbackFallbackInput)
	if b512FeedbackFallbackError != nil {
		return nil, false, b512FeedbackFallbackError
	}
	b512FeedbackFallbackIdentity, b512FeedbackFallbackError := b512FeedbackFallbackJSON.Marshal(b512FeedbackFallbackKey)
	if b512FeedbackFallbackError != nil {
		return nil, false, b512FeedbackFallbackError
	}
	b512FeedbackFallbackRow, b512FeedbackFallbackFound := b512FeedbackFallbackIndex[string(b512FeedbackFallbackIdentity)]
	return append(b512FeedbackFallbackJSON.RawMessage(nil), b512FeedbackFallbackRow...), b512FeedbackFallbackFound, nil
}

// B512FeedbackFallbackLookup returns the B512 feedback fallback row whose key equals the given key on
// every axis as caller-owned JSON. It returns an error if the embedded rows fail validation.
func B512FeedbackFallbackLookup(b512FeedbackFallbackKey B512FeedbackFallbackKey) (b512FeedbackFallbackJSON.RawMessage, bool, error) {
	return b512FeedbackFallbackLookup([]byte(b512FeedbackFallbackData), b512FeedbackFallbackKey)
}
