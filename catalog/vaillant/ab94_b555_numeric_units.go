package vaillant

import (
	_ "embed"
	"encoding/json"
	"reflect"
)

type Ab94Key struct {
	Ab94Family    string  `json:"family"`
	Ab94Variant   string  `json:"variant"`
	Ab94Primary   uint16  `json:"primary"`
	Ab94Secondary *uint8  `json:"secondary"`
	Ab94Operation *uint8  `json:"operation"`
	Ab94Direction *string `json:"direction"`
	Ab94Field     string  `json:"field"`
	Ab94Opcode    *uint8  `json:"opcode"`
	Ab94Group     *uint8  `json:"group"`
	Ab94Instance  *uint8  `json:"instance"`
	Ab94Register  *uint16 `json:"register"`
}

type Ab94Replacement struct {
	Ab94Raw     uint64    `json:"raw"`
	Ab94Role    string    `json:"role"`
	Ab94Context string    `json:"context"`
	Ab94Grade   Ab05Grade `json:"grade"`
}

type Ab94NumericField struct {
	Ab94Key             Ab94Key           `json:"key"`
	Ab94LogicalPosition *int              `json:"logical_position"`
	Ab94WidthBits       *int              `json:"width_bits"`
	Ab94Signedness      string            `json:"signedness"`
	Ab94ByteOrder       string            `json:"byte_order"`
	Ab94Conversion      string            `json:"conversion"`
	Ab94Divisor         *float64          `json:"divisor"`
	Ab94Unit            *string           `json:"unit"`
	Ab94UnitCode        *uint8            `json:"unit_code"`
	Ab94Default         *float64          `json:"default"`
	Ab94RawMin          *uint64           `json:"raw_min"`
	Ab94RawMax          *uint64           `json:"raw_max"`
	Ab94NumericMin      *float64          `json:"numeric_min"`
	Ab94NumericMax      *float64          `json:"numeric_max"`
	Ab94Enumeration     map[uint64]string `json:"enumeration"`
	Ab94Replacement     []Ab94Replacement `json:"replacement"`
	Ab94Reserved        []uint64          `json:"reserved"`
	Ab94Grade           Ab05Grade         `json:"grade"`
	Ab94ConversionGrade Ab05Grade         `json:"conversion_grade"`
	Ab94UnitGrade       Ab05Grade         `json:"unit_grade"`
	Ab94RawFallback     string            `json:"raw_fallback"`
	Ab94Enforcement     Ab05Grade         `json:"enforcement"`
	Ab94PeerEffects     Ab05Grade         `json:"peer_effects"`
}

type Ab94NumericMetadata struct {
	Ab94RepresentationGrade   Ab05Grade          `json:"representation_grade"`
	Ab94PhysicalQualification Ab05Grade          `json:"physical_qualification"`
	Ab94NativeSupport         Ab05Grade          `json:"native_support"`
	Ab94Freshness             Ab05Grade          `json:"freshness"`
	Ab94Persistence           Ab05Grade          `json:"persistence"`
	Ab94Execution             Ab05Grade          `json:"execution"`
	Ab94Complete              bool               `json:"complete"`
	Ab94SourceContract        string             `json:"source_contract"`
	Ab94Limitations           []string           `json:"limitations"`
	Ab94Entries               []Ab94NumericField `json:"entries"`
}

//go:embed data/ab94_b555_numeric_units.json
var ab94NumericJSON string

func Ab94B555NumericUnits() (Ab94NumericMetadata, error) {
	var ab94Metadata Ab94NumericMetadata
	ab94Error := json.Unmarshal([]byte(ab94NumericJSON), &ab94Metadata)
	return ab94Metadata, ab94Error
}

func ab94Lookup(ab94Entries []Ab94NumericField, ab94Key Ab94Key) (Ab94NumericField, bool) {
	for _, ab94Entry := range ab94Entries {
		if reflect.DeepEqual(ab94Entry.Ab94Key, ab94Key) {
			return ab94Entry, true
		}
	}
	return Ab94NumericField{}, false
}

func Ab94LookupB555NumericUnits(ab94Key Ab94Key) (Ab94NumericField, bool, error) {
	ab94Metadata, ab94Error := Ab94B555NumericUnits()
	if ab94Error != nil {
		return Ab94NumericField{}, false, ab94Error
	}
	ab94Field, ab94Found := ab94Lookup(ab94Metadata.Ab94Entries, ab94Key)
	return ab94Field, ab94Found, nil
}
