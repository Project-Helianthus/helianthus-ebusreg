package solar

import (
	"strings"

	"github.com/d3vi1/helianthus-ebusreg/registry"
	"github.com/d3vi1/helianthus-ebusreg/schema"
)

const (
	methodGetStatus     = "get_status"
	methodGetSolarYield = "get_solar_yield"
	methodGetParameters = "get_parameters"
)

type Provider struct{}

func NewProvider() Provider {
	return Provider{}
}

func (Provider) Name() string {
	return "vaillant_solar"
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
		newSolarPlane(),
	}
}

type plane struct {
	name    string
	methods []registry.Method
}

func newSolarPlane() *plane {
	return &plane{
		name: "solar",
		methods: []registry.Method{
			method{name: methodGetStatus, readOnly: true, template: template{primary: 0xB5, secondary: 0x04}, response: statusSchemaSelector()},
			method{name: methodGetSolarYield, readOnly: true, template: template{primary: 0xB5, secondary: 0x04}, response: solarYieldSchemaSelector()},
			method{name: methodGetParameters, readOnly: true, template: template{primary: 0xB5, secondary: 0x04}, response: parametersSchemaSelector()},
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
	response schema.SchemaSelector
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

func (method method) ResponseSchema() schema.SchemaSelector {
	return method.response
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
