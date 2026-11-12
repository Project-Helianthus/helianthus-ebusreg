package ebusstandard

import (
	_ "embed"
	standard09CandidateJSON "encoding/json"
)

type Standard09CandidateKey struct {
	Standard09CandidateFamily    string `json:"family"`
	Standard09CandidateVariant   string `json:"variant"`
	Standard09CandidatePrimary   *uint8 `json:"primary"`
	Standard09CandidateSecondary *uint8 `json:"secondary"`
	Standard09CandidateOperation *uint8 `json:"operation"`
	Standard09CandidateDirection string `json:"direction"`
	Standard09CandidateField     string `json:"field"`
}

type Standard09CandidateEntry struct {
	Standard09CandidateKey                Standard09CandidateKey `json:"key"`
	Standard09CandidateGrade              string                 `json:"grade"`
	Standard09CandidateAbstractOffset     *int                   `json:"abstract_offset"`
	Standard09CandidatePhysicalOffset     *int                   `json:"physical_offset"`
	Standard09CandidateWidth              *int                   `json:"width"`
	Standard09CandidateEncoding           *string                `json:"encoding"`
	Standard09CandidateDefault            *int                   `json:"default"`
	Standard09CandidateMinimum            *int                   `json:"minimum"`
	Standard09CandidateMaximum            *int                   `json:"maximum"`
	Standard09CandidateAlternativeMaximum *int                   `json:"alternative_maximum"`
	Standard09CandidateUnit               *string                `json:"unit"`
	Standard09CandidateEnum               []string               `json:"enum"`
	Standard09CandidateNotes              []string               `json:"notes"`
}

type Standard09CandidateCatalog struct {
	Standard09CandidatePresenceGrade  string                     `json:"presence_grade"`
	Standard09CandidateScope          string                     `json:"scope"`
	Standard09CandidateNativeContract string                     `json:"native_contract"`
	Standard09CandidateUnknownGrade   string                     `json:"unknown_grade"`
	Standard09CandidateUnknowns       []string                   `json:"unknowns"`
	Standard09CandidateEntries        []Standard09CandidateEntry `json:"entries"`
}

//go:embed testdata/ab79_09_candidates.json
var standard09CandidateCandidatesJSON []byte

func Standard09Candidates() Standard09CandidateCatalog {
	var standard09CandidateCatalog Standard09CandidateCatalog
	if standard09CandidateErr := standard09CandidateJSON.Unmarshal(standard09CandidateCandidatesJSON, &standard09CandidateCatalog); standard09CandidateErr != nil {
		panic(standard09CandidateErr)
	}
	return standard09CandidateCatalog
}

func standard09CandidateEqualByte(standard09CandidateLeft, standard09CandidateRight *uint8) bool {
	if standard09CandidateLeft == nil || standard09CandidateRight == nil {
		return standard09CandidateLeft == nil && standard09CandidateRight == nil
	}
	return *standard09CandidateLeft == *standard09CandidateRight
}

func Standard09CandidateLookup(standard09CandidateKey Standard09CandidateKey) (Standard09CandidateEntry, bool) {
	for _, standard09CandidateEntry := range Standard09Candidates().Standard09CandidateEntries {
		standard09CandidateCandidate := standard09CandidateEntry.Standard09CandidateKey
		if standard09CandidateCandidate.Standard09CandidateFamily == standard09CandidateKey.Standard09CandidateFamily &&
			standard09CandidateCandidate.Standard09CandidateVariant == standard09CandidateKey.Standard09CandidateVariant &&
			standard09CandidateEqualByte(standard09CandidateCandidate.Standard09CandidatePrimary, standard09CandidateKey.Standard09CandidatePrimary) &&
			standard09CandidateEqualByte(standard09CandidateCandidate.Standard09CandidateSecondary, standard09CandidateKey.Standard09CandidateSecondary) &&
			standard09CandidateEqualByte(standard09CandidateCandidate.Standard09CandidateOperation, standard09CandidateKey.Standard09CandidateOperation) &&
			standard09CandidateCandidate.Standard09CandidateDirection == standard09CandidateKey.Standard09CandidateDirection &&
			standard09CandidateCandidate.Standard09CandidateField == standard09CandidateKey.Standard09CandidateField {
			return standard09CandidateEntry, true
		}
	}
	return Standard09CandidateEntry{}, false
}
