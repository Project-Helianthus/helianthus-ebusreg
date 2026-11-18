package nativepoints_test

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"

	catalog "github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/nativepoints"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/system"
)

var controllerFamilies = []string{"controller-700", "controller-720"}

func mustLookup(t *testing.T, class nativepoints.DeviceClass, name string) nativepoints.Point {
	t.Helper()
	point, found := nativepoints.Lookup(class, name)
	if !found {
		t.Fatalf("Lookup(%s, %s) not found", class, name)
	}
	return point
}

func TestPointKeysAreUniqueAndOrdered(t *testing.T) {
	points := nativepoints.Points()
	if len(points) == 0 {
		t.Fatal("empty catalog")
	}
	seen := map[string]bool{}
	for index, point := range points {
		if point.DeviceClass == "" || point.Name == "" || len(point.Locations) == 0 {
			t.Fatalf("incomplete row %d: %+v", index, point)
		}
		key := string(point.DeviceClass) + "\x00" + point.Name
		if seen[key] {
			t.Fatalf("duplicate (device class, name) %s/%s", point.DeviceClass, point.Name)
		}
		seen[key] = true
		if index > 0 {
			previous := points[index-1]
			if previous.DeviceClass > point.DeviceClass || (previous.DeviceClass == point.DeviceClass && previous.Name >= point.Name) {
				t.Fatalf("rows not ordered at %d: %s/%s after %s/%s", index, point.DeviceClass, point.Name, previous.DeviceClass, previous.Name)
			}
		}
	}
}

func TestNonFinalNamesAreNeverExportable(t *testing.T) {
	counts := map[nativepoints.NameStatus]int{}
	for _, point := range nativepoints.Points() {
		counts[point.NameStatus]++
		if point.NameStatus == nativepoints.NameFinal {
			continue
		}
		if point.Exportable || point.Withheld != nativepoints.WithheldNameNotFinal {
			t.Errorf("%s/%s status %s: exportable=%v withheld=%q", point.DeviceClass, point.Name, point.NameStatus, point.Exportable, point.Withheld)
		}
	}
	for _, status := range []nativepoints.NameStatus{nativepoints.NameFinal, nativepoints.NameProvisional, nativepoints.NameAmbiguous} {
		if counts[status] == 0 {
			t.Errorf("no row with name status %s", status)
		}
	}
}

type controllerAnnotation struct {
	Key struct {
		Operation uint8  `json:"operation"`
		Group     uint8  `json:"group"`
		Register  uint16 `json:"register"`
	} `json:"key"`
	Name       *string `json:"canonical_name"`
	Unresolved *bool   `json:"unresolved"`
	Wire       *string `json:"wire"`
	Decode     *string `json:"decode"`
}

type symbolicAnnotation struct {
	Class      string  `json:"class"`
	Operation  uint8   `json:"operation"`
	Selector   *string `json:"selector"`
	Name       string  `json:"name"`
	Unresolved bool    `json:"unresolved"`
	Unit       *string `json:"unit"`
}

func controllerAnnotations(t *testing.T) []controllerAnnotation {
	t.Helper()
	raw, err := catalog.B524NameAnnotationEntries()
	if err != nil {
		t.Fatal(err)
	}
	var rows []controllerAnnotation
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatal(err)
	}
	return rows
}

func symbolicAnnotations(t *testing.T) []symbolicAnnotation {
	t.Helper()
	raw, err := catalog.B509SymbolicAnnotationEntries()
	if err != nil {
		t.Fatal(err)
	}
	var rows []symbolicAnnotation
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatal(err)
	}
	return rows
}

func selectorText(selector string) string {
	digits := strings.ToUpper(strings.ReplaceAll(strings.TrimPrefix(selector, "0x"), " ", ""))
	var pairs []string
	for index := 0; index+2 <= len(digits); index += 2 {
		pairs = append(pairs, digits[index:index+2])
	}
	return strings.Join(pairs, " ")
}

func TestEveryAnnotatedNameHasItsLocationAndStatus(t *testing.T) {
	type location struct {
		class    nativepoints.DeviceClass
		name     string
		opcode   uint8
		group    uint8
		register uint16
		selector string
	}
	want := map[location]bool{}
	unresolved := map[string]bool{}
	for _, row := range controllerAnnotations(t) {
		if row.Name == nil {
			continue
		}
		want[location{nativepoints.DeviceClassController, *row.Name, row.Key.Operation, row.Key.Group, row.Key.Register, ""}] = true
		if row.Unresolved == nil || *row.Unresolved {
			unresolved[string(nativepoints.DeviceClassController)+"/"+*row.Name] = true
		}
	}
	for _, row := range symbolicAnnotations(t) {
		if row.Class == "any" {
			continue
		}
		want[location{nativepoints.DeviceClass(row.Class), row.Name, row.Operation, 0, 0, selectorText(*row.Selector)}] = true
		if row.Unresolved {
			unresolved[row.Class+"/"+row.Name] = true
		}
	}
	got := map[location]bool{}
	for _, point := range nativepoints.Points() {
		for _, place := range point.Locations {
			got[location{point.DeviceClass, point.Name, place.Opcode, place.Group, place.Register, place.Selector}] = true
		}
		if unresolved[string(point.DeviceClass)+"/"+point.Name] && point.NameStatus == nativepoints.NameFinal {
			t.Errorf("%s/%s is unresolved in its annotation but final in the catalog", point.DeviceClass, point.Name)
		}
	}
	for place := range want {
		if !got[place] {
			t.Errorf("annotated location missing: %+v", place)
		}
	}
	for place := range got {
		if !want[place] {
			t.Errorf("catalog location without annotation: %+v", place)
		}
	}
}

func TestCodesAndIdentifiersAreNeverExportable(t *testing.T) {
	for _, point := range nativepoints.Points() {
		codeLike := strings.HasSuffix(point.Name, "_code") || strings.Contains(point.Name, "access_code") ||
			strings.Contains(point.Name, "menu_code") || strings.Contains(point.Name, "phone_number") ||
			strings.HasSuffix(point.Name, "_number")
		if point.ValueClass == nativepoints.ValueCode || point.ValueClass == nativepoints.ValueIdentifier || codeLike {
			if point.Exportable {
				t.Errorf("code or identifier %s/%s is exportable", point.DeviceClass, point.Name)
			}
		}
		if codeLike && point.NameStatus == nativepoints.NameFinal && point.Withheld != nativepoints.WithheldIdentifierOrCode {
			t.Errorf("%s/%s withheld %q, want %q", point.DeviceClass, point.Name, point.Withheld, nativepoints.WithheldIdentifierOrCode)
		}
	}
	for _, key := range []struct {
		class nativepoints.DeviceClass
		name  string
	}{
		{nativepoints.DeviceClassController, "installer_access_code"},
		{nativepoints.DeviceClassController, "hydraulic_scheme_code"},
		{nativepoints.DeviceClassController, "installer_phone_number_prefix"},
		{nativepoints.DeviceClassBoiler, "installer_phone_number"},
		{nativepoints.DeviceClassBoiler, "device_specific_number"},
		{nativepoints.DeviceClassBoiler, "boiler_state_code"},
	} {
		point := mustLookup(t, key.class, key.name)
		if point.Exportable || (point.ValueClass != nativepoints.ValueCode && point.ValueClass != nativepoints.ValueIdentifier) {
			t.Errorf("%s/%s: exportable=%v class=%s", key.class, key.name, point.Exportable, point.ValueClass)
		}
	}
}

func TestExportableRowsAreComplete(t *testing.T) {
	exportable := 0
	for _, point := range nativepoints.Points() {
		if !point.Exportable {
			if point.Withheld == nativepoints.WithheldNone {
				t.Errorf("%s/%s withheld without a reason", point.DeviceClass, point.Name)
			}
			continue
		}
		exportable++
		if point.Withheld != nativepoints.WithheldNone || point.NameStatus != nativepoints.NameFinal ||
			point.ValueClass != nativepoints.ValueNumeric || point.Codec == nativepoints.CodecUnknown ||
			point.NativeUnit == "" || point.NativeUnit == "bar" || point.ScaleDivisor == 0 {
			t.Errorf("exportable row incomplete: %+v", point)
		}
		raw := make([]byte, widthOf(point.Codec))
		if _, err := nativepoints.Decode(point, raw); err != nil {
			t.Errorf("exportable %s/%s does not decode a zero reply: %v", point.DeviceClass, point.Name, err)
		}
	}
	if exportable == 0 {
		t.Fatal("no exportable row")
	}
}

func widthOf(codec nativepoints.Codec) int {
	if codec == nativepoints.CodecFloat32LE {
		return 4
	}
	return 2
}

func TestHeatPumpRowsKeepTheirRealStatus(t *testing.T) {
	cases := []struct {
		class    nativepoints.DeviceClass
		name     string
		unit     string
		withheld nativepoints.WithheldReason
	}{
		{nativepoints.DeviceClassHeatPump, "refrigerant_high_pressure", "bar", nativepoints.WithheldPressureReferenceUnknown},
		{nativepoints.DeviceClassHeatPump, "refrigerant_low_pressure", "bar", nativepoints.WithheldPressureReferenceUnknown},
		{nativepoints.DeviceClassHeatPumpInterfaceUnit, "refrigerant_high_pressure", "bar", nativepoints.WithheldPressureReferenceUnknown},
		{nativepoints.DeviceClassHeatPump, "refrigerant_low_pressure_maximum", "", nativepoints.WithheldNameNotFinal},
		{nativepoints.DeviceClassHeatPump, "refrigerant_superheat", "K", nativepoints.WithheldCodecUnknown},
		{nativepoints.DeviceClassHeatPump, "refrigerant_subcooling", "K", nativepoints.WithheldCodecUnknown},
		{nativepoints.DeviceClassHeatPumpInterfaceUnit, "refrigerant_superheat", "K", nativepoints.WithheldNameNotFinal},
		{nativepoints.DeviceClassHeatPump, "compressor_dhw_start_count", "count", nativepoints.WithheldCodecUnknown},
		{nativepoints.DeviceClassHeatPump, "compressor_minimum_on_duration", "s", nativepoints.WithheldCodecUnknown},
		{nativepoints.DeviceClassHeatPumpInterfaceUnit, "compressor_start_count", "count", nativepoints.WithheldCodecUnknown},
		{nativepoints.DeviceClassHeatPumpInterfaceUnit, "compressor_operating_time", "h", nativepoints.WithheldCodecUnknown},
		{nativepoints.DeviceClassHeatPumpInterfaceUnit, "compressor_speed", "", nativepoints.WithheldNameNotFinal},
		{nativepoints.DeviceClassController, "ventilation_fan_speed_limit_maximum_at_day", "", nativepoints.WithheldUnitUnknown},
		{nativepoints.DeviceClassController, "water_pressure", "bar", nativepoints.WithheldPressureReferenceUnknown},
		{nativepoints.DeviceClassBoiler, "water_pressure", "bar", nativepoints.WithheldPressureReferenceUnknown},
	}
	for _, tc := range cases {
		point := mustLookup(t, tc.class, tc.name)
		if point.Exportable || point.Withheld != tc.withheld || point.NativeUnit != tc.unit {
			t.Errorf("%s/%s: exportable=%v withheld=%q unit=%q; want withheld %q unit %q", tc.class, tc.name, point.Exportable, point.Withheld, point.NativeUnit, tc.withheld, tc.unit)
		}
	}
}

func TestKnownExportableRows(t *testing.T) {
	cases := []struct {
		class    nativepoints.DeviceClass
		name     string
		codec    nativepoints.Codec
		divisor  uint16
		unit     string
		families []string
	}{
		{nativepoints.DeviceClassController, "outdoor_temperature", nativepoints.CodecFloat32LE, 1, "°C", controllerFamilies},
		{nativepoints.DeviceClassController, "preheating_time_maximum", nativepoints.CodecUint16LE, 1, "min", []string{"controller-720"}},
		{nativepoints.DeviceClassController, "zone_room_humidity", nativepoints.CodecFloat32LE, 1, "%", controllerFamilies},
		{nativepoints.DeviceClassBoiler, "flow_temperature", nativepoints.CodecData2c, 1, "°C", nil},
		{nativepoints.DeviceClassBoiler, "dhw_water_flow_rate", nativepoints.CodecUIN, 100, "L/min", nil},
		{nativepoints.DeviceClassBoiler, "fan_speed", nativepoints.CodecUIN, 1, "rpm", nil},
	}
	for _, tc := range cases {
		point := mustLookup(t, tc.class, tc.name)
		if !point.Exportable || point.Codec != tc.codec || point.ScaleDivisor != tc.divisor || point.NativeUnit != tc.unit ||
			fmt.Sprint(point.ControllerFamilies) != fmt.Sprint(tc.families) {
			t.Errorf("%s/%s = %+v", tc.class, tc.name, point)
		}
	}
}

type probe struct {
	codec string
	unit  string
}

func probeController(family string, group uint8, register uint16) probe {
	present := func(value string) system.B524SystemScalarField[string] {
		return system.B524SystemScalarField[string]{B524SystemScalarPresent: true, B524SystemScalarValue: value}
	}
	switch group {
	case 0x00:
		result := system.B524SystemInterpretScalars(system.B524SystemScalarInput{B524SystemScalarContext: system.B524SystemScalarContext{
			B524SystemScalarFamily: present(family), B524SystemScalarSoftware: present("s"), B524SystemScalarHardware: present("h"),
			B524SystemScalarAPIVersion: present("v"), B524SystemScalarAPIRevision: present("r"),
			B524SystemScalarOpcode:   system.B524SystemScalarField[byte]{B524SystemScalarPresent: true, B524SystemScalarValue: 0x02},
			B524SystemScalarGroup:    system.B524SystemScalarField[byte]{B524SystemScalarPresent: true, B524SystemScalarValue: group},
			B524SystemScalarInstance: system.B524SystemScalarField[byte]{B524SystemScalarPresent: true},
			B524SystemScalarRegister: system.B524SystemScalarField[uint16]{B524SystemScalarPresent: true, B524SystemScalarValue: register},
		}})
		return probe{result.B524SystemScalarCodec, result.B524SystemScalarDeclaredUnit}
	case 0x01, 0x09:
		text := func(value string) system.B524DHWVentilationScalarField[string] {
			return system.B524DHWVentilationScalarField[string]{B524DHWVentilationScalarPresent: true, B524DHWVentilationScalarValue: value}
		}
		result := system.B524DHWVentilationInterpretScalars(system.B524DHWVentilationScalarInput{B524DHWVentilationScalarContext: system.B524DHWVentilationScalarContext{
			B524DHWVentilationScalarFamily: text(family), B524DHWVentilationScalarSoftware: text("s"), B524DHWVentilationScalarHardware: text("h"),
			B524DHWVentilationScalarAPIVersion: text("v"), B524DHWVentilationScalarAPIRevision: text("r"),
			B524DHWVentilationScalarOpcode:   system.B524DHWVentilationScalarField[byte]{B524DHWVentilationScalarPresent: true, B524DHWVentilationScalarValue: 0x02},
			B524DHWVentilationScalarGroup:    system.B524DHWVentilationScalarField[byte]{B524DHWVentilationScalarPresent: true, B524DHWVentilationScalarValue: group},
			B524DHWVentilationScalarInstance: system.B524DHWVentilationScalarField[byte]{B524DHWVentilationScalarPresent: true},
			B524DHWVentilationScalarRegister: system.B524DHWVentilationScalarField[uint16]{B524DHWVentilationScalarPresent: true, B524DHWVentilationScalarValue: register},
		}})
		return probe{result.B524DHWVentilationScalarCodec, result.B524DHWVentilationScalarDeclaredUnit}
	case 0x02:
		text := func(value string) system.B524CircuitScalarField[string] {
			return system.B524CircuitScalarField[string]{B524CircuitScalarPresent: true, B524CircuitScalarValue: value}
		}
		result := system.B524InterpretCircuitScalar(system.B524CircuitScalarInput{B524CircuitScalarContext: system.B524CircuitScalarContext{
			B524CircuitScalarFamily: text(family), B524CircuitScalarSoftware: text("s"), B524CircuitScalarHardware: text("h"),
			B524CircuitScalarAPIVersion: text("v"), B524CircuitScalarAPIRevision: text("r"),
			B524CircuitScalarOpcode:   system.B524CircuitScalarField[byte]{B524CircuitScalarPresent: true, B524CircuitScalarValue: 0x02},
			B524CircuitScalarGroup:    system.B524CircuitScalarField[byte]{B524CircuitScalarPresent: true, B524CircuitScalarValue: group},
			B524CircuitScalarInstance: system.B524CircuitScalarField[byte]{B524CircuitScalarPresent: true},
			B524CircuitScalarRegister: system.B524CircuitScalarField[uint16]{B524CircuitScalarPresent: true, B524CircuitScalarValue: register},
		}})
		return probe{result.B524CircuitScalarCodec, result.B524CircuitScalarDeclaredUnit}
	case 0x03:
		text := func(value string) system.B524ZoneScalarField[string] {
			return system.B524ZoneScalarField[string]{B524ZoneScalarPresent: true, B524ZoneScalarValue: value}
		}
		result := system.B524ZoneScalarInterpret(system.B524ZoneScalarInput{B524ZoneScalarContext: system.B524ZoneScalarContext{
			B524ZoneScalarFamily: text(family), B524ZoneScalarSoftware: text("s"), B524ZoneScalarHardware: text("h"),
			B524ZoneScalarAPIVersion: text("v"), B524ZoneScalarAPIRevision: text("r"),
			B524ZoneScalarOpcode:   system.B524ZoneScalarField[byte]{B524ZoneScalarPresent: true, B524ZoneScalarValue: 0x02},
			B524ZoneScalarGroup:    system.B524ZoneScalarField[byte]{B524ZoneScalarPresent: true, B524ZoneScalarValue: group},
			B524ZoneScalarInstance: system.B524ZoneScalarField[byte]{B524ZoneScalarPresent: true},
			B524ZoneScalarRegister: system.B524ZoneScalarField[uint16]{B524ZoneScalarPresent: true, B524ZoneScalarValue: register},
		}})
		return probe{result.B524ZoneScalarCodec, result.B524ZoneScalarUnit}
	}
	return probe{}
}

func TestControllerCodecsMatchTheScalarInterpreters(t *testing.T) {
	type place struct {
		group    uint8
		register uint16
	}
	expected := map[place]map[string]probe{}
	for _, family := range controllerFamilies {
		for _, group := range []uint8{0x00, 0x01, 0x02, 0x03, 0x09} {
			for register := uint16(0); register <= 0x00FF; register++ {
				result := probeController(family, group, register)
				if result.codec == "" {
					continue
				}
				if result.unit == "Unknown" {
					result.unit = ""
				}
				key := place{group, register}
				if expected[key] == nil {
					expected[key] = map[string]probe{}
				}
				expected[key][family] = result
			}
		}
	}
	if len(expected) == 0 {
		t.Fatal("interpreter probe found no register")
	}
	located := map[place]bool{}
	for _, point := range nativepoints.Points() {
		if point.DeviceClass != nativepoints.DeviceClassController {
			continue
		}
		pointFamilies := map[string]bool{}
		covered := map[place]map[string]bool{}
		for _, location := range point.Locations {
			for _, family := range location.ControllerFamilies {
				pointFamilies[family] = true
			}
			if location.Opcode != 0x02 {
				continue
			}
			key := place{location.Group, location.Register}
			located[key] = true
			if covered[key] == nil {
				covered[key] = map[string]bool{}
			}
			for _, family := range location.ControllerFamilies {
				covered[key][family] = true
				result, interpreted := expected[key][family]
				if interpreted && (location.Codec != nativepoints.Codec(result.codec) || location.NativeUnit != result.unit) {
					t.Errorf("%s at %+v for %s: codec %q unit %q, interpreter %q %q", point.Name, key, family, location.Codec, location.NativeUnit, result.codec, result.unit)
				}
			}
		}
		for key, families := range covered {
			for family := range expected[key] {
				if !families[family] {
					t.Errorf("%s at %+v: interpreted family %s has no location", point.Name, key, family)
				}
			}
		}
		var want []string
		for family := range pointFamilies {
			want = append(want, family)
		}
		sort.Strings(want)
		if fmt.Sprint(point.ControllerFamilies) != fmt.Sprint(want) {
			t.Errorf("%s: families %v, want %v", point.Name, point.ControllerFamilies, want)
		}
	}
	annotated := map[place]bool{}
	for _, row := range controllerAnnotations(t) {
		if row.Name != nil && row.Key.Operation == 0x02 {
			annotated[place{row.Key.Group, row.Key.Register}] = true
		}
	}
	for key := range expected {
		if annotated[key] && !located[key] {
			t.Errorf("interpreted and named register %+v has no catalog location", key)
		}
	}
}

var mapWireCodecs = map[string]nativepoints.Codec{
	"u8": nativepoints.CodecUint8, "i8": nativepoints.CodecSignedByte, "u16": nativepoints.CodecUint16LE,
	"u32": nativepoints.CodecUint32LE, "f32": nativepoints.CodecFloat32LE, "string": nativepoints.CodecText,
}

var mapDecodeUnits = map[string]struct {
	unit    string
	divisor uint16
}{
	"°C": {"°C", 1}, "°C/2": {"°C", 2}, "K": {"K", 1}, "bar": {"bar", 1}, "%": {"%", 1},
	"kWh": {"kWh", 1}, "hrs": {"h", 1}, "min": {"min", 1}, "count": {"count", 1}, "days": {"d", 1},
}

func matrixLabels(t *testing.T, text string) []nativepoints.EnumerationLabel {
	t.Helper()
	if text == "" || strings.HasPrefix(text, "date:") || strings.HasPrefix(text, "time:") {
		return nil
	}
	var labels []nativepoints.EnumerationLabel
	for _, part := range strings.Split(strings.TrimPrefix(text, "candidate labels: "), "; ") {
		code, label, found := strings.Cut(part, "=")
		var value int64
		if _, err := fmt.Sscan(code, &value); !found || err != nil {
			t.Fatalf("malformed matrix enumeration %q", text)
		}
		labels = append(labels, nativepoints.EnumerationLabel{Value: value, Label: label})
	}
	sort.Slice(labels, func(i, j int) bool { return labels[i].Value < labels[j].Value })
	return labels
}

type layoutKey struct {
	opcode, group uint8
	register      uint16
}

type documentedLayout struct {
	codec   nativepoints.Codec
	unit    string
	divisor uint16
	labels  string
}

func matrixLayouts(t *testing.T) map[layoutKey]map[string]documentedLayout {
	t.Helper()
	matrix := map[layoutKey]map[string]documentedLayout{}
	for _, row := range registerMatrix(t) {
		codec := matrixCodec(row)
		if codec == nativepoints.CodecUnknown {
			continue
		}
		key := layoutKey{row.Opcode, row.Group, row.Register}
		if matrix[key] == nil {
			matrix[key] = map[string]documentedLayout{}
		}
		matrix[key][row.Family] = documentedLayout{codec, row.Unit, 1, fmt.Sprint(matrixLabels(t, row.Enumeration))}
	}
	return matrix
}

func mapLayouts(t *testing.T) map[layoutKey]documentedLayout {
	t.Helper()
	documented := map[layoutKey]documentedLayout{}
	stated := 0
	for _, row := range controllerAnnotations(t) {
		if row.Name == nil {
			continue
		}
		want := documentedLayout{divisor: 1, labels: fmt.Sprint([]nativepoints.EnumerationLabel(nil))}
		if row.Wire != nil {
			want.codec = mapWireCodecs[*row.Wire]
			if want.codec != nativepoints.CodecUnknown {
				stated++
			}
		}
		if row.Decode != nil {
			if unit, found := mapDecodeUnits[*row.Decode]; found {
				want.unit, want.divisor = unit.unit, unit.divisor
			}
		}
		documented[layoutKey{row.Key.Operation, row.Key.Group, row.Key.Register}] = want
	}
	if stated == 0 {
		t.Fatal("no annotation carries a map wire type with a layout")
	}
	return documented
}

func TestControllerLayoutsFollowTheMatrixAndTheMap(t *testing.T) {
	matrix, documented := matrixLayouts(t), mapLayouts(t)
	located := map[layoutKey]map[string]bool{}
	checked := 0
	for _, point := range nativepoints.Points() {
		if point.DeviceClass != nativepoints.DeviceClassController {
			continue
		}
		for _, location := range point.Locations {
			key := layoutKey{location.Opcode, location.Group, location.Register}
			got := documentedLayout{location.Codec, location.NativeUnit, location.ScaleDivisor, fmt.Sprint(location.Enumeration)}
			if documented[key].codec != nativepoints.CodecUnknown && location.Codec == nativepoints.CodecUnknown {
				t.Errorf("%s at %+v: the map states the wire type of codec %q but the location has none", point.Name, key, documented[key].codec)
			}
			if len(location.ControllerFamilies) == 0 {
				if len(matrix[key]) != 0 || got != documented[key] {
					t.Errorf("%s at %+v: family-independent layout %+v, map %+v, matrix families %d", point.Name, key, got, documented[key], len(matrix[key]))
				}
				checked++
				continue
			}
			if located[key] == nil {
				located[key] = map[string]bool{}
			}
			for _, family := range location.ControllerFamilies {
				located[key][family] = true
				want, inMatrix := matrix[key][family]
				if !inMatrix {
					want = documented[key]
				}
				if got != want {
					t.Errorf("%s at %+v for %s: layout %+v, documented %+v", point.Name, key, family, got, want)
				}
				checked++
			}
		}
	}
	for key, families := range matrix {
		if _, annotated := documented[key]; !annotated {
			continue
		}
		for family := range families {
			if !located[key][family] {
				t.Errorf("matrix row %s %+v has no location", family, key)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no controller location checked")
	}
}

var opcode06Slots = map[uint8]bool{0x01: true, 0x02: true, 0x03: true, 0x04: true, 0x05: true, 0x06: true, 0x07: true,
	0x08: true, 0x09: true, 0x0A: true, 0x0B: true, 0x0C: true, 0x0D: true, 0x0E: true, 0x0F: true}

func TestControllerInstanceRangesMatchTheDiscoveryProfiles(t *testing.T) {
	type scope struct{ group, opcode uint8 }
	profiles := map[scope]system.B524GroupProfile{}
	for _, profile := range system.B524DiscoveryProfiles {
		profiles[scope{profile.Group, profile.Opcode}] = profile
	}
	for _, point := range nativepoints.Points() {
		for _, location := range point.Locations {
			if point.DeviceClass != nativepoints.DeviceClassController {
				if location.InstanceRule != nativepoints.InstanceSingle || location.InstanceMax != 0 {
					t.Errorf("%s/%s: register service row is instanced: %+v", point.DeviceClass, point.Name, location)
				}
				continue
			}
			if location.Opcode == 0x06 {
				if !opcode06Slots[location.Group] || location.InstanceRule != nativepoints.InstanceIndexed || location.InstanceMin != 1 || location.InstanceMax != 8 {
					t.Errorf("%s: opcode 0x06 group 0x%02x has %+v, the map's slot interval is 1 to 8", point.Name, location.Group, location)
				}
				continue
			}
			profile, found := profiles[scope{location.Group, location.Opcode}]
			switch {
			case !found:
				if location.InstanceRule != nativepoints.InstanceUnknown {
					t.Errorf("%s: unprofiled group 0x%02x has rule %s", point.Name, location.Group, location.InstanceRule)
				}
			case profile.InstanceMax == 0:
				if location.InstanceRule != nativepoints.InstanceSingle || location.InstanceMax != 0 {
					t.Errorf("%s: singleton group 0x%02x has %+v", point.Name, location.Group, location)
				}
			default:
				if location.InstanceRule != nativepoints.InstanceIndexed || location.InstanceMin != 0 || location.InstanceMax != profile.InstanceMax {
					t.Errorf("%s: indexed group 0x%02x has %+v, profile max 0x%02x", point.Name, location.Group, location, profile.InstanceMax)
				}
			}
		}
	}
}

func TestRegisterServiceCodecsMatchTheDictionaries(t *testing.T) {
	declarations := map[string]string{}
	boiler, err := catalog.B509BoilerCodecDictionaryEntries()
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range boiler {
		selector := entry.B509BoilerCodecDictionaryKey.B509BoilerCodecDictionarySelector
		declarations[fmt.Sprintf("boiler/%02X %02X", selector[0], selector[1])] = entry.B509BoilerCodecDeclaration
	}
	heatPump, err := catalog.B509HeatpumpCodecDictionaryEntries()
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range heatPump {
		selector := entry.B509HeatpumpCodecDictionaryKey.B509HeatpumpCodecDictionarySelector
		declarations[fmt.Sprintf("heat_pump/%02X %02X", selector[0], selector[1])] = entry.B509HeatpumpCodecDeclaration
	}
	units := map[string]string{}
	for _, row := range symbolicAnnotations(t) {
		if row.Class != "any" && row.Unit != nil {
			units[row.Class+"/"+selectorText(*row.Selector)] = *row.Unit
		}
	}
	matched := 0
	for _, point := range nativepoints.Points() {
		if point.DeviceClass == nativepoints.DeviceClassController {
			continue
		}
		if point.Primary != 0xB5 || point.Secondary != 0x09 {
			t.Errorf("%s/%s service %02x%02x", point.DeviceClass, point.Name, point.Primary, point.Secondary)
		}
		for _, location := range point.Locations {
			key := string(point.DeviceClass) + "/" + location.Selector
			if declarations[key] != location.CodecDeclaration {
				t.Errorf("%s: declaration %q, dictionary %q", key, location.CodecDeclaration, declarations[key])
			}
			if units[key] != location.NativeUnit {
				t.Errorf("%s: unit %q, annotation %q", key, location.NativeUnit, units[key])
			}
			if declarations[key] != "" {
				matched++
			}
		}
	}
	if matched != len(declarations) {
		t.Errorf("joined %d dictionary declarations of %d", matched, len(declarations))
	}
}

func TestPointsAndLookupReturnCopies(t *testing.T) {
	point := mustLookup(t, nativepoints.DeviceClassController, "outdoor_temperature")
	point.Locations[0].ControllerFamilies[0] = "changed"
	point.ControllerFamilies[0] = "changed"
	point.Locations[0].Register = 0
	points := nativepoints.Points()
	points[0].Name = "changed"
	again := mustLookup(t, nativepoints.DeviceClassController, "outdoor_temperature")
	if again.Locations[0].ControllerFamilies[0] == "changed" || again.ControllerFamilies[0] == "changed" || again.Locations[0].Register == 0 {
		t.Fatal("Lookup shares mutable state")
	}
	if nativepoints.Points()[0].Name == "changed" {
		t.Fatal("Points shares mutable state")
	}
	if _, found := nativepoints.Lookup(nativepoints.DeviceClassBoiler, "outdoor_temperature"); found {
		t.Fatal("Lookup ignored the device class")
	}
}
