package ebusstandard

import (
	standard05BurnerBodyContext "context"
	standard05BurnerBodyErrors "errors"
	standard05BurnerBodyReflect "reflect"
	standard05BurnerBodyTesting "testing"
)

func standard05BurnerBodyInt(standard05BurnerBodyValue int) *int { return &standard05BurnerBodyValue }

func standard05BurnerBodyKey(standard05BurnerBodyName, standard05BurnerBodyDirection string) Standard05BurnerBodyKey {
	return Standard05BurnerBodyKey{Standard05BurnerBodyFamily: "burner", Standard05BurnerBodyVariant: "05_01_candidate", Standard05BurnerBodyPrimary: Standard05BurnerBodyAxis{true, 5}, Standard05BurnerBodySecondary: Standard05BurnerBodyAxis{true, 1}, Standard05BurnerBodyDirection: standard05BurnerBodyDirection, Standard05BurnerBodyField: standard05BurnerBodyName}
}

func standard05BurnerBodyExpectedField(standard05BurnerBodyName, standard05BurnerBodyDirection string, standard05BurnerBodyPosition int, standard05BurnerBodySignedness, standard05BurnerBodyUnit string) Standard05BurnerField {
	return Standard05BurnerField{Standard05BurnerBodyKey: standard05BurnerBodyKey(standard05BurnerBodyName, standard05BurnerBodyDirection), Standard05BurnerBodyLogicalPosition: standard05BurnerBodyInt(standard05BurnerBodyPosition), Standard05BurnerBodyWidthBits: standard05BurnerBodyInt(8), Standard05BurnerBodySignedness: standard05BurnerBodySignedness, Standard05BurnerBodyByteOrder: "not_applicable", Standard05BurnerBodyConversion: "identity", Standard05BurnerBodyUnit: standard05BurnerBodyUnit, Standard05BurnerBodyReplacement: []int{}, Standard05BurnerBodyGrade: "Hypothesis", Standard05BurnerBodyRawFallback: "retain_raw_bytes"}
}

func standard05BurnerBodyExpectedFields() []Standard05BurnerField {
	standard05BurnerBodyFields := []Standard05BurnerField{
		standard05BurnerBodyExpectedField("primary", "request", 0, "unsigned", "none"),
		standard05BurnerBodyExpectedField("secondary", "request", 1, "unsigned", "none"),
		standard05BurnerBodyExpectedField("heat_state", "request", 2, "unsigned", "none"),
		standard05BurnerBodyExpectedField("boiler_target", "request", 3, "unsigned", "degrees_Celsius"),
		standard05BurnerBodyExpectedField("dhw_target", "request", 4, "unsigned", "degrees_Celsius"),
		standard05BurnerBodyExpectedField("outside_temperature", "request", 5, "signed", "degrees_Celsius"),
		standard05BurnerBodyExpectedField("setting_percentage", "request", 6, "unsigned", "percent"),
		standard05BurnerBodyExpectedField("primary", "positive_declaration", 0, "unsigned", "none"),
		standard05BurnerBodyExpectedField("secondary", "positive_declaration", 1, "unsigned", "none"),
	}
	standard05BurnerBodyFields[0].Standard05BurnerBodyConstant = standard05BurnerBodyInt(5)
	standard05BurnerBodyFields[1].Standard05BurnerBodyConstant = standard05BurnerBodyInt(1)
	standard05BurnerBodyFields[7].Standard05BurnerBodyConstant = standard05BurnerBodyInt(5)
	standard05BurnerBodyFields[8].Standard05BurnerBodyConstant = standard05BurnerBodyInt(1)
	for standard05BurnerBodyIndex := 2; standard05BurnerBodyIndex < 7; standard05BurnerBodyIndex++ {
		standard05BurnerBodyFields[standard05BurnerBodyIndex].Standard05BurnerBodyDefault = standard05BurnerBodyInt(0)
	}
	standard05BurnerBodyFields[2].Standard05BurnerBodyEnumeration = map[int]string{0: "shutdown", 85: "DHW", 170: "heating", 204: "emission_check", 238: "controller_stop"}
	standard05BurnerBodyFields[3].Standard05BurnerBodyBounds = &Standard05BurnerBodyBounds{0, 90}
	standard05BurnerBodyFields[4].Standard05BurnerBodyBounds = &Standard05BurnerBodyBounds{0, 90}
	standard05BurnerBodyFields[5].Standard05BurnerBodyBounds = &Standard05BurnerBodyBounds{-128, 127}
	standard05BurnerBodyFields[6].Standard05BurnerBodyBounds = &Standard05BurnerBodyBounds{0, 255}
	for _, standard05BurnerBodyName := range []string{"physical_positive_body", "acknowledgement", "outside_3f_replacement", "control_permission", "scheduling", "other_05_candidate_bodies", "dd_check_semantics"} {
		standard05BurnerBodyFields = append(standard05BurnerBodyFields, Standard05BurnerField{Standard05BurnerBodyKey: standard05BurnerBodyKey(standard05BurnerBodyName, "unprovided"), Standard05BurnerBodySignedness: "Unknown", Standard05BurnerBodyByteOrder: "Unknown", Standard05BurnerBodyConversion: "Unknown", Standard05BurnerBodyUnit: "Unknown", Standard05BurnerBodyGrade: "Unknown", Standard05BurnerBodyRawFallback: "retain_raw_bytes"})
	}
	return standard05BurnerBodyFields
}

func TestAb84ProductionMetadata(standard05BurnerBodyT *standard05BurnerBodyTesting.T) {
	standard05BurnerBodyExpected := Standard05BurnerBodies{Standard05BurnerBodyRepresentationGrade: "Proven", Standard05BurnerBodyCoordinates: "logical byte coordinates only; no physical offsets or positive-body layout", Standard05BurnerBodyNativeSupport: "Unknown", Standard05BurnerBodyQualification: "Unknown", Standard05BurnerBodyFreshness: "Unknown", Standard05BurnerBodyPersistence: "Unknown", Standard05BurnerBodyExecutionPermission: "Unknown", Standard05BurnerBodyFields: standard05BurnerBodyExpectedFields()}
	if standard05BurnerBodyActual := Standard05BurnerMetadata(); !standard05BurnerBodyReflect.DeepEqual(standard05BurnerBodyActual, standard05BurnerBodyExpected) {
		standard05BurnerBodyT.Fatalf("production metadata differs from independent expectation:\n%+v", standard05BurnerBodyActual)
	}
	for _, standard05BurnerBodyExpectedField := range standard05BurnerBodyExpected.Standard05BurnerBodyFields {
		standard05BurnerBodyActual, standard05BurnerBodyOK := Standard05LookupBurnerField(standard05BurnerBodyExpectedField.Standard05BurnerBodyKey)
		if !standard05BurnerBodyOK || !standard05BurnerBodyReflect.DeepEqual(standard05BurnerBodyActual, standard05BurnerBodyExpectedField) {
			standard05BurnerBodyT.Fatalf("lookup %s: %+v, %v", standard05BurnerBodyExpectedField.Standard05BurnerBodyKey.Standard05BurnerBodyField, standard05BurnerBodyActual, standard05BurnerBodyOK)
		}
	}
}

func standard05BurnerBodyEvaluate(standard05BurnerBodyField Standard05BurnerField, standard05BurnerBodyLogical []byte) (int, string, bool) {
	if standard05BurnerBodyField.Standard05BurnerBodyLogicalPosition == nil || standard05BurnerBodyField.Standard05BurnerBodyWidthBits == nil || *standard05BurnerBodyField.Standard05BurnerBodyWidthBits != 8 || standard05BurnerBodyField.Standard05BurnerBodyConversion != "identity" || standard05BurnerBodyField.Standard05BurnerBodyByteOrder != "not_applicable" {
		return 0, "", false
	}
	standard05BurnerBodyPosition := *standard05BurnerBodyField.Standard05BurnerBodyLogicalPosition
	if standard05BurnerBodyPosition < 0 || standard05BurnerBodyPosition >= len(standard05BurnerBodyLogical) {
		return 0, "", false
	}
	standard05BurnerBodyValue := int(standard05BurnerBodyLogical[standard05BurnerBodyPosition])
	for _, standard05BurnerBodyReplacement := range standard05BurnerBodyField.Standard05BurnerBodyReplacement {
		if standard05BurnerBodyValue == standard05BurnerBodyReplacement {
			return standard05BurnerBodyValue, "", false
		}
	}
	if standard05BurnerBodyField.Standard05BurnerBodySignedness == "signed" && standard05BurnerBodyValue >= 128 {
		standard05BurnerBodyValue -= 256
	} else if standard05BurnerBodyField.Standard05BurnerBodySignedness != "signed" && standard05BurnerBodyField.Standard05BurnerBodySignedness != "unsigned" {
		return 0, "", false
	}
	standard05BurnerBodyValid := standard05BurnerBodyField.Standard05BurnerBodyBounds == nil || (standard05BurnerBodyValue >= standard05BurnerBodyField.Standard05BurnerBodyBounds.Standard05BurnerBodyMin && standard05BurnerBodyValue <= standard05BurnerBodyField.Standard05BurnerBodyBounds.Standard05BurnerBodyMax)
	return standard05BurnerBodyValue, standard05BurnerBodyField.Standard05BurnerBodyEnumeration[standard05BurnerBodyValue], standard05BurnerBodyValid
}

func TestAb84SyntheticCandidate(standard05BurnerBodyT *standard05BurnerBodyTesting.T) {
	standard05BurnerBodyLogical := []byte{5, 1, 0xAA, 90, 0, 0xFF, 255}
	for _, standard05BurnerBodyCase := range []struct {
		Standard05BurnerBodyName  string
		Standard05BurnerBodyValue int
		Standard05BurnerBodyLabel string
	}{
		{"heat_state", 170, "heating"}, {"boiler_target", 90, ""}, {"dhw_target", 0, ""}, {"outside_temperature", -1, ""}, {"setting_percentage", 255, ""},
	} {
		standard05BurnerBodyField, standard05BurnerBodyOK := Standard05LookupBurnerField(standard05BurnerBodyKey(standard05BurnerBodyCase.Standard05BurnerBodyName, "request"))
		if !standard05BurnerBodyOK || standard05BurnerBodyField.Standard05BurnerBodyGrade != "Hypothesis" {
			standard05BurnerBodyT.Fatal("missing hypothesis field", standard05BurnerBodyCase.Standard05BurnerBodyName)
		}
		standard05BurnerBodyValue, standard05BurnerBodyLabel, standard05BurnerBodyValid := standard05BurnerBodyEvaluate(standard05BurnerBodyField, standard05BurnerBodyLogical)
		if !standard05BurnerBodyValid || standard05BurnerBodyValue != standard05BurnerBodyCase.Standard05BurnerBodyValue || standard05BurnerBodyLabel != standard05BurnerBodyCase.Standard05BurnerBodyLabel {
			standard05BurnerBodyT.Fatalf("%s: %d %q %v", standard05BurnerBodyCase.Standard05BurnerBodyName, standard05BurnerBodyValue, standard05BurnerBodyLabel, standard05BurnerBodyValid)
		}
	}
	for _, standard05BurnerBodyName := range []string{"boiler_target", "dhw_target"} {
		standard05BurnerBodyField, _ := Standard05LookupBurnerField(standard05BurnerBodyKey(standard05BurnerBodyName, "request"))
		standard05BurnerBodyLogical[*standard05BurnerBodyField.Standard05BurnerBodyLogicalPosition] = 91
		if _, _, standard05BurnerBodyValid := standard05BurnerBodyEvaluate(standard05BurnerBodyField, standard05BurnerBodyLogical); standard05BurnerBodyValid {
			standard05BurnerBodyT.Fatal("target 91 accepted", standard05BurnerBodyName)
		}
	}
	standard05BurnerBodyOutside, _ := Standard05LookupBurnerField(standard05BurnerBodyKey("outside_temperature", "request"))
	standard05BurnerBodyLogical[*standard05BurnerBodyOutside.Standard05BurnerBodyLogicalPosition] = 0x3F
	if standard05BurnerBodyValue, _, standard05BurnerBodyValid := standard05BurnerBodyEvaluate(standard05BurnerBodyOutside, standard05BurnerBodyLogical); !standard05BurnerBodyValid || standard05BurnerBodyValue != 63 {
		standard05BurnerBodyT.Fatal("3F is numeric raw 63", standard05BurnerBodyValue, standard05BurnerBodyValid)
	}
	standard05BurnerBodyHeat, _ := Standard05LookupBurnerField(standard05BurnerBodyKey("heat_state", "request"))
	if _, standard05BurnerBodyOK := standard05BurnerBodyHeat.Standard05BurnerBodyEnumeration[0xDD]; standard05BurnerBodyOK {
		standard05BurnerBodyT.Fatal("DD semantics enabled")
	}
	for standard05BurnerBodyRaw, standard05BurnerBodyLabel := range map[byte]string{0: "shutdown", 0x55: "DHW", 0xAA: "heating", 0xCC: "emission_check", 0xEE: "controller_stop"} {
		standard05BurnerBodyLogical[*standard05BurnerBodyHeat.Standard05BurnerBodyLogicalPosition] = standard05BurnerBodyRaw
		if _, standard05BurnerBodyActual, standard05BurnerBodyValid := standard05BurnerBodyEvaluate(standard05BurnerBodyHeat, standard05BurnerBodyLogical); !standard05BurnerBodyValid || standard05BurnerBodyActual != standard05BurnerBodyLabel {
			standard05BurnerBodyT.Fatal("heat alternative", standard05BurnerBodyRaw, standard05BurnerBodyActual)
		}
	}
}

func standard05BurnerBodyPositiveComplete(standard05BurnerBodyFields []Standard05BurnerField, standard05BurnerBodyLogical []byte) bool {
	standard05BurnerBodyCount := 0
	for _, standard05BurnerBodyField := range standard05BurnerBodyFields {
		if standard05BurnerBodyField.Standard05BurnerBodyKey.Standard05BurnerBodyDirection != "positive_declaration" {
			continue
		}
		standard05BurnerBodyCount++
		if standard05BurnerBodyField.Standard05BurnerBodyConstant == nil {
			return false
		}
		standard05BurnerBodyValue, _, standard05BurnerBodyValid := standard05BurnerBodyEvaluate(standard05BurnerBodyField, standard05BurnerBodyLogical)
		if !standard05BurnerBodyValid || standard05BurnerBodyValue != *standard05BurnerBodyField.Standard05BurnerBodyConstant {
			return false
		}
	}
	return standard05BurnerBodyCount == 2 && len(standard05BurnerBodyLogical) == 2
}

func TestAb84PositiveDeclaration(standard05BurnerBodyT *standard05BurnerBodyTesting.T) {
	standard05BurnerBodyFields := Standard05BurnerMetadata().Standard05BurnerBodyFields
	for _, standard05BurnerBodyCase := range []struct {
		Standard05BurnerBodyBytes    []byte
		Standard05BurnerBodyComplete bool
	}{
		{[]byte{5, 1}, true}, {[]byte{5}, false}, {[]byte{1}, false}, {[]byte{5, 4}, false}, {[]byte{5, 1, 0}, false}, {nil, false},
	} {
		if standard05BurnerBodyActual := standard05BurnerBodyPositiveComplete(standard05BurnerBodyFields, standard05BurnerBodyCase.Standard05BurnerBodyBytes); standard05BurnerBodyActual != standard05BurnerBodyCase.Standard05BurnerBodyComplete {
			standard05BurnerBodyT.Fatalf("logical constants %v: %v", standard05BurnerBodyCase.Standard05BurnerBodyBytes, standard05BurnerBodyActual)
		}
	}
}

func standard05BurnerBodyKeyMutations(standard05BurnerBodyControl Standard05BurnerBodyKey) []Standard05BurnerBodyKey {
	standard05BurnerBodyKeys := []Standard05BurnerBodyKey{}
	for standard05BurnerBodyIndex := 0; standard05BurnerBodyIndex < 11; standard05BurnerBodyIndex++ {
		standard05BurnerBodyChanged := standard05BurnerBodyControl
		switch standard05BurnerBodyIndex {
		case 0:
			standard05BurnerBodyChanged.Standard05BurnerBodyFamily += "_other"
		case 1:
			standard05BurnerBodyChanged.Standard05BurnerBodyVariant += "_other"
		case 2:
			standard05BurnerBodyChanged.Standard05BurnerBodyPrimary.Standard05BurnerBodyValue++
		case 3:
			standard05BurnerBodyChanged.Standard05BurnerBodySecondary.Standard05BurnerBodyValue++
		case 4:
			standard05BurnerBodyChanged.Standard05BurnerBodyOperation.Standard05BurnerBodyValue++
		case 5:
			standard05BurnerBodyChanged.Standard05BurnerBodyDirection += "_other"
		case 6:
			standard05BurnerBodyChanged.Standard05BurnerBodyField += "_other"
		case 7:
			standard05BurnerBodyChanged.Standard05BurnerBodyOpcode.Standard05BurnerBodyValue++
		case 8:
			standard05BurnerBodyChanged.Standard05BurnerBodyGroup.Standard05BurnerBodyValue++
		case 9:
			standard05BurnerBodyChanged.Standard05BurnerBodyInstance.Standard05BurnerBodyValue++
		case 10:
			standard05BurnerBodyChanged.Standard05BurnerBodyRegister.Standard05BurnerBodyValue++
		}
		standard05BurnerBodyKeys = append(standard05BurnerBodyKeys, standard05BurnerBodyChanged)
	}
	for standard05BurnerBodyIndex := 0; standard05BurnerBodyIndex < 7; standard05BurnerBodyIndex++ {
		standard05BurnerBodyChanged := standard05BurnerBodyControl
		standard05BurnerBodyAxes := []*Standard05BurnerBodyAxis{&standard05BurnerBodyChanged.Standard05BurnerBodyPrimary, &standard05BurnerBodyChanged.Standard05BurnerBodySecondary, &standard05BurnerBodyChanged.Standard05BurnerBodyOperation, &standard05BurnerBodyChanged.Standard05BurnerBodyOpcode, &standard05BurnerBodyChanged.Standard05BurnerBodyGroup, &standard05BurnerBodyChanged.Standard05BurnerBodyInstance, &standard05BurnerBodyChanged.Standard05BurnerBodyRegister}
		standard05BurnerBodyAxes[standard05BurnerBodyIndex].Standard05BurnerBodyPresent = !standard05BurnerBodyAxes[standard05BurnerBodyIndex].Standard05BurnerBodyPresent
		standard05BurnerBodyKeys = append(standard05BurnerBodyKeys, standard05BurnerBodyChanged)
	}
	return standard05BurnerBodyKeys
}

func TestAb84ExactKeys(standard05BurnerBodyT *standard05BurnerBodyTesting.T) {
	standard05BurnerBodyControl := standard05BurnerBodyKey("heat_state", "request")
	if _, standard05BurnerBodyOK := Standard05LookupBurnerField(standard05BurnerBodyControl); !standard05BurnerBodyOK {
		standard05BurnerBodyT.Fatal("valid control missing")
	}
	for standard05BurnerBodyIndex, standard05BurnerBodyChanged := range standard05BurnerBodyKeyMutations(standard05BurnerBodyControl) {
		if _, standard05BurnerBodyOK := Standard05LookupBurnerField(standard05BurnerBodyChanged); standard05BurnerBodyOK {
			standard05BurnerBodyT.Fatal("identity axis ignored", standard05BurnerBodyIndex)
		}
	}
	standard05BurnerBodyStandard := standard05BurnerBodyControl
	standard05BurnerBodyStandard.Standard05BurnerBodySecondary.Standard05BurnerBodyValue = 4
	standard05BurnerBodyStandard.Standard05BurnerBodyDirection = "response"
	if _, standard05BurnerBodyOK := Standard05LookupBurnerField(standard05BurnerBodyStandard); standard05BurnerBodyOK {
		standard05BurnerBodyT.Fatal("standard 05/04 aliases candidate")
	}
	// Synthetic keys exercise the same lookup without declaring B524 semantics.
	standard05BurnerBodyB524 := Standard05BurnerBodyKey{Standard05BurnerBodyFamily: "synthetic_family", Standard05BurnerBodyVariant: "synthetic_variant", Standard05BurnerBodyDirection: "response", Standard05BurnerBodyField: "synthetic_field", Standard05BurnerBodyPrimary: Standard05BurnerBodyAxis{true, 0xB5}, Standard05BurnerBodySecondary: Standard05BurnerBodyAxis{true, 0x24}, Standard05BurnerBodyOperation: Standard05BurnerBodyAxis{true, 0}, Standard05BurnerBodyOpcode: Standard05BurnerBodyAxis{true, 0}, Standard05BurnerBodyGroup: Standard05BurnerBodyAxis{true, 0}, Standard05BurnerBodyInstance: Standard05BurnerBodyAxis{true, 0}, Standard05BurnerBodyRegister: Standard05BurnerBodyAxis{true, 0}}
	for _, standard05BurnerBodyControl := range []Standard05BurnerBodyKey{standard05BurnerBodyB524, {Standard05BurnerBodyFamily: "synthetic_family", Standard05BurnerBodyVariant: "synthetic_variant", Standard05BurnerBodyDirection: "response", Standard05BurnerBodyField: "synthetic_field"}} {
		standard05BurnerBodyFields := []Standard05BurnerField{{Standard05BurnerBodyKey: standard05BurnerBodyControl, Standard05BurnerBodyGrade: "Unknown"}}
		if _, standard05BurnerBodyOK := standard05BurnerBodyLookupBurnerField(standard05BurnerBodyFields, standard05BurnerBodyControl); !standard05BurnerBodyOK {
			standard05BurnerBodyT.Fatal("synthetic control missing")
		}
		for standard05BurnerBodyIndex, standard05BurnerBodyChanged := range standard05BurnerBodyKeyMutations(standard05BurnerBodyControl) {
			if _, standard05BurnerBodyOK := standard05BurnerBodyLookupBurnerField(standard05BurnerBodyFields, standard05BurnerBodyChanged); standard05BurnerBodyOK {
				standard05BurnerBodyT.Fatal("synthetic axis ignored", standard05BurnerBodyIndex)
			}
		}
	}
}

func standard05BurnerBodyMutate(standard05BurnerBodyField *Standard05BurnerField) {
	standard05BurnerBodyField.Standard05BurnerBodyKey.Standard05BurnerBodyFamily = "mutated"
	if standard05BurnerBodyField.Standard05BurnerBodyLogicalPosition != nil {
		*standard05BurnerBodyField.Standard05BurnerBodyLogicalPosition = 999
	}
	if standard05BurnerBodyField.Standard05BurnerBodyWidthBits != nil {
		*standard05BurnerBodyField.Standard05BurnerBodyWidthBits = 999
	}
	if standard05BurnerBodyField.Standard05BurnerBodyDefault != nil {
		*standard05BurnerBodyField.Standard05BurnerBodyDefault = 999
	}
	if standard05BurnerBodyField.Standard05BurnerBodyConstant != nil {
		*standard05BurnerBodyField.Standard05BurnerBodyConstant = 999
	}
	if standard05BurnerBodyField.Standard05BurnerBodyBounds != nil {
		standard05BurnerBodyField.Standard05BurnerBodyBounds.Standard05BurnerBodyMin = 999
		standard05BurnerBodyField.Standard05BurnerBodyBounds.Standard05BurnerBodyMax = 999
	}
	if standard05BurnerBodyField.Standard05BurnerBodyEnumeration != nil {
		standard05BurnerBodyField.Standard05BurnerBodyEnumeration[0] = "mutated"
	}
	standard05BurnerBodyField.Standard05BurnerBodyReplacement = append(standard05BurnerBodyField.Standard05BurnerBodyReplacement, 999)
	standard05BurnerBodyField.Standard05BurnerBodyReserved = append(standard05BurnerBodyField.Standard05BurnerBodyReserved, 999)
}

func TestAb84DetachedCopies(standard05BurnerBodyT *standard05BurnerBodyTesting.T) {
	standard05BurnerBodyBodies := Standard05BurnerMetadata()
	for standard05BurnerBodyIndex := range standard05BurnerBodyBodies.Standard05BurnerBodyFields {
		standard05BurnerBodyMutate(&standard05BurnerBodyBodies.Standard05BurnerBodyFields[standard05BurnerBodyIndex])
	}
	if !standard05BurnerBodyReflect.DeepEqual(Standard05BurnerMetadata().Standard05BurnerBodyFields, standard05BurnerBodyExpectedFields()) {
		standard05BurnerBodyT.Fatal("accessor storage shared")
	}
	for _, standard05BurnerBodyExpected := range standard05BurnerBodyExpectedFields() {
		standard05BurnerBodyField, standard05BurnerBodyOK := Standard05LookupBurnerField(standard05BurnerBodyExpected.Standard05BurnerBodyKey)
		if !standard05BurnerBodyOK {
			standard05BurnerBodyT.Fatal("lookup missing")
		}
		standard05BurnerBodyMutate(&standard05BurnerBodyField)
		standard05BurnerBodyAgain, standard05BurnerBodyOK := Standard05LookupBurnerField(standard05BurnerBodyExpected.Standard05BurnerBodyKey)
		if !standard05BurnerBodyOK || !standard05BurnerBodyReflect.DeepEqual(standard05BurnerBodyAgain, standard05BurnerBodyExpected) {
			standard05BurnerBodyT.Fatal("lookup storage shared", standard05BurnerBodyExpected.Standard05BurnerBodyKey)
		}
	}
}

func TestAb84ProviderUnchanged(standard05BurnerBodyT *standard05BurnerBodyTesting.T) {
	standard05BurnerBodyProvider, standard05BurnerBodyError := NewProvider(MustEmbeddedCatalog(), true)
	if standard05BurnerBodyError != nil {
		standard05BurnerBodyT.Fatal(standard05BurnerBodyError)
	}
	_ = Standard05BurnerMetadata()
	for _, standard05BurnerBodyCaller := range []CallerContext{CallerContextUserFacing, CallerContextSystemNMRuntime} {
		if _, standard05BurnerBodyError := standard05BurnerBodyProvider.Invoke(standard05BurnerBodyContext.Background(), "ebus_standard.burner.op_data_c2b", nil, standard05BurnerBodyCaller); !standard05BurnerBodyErrors.Is(standard05BurnerBodyError, ErrSafetyClassDenied) {
			standard05BurnerBodyT.Fatal("burner became callable", standard05BurnerBodyError)
		}
	}
	standard05BurnerBodyResult, standard05BurnerBodyError := standard05BurnerBodyProvider.Invoke(standard05BurnerBodyContext.Background(), "ebus_standard.system_data.identification_addressed", nil, CallerContextUserFacing)
	if standard05BurnerBodyError != nil {
		standard05BurnerBodyT.Fatal(standard05BurnerBodyError)
	}
	if standard05BurnerBodyExecuted, standard05BurnerBodyOK := standard05BurnerBodyResult["executed"].(bool); !standard05BurnerBodyOK || standard05BurnerBodyExecuted {
		standard05BurnerBodyT.Fatal("existing stub executed", standard05BurnerBodyResult)
	}
	standard05BurnerBodyIdentity, standard05BurnerBodyError := standard05BurnerBodyProvider.Identification(standard05BurnerBodyContext.Background())
	if standard05BurnerBodyError != nil || standard05BurnerBodyIdentity.Valid || standard05BurnerBodyIdentity.Manufacturer != "" || standard05BurnerBodyIdentity.DeviceID != "" {
		standard05BurnerBodyT.Fatal("identification gained evidence", standard05BurnerBodyIdentity, standard05BurnerBodyError)
	}
}
