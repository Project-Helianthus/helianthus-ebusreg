package vaillant

import (
	_ "embed"
	b523CandidateJSON "encoding/json"
	b523CandidateReflect "reflect"
)

//go:embed data/ab71_b523_candidates.json
var b523CandidateB523Data []byte

type B523CandidateKey struct {
	B523CandidateFamily    string  `json:"family"`
	B523CandidateVariant   *string `json:"variant"`
	B523CandidatePrimary   byte    `json:"primary"`
	B523CandidateSecondary *byte   `json:"secondary"`
	B523CandidateOperation *byte   `json:"operation"`
	B523CandidateDirection string  `json:"direction"`
	B523CandidateField     string  `json:"field"`
}

type B523CandidateRule struct {
	B523CandidateRawMin       byte     `json:"raw_min"`
	B523CandidateRawMax       byte     `json:"raw_max"`
	B523CandidateResult       string   `json:"result"`
	B523CandidateDivideBy     *float64 `json:"divide_by"`
	B523CandidateDurationUnit *string  `json:"duration_unit"`
	B523CandidateGrade        string   `json:"grade"`
}

type B523CandidateAlternative struct {
	B523CandidateName  string              `json:"name"`
	B523CandidateGrade string              `json:"grade"`
	B523CandidateRules []B523CandidateRule `json:"rules"`
}

type B523CandidateEntry struct {
	B523CandidateKey             B523CandidateKey           `json:"key"`
	B523CandidateLogicalPosition *int                       `json:"logical_position"`
	B523CandidateWidthBits       *int                       `json:"width_bits"`
	B523CandidateRequestExtent   *int                       `json:"request_extent"`
	B523CandidatePositiveExtent  *int                       `json:"positive_extent"`
	B523CandidateDefaultRaw      *byte                      `json:"default_raw"`
	B523CandidateDefaultResult   *string                    `json:"default_result"`
	B523CandidateGrade           string                     `json:"grade"`
	B523CandidateRawDomain       []byte                     `json:"raw_domain"`
	B523CandidateAlternatives    []B523CandidateAlternative `json:"alternatives"`
}

type B523CandidateCatalog struct {
	B523CandidateRepresentationGrade string               `json:"representation_grade"`
	B523CandidateNativeSupport       string               `json:"native_support"`
	B523CandidateRuntimeEffects      string               `json:"runtime_effects"`
	B523CandidateOtherModuleLayouts  string               `json:"other_module_layouts"`
	B523CandidateLimitations         []string             `json:"limitations"`
	B523CandidateEntries             []B523CandidateEntry `json:"entries"`
}

func B523Candidates() (B523CandidateCatalog, error) {
	var b523CandidateCatalog B523CandidateCatalog
	b523CandidateError := b523CandidateJSON.Unmarshal(b523CandidateB523Data, &b523CandidateCatalog)
	return b523CandidateCatalog, b523CandidateError
}

func B523CandidateLookup(b523CandidateKey B523CandidateKey) (B523CandidateEntry, bool, error) {
	b523CandidateCatalog, b523CandidateError := B523Candidates()
	if b523CandidateError != nil {
		return B523CandidateEntry{}, false, b523CandidateError
	}
	for _, b523CandidateEntry := range b523CandidateCatalog.B523CandidateEntries {
		if b523CandidateReflect.DeepEqual(b523CandidateEntry.B523CandidateKey, b523CandidateKey) {
			return b523CandidateEntry, true, nil
		}
	}
	return B523CandidateEntry{}, false, nil
}
