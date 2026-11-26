package input

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/nativepoints"
)

// Purpose is the donor call-site purpose a reply was acquired for.
type Purpose string

// Purposes of the inventory rows.
const (
	PurposeValueRead      Purpose = "value_read"
	PurposeStartupRead    Purpose = "startup_read"
	PurposeDiscoveryProbe Purpose = "discovery_probe"
	PurposeWrite          Purpose = "write"
	PurposeNotSent        Purpose = "not_sent"
)

// Codec names the donor decoder a call site applies to the reply payload.
type Codec string

// Codecs stated from the donor call sites; a codec never follows from catalog metadata alone.
const (
	CodecFloatStrict    Codec = "float32_le_exact_width"
	CodecFloatRaw       Codec = "float32_le_min_width"
	CodecUintWidths     Codec = "uint8_or_uint16_le_exact_width"
	CodecUintWidthsSent Codec = "uint8_or_uint16_le_exact_width_no_ffff"
	CodecUintRaw        Codec = "uint8_or_uint16_le_any_width"
	CodecByteStrict     Codec = "uint8_exact_width_no_ff"
	CodecByteRaw        Codec = "uint8_first_byte"
	CodecByteClass      Codec = "uint8_first_byte_no_ff"
	CodecBoolStrict     Codec = "bool_exact_width"
	CodecBoolRaw        Codec = "bool_first_byte"
	CodecBoolConnected  Codec = "bool_first_byte_is_one"
	CodecBoolExact      Codec = "bool_first_byte_zero_or_one"
	CodecCounterStrict  Codec = "uint32_le_exact_width"
	CodecCounterRaw     Codec = "uint32_le_min_width"
	CodecVersion        Codec = "major_minor_patch"
	CodecTextSanitized  Codec = "text_ascii"
	CodecTextTrimmed    Codec = "text_trimmed"
	CodecDateChecked    Codec = "day_month_year_checked"
	CodecDatePlain      Codec = "day_month_year"
	CodecDateSuppressed Codec = "day_month_year_default_suppressed"
	CodecClock          Codec = "hour_minute"
	CodecPresence       Codec = "first_byte_presence"
	CodecHardwareID     Codec = "uint16_le_min_width_no_ffff"
	CodecProbe          Codec = "reply_header_coherence"
	CodecWriteAck       Codec = "reply_received"
	CodecNotSent        Codec = "not_sent"
)

const (
	extendedMinimumReply      = 4
	extendedLayoutAReply      = 5
	defaultDateSuppressedYear = 2015
)

// NoValueReason is why a reply produced no value.
type NoValueReason string

// Reasons a read produces no value.
const (
	ReasonNoReply    NoValueReason = "no_reply"
	ReasonRejected   NoValueReason = "reply_rejected"
	ReasonWidth      NoValueReason = "width_mismatch"
	ReasonShort      NoValueReason = "too_short"
	ReasonSentinel   NoValueReason = "sentinel"
	ReasonNonFinite  NoValueReason = "non_finite"
	ReasonRange      NoValueReason = "out_of_range"
	ReasonInvalid    NoValueReason = "invalid_value"
	ReasonNotSent    NoValueReason = "not_sent"
	ReasonOutsideRow NoValueReason = "selector_outside_row"
)

// Value is a decoded native value as its kind and canonical text.
type Value struct{ Kind, Text string }

// Decode is one purpose of a row with the codec the donor applies for it.
type Decode struct {
	Purpose Purpose
	Codec   Codec
}

// Row is one extended register selector of the inventory: its opcode, operation, group, instance range and register.
type Row struct {
	ID                                                    string
	Opcode, Operation, Group, FirstInstance, LastInstance byte
	Register                                              uint16
	Decodes                                               []Decode
}

// Result is the native outcome of one acquired reply decoded by one row.
type Result struct {
	RowID       string
	Selector    B524Selector
	Purpose     Purpose
	Codec       Codec
	Outcome     ReadOutcome
	Reason      NoValueReason
	Value       Value
	Unit        string
	Raw         []byte
	Reply       []byte
	AcquiredAt  time.Time
	Diagnostics Diagnostics
	WidthWant   int
}

// CatalogRef is the native point catalog metadata of the row's location.
type CatalogRef struct {
	Name  string
	Unit  string
	Codec nativepoints.Codec
}

var catalogIndex = sync.OnceValue(func() map[[3]uint32]CatalogRef {
	index := map[[3]uint32]CatalogRef{}
	for _, point := range nativepoints.Points() {
		if point.DeviceClass != nativepoints.DeviceClassController {
			continue
		}
		for _, loc := range point.Locations {
			key := [3]uint32{uint32(loc.Opcode), uint32(loc.Group), uint32(loc.Register)}
			if _, seen := index[key]; !seen {
				index[key] = CatalogRef{Name: point.Name, Unit: loc.NativeUnit, Codec: loc.Codec}
			}
		}
	}
	return index
})

// Catalog returns the catalog metadata of the row's opcode, group and register; Name is empty when the catalog has none.
func (r Row) Catalog() CatalogRef {
	return catalogIndex()[[3]uint32{uint32(r.Opcode), uint32(r.Group), uint32(r.Register)}]
}

// Matches reports whether a selector is this row's opcode, operation, group, register and an instance of its range.
func (r Row) Matches(s B524Selector) bool {
	return s.Opcode == r.Opcode && s.Operation == r.Operation && s.Group == r.Group && s.Register == r.Register &&
		s.Instance >= r.FirstInstance && s.Instance <= r.LastInstance
}

type decoded struct {
	value  Value
	reason NoValueReason
	want   int
}

func ok(kind, text string) decoded { return decoded{value: Value{kind, text}} }
func miss(r NoValueReason) decoded { return decoded{reason: r} }
func wrongWidth(want int) decoded  { return decoded{reason: ReasonWidth, want: want} }
func plainDate(b []byte) string    { return fmt.Sprintf("%04d-%02d-%02d", 2000+int(b[2]), b[1], b[0]) }
func floatOf(b []byte) decoded {
	f := float64(math.Float32frombits(binary.LittleEndian.Uint32(b[:4])))
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return miss(ReasonNonFinite)
	}
	return ok("float", strconv.FormatFloat(f, 'g', -1, 64))
}

func uint16Of(b []byte) uint16 {
	if len(b) == 1 {
		return uint16(b[0])
	}
	return binary.LittleEndian.Uint16(b[:2])
}

func cstring(b []byte) []byte {
	if i := bytes.IndexByte(b, 0); i >= 0 {
		return b[:i]
	}
	return b
}

func number(kind string, v uint64) decoded { return ok(kind, strconv.FormatUint(v, 10)) }

func boolOf(b []byte, strictFF bool) decoded {
	if b[0] > 1 {
		if b[0] == 0xff && strictFF {
			return miss(ReasonSentinel)
		}
		return miss(ReasonInvalid)
	}
	return ok("bool", strconv.FormatBool(b[0] == 1))
}

func checkedDate(b []byte) decoded {
	day, month, offset := int(b[0]), int(b[1]), int(b[2])
	if day < 1 || day > 31 || month < 1 || month > 12 || offset > 99 {
		return miss(ReasonRange)
	}
	return ok("date", plainDate(b))
}

func suppressedDate(b []byte) decoded {
	if int(b[2])+2000 == defaultDateSuppressedYear && b[1] == 1 && b[0] == 1 {
		return ok("date", "")
	}
	return ok("date", plainDate(b))
}

func asciiText(b []byte) decoded {
	text := append([]byte(nil), cstring(b)...)
	for i, c := range text {
		if c > 0x7f {
			text[i] = '?'
		}
	}
	return ok("text", strconv.Quote(string(text)))
}

// widthOf is the width a codec needs and whether the donor requires exactly that width.
func widthOf(codec Codec) (need int, exact bool) {
	switch codec {
	case CodecFloatStrict, CodecCounterStrict:
		return 4, true
	case CodecFloatRaw, CodecCounterRaw:
		return 4, false
	case CodecUintWidths, CodecUintWidthsSent:
		return 2, true
	case CodecByteStrict, CodecBoolStrict:
		return 1, true
	case CodecDateChecked:
		return 3, true
	case CodecVersion, CodecDatePlain, CodecDateSuppressed:
		return 3, false
	case CodecClock, CodecHardwareID:
		return 2, false
	}
	return 1, false
}

func decodeBytes(codec Codec, b []byte) decoded {
	need, exact := widthOf(codec)
	if exact && len(b) != need && (need != 2 || len(b) != 1) {
		return wrongWidth(need)
	}
	if !exact && len(b) < need {
		return miss(ReasonShort)
	}
	switch codec {
	case CodecFloatStrict, CodecFloatRaw:
		return floatOf(b)
	case CodecUintWidths, CodecUintRaw:
		return number("uint", uint64(uint16Of(b)))
	case CodecUintWidthsSent, CodecHardwareID:
		if uint16Of(b) == 0xffff {
			return miss(ReasonSentinel)
		}
		return number("uint", uint64(uint16Of(b)))
	case CodecByteStrict, CodecByteClass:
		if b[0] == 0xff {
			return miss(ReasonSentinel)
		}
		return number("uint", uint64(b[0]))
	case CodecByteRaw:
		return number("uint", uint64(b[0]))
	case CodecBoolStrict, CodecBoolRaw:
		return boolOf(b, true)
	case CodecBoolExact:
		return boolOf(b, false)
	case CodecBoolConnected:
		return ok("bool", strconv.FormatBool(b[0] == 1))
	case CodecCounterStrict, CodecCounterRaw:
		if v := binary.LittleEndian.Uint32(b[:4]); v != 0xffffffff {
			return number("uint", uint64(v))
		}
		return miss(ReasonSentinel)
	}
	return decodeStructured(codec, b)
}

func decodeStructured(codec Codec, b []byte) decoded {
	switch codec {
	case CodecVersion:
		if b[0] == 0xff && b[1] == 0xff && b[2] == 0xff {
			return miss(ReasonSentinel)
		}
		return ok("version", fmt.Sprintf("%02d.%02d.%02d", b[0], b[1], b[2]))
	case CodecTextSanitized:
		return asciiText(b)
	case CodecTextTrimmed:
		return ok("text", strconv.Quote(strings.TrimSpace(string(cstring(b)))))
	case CodecDateChecked:
		return checkedDate(b)
	case CodecDatePlain:
		return ok("date", plainDate(b))
	case CodecDateSuppressed:
		return suppressedDate(b)
	case CodecClock:
		return ok("clock", fmt.Sprintf("%02d:%02d", b[0], b[1]))
	case CodecPresence:
		return ok("presence", strconv.FormatBool(b[0] != 0xff))
	}
	return miss(ReasonInvalid)
}

// ParseReply applies the donor's reply header rules to the frame data of a read: it returns the value payload, the
// mismatch line the donor logs for it, and whether the reply is accepted.
func ParseReply(reply []byte, s B524Selector) ([]byte, string, bool) {
	if len(reply) < extendedMinimumReply {
		return nil, "", false
	}
	if len(reply) >= extendedLayoutAReply {
		replyAddr := binary.LittleEndian.Uint16(reply[3:5])
		if reply[2] == s.Group && replyAddr == s.Register {
			if reply[1] != s.Instance && (s.Instance == 0xff || reply[1] != s.Instance+1) {
				return nil, ExtendedInstanceMismatch(s.Opcode, s.Group, s.Instance, s.Register, reply[1], reply[2], replyAddr), false
			}
			return reply[extendedLayoutAReply:], "", len(reply) > extendedLayoutAReply
		}
	}
	replyAddr := binary.LittleEndian.Uint16(reply[2:4])
	if reply[1] != s.Group || replyAddr != s.Register {
		return nil, ExtendedHeaderMismatch(s.Opcode, s.Group, s.Instance, s.Register, reply[0], reply[1], replyAddr, len(reply)), false
	}
	return reply[extendedMinimumReply:], "", len(reply) > extendedMinimumReply
}

// Coherent reports whether a reply names the probed group and register in either header layout, a header-only reply included.
func Coherent(reply []byte, group byte, register uint16) bool {
	return len(reply) >= extendedMinimumReply && (len(reply) >= extendedLayoutAReply && reply[2] == group &&
		binary.LittleEndian.Uint16(reply[3:5]) == register || reply[1] == group && binary.LittleEndian.Uint16(reply[2:4]) == register)
}

type def struct {
	name     string
	register uint16
	decodes  []Decode
}

func d(name string, register uint16, decodes ...Decode) def { return def{name, register, decodes} }
func rv(c Codec) Decode                                     { return Decode{PurposeValueRead, c} }
func rs(c Codec) Decode                                     { return Decode{PurposeStartupRead, c} }

func localRows(add func(prefix string, opcode, group, first, last, operation byte, defs ...def)) {
	add("b524.regulator.", OpcodeLocal, 0x00, 0, 0, OperationRead,
		d("system_off", 0x0007, rv(CodecUintWidths)), d("system_water_pressure", 0x0039, rv(CodecFloatStrict), rs(CodecFloatRaw)),
		d("system_flow_temperature", 0x004B, rv(CodecFloatStrict)), d("system_outdoor_temperature", 0x0073, rv(CodecFloatStrict)),
		d("system_outdoor_temperature_avg_24h", 0x0095, rv(CodecFloatStrict)), d("system_maintenance_due", 0x0096, rv(CodecUintWidths)),
		d("system_hwc_cylinder_temperature_top", 0x009D, rv(CodecFloatStrict)),
		d("system_hwc_cylinder_temperature_bottom", 0x009E, rv(CodecFloatStrict)), d("system_adaptive_heating_curve", 0x0014, rv(CodecUintWidths)),
		d("system_alternative_point", 0x0022, rv(CodecFloatStrict)), d("system_heating_circuit_bivalence_point", 0x0023, rv(CodecFloatStrict)),
		d("system_dhw_bivalence_point", 0x0001, rv(CodecFloatStrict)), d("system_hc_emergency_temperature", 0x0026, rv(CodecFloatStrict)),
		d("system_hwc_max_flow_temp_desired", 0x0046, rv(CodecFloatStrict)), d("system_max_room_humidity", 0x000E, rv(CodecUintWidths)),
		d("system_maintenance_date", 0x002C, rv(CodecDateChecked)), d("system_installer_name_1", 0x006C, rv(CodecTextSanitized)),
		d("system_installer_name_2", 0x006D, rv(CodecTextSanitized)), d("system_installer_phone_1", 0x006F, rv(CodecTextSanitized)),
		d("system_installer_phone_2", 0x0070, rv(CodecTextSanitized)), d("system_installer_menu_code", 0x0076, rv(CodecUintWidths)),
		d("system_scheme", 0x0036, rv(CodecUintWidths), rs(CodecUintRaw)), d("system_module_configuration_vr71", 0x002F, rv(CodecUintWidths), rs(CodecUintRaw)),
	)
	add("b524.energy.", OpcodeLocal, 0x00, 0, 0, OperationRead,
		d("energy_fuel_sum_hc", 0x0056, rv(CodecCounterStrict)), d("energy_electricity_sum_hc", 0x0057, rv(CodecCounterStrict)),
		d("energy_electricity_sum_hwc", 0x0058, rv(CodecCounterStrict)), d("energy_fuel_sum_hwc", 0x0059, rv(CodecCounterStrict)),
		d("energy_fuel_sum_hc_this_month", 0x004E, rv(CodecCounterStrict)), d("energy_electricity_sum_hc_this_month", 0x004F, rv(CodecCounterStrict)),
		d("energy_electricity_sum_hwc_this_month", 0x0050, rv(CodecCounterStrict)),
		d("energy_fuel_sum_hwc_this_month", 0x0051, rv(CodecCounterStrict)), d("energy_fuel_sum_hc_last_month", 0x0052, rv(CodecCounterStrict)),
		d("energy_electricity_sum_hc_last_month", 0x0053, rv(CodecCounterStrict)), d("energy_electricity_sum_hwc_last_month", 0x0054, rv(CodecCounterStrict)),
		d("energy_fuel_sum_hwc_last_month", 0x0055, rv(CodecCounterStrict)),
	)
	add("b524.regulator.write.", OpcodeLocal, 0x00, 0, 0, OperationWrite,
		d("system_installer_name_1", 0x006C, Decode{PurposeWrite, CodecWriteAck}), d("system_installer_name_2", 0x006D, Decode{PurposeWrite, CodecWriteAck}),
		d("system_installer_phone_1", 0x006F, Decode{PurposeWrite, CodecWriteAck}), d("system_installer_phone_2", 0x0070, Decode{PurposeWrite, CodecWriteAck}),
	)
	add("b524.zone.", OpcodeLocal, 0x03, 0, 8, OperationRead,
		d("zone_index", 0x001C, rv(CodecPresence), rs(CodecPresence)), d("zone_name", 0x0016, rv(CodecTextTrimmed)),
		d("zone_name_prefix", 0x0017, rv(CodecTextTrimmed)), d("zone_name_suffix", 0x0018, rv(CodecTextTrimmed)),
		d("zone_target_temp", 0x0022, rv(CodecFloatStrict)), d("zone_fallback_manual_temp", 0x0014, rv(CodecFloatStrict)),
		d("zone_current_humidity", 0x0028, rv(CodecFloatStrict)), d("zone_heating_op_mode", 0x0006, rv(CodecUintWidths)),
		d("zone_quick_veto_temp", 0x0008, rv(CodecFloatStrict)), d("zone_quick_veto_duration", 0x0026, rv(CodecFloatStrict)),
		d("zone_quick_veto_end_time", 0x001E, rv(CodecClock)), d("zone_quick_veto_end_date", 0x0024, rv(CodecDatePlain)),
		d("zone_holiday_start_date", 0x0003, rv(CodecDateSuppressed)), d("zone_holiday_end_date", 0x0004, rv(CodecDateSuppressed)),
		d("zone_holiday_setpoint", 0x0005, rv(CodecFloatStrict)), d("zone_holiday_end_time", 0x0020, rv(CodecClock)),
		d("zone_holiday_start_time", 0x0021, rv(CodecClock)), d("zone_room_temperature_zone_mapping_raw", 0x0013, rv(CodecUintWidths)),
		d("zone_current_temp", 0x000F, rv(CodecFloatStrict)), d("zone_special_function", 0x000E, rv(CodecUintWidths)), d("zone_valve_status", 0x0012, rv(CodecUintWidths)),
	)
	add("b524.circuit.", OpcodeLocal, 0x02, 0, 0x0a, OperationRead,
		d("circuit_type", 0x0002, rv(CodecUintWidths), rs(CodecUintRaw)), d("circuit_cooling_enabled", 0x0006, rv(CodecUintWidths)),
		d("circuit_flow_setpoint", 0x0007, rv(CodecFloatStrict)), d("circuit_flow_temp", 0x0008, rv(CodecFloatStrict), rv(CodecFloatRaw), rs(CodecFloatRaw)),
		d("circuit_heating_curve", 0x000F, rv(CodecFloatStrict)), d("circuit_flow_temp_max", 0x0010, rv(CodecFloatStrict)),
		d("circuit_flow_temp_min", 0x0012, rv(CodecFloatStrict)), d("circuit_summer_limit", 0x0014, rv(CodecFloatStrict)),
		d("circuit_room_temp_control", 0x0015, rv(CodecUintWidths)), d("circuit_circuit_state", 0x001E, rv(CodecUintWidths)),
		d("circuit_special_function_mode", 0x001D, rv(CodecByteStrict)), d("circuit_pump_status", 0x001B, rv(CodecUintWidths)),
		d("circuit_mixer_position", 0x0021, rv(CodecFloatStrict)), d("circuit_humidity", 0x0022, rv(CodecFloatStrict)),
		d("circuit_dew_point", 0x0023, rv(CodecFloatStrict)), d("circuit_pump_hours", 0x0024, rv(CodecCounterStrict)),
		d("circuit_pump_starts", 0x0025, rv(CodecCounterStrict)), d("circuit_calc_flow_temp", 0x0020, Decode{PurposeNotSent, CodecNotSent}),
	)
	add("b524.dhw.", OpcodeLocal, 0x01, 0, 0, OperationRead,
		d("dhw_current_temp", 0x0005, rv(CodecFloatStrict), rs(CodecFloatRaw)), d("dhw_special_function", 0x000D, rv(CodecUintWidths)),
		d("dhw_operation_mode", 0x0003, rv(CodecUintWidths), rs(CodecUintRaw)), d("dhw_target_temp", 0x0004, rv(CodecFloatStrict)),
		d("dhw_holiday_start_date", 0x0009, rv(CodecDateSuppressed)), d("dhw_holiday_end_date", 0x000A, rv(CodecDateSuppressed)),
	)
	add("b524.solar.", OpcodeLocal, 0x04, 0, 0, OperationRead,
		d("solar_enabled", 0x0001, rv(CodecBoolRaw), rs(CodecBoolRaw)), d("solar_function_mode", 0x0002, rv(CodecBoolRaw)),
		d("solar_collector_temp", 0x0003, rv(CodecFloatRaw), rs(CodecFloatRaw)), d("solar_return_temp", 0x0007, rv(CodecFloatRaw)),
		d("solar_pump_active", 0x0008, rv(CodecBoolRaw), rs(CodecBoolRaw)), d("solar_current_yield", 0x0009, rv(CodecFloatRaw)), d("solar_pump_hours", 0x000B, rv(CodecCounterRaw)),
	)
	add("b524.cylinder.", OpcodeLocal, 0x05, 0, 1, OperationRead,
		d("cylinder_max_setpoint", 0x0001, rv(CodecFloatRaw)), d("cylinder_charge_hysteresis", 0x0002, rv(CodecFloatRaw)),
		d("cylinder_charge_offset", 0x0003, rv(CodecFloatRaw)), d("cylinder_temperature", 0x0004, rv(CodecFloatRaw), rs(CodecFloatRaw)),
	)
}

func remoteDefs(forModule bool) []def {
	connected := []Decode{rv(CodecByteStrict), rs(CodecBoolConnected)}
	class := []Decode{rv(CodecByteStrict), rs(CodecByteRaw)}
	firmware, hardware := []Decode{rv(CodecVersion)}, []Decode{rv(CodecUintWidthsSent)}
	if forModule {
		connected, class = []Decode{rv(CodecBoolExact), rs(CodecBoolExact)}, []Decode{rv(CodecByteClass), rs(CodecByteClass)}
		firmware, hardware = []Decode{rv(CodecVersion), rs(CodecVersion)}, []Decode{rv(CodecHardwareID), rs(CodecHardwareID)}
	}
	return []def{
		d("device_slot_connected", 0x0001, connected...), d("device_slot_class_address", 0x0002, class...),
		d("device_slot_firmware", 0x0004, firmware...), d("device_slot_hardware_identifier", 0x0023, hardware...),
		d("device_slot_remote_control_address", 0x0019, rv(CodecByteStrict)), d("device_slot_paired", 0x001e, rv(CodecBoolStrict)),
		d("device_slot_reception_strength", 0x001f, rv(CodecByteStrict)), d("device_slot_zone_assignment", 0x0025, rv(CodecByteStrict)),
		d("device_slot_room_temperature", 0x000f, rv(CodecFloatStrict)), d("device_slot_room_humidity", 0x0007, rv(CodecFloatStrict)),
	}
}

func buildRows() []Row {
	var rows []Row
	add := func(prefix string, opcode, group, first, last, operation byte, defs ...def) {
		for _, df := range defs {
			rows = append(rows, Row{ID: prefix + df.name, Opcode: opcode, Operation: operation, Group: group,
				FirstInstance: first, LastInstance: last, Register: df.register, Decodes: df.decodes})
		}
	}
	localRows(add)
	for _, g := range []struct {
		name  string
		group byte
	}{{"remoteRegulators", 0x09}, {"remoteThermostats", 0x0a}, {"remoteFunctionalModules", 0x0c}} {
		add("b524.remote."+g.name+".", OpcodeRemote, g.group, 0, 0x0a, OperationRead, remoteDefs(g.group == 0x0c)...)
	}
	probe := Decode{PurposeDiscoveryProbe, CodecProbe}
	add("b524.discovery.", OpcodeLocal, 0x00, 0, 0, OperationRead, d("capability_probe_regulator", 0x0001, probe))
	add("b524.discovery.", OpcodeLocal, 0x01, 0, 0, OperationRead, d("capability_probe_hot_water_mode", 0x0003, probe))
	return rows
}

var extendedRows = buildRows()

// ExtendedRows returns a copy of every extended register row of the inventory.
func ExtendedRows() []Row { return append([]Row(nil), extendedRows...) }

// Decode decodes one acquired reply with one of the row's decodes.
func (r Row) Decode(rec Record, dec Decode) Result {
	res := Result{RowID: r.ID, Purpose: dec.Purpose, Codec: dec.Codec, Outcome: OutcomeNoValue, Reply: rec.Reply,
		AcquiredAt: rec.Clock, Unit: r.Catalog().Unit, Diagnostics: Diagnostics{Outcome: OutcomeNoValue, SendError: rec.Diagnostics.SendError}}
	sel, isExtended := rec.Selector.(B524Selector)
	if !isExtended || !r.Matches(sel) {
		res.Reason = ReasonOutsideRow
		return res
	}
	res.Selector = sel
	switch {
	case dec.Codec == CodecNotSent:
		res.Reason = ReasonNotSent
	case !rec.Replied:
		res.Reason = ReasonNoReply
	case dec.Codec == CodecWriteAck:
		res.succeed(Value{"ack", "received"})
	case dec.Codec == CodecProbe:
		res.succeed(Value{"coherent", strconv.FormatBool(Coherent(rec.Reply, sel.Group, sel.Register))})
	default:
		res.decodePayload(sel, dec.Codec)
	}
	return res
}

func (res *Result) succeed(v Value) {
	res.Outcome, res.Diagnostics.Outcome, res.Value = OutcomeValue, OutcomeValue, v
}

func (res *Result) decodePayload(sel B524Selector, codec Codec) {
	payload, line, accepted := ParseReply(res.Reply, sel)
	res.Diagnostics.MismatchText = line
	if !accepted {
		res.Reason = ReasonRejected
		return
	}
	res.Raw = payload
	out := decodeBytes(codec, payload)
	if out.reason != "" {
		res.Reason, res.WidthWant = out.reason, out.want
		return
	}
	res.succeed(out.value)
}
