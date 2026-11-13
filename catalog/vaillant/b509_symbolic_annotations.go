package vaillant

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed data/ab40_b509_symbolic_annotations.json
var b509SymbolicAnnotationData []byte

type B509SymbolicAnnotationKey struct {
	B509SymbolicAnnotationFamily        string  `json:"family"`
	B509SymbolicAnnotationVariant       string  `json:"variant"`
	B509SymbolicAnnotationPrimary       uint8   `json:"primary"`
	B509SymbolicAnnotationSecondary     *uint8  `json:"secondary,omitempty"`
	B509SymbolicAnnotationOperation     *uint8  `json:"operation,omitempty"`
	B509SymbolicAnnotationDirection     string  `json:"direction"`
	B509SymbolicAnnotationField         string  `json:"field"`
	B509SymbolicAnnotationAlternative   string  `json:"alternative"`
	B509SymbolicAnnotationGroup         *uint8  `json:"group,omitempty"`
	B509SymbolicAnnotationInstanceScope string  `json:"instance_scope"`
	B509SymbolicAnnotationRegister      *uint16 `json:"register,omitempty"`
	B509SymbolicAnnotationClass         string  `json:"class"`
	B509SymbolicAnnotationSelector      *string `json:"selector,omitempty"`
}

func b509SymbolicAnnotationIdentity(b509SymbolicAnnotationKey B509SymbolicAnnotationKey) string {
	b509SymbolicAnnotationEncoded, _ := json.Marshal(b509SymbolicAnnotationKey)
	return string(b509SymbolicAnnotationEncoded)
}

func b509SymbolicAnnotationValidate(b509SymbolicAnnotationSource []byte) ([]json.RawMessage, map[string]json.RawMessage, error) {
	var b509SymbolicAnnotationRows []json.RawMessage
	if b509SymbolicAnnotationError := json.Unmarshal(b509SymbolicAnnotationSource, &b509SymbolicAnnotationRows); b509SymbolicAnnotationError != nil {
		return nil, nil, b509SymbolicAnnotationError
	}
	if b509SymbolicAnnotationRows == nil {
		return nil, nil, fmt.Errorf("ab40: expected array")
	}
	b509SymbolicAnnotationIndex := make(map[string]json.RawMessage, len(b509SymbolicAnnotationRows))
	for _, b509SymbolicAnnotationRaw := range b509SymbolicAnnotationRows {
		var b509SymbolicAnnotationKey B509SymbolicAnnotationKey
		if b509SymbolicAnnotationError := json.Unmarshal(b509SymbolicAnnotationRaw, &b509SymbolicAnnotationKey); b509SymbolicAnnotationError != nil {
			return nil, nil, b509SymbolicAnnotationError
		}
		var b509SymbolicAnnotationObject map[string]json.RawMessage
		if b509SymbolicAnnotationError := json.Unmarshal(b509SymbolicAnnotationRaw, &b509SymbolicAnnotationObject); b509SymbolicAnnotationError != nil {
			return nil, nil, b509SymbolicAnnotationError
		}
		for _, b509SymbolicAnnotationAxis := range []string{"family", "variant", "primary", "direction", "field", "alternative", "instance_scope", "class"} {
			if b509SymbolicAnnotationValue, b509SymbolicAnnotationPresent := b509SymbolicAnnotationObject[b509SymbolicAnnotationAxis]; !b509SymbolicAnnotationPresent || string(b509SymbolicAnnotationValue) == "null" {
				return nil, nil, fmt.Errorf("ab40: missing axis %s", b509SymbolicAnnotationAxis)
			}
		}
		for _, b509SymbolicAnnotationAxis := range []string{"secondary", "operation", "group", "register", "selector"} {
			if b509SymbolicAnnotationValue, b509SymbolicAnnotationPresent := b509SymbolicAnnotationObject[b509SymbolicAnnotationAxis]; b509SymbolicAnnotationPresent && string(b509SymbolicAnnotationValue) == "null" {
				return nil, nil, fmt.Errorf("ab40: null optional axis %s", b509SymbolicAnnotationAxis)
			}
		}
		if b509SymbolicAnnotationKey.B509SymbolicAnnotationFamily == "" || b509SymbolicAnnotationKey.B509SymbolicAnnotationVariant == "" || b509SymbolicAnnotationKey.B509SymbolicAnnotationDirection == "" || b509SymbolicAnnotationKey.B509SymbolicAnnotationField == "" || b509SymbolicAnnotationKey.B509SymbolicAnnotationAlternative == "" || b509SymbolicAnnotationKey.B509SymbolicAnnotationInstanceScope == "" || b509SymbolicAnnotationKey.B509SymbolicAnnotationClass == "" {
			return nil, nil, fmt.Errorf("ab40: empty string axis; use none when inapplicable")
		}
		b509SymbolicAnnotationID := b509SymbolicAnnotationIdentity(b509SymbolicAnnotationKey)
		if _, b509SymbolicAnnotationPresent := b509SymbolicAnnotationIndex[b509SymbolicAnnotationID]; b509SymbolicAnnotationPresent {
			return nil, nil, fmt.Errorf("ab40: duplicate key")
		}
		b509SymbolicAnnotationIndex[b509SymbolicAnnotationID] = append(json.RawMessage(nil), b509SymbolicAnnotationRaw...)
	}
	return b509SymbolicAnnotationRows, b509SymbolicAnnotationIndex, nil
}

func B509SymbolicAnnotationEntries() (json.RawMessage, error) {
	if _, _, b509SymbolicAnnotationError := b509SymbolicAnnotationValidate(b509SymbolicAnnotationData); b509SymbolicAnnotationError != nil {
		return nil, b509SymbolicAnnotationError
	}
	return append(json.RawMessage(nil), b509SymbolicAnnotationData...), nil
}

func B509SymbolicAnnotationLookup(b509SymbolicAnnotationKey B509SymbolicAnnotationKey) (json.RawMessage, bool, error) {
	_, b509SymbolicAnnotationIndex, b509SymbolicAnnotationError := b509SymbolicAnnotationValidate(b509SymbolicAnnotationData)
	if b509SymbolicAnnotationError != nil {
		return nil, false, b509SymbolicAnnotationError
	}
	b509SymbolicAnnotationRaw, b509SymbolicAnnotationPresent := b509SymbolicAnnotationIndex[b509SymbolicAnnotationIdentity(b509SymbolicAnnotationKey)]
	return append(json.RawMessage(nil), b509SymbolicAnnotationRaw...), b509SymbolicAnnotationPresent, nil
}
