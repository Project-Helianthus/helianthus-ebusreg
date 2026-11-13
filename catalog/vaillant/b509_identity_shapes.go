package vaillant

import (
	b509IdentityShapeBytes "bytes"
	_ "embed"
	b509IdentityShapeJSON "encoding/json"
	b509IdentityShapeFmt "fmt"
	b509IdentityShapeIO "io"
	b509IdentityShapeReflect "reflect"
)

type B509IdentityShapeKey struct {
	B509IdentityShapeFamily        string  `json:"family"`
	B509IdentityShapeVariant       string  `json:"variant"`
	B509IdentityShapePrimary       *uint16 `json:"primary"`
	B509IdentityShapeSecondary     *uint16 `json:"secondary"`
	B509IdentityShapeOperation     *uint16 `json:"operation"`
	B509IdentityShapeDirection     string  `json:"direction"`
	B509IdentityShapeField         string  `json:"field"`
	B509IdentityShapeAlternative   string  `json:"alternative"`
	B509IdentityShapeGroup         *uint16 `json:"group"`
	B509IdentityShapeInstanceScope string  `json:"instance_scope"`
	B509IdentityShapeRegister      *uint16 `json:"register"`
}

//go:embed data/ab46_b509_identity_shapes.json
var b509IdentityShapeData string

func b509IdentityShapeValidate(b509IdentityShapeSource []byte) ([]b509IdentityShapeJSON.RawMessage, []B509IdentityShapeKey, error) {
	var b509IdentityShapeRows []b509IdentityShapeJSON.RawMessage
	if b509IdentityShapeError := b509IdentityShapeJSON.Unmarshal(b509IdentityShapeSource, &b509IdentityShapeRows); b509IdentityShapeError != nil {
		return nil, nil, b509IdentityShapeError
	}
	if b509IdentityShapeRows == nil {
		return nil, nil, b509IdentityShapeFmt.Errorf("ab46: expected an array")
	}
	b509IdentityShapeKeys := make([]B509IdentityShapeKey, 0, len(b509IdentityShapeRows))
	for b509IdentityShapeIndex, b509IdentityShapeRaw := range b509IdentityShapeRows {
		var b509IdentityShapeObject map[string]b509IdentityShapeJSON.RawMessage
		if b509IdentityShapeError := b509IdentityShapeJSON.Unmarshal(b509IdentityShapeRaw, &b509IdentityShapeObject); b509IdentityShapeError != nil || b509IdentityShapeObject == nil {
			return nil, nil, b509IdentityShapeFmt.Errorf("ab46: row %d must be an object", b509IdentityShapeIndex)
		}
		var b509IdentityShapeKey B509IdentityShapeKey
		b509IdentityShapeDecoder := b509IdentityShapeJSON.NewDecoder(b509IdentityShapeBytes.NewReader(b509IdentityShapeObject["key"]))
		b509IdentityShapeDecoder.DisallowUnknownFields()
		if b509IdentityShapeError := b509IdentityShapeDecoder.Decode(&b509IdentityShapeKey); b509IdentityShapeError != nil {
			return nil, nil, b509IdentityShapeError
		}
		if b509IdentityShapeError := b509IdentityShapeDecoder.Decode(new(any)); b509IdentityShapeError != b509IdentityShapeIO.EOF {
			return nil, nil, b509IdentityShapeFmt.Errorf("ab46: invalid key")
		}
		if b509IdentityShapeKey.B509IdentityShapeFamily == "" || b509IdentityShapeKey.B509IdentityShapeVariant == "" || b509IdentityShapeKey.B509IdentityShapePrimary == nil || b509IdentityShapeKey.B509IdentityShapeDirection == "" || b509IdentityShapeKey.B509IdentityShapeField == "" || b509IdentityShapeKey.B509IdentityShapeAlternative == "" || b509IdentityShapeKey.B509IdentityShapeInstanceScope == "" {
			return nil, nil, b509IdentityShapeFmt.Errorf("ab46: incomplete key at row %d", b509IdentityShapeIndex)
		}
		for _, b509IdentityShapePrevious := range b509IdentityShapeKeys {
			if b509IdentityShapeReflect.DeepEqual(b509IdentityShapePrevious, b509IdentityShapeKey) {
				return nil, nil, b509IdentityShapeFmt.Errorf("ab46: duplicate key at row %d", b509IdentityShapeIndex)
			}
		}
		b509IdentityShapeKeys = append(b509IdentityShapeKeys, b509IdentityShapeKey)
	}
	return b509IdentityShapeRows, b509IdentityShapeKeys, nil
}

func B509IdentityShapeEntries() (b509IdentityShapeJSON.RawMessage, error) {
	if _, _, b509IdentityShapeError := b509IdentityShapeValidate([]byte(b509IdentityShapeData)); b509IdentityShapeError != nil {
		return nil, b509IdentityShapeError
	}
	return b509IdentityShapeJSON.RawMessage([]byte(b509IdentityShapeData)), nil
}

func b509IdentityShapeLookup(b509IdentityShapeSource []byte, b509IdentityShapeKey B509IdentityShapeKey) (b509IdentityShapeJSON.RawMessage, bool, error) {
	b509IdentityShapeRows, b509IdentityShapeKeys, b509IdentityShapeError := b509IdentityShapeValidate(b509IdentityShapeSource)
	if b509IdentityShapeError != nil {
		return nil, false, b509IdentityShapeError
	}
	for b509IdentityShapeIndex, b509IdentityShapeCandidate := range b509IdentityShapeKeys {
		if b509IdentityShapeReflect.DeepEqual(b509IdentityShapeCandidate, b509IdentityShapeKey) {
			return append(b509IdentityShapeJSON.RawMessage(nil), b509IdentityShapeRows[b509IdentityShapeIndex]...), true, nil
		}
	}
	return nil, false, nil
}

func B509IdentityShapeLookup(b509IdentityShapeKey B509IdentityShapeKey) (b509IdentityShapeJSON.RawMessage, bool, error) {
	return b509IdentityShapeLookup([]byte(b509IdentityShapeData), b509IdentityShapeKey)
}
