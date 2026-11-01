package vaillant

import (
	ab95bytes "bytes"
	_ "embed"
	ab95json "encoding/json"
	ab95fmt "fmt"
	ab95io "io"
	ab95reflect "reflect"
)

//go:embed data/ab95_b524_energy_roles.json
var ab95Data []byte

type Ab95Key struct {
	Ab95Family            string  `json:"family"`
	Ab95Variant           string  `json:"variant"`
	Ab95Primary           uint8   `json:"primary"`
	Ab95Secondary         *uint8  `json:"secondary"`
	Ab95Operation         *uint8  `json:"operation"`
	Ab95Direction         string  `json:"direction"`
	Ab95Field             string  `json:"field"`
	Ab95Group             *uint8  `json:"group"`
	Ab95RequestedInstance *uint8  `json:"requested_instance"`
	Ab95Register          *uint16 `json:"register"`
	Ab95CommandRole       string  `json:"command_role"`
}

type Ab95Entry struct {
	Ab95ReadbackInstance      *uint8            `json:"readback_instance"`
	Ab95Key                   Ab95Key           `json:"key"`
	Ab95Grade                 string            `json:"grade"`
	Ab95Codec                 *string           `json:"codec"`
	Ab95WidthBytes            *uint8            `json:"width_bytes"`
	Ab95FieldOffset           *uint8            `json:"field_offset"`
	Ab95ByteOrder             *string           `json:"byte_order"`
	Ab95NativeUnit            *string           `json:"native_unit"`
	Ab95ControlEnumeration    map[string]string `json:"control_enumeration"`
	Ab95ResponseForm          *string           `json:"response_form"`
	Ab95Unknown               []string          `json:"unknown"`
	Ab95RawPreservation       string            `json:"raw_preservation"`
	Ab95MetadataBoundaryGrade string            `json:"metadata_boundary_grade"`
	Ab95MetadataBoundary      string            `json:"metadata_boundary"`
	Ab95PhysicalQualification string            `json:"physical_qualification"`
	Ab95NativeContract        string            `json:"native_contract"`
}

func Ab95Entries() ([]Ab95Entry, error) {
	var ab95Entries []Ab95Entry
	ab95Decoder := ab95json.NewDecoder(ab95bytes.NewReader(ab95Data))
	ab95Decoder.DisallowUnknownFields()
	if ab95Error := ab95Decoder.Decode(&ab95Entries); ab95Error != nil {
		return nil, ab95fmt.Errorf("ab95 metadata: %w", ab95Error)
	}
	var ab95Trailing any
	if ab95Error := ab95Decoder.Decode(&ab95Trailing); ab95Error != ab95io.EOF {
		return nil, ab95fmt.Errorf("ab95 metadata: trailing data")
	}
	if len(ab95Entries) != 5 {
		return nil, ab95fmt.Errorf("ab95 metadata: expected five targets")
	}
	for ab95Index, ab95Entry := range ab95Entries {
		if ab95Entry.Ab95Key.Ab95Family == "" || ab95Entry.Ab95Key.Ab95Variant == "" || ab95Entry.Ab95Key.Ab95Direction == "" || ab95Entry.Ab95Key.Ab95Field == "" || ab95Entry.Ab95Key.Ab95CommandRole == "" || ab95Entry.Ab95Grade != "Hypothesis" || ab95Entry.Ab95MetadataBoundaryGrade != "Proven" || ab95Entry.Ab95PhysicalQualification != "Unknown" || ab95Entry.Ab95NativeContract == "" || ab95Entry.Ab95MetadataBoundary == "" || ab95Entry.Ab95RawPreservation == "" || len(ab95Entry.Ab95Unknown) == 0 {
			return nil, ab95fmt.Errorf("ab95 metadata: invalid entry %d", ab95Index)
		}
		for _, ab95Previous := range ab95Entries[:ab95Index] {
			if ab95reflect.DeepEqual(ab95Entry.Ab95Key, ab95Previous.Ab95Key) {
				return nil, ab95fmt.Errorf("ab95 metadata: duplicate key at %d", ab95Index)
			}
		}
	}
	return ab95Entries, nil
}

func Ab95Lookup(ab95Key Ab95Key) (Ab95Entry, bool, error) {
	ab95Entries, ab95Error := Ab95Entries()
	if ab95Error != nil {
		return Ab95Entry{}, false, ab95Error
	}
	for _, ab95Entry := range ab95Entries {
		if ab95reflect.DeepEqual(ab95Key, ab95Entry.Ab95Key) {
			return ab95Entry, true, nil
		}
	}
	return Ab95Entry{}, false, nil
}
