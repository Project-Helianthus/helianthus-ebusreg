package wolf

import (
	"fmt"
	"math"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	wolfcat "github.com/Project-Helianthus/helianthus-ebusreg/catalog/wolf"
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
		return protocol.Frame{}, fmt.Errorf("wolf BuildRequest missing method: %w", ebuserrors.ErrInvalidPayload)
	}
	if params == nil {
		return protocol.Frame{}, fmt.Errorf("wolf BuildRequest missing params: %w", ebuserrors.ErrInvalidPayload)
	}
	source, ok := uint8Param(params, "source")
	if !ok {
		return protocol.Frame{}, fmt.Errorf("wolf BuildRequest source: %w", ebuserrors.ErrInvalidPayload)
	}
	target := p.address
	if target == 0 {
		target, ok = uint8Param(params, "target")
		if !ok {
			return protocol.Frame{}, fmt.Errorf("wolf BuildRequest target: %w", ebuserrors.ErrInvalidPayload)
		}
	}
	template := method.Template()
	if template == nil {
		return protocol.Frame{}, fmt.Errorf("wolf BuildRequest missing template: %w", ebuserrors.ErrInvalidPayload)
	}
	builder, ok := template.(interface {
		Build(params map[string]any) ([]byte, error)
	})
	if !ok {
		return protocol.Frame{}, fmt.Errorf("wolf BuildRequest unspecified payload: %w", ebuserrors.ErrInvalidPayload)
	}
	payload, err := builder.Build(params)
	if err != nil {
		return protocol.Frame{}, fmt.Errorf("wolf BuildRequest build payload: %w", err)
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
		return nil, fmt.Errorf("wolf DecodeResponse missing method: %w", ebuserrors.ErrInvalidPayload)
	}
	template := method.Template()
	if template == nil {
		return nil, fmt.Errorf("wolf DecodeResponse missing template: %w", ebuserrors.ErrInvalidPayload)
	}
	if response.Primary != template.Primary() || response.Secondary != template.Secondary() {
		return nil, fmt.Errorf("wolf DecodeResponse unexpected response type: %w", ebuserrors.ErrInvalidPayload)
	}
	switch method.Name() {
	case methodServiceRead:
		return decodeServiceRead(params, response.Data)
	case methodBurnerStatusSend, methodModulationSend, methodHeatRequestSend,
		methodBoilerSetpointSend, methodSolarStatusSend, methodSolarPumpSend,
		methodSolarPowerSend, methodServiceWrite:
		return append([]byte(nil), response.Data...), nil
	default:
		return nil, fmt.Errorf("wolf DecodeResponse unknown method %q: %w", method.Name(), ebuserrors.ErrInvalidPayload)
	}
}

// Build returns the documented request payload for this method.
func (m method) Build(params map[string]any) ([]byte, error) {
	switch m.name {
	case methodBurnerStatusSend:
		return buildBurnerStatus(params)
	case methodModulationSend:
		return buildModulation(params)
	case methodHeatRequestSend:
		return buildHeatRequest(params)
	case methodBoilerSetpointSend:
		return buildBoilerSetpoint(params)
	case methodSolarStatusSend:
		return buildSolarStatus(params)
	case methodSolarPumpSend:
		return buildSolarPump(params)
	case methodSolarPowerSend:
		return buildSolarPower(params)
	case methodServiceRead:
		return buildServiceRead(params)
	case methodServiceWrite:
		return buildServiceWrite(params)
	default:
		return nil, fmt.Errorf("wolf BuildRequest unspecified method %q: %w", m.name, ebuserrors.ErrInvalidPayload)
	}
}

func buildBurnerStatus(params map[string]any) ([]byte, error) {
	state, err := requireByte(params, "state")
	if err != nil {
		return nil, err
	}
	modulation, err := requireByte(params, "modulation")
	if err != nil {
		return nil, err
	}
	temperature, err := optionalBytes(params, "kessel_temperature")
	if err != nil {
		return nil, err
	}
	return wolfcat.EncodeBurnerStatus(state, modulation, temperature), nil
}

func buildModulation(params map[string]any) ([]byte, error) {
	actual, err := requireByte(params, "actual")
	if err != nil {
		return nil, err
	}
	maximum, err := requireByte(params, "maximum")
	if err != nil {
		return nil, err
	}
	payload, err := wolfcat.EncodeModulation(actual, maximum)
	if err != nil {
		return nil, fmt.Errorf("wolf modulation: %w", err)
	}
	return payload, nil
}

func buildHeatRequest(params map[string]any) ([]byte, error) {
	state, err := requireByte(params, "state")
	if err != nil {
		return nil, err
	}
	setpoint, err := requireByte(params, "setpoint")
	if err != nil {
		return nil, err
	}
	pressure, err := optionalBytes(params, "pressure")
	if err != nil {
		return nil, err
	}
	return wolfcat.EncodeHeatRequest(state, setpoint, pressure), nil
}

func buildBoilerSetpoint(params map[string]any) ([]byte, error) {
	celsius, ok := floatParam(params, "celsius")
	if !ok {
		return nil, invalidParam("celsius")
	}
	forced, err := requireByte(params, "forced_output")
	if err != nil {
		return nil, err
	}
	status, err := requireByte(params, "status")
	if err != nil {
		return nil, err
	}
	payload, err := wolfcat.EncodeBoilerSetpoint(celsius, forced, status)
	if err != nil {
		return nil, fmt.Errorf("wolf boiler setpoint: %w", err)
	}
	return payload, nil
}

func buildSolarStatus(params map[string]any) ([]byte, error) {
	status, err := requireByte(params, "status")
	if err != nil {
		return nil, err
	}
	ignored, err := requireByte(params, "ignored")
	if err != nil {
		return nil, err
	}
	rest, err := optionalBytes(params, "rest")
	if err != nil {
		return nil, err
	}
	return wolfcat.EncodeSolarStatus(status, ignored, rest), nil
}

func buildSolarPump(params map[string]any) ([]byte, error) {
	on, ok := boolParam(params, "on")
	if !ok {
		return nil, invalidParam("on")
	}
	ignored, err := requireByte(params, "ignored")
	if err != nil {
		return nil, err
	}
	temperatures, err := optionalBytes(params, "temperatures")
	if err != nil {
		return nil, err
	}
	return wolfcat.EncodeSolarPump(on, ignored, temperatures), nil
}

func buildSolarPower(params map[string]any) ([]byte, error) {
	watts, ok := floatParam(params, "watts")
	if !ok {
		return nil, invalidParam("watts")
	}
	yieldHigh, err := requireByte(params, "yield_high")
	if err != nil {
		return nil, err
	}
	yieldSumLow, err := requireByte(params, "yield_sum_low")
	if err != nil {
		return nil, err
	}
	payload, err := wolfcat.EncodeSolarPower(watts, yieldHigh, yieldSumLow)
	if err != nil {
		return nil, fmt.Errorf("wolf solar power: %w", err)
	}
	return payload, nil
}

func buildServiceRead(params map[string]any) ([]byte, error) {
	reg, err := serviceRegister(params)
	if err != nil {
		return nil, err
	}
	return wolfcat.EncodeServiceRead(reg), nil
}

func buildServiceWrite(params map[string]any) ([]byte, error) {
	reg, err := serviceRegister(params)
	if err != nil {
		return nil, err
	}
	raw, ok := uint16Param(params, "value")
	if !ok {
		return nil, invalidParam("value")
	}
	trailer, present, err := bytesParam(params, "trailer")
	if err != nil {
		return nil, fmt.Errorf("wolf trailer: %w", err)
	}
	if !present || len(trailer) != 4 {
		return nil, invalidParam("trailer")
	}
	var fixed [4]byte
	copy(fixed[:], trailer)
	payload, err := wolfcat.EncodeServiceWrite(reg, raw, fixed)
	if err != nil {
		return nil, fmt.Errorf("wolf service write: %w", err)
	}
	return payload, nil
}

func decodeServiceRead(params map[string]any, payload []byte) (wolfcat.Field, error) {
	reg, err := serviceRegister(params)
	if err != nil {
		return wolfcat.Field{}, err
	}
	field, err := wolfcat.DecodeServiceRead(reg, payload)
	if err != nil {
		return wolfcat.Field{}, fmt.Errorf("wolf service read: %w", err)
	}
	return field, nil
}

func serviceRegister(params map[string]any) (wolfcat.ServiceRegister, error) {
	id, present, err := bytesParam(params, "id")
	if err != nil {
		return wolfcat.ServiceRegister{}, fmt.Errorf("wolf service id: %w", err)
	}
	if !present || len(id) != 3 {
		return wolfcat.ServiceRegister{}, invalidParam("id")
	}
	var key [3]byte
	copy(key[:], id)
	reg, found := wolfcat.ServiceRegisterByID(key)
	if !found {
		return wolfcat.ServiceRegister{}, fmt.Errorf("wolf service id not in the documented table: %w", ebuserrors.ErrInvalidPayload)
	}
	return reg, nil
}

func requireByte(params map[string]any, key string) (byte, error) {
	value, ok := uint8Param(params, key)
	if !ok {
		return 0, invalidParam(key)
	}
	return value, nil
}

func optionalBytes(params map[string]any, key string) ([]byte, error) {
	value, present, err := bytesParam(params, key)
	if err != nil {
		return nil, fmt.Errorf("wolf %s: %w", key, err)
	}
	if !present {
		return nil, nil
	}
	return value, nil
}

func invalidParam(name string) error {
	return fmt.Errorf("wolf %s: %w", name, ebuserrors.ErrInvalidPayload)
}

func uint8Param(params map[string]any, key string) (byte, bool) {
	value, ok := uintParam(params, key, 0xFF)
	return byte(value), ok
}

func uint16Param(params map[string]any, key string) (uint16, bool) {
	value, ok := uintParam(params, key, 0xFFFF)
	return uint16(value), ok
}

func uintParam(params map[string]any, key string, max uint64) (uint64, bool) {
	value, ok := params[key]
	if !ok || value == nil {
		return 0, false
	}
	return uintValue(value, max)
}

func uintValue(value any, max uint64) (uint64, bool) {
	switch typed := value.(type) {
	case byte:
		return fitsUint(uint64(typed), max)
	case int:
		if typed < 0 {
			return 0, false
		}
		return fitsUint(uint64(typed), max)
	case int64:
		if typed < 0 {
			return 0, false
		}
		return fitsUint(uint64(typed), max)
	case uint16:
		return fitsUint(uint64(typed), max)
	case uint32:
		return fitsUint(uint64(typed), max)
	case float64:
		if math.IsNaN(typed) || math.IsInf(typed, 0) || typed < 0 || typed > float64(max) || math.Trunc(typed) != typed {
			return 0, false
		}
		return uint64(typed), true
	default:
		return 0, false
	}
}

func fitsUint(value, max uint64) (uint64, bool) {
	if value > max {
		return 0, false
	}
	return value, true
}

func floatParam(params map[string]any, key string) (float64, bool) {
	value, ok := params[key]
	if !ok || value == nil {
		return 0, false
	}
	switch typed := value.(type) {
	case float64:
		if math.IsNaN(typed) || math.IsInf(typed, 0) {
			return 0, false
		}
		return typed, true
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	default:
		return 0, false
	}
}

func boolParam(params map[string]any, key string) (bool, bool) {
	value, ok := params[key]
	if !ok {
		return false, false
	}
	typed, ok := value.(bool)
	return typed, ok
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
		for i, item := range typed {
			if item < 0 || item > 255 {
				return nil, true, ebuserrors.ErrInvalidPayload
			}
			out[i] = byte(item)
		}
		return out, true, nil
	case []any:
		// JSON arrays decode as []any of float64. Each element must be an
		// integral byte; a fraction or an out-of-range number is not sent.
		out := make([]byte, len(typed))
		for i, item := range typed {
			parsed, ok := uintValue(item, 0xFF)
			if !ok {
				return nil, true, ebuserrors.ErrInvalidPayload
			}
			out[i] = byte(parsed)
		}
		return out, true, nil
	default:
		return nil, true, ebuserrors.ErrInvalidPayload
	}
}
