package system

import (
	"fmt"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func (plane *plane) OnBroadcast(frame protocol.Frame) error {
	return nil
}

func (plane *plane) BuildRequest(method registry.Method, params map[string]any) (protocol.Frame, error) {
	if method == nil {
		return protocol.Frame{}, fmt.Errorf("system BuildRequest missing method: %w", ebuserrors.ErrInvalidPayload)
	}
	if params == nil {
		return protocol.Frame{}, fmt.Errorf("system BuildRequest missing params: %w", ebuserrors.ErrInvalidPayload)
	}

	source, ok := uint8Param(params, "source")
	if !ok {
		return protocol.Frame{}, fmt.Errorf("system BuildRequest source: %w", ebuserrors.ErrInvalidPayload)
	}

	target := plane.address
	if target == 0 {
		targetParam, ok := uint8Param(params, "target")
		if !ok {
			return protocol.Frame{}, fmt.Errorf("system BuildRequest target: %w", ebuserrors.ErrInvalidPayload)
		}
		target = targetParam
	}

	template := method.Template()
	if template == nil {
		return protocol.Frame{}, fmt.Errorf("system BuildRequest missing template: %w", ebuserrors.ErrInvalidPayload)
	}

	var payload []byte
	if builder, ok := template.(interface {
		Build(params map[string]any) ([]byte, error)
	}); ok {
		value, err := builder.Build(params)
		if err != nil {
			return protocol.Frame{}, fmt.Errorf("system BuildRequest build payload: %w", err)
		}
		payload = value
	}

	return protocol.Frame{
		Source:    source,
		Target:    target,
		Primary:   template.Primary(),
		Secondary: template.Secondary(),
		Data:      payload,
	}, nil
}

func (plane *plane) DecodeResponse(method registry.Method, response protocol.Frame, params map[string]any) (any, error) {
	if method == nil {
		return nil, fmt.Errorf("system DecodeResponse missing method: %w", ebuserrors.ErrInvalidPayload)
	}
	template := method.Template()
	if template == nil {
		return nil, fmt.Errorf("system DecodeResponse missing template: %w", ebuserrors.ErrInvalidPayload)
	}
	if response.Primary != template.Primary() || response.Secondary != template.Secondary() {
		return nil, fmt.Errorf("system DecodeResponse unexpected response type: %w", ebuserrors.ErrInvalidPayload)
	}

	switch method.Name() {
	case methodGetOperationalData:
		op, ok := uint8Param(params, "op")
		if !ok {
			return nil, fmt.Errorf("system DecodeResponse op: %w", ebuserrors.ErrInvalidPayload)
		}
		decoded, err := decodeOperationalData(op, response.Data)
		if err != nil {
			return nil, fmt.Errorf("system DecodeResponse operational: %w", err)
		}
		return decoded, nil
	case methodSetOperationalData:
		op, ok := uint8Param(params, "op")
		if !ok {
			return nil, fmt.Errorf("system DecodeResponse op: %w", ebuserrors.ErrInvalidPayload)
		}
		return decodeOperationalWriteResponse(op, response.Data), nil
	case methodGetRegister:
		addr, ok, err := registerAddrParam(params, "addr")
		if err != nil {
			return nil, fmt.Errorf("system DecodeResponse addr: %w", err)
		}
		if !ok {
			return nil, fmt.Errorf("system DecodeResponse addr: %w", ebuserrors.ErrInvalidPayload)
		}
		return decodeRegisterResponse(registerCmdRead, addr, response.Data), nil
	case methodSetRegister:
		addr, ok, err := registerAddrParam(params, "addr")
		if err != nil {
			return nil, fmt.Errorf("system DecodeResponse addr: %w", err)
		}
		if !ok {
			return nil, fmt.Errorf("system DecodeResponse addr: %w", ebuserrors.ErrInvalidPayload)
		}
		return decodeRegisterResponse(registerCmdWrite, addr, response.Data), nil
	case methodGetExtRegister:
		group, ok := uint8Param(params, "group")
		if !ok {
			return nil, fmt.Errorf("system DecodeResponse group: %w", ebuserrors.ErrInvalidPayload)
		}
		instance, ok := uint8Param(params, "instance")
		if !ok {
			return nil, fmt.Errorf("system DecodeResponse instance: %w", ebuserrors.ErrInvalidPayload)
		}
		opcode, err := extRegisterOpcode(params)
		if err != nil {
			return nil, fmt.Errorf("system DecodeResponse opcode: %w", err)
		}
		addr, ok, err := registerAddrParam(params, "addr")
		if err != nil {
			return nil, fmt.Errorf("system DecodeResponse addr: %w", err)
		}
		if !ok {
			return nil, fmt.Errorf("system DecodeResponse addr: %w", ebuserrors.ErrInvalidPayload)
		}
		return decodeExtRegisterResponse(extRegisterOpRead, opcode, group, instance, addr, response.Data), nil
	case methodSetExtRegister:
		group, ok := uint8Param(params, "group")
		if !ok {
			return nil, fmt.Errorf("system DecodeResponse group: %w", ebuserrors.ErrInvalidPayload)
		}
		instance, ok := uint8Param(params, "instance")
		if !ok {
			return nil, fmt.Errorf("system DecodeResponse instance: %w", ebuserrors.ErrInvalidPayload)
		}
		opcode, err := extRegisterOpcode(params)
		if err != nil {
			return nil, fmt.Errorf("system DecodeResponse opcode: %w", err)
		}
		addr, ok, err := registerAddrParam(params, "addr")
		if err != nil {
			return nil, fmt.Errorf("system DecodeResponse addr: %w", err)
		}
		if !ok {
			return nil, fmt.Errorf("system DecodeResponse addr: %w", ebuserrors.ErrInvalidPayload)
		}
		return decodeExtRegisterResponse(extRegisterOpWrite, opcode, group, instance, addr, response.Data), nil
	default:
		return nil, fmt.Errorf("system DecodeResponse unknown method %q: %w", method.Name(), ebuserrors.ErrInvalidPayload)
	}
}
