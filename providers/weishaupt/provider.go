// Package weishaupt provides the registry plane of the Weishaupt WRSol and
// Kromschroeder/Weishaupt memory messages. Every documented active write or command
// is a method a caller can invoke by hand. BuildRequest uses the catalogue
// encoders. A request the page does not layout returns an error and is not sent.
package weishaupt

import (
	"strings"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
)

// Plane and method names.
const (
	planeName = "weishaupt"

	methodIdentificationRead = "wrsol_identification_read"
	methodRAMRead            = "wrsol_ram_read"
	methodDateTimeSet        = "wrsol_datetime_set"
	methodMemoryRead         = "kromschroeder_memory_read"
	methodMemoryWrite        = "kromschroeder_memory_write"
)

// matchedManufacturers are the manufacturer bytes the pages state, as the
// registry formats them: 0x10 (WRSol identification), 0x50 (Kromschroeder) and
// 0xC5 (Weishaupt).
var matchedManufacturers = map[string]bool{"0X10": true, "0X50": true, "0XC5": true}

// Provider is the stateless registry provider of the Weishaupt plane.
type Provider struct{}

// NewProvider returns a Provider.
func NewProvider() Provider {
	return Provider{}
}

// Name returns the provider name "weishaupt".
func (Provider) Name() string {
	return planeName
}

// Match reports whether the manufacturer byte is one the pages state.
func (Provider) Match(info registry.DeviceInfo) bool {
	return matchedManufacturers[strings.ToUpper(strings.TrimSpace(info.Manufacturer))]
}

// CreatePlanes returns the Weishaupt plane bound to info.Address, with every
// documented method. A zero address takes the target parameter when the method
// is invoked. wrsol_datetime_set always addresses FEh.
func (Provider) CreatePlanes(info registry.DeviceInfo) []registry.Plane {
	return []registry.Plane{newPlane(info.Address)}
}

type plane struct {
	address byte
	methods []registry.Method
}

func newPlane(address byte) *plane {
	return &plane{address: address, methods: []registry.Method{
		method{name: methodIdentificationRead, primary: 0x07, secondary: 0x04, readOnly: true},
		method{name: methodRAMRead, primary: 0x09, secondary: 0x00, readOnly: true},
		method{name: methodDateTimeSet, primary: 0x07, secondary: 0x00},
		method{name: methodMemoryRead, primary: 0x50, secondary: 0x00, readOnly: true},
		method{name: methodMemoryWrite, primary: 0x50, secondary: 0x01},
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
	return m
}

// Primary returns the method's primary command byte.
func (m method) Primary() byte {
	return m.primary
}

// Secondary returns the method's secondary command byte.
func (m method) Secondary() byte {
	return m.secondary
}

// ResponseSchema is empty: replies are decoded by the catalogue package
// (catalog/weishaupt), which keeps undocumented replies raw.
func (m method) ResponseSchema() schema.SchemaSelector {
	return schema.SchemaSelector{}
}
