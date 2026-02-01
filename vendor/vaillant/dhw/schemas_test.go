package dhw

import (
	"testing"

	"github.com/d3vi1/helianthus-ebusgo/types"
)

func TestParametersSchemaDecode(t *testing.T) {
	t.Parallel()

	selector := parametersSchemaSelector()
	schemaValue := selector.Select(0x08, "")
	payload := []byte{0x00, 0x2D, 0x00, 0x32, 0x01}

	values, err := schemaValue.Decode(payload)
	if err != nil {
		t.Fatalf("Decode error = %v", err)
	}

	assertValue(t, values, "dhw_temp", 45.0)
	assertValue(t, values, "target_temp", 50.0)
	assertValue(t, values, "mode", int8(1))
}

func assertValue(t *testing.T, values map[string]types.Value, key string, want any) {
	t.Helper()

	value, ok := values[key]
	if !ok {
		t.Fatalf("missing value %s", key)
	}
	if !value.Valid {
		t.Fatalf("value %s invalid", key)
	}
	if value.Value != want {
		t.Fatalf("value %s = %v; want %v", key, value.Value, want)
	}
}
