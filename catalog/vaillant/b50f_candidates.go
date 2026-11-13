package vaillant

import (
	_ "embed"
	b50FCandidateJSON "encoding/json"
	b50FCandidateReflect "reflect"
)

//go:embed data/ab63_b50f_candidates.json
var b50FCandidateB50FData string

type B50FCandidateKey struct {
	B50FCandidateFamily    string  `json:"family"`
	B50FCandidateVariant   string  `json:"variant"`
	B50FCandidatePrimary   byte    `json:"primary"`
	B50FCandidateSecondary *byte   `json:"secondary"`
	B50FCandidateOperation *byte   `json:"operation"`
	B50FCandidateDirection *string `json:"direction"`
	B50FCandidateField     string  `json:"field"`
}

type B50FCandidateEntry struct {
	B50FCandidateKey            B50FCandidateKey `json:"key"`
	B50FCandidateField          SchemaField      `json:"metadata"`
	B50FCandidateMasks          map[string]byte  `json:"demand_masks"`
	B50FCandidateReserved       *byte            `json:"reserved_mask"`
	B50FCandidateDivisor        *int             `json:"divisor"`
	B50FCandidateConflict       []string         `json:"conversion_conflict"`
	B50FCandidateInterpretation string           `json:"interpretation"`
	B50FCandidateCoordinates    string           `json:"coordinates"`
	B50FCandidateExclusions     []string         `json:"exclusions"`
}

type B50FCandidateCatalog struct {
	B50FCandidatePresence      SchemaGrade          `json:"metadata_presence"`
	B50FCandidateQualification SchemaGrade          `json:"physical_qualification"`
	B50FCandidateInvocation    SchemaGrade          `json:"routable_invocation"`
	B50FCandidateLimitations   []string             `json:"limitations"`
	B50FCandidateEntries       []B50FCandidateEntry `json:"entries"`
}

func B50FCandidates() B50FCandidateCatalog {
	var b50FCandidateCatalog B50FCandidateCatalog
	if b50FCandidateError := b50FCandidateJSON.Unmarshal([]byte(b50FCandidateB50FData), &b50FCandidateCatalog); b50FCandidateError != nil {
		panic(b50FCandidateError)
	}
	return b50FCandidateCatalog
}

func B50FCandidateLookup(b50FCandidateKey B50FCandidateKey) (B50FCandidateEntry, bool) {
	for _, b50FCandidateEntry := range B50FCandidates().B50FCandidateEntries {
		if b50FCandidateReflect.DeepEqual(b50FCandidateEntry.B50FCandidateKey, b50FCandidateKey) {
			return b50FCandidateEntry, true
		}
	}
	return B50FCandidateEntry{}, false
}
