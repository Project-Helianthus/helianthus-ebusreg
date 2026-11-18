package nativepoints

// DeviceClass is the class of bus participant that answers a point's request.
type DeviceClass string

// Device classes of the catalog rows.
const (
	DeviceClassController            DeviceClass = "controller"
	DeviceClassBoiler                DeviceClass = "boiler"
	DeviceClassHeatPump              DeviceClass = "heat_pump"
	DeviceClassHeatPumpInterfaceUnit DeviceClass = "heat_pump_interface_unit"
)

// NameStatus says whether a canonical name identifies one register.
type NameStatus string

// Name statuses: final names are resolved at every location and their
// locations apply to disjoint controller families; provisional names are
// unresolved at one location at least; ambiguous names are resolved but name
// several locations that may apply to the same controller.
const (
	NameFinal       NameStatus = "final"
	NameProvisional NameStatus = "provisional"
	NameAmbiguous   NameStatus = "ambiguous"
)

// ValueClass is the kind of value a point carries.
type ValueClass string

// Value classes; only numeric values can be exported.
const (
	ValueNumeric     ValueClass = "numeric"
	ValueBoolean     ValueClass = "boolean"
	ValueEnumeration ValueClass = "enumeration"
	ValueFlags       ValueClass = "flags"
	ValueCode        ValueClass = "code"
	ValueIdentifier  ValueClass = "identifier"
	ValueText        ValueClass = "text"
	ValueDate        ValueClass = "date"
	ValueTime        ValueClass = "time"
	ValueCommand     ValueClass = "command"
	ValueUnknown     ValueClass = "unknown"
)

// Codec is a machine-readable value encoding. DecodeNative implements every
// codec; Decode implements the float32, uint16, DATA2c, DATA2b and UIN codecs.
type Codec string

// Codecs: little-endian float32, uint16 and uint32, uint8, int8 and the
// unsigned integer of one or two little-endian bytes of the controller
// registers; three plain bytes of day, month and year and of hour, minute and
// second; a NUL-terminated text; and the signed 16-bit DATA2c (1/16) and
// DATA2b (1/256) and the unsigned 16-bit UIN word of the register service.
// CodecUnknown marks a missing, literal-only or undocumented layout.
const (
	CodecUnknown          Codec = ""
	CodecFloat32LE        Codec = "float32_le"
	CodecUint16LE         Codec = "uint16_le"
	CodecUint32LE         Codec = "uint32_le"
	CodecUint8            Codec = "uint8"
	CodecSignedByte       Codec = "int8"
	CodecUint8OrUint16LE  Codec = "uint8_or_uint16_le"
	CodecDayMonthYear     Codec = "day_month_year"
	CodecHourMinuteSecond Codec = "hour_minute_second"
	CodecText             Codec = "text"
	CodecData2c           Codec = "DATA2c"
	CodecData2b           Codec = "DATA2b"
	CodecUIN              Codec = "UIN"
)

// EnumerationLabel is one documented code of an enumeration or Boolean value
// and its label.
type EnumerationLabel struct {
	Value int64  `json:"value"`
	Label string `json:"label"`
}

// InstanceRule says how a location is instanced.
type InstanceRule string

// Instance rules: single has instance 0 only or no instance byte, indexed
// spans InstanceMin to InstanceMax, unknown has no documented instance range.
const (
	InstanceSingle  InstanceRule = "single"
	InstanceIndexed InstanceRule = "indexed"
	InstanceUnknown InstanceRule = "unknown"
)

// WithheldReason is the first rule a non-exportable point fails.
type WithheldReason string

// Withheld reasons, in the order they are checked. WithheldNone marks an
// exportable point.
const (
	WithheldNone                     WithheldReason = ""
	WithheldNameNotFinal             WithheldReason = "name_not_final"
	WithheldIdentifierOrCode         WithheldReason = "identifier_or_code"
	WithheldNotNumeric               WithheldReason = "not_numeric"
	WithheldUnitUnknown              WithheldReason = "unit_unknown"
	WithheldPressureReferenceUnknown WithheldReason = "pressure_reference_unknown"
	WithheldCodecUnknown             WithheldReason = "codec_unknown"
)

// Location is one request that reads a point. Controller locations carry the
// group and register; register-service locations carry the selector bytes
// the symbolic annotation records, as upper-case hexadecimal pairs separated
// by spaces. Codec,
// declaration, scale, unit and enumeration are those declared for this
// location; ControllerFamilies lists the controller families whose register
// matrix row or scalar interpreter declares them, and is empty when only the
// family-independent register map declares them.
type Location struct {
	Opcode             uint8              `json:"opcode"`
	Group              uint8              `json:"group"`
	Register           uint16             `json:"register"`
	Selector           string             `json:"selector"`
	InstanceRule       InstanceRule       `json:"instance_rule"`
	InstanceMin        uint8              `json:"instance_min"`
	InstanceMax        uint8              `json:"instance_max"`
	ControllerFamilies []string           `json:"controller_families"`
	Codec              Codec              `json:"codec"`
	CodecDeclaration   string             `json:"codec_declaration"`
	ScaleDivisor       uint16             `json:"scale_divisor"`
	NativeUnit         string             `json:"native_unit"`
	Enumeration        []EnumerationLabel `json:"enumeration"`
}

// Point is one catalog row, keyed by device class and canonical name. Codec,
// scale and native unit are set only when every location agrees on them, and
// the enumeration only when every location also agrees on it; the native value
// is the decoded integer or float divided by ScaleDivisor.
// ControllerFamilies is the union of the location families. Exportable is
// true only when Withheld is WithheldNone.
type Point struct {
	DeviceClass        DeviceClass        `json:"device_class"`
	Name               string             `json:"name"`
	NameStatus         NameStatus         `json:"name_status"`
	Primary            uint8              `json:"primary"`
	Secondary          uint8              `json:"secondary"`
	ValueClass         ValueClass         `json:"value_class"`
	Codec              Codec              `json:"codec"`
	CodecDeclaration   string             `json:"codec_declaration"`
	ScaleDivisor       uint16             `json:"scale_divisor"`
	NativeUnit         string             `json:"native_unit"`
	Enumeration        []EnumerationLabel `json:"enumeration"`
	ControllerFamilies []string           `json:"controller_families"`
	Locations          []Location         `json:"locations"`
	Exportable         bool               `json:"exportable"`
	Withheld           WithheldReason     `json:"withheld"`
}

func (point Point) clone() Point {
	point.ControllerFamilies = append([]string(nil), point.ControllerFamilies...)
	point.Enumeration = append([]EnumerationLabel(nil), point.Enumeration...)
	locations := make([]Location, len(point.Locations))
	for index, location := range point.Locations {
		location.ControllerFamilies = append([]string(nil), location.ControllerFamilies...)
		location.Enumeration = append([]EnumerationLabel(nil), location.Enumeration...)
		locations[index] = location
	}
	point.Locations = locations
	return point
}

// ForLocation returns a copy of the point narrowed to the location at index:
// its only location is that location, and its codec, declaration, scale,
// native unit, enumeration and controller families are that location's. Name,
// name status, value class and the export decision stay those of the point.
// It reports false when index is outside the locations.
func (point Point) ForLocation(index int) (Point, bool) {
	if index < 0 || index >= len(point.Locations) {
		return Point{}, false
	}
	narrowed := point.clone()
	location := narrowed.Locations[index]
	narrowed.Locations = []Location{location}
	narrowed.Codec, narrowed.CodecDeclaration = location.Codec, location.CodecDeclaration
	narrowed.ScaleDivisor, narrowed.NativeUnit = location.ScaleDivisor, location.NativeUnit
	narrowed.Enumeration = append([]EnumerationLabel(nil), location.Enumeration...)
	narrowed.ControllerFamilies = append([]string(nil), location.ControllerFamilies...)
	return narrowed, true
}

// Points returns a copy of every catalog row, ordered by device class and
// name.
func Points() []Point {
	points := make([]Point, len(generatedPoints))
	for index, point := range generatedPoints {
		points[index] = point.clone()
	}
	return points
}

// Lookup returns a copy of the row for a device class and canonical name.
func Lookup(deviceClass DeviceClass, name string) (Point, bool) {
	for _, point := range generatedPoints {
		if point.DeviceClass == deviceClass && point.Name == name {
			return point.clone(), true
		}
	}
	return Point{}, false
}
