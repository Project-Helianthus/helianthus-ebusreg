package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/simulator_adc.json
var simulatorADCSimulatorADC []byte

func SimulatorADC() json.RawMessage {
	return append(json.RawMessage(nil), simulatorADCSimulatorADC...)
}
