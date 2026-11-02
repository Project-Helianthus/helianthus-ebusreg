package vaillant

import (
	ab93Binary "encoding/binary"
	ab93JSON "encoding/json"
	ab93Fmt "fmt"
	ab93Reflect "reflect"
	ab93Testing "testing"
)

func ab93Ptr[Ab93Value any](ab93Value Ab93Value) *Ab93Value { return &ab93Value }

func ab93Catalog(ab93T *ab93Testing.T) Ab93SensorCatalog {
	ab93T.Helper()
	ab93Catalog, ab93Error := Ab93B523SensorFields()
	if ab93Error != nil {
		ab93T.Fatal(ab93Error)
	}
	return ab93Catalog
}

func ab93ExpectedBase(ab93Variant string, ab93Operation uint16, ab93Name string, ab93Position, ab93Width int) Ab93SensorField {
	return Ab93SensorField{
		Ab93Key:      Ab93SensorKey{Ab93Family: "Vaillant", Ab93Variant: ab93Variant, Ab93Primary: ab93Ptr(uint16(181)), Ab93Secondary: ab93Ptr(uint16(35)), Ab93Operation: ab93Ptr(ab93Operation), Ab93Field: ab93Name},
		Ab93Position: ab93Ptr(ab93Position), Ab93Width: ab93Ptr(ab93Width), Ab93Signedness: "Unknown", Ab93ByteOrder: "Unknown", Ab93Enumeration: map[string]string{}, Ab93Selection: "none", Ab93Branches: []Ab93SensorBranch{}, Ab93Grade: "Hypothesis", Ab93RawFallback: "preserve original bytes including reserved bits",
	}
}

func ab93ExpectedFields() []Ab93SensorField {
	var ab93Fields []Ab93SensorField
	for ab93Variant, ab93Zero := range map[string]string{"zero-valid-value": "valid-value", "zero-valid-value-or-low-level": "valid-value-or-low-level"} {
		for _, ab93Operation := range []struct {
			Ab93Op                            uint16
			Ab93Sensors, Ab93Positions        []int
			Ab93First, Ab93Last, Ab93Reserved int
		}{
			{3, []int{1, 2, 3, 4, 5, 6}, []int{3, 5, 7, 9, 11, 13}, 16, 17, 4},
			{6, []int{1, 2, 3, 4, 5, 6, 7}, []int{3, 5, 7, 9, 11, 13, 15}, 17, 18, 6},
			{7, []int{8, 9, 10, 11, 12, 13}, []int{3, 5, 7, 9, 12, 14}, 16, 17, 4},
		} {
			for ab93Index, ab93Sensor := range ab93Operation.Ab93Sensors {
				ab93Field := ab93ExpectedBase(ab93Variant, ab93Operation.Ab93Op, ab93Fmt.Sprintf("S%d", ab93Sensor), ab93Operation.Ab93Positions[ab93Index], 16)
				ab93Field.Ab93ReplacementRaw = ab93Ptr(uint16(32768))
				ab93Field.Ab93Selection = "unsigned raw equality before signed conversion"
				ab93Field.Ab93Branches = []Ab93SensorBranch{
					{Ab93Name: "replacement", Ab93Signedness: "unsigned", Ab93ByteOrder: "little-endian", Ab93DivideBy: ab93Ptr(1.0), Ab93Default: ab93Ptr(32768.0), Ab93Bounds: []float64{0, 65535}, Ab93Grade: "Hypothesis"},
					{Ab93Name: "normal", Ab93Signedness: "signed", Ab93ByteOrder: "little-endian", Ab93DivideBy: ab93Ptr(16.0), Ab93Unit: ab93Ptr("Celsius"), Ab93Default: ab93Ptr(0.0), Ab93Bounds: []float64{-32768, 32767}, Ab93Grade: "Hypothesis"},
				}
				ab93Fields = append(ab93Fields, ab93Field)
				ab93Position := ab93Operation.Ab93First
				if ab93Index >= 4 {
					ab93Position = ab93Operation.Ab93Last
				}
				ab93Field = ab93ExpectedBase(ab93Variant, ab93Operation.Ab93Op, ab93Fmt.Sprintf("S%d.status", ab93Sensor), ab93Position, 2)
				ab93Field.Ab93BitOffset = ab93Ptr(ab93Index % 4 * 2)
				ab93Field.Ab93Signedness = "unsigned"
				ab93Field.Ab93ByteOrder = "not_applicable"
				ab93Field.Ab93DivideBy = ab93Ptr(1.0)
				ab93Field.Ab93Default = ab93Ptr(0.0)
				ab93Field.Ab93Bounds = []float64{0, 3}
				ab93Field.Ab93Enumeration = map[string]string{"0": ab93Zero, "1": "open-circuit", "2": "closed-circuit", "3": "high-level"}
				ab93Fields = append(ab93Fields, ab93Field)
			}
			ab93Field := ab93ExpectedBase(ab93Variant, ab93Operation.Ab93Op, "status.reserved", ab93Operation.Ab93Last, 8-ab93Operation.Ab93Reserved)
			ab93Field.Ab93BitOffset = ab93Ptr(ab93Operation.Ab93Reserved)
			ab93Field.Ab93Reserved = true
			ab93Field.Ab93Signedness = "unsigned"
			ab93Field.Ab93ByteOrder = "not_applicable"
			ab93Fields = append(ab93Fields, ab93Field)
			if ab93Operation.Ab93Op == 3 {
				ab93Field = ab93ExpectedBase(ab93Variant, 3, "PWM", 15, 8)
				ab93Field.Ab93Signedness = "unsigned"
				ab93Field.Ab93ByteOrder = "not_applicable"
				ab93Field.Ab93DivideBy = ab93Ptr(2.0)
				ab93Field.Ab93Unit = ab93Ptr("percent")
				ab93Field.Ab93Bounds = []float64{0, 255}
				ab93Fields = append(ab93Fields, ab93Field)
			}
			if ab93Operation.Ab93Op == 7 {
				ab93Field = ab93ExpectedBase(ab93Variant, 7, "opaque", 11, 8)
				ab93Field.Ab93Signedness = "unsigned"
				ab93Field.Ab93ByteOrder = "not_applicable"
				ab93Field.Ab93Bounds = []float64{0, 255}
				ab93Fields = append(ab93Fields, ab93Field)
			}
		}
	}
	return ab93Fields
}

func TestAb93ProductionMetadata(ab93T *ab93Testing.T) {
	ab93Catalog := ab93Catalog(ab93T)
	if ab93Catalog.Ab93RepresentationGrade != "Proven" || ab93Catalog.Ab93Qualification != "Unknown" || ab93Catalog.Ab93Complete || len(ab93Catalog.Ab93Limitations) != 6 {
		ab93T.Fatal("qualification/completeness boundary")
	}
	ab93Expected := ab93ExpectedFields()
	if len(ab93Catalog.Ab93Fields) != 86 || len(ab93Expected) != 86 {
		ab93T.Fatal("finite field count")
	}
	ab93Seen := map[string]bool{}
	for _, ab93Want := range ab93Expected {
		ab93Got, ab93Found, ab93Error := Ab93LookupB523SensorField(ab93Want.Ab93Key)
		if ab93Error != nil || !ab93Found || !ab93Reflect.DeepEqual(ab93Got, ab93Want) {
			ab93T.Fatalf("metadata mismatch %v: got %+v want %+v, %v", ab93Want.Ab93Key, ab93Got, ab93Want, ab93Error)
		}
		ab93Bytes, _ := ab93JSON.Marshal(ab93Want.Ab93Key)
		if ab93Seen[string(ab93Bytes)] {
			ab93T.Fatal("duplicate key")
		}
		ab93Seen[string(ab93Bytes)] = true
	}
	for _, ab93Field := range ab93Catalog.Ab93Fields {
		ab93Bytes, _ := ab93JSON.Marshal(ab93Field.Ab93Key)
		if !ab93Seen[string(ab93Bytes)] {
			ab93T.Fatal("unexpected key")
		}
	}
}

func ab93Evaluate(ab93T *ab93Testing.T, ab93Field Ab93SensorField, ab93Bytes []byte) (string, float64, string, []byte) {
	ab93T.Helper()
	ab93Raw := append([]byte(nil), ab93Bytes...)
	if ab93Field.Ab93RawFallback != "preserve original bytes including reserved bits" {
		ab93T.Fatal("raw preservation")
	}
	if ab93Field.Ab93ReplacementRaw != nil {
		ab93Word := ab93Binary.LittleEndian.Uint16(ab93Bytes[*ab93Field.Ab93Position:])
		ab93Branch := ab93Field.Ab93Branches[1]
		if ab93Word == *ab93Field.Ab93ReplacementRaw {
			ab93Branch = ab93Field.Ab93Branches[0]
		}
		ab93Value := float64(ab93Word)
		if ab93Branch.Ab93Signedness == "signed" {
			ab93Value = float64(int16(ab93Word))
		}
		return ab93Branch.Ab93Name, ab93Value / *ab93Branch.Ab93DivideBy, "", ab93Raw
	}
	ab93Value := uint16(ab93Bytes[*ab93Field.Ab93Position])
	if ab93Field.Ab93BitOffset != nil {
		ab93Value = (ab93Value >> *ab93Field.Ab93BitOffset) & ((1 << *ab93Field.Ab93Width) - 1)
	}
	if ab93Field.Ab93Reserved {
		return "reserved", float64(ab93Value), "", ab93Raw
	}
	ab93Number := float64(ab93Value)
	if ab93Field.Ab93DivideBy != nil {
		ab93Number /= *ab93Field.Ab93DivideBy
	}
	return "value", ab93Number, ab93Field.Ab93Enumeration[ab93Fmt.Sprint(ab93Value)], ab93Raw
}

func TestAb93SyntheticSensorWords(ab93T *ab93Testing.T) {
	ab93Fields := ab93Catalog(ab93T).Ab93Fields
	ab93Count := 0
	for _, ab93Target := range ab93Fields {
		if ab93Target.Ab93ReplacementRaw == nil {
			continue
		}
		ab93Count++
		for _, ab93Case := range []struct {
			Ab93Bytes  []byte
			Ab93Branch string
			Ab93Value  float64
		}{{[]byte{0, 128}, "replacement", 32768}, {[]byte{16, 0}, "normal", 1}, {[]byte{240, 255}, "normal", -1}} {
			ab93T.Run(ab93Fmt.Sprintf("ab93_%s_%d_%s_%s", ab93Target.Ab93Key.Ab93Variant, *ab93Target.Ab93Key.Ab93Operation, ab93Target.Ab93Key.Ab93Field, ab93Case.Ab93Branch+ab93Fmt.Sprint(ab93Case.Ab93Value)), func(ab93T *ab93Testing.T) {
				ab93Bytes := make([]byte, 20)
				copy(ab93Bytes[*ab93Target.Ab93Position:], ab93Case.Ab93Bytes)
				for _, ab93Field := range ab93Fields {
					if ab93Field.Ab93ReplacementRaw == nil || ab93Field.Ab93Key.Ab93Variant != ab93Target.Ab93Key.Ab93Variant || *ab93Field.Ab93Key.Ab93Operation != *ab93Target.Ab93Key.Ab93Operation {
						continue
					}
					ab93Branch, ab93Value, _, ab93Raw := ab93Evaluate(ab93T, ab93Field, ab93Bytes)
					ab93WantBranch, ab93WantValue := "normal", 0.0
					if ab93Field.Ab93Key.Ab93Field == ab93Target.Ab93Key.Ab93Field {
						ab93WantBranch, ab93WantValue = ab93Case.Ab93Branch, ab93Case.Ab93Value
					}
					if ab93Branch != ab93WantBranch || ab93Value != ab93WantValue || !ab93Reflect.DeepEqual(ab93Raw, ab93Bytes) {
						ab93T.Fatalf("word isolation: %s %s %v", ab93Field.Ab93Key.Ab93Field, ab93Branch, ab93Value)
					}
				}
			})
		}
	}
	if ab93Count != 38 {
		ab93T.Fatalf("executed %d sensor occurrences, want 19 per variant", ab93Count)
	}
}

func TestAb93SyntheticStatusAndReserved(ab93T *ab93Testing.T) {
	ab93Fields := ab93Catalog(ab93T).Ab93Fields
	ab93Count, ab93ReservedCount := 0, 0
	for _, ab93Target := range ab93Fields {
		if len(ab93Target.Ab93Enumeration) == 0 && !ab93Target.Ab93Reserved {
			continue
		}
		if ab93Target.Ab93Reserved {
			for ab93Bit := *ab93Target.Ab93BitOffset; ab93Bit < 8; ab93Bit++ {
				ab93ReservedCount++
				ab93Bytes := make([]byte, 20)
				ab93Bytes[*ab93Target.Ab93Position] = 1 << ab93Bit
				ab93Branch, ab93Value, ab93Label, ab93Raw := ab93Evaluate(ab93T, ab93Target, ab93Bytes)
				if ab93Branch != "reserved" || ab93Value != float64(uint16(1)<<(ab93Bit-*ab93Target.Ab93BitOffset)) || ab93Label != "" || !ab93Reflect.DeepEqual(ab93Raw, ab93Bytes) {
					ab93T.Fatal("reserved bit not preserved")
				}
				for _, ab93Field := range ab93Fields {
					if len(ab93Field.Ab93Enumeration) == 0 || ab93Field.Ab93Key.Ab93Variant != ab93Target.Ab93Key.Ab93Variant || *ab93Field.Ab93Key.Ab93Operation != *ab93Target.Ab93Key.Ab93Operation {
						continue
					}
					_, ab93Value, _, ab93Raw := ab93Evaluate(ab93T, ab93Field, ab93Bytes)
					if ab93Value != 0 || !ab93Reflect.DeepEqual(ab93Raw, ab93Bytes) {
						ab93T.Fatal("reserved bit became sensor state")
					}
				}
			}
			continue
		}
		ab93Count++
		for ab93State := 0; ab93State < 4; ab93State++ {
			ab93Bytes := make([]byte, 20)
			ab93Bytes[*ab93Target.Ab93Position] = byte(ab93State << *ab93Target.Ab93BitOffset)
			for _, ab93Field := range ab93Fields {
				if len(ab93Field.Ab93Enumeration) == 0 || ab93Field.Ab93Key.Ab93Variant != ab93Target.Ab93Key.Ab93Variant || *ab93Field.Ab93Key.Ab93Operation != *ab93Target.Ab93Key.Ab93Operation {
					continue
				}
				ab93Want := 0
				if ab93Field.Ab93Key.Ab93Field == ab93Target.Ab93Key.Ab93Field {
					ab93Want = ab93State
				}
				ab93Labels := []string{"valid-value", "open-circuit", "closed-circuit", "high-level"}
				if ab93Field.Ab93Key.Ab93Variant == "zero-valid-value-or-low-level" {
					ab93Labels[0] = "valid-value-or-low-level"
				}
				_, ab93Value, ab93Label, ab93Raw := ab93Evaluate(ab93T, ab93Field, ab93Bytes)
				if ab93Value != float64(ab93Want) || ab93Label != ab93Labels[ab93Want] || !ab93Reflect.DeepEqual(ab93Raw, ab93Bytes) {
					ab93T.Fatal("status isolation/alternative label")
				}
			}
		}
	}
	if ab93Count != 38 || ab93ReservedCount != 20 {
		ab93T.Fatalf("status/reserved coverage %d/%d", ab93Count, ab93ReservedCount)
	}
}

func TestAb93PWMAndOpaque(ab93T *ab93Testing.T) {
	for _, ab93Field := range ab93Catalog(ab93T).Ab93Fields {
		if ab93Field.Ab93Key.Ab93Field == "PWM" {
			ab93Bytes := make([]byte, 20)
			ab93Bytes[15] = 255
			_, ab93Value, _, _ := ab93Evaluate(ab93T, ab93Field, ab93Bytes)
			if ab93Value != 127.5 {
				ab93T.Fatal("PWM clipped")
			}
		}
		if *ab93Field.Ab93Key.Ab93Operation == 7 && (ab93Field.Ab93Key.Ab93Field == "opaque" || ab93Field.Ab93Key.Ab93Field == "S12" || ab93Field.Ab93Key.Ab93Field == "S13") {
			ab93Bytes := make([]byte, 20)
			ab93Bytes[11] = 255
			ab93Bytes[12] = 16
			ab93Bytes[14] = 240
			ab93Bytes[15] = 255
			_, ab93Value, _, _ := ab93Evaluate(ab93T, ab93Field, ab93Bytes)
			ab93Want := map[string]float64{"opaque": 255, "S12": 1, "S13": -1}[ab93Field.Ab93Key.Ab93Field]
			if ab93Value != ab93Want {
				ab93T.Fatal("OP07 intermediate byte conflated with sensor")
			}
		}
	}
}

func ab93CopyKey(ab93Key Ab93SensorKey) Ab93SensorKey {
	ab93Data, _ := ab93JSON.Marshal(ab93Key)
	var ab93Copy Ab93SensorKey
	if ab93Error := ab93JSON.Unmarshal(ab93Data, &ab93Copy); ab93Error != nil {
		panic(ab93Error)
	}
	return ab93Copy
}

func TestAb93ExactIdentityAxes(ab93T *ab93Testing.T) {
	ab93Original := ab93Catalog(ab93T).Ab93Fields[0]
	for _, ab93B524 := range []bool{false, true} {
		ab93Control := ab93Original
		ab93Control.Ab93Key = ab93CopyKey(ab93Original.Ab93Key)
		if ab93B524 {
			ab93Control.Ab93Key.Ab93Secondary = ab93Ptr(uint16(36))
			ab93Control.Ab93Key.Ab93Operation = nil
			ab93Control.Ab93Key.Ab93Opcode = ab93Ptr(uint16(0))
			ab93Control.Ab93Key.Ab93Group = ab93Ptr(uint16(0))
			ab93Control.Ab93Key.Ab93Instance = ab93Ptr(uint16(0))
			ab93Control.Ab93Key.Ab93Register = ab93Ptr(uint16(0))
		}
		ab93Fields := []Ab93SensorField{ab93Control}
		ab93Equal := ab93CopyKey(ab93Control.Ab93Key)
		if _, ab93Found, ab93Error := ab93FindSensorField(ab93Fields, ab93Equal); !ab93Found || ab93Error != nil {
			ab93T.Fatal("valid control failed", ab93Error)
		}
		for _, ab93Axis := range []string{"Ab93Family", "Ab93Variant", "Ab93Primary", "Ab93Secondary", "Ab93Operation", "Ab93Direction", "Ab93Field", "Ab93Opcode", "Ab93Group", "Ab93Instance", "Ab93Register"} {
			ab93T.Run(ab93Fmt.Sprintf("ab93_%t_%s", ab93B524, ab93Axis), func(ab93T *ab93Testing.T) {
				ab93Key := ab93CopyKey(ab93Control.Ab93Key)
				ab93Value := ab93Reflect.ValueOf(&ab93Key).Elem().FieldByName(ab93Axis)
				if ab93Value.Kind() == ab93Reflect.String {
					ab93Value.SetString("ab93-other")
					if _, ab93Found, ab93Error := ab93FindSensorField(ab93Fields, ab93Key); ab93Found || ab93Error != nil {
						ab93T.Fatal("ignored string axis", ab93Error)
					}
					return
				}
				// Each nullable coordinate is tested absent, explicitly zero and nonzero.
				// The current control must match; every other value must miss.
				ab93Candidates := []ab93Reflect.Value{ab93Reflect.Zero(ab93Value.Type())}
				if ab93Value.Type().Elem().Kind() == ab93Reflect.String {
					ab93Candidates = append(ab93Candidates, ab93Reflect.ValueOf(ab93Ptr("")), ab93Reflect.ValueOf(ab93Ptr("response")))
				} else {
					ab93Candidates = append(ab93Candidates, ab93Reflect.ValueOf(ab93Ptr(uint16(0))), ab93Reflect.ValueOf(ab93Ptr(uint16(123))))
				}
				for _, ab93Candidate := range ab93Candidates {
					ab93Value.Set(ab93Candidate)
					ab93Want := ab93Reflect.DeepEqual(ab93Key, ab93Control.Ab93Key)
					_, ab93Found, ab93Error := ab93FindSensorField(ab93Fields, ab93Key)
					if ab93Found != ab93Want || ab93Error != nil {
						ab93T.Fatal("nullable exact identity failed", ab93Error)
					}
					if !ab93B524 {
						_, ab93Found, ab93Error = Ab93LookupB523SensorField(ab93Key)
						if ab93Found != ab93Want || ab93Error != nil {
							ab93T.Fatal("production nullable identity failed", ab93Error)
						}
					}
				}
			})
		}
	}
}

func ab93MutateField(ab93Field *Ab93SensorField) {
	ab93Field.Ab93Key.Ab93Family = "ab93-mutated"
	*ab93Field.Ab93Key.Ab93Primary = 0
	*ab93Field.Ab93Key.Ab93Secondary = 0
	*ab93Field.Ab93Key.Ab93Operation = 0
	*ab93Field.Ab93Position = 0
	*ab93Field.Ab93Width = 0
	if ab93Field.Ab93BitOffset != nil {
		*ab93Field.Ab93BitOffset = 99
	}
	if ab93Field.Ab93Default != nil {
		*ab93Field.Ab93Default = 99
	}
	if ab93Field.Ab93DivideBy != nil {
		*ab93Field.Ab93DivideBy = 99
	}
	if ab93Field.Ab93Unit != nil {
		*ab93Field.Ab93Unit = "ab93-mutated"
	}
	if ab93Field.Ab93ReplacementRaw != nil {
		*ab93Field.Ab93ReplacementRaw = 0
	}
	for ab93Index := range ab93Field.Ab93Bounds {
		ab93Field.Ab93Bounds[ab93Index] = 99
	}
	ab93Field.Ab93Enumeration["0"] = "ab93-mutated"
	for ab93Index := range ab93Field.Ab93Branches {
		ab93Branch := &ab93Field.Ab93Branches[ab93Index]
		ab93Branch.Ab93Name = "ab93-mutated"
		*ab93Branch.Ab93Default = 99
		*ab93Branch.Ab93DivideBy = 99
		if ab93Branch.Ab93Unit != nil {
			*ab93Branch.Ab93Unit = "ab93-mutated"
		}
		ab93Branch.Ab93Bounds[0] = 99
	}
}

func TestAb93DetachedAccessors(ab93T *ab93Testing.T) {
	ab93Original := ab93Catalog(ab93T)
	ab93Changed := ab93Catalog(ab93T)
	ab93Changed.Ab93Limitations[0] = "ab93-mutated"
	for ab93Index := range ab93Changed.Ab93Fields {
		ab93MutateField(&ab93Changed.Ab93Fields[ab93Index])
	}
	if !ab93Reflect.DeepEqual(ab93Original, ab93Catalog(ab93T)) {
		ab93T.Fatal("catalog shares backing storage")
	}
	for _, ab93OriginalField := range ab93Original.Ab93Fields {
		ab93ChangedField, ab93Found, ab93Error := Ab93LookupB523SensorField(ab93OriginalField.Ab93Key)
		if !ab93Found || ab93Error != nil {
			ab93T.Fatal(ab93Error)
		}
		ab93MutateField(&ab93ChangedField)
		ab93Fresh, ab93Found, ab93Error := Ab93LookupB523SensorField(ab93OriginalField.Ab93Key)
		if !ab93Found || ab93Error != nil || !ab93Reflect.DeepEqual(ab93Fresh, ab93OriginalField) {
			ab93T.Fatal("lookup shares backing storage", ab93Error)
		}
		ab93ChangedField, ab93Found, ab93Error = ab93FindSensorField(ab93Original.Ab93Fields, ab93OriginalField.Ab93Key)
		if !ab93Found || ab93Error != nil {
			ab93T.Fatal(ab93Error)
		}
		ab93MutateField(&ab93ChangedField)
		ab93Fresh, ab93Found, ab93Error = ab93FindSensorField(ab93Original.Ab93Fields, ab93OriginalField.Ab93Key)
		if !ab93Found || ab93Error != nil || !ab93Reflect.DeepEqual(ab93Fresh, ab93OriginalField) {
			ab93T.Fatal("find aliases input", ab93Error)
		}
	}
	if !ab93Reflect.DeepEqual(ab93Original, ab93Catalog(ab93T)) {
		ab93T.Fatal("nested lookup mutation reached catalog")
	}
}
