package vaillant

import (
	_ "embed"
	b511CandidateJSON "encoding/json"
)

//go:embed data/ab67_b511_candidates.json
var b511CandidateB511Data []byte

type B511CandidateB511Key struct {
	B511CandidateFamily    string `json:"family"`
	B511CandidateVariant   string `json:"variant"`
	B511CandidatePrimary   *uint8 `json:"primary"`
	B511CandidateSecondary *uint8 `json:"secondary"`
	B511CandidateOperation *uint8 `json:"operation"`
	B511CandidateDirection string `json:"direction"`
	B511CandidateField     string `json:"field"`
}

type B511CandidateEntry struct {
	B511CandidateKey              B511CandidateB511Key `json:"key"`
	B511CandidatePosition         *int                 `json:"position"`
	B511CandidateWidthBits        *int                 `json:"width_bits"`
	B511CandidateExtent           *int                 `json:"extent"`
	B511CandidateGrade            SchemaGrade          `json:"grade"`
	B511CandidateDefault          *int                 `json:"default"`
	B511CandidateValues           map[uint8]string     `json:"values"`
	B511CandidateBits             map[uint8]string     `json:"bits"`
	B511CandidateReservedBits     []uint8              `json:"reserved_bits"`
	B511CandidateMembership       []uint8              `json:"membership"`
	B511CandidateRawDomain        []int                `json:"raw_domain"`
	B511CandidateScale            *float64             `json:"scale"`
	B511CandidateUnit             *string              `json:"unit"`
	B511CandidateClockValidation  SchemaGrade          `json:"clock_validation"`
	B511CandidatePhysicalPrefix   SchemaGrade          `json:"physical_prefix"`
	B511CandidateInstalledSupport SchemaGrade          `json:"installed_support"`
	B511CandidateRawPolicy        string               `json:"raw_policy"`
}

type B511CandidateDeclarations struct {
	B511CandidateMetadataPresence      SchemaGrade          `json:"metadata_presence"`
	B511CandidatePhysicalQualification SchemaGrade          `json:"physical_qualification"`
	B511CandidateNativeContract        string               `json:"native_contract"`
	B511CandidateLimitations           []string             `json:"limitations"`
	B511CandidateEntries               []B511CandidateEntry `json:"entries"`
}

func B511CandidateCatalog() (B511CandidateDeclarations, error) {
	var b511CandidateCatalog B511CandidateDeclarations
	b511CandidateError := b511CandidateJSON.Unmarshal(b511CandidateB511Data, &b511CandidateCatalog)
	return b511CandidateCatalog, b511CandidateError
}

func b511CandidateEqualByte(b511CandidateLeft, b511CandidateRight *uint8) bool {
	if b511CandidateLeft == nil || b511CandidateRight == nil {
		return b511CandidateLeft == nil && b511CandidateRight == nil
	}
	return *b511CandidateLeft == *b511CandidateRight
}

func B511CandidateLookup(b511CandidateKey B511CandidateB511Key) (B511CandidateEntry, bool, error) {
	b511CandidateCatalog, b511CandidateError := B511CandidateCatalog()
	if b511CandidateError != nil {
		return B511CandidateEntry{}, false, b511CandidateError
	}
	for _, b511CandidateEntry := range b511CandidateCatalog.B511CandidateEntries {
		b511CandidateCandidate := b511CandidateEntry.B511CandidateKey
		if b511CandidateCandidate.B511CandidateFamily == b511CandidateKey.B511CandidateFamily &&
			b511CandidateCandidate.B511CandidateVariant == b511CandidateKey.B511CandidateVariant &&
			b511CandidateEqualByte(b511CandidateCandidate.B511CandidatePrimary, b511CandidateKey.B511CandidatePrimary) &&
			b511CandidateEqualByte(b511CandidateCandidate.B511CandidateSecondary, b511CandidateKey.B511CandidateSecondary) &&
			b511CandidateEqualByte(b511CandidateCandidate.B511CandidateOperation, b511CandidateKey.B511CandidateOperation) &&
			b511CandidateCandidate.B511CandidateDirection == b511CandidateKey.B511CandidateDirection &&
			b511CandidateCandidate.B511CandidateField == b511CandidateKey.B511CandidateField {
			return b511CandidateEntry, true, nil
		}
	}
	return B511CandidateEntry{}, false, nil
}
