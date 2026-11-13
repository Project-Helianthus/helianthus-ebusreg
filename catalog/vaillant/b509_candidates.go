package vaillant

import (
	_ "embed"
	b509CandidateJSON "encoding/json"
	b509CandidateReflect "reflect"
)

//go:embed data/ab59_b509_candidates.json
var b509CandidateB509CandidatesJSON string

type B509CandidateKey struct {
	B509CandidateFamily    string  `json:"family"`
	B509CandidateVariant   string  `json:"variant"`
	B509CandidatePrimary   byte    `json:"primary"`
	B509CandidateSecondary *byte   `json:"secondary"`
	B509CandidateOperation *byte   `json:"operation"`
	B509CandidateDirection *string `json:"direction"`
	B509CandidateField     string  `json:"field"`
}

type B509CandidateEntry struct {
	B509CandidateKey           B509CandidateKey  `json:"key"`
	B509CandidatePosition      *int              `json:"position"`
	B509CandidateWidthBits     *int              `json:"width_bits"`
	B509CandidateSignedness    string            `json:"signedness"`
	B509CandidateByteOrder     string            `json:"byte_order"`
	B509CandidateConversion    string            `json:"conversion"`
	B509CandidateDivisor       *int              `json:"divisor"`
	B509CandidateUnit          *string           `json:"unit"`
	B509CandidateDefault       *int64            `json:"default"`
	B509CandidateLower         *int64            `json:"lower"`
	B509CandidateUpper         *int64            `json:"upper"`
	B509CandidateMappings      map[string]string `json:"mappings"`
	B509CandidateOtherFeedback *string           `json:"other_feedback"`
	B509CandidateState         string            `json:"state"`
	B509CandidateGrade         SchemaGrade       `json:"grade"`
	B509CandidateNotes         []string          `json:"notes"`
}

type B509CandidateCatalog struct {
	B509CandidateMetadataPresence SchemaGrade          `json:"metadata_presence"`
	B509CandidateWireAdaptation   SchemaGrade          `json:"wire_adaptation"`
	B509CandidateInstalledSupport SchemaGrade          `json:"installed_support"`
	B509CandidateInvocation       SchemaGrade          `json:"invocation"`
	B509CandidateLimitations      []string             `json:"limitations"`
	B509CandidateEntries          []B509CandidateEntry `json:"entries"`
}

func B509Candidates() (B509CandidateCatalog, error) {
	var b509CandidateCatalog B509CandidateCatalog
	b509CandidateError := b509CandidateJSON.Unmarshal([]byte(b509CandidateB509CandidatesJSON), &b509CandidateCatalog)
	return b509CandidateCatalog, b509CandidateError
}

func B509CandidateLookup(b509CandidateKey B509CandidateKey) (B509CandidateEntry, bool, error) {
	b509CandidateCatalog, b509CandidateError := B509Candidates()
	if b509CandidateError != nil {
		return B509CandidateEntry{}, false, b509CandidateError
	}
	b509CandidateEntry, b509CandidateFound := b509CandidateFind(b509CandidateCatalog.B509CandidateEntries, b509CandidateKey)
	return b509CandidateEntry, b509CandidateFound, nil
}

func b509CandidateFind(b509CandidateEntries []B509CandidateEntry, b509CandidateKey B509CandidateKey) (B509CandidateEntry, bool) {
	for _, b509CandidateEntry := range b509CandidateEntries {
		if b509CandidateReflect.DeepEqual(b509CandidateEntry.B509CandidateKey, b509CandidateKey) {
			return b509CandidateEntry, true
		}
	}
	return B509CandidateEntry{}, false
}
