package ebusstandard

import (
	_ "embed"
	standard07ResidualJSON "encoding/json"
	standard07ResidualFmt "fmt"
)

type Standard07ResidualKey struct {
	Standard07ResidualFamily        string  `json:"family"`
	Standard07ResidualVariant       string  `json:"variant"`
	Standard07ResidualPrimary       uint8   `json:"primary"`
	Standard07ResidualSecondary     *uint8  `json:"secondary"`
	Standard07ResidualOperation     *uint8  `json:"operation"`
	Standard07ResidualDirection     string  `json:"direction"`
	Standard07ResidualField         string  `json:"field"`
	Standard07ResidualAlternative   string  `json:"alternative"`
	Standard07ResidualGroup         *uint8  `json:"group"`
	Standard07ResidualInstanceScope string  `json:"instance_scope"`
	Standard07ResidualRegister      *uint16 `json:"register"`
}

//go:embed testdata/ab28_standard07_residual_fields.json
var standard07ResidualEmbedded string

func standard07ResidualValidate(standard07ResidualData []byte) ([]standard07ResidualJSON.RawMessage, error) {
	var standard07ResidualRows []standard07ResidualJSON.RawMessage
	if standard07ResidualErr := standard07ResidualJSON.Unmarshal(standard07ResidualData, &standard07ResidualRows); standard07ResidualErr != nil {
		return nil, standard07ResidualErr
	}
	if standard07ResidualRows == nil {
		return nil, standard07ResidualFmt.Errorf("ab28: expected JSON array")
	}
	standard07ResidualSeen := make(map[string]bool)
	for _, standard07ResidualRow := range standard07ResidualRows {
		var standard07ResidualEnvelope struct {
			Standard07ResidualKey *Standard07ResidualKey `json:"key"`
		}
		if standard07ResidualErr := standard07ResidualJSON.Unmarshal(standard07ResidualRow, &standard07ResidualEnvelope); standard07ResidualErr != nil {
			return nil, standard07ResidualErr
		}
		if standard07ResidualEnvelope.Standard07ResidualKey == nil {
			return nil, standard07ResidualFmt.Errorf("ab28: missing key")
		}
		var standard07ResidualObject map[string]standard07ResidualJSON.RawMessage
		if standard07ResidualErr := standard07ResidualJSON.Unmarshal(standard07ResidualRow, &standard07ResidualObject); standard07ResidualErr != nil {
			return nil, standard07ResidualErr
		}
		var standard07ResidualAxes map[string]standard07ResidualJSON.RawMessage
		if standard07ResidualErr := standard07ResidualJSON.Unmarshal(standard07ResidualObject["key"], &standard07ResidualAxes); standard07ResidualErr != nil {
			return nil, standard07ResidualErr
		}
		for _, standard07ResidualAxis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			if _, standard07ResidualOK := standard07ResidualAxes[standard07ResidualAxis]; !standard07ResidualOK {
				return nil, standard07ResidualFmt.Errorf("ab28: missing axis %s", standard07ResidualAxis)
			}
		}
		standard07ResidualKey := *standard07ResidualEnvelope.Standard07ResidualKey
		if standard07ResidualKey.Standard07ResidualFamily == "" || standard07ResidualKey.Standard07ResidualVariant == "" || standard07ResidualKey.Standard07ResidualDirection == "" || standard07ResidualKey.Standard07ResidualField == "" || standard07ResidualKey.Standard07ResidualAlternative == "" || standard07ResidualKey.Standard07ResidualInstanceScope == "" || string(standard07ResidualAxes["primary"]) == "null" {
			return nil, standard07ResidualFmt.Errorf("ab28: incomplete identity")
		}
		standard07ResidualIdentity, standard07ResidualErr := standard07ResidualJSON.Marshal(standard07ResidualKey)
		if standard07ResidualErr != nil {
			return nil, standard07ResidualErr
		}
		if standard07ResidualSeen[string(standard07ResidualIdentity)] {
			return nil, standard07ResidualFmt.Errorf("ab28: duplicate key")
		}
		standard07ResidualSeen[string(standard07ResidualIdentity)] = true
	}
	return standard07ResidualRows, nil
}

func Standard07ResidualEntries() (standard07ResidualJSON.RawMessage, error) {
	if _, standard07ResidualErr := standard07ResidualValidate([]byte(standard07ResidualEmbedded)); standard07ResidualErr != nil {
		return nil, standard07ResidualErr
	}
	return standard07ResidualJSON.RawMessage([]byte(standard07ResidualEmbedded)), nil
}

func standard07ResidualLookup(standard07ResidualData []byte, standard07ResidualKey Standard07ResidualKey) (standard07ResidualJSON.RawMessage, bool, error) {
	standard07ResidualRows, standard07ResidualErr := standard07ResidualValidate(standard07ResidualData)
	if standard07ResidualErr != nil {
		return nil, false, standard07ResidualErr
	}
	standard07ResidualWanted, standard07ResidualErr := standard07ResidualJSON.Marshal(standard07ResidualKey)
	if standard07ResidualErr != nil {
		return nil, false, standard07ResidualErr
	}
	for _, standard07ResidualRow := range standard07ResidualRows {
		var standard07ResidualEnvelope struct {
			Standard07ResidualKey Standard07ResidualKey `json:"key"`
		}
		if standard07ResidualErr := standard07ResidualJSON.Unmarshal(standard07ResidualRow, &standard07ResidualEnvelope); standard07ResidualErr != nil {
			return nil, false, standard07ResidualErr
		}
		standard07ResidualIdentity, standard07ResidualErr := standard07ResidualJSON.Marshal(standard07ResidualEnvelope.Standard07ResidualKey)
		if standard07ResidualErr != nil {
			return nil, false, standard07ResidualErr
		}
		if string(standard07ResidualIdentity) == string(standard07ResidualWanted) {
			return append(standard07ResidualJSON.RawMessage(nil), standard07ResidualRow...), true, nil
		}
	}
	return nil, false, nil
}

func Standard07ResidualLookup(standard07ResidualKey Standard07ResidualKey) (standard07ResidualJSON.RawMessage, bool, error) {
	return standard07ResidualLookup([]byte(standard07ResidualEmbedded), standard07ResidualKey)
}
