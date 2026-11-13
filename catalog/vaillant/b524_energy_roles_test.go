package vaillant

import (
	b524EnergyRoleJSON "encoding/json"
	b524EnergyRoleReflect "reflect"
	b524EnergyRoleTesting "testing"
)

func b524EnergyRolePointer[b524EnergyRoleT any](b524EnergyRoleValue b524EnergyRoleT) *b524EnergyRoleT {
	return &b524EnergyRoleValue
}

func b524EnergyRoleExpectedKeys() []B524EnergyRoleKey {
	return []B524EnergyRoleKey{
		{"vaillant", "B524", 0xb5, b524EnergyRolePointer(uint8(0x24)), b524EnergyRolePointer(uint8(2)), "request", "setpoint", b524EnergyRolePointer(uint8(1)), b524EnergyRolePointer(uint8(0xff)), b524EnergyRolePointer(uint16(4)), "setpoint"},
		{"vaillant", "B524", 0xb5, b524EnergyRolePointer(uint8(0x24)), b524EnergyRolePointer(uint8(2)), "request", "hysteresis", b524EnergyRolePointer(uint8(0)), b524EnergyRolePointer(uint8(0xff)), b524EnergyRolePointer(uint16(0x27)), "hysteresis"},
		{"vaillant", "B524", 0xb5, b524EnergyRolePointer(uint8(0x24)), b524EnergyRolePointer(uint8(2)), "request", "control", b524EnergyRolePointer(uint8(1)), b524EnergyRolePointer(uint8(0xff)), b524EnergyRolePointer(uint16(0xd)), "control"},
		{"vaillant", "B524", 0xb5, b524EnergyRolePointer(uint8(0x24)), b524EnergyRolePointer(uint8(2)), "request", "offset", b524EnergyRolePointer(uint8(0)), b524EnergyRolePointer(uint8(0xff)), b524EnergyRolePointer(uint16(0xa2)), "offset"},
		{"vaillant", "B524", 0xb5, b524EnergyRolePointer(uint8(0x24)), b524EnergyRolePointer(uint8(2)), "request", "activation", b524EnergyRolePointer(uint8(0)), b524EnergyRolePointer(uint8(0xff)), b524EnergyRolePointer(uint16(0xb6)), "activation"},
	}
}

func TestAb95FiveTargets(b524EnergyRoleT *b524EnergyRoleTesting.T) {
	b524EnergyRoleEntries, b524EnergyRoleError := B524EnergyRoleEntries()
	if b524EnergyRoleError != nil || len(b524EnergyRoleEntries) != 5 {
		b524EnergyRoleT.Fatalf("entries: %v, %d", b524EnergyRoleError, len(b524EnergyRoleEntries))
	}
	for b524EnergyRoleIndex, b524EnergyRoleKey := range b524EnergyRoleExpectedKeys() {
		b524EnergyRoleEntry, b524EnergyRoleOK, b524EnergyRoleError := B524EnergyRoleLookup(b524EnergyRoleKey)
		if b524EnergyRoleError != nil || !b524EnergyRoleOK || !b524EnergyRoleReflect.DeepEqual(b524EnergyRoleEntry.B524EnergyRoleKey, b524EnergyRoleKey) {
			b524EnergyRoleT.Fatalf("target %d: %v %v", b524EnergyRoleIndex, b524EnergyRoleOK, b524EnergyRoleError)
		}
		if b524EnergyRoleEntry.B524EnergyRoleReadbackInstance == nil || *b524EnergyRoleEntry.B524EnergyRoleReadbackInstance != 0 {
			b524EnergyRoleT.Fatal("readback context must remain separate")
		}
		b524EnergyRoleWidth := []*uint8{b524EnergyRolePointer(uint8(4)), b524EnergyRolePointer(uint8(4)), nil, b524EnergyRolePointer(uint8(4)), b524EnergyRolePointer(uint8(1))}[b524EnergyRoleIndex]
		b524EnergyRoleCodec := []*string{b524EnergyRolePointer("binary32"), b524EnergyRolePointer("binary32"), nil, b524EnergyRolePointer("binary32"), nil}[b524EnergyRoleIndex]
		b524EnergyRoleUnknown := []string{"binary32 byte order", "field offsets", "control enumeration", "native units", "response forms"}
		if b524EnergyRoleCodec == nil {
			b524EnergyRoleUnknown = append(b524EnergyRoleUnknown, "codec")
		}
		if b524EnergyRoleWidth == nil {
			b524EnergyRoleUnknown = append(b524EnergyRoleUnknown, "width")
		}
		if !b524EnergyRoleReflect.DeepEqual(b524EnergyRoleEntry.B524EnergyRoleWidthBytes, b524EnergyRoleWidth) || !b524EnergyRoleReflect.DeepEqual(b524EnergyRoleEntry.B524EnergyRoleCodec, b524EnergyRoleCodec) || !b524EnergyRoleReflect.DeepEqual(b524EnergyRoleEntry.B524EnergyRoleUnknown, b524EnergyRoleUnknown) {
			b524EnergyRoleT.Fatalf("target %d: codec/width/unknown changed", b524EnergyRoleIndex)
		}
		if b524EnergyRoleEntry.B524EnergyRoleGrade != "Hypothesis" || b524EnergyRoleEntry.B524EnergyRoleMetadataBoundaryGrade != "Proven" || b524EnergyRoleEntry.B524EnergyRolePhysicalQualification != "Unknown" || b524EnergyRoleEntry.B524EnergyRoleMetadataBoundary != "Metadata lookup cannot establish applied state or installed capability." {
			b524EnergyRoleT.Fatal("grade boundary changed")
		}
		if b524EnergyRoleEntry.B524EnergyRoleFieldOffset != nil || b524EnergyRoleEntry.B524EnergyRoleByteOrder != nil || b524EnergyRoleEntry.B524EnergyRoleNativeUnit != nil || b524EnergyRoleEntry.B524EnergyRoleControlEnumeration != nil || b524EnergyRoleEntry.B524EnergyRoleResponseForm != nil {
			b524EnergyRoleT.Fatal("unprovided metadata invented")
		}
		if b524EnergyRoleEntry.B524EnergyRoleRawPreservation != "Preserve native raw values and partial last-known-good fields independently." || b524EnergyRoleEntry.B524EnergyRoleNativeContract != "https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/protocols/vaillant/ebus-vaillant-b524.md" {
			b524EnergyRoleT.Fatal("native contract/raw boundary changed")
		}
	}
}

func TestAb95ExactKeyAxes(b524EnergyRoleT *b524EnergyRoleTesting.T) {
	b524EnergyRoleChanges := map[string]func(*B524EnergyRoleKey){
		"ab95_family":  func(b524EnergyRoleKey *B524EnergyRoleKey) { b524EnergyRoleKey.B524EnergyRoleFamily = "other" },
		"ab95_variant": func(b524EnergyRoleKey *B524EnergyRoleKey) { b524EnergyRoleKey.B524EnergyRoleVariant = "other" },
		"ab95_primary": func(b524EnergyRoleKey *B524EnergyRoleKey) { b524EnergyRoleKey.B524EnergyRolePrimary = 0 },
		"ab95_secondary": func(b524EnergyRoleKey *B524EnergyRoleKey) {
			b524EnergyRoleKey.B524EnergyRoleSecondary = b524EnergyRolePointer(uint8(0))
		},
		"ab95_secondary_absent": func(b524EnergyRoleKey *B524EnergyRoleKey) { b524EnergyRoleKey.B524EnergyRoleSecondary = nil },
		"ab95_operation": func(b524EnergyRoleKey *B524EnergyRoleKey) {
			b524EnergyRoleKey.B524EnergyRoleOperation = b524EnergyRolePointer(uint8(6))
		},
		"ab95_operation_zero": func(b524EnergyRoleKey *B524EnergyRoleKey) {
			b524EnergyRoleKey.B524EnergyRoleOperation = b524EnergyRolePointer(uint8(0))
		},
		"ab95_operation_absent": func(b524EnergyRoleKey *B524EnergyRoleKey) { b524EnergyRoleKey.B524EnergyRoleOperation = nil },
		"ab95_direction":        func(b524EnergyRoleKey *B524EnergyRoleKey) { b524EnergyRoleKey.B524EnergyRoleDirection = "response" },
		"ab95_field":            func(b524EnergyRoleKey *B524EnergyRoleKey) { b524EnergyRoleKey.B524EnergyRoleField = "other" },
		"ab95_group": func(b524EnergyRoleKey *B524EnergyRoleKey) {
			b524EnergyRoleKey.B524EnergyRoleGroup = b524EnergyRolePointer(*b524EnergyRoleKey.B524EnergyRoleGroup ^ 1)
		},
		"ab95_group_absent": func(b524EnergyRoleKey *B524EnergyRoleKey) { b524EnergyRoleKey.B524EnergyRoleGroup = nil },
		"ab95_requested_instance": func(b524EnergyRoleKey *B524EnergyRoleKey) {
			b524EnergyRoleKey.B524EnergyRoleRequestedInstance = b524EnergyRolePointer(uint8(0))
		},
		"ab95_requested_instance_absent": func(b524EnergyRoleKey *B524EnergyRoleKey) { b524EnergyRoleKey.B524EnergyRoleRequestedInstance = nil },
		"ab95_register": func(b524EnergyRoleKey *B524EnergyRoleKey) {
			b524EnergyRoleKey.B524EnergyRoleRegister = b524EnergyRolePointer(*b524EnergyRoleKey.B524EnergyRoleRegister + 1)
		},
		"ab95_register_zero": func(b524EnergyRoleKey *B524EnergyRoleKey) {
			b524EnergyRoleKey.B524EnergyRoleRegister = b524EnergyRolePointer(uint16(0))
		},
		"ab95_register_absent": func(b524EnergyRoleKey *B524EnergyRoleKey) { b524EnergyRoleKey.B524EnergyRoleRegister = nil },
		"ab95_command_role":    func(b524EnergyRoleKey *B524EnergyRoleKey) { b524EnergyRoleKey.B524EnergyCommandRole = "other" },
	}
	for b524EnergyRoleName, b524EnergyRoleChange := range b524EnergyRoleChanges {
		b524EnergyRoleT.Run(b524EnergyRoleName, func(b524EnergyRoleT *b524EnergyRoleTesting.T) {
			for _, b524EnergyRoleKey := range b524EnergyRoleExpectedKeys() {
				if _, b524EnergyRoleOK, b524EnergyRoleError := B524EnergyRoleLookup(b524EnergyRoleKey); !b524EnergyRoleOK || b524EnergyRoleError != nil {
					b524EnergyRoleT.Fatal("invalid positive control")
				}
				b524EnergyRoleChange(&b524EnergyRoleKey)
				if b524EnergyRoleEntry, b524EnergyRoleOK, b524EnergyRoleError := B524EnergyRoleLookup(b524EnergyRoleKey); b524EnergyRoleOK || b524EnergyRoleError != nil || !b524EnergyRoleReflect.DeepEqual(b524EnergyRoleEntry, B524EnergyRoleEntry{}) {
					b524EnergyRoleT.Fatalf("axis inherited a target: %v %v", b524EnergyRoleOK, b524EnergyRoleError)
				}
			}
		})
	}
}

func b524EnergyRoleInstall(b524EnergyRoleT *b524EnergyRoleTesting.T, b524EnergyRoleEntries []B524EnergyRoleEntry) {
	b524EnergyRoleT.Helper()
	b524EnergyRoleEncoded, b524EnergyRoleError := b524EnergyRoleJSON.Marshal(b524EnergyRoleEntries)
	if b524EnergyRoleError != nil {
		b524EnergyRoleT.Fatal(b524EnergyRoleError)
	}
	b524EnergyRoleData = b524EnergyRoleEncoded
}

func TestAb95AbsentVersusZero(b524EnergyRoleT *b524EnergyRoleTesting.T) {
	b524EnergyRoleOriginal := b524EnergyRoleData
	defer func() { b524EnergyRoleData = b524EnergyRoleOriginal }()
	b524EnergyRoleChanges := map[string]func(*B524EnergyRoleKey, bool){
		"ab95_secondary": func(b524EnergyRoleKey *B524EnergyRoleKey, b524EnergyRolePresent bool) {
			b524EnergyRoleKey.B524EnergyRoleSecondary = nil
			if b524EnergyRolePresent {
				b524EnergyRoleKey.B524EnergyRoleSecondary = b524EnergyRolePointer(uint8(0))
			}
		},
		"ab95_operation": func(b524EnergyRoleKey *B524EnergyRoleKey, b524EnergyRolePresent bool) {
			b524EnergyRoleKey.B524EnergyRoleOperation = nil
			if b524EnergyRolePresent {
				b524EnergyRoleKey.B524EnergyRoleOperation = b524EnergyRolePointer(uint8(0))
			}
		},
		"ab95_group": func(b524EnergyRoleKey *B524EnergyRoleKey, b524EnergyRolePresent bool) {
			b524EnergyRoleKey.B524EnergyRoleGroup = nil
			if b524EnergyRolePresent {
				b524EnergyRoleKey.B524EnergyRoleGroup = b524EnergyRolePointer(uint8(0))
			}
		},
		"ab95_requested_instance": func(b524EnergyRoleKey *B524EnergyRoleKey, b524EnergyRolePresent bool) {
			b524EnergyRoleKey.B524EnergyRoleRequestedInstance = nil
			if b524EnergyRolePresent {
				b524EnergyRoleKey.B524EnergyRoleRequestedInstance = b524EnergyRolePointer(uint8(0))
			}
		},
		"ab95_register": func(b524EnergyRoleKey *B524EnergyRoleKey, b524EnergyRolePresent bool) {
			b524EnergyRoleKey.B524EnergyRoleRegister = nil
			if b524EnergyRolePresent {
				b524EnergyRoleKey.B524EnergyRoleRegister = b524EnergyRolePointer(uint16(0))
			}
		},
	}
	for b524EnergyRoleName, b524EnergyRoleChange := range b524EnergyRoleChanges {
		b524EnergyRoleT.Run(b524EnergyRoleName, func(b524EnergyRoleT *b524EnergyRoleTesting.T) {
			for _, b524EnergyRolePresent := range []bool{false, true} {
				b524EnergyRoleData = b524EnergyRoleOriginal
				b524EnergyRoleEntries, b524EnergyRoleError := B524EnergyRoleEntries()
				if b524EnergyRoleError != nil {
					b524EnergyRoleT.Fatal(b524EnergyRoleError)
				}
				b524EnergyRoleChange(&b524EnergyRoleEntries[0].B524EnergyRoleKey, b524EnergyRolePresent)
				b524EnergyRoleInstall(b524EnergyRoleT, b524EnergyRoleEntries)
				b524EnergyRoleKey := b524EnergyRoleEntries[0].B524EnergyRoleKey
				if _, b524EnergyRoleOK, b524EnergyRoleError := B524EnergyRoleLookup(b524EnergyRoleKey); !b524EnergyRoleOK || b524EnergyRoleError != nil {
					b524EnergyRoleT.Fatal("synthetic positive control")
				}
				b524EnergyRoleChange(&b524EnergyRoleKey, !b524EnergyRolePresent)
				if _, b524EnergyRoleOK, b524EnergyRoleError := B524EnergyRoleLookup(b524EnergyRoleKey); b524EnergyRoleOK || b524EnergyRoleError != nil {
					b524EnergyRoleT.Fatal("absent/zero collapsed")
				}
			}
		})
	}
}

func b524EnergyRoleMutate(b524EnergyRoleValue b524EnergyRoleReflect.Value) {
	switch b524EnergyRoleValue.Kind() {
	case b524EnergyRoleReflect.Pointer:
		if !b524EnergyRoleValue.IsNil() {
			b524EnergyRoleMutate(b524EnergyRoleValue.Elem())
		}
	case b524EnergyRoleReflect.Struct:
		for b524EnergyRoleIndex := 0; b524EnergyRoleIndex < b524EnergyRoleValue.NumField(); b524EnergyRoleIndex++ {
			b524EnergyRoleMutate(b524EnergyRoleValue.Field(b524EnergyRoleIndex))
		}
	case b524EnergyRoleReflect.Slice:
		for b524EnergyRoleIndex := 0; b524EnergyRoleIndex < b524EnergyRoleValue.Len(); b524EnergyRoleIndex++ {
			b524EnergyRoleMutate(b524EnergyRoleValue.Index(b524EnergyRoleIndex))
		}
	case b524EnergyRoleReflect.String:
		b524EnergyRoleValue.SetString("ab95_mutated")
	case b524EnergyRoleReflect.Uint8, b524EnergyRoleReflect.Uint16:
		b524EnergyRoleValue.SetUint(99)
	}
}

func TestAb95IndependentStorage(b524EnergyRoleT *b524EnergyRoleTesting.T) {
	b524EnergyRoleBefore, b524EnergyRoleError := B524EnergyRoleEntries()
	if b524EnergyRoleError != nil {
		b524EnergyRoleT.Fatal(b524EnergyRoleError)
	}
	b524EnergyRoleReturned, b524EnergyRoleError := B524EnergyRoleEntries()
	if b524EnergyRoleError != nil {
		b524EnergyRoleT.Fatal(b524EnergyRoleError)
	}
	b524EnergyRoleMutate(b524EnergyRoleReflect.ValueOf(&b524EnergyRoleReturned).Elem())
	b524EnergyRoleFresh, b524EnergyRoleError := B524EnergyRoleEntries()
	if b524EnergyRoleError != nil || !b524EnergyRoleReflect.DeepEqual(b524EnergyRoleFresh, b524EnergyRoleBefore) {
		b524EnergyRoleT.Fatal("shared accessor storage")
	}
	for b524EnergyRoleIndex, b524EnergyRoleKey := range b524EnergyRoleExpectedKeys() {
		b524EnergyRoleEntry, b524EnergyRoleOK, b524EnergyRoleError := B524EnergyRoleLookup(b524EnergyRoleKey)
		if !b524EnergyRoleOK || b524EnergyRoleError != nil {
			b524EnergyRoleT.Fatal("lookup positive control")
		}
		b524EnergyRoleMutate(b524EnergyRoleReflect.ValueOf(&b524EnergyRoleEntry).Elem())
		b524EnergyRoleFreshEntry, b524EnergyRoleOK, b524EnergyRoleError := B524EnergyRoleLookup(b524EnergyRoleKey)
		if !b524EnergyRoleOK || b524EnergyRoleError != nil || !b524EnergyRoleReflect.DeepEqual(b524EnergyRoleFreshEntry, b524EnergyRoleBefore[b524EnergyRoleIndex]) {
			b524EnergyRoleT.Fatal("shared lookup storage")
		}
	}
}

func TestAb95MalformedData(b524EnergyRoleT *b524EnergyRoleTesting.T) {
	b524EnergyRoleOriginal := b524EnergyRoleData
	defer func() { b524EnergyRoleData = b524EnergyRoleOriginal }()
	b524EnergyRoleCases := map[string][]byte{"ab95_syntax": []byte("["), "ab95_empty": []byte("[]"), "ab95_null": []byte("null"), "ab95_trailing": append(append([]byte(nil), b524EnergyRoleOriginal...), []byte(" null")...)}
	for b524EnergyRoleName, b524EnergyRoleChange := range map[string]func([]B524EnergyRoleEntry){
		"ab95_duplicate": func(b524EnergyRoleEntries []B524EnergyRoleEntry) {
			b524EnergyRoleEntries[1].B524EnergyRoleKey = b524EnergyRoleEntries[0].B524EnergyRoleKey
		},
		"ab95_missing_family": func(b524EnergyRoleEntries []B524EnergyRoleEntry) {
			b524EnergyRoleEntries[4].B524EnergyRoleKey.B524EnergyRoleFamily = ""
		},
		"ab95_grade": func(b524EnergyRoleEntries []B524EnergyRoleEntry) {
			b524EnergyRoleEntries[4].B524EnergyRoleGrade = "Proven"
		},
	} {
		b524EnergyRoleData = b524EnergyRoleOriginal
		b524EnergyRoleEntries, b524EnergyRoleError := B524EnergyRoleEntries()
		if b524EnergyRoleError != nil {
			b524EnergyRoleT.Fatal(b524EnergyRoleError)
		}
		b524EnergyRoleChange(b524EnergyRoleEntries)
		b524EnergyRoleInstall(b524EnergyRoleT, b524EnergyRoleEntries)
		b524EnergyRoleCases[b524EnergyRoleName] = b524EnergyRoleData
	}
	for b524EnergyRoleName, b524EnergyRoleBytes := range b524EnergyRoleCases {
		b524EnergyRoleT.Run(b524EnergyRoleName, func(b524EnergyRoleT *b524EnergyRoleTesting.T) {
			b524EnergyRoleData = b524EnergyRoleBytes
			if b524EnergyRoleEntries, b524EnergyRoleError := B524EnergyRoleEntries(); b524EnergyRoleError == nil || b524EnergyRoleEntries != nil {
				b524EnergyRoleT.Fatal("invalid data returned entries")
			}
			if b524EnergyRoleEntry, b524EnergyRoleOK, b524EnergyRoleError := B524EnergyRoleLookup(b524EnergyRoleExpectedKeys()[0]); b524EnergyRoleError == nil || b524EnergyRoleOK || !b524EnergyRoleReflect.DeepEqual(b524EnergyRoleEntry, B524EnergyRoleEntry{}) {
				b524EnergyRoleT.Fatal("invalid data returned successful lookup")
			}
		})
	}
}
