package vaillant

import (
	b509HeatpumpCodecDictionaryBytes "bytes"
	_ "embed"
	b509HeatpumpCodecDictionaryJSONPackage "encoding/json"
	b509HeatpumpCodecDictionaryFmt "fmt"
	b509HeatpumpCodecDictionaryIO "io"
	b509HeatpumpCodecDictionaryReflect "reflect"
)

type B509HeatpumpCodecDictionaryKey struct {
	B509HeatpumpCodecDictionaryFamily    string   `json:"family"`
	B509HeatpumpCodecDictionaryVariant   string   `json:"variant"`
	B509HeatpumpCodecDictionaryPrimary   *byte    `json:"primary"`
	B509HeatpumpCodecDictionarySecondary *byte    `json:"secondary"`
	B509HeatpumpCodecDictionaryOperation *byte    `json:"operation"`
	B509HeatpumpCodecDictionaryDirection string   `json:"direction"`
	B509HeatpumpCodecDictionarySelector  []uint16 `json:"selector"`
}

type B509HeatpumpCodecDictionaryEntry struct {
	B509HeatpumpCodecDictionaryKey              B509HeatpumpCodecDictionaryKey `json:"key"`
	B509HeatpumpCodecDictionaryQuantity         string                         `json:"quantity"`
	B509HeatpumpCodecDeclaration                string                         `json:"codec_declaration"`
	B509HeatpumpCodecDictionaryUnit             *string                        `json:"unit"`
	B509HeatpumpCodecDictionaryWidthBits        *int                           `json:"width_bits"`
	B509HeatpumpCodecDictionarySignedness       *string                        `json:"signedness"`
	B509HeatpumpCodecDictionaryByteOrder        *string                        `json:"byte_order"`
	B509HeatpumpCodecDictionaryDeclaredScale    *int                           `json:"declared_scale"`
	B509HeatpumpCodecDictionaryConverter        *string                        `json:"converter"`
	B509HeatpumpCodecDictionaryMapping          map[string]string              `json:"mapping"`
	B509HeatpumpCodecDictionarySentinel         *int                           `json:"sentinel"`
	B509HeatpumpCodecDictionaryGrade            string                         `json:"grade"`
	B509HeatpumpCodecDictionarySoftwareContract string                         `json:"software_contract"`
	B509HeatpumpCodecDictionaryQualification    string                         `json:"qualification"`
	B509HeatpumpCodecDictionaryPresence         string                         `json:"presence"`
	B509HeatpumpCodecDictionaryFreshness        string                         `json:"freshness"`
	B509HeatpumpCodecDictionaryExecution        string                         `json:"execution"`
	B509HeatpumpCodecDictionaryPhysicalEffects  string                         `json:"physical_effects"`
	B509HeatpumpCodecDictionaryPersistence      string                         `json:"persistence"`
	B509HeatpumpCodecDictionaryDelivery         string                         `json:"delivery"`
	B509HeatpumpCodecDictionaryNotes            []string                       `json:"notes"`
}

//go:embed data/ab70_b509_heatpump_codec_dictionary.json
var b509HeatpumpCodecDictionaryJSON []byte

func b509HeatpumpCodecDictionaryParse(b509HeatpumpCodecDictionaryData []byte) ([]B509HeatpumpCodecDictionaryEntry, error) {
	var b509HeatpumpCodecDictionaryEntries []B509HeatpumpCodecDictionaryEntry
	b509HeatpumpCodecDictionaryDecoder := b509HeatpumpCodecDictionaryJSONPackage.NewDecoder(b509HeatpumpCodecDictionaryBytes.NewReader(b509HeatpumpCodecDictionaryData))
	b509HeatpumpCodecDictionaryDecoder.DisallowUnknownFields()
	if b509HeatpumpCodecDictionaryError := b509HeatpumpCodecDictionaryDecoder.Decode(&b509HeatpumpCodecDictionaryEntries); b509HeatpumpCodecDictionaryError != nil {
		return nil, b509HeatpumpCodecDictionaryError
	}
	if b509HeatpumpCodecDictionaryError := b509HeatpumpCodecDictionaryDecoder.Decode(new(any)); b509HeatpumpCodecDictionaryError != b509HeatpumpCodecDictionaryIO.EOF {
		return nil, b509HeatpumpCodecDictionaryFmt.Errorf("ab70: trailing data")
	}
	if len(b509HeatpumpCodecDictionaryEntries) == 0 {
		return nil, b509HeatpumpCodecDictionaryFmt.Errorf("ab70: empty dictionary")
	}
	b509HeatpumpCodecDictionarySeen := make(map[string]bool)
	for _, b509HeatpumpCodecDictionaryEntry := range b509HeatpumpCodecDictionaryEntries {
		b509HeatpumpCodecDictionaryKey := b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryKey
		if b509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionaryFamily == "" || b509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionaryVariant == "" || b509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionaryPrimary == nil || b509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionarySecondary == nil || b509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionaryOperation == nil || b509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionaryDirection == "" || len(b509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionarySelector) != 2 || b509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionarySelector[0] > 255 || b509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionarySelector[1] > 255 || b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryQuantity == "" || b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDeclaration == "" || b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryGrade != "Hypothesis" || b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionarySoftwareContract != "Proven" || b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryQualification != "Unknown" || b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryPresence != "Unknown" || b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryFreshness != "Unknown" || b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryExecution != "Unknown" || b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryPhysicalEffects != "Unknown" || b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryPersistence != "Unknown" || b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryDelivery != "Unknown" {
			return nil, b509HeatpumpCodecDictionaryFmt.Errorf("ab70: incomplete or invalid entry")
		}
		b509HeatpumpCodecDictionaryEncoded, b509HeatpumpCodecDictionaryError := b509HeatpumpCodecDictionaryJSONPackage.Marshal(b509HeatpumpCodecDictionaryKey)
		if b509HeatpumpCodecDictionaryError != nil {
			return nil, b509HeatpumpCodecDictionaryError
		}
		if b509HeatpumpCodecDictionarySeen[string(b509HeatpumpCodecDictionaryEncoded)] {
			return nil, b509HeatpumpCodecDictionaryFmt.Errorf("ab70: duplicate key")
		}
		b509HeatpumpCodecDictionarySeen[string(b509HeatpumpCodecDictionaryEncoded)] = true
	}
	return b509HeatpumpCodecDictionaryEntries, nil
}

func B509HeatpumpCodecDictionaryEntries() ([]B509HeatpumpCodecDictionaryEntry, error) {
	return b509HeatpumpCodecDictionaryParse(b509HeatpumpCodecDictionaryJSON)
}

func b509HeatpumpCodecDictionaryFind(b509HeatpumpCodecDictionaryEntries []B509HeatpumpCodecDictionaryEntry, b509HeatpumpCodecDictionaryKey B509HeatpumpCodecDictionaryKey) (B509HeatpumpCodecDictionaryEntry, bool, error) {
	for _, b509HeatpumpCodecDictionaryEntry := range b509HeatpumpCodecDictionaryEntries {
		if b509HeatpumpCodecDictionaryReflect.DeepEqual(b509HeatpumpCodecDictionaryEntry.B509HeatpumpCodecDictionaryKey, b509HeatpumpCodecDictionaryKey) {
			b509HeatpumpCodecDictionaryData, b509HeatpumpCodecDictionaryError := b509HeatpumpCodecDictionaryJSONPackage.Marshal(b509HeatpumpCodecDictionaryEntry)
			if b509HeatpumpCodecDictionaryError != nil {
				return B509HeatpumpCodecDictionaryEntry{}, false, b509HeatpumpCodecDictionaryError
			}
			var b509HeatpumpCodecDictionaryCopy B509HeatpumpCodecDictionaryEntry
			if b509HeatpumpCodecDictionaryError = b509HeatpumpCodecDictionaryJSONPackage.Unmarshal(b509HeatpumpCodecDictionaryData, &b509HeatpumpCodecDictionaryCopy); b509HeatpumpCodecDictionaryError != nil {
				return B509HeatpumpCodecDictionaryEntry{}, false, b509HeatpumpCodecDictionaryError
			}
			return b509HeatpumpCodecDictionaryCopy, true, nil
		}
	}
	return B509HeatpumpCodecDictionaryEntry{}, false, nil
}

func B509HeatpumpCodecDictionaryLookup(b509HeatpumpCodecDictionaryKey B509HeatpumpCodecDictionaryKey) (B509HeatpumpCodecDictionaryEntry, bool, error) {
	b509HeatpumpCodecDictionaryEntries, b509HeatpumpCodecDictionaryError := B509HeatpumpCodecDictionaryEntries()
	if b509HeatpumpCodecDictionaryError != nil {
		return B509HeatpumpCodecDictionaryEntry{}, false, b509HeatpumpCodecDictionaryError
	}
	return b509HeatpumpCodecDictionaryFind(b509HeatpumpCodecDictionaryEntries, b509HeatpumpCodecDictionaryKey)
}
