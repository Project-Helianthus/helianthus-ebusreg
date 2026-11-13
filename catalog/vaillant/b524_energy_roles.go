package vaillant

import (
	b524EnergyRoleBytesPackage "bytes"
	_ "embed"
	b524EnergyRoleJSON "encoding/json"
	b524EnergyRoleFmt "fmt"
	b524EnergyRoleIO "io"
	b524EnergyRoleReflect "reflect"
)

//go:embed data/ab95_b524_energy_roles.json
var b524EnergyRoleData []byte

type B524EnergyRoleKey struct {
	B524EnergyRoleFamily            string  `json:"family"`
	B524EnergyRoleVariant           string  `json:"variant"`
	B524EnergyRolePrimary           uint8   `json:"primary"`
	B524EnergyRoleSecondary         *uint8  `json:"secondary"`
	B524EnergyRoleOperation         *uint8  `json:"operation"`
	B524EnergyRoleDirection         string  `json:"direction"`
	B524EnergyRoleField             string  `json:"field"`
	B524EnergyRoleGroup             *uint8  `json:"group"`
	B524EnergyRoleRequestedInstance *uint8  `json:"requested_instance"`
	B524EnergyRoleRegister          *uint16 `json:"register"`
	B524EnergyCommandRole           string  `json:"command_role"`
}

type B524EnergyRoleEntry struct {
	B524EnergyRoleReadbackInstance      *uint8            `json:"readback_instance"`
	B524EnergyRoleKey                   B524EnergyRoleKey `json:"key"`
	B524EnergyRoleGrade                 string            `json:"grade"`
	B524EnergyRoleCodec                 *string           `json:"codec"`
	B524EnergyRoleWidthBytes            *uint8            `json:"width_bytes"`
	B524EnergyRoleFieldOffset           *uint8            `json:"field_offset"`
	B524EnergyRoleByteOrder             *string           `json:"byte_order"`
	B524EnergyRoleNativeUnit            *string           `json:"native_unit"`
	B524EnergyRoleControlEnumeration    map[string]string `json:"control_enumeration"`
	B524EnergyRoleResponseForm          *string           `json:"response_form"`
	B524EnergyRoleUnknown               []string          `json:"unknown"`
	B524EnergyRoleRawPreservation       string            `json:"raw_preservation"`
	B524EnergyRoleMetadataBoundaryGrade string            `json:"metadata_boundary_grade"`
	B524EnergyRoleMetadataBoundary      string            `json:"metadata_boundary"`
	B524EnergyRolePhysicalQualification string            `json:"physical_qualification"`
	B524EnergyRoleNativeContract        string            `json:"native_contract"`
}

func B524EnergyRoleEntries() ([]B524EnergyRoleEntry, error) {
	var b524EnergyRoleEntries []B524EnergyRoleEntry
	b524EnergyRoleDecoder := b524EnergyRoleJSON.NewDecoder(b524EnergyRoleBytesPackage.NewReader(b524EnergyRoleData))
	b524EnergyRoleDecoder.DisallowUnknownFields()
	if b524EnergyRoleError := b524EnergyRoleDecoder.Decode(&b524EnergyRoleEntries); b524EnergyRoleError != nil {
		return nil, b524EnergyRoleFmt.Errorf("ab95 metadata: %w", b524EnergyRoleError)
	}
	var b524EnergyRoleTrailing any
	if b524EnergyRoleError := b524EnergyRoleDecoder.Decode(&b524EnergyRoleTrailing); b524EnergyRoleError != b524EnergyRoleIO.EOF {
		return nil, b524EnergyRoleFmt.Errorf("ab95 metadata: trailing data")
	}
	if len(b524EnergyRoleEntries) != 5 {
		return nil, b524EnergyRoleFmt.Errorf("ab95 metadata: expected five targets")
	}
	for b524EnergyRoleIndex, b524EnergyRoleEntry := range b524EnergyRoleEntries {
		if b524EnergyRoleEntry.B524EnergyRoleKey.B524EnergyRoleFamily == "" || b524EnergyRoleEntry.B524EnergyRoleKey.B524EnergyRoleVariant == "" || b524EnergyRoleEntry.B524EnergyRoleKey.B524EnergyRoleDirection == "" || b524EnergyRoleEntry.B524EnergyRoleKey.B524EnergyRoleField == "" || b524EnergyRoleEntry.B524EnergyRoleKey.B524EnergyCommandRole == "" || b524EnergyRoleEntry.B524EnergyRoleGrade != "Hypothesis" || b524EnergyRoleEntry.B524EnergyRoleMetadataBoundaryGrade != "Proven" || b524EnergyRoleEntry.B524EnergyRolePhysicalQualification != "Unknown" || b524EnergyRoleEntry.B524EnergyRoleNativeContract == "" || b524EnergyRoleEntry.B524EnergyRoleMetadataBoundary == "" || b524EnergyRoleEntry.B524EnergyRoleRawPreservation == "" || len(b524EnergyRoleEntry.B524EnergyRoleUnknown) == 0 {
			return nil, b524EnergyRoleFmt.Errorf("ab95 metadata: invalid entry %d", b524EnergyRoleIndex)
		}
		for _, b524EnergyRolePrevious := range b524EnergyRoleEntries[:b524EnergyRoleIndex] {
			if b524EnergyRoleReflect.DeepEqual(b524EnergyRoleEntry.B524EnergyRoleKey, b524EnergyRolePrevious.B524EnergyRoleKey) {
				return nil, b524EnergyRoleFmt.Errorf("ab95 metadata: duplicate key at %d", b524EnergyRoleIndex)
			}
		}
	}
	return b524EnergyRoleEntries, nil
}

func B524EnergyRoleLookup(b524EnergyRoleKey B524EnergyRoleKey) (B524EnergyRoleEntry, bool, error) {
	b524EnergyRoleEntries, b524EnergyRoleError := B524EnergyRoleEntries()
	if b524EnergyRoleError != nil {
		return B524EnergyRoleEntry{}, false, b524EnergyRoleError
	}
	for _, b524EnergyRoleEntry := range b524EnergyRoleEntries {
		if b524EnergyRoleReflect.DeepEqual(b524EnergyRoleKey, b524EnergyRoleEntry.B524EnergyRoleKey) {
			return b524EnergyRoleEntry, true, nil
		}
	}
	return B524EnergyRoleEntry{}, false, nil
}
