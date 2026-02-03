package system

import (
	"strings"

	"github.com/d3vi1/helianthus-ebusreg/registry"
	"github.com/d3vi1/helianthus-ebusreg/router"
	"github.com/d3vi1/helianthus-ebusreg/schema"
)

type Provider struct{}

func NewProvider() Provider {
	return Provider{}
}

func (Provider) Name() string {
	return "vaillant_system"
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
		newSystemPlane(info),
	}
}

type plane struct {
	name          string
	address       byte
	hwVersion     string
	methods       []registry.Method
	subscriptions []router.Subscription
}

func newSystemPlane(info registry.DeviceInfo) *plane {
	return &plane{
		name:      "system",
		address:   info.Address,
		hwVersion: info.HardwareVersion,
		methods: []registry.Method{
			method{name: methodGetOperationalData, readOnly: true, template: operationalTemplate{primary: 0xB5, secondary: 0x04}, response: schema.SchemaSelector{}},
			method{name: methodSetOperationalData, readOnly: false, template: operationalWriteTemplate{primary: 0xB5, secondary: 0x05}, response: schema.SchemaSelector{}},
		},
		subscriptions: []router.Subscription{
			{Primary: 0xB5, Secondary: 0x16},
			{Primary: 0xB5, Secondary: 0x16},
			{Primary: 0xFE, Secondary: 0x01},
		},
	}
}

func (plane *plane) Name() string {
	return plane.name
}

func (plane *plane) Methods() []registry.Method {
	return plane.methods
}

func (plane *plane) Subscriptions() []router.Subscription {
	return plane.subscriptions
}
