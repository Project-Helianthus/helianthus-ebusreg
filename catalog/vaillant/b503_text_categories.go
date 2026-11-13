package vaillant

import (
	b503TextCategoryBytesPackage "bytes"
	_ "embed"
	b503TextCategoryJSON "encoding/json"
	b503TextCategoryFmt "fmt"
	b503TextCategoryIO "io"
	b503TextCategoryReflect "reflect"
)

//go:embed data/ab96_b503_text_categories.json
var b503TextCategoryData []byte

type B503TextCategoryKey struct {
	B503TextCategoryFamily            string  `json:"family"`
	B503TextCategoryVariant           string  `json:"variant"`
	B503TextCategoryPrimary           uint8   `json:"primary"`
	B503TextCategorySecondary         *uint8  `json:"secondary"`
	B503TextCategoryOperation         *uint8  `json:"operation"`
	B503TextCategoryDirection         *string `json:"direction"`
	B503TextCategoryField             string  `json:"field"`
	B503TextCategoryGroup             *uint8  `json:"group"`
	B503TextCategoryRequestedInstance *uint8  `json:"requested_instance"`
	B503TextCategoryRegister          *uint16 `json:"register"`
	B503TextCategoryCommandRole       *string `json:"command_role"`
}

type B503TextCategoryEntry struct {
	B503TextCategoryKey                   B503TextCategoryKey `json:"key"`
	B503TextCategoryGrade                 string              `json:"grade"`
	B503TextCategoryPosition              *int                `json:"position"`
	B503TextCategoryWidthCharacters       *int                `json:"width_characters"`
	B503TextCategoryWidthBits             *int                `json:"width_bits"`
	B503TextCategoryUnsigned              *bool               `json:"unsigned"`
	B503TextCategoryDefault               *uint8              `json:"default"`
	B503TextCategoryEncoding              *string             `json:"encoding"`
	B503TextCategoryPadding               *string             `json:"padding"`
	B503TextCategoryExtent                *string             `json:"extent"`
	B503TextCategoryFollows               *string             `json:"follows"`
	B503TextCategoryLabels                map[uint8]*string   `json:"labels"`
	B503TextCategoryMetadataPresence      string              `json:"metadata_presence"`
	B503TextCategoryPhysicalQualification string              `json:"physical_qualification"`
	B503TextCategoryRawPreservation       string              `json:"raw_preservation"`
	B503TextCategoryUnknown               []string            `json:"unknown"`
	// B503TextCategoryNativeContract is empty until the companion contract is published.
	B503TextCategoryNativeContract string `json:"native_contract"`
}

func B503TextCategoryEntries() ([]B503TextCategoryEntry, error) {
	var b503TextCategoryDecoded []struct {
		B503TextCategoryEntry
		B503TextCategoryKey struct {
			B503TextCategoryKey
			B503TextCategoryPrimary *uint8 `json:"primary"`
		} `json:"key"`
	}
	b503TextCategoryDecoder := b503TextCategoryJSON.NewDecoder(b503TextCategoryBytesPackage.NewReader(b503TextCategoryData))
	b503TextCategoryDecoder.DisallowUnknownFields()
	if b503TextCategoryError := b503TextCategoryDecoder.Decode(&b503TextCategoryDecoded); b503TextCategoryError != nil {
		return nil, b503TextCategoryError
	}
	if b503TextCategoryError := b503TextCategoryDecoder.Decode(new(any)); b503TextCategoryError != b503TextCategoryIO.EOF {
		return nil, b503TextCategoryFmt.Errorf("ab96: trailing embedded data: %v", b503TextCategoryError)
	}
	if len(b503TextCategoryDecoded) == 0 {
		return nil, b503TextCategoryFmt.Errorf("ab96: empty embedded catalog")
	}
	b503TextCategoryEntries := make([]B503TextCategoryEntry, 0, len(b503TextCategoryDecoded))
	for b503TextCategoryIndex, b503TextCategoryWire := range b503TextCategoryDecoded {
		if b503TextCategoryWire.B503TextCategoryKey.B503TextCategoryPrimary == nil {
			return nil, b503TextCategoryFmt.Errorf("ab96: missing embedded primary in entry %d", b503TextCategoryIndex)
		}
		b503TextCategoryEntry := b503TextCategoryWire.B503TextCategoryEntry
		b503TextCategoryEntry.B503TextCategoryKey = b503TextCategoryWire.B503TextCategoryKey.B503TextCategoryKey
		b503TextCategoryEntry.B503TextCategoryKey.B503TextCategoryPrimary = *b503TextCategoryWire.B503TextCategoryKey.B503TextCategoryPrimary
		if b503TextCategoryEntry.B503TextCategoryKey.B503TextCategoryFamily == "" || b503TextCategoryEntry.B503TextCategoryKey.B503TextCategoryVariant == "" || b503TextCategoryEntry.B503TextCategoryKey.B503TextCategoryField == "" ||
			(b503TextCategoryEntry.B503TextCategoryGrade != "Hypothesis" && b503TextCategoryEntry.B503TextCategoryGrade != "Proven" && b503TextCategoryEntry.B503TextCategoryGrade != "Unknown") ||
			b503TextCategoryEntry.B503TextCategoryMetadataPresence != "Proven" || b503TextCategoryEntry.B503TextCategoryPhysicalQualification != "Unknown" ||
			b503TextCategoryEntry.B503TextCategoryRawPreservation == "" || len(b503TextCategoryEntry.B503TextCategoryUnknown) == 0 {
			return nil, b503TextCategoryFmt.Errorf("ab96: invalid embedded entry %d", b503TextCategoryIndex)
		}
		if (b503TextCategoryEntry.B503TextCategoryWidthCharacters != nil && *b503TextCategoryEntry.B503TextCategoryWidthCharacters <= 0) ||
			(b503TextCategoryEntry.B503TextCategoryWidthBits != nil && *b503TextCategoryEntry.B503TextCategoryWidthBits <= 0) ||
			(b503TextCategoryEntry.B503TextCategoryPosition != nil && *b503TextCategoryEntry.B503TextCategoryPosition < 0) {
			return nil, b503TextCategoryFmt.Errorf("ab96: invalid embedded width or position")
		}
		for _, b503TextCategoryPrevious := range b503TextCategoryEntries {
			if b503TextCategoryReflect.DeepEqual(b503TextCategoryPrevious.B503TextCategoryKey, b503TextCategoryEntry.B503TextCategoryKey) {
				return nil, b503TextCategoryFmt.Errorf("ab96: duplicate embedded key")
			}
		}
		b503TextCategoryEntries = append(b503TextCategoryEntries, b503TextCategoryEntry)
	}
	return b503TextCategoryEntries, nil
}

func B503TextCategoryLookup(b503TextCategoryKey B503TextCategoryKey) (B503TextCategoryEntry, bool, error) {
	b503TextCategoryEntries, b503TextCategoryError := B503TextCategoryEntries()
	if b503TextCategoryError != nil {
		return B503TextCategoryEntry{}, false, b503TextCategoryError
	}
	for _, b503TextCategoryEntry := range b503TextCategoryEntries {
		if b503TextCategoryReflect.DeepEqual(b503TextCategoryKey, b503TextCategoryEntry.B503TextCategoryKey) {
			return b503TextCategoryEntry, true, nil
		}
	}
	return B503TextCategoryEntry{}, false, nil
}
