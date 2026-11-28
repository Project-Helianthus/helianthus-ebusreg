package owner

import (
	"context"
	"sort"

	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/heat"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
)

// thermalGraphPublisher is the registration point the thermal publisher
// attaches to. A nil publisher publishes nothing.
type thermalGraphPublisher interface {
	Publish(ctx context.Context, session semdriver.Session, statement heat.Statement) error
}

// RegisterThermalPublisher attaches the publisher that runs after each native
// view is published. A nil publisher clears the attachment.
func (p *vaillantSemanticPoller) RegisterThermalPublisher(publisher thermalGraphPublisher) {
	if p == nil {
		return
	}
	p.mu.Lock()
	p.thermalGraph = publisher
	p.mu.Unlock()
}

func newThermalGraphPublisher() thermalGraphPublisher {
	return &heat.Publisher{}
}

type thermalDraft struct {
	controller    byte
	boilerAddress byte
	boilerPresent bool
	hotWater      bool
	zones         []heat.Zone
	circuits      []heat.Circuit
	cylinders     []heat.Cylinder
	catalog       productids.Catalog
	catalogErr    error
}

func (p *vaillantSemanticPoller) publishThermalGraph() {
	if p == nil {
		return
	}
	p.mu.Lock()
	publisher := p.thermalGraph
	session := p.installationSession
	ctx := p.installationCtx
	draft := p.thermalDraftLocked()
	p.mu.Unlock()
	if publisher == nil || session == nil || ctx == nil {
		return
	}
	_ = publisher.Publish(ctx, session, p.thermalStatement(draft))
}

func (p *vaillantSemanticPoller) thermalDraftLocked() thermalDraft {
	draft := thermalDraft{
		controller:    p.controller,
		boilerAddress: p.boilerAddress,
		boilerPresent: p.boiler != nil || p.boilerAddress != 0,
		hotWater:      p.dhw != nil,
		catalog:       p.catalog,
		catalogErr:    p.catalogErr,
	}
	for _, zone := range p.zones {
		if zone == nil || !zone.Present {
			continue
		}
		draft.zones = append(draft.zones, heat.Zone{
			Instance:      zone.Instance,
			Assignment:    cloneUint16Ptr(zone.ConfigurationAssociatedCircuitRaw),
			ValveReported: zone.StateValveStatusRaw != nil,
		})
	}
	for _, circuit := range p.circuits {
		if circuit == nil || !circuit.Active {
			continue
		}
		draft.circuits = append(draft.circuits, heat.Circuit{Instance: circuit.Instance})
	}
	for _, cylinder := range p.solarCylinders {
		if cylinder == nil {
			continue
		}
		draft.cylinders = append(draft.cylinders, heat.Cylinder{Instance: cylinder.Instance})
	}
	return draft
}

func (p *vaillantSemanticPoller) thermalStatement(draft thermalDraft) heat.Statement {
	statement := heat.Statement{
		ControllerPresent: draft.controller != 0,
		GeneratorPresent:  draft.boilerPresent,
		HotWater:          draft.hotWater,
		Zones:             draft.zones,
		Circuits:          draft.circuits,
		Cylinders:         draft.cylinders,
	}
	if p != nil && p.reg != nil && draft.controller != 0 {
		if snap, ok := p.reg.LookupEntrySnapshot(draft.controller); ok {
			statement.ControllerSerial = snap.SerialNumber
		}
	}
	if p != nil && p.reg != nil && draft.boilerAddress != 0 {
		if snap, ok := p.reg.LookupEntrySnapshot(draft.boilerAddress); ok {
			statement.GeneratorSerial = snap.SerialNumber
		}
	}
	statement.HeatPumps = p.heatPumps(draft.catalog, draft.catalogErr)
	sort.Slice(statement.Zones, func(i, j int) bool { return statement.Zones[i].Instance < statement.Zones[j].Instance })
	sort.Slice(statement.Circuits, func(i, j int) bool { return statement.Circuits[i].Instance < statement.Circuits[j].Instance })
	sort.Slice(statement.Cylinders, func(i, j int) bool { return statement.Cylinders[i].Instance < statement.Cylinders[j].Instance })
	return statement
}

func (p *vaillantSemanticPoller) heatPumps(catalog productids.Catalog, catalogErr error) []heat.Pump {
	if p == nil || p.reg == nil {
		return nil
	}
	var pumps []heat.Pump
	p.reg.IterateSnapshots(func(snap registry.DeviceEntrySnapshot) bool {
		if snap.SerialNumber == "" || !isHeatPumpDeviceSnapshotCandidate(snap, catalog, catalogErr) {
			return true
		}
		pumps = append(pumps, heat.Pump{Serial: snap.SerialNumber})
		return true
	})
	sort.Slice(pumps, func(i, j int) bool { return pumps[i].Serial < pumps[j].Serial })
	return pumps
}
