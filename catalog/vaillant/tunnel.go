package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/tunnel.json
var tunnelTunnel []byte

// TunnelCatalog returns a caller-owned copy of the embedded tunnel catalog JSON.
func TunnelCatalog() json.RawMessage {
	return append(json.RawMessage(nil), tunnelTunnel...)
}
