package vaillant

import (
	b516DirectionalTypeBytes "bytes"
	b516DirectionalTypeJSON "encoding/json"
	b516DirectionalTypeReflect "reflect"
	b516DirectionalTypeTesting "testing"
)

func b516DirectionalTypeKeyFor(b516DirectionalTypeVariant, b516DirectionalTypeDirection, b516DirectionalTypeField string) B516DirectionalTypeKey {
	b516DirectionalTypeSecondary := uint8(22)
	return B516DirectionalTypeKey{B516DirectionalTypeFamily: "vaillant", B516DirectionalTypeVariant: b516DirectionalTypeVariant, B516DirectionalTypePrimary: 181, B516DirectionalTypeSecondary: &b516DirectionalTypeSecondary, B516DirectionalTypeDirection: b516DirectionalTypeDirection, B516DirectionalTypeField: b516DirectionalTypeField, B516DirectionalTypeAlternative: "none", B516DirectionalTypeInstanceScope: "none"}
}

func b516DirectionalTypeRowFor(b516DirectionalTypeT *b516DirectionalTypeTesting.T, b516DirectionalTypeKey B516DirectionalTypeKey) map[string]b516DirectionalTypeJSON.RawMessage {
	b516DirectionalTypeT.Helper()
	b516DirectionalTypeRaw, b516DirectionalTypeFound, b516DirectionalTypeErr := B516DirectionalTypeLookup(b516DirectionalTypeKey)
	if b516DirectionalTypeErr != nil || !b516DirectionalTypeFound {
		b516DirectionalTypeT.Fatalf("ab20 lookup: %v %v", b516DirectionalTypeFound, b516DirectionalTypeErr)
	}
	var b516DirectionalTypeRow map[string]b516DirectionalTypeJSON.RawMessage
	if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeRaw, &b516DirectionalTypeRow); b516DirectionalTypeErr != nil {
		b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
	}
	return b516DirectionalTypeRow
}

func b516DirectionalTypeWant(b516DirectionalTypeT *b516DirectionalTypeTesting.T, b516DirectionalTypeRow map[string]b516DirectionalTypeJSON.RawMessage, b516DirectionalTypeField, b516DirectionalTypeExpected string) {
	b516DirectionalTypeT.Helper()
	var b516DirectionalTypeGot, b516DirectionalTypeWant any
	if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeRow[b516DirectionalTypeField], &b516DirectionalTypeGot); b516DirectionalTypeErr != nil {
		b516DirectionalTypeT.Fatalf("ab20 %s: %v", b516DirectionalTypeField, b516DirectionalTypeErr)
	}
	if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal([]byte(b516DirectionalTypeExpected), &b516DirectionalTypeWant); b516DirectionalTypeErr != nil {
		b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
	}
	if !b516DirectionalTypeReflect.DeepEqual(b516DirectionalTypeGot, b516DirectionalTypeWant) {
		b516DirectionalTypeT.Fatalf("ab20 %s = %v, want %v", b516DirectionalTypeField, b516DirectionalTypeGot, b516DirectionalTypeWant)
	}
}

func TestAb20DirectionalRepresentations(b516DirectionalTypeT *b516DirectionalTypeTesting.T) {
	b516DirectionalTypeRequest := b516DirectionalTypeRowFor(b516DirectionalTypeT, b516DirectionalTypeKeyFor("energy-type", "request", "energy-type"))
	b516DirectionalTypeRead := b516DirectionalTypeRowFor(b516DirectionalTypeT, b516DirectionalTypeKeyFor("energy-type", "read-positive", "energy-type"))
	b516DirectionalTypeWrite := b516DirectionalTypeRowFor(b516DirectionalTypeT, b516DirectionalTypeKeyFor("energy-type", "write-positive", "energy-type"))
	b516DirectionalTypeWant(b516DirectionalTypeT, b516DirectionalTypeRequest, "representation", `"raw-u8-enum"`)
	b516DirectionalTypeWant(b516DirectionalTypeT, b516DirectionalTypeRequest, "enum_labels", `"Unknown"`)
	b516DirectionalTypeWant(b516DirectionalTypeT, b516DirectionalTypeRead, "representation", `"low-six-and-independent-flags"`)
	b516DirectionalTypeWant(b516DirectionalTypeT, b516DirectionalTypeRead, "flags", `[{"bit":6,"mask":64},{"bit":7,"mask":128}]`)
	var b516DirectionalTypeFlags []struct {
		B516DirectionalTypeBit  uint8 `json:"bit"`
		B516DirectionalTypeMask uint8 `json:"mask"`
	}
	if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeRead["flags"], &b516DirectionalTypeFlags); b516DirectionalTypeErr != nil {
		b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
	}
	for _, b516DirectionalTypeFlag := range b516DirectionalTypeFlags {
		if uint8(0xC1)&b516DirectionalTypeFlag.B516DirectionalTypeMask != b516DirectionalTypeFlag.B516DirectionalTypeMask {
			b516DirectionalTypeT.Fatalf("ab20 missing high flag %d", b516DirectionalTypeFlag.B516DirectionalTypeBit)
		}
	}
	b516DirectionalTypeWant(b516DirectionalTypeT, b516DirectionalTypeWrite, "representation", `"scalar-six-bit"`)
	for _, b516DirectionalTypeRow := range []map[string]b516DirectionalTypeJSON.RawMessage{b516DirectionalTypeRequest, b516DirectionalTypeWrite} {
		if _, b516DirectionalTypeExists := b516DirectionalTypeRow["flags"]; b516DirectionalTypeExists {
			b516DirectionalTypeT.Fatal("ab20 flags assigned outside read-positive")
		}
	}
	if _, b516DirectionalTypeExists := b516DirectionalTypeRequest["value_mask"]; b516DirectionalTypeExists {
		b516DirectionalTypeT.Fatal("ab20 request flattened into scalar")
	}
	for _, b516DirectionalTypeCase := range []struct {
		b516DirectionalTypeName     string
		b516DirectionalTypeRow      map[string]b516DirectionalTypeJSON.RawMessage
		b516DirectionalTypeMasks    map[string]uint8
		b516DirectionalTypeExpected map[string]uint8
	}{
		{"ab20-request", b516DirectionalTypeRequest, map[string]uint8{"raw_mask": 255}, map[string]uint8{"raw_mask": 193}},
		{"ab20-read", b516DirectionalTypeRead, map[string]uint8{"raw_mask": 255, "value_mask": 63}, map[string]uint8{"raw_mask": 193, "value_mask": 1}},
		{"ab20-write", b516DirectionalTypeWrite, map[string]uint8{"raw_mask": 255, "value_mask": 63, "unmodeled_mask": 192}, map[string]uint8{"raw_mask": 193, "value_mask": 1, "unmodeled_mask": 192}},
	} {
		b516DirectionalTypeT.Run(b516DirectionalTypeCase.b516DirectionalTypeName, func(b516DirectionalTypeT *b516DirectionalTypeTesting.T) {
			for b516DirectionalTypeField, b516DirectionalTypeExpectedMask := range b516DirectionalTypeCase.b516DirectionalTypeMasks {
				var b516DirectionalTypeMask uint8
				if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeCase.b516DirectionalTypeRow[b516DirectionalTypeField], &b516DirectionalTypeMask); b516DirectionalTypeErr != nil {
					b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
				}
				if b516DirectionalTypeMask != b516DirectionalTypeExpectedMask || uint8(0xC1)&b516DirectionalTypeMask != b516DirectionalTypeCase.b516DirectionalTypeExpected[b516DirectionalTypeField] {
					b516DirectionalTypeT.Fatalf("ab20 %s: mask %d", b516DirectionalTypeField, b516DirectionalTypeMask)
				}
			}
		})
	}
}

func TestAb20TimeAndEvidenceBoundaries(b516DirectionalTypeT *b516DirectionalTypeTesting.T) {
	for _, b516DirectionalTypeCase := range []struct{ b516DirectionalTypeVariant, b516DirectionalTypeState string }{{"generic", "header-only"}, {"other", "absent"}} {
		b516DirectionalTypeRow := b516DirectionalTypeRowFor(b516DirectionalTypeT, b516DirectionalTypeKeyFor(b516DirectionalTypeCase.b516DirectionalTypeVariant, "positive", "system-time"))
		b516DirectionalTypeWant(b516DirectionalTypeT, b516DirectionalTypeRow, "positive_state", `"`+b516DirectionalTypeCase.b516DirectionalTypeState+`"`)
	}
	b516DirectionalTypeData, b516DirectionalTypeErr := B516DirectionalTypeEntries()
	if b516DirectionalTypeErr != nil {
		b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
	}
	var b516DirectionalTypeRows []map[string]b516DirectionalTypeJSON.RawMessage
	if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeData, &b516DirectionalTypeRows); b516DirectionalTypeErr != nil {
		b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
	}
	if len(b516DirectionalTypeRows) != 5 {
		b516DirectionalTypeT.Fatalf("ab20 finite area: %d", len(b516DirectionalTypeRows))
	}
	for _, b516DirectionalTypeRow := range b516DirectionalTypeRows {
		var b516DirectionalTypeIdentity B516DirectionalTypeKey
		if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeRow["key"], &b516DirectionalTypeIdentity); b516DirectionalTypeErr != nil {
			b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
		}
		b516DirectionalTypeLookedUp := b516DirectionalTypeRowFor(b516DirectionalTypeT, b516DirectionalTypeIdentity)
		if !b516DirectionalTypeReflect.DeepEqual(b516DirectionalTypeLookedUp, b516DirectionalTypeRow) {
			b516DirectionalTypeT.Fatal("ab20 lookup did not return whole entry")
		}
		b516DirectionalTypeWant(b516DirectionalTypeT, b516DirectionalTypeRow, "grade", `"Hypothesis"`)
		b516DirectionalTypeWant(b516DirectionalTypeT, b516DirectionalTypeRow, "metadata_presence", `"Proven"`)
		b516DirectionalTypeWant(b516DirectionalTypeT, b516DirectionalTypeRow, "raw_form", `"retained without normalization"`)
		for _, b516DirectionalTypeField := range []string{"physical_length", "native_support", "physical_effects", "persistence", "delivery", "timeout", "broadcast"} {
			b516DirectionalTypeWant(b516DirectionalTypeT, b516DirectionalTypeRow, b516DirectionalTypeField, `"Unknown"`)
		}
		b516DirectionalTypeWant(b516DirectionalTypeT, b516DirectionalTypeRow, "unknown", `["encoded operation","exact enum tables","missing positions","native qualification"]`)
		var b516DirectionalTypeKey map[string]b516DirectionalTypeJSON.RawMessage
		if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeRow["key"], &b516DirectionalTypeKey); b516DirectionalTypeErr != nil {
			b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
		}
		for _, b516DirectionalTypeField := range []string{"operation", "position", "enum_table"} {
			if _, b516DirectionalTypeExists := b516DirectionalTypeKey[b516DirectionalTypeField]; b516DirectionalTypeExists {
				b516DirectionalTypeT.Fatalf("ab20 invented key %s", b516DirectionalTypeField)
			}
			if _, b516DirectionalTypeExists := b516DirectionalTypeRow[b516DirectionalTypeField]; b516DirectionalTypeExists {
				b516DirectionalTypeT.Fatalf("ab20 invented metadata %s", b516DirectionalTypeField)
			}
		}
	}
}

func b516DirectionalTypeAxisChanges() map[string]func(*B516DirectionalTypeKey) {
	return map[string]func(*B516DirectionalTypeKey){
		"ab20-family": func(b516DirectionalTypeKey *B516DirectionalTypeKey) {
			b516DirectionalTypeKey.B516DirectionalTypeFamily = "synthetic-family"
		},
		"ab20-variant": func(b516DirectionalTypeKey *B516DirectionalTypeKey) {
			b516DirectionalTypeKey.B516DirectionalTypeVariant = "synthetic-variant"
		},
		"ab20-primary": func(b516DirectionalTypeKey *B516DirectionalTypeKey) {
			b516DirectionalTypeKey.B516DirectionalTypePrimary = 0
		},
		"ab20-secondary": func(b516DirectionalTypeKey *B516DirectionalTypeKey) {
			b516DirectionalTypeValue := uint8(0)
			b516DirectionalTypeKey.B516DirectionalTypeSecondary = &b516DirectionalTypeValue
		},
		"ab20-operation": func(b516DirectionalTypeKey *B516DirectionalTypeKey) {
			b516DirectionalTypeValue := uint8(0)
			b516DirectionalTypeKey.B516DirectionalTypeOperation = &b516DirectionalTypeValue
		},
		"ab20-direction": func(b516DirectionalTypeKey *B516DirectionalTypeKey) {
			b516DirectionalTypeKey.B516DirectionalTypeDirection = "synthetic-direction"
		},
		"ab20-field": func(b516DirectionalTypeKey *B516DirectionalTypeKey) {
			b516DirectionalTypeKey.B516DirectionalTypeField = "synthetic-field"
		},
		"ab20-alternative": func(b516DirectionalTypeKey *B516DirectionalTypeKey) {
			b516DirectionalTypeKey.B516DirectionalTypeAlternative = "synthetic-alternative"
		},
		"ab20-group": func(b516DirectionalTypeKey *B516DirectionalTypeKey) {
			b516DirectionalTypeValue := uint8(0)
			b516DirectionalTypeKey.B516DirectionalTypeGroup = &b516DirectionalTypeValue
		},
		"ab20-instance-scope": func(b516DirectionalTypeKey *B516DirectionalTypeKey) {
			b516DirectionalTypeKey.B516DirectionalTypeInstanceScope = "dynamic"
		},
		"ab20-register": func(b516DirectionalTypeKey *B516DirectionalTypeKey) {
			b516DirectionalTypeValue := uint16(0)
			b516DirectionalTypeKey.B516DirectionalTypeRegister = &b516DirectionalTypeValue
		},
		"ab20-secondary-absent": func(b516DirectionalTypeKey *B516DirectionalTypeKey) {
			b516DirectionalTypeKey.B516DirectionalTypeSecondary = nil
		},
	}
}

func TestAb20ExactAxes(b516DirectionalTypeT *b516DirectionalTypeTesting.T) {
	b516DirectionalTypeBase := b516DirectionalTypeKeyFor("energy-type", "read-positive", "energy-type")
	b516DirectionalTypeOriginal, b516DirectionalTypeFound, b516DirectionalTypeErr := B516DirectionalTypeLookup(b516DirectionalTypeBase)
	if b516DirectionalTypeErr != nil || !b516DirectionalTypeFound {
		b516DirectionalTypeT.Fatalf("ab20 base: %v %v", b516DirectionalTypeFound, b516DirectionalTypeErr)
	}
	for b516DirectionalTypeName, b516DirectionalTypeChange := range b516DirectionalTypeAxisChanges() {
		b516DirectionalTypeT.Run(b516DirectionalTypeName, func(b516DirectionalTypeT *b516DirectionalTypeTesting.T) {
			b516DirectionalTypeChanged := b516DirectionalTypeBase
			b516DirectionalTypeChange(&b516DirectionalTypeChanged)
			b516DirectionalTypeRaw, b516DirectionalTypeFound, b516DirectionalTypeErr := B516DirectionalTypeLookup(b516DirectionalTypeChanged)
			if b516DirectionalTypeErr != nil || b516DirectionalTypeFound || b516DirectionalTypeRaw != nil {
				b516DirectionalTypeT.Fatalf("ab20 fallback: %v %v", b516DirectionalTypeFound, b516DirectionalTypeErr)
			}
			b516DirectionalTypeFixture, b516DirectionalTypeErr := b516DirectionalTypeJSON.Marshal([]any{map[string]any{"id": "ab20-control-original", "key": b516DirectionalTypeBase}, map[string]any{"id": "ab20-control-changed", "key": b516DirectionalTypeChanged}})
			if b516DirectionalTypeErr != nil {
				b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
			}
			for _, b516DirectionalTypeKey := range []B516DirectionalTypeKey{b516DirectionalTypeBase, b516DirectionalTypeChanged} {
				b516DirectionalTypeRaw, b516DirectionalTypeFound, b516DirectionalTypeErr := b516DirectionalTypeLookup(b516DirectionalTypeFixture, b516DirectionalTypeKey)
				if b516DirectionalTypeErr != nil || !b516DirectionalTypeFound {
					b516DirectionalTypeT.Fatalf("ab20 independent control: %v %v", b516DirectionalTypeFound, b516DirectionalTypeErr)
				}
				var b516DirectionalTypeEntry struct {
					B516DirectionalTypeKey B516DirectionalTypeKey `json:"key"`
				}
				if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeRaw, &b516DirectionalTypeEntry); b516DirectionalTypeErr != nil {
					b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
				}
				if !b516DirectionalTypeReflect.DeepEqual(b516DirectionalTypeEntry.B516DirectionalTypeKey, b516DirectionalTypeKey) {
					b516DirectionalTypeT.Fatal("ab20 wrong whole row")
				}
			}
		})
	}
	b516DirectionalTypeAgain, _, _ := B516DirectionalTypeLookup(b516DirectionalTypeBase)
	if !b516DirectionalTypeBytes.Equal(b516DirectionalTypeOriginal, b516DirectionalTypeAgain) {
		b516DirectionalTypeT.Fatal("ab20 key mutation changed catalog")
	}
}

func TestAb20ValidationAndCopies(b516DirectionalTypeT *b516DirectionalTypeTesting.T) {
	b516DirectionalTypeKey := b516DirectionalTypeKeyFor("generic", "positive", "system-time")
	for _, b516DirectionalTypeInvalid := range []string{`{`, `null`, `[]`, `{}`, `[null]`, `[{}]`, `[] []`, `[{"id":"ab20-bad","key":{"primary":256}}]`} {
		if _, b516DirectionalTypeErr := b516DirectionalTypeValidate([]byte(b516DirectionalTypeInvalid)); b516DirectionalTypeErr == nil {
			b516DirectionalTypeT.Fatalf("ab20 accepted %s", b516DirectionalTypeInvalid)
		}
		if _, b516DirectionalTypeFound, b516DirectionalTypeErr := b516DirectionalTypeLookup([]byte(b516DirectionalTypeInvalid), b516DirectionalTypeKey); b516DirectionalTypeErr == nil || b516DirectionalTypeFound {
			b516DirectionalTypeT.Fatal("ab20 lookup did not fail closed")
		}
	}
	b516DirectionalTypeDuplicate, b516DirectionalTypeErr := b516DirectionalTypeJSON.Marshal([]any{map[string]any{"id": "ab20-first", "key": b516DirectionalTypeKey}, map[string]any{"id": "ab20-second", "key": b516DirectionalTypeKey}})
	if b516DirectionalTypeErr != nil {
		b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
	}
	if _, b516DirectionalTypeErr := b516DirectionalTypeValidate(b516DirectionalTypeDuplicate); b516DirectionalTypeErr == nil {
		b516DirectionalTypeT.Fatal("ab20 duplicate key admitted")
	}
	if _, b516DirectionalTypeFound, b516DirectionalTypeErr := b516DirectionalTypeLookup(b516DirectionalTypeDuplicate, b516DirectionalTypeKey); b516DirectionalTypeErr == nil || b516DirectionalTypeFound {
		b516DirectionalTypeT.Fatal("ab20 duplicate lookup admitted")
	}
	b516DirectionalTypeEntries, b516DirectionalTypeErr := B516DirectionalTypeEntries()
	if b516DirectionalTypeErr != nil {
		b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
	}
	b516DirectionalTypeRow, b516DirectionalTypeFound, b516DirectionalTypeErr := B516DirectionalTypeLookup(b516DirectionalTypeKey)
	if b516DirectionalTypeErr != nil || !b516DirectionalTypeFound {
		b516DirectionalTypeT.Fatalf("ab20 lookup: %v %v", b516DirectionalTypeFound, b516DirectionalTypeErr)
	}
	b516DirectionalTypeExpectedEntries := append([]byte(nil), b516DirectionalTypeEntries...)
	b516DirectionalTypeExpectedRow := append([]byte(nil), b516DirectionalTypeRow...)
	b516DirectionalTypeSecondEntries, _ := B516DirectionalTypeEntries()
	b516DirectionalTypeSecondRow, _, _ := B516DirectionalTypeLookup(b516DirectionalTypeKey)
	for b516DirectionalTypeIndex := range b516DirectionalTypeEntries {
		b516DirectionalTypeEntries[b516DirectionalTypeIndex] = 'x'
	}
	for b516DirectionalTypeIndex := range b516DirectionalTypeRow {
		b516DirectionalTypeRow[b516DirectionalTypeIndex] = 'x'
	}
	b516DirectionalTypeLaterEntries, b516DirectionalTypeErr := B516DirectionalTypeEntries()
	if b516DirectionalTypeErr != nil {
		b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
	}
	b516DirectionalTypeLaterRow, b516DirectionalTypeFound, b516DirectionalTypeErr := B516DirectionalTypeLookup(b516DirectionalTypeKey)
	if b516DirectionalTypeErr != nil || !b516DirectionalTypeFound {
		b516DirectionalTypeT.Fatalf("ab20 later: %v %v", b516DirectionalTypeFound, b516DirectionalTypeErr)
	}
	if !b516DirectionalTypeBytes.Equal(b516DirectionalTypeExpectedEntries, b516DirectionalTypeSecondEntries) || !b516DirectionalTypeBytes.Equal(b516DirectionalTypeExpectedEntries, b516DirectionalTypeLaterEntries) || !b516DirectionalTypeBytes.Equal(b516DirectionalTypeExpectedRow, b516DirectionalTypeSecondRow) || !b516DirectionalTypeBytes.Equal(b516DirectionalTypeExpectedRow, b516DirectionalTypeLaterRow) {
		b516DirectionalTypeT.Fatal("ab20 shared return memory")
	}
}

func TestAb20OptionalPresence(b516DirectionalTypeT *b516DirectionalTypeTesting.T) {
	for _, b516DirectionalTypeName := range []string{"ab20-secondary", "ab20-operation", "ab20-group", "ab20-register"} {
		b516DirectionalTypeT.Run(b516DirectionalTypeName, func(b516DirectionalTypeT *b516DirectionalTypeTesting.T) {
			b516DirectionalTypeAbsent := b516DirectionalTypeKeyFor("generic", "positive", "system-time")
			b516DirectionalTypeAbsent.B516DirectionalTypeSecondary = nil
			b516DirectionalTypeZero := b516DirectionalTypeAbsent
			b516DirectionalTypeAxisChanges()[b516DirectionalTypeName](&b516DirectionalTypeZero)
			b516DirectionalTypeFixture, b516DirectionalTypeErr := b516DirectionalTypeJSON.Marshal([]any{map[string]any{"id": "ab20-absent", "key": b516DirectionalTypeAbsent}, map[string]any{"id": "ab20-zero", "key": b516DirectionalTypeZero}})
			if b516DirectionalTypeErr != nil {
				b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
			}
			for b516DirectionalTypeIndex, b516DirectionalTypeKey := range []B516DirectionalTypeKey{b516DirectionalTypeAbsent, b516DirectionalTypeZero} {
				b516DirectionalTypeRaw, b516DirectionalTypeFound, b516DirectionalTypeErr := b516DirectionalTypeLookup(b516DirectionalTypeFixture, b516DirectionalTypeKey)
				if b516DirectionalTypeErr != nil || !b516DirectionalTypeFound {
					b516DirectionalTypeT.Fatalf("ab20 optional presence: %v %v", b516DirectionalTypeFound, b516DirectionalTypeErr)
				}
				var b516DirectionalTypeEntry struct {
					B516DirectionalTypeID string `json:"id"`
				}
				if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeRaw, &b516DirectionalTypeEntry); b516DirectionalTypeErr != nil {
					b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
				}
				if b516DirectionalTypeEntry.B516DirectionalTypeID != []string{"ab20-absent", "ab20-zero"}[b516DirectionalTypeIndex] {
					b516DirectionalTypeT.Fatal("ab20 absent collapsed with zero")
				}
			}
		})
	}
}

func TestAb20IdentityValidation(b516DirectionalTypeT *b516DirectionalTypeTesting.T) {
	b516DirectionalTypeEncoded, b516DirectionalTypeErr := b516DirectionalTypeJSON.Marshal(b516DirectionalTypeKeyFor("generic", "positive", "system-time"))
	if b516DirectionalTypeErr != nil {
		b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
	}
	for _, b516DirectionalTypeField := range []string{"family", "variant", "primary", "direction", "field", "alternative", "instance_scope"} {
		b516DirectionalTypeT.Run("ab20-"+b516DirectionalTypeField, func(b516DirectionalTypeT *b516DirectionalTypeTesting.T) {
			for _, b516DirectionalTypeNull := range []bool{false, true} {
				var b516DirectionalTypeIdentity map[string]b516DirectionalTypeJSON.RawMessage
				if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeEncoded, &b516DirectionalTypeIdentity); b516DirectionalTypeErr != nil {
					b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
				}
				delete(b516DirectionalTypeIdentity, b516DirectionalTypeField)
				if b516DirectionalTypeNull {
					b516DirectionalTypeIdentity[b516DirectionalTypeField] = b516DirectionalTypeJSON.RawMessage(`null`)
				}
				b516DirectionalTypeFixture, b516DirectionalTypeErr := b516DirectionalTypeJSON.Marshal([]any{map[string]any{"id": "ab20-incomplete", "key": b516DirectionalTypeIdentity}})
				if b516DirectionalTypeErr != nil {
					b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
				}
				if _, b516DirectionalTypeErr := b516DirectionalTypeValidate(b516DirectionalTypeFixture); b516DirectionalTypeErr == nil {
					b516DirectionalTypeT.Fatal("ab20 incomplete identity admitted")
				}
			}
		})
	}
}

func TestAb20UnknownIdentityAxes(b516DirectionalTypeT *b516DirectionalTypeTesting.T) {
	b516DirectionalTypeKey := b516DirectionalTypeKeyFor("generic", "positive", "system-time")
	b516DirectionalTypeEncoded, b516DirectionalTypeErr := b516DirectionalTypeJSON.Marshal(b516DirectionalTypeKey)
	if b516DirectionalTypeErr != nil {
		b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
	}
	for _, b516DirectionalTypeField := range []string{"operaton", "unknown_axis"} {
		b516DirectionalTypeT.Run("ab20-"+b516DirectionalTypeField, func(b516DirectionalTypeT *b516DirectionalTypeTesting.T) {
			var b516DirectionalTypeIdentity map[string]b516DirectionalTypeJSON.RawMessage
			if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeEncoded, &b516DirectionalTypeIdentity); b516DirectionalTypeErr != nil {
				b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
			}
			b516DirectionalTypeIdentity[b516DirectionalTypeField] = b516DirectionalTypeJSON.RawMessage(`0`)
			b516DirectionalTypeFixture, b516DirectionalTypeErr := b516DirectionalTypeJSON.Marshal([]any{map[string]any{"id": "ab20-unknown-axis", "key": b516DirectionalTypeIdentity}})
			if b516DirectionalTypeErr != nil {
				b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
			}
			if _, b516DirectionalTypeErr := b516DirectionalTypeValidate(b516DirectionalTypeFixture); b516DirectionalTypeErr == nil {
				b516DirectionalTypeT.Fatal("ab20 unknown identity axis admitted")
			}
			if b516DirectionalTypeRaw, b516DirectionalTypeFound, b516DirectionalTypeErr := b516DirectionalTypeLookup(b516DirectionalTypeFixture, b516DirectionalTypeKey); b516DirectionalTypeErr == nil || b516DirectionalTypeFound || b516DirectionalTypeRaw != nil {
				b516DirectionalTypeT.Fatal("ab20 malformed identity lookup did not fail closed")
			}
		})
	}
	b516DirectionalTypeFixture, b516DirectionalTypeErr := b516DirectionalTypeJSON.Marshal([]any{map[string]any{"id": "ab20-variable-metadata", "key": b516DirectionalTypeKey, "ab20_extra_metadata": "retained"}})
	if b516DirectionalTypeErr != nil {
		b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
	}
	if _, b516DirectionalTypeErr := b516DirectionalTypeValidate(b516DirectionalTypeFixture); b516DirectionalTypeErr != nil {
		b516DirectionalTypeT.Fatalf("ab20 variable row metadata rejected: %v", b516DirectionalTypeErr)
	}
	b516DirectionalTypeRaw, b516DirectionalTypeFound, b516DirectionalTypeErr := b516DirectionalTypeLookup(b516DirectionalTypeFixture, b516DirectionalTypeKey)
	if b516DirectionalTypeErr != nil || !b516DirectionalTypeFound {
		b516DirectionalTypeT.Fatalf("ab20 variable metadata lookup: %v %v", b516DirectionalTypeFound, b516DirectionalTypeErr)
	}
	var b516DirectionalTypeRow map[string]b516DirectionalTypeJSON.RawMessage
	if b516DirectionalTypeErr := b516DirectionalTypeJSON.Unmarshal(b516DirectionalTypeRaw, &b516DirectionalTypeRow); b516DirectionalTypeErr != nil {
		b516DirectionalTypeT.Fatal(b516DirectionalTypeErr)
	}
	b516DirectionalTypeWant(b516DirectionalTypeT, b516DirectionalTypeRow, "ab20_extra_metadata", `"retained"`)
}
