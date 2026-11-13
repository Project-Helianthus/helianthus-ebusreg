package vaillant

import (
	_ "embed"
	applianceInfoJSONPackage "encoding/json"
)

//go:embed data/ab33_appliance_info.json
var applianceInfoApplianceInfo []byte

func ApplianceInfoCatalog() applianceInfoJSONPackage.RawMessage {
	return append(applianceInfoJSONPackage.RawMessage(nil), applianceInfoApplianceInfo...)
}
