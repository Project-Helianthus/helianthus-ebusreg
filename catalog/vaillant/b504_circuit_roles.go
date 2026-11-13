package vaillant

import (
	b504CircuitRoleEmbed "embed"
	b504CircuitRoleJSON "encoding/json"
	b504CircuitRoleFmt "fmt"
	b504CircuitRoleReflect "reflect"
	b504CircuitRoleStrings "strings"
)

//go:embed data/ab15_b504_circuit_roles.json
var b504CircuitRoleFiles b504CircuitRoleEmbed.FS

type B504CircuitRoleKey struct {
	B504CircuitRoleFamily        string  `json:"family"`
	B504CircuitRoleVariant       string  `json:"variant"`
	B504CircuitRolePrimary       uint8   `json:"primary"`
	B504CircuitRoleSecondary     *uint8  `json:"secondary"`
	B504CircuitRoleOperation     *uint8  `json:"operation"`
	B504CircuitRoleDirection     string  `json:"direction"`
	B504CircuitRoleField         string  `json:"field"`
	B504CircuitRoleAlternative   string  `json:"alternative"`
	B504CircuitRoleGroup         *uint8  `json:"group"`
	B504CircuitRoleInstanceScope string  `json:"instance_scope"`
	B504CircuitRoleRegister      *uint16 `json:"register"`
}

type b504CircuitRoleRow struct {
	B504CircuitRoleID                    string                         `json:"id"`
	B504CircuitRoleKey                   B504CircuitRoleKey             `json:"key"`
	B504CircuitRoleGrade                 string                         `json:"grade"`
	B504CircuitRoleMetadataPresence      string                         `json:"metadata_presence"`
	B504CircuitRoleNativeContract        string                         `json:"native_contract"`
	B504CircuitRolePhysicalQualification string                         `json:"physical_qualification"`
	B504CircuitRole                      string                         `json:"role"`
	B504CircuitRoleAppliesTo             []string                       `json:"applies_to"`
	B504CircuitRoleRelations             []string                       `json:"relations"`
	B504CircuitRoleCondition             string                         `json:"condition"`
	B504CircuitRoleBit                   *uint8                         `json:"bit"`
	B504CircuitRoleUnknown               []string                       `json:"unknown"`
	B504CircuitRoleRaw                   b504CircuitRoleJSON.RawMessage `json:"raw"`
}

func B504CircuitRoleEntries() (b504CircuitRoleJSON.RawMessage, error) {
	b504CircuitRoleData, b504CircuitRoleErr := b504CircuitRoleFiles.ReadFile("data/ab15_b504_circuit_roles.json")
	if b504CircuitRoleErr != nil {
		return nil, b504CircuitRoleErr
	}
	return b504CircuitRoleEntries(b504CircuitRoleData)
}

func B504CircuitRoleLookup(b504CircuitRoleKey B504CircuitRoleKey) (b504CircuitRoleJSON.RawMessage, bool, error) {
	b504CircuitRoleData, b504CircuitRoleErr := b504CircuitRoleFiles.ReadFile("data/ab15_b504_circuit_roles.json")
	if b504CircuitRoleErr != nil {
		return nil, false, b504CircuitRoleErr
	}
	return b504CircuitRoleLookup(b504CircuitRoleData, b504CircuitRoleKey)
}

func b504CircuitRoleEntries(b504CircuitRoleData []byte) (b504CircuitRoleJSON.RawMessage, error) {
	var b504CircuitRoleRows []b504CircuitRoleJSON.RawMessage
	if b504CircuitRoleErr := b504CircuitRoleJSON.Unmarshal(b504CircuitRoleData, &b504CircuitRoleRows); b504CircuitRoleErr != nil {
		return nil, b504CircuitRoleErr
	}
	if len(b504CircuitRoleRows) == 0 {
		return nil, b504CircuitRoleFmt.Errorf("ab15: empty catalog")
	}
	b504CircuitRoleKeys := make([]B504CircuitRoleKey, 0, len(b504CircuitRoleRows))
	b504CircuitRoleIDs := map[string]bool{}
	for _, b504CircuitRoleRaw := range b504CircuitRoleRows {
		var b504CircuitRoleEntry b504CircuitRoleRow
		if b504CircuitRoleErr := b504CircuitRoleJSON.Unmarshal(b504CircuitRoleRaw, &b504CircuitRoleEntry); b504CircuitRoleErr != nil {
			return nil, b504CircuitRoleErr
		}
		var b504CircuitRoleObject map[string]b504CircuitRoleJSON.RawMessage
		if b504CircuitRoleErr := b504CircuitRoleJSON.Unmarshal(b504CircuitRoleRaw, &b504CircuitRoleObject); b504CircuitRoleErr != nil {
			return nil, b504CircuitRoleErr
		}
		var b504CircuitRoleIdentity map[string]b504CircuitRoleJSON.RawMessage
		if b504CircuitRoleErr := b504CircuitRoleJSON.Unmarshal(b504CircuitRoleObject["key"], &b504CircuitRoleIdentity); b504CircuitRoleErr != nil {
			return nil, b504CircuitRoleErr
		}
		for _, b504CircuitRoleAxis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
			b504CircuitRoleValue, b504CircuitRolePresent := b504CircuitRoleIdentity[b504CircuitRoleAxis]
			if !b504CircuitRolePresent || (string(b504CircuitRoleValue) == "null" && b504CircuitRoleAxis != "secondary" && b504CircuitRoleAxis != "operation" && b504CircuitRoleAxis != "group" && b504CircuitRoleAxis != "register") {
				return nil, b504CircuitRoleFmt.Errorf("ab15: missing identity axis %s", b504CircuitRoleAxis)
			}
		}
		if len(b504CircuitRoleIdentity) != 11 || b504CircuitRoleEntry.B504CircuitRoleKey.B504CircuitRoleFamily == "" || b504CircuitRoleEntry.B504CircuitRoleKey.B504CircuitRoleVariant == "" || b504CircuitRoleEntry.B504CircuitRoleKey.B504CircuitRoleDirection == "" || b504CircuitRoleEntry.B504CircuitRoleKey.B504CircuitRoleField == "" || b504CircuitRoleEntry.B504CircuitRoleKey.B504CircuitRoleAlternative == "" || b504CircuitRoleEntry.B504CircuitRoleKey.B504CircuitRoleInstanceScope == "" {
			return nil, b504CircuitRoleFmt.Errorf("ab15: incomplete identity")
		}
		if !b504CircuitRoleStrings.HasPrefix(b504CircuitRoleEntry.B504CircuitRoleID, "ab15-") || b504CircuitRoleIDs[b504CircuitRoleEntry.B504CircuitRoleID] {
			return nil, b504CircuitRoleFmt.Errorf("ab15: invalid or duplicate id")
		}
		b504CircuitRoleIDs[b504CircuitRoleEntry.B504CircuitRoleID] = true
		if b504CircuitRoleEntry.B504CircuitRoleGrade != "Hypothesis" || b504CircuitRoleEntry.B504CircuitRoleMetadataPresence != "Proven" || b504CircuitRoleEntry.B504CircuitRolePhysicalQualification != "Unknown" || b504CircuitRoleEntry.B504CircuitRoleNativeContract == "" || b504CircuitRoleEntry.B504CircuitRole == "" || len(b504CircuitRoleEntry.B504CircuitRoleAppliesTo) == 0 || len(b504CircuitRoleEntry.B504CircuitRoleRelations) == 0 || len(b504CircuitRoleEntry.B504CircuitRoleUnknown) == 0 || len(b504CircuitRoleEntry.B504CircuitRoleRaw) == 0 || string(b504CircuitRoleEntry.B504CircuitRoleRaw) == "null" {
			return nil, b504CircuitRoleFmt.Errorf("ab15: incomplete evidence boundary")
		}
		for _, b504CircuitRolePrevious := range b504CircuitRoleKeys {
			if b504CircuitRoleReflect.DeepEqual(b504CircuitRolePrevious, b504CircuitRoleEntry.B504CircuitRoleKey) {
				return nil, b504CircuitRoleFmt.Errorf("ab15: duplicate key")
			}
		}
		b504CircuitRoleKeys = append(b504CircuitRoleKeys, b504CircuitRoleEntry.B504CircuitRoleKey)
	}
	return append(b504CircuitRoleJSON.RawMessage(nil), b504CircuitRoleData...), nil
}

func b504CircuitRoleLookup(b504CircuitRoleData []byte, b504CircuitRoleKey B504CircuitRoleKey) (b504CircuitRoleJSON.RawMessage, bool, error) {
	b504CircuitRoleValidated, b504CircuitRoleErr := b504CircuitRoleEntries(b504CircuitRoleData)
	if b504CircuitRoleErr != nil {
		return nil, false, b504CircuitRoleErr
	}
	var b504CircuitRoleRows []b504CircuitRoleJSON.RawMessage
	if b504CircuitRoleErr = b504CircuitRoleJSON.Unmarshal(b504CircuitRoleValidated, &b504CircuitRoleRows); b504CircuitRoleErr != nil {
		return nil, false, b504CircuitRoleErr
	}
	for _, b504CircuitRoleRaw := range b504CircuitRoleRows {
		var b504CircuitRoleEntry b504CircuitRoleRow
		if b504CircuitRoleErr = b504CircuitRoleJSON.Unmarshal(b504CircuitRoleRaw, &b504CircuitRoleEntry); b504CircuitRoleErr != nil {
			return nil, false, b504CircuitRoleErr
		}
		if b504CircuitRoleReflect.DeepEqual(b504CircuitRoleKey, b504CircuitRoleEntry.B504CircuitRoleKey) {
			return append(b504CircuitRoleJSON.RawMessage(nil), b504CircuitRoleRaw...), true, nil
		}
	}
	return nil, false, nil
}
