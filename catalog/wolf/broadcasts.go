package wolf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// Primary and secondary bytes of the documented Wolf messages.
const (
	PBBurner byte = 0x05
	PBBoiler byte = 0x08
	PBWolf   byte = 0x50

	SBBurnerStatus       byte  = 0x03
	SBModulationActual   byte  = 0x04
	SBHeatRequestStatus  byte  = 0x07
	SBBoilerSetpoint     byte  = 0x00
	SBSolarStatus        byte  = 0x14
	SBSolarPump          byte  = 0x17
	SBSolarPower         byte  = 0x18
	SBServiceRead        byte  = 0x22
	SBServiceWrite       byte  = 0x23
	solarPumpOff         byte  = 0xBC
	solarPumpOn          byte  = 0xBD
	maxPercent           byte  = 100
	boilerSetpointDivide int16 = 16
)

// Errors of the Wolf catalogue.
var (
	ErrWidth     = errors.New("wolf: payload width does not match the documented layout")
	ErrRange     = errors.New("wolf: value outside the documented range")
	ErrReadOnly  = errors.New("wolf: register is read-only in the documented table")
	ErrValueSize = errors.New("wolf: value does not fit the documented type")
)

// Field is one named decoded value. Unit is empty when the page states none.
type Field struct {
	Name  string
	Value any
	Unit  string
}

// Message is the decode of one documented frame. Known is false when the
// PB/SB pair is not a documented Wolf message; Raw keeps the payload.
type Message struct {
	Name   string
	Known  bool
	Fields []Field
	Raw    []byte
}

// heatRequestNames are the ten states of the heat-request machine on the page.
var heatRequestNames = map[byte]string{
	0x00: "aus",
	0x01: "keine",
	0x44: "Reglerstop_stufig",
	0x55: "Brauchwasser",
	0x66: "Brauchwasser_Reglerstop",
	0xAA: "Heizen",
	0xBB: "Brauchwasser_Heizen",
	0xCC: "Emissionskontrolle",
	0xDD: "TUeV",
	0xEE: "Reglerstop",
}

// DecodeFrame dispatches a payload by PB/SB to its documented decoder. A pair
// the page does not document is kept raw.
func DecodeFrame(pb, sb byte, payload []byte) (Message, error) {
	switch {
	case pb == PBBurner && sb == SBBurnerStatus:
		return DecodeBurnerStatus(payload)
	case pb == PBBurner && sb == SBModulationActual:
		return DecodeModulation(payload)
	case pb == PBBurner && sb == SBHeatRequestStatus:
		return DecodeHeatRequest(payload)
	case pb == PBBoiler && sb == SBBoilerSetpoint:
		return DecodeBoilerSetpoint(payload)
	case pb == PBWolf && sb == SBSolarStatus:
		return DecodeSolarStatus(payload)
	case pb == PBWolf && sb == SBSolarPump:
		return DecodeSolarPump(payload)
	case pb == PBWolf && sb == SBSolarPower:
		return DecodeSolarPower(payload)
	}
	return Message{Raw: append([]byte(nil), payload...)}, nil
}

// DecodeBurnerStatus decodes PB 05 SB 03: block number (0x01), burner
// automaton state, modulation (0xFF is the gas-thermal dummy value), and the
// supply temperature bytes, which the page gives no width or unit for.
func DecodeBurnerStatus(payload []byte) (Message, error) {
	if len(payload) < 3 {
		return Message{}, fmt.Errorf("%w: burner status wants at least 3 bytes", ErrWidth)
	}
	return Message{Name: "burner_status", Known: true, Raw: copyOf(payload), Fields: []Field{
		{Name: "block_number", Value: payload[0]},
		{Name: "burner_state", Value: payload[1]},
		{Name: "modulation_raw", Value: payload[2]},
		{Name: "kessel_temperature_raw", Value: copyOf(payload[3:])},
	}}, nil
}

// DecodeModulation decodes PB 05 SB 04: actual and maximum modulation as
// percent0 values. A value above 100 is a decode error by the page's claim.
func DecodeModulation(payload []byte) (Message, error) {
	if len(payload) != 2 {
		return Message{}, fmt.Errorf("%w: modulation wants 2 bytes", ErrWidth)
	}
	for _, b := range payload {
		if b > maxPercent {
			return Message{}, fmt.Errorf("%w: modulation %d", ErrRange, b)
		}
	}
	return Message{Name: "modulation_actual", Known: true, Raw: copyOf(payload), Fields: []Field{
		{Name: "modulation_actual", Value: payload[0], Unit: "%"},
		{Name: "modulation_max", Value: payload[1], Unit: "%"},
	}}, nil
}

// DecodeHeatRequest decodes PB 05 SB 07: the heat-request state (named when it
// is one of the ten documented values), the boiler setpoint byte, and the
// remaining pressure bytes raw.
func DecodeHeatRequest(payload []byte) (Message, error) {
	if len(payload) < 2 {
		return Message{}, fmt.Errorf("%w: heat request wants at least 2 bytes", ErrWidth)
	}
	name, named := heatRequestNames[payload[0]]
	state := Field{Name: "heat_request_state", Value: payload[0]}
	if named {
		state.Value = name
	}
	return Message{Name: "heat_request_status", Known: true, Raw: copyOf(payload), Fields: []Field{
		state,
		{Name: "boiler_setpoint_raw", Value: payload[1]},
		{Name: "target_pressure_raw", Value: copyOf(payload[2:])},
	}}, nil
}

// DecodeBoilerSetpoint decodes PB 08 SB 00: the boiler setpoint as temp2 (a
// signed 16-bit little-endian value divided by 16, degrees C), the forced-output
// byte and the status byte.
func DecodeBoilerSetpoint(payload []byte) (Message, error) {
	if len(payload) != 4 {
		return Message{}, fmt.Errorf("%w: boiler setpoint wants 4 bytes", ErrWidth)
	}
	setpoint := float64(int16(binary.LittleEndian.Uint16(payload[0:2]))) / float64(boilerSetpointDivide)
	return Message{Name: "boiler_setpoint_transfer", Known: true, Raw: copyOf(payload), Fields: []Field{
		{Name: "boiler_setpoint", Value: setpoint, Unit: "degC"},
		{Name: "forced_output_raw", Value: payload[2]},
		{Name: "status_raw", Value: payload[3]},
	}}, nil
}

// DecodeSolarStatus decodes PB 50 SB 14: the solar status byte and the ignored
// byte; the mixer setpoint and room temperature bytes are kept raw.
func DecodeSolarStatus(payload []byte) (Message, error) {
	if len(payload) < 2 {
		return Message{}, fmt.Errorf("%w: solar status wants at least 2 bytes", ErrWidth)
	}
	return Message{Name: "solar_status", Known: true, Raw: copyOf(payload), Fields: []Field{
		{Name: "solar_status_raw", Value: payload[0]},
		{Name: "ignored_byte", Value: payload[1]},
		{Name: "mixer_setpoint_raw", Value: copyOf(payload[2:])},
	}}, nil
}

// DecodeSolarPump decodes PB 50 SB 17: the pump state, where 0xBC is off and
// 0xBD is on and any other byte is kept raw as unknown.
func DecodeSolarPump(payload []byte) (Message, error) {
	if len(payload) < 2 {
		return Message{}, fmt.Errorf("%w: solar pump wants at least 2 bytes", ErrWidth)
	}
	pump := Field{Name: "pump_state", Value: payload[0]}
	switch payload[0] {
	case solarPumpOff:
		pump.Value = "off"
	case solarPumpOn:
		pump.Value = "on"
	}
	return Message{Name: "solar_pump", Known: true, Raw: copyOf(payload), Fields: []Field{
		pump,
		{Name: "ignored_byte", Value: payload[1]},
		{Name: "temperatures_raw", Value: copyOf(payload[2:])},
	}}, nil
}

// DecodeSolarPower decodes PB 50 SB 18: solar power as D2B (a signed 16-bit
// value divided by 256, watts) and the two yield bytes.
func DecodeSolarPower(payload []byte) (Message, error) {
	if len(payload) != 4 {
		return Message{}, fmt.Errorf("%w: solar power wants 4 bytes", ErrWidth)
	}
	power := float64(int16(binary.LittleEndian.Uint16(payload[0:2]))) / 256
	return Message{Name: "solar_power", Known: true, Raw: copyOf(payload), Fields: []Field{
		{Name: "solar_power", Value: power, Unit: "W"},
		{Name: "yield_high_raw", Value: payload[2]},
		{Name: "yield_sum_low_raw", Value: payload[3]},
	}}, nil
}

// EncodeBurnerStatus builds the PB 05 SB 03 payload: block 0x01, the burner
// state, the modulation byte and the supply temperature bytes.
func EncodeBurnerStatus(state, modulation byte, kesselTemperature []byte) []byte {
	return append([]byte{0x01, state, modulation}, kesselTemperature...)
}

// EncodeModulation builds the PB 05 SB 04 payload with the two percent values.
func EncodeModulation(actual, maximum byte) ([]byte, error) {
	if actual > maxPercent || maximum > maxPercent {
		return nil, ErrRange
	}
	return []byte{actual, maximum}, nil
}

// EncodeHeatRequest builds the PB 05 SB 07 payload: the state, the boiler
// setpoint byte and the pressure bytes.
func EncodeHeatRequest(state, setpoint byte, pressure []byte) []byte {
	return append([]byte{state, setpoint}, pressure...)
}

// EncodeBoilerSetpoint builds the PB 08 SB 00 payload: the setpoint as temp2
// little-endian, the forced-output byte and the status byte.
func EncodeBoilerSetpoint(celsius float64, forcedOutput, status byte) ([]byte, error) {
	raw := math.Round(celsius * float64(boilerSetpointDivide))
	if raw < math.MinInt16 || raw > math.MaxInt16 {
		return nil, ErrValueSize
	}
	out := make([]byte, 4)
	binary.LittleEndian.PutUint16(out[0:2], uint16(int16(raw)))
	out[2] = forcedOutput
	out[3] = status
	return out, nil
}

// EncodeSolarStatus builds the PB 50 SB 14 payload: the status byte, the ignored
// byte (caller-supplied, since the page gives no value) and the rest.
func EncodeSolarStatus(status, ignored byte, rest []byte) []byte {
	return append([]byte{status, ignored}, rest...)
}

// EncodeSolarPump builds the PB 50 SB 17 payload: the pump state byte, the
// ignored byte and the temperature bytes.
func EncodeSolarPump(on bool, ignored byte, temperatures []byte) []byte {
	state := solarPumpOff
	if on {
		state = solarPumpOn
	}
	return append([]byte{state, ignored}, temperatures...)
}

// EncodeSolarPower builds the PB 50 SB 18 payload from watts as D2B.
func EncodeSolarPower(watts float64, yieldHigh, yieldSumLow byte) ([]byte, error) {
	raw := math.Round(watts * 256)
	if raw < math.MinInt16 || raw > math.MaxInt16 {
		return nil, ErrValueSize
	}
	out := make([]byte, 4)
	binary.LittleEndian.PutUint16(out[0:2], uint16(int16(raw)))
	out[2] = yieldHigh
	out[3] = yieldSumLow
	return out, nil
}

func copyOf(b []byte) []byte {
	return append([]byte(nil), b...)
}

// ServiceRegister is one entry of the 5022/5023 table. ID holds the first two
// bytes of the read ID (the write ID) and the parameter index as the third
// byte. Scale is 10 for UIN/10 and 1 for UIN. Writable reflects the table's R+W mode.
type ServiceRegister struct {
	ID       [3]byte
	Writable bool
	Name     string
	Scale    int
	Notes    string
}

// ServiceRegisterByID returns the table entry for an ID.
func ServiceRegisterByID(id [3]byte) (ServiceRegister, bool) {
	for _, reg := range serviceRegisterTable {
		if reg.ID == id {
			return reg, true
		}
	}
	return ServiceRegister{}, false
}

// EncodeServiceRead builds the 5022 request: id_lo, id_hi, param_idx.
func EncodeServiceRead(reg ServiceRegister) []byte {
	return append([]byte(nil), reg.ID[:]...)
}

// DecodeServiceRead decodes the 5022 reply, a two-byte little-endian UIN value
// divided by the register scale.
func DecodeServiceRead(reg ServiceRegister, payload []byte) (Field, error) {
	if len(payload) != 2 {
		return Field{}, fmt.Errorf("%w: service read wants 2 bytes", ErrWidth)
	}
	raw := binary.LittleEndian.Uint16(payload)
	return Field{Name: reg.Name, Value: float64(raw) / float64(reg.Scale)}, nil
}

// EncodeServiceWrite builds the 5023 request: the write ID (the first two ID
// bytes), the UIN wire value (already scaled, little-endian) and the four
// trailing HEX bytes, which the page states without a value and the caller
// supplies. A register the table marks read-only is refused.
func EncodeServiceWrite(reg ServiceRegister, raw uint16, trailer [4]byte) ([]byte, error) {
	if !reg.Writable {
		return nil, ErrReadOnly
	}
	out := make([]byte, 8)
	out[0] = reg.ID[0]
	out[1] = reg.ID[1]
	binary.LittleEndian.PutUint16(out[2:4], raw)
	copy(out[4:], trailer[:])
	return out, nil
}

// serviceRegisterTable is the 5022/5023 register table transcribed from the Wolf page.
var serviceRegisterTable = []ServiceRegister{
	{ID: [3]byte{0x84, 0x22, 0x00}, Writable: true, Name: "hg01 hysterese_vorlauf", Scale: 10, Notes: "5--30 K -- flow hysteresis"},
	{ID: [3]byte{0x29, 0x5a, 0x01}, Writable: true, Name: "hg02 gebl_unten", Scale: 1, Notes: "25--100% -- blower minimum RPM"},
	{ID: [3]byte{0xcd, 0x59, 0x01}, Writable: false, Name: "hg03 gebl_oben_ww", Scale: 1, Notes: "25--100% -- blower max DHW"},
	{ID: [3]byte{0x1d, 0x3f, 0x01}, Writable: false, Name: "hg04 gebl_oben_hz", Scale: 1, Notes: "1--100% -- blower max heating"},
	{ID: [3]byte{0x25, 0x41, 0x01}, Writable: true, Name: "hg06 betriebsart", Scale: 1, Notes: "0--2 -- pump mode"},
	{ID: [3]byte{0xc1, 0x42, 0x01}, Writable: true, Name: "hg07 pumpen_nachlauf", Scale: 1, Notes: "0--30 min -- pump rundown"},
	{ID: [3]byte{0xde, 0x84, 0x02}, Writable: true, Name: "hg08 max_vorlauf_temp", Scale: 10, Notes: "40--90 degrees C -- max flow setpoint"},
	{ID: [3]byte{0x9d, 0x43, 0x01}, Writable: true, Name: "hg09 taktsperre", Scale: 1, Notes: "1--30 min -- anti-cycling time"},
	{ID: [3]byte{0xad, 0x78, 0x01}, Writable: false, Name: "hg10 ebus_addr", Scale: 1, Notes: "eBUS bus address"},
	{ID: [3]byte{0x79, 0x40, 0x01}, Writable: true, Name: "hg15 hysterese_speicher", Scale: 10, Notes: "1--30 K -- storage hysteresis"},
	{ID: [3]byte{0xb9, 0x55, 0x01}, Writable: true, Name: "hg16 pumpenleistung_hk_min", Scale: 1, Notes: "20--100% -- HC pump min"},
	{ID: [3]byte{0x5d, 0x56, 0x01}, Writable: true, Name: "hg17 pumpenleistung_hk_max", Scale: 1, Notes: "20--100% -- HC pump max"},
	{ID: [3]byte{0x20, 0x1f, 0x00}, Writable: true, Name: "hg21 kessel_min", Scale: 10, Notes: "Boiler minimum temperature"},
	{ID: [3]byte{0xf4, 0x27, 0x00}, Writable: true, Name: "hg22 kessel_max", Scale: 10, Notes: "Boiler maximum temperature"},
	{ID: [3]byte{0xf5, 0x70, 0x01}, Writable: false, Name: "hg73 io_istwert", Scale: 1, Notes: "IO current value"},
	{ID: [3]byte{0xd5, 0xf6, 0x01}, Writable: false, Name: "hg74 gebl_drehzahl", Scale: 1, Notes: "Blower RPM actual"},
	{ID: [3]byte{0x31, 0x6c, 0x01}, Writable: false, Name: "hg75 ww_durchsatz", Scale: 1, Notes: "DHW flow rate"},
	{ID: [3]byte{0x66, 0x28, 0x02}, Writable: false, Name: "hours_hg9x betriebsstunden_netz", Scale: 1, Notes: "Mains operating hours"},
	{ID: [3]byte{0xde, 0x2a, 0x02}, Writable: false, Name: "hours_hg90 betriebsstunden_brenner", Scale: 1, Notes: "Burner operating hours"},
	{ID: [3]byte{0xaa, 0x26, 0x02}, Writable: false, Name: "starts_hg91 brennerstarts", Scale: 1, Notes: "Burner start count"},
	{ID: [3]byte{0xb8, 0x02, 0x00}, Writable: false, Name: "temp_vorlauf_soll", Scale: 10, Notes: "Flow temperature setpoint"},
	{ID: [3]byte{0x28, 0x0d, 0x00}, Writable: false, Name: "temp_vorlauf_ist", Scale: 10, Notes: "Flow temperature actual"},
	{ID: [3]byte{0x24, 0x16, 0x00}, Writable: false, Name: "temp_ruecklauf_ist", Scale: 10, Notes: "Return temperature actual"},
	{ID: [3]byte{0xe4, 0x03, 0x00}, Writable: false, Name: "temp_warmwasser_soll", Scale: 10, Notes: "DHW setpoint"},
	{ID: [3]byte{0xcc, 0x0e, 0x00}, Writable: false, Name: "temp_warmwasser_ist", Scale: 10, Notes: "DHW temperature actual"},
	{ID: [3]byte{0x6d, 0x6d, 0x01}, Writable: false, Name: "status_pwm_pumpe", Scale: 1, Notes: "PWM pump status"},
}
