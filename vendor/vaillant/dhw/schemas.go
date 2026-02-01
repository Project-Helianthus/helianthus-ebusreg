package dhw

import (
	"github.com/d3vi1/helianthus-ebusgo/types"
	"github.com/d3vi1/helianthus-ebusreg/schema"
)

func statusSchemaSelector() schema.SchemaSelector {
	return schema.SchemaSelector{
		Default: schema.Schema{
			Fields: []schema.SchemaField{
				{Name: "status", Type: types.DATA1b{}},
			},
		},
	}
}

func setTargetTempSchemaSelector() schema.SchemaSelector {
	return schema.SchemaSelector{
		Default: schema.Schema{},
	}
}

func parametersSchemaSelector() schema.SchemaSelector {
	return schema.SchemaSelector{
		Default: schema.Schema{
			Fields: []schema.SchemaField{
				{Name: "dhw_temp", Type: types.DATA2b{}},
				{Name: "target_temp", Type: types.DATA2b{}},
				{Name: "mode", Type: types.DATA1b{}},
			},
		},
	}
}
