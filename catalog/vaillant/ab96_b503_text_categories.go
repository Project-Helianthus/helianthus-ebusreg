package vaillant

import (
	ab96bytes "bytes"
	_ "embed"
	ab96json "encoding/json"
	ab96fmt "fmt"
	ab96io "io"
	ab96reflect "reflect"
)

//go:embed data/ab96_b503_text_categories.json
var ab96Data []byte

type Ab96Key struct {
	Ab96Family            string  `json:"family"`
	Ab96Variant           string  `json:"variant"`
	Ab96Primary           uint8   `json:"primary"`
	Ab96Secondary         *uint8  `json:"secondary"`
	Ab96Operation         *uint8  `json:"operation"`
	Ab96Direction         *string `json:"direction"`
	Ab96Field             string  `json:"field"`
	Ab96Group             *uint8  `json:"group"`
	Ab96RequestedInstance *uint8  `json:"requested_instance"`
	Ab96Register          *uint16 `json:"register"`
	Ab96CommandRole       *string `json:"command_role"`
}

type Ab96Entry struct {
	Ab96Key                   Ab96Key           `json:"key"`
	Ab96Grade                 string            `json:"grade"`
	Ab96Position              *int              `json:"position"`
	Ab96WidthCharacters       *int              `json:"width_characters"`
	Ab96WidthBits             *int              `json:"width_bits"`
	Ab96Unsigned              *bool             `json:"unsigned"`
	Ab96Default               *uint8            `json:"default"`
	Ab96Encoding              *string           `json:"encoding"`
	Ab96Padding               *string           `json:"padding"`
	Ab96Extent                *string           `json:"extent"`
	Ab96Follows               *string           `json:"follows"`
	Ab96Labels                map[uint8]*string `json:"labels"`
	Ab96MetadataPresence      string            `json:"metadata_presence"`
	Ab96PhysicalQualification string            `json:"physical_qualification"`
	Ab96RawPreservation       string            `json:"raw_preservation"`
	Ab96Unknown               []string          `json:"unknown"`
	Ab96NativeContract        string            `json:"native_contract"`
}

func Ab96Entries() ([]Ab96Entry, error) {
	var ab96Entries []Ab96Entry
	ab96Decoder := ab96json.NewDecoder(ab96bytes.NewReader(ab96Data))
	ab96Decoder.DisallowUnknownFields()
	if ab96Error := ab96Decoder.Decode(&ab96Entries); ab96Error != nil {
		return nil, ab96Error
	}
	if ab96Error := ab96Decoder.Decode(new(any)); ab96Error != ab96io.EOF {
		return nil, ab96fmt.Errorf("ab96: trailing embedded data: %v", ab96Error)
	}
	if len(ab96Entries) == 0 {
		return nil, ab96fmt.Errorf("ab96: empty embedded catalog")
	}
	for ab96Index, ab96Entry := range ab96Entries {
		if ab96Entry.Ab96Key.Ab96Family == "" || ab96Entry.Ab96Key.Ab96Variant == "" || ab96Entry.Ab96Key.Ab96Field == "" ||
			(ab96Entry.Ab96Grade != "Hypothesis" && ab96Entry.Ab96Grade != "Proven" && ab96Entry.Ab96Grade != "Unknown") ||
			ab96Entry.Ab96MetadataPresence != "Proven" || ab96Entry.Ab96PhysicalQualification != "Unknown" ||
			ab96Entry.Ab96NativeContract == "" || ab96Entry.Ab96RawPreservation == "" || len(ab96Entry.Ab96Unknown) == 0 {
			return nil, ab96fmt.Errorf("ab96: invalid embedded entry %d", ab96Index)
		}
		if (ab96Entry.Ab96WidthCharacters != nil && *ab96Entry.Ab96WidthCharacters <= 0) ||
			(ab96Entry.Ab96WidthBits != nil && *ab96Entry.Ab96WidthBits <= 0) ||
			(ab96Entry.Ab96Position != nil && *ab96Entry.Ab96Position < 0) {
			return nil, ab96fmt.Errorf("ab96: invalid embedded width or position")
		}
		for _, ab96Previous := range ab96Entries[:ab96Index] {
			if ab96reflect.DeepEqual(ab96Previous.Ab96Key, ab96Entry.Ab96Key) {
				return nil, ab96fmt.Errorf("ab96: duplicate embedded key")
			}
		}
	}
	return ab96Entries, nil
}

func Ab96Lookup(ab96Key Ab96Key) (Ab96Entry, bool, error) {
	ab96Entries, ab96Error := Ab96Entries()
	if ab96Error != nil {
		return Ab96Entry{}, false, ab96Error
	}
	for _, ab96Entry := range ab96Entries {
		if ab96reflect.DeepEqual(ab96Key, ab96Entry.Ab96Key) {
			return ab96Entry, true, nil
		}
	}
	return Ab96Entry{}, false, nil
}
