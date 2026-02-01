package schema

import (
	"fmt"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
	"github.com/d3vi1/helianthus-ebusgo/types"
)

// SchemaField describes a single named value in a structured payload.
type SchemaField struct {
	Name string
	Type types.DataType
}

// Schema defines a sequential payload structure.
type Schema struct {
	Fields []SchemaField
}

// Decode parses a payload into named values.
func (s Schema) Decode(payload []byte) (map[string]types.Value, error) {
	values := make(map[string]types.Value, len(s.Fields))
	offset := 0
	for _, field := range s.Fields {
		if field.Type == nil {
			return nil, fmt.Errorf("field %q missing type: %w", field.Name, ebuserrors.ErrInvalidPayload)
		}
		size := field.Type.Size()
		if size <= 0 {
			return nil, fmt.Errorf("field %q invalid size: %w", field.Name, ebuserrors.ErrInvalidPayload)
		}
		if offset+size > len(payload) {
			return nil, fmt.Errorf("field %q short payload: %w", field.Name, ebuserrors.ErrInvalidPayload)
		}
		value, err := field.Type.Decode(payload[offset : offset+size])
		if err != nil {
			return nil, fmt.Errorf("field %q decode: %w", field.Name, err)
		}
		values[field.Name] = value
		offset += size
	}
	return values, nil
}

// Encode builds a payload from named values.
func (s Schema) Encode(values map[string]any) ([]byte, error) {
	if len(s.Fields) == 0 {
		return nil, nil
	}
	if values == nil {
		return nil, fmt.Errorf("missing values: %w", ebuserrors.ErrInvalidPayload)
	}

	payload := make([]byte, 0, s.Size())
	for _, field := range s.Fields {
		if field.Type == nil {
			return nil, fmt.Errorf("field %q missing type: %w", field.Name, ebuserrors.ErrInvalidPayload)
		}
		value, ok := values[field.Name]
		if !ok {
			return nil, fmt.Errorf("field %q missing value: %w", field.Name, ebuserrors.ErrInvalidPayload)
		}
		encoded, err := field.Type.Encode(value)
		if err != nil {
			return nil, fmt.Errorf("field %q encode: %w", field.Name, err)
		}
		if size := field.Type.Size(); size > 0 && len(encoded) != size {
			return nil, fmt.Errorf("field %q encoded size mismatch: %w", field.Name, ebuserrors.ErrInvalidPayload)
		}
		payload = append(payload, encoded...)
	}
	return payload, nil
}

// Size returns the total byte size required for the schema.
func (s Schema) Size() int {
	total := 0
	for _, field := range s.Fields {
		if field.Type == nil {
			continue
		}
		size := field.Type.Size()
		if size > 0 {
			total += size
		}
	}
	return total
}
