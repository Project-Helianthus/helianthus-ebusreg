// Package b555 provides the timer plane for the Vaillant B555 weekly
// timer and schedule service: reading a timer's configuration, its
// per-weekday slot counts and single slots, and writing a single slot.
//
// Matching. The provider matches on the manufacturer name only. It does
// not decide whether a device implements B555: the service is
// limited to one controller family, and another family keeps its timers
// in B524, so callers check the
// device family before choosing the timer transport.
//
// Requests and replies. The plane implements the router's plane contract
// and has no broadcast subscriptions. A request goes to the device's own
// address and needs the initiator address as the source parameter; the
// target parameter is consulted only when the device has no address. A
// reply decodes to a value map that always carries the raw reply payload
// and the request coordinates next to the decoded fields. A reply shorter
// than its layout is not an error: the fields it lacks are marked invalid,
// and an empty write reply reports the write as not accepted. Missing
// parameters, an unknown method or a reply with other command bytes fail
// with the transport library's invalid-payload error.
//
// The package also exports helpers that bound reply extents and narrow
// numeric values into graded views. They only classify the values they are
// given and send nothing; each result keeps the raw input and carries an
// explicit grade (Proven, Hypothesis or Unknown) instead of a guessed
// value.
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
