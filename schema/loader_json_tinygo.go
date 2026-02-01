//go:build tinygo

package schema

import (
	"fmt"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
)

// LoadSchemaJSON is not available on TinyGo targets.
func LoadSchemaJSON(data []byte) (Schema, error) {
	return Schema{}, fmt.Errorf("schema json unavailable: %w", ebuserrors.ErrInvalidPayload)
}

// LoadSchemaFile is not available on TinyGo targets.
func LoadSchemaFile(path string) (Schema, error) {
	return Schema{}, fmt.Errorf("schema file unavailable: %w", ebuserrors.ErrInvalidPayload)
}

// LoadSchemaSelectorJSON is not available on TinyGo targets.
func LoadSchemaSelectorJSON(data []byte) (SchemaSelector, error) {
	return SchemaSelector{}, fmt.Errorf("schema selector json unavailable: %w", ebuserrors.ErrInvalidPayload)
}

// LoadSchemaSelectorFile is not available on TinyGo targets.
func LoadSchemaSelectorFile(path string) (SchemaSelector, error) {
	return SchemaSelector{}, fmt.Errorf("schema selector file unavailable: %w", ebuserrors.ErrInvalidPayload)
}
