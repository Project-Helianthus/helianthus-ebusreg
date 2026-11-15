package vaillant

import (
	_ "embed"
	b50ECandidateJSON "encoding/json"
)

//go:embed data/b50e_candidates.json
var b50ECandidateB50ECandidatesJSON string

// B50ECandidateKey identifies one B50E candidate field; absent optional axes are nil and lookups
// compare every axis exactly.
type B50ECandidateKey struct {
	B50ECandidateFamily    string  `json:"family"`
	B50ECandidateVariant   string  `json:"variant"`
	B50ECandidatePrimary   byte    `json:"primary"`
	B50ECandidateSecondary *byte   `json:"secondary"`
	B50ECandidateOperation *byte   `json:"operation"`
	B50ECandidateDirection *string `json:"direction"`
	B50ECandidateField     string  `json:"field"`
}

// B50ECandidateComparison holds the bit masks that split a B50E comparison byte into mode,
// stop-override and unresolved bits.
type B50ECandidateComparison struct {
	B50ECandidateModeMask             byte `json:"mode_mask"`
	B50ECandidateStopOverrideFIFOMask byte `json:"stop_override_fifo_mask"`
	B50ECandidateUnresolvedMask       byte `json:"unresolved_mask"`
}

// B50ECandidateClearIndex declares the encoding of a rule clear index and the value that selects all
// rules.
type B50ECandidateClearIndex struct {
	B50ECandidateByteOrder string `json:"byte_order"`
	B50ECandidateUnsigned  bool   `json:"unsigned"`
	B50EAllRulesCandidate  uint16 `json:"all_rules_candidate"`
	B50ECandidateMeaning   string `json:"meaning"`
}

// B50ECandidateEntry declares one B50E candidate field with its position, width and optional
// comparison or clear-index declaration.
type B50ECandidateEntry struct {
	B50ECandidateKey         B50ECandidateKey         `json:"key"`
	B50ECandidatePosition    *int                     `json:"position"`
	B50ECandidateWidthBits   *int                     `json:"width_bits"`
	B50ECandidateGrade       string                   `json:"grade"`
	B50ECandidateComparison  *B50ECandidateComparison `json:"comparison"`
	B50ECandidateClearIndex  *B50ECandidateClearIndex `json:"clear_index"`
	B50ECandidateRawFallback string                   `json:"raw_fallback"`
}

// B50ECandidateCatalog is the B50E candidate catalog with its qualification and limitations.
type B50ECandidateCatalog struct {
	B50ECandidateMetadataPresence      string               `json:"metadata_presence"`
	B50ECandidatePhysicalQualification string               `json:"physical_qualification"`
	B50ECandidateLimitations           []string             `json:"limitations"`
	B50ECandidateEntries               []B50ECandidateEntry `json:"entries"`
}

// B50ECandidates decodes the embedded B50E candidate catalog on each call and returns a detached copy
// or the decoding error.
func B50ECandidates() (B50ECandidateCatalog, error) {
	var b50ECandidateCatalog B50ECandidateCatalog
	b50ECandidateError := b50ECandidateJSON.Unmarshal([]byte(b50ECandidateB50ECandidatesJSON), &b50ECandidateCatalog)
	return b50ECandidateCatalog, b50ECandidateError
}

func b50ECandidateOptionalEqual[b50ECandidateValue comparable](b50ECandidateLeft, b50ECandidateRight *b50ECandidateValue) bool {
	if b50ECandidateLeft == nil || b50ECandidateRight == nil {
		return b50ECandidateLeft == nil && b50ECandidateRight == nil
	}
	return *b50ECandidateLeft == *b50ECandidateRight
}

// B50ELookupCandidate returns the candidate entry whose key equals the given key on every axis, or the
// catalog decoding error.
func B50ELookupCandidate(b50ECandidateKey B50ECandidateKey) (B50ECandidateEntry, bool, error) {
	b50ECandidateCatalog, b50ECandidateError := B50ECandidates()
	if b50ECandidateError != nil {
		return B50ECandidateEntry{}, false, b50ECandidateError
	}
	for _, b50ECandidateEntry := range b50ECandidateCatalog.B50ECandidateEntries {
		b50ECandidateStored := b50ECandidateEntry.B50ECandidateKey
		if b50ECandidateStored.B50ECandidateFamily == b50ECandidateKey.B50ECandidateFamily &&
			b50ECandidateStored.B50ECandidateVariant == b50ECandidateKey.B50ECandidateVariant &&
			b50ECandidateStored.B50ECandidatePrimary == b50ECandidateKey.B50ECandidatePrimary &&
			b50ECandidateOptionalEqual(b50ECandidateStored.B50ECandidateSecondary, b50ECandidateKey.B50ECandidateSecondary) &&
			b50ECandidateOptionalEqual(b50ECandidateStored.B50ECandidateOperation, b50ECandidateKey.B50ECandidateOperation) &&
			b50ECandidateOptionalEqual(b50ECandidateStored.B50ECandidateDirection, b50ECandidateKey.B50ECandidateDirection) &&
			b50ECandidateStored.B50ECandidateField == b50ECandidateKey.B50ECandidateField {
			return b50ECandidateEntry, true, nil
		}
	}
	return B50ECandidateEntry{}, false, nil
}
