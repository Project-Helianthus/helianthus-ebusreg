package vaillant

import (
	ab72bytes "bytes"
	_ "embed"
	ab72json "encoding/json"
	ab72fmt "fmt"
	ab72io "io"
	ab72reflect "reflect"
)

//go:embed data/ab72_b516_statistics_directions.json
var ab72Data []byte

type Ab72Key struct {
	Ab72Family    string `json:"family"`
	Ab72Variant   string `json:"variant"`
	Ab72Operation *uint8 `json:"operation"`
	Ab72Direction string `json:"direction"`
	Ab72Selector  *uint8 `json:"selector"`
	Ab72Field     string `json:"field"`
}

type Ab72Entry struct {
	Ab72Key              Ab72Key           `json:"key"`
	Ab72Grade            string            `json:"grade"`
	Ab72SoftwareContract string            `json:"software_contract"`
	Ab72PhysicalSupport  string            `json:"physical_support"`
	Ab72Presence         string            `json:"presence"`
	Ab72Freshness        string            `json:"freshness"`
	Ab72Execution        string            `json:"execution"`
	Ab72Effects          string            `json:"effects"`
	Ab72Persistence      string            `json:"persistence"`
	Ab72Delivery         string            `json:"delivery"`
	Ab72NativeContract   string            `json:"native_contract"`
	Ab72Length           *uint8            `json:"length"`
	Ab72Position         *uint8            `json:"position"`
	Ab72Width            *uint8            `json:"width"`
	Ab72Encoding         *string           `json:"encoding"`
	Ab72Unit             *string           `json:"unit"`
	Ab72Bits             map[string]uint8  `json:"bits"`
	Ab72Values           map[string]string `json:"values"`
	Ab72Order            []string          `json:"order"`
	Ab72Rules            []string          `json:"rules"`
	Ab72Conflicts        []string          `json:"conflicts"`
	Ab72Unknown          []string          `json:"unknown"`
}

func ab72Parse(ab72Raw []byte) ([]Ab72Entry, error) {
	var ab72Entries []Ab72Entry
	ab72Decoder := ab72json.NewDecoder(ab72bytes.NewReader(ab72Raw))
	ab72Decoder.DisallowUnknownFields()
	if ab72Err := ab72Decoder.Decode(&ab72Entries); ab72Err != nil {
		return nil, ab72fmt.Errorf("ab72 catalog: %w", ab72Err)
	}
	if ab72Err := ab72Decoder.Decode(new(any)); ab72Err != ab72io.EOF {
		return nil, ab72fmt.Errorf("ab72 catalog: trailing input")
	}
	if len(ab72Entries) == 0 {
		return nil, ab72fmt.Errorf("ab72 catalog: empty")
	}
	for ab72Index, ab72Entry := range ab72Entries {
		if ab72Entry.Ab72Key.Ab72Family == "" || ab72Entry.Ab72Key.Ab72Variant == "" || ab72Entry.Ab72Key.Ab72Field == "" || ab72Entry.Ab72Key.Ab72Operation == nil || (ab72Entry.Ab72Key.Ab72Direction != "request" && ab72Entry.Ab72Key.Ab72Direction != "positive") {
			return nil, ab72fmt.Errorf("ab72 catalog: incomplete key")
		}
		if ab72Entry.Ab72Grade != "Hypothesis" || ab72Entry.Ab72SoftwareContract != "Proven" || ab72Entry.Ab72NativeContract != "https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/protocols/ebus-services/ebus-overview.md" {
			return nil, ab72fmt.Errorf("ab72 catalog: invalid evidence")
		}
		for _, ab72Grade := range []string{ab72Entry.Ab72PhysicalSupport, ab72Entry.Ab72Presence, ab72Entry.Ab72Freshness, ab72Entry.Ab72Execution, ab72Entry.Ab72Effects, ab72Entry.Ab72Persistence, ab72Entry.Ab72Delivery} {
			if ab72Grade != "Unknown" {
				return nil, ab72fmt.Errorf("ab72 catalog: invalid qualification")
			}
		}
		for _, ab72Bit := range ab72Entry.Ab72Bits {
			if ab72Bit > 7 {
				return nil, ab72fmt.Errorf("ab72 catalog: invalid bit")
			}
		}
		if ab72Entry.Ab72Width != nil && (*ab72Entry.Ab72Width == 0 || *ab72Entry.Ab72Width > 32) {
			return nil, ab72fmt.Errorf("ab72 catalog: invalid width")
		}
		for ab72Previous := 0; ab72Previous < ab72Index; ab72Previous++ {
			if ab72reflect.DeepEqual(ab72Entries[ab72Previous].Ab72Key, ab72Entry.Ab72Key) {
				return nil, ab72fmt.Errorf("ab72 catalog: duplicate key")
			}
		}
	}
	return ab72Entries, nil
}

func Ab72Entries() ([]Ab72Entry, error) { return ab72Parse(ab72Data) }

func Ab72Lookup(ab72Key Ab72Key) (Ab72Entry, bool, error) {
	ab72Entries, ab72Err := Ab72Entries()
	if ab72Err != nil {
		return Ab72Entry{}, false, ab72Err
	}
	for _, ab72Entry := range ab72Entries {
		if ab72reflect.DeepEqual(ab72Key, ab72Entry.Ab72Key) {
			return ab72Entry, true, nil
		}
	}
	return Ab72Entry{}, false, nil
}
