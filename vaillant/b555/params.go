package b555

import (
	"encoding/json"
	"math"
)

func uint8Param(params map[string]any, key string) (byte, bool) {
	value, ok := params[key]
	if !ok {
		return 0, false
	}
	return toByte(value)
}

func toByte(value any) (byte, bool) {
	parsed, ok := parseUint(value, 0xFF)
	if !ok {
		return 0, false
	}
	return byte(parsed), true
}

func uint16Param(params map[string]any, key string) (uint16, bool) {
	value, ok := params[key]
	if !ok {
		return 0, false
	}
	parsed, ok := parseUint(value, 0xFFFF)
	if !ok {
		return 0, false
	}
	return uint16(parsed), true
}

func parseUint(value any, max uint64) (uint64, bool) {
	switch typed := value.(type) {
	case int:
		if typed < 0 || uint64(typed) > max {
			return 0, false
		}
		return uint64(typed), true
	case int8:
		if typed < 0 || uint64(typed) > max {
			return 0, false
		}
		return uint64(typed), true
	case int16:
		if typed < 0 || uint64(typed) > max {
			return 0, false
		}
		return uint64(typed), true
	case int32:
		if typed < 0 || uint64(typed) > max {
			return 0, false
		}
		return uint64(typed), true
	case int64:
		if typed < 0 || uint64(typed) > max {
			return 0, false
		}
		return uint64(typed), true
	case uint:
		if uint64(typed) > max {
			return 0, false
		}
		return uint64(typed), true
	case uint8:
		if uint64(typed) > max {
			return 0, false
		}
		return uint64(typed), true
	case uint16:
		if uint64(typed) > max {
			return 0, false
		}
		return uint64(typed), true
	case uint32:
		if uint64(typed) > max {
			return 0, false
		}
		return uint64(typed), true
	case uint64:
		if typed > max {
			return 0, false
		}
		return typed, true
	case float32:
		return parseFloatUint(float64(typed), max)
	case float64:
		return parseFloatUint(typed, max)
	case json.Number:
		if parsed, err := typed.Int64(); err == nil {
			if parsed < 0 || uint64(parsed) > max {
				return 0, false
			}
			return uint64(parsed), true
		}
		floatParsed, err := typed.Float64()
		if err != nil {
			return 0, false
		}
		return parseFloatUint(floatParsed, max)
	default:
		return 0, false
	}
}

func parseFloatUint(value float64, max uint64) (uint64, bool) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, false
	}
	if value < 0 || value > float64(max) {
		return 0, false
	}
	if math.Trunc(value) != value {
		return 0, false
	}
	return uint64(value), true
}
