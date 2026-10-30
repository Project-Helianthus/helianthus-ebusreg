package vaillant

import (
	ab96json "encoding/json"
	ab96reflect "reflect"
	ab96testing "testing"
)

func ab96Pointer[ab96Type any](ab96Value ab96Type) *ab96Type { return &ab96Value }

func ab96Find(ab96Test *ab96testing.T, ab96Variant string, ab96Operation *uint8, ab96Field string) Ab96Entry {
	ab96Test.Helper()
	ab96Entry, ab96Found, ab96Error := Ab96Lookup(Ab96Key{Ab96Family: "vaillant", Ab96Variant: ab96Variant, Ab96Primary: 181, Ab96Secondary: ab96Pointer(uint8(3)), Ab96Operation: ab96Operation, Ab96Field: ab96Field})
	if ab96Error != nil || !ab96Found {
		ab96Test.Fatalf("lookup %s/%v/%s: %v %v", ab96Variant, ab96Operation, ab96Field, ab96Found, ab96Error)
	}
	return ab96Entry
}

func TestAb96Facts(ab96Test *ab96testing.T) {
	ab96Entries, ab96Error := Ab96Entries()
	if ab96Error != nil || len(ab96Entries) != 9 {
		ab96Test.Fatalf("finite entries: %d %v", len(ab96Entries), ab96Error)
	}
	for _, ab96Entry := range ab96Entries {
		if ab96Entry.Ab96Grade != "Hypothesis" || ab96Entry.Ab96MetadataPresence != "Proven" || ab96Entry.Ab96PhysicalQualification != "Unknown" || ab96Entry.Ab96Position != nil || ab96Entry.Ab96Padding != nil || ab96Entry.Ab96Key.Ab96Direction != nil || ab96Entry.Ab96RawPreservation != "Retain raw values and partial last-known-good fields independently." {
			ab96Test.Fatalf("qualification/unknown/raw boundary: %+v", ab96Entry)
		}
		if !ab96reflect.DeepEqual(ab96Entry.Ab96Unknown, []string{"absolute coordinates", "termination on the wire", "success", "session exclusivity", "physical effects", "prefix-only positive form does not establish an empty physical response", "direction"}) {
			ab96Test.Fatalf("unknown boundaries: %v", ab96Entry.Ab96Unknown)
		}
	}
	for _, ab96Operation := range []uint8{3, 6} {
		ab96Unlock := ab96Find(ab96Test, "common", ab96Pointer(ab96Operation), "unlock")
		if ab96Unlock.Ab96WidthCharacters == nil || *ab96Unlock.Ab96WidthCharacters != 8 || ab96Unlock.Ab96Encoding != nil || ab96Unlock.Ab96WidthBits != nil {
			ab96Test.Fatal("unlock must declare eight characters without encoding or byte width")
		}
		for _, ab96Text := range []string{"1234567", "12345678", "123456789"} {
			if (len(ab96Text) == *ab96Unlock.Ab96WidthCharacters) != (ab96Text == "12345678") {
				ab96Test.Fatal("declared character width")
			}
		}
		ab96Trailing := ab96Find(ab96Test, "common", ab96Pointer(ab96Operation), "unlock_trailing_byte")
		if ab96Trailing.Ab96WidthBits == nil || *ab96Trailing.Ab96WidthBits != 8 || ab96Trailing.Ab96Unsigned == nil || !*ab96Trailing.Ab96Unsigned || ab96Trailing.Ab96Default == nil || *ab96Trailing.Ab96Default != 0 || ab96Trailing.Ab96Follows == nil || *ab96Trailing.Ab96Follows != "unlock" {
			ab96Test.Fatal("separate trailing unsigned byte/default")
		}
	}
	for _, ab96Operation := range []uint8{4, 5} {
		ab96Text := ab96Find(ab96Test, "common", ab96Pointer(ab96Operation), "text")
		if ab96Text.Ab96Encoding == nil || *ab96Text.Ab96Encoding != "ASCII" || ab96Text.Ab96Extent == nil || *ab96Text.Ab96Extent != "through end of data extent" || ab96Text.Ab96WidthBits != nil || ab96Text.Ab96WidthCharacters != nil {
			ab96Test.Fatal("variable ASCII extent")
		}
		ab96Bytes := []byte{0x41, 0x42}
		if *ab96Text.Ab96Encoding != "ASCII" || string(ab96Bytes) != "AB" {
			ab96Test.Fatal("ASCII AB fixture")
		}
	}
	for _, ab96Text := range []string{"", "AB", "ABCDE"} {
		ab96Mode := ab96Find(ab96Test, "common", ab96Pointer(uint8(5)), "mode")
		if ab96Mode.Ab96Position != nil || ab96Mode.Ab96Follows != nil || ab96Mode.Ab96WidthBits == nil || *ab96Mode.Ab96WidthBits != 8 || ab96Mode.Ab96Unsigned == nil || !*ab96Mode.Ab96Unsigned || ab96Mode.Ab96Default != nil || ab96Mode.Ab96Labels != nil {
			ab96Test.Fatalf("mode byte/unknown coordinate for text length %d", len(ab96Text))
		}
	}
	ab96Common := ab96Find(ab96Test, "common", nil, "categories")
	ab96Reset := ab96Find(ab96Test, "reset", nil, "categories")
	if !ab96reflect.DeepEqual(ab96Common.Ab96Labels, map[uint8]*string{1: nil, 2: ab96Pointer("predictive maintenance"), 3: nil, 4: nil, 5: nil}) || !ab96reflect.DeepEqual(ab96Reset.Ab96Labels, map[uint8]*string{1: ab96Pointer("error"), 3: ab96Pointer("status")}) {
		ab96Test.Fatal("variant-specific category labels")
	}
	if ab96Common.Ab96WidthBits != nil || ab96Reset.Ab96WidthBits != nil {
		ab96Test.Fatal("unprovided category width")
	}
}

func TestAb96ExactKeyAxes(ab96Test *ab96testing.T) {
	ab96Saved := ab96Data
	ab96Test.Cleanup(func() { ab96Data = ab96Saved })
	ab96Key := Ab96Key{Ab96Family: "vaillant", Ab96Variant: "synthetic", Ab96Primary: 181, Ab96Secondary: ab96Pointer(uint8(36)), Ab96Operation: ab96Pointer(uint8(3)), Ab96Direction: ab96Pointer("request"), Ab96Field: "synthetic", Ab96Group: ab96Pointer(uint8(1)), Ab96RequestedInstance: ab96Pointer(uint8(2)), Ab96Register: ab96Pointer(uint16(3)), Ab96CommandRole: ab96Pointer("read")}
	ab96Entry := ab96Find(ab96Test, "common", ab96Pointer(uint8(3)), "unlock")
	ab96Entry.Ab96Key = ab96Key
	var ab96Error error
	ab96Data, ab96Error = ab96json.Marshal([]Ab96Entry{ab96Entry})
	if ab96Error != nil {
		ab96Test.Fatal(ab96Error)
	}
	ab96Got, ab96Found, ab96Error := Ab96Lookup(ab96Key)
	if ab96Error != nil || !ab96Found || !ab96reflect.DeepEqual(ab96Got.Ab96Key, ab96Key) {
		ab96Test.Fatal("synthetic B524 exact key", ab96Error)
	}
	for _, ab96Case := range []struct {
		ab96Name   string
		ab96Change func(*Ab96Key)
	}{
		{"family", func(ab96Key *Ab96Key) { ab96Key.Ab96Family = "other" }},
		{"variant", func(ab96Key *Ab96Key) { ab96Key.Ab96Variant = "other" }},
		{"primary", func(ab96Key *Ab96Key) { ab96Key.Ab96Primary = 180 }},
		{"secondary", func(ab96Key *Ab96Key) { ab96Key.Ab96Secondary = ab96Pointer(uint8(35)) }},
		{"operation", func(ab96Key *Ab96Key) { ab96Key.Ab96Operation = ab96Pointer(uint8(4)) }},
		{"direction", func(ab96Key *Ab96Key) { ab96Key.Ab96Direction = ab96Pointer("response") }},
		{"field", func(ab96Key *Ab96Key) { ab96Key.Ab96Field = "other" }},
		{"group", func(ab96Key *Ab96Key) { ab96Key.Ab96Group = ab96Pointer(uint8(2)) }},
		{"requested_instance", func(ab96Key *Ab96Key) { ab96Key.Ab96RequestedInstance = ab96Pointer(uint8(3)) }},
		{"register", func(ab96Key *Ab96Key) { ab96Key.Ab96Register = ab96Pointer(uint16(4)) }},
		{"command_role", func(ab96Key *Ab96Key) { ab96Key.Ab96CommandRole = ab96Pointer("write") }},
	} {
		ab96Test.Run(ab96Case.ab96Name, func(ab96Test *ab96testing.T) {
			ab96Changed := ab96Key
			ab96Case.ab96Change(&ab96Changed)
			ab96Got, ab96Found, ab96Error := Ab96Lookup(ab96Changed)
			if ab96Error != nil || ab96Found || !ab96reflect.DeepEqual(ab96Got, Ab96Entry{}) {
				ab96Test.Fatal("changed axis matched", ab96Error)
			}
		})
	}
}

func TestAb96AbsentZero(ab96Test *ab96testing.T) {
	ab96Saved := ab96Data
	ab96Test.Cleanup(func() { ab96Data = ab96Saved })
	ab96Entry := ab96Find(ab96Test, "common", nil, "categories")
	ab96Axes := []string{"Ab96Secondary", "Ab96Operation", "Ab96Group", "Ab96RequestedInstance", "Ab96Register", "Ab96Direction", "Ab96CommandRole"}
	for _, ab96Axis := range ab96Axes {
		ab96Test.Run(ab96Axis, func(ab96Test *ab96testing.T) {
			ab96Absent := ab96Entry
			ab96reflect.ValueOf(&ab96Absent.Ab96Key).Elem().FieldByName(ab96Axis).SetZero()
			ab96Zero := ab96Absent
			ab96Field := ab96reflect.ValueOf(&ab96Zero.Ab96Key).Elem().FieldByName(ab96Axis)
			ab96Field.Set(ab96reflect.New(ab96Field.Type().Elem()))
			for _, ab96Stored := range []Ab96Entry{ab96Absent, ab96Zero} {
				var ab96Error error
				ab96Data, ab96Error = ab96json.Marshal([]Ab96Entry{ab96Stored})
				if ab96Error != nil {
					ab96Test.Fatal(ab96Error)
				}
				for _, ab96Query := range []Ab96Key{ab96Absent.Ab96Key, ab96Zero.Ab96Key} {
					_, ab96Found, ab96Error := Ab96Lookup(ab96Query)
					if ab96Error != nil || ab96Found != ab96reflect.DeepEqual(ab96Query, ab96Stored.Ab96Key) {
						ab96Test.Fatal("absent/zero collapsed", ab96Error)
					}
				}
			}
		})
	}
}

func ab96Mutate(ab96Value ab96reflect.Value) {
	switch ab96Value.Kind() {
	case ab96reflect.Pointer:
		if !ab96Value.IsNil() {
			ab96Mutate(ab96Value.Elem())
		}
	case ab96reflect.Struct:
		for ab96Index := 0; ab96Index < ab96Value.NumField(); ab96Index++ {
			ab96Mutate(ab96Value.Field(ab96Index))
		}
	case ab96reflect.Slice:
		for ab96Index := 0; ab96Index < ab96Value.Len(); ab96Index++ {
			ab96Mutate(ab96Value.Index(ab96Index))
		}
	case ab96reflect.Map:
		for _, ab96Key := range ab96Value.MapKeys() {
			ab96Mutate(ab96Value.MapIndex(ab96Key))
		}
		ab96Value.Clear()
	case ab96reflect.String:
		if ab96Value.CanSet() {
			ab96Value.SetString("changed")
		}
	case ab96reflect.Uint8, ab96reflect.Uint16:
		if ab96Value.CanSet() {
			ab96Value.SetUint(0)
		}
	case ab96reflect.Int:
		if ab96Value.CanSet() {
			ab96Value.SetInt(99)
		}
	case ab96reflect.Bool:
		if ab96Value.CanSet() {
			ab96Value.SetBool(false)
		}
	}
}

func TestAb96IndependentStorage(ab96Test *ab96testing.T) {
	ab96Expected, ab96Error := Ab96Entries()
	if ab96Error != nil {
		ab96Test.Fatal(ab96Error)
	}
	ab96Entries, ab96Error := Ab96Entries()
	if ab96Error != nil {
		ab96Test.Fatal(ab96Error)
	}
	for ab96Index := range ab96Entries {
		ab96Lookup, ab96Found, ab96Error := Ab96Lookup(ab96Expected[ab96Index].Ab96Key)
		if ab96Error != nil || !ab96Found {
			ab96Test.Fatal(ab96Error)
		}
		ab96Mutate(ab96reflect.ValueOf(&ab96Lookup).Elem())
		ab96Again, ab96Found, ab96Error := Ab96Lookup(ab96Expected[ab96Index].Ab96Key)
		if ab96Error != nil || !ab96Found || !ab96reflect.DeepEqual(ab96Again, ab96Expected[ab96Index]) {
			ab96Test.Fatal("lookup aliases nested storage")
		}
	}
	ab96Mutate(ab96reflect.ValueOf(&ab96Entries).Elem())
	ab96Again, ab96Error := Ab96Entries()
	if ab96Error != nil || !ab96reflect.DeepEqual(ab96Again, ab96Expected) {
		ab96Test.Fatal("entries alias nested storage")
	}
}

func TestAb96MalformedData(ab96Test *ab96testing.T) {
	ab96Saved := ab96Data
	ab96Test.Cleanup(func() { ab96Data = ab96Saved })
	ab96Key := ab96Find(ab96Test, "common", nil, "categories").Ab96Key
	for _, ab96Invalid := range []string{"{", "null", "[]", "[{}]", string(ab96Saved) + " []", "[{\"unexpected\":1}]"} {
		ab96Data = []byte(ab96Invalid)
		ab96Entries, ab96Error := Ab96Entries()
		if ab96Error == nil || ab96Entries != nil {
			ab96Test.Fatal("invalid entries succeeded")
		}
		ab96Entry, ab96Found, ab96Error := Ab96Lookup(ab96Key)
		if ab96Error == nil || ab96Found || !ab96reflect.DeepEqual(ab96Entry, Ab96Entry{}) {
			ab96Test.Fatal("invalid lookup succeeded")
		}
	}
	ab96Data = ab96Saved
	ab96Entries, ab96Error := Ab96Entries()
	if ab96Error != nil {
		ab96Test.Fatal(ab96Error)
	}
	ab96Data, ab96Error = ab96json.Marshal([]Ab96Entry{ab96Entries[0], ab96Entries[0]})
	if ab96Error != nil {
		ab96Test.Fatal(ab96Error)
	}
	if _, ab96Error = Ab96Entries(); ab96Error == nil {
		ab96Test.Fatal("duplicate accepted")
	}
}
