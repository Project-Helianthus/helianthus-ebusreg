package vaillant

import (
	_ "embed"
	ab57JSON "encoding/json"
	ab57Reflect "reflect"
)

//go:embed data/ab57_b506_candidates.json
var ab57B506JSON []byte

type Ab57Key struct {
	Ab57Family    string  `json:"family"`
	Ab57Variant   string  `json:"variant"`
	Ab57Primary   byte    `json:"primary"`
	Ab57Secondary *byte   `json:"secondary"`
	Ab57Operation *byte   `json:"operation"`
	Ab57Direction *string `json:"direction"`
	Ab57Field     string  `json:"field"`
}

type Ab57Parameter struct {
	Ab57Selector       byte      `json:"selector"`
	Ab57Representation string    `json:"representation"`
	Ab57Metadata       Ab05Field `json:"metadata"`
	Ab57MinimumBytes   *int      `json:"minimum_bytes"`
	Ab57MaximumBytes   *int      `json:"maximum_bytes"`
}

type Ab57Entry struct {
	Ab57Key                    Ab57Key         `json:"key"`
	Ab57Grade                  Ab05Grade       `json:"grade"`
	Ab57Parameters             []Ab57Parameter `json:"parameters"`
	Ab57OutputMask             *byte           `json:"output_mask"`
	Ab57ReservedMask           *byte           `json:"reserved_mask"`
	Ab57ChecksumValues         map[byte]string `json:"checksum_values"`
	Ab57ChecksumInterpretation string          `json:"checksum_interpretation"`
	Ab57ChecksumWidth          *int            `json:"checksum_width"`
	Ab57Algorithm              Ab05Grade       `json:"algorithm"`
	Ab57Completion             Ab05Grade       `json:"completion"`
	Ab57ChecksumData           Ab05Grade       `json:"checksum_data"`
}

type Ab57Catalog struct {
	Ab57Presence      Ab05Grade   `json:"metadata_presence"`
	Ab57Qualification Ab05Grade   `json:"physical_qualification"`
	Ab57Limitations   []string    `json:"limitations"`
	Ab57Entries       []Ab57Entry `json:"entries"`
}

func Ab57B506Candidates() Ab57Catalog {
	var ab57Catalog Ab57Catalog
	if ab57Error := ab57JSON.Unmarshal(ab57B506JSON, &ab57Catalog); ab57Error != nil {
		panic(ab57Error)
	}
	return ab57Catalog
}

func Ab57B506Lookup(ab57Key Ab57Key) (Ab57Entry, bool) {
	for _, ab57Entry := range Ab57B506Candidates().Ab57Entries {
		if ab57Reflect.DeepEqual(ab57Entry.Ab57Key, ab57Key) {
			return ab57Entry, true
		}
	}
	return Ab57Entry{}, false
}

func Ab57ParameterLookup(ab57Entry Ab57Entry, ab57Selector byte) (Ab57Parameter, bool) {
	for _, ab57Parameter := range ab57Entry.Ab57Parameters {
		if ab57Parameter.Ab57Selector == ab57Selector {
			return ab57Parameter, true
		}
	}
	return Ab57Parameter{}, false
}
