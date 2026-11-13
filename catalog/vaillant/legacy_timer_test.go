package vaillant

import (
	legacyTimerBytes "bytes"
	legacyTimerJSON "encoding/json"
	legacyTimerFmt "fmt"
	legacyTimerReflect "reflect"
	legacyTimerTesting "testing"
)

func legacyTimerDecode(legacyTimerT *legacyTimerTesting.T, legacyTimerData []byte) map[string]any {
	legacyTimerT.Helper()
	var legacyTimerValue map[string]any
	if legacyTimerError := legacyTimerJSON.Unmarshal(legacyTimerData, &legacyTimerValue); legacyTimerError != nil {
		legacyTimerT.Fatal(legacyTimerError)
	}
	return legacyTimerValue
}

func legacyTimerExpected(legacyTimerT *legacyTimerTesting.T) map[string]any {
	legacyTimerT.Helper()
	legacyTimerWant := legacyTimerDecode(legacyTimerT, []byte(`{
 "id":"ab35_legacy_timer","representation_grade":"Proven","candidate_grade":"Hypothesis","physical_qualification":"Unknown",
 "opt_in":true,"runtime_registration":false,"wire_builder":false,"automatic_selection":false,"native_codec_promotion":false,"write_admission":false,
 "coordinates":"logical coordinates include prefix positions; they are not physical offsets",
 "limitations":{"physical_forms":"Unknown","competing_two_position_gap":"Unknown","mode_labels":"Unknown","enforcement":"Unknown","deletion_reset":"Unknown","atomicity":"Unknown","overlap":"Unknown","midnight_crossing":"Unknown","timezone_DST":"Unknown","physical_schedule_effects":"Unknown"},
 "ordering_candidate":{"grade":"Hypothesis","written_start_equal_end":"disabled interval","interval_order":["first","second","third"],"out_of_order_qualification":"Unknown"}
 }`))
	legacyTimerEntries := []any{}
	for legacyTimerIndex, legacyTimerIdentity := range []struct {
		legacyTimerOperation, legacyTimerAlternative, legacyTimerDirection string
		legacyTimerExtent                                                  any
	}{
		{"read", "request", "request", float64(4)}, {"read", "positive", "response", float64(11)}, {"read", "negative", "response", nil},
		{"write", "request", "request", float64(11)}, {"write", "positive", "response", float64(3)}, {"write", "negative", "response", nil},
	} {
		legacyTimerRow := legacyTimerDecode(legacyTimerT, []byte(`{"family":"Vaillant legacy timer","variant":"periods","service":"B515","logical_positions":null,"physical_extent":null,"width":null,"raw_domain":null,"signedness":"Unknown","byte_order":"Unknown","conversion":null,"unit":null,"default":null,"finite_mappings":null,"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve complete native bytes and independent useful fields"}`))
		legacyTimerRow["id"] = "ab35_" + legacyTimerIdentity.legacyTimerOperation + "_" + legacyTimerIdentity.legacyTimerAlternative
		legacyTimerRow["operation"], legacyTimerRow["alternative"], legacyTimerRow["direction"], legacyTimerRow["logical_extent"] = legacyTimerIdentity.legacyTimerOperation, legacyTimerIdentity.legacyTimerAlternative, legacyTimerIdentity.legacyTimerDirection, legacyTimerIdentity.legacyTimerExtent
		legacyTimerFields := map[string]any{}
		if legacyTimerIndex != 2 && legacyTimerIndex != 4 && legacyTimerIndex != 5 {
			for _, legacyTimerName := range []string{"weekday", "channel", "mode"} {
				legacyTimerField := legacyTimerDecode(legacyTimerT, []byte(`{"logical_positions":null,"width":null,"raw_domain":[0,255],"signedness":"Unknown","byte_order":"Unknown","conversion":null,"unit":null,"default":null,"finite_mappings":null,"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve raw carrier"}`))
				if legacyTimerName == "weekday" {
					legacyTimerField["default"] = "Monday"
				}
				if legacyTimerName == "channel" {
					legacyTimerField["default"] = float64(0)
					if legacyTimerIndex == 0 || legacyTimerIndex == 3 {
						legacyTimerField["default"] = map[string]any{"circulation": float64(2), "DHW": float64(1), "mixer_one": float64(0), "mixer_two": float64(0)}
					}
				}
				legacyTimerFields[legacyTimerName] = legacyTimerField
			}
		}
		if legacyTimerIndex == 1 || legacyTimerIndex == 3 {
			legacyTimerPeriod := legacyTimerDecode(legacyTimerT, []byte(`{"logical_positions":null,"width":null,"raw_domain":null,"signedness":"Unknown","byte_order":"Unknown","conversion":"raw * 10 minutes for raw 0..144 only","unit":"minute","default":null,"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve raw label; 145 and 255 have Unknown roles","undeclared_time_conversion":"09:05","rounding":"Unknown"}`))
			legacyTimerMappings := map[string]any{}
			for legacyTimerMinutes := 0; legacyTimerMinutes <= 1440; legacyTimerMinutes += 10 {
				legacyTimerMappings[legacyTimerFmt.Sprint(legacyTimerMinutes/10)] = legacyTimerFmt.Sprintf("%02d:%02d", legacyTimerMinutes/60, legacyTimerMinutes%60)
			}
			legacyTimerPeriod["finite_mappings"] = legacyTimerMappings
			legacyTimerFields["period"] = legacyTimerPeriod
		}
		legacyTimerRow["fields"] = legacyTimerFields
		legacyTimerEntries = append(legacyTimerEntries, legacyTimerRow)
	}
	for legacyTimerIndex := 1; legacyTimerIndex <= 5; legacyTimerIndex++ {
		legacyTimerRow := legacyTimerDecode(legacyTimerT, []byte(`{"family":"Vaillant legacy timer","direction":null,"service":"B515","operation":null,"logical_positions":null,"logical_extent":null,"physical_extent":null,"width":null,"raw_domain":[0,3],"signedness":"Unknown","byte_order":"Unknown","conversion":null,"unit":null,"default":null,"finite_mappings":null,"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve raw mode and unresolved variant","mapping_injective":false,"inverse_default":null,"fields":{}}`))
		legacyTimerRow["id"] = legacyTimerFmt.Sprintf("ab35_mode_candidate_%d", legacyTimerIndex)
		legacyTimerRow["variant"] = legacyTimerFmt.Sprintf("mode_candidate_%d", legacyTimerIndex)
		legacyTimerRow["alternative"] = legacyTimerRow["variant"]
		legacyTimerEntries = append(legacyTimerEntries, legacyTimerRow)
	}
	legacyTimerWant["entries"] = legacyTimerEntries
	return legacyTimerWant
}

func legacyTimerAssert(legacyTimerT *legacyTimerTesting.T, legacyTimerPath string, legacyTimerGot, legacyTimerWant any) {
	legacyTimerT.Helper()
	switch legacyTimerValue := legacyTimerWant.(type) {
	case map[string]any:
		legacyTimerActual, legacyTimerOK := legacyTimerGot.(map[string]any)
		if !legacyTimerOK || len(legacyTimerActual) != len(legacyTimerValue) {
			legacyTimerT.Fatalf("%s: object keys differ", legacyTimerPath)
		}
		for legacyTimerKey, legacyTimerChild := range legacyTimerValue {
			legacyTimerActualChild, legacyTimerPresent := legacyTimerActual[legacyTimerKey]
			if !legacyTimerPresent {
				legacyTimerT.Fatalf("%s/%s: missing explicit property", legacyTimerPath, legacyTimerKey)
			}
			legacyTimerAssert(legacyTimerT, legacyTimerPath+"/"+legacyTimerKey, legacyTimerActualChild, legacyTimerChild)
		}
	case []any:
		legacyTimerActual, legacyTimerOK := legacyTimerGot.([]any)
		if !legacyTimerOK || len(legacyTimerActual) != len(legacyTimerValue) {
			legacyTimerT.Fatalf("%s: alternative/domain count differs", legacyTimerPath)
		}
		for legacyTimerIndex, legacyTimerChild := range legacyTimerValue {
			legacyTimerAssert(legacyTimerT, legacyTimerFmt.Sprintf("%s/%d", legacyTimerPath, legacyTimerIndex), legacyTimerActual[legacyTimerIndex], legacyTimerChild)
		}
	default:
		if !legacyTimerReflect.DeepEqual(legacyTimerGot, legacyTimerWant) {
			legacyTimerT.Errorf("%s: got %v, want %v", legacyTimerPath, legacyTimerGot, legacyTimerWant)
		}
	}
}

func TestAb35CatalogContract(legacyTimerT *legacyTimerTesting.T) {
	legacyTimerAssert(legacyTimerT, "ab35", legacyTimerDecode(legacyTimerT, LegacyTimerCatalog()), legacyTimerExpected(legacyTimerT))
}

func TestAb35CopyOwnership(legacyTimerT *legacyTimerTesting.T) {
	legacyTimerFirst, legacyTimerSecond := LegacyTimerCatalog(), LegacyTimerCatalog()
	if len(legacyTimerFirst) == 0 || !legacyTimerBytes.Equal(legacyTimerFirst, legacyTimerSecond) {
		legacyTimerT.Fatal("ab35 inconsistent catalog")
	}
	legacyTimerSnapshot := legacyTimerBytes.Clone(legacyTimerSecond)
	for legacyTimerIndex := range legacyTimerFirst {
		legacyTimerFirst[legacyTimerIndex] ^= 255
	}
	if !legacyTimerBytes.Equal(legacyTimerSnapshot, legacyTimerSecond) || !legacyTimerBytes.Equal(legacyTimerSnapshot, LegacyTimerCatalog()) {
		legacyTimerT.Fatal("ab35 caller mutation changed catalog")
	}
}

func TestAb35PeriodVocabulary(legacyTimerT *legacyTimerTesting.T) {
	legacyTimerRows := legacyTimerDecode(legacyTimerT, LegacyTimerCatalog())["entries"].([]any)
	for _, legacyTimerIndex := range []int{1, 3} {
		legacyTimerPeriod := legacyTimerRows[legacyTimerIndex].(map[string]any)["fields"].(map[string]any)["period"].(map[string]any)
		legacyTimerMappings := legacyTimerPeriod["finite_mappings"].(map[string]any)
		if len(legacyTimerMappings) != 145 {
			legacyTimerT.Fatal("ab35 requires 145 labels")
		}
		for legacyTimerRaw, legacyTimerWant := range map[string]string{"0": "00:00", "54": "09:00", "144": "24:00"} {
			if legacyTimerMappings[legacyTimerRaw] != legacyTimerWant {
				legacyTimerT.Fatalf("ab35 %s: %v", legacyTimerRaw, legacyTimerMappings[legacyTimerRaw])
			}
		}
		for _, legacyTimerRaw := range []string{"145", "255"} {
			if _, legacyTimerPresent := legacyTimerMappings[legacyTimerRaw]; legacyTimerPresent {
				legacyTimerT.Fatal("ab35 unknown raw role assigned")
			}
		}
		for _, legacyTimerLabel := range legacyTimerMappings {
			if legacyTimerLabel == "09:05" {
				legacyTimerT.Fatal("ab35 undeclared rounding")
			}
		}
		if legacyTimerPeriod["rounding"] != "Unknown" {
			legacyTimerT.Fatal("ab35 rounding promoted")
		}
	}
}

func TestAb35IdentityAxes(legacyTimerT *legacyTimerTesting.T) {
	legacyTimerRows := legacyTimerDecode(legacyTimerT, LegacyTimerCatalog())["entries"].([]any)
	legacyTimerWant := legacyTimerExpected(legacyTimerT)["entries"].([]any)
	for legacyTimerIndex, legacyTimerEntry := range legacyTimerRows {
		for _, legacyTimerAxis := range []string{"id", "family", "variant", "direction", "service", "operation", "alternative"} {
			legacyTimerT.Run(legacyTimerFmt.Sprintf("ab35_%d_%s", legacyTimerIndex, legacyTimerAxis), func(legacyTimerT *legacyTimerTesting.T) {
				legacyTimerRow := legacyTimerEntry.(map[string]any)
				if !legacyTimerReflect.DeepEqual(legacyTimerRow[legacyTimerAxis], legacyTimerWant[legacyTimerIndex].(map[string]any)[legacyTimerAxis]) {
					legacyTimerT.Fatal("ab35 identity changed")
				}
				legacyTimerPair := []any{legacyTimerRow, legacyTimerRow}
				legacyTimerWire, legacyTimerError := legacyTimerJSON.Marshal(legacyTimerPair)
				if legacyTimerError != nil {
					legacyTimerT.Fatal(legacyTimerError)
				}
				var legacyTimerRoundTrip []map[string]any
				if legacyTimerError = legacyTimerJSON.Unmarshal(legacyTimerWire, &legacyTimerRoundTrip); legacyTimerError != nil {
					legacyTimerT.Fatal(legacyTimerError)
				}
				legacyTimerRoundTrip[1][legacyTimerAxis] = "ab35_distinct"
				legacyTimerWire, legacyTimerError = legacyTimerJSON.Marshal(legacyTimerRoundTrip)
				if legacyTimerError != nil {
					legacyTimerT.Fatal(legacyTimerError)
				}
				if legacyTimerError = legacyTimerJSON.Unmarshal(legacyTimerWire, &legacyTimerRoundTrip); legacyTimerError != nil {
					legacyTimerT.Fatal(legacyTimerError)
				}
				if legacyTimerReflect.DeepEqual(legacyTimerRoundTrip[0], legacyTimerRoundTrip[1]) || legacyTimerRoundTrip[1][legacyTimerAxis] != "ab35_distinct" {
					legacyTimerT.Fatal("ab35 identity axis collapsed")
				}
				delete(legacyTimerRoundTrip[1], legacyTimerAxis)
				if legacyTimerReflect.DeepEqual(legacyTimerRoundTrip[0], legacyTimerRoundTrip[1]) {
					legacyTimerT.Fatal("ab35 absent axis equated with explicit null/value")
				}
			})
		}
	}
}

func legacyTimerInvert(legacyTimerMapping map[int]string) (map[string]int, bool) {
	legacyTimerInverse := map[string]int{}
	for legacyTimerRaw, legacyTimerLabel := range legacyTimerMapping {
		if _, legacyTimerPresent := legacyTimerInverse[legacyTimerLabel]; legacyTimerPresent {
			return nil, false
		}
		legacyTimerInverse[legacyTimerLabel] = legacyTimerRaw
	}
	return legacyTimerInverse, true
}

func TestAb35ModeAmbiguity(legacyTimerT *legacyTimerTesting.T) {
	legacyTimerControl := map[int]string{0: "ab35_a", 1: "ab35_b", 2: "ab35_c", 3: "ab35_d"}
	legacyTimerInverse, legacyTimerOK := legacyTimerInvert(legacyTimerControl)
	if !legacyTimerOK || len(legacyTimerInverse) != 4 {
		legacyTimerT.Fatal("ab35 injective control rejected")
	}
	for legacyTimerRaw, legacyTimerLabel := range legacyTimerControl {
		if legacyTimerInverse[legacyTimerLabel] != legacyTimerRaw {
			legacyTimerT.Fatal("ab35 control inverse changed")
		}
	}
	legacyTimerInverse, legacyTimerOK = legacyTimerInvert(map[int]string{0: "ab35_same", 1: "ab35_same", 2: "ab35_same", 3: "ab35_same"})
	if legacyTimerOK || legacyTimerInverse != nil {
		legacyTimerT.Fatal("ab35 ambiguous inversion accepted")
	}
	legacyTimerRows := legacyTimerDecode(legacyTimerT, LegacyTimerCatalog())["entries"].([]any)
	for _, legacyTimerEntry := range legacyTimerRows[6:] {
		legacyTimerRow := legacyTimerEntry.(map[string]any)
		if legacyTimerRow["mapping_injective"] != false || legacyTimerRow["finite_mappings"] != nil || legacyTimerRow["inverse_default"] != nil {
			legacyTimerT.Fatal("ab35 candidate mode table invented or invertible")
		}
	}
}

func TestAb35OrderingCandidate(legacyTimerT *legacyTimerTesting.T) {
	legacyTimerOrder := legacyTimerDecode(legacyTimerT, LegacyTimerCatalog())["ordering_candidate"].(map[string]any)
	if legacyTimerOrder["written_start_equal_end"] != "disabled interval" || legacyTimerOrder["out_of_order_qualification"] != "Unknown" {
		legacyTimerT.Fatal("ab35 ordering promoted")
	}
	for _, legacyTimerExample := range []struct {
		legacyTimerName                  string
		legacyTimerStart, legacyTimerEnd int
		legacyTimerWant                  string
	}{
		{"ab35_equal", 54, 54, "disabled interval"},
		{"ab35_unequal", 54, 60, "Unknown"},
	} {
		legacyTimerT.Run(legacyTimerExample.legacyTimerName, func(legacyTimerT *legacyTimerTesting.T) {
			legacyTimerClaim := "Unknown"
			if legacyTimerExample.legacyTimerStart == legacyTimerExample.legacyTimerEnd {
				legacyTimerClaim = legacyTimerOrder["written_start_equal_end"].(string)
			}
			if legacyTimerClaim != legacyTimerExample.legacyTimerWant {
				legacyTimerT.Fatal("ab35 equality candidate changed")
			}
		})
	}
	for legacyTimerIndex, legacyTimerIntervals := range [][3][2]int{
		{{90, 96}, {54, 60}, {108, 114}},
		{{54, 60}, {108, 114}, {90, 96}},
	} {
		legacyTimerT.Run(legacyTimerFmt.Sprintf("ab35_out_of_order_%d", legacyTimerIndex), func(legacyTimerT *legacyTimerTesting.T) {
			if legacyTimerIntervals[0][0] < legacyTimerIntervals[1][0] && legacyTimerIntervals[1][0] < legacyTimerIntervals[2][0] {
				legacyTimerT.Fatal("ab35 invalid out-of-order fixture")
			}
			if legacyTimerOrder["out_of_order_qualification"] != "Unknown" {
				legacyTimerT.Fatal("ab35 synthetic ordering qualified")
			}
		})
	}
	legacyTimerAssert(legacyTimerT, "ab35_order", legacyTimerOrder["interval_order"], []any{"first", "second", "third"})
}
