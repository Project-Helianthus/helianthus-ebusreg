package ebusstandard

import (
	_ "embed"
	standard07CandidateJSON "encoding/json"
)

//go:embed testdata/standard07_candidates.json
var standard07CandidateData []byte

// Standard07CandidateKey identifies one service 07 candidate field; absent optional axes are nil and
// lookups compare every axis exactly.
type Standard07CandidateKey struct {
	Standard07CandidateFamily    string `json:"family"`
	Standard07CandidateVariant   string `json:"variant"`
	Standard07CandidatePrimary   byte   `json:"primary"`
	Standard07CandidateSecondary *byte  `json:"secondary"`
	Standard07CandidateOperation *byte  `json:"operation"`
	Standard07CandidateDirection string `json:"direction"`
	Standard07CandidateField     string `json:"field"`
}

// Standard07CandidateControl declares the not-okay range, the default and the okay values of a control
// field.
type Standard07CandidateControl struct {
	Standard07CandidateNotOkayMinimum int    `json:"not_okay_minimum"`
	Standard07CandidateNotOkayMaximum int    `json:"not_okay_maximum"`
	Standard07CandidateDefault        string `json:"default"`
	Standard07CandidateOkayValues     []int  `json:"okay_values"`
}

// Standard07CandidateLocalParser declares the offsets, widths and limits the local identification
// parser applies to a candidate body.
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

// Standard07CandidateEntry declares one service 07 candidate field with its grade, position, width and
// optional control or parser declarations.
type Standard07CandidateEntry struct {
	Standard07CandidateKey                 Standard07CandidateKey          `json:"key"`
	Standard07CandidateGrade               string                          `json:"grade"`
	Standard07CandidateLogicalPosition     *int                            `json:"logical_position"`
	Standard07CandidateWidth               *int                            `json:"width"`
	Standard07CandidateDirectReplyRequired *bool                           `json:"direct_reply_required"`
	Standard07CandidateControl             *Standard07CandidateControl     `json:"control"`
	Standard07CandidateLocalParser         *Standard07CandidateLocalParser `json:"local_parser"`
}

// Standard07Candidates is the service 07 candidate catalog with its catalog-wide grades, scope and
// unknowns.
type Standard07Candidates struct {
	Standard07CandidateUnknownGrade  string                     `json:"unknown_grade"`
	Standard07CandidatePresenceGrade string                     `json:"presence_grade"`
	Standard07CandidateScope         string                     `json:"scope"`
	Standard07CandidateUnknown       []string                   `json:"unknown"`
	Standard07CandidateEntries       []Standard07CandidateEntry `json:"entries"`
}

// Standard07CandidateMetadata decodes the embedded service 07 candidate catalog on each call and
// returns a detached copy. It panics if the embedded data is malformed.
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

// Standard07CandidateLookup returns the candidate entry whose key equals the given key on every axis.
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
