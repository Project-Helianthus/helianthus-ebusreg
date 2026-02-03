package system

import (
	"fmt"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
	"github.com/d3vi1/helianthus-ebusgo/types"
)

const (
	methodGetExtRegister = "get_ext_register"
	methodSetExtRegister = "set_ext_register"
)

const (
	extRegisterCmdRead  = byte(0x00)
	extRegisterCmdWrite = byte(0x01)
	extRegisterPrefix   = byte(0x02)
)

type extRegisterReadTemplate struct {
	primary   byte
	secondary byte
}

func (template extRegisterReadTemplate) Primary() byte {
	return template.primary
}

func (template extRegisterReadTemplate) Secondary() byte {
	return template.secondary
}

func (template extRegisterReadTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("ext register read template missing params: %w", ebuserrors.ErrInvalidPayload)
	}

	group, ok := uint8Param(params, "group")
	if !ok {
		return nil, fmt.Errorf("ext register read template group: %w", ebuserrors.ErrInvalidPayload)
	}
	instance, ok := uint8Param(params, "instance")
	if !ok {
		return nil, fmt.Errorf("ext register read template instance: %w", ebuserrors.ErrInvalidPayload)
	}
	addr, ok, err := registerAddrParam(params, "addr")
	if err != nil {
		return nil, fmt.Errorf("ext register read template addr: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("ext register read template addr: %w", ebuserrors.ErrInvalidPayload)
	}

	return []byte{extRegisterPrefix, extRegisterCmdRead, group, instance, byte(addr >> 8), byte(addr)}, nil
}

type extRegisterWriteTemplate struct {
	primary   byte
	secondary byte
}

func (template extRegisterWriteTemplate) Primary() byte {
	return template.primary
}

func (template extRegisterWriteTemplate) Secondary() byte {
	return template.secondary
}

func (template extRegisterWriteTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("ext register write template missing params: %w", ebuserrors.ErrInvalidPayload)
	}

	group, ok := uint8Param(params, "group")
	if !ok {
		return nil, fmt.Errorf("ext register write template group: %w", ebuserrors.ErrInvalidPayload)
	}
	instance, ok := uint8Param(params, "instance")
	if !ok {
		return nil, fmt.Errorf("ext register write template instance: %w", ebuserrors.ErrInvalidPayload)
	}
	addr, ok, err := registerAddrParam(params, "addr")
	if err != nil {
		return nil, fmt.Errorf("ext register write template addr: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("ext register write template addr: %w", ebuserrors.ErrInvalidPayload)
	}

	data, hasData, err := bytesParam(params, "data")
	if err != nil {
		return nil, fmt.Errorf("ext register write template data: %w", err)
	}
	if !hasData {
		data, _, err = bytesParam(params, "payload")
		if err != nil {
			return nil, fmt.Errorf("ext register write template payload: %w", err)
		}
	}

	if len(data) > 0xF9 {
		return nil, fmt.Errorf("ext register write template data too long: %w", ebuserrors.ErrInvalidPayload)
	}

	payload := make([]byte, 0, 6+len(data))
	payload = append(payload, extRegisterPrefix, extRegisterCmdWrite, group, instance, byte(addr>>8), byte(addr))
	payload = append(payload, data...)
	return payload, nil
}

func decodeExtRegisterResponse(cmd byte, group, instance byte, addr uint16, payload []byte) map[string]types.Value {
	values := map[string]types.Value{
		"cmd":      {Value: cmd, Valid: true},
		"group":    {Value: group, Valid: true},
		"instance": {Value: instance, Valid: true},
		"addr":     {Value: addr, Valid: true},
		"addr_hex": {Value: fmt.Sprintf("%04X", addr), Valid: true},
		"payload":  {Value: append([]byte(nil), payload...), Valid: true},
	}

	if len(payload) >= 4 {
		values["prefix"] = types.Value{Value: append([]byte(nil), payload[:4]...), Valid: true}
	} else {
		values["prefix"] = types.Value{Valid: false}
	}

	if len(payload) == 1 && payload[0] == 0x00 {
		values["value"] = types.Value{Valid: false}
		return values
	}
	if len(payload) <= 4 {
		values["value"] = types.Value{Valid: false}
		return values
	}

	values["value"] = types.Value{Value: append([]byte(nil), payload[4:]...), Valid: true}
	return values
}
