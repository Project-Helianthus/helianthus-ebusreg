package vaillant

import (
	_ "embed"
	b512CandidateJSON "encoding/json"
)

//go:embed data/ab69_b512_candidates.json
var b512CandidateB512JSON []byte

type B512CandidateKey struct {
	B512CandidateFamily    string `json:"family"`
	B512CandidateVariant   string `json:"variant"`
	B512CandidatePrimary   byte   `json:"primary"`
	B512CandidateSecondary *byte  `json:"secondary"`
	B512CandidateOperation *byte  `json:"operation"`
	B512CandidateDirection string `json:"direction"`
	B512CandidateField     string `json:"field"`
}

type B512CandidateRange struct {
	B512CandidateMin     int    `json:"min"`
	B512CandidateMax     int    `json:"max"`
	B512CandidateMeaning string `json:"meaning"`
}

type B512Candidate struct {
	B512CandidateKey            B512CandidateKey     `json:"key"`
	B512CandidateGrade          string               `json:"grade"`
	B512CandidatePosition       *int                 `json:"logical_position"`
	B512CandidateWidth          *int                 `json:"width"`
	B512CandidateRawDomain      []int                `json:"raw_domain"`
	B512CandidateSigned         *bool                `json:"signed"`
	B512CandidateByteOrder      *string              `json:"byte_order"`
	B512CandidateDivisor        *float64             `json:"divisor"`
	B512CandidateUnit           *string              `json:"unit"`
	B512CandidateDefault        *int                 `json:"default"`
	B512CandidateMappings       map[string]string    `json:"mappings"`
	B512CandidateRanges         []B512CandidateRange `json:"ranges"`
	B512CandidateAlternatives   []string             `json:"alternatives"`
	B512CandidateTimeoutSeconds *int                 `json:"timeout_seconds"`
	B512CandidateTimeoutRestart *string              `json:"timeout_restart"`
	B512CandidateExternalMode   *int                 `json:"external_mode"`
	B512CandidateNotes          []string             `json:"notes"`
}

type B512CandidateCatalog struct {
	B512CandidateRepresentationGrade   string          `json:"representation_grade"`
	B512CandidatePhysicalQualification string          `json:"physical_qualification"`
	B512CandidateInstalledSupport      string          `json:"installed_support"`
	B512CandidateResponseAdaptation    string          `json:"response_adaptation"`
	B512CandidateActualSelection       string          `json:"actual_versus_fallback_selection"`
	B512CandidateScope                 string          `json:"scope"`
	B512CandidateEntries               []B512Candidate `json:"entries"`
}

func B512Candidates() B512CandidateCatalog {
	var b512CandidateCatalog B512CandidateCatalog
	if b512CandidateErr := b512CandidateJSON.Unmarshal(b512CandidateB512JSON, &b512CandidateCatalog); b512CandidateErr != nil {
		panic(b512CandidateErr)
	}
	return b512CandidateCatalog
}

func b512CandidateOptionalEqual(b512CandidateLeft, b512CandidateRight *byte) bool {
	if b512CandidateLeft == nil || b512CandidateRight == nil {
		return b512CandidateLeft == nil && b512CandidateRight == nil
	}
	return *b512CandidateLeft == *b512CandidateRight
}

func B512CandidateLookup(b512CandidateKey B512CandidateKey) (B512Candidate, bool) {
	for _, b512CandidateEntry := range B512Candidates().B512CandidateEntries {
		b512CandidateStored := b512CandidateEntry.B512CandidateKey
		if b512CandidateStored.B512CandidateFamily == b512CandidateKey.B512CandidateFamily && b512CandidateStored.B512CandidateVariant == b512CandidateKey.B512CandidateVariant && b512CandidateStored.B512CandidatePrimary == b512CandidateKey.B512CandidatePrimary && b512CandidateOptionalEqual(b512CandidateStored.B512CandidateSecondary, b512CandidateKey.B512CandidateSecondary) && b512CandidateOptionalEqual(b512CandidateStored.B512CandidateOperation, b512CandidateKey.B512CandidateOperation) && b512CandidateStored.B512CandidateDirection == b512CandidateKey.B512CandidateDirection && b512CandidateStored.B512CandidateField == b512CandidateKey.B512CandidateField {
			return b512CandidateEntry, true
		}
	}
	return B512Candidate{}, false
}
