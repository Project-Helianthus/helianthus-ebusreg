package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func fanDemandCatalog(fanDemandT *testing.T) map[string]any {
	fanDemandT.Helper()
	var fanDemandValue map[string]any
	if fanDemandError := json.Unmarshal(FanDemandCatalog(), &fanDemandValue); fanDemandError != nil {
		fanDemandT.Fatal(fanDemandError)
	}
	return fanDemandValue
}

func fanDemandField(fanDemandName any, fanDemandWidth int, fanDemandMax int, fanDemandDefault any, fanDemandMappings map[string]any, fanDemandConversion string, fanDemandUnit any, fanDemandBits any, fanDemandReserved []any) map[string]any {
	if fanDemandMappings == nil {
		fanDemandMappings = map[string]any{}
	}
	if fanDemandReserved == nil {
		fanDemandReserved = []any{}
	}
	fanDemandOrder := "not_applicable"
	if fanDemandWidth > 8 {
		fanDemandOrder = "Unknown"
	}
	return map[string]any{"name": fanDemandName, "logical_position": nil, "width_bits": fanDemandWidth, "raw_domain": []int{0, fanDemandMax}, "signedness": "unsigned", "byte_order": fanDemandOrder, "conversion": fanDemandConversion, "unit": fanDemandUnit, "default": fanDemandDefault, "mappings": fanDemandMappings, "bits": fanDemandBits, "reserved": fanDemandReserved, "grade": "Hypothesis", "raw_fallback": "preserve original bytes including reserved bits"}
}

func fanDemandExpectedRows() []any {
	fanDemandRows := []any{}
	fanDemandAdd := func(fanDemandID string, fanDemandDirection any, fanDemandAlternative string, fanDemandExtent any, fanDemandFields []any) {
		fanDemandRows = append(fanDemandRows, map[string]any{"id": fanDemandID, "family": "vaillant", "variant": "B522", "direction": fanDemandDirection, "service": "B522", "operation": nil, "alternative": fanDemandAlternative, "logical_extent_bytes": fanDemandExtent, "logical_coordinates": "include prefix positions; not physical offsets", "fields": fanDemandFields, "grade": "Hypothesis", "raw_fallback": "preserve full original carrier", "physical_mapping": "Unknown", "installed_support": "Unknown", "execution": "Unknown", "semantic_equivalence": "Unknown"})
	}
	fanDemandAdd("ab47_current_mode", nil, "current-mode", nil, []any{fanDemandField("current_mode", 8, 255, 0, map[string]any{
		"0": []any{"no heating demand to source/circuits"},
		"1": []any{"no source demand but demand from a sink", nil},
		"2": []any{"source heating demand"},
		"3": []any{"no source cooling demand but demand from a sink"},
		"4": []any{"source cooling demand"},
	}, "identity", nil, nil, nil)})
	fanDemandAdd("ab47_request", "request", "request", 7, []any{
		fanDemandField("fan", 8, 255, 0, nil, "identity", nil, nil, nil),
		fanDemandField("automatic", 1, 1, 0, map[string]any{"0": []any{"disabled"}, "1": []any{"enabled"}}, "identity", nil, []int{0, 0}, []any{map[string]any{"bits": []int{1, 7}, "expected": 0}}),
		fanDemandField("room_temperature", 8, 255, 0, nil, "raw * 0.5", "degrees C", nil, nil),
		fanDemandField("mode", 8, 255, 0, nil, "identity", nil, nil, nil),
	})
	fanDemandPositive := []any{}
	for fanDemandBit := 0; fanDemandBit < 3; fanDemandBit++ {
		fanDemandPositive = append(fanDemandPositive, fanDemandField(nil, 1, 1, false, map[string]any{"0": []any{false}, "1": []any{true}}, "identity", nil, []int{fanDemandBit, fanDemandBit}, nil))
	}
	fanDemandAdd("ab47_positive", "response", "positive", 4, fanDemandPositive)
	fanDemandRows[2].(map[string]any)["reserved"] = []any{map[string]any{"bits": []int{3, 7}, "expected": 0}}
	fanDemandAdd("ab47_negative", "response", "negative", nil, []any{})
	fanDemandUnset := fanDemandField(nil, 16, 65535, nil, map[string]any{"FF FF": []any{"candidate unset"}}, "Unknown", nil, nil, nil)
	fanDemandUnset["width_bits"] = nil
	fanDemandUnset["raw_domain"] = nil
	fanDemandUnset["signedness"] = "Unknown"
	fanDemandAdd("ab47_unset", nil, "FF FF unset", nil, []any{fanDemandUnset})
	return fanDemandRows
}

func fanDemandNormalized(fanDemandT *testing.T, fanDemandValue any) any {
	fanDemandT.Helper()
	fanDemandBytes, fanDemandError := json.Marshal(fanDemandValue)
	if fanDemandError != nil {
		fanDemandT.Fatal(fanDemandError)
	}
	var fanDemandResult any
	if fanDemandError = json.Unmarshal(fanDemandBytes, &fanDemandResult); fanDemandError != nil {
		fanDemandT.Fatal(fanDemandError)
	}
	return fanDemandResult
}

func TestAb47ExactCatalog(fanDemandT *testing.T) {
	fanDemandGot := fanDemandCatalog(fanDemandT)
	fanDemandWant := fanDemandNormalized(fanDemandT, fanDemandExpectedRows()).([]any)
	fanDemandRows, fanDemandOK := fanDemandGot["entries"].([]any)
	if !fanDemandOK || len(fanDemandRows) != 5 {
		fanDemandT.Fatal("ab47 expected exactly five separate rows")
	}
	for fanDemandIndex, fanDemandRow := range fanDemandWant {
		for fanDemandKey, fanDemandValue := range fanDemandRow.(map[string]any) {
			fanDemandT.Run(fmt.Sprintf("ab47_row_%d/%s", fanDemandIndex, fanDemandKey), func(fanDemandT *testing.T) {
				if !reflect.DeepEqual(fanDemandRows[fanDemandIndex].(map[string]any)[fanDemandKey], fanDemandValue) {
					fanDemandT.Fatalf("ab47 %s: got %#v want %#v", fanDemandKey, fanDemandRows[fanDemandIndex].(map[string]any)[fanDemandKey], fanDemandValue)
				}
			})
		}
		if !reflect.DeepEqual(fanDemandRows[fanDemandIndex], fanDemandRow) {
			fanDemandT.Errorf("ab47 row %d differs from finite contract", fanDemandIndex)
		}
	}
	if fanDemandGot["software_representation"] != "Proven" || fanDemandGot["candidate_semantics"] != "Hypothesis" || fanDemandGot["physical_qualification"] != "Unknown" {
		fanDemandT.Fatal("ab47 grades conflated")
	}
	fanDemandLimitations := []string{
		"Partial field catalog; absent names and coordinates remain unknown.",
		"Raw current modes 5..255 and all mode values remain unlabeled.",
		"Current-mode raw 1 has two unresolved labels; the competing label is Unknown.",
		"Heating priority during simultaneous cooling and bypass behavior are candidate semantics only.",
		"FF is not a sentinel for room temperature, mode or fan.",
		"Separate FF FF unset interpretation does not promote individual FF carriers.",
		"No runtime registration, wire builder, automatic selection, native codec promotion or executable capability.",
	}
	if !reflect.DeepEqual(fanDemandGot["limitations"], fanDemandNormalized(fanDemandT, fanDemandLimitations)) {
		fanDemandT.Fatal("ab47 limitations changed")
	}
}

func TestAb47IdentityAxes(fanDemandT *testing.T) {
	fanDemandRows := fanDemandCatalog(fanDemandT)["entries"].([]any)
	for fanDemandIndex, fanDemandValue := range fanDemandRows {
		for _, fanDemandAxis := range []string{"family", "variant", "direction", "service", "operation", "alternative"} {
			fanDemandT.Run(fmt.Sprintf("ab47_row_%d/%s", fanDemandIndex, fanDemandAxis), func(fanDemandT *testing.T) {
				fanDemandOriginal := fanDemandValue.(map[string]any)
				fanDemandChanged := fanDemandNormalized(fanDemandT, fanDemandOriginal).(map[string]any)
				fanDemandChanged[fanDemandAxis] = "ab47_distinct"
				fanDemandRoundTrip := fanDemandNormalized(fanDemandT, []any{fanDemandOriginal, fanDemandChanged}).([]any)
				if !reflect.DeepEqual(fanDemandRoundTrip[0], fanDemandOriginal) || reflect.DeepEqual(fanDemandRoundTrip[0], fanDemandRoundTrip[1]) || fanDemandRoundTrip[1].(map[string]any)[fanDemandAxis] != "ab47_distinct" {
					fanDemandT.Fatal("ab47 identity axis lost")
				}
				fanDemandChanged[fanDemandAxis] = fanDemandOriginal[fanDemandAxis]
				if !reflect.DeepEqual(fanDemandChanged, fanDemandOriginal) {
					fanDemandT.Fatal("ab47 invalid control")
				}
			})
		}
	}
}

func TestAb47RawCarriers(fanDemandT *testing.T) {
	fanDemandRows := fanDemandCatalog(fanDemandT)["entries"].([]any)
	fanDemandFields := fanDemandRows[1].(map[string]any)["fields"].([]any)
	fanDemandRoom := fanDemandFields[2].(map[string]any)
	if fanDemandRoom["conversion"] != "raw * 0.5" {
		fanDemandT.Fatal("ab47 temperature conversion changed")
	}
	for fanDemandIndex, fanDemandRaw := range []float64{0, 1, 255} {
		if fanDemandRaw*0.5 != []float64{0, 0.5, 127.5}[fanDemandIndex] {
			fanDemandT.Fatal("ab47 temperature examples")
		}
	}
	fanDemandFan := fanDemandFields[0].(map[string]any)
	if fanDemandFan["raw_domain"].([]any)[1] != float64(255) || fanDemandFan["conversion"] != "identity" || len(fanDemandFan["mappings"].(map[string]any)) != 0 {
		fanDemandT.Fatal("ab47 fan255 lost")
	}
	fanDemandModes := fanDemandRows[0].(map[string]any)["fields"].([]any)[0].(map[string]any)["mappings"].(map[string]any)
	if len(fanDemandModes["1"].([]any)) != 2 || fanDemandModes["1"].([]any)[1] != nil || fanDemandModes["5"] != nil {
		fanDemandT.Fatal("ab47 unresolved mode alternatives lost")
	}
	fanDemandAutomatic := fanDemandFields[1].(map[string]any)
	fanDemandPositive := fanDemandRows[2].(map[string]any)
	if !reflect.DeepEqual(fanDemandAutomatic["bits"], []any{float64(0), float64(0)}) || !reflect.DeepEqual(fanDemandAutomatic["reserved"], fanDemandNormalized(fanDemandT, []any{map[string]any{"bits": []int{1, 7}, "expected": 0}})) || len(fanDemandPositive["fields"].([]any)) != 3 || !reflect.DeepEqual(fanDemandPositive["reserved"], fanDemandNormalized(fanDemandT, []any{map[string]any{"bits": []int{3, 7}, "expected": 0}})) {
		fanDemandT.Fatal("ab47 bit roles changed")
	}
	for _, fanDemandRaw := range []byte{1, 128} {
		fanDemandEnabled := fanDemandRaw&1 != 0
		fanDemandDiscrepancy := fanDemandRaw&254 != 0
		if fanDemandEnabled != (fanDemandRaw == 1) || fanDemandDiscrepancy != (fanDemandRaw == 128) {
			fanDemandT.Fatal("ab47 automatic examples")
		}
	}
	if byte(8)&7 != 0 || byte(8)&248 != 8 {
		fanDemandT.Fatal("ab47 positive08 must remain reserved")
	}
}

func TestAb47CopyIsolation(fanDemandT *testing.T) {
	fanDemandBefore := bytes.Clone(FanDemandCatalog())
	fanDemandReturned := FanDemandCatalog()
	if len(fanDemandReturned) == 0 {
		fanDemandT.Fatal("ab47 empty catalog")
	}
	for fanDemandIndex := range fanDemandReturned {
		fanDemandReturned[fanDemandIndex] ^= 255
	}
	if !bytes.Equal(FanDemandCatalog(), fanDemandBefore) {
		fanDemandT.Fatal("ab47 returned bytes alias catalog")
	}
}
