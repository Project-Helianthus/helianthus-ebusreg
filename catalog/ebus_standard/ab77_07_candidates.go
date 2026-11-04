package ebus_standard_catalog

import (
	_ "embed"
	ab77json "encoding/json"
)

//go:embed testdata/ab77_07_candidates.json
var ab77Data []byte

type Ab77Key struct {
	Ab77Family    string `json:"family"`
	Ab77Variant   string `json:"variant"`
	Ab77Primary   byte   `json:"primary"`
	Ab77Secondary *byte  `json:"secondary"`
	Ab77Operation *byte  `json:"operation"`
	Ab77Direction string `json:"direction"`
	Ab77Field     string `json:"field"`
}

type Ab77Control struct {
	Ab77NotOkayMinimum int    `json:"not_okay_minimum"`
	Ab77NotOkayMaximum int    `json:"not_okay_maximum"`
	Ab77Default        string `json:"default"`
	Ab77OkayValues     []int  `json:"okay_values"`
}

type Ab77LocalParser struct {
	Ab77MinimumDeclaredCount   int    `json:"minimum_declared_count"`
	Ab77MaximumDeclaredCount   *int   `json:"maximum_declared_count"`
	Ab77MinimumAvailableCount  int    `json:"minimum_available_count"`
	Ab77IdentifierOffset       int    `json:"identifier_offset"`
	Ab77IdentifierWidth        int    `json:"identifier_width"`
	Ab77FirstIdentifierNonzero bool   `json:"first_identifier_nonzero"`
	Ab77CharacterZeroAllowed   bool   `json:"character_zero_allowed"`
	Ab77CharacterMinimum       byte   `json:"character_minimum"`
	Ab77CharacterMaximum       byte   `json:"character_maximum"`
	Ab77NumericOffsets         []int  `json:"numeric_offsets"`
	Ab77LowNibbleMaximum       byte   `json:"low_nibble_maximum"`
	Ab77HighNibbleMaximum      *byte  `json:"high_nibble_maximum"`
	Ab77DecimalMultiplier      int    `json:"decimal_multiplier"`
	Ab77MajorOffset            int    `json:"major_offset"`
	Ab77MinorOffset            int    `json:"minor_offset"`
	Ab77ProductOffsets         []int  `json:"product_offsets"`
	Ab77ProductMultiplier      int    `json:"product_multiplier"`
	Ab77ConsumesByteZero       bool   `json:"consumes_byte_zero"`
	Ab77Grade                  string `json:"grade"`
	Ab77Scope                  string `json:"scope"`
}

type Ab77Entry struct {
	Ab77Key                 Ab77Key          `json:"key"`
	Ab77Grade               string           `json:"grade"`
	Ab77LogicalPosition     *int             `json:"logical_position"`
	Ab77Width               *int             `json:"width"`
	Ab77DirectReplyRequired *bool            `json:"direct_reply_required"`
	Ab77Control             *Ab77Control     `json:"control"`
	Ab77LocalParser         *Ab77LocalParser `json:"local_parser"`
}

type Ab77Candidates struct {
	Ab77UnknownGrade  string      `json:"unknown_grade"`
	Ab77PresenceGrade string      `json:"presence_grade"`
	Ab77Scope         string      `json:"scope"`
	Ab77Unknown       []string    `json:"unknown"`
	Ab77Entries       []Ab77Entry `json:"entries"`
}

func Ab77CandidateMetadata() Ab77Candidates {
	var ab77Result Ab77Candidates
	if ab77Err := ab77json.Unmarshal(ab77Data, &ab77Result); ab77Err != nil {
		panic(ab77Err)
	}
	return ab77Result
}

func ab77EqualByte(ab77Left, ab77Right *byte) bool {
	return ab77Left == nil && ab77Right == nil || ab77Left != nil && ab77Right != nil && *ab77Left == *ab77Right
}

func Ab77Lookup(ab77Key Ab77Key) (Ab77Entry, bool) {
	for _, ab77Entry := range Ab77CandidateMetadata().Ab77Entries {
		ab77Candidate := ab77Entry.Ab77Key
		if ab77Candidate.Ab77Family == ab77Key.Ab77Family &&
			ab77Candidate.Ab77Variant == ab77Key.Ab77Variant &&
			ab77Candidate.Ab77Primary == ab77Key.Ab77Primary &&
			ab77EqualByte(ab77Candidate.Ab77Secondary, ab77Key.Ab77Secondary) &&
			ab77EqualByte(ab77Candidate.Ab77Operation, ab77Key.Ab77Operation) &&
			ab77Candidate.Ab77Direction == ab77Key.Ab77Direction &&
			ab77Candidate.Ab77Field == ab77Key.Ab77Field {
			return ab77Entry, true
		}
	}
	return Ab77Entry{}, false
}
