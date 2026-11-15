package vaillant

import (
	_ "embed"
	b504CandidateJSON "encoding/json"
	b504CandidateReflect "reflect"
)

// B504CandidateKey identifies one B504 candidate field; absent optional axes are nil and lookups
// compare every axis exactly.
type B504CandidateKey struct {
	B504CandidateFamily    string `json:"family"`
	B504CandidateVariant   string `json:"variant"`
	B504CandidatePrimary   uint8  `json:"primary"`
	B504CandidateSecondary *uint8 `json:"secondary"`
	B504CandidateOperation *uint8 `json:"operation"`
	B504CandidateDirection string `json:"direction"`
	B504CandidateField     string `json:"field"`
}

// B504CandidateDomain is the value domain of a candidate field: listed values or a stepped range.
type B504CandidateDomain struct {
	B504CandidateValues  []int `json:"values"`
	B504CandidateMinimum *int  `json:"minimum"`
	B504CandidateMaximum *int  `json:"maximum"`
	B504CandidateStep    *int  `json:"step"`
}

// B504CandidateEntry declares one B504 candidate field with its position, width, scaling, unit, labels
// and domain.
type B504CandidateEntry struct {
	B504CandidateKey           B504CandidateKey     `json:"key"`
	B504CandidateGrade         SchemaGrade          `json:"grade"`
	B504CandidatePosition      *int                 `json:"position"`
	B504CandidateWidth         *int                 `json:"width_bytes"`
	B504CandidateSignedness    string               `json:"signedness"`
	B504CandidateByteOrder     string               `json:"byte_order"`
	B504CandidateDivisor       *int                 `json:"divisor"`
	B504CandidateUnit          *string              `json:"unit"`
	B504CandidateLabels        map[int]string       `json:"labels"`
	B504CandidatePositiveLabel *string              `json:"positive_label"`
	B504CandidateDomain        *B504CandidateDomain `json:"domain"`
	B504CandidateNotes         []string             `json:"notes"`
}

// B504CandidateCatalog is the B504 candidate catalog with its catalog-wide grades and limitations.
type B504CandidateCatalog struct {
	B504CandidateMetadataPresence   SchemaGrade          `json:"metadata_presence"`
	B504CandidatePhysicalAdaptation SchemaGrade          `json:"physical_adaptation"`
	B504CandidateInstalledSupport   SchemaGrade          `json:"installed_support"`
	B504CandidatePollingPermission  SchemaGrade          `json:"polling_permission"`
	B504CandidateRoutableInvocation SchemaGrade          `json:"routable_invocation"`
	B504CandidateLimitations        []string             `json:"limitations"`
	B504CandidateEntries            []B504CandidateEntry `json:"entries"`
}

//go:embed data/b504_candidates.json
var b504CandidateB504JSON []byte

// B504Candidates decodes the embedded B504 candidate catalog on each call and returns a detached copy
// or the decoding error.
func B504Candidates() (B504CandidateCatalog, error) {
	var b504CandidateCatalog B504CandidateCatalog
	b504CandidateError := b504CandidateJSON.Unmarshal(b504CandidateB504JSON, &b504CandidateCatalog)
	return b504CandidateCatalog, b504CandidateError
}

// B504LookupCandidate returns the candidate entry whose key equals the given key on every axis, or the
// catalog decoding error.
func B504LookupCandidate(b504CandidateKey B504CandidateKey) (B504CandidateEntry, bool, error) {
	b504CandidateCatalog, b504CandidateError := B504Candidates()
	if b504CandidateError != nil {
		return B504CandidateEntry{}, false, b504CandidateError
	}
	for _, b504CandidateEntry := range b504CandidateCatalog.B504CandidateEntries {
		if b504CandidateReflect.DeepEqual(b504CandidateEntry.B504CandidateKey, b504CandidateKey) {
			return b504CandidateEntry, true, nil
		}
	}
	return B504CandidateEntry{}, false, nil
}
