package vaillant

import (
	_ "embed"
	b512HeatpumpIOJSON "encoding/json"
	b512HeatpumpIOFmt "fmt"
	b512HeatpumpIOReflect "reflect"
)

//go:embed data/ab19_b512_heatpump_io.json
var b512HeatpumpIOData string

type B512HeatpumpIOKey struct {
	B512HeatpumpIOFamily        string  `json:"family"`
	B512HeatpumpIOVariant       string  `json:"variant"`
	B512HeatpumpIOPrimary       uint8   `json:"primary"`
	B512HeatpumpIOSecondary     *uint8  `json:"secondary"`
	B512HeatpumpIOOperation     *uint8  `json:"operation"`
	B512HeatpumpIODirection     string  `json:"direction"`
	B512HeatpumpIOField         string  `json:"field"`
	B512HeatpumpIOAlternative   string  `json:"alternative"`
	B512HeatpumpIOGroup         *uint8  `json:"group"`
	B512HeatpumpIOInstanceScope string  `json:"instance_scope"`
	B512HeatpumpIORegister      *uint16 `json:"register"`
}

func b512HeatpumpIOValidate(b512HeatpumpIOInput string) ([]b512HeatpumpIOJSON.RawMessage, []B512HeatpumpIOKey, error) {
	var b512HeatpumpIORows []b512HeatpumpIOJSON.RawMessage
	if b512HeatpumpIOErr := b512HeatpumpIOJSON.Unmarshal([]byte(b512HeatpumpIOInput), &b512HeatpumpIORows); b512HeatpumpIOErr != nil {
		return nil, nil, b512HeatpumpIOFmt.Errorf("ab19 catalog: %w", b512HeatpumpIOErr)
	}
	if len(b512HeatpumpIORows) == 0 {
		return nil, nil, b512HeatpumpIOFmt.Errorf("ab19 catalog: empty array")
	}
	b512HeatpumpIOKeys := make([]B512HeatpumpIOKey, 0, len(b512HeatpumpIORows))
	b512HeatpumpIOIDs := make(map[string]bool)
	for _, b512HeatpumpIORow := range b512HeatpumpIORows {
		var b512HeatpumpIOFields map[string]b512HeatpumpIOJSON.RawMessage
		if b512HeatpumpIOErr := b512HeatpumpIOJSON.Unmarshal(b512HeatpumpIORow, &b512HeatpumpIOFields); b512HeatpumpIOErr != nil {
			return nil, nil, b512HeatpumpIOErr
		}
		var b512HeatpumpIOKeyFields map[string]b512HeatpumpIOJSON.RawMessage
		if b512HeatpumpIOErr := b512HeatpumpIOJSON.Unmarshal(b512HeatpumpIOFields["key"], &b512HeatpumpIOKeyFields); b512HeatpumpIOErr != nil {
			return nil, nil, b512HeatpumpIOErr
		}
		for _, b512HeatpumpIOAxis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			if _, b512HeatpumpIOPresent := b512HeatpumpIOKeyFields[b512HeatpumpIOAxis]; !b512HeatpumpIOPresent {
				return nil, nil, b512HeatpumpIOFmt.Errorf("ab19 catalog: missing %s", b512HeatpumpIOAxis)
			}
		}
		var b512HeatpumpIOKey B512HeatpumpIOKey
		if b512HeatpumpIOErr := b512HeatpumpIOJSON.Unmarshal(b512HeatpumpIOFields["key"], &b512HeatpumpIOKey); b512HeatpumpIOErr != nil {
			return nil, nil, b512HeatpumpIOErr
		}
		if len(b512HeatpumpIOKeyFields) != 11 || b512HeatpumpIOKey.B512HeatpumpIOFamily == "" || b512HeatpumpIOKey.B512HeatpumpIOVariant == "" || b512HeatpumpIOKey.B512HeatpumpIODirection == "" || b512HeatpumpIOKey.B512HeatpumpIOField == "" || b512HeatpumpIOKey.B512HeatpumpIOAlternative == "" || b512HeatpumpIOKey.B512HeatpumpIOInstanceScope == "" || string(b512HeatpumpIOKeyFields["primary"]) == "null" {
			return nil, nil, b512HeatpumpIOFmt.Errorf("ab19 catalog: incomplete identity")
		}
		for _, b512HeatpumpIOPrevious := range b512HeatpumpIOKeys {
			if b512HeatpumpIOReflect.DeepEqual(b512HeatpumpIOPrevious, b512HeatpumpIOKey) {
				return nil, nil, b512HeatpumpIOFmt.Errorf("ab19 catalog: duplicate key")
			}
		}
		var b512HeatpumpIOID string
		if b512HeatpumpIOErr := b512HeatpumpIOJSON.Unmarshal(b512HeatpumpIOFields["id"], &b512HeatpumpIOID); b512HeatpumpIOErr != nil {
			return nil, nil, b512HeatpumpIOErr
		}
		if len(b512HeatpumpIOID) < 5 || b512HeatpumpIOID[:5] != "ab19_" || b512HeatpumpIOIDs[b512HeatpumpIOID] {
			return nil, nil, b512HeatpumpIOFmt.Errorf("ab19 catalog: invalid ID")
		}
		b512HeatpumpIOIDs[b512HeatpumpIOID] = true
		b512HeatpumpIOKeys = append(b512HeatpumpIOKeys, b512HeatpumpIOKey)
	}
	return b512HeatpumpIORows, b512HeatpumpIOKeys, nil
}

func b512HeatpumpIOEntries(b512HeatpumpIOInput string) (b512HeatpumpIOJSON.RawMessage, error) {
	if _, _, b512HeatpumpIOErr := b512HeatpumpIOValidate(b512HeatpumpIOInput); b512HeatpumpIOErr != nil {
		return nil, b512HeatpumpIOErr
	}
	return append(b512HeatpumpIOJSON.RawMessage(nil), b512HeatpumpIOInput...), nil
}

func B512HeatpumpIOEntries() (b512HeatpumpIOJSON.RawMessage, error) {
	return b512HeatpumpIOEntries(b512HeatpumpIOData)
}

func b512HeatpumpIOLookup(b512HeatpumpIOInput string, b512HeatpumpIOKey B512HeatpumpIOKey) (b512HeatpumpIOJSON.RawMessage, bool, error) {
	b512HeatpumpIORows, b512HeatpumpIOKeys, b512HeatpumpIOErr := b512HeatpumpIOValidate(b512HeatpumpIOInput)
	if b512HeatpumpIOErr != nil {
		return nil, false, b512HeatpumpIOErr
	}
	for b512HeatpumpIOIndex, b512HeatpumpIOCandidate := range b512HeatpumpIOKeys {
		if b512HeatpumpIOReflect.DeepEqual(b512HeatpumpIOCandidate, b512HeatpumpIOKey) {
			return append(b512HeatpumpIOJSON.RawMessage(nil), b512HeatpumpIORows[b512HeatpumpIOIndex]...), true, nil
		}
	}
	return nil, false, nil
}

func B512HeatpumpIOLookup(b512HeatpumpIOKey B512HeatpumpIOKey) (b512HeatpumpIOJSON.RawMessage, bool, error) {
	return b512HeatpumpIOLookup(b512HeatpumpIOData, b512HeatpumpIOKey)
}
