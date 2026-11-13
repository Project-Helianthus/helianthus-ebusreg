package system

// StatusRequestContext is caller-supplied request context, not a response echo.
// StatusFamily records a family label only; it does not qualify a codec or device.
type StatusRequestContext struct {
	StatusFamily string `json:"Ab22Family"`
	StatusOpcode byte   `json:"Ab22Opcode"`
}

// StatusStatusFields preserves the eight OP08 bytes without interpreting them.
type StatusStatusFields struct {
	StatusBindingZone           byte `json:"Ab22BindingZone"`
	SpecialFunctionStatus       byte `json:"Ab22SpecialFunctionStatus"`
	StatusHeatingOperatingMode  byte `json:"Ab22HeatingOperatingMode"`
	StatusCoolingOperatingMode  byte `json:"Ab22CoolingOperatingMode"`
	StatusInfo                  byte `json:"Ab22StatusInfo"`
	StatusFrostProtection       byte `json:"Ab22FrostProtection"`
	StatusHeatingTemperatureRaw byte `json:"Ab22HeatingTemperatureRaw"`
	StatusCoolingTemperatureRaw byte `json:"Ab22CoolingTemperatureRaw"`
}

// StatusStatus retains owned native bytes and their explicit request context.
// Nil StatusFields means incomplete or non-OP08 input; callers must retain useful
// prior state rather than treating absent fields as zero-valued updates.
// Fields are raw structure only for every family, including unknown families.
// Semantic units, enum meanings, presence and control authority remain Unknown.
type StatusStatus struct {
	StatusContext StatusRequestContext `json:"Ab22Context"`
	StatusPayload []byte               `json:"Ab22Payload"`
	StatusFields  *StatusStatusFields  `json:"Ab22Fields"`
}

// DecodeStatus is an explicitly called, pure B524 OP08 raw decoder.
// Completeness qualifies only the eight-byte software shape, not physical support.
func DecodeStatus(statusContext StatusRequestContext, statusPayload []byte) StatusStatus {
	statusResult := StatusStatus{
		StatusContext: statusContext,
		StatusPayload: append([]byte(nil), statusPayload...),
	}
	if statusContext.StatusOpcode != 0x08 || len(statusPayload) != 8 {
		return statusResult
	}
	statusResult.StatusFields = &StatusStatusFields{
		StatusBindingZone:           statusPayload[0],
		SpecialFunctionStatus:       statusPayload[1],
		StatusHeatingOperatingMode:  statusPayload[2],
		StatusCoolingOperatingMode:  statusPayload[3],
		StatusInfo:                  statusPayload[4],
		StatusFrostProtection:       statusPayload[5],
		StatusHeatingTemperatureRaw: statusPayload[6],
		StatusCoolingTemperatureRaw: statusPayload[7],
	}
	return statusResult
}
