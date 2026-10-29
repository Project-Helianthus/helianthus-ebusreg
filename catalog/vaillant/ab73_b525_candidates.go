package vaillant

import (
	_ "embed"
	ab73json "encoding/json"
)

//go:embed data/ab73_b525_candidates.json
var ab73B525Data []byte

type Ab73Key struct {
	Ab73Family    string  `json:"family"`
	Ab73Variant   *string `json:"variant"`
	Ab73Primary   byte    `json:"primary"`
	Ab73Secondary *byte   `json:"secondary"`
	Ab73Operation *byte   `json:"operation"`
	Ab73Direction string  `json:"direction"`
	Ab73Field     string  `json:"field"`
}

type Ab73Entry struct {
	Ab73Key        Ab73Key        `json:"key"`
	Ab73Grade      string         `json:"grade"`
	Ab73Position   *int           `json:"position"`
	Ab73Width      *int           `json:"width"`
	Ab73Codec      string         `json:"codec"`
	Ab73Minimum    *int           `json:"minimum"`
	Ab73Maximum    *int           `json:"maximum"`
	Ab73Default    *int           `json:"default"`
	Ab73Count      *int           `json:"count"`
	Ab73StartIndex *int           `json:"start_index"`
	Ab73Unit       *string        `json:"unit"`
	Ab73Mappings   map[int]string `json:"mappings"`
	Ab73Bits       map[int]string `json:"bits"`
	Ab73Unknowns   []string       `json:"unknowns"`
	Ab73RawPolicy  string         `json:"raw_policy"`
	Ab73Relation   string         `json:"relation"`
}

type Ab73Catalog struct {
	Ab73RepresentationGrade   string      `json:"representation_grade"`
	Ab73PhysicalQualification string      `json:"physical_qualification"`
	Ab73CoordinateSystem      string      `json:"coordinate_system"`
	Ab73Limitations           []string    `json:"limitations"`
	Ab73Entries               []Ab73Entry `json:"entries"`
}

func Ab73B525Candidates() Ab73Catalog {
	var ab73Catalog Ab73Catalog
	if ab73Error := ab73json.Unmarshal(ab73B525Data, &ab73Catalog); ab73Error != nil {
		panic(ab73Error)
	}
	return ab73Catalog
}

func ab73EqualOptional[ab73T comparable](ab73Left, ab73Right *ab73T) bool {
	return ab73Left == nil && ab73Right == nil || ab73Left != nil && ab73Right != nil && *ab73Left == *ab73Right
}

func Ab73Lookup(ab73Key Ab73Key) (Ab73Entry, bool) {
	for _, ab73Entry := range Ab73B525Candidates().Ab73Entries {
		ab73Candidate := ab73Entry.Ab73Key
		if ab73Candidate.Ab73Family == ab73Key.Ab73Family &&
			ab73EqualOptional(ab73Candidate.Ab73Variant, ab73Key.Ab73Variant) &&
			ab73Candidate.Ab73Primary == ab73Key.Ab73Primary &&
			ab73EqualOptional(ab73Candidate.Ab73Secondary, ab73Key.Ab73Secondary) &&
			ab73EqualOptional(ab73Candidate.Ab73Operation, ab73Key.Ab73Operation) &&
			ab73Candidate.Ab73Direction == ab73Key.Ab73Direction &&
			ab73Candidate.Ab73Field == ab73Key.Ab73Field {
			return ab73Entry, true
		}
	}
	return Ab73Entry{}, false
}
