package vaillant

import (
	ab99bytes "bytes"
	_ "embed"
	ab99json "encoding/json"
	ab99fmt "fmt"
	ab99io "io"
	ab99reflect "reflect"
)

//go:embed data/ab99_b523_config_mixer.json
var ab99Data []byte

type Ab99Key struct {
	Ab99Family            string  `json:"family"`
	Ab99Variant           string  `json:"variant"`
	Ab99Primary           uint8   `json:"primary"`
	Ab99Secondary         *uint8  `json:"secondary"`
	Ab99Operation         *uint8  `json:"operation"`
	Ab99Direction         string  `json:"direction"`
	Ab99Field             string  `json:"field"`
	Ab99Group             *uint8  `json:"group"`
	Ab99RequestedInstance *uint8  `json:"requested_instance"`
	Ab99Register          *uint16 `json:"register"`
	Ab99CommandRole       *string `json:"command_role"`
}

type Ab99Bit struct {
	Ab99Index    uint8   `json:"index"`
	Ab99Label    *string `json:"label"`
	Ab99Reserved bool    `json:"reserved"`
	Ab99Default  uint8   `json:"default"`
}

type Ab99Candidate struct {
	Ab99Label   string   `json:"label"`
	Ab99Signed  *bool    `json:"signed"`
	Ab99Scale   *float64 `json:"scale"`
	Ab99Minimum *float64 `json:"minimum"`
	Ab99Maximum *float64 `json:"maximum"`
	Ab99Raw     *uint8   `json:"raw"`
	Ab99Value   *float64 `json:"value"`
	Ab99Clamp   *float64 `json:"clamp"`
}

type Ab99Interval struct {
	Ab99Minimum int    `json:"minimum"`
	Ab99Maximum int    `json:"maximum"`
	Ab99Label   string `json:"label"`
}

type Ab99Entry struct {
	Ab99Signed                 *bool            `json:"signed"`
	Ab99Key                    Ab99Key          `json:"key"`
	Ab99Grade                  string           `json:"grade"`
	Ab99PhysicalAdaptation     string           `json:"physical_adaptation"`
	Ab99Support                string           `json:"support"`
	Ab99ExecutedMovement       string           `json:"executed_movement"`
	Ab99SelectedInterpretation string           `json:"selected_interpretation"`
	Ab99UnknownMetadata        string           `json:"unknown_metadata"`
	Ab99PreserveRaw            bool             `json:"preserve_raw"`
	Ab99Position               *uint8           `json:"position"`
	Ab99WidthBits              *uint8           `json:"width_bits"`
	Ab99RequestExtent          *uint8           `json:"request_extent"`
	Ab99PositiveExtent         *uint8           `json:"positive_extent"`
	Ab99Default                *float64         `json:"default"`
	Ab99Values                 map[uint8]string `json:"values"`
	Ab99OtherValues            string           `json:"other_values"`
	Ab99Bits                   []Ab99Bit        `json:"bits"`
	Ab99OneHot                 *bool            `json:"one_hot"`
	Ab99Alternatives           []string         `json:"alternatives"`
	Ab99RelevantModes          []uint8          `json:"relevant_modes"`
	Ab99IrrelevantModes        []uint8          `json:"irrelevant_modes"`
	Ab99Candidates             []Ab99Candidate  `json:"candidates"`
	Ab99Intervals              []Ab99Interval   `json:"intervals"`
	Ab99Unit                   *string          `json:"unit"`
	Ab99Notes                  []string         `json:"notes"`
}

func Ab99Entries() ([]Ab99Entry, error) {
	var ab99Entries []Ab99Entry
	ab99Decoder := ab99json.NewDecoder(ab99bytes.NewReader(ab99Data))
	ab99Decoder.DisallowUnknownFields()
	if ab99Error := ab99Decoder.Decode(&ab99Entries); ab99Error != nil {
		return nil, ab99Error
	}
	var ab99Trailing any
	if ab99Error := ab99Decoder.Decode(&ab99Trailing); ab99Error != ab99io.EOF {
		return nil, ab99fmt.Errorf("ab99: trailing data")
	}
	if len(ab99Entries) == 0 {
		return nil, ab99fmt.Errorf("ab99: empty catalog")
	}
	for ab99Index, ab99Entry := range ab99Entries {
		if ab99Entry.Ab99Key.Ab99Family == "" || ab99Entry.Ab99Key.Ab99Variant == "" || ab99Entry.Ab99Key.Ab99Direction == "" || ab99Entry.Ab99Key.Ab99Field == "" || ab99Entry.Ab99Grade != "Hypothesis" || !ab99Entry.Ab99PreserveRaw || ab99Entry.Ab99PhysicalAdaptation != "Unknown" || ab99Entry.Ab99Support != "Unknown" || ab99Entry.Ab99ExecutedMovement != "Unknown" || ab99Entry.Ab99SelectedInterpretation != "Unknown" || ab99Entry.Ab99UnknownMetadata != "Unknown" {
			return nil, ab99fmt.Errorf("ab99: invalid entry %d", ab99Index)
		}
		for _, ab99Previous := range ab99Entries[:ab99Index] {
			if ab99reflect.DeepEqual(ab99Previous.Ab99Key, ab99Entry.Ab99Key) {
				return nil, ab99fmt.Errorf("ab99: duplicate key")
			}
		}
	}
	return ab99Entries, nil
}

func Ab99Lookup(ab99Key Ab99Key) (Ab99Entry, bool, error) {
	ab99Entries, ab99Error := Ab99Entries()
	if ab99Error != nil {
		return Ab99Entry{}, false, ab99Error
	}
	for _, ab99Entry := range ab99Entries {
		if ab99reflect.DeepEqual(ab99Entry.Ab99Key, ab99Key) {
			return ab99Entry, true, nil
		}
	}
	return Ab99Entry{}, false, nil
}
