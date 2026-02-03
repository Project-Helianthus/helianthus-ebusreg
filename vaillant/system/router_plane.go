package system

import (
	"fmt"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
	"github.com/d3vi1/helianthus-ebusgo/protocol"
	"github.com/d3vi1/helianthus-ebusreg/registry"
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
	default:
		return nil, fmt.Errorf("system DecodeResponse unknown method %q: %w", method.Name(), ebuserrors.ErrInvalidPayload)
	}
}
