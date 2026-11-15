package vaillant

import (
	_ "embed"
	b506CandidateJSON "encoding/json"
	b506CandidateReflect "reflect"
)

//go:embed data/b506_candidates.json
var b506CandidateB506JSON []byte

// B506CandidateKey identifies one B506 candidate; absent optional axes are nil and lookups compare
// every axis exactly.
type B506CandidateKey struct {
	B506CandidateFamily    string  `json:"family"`
	B506CandidateVariant   string  `json:"variant"`
	B506CandidatePrimary   byte    `json:"primary"`
	B506CandidateSecondary *byte   `json:"secondary"`
	B506CandidateOperation *byte   `json:"operation"`
	B506CandidateDirection *string `json:"direction"`
	B506CandidateField     string  `json:"field"`
}

// B506CandidateParameter declares one selector-addressed parameter of a B506 candidate with its
// representation and length bounds.
type B506CandidateParameter struct {
	B506CandidateSelector       byte        `json:"selector"`
	B506CandidateRepresentation string      `json:"representation"`
	B506CandidateMetadata       SchemaField `json:"metadata"`
	B506CandidateMinimumBytes   *int        `json:"minimum_bytes"`
	B506CandidateMaximumBytes   *int        `json:"maximum_bytes"`
}

// B506CandidateEntry declares one B506 candidate with its parameters, masks and checksum declarations.
type B506CandidateEntry struct {
	B506CandidateKey                    B506CandidateKey         `json:"key"`
	B506CandidateGrade                  SchemaGrade              `json:"grade"`
	B506CandidateParameters             []B506CandidateParameter `json:"parameters"`
	B506CandidateOutputMask             *byte                    `json:"output_mask"`
	B506CandidateReservedMask           *byte                    `json:"reserved_mask"`
	B506CandidateChecksumValues         map[byte]string          `json:"checksum_values"`
	B506CandidateChecksumInterpretation string                   `json:"checksum_interpretation"`
	B506CandidateChecksumWidth          *int                     `json:"checksum_width"`
	B506CandidateAlgorithm              SchemaGrade              `json:"algorithm"`
	B506CandidateCompletion             SchemaGrade              `json:"completion"`
	B506CandidateChecksumData           SchemaGrade              `json:"checksum_data"`
}

// B506CandidateCatalog is the B506 candidate catalog with its presence, qualification and limitations.
type B506CandidateCatalog struct {
	B506CandidatePresence      SchemaGrade          `json:"metadata_presence"`
	B506CandidateQualification SchemaGrade          `json:"physical_qualification"`
	B506CandidateLimitations   []string             `json:"limitations"`
	B506CandidateEntries       []B506CandidateEntry `json:"entries"`
}

// B506Candidates decodes the embedded B506 candidate catalog on each call and returns a detached copy.
// It panics if the embedded data is malformed.
func B506Candidates() B506CandidateCatalog {
	var b506CandidateCatalog B506CandidateCatalog
	if b506CandidateError := b506CandidateJSON.Unmarshal(b506CandidateB506JSON, &b506CandidateCatalog); b506CandidateError != nil {
		panic(b506CandidateError)
	}
	return b506CandidateCatalog
}

// B506CandidateLookup returns the candidate entry whose key equals the given key on every axis.
func B506CandidateLookup(b506CandidateKey B506CandidateKey) (B506CandidateEntry, bool) {
	for _, b506CandidateEntry := range B506Candidates().B506CandidateEntries {
		if b506CandidateReflect.DeepEqual(b506CandidateEntry.B506CandidateKey, b506CandidateKey) {
			return b506CandidateEntry, true
		}
	}
	return B506CandidateEntry{}, false
}

// B506CandidateParameterLookup returns the parameter of the entry declared for the selector.
func B506CandidateParameterLookup(b506CandidateEntry B506CandidateEntry, b506CandidateSelector byte) (B506CandidateParameter, bool) {
	for _, b506CandidateParameter := range b506CandidateEntry.B506CandidateParameters {
		if b506CandidateParameter.B506CandidateSelector == b506CandidateSelector {
			return b506CandidateParameter, true
		}
	}
	return B506CandidateParameter{}, false
}
