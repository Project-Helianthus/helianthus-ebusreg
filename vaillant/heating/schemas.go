package heating

import (
	"github.com/Project-Helianthus/helianthus-ebusgo/types"
	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
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
	boilerSchema := schema.Schema{
		Fields: []schema.SchemaField{
			{Name: "flow_temp", Type: types.DATA2b{}},
			{Name: "return_temp", Type: types.DATA2b{}},
			{Name: "pump_status", Type: types.DATA1b{}},
		},
	}
	controllerSchema := schema.Schema{
		Fields: []schema.SchemaField{
			{Name: "room_temp", Type: types.DATA2b{}},
			{Name: "target_temp", Type: types.DATA2b{}},
			{Name: "mode", Type: types.DATA1b{}},
		},
	}

	return schema.SchemaSelector{
		Default: boilerSchema,
		Conditions: []schema.SchemaCondition{
			{
				Target:    0x10,
				HasTarget: true,
				Schema:    controllerSchema,
			},
		},
	}
}

func energySchemaSelector() schema.SchemaSelector {
	return schema.SchemaSelector{
		Default: schema.Schema{
			Fields: []schema.SchemaField{
				{Name: "energy", Type: types.WORD{}},
			},
		},
	}
}
