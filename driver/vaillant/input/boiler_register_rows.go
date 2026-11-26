package input

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

// BoilerCodec names the byte layout the poller decodes a boiler register tail with.
type BoilerCodec string

const (
	// BoilerCodecUnsignedFixed is an unsigned little-endian word divided by 256.
	BoilerCodecUnsignedFixed BoilerCodec = "DATA2b"
	// BoilerCodecSignedFixed is a signed little-endian word divided by 16.
	BoilerCodecSignedFixed BoilerCodec = "DATA2c"
	// BoilerCodecUCH is the first tail byte.
	BoilerCodecUCH BoilerCodec = "UCH"
	// BoilerCodecUIN is an unsigned little-endian word.
	BoilerCodecUIN BoilerCodec = "UIN"
	// BoilerCodecSIN is a signed little-endian word, divided by the row divisor when one is set.
	BoilerCodecSIN BoilerCodec = "SIN"
	// BoilerCodecOnOff is true for every first tail byte except 0xF0.
	BoilerCodecOnOff BoilerCodec = "OnOff"
	// BoilerCodecHoursum2 is an unsigned little-endian word times two.
	BoilerCodecHoursum2 BoilerCodec = "Hoursum2"
	// BoilerCodecHundredths is an unsigned little-endian word divided by 100.
	BoilerCodecHundredths BoilerCodec = "UIN100"
	// BoilerCodecPercent0 is the first tail byte as a percentage.
	BoilerCodecPercent0 BoilerCodec = "Percent0"
	// BoilerCodecBCDPhone is packed decimal digits that end at the first nibble above nine.
	BoilerCodecBCDPhone BoilerCodec = "BCDPhone"
	// BoilerCodecRaw keeps the tail undecoded.
	BoilerCodecRaw BoilerCodec = "Raw"
)

// BoilerValueKind is the Go type of the value the poller publishes for a row.
type BoilerValueKind string

const (
	// BoilerKindFloat is a float64 value.
	BoilerKindFloat BoilerValueKind = "float"
	// BoilerKindInt is an int value truncated from the decoded number.
	BoilerKindInt BoilerValueKind = "int"
	// BoilerKindBool is a bool value.
	BoilerKindBool BoilerValueKind = "bool"
	// BoilerKindText is a string value.
	BoilerKindText BoilerValueKind = "text"
	// BoilerKindRaw is the tail as lower-case hexadecimal, undecoded.
	BoilerKindRaw BoilerValueKind = "raw"
)

// CodecSource says where a row's codec comes from.
type CodecSource string

const (
	// CodecSourceCatalog marks a codec the native catalog declares in the same terms as the poller decodes.
	CodecSourceCatalog CodecSource = "catalog"
	// CodecSourceGatewayDecoder marks a codec taken from the poller's call-site decoder, either because the
	// catalog has no entry for the register or because its declaration does not state the poller's conversion.
	CodecSourceGatewayDecoder CodecSource = "gateway_decoder"
)

// BoilerRegisterRow is one B509 read register: layout, codec, scale and the unit the poller names for its field.
// Catalog is the quantity of the native catalog entry for the same selector and is empty when the catalog has none;
// the catalog's own metadata is not copied here. Unit is empty when the poller's field names none.
type BoilerRegisterRow struct {
	ID, Name    string
	Register    uint16
	Codec       BoilerCodec
	Kind        BoilerValueKind
	Divisor     float64
	Unit        string
	Catalog     string
	CodecSource CodecSource
	NativeField string
}

// BoilerWriteRow is one B509 write register. Encoder names the encoder of the value bytes the poller's callers pass.
type BoilerWriteRow struct {
	ID, Name    string
	Register    uint16
	Encoder     BoilerCodec
	NativeField string
}

// Selector returns the B509 read request of the row.
func (r BoilerRegisterRow) Selector() B509Selector {
	return B509Selector{Opcode: BoilerRead, Register: r.Register}
}

// Selector returns the B509 write request header of the row.
func (r BoilerWriteRow) Selector() B509Selector {
	return B509Selector{Opcode: BoilerWrite, Register: r.Register}
}

// BoilerRegisterRows lists every B509 register the poller reads.
var BoilerRegisterRows = []BoilerRegisterRow{
	{ID: "b509.boiler_b509_water_pressure", Name: "boiler_b509_water_pressure", Register: 0x0200, Codec: BoilerCodecUnsignedFixed, Kind: BoilerKindFloat, Unit: "bar", Catalog: "water_pressure", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.WaterPressureBar"},
	{ID: "b509.boiler_b509_flame_active", Name: "boiler_b509_flame_active", Register: 0x0500, Codec: BoilerCodecOnOff, Kind: BoilerKindBool, Unit: "", Catalog: "flame_active", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.FlameActive"},
	{ID: "b509.boiler_b509_partload_hc_kw", Name: "boiler_b509_partload_hc_kw", Register: 0x0704, Codec: BoilerCodecUCH, Kind: BoilerKindFloat, Unit: "kW", Catalog: "heating_partial_load_power", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.PartloadHcKW"},
	{ID: "b509.boiler_b509_partload_hwc_kw", Name: "boiler_b509_partload_hwc_kw", Register: 0x0804, Codec: BoilerCodecUCH, Kind: BoilerKindFloat, Unit: "kW", Catalog: "dhw_partial_load_power", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.PartloadHwcKW"},
	{ID: "b509.boiler_b509_flowset_hc_max_c", Name: "boiler_b509_flowset_hc_max_c", Register: 0x0E04, Codec: BoilerCodecSignedFixed, Kind: BoilerKindFloat, Unit: "degC", Catalog: "maximum_heating_flow_temperature_setpoint", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.FlowsetHcMaxC"},
	{ID: "b509.boiler_b509_flowset_hc_max_c_fallback", Name: "boiler_b509_flowset_hc_max_c_fallback", Register: 0xA500, Codec: BoilerCodecSignedFixed, Kind: BoilerKindFloat, Unit: "degC", Catalog: "maximum_heating_flow_temperature_setpoint", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.FlowsetHcMaxC"},
	{ID: "b509.boiler_b509_flowset_hwc_max_c", Name: "boiler_b509_flowset_hwc_max_c", Register: 0x0F04, Codec: BoilerCodecSignedFixed, Kind: BoilerKindFloat, Unit: "degC", Catalog: "maximum_dhw_flow_temperature_setpoint", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.FlowsetHwcMaxC"},
	{ID: "b509.boiler_b509_pump_hours", Name: "boiler_b509_pump_hours", Register: 0x1400, Codec: BoilerCodecHoursum2, Kind: BoilerKindFloat, Unit: "h", Catalog: "", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.PumpHours"},
	{ID: "b509.boiler_b509_flow_temperature", Name: "boiler_b509_flow_temperature", Register: 0x1800, Codec: BoilerCodecRaw, Kind: BoilerKindRaw, Unit: "", Catalog: "flow_temperature", CodecSource: CodecSourceGatewayDecoder, NativeField: ""},
	{ID: "b509.boiler_b509_fan_hours", Name: "boiler_b509_fan_hours", Register: 0x1B00, Codec: BoilerCodecHoursum2, Kind: BoilerKindFloat, Unit: "h", Catalog: "", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.FanHours"},
	{ID: "b509.boiler_b509_deactivations_ifc", Name: "boiler_b509_deactivations_ifc", Register: 0x1F00, Codec: BoilerCodecUCH, Kind: BoilerKindInt, Unit: "", Catalog: "", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.DeactivationsIFC"},
	{ID: "b509.boiler_b509_hours_till_service", Name: "boiler_b509_hours_till_service", Register: 0x2004, Codec: BoilerCodecHoursum2, Kind: BoilerKindInt, Unit: "h", Catalog: "", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.HoursTillService"},
	{ID: "b509.boiler_b509_deactivations_limit", Name: "boiler_b509_deactivations_limit", Register: 0x2000, Codec: BoilerCodecUCH, Kind: BoilerKindInt, Unit: "", Catalog: "", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.DeactivationsTemplimiter"},
	{ID: "b509.boiler_b509_dhw_hours", Name: "boiler_b509_dhw_hours", Register: 0x2200, Codec: BoilerCodecHoursum2, Kind: BoilerKindFloat, Unit: "h", Catalog: "", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.DhwHours"},
	{ID: "b509.boiler_b509_dhw_starts", Name: "boiler_b509_dhw_starts", Register: 0x2300, Codec: BoilerCodecUIN, Kind: BoilerKindInt, Unit: "", Catalog: "", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.DhwStarts"},
	{ID: "b509.boiler_b509_target_fan_speed_rpm", Name: "boiler_b509_target_fan_speed_rpm", Register: 0x2400, Codec: BoilerCodecUIN, Kind: BoilerKindInt, Unit: "rpm", Catalog: "fan_speed_setpoint", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.TargetFanSpeedRpm"},
	{ID: "b509.boiler_b509_central_heating_hours", Name: "boiler_b509_central_heating_hours", Register: 0x2800, Codec: BoilerCodecHoursum2, Kind: BoilerKindFloat, Unit: "h", Catalog: "", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.CentralHeatingHours"},
	{ID: "b509.boiler_b509_central_heating_starts", Name: "boiler_b509_central_heating_starts", Register: 0x2900, Codec: BoilerCodecUIN, Kind: BoilerKindInt, Unit: "", Catalog: "", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.CentralHeatingStarts"},
	{ID: "b509.boiler_b509_modulation_pct", Name: "boiler_b509_modulation_pct", Register: 0x2E00, Codec: BoilerCodecSIN, Kind: BoilerKindFloat, Divisor: 10, Unit: "%", Catalog: "modulation_setpoint", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.ModulationPct"},
	{ID: "b509.boiler_b509_flow_temp_desired_c", Name: "boiler_b509_flow_temp_desired_c", Register: 0x3900, Codec: BoilerCodecSignedFixed, Kind: BoilerKindFloat, Unit: "degC", Catalog: "flow_temperature_setpoint", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.FlowTempDesiredC"},
	{ID: "b509.boiler_b509_external_pump_active", Name: "boiler_b509_external_pump_active", Register: 0x3F00, Codec: BoilerCodecOnOff, Kind: BoilerKindBool, Unit: "", Catalog: "external_pump_active", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.ExternalPumpActive"},
	{ID: "b509.boiler_b509_installer_menu_code", Name: "boiler_b509_installer_menu_code", Register: 0x4904, Codec: BoilerCodecUCH, Kind: BoilerKindInt, Unit: "", Catalog: "", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.InstallerMenuCode"},
	{ID: "b509.boiler_b509_central_heating_pump", Name: "boiler_b509_central_heating_pump", Register: 0x4400, Codec: BoilerCodecOnOff, Kind: BoilerKindBool, Unit: "", Catalog: "central_heating_pump_active", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.CentralHeatingPumpActive"},
	{ID: "b509.boiler_b509_diverter_valve_position", Name: "boiler_b509_diverter_valve_position", Register: 0x5400, Codec: BoilerCodecUCH, Kind: BoilerKindFloat, Unit: "%", Catalog: "diverter_valve_position", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.DiverterValvePositionPct"},
	{ID: "b509.boiler_b509_dhw_water_flow_lpm", Name: "boiler_b509_dhw_water_flow_lpm", Register: 0x5500, Codec: BoilerCodecHundredths, Kind: BoilerKindFloat, Unit: "l/min", Catalog: "dhw_water_flow_rate", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.DhwWaterFlowLpm"},
	{ID: "b509.boiler_b509_dhw_demand_active", Name: "boiler_b509_dhw_demand_active", Register: 0x5800, Codec: BoilerCodecOnOff, Kind: BoilerKindBool, Unit: "", Catalog: "dhw_demand_active", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.DhwDemandActive"},
	{ID: "b509.boiler_b509_circulation_pump_active", Name: "boiler_b509_circulation_pump_active", Register: 0x7B00, Codec: BoilerCodecOnOff, Kind: BoilerKindBool, Unit: "", Catalog: "circulation_pump_active", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.CirculationPumpActive"},
	{ID: "b509.boiler_b509_phone_number", Name: "boiler_b509_phone_number", Register: 0x8104, Codec: BoilerCodecBCDPhone, Kind: BoilerKindText, Unit: "", Catalog: "", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.PhoneNumber"},
	{ID: "b509.boiler_b509_fan_speed_rpm", Name: "boiler_b509_fan_speed_rpm", Register: 0x8300, Codec: BoilerCodecUIN, Kind: BoilerKindInt, Unit: "rpm", Catalog: "fan_speed", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.FanSpeedRpm"},
	{ID: "b509.boiler_b509_storage_load_pump_pct", Name: "boiler_b509_storage_load_pump_pct", Register: 0x9E00, Codec: BoilerCodecPercent0, Kind: BoilerKindFloat, Unit: "%", Catalog: "storage_loading_pump_level", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.StorageLoadPumpPct"},
	{ID: "b509.boiler_b509_ionisation_voltage_ua", Name: "boiler_b509_ionisation_voltage_ua", Register: 0xA400, Codec: BoilerCodecSIN, Kind: BoilerKindFloat, Divisor: 10, Unit: "uA", Catalog: "flame_ionization_current", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.IonisationVoltageUa"},
	{ID: "b509.boiler_b509_state_number", Name: "boiler_b509_state_number", Register: 0xAB00, Codec: BoilerCodecUCH, Kind: BoilerKindInt, Unit: "", Catalog: "boiler_state_code", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.StateNumber"},
	{ID: "b509.boiler_b509_gas_valve_active", Name: "boiler_b509_gas_valve_active", Register: 0xBB00, Codec: BoilerCodecOnOff, Kind: BoilerKindBool, Unit: "", Catalog: "gas_valve_active", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.GasValveActive"},
	{ID: "b509.boiler_b509_heating_switch_active", Name: "boiler_b509_heating_switch_active", Register: 0xF203, Codec: BoilerCodecOnOff, Kind: BoilerKindBool, Unit: "", Catalog: "heating_switch_active", CodecSource: CodecSourceGatewayDecoder, NativeField: "vaillantBoilerSnapshot.HeatingSwitchActive"},
	{ID: "b509.boiler_b509_primary_circuit_flow_lpm", Name: "boiler_b509_primary_circuit_flow_lpm", Register: 0xFB00, Codec: BoilerCodecHundredths, Kind: BoilerKindFloat, Unit: "l/min", Catalog: "primary_circuit_flow_rate", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.PrimaryCircuitFlowLpm"},
	{ID: "b509.boiler_b509_dhw_temp_desired_c", Name: "boiler_b509_dhw_temp_desired_c", Register: 0xEA03, Codec: BoilerCodecSignedFixed, Kind: BoilerKindFloat, Unit: "degC", Catalog: "dhw_temperature_setpoint", CodecSource: CodecSourceCatalog, NativeField: "vaillantBoilerSnapshot.DhwTempDesiredC"},
}

// BoilerWriteRows lists every B509 register the poller writes.
var BoilerWriteRows = []BoilerWriteRow{
	{ID: "b509.write.boiler_b509_flowset_hc_max_c", Name: "boiler_b509_flowset_hc_max_c", Register: 0x0E04, Encoder: BoilerCodecSignedFixed, NativeField: "vaillantBoilerSnapshot.FlowsetHcMaxC"},
	{ID: "b509.write.boiler_b509_flowset_hc_max_c_fallback", Name: "boiler_b509_flowset_hc_max_c_fallback", Register: 0xA500, Encoder: BoilerCodecSignedFixed, NativeField: "vaillantBoilerSnapshot.FlowsetHcMaxC"},
	{ID: "b509.write.boiler_b509_flowset_hwc_max_c", Name: "boiler_b509_flowset_hwc_max_c", Register: 0x0F04, Encoder: BoilerCodecSignedFixed, NativeField: "vaillantBoilerSnapshot.FlowsetHwcMaxC"},
	{ID: "b509.write.boiler_b509_partload_hc_kw", Name: "boiler_b509_partload_hc_kw", Register: 0x0704, Encoder: BoilerCodecUCH, NativeField: "vaillantBoilerSnapshot.PartloadHcKW"},
	{ID: "b509.write.boiler_b509_partload_hwc_kw", Name: "boiler_b509_partload_hwc_kw", Register: 0x0804, Encoder: BoilerCodecUCH, NativeField: "vaillantBoilerSnapshot.PartloadHwcKW"},
	{ID: "b509.write.boiler_b509_phone_number", Name: "boiler_b509_phone_number", Register: 0x8104, Encoder: BoilerCodecBCDPhone, NativeField: "vaillantBoilerSnapshot.PhoneNumber"},
	{ID: "b509.write.boiler_b509_installer_menu_code", Name: "boiler_b509_installer_menu_code", Register: 0x4904, Encoder: BoilerCodecUCH, NativeField: "vaillantBoilerSnapshot.InstallerMenuCode"},
}

// BoilerReadReason is why a boiler register read or write produced no value or no acknowledgement.
type BoilerReadReason string

const (
	// BoilerReasonNone is the reason of a read that produced a value or a write that was acknowledged.
	BoilerReasonNone BoilerReadReason = ""
	// BoilerReasonSendError is a send that returned an error.
	BoilerReasonSendError BoilerReadReason = "send_error"
	// BoilerReasonNoReply is a send that returned neither a frame nor an error.
	BoilerReasonNoReply BoilerReadReason = "no_reply"
	// BoilerReasonPayloadRejected is a read reply the poller's payload parser rejects.
	BoilerReasonPayloadRejected BoilerReadReason = "payload_rejected"
	// BoilerReasonDecodeRejected is a tail the codec cannot decode.
	BoilerReasonDecodeRejected BoilerReadReason = "decode_rejected"
	// BoilerReasonAckInvalid is a write reply that is not an acknowledgement.
	BoilerReasonAckInvalid BoilerReadReason = "ack_invalid"
)

const (
	boilerWriteNilResponseText = "b509 write returned nil response"
	boilerWriteAckInvalidText  = "b509 write ack invalid: target=0x%02x addr=0x%04x"
)

// ErrBoilerSelectorMismatch is returned when a record was acquired for another selector than the row's.
var ErrBoilerSelectorMismatch = errors.New("record selector is not the selector of the row")

// BoilerReadResult is the native result of one acquired B509 read: the evidence as received and either a value or the reason there is none.
type BoilerReadResult struct {
	Row        string
	Request    []byte
	Raw        []byte
	Tail       []byte
	Outcome    ReadOutcome
	Reason     BoilerReadReason
	Value      any
	Kind       BoilerValueKind
	Unit       string
	AcquiredAt time.Time
	Source     byte
	Target     byte
}

// BoilerWriteResult is the native result of one acquired B509 write reply.
type BoilerWriteResult struct {
	Row        string
	Request    []byte
	Raw        []byte
	Accepted   bool
	Reason     BoilerReadReason
	Error      string
	AcquiredAt time.Time
	Source     byte
	Target     byte
}

// ParseBoilerReadPayload strips the opcode or register header from a read reply. It returns false for an empty
// reply and for a reply that is only a header; a reply with another header is returned unchanged.
func ParseBoilerReadPayload(reply []byte, register uint16) ([]byte, bool) {
	if len(reply) == 0 {
		return nil, false
	}
	hi, lo := byte(register>>8), byte(register)
	if len(reply) >= 3 && reply[0] == BoilerRead && reply[1] == hi && reply[2] == lo {
		return reply[3:], len(reply) > 3
	}
	if len(reply) >= 2 && reply[0] == hi && reply[1] == lo {
		return reply[2:], len(reply) > 2
	}
	return reply, true
}

// ParseBoilerWriteAck reports whether a write reply acknowledges the write: an empty reply does, and so does a reply
// whose first byte after an optional header is not zero.
func ParseBoilerWriteAck(reply []byte, register uint16) bool {
	if len(reply) == 0 {
		return true
	}
	if len(reply) == 1 {
		return reply[0] != 0
	}
	hi, lo := byte(register>>8), byte(register)
	if len(reply) >= 3 && reply[0] == BoilerWrite && reply[1] == hi && reply[2] == lo {
		return len(reply) == 3 || reply[3] != 0
	}
	if len(reply) >= 2 && reply[0] == hi && reply[1] == lo {
		return len(reply) == 2 || reply[2] != 0
	}
	return reply[0] != 0
}

// Read evaluates one acquired reply against the row.
func (r BoilerRegisterRow) Read(rec Record) (BoilerReadResult, error) {
	if s, ok := rec.Selector.(B509Selector); !ok || s != r.Selector() {
		return BoilerReadResult{}, ErrBoilerSelectorMismatch
	}
	out := BoilerReadResult{Row: r.ID, Request: rec.Request, Raw: rec.Reply, Outcome: OutcomeNoValue, Kind: r.Kind, Unit: r.Unit, AcquiredAt: rec.Clock, Source: rec.Source, Target: rec.Target}
	switch {
	case rec.Diagnostics.SendError != "":
		out.Reason = BoilerReasonSendError
		return out, nil
	case !rec.Replied:
		out.Reason = BoilerReasonNoReply
		return out, nil
	}
	tail, ok := ParseBoilerReadPayload(rec.Reply, r.Register)
	if !ok {
		out.Reason = BoilerReasonPayloadRejected
		return out, nil
	}
	out.Tail = tail
	value, ok := r.decode(tail)
	if !ok {
		out.Reason = BoilerReasonDecodeRejected
		return out, nil
	}
	out.Value, out.Outcome = value, OutcomeValue
	return out, nil
}

func (r BoilerRegisterRow) decode(tail []byte) (any, bool) {
	switch r.Codec {
	case BoilerCodecRaw:
		return hex.EncodeToString(tail), true
	case BoilerCodecOnOff:
		return tail[0] != 0xf0, true
	case BoilerCodecBCDPhone:
		text := decodeBoilerBCD(tail)
		return text, text != ""
	}
	width := 2
	if r.Codec == BoilerCodecUCH || r.Codec == BoilerCodecPercent0 {
		width = 1
	}
	if len(tail) < width {
		return nil, false
	}
	var number float64
	switch r.Codec {
	case BoilerCodecUnsignedFixed:
		number = float64(binary.LittleEndian.Uint16(tail)) / 256
	case BoilerCodecSignedFixed:
		number = float64(int16(binary.LittleEndian.Uint16(tail))) / 16
	case BoilerCodecUCH, BoilerCodecPercent0:
		number = float64(tail[0])
	case BoilerCodecUIN:
		number = float64(binary.LittleEndian.Uint16(tail))
	case BoilerCodecSIN:
		number = float64(int16(binary.LittleEndian.Uint16(tail)))
		if r.Divisor != 0 {
			number /= r.Divisor
		}
	case BoilerCodecHoursum2:
		number = float64(binary.LittleEndian.Uint16(tail)) * 2
	case BoilerCodecHundredths:
		number = float64(binary.LittleEndian.Uint16(tail)) / 100
	default:
		return nil, false
	}
	if r.Kind == BoilerKindInt {
		return int(number), true
	}
	return number, true
}

func decodeBoilerBCD(tail []byte) string {
	var digits strings.Builder
	for _, b := range tail {
		hi, lo := b>>4, b&0x0f
		if hi > 9 {
			break
		}
		digits.WriteByte('0' + hi)
		if lo > 9 {
			break
		}
		digits.WriteByte('0' + lo)
	}
	return digits.String()
}

// Acknowledge evaluates one acquired write reply against the row. Error is the text the poller's writer returns for the attempt.
func (r BoilerWriteRow) Acknowledge(rec Record) (BoilerWriteResult, error) {
	if s, ok := rec.Selector.(B509Selector); !ok || s != r.Selector() {
		return BoilerWriteResult{}, ErrBoilerSelectorMismatch
	}
	out := BoilerWriteResult{Row: r.ID, Request: rec.Request, Raw: rec.Reply, AcquiredAt: rec.Clock, Source: rec.Source, Target: rec.Target}
	switch {
	case rec.Diagnostics.SendError != "":
		out.Reason, out.Error = BoilerReasonSendError, rec.Diagnostics.SendError
	case !rec.Replied:
		out.Reason, out.Error = BoilerReasonNoReply, boilerWriteNilResponseText
	case !ParseBoilerWriteAck(rec.Reply, r.Register):
		out.Reason, out.Error = BoilerReasonAckInvalid, fmt.Sprintf(boilerWriteAckInvalidText, rec.Target, r.Register)
	default:
		out.Accepted = true
	}
	return out, nil
}
