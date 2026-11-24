// Package brink provides the registry plane of the Brink 40xx family. It
// exposes every documented read and write as a method a caller can invoke by
// hand. BuildRequest uses the catalogue encoders. A request the page does not
// layout returns an error and is not sent.
package brink

import (
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
)

// Plane and method names.
const (
	planeName = "brink_40xx"

	methodSensorRead         = "sensor_read"
	methodParameterRead      = "parameter_read"
	methodErrorHistoryRead   = "error_history_read"
	methodFactoryReset       = "factory_reset"
	methodResetNotifications = "reset_notifications"
	methodFanModeTemporary   = "fan_mode_temporary_set"
	methodFanModePersistent  = "fan_mode_persistent_set"
	methodParameterWrite     = "parameter_write"
)

// Provider is the stateless registry provider of the Brink 40xx family.
type Provider struct{}

// NewProvider returns a Provider.
func NewProvider() Provider {
	return Provider{}
}

// Name returns the provider name "brink_40xx".
func (Provider) Name() string {
	return planeName
}

// Match reports false for every identity. The page identifies the family by
// primary byte 0x40 and states no manufacturer byte, and its default address
// 0x3C can be changed, so neither value is a family key on DeviceInfo.
func (Provider) Match(info registry.DeviceInfo) bool {
	return false
}

// CreatePlanes returns the family plane for a caller-supplied identity. The
// method set does not depend on info. The plane binds info.Address as the
// request target; a zero address takes the target parameter at invocation.
func (Provider) CreatePlanes(info registry.DeviceInfo) []registry.Plane {
	return []registry.Plane{newPlane(info.Address)}
}

type plane struct {
	address byte
	methods []registry.Method
}

func newPlane(address byte) *plane {
	return &plane{address: address, methods: []registry.Method{
		method{name: methodSensorRead, readOnly: true, sb: 0x22},
		method{name: methodParameterRead, readOnly: true, sb: 0x50},
		method{name: methodErrorHistoryRead, readOnly: true, sb: 0x90},
		method{name: methodFactoryReset, readOnly: false, sb: 0xFF},
		method{name: methodResetNotifications, readOnly: false, sb: 0x91},
		method{name: methodFanModeTemporary, readOnly: false, sb: 0xA1},
		method{name: methodFanModePersistent, readOnly: false, sb: 0xCB},
		method{name: methodParameterWrite, readOnly: false, sb: 0x80},
	}}
}

func (plane *plane) Name() string {
	return planeName
}

func (plane *plane) Methods() []registry.Method {
	return plane.methods
}

type method struct {
	name     string
	readOnly bool
	sb       byte
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

// Primary returns the family process byte, 0x40.
func (m method) Primary() byte {
	return 0x40
}

// Secondary returns the method's sub-protocol byte.
func (m method) Secondary() byte {
	return m.sb
}

// ResponseSchema is empty: replies are decoded by the catalogue package
// (catalog/brink), which keeps every reply the page does not define raw.
func (m method) ResponseSchema() schema.SchemaSelector {
	return schema.SchemaSelector{}
}
