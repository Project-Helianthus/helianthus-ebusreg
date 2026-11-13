package vaillant

import (
	b523ConfigMixerBytes "bytes"
	_ "embed"
	b523ConfigMixerJSON "encoding/json"
	b523ConfigMixerFmt "fmt"
	b523ConfigMixerIO "io"
	b523ConfigMixerReflect "reflect"
)

//go:embed data/ab99_b523_config_mixer.json
var b523ConfigMixerData []byte

type B523ConfigMixerKey struct {
	B523ConfigMixerFamily            string  `json:"family"`
	B523ConfigMixerVariant           string  `json:"variant"`
	B523ConfigMixerPrimary           uint8   `json:"primary"`
	B523ConfigMixerSecondary         *uint8  `json:"secondary"`
	B523ConfigMixerOperation         *uint8  `json:"operation"`
	B523ConfigMixerDirection         string  `json:"direction"`
	B523ConfigMixerField             string  `json:"field"`
	B523ConfigMixerGroup             *uint8  `json:"group"`
	B523ConfigMixerRequestedInstance *uint8  `json:"requested_instance"`
	B523ConfigMixerRegister          *uint16 `json:"register"`
	B523ConfigMixerCommandRole       *string `json:"command_role"`
}

type B523ConfigMixerBit struct {
	B523ConfigMixerIndex    uint8   `json:"index"`
	B523ConfigMixerLabel    *string `json:"label"`
	B523ConfigMixerReserved bool    `json:"reserved"`
	B523ConfigMixerDefault  uint8   `json:"default"`
}

type B523ConfigMixerCandidate struct {
	B523ConfigMixerLabel   string   `json:"label"`
	B523ConfigMixerSigned  *bool    `json:"signed"`
	B523ConfigMixerScale   *float64 `json:"scale"`
	B523ConfigMixerMinimum *float64 `json:"minimum"`
	B523ConfigMixerMaximum *float64 `json:"maximum"`
	B523ConfigMixerRaw     *uint8   `json:"raw"`
	B523ConfigMixerValue   *float64 `json:"value"`
	B523ConfigMixerClamp   *float64 `json:"clamp"`
}

type B523ConfigMixerInterval struct {
	B523ConfigMixerMinimum int    `json:"minimum"`
	B523ConfigMixerMaximum int    `json:"maximum"`
	B523ConfigMixerLabel   string `json:"label"`
}

type B523ConfigMixerEntry struct {
	B523ConfigMixerSigned                 *bool                      `json:"signed"`
	B523ConfigMixerKey                    B523ConfigMixerKey         `json:"key"`
	B523ConfigMixerGrade                  string                     `json:"grade"`
	B523ConfigMixerPhysicalAdaptation     string                     `json:"physical_adaptation"`
	B523ConfigMixerSupport                string                     `json:"support"`
	B523ConfigMixerExecutedMovement       string                     `json:"executed_movement"`
	B523ConfigMixerSelectedInterpretation string                     `json:"selected_interpretation"`
	B523ConfigMixerUnknownMetadata        string                     `json:"unknown_metadata"`
	B523ConfigMixerPreserveRaw            bool                       `json:"preserve_raw"`
	B523ConfigMixerPosition               *uint8                     `json:"position"`
	B523ConfigMixerWidthBits              *uint8                     `json:"width_bits"`
	B523ConfigMixerRequestExtent          *uint8                     `json:"request_extent"`
	B523ConfigMixerPositiveExtent         *uint8                     `json:"positive_extent"`
	B523ConfigMixerDefault                *float64                   `json:"default"`
	B523ConfigMixerValues                 map[uint8]string           `json:"values"`
	B523ConfigMixerOtherValues            string                     `json:"other_values"`
	B523ConfigMixerBits                   []B523ConfigMixerBit       `json:"bits"`
	B523ConfigMixerOneHot                 *bool                      `json:"one_hot"`
	B523ConfigMixerAlternatives           []string                   `json:"alternatives"`
	B523ConfigMixerRelevantModes          []uint8                    `json:"relevant_modes"`
	B523ConfigMixerIrrelevantModes        []uint8                    `json:"irrelevant_modes"`
	B523ConfigMixerCandidates             []B523ConfigMixerCandidate `json:"candidates"`
	B523ConfigMixerIntervals              []B523ConfigMixerInterval  `json:"intervals"`
	B523ConfigMixerUnit                   *string                    `json:"unit"`
	B523ConfigMixerNotes                  []string                   `json:"notes"`
}

func B523ConfigMixerEntries() ([]B523ConfigMixerEntry, error) {
	var b523ConfigMixerEntries []B523ConfigMixerEntry
	b523ConfigMixerDecoder := b523ConfigMixerJSON.NewDecoder(b523ConfigMixerBytes.NewReader(b523ConfigMixerData))
	b523ConfigMixerDecoder.DisallowUnknownFields()
	if b523ConfigMixerError := b523ConfigMixerDecoder.Decode(&b523ConfigMixerEntries); b523ConfigMixerError != nil {
		return nil, b523ConfigMixerError
	}
	var b523ConfigMixerTrailing any
	if b523ConfigMixerError := b523ConfigMixerDecoder.Decode(&b523ConfigMixerTrailing); b523ConfigMixerError != b523ConfigMixerIO.EOF {
		return nil, b523ConfigMixerFmt.Errorf("ab99: trailing data")
	}
	if len(b523ConfigMixerEntries) == 0 {
		return nil, b523ConfigMixerFmt.Errorf("ab99: empty catalog")
	}
	for b523ConfigMixerIndex, b523ConfigMixerEntry := range b523ConfigMixerEntries {
		if b523ConfigMixerEntry.B523ConfigMixerKey.B523ConfigMixerPrimary != 0xb5 || b523ConfigMixerEntry.B523ConfigMixerKey.B523ConfigMixerSecondary == nil || *b523ConfigMixerEntry.B523ConfigMixerKey.B523ConfigMixerSecondary != 0x23 || b523ConfigMixerEntry.B523ConfigMixerKey.B523ConfigMixerOperation == nil {
			return nil, b523ConfigMixerFmt.Errorf("ab99: invalid key coordinates in entry %d", b523ConfigMixerIndex)
		}
		if b523ConfigMixerEntry.B523ConfigMixerKey.B523ConfigMixerFamily == "" || b523ConfigMixerEntry.B523ConfigMixerKey.B523ConfigMixerVariant == "" || b523ConfigMixerEntry.B523ConfigMixerKey.B523ConfigMixerDirection == "" || b523ConfigMixerEntry.B523ConfigMixerKey.B523ConfigMixerField == "" || b523ConfigMixerEntry.B523ConfigMixerGrade != "Hypothesis" || !b523ConfigMixerEntry.B523ConfigMixerPreserveRaw || b523ConfigMixerEntry.B523ConfigMixerPhysicalAdaptation != "Unknown" || b523ConfigMixerEntry.B523ConfigMixerSupport != "Unknown" || b523ConfigMixerEntry.B523ConfigMixerExecutedMovement != "Unknown" || b523ConfigMixerEntry.B523ConfigMixerSelectedInterpretation != "Unknown" || b523ConfigMixerEntry.B523ConfigMixerUnknownMetadata != "Unknown" {
			return nil, b523ConfigMixerFmt.Errorf("ab99: invalid entry %d", b523ConfigMixerIndex)
		}
		for _, b523ConfigMixerPrevious := range b523ConfigMixerEntries[:b523ConfigMixerIndex] {
			if b523ConfigMixerReflect.DeepEqual(b523ConfigMixerPrevious.B523ConfigMixerKey, b523ConfigMixerEntry.B523ConfigMixerKey) {
				return nil, b523ConfigMixerFmt.Errorf("ab99: duplicate key")
			}
		}
	}
	return b523ConfigMixerEntries, nil
}

func B523ConfigMixerLookup(b523ConfigMixerKey B523ConfigMixerKey) (B523ConfigMixerEntry, bool, error) {
	b523ConfigMixerEntries, b523ConfigMixerError := B523ConfigMixerEntries()
	if b523ConfigMixerError != nil {
		return B523ConfigMixerEntry{}, false, b523ConfigMixerError
	}
	for _, b523ConfigMixerEntry := range b523ConfigMixerEntries {
		if b523ConfigMixerReflect.DeepEqual(b523ConfigMixerEntry.B523ConfigMixerKey, b523ConfigMixerKey) {
			return b523ConfigMixerEntry, true, nil
		}
	}
	return B523ConfigMixerEntry{}, false, nil
}
