package vaillant

import (
	_ "embed"
	"encoding/json"
	"reflect"
)

//go:embed data/ab91_b51cd_partial_labels.json
var ab91PartialLabelsJSON string

type Ab91PartialKey struct {
	Ab91Family    string  `json:"family"`
	Ab91Variant   string  `json:"variant"`
	Ab91Label     string  `json:"label"`
	Ab91Primary   uint8   `json:"primary"`
	Ab91Secondary *uint16 `json:"secondary"`
	Ab91Operation *uint16 `json:"operation"`
	Ab91Direction string  `json:"direction"`
	Ab91Field     string  `json:"field"`
	Ab91Opcode    *uint16 `json:"opcode"`
	Ab91Group     *uint16 `json:"group"`
	Ab91Instance  *uint16 `json:"instance"`
	Ab91Register  *uint16 `json:"register"`
}

type Ab91PartialField struct {
	Ab91LogicalByte *int              `json:"logical_byte"`
	Ab91LogicalBit  *int              `json:"logical_bit"`
	Ab91WidthBits   *int              `json:"width_bits"`
	Ab91Constant    *uint8            `json:"constant"`
	Ab91Signedness  string            `json:"signedness"`
	Ab91ByteOrder   string            `json:"byte_order"`
	Ab91Conversion  string            `json:"conversion"`
	Ab91Unit        *string           `json:"unit"`
	Ab91Default     *int64            `json:"default"`
	Ab91Bounds      []int64           `json:"bounds"`
	Ab91Enumeration map[string]string `json:"enumeration"`
	Ab91Replacement []int64           `json:"replacement"`
	Ab91Reserved    []int64           `json:"reserved"`
	Ab91Grade       string            `json:"grade"`
	Ab91RawFallback string            `json:"raw_fallback"`
}

type Ab91PartialLabel struct {
	Ab91Key         Ab91PartialKey   `json:"key"`
	Ab91Field       Ab91PartialField `json:"metadata"`
	Ab91BodyWidth   *int             `json:"body_width"`
	Ab91Grade       string           `json:"grade"`
	Ab91RawFallback string           `json:"raw_fallback"`
}

type Ab91PartialLabels struct {
	Ab91Representation        string             `json:"software_representation"`
	Ab91PhysicalQualification string             `json:"physical_qualification"`
	Ab91InstalledSupport      string             `json:"installed_support"`
	Ab91Equivalence           string             `json:"ebus_pebus_equivalence"`
	Ab91BodyLayout            string             `json:"body_layout"`
	Ab91Codec                 string             `json:"codec"`
	Ab91Gating                string             `json:"gating"`
	Ab91Confirmation          string             `json:"confirmation"`
	Ab91Effect                string             `json:"effect"`
	Ab91PhysicalFraming       string             `json:"physical_framing"`
	Ab91Controls              string             `json:"uninterpreted_controls"`
	Ab91Completeness          string             `json:"family_completeness"`
	Ab91Limitations           []string           `json:"limitations"`
	Ab91Entries               []Ab91PartialLabel `json:"entries"`
}

func Ab91B51CDPartialLabels() (Ab91PartialLabels, error) {
	var ab91Catalog Ab91PartialLabels
	ab91Error := json.Unmarshal([]byte(ab91PartialLabelsJSON), &ab91Catalog)
	return ab91Catalog, ab91Error
}

func ab91FindPartialLabel(ab91Entries []Ab91PartialLabel, ab91Key Ab91PartialKey) (Ab91PartialLabel, bool) {
	for _, ab91Entry := range ab91Entries {
		if reflect.DeepEqual(ab91Entry.Ab91Key, ab91Key) {
			return ab91Entry, true
		}
	}
	return Ab91PartialLabel{}, false
}

func Ab91LookupB51CDPartialLabel(ab91Key Ab91PartialKey) (Ab91PartialLabel, bool, error) {
	ab91Catalog, ab91Error := Ab91B51CDPartialLabels()
	if ab91Error != nil {
		return Ab91PartialLabel{}, false, ab91Error
	}
	ab91Entry, ab91Found := ab91FindPartialLabel(ab91Catalog.Ab91Entries, ab91Key)
	return ab91Entry, ab91Found, nil
}
