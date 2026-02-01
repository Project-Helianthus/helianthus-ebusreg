package system

import (
	"strings"

	"github.com/d3vi1/helianthus-ebusreg/registry"
	"github.com/d3vi1/helianthus-ebusreg/router"
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
		newSystemPlane(),
	}
}

type plane struct {
	name          string
	subscriptions []router.Subscription
}

func newSystemPlane() *plane {
	return &plane{
		name: "system",
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
	return nil
}

func (plane *plane) Subscriptions() []router.Subscription {
	return plane.subscriptions
}
