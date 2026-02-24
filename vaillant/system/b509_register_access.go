package system

import (
	"fmt"
	"strconv"
	"strings"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
	"github.com/d3vi1/helianthus-ebusgo/types"
)

const (
	methodGetRegister = "get_register"
	methodSetRegister = "set_register"
)

const (
	registerCmdRead  = byte(0x0D)
	registerCmdWrite = byte(0x0E)
)

type registerReadTemplate struct {
	primary   byte
	secondary byte
}

func (template registerReadTemplate) Primary() byte {
	return template.primary
}

func (template registerReadTemplate) Secondary() byte {
	return template.secondary
}

func (template registerReadTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("register read template missing params: %w", ebuserrors.ErrInvalidPayload)
	}

	addr, ok, err := registerAddrParam(params, "addr")
	if err != nil {
		return nil, fmt.Errorf("register read template addr: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("register read template addr: %w", ebuserrors.ErrInvalidPayload)
	}

	return []byte{registerCmdRead, byte(addr >> 8), byte(addr)}, nil
}

type registerWriteTemplate struct {
	primary   byte
	secondary byte
}

func (template registerWriteTemplate) Primary() byte {
	return template.primary
}

func (template registerWriteTemplate) Secondary() byte {
	return template.secondary
}

func (template registerWriteTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("register write template missing params: %w", ebuserrors.ErrInvalidPayload)
	}

	addr, ok, err := registerAddrParam(params, "addr")
	if err != nil {
		return nil, fmt.Errorf("register write template addr: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("register write template addr: %w", ebuserrors.ErrInvalidPayload)
	}

	data, hasData, err := bytesParam(params, "data")
	if err != nil {
		return nil, fmt.Errorf("register write template data: %w", err)
	}
	if !hasData {
		data, _, err = bytesParam(params, "payload")
		if err != nil {
			return nil, fmt.Errorf("register write template payload: %w", err)
		}
	}

	if len(data) > 0xFC {
		return nil, fmt.Errorf("register write template data too long: %w", ebuserrors.ErrInvalidPayload)
	}

	payload := make([]byte, 0, 3+len(data))
	payload = append(payload, registerCmdWrite, byte(addr>>8), byte(addr))
	payload = append(payload, data...)
	return payload, nil
}

func decodeRegisterResponse(cmd byte, addr uint16, payload []byte) map[string]types.Value {
	values := map[string]types.Value{
		"cmd":      {Value: cmd, Valid: true},
		"addr":     {Value: addr, Valid: true},
		"addr_hex": {Value: fmt.Sprintf("%04X", addr), Valid: true},
		"payload":  {Value: append([]byte(nil), payload...), Valid: true},
	}

	if len(payload) == 1 && payload[0] == 0x00 {
		values["value"] = types.Value{Valid: false}
		return values
	}

	values["value"] = types.Value{Value: append([]byte(nil), payload...), Valid: true}
	return values
}

func registerAddrParam(params map[string]any, key string) (uint16, bool, error) {
	value, ok := params[key]
	if !ok {
		return 0, false, nil
	}
	if value == nil {
		return 0, true, ebuserrors.ErrInvalidPayload
	}
	if parsed, ok := parseUint(value, 0xFFFF); ok {
		return uint16(parsed), true, nil
	}

	switch typed := value.(type) {
	case string:
		s := strings.TrimSpace(typed)
		s = strings.TrimPrefix(s, "0x")
		s = strings.TrimPrefix(s, "0X")
		if len(s) != 4 {
			return 0, true, ebuserrors.ErrInvalidPayload
		}
		parsed, err := strconv.ParseUint(s, 16, 16)
		if err != nil {
			return 0, true, ebuserrors.ErrInvalidPayload
		}
		return uint16(parsed), true, nil
	case []byte:
		if len(typed) != 2 {
			return 0, true, ebuserrors.ErrInvalidPayload
		}
		return uint16(typed[0])<<8 | uint16(typed[1]), true, nil
	case []int:
		if len(typed) != 2 {
			return 0, true, ebuserrors.ErrInvalidPayload
		}
		for _, v := range typed {
			if v < 0 || v > 255 {
				return 0, true, ebuserrors.ErrInvalidPayload
			}
		}
		return uint16(byte(typed[0]))<<8 | uint16(byte(typed[1])), true, nil
	case []any:
		if len(typed) != 2 {
			return 0, true, ebuserrors.ErrInvalidPayload
		}
		hi, ok := toByte(typed[0])
		if !ok {
			return 0, true, ebuserrors.ErrInvalidPayload
		}
		lo, ok := toByte(typed[1])
		if !ok {
			return 0, true, ebuserrors.ErrInvalidPayload
		}
		return uint16(hi)<<8 | uint16(lo), true, nil
	default:
		return 0, true, ebuserrors.ErrInvalidPayload
	}
}
