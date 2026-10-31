package vaillant

import (
	ab70bytes "bytes"
	_ "embed"
	ab70json "encoding/json"
	ab70fmt "fmt"
	ab70io "io"
	ab70reflect "reflect"
)

type Ab70Key struct {
	Ab70Family    string   `json:"family"`
	Ab70Variant   string   `json:"variant"`
	Ab70Primary   *byte    `json:"primary"`
	Ab70Secondary *byte    `json:"secondary"`
	Ab70Operation *byte    `json:"operation"`
	Ab70Direction string   `json:"direction"`
	Ab70Selector  []uint16 `json:"selector"`
}

type Ab70Entry struct {
	Ab70Key              Ab70Key           `json:"key"`
	Ab70Quantity         string            `json:"quantity"`
	Ab70CodecDeclaration string            `json:"codec_declaration"`
	Ab70Unit             *string           `json:"unit"`
	Ab70WidthBits        *int              `json:"width_bits"`
	Ab70Signedness       *string           `json:"signedness"`
	Ab70ByteOrder        *string           `json:"byte_order"`
	Ab70DeclaredScale    *int              `json:"declared_scale"`
	Ab70Converter        *string           `json:"converter"`
	Ab70Mapping          map[string]string `json:"mapping"`
	Ab70Sentinel         *int              `json:"sentinel"`
	Ab70Grade            string            `json:"grade"`
	Ab70SoftwareContract string            `json:"software_contract"`
	Ab70Qualification    string            `json:"qualification"`
	Ab70Presence         string            `json:"presence"`
	Ab70Freshness        string            `json:"freshness"`
	Ab70Execution        string            `json:"execution"`
	Ab70PhysicalEffects  string            `json:"physical_effects"`
	Ab70Persistence      string            `json:"persistence"`
	Ab70Delivery         string            `json:"delivery"`
	Ab70Notes            []string          `json:"notes"`
}

//go:embed data/ab70_b509_heatpump_codec_dictionary.json
var ab70JSON []byte

func ab70Parse(ab70Data []byte) ([]Ab70Entry, error) {
	var ab70Entries []Ab70Entry
	ab70Decoder := ab70json.NewDecoder(ab70bytes.NewReader(ab70Data))
	ab70Decoder.DisallowUnknownFields()
	if ab70Error := ab70Decoder.Decode(&ab70Entries); ab70Error != nil {
		return nil, ab70Error
	}
	if ab70Error := ab70Decoder.Decode(new(any)); ab70Error != ab70io.EOF {
		return nil, ab70fmt.Errorf("ab70: trailing data")
	}
	if len(ab70Entries) == 0 {
		return nil, ab70fmt.Errorf("ab70: empty dictionary")
	}
	ab70Seen := make(map[string]bool)
	for _, ab70Entry := range ab70Entries {
		ab70Key := ab70Entry.Ab70Key
		if ab70Key.Ab70Family == "" || ab70Key.Ab70Variant == "" || ab70Key.Ab70Primary == nil || ab70Key.Ab70Secondary == nil || ab70Key.Ab70Operation == nil || ab70Key.Ab70Direction == "" || len(ab70Key.Ab70Selector) != 2 || ab70Key.Ab70Selector[0] > 255 || ab70Key.Ab70Selector[1] > 255 || ab70Entry.Ab70Quantity == "" || ab70Entry.Ab70CodecDeclaration == "" || ab70Entry.Ab70Grade != "Hypothesis" || ab70Entry.Ab70SoftwareContract != "Proven" || ab70Entry.Ab70Qualification != "Unknown" || ab70Entry.Ab70Presence != "Unknown" || ab70Entry.Ab70Freshness != "Unknown" || ab70Entry.Ab70Execution != "Unknown" || ab70Entry.Ab70PhysicalEffects != "Unknown" || ab70Entry.Ab70Persistence != "Unknown" || ab70Entry.Ab70Delivery != "Unknown" {
			return nil, ab70fmt.Errorf("ab70: incomplete or invalid entry")
		}
		ab70Encoded, ab70Error := ab70json.Marshal(ab70Key)
		if ab70Error != nil {
			return nil, ab70Error
		}
		if ab70Seen[string(ab70Encoded)] {
			return nil, ab70fmt.Errorf("ab70: duplicate key")
		}
		ab70Seen[string(ab70Encoded)] = true
	}
	return ab70Entries, nil
}

func Ab70Entries() ([]Ab70Entry, error) { return ab70Parse(ab70JSON) }

func ab70Find(ab70Entries []Ab70Entry, ab70Key Ab70Key) (Ab70Entry, bool, error) {
	for _, ab70Entry := range ab70Entries {
		if ab70reflect.DeepEqual(ab70Entry.Ab70Key, ab70Key) {
			ab70Data, ab70Error := ab70json.Marshal(ab70Entry)
			if ab70Error != nil {
				return Ab70Entry{}, false, ab70Error
			}
			var ab70Copy Ab70Entry
			if ab70Error = ab70json.Unmarshal(ab70Data, &ab70Copy); ab70Error != nil {
				return Ab70Entry{}, false, ab70Error
			}
			return ab70Copy, true, nil
		}
	}
	return Ab70Entry{}, false, nil
}

func Ab70Lookup(ab70Key Ab70Key) (Ab70Entry, bool, error) {
	ab70Entries, ab70Error := Ab70Entries()
	if ab70Error != nil {
		return Ab70Entry{}, false, ab70Error
	}
	return ab70Find(ab70Entries, ab70Key)
}
