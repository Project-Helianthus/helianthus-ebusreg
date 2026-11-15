package vaillant

import (
	_ "embed"
	b505CandidateJSON "encoding/json"
	b505CandidateReflect "reflect"
)

//go:embed data/b505_candidates.json
var b505CandidateCandidatesJSON string

// B505CandidateKey identifies one B505 candidate field; absent optional axes are nil and lookups
// compare every axis exactly.
type B505CandidateKey struct {
	B505CandidateFamily    string  `json:"family"`
	B505CandidateVariant   string  `json:"variant"`
	B505CandidatePrimary   byte    `json:"primary"`
	B505CandidateSecondary *byte   `json:"secondary"`
	B505CandidateOperation *byte   `json:"operation"`
	B505CandidateDirection *string `json:"direction"`
	B505CandidateField     string  `json:"field"`
}

// B505CandidateEntry declares one B505 candidate field with its position, width, conversion, bounds
// and raw-value branches.
type B505CandidateEntry struct {
	B505CandidateSignedness          string            `json:"signedness"`
	B505CandidateUnit                *string           `json:"unit"`
	B505CandidateFreshness           string            `json:"freshness"`
	B505CandidateExecution           string            `json:"execution"`
	B505CandidateKey                 B505CandidateKey  `json:"key"`
	B505CandidatePosition            *int              `json:"position"`
	B505CandidateWidthBits           *int              `json:"width_bits"`
	B505CandidateGrade               string            `json:"grade"`
	B505CandidateConversion          string            `json:"conversion"`
	B505CandidateBounds              string            `json:"bounds"`
	B505CandidateDurationSeconds     *int              `json:"duration_seconds"`
	B505CandidateBranches            map[string]string `json:"branches"`
	B505CandidateFallback            string            `json:"fallback"`
	B505CandidateDisplaySteps        []int             `json:"display_steps"`
	B505CandidatePriority            map[string]string `json:"priority"`
	B505CandidatePhysicalSupport     string            `json:"physical_support"`
	B505CandidateAcknowledgement     string            `json:"acknowledgement"`
	B505CandidateAcceptance          string            `json:"acceptance"`
	B505CandidatePersistence         string            `json:"persistence"`
	B505CandidateDurationEnforcement string            `json:"duration_enforcement"`
}

// B505CandidateCatalog is the B505 candidate catalog with its qualification and limitations.
type B505CandidateCatalog struct {
	B505CandidateMetadataPresence string               `json:"metadata_presence"`
	B505CandidateQualification    string               `json:"qualification"`
	B505CandidateLimitations      []string             `json:"limitations"`
	B505CandidateEntries          []B505CandidateEntry `json:"entries"`
}

// B505Candidates decodes the embedded B505 candidate catalog on each call and returns a detached copy.
// It panics if the embedded data is malformed.
func B505Candidates() B505CandidateCatalog {
	var b505CandidateCatalog B505CandidateCatalog
	if b505CandidateError := b505CandidateJSON.Unmarshal([]byte(b505CandidateCandidatesJSON), &b505CandidateCatalog); b505CandidateError != nil {
		panic(b505CandidateError)
	}
	return b505CandidateCatalog
}

// B505CandidateLookup returns the candidate entry whose key equals the given key on every axis.
func B505CandidateLookup(b505CandidateKey B505CandidateKey) (B505CandidateEntry, bool) {
	for _, b505CandidateEntry := range B505Candidates().B505CandidateEntries {
		if b505CandidateKeysEqual(b505CandidateEntry.B505CandidateKey, b505CandidateKey) {
			return b505CandidateEntry, true
		}
	}
	return B505CandidateEntry{}, false
}

// B505CandidateSelectRaw returns the branch declared for the raw value, or the fallback when no branch
// matches.
func (b505CandidateEntry B505CandidateEntry) B505CandidateSelectRaw(b505CandidateRaw string) string {
	if b505CandidateBranch, b505CandidateFound := b505CandidateEntry.B505CandidateBranches[b505CandidateRaw]; b505CandidateFound {
		return b505CandidateBranch
	}
	return b505CandidateEntry.B505CandidateFallback
}

func b505CandidateKeysEqual(b505CandidateLeft, b505CandidateRight B505CandidateKey) bool {
	return b505CandidateReflect.DeepEqual(b505CandidateLeft, b505CandidateRight)
}
