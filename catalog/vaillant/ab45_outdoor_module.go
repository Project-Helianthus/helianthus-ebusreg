package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/ab45_outdoor_module.json
var ab45OutdoorModule []byte

func Ab45OutdoorModuleCatalog() json.RawMessage {
	return append(json.RawMessage(nil), ab45OutdoorModule...)
}
