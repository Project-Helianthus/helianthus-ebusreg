package heating

import (
	"strings"

	"github.com/d3vi1/helianthus-ebusreg/registry"
)

const (
	methodGetStatus     = "get_status"
	methodSetTargetTemp = "set_target_temp"
	methodGetParameters = "get_parameters"
)

type Provider struct{}

func NewProvider() Provider {
	return Provider{}
}

func (Provider) Name() string {
	return "vaillant_heating"
}

func (Provider) Match(info registry.DeviceInfo) bool {
	manufacturer := strings.TrimSpace(info.Manufacturer)
	if manufacturer == "" {
		return false
	}
	normalized := strings.ToLower(manufacturer)
	return strings.Contains(normalized, "vaillant") ||
		strings.Contains(normalized, "saunier") ||
		strings.Contains(normalized, "awb")
}

func (Provider) CreatePlanes(info registry.DeviceInfo) []registry.Plane {
	return []registry.Plane{
		newHeatingPlane(),
	}
}

type plane struct {
	name    string
	methods []registry.Method
}

func newHeatingPlane() *plane {
	return &plane{
		name: "heating",
		methods: []registry.Method{
			method{name: methodGetStatus, readOnly: true, template: template{primary: 0xB5, secondary: 0x04}},
			method{name: methodSetTargetTemp, readOnly: false, template: template{primary: 0xB5, secondary: 0x05}},
			method{name: methodGetParameters, readOnly: true, template: template{primary: 0xB5, secondary: 0x04}},
		},
	}
}

func (plane *plane) Name() string {
	return plane.name
}

func (plane *plane) Methods() []registry.Method {
	return plane.methods
}

type method struct {
	name     string
	readOnly bool
	template registry.FrameTemplate
}

func (method method) Name() string {
	return method.name
}

func (method method) ReadOnly() bool {
	return method.readOnly
}

func (method method) Template() registry.FrameTemplate {
	return method.template
}

type template struct {
	primary   byte
	secondary byte
}

func (template template) Primary() byte {
	return template.primary
}

func (template template) Secondary() byte {
	return template.secondary
}
