package ebusstandard

import (
	standardFFBodyStateBytes "bytes"
	_ "embed"
	standardFFBodyStateJSON "encoding/json"
	standardFFBodyStateFmt "fmt"
)

type StandardFFBodyStateKey struct {
	StandardFFBodyStateFamily        string  `json:"family"`
	StandardFFBodyStateVariant       string  `json:"variant"`
	StandardFFBodyStatePrimary       uint8   `json:"primary"`
	StandardFFBodyStateSecondary     *uint8  `json:"secondary"`
	StandardFFBodyStateOperation     *uint8  `json:"operation"`
	StandardFFBodyStateDirection     string  `json:"direction"`
	StandardFFBodyStateField         string  `json:"field"`
	StandardFFBodyStateAlternative   string  `json:"alternative"`
	StandardFFBodyStateGroup         *uint16 `json:"group"`
	StandardFFBodyStateInstanceScope string  `json:"instance_scope"`
	StandardFFBodyStateRegister      *uint16 `json:"register"`
}

//go:embed testdata/ab54_standard_ff_body_states.json
var standardFFBodyStateBodyStatesJSON []byte

func standardFFBodyStateUniqueFields(standardFFBodyStateDecoder *standardFFBodyStateJSON.Decoder) error {
	standardFFBodyStateToken, standardFFBodyStateErr := standardFFBodyStateDecoder.Token()
	if standardFFBodyStateErr != nil {
		return standardFFBodyStateErr
	}
	standardFFBodyStateDelimiter, standardFFBodyStateComposite := standardFFBodyStateToken.(standardFFBodyStateJSON.Delim)
	if !standardFFBodyStateComposite {
		return nil
	}
	standardFFBodyStateSeen := make(map[string]bool)
	for standardFFBodyStateDecoder.More() {
		if standardFFBodyStateDelimiter == '{' {
			standardFFBodyStateNameToken, standardFFBodyStateErr := standardFFBodyStateDecoder.Token()
			if standardFFBodyStateErr != nil {
				return standardFFBodyStateErr
			}
			standardFFBodyStateName, standardFFBodyStateOK := standardFFBodyStateNameToken.(string)
			if !standardFFBodyStateOK || standardFFBodyStateSeen[standardFFBodyStateName] {
				return standardFFBodyStateFmt.Errorf("standardff body states: duplicate JSON member")
			}
			standardFFBodyStateSeen[standardFFBodyStateName] = true
		}
		if standardFFBodyStateErr := standardFFBodyStateUniqueFields(standardFFBodyStateDecoder); standardFFBodyStateErr != nil {
			return standardFFBodyStateErr
		}
	}
	_, standardFFBodyStateErr = standardFFBodyStateDecoder.Token()
	return standardFFBodyStateErr
}

func standardFFBodyStateValidate(standardFFBodyStateData []byte) (standardFFBodyStateJSON.RawMessage, map[string]standardFFBodyStateJSON.RawMessage, error) {
	standardFFBodyStateData = standardFFBodyStateBytes.TrimSpace(standardFFBodyStateData)
	if len(standardFFBodyStateData) == 0 || standardFFBodyStateData[0] != '[' {
		return nil, nil, standardFFBodyStateFmt.Errorf("standardff body states: expected JSON array")
	}
	var standardFFBodyStateRows []standardFFBodyStateJSON.RawMessage
	if standardFFBodyStateErr := standardFFBodyStateJSON.Unmarshal(standardFFBodyStateData, &standardFFBodyStateRows); standardFFBodyStateErr != nil {
		return nil, nil, standardFFBodyStateErr
	}
	if standardFFBodyStateErr := standardFFBodyStateUniqueFields(standardFFBodyStateJSON.NewDecoder(standardFFBodyStateBytes.NewReader(standardFFBodyStateData))); standardFFBodyStateErr != nil {
		return nil, nil, standardFFBodyStateErr
	}
	standardFFBodyStateIndex := make(map[string]standardFFBodyStateJSON.RawMessage, len(standardFFBodyStateRows))
	for _, standardFFBodyStateRow := range standardFFBodyStateRows {
		var standardFFBodyStateEnvelope map[string]standardFFBodyStateJSON.RawMessage
		if standardFFBodyStateErr := standardFFBodyStateJSON.Unmarshal(standardFFBodyStateRow, &standardFFBodyStateEnvelope); standardFFBodyStateErr != nil {
			return nil, nil, standardFFBodyStateErr
		}
		var standardFFBodyStateKey StandardFFBodyStateKey
		standardFFBodyStateDecoder := standardFFBodyStateJSON.NewDecoder(standardFFBodyStateBytes.NewReader(standardFFBodyStateEnvelope["key"]))
		standardFFBodyStateDecoder.DisallowUnknownFields()
		if standardFFBodyStateErr := standardFFBodyStateDecoder.Decode(&standardFFBodyStateKey); standardFFBodyStateErr != nil {
			return nil, nil, standardFFBodyStateErr
		}
		var standardFFBodyStateFields map[string]standardFFBodyStateJSON.RawMessage
		if standardFFBodyStateErr := standardFFBodyStateJSON.Unmarshal(standardFFBodyStateEnvelope["key"], &standardFFBodyStateFields); standardFFBodyStateErr != nil {
			return nil, nil, standardFFBodyStateErr
		}
		for standardFFBodyStateName := range standardFFBodyStateFields {
			switch standardFFBodyStateName {
			case "family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register":
			default:
				return nil, nil, standardFFBodyStateFmt.Errorf("standardff body states: unknown key axis")
			}
		}
		if len(standardFFBodyStateFields["primary"]) == 0 || standardFFBodyStateBytes.Equal(standardFFBodyStateFields["primary"], []byte("null")) ||
			len(standardFFBodyStateFields["secondary"]) == 0 || len(standardFFBodyStateFields["operation"]) == 0 ||
			len(standardFFBodyStateFields["group"]) == 0 || len(standardFFBodyStateFields["register"]) == 0 ||
			standardFFBodyStateKey.StandardFFBodyStateFamily == "" || standardFFBodyStateKey.StandardFFBodyStateVariant == "" || standardFFBodyStateKey.StandardFFBodyStateDirection == "" ||
			standardFFBodyStateKey.StandardFFBodyStateField == "" || standardFFBodyStateKey.StandardFFBodyStateAlternative == "" || standardFFBodyStateKey.StandardFFBodyStateInstanceScope == "" {
			return nil, nil, standardFFBodyStateFmt.Errorf("standardff body states: incomplete key; use none for non-applicable string axes")
		}
		standardFFBodyStateIdentity, standardFFBodyStateErr := standardFFBodyStateJSON.Marshal(standardFFBodyStateKey)
		if standardFFBodyStateErr != nil {
			return nil, nil, standardFFBodyStateErr
		}
		if _, standardFFBodyStateExists := standardFFBodyStateIndex[string(standardFFBodyStateIdentity)]; standardFFBodyStateExists {
			return nil, nil, standardFFBodyStateFmt.Errorf("standardff body states: duplicate key")
		}
		standardFFBodyStateIndex[string(standardFFBodyStateIdentity)] = append(standardFFBodyStateJSON.RawMessage(nil), standardFFBodyStateRow...)
	}
	return append(standardFFBodyStateJSON.RawMessage(nil), standardFFBodyStateData...), standardFFBodyStateIndex, nil
}

func StandardFFBodyStateEntries() (standardFFBodyStateJSON.RawMessage, error) {
	standardFFBodyStateEntries, _, standardFFBodyStateErr := standardFFBodyStateValidate(standardFFBodyStateBodyStatesJSON)
	return standardFFBodyStateEntries, standardFFBodyStateErr
}

func standardFFBodyStateLookup(standardFFBodyStateData []byte, standardFFBodyStateKey StandardFFBodyStateKey) (standardFFBodyStateJSON.RawMessage, bool, error) {
	_, standardFFBodyStateIndex, standardFFBodyStateErr := standardFFBodyStateValidate(standardFFBodyStateData)
	if standardFFBodyStateErr != nil {
		return nil, false, standardFFBodyStateErr
	}
	standardFFBodyStateIdentity, standardFFBodyStateErr := standardFFBodyStateJSON.Marshal(standardFFBodyStateKey)
	if standardFFBodyStateErr != nil {
		return nil, false, standardFFBodyStateErr
	}
	standardFFBodyStateRow, standardFFBodyStateFound := standardFFBodyStateIndex[string(standardFFBodyStateIdentity)]
	return append(standardFFBodyStateJSON.RawMessage(nil), standardFFBodyStateRow...), standardFFBodyStateFound, nil
}

func StandardFFBodyStateLookup(standardFFBodyStateKey StandardFFBodyStateKey) (standardFFBodyStateJSON.RawMessage, bool, error) {
	return standardFFBodyStateLookup(standardFFBodyStateBodyStatesJSON, standardFFBodyStateKey)
}
