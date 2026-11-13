package vaillant

import (
	b509Op59QuantityBinary "encoding/binary"
	b509Op59QuantityJSON "encoding/json"
	b509Op59QuantityReflect "reflect"
	b509Op59QuantityTesting "testing"
)

func b509Op59QuantityPtr[b509Op59QuantityType any](b509Op59QuantityValue b509Op59QuantityType) *b509Op59QuantityType {
	return &b509Op59QuantityValue
}

func b509Op59QuantityKey(b509Op59QuantityField string) B509Op59QuantityKey {
	return B509Op59QuantityKey{B509Op59QuantityFamily: "vaillant", B509Op59QuantityVariant: "B509", B509Op59QuantityPrimary: 181, B509Op59QuantitySecondary: b509Op59QuantityPtr(uint8(9)), B509Op59QuantityOperation: b509Op59QuantityPtr(uint8(89)), B509Op59QuantityField: b509Op59QuantityField}
}

func b509Op59QuantityGet(b509Op59QuantityT *b509Op59QuantityTesting.T, b509Op59QuantityKey B509Op59QuantityKey) B509Op59QuantityEntry {
	b509Op59QuantityT.Helper()
	b509Op59QuantityEntry, b509Op59QuantityOK, b509Op59QuantityError := B509Op59QuantityLookup(b509Op59QuantityKey)
	if b509Op59QuantityError != nil || !b509Op59QuantityOK {
		b509Op59QuantityT.Fatalf("lookup: %v %v", b509Op59QuantityOK, b509Op59QuantityError)
	}
	return b509Op59QuantityEntry
}

func TestAb14Quantities(b509Op59QuantityT *b509Op59QuantityTesting.T) {
	b509Op59QuantityEntries, b509Op59QuantityError := B509Op59QuantityEntries()
	if b509Op59QuantityError != nil || len(b509Op59QuantityEntries) != 5 {
		b509Op59QuantityT.Fatalf("entries: %v %v", b509Op59QuantityEntries, b509Op59QuantityError)
	}
	for _, b509Op59QuantityCase := range []struct {
		b509Op59QuantityField          string
		b509Op59QuantityRepresentation string
		b509Op59QuantityDivisor        float64
		b509Op59QuantityUnit           string
		b509Op59QuantityFragments      [][2]byte
		b509Op59QuantityExpected       []float64
	}{
		{"temperature", "signedLE16", 16, "Celsius", [][2]byte{{0xf0, 0xff}, {0, 0x80}}, []float64{-1, -2048}},
		{"flow", "unsignedLE16", 100, "litres per minute", [][2]byte{{0x7b, 0}, {0xff, 0xff}}, []float64{1.23, 655.35}},
		{"volume", "unsignedLE16", 10, "litres", [][2]byte{{0x7b, 0}, {0xff, 0xff}}, []float64{12.3, 6553.5}},
	} {
		b509Op59QuantityT.Run("ab14_"+b509Op59QuantityCase.b509Op59QuantityField, func(b509Op59QuantityT *b509Op59QuantityTesting.T) {
			b509Op59QuantityEntry := b509Op59QuantityGet(b509Op59QuantityT, b509Op59QuantityKey(b509Op59QuantityCase.b509Op59QuantityField))
			b509Op59QuantityWant := &B509Op59QuantityCodec{B509Op59QuantityRepresentation: b509Op59QuantityCase.b509Op59QuantityRepresentation, B509Op59QuantityWidth: b509Op59QuantityPtr(uint8(16)), B509Op59QuantityDivisor: b509Op59QuantityPtr(b509Op59QuantityCase.b509Op59QuantityDivisor), B509Op59QuantityUnit: b509Op59QuantityPtr(b509Op59QuantityCase.b509Op59QuantityUnit)}
			if !b509Op59QuantityReflect.DeepEqual(b509Op59QuantityEntry.B509Op59QuantityCodec, b509Op59QuantityWant) {
				b509Op59QuantityT.Fatalf("codec: %#v", b509Op59QuantityEntry.B509Op59QuantityCodec)
			}
			for b509Op59QuantityIndex, b509Op59QuantityFragment := range b509Op59QuantityCase.b509Op59QuantityFragments {
				b509Op59QuantityRaw := b509Op59QuantityBinary.LittleEndian.Uint16(b509Op59QuantityFragment[:])
				b509Op59QuantityValue := float64(b509Op59QuantityRaw)
				if b509Op59QuantityEntry.B509Op59QuantityCodec.B509Op59QuantityRepresentation == "signedLE16" {
					b509Op59QuantityValue = float64(int16(b509Op59QuantityRaw))
				}
				if b509Op59QuantityValue / *b509Op59QuantityEntry.B509Op59QuantityCodec.B509Op59QuantityDivisor != b509Op59QuantityCase.b509Op59QuantityExpected[b509Op59QuantityIndex] {
					b509Op59QuantityT.Fatal("independent fragment expectation failed")
				}
			}
		})
	}
	for _, b509Op59QuantityField := range []string{"minimum_volume", "maximum_volume"} {
		b509Op59QuantityEntry := b509Op59QuantityGet(b509Op59QuantityT, b509Op59QuantityKey(b509Op59QuantityField))
		if !b509Op59QuantityReflect.DeepEqual(b509Op59QuantityEntry.B509Op59QuantityCodec, &B509Op59QuantityCodec{B509Op59QuantityRepresentation: "raw"}) {
			b509Op59QuantityT.Fatal("raw limit acquired a scale, width or unit")
		}
		for _, b509Op59QuantityRaw := range []uint16{123, 456} {
			b509Op59QuantityValue := float64(b509Op59QuantityRaw)
			if b509Op59QuantityEntry.B509Op59QuantityCodec.B509Op59QuantityDivisor != nil {
				b509Op59QuantityValue /= *b509Op59QuantityEntry.B509Op59QuantityCodec.B509Op59QuantityDivisor
			}
			if b509Op59QuantityValue != float64(b509Op59QuantityRaw) {
				b509Op59QuantityT.Fatal("raw limit scaled")
			}
		}
	}
	for _, b509Op59QuantityEntry := range b509Op59QuantityEntries {
		b509Op59QuantityUnknown := []string{"absolute_field_positions", "native_support", "sentinels", "response_adaptation", "direction", "physical_qualification", "broad_catalog_completeness"}
		if b509Op59QuantityEntry.B509Op59QuantityCodec.B509Op59QuantityRepresentation == "raw" {
			b509Op59QuantityUnknown = append(b509Op59QuantityUnknown, "width", "scale", "coordinates", "enforcement", "unit")
		}
		if b509Op59QuantityEntry.B509Op59QuantityGrade != "Hypothesis" || b509Op59QuantityEntry.B509Op59QuantityMetadataPresence != "Proven" || b509Op59QuantityEntry.B509Op59QuantityPhysicalQualification != "Unknown" || !b509Op59QuantityEntry.B509Op59QuantityPreserveRaw || b509Op59QuantityEntry.B509Op59QuantityPosition != nil || b509Op59QuantityEntry.B509Op59QuantityLabel != nil || !b509Op59QuantityReflect.DeepEqual(b509Op59QuantityEntry.B509Op59QuantityUnknown, b509Op59QuantityUnknown) {
			b509Op59QuantityT.Fatal("evidence boundary changed")
		}
	}
}

func b509Op59QuantityInstall(b509Op59QuantityT *b509Op59QuantityTesting.T, b509Op59QuantityEntries []B509Op59QuantityEntry) {
	b509Op59QuantityT.Helper()
	b509Op59QuantityBytes, b509Op59QuantityError := b509Op59QuantityJSON.Marshal(b509Op59QuantityEntries)
	if b509Op59QuantityError != nil {
		b509Op59QuantityT.Fatal(b509Op59QuantityError)
	}
	b509Op59QuantityOriginal := b509Op59QuantityData
	b509Op59QuantityT.Cleanup(func() { b509Op59QuantityData = b509Op59QuantityOriginal })
	b509Op59QuantityData = string(b509Op59QuantityBytes)
}

func TestAb14ExactKeyAxes(b509Op59QuantityT *b509Op59QuantityTesting.T) {
	for _, b509Op59QuantityB524 := range []bool{false, true} {
		b509Op59QuantityT.Run(map[bool]string{false: "ab14_B509", true: "ab14_B524_synthetic"}[b509Op59QuantityB524], func(b509Op59QuantityT *b509Op59QuantityTesting.T) {
			b509Op59QuantityBase := b509Op59QuantityKey("temperature")
			if b509Op59QuantityB524 {
				b509Op59QuantityBase.B509Op59QuantityVariant = "B524"
				b509Op59QuantityBase.B509Op59QuantitySecondary = b509Op59QuantityPtr(uint8(36))
				b509Op59QuantityBase.B509Op59QuantityGroup = b509Op59QuantityPtr(uint8(1))
				b509Op59QuantityBase.B509Op59QuantityRequestedInstance = b509Op59QuantityPtr(uint8(2))
				b509Op59QuantityBase.B509Op59QuantityRegister = b509Op59QuantityPtr(uint16(257))
				b509Op59QuantityBase.B509Op59QuantityCommandRole = b509Op59QuantityPtr("synthetic-role")
				b509Op59QuantityEntries, b509Op59QuantityError := B509Op59QuantityEntries()
				if b509Op59QuantityError != nil {
					b509Op59QuantityT.Fatal(b509Op59QuantityError)
				}
				b509Op59QuantityEntries[0].B509Op59QuantityKey = b509Op59QuantityBase
				b509Op59QuantityInstall(b509Op59QuantityT, b509Op59QuantityEntries)
			}
			b509Op59QuantityExpected := b509Op59QuantityGet(b509Op59QuantityT, b509Op59QuantityBase)
			for _, b509Op59QuantityCase := range []struct {
				b509Op59QuantityName   string
				b509Op59QuantityChange func(*B509Op59QuantityKey)
			}{
				{"family", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantityFamily = "synthetic-family"
				}},
				{"variant", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantityVariant = "synthetic-variant"
				}},
				{"primary", func(b509Op59QuantityK *B509Op59QuantityKey) { b509Op59QuantityK.B509Op59QuantityPrimary = 182 }},
				{"secondary", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantitySecondary = b509Op59QuantityPtr(uint8(10))
				}},
				{"secondary_absent", func(b509Op59QuantityK *B509Op59QuantityKey) { b509Op59QuantityK.B509Op59QuantitySecondary = nil }},
				{"secondary_zero", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantitySecondary = b509Op59QuantityPtr(uint8(0))
				}},
				{"operation", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantityOperation = b509Op59QuantityPtr(uint8(88))
				}},
				{"operation_absent", func(b509Op59QuantityK *B509Op59QuantityKey) { b509Op59QuantityK.B509Op59QuantityOperation = nil }},
				{"operation_zero", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantityOperation = b509Op59QuantityPtr(uint8(0))
				}},
				{"direction", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantityDirection = b509Op59QuantityPtr("response")
				}},
				{"field", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantityField = "synthetic-field"
				}},
				{"group", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantityGroup = b509Op59QuantityPtr(uint8(3))
				}},
				{"group_zero", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantityGroup = b509Op59QuantityPtr(uint8(0))
				}},
				{"requested_instance", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantityRequestedInstance = b509Op59QuantityPtr(uint8(3))
				}},
				{"requested_instance_zero", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantityRequestedInstance = b509Op59QuantityPtr(uint8(0))
				}},
				{"register", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantityRegister = b509Op59QuantityPtr(uint16(513))
				}},
				{"register_zero", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantityRegister = b509Op59QuantityPtr(uint16(0))
				}},
				{"command_role", func(b509Op59QuantityK *B509Op59QuantityKey) {
					b509Op59QuantityK.B509Op59QuantityCommandRole = b509Op59QuantityPtr("different-role")
				}},
			} {
				b509Op59QuantityT.Run("ab14_"+b509Op59QuantityCase.b509Op59QuantityName, func(b509Op59QuantityT *b509Op59QuantityTesting.T) {
					b509Op59QuantityChanged := b509Op59QuantityBase
					b509Op59QuantityCase.b509Op59QuantityChange(&b509Op59QuantityChanged)
					b509Op59QuantityGot, b509Op59QuantityOK, b509Op59QuantityError := B509Op59QuantityLookup(b509Op59QuantityChanged)
					if b509Op59QuantityError != nil || b509Op59QuantityOK || !b509Op59QuantityReflect.DeepEqual(b509Op59QuantityGot, B509Op59QuantityEntry{}) {
						b509Op59QuantityT.Fatalf("axis matched: %v %v", b509Op59QuantityOK, b509Op59QuantityError)
					}
					if !b509Op59QuantityReflect.DeepEqual(b509Op59QuantityGet(b509Op59QuantityT, b509Op59QuantityBase), b509Op59QuantityExpected) {
						b509Op59QuantityT.Fatal("baseline changed")
					}
				})
			}
		})
	}
}

func TestAb14AbsentZeroAndCopies(b509Op59QuantityT *b509Op59QuantityTesting.T) {
	b509Op59QuantityEntries, b509Op59QuantityError := B509Op59QuantityEntries()
	if b509Op59QuantityError != nil {
		b509Op59QuantityT.Fatal(b509Op59QuantityError)
	}
	b509Op59QuantityBaseline := b509Op59QuantityGet(b509Op59QuantityT, b509Op59QuantityKey("temperature"))
	*b509Op59QuantityEntries[0].B509Op59QuantityKey.B509Op59QuantitySecondary = 0
	*b509Op59QuantityEntries[0].B509Op59QuantityKey.B509Op59QuantityOperation = 0
	*b509Op59QuantityEntries[0].B509Op59QuantityCodec.B509Op59QuantityWidth = 1
	*b509Op59QuantityEntries[0].B509Op59QuantityCodec.B509Op59QuantityDivisor = 1
	*b509Op59QuantityEntries[0].B509Op59QuantityCodec.B509Op59QuantityUnit = "changed"
	b509Op59QuantityEntries[0].B509Op59QuantityUnknown[0] = "changed"
	b509Op59QuantityLookup := b509Op59QuantityGet(b509Op59QuantityT, b509Op59QuantityKey("temperature"))
	if !b509Op59QuantityReflect.DeepEqual(b509Op59QuantityLookup, b509Op59QuantityBaseline) {
		b509Op59QuantityT.Fatal("entries storage shared")
	}
	*b509Op59QuantityLookup.B509Op59QuantityCodec.B509Op59QuantityDivisor = 999
	*b509Op59QuantityLookup.B509Op59QuantityKey.B509Op59QuantitySecondary = 0
	b509Op59QuantityLookup.B509Op59QuantityUnknown[0] = "changed"
	b509Op59QuantityAgain, b509Op59QuantityError := B509Op59QuantityEntries()
	if b509Op59QuantityError != nil || !b509Op59QuantityReflect.DeepEqual(b509Op59QuantityAgain[0], b509Op59QuantityBaseline) || !b509Op59QuantityReflect.DeepEqual(b509Op59QuantityGet(b509Op59QuantityT, b509Op59QuantityKey("temperature")), b509Op59QuantityBaseline) {
		b509Op59QuantityT.Fatal("lookup storage shared")
	}
	for _, b509Op59QuantityAxis := range []string{"secondary", "operation", "group", "requested_instance", "register"} {
		b509Op59QuantityT.Run("ab14_"+b509Op59QuantityAxis, func(b509Op59QuantityT *b509Op59QuantityTesting.T) {
			b509Op59QuantityRows, b509Op59QuantityError := B509Op59QuantityEntries()
			if b509Op59QuantityError != nil {
				b509Op59QuantityT.Fatal(b509Op59QuantityError)
			}
			b509Op59QuantityNil := b509Op59QuantityKey("temperature")
			b509Op59QuantityZero := b509Op59QuantityNil
			switch b509Op59QuantityAxis {
			case "secondary":
				b509Op59QuantityNil.B509Op59QuantitySecondary = nil
				b509Op59QuantityZero.B509Op59QuantitySecondary = b509Op59QuantityPtr(uint8(0))
			case "operation":
				b509Op59QuantityNil.B509Op59QuantityOperation = nil
				b509Op59QuantityZero.B509Op59QuantityOperation = b509Op59QuantityPtr(uint8(0))
			case "group":
				b509Op59QuantityZero.B509Op59QuantityGroup = b509Op59QuantityPtr(uint8(0))
			case "requested_instance":
				b509Op59QuantityZero.B509Op59QuantityRequestedInstance = b509Op59QuantityPtr(uint8(0))
			case "register":
				b509Op59QuantityZero.B509Op59QuantityRegister = b509Op59QuantityPtr(uint16(0))
			}
			b509Op59QuantityRows[0].B509Op59QuantityKey = b509Op59QuantityNil
			b509Op59QuantityRows[1].B509Op59QuantityKey = b509Op59QuantityZero
			b509Op59QuantityInstall(b509Op59QuantityT, b509Op59QuantityRows)
			if b509Op59QuantityGet(b509Op59QuantityT, b509Op59QuantityNil).B509Op59QuantityCodec.B509Op59QuantityRepresentation != "signedLE16" || b509Op59QuantityGet(b509Op59QuantityT, b509Op59QuantityZero).B509Op59QuantityCodec.B509Op59QuantityRepresentation != "unsignedLE16" {
				b509Op59QuantityT.Fatal("absent and zero collapsed")
			}
		})
	}
}

func TestAb14MalformedData(b509Op59QuantityT *b509Op59QuantityTesting.T) {
	b509Op59QuantityOriginal := b509Op59QuantityData
	for _, b509Op59QuantityBad := range []string{"{", "null", "[]", b509Op59QuantityOriginal + " []"} {
		b509Op59QuantityT.Run("ab14_invalid", func(b509Op59QuantityT *b509Op59QuantityTesting.T) {
			b509Op59QuantityData = b509Op59QuantityBad
			b509Op59QuantityT.Cleanup(func() { b509Op59QuantityData = b509Op59QuantityOriginal })
			b509Op59QuantityEntries, b509Op59QuantityError := B509Op59QuantityEntries()
			if b509Op59QuantityError == nil || b509Op59QuantityEntries != nil {
				b509Op59QuantityT.Fatal("malformed entries accepted")
			}
			b509Op59QuantityEntry, b509Op59QuantityOK, b509Op59QuantityError := B509Op59QuantityLookup(b509Op59QuantityKey("temperature"))
			if b509Op59QuantityError == nil || b509Op59QuantityOK || !b509Op59QuantityReflect.DeepEqual(b509Op59QuantityEntry, B509Op59QuantityEntry{}) {
				b509Op59QuantityT.Fatal("malformed lookup succeeded")
			}
		})
	}
	for _, b509Op59QuantityMutation := range []func([]B509Op59QuantityEntry){
		func(b509Op59QuantityRows []B509Op59QuantityEntry) {
			b509Op59QuantityRows[1].B509Op59QuantityKey = b509Op59QuantityRows[0].B509Op59QuantityKey
		},
		func(b509Op59QuantityRows []B509Op59QuantityEntry) {
			b509Op59QuantityRows[0].B509Op59QuantityGrade = "Proven"
		},
		func(b509Op59QuantityRows []B509Op59QuantityEntry) {
			b509Op59QuantityRows[0].B509Op59QuantityCodec = nil
		},
		func(b509Op59QuantityRows []B509Op59QuantityEntry) {
			b509Op59QuantityRows[0].B509Op59QuantityKey.B509Op59QuantityFamily = ""
		},
	} {
		b509Op59QuantityT.Run("ab14_invalid_entry", func(b509Op59QuantityT *b509Op59QuantityTesting.T) {
			b509Op59QuantityRows, b509Op59QuantityError := B509Op59QuantityEntries()
			if b509Op59QuantityError != nil {
				b509Op59QuantityT.Fatal(b509Op59QuantityError)
			}
			b509Op59QuantityMutation(b509Op59QuantityRows)
			b509Op59QuantityInstall(b509Op59QuantityT, b509Op59QuantityRows)
			b509Op59QuantityEntry, b509Op59QuantityOK, b509Op59QuantityError := B509Op59QuantityLookup(b509Op59QuantityKey("temperature"))
			if b509Op59QuantityError == nil || b509Op59QuantityOK || !b509Op59QuantityReflect.DeepEqual(b509Op59QuantityEntry, B509Op59QuantityEntry{}) {
				b509Op59QuantityT.Fatal("invalid entry accepted")
			}
		})
	}
}
