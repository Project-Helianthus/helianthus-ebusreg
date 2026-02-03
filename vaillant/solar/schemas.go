package solar

import (
	"github.com/d3vi1/helianthus-ebusgo/types"
	"github.com/d3vi1/helianthus-ebusreg/schema"
)

func statusSchemaSelector() schema.SchemaSelector {
	return schema.SchemaSelector{
		Default: schema.Schema{
			Fields: []schema.SchemaField{
				{Name: "collector_temp", Type: types.DATA2b{}},
				{Name: "tank_temp", Type: types.DATA2b{}},
			},
		},
	}
}

func solarYieldSchemaSelector() schema.SchemaSelector {
	return schema.SchemaSelector{
		Default: schema.Schema{
			Fields: []schema.SchemaField{
				{Name: "yield", Type: types.WORD{}},
			},
		},
	}
}

func parametersSchemaSelector() schema.SchemaSelector {
	return schema.SchemaSelector{
		Default: schema.Schema{
			Fields: []schema.SchemaField{
				{Name: "pump_speed", Type: types.DATA1b{}},
				{Name: "delta_temp", Type: types.DATA2b{}},
			},
		},
	}
}
