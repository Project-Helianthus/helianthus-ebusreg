package solar

import (
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/types"
)

func TestParametersSchemaDecode(t *testing.T) {
	t.Parallel()

	selector := parametersSchemaSelector()
	schemaValue := selector.Select(0x08, "")
	payload := []byte{0x03, 0x80, 0x06}

	values, err := schemaValue.Decode(payload)
	if err != nil {
		t.Fatalf("Decode error = %v", err)
	}

	assertValue(t, values, "pump_speed", int8(3))
	assertValue(t, values, "delta_temp", 6.5)
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
