package vaillant

import (
	b50FStatusSpanJSON "encoding/json"
	b50FStatusSpanFmt "fmt"
	b50FStatusSpanReflect "reflect"
	b50FStatusSpanTesting "testing"
)

func b50FStatusSpanPointer[b50FStatusSpanValue any](b50FStatusSpanInput b50FStatusSpanValue) *b50FStatusSpanValue {
	return &b50FStatusSpanInput
}

func b50FStatusSpanKey(b50FStatusSpanField string) B50FStatusSpanKey {
	return B50FStatusSpanKey{B50FStatusSpanFamily: "vaillant", B50FStatusSpanVariant: "B50F", B50FStatusSpanPrimary: 181, B50FStatusSpanSecondary: b50FStatusSpanPointer(byte(15)), B50FStatusSpanOperation: b50FStatusSpanPointer(byte(2)), B50FStatusSpanDirection: b50FStatusSpanPointer("request"), B50FStatusSpanField: b50FStatusSpanField}
}

func b50FStatusSpanRequire(b50FStatusSpanT *b50FStatusSpanTesting.T, b50FStatusSpanKey B50FStatusSpanKey) B50FStatusSpanEntry {
	b50FStatusSpanT.Helper()
	b50FStatusSpanEntry, b50FStatusSpanFound, b50FStatusSpanError := B50FStatusSpanLookup(b50FStatusSpanKey)
	if b50FStatusSpanError != nil || !b50FStatusSpanFound {
		b50FStatusSpanT.Fatalf("lookup: found=%v error=%v", b50FStatusSpanFound, b50FStatusSpanError)
	}
	return b50FStatusSpanEntry
}

func TestAb12FiniteContract(b50FStatusSpanT *b50FStatusSpanTesting.T) {
	b50FStatusSpanEntries, b50FStatusSpanError := B50FStatusSpanEntries()
	if b50FStatusSpanError != nil || len(b50FStatusSpanEntries) != 8 {
		b50FStatusSpanT.Fatalf("entries: %v %d", b50FStatusSpanError, len(b50FStatusSpanEntries))
	}
	b50FStatusSpanMasks := []byte{127, 127, 127, 255, 15, 7}
	for b50FStatusSpanIndex, b50FStatusSpanListed := range b50FStatusSpanEntries {
		b50FStatusSpanField := b50FStatusSpanFmt.Sprintf("status_container_%d", b50FStatusSpanIndex+1)
		if b50FStatusSpanIndex == 6 {
			b50FStatusSpanField = "burner_status"
		}
		if b50FStatusSpanIndex == 7 {
			b50FStatusSpanField = "gas_pump_status"
		}
		b50FStatusSpanEntry := b50FStatusSpanRequire(b50FStatusSpanT, b50FStatusSpanKey(b50FStatusSpanField))
		if !b50FStatusSpanReflect.DeepEqual(b50FStatusSpanEntry, b50FStatusSpanListed) || !b50FStatusSpanReflect.DeepEqual(b50FStatusSpanEntry.B50FStatusSpanKey, b50FStatusSpanKey(b50FStatusSpanField)) {
			b50FStatusSpanT.Fatal("entry identity")
		}
		if b50FStatusSpanEntry.B50FStatusSpanGrade != "Hypothesis" || b50FStatusSpanEntry.B50FStatusSpanQualification != "Unknown" || b50FStatusSpanEntry.B50FStatusSpanCoordinateGrade != "Unknown" || b50FStatusSpanEntry.B50FStatusSpanPosition != nil || b50FStatusSpanEntry.B50FStatusSpanWidthBits == nil || *b50FStatusSpanEntry.B50FStatusSpanWidthBits != 8 || b50FStatusSpanEntry.B50FStatusSpanRawFallback != "Preserve the complete independent raw byte, including reserved bits; no physical status inference." {
			b50FStatusSpanT.Fatal("evidence boundary")
		}
		b50FStatusSpanUnknown := []string{"absolute byte position", "physical qualification"}
		if b50FStatusSpanIndex < 4 {
			b50FStatusSpanUnknown = append(b50FStatusSpanUnknown, "individual allocated flag names")
		}
		if !b50FStatusSpanReflect.DeepEqual(b50FStatusSpanEntry.B50FStatusSpanUnknown, b50FStatusSpanUnknown) {
			b50FStatusSpanT.Fatal("unknown metadata")
		}
		if b50FStatusSpanIndex >= 6 {
			if b50FStatusSpanEntry.B50FStatusSpanAllocatedMask != nil || b50FStatusSpanEntry.B50FStatusSpanReservedMask != nil || len(b50FStatusSpanEntry.B50FStatusSpanBits) != 0 {
				b50FStatusSpanT.Fatal("raw byte assigned flags")
			}
			continue
		}
		if b50FStatusSpanEntry.B50FStatusSpanAllocatedMask == nil || *b50FStatusSpanEntry.B50FStatusSpanAllocatedMask != b50FStatusSpanMasks[b50FStatusSpanIndex] || b50FStatusSpanEntry.B50FStatusSpanReservedMask == nil || *b50FStatusSpanEntry.B50FStatusSpanReservedMask != ^b50FStatusSpanMasks[b50FStatusSpanIndex] || len(b50FStatusSpanEntry.B50FStatusSpanBits) != 8 {
			b50FStatusSpanT.Fatal("independent masks")
		}
		for b50FStatusSpanBitIndex, b50FStatusSpanBit := range b50FStatusSpanEntry.B50FStatusSpanBits {
			var b50FStatusSpanLabel, b50FStatusSpanCategory *string
			if b50FStatusSpanIndex == 2 && b50FStatusSpanBitIndex < 7 {
				b50FStatusSpanCategory = b50FStatusSpanPointer("temperature_limit")
			}
			if b50FStatusSpanIndex == 3 {
				b50FStatusSpanCategory = b50FStatusSpanPointer("plausibility/channel")
			}
			if b50FStatusSpanIndex == 4 && b50FStatusSpanBitIndex < 4 {
				b50FStatusSpanLabel = b50FStatusSpanPointer([]string{"valve_1_open", "valve_1_error", "valve_2_open", "valve_2_error"}[b50FStatusSpanBitIndex])
			}
			if b50FStatusSpanIndex == 5 && b50FStatusSpanBitIndex < 3 {
				b50FStatusSpanLabel = b50FStatusSpanPointer([]string{"valve_open", "error", "air_check"}[b50FStatusSpanBitIndex])
			}
			b50FStatusSpanExpected := B50FStatusSpanBit{B50FStatusSpanBit: byte(b50FStatusSpanBitIndex), B50FStatusSpanLabel: b50FStatusSpanLabel, B50FStatusSpanCategory: b50FStatusSpanCategory, B50FStatusSpanReserved: b50FStatusSpanMasks[b50FStatusSpanIndex]&(1<<b50FStatusSpanBitIndex) == 0}
			if !b50FStatusSpanReflect.DeepEqual(b50FStatusSpanBit, b50FStatusSpanExpected) {
				b50FStatusSpanT.Fatalf("container %d bit %d: %+v", b50FStatusSpanIndex+1, b50FStatusSpanBitIndex, b50FStatusSpanBit)
			}
		}
	}
}

func TestAb12SyntheticBits(b50FStatusSpanT *b50FStatusSpanTesting.T) {
	b50FStatusSpanAllocated := []byte{127, 127, 127, 255, 15, 7}
	b50FStatusSpanReserved := []byte{128, 128, 128, 0, 240, 248}
	for b50FStatusSpanContainer := 0; b50FStatusSpanContainer < 6; b50FStatusSpanContainer++ {
		for b50FStatusSpanBit := 0; b50FStatusSpanBit < 8; b50FStatusSpanBit++ {
			b50FStatusSpanT.Run(b50FStatusSpanFmt.Sprintf("ab12_container_%d_bit_%d", b50FStatusSpanContainer+1, b50FStatusSpanBit), func(b50FStatusSpanT *b50FStatusSpanTesting.T) {
				b50FStatusSpanEntry := b50FStatusSpanRequire(b50FStatusSpanT, b50FStatusSpanKey(b50FStatusSpanFmt.Sprintf("status_container_%d", b50FStatusSpanContainer+1)))
				b50FStatusSpanRaw := byte(1 << b50FStatusSpanBit)
				b50FStatusSpanTuple := []byte{0, 0, 0, 0, 0, 0, 0, 0}
				b50FStatusSpanTuple[b50FStatusSpanContainer] = b50FStatusSpanRaw
				for _, b50FStatusSpanRawIndex := range []int{6, 7} {
					b50FStatusSpanTuple[6], b50FStatusSpanTuple[7] = 0, 0
					b50FStatusSpanTuple[b50FStatusSpanRawIndex] = ^b50FStatusSpanRaw
					for b50FStatusSpanIndex, b50FStatusSpanByte := range b50FStatusSpanTuple {
						b50FStatusSpanField := b50FStatusSpanFmt.Sprintf("status_container_%d", b50FStatusSpanIndex+1)
						if b50FStatusSpanIndex == 6 {
							b50FStatusSpanField = "burner_status"
						}
						if b50FStatusSpanIndex == 7 {
							b50FStatusSpanField = "gas_pump_status"
						}
						b50FStatusSpanMetadata := b50FStatusSpanRequire(b50FStatusSpanT, b50FStatusSpanKey(b50FStatusSpanField))
						b50FStatusSpanRetained := b50FStatusSpanByte
						if b50FStatusSpanIndex < 6 {
							b50FStatusSpanAllocatedByte := b50FStatusSpanByte & *b50FStatusSpanMetadata.B50FStatusSpanAllocatedMask
							b50FStatusSpanReservedByte := b50FStatusSpanByte & *b50FStatusSpanMetadata.B50FStatusSpanReservedMask
							b50FStatusSpanExpectedByte := byte(0)
							if b50FStatusSpanIndex == b50FStatusSpanContainer {
								b50FStatusSpanExpectedByte = b50FStatusSpanRaw
							}
							if b50FStatusSpanAllocatedByte != b50FStatusSpanExpectedByte&b50FStatusSpanAllocated[b50FStatusSpanIndex] || b50FStatusSpanReservedByte != b50FStatusSpanExpectedByte&b50FStatusSpanReserved[b50FStatusSpanIndex] {
								b50FStatusSpanT.Fatal("independent tuple partitions")
							}
							b50FStatusSpanRetained = b50FStatusSpanAllocatedByte | b50FStatusSpanReservedByte
						} else if b50FStatusSpanMetadata.B50FStatusSpanKey.B50FStatusSpanField != b50FStatusSpanField || b50FStatusSpanMetadata.B50FStatusSpanAllocatedMask != nil || b50FStatusSpanMetadata.B50FStatusSpanReservedMask != nil {
							b50FStatusSpanT.Fatal("raw tuple identity")
						}
						b50FStatusSpanExpectedRaw := byte(0)
						if b50FStatusSpanIndex == b50FStatusSpanContainer {
							b50FStatusSpanExpectedRaw = b50FStatusSpanRaw
						}
						if b50FStatusSpanIndex == b50FStatusSpanRawIndex {
							b50FStatusSpanExpectedRaw = ^b50FStatusSpanRaw
						}
						if b50FStatusSpanRetained != b50FStatusSpanExpectedRaw {
							b50FStatusSpanT.Fatal("independent raw tuple lost")
						}
					}
				}
				b50FStatusSpanAllocatedValue := b50FStatusSpanRaw & *b50FStatusSpanEntry.B50FStatusSpanAllocatedMask
				b50FStatusSpanReservedValue := b50FStatusSpanRaw & *b50FStatusSpanEntry.B50FStatusSpanReservedMask
				if b50FStatusSpanAllocatedValue != b50FStatusSpanRaw&b50FStatusSpanAllocated[b50FStatusSpanContainer] || b50FStatusSpanReservedValue != b50FStatusSpanRaw&b50FStatusSpanReserved[b50FStatusSpanContainer] || b50FStatusSpanAllocatedValue|b50FStatusSpanReservedValue != b50FStatusSpanRaw {
					b50FStatusSpanT.Fatal("single bit lost or misclassified")
				}
				for _, b50FStatusSpanRawField := range []string{"burner_status", "gas_pump_status"} {
					b50FStatusSpanOther := b50FStatusSpanRequire(b50FStatusSpanT, b50FStatusSpanKey(b50FStatusSpanRawField))
					if b50FStatusSpanOther.B50FStatusSpanKey.B50FStatusSpanField != b50FStatusSpanRawField || b50FStatusSpanOther.B50FStatusSpanAllocatedMask != nil || b50FStatusSpanOther.B50FStatusSpanReservedMask != nil || len(b50FStatusSpanOther.B50FStatusSpanBits) != 0 {
						b50FStatusSpanT.Fatal("raw byte aliased with a container")
					}
					*b50FStatusSpanOther.B50FStatusSpanWidthBits = 0
					if *b50FStatusSpanEntry.B50FStatusSpanWidthBits != 8 {
						b50FStatusSpanT.Fatal("raw byte shares storage")
					}
				}
			})
		}
		for _, b50FStatusSpanRaw := range []byte{128, 127, 255, 240, 248} {
			b50FStatusSpanEntry := b50FStatusSpanRequire(b50FStatusSpanT, b50FStatusSpanKey(b50FStatusSpanFmt.Sprintf("status_container_%d", b50FStatusSpanContainer+1)))
			if b50FStatusSpanRaw&*b50FStatusSpanEntry.B50FStatusSpanAllocatedMask != b50FStatusSpanRaw&b50FStatusSpanAllocated[b50FStatusSpanContainer] || b50FStatusSpanRaw&*b50FStatusSpanEntry.B50FStatusSpanReservedMask != b50FStatusSpanRaw&b50FStatusSpanReserved[b50FStatusSpanContainer] {
				b50FStatusSpanT.Fatal("synthetic partition")
			}
		}
	}
}

func b50FStatusSpanInstall(b50FStatusSpanT *b50FStatusSpanTesting.T, b50FStatusSpanEntries []B50FStatusSpanEntry) {
	b50FStatusSpanT.Helper()
	b50FStatusSpanBefore := b50FStatusSpanStatusJSON
	b50FStatusSpanEncoded, b50FStatusSpanError := b50FStatusSpanJSON.Marshal(b50FStatusSpanEntries)
	if b50FStatusSpanError != nil {
		b50FStatusSpanT.Fatal(b50FStatusSpanError)
	}
	b50FStatusSpanStatusJSON = string(b50FStatusSpanEncoded)
	b50FStatusSpanT.Cleanup(func() { b50FStatusSpanStatusJSON = b50FStatusSpanBefore })
}

func TestAb12EveryKeyAxis(b50FStatusSpanT *b50FStatusSpanTesting.T) {
	for _, b50FStatusSpanB524Case := range []bool{false, true} {
		b50FStatusSpanT.Run(b50FStatusSpanFmt.Sprintf("ab12_b524_%v", b50FStatusSpanB524Case), func(b50FStatusSpanT *b50FStatusSpanTesting.T) {
			b50FStatusSpanBase := b50FStatusSpanKey("status_container_1")
			if b50FStatusSpanB524Case {
				b50FStatusSpanBase.B50FStatusSpanVariant = "B524"
				b50FStatusSpanBase.B50FStatusSpanSecondary = b50FStatusSpanPointer(byte(36))
				b50FStatusSpanBase.B50FStatusSpanGroup = b50FStatusSpanPointer(byte(1))
				b50FStatusSpanBase.B50FStatusSpanRequestedInstance = b50FStatusSpanPointer(byte(2))
				b50FStatusSpanBase.B50FStatusSpanRegister = b50FStatusSpanPointer(uint16(3))
				b50FStatusSpanBase.B50FStatusSpanCommandRole = b50FStatusSpanPointer("synthetic_role")
				b50FStatusSpanEntries, b50FStatusSpanError := B50FStatusSpanEntries()
				if b50FStatusSpanError != nil {
					b50FStatusSpanT.Fatal(b50FStatusSpanError)
				}
				b50FStatusSpanEntries[0].B50FStatusSpanKey = b50FStatusSpanBase
				b50FStatusSpanInstall(b50FStatusSpanT, b50FStatusSpanEntries)
			}
			b50FStatusSpanRequire(b50FStatusSpanT, b50FStatusSpanBase)
			b50FStatusSpanType := b50FStatusSpanReflect.TypeOf(b50FStatusSpanBase)
			for b50FStatusSpanIndex := 0; b50FStatusSpanIndex < b50FStatusSpanType.NumField(); b50FStatusSpanIndex++ {
				b50FStatusSpanField := b50FStatusSpanType.Field(b50FStatusSpanIndex)
				b50FStatusSpanT.Run("ab12_"+b50FStatusSpanField.Name, func(b50FStatusSpanT *b50FStatusSpanTesting.T) {
					b50FStatusSpanSource := b50FStatusSpanReflect.ValueOf(b50FStatusSpanBase).Field(b50FStatusSpanIndex)
					b50FStatusSpanVariants := []b50FStatusSpanReflect.Value{}
					if b50FStatusSpanSource.Kind() == b50FStatusSpanReflect.Pointer {
						b50FStatusSpanZero := b50FStatusSpanReflect.New(b50FStatusSpanSource.Type().Elem())
						b50FStatusSpanOther := b50FStatusSpanReflect.New(b50FStatusSpanSource.Type().Elem())
						if b50FStatusSpanOther.Elem().Kind() == b50FStatusSpanReflect.String {
							b50FStatusSpanOther.Elem().SetString("synthetic_other")
						} else {
							b50FStatusSpanOther.Elem().SetUint(99)
						}
						b50FStatusSpanVariants = append(b50FStatusSpanVariants, b50FStatusSpanReflect.Zero(b50FStatusSpanSource.Type()), b50FStatusSpanZero, b50FStatusSpanOther)
					} else {
						b50FStatusSpanOther := b50FStatusSpanReflect.New(b50FStatusSpanSource.Type()).Elem()
						if b50FStatusSpanSource.Kind() == b50FStatusSpanReflect.String {
							b50FStatusSpanOther.SetString("synthetic_other")
						} else {
							b50FStatusSpanOther.SetUint(99)
						}
						b50FStatusSpanVariants = append(b50FStatusSpanVariants, b50FStatusSpanOther)
					}
					for _, b50FStatusSpanVariant := range b50FStatusSpanVariants {
						if b50FStatusSpanReflect.DeepEqual(b50FStatusSpanSource.Interface(), b50FStatusSpanVariant.Interface()) {
							continue
						}
						b50FStatusSpanChanged := b50FStatusSpanBase
						b50FStatusSpanReflect.ValueOf(&b50FStatusSpanChanged).Elem().Field(b50FStatusSpanIndex).Set(b50FStatusSpanVariant)
						b50FStatusSpanEntry, b50FStatusSpanFound, b50FStatusSpanError := B50FStatusSpanLookup(b50FStatusSpanChanged)
						if b50FStatusSpanFound || b50FStatusSpanError != nil || !b50FStatusSpanReflect.DeepEqual(b50FStatusSpanEntry, B50FStatusSpanEntry{}) {
							b50FStatusSpanT.Fatalf("axis ignored: %s", b50FStatusSpanField.Name)
						}
					}
				})
			}
		})
	}
}

func TestAb12IndependentStorage(b50FStatusSpanT *b50FStatusSpanTesting.T) {
	b50FStatusSpanBefore, b50FStatusSpanError := B50FStatusSpanEntries()
	if b50FStatusSpanError != nil {
		b50FStatusSpanT.Fatal(b50FStatusSpanError)
	}
	b50FStatusSpanChanged, b50FStatusSpanError := B50FStatusSpanEntries()
	if b50FStatusSpanError != nil {
		b50FStatusSpanT.Fatal(b50FStatusSpanError)
	}
	for b50FStatusSpanIndex := range b50FStatusSpanChanged {
		b50FStatusSpanEntry := &b50FStatusSpanChanged[b50FStatusSpanIndex]
		*b50FStatusSpanEntry.B50FStatusSpanKey.B50FStatusSpanSecondary = 0
		*b50FStatusSpanEntry.B50FStatusSpanKey.B50FStatusSpanOperation = 0
		*b50FStatusSpanEntry.B50FStatusSpanKey.B50FStatusSpanDirection = "changed"
		*b50FStatusSpanEntry.B50FStatusSpanWidthBits = 0
		b50FStatusSpanEntry.B50FStatusSpanUnknown[0] = "changed"
		if b50FStatusSpanEntry.B50FStatusSpanAllocatedMask != nil {
			*b50FStatusSpanEntry.B50FStatusSpanAllocatedMask = 0
			*b50FStatusSpanEntry.B50FStatusSpanReservedMask = 0
		}
		for b50FStatusSpanBit := range b50FStatusSpanEntry.B50FStatusSpanBits {
			if b50FStatusSpanEntry.B50FStatusSpanBits[b50FStatusSpanBit].B50FStatusSpanLabel != nil {
				*b50FStatusSpanEntry.B50FStatusSpanBits[b50FStatusSpanBit].B50FStatusSpanLabel = "changed"
			}
			if b50FStatusSpanEntry.B50FStatusSpanBits[b50FStatusSpanBit].B50FStatusSpanCategory != nil {
				*b50FStatusSpanEntry.B50FStatusSpanBits[b50FStatusSpanBit].B50FStatusSpanCategory = "changed"
			}
			b50FStatusSpanEntry.B50FStatusSpanBits[b50FStatusSpanBit].B50FStatusSpanReserved = true
		}
	}
	for _, b50FStatusSpanExpected := range b50FStatusSpanBefore {
		b50FStatusSpanLookup := b50FStatusSpanRequire(b50FStatusSpanT, b50FStatusSpanExpected.B50FStatusSpanKey)
		if !b50FStatusSpanReflect.DeepEqual(b50FStatusSpanLookup, b50FStatusSpanExpected) {
			b50FStatusSpanT.Fatal("accessor storage shared")
		}
		*b50FStatusSpanLookup.B50FStatusSpanKey.B50FStatusSpanSecondary = 0
		*b50FStatusSpanLookup.B50FStatusSpanKey.B50FStatusSpanOperation = 0
		*b50FStatusSpanLookup.B50FStatusSpanKey.B50FStatusSpanDirection = "changed"
		*b50FStatusSpanLookup.B50FStatusSpanWidthBits = 0
		b50FStatusSpanLookup.B50FStatusSpanUnknown[0] = "changed"
		if b50FStatusSpanLookup.B50FStatusSpanAllocatedMask != nil {
			*b50FStatusSpanLookup.B50FStatusSpanAllocatedMask = 0
			*b50FStatusSpanLookup.B50FStatusSpanReservedMask = 0
		}
		for b50FStatusSpanBit := range b50FStatusSpanLookup.B50FStatusSpanBits {
			if b50FStatusSpanLookup.B50FStatusSpanBits[b50FStatusSpanBit].B50FStatusSpanLabel != nil {
				*b50FStatusSpanLookup.B50FStatusSpanBits[b50FStatusSpanBit].B50FStatusSpanLabel = "changed"
			}
			if b50FStatusSpanLookup.B50FStatusSpanBits[b50FStatusSpanBit].B50FStatusSpanCategory != nil {
				*b50FStatusSpanLookup.B50FStatusSpanBits[b50FStatusSpanBit].B50FStatusSpanCategory = "changed"
			}
		}
		if len(b50FStatusSpanLookup.B50FStatusSpanBits) > 0 {
			b50FStatusSpanLookup.B50FStatusSpanBits[0].B50FStatusSpanLabel = b50FStatusSpanPointer("changed")
		}
		if !b50FStatusSpanReflect.DeepEqual(b50FStatusSpanRequire(b50FStatusSpanT, b50FStatusSpanExpected.B50FStatusSpanKey), b50FStatusSpanExpected) {
			b50FStatusSpanT.Fatal("lookup storage shared")
		}
	}
	b50FStatusSpanAgain, b50FStatusSpanError := B50FStatusSpanEntries()
	if b50FStatusSpanError != nil || !b50FStatusSpanReflect.DeepEqual(b50FStatusSpanAgain, b50FStatusSpanBefore) {
		b50FStatusSpanT.Fatal("fresh entries mutated")
	}
}

func TestAb12MalformedData(b50FStatusSpanT *b50FStatusSpanTesting.T) {
	b50FStatusSpanBefore := b50FStatusSpanStatusJSON
	b50FStatusSpanT.Cleanup(func() { b50FStatusSpanStatusJSON = b50FStatusSpanBefore })
	b50FStatusSpanCases := []string{"{", "null", "[]", b50FStatusSpanBefore + " {}"}
	for _, b50FStatusSpanMutation := range []func([]B50FStatusSpanEntry){
		func(b50FStatusSpanEntries []B50FStatusSpanEntry) { b50FStatusSpanEntries[0] = b50FStatusSpanEntries[1] },
		func(b50FStatusSpanEntries []B50FStatusSpanEntry) {
			b50FStatusSpanEntries[0].B50FStatusSpanGrade = "Proven"
		},
		func(b50FStatusSpanEntries []B50FStatusSpanEntry) {
			b50FStatusSpanEntries[0].B50FStatusSpanAllocatedMask = b50FStatusSpanPointer(byte(255))
		},
		func(b50FStatusSpanEntries []B50FStatusSpanEntry) {
			b50FStatusSpanEntries[0].B50FStatusSpanBits[7].B50FStatusSpanLabel = b50FStatusSpanPointer("demand")
		},
		func(b50FStatusSpanEntries []B50FStatusSpanEntry) {
			b50FStatusSpanEntries[0].B50FStatusSpanWidthBits = nil
		},
		func(b50FStatusSpanEntries []B50FStatusSpanEntry) {
			b50FStatusSpanEntries[0].B50FStatusSpanAllocatedMask = nil
			b50FStatusSpanEntries[0].B50FStatusSpanReservedMask = nil
			b50FStatusSpanEntries[0].B50FStatusSpanBits = nil
		},
		func(b50FStatusSpanEntries []B50FStatusSpanEntry) {
			b50FStatusSpanEntries[0].B50FStatusSpanKey.B50FStatusSpanPrimary = 0
		},
	} {
		b50FStatusSpanStatusJSON = b50FStatusSpanBefore
		b50FStatusSpanEntries, b50FStatusSpanError := B50FStatusSpanEntries()
		if b50FStatusSpanError != nil {
			b50FStatusSpanT.Fatal(b50FStatusSpanError)
		}
		b50FStatusSpanMutation(b50FStatusSpanEntries)
		b50FStatusSpanEncoded, b50FStatusSpanError := b50FStatusSpanJSON.Marshal(b50FStatusSpanEntries)
		if b50FStatusSpanError != nil {
			b50FStatusSpanT.Fatal(b50FStatusSpanError)
		}
		b50FStatusSpanCases = append(b50FStatusSpanCases, string(b50FStatusSpanEncoded))
	}
	for b50FStatusSpanIndex, b50FStatusSpanCase := range b50FStatusSpanCases {
		b50FStatusSpanT.Run(b50FStatusSpanFmt.Sprintf("ab12_invalid_%d", b50FStatusSpanIndex), func(b50FStatusSpanT *b50FStatusSpanTesting.T) {
			b50FStatusSpanStatusJSON = b50FStatusSpanCase
			if b50FStatusSpanEntries, b50FStatusSpanError := B50FStatusSpanEntries(); b50FStatusSpanError == nil || b50FStatusSpanEntries != nil {
				b50FStatusSpanT.Fatal("invalid catalog returned entries")
			}
			if b50FStatusSpanEntry, b50FStatusSpanFound, b50FStatusSpanError := B50FStatusSpanLookup(b50FStatusSpanKey("status_container_1")); b50FStatusSpanError == nil || b50FStatusSpanFound || !b50FStatusSpanReflect.DeepEqual(b50FStatusSpanEntry, B50FStatusSpanEntry{}) {
				b50FStatusSpanT.Fatal("invalid catalog returned success")
			}
		})
	}
}

func TestAb12AbsentZeroAndB524Storage(b50FStatusSpanT *b50FStatusSpanTesting.T) {
	b50FStatusSpanEntries, b50FStatusSpanError := B50FStatusSpanEntries()
	if b50FStatusSpanError != nil {
		b50FStatusSpanT.Fatal(b50FStatusSpanError)
	}
	b50FStatusSpanBase := b50FStatusSpanKey("synthetic_identity")
	b50FStatusSpanBase.B50FStatusSpanVariant = "B524"
	b50FStatusSpanBase.B50FStatusSpanSecondary = b50FStatusSpanPointer(byte(0))
	b50FStatusSpanBase.B50FStatusSpanOperation = b50FStatusSpanPointer(byte(0))
	b50FStatusSpanBase.B50FStatusSpanGroup = b50FStatusSpanPointer(byte(0))
	b50FStatusSpanBase.B50FStatusSpanRequestedInstance = b50FStatusSpanPointer(byte(0))
	b50FStatusSpanBase.B50FStatusSpanRegister = b50FStatusSpanPointer(uint16(0))
	b50FStatusSpanBase.B50FStatusSpanCommandRole = b50FStatusSpanPointer("")
	b50FStatusSpanBase.B50FStatusSpanDirection = b50FStatusSpanPointer("")
	b50FStatusSpanEntries[0].B50FStatusSpanKey = b50FStatusSpanBase
	b50FStatusSpanInstall(b50FStatusSpanT, b50FStatusSpanEntries)
	for _, b50FStatusSpanName := range []string{"B50FStatusSpanSecondary", "B50FStatusSpanOperation", "B50FStatusSpanDirection", "B50FStatusSpanGroup", "B50FStatusSpanRequestedInstance", "B50FStatusSpanRegister", "B50FStatusSpanCommandRole"} {
		b50FStatusSpanT.Run("ab12_"+b50FStatusSpanName, func(b50FStatusSpanT *b50FStatusSpanTesting.T) {
			b50FStatusSpanRequire(b50FStatusSpanT, b50FStatusSpanBase)
			b50FStatusSpanAbsent := b50FStatusSpanBase
			b50FStatusSpanField := b50FStatusSpanReflect.ValueOf(&b50FStatusSpanAbsent).Elem().FieldByName(b50FStatusSpanName)
			b50FStatusSpanField.Set(b50FStatusSpanReflect.Zero(b50FStatusSpanField.Type()))
			if _, b50FStatusSpanFound, b50FStatusSpanError := B50FStatusSpanLookup(b50FStatusSpanAbsent); b50FStatusSpanFound || b50FStatusSpanError != nil {
				b50FStatusSpanT.Fatal("explicit zero matched absent")
			}
			b50FStatusSpanLookup := b50FStatusSpanRequire(b50FStatusSpanT, b50FStatusSpanBase)
			b50FStatusSpanPointerField := b50FStatusSpanReflect.ValueOf(&b50FStatusSpanLookup.B50FStatusSpanKey).Elem().FieldByName(b50FStatusSpanName).Elem()
			if b50FStatusSpanPointerField.Kind() == b50FStatusSpanReflect.String {
				b50FStatusSpanPointerField.SetString("changed")
			} else {
				b50FStatusSpanPointerField.SetUint(9)
			}
			if !b50FStatusSpanReflect.DeepEqual(b50FStatusSpanRequire(b50FStatusSpanT, b50FStatusSpanBase), b50FStatusSpanEntries[0]) {
				b50FStatusSpanT.Fatal("nested key storage shared")
			}
		})
	}
	b50FStatusSpanChanged, b50FStatusSpanError := B50FStatusSpanEntries()
	if b50FStatusSpanError != nil {
		b50FStatusSpanT.Fatal(b50FStatusSpanError)
	}
	*b50FStatusSpanChanged[0].B50FStatusSpanKey.B50FStatusSpanGroup = 7
	*b50FStatusSpanChanged[0].B50FStatusSpanKey.B50FStatusSpanRequestedInstance = 7
	*b50FStatusSpanChanged[0].B50FStatusSpanKey.B50FStatusSpanRegister = 7
	*b50FStatusSpanChanged[0].B50FStatusSpanKey.B50FStatusSpanCommandRole = "changed"
	if !b50FStatusSpanReflect.DeepEqual(b50FStatusSpanRequire(b50FStatusSpanT, b50FStatusSpanBase), b50FStatusSpanEntries[0]) {
		b50FStatusSpanT.Fatal("entries B524 storage shared")
	}
}
