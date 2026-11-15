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

// Codec is a machine-readable value encoding that Decode implements.
type Codec string

// Codecs: little-endian float32 and uint16 of the controller registers, the
// signed 16-bit DATA2c (1/16) and DATA2b (1/256) and the unsigned 16-bit UIN
// word of the register service. CodecUnknown marks a missing or literal-only
// declaration.
const (
	CodecUnknown   Codec = ""
	CodecFloat32LE Codec = "float32_le"
	CodecUint16LE  Codec = "uint16_le"
	CodecData2c    Codec = "DATA2c"
	CodecData2b    Codec = "DATA2b"
	CodecUIN       Codec = "UIN"
)

// InstanceRule says how a location is instanced.
type InstanceRule string

// Instance rules: single has instance 0 only or no instance byte, indexed
// spans InstanceMin to InstanceMax, unknown has no discovery profile.
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
// declaration, scale and unit are those declared for this location;
// ControllerFamilies lists the controller families whose interpreter declares
// the codec.
type Location struct {
	Opcode             uint8        `json:"opcode"`
	Group              uint8        `json:"group"`
	Register           uint16       `json:"register"`
	Selector           string       `json:"selector"`
	InstanceRule       InstanceRule `json:"instance_rule"`
	InstanceMin        uint8        `json:"instance_min"`
	InstanceMax        uint8        `json:"instance_max"`
	ControllerFamilies []string     `json:"controller_families"`
	Codec              Codec        `json:"codec"`
	CodecDeclaration   string       `json:"codec_declaration"`
	ScaleDivisor       uint16       `json:"scale_divisor"`
	NativeUnit         string       `json:"native_unit"`
}

// Point is one catalog row, keyed by device class and canonical name. Codec,
// scale and native unit are set only when every location agrees on them; the
// native value is the decoded integer or float divided by ScaleDivisor.
// ControllerFamilies is the union of the location families. Exportable is
// true only when Withheld is WithheldNone.
type Point struct {
	DeviceClass        DeviceClass    `json:"device_class"`
	Name               string         `json:"name"`
	NameStatus         NameStatus     `json:"name_status"`
	Primary            uint8          `json:"primary"`
	Secondary          uint8          `json:"secondary"`
	ValueClass         ValueClass     `json:"value_class"`
	Codec              Codec          `json:"codec"`
	CodecDeclaration   string         `json:"codec_declaration"`
	ScaleDivisor       uint16         `json:"scale_divisor"`
	NativeUnit         string         `json:"native_unit"`
	ControllerFamilies []string       `json:"controller_families"`
	Locations          []Location     `json:"locations"`
	Exportable         bool           `json:"exportable"`
	Withheld           WithheldReason `json:"withheld"`
}

func (point Point) clone() Point {
	point.ControllerFamilies = append([]string(nil), point.ControllerFamilies...)
	locations := make([]Location, len(point.Locations))
	for index, location := range point.Locations {
		location.ControllerFamilies = append([]string(nil), location.ControllerFamilies...)
		locations[index] = location
	}
	point.Locations = locations
	return point
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
