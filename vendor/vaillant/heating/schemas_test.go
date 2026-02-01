package heating

import (
	"testing"

	"github.com/d3vi1/helianthus-ebusgo/types"
)

func TestParametersSchemaSelector_BoilerVariant(t *testing.T) {
	t.Parallel()

	selector := parametersSchemaSelector()
	schema := selector.Select(0x08, "")
	payload := []byte{0x00, 0x14, 0x00, 0x12, 0x01}

	values, err := schema.Decode(payload)
	if err != nil {
		t.Fatalf("Decode error = %v", err)
	}

	assertValue(t, values, "flow_temp", 20.0)
	assertValue(t, values, "return_temp", 18.0)
	assertValue(t, values, "pump_status", int8(1))
}

func TestParametersSchemaSelector_ControllerVariant(t *testing.T) {
	t.Parallel()

	selector := parametersSchemaSelector()
	schema := selector.Select(0x10, "")
	payload := []byte{0x80, 0x15, 0x00, 0x16, 0x02}

	values, err := schema.Decode(payload)
	if err != nil {
		t.Fatalf("Decode error = %v", err)
	}

	assertValue(t, values, "room_temp", 21.5)
	assertValue(t, values, "target_temp", 22.0)
	assertValue(t, values, "mode", int8(2))
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
