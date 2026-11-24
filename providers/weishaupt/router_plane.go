package weishaupt

import (
	"fmt"
	"math"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	weishauptcat "github.com/Project-Helianthus/helianthus-ebusreg/catalog/weishaupt"
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
		return protocol.Frame{}, fmt.Errorf("weishaupt BuildRequest missing method: %w", ebuserrors.ErrInvalidPayload)
	}
	if params == nil {
		return protocol.Frame{}, fmt.Errorf("weishaupt BuildRequest missing params: %w", ebuserrors.ErrInvalidPayload)
	}
	source, ok := uint8Param(params, "source")
	if !ok {
		return protocol.Frame{}, fmt.Errorf("weishaupt BuildRequest source: %w", ebuserrors.ErrInvalidPayload)
	}
	// wrsol_datetime_set is the documented 07 00 broadcast: ZZ is FEh for every
	// caller. The bound device address and a supplied target do not select it.
	target := p.address
	if method.Name() == methodDateTimeSet {
		target = 0xFE
	} else if target == 0 {
		target, ok = uint8Param(params, "target")
		if !ok {
			return protocol.Frame{}, fmt.Errorf("weishaupt BuildRequest target: %w", ebuserrors.ErrInvalidPayload)
		}
	}
	template := method.Template()
	if template == nil {
		return protocol.Frame{}, fmt.Errorf("weishaupt BuildRequest missing template: %w", ebuserrors.ErrInvalidPayload)
	}
	builder, ok := template.(interface {
		Build(params map[string]any) ([]byte, error)
	})
	if !ok {
		return protocol.Frame{}, fmt.Errorf("weishaupt BuildRequest unspecified payload: %w", ebuserrors.ErrInvalidPayload)
	}
	payload, err := builder.Build(params)
	if err != nil {
		return protocol.Frame{}, fmt.Errorf("weishaupt BuildRequest build payload: %w", err)
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
		return nil, fmt.Errorf("weishaupt DecodeResponse missing method: %w", ebuserrors.ErrInvalidPayload)
	}
	template := method.Template()
	if template == nil {
		return nil, fmt.Errorf("weishaupt DecodeResponse missing template: %w", ebuserrors.ErrInvalidPayload)
	}
	if response.Primary != template.Primary() || response.Secondary != template.Secondary() {
		return nil, fmt.Errorf("weishaupt DecodeResponse unexpected response type: %w", ebuserrors.ErrInvalidPayload)
	}
	switch method.Name() {
	case methodIdentificationRead:
		return weishauptcat.DecodeIdentification(response.Data)
	case methodRAMRead:
		return decodeRAM(params, response.Data)
	case methodMemoryRead:
		return weishauptcat.DecodeMemoryReply(response.Data)
	case methodDateTimeSet, methodMemoryWrite:
		return append([]byte(nil), response.Data...), nil
	case methodSecondaryRead, methodSecondaryWrite:
		return nil, unspecified(method.Name())
	default:
		return nil, fmt.Errorf("weishaupt DecodeResponse unknown method %q: %w", method.Name(), ebuserrors.ErrInvalidPayload)
	}
}

// Build returns the documented request payload for this method.
func (m method) Build(params map[string]any) ([]byte, error) {
	switch m.name {
	case methodIdentificationRead:
		return buildIdentification()
	case methodRAMRead:
		return buildRAMRead(params)
	case methodDateTimeSet:
		return buildDateTime(params)
	case methodMemoryRead:
		return buildMemoryRead(params)
	case methodMemoryWrite:
		return buildMemoryWrite(params)
	case methodSecondaryRead, methodSecondaryWrite:
		return nil, unspecified(m.name)
	default:
		return nil, fmt.Errorf("weishaupt BuildRequest unspecified method %q: %w", m.name, ebuserrors.ErrInvalidPayload)
	}
}

// buildIdentification returns the 07 04 query. The catalogue documents the
// ten-byte reply and no request fields, so the payload is empty.
func buildIdentification() ([]byte, error) {
	return []byte{}, nil
}

func buildRAMRead(params map[string]any) ([]byte, error) {
	address, ok := uint16Param(params, "address")
	if !ok {
		return nil, invalidParam("address")
	}
	return weishauptcat.EncodeRAMRead(address), nil
}

func buildDateTime(params map[string]any) ([]byte, error) {
	outdoor, ok := floatParam(params, "outdoor_c")
	if !ok {
		return nil, invalidParam("outdoor_c")
	}
	when, err := dateTimeParam(params)
	if err != nil {
		return nil, err
	}
	payload, err := weishauptcat.EncodeDateTime(outdoor, when)
	if err != nil {
		return nil, fmt.Errorf("weishaupt datetime: %w", err)
	}
	return payload, nil
}

func dateTimeParam(params map[string]any) (weishauptcat.DateTime, error) {
	var when weishauptcat.DateTime
	fields := []struct {
		key  string
		dest *int
	}{
		{"seconds", &when.Seconds},
		{"minutes", &when.Minutes},
		{"hours", &when.Hours},
		{"day", &when.Day},
		{"month", &when.Month},
		{"weekday", &when.Weekday},
		{"year", &when.Year},
	}
	for _, field := range fields {
		value, ok := smallInt(params, field.key)
		if !ok {
			return weishauptcat.DateTime{}, invalidParam(field.key)
		}
		*field.dest = value
	}
	return when, nil
}

func buildMemoryRead(params map[string]any) ([]byte, error) {
	id, err := memoryID(params)
	if err != nil {
		return nil, err
	}
	return weishauptcat.EncodeMemoryRead(id), nil
}

func buildMemoryWrite(params map[string]any) ([]byte, error) {
	id, err := memoryID(params)
	if err != nil {
		return nil, err
	}
	value, present, err := bytesParam(params, "value")
	if err != nil {
		return nil, fmt.Errorf("weishaupt memory value: %w", err)
	}
	if !present {
		return nil, invalidParam("value")
	}
	payload, err := weishauptcat.EncodeMemoryWrite(id, value)
	if err != nil {
		return nil, fmt.Errorf("weishaupt memory write: %w", err)
	}
	return payload, nil
}

func memoryID(params map[string]any) (weishauptcat.MemoryID, error) {
	id, present, err := bytesParam(params, "id")
	if err != nil {
		return weishauptcat.MemoryID{}, fmt.Errorf("weishaupt memory id: %w", err)
	}
	if !present || len(id) != 3 {
		return weishauptcat.MemoryID{}, invalidParam("id")
	}
	var key weishauptcat.MemoryID
	copy(key[:], id)
	return key, nil
}

func decodeRAM(params map[string]any, payload []byte) (any, error) {
	address, ok := uint16Param(params, "address")
	if !ok {
		return nil, invalidParam("address")
	}
	reg, found := weishauptcat.RAMByAddress(address)
	if !found {
		if len(payload) != 2 {
			return nil, fmt.Errorf("%w: RAM read wants 2 bytes", weishauptcat.ErrWidth)
		}
		return append([]byte(nil), payload...), nil
	}
	field, err := weishauptcat.DecodeRAMRead(reg, payload)
	if err != nil {
		return nil, fmt.Errorf("weishaupt ram: %w", err)
	}
	return field, nil
}

func unspecified(name string) error {
	return fmt.Errorf("weishaupt %s has no documented request layout: %w", name, ebuserrors.ErrInvalidPayload)
}

func invalidParam(name string) error {
	return fmt.Errorf("weishaupt %s: %w", name, ebuserrors.ErrInvalidPayload)
}

func uint8Param(params map[string]any, key string) (byte, bool) {
	value, ok := uintParam(params, key, 0xFF)
	return byte(value), ok
}

func uint16Param(params map[string]any, key string) (uint16, bool) {
	value, ok := uintParam(params, key, 0xFFFF)
	return uint16(value), ok
}

func smallInt(params map[string]any, key string) (int, bool) {
	value, ok := uintParam(params, key, 255)
	return int(value), ok
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
