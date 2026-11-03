package vaillant

import (
	ab14Bytes "bytes"
	_ "embed"
	ab14JSON "encoding/json"
	ab14Fmt "fmt"
	ab14IO "io"
	ab14Reflect "reflect"
)

//go:embed data/ab14_b509_op59_quantities.json
var ab14Data string

type Ab14Key struct {
	Ab14Family            string  `json:"family"`
	Ab14Variant           string  `json:"variant"`
	Ab14Primary           uint8   `json:"primary"`
	Ab14Secondary         *uint8  `json:"secondary"`
	Ab14Operation         *uint8  `json:"operation"`
	Ab14Direction         *string `json:"direction"`
	Ab14Field             string  `json:"field"`
	Ab14Group             *uint8  `json:"group"`
	Ab14RequestedInstance *uint8  `json:"requested_instance"`
	Ab14Register          *uint16 `json:"register"`
	Ab14CommandRole       *string `json:"command_role"`
}

type Ab14Codec struct {
	Ab14Representation string   `json:"representation"`
	Ab14Width          *uint8   `json:"width_bits"`
	Ab14Divisor        *float64 `json:"divisor"`
	Ab14Unit           *string  `json:"unit"`
}

type Ab14Entry struct {
	Ab14Key                   Ab14Key    `json:"key"`
	Ab14Grade                 string     `json:"grade"`
	Ab14MetadataPresence      string     `json:"metadata_presence"`
	Ab14PhysicalQualification string     `json:"physical_qualification"`
	Ab14Codec                 *Ab14Codec `json:"codec"`
	Ab14Position              *uint16    `json:"position"`
	Ab14Label                 *string    `json:"label"`
	Ab14PreserveRaw           bool       `json:"preserve_raw"`
	Ab14Unknown               []string   `json:"unknown"`
}

func Ab14Entries() ([]Ab14Entry, error) {
	var ab14Entries []Ab14Entry
	ab14Decoder := ab14JSON.NewDecoder(ab14Bytes.NewBufferString(ab14Data))
	ab14Decoder.DisallowUnknownFields()
	if ab14Error := ab14Decoder.Decode(&ab14Entries); ab14Error != nil {
		return nil, ab14Error
	}
	if ab14Error := ab14Decoder.Decode(new(any)); ab14Error != ab14IO.EOF {
		return nil, ab14Fmt.Errorf("ab14: trailing embedded data")
	}
	if len(ab14Entries) != 5 {
		return nil, ab14Fmt.Errorf("ab14: expected five finite entries")
	}
	for ab14Index, ab14Entry := range ab14Entries {
		if ab14Entry.Ab14Key.Ab14Family == "" || ab14Entry.Ab14Key.Ab14Variant == "" || ab14Entry.Ab14Key.Ab14Primary == 0 || ab14Entry.Ab14Key.Ab14Field == "" || ab14Entry.Ab14Grade != "Hypothesis" || ab14Entry.Ab14MetadataPresence != "Proven" || ab14Entry.Ab14PhysicalQualification != "Unknown" || !ab14Entry.Ab14PreserveRaw || len(ab14Entry.Ab14Unknown) == 0 {
			return nil, ab14Fmt.Errorf("ab14: invalid entry %d", ab14Index)
		}
		if ab14Entry.Ab14Codec == nil || (ab14Entry.Ab14Codec.Ab14Representation != "raw" && ab14Entry.Ab14Codec.Ab14Representation != "signedLE16" && ab14Entry.Ab14Codec.Ab14Representation != "unsignedLE16") || (ab14Entry.Ab14Codec.Ab14Divisor != nil && *ab14Entry.Ab14Codec.Ab14Divisor <= 0) {
			return nil, ab14Fmt.Errorf("ab14: invalid codec %d", ab14Index)
		}
		ab14Codec := ab14Entry.Ab14Codec
		if ab14Codec.Ab14Representation == "raw" {
			if ab14Codec.Ab14Width != nil || ab14Codec.Ab14Divisor != nil || ab14Codec.Ab14Unit != nil {
				return nil, ab14Fmt.Errorf("ab14: raw limit has unprovided metadata")
			}
		} else if ab14Codec.Ab14Width == nil || *ab14Codec.Ab14Width != 16 || ab14Codec.Ab14Divisor == nil || ab14Codec.Ab14Unit == nil || *ab14Codec.Ab14Unit == "" {
			return nil, ab14Fmt.Errorf("ab14: incomplete quantity codec")
		}
		for _, ab14Prior := range ab14Entries[:ab14Index] {
			if ab14Reflect.DeepEqual(ab14Prior.Ab14Key, ab14Entry.Ab14Key) {
				return nil, ab14Fmt.Errorf("ab14: duplicate key")
			}
		}
	}
	return ab14Entries, nil
}

func Ab14Lookup(ab14Key Ab14Key) (Ab14Entry, bool, error) {
	ab14Entries, ab14Error := Ab14Entries()
	if ab14Error != nil {
		return Ab14Entry{}, false, ab14Error
	}
	for _, ab14Entry := range ab14Entries {
		if ab14Reflect.DeepEqual(ab14Entry.Ab14Key, ab14Key) {
			return ab14Entry, true, nil
		}
	}
	return Ab14Entry{}, false, nil
}
