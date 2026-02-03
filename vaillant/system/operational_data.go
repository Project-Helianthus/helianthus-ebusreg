package system

import (
	"fmt"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
	"github.com/d3vi1/helianthus-ebusgo/types"
	"github.com/d3vi1/helianthus-ebusreg/registry"
	"github.com/d3vi1/helianthus-ebusreg/schema"
)

const methodGetOperationalData = "get_operational_data"

type method struct {
	name     string
	readOnly bool
	template registry.FrameTemplate
	response schema.SchemaSelector
}

func (method method) Name() string {
	return method.name
}

func (method method) ReadOnly() bool {
	return method.readOnly
}

func (method method) Template() registry.FrameTemplate {
	return method.template
}

func (method method) ResponseSchema() schema.SchemaSelector {
	return method.response
}

type operationalTemplate struct {
	primary   byte
	secondary byte
}

func (template operationalTemplate) Primary() byte {
	return template.primary
}

func (template operationalTemplate) Secondary() byte {
	return template.secondary
}

func (template operationalTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("operational template missing params: %w", ebuserrors.ErrInvalidPayload)
	}

	op, ok := uint8Param(params, "op")
	if !ok {
		return nil, fmt.Errorf("operational template op: %w", ebuserrors.ErrInvalidPayload)
	}
	return []byte{op}, nil
}

func decodeOperationalData(op byte, payload []byte) (map[string]types.Value, error) {
	values := map[string]types.Value{
		"op":      {Value: op, Valid: true},
		"payload": {Value: append([]byte(nil), payload...), Valid: true},
	}

	if op != 0x00 {
		return values, nil
	}

	decoded, err := decodeOperationalDateTime(payload)
	if err != nil {
		return nil, err
	}
	for key, value := range decoded {
		values[key] = value
	}
	return values, nil
}

func decodeOperationalDateTime(payload []byte) (map[string]types.Value, error) {
	const expectedSize = 8
	if len(payload) < expectedSize {
		return nil, fmt.Errorf("operational datetime short payload: %w", ebuserrors.ErrInvalidPayload)
	}

	hour, err := types.BCD{}.Decode(payload[1:2])
	if err != nil {
		return nil, fmt.Errorf("operational datetime hour: %w", err)
	}
	minute, err := types.BCD{}.Decode(payload[2:3])
	if err != nil {
		return nil, fmt.Errorf("operational datetime minute: %w", err)
	}
	day, err := types.BCD{}.Decode(payload[3:4])
	if err != nil {
		return nil, fmt.Errorf("operational datetime day: %w", err)
	}
	month, err := types.BCD{}.Decode(payload[4:5])
	if err != nil {
		return nil, fmt.Errorf("operational datetime month: %w", err)
	}
	year, err := types.BCD{}.Decode(payload[5:6])
	if err != nil {
		return nil, fmt.Errorf("operational datetime year: %w", err)
	}
	temp, err := types.DATA2b{}.Decode(payload[6:8])
	if err != nil {
		return nil, fmt.Errorf("operational datetime temp: %w", err)
	}

	return map[string]types.Value{
		"dcfstate":    {Value: uint8(payload[0]), Valid: true},
		"time_hour":   hour,
		"time_minute": minute,
		"date_day":    day,
		"date_month":  month,
		"date_year":   year,
		"temp":        temp,
	}, nil
}
