package system

import (
	"fmt"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
	"github.com/d3vi1/helianthus-ebusgo/types"
	"github.com/d3vi1/helianthus-ebusreg/registry"
	"github.com/d3vi1/helianthus-ebusreg/schema"
)

const (
	methodGetOperationalData = "get_operational_data"
	methodSetOperationalData = "set_operational_data"
)

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

type operationalWriteTemplate struct {
	primary   byte
	secondary byte
}

func (template operationalWriteTemplate) Primary() byte {
	return template.primary
}

func (template operationalWriteTemplate) Secondary() byte {
	return template.secondary
}

func (template operationalWriteTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("operational write template missing params: %w", ebuserrors.ErrInvalidPayload)
	}

	op, ok := uint8Param(params, "op")
	if !ok {
		return nil, fmt.Errorf("operational write template op: %w", ebuserrors.ErrInvalidPayload)
	}

	data, hasData, err := bytesParam(params, "data")
	if err != nil {
		return nil, fmt.Errorf("operational write template data: %w", err)
	}
	if !hasData {
		data, _, err = bytesParam(params, "payload")
		if err != nil {
			return nil, fmt.Errorf("operational write template payload: %w", err)
		}
	}

	if len(data) > 0xFE {
		return nil, fmt.Errorf("operational write template data too long: %w", ebuserrors.ErrInvalidPayload)
	}

	payload := make([]byte, 0, 1+len(data))
	payload = append(payload, op)
	payload = append(payload, data...)
	return payload, nil
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

func decodeOperationalWriteResponse(op byte, payload []byte) map[string]types.Value {
	return map[string]types.Value{
		"op":      {Value: op, Valid: true},
		"payload": {Value: append([]byte(nil), payload...), Valid: true},
	}
}

func decodeOperationalDateTime(payload []byte) (map[string]types.Value, error) {
	const legacySize = 8
	const btiBdaSize = 10 // BTI (3 bytes) + BDA (4 bytes) + temp2 (2 bytes) + dcfstate

	if len(payload) < legacySize {
		return nil, fmt.Errorf("operational datetime short payload: %w", ebuserrors.ErrInvalidPayload)
	}

	var (
		hour, minute, day, month, year types.Value
		temp                           types.Value
		err                            error
	)
	if len(payload) >= btiBdaSize {
		// ebusd types:
		// - BTI (24-bit, BCD|REV): wire order is SS,MM,HH
		// - BDA (32-bit, BCD): DD,MM,<weekday>,YY
		hour, err = types.BCD{}.Decode(payload[3:4])
		if err != nil {
			return nil, fmt.Errorf("operational datetime hour: %w", err)
		}
		minute, err = types.BCD{}.Decode(payload[2:3])
		if err != nil {
			return nil, fmt.Errorf("operational datetime minute: %w", err)
		}
		day, err = types.BCD{}.Decode(payload[4:5])
		if err != nil {
			return nil, fmt.Errorf("operational datetime day: %w", err)
		}
		month, err = types.BCD{}.Decode(payload[5:6])
		if err != nil {
			return nil, fmt.Errorf("operational datetime month: %w", err)
		}
		year, err = types.BCD{}.Decode(payload[7:8])
		if err != nil {
			return nil, fmt.Errorf("operational datetime year: %w", err)
		}
		temp, err = types.DATA2b{}.Decode(payload[8:10])
		if err != nil {
			return nil, fmt.Errorf("operational datetime temp: %w", err)
		}
	} else {
		// Legacy layout (older configs): HH,MM,DD,MM,YY,temp2.
		hour, err = types.BCD{}.Decode(payload[1:2])
		if err != nil {
			return nil, fmt.Errorf("operational datetime hour: %w", err)
		}
		minute, err = types.BCD{}.Decode(payload[2:3])
		if err != nil {
			return nil, fmt.Errorf("operational datetime minute: %w", err)
		}
		day, err = types.BCD{}.Decode(payload[3:4])
		if err != nil {
			return nil, fmt.Errorf("operational datetime day: %w", err)
		}
		month, err = types.BCD{}.Decode(payload[4:5])
		if err != nil {
			return nil, fmt.Errorf("operational datetime month: %w", err)
		}
		year, err = types.BCD{}.Decode(payload[5:6])
		if err != nil {
			return nil, fmt.Errorf("operational datetime year: %w", err)
		}
		temp, err = types.DATA2b{}.Decode(payload[6:8])
		if err != nil {
			return nil, fmt.Errorf("operational datetime temp: %w", err)
		}
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
