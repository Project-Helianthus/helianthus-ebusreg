package ebusstandard

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

type standard07FieldViewCase struct {
	Standard07FieldViewID       string                    `json:"id"`
	Standard07FieldViewInput    json.RawMessage           `json:"input"`
	Standard07FieldViewExpected map[string]map[string]any `json:"expected"`
}

func standard07FieldViewCases(standard07FieldViewT *testing.T) []standard07FieldViewCase {
	standard07FieldViewT.Helper()
	standard07FieldViewData, standard07FieldViewErr := os.ReadFile("testdata/ab15aq15_standard07_field_views.json")
	if standard07FieldViewErr != nil {
		standard07FieldViewT.Fatal(standard07FieldViewErr)
	}
	var standard07FieldViewCases []standard07FieldViewCase
	if standard07FieldViewErr := json.Unmarshal(standard07FieldViewData, &standard07FieldViewCases); standard07FieldViewErr != nil {
		standard07FieldViewT.Fatal(standard07FieldViewErr)
	}
	if len(standard07FieldViewCases) != 12 {
		standard07FieldViewT.Fatal("fixture count", len(standard07FieldViewCases))
	}
	return standard07FieldViewCases
}

func standard07FieldViewEvaluate(standard07FieldViewT *testing.T, standard07FieldViewInput any) map[string]any {
	standard07FieldViewT.Helper()
	standard07FieldViewData, standard07FieldViewErr := json.Marshal(standard07FieldViewInput)
	if standard07FieldViewErr != nil {
		standard07FieldViewT.Fatal(standard07FieldViewErr)
	}
	return standard07FieldViewDecode(standard07FieldViewT, standard07FieldViewData)
}

func standard07FieldViewDecode(standard07FieldViewT *testing.T, standard07FieldViewInput []byte) map[string]any {
	standard07FieldViewT.Helper()
	standard07FieldViewData, standard07FieldViewErr := Standard07FieldViewEvaluate(standard07FieldViewInput)
	if standard07FieldViewErr != nil {
		standard07FieldViewT.Fatal(standard07FieldViewErr)
	}
	var standard07FieldViewOutput map[string]any
	if standard07FieldViewErr := json.Unmarshal(standard07FieldViewData, &standard07FieldViewOutput); standard07FieldViewErr != nil {
		standard07FieldViewT.Fatal(standard07FieldViewErr)
	}
	return standard07FieldViewOutput
}

func standard07FieldViewInput(standard07FieldViewT *testing.T, standard07FieldViewRaw []byte) map[string]any {
	standard07FieldViewT.Helper()
	var standard07FieldViewInput map[string]any
	if standard07FieldViewErr := json.Unmarshal(standard07FieldViewRaw, &standard07FieldViewInput); standard07FieldViewErr != nil {
		standard07FieldViewT.Fatal(standard07FieldViewErr)
	}
	return standard07FieldViewInput
}

func TestAb15Aq15Fixture(standard07FieldViewT *testing.T) {
	for _, standard07FieldViewCase := range standard07FieldViewCases(standard07FieldViewT) {
		standard07FieldViewT.Run(standard07FieldViewCase.Standard07FieldViewID, func(standard07FieldViewT *testing.T) {
			standard07FieldViewInput := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCase.Standard07FieldViewInput)
			standard07FieldViewOutput := standard07FieldViewDecode(standard07FieldViewT, standard07FieldViewCase.Standard07FieldViewInput)
			if standard07FieldViewOutput["state"] != "complete" || standard07FieldViewOutput["physical"] != "Unknown" || standard07FieldViewOutput["qualified"] != false {
				standard07FieldViewT.Fatal(standard07FieldViewOutput)
			}
			if !reflect.DeepEqual(standard07FieldViewInput["context"], standard07FieldViewOutput["context"]) || !reflect.DeepEqual(standard07FieldViewInput["fields"], standard07FieldViewOutput["fields"]) {
				standard07FieldViewT.Fatal("raw input lost")
			}
			for standard07FieldViewKey := range standard07FieldViewOutput {
				if _, standard07FieldViewPresent := standard07FieldViewInput[standard07FieldViewKey]; !standard07FieldViewPresent && standard07FieldViewKey != "state" && standard07FieldViewKey != "results" && standard07FieldViewKey != "physical" && standard07FieldViewKey != "qualified" {
					standard07FieldViewT.Fatal("invented output", standard07FieldViewKey)
				}
			}
			standard07FieldViewResults := standard07FieldViewOutput["results"].(map[string]any)
			for standard07FieldViewKey, standard07FieldViewExpected := range standard07FieldViewCase.Standard07FieldViewExpected {
				standard07FieldViewResult := standard07FieldViewResults[standard07FieldViewKey].(map[string]any)
				if standard07FieldViewResult["state"] != "complete" || standard07FieldViewResult["grade"] != "Hypothesis" || !reflect.DeepEqual(standard07FieldViewResult["raw"], standard07FieldViewInput["fields"].(map[string]any)[standard07FieldViewKey]) {
					standard07FieldViewT.Fatal(standard07FieldViewResult)
				}
				for standard07FieldViewProperty, standard07FieldViewValue := range standard07FieldViewExpected {
					if !reflect.DeepEqual(standard07FieldViewResult[standard07FieldViewProperty], standard07FieldViewValue) {
						standard07FieldViewT.Fatalf("%s.%s got %v want %v", standard07FieldViewKey, standard07FieldViewProperty, standard07FieldViewResult[standard07FieldViewProperty], standard07FieldViewValue)
					}
				}
				standard07FieldViewAllowed := map[string]bool{"raw": true, "state": true, "grade": true}
				switch standard07FieldViewKey {
				case "outdoor_temperature":
					for _, standard07FieldViewProperty := range []string{"number", "unit", "calculation", "candidate_domain"} {
						standard07FieldViewAllowed[standard07FieldViewProperty] = true
					}
				case "duration":
					for _, standard07FieldViewProperty := range []string{"number", "unit", "calculation"} {
						standard07FieldViewAllowed[standard07FieldViewProperty] = true
					}
				case "weekday":
					standard07FieldViewAllowed["label"] = true
				case "second", "minute", "hour", "day", "month", "year":
					standard07FieldViewAllowed["candidate_domain"] = true
					standard07FieldViewAllowed["conversion"] = true
				}
				for standard07FieldViewProperty := range standard07FieldViewResult {
					if !standard07FieldViewAllowed[standard07FieldViewProperty] {
						standard07FieldViewT.Fatalf("invented property %s", standard07FieldViewProperty)
					}
				}
			}
		})
	}
}

func TestAb15Aq15Axes(standard07FieldViewT *testing.T) {
	for _, standard07FieldViewCase := range standard07FieldViewCases(standard07FieldViewT) {
		for standard07FieldViewKey, standard07FieldViewValue := range map[string]any{"primary": 8, "secondary": 1, "operation": 0, "direction": "negative", "view": "other"} {
			standard07FieldViewT.Run(standard07FieldViewCase.Standard07FieldViewID+"_"+standard07FieldViewKey, func(standard07FieldViewT *testing.T) {
				standard07FieldViewInput := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCase.Standard07FieldViewInput)
				standard07FieldViewInput["context"].(map[string]any)[standard07FieldViewKey] = standard07FieldViewValue
				standard07FieldViewOutput := standard07FieldViewEvaluate(standard07FieldViewT, standard07FieldViewInput)
				if standard07FieldViewOutput["state"] != "Unknown" {
					standard07FieldViewT.Fatal(standard07FieldViewOutput)
				}
				for _, standard07FieldViewValue := range standard07FieldViewOutput["results"].(map[string]any) {
					standard07FieldViewResult := standard07FieldViewValue.(map[string]any)
					if len(standard07FieldViewResult) != 2 || standard07FieldViewResult["state"] != "Unknown" {
						standard07FieldViewT.Fatal(standard07FieldViewResult)
					}
				}
			})
		}
		for _, standard07FieldViewKey := range []string{"family", "variant"} {
			for _, standard07FieldViewValue := range []any{nil, false, float64(0), "", "other-family"} {
				standard07FieldViewInput := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCase.Standard07FieldViewInput)
				standard07FieldViewInput["context"].(map[string]any)[standard07FieldViewKey] = standard07FieldViewValue
				standard07FieldViewOutput := standard07FieldViewEvaluate(standard07FieldViewT, standard07FieldViewInput)
				if standard07FieldViewOutput["qualified"] != false || standard07FieldViewOutput["physical"] != "Unknown" || !reflect.DeepEqual(standard07FieldViewOutput["results"], standard07FieldViewDecode(standard07FieldViewT, standard07FieldViewCase.Standard07FieldViewInput)["results"]) || !reflect.DeepEqual(standard07FieldViewOutput["context"], standard07FieldViewInput["context"]) {
					standard07FieldViewT.Fatal(standard07FieldViewOutput)
				}
			}
		}
	}
}

func TestAb15Aq15WidthsAndIndependence(standard07FieldViewT *testing.T) {
	for _, standard07FieldViewCase := range standard07FieldViewCases(standard07FieldViewT) {
		standard07FieldViewOriginal := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCase.Standard07FieldViewInput)
		for standard07FieldViewKey, standard07FieldViewRaw := range standard07FieldViewOriginal["fields"].(map[string]any) {
			for standard07FieldViewIndex, standard07FieldViewReplacement := range []any{nil, []any{}, append(standard07FieldViewRaw.([]any), float64(0)), standard07FieldViewRaw.([]any)[:len(standard07FieldViewRaw.([]any))-1]} {
				standard07FieldViewInput := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCase.Standard07FieldViewInput)
				standard07FieldViewFields := standard07FieldViewInput["fields"].(map[string]any)
				if standard07FieldViewIndex == 0 {
					delete(standard07FieldViewFields, standard07FieldViewKey)
				} else {
					standard07FieldViewFields[standard07FieldViewKey] = standard07FieldViewReplacement
				}
				standard07FieldViewOutput := standard07FieldViewEvaluate(standard07FieldViewT, standard07FieldViewInput)
				standard07FieldViewResults := standard07FieldViewOutput["results"].(map[string]any)
				standard07FieldViewResult := standard07FieldViewResults[standard07FieldViewKey].(map[string]any)
				if standard07FieldViewOutput["state"] != "partial" || standard07FieldViewResult["state"] != "incomplete" || standard07FieldViewResult["number"] != nil {
					standard07FieldViewT.Fatal(standard07FieldViewOutput)
				}
				if standard07FieldViewIndex == 0 {
					if standard07FieldViewResult["absent"] != true {
						standard07FieldViewT.Fatal(standard07FieldViewResult)
					}
					if _, standard07FieldViewFound := standard07FieldViewResult["raw"]; standard07FieldViewFound {
						standard07FieldViewT.Fatal("invented raw")
					}
				} else if !reflect.DeepEqual(standard07FieldViewResult["raw"], standard07FieldViewReplacement) {
					standard07FieldViewT.Fatal("raw changed")
				}
				standard07FieldViewControl := standard07FieldViewDecode(standard07FieldViewT, standard07FieldViewCase.Standard07FieldViewInput)["results"].(map[string]any)
				for standard07FieldViewSibling, standard07FieldViewValue := range standard07FieldViewControl {
					if standard07FieldViewSibling != standard07FieldViewKey && !reflect.DeepEqual(standard07FieldViewResults[standard07FieldViewSibling], standard07FieldViewValue) {
						standard07FieldViewT.Fatal("sibling lost")
					}
				}
			}
			standard07FieldViewInput := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCase.Standard07FieldViewInput)
			standard07FieldViewArray := standard07FieldViewInput["fields"].(map[string]any)[standard07FieldViewKey].([]any)
			standard07FieldViewArray[0] = float64(129)
			standard07FieldViewOutput := standard07FieldViewEvaluate(standard07FieldViewT, standard07FieldViewInput)
			standard07FieldViewControl := standard07FieldViewDecode(standard07FieldViewT, standard07FieldViewCase.Standard07FieldViewInput)["results"].(map[string]any)
			for standard07FieldViewSibling, standard07FieldViewValue := range standard07FieldViewControl {
				if standard07FieldViewSibling != standard07FieldViewKey && !reflect.DeepEqual(standard07FieldViewOutput["results"].(map[string]any)[standard07FieldViewSibling], standard07FieldViewValue) {
					standard07FieldViewT.Fatal("field coupling")
				}
			}
		}
	}
}

func TestAb15Aq15Reject(standard07FieldViewT *testing.T) {
	for _, standard07FieldViewRaw := range []string{``, `null`, `[]`, `{} {} `, `{"context":{},"context":{},"fields":{}}`, `{"context":{"x":{"a":1,"a":2}},"fields":{}}`, `{"context":{},"fields":{},"x":[{"a":1,"a":2}]}`, `{"context":{},"fields":{},"x":1e999}`, `{"context":{},"fields":{}} true`, `{"context":{},"fields":{},}`} {
		standard07FieldViewResult, standard07FieldViewErr := Standard07FieldViewEvaluate([]byte(standard07FieldViewRaw))
		if standard07FieldViewErr == nil || standard07FieldViewResult != nil {
			standard07FieldViewT.Fatalf("accepted %s", standard07FieldViewRaw)
		}
	}
	for _, standard07FieldViewValue := range []any{nil, true, "1", 1.5, -1, 256, []any{}, map[string]any{}} {
		for _, standard07FieldViewKey := range []string{"primary", "secondary", "operation", "field"} {
			standard07FieldViewInput := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCases(standard07FieldViewT)[0].Standard07FieldViewInput)
			if standard07FieldViewKey == "field" {
				standard07FieldViewInput["fields"].(map[string]any)["outdoor_temperature"] = []any{standard07FieldViewValue, 0}
			} else {
				standard07FieldViewInput["context"].(map[string]any)[standard07FieldViewKey] = standard07FieldViewValue
			}
			standard07FieldViewData, _ := json.Marshal(standard07FieldViewInput)
			standard07FieldViewResult, standard07FieldViewErr := Standard07FieldViewEvaluate(standard07FieldViewData)
			if standard07FieldViewErr == nil || standard07FieldViewResult != nil {
				standard07FieldViewT.Fatalf("accepted %s", standard07FieldViewData)
			}
		}
	}
	for _, standard07FieldViewKey := range []string{"context", "fields", "arguments"} {
		for _, standard07FieldViewValue := range []any{nil, false, float64(0), "", []any{}} {
			standard07FieldViewInput := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCases(standard07FieldViewT)[0].Standard07FieldViewInput)
			standard07FieldViewInput[standard07FieldViewKey] = standard07FieldViewValue
			standard07FieldViewData, _ := json.Marshal(standard07FieldViewInput)
			if standard07FieldViewResult, standard07FieldViewErr := Standard07FieldViewEvaluate(standard07FieldViewData); standard07FieldViewErr == nil || standard07FieldViewResult != nil {
				standard07FieldViewT.Fatal("wrong kind accepted")
			}
		}
	}
	standard07FieldViewInput := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCases(standard07FieldViewT)[0].Standard07FieldViewInput)
	standard07FieldViewInput["arguments"] = map[string]any{"unexpected": 0}
	standard07FieldViewData, _ := json.Marshal(standard07FieldViewInput)
	if standard07FieldViewResult, standard07FieldViewErr := Standard07FieldViewEvaluate(standard07FieldViewData); standard07FieldViewErr == nil || standard07FieldViewResult != nil {
		standard07FieldViewT.Fatal("argument accepted")
	}
}

func TestAb15Aq15Storage(standard07FieldViewT *testing.T) {
	standard07FieldViewInput := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCases(standard07FieldViewT)[0].Standard07FieldViewInput)
	standard07FieldViewInput["arguments"] = map[string]any{}
	standard07FieldViewInput["extra"] = []any{nil, false, float64(0), "", []any{}}
	standard07FieldViewInput["fields"].(map[string]any)["opaque"] = []any{}
	standard07FieldViewOutput := standard07FieldViewEvaluate(standard07FieldViewT, standard07FieldViewInput)
	for _, standard07FieldViewKey := range []string{"arguments", "extra", "fields"} {
		if !reflect.DeepEqual(standard07FieldViewOutput[standard07FieldViewKey], standard07FieldViewInput[standard07FieldViewKey]) {
			standard07FieldViewT.Fatal("retention", standard07FieldViewKey)
		}
	}
	standard07FieldViewData, _ := json.Marshal(standard07FieldViewInput)
	standard07FieldViewFirst, standard07FieldViewErr := Standard07FieldViewEvaluate(standard07FieldViewData)
	if standard07FieldViewErr != nil {
		standard07FieldViewT.Fatal(standard07FieldViewErr)
	}
	standard07FieldViewSaved := append([]byte(nil), standard07FieldViewFirst...)
	for standard07FieldViewIndex := range standard07FieldViewData {
		standard07FieldViewData[standard07FieldViewIndex] = ' '
	}
	if !reflect.DeepEqual([]byte(standard07FieldViewFirst), standard07FieldViewSaved) {
		standard07FieldViewT.Fatal("input alias")
	}
	for standard07FieldViewIndex := range standard07FieldViewFirst {
		standard07FieldViewFirst[standard07FieldViewIndex] = ' '
	}
	if !reflect.DeepEqual(standard07FieldViewEvaluate(standard07FieldViewT, standard07FieldViewInput), standard07FieldViewOutput) {
		standard07FieldViewT.Fatal("output alias")
	}
}

func TestAb15Aq15ExactNumbersAndMissingSelectors(standard07FieldViewT *testing.T) {
	for _, standard07FieldViewNumber := range []string{"1e-999", "255.000000000000000001", "0.000000000000000001"} {
		for _, standard07FieldViewRaw := range []string{`{"context":{"primary":` + standard07FieldViewNumber + `},"fields":{}}`, `{"context":{},"fields":{"x":[` + standard07FieldViewNumber + `]}}`} {
			if standard07FieldViewResult, standard07FieldViewErr := Standard07FieldViewEvaluate([]byte(standard07FieldViewRaw)); standard07FieldViewErr == nil || standard07FieldViewResult != nil {
				standard07FieldViewT.Fatal("fraction accepted", standard07FieldViewRaw)
			}
		}
	}
	for _, standard07FieldViewKey := range []string{"primary", "secondary", "direction", "view"} {
		standard07FieldViewInput := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCases(standard07FieldViewT)[0].Standard07FieldViewInput)
		delete(standard07FieldViewInput["context"].(map[string]any), standard07FieldViewKey)
		if standard07FieldViewOutput := standard07FieldViewEvaluate(standard07FieldViewT, standard07FieldViewInput); standard07FieldViewOutput["state"] != "Unknown" {
			standard07FieldViewT.Fatal(standard07FieldViewOutput)
		}
	}
	for _, standard07FieldViewSecondary := range []int{1, 4, 254, 255} {
		standard07FieldViewInput := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCases(standard07FieldViewT)[0].Standard07FieldViewInput)
		standard07FieldViewInput["context"].(map[string]any)["secondary"] = standard07FieldViewSecondary
		if standard07FieldViewOutput := standard07FieldViewEvaluate(standard07FieldViewT, standard07FieldViewInput); standard07FieldViewOutput["state"] != "Unknown" {
			standard07FieldViewT.Fatal(standard07FieldViewOutput)
		}
	}
	for _, standard07FieldViewKey := range []string{"direction", "view"} {
		for _, standard07FieldViewValue := range []any{nil, 0, false, []any{}, map[string]any{}} {
			standard07FieldViewInput := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCases(standard07FieldViewT)[0].Standard07FieldViewInput)
			standard07FieldViewInput["context"].(map[string]any)[standard07FieldViewKey] = standard07FieldViewValue
			standard07FieldViewData, _ := json.Marshal(standard07FieldViewInput)
			if standard07FieldViewResult, standard07FieldViewErr := Standard07FieldViewEvaluate(standard07FieldViewData); standard07FieldViewErr == nil || standard07FieldViewResult != nil {
				standard07FieldViewT.Fatal("identifier kind accepted")
			}
		}
	}
}

func TestAb15Aq15ReservedKeysAndUnsupportedArguments(standard07FieldViewT *testing.T) {
	for _, standard07FieldViewKey := range []string{"results", "state", "physical", "qualified"} {
		standard07FieldViewInput := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCases(standard07FieldViewT)[0].Standard07FieldViewInput)
		standard07FieldViewInput[standard07FieldViewKey] = "x"
		standard07FieldViewData, _ := json.Marshal(standard07FieldViewInput)
		if standard07FieldViewResult, standard07FieldViewErr := Standard07FieldViewEvaluate(standard07FieldViewData); standard07FieldViewErr == nil || standard07FieldViewResult != nil {
			standard07FieldViewT.Fatal("reserved key accepted", standard07FieldViewKey)
		}
	}
	standard07FieldViewInput := standard07FieldViewInput(standard07FieldViewT, standard07FieldViewCases(standard07FieldViewT)[0].Standard07FieldViewInput)
	standard07FieldViewInput["context"].(map[string]any)["secondary"] = float64(1)
	standard07FieldViewInput["arguments"] = map[string]any{"unexpected": float64(0)}
	standard07FieldViewOutput := standard07FieldViewEvaluate(standard07FieldViewT, standard07FieldViewInput)
	if standard07FieldViewOutput["state"] != "Unknown" || !reflect.DeepEqual(standard07FieldViewOutput["arguments"], standard07FieldViewInput["arguments"]) {
		standard07FieldViewT.Fatal("unsupported context must stay raw-only Unknown")
	}
}

func TestAb15Aq15DepthLimit(standard07FieldViewT *testing.T) {
	standard07FieldViewDeep := `{"context":{},"fields":{},"x":` + strings.Repeat("[", 100) + strings.Repeat("]", 100) + `}`
	if standard07FieldViewResult, standard07FieldViewErr := Standard07FieldViewEvaluate([]byte(standard07FieldViewDeep)); standard07FieldViewErr == nil || standard07FieldViewResult != nil {
		standard07FieldViewT.Fatal("deep nesting accepted")
	}
}

func TestAb15Aq15RejectArgumentsOnUnsupportedContext(standard07FieldViewT *testing.T) {
	for _, standard07FieldViewRaw := range []string{`{"context":{"primary":8},"fields":{},"arguments":5}`, `{"context":{},"fields":{},"arguments":[]}`} {
		if standard07FieldViewResult, standard07FieldViewErr := Standard07FieldViewEvaluate([]byte(standard07FieldViewRaw)); standard07FieldViewErr == nil || standard07FieldViewResult != nil {
			standard07FieldViewT.Fatal("arguments accepted", standard07FieldViewRaw)
		}
	}
}
