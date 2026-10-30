package vaillant

import (
	ab12JSON "encoding/json"
	ab12Fmt "fmt"
	ab12Reflect "reflect"
	ab12Testing "testing"
)

func ab12Pointer[ab12Value any](ab12Input ab12Value) *ab12Value { return &ab12Input }

func ab12Key(ab12Field string) Ab12Key {
	return Ab12Key{Ab12Family: "vaillant", Ab12Variant: "B50F", Ab12Primary: 181, Ab12Secondary: ab12Pointer(byte(15)), Ab12Operation: ab12Pointer(byte(2)), Ab12Direction: ab12Pointer("request"), Ab12Field: ab12Field}
}

func ab12Require(ab12T *ab12Testing.T, ab12Key Ab12Key) Ab12Entry {
	ab12T.Helper()
	ab12Entry, ab12Found, ab12Error := Ab12Lookup(ab12Key)
	if ab12Error != nil || !ab12Found {
		ab12T.Fatalf("lookup: found=%v error=%v", ab12Found, ab12Error)
	}
	return ab12Entry
}

func TestAb12FiniteContract(ab12T *ab12Testing.T) {
	ab12Entries, ab12Error := Ab12Entries()
	if ab12Error != nil || len(ab12Entries) != 8 {
		ab12T.Fatalf("entries: %v %d", ab12Error, len(ab12Entries))
	}
	ab12Masks := []byte{127, 127, 127, 255, 15, 7}
	for ab12Index, ab12Listed := range ab12Entries {
		ab12Field := ab12Fmt.Sprintf("status_container_%d", ab12Index+1)
		if ab12Index == 6 {
			ab12Field = "burner_status"
		}
		if ab12Index == 7 {
			ab12Field = "gas_pump_status"
		}
		ab12Entry := ab12Require(ab12T, ab12Key(ab12Field))
		if !ab12Reflect.DeepEqual(ab12Entry, ab12Listed) || !ab12Reflect.DeepEqual(ab12Entry.Ab12Key, ab12Key(ab12Field)) {
			ab12T.Fatal("entry identity")
		}
		if ab12Entry.Ab12Grade != "Hypothesis" || ab12Entry.Ab12Qualification != "Unknown" || ab12Entry.Ab12CoordinateGrade != "Unknown" || ab12Entry.Ab12Position != nil || ab12Entry.Ab12WidthBits == nil || *ab12Entry.Ab12WidthBits != 8 || ab12Entry.Ab12RawFallback != "Preserve the complete independent raw byte, including reserved bits; no physical status inference." {
			ab12T.Fatal("evidence boundary")
		}
		ab12Unknown := []string{"absolute byte position", "physical qualification"}
		if ab12Index < 4 {
			ab12Unknown = append(ab12Unknown, "individual allocated flag names")
		}
		if !ab12Reflect.DeepEqual(ab12Entry.Ab12Unknown, ab12Unknown) {
			ab12T.Fatal("unknown metadata")
		}
		if ab12Index >= 6 {
			if ab12Entry.Ab12AllocatedMask != nil || ab12Entry.Ab12ReservedMask != nil || len(ab12Entry.Ab12Bits) != 0 {
				ab12T.Fatal("raw byte assigned flags")
			}
			continue
		}
		if ab12Entry.Ab12AllocatedMask == nil || *ab12Entry.Ab12AllocatedMask != ab12Masks[ab12Index] || ab12Entry.Ab12ReservedMask == nil || *ab12Entry.Ab12ReservedMask != ^ab12Masks[ab12Index] || len(ab12Entry.Ab12Bits) != 8 {
			ab12T.Fatal("independent masks")
		}
		for ab12BitIndex, ab12Bit := range ab12Entry.Ab12Bits {
			var ab12Label, ab12Category *string
			if ab12Index == 2 && ab12BitIndex < 7 {
				ab12Category = ab12Pointer("temperature_limit")
			}
			if ab12Index == 3 {
				ab12Category = ab12Pointer("plausibility/channel")
			}
			if ab12Index == 4 && ab12BitIndex < 4 {
				ab12Label = ab12Pointer([]string{"valve_1_open", "valve_1_error", "valve_2_open", "valve_2_error"}[ab12BitIndex])
			}
			if ab12Index == 5 && ab12BitIndex < 3 {
				ab12Label = ab12Pointer([]string{"valve_open", "error", "air_check"}[ab12BitIndex])
			}
			ab12Expected := Ab12Bit{Ab12Bit: byte(ab12BitIndex), Ab12Label: ab12Label, Ab12Category: ab12Category, Ab12Reserved: ab12Masks[ab12Index]&(1<<ab12BitIndex) == 0}
			if !ab12Reflect.DeepEqual(ab12Bit, ab12Expected) {
				ab12T.Fatalf("container %d bit %d: %+v", ab12Index+1, ab12BitIndex, ab12Bit)
			}
		}
	}
}

func TestAb12SyntheticBits(ab12T *ab12Testing.T) {
	ab12Allocated := []byte{127, 127, 127, 255, 15, 7}
	ab12Reserved := []byte{128, 128, 128, 0, 240, 248}
	for ab12Container := 0; ab12Container < 6; ab12Container++ {
		for ab12Bit := 0; ab12Bit < 8; ab12Bit++ {
			ab12T.Run(ab12Fmt.Sprintf("ab12_container_%d_bit_%d", ab12Container+1, ab12Bit), func(ab12T *ab12Testing.T) {
				ab12Entry := ab12Require(ab12T, ab12Key(ab12Fmt.Sprintf("status_container_%d", ab12Container+1)))
				ab12Raw := byte(1 << ab12Bit)
				ab12Tuple := []byte{0, 0, 0, 0, 0, 0, 0, 0}
				ab12Tuple[ab12Container] = ab12Raw
				for _, ab12RawIndex := range []int{6, 7} {
					ab12Tuple[6], ab12Tuple[7] = 0, 0
					ab12Tuple[ab12RawIndex] = ^ab12Raw
					for ab12Index, ab12Byte := range ab12Tuple {
						ab12Field := ab12Fmt.Sprintf("status_container_%d", ab12Index+1)
						if ab12Index == 6 {
							ab12Field = "burner_status"
						}
						if ab12Index == 7 {
							ab12Field = "gas_pump_status"
						}
						ab12Metadata := ab12Require(ab12T, ab12Key(ab12Field))
						ab12Retained := ab12Byte
						if ab12Index < 6 {
							ab12AllocatedByte := ab12Byte & *ab12Metadata.Ab12AllocatedMask
							ab12ReservedByte := ab12Byte & *ab12Metadata.Ab12ReservedMask
							ab12ExpectedByte := byte(0)
							if ab12Index == ab12Container {
								ab12ExpectedByte = ab12Raw
							}
							if ab12AllocatedByte != ab12ExpectedByte&ab12Allocated[ab12Index] || ab12ReservedByte != ab12ExpectedByte&ab12Reserved[ab12Index] {
								ab12T.Fatal("independent tuple partitions")
							}
							ab12Retained = ab12AllocatedByte | ab12ReservedByte
						} else if ab12Metadata.Ab12Key.Ab12Field != ab12Field || ab12Metadata.Ab12AllocatedMask != nil || ab12Metadata.Ab12ReservedMask != nil {
							ab12T.Fatal("raw tuple identity")
						}
						ab12ExpectedRaw := byte(0)
						if ab12Index == ab12Container {
							ab12ExpectedRaw = ab12Raw
						}
						if ab12Index == ab12RawIndex {
							ab12ExpectedRaw = ^ab12Raw
						}
						if ab12Retained != ab12ExpectedRaw {
							ab12T.Fatal("independent raw tuple lost")
						}
					}
				}
				ab12AllocatedValue := ab12Raw & *ab12Entry.Ab12AllocatedMask
				ab12ReservedValue := ab12Raw & *ab12Entry.Ab12ReservedMask
				if ab12AllocatedValue != ab12Raw&ab12Allocated[ab12Container] || ab12ReservedValue != ab12Raw&ab12Reserved[ab12Container] || ab12AllocatedValue|ab12ReservedValue != ab12Raw {
					ab12T.Fatal("single bit lost or misclassified")
				}
				for _, ab12RawField := range []string{"burner_status", "gas_pump_status"} {
					ab12Other := ab12Require(ab12T, ab12Key(ab12RawField))
					if ab12Other.Ab12Key.Ab12Field != ab12RawField || ab12Other.Ab12AllocatedMask != nil || ab12Other.Ab12ReservedMask != nil || len(ab12Other.Ab12Bits) != 0 {
						ab12T.Fatal("raw byte aliased with a container")
					}
					*ab12Other.Ab12WidthBits = 0
					if *ab12Entry.Ab12WidthBits != 8 {
						ab12T.Fatal("raw byte shares storage")
					}
				}
			})
		}
		for _, ab12Raw := range []byte{128, 127, 255, 240, 248} {
			ab12Entry := ab12Require(ab12T, ab12Key(ab12Fmt.Sprintf("status_container_%d", ab12Container+1)))
			if ab12Raw&*ab12Entry.Ab12AllocatedMask != ab12Raw&ab12Allocated[ab12Container] || ab12Raw&*ab12Entry.Ab12ReservedMask != ab12Raw&ab12Reserved[ab12Container] {
				ab12T.Fatal("synthetic partition")
			}
		}
	}
}

func ab12Install(ab12T *ab12Testing.T, ab12Entries []Ab12Entry) {
	ab12T.Helper()
	ab12Before := ab12StatusJSON
	ab12Encoded, ab12Error := ab12JSON.Marshal(ab12Entries)
	if ab12Error != nil {
		ab12T.Fatal(ab12Error)
	}
	ab12StatusJSON = string(ab12Encoded)
	ab12T.Cleanup(func() { ab12StatusJSON = ab12Before })
}

func TestAb12EveryKeyAxis(ab12T *ab12Testing.T) {
	for _, ab12B524Case := range []bool{false, true} {
		ab12T.Run(ab12Fmt.Sprintf("ab12_b524_%v", ab12B524Case), func(ab12T *ab12Testing.T) {
			ab12Base := ab12Key("status_container_1")
			if ab12B524Case {
				ab12Base.Ab12Variant = "B524"
				ab12Base.Ab12Secondary = ab12Pointer(byte(36))
				ab12Base.Ab12Group = ab12Pointer(byte(1))
				ab12Base.Ab12RequestedInstance = ab12Pointer(byte(2))
				ab12Base.Ab12Register = ab12Pointer(uint16(3))
				ab12Base.Ab12CommandRole = ab12Pointer("synthetic_role")
				ab12Entries, ab12Error := Ab12Entries()
				if ab12Error != nil {
					ab12T.Fatal(ab12Error)
				}
				ab12Entries[0].Ab12Key = ab12Base
				ab12Install(ab12T, ab12Entries)
			}
			ab12Require(ab12T, ab12Base)
			ab12Type := ab12Reflect.TypeOf(ab12Base)
			for ab12Index := 0; ab12Index < ab12Type.NumField(); ab12Index++ {
				ab12Field := ab12Type.Field(ab12Index)
				ab12T.Run("ab12_"+ab12Field.Name, func(ab12T *ab12Testing.T) {
					ab12Source := ab12Reflect.ValueOf(ab12Base).Field(ab12Index)
					ab12Variants := []ab12Reflect.Value{}
					if ab12Source.Kind() == ab12Reflect.Pointer {
						ab12Zero := ab12Reflect.New(ab12Source.Type().Elem())
						ab12Other := ab12Reflect.New(ab12Source.Type().Elem())
						if ab12Other.Elem().Kind() == ab12Reflect.String {
							ab12Other.Elem().SetString("synthetic_other")
						} else {
							ab12Other.Elem().SetUint(99)
						}
						ab12Variants = append(ab12Variants, ab12Reflect.Zero(ab12Source.Type()), ab12Zero, ab12Other)
					} else {
						ab12Other := ab12Reflect.New(ab12Source.Type()).Elem()
						if ab12Source.Kind() == ab12Reflect.String {
							ab12Other.SetString("synthetic_other")
						} else {
							ab12Other.SetUint(99)
						}
						ab12Variants = append(ab12Variants, ab12Other)
					}
					for _, ab12Variant := range ab12Variants {
						if ab12Reflect.DeepEqual(ab12Source.Interface(), ab12Variant.Interface()) {
							continue
						}
						ab12Changed := ab12Base
						ab12Reflect.ValueOf(&ab12Changed).Elem().Field(ab12Index).Set(ab12Variant)
						ab12Entry, ab12Found, ab12Error := Ab12Lookup(ab12Changed)
						if ab12Found || ab12Error != nil || !ab12Reflect.DeepEqual(ab12Entry, Ab12Entry{}) {
							ab12T.Fatalf("axis ignored: %s", ab12Field.Name)
						}
					}
				})
			}
		})
	}
}

func TestAb12IndependentStorage(ab12T *ab12Testing.T) {
	ab12Before, ab12Error := Ab12Entries()
	if ab12Error != nil {
		ab12T.Fatal(ab12Error)
	}
	ab12Changed, ab12Error := Ab12Entries()
	if ab12Error != nil {
		ab12T.Fatal(ab12Error)
	}
	for ab12Index := range ab12Changed {
		ab12Entry := &ab12Changed[ab12Index]
		*ab12Entry.Ab12Key.Ab12Secondary = 0
		*ab12Entry.Ab12Key.Ab12Operation = 0
		*ab12Entry.Ab12Key.Ab12Direction = "changed"
		*ab12Entry.Ab12WidthBits = 0
		ab12Entry.Ab12Unknown[0] = "changed"
		if ab12Entry.Ab12AllocatedMask != nil {
			*ab12Entry.Ab12AllocatedMask = 0
			*ab12Entry.Ab12ReservedMask = 0
		}
		for ab12Bit := range ab12Entry.Ab12Bits {
			if ab12Entry.Ab12Bits[ab12Bit].Ab12Label != nil {
				*ab12Entry.Ab12Bits[ab12Bit].Ab12Label = "changed"
			}
			if ab12Entry.Ab12Bits[ab12Bit].Ab12Category != nil {
				*ab12Entry.Ab12Bits[ab12Bit].Ab12Category = "changed"
			}
			ab12Entry.Ab12Bits[ab12Bit].Ab12Reserved = true
		}
	}
	for _, ab12Expected := range ab12Before {
		ab12Lookup := ab12Require(ab12T, ab12Expected.Ab12Key)
		if !ab12Reflect.DeepEqual(ab12Lookup, ab12Expected) {
			ab12T.Fatal("accessor storage shared")
		}
		*ab12Lookup.Ab12Key.Ab12Secondary = 0
		*ab12Lookup.Ab12Key.Ab12Operation = 0
		*ab12Lookup.Ab12Key.Ab12Direction = "changed"
		*ab12Lookup.Ab12WidthBits = 0
		ab12Lookup.Ab12Unknown[0] = "changed"
		if ab12Lookup.Ab12AllocatedMask != nil {
			*ab12Lookup.Ab12AllocatedMask = 0
			*ab12Lookup.Ab12ReservedMask = 0
		}
		for ab12Bit := range ab12Lookup.Ab12Bits {
			if ab12Lookup.Ab12Bits[ab12Bit].Ab12Label != nil {
				*ab12Lookup.Ab12Bits[ab12Bit].Ab12Label = "changed"
			}
			if ab12Lookup.Ab12Bits[ab12Bit].Ab12Category != nil {
				*ab12Lookup.Ab12Bits[ab12Bit].Ab12Category = "changed"
			}
		}
		if len(ab12Lookup.Ab12Bits) > 0 {
			ab12Lookup.Ab12Bits[0].Ab12Label = ab12Pointer("changed")
		}
		if !ab12Reflect.DeepEqual(ab12Require(ab12T, ab12Expected.Ab12Key), ab12Expected) {
			ab12T.Fatal("lookup storage shared")
		}
	}
	ab12Again, ab12Error := Ab12Entries()
	if ab12Error != nil || !ab12Reflect.DeepEqual(ab12Again, ab12Before) {
		ab12T.Fatal("fresh entries mutated")
	}
}

func TestAb12MalformedData(ab12T *ab12Testing.T) {
	ab12Before := ab12StatusJSON
	ab12T.Cleanup(func() { ab12StatusJSON = ab12Before })
	ab12Cases := []string{"{", "null", "[]", ab12Before + " {}"}
	for _, ab12Mutation := range []func([]Ab12Entry){
		func(ab12Entries []Ab12Entry) { ab12Entries[0] = ab12Entries[1] },
		func(ab12Entries []Ab12Entry) { ab12Entries[0].Ab12Grade = "Proven" },
		func(ab12Entries []Ab12Entry) { ab12Entries[0].Ab12AllocatedMask = ab12Pointer(byte(255)) },
		func(ab12Entries []Ab12Entry) { ab12Entries[0].Ab12Bits[7].Ab12Label = ab12Pointer("demand") },
		func(ab12Entries []Ab12Entry) { ab12Entries[0].Ab12WidthBits = nil },
		func(ab12Entries []Ab12Entry) {
			ab12Entries[0].Ab12AllocatedMask = nil
			ab12Entries[0].Ab12ReservedMask = nil
			ab12Entries[0].Ab12Bits = nil
		},
		func(ab12Entries []Ab12Entry) { ab12Entries[0].Ab12Key.Ab12Primary = 0 },
	} {
		ab12StatusJSON = ab12Before
		ab12Entries, ab12Error := Ab12Entries()
		if ab12Error != nil {
			ab12T.Fatal(ab12Error)
		}
		ab12Mutation(ab12Entries)
		ab12Encoded, ab12Error := ab12JSON.Marshal(ab12Entries)
		if ab12Error != nil {
			ab12T.Fatal(ab12Error)
		}
		ab12Cases = append(ab12Cases, string(ab12Encoded))
	}
	for ab12Index, ab12Case := range ab12Cases {
		ab12T.Run(ab12Fmt.Sprintf("ab12_invalid_%d", ab12Index), func(ab12T *ab12Testing.T) {
			ab12StatusJSON = ab12Case
			if ab12Entries, ab12Error := Ab12Entries(); ab12Error == nil || ab12Entries != nil {
				ab12T.Fatal("invalid catalog returned entries")
			}
			if ab12Entry, ab12Found, ab12Error := Ab12Lookup(ab12Key("status_container_1")); ab12Error == nil || ab12Found || !ab12Reflect.DeepEqual(ab12Entry, Ab12Entry{}) {
				ab12T.Fatal("invalid catalog returned success")
			}
		})
	}
}

func TestAb12AbsentZeroAndB524Storage(ab12T *ab12Testing.T) {
	ab12Entries, ab12Error := Ab12Entries()
	if ab12Error != nil {
		ab12T.Fatal(ab12Error)
	}
	ab12Base := ab12Key("synthetic_identity")
	ab12Base.Ab12Variant = "B524"
	ab12Base.Ab12Secondary = ab12Pointer(byte(0))
	ab12Base.Ab12Operation = ab12Pointer(byte(0))
	ab12Base.Ab12Group = ab12Pointer(byte(0))
	ab12Base.Ab12RequestedInstance = ab12Pointer(byte(0))
	ab12Base.Ab12Register = ab12Pointer(uint16(0))
	ab12Base.Ab12CommandRole = ab12Pointer("")
	ab12Base.Ab12Direction = ab12Pointer("")
	ab12Entries[0].Ab12Key = ab12Base
	ab12Install(ab12T, ab12Entries)
	for _, ab12Name := range []string{"Ab12Secondary", "Ab12Operation", "Ab12Direction", "Ab12Group", "Ab12RequestedInstance", "Ab12Register", "Ab12CommandRole"} {
		ab12T.Run("ab12_"+ab12Name, func(ab12T *ab12Testing.T) {
			ab12Require(ab12T, ab12Base)
			ab12Absent := ab12Base
			ab12Field := ab12Reflect.ValueOf(&ab12Absent).Elem().FieldByName(ab12Name)
			ab12Field.Set(ab12Reflect.Zero(ab12Field.Type()))
			if _, ab12Found, ab12Error := Ab12Lookup(ab12Absent); ab12Found || ab12Error != nil {
				ab12T.Fatal("explicit zero matched absent")
			}
			ab12Lookup := ab12Require(ab12T, ab12Base)
			ab12PointerField := ab12Reflect.ValueOf(&ab12Lookup.Ab12Key).Elem().FieldByName(ab12Name).Elem()
			if ab12PointerField.Kind() == ab12Reflect.String {
				ab12PointerField.SetString("changed")
			} else {
				ab12PointerField.SetUint(9)
			}
			if !ab12Reflect.DeepEqual(ab12Require(ab12T, ab12Base), ab12Entries[0]) {
				ab12T.Fatal("nested key storage shared")
			}
		})
	}
	ab12Changed, ab12Error := Ab12Entries()
	if ab12Error != nil {
		ab12T.Fatal(ab12Error)
	}
	*ab12Changed[0].Ab12Key.Ab12Group = 7
	*ab12Changed[0].Ab12Key.Ab12RequestedInstance = 7
	*ab12Changed[0].Ab12Key.Ab12Register = 7
	*ab12Changed[0].Ab12Key.Ab12CommandRole = "changed"
	if !ab12Reflect.DeepEqual(ab12Require(ab12T, ab12Base), ab12Entries[0]) {
		ab12T.Fatal("entries B524 storage shared")
	}
}
