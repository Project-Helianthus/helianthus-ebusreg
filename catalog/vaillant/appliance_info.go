package vaillant

import (
	_ "embed"
	applianceInfoJSONPackage "encoding/json"
)

//go:embed data/appliance_info.json
var applianceInfoApplianceInfo []byte

// ApplianceInfoCatalog returns a caller-owned copy of the embedded appliance information catalog JSON.
func ApplianceInfoCatalog() applianceInfoJSONPackage.RawMessage {
	return append(applianceInfoJSONPackage.RawMessage(nil), applianceInfoApplianceInfo...)
}
