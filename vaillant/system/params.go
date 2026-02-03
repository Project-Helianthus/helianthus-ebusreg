package system

func uint8Param(params map[string]any, key string) (byte, bool) {
	value, ok := params[key]
	if !ok {
		return 0, false
	}
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
