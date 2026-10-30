package vaillant

import (
	ab85bytes "bytes"
	ab85context "context"
	ab85binary "encoding/binary"
	ab85json "encoding/json"
	ab85errors "errors"
	ab85math "math"
	ab85reflect "reflect"
	ab85testing "testing"

	ab85standard "github.com/Project-Helianthus/helianthus-ebusreg/catalog/ebus_standard"
)

func ab85Pointer[ab85Value any](ab85Input ab85Value) *ab85Value { return &ab85Input }

func ab85ExpectedKey(ab85Family *string, ab85Direction string, ab85Group uint8, ab85Register uint16, ab85Field string) Ab85Key {
	return Ab85Key{Ab85Family: ab85Family, Ab85Primary: ab85Pointer(uint8(0xB5)), Ab85Secondary: ab85Pointer(uint8(0x24)), Ab85Direction: ab85Pointer(ab85Direction), Ab85Field: ab85Pointer(ab85Field), Ab85Opcode: ab85Pointer(uint8(2)), Ab85Group: ab85Pointer(ab85Group), Ab85Register: ab85Pointer(ab85Register)}
}

func ab85ExpectedEntries() []Ab85ScalarShape {
	var ab85Entries []Ab85ScalarShape
	for _, ab85Case := range []struct {
		ab85Family           *string
		ab85Direction        string
		ab85Group            uint8
		ab85Register         uint16
		ab85Field, ab85Codec string
		ab85Width            int
	}{
		{nil, "read", 3, 0x001B, "scalar", "u8", 1},
		{nil, "read", 3, 0x001B, "scalar", "u16", 2},
		{nil, "read", 1, 0x0014, "scalar", "u8", 1},
		{nil, "read", 1, 0x0019, "scalar", "u8", 1},
		{nil, "read", 1, 0x0015, "scalar", "binary32", 4},
		{ab85Pointer("controller-family"), "write", 3, 0x0006, "mode", "u16", 2},
		{ab85Pointer("controller-family"), "write", 3, 0x0022, "manual-temperature", "binary32", 4},
		{ab85Pointer("controller-family"), "write", 3, 0x0007, "comfort-temperature", "binary32", 4},
	} {
		ab85Signedness, ab85Order := "unsigned", "little-endian"
		if ab85Case.ab85Codec == "binary32" {
			ab85Signedness = "signed"
		}
		if ab85Case.ab85Codec == "u8" {
			ab85Order = "not applicable"
		}
		ab85Entries = append(ab85Entries, Ab85ScalarShape{
			Ab85Key:   ab85ExpectedKey(ab85Case.ab85Family, ab85Case.ab85Direction, ab85Case.ab85Group, ab85Case.ab85Register, ab85Case.ab85Field),
			Ab85Shape: ab85Case.ab85Codec, Ab85WidthBytes: ab85Pointer(ab85Case.ab85Width),
			Ab85Signedness: ab85Signedness, Ab85ByteOrder: ab85Order,
			Ab85Conversion: &Ab85Conversion{Ab85Encoding: ab85Case.ab85Codec, Ab85Factor: 1, Ab85Offset: 0},
			Ab85Grade:      "Hypothesis", Ab85RawFallback: "preserve native field bytes on every unsuccessful typed interpretation",
		})
	}
	return ab85Entries
}

func ab85Catalog(ab85T *ab85testing.T) Ab85Catalog {
	ab85T.Helper()
	ab85Catalog, ab85Error := Ab85ScalarShapes()
	if ab85Error != nil {
		ab85T.Fatal(ab85Error)
	}
	return ab85Catalog
}

func TestAb85ProductionMetadata(ab85T *ab85testing.T) {
	var ab85Native struct {
		Ab85Entries []map[string]ab85json.RawMessage `json:"entries"`
	}
	if ab85Error := ab85json.Unmarshal(ab85ScalarShapesJSON, &ab85Native); ab85Error != nil {
		ab85T.Fatal(ab85Error)
	}
	for _, ab85Entry := range ab85Native.Ab85Entries {
		for _, ab85Label := range []string{"key", "shape", "logical_position", "width_bytes", "signedness", "byte_order", "conversion", "unit", "default", "bounds", "enumeration", "replacement_values", "reserved_values", "grade", "raw_fallback"} {
			if _, ab85Present := ab85Entry[ab85Label]; !ab85Present {
				ab85T.Fatalf("ab85 missing explicit field %s", ab85Label)
			}
		}
		var ab85Key map[string]ab85json.RawMessage
		if ab85Error := ab85json.Unmarshal(ab85Entry["key"], &ab85Key); ab85Error != nil {
			ab85T.Fatal(ab85Error)
		}
		for _, ab85Label := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "opcode", "group", "instance", "register"} {
			if _, ab85Present := ab85Key[ab85Label]; !ab85Present {
				ab85T.Fatalf("ab85 missing explicit identity axis %s", ab85Label)
			}
		}
	}
	ab85Got := ab85Catalog(ab85T)
	ab85Want := Ab85Catalog{
		Ab85Catalog: "ab85_b524_scalar_shapes", Ab85RepresentationGrade: "Proven", Ab85PhysicalQualification: "Unknown", Ab85Complete: false,
		Ab85Limitations: []string{
			"Finite supplemental declarations; broader family coverage is incomplete.",
			"Family and instance absence is exact missing qualification, never a wildcard.",
			"Request length and reply length do not choose a semantic alternative or family qualification.",
			"Decode failure retains a prior cached value with uncertain freshness; no cache updater is supplied.",
			"Installed support, omitted instance predicates, complete enumerations, units, persistence and wire adaptation are Unknown.",
			"No builder, executor, provider registration, invocation, admission or runtime cache is supplied.",
			"Logical positions including prefixes are Unknown; widths describe isolated candidate field bytes only.",
		}, Ab85Entries: ab85ExpectedEntries(),
	}
	if !ab85reflect.DeepEqual(ab85Got, ab85Want) {
		ab85T.Fatalf("ab85 production metadata differs from independent contract: %+v", ab85Got)
	}
}

func TestAb85ExactIdentity(ab85T *ab85testing.T) {
	for ab85Index, ab85Expected := range ab85ExpectedEntries() {
		ab85Control, ab85Error := Ab85Lookup(ab85Expected.Ab85Key)
		if ab85Error != nil {
			ab85T.Fatal(ab85Error)
		}
		ab85Count := 1
		if ab85Index < 2 {
			ab85Count = 2
		}
		if len(ab85Control) != ab85Count {
			ab85T.Fatalf("ab85 valid control %d: %v", ab85Index, ab85Control)
		}
		for ab85Axis := 0; ab85Axis < ab85reflect.TypeOf(Ab85Key{}).NumField(); ab85Axis++ {
			for ab85Mutation := 0; ab85Mutation < 3; ab85Mutation++ {
				ab85Key := ab85Expected.Ab85Key
				ab85Slot := ab85reflect.ValueOf(&ab85Key).Elem().Field(ab85Axis)
				ab85Old := ab85reflect.ValueOf(ab85Expected.Ab85Key).Field(ab85Axis)
				if ab85Mutation == 0 {
					ab85Slot.SetZero()
				} else {
					ab85Value := ab85reflect.New(ab85Slot.Type().Elem())
					if ab85Value.Elem().Kind() == ab85reflect.String {
						if ab85Mutation == 2 {
							ab85Value.Elem().SetString("ab85_other")
						}
					} else if ab85Mutation == 2 {
						ab85Number := uint64(6)
						if !ab85Old.IsNil() {
							ab85Number = ab85Old.Elem().Uint() ^ 0x80
						}
						if ab85reflect.TypeOf(ab85Key).Field(ab85Axis).Name == "Ab85Opcode" {
							ab85Number = 6
						}
						if ab85reflect.TypeOf(ab85Key).Field(ab85Axis).Name == "Ab85Register" {
							ab85Number = ab85Old.Elem().Uint() ^ 0x100
						}
						ab85Value.Elem().SetUint(ab85Number)
					}
					ab85Slot.Set(ab85Value)
				}
				if ab85reflect.DeepEqual(ab85Key, ab85Expected.Ab85Key) {
					continue
				}
				ab85Matches, ab85Error := Ab85Lookup(ab85Key)
				if ab85Error != nil || len(ab85Matches) != 0 {
					ab85T.Fatalf("ab85 identity axis %s mutation %d inherited entry: %v %v", ab85reflect.TypeOf(ab85Key).Field(ab85Axis).Name, ab85Mutation, ab85Matches, ab85Error)
				}
			}
		}
		ab85WrongDirection := ab85Expected.Ab85Key
		if *ab85WrongDirection.Ab85Direction == "read" {
			ab85WrongDirection.Ab85Direction = ab85Pointer("write")
		} else {
			ab85WrongDirection.Ab85Direction = ab85Pointer("read")
		}
		ab85Matches, ab85Error := Ab85Lookup(ab85WrongDirection)
		if ab85Error != nil || len(ab85Matches) != 0 {
			ab85T.Fatal("ab85 read/write leaked", ab85Error)
		}
	}
}

func ab85Evaluate(ab85Shape Ab85ScalarShape, ab85Raw []byte) (float64, bool, []byte) {
	ab85Evidence := ab85bytes.Clone(ab85Raw)
	if ab85Shape.Ab85WidthBytes == nil || len(ab85Raw) != *ab85Shape.Ab85WidthBytes || ab85Shape.Ab85Conversion == nil {
		return 0, false, ab85Evidence
	}
	var ab85Value float64
	switch ab85Shape.Ab85Conversion.Ab85Encoding {
	case "u8":
		if len(ab85Raw) != 1 || ab85Shape.Ab85Signedness != "unsigned" || ab85Shape.Ab85ByteOrder != "not applicable" {
			return 0, false, ab85Evidence
		}
		ab85Value = float64(ab85Raw[0])
	case "u16":
		if len(ab85Raw) != 2 || ab85Shape.Ab85Signedness != "unsigned" || ab85Shape.Ab85ByteOrder != "little-endian" {
			return 0, false, ab85Evidence
		}
		ab85Value = float64(ab85binary.LittleEndian.Uint16(ab85Raw))
	case "binary32":
		if len(ab85Raw) != 4 || ab85Shape.Ab85Signedness != "signed" || ab85Shape.Ab85ByteOrder != "little-endian" {
			return 0, false, ab85Evidence
		}
		ab85Value = float64(ab85math.Float32frombits(ab85binary.LittleEndian.Uint32(ab85Raw)))
	default:
		return 0, false, ab85Evidence
	}
	return ab85Value*ab85Shape.Ab85Conversion.Ab85Factor + ab85Shape.Ab85Conversion.Ab85Offset, true, ab85Evidence
}

func TestAb85SyntheticFields(ab85T *ab85testing.T) {
	ab85Entries := ab85Catalog(ab85T).Ab85Entries
	for _, ab85Case := range []struct {
		ab85Index int
		ab85Raw   []byte
		ab85Want  float64
		ab85OK    bool
	}{
		{0, []byte{1}, 1, true}, {1, []byte{1}, 0, false},
		{0, []byte{1, 0}, 0, false}, {1, []byte{1, 0}, 1, true},
		{0, []byte{0, 1}, 0, false}, {1, []byte{0, 1}, 256, true},
		{2, []byte{14}, 14, true}, {3, []byte{19}, 19, true},
		{4, []byte{0, 0, 128, 63}, 1, true},
		{5, []byte{1, 1}, 257, true},
		{6, []byte{0, 0, 128, 63}, 1, true}, {7, []byte{0, 0, 128, 63}, 1, true},
	} {
		ab85Value, ab85OK, ab85Raw := ab85Evaluate(ab85Entries[ab85Case.ab85Index], ab85Case.ab85Raw)
		if ab85Value != ab85Case.ab85Want || ab85OK != ab85Case.ab85OK || !ab85bytes.Equal(ab85Raw, ab85Case.ab85Raw) {
			ab85T.Fatalf("ab85 synthetic case %+v got %v %v %x", ab85Case, ab85Value, ab85OK, ab85Raw)
		}
	}
	for _, ab85Entry := range ab85Entries {
		for _, ab85Length := range []int{0, *ab85Entry.Ab85WidthBytes - 1, *ab85Entry.Ab85WidthBytes + 1, 6} {
			if ab85Length == *ab85Entry.Ab85WidthBytes {
				continue
			}
			ab85Input := ab85bytes.Repeat([]byte{1}, ab85Length)
			_, ab85OK, ab85Raw := ab85Evaluate(ab85Entry, ab85Input)
			if ab85OK || !ab85bytes.Equal(ab85Raw, ab85Input) {
				ab85T.Fatal("ab85 short/extra bytes gained typed claim or lost raw")
			}
		}
		ab85Matches, ab85Error := Ab85Lookup(ab85Entry.Ab85Key)
		if ab85Error != nil || len(ab85Matches) == 0 {
			ab85T.Fatal("ab85 exact identity lost", ab85Error)
		}
		for _, ab85Match := range ab85Matches {
			if ab85Match.Ab85Grade != "Hypothesis" {
				ab85T.Fatal("ab85 length promoted candidate grade")
			}
		}
	}
	if ab85reflect.DeepEqual(ab85Entries[2].Ab85Key, ab85Entries[3].Ab85Key) || ab85reflect.DeepEqual(ab85Entries[6].Ab85Key, ab85Entries[7].Ab85Key) {
		ab85T.Fatal("ab85 separate field roles merged")
	}
}

func ab85MutateShape(ab85Shape *Ab85ScalarShape) {
	*ab85Shape.Ab85Key.Ab85Primary = 0
	*ab85Shape.Ab85Key.Ab85Secondary = 0
	*ab85Shape.Ab85Key.Ab85Direction = "ab85_mutated"
	*ab85Shape.Ab85Key.Ab85Opcode = 6
	*ab85Shape.Ab85Key.Ab85Group = 0
	*ab85Shape.Ab85Key.Ab85Register = 0
	*ab85Shape.Ab85Key.Ab85Field = "ab85_mutated"
	if ab85Shape.Ab85Key.Ab85Family != nil {
		*ab85Shape.Ab85Key.Ab85Family = "ab85_mutated"
	}
	*ab85Shape.Ab85WidthBytes = 99
	ab85Shape.Ab85Conversion.Ab85Offset = 273.15
	ab85Shape.Ab85Enumeration = map[string]string{"ab85_mutated": "ab85_mutated"}
	ab85Shape.Ab85ReplacementValues = []string{"ab85_mutated"}
	ab85Shape.Ab85ReservedValues = []string{"ab85_mutated"}
	ab85Shape.Ab85Bounds = &Ab85Bounds{Ab85Minimum: ab85Pointer(99.0)}
}

func TestAb85DetachedCopies(ab85T *ab85testing.T) {
	ab85Before := ab85Catalog(ab85T)
	ab85Copy := ab85Catalog(ab85T)
	ab85Copy.Ab85Limitations[0] = "ab85_mutated"
	for ab85Index := range ab85Copy.Ab85Entries {
		ab85MutateShape(&ab85Copy.Ab85Entries[ab85Index])
	}
	if !ab85reflect.DeepEqual(ab85Catalog(ab85T), ab85Before) {
		ab85T.Fatal("ab85 accessor aliases nested data")
	}
	for _, ab85Expected := range ab85ExpectedEntries() {
		ab85First, ab85Error := Ab85Lookup(ab85Expected.Ab85Key)
		if ab85Error != nil {
			ab85T.Fatal(ab85Error)
		}
		ab85Second, ab85Error := Ab85Lookup(ab85Expected.Ab85Key)
		if ab85Error != nil {
			ab85T.Fatal(ab85Error)
		}
		for ab85Index := range ab85First {
			ab85MutateShape(&ab85First[ab85Index])
		}
		ab85Again, ab85Error := Ab85Lookup(ab85Expected.Ab85Key)
		if ab85Error != nil || !ab85reflect.DeepEqual(ab85Again, ab85Second) {
			ab85T.Fatal("ab85 lookup aliases nested data", ab85Error)
		}
	}
}

func TestAb85MetadataDoesNotEnableInvocation(ab85T *ab85testing.T) {
	ab85Provider, ab85Error := ab85standard.NewProvider(ab85standard.MustEmbeddedCatalog(), true)
	if ab85Error != nil {
		ab85T.Fatal(ab85Error)
	}
	for ab85Pass := 0; ab85Pass < 2; ab85Pass++ {
		if ab85Pass == 1 {
			_ = ab85Catalog(ab85T)
		}
		ab85Result, ab85Error := ab85Provider.Invoke(ab85context.Background(), "ebus_standard.service_data.start_counts", nil, ab85standard.CallerContextUserFacing)
		if ab85Error != nil || ab85Result["executed"] != false {
			ab85T.Fatal("ab85 existing stub changed", ab85Result, ab85Error)
		}
		_, ab85Error = ab85Provider.Invoke(ab85context.Background(), "ab85_b524_scalar_shapes", nil, ab85standard.CallerContextUserFacing)
		if !ab85errors.Is(ab85Error, ab85standard.ErrUnknownMethod) {
			ab85T.Fatal("ab85 metadata gained method registration", ab85Error)
		}
	}
}
