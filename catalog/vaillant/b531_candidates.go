package vaillant

import (
	_ "embed"
	b531CandidateJSON "encoding/json"
)

type B531CandidateKey struct {
	B531CandidateFamily    string  `json:"family"`
	B531CandidateVariant   string  `json:"variant"`
	B531CandidatePrimary   byte    `json:"primary"`
	B531CandidateSecondary *byte   `json:"secondary"`
	B531CandidateOperation *byte   `json:"operation"`
	B531CandidateDirection *string `json:"direction"`
	B531CandidateField     string  `json:"field"`
}

type B531CandidateProjection struct {
	B531CandidateConversion string   `json:"conversion"`
	B531CandidateMultiplier uint32   `json:"multiplier"`
	B531CandidateWidthBits  *int     `json:"width_bits"`
	B531CandidateDomain     []uint32 `json:"domain"`
	B531CandidateUnit       string   `json:"unit"`
	B531CandidateGrade      string   `json:"grade"`
}

type B531CandidateMaximumTable struct {
	B531CandidateSourceField    string `json:"source_field"`
	B531CandidateIndexDomain    []int  `json:"index_domain"`
	B531CandidateInstalledCount *int   `json:"installed_count"`
	B531CandidateGrade          string `json:"grade"`
}

type B531CandidateEntry struct {
	B531CandidateKey             B531CandidateKey           `json:"key"`
	B531CandidatePosition        *int                       `json:"position"`
	B531CandidateWidthBits       *int                       `json:"width_bits"`
	B531CandidateRawDomain       []uint32                   `json:"raw_domain"`
	B531CandidateSignedness      string                     `json:"signedness"`
	B531CandidateByteOrder       string                     `json:"byte_order"`
	B531CandidateNativeUnit      *string                    `json:"native_unit"`
	B531CandidateNativeUnitGrade string                     `json:"native_unit_grade"`
	B531CandidateSentinel        *uint16                    `json:"sentinel"`
	B531CandidateGrade           string                     `json:"grade"`
	B531CandidateLocalProjection B531CandidateProjection    `json:"local_projection"`
	B531CandidateMaximumTable    *B531CandidateMaximumTable `json:"maximum_table"`
}

type B531CandidateCatalog struct {
	B531CandidateMetadataPresence   string               `json:"metadata_presence"`
	B531CandidatePhysicalSupport    string               `json:"physical_support"`
	B531CandidateEndpointSupport    string               `json:"endpoint_support"`
	B531CandidateCorrelation        string               `json:"correlation"`
	B531CandidateReplyAdaptation    string               `json:"reply_adaptation"`
	B531CandidateCacheFreshness     string               `json:"cache_freshness"`
	B531CandidateTargetEligibility  string               `json:"target_eligibility"`
	B531CandidateAppliedLimits      string               `json:"applied_limits"`
	B531CandidateOrderingConstraint *string              `json:"ordering_constraint"`
	B531CandidateLimitations        []string             `json:"limitations"`
	B531CandidateEntries            []B531CandidateEntry `json:"entries"`
}

//go:embed data/ab75_b531_candidates.json
var b531CandidateB531JSON []byte

func B531Candidates() (B531CandidateCatalog, error) {
	var b531CandidateCatalog B531CandidateCatalog
	b531CandidateError := b531CandidateJSON.Unmarshal(b531CandidateB531JSON, &b531CandidateCatalog)
	return b531CandidateCatalog, b531CandidateError
}

func b531CandidateEqualOptional[b531CandidateValue comparable](b531CandidateLeft, b531CandidateRight *b531CandidateValue) bool {
	if b531CandidateLeft == nil || b531CandidateRight == nil {
		return b531CandidateLeft == nil && b531CandidateRight == nil
	}
	return *b531CandidateLeft == *b531CandidateRight
}

func B531CandidateLookup(b531CandidateKey B531CandidateKey) (B531CandidateEntry, bool, error) {
	b531CandidateCatalog, b531CandidateError := B531Candidates()
	if b531CandidateError != nil {
		return B531CandidateEntry{}, false, b531CandidateError
	}
	for _, b531CandidateEntry := range b531CandidateCatalog.B531CandidateEntries {
		b531CandidateCandidate := b531CandidateEntry.B531CandidateKey
		if b531CandidateCandidate.B531CandidateFamily == b531CandidateKey.B531CandidateFamily &&
			b531CandidateCandidate.B531CandidateVariant == b531CandidateKey.B531CandidateVariant &&
			b531CandidateCandidate.B531CandidatePrimary == b531CandidateKey.B531CandidatePrimary &&
			b531CandidateEqualOptional(b531CandidateCandidate.B531CandidateSecondary, b531CandidateKey.B531CandidateSecondary) &&
			b531CandidateEqualOptional(b531CandidateCandidate.B531CandidateOperation, b531CandidateKey.B531CandidateOperation) &&
			b531CandidateEqualOptional(b531CandidateCandidate.B531CandidateDirection, b531CandidateKey.B531CandidateDirection) &&
			b531CandidateCandidate.B531CandidateField == b531CandidateKey.B531CandidateField {
			return b531CandidateEntry, true, nil
		}
	}
	return B531CandidateEntry{}, false, nil
}
