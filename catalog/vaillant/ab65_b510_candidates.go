package vaillant

import (
	_ "embed"
	ab65json "encoding/json"
	ab65reflect "reflect"
)

//go:embed data/ab65_b510_candidates.json
var ab65B510Data []byte

type Ab65Key struct {
	Ab65Family    string  `json:"family"`
	Ab65Variant   string  `json:"variant"`
	Ab65Primary   uint8   `json:"primary"`
	Ab65Secondary *uint8  `json:"secondary"`
	Ab65Operation *uint8  `json:"operation"`
	Ab65Direction *string `json:"direction"`
	Ab65Field     string  `json:"field"`
}

type Ab65Bit struct {
	Ab65Bit                 uint8            `json:"bit"`
	Ab65Grade               string           `json:"grade"`
	Ab65Values              map[uint8]string `json:"values"`
	Ab65Reserved            bool             `json:"reserved"`
	Ab65Expected            *uint8           `json:"expected"`
	Ab65ActiveModes         []string         `json:"active_modes"`
	Ab65IgnoredModes        []string         `json:"ignored_modes"`
	Ab65Meaning             *string          `json:"meaning"`
	Ab65PersistsUntilStatus *bool            `json:"persists_until_status"`
	Ab65CommandsTermination *bool            `json:"commands_termination"`
}

type Ab65Entry struct {
	Ab65Key            Ab65Key          `json:"key"`
	Ab65Position       *int             `json:"position"`
	Ab65WidthBits      *int             `json:"width_bits"`
	Ab65Grade          string           `json:"grade"`
	Ab65RawDomain      []uint8          `json:"raw_domain"`
	Ab65Divisor        *float64         `json:"divisor"`
	Ab65Unit           *string          `json:"unit"`
	Ab65Values         map[uint8]string `json:"values"`
	Ab65Bits           []Ab65Bit        `json:"bits"`
	Ab65Interpretation string           `json:"interpretation"`
}

type Ab65Catalog struct {
	Ab65MetadataPresence      string      `json:"metadata_presence"`
	Ab65PhysicalQualification string      `json:"physical_qualification"`
	Ab65RoutableInvocation    string      `json:"routable_invocation"`
	Ab65NativeContract        string      `json:"native_contract"`
	Ab65Limitations           []string    `json:"limitations"`
	Ab65Entries               []Ab65Entry `json:"entries"`
}

func Ab65B510Candidates() Ab65Catalog {
	var ab65Catalog Ab65Catalog
	if ab65Error := ab65json.Unmarshal(ab65B510Data, &ab65Catalog); ab65Error != nil {
		panic(ab65Error)
	}
	return ab65Catalog
}

func Ab65Lookup(ab65Key Ab65Key) (Ab65Entry, bool) {
	for _, ab65Entry := range Ab65B510Candidates().Ab65Entries {
		if ab65reflect.DeepEqual(ab65Key, ab65Entry.Ab65Key) {
			return ab65Entry, true
		}
	}
	return Ab65Entry{}, false
}
