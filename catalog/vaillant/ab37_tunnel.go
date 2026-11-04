package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/ab37_tunnel.json
var ab37Tunnel []byte

func Ab37TunnelCatalog() json.RawMessage {
	return append(json.RawMessage(nil), ab37Tunnel...)
}
