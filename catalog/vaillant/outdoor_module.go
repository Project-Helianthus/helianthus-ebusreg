package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/ab45_outdoor_module.json
var outdoorModuleOutdoorModule []byte

func OutdoorModuleCatalog() json.RawMessage {
	return append(json.RawMessage(nil), outdoorModuleOutdoorModule...)
}
