package vaillant

import (
	"encoding/binary"
	"encoding/json"
	"reflect"
	"testing"
)

func ab94Pointer[ab94Type any](ab94Value ab94Type) *ab94Type { return &ab94Value }

func ab94ExpectedField(ab94Name string, ab94Width int, ab94Order, ab94Conversion string, ab94RawMax uint64, ab94NumericMax float64) Ab94NumericField {
	return Ab94NumericField{
		Ab94Key:       Ab94Key{Ab94Family: "vaillant", Ab94Variant: "generic", Ab94Primary: 0xB555, Ab94Field: ab94Name},
		Ab94WidthBits: ab94Pointer(ab94Width), Ab94Signedness: "unsigned", Ab94ByteOrder: ab94Order,
		Ab94Conversion: ab94Conversion, Ab94RawMin: ab94Pointer(uint64(0)), Ab94RawMax: ab94Pointer(ab94RawMax),
		Ab94NumericMin: ab94Pointer(float64(0)), Ab94NumericMax: ab94Pointer(ab94NumericMax),
		Ab94Enumeration: map[uint64]string{}, Ab94Replacement: []Ab94Replacement{}, Ab94Reserved: []uint64{},
		Ab94Grade: Ab05Hypothesis, Ab94ConversionGrade: Ab05Hypothesis, Ab94UnitGrade: Ab05Unknown,
		Ab94RawFallback: "preserve original bytes; no implicit sentinel selection", Ab94Enforcement: Ab05Unknown, Ab94PeerEffects: Ab05Unknown,
	}
}

func ab94Expected() []Ab94NumericField {
	ab94Numeric := ab94ExpectedField("numeric", 16, "little_endian", "divide", 65535, 6553.5)
	ab94Numeric.Ab94Divisor = ab94Pointer(float64(10))
	ab94Context := ab94Numeric
	ab94Context.Ab94Key.Ab94Variant = "synthetic-contextual-replacement"
	ab94Context.Ab94Replacement = []Ab94Replacement{{Ab94Raw: 65535, Ab94Role: "contextual replacement", Ab94Context: "synthetic-only qualified condition", Ab94Grade: Ab05Hypothesis}}
	ab94Min := ab94ExpectedField("minimum_setpoint", 8, "not_applicable", "identity", 255, 255)
	ab94Min.Ab94Key.Ab94Secondary = ab94Pointer(uint8(0xA3))
	ab94Max := ab94Min
	ab94Max.Ab94Key.Ab94Field = "maximum_setpoint"
	ab94Celsius := ab94Min
	ab94Celsius.Ab94Key.Ab94Variant = "synthetic-qualified-celsius"
	ab94Celsius.Ab94Unit = ab94Pointer("degrees C")
	ab94Celsius.Ab94UnitGrade = Ab05Hypothesis
	ab94Celsius.Ab94Replacement = []Ab94Replacement{{Ab94Raw: 255, Ab94Role: "unavailable", Ab94Context: "synthetic-only qualified condition", Ab94Grade: Ab05Hypothesis}}
	ab94Unit := ab94ExpectedField("raw_numeric", 0, "Unknown", "Unknown", 0, 0)
	ab94Unit.Ab94Key.Ab94Variant = "unqualified-unit-2"
	ab94Unit.Ab94WidthBits = nil
	ab94Unit.Ab94Signedness = "Unknown"
	ab94Unit.Ab94UnitCode = ab94Pointer(uint8(2))
	ab94Unit.Ab94RawMin = nil
	ab94Unit.Ab94RawMax = nil
	ab94Unit.Ab94NumericMin = nil
	ab94Unit.Ab94NumericMax = nil
	ab94Unit.Ab94Grade = Ab05Unknown
	ab94Unit.Ab94ConversionGrade = Ab05Unknown
	return []Ab94NumericField{ab94Numeric, ab94Context, ab94Min, ab94Max, ab94Celsius, ab94Unit}
}

func TestAb94ProductionMetadata(ab94T *testing.T) {
	ab94Got, ab94Error := Ab94B555NumericUnits()
	if ab94Error != nil {
		ab94T.Fatal(ab94Error)
	}
	ab94Want := Ab94NumericMetadata{
		Ab94RepresentationGrade: Ab05Proven, Ab94PhysicalQualification: Ab05Unknown, Ab94NativeSupport: Ab05Unknown,
		Ab94Freshness: Ab05Unknown, Ab94Persistence: Ab05Unknown, Ab94Execution: Ab05Unknown,
		Ab94SourceContract: "AI23 supplied functional facts; supplemental declarations only",
		Ab94Limitations:    []string{"Logical positions are unknown; no physical layout is asserted.", "Synthetic contextual variants are examples, not identified device variants.", "No invocation, transport, group/day mapping, persistence or admission is enabled.", "Unit-code 2 has raw numeric material but no qualified semantic conversion."},
		Ab94Entries:        ab94Expected(),
	}
	if !reflect.DeepEqual(ab94Got, ab94Want) {
		ab94T.Fatalf("ab94 production metadata differs: got %+v", ab94Got)
	}
	for _, ab94Expected := range ab94Want.Ab94Entries {
		ab94GotField, ab94Found, ab94Error := Ab94LookupB555NumericUnits(ab94Expected.Ab94Key)
		if ab94Error != nil || !ab94Found || !reflect.DeepEqual(ab94GotField, ab94Expected) {
			ab94T.Fatalf("ab94 exact lookup: %v %v %+v", ab94Error, ab94Found, ab94GotField)
		}
	}
}

func ab94Evaluate(ab94Field Ab94NumericField, ab94Raw []byte) (float64, bool) {
	if ab94Field.Ab94WidthBits == nil || ab94Field.Ab94RawMin == nil || ab94Field.Ab94RawMax == nil {
		return 0, false
	}
	var ab94Number uint64
	switch {
	case *ab94Field.Ab94WidthBits == 16 && ab94Field.Ab94ByteOrder == "little_endian" && len(ab94Raw) == 2:
		ab94Number = uint64(binary.LittleEndian.Uint16(ab94Raw))
	case *ab94Field.Ab94WidthBits == 8 && ab94Field.Ab94ByteOrder == "not_applicable" && len(ab94Raw) == 1:
		ab94Number = uint64(ab94Raw[0])
	default:
		return 0, false
	}
	if ab94Field.Ab94Signedness != "unsigned" || ab94Number < *ab94Field.Ab94RawMin || ab94Number > *ab94Field.Ab94RawMax {
		return 0, false
	}
	ab94Value := float64(ab94Number)
	switch ab94Field.Ab94Conversion {
	case "divide":
		if ab94Field.Ab94Divisor == nil || *ab94Field.Ab94Divisor <= 0 {
			return 0, false
		}
		ab94Value /= *ab94Field.Ab94Divisor
	case "identity":
	default:
		return 0, false
	}
	if ab94Field.Ab94NumericMin == nil || ab94Field.Ab94NumericMax == nil || ab94Value < *ab94Field.Ab94NumericMin || ab94Value > *ab94Field.Ab94NumericMax {
		return 0, false
	}
	return ab94Value, true
}

func TestAb94SyntheticNumbersAndContexts(ab94T *testing.T) {
	ab94Metadata, ab94Error := Ab94B555NumericUnits()
	if ab94Error != nil {
		ab94T.Fatal(ab94Error)
	}
	ab94Rows := ab94Metadata.Ab94Entries
	for ab94Index, ab94Raw := range [][]byte{{0, 0}, {1, 0}, {255, 255}} {
		ab94Value, ab94OK := ab94Evaluate(ab94Rows[0], ab94Raw)
		if !ab94OK || ab94Value != []float64{0, 0.1, 6553.5}[ab94Index] {
			ab94T.Fatalf("ab94 numeric %v: %v %v", ab94Raw, ab94Value, ab94OK)
		}
	}
	for _, ab94Index := range []int{0, 1} {
		ab94Value, ab94OK := ab94Evaluate(ab94Rows[ab94Index], []byte{255, 255})
		if !ab94OK || ab94Value != 6553.5 {
			ab94T.Fatal("ab94 context erased formal numeric alternative")
		}
	}
	if len(ab94Rows[0].Ab94Replacement) != 0 || len(ab94Rows[1].Ab94Replacement) != 1 || ab94Rows[1].Ab94Replacement[0].Ab94Raw != 65535 {
		ab94T.Fatal("ab94 replacement branch conflated")
	}
	for _, ab94Index := range []int{2, 3} {
		ab94Value, ab94OK := ab94Evaluate(ab94Rows[ab94Index], []byte{255})
		if !ab94OK || ab94Value != 255 || ab94Rows[ab94Index].Ab94Unit != nil || len(ab94Rows[ab94Index].Ab94Replacement) != 0 {
			ab94T.Fatal("ab94 generic A3 acquired qualified semantics")
		}
	}
	if ab94Rows[4].Ab94Unit == nil || *ab94Rows[4].Ab94Unit != "degrees C" || ab94Rows[4].Ab94Replacement[0].Ab94Raw != 255 {
		ab94T.Fatal("ab94 synthetic qualification missing")
	}
	ab94Raw := []byte{1, 0}
	ab94Before := append([]byte(nil), ab94Raw...)
	if _, ab94OK := ab94Evaluate(ab94Rows[5], ab94Raw); ab94OK {
		ab94T.Fatal("ab94 unit 2 inherited conversion")
	}
	if !reflect.DeepEqual(ab94Raw, ab94Before) || ab94Rows[5].Ab94Unit != nil || ab94Rows[5].Ab94ConversionGrade != Ab05Unknown || *ab94Rows[5].Ab94UnitCode != 2 {
		ab94T.Fatal("ab94 unit 2 raw/unknown lost")
	}
}

func TestAb94ExactIdentityAxes(ab94T *testing.T) {
	ab94Metadata, ab94Error := Ab94B555NumericUnits()
	if ab94Error != nil {
		ab94T.Fatal(ab94Error)
	}
	for _, ab94Row := range ab94Metadata.Ab94Entries {
		for ab94Index := 0; ab94Index < reflect.TypeOf(ab94Row.Ab94Key).NumField(); ab94Index++ {
			ab94Name := reflect.TypeOf(ab94Row.Ab94Key).Field(ab94Index).Name
			ab94T.Run(ab94Row.Ab94Key.Ab94Variant+"/"+ab94Row.Ab94Key.Ab94Field+"/"+ab94Name, func(ab94T *testing.T) {
				ab94Key := ab94Row.Ab94Key
				ab94Axis := reflect.ValueOf(&ab94Key).Elem().Field(ab94Index)
				switch ab94Axis.Kind() {
				case reflect.String:
					ab94Axis.SetString("ab94_distinct")
				case reflect.Uint16:
					ab94Axis.SetUint(0)
				case reflect.Ptr:
					ab94Axis.Set(reflect.New(ab94Axis.Type().Elem()))
				}
				if reflect.DeepEqual(ab94Key, ab94Row.Ab94Key) {
					ab94T.Fatal("ab94 axis not varied")
				}
				if _, ab94Found, ab94Error := Ab94LookupB555NumericUnits(ab94Key); ab94Error != nil || ab94Found {
					ab94T.Fatalf("ab94 axis ignored: %v %v", ab94Found, ab94Error)
				}
				if _, ab94Found, ab94Error := Ab94LookupB555NumericUnits(ab94Row.Ab94Key); ab94Error != nil || !ab94Found {
					ab94T.Fatal("ab94 valid control failed")
				}
			})
		}
	}
	// Synthetic B524 coordinates exercise the same matcher without declaring a B524 protocol row.
	ab94Key := Ab94Key{Ab94Family: "ab94_family", Ab94Variant: "ab94_variant", Ab94Primary: 0xB524, Ab94Secondary: ab94Pointer(uint8(0)), Ab94Operation: ab94Pointer(uint8(0)), Ab94Direction: ab94Pointer("response"), Ab94Field: "ab94_field", Ab94Opcode: ab94Pointer(uint8(0)), Ab94Group: ab94Pointer(uint8(0)), Ab94Instance: ab94Pointer(uint8(0)), Ab94Register: ab94Pointer(uint16(0))}
	for ab94Index := 3; ab94Index < reflect.TypeOf(ab94Key).NumField(); ab94Index++ {
		if reflect.ValueOf(ab94Key).Field(ab94Index).Kind() != reflect.Ptr {
			continue
		}
		for _, ab94Absent := range []bool{false, true} {
			ab94Changed := ab94Key
			ab94Axis := reflect.ValueOf(&ab94Changed).Elem().Field(ab94Index)
			if ab94Absent {
				ab94Axis.SetZero()
			} else {
				ab94Axis.Set(reflect.New(ab94Axis.Type().Elem()))
				if ab94Axis.Elem().Kind() == reflect.String {
					ab94Axis.Elem().SetString("ab94_other")
				} else {
					ab94Axis.Elem().SetUint(1)
				}
			}
			ab94Rows := []Ab94NumericField{{Ab94Key: ab94Key}}
			if _, ab94Found := ab94Lookup(ab94Rows, ab94Changed); ab94Found {
				ab94T.Fatalf("ab94 B524 axis %s ignored (absent=%v)", reflect.TypeOf(ab94Key).Field(ab94Index).Name, ab94Absent)
			}
			if _, ab94Found := ab94Lookup(ab94Rows, ab94Key); !ab94Found {
				ab94T.Fatal("ab94 B524 control failed")
			}
		}
	}
}

func ab94Mutate(ab94Value reflect.Value) {
	switch ab94Value.Kind() {
	case reflect.Struct:
		for ab94Index := 0; ab94Index < ab94Value.NumField(); ab94Index++ {
			ab94Mutate(ab94Value.Field(ab94Index))
		}
	case reflect.Ptr:
		if !ab94Value.IsNil() {
			ab94Mutate(ab94Value.Elem())
		}
	case reflect.Slice:
		for ab94Index := 0; ab94Index < ab94Value.Len(); ab94Index++ {
			ab94Mutate(ab94Value.Index(ab94Index))
		}
		if ab94Value.Type().Elem().Kind() == reflect.Uint64 {
			ab94Value.Set(reflect.Append(ab94Value, reflect.ValueOf(uint64(7))))
		}
	case reflect.Map:
		ab94Value.SetMapIndex(reflect.ValueOf(uint64(7)), reflect.ValueOf("ab94_changed"))
	case reflect.String:
		ab94Value.SetString("ab94_changed")
	case reflect.Uint8, reflect.Uint16, reflect.Uint64:
		ab94Value.SetUint(7)
	case reflect.Int:
		ab94Value.SetInt(7)
	case reflect.Float64:
		ab94Value.SetFloat(7)
	case reflect.Bool:
		ab94Value.SetBool(true)
	}
}

func TestAb94DeepCopyIsolation(ab94T *testing.T) {
	ab94Before, ab94Error := Ab94B555NumericUnits()
	if ab94Error != nil {
		ab94T.Fatal(ab94Error)
	}
	ab94Snapshot, ab94Error := json.Marshal(ab94Before)
	if ab94Error != nil {
		ab94T.Fatal(ab94Error)
	}
	ab94Returned, ab94Error := Ab94B555NumericUnits()
	if ab94Error != nil {
		ab94T.Fatal(ab94Error)
	}
	ab94Mutate(reflect.ValueOf(&ab94Returned).Elem())
	for _, ab94Row := range ab94Before.Ab94Entries {
		ab94LookupResult, ab94Found, ab94Error := Ab94LookupB555NumericUnits(ab94Row.Ab94Key)
		if ab94Error != nil || !ab94Found {
			ab94T.Fatal("ab94 lookup failed")
		}
		ab94Mutate(reflect.ValueOf(&ab94LookupResult).Elem())
		ab94Again, ab94Found, ab94Error := Ab94LookupB555NumericUnits(ab94Row.Ab94Key)
		if ab94Error != nil || !ab94Found || !reflect.DeepEqual(ab94Again, ab94Row) {
			ab94T.Fatal("ab94 lookup shared nested storage")
		}
	}
	ab94After, ab94Error := Ab94B555NumericUnits()
	if ab94Error != nil {
		ab94T.Fatal(ab94Error)
	}
	ab94AfterBytes, ab94Error := json.Marshal(ab94After)
	if ab94Error != nil {
		ab94T.Fatal(ab94Error)
	}
	if string(ab94Snapshot) != string(ab94AfterBytes) {
		ab94T.Fatal("ab94 accessor shared nested storage")
	}
}
