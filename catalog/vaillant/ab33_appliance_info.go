package vaillant

import (
	_ "embed"
	ab33json "encoding/json"
)

//go:embed data/ab33_appliance_info.json
var ab33ApplianceInfo []byte

func Ab33ApplianceInfoCatalog() ab33json.RawMessage {
	return append(ab33json.RawMessage(nil), ab33ApplianceInfo...)
}
