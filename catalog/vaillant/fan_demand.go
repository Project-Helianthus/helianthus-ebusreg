package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/fan_demand.json
var fanDemandFanDemandJSON string

// FanDemandCatalog returns a caller-owned copy of the embedded fan demand catalog JSON.
func FanDemandCatalog() json.RawMessage {
	return json.RawMessage([]byte(fanDemandFanDemandJSON))
}
