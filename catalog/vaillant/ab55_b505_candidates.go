package vaillant

import (
	_ "embed"
	ab55JSON "encoding/json"
	ab55Reflect "reflect"
)

//go:embed data/ab55_b505_candidates.json
var ab55CandidatesJSON string

type Ab55Key struct {
	Ab55Family    string  `json:"family"`
	Ab55Variant   string  `json:"variant"`
	Ab55Primary   byte    `json:"primary"`
	Ab55Secondary *byte   `json:"secondary"`
	Ab55Operation *byte   `json:"operation"`
	Ab55Direction *string `json:"direction"`
	Ab55Field     string  `json:"field"`
}

type Ab55Entry struct {
	Ab55Signedness          string            `json:"signedness"`
	Ab55Unit                *string           `json:"unit"`
	Ab55Freshness           string            `json:"freshness"`
	Ab55Execution           string            `json:"execution"`
	Ab55Key                 Ab55Key           `json:"key"`
	Ab55Position            *int              `json:"position"`
	Ab55WidthBits           *int              `json:"width_bits"`
	Ab55Grade               string            `json:"grade"`
	Ab55Conversion          string            `json:"conversion"`
	Ab55Bounds              string            `json:"bounds"`
	Ab55DurationSeconds     *int              `json:"duration_seconds"`
	Ab55Branches            map[string]string `json:"branches"`
	Ab55Fallback            string            `json:"fallback"`
	Ab55DisplaySteps        []int             `json:"display_steps"`
	Ab55Priority            map[string]string `json:"priority"`
	Ab55PhysicalSupport     string            `json:"physical_support"`
	Ab55Acknowledgement     string            `json:"acknowledgement"`
	Ab55Acceptance          string            `json:"acceptance"`
	Ab55Persistence         string            `json:"persistence"`
	Ab55DurationEnforcement string            `json:"duration_enforcement"`
}

type Ab55Catalog struct {
	Ab55MetadataPresence string      `json:"metadata_presence"`
	Ab55Qualification    string      `json:"qualification"`
	Ab55Limitations      []string    `json:"limitations"`
	Ab55Entries          []Ab55Entry `json:"entries"`
}

func Ab55B505Candidates() Ab55Catalog {
	var ab55Catalog Ab55Catalog
	if ab55Error := ab55JSON.Unmarshal([]byte(ab55CandidatesJSON), &ab55Catalog); ab55Error != nil {
		panic(ab55Error)
	}
	return ab55Catalog
}

func Ab55Lookup(ab55Key Ab55Key) (Ab55Entry, bool) {
	for _, ab55Entry := range Ab55B505Candidates().Ab55Entries {
		if ab55KeysEqual(ab55Entry.Ab55Key, ab55Key) {
			return ab55Entry, true
		}
	}
	return Ab55Entry{}, false
}

func (ab55Entry Ab55Entry) Ab55SelectRaw(ab55Raw string) string {
	if ab55Branch, ab55Found := ab55Entry.Ab55Branches[ab55Raw]; ab55Found {
		return ab55Branch
	}
	return ab55Entry.Ab55Fallback
}

func ab55KeysEqual(ab55Left, ab55Right Ab55Key) bool {
	return ab55Reflect.DeepEqual(ab55Left, ab55Right)
}
