package vaillant

import (
	_ "embed"
	b514MenuFeedbackJSON "encoding/json"
)

type B514MenuFeedbackAxis struct {
	B514MenuFeedbackPresent bool   `json:"present"`
	B514MenuFeedbackValue   uint16 `json:"value"`
}

type B514MenuFeedbackKey struct {
	B514MenuFeedbackFamily    string               `json:"family"`
	B514MenuFeedbackVariant   string               `json:"variant"`
	B514MenuFeedbackPrimary   uint16               `json:"primary"`
	B514MenuFeedbackSecondary B514MenuFeedbackAxis `json:"secondary"`
	B514MenuFeedbackOperation B514MenuFeedbackAxis `json:"operation"`
	B514MenuFeedbackDirection string               `json:"direction"`
	B514MenuFeedbackField     string               `json:"field"`
	B514MenuFeedbackOpcode    B514MenuFeedbackAxis `json:"opcode"`
	B514MenuFeedbackGroup     B514MenuFeedbackAxis `json:"group"`
	B514MenuFeedbackInstance  B514MenuFeedbackAxis `json:"instance"`
	B514MenuFeedbackRegister  B514MenuFeedbackAxis `json:"register"`
}

type B514MenuFeedbackConversion struct {
	B514MenuFeedbackFactor float64 `json:"factor"`
	B514MenuFeedbackOffset float64 `json:"offset"`
}

type B514MenuFeedbackRange struct {
	B514MenuFeedbackMinimum int    `json:"minimum"`
	B514MenuFeedbackMaximum int    `json:"maximum"`
	B514MenuFeedbackMeaning string `json:"meaning"`
}

type B514MenuFeedbackField struct {
	B514MenuFeedbackKey         B514MenuFeedbackKey         `json:"key"`
	B514MenuFeedbackPosition    *int                        `json:"logical_position"`
	B514MenuFeedbackWidth       *int                        `json:"width"`
	B514MenuFeedbackExtent      *int                        `json:"logical_extent"`
	B514MenuFeedbackSignedness  string                      `json:"signedness"`
	B514MenuFeedbackByteOrder   string                      `json:"byte_order"`
	B514MenuFeedbackConversion  *B514MenuFeedbackConversion `json:"conversion"`
	B514MenuFeedbackUnit        *string                     `json:"unit"`
	B514MenuFeedbackDefault     *int                        `json:"default"`
	B514MenuFeedbackBounds      []int                       `json:"bounds"`
	B514MenuFeedbackEnumeration map[int]string              `json:"enumeration"`
	B514MenuFeedbackRanges      []B514MenuFeedbackRange     `json:"ranges"`
	B514MenuFeedbackReplacement *string                     `json:"replacement"`
	B514MenuFeedbackReserved    []B514MenuFeedbackRange     `json:"reserved"`
	B514MenuFeedbackConstant    *int                        `json:"constant"`
	B514MenuFeedbackConflict    *string                     `json:"conflict"`
	B514MenuFeedbackGrade       string                      `json:"grade"`
	B514MenuFeedbackRawFallback string                      `json:"raw_fallback"`
}

type B514MenuFeedbackCatalog struct {
	B514MenuFeedbackRepresentationGrade string                  `json:"representation_grade"`
	B514MenuFeedbackQualification       string                  `json:"physical_qualification"`
	B514MenuFeedbackLimitations         []string                `json:"limitations"`
	B514MenuFeedbackFields              []B514MenuFeedbackField `json:"fields"`
}

//go:embed data/ab89_b514_menu_feedback.json
var b514MenuFeedbackMenuFeedbackJSON string

func B514MenuFeedback() (B514MenuFeedbackCatalog, error) {
	var b514MenuFeedbackCatalog B514MenuFeedbackCatalog
	b514MenuFeedbackError := b514MenuFeedbackJSON.Unmarshal([]byte(b514MenuFeedbackMenuFeedbackJSON), &b514MenuFeedbackCatalog)
	return b514MenuFeedbackCatalog, b514MenuFeedbackError
}

func B514MenuFeedbackLookup(b514MenuFeedbackKey B514MenuFeedbackKey) (B514MenuFeedbackField, bool, error) {
	b514MenuFeedbackCatalog, b514MenuFeedbackError := B514MenuFeedback()
	if b514MenuFeedbackError != nil {
		return B514MenuFeedbackField{}, false, b514MenuFeedbackError
	}
	return b514MenuFeedbackLookupFields(b514MenuFeedbackCatalog.B514MenuFeedbackFields, b514MenuFeedbackKey)
}

func b514MenuFeedbackLookupFields(b514MenuFeedbackFields []B514MenuFeedbackField, b514MenuFeedbackKey B514MenuFeedbackKey) (B514MenuFeedbackField, bool, error) {
	for _, b514MenuFeedbackField := range b514MenuFeedbackFields {
		if b514MenuFeedbackField.B514MenuFeedbackKey == b514MenuFeedbackKey {
			return b514MenuFeedbackField, true, nil
		}
	}
	return B514MenuFeedbackField{}, false, nil
}
