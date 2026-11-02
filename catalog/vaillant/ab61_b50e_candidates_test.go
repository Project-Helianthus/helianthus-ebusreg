package vaillant

import (
	ab61Reflect "reflect"
	ab61Testing "testing"
)

func ab61Pointer[ab61Value any](ab61Input ab61Value) *ab61Value { return &ab61Input }

func ab61ComparisonKey() Ab61Key {
	return Ab61Key{Ab61Family: "vaillant", Ab61Variant: "B50E", Ab61Primary: 0xB5, Ab61Secondary: ab61Pointer(byte(0x0E)), Ab61Field: "rule_comparison"}
}

func ab61ClearKey() Ab61Key {
	return Ab61Key{Ab61Family: "vaillant", Ab61Variant: "B50E", Ab61Primary: 0xB5, Ab61Secondary: ab61Pointer(byte(0x0E)), Ab61Operation: ab61Pointer(byte(1)), Ab61Direction: ab61Pointer("response"), Ab61Field: "clear_rule_positive_index"}
}

func ab61RequireEntry(ab61T *ab61Testing.T, ab61Key Ab61Key) Ab61Entry {
	ab61T.Helper()
	ab61Entry, ab61Found, ab61Error := Ab61LookupB50ECandidate(ab61Key)
	if ab61Error != nil || !ab61Found {
		ab61T.Fatalf("lookup: found=%v error=%v", ab61Found, ab61Error)
	}
	return ab61Entry
}

func TestAb61ProductionMetadata(ab61T *ab61Testing.T) {
	ab61Catalog, ab61Error := Ab61B50ECandidates()
	if ab61Error != nil {
		ab61T.Fatal(ab61Error)
	}
	if ab61Catalog.Ab61MetadataPresence != "Proven" || ab61Catalog.Ab61PhysicalQualification != "Unknown" || len(ab61Catalog.Ab61Entries) != 2 || len(ab61Catalog.Ab61Limitations) != 6 {
		ab61T.Fatalf("unexpected bounded catalog: %+v", ab61Catalog)
	}
	for ab61Index, ab61Key := range []Ab61Key{ab61ComparisonKey(), ab61ClearKey()} {
		ab61Entry := ab61RequireEntry(ab61T, ab61Key)
		if !ab61Reflect.DeepEqual(ab61Entry, ab61Catalog.Ab61Entries[ab61Index]) {
			ab61T.Fatal("lookup differs from accessor")
		}
		if ab61Entry.Ab61Grade != "Hypothesis" || ab61Entry.Ab61Position != nil || ab61Entry.Ab61WidthBits == nil || ab61Entry.Ab61RawFallback == "" {
			ab61T.Fatalf("lost grade, unknown position, carrier width or raw fallback: %+v", ab61Entry)
		}
	}
	ab61Comparison := ab61Catalog.Ab61Entries[0]
	ab61Clear := ab61Catalog.Ab61Entries[1]
	if *ab61Comparison.Ab61WidthBits != 8 || ab61Comparison.Ab61ClearIndex != nil || ab61Comparison.Ab61Comparison == nil || *ab61Clear.Ab61WidthBits != 16 || ab61Clear.Ab61Comparison != nil || ab61Clear.Ab61ClearIndex == nil {
		ab61T.Fatal("candidate alternatives merged or widths changed")
	}
	if ab61Clear.Ab61ClearIndex.Ab61ByteOrder != "little_endian" || !ab61Clear.Ab61ClearIndex.Ab61Unsigned || ab61Clear.Ab61ClearIndex.Ab61AllRulesCandidate != 65535 {
		ab61T.Fatal("clear-index metadata changed")
	}
}

func TestAb61ExactKeyAxes(ab61T *ab61Testing.T) {
	for _, ab61Control := range []Ab61Key{ab61ComparisonKey(), ab61ClearKey()} {
		ab61RequireEntry(ab61T, ab61Control)
		ab61Changes := []struct {
			ab61Name   string
			ab61Change func(*Ab61Key)
		}{
			{"ab61Family", func(ab61Key *Ab61Key) { ab61Key.Ab61Family = "other_family" }},
			{"ab61FamilyAbsent", func(ab61Key *Ab61Key) { ab61Key.Ab61Family = "" }},
			{"ab61Variant", func(ab61Key *Ab61Key) { ab61Key.Ab61Variant = "other_variant" }},
			{"ab61VariantAbsent", func(ab61Key *Ab61Key) { ab61Key.Ab61Variant = "" }},
			{"ab61Primary", func(ab61Key *Ab61Key) { ab61Key.Ab61Primary = 0 }},
			{"ab61Secondary", func(ab61Key *Ab61Key) { ab61Key.Ab61Secondary = ab61Pointer(byte(2)) }},
			{"ab61SecondaryAbsent", func(ab61Key *Ab61Key) { ab61Key.Ab61Secondary = nil }},
			{"ab61SecondaryZero", func(ab61Key *Ab61Key) { ab61Key.Ab61Secondary = ab61Pointer(byte(0)) }},
			{"ab61Operation", func(ab61Key *Ab61Key) { ab61Key.Ab61Operation = ab61Pointer(byte(2)) }},
			{"ab61OperationZero", func(ab61Key *Ab61Key) { ab61Key.Ab61Operation = ab61Pointer(byte(0)) }},
			{"ab61Direction", func(ab61Key *Ab61Key) { ab61Key.Ab61Direction = ab61Pointer("request") }},
			{"ab61DirectionEmpty", func(ab61Key *Ab61Key) { ab61Key.Ab61Direction = ab61Pointer("") }},
			{"ab61Field", func(ab61Key *Ab61Key) { ab61Key.Ab61Field = "other_field" }},
			{"ab61FieldAbsent", func(ab61Key *Ab61Key) { ab61Key.Ab61Field = "" }},
		}
		if ab61Control.Ab61Operation != nil {
			ab61Changes = append(ab61Changes, struct {
				ab61Name   string
				ab61Change func(*Ab61Key)
			}{"ab61OperationAbsent", func(ab61Key *Ab61Key) { ab61Key.Ab61Operation = nil }})
			ab61Changes = append(ab61Changes, struct {
				ab61Name   string
				ab61Change func(*Ab61Key)
			}{"ab61DirectionAbsent", func(ab61Key *Ab61Key) { ab61Key.Ab61Direction = nil }})
		}
		for _, ab61Case := range ab61Changes {
			ab61T.Run("ab61"+ab61Control.Ab61Field+"/"+ab61Case.ab61Name, func(ab61T *ab61Testing.T) {
				ab61Changed := ab61Control
				ab61Case.ab61Change(&ab61Changed)
				ab61Entry, ab61Found, ab61Error := Ab61LookupB50ECandidate(ab61Changed)
				if ab61Error != nil || ab61Found || !ab61Reflect.DeepEqual(ab61Entry, Ab61Entry{}) {
					ab61T.Fatalf("non-exact key matched: %+v found=%v error=%v", ab61Changed, ab61Found, ab61Error)
				}
			})
		}
	}
}

func TestAb61ComparisonSynthetic(ab61T *ab61Testing.T) {
	ab61Entry := ab61RequireEntry(ab61T, ab61ComparisonKey())
	if ab61Entry.Ab61Comparison == nil {
		ab61T.Fatal("missing comparison metadata")
	}
	ab61Metadata := ab61Entry.Ab61Comparison
	for _, ab61Case := range []struct {
		ab61Raw        byte
		ab61Mode       byte
		ab61Stop       bool
		ab61Unresolved byte
	}{
		{0x00, 0, false, 0}, {0x07, 7, false, 0}, {0x80, 0, true, 0}, {0x87, 7, true, 0}, {0x88, 0, true, 8},
	} {
		ab61Mode := ab61Case.ab61Raw & ab61Metadata.Ab61ModeMask
		ab61Stop := ab61Case.ab61Raw&ab61Metadata.Ab61StopOverrideFIFOMask != 0
		ab61Unresolved := ab61Case.ab61Raw & ab61Metadata.Ab61UnresolvedMask
		if ab61Mode != ab61Case.ab61Mode || ab61Stop != ab61Case.ab61Stop || ab61Unresolved != ab61Case.ab61Unresolved {
			ab61T.Fatalf("raw %02x: mode=%d stop=%v unresolved=%02x", ab61Case.ab61Raw, ab61Mode, ab61Stop, ab61Unresolved)
		}
		if ab61Mode|ab61Unresolved|(ab61Case.ab61Raw&ab61Metadata.Ab61StopOverrideFIFOMask) != ab61Case.ab61Raw {
			ab61T.Fatal("raw bits lost")
		}
	}
	for ab61Raw := 0; ab61Raw < 256; ab61Raw++ {
		if byte(ab61Raw)&ab61Metadata.Ab61ModeMask != byte(ab61Raw%8) || (byte(ab61Raw)&ab61Metadata.Ab61StopOverrideFIFOMask != 0) != (ab61Raw >= 128) || byte(ab61Raw)&ab61Metadata.Ab61UnresolvedMask != byte(ab61Raw)&0x78 {
			ab61T.Fatalf("independent scalar/flag/reserved bits changed for %02x", ab61Raw)
		}
	}
}

func TestAb61ClearEchoSynthetic(ab61T *ab61Testing.T) {
	ab61Entry := ab61RequireEntry(ab61T, ab61ClearKey())
	if ab61Entry.Ab61ClearIndex == nil || ab61Entry.Ab61WidthBits == nil {
		ab61T.Fatal("missing clear metadata")
	}
	for _, ab61Case := range []struct {
		ab61Raw      []byte
		ab61Index    uint16
		ab61Complete bool
		ab61All      bool
	}{
		{[]byte{0x34, 0x12}, 4660, true, false}, {[]byte{0xFF, 0xFF}, 65535, true, true}, {[]byte{0xFF}, 0, false, false},
	} {
		ab61Saved := append([]byte(nil), ab61Case.ab61Raw...)
		ab61Complete := len(ab61Case.ab61Raw)*8 == *ab61Entry.Ab61WidthBits
		var ab61Index uint16
		if ab61Complete {
			if ab61Entry.Ab61ClearIndex.Ab61ByteOrder != "little_endian" {
				ab61T.Fatal("not little endian")
			}
			ab61Index = uint16(ab61Case.ab61Raw[0]) | uint16(ab61Case.ab61Raw[1])<<8
		}
		ab61All := ab61Complete && ab61Index == ab61Entry.Ab61ClearIndex.Ab61AllRulesCandidate
		if ab61Complete != ab61Case.ab61Complete || ab61Index != ab61Case.ab61Index || ab61All != ab61Case.ab61All || !ab61Reflect.DeepEqual(ab61Saved, ab61Case.ab61Raw) {
			ab61T.Fatal("echo value, completeness, sentinel or raw preservation changed")
		}
	}
	ab61OtherOperation := ab61ClearKey()
	ab61OtherOperation.Ab61Operation = ab61Pointer(byte(2))
	if _, ab61Found, ab61Error := Ab61LookupB50ECandidate(ab61OtherOperation); ab61Error != nil || ab61Found {
		ab61T.Fatal("FFFF clear candidate escaped OP01")
	}
}

func TestAb61FreshDeepCopies(ab61T *ab61Testing.T) {
	ab61Before, ab61Error := Ab61B50ECandidates()
	if ab61Error != nil {
		ab61T.Fatal(ab61Error)
	}
	ab61Changed, ab61Error := Ab61B50ECandidates()
	if ab61Error != nil {
		ab61T.Fatal(ab61Error)
	}
	ab61Changed.Ab61Limitations[0] = "changed"
	*ab61Changed.Ab61Entries[0].Ab61Key.Ab61Secondary = 0
	*ab61Changed.Ab61Entries[0].Ab61WidthBits = 0
	ab61Changed.Ab61Entries[0].Ab61Comparison.Ab61ModeMask = 0
	*ab61Changed.Ab61Entries[1].Ab61Key.Ab61Operation = 0
	*ab61Changed.Ab61Entries[1].Ab61Key.Ab61Direction = "changed"
	ab61Changed.Ab61Entries[1].Ab61ClearIndex.Ab61AllRulesCandidate = 0
	ab61Changed.Ab61Entries[1].Ab61RawFallback = "changed"
	ab61Lookup := ab61RequireEntry(ab61T, ab61ClearKey())
	*ab61Lookup.Ab61WidthBits = 0
	*ab61Lookup.Ab61Key.Ab61Secondary = 0
	*ab61Lookup.Ab61Key.Ab61Operation = 0
	*ab61Lookup.Ab61Key.Ab61Direction = "changed"
	ab61Lookup.Ab61ClearIndex.Ab61ByteOrder = "changed"
	ab61ComparisonLookup := ab61RequireEntry(ab61T, ab61ComparisonKey())
	ab61ComparisonLookup.Ab61Comparison.Ab61UnresolvedMask = 0
	ab61After, ab61Error := Ab61B50ECandidates()
	if ab61Error != nil || !ab61Reflect.DeepEqual(ab61Before, ab61After) {
		ab61T.Fatal("accessor shares backing storage")
	}
	if !ab61Reflect.DeepEqual(ab61RequireEntry(ab61T, ab61ClearKey()), ab61Before.Ab61Entries[1]) || !ab61Reflect.DeepEqual(ab61RequireEntry(ab61T, ab61ComparisonKey()), ab61Before.Ab61Entries[0]) {
		ab61T.Fatal("lookup shares backing storage")
	}
}
