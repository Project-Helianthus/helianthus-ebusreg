package ebus_standard_catalog

import (
	ab83context "context"
	ab83reflect "reflect"
	ab83testing "testing"
)

func ab83Ptr[ab83T any](ab83Value ab83T) *ab83T { return &ab83Value }

func ab83TestKey(ab83Secondary uint8, ab83Variant string) Ab83Key {
	ab83Direction := "request"
	if ab83Variant == "standard_response" {
		ab83Direction = "response"
	}
	return Ab83Key{Ab83Family: "standard_03", Ab83Variant: ab83Variant, Ab83Primary: ab83Ptr(uint16(3)), Ab83Secondary: ab83Ptr(ab83Secondary), Ab83Direction: ab83Direction}
}

func TestAb83ProductionMetadata(ab83T *ab83testing.T) {
	ab83Bodies := Ab83CounterBodies()
	if len(ab83Bodies) != 12 {
		ab83T.Fatalf("body count = %d, want 12", len(ab83Bodies))
	}
	for _, ab83Variant := range []string{"candidate_request", "standard_response"} {
		for _, ab83Case := range []struct {
			ab83Secondary uint8
			ab83Names     []string
			ab83Maxima    []int
		}{
			{4, []string{"start_count_1s", "start_count_100s", "start_count_10000s"}, []int{99, 99, 99}},
			{5, []string{"minutes", "hours_100s", "hours_10000s"}, []int{59, 59, 59}},
			{6, []string{"minutes", "hours_1s", "hours_100s", "hours_10000s"}, []int{59, 99, 99, 99}},
			{7, nil, nil}, {8, nil, nil}, {16, nil, nil},
		} {
			ab83Key := ab83TestKey(ab83Case.ab83Secondary, ab83Variant)
			ab83Body, ab83OK := Ab83LookupBody(ab83Key)
			if !ab83OK {
				ab83T.Fatalf("missing body: %+v", ab83Key)
			}
			ab83Expected := Ab83Body{Ab83Key: ab83Key, Ab83BodyState: "candidate", Ab83ClaimGrade: "Hypothesis", Ab83PhysicalDirection: "Unknown", Ab83RawFallback: true}
			if len(ab83Case.ab83Names) == 0 {
				ab83Expected.Ab83BodyState = "missing"
				ab83Expected.Ab83ClaimGrade = "Unknown"
			}
			for ab83Position, ab83Name := range ab83Case.ab83Names {
				ab83ExpectedField := Ab83Field{Ab83Identity: ab83Name, Ab83LogicalPosition: ab83Ptr(ab83Position), Ab83WidthBits: ab83Ptr(8), Ab83Signedness: ab83Ptr("unsigned"), Ab83Conversion: ab83Ptr("identity"), Ab83Unit: ab83Ptr("raw"), Ab83Minimum: ab83Ptr(0), Ab83Maximum: ab83Ptr(ab83Case.ab83Maxima[ab83Position]), Ab83ClaimGrade: "Hypothesis", Ab83RawFallback: true}
				if ab83Variant == "candidate_request" {
					ab83ExpectedField.Ab83Encoding = ab83Ptr("raw")
					if ab83Case.ab83Secondary == 4 {
						ab83ExpectedField.Ab83Encoding = ab83Ptr("raw_decimal")
					}
				}
				if ab83Case.ab83Secondary == 4 {
					ab83ExpectedField.Ab83Default = ab83Ptr(0)
				}
				ab83Expected.Ab83Fields = append(ab83Expected.Ab83Fields, ab83ExpectedField)
				ab83Key.Ab83Field = ab83Name
				ab83Field, ab83Found := Ab83LookupField(ab83Key)
				if !ab83Found || !ab83reflect.DeepEqual(ab83Field, ab83ExpectedField) {
					ab83T.Fatalf("field %d/%s/%s: got %+v, want %+v", ab83Case.ab83Secondary, ab83Variant, ab83Name, ab83Field, ab83ExpectedField)
				}
			}
			if !ab83reflect.DeepEqual(ab83Body, ab83Expected) {
				ab83T.Fatalf("body %d/%s:\ngot %+v\nwant %+v", ab83Case.ab83Secondary, ab83Variant, ab83Body, ab83Expected)
			}
		}
	}
}

func ab83SyntheticValid(ab83Body Ab83Body, ab83Bytes []byte) bool {
	if ab83Body.Ab83BodyState != "candidate" || len(ab83Body.Ab83Fields) == 0 || len(ab83Bytes) != len(ab83Body.Ab83Fields) {
		return false
	}
	for ab83Index, ab83Field := range ab83Body.Ab83Fields {
		if ab83Field.Ab83LogicalPosition == nil || *ab83Field.Ab83LogicalPosition != ab83Index || ab83Field.Ab83WidthBits == nil || *ab83Field.Ab83WidthBits != 8 || ab83Field.Ab83Signedness == nil || *ab83Field.Ab83Signedness != "unsigned" || ab83Field.Ab83Conversion == nil || *ab83Field.Ab83Conversion != "identity" || ab83Field.Ab83Encoding == nil {
			return false
		}
		if *ab83Field.Ab83Encoding != "raw" && *ab83Field.Ab83Encoding != "raw_decimal" {
			return false
		}
		ab83Value := int(ab83Bytes[ab83Index])
		if ab83Field.Ab83Minimum == nil || ab83Field.Ab83Maximum == nil || ab83Value < *ab83Field.Ab83Minimum || ab83Value > *ab83Field.Ab83Maximum {
			return false
		}
	}
	return true
}

func TestAb83SyntheticCases(ab83T *ab83testing.T) {
	ab83Four, _ := Ab83LookupBody(ab83TestKey(4, "candidate_request"))
	ab83Bytes := []byte{12, 34, 56}
	if !ab83SyntheticValid(ab83Four, ab83Bytes) {
		ab83T.Fatal("valid decimal chunks rejected")
	}
	ab83Total, ab83Multiplier := 0, 1
	for _, ab83Field := range ab83Four.Ab83Fields {
		ab83Total += int(ab83Bytes[*ab83Field.Ab83LogicalPosition]) * ab83Multiplier
		ab83Multiplier *= 100
	}
	if ab83Total != 563412 {
		ab83T.Fatalf("decimal total=%d", ab83Total)
	}
	if !ab83SyntheticValid(ab83Four, []byte{0x63, 0, 0}) || int(byte(0x63)) != 99 {
		ab83T.Fatal("decimal 99 must not be treated as BCD 63")
	}
	if ab83SyntheticValid(ab83Four, []byte{100, 0, 0}) {
		ab83T.Fatal("decimal bound ignored")
	}
	ab83Five, _ := Ab83LookupBody(ab83TestKey(5, "candidate_request"))
	if len(ab83Five.Ab83Fields) != 3 || !ab83SyntheticValid(ab83Five, []byte{59, 59, 59}) || ab83SyntheticValid(ab83Five, []byte{0, 60, 0}) || ab83SyntheticValid(ab83Five, []byte{0, 0, 0, 0}) {
		ab83T.Fatal("03/05 domain or missing hours-one distinction lost")
	}
	ab83Six, _ := Ab83LookupBody(ab83TestKey(6, "candidate_request"))
	if len(ab83Six.Ab83Fields) != 4 || !ab83SyntheticValid(ab83Six, []byte{59, 99, 99, 99}) || ab83SyntheticValid(ab83Six, []byte{59, 99, 99}) {
		ab83T.Fatal("03/06 completeness lost")
	}
	ab83Response, _ := Ab83LookupBody(ab83TestKey(6, "standard_response"))
	if ab83Response.Ab83Fields[0].Ab83Encoding != nil || ab83Six.Ab83Fields[0].Ab83Encoding == nil || ab83SyntheticValid(ab83Response, []byte{59, 99, 99, 99}) {
		ab83T.Fatal("unspecified encoding became explicit raw")
	}
	if ab83Six.Ab83Fields[0].Ab83NativePosition != nil || ab83Six.Ab83Fields[0].Ab83LogicalPosition == nil || *ab83Six.Ab83Fields[0].Ab83LogicalPosition != 0 {
		ab83T.Fatal("unknown native position became zero")
	}
	for _, ab83Secondary := range []uint8{7, 8, 16} {
		for _, ab83Variant := range []string{"candidate_request", "standard_response"} {
			ab83Unknown, _ := Ab83LookupBody(ab83TestKey(ab83Secondary, ab83Variant))
			if ab83Unknown.Ab83BodyState != "missing" || ab83Unknown.Ab83ClaimGrade != "Unknown" || ab83Unknown.Ab83Callable || ab83Unknown.Ab83Fields != nil || ab83SyntheticValid(ab83Unknown, nil) {
				ab83T.Fatalf("unknown %d advertised as body", ab83Secondary)
			}
		}
	}
}

func TestAb83ExactIdentityAxes(ab83T *ab83testing.T) {
	ab83Control := ab83TestKey(4, "candidate_request")
	ab83Control.Ab83Field = "start_count_1s"
	if _, ab83OK := Ab83LookupField(ab83Control); !ab83OK {
		ab83T.Fatal("valid control missing")
	}
	ab83Mutations := []struct {
		ab83Name   string
		ab83Change func(*Ab83Key)
	}{
		{"family", func(ab83Key *Ab83Key) { ab83Key.Ab83Family = "other" }},
		{"variant", func(ab83Key *Ab83Key) { ab83Key.Ab83Variant = "other" }},
		{"primary", func(ab83Key *Ab83Key) { ab83Key.Ab83Primary = ab83Ptr(uint16(4)) }},
		{"primary_absent", func(ab83Key *Ab83Key) { ab83Key.Ab83Primary = nil }},
		{"secondary", func(ab83Key *Ab83Key) { ab83Key.Ab83Secondary = ab83Ptr(uint8(6)) }},
		{"secondary_absent", func(ab83Key *Ab83Key) { ab83Key.Ab83Secondary = nil }},
		{"operation_zero", func(ab83Key *Ab83Key) { ab83Key.Ab83Operation = ab83Ptr(uint8(0)) }},
		{"direction", func(ab83Key *Ab83Key) { ab83Key.Ab83Direction = "response" }},
		{"field", func(ab83Key *Ab83Key) { ab83Key.Ab83Field = "other" }},
		{"opcode_zero", func(ab83Key *Ab83Key) { ab83Key.Ab83Opcode = ab83Ptr(uint8(0)) }},
		{"group_zero", func(ab83Key *Ab83Key) { ab83Key.Ab83Group = ab83Ptr(uint8(0)) }},
		{"instance_zero", func(ab83Key *Ab83Key) { ab83Key.Ab83Instance = ab83Ptr(uint8(0)) }},
		{"register_zero", func(ab83Key *Ab83Key) { ab83Key.Ab83Register = ab83Ptr(uint16(0)) }},
	}
	for _, ab83Mutation := range ab83Mutations {
		ab83T.Run("ab83_"+ab83Mutation.ab83Name, func(ab83T *ab83testing.T) {
			ab83Changed := ab83Control
			ab83Mutation.ab83Change(&ab83Changed)
			if ab83EqualKey(ab83Control, ab83Changed) {
				ab83T.Fatal("changed axis ignored")
			}
			if _, ab83OK := Ab83LookupField(ab83Changed); ab83OK {
				ab83T.Fatal("non-exact field match")
			}
			ab83Changed.Ab83Field = ""
			if _, ab83OK := Ab83LookupBody(ab83Changed); ab83OK && ab83Mutation.ab83Name != "field" && ab83Mutation.ab83Name != "secondary" {
				ab83T.Fatal("non-exact body match")
			}
		})
	}
	ab83BodyKey := ab83Control
	if _, ab83OK := Ab83LookupBody(ab83BodyKey); ab83OK {
		ab83T.Fatal("field-qualified key matched body")
	}
	// Synthetic B524 keys exercise matching only; no B524 declaration is published.
	ab83B524 := Ab83Key{Ab83Family: "B524", Ab83Variant: "synthetic", Ab83Primary: ab83Ptr(uint16(0xb524)), Ab83Secondary: ab83Ptr(uint8(0)), Ab83Operation: ab83Ptr(uint8(0)), Ab83Direction: "request", Ab83Field: "synthetic", Ab83Opcode: ab83Ptr(uint8(0)), Ab83Group: ab83Ptr(uint8(0)), Ab83Instance: ab83Ptr(uint8(0)), Ab83Register: ab83Ptr(uint16(0))}
	ab83Type := ab83reflect.TypeOf(ab83B524)
	for ab83Index := 0; ab83Index < ab83Type.NumField(); ab83Index++ {
		ab83Name := ab83Type.Field(ab83Index).Name
		for _, ab83Mode := range []string{"value", "absent"} {
			ab83T.Run("ab83_B524_"+ab83Name+"_"+ab83Mode, func(ab83T *ab83testing.T) {
				ab83Changed := ab83B524
				ab83Value := ab83reflect.ValueOf(&ab83Changed).Elem().Field(ab83Index)
				if ab83Value.Kind() == ab83reflect.String {
					ab83Value.SetString("")
				} else if ab83Mode == "absent" {
					ab83Value.Set(ab83reflect.Zero(ab83Value.Type()))
				} else {
					ab83Pointer := ab83reflect.New(ab83Value.Type().Elem())
					ab83Pointer.Elem().SetUint(ab83Value.Elem().Uint() + 1)
					ab83Value.Set(ab83Pointer)
				}
				if !ab83EqualKey(ab83B524, ab83B524) || ab83EqualKey(ab83B524, ab83Changed) || ab83EqualKey(ab83Changed, ab83B524) {
					ab83T.Fatal("synthetic exact axis equality failed")
				}
				if _, ab83OK := Ab83LookupField(ab83B524); ab83OK {
					ab83T.Fatal("synthetic B524 key became published support")
				}
			})
		}
	}
	for ab83Index := 0; ab83Index < ab83Type.NumField(); ab83Index++ {
		if ab83Type.Field(ab83Index).Type.Kind() != ab83reflect.Pointer {
			continue
		}
		ab83Zero := ab83B524
		ab83Value := ab83reflect.ValueOf(&ab83Zero).Elem().Field(ab83Index)
		ab83Value.Set(ab83reflect.New(ab83Value.Type().Elem()))
		ab83Absent := ab83Zero
		ab83reflect.ValueOf(&ab83Absent).Elem().Field(ab83Index).Set(ab83reflect.Zero(ab83Value.Type()))
		if !ab83EqualKey(ab83Zero, ab83Zero) || ab83EqualKey(ab83Zero, ab83Absent) || ab83EqualKey(ab83Absent, ab83Zero) {
			ab83T.Fatalf("absent versus zero collapsed for %s", ab83Type.Field(ab83Index).Name)
		}
	}

}

func ab83MutateNested(ab83Value ab83reflect.Value) {
	switch ab83Value.Kind() {
	case ab83reflect.Pointer:
		if !ab83Value.IsNil() {
			ab83MutateNested(ab83Value.Elem())
		}
	case ab83reflect.Struct:
		for ab83Index := 0; ab83Index < ab83Value.NumField(); ab83Index++ {
			ab83MutateNested(ab83Value.Field(ab83Index))
		}
	case ab83reflect.Slice:
		for ab83Index := 0; ab83Index < ab83Value.Len(); ab83Index++ {
			ab83MutateNested(ab83Value.Index(ab83Index))
		}
		if ab83Value.CanSet() {
			ab83Value.Set(ab83reflect.MakeSlice(ab83Value.Type(), 1, 1))
		}
	case ab83reflect.Map:
		if ab83Value.CanSet() {
			ab83Value.Set(ab83reflect.MakeMap(ab83Value.Type()))
			ab83Value.SetMapIndex(ab83reflect.ValueOf("synthetic"), ab83reflect.ValueOf("changed"))
		}
	case ab83reflect.String:
		if ab83Value.CanSet() {
			ab83Value.SetString("changed")
		}
	case ab83reflect.Int:
		if ab83Value.CanSet() {
			ab83Value.SetInt(17)
		}
	case ab83reflect.Uint8, ab83reflect.Uint16:
		if ab83Value.CanSet() {
			ab83Value.SetUint(17)
		}
	case ab83reflect.Bool:
		if ab83Value.CanSet() {
			ab83Value.SetBool(!ab83Value.Bool())
		}
	}
}

func TestAb83DetachedResults(ab83T *ab83testing.T) {
	ab83Expected := Ab83CounterBodies()
	ab83Changed := Ab83CounterBodies()
	ab83MutateNested(ab83reflect.ValueOf(&ab83Changed).Elem())
	if !ab83reflect.DeepEqual(Ab83CounterBodies(), ab83Expected) {
		ab83T.Fatal("accessor shared storage")
	}
	for _, ab83Body := range ab83Expected {
		ab83ChangedBody, ab83OK := Ab83LookupBody(ab83Body.Ab83Key)
		if !ab83OK {
			ab83T.Fatal("body missing")
		}
		ab83MutateNested(ab83reflect.ValueOf(&ab83ChangedBody).Elem())
		ab83Again, _ := Ab83LookupBody(ab83Body.Ab83Key)
		if !ab83reflect.DeepEqual(ab83Again, ab83Body) {
			ab83T.Fatal("body lookup shared storage")
		}
		for _, ab83Field := range ab83Body.Ab83Fields {
			ab83Key := ab83Body.Ab83Key
			ab83Key.Ab83Field = ab83Field.Ab83Identity
			ab83ChangedField, ab83Found := Ab83LookupField(ab83Key)
			if !ab83Found {
				ab83T.Fatal("field missing")
			}
			ab83MutateNested(ab83reflect.ValueOf(&ab83ChangedField).Elem())
			ab83AgainField, _ := Ab83LookupField(ab83Key)
			if !ab83reflect.DeepEqual(ab83AgainField, ab83Field) {
				ab83T.Fatal("field lookup shared storage")
			}
		}
	}
}

func TestAb83MetadataCannotEnableStub(ab83T *ab83testing.T) {
	ab83Provider, ab83Err := NewProvider(MustEmbeddedCatalog(), true)
	if ab83Err != nil {
		ab83T.Fatal(ab83Err)
	}
	ab83IDs := []string{"ebus_standard.service_data.start_counts", "ebus_standard.service_data.operating_time_l1", "ebus_standard.service_data.operating_time_l2", "ebus_standard.service_data.operating_time_l3", "ebus_standard.service_data.fuel_quantity_counter", "ebus_standard.service_data.meter_reading"}
	for _, ab83ID := range ab83IDs {
		ab83Before, ab83Err := ab83Provider.Invoke(ab83context.Background(), ab83ID, nil, CallerContext(0))
		if ab83Err != nil {
			ab83T.Fatal(ab83Err)
		}
		for _, ab83Body := range Ab83CounterBodies() {
			ab83Selected, ab83OK := Ab83LookupBody(ab83Body.Ab83Key)
			if !ab83OK || ab83Selected.Ab83Callable {
				ab83T.Fatal("metadata advertised callable body")
			}
			ab83Selected.Ab83Callable = true
		}
		ab83After, ab83Err := ab83Provider.Invoke(ab83context.Background(), ab83ID, nil, CallerContext(0))
		if ab83Err != nil || ab83After["executed"] != false || !ab83reflect.DeepEqual(ab83Before, ab83After) {
			ab83T.Fatalf("stub behavior changed: %v %v", ab83After, ab83Err)
		}
	}
	ab83Identification, ab83Err := ab83Provider.Identification(ab83context.Background())
	if ab83Err != nil || ab83Identification.Valid {
		ab83T.Fatal("metadata qualified identification")
	}
}
