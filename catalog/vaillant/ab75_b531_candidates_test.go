package vaillant

import (
	ab75Binary "encoding/binary"
	ab75Reflect "reflect"
	ab75Testing "testing"
)

func ab75Pointer[ab75Value any](ab75Input ab75Value) *ab75Value { return &ab75Input }

func ab75ExpectedKey(ab75Field string) Ab75Key {
	return Ab75Key{Ab75Family: "vaillant", Ab75Variant: "B531", Ab75Primary: 0xB5,
		Ab75Secondary: ab75Pointer(byte(0x31)), Ab75Operation: ab75Pointer(byte(2)), Ab75Field: ab75Field}
}

func ab75ExpectedEntry(ab75Field string) Ab75Entry {
	ab75Entry := Ab75Entry{Ab75Key: ab75ExpectedKey(ab75Field), Ab75WidthBits: ab75Pointer(16),
		Ab75RawDomain: []uint32{0, 65535}, Ab75Signedness: "unsigned", Ab75ByteOrder: "little_endian",
		Ab75NativeUnitGrade: "Unknown", Ab75Grade: "Hypothesis",
		Ab75LocalProjection: Ab75Projection{Ab75Conversion: "raw * 10", Ab75Multiplier: 10, Ab75WidthBits: ab75Pointer(32),
			Ab75Domain: []uint32{0, 655350}, Ab75Unit: "W", Ab75Grade: "Hypothesis"}}
	if ab75Field == "maximum" {
		ab75Entry.Ab75MaximumTable = &Ab75MaximumTable{Ab75SourceField: "maximum", Ab75IndexDomain: []int{0, 8}, Ab75Grade: "Hypothesis"}
	}
	return ab75Entry
}

func TestAb75ProductionMetadata(ab75T *ab75Testing.T) {
	ab75Catalog, ab75Error := Ab75B531Candidates()
	if ab75Error != nil {
		ab75T.Fatal(ab75Error)
	}
	ab75Expected := Ab75Catalog{Ab75MetadataPresence: "Proven", Ab75PhysicalSupport: "Unknown",
		Ab75EndpointSupport: "Unknown", Ab75Correlation: "Unknown", Ab75ReplyAdaptation: "Unknown",
		Ab75CacheFreshness: "Unknown", Ab75TargetEligibility: "Unknown", Ab75AppliedLimits: "Unknown",
		Ab75Limitations: []string{
			"Finite candidate declarations; the broader family catalog remains incomplete.",
			"Native units remain Unknown; local W metadata does not qualify a wire unit.",
			"FFFF has no established sentinel; independent minimum and maximum permit minimum greater than maximum.",
			"The maximum-only local slot index domain 0..8 is not an installed count.",
			"No routing, invocation, polling, capability predicate, cache writer, publication or physical qualification.",
			"Pre-1.0 declarations may change.",
		}, Ab75Entries: []Ab75Entry{ab75ExpectedEntry("minimum"), ab75ExpectedEntry("maximum")}}
	if !ab75Reflect.DeepEqual(ab75Catalog, ab75Expected) {
		ab75T.Fatalf("ab75 catalog differs: %#v", ab75Catalog)
	}
	for _, ab75Field := range []string{"minimum", "maximum"} {
		ab75Entry, ab75Found, ab75Error := Ab75LookupB531(ab75ExpectedKey(ab75Field))
		if ab75Error != nil || !ab75Found || !ab75Reflect.DeepEqual(ab75Entry, ab75ExpectedEntry(ab75Field)) {
			ab75T.Fatalf("ab75 %s exact lookup: %#v %v %v", ab75Field, ab75Entry, ab75Found, ab75Error)
		}
	}
}

func TestAb75ExactKeyAxes(ab75T *ab75Testing.T) {
	ab75Cases := []struct {
		ab75Name   string
		ab75Change func(*Ab75Key)
	}{
		{"ab75_family", func(ab75Key *Ab75Key) { ab75Key.Ab75Family = "ab75_other" }},
		{"ab75_variant", func(ab75Key *Ab75Key) { ab75Key.Ab75Variant = "ab75_other" }},
		{"ab75_primary", func(ab75Key *Ab75Key) { ab75Key.Ab75Primary = 0 }},
		{"ab75_secondary_absent", func(ab75Key *Ab75Key) { ab75Key.Ab75Secondary = nil }},
		{"ab75_secondary_zero", func(ab75Key *Ab75Key) { ab75Key.Ab75Secondary = ab75Pointer(byte(0)) }},
		{"ab75_secondary_other", func(ab75Key *Ab75Key) { ab75Key.Ab75Secondary = ab75Pointer(byte(0x32)) }},
		{"ab75_operation_absent", func(ab75Key *Ab75Key) { ab75Key.Ab75Operation = nil }},
		{"ab75_operation_zero", func(ab75Key *Ab75Key) { ab75Key.Ab75Operation = ab75Pointer(byte(0)) }},
		{"ab75_operation_other", func(ab75Key *Ab75Key) { ab75Key.Ab75Operation = ab75Pointer(byte(3)) }},
		{"ab75_direction_empty", func(ab75Key *Ab75Key) { ab75Key.Ab75Direction = ab75Pointer("") }},
		{"ab75_direction_request", func(ab75Key *Ab75Key) { ab75Key.Ab75Direction = ab75Pointer("request") }},
		{"ab75_field", func(ab75Key *Ab75Key) { ab75Key.Ab75Field = "ab75_other" }},
	}
	for _, ab75Field := range []string{"minimum", "maximum"} {
		for _, ab75Case := range ab75Cases {
			ab75T.Run("ab75_"+ab75Field+"/"+ab75Case.ab75Name, func(ab75T *ab75Testing.T) {
				ab75Key := ab75ExpectedKey(ab75Field)
				ab75Entry, ab75Found, ab75Error := Ab75LookupB531(ab75Key)
				if ab75Error != nil || !ab75Found || !ab75Reflect.DeepEqual(ab75Entry, ab75ExpectedEntry(ab75Field)) {
					ab75T.Fatal("ab75 invalid control")
				}
				ab75Case.ab75Change(&ab75Key)
				ab75Entry, ab75Found, ab75Error = Ab75LookupB531(ab75Key)
				if ab75Error != nil || ab75Found || !ab75Reflect.DeepEqual(ab75Entry, Ab75Entry{}) {
					ab75T.Fatal("ab75 nonexact key matched")
				}
			})
		}
	}
	if ab75EqualOptional[byte](nil, ab75Pointer(byte(0))) || ab75EqualOptional[byte](ab75Pointer(byte(0)), nil) ||
		!ab75EqualOptional(ab75Pointer(byte(0)), ab75Pointer(byte(0))) {
		ab75T.Fatal("ab75 absent and explicit zero conflated")
	}
}

func ab75SyntheticProjection(ab75T *ab75Testing.T, ab75Entry Ab75Entry, ab75Bytes [2]byte) (uint16, uint32) {
	ab75T.Helper()
	if ab75Entry.Ab75Signedness != "unsigned" || ab75Entry.Ab75ByteOrder != "little_endian" ||
		ab75Entry.Ab75WidthBits == nil || *ab75Entry.Ab75WidthBits != 16 || ab75Entry.Ab75Sentinel != nil ||
		ab75Entry.Ab75NativeUnit != nil || ab75Entry.Ab75NativeUnitGrade != "Unknown" {
		ab75T.Fatal("ab75 native declaration changed")
	}
	ab75Raw := ab75Binary.LittleEndian.Uint16(ab75Bytes[:])
	return ab75Raw, uint32(ab75Raw) * ab75Entry.Ab75LocalProjection.Ab75Multiplier
}

func TestAb75SyntheticBoundaries(ab75T *ab75Testing.T) {
	ab75Minimum, ab75Found, ab75Error := Ab75LookupB531(ab75ExpectedKey("minimum"))
	if ab75Error != nil || !ab75Found {
		ab75T.Fatal("ab75 missing minimum")
	}
	ab75Maximum, ab75Found, ab75Error := Ab75LookupB531(ab75ExpectedKey("maximum"))
	if ab75Error != nil || !ab75Found {
		ab75T.Fatal("ab75 missing maximum")
	}
	ab75Cases := []struct {
		ab75Name                   string
		ab75Min, ab75Max           [2]byte
		ab75RawMin, ab75RawMax     uint16
		ab75LocalMin, ab75LocalMax uint32
	}{
		{"ab75_extremes", [2]byte{0, 0}, [2]byte{255, 255}, 0, 65535, 0, 655350},
		{"ab75_independent", [2]byte{2, 0}, [2]byte{1, 0}, 2, 1, 20, 10},
		{"ab75_reversed_extremes", [2]byte{255, 255}, [2]byte{0, 0}, 65535, 0, 655350, 0},
		{"ab75_byte_order", [2]byte{0, 1}, [2]byte{1, 128}, 256, 32769, 2560, 327690},
	}
	for _, ab75Case := range ab75Cases {
		ab75T.Run(ab75Case.ab75Name, func(ab75T *ab75Testing.T) {
			ab75MinBytes, ab75MaxBytes := ab75Case.ab75Min, ab75Case.ab75Max
			ab75RawMin, ab75LocalMin := ab75SyntheticProjection(ab75T, ab75Minimum, ab75MinBytes)
			ab75RawMax, ab75LocalMax := ab75SyntheticProjection(ab75T, ab75Maximum, ab75MaxBytes)
			if ab75RawMin != ab75Case.ab75RawMin || ab75RawMax != ab75Case.ab75RawMax || ab75LocalMin != ab75Case.ab75LocalMin || ab75LocalMax != ab75Case.ab75LocalMax {
				ab75T.Fatal("ab75 independent raw/projected boundaries changed")
			}
			if ab75MinBytes != ab75Case.ab75Min || ab75MaxBytes != ab75Case.ab75Max {
				ab75T.Fatal("ab75 raw bytes altered")
			}
		})
	}
	ab75Catalog, ab75Error := Ab75B531Candidates()
	if ab75Error != nil || ab75Catalog.Ab75OrderingConstraint != nil {
		ab75T.Fatal("ab75 ordering restriction introduced")
	}
	if ab75Minimum.Ab75MaximumTable != nil || ab75Maximum.Ab75MaximumTable == nil {
		ab75T.Fatal("ab75 maximum-only association lost")
	}
	ab75Table := ab75Maximum.Ab75MaximumTable
	if ab75Table.Ab75SourceField != "maximum" || !ab75Reflect.DeepEqual(ab75Table.Ab75IndexDomain, []int{0, 8}) || ab75Table.Ab75InstalledCount != nil {
		ab75T.Fatal("ab75 slot domain changed into count or minimum")
	}
	for _, ab75Index := range []int{0, 8} {
		if ab75Index < ab75Table.Ab75IndexDomain[0] || ab75Index > ab75Table.Ab75IndexDomain[1] {
			ab75T.Fatal("ab75 endpoint outside declared index domain")
		}
	}
	if 9 <= ab75Table.Ab75IndexDomain[1] {
		ab75T.Fatal("ab75 index 9 admitted")
	}
}

func ab75MutateEntry(ab75Entry *Ab75Entry) {
	*ab75Entry.Ab75Key.Ab75Secondary = 0
	*ab75Entry.Ab75Key.Ab75Operation = 0
	*ab75Entry.Ab75WidthBits = 1
	ab75Entry.Ab75RawDomain[1] = 1
	*ab75Entry.Ab75LocalProjection.Ab75WidthBits = 1
	ab75Entry.Ab75LocalProjection.Ab75Domain[1] = 1
	if ab75Entry.Ab75MaximumTable != nil {
		ab75Entry.Ab75MaximumTable.Ab75IndexDomain[1] = 1
	}
}

func TestAb75DeepCopy(ab75T *ab75Testing.T) {
	ab75Before, ab75Error := Ab75B531Candidates()
	if ab75Error != nil {
		ab75T.Fatal(ab75Error)
	}
	ab75Changed, ab75Error := Ab75B531Candidates()
	if ab75Error != nil {
		ab75T.Fatal(ab75Error)
	}
	ab75Changed.Ab75Limitations[0] = "ab75_mutated"
	for ab75Index := range ab75Changed.Ab75Entries {
		ab75MutateEntry(&ab75Changed.Ab75Entries[ab75Index])
	}
	for _, ab75Field := range []string{"minimum", "maximum"} {
		ab75Entry, ab75Found, ab75Error := Ab75LookupB531(ab75ExpectedKey(ab75Field))
		if ab75Error != nil || !ab75Found {
			ab75T.Fatal("ab75 lookup missing")
		}
		ab75MutateEntry(&ab75Entry)
		ab75Fresh, ab75Found, ab75Error := Ab75LookupB531(ab75ExpectedKey(ab75Field))
		if ab75Error != nil || !ab75Found || !ab75Reflect.DeepEqual(ab75Fresh, ab75ExpectedEntry(ab75Field)) {
			ab75T.Fatal("ab75 lookup shares nested storage")
		}
	}
	ab75After, ab75Error := Ab75B531Candidates()
	if ab75Error != nil || !ab75Reflect.DeepEqual(ab75Before, ab75After) {
		ab75T.Fatal("ab75 accessor shares nested storage")
	}
}
