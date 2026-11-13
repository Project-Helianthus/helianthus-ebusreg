package vaillant

import (
	b516StatisticsDirectionBytes "bytes"
	_ "embed"
	b516StatisticsDirectionJSON "encoding/json"
	b516StatisticsDirectionFmt "fmt"
	b516StatisticsDirectionIO "io"
	b516StatisticsDirectionReflect "reflect"
)

//go:embed data/ab72_b516_statistics_directions.json
var b516StatisticsDirectionData []byte

type B516StatisticsDirectionKey struct {
	B516StatisticsDirectionFamily    string `json:"family"`
	B516StatisticsDirectionVariant   string `json:"variant"`
	B516StatisticsDirectionOperation *uint8 `json:"operation"`
	B516StatisticsDirection          string `json:"direction"`
	B516StatisticsDirectionSelector  *uint8 `json:"selector"`
	B516StatisticsDirectionField     string `json:"field"`
}

type B516StatisticsDirectionEntry struct {
	B516StatisticsDirectionKey              B516StatisticsDirectionKey `json:"key"`
	B516StatisticsDirectionGrade            string                     `json:"grade"`
	B516StatisticsDirectionSoftwareContract string                     `json:"software_contract"`
	B516StatisticsDirectionPhysicalSupport  string                     `json:"physical_support"`
	B516StatisticsDirectionPresence         string                     `json:"presence"`
	B516StatisticsDirectionFreshness        string                     `json:"freshness"`
	B516StatisticsDirectionExecution        string                     `json:"execution"`
	B516StatisticsDirectionEffects          string                     `json:"effects"`
	B516StatisticsDirectionPersistence      string                     `json:"persistence"`
	B516StatisticsDirectionDelivery         string                     `json:"delivery"`
	B516StatisticsDirectionNativeContract   string                     `json:"native_contract"`
	B516StatisticsDirectionLength           *uint8                     `json:"length"`
	B516StatisticsDirectionPosition         *uint8                     `json:"position"`
	B516StatisticsDirectionWidth            *uint8                     `json:"width"`
	B516StatisticsDirectionEncoding         *string                    `json:"encoding"`
	B516StatisticsDirectionUnit             *string                    `json:"unit"`
	B516StatisticsDirectionBits             map[string]uint8           `json:"bits"`
	B516StatisticsDirectionValues           map[string]string          `json:"values"`
	B516StatisticsDirectionOrder            []string                   `json:"order"`
	B516StatisticsDirectionRules            []string                   `json:"rules"`
	B516StatisticsDirectionConflicts        []string                   `json:"conflicts"`
	B516StatisticsDirectionUnknown          []string                   `json:"unknown"`
}

func b516StatisticsDirectionParse(b516StatisticsDirectionRaw []byte) ([]B516StatisticsDirectionEntry, error) {
	var b516StatisticsDirectionEntries []B516StatisticsDirectionEntry
	b516StatisticsDirectionDecoder := b516StatisticsDirectionJSON.NewDecoder(b516StatisticsDirectionBytes.NewReader(b516StatisticsDirectionRaw))
	b516StatisticsDirectionDecoder.DisallowUnknownFields()
	if b516StatisticsDirectionErr := b516StatisticsDirectionDecoder.Decode(&b516StatisticsDirectionEntries); b516StatisticsDirectionErr != nil {
		return nil, b516StatisticsDirectionFmt.Errorf("ab72 catalog: %w", b516StatisticsDirectionErr)
	}
	if b516StatisticsDirectionErr := b516StatisticsDirectionDecoder.Decode(new(any)); b516StatisticsDirectionErr != b516StatisticsDirectionIO.EOF {
		return nil, b516StatisticsDirectionFmt.Errorf("ab72 catalog: trailing input")
	}
	if len(b516StatisticsDirectionEntries) == 0 {
		return nil, b516StatisticsDirectionFmt.Errorf("ab72 catalog: empty")
	}
	for b516StatisticsDirectionIndex, b516StatisticsDirectionEntry := range b516StatisticsDirectionEntries {
		if b516StatisticsDirectionEntry.B516StatisticsDirectionKey.B516StatisticsDirectionFamily == "" || b516StatisticsDirectionEntry.B516StatisticsDirectionKey.B516StatisticsDirectionVariant == "" || b516StatisticsDirectionEntry.B516StatisticsDirectionKey.B516StatisticsDirectionField == "" || b516StatisticsDirectionEntry.B516StatisticsDirectionKey.B516StatisticsDirectionOperation == nil || (b516StatisticsDirectionEntry.B516StatisticsDirectionKey.B516StatisticsDirection != "request" && b516StatisticsDirectionEntry.B516StatisticsDirectionKey.B516StatisticsDirection != "positive") {
			return nil, b516StatisticsDirectionFmt.Errorf("ab72 catalog: incomplete key")
		}
		if b516StatisticsDirectionEntry.B516StatisticsDirectionGrade != "Hypothesis" || b516StatisticsDirectionEntry.B516StatisticsDirectionSoftwareContract != "Proven" || b516StatisticsDirectionEntry.B516StatisticsDirectionNativeContract != "https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/protocols/ebus-services/ebus-overview.md" {
			return nil, b516StatisticsDirectionFmt.Errorf("ab72 catalog: invalid evidence")
		}
		for _, b516StatisticsDirectionGrade := range []string{b516StatisticsDirectionEntry.B516StatisticsDirectionPhysicalSupport, b516StatisticsDirectionEntry.B516StatisticsDirectionPresence, b516StatisticsDirectionEntry.B516StatisticsDirectionFreshness, b516StatisticsDirectionEntry.B516StatisticsDirectionExecution, b516StatisticsDirectionEntry.B516StatisticsDirectionEffects, b516StatisticsDirectionEntry.B516StatisticsDirectionPersistence, b516StatisticsDirectionEntry.B516StatisticsDirectionDelivery} {
			if b516StatisticsDirectionGrade != "Unknown" {
				return nil, b516StatisticsDirectionFmt.Errorf("ab72 catalog: invalid qualification")
			}
		}
		for _, b516StatisticsDirectionBit := range b516StatisticsDirectionEntry.B516StatisticsDirectionBits {
			if b516StatisticsDirectionBit > 7 {
				return nil, b516StatisticsDirectionFmt.Errorf("ab72 catalog: invalid bit")
			}
		}
		if b516StatisticsDirectionEntry.B516StatisticsDirectionWidth != nil && (*b516StatisticsDirectionEntry.B516StatisticsDirectionWidth == 0 || *b516StatisticsDirectionEntry.B516StatisticsDirectionWidth > 32) {
			return nil, b516StatisticsDirectionFmt.Errorf("ab72 catalog: invalid width")
		}
		for b516StatisticsDirectionPrevious := 0; b516StatisticsDirectionPrevious < b516StatisticsDirectionIndex; b516StatisticsDirectionPrevious++ {
			if b516StatisticsDirectionReflect.DeepEqual(b516StatisticsDirectionEntries[b516StatisticsDirectionPrevious].B516StatisticsDirectionKey, b516StatisticsDirectionEntry.B516StatisticsDirectionKey) {
				return nil, b516StatisticsDirectionFmt.Errorf("ab72 catalog: duplicate key")
			}
		}
	}
	return b516StatisticsDirectionEntries, nil
}

func B516StatisticsDirectionEntries() ([]B516StatisticsDirectionEntry, error) {
	return b516StatisticsDirectionParse(b516StatisticsDirectionData)
}

func B516StatisticsDirectionLookup(b516StatisticsDirectionKey B516StatisticsDirectionKey) (B516StatisticsDirectionEntry, bool, error) {
	b516StatisticsDirectionEntries, b516StatisticsDirectionErr := B516StatisticsDirectionEntries()
	if b516StatisticsDirectionErr != nil {
		return B516StatisticsDirectionEntry{}, false, b516StatisticsDirectionErr
	}
	for _, b516StatisticsDirectionEntry := range b516StatisticsDirectionEntries {
		if b516StatisticsDirectionReflect.DeepEqual(b516StatisticsDirectionKey, b516StatisticsDirectionEntry.B516StatisticsDirectionKey) {
			return b516StatisticsDirectionEntry, true, nil
		}
	}
	return B516StatisticsDirectionEntry{}, false, nil
}
