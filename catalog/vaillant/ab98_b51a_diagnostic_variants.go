package vaillant

import (
	ab98bytes "bytes"
	_ "embed"
	ab98json "encoding/json"
	ab98fmt "fmt"
	ab98io "io"
	ab98reflect "reflect"
)

//go:embed data/ab98_b51a_diagnostic_variants.json
var ab98Data []byte

type Ab98Key struct {
	Ab98Family            string  `json:"family"`
	Ab98Variant           string  `json:"variant"`
	Ab98Primary           uint8   `json:"primary"`
	Ab98Secondary         *uint8  `json:"secondary"`
	Ab98Operation         *uint8  `json:"operation"`
	Ab98Direction         string  `json:"direction"`
	Ab98Field             string  `json:"field"`
	Ab98Group             *uint8  `json:"group"`
	Ab98RequestedInstance *uint8  `json:"requested_instance"`
	Ab98Register          *uint16 `json:"register"`
	Ab98CommandRole       *string `json:"command_role"`
}

type Ab98Condition struct {
	Ab98Field string `json:"field"`
	Ab98Bit   uint8  `json:"bit"`
	Ab98Set   bool   `json:"set"`
}

type Ab98Bit struct {
	Ab98Bit      uint8    `json:"bit"`
	Ab98Labels   []string `json:"labels"`
	Ab98Reserved bool     `json:"reserved"`
	Ab98Default  bool     `json:"default"`
}

type Ab98Entry struct {
	Ab98Key                   Ab98Key           `json:"key"`
	Ab98Grade                 string            `json:"grade"`
	Ab98MetadataPresence      string            `json:"metadata_presence"`
	Ab98PhysicalQualification string            `json:"physical_qualification"`
	Ab98Invocation            string            `json:"invocation"`
	Ab98NativeContract        string            `json:"native_contract"`
	Ab98Position              *uint16           `json:"position"`
	Ab98WidthBits             *uint16           `json:"width_bits"`
	Ab98ByteOrder             *string           `json:"byte_order"`
	Ab98Domain                []uint64          `json:"domain"`
	Ab98Default               *uint64           `json:"default"`
	Ab98Values                map[uint64]string `json:"values"`
	Ab98Bits                  []Ab98Bit         `json:"bits"`
	Ab98Conditions            []Ab98Condition   `json:"conditions"`
	Ab98Alternatives          []Ab98Alternative `json:"alternatives"`
	Ab98Interpretations       []string          `json:"interpretations"`
	Ab98Unknown               []string          `json:"unknown"`
	Ab98Exclusive             *bool             `json:"exclusive"`
	Ab98Count                 *uint16           `json:"count"`
	Ab98LogicalFootprint      *uint16           `json:"logical_footprint"`
	Ab98Tail                  []uint16          `json:"tail"`
}

type Ab98Alternative struct {
	Ab98Label           string          `json:"label"`
	Ab98Conditions      []Ab98Condition `json:"conditions"`
	Ab98Fields          []Ab98Entry     `json:"fields"`
	Ab98Interpretations []string        `json:"interpretations"`
}

func Ab98Entries() ([]Ab98Entry, error) {
	var ab98Entries []Ab98Entry
	ab98Decoder := ab98json.NewDecoder(ab98bytes.NewReader(ab98Data))
	ab98Decoder.DisallowUnknownFields()
	if ab98Err := ab98Decoder.Decode(&ab98Entries); ab98Err != nil {
		return nil, ab98fmt.Errorf("ab98 catalog: %w", ab98Err)
	}
	if ab98Err := ab98Decoder.Decode(new(any)); ab98Err != ab98io.EOF {
		return nil, ab98fmt.Errorf("ab98 catalog: trailing data")
	}
	if len(ab98Entries) == 0 {
		return nil, ab98fmt.Errorf("ab98 catalog: empty entries")
	}
	for ab98Index, ab98Entry := range ab98Entries {
		if ab98Err := ab98Validate(ab98Entry); ab98Err != nil {
			return nil, ab98Err
		}
		for ab98Previous := 0; ab98Previous < ab98Index; ab98Previous++ {
			if ab98reflect.DeepEqual(ab98Entries[ab98Previous].Ab98Key, ab98Entry.Ab98Key) {
				return nil, ab98fmt.Errorf("ab98 catalog: duplicate key")
			}
		}
	}
	return ab98Entries, nil
}

func ab98Validate(ab98Entry Ab98Entry) error {
	if ab98Entry.Ab98Key.Ab98Family == "" || ab98Entry.Ab98Key.Ab98Variant == "" || ab98Entry.Ab98Key.Ab98Field == "" || (ab98Entry.Ab98Key.Ab98Direction != "request" && ab98Entry.Ab98Key.Ab98Direction != "positive") {
		return ab98fmt.Errorf("ab98 catalog: incomplete key")
	}
	if ab98Entry.Ab98Grade != "Hypothesis" || ab98Entry.Ab98MetadataPresence != "Proven" || ab98Entry.Ab98PhysicalQualification != "Unknown" || ab98Entry.Ab98Invocation != "Unknown" || ab98Entry.Ab98NativeContract == "" {
		return ab98fmt.Errorf("ab98 catalog: invalid evidence boundary")
	}
	if len(ab98Entry.Ab98Domain) != 0 && (len(ab98Entry.Ab98Domain) != 2 || ab98Entry.Ab98Domain[0] > ab98Entry.Ab98Domain[1]) {
		return ab98fmt.Errorf("ab98 catalog: invalid domain")
	}
	if ab98Entry.Ab98Default != nil && len(ab98Entry.Ab98Domain) == 2 && (*ab98Entry.Ab98Default < ab98Entry.Ab98Domain[0] || *ab98Entry.Ab98Default > ab98Entry.Ab98Domain[1]) {
		return ab98fmt.Errorf("ab98 catalog: default outside domain")
	}
	if ab98Entry.Ab98WidthBits != nil && (*ab98Entry.Ab98WidthBits == 0 || *ab98Entry.Ab98WidthBits > 64) {
		return ab98fmt.Errorf("ab98 catalog: invalid width")
	}
	if ab98Entry.Ab98ByteOrder != nil && *ab98Entry.Ab98ByteOrder != "little-endian" {
		return ab98fmt.Errorf("ab98 catalog: invalid byte order")
	}
	ab98SeenBits := map[uint8]bool{}
	for _, ab98Bit := range ab98Entry.Ab98Bits {
		if ab98Bit.Ab98Bit > 7 || ab98SeenBits[ab98Bit.Ab98Bit] {
			return ab98fmt.Errorf("ab98 catalog: invalid bit")
		}
		ab98SeenBits[ab98Bit.Ab98Bit] = true
	}
	if ab98Err := ab98ValidateConditions(ab98Entry.Ab98Conditions); ab98Err != nil {
		return ab98Err
	}
	for _, ab98Alternative := range ab98Entry.Ab98Alternatives {
		if ab98Err := ab98ValidateConditions(ab98Alternative.Ab98Conditions); ab98Err != nil {
			return ab98Err
		}
	}
	for _, ab98Alternative := range ab98Entry.Ab98Alternatives {
		if ab98Alternative.Ab98Label == "" || len(ab98Alternative.Ab98Fields) == 0 {
			return ab98fmt.Errorf("ab98 catalog: invalid alternative")
		}
		for _, ab98Field := range ab98Alternative.Ab98Fields {
			if ab98Err := ab98Validate(ab98Field); ab98Err != nil {
				return ab98Err
			}
		}
	}
	return nil
}

func ab98ValidateConditions(ab98Conditions []Ab98Condition) error {
	for _, ab98Condition := range ab98Conditions {
		if ab98Condition.Ab98Bit > 7 || ab98Condition.Ab98Field == "" {
			return ab98fmt.Errorf("ab98 catalog: invalid condition")
		}
	}
	return nil
}

func Ab98Lookup(ab98Key Ab98Key) (Ab98Entry, bool, error) {
	ab98Entries, ab98Err := Ab98Entries()
	if ab98Err != nil {
		return Ab98Entry{}, false, ab98Err
	}
	for _, ab98Entry := range ab98Entries {
		if ab98reflect.DeepEqual(ab98Key, ab98Entry.Ab98Key) {
			return ab98Entry, true, nil
		}
	}
	return Ab98Entry{}, false, nil
}
