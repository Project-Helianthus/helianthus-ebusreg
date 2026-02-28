package schema_test

import (
	"errors"
	"reflect"
	"testing"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/types"
	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
)

func TestSchemaDecode(t *testing.T) {
	t.Parallel()

	s := schema.Schema{
		Fields: []schema.SchemaField{
			{Name: "first", Type: types.DATA1b{}},
			{Name: "second", Type: types.WORD{}},
		},
	}

	got, err := s.Decode([]byte{0x01, 0x34, 0x12})
	if err != nil {
		t.Fatalf("Decode error = %v", err)
	}

	want := map[string]types.Value{
		"first":  {Value: int8(1), Valid: true},
		"second": {Value: uint16(0x1234), Valid: true},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Decode = %#v; want %#v", got, want)
	}
}

func TestSchemaEncode(t *testing.T) {
	t.Parallel()

	s := schema.Schema{
		Fields: []schema.SchemaField{
			{Name: "first", Type: types.DATA1b{}},
			{Name: "second", Type: types.WORD{}},
		},
	}

	got, err := s.Encode(map[string]any{
		"first":  int8(1),
		"second": uint16(0x1234),
	})
	if err != nil {
		t.Fatalf("Encode error = %v", err)
	}

	want := []byte{0x01, 0x34, 0x12}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Encode = %#v; want %#v", got, want)
	}
}

func TestSchemaDecodeShortPayload(t *testing.T) {
	t.Parallel()

	s := schema.Schema{
		Fields: []schema.SchemaField{
			{Name: "first", Type: types.DATA1b{}},
			{Name: "second", Type: types.WORD{}},
		},
	}

	_, err := s.Decode([]byte{0x01})
	assertInvalidPayload(t, err)
}

func TestSchemaEncodeMissingField(t *testing.T) {
	t.Parallel()

	s := schema.Schema{
		Fields: []schema.SchemaField{
			{Name: "first", Type: types.DATA1b{}},
			{Name: "second", Type: types.WORD{}},
		},
	}

	_, err := s.Encode(map[string]any{
		"first": int8(1),
	})
	assertInvalidPayload(t, err)
}

func TestSchemaEncodeInvalidValue(t *testing.T) {
	t.Parallel()

	s := schema.Schema{
		Fields: []schema.SchemaField{
			{Name: "first", Type: types.DATA1b{}},
		},
	}

	_, err := s.Encode(map[string]any{
		"first": "not-a-number",
	})
	assertInvalidPayload(t, err)
}

func TestSchemaDecodeMissingType(t *testing.T) {
	t.Parallel()

	s := schema.Schema{
		Fields: []schema.SchemaField{
			{Name: "first"},
		},
	}

	_, err := s.Decode([]byte{0x01})
	assertInvalidPayload(t, err)
}

func TestSchemaEncodeMissingType(t *testing.T) {
	t.Parallel()

	s := schema.Schema{
		Fields: []schema.SchemaField{
			{Name: "first"},
		},
	}

	_, err := s.Encode(map[string]any{"first": int8(1)})
	assertInvalidPayload(t, err)
}

func TestSchemaSize(t *testing.T) {
	t.Parallel()

	s := schema.Schema{
		Fields: []schema.SchemaField{
			{Name: "first", Type: types.DATA1b{}},
			{Name: "second", Type: types.WORD{}},
		},
	}

	if got := s.Size(); got != 3 {
		t.Fatalf("Size = %d; want 3", got)
	}
}

func assertInvalidPayload(t *testing.T, err error) {
	t.Helper()

	if !errors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("error = %v; want ErrInvalidPayload", err)
	}
}
