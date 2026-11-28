package owner

import (
	"context"
	"testing"

	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/heat"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

// Statement fields, and the production mutation each assertion catches:
//   - controller serial: read from the registry entry at the controller address. Using the address as the serial keeps "serial-alpha" out of the statement.
//   - zone presence: an absent zone is dropped. Publishing it keeps instance 4.
//   - zone assignment and valve: copied from the zone snapshot. Dropping either leaves the assignment nil or the valve unreported.
//   - circuit activity: an inactive circuit is dropped. Publishing it keeps instance 3.
//   - cylinder instance: copied from the reported store. Dropping it leaves the cylinder list empty.
//   - generator serial: read from the boiler address. Using the controller serial instead keeps "generator-serial" out.
//   - heat-pump serial: a heat-pump device id with a serial. Using the controller serial instead keeps "pump-serial" out.
//   - nil publisher: publishThermalGraph returns without a statement.

type captureThermalPublisher struct {
	statements []heat.Statement
}

func (c *captureThermalPublisher) Publish(context.Context, semdriver.Session, heat.Statement) error {
	return nil
}

func (c *captureThermalPublisher) record(statement heat.Statement) {
	c.statements = append(c.statements, statement)
}

func TestThermalRegistration(t *testing.T) {
	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x10, Manufacturer: "Family", DeviceID: "REG00", SerialNumber: "serial-alpha"})
	reg.Register(registry.DeviceInfo{Address: 0x11, Manufacturer: "Family", DeviceID: "EHP00", SerialNumber: "pump-serial"})
	reg.Register(registry.DeviceInfo{Address: 0x12, Manufacturer: "Family", DeviceID: "REG00", SerialNumber: "generator-serial"})
	assignment := uint16(0)
	valve := uint16(1)
	poller := &vaillantSemanticPoller{
		reg:                 reg,
		controller:          0x10,
		boilerAddress:       0x12,
		boiler:              &vaillantBoilerSnapshot{},
		dhw:                 &vaillantDhwSnapshot{},
		zones:               map[byte]*vaillantZoneSnapshot{2: {Instance: 2, Present: true, ConfigurationAssociatedCircuitRaw: &assignment, StateValveStatusRaw: &valve}, 4: {Instance: 4, Present: false}},
		circuits:            map[byte]*vaillantCircuitSnapshot{0: {Instance: 0, Active: true}, 3: {Instance: 3, Active: false}},
		solarCylinders:      map[byte]*vaillantCylinderSnapshot{7: {Instance: 7}},
		installationSession: &energyCatalogueSession{},
		installationCtx:     context.Background(),
	}
	captured := &captureThermalPublisher{}
	poller.RegisterThermalPublisher(&recordingThermalPublisher{capture: captured})
	poller.publishThermalGraph()
	if len(captured.statements) != 1 {
		t.Fatalf("statements = %d", len(captured.statements))
	}
	got := captured.statements[0]
	if !got.ControllerPresent || got.ControllerSerial != "serial-alpha" || !got.GeneratorPresent || got.GeneratorSerial != "generator-serial" || !got.HotWater {
		t.Fatalf("statement identity = %+v", got)
	}
	if len(got.Zones) != 1 || got.Zones[0].Instance != 2 || !got.Zones[0].ValveReported || got.Zones[0].Assignment == nil || *got.Zones[0].Assignment != 0 {
		t.Fatalf("zones = %+v", got.Zones)
	}
	if len(got.Circuits) != 1 || got.Circuits[0].Instance != 0 {
		t.Fatalf("circuits = %+v", got.Circuits)
	}
	if len(got.Cylinders) != 1 || got.Cylinders[0].Instance != 7 {
		t.Fatalf("cylinders = %+v", got.Cylinders)
	}
	if len(got.HeatPumps) != 1 || got.HeatPumps[0].Serial != "pump-serial" {
		t.Fatalf("heat pumps = %+v", got.HeatPumps)
	}
	poller.thermalGraph = nil
	poller.publishThermalGraph()
	if len(captured.statements) != 1 {
		t.Fatal("nil publisher published another statement")
	}
	(&vaillantSemanticPoller{}).publishThermalGraph()
}

type recordingThermalPublisher struct {
	capture *captureThermalPublisher
}

func (r *recordingThermalPublisher) Publish(_ context.Context, _ semdriver.Session, statement heat.Statement) error {
	r.capture.record(statement)
	return nil
}
