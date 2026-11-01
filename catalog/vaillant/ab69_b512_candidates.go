package vaillant

import (
	_ "embed"
	ab69json "encoding/json"
)

//go:embed data/ab69_b512_candidates.json
var ab69B512JSON []byte

type Ab69Key struct {
	Ab69Family    string `json:"family"`
	Ab69Variant   string `json:"variant"`
	Ab69Primary   byte   `json:"primary"`
	Ab69Secondary *byte  `json:"secondary"`
	Ab69Operation *byte  `json:"operation"`
	Ab69Direction string `json:"direction"`
	Ab69Field     string `json:"field"`
}

type Ab69Range struct {
	Ab69Min     int    `json:"min"`
	Ab69Max     int    `json:"max"`
	Ab69Meaning string `json:"meaning"`
}

type Ab69Candidate struct {
	Ab69Key            Ab69Key           `json:"key"`
	Ab69Grade          string            `json:"grade"`
	Ab69Position       *int              `json:"logical_position"`
	Ab69Width          *int              `json:"width"`
	Ab69RawDomain      []int             `json:"raw_domain"`
	Ab69Signed         *bool             `json:"signed"`
	Ab69ByteOrder      *string           `json:"byte_order"`
	Ab69Divisor        *float64          `json:"divisor"`
	Ab69Unit           *string           `json:"unit"`
	Ab69Default        *int              `json:"default"`
	Ab69Mappings       map[string]string `json:"mappings"`
	Ab69Ranges         []Ab69Range       `json:"ranges"`
	Ab69Alternatives   []string          `json:"alternatives"`
	Ab69TimeoutSeconds *int              `json:"timeout_seconds"`
	Ab69TimeoutRestart *string           `json:"timeout_restart"`
	Ab69ExternalMode   *int              `json:"external_mode"`
	Ab69Notes          []string          `json:"notes"`
}

type Ab69Catalog struct {
	Ab69RepresentationGrade   string          `json:"representation_grade"`
	Ab69PhysicalQualification string          `json:"physical_qualification"`
	Ab69InstalledSupport      string          `json:"installed_support"`
	Ab69ResponseAdaptation    string          `json:"response_adaptation"`
	Ab69ActualSelection       string          `json:"actual_versus_fallback_selection"`
	Ab69Scope                 string          `json:"scope"`
	Ab69Entries               []Ab69Candidate `json:"entries"`
}

func Ab69B512Candidates() Ab69Catalog {
	var ab69Catalog Ab69Catalog
	if ab69Err := ab69json.Unmarshal(ab69B512JSON, &ab69Catalog); ab69Err != nil {
		panic(ab69Err)
	}
	return ab69Catalog
}

func ab69OptionalEqual(ab69Left, ab69Right *byte) bool {
	if ab69Left == nil || ab69Right == nil {
		return ab69Left == nil && ab69Right == nil
	}
	return *ab69Left == *ab69Right
}

func Ab69Lookup(ab69Key Ab69Key) (Ab69Candidate, bool) {
	for _, ab69Entry := range Ab69B512Candidates().Ab69Entries {
		ab69Stored := ab69Entry.Ab69Key
		if ab69Stored.Ab69Family == ab69Key.Ab69Family && ab69Stored.Ab69Variant == ab69Key.Ab69Variant && ab69Stored.Ab69Primary == ab69Key.Ab69Primary && ab69OptionalEqual(ab69Stored.Ab69Secondary, ab69Key.Ab69Secondary) && ab69OptionalEqual(ab69Stored.Ab69Operation, ab69Key.Ab69Operation) && ab69Stored.Ab69Direction == ab69Key.Ab69Direction && ab69Stored.Ab69Field == ab69Key.Ab69Field {
			return ab69Entry, true
		}
	}
	return Ab69Candidate{}, false
}
