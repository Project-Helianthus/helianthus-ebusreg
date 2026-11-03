package vaillant

import (
	ab14Binary "encoding/binary"
	ab14JSON "encoding/json"
	ab14Reflect "reflect"
	ab14Testing "testing"
)

func ab14Ptr[ab14Type any](ab14Value ab14Type) *ab14Type { return &ab14Value }

func ab14Key(ab14Field string) Ab14Key {
	return Ab14Key{Ab14Family: "vaillant", Ab14Variant: "B509", Ab14Primary: 181, Ab14Secondary: ab14Ptr(uint8(9)), Ab14Operation: ab14Ptr(uint8(89)), Ab14Field: ab14Field}
}

func ab14Get(ab14T *ab14Testing.T, ab14Key Ab14Key) Ab14Entry {
	ab14T.Helper()
	ab14Entry, ab14OK, ab14Error := Ab14Lookup(ab14Key)
	if ab14Error != nil || !ab14OK {
		ab14T.Fatalf("lookup: %v %v", ab14OK, ab14Error)
	}
	return ab14Entry
}

func TestAb14Quantities(ab14T *ab14Testing.T) {
	ab14Entries, ab14Error := Ab14Entries()
	if ab14Error != nil || len(ab14Entries) != 5 {
		ab14T.Fatalf("entries: %v %v", ab14Entries, ab14Error)
	}
	for _, ab14Case := range []struct {
		ab14Field          string
		ab14Representation string
		ab14Divisor        float64
		ab14Unit           string
		ab14Fragments      [][2]byte
		ab14Expected       []float64
	}{
		{"temperature", "signedLE16", 16, "Celsius", [][2]byte{{0xf0, 0xff}, {0, 0x80}}, []float64{-1, -2048}},
		{"flow", "unsignedLE16", 100, "litres per minute", [][2]byte{{0x7b, 0}, {0xff, 0xff}}, []float64{1.23, 655.35}},
		{"volume", "unsignedLE16", 10, "litres", [][2]byte{{0x7b, 0}, {0xff, 0xff}}, []float64{12.3, 6553.5}},
	} {
		ab14T.Run("ab14_"+ab14Case.ab14Field, func(ab14T *ab14Testing.T) {
			ab14Entry := ab14Get(ab14T, ab14Key(ab14Case.ab14Field))
			ab14Want := &Ab14Codec{Ab14Representation: ab14Case.ab14Representation, Ab14Width: ab14Ptr(uint8(16)), Ab14Divisor: ab14Ptr(ab14Case.ab14Divisor), Ab14Unit: ab14Ptr(ab14Case.ab14Unit)}
			if !ab14Reflect.DeepEqual(ab14Entry.Ab14Codec, ab14Want) {
				ab14T.Fatalf("codec: %#v", ab14Entry.Ab14Codec)
			}
			for ab14Index, ab14Fragment := range ab14Case.ab14Fragments {
				ab14Raw := ab14Binary.LittleEndian.Uint16(ab14Fragment[:])
				ab14Value := float64(ab14Raw)
				if ab14Entry.Ab14Codec.Ab14Representation == "signedLE16" {
					ab14Value = float64(int16(ab14Raw))
				}
				if ab14Value / *ab14Entry.Ab14Codec.Ab14Divisor != ab14Case.ab14Expected[ab14Index] {
					ab14T.Fatal("independent fragment expectation failed")
				}
			}
		})
	}
	for _, ab14Field := range []string{"minimum_volume", "maximum_volume"} {
		ab14Entry := ab14Get(ab14T, ab14Key(ab14Field))
		if !ab14Reflect.DeepEqual(ab14Entry.Ab14Codec, &Ab14Codec{Ab14Representation: "raw"}) {
			ab14T.Fatal("raw limit acquired a scale, width or unit")
		}
		for _, ab14Raw := range []uint16{123, 456} {
			ab14Value := float64(ab14Raw)
			if ab14Entry.Ab14Codec.Ab14Divisor != nil {
				ab14Value /= *ab14Entry.Ab14Codec.Ab14Divisor
			}
			if ab14Value != float64(ab14Raw) {
				ab14T.Fatal("raw limit scaled")
			}
		}
	}
	for _, ab14Entry := range ab14Entries {
		ab14Unknown := []string{"absolute_field_positions", "native_support", "sentinels", "response_adaptation", "direction", "physical_qualification", "broad_catalog_completeness"}
		if ab14Entry.Ab14Codec.Ab14Representation == "raw" {
			ab14Unknown = append(ab14Unknown, "width", "scale", "coordinates", "enforcement", "unit")
		}
		if ab14Entry.Ab14Grade != "Hypothesis" || ab14Entry.Ab14MetadataPresence != "Proven" || ab14Entry.Ab14PhysicalQualification != "Unknown" || !ab14Entry.Ab14PreserveRaw || ab14Entry.Ab14Position != nil || ab14Entry.Ab14Label != nil || !ab14Reflect.DeepEqual(ab14Entry.Ab14Unknown, ab14Unknown) {
			ab14T.Fatal("evidence boundary changed")
		}
	}
}

func ab14Install(ab14T *ab14Testing.T, ab14Entries []Ab14Entry) {
	ab14T.Helper()
	ab14Bytes, ab14Error := ab14JSON.Marshal(ab14Entries)
	if ab14Error != nil {
		ab14T.Fatal(ab14Error)
	}
	ab14Original := ab14Data
	ab14T.Cleanup(func() { ab14Data = ab14Original })
	ab14Data = string(ab14Bytes)
}

func TestAb14ExactKeyAxes(ab14T *ab14Testing.T) {
	for _, ab14B524 := range []bool{false, true} {
		ab14T.Run(map[bool]string{false: "ab14_B509", true: "ab14_B524_synthetic"}[ab14B524], func(ab14T *ab14Testing.T) {
			ab14Base := ab14Key("temperature")
			if ab14B524 {
				ab14Base.Ab14Variant = "B524"
				ab14Base.Ab14Secondary = ab14Ptr(uint8(36))
				ab14Base.Ab14Group = ab14Ptr(uint8(1))
				ab14Base.Ab14RequestedInstance = ab14Ptr(uint8(2))
				ab14Base.Ab14Register = ab14Ptr(uint16(257))
				ab14Base.Ab14CommandRole = ab14Ptr("synthetic-role")
				ab14Entries, ab14Error := Ab14Entries()
				if ab14Error != nil {
					ab14T.Fatal(ab14Error)
				}
				ab14Entries[0].Ab14Key = ab14Base
				ab14Install(ab14T, ab14Entries)
			}
			ab14Expected := ab14Get(ab14T, ab14Base)
			for _, ab14Case := range []struct {
				ab14Name   string
				ab14Change func(*Ab14Key)
			}{
				{"family", func(ab14K *Ab14Key) { ab14K.Ab14Family = "synthetic-family" }},
				{"variant", func(ab14K *Ab14Key) { ab14K.Ab14Variant = "synthetic-variant" }},
				{"primary", func(ab14K *Ab14Key) { ab14K.Ab14Primary = 182 }},
				{"secondary", func(ab14K *Ab14Key) { ab14K.Ab14Secondary = ab14Ptr(uint8(10)) }},
				{"secondary_absent", func(ab14K *Ab14Key) { ab14K.Ab14Secondary = nil }},
				{"secondary_zero", func(ab14K *Ab14Key) { ab14K.Ab14Secondary = ab14Ptr(uint8(0)) }},
				{"operation", func(ab14K *Ab14Key) { ab14K.Ab14Operation = ab14Ptr(uint8(88)) }},
				{"operation_absent", func(ab14K *Ab14Key) { ab14K.Ab14Operation = nil }},
				{"operation_zero", func(ab14K *Ab14Key) { ab14K.Ab14Operation = ab14Ptr(uint8(0)) }},
				{"direction", func(ab14K *Ab14Key) { ab14K.Ab14Direction = ab14Ptr("response") }},
				{"field", func(ab14K *Ab14Key) { ab14K.Ab14Field = "synthetic-field" }},
				{"group", func(ab14K *Ab14Key) { ab14K.Ab14Group = ab14Ptr(uint8(3)) }},
				{"group_zero", func(ab14K *Ab14Key) { ab14K.Ab14Group = ab14Ptr(uint8(0)) }},
				{"requested_instance", func(ab14K *Ab14Key) { ab14K.Ab14RequestedInstance = ab14Ptr(uint8(3)) }},
				{"requested_instance_zero", func(ab14K *Ab14Key) { ab14K.Ab14RequestedInstance = ab14Ptr(uint8(0)) }},
				{"register", func(ab14K *Ab14Key) { ab14K.Ab14Register = ab14Ptr(uint16(513)) }},
				{"register_zero", func(ab14K *Ab14Key) { ab14K.Ab14Register = ab14Ptr(uint16(0)) }},
				{"command_role", func(ab14K *Ab14Key) { ab14K.Ab14CommandRole = ab14Ptr("different-role") }},
			} {
				ab14T.Run("ab14_"+ab14Case.ab14Name, func(ab14T *ab14Testing.T) {
					ab14Changed := ab14Base
					ab14Case.ab14Change(&ab14Changed)
					ab14Got, ab14OK, ab14Error := Ab14Lookup(ab14Changed)
					if ab14Error != nil || ab14OK || !ab14Reflect.DeepEqual(ab14Got, Ab14Entry{}) {
						ab14T.Fatalf("axis matched: %v %v", ab14OK, ab14Error)
					}
					if !ab14Reflect.DeepEqual(ab14Get(ab14T, ab14Base), ab14Expected) {
						ab14T.Fatal("baseline changed")
					}
				})
			}
		})
	}
}

func TestAb14AbsentZeroAndCopies(ab14T *ab14Testing.T) {
	ab14Entries, ab14Error := Ab14Entries()
	if ab14Error != nil {
		ab14T.Fatal(ab14Error)
	}
	ab14Baseline := ab14Get(ab14T, ab14Key("temperature"))
	*ab14Entries[0].Ab14Key.Ab14Secondary = 0
	*ab14Entries[0].Ab14Key.Ab14Operation = 0
	*ab14Entries[0].Ab14Codec.Ab14Width = 1
	*ab14Entries[0].Ab14Codec.Ab14Divisor = 1
	*ab14Entries[0].Ab14Codec.Ab14Unit = "changed"
	ab14Entries[0].Ab14Unknown[0] = "changed"
	ab14Lookup := ab14Get(ab14T, ab14Key("temperature"))
	if !ab14Reflect.DeepEqual(ab14Lookup, ab14Baseline) {
		ab14T.Fatal("entries storage shared")
	}
	*ab14Lookup.Ab14Codec.Ab14Divisor = 999
	*ab14Lookup.Ab14Key.Ab14Secondary = 0
	ab14Lookup.Ab14Unknown[0] = "changed"
	ab14Again, ab14Error := Ab14Entries()
	if ab14Error != nil || !ab14Reflect.DeepEqual(ab14Again[0], ab14Baseline) || !ab14Reflect.DeepEqual(ab14Get(ab14T, ab14Key("temperature")), ab14Baseline) {
		ab14T.Fatal("lookup storage shared")
	}
	for _, ab14Axis := range []string{"secondary", "operation", "group", "requested_instance", "register"} {
		ab14T.Run("ab14_"+ab14Axis, func(ab14T *ab14Testing.T) {
			ab14Rows, ab14Error := Ab14Entries()
			if ab14Error != nil {
				ab14T.Fatal(ab14Error)
			}
			ab14Nil := ab14Key("temperature")
			ab14Zero := ab14Nil
			switch ab14Axis {
			case "secondary":
				ab14Nil.Ab14Secondary = nil
				ab14Zero.Ab14Secondary = ab14Ptr(uint8(0))
			case "operation":
				ab14Nil.Ab14Operation = nil
				ab14Zero.Ab14Operation = ab14Ptr(uint8(0))
			case "group":
				ab14Zero.Ab14Group = ab14Ptr(uint8(0))
			case "requested_instance":
				ab14Zero.Ab14RequestedInstance = ab14Ptr(uint8(0))
			case "register":
				ab14Zero.Ab14Register = ab14Ptr(uint16(0))
			}
			ab14Rows[0].Ab14Key = ab14Nil
			ab14Rows[1].Ab14Key = ab14Zero
			ab14Install(ab14T, ab14Rows)
			if ab14Get(ab14T, ab14Nil).Ab14Codec.Ab14Representation != "signedLE16" || ab14Get(ab14T, ab14Zero).Ab14Codec.Ab14Representation != "unsignedLE16" {
				ab14T.Fatal("absent and zero collapsed")
			}
		})
	}
}

func TestAb14MalformedData(ab14T *ab14Testing.T) {
	ab14Original := ab14Data
	for _, ab14Bad := range []string{"{", "null", "[]", ab14Original + " []"} {
		ab14T.Run("ab14_invalid", func(ab14T *ab14Testing.T) {
			ab14Data = ab14Bad
			ab14T.Cleanup(func() { ab14Data = ab14Original })
			ab14Entries, ab14Error := Ab14Entries()
			if ab14Error == nil || ab14Entries != nil {
				ab14T.Fatal("malformed entries accepted")
			}
			ab14Entry, ab14OK, ab14Error := Ab14Lookup(ab14Key("temperature"))
			if ab14Error == nil || ab14OK || !ab14Reflect.DeepEqual(ab14Entry, Ab14Entry{}) {
				ab14T.Fatal("malformed lookup succeeded")
			}
		})
	}
	for _, ab14Mutation := range []func([]Ab14Entry){
		func(ab14Rows []Ab14Entry) { ab14Rows[1].Ab14Key = ab14Rows[0].Ab14Key },
		func(ab14Rows []Ab14Entry) { ab14Rows[0].Ab14Grade = "Proven" },
		func(ab14Rows []Ab14Entry) { ab14Rows[0].Ab14Codec = nil },
		func(ab14Rows []Ab14Entry) { ab14Rows[0].Ab14Key.Ab14Family = "" },
	} {
		ab14T.Run("ab14_invalid_entry", func(ab14T *ab14Testing.T) {
			ab14Rows, ab14Error := Ab14Entries()
			if ab14Error != nil {
				ab14T.Fatal(ab14Error)
			}
			ab14Mutation(ab14Rows)
			ab14Install(ab14T, ab14Rows)
			ab14Entry, ab14OK, ab14Error := Ab14Lookup(ab14Key("temperature"))
			if ab14Error == nil || ab14OK || !ab14Reflect.DeepEqual(ab14Entry, Ab14Entry{}) {
				ab14T.Fatal("invalid entry accepted")
			}
		})
	}
}
