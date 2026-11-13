package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func b512HeatpumpIOEvaluatorRead(b512HeatpumpIOEvaluatorDecoder *json.Decoder) (any, error) {
	b512HeatpumpIOEvaluatorToken, b512HeatpumpIOEvaluatorErr := b512HeatpumpIOEvaluatorDecoder.Token()
	if b512HeatpumpIOEvaluatorErr != nil {
		return nil, b512HeatpumpIOEvaluatorErr
	}
	if b512HeatpumpIOEvaluatorDelimiter, b512HeatpumpIOEvaluatorOK := b512HeatpumpIOEvaluatorToken.(json.Delim); b512HeatpumpIOEvaluatorOK {
		switch b512HeatpumpIOEvaluatorDelimiter {
		case '{':
			b512HeatpumpIOEvaluatorObject := map[string]any{}
			for b512HeatpumpIOEvaluatorDecoder.More() {
				b512HeatpumpIOEvaluatorKeyToken, b512HeatpumpIOEvaluatorErr := b512HeatpumpIOEvaluatorDecoder.Token()
				if b512HeatpumpIOEvaluatorErr != nil {
					return nil, b512HeatpumpIOEvaluatorErr
				}
				b512HeatpumpIOEvaluatorKey, b512HeatpumpIOEvaluatorOK := b512HeatpumpIOEvaluatorKeyToken.(string)
				if !b512HeatpumpIOEvaluatorOK {
					return nil, fmt.Errorf("object key required")
				}
				if _, b512HeatpumpIOEvaluatorExists := b512HeatpumpIOEvaluatorObject[b512HeatpumpIOEvaluatorKey]; b512HeatpumpIOEvaluatorExists {
					return nil, fmt.Errorf("duplicate key %q", b512HeatpumpIOEvaluatorKey)
				}
				b512HeatpumpIOEvaluatorValue, b512HeatpumpIOEvaluatorErr := b512HeatpumpIOEvaluatorRead(b512HeatpumpIOEvaluatorDecoder)
				if b512HeatpumpIOEvaluatorErr != nil {
					return nil, b512HeatpumpIOEvaluatorErr
				}
				b512HeatpumpIOEvaluatorObject[b512HeatpumpIOEvaluatorKey] = b512HeatpumpIOEvaluatorValue
			}
			_, b512HeatpumpIOEvaluatorErr = b512HeatpumpIOEvaluatorDecoder.Token()
			return b512HeatpumpIOEvaluatorObject, b512HeatpumpIOEvaluatorErr
		case '[':
			b512HeatpumpIOEvaluatorArray := []any{}
			for b512HeatpumpIOEvaluatorDecoder.More() {
				b512HeatpumpIOEvaluatorValue, b512HeatpumpIOEvaluatorErr := b512HeatpumpIOEvaluatorRead(b512HeatpumpIOEvaluatorDecoder)
				if b512HeatpumpIOEvaluatorErr != nil {
					return nil, b512HeatpumpIOEvaluatorErr
				}
				b512HeatpumpIOEvaluatorArray = append(b512HeatpumpIOEvaluatorArray, b512HeatpumpIOEvaluatorValue)
			}
			_, b512HeatpumpIOEvaluatorErr = b512HeatpumpIOEvaluatorDecoder.Token()
			return b512HeatpumpIOEvaluatorArray, b512HeatpumpIOEvaluatorErr
		default:
			return nil, fmt.Errorf("unexpected delimiter")
		}
	}
	return b512HeatpumpIOEvaluatorToken, nil
}

func b512HeatpumpIOEvaluatorUint(b512HeatpumpIOEvaluatorValue any, b512HeatpumpIOEvaluatorBits int) (uint64, error) {
	b512HeatpumpIOEvaluatorNumber, b512HeatpumpIOEvaluatorOK := b512HeatpumpIOEvaluatorValue.(json.Number)
	if !b512HeatpumpIOEvaluatorOK {
		return 0, fmt.Errorf("unsigned integer required")
	}
	return strconv.ParseUint(string(b512HeatpumpIOEvaluatorNumber), 10, b512HeatpumpIOEvaluatorBits)
}

func b512HeatpumpIOEvaluatorValidateArguments(b512HeatpumpIOEvaluatorValue any) error {
	switch b512HeatpumpIOEvaluatorItem := b512HeatpumpIOEvaluatorValue.(type) {
	case json.Number:
		if _, b512HeatpumpIOEvaluatorErr := b512HeatpumpIOEvaluatorUint(b512HeatpumpIOEvaluatorItem, 32); b512HeatpumpIOEvaluatorErr != nil {
			return fmt.Errorf("argument number: %w", b512HeatpumpIOEvaluatorErr)
		}
	case map[string]any:
		for _, b512HeatpumpIOEvaluatorChild := range b512HeatpumpIOEvaluatorItem {
			if b512HeatpumpIOEvaluatorErr := b512HeatpumpIOEvaluatorValidateArguments(b512HeatpumpIOEvaluatorChild); b512HeatpumpIOEvaluatorErr != nil {
				return b512HeatpumpIOEvaluatorErr
			}
		}
	case []any:
		for _, b512HeatpumpIOEvaluatorChild := range b512HeatpumpIOEvaluatorItem {
			if b512HeatpumpIOEvaluatorErr := b512HeatpumpIOEvaluatorValidateArguments(b512HeatpumpIOEvaluatorChild); b512HeatpumpIOEvaluatorErr != nil {
				return b512HeatpumpIOEvaluatorErr
			}
		}
	}
	return nil
}

func B512HeatpumpIOEvaluate(b512HeatpumpIOEvaluatorInput json.RawMessage) (json.RawMessage, error) {
	b512HeatpumpIOEvaluatorDecoder := json.NewDecoder(bytes.NewReader(b512HeatpumpIOEvaluatorInput))
	b512HeatpumpIOEvaluatorDecoder.UseNumber()
	b512HeatpumpIOEvaluatorValue, b512HeatpumpIOEvaluatorErr := b512HeatpumpIOEvaluatorRead(b512HeatpumpIOEvaluatorDecoder)
	if b512HeatpumpIOEvaluatorErr != nil {
		return nil, b512HeatpumpIOEvaluatorErr
	}
	if _, b512HeatpumpIOEvaluatorErr = b512HeatpumpIOEvaluatorDecoder.Token(); b512HeatpumpIOEvaluatorErr != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b512HeatpumpIOEvaluatorRoot, b512HeatpumpIOEvaluatorOK := b512HeatpumpIOEvaluatorValue.(map[string]any)
	if !b512HeatpumpIOEvaluatorOK {
		return nil, fmt.Errorf("input object required")
	}
	b512HeatpumpIOEvaluatorContext, b512HeatpumpIOEvaluatorOK := b512HeatpumpIOEvaluatorRoot["context"].(map[string]any)
	if !b512HeatpumpIOEvaluatorOK {
		return nil, fmt.Errorf("context object required")
	}
	b512HeatpumpIOEvaluatorRaw, b512HeatpumpIOEvaluatorOK := b512HeatpumpIOEvaluatorRoot["fields"].(map[string]any)
	if !b512HeatpumpIOEvaluatorOK {
		return nil, fmt.Errorf("fields object required")
	}
	for _, b512HeatpumpIOEvaluatorKey := range []string{"primary", "secondary", "operation"} {
		if b512HeatpumpIOEvaluatorItem, b512HeatpumpIOEvaluatorExists := b512HeatpumpIOEvaluatorContext[b512HeatpumpIOEvaluatorKey]; b512HeatpumpIOEvaluatorExists && b512HeatpumpIOEvaluatorItem != nil {
			if _, b512HeatpumpIOEvaluatorErr := b512HeatpumpIOEvaluatorUint(b512HeatpumpIOEvaluatorItem, 8); b512HeatpumpIOEvaluatorErr != nil {
				return nil, fmt.Errorf("context %s: %w", b512HeatpumpIOEvaluatorKey, b512HeatpumpIOEvaluatorErr)
			}
		}
	}
	for _, b512HeatpumpIOEvaluatorKey := range []string{"view", "direction"} {
		if b512HeatpumpIOEvaluatorItem, b512HeatpumpIOEvaluatorExists := b512HeatpumpIOEvaluatorContext[b512HeatpumpIOEvaluatorKey]; b512HeatpumpIOEvaluatorExists && b512HeatpumpIOEvaluatorItem != nil {
			if _, b512HeatpumpIOEvaluatorOK := b512HeatpumpIOEvaluatorItem.(string); !b512HeatpumpIOEvaluatorOK {
				return nil, fmt.Errorf("context %s: string required", b512HeatpumpIOEvaluatorKey)
			}
		}
	}
	for b512HeatpumpIOEvaluatorKey, b512HeatpumpIOEvaluatorItem := range b512HeatpumpIOEvaluatorRaw {
		b512HeatpumpIOEvaluatorArray, b512HeatpumpIOEvaluatorOK := b512HeatpumpIOEvaluatorItem.([]any)
		if !b512HeatpumpIOEvaluatorOK {
			return nil, fmt.Errorf("field %s: byte array required", b512HeatpumpIOEvaluatorKey)
		}
		for _, b512HeatpumpIOEvaluatorByte := range b512HeatpumpIOEvaluatorArray {
			if _, b512HeatpumpIOEvaluatorErr := b512HeatpumpIOEvaluatorUint(b512HeatpumpIOEvaluatorByte, 8); b512HeatpumpIOEvaluatorErr != nil {
				return nil, fmt.Errorf("field %s: %w", b512HeatpumpIOEvaluatorKey, b512HeatpumpIOEvaluatorErr)
			}
		}
	}
	b512HeatpumpIOEvaluatorFields := map[string]any{}
	b512HeatpumpIOEvaluatorOutput := map[string]any{"context": b512HeatpumpIOEvaluatorContext, "raw": b512HeatpumpIOEvaluatorRaw, "state": "Unknown", "fields": b512HeatpumpIOEvaluatorFields, "physical": "Unknown", "qualified": false}
	if b512HeatpumpIOEvaluatorArguments, b512HeatpumpIOEvaluatorExists := b512HeatpumpIOEvaluatorRoot["arguments"]; b512HeatpumpIOEvaluatorExists {
		if _, b512HeatpumpIOEvaluatorOK := b512HeatpumpIOEvaluatorArguments.(map[string]any); !b512HeatpumpIOEvaluatorOK {
			return nil, fmt.Errorf("arguments object required")
		}
		if b512HeatpumpIOEvaluatorErr := b512HeatpumpIOEvaluatorValidateArguments(b512HeatpumpIOEvaluatorArguments); b512HeatpumpIOEvaluatorErr != nil {
			return nil, b512HeatpumpIOEvaluatorErr
		}
		b512HeatpumpIOEvaluatorOutput["arguments"] = b512HeatpumpIOEvaluatorArguments
	}
	b512HeatpumpIOEvaluatorPrimary, _ := b512HeatpumpIOEvaluatorUint(b512HeatpumpIOEvaluatorContext["primary"], 8)
	b512HeatpumpIOEvaluatorSecondary, _ := b512HeatpumpIOEvaluatorUint(b512HeatpumpIOEvaluatorContext["secondary"], 8)
	b512HeatpumpIOEvaluatorOperation, _ := b512HeatpumpIOEvaluatorUint(b512HeatpumpIOEvaluatorContext["operation"], 8)
	b512HeatpumpIOEvaluatorDirection, _ := b512HeatpumpIOEvaluatorContext["direction"].(string)
	b512HeatpumpIOEvaluatorView, _ := b512HeatpumpIOEvaluatorContext["view"].(string)
	b512HeatpumpIOEvaluatorSupported := b512HeatpumpIOEvaluatorPrimary == 181 && b512HeatpumpIOEvaluatorSecondary == 18 && (b512HeatpumpIOEvaluatorDirection == "request" || b512HeatpumpIOEvaluatorDirection == "positive") && ((b512HeatpumpIOEvaluatorOperation == 15 && b512HeatpumpIOEvaluatorView == "valve-io") || (b512HeatpumpIOEvaluatorOperation == 16 && b512HeatpumpIOEvaluatorView == "board-io"))
	if !b512HeatpumpIOEvaluatorSupported {
		return json.Marshal(b512HeatpumpIOEvaluatorOutput)
	}
	b512HeatpumpIOEvaluatorOutput["state"] = "complete"
	b512HeatpumpIOEvaluatorNames := []string{}
	if b512HeatpumpIOEvaluatorOperation == 15 && b512HeatpumpIOEvaluatorDirection == "request" {
		b512HeatpumpIOEvaluatorNames = []string{"output", "valve"}
	}
	if b512HeatpumpIOEvaluatorOperation == 15 && b512HeatpumpIOEvaluatorDirection == "positive" {
		b512HeatpumpIOEvaluatorNames = []string{"temperature1", "temperature2", "input", "pressure", "backup"}
	}
	if b512HeatpumpIOEvaluatorOperation == 16 && b512HeatpumpIOEvaluatorDirection == "request" {
		b512HeatpumpIOEvaluatorNames = []string{"actuator1", "actuator2", "inverse"}
	}
	if b512HeatpumpIOEvaluatorOperation == 16 && b512HeatpumpIOEvaluatorDirection == "positive" {
		b512HeatpumpIOEvaluatorNames = []string{"status", "phase1", "phase2", "adc0", "adc1", "adc2", "adc3", "adc4", "adc5", "board"}
	}
	for _, b512HeatpumpIOEvaluatorName := range b512HeatpumpIOEvaluatorNames {
		b512HeatpumpIOEvaluatorResult := map[string]any{"state": "incomplete"}
		b512HeatpumpIOEvaluatorFields[b512HeatpumpIOEvaluatorName] = b512HeatpumpIOEvaluatorResult
		b512HeatpumpIOEvaluatorItem, b512HeatpumpIOEvaluatorExists := b512HeatpumpIOEvaluatorRaw[b512HeatpumpIOEvaluatorName]
		if b512HeatpumpIOEvaluatorExists {
			b512HeatpumpIOEvaluatorResult["raw"] = b512HeatpumpIOEvaluatorItem
		}
		b512HeatpumpIOEvaluatorArray, _ := b512HeatpumpIOEvaluatorItem.([]any)
		b512HeatpumpIOEvaluatorWidth := 1
		if b512HeatpumpIOEvaluatorName == "temperature1" || b512HeatpumpIOEvaluatorName == "temperature2" || b512HeatpumpIOEvaluatorName == "inverse" || (len(b512HeatpumpIOEvaluatorName) == 4 && b512HeatpumpIOEvaluatorName[:3] == "adc") {
			b512HeatpumpIOEvaluatorWidth = 2
		}
		if !b512HeatpumpIOEvaluatorExists || len(b512HeatpumpIOEvaluatorArray) != b512HeatpumpIOEvaluatorWidth {
			b512HeatpumpIOEvaluatorOutput["state"] = "partial"
			continue
		}
		b512HeatpumpIOEvaluatorFirst, _ := b512HeatpumpIOEvaluatorUint(b512HeatpumpIOEvaluatorArray[0], 8)
		b512HeatpumpIOEvaluatorWord := b512HeatpumpIOEvaluatorFirst
		if b512HeatpumpIOEvaluatorWidth == 2 {
			b512HeatpumpIOEvaluatorSecond, _ := b512HeatpumpIOEvaluatorUint(b512HeatpumpIOEvaluatorArray[1], 8)
			b512HeatpumpIOEvaluatorWord |= b512HeatpumpIOEvaluatorSecond << 8
		}
		b512HeatpumpIOEvaluatorResult["state"] = "Hypothesis"
		switch b512HeatpumpIOEvaluatorName {
		case "temperature1", "temperature2":
			if b512HeatpumpIOEvaluatorWord == 0x8000 {
				b512HeatpumpIOEvaluatorResult["label"] = "replacement"
			} else {
				b512HeatpumpIOEvaluatorResult["number"] = float64(int16(b512HeatpumpIOEvaluatorWord)) / 16
				b512HeatpumpIOEvaluatorResult["calculation"] = "Proven"
				b512HeatpumpIOEvaluatorResult["unit"] = "Celsius"
			}
		case "pressure":
			b512HeatpumpIOEvaluatorResult["number"] = float64(b512HeatpumpIOEvaluatorFirst) / 10
			b512HeatpumpIOEvaluatorResult["calculation"] = "Proven"
			b512HeatpumpIOEvaluatorResult["unit"] = "bar"
			b512HeatpumpIOEvaluatorResult["candidate_domain"] = true
			b512HeatpumpIOEvaluatorResult["alternatives"] = []any{map[string]any{"state": "Hypothesis", "label": "candidate-default", "number": 255.0, "candidate_domain": false, "conflict": true}}
		case "phase1", "phase2":
			b512HeatpumpIOEvaluatorResult["number"] = b512HeatpumpIOEvaluatorFirst
			b512HeatpumpIOEvaluatorResult["unit"] = "percent"
			b512HeatpumpIOEvaluatorResult["alternatives"] = []any{map[string]any{"state": "Hypothesis", "label": "single-phase", "approximate": []int{50}}, map[string]any{"state": "Hypothesis", "label": "three-phase", "approximate": []int{0, 33, 66, 100}}}
		case "adc0", "adc1", "adc2", "adc3", "adc4", "adc5":
			b512HeatpumpIOEvaluatorResult["number"] = b512HeatpumpIOEvaluatorWord
			b512HeatpumpIOEvaluatorResult["calculation"] = "Proven"
		case "inverse":
			b512HeatpumpIOEvaluatorResult["state"] = "Unknown"
			b512HeatpumpIOEvaluatorResult["unmodeled"] = map[string]any{"order": "Unknown", "algorithm": "Unknown"}
		case "valve", "input", "backup", "board":
			b512HeatpumpIOEvaluatorLabel := ""
			switch b512HeatpumpIOEvaluatorName {
			case "valve":
				if b512HeatpumpIOEvaluatorFirst < 3 {
					b512HeatpumpIOEvaluatorLabel = []string{"off", "CH", "DHW"}[b512HeatpumpIOEvaluatorFirst]
				}
			case "input":
				if b512HeatpumpIOEvaluatorFirst < 2 {
					b512HeatpumpIOEvaluatorLabel = []string{"multipurpose", "utility-contact"}[b512HeatpumpIOEvaluatorFirst]
				}
			case "backup":
				if b512HeatpumpIOEvaluatorFirst == 0 {
					b512HeatpumpIOEvaluatorLabel = "off"
				}
				if b512HeatpumpIOEvaluatorFirst >= 1 && b512HeatpumpIOEvaluatorFirst <= 3 {
					b512HeatpumpIOEvaluatorLabel = "step"
					b512HeatpumpIOEvaluatorResult["number"] = b512HeatpumpIOEvaluatorFirst
				}
				if b512HeatpumpIOEvaluatorFirst == 255 {
					b512HeatpumpIOEvaluatorLabel = "replacement"
				}
			case "board":
				if b512HeatpumpIOEvaluatorFirst <= 3 || b512HeatpumpIOEvaluatorFirst == 11 || b512HeatpumpIOEvaluatorFirst == 12 {
					b512HeatpumpIOEvaluatorLabel = "opaque-code-" + strconv.FormatUint(b512HeatpumpIOEvaluatorFirst, 10)
				}
			}
			if b512HeatpumpIOEvaluatorLabel == "" {
				b512HeatpumpIOEvaluatorResult["state"] = "Unknown"
				b512HeatpumpIOEvaluatorResult["unmodeled"] = true
			} else {
				b512HeatpumpIOEvaluatorResult["label"] = b512HeatpumpIOEvaluatorLabel
			}
		default:
			b512HeatpumpIOEvaluatorLabels := []string{}
			switch b512HeatpumpIOEvaluatorName {
			case "output":
				b512HeatpumpIOEvaluatorLabels = []string{"output1", "output2"}
			case "actuator1":
				b512HeatpumpIOEvaluatorLabels = []string{"cooling-valve", "passive-mixer1", "passive-mixer2", "cooling-active", "well-pump", "backup-safety", "backup1", "backup2"}
			case "actuator2":
				b512HeatpumpIOEvaluatorLabels = []string{"backup3", "external-valve", "external-backup"}
			case "status":
				b512HeatpumpIOEvaluatorLabels = []string{"anode1", "anode2", "smart-grid", "well", "safety-limiter", "safety"}
			}
			b512HeatpumpIOEvaluatorFlags := map[string]bool{}
			for b512HeatpumpIOEvaluatorBit, b512HeatpumpIOEvaluatorLabel := range b512HeatpumpIOEvaluatorLabels {
				b512HeatpumpIOEvaluatorFlags[b512HeatpumpIOEvaluatorLabel] = b512HeatpumpIOEvaluatorFirst&(1<<b512HeatpumpIOEvaluatorBit) != 0
			}
			b512HeatpumpIOEvaluatorResult["flags"] = b512HeatpumpIOEvaluatorFlags
			b512HeatpumpIOEvaluatorResult["reserved"] = b512HeatpumpIOEvaluatorFirst & (255 << len(b512HeatpumpIOEvaluatorLabels))
		}
	}
	return json.Marshal(b512HeatpumpIOEvaluatorOutput)
}
