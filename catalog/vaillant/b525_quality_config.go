package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
	"unicode/utf8"
)

// b525QualityConfigMaxDepth bounds container nesting so hostile input cannot exhaust the stack.
const b525QualityConfigMaxDepth = 256

func b525QualityConfigRead(b525QualityConfigDecoder *json.Decoder, b525QualityConfigDepth int) (any, error) {
	if b525QualityConfigDepth > b525QualityConfigMaxDepth {
		return nil, fmt.Errorf("nesting exceeds maximum depth")
	}
	b525QualityConfigToken, b525QualityConfigError := b525QualityConfigDecoder.Token()
	if b525QualityConfigError != nil {
		return nil, b525QualityConfigError
	}
	if b525QualityConfigDelimiter, b525QualityConfigOK := b525QualityConfigToken.(json.Delim); b525QualityConfigOK {
		switch b525QualityConfigDelimiter {
		case '{':
			b525QualityConfigObject := map[string]any{}
			for b525QualityConfigDecoder.More() {
				b525QualityConfigKeyToken, b525QualityConfigError := b525QualityConfigDecoder.Token()
				if b525QualityConfigError != nil {
					return nil, b525QualityConfigError
				}
				b525QualityConfigKey := b525QualityConfigKeyToken.(string)
				if _, b525QualityConfigExists := b525QualityConfigObject[b525QualityConfigKey]; b525QualityConfigExists {
					return nil, fmt.Errorf("duplicate key %q", b525QualityConfigKey)
				}
				b525QualityConfigValue, b525QualityConfigError := b525QualityConfigRead(b525QualityConfigDecoder, b525QualityConfigDepth+1)
				if b525QualityConfigError != nil {
					return nil, b525QualityConfigError
				}
				b525QualityConfigObject[b525QualityConfigKey] = b525QualityConfigValue
			}
			if _, b525QualityConfigError := b525QualityConfigDecoder.Token(); b525QualityConfigError != nil {
				return nil, b525QualityConfigError
			}
			return b525QualityConfigObject, nil
		case '[':
			b525QualityConfigArray := []any{}
			for b525QualityConfigDecoder.More() {
				b525QualityConfigValue, b525QualityConfigError := b525QualityConfigRead(b525QualityConfigDecoder, b525QualityConfigDepth+1)
				if b525QualityConfigError != nil {
					return nil, b525QualityConfigError
				}
				b525QualityConfigArray = append(b525QualityConfigArray, b525QualityConfigValue)
			}
			if _, b525QualityConfigError := b525QualityConfigDecoder.Token(); b525QualityConfigError != nil {
				return nil, b525QualityConfigError
			}
			return b525QualityConfigArray, nil
		}
		return nil, fmt.Errorf("unexpected delimiter")
	}
	return b525QualityConfigToken, nil
}

func b525QualityConfigByte(b525QualityConfigValue any) (int, bool) {
	b525QualityConfigNumber, b525QualityConfigOK := b525QualityConfigValue.(json.Number)
	if !b525QualityConfigOK {
		return 0, false
	}
	// Bound the token before big.Rat materializes any power of ten.
	if len(b525QualityConfigNumber) > 64 {
		return 0, false
	}
	if b525QualityConfigMarker := strings.IndexAny(string(b525QualityConfigNumber), "eE"); b525QualityConfigMarker >= 0 {
		b525QualityConfigExponent, b525QualityConfigError := strconv.Atoi(string(b525QualityConfigNumber[b525QualityConfigMarker+1:]))
		if b525QualityConfigError != nil || b525QualityConfigExponent < -64 || b525QualityConfigExponent > 64 {
			return 0, false
		}
	}
	b525QualityConfigRational, b525QualityConfigOK := new(big.Rat).SetString(string(b525QualityConfigNumber))
	if !b525QualityConfigOK || !b525QualityConfigRational.IsInt() || !b525QualityConfigRational.Num().IsInt64() {
		return 0, false
	}
	b525QualityConfigInteger := b525QualityConfigRational.Num().Int64()
	if b525QualityConfigInteger < 0 || b525QualityConfigInteger > 255 {
		return 0, false
	}
	return int(b525QualityConfigInteger), true
}

// b525QualityConfigValidateEscapes rejects surrogate escapes that encoding/json would
// silently replace, while leaving escaped backslashes and valid pairs intact.
func b525QualityConfigValidateEscapes(b525QualityConfigInput []byte) error {
	for b525QualityConfigIndex := 0; b525QualityConfigIndex < len(b525QualityConfigInput); b525QualityConfigIndex++ {
		if b525QualityConfigInput[b525QualityConfigIndex] != '\\' {
			continue
		}
		b525QualityConfigIndex++
		if b525QualityConfigIndex >= len(b525QualityConfigInput) || b525QualityConfigInput[b525QualityConfigIndex] != 'u' {
			continue
		}
		if len(b525QualityConfigInput)-b525QualityConfigIndex < 5 {
			return fmt.Errorf("invalid Unicode escape")
		}
		b525QualityConfigCode, b525QualityConfigError := strconv.ParseUint(string(b525QualityConfigInput[b525QualityConfigIndex+1:b525QualityConfigIndex+5]), 16, 16)
		if b525QualityConfigError != nil {
			return fmt.Errorf("invalid Unicode escape")
		}
		b525QualityConfigIndex += 4
		if b525QualityConfigCode >= 0xdc00 && b525QualityConfigCode <= 0xdfff {
			return fmt.Errorf("unpaired low surrogate escape")
		}
		if b525QualityConfigCode >= 0xd800 && b525QualityConfigCode <= 0xdbff {
			if len(b525QualityConfigInput)-b525QualityConfigIndex < 7 || b525QualityConfigInput[b525QualityConfigIndex+1] != '\\' || b525QualityConfigInput[b525QualityConfigIndex+2] != 'u' {
				return fmt.Errorf("unpaired high surrogate escape")
			}
			b525QualityConfigLow, b525QualityConfigError := strconv.ParseUint(string(b525QualityConfigInput[b525QualityConfigIndex+3:b525QualityConfigIndex+7]), 16, 16)
			if b525QualityConfigError != nil || b525QualityConfigLow < 0xdc00 || b525QualityConfigLow > 0xdfff {
				return fmt.Errorf("unpaired high surrogate escape")
			}
			b525QualityConfigIndex += 6
		}
	}
	return nil
}

func B525QualityConfigEvaluate(b525QualityConfigInput json.RawMessage) (json.RawMessage, error) {
	if !utf8.Valid(b525QualityConfigInput) {
		return nil, fmt.Errorf("input must be valid UTF-8")
	}
	if b525QualityConfigError := b525QualityConfigValidateEscapes(b525QualityConfigInput); b525QualityConfigError != nil {
		return nil, b525QualityConfigError
	}
	b525QualityConfigDecoder := json.NewDecoder(bytes.NewReader(b525QualityConfigInput))
	b525QualityConfigDecoder.UseNumber()
	b525QualityConfigValue, b525QualityConfigError := b525QualityConfigRead(b525QualityConfigDecoder, 0)
	if b525QualityConfigError != nil {
		return nil, b525QualityConfigError
	}
	if _, b525QualityConfigError = b525QualityConfigDecoder.Token(); b525QualityConfigError != io.EOF {
		return nil, fmt.Errorf("trailing JSON")
	}
	b525QualityConfigInputObject, b525QualityConfigOK := b525QualityConfigValue.(map[string]any)
	if !b525QualityConfigOK {
		return nil, fmt.Errorf("input must be object")
	}
	b525QualityConfigContext, b525QualityConfigOK := b525QualityConfigInputObject["context"].(map[string]any)
	if !b525QualityConfigOK {
		return nil, fmt.Errorf("context must be object")
	}
	b525QualityConfigRaw, b525QualityConfigOK := b525QualityConfigInputObject["fields"].(map[string]any)
	if !b525QualityConfigOK {
		return nil, fmt.Errorf("fields must be object")
	}
	for _, b525QualityConfigName := range []string{"primary", "secondary", "operation"} {
		if b525QualityConfigValue, b525QualityConfigExists := b525QualityConfigContext[b525QualityConfigName]; b525QualityConfigExists && b525QualityConfigValue != nil {
			if _, b525QualityConfigOK := b525QualityConfigByte(b525QualityConfigValue); !b525QualityConfigOK {
				return nil, fmt.Errorf("invalid context %s", b525QualityConfigName)
			}
		}
	}
	for _, b525QualityConfigName := range []string{"direction", "view"} {
		if b525QualityConfigValue, b525QualityConfigExists := b525QualityConfigContext[b525QualityConfigName]; b525QualityConfigExists && b525QualityConfigValue != nil {
			if _, b525QualityConfigOK := b525QualityConfigValue.(string); !b525QualityConfigOK {
				return nil, fmt.Errorf("invalid context %s", b525QualityConfigName)
			}
		}
	}
	for b525QualityConfigName, b525QualityConfigValue := range b525QualityConfigRaw {
		b525QualityConfigArray, b525QualityConfigOK := b525QualityConfigValue.([]any)
		if !b525QualityConfigOK {
			return nil, fmt.Errorf("field %s must be byte array", b525QualityConfigName)
		}
		for _, b525QualityConfigItem := range b525QualityConfigArray {
			if _, b525QualityConfigOK := b525QualityConfigByte(b525QualityConfigItem); !b525QualityConfigOK {
				return nil, fmt.Errorf("invalid byte in %s", b525QualityConfigName)
			}
		}
	}
	b525QualityConfigFields := map[string]any{}
	b525QualityConfigOutput := map[string]any{"context": b525QualityConfigContext, "raw": b525QualityConfigRaw, "state": "Unknown", "fields": b525QualityConfigFields, "physical": "Unknown", "qualified": false}
	if b525QualityConfigArguments, b525QualityConfigExists := b525QualityConfigInputObject["arguments"]; b525QualityConfigExists {
		if _, b525QualityConfigOK := b525QualityConfigArguments.(map[string]any); !b525QualityConfigOK {
			return nil, fmt.Errorf("arguments must be object")
		}
		b525QualityConfigOutput["arguments"] = b525QualityConfigArguments
	}
	b525QualityConfigPrimary, _ := b525QualityConfigByte(b525QualityConfigContext["primary"])
	b525QualityConfigSecondary, _ := b525QualityConfigByte(b525QualityConfigContext["secondary"])
	b525QualityConfigOperation, _ := b525QualityConfigByte(b525QualityConfigContext["operation"])
	b525QualityConfigQuality := b525QualityConfigOperation == 3 && b525QualityConfigContext["view"] == "quality" && b525QualityConfigContext["direction"] == "positive"
	b525QualityConfigConfiguration := b525QualityConfigOperation >= 24 && b525QualityConfigOperation <= 27 && b525QualityConfigContext["view"] == "configuration" && (b525QualityConfigContext["direction"] == "request" || b525QualityConfigContext["direction"] == "positive")
	if b525QualityConfigPrimary != 181 || b525QualityConfigSecondary != 37 || (!b525QualityConfigQuality && !b525QualityConfigConfiguration) {
		return json.Marshal(b525QualityConfigOutput)
	}
	b525QualityConfigOutput["state"] = "complete"
	b525QualityConfigNames := []string{"rssi", "quality"}
	if b525QualityConfigConfiguration {
		b525QualityConfigNames = []string{"configuration", "data"}
	}
	for _, b525QualityConfigName := range b525QualityConfigNames {
		b525QualityConfigResult := map[string]any{"state": "Unknown"}
		b525QualityConfigFields[b525QualityConfigName] = b525QualityConfigResult
		b525QualityConfigArray, b525QualityConfigExists := b525QualityConfigRaw[b525QualityConfigName].([]any)
		if b525QualityConfigExists {
			b525QualityConfigResult["raw"] = b525QualityConfigArray
		}
		if !b525QualityConfigExists || (b525QualityConfigName != "data" && len(b525QualityConfigArray) != 1) {
			b525QualityConfigOutput["state"] = "partial"
			continue
		}
		b525QualityConfigResult["state"] = "Hypothesis"
		if b525QualityConfigName == "data" {
			continue
		}
		b525QualityConfigNumber, _ := b525QualityConfigByte(b525QualityConfigArray[0])
		switch b525QualityConfigName {
		case "rssi":
			if b525QualityConfigNumber >= 128 {
				b525QualityConfigNumber -= 256
			}
			b525QualityConfigResult["number"] = b525QualityConfigNumber
			b525QualityConfigResult["unit"] = "Unknown"
			b525QualityConfigResult["calculation"] = "Proven"
		case "quality":
			b525QualityConfigFlags := map[string]bool{}
			for b525QualityConfigBit := 0; b525QualityConfigBit < 7; b525QualityConfigBit++ {
				b525QualityConfigFlags[fmt.Sprintf("bit%d", b525QualityConfigBit)] = b525QualityConfigNumber&(1<<b525QualityConfigBit) != 0
			}
			b525QualityConfigFlags["crc_okay"] = b525QualityConfigNumber&128 != 0
			b525QualityConfigResult["flags"] = b525QualityConfigFlags
			b525QualityConfigResult["checksum_verified"] = false
		case "configuration":
			b525QualityConfigResult["flags"] = map[string]bool{"bit0": b525QualityConfigNumber&1 != 0}
			b525QualityConfigResult["unmodeled"] = b525QualityConfigNumber & 254
			b525QualityConfigResult["candidate_default"] = map[string]any{"number": 0, "label": "read", "state": "Hypothesis"}
		}
	}
	return json.Marshal(b525QualityConfigOutput)
}
