package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/ab43_simulator_adc.json
var ab43SimulatorADC []byte

func Ab43SimulatorADC() json.RawMessage {
	return append(json.RawMessage(nil), ab43SimulatorADC...)
}
