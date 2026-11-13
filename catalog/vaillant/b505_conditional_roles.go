package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"strconv"
	"strings"
)

const b505ConditionalRoleMaxDepth = 64

func b505ConditionalRoleRead(b505ConditionalRoleDecoder *json.Decoder, b505ConditionalRoleDepth int) (any, error) {
	if b505ConditionalRoleDepth > b505ConditionalRoleMaxDepth {
		return nil, fmt.Errorf("nesting too deep")
	}
	b505ConditionalRoleToken, b505ConditionalRoleError := b505ConditionalRoleDecoder.Token()
	if b505ConditionalRoleError != nil {
		return nil, b505ConditionalRoleError
	}
	switch b505ConditionalRoleValue := b505ConditionalRoleToken.(type) {
	case json.Delim:
		switch b505ConditionalRoleValue {
		case '{':
			b505ConditionalRoleObject := map[string]any{}
			for b505ConditionalRoleDecoder.More() {
				b505ConditionalRoleKey, b505ConditionalRoleErr := b505ConditionalRoleDecoder.Token()
				if b505ConditionalRoleErr != nil {
					return nil, b505ConditionalRoleErr
				}
				b505ConditionalRoleName := b505ConditionalRoleKey.(string)
				if _, b505ConditionalRoleExists := b505ConditionalRoleObject[b505ConditionalRoleName]; b505ConditionalRoleExists {
					return nil, fmt.Errorf("duplicate key %q", b505ConditionalRoleName)
				}
				b505ConditionalRoleChild, b505ConditionalRoleErr := b505ConditionalRoleRead(b505ConditionalRoleDecoder, b505ConditionalRoleDepth+1)
				if b505ConditionalRoleErr != nil {
					return nil, b505ConditionalRoleErr
				}
				b505ConditionalRoleObject[b505ConditionalRoleName] = b505ConditionalRoleChild
			}
			_, b505ConditionalRoleError = b505ConditionalRoleDecoder.Token()
			return b505ConditionalRoleObject, b505ConditionalRoleError
		case '[':
			b505ConditionalRoleArray := []any{}
			for b505ConditionalRoleDecoder.More() {
				b505ConditionalRoleChild, b505ConditionalRoleErr := b505ConditionalRoleRead(b505ConditionalRoleDecoder, b505ConditionalRoleDepth+1)
				if b505ConditionalRoleErr != nil {
					return nil, b505ConditionalRoleErr
				}
				b505ConditionalRoleArray = append(b505ConditionalRoleArray, b505ConditionalRoleChild)
			}
			_, b505ConditionalRoleError = b505ConditionalRoleDecoder.Token()
			return b505ConditionalRoleArray, b505ConditionalRoleError
		}
	case json.Number:
		if b505ConditionalRoleErr := b505ConditionalRoleBound(string(b505ConditionalRoleValue)); b505ConditionalRoleErr != nil {
			return nil, b505ConditionalRoleErr
		}
		b505ConditionalRoleNumber, b505ConditionalRoleErr := b505ConditionalRoleValue.Float64()
		if b505ConditionalRoleErr != nil || math.IsInf(b505ConditionalRoleNumber, 0) || math.IsNaN(b505ConditionalRoleNumber) {
			return nil, fmt.Errorf("nonfinite number")
		}
	}
	return b505ConditionalRoleToken, nil
}

// b505ConditionalRoleBound rejects numeric tokens whose length or decimal exponent would
// make exact rational arithmetic and rendering unbounded.
func b505ConditionalRoleBound(b505ConditionalRoleToken string) error {
	if len(b505ConditionalRoleToken) > 64 {
		return fmt.Errorf("number token too long")
	}
	if b505ConditionalRoleIndex := strings.IndexAny(b505ConditionalRoleToken, "eE"); b505ConditionalRoleIndex >= 0 {
		b505ConditionalRoleExponent, b505ConditionalRoleErr := strconv.Atoi(b505ConditionalRoleToken[b505ConditionalRoleIndex+1:])
		if b505ConditionalRoleErr != nil || b505ConditionalRoleExponent > 400 || b505ConditionalRoleExponent < -400 {
			return fmt.Errorf("number exponent out of range")
		}
	}
	return nil
}

func b505ConditionalRoleByte(b505ConditionalRoleValue any) (int, bool) {
	b505ConditionalRoleNumber, b505ConditionalRoleOK := b505ConditionalRoleValue.(json.Number)
	if !b505ConditionalRoleOK {
		return 0, false
	}
	b505ConditionalRoleRational, b505ConditionalRoleValid := new(big.Rat).SetString(string(b505ConditionalRoleNumber))
	if !b505ConditionalRoleValid || !b505ConditionalRoleRational.IsInt() || !b505ConditionalRoleRational.Num().IsInt64() {
		return 0, false
	}
	b505ConditionalRoleInteger := b505ConditionalRoleRational.Num().Int64()
	return int(b505ConditionalRoleInteger), b505ConditionalRoleInteger >= 0 && b505ConditionalRoleInteger <= 255
}

func B505ConditionalRoleEvaluate(b505ConditionalRoleInput json.RawMessage) (json.RawMessage, error) {
	b505ConditionalRoleDecoder := json.NewDecoder(bytes.NewReader(b505ConditionalRoleInput))
	b505ConditionalRoleDecoder.UseNumber()
	b505ConditionalRoleParsed, b505ConditionalRoleError := b505ConditionalRoleRead(b505ConditionalRoleDecoder, 0)
	if b505ConditionalRoleError != nil {
		return nil, b505ConditionalRoleError
	}
	if _, b505ConditionalRoleError = b505ConditionalRoleDecoder.Token(); b505ConditionalRoleError != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b505ConditionalRoleOutput, b505ConditionalRoleOK := b505ConditionalRoleParsed.(map[string]any)
	if !b505ConditionalRoleOK {
		return nil, fmt.Errorf("input must be object")
	}
	b505ConditionalRoleContext, b505ConditionalRoleOK := b505ConditionalRoleOutput["context"].(map[string]any)
	if !b505ConditionalRoleOK {
		return nil, fmt.Errorf("context must be object")
	}
	b505ConditionalRoleFields, b505ConditionalRoleOK := b505ConditionalRoleOutput["fields"].(map[string]any)
	if !b505ConditionalRoleOK {
		return nil, fmt.Errorf("fields must be object")
	}
	b505ConditionalRoleArguments := map[string]any{}
	if b505ConditionalRoleValue, b505ConditionalRolePresent := b505ConditionalRoleOutput["arguments"]; b505ConditionalRolePresent {
		b505ConditionalRoleArguments, b505ConditionalRoleOK = b505ConditionalRoleValue.(map[string]any)
		if !b505ConditionalRoleOK {
			return nil, fmt.Errorf("arguments must be object")
		}
	}
	for _, b505ConditionalRoleKey := range []string{"primary", "secondary", "operation"} {
		if b505ConditionalRoleValue, b505ConditionalRolePresent := b505ConditionalRoleContext[b505ConditionalRoleKey]; b505ConditionalRolePresent {
			if _, b505ConditionalRoleOK = b505ConditionalRoleByte(b505ConditionalRoleValue); !b505ConditionalRoleOK {
				return nil, fmt.Errorf("invalid %s", b505ConditionalRoleKey)
			}
		}
	}
	for _, b505ConditionalRoleKey := range []string{"direction", "view", "family", "variant"} {
		if b505ConditionalRoleValue, b505ConditionalRolePresent := b505ConditionalRoleContext[b505ConditionalRoleKey]; b505ConditionalRolePresent {
			if _, b505ConditionalRoleOK = b505ConditionalRoleValue.(string); !b505ConditionalRoleOK {
				return nil, fmt.Errorf("invalid %s", b505ConditionalRoleKey)
			}
		}
	}
	for _, b505ConditionalRoleReserved := range []string{"state", "results", "physical", "qualified"} {
		if _, b505ConditionalRolePresent := b505ConditionalRoleOutput[b505ConditionalRoleReserved]; b505ConditionalRolePresent {
			return nil, fmt.Errorf("reserved top-level key %q", b505ConditionalRoleReserved)
		}
	}
	b505ConditionalRoleResults := map[string]any{}
	for b505ConditionalRoleName, b505ConditionalRoleRaw := range b505ConditionalRoleFields {
		b505ConditionalRoleArray, b505ConditionalRoleValid := b505ConditionalRoleRaw.([]any)
		if !b505ConditionalRoleValid {
			return nil, fmt.Errorf("field %s must be byte array", b505ConditionalRoleName)
		}
		for _, b505ConditionalRoleValue := range b505ConditionalRoleArray {
			if _, b505ConditionalRoleValid = b505ConditionalRoleByte(b505ConditionalRoleValue); !b505ConditionalRoleValid {
				return nil, fmt.Errorf("invalid field byte")
			}
		}
		b505ConditionalRoleResults[b505ConditionalRoleName] = map[string]any{"raw": b505ConditionalRoleRaw, "state": "Unknown"}
	}
	b505ConditionalRoleOutput["state"] = "Unknown"
	b505ConditionalRoleOutput["results"] = b505ConditionalRoleResults
	b505ConditionalRoleOutput["physical"] = "Unknown"
	b505ConditionalRoleOutput["qualified"] = false
	b505ConditionalRolePrimary, _ := b505ConditionalRoleByte(b505ConditionalRoleContext["primary"])
	b505ConditionalRoleSecondary, _ := b505ConditionalRoleByte(b505ConditionalRoleContext["secondary"])
	b505ConditionalRoleOperation, b505ConditionalRoleHasOperation := b505ConditionalRoleByte(b505ConditionalRoleContext["operation"])
	b505ConditionalRoleView, _ := b505ConditionalRoleContext["view"].(string)
	b505ConditionalRoleSupported := b505ConditionalRolePrimary == 181 && b505ConditionalRoleSecondary == 5 && b505ConditionalRoleContext["direction"] == "request" && b505ConditionalRoleHasOperation
	b505ConditionalRoleSupported = b505ConditionalRoleSupported && ((b505ConditionalRoleView == "unchanged" && (b505ConditionalRoleOperation == 90 || b505ConditionalRoleOperation == 61)) || (b505ConditionalRoleView == "dhw-enable" && b505ConditionalRoleOperation == 49) || (b505ConditionalRoleView == "parallel-charge" && (b505ConditionalRoleOperation == 24 || b505ConditionalRoleOperation == 39)) || (b505ConditionalRoleView == "flow-excess" && b505ConditionalRoleOperation == 23) || (b505ConditionalRoleView == "programmed-start" && b505ConditionalRoleOperation == 9))
	if !b505ConditionalRoleSupported {
		return json.Marshal(b505ConditionalRoleOutput)
	}
	b505ConditionalRoleOutput["state"] = "complete"
	b505ConditionalRoleResult := map[string]any{"state": "complete", "grade": "Hypothesis"}
	b505ConditionalRoleName := b505ConditionalRoleView
	switch b505ConditionalRoleView {
	case "unchanged":
		b505ConditionalRoleReplacement, b505ConditionalRoleValid := b505ConditionalRoleArguments["replacement"].(bool)
		if !b505ConditionalRoleValid {
			return nil, fmt.Errorf("replacement must be boolean")
		}
		for _, b505ConditionalRoleKey := range []string{"previous", "requested"} {
			b505ConditionalRoleValue, b505ConditionalRolePresent := b505ConditionalRoleArguments[b505ConditionalRoleKey]
			if !b505ConditionalRolePresent {
				return nil, fmt.Errorf("missing %s", b505ConditionalRoleKey)
			}
			switch b505ConditionalRoleValue.(type) {
			case map[string]any, []any:
				return nil, fmt.Errorf("%s must be scalar", b505ConditionalRoleKey)
			}
		}
		b505ConditionalRoleResult["raw"] = b505ConditionalRoleArguments["requested"]
		b505ConditionalRoleResult["value"] = b505ConditionalRoleArguments["requested"]
		b505ConditionalRoleResult["label"] = "numeric-candidate"
		b505ConditionalRoleResult["replacement_code"] = "Unknown"
		if b505ConditionalRoleReplacement {
			b505ConditionalRoleResult["value"] = b505ConditionalRoleArguments["previous"]
			b505ConditionalRoleResult["label"] = "unchanged"
		}
		if b505ConditionalRoleNumber, b505ConditionalRoleNumeric := b505ConditionalRoleResult["value"].(json.Number); b505ConditionalRoleNumeric {
			b505ConditionalRoleResult["number"] = b505ConditionalRoleNumber
		}
	case "flow-excess":
		b505ConditionalRoleDemand, b505ConditionalRoleValid := b505ConditionalRoleArguments["demand"].(json.Number)
		if !b505ConditionalRoleValid {
			return nil, fmt.Errorf("demand must be numeric")
		}
		b505ConditionalRoleExcess, b505ConditionalRoleValid := b505ConditionalRoleArguments["excess"].(json.Number)
		if !b505ConditionalRoleValid {
			return nil, fmt.Errorf("excess must be numeric")
		}
		b505ConditionalRoleLeft, b505ConditionalRoleValid := new(big.Rat).SetString(string(b505ConditionalRoleDemand))
		if !b505ConditionalRoleValid {
			return nil, fmt.Errorf("invalid demand")
		}
		b505ConditionalRoleRight, b505ConditionalRoleValid := new(big.Rat).SetString(string(b505ConditionalRoleExcess))
		if !b505ConditionalRoleValid {
			return nil, fmt.Errorf("invalid excess")
		}
		b505ConditionalRoleSum := new(big.Rat).Add(b505ConditionalRoleLeft, b505ConditionalRoleRight)
		b505ConditionalRoleDenominator := new(big.Int).Set(b505ConditionalRoleSum.Denom())
		b505ConditionalRoleTwos, b505ConditionalRoleFives := 0, 0
		for b505ConditionalRoleDenominator.Cmp(big.NewInt(1)) != 0 {
			b505ConditionalRoleDivisor := int64(2)
			if b505ConditionalRoleDenominator.Bit(0) != 0 {
				b505ConditionalRoleDivisor = 5
				b505ConditionalRoleFives++
			} else {
				b505ConditionalRoleTwos++
			}
			b505ConditionalRoleDenominator.Quo(b505ConditionalRoleDenominator, big.NewInt(b505ConditionalRoleDivisor))
		}
		b505ConditionalRoleScale := max(b505ConditionalRoleTwos, b505ConditionalRoleFives)
		b505ConditionalRoleResult["raw"] = map[string]any{"demand": b505ConditionalRoleDemand, "excess": b505ConditionalRoleExcess}
		b505ConditionalRoleResult["number"] = json.Number(b505ConditionalRoleSum.FloatString(b505ConditionalRoleScale))
		b505ConditionalRoleResult["calculation"] = "Proven"
	default:
		b505ConditionalRoleName = "enabled"
		if b505ConditionalRoleView == "parallel-charge" {
			b505ConditionalRoleName = "disable"
		}
		if b505ConditionalRoleView == "programmed-start" {
			b505ConditionalRoleName = "start"
		}
		b505ConditionalRoleRole := ""
		b505ConditionalRoleParallel := false
		b505ConditionalRoleMode := 0
		if b505ConditionalRoleView == "parallel-charge" {
			b505ConditionalRoleRole, b505ConditionalRoleOK = b505ConditionalRoleArguments["role"].(string)
			if !b505ConditionalRoleOK || (b505ConditionalRoleRole != "mixer" && b505ConditionalRoleRole != "heating" && b505ConditionalRoleRole != "dhw") {
				return nil, fmt.Errorf("invalid role")
			}
			b505ConditionalRoleParallel, b505ConditionalRoleOK = b505ConditionalRoleArguments["parallel"].(bool)
			if !b505ConditionalRoleOK {
				return nil, fmt.Errorf("parallel must be boolean")
			}
		}
		if b505ConditionalRoleView == "programmed-start" {
			b505ConditionalRoleMode, b505ConditionalRoleOK = b505ConditionalRoleByte(b505ConditionalRoleArguments["program_mode"])
			if !b505ConditionalRoleOK {
				return nil, fmt.Errorf("invalid program_mode")
			}
		}
		b505ConditionalRoleRaw, b505ConditionalRolePresent := b505ConditionalRoleFields[b505ConditionalRoleName]
		if b505ConditionalRolePresent {
			b505ConditionalRoleResult["raw"] = b505ConditionalRoleRaw
		}
		b505ConditionalRoleArray, _ := b505ConditionalRoleRaw.([]any)
		if !b505ConditionalRolePresent || (b505ConditionalRoleView != "programmed-start" && len(b505ConditionalRoleArray) != 1) {
			b505ConditionalRoleResult["state"] = "incomplete"
			b505ConditionalRoleOutput["state"] = "partial"
			break
		}
		if b505ConditionalRoleView == "programmed-start" {
			b505ConditionalRoleResult["label"] = "ignored"
			if b505ConditionalRoleMode == 0 {
				b505ConditionalRoleResult["label"] = "applicable"
			}
			break
		}
		b505ConditionalRoleByteValue, _ := b505ConditionalRoleByte(b505ConditionalRoleArray[0])
		b505ConditionalRoleResult["label"] = "unmapped"
		if b505ConditionalRoleView == "dhw-enable" {
			if b505ConditionalRoleByteValue == 0 {
				b505ConditionalRoleResult["label"] = "frost-only-eligible"
			}
			if b505ConditionalRoleByteValue == 1 {
				b505ConditionalRoleResult["label"] = "ordinary-DHW-eligible"
			}
			break
		}
		switch {
		case b505ConditionalRoleRole == "dhw":
			b505ConditionalRoleResult["label"] = "not-applicable"
		case b505ConditionalRoleByteValue == 0:
			b505ConditionalRoleResult["label"] = "suppression-not-requested"
		case b505ConditionalRoleByteValue == 1:
			b505ConditionalRoleResult["label"] = "heating-suppressed"
			if b505ConditionalRoleParallel {
				b505ConditionalRoleResult["label"] = "exemption-unknown"
				if b505ConditionalRoleRole == "mixer" {
					b505ConditionalRoleResult["label"] = "mixer-parallel-exception"
				}
			}
		}
	}
	if b505ConditionalRolePrior, b505ConditionalRoleExists := b505ConditionalRoleResults[b505ConditionalRoleName].(map[string]any); b505ConditionalRoleExists && b505ConditionalRoleName == b505ConditionalRoleView {
		b505ConditionalRoleResult["field_raw"] = b505ConditionalRolePrior["raw"]
	}
	b505ConditionalRoleResults[b505ConditionalRoleName] = b505ConditionalRoleResult
	return json.Marshal(b505ConditionalRoleOutput)
}
