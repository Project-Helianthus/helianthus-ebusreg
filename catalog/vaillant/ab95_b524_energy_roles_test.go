package vaillant

import (
	ab95json "encoding/json"
	ab95reflect "reflect"
	ab95testing "testing"
)

func ab95Pointer[ab95T any](ab95Value ab95T) *ab95T { return &ab95Value }

func ab95ExpectedKeys() []Ab95Key {
	return []Ab95Key{
		{"vaillant", "B524", 0xb5, ab95Pointer(uint8(0x24)), ab95Pointer(uint8(2)), "request", "setpoint", ab95Pointer(uint8(1)), ab95Pointer(uint8(0xff)), ab95Pointer(uint16(4)), "setpoint"},
		{"vaillant", "B524", 0xb5, ab95Pointer(uint8(0x24)), ab95Pointer(uint8(2)), "request", "hysteresis", ab95Pointer(uint8(0)), ab95Pointer(uint8(0xff)), ab95Pointer(uint16(0x27)), "hysteresis"},
		{"vaillant", "B524", 0xb5, ab95Pointer(uint8(0x24)), ab95Pointer(uint8(2)), "request", "control", ab95Pointer(uint8(1)), ab95Pointer(uint8(0xff)), ab95Pointer(uint16(0xd)), "control"},
		{"vaillant", "B524", 0xb5, ab95Pointer(uint8(0x24)), ab95Pointer(uint8(2)), "request", "offset", ab95Pointer(uint8(0)), ab95Pointer(uint8(0xff)), ab95Pointer(uint16(0xa2)), "offset"},
		{"vaillant", "B524", 0xb5, ab95Pointer(uint8(0x24)), ab95Pointer(uint8(2)), "request", "activation", ab95Pointer(uint8(0)), ab95Pointer(uint8(0xff)), ab95Pointer(uint16(0xb6)), "activation"},
	}
}

func TestAb95FiveTargets(ab95T *ab95testing.T) {
	ab95Entries, ab95Error := Ab95Entries()
	if ab95Error != nil || len(ab95Entries) != 5 {
		ab95T.Fatalf("entries: %v, %d", ab95Error, len(ab95Entries))
	}
	for ab95Index, ab95Key := range ab95ExpectedKeys() {
		ab95Entry, ab95OK, ab95Error := Ab95Lookup(ab95Key)
		if ab95Error != nil || !ab95OK || !ab95reflect.DeepEqual(ab95Entry.Ab95Key, ab95Key) {
			ab95T.Fatalf("target %d: %v %v", ab95Index, ab95OK, ab95Error)
		}
		if ab95Entry.Ab95ReadbackInstance == nil || *ab95Entry.Ab95ReadbackInstance != 0 {
			ab95T.Fatal("readback context must remain separate")
		}
		ab95Width := []*uint8{ab95Pointer(uint8(4)), ab95Pointer(uint8(4)), nil, ab95Pointer(uint8(4)), ab95Pointer(uint8(1))}[ab95Index]
		ab95Codec := []*string{ab95Pointer("binary32"), ab95Pointer("binary32"), nil, ab95Pointer("binary32"), nil}[ab95Index]
		ab95Unknown := []string{"binary32 byte order", "field offsets", "control enumeration", "native units", "response forms"}
		if ab95Codec == nil {
			ab95Unknown = append(ab95Unknown, "codec")
		}
		if ab95Width == nil {
			ab95Unknown = append(ab95Unknown, "width")
		}
		if !ab95reflect.DeepEqual(ab95Entry.Ab95WidthBytes, ab95Width) || !ab95reflect.DeepEqual(ab95Entry.Ab95Codec, ab95Codec) || !ab95reflect.DeepEqual(ab95Entry.Ab95Unknown, ab95Unknown) {
			ab95T.Fatalf("target %d: codec/width/unknown changed", ab95Index)
		}
		if ab95Entry.Ab95Grade != "Hypothesis" || ab95Entry.Ab95MetadataBoundaryGrade != "Proven" || ab95Entry.Ab95PhysicalQualification != "Unknown" || ab95Entry.Ab95MetadataBoundary != "Metadata lookup cannot establish applied state or installed capability." {
			ab95T.Fatal("grade boundary changed")
		}
		if ab95Entry.Ab95FieldOffset != nil || ab95Entry.Ab95ByteOrder != nil || ab95Entry.Ab95NativeUnit != nil || ab95Entry.Ab95ControlEnumeration != nil || ab95Entry.Ab95ResponseForm != nil {
			ab95T.Fatal("unprovided metadata invented")
		}
		if ab95Entry.Ab95RawPreservation != "Preserve native raw values and partial last-known-good fields independently." || ab95Entry.Ab95NativeContract != "https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/protocols/vaillant/ebus-vaillant-B524.md" {
			ab95T.Fatal("native contract/raw boundary changed")
		}
	}
}

func TestAb95ExactKeyAxes(ab95T *ab95testing.T) {
	ab95Changes := map[string]func(*Ab95Key){
		"ab95_family":                    func(ab95Key *Ab95Key) { ab95Key.Ab95Family = "other" },
		"ab95_variant":                   func(ab95Key *Ab95Key) { ab95Key.Ab95Variant = "other" },
		"ab95_primary":                   func(ab95Key *Ab95Key) { ab95Key.Ab95Primary = 0 },
		"ab95_secondary":                 func(ab95Key *Ab95Key) { ab95Key.Ab95Secondary = ab95Pointer(uint8(0)) },
		"ab95_secondary_absent":          func(ab95Key *Ab95Key) { ab95Key.Ab95Secondary = nil },
		"ab95_operation":                 func(ab95Key *Ab95Key) { ab95Key.Ab95Operation = ab95Pointer(uint8(6)) },
		"ab95_operation_zero":            func(ab95Key *Ab95Key) { ab95Key.Ab95Operation = ab95Pointer(uint8(0)) },
		"ab95_operation_absent":          func(ab95Key *Ab95Key) { ab95Key.Ab95Operation = nil },
		"ab95_direction":                 func(ab95Key *Ab95Key) { ab95Key.Ab95Direction = "response" },
		"ab95_field":                     func(ab95Key *Ab95Key) { ab95Key.Ab95Field = "other" },
		"ab95_group":                     func(ab95Key *Ab95Key) { ab95Key.Ab95Group = ab95Pointer(*ab95Key.Ab95Group ^ 1) },
		"ab95_group_absent":              func(ab95Key *Ab95Key) { ab95Key.Ab95Group = nil },
		"ab95_requested_instance":        func(ab95Key *Ab95Key) { ab95Key.Ab95RequestedInstance = ab95Pointer(uint8(0)) },
		"ab95_requested_instance_absent": func(ab95Key *Ab95Key) { ab95Key.Ab95RequestedInstance = nil },
		"ab95_register":                  func(ab95Key *Ab95Key) { ab95Key.Ab95Register = ab95Pointer(*ab95Key.Ab95Register + 1) },
		"ab95_register_zero":             func(ab95Key *Ab95Key) { ab95Key.Ab95Register = ab95Pointer(uint16(0)) },
		"ab95_register_absent":           func(ab95Key *Ab95Key) { ab95Key.Ab95Register = nil },
		"ab95_command_role":              func(ab95Key *Ab95Key) { ab95Key.Ab95CommandRole = "other" },
	}
	for ab95Name, ab95Change := range ab95Changes {
		ab95T.Run(ab95Name, func(ab95T *ab95testing.T) {
			for _, ab95Key := range ab95ExpectedKeys() {
				if _, ab95OK, ab95Error := Ab95Lookup(ab95Key); !ab95OK || ab95Error != nil {
					ab95T.Fatal("invalid positive control")
				}
				ab95Change(&ab95Key)
				if ab95Entry, ab95OK, ab95Error := Ab95Lookup(ab95Key); ab95OK || ab95Error != nil || !ab95reflect.DeepEqual(ab95Entry, Ab95Entry{}) {
					ab95T.Fatalf("axis inherited a target: %v %v", ab95OK, ab95Error)
				}
			}
		})
	}
}

func ab95Install(ab95T *ab95testing.T, ab95Entries []Ab95Entry) {
	ab95T.Helper()
	ab95Encoded, ab95Error := ab95json.Marshal(ab95Entries)
	if ab95Error != nil {
		ab95T.Fatal(ab95Error)
	}
	ab95Data = ab95Encoded
}

func TestAb95AbsentVersusZero(ab95T *ab95testing.T) {
	ab95Original := ab95Data
	defer func() { ab95Data = ab95Original }()
	ab95Changes := map[string]func(*Ab95Key, bool){
		"ab95_secondary": func(ab95Key *Ab95Key, ab95Present bool) {
			ab95Key.Ab95Secondary = nil
			if ab95Present {
				ab95Key.Ab95Secondary = ab95Pointer(uint8(0))
			}
		},
		"ab95_operation": func(ab95Key *Ab95Key, ab95Present bool) {
			ab95Key.Ab95Operation = nil
			if ab95Present {
				ab95Key.Ab95Operation = ab95Pointer(uint8(0))
			}
		},
		"ab95_group": func(ab95Key *Ab95Key, ab95Present bool) {
			ab95Key.Ab95Group = nil
			if ab95Present {
				ab95Key.Ab95Group = ab95Pointer(uint8(0))
			}
		},
		"ab95_requested_instance": func(ab95Key *Ab95Key, ab95Present bool) {
			ab95Key.Ab95RequestedInstance = nil
			if ab95Present {
				ab95Key.Ab95RequestedInstance = ab95Pointer(uint8(0))
			}
		},
		"ab95_register": func(ab95Key *Ab95Key, ab95Present bool) {
			ab95Key.Ab95Register = nil
			if ab95Present {
				ab95Key.Ab95Register = ab95Pointer(uint16(0))
			}
		},
	}
	for ab95Name, ab95Change := range ab95Changes {
		ab95T.Run(ab95Name, func(ab95T *ab95testing.T) {
			for _, ab95Present := range []bool{false, true} {
				ab95Data = ab95Original
				ab95Entries, ab95Error := Ab95Entries()
				if ab95Error != nil {
					ab95T.Fatal(ab95Error)
				}
				ab95Change(&ab95Entries[0].Ab95Key, ab95Present)
				ab95Install(ab95T, ab95Entries)
				ab95Key := ab95Entries[0].Ab95Key
				if _, ab95OK, ab95Error := Ab95Lookup(ab95Key); !ab95OK || ab95Error != nil {
					ab95T.Fatal("synthetic positive control")
				}
				ab95Change(&ab95Key, !ab95Present)
				if _, ab95OK, ab95Error := Ab95Lookup(ab95Key); ab95OK || ab95Error != nil {
					ab95T.Fatal("absent/zero collapsed")
				}
			}
		})
	}
}

func ab95Mutate(ab95Value ab95reflect.Value) {
	switch ab95Value.Kind() {
	case ab95reflect.Pointer:
		if !ab95Value.IsNil() {
			ab95Mutate(ab95Value.Elem())
		}
	case ab95reflect.Struct:
		for ab95Index := 0; ab95Index < ab95Value.NumField(); ab95Index++ {
			ab95Mutate(ab95Value.Field(ab95Index))
		}
	case ab95reflect.Slice:
		for ab95Index := 0; ab95Index < ab95Value.Len(); ab95Index++ {
			ab95Mutate(ab95Value.Index(ab95Index))
		}
	case ab95reflect.String:
		ab95Value.SetString("ab95_mutated")
	case ab95reflect.Uint8, ab95reflect.Uint16:
		ab95Value.SetUint(99)
	}
}

func TestAb95IndependentStorage(ab95T *ab95testing.T) {
	ab95Before, ab95Error := Ab95Entries()
	if ab95Error != nil {
		ab95T.Fatal(ab95Error)
	}
	ab95Returned, ab95Error := Ab95Entries()
	if ab95Error != nil {
		ab95T.Fatal(ab95Error)
	}
	ab95Mutate(ab95reflect.ValueOf(&ab95Returned).Elem())
	ab95Fresh, ab95Error := Ab95Entries()
	if ab95Error != nil || !ab95reflect.DeepEqual(ab95Fresh, ab95Before) {
		ab95T.Fatal("shared accessor storage")
	}
	for ab95Index, ab95Key := range ab95ExpectedKeys() {
		ab95Entry, ab95OK, ab95Error := Ab95Lookup(ab95Key)
		if !ab95OK || ab95Error != nil {
			ab95T.Fatal("lookup positive control")
		}
		ab95Mutate(ab95reflect.ValueOf(&ab95Entry).Elem())
		ab95FreshEntry, ab95OK, ab95Error := Ab95Lookup(ab95Key)
		if !ab95OK || ab95Error != nil || !ab95reflect.DeepEqual(ab95FreshEntry, ab95Before[ab95Index]) {
			ab95T.Fatal("shared lookup storage")
		}
	}
}

func TestAb95MalformedData(ab95T *ab95testing.T) {
	ab95Original := ab95Data
	defer func() { ab95Data = ab95Original }()
	ab95Cases := map[string][]byte{"ab95_syntax": []byte("["), "ab95_empty": []byte("[]"), "ab95_null": []byte("null"), "ab95_trailing": append(append([]byte(nil), ab95Original...), []byte(" null")...)}
	for ab95Name, ab95Change := range map[string]func([]Ab95Entry){
		"ab95_duplicate":      func(ab95Entries []Ab95Entry) { ab95Entries[1].Ab95Key = ab95Entries[0].Ab95Key },
		"ab95_missing_family": func(ab95Entries []Ab95Entry) { ab95Entries[4].Ab95Key.Ab95Family = "" },
		"ab95_grade":          func(ab95Entries []Ab95Entry) { ab95Entries[4].Ab95Grade = "Proven" },
	} {
		ab95Data = ab95Original
		ab95Entries, ab95Error := Ab95Entries()
		if ab95Error != nil {
			ab95T.Fatal(ab95Error)
		}
		ab95Change(ab95Entries)
		ab95Install(ab95T, ab95Entries)
		ab95Cases[ab95Name] = ab95Data
	}
	for ab95Name, ab95Bytes := range ab95Cases {
		ab95T.Run(ab95Name, func(ab95T *ab95testing.T) {
			ab95Data = ab95Bytes
			if ab95Entries, ab95Error := Ab95Entries(); ab95Error == nil || ab95Entries != nil {
				ab95T.Fatal("invalid data returned entries")
			}
			if ab95Entry, ab95OK, ab95Error := Ab95Lookup(ab95ExpectedKeys()[0]); ab95Error == nil || ab95OK || !ab95reflect.DeepEqual(ab95Entry, Ab95Entry{}) {
				ab95T.Fatal("invalid data returned successful lookup")
			}
		})
	}
}
