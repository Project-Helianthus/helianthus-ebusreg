package system

import (
	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
)

func uint8Param(params map[string]any, key string) (byte, bool) {
	value, ok := params[key]
	if !ok {
		return 0, false
	}
	return toByte(value)
}

func toByte(value any) (byte, bool) {
	switch typed := value.(type) {
	case int:
		if typed < 0 || typed > 255 {
			return 0, false
		}
		return byte(typed), true
	case int8:
		if typed < 0 {
			return 0, false
		}
		return byte(typed), true
	case int16:
		if typed < 0 || typed > 255 {
			return 0, false
		}
		return byte(typed), true
	case int32:
		if typed < 0 || typed > 255 {
			return 0, false
		}
		return byte(typed), true
	case int64:
		if typed < 0 || typed > 255 {
			return 0, false
		}
		return byte(typed), true
	case uint:
		if typed > 255 {
			return 0, false
		}
		return byte(typed), true
	case uint8:
		return typed, true
	case uint16:
		if typed > 255 {
			return 0, false
		}
		return byte(typed), true
	case uint32:
		if typed > 255 {
			return 0, false
		}
		return byte(typed), true
	case uint64:
		if typed > 255 {
			return 0, false
		}
		return byte(typed), true
	default:
		return 0, false
	}
}

func bytesParam(params map[string]any, key string) ([]byte, bool, error) {
	value, ok := params[key]
	if !ok {
		return nil, false, nil
	}
	if value == nil {
		return nil, true, nil
	}

	switch typed := value.(type) {
	case []byte:
		return append([]byte(nil), typed...), true, nil
	case []int:
		out := make([]byte, len(typed))
		for i, v := range typed {
			if v < 0 || v > 255 {
				return nil, true, ebuserrors.ErrInvalidPayload
			}
			out[i] = byte(v)
		}
		return out, true, nil
	case []any:
		out := make([]byte, len(typed))
		for i, v := range typed {
			b, ok := toByte(v)
			if !ok {
				return nil, true, ebuserrors.ErrInvalidPayload
			}
			out[i] = b
		}
		return out, true, nil
	default:
		if b, ok := toByte(value); ok {
			return []byte{b}, true, nil
		}
		return nil, true, ebuserrors.ErrInvalidPayload
	}
}
