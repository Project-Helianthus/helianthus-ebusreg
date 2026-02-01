package schema_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/d3vi1/helianthus-ebusreg/schema"
)

func TestLoadSchemaJSON(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/basic_schema.json")
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}

	loaded, err := schema.LoadSchemaJSON(data)
	if err != nil {
		t.Fatalf("LoadSchemaJSON error = %v", err)
	}

	if got := loaded.Size(); got != 5 {
		t.Fatalf("Size = %d; want 5", got)
	}
	if got := fieldNames(loaded); !reflect.DeepEqual(got, []string{"status", "temp", "flags"}) {
		t.Fatalf("Fields = %v; want [status temp flags]", got)
	}
}

func TestLoadSchemaSelectorJSON(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/selector_schema.json")
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}

	selector, err := schema.LoadSchemaSelectorJSON(data)
	if err != nil {
		t.Fatalf("LoadSchemaSelectorJSON error = %v", err)
	}

	got := selector.Select(0x10, "7000")
	if gotNames := fieldNames(got); !reflect.DeepEqual(gotNames, []string{"target"}) {
		t.Fatalf("Select target fields = %v; want [target]", gotNames)
	}

	got = selector.Select(0x08, "7603")
	if gotNames := fieldNames(got); !reflect.DeepEqual(gotNames, []string{"hw"}) {
		t.Fatalf("Select hw fields = %v; want [hw]", gotNames)
	}

	got = selector.Select(0x08, "7000")
	if gotNames := fieldNames(got); !reflect.DeepEqual(gotNames, []string{"default"}) {
		t.Fatalf("Select default fields = %v; want [default]", gotNames)
	}
}

func fieldNames(s schema.Schema) []string {
	names := make([]string, 0, len(s.Fields))
	for _, field := range s.Fields {
		names = append(names, field.Name)
	}
	return names
}
