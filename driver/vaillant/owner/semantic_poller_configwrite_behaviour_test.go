package owner

import (
	"context"
	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
	"testing"
	"time"
)

func TestVaillantOwnerMergeCylinderSnapshotMapNonDestructivePreservesTemperatureWhileRefreshingConfig(t *testing.T) {
	t.Parallel()

	oldTemp := 49.0
	oldSetpoint := 55.0
	newSetpoint := 60.0
	existing := map[byte]*vaillantCylinderSnapshot{
		0x00: {
			Instance:     0x00,
			TemperatureC: &oldTemp,
			MaxSetpointC: &oldSetpoint,
		},
	}
	incoming := map[byte]*vaillantCylinderSnapshot{
		0x00: {
			Instance:     0x00,
			MaxSetpointC: &newSetpoint,
		},
	}

	merged := mergeCylinderSnapshotMapNonDestructive(existing, incoming)
	snapshot := merged[0x00]
	if snapshot == nil {
		t.Fatal("merged[0x00] = nil; want preserved cylinder snapshot")
	}
	if snapshot.TemperatureC == nil || *snapshot.TemperatureC != oldTemp {
		t.Fatalf("merged temperature = %#v; want preserved %v", snapshot.TemperatureC, oldTemp)
	}
	if snapshot.MaxSetpointC == nil || *snapshot.MaxSetpointC != newSetpoint {
		t.Fatalf("merged max setpoint = %#v; want refreshed %v", snapshot.MaxSetpointC, newSetpoint)
	}
}

func TestVaillantOwnerVaillantSemanticPollerDHWTransientCacheFailurePreservesLastKnown(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.January, 10, 8, 0, 0, 0, time.UTC)
	current := 48.5

	provider := newRetainedSemanticState()
	poller := &vaillantSemanticPoller{
		publication: newVaillantSemanticPublication(provider, nil, nil),
		dhw: &vaillantDhwSnapshot{
			OperatingMode: "auto",
			Preset:        "schedule",
			CurrentTempC:  &current,
		},
		dhwStaleTTL: 10 * time.Minute,
	}
	poller.nowFn = func() time.Time { return now }
	poller.dhwLastUpdateAt = now

	poller.publishDHW(semanticSnapshotSourceLive)

	now = now.Add(5 * time.Minute)
	poller.publishDHW(semanticSnapshotSourceCache)

	dhw := provider.DHW()
	if dhw == nil {
		t.Fatalf("provider.DHW() = nil; want preserved last-known DHW before TTL expiry")
	}
	if dhw.Config.OperatingMode != "auto" {
		t.Fatalf("provider.DHW().Config.OperatingMode = %q; want auto", dhw.Config.OperatingMode)
	}
	if dhw.State.CurrentTempC == nil || *dhw.State.CurrentTempC != 48.5 {
		t.Fatalf("provider.DHW().State.CurrentTempC = %v; want 48.5", dhw.State.CurrentTempC)
	}
}

func TestVaillantOwnerVaillantSemanticPollerDHWCacheFailureExpiresAfterTTL(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.January, 10, 8, 0, 0, 0, time.UTC)
	current := 47.2

	provider := newRetainedSemanticState()
	cacheSpy := &semanticSnapshotCaptureSpy{}
	poller := &vaillantSemanticPoller{
		publication: newVaillantSemanticPublication(provider, nil, cacheSpy),
		dhw: &vaillantDhwSnapshot{
			OperatingMode: "auto",
			Preset:        "schedule",
			CurrentTempC:  &current,
		},
		dhwStaleTTL: 10 * time.Minute,
	}
	poller.nowFn = func() time.Time { return now }
	poller.dhwLastUpdateAt = now

	poller.publishDHW(semanticSnapshotSourceLive)

	now = now.Add(11 * time.Minute)
	poller.publishDHW(semanticSnapshotSourceCache)

	if dhw := provider.DHW(); dhw != nil {
		t.Fatalf("provider.DHW() = %#v; want nil after TTL expiry", dhw)
	}
	if poller.dhw != nil {
		t.Fatalf("poller.dhw = %#v; want nil after TTL expiry", poller.dhw)
	}
	if cacheSpy.calls < 2 {
		t.Fatalf("cache Save calls = %d; want at least 2 (live persist + expiry persist)", cacheSpy.calls)
	}
	if cacheSpy.last.DHW != nil {
		t.Fatalf("cache last DHW = %#v; want nil after TTL expiry persist", cacheSpy.last.DHW)
	}
}

func TestVaillantOwnerVaillantSemanticPollerHydrateFromCacheSeedsDHWStalenessFromPersistedAt(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.January, 10, 8, 0, 0, 0, time.UTC)
	current := 45.0

	provider := newRetainedSemanticState()
	cacheSpy := &semanticSnapshotCaptureSpy{}
	poller := &vaillantSemanticPoller{
		publication: newVaillantSemanticPublication(provider, nil, cacheSpy),
		dhwStaleTTL: 15 * time.Minute,
	}
	poller.nowFn = func() time.Time { return now }

	// Hydrate from a cache that was persisted 20 minutes ago — already past the 15m TTL.
	poller.HydrateFromCache(cacheHydrationFromSnapshot(output.SemanticCacheSnapshot{
		DHW: &output.DhwStatus{
			Config: output.DhwConfig{
				OperatingMode: "auto",
				Preset:        "schedule",
			},
			State: output.DhwState{
				CurrentTempC: &current,
			},
		},
		PersistedAt: now.Add(-20 * time.Minute),
	}))

	// The poller should have seeded dhwLastUpdateAt from PersistedAt, not now().
	if poller.dhwLastUpdateAt.Equal(now) {
		t.Fatalf("dhwLastUpdateAt = now; want PersistedAt (20m ago)")
	}

	// First cache-sourced publish should expire the DHW immediately since age > TTL.
	poller.publishDHW(semanticSnapshotSourceCache)

	if dhw := provider.DHW(); dhw != nil {
		t.Fatalf("provider.DHW() = %#v; want nil — cache was already past TTL at hydration", dhw)
	}
	if poller.dhw != nil {
		t.Fatalf("poller.dhw = %#v; want nil after immediate TTL expiry", poller.dhw)
	}
}

func TestVaillantOwnerVaillantSemanticPollerHydrateFromCacheFallsBackToNowWhenPersistedAtZero(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.January, 10, 8, 0, 0, 0, time.UTC)
	current := 45.0

	poller := &vaillantSemanticPoller{
		publication: newVaillantSemanticPublication(newRetainedSemanticState(), nil, nil),
		dhwStaleTTL: 15 * time.Minute,
	}
	poller.nowFn = func() time.Time { return now }

	// Hydrate with zero PersistedAt — should fall back to now().
	poller.HydrateFromCache(cacheHydrationFromSnapshot(output.SemanticCacheSnapshot{
		DHW: &output.DhwStatus{
			Config: output.DhwConfig{
				OperatingMode: "auto",
			},
			State: output.DhwState{
				CurrentTempC: &current,
			},
		},
	}))

	if !poller.dhwLastUpdateAt.Equal(now) {
		t.Fatalf("dhwLastUpdateAt = %v; want %v (fallback to now when PersistedAt is zero)", poller.dhwLastUpdateAt, now)
	}
}

func TestVaillantOwnerVaillantSemanticPollerDiscoveryLossTTLExpiresDHW(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.January, 11, 9, 0, 0, 0, time.UTC)
	current := 49.1

	provider := newRetainedSemanticState()
	cacheSpy := &semanticSnapshotCaptureSpy{}
	poller := &vaillantSemanticPoller{
		publication: newVaillantSemanticPublication(provider, nil, cacheSpy),
		dhwStaleTTL: 10 * time.Minute,
		dhw: &vaillantDhwSnapshot{
			OperatingMode: "auto",
			Preset:        "schedule",
			CurrentTempC:  &current,
		},
	}
	poller.nowFn = func() time.Time { return now }
	poller.dhwLastUpdateAt = now

	poller.publishDHW(semanticSnapshotSourceLive)

	now = now.Add(5 * time.Minute)
	poller.refreshDiscovery(context.Background())
	if dhw := provider.DHW(); dhw == nil {
		t.Fatalf("provider.DHW() = nil; want DHW preserved before TTL expiry")
	}

	now = now.Add(6 * time.Minute)
	poller.refreshDiscovery(context.Background())

	if dhw := provider.DHW(); dhw != nil {
		t.Fatalf("provider.DHW() = %#v; want nil after TTL expiry on discovery loss", dhw)
	}
	if poller.dhw != nil {
		t.Fatalf("poller.dhw = %#v; want nil after TTL expiry on discovery loss", poller.dhw)
	}
	if cacheSpy.calls < 2 {
		t.Fatalf("cache Save calls = %d; want at least 2 (live persist + expiry persist)", cacheSpy.calls)
	}
	if cacheSpy.last.DHW != nil {
		t.Fatalf("cache last DHW = %#v; want nil after TTL expiry persist", cacheSpy.last.DHW)
	}
}

func TestVaillantOwnerVaillantSemanticPollerDiscoveryLossPreservesDHWBeforeTTL(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.January, 11, 10, 0, 0, 0, time.UTC)
	current := 47.8

	provider := newRetainedSemanticState()
	cacheSpy := &semanticSnapshotCaptureSpy{}
	poller := &vaillantSemanticPoller{
		publication: newVaillantSemanticPublication(provider, nil, cacheSpy),
		dhwStaleTTL: 15 * time.Minute,
		dhw: &vaillantDhwSnapshot{
			OperatingMode: "auto",
			Preset:        "schedule",
			CurrentTempC:  &current,
		},
	}
	poller.nowFn = func() time.Time { return now }
	poller.dhwLastUpdateAt = now

	poller.publishDHW(semanticSnapshotSourceLive)

	now = now.Add(14 * time.Minute)
	poller.refreshDiscovery(context.Background())

	dhw := provider.DHW()
	if dhw == nil {
		t.Fatalf("provider.DHW() = nil; want DHW preserved before TTL expiry")
	}
	if dhw.Config.OperatingMode != "auto" {
		t.Fatalf("provider.DHW().Config.OperatingMode = %q; want auto", dhw.Config.OperatingMode)
	}
	if dhw.State.CurrentTempC == nil || *dhw.State.CurrentTempC != 47.8 {
		t.Fatalf("provider.DHW().State.CurrentTempC = %v; want 47.8", dhw.State.CurrentTempC)
	}
	if poller.dhw == nil {
		t.Fatalf("poller.dhw = nil; want preserved snapshot before TTL expiry")
	}
	if cacheSpy.last.DHW == nil {
		t.Fatalf("cache last DHW = nil; want last persisted live DHW snapshot")
	}
}

func TestVaillantOwnerVaillantSemanticPollerDiscoveryFlapDuringStartup(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.January, 11, 11, 0, 0, 0, time.UTC)
	current := 46.6
	cached := &output.DhwStatus{
		Config: output.DhwConfig{
			OperatingMode: "auto",
			Preset:        "schedule",
		},
		State: output.DhwState{
			CurrentTempC: &current,
		},
	}

	provider := newRetainedSemanticState()
	provider.SetDHWFromCache(cached)
	if got := provider.StartupPhase(); got != semanticStartupPhaseCacheLoadedStale {
		t.Fatalf("phase after cache load = %s; want %s", got, semanticStartupPhaseCacheLoadedStale)
	}

	poller := &vaillantSemanticPoller{
		publication: newVaillantSemanticPublication(provider, nil, nil),
		dhwStaleTTL: 30 * time.Minute,
	}
	poller.nowFn = func() time.Time { return now }
	poller.HydrateFromCache(cacheHydrationFromSnapshot(output.SemanticCacheSnapshot{
		DHW:         cached,
		PersistedAt: now.Add(-5 * time.Minute),
	}))

	poller.refreshDiscovery(context.Background())
	if got := provider.StartupPhase(); got != semanticStartupPhaseCacheLoadedStale {
		t.Fatalf("phase after cache-phase discovery flap = %s; want %s", got, semanticStartupPhaseCacheLoadedStale)
	}
	if dhw := provider.DHW(); dhw == nil {
		t.Fatalf("provider.DHW() = nil; want preserved cached DHW during CACHE_LOADED_STALE discovery flap")
	}

	now = now.Add(1 * time.Minute)
	poller.publishDHW(semanticSnapshotSourceLive)
	if got := provider.StartupPhase(); got != semanticStartupPhaseLiveWarmup {
		t.Fatalf("phase after first live DHW publish = %s; want %s", got, semanticStartupPhaseLiveWarmup)
	}

	now = now.Add(1 * time.Minute)
	poller.refreshDiscovery(context.Background())
	if got := provider.StartupPhase(); got != semanticStartupPhaseLiveWarmup {
		t.Fatalf("phase after warmup discovery flap = %s; want %s", got, semanticStartupPhaseLiveWarmup)
	}
	dhw := provider.DHW()
	if dhw == nil {
		t.Fatalf("provider.DHW() = nil; want preserved DHW during LIVE_WARMUP discovery flap")
	}
	if dhw.State.CurrentTempC == nil || *dhw.State.CurrentTempC != 46.6 {
		t.Fatalf("provider.DHW().State.CurrentTempC = %v; want 46.6", dhw.State.CurrentTempC)
	}
}

func TestVaillantOwnerMergeZoneSnapshotFieldsPartialLiveUpdatePreservesLastKnownAndFreshness(t *testing.T) {
	t.Parallel()

	current := 21.0
	target := 22.5
	humidity := 44.0
	roomTemperatureZoneMapping := uint16(2)
	circuitIndex := uint16(2)
	circuitType := uint16(1)
	valve := uint16(1)
	entry := &vaillantZoneSnapshot{
		Name:                              "Zone 1",
		OperatingMode:                     "heat",
		Preset:                            "manual",
		HvacAction:                        "heating",
		AllowedModes:                      []string{"off", "auto", "heat"},
		CurrentTempC:                      &current,
		TargetTempC:                       &target,
		HumidityPct:                       &humidity,
		ConfigurationHeatingOperationMode: "1",
		StateSpecialFunction:              "0",
		ConfigurationRoomTemperatureZoneMappingRaw: &roomTemperatureZoneMapping,
		ConfigurationAssociatedCircuitRaw:          &circuitIndex,
		ConfigurationCircuitTypeRaw:                &circuitType,
		StateValveStatusRaw:                        &valve,
	}
	seedZoneFreshness(entry, semanticSnapshotSourceCache, true)

	updatedCurrent := 21.8
	incoming := &vaillantZoneSnapshot{
		CurrentTempC:                      &updatedCurrent,
		ConfigurationHeatingOperationMode: "2",
	}
	mergeZoneSnapshotFields(entry, incoming, semanticSnapshotSourceLive, newSemanticFieldSet(
		zoneFieldOperatingMode,
		zoneFieldPreset,
		zoneFieldHvacAction,
		zoneFieldAllowedModes,
		zoneFieldCurrentTempC,
		zoneFieldTargetTempC,
		zoneFieldCurrentHumidityPct,
		zoneFieldZoneOperationModeRaw,
		zoneFieldSpecialFunctionRaw,
		zoneFieldRoomTemperatureZoneMappingRaw,
		zoneFieldZoneCircuitIndexRaw,
		zoneFieldZoneValveStatusRaw,
		zoneFieldCircuitTypeRaw,
		zoneFieldQuickVetoTempC,
		zoneFieldQuickVetoDurationH,
		zoneFieldQuickVetoEndTime,
		zoneFieldQuickVetoEndDate,
		zoneFieldHolidayStartDate,
		zoneFieldHolidayEndDate,
		zoneFieldHolidaySetpointC,
		zoneFieldHolidayStartTime,
		zoneFieldHolidayEndTime,
	))

	if entry.CurrentTempC == nil || *entry.CurrentTempC != 21.8 {
		t.Fatalf("entry.CurrentTempC = %v; want 21.8", entry.CurrentTempC)
	}
	if entry.TargetTempC == nil || *entry.TargetTempC != 22.5 {
		t.Fatalf("entry.TargetTempC = %v; want preserved 22.5", entry.TargetTempC)
	}
	if entry.ConfigurationAssociatedCircuitRaw == nil || *entry.ConfigurationAssociatedCircuitRaw != 2 {
		t.Fatalf("entry.ConfigurationAssociatedCircuitRaw = %v; want preserved 2", entry.ConfigurationAssociatedCircuitRaw)
	}
	if entry.ConfigurationRoomTemperatureZoneMappingRaw == nil || *entry.ConfigurationRoomTemperatureZoneMappingRaw != 2 {
		t.Fatalf("entry.ConfigurationRoomTemperatureZoneMappingRaw = %v; want preserved 2", entry.ConfigurationRoomTemperatureZoneMappingRaw)
	}
	if entry.ConfigurationHeatingOperationMode != "2" {
		t.Fatalf("entry.ConfigurationHeatingOperationMode = %q; want 2", entry.ConfigurationHeatingOperationMode)
	}

	currentFreshness, ok := entry.FieldFreshness[zoneFieldCurrentTempC]
	if !ok || currentFreshness.Source != semanticSnapshotSourceLive || currentFreshness.Stale {
		t.Fatalf("current_temp freshness = %+v (ok=%v); want source=live stale=false", currentFreshness, ok)
	}
	targetFreshness, ok := entry.FieldFreshness[zoneFieldTargetTempC]
	if !ok || targetFreshness.Source != semanticSnapshotSourceCache || !targetFreshness.Stale {
		t.Fatalf("target_temp freshness = %+v (ok=%v); want source=cache stale=true", targetFreshness, ok)
	}
	opModeFreshness, ok := entry.FieldFreshness[zoneFieldZoneOperationModeRaw]
	if !ok || opModeFreshness.Source != semanticSnapshotSourceLive || opModeFreshness.Stale {
		t.Fatalf("operation_mode_raw freshness = %+v (ok=%v); want source=live stale=false", opModeFreshness, ok)
	}
	roomMappingFreshness, ok := entry.FieldFreshness[zoneFieldRoomTemperatureZoneMappingRaw]
	if !ok || roomMappingFreshness.Source != semanticSnapshotSourceCache || !roomMappingFreshness.Stale {
		t.Fatalf("room_temperature_zone_mapping_raw freshness = %+v (ok=%v); want source=cache stale=true", roomMappingFreshness, ok)
	}
}

func TestVaillantOwnerVaillantSemanticPollerPublishZonesSeparatesRoomTemperatureZoneMappingFromAssociatedCircuit(t *testing.T) {
	t.Parallel()

	provider := newRetainedSemanticState()
	roomTemperatureZoneMapping := uint16(2)
	associatedCircuit := uint16(resolveAssociatedCircuitInstance(&roomTemperatureZoneMapping, 1))

	poller := &vaillantSemanticPoller{
		publication: newVaillantSemanticPublication(provider, nil, nil),
		zones: map[byte]*vaillantZoneSnapshot{
			0x01: {
				Instance:      0x01,
				Present:       true,
				Name:          "Etaj",
				OperatingMode: "heat",
				Preset:        "manual",
				AllowedModes:  []string{"off", "auto", "heat"},
				ConfigurationRoomTemperatureZoneMappingRaw: &roomTemperatureZoneMapping,
				ConfigurationAssociatedCircuitRaw:          &associatedCircuit,
			},
		},
	}

	poller.publishZones(semanticSnapshotSourceLive)

	zones := provider.Zones()
	if len(zones) != 1 {
		t.Fatalf("len(provider.Zones()) = %d; want 1", len(zones))
	}
	if zones[0].Config.RoomTemperatureZoneMapping == nil || *zones[0].Config.RoomTemperatureZoneMapping != 2 {
		t.Fatalf("roomTemperatureZoneMapping = %#v; want 2", zones[0].Config.RoomTemperatureZoneMapping)
	}
	if zones[0].Config.AssociatedCircuit == nil || *zones[0].Config.AssociatedCircuit != 1 {
		t.Fatalf("associatedCircuit = %#v; want 1", zones[0].Config.AssociatedCircuit)
	}
}

func TestVaillantOwnerVaillantSemanticPollerPublishZonesFillsRoomStateFromMappedRadioSensors(t *testing.T) {
	t.Parallel()

	provider := newRetainedSemanticState()
	regulatorMapping := uint16(1)
	thermostatMapping := uint16(2)
	regulatorConnected := true
	thermostatConnected := true
	regulatorTemp := 14.8125
	regulatorHumidity := 62.0
	thermostatTemp := 17.6875
	thermostatHumidity := 53.0

	poller := &vaillantSemanticPoller{
		publication: newVaillantSemanticPublication(provider, nil, nil),
		zones: map[byte]*vaillantZoneSnapshot{
			0x00: {
				Instance: 0x00,
				Present:  true,
				Name:     "Parter",
				ConfigurationRoomTemperatureZoneMappingRaw: &regulatorMapping,
			},
			0x01: {
				Instance: 0x01,
				Present:  true,
				Name:     "Etaj",
				ConfigurationRoomTemperatureZoneMappingRaw: &thermostatMapping,
			},
		},
		radioDevices: map[radioDeviceKey]*vaillantRadioDeviceSnapshot{
			{Group: remoteRegulators.group, Instance: 0x01}: {
				Group:            remoteRegulators.group,
				Instance:         0x01,
				DeviceConnected:  &regulatorConnected,
				RoomTemperatureC: &regulatorTemp,
				RoomHumidityPct:  &regulatorHumidity,
			},
			{Group: remoteThermostats.group, Instance: 0x01}: {
				Group:            remoteThermostats.group,
				Instance:         0x01,
				DeviceConnected:  &thermostatConnected,
				RoomTemperatureC: &thermostatTemp,
				RoomHumidityPct:  &thermostatHumidity,
			},
		},
	}

	poller.publishZones(semanticSnapshotSourceLive)

	zones := provider.Zones()
	if len(zones) != 2 {
		t.Fatalf("len(provider.Zones()) = %d; want 2", len(zones))
	}
	if zones[0].State.CurrentTempC == nil || *zones[0].State.CurrentTempC != regulatorTemp {
		t.Fatalf("zone 1 current temp = %#v; want %.4f", zones[0].State.CurrentTempC, regulatorTemp)
	}
	if zones[0].State.CurrentHumidityPct == nil || *zones[0].State.CurrentHumidityPct != regulatorHumidity {
		t.Fatalf("zone 1 humidity = %#v; want %.1f", zones[0].State.CurrentHumidityPct, regulatorHumidity)
	}
	if zones[1].State.CurrentTempC == nil || *zones[1].State.CurrentTempC != thermostatTemp {
		t.Fatalf("zone 2 current temp = %#v; want %.4f", zones[1].State.CurrentTempC, thermostatTemp)
	}
	if zones[1].State.CurrentHumidityPct == nil || *zones[1].State.CurrentHumidityPct != thermostatHumidity {
		t.Fatalf("zone 2 humidity = %#v; want %.1f", zones[1].State.CurrentHumidityPct, thermostatHumidity)
	}
}

func TestVaillantOwnerVaillantSemanticPollerPublishZonesInfersThermostatMappingFromRadioAssignment(t *testing.T) {
	t.Parallel()

	provider := newRetainedSemanticState()
	cache := &semanticSnapshotCaptureSpy{}
	connected := true
	zoneAssignment := uint8(2)
	radioTemp := 17.8125
	radioHumidity := 53.0

	poller := &vaillantSemanticPoller{
		publication: newVaillantSemanticPublication(provider, nil, cache),
		zones: map[byte]*vaillantZoneSnapshot{
			0x01: {
				Instance: 0x01,
				Present:  true,
				Name:     "Etaj",
			},
		},
		radioDevices: map[radioDeviceKey]*vaillantRadioDeviceSnapshot{
			{Group: remoteThermostats.group, Instance: 0x01}: {
				Group:            remoteThermostats.group,
				Instance:         0x01,
				DeviceConnected:  &connected,
				ZoneAssignment:   &zoneAssignment,
				RoomTemperatureC: &radioTemp,
				RoomHumidityPct:  &radioHumidity,
			},
		},
	}

	poller.publishZones(semanticSnapshotSourceLive)

	zones := provider.Zones()
	if len(zones) != 1 {
		t.Fatalf("len(provider.Zones()) = %d; want 1", len(zones))
	}
	if zones[0].State.CurrentTempC == nil || *zones[0].State.CurrentTempC != radioTemp {
		t.Fatalf("zone current temp = %#v; want radio %.4f", zones[0].State.CurrentTempC, radioTemp)
	}
	if zones[0].State.CurrentHumidityPct == nil || *zones[0].State.CurrentHumidityPct != radioHumidity {
		t.Fatalf("zone humidity = %#v; want radio %.1f", zones[0].State.CurrentHumidityPct, radioHumidity)
	}
	if zones[0].Config.RoomTemperatureZoneMapping == nil || *zones[0].Config.RoomTemperatureZoneMapping != 2 {
		t.Fatalf("roomTemperatureZoneMapping = %#v; want inferred 2", zones[0].Config.RoomTemperatureZoneMapping)
	}
	if zones[0].Config.AssociatedCircuit == nil || *zones[0].Config.AssociatedCircuit != 1 {
		t.Fatalf("associatedCircuit = %#v; want inferred 1", zones[0].Config.AssociatedCircuit)
	}
	if cache.calls == 0 || len(cache.last.Zones) != 1 {
		t.Fatalf("cache calls=%d zones=%d; want persisted zone snapshot", cache.calls, len(cache.last.Zones))
	}
	if cache.last.Zones[0].State.CurrentTempC != nil {
		t.Fatalf("cached current temp = %#v; want nil direct zone value", cache.last.Zones[0].State.CurrentTempC)
	}
	if cache.last.Zones[0].Config.RoomTemperatureZoneMapping != nil {
		t.Fatalf("cached inferred mapping = %#v; want nil direct zone mapping", cache.last.Zones[0].Config.RoomTemperatureZoneMapping)
	}
	if cache.last.Zones[0].Config.AssociatedCircuit != nil {
		t.Fatalf("cached inferred associated circuit = %#v; want nil direct associated circuit", cache.last.Zones[0].Config.AssociatedCircuit)
	}
}

func TestVaillantOwnerVaillantSemanticPollerPublishZonesDoesNotInferAmbiguousThermostatAssignment(t *testing.T) {
	t.Parallel()

	provider := newRetainedSemanticState()
	cache := &semanticSnapshotCaptureSpy{}
	connected := true
	zoneAssignment := uint8(2)
	firstTemp := 17.8125
	secondTemp := 18.25

	poller := &vaillantSemanticPoller{
		publication: newVaillantSemanticPublication(provider, nil, cache),
		zones: map[byte]*vaillantZoneSnapshot{
			0x01: {
				Instance: 0x01,
				Present:  true,
				Name:     "Etaj",
			},
		},
		radioDevices: map[radioDeviceKey]*vaillantRadioDeviceSnapshot{
			{Group: remoteThermostats.group, Instance: 0x01}: {
				Group:            remoteThermostats.group,
				Instance:         0x01,
				DeviceConnected:  &connected,
				ZoneAssignment:   &zoneAssignment,
				RoomTemperatureC: &firstTemp,
			},
			{Group: remoteThermostats.group, Instance: 0x02}: {
				Group:            remoteThermostats.group,
				Instance:         0x02,
				DeviceConnected:  &connected,
				ZoneAssignment:   &zoneAssignment,
				RoomTemperatureC: &secondTemp,
			},
		},
	}

	poller.publishZones(semanticSnapshotSourceLive)

	zones := provider.Zones()
	if len(zones) != 1 {
		t.Fatalf("len(provider.Zones()) = %d; want 1", len(zones))
	}
	if zones[0].State.CurrentTempC != nil {
		t.Fatalf("zone current temp = %#v; want nil for ambiguous inferred mapping", zones[0].State.CurrentTempC)
	}
	if zones[0].Config.RoomTemperatureZoneMapping != nil {
		t.Fatalf("roomTemperatureZoneMapping = %#v; want nil for ambiguous inferred mapping", zones[0].Config.RoomTemperatureZoneMapping)
	}
	if zones[0].Config.AssociatedCircuit != nil {
		t.Fatalf("associatedCircuit = %#v; want nil for ambiguous inferred mapping", zones[0].Config.AssociatedCircuit)
	}
	if cache.calls == 0 || len(cache.last.Zones) != 1 {
		t.Fatalf("cache calls=%d zones=%d; want persisted zone snapshot", cache.calls, len(cache.last.Zones))
	}
	if cache.last.Zones[0].Config.RoomTemperatureZoneMapping != nil {
		t.Fatalf("cached ambiguous mapping = %#v; want nil", cache.last.Zones[0].Config.RoomTemperatureZoneMapping)
	}
	if cache.last.Zones[0].Config.AssociatedCircuit != nil {
		t.Fatalf("cached ambiguous associated circuit = %#v; want nil", cache.last.Zones[0].Config.AssociatedCircuit)
	}
}

func TestVaillantOwnerVaillantSemanticPollerPublishZonesKeepsDirectRoomStateOverRadioFallback(t *testing.T) {
	t.Parallel()

	provider := newRetainedSemanticState()
	mapping := uint16(2)
	connected := true
	directTemp := 21.0
	directHumidity := 44.0
	radioTemp := 17.6875
	radioHumidity := 53.0

	poller := &vaillantSemanticPoller{
		publication: newVaillantSemanticPublication(provider, nil, nil),
		zones: map[byte]*vaillantZoneSnapshot{
			0x01: {
				Instance:     0x01,
				Present:      true,
				Name:         "Etaj",
				CurrentTempC: &directTemp,
				HumidityPct:  &directHumidity,
				ConfigurationRoomTemperatureZoneMappingRaw: &mapping,
			},
		},
		radioDevices: map[radioDeviceKey]*vaillantRadioDeviceSnapshot{
			{Group: remoteThermostats.group, Instance: 0x01}: {
				Group:            remoteThermostats.group,
				Instance:         0x01,
				DeviceConnected:  &connected,
				RoomTemperatureC: &radioTemp,
				RoomHumidityPct:  &radioHumidity,
			},
		},
	}

	poller.publishZones(semanticSnapshotSourceLive)

	zones := provider.Zones()
	if len(zones) != 1 {
		t.Fatalf("len(provider.Zones()) = %d; want 1", len(zones))
	}
	if zones[0].State.CurrentTempC == nil || *zones[0].State.CurrentTempC != directTemp {
		t.Fatalf("zone current temp = %#v; want direct %.1f", zones[0].State.CurrentTempC, directTemp)
	}
	if zones[0].State.CurrentHumidityPct == nil || *zones[0].State.CurrentHumidityPct != directHumidity {
		t.Fatalf("zone humidity = %#v; want direct %.1f", zones[0].State.CurrentHumidityPct, directHumidity)
	}
}
