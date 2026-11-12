package ebusstandard

import (
	standard03CounterBodyContext "context"
	standard03CounterBodyReflect "reflect"
	standard03CounterBodyTesting "testing"
)

func standard03CounterBodyPtr[standard03CounterBodyT any](standard03CounterBodyValue standard03CounterBodyT) *standard03CounterBodyT {
	return &standard03CounterBodyValue
}

func standard03CounterBodyTestKey(standard03CounterBodySecondary uint8, standard03CounterBodyVariant string) Standard03CounterBodyKey {
	standard03CounterBodyDirection := "request"
	if standard03CounterBodyVariant == "standard_response" {
		standard03CounterBodyDirection = "response"
	}
	return Standard03CounterBodyKey{Standard03CounterBodyFamily: "standard_03", Standard03CounterBodyVariant: standard03CounterBodyVariant, Standard03CounterBodyPrimary: standard03CounterBodyPtr(uint16(3)), Standard03CounterBodySecondary: standard03CounterBodyPtr(standard03CounterBodySecondary), Standard03CounterBodyDirection: standard03CounterBodyDirection}
}

func standard03CounterBodyFive(standard03CounterBodyVariant string) struct {
	standard03CounterBodySecondary uint8
	standard03CounterBodyNames     []string
	standard03CounterBodyMaxima    []int
} {
	if standard03CounterBodyVariant == "standard_response" {
		return struct {
			standard03CounterBodySecondary uint8
			standard03CounterBodyNames     []string
			standard03CounterBodyMaxima    []int
		}{5, []string{"minutes", "hours_1s", "hours_100s", "hours_10000s"}, []int{59, 99, 99, 99}}
	}
	return struct {
		standard03CounterBodySecondary uint8
		standard03CounterBodyNames     []string
		standard03CounterBodyMaxima    []int
	}{5, []string{"minutes", "hours_100s", "hours_10000s"}, []int{59, 59, 59}}
}

func TestAb83ProductionMetadata(standard03CounterBodyT *standard03CounterBodyTesting.T) {
	standard03CounterBodyBodies := Standard03CounterBodies()
	if len(standard03CounterBodyBodies) != 12 {
		standard03CounterBodyT.Fatalf("body count = %d, want 12", len(standard03CounterBodyBodies))
	}
	for _, standard03CounterBodyVariant := range []string{"candidate_request", "standard_response"} {
		for _, standard03CounterBodyCase := range []struct {
			standard03CounterBodySecondary uint8
			standard03CounterBodyNames     []string
			standard03CounterBodyMaxima    []int
		}{
			{4, []string{"start_count_1s", "start_count_100s", "start_count_10000s"}, []int{99, 99, 99}},
			standard03CounterBodyFive(standard03CounterBodyVariant),
			{6, []string{"minutes", "hours_1s", "hours_100s", "hours_10000s"}, []int{59, 99, 99, 99}},
			{7, nil, nil}, {8, nil, nil}, {16, nil, nil},
		} {
			standard03CounterBodyKey := standard03CounterBodyTestKey(standard03CounterBodyCase.standard03CounterBodySecondary, standard03CounterBodyVariant)
			standard03CounterBodyBody, standard03CounterBodyOK := Standard03CounterLookupBody(standard03CounterBodyKey)
			if !standard03CounterBodyOK {
				standard03CounterBodyT.Fatalf("missing body: %+v", standard03CounterBodyKey)
			}
			standard03CounterBodyExpected := Standard03CounterBody{Standard03CounterBodyKey: standard03CounterBodyKey, Standard03CounterBodyState: "candidate", Standard03CounterBodyClaimGrade: "Hypothesis", Standard03CounterBodyPhysicalDirection: "Unknown", Standard03CounterBodyRawFallback: true}
			if len(standard03CounterBodyCase.standard03CounterBodyNames) == 0 {
				standard03CounterBodyExpected.Standard03CounterBodyState = "missing"
				standard03CounterBodyExpected.Standard03CounterBodyClaimGrade = "Unknown"
			}
			for standard03CounterBodyPosition, standard03CounterBodyName := range standard03CounterBodyCase.standard03CounterBodyNames {
				standard03CounterBodyExpectedField := Standard03CounterBodyField{Standard03CounterBodyIdentity: standard03CounterBodyName, Standard03CounterBodyLogicalPosition: standard03CounterBodyPtr(standard03CounterBodyPosition), Standard03CounterBodyWidthBits: standard03CounterBodyPtr(8), Standard03CounterBodySignedness: standard03CounterBodyPtr("unsigned"), Standard03CounterBodyConversion: standard03CounterBodyPtr("identity"), Standard03CounterBodyUnit: standard03CounterBodyPtr("raw"), Standard03CounterBodyMinimum: standard03CounterBodyPtr(0), Standard03CounterBodyMaximum: standard03CounterBodyPtr(standard03CounterBodyCase.standard03CounterBodyMaxima[standard03CounterBodyPosition]), Standard03CounterBodyClaimGrade: "Hypothesis", Standard03CounterBodyRawFallback: true}
				if standard03CounterBodyVariant == "candidate_request" {
					standard03CounterBodyExpectedField.Standard03CounterBodyEncoding = standard03CounterBodyPtr("raw")
					if standard03CounterBodyCase.standard03CounterBodySecondary == 4 {
						standard03CounterBodyExpectedField.Standard03CounterBodyEncoding = standard03CounterBodyPtr("raw_decimal")
					}
				}
				if standard03CounterBodyCase.standard03CounterBodySecondary == 4 {
					standard03CounterBodyExpectedField.Standard03CounterBodyDefault = standard03CounterBodyPtr(0)
				}
				standard03CounterBodyExpected.Standard03CounterBodyFields = append(standard03CounterBodyExpected.Standard03CounterBodyFields, standard03CounterBodyExpectedField)
				standard03CounterBodyKey.Standard03CounterBodyField = standard03CounterBodyName
				standard03CounterBodyField, standard03CounterBodyFound := Standard03CounterBodyLookupField(standard03CounterBodyKey)
				if !standard03CounterBodyFound || !standard03CounterBodyReflect.DeepEqual(standard03CounterBodyField, standard03CounterBodyExpectedField) {
					standard03CounterBodyT.Fatalf("field %d/%s/%s: got %+v, want %+v", standard03CounterBodyCase.standard03CounterBodySecondary, standard03CounterBodyVariant, standard03CounterBodyName, standard03CounterBodyField, standard03CounterBodyExpectedField)
				}
			}
			if !standard03CounterBodyReflect.DeepEqual(standard03CounterBodyBody, standard03CounterBodyExpected) {
				standard03CounterBodyT.Fatalf("body %d/%s:\ngot %+v\nwant %+v", standard03CounterBodyCase.standard03CounterBodySecondary, standard03CounterBodyVariant, standard03CounterBodyBody, standard03CounterBodyExpected)
			}
		}
	}
}

func standard03CounterBodySyntheticValid(standard03CounterBodyBody Standard03CounterBody, standard03CounterBodyBytes []byte) bool {
	if standard03CounterBodyBody.Standard03CounterBodyState != "candidate" || len(standard03CounterBodyBody.Standard03CounterBodyFields) == 0 || len(standard03CounterBodyBytes) != len(standard03CounterBodyBody.Standard03CounterBodyFields) {
		return false
	}
	for standard03CounterBodyIndex, standard03CounterBodyField := range standard03CounterBodyBody.Standard03CounterBodyFields {
		if standard03CounterBodyField.Standard03CounterBodyLogicalPosition == nil || *standard03CounterBodyField.Standard03CounterBodyLogicalPosition != standard03CounterBodyIndex || standard03CounterBodyField.Standard03CounterBodyWidthBits == nil || *standard03CounterBodyField.Standard03CounterBodyWidthBits != 8 || standard03CounterBodyField.Standard03CounterBodySignedness == nil || *standard03CounterBodyField.Standard03CounterBodySignedness != "unsigned" || standard03CounterBodyField.Standard03CounterBodyConversion == nil || *standard03CounterBodyField.Standard03CounterBodyConversion != "identity" || standard03CounterBodyField.Standard03CounterBodyEncoding == nil {
			return false
		}
		if *standard03CounterBodyField.Standard03CounterBodyEncoding != "raw" && *standard03CounterBodyField.Standard03CounterBodyEncoding != "raw_decimal" {
			return false
		}
		standard03CounterBodyValue := int(standard03CounterBodyBytes[standard03CounterBodyIndex])
		if standard03CounterBodyField.Standard03CounterBodyMinimum == nil || standard03CounterBodyField.Standard03CounterBodyMaximum == nil || standard03CounterBodyValue < *standard03CounterBodyField.Standard03CounterBodyMinimum || standard03CounterBodyValue > *standard03CounterBodyField.Standard03CounterBodyMaximum {
			return false
		}
	}
	return true
}

func TestAb83SyntheticCases(standard03CounterBodyT *standard03CounterBodyTesting.T) {
	standard03CounterBodyFour, _ := Standard03CounterLookupBody(standard03CounterBodyTestKey(4, "candidate_request"))
	standard03CounterBodyBytes := []byte{12, 34, 56}
	if !standard03CounterBodySyntheticValid(standard03CounterBodyFour, standard03CounterBodyBytes) {
		standard03CounterBodyT.Fatal("valid decimal chunks rejected")
	}
	standard03CounterBodyTotal, standard03CounterBodyMultiplier := 0, 1
	for _, standard03CounterBodyField := range standard03CounterBodyFour.Standard03CounterBodyFields {
		standard03CounterBodyTotal += int(standard03CounterBodyBytes[*standard03CounterBodyField.Standard03CounterBodyLogicalPosition]) * standard03CounterBodyMultiplier
		standard03CounterBodyMultiplier *= 100
	}
	if standard03CounterBodyTotal != 563412 {
		standard03CounterBodyT.Fatalf("decimal total=%d", standard03CounterBodyTotal)
	}
	if !standard03CounterBodySyntheticValid(standard03CounterBodyFour, []byte{0x63, 0, 0}) || int(byte(0x63)) != 99 {
		standard03CounterBodyT.Fatal("decimal 99 must not be treated as BCD 63")
	}
	if standard03CounterBodySyntheticValid(standard03CounterBodyFour, []byte{100, 0, 0}) {
		standard03CounterBodyT.Fatal("decimal bound ignored")
	}
	standard03CounterBodyFive, _ := Standard03CounterLookupBody(standard03CounterBodyTestKey(5, "candidate_request"))
	if len(standard03CounterBodyFive.Standard03CounterBodyFields) != 3 || !standard03CounterBodySyntheticValid(standard03CounterBodyFive, []byte{59, 59, 59}) || standard03CounterBodySyntheticValid(standard03CounterBodyFive, []byte{0, 60, 0}) || standard03CounterBodySyntheticValid(standard03CounterBodyFive, []byte{0, 0, 0, 0}) {
		standard03CounterBodyT.Fatal("03/05 domain or missing hours-one distinction lost")
	}
	standard03CounterBodySix, _ := Standard03CounterLookupBody(standard03CounterBodyTestKey(6, "candidate_request"))
	if len(standard03CounterBodySix.Standard03CounterBodyFields) != 4 || !standard03CounterBodySyntheticValid(standard03CounterBodySix, []byte{59, 99, 99, 99}) || standard03CounterBodySyntheticValid(standard03CounterBodySix, []byte{59, 99, 99}) {
		standard03CounterBodyT.Fatal("03/06 completeness lost")
	}
	standard03CounterBodyResponse, _ := Standard03CounterLookupBody(standard03CounterBodyTestKey(6, "standard_response"))
	if standard03CounterBodyResponse.Standard03CounterBodyFields[0].Standard03CounterBodyEncoding != nil || standard03CounterBodySix.Standard03CounterBodyFields[0].Standard03CounterBodyEncoding == nil || standard03CounterBodySyntheticValid(standard03CounterBodyResponse, []byte{59, 99, 99, 99}) {
		standard03CounterBodyT.Fatal("unspecified encoding became explicit raw")
	}
	if standard03CounterBodySix.Standard03CounterBodyFields[0].Standard03CounterBodyNativePosition != nil || standard03CounterBodySix.Standard03CounterBodyFields[0].Standard03CounterBodyLogicalPosition == nil || *standard03CounterBodySix.Standard03CounterBodyFields[0].Standard03CounterBodyLogicalPosition != 0 {
		standard03CounterBodyT.Fatal("unknown native position became zero")
	}
	for _, standard03CounterBodySecondary := range []uint8{7, 8, 16} {
		for _, standard03CounterBodyVariant := range []string{"candidate_request", "standard_response"} {
			standard03CounterBodyUnknown, _ := Standard03CounterLookupBody(standard03CounterBodyTestKey(standard03CounterBodySecondary, standard03CounterBodyVariant))
			if standard03CounterBodyUnknown.Standard03CounterBodyState != "missing" || standard03CounterBodyUnknown.Standard03CounterBodyClaimGrade != "Unknown" || standard03CounterBodyUnknown.Standard03CounterBodyCallable || standard03CounterBodyUnknown.Standard03CounterBodyFields != nil || standard03CounterBodySyntheticValid(standard03CounterBodyUnknown, nil) {
				standard03CounterBodyT.Fatalf("unknown %d advertised as body", standard03CounterBodySecondary)
			}
		}
	}
}

func TestAb83ExactIdentityAxes(standard03CounterBodyT *standard03CounterBodyTesting.T) {
	standard03CounterBodyControl := standard03CounterBodyTestKey(4, "candidate_request")
	standard03CounterBodyControl.Standard03CounterBodyField = "start_count_1s"
	if _, standard03CounterBodyOK := Standard03CounterBodyLookupField(standard03CounterBodyControl); !standard03CounterBodyOK {
		standard03CounterBodyT.Fatal("valid control missing")
	}
	standard03CounterBodyMutations := []struct {
		standard03CounterBodyName   string
		standard03CounterBodyChange func(*Standard03CounterBodyKey)
	}{
		{"family", func(standard03CounterBodyKey *Standard03CounterBodyKey) {
			standard03CounterBodyKey.Standard03CounterBodyFamily = "other"
		}},
		{"variant", func(standard03CounterBodyKey *Standard03CounterBodyKey) {
			standard03CounterBodyKey.Standard03CounterBodyVariant = "other"
		}},
		{"primary", func(standard03CounterBodyKey *Standard03CounterBodyKey) {
			standard03CounterBodyKey.Standard03CounterBodyPrimary = standard03CounterBodyPtr(uint16(4))
		}},
		{"primary_absent", func(standard03CounterBodyKey *Standard03CounterBodyKey) {
			standard03CounterBodyKey.Standard03CounterBodyPrimary = nil
		}},
		{"secondary", func(standard03CounterBodyKey *Standard03CounterBodyKey) {
			standard03CounterBodyKey.Standard03CounterBodySecondary = standard03CounterBodyPtr(uint8(6))
		}},
		{"secondary_absent", func(standard03CounterBodyKey *Standard03CounterBodyKey) {
			standard03CounterBodyKey.Standard03CounterBodySecondary = nil
		}},
		{"operation_zero", func(standard03CounterBodyKey *Standard03CounterBodyKey) {
			standard03CounterBodyKey.Standard03CounterBodyOperation = standard03CounterBodyPtr(uint8(0))
		}},
		{"direction", func(standard03CounterBodyKey *Standard03CounterBodyKey) {
			standard03CounterBodyKey.Standard03CounterBodyDirection = "response"
		}},
		{"field", func(standard03CounterBodyKey *Standard03CounterBodyKey) {
			standard03CounterBodyKey.Standard03CounterBodyField = "other"
		}},
		{"opcode_zero", func(standard03CounterBodyKey *Standard03CounterBodyKey) {
			standard03CounterBodyKey.Standard03CounterBodyOpcode = standard03CounterBodyPtr(uint8(0))
		}},
		{"group_zero", func(standard03CounterBodyKey *Standard03CounterBodyKey) {
			standard03CounterBodyKey.Standard03CounterBodyGroup = standard03CounterBodyPtr(uint8(0))
		}},
		{"instance_zero", func(standard03CounterBodyKey *Standard03CounterBodyKey) {
			standard03CounterBodyKey.Standard03CounterBodyInstance = standard03CounterBodyPtr(uint8(0))
		}},
		{"register_zero", func(standard03CounterBodyKey *Standard03CounterBodyKey) {
			standard03CounterBodyKey.Standard03CounterBodyRegister = standard03CounterBodyPtr(uint16(0))
		}},
	}
	for _, standard03CounterBodyMutation := range standard03CounterBodyMutations {
		standard03CounterBodyT.Run("ab83_"+standard03CounterBodyMutation.standard03CounterBodyName, func(standard03CounterBodyT *standard03CounterBodyTesting.T) {
			standard03CounterBodyChanged := standard03CounterBodyControl
			standard03CounterBodyMutation.standard03CounterBodyChange(&standard03CounterBodyChanged)
			if standard03CounterBodyEqualKey(standard03CounterBodyControl, standard03CounterBodyChanged) {
				standard03CounterBodyT.Fatal("changed axis ignored")
			}
			if _, standard03CounterBodyOK := Standard03CounterBodyLookupField(standard03CounterBodyChanged); standard03CounterBodyOK {
				standard03CounterBodyT.Fatal("non-exact field match")
			}
			standard03CounterBodyChanged.Standard03CounterBodyField = ""
			if _, standard03CounterBodyOK := Standard03CounterLookupBody(standard03CounterBodyChanged); standard03CounterBodyOK && standard03CounterBodyMutation.standard03CounterBodyName != "field" && standard03CounterBodyMutation.standard03CounterBodyName != "secondary" {
				standard03CounterBodyT.Fatal("non-exact body match")
			}
		})
	}
	standard03CounterBodyBodyKey := standard03CounterBodyControl
	if _, standard03CounterBodyOK := Standard03CounterLookupBody(standard03CounterBodyBodyKey); standard03CounterBodyOK {
		standard03CounterBodyT.Fatal("field-qualified key matched body")
	}
	// Synthetic B524 keys exercise matching only; no B524 declaration is published.
	standard03CounterBodyB524 := Standard03CounterBodyKey{Standard03CounterBodyFamily: "B524", Standard03CounterBodyVariant: "synthetic", Standard03CounterBodyPrimary: standard03CounterBodyPtr(uint16(0xb524)), Standard03CounterBodySecondary: standard03CounterBodyPtr(uint8(0)), Standard03CounterBodyOperation: standard03CounterBodyPtr(uint8(0)), Standard03CounterBodyDirection: "request", Standard03CounterBodyField: "synthetic", Standard03CounterBodyOpcode: standard03CounterBodyPtr(uint8(0)), Standard03CounterBodyGroup: standard03CounterBodyPtr(uint8(0)), Standard03CounterBodyInstance: standard03CounterBodyPtr(uint8(0)), Standard03CounterBodyRegister: standard03CounterBodyPtr(uint16(0))}
	standard03CounterBodyType := standard03CounterBodyReflect.TypeOf(standard03CounterBodyB524)
	for standard03CounterBodyIndex := 0; standard03CounterBodyIndex < standard03CounterBodyType.NumField(); standard03CounterBodyIndex++ {
		standard03CounterBodyName := standard03CounterBodyType.Field(standard03CounterBodyIndex).Name
		for _, standard03CounterBodyMode := range []string{"value", "absent"} {
			standard03CounterBodyT.Run("ab83_B524_"+standard03CounterBodyName+"_"+standard03CounterBodyMode, func(standard03CounterBodyT *standard03CounterBodyTesting.T) {
				standard03CounterBodyChanged := standard03CounterBodyB524
				standard03CounterBodyValue := standard03CounterBodyReflect.ValueOf(&standard03CounterBodyChanged).Elem().Field(standard03CounterBodyIndex)
				if standard03CounterBodyValue.Kind() == standard03CounterBodyReflect.String {
					standard03CounterBodyValue.SetString("")
				} else if standard03CounterBodyMode == "absent" {
					standard03CounterBodyValue.Set(standard03CounterBodyReflect.Zero(standard03CounterBodyValue.Type()))
				} else {
					standard03CounterBodyPointer := standard03CounterBodyReflect.New(standard03CounterBodyValue.Type().Elem())
					standard03CounterBodyPointer.Elem().SetUint(standard03CounterBodyValue.Elem().Uint() + 1)
					standard03CounterBodyValue.Set(standard03CounterBodyPointer)
				}
				if !standard03CounterBodyEqualKey(standard03CounterBodyB524, standard03CounterBodyB524) || standard03CounterBodyEqualKey(standard03CounterBodyB524, standard03CounterBodyChanged) || standard03CounterBodyEqualKey(standard03CounterBodyChanged, standard03CounterBodyB524) {
					standard03CounterBodyT.Fatal("synthetic exact axis equality failed")
				}
				if _, standard03CounterBodyOK := Standard03CounterBodyLookupField(standard03CounterBodyB524); standard03CounterBodyOK {
					standard03CounterBodyT.Fatal("synthetic B524 key became published support")
				}
			})
		}
	}
	for standard03CounterBodyIndex := 0; standard03CounterBodyIndex < standard03CounterBodyType.NumField(); standard03CounterBodyIndex++ {
		if standard03CounterBodyType.Field(standard03CounterBodyIndex).Type.Kind() != standard03CounterBodyReflect.Pointer {
			continue
		}
		standard03CounterBodyZero := standard03CounterBodyB524
		standard03CounterBodyValue := standard03CounterBodyReflect.ValueOf(&standard03CounterBodyZero).Elem().Field(standard03CounterBodyIndex)
		standard03CounterBodyValue.Set(standard03CounterBodyReflect.New(standard03CounterBodyValue.Type().Elem()))
		standard03CounterBodyAbsent := standard03CounterBodyZero
		standard03CounterBodyReflect.ValueOf(&standard03CounterBodyAbsent).Elem().Field(standard03CounterBodyIndex).Set(standard03CounterBodyReflect.Zero(standard03CounterBodyValue.Type()))
		if !standard03CounterBodyEqualKey(standard03CounterBodyZero, standard03CounterBodyZero) || standard03CounterBodyEqualKey(standard03CounterBodyZero, standard03CounterBodyAbsent) || standard03CounterBodyEqualKey(standard03CounterBodyAbsent, standard03CounterBodyZero) {
			standard03CounterBodyT.Fatalf("absent versus zero collapsed for %s", standard03CounterBodyType.Field(standard03CounterBodyIndex).Name)
		}
	}

}

func standard03CounterBodyMutateNested(standard03CounterBodyValue standard03CounterBodyReflect.Value) {
	switch standard03CounterBodyValue.Kind() {
	case standard03CounterBodyReflect.Pointer:
		if !standard03CounterBodyValue.IsNil() {
			standard03CounterBodyMutateNested(standard03CounterBodyValue.Elem())
		}
	case standard03CounterBodyReflect.Struct:
		for standard03CounterBodyIndex := 0; standard03CounterBodyIndex < standard03CounterBodyValue.NumField(); standard03CounterBodyIndex++ {
			standard03CounterBodyMutateNested(standard03CounterBodyValue.Field(standard03CounterBodyIndex))
		}
	case standard03CounterBodyReflect.Slice:
		for standard03CounterBodyIndex := 0; standard03CounterBodyIndex < standard03CounterBodyValue.Len(); standard03CounterBodyIndex++ {
			standard03CounterBodyMutateNested(standard03CounterBodyValue.Index(standard03CounterBodyIndex))
		}
		if standard03CounterBodyValue.CanSet() {
			standard03CounterBodyValue.Set(standard03CounterBodyReflect.MakeSlice(standard03CounterBodyValue.Type(), 1, 1))
		}
	case standard03CounterBodyReflect.Map:
		if standard03CounterBodyValue.CanSet() {
			standard03CounterBodyValue.Set(standard03CounterBodyReflect.MakeMap(standard03CounterBodyValue.Type()))
			standard03CounterBodyValue.SetMapIndex(standard03CounterBodyReflect.ValueOf("synthetic"), standard03CounterBodyReflect.ValueOf("changed"))
		}
	case standard03CounterBodyReflect.String:
		if standard03CounterBodyValue.CanSet() {
			standard03CounterBodyValue.SetString("changed")
		}
	case standard03CounterBodyReflect.Int:
		if standard03CounterBodyValue.CanSet() {
			standard03CounterBodyValue.SetInt(17)
		}
	case standard03CounterBodyReflect.Uint8, standard03CounterBodyReflect.Uint16:
		if standard03CounterBodyValue.CanSet() {
			standard03CounterBodyValue.SetUint(17)
		}
	case standard03CounterBodyReflect.Bool:
		if standard03CounterBodyValue.CanSet() {
			standard03CounterBodyValue.SetBool(!standard03CounterBodyValue.Bool())
		}
	}
}

func TestAb83DetachedResults(standard03CounterBodyT *standard03CounterBodyTesting.T) {
	standard03CounterBodyExpected := Standard03CounterBodies()
	standard03CounterBodyChanged := Standard03CounterBodies()
	standard03CounterBodyMutateNested(standard03CounterBodyReflect.ValueOf(&standard03CounterBodyChanged).Elem())
	if !standard03CounterBodyReflect.DeepEqual(Standard03CounterBodies(), standard03CounterBodyExpected) {
		standard03CounterBodyT.Fatal("accessor shared storage")
	}
	for _, standard03CounterBodyBody := range standard03CounterBodyExpected {
		standard03CounterBodyChangedBody, standard03CounterBodyOK := Standard03CounterLookupBody(standard03CounterBodyBody.Standard03CounterBodyKey)
		if !standard03CounterBodyOK {
			standard03CounterBodyT.Fatal("body missing")
		}
		standard03CounterBodyMutateNested(standard03CounterBodyReflect.ValueOf(&standard03CounterBodyChangedBody).Elem())
		standard03CounterBodyAgain, _ := Standard03CounterLookupBody(standard03CounterBodyBody.Standard03CounterBodyKey)
		if !standard03CounterBodyReflect.DeepEqual(standard03CounterBodyAgain, standard03CounterBodyBody) {
			standard03CounterBodyT.Fatal("body lookup shared storage")
		}
		for _, standard03CounterBodyField := range standard03CounterBodyBody.Standard03CounterBodyFields {
			standard03CounterBodyKey := standard03CounterBodyBody.Standard03CounterBodyKey
			standard03CounterBodyKey.Standard03CounterBodyField = standard03CounterBodyField.Standard03CounterBodyIdentity
			standard03CounterBodyChangedField, standard03CounterBodyFound := Standard03CounterBodyLookupField(standard03CounterBodyKey)
			if !standard03CounterBodyFound {
				standard03CounterBodyT.Fatal("field missing")
			}
			standard03CounterBodyMutateNested(standard03CounterBodyReflect.ValueOf(&standard03CounterBodyChangedField).Elem())
			standard03CounterBodyAgainField, _ := Standard03CounterBodyLookupField(standard03CounterBodyKey)
			if !standard03CounterBodyReflect.DeepEqual(standard03CounterBodyAgainField, standard03CounterBodyField) {
				standard03CounterBodyT.Fatal("field lookup shared storage")
			}
		}
	}
}

func TestAb83MetadataCannotEnableStub(standard03CounterBodyT *standard03CounterBodyTesting.T) {
	standard03CounterBodyProvider, standard03CounterBodyErr := NewProvider(MustEmbeddedCatalog(), true)
	if standard03CounterBodyErr != nil {
		standard03CounterBodyT.Fatal(standard03CounterBodyErr)
	}
	standard03CounterBodyIDs := []string{"ebus_standard.service_data.start_counts", "ebus_standard.service_data.operating_time_l1", "ebus_standard.service_data.operating_time_l2", "ebus_standard.service_data.operating_time_l3", "ebus_standard.service_data.fuel_quantity_counter", "ebus_standard.service_data.meter_reading"}
	for _, standard03CounterBodyID := range standard03CounterBodyIDs {
		standard03CounterBodyBefore, standard03CounterBodyErr := standard03CounterBodyProvider.Invoke(standard03CounterBodyContext.Background(), standard03CounterBodyID, nil, CallerContext(0))
		if standard03CounterBodyErr != nil {
			standard03CounterBodyT.Fatal(standard03CounterBodyErr)
		}
		for _, standard03CounterBodyBody := range Standard03CounterBodies() {
			standard03CounterBodySelected, standard03CounterBodyOK := Standard03CounterLookupBody(standard03CounterBodyBody.Standard03CounterBodyKey)
			if !standard03CounterBodyOK || standard03CounterBodySelected.Standard03CounterBodyCallable {
				standard03CounterBodyT.Fatal("metadata advertised callable body")
			}
			standard03CounterBodySelected.Standard03CounterBodyCallable = true
		}
		standard03CounterBodyAfter, standard03CounterBodyErr := standard03CounterBodyProvider.Invoke(standard03CounterBodyContext.Background(), standard03CounterBodyID, nil, CallerContext(0))
		if standard03CounterBodyErr != nil || standard03CounterBodyAfter["executed"] != false || !standard03CounterBodyReflect.DeepEqual(standard03CounterBodyBefore, standard03CounterBodyAfter) {
			standard03CounterBodyT.Fatalf("stub behavior changed: %v %v", standard03CounterBodyAfter, standard03CounterBodyErr)
		}
	}
	standard03CounterBodyIdentification, standard03CounterBodyErr := standard03CounterBodyProvider.Identification(standard03CounterBodyContext.Background())
	if standard03CounterBodyErr != nil || standard03CounterBodyIdentification.Valid {
		standard03CounterBodyT.Fatal("metadata qualified identification")
	}
}
