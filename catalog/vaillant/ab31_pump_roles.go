package vaillant

import (
	_ "embed"
	ab31json "encoding/json"
)

//go:embed data/ab31_pump_roles.json
var ab31PumpRoles []byte

func Ab31PumpRoles() ab31json.RawMessage {
	return append(ab31json.RawMessage(nil), ab31PumpRoles...)
}
