package brink

import (
	"encoding/json"
	"fmt"
	"math"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	brinkcat "github.com/Project-Helianthus/helianthus-ebusreg/catalog/brink"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/router"
)

var (
	_ router.Plane   = (*plane)(nil)
	_ registry.Plane = (*plane)(nil)
)

// Subscriptions returns no broadcast list. Hand invocation sends a directed
// request; it does not register a passive listener.
func (p *plane) Subscriptions() []router.Subscription {
	return nil
}

// OnBroadcast implements router.Plane. This plane does not consume broadcasts.
func (p *plane) OnBroadcast(protocol.Frame) error {
	return nil
}

// BuildRequest builds one documented request. An unspecified layout, a missing
// parameter, or a value the catalogue rejects is returned before any bus send.
func (p *plane) BuildRequest(method registry.Method, params map[string]any) (protocol.Frame, error) {
	if method == nil {
		return protocol.Frame{}, fmt.Errorf("brink BuildRequest missing method: %w", ebuserrors.ErrInvalidPayload)
	}
	if params == nil {
		return protocol.Frame{}, fmt.Errorf("brink BuildRequest missing params: %w", ebuserrors.ErrInvalidPayload)
	}
	source, ok := uint8Param(params, "source")
	if !ok {
		return protocol.Frame{}, fmt.Errorf("brink BuildRequest source: %w", ebuserrors.ErrInvalidPayload)
	}
	target := p.address
	if target == 0 {
		target, ok = uint8Param(params, "target")
		if !ok {
			return protocol.Frame{}, fmt.Errorf("brink BuildRequest target: %w", ebuserrors.ErrInvalidPayload)
		}
	}
	template := method.Template()
	if template == nil {
		return protocol.Frame{}, fmt.Errorf("brink BuildRequest missing template: %w", ebuserrors.ErrInvalidPayload)
	}
	builder, ok := template.(interface {
		Build(params map[string]any) ([]byte, error)
	})
	if !ok {
		return protocol.Frame{}, fmt.Errorf("brink BuildRequest unspecified payload: %w", ebuserrors.ErrInvalidPayload)
	}
	payload, err := builder.Build(params)
	if err != nil {
		return protocol.Frame{}, fmt.Errorf("brink BuildRequest build payload: %w", err)
	}
	return protocol.Frame{
		Source:    source,
		Target:    target,
		Primary:   template.Primary(),
		Secondary: template.Secondary(),
		Data:      payload,
	}, nil
}

// DecodeResponse decodes a documented reply. A reply the page does not define
// is returned raw. The primary and secondary bytes must match the method.
func (p *plane) DecodeResponse(method registry.Method, response protocol.Frame, params map[string]any) (any, error) {
	if method == nil {
		return nil, fmt.Errorf("brink DecodeResponse missing method: %w", ebuserrors.ErrInvalidPayload)
	}
	template := method.Template()
	if template == nil {
		return nil, fmt.Errorf("brink DecodeResponse missing template: %w", ebuserrors.ErrInvalidPayload)
	}
	if response.Primary != template.Primary() || response.Secondary != template.Secondary() {
		return nil, fmt.Errorf("brink DecodeResponse unexpected response type: %w", ebuserrors.ErrInvalidPayload)
	}
	switch method.Name() {
	case methodSensorRead:
		return decodeSensor(params, response.Data)
	case methodParameterRead:
		return decodeParameter(params, response.Data)
	case methodErrorHistoryRead, methodFactoryReset, methodResetNotifications,
		methodFanModeTemporary, methodFanModePersistent, methodParameterWrite:
		return append([]byte(nil), response.Data...), nil
	default:
		return nil, fmt.Errorf("brink DecodeResponse unknown method %q: %w", method.Name(), ebuserrors.ErrInvalidPayload)
	}
}

// Build returns the documented request payload for this method.
func (m method) Build(params map[string]any) ([]byte, error) {
	switch m.name {
	case methodSensorRead:
		return buildSensorRead(params)
	case methodParameterRead:
		return buildParameterRead(params)
	case methodErrorHistoryRead:
		return brinkcat.EncodeErrorHistoryRequest(), nil
	case methodFactoryReset:
		return brinkcat.EncodeFactoryReset(), nil
	case methodResetNotifications:
		return brinkcat.EncodeResetNotifications(), nil
	case methodFanModeTemporary:
		return buildFanMode(params, brinkcat.EncodeFanModeTemporary)
	case methodFanModePersistent:
		return buildFanMode(params, brinkcat.EncodeFanModePersistent)
	case methodParameterWrite:
		return buildParameterWrite(params)
	default:
		return nil, fmt.Errorf("brink BuildRequest unspecified method %q: %w", m.name, ebuserrors.ErrInvalidPayload)
	}
}

func buildSensorRead(params map[string]any) ([]byte, error) {
	id, err := requireByte(params, "id")
	if err != nil {
		return nil, err
	}
	return brinkcat.EncodeSensorRequest(id), nil
}

func buildParameterRead(params map[string]any) ([]byte, error) {
	id, err := knownParameter(params)
	if err != nil {
		return nil, err
	}
	return brinkcat.EncodeParameterRead(id), nil
}

func buildFanMode(params map[string]any, encode func(byte) []byte) ([]byte, error) {
	mode, err := requireByte(params, "mode")
	if err != nil {
		return nil, err
	}
	return encode(mode), nil
}

func buildParameterWrite(params map[string]any) ([]byte, error) {
	id, err := knownParameter(params)
	if err != nil {
		return nil, err
	}
	raw, ok := int64Param(params, "value")
	if !ok {
		return nil, invalidParam("value")
	}
	payload, err := brinkcat.EncodeParameterWrite(id, raw)
	if err != nil {
		return nil, fmt.Errorf("brink parameter write: %w", err)
	}
	return payload, nil
}

func knownParameter(params map[string]any) (byte, error) {
	id, err := requireByte(params, "id")
	if err != nil {
		return 0, err
	}
	if _, found := brinkcat.ParameterByID(id); !found {
		return 0, fmt.Errorf("%w: parameter 0x%02X", brinkcat.ErrUnknownID, id)
	}
	return id, nil
}

func decodeSensor(params map[string]any, payload []byte) (brinkcat.Reading, error) {
	id, err := requireByte(params, "id")
	if err != nil {
		return brinkcat.Reading{}, err
	}
	reading, err := brinkcat.DecodeSensorReply(id, payload)
	if err != nil {
		return brinkcat.Reading{}, fmt.Errorf("brink sensor read: %w", err)
	}
	return reading, nil
}

func decodeParameter(params map[string]any, payload []byte) (brinkcat.ParameterReading, error) {
	id, err := requireByte(params, "id")
	if err != nil {
		return brinkcat.ParameterReading{}, err
	}
	reading, err := brinkcat.DecodeParameterReply(id, payload)
	if err != nil {
		return brinkcat.ParameterReading{}, fmt.Errorf("brink parameter read: %w", err)
	}
	return reading, nil
}

func requireByte(params map[string]any, key string) (byte, error) {
	value, ok := uint8Param(params, key)
	if !ok {
		return 0, invalidParam(key)
	}
	return value, nil
}

func invalidParam(name string) error {
	return fmt.Errorf("brink %s: %w", name, ebuserrors.ErrInvalidPayload)
}

func uint8Param(params map[string]any, key string) (byte, bool) {
	value, ok := params[key]
	if !ok || value == nil {
		return 0, false
	}
	return integralByte(value)
}

// integralByte accepts a finite nonnegative integer in 0..255. JSON object
// decoding yields float64 for numbers, so that width is required for a hand
// invocation. Fractions and values outside the byte range are rejected.
func integralByte(value any) (byte, bool) {
	switch typed := value.(type) {
	case byte:
		return typed, true
	case int:
		return byteFromSigned(int64(typed))
	case int8:
		return byteFromSigned(int64(typed))
	case int16:
		return byteFromSigned(int64(typed))
	case int32:
		return byteFromSigned(int64(typed))
	case int64:
		return byteFromSigned(typed)
	case uint:
		return byteFromUint64(uint64(typed))
	case uint16:
		return byteFromUint64(uint64(typed))
	case uint32:
		return byteFromUint64(uint64(typed))
	case uint64:
		return byteFromUint64(typed)
	case float32:
		return byteFromFloat(float64(typed))
	case float64:
		return byteFromFloat(typed)
	case json.Number:
		if parsed, err := typed.Int64(); err == nil {
			return byteFromSigned(parsed)
		}
		floatParsed, err := typed.Float64()
		if err != nil {
			return 0, false
		}
		return byteFromFloat(floatParsed)
	default:
		return 0, false
	}
}

func byteFromSigned(value int64) (byte, bool) {
	if value < 0 || value > 255 {
		return 0, false
	}
	return byte(value), true
}

func byteFromUint64(value uint64) (byte, bool) {
	if value > 255 {
		return 0, false
	}
	return byte(value), true
}

func byteFromFloat(value float64) (byte, bool) {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || value > 255 || math.Trunc(value) != value {
		return 0, false
	}
	return byte(value), true
}

func int64Param(params map[string]any, key string) (int64, bool) {
	value, ok := params[key]
	if !ok || value == nil {
		return 0, false
	}
	switch typed := value.(type) {
	case int:
		return int64(typed), true
	case int64:
		return typed, true
	case int32:
		return int64(typed), true
	case byte:
		return int64(typed), true
	case uint16:
		return int64(typed), true
	case float64:
		if math.IsNaN(typed) || math.IsInf(typed, 0) || math.Trunc(typed) != typed {
			return 0, false
		}
		if typed < math.MinInt64 || typed > float64(math.MaxInt64) {
			return 0, false
		}
		return int64(typed), true
	default:
		return 0, false
	}
}
