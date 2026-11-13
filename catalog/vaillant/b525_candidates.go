package vaillant

import (
	_ "embed"
	b525CandidateJSON "encoding/json"
)

//go:embed data/ab73_b525_candidates.json
var b525CandidateB525Data []byte

type B525CandidateKey struct {
	B525CandidateFamily    string  `json:"family"`
	B525CandidateVariant   *string `json:"variant"`
	B525CandidatePrimary   byte    `json:"primary"`
	B525CandidateSecondary *byte   `json:"secondary"`
	B525CandidateOperation *byte   `json:"operation"`
	B525CandidateDirection string  `json:"direction"`
	B525CandidateField     string  `json:"field"`
}

type B525CandidateEntry struct {
	B525CandidateKey        B525CandidateKey `json:"key"`
	B525CandidateGrade      string           `json:"grade"`
	B525CandidatePosition   *int             `json:"position"`
	B525CandidateWidth      *int             `json:"width"`
	B525CandidateCodec      string           `json:"codec"`
	B525CandidateMinimum    *int             `json:"minimum"`
	B525CandidateMaximum    *int             `json:"maximum"`
	B525CandidateDefault    *int             `json:"default"`
	B525CandidateCount      *int             `json:"count"`
	B525CandidateStartIndex *int             `json:"start_index"`
	B525CandidateUnit       *string          `json:"unit"`
	B525CandidateMappings   map[int]string   `json:"mappings"`
	B525CandidateBits       map[int]string   `json:"bits"`
	B525CandidateUnknowns   []string         `json:"unknowns"`
	B525CandidateRawPolicy  string           `json:"raw_policy"`
	B525CandidateRelation   string           `json:"relation"`
}

type B525CandidateCatalog struct {
	B525CandidateRepresentationGrade   string               `json:"representation_grade"`
	B525CandidatePhysicalQualification string               `json:"physical_qualification"`
	B525CandidateCoordinateSystem      string               `json:"coordinate_system"`
	B525CandidateLimitations           []string             `json:"limitations"`
	B525CandidateEntries               []B525CandidateEntry `json:"entries"`
}

func B525Candidates() B525CandidateCatalog {
	var b525CandidateCatalog B525CandidateCatalog
	if b525CandidateError := b525CandidateJSON.Unmarshal(b525CandidateB525Data, &b525CandidateCatalog); b525CandidateError != nil {
		panic(b525CandidateError)
	}
	return b525CandidateCatalog
}

func b525CandidateEqualOptional[b525CandidateT comparable](b525CandidateLeft, b525CandidateRight *b525CandidateT) bool {
	return b525CandidateLeft == nil && b525CandidateRight == nil || b525CandidateLeft != nil && b525CandidateRight != nil && *b525CandidateLeft == *b525CandidateRight
}

func B525CandidateLookup(b525CandidateKey B525CandidateKey) (B525CandidateEntry, bool) {
	for _, b525CandidateEntry := range B525Candidates().B525CandidateEntries {
		b525CandidateCandidate := b525CandidateEntry.B525CandidateKey
		if b525CandidateCandidate.B525CandidateFamily == b525CandidateKey.B525CandidateFamily &&
			b525CandidateEqualOptional(b525CandidateCandidate.B525CandidateVariant, b525CandidateKey.B525CandidateVariant) &&
			b525CandidateCandidate.B525CandidatePrimary == b525CandidateKey.B525CandidatePrimary &&
			b525CandidateEqualOptional(b525CandidateCandidate.B525CandidateSecondary, b525CandidateKey.B525CandidateSecondary) &&
			b525CandidateEqualOptional(b525CandidateCandidate.B525CandidateOperation, b525CandidateKey.B525CandidateOperation) &&
			b525CandidateCandidate.B525CandidateDirection == b525CandidateKey.B525CandidateDirection &&
			b525CandidateCandidate.B525CandidateField == b525CandidateKey.B525CandidateField {
			return b525CandidateEntry, true
		}
	}
	return B525CandidateEntry{}, false
}
