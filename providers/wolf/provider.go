// Package wolf provides the registry plane of the Wolf burner and solar
// messages. Every documented write or command is a method a caller can invoke
// by hand; nothing is sent unless a method is called.
package wolf

import (
	"strings"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
)

// Plane and method names.
const (
	planeName = "wolf"

	methodBurnerStatusSend   = "burner_status_send"
	methodModulationSend     = "modulation_send"
	methodHeatRequestSend    = "heat_request_send"
	methodBoilerSetpointSend = "boiler_setpoint_send"
	methodSolarStatusSend    = "solar_status_send"
	methodSolarPumpSend      = "solar_pump_send"
	methodSolarPowerSend     = "solar_power_send"
	methodServiceRead        = "service_parameter_read"
	methodServiceWrite       = "service_parameter_write"
)

// wolfManufacturer is the manufacturer byte the page states, as the registry
// formats it.
const wolfManufacturer = "0x19"

// Provider is the stateless registry provider of the Wolf plane.
type Provider struct{}

// NewProvider returns a Provider.
func NewProvider() Provider {
	return Provider{}
}

// Name returns the provider name "wolf".
func (Provider) Name() string {
	return planeName
}

// Match reports whether the manufacturer is the Wolf byte 0x19. Matching is
// case-insensitive on the hexadecimal form.
func (Provider) Match(info registry.DeviceInfo) bool {
	return strings.EqualFold(strings.TrimSpace(info.Manufacturer), wolfManufacturer)
}

// CreatePlanes returns the Wolf plane with every documented write and command,
// and the service-parameter read.
func (Provider) CreatePlanes(info registry.DeviceInfo) []registry.Plane {
	return []registry.Plane{newPlane()}
}

type plane struct {
	methods []registry.Method
}

func newPlane() *plane {
	return &plane{methods: []registry.Method{
		method{name: methodBurnerStatusSend, primary: 0x05, secondary: 0x03},
		method{name: methodModulationSend, primary: 0x05, secondary: 0x04},
		method{name: methodHeatRequestSend, primary: 0x05, secondary: 0x07},
		method{name: methodBoilerSetpointSend, primary: 0x08, secondary: 0x00},
		method{name: methodSolarStatusSend, primary: 0x50, secondary: 0x14},
		method{name: methodSolarPumpSend, primary: 0x50, secondary: 0x17},
		method{name: methodSolarPowerSend, primary: 0x50, secondary: 0x18},
		method{name: methodServiceRead, primary: 0x50, secondary: 0x22, readOnly: true},
		method{name: methodServiceWrite, primary: 0x50, secondary: 0x23},
	}}
}

func (plane *plane) Name() string {
	return planeName
}

func (plane *plane) Methods() []registry.Method {
	return plane.methods
}

type method struct {
	name      string
	primary   byte
	secondary byte
	readOnly  bool
}

func (m method) Name() string {
	return m.name
}

func (m method) ReadOnly() bool {
	return m.readOnly
}

func (m method) Template() registry.FrameTemplate {
	return template{primary: m.primary, secondary: m.secondary}
}

// ResponseSchema is empty: replies are decoded by the catalogue package
// (catalog/wolf), which keeps undocumented frames raw.
func (m method) ResponseSchema() schema.SchemaSelector {
	return schema.SchemaSelector{}
}

type template struct {
	primary   byte
	secondary byte
}

func (t template) Primary() byte {
	return t.primary
}

func (t template) Secondary() byte {
	return t.secondary
}
