package vaillant

import (
	_ "embed"
	pumpRoleJSON "encoding/json"
)

//go:embed data/pump_roles.json
var pumpRolePumpRoles []byte

// PumpRoles returns a caller-owned copy of the embedded pump role catalog JSON.
func PumpRoles() pumpRoleJSON.RawMessage {
	return append(pumpRoleJSON.RawMessage(nil), pumpRolePumpRoles...)
}
