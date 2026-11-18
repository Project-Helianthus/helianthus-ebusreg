package nativepoints_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/nativepoints"
)

func findLocation(point nativepoints.Point, opcode, group uint8, register uint16) (nativepoints.Location, bool) {
	for _, location := range point.Locations {
		if location.Opcode == opcode && location.Group == group && location.Register == register {
			return location, true
		}
	}
	return nativepoints.Location{}, false
}

func TestNativeLocalParameterRowHasAPoint(t *testing.T) {
	point := mustLookup(t, nativepoints.DeviceClassController, "sensor_mode")
	location, found := findLocation(point, 0x02, 0x0A, 0x0001)
	if !found || len(point.Locations) != 1 {
		t.Fatalf("sensor_mode locations %+v", point.Locations)
	}
	if location.InstanceRule != nativepoints.InstanceIndexed || location.InstanceMin != 0 || location.InstanceMax != 0x0A {
		t.Errorf("local parameter instance rule %+v", location)
	}
	if point.Codec != nativepoints.CodecUint8 || point.ValueClass != nativepoints.ValueEnumeration || point.NameStatus != nativepoints.NameProvisional {
		t.Errorf("sensor_mode = codec %q class %q status %q", point.Codec, point.ValueClass, point.NameStatus)
	}
}

func TestNativeRemoteThermostatRoomTemperatureOffsetHasAPoint(t *testing.T) {
	point := mustLookup(t, nativepoints.DeviceClassController, "room_temperature_offset")
	if point.Codec != nativepoints.CodecFloat32LE || point.NativeUnit != "°C" || point.ValueClass != nativepoints.ValueNumeric {
		t.Fatalf("room_temperature_offset = %+v", point)
	}
	for _, group := range []uint8{0x09, 0x0A} {
		location, found := findLocation(point, 0x06, group, 0x000E)
		if !found {
			t.Fatalf("no opcode 0x06 group 0x%02x location: %+v", group, point.Locations)
		}
		if location.InstanceRule != nativepoints.InstanceIndexed || location.InstanceMin != 1 || location.InstanceMax != 8 {
			t.Errorf("group 0x%02x slot rule %+v", group, location)
		}
	}
}

type matrixRow struct {
	Family      string `json:"family"`
	Opcode      uint8  `json:"opcode"`
	Group       uint8  `json:"group"`
	Register    uint16 `json:"register"`
	Name        string `json:"name"`
	ReplyWidths []int  `json:"reply_widths"`
	ValueType   string `json:"value_type"`
	Unit        string `json:"unit"`
	Enumeration string `json:"enumeration"`
}

func registerMatrix(t *testing.T) []matrixRow {
	t.Helper()
	raw, err := os.ReadFile("../../catalog/vaillant/data/b524_register_matrix.json")
	if err != nil {
		t.Fatal(err)
	}
	var rows []matrixRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 {
		t.Fatal("empty register matrix")
	}
	return rows
}

func matrixCodec(row matrixRow) nativepoints.Codec {
	switch row.ValueType {
	case "float32 LE":
		return nativepoints.CodecFloat32LE
	case "uint16 LE":
		return nativepoints.CodecUint16LE
	case "uint8":
		return nativepoints.CodecUint8
	case "int8":
		return nativepoints.CodecSignedByte
	case "uint8 or uint16 LE":
		return nativepoints.CodecUint8OrUint16LE
	case "3 x uint8":
		if len(row.Enumeration) >= 5 && row.Enumeration[:5] == "date:" {
			return nativepoints.CodecDayMonthYear
		}
		return nativepoints.CodecHourMinuteSecond
	}
	return nativepoints.CodecUnknown
}

func TestNativeZoneCoolingSetpointTakesTheMatrixLayout(t *testing.T) {
	point := mustLookup(t, nativepoints.DeviceClassController, "cooling_zone_room_temperature_setpoint")
	if point.Codec != nativepoints.CodecFloat32LE || point.NativeUnit != "°C" || fmt.Sprint(point.ControllerFamilies) != fmt.Sprint(controllerFamilies) {
		t.Errorf("cooling_zone_room_temperature_setpoint = %+v", point)
	}
	checked, mapped := 0, mapLayouts(t)
	for _, row := range registerMatrix(t) {
		want := matrixCodec(row)
		if want == nativepoints.CodecUnknown {
			continue
		}
		checked++
		// A blank matrix unit keeps the unit the map's Decode column states.
		unit := row.Unit
		if unit == "" {
			unit = mapped[layoutKey{row.Opcode, row.Group, row.Register}].unit
		}
		matched := false
		for _, candidate := range nativepoints.Points() {
			if candidate.DeviceClass != nativepoints.DeviceClassController {
				continue
			}
			for _, location := range candidate.Locations {
				if location.Opcode != row.Opcode || location.Group != row.Group || location.Register != row.Register {
					continue
				}
				for _, family := range location.ControllerFamilies {
					if family == row.Family {
						matched = location.Codec == want && location.NativeUnit == unit
					}
				}
			}
		}
		if !matched {
			t.Errorf("matrix row %s 0x%02x/0x%04x %s: no location with codec %q unit %q", row.Family, row.Group, row.Register, row.Name, want, unit)
		}
	}
	if checked == 0 {
		t.Fatal("no matrix row with a stated value type")
	}
}

func nativeOf(t *testing.T, point nativepoints.Point, raw []byte) nativepoints.NativeValue {
	t.Helper()
	before := append([]byte(nil), raw...)
	value, err := nativepoints.DecodeNative(point, raw)
	if err != nil {
		t.Fatalf("DecodeNative(%s, % X): %v", point.Name, raw, err)
	}
	if !bytes.Equal(before, raw) {
		t.Fatal("DecodeNative modified the reply")
	}
	return value
}

func TestNativeDecodeBooleanEnumerationDateAndTime(t *testing.T) {
	cases := []struct {
		name  string
		raw   []byte
		class nativepoints.ValueClass
		value string
		label string
	}{
		{"adaptive_heating_curve_enabled", []byte{0x01}, nativepoints.ValueBoolean, "true", "Enabled"},
		{"adaptive_heating_curve_enabled", []byte{0x00}, nativepoints.ValueBoolean, "false", "Disabled"},
		{"dhw_legionella_protection_start_weekday", []byte{0x03, 0x00}, nativepoints.ValueEnumeration, "3", "Wednesday"},
		{"date", []byte{18, 11, 26}, nativepoints.ValueDate, "18.11.26", ""},
		{"clock_time", []byte{7, 5, 9}, nativepoints.ValueTime, "07:05:09", ""},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s % X", tc.name, tc.raw), func(t *testing.T) {
			value := nativeOf(t, mustLookup(t, nativepoints.DeviceClassController, tc.name), tc.raw)
			if value.ValueClass != tc.class || value.Native != tc.value || value.Label != tc.label {
				t.Errorf("got %+v, want class %s native %q label %q", value, tc.class, tc.value, tc.label)
			}
		})
	}
}

func TestNativeDecodeEveryDocumentedWireType(t *testing.T) {
	controller := func(name string) nativepoints.Point {
		return mustLookup(t, nativepoints.DeviceClassController, name)
	}
	cases := []struct {
		label  string
		point  nativepoints.Point
		raw    []byte
		class  nativepoints.ValueClass
		native string
		unit   string
	}{
		{"float32", controller("room_temperature_offset"), float32Reply(-1.5), nativepoints.ValueNumeric, "-1.5", "°C"},
		{"uint8 halves", controller("time_program_temperature_setpoint"), []byte{45}, nativepoints.ValueNumeric, "22.5", "°C"},
		{"uint8 count", controller("reception_level"), []byte{7}, nativepoints.ValueNumeric, "7", "count"},
		{"uint16", controller("preheating_time_maximum"), []byte{0x5A, 0x00}, nativepoints.ValueNumeric, "90", "min"},
		{"uint32", synthetic(nativepoints.CodecUint32LE, 1, "kWh"), []byte{0x01, 0x02, 0x03, 0x04}, nativepoints.ValueNumeric, "67305985", "kWh"},
		{"one byte of one or two", controller("zone_enabled"), []byte{0x01}, nativepoints.ValueBoolean, "true", ""},
		{"two bytes of one or two", controller("zone_enabled"), []byte{0x00, 0x00}, nativepoints.ValueBoolean, "false", ""},
		{"text", controller("zone_name"), []byte("Living\x00\x00"), nativepoints.ValueText, "Living", ""},
		{"text without terminator", controller("zone_name"), []byte("Hall"), nativepoints.ValueText, "Hall", ""},
		{"DATA2c", mustLookup(t, nativepoints.DeviceClassBoiler, "flow_temperature"), []byte{0x58, 0x02}, nativepoints.ValueNumeric, "37.5", "°C"},
		{"scaled UIN", mustLookup(t, nativepoints.DeviceClassBoiler, "dhw_water_flow_rate"), []byte{0x7B, 0x00}, nativepoints.ValueNumeric, "1.23", "L/min"},
	}
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			value := nativeOf(t, tc.point, tc.raw)
			if value.ValueClass != tc.class || value.Native != tc.native || value.NativeUnit != tc.unit {
				t.Errorf("got %+v, want class %s native %q unit %q", value, tc.class, tc.native, tc.unit)
			}
		})
	}
}

func TestNativeDecodeSignedEnumerationPerFamily(t *testing.T) {
	point := mustLookup(t, nativepoints.DeviceClassController, "zone_circuit_assignment")
	if point.Codec != nativepoints.CodecSignedByte || len(point.Enumeration) != 0 {
		t.Fatalf("families disagree on the labels, so the point has none: %+v", point)
	}
	if value := nativeOf(t, point, []byte{0xFF}); value.Native != "-1" || value.Label != "" {
		t.Errorf("point-level decode %+v", value)
	}
	cases := []struct {
		family string
		raw    byte
		native string
		label  string
	}{
		{"controller-720", 0xFF, "-1", "notBound"},
		{"controller-720", 0x08, "8", "Control_CC9"},
		{"controller-700", 0x02, "2", "Control_CC3"},
	}
	for _, tc := range cases {
		narrowed := locationForFamily(t, point, tc.family)
		if value := nativeOf(t, narrowed, []byte{tc.raw}); value.Native != tc.native || value.Label != tc.label || value.ValueClass != nativepoints.ValueEnumeration {
			t.Errorf("%s % X: %+v", tc.family, tc.raw, value)
		}
	}
	for _, tc := range []struct {
		family string
		raw    byte
	}{{"controller-700", 0x08}, {"controller-720", 0x09}} {
		value := nativeOf(t, locationForFamily(t, point, tc.family), []byte{tc.raw})
		if value.Native != fmt.Sprint(tc.raw) || value.Label != "" || value.ValueClass != nativepoints.ValueEnumeration {
			t.Errorf("%s % X without a candidate label: %+v", tc.family, tc.raw, value)
		}
	}
}

// The register map states that a code with no candidate label is not shown to
// be invalid on the wire, naming codes 2 and 9 of the zone special function.
func TestNativeDecodeKeepsACodeWithoutACandidateLabel(t *testing.T) {
	zone := mustLookup(t, nativepoints.DeviceClassController, "zone_special_function_state")
	for _, family := range controllerFamilies {
		narrowed := locationForFamily(t, zone, family)
		if len(narrowed.Enumeration) == 0 {
			t.Fatalf("%s location has no candidate labels: %+v", family, narrowed)
		}
		for _, raw := range [][]byte{{2}, {9}, {2, 0}, {9, 0}} {
			value := nativeOf(t, narrowed, raw)
			if value.ValueClass != nativepoints.ValueEnumeration || value.Native != fmt.Sprint(raw[0]) || value.Label != "" {
				t.Errorf("%s % X = %+v, want code %d without a label", family, raw, value, raw[0])
			}
		}
		if value := nativeOf(t, narrowed, []byte{7}); value.Label != "Holiday" {
			t.Errorf("%s labelled code 7 = %+v", family, value)
		}
	}
	dhw := mustLookup(t, nativepoints.DeviceClassController, "dhw_special_function_state")
	if value := nativeOf(t, dhw, []byte{2, 0}); value.Native != "2" || value.Label != "" || value.ValueClass != nativepoints.ValueEnumeration {
		t.Errorf("dhw_special_function_state 02 00 = %+v", value)
	}
	weekday := mustLookup(t, nativepoints.DeviceClassController, "dhw_legionella_protection_start_weekday")
	if value := nativeOf(t, weekday, []byte{9, 0}); value.Native != "9" || value.Label != "" {
		t.Errorf("dhw_legionella_protection_start_weekday 09 00 = %+v", value)
	}
}

func locationAt(t *testing.T, point nativepoints.Point, opcode, group uint8, register uint16) nativepoints.Point {
	t.Helper()
	for index, location := range point.Locations {
		if location.Opcode == opcode && location.Group == group && location.Register == register {
			narrowed, found := point.ForLocation(index)
			if !found {
				t.Fatalf("ForLocation(%d) not found", index)
			}
			return narrowed
		}
	}
	t.Fatalf("%s has no location 0x%02x/0x%02x/0x%04x", point.Name, opcode, group, register)
	return nativepoints.Point{}
}

func TestNativeDecodeLabelsTheMapStatesOutright(t *testing.T) {
	influence := mustLookup(t, nativepoints.DeviceClassController, "circuit_room_influence_mode")
	for raw, label := range map[byte]string{0: "inactive", 1: "active", 2: "extended"} {
		if value := nativeOf(t, influence, []byte{raw}); value.Native != fmt.Sprint(raw) || value.Label != label || value.ValueClass != nativepoints.ValueEnumeration {
			t.Errorf("circuit_room_influence_mode % X = %+v, want label %q", raw, value, label)
		}
	}
	address := mustLookup(t, nativepoints.DeviceClassController, "device_class_address")
	for _, group := range []uint8{0x09, 0x0A} {
		narrowed := locationAt(t, address, 0x06, group, 0x0002)
		for raw, label := range map[byte]string{0x15: "VRC720", 0x35: "VR92", 0x26: "VR71"} {
			if value := nativeOf(t, narrowed, []byte{raw}); value.Native != fmt.Sprint(raw) || value.Label != label {
				t.Errorf("group 0x%02x device_class_address % X = %+v, want label %q", group, raw, value, label)
			}
		}
	}
	if narrowed := locationAt(t, address, 0x06, 0x01, 0x0002); len(narrowed.Enumeration) != 0 {
		t.Errorf("a group whose map table states no label list has labels: %+v", narrowed.Enumeration)
	}
}

func TestNativeDecodeFirmwareVersionTriple(t *testing.T) {
	point := mustLookup(t, nativepoints.DeviceClassController, "device_firmware_version")
	for _, group := range []uint8{0x09, 0x0A, 0x0C} {
		narrowed := locationAt(t, point, 0x06, group, 0x0004)
		if narrowed.Codec != nativepoints.CodecMajorMinorPatch || narrowed.ValueClass != nativepoints.ValueVersion {
			t.Fatalf("group 0x%02x firmware version = codec %q class %q", group, narrowed.Codec, narrowed.ValueClass)
		}
		value := nativeOf(t, narrowed, []byte{0x08, 0x05, 0x00})
		if value.ValueClass != nativepoints.ValueVersion || value.Native != "08.05.00" || value.Label != "" {
			t.Errorf("group 0x%02x 08 05 00 = %+v", group, value)
		}
		if value := nativeOf(t, narrowed, []byte{0x10, 0x00, 0x63}); value.Native != "16.00.99" {
			t.Errorf("group 0x%02x byte-decimal 10 00 63 = %+v", group, value)
		}
		refuseNative(t, narrowed, []byte{0xFF, 0xFF, 0xFF}, nativepoints.DecodeSentinel)
		refuseNative(t, narrowed, []byte{0x08, 0x05}, nativepoints.DecodeShortReply)
		refuseNative(t, narrowed, []byte{0x08, 0x05, 0x00, 0x00}, nativepoints.DecodeLongReply)
	}
}

func locationForFamily(t *testing.T, point nativepoints.Point, family string) nativepoints.Point {
	t.Helper()
	for index, location := range point.Locations {
		for _, candidate := range location.ControllerFamilies {
			if candidate == family {
				narrowed, found := point.ForLocation(index)
				if !found || len(narrowed.Locations) != 1 || fmt.Sprint(narrowed.ControllerFamilies) != fmt.Sprint(location.ControllerFamilies) {
					t.Fatalf("ForLocation(%d) = %+v, %v", index, narrowed, found)
				}
				return narrowed
			}
		}
	}
	t.Fatalf("%s has no location for %s", point.Name, family)
	return nativepoints.Point{}
}

// A matrix row that states a value type and leaves the unit blank keeps the
// unit the map's Decode column states, on every family.
func TestNativeDecodeKeepsTheMapUnitBehindABlankMatrixUnit(t *testing.T) {
	cases := []struct {
		name string
		unit string
	}{
		{"manual_cooling_time", "d"},
		{"hydraulic_scheme_code", "count"},
		{"functional_module_configuration_code", "count"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			point := mustLookup(t, nativepoints.DeviceClassController, tc.name)
			narrowed := []nativepoints.Point{point}
			for _, family := range controllerFamilies {
				narrowed = append(narrowed, locationForFamily(t, point, family))
			}
			for _, candidate := range narrowed {
				value := nativeOf(t, candidate, []byte{0x02, 0x00})
				if value.ValueClass != nativepoints.ValueNumeric || value.Native != "2" || value.NativeUnit != tc.unit {
					t.Errorf("%v: got %+v, want numeric 2 %s", candidate.ControllerFamilies, value, tc.unit)
				}
			}
		})
	}
}

// A bare wire type does not make a register numeric when its Decode column
// states no class and its name is not a quantity; an identifier stays one.
func TestNativeUnclassifiedRegistersAreNotNumeric(t *testing.T) {
	for name, want := range map[string]nativepoints.ValueClass{
		"language": nativepoints.ValueUnknown, "network_identifier": nativepoints.ValueIdentifier,
		"hardware_identifier": nativepoints.ValueIdentifier,
	} {
		point := mustLookup(t, nativepoints.DeviceClassController, name)
		if point.ValueClass != want || point.Codec == nativepoints.CodecUnknown {
			t.Errorf("%s: class %q codec %q, want class %q with a codec", name, point.ValueClass, point.Codec, want)
		}
	}
}

func refuseNative(t *testing.T, point nativepoints.Point, raw []byte, failure nativepoints.DecodeFailure) {
	t.Helper()
	_, err := nativepoints.DecodeNative(point, raw)
	var decodeError *nativepoints.DecodeError
	if !errors.As(err, &decodeError) || decodeError.Failure != failure || decodeError.Name != point.Name || decodeError.Length != len(raw) {
		t.Errorf("DecodeNative(%s, % X) = %v, want %s", point.Name, raw, err, failure)
	}
}

func TestNativeDecodeRefusals(t *testing.T) {
	controller := func(name string) nativepoints.Point {
		return mustLookup(t, nativepoints.DeviceClassController, name)
	}
	nan := binary32(math.Float32bits(float32(math.NaN())))
	cases := []struct {
		label   string
		point   nativepoints.Point
		raw     []byte
		failure nativepoints.DecodeFailure
	}{
		{"float32 short", controller("room_temperature_offset"), []byte{0, 0, 0}, nativepoints.DecodeShortReply},
		{"float32 long", controller("room_temperature_offset"), []byte{0, 0, 0, 0, 0}, nativepoints.DecodeLongReply},
		{"float32 NaN", controller("room_temperature_offset"), nan, nativepoints.DecodeNotANumber},
		{"float32 infinite", controller("room_temperature_offset"), float32Reply(float32(math.Inf(-1))), nativepoints.DecodeInfinite},
		{"uint8 long", controller("reception_level"), []byte{1, 0}, nativepoints.DecodeLongReply},
		{"one or two empty", controller("zone_enabled"), nil, nativepoints.DecodeShortReply},
		{"one or two long", controller("zone_enabled"), []byte{1, 0, 0}, nativepoints.DecodeLongReply},
		{"boolean out of range", controller("adaptive_heating_curve_enabled"), []byte{2}, nativepoints.DecodeOutOfRange},
		{"date short", controller("date"), []byte{18, 11}, nativepoints.DecodeShortReply},
		{"date day", controller("date"), []byte{32, 11, 26}, nativepoints.DecodeOutOfRange},
		{"date month", controller("date"), []byte{18, 13, 26}, nativepoints.DecodeOutOfRange},
		{"date year", controller("date"), []byte{18, 11, 100}, nativepoints.DecodeOutOfRange},
		{"time long", controller("clock_time"), []byte{7, 5, 9, 0}, nativepoints.DecodeLongReply},
		{"text not UTF-8", controller("zone_name"), []byte{0xFF, 0xFE, 0x00}, nativepoints.DecodeInvalidText},
		{"DATA2c sentinel", mustLookup(t, nativepoints.DeviceClassBoiler, "flow_temperature"), []byte{0x00, 0x80}, nativepoints.DecodeSentinel},
		{"UIN sentinel", mustLookup(t, nativepoints.DeviceClassBoiler, "fan_speed"), []byte{0xFF, 0xFF}, nativepoints.DecodeSentinel},
		{"layout not documented", controller("heat_generator_status"), []byte{0x00}, nativepoints.DecodeCodecUnknown},
		{"zero scale", synthetic(nativepoints.CodecUint16LE, 0, "min"), []byte{1, 0}, nativepoints.DecodeCodecUnknown},
	}
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			refuseNative(t, tc.point, tc.raw, tc.failure)
		})
	}
}

func binary32(bits uint32) []byte {
	return []byte{byte(bits), byte(bits >> 8), byte(bits >> 16), byte(bits >> 24)}
}

func TestNativeDecodeKeepsTheSIExportRule(t *testing.T) {
	point := mustLookup(t, nativepoints.DeviceClassController, "adaptive_heating_curve_enabled")
	if point.Exportable || point.Withheld != nativepoints.WithheldNotNumeric {
		t.Fatalf("boolean row became exportable: %+v", point)
	}
	if _, err := nativepoints.Decode(point, []byte{1}); err == nil {
		t.Fatal("Decode exported a boolean")
	}
	if value := nativeOf(t, point, []byte{1}); value.Native != "true" {
		t.Fatalf("DecodeNative = %+v", value)
	}
}

func TestNativeForLocationRejectsAnIndexOutsideTheLocations(t *testing.T) {
	point := mustLookup(t, nativepoints.DeviceClassController, "device_connected")
	for _, index := range []int{-1, len(point.Locations)} {
		if _, found := point.ForLocation(index); found {
			t.Errorf("ForLocation(%d) found a location", index)
		}
	}
	narrowed, found := point.ForLocation(0)
	if !found {
		t.Fatal("ForLocation(0) not found")
	}
	narrowed.Locations[0].ControllerFamilies = append(narrowed.Locations[0].ControllerFamilies, "changed")
	if again := mustLookup(t, nativepoints.DeviceClassController, "device_connected"); len(again.Locations[0].ControllerFamilies) != 0 {
		t.Fatal("ForLocation shares mutable state")
	}
}

// A location whose own documentation states no Boolean or enumeration keeps
// its integer even when the point's name or a sibling location is Boolean or
// an enumeration.
func TestNativeDecodeFollowsTheClassOfEachLocation(t *testing.T) {
	connected := mustLookup(t, nativepoints.DeviceClassController, "device_connected")
	status := locationAt(t, connected, 0x06, 0x08, 0x0001)
	if value := nativeOf(t, status, []byte{2}); value.ValueClass != nativepoints.ValueNumeric || value.Native != "2" || value.Label != "" {
		t.Errorf("group 0x08 device_connected status byte 02 = %+v", value)
	}
	refuseNative(t, locationAt(t, connected, 0x06, 0x09, 0x0001), []byte{2}, nativepoints.DecodeOutOfRange)

	frost := mustLookup(t, nativepoints.DeviceClassController, "zone_frost_protection_active")
	if value := nativeOf(t, locationForFamily(t, frost, "controller-700"), []byte{2, 0}); value.ValueClass != nativepoints.ValueNumeric || value.Native != "2" {
		t.Errorf("controller-700 zone_frost_protection_active 02 00 = %+v", value)
	}
	narrowed := locationForFamily(t, frost, "controller-720")
	if value := nativeOf(t, narrowed, []byte{1}); value.ValueClass != nativepoints.ValueBoolean || value.Native != "true" || value.Label != "True" {
		t.Errorf("controller-720 zone_frost_protection_active 01 = %+v", value)
	}
	refuseNative(t, narrowed, []byte{2}, nativepoints.DecodeOutOfRange)

	address := mustLookup(t, nativepoints.DeviceClassController, "device_class_address")
	if value := nativeOf(t, locationAt(t, address, 0x06, 0x08, 0x0002), []byte{0x26}); value.ValueClass != nativepoints.ValueNumeric || value.Native != "38" || value.Label != "" {
		t.Errorf("group 0x08 device_class_address 26 = %+v", value)
	}
	if value := nativeOf(t, locationAt(t, address, 0x06, 0x09, 0x0002), []byte{0x26}); value.ValueClass != nativepoints.ValueEnumeration || value.Label != "VR71" {
		t.Errorf("group 0x09 device_class_address 26 = %+v", value)
	}
}

// Each group states its own pair of labels for device_connected.
func TestNativeDecodeDeviceConnectedLabelsOfEachGroup(t *testing.T) {
	connected := mustLookup(t, nativepoints.DeviceClassController, "device_connected")
	cases := []struct {
		group uint8
		raw   byte
		value string
		label string
	}{
		{0x09, 0, "false", "empty"}, {0x09, 1, "true", "paired"},
		{0x0A, 0, "false", "empty"}, {0x0A, 1, "true", "paired"},
		{0x0C, 0, "false", "not_connected"}, {0x0C, 1, "true", "connected"},
	}
	for _, tc := range cases {
		value := nativeOf(t, locationAt(t, connected, 0x06, tc.group, 0x0001), []byte{tc.raw})
		if value.ValueClass != nativepoints.ValueBoolean || value.Native != tc.value || value.Label != tc.label {
			t.Errorf("group 0x%02x device_connected % X = %+v, want %s %q", tc.group, tc.raw, value, tc.value, tc.label)
		}
	}
}

// Every opcode 0x02 group 0x09 location names the controller-720 family only,
// including the rows the register matrix does not list.
func TestNativeVentilationGroupNamesItsOneFamily(t *testing.T) {
	checked := 0
	for _, point := range nativepoints.Points() {
		for _, location := range point.Locations {
			if point.DeviceClass != nativepoints.DeviceClassController || location.Opcode != 0x02 || location.Group != 0x09 {
				continue
			}
			checked++
			if fmt.Sprint(location.ControllerFamilies) != "[controller-720]" {
				t.Errorf("%s at group 0x09 register 0x%04x names families %v", point.Name, location.Register, location.ControllerFamilies)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no opcode 0x02 group 0x09 location")
	}
	point := mustLookup(t, nativepoints.DeviceClassController, "ventilation_special_function_mode")
	if fmt.Sprint(point.ControllerFamilies) != "[controller-720]" {
		t.Errorf("ventilation_special_function_mode families %v", point.ControllerFamilies)
	}
}
