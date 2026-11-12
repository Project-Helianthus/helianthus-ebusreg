package ebusstandard

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func standardFFExtentViewRun(standardFFExtentViewT *testing.T, standardFFExtentViewInput any) map[string]any {
	standardFFExtentViewT.Helper()
	standardFFExtentViewRaw, standardFFExtentViewErr := json.Marshal(standardFFExtentViewInput)
	if standardFFExtentViewErr != nil {
		standardFFExtentViewT.Fatal(standardFFExtentViewErr)
	}
	standardFFExtentViewOutput, standardFFExtentViewErr := StandardFFExtentViewEvaluate(standardFFExtentViewRaw)
	if standardFFExtentViewErr != nil {
		standardFFExtentViewT.Fatal(standardFFExtentViewErr)
	}
	var standardFFExtentViewResult map[string]any
	if standardFFExtentViewErr := json.Unmarshal(standardFFExtentViewOutput, &standardFFExtentViewResult); standardFFExtentViewErr != nil {
		standardFFExtentViewT.Fatal(standardFFExtentViewErr)
	}
	if standardFFExtentViewResult["physical"] != "Unknown" || standardFFExtentViewResult["qualified"] != false {
		standardFFExtentViewT.Fatal("arithmetic qualified physical support", standardFFExtentViewResult)
	}
	return standardFFExtentViewResult
}

func standardFFExtentViewBase() map[string]any {
	return map[string]any{"context": map[string]any{"primary": 255, "secondary": 4, "direction": "positive", "view": "body-extent", "family": "synthetic-family", "variant": ""}, "fields": map[string]any{"body": []int{0, 1}, "opaque": []int{}}, "arguments": map[string]any{"declared_length": 2}}
}

func TestAb17Aq17Fixture(standardFFExtentViewT *testing.T) {
	standardFFExtentViewRaw, standardFFExtentViewErr := os.ReadFile("testdata/ab17aq17_standardff_extent_view.json")
	if standardFFExtentViewErr != nil {
		standardFFExtentViewT.Fatal(standardFFExtentViewErr)
	}
	var standardFFExtentViewCases []struct {
		StandardFFExtentViewID       string         `json:"id"`
		StandardFFExtentViewInput    map[string]any `json:"input"`
		StandardFFExtentViewExpected map[string]any `json:"expected"`
	}
	if standardFFExtentViewErr := json.Unmarshal(standardFFExtentViewRaw, &standardFFExtentViewCases); standardFFExtentViewErr != nil {
		standardFFExtentViewT.Fatal(standardFFExtentViewErr)
	}
	if len(standardFFExtentViewCases) != 11 {
		standardFFExtentViewT.Fatal("fixture count")
	}
	for _, standardFFExtentViewCase := range standardFFExtentViewCases {
		standardFFExtentViewT.Run(standardFFExtentViewCase.StandardFFExtentViewID, func(standardFFExtentViewT *testing.T) {
			standardFFExtentViewOut := standardFFExtentViewRun(standardFFExtentViewT, standardFFExtentViewCase.StandardFFExtentViewInput)
			for _, standardFFExtentViewKey := range []string{"context", "fields", "arguments", "extra"} {
				if !reflect.DeepEqual(standardFFExtentViewOut[standardFFExtentViewKey], standardFFExtentViewCase.StandardFFExtentViewInput[standardFFExtentViewKey]) {
					standardFFExtentViewT.Fatal("raw input changed", standardFFExtentViewKey)
				}
			}
			standardFFExtentViewResults := standardFFExtentViewOut["results"].(map[string]any)
			if len(standardFFExtentViewResults) != 2 {
				standardFFExtentViewT.Fatal("decoded unprovided fields")
			}
			standardFFExtentViewBody := standardFFExtentViewResults["body"].(map[string]any)
			if standardFFExtentViewBody["state"] != "Unknown" || standardFFExtentViewBody["calculation"] != "Proven" {
				standardFFExtentViewT.Fatal("semantic grade", standardFFExtentViewBody)
			}
			if !reflect.DeepEqual(standardFFExtentViewBody["raw"], standardFFExtentViewCase.StandardFFExtentViewInput["fields"].(map[string]any)["body"]) {
				standardFFExtentViewT.Fatal("body raw")
			}
			for _, standardFFExtentViewKey := range []string{"flags", "declared_pair_count", "available_pair_count"} {
				if !reflect.DeepEqual(standardFFExtentViewBody[standardFFExtentViewKey], standardFFExtentViewCase.StandardFFExtentViewExpected[standardFFExtentViewKey]) {
					standardFFExtentViewT.Fatalf("%s: got %v want %v", standardFFExtentViewKey, standardFFExtentViewBody[standardFFExtentViewKey], standardFFExtentViewCase.StandardFFExtentViewExpected[standardFFExtentViewKey])
				}
			}
			standardFFExtentViewFlags := standardFFExtentViewBody["flags"].(map[string]any)
			standardFFExtentViewComplete := standardFFExtentViewFlags["bounded_minimum"] == true || standardFFExtentViewFlags["whole_declared_span"] == true
			if (standardFFExtentViewOut["state"] == "complete") != standardFFExtentViewComplete {
				standardFFExtentViewT.Fatal("extent completeness", standardFFExtentViewOut)
			}
		})
	}
}

func TestAb17Aq17Identity(standardFFExtentViewT *testing.T) {
	for standardFFExtentViewKey, standardFFExtentViewValues := range map[string][]any{"primary": {254, nil}, "secondary": {5, nil}, "operation": {0, 1}, "direction": {"negative", ""}, "view": {"other", ""}} {
		for _, standardFFExtentViewValue := range standardFFExtentViewValues {
			standardFFExtentViewInput := standardFFExtentViewBase()
			standardFFExtentViewContext := standardFFExtentViewInput["context"].(map[string]any)
			if standardFFExtentViewValue == nil {
				delete(standardFFExtentViewContext, standardFFExtentViewKey)
			} else {
				standardFFExtentViewContext[standardFFExtentViewKey] = standardFFExtentViewValue
			}
			standardFFExtentViewOut := standardFFExtentViewRun(standardFFExtentViewT, standardFFExtentViewInput)
			standardFFExtentViewBody := standardFFExtentViewOut["results"].(map[string]any)["body"].(map[string]any)
			if standardFFExtentViewOut["state"] != "Unknown" || len(standardFFExtentViewBody) != 2 || standardFFExtentViewBody["state"] != "Unknown" {
				standardFFExtentViewT.Fatal("unsupported identity calculated", standardFFExtentViewKey, standardFFExtentViewOut)
			}
		}
	}
	for _, standardFFExtentViewKey := range []string{"family", "variant"} {
		for _, standardFFExtentViewValue := range []any{nil, false, 0, "another-family", "", []any{}} {
			standardFFExtentViewInput := standardFFExtentViewBase()
			standardFFExtentViewInput["context"].(map[string]any)[standardFFExtentViewKey] = standardFFExtentViewValue
			standardFFExtentViewOut := standardFFExtentViewRun(standardFFExtentViewT, standardFFExtentViewInput)
			if standardFFExtentViewOut["state"] != "complete" || !reflect.DeepEqual(standardFFExtentViewOut["context"].(map[string]any)[standardFFExtentViewKey], func() any {
				standardFFExtentViewRaw, _ := json.Marshal(standardFFExtentViewValue)
				var standardFFExtentViewNormalized any
				_ = json.Unmarshal(standardFFExtentViewRaw, &standardFFExtentViewNormalized)
				return standardFFExtentViewNormalized
			}()) {
				standardFFExtentViewT.Fatal("opaque identity", standardFFExtentViewKey)
			}
		}
	}
}

func TestAb17Aq17Invalid(standardFFExtentViewT *testing.T) {
	for _, standardFFExtentViewRaw := range []string{``, `{`, `[]`, `null`, `{} {}`, `{"x":{"a":1,"a":2}}`, `{"x":[{"a":1,"\u0061":2}]}`, `{"x":NaN}`} {
		if standardFFExtentViewOut, standardFFExtentViewErr := StandardFFExtentViewEvaluate(json.RawMessage(standardFFExtentViewRaw)); standardFFExtentViewErr == nil || standardFFExtentViewOut != nil {
			standardFFExtentViewT.Fatal("accepted malformed", standardFFExtentViewRaw)
		}
	}
	for _, standardFFExtentViewKey := range []string{"primary", "secondary", "operation", "declared_length", "body", "opaque", "direction", "view", "arguments", "fields", "context"} {
		for _, standardFFExtentViewBad := range []any{nil, true, "2", -1, 256, 1.5} {
			standardFFExtentViewInput := standardFFExtentViewBase()
			switch standardFFExtentViewKey {
			case "primary", "secondary", "operation":
				standardFFExtentViewInput["context"].(map[string]any)[standardFFExtentViewKey] = standardFFExtentViewBad
			case "direction", "view":
				if _, standardFFExtentViewOK := standardFFExtentViewBad.(string); standardFFExtentViewOK {
					continue
				}
				standardFFExtentViewInput["context"].(map[string]any)[standardFFExtentViewKey] = standardFFExtentViewBad
			case "declared_length":
				standardFFExtentViewInput["arguments"].(map[string]any)[standardFFExtentViewKey] = standardFFExtentViewBad
			case "body", "opaque":
				standardFFExtentViewInput["fields"].(map[string]any)[standardFFExtentViewKey] = []any{standardFFExtentViewBad}
			default:
				standardFFExtentViewInput[standardFFExtentViewKey] = standardFFExtentViewBad
			}
			standardFFExtentViewRaw, _ := json.Marshal(standardFFExtentViewInput)
			if standardFFExtentViewOut, standardFFExtentViewErr := StandardFFExtentViewEvaluate(standardFFExtentViewRaw); standardFFExtentViewErr == nil || standardFFExtentViewOut != nil {
				standardFFExtentViewT.Fatal("accepted wrong kind", standardFFExtentViewKey, standardFFExtentViewBad)
			}
		}
	}
	for _, standardFFExtentViewRaw := range []string{`null`, `true`, `{}`, `""`} {
		standardFFExtentViewRaw = []string{`{"context":{"primary":255,"secondary":4,"direction":"positive","view":"body-extent"},"fields":{"body":` + standardFFExtentViewRaw + `},"arguments":{"declared_length":2}}`}[0]
		if standardFFExtentViewOut, standardFFExtentViewErr := StandardFFExtentViewEvaluate([]byte(standardFFExtentViewRaw)); standardFFExtentViewErr == nil || standardFFExtentViewOut != nil {
			standardFFExtentViewT.Fatal("nonarray body")
		}
	}
	for _, standardFFExtentViewKey := range []string{"arguments", "declared_length"} {
		standardFFExtentViewInput := standardFFExtentViewBase()
		if standardFFExtentViewKey == "arguments" {
			delete(standardFFExtentViewInput, standardFFExtentViewKey)
		} else {
			delete(standardFFExtentViewInput["arguments"].(map[string]any), standardFFExtentViewKey)
		}
		standardFFExtentViewRaw, _ := json.Marshal(standardFFExtentViewInput)
		if _, standardFFExtentViewErr := StandardFFExtentViewEvaluate(standardFFExtentViewRaw); standardFFExtentViewErr == nil {
			standardFFExtentViewT.Fatal("missing argument")
		}
	}
}

func TestAb17Aq17PresenceAndMemory(standardFFExtentViewT *testing.T) {
	for _, standardFFExtentViewBody := range [][]int{nil, {}, {0}, {0, 1}, {0, 1, 2}} {
		standardFFExtentViewInput := standardFFExtentViewBase()
		standardFFExtentViewFields := standardFFExtentViewInput["fields"].(map[string]any)
		if standardFFExtentViewBody == nil {
			delete(standardFFExtentViewFields, "body")
		} else {
			standardFFExtentViewFields["body"] = standardFFExtentViewBody
		}
		standardFFExtentViewOut := standardFFExtentViewRun(standardFFExtentViewT, standardFFExtentViewInput)
		standardFFExtentViewResults := standardFFExtentViewOut["results"].(map[string]any)
		if !reflect.DeepEqual(standardFFExtentViewResults["opaque"], map[string]any{"raw": []any{}, "state": "Unknown"}) {
			standardFFExtentViewT.Fatal("independent sibling lost")
		}
		standardFFExtentViewResult := standardFFExtentViewResults["body"].(map[string]any)
		if standardFFExtentViewBody == nil {
			if _, standardFFExtentViewOK := standardFFExtentViewResult["raw"]; standardFFExtentViewOK {
				standardFFExtentViewT.Fatal("invented raw")
			}
			if _, standardFFExtentViewOK := standardFFExtentViewResult["available_length"]; standardFFExtentViewOK {
				standardFFExtentViewT.Fatal("invented length")
			}
		} else if standardFFExtentViewResult["available_length"] != float64(len(standardFFExtentViewBody)) {
			standardFFExtentViewT.Fatal("available length")
		}
	}
	standardFFExtentViewRaw, _ := json.Marshal(standardFFExtentViewBase())
	standardFFExtentViewOriginal := append([]byte(nil), standardFFExtentViewRaw...)
	standardFFExtentViewFirst, standardFFExtentViewErr := StandardFFExtentViewEvaluate(standardFFExtentViewRaw)
	if standardFFExtentViewErr != nil {
		standardFFExtentViewT.Fatal(standardFFExtentViewErr)
	}
	standardFFExtentViewSaved := string(standardFFExtentViewFirst)
	for standardFFExtentViewIndex := range standardFFExtentViewRaw {
		standardFFExtentViewRaw[standardFFExtentViewIndex] = ' '
	}
	if string(standardFFExtentViewFirst) != standardFFExtentViewSaved {
		standardFFExtentViewT.Fatal("input aliases output")
	}
	for standardFFExtentViewIndex := range standardFFExtentViewFirst {
		standardFFExtentViewFirst[standardFFExtentViewIndex] = ' '
	}
	standardFFExtentViewNext, standardFFExtentViewErr := StandardFFExtentViewEvaluate(standardFFExtentViewOriginal)
	if standardFFExtentViewErr != nil || string(standardFFExtentViewNext) != standardFFExtentViewSaved {
		standardFFExtentViewT.Fatal("cross call mutation")
	}
	if strings.Contains(standardFFExtentViewSaved, "address") || strings.Contains(standardFFExtentViewSaved, "status") {
		standardFFExtentViewT.Fatal("invented decoder")
	}
}

func TestAb17Aq17AdditionalBoundaries(standardFFExtentViewT *testing.T) {
	for _, standardFFExtentViewKey := range []string{"family", "variant"} {
		standardFFExtentViewInput := standardFFExtentViewBase()
		delete(standardFFExtentViewInput["context"].(map[string]any), standardFFExtentViewKey)
		if standardFFExtentViewRun(standardFFExtentViewT, standardFFExtentViewInput)["state"] != "complete" {
			standardFFExtentViewT.Fatal("opaque absence changes arithmetic")
		}
	}
	for _, standardFFExtentViewLength := range []int{0, 1, 11, 255} {
		standardFFExtentViewInput := standardFFExtentViewBase()
		standardFFExtentViewInput["arguments"].(map[string]any)["declared_length"] = standardFFExtentViewLength
		if standardFFExtentViewRun(standardFFExtentViewT, standardFFExtentViewInput)["state"] != "partial" {
			standardFFExtentViewT.Fatal("invalid minimum admitted", standardFFExtentViewLength)
		}
	}
	standardFFExtentViewInput := standardFFExtentViewBase()
	standardFFExtentViewInput["fields"].(map[string]any)["body"] = []int{0, 255, 254}
	standardFFExtentViewOut := standardFFExtentViewRun(standardFFExtentViewT, standardFFExtentViewInput)
	if standardFFExtentViewOut["state"] != "complete" {
		standardFFExtentViewT.Fatal("extra byte repairs or breaks minimum")
	}
	delete(standardFFExtentViewInput["fields"].(map[string]any), "opaque")
	standardFFExtentViewOut = standardFFExtentViewRun(standardFFExtentViewT, standardFFExtentViewInput)
	if len(standardFFExtentViewOut["results"].(map[string]any)) != 1 {
		standardFFExtentViewT.Fatal("invented absent sibling")
	}
	standardFFExtentViewInput["arguments"].(map[string]any)["extra"] = 0
	standardFFExtentViewRaw, _ := json.Marshal(standardFFExtentViewInput)
	if _, standardFFExtentViewErr := StandardFFExtentViewEvaluate(standardFFExtentViewRaw); standardFFExtentViewErr == nil {
		standardFFExtentViewT.Fatal("extra typed argument admitted")
	}
	for _, standardFFExtentViewBody := range [][]int{nil, {}} {
		standardFFExtentViewInput = standardFFExtentViewBase()
		standardFFExtentViewInput["context"].(map[string]any)["secondary"] = 6
		if standardFFExtentViewBody == nil {
			delete(standardFFExtentViewInput["fields"].(map[string]any), "body")
		} else {
			standardFFExtentViewInput["fields"].(map[string]any)["body"] = standardFFExtentViewBody
		}
		standardFFExtentViewResult := standardFFExtentViewRun(standardFFExtentViewT, standardFFExtentViewInput)["results"].(map[string]any)["body"].(map[string]any)
		if _, standardFFExtentViewOK := standardFFExtentViewResult["available_pair_count"]; standardFFExtentViewOK {
			standardFFExtentViewT.Fatal("empty or absent data invented pair count")
		}
	}
}

// TestAb17Aq17EvaluateExportedContract drives the exported entry point with
// literal JSON and compares the whole decoded result (or the error) with
// independently written expectations, so a change to the extent arithmetic,
// the identity gate, raw preservation or input rejection fails here.
func TestAb17Aq17EvaluateExportedContract(standardFFExtentViewT *testing.T) {
	standardFFExtentViewCases := []struct {
		standardFFExtentViewName  string
		standardFFExtentViewInput string
		standardFFExtentViewWant  string
		standardFFExtentViewErr   string
	}{
		{
			standardFFExtentViewName:  "ff06-declared5-available3",
			standardFFExtentViewInput: `{"context":{"primary":255,"secondary":6,"direction":"positive","view":"body-extent"},"fields":{"body":[1,2,3]},"arguments":{"declared_length":5}}`,
			standardFFExtentViewWant:  `{"context":{"primary":255,"secondary":6,"direction":"positive","view":"body-extent"},"fields":{"body":[1,2,3]},"arguments":{"declared_length":5},"state":"partial","physical":"Unknown","qualified":false,"results":{"body":{"state":"Unknown","calculation":"Proven","declared_length":5,"raw":[1,2,3],"available_length":3,"declared_pair_count":2,"available_pair_count":1,"flags":{"declared_incomplete_tail":false,"whole_declared_span":false}}}}`,
		},
		{
			standardFFExtentViewName:  "ff06-declared4-available6",
			standardFFExtentViewInput: `{"context":{"primary":255,"secondary":6,"direction":"positive","view":"body-extent"},"fields":{"body":[1,2,3,4,5,6]},"arguments":{"declared_length":4}}`,
			standardFFExtentViewWant:  `{"context":{"primary":255,"secondary":6,"direction":"positive","view":"body-extent"},"fields":{"body":[1,2,3,4,5,6]},"arguments":{"declared_length":4},"state":"partial","physical":"Unknown","qualified":false,"results":{"body":{"state":"Unknown","calculation":"Proven","declared_length":4,"raw":[1,2,3,4,5,6],"available_length":6,"declared_pair_count":1,"available_pair_count":1,"flags":{"declared_incomplete_tail":true,"whole_declared_span":false}}}}`,
		},
		{
			standardFFExtentViewName:  "ff04-declared10-available10",
			standardFFExtentViewInput: `{"context":{"primary":255,"secondary":4,"direction":"positive","view":"body-extent"},"fields":{"body":[0,1,2,3,4,5,6,7,8,9]},"arguments":{"declared_length":10}}`,
			standardFFExtentViewWant:  `{"context":{"primary":255,"secondary":4,"direction":"positive","view":"body-extent"},"fields":{"body":[0,1,2,3,4,5,6,7,8,9]},"arguments":{"declared_length":10},"state":"complete","physical":"Unknown","qualified":false,"results":{"body":{"state":"Unknown","calculation":"Proven","declared_length":10,"raw":[0,1,2,3,4,5,6,7,8,9],"available_length":10,"flags":{"declared_minimum_satisfied":true,"initialized_minimum_satisfied":true,"bounded_minimum":true}}}}`,
		},
		{
			standardFFExtentViewName:  "ff04-declared1-available1",
			standardFFExtentViewInput: `{"context":{"primary":255,"secondary":4,"direction":"positive","view":"body-extent"},"fields":{"body":[7]},"arguments":{"declared_length":1}}`,
			standardFFExtentViewWant:  `{"context":{"primary":255,"secondary":4,"direction":"positive","view":"body-extent"},"fields":{"body":[7]},"arguments":{"declared_length":1},"state":"partial","physical":"Unknown","qualified":false,"results":{"body":{"state":"Unknown","calculation":"Proven","declared_length":1,"raw":[7],"available_length":1,"flags":{"declared_minimum_satisfied":false,"initialized_minimum_satisfied":false,"bounded_minimum":false}}}}`,
		},
		{
			standardFFExtentViewName:  "zero-operation-stays-raw-only",
			standardFFExtentViewInput: `{"context":{"primary":255,"secondary":4,"operation":0,"direction":"positive","view":"body-extent"},"fields":{"body":[0,1]},"arguments":{"declared_length":2}}`,
			standardFFExtentViewWant:  `{"context":{"primary":255,"secondary":4,"operation":0,"direction":"positive","view":"body-extent"},"fields":{"body":[0,1]},"arguments":{"declared_length":2},"state":"Unknown","physical":"Unknown","qualified":false,"results":{"body":{"raw":[0,1],"state":"Unknown"}}}`,
		},
		{
			standardFFExtentViewName:  "trailing-json",
			standardFFExtentViewInput: `{"context":{"primary":255,"secondary":6,"direction":"positive","view":"body-extent"},"fields":{"body":[1,2,3]},"arguments":{"declared_length":5}} {}`,
			standardFFExtentViewErr:   "trailing JSON",
		},
		{
			standardFFExtentViewName:  "declared-length-outside-byte",
			standardFFExtentViewInput: `{"context":{"primary":255,"secondary":6,"direction":"positive","view":"body-extent"},"fields":{"body":[1,2,3]},"arguments":{"declared_length":256}}`,
			standardFFExtentViewErr:   "integer outside byte range",
		},
	}
	for _, standardFFExtentViewCase := range standardFFExtentViewCases {
		standardFFExtentViewT.Run(standardFFExtentViewCase.standardFFExtentViewName, func(standardFFExtentViewT *testing.T) {
			standardFFExtentViewOut, standardFFExtentViewErr := StandardFFExtentViewEvaluate(json.RawMessage(standardFFExtentViewCase.standardFFExtentViewInput))
			if standardFFExtentViewCase.standardFFExtentViewErr != "" {
				if standardFFExtentViewErr == nil || !strings.Contains(standardFFExtentViewErr.Error(), standardFFExtentViewCase.standardFFExtentViewErr) || standardFFExtentViewOut != nil {
					standardFFExtentViewT.Fatalf("got output %s error %v, want nil output and error containing %q", standardFFExtentViewOut, standardFFExtentViewErr, standardFFExtentViewCase.standardFFExtentViewErr)
				}
				return
			}
			if standardFFExtentViewErr != nil {
				standardFFExtentViewT.Fatalf("unexpected error: %v", standardFFExtentViewErr)
			}
			var standardFFExtentViewGot, standardFFExtentViewWant any
			if standardFFExtentViewErr := json.Unmarshal(standardFFExtentViewOut, &standardFFExtentViewGot); standardFFExtentViewErr != nil {
				standardFFExtentViewT.Fatal(standardFFExtentViewErr)
			}
			if standardFFExtentViewErr := json.Unmarshal([]byte(standardFFExtentViewCase.standardFFExtentViewWant), &standardFFExtentViewWant); standardFFExtentViewErr != nil {
				standardFFExtentViewT.Fatal(standardFFExtentViewErr)
			}
			if !reflect.DeepEqual(standardFFExtentViewGot, standardFFExtentViewWant) {
				standardFFExtentViewT.Fatalf("got %s\nwant %s", standardFFExtentViewOut, standardFFExtentViewCase.standardFFExtentViewWant)
			}
		})
	}
}

func TestAb17Aq17OutputKeyCollision(standardFFExtentViewT *testing.T) {
	for _, standardFFExtentViewKey := range []string{"state", "results", "physical", "qualified"} {
		standardFFExtentViewInput := standardFFExtentViewBase()
		standardFFExtentViewInput[standardFFExtentViewKey] = 7
		standardFFExtentViewRaw, standardFFExtentViewErr := json.Marshal(standardFFExtentViewInput)
		if standardFFExtentViewErr != nil {
			standardFFExtentViewT.Fatal(standardFFExtentViewErr)
		}
		if standardFFExtentViewOut, standardFFExtentViewErr := StandardFFExtentViewEvaluate(standardFFExtentViewRaw); standardFFExtentViewErr == nil || standardFFExtentViewOut != nil {
			standardFFExtentViewT.Fatal("accepted colliding key", standardFFExtentViewKey)
		}
	}
}
