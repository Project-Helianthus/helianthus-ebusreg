package vaillant

import (
	_ "embed"
	ab61JSON "encoding/json"
)

//go:embed data/ab61_b50e_candidates.json
var ab61B50ECandidatesJSON string

type Ab61Key struct {
	Ab61Family    string  `json:"family"`
	Ab61Variant   string  `json:"variant"`
	Ab61Primary   byte    `json:"primary"`
	Ab61Secondary *byte   `json:"secondary"`
	Ab61Operation *byte   `json:"operation"`
	Ab61Direction *string `json:"direction"`
	Ab61Field     string  `json:"field"`
}

type Ab61Comparison struct {
	Ab61ModeMask             byte `json:"mode_mask"`
	Ab61StopOverrideFIFOMask byte `json:"stop_override_fifo_mask"`
	Ab61UnresolvedMask       byte `json:"unresolved_mask"`
}

type Ab61ClearIndex struct {
	Ab61ByteOrder         string `json:"byte_order"`
	Ab61Unsigned          bool   `json:"unsigned"`
	Ab61AllRulesCandidate uint16 `json:"all_rules_candidate"`
	Ab61Meaning           string `json:"meaning"`
}

type Ab61Entry struct {
	Ab61Key         Ab61Key         `json:"key"`
	Ab61Position    *int            `json:"position"`
	Ab61WidthBits   *int            `json:"width_bits"`
	Ab61Grade       string          `json:"grade"`
	Ab61Comparison  *Ab61Comparison `json:"comparison"`
	Ab61ClearIndex  *Ab61ClearIndex `json:"clear_index"`
	Ab61RawFallback string          `json:"raw_fallback"`
}

type Ab61Catalog struct {
	Ab61MetadataPresence      string      `json:"metadata_presence"`
	Ab61PhysicalQualification string      `json:"physical_qualification"`
	Ab61Limitations           []string    `json:"limitations"`
	Ab61Entries               []Ab61Entry `json:"entries"`
}

func Ab61B50ECandidates() (Ab61Catalog, error) {
	var ab61Catalog Ab61Catalog
	ab61Error := ab61JSON.Unmarshal([]byte(ab61B50ECandidatesJSON), &ab61Catalog)
	return ab61Catalog, ab61Error
}

func ab61OptionalEqual[ab61Value comparable](ab61Left, ab61Right *ab61Value) bool {
	if ab61Left == nil || ab61Right == nil {
		return ab61Left == nil && ab61Right == nil
	}
	return *ab61Left == *ab61Right
}

func Ab61LookupB50ECandidate(ab61Key Ab61Key) (Ab61Entry, bool, error) {
	ab61Catalog, ab61Error := Ab61B50ECandidates()
	if ab61Error != nil {
		return Ab61Entry{}, false, ab61Error
	}
	for _, ab61Entry := range ab61Catalog.Ab61Entries {
		ab61Stored := ab61Entry.Ab61Key
		if ab61Stored.Ab61Family == ab61Key.Ab61Family &&
			ab61Stored.Ab61Variant == ab61Key.Ab61Variant &&
			ab61Stored.Ab61Primary == ab61Key.Ab61Primary &&
			ab61OptionalEqual(ab61Stored.Ab61Secondary, ab61Key.Ab61Secondary) &&
			ab61OptionalEqual(ab61Stored.Ab61Operation, ab61Key.Ab61Operation) &&
			ab61OptionalEqual(ab61Stored.Ab61Direction, ab61Key.Ab61Direction) &&
			ab61Stored.Ab61Field == ab61Key.Ab61Field {
			return ab61Entry, true, nil
		}
	}
	return Ab61Entry{}, false, nil
}
