package owner

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

// extractPartNumberFromSerial extracts a 10-digit Vaillant part number from a
// device serial number string. The logic mirrors graphql.extractVaillantPartNumber.
func extractPartNumberFromSerial(serial string) string {
	serial = strings.TrimSpace(serial)
	if serial == "" {
		return ""
	}

	// Try dash-separated format: "YY-WW-CC-PPPPPPPPPP-SSSS-RRRRRR-XX"
	parts := strings.Split(serial, "-")
	if len(parts) >= 4 {
		partNumber := strings.TrimSpace(parts[3])
		if len(partNumber) == 10 && isAllDigits(partNumber) {
			return partNumber
		}
	}

	// Fallback: strip non-alphanumeric, take digits 6..16.
	compact := make([]byte, 0, len(serial))
	for i := 0; i < len(serial); i++ {
		ch := serial[i]
		if (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') {
			compact = append(compact, ch)
		}
	}
	if len(compact) >= 16 {
		candidate := string(compact[6:16])
		if len(candidate) == 10 && isAllDigits(candidate) {
			return candidate
		}
	}

	return ""
}

// isAllDigits returns true if the string is non-empty and contains only ASCII digits.
func isAllDigits(value string) bool {
	if value == "" {
		return false
	}
	for i := 0; i < len(value); i++ {
		if value[i] < '0' || value[i] > '9' {
			return false
		}
	}
	return true
}

// b524DiscoveryOptions controls how discoverB524RootWithOptions builds its
// candidate set. The default zero value (augmentStructural=false) is the
// registry-health-check semantic — answers "does the existing registry
// contain a coherent B524 root?" without probing addresses outside the
// registry. Set augmentStructural=true for the production discovery
// path used by refreshDiscovery, which unconditionally probes the
// bounded Vaillant structural target set even when the registry has
// not yet learned about the regulator.
type b524DiscoveryOptions struct {
	augmentStructural bool
}

// discoverB524Root is the production discovery path. It augments the
// registry-derived candidate list with the bounded Vaillant structural
// target set {0x08, 0x15, 0x26}, probes D1-ordered (0x15 first), and
// stops at the first coherent candidate.
//
// Structural augmentation is unconditional and idempotent: the dedup
// map collapses overlap, and the structural set is bounded to three
// non-broadcast unicast targets. The alternative — gating on registry
// size — has a hostile-registry hole where any stray non-Vaillant
// entry (e.g. {0x08 boiler, 0xEC passive-only responder}) suppresses
// the fallback exactly when it is needed.
//
// Probing structural addresses that are not yet in the registry is
// safe: the unicast probe either elicits a coherent response (in
// which case refreshDiscovery registers the controller post-
// discovery via registerStructuralControllerIfMissing) or it does
// not. Source-address invariant is upheld — the probe uses the
// existing admitted source via probeFn.
func (p *vaillantSemanticPoller) discoverB524Root(ctx context.Context) (byte, error) {
	return p.discoverB524RootWithOptions(ctx, b524DiscoveryOptions{augmentStructural: true})
}

// DiscoverB524RootInRegistry is the registry-health-check path. It
// considers ONLY candidates already in the registry. Used by startup
// scan helpers to ask "is the existing registry coherent?", which is
// a distinct question from "where is the regulator?".
func (p *vaillantSemanticPoller) DiscoverB524RootInRegistry(ctx context.Context) (byte, error) {
	return p.discoverB524RootWithOptions(ctx, b524DiscoveryOptions{augmentStructural: false})
}

func (p *vaillantSemanticPoller) discoverB524RootWithOptions(ctx context.Context, opts b524DiscoveryOptions) (byte, error) {
	if p == nil || p.reg == nil {
		return 0, fmt.Errorf("gateway.discoverB524Root: nil poller or registry")
	}

	probeFn := p.b524ProbeFn
	if probeFn == nil {
		probeFn = p.probeB524Register
	}

	// Source-address invariant applies uniformly across both
	// candidate entry points (registry iteration AND structural
	// fallback). A configured source / companion target that happens
	// to be already present in the registry — e.g. an ebusd preload
	// imports 0x26 before semantic discovery runs, or a prior
	// session left the companion address registered — would otherwise
	// bypass the structural-augmentation guard via the dedup map and
	// be probed unicast. The shared filter ensures any candidate
	// matching the admitted source or companion is dropped regardless
	// of how it entered the candidate set.
	skipReservedSourceCompanion := func(addr byte) bool {
		reservedSource, reservedCompanion := p.sourceAndCompanion()
		if addr == reservedSource {
			return true
		}
		if reservedCompanion != 0 && addr == reservedCompanion {
			return true
		}
		return false
	}
	candidateSet := make(map[byte]struct{})
	var candidates []byte
	// P9.3 — race-free routing-address enumeration via
	// IterateSnapshots + SnapshotTargetAddressForRouting. The
	// previous Iterate path called entry.AddressByRole and
	// entry.PrimaryDisplayAddress lock-free.
	p.reg.IterateSnapshots(func(snap registry.DeviceEntrySnapshot) bool {
		// Phase C M-C6b: B524 candidates are routing targets, not
		// display addresses. An aliased regulator entry (e.g.
		// 0x10↔0x15 displayed as 0x10) must be probed at 0x15
		// where the responder lives — probing 0x10 misses the
		// coherent responder and reports no B524 root.
		addr := snapshotTargetAddressForRouting(snap)
		if skipReservedSourceCompanion(addr) {
			return true
		}
		if _, dup := candidateSet[addr]; !dup {
			candidateSet[addr] = struct{}{}
			candidates = append(candidates, addr)
		}
		return true
	})

	if opts.augmentStructural {
		for _, addr := range vaillantStructuralStartupProbeTargets {
			if _, dup := candidateSet[addr]; dup {
				continue
			}
			if skipReservedSourceCompanion(addr) {
				continue
			}
			candidateSet[addr] = struct{}{}
			candidates = append(candidates, addr)
		}
	}

	if len(candidates) == 0 {
		return 0, fmt.Errorf("gateway.discoverB524Root: no devices in registry")
	}

	slices.Sort(candidates)

	if _, hasRegulatorAddress := candidateSet[0x15]; hasRegulatorAddress {
		ordered := make([]byte, 0, len(candidates))
		ordered = append(ordered, 0x15)
		for _, c := range candidates {
			if c != 0x15 {
				ordered = append(ordered, c)
			}
		}
		candidates = ordered
	}

	for _, addr := range candidates {
		coherent := true
		for _, probe := range b524CapabilityProbes {
			if !probeFn(ctx, addr, probe.opcode, probe.group, probe.instance, probe.addr) {
				coherent = false
				break
			}
		}
		if coherent {
			log.Printf("semantic_b524_root_discovery address=0x%02x probed=%d", addr, len(candidates))
			return addr, nil
		}
	}

	addrs := make([]string, len(candidates))
	for i, c := range candidates {
		addrs[i] = fmt.Sprintf("0x%02x", c)
	}
	return 0, fmt.Errorf("gateway.discoverB524Root: no coherent responder among [%s]", strings.Join(addrs, ", "))
}
