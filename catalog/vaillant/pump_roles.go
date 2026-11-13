package vaillant

import (
	_ "embed"
	pumpRoleJSON "encoding/json"
)

//go:embed data/ab31_pump_roles.json
var pumpRolePumpRoles []byte

func PumpRoles() pumpRoleJSON.RawMessage {
	return append(pumpRoleJSON.RawMessage(nil), pumpRolePumpRoles...)
}
