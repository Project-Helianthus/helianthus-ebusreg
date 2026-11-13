package vaillant

import (
	b507EnergyCandidateBytes "bytes"
	b507EnergyCandidateEmbed "embed"
	b507EnergyCandidateJSON "encoding/json"
	b507EnergyCandidateFmt "fmt"
	b507EnergyCandidateIO "io"
	b507EnergyCandidateReflect "reflect"
)

//go:embed data/ab13_b507_energy_candidates.json
var b507EnergyCandidateFiles b507EnergyCandidateEmbed.FS

var b507EnergyCandidateData, b507EnergyCandidateReadError = b507EnergyCandidateFiles.ReadFile("data/ab13_b507_energy_candidates.json")

type B507EnergyCandidateKey struct {
	B507EnergyCandidateFamily            string  `json:"family"`
	B507EnergyCandidateVariant           string  `json:"variant"`
	B507EnergyCandidatePrimary           uint8   `json:"primary"`
	B507EnergyCandidateSecondary         *uint8  `json:"secondary"`
	B507EnergyCandidateOperation         *uint8  `json:"operation"`
	B507EnergyCandidateDirection         *string `json:"direction"`
	B507EnergyCandidateField             string  `json:"field"`
	B507EnergyCandidateGroup             *uint8  `json:"group"`
	B507EnergyCandidateRequestedInstance *uint8  `json:"requested_instance"`
	B507EnergyCandidateRegister          *uint16 `json:"register"`
	B507EnergyCandidateCommandRole       *string `json:"command_role"`
}

type B507EnergyCandidateEntry struct {
	B507EnergyCandidateKey                   B507EnergyCandidateKey `json:"key"`
	B507EnergyCandidateGrade                 string                 `json:"grade"`
	B507EnergyCandidateMetadataPresence      string                 `json:"metadata_presence"`
	B507EnergyCandidatePhysicalQualification string                 `json:"physical_qualification"`
	B507EnergyCandidateInvocation            string                 `json:"invocation"`
	B507EnergyCandidateNativeContract        string                 `json:"native_contract"`
	B507EnergyCandidatePosition              *uint16                `json:"position"`
	B507EnergyCandidateWidthBits             *uint8                 `json:"width_bits"`
	B507EnergyCandidateEncoding              *string                `json:"encoding"`
	B507EnergyCandidateRawDomain             []uint16               `json:"raw_domain"`
	B507EnergyCandidateDivisor               *float64               `json:"divisor"`
	B507EnergyCandidateUnit                  *string                `json:"unit"`
	B507EnergyCandidateValues                map[uint16]string      `json:"values"`
	B507EnergyCandidateDependencies          []string               `json:"dependencies"`
	B507EnergyCandidateReference             string                 `json:"reference"`
	B507EnergyCandidateQuantityOrigin        string                 `json:"quantity_origin"`
	B507EnergyCandidatePreserveRaw           bool                   `json:"preserve_raw"`
	B507EnergyCandidateUnknowns              []string               `json:"unknowns"`
	B507EnergyCandidateAlternatives          []string               `json:"alternatives"`
}

func B507EnergyCandidateEntries() ([]B507EnergyCandidateEntry, error) {
	if b507EnergyCandidateReadError != nil {
		return nil, b507EnergyCandidateReadError
	}
	var b507EnergyCandidateEntries []B507EnergyCandidateEntry
	b507EnergyCandidateDecoder := b507EnergyCandidateJSON.NewDecoder(b507EnergyCandidateBytes.NewReader(b507EnergyCandidateData))
	b507EnergyCandidateDecoder.DisallowUnknownFields()
	if b507EnergyCandidateError := b507EnergyCandidateDecoder.Decode(&b507EnergyCandidateEntries); b507EnergyCandidateError != nil {
		return nil, b507EnergyCandidateError
	}
	var b507EnergyCandidateExtra any
	if b507EnergyCandidateError := b507EnergyCandidateDecoder.Decode(&b507EnergyCandidateExtra); b507EnergyCandidateError != b507EnergyCandidateIO.EOF {
		return nil, b507EnergyCandidateFmt.Errorf("ab13: trailing embedded data")
	}
	if len(b507EnergyCandidateEntries) == 0 {
		return nil, b507EnergyCandidateFmt.Errorf("ab13: empty embedded catalog")
	}
	for b507EnergyCandidateIndex, b507EnergyCandidateEntry := range b507EnergyCandidateEntries {
		if b507EnergyCandidateEntry.B507EnergyCandidateKey.B507EnergyCandidateFamily == "" || b507EnergyCandidateEntry.B507EnergyCandidateKey.B507EnergyCandidateVariant == "" || b507EnergyCandidateEntry.B507EnergyCandidateKey.B507EnergyCandidateField == "" ||
			(b507EnergyCandidateEntry.B507EnergyCandidateGrade != "Proven" && b507EnergyCandidateEntry.B507EnergyCandidateGrade != "Hypothesis" && b507EnergyCandidateEntry.B507EnergyCandidateGrade != "Unknown") ||
			b507EnergyCandidateEntry.B507EnergyCandidateMetadataPresence != "Proven" || b507EnergyCandidateEntry.B507EnergyCandidatePhysicalQualification != "Unknown" ||
			b507EnergyCandidateEntry.B507EnergyCandidateInvocation != "Unknown" || b507EnergyCandidateEntry.B507EnergyCandidateNativeContract == "" || !b507EnergyCandidateEntry.B507EnergyCandidatePreserveRaw {
			return nil, b507EnergyCandidateFmt.Errorf("ab13: invalid embedded entry %d", b507EnergyCandidateIndex)
		}
		if b507EnergyCandidateEntry.B507EnergyCandidateDivisor != nil && *b507EnergyCandidateEntry.B507EnergyCandidateDivisor <= 0 {
			return nil, b507EnergyCandidateFmt.Errorf("ab13: invalid divisor at %d", b507EnergyCandidateIndex)
		}
		if b507EnergyCandidateEntry.B507EnergyCandidateRawDomain != nil && (len(b507EnergyCandidateEntry.B507EnergyCandidateRawDomain) != 2 || b507EnergyCandidateEntry.B507EnergyCandidateRawDomain[0] > b507EnergyCandidateEntry.B507EnergyCandidateRawDomain[1]) {
			return nil, b507EnergyCandidateFmt.Errorf("ab13: invalid domain at %d", b507EnergyCandidateIndex)
		}
		for _, b507EnergyCandidatePrevious := range b507EnergyCandidateEntries[:b507EnergyCandidateIndex] {
			if b507EnergyCandidateReflect.DeepEqual(b507EnergyCandidatePrevious.B507EnergyCandidateKey, b507EnergyCandidateEntry.B507EnergyCandidateKey) {
				return nil, b507EnergyCandidateFmt.Errorf("ab13: duplicate embedded key at %d", b507EnergyCandidateIndex)
			}
		}
	}
	return b507EnergyCandidateEntries, nil
}

func B507EnergyCandidateLookup(b507EnergyCandidateKey B507EnergyCandidateKey) (B507EnergyCandidateEntry, bool, error) {
	b507EnergyCandidateEntries, b507EnergyCandidateError := B507EnergyCandidateEntries()
	if b507EnergyCandidateError != nil {
		return B507EnergyCandidateEntry{}, false, b507EnergyCandidateError
	}
	for _, b507EnergyCandidateEntry := range b507EnergyCandidateEntries {
		if b507EnergyCandidateReflect.DeepEqual(b507EnergyCandidateKey, b507EnergyCandidateEntry.B507EnergyCandidateKey) {
			return b507EnergyCandidateEntry, true, nil
		}
	}
	return B507EnergyCandidateEntry{}, false, nil
}
