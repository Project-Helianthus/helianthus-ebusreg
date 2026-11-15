package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/outdoor_module.json
var outdoorModuleOutdoorModule []byte

// OutdoorModuleCatalog returns a caller-owned copy of the embedded outdoor module catalog JSON.
func OutdoorModuleCatalog() json.RawMessage {
	return append(json.RawMessage(nil), outdoorModuleOutdoorModule...)
}
