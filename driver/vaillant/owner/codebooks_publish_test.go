package owner

import (
	"bytes"
	"context"
	"encoding/binary"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func TestVaillantPublishedCataloguesFillOrdinaryViews(t *testing.T) {
	bus := &publishedCatalogueBus{
		b524: &catalogueReplyBus{payloads: publishedCataloguePayloads()},
		boiler: &boilerReplyBus{values: map[uint16][]byte{
			boilerB509FlameActive:    {0x00},
			boilerB509GasValveActive: {0x00},
			boilerB509StateNumber:    {0x07},
			boilerB509PumpHours:      {0x05, 0x00},
		}},
	}
	poller, provider := newCataloguePoller(bus.send)
	poller.boilerAddress = 0x08
	poller.configInterval = time.Minute
	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: poller.controller, Manufacturer: "Vaillant", DeviceID: "controller-700"})
	poller.reg = reg
	poller.zones = map[byte]*vaillantZoneSnapshot{0: {}}
	poller.dhw = &vaillantDhwSnapshot{}
	poller.circuits = map[byte]*vaillantCircuitSnapshot{0: {Instance: 1, Active: true}}
	poller.system = &vaillantSystemSnapshot{Controller: poller.controller}
	poller.boiler = &vaillantBoilerSnapshot{}

	poller.publishZones(semanticSnapshotSourceLive)
	poller.publishDHW(semanticSnapshotSourceLive)
	poller.publishCircuits(semanticSnapshotSourceLive)
	poller.publishSystem(semanticSnapshotSourceLive)
	poller.publishBoilerStatus(semanticSnapshotSourceLive)
	if zones := provider.Zones(); len(zones) != 1 || len(zones[0].State.Enumerations()) != 0 {
		t.Fatal("ordinary publication filled zone enumerations before the catalogue task")
	}
	if totals := provider.EnergyTotals(); totals != nil && len(totals.CatalogueDirections()) != 0 {
		t.Fatal("ordinary publication filled energy directions before the catalogue task")
	}

	var task semanticPollTask
	found := false
	for _, item := range poller.pollSchedule().recurring {
		if item.key == semanticTaskRefreshCatalogues {
			task = item
			found = true
			break
		}
	}
	if !found || task.refresh == nil || task.interval != poller.configInterval {
		t.Fatalf("catalogue task found=%v interval=%v", found, task.interval)
	}
	task.refresh(context.Background())
	assertPublishedCatalogueViews(t, provider)

	poller.publishZones(semanticSnapshotSourceLive)
	poller.publishDHW(semanticSnapshotSourceLive)
	poller.publishCircuits(semanticSnapshotSourceLive)
	poller.publishSystem(semanticSnapshotSourceLive)
	poller.publishBoilerStatus(semanticSnapshotSourceLive)
	assertPublishedCatalogueViews(t, provider)
}

func assertPublishedCatalogueViews(t *testing.T, provider *retainedSemanticState) {
	t.Helper()
	zones := provider.Zones()
	if len(zones) != 1 {
		t.Fatalf("zones = %d; want the published zone", len(zones))
	}
	assertEnumeration(t, zones[0].State.Enumerations(), localZones.group, 27, "u8", []byte{1, 0}, 1, "ventilation-boost", "Hypothesis")
	assertEnumeration(t, zones[0].State.Enumerations(), localZones.group, 0x000E, "u16", []byte{4, 0}, 4, catalogueUnresolved, "Unknown")
	dhw := provider.DHW()
	if dhw == nil {
		t.Fatal("hot water view missing")
	}
	assertEnumeration(t, dhw.State.Enumerations(), localDHW.group, 3, "u16", []byte{1, 0}, 1, "automatic", "Hypothesis")
	circuits := provider.Circuits()
	if len(circuits) != 1 {
		t.Fatalf("circuits = %d; want the published circuit", len(circuits))
	}
	assertEnumeration(t, circuits[0].Config.Enumerations(), localCircuits.group, 29, "u8", []byte{2}, 2, "party", "Hypothesis")
	system := provider.System()
	if system == nil || system.Ventilation() == nil || system.Ventilation().DayFanSpeedLimit == nil || *system.Ventilation().DayFanSpeedLimit != 300 {
		t.Fatalf("system ventilation missing; want day word 300")
	}
	status := provider.BoilerStatus()
	if status == nil {
		t.Fatal("boiler status missing")
	}
	flame := boilerSelector(status.Selectors(), boilerB509FlameActive)
	if !bytes.Equal(flame.Raw, []byte{0x00}) || flame.CodecUsed != "onoff" || flame.Label != "flame_active" {
		t.Fatalf("published flame = %+v; want raw 00 and codec onoff", flame)
	}
	pump := boilerSelector(status.Selectors(), boilerB509PumpHours)
	if !bytes.Equal(pump.Raw, []byte{0x05, 0x00}) || pump.CodecUsed != "hoursum2" {
		t.Fatalf("published pump hours = %+v; want raw 05 00 and hoursum2", pump)
	}
	totals := provider.EnergyTotals()
	if totals == nil || len(totals.RegisterDirections()) != len(b524EnergyQueries) || len(totals.CatalogueDirections()) == 0 || len(totals.ResidualFields()) == 0 || len(totals.CandidateFields()) == 0 {
		t.Fatalf("published energy directions missing")
	}
}

func publishedCataloguePayloads() map[[2]uint16][]byte {
	payloads := cataloguePayloads()
	for _, query := range b524EnergyQueries {
		value, ok := energyCatalogueValue(query.addr)
		if !ok {
			continue
		}
		payloads[[2]uint16{uint16(localRegulator.group), query.addr}] = binary.LittleEndian.AppendUint32(nil, value)
	}
	return payloads
}

type publishedCatalogueBus struct {
	b524   *catalogueReplyBus
	boiler *boilerReplyBus
}

func (bus *publishedCatalogueBus) send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	if len(frame.Data) > 0 && frame.Data[0] == vaillantB509OpcodeRead {
		return bus.boiler.send(ctx, frame)
	}
	return bus.b524.send(ctx, frame)
}
