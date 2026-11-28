package owner

import (
	"sync"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

// vaillantSemanticPublication is the publication owner of the Vaillant
// poller: it holds the retained semantic state, the publication events sink
// and the semantic cache persister, and no bus access. The acquisition owner
// reaches it only through vaillantSemanticPublisher, whose methods speak
// native values. Zones, hot water, circuits, radio devices, the boiler
// status and the system reach it as native views that the poller builds
// under its snapshot lock; the owner maps them to the output records with
// the mappings of the output package (publication_views.go,
// publication_native.go). The function-module verdict, solar and cylinders
// are mapped in publication_fm5.go, the schedules in publication_schedules.go,
// and the cache hydration and persistence in publication_cache.go. It is also
// the energy ledger of a register acquisition: it derives the retained
// state's energy keys from native channel, usage and window values.
//
// boilerPublishMu orders boiler-status publication. publishBoilerStatus holds
// it across the native snapshot copy and the provider update. The thermal
// control callback does not clear the native flow temperature. The lock
// guards no other state.
type vaillantSemanticPublication struct {
	state  *retainedSemanticState
	events publicationEvents
	cache  output.SemanticCachePersister

	boilerPublishMu sync.Mutex
}

func newVaillantSemanticPublication(state *retainedSemanticState, events publicationEvents, cache output.SemanticCachePersister) *vaillantSemanticPublication {
	return &vaillantSemanticPublication{state: state, events: events, cache: cache}
}

// newVaillantSemanticPoller builds the acquisition owner with a publication
// owner over state, events and cache.
func newVaillantSemanticPoller(cfg Config, substrate Substrate, state *retainedSemanticState, events publicationEvents, cache output.SemanticCachePersister) *vaillantSemanticPoller {
	return newVaillantSemanticAcquisition(cfg, substrate, newVaillantSemanticPublication(state, events, cache))
}

var _ vaillantSemanticPublisher = (*vaillantSemanticPublication)(nil)

// publicationOwner returns the poller's publication owner, or nil when the
// poller has none.
func publicationOwner(p *vaillantSemanticPoller) *vaillantSemanticPublication {
	if p == nil {
		return nil
	}
	owner, _ := p.publication.(*vaillantSemanticPublication)
	return owner
}

// live reports whether the owner holds a provider to publish to.
func (o *vaillantSemanticPublication) live() bool {
	return o != nil && o.state != nil
}

// boilerPublishLock returns the owner's boiler publication lock.
func (o *vaillantSemanticPublication) boilerPublishLock() sync.Locker {
	return &o.boilerPublishMu
}

// setNoSupportedHeatSourcePath sets or clears the provider's startup reason
// for a bus without a supported heat source path and returns the reason and
// whether the provider changed.
func (o *vaillantSemanticPublication) setNoSupportedHeatSourcePath(active bool) (string, bool) {
	return semanticStartupReasonNoSupportedHeatSourcePath, o.state.SetNoSupportedHeatSourcePath(active)
}

// publishedPlanes reports which planes the provider publishes.
func (o *vaillantSemanticPublication) publishedPlanes() vaillantPublishedPlanes {
	cylinders := o.state.Cylinders()
	return vaillantPublishedPlanes{
		zones:            o.state.Zones() != nil,
		dhw:              o.state.DHW() != nil,
		circuits:         o.state.Circuits() != nil,
		system:           o.state.System() != nil,
		radioDeviceCount: len(o.state.RadioDevices()),
		solar:            o.state.Solar() != nil,
		cylinders:        cylinders != nil,
		cylinderCount:    len(cylinders),
		boilerStatus:     o.state.BoilerStatus() != nil,
	}
}

// publishedZoneCount returns the number of zones the provider publishes.
func (o *vaillantSemanticPublication) publishedZoneCount() int {
	return len(o.state.Zones())
}

// energyReferenceTime returns the provider's energy reference instant.
func (o *vaillantSemanticPublication) energyReferenceTime() time.Time {
	return o.state.EnergyReferenceTime()
}

// applyLifetimeEnergy applies an undated lifetime register total under the
// provider's lifetime key.
func (o *vaillantSemanticPublication) applyLifetimeEnergy(channel, usage string, kwh float64) bool {
	return o.state.ApplyEnergyFromRegister(o.state.EnergyLifetimeKey(channel, usage), kwh)
}

// applyWindowEnergy applies a relative window register value under the
// provider's key dated against ref.
func (o *vaillantSemanticPublication) applyWindowEnergy(channel, usage, period, label string, ref time.Time, kwh float64) bool {
	return o.state.ApplyEnergyFromRegister(o.state.EnergyWindowKeyAt(channel, usage, period, label, ref), kwh)
}

// publishEnergyTotals pushes the provider's energy totals to the hub.
func (o *vaillantSemanticPublication) publishEnergyTotals() {
	if o == nil || o.state == nil || o.events == nil {
		return
	}
	totals := o.state.EnergyTotals()
	if totals == nil {
		return
	}
	o.events.PublishEnergyUpdate(totals)
}

// publishStartupPlaneDefaults publishes an empty circuits, solar and
// cylinders plane for each of them that the provider does not hold yet, so
// that startup readers see published planes before the first acquisition.
func (o *vaillantSemanticPublication) publishStartupPlaneDefaults() {
	if o.state.Circuits() == nil {
		o.state.SetCircuits([]output.CircuitStatus{})
	}
	if o.state.Solar() == nil {
		o.state.SetSolar(&output.SolarStatus{})
	}
	if o.state.Cylinders() == nil {
		o.state.SetCylinders([]output.CylinderStatus{})
	}
}

// publishStartupSchedulePlaceholder publishes an empty schedule program list
// when the provider holds no schedules yet.
func (o *vaillantSemanticPublication) publishStartupSchedulePlaceholder() {
	if o.state.Schedules() != nil {
		return
	}
	o.state.SetSchedules(&output.ScheduleStatus{Programs: []output.ScheduleProgram{}})
}

// publishAdapterInfo maps the native adapter hardware info to its GraphQL
// form and hands it to the provider.
func (o *vaillantSemanticPublication) publishAdapterInfo(info vaillantAdapterInfo) {
	if o == nil || o.state == nil {
		return
	}
	o.state.SetAdapterHardwareInfo(output.AdapterHardwareInfoFromNative(output.NativeAdapterInfo{
		Identity:       info.identity,
		HardwareID:     info.hardwareID,
		HardwareConfig: info.hardwareConfig,
		TemperatureC:   info.temperatureC,
		SupplyMV:       info.supplyMV,
		BusMaxDV:       info.busMaxDV,
		BusMinDV:       info.busMinDV,
		ResetCause:     info.resetCause,
		ResetCode:      info.resetCode,
		RestartCount:   info.restartCount,
		WifiRSSI:       info.wifiRSSI,
		LastIdentity:   info.lastIdentity,
		LastTelemetry:  info.lastTelemetry,
	}))
}
