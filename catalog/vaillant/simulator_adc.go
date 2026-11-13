package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/ab43_simulator_adc.json
var simulatorADCSimulatorADC []byte

func SimulatorADC() json.RawMessage {
	return append(json.RawMessage(nil), simulatorADCSimulatorADC...)
}
