package vaillant

import (
	ab92Reflect "reflect"
	ab92Testing "testing"
)

func ab92Ptr[Ab92Value any](ab92Value Ab92Value) *Ab92Value { return &ab92Value }

func ab92Key(ab92Operation uint16, ab92Alternative, ab92Name string) Ab92Key {
	ab92Direction := "response"
	if ab92Alternative == "request" {
		ab92Direction = "request"
	}
	return Ab92Key{Ab92Family: "Vaillant", Ab92Variant: "Unknown", Ab92Primary: 181,
		Ab92Secondary: ab92Ptr(uint16(32)), Ab92Operation: ab92Ptr(ab92Operation),
		Ab92Direction: ab92Direction, Ab92Alternative: ab92Alternative, Ab92Field: ab92Name}
}

func ab92Get(ab92T *ab92Testing.T, ab92Operation uint16, ab92Alternative, ab92Name string) Ab92B520Field {
	ab92T.Helper()
	ab92Field, ab92Found, ab92Error := Ab92LookupB520(ab92Key(ab92Operation, ab92Alternative, ab92Name))
	if ab92Error != nil || !ab92Found {
		ab92T.Fatalf("ab92 lookup %d/%s/%s: %v %v", ab92Operation, ab92Alternative, ab92Name, ab92Found, ab92Error)
	}
	return ab92Field
}

func TestAb92ProductionMetadata(ab92T *ab92Testing.T) {
	ab92Catalog, ab92Error := Ab92B520Blocks()
	if ab92Error != nil {
		ab92T.Fatal(ab92Error)
	}
	if ab92Catalog.Ab92RepresentationGrade != "Proven" || ab92Catalog.Ab92PhysicalQualification != "Unknown" ||
		ab92Catalog.Ab92FamilyCompleteness != "incomplete" || ab92Catalog.Ab92Invocation != "unchanged" ||
		ab92Catalog.Ab92Selection != "explicit-key-only" || ab92Catalog.Ab92NegativeForm != "object-not-physical-NAK" ||
		ab92Catalog.Ab92LogicalPositions != "include-prefix;not-physical-offsets" {
		ab92T.Fatal("ab92 admission/representation boundary")
	}
	ab92Unknowns := map[string]string{"status_semantics": "Unknown", "block_correlation": "Unknown", "atomicity": "Unknown", "physical_adaptation": "Unknown", "effects": "Unknown", "support": "Unknown", "freshness": "Unknown", "persistence": "Unknown", "execution_permission": "Unknown"}
	if !ab92Reflect.DeepEqual(ab92Catalog.Ab92Unknowns, ab92Unknowns) {
		ab92T.Fatal("ab92 unknown claims")
	}
	ab92Expected := []Ab92B520Field{}
	ab92Add := func(ab92Operation uint16, ab92Alternative, ab92Name string, ab92Position, ab92Width int, ab92Max uint64, ab92Reserved bool) {
		ab92Order := "not-applicable"
		if ab92Width > 1 {
			ab92Order = "little-endian"
		}
		ab92Expected = append(ab92Expected, Ab92B520Field{Ab92Key: ab92Key(ab92Operation, ab92Alternative, ab92Name), Ab92Prefix: []int{181, 32, int(ab92Operation)},
			Ab92LogicalPosition: ab92Ptr(ab92Position), Ab92WidthBytes: ab92Ptr(ab92Width), Ab92Signedness: "unsigned", Ab92ByteOrder: ab92Order,
			Ab92Conversion: "unchanged", Ab92Default: ab92Ptr(uint64(0)), Ab92Minimum: ab92Ptr(uint64(0)), Ab92Maximum: ab92Ptr(ab92Max),
			Ab92Enumeration: map[string]string{}, Ab92Reserved: ab92Reserved, Ab92Grade: "Hypothesis", Ab92Semantics: "Unknown", Ab92RawFallback: "preserve-raw-number"})
	}
	ab92Data := func(ab92Operation uint16, ab92Alternative string, ab92Position int, ab92Minimum, ab92Maximum *int, ab92Default *uint64) {
		ab92Expected = append(ab92Expected, Ab92B520Field{Ab92Key: ab92Key(ab92Operation, ab92Alternative, "data"), Ab92Prefix: []int{181, 32, int(ab92Operation)},
			Ab92LogicalPosition: ab92Ptr(ab92Position), Ab92Signedness: "not-applicable", Ab92ByteOrder: "not-applicable", Ab92Conversion: "unchanged",
			Ab92Default: ab92Default, Ab92MinimumBytes: ab92Minimum, Ab92MaximumBytes: ab92Maximum, Ab92Extent: ab92Ptr("through-end-of-PDU"),
			Ab92Enumeration: map[string]string{}, Ab92Grade: "Hypothesis", Ab92Semantics: "Unknown", Ab92RawFallback: "preserve-raw-bytes"})
	}
	for _, ab92Alternative := range []string{"request", "positive"} {
		ab92Add(1, ab92Alternative, "status", 3, 1, 255, false)
	}
	ab92Add(2, "request", "status", 3, 1, 255, false)
	ab92Add(2, "request", "address", 4, 4, 4294967295, false)
	ab92Add(2, "request", "flags", 8, 1, 255, false)
	ab92Add(2, "request", "size-bits", 9, 1, 255, false)
	ab92Add(2, "request", "block", 10, 1, 254, false)
	ab92Add(2, "request", "start-bit", 11, 1, 255, false)
	ab92Add(2, "positive", "status", 3, 1, 255, false)
	for _, ab92Alternative := range []string{"request", "positive"} {
		ab92Add(3, ab92Alternative, "status", 3, 1, 255, false)
		ab92Add(3, ab92Alternative, "block", 4, 1, 254, false)
	}
	ab92Data(3, "positive", 5, ab92Ptr(0), nil, nil)
	ab92Add(4, "request", "status", 3, 1, 255, false)
	ab92Add(4, "request", "block", 4, 1, 255, false)
	ab92Data(4, "request", 5, nil, nil, nil)
	ab92Add(4, "positive", "status", 3, 1, 255, false)
	ab92Add(5, "request", "status", 3, 1, 255, false)
	for ab92Index, ab92Name := range []string{"status", "version", "max-block-length", "max-variable-count", "current-variable-count", "max-block-count", "current-block-count", "error-code"} {
		ab92Add(5, "positive", ab92Name, 3+ab92Index, 1, 255, false)
	}
	for _, ab92Alternative := range []string{"request", "positive"} {
		ab92Max := uint64(254)
		if ab92Alternative == "positive" {
			ab92Max = 255
		}
		ab92Add(6, ab92Alternative, "status", 3, 1, 255, false)
		ab92Add(6, ab92Alternative, "write-block", 4, 1, ab92Max, false)
		ab92Add(6, ab92Alternative, "read-block", 5, 1, ab92Max, false)
		ab92Data(6, ab92Alternative, 6, nil, nil, nil)
	}
	for _, ab92Alternative := range []string{"request", "positive"} {
		ab92Add(7, ab92Alternative, "reserved-identifier", 3, 4, 4294967295, true)
		ab92Add(7, ab92Alternative, "reserved-offset", 7, 2, 65535, true)
		ab92Data(7, ab92Alternative, 9, ab92Ptr(0), ab92Ptr(23), ab92Ptr(uint64(0)))
	}
	for _, ab92Operation := range []uint16{2, 3, 4, 6} {
		ab92Add(ab92Operation, "negative", "status", 3, 1, 255, false)
		if ab92Operation == 6 {
			ab92Add(6, "negative", "write-error", 4, 1, 255, false)
			ab92Add(6, "negative", "read-error", 5, 1, 255, false)
		} else {
			ab92Add(ab92Operation, "negative", "error", 4, 1, 255, false)
		}
	}
	if len(ab92Expected) != 50 || len(ab92Catalog.Ab92Fields) != 50 {
		ab92T.Fatal("ab92 exact finite field count")
	}
	for _, ab92Want := range ab92Expected {
		ab92Got, ab92Found, ab92Error := Ab92LookupB520(ab92Want.Ab92Key)
		if ab92Error != nil || !ab92Found || !ab92Reflect.DeepEqual(ab92Got, ab92Want) {
			ab92T.Errorf("ab92 field mismatch %+v: got %+v, found %v, error %v", ab92Want.Ab92Key, ab92Got, ab92Found, ab92Error)
		}
	}
}

func ab92Fits(ab92Field Ab92B520Field, ab92Value uint64) bool {
	return ab92Field.Ab92Minimum != nil && ab92Field.Ab92Maximum != nil && ab92Value >= *ab92Field.Ab92Minimum && ab92Value <= *ab92Field.Ab92Maximum
}

func ab92Read(ab92T *ab92Testing.T, ab92Field Ab92B520Field, ab92Raw []byte) uint64 {
	ab92T.Helper()
	if ab92Field.Ab92LogicalPosition == nil || ab92Field.Ab92WidthBytes == nil {
		ab92T.Fatal("ab92 no scalar extent")
	}
	ab92Position, ab92Width := *ab92Field.Ab92LogicalPosition, *ab92Field.Ab92WidthBytes
	if ab92Position < 0 || ab92Width < 1 || ab92Width > 8 || ab92Position+ab92Width > len(ab92Raw) {
		ab92T.Fatal("ab92 invalid synthetic scalar extent")
	}
	if ab92Width > 1 && ab92Field.Ab92ByteOrder != "little-endian" {
		ab92T.Fatal("ab92 unexpected scalar order")
	}
	ab92Value := uint64(0)
	for ab92Index := 0; ab92Index < ab92Width; ab92Index++ {
		ab92Value |= uint64(ab92Raw[ab92Position+ab92Index]) << uint(8*ab92Index)
	}
	return ab92Value
}

func TestAb92SyntheticCases(ab92T *ab92Testing.T) {
	for _, ab92Case := range []struct {
		ab92Op                    uint16
		ab92Alternative, ab92Name string
		ab92Value                 uint64
		ab92Want                  bool
	}{
		{3, "request", "block", 254, true}, {3, "request", "block", 255, false}, {4, "request", "block", 255, true},
		{6, "request", "write-block", 255, false}, {6, "request", "read-block", 255, false}, {6, "positive", "write-block", 255, true}, {6, "positive", "read-block", 255, true},
	} {
		ab92Field := ab92Get(ab92T, ab92Case.ab92Op, ab92Case.ab92Alternative, ab92Case.ab92Name)
		if ab92Fits(ab92Field, ab92Case.ab92Value) != ab92Case.ab92Want {
			ab92T.Errorf("ab92 domain %+v", ab92Case)
		}
		if len(ab92Field.Ab92Enumeration) != 0 || ab92Field.Ab92Replacement != nil || ab92Field.Ab92Conversion != "unchanged" {
			ab92T.Fatal("ab92 invented sentinel")
		}
	}
	ab92Address := ab92Get(ab92T, 2, "request", "address")
	if ab92Read(ab92T, ab92Address, []byte{181, 32, 2, 0, 0x78, 0x56, 0x34, 0x12}) != 0x12345678 || ab92Address.Ab92Semantics != "Unknown" || ab92Address.Ab92Conversion != "unchanged" {
		ab92T.Fatal("ab92 raw address")
	}
	ab92Raw := []byte{181, 32, 5, 0, 9, 2, 3, 200, 4, 201, 255}
	for ab92Index, ab92Name := range []string{"status", "version", "max-block-length", "max-variable-count", "current-variable-count", "max-block-count", "current-block-count", "error-code"} {
		ab92Field := ab92Get(ab92T, 5, "positive", ab92Name)
		ab92Want := []uint64{0, 9, 2, 3, 200, 4, 201, 255}[ab92Index]
		if ab92Read(ab92T, ab92Field, ab92Raw) != ab92Want || !ab92Fits(ab92Field, ab92Want) || ab92Field.Ab92Conversion != "unchanged" || len(ab92Field.Ab92Enumeration) != 0 {
			ab92T.Fatal("ab92 independent OP05 raw fields")
		}
	}
	if !ab92Fits(ab92Get(ab92T, 4, "request", "block"), 255) {
		ab92T.Fatal("ab92 OP05 cannot impose a global maximum")
	}
	for _, ab92Alternative := range []string{"request", "positive"} {
		ab92Field := ab92Get(ab92T, 7, ab92Alternative, "data")
		for _, ab92Length := range []int{23, 24} {
			ab92Raw := make([]byte, 9+ab92Length)
			ab92Count := len(ab92Raw) - *ab92Field.Ab92LogicalPosition
			ab92Fits := ab92Count >= *ab92Field.Ab92MinimumBytes && ab92Count <= *ab92Field.Ab92MaximumBytes
			if ab92Fits != (ab92Length == 23) {
				ab92T.Fatal("ab92 OP07 abstract extent")
			}
		}
		for _, ab92Name := range []string{"reserved-identifier", "reserved-offset"} {
			ab92Field := ab92Get(ab92T, 7, ab92Alternative, ab92Name)
			if !ab92Field.Ab92Reserved || ab92Field.Ab92Semantics != "Unknown" || ab92Field.Ab92Conversion != "unchanged" {
				ab92T.Fatal("ab92 reserved is not a pointer")
			}
		}
	}
	ab92Negative := []byte{181, 32, 6, 0, 1, 2}
	if ab92Read(ab92T, ab92Get(ab92T, 6, "negative", "write-error"), ab92Negative) != 1 || ab92Read(ab92T, ab92Get(ab92T, 6, "negative", "read-error"), ab92Negative) != 2 {
		ab92T.Fatal("ab92 independent error bytes")
	}
	ab92Positive := ab92Get(ab92T, 6, "positive", "status")
	ab92NegativeField := ab92Get(ab92T, 6, "negative", "status")
	if !ab92Reflect.DeepEqual(ab92Positive.Ab92Prefix, ab92NegativeField.Ab92Prefix) || ab92SameKey(ab92Positive.Ab92Key, ab92NegativeField.Ab92Key) {
		ab92T.Fatal("ab92 prefix does not select branch")
	}
	ab92Missing := ab92Positive.Ab92Key
	ab92Missing.Ab92Alternative = ""
	if _, ab92Found, _ := Ab92LookupB520(ab92Missing); ab92Found {
		ab92T.Fatal("ab92 automatic branch selection")
	}
}

func TestAb92ExactIdentity(ab92T *ab92Testing.T) {
	ab92Control := ab92Key(2, "request", "status")
	ab92Mutations := []struct {
		ab92Name   string
		ab92Change func(*Ab92Key)
	}{
		{"family", func(ab92Key *Ab92Key) { ab92Key.Ab92Family = "other" }},
		{"variant", func(ab92Key *Ab92Key) { ab92Key.Ab92Variant = "other" }},
		{"primary", func(ab92Key *Ab92Key) { ab92Key.Ab92Primary = 182 }},
		{"secondary", func(ab92Key *Ab92Key) { ab92Key.Ab92Secondary = ab92Ptr(uint16(33)) }},
		{"secondary-absent", func(ab92Key *Ab92Key) { ab92Key.Ab92Secondary = nil }},
		{"secondary-zero", func(ab92Key *Ab92Key) { ab92Key.Ab92Secondary = ab92Ptr(uint16(0)) }},
		{"operation", func(ab92Key *Ab92Key) { ab92Key.Ab92Operation = ab92Ptr(uint16(0)) }},
		{"operation-absent", func(ab92Key *Ab92Key) { ab92Key.Ab92Operation = nil }},
		{"direction", func(ab92Key *Ab92Key) { ab92Key.Ab92Direction = "response" }},
		{"alternative", func(ab92Key *Ab92Key) { ab92Key.Ab92Alternative = "negative" }},
		{"field", func(ab92Key *Ab92Key) { ab92Key.Ab92Field = "other" }},
		{"opcode-zero", func(ab92Key *Ab92Key) { ab92Key.Ab92Opcode = ab92Ptr(uint16(0)) }},
		{"group-zero", func(ab92Key *Ab92Key) { ab92Key.Ab92Group = ab92Ptr(uint16(0)) }},
		{"instance-zero", func(ab92Key *Ab92Key) { ab92Key.Ab92Instance = ab92Ptr(uint16(0)) }},
		{"register-zero", func(ab92Key *Ab92Key) { ab92Key.Ab92Register = ab92Ptr(uint16(0)) }},
	}
	if _, ab92Found, ab92Error := Ab92LookupB520(ab92Control); !ab92Found || ab92Error != nil {
		ab92T.Fatal("ab92 valid control")
	}
	for _, ab92Mutation := range ab92Mutations {
		ab92T.Run("ab92-"+ab92Mutation.ab92Name, func(ab92T *ab92Testing.T) {
			ab92Changed := ab92Control
			ab92Mutation.ab92Change(&ab92Changed)
			if _, ab92Found, ab92Error := Ab92LookupB520(ab92Changed); ab92Found || ab92Error != nil {
				ab92T.Fatal("ab92 axis ignored")
			}
		})
	}
	ab92B524 := ab92Control
	ab92B524.Ab92Secondary = ab92Ptr(uint16(36))
	ab92B524.Ab92Opcode = ab92Ptr(uint16(1))
	ab92B524.Ab92Group = ab92Ptr(uint16(2))
	ab92B524.Ab92Instance = ab92Ptr(uint16(3))
	ab92B524.Ab92Register = ab92Ptr(uint16(4))
	ab92B524Copy := ab92B524
	ab92B524Copy.Ab92Opcode = ab92Ptr(uint16(1))
	if !ab92SameKey(ab92B524, ab92B524Copy) {
		ab92T.Fatal("ab92 pointer identity used instead of value")
	}
	for ab92Axis := 0; ab92Axis < 4; ab92Axis++ {
		for _, ab92Value := range []*uint16{nil, ab92Ptr(uint16(0)), ab92Ptr(uint16(99))} {
			ab92Changed := ab92B524
			ab92Axes := []**uint16{&ab92Changed.Ab92Opcode, &ab92Changed.Ab92Group, &ab92Changed.Ab92Instance, &ab92Changed.Ab92Register}
			*ab92Axes[ab92Axis] = ab92Value
			if ab92SameKey(ab92B524, ab92Changed) {
				ab92T.Fatalf("ab92 qualified B524 axis %d ignored", ab92Axis)
			}
		}
	}
	ab92Absent := ab92Control
	ab92Absent.Ab92Operation = nil
	ab92Zero := ab92Absent
	ab92Zero.Ab92Operation = ab92Ptr(uint16(0))
	if ab92SameKey(ab92Absent, ab92Zero) {
		ab92T.Fatal("ab92 absent operation equals zero")
	}
	ab92Absent = ab92Control
	ab92Absent.Ab92Secondary = nil
	ab92Zero = ab92Absent
	ab92Zero.Ab92Secondary = ab92Ptr(uint16(0))
	if ab92SameKey(ab92Absent, ab92Zero) {
		ab92T.Fatal("ab92 absent secondary equals zero")
	}
}

func TestAb92DeepCopies(ab92T *ab92Testing.T) {
	ab92Before, ab92Error := Ab92B520Blocks()
	if ab92Error != nil {
		ab92T.Fatal(ab92Error)
	}
	ab92Changed, _ := Ab92B520Blocks()
	ab92Changed.Ab92Unknowns["support"] = "Proven"
	for ab92Index := range ab92Changed.Ab92Fields {
		ab92Mutate(&ab92Changed.Ab92Fields[ab92Index])
	}
	ab92After, ab92Error := Ab92B520Blocks()
	if ab92Error != nil || !ab92Reflect.DeepEqual(ab92Before, ab92After) {
		ab92T.Fatal("ab92 accessor shares storage")
	}
	for _, ab92Original := range ab92Before.Ab92Fields {
		ab92Changed, ab92Found, ab92Error := Ab92LookupB520(ab92Original.Ab92Key)
		if !ab92Found || ab92Error != nil {
			ab92T.Fatal("ab92 lookup control")
		}
		ab92Mutate(&ab92Changed)
		ab92Again, ab92Found, ab92Error := Ab92LookupB520(ab92Original.Ab92Key)
		if !ab92Found || ab92Error != nil || !ab92Reflect.DeepEqual(ab92Original, ab92Again) {
			ab92T.Fatal("ab92 lookup shares storage")
		}
	}
}

func ab92Mutate(ab92Field *Ab92B520Field) {
	ab92Field.Ab92Prefix[0] = 0
	*ab92Field.Ab92Key.Ab92Secondary = 0
	*ab92Field.Ab92Key.Ab92Operation = 0
	*ab92Field.Ab92LogicalPosition = 99
	if ab92Field.Ab92WidthBytes != nil {
		*ab92Field.Ab92WidthBytes = 99
	}
	if ab92Field.Ab92Default != nil {
		*ab92Field.Ab92Default = 99
	}
	if ab92Field.Ab92Minimum != nil {
		*ab92Field.Ab92Minimum = 99
	}
	if ab92Field.Ab92Maximum != nil {
		*ab92Field.Ab92Maximum = 99
	}
	if ab92Field.Ab92MinimumBytes != nil {
		*ab92Field.Ab92MinimumBytes = 99
	}
	if ab92Field.Ab92MaximumBytes != nil {
		*ab92Field.Ab92MaximumBytes = 99
	}
	if ab92Field.Ab92Extent != nil {
		*ab92Field.Ab92Extent = "changed"
	}
	ab92Field.Ab92Enumeration["255"] = "sentinel"
	ab92Field.Ab92RawFallback = "changed"
}
