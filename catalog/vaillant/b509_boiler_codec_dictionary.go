package vaillant

import (
	b509BoilerCodecDictionaryBytes "bytes"
	_ "embed"
	b509BoilerCodecDictionaryJSONPackage "encoding/json"
	b509BoilerCodecDictionaryFmt "fmt"
	b509BoilerCodecDictionaryIO "io"
	b509BoilerCodecDictionaryReflect "reflect"
)

type B509BoilerCodecDictionaryKey struct {
	B509BoilerCodecDictionaryFamily    string   `json:"family"`
	B509BoilerCodecDictionaryVariant   string   `json:"variant"`
	B509BoilerCodecDictionaryPrimary   *byte    `json:"primary"`
	B509BoilerCodecDictionarySecondary *byte    `json:"secondary"`
	B509BoilerCodecDictionaryOperation *byte    `json:"operation"`
	B509BoilerCodecDictionaryDirection string   `json:"direction"`
	B509BoilerCodecDictionarySelector  []uint16 `json:"selector"`
}

type B509BoilerCodecDictionaryEntry struct {
	B509BoilerCodecDictionaryKey       B509BoilerCodecDictionaryKey `json:"key"`
	B509BoilerCodecDictionaryQuantity  string                       `json:"quantity"`
	B509BoilerCodecDeclaration         string                       `json:"codec_declaration"`
	B509BoilerCodecDictionaryWidthBits *int                         `json:"width_bits"`
	B509BoilerCodecDictionaryByteOrder *string                      `json:"byte_order"`
	B509BoilerCodecDictionaryUnit      *string                      `json:"unit"`
	B509BoilerCodecDictionaryConverter *string                      `json:"converter"`
	B509BoilerCodecDictionaryMapping   map[string]string            `json:"mapping"`
	B509BoilerCodecDictionarySentinel  *int64                       `json:"sentinel"`
	B509BoilerCodecDictionaryGrade     string                       `json:"grade"`
	B509BoilerCodecDictionaryConflicts []string                     `json:"conflicts"`
	B509BoilerCodecDictionaryEvidence  map[string]string            `json:"evidence"`
}

//go:embed data/ab68_b509_boiler_codec_dictionary.json
var b509BoilerCodecDictionaryJSON []byte

func b509BoilerCodecDictionaryParse(b509BoilerCodecDictionaryData []byte) ([]B509BoilerCodecDictionaryEntry, error) {
	var b509BoilerCodecDictionaryEntries []B509BoilerCodecDictionaryEntry
	b509BoilerCodecDictionaryDecoder := b509BoilerCodecDictionaryJSONPackage.NewDecoder(b509BoilerCodecDictionaryBytes.NewReader(b509BoilerCodecDictionaryData))
	b509BoilerCodecDictionaryDecoder.DisallowUnknownFields()
	if b509BoilerCodecDictionaryError := b509BoilerCodecDictionaryDecoder.Decode(&b509BoilerCodecDictionaryEntries); b509BoilerCodecDictionaryError != nil {
		return nil, b509BoilerCodecDictionaryError
	}
	if b509BoilerCodecDictionaryError := b509BoilerCodecDictionaryDecoder.Decode(new(any)); b509BoilerCodecDictionaryError != b509BoilerCodecDictionaryIO.EOF {
		return nil, b509BoilerCodecDictionaryFmt.Errorf("ab68: trailing JSON: %v", b509BoilerCodecDictionaryError)
	}
	if len(b509BoilerCodecDictionaryEntries) == 0 {
		return nil, b509BoilerCodecDictionaryFmt.Errorf("ab68: empty dictionary")
	}
	b509BoilerCodecDictionarySeen := make(map[string]bool)
	for b509BoilerCodecDictionaryIndex, b509BoilerCodecDictionaryEntry := range b509BoilerCodecDictionaryEntries {
		b509BoilerCodecDictionaryKey := b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryKey
		if b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryFamily == "" || b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryVariant == "" || b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryPrimary == nil || b509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySecondary == nil || b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryOperation == nil || b509BoilerCodecDictionaryKey.B509BoilerCodecDictionaryDirection == "" || len(b509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector) != 2 {
			return nil, b509BoilerCodecDictionaryFmt.Errorf("ab68: incomplete key at row %d", b509BoilerCodecDictionaryIndex)
		}
		for _, b509BoilerCodecDictionaryByte := range b509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector {
			if b509BoilerCodecDictionaryByte > 255 {
				return nil, b509BoilerCodecDictionaryFmt.Errorf("ab68: selector byte out of range at row %d", b509BoilerCodecDictionaryIndex)
			}
		}
		if b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryQuantity == "" || b509BoilerCodecDictionaryEntry.B509BoilerCodecDeclaration == "" || b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryGrade != "Hypothesis" {
			return nil, b509BoilerCodecDictionaryFmt.Errorf("ab68: incomplete declaration at row %d", b509BoilerCodecDictionaryIndex)
		}
		if b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryEvidence["metadata"] != "Proven" {
			return nil, b509BoilerCodecDictionaryFmt.Errorf("ab68: missing metadata grade at row %d", b509BoilerCodecDictionaryIndex)
		}
		for _, b509BoilerCodecDictionaryAxis := range []string{"qualification", "physical_support", "presence", "freshness", "effects", "persistence", "delivery", "execution"} {
			if b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryEvidence[b509BoilerCodecDictionaryAxis] != "Unknown" {
				return nil, b509BoilerCodecDictionaryFmt.Errorf("ab68: invalid %s grade at row %d", b509BoilerCodecDictionaryAxis, b509BoilerCodecDictionaryIndex)
			}
		}
		b509BoilerCodecDictionaryEncoded, b509BoilerCodecDictionaryError := b509BoilerCodecDictionaryJSONPackage.Marshal(b509BoilerCodecDictionaryKey)
		if b509BoilerCodecDictionaryError != nil {
			return nil, b509BoilerCodecDictionaryError
		}
		if b509BoilerCodecDictionarySeen[string(b509BoilerCodecDictionaryEncoded)] {
			return nil, b509BoilerCodecDictionaryFmt.Errorf("ab68: duplicate key at row %d", b509BoilerCodecDictionaryIndex)
		}
		b509BoilerCodecDictionarySeen[string(b509BoilerCodecDictionaryEncoded)] = true
	}
	return b509BoilerCodecDictionaryEntries, nil
}

func B509BoilerCodecDictionaryEntries() ([]B509BoilerCodecDictionaryEntry, error) {
	return b509BoilerCodecDictionaryParse(b509BoilerCodecDictionaryJSON)
}

func b509BoilerCodecDictionaryLookup(b509BoilerCodecDictionaryData []byte, b509BoilerCodecDictionaryKey B509BoilerCodecDictionaryKey) (B509BoilerCodecDictionaryEntry, bool, error) {
	b509BoilerCodecDictionaryEntries, b509BoilerCodecDictionaryError := b509BoilerCodecDictionaryParse(b509BoilerCodecDictionaryData)
	if b509BoilerCodecDictionaryError != nil {
		return B509BoilerCodecDictionaryEntry{}, false, b509BoilerCodecDictionaryError
	}
	for _, b509BoilerCodecDictionaryEntry := range b509BoilerCodecDictionaryEntries {
		if b509BoilerCodecDictionaryReflect.DeepEqual(b509BoilerCodecDictionaryEntry.B509BoilerCodecDictionaryKey, b509BoilerCodecDictionaryKey) {
			return b509BoilerCodecDictionaryEntry, true, nil
		}
	}
	return B509BoilerCodecDictionaryEntry{}, false, nil
}

func B509BoilerCodecDictionaryLookup(b509BoilerCodecDictionaryKey B509BoilerCodecDictionaryKey) (B509BoilerCodecDictionaryEntry, bool, error) {
	return b509BoilerCodecDictionaryLookup(b509BoilerCodecDictionaryJSON, b509BoilerCodecDictionaryKey)
}
