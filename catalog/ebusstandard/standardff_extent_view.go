package ebusstandard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func standardFFExtentViewValue(standardFFExtentViewDecoder *json.Decoder) error {
	standardFFExtentViewToken, standardFFExtentViewErr := standardFFExtentViewDecoder.Token()
	if standardFFExtentViewErr != nil {
		return standardFFExtentViewErr
	}
	if standardFFExtentViewDelimiter, standardFFExtentViewOK := standardFFExtentViewToken.(json.Delim); standardFFExtentViewOK {
		switch standardFFExtentViewDelimiter {
		case '{':
			standardFFExtentViewSeen := map[string]bool{}
			for standardFFExtentViewDecoder.More() {
				standardFFExtentViewKey, standardFFExtentViewErr := standardFFExtentViewDecoder.Token()
				if standardFFExtentViewErr != nil {
					return standardFFExtentViewErr
				}
				standardFFExtentViewName, standardFFExtentViewOK := standardFFExtentViewKey.(string)
				if !standardFFExtentViewOK || standardFFExtentViewSeen[standardFFExtentViewName] {
					return fmt.Errorf("duplicate or invalid object key")
				}
				standardFFExtentViewSeen[standardFFExtentViewName] = true
				if standardFFExtentViewErr := standardFFExtentViewValue(standardFFExtentViewDecoder); standardFFExtentViewErr != nil {
					return standardFFExtentViewErr
				}
			}
		case '[':
			for standardFFExtentViewDecoder.More() {
				if standardFFExtentViewErr := standardFFExtentViewValue(standardFFExtentViewDecoder); standardFFExtentViewErr != nil {
					return standardFFExtentViewErr
				}
			}
		default:
			return fmt.Errorf("unexpected delimiter")
		}
		_, standardFFExtentViewErr = standardFFExtentViewDecoder.Token()
	}
	return standardFFExtentViewErr
}

func standardFFExtentViewObject(standardFFExtentViewRaw json.RawMessage) (map[string]json.RawMessage, error) {
	standardFFExtentViewMap := map[string]json.RawMessage{}
	if len(bytes.TrimSpace(standardFFExtentViewRaw)) == 0 || bytes.TrimSpace(standardFFExtentViewRaw)[0] != '{' {
		return nil, fmt.Errorf("expected object")
	}
	standardFFExtentViewErr := json.Unmarshal(standardFFExtentViewRaw, &standardFFExtentViewMap)
	return standardFFExtentViewMap, standardFFExtentViewErr
}

func standardFFExtentViewByte(standardFFExtentViewRaw json.RawMessage) (int, error) {
	var standardFFExtentViewNumber int
	if len(standardFFExtentViewRaw) == 0 || bytes.Equal(bytes.TrimSpace(standardFFExtentViewRaw), []byte("null")) {
		return 0, fmt.Errorf("expected integer byte")
	}
	if standardFFExtentViewErr := json.Unmarshal(standardFFExtentViewRaw, &standardFFExtentViewNumber); standardFFExtentViewErr != nil {
		return 0, standardFFExtentViewErr
	}
	if standardFFExtentViewNumber < 0 || standardFFExtentViewNumber > 255 {
		return 0, fmt.Errorf("integer outside byte range")
	}
	return standardFFExtentViewNumber, nil
}

// StandardFFExtentViewEvaluate is the pure explicit-input evaluator for standard FF04 and
// FF06 body extents. It reports the declared and available body lengths
// independently, preserves every raw field and the opaque context, and grades
// only the extent arithmetic (calculation Proven); body semantics and physical
// support stay Unknown and the result is never qualified.
func StandardFFExtentViewEvaluate(standardFFExtentViewInputRaw json.RawMessage) (json.RawMessage, error) {
	standardFFExtentViewDecoder := json.NewDecoder(bytes.NewReader(standardFFExtentViewInputRaw))
	standardFFExtentViewDecoder.UseNumber()
	if standardFFExtentViewErr := standardFFExtentViewValue(standardFFExtentViewDecoder); standardFFExtentViewErr != nil {
		return nil, standardFFExtentViewErr
	}
	if _, standardFFExtentViewErr := standardFFExtentViewDecoder.Token(); standardFFExtentViewErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	standardFFExtentViewInput, standardFFExtentViewErr := standardFFExtentViewObject(standardFFExtentViewInputRaw)
	if standardFFExtentViewErr != nil {
		return nil, standardFFExtentViewErr
	}
	standardFFExtentViewContext, standardFFExtentViewErr := standardFFExtentViewObject(standardFFExtentViewInput["context"])
	if standardFFExtentViewErr != nil {
		return nil, standardFFExtentViewErr
	}
	standardFFExtentViewFields, standardFFExtentViewErr := standardFFExtentViewObject(standardFFExtentViewInput["fields"])
	if standardFFExtentViewErr != nil {
		return nil, standardFFExtentViewErr
	}
	standardFFExtentViewNumbers := map[string]int{}
	for _, standardFFExtentViewKey := range []string{"primary", "secondary", "operation"} {
		if standardFFExtentViewRaw, standardFFExtentViewOK := standardFFExtentViewContext[standardFFExtentViewKey]; standardFFExtentViewOK {
			standardFFExtentViewNumber, standardFFExtentViewErr := standardFFExtentViewByte(standardFFExtentViewRaw)
			if standardFFExtentViewErr != nil {
				return nil, standardFFExtentViewErr
			}
			standardFFExtentViewNumbers[standardFFExtentViewKey] = standardFFExtentViewNumber
		}
	}
	for _, standardFFExtentViewKey := range []string{"direction", "view"} {
		if standardFFExtentViewRaw, standardFFExtentViewOK := standardFFExtentViewContext[standardFFExtentViewKey]; standardFFExtentViewOK {
			var standardFFExtentViewString string
			if bytes.Equal(bytes.TrimSpace(standardFFExtentViewRaw), []byte("null")) {
				return nil, fmt.Errorf("expected context string")
			}
			if standardFFExtentViewErr := json.Unmarshal(standardFFExtentViewRaw, &standardFFExtentViewString); standardFFExtentViewErr != nil {
				return nil, standardFFExtentViewErr
			}
		}
	}
	standardFFExtentViewResults := map[string]any{}
	standardFFExtentViewLengths := map[string]int{}
	for standardFFExtentViewKey, standardFFExtentViewRaw := range standardFFExtentViewFields {
		if len(bytes.TrimSpace(standardFFExtentViewRaw)) == 0 || bytes.TrimSpace(standardFFExtentViewRaw)[0] != '[' {
			return nil, fmt.Errorf("expected byte array")
		}
		var standardFFExtentViewBytes []json.RawMessage
		if standardFFExtentViewErr := json.Unmarshal(standardFFExtentViewRaw, &standardFFExtentViewBytes); standardFFExtentViewErr != nil {
			return nil, standardFFExtentViewErr
		}
		for _, standardFFExtentViewByteRaw := range standardFFExtentViewBytes {
			if _, standardFFExtentViewErr := standardFFExtentViewByte(standardFFExtentViewByteRaw); standardFFExtentViewErr != nil {
				return nil, standardFFExtentViewErr
			}
		}
		standardFFExtentViewLengths[standardFFExtentViewKey] = len(standardFFExtentViewBytes)
		standardFFExtentViewResults[standardFFExtentViewKey] = map[string]any{"raw": standardFFExtentViewRaw, "state": "Unknown"}
	}
	standardFFExtentViewArguments := map[string]json.RawMessage{}
	if standardFFExtentViewRaw, standardFFExtentViewOK := standardFFExtentViewInput["arguments"]; standardFFExtentViewOK {
		standardFFExtentViewArguments, standardFFExtentViewErr = standardFFExtentViewObject(standardFFExtentViewRaw)
		if standardFFExtentViewErr != nil {
			return nil, standardFFExtentViewErr
		}
	}
	standardFFExtentViewDeclared := 0
	if standardFFExtentViewRaw, standardFFExtentViewOK := standardFFExtentViewArguments["declared_length"]; standardFFExtentViewOK {
		standardFFExtentViewDeclared, standardFFExtentViewErr = standardFFExtentViewByte(standardFFExtentViewRaw)
		if standardFFExtentViewErr != nil {
			return nil, standardFFExtentViewErr
		}
	}
	for _, standardFFExtentViewKey := range []string{"state", "results", "physical", "qualified"} {
		if _, standardFFExtentViewOK := standardFFExtentViewInput[standardFFExtentViewKey]; standardFFExtentViewOK {
			return nil, fmt.Errorf("input key collides with output key %q", standardFFExtentViewKey)
		}
	}
	standardFFExtentViewOutput := map[string]any{}
	for standardFFExtentViewKey, standardFFExtentViewRaw := range standardFFExtentViewInput {
		standardFFExtentViewOutput[standardFFExtentViewKey] = standardFFExtentViewRaw
	}
	standardFFExtentViewOutput["state"] = "Unknown"
	standardFFExtentViewOutput["results"] = standardFFExtentViewResults
	standardFFExtentViewOutput["physical"] = "Unknown"
	standardFFExtentViewOutput["qualified"] = false
	_, standardFFExtentViewOperation := standardFFExtentViewContext["operation"]
	var standardFFExtentViewDirection, standardFFExtentViewView string
	_ = json.Unmarshal(standardFFExtentViewContext["direction"], &standardFFExtentViewDirection)
	_ = json.Unmarshal(standardFFExtentViewContext["view"], &standardFFExtentViewView)
	if standardFFExtentViewNumbers["primary"] != 255 || (standardFFExtentViewNumbers["secondary"] != 4 && standardFFExtentViewNumbers["secondary"] != 6) || standardFFExtentViewOperation || standardFFExtentViewDirection != "positive" || standardFFExtentViewView != "body-extent" {
		return json.Marshal(standardFFExtentViewOutput)
	}
	if _, standardFFExtentViewOK := standardFFExtentViewArguments["declared_length"]; !standardFFExtentViewOK || len(standardFFExtentViewArguments) != 1 {
		return nil, fmt.Errorf("expected declared_length argument only")
	}
	standardFFExtentViewAvailable, standardFFExtentViewPresent := standardFFExtentViewLengths["body"]
	standardFFExtentViewBody := map[string]any{"state": "Unknown", "calculation": "Proven", "declared_length": standardFFExtentViewDeclared}
	if standardFFExtentViewPresent {
		standardFFExtentViewBody["raw"] = standardFFExtentViewFields["body"]
		standardFFExtentViewBody["available_length"] = standardFFExtentViewAvailable
	} else {
		standardFFExtentViewBody["state"] = "partial"
	}
	standardFFExtentViewFlags := map[string]bool{}
	standardFFExtentViewComplete := false
	if standardFFExtentViewNumbers["secondary"] == 4 {
		standardFFExtentViewFlags["declared_minimum_satisfied"] = standardFFExtentViewDeclared >= 2
		standardFFExtentViewFlags["initialized_minimum_satisfied"] = standardFFExtentViewPresent && standardFFExtentViewAvailable >= 2
		standardFFExtentViewComplete = standardFFExtentViewDeclared >= 2 && standardFFExtentViewDeclared <= 10 && standardFFExtentViewPresent && standardFFExtentViewAvailable >= standardFFExtentViewDeclared
		standardFFExtentViewFlags["bounded_minimum"] = standardFFExtentViewComplete
	} else {
		standardFFExtentViewTail := standardFFExtentViewDeclared >= 1 && (standardFFExtentViewDeclared-1)%2 == 1
		standardFFExtentViewFlags["declared_incomplete_tail"] = standardFFExtentViewTail
		if standardFFExtentViewDeclared >= 1 {
			standardFFExtentViewBody["declared_pair_count"] = (standardFFExtentViewDeclared - 1) / 2
		}
		standardFFExtentViewBounded := standardFFExtentViewAvailable
		if standardFFExtentViewBounded > standardFFExtentViewDeclared {
			standardFFExtentViewBounded = standardFFExtentViewDeclared
		}
		if standardFFExtentViewPresent && standardFFExtentViewBounded >= 1 {
			standardFFExtentViewBody["available_pair_count"] = (standardFFExtentViewBounded - 1) / 2
		}
		standardFFExtentViewComplete = standardFFExtentViewDeclared >= 1 && standardFFExtentViewPresent && standardFFExtentViewAvailable >= standardFFExtentViewDeclared && !standardFFExtentViewTail
		standardFFExtentViewFlags["whole_declared_span"] = standardFFExtentViewComplete
	}
	standardFFExtentViewBody["flags"] = standardFFExtentViewFlags
	standardFFExtentViewResults["body"] = standardFFExtentViewBody
	standardFFExtentViewOutput["state"] = "partial"
	if standardFFExtentViewComplete {
		standardFFExtentViewOutput["state"] = "complete"
	}
	return json.Marshal(standardFFExtentViewOutput)
}
