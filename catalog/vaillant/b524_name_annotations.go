package vaillant

import (
	b524NameAnnotationBytes "bytes"
	_ "embed"
	b524NameAnnotationJSON "encoding/json"
	b524NameAnnotationFmt "fmt"
	b524NameAnnotationSync "sync"
)

//go:embed data/ab36_b524_name_annotations.json
var b524NameAnnotationData string

type b524NameAnnotationCatalog struct {
	b524NameAnnotationSource string
	b524NameAnnotationOnce   b524NameAnnotationSync.Once
	b524NameAnnotationIndex  map[string]b524NameAnnotationJSON.RawMessage
	b524NameAnnotationError  error
}

var b524NameAnnotationCachedCatalog = &b524NameAnnotationCatalog{b524NameAnnotationSource: b524NameAnnotationData}

func (b524NameAnnotationCatalogValue *b524NameAnnotationCatalog) b524NameAnnotationLoad() (map[string]b524NameAnnotationJSON.RawMessage, error) {
	b524NameAnnotationCatalogValue.b524NameAnnotationOnce.Do(func() {
		_, b524NameAnnotationCatalogValue.b524NameAnnotationIndex, b524NameAnnotationCatalogValue.b524NameAnnotationError = b524NameAnnotationValidate([]byte(b524NameAnnotationCatalogValue.b524NameAnnotationSource))
	})
	return b524NameAnnotationCatalogValue.b524NameAnnotationIndex, b524NameAnnotationCatalogValue.b524NameAnnotationError
}

type B524NameAnnotationKey struct {
	B524NameAnnotationFamily        string  `json:"family"`
	B524NameAnnotationVariant       string  `json:"variant"`
	B524NameAnnotationPrimary       uint8   `json:"primary"`
	B524NameAnnotationSecondary     *uint8  `json:"secondary"`
	B524NameAnnotationOperation     *uint8  `json:"operation"`
	B524NameAnnotationDirection     string  `json:"direction"`
	B524NameAnnotationField         string  `json:"field"`
	B524NameAnnotationAlternative   string  `json:"alternative"`
	B524NameAnnotationGroup         *uint8  `json:"group"`
	B524NameAnnotationInstanceScope string  `json:"instance_scope"`
	B524NameAnnotationRegister      *uint16 `json:"register"`
}

func b524NameAnnotationValidate(b524NameAnnotationSource []byte) ([]b524NameAnnotationJSON.RawMessage, map[string]b524NameAnnotationJSON.RawMessage, error) {
	var b524NameAnnotationRows []b524NameAnnotationJSON.RawMessage
	if b524NameAnnotationError := b524NameAnnotationJSON.Unmarshal(b524NameAnnotationSource, &b524NameAnnotationRows); b524NameAnnotationError != nil {
		return nil, nil, b524NameAnnotationError
	}
	if b524NameAnnotationRows == nil {
		return nil, nil, b524NameAnnotationFmt.Errorf("ab36: expected an array")
	}
	b524NameAnnotationIndex := make(map[string]b524NameAnnotationJSON.RawMessage, len(b524NameAnnotationRows))
	for _, b524NameAnnotationRow := range b524NameAnnotationRows {
		var b524NameAnnotationObject map[string]b524NameAnnotationJSON.RawMessage
		if b524NameAnnotationError := b524NameAnnotationJSON.Unmarshal(b524NameAnnotationRow, &b524NameAnnotationObject); b524NameAnnotationError != nil {
			return nil, nil, b524NameAnnotationError
		}
		for _, documentaryField := range []string{"legacy_name", "display_alias", "previous_name"} {
			if _, b524NameAnnotationPresent := b524NameAnnotationObject[documentaryField]; b524NameAnnotationPresent {
				return nil, nil, b524NameAnnotationFmt.Errorf("ab36: documentary name field is not runtime data")
			}
		}
		b524NameAnnotationRawKey, b524NameAnnotationPresent := b524NameAnnotationObject["key"]
		if !b524NameAnnotationPresent || b524NameAnnotationBytes.Equal(b524NameAnnotationBytes.TrimSpace(b524NameAnnotationRawKey), []byte("null")) {
			return nil, nil, b524NameAnnotationFmt.Errorf("ab36: missing key")
		}
		var b524NameAnnotationKey B524NameAnnotationKey
		b524NameAnnotationDecoder := b524NameAnnotationJSON.NewDecoder(b524NameAnnotationBytes.NewReader(b524NameAnnotationRawKey))
		b524NameAnnotationDecoder.DisallowUnknownFields()
		if b524NameAnnotationError := b524NameAnnotationDecoder.Decode(&b524NameAnnotationKey); b524NameAnnotationError != nil {
			return nil, nil, b524NameAnnotationError
		}
		var b524NameAnnotationFields map[string]b524NameAnnotationJSON.RawMessage
		if b524NameAnnotationError := b524NameAnnotationJSON.Unmarshal(b524NameAnnotationRawKey, &b524NameAnnotationFields); b524NameAnnotationError != nil {
			return nil, nil, b524NameAnnotationError
		}
		b524NameAnnotationRawPrimary, b524NameAnnotationPresent := b524NameAnnotationFields["primary"]
		if !b524NameAnnotationPresent || b524NameAnnotationKey.B524NameAnnotationFamily == "" || b524NameAnnotationKey.B524NameAnnotationVariant == "" || b524NameAnnotationKey.B524NameAnnotationDirection == "" || b524NameAnnotationKey.B524NameAnnotationField == "" || b524NameAnnotationKey.B524NameAnnotationAlternative == "" || b524NameAnnotationKey.B524NameAnnotationInstanceScope == "" || b524NameAnnotationBytes.Equal(b524NameAnnotationBytes.TrimSpace(b524NameAnnotationRawPrimary), []byte("null")) {
			return nil, nil, b524NameAnnotationFmt.Errorf("ab36: incomplete key")
		}
		b524NameAnnotationEncoded, b524NameAnnotationError := b524NameAnnotationJSON.Marshal(b524NameAnnotationKey)
		if b524NameAnnotationError != nil {
			return nil, nil, b524NameAnnotationError
		}
		if _, b524NameAnnotationPresent := b524NameAnnotationIndex[string(b524NameAnnotationEncoded)]; b524NameAnnotationPresent {
			return nil, nil, b524NameAnnotationFmt.Errorf("ab36: duplicate key")
		}
		b524NameAnnotationIndex[string(b524NameAnnotationEncoded)] = b524NameAnnotationRow
	}
	return b524NameAnnotationRows, b524NameAnnotationIndex, nil
}

func B524NameAnnotationEntries() (b524NameAnnotationJSON.RawMessage, error) {
	if _, b524NameAnnotationError := b524NameAnnotationCachedCatalog.b524NameAnnotationLoad(); b524NameAnnotationError != nil {
		return nil, b524NameAnnotationError
	}
	return b524NameAnnotationJSON.RawMessage(b524NameAnnotationCachedCatalog.b524NameAnnotationSource), nil
}

func B524NameAnnotationLookup(b524NameAnnotationKey B524NameAnnotationKey) (b524NameAnnotationJSON.RawMessage, bool, error) {
	b524NameAnnotationIndex, b524NameAnnotationError := b524NameAnnotationCachedCatalog.b524NameAnnotationLoad()
	if b524NameAnnotationError != nil {
		return nil, false, b524NameAnnotationError
	}
	b524NameAnnotationEncoded, b524NameAnnotationError := b524NameAnnotationJSON.Marshal(b524NameAnnotationKey)
	if b524NameAnnotationError != nil {
		return nil, false, b524NameAnnotationError
	}
	b524NameAnnotationRow, b524NameAnnotationFound := b524NameAnnotationIndex[string(b524NameAnnotationEncoded)]
	if !b524NameAnnotationFound {
		return nil, false, nil
	}
	return append(b524NameAnnotationJSON.RawMessage(nil), b524NameAnnotationRow...), true, nil
}
