package vaillant

import (
	b51ADiagnosticVariantBytes "bytes"
	_ "embed"
	b51ADiagnosticVariantJSON "encoding/json"
	b51ADiagnosticVariantFmt "fmt"
	b51ADiagnosticVariantIO "io"
	b51ADiagnosticVariantReflect "reflect"
)

//go:embed data/ab98_b51a_diagnostic_variants.json
var b51ADiagnosticVariantData []byte

type B51ADiagnosticVariantKey struct {
	B51ADiagnosticVariantFamily            string  `json:"family"`
	B51ADiagnosticVariant                  string  `json:"variant"`
	B51ADiagnosticVariantPrimary           uint8   `json:"primary"`
	B51ADiagnosticVariantSecondary         *uint8  `json:"secondary"`
	B51ADiagnosticVariantOperation         *uint8  `json:"operation"`
	B51ADiagnosticVariantDirection         string  `json:"direction"`
	B51ADiagnosticVariantField             string  `json:"field"`
	B51ADiagnosticVariantGroup             *uint8  `json:"group"`
	B51ADiagnosticVariantRequestedInstance *uint8  `json:"requested_instance"`
	B51ADiagnosticVariantRegister          *uint16 `json:"register"`
	B51ADiagnosticVariantCommandRole       *string `json:"command_role"`
}

type B51ADiagnosticVariantCondition struct {
	B51ADiagnosticVariantField string `json:"field"`
	B51ADiagnosticVariantBit   uint8  `json:"bit"`
	B51ADiagnosticVariantSet   bool   `json:"set"`
}

type B51ADiagnosticVariantBit struct {
	B51ADiagnosticVariantBit      uint8    `json:"bit"`
	B51ADiagnosticVariantLabels   []string `json:"labels"`
	B51ADiagnosticVariantReserved bool     `json:"reserved"`
	B51ADiagnosticVariantDefault  bool     `json:"default"`
}

type B51ADiagnosticVariantEntry struct {
	B51ADiagnosticVariantKey                   B51ADiagnosticVariantKey           `json:"key"`
	B51ADiagnosticVariantGrade                 string                             `json:"grade"`
	B51ADiagnosticVariantMetadataPresence      string                             `json:"metadata_presence"`
	B51ADiagnosticVariantPhysicalQualification string                             `json:"physical_qualification"`
	B51ADiagnosticVariantInvocation            string                             `json:"invocation"`
	B51ADiagnosticVariantNativeContract        string                             `json:"native_contract"`
	B51ADiagnosticVariantPosition              *uint16                            `json:"position"`
	B51ADiagnosticVariantWidthBits             *uint16                            `json:"width_bits"`
	B51ADiagnosticVariantByteOrder             *string                            `json:"byte_order"`
	B51ADiagnosticVariantDomain                []uint64                           `json:"domain"`
	B51ADiagnosticVariantDefault               *uint64                            `json:"default"`
	B51ADiagnosticVariantValues                map[uint64]string                  `json:"values"`
	B51ADiagnosticVariantBits                  []B51ADiagnosticVariantBit         `json:"bits"`
	B51ADiagnosticVariantConditions            []B51ADiagnosticVariantCondition   `json:"conditions"`
	B51ADiagnosticVariantAlternatives          []B51ADiagnosticVariantAlternative `json:"alternatives"`
	B51ADiagnosticVariantInterpretations       []string                           `json:"interpretations"`
	B51ADiagnosticVariantUnknown               []string                           `json:"unknown"`
	B51ADiagnosticVariantExclusive             *bool                              `json:"exclusive"`
	B51ADiagnosticVariantCount                 *uint16                            `json:"count"`
	B51ADiagnosticVariantLogicalFootprint      *uint16                            `json:"logical_footprint"`
	B51ADiagnosticVariantTail                  []uint16                           `json:"tail"`
}

type B51ADiagnosticVariantAlternative struct {
	B51ADiagnosticVariantLabel           string                           `json:"label"`
	B51ADiagnosticVariantConditions      []B51ADiagnosticVariantCondition `json:"conditions"`
	B51ADiagnosticVariantFields          []B51ADiagnosticVariantEntry     `json:"fields"`
	B51ADiagnosticVariantInterpretations []string                         `json:"interpretations"`
}

func B51ADiagnosticVariantEntries() ([]B51ADiagnosticVariantEntry, error) {
	var b51ADiagnosticVariantEntries []B51ADiagnosticVariantEntry
	b51ADiagnosticVariantDecoder := b51ADiagnosticVariantJSON.NewDecoder(b51ADiagnosticVariantBytes.NewReader(b51ADiagnosticVariantData))
	b51ADiagnosticVariantDecoder.DisallowUnknownFields()
	if b51ADiagnosticVariantErr := b51ADiagnosticVariantDecoder.Decode(&b51ADiagnosticVariantEntries); b51ADiagnosticVariantErr != nil {
		return nil, b51ADiagnosticVariantFmt.Errorf("ab98 catalog: %w", b51ADiagnosticVariantErr)
	}
	if b51ADiagnosticVariantErr := b51ADiagnosticVariantDecoder.Decode(new(any)); b51ADiagnosticVariantErr != b51ADiagnosticVariantIO.EOF {
		return nil, b51ADiagnosticVariantFmt.Errorf("ab98 catalog: trailing data")
	}
	if len(b51ADiagnosticVariantEntries) == 0 {
		return nil, b51ADiagnosticVariantFmt.Errorf("ab98 catalog: empty entries")
	}
	for b51ADiagnosticVariantIndex, b51ADiagnosticVariantEntry := range b51ADiagnosticVariantEntries {
		if b51ADiagnosticVariantErr := b51ADiagnosticVariantValidate(b51ADiagnosticVariantEntry); b51ADiagnosticVariantErr != nil {
			return nil, b51ADiagnosticVariantErr
		}
		for b51ADiagnosticVariantPrevious := 0; b51ADiagnosticVariantPrevious < b51ADiagnosticVariantIndex; b51ADiagnosticVariantPrevious++ {
			if b51ADiagnosticVariantReflect.DeepEqual(b51ADiagnosticVariantEntries[b51ADiagnosticVariantPrevious].B51ADiagnosticVariantKey, b51ADiagnosticVariantEntry.B51ADiagnosticVariantKey) {
				return nil, b51ADiagnosticVariantFmt.Errorf("ab98 catalog: duplicate key")
			}
		}
	}
	return b51ADiagnosticVariantEntries, nil
}

func b51ADiagnosticVariantValidate(b51ADiagnosticVariantEntry B51ADiagnosticVariantEntry) error {
	if b51ADiagnosticVariantEntry.B51ADiagnosticVariantKey.B51ADiagnosticVariantFamily == "" || b51ADiagnosticVariantEntry.B51ADiagnosticVariantKey.B51ADiagnosticVariant == "" || b51ADiagnosticVariantEntry.B51ADiagnosticVariantKey.B51ADiagnosticVariantField == "" || (b51ADiagnosticVariantEntry.B51ADiagnosticVariantKey.B51ADiagnosticVariantDirection != "request" && b51ADiagnosticVariantEntry.B51ADiagnosticVariantKey.B51ADiagnosticVariantDirection != "positive") {
		return b51ADiagnosticVariantFmt.Errorf("ab98 catalog: incomplete key")
	}
	if b51ADiagnosticVariantEntry.B51ADiagnosticVariantGrade != "Hypothesis" || b51ADiagnosticVariantEntry.B51ADiagnosticVariantMetadataPresence != "Proven" || b51ADiagnosticVariantEntry.B51ADiagnosticVariantPhysicalQualification != "Unknown" || b51ADiagnosticVariantEntry.B51ADiagnosticVariantInvocation != "Unknown" || b51ADiagnosticVariantEntry.B51ADiagnosticVariantNativeContract == "" {
		return b51ADiagnosticVariantFmt.Errorf("ab98 catalog: invalid evidence boundary")
	}
	if len(b51ADiagnosticVariantEntry.B51ADiagnosticVariantDomain) != 0 && (len(b51ADiagnosticVariantEntry.B51ADiagnosticVariantDomain) != 2 || b51ADiagnosticVariantEntry.B51ADiagnosticVariantDomain[0] > b51ADiagnosticVariantEntry.B51ADiagnosticVariantDomain[1]) {
		return b51ADiagnosticVariantFmt.Errorf("ab98 catalog: invalid domain")
	}
	if b51ADiagnosticVariantEntry.B51ADiagnosticVariantDefault != nil && len(b51ADiagnosticVariantEntry.B51ADiagnosticVariantDomain) == 2 && (*b51ADiagnosticVariantEntry.B51ADiagnosticVariantDefault < b51ADiagnosticVariantEntry.B51ADiagnosticVariantDomain[0] || *b51ADiagnosticVariantEntry.B51ADiagnosticVariantDefault > b51ADiagnosticVariantEntry.B51ADiagnosticVariantDomain[1]) {
		return b51ADiagnosticVariantFmt.Errorf("ab98 catalog: default outside domain")
	}
	if b51ADiagnosticVariantEntry.B51ADiagnosticVariantWidthBits != nil && (*b51ADiagnosticVariantEntry.B51ADiagnosticVariantWidthBits == 0 || *b51ADiagnosticVariantEntry.B51ADiagnosticVariantWidthBits > 64) {
		return b51ADiagnosticVariantFmt.Errorf("ab98 catalog: invalid width")
	}
	if b51ADiagnosticVariantEntry.B51ADiagnosticVariantByteOrder != nil && *b51ADiagnosticVariantEntry.B51ADiagnosticVariantByteOrder != "little-endian" {
		return b51ADiagnosticVariantFmt.Errorf("ab98 catalog: invalid byte order")
	}
	b51ADiagnosticVariantSeenBits := map[uint8]bool{}
	for _, b51ADiagnosticVariantBit := range b51ADiagnosticVariantEntry.B51ADiagnosticVariantBits {
		if b51ADiagnosticVariantBit.B51ADiagnosticVariantBit > 7 || b51ADiagnosticVariantSeenBits[b51ADiagnosticVariantBit.B51ADiagnosticVariantBit] {
			return b51ADiagnosticVariantFmt.Errorf("ab98 catalog: invalid bit")
		}
		b51ADiagnosticVariantSeenBits[b51ADiagnosticVariantBit.B51ADiagnosticVariantBit] = true
	}
	if b51ADiagnosticVariantErr := b51ADiagnosticVariantValidateConditions(b51ADiagnosticVariantEntry.B51ADiagnosticVariantConditions); b51ADiagnosticVariantErr != nil {
		return b51ADiagnosticVariantErr
	}
	for _, b51ADiagnosticVariantAlternative := range b51ADiagnosticVariantEntry.B51ADiagnosticVariantAlternatives {
		if b51ADiagnosticVariantErr := b51ADiagnosticVariantValidateConditions(b51ADiagnosticVariantAlternative.B51ADiagnosticVariantConditions); b51ADiagnosticVariantErr != nil {
			return b51ADiagnosticVariantErr
		}
	}
	for _, b51ADiagnosticVariantAlternative := range b51ADiagnosticVariantEntry.B51ADiagnosticVariantAlternatives {
		if b51ADiagnosticVariantAlternative.B51ADiagnosticVariantLabel == "" || len(b51ADiagnosticVariantAlternative.B51ADiagnosticVariantFields) == 0 {
			return b51ADiagnosticVariantFmt.Errorf("ab98 catalog: invalid alternative")
		}
		for _, b51ADiagnosticVariantField := range b51ADiagnosticVariantAlternative.B51ADiagnosticVariantFields {
			if b51ADiagnosticVariantErr := b51ADiagnosticVariantValidate(b51ADiagnosticVariantField); b51ADiagnosticVariantErr != nil {
				return b51ADiagnosticVariantErr
			}
		}
	}
	return nil
}

func b51ADiagnosticVariantValidateConditions(b51ADiagnosticVariantConditions []B51ADiagnosticVariantCondition) error {
	for _, b51ADiagnosticVariantCondition := range b51ADiagnosticVariantConditions {
		if b51ADiagnosticVariantCondition.B51ADiagnosticVariantBit > 7 || b51ADiagnosticVariantCondition.B51ADiagnosticVariantField == "" {
			return b51ADiagnosticVariantFmt.Errorf("ab98 catalog: invalid condition")
		}
	}
	return nil
}

func B51ADiagnosticVariantLookup(b51ADiagnosticVariantKey B51ADiagnosticVariantKey) (B51ADiagnosticVariantEntry, bool, error) {
	b51ADiagnosticVariantEntries, b51ADiagnosticVariantErr := B51ADiagnosticVariantEntries()
	if b51ADiagnosticVariantErr != nil {
		return B51ADiagnosticVariantEntry{}, false, b51ADiagnosticVariantErr
	}
	for _, b51ADiagnosticVariantEntry := range b51ADiagnosticVariantEntries {
		if b51ADiagnosticVariantReflect.DeepEqual(b51ADiagnosticVariantKey, b51ADiagnosticVariantEntry.B51ADiagnosticVariantKey) {
			return b51ADiagnosticVariantEntry, true, nil
		}
	}
	return B51ADiagnosticVariantEntry{}, false, nil
}
