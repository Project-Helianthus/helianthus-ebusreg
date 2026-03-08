package b555

import (
	"strings"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/router"
	"github.com/Project-Helianthus/helianthus-ebusreg/schema"
)

const (
	methodReadTimerConfig = "read_timer_config"
	methodReadTimerSlots  = "read_timer_slots"
	methodReadTimer       = "read_timer"
	methodWriteTimer      = "write_timer"
)

// Provider creates B555 timer/schedule planes for Vaillant controllers.
type Provider struct{}

func NewProvider() Provider {
	return Provider{}
}

func (Provider) Name() string {
	return "vaillant_b555"
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
		newTimerPlane(info),
	}
}

type plane struct {
	name    string
	address byte
	methods []registry.Method
}

func newTimerPlane(info registry.DeviceInfo) *plane {
	return &plane{
		name:    "timer",
		address: info.Address,
		methods: []registry.Method{
			method{name: methodReadTimerConfig, readOnly: true, template: configReadTemplate{primary: 0xB5, secondary: 0x55}, response: schema.SchemaSelector{}},
			method{name: methodReadTimerSlots, readOnly: true, template: slotsReadTemplate{primary: 0xB5, secondary: 0x55}, response: schema.SchemaSelector{}},
			method{name: methodReadTimer, readOnly: true, template: timerReadTemplate{primary: 0xB5, secondary: 0x55}, response: schema.SchemaSelector{}},
			method{name: methodWriteTimer, readOnly: false, template: timerWriteTemplate{primary: 0xB5, secondary: 0x55}, response: schema.SchemaSelector{}},
		},
	}
}

func (p *plane) Name() string {
	return p.name
}

func (p *plane) Methods() []registry.Method {
	return p.methods
}

func (p *plane) Subscriptions() []router.Subscription {
	return nil
}

type method struct {
	name     string
	readOnly bool
	template registry.FrameTemplate
	response schema.SchemaSelector
}

func (m method) Name() string {
	return m.name
}

func (m method) ReadOnly() bool {
	return m.readOnly
}

func (m method) Template() registry.FrameTemplate {
	return m.template
}

func (m method) ResponseSchema() schema.SchemaSelector {
	return m.response
}
