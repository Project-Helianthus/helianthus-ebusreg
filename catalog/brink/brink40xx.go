package brink

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// Process byte and sub-protocol bytes of the 40xx family.
const (
	PB                   byte = 0x40
	SBFactoryReset       byte = 0xFF
	SBResetNotifications byte = 0x91
	SBErrorHistory       byte = 0x90
	SBFanModeTemporary   byte = 0xA1
	SBFanModePersistent  byte = 0xCB
	SBSensorStatus       byte = 0x22
	SBParameterRead      byte = 0x50
	SBParameterWrite     byte = 0x80
)

// DefaultAddress is the device address the page states. The address can be
// changed with parameter 0xD0, so it is not a family key.
const DefaultAddress byte = 0x3C

// ErrUnknownID reports an ID the page does not define for the sub-protocol.
var ErrUnknownID = errors.New("brink 40xx: unknown ID")

// ErrWidth reports a payload whose length does not match the declared wire type.
var ErrWidth = errors.New("brink 40xx: payload width does not match the declared type")

// ErrValueWidth reports a write value that does not fit the declared wire type.
var ErrValueWidth = errors.New("brink 40xx: value does not fit the declared wire type")

// Kind is the wire type declared by the page.
type Kind int

// Wire types of the family.
const (
	KindUCH Kind = iota
	KindUIR
	KindSIR
	KindULR
	KindSTR
)

func (k Kind) width() int {
	switch k {
	case KindUCH:
		return 1
	case KindUIR, KindSIR:
		return 2
	case KindULR:
		return 4
	}
	return 0
}

// Register is one 4022 sensor or status register. Divisor is the post-decode
// scale (physical = wire / Divisor); 0 means the page states no divisor.
// Length is the byte count of an STR register.
type Register struct {
	ID       byte
	Name     string
	Kind     Kind
	Length   int
	Divisor  int
	Unit     string
	Clusters string
}

// Parameter is one 4050/4080 configurable parameter. The range and default are
// physical values. RangeKnown is false when the page states no numeric range.
type Parameter struct {
	ID         byte
	Name       string
	Kind       Kind
	Divisor    int
	Unit       string
	RangeKnown bool
	Min        float64
	Max        float64
	Default    float64
}

// Reading is one decoded 4022 reply. Known is false for an ID the page does
// not define; Raw always keeps the payload. Value is set for numeric kinds and
// Text for STR kinds.
type Reading struct {
	ID    byte
	Name  string
	Known bool
	Raw   []byte
	Value float64
	Text  string
	Unit  string
}

// ParameterReading is one decoded 4050 reply: the five fields of the page,
// current, minimum, maximum, step and default, in physical values.
type ParameterReading struct {
	ID      byte
	Name    string
	Raw     []byte
	Current float64
	Min     float64
	Max     float64
	Step    float64
	Default float64
	Unit    string
}

// SensorRegister returns the 4022 register with the given ID.
func SensorRegister(id byte) (Register, bool) {
	for _, reg := range sensorRegisterTable {
		if reg.ID == id {
			return reg, true
		}
	}
	return Register{}, false
}

// ParameterByID returns the 4050/4080 parameter with the given ID.
func ParameterByID(id byte) (Parameter, bool) {
	for _, p := range parameterTable {
		if p.ID == id {
			return p, true
		}
	}
	return Parameter{}, false
}

// EncodeSensorRequest returns the 4022 request payload for ID: one byte (NN = 0x01).
func EncodeSensorRequest(id byte) []byte {
	return []byte{id}
}

// DecodeSensorReply decodes a 4022 reply for ID. An ID the page does not define
// returns the raw payload with Known false and no error.
func DecodeSensorReply(id byte, payload []byte) (Reading, error) {
	reg, ok := SensorRegister(id)
	if !ok {
		return Reading{ID: id, Raw: append([]byte(nil), payload...)}, nil
	}
	out := Reading{ID: id, Name: reg.Name, Known: true, Raw: append([]byte(nil), payload...), Unit: reg.Unit}
	if reg.Kind == KindSTR {
		if len(payload) != reg.Length {
			return out, fmt.Errorf("%w: ID 0x%02X wants %d bytes, got %d", ErrWidth, id, reg.Length, len(payload))
		}
		out.Text = string(payload)
		return out, nil
	}
	if len(payload) != reg.Kind.width() {
		return out, fmt.Errorf("%w: ID 0x%02X wants %d bytes, got %d", ErrWidth, id, reg.Kind.width(), len(payload))
	}
	value, err := decodeNumeric(reg.Kind, payload)
	if err != nil {
		return out, err
	}
	out.Value = scale(value, reg.Divisor)
	return out, nil
}

// EncodeParameterRead returns the 4050 request payload for ID: one byte (NN = 0x01).
func EncodeParameterRead(id byte) []byte {
	return []byte{id}
}

// DecodeParameterReply decodes a 4050 reply. The reply is five fields of the
// parameter's width, so its length is 5 x width. An ID the page does not
// define is an error, because the field layout is then unknown.
func DecodeParameterReply(id byte, payload []byte) (ParameterReading, error) {
	p, ok := ParameterByID(id)
	if !ok {
		return ParameterReading{ID: id, Raw: append([]byte(nil), payload...)}, fmt.Errorf("%w: parameter 0x%02X", ErrUnknownID, id)
	}
	w := p.Kind.width()
	if len(payload) != 5*w {
		return ParameterReading{}, fmt.Errorf("%w: parameter 0x%02X wants %d bytes, got %d", ErrWidth, id, 5*w, len(payload))
	}
	fields := make([]float64, 5)
	for i := range fields {
		v, err := decodeNumeric(p.Kind, payload[i*w:(i+1)*w])
		if err != nil {
			return ParameterReading{}, err
		}
		fields[i] = scale(v, p.Divisor)
	}
	return ParameterReading{
		ID: id, Name: p.Name, Raw: append([]byte(nil), payload...), Unit: p.Unit,
		Current: fields[0], Min: fields[1], Max: fields[2], Step: fields[3], Default: fields[4],
	}, nil
}

// EncodeParameterWrite returns the 4080 request payload: the ID and the wire
// value of the parameter's width, little-endian. raw is the wire integer
// (physical x Divisor). The page states that the wire enforces no range, so none
// is checked here; only the width is.
func EncodeParameterWrite(id byte, raw int64) ([]byte, error) {
	p, ok := ParameterByID(id)
	if !ok {
		return nil, fmt.Errorf("%w: parameter 0x%02X", ErrUnknownID, id)
	}
	w := p.Kind.width()
	out := make([]byte, 1+w)
	out[0] = id
	switch p.Kind {
	case KindUCH:
		if raw < 0 || raw > math.MaxUint8 {
			return nil, ErrValueWidth
		}
		out[1] = byte(raw)
	case KindUIR:
		if raw < 0 || raw > math.MaxUint16 {
			return nil, ErrValueWidth
		}
		binary.LittleEndian.PutUint16(out[1:], uint16(raw))
	case KindSIR:
		if raw < math.MinInt16 || raw > math.MaxInt16 {
			return nil, ErrValueWidth
		}
		binary.LittleEndian.PutUint16(out[1:], uint16(int16(raw)))
	default:
		return nil, ErrValueWidth
	}
	return out, nil
}

// EncodeFactoryReset returns the 0x40FF request: no data (NN = 0x00) and the
// reset is triggered on the device. It is sent only when called.
func EncodeFactoryReset() []byte {
	return []byte{}
}

// EncodeResetNotifications returns the 0x4091 request: ID 0x00 (NN = 0x01).
func EncodeResetNotifications() []byte {
	return []byte{0x00}
}

// EncodeErrorHistoryRequest returns the 0x4090 request: ID 0x00 (NN = 0x01).
// The reply layout is not documented on the page; it is kept raw.
func EncodeErrorHistoryRequest() []byte {
	return []byte{0x00}
}

// EncodeFanModeTemporary returns the 0x40A1 request: ID 0x00 and the FanMode
// value (NN = 0x02). The mode reverts on the next power cycle.
func EncodeFanModeTemporary(mode byte) []byte {
	return []byte{0x00, mode}
}

// EncodeFanModePersistent returns the 0x40CB request: ID 0x00 and the FanMode
// value (NN = 0x02). The mode survives power cycles.
func EncodeFanModePersistent(mode byte) []byte {
	return []byte{0x00, mode}
}

func decodeNumeric(kind Kind, payload []byte) (float64, error) {
	switch kind {
	case KindUCH:
		return float64(payload[0]), nil
	case KindUIR:
		return float64(binary.LittleEndian.Uint16(payload)), nil
	case KindSIR:
		return float64(int16(binary.LittleEndian.Uint16(payload))), nil
	case KindULR:
		return float64(binary.LittleEndian.Uint32(payload)), nil
	}
	return 0, fmt.Errorf("%w: kind %d is not numeric", ErrWidth, kind)
}

func scale(v float64, divisor int) float64 {
	if divisor <= 0 {
		return v
	}
	return v / float64(divisor)
}

// sensorRegisterTable is the 4022 register table transcribed from the Brink 40xx page.
var sensorRegisterTable = []Register{
	{ID: 0x00, Name: "software_version", Kind: KindSTR, Length: 13, Divisor: 0, Unit: "", Clusters: "All"},
	{ID: 0x4E, Name: "firmware_date", Kind: KindSTR, Length: 13, Divisor: 0, Unit: "", Clusters: "All"},
	{ID: 0x4F, Name: "hardware_version", Kind: KindUCH, Length: 0, Divisor: 1, Unit: "", Clusters: "All"},
	{ID: 0x50, Name: "hardware_sub_version", Kind: KindUCH, Length: 0, Divisor: 1, Unit: "", Clusters: "All"},
	{ID: 0x51, Name: "serial_number", Kind: KindSTR, Length: 12, Divisor: 0, Unit: "", Clusters: "1,2,3"},
	{ID: 0x52, Name: "model_name", Kind: KindSTR, Length: 13, Divisor: 0, Unit: "", Clusters: "All"},
	{ID: 0x53, Name: "postal_code", Kind: KindSTR, Length: 7, Divisor: 0, Unit: "", Clusters: "1,2,3"},
	{ID: 0x54, Name: "installation_date", Kind: KindSTR, Length: 11, Divisor: 0, Unit: "", Clusters: "1,2,3"},
	{ID: 0x01, Name: "fan_mode", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "enum", Clusters: "1,2,3,4,5"},
	{ID: 0x02, Name: "inlet_fan_speed", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "rpm", Clusters: "1,2,3,5"},
	{ID: 0x03, Name: "outlet_fan_speed", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "rpm", Clusters: "1,2,5"},
	{ID: 0x04, Name: "inlet_fan_speed_setpoint", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "rpm", Clusters: "1,2,3,5"},
	{ID: 0x05, Name: "outlet_fan_speed_setpoint", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "rpm", Clusters: "1,2,5"},
	{ID: 0x06, Name: "inlet_fan_flow", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "m3/h", Clusters: "1,2,3,5"},
	{ID: 0x0B, Name: "inlet_flow", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "m3/h", Clusters: "1,2,3,5"},
	{ID: 0x0C, Name: "outlet_flow", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "m3/h", Clusters: "1,2,5"},
	{ID: 0x0D, Name: "inlet_flow_setpoint", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "m3/h", Clusters: "1,2,3,5"},
	{ID: 0x12, Name: "fan_mode_0_flow", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "m3/h", Clusters: "1,2,3,5"},
	{ID: 0x13, Name: "fan_mode_1_flow", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "m3/h", Clusters: "1,2,3,5"},
	{ID: 0x27, Name: "outlet_fan_flow", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "m3/h", Clusters: "1,2,5"},
	{ID: 0x28, Name: "outlet_flow_setpoint", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "m3/h", Clusters: "1,2,5"},
	{ID: 0x40, Name: "fan_mode_2_flow", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "m3/h", Clusters: "1,2,3,5"},
	{ID: 0x41, Name: "fan_mode_3_flow", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "m3/h", Clusters: "1,2,5"},
	{ID: 0x07, Name: "inside_temperature", Kind: KindSIR, Length: 0, Divisor: 10, Unit: "deg C", Clusters: "1,2,3,4,5"},
	{ID: 0x08, Name: "outside_temperature", Kind: KindSIR, Length: 0, Divisor: 10, Unit: "deg C", Clusters: "1,2,3,5"},
	{ID: 0x09, Name: "supply_temperature", Kind: KindSIR, Length: 0, Divisor: 10, Unit: "deg C", Clusters: "1,2,3,5"},
	{ID: 0x0A, Name: "exhaust_temperature", Kind: KindSIR, Length: 0, Divisor: 10, Unit: "deg C", Clusters: "1,2"},
	{ID: 0x11, Name: "preheater_temp", Kind: KindSIR, Length: 0, Divisor: 10, Unit: "deg C", Clusters: "1,2,3"},
	{ID: 0x42, Name: "supply_temp_post_bypass", Kind: KindSIR, Length: 0, Divisor: 10, Unit: "deg C", Clusters: "1,2"},
	{ID: 0x43, Name: "exhaust_temp_post_he", Kind: KindSIR, Length: 0, Divisor: 10, Unit: "deg C", Clusters: "1,2"},
	{ID: 0x14, Name: "pressure_inlet", Kind: KindUIR, Length: 0, Divisor: 10, Unit: "Pa", Clusters: "1,2,3,5"},
	{ID: 0x15, Name: "pressure_outlet", Kind: KindUIR, Length: 0, Divisor: 10, Unit: "Pa", Clusters: "1,2,5"},
	{ID: 0x0E, Name: "bypass_status", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "enum", Clusters: "1,2,3"},
	{ID: 0x0F, Name: "preheater_status", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "enum", Clusters: "1,2,3"},
	{ID: 0x10, Name: "preheater_duty_cycle", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "%", Clusters: "1,2,3"},
	{ID: 0x16, Name: "frost_status", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "enum", Clusters: "1,2,3,4,5"},
	{ID: 0x17, Name: "frost_bypass", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "enum", Clusters: "1,2,3"},
	{ID: 0x18, Name: "filter_status", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "enum", Clusters: "1,2,3,5"},
	{ID: 0x19, Name: "filter_hours", Kind: KindULR, Length: 0, Divisor: 1, Unit: "h", Clusters: "1,2,3,5"},
	{ID: 0x1A, Name: "notification_flags_0", Kind: KindUCH, Length: 0, Divisor: 1, Unit: "bitmask", Clusters: "All"},
	{ID: 0x1B, Name: "notification_flags_1", Kind: KindUCH, Length: 0, Divisor: 1, Unit: "bitmask", Clusters: "All"},
	{ID: 0x1C, Name: "notification_flags_2", Kind: KindUCH, Length: 0, Divisor: 1, Unit: "bitmask", Clusters: "All"},
	{ID: 0x1D, Name: "notification_flags_3", Kind: KindUCH, Length: 0, Divisor: 1, Unit: "bitmask", Clusters: "All"},
	{ID: 0x1E, Name: "notification_flags_4", Kind: KindUCH, Length: 0, Divisor: 1, Unit: "bitmask", Clusters: "All"},
	{ID: 0x1F, Name: "notification_flags_5", Kind: KindUCH, Length: 0, Divisor: 1, Unit: "bitmask", Clusters: "All"},
	{ID: 0x57, Name: "notification_count", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "count", Clusters: "All"},
	{ID: 0x58, Name: "active_notification_code", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "code", Clusters: "All"},
	{ID: 0x20, Name: "relative_humidity", Kind: KindSIR, Length: 0, Divisor: 10, Unit: "%", Clusters: "1,2,3"},
	{ID: 0x21, Name: "absolute_humidity", Kind: KindSIR, Length: 0, Divisor: 10, Unit: "g/m3", Clusters: "1,2"},
	{ID: 0x22, Name: "fan_inlet_temp", Kind: KindSIR, Length: 0, Divisor: 10, Unit: "deg C", Clusters: "2"},
	{ID: 0x23, Name: "fan_outlet_temp", Kind: KindSIR, Length: 0, Divisor: 10, Unit: "deg C", Clusters: "2"},
	{ID: 0x24, Name: "relay_output", Kind: KindUCH, Length: 0, Divisor: 1, Unit: "bitmask", Clusters: "2"},
	{ID: 0x25, Name: "analogue_output", Kind: KindUIR, Length: 0, Divisor: 10, Unit: "V", Clusters: "2"},
	{ID: 0x26, Name: "analogue_input", Kind: KindUIR, Length: 0, Divisor: 10, Unit: "V", Clusters: "2"},
	{ID: 0x29, Name: "co2_sensor_1_value", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "ppm", Clusters: "1,2"},
	{ID: 0x2A, Name: "co2_sensor_1_status", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "enum", Clusters: "1,2"},
	{ID: 0x2B, Name: "co2_sensor_2_value", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "ppm", Clusters: "1,2"},
	{ID: 0x2C, Name: "co2_sensor_2_status", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "enum", Clusters: "1,2"},
	{ID: 0x2D, Name: "co2_sensor_3_value", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "ppm", Clusters: "1"},
	{ID: 0x2E, Name: "co2_sensor_3_status", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "enum", Clusters: "1"},
	{ID: 0x2F, Name: "co2_sensor_4_value", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "ppm", Clusters: "1"},
	{ID: 0x30, Name: "co2_sensor_4_status", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "enum", Clusters: "1"},
	{ID: 0x31, Name: "da70_preset", Kind: KindUCH, Length: 0, Divisor: 1, Unit: "enum", Clusters: "3"},
	{ID: 0x32, Name: "da70_co2_value", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "ppm", Clusters: "3"},
	{ID: 0x33, Name: "da70_co2_status", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "enum", Clusters: "3"},
	{ID: 0x34, Name: "elan_setpoint_temp", Kind: KindSIR, Length: 0, Divisor: 10, Unit: "deg C", Clusters: "4"},
	{ID: 0x35, Name: "elan_current_temp", Kind: KindSIR, Length: 0, Divisor: 10, Unit: "deg C", Clusters: "4"},
	{ID: 0x36, Name: "elan_heat_output", Kind: KindUCH, Length: 0, Divisor: 1, Unit: "%", Clusters: "4"},
	{ID: 0x37, Name: "elan_fan_speed", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "rpm", Clusters: "4"},
	{ID: 0x38, Name: "const_rpm_setpoint", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "%", Clusters: "5"},
	{ID: 0x39, Name: "const_rpm_actual", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "%", Clusters: "5"},
	{ID: 0x3A, Name: "mrc_zone_1_co2", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "ppm", Clusters: "6"},
	{ID: 0x3B, Name: "mrc_zone_2_co2", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "ppm", Clusters: "6"},
	{ID: 0x3C, Name: "mrc_zone_3_co2", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "ppm", Clusters: "6"},
	{ID: 0x3D, Name: "mrc_zone_4_co2", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "ppm", Clusters: "6"},
	{ID: 0x3E, Name: "valve_position", Kind: KindUCH, Length: 0, Divisor: 1, Unit: "%", Clusters: "7"},
	{ID: 0x3F, Name: "valve_status", Kind: KindUCH, Length: 0, Divisor: 1, Unit: "enum", Clusters: "7"},
	{ID: 0x44, Name: "heat_exchanger_efficiency", Kind: KindUIR, Length: 0, Divisor: 10, Unit: "%", Clusters: "1,2"},
	{ID: 0x45, Name: "power_consumption_inlet", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "W", Clusters: "1,2,5"},
	{ID: 0x46, Name: "power_consumption_outlet", Kind: KindUIR, Length: 0, Divisor: 1, Unit: "W", Clusters: "1,2,5"},
	{ID: 0x47, Name: "operating_hours_total", Kind: KindULR, Length: 0, Divisor: 1, Unit: "h", Clusters: "1,2,3,4,5"},
	{ID: 0x48, Name: "operating_hours_fan_mode_0", Kind: KindULR, Length: 0, Divisor: 1, Unit: "h", Clusters: "1,2,3,5"},
	{ID: 0x49, Name: "operating_hours_fan_mode_1", Kind: KindULR, Length: 0, Divisor: 1, Unit: "h", Clusters: "1,2,3,5"},
	{ID: 0x4A, Name: "operating_hours_fan_mode_2", Kind: KindULR, Length: 0, Divisor: 1, Unit: "h", Clusters: "1,2,3,5"},
	{ID: 0x4B, Name: "operating_hours_fan_mode_3", Kind: KindULR, Length: 0, Divisor: 1, Unit: "h", Clusters: "1,2,5"},
	{ID: 0x4C, Name: "operating_hours_bypass_open", Kind: KindULR, Length: 0, Divisor: 1, Unit: "h", Clusters: "1,2,3"},
	{ID: 0x4D, Name: "operating_hours_preheater", Kind: KindULR, Length: 0, Divisor: 1, Unit: "h", Clusters: "1,2,3"},
	{ID: 0x55, Name: "bypass_openings", Kind: KindULR, Length: 0, Divisor: 1, Unit: "count", Clusters: "1,2,3"},
	{ID: 0x56, Name: "frost_events", Kind: KindULR, Length: 0, Divisor: 1, Unit: "count", Clusters: "1,2,3,4,5"},
}

// parameterTable is the 4050 parameter table transcribed from the Brink 40xx page; 4080 writes the same IDs.
var parameterTable = []Parameter{
	{ID: 0x01, Name: "flow_mode_1", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 50.0, Max: 300.0, Default: 100.0},
	{ID: 0x02, Name: "flow_mode_2", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 50.0, Max: 300.0, Default: 150.0},
	{ID: 0x03, Name: "flow_mode_3", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 50.0, Max: 300.0, Default: 225.0},
	{ID: 0x04, Name: "bypass_open_temp", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 15.0, Max: 35.0, Default: 24.0},
	{ID: 0x05, Name: "bypass_close_temp", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 10.0, Max: 30.0, Default: 12.0},
	{ID: 0x06, Name: "bypass_temp_high", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 20.0, Max: 40.0, Default: 27.0},
	{ID: 0x07, Name: "bypass_temp_low", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 5.0, Max: 20.0, Default: 12.0},
	{ID: 0x08, Name: "preheater_frost_temp", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: -15.0, Max: 5.0, Default: 0.0},
	{ID: 0x09, Name: "preheater_max_temp", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 0.0, Max: 40.0, Default: 25.0},
	{ID: 0x0A, Name: "filter_dirty_hours", Kind: KindUIR, Divisor: 1, Unit: "h", RangeKnown: true, Min: 500.0, Max: 8760.0, Default: 2160.0},
	{ID: 0x0B, Name: "filter_reset_hours", Kind: KindUIR, Divisor: 1, Unit: "h", RangeKnown: true, Min: 0.0, Max: 8760.0, Default: 0.0},
	{ID: 0x0C, Name: "bypass_min_outside_temp", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 0.0, Max: 30.0, Default: 16.0},
	{ID: 0x0D, Name: "bypass_min_inside_temp", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 15.0, Max: 30.0, Default: 22.0},
	{ID: 0x0E, Name: "bypass_hysteresis", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 0.0, Max: 5.0, Default: 2.0},
	{ID: 0x0F, Name: "preheater_enable", Kind: KindUCH, Divisor: 1, Unit: "bool", RangeKnown: true, Min: 0.0, Max: 1.0, Default: 1.0},
	{ID: 0x10, Name: "bypass_enable", Kind: KindUCH, Divisor: 1, Unit: "bool", RangeKnown: true, Min: 0.0, Max: 1.0, Default: 1.0},
	{ID: 0x11, Name: "frost_bypass_enable", Kind: KindUCH, Divisor: 1, Unit: "bool", RangeKnown: true, Min: 0.0, Max: 1.0, Default: 1.0},
	{ID: 0x12, Name: "fan_balance_enable", Kind: KindUCH, Divisor: 1, Unit: "bool", RangeKnown: true, Min: 0.0, Max: 1.0, Default: 1.0},
	{ID: 0x13, Name: "mode_0_flow_inlet", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 0.0, Max: 50.0, Default: 50.0},
	{ID: 0x14, Name: "mode_0_flow_outlet", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 0.0, Max: 50.0, Default: 50.0},
	{ID: 0x15, Name: "outlet_flow_mode_1", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 50.0, Max: 300.0, Default: 100.0},
	{ID: 0x16, Name: "outlet_flow_mode_2", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 50.0, Max: 300.0, Default: 150.0},
	{ID: 0x17, Name: "outlet_flow_mode_3", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 50.0, Max: 300.0, Default: 225.0},
	{ID: 0x18, Name: "max_flow_inlet", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 50.0, Max: 400.0, Default: 300.0},
	{ID: 0x19, Name: "max_flow_outlet", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 50.0, Max: 400.0, Default: 300.0},
	{ID: 0x1A, Name: "min_flow_inlet", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 10.0, Max: 100.0, Default: 50.0},
	{ID: 0x1B, Name: "min_flow_outlet", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 10.0, Max: 100.0, Default: 50.0},
	{ID: 0x1C, Name: "flow_balance_offset", Kind: KindSIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: -50.0, Max: 50.0, Default: 0.0},
	{ID: 0x1D, Name: "speed_control_p", Kind: KindUIR, Divisor: 1, Unit: "", RangeKnown: true, Min: 1.0, Max: 100.0, Default: 20.0},
	{ID: 0x1E, Name: "speed_control_i", Kind: KindUIR, Divisor: 1, Unit: "", RangeKnown: true, Min: 1.0, Max: 200.0, Default: 60.0},
	{ID: 0x1F, Name: "speed_control_d", Kind: KindUIR, Divisor: 1, Unit: "", RangeKnown: true, Min: 0.0, Max: 100.0, Default: 0.0},
	{ID: 0x20, Name: "pressure_mode", Kind: KindUCH, Divisor: 1, Unit: "enum", RangeKnown: true, Min: 0.0, Max: 2.0, Default: 0.0},
	{ID: 0x21, Name: "flow_mode_0", Kind: KindSIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 0.0, Max: 50.0, Default: 50.0},
	{ID: 0x22, Name: "frost_start_temp", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: -5.0, Max: 5.0, Default: 0.0},
	{ID: 0x23, Name: "frost_stop_temp", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: -5.0, Max: 5.0, Default: 2.0},
	{ID: 0x24, Name: "frost_bypass_mode", Kind: KindUCH, Divisor: 1, Unit: "enum", RangeKnown: true, Min: 0.0, Max: 3.0, Default: 0.0},
	{ID: 0x25, Name: "frost_preheater_mode", Kind: KindUCH, Divisor: 1, Unit: "enum", RangeKnown: true, Min: 0.0, Max: 2.0, Default: 0.0},
	{ID: 0x26, Name: "frost_fan_speed_reduction", Kind: KindUIR, Divisor: 1, Unit: "%", RangeKnown: true, Min: 0.0, Max: 100.0, Default: 50.0},
	{ID: 0x27, Name: "frost_recovery_temp", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 0.0, Max: 10.0, Default: 5.0},
	{ID: 0x28, Name: "bypass_motor_time", Kind: KindUIR, Divisor: 1, Unit: "s", RangeKnown: true, Min: 30.0, Max: 300.0, Default: 150.0},
	{ID: 0x29, Name: "bypass_motor_power", Kind: KindUIR, Divisor: 1, Unit: "%", RangeKnown: true, Min: 10.0, Max: 100.0, Default: 100.0},
	{ID: 0x2A, Name: "preheater_hysteresis", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 0.0, Max: 5.0, Default: 1.0},
	{ID: 0x2B, Name: "preheater_over_temp", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 20.0, Max: 60.0, Default: 40.0},
	{ID: 0x2C, Name: "humidity_enable", Kind: KindUCH, Divisor: 1, Unit: "bool", RangeKnown: true, Min: 0.0, Max: 1.0, Default: 0.0},
	{ID: 0x2D, Name: "humidity_high_threshold", Kind: KindSIR, Divisor: 10, Unit: "%", RangeKnown: true, Min: 50.0, Max: 90.0, Default: 70.0},
	{ID: 0x2E, Name: "humidity_low_threshold", Kind: KindSIR, Divisor: 10, Unit: "%", RangeKnown: true, Min: 30.0, Max: 60.0, Default: 40.0},
	{ID: 0x2F, Name: "humidity_boost_flow", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 50.0, Max: 300.0, Default: 200.0},
	{ID: 0x30, Name: "bypass_temp_hyst", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 0.0, Max: 5.0, Default: 2.0},
	{ID: 0x31, Name: "bypass_night_cooling_enable", Kind: KindUCH, Divisor: 1, Unit: "bool", RangeKnown: true, Min: 0.0, Max: 1.0, Default: 1.0},
	{ID: 0x32, Name: "night_cooling_start_hour", Kind: KindUCH, Divisor: 1, Unit: "h", RangeKnown: true, Min: 0.0, Max: 23.0, Default: 22.0},
	{ID: 0x33, Name: "night_cooling_end_hour", Kind: KindUCH, Divisor: 1, Unit: "h", RangeKnown: true, Min: 0.0, Max: 23.0, Default: 7.0},
	{ID: 0x34, Name: "co2_limit_sensor_1_low", Kind: KindUIR, Divisor: 1, Unit: "ppm", RangeKnown: true, Min: 400.0, Max: 2000.0, Default: 400.0},
	{ID: 0x35, Name: "co2_limit_sensor_1_high", Kind: KindUIR, Divisor: 1, Unit: "ppm", RangeKnown: true, Min: 400.0, Max: 2000.0, Default: 1200.0},
	{ID: 0x36, Name: "co2_limit_sensor_2_low", Kind: KindUIR, Divisor: 1, Unit: "ppm", RangeKnown: true, Min: 400.0, Max: 2000.0, Default: 400.0},
	{ID: 0x37, Name: "co2_limit_sensor_2_high", Kind: KindUIR, Divisor: 1, Unit: "ppm", RangeKnown: true, Min: 400.0, Max: 2000.0, Default: 1200.0},
	{ID: 0x38, Name: "co2_limit_sensor_3_low", Kind: KindUIR, Divisor: 1, Unit: "ppm", RangeKnown: true, Min: 400.0, Max: 2000.0, Default: 400.0},
	{ID: 0x39, Name: "co2_limit_sensor_3_high", Kind: KindUIR, Divisor: 1, Unit: "ppm", RangeKnown: true, Min: 400.0, Max: 2000.0, Default: 1200.0},
	{ID: 0x3A, Name: "co2_limit_sensor_4_low", Kind: KindUIR, Divisor: 1, Unit: "ppm", RangeKnown: true, Min: 400.0, Max: 2000.0, Default: 400.0},
	{ID: 0x3B, Name: "co2_limit_sensor_4_high", Kind: KindUIR, Divisor: 1, Unit: "ppm", RangeKnown: true, Min: 400.0, Max: 2000.0, Default: 1200.0},
	{ID: 0x3C, Name: "co2_boost_flow", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 50.0, Max: 300.0, Default: 225.0},
	{ID: 0x3D, Name: "co2_control_enable", Kind: KindUCH, Divisor: 1, Unit: "bool", RangeKnown: true, Min: 0.0, Max: 1.0, Default: 0.0},
	{ID: 0x3E, Name: "co2_sensor_count", Kind: KindUCH, Divisor: 1, Unit: "count", RangeKnown: true, Min: 0.0, Max: 4.0, Default: 0.0},
	{ID: 0x3F, Name: "da70_preset_count", Kind: KindUCH, Divisor: 1, Unit: "count", RangeKnown: true, Min: 2.0, Max: 5.0, Default: 3.0},
	{ID: 0x40, Name: "da70_preset_1_flow", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 10.0, Max: 200.0, Default: 50.0},
	{ID: 0x41, Name: "da70_preset_2_flow", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 10.0, Max: 200.0, Default: 100.0},
	{ID: 0x42, Name: "da70_preset_3_flow", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 10.0, Max: 200.0, Default: 150.0},
	{ID: 0x43, Name: "da70_preset_4_flow", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 10.0, Max: 200.0, Default: 175.0},
	{ID: 0x44, Name: "da70_preset_5_flow", Kind: KindUIR, Divisor: 1, Unit: "m3/h", RangeKnown: true, Min: 10.0, Max: 200.0, Default: 200.0},
	{ID: 0x45, Name: "elan_setpoint", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 15.0, Max: 30.0, Default: 21.0},
	{ID: 0x46, Name: "elan_hysteresis", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 0.0, Max: 5.0, Default: 1.5},
	{ID: 0x47, Name: "elan_max_power", Kind: KindUIR, Divisor: 1, Unit: "%", RangeKnown: true, Min: 10.0, Max: 100.0, Default: 100.0},
	{ID: 0x48, Name: "const_rpm_setpoint_mode_1", Kind: KindUIR, Divisor: 1, Unit: "%", RangeKnown: true, Min: 10.0, Max: 100.0, Default: 40.0},
	{ID: 0x49, Name: "const_rpm_setpoint_mode_2", Kind: KindUIR, Divisor: 1, Unit: "%", RangeKnown: true, Min: 10.0, Max: 100.0, Default: 60.0},
	{ID: 0x4A, Name: "const_rpm_setpoint_mode_3", Kind: KindUIR, Divisor: 1, Unit: "%", RangeKnown: true, Min: 10.0, Max: 100.0, Default: 80.0},
	{ID: 0x4B, Name: "const_rpm_setpoint_mode_0", Kind: KindUIR, Divisor: 1, Unit: "%", RangeKnown: true, Min: 0.0, Max: 30.0, Default: 20.0},
	{ID: 0x4C, Name: "humidity_control_enable", Kind: KindUCH, Divisor: 1, Unit: "bool", RangeKnown: true, Min: 0.0, Max: 1.0, Default: 0.0},
	{ID: 0x4D, Name: "humidity_boost_enable", Kind: KindUCH, Divisor: 1, Unit: "bool", RangeKnown: true, Min: 0.0, Max: 1.0, Default: 0.0},
	{ID: 0x4E, Name: "frost_start_temp_alt", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: -1.5, Max: 1.5, Default: 0.0},
	{ID: 0x4F, Name: "frost_stop_offset", Kind: KindSIR, Divisor: 10, Unit: "deg C", RangeKnown: true, Min: 0.0, Max: 5.0, Default: 2.0},
	{ID: 0xD0, Name: "ebus_address", Kind: KindUIR, Divisor: 1, Unit: "", RangeKnown: true, Min: 0.0, Max: 254.0, Default: 60.0},
	{ID: 0xD1, Name: "ebus_master_address", Kind: KindUIR, Divisor: 1, Unit: "", RangeKnown: true, Min: 0.0, Max: 254.0, Default: 0.0},
	{ID: 0xD2, Name: "ebus_speed", Kind: KindUIR, Divisor: 1, Unit: "baud", RangeKnown: false, Default: 2400.0},
	{ID: 0xD3, Name: "ebus_protocol_version", Kind: KindUIR, Divisor: 1, Unit: "", RangeKnown: true, Min: 1.0, Max: 2.0, Default: 1.0},
}
