package ebusstandard

import (
	_ "embed"
	standard0fCandidateJSON "encoding/json"
)

//go:embed testdata/ab81_0f_candidates.json
var standard0fCandidateCandidatesJSON []byte

type Standard0fCandidateKey struct {
	Standard0fCandidateFamily    string `json:"family"`
	Standard0fCandidateVariant   string `json:"variant"`
	Standard0fCandidatePrimary   byte   `json:"primary"`
	Standard0fCandidateSecondary *byte  `json:"secondary"`
	Standard0fCandidateOperation *byte  `json:"operation"`
	Standard0fCandidateDirection string `json:"direction"`
	Standard0fCandidateField     string `json:"field"`
}

type Standard0fCandidatePositive struct {
	Standard0fCandidatePrefix           string `json:"prefix"`
	Standard0fCandidateTypedDataPresent *bool  `json:"typed_data_present"`
	Standard0fCandidateGrade            string `json:"grade"`
}

type Standard0fCandidate struct {
	Standard0fCandidateKey           Standard0fCandidateKey       `json:"key"`
	Standard0fCandidateMetadataGrade string                       `json:"metadata_grade"`
	Standard0fCandidateGrade         string                       `json:"grade"`
	Standard0fCandidateBody          string                       `json:"body"`
	Standard0fCandidateOffset        *int                         `json:"offset"`
	Standard0fCandidateWidth         *int                         `json:"width"`
	Standard0fCandidatePhysicalNN    *int                         `json:"physical_nn"`
	Standard0fCandidateACK           string                       `json:"ack"`
	Standard0fCandidateSupport       string                       `json:"support"`
	Standard0fCandidateEffects       string                       `json:"effects"`
	Standard0fCandidateExecutable    bool                         `json:"executable"`
	Standard0fCandidateSendable      bool                         `json:"sendable"`
	Standard0fCandidateRaw           []byte                       `json:"raw"`
	Standard0fCandidatePositive      *Standard0fCandidatePositive `json:"positive"`
	Standard0fCandidateLimits        []string                     `json:"limits"`
}

func Standard0fCandidates() []Standard0fCandidate {
	var standard0fCandidateEntries []Standard0fCandidate
	if standard0fCandidateErr := standard0fCandidateJSON.Unmarshal(standard0fCandidateCandidatesJSON, &standard0fCandidateEntries); standard0fCandidateErr != nil {
		panic("standard0f embedded candidate metadata: " + standard0fCandidateErr.Error())
	}
	return standard0fCandidateEntries
}

func standard0fCandidateOptionalEqual(standard0fCandidateLeft, standard0fCandidateRight *byte) bool {
	return standard0fCandidateLeft == nil && standard0fCandidateRight == nil || standard0fCandidateLeft != nil && standard0fCandidateRight != nil && *standard0fCandidateLeft == *standard0fCandidateRight
}

func Standard0fCandidateLookup(standard0fCandidateKey Standard0fCandidateKey) (Standard0fCandidate, bool) {
	for _, standard0fCandidateEntry := range Standard0fCandidates() {
		standard0fCandidateStored := standard0fCandidateEntry.Standard0fCandidateKey
		if standard0fCandidateStored.Standard0fCandidateFamily == standard0fCandidateKey.Standard0fCandidateFamily &&
			standard0fCandidateStored.Standard0fCandidateVariant == standard0fCandidateKey.Standard0fCandidateVariant &&
			standard0fCandidateStored.Standard0fCandidatePrimary == standard0fCandidateKey.Standard0fCandidatePrimary &&
			standard0fCandidateOptionalEqual(standard0fCandidateStored.Standard0fCandidateSecondary, standard0fCandidateKey.Standard0fCandidateSecondary) &&
			standard0fCandidateOptionalEqual(standard0fCandidateStored.Standard0fCandidateOperation, standard0fCandidateKey.Standard0fCandidateOperation) &&
			standard0fCandidateStored.Standard0fCandidateDirection == standard0fCandidateKey.Standard0fCandidateDirection &&
			standard0fCandidateStored.Standard0fCandidateField == standard0fCandidateKey.Standard0fCandidateField {
			return standard0fCandidateEntry, true
		}
	}
	return Standard0fCandidate{}, false
}
