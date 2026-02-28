//go:build !tinygo

package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/types"
)

type jsonSchema struct {
	Fields []jsonField `json:"fields"`
}

type jsonField struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size,omitempty"`
}

type jsonSelector struct {
	Default    jsonSchema      `json:"default"`
	Conditions []jsonCondition `json:"conditions"`
}

type jsonCondition struct {
	Target *uint8     `json:"target,omitempty"`
	MinHW  *uint16    `json:"min_hw,omitempty"`
	MaxHW  *uint16    `json:"max_hw,omitempty"`
	Schema jsonSchema `json:"schema"`
}

// LoadSchemaJSON parses a JSON schema definition.
func LoadSchemaJSON(data []byte) (Schema, error) {
	var raw jsonSchema
	if err := json.Unmarshal(data, &raw); err != nil {
		return Schema{}, fmt.Errorf("schema json: %w", err)
	}
	return decodeJSONSchema(raw)
}

// LoadSchemaFile parses a JSON schema file.
func LoadSchemaFile(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Schema{}, fmt.Errorf("schema file %s: %w", path, err)
	}
	return LoadSchemaJSON(data)
}

// LoadSchemaSelectorJSON parses a JSON schema selector definition.
func LoadSchemaSelectorJSON(data []byte) (SchemaSelector, error) {
	var raw jsonSelector
	if err := json.Unmarshal(data, &raw); err != nil {
		return SchemaSelector{}, fmt.Errorf("schema selector json: %w", err)
	}
	defaultSchema, err := decodeJSONSchema(raw.Default)
	if err != nil {
		return SchemaSelector{}, err
	}

	conditions := make([]SchemaCondition, 0, len(raw.Conditions))
	for _, entry := range raw.Conditions {
		schemaValue, err := decodeJSONSchema(entry.Schema)
		if err != nil {
			return SchemaSelector{}, err
		}
		condition := SchemaCondition{Schema: schemaValue}
		if entry.Target != nil {
			condition.Target = *entry.Target
			condition.HasTarget = true
		}
		if entry.MinHW != nil {
			condition.MinHW = *entry.MinHW
			condition.HasMinHW = true
		}
		if entry.MaxHW != nil {
			condition.MaxHW = *entry.MaxHW
			condition.HasMaxHW = true
		}
		conditions = append(conditions, condition)
	}

	return SchemaSelector{
		Default:    defaultSchema,
		Conditions: conditions,
	}, nil
}

// LoadSchemaSelectorFile parses a JSON schema selector file.
func LoadSchemaSelectorFile(path string) (SchemaSelector, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SchemaSelector{}, fmt.Errorf("schema selector file %s: %w", path, err)
	}
	return LoadSchemaSelectorJSON(data)
}

func decodeJSONSchema(raw jsonSchema) (Schema, error) {
	fields := make([]SchemaField, 0, len(raw.Fields))
	for _, field := range raw.Fields {
		if field.Name == "" {
			return Schema{}, fmt.Errorf("schema field missing name: %w", ebuserrors.ErrInvalidPayload)
		}
		if field.Type == "" {
			return Schema{}, fmt.Errorf("schema field %q missing type: %w", field.Name, ebuserrors.ErrInvalidPayload)
		}
		dataType, err := dataTypeByName(field.Type, field.Size)
		if err != nil {
			return Schema{}, fmt.Errorf("schema field %q type %q: %w", field.Name, field.Type, err)
		}
		fields = append(fields, SchemaField{
			Name: field.Name,
			Type: dataType,
		})
	}
	return Schema{Fields: fields}, nil
}

func dataTypeByName(name string, size int) (types.DataType, error) {
	upper := strings.ToUpper(strings.TrimSpace(name))
	switch upper {
	case "DATA1B":
		return types.DATA1b{}, nil
	case "DATA2C":
		return types.DATA2c{}, nil
	case "DATA2B":
		return types.DATA2b{}, nil
	case "EXP":
		return types.EXP{}, nil
	case "WORD":
		return types.WORD{}, nil
	case "BCD":
		return types.BCD{}, nil
	case "BITFIELD":
		if size <= 0 {
			return nil, fmt.Errorf("bitfield missing size: %w", ebuserrors.ErrInvalidPayload)
		}
		return types.BITFIELD{SizeBytes: size}, nil
	default:
		return nil, fmt.Errorf("unknown data type: %w", ebuserrors.ErrInvalidPayload)
	}
}
