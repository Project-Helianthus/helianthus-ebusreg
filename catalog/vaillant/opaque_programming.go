package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/opaque_programming.json
var opaqueProgrammingOpaqueProgramming []byte

func OpaqueProgramming() json.RawMessage {
	return append(json.RawMessage(nil), opaqueProgrammingOpaqueProgramming...)
}
