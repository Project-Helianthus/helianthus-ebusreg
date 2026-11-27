package owner

import (
	"fmt"
	"log"
	"time"
)

func (provider *retainedSemanticState) recordEpochUpdateLocked(source semanticDataSource, stream semanticStreamKind, reason string) *phaseTransitionLog {
	switch source {
	case semanticDataSourceCache:
		provider.cacheEpoch++
		semanticCacheEpoch.Set(int64(provider.cacheEpoch))
		provider.touchStartupUpdatedAtLocked()
		if provider.liveEpoch == 0 && provider.phase != semanticStartupPhaseDegraded {
			return provider.transitionPhaseLocked(semanticStartupPhaseCacheLoadedStale, reason)
		}
	case semanticDataSourceLive:
		switch stream {
		case semanticStreamZones:
			provider.zoneLiveSeen = true
		case semanticStreamDHW:
			provider.dhwLiveSeen = true
		case semanticStreamBoiler:
			provider.boilerLiveSeen = true
		case semanticStreamCircuits:
			provider.circuitLiveSeen = true
		case semanticStreamRadioDevices:
			provider.radioLiveSeen = true
		}
		provider.liveEpoch++
		semanticLiveEpoch.Set(int64(provider.liveEpoch))
		provider.touchStartupUpdatedAtLocked()
		if provider.liveEpoch == 1 {
			return provider.transitionPhaseLocked(semanticStartupPhaseLiveWarmup, reason)
		}
		if provider.liveReadyCriteriaLocked() {
			return provider.transitionPhaseLocked(semanticStartupPhaseLiveReady, reason)
		}
	}
	return nil
}

func (provider *retainedSemanticState) liveReadyCriteriaLocked() bool {
	if provider.zonePublished && !provider.zoneLiveSeen {
		return false
	}
	if provider.dhwPublished && !provider.dhwLiveSeen {
		return false
	}
	if provider.boilerPublished && !provider.boilerLiveSeen {
		return false
	}
	if provider.circuitPublished && !provider.circuitLiveSeen {
		return false
	}
	if provider.radioPublished && !provider.radioLiveSeen {
		return false
	}
	return provider.zoneLiveSeen || provider.dhwLiveSeen || provider.boilerLiveSeen || provider.circuitLiveSeen || provider.radioLiveSeen
}

func (provider *retainedSemanticState) transitionPhaseLocked(next semanticStartupPhase, reason string) *phaseTransitionLog {
	if next == "" || provider.phase == next {
		return nil
	}
	previous := provider.phase
	provider.phase = next
	provider.touchStartupUpdatedAtLocked()
	semanticStartupPhaseTransitionsTotal.Add(fmt.Sprintf("%s->%s", previous, next), 1)
	semanticStartupCurrentPhase.Set(string(next))
	return &phaseTransitionLog{
		previous:   previous,
		next:       next,
		reason:     reason,
		cacheEpoch: provider.cacheEpoch,
		liveEpoch:  provider.liveEpoch,
	}
}

func (provider *retainedSemanticState) touchStartupUpdatedAtLocked() {
	if provider == nil {
		return
	}
	now := time.Now().UTC()
	if provider.startupUpdatedAt.IsZero() || now.After(provider.startupUpdatedAt) {
		provider.startupUpdatedAt = now
	}
}

func (provider *retainedSemanticState) logPhaseTransition(transition *phaseTransitionLog) {
	if provider == nil || transition == nil {
		return
	}
	logger := provider.phaseLogger
	if logger == nil {
		logger = log.Printf
	}
	logger(
		"semantic_startup_phase_transition from=%s to=%s reason=%q cache_epoch=%d live_epoch=%d",
		transition.previous,
		transition.next,
		transition.reason,
		transition.cacheEpoch,
		transition.liveEpoch,
	)
}
