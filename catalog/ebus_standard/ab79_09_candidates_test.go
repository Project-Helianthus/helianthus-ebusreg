package ebus_standard_catalog

import (
	ab79binary "encoding/binary"
	ab79reflect "reflect"
	ab79testing "testing"
)

func ab79Byte(ab79Value uint8) *uint8 { return &ab79Value }

func ab79Key(ab79Variant string, ab79Secondary uint8, ab79Direction, ab79Field string) Ab79Key {
	return Ab79Key{Ab79Family: "memory-service", Ab79Variant: ab79Variant, Ab79Primary: ab79Byte(9), Ab79Secondary: ab79Byte(ab79Secondary), Ab79Direction: ab79Direction, Ab79Field: ab79Field}
}

func ab79Get(ab79T *ab79testing.T, ab79Key Ab79Key) Ab79Entry {
	ab79T.Helper()
	ab79Entry, ab79OK := Ab79Lookup(ab79Key)
	if !ab79OK {
		ab79T.Fatalf("missing key: %+v", ab79Key)
	}
	return ab79Entry
}

func ab79Int(ab79T *ab79testing.T, ab79Value *int, ab79Want int) {
	ab79T.Helper()
	if ab79Value == nil || *ab79Value != ab79Want {
		ab79T.Fatalf("integer %v, want %d", ab79Value, ab79Want)
	}
}

func TestAb79ExactKeys(ab79T *ab79testing.T) {
	ab79Control := ab79Key("0902", 2, "abstract-positive", "data")
	ab79Get(ab79T, ab79Control)
	ab79Cases := []struct {
		Ab79Name   string
		Ab79Mutate func(*Ab79Key)
	}{
		{"family", func(ab79K *Ab79Key) { ab79K.Ab79Family = "unknown" }},
		{"variant", func(ab79K *Ab79Key) { ab79K.Ab79Variant = "0901" }},
		{"primary", func(ab79K *Ab79Key) { ab79K.Ab79Primary = ab79Byte(8) }},
		{"primary-absent", func(ab79K *Ab79Key) { ab79K.Ab79Primary = nil }},
		{"primary-zero", func(ab79K *Ab79Key) { ab79K.Ab79Primary = ab79Byte(0) }},
		{"secondary", func(ab79K *Ab79Key) { ab79K.Ab79Secondary = ab79Byte(3) }},
		{"secondary-absent", func(ab79K *Ab79Key) { ab79K.Ab79Secondary = nil }},
		{"secondary-zero", func(ab79K *Ab79Key) { ab79K.Ab79Secondary = ab79Byte(0) }},
		{"operation-zero", func(ab79K *Ab79Key) { ab79K.Ab79Operation = ab79Byte(0) }},
		{"operation", func(ab79K *Ab79Key) { ab79K.Ab79Operation = ab79Byte(1) }},
		{"direction", func(ab79K *Ab79Key) { ab79K.Ab79Direction = "abstract-request" }},
		{"field", func(ab79K *Ab79Key) { ab79K.Ab79Field = "address" }},
	}
	for _, ab79Case := range ab79Cases {
		ab79T.Run("ab79-"+ab79Case.Ab79Name, func(ab79T *ab79testing.T) {
			ab79K := ab79Control
			ab79Case.Ab79Mutate(&ab79K)
			if _, ab79OK := Ab79Lookup(ab79K); ab79OK {
				ab79T.Fatal("changed axis matched")
			}
		})
	}
	ab79Zero := ab79Key("0900", 0, "abstract-positive", "data")
	ab79Get(ab79T, ab79Zero)
	ab79Zero.Ab79Secondary = nil
	if _, ab79OK := Ab79Lookup(ab79Zero); ab79OK {
		ab79T.Fatal("absent secondary matched explicit zero")
	}
}

func TestAb79Inventory(ab79T *ab79testing.T) {
	ab79Catalog := Ab79Candidates()
	if len(ab79Catalog.Ab79Entries) != 25 || ab79Catalog.Ab79PresenceGrade != "Proven" || ab79Catalog.Ab79UnknownGrade != "Unknown" || len(ab79Catalog.Ab79Unknowns) != 8 {
		ab79T.Fatal("inventory or grade changed")
	}
	ab79Seen := make(map[string]bool)
	for _, ab79Entry := range ab79Catalog.Ab79Entries {
		if ab79Entry.Ab79Grade != "Hypothesis" || ab79Entry.Ab79PhysicalOffset != nil || ab79Entry.Ab79Unit != nil || ab79Entry.Ab79Enum != nil || ab79Entry.Ab79Key.Ab79Operation != nil {
			ab79T.Fatal("invented qualification or metadata")
		}
		ab79Found := ab79Get(ab79T, ab79Entry.Ab79Key)
		if !ab79reflect.DeepEqual(ab79Found, ab79Entry) {
			ab79T.Fatal("lookup differs from accessor")
		}
		ab79ID := ab79Entry.Ab79Key.Ab79Variant + ab79Entry.Ab79Key.Ab79Direction + ab79Entry.Ab79Key.Ab79Field
		if ab79Seen[ab79ID] {
			ab79T.Fatal("duplicate identity")
		}
		ab79Seen[ab79ID] = true
	}
	for ab79V, ab79Variant := range []string{"0900", "0901", "0902", "0903"} {
		for ab79Offset, ab79Field := range []string{"PB", "SB"} {
			ab79Entry := ab79Get(ab79T, ab79Key(ab79Variant, uint8(ab79V), "abstract-request", ab79Field))
			ab79Int(ab79T, ab79Entry.Ab79AbstractOffset, ab79Offset)
			ab79Int(ab79T, ab79Entry.Ab79Width, 1)
			if ab79V != 2 {
				ab79Entry = ab79Get(ab79T, ab79Key(ab79Variant, uint8(ab79V), "abstract-positive", ab79Field))
				ab79Int(ab79T, ab79Entry.Ab79AbstractOffset, ab79Offset)
				ab79Int(ab79T, ab79Entry.Ab79Width, 1)
			} else if _, ab79OK := Ab79Lookup(ab79Key(ab79Variant, 2, "abstract-positive", ab79Field)); ab79OK {
				ab79T.Fatal("0902 invented prefix")
			}
		}
	}
}

func TestAb79AddressAlternatives(ab79T *ab79testing.T) {
	ab79Raw := []byte{0x34, 0x12}
	for ab79V, ab79Variant := range []string{"0900", "0901", "0902", "0903"} {
		if ab79V < 2 {
			ab79Entry := ab79Get(ab79T, ab79Key(ab79Variant, uint8(ab79V), "abstract-request", "address"))
			if ab79Entry.Ab79Encoding == nil || *ab79Entry.Ab79Encoding != "unsigned-le16" || ab79binary.LittleEndian.Uint16(ab79Raw) != 4660 {
				ab79T.Fatal("LE16 address changed")
			}
			ab79Int(ab79T, ab79Entry.Ab79Width, 2)
			ab79Int(ab79T, ab79Entry.Ab79Maximum, 65535)
			ab79Int(ab79T, ab79Entry.Ab79Default, 0)
			ab79Int(ab79T, ab79Entry.Ab79Minimum, 0)
			if ab79Entry.Ab79AbstractOffset != nil {
				ab79T.Fatal("invented address offset")
			}
		} else {
			for ab79Index, ab79Field := range []string{"address_low", "address_high"} {
				ab79Entry := ab79Get(ab79T, ab79Key(ab79Variant, uint8(ab79V), "abstract-request", ab79Field))
				if ab79Entry.Ab79Encoding == nil || *ab79Entry.Ab79Encoding != "unsigned-u8" || ab79Raw[ab79Index] != []byte{52, 18}[ab79Index] {
					ab79T.Fatal("split address changed")
				}
				ab79Int(ab79T, ab79Entry.Ab79Width, 1)
				ab79Int(ab79T, ab79Entry.Ab79Maximum, 255)
				ab79Int(ab79T, ab79Entry.Ab79Default, 0)
				ab79Int(ab79T, ab79Entry.Ab79Minimum, 0)
				if ab79Entry.Ab79AbstractOffset != nil {
					ab79T.Fatal("invented split address offset")
				}
			}
			if _, ab79OK := Ab79Lookup(ab79Key(ab79Variant, uint8(ab79V), "abstract-request", "address")); ab79OK {
				ab79T.Fatal("merged split address")
			}
		}
	}
}

func TestAb79PositiveAlternatives(ab79T *ab79testing.T) {
	for _, ab79Case := range []struct {
		Ab79Variant   string
		Ab79Secondary uint8
		Ab79Offset    int
		Ab79Raw       []byte
	}{
		{"0900", 0, 2, []byte{9, 0, 0xaa, 0x55}}, {"0902", 2, 0, []byte{0xaa, 0x55}}, {"0903", 3, 2, []byte{9, 3, 0xaa, 0x55}},
	} {
		ab79Entry := ab79Get(ab79T, ab79Key(ab79Case.Ab79Variant, ab79Case.Ab79Secondary, "abstract-positive", "data"))
		ab79Int(ab79T, ab79Entry.Ab79AbstractOffset, ab79Case.Ab79Offset)
		ab79Int(ab79T, ab79Entry.Ab79Minimum, 0)
		if ab79Entry.Ab79Width != nil || ab79Entry.Ab79Maximum != nil || ab79Entry.Ab79Encoding == nil || *ab79Entry.Ab79Encoding != "opaque-to-end" {
			ab79T.Fatal("opaque metadata changed")
		}
		if !ab79reflect.DeepEqual(ab79Case.Ab79Raw[*ab79Entry.Ab79AbstractOffset:], []byte{0xaa, 0x55}) {
			ab79T.Fatal("incorrect stripping")
		}
		ab79Empty := ab79Case.Ab79Raw[:ab79Case.Ab79Offset]
		if len(ab79Empty[*ab79Entry.Ab79AbstractOffset:]) != 0 {
			ab79T.Fatal("empty opaque data lost")
		}
	}
	for _, ab79Field := range []string{"data", "success"} {
		if _, ab79OK := Ab79Lookup(ab79Key("0901", 1, "abstract-positive", ab79Field)); ab79OK {
			ab79T.Fatal("prefix-only invented body or success")
		}
	}
}

func TestAb79CountConflict(ab79T *ab79testing.T) {
	for _, ab79V := range []uint8{0, 2} {
		ab79Variant := map[uint8]string{0: "0900", 2: "0902"}[ab79V]
		ab79Entry := ab79Get(ab79T, ab79Key(ab79Variant, ab79V, "abstract-request", "count"))
		ab79Int(ab79T, ab79Entry.Ab79Maximum, 255)
		ab79Int(ab79T, ab79Entry.Ab79Minimum, 0)
		ab79Int(ab79T, ab79Entry.Ab79Default, 0)
		ab79Int(ab79T, ab79Entry.Ab79Width, 1)
		if ab79Entry.Ab79AbstractOffset != nil || ab79Entry.Ab79Encoding == nil || *ab79Entry.Ab79Encoding != "unsigned-u8" {
			ab79T.Fatal("count metadata changed")
		}
		if ab79V == 2 {
			ab79Int(ab79T, ab79Entry.Ab79AlternativeMaximum, 10)
		} else if ab79Entry.Ab79AlternativeMaximum != nil {
			ab79T.Fatal("conflict transferred")
		}
		for _, ab79Raw := range []byte{0, 10, 11, 255} {
			if int(ab79Raw) < *ab79Entry.Ab79Minimum || int(ab79Raw) > *ab79Entry.Ab79Maximum {
				ab79T.Fatalf("raw count %d clipped", ab79Raw)
			}
		}
	}
}

func TestAb79DeepCopies(ab79T *ab79testing.T) {
	ab79Baseline := Ab79Candidates()
	ab79Changed := Ab79Candidates()
	ab79Changed.Ab79Unknowns[0] = "changed"
	for ab79Index := range ab79Changed.Ab79Entries {
		ab79Entry := &ab79Changed.Ab79Entries[ab79Index]
		*ab79Entry.Ab79Key.Ab79Primary = 0
		*ab79Entry.Ab79Key.Ab79Secondary = 99
		for _, ab79Pointer := range []*int{ab79Entry.Ab79AbstractOffset, ab79Entry.Ab79Width, ab79Entry.Ab79Default, ab79Entry.Ab79Minimum, ab79Entry.Ab79Maximum, ab79Entry.Ab79AlternativeMaximum} {
			if ab79Pointer != nil {
				*ab79Pointer = 99
			}
		}
		if ab79Entry.Ab79Encoding != nil {
			*ab79Entry.Ab79Encoding = "changed"
		}
		if len(ab79Entry.Ab79Notes) > 0 {
			ab79Entry.Ab79Notes[0] = "changed"
		}
	}
	ab79Lookup := ab79Get(ab79T, ab79Key("0902", 2, "abstract-request", "count"))
	*ab79Lookup.Ab79Key.Ab79Secondary = 0
	*ab79Lookup.Ab79AlternativeMaximum = 99
	ab79Lookup.Ab79Notes[0] = "changed"
	if !ab79reflect.DeepEqual(Ab79Candidates(), ab79Baseline) {
		ab79T.Fatal("shared nested storage")
	}
	ab79Fresh := ab79Get(ab79T, ab79Key("0902", 2, "abstract-request", "count"))
	ab79Int(ab79T, ab79Fresh.Ab79AlternativeMaximum, 10)
	if ab79Fresh.Ab79Notes[0] == "changed" {
		ab79T.Fatal("lookup shares notes")
	}
}
