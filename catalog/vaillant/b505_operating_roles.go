package vaillant

import (
	_ "embed"
	b505OperatingRoleJSON "encoding/json"
	b505OperatingRoleFmt "fmt"
)

//go:embed data/ab16_b505_operating_roles.json
var b505OperatingRoleData string

type B505OperatingRoleKey struct {
	B505OperatingRoleFamily        string  `json:"Ab16Family"`
	B505OperatingRoleVariant       string  `json:"Ab16Variant"`
	B505OperatingRolePrimary       uint8   `json:"Ab16Primary"`
	B505OperatingRoleSecondary     *uint8  `json:"Ab16Secondary"`
	B505OperatingRoleOperation     *uint8  `json:"Ab16Operation"`
	B505OperatingRoleDirection     string  `json:"Ab16Direction"`
	B505OperatingRoleField         string  `json:"Ab16Field"`
	B505OperatingRoleAlternative   string  `json:"Ab16Alternative"`
	B505OperatingRoleGroup         *uint8  `json:"Ab16Group"`
	B505OperatingRoleInstanceScope string  `json:"Ab16InstanceScope"`
	B505OperatingRoleRegister      *uint16 `json:"Ab16Register"`
}

func b505OperatingRoleValidate(b505OperatingRoleData []byte) (map[string]b505OperatingRoleJSON.RawMessage, error) {
	var b505OperatingRoleRows []b505OperatingRoleJSON.RawMessage
	if b505OperatingRoleErr := b505OperatingRoleJSON.Unmarshal(b505OperatingRoleData, &b505OperatingRoleRows); b505OperatingRoleErr != nil {
		return nil, b505OperatingRoleErr
	}
	if b505OperatingRoleRows == nil {
		return nil, b505OperatingRoleFmt.Errorf("ab16: expected array")
	}
	b505OperatingRoleIndex := make(map[string]b505OperatingRoleJSON.RawMessage, len(b505OperatingRoleRows))
	for _, b505OperatingRoleRow := range b505OperatingRoleRows {
		var b505OperatingRoleEnvelope struct {
			B505OperatingRoleIdentity *B505OperatingRoleKey `json:"identity"`
		}
		if b505OperatingRoleErr := b505OperatingRoleJSON.Unmarshal(b505OperatingRoleRow, &b505OperatingRoleEnvelope); b505OperatingRoleErr != nil {
			return nil, b505OperatingRoleErr
		}
		if b505OperatingRoleEnvelope.B505OperatingRoleIdentity == nil {
			return nil, b505OperatingRoleFmt.Errorf("ab16: missing identity")
		}
		b505OperatingRoleKey := b505OperatingRoleEnvelope.B505OperatingRoleIdentity
		if b505OperatingRoleKey.B505OperatingRoleFamily == "" || b505OperatingRoleKey.B505OperatingRoleVariant == "" || b505OperatingRoleKey.B505OperatingRoleDirection == "" || b505OperatingRoleKey.B505OperatingRoleField == "" || b505OperatingRoleKey.B505OperatingRoleAlternative == "" || b505OperatingRoleKey.B505OperatingRoleInstanceScope == "" {
			return nil, b505OperatingRoleFmt.Errorf("ab16: incomplete string axes")
		}
		b505OperatingRoleEncoded, b505OperatingRoleErr := b505OperatingRoleJSON.Marshal(b505OperatingRoleKey)
		if b505OperatingRoleErr != nil {
			return nil, b505OperatingRoleErr
		}
		if _, b505OperatingRoleFound := b505OperatingRoleIndex[string(b505OperatingRoleEncoded)]; b505OperatingRoleFound {
			return nil, b505OperatingRoleFmt.Errorf("ab16: duplicate identity")
		}
		b505OperatingRoleIndex[string(b505OperatingRoleEncoded)] = append(b505OperatingRoleJSON.RawMessage(nil), b505OperatingRoleRow...)
	}
	return b505OperatingRoleIndex, nil
}

func B505OperatingRoleEntries() (b505OperatingRoleJSON.RawMessage, error) {
	if _, b505OperatingRoleErr := b505OperatingRoleValidate([]byte(b505OperatingRoleData)); b505OperatingRoleErr != nil {
		return nil, b505OperatingRoleErr
	}
	return b505OperatingRoleJSON.RawMessage(b505OperatingRoleData), nil
}

func b505OperatingRoleLookup(b505OperatingRoleData []byte, b505OperatingRoleKey B505OperatingRoleKey) (b505OperatingRoleJSON.RawMessage, bool, error) {
	b505OperatingRoleIndex, b505OperatingRoleErr := b505OperatingRoleValidate(b505OperatingRoleData)
	if b505OperatingRoleErr != nil {
		return nil, false, b505OperatingRoleErr
	}
	b505OperatingRoleEncoded, b505OperatingRoleErr := b505OperatingRoleJSON.Marshal(b505OperatingRoleKey)
	if b505OperatingRoleErr != nil {
		return nil, false, b505OperatingRoleErr
	}
	b505OperatingRoleRow, b505OperatingRoleFound := b505OperatingRoleIndex[string(b505OperatingRoleEncoded)]
	return append(b505OperatingRoleJSON.RawMessage(nil), b505OperatingRoleRow...), b505OperatingRoleFound, nil
}

func B505OperatingRoleLookup(b505OperatingRoleKey B505OperatingRoleKey) (b505OperatingRoleJSON.RawMessage, bool, error) {
	return b505OperatingRoleLookup([]byte(b505OperatingRoleData), b505OperatingRoleKey)
}
