package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/fan_demand.json
var fanDemandFanDemandJSON string

func FanDemandCatalog() json.RawMessage {
	return json.RawMessage([]byte(fanDemandFanDemandJSON))
}
