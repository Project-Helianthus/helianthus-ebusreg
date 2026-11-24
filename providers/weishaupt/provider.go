// Package weishaupt provides the registry plane of the Weishaupt WRSol and
// Kromschroeder/Weishaupt memory messages. Every documented write or command
// is a method a caller can invoke by hand; nothing is sent unless a method is
// called.
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
	methodSecondaryRead      = "kromschroeder_secondary_read"
	methodSecondaryWrite     = "kromschroeder_secondary_write"
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

// CreatePlanes returns the Weishaupt plane with every documented method.
func (Provider) CreatePlanes(info registry.DeviceInfo) []registry.Plane {
	return []registry.Plane{newPlane()}
}

type plane struct {
	methods []registry.Method
}

func newPlane() *plane {
	return &plane{methods: []registry.Method{
		method{name: methodIdentificationRead, primary: 0x07, secondary: 0x04, readOnly: true},
		method{name: methodRAMRead, primary: 0x09, secondary: 0x00, readOnly: true},
		method{name: methodDateTimeSet, primary: 0x07, secondary: 0x00},
		method{name: methodMemoryRead, primary: 0x50, secondary: 0x00, readOnly: true},
		method{name: methodMemoryWrite, primary: 0x50, secondary: 0x01},
		// The page states PB 0x09 with SB 0x02 for the secondary-channel read and
		// SB 0x03 for its write; it states no request layout, so these methods
		// carry the bytes only and no payload is built here.
		method{name: methodSecondaryRead, primary: 0x09, secondary: 0x02, readOnly: true},
		method{name: methodSecondaryWrite, primary: 0x09, secondary: 0x03},
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
// (catalog/weishaupt), which keeps undocumented replies raw.
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
