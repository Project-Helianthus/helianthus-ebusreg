package solar

import (
	"strings"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
)

const (
	methodGetStatus     = "get_status"
	methodGetSolarYield = "get_solar_yield"
	methodGetParameters = "get_parameters"
)

// Provider is the stateless registry provider of the Vaillant solar plane.
type Provider struct{}

// NewProvider returns a Provider.
func NewProvider() Provider {
	return Provider{}
}

// Name returns the provider name "vaillant_solar".
func (Provider) Name() string {
	return "vaillant_solar"
}

// Match reports whether the manufacturer names the Vaillant group: a
// case-insensitive match of "vaillant", "saunier" or "awb". An empty
// manufacturer never matches.
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

// CreatePlanes returns the single read-only "solar" plane with get_status,
// get_solar_yield and get_parameters; info does not change the plane.
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
