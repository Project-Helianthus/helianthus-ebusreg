package vaillant

import (
	b513SymbolicDictionaryBytes "bytes"
	_ "embed"
	b513SymbolicDictionaryJSON "encoding/json"
	b513SymbolicDictionaryFmt "fmt"
	b513SymbolicDictionaryIO "io"
	b513SymbolicDictionarySync "sync"
)

//go:embed data/ab38_b513_symbolic_dictionary.json
var b513SymbolicDictionaryData []byte

type b513SymbolicDictionaryCache struct {
	b513SymbolicDictionaryOnce    b513SymbolicDictionarySync.Once
	b513SymbolicDictionaryEntries b513SymbolicDictionaryJSON.RawMessage
	b513SymbolicDictionaryIndex   map[string]b513SymbolicDictionaryJSON.RawMessage
	b513SymbolicDictionaryErr     error
}

var b513SymbolicDictionaryDictionary = new(b513SymbolicDictionaryCache)

func b513SymbolicDictionaryLoad() *b513SymbolicDictionaryCache {
	b513SymbolicDictionaryDictionary.b513SymbolicDictionaryOnce.Do(func() {
		_, b513SymbolicDictionaryDictionary.b513SymbolicDictionaryIndex, b513SymbolicDictionaryDictionary.b513SymbolicDictionaryErr = b513SymbolicDictionaryValidate(b513SymbolicDictionaryData)
		if b513SymbolicDictionaryDictionary.b513SymbolicDictionaryErr == nil {
			b513SymbolicDictionaryDictionary.b513SymbolicDictionaryEntries = b513SymbolicDictionaryBytes.Clone(b513SymbolicDictionaryData)
		}
	})
	return b513SymbolicDictionaryDictionary
}

type B513SymbolicDictionaryKey struct {
	B513SymbolicDictionaryFamily        string  `json:"family"`
	B513SymbolicDictionaryVariant       string  `json:"variant"`
	B513SymbolicDictionaryPrimary       uint8   `json:"primary"`
	B513SymbolicDictionarySecondary     *uint8  `json:"secondary"`
	B513SymbolicDictionaryOperation     *uint8  `json:"operation"`
	B513SymbolicDictionaryDirection     string  `json:"direction"`
	B513SymbolicDictionaryField         string  `json:"field"`
	B513SymbolicDictionaryAlternative   string  `json:"alternative"`
	B513SymbolicDictionaryGroup         *uint8  `json:"group"`
	B513SymbolicDictionaryInstanceScope string  `json:"instance_scope"`
	B513SymbolicDictionaryRegister      *uint16 `json:"register"`
	B513SymbolicDictionaryNamespace     string  `json:"dictionary_namespace"`
	B513SymbolicDictionaryClass         string  `json:"dictionary_class"`
	B513SymbolicDictionaryRawAddress    *uint16 `json:"raw_address"`
	B513SymbolicDictionaryBusFamily     string  `json:"bus_family,omitempty"`
}

type b513SymbolicDictionaryRow struct {
	B513SymbolicDictionaryKey
	B513SymbolicDictionaryNativeLabel           string   `json:"native_label"`
	B513SymbolicDictionaryRawForm               string   `json:"raw_form"`
	B513SymbolicDictionaryBusFamilies           []string `json:"bus_families"`
	B513SymbolicDictionaryGrade                 string   `json:"grade"`
	B513SymbolicDictionaryPhysicalQualification string   `json:"physical_qualification"`
	B513SymbolicDictionaryValueCodec            *string  `json:"value_codec"`
	B513SymbolicDictionaryExecutableCapability  bool     `json:"executable_capability"`
	B513SymbolicDictionaryAddressWidthBits      uint8    `json:"address_width_bits"`
	B513SymbolicDictionaryValueWidthBits        *uint8   `json:"value_width_bits"`
	B513SymbolicDictionaryUnknown               []string `json:"unknown"`
	B513SymbolicDictionaryUnknownAddress        string   `json:"unknown_address"`
	B513SymbolicDictionaryLabelMeaning          string   `json:"label_meaning"`
	B513SymbolicDictionaryNativeContract        string   `json:"native_contract"`
}

func b513SymbolicDictionaryValidate(b513SymbolicDictionarySource []byte) ([]b513SymbolicDictionaryJSON.RawMessage, map[string]b513SymbolicDictionaryJSON.RawMessage, error) {
	var b513SymbolicDictionaryRows []b513SymbolicDictionaryJSON.RawMessage
	if b513SymbolicDictionaryErr := b513SymbolicDictionaryJSON.Unmarshal(b513SymbolicDictionarySource, &b513SymbolicDictionaryRows); b513SymbolicDictionaryErr != nil {
		return nil, nil, b513SymbolicDictionaryFmt.Errorf("ab38 dictionary: %w", b513SymbolicDictionaryErr)
	}
	if len(b513SymbolicDictionaryRows) == 0 {
		return nil, nil, b513SymbolicDictionaryFmt.Errorf("ab38 dictionary: expected nonempty array")
	}
	b513SymbolicDictionaryIndex := make(map[string]b513SymbolicDictionaryJSON.RawMessage)
	for _, b513SymbolicDictionaryRaw := range b513SymbolicDictionaryRows {
		if b513SymbolicDictionaryErr := b513SymbolicDictionaryValidateObject(b513SymbolicDictionaryRaw); b513SymbolicDictionaryErr != nil {
			return nil, nil, b513SymbolicDictionaryErr
		}
		var b513SymbolicDictionaryEntry b513SymbolicDictionaryRow
		b513SymbolicDictionaryDecoder := b513SymbolicDictionaryJSON.NewDecoder(b513SymbolicDictionaryBytes.NewReader(b513SymbolicDictionaryRaw))
		b513SymbolicDictionaryDecoder.DisallowUnknownFields()
		if b513SymbolicDictionaryErr := b513SymbolicDictionaryDecoder.Decode(&b513SymbolicDictionaryEntry); b513SymbolicDictionaryErr != nil {
			return nil, nil, b513SymbolicDictionaryFmt.Errorf("ab38 row: %w", b513SymbolicDictionaryErr)
		}
		if b513SymbolicDictionaryErr := b513SymbolicDictionaryDecoder.Decode(new(any)); b513SymbolicDictionaryErr != b513SymbolicDictionaryIO.EOF {
			return nil, nil, b513SymbolicDictionaryFmt.Errorf("ab38 row: trailing data")
		}
		if b513SymbolicDictionaryEntry.B513SymbolicDictionaryFamily == "" || b513SymbolicDictionaryEntry.B513SymbolicDictionaryVariant == "" || b513SymbolicDictionaryEntry.B513SymbolicDictionaryDirection == "" || b513SymbolicDictionaryEntry.B513SymbolicDictionaryField == "" || b513SymbolicDictionaryEntry.B513SymbolicDictionaryAlternative == "" || b513SymbolicDictionaryEntry.B513SymbolicDictionaryInstanceScope == "" || b513SymbolicDictionaryEntry.B513SymbolicDictionaryNamespace == "" || b513SymbolicDictionaryEntry.B513SymbolicDictionaryClass == "" || b513SymbolicDictionaryEntry.B513SymbolicDictionaryRawAddress == nil || b513SymbolicDictionaryEntry.B513SymbolicDictionaryBusFamily != "" || b513SymbolicDictionaryEntry.B513SymbolicDictionaryNativeLabel == "" {
			return nil, nil, b513SymbolicDictionaryFmt.Errorf("ab38 row: incomplete dictionary identity")
		}
		if b513SymbolicDictionaryEntry.B513SymbolicDictionaryGrade != "Hypothesis" || b513SymbolicDictionaryEntry.B513SymbolicDictionaryPhysicalQualification != "Unknown" || b513SymbolicDictionaryEntry.B513SymbolicDictionaryValueCodec != nil || b513SymbolicDictionaryEntry.B513SymbolicDictionaryValueWidthBits != nil || b513SymbolicDictionaryEntry.B513SymbolicDictionaryExecutableCapability || b513SymbolicDictionaryEntry.B513SymbolicDictionaryAddressWidthBits != 16 || len(b513SymbolicDictionaryEntry.B513SymbolicDictionaryBusFamilies) != 2 || b513SymbolicDictionaryEntry.B513SymbolicDictionaryBusFamilies[0] != "eBUS" || b513SymbolicDictionaryEntry.B513SymbolicDictionaryBusFamilies[1] != "PeBUS" || b513SymbolicDictionaryEntry.B513SymbolicDictionaryRawForm != b513SymbolicDictionaryFmt.Sprintf("%04X|%s", *b513SymbolicDictionaryEntry.B513SymbolicDictionaryRawAddress, b513SymbolicDictionaryEntry.B513SymbolicDictionaryNativeLabel) {
			return nil, nil, b513SymbolicDictionaryFmt.Errorf("ab38 row: invalid candidate boundary or raw form")
		}
		for _, b513SymbolicDictionaryBusFamily := range b513SymbolicDictionaryEntry.B513SymbolicDictionaryBusFamilies {
			b513SymbolicDictionaryEntry.B513SymbolicDictionaryBusFamily = b513SymbolicDictionaryBusFamily
			b513SymbolicDictionaryEncoded, b513SymbolicDictionaryErr := b513SymbolicDictionaryJSON.Marshal(b513SymbolicDictionaryEntry.B513SymbolicDictionaryKey)
			if b513SymbolicDictionaryErr != nil {
				return nil, nil, b513SymbolicDictionaryErr
			}
			if _, b513SymbolicDictionaryExists := b513SymbolicDictionaryIndex[string(b513SymbolicDictionaryEncoded)]; b513SymbolicDictionaryExists {
				return nil, nil, b513SymbolicDictionaryFmt.Errorf("ab38 dictionary: duplicate key")
			}
			b513SymbolicDictionaryIndex[string(b513SymbolicDictionaryEncoded)] = b513SymbolicDictionaryBytes.Clone(b513SymbolicDictionaryRaw)
		}
	}
	return b513SymbolicDictionaryRows, b513SymbolicDictionaryIndex, nil
}

func B513SymbolicDictionaryEntries() (b513SymbolicDictionaryJSON.RawMessage, error) {
	b513SymbolicDictionaryCached := b513SymbolicDictionaryLoad()
	if b513SymbolicDictionaryCached.b513SymbolicDictionaryErr != nil {
		return nil, b513SymbolicDictionaryCached.b513SymbolicDictionaryErr
	}
	return b513SymbolicDictionaryBytes.Clone(b513SymbolicDictionaryCached.b513SymbolicDictionaryEntries), nil
}

func B513SymbolicDictionaryLookup(b513SymbolicDictionaryKey B513SymbolicDictionaryKey) (b513SymbolicDictionaryJSON.RawMessage, bool, error) {
	b513SymbolicDictionaryCached := b513SymbolicDictionaryLoad()
	if b513SymbolicDictionaryCached.b513SymbolicDictionaryErr != nil {
		return nil, false, b513SymbolicDictionaryCached.b513SymbolicDictionaryErr
	}
	b513SymbolicDictionaryEncoded, b513SymbolicDictionaryErr := b513SymbolicDictionaryJSON.Marshal(b513SymbolicDictionaryKey)
	if b513SymbolicDictionaryErr != nil {
		return nil, false, b513SymbolicDictionaryErr
	}
	b513SymbolicDictionaryRaw, b513SymbolicDictionaryFound := b513SymbolicDictionaryCached.b513SymbolicDictionaryIndex[string(b513SymbolicDictionaryEncoded)]
	return b513SymbolicDictionaryBytes.Clone(b513SymbolicDictionaryRaw), b513SymbolicDictionaryFound, nil
}

func b513SymbolicDictionaryValidateObject(b513SymbolicDictionaryRaw []byte) error {
	b513SymbolicDictionaryDecoder := b513SymbolicDictionaryJSON.NewDecoder(b513SymbolicDictionaryBytes.NewReader(b513SymbolicDictionaryRaw))
	b513SymbolicDictionaryToken, b513SymbolicDictionaryErr := b513SymbolicDictionaryDecoder.Token()
	if b513SymbolicDictionaryErr != nil || b513SymbolicDictionaryToken != b513SymbolicDictionaryJSON.Delim('{') {
		return b513SymbolicDictionaryFmt.Errorf("ab38 row: expected object")
	}
	b513SymbolicDictionarySeen := map[string]bool{}
	for b513SymbolicDictionaryDecoder.More() {
		b513SymbolicDictionaryToken, b513SymbolicDictionaryErr = b513SymbolicDictionaryDecoder.Token()
		if b513SymbolicDictionaryErr != nil {
			return b513SymbolicDictionaryErr
		}
		b513SymbolicDictionaryName, b513SymbolicDictionaryOK := b513SymbolicDictionaryToken.(string)
		if !b513SymbolicDictionaryOK || b513SymbolicDictionarySeen[b513SymbolicDictionaryName] {
			return b513SymbolicDictionaryFmt.Errorf("ab38 row: duplicate member")
		}
		b513SymbolicDictionarySeen[b513SymbolicDictionaryName] = true
		var b513SymbolicDictionaryValue b513SymbolicDictionaryJSON.RawMessage
		if b513SymbolicDictionaryErr = b513SymbolicDictionaryDecoder.Decode(&b513SymbolicDictionaryValue); b513SymbolicDictionaryErr != nil {
			return b513SymbolicDictionaryErr
		}
		if b513SymbolicDictionaryBytes.Equal(b513SymbolicDictionaryValue, []byte("null")) && b513SymbolicDictionaryName != "secondary" && b513SymbolicDictionaryName != "operation" && b513SymbolicDictionaryName != "group" && b513SymbolicDictionaryName != "register" && b513SymbolicDictionaryName != "value_codec" && b513SymbolicDictionaryName != "value_width_bits" {
			return b513SymbolicDictionaryFmt.Errorf("ab38 row: null required member")
		}
	}
	for _, b513SymbolicDictionaryName := range []string{"family", "variant", "primary", "direction", "field", "alternative", "instance_scope", "dictionary_namespace", "dictionary_class", "raw_address", "native_label", "raw_form", "bus_families", "grade", "physical_qualification", "value_codec", "executable_capability", "address_width_bits", "value_width_bits", "unknown", "unknown_address", "label_meaning", "native_contract"} {
		if !b513SymbolicDictionarySeen[b513SymbolicDictionaryName] {
			return b513SymbolicDictionaryFmt.Errorf("ab38 row: missing member %s", b513SymbolicDictionaryName)
		}
	}
	return nil
}
