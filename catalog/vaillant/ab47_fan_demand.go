package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/ab47_fan_demand.json
var ab47FanDemandJSON string

func Ab47FanDemandCatalog() json.RawMessage {
	return json.RawMessage([]byte(ab47FanDemandJSON))
}
