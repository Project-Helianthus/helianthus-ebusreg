package vaillant

import (
	ab13bytes "bytes"
	ab13embed "embed"
	ab13json "encoding/json"
	ab13fmt "fmt"
	ab13io "io"
	ab13reflect "reflect"
)

//go:embed data/ab13_b507_energy_candidates.json
var ab13Files ab13embed.FS

var ab13Data, ab13ReadError = ab13Files.ReadFile("data/ab13_b507_energy_candidates.json")

type Ab13Key struct {
	Ab13Family            string  `json:"family"`
	Ab13Variant           string  `json:"variant"`
	Ab13Primary           uint8   `json:"primary"`
	Ab13Secondary         *uint8  `json:"secondary"`
	Ab13Operation         *uint8  `json:"operation"`
	Ab13Direction         *string `json:"direction"`
	Ab13Field             string  `json:"field"`
	Ab13Group             *uint8  `json:"group"`
	Ab13RequestedInstance *uint8  `json:"requested_instance"`
	Ab13Register          *uint16 `json:"register"`
	Ab13CommandRole       *string `json:"command_role"`
}

type Ab13Entry struct {
	Ab13Key                   Ab13Key           `json:"key"`
	Ab13Grade                 string            `json:"grade"`
	Ab13MetadataPresence      string            `json:"metadata_presence"`
	Ab13PhysicalQualification string            `json:"physical_qualification"`
	Ab13Invocation            string            `json:"invocation"`
	Ab13NativeContract        string            `json:"native_contract"`
	Ab13Position              *uint16           `json:"position"`
	Ab13WidthBits             *uint8            `json:"width_bits"`
	Ab13Encoding              *string           `json:"encoding"`
	Ab13RawDomain             []uint16          `json:"raw_domain"`
	Ab13Divisor               *float64          `json:"divisor"`
	Ab13Unit                  *string           `json:"unit"`
	Ab13Values                map[uint16]string `json:"values"`
	Ab13Dependencies          []string          `json:"dependencies"`
	Ab13Reference             string            `json:"reference"`
	Ab13QuantityOrigin        string            `json:"quantity_origin"`
	Ab13PreserveRaw           bool              `json:"preserve_raw"`
	Ab13Unknowns              []string          `json:"unknowns"`
	Ab13Alternatives          []string          `json:"alternatives"`
}

func Ab13Entries() ([]Ab13Entry, error) {
	if ab13ReadError != nil {
		return nil, ab13ReadError
	}
	var ab13Entries []Ab13Entry
	ab13Decoder := ab13json.NewDecoder(ab13bytes.NewReader(ab13Data))
	ab13Decoder.DisallowUnknownFields()
	if ab13Error := ab13Decoder.Decode(&ab13Entries); ab13Error != nil {
		return nil, ab13Error
	}
	var ab13Extra any
	if ab13Error := ab13Decoder.Decode(&ab13Extra); ab13Error != ab13io.EOF {
		return nil, ab13fmt.Errorf("ab13: trailing embedded data")
	}
	if len(ab13Entries) == 0 {
		return nil, ab13fmt.Errorf("ab13: empty embedded catalog")
	}
	for ab13Index, ab13Entry := range ab13Entries {
		if ab13Entry.Ab13Key.Ab13Family == "" || ab13Entry.Ab13Key.Ab13Variant == "" || ab13Entry.Ab13Key.Ab13Field == "" ||
			(ab13Entry.Ab13Grade != "Proven" && ab13Entry.Ab13Grade != "Hypothesis" && ab13Entry.Ab13Grade != "Unknown") ||
			ab13Entry.Ab13MetadataPresence != "Proven" || ab13Entry.Ab13PhysicalQualification != "Unknown" ||
			ab13Entry.Ab13Invocation != "Unknown" || ab13Entry.Ab13NativeContract == "" || !ab13Entry.Ab13PreserveRaw {
			return nil, ab13fmt.Errorf("ab13: invalid embedded entry %d", ab13Index)
		}
		if ab13Entry.Ab13Divisor != nil && *ab13Entry.Ab13Divisor <= 0 {
			return nil, ab13fmt.Errorf("ab13: invalid divisor at %d", ab13Index)
		}
		if ab13Entry.Ab13RawDomain != nil && (len(ab13Entry.Ab13RawDomain) != 2 || ab13Entry.Ab13RawDomain[0] > ab13Entry.Ab13RawDomain[1]) {
			return nil, ab13fmt.Errorf("ab13: invalid domain at %d", ab13Index)
		}
		for _, ab13Previous := range ab13Entries[:ab13Index] {
			if ab13reflect.DeepEqual(ab13Previous.Ab13Key, ab13Entry.Ab13Key) {
				return nil, ab13fmt.Errorf("ab13: duplicate embedded key at %d", ab13Index)
			}
		}
	}
	return ab13Entries, nil
}

func Ab13Lookup(ab13Key Ab13Key) (Ab13Entry, bool, error) {
	ab13Entries, ab13Error := Ab13Entries()
	if ab13Error != nil {
		return Ab13Entry{}, false, ab13Error
	}
	for _, ab13Entry := range ab13Entries {
		if ab13reflect.DeepEqual(ab13Key, ab13Entry.Ab13Key) {
			return ab13Entry, true, nil
		}
	}
	return Ab13Entry{}, false, nil
}
