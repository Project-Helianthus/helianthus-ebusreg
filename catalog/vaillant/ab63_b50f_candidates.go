package vaillant

import (
	_ "embed"
	ab63json "encoding/json"
	ab63reflect "reflect"
)

//go:embed data/ab63_b50f_candidates.json
var ab63B50FData string

type Ab63Key struct {
	Ab63Family    string  `json:"family"`
	Ab63Variant   string  `json:"variant"`
	Ab63Primary   byte    `json:"primary"`
	Ab63Secondary *byte   `json:"secondary"`
	Ab63Operation *byte   `json:"operation"`
	Ab63Direction *string `json:"direction"`
	Ab63Field     string  `json:"field"`
}

type Ab63Entry struct {
	Ab63Key            Ab63Key         `json:"key"`
	Ab63Field          Ab05Field       `json:"metadata"`
	Ab63Masks          map[string]byte `json:"demand_masks"`
	Ab63Reserved       *byte           `json:"reserved_mask"`
	Ab63Divisor        *int            `json:"divisor"`
	Ab63Conflict       []string        `json:"conversion_conflict"`
	Ab63Interpretation string          `json:"interpretation"`
	Ab63Coordinates    string          `json:"coordinates"`
	Ab63Exclusions     []string        `json:"exclusions"`
}

type Ab63Catalog struct {
	Ab63Presence      Ab05Grade   `json:"metadata_presence"`
	Ab63Qualification Ab05Grade   `json:"physical_qualification"`
	Ab63Invocation    Ab05Grade   `json:"routable_invocation"`
	Ab63Limitations   []string    `json:"limitations"`
	Ab63Entries       []Ab63Entry `json:"entries"`
}

func Ab63B50FCandidates() Ab63Catalog {
	var ab63Catalog Ab63Catalog
	if ab63Error := ab63json.Unmarshal([]byte(ab63B50FData), &ab63Catalog); ab63Error != nil {
		panic(ab63Error)
	}
	return ab63Catalog
}

func Ab63B50FLookup(ab63Key Ab63Key) (Ab63Entry, bool) {
	for _, ab63Entry := range Ab63B50FCandidates().Ab63Entries {
		if ab63reflect.DeepEqual(ab63Entry.Ab63Key, ab63Key) {
			return ab63Entry, true
		}
	}
	return Ab63Entry{}, false
}
