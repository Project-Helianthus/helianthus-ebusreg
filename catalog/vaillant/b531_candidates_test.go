package vaillant

import (
	b531CandidateBinary "encoding/binary"
	b531CandidateReflect "reflect"
	b531CandidateTesting "testing"
)

func b531CandidatePointer[b531CandidateValue any](b531CandidateInput b531CandidateValue) *b531CandidateValue {
	return &b531CandidateInput
}

func b531CandidateExpectedKey(b531CandidateField string) B531CandidateKey {
	return B531CandidateKey{B531CandidateFamily: "vaillant", B531CandidateVariant: "B531", B531CandidatePrimary: 0xB5,
		B531CandidateSecondary: b531CandidatePointer(byte(0x31)), B531CandidateOperation: b531CandidatePointer(byte(2)), B531CandidateField: b531CandidateField}
}

func b531CandidateExpectedEntry(b531CandidateField string) B531CandidateEntry {
	b531CandidateEntry := B531CandidateEntry{B531CandidateKey: b531CandidateExpectedKey(b531CandidateField), B531CandidateWidthBits: b531CandidatePointer(16),
		B531CandidateRawDomain: []uint32{0, 65535}, B531CandidateSignedness: "unsigned", B531CandidateByteOrder: "little_endian",
		B531CandidateNativeUnitGrade: "Unknown", B531CandidateGrade: "Hypothesis",
		B531CandidateLocalProjection: B531CandidateProjection{B531CandidateConversion: "raw * 10", B531CandidateMultiplier: 10, B531CandidateWidthBits: b531CandidatePointer(32),
			B531CandidateDomain: []uint32{0, 655350}, B531CandidateUnit: "W", B531CandidateGrade: "Hypothesis"}}
	if b531CandidateField == "maximum" {
		b531CandidateEntry.B531CandidateMaximumTable = &B531CandidateMaximumTable{B531CandidateSourceField: "maximum", B531CandidateIndexDomain: []int{0, 8}, B531CandidateGrade: "Hypothesis"}
	}
	return b531CandidateEntry
}

func TestAb75ProductionMetadata(b531CandidateT *b531CandidateTesting.T) {
	b531CandidateCatalog, b531CandidateError := B531Candidates()
	if b531CandidateError != nil {
		b531CandidateT.Fatal(b531CandidateError)
	}
	b531CandidateExpected := B531CandidateCatalog{B531CandidateMetadataPresence: "Proven", B531CandidatePhysicalSupport: "Unknown",
		B531CandidateEndpointSupport: "Unknown", B531CandidateCorrelation: "Unknown", B531CandidateReplyAdaptation: "Unknown",
		B531CandidateCacheFreshness: "Unknown", B531CandidateTargetEligibility: "Unknown", B531CandidateAppliedLimits: "Unknown",
		B531CandidateLimitations: []string{
			"Finite candidate declarations; the broader family catalog remains incomplete.",
			"Native units remain Unknown; local W metadata does not qualify a wire unit.",
			"FFFF has no established sentinel; independent minimum and maximum permit minimum greater than maximum.",
			"The maximum-only local slot index domain 0..8 is not an installed count.",
			"No routing, invocation, polling, capability predicate, cache writer, publication or physical qualification.",
			"Pre-1.0 declarations may change.",
		}, B531CandidateEntries: []B531CandidateEntry{b531CandidateExpectedEntry("minimum"), b531CandidateExpectedEntry("maximum")}}
	if !b531CandidateReflect.DeepEqual(b531CandidateCatalog, b531CandidateExpected) {
		b531CandidateT.Fatalf("ab75 catalog differs: %#v", b531CandidateCatalog)
	}
	for _, b531CandidateField := range []string{"minimum", "maximum"} {
		b531CandidateEntry, b531CandidateFound, b531CandidateError := B531CandidateLookup(b531CandidateExpectedKey(b531CandidateField))
		if b531CandidateError != nil || !b531CandidateFound || !b531CandidateReflect.DeepEqual(b531CandidateEntry, b531CandidateExpectedEntry(b531CandidateField)) {
			b531CandidateT.Fatalf("ab75 %s exact lookup: %#v %v %v", b531CandidateField, b531CandidateEntry, b531CandidateFound, b531CandidateError)
		}
	}
}

func TestAb75ExactKeyAxes(b531CandidateT *b531CandidateTesting.T) {
	b531CandidateCases := []struct {
		b531CandidateName   string
		b531CandidateChange func(*B531CandidateKey)
	}{
		{"ab75_family", func(b531CandidateKey *B531CandidateKey) { b531CandidateKey.B531CandidateFamily = "ab75_other" }},
		{"ab75_variant", func(b531CandidateKey *B531CandidateKey) { b531CandidateKey.B531CandidateVariant = "ab75_other" }},
		{"ab75_primary", func(b531CandidateKey *B531CandidateKey) { b531CandidateKey.B531CandidatePrimary = 0 }},
		{"ab75_secondary_absent", func(b531CandidateKey *B531CandidateKey) { b531CandidateKey.B531CandidateSecondary = nil }},
		{"ab75_secondary_zero", func(b531CandidateKey *B531CandidateKey) {
			b531CandidateKey.B531CandidateSecondary = b531CandidatePointer(byte(0))
		}},
		{"ab75_secondary_other", func(b531CandidateKey *B531CandidateKey) {
			b531CandidateKey.B531CandidateSecondary = b531CandidatePointer(byte(0x32))
		}},
		{"ab75_operation_absent", func(b531CandidateKey *B531CandidateKey) { b531CandidateKey.B531CandidateOperation = nil }},
		{"ab75_operation_zero", func(b531CandidateKey *B531CandidateKey) {
			b531CandidateKey.B531CandidateOperation = b531CandidatePointer(byte(0))
		}},
		{"ab75_operation_other", func(b531CandidateKey *B531CandidateKey) {
			b531CandidateKey.B531CandidateOperation = b531CandidatePointer(byte(3))
		}},
		{"ab75_direction_empty", func(b531CandidateKey *B531CandidateKey) {
			b531CandidateKey.B531CandidateDirection = b531CandidatePointer("")
		}},
		{"ab75_direction_request", func(b531CandidateKey *B531CandidateKey) {
			b531CandidateKey.B531CandidateDirection = b531CandidatePointer("request")
		}},
		{"ab75_field", func(b531CandidateKey *B531CandidateKey) { b531CandidateKey.B531CandidateField = "ab75_other" }},
	}
	for _, b531CandidateField := range []string{"minimum", "maximum"} {
		for _, b531CandidateCase := range b531CandidateCases {
			b531CandidateT.Run("ab75_"+b531CandidateField+"/"+b531CandidateCase.b531CandidateName, func(b531CandidateT *b531CandidateTesting.T) {
				b531CandidateKey := b531CandidateExpectedKey(b531CandidateField)
				b531CandidateEntry, b531CandidateFound, b531CandidateError := B531CandidateLookup(b531CandidateKey)
				if b531CandidateError != nil || !b531CandidateFound || !b531CandidateReflect.DeepEqual(b531CandidateEntry, b531CandidateExpectedEntry(b531CandidateField)) {
					b531CandidateT.Fatal("ab75 invalid control")
				}
				b531CandidateCase.b531CandidateChange(&b531CandidateKey)
				b531CandidateEntry, b531CandidateFound, b531CandidateError = B531CandidateLookup(b531CandidateKey)
				if b531CandidateError != nil || b531CandidateFound || !b531CandidateReflect.DeepEqual(b531CandidateEntry, B531CandidateEntry{}) {
					b531CandidateT.Fatal("ab75 nonexact key matched")
				}
			})
		}
	}
	if b531CandidateEqualOptional[byte](nil, b531CandidatePointer(byte(0))) || b531CandidateEqualOptional[byte](b531CandidatePointer(byte(0)), nil) ||
		!b531CandidateEqualOptional(b531CandidatePointer(byte(0)), b531CandidatePointer(byte(0))) {
		b531CandidateT.Fatal("ab75 absent and explicit zero conflated")
	}
}

func b531CandidateSyntheticProjection(b531CandidateT *b531CandidateTesting.T, b531CandidateEntry B531CandidateEntry, b531CandidateBytes [2]byte) (uint16, uint32) {
	b531CandidateT.Helper()
	if b531CandidateEntry.B531CandidateSignedness != "unsigned" || b531CandidateEntry.B531CandidateByteOrder != "little_endian" ||
		b531CandidateEntry.B531CandidateWidthBits == nil || *b531CandidateEntry.B531CandidateWidthBits != 16 || b531CandidateEntry.B531CandidateSentinel != nil ||
		b531CandidateEntry.B531CandidateNativeUnit != nil || b531CandidateEntry.B531CandidateNativeUnitGrade != "Unknown" {
		b531CandidateT.Fatal("ab75 native declaration changed")
	}
	b531CandidateRaw := b531CandidateBinary.LittleEndian.Uint16(b531CandidateBytes[:])
	return b531CandidateRaw, uint32(b531CandidateRaw) * b531CandidateEntry.B531CandidateLocalProjection.B531CandidateMultiplier
}

func TestAb75SyntheticBoundaries(b531CandidateT *b531CandidateTesting.T) {
	b531CandidateMinimum, b531CandidateFound, b531CandidateError := B531CandidateLookup(b531CandidateExpectedKey("minimum"))
	if b531CandidateError != nil || !b531CandidateFound {
		b531CandidateT.Fatal("ab75 missing minimum")
	}
	b531CandidateMaximum, b531CandidateFound, b531CandidateError := B531CandidateLookup(b531CandidateExpectedKey("maximum"))
	if b531CandidateError != nil || !b531CandidateFound {
		b531CandidateT.Fatal("ab75 missing maximum")
	}
	b531CandidateCases := []struct {
		b531CandidateName                            string
		b531CandidateMin, b531CandidateMax           [2]byte
		b531CandidateRawMin, b531CandidateRawMax     uint16
		b531CandidateLocalMin, b531CandidateLocalMax uint32
	}{
		{"ab75_extremes", [2]byte{0, 0}, [2]byte{255, 255}, 0, 65535, 0, 655350},
		{"ab75_independent", [2]byte{2, 0}, [2]byte{1, 0}, 2, 1, 20, 10},
		{"ab75_reversed_extremes", [2]byte{255, 255}, [2]byte{0, 0}, 65535, 0, 655350, 0},
		{"ab75_byte_order", [2]byte{0, 1}, [2]byte{1, 128}, 256, 32769, 2560, 327690},
	}
	for _, b531CandidateCase := range b531CandidateCases {
		b531CandidateT.Run(b531CandidateCase.b531CandidateName, func(b531CandidateT *b531CandidateTesting.T) {
			b531CandidateMinBytes, b531CandidateMaxBytes := b531CandidateCase.b531CandidateMin, b531CandidateCase.b531CandidateMax
			b531CandidateRawMin, b531CandidateLocalMin := b531CandidateSyntheticProjection(b531CandidateT, b531CandidateMinimum, b531CandidateMinBytes)
			b531CandidateRawMax, b531CandidateLocalMax := b531CandidateSyntheticProjection(b531CandidateT, b531CandidateMaximum, b531CandidateMaxBytes)
			if b531CandidateRawMin != b531CandidateCase.b531CandidateRawMin || b531CandidateRawMax != b531CandidateCase.b531CandidateRawMax || b531CandidateLocalMin != b531CandidateCase.b531CandidateLocalMin || b531CandidateLocalMax != b531CandidateCase.b531CandidateLocalMax {
				b531CandidateT.Fatal("ab75 independent raw/projected boundaries changed")
			}
			if b531CandidateMinBytes != b531CandidateCase.b531CandidateMin || b531CandidateMaxBytes != b531CandidateCase.b531CandidateMax {
				b531CandidateT.Fatal("ab75 raw bytes altered")
			}
		})
	}
	b531CandidateCatalog, b531CandidateError := B531Candidates()
	if b531CandidateError != nil || b531CandidateCatalog.B531CandidateOrderingConstraint != nil {
		b531CandidateT.Fatal("ab75 ordering restriction introduced")
	}
	if b531CandidateMinimum.B531CandidateMaximumTable != nil || b531CandidateMaximum.B531CandidateMaximumTable == nil {
		b531CandidateT.Fatal("ab75 maximum-only association lost")
	}
	b531CandidateTable := b531CandidateMaximum.B531CandidateMaximumTable
	if b531CandidateTable.B531CandidateSourceField != "maximum" || !b531CandidateReflect.DeepEqual(b531CandidateTable.B531CandidateIndexDomain, []int{0, 8}) || b531CandidateTable.B531CandidateInstalledCount != nil {
		b531CandidateT.Fatal("ab75 slot domain changed into count or minimum")
	}
	for _, b531CandidateIndex := range []int{0, 8} {
		if b531CandidateIndex < b531CandidateTable.B531CandidateIndexDomain[0] || b531CandidateIndex > b531CandidateTable.B531CandidateIndexDomain[1] {
			b531CandidateT.Fatal("ab75 endpoint outside declared index domain")
		}
	}
	if 9 <= b531CandidateTable.B531CandidateIndexDomain[1] {
		b531CandidateT.Fatal("ab75 index 9 admitted")
	}
}

func b531CandidateMutateEntry(b531CandidateEntry *B531CandidateEntry) {
	*b531CandidateEntry.B531CandidateKey.B531CandidateSecondary = 0
	*b531CandidateEntry.B531CandidateKey.B531CandidateOperation = 0
	*b531CandidateEntry.B531CandidateWidthBits = 1
	b531CandidateEntry.B531CandidateRawDomain[1] = 1
	*b531CandidateEntry.B531CandidateLocalProjection.B531CandidateWidthBits = 1
	b531CandidateEntry.B531CandidateLocalProjection.B531CandidateDomain[1] = 1
	if b531CandidateEntry.B531CandidateMaximumTable != nil {
		b531CandidateEntry.B531CandidateMaximumTable.B531CandidateIndexDomain[1] = 1
	}
}

func TestAb75DeepCopy(b531CandidateT *b531CandidateTesting.T) {
	b531CandidateBefore, b531CandidateError := B531Candidates()
	if b531CandidateError != nil {
		b531CandidateT.Fatal(b531CandidateError)
	}
	b531CandidateChanged, b531CandidateError := B531Candidates()
	if b531CandidateError != nil {
		b531CandidateT.Fatal(b531CandidateError)
	}
	b531CandidateChanged.B531CandidateLimitations[0] = "ab75_mutated"
	for b531CandidateIndex := range b531CandidateChanged.B531CandidateEntries {
		b531CandidateMutateEntry(&b531CandidateChanged.B531CandidateEntries[b531CandidateIndex])
	}
	for _, b531CandidateField := range []string{"minimum", "maximum"} {
		b531CandidateEntry, b531CandidateFound, b531CandidateError := B531CandidateLookup(b531CandidateExpectedKey(b531CandidateField))
		if b531CandidateError != nil || !b531CandidateFound {
			b531CandidateT.Fatal("ab75 lookup missing")
		}
		b531CandidateMutateEntry(&b531CandidateEntry)
		b531CandidateFresh, b531CandidateFound, b531CandidateError := B531CandidateLookup(b531CandidateExpectedKey(b531CandidateField))
		if b531CandidateError != nil || !b531CandidateFound || !b531CandidateReflect.DeepEqual(b531CandidateFresh, b531CandidateExpectedEntry(b531CandidateField)) {
			b531CandidateT.Fatal("ab75 lookup shares nested storage")
		}
	}
	b531CandidateAfter, b531CandidateError := B531Candidates()
	if b531CandidateError != nil || !b531CandidateReflect.DeepEqual(b531CandidateBefore, b531CandidateAfter) {
		b531CandidateT.Fatal("ab75 accessor shares nested storage")
	}
}
