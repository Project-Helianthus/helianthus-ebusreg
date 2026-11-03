package system

// Ab22RequestContext is caller-supplied request context, not a response echo.
// Ab22Family records a family label only; it does not qualify a codec or device.
type Ab22RequestContext struct {
	Ab22Family string
	Ab22Opcode byte
}

// Ab22StatusFields preserves the eight OP08 bytes without interpreting them.
type Ab22StatusFields struct {
	Ab22BindingZone           byte
	Ab22SpecialFunctionStatus byte
	Ab22HeatingOperatingMode  byte
	Ab22CoolingOperatingMode  byte
	Ab22StatusInfo            byte
	Ab22FrostProtection       byte
	Ab22HeatingTemperatureRaw byte
	Ab22CoolingTemperatureRaw byte
}

// Ab22Status retains owned native bytes and their explicit request context.
// Nil Ab22Fields means incomplete or non-OP08 input; callers must retain useful
// prior state rather than treating absent fields as zero-valued updates.
// Fields are raw structure only for every family, including unknown families.
// Semantic units, enum meanings, presence and control authority remain Unknown.
type Ab22Status struct {
	Ab22Context Ab22RequestContext
	Ab22Payload []byte
	Ab22Fields  *Ab22StatusFields
}

// Ab22DecodeStatus is an explicitly called, pure B524 OP08 raw decoder.
// Completeness qualifies only the eight-byte software shape, not physical support.
func Ab22DecodeStatus(ab22Context Ab22RequestContext, ab22Payload []byte) Ab22Status {
	ab22Result := Ab22Status{
		Ab22Context: ab22Context,
		Ab22Payload: append([]byte(nil), ab22Payload...),
	}
	if ab22Context.Ab22Opcode != 0x08 || len(ab22Payload) != 8 {
		return ab22Result
	}
	ab22Result.Ab22Fields = &Ab22StatusFields{
		Ab22BindingZone:           ab22Payload[0],
		Ab22SpecialFunctionStatus: ab22Payload[1],
		Ab22HeatingOperatingMode:  ab22Payload[2],
		Ab22CoolingOperatingMode:  ab22Payload[3],
		Ab22StatusInfo:            ab22Payload[4],
		Ab22FrostProtection:       ab22Payload[5],
		Ab22HeatingTemperatureRaw: ab22Payload[6],
		Ab22CoolingTemperatureRaw: ab22Payload[7],
	}
	return ab22Result
}
