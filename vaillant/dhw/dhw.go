package dhw

import (
	"strings"

	"github.com/d3vi1/helianthus-ebusreg/registry"
	"github.com/d3vi1/helianthus-ebusreg/schema"
)

const (
	methodGetStatus     = "get_status"
	methodSetTargetTemp = "set_target_temp"
	methodGetParameters = "get_parameters"
	methodEnergyStats   = "get_energy_stats"

	opStatus     = byte(0x0D)
	opParameters = byte(0x09)
)

type Provider struct{}

func NewProvider() Provider {
	return Provider{}
}

func (Provider) Name() string {
	return "vaillant_dhw"
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
		newDHWPlane(info),
	}
}

type plane struct {
	name    string
	methods []registry.Method
}

func newDHWPlane(info registry.DeviceInfo) *plane {
	methods := []registry.Method{
		method{name: methodGetStatus, readOnly: true, template: opTemplate{primary: 0xB5, secondary: 0x04, op: opStatus}, response: statusSchemaSelector()},
		method{name: methodSetTargetTemp, readOnly: false, template: template{primary: 0xB5, secondary: 0x05}, response: setTargetTempSchemaSelector()},
		method{name: methodGetParameters, readOnly: true, template: opTemplate{primary: 0xB5, secondary: 0x04, op: opParameters}, response: parametersSchemaSelector()},
	}
	if supportsEnergyStats(info) {
		methods = append(methods, method{
			name:     methodEnergyStats,
			readOnly: true,
			template: energyTemplate{primary: 0xB5, secondary: 0x16},
			response: energySchemaSelector(),
		})
	}

	return &plane{
		name:    "dhw",
		methods: methods,
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

type opTemplate struct {
	primary   byte
	secondary byte
	op        byte
}

func (template opTemplate) Primary() byte {
	return template.primary
}

func (template opTemplate) Secondary() byte {
	return template.secondary
}

func (template opTemplate) Build(map[string]any) ([]byte, error) {
	return []byte{template.op}, nil
}
