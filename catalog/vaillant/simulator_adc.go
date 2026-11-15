package vaillant

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/simulator_adc.json
var simulatorADCSimulatorADC []byte

// SimulatorADC returns a caller-owned copy of the embedded simulator ADC catalog JSON.
func SimulatorADC() json.RawMessage {
	return append(json.RawMessage(nil), simulatorADCSimulatorADC...)
}
