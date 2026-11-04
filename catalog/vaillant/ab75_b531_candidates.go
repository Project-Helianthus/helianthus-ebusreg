package vaillant

import (
	_ "embed"
	ab75JSON "encoding/json"
)

type Ab75Key struct {
	Ab75Family    string  `json:"family"`
	Ab75Variant   string  `json:"variant"`
	Ab75Primary   byte    `json:"primary"`
	Ab75Secondary *byte   `json:"secondary"`
	Ab75Operation *byte   `json:"operation"`
	Ab75Direction *string `json:"direction"`
	Ab75Field     string  `json:"field"`
}

type Ab75Projection struct {
	Ab75Conversion string   `json:"conversion"`
	Ab75Multiplier uint32   `json:"multiplier"`
	Ab75WidthBits  *int     `json:"width_bits"`
	Ab75Domain     []uint32 `json:"domain"`
	Ab75Unit       string   `json:"unit"`
	Ab75Grade      string   `json:"grade"`
}

type Ab75MaximumTable struct {
	Ab75SourceField    string `json:"source_field"`
	Ab75IndexDomain    []int  `json:"index_domain"`
	Ab75InstalledCount *int   `json:"installed_count"`
	Ab75Grade          string `json:"grade"`
}

type Ab75Entry struct {
	Ab75Key             Ab75Key           `json:"key"`
	Ab75Position        *int              `json:"position"`
	Ab75WidthBits       *int              `json:"width_bits"`
	Ab75RawDomain       []uint32          `json:"raw_domain"`
	Ab75Signedness      string            `json:"signedness"`
	Ab75ByteOrder       string            `json:"byte_order"`
	Ab75NativeUnit      *string           `json:"native_unit"`
	Ab75NativeUnitGrade string            `json:"native_unit_grade"`
	Ab75Sentinel        *uint16           `json:"sentinel"`
	Ab75Grade           string            `json:"grade"`
	Ab75LocalProjection Ab75Projection    `json:"local_projection"`
	Ab75MaximumTable    *Ab75MaximumTable `json:"maximum_table"`
}

type Ab75Catalog struct {
	Ab75MetadataPresence   string      `json:"metadata_presence"`
	Ab75PhysicalSupport    string      `json:"physical_support"`
	Ab75EndpointSupport    string      `json:"endpoint_support"`
	Ab75Correlation        string      `json:"correlation"`
	Ab75ReplyAdaptation    string      `json:"reply_adaptation"`
	Ab75CacheFreshness     string      `json:"cache_freshness"`
	Ab75TargetEligibility  string      `json:"target_eligibility"`
	Ab75AppliedLimits      string      `json:"applied_limits"`
	Ab75OrderingConstraint *string     `json:"ordering_constraint"`
	Ab75Limitations        []string    `json:"limitations"`
	Ab75Entries            []Ab75Entry `json:"entries"`
}

//go:embed data/ab75_b531_candidates.json
var ab75B531JSON []byte

func Ab75B531Candidates() (Ab75Catalog, error) {
	var ab75Catalog Ab75Catalog
	ab75Error := ab75JSON.Unmarshal(ab75B531JSON, &ab75Catalog)
	return ab75Catalog, ab75Error
}

func ab75EqualOptional[ab75Value comparable](ab75Left, ab75Right *ab75Value) bool {
	if ab75Left == nil || ab75Right == nil {
		return ab75Left == nil && ab75Right == nil
	}
	return *ab75Left == *ab75Right
}

func Ab75LookupB531(ab75Key Ab75Key) (Ab75Entry, bool, error) {
	ab75Catalog, ab75Error := Ab75B531Candidates()
	if ab75Error != nil {
		return Ab75Entry{}, false, ab75Error
	}
	for _, ab75Entry := range ab75Catalog.Ab75Entries {
		ab75Candidate := ab75Entry.Ab75Key
		if ab75Candidate.Ab75Family == ab75Key.Ab75Family &&
			ab75Candidate.Ab75Variant == ab75Key.Ab75Variant &&
			ab75Candidate.Ab75Primary == ab75Key.Ab75Primary &&
			ab75EqualOptional(ab75Candidate.Ab75Secondary, ab75Key.Ab75Secondary) &&
			ab75EqualOptional(ab75Candidate.Ab75Operation, ab75Key.Ab75Operation) &&
			ab75EqualOptional(ab75Candidate.Ab75Direction, ab75Key.Ab75Direction) &&
			ab75Candidate.Ab75Field == ab75Key.Ab75Field {
			return ab75Entry, true, nil
		}
	}
	return Ab75Entry{}, false, nil
}
