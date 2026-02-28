package schema_test

import (
	"reflect"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusgo/types"
	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
)

func TestSchemaSelector_SelectsByTarget(t *testing.T) {
	t.Parallel()

	defaultSchema := schema.Schema{Fields: []schema.SchemaField{{Name: "default", Type: types.DATA1b{}}}}
	targetSchema := schema.Schema{Fields: []schema.SchemaField{{Name: "target", Type: types.WORD{}}}}

	selector := schema.SchemaSelector{
		Default: defaultSchema,
		Conditions: []schema.SchemaCondition{
			{
				Target:    0x10,
				HasTarget: true,
				Schema:    targetSchema,
			},
		},
	}

	got := selector.Select(0x10, "")
	if !reflect.DeepEqual(got, targetSchema) {
		t.Fatalf("Select target = %#v; want %#v", got, targetSchema)
	}

	got = selector.Select(0x08, "")
	if !reflect.DeepEqual(got, defaultSchema) {
		t.Fatalf("Select default = %#v; want %#v", got, defaultSchema)
	}
}

func TestSchemaSelector_SelectsByHardwareVersion(t *testing.T) {
	t.Parallel()

	defaultSchema := schema.Schema{Fields: []schema.SchemaField{{Name: "default", Type: types.DATA1b{}}}}
	hwSchema := schema.Schema{Fields: []schema.SchemaField{{Name: "hw", Type: types.WORD{}}}}

	selector := schema.SchemaSelector{
		Default: defaultSchema,
		Conditions: []schema.SchemaCondition{
			{
				MinHW:    7603,
				HasMinHW: true,
				Schema:   hwSchema,
			},
		},
	}

	got := selector.Select(0x08, "7603")
	if !reflect.DeepEqual(got, hwSchema) {
		t.Fatalf("Select hw = %#v; want %#v", got, hwSchema)
	}

	got = selector.Select(0x08, "7602")
	if !reflect.DeepEqual(got, defaultSchema) {
		t.Fatalf("Select default = %#v; want %#v", got, defaultSchema)
	}
}

func TestSchemaSelector_SelectsFirstMatch(t *testing.T) {
	t.Parallel()

	defaultSchema := schema.Schema{Fields: []schema.SchemaField{{Name: "default", Type: types.DATA1b{}}}}
	targetSchema := schema.Schema{Fields: []schema.SchemaField{{Name: "target", Type: types.WORD{}}}}
	hwSchema := schema.Schema{Fields: []schema.SchemaField{{Name: "hw", Type: types.BCD{}}}}

	selector := schema.SchemaSelector{
		Default: defaultSchema,
		Conditions: []schema.SchemaCondition{
			{
				Target:    0x10,
				HasTarget: true,
				Schema:    targetSchema,
			},
			{
				MinHW:    7603,
				HasMinHW: true,
				Schema:   hwSchema,
			},
		},
	}

	got := selector.Select(0x10, "8000")
	if !reflect.DeepEqual(got, targetSchema) {
		t.Fatalf("Select first match = %#v; want %#v", got, targetSchema)
	}
}

func TestSchemaSelector_ParsesHardwareVersionDigits(t *testing.T) {
	t.Parallel()

	defaultSchema := schema.Schema{Fields: []schema.SchemaField{{Name: "default", Type: types.DATA1b{}}}}
	hwSchema := schema.Schema{Fields: []schema.SchemaField{{Name: "hw", Type: types.WORD{}}}}

	selector := schema.SchemaSelector{
		Default: defaultSchema,
		Conditions: []schema.SchemaCondition{
			{
				MinHW:    7603,
				HasMinHW: true,
				Schema:   hwSchema,
			},
		},
	}

	got := selector.Select(0x08, "7.603")
	if !reflect.DeepEqual(got, hwSchema) {
		t.Fatalf("Select parsed hw = %#v; want %#v", got, hwSchema)
	}
}
