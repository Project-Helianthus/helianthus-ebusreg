package vaillant

import (
	ab87JSON "encoding/json"
	ab87Reflect "reflect"
	ab87Testing "testing"
)

func ab87Number(ab87Value uint16) *uint16 { return &ab87Value }
func ab87Int(ab87Value int) *int          { return &ab87Value }
func ab87Bound(ab87Value int64) *int64    { return &ab87Value }

func ab87Key(ab87Operation uint16, ab87Direction, ab87Field string) Ab87DiagnosticKey {
	return Ab87DiagnosticKey{Ab87Family: "vaillant", Ab87Variant: "B509", Ab87Primary: ab87Number(0xB5), Ab87Secondary: ab87Number(9), Ab87Operation: ab87Number(ab87Operation), Ab87Direction: ab87Direction, Ab87Field: ab87Field}
}

func ab87Get(ab87T *ab87Testing.T, ab87Operation uint16, ab87Direction, ab87Name string) Ab87DiagnosticField {
	ab87T.Helper()
	ab87Field, ab87OK, ab87Error := Ab87LookupB509DiagnosticField(ab87Key(ab87Operation, ab87Direction, ab87Name))
	if ab87Error != nil || !ab87OK {
		ab87T.Fatalf("ab87 missing %s: %v", ab87Name, ab87Error)
	}
	return ab87Field
}

func ab87Expected(ab87Operation uint16, ab87Direction, ab87Name string, ab87Width *int, ab87Sign, ab87Order, ab87Conversion string, ab87Min, ab87Max *int64, ab87Mapping string) Ab87DiagnosticField {
	return Ab87DiagnosticField{Ab87Key: ab87Key(ab87Operation, ab87Direction, ab87Name), Ab87WidthBits: ab87Width, Ab87Signedness: ab87Sign, Ab87ByteOrder: ab87Order, Ab87Conversion: ab87Conversion, Ab87Minimum: ab87Min, Ab87Maximum: ab87Max, Ab87Grade: "Hypothesis", Ab87MappingGrade: ab87Mapping, Ab87RawFallback: "preserve original bytes"}
}

func TestAb87ProductionContract(ab87T *ab87Testing.T) {
	ab87Catalog, ab87Error := Ab87B509DiagnosticShapes()
	if ab87Error != nil {
		ab87T.Fatal(ab87Error)
	}
	ab87ExpectedFields := []Ab87DiagnosticField{
		ab87Expected(0x13, "request", "arguments", ab87Int(0), "Unknown", "Unknown", "no_arguments", nil, nil, "Unknown"),
		ab87Expected(0x13, "response", "result", nil, "Unknown", "Unknown", "Boolean", nil, nil, "Unknown"),
		ab87Expected(0x34, "request", "address", ab87Int(32), "unsigned", "little_endian", "identity", ab87Bound(0), ab87Bound(4294967295), "Hypothesis"),
		ab87Expected(0x34, "request", "requested_count", ab87Int(8), "unsigned", "not_applicable", "identity", ab87Bound(0), ab87Bound(255), "Hypothesis"),
		ab87Expected(0x34, "response", "address", ab87Int(32), "unsigned", "little_endian", "identity", ab87Bound(0), ab87Bound(4294967295), "Hypothesis"),
		ab87Expected(0x34, "response", "data", nil, "Unknown", "Unknown", "raw_bytes", nil, nil, "Unknown"),
		ab87Expected(0x35, "request", "address", ab87Int(32), "unsigned", "little_endian", "identity", ab87Bound(0), ab87Bound(4294967295), "Hypothesis"),
		ab87Expected(0x35, "request", "data", nil, "Unknown", "Unknown", "raw_bytes", nil, nil, "Unknown"),
		ab87Expected(0x35, "response", "write_status", nil, "Unknown", "Unknown", "status_domain", ab87Bound(0), ab87Bound(6), "Unknown"),
		ab87Expected(0x51, "Unknown", "component", nil, "Unknown", "Unknown", "identity", ab87Bound(0), ab87Bound(7), "Hypothesis"),
		ab87Expected(0x51, "Unknown", "master", ab87Int(8), "unsigned", "not_applicable", "identity", ab87Bound(0), ab87Bound(255), "Hypothesis"),
		ab87Expected(0x51, "Unknown", "identity", ab87Int(40), "not_applicable", "not_applicable", "ASCII_exact", nil, nil, "Hypothesis"),
		ab87Expected(0x51, "Unknown", "software_component_1", nil, "Unknown", "Unknown", "decimal_scalar", nil, nil, "Unknown"),
		ab87Expected(0x51, "Unknown", "software_component_2", nil, "Unknown", "Unknown", "decimal_scalar", nil, nil, "Unknown"),
		ab87Expected(0x51, "Unknown", "software_component_3", nil, "Unknown", "Unknown", "decimal_scalar", nil, nil, "Unknown"),
		ab87Expected(0x51, "Unknown", "software_component_4", nil, "Unknown", "Unknown", "decimal_scalar", nil, nil, "Unknown"),
	}
	ab87ExpectedFields[5].Ab87RecommendedBytes = ab87Int(10)
	ab87ExpectedFields[8].Ab87Enumeration = map[string]*string{"0": nil, "1": nil, "2": nil, "3": nil, "4": nil, "5": nil, "6": nil}
	if !ab87Reflect.DeepEqual(ab87Catalog.Ab87Fields, ab87ExpectedFields) {
		ab87T.Fatal("ab87 finite metadata differs from independent contract")
	}
	if ab87Catalog.Ab87Representation != "Proven" || ab87Catalog.Ab87Qualification != "Unknown" || ab87Catalog.Ab87Execution != "Unknown" || ab87Catalog.Ab87PhysicalEffects != "Unknown" || ab87Catalog.Ab87ResponseCorrelation != "Unknown" || len(ab87Catalog.Ab87Limitations) != 6 {
		ab87T.Fatal("ab87 evidence boundaries changed")
	}
	for _, ab87Field := range ab87ExpectedFields {
		ab87Got, ab87OK, ab87Error := Ab87LookupB509DiagnosticField(ab87Field.Ab87Key)
		if ab87Error != nil || !ab87OK || !ab87Reflect.DeepEqual(ab87Got, ab87Field) {
			ab87T.Fatal("ab87 exact lookup changed")
		}
	}
	if _, ab87OK, ab87Error := Ab87LookupB509DiagnosticField(ab87Key(0x35, "request", "requested_count")); ab87Error != nil || ab87OK {
		ab87T.Fatal("ab87 OP35 must not gain OP34 count")
	}
}

func ab87Unsigned(ab87Field Ab87DiagnosticField, ab87Raw []byte) (uint64, bool) {
	if ab87Field.Ab87WidthBits == nil || *ab87Field.Ab87WidthBits != 8*len(ab87Raw) || ab87Field.Ab87Signedness != "unsigned" || ab87Field.Ab87Conversion != "identity" {
		return 0, false
	}
	var ab87Value uint64
	switch ab87Field.Ab87ByteOrder {
	case "little_endian":
		for ab87Index, ab87Byte := range ab87Raw {
			ab87Value |= uint64(ab87Byte) << uint(8*ab87Index)
		}
	case "not_applicable":
		if len(ab87Raw) != 1 {
			return 0, false
		}
		ab87Value = uint64(ab87Raw[0])
	default:
		return 0, false
	}
	return ab87Value, true
}

func ab87InBounds(ab87Field Ab87DiagnosticField, ab87Value int64) bool {
	return (ab87Field.Ab87Minimum == nil || ab87Value >= *ab87Field.Ab87Minimum) && (ab87Field.Ab87Maximum == nil || ab87Value <= *ab87Field.Ab87Maximum)
}

func ab87FitsASCII(ab87Field Ab87DiagnosticField, ab87Value string) bool {
	if ab87Field.Ab87Conversion != "ASCII_exact" || ab87Field.Ab87WidthBits == nil || len(ab87Value)*8 != *ab87Field.Ab87WidthBits {
		return false
	}
	for _, ab87Byte := range []byte(ab87Value) {
		if ab87Byte > 127 {
			return false
		}
	}
	return true
}

func TestAb87SyntheticCases(ab87T *ab87Testing.T) {
	for _, ab87Case := range []struct {
		Ab87Operation uint16
		Ab87Direction string
	}{{0x34, "request"}, {0x34, "response"}, {0x35, "request"}} {
		ab87Address := ab87Get(ab87T, ab87Case.Ab87Operation, ab87Case.Ab87Direction, "address")
		ab87Value, ab87OK := ab87Unsigned(ab87Address, []byte{0x78, 0x56, 0x34, 0x12})
		if !ab87OK || ab87Value != 0x12345678 {
			ab87T.Fatal("ab87 u32 little endian address")
		}
	}
	ab87Count := ab87Get(ab87T, 0x34, "request", "requested_count")
	for _, ab87Raw := range []byte{0, 255} {
		ab87Value, ab87OK := ab87Unsigned(ab87Count, []byte{ab87Raw})
		if !ab87OK || ab87Value != uint64(ab87Raw) || !ab87InBounds(ab87Count, int64(ab87Value)) {
			ab87T.Fatal("ab87 raw count changed")
		}
	}
	ab87Data := ab87Get(ab87T, 0x34, "response", "data")
	ab87Synthetic := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if ab87Data.Ab87Conversion != "raw_bytes" || ab87Data.Ab87RecommendedBytes == nil || *ab87Data.Ab87RecommendedBytes != 10 || ab87Data.Ab87MaximumBytes != nil || (ab87Data.Ab87WidthBits != nil && *ab87Data.Ab87WidthBits != 8*len(ab87Synthetic)) {
		ab87T.Fatal("ab87 recommendation became capacity")
	}
	ab87Status := ab87Get(ab87T, 0x35, "response", "write_status")
	if !ab87InBounds(ab87Status, 6) || ab87InBounds(ab87Status, 7) {
		ab87T.Fatal("ab87 status domain")
	}
	if _, ab87OK := ab87Status.Ab87Enumeration["7"]; ab87OK {
		ab87T.Fatal("ab87 status 7 labeled")
	}
	ab87Identity := ab87Get(ab87T, 0x51, "Unknown", "identity")
	if !ab87FitsASCII(ab87Identity, "ABCDE") || ab87FitsASCII(ab87Identity, "AB") || ab87FitsASCII(ab87Identity, "ABCDEF") || ab87FitsASCII(ab87Identity, "ABCDé") || ab87Identity.Ab87Padding != nil {
		ab87T.Fatal("ab87 exact ASCII identity")
	}
	ab87Component := ab87Get(ab87T, 0x51, "Unknown", "component")
	if !ab87InBounds(ab87Component, 7) || ab87InBounds(ab87Component, 8) {
		ab87T.Fatal("ab87 component range")
	}
	ab87Scalars := []int64{1, 20, 3, 40}
	for ab87Index, ab87Name := range []string{"software_component_1", "software_component_2", "software_component_3", "software_component_4"} {
		ab87Field := ab87Get(ab87T, 0x51, "Unknown", ab87Name)
		if ab87Field.Ab87Conversion != "decimal_scalar" || !ab87InBounds(ab87Field, ab87Scalars[ab87Index]) || ab87Field.Ab87WidthBits != nil {
			ab87T.Fatal("ab87 four independent decimal scalars")
		}
	}
}

func TestAb87ExactIdentityAxes(ab87T *ab87Testing.T) {
	ab87Control := ab87Get(ab87T, 0x34, "request", "address")
	for _, ab87B524 := range []bool{false, true} {
		ab87Row := ab87Control
		if ab87B524 {
			ab87Row.Ab87Key.Ab87Variant = "B524"
			ab87Row.Ab87Key.Ab87Secondary = ab87Number(0x24)
			ab87Row.Ab87Key.Ab87Opcode = ab87Number(0)
			ab87Row.Ab87Key.Ab87Group = ab87Number(0)
			ab87Row.Ab87Key.Ab87Instance = ab87Number(0)
			ab87Row.Ab87Key.Ab87Register = ab87Number(0)
		}
		ab87Rows := []Ab87DiagnosticField{ab87Row}
		ab87Type := ab87Reflect.TypeOf(ab87Row.Ab87Key)
		for ab87Index := 0; ab87Index < ab87Type.NumField(); ab87Index++ {
			ab87T.Run("ab87_"+ab87Row.Ab87Key.Ab87Variant+"_"+ab87Type.Field(ab87Index).Name, func(ab87T *ab87Testing.T) {
				if _, ab87OK, ab87Error := ab87LookupField(ab87Rows, ab87Row.Ab87Key); ab87Error != nil || !ab87OK {
					ab87T.Fatal("ab87 invalid control")
				}
				ab87Changed := ab87Row.Ab87Key
				ab87Axis := ab87Reflect.ValueOf(&ab87Changed).Elem().Field(ab87Index)
				if ab87Axis.Kind() == ab87Reflect.String {
					ab87Axis.SetString("ab87_distinct")
				} else if ab87Axis.IsNil() {
					ab87Axis.Set(ab87Reflect.ValueOf(ab87Number(0)))
				} else {
					ab87Axis.Set(ab87Reflect.Zero(ab87Axis.Type()))
				}
				if _, ab87OK, ab87Error := ab87LookupField(ab87Rows, ab87Changed); ab87Error != nil || ab87OK {
					ab87T.Fatal("ab87 axis aliased absent/zero")
				}
				if ab87Axis.Kind() == ab87Reflect.Ptr {
					ab87Axis.Set(ab87Reflect.ValueOf(ab87Number(0)))
					if !ab87Reflect.DeepEqual(ab87Changed, ab87Row.Ab87Key) {
						if _, ab87OK, _ := ab87LookupField(ab87Rows, ab87Changed); ab87OK {
							ab87T.Fatal("ab87 zero wildcard")
						}
					}
					ab87Axis.Set(ab87Reflect.ValueOf(ab87Number(65535)))
					if _, ab87OK, _ := ab87LookupField(ab87Rows, ab87Changed); ab87OK {
						ab87T.Fatal("ab87 numeric axis ignored")
					}
				}
			})
		}
	}
}

func TestAb87DeepCopy(ab87T *ab87Testing.T) {
	ab87Before, ab87Error := Ab87B509DiagnosticShapes()
	if ab87Error != nil {
		ab87T.Fatal(ab87Error)
	}
	ab87Copy, _ := Ab87B509DiagnosticShapes()
	ab87Copy.Ab87Limitations[0] = "ab87_changed"
	*ab87Copy.Ab87Fields[2].Ab87Key.Ab87Operation = 0
	*ab87Copy.Ab87Fields[2].Ab87WidthBits = 16
	*ab87Copy.Ab87Fields[2].Ab87Maximum = 1
	*ab87Copy.Ab87Fields[5].Ab87RecommendedBytes = 1
	ab87Copy.Ab87Fields[8].Ab87Enumeration["7"] = nil
	ab87Copy.Ab87Fields[0] = Ab87DiagnosticField{}
	ab87After, _ := Ab87B509DiagnosticShapes()
	if !ab87Reflect.DeepEqual(ab87After, ab87Before) {
		ab87T.Fatal("ab87 accessor shared backing")
	}
	ab87Lookup := ab87Get(ab87T, 0x35, "response", "write_status")
	*ab87Lookup.Ab87Key.Ab87Primary = 0
	*ab87Lookup.Ab87Maximum = 7
	ab87Lookup.Ab87Enumeration["7"] = nil
	if ab87Reflect.DeepEqual(ab87Get(ab87T, 0x35, "response", "write_status"), ab87Lookup) {
		ab87T.Fatal("ab87 lookup shared backing")
	}
	ab87Source := ab87Before.Ab87Fields[8]
	ab87Label := "ab87_label"
	ab87Unit := "ab87_unit"
	ab87Padding := "ab87_padding"
	ab87Source.Ab87Enumeration = map[string]*string{"6": &ab87Label}
	ab87Source.Ab87Replacement = []int64{1}
	ab87Source.Ab87Reserved = []int64{2}
	ab87Source.Ab87LogicalPosition = ab87Int(0)
	ab87Source.Ab87Unit = &ab87Unit
	ab87Source.Ab87Default = ab87Bound(0)
	ab87Source.Ab87MaximumBytes = ab87Int(11)
	ab87Source.Ab87Padding = &ab87Padding
	ab87Snapshot, _ := ab87JSON.Marshal(ab87Source)
	ab87Detached, ab87OK, ab87Error := ab87LookupField([]Ab87DiagnosticField{ab87Source}, ab87Source.Ab87Key)
	if ab87Error != nil || !ab87OK {
		ab87T.Fatal("ab87 synthetic copy control")
	}
	*ab87Detached.Ab87Enumeration["6"] = "ab87_changed"
	ab87Detached.Ab87Replacement[0] = 9
	ab87Detached.Ab87Reserved[0] = 9
	*ab87Detached.Ab87LogicalPosition = 9
	*ab87Detached.Ab87Unit = "ab87_changed"
	*ab87Detached.Ab87Default = 9
	*ab87Detached.Ab87MaximumBytes = 9
	*ab87Detached.Ab87Padding = "ab87_changed"
	ab87Again, _, _ := ab87LookupField([]Ab87DiagnosticField{ab87Source}, ab87Source.Ab87Key)
	ab87Actual, _ := ab87JSON.Marshal(ab87Again)
	if string(ab87Actual) != string(ab87Snapshot) {
		ab87T.Fatal("ab87 nested lookup storage shared")
	}
}
