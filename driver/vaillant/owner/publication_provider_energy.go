package owner

import (
	"strings"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

// SetEnergyClock injects the reference clock used to ingest broadcasts and to
func (provider *retainedSemanticState) SetEnergyClock(clock func() time.Time) {
	if provider == nil {
		return
	}
	provider.mu.Lock()
	provider.energyClock = clock
	provider.mu.Unlock()
}

func (provider *retainedSemanticState) energyNow() time.Time {
	provider.mu.RLock()
	clock := provider.energyClock
	provider.mu.RUnlock()
	if clock == nil {
		return time.Now()
	}
	return clock()
}

// EnergyWindowKey resolves a relative register window ("current"/"previous"
// for year and month periods; the reference day for "day") to the dated merge
// key against the provider's reference clock.
func (provider *retainedSemanticState) EnergyWindowKey(channel, usage, period, label string) energyMergeKey {
	return energyWindowKey(channel, usage, period, label, provider.energyNow())
}

// EnergyReferenceTime returns the provider's current reference instant. A
// register acquisition captures it once and derives every key from it.
func (provider *retainedSemanticState) EnergyReferenceTime() time.Time {
	return provider.energyNow()
}

// EnergyWindowKeyAt is EnergyWindowKey against an explicit reference instant.
func (provider *retainedSemanticState) EnergyWindowKeyAt(channel, usage, period, label string, ref time.Time) energyMergeKey {
	return energyWindowKey(channel, usage, period, label, ref)
}

// EnergyLifetimeKey is the undated key of a lifetime (all-time) total. It is
// never keyed by a reference year; the public series maps it to the current
// yearly slot with the previous slot and today held at zero.
func (provider *retainedSemanticState) EnergyLifetimeKey(channel, usage string) energyMergeKey {
	return energyLifetimeKey(channel, usage)
}

func (provider *retainedSemanticState) EnergyTotals() *output.EnergyTotals {
	if provider == nil || provider.energyMerge == nil {
		return nil
	}
	snapshot := provider.energyMerge.SnapshotWithContext(provider.energyNow(), provider.currentEnergyPassiveState())
	if snapshot == nil {
		return nil
	}
	copy := *snapshot
	copy.Gas = cloneEnergyChannel(copy.Gas)
	copy.Electric = cloneEnergyChannel(copy.Electric)
	copy.Solar = cloneEnergyChannel(copy.Solar)
	return &copy
}

func (provider *retainedSemanticState) SetEnergyPassiveState(state string) {
	if provider == nil {
		return
	}
	provider.mu.Lock()
	provider.energyPassiveState = strings.TrimSpace(state)
	provider.mu.Unlock()
}

func (provider *retainedSemanticState) SetEnergyPassiveStateProvider(providerFn func() string) {
	if provider == nil {
		return
	}
	provider.mu.Lock()
	provider.energyPassiveStateProvider = providerFn
	provider.mu.Unlock()
}

func (provider *retainedSemanticState) RefreshEnergyFreshnessMetrics(now time.Time, passiveState string) {
	if provider == nil || provider.energyMerge == nil {
		return
	}
	state := strings.TrimSpace(passiveState)
	if state == "" {
		state = provider.currentEnergyPassiveState()
	}
	provider.energyMerge.RefreshFreshnessMetricsWithContext(now, state)
}

func (provider *retainedSemanticState) currentEnergyPassiveState() string {
	if provider == nil {
		return ""
	}
	provider.mu.RLock()
	providerFn := provider.energyPassiveStateProvider
	fallback := provider.energyPassiveState
	provider.mu.RUnlock()
	if providerFn == nil {
		return fallback
	}
	if value := strings.TrimSpace(providerFn()); value != "" {
		return value
	}
	return fallback
}

// ApplyEnergyFromRegister updates semantic energy snapshots from register reads.
// Returns true when the merge store accepted the value.
func (provider *retainedSemanticState) ApplyEnergyFromRegister(key energyMergeKey, kwh float64) bool {
	if provider == nil || provider.energyMerge == nil {
		return false
	}
	return provider.applyEnergyPoint(key, kwh, energySourceRegister, provider.energyNow())
}
func (provider *retainedSemanticState) applyEnergyPoint(key energyMergeKey, kwh float64, source energyDataSource, now time.Time) bool {
	if provider == nil || provider.energyMerge == nil {
		return false
	}
	if !provider.energyMerge.Apply(key, kwh, source, now) {
		return false
	}

	// Rebuild the snapshot from the merge store and publish atomically.
	// Only publish if our revision is still the latest (prevents stale overwrites
	// from concurrent callers).
	rev := provider.energyMerge.Revision()
	snapshot := provider.energyMerge.SnapshotWithContext(now, provider.currentEnergyPassiveState())
	provider.mu.Lock()
	if rev >= provider.energyRevision {
		provider.energy = snapshot
		provider.energyRevision = rev
	}
	provider.mu.Unlock()
	return true
}

// matchesToday requires a valid absolute year, month and day equal to the
// reference date.
