package ebusstandard

import (
	standard07ResidualBytes "bytes"
	standard07ResidualJSON "encoding/json"
	standard07ResidualFmt "fmt"
	standard07ResidualReflect "reflect"
	standard07ResidualStrings "strings"
	standard07ResidualTesting "testing"
)

func standard07ResidualTestKey(standard07ResidualSecondary uint8, standard07ResidualDirection, standard07ResidualField string) Standard07ResidualKey {
	return Standard07ResidualKey{Standard07ResidualFamily: "standard", Standard07ResidualVariant: "standard07", Standard07ResidualPrimary: 7, Standard07ResidualSecondary: &standard07ResidualSecondary, Standard07ResidualDirection: standard07ResidualDirection, Standard07ResidualField: standard07ResidualField, Standard07ResidualAlternative: "none", Standard07ResidualInstanceScope: "none"}
}

func standard07ResidualTestRow(standard07ResidualT *standard07ResidualTesting.T, standard07ResidualKey Standard07ResidualKey) map[string]standard07ResidualJSON.RawMessage {
	standard07ResidualT.Helper()
	standard07ResidualRow, standard07ResidualFound, standard07ResidualErr := Standard07ResidualLookup(standard07ResidualKey)
	if standard07ResidualErr != nil || !standard07ResidualFound {
		standard07ResidualT.Fatalf("lookup %+v: found=%v error=%v", standard07ResidualKey, standard07ResidualFound, standard07ResidualErr)
	}
	var standard07ResidualObject map[string]standard07ResidualJSON.RawMessage
	if standard07ResidualErr := standard07ResidualJSON.Unmarshal(standard07ResidualRow, &standard07ResidualObject); standard07ResidualErr != nil {
		standard07ResidualT.Fatal(standard07ResidualErr)
	}
	return standard07ResidualObject
}

func standard07ResidualAssert(standard07ResidualT *standard07ResidualTesting.T, standard07ResidualObject map[string]standard07ResidualJSON.RawMessage, standard07ResidualField, standard07ResidualExpected string) {
	standard07ResidualT.Helper()
	var standard07ResidualActualValue, standard07ResidualExpectedValue any
	if standard07ResidualErr := standard07ResidualJSON.Unmarshal(standard07ResidualObject[standard07ResidualField], &standard07ResidualActualValue); standard07ResidualErr != nil {
		standard07ResidualT.Fatalf("%s: %v", standard07ResidualField, standard07ResidualErr)
	}
	if standard07ResidualErr := standard07ResidualJSON.Unmarshal([]byte(standard07ResidualExpected), &standard07ResidualExpectedValue); standard07ResidualErr != nil {
		standard07ResidualT.Fatal(standard07ResidualErr)
	}
	if !standard07ResidualReflect.DeepEqual(standard07ResidualActualValue, standard07ResidualExpectedValue) {
		standard07ResidualT.Fatalf("%s: got %s, want %s", standard07ResidualField, standard07ResidualObject[standard07ResidualField], standard07ResidualExpected)
	}
}

func TestAb28FiniteInventory(standard07ResidualT *standard07ResidualTesting.T) {
	standard07ResidualData, standard07ResidualErr := Standard07ResidualEntries()
	if standard07ResidualErr != nil {
		standard07ResidualT.Fatal(standard07ResidualErr)
	}
	var standard07ResidualRows []struct {
		Standard07ResidualID  string                `json:"id"`
		Standard07ResidualKey Standard07ResidualKey `json:"key"`
	}
	if standard07ResidualErr := standard07ResidualJSON.Unmarshal(standard07ResidualData, &standard07ResidualRows); standard07ResidualErr != nil {
		standard07ResidualT.Fatal(standard07ResidualErr)
	}
	if len(standard07ResidualRows) != 37 {
		standard07ResidualT.Fatalf("finite inventory: %d", len(standard07ResidualRows))
	}
	standard07ResidualSeen := make(map[string]bool)
	for _, standard07ResidualRow := range standard07ResidualRows {
		if !standard07ResidualStrings.HasPrefix(standard07ResidualRow.Standard07ResidualID, "ab28-") || standard07ResidualSeen[standard07ResidualRow.Standard07ResidualID] {
			standard07ResidualT.Fatalf("invalid or repeated ID: %s", standard07ResidualRow.Standard07ResidualID)
		}
		standard07ResidualSeen[standard07ResidualRow.Standard07ResidualID] = true
		if standard07ResidualRow.Standard07ResidualKey.Standard07ResidualSecondary == nil || !standard07ResidualBytes.Contains([]byte{0, 2, 3, 5, 254, 255}, []byte{*standard07ResidualRow.Standard07ResidualKey.Standard07ResidualSecondary}) {
			standard07ResidualT.Fatal("invented or out-of-slice service")
		}
		standard07ResidualObject := standard07ResidualTestRow(standard07ResidualT, standard07ResidualRow.Standard07ResidualKey)
		standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, "grade", `"Hypothesis"`)
		standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, "unknown", `{"grade":"Unknown","native_support":"Unknown","physical_effects":"Unknown","persistence":"Unknown","delivery":"Unknown","qualification":"Unknown"}`)
		standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, "id", standard07ResidualFmt.Sprintf("%q", standard07ResidualRow.Standard07ResidualID))
		var standard07ResidualRaw map[string]standard07ResidualJSON.RawMessage
		if standard07ResidualErr := standard07ResidualJSON.Unmarshal(standard07ResidualObject["raw"], &standard07ResidualRaw); standard07ResidualErr != nil {
			standard07ResidualT.Fatal(standard07ResidualErr)
		}
		standard07ResidualAssert(standard07ResidualT, standard07ResidualRaw, "retained", `true`)
		standard07ResidualAssert(standard07ResidualT, standard07ResidualRaw, "coordinates", `"after PB/SB"`)
	}
	for _, standard07ResidualSecondary := range []uint8{1, 4} {
		if _, standard07ResidualFound, standard07ResidualErr := Standard07ResidualLookup(standard07ResidualTestKey(standard07ResidualSecondary, "request", "temperature")); standard07ResidualErr != nil || standard07ResidualFound {
			standard07ResidualT.Fatal("out-of-slice body admitted")
		}
	}
}

func TestAb28ResidualFields(standard07ResidualT *standard07ResidualTesting.T) {
	standard07ResidualWeekday := standard07ResidualTestRow(standard07ResidualT, standard07ResidualTestKey(0, "request", "weekday"))
	standard07ResidualAssert(standard07ResidualT, standard07ResidualWeekday, "encoding", `"raw u8"`)
	standard07ResidualAssert(standard07ResidualT, standard07ResidualWeekday, "offset", `7`)
	standard07ResidualAssert(standard07ResidualT, standard07ResidualWeekday, "width", `1`)
	standard07ResidualAssert(standard07ResidualT, standard07ResidualWeekday, "labels", `{"0":"Monday","1":"Tuesday","2":"Wednesday","3":"Thursday","4":"Friday","5":"Saturday","6":"Sunday"}`)
	standard07ResidualAssert(standard07ResidualT, standard07ResidualWeekday, "unlabelled_minimum", `7`)
	standard07ResidualAssert(standard07ResidualT, standard07ResidualWeekday, "unlabelled_maximum", `255`)
	standard07ResidualAssert(standard07ResidualT, standard07ResidualWeekday, "examples", `[{"raw":[6],"label":"Sunday"},{"raw":[7],"label":null},{"raw":[255],"label":null}]`)
	for _, standard07ResidualCase := range []struct {
		Standard07ResidualField                                                        string
		Standard07ResidualOffset, Standard07ResidualMinimum, Standard07ResidualMaximum int
	}{{"second", 2, 0, 59}, {"minute", 3, 0, 59}, {"hour", 4, 0, 23}, {"day", 5, 1, 31}, {"month", 6, 1, 12}, {"year", 8, 0, 99}} {
		standard07ResidualObject := standard07ResidualTestRow(standard07ResidualT, standard07ResidualTestKey(0, "request", standard07ResidualCase.Standard07ResidualField))
		standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, "encoding", `"packed-decimal"`)
		standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, "width", `1`)
		for standard07ResidualField, standard07ResidualValue := range map[string]int{"offset": standard07ResidualCase.Standard07ResidualOffset, "minimum": standard07ResidualCase.Standard07ResidualMinimum, "maximum": standard07ResidualCase.Standard07ResidualMaximum} {
			standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, standard07ResidualField, standard07ResidualFmt.Sprint(standard07ResidualValue))
		}
		for _, standard07ResidualField := range []string{"century", "timezone", "calendar_validation", "complete_native_decimal_conversion"} {
			standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, standard07ResidualField, `"Unknown"`)
		}
		standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, "executable_schema", `false`)
	}
	for _, standard07ResidualSecondary := range []uint8{0, 2} {
		standard07ResidualField := "outdoor_temperature"
		standard07ResidualObject := standard07ResidualTestRow(standard07ResidualT, standard07ResidualTestKey(standard07ResidualSecondary, "request", standard07ResidualField))
		for standard07ResidualField, standard07ResidualExpected := range map[string]string{"offset": "0", "width": "2", "encoding": `"signed LE16 /256"`, "unit": `"degrees Celsius"`, "raw_minimum": "-12800", "raw_maximum": "12800", "minimum": "-50", "maximum": "50", "sentinel_8000": `"Unknown"`, "examples": `[{"raw":[0,50],"value":50},{"raw":[0,206],"value":-50}]`} {
			standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, standard07ResidualField, standard07ResidualExpected)
		}
	}
	standard07ResidualDuration := standard07ResidualTestRow(standard07ResidualT, standard07ResidualTestKey(2, "request", "duration"))
	for standard07ResidualField, standard07ResidualExpected := range map[string]string{"offset": "2", "width": "1", "encoding": `"raw u8"`, "unit": `"minutes"`, "default": "0", "zero_permanent_interpretation": `"Unknown"`, "examples": `[{"raw":[0],"value":0}]`} {
		standard07ResidualAssert(standard07ResidualT, standard07ResidualDuration, standard07ResidualField, standard07ResidualExpected)
	}
	for _, standard07ResidualCase := range []struct {
		Standard07ResidualSecondary uint8
		Standard07ResidualLiteral   int
	}{{0, 6}, {2, 1}} {
		standard07ResidualObject := standard07ResidualTestRow(standard07ResidualT, standard07ResidualTestKey(standard07ResidualCase.Standard07ResidualSecondary, "none", "unassociated_default"))
		standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, "literal", standard07ResidualFmt.Sprint(standard07ResidualCase.Standard07ResidualLiteral))
		standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, "field_association", `"Unknown"`)
	}
	for _, standard07ResidualKey := range []Standard07ResidualKey{standard07ResidualTestKey(0, "request", "weekday"), standard07ResidualTestKey(0, "request", "outdoor_temperature"), standard07ResidualTestKey(2, "request", "outdoor_temperature")} {
		if _, standard07ResidualOK := standard07ResidualTestRow(standard07ResidualT, standard07ResidualKey)["default"]; standard07ResidualOK {
			standard07ResidualT.Fatal("unassociated default attached to field")
		}
	}
}

func TestAb28BlocksAndPairs(standard07ResidualT *standard07ResidualTesting.T) {
	for _, standard07ResidualCase := range []struct {
		Standard07ResidualSecondary                       uint8
		Standard07ResidualField                           string
		Standard07ResidualOffset, Standard07ResidualWidth int
	}{{3, "block", 0, 1}, {5, "secondary_block", 0, 2}, {5, "primary_block", 1, 2}} {
		standard07ResidualObject := standard07ResidualTestRow(standard07ResidualT, standard07ResidualTestKey(standard07ResidualCase.Standard07ResidualSecondary, "request", standard07ResidualCase.Standard07ResidualField))
		for standard07ResidualField, standard07ResidualExpected := range map[string]string{"offset": standard07ResidualFmt.Sprint(standard07ResidualCase.Standard07ResidualOffset), "width": "1", "encoding": `"raw u8"`, "request_body_width": standard07ResidualFmt.Sprint(standard07ResidualCase.Standard07ResidualWidth)} {
			standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, standard07ResidualField, standard07ResidualExpected)
		}
	}
	for _, standard07ResidualSecondary := range []uint8{3, 5} {
		for _, standard07ResidualCase := range []struct {
			Standard07ResidualField    string
			Standard07ResidualOffset   int
			Standard07ResidualExamples string
		}{{"version", 0, `[{"raw":[18],"value":12}]`}, {"revision", 1, `[{"raw":[52],"value":34}]`}} {
			standard07ResidualObject := standard07ResidualTestRow(standard07ResidualT, standard07ResidualTestKey(standard07ResidualSecondary, "positive", standard07ResidualCase.Standard07ResidualField))
			for standard07ResidualField, standard07ResidualExpected := range map[string]string{"offset": standard07ResidualFmt.Sprint(standard07ResidualCase.Standard07ResidualOffset), "width": "1", "encoding": `"packed-decimal"`, "independent": "true", "combined_api_number": "false", "positive_body_width": "10", "version_control_overlap_0704": `"Unknown"`, "complete_native_decimal_conversion": `"Unknown"`, "executable_schema": "false", "examples": standard07ResidualCase.Standard07ResidualExamples} {
				standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, standard07ResidualField, standard07ResidualExpected)
			}
		}
		for standard07ResidualIndex := 0; standard07ResidualIndex < 8; standard07ResidualIndex++ {
			standard07ResidualObject := standard07ResidualTestRow(standard07ResidualT, standard07ResidualTestKey(standard07ResidualSecondary, "positive", standard07ResidualFmt.Sprintf("bitmap_%d", standard07ResidualIndex)))
			for standard07ResidualField, standard07ResidualExpected := range map[string]string{"offset": standard07ResidualFmt.Sprint(standard07ResidualIndex + 2), "width": "1", "encoding": `"raw bitmap u8"`, "independent": "true", "bit_selection_algorithm": `"Unknown"`, "admitted_methods": "[]", "examples": `[{"raw":[128],"value":128,"selected_service":"Unknown"}]`} {
				standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, standard07ResidualField, standard07ResidualExpected)
			}
		}
	}
	for _, standard07ResidualCase := range []struct {
		Standard07ResidualSecondary uint8
		Standard07ResidualField     string
	}{{254, "inquiry"}, {255, "sign_of_life"}} {
		standard07ResidualObject := standard07ResidualTestRow(standard07ResidualT, standard07ResidualTestKey(standard07ResidualCase.Standard07ResidualSecondary, "request", standard07ResidualCase.Standard07ResidualField))
		for standard07ResidualField, standard07ResidualExpected := range map[string]string{"argument_body_width": "0", "emission": `"Unknown"`, "physical_response": `"Unknown"`, "direct_inquiry_reply": `"Unknown"`, "raw": standard07ResidualFmt.Sprintf(`{"retained":true,"coordinates":"after PB/SB","pair":[7,%d],"body":[]}`, standard07ResidualCase.Standard07ResidualSecondary)} {
			standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, standard07ResidualField, standard07ResidualExpected)
		}
	}
}

func TestAb28IndependentSelectors(standard07ResidualT *standard07ResidualTesting.T) {
	for _, standard07ResidualCase := range []struct {
		Standard07ResidualSecondary                         uint8
		Standard07ResidualField, Standard07ResidualExamples string
	}{
		{3, "block", `[{"raw":[1],"value":1},{"raw":[2],"value":2}]`},
		{5, "secondary_block", `[{"raw":[1,2],"value":1},{"raw":[3,2],"value":3}]`},
		{5, "primary_block", `[{"raw":[1,2],"value":2},{"raw":[1,3],"value":3}]`},
	} {
		standard07ResidualAssert(standard07ResidualT, standard07ResidualTestRow(standard07ResidualT, standard07ResidualTestKey(standard07ResidualCase.Standard07ResidualSecondary, "request", standard07ResidualCase.Standard07ResidualField)), "examples", standard07ResidualCase.Standard07ResidualExamples)
	}
}

func TestAb28ExactAxes(standard07ResidualT *standard07ResidualTesting.T) {
	standard07ResidualZero8 := uint8(0)
	standard07ResidualZero16 := uint16(0)
	standard07ResidualOne8 := uint8(1)
	standard07ResidualChanges := []struct {
		Standard07ResidualName   string
		Standard07ResidualChange func(*Standard07ResidualKey)
	}{
		{"family", func(standard07ResidualKey *Standard07ResidualKey) {
			standard07ResidualKey.Standard07ResidualFamily = "synthetic"
		}},
		{"variant", func(standard07ResidualKey *Standard07ResidualKey) {
			standard07ResidualKey.Standard07ResidualVariant = "synthetic"
		}},
		{"primary", func(standard07ResidualKey *Standard07ResidualKey) {
			standard07ResidualKey.Standard07ResidualPrimary = 8
		}},
		{"secondary", func(standard07ResidualKey *Standard07ResidualKey) {
			standard07ResidualKey.Standard07ResidualSecondary = &standard07ResidualOne8
		}},
		{"operation", func(standard07ResidualKey *Standard07ResidualKey) {
			standard07ResidualKey.Standard07ResidualOperation = &standard07ResidualZero8
		}},
		{"direction", func(standard07ResidualKey *Standard07ResidualKey) {
			standard07ResidualKey.Standard07ResidualDirection = "positive"
		}},
		{"field", func(standard07ResidualKey *Standard07ResidualKey) {
			standard07ResidualKey.Standard07ResidualField = "synthetic"
		}},
		{"alternative", func(standard07ResidualKey *Standard07ResidualKey) {
			standard07ResidualKey.Standard07ResidualAlternative = "synthetic"
		}},
		{"group", func(standard07ResidualKey *Standard07ResidualKey) {
			standard07ResidualKey.Standard07ResidualGroup = &standard07ResidualZero8
		}},
		{"instance_scope", func(standard07ResidualKey *Standard07ResidualKey) {
			standard07ResidualKey.Standard07ResidualInstanceScope = "dynamic"
		}},
		{"register", func(standard07ResidualKey *Standard07ResidualKey) {
			standard07ResidualKey.Standard07ResidualRegister = &standard07ResidualZero16
		}},
		{"secondary_absent_zero", func(standard07ResidualKey *Standard07ResidualKey) {
			standard07ResidualKey.Standard07ResidualSecondary = nil
		}},
	}
	for _, standard07ResidualCase := range standard07ResidualChanges {
		standard07ResidualT.Run("ab28_"+standard07ResidualCase.Standard07ResidualName, func(standard07ResidualT *standard07ResidualTesting.T) {
			standard07ResidualBase := standard07ResidualTestKey(0, "request", "outdoor_temperature")
			standard07ResidualChanged := standard07ResidualBase
			standard07ResidualCase.Standard07ResidualChange(&standard07ResidualChanged)
			standard07ResidualOriginal, standard07ResidualFound, standard07ResidualErr := Standard07ResidualLookup(standard07ResidualBase)
			if standard07ResidualErr != nil || !standard07ResidualFound {
				standard07ResidualT.Fatalf("base control: %v %v", standard07ResidualFound, standard07ResidualErr)
			}
			if standard07ResidualRow, standard07ResidualFound, standard07ResidualErr := Standard07ResidualLookup(standard07ResidualChanged); standard07ResidualErr != nil || standard07ResidualFound || standard07ResidualRow != nil {
				standard07ResidualT.Fatalf("non-exact key matched: %s %v %v", standard07ResidualRow, standard07ResidualFound, standard07ResidualErr)
			}
			var standard07ResidualCompact standard07ResidualBytes.Buffer
			if standard07ResidualErr := standard07ResidualJSON.Compact(&standard07ResidualCompact, standard07ResidualOriginal); standard07ResidualErr != nil {
				standard07ResidualT.Fatal(standard07ResidualErr)
			}
			standard07ResidualOriginal = standard07ResidualCompact.Bytes()
			var standard07ResidualObject map[string]standard07ResidualJSON.RawMessage
			if standard07ResidualErr := standard07ResidualJSON.Unmarshal(standard07ResidualOriginal, &standard07ResidualObject); standard07ResidualErr != nil {
				standard07ResidualT.Fatal(standard07ResidualErr)
			}
			standard07ResidualChangedJSON, standard07ResidualErr := standard07ResidualJSON.Marshal(standard07ResidualChanged)
			if standard07ResidualErr != nil {
				standard07ResidualT.Fatal(standard07ResidualErr)
			}
			standard07ResidualObject["key"] = standard07ResidualChangedJSON
			standard07ResidualObject["id"] = standard07ResidualJSON.RawMessage(`"ab28-synthetic-control"`)
			standard07ResidualOther, standard07ResidualErr := standard07ResidualJSON.Marshal(standard07ResidualObject)
			if standard07ResidualErr != nil {
				standard07ResidualT.Fatal(standard07ResidualErr)
			}
			standard07ResidualControls, standard07ResidualErr := standard07ResidualJSON.Marshal([]standard07ResidualJSON.RawMessage{standard07ResidualOriginal, standard07ResidualOther})
			if standard07ResidualErr != nil {
				standard07ResidualT.Fatal(standard07ResidualErr)
			}
			for _, standard07ResidualControl := range []struct {
				Standard07ResidualKey      Standard07ResidualKey
				Standard07ResidualExpected standard07ResidualJSON.RawMessage
			}{{standard07ResidualBase, standard07ResidualOriginal}, {standard07ResidualChanged, standard07ResidualOther}} {
				standard07ResidualRow, standard07ResidualFound, standard07ResidualErr := standard07ResidualLookup(standard07ResidualControls, standard07ResidualControl.Standard07ResidualKey)
				if standard07ResidualErr != nil || !standard07ResidualFound || !standard07ResidualBytes.Equal(standard07ResidualRow, standard07ResidualControl.Standard07ResidualExpected) {
					standard07ResidualT.Fatalf("distinct valid control: %s %v %v", standard07ResidualRow, standard07ResidualFound, standard07ResidualErr)
				}
			}
		})
	}
	// ab28: independent nonzero optional controls prevent presence-only matching.
	for _, standard07ResidualAxis := range []string{"secondary", "operation", "group", "register"} {
		standard07ResidualT.Run("ab28_"+standard07ResidualAxis+"_value", func(standard07ResidualT *standard07ResidualTesting.T) {
			standard07ResidualBase := standard07ResidualTestKey(0, "request", "outdoor_temperature")
			standard07ResidualBase.Standard07ResidualOperation = &standard07ResidualZero8
			standard07ResidualBase.Standard07ResidualGroup = &standard07ResidualZero8
			standard07ResidualBase.Standard07ResidualRegister = &standard07ResidualZero16
			standard07ResidualChanged := standard07ResidualBase
			standard07ResidualOne16 := uint16(1)
			switch standard07ResidualAxis {
			case "secondary":
				standard07ResidualChanged.Standard07ResidualSecondary = &standard07ResidualOne8
			case "operation":
				standard07ResidualChanged.Standard07ResidualOperation = &standard07ResidualOne8
			case "group":
				standard07ResidualChanged.Standard07ResidualGroup = &standard07ResidualOne8
			case "register":
				standard07ResidualChanged.Standard07ResidualRegister = &standard07ResidualOne16
			}
			standard07ResidualControls, standard07ResidualErr := standard07ResidualJSON.Marshal([]map[string]any{{"id": "ab28-zero", "key": standard07ResidualBase}, {"id": "ab28-one", "key": standard07ResidualChanged}})
			if standard07ResidualErr != nil {
				standard07ResidualT.Fatal(standard07ResidualErr)
			}
			for _, standard07ResidualControl := range []struct {
				Standard07ResidualKey Standard07ResidualKey
				Standard07ResidualID  string
			}{{standard07ResidualBase, "ab28-zero"}, {standard07ResidualChanged, "ab28-one"}} {
				standard07ResidualRow, standard07ResidualFound, standard07ResidualErr := standard07ResidualLookup(standard07ResidualControls, standard07ResidualControl.Standard07ResidualKey)
				if standard07ResidualErr != nil || !standard07ResidualFound {
					standard07ResidualT.Fatalf("optional value: %v %v", standard07ResidualFound, standard07ResidualErr)
				}
				var standard07ResidualObject map[string]standard07ResidualJSON.RawMessage
				if standard07ResidualErr := standard07ResidualJSON.Unmarshal(standard07ResidualRow, &standard07ResidualObject); standard07ResidualErr != nil {
					standard07ResidualT.Fatal(standard07ResidualErr)
				}
				standard07ResidualAssert(standard07ResidualT, standard07ResidualObject, "id", standard07ResidualFmt.Sprintf("%q", standard07ResidualControl.Standard07ResidualID))
			}
		})
	}
}

func TestAb28ValidationAndCopies(standard07ResidualT *standard07ResidualTesting.T) {
	standard07ResidualKey := standard07ResidualTestKey(0, "request", "outdoor_temperature")
	standard07ResidualRow, standard07ResidualFound, standard07ResidualErr := Standard07ResidualLookup(standard07ResidualKey)
	if standard07ResidualErr != nil || !standard07ResidualFound {
		standard07ResidualT.Fatal(standard07ResidualErr)
	}
	standard07ResidualDuplicate, standard07ResidualErr := standard07ResidualJSON.Marshal([]standard07ResidualJSON.RawMessage{standard07ResidualRow, standard07ResidualRow})
	if standard07ResidualErr != nil {
		standard07ResidualT.Fatal(standard07ResidualErr)
	}
	for _, standard07ResidualBad := range [][]byte{[]byte(`{`), []byte(`{}`), []byte(`null`), []byte(`[null]`), []byte(`[{}]`), []byte(`[{"key":{"primary":256}}]`), standard07ResidualDuplicate} {
		if _, standard07ResidualErr := standard07ResidualValidate(standard07ResidualBad); standard07ResidualErr == nil {
			standard07ResidualT.Fatalf("invalid data admitted: %s", standard07ResidualBad)
		}
		if standard07ResidualResult, standard07ResidualFound, standard07ResidualErr := standard07ResidualLookup(standard07ResidualBad, standard07ResidualKey); standard07ResidualErr == nil || standard07ResidualFound || standard07ResidualResult != nil {
			standard07ResidualT.Fatal("invalid catalog did not fail closed")
		}
	}
	var standard07ResidualObject map[string]standard07ResidualJSON.RawMessage
	if standard07ResidualErr := standard07ResidualJSON.Unmarshal(standard07ResidualRow, &standard07ResidualObject); standard07ResidualErr != nil {
		standard07ResidualT.Fatal(standard07ResidualErr)
	}
	var standard07ResidualAxes map[string]standard07ResidualJSON.RawMessage
	if standard07ResidualErr := standard07ResidualJSON.Unmarshal(standard07ResidualObject["key"], &standard07ResidualAxes); standard07ResidualErr != nil {
		standard07ResidualT.Fatal(standard07ResidualErr)
	}
	for standard07ResidualAxis, standard07ResidualValue := range standard07ResidualAxes {
		delete(standard07ResidualAxes, standard07ResidualAxis)
		standard07ResidualKeyJSON, standard07ResidualErr := standard07ResidualJSON.Marshal(standard07ResidualAxes)
		if standard07ResidualErr != nil {
			standard07ResidualT.Fatal(standard07ResidualErr)
		}
		standard07ResidualObject["key"] = standard07ResidualKeyJSON
		standard07ResidualBad, standard07ResidualErr := standard07ResidualJSON.Marshal([]map[string]standard07ResidualJSON.RawMessage{standard07ResidualObject})
		if standard07ResidualErr != nil {
			standard07ResidualT.Fatal(standard07ResidualErr)
		}
		if _, standard07ResidualErr := standard07ResidualValidate(standard07ResidualBad); standard07ResidualErr == nil {
			standard07ResidualT.Fatalf("missing %s admitted", standard07ResidualAxis)
		}
		standard07ResidualAxes[standard07ResidualAxis] = standard07ResidualValue
	}
	standard07ResidualEntries, standard07ResidualErr := Standard07ResidualEntries()
	if standard07ResidualErr != nil {
		standard07ResidualT.Fatal(standard07ResidualErr)
	}
	standard07ResidualEntriesAgain, standard07ResidualErr := Standard07ResidualEntries()
	if standard07ResidualErr != nil {
		standard07ResidualT.Fatal(standard07ResidualErr)
	}
	standard07ResidualRowAgain, standard07ResidualFound, standard07ResidualErr := Standard07ResidualLookup(standard07ResidualKey)
	if standard07ResidualErr != nil || !standard07ResidualFound {
		standard07ResidualT.Fatal(standard07ResidualErr)
	}
	standard07ResidualExpectedRow := append([]byte(nil), standard07ResidualRow...)
	standard07ResidualExpectedEntries := append([]byte(nil), standard07ResidualEntries...)
	for standard07ResidualIndex := range standard07ResidualRow {
		standard07ResidualRow[standard07ResidualIndex] = '!'
	}
	for standard07ResidualIndex := range standard07ResidualEntries {
		standard07ResidualEntries[standard07ResidualIndex] = '!'
	}
	if !standard07ResidualBytes.Equal(standard07ResidualRowAgain, standard07ResidualExpectedRow) || !standard07ResidualBytes.Equal(standard07ResidualEntriesAgain, standard07ResidualExpectedEntries) {
		standard07ResidualT.Fatal("simultaneous results share memory")
	}
	standard07ResidualLater, standard07ResidualFound, standard07ResidualErr := Standard07ResidualLookup(standard07ResidualKey)
	if standard07ResidualErr != nil || !standard07ResidualFound || !standard07ResidualBytes.Equal(standard07ResidualLater, standard07ResidualExpectedRow) {
		standard07ResidualT.Fatal("row mutation affected later lookup")
	}
	standard07ResidualLaterEntries, standard07ResidualErr := Standard07ResidualEntries()
	if standard07ResidualErr != nil || !standard07ResidualBytes.Equal(standard07ResidualLaterEntries, standard07ResidualExpectedEntries) {
		standard07ResidualT.Fatal("buffer mutation affected later entries")
	}
	standard07ResidualSource := append([]byte{'['}, standard07ResidualExpectedRow...)
	standard07ResidualSource = append(standard07ResidualSource, ']')
	standard07ResidualIsolated, standard07ResidualFound, standard07ResidualErr := standard07ResidualLookup(standard07ResidualSource, standard07ResidualKey)
	if standard07ResidualErr != nil || !standard07ResidualFound {
		standard07ResidualT.Fatal(standard07ResidualErr)
	}
	for standard07ResidualIndex := range standard07ResidualSource {
		standard07ResidualSource[standard07ResidualIndex] = '!'
	}
	if !standard07ResidualBytes.Equal(standard07ResidualIsolated, standard07ResidualExpectedRow) {
		standard07ResidualT.Fatal("lookup aliases input")
	}
}
