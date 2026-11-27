package owner

import (
	"context"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
	"strings"
	"testing"
	"time"
)

func TestVaillantOwnerCapabilityFirstDiscoveryPrefersPreferredControllerWhenCoherent(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})
	reg.Register(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"})
	reg.Register(registry.DeviceInfo{Address: 0x26, Manufacturer: "Vaillant", DeviceID: "VR_71"})

	var probed []byte
	poller := newTestPoller(reg)
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x08: true, 0x15: true, 0x26: true}, &probed)

	addr, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v", err)
	}
	if addr != 0x15 {
		t.Fatalf("discoverB524Root = 0x%02x; want 0x15 (D1-ordered preference)", addr)
	}
	// 0x15 should be probed first and discovery should stop there
	if len(probed) != 1 || probed[0] != 0x15 {
		t.Fatalf("probed = %v; want [0x15] (stop-at-first)", probed)
	}
}

func TestVaillantOwnerCapabilityFirstDiscoveryFallsBackWhenPreferredControllerNonCoherent(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})
	reg.Register(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"})
	reg.Register(registry.DeviceInfo{Address: 0x26, Manufacturer: "Vaillant", DeviceID: "VR_71"})

	var probed []byte
	// 0x15 non-coherent, 0x26 coherent
	poller := newTestPoller(reg)
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x08: false, 0x15: false, 0x26: true}, &probed)

	addr, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v", err)
	}
	if addr != 0x26 {
		t.Fatalf("discoverB524Root = 0x%02x; want 0x26 (fallback after 0x15 fails)", addr)
	}
	// Should probe 0x15 first, then 0x08, then 0x26 (ascending after 0x15)
	if len(probed) < 2 {
		t.Fatalf("probed = %v; want at least [0x15, ...0x26]", probed)
	}
	if probed[0] != 0x15 {
		t.Fatalf("probed[0] = 0x%02x; want 0x15 (probed first per D1)", probed[0])
	}
}

func TestVaillantOwnerCapabilityFirstDiscoveryRejectsShortResponse(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"})

	// Probe function that simulates short response (< 4 bytes) → non-coherent
	poller := newTestPoller(reg)
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x15: false}, nil)

	_, err := poller.discoverB524Root(context.Background())
	if err == nil {
		t.Fatal("discoverB524Root should return error when all responses are short/non-coherent")
	}
	if !strings.Contains(err.Error(), "no coherent responder") {
		t.Fatalf("error = %q; want substring 'no coherent responder'", err.Error())
	}
}

func TestVaillantOwnerCapabilityFirstDiscoveryReturnsDefinitiveWhenNoneQualify(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})
	reg.Register(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"})

	poller := newTestPoller(reg)
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x08: false, 0x15: false}, nil)

	addr, err := poller.discoverB524Root(context.Background())
	if err == nil {
		t.Fatal("discoverB524Root should return error when no candidate qualifies")
	}
	if addr != 0 {
		t.Fatalf("discoverB524Root addr = 0x%02x; want 0 on failure", addr)
	}
	if !strings.Contains(err.Error(), "no coherent responder") {
		t.Fatalf("error = %q; want substring 'no coherent responder'", err.Error())
	}
}

func TestVaillantOwnerCapabilityFirstDiscoverySucceedsWhenEnrichmentFails(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	// Register device with unknown identity (no known family prefix)
	reg.Register(registry.DeviceInfo{Address: 0x15, Manufacturer: "Unknown", DeviceID: "XYZABC"})

	poller := newTestPoller(reg)
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x15: true}, nil)

	addr, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v; discovery should succeed regardless of enrichment", err)
	}
	if addr != 0x15 {
		t.Fatalf("discoverB524Root = 0x%02x; want 0x15", addr)
	}

	// Enrichment should return empty family but not nil
	enrichment := poller.enrichRegulatorIdentity(addr)
	if enrichment == nil {
		t.Fatal("enrichRegulatorIdentity should return non-nil for registered device")
	}
	if enrichment.family != "" {
		t.Fatalf("enrichment.family = %q; want empty for unknown device", enrichment.family)
	}
}

func TestVaillantOwnerCapabilityFirstDiscoverySkipsNonCoherentDevice(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})
	reg.Register(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"})

	var probed []byte
	// 0x08 responds but non-coherent (e.g., boiler, not a B524 controller)
	// 0x15 responds coherently
	poller := newTestPoller(reg)
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x08: false, 0x15: true}, &probed)

	addr, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v", err)
	}
	if addr != 0x15 {
		t.Fatalf("discoverB524Root = 0x%02x; want 0x15", addr)
	}
	// Both should be probed (0x15 first per D1, then stop)
	if len(probed) != 1 || probed[0] != 0x15 {
		t.Fatalf("probed = %v; want [0x15] (0x15 probed first, coherent, stop)", probed)
	}
}

// TestCapabilityFirstDiscovery_StructuralFallbackWhenRegistryHasOnlyBoiler
// reproduces the post-source-selection regression where startup admission
// completes with only the boiler in the registry. discoverB524Root must
// augment its candidate list with the bounded Vaillant structural target
// set so the regulator (0x15) can be probed even before it appears in
// the registry.
func TestVaillantOwnerCapabilityFirstDiscoveryStructuralFallbackWhenRegistryHasOnlyBoiler(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})

	var probed []byte
	// Only 0x15 is coherent. Without structural fallback, 0x15 is not in
	// the registry and discoverB524Root would return "no coherent
	// responder among [0x08]".
	poller := newTestPoller(reg)
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x08: false, 0x15: true, 0x26: false}, &probed)

	addr, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v; structural fallback must add 0x15/0x26 to candidates", err)
	}
	if addr != 0x15 {
		t.Fatalf("discoverB524Root = 0x%02x; want 0x15 (structural fallback regulator)", addr)
	}
	if len(probed) == 0 || probed[0] != 0x15 {
		t.Fatalf("probed = %v; want 0x15 first per D1 ordering", probed)
	}
}

// TestCapabilityFirstDiscovery_RichRegistryStopsAt0x15CoherentRoot pins the
// stop-at-first-coherent invariant when the registry already contains the
// regulator. Structural augmentation is now unconditional (idempotent — see
// the hostile-registry rationale in the next test) but D1 ordering must
// still terminate on the first coherent candidate. With 0x15 coherent,
// 0x26 must never be probed.
func TestVaillantOwnerCapabilityFirstDiscoveryRichRegistryStopsAtPreferredControllerCoherentRoot(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})
	reg.Register(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"})

	var probed []byte
	poller := newTestPoller(reg)
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x08: false, 0x15: true, 0x26: true}, &probed)

	addr, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v", err)
	}
	if addr != 0x15 {
		t.Fatalf("discoverB524Root = 0x%02x; want 0x15", addr)
	}
	for _, p := range probed {
		if p == 0x26 {
			t.Fatalf("D1 stop-at-first violated: probed includes 0x26 after 0x15 coherent: %v", probed)
		}
	}
}

// TestCapabilityFirstDiscovery_HostileRegistryWithStaleSecondaryStillProbesRegulator
// closes the post-source-selection regression hole identified during
// adversarial review: a registry containing the boiler plus any stray
// non-Vaillant entry (e.g. 0xEC SOL00 from passive observation, or a
// phantom from a prior process) must NOT suppress the structural
// fallback. Always-augment guarantees 0x15 enters the candidate list
// regardless of registry size.
func TestVaillantOwnerCapabilityFirstDiscoveryHostileRegistryWithStaleSecondaryStillProbesRegulator(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})
	// Stale passive-only-responder entry from prior session — never
	// responds to active B524 probes.
	reg.Register(registry.DeviceInfo{Address: 0xEC, Manufacturer: "", DeviceID: ""})

	var probed []byte
	poller := newTestPoller(reg)
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x08: false, 0xEC: false, 0x15: true, 0x26: false}, &probed)

	addr, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v; structural fallback must still fire when registry contains stray passive-only entries", err)
	}
	if addr != 0x15 {
		t.Fatalf("discoverB524Root = 0x%02x; want 0x15", addr)
	}
	if len(probed) == 0 || probed[0] != 0x15 {
		t.Fatalf("probed = %v; want 0x15 first per D1 ordering", probed)
	}
}

// TestRefreshDiscovery_StructuralFallbackRegistersControllerInRegistry
// locks the live-evidence claim: when discoverB524Root converges on 0x15
// via the structural-fallback path (i.e. 0x15 was NOT in the registry
// when discovery started), refreshDiscovery must register 0x15 so the
// regulator surfaces in GraphQL devices and on the router plane.
// Without this side-effect the user-visible "GraphQL devices includes
// regulator 0x15" outcome is unattainable from the gateway-side fix.
func TestVaillantOwnerRefreshDiscoveryStructuralFallbackRegistersControllerInRegistry(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})

	poller := newTestPoller(reg)
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x08: false, 0x15: true, 0x26: false}, nil)

	addr, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v", err)
	}
	if addr != 0x15 {
		t.Fatalf("discoverB524Root = 0x%02x; want 0x15", addr)
	}

	// Pre-condition: 0x15 not in registry.
	var sawControllerBefore bool
	reg.Iterate(func(e registry.DeviceEntry) bool {
		if e != nil && e.PrimaryDisplayAddress() == 0x15 {
			sawControllerBefore = true
			return false
		}
		return true
	})
	if sawControllerBefore {
		t.Fatal("test setup invalid: 0x15 already in registry before structural-fallback registration")
	}

	// Trigger the registration helper directly. (refreshDiscovery wraps
	// this with snapshot rebuilds we don't need to exercise here.)
	poller.registerStructuralControllerIfMissing(addr)

	var sawControllerAfter bool
	var manufacturer string
	reg.Iterate(func(e registry.DeviceEntry) bool {
		if e != nil && e.PrimaryDisplayAddress() == 0x15 {
			sawControllerAfter = true
			manufacturer = e.Manufacturer()
			return false
		}
		return true
	})
	if !sawControllerAfter {
		t.Fatal("registry missing 0x15 after structural-fallback registration; GraphQL devices will not surface the regulator")
	}
	if manufacturer != "Vaillant" {
		t.Fatalf("registered 0x15 manufacturer = %q; want %q", manufacturer, "Vaillant")
	}
}

// TestCapabilityFirstDiscovery_StructuralFallbackSkipsAdmittedSource closes
// the source-address invariant hole found during Codex adversarial review
// (PR #560 P2): when the admitted semantic source equals one of the
// structural targets (e.g. configured source 0x15, or a source-selection
// result lands on 0x15), the structural augmentation must NOT add that
// address to the probe candidates — probeB524Register would issue
// Source=0x15 / Target=0x15, a self-directed unicast probe.
func TestVaillantOwnerCapabilityFirstDiscoveryStructuralFallbackSkipsAdmittedSource(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})

	var probed []byte
	poller := newTestPoller(reg)
	poller.source = 0x15
	// 0x26 is the only coherent structural target left after 0x15 is
	// dropped (admitted source). 0x08 stays in candidates from the
	// registry but is non-coherent.
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x08: false, 0x15: true, 0x26: true}, &probed)

	addr, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v", err)
	}
	if addr != 0x26 {
		t.Fatalf("discoverB524Root = 0x%02x; want 0x26 (0x15 dropped because it equals admitted source)", addr)
	}
	for _, p := range probed {
		if p == 0x15 {
			t.Fatalf("self-directed probe issued: 0x15 was probed but is admitted source. probed=%v", probed)
		}
	}
}

// TestCapabilityFirstDiscovery_StructuralFallbackSkipsAdmittedCompanion closes
// the source-address invariant on the companion side (Codex PR #560 P2
// re-review): under source-selection, a configured / source-selected
// companion target reserved for the gateway must NOT be probed by the
// runtime semantic-discovery entry point. The startup probe path
// sanitizes via sanitizeStartupProbeTargets(... source, companion);
// discoverB524Root must apply the same guard.
func TestVaillantOwnerCapabilityFirstDiscoveryStructuralFallbackSkipsAdmittedCompanion(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})

	var probed []byte
	poller := newTestPoller(reg)
	poller.source = 0x7F
	poller.companion = 0x26 // companion reserved for the gateway

	// Only 0x15 is coherent. Without the companion guard, 0x26 would
	// be probed as a structural target (it's in the structural set).
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x08: false, 0x15: true, 0x26: true}, &probed)

	addr, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v", err)
	}
	if addr != 0x15 {
		t.Fatalf("discoverB524Root = 0x%02x; want 0x15", addr)
	}
	for _, p := range probed {
		if p == 0x26 {
			t.Fatalf("admitted companion 0x26 was probed; source-address invariant violated. probed=%v", probed)
		}
	}
}

// TestCapabilityFirstDiscovery_DropsAdmittedSourceFromRegistryCandidates closes
// Codex PR #560 P2 finding: when the admitted source is already in the
// registry (e.g. preload imported it), the registry-iteration pass would
// add it to candidates and bypass the structural-augmentation guard.
// The unified filter must drop it regardless of how it entered.
func TestVaillantOwnerCapabilityFirstDiscoveryDropsAdmittedSourceFromRegistryCandidates(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	// Preload contains the configured source address.
	reg.Register(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"})
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})

	var probed []byte
	poller := newTestPoller(reg)
	poller.source = 0x15
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x08: true, 0x15: true}, &probed)

	addr, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v", err)
	}
	if addr != 0x08 {
		t.Fatalf("discoverB524Root = 0x%02x; want 0x08 (0x15 dropped — admitted source)", addr)
	}
	for _, p := range probed {
		if p == 0x15 {
			t.Fatalf("self-directed probe issued via registry-derived candidate: 0x15 was probed but is admitted source. probed=%v", probed)
		}
	}
}

// TestCapabilityFirstDiscovery_DropsAdmittedCompanionFromRegistryCandidates
// closes the registry-side companion guard hole identified by Codex.
func TestVaillantOwnerCapabilityFirstDiscoveryDropsAdmittedCompanionFromRegistryCandidates(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	// Registry preload includes the companion target.
	reg.Register(registry.DeviceInfo{Address: 0x26, Manufacturer: "Vaillant", DeviceID: "VR_71"})
	reg.Register(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"})

	var probed []byte
	poller := newTestPoller(reg)
	poller.source = 0x7F
	poller.companion = 0x26 // companion reserved for the gateway

	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x15: true, 0x26: true}, &probed)

	addr, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v", err)
	}
	if addr != 0x15 {
		t.Fatalf("discoverB524Root = 0x%02x; want 0x15 (0x26 dropped — admitted companion)", addr)
	}
	for _, p := range probed {
		if p == 0x26 {
			t.Fatalf("companion-targeted probe issued via registry-derived candidate: 0x26 was probed. probed=%v", probed)
		}
	}
}

// TestRefreshDiscovery_RecomputesRegulatorCapabilityAfterStructuralRegistration
// pins the capability-recompute fix (Codex PR #560 P2 re-review): when
// refreshDiscovery's discoverB524Root succeeds via the structural-
// fallback path and registers the controller as a side effect, the
// regulator capability lookup must run against the post-registration
// inventory. Otherwise p.regulatorCapability stays at ControllerNone
// until the next regulator-recheck tick.
func TestVaillantOwnerRefreshDiscoveryRecomputesRegulatorCapabilityAfterStructuralRegistration(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"})

	poller := newTestPoller(reg)
	// Initial capability is ControllerNone (registry has only boiler).
	if cap := poller.findRegulatorCapability(); cap == productids.ControllerPresent {
		t.Fatalf("test setup invalid: pre-registration capability already %v", cap)
	}

	// Discover 0x15 via structural fallback, then register it.
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x08: false, 0x15: true, 0x26: false}, nil)
	controller, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v", err)
	}
	registered := poller.registerStructuralControllerIfMissing(controller)
	if !registered {
		t.Fatal("registerStructuralControllerIfMissing returned false; expected true for newly-discovered 0x15")
	}

	// findRegulatorCapability must now reflect the new inventory.
	postCap := poller.findRegulatorCapability()
	if postCap == productids.ControllerNone {
		t.Fatalf("post-registration capability = %v; want non-ControllerNone (registry now contains 0x15 Vaillant)", postCap)
	}
}

// TestRegisterStructuralControllerIfMissing_NoOpWhenAlreadyRegistered pins
// the idempotence contract: repeated structural-fallback calls for an
// already-registered controller must not duplicate registry entries or
// thrash the router plane.
func TestVaillantOwnerRegisterStructuralControllerIfMissingNoOpWhenAlreadyRegistered(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)
	reg.Register(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV2"})

	poller := newTestPoller(reg)

	refreshCount := 0
	poller.routerPlanesRefreshFn = func() { refreshCount++ }

	poller.registerStructuralControllerIfMissing(0x15)

	if refreshCount != 0 {
		t.Fatalf("router plane refreshed %d times for already-registered controller; want 0", refreshCount)
	}

	count := 0
	reg.Iterate(func(e registry.DeviceEntry) bool {
		if e != nil && e.PrimaryDisplayAddress() == 0x15 {
			count++
		}
		return true
	})
	if count != 1 {
		t.Fatalf("registry has %d entries for 0x15; want 1 (idempotent)", count)
	}
}

// TestCapabilityFirstDiscovery_StructuralFallbackEmptyRegistry covers the
// degenerate case where registry is empty — structural fallback alone
// must permit discovery to attempt the structural set rather than
// returning "no devices in registry".
func TestVaillantOwnerCapabilityFirstDiscoveryStructuralFallbackEmptyRegistry(t *testing.T) {
	t.Parallel()

	reg := registry.NewDeviceRegistry(nil)

	var probed []byte
	poller := newTestPoller(reg)
	poller.b524ProbeFn = mockB524Probe(map[byte]bool{0x15: true}, &probed)

	addr, err := poller.discoverB524Root(context.Background())
	if err != nil {
		t.Fatalf("discoverB524Root error = %v; want structural fallback to attempt 0x15", err)
	}
	if addr != 0x15 {
		t.Fatalf("discoverB524Root = 0x%02x; want 0x15", addr)
	}
}

func TestVaillantOwnerEnrichRegulatorIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		deviceID   string
		wantFamily string
	}{
		{name: "BASV2 family", deviceID: "BASV2_xx", wantFamily: "BASV2"},
		{name: "BASS2 family", deviceID: "BASS2_yy", wantFamily: "BASS2"},
		{name: "CTLV2 family", deviceID: "CTLV2_zz", wantFamily: "CTLV2"},
		{name: "unknown family", deviceID: "XYZABC", wantFamily: ""},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			reg := registry.NewDeviceRegistry(nil)
			reg.Register(registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: test.deviceID})

			poller := newTestPoller(reg)
			enrichment := poller.enrichRegulatorIdentity(0x15)
			if enrichment == nil {
				t.Fatal("enrichRegulatorIdentity returned nil")
			}
			if enrichment.family != test.wantFamily {
				t.Fatalf("enrichment.family = %q; want %q", enrichment.family, test.wantFamily)
			}
		})
	}
}

func TestVaillantOwnerDecodeB524DateSuppressSentinel(t *testing.T) {
	tests := []struct {
		name string
		raw  []byte
		want string
	}{
		{"sentinel 01.01.2015", []byte{0x01, 0x01, 0x0F}, ""},
		{"valid 15.03.2026", []byte{0x0F, 0x03, 0x1A}, "2026-03-15"},
		{"valid 01.01.2026", []byte{0x01, 0x01, 0x1A}, "2026-01-01"},
		{"short payload", []byte{0x01, 0x01}, ""},
		{"empty payload", []byte{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeB524DateSuppressSentinel(tt.raw)
			if got != tt.want {
				t.Errorf("decodeB524DateSuppressSentinel(%v) = %q; want %q", tt.raw, got, tt.want)
			}
		})
	}
}

// --- P03 startup barrier regression tests (#362) ---

func TestVaillantOwnerVaillantSemanticPollerStartDefersPollingUntilBarrier(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := newSemanticTaskScheduler()
	poller := &vaillantSemanticPoller{
		tasks:                    scheduler,
		discoveryInterval:        time.Hour,
		configInterval:           time.Hour,
		stateInterval:            time.Hour,
		energyInterval:           time.Hour,
		scheduleInterval:         time.Hour,
		regulatorRecheckInterval: time.Hour,
		nowFn:                    time.Now,
	}

	barrier := make(chan struct{})
	poller.startupBarrier = barrier
	poller.Start(ctx)

	// Allow goroutines to schedule.
	time.Sleep(50 * time.Millisecond)

	scheduler.mu.Lock()
	seqBefore := scheduler.seq
	scheduler.mu.Unlock()

	if seqBefore != 0 {
		t.Fatalf("scheduler.seq = %d before barrier; want 0 (no tasks submitted)", seqBefore)
	}

	close(barrier)

	// Allow deferred tasks to be enqueued.
	time.Sleep(50 * time.Millisecond)

	scheduler.mu.Lock()
	seqAfter := scheduler.seq
	scheduler.mu.Unlock()

	if seqAfter == 0 {
		t.Fatal("scheduler.seq = 0 after barrier; want > 0 (tasks should be submitted)")
	}
}

func TestVaillantOwnerVaillantSemanticPollerStartWithoutBarrierSubmitsImmediately(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := newSemanticTaskScheduler()
	poller := &vaillantSemanticPoller{
		tasks:                    scheduler,
		discoveryInterval:        time.Hour,
		configInterval:           time.Hour,
		stateInterval:            time.Hour,
		energyInterval:           time.Hour,
		scheduleInterval:         time.Hour,
		regulatorRecheckInterval: time.Hour,
		nowFn:                    time.Now,
	}

	poller.Start(ctx)

	time.Sleep(50 * time.Millisecond)

	scheduler.mu.Lock()
	seq := scheduler.seq
	scheduler.mu.Unlock()

	if seq == 0 {
		t.Fatal("scheduler.seq = 0 without barrier; want > 0 (tasks should be submitted immediately)")
	}
}

func TestVaillantOwnerVaillantSemanticPollerStartBarrierExitsOnContextCancel(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())

	scheduler := newSemanticTaskScheduler()
	poller := &vaillantSemanticPoller{
		tasks:                    scheduler,
		discoveryInterval:        time.Hour,
		configInterval:           time.Hour,
		stateInterval:            time.Hour,
		energyInterval:           time.Hour,
		scheduleInterval:         time.Hour,
		regulatorRecheckInterval: time.Hour,
		nowFn:                    time.Now,
	}

	barrier := make(chan struct{})
	poller.startupBarrier = barrier
	poller.Start(ctx)

	cancel()
	time.Sleep(50 * time.Millisecond)

	scheduler.mu.Lock()
	seq := scheduler.seq
	scheduler.mu.Unlock()

	if seq != 0 {
		t.Fatalf("scheduler.seq = %d after cancel; want 0 (should not start polling on cancel)", seq)
	}
}
