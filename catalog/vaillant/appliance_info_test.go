package vaillant

import (
	applianceInfoJSONPackage "encoding/json"
	applianceInfoFmt "fmt"
	applianceInfoReflect "reflect"
	applianceInfoStrconv "strconv"
	applianceInfoTesting "testing"
)

func applianceInfoDecode(applianceInfoT *applianceInfoTesting.T, applianceInfoBytes []byte) map[string]any {
	applianceInfoT.Helper()
	var applianceInfoValue map[string]any
	if applianceInfoError := applianceInfoJSONPackage.Unmarshal(applianceInfoBytes, &applianceInfoValue); applianceInfoError != nil {
		applianceInfoT.Fatal(applianceInfoError)
	}
	return applianceInfoValue
}

func applianceInfoJSON(applianceInfoText string) any {
	var applianceInfoValue any
	if applianceInfoError := applianceInfoJSONPackage.Unmarshal([]byte(applianceInfoText), &applianceInfoValue); applianceInfoError != nil {
		panic(applianceInfoError)
	}
	return applianceInfoValue
}

func applianceInfoExpectedField(applianceInfoName string, applianceInfoPosition any, applianceInfoDefault any, applianceInfoMappings, applianceInfoBits, applianceInfoReserved, applianceInfoDomain string) any {
	return map[string]any{
		"name": applianceInfoName, "logical_position": applianceInfoPosition, "width_bytes": 1,
		"raw_domain": applianceInfoJSON(applianceInfoDomain), "signedness": "unsigned",
		"byte_order": "not-applicable", "conversion": "identity", "unit": nil,
		"default": applianceInfoDefault, "mappings": applianceInfoJSON(applianceInfoMappings),
		"bits": applianceInfoJSON(applianceInfoBits), "reserved_bits": applianceInfoJSON(applianceInfoReserved),
		"grade": "Hypothesis", "raw_fallback": "preserve-raw-byte",
	}
}

func applianceInfoExpectedRow(applianceInfoID, applianceInfoDirection, applianceInfoOperation, applianceInfoAlternative string, applianceInfoExtent any, applianceInfoFields ...any) any {
	if applianceInfoFields == nil {
		applianceInfoFields = []any{}
	}
	return map[string]any{
		"id": applianceInfoID, "family": "B513", "variant": nil,
		"direction": applianceInfoDirection, "service": "B513", "operation": applianceInfoOperation,
		"alternative": applianceInfoAlternative, "logical_positions_include_prefix": true,
		"physical_offsets": nil, "extent_bytes": applianceInfoExtent, "grade": "Hypothesis",
		"raw_fallback": "preserve-raw-bytes", "reserved_regions": []any{}, "fields": applianceInfoFields,
	}
}

func applianceInfoExpected(applianceInfoT *applianceInfoTesting.T) map[string]any {
	applianceInfoT.Helper()
	applianceInfoValue := map[string]any{
		"representation_grade": "Proven", "candidate_grade": "Hypothesis", "physical_qualification": "Unknown",
		"coordinate_system": "logical-including-prefix-not-physical-offsets",
		"limitations":       applianceInfoJSON(`["physical-adaptation-Unknown","OP00-complement-Unknown","OP04-address-dictionaries-Unknown","unspecified-positions-and-extents-Unknown","variant-Unknown","no-runtime-registration","no-wire-builder","no-automatic-selection","no-native-codec-promotion","no-executable-capability"]`),
		"entries": []any{
			applianceInfoExpectedRow("ab33_op05_request", "request", "OP05", "request", nil,
				applianceInfoExpectedField("request_bits", nil, 8, `{}`, `{"0":{"default":false,"role":"request-data"},"1":{"default":false,"role":"request-data"},"2":{"default":false,"role":"request-data"},"3":{"default":true,"role":"request-data"},"4":{"default":false,"role":"request-data"},"5":{"default":false,"role":"request-data"}}`, `[6,7]`, `[0,255]`)),
			applianceInfoExpectedRow("ab33_op05_positive", "response", "OP05", "positive", nil,
				applianceInfoExpectedField("feedback", nil, nil, `{"1":"accepted"}`, `{}`, `[]`, `[1,1]`)),
			applianceInfoExpectedRow("ab33_op05_negative", "response", "OP05", "negative", nil,
				applianceInfoExpectedField("feedback", nil, nil, `{"0":"not-accepted-or-not-supported","2..255":"not-accepted-or-not-supported"}`, `{}`, `[]`, `[0,255]`)),
			applianceInfoExpectedRow("ab33_op06_request", "request", "OP06", "request", nil,
				applianceInfoExpectedField("light", 3, 0, `{"0":"off","100":"on","255":"replacement"}`, `{}`, `[]`, `[0,255]`),
				applianceInfoExpectedField("colour", 4, 0, `{}`, `{}`, `[]`, `[0,255]`)),
			applianceInfoExpectedRow("ab33_op06_positive", "response", "OP06", "positive", nil,
				applianceInfoExpectedField("feedback", 3, nil, `{"1":"accepted"}`, `{}`, `[]`, `[1,1]`)),
			applianceInfoExpectedRow("ab33_op06_negative", "response", "OP06", "negative", nil,
				applianceInfoExpectedField("feedback", 3, nil, `{"0":"rejected","2..255":"rejected"}`, `{}`, `[]`, `[0,255]`)),
			applianceInfoExpectedRow("ab33_op07_request", "request", "OP07", "request", 0),
			applianceInfoExpectedRow("ab33_op07_positive", "response", "OP07", "positive", nil,
				applianceInfoExpectedField("appliance_type", 3, 0, `{"0":"modulating","1":"single-step","2":"two-step"}`, `{}`, `[]`, `[0,255]`),
				applianceInfoExpectedField("appliance_flags", 4, 0, `{}`, `{"0":{"label":"floor-standing","default":false},"1":{"label":"condensing","default":false},"2":{"label":"fan-assisted","default":false},"3":{"label":"open-vented","default":false},"4":{"label":"high-modulation-range","default":false},"5":{"label":"water-heater","default":false},"6":{"label":"step-modulation","default":false},"7":{"label":"ECS","default":false}}`, `[]`, `[0,255]`),
				applianceInfoExpectedField("fuel_flags", 5, 0, `{}`, `{"0":{"label":"gas","default":false},"1":{"label":"oil","default":false},"2":{"label":"ABC","default":false}}`, `[3,4,5,6,7]`, `[0,255]`),
				applianceInfoExpectedField("tank_flags", 6, 0, `{}`, `{"0":{"label":"layer-tank","default":false},"1":{"label":"warm-start","default":false},"2":{"label":"external-tank","default":false},"3":{"label":"combination","default":false}}`, `[4,5,6,7]`, `[0,255]`),
				applianceInfoExpectedField("feature_flags", 7, 0, `{}`, `{"0":{"label":"solar-compact","default":false},"1":{"label":"basic-solar","default":false},"2":{"label":"AAPS","default":false},"3":{"label":"advanced-power","default":false},"4":{"label":"partial-load","default":false},"5":{"label":"parameter-primary","default":false}}`, `[6,7]`, `[0,255]`),
				applianceInfoExpectedField("pump_mode", 8, 0, `{"0":false,"1":true}`, `{}`, `[]`, `[0,255]`),
				applianceInfoExpectedField("opaque_9", 9, 0, `{}`, `{}`, `[]`, `[0,255]`),
				applianceInfoExpectedField("opaque_10", 10, 0, `{}`, `{}`, `[]`, `[0,255]`),
				applianceInfoExpectedField("opaque_11", 11, 0, `{}`, `{}`, `[]`, `[0,255]`)),
		},
	}
	applianceInfoBytes, applianceInfoError := applianceInfoJSONPackage.Marshal(applianceInfoValue)
	if applianceInfoError != nil {
		applianceInfoT.Fatal(applianceInfoError)
	}
	return applianceInfoDecode(applianceInfoT, applianceInfoBytes)
}

func applianceInfoAssert(applianceInfoT *applianceInfoTesting.T, applianceInfoGot, applianceInfoWant any) {
	applianceInfoT.Helper()
	switch applianceInfoExpectedValue := applianceInfoWant.(type) {
	case map[string]any:
		applianceInfoActualValue, applianceInfoOK := applianceInfoGot.(map[string]any)
		if !applianceInfoOK || len(applianceInfoActualValue) != len(applianceInfoExpectedValue) {
			applianceInfoT.Fatalf("ab33 object keys: got %v, want %v", applianceInfoGot, applianceInfoWant)
		}
		for applianceInfoKey, applianceInfoValue := range applianceInfoExpectedValue {
			applianceInfoT.Run("ab33_"+applianceInfoKey, func(applianceInfoT *applianceInfoTesting.T) {
				applianceInfoActual, applianceInfoPresent := applianceInfoActualValue[applianceInfoKey]
				if !applianceInfoPresent {
					applianceInfoT.Fatal("ab33 missing explicit property")
				}
				applianceInfoAssert(applianceInfoT, applianceInfoActual, applianceInfoValue)
			})
		}
	case []any:
		applianceInfoActualValue, applianceInfoOK := applianceInfoGot.([]any)
		if !applianceInfoOK || len(applianceInfoActualValue) != len(applianceInfoExpectedValue) {
			applianceInfoT.Fatalf("ab33 array extent: got %v, want %v", applianceInfoGot, applianceInfoWant)
		}
		for applianceInfoIndex, applianceInfoValue := range applianceInfoExpectedValue {
			applianceInfoT.Run(applianceInfoFmt.Sprintf("ab33_%d", applianceInfoIndex), func(applianceInfoT *applianceInfoTesting.T) {
				applianceInfoAssert(applianceInfoT, applianceInfoActualValue[applianceInfoIndex], applianceInfoValue)
			})
		}
	default:
		if !applianceInfoReflect.DeepEqual(applianceInfoGot, applianceInfoWant) {
			applianceInfoT.Fatalf("ab33 got %#v, want %#v", applianceInfoGot, applianceInfoWant)
		}
	}
}

func TestAb33Contract(applianceInfoT *applianceInfoTesting.T) {
	applianceInfoAssert(applianceInfoT, applianceInfoDecode(applianceInfoT, ApplianceInfoCatalog()), applianceInfoExpected(applianceInfoT))
}

func TestAb33IndependentBytes(applianceInfoT *applianceInfoTesting.T) {
	applianceInfoFirst := ApplianceInfoCatalog()
	applianceInfoControl := append([]byte(nil), applianceInfoFirst...)
	for applianceInfoIndex := range applianceInfoFirst {
		applianceInfoFirst[applianceInfoIndex] = 0
	}
	if !applianceInfoReflect.DeepEqual([]byte(ApplianceInfoCatalog()), applianceInfoControl) {
		applianceInfoT.Fatal("ab33 accessor shares storage")
	}
}

func TestAb33IdentityAxes(applianceInfoT *applianceInfoTesting.T) {
	for applianceInfoIndex := 0; applianceInfoIndex < 8; applianceInfoIndex++ {
		for _, applianceInfoAxis := range []string{"id", "family", "variant", "direction", "service", "operation", "alternative"} {
			applianceInfoT.Run(applianceInfoFmt.Sprintf("ab33_%d_%s", applianceInfoIndex, applianceInfoAxis), func(applianceInfoT *applianceInfoTesting.T) {
				applianceInfoCatalog := applianceInfoDecode(applianceInfoT, ApplianceInfoCatalog())
				applianceInfoRows := applianceInfoCatalog["entries"].([]any)
				applianceInfoRow := applianceInfoRows[applianceInfoIndex].(map[string]any)
				applianceInfoRow[applianceInfoAxis] = "ab33_changed"
				applianceInfoBytes, applianceInfoError := applianceInfoJSONPackage.Marshal(applianceInfoCatalog)
				if applianceInfoError != nil {
					applianceInfoT.Fatal(applianceInfoError)
				}
				applianceInfoRoundTrip := applianceInfoDecode(applianceInfoT, applianceInfoBytes)
				if !applianceInfoReflect.DeepEqual(applianceInfoRoundTrip, applianceInfoCatalog) || applianceInfoReflect.DeepEqual(applianceInfoRoundTrip, applianceInfoExpected(applianceInfoT)) {
					applianceInfoT.Fatal("ab33 identity axis lost or comparison did not detect its mutation")
				}
				if !applianceInfoReflect.DeepEqual(applianceInfoDecode(applianceInfoT, ApplianceInfoCatalog()), applianceInfoExpected(applianceInfoT)) {
					applianceInfoT.Fatal("ab33 local identity mutation escaped into catalog")
				}
			})
		}
	}
}

func applianceInfoFields(applianceInfoT *applianceInfoTesting.T, applianceInfoIndex int) []any {
	applianceInfoT.Helper()
	return applianceInfoDecode(applianceInfoT, ApplianceInfoCatalog())["entries"].([]any)[applianceInfoIndex].(map[string]any)["fields"].([]any)
}

func applianceInfoLabel(applianceInfoField map[string]any, applianceInfoRaw int) any {
	applianceInfoMappings := applianceInfoField["mappings"].(map[string]any)
	if applianceInfoLabel, applianceInfoPresent := applianceInfoMappings[applianceInfoStrconv.Itoa(applianceInfoRaw)]; applianceInfoPresent {
		return applianceInfoLabel
	}
	if applianceInfoRaw >= 2 && applianceInfoRaw <= 255 {
		return applianceInfoMappings["2..255"]
	}
	return nil
}

func TestAb33FeedbackAndRequestBits(applianceInfoT *applianceInfoTesting.T) {
	applianceInfoRequest := applianceInfoFields(applianceInfoT, 0)[0].(map[string]any)
	applianceInfoBits := applianceInfoRequest["bits"].(map[string]any)
	applianceInfoPacked := 0
	for applianceInfoBit := 0; applianceInfoBit < 8; applianceInfoBit++ {
		applianceInfoDefinition, applianceInfoPresent := applianceInfoBits[applianceInfoStrconv.Itoa(applianceInfoBit)]
		if applianceInfoPresent != (applianceInfoBit < 6) {
			applianceInfoT.Fatalf("ab33 bit %d request/reserved boundary", applianceInfoBit)
		}
		if applianceInfoPresent {
			applianceInfoDefault := applianceInfoDefinition.(map[string]any)["default"].(bool)
			if applianceInfoDefault != (applianceInfoBit == 3) {
				applianceInfoT.Fatalf("ab33 bit %d default", applianceInfoBit)
			}
			if applianceInfoDefault {
				applianceInfoPacked |= 1 << applianceInfoBit
			}
		}
	}
	if applianceInfoPacked != 0x08 || applianceInfoRequest["default"] != float64(applianceInfoPacked) {
		applianceInfoT.Fatal("ab33 request packed default is not 08")
	}
	for _, applianceInfoCase := range []struct {
		applianceInfoIndex int
		applianceInfoWant  string
	}{{2, "not-accepted-or-not-supported"}, {5, "rejected"}} {
		for applianceInfoRaw := 0; applianceInfoRaw <= 255; applianceInfoRaw++ {
			applianceInfoWant := any(applianceInfoCase.applianceInfoWant)
			if applianceInfoRaw == 1 {
				applianceInfoWant = nil
			}
			if applianceInfoLabel(applianceInfoFields(applianceInfoT, applianceInfoCase.applianceInfoIndex)[0].(map[string]any), applianceInfoRaw) != applianceInfoWant {
				applianceInfoT.Fatalf("ab33 feedback row %d raw %d", applianceInfoCase.applianceInfoIndex, applianceInfoRaw)
			}
		}
	}
}

func TestAb33InformationWalkingBitsAndRaw(applianceInfoT *applianceInfoTesting.T) {
	applianceInfoFieldsValue := applianceInfoFields(applianceInfoT, 7)
	for applianceInfoIndex, applianceInfoLabels := range [][]string{
		{"floor-standing", "condensing", "fan-assisted", "open-vented", "high-modulation-range", "water-heater", "step-modulation", "ECS"},
		{"gas", "oil", "ABC"},
		{"layer-tank", "warm-start", "external-tank", "combination"},
		{"solar-compact", "basic-solar", "AAPS", "advanced-power", "partial-load", "parameter-primary"},
	} {
		applianceInfoField := applianceInfoFieldsValue[applianceInfoIndex+1].(map[string]any)
		applianceInfoBits := applianceInfoField["bits"].(map[string]any)
		for applianceInfoBit := 0; applianceInfoBit < 8; applianceInfoBit++ {
			applianceInfoRaw := 1 << applianceInfoBit
			applianceInfoActive := []string{}
			for applianceInfoKey, applianceInfoDefinition := range applianceInfoBits {
				applianceInfoPosition, applianceInfoError := applianceInfoStrconv.Atoi(applianceInfoKey)
				if applianceInfoError != nil {
					applianceInfoT.Fatal(applianceInfoError)
				}
				if applianceInfoRaw&(1<<applianceInfoPosition) != 0 {
					applianceInfoActive = append(applianceInfoActive, applianceInfoDefinition.(map[string]any)["label"].(string))
				}
			}
			applianceInfoWant := []string{}
			if applianceInfoBit < len(applianceInfoLabels) {
				applianceInfoWant = append(applianceInfoWant, applianceInfoLabels[applianceInfoBit])
			}
			if !applianceInfoReflect.DeepEqual(applianceInfoActive, applianceInfoWant) {
				applianceInfoT.Fatalf("ab33 flag byte %d walking bit %d: %v", applianceInfoIndex+4, applianceInfoBit, applianceInfoActive)
			}
		}
	}
	if applianceInfoLabel(applianceInfoFieldsValue[0].(map[string]any), 3) != nil {
		applianceInfoT.Fatal("ab33 type 3 acquired a fabricated label")
	}
	for applianceInfoRaw, applianceInfoWant := range []any{false, true, nil} {
		if applianceInfoLabel(applianceInfoFieldsValue[5].(map[string]any), applianceInfoRaw) != applianceInfoWant {
			applianceInfoT.Fatalf("ab33 Boolean raw %d", applianceInfoRaw)
		}
	}
	applianceInfoLogical := []byte{0, 0, 0, 3, 0, 0, 0, 0, 2, 7, 8, 9}
	applianceInfoTail := []byte{}
	for _, applianceInfoValue := range applianceInfoFieldsValue[6:] {
		applianceInfoField := applianceInfoValue.(map[string]any)
		if len(applianceInfoField["mappings"].(map[string]any)) != 0 || applianceInfoField["raw_fallback"] != "preserve-raw-byte" {
			applianceInfoT.Fatal("ab33 opaque field was interpreted")
		}
		applianceInfoTail = append(applianceInfoTail, applianceInfoLogical[int(applianceInfoField["logical_position"].(float64))])
	}
	if !applianceInfoReflect.DeepEqual(applianceInfoTail, []byte{7, 8, 9}) {
		applianceInfoT.Fatalf("ab33 opaque tail lost: %v", applianceInfoTail)
	}
}
