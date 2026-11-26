package owner

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
)

// ConfirmB524Coherent runs the multi-register B524 capability probe
// against target and returns true if every probe responds coherently.
// Used by the runtime passive-promotion pipeline to validate that a
// candidate address discovered via passive evidence actually
// implements the Vaillant extended-register protocol before
// registering it. The probe sources from the poller's admitted
// source — source-address invariant is preserved.
func (p *vaillantSemanticPoller) ConfirmB524Coherent(ctx context.Context, target byte) bool {
	if p == nil {
		return false
	}
	probeFn := p.b524ProbeFn
	if probeFn == nil {
		probeFn = p.probeB524Register
	}
	for _, probe := range b524CapabilityProbes {
		if !probeFn(ctx, target, probe.opcode, probe.group, probe.instance, probe.addr) {
			return false
		}
	}
	return true
}

// EnqueueDiscoveryRefresh signals the semantic poller to re-run
// discovery on the next available task slot. Used by the passive-
// promotion pipeline after a late-arrival candidate is registered,
// so the regulator surface populates without a gateway restart.
func (p *vaillantSemanticPoller) EnqueueDiscoveryRefresh() {
	if p == nil {
		return
	}
	p.enqueueTask(semanticTaskRefreshDiscovery, semanticTaskPriorityHigh, p.refreshDiscovery)
}

// EnqueueAddressIdentityProbe schedules a per-address identity
// probe for the given address. Runs registry.Scan against just
// that address using the gateway's current scan source — on
// success, the registry's M6 identity-merge path collapses
// canonical pairs that share identity into a single DeviceEntry.
//
// Phase post-C P5 (live validation 2026-05-08): closes the gap
// where passive-observed addresses (e.g. NETX3 0xF1, 0xF6, BASV2
// 0x10) sit in the registry forever with empty manufacturer /
// deviceID / serialNumber because no path probes them for
// identity post-insertion.
//
// Bounded + idempotent: probes each address at most once per
// gateway lifetime via p.identityProbedAddresses (sync.Map).
// Probes that fail (timeout, NACK, etc.) are NOT retried — the
// next gateway restart re-attempts. This is intentional: NETX3's
// 0x04 broadcast face does not respond to active probes by spec,
// so retry would just spam the bus.
func (p *vaillantSemanticPoller) EnqueueAddressIdentityProbe(addr byte) {
	if p == nil || p.bus == nil || p.reg == nil {
		return
	}
	if addr == 0 || addr == 0xFE || addr == 0xAA {
		return
	}
	// P5 round-8 (Codex P2 follow-up 2026-05-08): skip the gateway's
	// own admitted source + companion. The existing
	// DiscoverB524RootInRegistry candidate path explicitly drops
	// these (skipReservedSourceCompanion); this probe path must
	// honor the same protection so we don't emit identity traffic
	// to our own admitted address (or its companion).
	probeSource, probeCompanion := p.sourceAndCompanion()
	if addr == probeSource {
		return
	}
	if probeCompanion != 0 && addr == probeCompanion {
		return
	}
	// P5 round-6 (Codex P2 follow-up 2026-05-08): resolve the input
	// address to its responder/target byte before probing. Active
	// identity reads (0x07/0x04 + B5.09 ScanID) must be addressed
	// to the responder face — sending them to an initiator/source
	// byte (e.g. BASV2 0x10, NETX3 0xF1) produces no response and
	// just times out. The existing startup scan path filters
	// initiator-capable addresses out of target lists for the same
	// reason.
	//
	// Resolution order:
	// 1. If addr is non-initiator-capable, it IS already a
	//    responder/target — probe directly.
	// 2. Otherwise, look up the entry containing addr and use
	//    TargetAddressForRouting (prefers SlotRoleSlave, falls
	//    back to PrimaryDisplayAddress when no target-role face
	//    exists).
	// 3. If neither yields a valid target, skip — initiator-only
	//    devices have no probable identity face.
	probeAddr := addr
	if protocol.IsInitiatorCapableAddress(addr) {
		// P9.3 — race-free identity-probe target resolution via
		// IterateSnapshots + SnapshotTargetAddressForRouting. The
		// semantic poller runs concurrently with passive inserter /
		// startup scan registry writes; the previous Iterate path
		// dereferenced live entry.Addresses() and TargetAddressForRouting.
		var resolvedTarget byte
		p.reg.IterateSnapshots(func(snap registry.DeviceEntrySnapshot) bool {
			if snapshotContainsAddress(snap, addr) {
				resolvedTarget = snapshotTargetAddressForRouting(snap)
				return false
			}
			return true
		})
		if resolvedTarget == 0 || protocol.IsInitiatorCapableAddress(resolvedTarget) {
			// No target-role face on the entry, OR the entry's
			// only addresses are initiator-capable — skip.
			log.Printf("semantic_address_identity_probe address=0x%02X status=skipped reason=no_responder_face", addr)
			return
		}
		probeAddr = resolvedTarget
	}
	if _, alreadyProbed := p.identityProbedAddresses.LoadOrStore(probeAddr, struct{}{}); alreadyProbed {
		return
	}
	scheduler := p.tasks
	probeFn := func(ctx context.Context) {
		probeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		scanSource, _ := p.sourceAndCompanion()
		_, err := registryScanDirectedFn(probeCtx, p.bus, p.reg, scanSource, []byte{probeAddr})
		if err != nil {
			log.Printf("semantic_address_identity_probe address=0x%02X status=fail err=%v", probeAddr, err)
			return
		}
		log.Printf("semantic_address_identity_probe address=0x%02X status=ok", probeAddr)
		if probeAddr == circuitManagingDeviceFunctionModuleControllerAddress {
			recoveryCtx, recoveryCancel := context.WithTimeout(ctx, semanticIdentityRecoveryTimeout)
			defer recoveryCancel()
			p.refreshFunctionModuleConfigurationForLateIdentity(recoveryCtx)
			p.refreshFunctionModuleAfterLateIdentity(recoveryCtx)
		}
	}
	if scheduler == nil {
		// No scheduler wired (unit tests, etc.) — run inline.
		probeFn(context.Background())
		return
	}
	// Codex P2 follow-up on PR #583 (2026-05-08): submit BEFORE
	// committing the per-address probed marker. If the task queue
	// is overloaded and submit fails, roll back the sync.Map entry
	// so the address remains eligible for retry on the next passive
	// observation. Without rollback, an overloaded queue at startup
	// permanently blocks identity enrichment for the dropped
	// addresses until the next gateway restart.
	err := scheduler.submit(semanticTaskPriorityLow, func(taskCtx context.Context) {
		p.withPollLock(taskCtx, probeFn)
	})
	if err != nil {
		// Rollback so the (resolved) address is eligible for retry.
		p.identityProbedAddresses.Delete(probeAddr)
		if errors.Is(err, errSemanticTaskQueueOverloaded) {
			log.Printf("semantic_address_identity_probe address=0x%02X status=enqueue_dropped reason=queue_overloaded", probeAddr)
			return
		}
		log.Printf("semantic_address_identity_probe address=0x%02X status=enqueue_failed err=%v", probeAddr, err)
	}
}

// registerStructuralControllerIfMissing registers a minimal Vaillant
// device entry for a B524-coherent controller address that is not yet
// in the registry. This closes the gap where discoverB524Root converges
// via the structural-fallback path (probing 0x15 / 0x26 even when they
// are absent from the registry) — the address must surface in
// GraphQL devices and on the router plane after discovery.
//
// The minimal entry is safe: the address has just satisfied the B524
// multi-register coherency check, which is strong evidence of a
// Vaillant regulator. Registry de-dupes on address, so this is a
// no-op for already-known controllers. Manufacturer is the only
// identity field set; richer fields (DeviceID, SerialNumber,
// SoftwareVersion) are populated naturally by subsequent identity
// reads or by the discovery promotion pipeline once that lands.
// registerStructuralControllerIfMissing returns true iff the registry
// gained a new entry as a result of this call (false when the
// controller was already known, was zero, or the poller was nil).
// Callers use the return value to decide whether downstream caches
// (e.g. regulator capability) need recomputation.
func (p *vaillantSemanticPoller) registerStructuralControllerIfMissing(controller byte) bool {
	if p == nil || p.reg == nil || controller == 0 {
		return false
	}
	// Phase C M-C6b: containment must check the full address set,
	// not just PrimaryDisplayAddress. A controller (e.g. 0x15) that
	// is already an alias face on a registered entry (e.g. an
	// aliased 0x10↔0x15 regulator displayed as 0x10) must be
	// treated as already-known, otherwise the function violates
	// its idempotence contract and re-registers + re-refreshes
	// router planes on every discovery cycle.
	// P9.3 — race-free address-membership check via IterateSnapshots.
	// The semantic poller runs concurrently with passive inserter /
	// startup scan / identity probe goroutines that mutate the
	// registry; the previous Iterate path read entry.Addresses() on
	// a live *deviceEntry pointer, racing with concurrent Register /
	// AliasAddresses writes.
	already := false
	p.reg.IterateSnapshots(func(snap registry.DeviceEntrySnapshot) bool {
		if snapshotContainsAddress(snap, controller) {
			already = true
			return false
		}
		return true
	})
	if already {
		return false
	}
	p.reg.Register(registry.DeviceInfo{
		Address:      controller,
		Manufacturer: "Vaillant",
	})
	log.Printf("semantic_controller_registry_register address=0x%02x source=b524_structural_fallback", controller)
	if p.routerPlanesRefreshFn != nil {
		p.routerPlanesRefreshFn()
	}
	return true
}

// enrichRegulatorIdentity returns enrichment metadata for the B524 root device.
// Unknown identity is acceptable and does not block discovery.
func (p *vaillantSemanticPoller) enrichRegulatorIdentity(addr byte) *regulatorEnrichment {
	if p == nil || p.reg == nil {
		return nil
	}

	// Phase C M-C6b: discoverB524Root returns the routed controller
	// face (e.g. 0x15 for an aliased 0x10↔0x15 regulator). The
	// existing entry is displayed as 0x10, so matching only
	// PrimaryDisplayAddress would miss it and the enrichment
	// metadata/logging would be lost for the aliased path. Use the
	// full address-set membership check instead.
	//
	// P9.3 — race-free via IterateSnapshots; reads entry's DeviceID
	// from the value-typed snapshot (immune to concurrent Register
	// writes that would torn-read the deviceID string).
	var (
		matched  bool
		deviceID string
	)
	p.reg.IterateSnapshots(func(snap registry.DeviceEntrySnapshot) bool {
		if snapshotContainsAddress(snap, addr) {
			matched = true
			deviceID = snap.DeviceID
			return false
		}
		return true
	})
	if !matched {
		return nil
	}

	normalizedID := normalizeDeviceID(deviceID)
	families := []string{"BASV2", "BASS2", "CTLV2", "CTLS2", "E7C00"}
	var family string
	for _, f := range families {
		if strings.HasPrefix(normalizedID, f) {
			family = f
			break
		}
	}

	return &regulatorEnrichment{
		family:   family,
		deviceID: deviceID,
	}
}

// publishRegulatorCapability hands the controller capability to the
// publication owner; a poller without one publishes nothing.
func (p *vaillantSemanticPoller) publishRegulatorCapability(capability productids.ControllerCapability) {
	if p == nil || p.publication == nil {
		return
	}
	p.publication.publishRegulatorCapability(capability)
}
