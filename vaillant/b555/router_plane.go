package b555

import (
	"fmt"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func (p *plane) OnBroadcast(frame protocol.Frame) error {
	return nil
}

func (p *plane) BuildRequest(method registry.Method, params map[string]any) (protocol.Frame, error) {
	if method == nil {
		return protocol.Frame{}, fmt.Errorf("b555 BuildRequest missing method: %w", ebuserrors.ErrInvalidPayload)
	}
	if params == nil {
		return protocol.Frame{}, fmt.Errorf("b555 BuildRequest missing params: %w", ebuserrors.ErrInvalidPayload)
	}

	source, ok := uint8Param(params, "source")
	if !ok {
		return protocol.Frame{}, fmt.Errorf("b555 BuildRequest source: %w", ebuserrors.ErrInvalidPayload)
	}

	target := p.address
	if target == 0 {
		targetParam, ok := uint8Param(params, "target")
		if !ok {
			return protocol.Frame{}, fmt.Errorf("b555 BuildRequest target: %w", ebuserrors.ErrInvalidPayload)
		}
		target = targetParam
	}

	template := method.Template()
	if template == nil {
		return protocol.Frame{}, fmt.Errorf("b555 BuildRequest missing template: %w", ebuserrors.ErrInvalidPayload)
	}

	var payload []byte
	if builder, ok := template.(interface {
		Build(params map[string]any) ([]byte, error)
	}); ok {
		value, err := builder.Build(params)
		if err != nil {
			return protocol.Frame{}, fmt.Errorf("b555 BuildRequest build payload: %w", err)
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

func (p *plane) DecodeResponse(method registry.Method, response protocol.Frame, params map[string]any) (any, error) {
	if method == nil {
		return nil, fmt.Errorf("b555 DecodeResponse missing method: %w", ebuserrors.ErrInvalidPayload)
	}
	template := method.Template()
	if template == nil {
		return nil, fmt.Errorf("b555 DecodeResponse missing template: %w", ebuserrors.ErrInvalidPayload)
	}
	if response.Primary != template.Primary() || response.Secondary != template.Secondary() {
		return nil, fmt.Errorf("b555 DecodeResponse unexpected response type: %w", ebuserrors.ErrInvalidPayload)
	}

	switch method.Name() {
	case methodReadTimerConfig:
		zone, ok := uint8Param(params, "zone")
		if !ok {
			return nil, fmt.Errorf("b555 DecodeResponse zone: %w", ebuserrors.ErrInvalidPayload)
		}
		hc, ok := uint8Param(params, "hc")
		if !ok {
			return nil, fmt.Errorf("b555 DecodeResponse hc: %w", ebuserrors.ErrInvalidPayload)
		}
		return decodeConfigResponse(zone, hc, response.Data), nil

	case methodReadTimerSlots:
		zone, ok := uint8Param(params, "zone")
		if !ok {
			return nil, fmt.Errorf("b555 DecodeResponse zone: %w", ebuserrors.ErrInvalidPayload)
		}
		hc, ok := uint8Param(params, "hc")
		if !ok {
			return nil, fmt.Errorf("b555 DecodeResponse hc: %w", ebuserrors.ErrInvalidPayload)
		}
		return decodeSlotsResponse(zone, hc, response.Data), nil

	case methodReadTimer:
		zone, ok := uint8Param(params, "zone")
		if !ok {
			return nil, fmt.Errorf("b555 DecodeResponse zone: %w", ebuserrors.ErrInvalidPayload)
		}
		hc, ok := uint8Param(params, "hc")
		if !ok {
			return nil, fmt.Errorf("b555 DecodeResponse hc: %w", ebuserrors.ErrInvalidPayload)
		}
		weekday, ok := uint8Param(params, "weekday")
		if !ok {
			return nil, fmt.Errorf("b555 DecodeResponse weekday: %w", ebuserrors.ErrInvalidPayload)
		}
		slot, ok := uint8Param(params, "slot")
		if !ok {
			return nil, fmt.Errorf("b555 DecodeResponse slot: %w", ebuserrors.ErrInvalidPayload)
		}
		return decodeTimerReadResponse(zone, hc, weekday, slot, response.Data), nil

	case methodWriteTimer:
		zone, ok := uint8Param(params, "zone")
		if !ok {
			return nil, fmt.Errorf("b555 DecodeResponse zone: %w", ebuserrors.ErrInvalidPayload)
		}
		hc, ok := uint8Param(params, "hc")
		if !ok {
			return nil, fmt.Errorf("b555 DecodeResponse hc: %w", ebuserrors.ErrInvalidPayload)
		}
		weekday, ok := uint8Param(params, "weekday")
		if !ok {
			return nil, fmt.Errorf("b555 DecodeResponse weekday: %w", ebuserrors.ErrInvalidPayload)
		}
		slot, ok := uint8Param(params, "slot")
		if !ok {
			return nil, fmt.Errorf("b555 DecodeResponse slot: %w", ebuserrors.ErrInvalidPayload)
		}
		return decodeTimerWriteResponse(zone, hc, weekday, slot, response.Data), nil

	default:
		return nil, fmt.Errorf("b555 DecodeResponse unknown method %q: %w", method.Name(), ebuserrors.ErrInvalidPayload)
	}
}
