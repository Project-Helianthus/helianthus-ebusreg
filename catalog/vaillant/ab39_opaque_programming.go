package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/ab39_opaque_programming.json
var ab39OpaqueProgramming []byte

func Ab39OpaqueProgramming() json.RawMessage {
	return append(json.RawMessage(nil), ab39OpaqueProgramming...)
}
