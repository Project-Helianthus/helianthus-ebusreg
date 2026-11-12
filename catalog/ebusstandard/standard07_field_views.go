package ebusstandard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"strconv"
)

func standard07FieldViewRead(standard07FieldViewDecoder *json.Decoder, standard07FieldViewDepth int) (any, error) {
	if standard07FieldViewDepth > 64 {
		return nil, fmt.Errorf("nesting too deep")
	}
	standard07FieldViewToken, standard07FieldViewErr := standard07FieldViewDecoder.Token()
	if standard07FieldViewErr != nil {
		return nil, standard07FieldViewErr
	}
	if standard07FieldViewNumber, standard07FieldViewOK := standard07FieldViewToken.(json.Number); standard07FieldViewOK {
		standard07FieldViewValue, standard07FieldViewErr := strconv.ParseFloat(string(standard07FieldViewNumber), 64)
		if standard07FieldViewErr != nil || math.IsInf(standard07FieldViewValue, 0) || math.IsNaN(standard07FieldViewValue) {
			return nil, fmt.Errorf("invalid number")
		}
	}
	if standard07FieldViewDelimiter, standard07FieldViewOK := standard07FieldViewToken.(json.Delim); standard07FieldViewOK {
		if standard07FieldViewDelimiter == '{' {
			standard07FieldViewObject := map[string]any{}
			for standard07FieldViewDecoder.More() {
				standard07FieldViewKeyToken, standard07FieldViewErr := standard07FieldViewDecoder.Token()
				if standard07FieldViewErr != nil {
					return nil, standard07FieldViewErr
				}
				standard07FieldViewKey, standard07FieldViewOK := standard07FieldViewKeyToken.(string)
				if !standard07FieldViewOK {
					return nil, fmt.Errorf("invalid key")
				}
				if _, standard07FieldViewFound := standard07FieldViewObject[standard07FieldViewKey]; standard07FieldViewFound {
					return nil, fmt.Errorf("duplicate key %q", standard07FieldViewKey)
				}
				standard07FieldViewObject[standard07FieldViewKey], standard07FieldViewErr = standard07FieldViewRead(standard07FieldViewDecoder, standard07FieldViewDepth+1)
				if standard07FieldViewErr != nil {
					return nil, standard07FieldViewErr
				}
			}
			_, standard07FieldViewErr := standard07FieldViewDecoder.Token()
			return standard07FieldViewObject, standard07FieldViewErr
		}
		if standard07FieldViewDelimiter == '[' {
			standard07FieldViewArray := []any{}
			for standard07FieldViewDecoder.More() {
				standard07FieldViewValue, standard07FieldViewErr := standard07FieldViewRead(standard07FieldViewDecoder, standard07FieldViewDepth+1)
				if standard07FieldViewErr != nil {
					return nil, standard07FieldViewErr
				}
				standard07FieldViewArray = append(standard07FieldViewArray, standard07FieldViewValue)
			}
			_, standard07FieldViewErr := standard07FieldViewDecoder.Token()
			return standard07FieldViewArray, standard07FieldViewErr
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	return standard07FieldViewToken, nil
}

func standard07FieldViewByte(standard07FieldViewValue any) (int, bool) {
	standard07FieldViewNumber, standard07FieldViewOK := standard07FieldViewValue.(json.Number)
	if !standard07FieldViewOK {
		return 0, false
	}
	standard07FieldViewRational, standard07FieldViewOK := new(big.Rat).SetString(string(standard07FieldViewNumber))
	if !standard07FieldViewOK || !standard07FieldViewRational.IsInt() || !standard07FieldViewRational.Num().IsInt64() {
		return 0, false
	}
	standard07FieldViewInteger := standard07FieldViewRational.Num().Int64()
	return int(standard07FieldViewInteger), standard07FieldViewInteger >= 0 && standard07FieldViewInteger <= 255
}

func Standard07FieldViewEvaluate(standard07FieldViewInput json.RawMessage) (json.RawMessage, error) {
	standard07FieldViewDecoder := json.NewDecoder(bytes.NewReader(standard07FieldViewInput))
	standard07FieldViewDecoder.UseNumber()
	standard07FieldViewValue, standard07FieldViewErr := standard07FieldViewRead(standard07FieldViewDecoder, 0)
	if standard07FieldViewErr != nil {
		return nil, standard07FieldViewErr
	}
	if _, standard07FieldViewErr = standard07FieldViewDecoder.Token(); standard07FieldViewErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	standard07FieldViewOutput, standard07FieldViewOK := standard07FieldViewValue.(map[string]any)
	if !standard07FieldViewOK {
		return nil, fmt.Errorf("input must be object")
	}
	standard07FieldViewContext, standard07FieldViewOK := standard07FieldViewOutput["context"].(map[string]any)
	if !standard07FieldViewOK {
		return nil, fmt.Errorf("context must be object")
	}
	standard07FieldViewFields, standard07FieldViewOK := standard07FieldViewOutput["fields"].(map[string]any)
	if !standard07FieldViewOK {
		return nil, fmt.Errorf("fields must be object")
	}
	for _, standard07FieldViewKey := range []string{"primary", "secondary", "operation"} {
		if standard07FieldViewValue, standard07FieldViewFound := standard07FieldViewContext[standard07FieldViewKey]; standard07FieldViewFound {
			if _, standard07FieldViewOK := standard07FieldViewByte(standard07FieldViewValue); !standard07FieldViewOK {
				return nil, fmt.Errorf("invalid %s", standard07FieldViewKey)
			}
		}
	}
	for _, standard07FieldViewKey := range []string{"direction", "view"} {
		if standard07FieldViewValue, standard07FieldViewFound := standard07FieldViewContext[standard07FieldViewKey]; standard07FieldViewFound {
			if _, standard07FieldViewOK := standard07FieldViewValue.(string); !standard07FieldViewOK {
				return nil, fmt.Errorf("invalid %s", standard07FieldViewKey)
			}
		}
	}
	for _, standard07FieldViewKey := range []string{"results", "state", "physical", "qualified"} {
		if _, standard07FieldViewFound := standard07FieldViewOutput[standard07FieldViewKey]; standard07FieldViewFound {
			return nil, fmt.Errorf("reserved key %q in input", standard07FieldViewKey)
		}
	}
	standard07FieldViewHasArguments := false
	if standard07FieldViewValue, standard07FieldViewFound := standard07FieldViewOutput["arguments"]; standard07FieldViewFound {
		standard07FieldViewArguments, standard07FieldViewOK := standard07FieldViewValue.(map[string]any)
		if !standard07FieldViewOK {
			return nil, fmt.Errorf("arguments must be object")
		}
		standard07FieldViewHasArguments = len(standard07FieldViewArguments) != 0
	}
	standard07FieldViewResults := map[string]any{}
	for standard07FieldViewKey, standard07FieldViewValue := range standard07FieldViewFields {
		standard07FieldViewArray, standard07FieldViewOK := standard07FieldViewValue.([]any)
		if !standard07FieldViewOK {
			return nil, fmt.Errorf("field %s must be byte array", standard07FieldViewKey)
		}
		for _, standard07FieldViewValue := range standard07FieldViewArray {
			if _, standard07FieldViewOK := standard07FieldViewByte(standard07FieldViewValue); !standard07FieldViewOK {
				return nil, fmt.Errorf("invalid field byte")
			}
		}
		standard07FieldViewResults[standard07FieldViewKey] = map[string]any{"raw": standard07FieldViewValue, "state": "Unknown"}
	}
	standard07FieldViewOutput["results"] = standard07FieldViewResults
	standard07FieldViewOutput["state"] = "Unknown"
	standard07FieldViewOutput["physical"] = "Unknown"
	standard07FieldViewOutput["qualified"] = false
	standard07FieldViewPrimary, _ := standard07FieldViewByte(standard07FieldViewContext["primary"])
	standard07FieldViewSecondary, standard07FieldViewHasSecondary := standard07FieldViewByte(standard07FieldViewContext["secondary"])
	_, standard07FieldViewHasOperation := standard07FieldViewContext["operation"]
	standard07FieldViewDirection, _ := standard07FieldViewContext["direction"].(string)
	if standard07FieldViewPrimary != 7 || !standard07FieldViewHasSecondary || (standard07FieldViewSecondary != 0 && standard07FieldViewSecondary != 2 && standard07FieldViewSecondary != 3 && standard07FieldViewSecondary != 5) || standard07FieldViewHasOperation || standard07FieldViewContext["view"] != "field-candidates" || (standard07FieldViewDirection != "request" && standard07FieldViewDirection != "positive") {
		return json.Marshal(standard07FieldViewOutput)
	}
	if standard07FieldViewHasArguments {
		return nil, fmt.Errorf("no arguments defined for this view")
	}
	standard07FieldViewWidths := map[string]int{}
	if standard07FieldViewDirection == "positive" {
		standard07FieldViewWidths = map[string]int{"version": 1, "revision": 1, "bitmap": 8}
	} else {
		switch standard07FieldViewSecondary {
		case 0:
			standard07FieldViewWidths = map[string]int{"outdoor_temperature": 2, "weekday": 1, "second": 1, "minute": 1, "hour": 1, "day": 1, "month": 1, "year": 1}
		case 2:
			standard07FieldViewWidths = map[string]int{"outdoor_temperature": 2, "duration": 1}
		case 3:
			standard07FieldViewWidths = map[string]int{"block": 1}
		case 5:
			standard07FieldViewWidths = map[string]int{"secondary_block": 1, "primary_block": 1}
		}
	}
	standard07FieldViewOutput["state"] = "complete"
	for standard07FieldViewKey, standard07FieldViewWidth := range standard07FieldViewWidths {
		standard07FieldViewRaw, standard07FieldViewPresent := standard07FieldViewFields[standard07FieldViewKey]
		standard07FieldViewResult := map[string]any{"state": "incomplete", "grade": "Hypothesis"}
		if standard07FieldViewPresent {
			standard07FieldViewResult["raw"] = standard07FieldViewRaw
		} else {
			standard07FieldViewResult["absent"] = true
		}
		standard07FieldViewResults[standard07FieldViewKey] = standard07FieldViewResult
		standard07FieldViewArray, _ := standard07FieldViewRaw.([]any)
		if !standard07FieldViewPresent || len(standard07FieldViewArray) != standard07FieldViewWidth {
			standard07FieldViewOutput["state"] = "partial"
			continue
		}
		standard07FieldViewResult["state"] = "complete"
		standard07FieldViewFirst, _ := standard07FieldViewByte(standard07FieldViewArray[0])
		if standard07FieldViewDirection == "positive" {
			continue
		}
		switch standard07FieldViewKey {
		case "outdoor_temperature":
			standard07FieldViewSecond, _ := standard07FieldViewByte(standard07FieldViewArray[1])
			standard07FieldViewSigned := int(int16(uint16(standard07FieldViewFirst | standard07FieldViewSecond<<8)))
			standard07FieldViewResult["number"] = float64(standard07FieldViewSigned) / 256
			standard07FieldViewResult["unit"] = "Celsius"
			standard07FieldViewResult["calculation"] = "Proven"
			if standard07FieldViewSecondary == 2 {
				standard07FieldViewResult["candidate_domain"] = standard07FieldViewSigned >= -12800 && standard07FieldViewSigned <= 12800
			}
		case "weekday":
			standard07FieldViewLabels := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
			if standard07FieldViewFirst < 7 {
				standard07FieldViewResult["label"] = standard07FieldViewLabels[standard07FieldViewFirst]
			} else {
				standard07FieldViewResult["label"] = "unmapped"
			}
		case "duration":
			standard07FieldViewResult["number"] = standard07FieldViewFirst
			standard07FieldViewResult["unit"] = "minutes"
			standard07FieldViewResult["calculation"] = "Proven"
		case "second", "minute", "hour", "day", "month", "year":
			standard07FieldViewDomains := map[string][2]int{"second": {0, 59}, "minute": {0, 59}, "hour": {0, 23}, "day": {1, 31}, "month": {1, 12}, "year": {0, 99}}
			standard07FieldViewResult["candidate_domain"] = standard07FieldViewDomains[standard07FieldViewKey]
			standard07FieldViewResult["conversion"] = "Unknown"
		}
	}
	return json.Marshal(standard07FieldViewOutput)
}
