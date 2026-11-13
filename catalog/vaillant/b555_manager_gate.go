package vaillant

import (
	b555ManagerGateBytesPackage "bytes"
	b555ManagerGateJSON "encoding/json"
	b555ManagerGateFmt "fmt"
	b555ManagerGateIO "io"
	b555ManagerGateStrconv "strconv"
)

func b555ManagerGateValue(b555ManagerGateDecoder *b555ManagerGateJSON.Decoder) (any, error) {
	b555ManagerGateToken, b555ManagerGateErr := b555ManagerGateDecoder.Token()
	if b555ManagerGateErr != nil {
		return nil, b555ManagerGateErr
	}
	switch b555ManagerGateToken {
	case b555ManagerGateJSON.Delim('{'):
		b555ManagerGateObject := map[string]any{}
		for b555ManagerGateDecoder.More() {
			b555ManagerGateKeyToken, b555ManagerGateErr := b555ManagerGateDecoder.Token()
			if b555ManagerGateErr != nil {
				return nil, b555ManagerGateErr
			}
			b555ManagerGateKey := b555ManagerGateKeyToken.(string)
			if _, b555ManagerGateExists := b555ManagerGateObject[b555ManagerGateKey]; b555ManagerGateExists {
				return nil, b555ManagerGateFmt.Errorf("duplicate key %q", b555ManagerGateKey)
			}
			b555ManagerGateChild, b555ManagerGateErr := b555ManagerGateValue(b555ManagerGateDecoder)
			if b555ManagerGateErr != nil {
				return nil, b555ManagerGateErr
			}
			b555ManagerGateObject[b555ManagerGateKey] = b555ManagerGateChild
		}
		_, b555ManagerGateErr := b555ManagerGateDecoder.Token()
		return b555ManagerGateObject, b555ManagerGateErr
	case b555ManagerGateJSON.Delim('['):
		b555ManagerGateArray := []any{}
		for b555ManagerGateDecoder.More() {
			b555ManagerGateChild, b555ManagerGateErr := b555ManagerGateValue(b555ManagerGateDecoder)
			if b555ManagerGateErr != nil {
				return nil, b555ManagerGateErr
			}
			b555ManagerGateArray = append(b555ManagerGateArray, b555ManagerGateChild)
		}
		_, b555ManagerGateErr := b555ManagerGateDecoder.Token()
		return b555ManagerGateArray, b555ManagerGateErr
	}
	return b555ManagerGateToken, nil
}

func b555ManagerGateUnsigned(b555ManagerGateValue any, b555ManagerGateBits int) (uint64, error) {
	b555ManagerGateNumber, b555ManagerGateOK := b555ManagerGateValue.(b555ManagerGateJSON.Number)
	if !b555ManagerGateOK {
		return 0, b555ManagerGateFmt.Errorf("expected unsigned integer")
	}
	return b555ManagerGateStrconv.ParseUint(string(b555ManagerGateNumber), 10, b555ManagerGateBits)
}

func B555ManagerGateEvaluate(b555ManagerGateInput b555ManagerGateJSON.RawMessage) (b555ManagerGateJSON.RawMessage, error) {
	b555ManagerGateDecoder := b555ManagerGateJSON.NewDecoder(b555ManagerGateBytesPackage.NewReader(b555ManagerGateInput))
	b555ManagerGateDecoder.UseNumber()
	b555ManagerGateRoot, b555ManagerGateErr := b555ManagerGateValue(b555ManagerGateDecoder)
	if b555ManagerGateErr != nil {
		return nil, b555ManagerGateErr
	}
	if _, b555ManagerGateErr = b555ManagerGateDecoder.Token(); b555ManagerGateErr != b555ManagerGateIO.EOF {
		return nil, b555ManagerGateFmt.Errorf("trailing JSON")
	}
	b555ManagerGateObject, b555ManagerGateOK := b555ManagerGateRoot.(map[string]any)
	if !b555ManagerGateOK {
		return nil, b555ManagerGateFmt.Errorf("expected input object")
	}
	b555ManagerGateContext, b555ManagerGateOK := b555ManagerGateObject["context"].(map[string]any)
	if !b555ManagerGateOK {
		return nil, b555ManagerGateFmt.Errorf("expected context object")
	}
	b555ManagerGateRaw, b555ManagerGateOK := b555ManagerGateObject["fields"].(map[string]any)
	if !b555ManagerGateOK {
		return nil, b555ManagerGateFmt.Errorf("expected fields object")
	}
	for _, b555ManagerGateField := range b555ManagerGateRaw {
		b555ManagerGateArray, b555ManagerGateOK := b555ManagerGateField.([]any)
		if !b555ManagerGateOK {
			return nil, b555ManagerGateFmt.Errorf("expected byte array")
		}
		for _, b555ManagerGateByte := range b555ManagerGateArray {
			if _, b555ManagerGateErr := b555ManagerGateUnsigned(b555ManagerGateByte, 8); b555ManagerGateErr != nil {
				return nil, b555ManagerGateErr
			}
		}
	}
	b555ManagerGateArguments := map[string]any{}
	if b555ManagerGateValue, b555ManagerGatePresent := b555ManagerGateObject["arguments"]; b555ManagerGatePresent {
		b555ManagerGateArguments, b555ManagerGateOK = b555ManagerGateValue.(map[string]any)
		if !b555ManagerGateOK {
			return nil, b555ManagerGateFmt.Errorf("expected arguments object")
		}
	}
	for b555ManagerGateKey, b555ManagerGateBits := range map[string]int{"controller_type": 8, "spn": 16} {
		if b555ManagerGateValue, b555ManagerGatePresent := b555ManagerGateArguments[b555ManagerGateKey]; b555ManagerGatePresent {
			if _, b555ManagerGateErr := b555ManagerGateUnsigned(b555ManagerGateValue, b555ManagerGateBits); b555ManagerGateErr != nil {
				return nil, b555ManagerGateErr
			}
		}
	}
	if b555ManagerGateFamily, b555ManagerGatePresent := b555ManagerGateArguments["family_match"]; b555ManagerGatePresent && b555ManagerGateFamily != nil {
		if _, b555ManagerGateOK := b555ManagerGateFamily.(bool); !b555ManagerGateOK {
			return nil, b555ManagerGateFmt.Errorf("expected Boolean or null family_match")
		}
	}
	b555ManagerGateSupported := true
	for b555ManagerGateKey, b555ManagerGateExpected := range map[string]uint64{"primary": 181, "secondary": 85, "operation": 0} {
		b555ManagerGateValue := b555ManagerGateContext[b555ManagerGateKey]
		if b555ManagerGateValue == nil {
			b555ManagerGateSupported = false
			continue
		}
		b555ManagerGateNumber, b555ManagerGateErr := b555ManagerGateUnsigned(b555ManagerGateValue, 8)
		if b555ManagerGateErr != nil {
			return nil, b555ManagerGateErr
		}
		b555ManagerGateSupported = b555ManagerGateSupported && b555ManagerGateNumber == b555ManagerGateExpected
	}
	for b555ManagerGateKey, b555ManagerGateExpected := range map[string]string{"direction": "local", "view": "manager-gate"} {
		b555ManagerGateValue := b555ManagerGateContext[b555ManagerGateKey]
		if b555ManagerGateValue == nil {
			b555ManagerGateSupported = false
			continue
		}
		b555ManagerGateString, b555ManagerGateOK := b555ManagerGateValue.(string)
		if !b555ManagerGateOK {
			return nil, b555ManagerGateFmt.Errorf("expected string selector")
		}
		b555ManagerGateSupported = b555ManagerGateSupported && b555ManagerGateString == b555ManagerGateExpected
	}
	b555ManagerGateFields := map[string]any{}
	b555ManagerGateOutput := map[string]any{"context": b555ManagerGateContext, "raw": b555ManagerGateRaw, "state": "Unknown", "fields": b555ManagerGateFields, "physical": "Unknown", "qualified": false}
	if b555ManagerGateValue, b555ManagerGatePresent := b555ManagerGateObject["arguments"]; b555ManagerGatePresent {
		b555ManagerGateOutput["arguments"] = b555ManagerGateValue
	}
	if !b555ManagerGateSupported {
		return b555ManagerGateJSON.Marshal(b555ManagerGateOutput)
	}
	for _, b555ManagerGateKey := range []string{"controller_type", "spn", "family_match"} {
		if _, b555ManagerGatePresent := b555ManagerGateArguments[b555ManagerGateKey]; !b555ManagerGatePresent {
			return nil, b555ManagerGateFmt.Errorf("missing argument %s", b555ManagerGateKey)
		}
	}
	b555ManagerGateReply := map[string]any{"state": "incomplete"}
	if b555ManagerGateValue, b555ManagerGatePresent := b555ManagerGateRaw["reply"]; b555ManagerGatePresent {
		b555ManagerGateReply["raw"] = b555ManagerGateValue
	}
	b555ManagerGateFields["reply"] = b555ManagerGateReply
	b555ManagerGateOutput["state"] = "partial"
	b555ManagerGateArray, _ := b555ManagerGateRaw["reply"].([]any)
	if len(b555ManagerGateArray) != 1 {
		return b555ManagerGateJSON.Marshal(b555ManagerGateOutput)
	}
	b555ManagerGateByte, _ := b555ManagerGateUnsigned(b555ManagerGateArray[0], 8)
	b555ManagerGateType, _ := b555ManagerGateUnsigned(b555ManagerGateArguments["controller_type"], 8)
	b555ManagerGateSPN, _ := b555ManagerGateUnsigned(b555ManagerGateArguments["spn"], 16)
	b555ManagerGateBit0 := b555ManagerGateByte&1 != 0
	b555ManagerGateBit1 := b555ManagerGateByte&2 != 0
	b555ManagerGateFamily, _ := b555ManagerGateArguments["family_match"].(bool)
	b555ManagerGateReply["state"] = "Hypothesis"
	b555ManagerGateReply["calculation"] = "Proven"
	b555ManagerGateReply["flags"] = map[string]any{"bit0": b555ManagerGateBit0, "bit1": b555ManagerGateBit1}
	b555ManagerGateReply["unmodeled"] = b555ManagerGateByte & 252
	b555ManagerGateReply["eligible"] = b555ManagerGateType == 2 && b555ManagerGateSPN == 385 && b555ManagerGateFamily && b555ManagerGateBit0 && b555ManagerGateBit1
	b555ManagerGateReply["first_bit_policy"] = b555ManagerGateBit0
	b555ManagerGateReply["first_bit_policy_label"] = "caller predicate"
	b555ManagerGateReply["family_role"] = 0x15
	b555ManagerGateReply["route_role"] = 0x08
	b555ManagerGateReply["family_applicability"] = "Hypothesis"
	if b555ManagerGateArguments["family_match"] == nil {
		b555ManagerGateReply["family_applicability"] = "Unknown"
	}
	b555ManagerGateReply["unknown"] = []string{"exact family prefix", "reply adaptation", "actual manager state", "physical support"}
	b555ManagerGateOutput["state"] = "complete"
	return b555ManagerGateJSON.Marshal(b555ManagerGateOutput)
}
