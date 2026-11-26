package owner

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
)

type regulatorAbsenceState string

const (
	// regulatorAbsenceUnknown is the state before the first regulator
	// recheck has evaluated the capability.
	regulatorAbsenceUnknown regulatorAbsenceState = "UNKNOWN"
	regulatorPresent        regulatorAbsenceState = "PRESENT"
	regulatorAbsenceGrace   regulatorAbsenceState = "ABSENCE_GRACE"
	regulatorAbsent         regulatorAbsenceState = "ABSENT"
)

type regulatorEnrichment struct {
	family   string
	deviceID string
}

type b524ProbeSpec struct {
	opcode   byte
	group    byte
	instance byte
	addr     uint16
}

// b524CapabilityProbes lists registers from different GG groups used for
// multi-register coherency verification during capability-first discovery.
var b524CapabilityProbes = []b524ProbeSpec{
	{opcode: localRegulator.opcode, group: localRegulator.group, instance: regulatorInstance, addr: 0x0001},
	{opcode: localDHW.opcode, group: localDHW.group, instance: dhwInstance, addr: dhw_operation_mode},
}

func (p *vaillantSemanticPoller) refreshRegulatorCapability(_ context.Context) {
	if p == nil || p.reg == nil {
		return
	}

	regCap := p.findRegulatorCapability()
	deviceCount := countRegistryDevices(p.reg)
	now := p.now()

	p.mu.Lock()
	prevCap := p.regulatorCapability
	prevAbsence := p.regAbsenceState
	prevDeviceCount := p.registryDeviceCount
	p.regulatorCapability = regCap
	p.registryDeviceCount = deviceCount

	// Absence grace FSM transitions.
	switch {
	case regCap == productids.ControllerPresent:
		if p.regAbsenceState != regulatorPresent {
			p.regAbsenceState = regulatorPresent
			p.regAbsenceSince = time.Time{}
		}
	case prevAbsence == regulatorPresent || prevAbsence == regulatorAbsenceUnknown:
		// Was present, or not yet evaluated, and is not present: enter the
		// grace window from this evaluation.
		p.regAbsenceState = regulatorAbsenceGrace
		p.regAbsenceSince = now
	case prevAbsence == regulatorAbsenceGrace:
		if !p.regAbsenceSince.IsZero() && now.Sub(p.regAbsenceSince) >= p.regulatorAbsenceGrace {
			p.regAbsenceState = regulatorAbsent
		}
	}
	newAbsence := p.regAbsenceState
	p.mu.Unlock()
	p.publishRegulatorCapability(regCap)
	semanticRegulatorState.Set(string(newAbsence))

	// Log capability changes.
	if regCap != prevCap {
		log.Printf("semantic_regulator_capability capability=%s", regCap.String())
	}

	// Log absence state transitions.
	if newAbsence != prevAbsence {
		semanticRegulatorTransitionsTotal.Add(fmt.Sprintf("%s->%s", prevAbsence, newAbsence), 1)
		log.Printf("semantic_regulator_transition from=%s to=%s capability=%s", prevAbsence, newAbsence, regCap.String())
		log.Printf("semantic_regulator_absence state=%s capability=%s", newAbsence, regCap.String())
	}

	// The reason depends on the regulator capability: re-evaluate it here so
	// an enrichment that classifies a registered identity as a controller
	// clears it without waiting for a discovery pass.
	p.refreshHeatSourcePathReason()

	// If device count changed, enqueue an immediate full discovery refresh.
	if deviceCount != prevDeviceCount && prevDeviceCount > 0 {
		log.Printf("semantic_regulator_recheck inventory_change prev=%d curr=%d", prevDeviceCount, deviceCount)
		p.enqueueTask(semanticTaskRefreshDiscovery, semanticTaskPriorityLow, p.refreshDiscovery)
	}
}

// fallbackDiscoveryAllowed reports whether the state and configuration tasks
// may run their fallback root discovery. On a bus where this process has never
// discovered a B524 root, once the regulator absence state is ABSENT, root
// discovery runs only from the discovery task at the configured discovery
// interval and from event-driven refreshes (inventory change, passive
// insertion and promotion). ABSENT alone does not prove that no controller
// exists (an unclassified or structurally registered controller also reaches
// it), so a process that has discovered a root once keeps the per-minute
// recovery, and a known controller without zones keeps the fallback. The
// fallback is suspended only while the periodic discovery task runs: the poll
// schedule starts that task's loop only for a positive discovery interval, so
// with the interval disabled the fallback stays the only periodic search.
func (p *vaillantSemanticPoller) fallbackDiscoveryAllowed(controller byte) bool {
	p.mu.Lock()
	periodicDiscovery := p.discoveryInterval > 0
	suspend := controller == 0 && p.regAbsenceState == regulatorAbsent && !p.controllerEverDiscovered && periodicDiscovery
	logSuspension := suspend && !p.discoveryFallbackSuspended
	p.discoveryFallbackSuspended = suspend
	interval := p.discoveryInterval
	p.mu.Unlock()
	if logSuspension {
		log.Printf("semantic_discovery_fallback_suspended regulator_state=%s discovery_interval=%s", regulatorAbsent, interval)
	}
	return !suspend
}

func (p *vaillantSemanticPoller) refreshDiscovery(ctx context.Context) {
	// Regulator capability is always recomputed, even when no B524 root is found.
	regCap := p.findRegulatorCapability()
	boilerAddress := p.findBoilerAddress()

	controller, discoveryErr := p.discoverB524Root(ctx)
	if discoveryErr != nil {
		log.Printf("semantic_b524_root: %v", discoveryErr)
		p.mu.Lock()
		prev := p.regulatorCapability
		prevController := p.controller
		prevBoilerAddress := p.boilerAddress
		p.controller = 0
		p.boilerAddress = boilerAddress
		p.regulatorCapability = regCap
		p.zones = make(map[byte]*vaillantZoneSnapshot)
		p.presence = make(map[byte]*zonePresenceRecord)
		p.circuits = make(map[byte]*vaillantCircuitSnapshot)
		p.functionModuleEvidenceGeneration++
		p.startupSemanticPrimed = false
		semanticZoneCount.Set(0)
		p.mu.Unlock()
		p.publishRegulatorCapability(regCap)
		p.refreshHeatSourcePathReason()
		if regCap != prev {
			log.Printf("semantic_regulator_capability capability=%s", regCap.String())
		}
		if prevController != 0 {
			log.Printf("semantic_controller_discovery address=0x00 source=no_coherent_responder")
		}
		if boilerAddress != prevBoilerAddress && boilerAddress != 0 && boilerAddress != 0x08 {
			log.Printf("semantic_boiler_discovery address=0x%02x source=nonstandard", boilerAddress)
		}
		if boilerAddress != prevBoilerAddress && boilerAddress != 0 {
			p.enqueueBoilerStatusPriming(ctx)
		}
		p.publishZones(semanticSnapshotSourceCache)
		p.publishDHW(semanticSnapshotSourceCache)
		p.publishCircuits(semanticSnapshotSourceCache)
		p.publishRadioDevices(semanticSnapshotSourceCache)
		p.refreshFunctionModuleSemantic(ctx)
		return
	}

	// When B524 root discovery converges via the structural-fallback
	// path, the controller address may not yet be in the registry —
	// the directed startup probe that would normally register it
	// either has not run for this address (e.g. the regulator missed
	// the startup window) or did not yield enrichment. Register a
	// minimal Vaillant entry so GraphQL devices and the router plane
	// reflect the regulator. Registry de-dupes on address, so a
	// no-op when the controller is already known.
	registered := p.registerStructuralControllerIfMissing(controller)

	if enrichment := p.enrichRegulatorIdentity(controller); enrichment != nil {
		log.Printf("semantic_controller_enrichment address=0x%02x family=%s device=%s",
			controller, enrichment.family, enrichment.deviceID)
	}

	// Recompute regulator capability after structural registration:
	// the registry just gained a new Vaillant entry, so the
	// capability lookup must run against the post-registration
	// inventory. Without this recompute, p.regulatorCapability
	// stays at the pre-registration value (typically ControllerNone
	// when the registry held only the boiler) until the next
	// regulatorRecheckInterval tick — leaving status surfaces
	// reporting "no regulator" while semantic discovery already
	// has one.
	if registered {
		regCap = p.findRegulatorCapability()
	}

	p.mu.Lock()
	prev := p.regulatorCapability
	prevController := p.controller
	prevBoilerAddress := p.boilerAddress
	p.controller = controller
	p.controllerEverDiscovered = true
	p.discoveryFallbackSuspended = false
	if controller != prevController {
		p.functionModuleEvidenceGeneration++
	}
	p.boilerAddress = boilerAddress
	p.regulatorCapability = regCap
	primeStartup := !p.startupSemanticPrimed
	if primeStartup {
		p.startupSemanticPrimed = true
	}
	p.mu.Unlock()
	p.publishRegulatorCapability(regCap)
	p.refreshHeatSourcePathReason()

	if primeStartup {
		p.refreshStartupCriticalSemanticPlanes(ctx)
	}
	p.refreshZoneDiscovery(ctx, primeStartup)
	if primeStartup {
		p.refreshStartupSemanticPlanes(ctx)
	}
	if regCap != prev {
		log.Printf("semantic_regulator_capability capability=%s", regCap.String())
	}
	if controller != prevController && controller != 0 {
		log.Printf("semantic_controller_discovery address=0x%02x", controller)
		p.enqueueControllerSemanticPriming(ctx)
	}
	if boilerAddress != prevBoilerAddress && boilerAddress != 0 && boilerAddress != 0x08 {
		log.Printf("semantic_boiler_discovery address=0x%02x source=nonstandard", boilerAddress)
	}
	if boilerAddress != prevBoilerAddress && boilerAddress != 0 {
		p.enqueueBoilerStatusPriming(ctx)
	}
}

func (p *vaillantSemanticPoller) tryRefreshFromEbusdGrab(ctx context.Context) bool {
	_, ok := p.tryRefreshFromEbusdGrabWithInstances(ctx)
	return ok
}

func (p *vaillantSemanticPoller) tryRefreshFromEbusdGrabWithInstances(ctx context.Context) (map[byte]bool, bool) {
	if p == nil {
		return nil, false
	}
	if p.transportConfig.Protocol != engine.TransportEbusdTCP {
		return nil, false
	}
	if p.refreshFromEbusdGrabFn != nil {
		return p.refreshFromEbusdGrabFn(ctx)
	}
	return p.refreshFromEbusdGrabWithInstances(ctx)
}

func normalizeDeviceID(value string) string {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	normalized = strings.ReplaceAll(normalized, "_", "")
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, " ", "")
	return normalized
}

// findRegulatorCapability iterates all devices in the registry and determines
// the regulator capability strictly from the product_ids catalog classification.
// Returns ControllerPresent if any Vaillant device has role=="Regulator",
// ControllerNone if all Vaillant devices have known non-regulator roles,
// ControllerUnknown if any device has an unknown classification or no serial.
func (p *vaillantSemanticPoller) findRegulatorCapability() productids.ControllerCapability {
	if p == nil || p.reg == nil || p.catalogErr != nil {
		return productids.ControllerUnknown
	}

	// P9.3 — race-free Manufacturer/SerialNumber reads via snapshot.
	foundPresent := false
	hasUnknown := false
	hasAnyDevice := false
	p.reg.IterateSnapshots(func(snap registry.DeviceEntrySnapshot) bool {
		if !strings.EqualFold(snap.Manufacturer, "Vaillant") {
			return true
		}
		hasAnyDevice = true
		partNumber := extractPartNumberFromSerial(snap.SerialNumber)
		cap := p.catalog.ControllerCapability(partNumber)
		switch cap {
		case productids.ControllerPresent:
			foundPresent = true
			return false // stop iteration
		case productids.ControllerUnknown:
			hasUnknown = true
		}
		return true
	})

	if foundPresent {
		return productids.ControllerPresent
	}
	if !hasAnyDevice {
		return productids.ControllerUnknown
	}
	if hasUnknown {
		return productids.ControllerUnknown
	}
	return productids.ControllerNone
}
