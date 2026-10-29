package ebus_standard_catalog

import (
	_ "embed"
	ab79json "encoding/json"
)

type Ab79Key struct {
	Ab79Family    string `json:"family"`
	Ab79Variant   string `json:"variant"`
	Ab79Primary   *uint8 `json:"primary"`
	Ab79Secondary *uint8 `json:"secondary"`
	Ab79Operation *uint8 `json:"operation"`
	Ab79Direction string `json:"direction"`
	Ab79Field     string `json:"field"`
}

type Ab79Entry struct {
	Ab79Key                Ab79Key  `json:"key"`
	Ab79Grade              string   `json:"grade"`
	Ab79AbstractOffset     *int     `json:"abstract_offset"`
	Ab79PhysicalOffset     *int     `json:"physical_offset"`
	Ab79Width              *int     `json:"width"`
	Ab79Encoding           *string  `json:"encoding"`
	Ab79Default            *int     `json:"default"`
	Ab79Minimum            *int     `json:"minimum"`
	Ab79Maximum            *int     `json:"maximum"`
	Ab79AlternativeMaximum *int     `json:"alternative_maximum"`
	Ab79Unit               *string  `json:"unit"`
	Ab79Enum               []string `json:"enum"`
	Ab79Notes              []string `json:"notes"`
}

type Ab79Catalog struct {
	Ab79PresenceGrade  string      `json:"presence_grade"`
	Ab79Scope          string      `json:"scope"`
	Ab79NativeContract string      `json:"native_contract"`
	Ab79UnknownGrade   string      `json:"unknown_grade"`
	Ab79Unknowns       []string    `json:"unknowns"`
	Ab79Entries        []Ab79Entry `json:"entries"`
}

//go:embed testdata/ab79_09_candidates.json
var ab79CandidatesJSON []byte

func Ab79Candidates() Ab79Catalog {
	var ab79Catalog Ab79Catalog
	if ab79Err := ab79json.Unmarshal(ab79CandidatesJSON, &ab79Catalog); ab79Err != nil {
		panic(ab79Err)
	}
	return ab79Catalog
}

func ab79EqualByte(ab79Left, ab79Right *uint8) bool {
	if ab79Left == nil || ab79Right == nil {
		return ab79Left == nil && ab79Right == nil
	}
	return *ab79Left == *ab79Right
}

func Ab79Lookup(ab79Key Ab79Key) (Ab79Entry, bool) {
	for _, ab79Entry := range Ab79Candidates().Ab79Entries {
		ab79Candidate := ab79Entry.Ab79Key
		if ab79Candidate.Ab79Family == ab79Key.Ab79Family &&
			ab79Candidate.Ab79Variant == ab79Key.Ab79Variant &&
			ab79EqualByte(ab79Candidate.Ab79Primary, ab79Key.Ab79Primary) &&
			ab79EqualByte(ab79Candidate.Ab79Secondary, ab79Key.Ab79Secondary) &&
			ab79EqualByte(ab79Candidate.Ab79Operation, ab79Key.Ab79Operation) &&
			ab79Candidate.Ab79Direction == ab79Key.Ab79Direction &&
			ab79Candidate.Ab79Field == ab79Key.Ab79Field {
			return ab79Entry, true
		}
	}
	return Ab79Entry{}, false
}
