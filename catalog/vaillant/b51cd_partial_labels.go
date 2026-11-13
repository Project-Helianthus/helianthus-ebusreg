package vaillant

import (
	_ "embed"
	"encoding/json"
	"reflect"
)

//go:embed data/ab91_b51cd_partial_labels.json
var b51CDPartialLabelPartialLabelsJSON string

type B51CDPartialKey struct {
	B51CDPartialLabelFamily    string  `json:"family"`
	B51CDPartialLabelVariant   string  `json:"variant"`
	B51CDPartialLabelLabel     string  `json:"label"`
	B51CDPartialLabelPrimary   uint8   `json:"primary"`
	B51CDPartialLabelSecondary *uint16 `json:"secondary"`
	B51CDPartialLabelOperation *uint16 `json:"operation"`
	B51CDPartialLabelDirection string  `json:"direction"`
	B51CDPartialLabelField     string  `json:"field"`
	B51CDPartialLabelOpcode    *uint16 `json:"opcode"`
	B51CDPartialLabelGroup     *uint16 `json:"group"`
	B51CDPartialLabelInstance  *uint16 `json:"instance"`
	B51CDPartialLabelRegister  *uint16 `json:"register"`
}

type B51CDPartialField struct {
	B51CDPartialLabelLogicalByte *int              `json:"logical_byte"`
	B51CDPartialLabelLogicalBit  *int              `json:"logical_bit"`
	B51CDPartialLabelWidthBits   *int              `json:"width_bits"`
	B51CDPartialLabelConstant    *uint8            `json:"constant"`
	B51CDPartialLabelSignedness  string            `json:"signedness"`
	B51CDPartialLabelByteOrder   string            `json:"byte_order"`
	B51CDPartialLabelConversion  string            `json:"conversion"`
	B51CDPartialLabelUnit        *string           `json:"unit"`
	B51CDPartialLabelDefault     *int64            `json:"default"`
	B51CDPartialLabelBounds      []int64           `json:"bounds"`
	B51CDPartialLabelEnumeration map[string]string `json:"enumeration"`
	B51CDPartialLabelReplacement []int64           `json:"replacement"`
	B51CDPartialLabelReserved    []int64           `json:"reserved"`
	B51CDPartialLabelGrade       string            `json:"grade"`
	B51CDPartialLabelRawFallback string            `json:"raw_fallback"`
}

type B51CDPartialLabel struct {
	B51CDPartialLabelKey         B51CDPartialKey   `json:"key"`
	B51CDPartialLabelField       B51CDPartialField `json:"metadata"`
	B51CDPartialLabelBodyWidth   *int              `json:"body_width"`
	B51CDPartialLabelGrade       string            `json:"grade"`
	B51CDPartialLabelRawFallback string            `json:"raw_fallback"`
}

type B51CDPartialLabelPartialLabels struct {
	B51CDPartialLabelRepresentation        string              `json:"software_representation"`
	B51CDPartialLabelPhysicalQualification string              `json:"physical_qualification"`
	B51CDPartialLabelInstalledSupport      string              `json:"installed_support"`
	B51CDPartialLabelEquivalence           string              `json:"ebus_pebus_equivalence"`
	B51CDPartialLabelBodyLayout            string              `json:"body_layout"`
	B51CDPartialLabelCodec                 string              `json:"codec"`
	B51CDPartialLabelGating                string              `json:"gating"`
	B51CDPartialLabelConfirmation          string              `json:"confirmation"`
	B51CDPartialLabelEffect                string              `json:"effect"`
	B51CDPartialLabelPhysicalFraming       string              `json:"physical_framing"`
	B51CDPartialLabelControls              string              `json:"uninterpreted_controls"`
	B51CDPartialLabelCompleteness          string              `json:"family_completeness"`
	B51CDPartialLabelLimitations           []string            `json:"limitations"`
	B51CDPartialLabelEntries               []B51CDPartialLabel `json:"entries"`
}

func B51CDPartialLabels() (B51CDPartialLabelPartialLabels, error) {
	var b51CDPartialLabelCatalog B51CDPartialLabelPartialLabels
	b51CDPartialLabelError := json.Unmarshal([]byte(b51CDPartialLabelPartialLabelsJSON), &b51CDPartialLabelCatalog)
	return b51CDPartialLabelCatalog, b51CDPartialLabelError
}

func b51CDPartialLabelFindPartialLabel(b51CDPartialLabelEntries []B51CDPartialLabel, b51CDPartialLabelKey B51CDPartialKey) (B51CDPartialLabel, bool) {
	for _, b51CDPartialLabelEntry := range b51CDPartialLabelEntries {
		if reflect.DeepEqual(b51CDPartialLabelEntry.B51CDPartialLabelKey, b51CDPartialLabelKey) {
			return b51CDPartialLabelEntry, true
		}
	}
	return B51CDPartialLabel{}, false
}

func B51CDLookupPartialLabel(b51CDPartialLabelKey B51CDPartialKey) (B51CDPartialLabel, bool, error) {
	b51CDPartialLabelCatalog, b51CDPartialLabelError := B51CDPartialLabels()
	if b51CDPartialLabelError != nil {
		return B51CDPartialLabel{}, false, b51CDPartialLabelError
	}
	b51CDPartialLabelEntry, b51CDPartialLabelFound := b51CDPartialLabelFindPartialLabel(b51CDPartialLabelCatalog.B51CDPartialLabelEntries, b51CDPartialLabelKey)
	return b51CDPartialLabelEntry, b51CDPartialLabelFound, nil
}
