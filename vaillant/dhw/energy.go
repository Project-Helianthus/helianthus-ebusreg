package dhw

import (
	"fmt"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
	"github.com/d3vi1/helianthus-ebusreg/internal/match"
	"github.com/d3vi1/helianthus-ebusreg/registry"
)

type energyTemplate struct {
	primary   byte
	secondary byte
}

func (template energyTemplate) Primary() byte {
	return template.primary
}

func (template energyTemplate) Secondary() byte {
	return template.secondary
}

func (template energyTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("energy template missing params: %w", ebuserrors.ErrInvalidPayload)
	}

	period, ok := uint8Param(params, "period")
	if !ok {
		return nil, fmt.Errorf("energy template period: %w", ebuserrors.ErrInvalidPayload)
	}
	source, ok := uint8Param(params, "source")
	if !ok {
		return nil, fmt.Errorf("energy template source: %w", ebuserrors.ErrInvalidPayload)
	}
	usage, ok := uint8Param(params, "usage")
	if !ok {
		return nil, fmt.Errorf("energy template usage: %w", ebuserrors.ErrInvalidPayload)
	}

	return []byte{period, source, usage}, nil
}

func supportsEnergyStats(info registry.DeviceInfo) bool {
	return match.HardwareVersionAtLeast(info.HardwareVersion, 7603)
}

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
