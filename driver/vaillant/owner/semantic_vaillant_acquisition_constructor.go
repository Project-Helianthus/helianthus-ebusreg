package owner

import (
	"fmt"
	"log"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
)

// newVaillantSemanticAcquisition builds the acquisition owner over the
// gateway's bus and registry; it reaches its publication owner only through
// publication.
func newVaillantSemanticAcquisition(cfg Config, substrate Substrate, publication vaillantSemanticPublisher) *vaillantSemanticPoller {
	catalog, catalogErr := productids.LoadCatalog()
	observeFirstFlags := normalizeObserveFirstFeatureFlagsForPoller(cfg)
	watchCatalog, _ := engine.NewWatchCatalog(nil, EngineCollaborator{})
	shadow := engine.NewShadowCache(engine.ShadowCacheOptions{
		FeatureFlags: observeFirstFlags,
		Catalog:      watchCatalog,
	})
	poller := &vaillantSemanticPoller{
		scheduler:       engine.NewSemanticReadScheduler(),
		tasks:           newSemanticTaskScheduler(),
		reg:             substrate.Registry,
		bus:             substrate.Bus,
		publication:     publication,
		shadow:          shadow,
		transportConfig: cfg.TransportConfig,
		watchObserver:   cfg.WatchObserver,
		watchEfficiency: cfg.WatchEfficiencyObserver,
		source:          cfg.ScanSource,
		companion:       cfg.StartupCompanionTarget,

		requestTimeout:            cfg.SemanticRequestTimeout,
		discoveryInterval:         cfg.SemanticDiscoveryInterval,
		configInterval:            cfg.SemanticConfigInterval,
		stateInterval:             cfg.SemanticStateInterval,
		energyInterval:            cfg.SemanticEnergyInterval,
		scheduleInterval:          10 * time.Minute,
		boilerFastInterval:        30 * time.Second,
		boilerMediumInterval:      5 * time.Minute,
		boilerSlowInterval:        10 * time.Minute,
		zoneMissThreshold:         cfg.SemanticZonePresenceMissThreshold,
		zoneHitThreshold:          cfg.SemanticZonePresenceHitThreshold,
		dhwStaleTTL:               cfg.SemanticDHWStaleTTL,
		deviceSlotRediscoveryTTL:  30 * time.Minute,
		functionModuleEvidenceTTL: defaultFunctionModuleEvidenceTTL,
		circuitFullScanInterval:   semanticCircuitFullScanInterval,

		catalog:    catalog,
		catalogErr: catalogErr,

		regulatorRecheckInterval: cfg.SemanticRegulatorRecheckInterval,
		regulatorAbsenceGrace:    cfg.SemanticRegulatorAbsenceGrace,
		regAbsenceState:          regulatorAbsenceUnknown,
		functionModuleMode:       vaillantFunctionModuleModeAbsent,

		zones:          make(map[byte]*vaillantZoneSnapshot),
		presence:       make(map[byte]*zonePresenceRecord),
		circuits:       make(map[byte]*vaillantCircuitSnapshot),
		radioDevices:   make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot),
		solarCylinders: make(map[byte]*vaillantCylinderSnapshot),
		nowFn:          time.Now,
	}
	if poller.zoneMissThreshold <= 0 {
		poller.zoneMissThreshold = DefaultSemanticZonePresenceMissThreshold
	}
	if poller.zoneHitThreshold <= 0 {
		poller.zoneHitThreshold = DefaultSemanticZonePresenceHitThreshold
	}
	if poller.dhwStaleTTL <= 0 {
		poller.dhwStaleTTL = DefaultSemanticDHWStaleTTL
	}
	if poller.energyInterval <= 0 {
		poller.energyInterval = DefaultSemanticEnergyInterval
	}
	if poller.regulatorRecheckInterval <= 0 {
		poller.regulatorRecheckInterval = DefaultSemanticRegulatorRecheckInterval
	}
	if poller.regulatorAbsenceGrace <= 0 {
		poller.regulatorAbsenceGrace = DefaultSemanticRegulatorAbsenceGrace
	}
	semanticZoneCount.Set(0)
	semanticRegulatorState.Set(string(poller.regAbsenceState))
	poller.scheduler.SetShadowCache(shadow)
	poller.scheduler.SetCircuitBreaker(engine.SemanticReadCircuitBreakerOptions{
		FailureBudget:      cfg.SemanticReadBreakerFailureBudget,
		OpenCooldown:       cfg.SemanticReadBreakerOpenCooldown,
		HalfOpenProbeLimit: cfg.SemanticReadBreakerHalfOpenProbeLimit,
		OnTransition:       poller.onSemanticReadBreakerTransition,
		OnSuppressed:       poller.onSemanticReadBreakerSuppressed,
	})
	poller.adapterInfo = newVaillantAdapterInfoState(substrate.Bus, substrate.Transport, poller.publication)
	if substrate.RefreshRouterPlanes != nil {
		poller.routerPlanesRefreshFn = func() { _ = substrate.RefreshRouterPlanes() }
	}
	poller.thermalGraph = newThermalGraphPublisher()
	return poller
}

func normalizeObserveFirstFeatureFlagsForPoller(cfg Config) engine.ObserveFirstFeatureFlags {
	if cfg.ExternalWritePolicy == "" &&
		!cfg.ObserveFirstEnabled &&
		!cfg.PassiveStateDirectApply &&
		!cfg.PassiveConfigDirectApply {
		return engine.NormalizeObserveFirstFeatureFlagsFromView(cfg.ObserveFirstFlags)
	}
	return engine.NormalizeObserveFirstFeatureFlags(
		cfg.ObserveFirstEnabled,
		cfg.PassiveStateDirectApply,
		cfg.PassiveConfigDirectApply,
		cfg.ExternalWritePolicy,
	)
}

func (p *vaillantSemanticPoller) onSemanticReadBreakerTransition(event engine.SemanticReadCircuitBreakerTransition) {
	if p == nil {
		return
	}
	semanticReadBreakerTransitionsTotal.Add(fmt.Sprintf("%s->%s", event.From, event.To), 1)
	log.Printf(
		"semantic_read_breaker_transition key=%q from=%s to=%s consecutive_failures=%d",
		event.Key,
		event.From,
		event.To,
		event.ConsecutiveFailures,
	)
}

func (p *vaillantSemanticPoller) onSemanticReadBreakerSuppressed(event engine.SemanticReadCircuitBreakerSuppression) {
	if p == nil {
		return
	}
	semanticReadBreakerSuppressedTotal.Add(string(event.State), 1)
	log.Printf(
		"semantic_read_breaker_suppressed key=%q state=%s retry_after=%s suppressed_total=%d",
		event.Key,
		event.State,
		event.RetryAfter.Round(time.Millisecond),
		event.SuppressedTotal,
	)
}

// HydrateFromCache seeds the retained zones, hot water and boiler from the
// native cache hydration value and takes ownership of its snapshots.
func (p *vaillantSemanticPoller) HydrateFromCache(hydration VaillantCacheHydration) {
	if p == nil {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.zones == nil {
		p.zones = make(map[byte]*vaillantZoneSnapshot)
	}
	if p.presence == nil {
		p.presence = make(map[byte]*zonePresenceRecord)
	}
	for _, zone := range hydration.zones {
		p.zones[zone.Instance] = zone
		p.presence[zone.Instance] = &zonePresenceRecord{
			State:     zonePresencePresent,
			HitStreak: p.zoneHitThresholdValue(),
		}
	}
	if hydration.dhw != nil {
		p.dhw = hydration.dhw
		if !hydration.persistedAt.IsZero() {
			p.dhwLastUpdateAt = hydration.persistedAt.UTC()
		} else {
			p.markDHWUpdatedNowLocked()
		}
	}
	if hydration.boiler != nil {
		p.boiler = hydration.boiler
	}
	semanticZoneCount.Set(int64(len(p.zones)))
}
