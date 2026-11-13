package vaillant

import (
	_ "embed"
	b510CandidateJSON "encoding/json"
	b510CandidateReflect "reflect"
)

//go:embed data/ab65_b510_candidates.json
var b510CandidateB510Data []byte

type B510CandidateKey struct {
	B510CandidateFamily    string  `json:"family"`
	B510CandidateVariant   string  `json:"variant"`
	B510CandidatePrimary   uint8   `json:"primary"`
	B510CandidateSecondary *uint8  `json:"secondary"`
	B510CandidateOperation *uint8  `json:"operation"`
	B510CandidateDirection *string `json:"direction"`
	B510CandidateField     string  `json:"field"`
}

type B510CandidateBit struct {
	B510CandidateBit                 uint8            `json:"bit"`
	B510CandidateGrade               string           `json:"grade"`
	B510CandidateValues              map[uint8]string `json:"values"`
	B510CandidateReserved            bool             `json:"reserved"`
	B510CandidateExpected            *uint8           `json:"expected"`
	B510CandidateActiveModes         []string         `json:"active_modes"`
	B510CandidateIgnoredModes        []string         `json:"ignored_modes"`
	B510CandidateMeaning             *string          `json:"meaning"`
	B510CandidatePersistsUntilStatus *bool            `json:"persists_until_status"`
	B510CandidateCommandsTermination *bool            `json:"commands_termination"`
}

type B510CandidateEntry struct {
	B510CandidateKey            B510CandidateKey   `json:"key"`
	B510CandidatePosition       *int               `json:"position"`
	B510CandidateWidthBits      *int               `json:"width_bits"`
	B510CandidateGrade          string             `json:"grade"`
	B510CandidateRawDomain      []uint8            `json:"raw_domain"`
	B510CandidateDivisor        *float64           `json:"divisor"`
	B510CandidateUnit           *string            `json:"unit"`
	B510CandidateValues         map[uint8]string   `json:"values"`
	B510CandidateBits           []B510CandidateBit `json:"bits"`
	B510CandidateInterpretation string             `json:"interpretation"`
}

type B510CandidateCatalog struct {
	B510CandidateMetadataPresence      string               `json:"metadata_presence"`
	B510CandidatePhysicalQualification string               `json:"physical_qualification"`
	B510CandidateRoutableInvocation    string               `json:"routable_invocation"`
	B510CandidateNativeContract        string               `json:"native_contract"`
	B510CandidateLimitations           []string             `json:"limitations"`
	B510CandidateEntries               []B510CandidateEntry `json:"entries"`
}

func B510Candidates() B510CandidateCatalog {
	var b510CandidateCatalog B510CandidateCatalog
	if b510CandidateError := b510CandidateJSON.Unmarshal(b510CandidateB510Data, &b510CandidateCatalog); b510CandidateError != nil {
		panic(b510CandidateError)
	}
	return b510CandidateCatalog
}

func B510CandidateLookup(b510CandidateKey B510CandidateKey) (B510CandidateEntry, bool) {
	for _, b510CandidateEntry := range B510Candidates().B510CandidateEntries {
		if b510CandidateReflect.DeepEqual(b510CandidateKey, b510CandidateEntry.B510CandidateKey) {
			return b510CandidateEntry, true
		}
	}
	return B510CandidateEntry{}, false
}
