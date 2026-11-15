package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/opaque_programming.json
var opaqueProgrammingOpaqueProgramming []byte

// OpaqueProgramming returns a caller-owned copy of the embedded opaque programming catalog JSON.
func OpaqueProgramming() json.RawMessage {
	return append(json.RawMessage(nil), opaqueProgrammingOpaqueProgramming...)
}
