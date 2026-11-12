package ebusstandard

import (
	_ "embed"
	standard07CandidateJSON "encoding/json"
)

//go:embed testdata/ab77_07_candidates.json
var standard07CandidateData []byte

type Standard07CandidateKey struct {
	Standard07CandidateFamily    string `json:"family"`
	Standard07CandidateVariant   string `json:"variant"`
	Standard07CandidatePrimary   byte   `json:"primary"`
	Standard07CandidateSecondary *byte  `json:"secondary"`
	Standard07CandidateOperation *byte  `json:"operation"`
	Standard07CandidateDirection string `json:"direction"`
	Standard07CandidateField     string `json:"field"`
}

type Standard07CandidateControl struct {
	Standard07CandidateNotOkayMinimum int    `json:"not_okay_minimum"`
	Standard07CandidateNotOkayMaximum int    `json:"not_okay_maximum"`
	Standard07CandidateDefault        string `json:"default"`
	Standard07CandidateOkayValues     []int  `json:"okay_values"`
}

type Standard07CandidateLocalParser struct {
	Standard07CandidateMinimumDeclaredCount   int    `json:"minimum_declared_count"`
	Standard07CandidateMaximumDeclaredCount   *int   `json:"maximum_declared_count"`
	Standard07CandidateMinimumAvailableCount  int    `json:"minimum_available_count"`
	Standard07CandidateIdentifierOffset       int    `json:"identifier_offset"`
	Standard07CandidateIdentifierWidth        int    `json:"identifier_width"`
	Standard07CandidateFirstIdentifierNonzero bool   `json:"first_identifier_nonzero"`
	Standard07CandidateCharacterZeroAllowed   bool   `json:"character_zero_allowed"`
	Standard07CandidateCharacterMinimum       byte   `json:"character_minimum"`
	Standard07CandidateCharacterMaximum       byte   `json:"character_maximum"`
	Standard07CandidateNumericOffsets         []int  `json:"numeric_offsets"`
	Standard07CandidateLowNibbleMaximum       byte   `json:"low_nibble_maximum"`
	Standard07CandidateHighNibbleMaximum      *byte  `json:"high_nibble_maximum"`
	Standard07CandidateDecimalMultiplier      int    `json:"decimal_multiplier"`
	Standard07CandidateMajorOffset            int    `json:"major_offset"`
	Standard07CandidateMinorOffset            int    `json:"minor_offset"`
	Standard07CandidateProductOffsets         []int  `json:"product_offsets"`
	Standard07CandidateProductMultiplier      int    `json:"product_multiplier"`
	Standard07CandidateConsumesByteZero       bool   `json:"consumes_byte_zero"`
	Standard07CandidateGrade                  string `json:"grade"`
	Standard07CandidateScope                  string `json:"scope"`
}

type Standard07CandidateEntry struct {
	Standard07CandidateKey                 Standard07CandidateKey          `json:"key"`
	Standard07CandidateGrade               string                          `json:"grade"`
	Standard07CandidateLogicalPosition     *int                            `json:"logical_position"`
	Standard07CandidateWidth               *int                            `json:"width"`
	Standard07CandidateDirectReplyRequired *bool                           `json:"direct_reply_required"`
	Standard07CandidateControl             *Standard07CandidateControl     `json:"control"`
	Standard07CandidateLocalParser         *Standard07CandidateLocalParser `json:"local_parser"`
}

type Standard07Candidates struct {
	Standard07CandidateUnknownGrade  string                     `json:"unknown_grade"`
	Standard07CandidatePresenceGrade string                     `json:"presence_grade"`
	Standard07CandidateScope         string                     `json:"scope"`
	Standard07CandidateUnknown       []string                   `json:"unknown"`
	Standard07CandidateEntries       []Standard07CandidateEntry `json:"entries"`
}

func Standard07CandidateMetadata() Standard07Candidates {
	var standard07CandidateResult Standard07Candidates
	if standard07CandidateErr := standard07CandidateJSON.Unmarshal(standard07CandidateData, &standard07CandidateResult); standard07CandidateErr != nil {
		panic(standard07CandidateErr)
	}
	return standard07CandidateResult
}

func standard07CandidateEqualByte(standard07CandidateLeft, standard07CandidateRight *byte) bool {
	return standard07CandidateLeft == nil && standard07CandidateRight == nil || standard07CandidateLeft != nil && standard07CandidateRight != nil && *standard07CandidateLeft == *standard07CandidateRight
}

func Standard07CandidateLookup(standard07CandidateKey Standard07CandidateKey) (Standard07CandidateEntry, bool) {
	for _, standard07CandidateEntry := range Standard07CandidateMetadata().Standard07CandidateEntries {
		standard07CandidateCandidate := standard07CandidateEntry.Standard07CandidateKey
		if standard07CandidateCandidate.Standard07CandidateFamily == standard07CandidateKey.Standard07CandidateFamily &&
			standard07CandidateCandidate.Standard07CandidateVariant == standard07CandidateKey.Standard07CandidateVariant &&
			standard07CandidateCandidate.Standard07CandidatePrimary == standard07CandidateKey.Standard07CandidatePrimary &&
			standard07CandidateEqualByte(standard07CandidateCandidate.Standard07CandidateSecondary, standard07CandidateKey.Standard07CandidateSecondary) &&
			standard07CandidateEqualByte(standard07CandidateCandidate.Standard07CandidateOperation, standard07CandidateKey.Standard07CandidateOperation) &&
			standard07CandidateCandidate.Standard07CandidateDirection == standard07CandidateKey.Standard07CandidateDirection &&
			standard07CandidateCandidate.Standard07CandidateField == standard07CandidateKey.Standard07CandidateField {
			return standard07CandidateEntry, true
		}
	}
	return Standard07CandidateEntry{}, false
}
