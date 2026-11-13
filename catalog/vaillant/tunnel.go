package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/ab37_tunnel.json
var tunnelTunnel []byte

func TunnelCatalog() json.RawMessage {
	return append(json.RawMessage(nil), tunnelTunnel...)
}
