package owner

import (
	"context"
	"fmt"
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/productids"
	"slices"
	"testing"
)

func TestVaillantOwnerVaillantSemanticPollerPublishZonesRadioFallbackDoesNotPersistAsDirectState(t *testing.T) {
	t.Parallel()

	provider := newRetainedSemanticState()
	cache := &semanticSnapshotCaptureSpy{}
	mapping := uint16(2)
	connected := true
	radioTemp := 17.6875
	radioHumidity := 53.0

	poller := &vaillantSemanticPoller{
		publication: newVaillantSemanticPublication(provider, nil, cache),
		zones: map[byte]*vaillantZoneSnapshot{
			0x01: {
				Instance: 0x01,
				Present:  true,
				Name:     "Etaj",
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

	published := provider.Zones()
	if len(published) != 1 {
		t.Fatalf("len(provider.Zones()) = %d; want 1", len(published))
	}
	if published[0].State.CurrentTempC == nil || *published[0].State.CurrentTempC != radioTemp {
		t.Fatalf("published current temp = %#v; want radio fallback %.4f", published[0].State.CurrentTempC, radioTemp)
	}
	if cache.calls == 0 || len(cache.last.Zones) != 1 {
		t.Fatalf("cache calls=%d zones=%d; want persisted zone snapshot", cache.calls, len(cache.last.Zones))
	}
	if cache.last.Zones[0].State.CurrentTempC != nil {
		t.Fatalf("cached current temp = %#v; want nil direct zone value", cache.last.Zones[0].State.CurrentTempC)
	}
	if cache.last.Zones[0].State.CurrentHumidityPct != nil {
		t.Fatalf("cached humidity = %#v; want nil direct zone value", cache.last.Zones[0].State.CurrentHumidityPct)
	}
}

func TestVaillantOwnerVaillantSemanticPollerPublishZonesSkipsUnknownRoomSensorMappingFallback(t *testing.T) {
	t.Parallel()

	for _, raw := range []uint16{0, 0xFF, 0xFFFF} {
		raw := raw
		t.Run(fmt.Sprintf("mapping_0x%04x", raw), func(t *testing.T) {
			t.Parallel()

			provider := newRetainedSemanticState()
			connected := true
			radioTemp := 17.6875
			poller := &vaillantSemanticPoller{
				publication: newVaillantSemanticPublication(provider, nil, nil),
				zones: map[byte]*vaillantZoneSnapshot{
					0x01: {
						Instance: 0x01,
						Present:  true,
						Name:     "Etaj",
						ConfigurationRoomTemperatureZoneMappingRaw: &raw,
					},
				},
				radioDevices: map[radioDeviceKey]*vaillantRadioDeviceSnapshot{
					{Group: remoteThermostats.group, Instance: 0xFE}: {
						Group:            remoteThermostats.group,
						Instance:         0xFE,
						DeviceConnected:  &connected,
						RoomTemperatureC: &radioTemp,
					},
				},
			}

			poller.publishZones(semanticSnapshotSourceLive)

			zones := provider.Zones()
			if len(zones) != 1 {
				t.Fatalf("len(provider.Zones()) = %d; want 1", len(zones))
			}
			if zones[0].State.CurrentTempC != nil {
				t.Fatalf("zone current temp = %#v; want nil for unknown mapping 0x%04x", zones[0].State.CurrentTempC, raw)
			}
		})
	}
}

func TestVaillantOwnerMergeDhwSnapshotFieldsPartialLiveUpdatePreservesLastKnownAndFreshness(t *testing.T) {
	t.Parallel()

	current := 48.0
	target := 50.0
	entry := &vaillantDhwSnapshot{
		OperatingMode:                 "auto",
		Preset:                        "schedule",
		CurrentTempC:                  &current,
		TargetTempC:                   &target,
		ConfigurationDHWOperationMode: "1",
		StateSpecialFunction:          "0",
	}
	seedDhwFreshness(entry, semanticSnapshotSourceCache, true)

	updatedCurrent := 49.5
	incoming := &vaillantDhwSnapshot{
		CurrentTempC: &updatedCurrent,
	}
	mergeDhwSnapshotFields(entry, incoming, semanticSnapshotSourceLive, dhwFieldSet)

	if entry.CurrentTempC == nil || *entry.CurrentTempC != 49.5 {
		t.Fatalf("entry.CurrentTempC = %v; want 49.5", entry.CurrentTempC)
	}
	if entry.TargetTempC == nil || *entry.TargetTempC != 50.0 {
		t.Fatalf("entry.TargetTempC = %v; want preserved 50.0", entry.TargetTempC)
	}
	if entry.OperatingMode != "auto" {
		t.Fatalf("entry.OperatingMode = %q; want preserved auto", entry.OperatingMode)
	}

	currentFreshness, ok := entry.FieldFreshness[dhwFieldCurrentTempC]
	if !ok || currentFreshness.Source != semanticSnapshotSourceLive || currentFreshness.Stale {
		t.Fatalf("current_temp freshness = %+v (ok=%v); want source=live stale=false", currentFreshness, ok)
	}
	targetFreshness, ok := entry.FieldFreshness[dhwFieldTargetTempC]
	if !ok || targetFreshness.Source != semanticSnapshotSourceCache || !targetFreshness.Stale {
		t.Fatalf("target_temp freshness = %+v (ok=%v); want source=cache stale=true", targetFreshness, ok)
	}
}

func TestVaillantOwnerZonePresenceFSMNoRegressionWithFreshnessMetadata(t *testing.T) {
	t.Parallel()

	instance := byte(0x02)
	poller := &vaillantSemanticPoller{
		zones:             make(map[byte]*vaillantZoneSnapshot),
		presence:          make(map[byte]*zonePresenceRecord),
		zoneMissThreshold: 3,
		zoneHitThreshold:  2,
	}

	current := 20.0
	entry := &vaillantZoneSnapshot{Instance: instance, Present: true, CurrentTempC: &current}
	seedZoneFreshness(entry, semanticSnapshotSourceCache, true)
	poller.zones[instance] = entry
	poller.presence[instance] = &zonePresenceRecord{
		State:     zonePresencePresent,
		HitStreak: poller.zoneHitThresholdValue(),
	}

	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{})
	if _, exists := poller.zones[instance]; !exists {
		t.Fatalf("zone removed too early on first miss")
	}
	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{})
	if _, exists := poller.zones[instance]; !exists {
		t.Fatalf("zone removed too early on second miss")
	}
	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{})
	if _, exists := poller.zones[instance]; exists {
		t.Fatalf("zone should be absent after third miss")
	}
}

func TestVaillantOwnerZonePresenceFSMAbsenceRequiresConsecutiveMisses(t *testing.T) {
	t.Parallel()

	instance := byte(0x02)
	poller := &vaillantSemanticPoller{
		zones:             make(map[byte]*vaillantZoneSnapshot),
		presence:          make(map[byte]*zonePresenceRecord),
		zoneMissThreshold: 3,
		zoneHitThreshold:  2,
	}

	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{instance: true})
	if _, exists := poller.zones[instance]; exists {
		t.Fatalf("zone should not be present after first hit when hit threshold is 2")
	}
	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{instance: true})
	if _, exists := poller.zones[instance]; !exists {
		t.Fatalf("zone should be present after second consecutive hit")
	}

	for miss := 1; miss <= 2; miss++ {
		poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{})
		if _, exists := poller.zones[instance]; !exists {
			t.Fatalf("zone removed too early after miss %d", miss)
		}
	}

	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{})
	if _, exists := poller.zones[instance]; exists {
		t.Fatalf("zone should be absent after third consecutive miss")
	}
	record := poller.presence[instance]
	if record == nil || record.State != zonePresenceAbsent {
		t.Fatalf("presence state = %+v; want ABSENT", record)
	}
}

func TestVaillantOwnerZonePresenceFSMResurrectionRequiresConsecutiveHits(t *testing.T) {
	t.Parallel()

	instance := byte(0x01)
	poller := &vaillantSemanticPoller{
		zones:             make(map[byte]*vaillantZoneSnapshot),
		presence:          make(map[byte]*zonePresenceRecord),
		zoneMissThreshold: 3,
		zoneHitThreshold:  2,
	}

	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{instance: true})
	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{instance: true})
	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{})
	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{})
	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{})
	if _, exists := poller.zones[instance]; exists {
		t.Fatalf("zone should be absent after miss threshold")
	}

	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{instance: true})
	if _, exists := poller.zones[instance]; exists {
		t.Fatalf("zone should stay absent until hit threshold is reached")
	}

	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{})
	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{instance: true})
	if _, exists := poller.zones[instance]; exists {
		t.Fatalf("single hit after miss should not resurrect zone")
	}

	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{instance: true})
	if _, exists := poller.zones[instance]; !exists {
		t.Fatalf("zone should resurrect after two consecutive hits")
	}
	record := poller.presence[instance]
	if record == nil || record.State != zonePresencePresent {
		t.Fatalf("presence state = %+v; want PRESENT", record)
	}
}

func TestVaillantOwnerZonePresenceFSMAlternatingProbesDoNotFlapPresentZone(t *testing.T) {
	t.Parallel()

	instance := byte(0x00)
	poller := &vaillantSemanticPoller{
		zones:             make(map[byte]*vaillantZoneSnapshot),
		presence:          make(map[byte]*zonePresenceRecord),
		zoneMissThreshold: 3,
		zoneHitThreshold:  2,
	}

	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{instance: true})
	poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{instance: true})

	for iteration := 0; iteration < 6; iteration++ {
		poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{})
		if _, exists := poller.zones[instance]; !exists {
			t.Fatalf("zone disappeared on alternating miss at iteration %d", iteration)
		}
		record := poller.presence[instance]
		if record == nil || record.State != zonePresenceSuspectMissing {
			t.Fatalf("presence state after miss = %+v; want SUSPECT_MISSING", record)
		}

		poller.applyZonePresenceProbes(map[byte]bool{instance: true}, map[byte]bool{instance: true})
		if _, exists := poller.zones[instance]; !exists {
			t.Fatalf("zone disappeared on alternating hit at iteration %d", iteration)
		}
		record = poller.presence[instance]
		if record == nil || record.State != zonePresencePresent {
			t.Fatalf("presence state after hit = %+v; want PRESENT", record)
		}
	}
}

func TestVaillantOwnerReconcileDiscoveryPresenceSkipsMissDemotionWhenFallbackHydrates(t *testing.T) {
	t.Parallel()

	instance := byte(0x03)
	poller := &vaillantSemanticPoller{
		zones:             make(map[byte]*vaillantZoneSnapshot),
		presence:          make(map[byte]*zonePresenceRecord),
		zoneMissThreshold: 3,
		zoneHitThreshold:  2,
		transportConfig:   engine.TransportConfig{Protocol: engine.TransportEbusdTCP},
	}

	poller.presence[instance] = &zonePresenceRecord{
		State:     zonePresenceSuspectResurrect,
		HitStreak: 1,
	}

	poller.refreshFromEbusdGrabFn = func(_ context.Context) (map[byte]bool, bool) {
		poller.applyZonePresenceProbes(
			map[byte]bool{instance: true},
			map[byte]bool{instance: true},
		)
		return map[byte]bool{instance: true}, true
	}

	source := poller.reconcileDiscoveryPresence(
		context.Background(),
		map[byte]bool{instance: true},
		map[byte]bool{},
	)
	if source != semanticSnapshotSourceLive {
		t.Fatalf("source = %v; want LIVE", source)
	}
	record := poller.presence[instance]
	if record == nil || record.State != zonePresencePresent {
		t.Fatalf("presence state = %+v; want PRESENT", record)
	}
	if _, exists := poller.zones[instance]; !exists {
		t.Fatalf("zone should be present after fallback hydration hit reaches threshold")
	}
}

func TestVaillantOwnerReconcileDiscoveryPresenceAppliesMissesWhenFallbackUnavailable(t *testing.T) {
	t.Parallel()

	instance := byte(0x01)
	poller := &vaillantSemanticPoller{
		zones:             make(map[byte]*vaillantZoneSnapshot),
		presence:          make(map[byte]*zonePresenceRecord),
		zoneMissThreshold: 3,
		zoneHitThreshold:  2,
	}
	poller.zones[instance] = &vaillantZoneSnapshot{Instance: instance, Present: true}
	poller.presence[instance] = &zonePresenceRecord{
		State:     zonePresencePresent,
		HitStreak: poller.zoneHitThresholdValue(),
	}
	poller.refreshFromEbusdGrabFn = func(_ context.Context) (map[byte]bool, bool) { return nil, false }

	source := poller.reconcileDiscoveryPresence(
		context.Background(),
		map[byte]bool{instance: true},
		map[byte]bool{},
	)
	if source != semanticSnapshotSourceCache {
		t.Fatalf("source = %v; want CACHE", source)
	}
	record := poller.presence[instance]
	if record == nil || record.State != zonePresenceSuspectMissing {
		t.Fatalf("presence state = %+v; want SUSPECT_MISSING", record)
	}
	if _, exists := poller.zones[instance]; !exists {
		t.Fatalf("zone should remain present before miss threshold")
	}
}

func TestVaillantOwnerReconcileDiscoveryPresenceTracksMissesWithPartialFallback(t *testing.T) {
	t.Parallel()

	keep := byte(0x00)
	missing := byte(0x01)
	poller := &vaillantSemanticPoller{
		zones:             make(map[byte]*vaillantZoneSnapshot),
		presence:          make(map[byte]*zonePresenceRecord),
		zoneMissThreshold: 3,
		zoneHitThreshold:  2,
		transportConfig:   engine.TransportConfig{Protocol: engine.TransportEbusdTCP},
	}
	poller.zones[keep] = &vaillantZoneSnapshot{Instance: keep, Present: true}
	poller.zones[missing] = &vaillantZoneSnapshot{Instance: missing, Present: true}
	poller.presence[keep] = &zonePresenceRecord{State: zonePresencePresent, HitStreak: poller.zoneHitThresholdValue()}
	poller.presence[missing] = &zonePresenceRecord{State: zonePresencePresent, HitStreak: poller.zoneHitThresholdValue()}
	poller.refreshFromEbusdGrabFn = func(_ context.Context) (map[byte]bool, bool) {
		return map[byte]bool{keep: true}, true
	}

	source := poller.reconcileDiscoveryPresence(
		context.Background(),
		map[byte]bool{keep: true, missing: true},
		map[byte]bool{},
	)
	if source != semanticSnapshotSourceLive {
		t.Fatalf("source = %v; want LIVE", source)
	}
	if _, exists := poller.zones[keep]; !exists {
		t.Fatalf("fallback-present zone must remain present")
	}
	if _, exists := poller.zones[missing]; !exists {
		t.Fatalf("missing zone must not be removed on first miss")
	}
	record := poller.presence[missing]
	if record == nil || record.State != zonePresenceSuspectMissing {
		t.Fatalf("missing zone state = %+v; want SUSPECT_MISSING", record)
	}
}

func TestVaillantOwnerReconcileDiscoveryPresenceDoesNotDoubleCountFallbackHit(t *testing.T) {
	t.Parallel()

	instance := byte(0x04)
	poller := &vaillantSemanticPoller{
		zones:             make(map[byte]*vaillantZoneSnapshot),
		presence:          make(map[byte]*zonePresenceRecord),
		zoneMissThreshold: 3,
		zoneHitThreshold:  2,
		transportConfig:   engine.TransportConfig{Protocol: engine.TransportEbusdTCP},
	}
	poller.refreshFromEbusdGrabFn = func(_ context.Context) (map[byte]bool, bool) {
		poller.applyZonePresenceProbes(
			map[byte]bool{instance: true},
			map[byte]bool{instance: true},
		)
		return map[byte]bool{instance: true}, true
	}

	source := poller.reconcileDiscoveryPresence(
		context.Background(),
		map[byte]bool{instance: true},
		map[byte]bool{},
	)
	if source != semanticSnapshotSourceLive {
		t.Fatalf("source = %v; want LIVE", source)
	}
	if _, exists := poller.zones[instance]; exists {
		t.Fatalf("zone should not be promoted to present after a single fallback hit when threshold is 2")
	}
	record := poller.presence[instance]
	if record == nil || record.State != zonePresenceSuspectResurrect || record.HitStreak != 1 {
		t.Fatalf("presence state = %+v; want SUSPECT_RESURRECT with hit_streak=1", record)
	}
}

func TestVaillantOwnerSemanticReadBreakerKeyIncludesOpcode(t *testing.T) {
	t.Parallel()

	localKey := SemanticReadBreakerKey(0x15, vaillantB524OpcodeLocal, localDHW.group, dhwInstance, dhwCurrentTemp)
	readKey := SemanticReadBreakerKey(0x15, vaillantB524OpcodeRead, localDHW.group, dhwInstance, dhwCurrentTemp)
	if localKey == readKey {
		t.Fatalf("semanticReadBreakerKey must include opcode; got equal keys %q", localKey)
	}

	otherTarget := SemanticReadBreakerKey(0x16, vaillantB524OpcodeRead, localDHW.group, dhwInstance, dhwCurrentTemp)
	if readKey == otherTarget {
		t.Fatalf("semanticReadBreakerKey must include target; got equal keys %q", readKey)
	}
}

func TestVaillantOwnerParseB524ZonesFromGrabFiltersAbsentInstances(t *testing.T) {
	t.Parallel()

	lines := []string{
		"f115b52406020003001c00 / 0501031c0000 = 1: basv Z1ZoneCircuitIndex",
		"3115b52406020003001600 / 0b0303160050617274657200 = 2: basv Z1Shortname",
		"f115b52406020003000f00 / 0801030f0000007041 = 1: basv Z1RoomTemp",
		"f115b52406020003012200 / 080303220000002041 = 1: basv Z2HeatingManualTemp",
		"f115b52406020003011c00 / 0501031c0001 = 2: basv Z2ZoneCircuitIndex",
		"3115b52406020003011600 / 09030316004574616a00 = 1: basv Z2Shortname",
		"f115b52406020003021c00 / 0500031c00ff = 2: basv Z3ZoneCircuitIndex",
		"3115b52406020003021600 / 0b020316005a6f6e65203300 = 1: basv Z3Shortname",
	}

	zones := parseB524ZonesFromGrab(lines, 0x15)
	if len(zones) != 2 {
		t.Fatalf("len(zones) = %d; want 2", len(zones))
	}

	zone1, ok := zones[0]
	if !ok {
		t.Fatalf("zone instance 0 missing")
	}
	if zone1.Name != "Parter" {
		t.Fatalf("zone1 name = %q; want %q", zone1.Name, "Parter")
	}
	if zone1.CurrentTempC == nil || *zone1.CurrentTempC != 15.0 {
		t.Fatalf("zone1 current = %v; want 15.0", zone1.CurrentTempC)
	}

	zone2, ok := zones[1]
	if !ok {
		t.Fatalf("zone instance 1 missing")
	}
	if zone2.Name != "Etaj" {
		t.Fatalf("zone2 name = %q; want %q", zone2.Name, "Etaj")
	}
	if zone2.TargetTempC == nil || *zone2.TargetTempC != 10.0 {
		t.Fatalf("zone2 target = %v; want 10.0", zone2.TargetTempC)
	}

	if _, exists := zones[2]; exists {
		t.Fatalf("zone instance 2 should be filtered as absent")
	}
}

func TestVaillantOwnerParseB524ZonesFromGrabIncludesHumidityAndAllowedModes(t *testing.T) {
	t.Parallel()

	lines := []string{
		"f115b52406020003001c00 / 0501031c0000 = 1",
		"f115b52406020003000600 / 010306000200 = 1",
		"f115b52406020003000e00 / 01030e000000 = 1",
		"f115b52406020003001300 / 010313000000 = 1",
		"f115b52406020003001200 / 010312000100 = 1",
		"f115b52406020003002800 / 0103280000002042 = 1",
		"f115b52406020003002200 / 0103220000003041 = 1",
		"f115b52406020002000200 / 010202000100 = 1",
	}

	zones := parseB524ZonesFromGrab(lines, 0x15)
	zone, ok := zones[0]
	if !ok {
		t.Fatalf("zone instance 0 missing")
	}
	if zone.HumidityPct == nil || *zone.HumidityPct != 40.0 {
		t.Fatalf("zone humidity = %v; want 40.0", zone.HumidityPct)
	}
	if zone.Preset != "manual" {
		t.Fatalf("zone preset = %q; want manual", zone.Preset)
	}
	if got, want := zone.AllowedModes, []string{"off", "auto", "heat"}; !slices.Equal(got, want) {
		t.Fatalf("allowed modes = %v; want %v", got, want)
	}
}

func TestVaillantOwnerParseB524DhwFromGrab(t *testing.T) {
	t.Parallel()

	lines := []string{
		"f115b52406020001000300 / 010103000100 = 1",
		"f115b52406020001000400 / 0101040000005c42 = 1",
		"f115b52406020001000500 / 0101050000004842 = 1",
		"f115b52406020001000d00 / 01010d000000 = 1",
	}

	dhw, ok := parseB524DhwFromGrab(lines, 0x15)
	if !ok || dhw == nil {
		t.Fatalf("parseB524DhwFromGrab returned nil")
	}
	if dhw.OperatingMode != "auto" {
		t.Fatalf("operating mode = %q; want auto", dhw.OperatingMode)
	}
	if dhw.Preset != "schedule" {
		t.Fatalf("preset = %q; want schedule", dhw.Preset)
	}
	if dhw.TargetTempC == nil || *dhw.TargetTempC != 55.0 {
		t.Fatalf("target temp = %v; want 55", dhw.TargetTempC)
	}
	if dhw.CurrentTempC == nil || *dhw.CurrentTempC != 50.0 {
		t.Fatalf("current temp = %v; want 50", dhw.CurrentTempC)
	}
}

func TestVaillantOwnerParseB524DhwFromGrabAcceptsRemoteReadOpcode(t *testing.T) {
	t.Parallel()

	lines := []string{
		"f115b52406060001000300 / 010103000100 = 1",
		"f115b52406060001000400 / 0101040000005c42 = 1",
		"f115b52406060001000500 / 0101050000004842 = 1",
		"f115b52406060001000d00 / 01010d000000 = 1",
	}

	dhw, ok := parseB524DhwFromGrab(lines, 0x15)
	if !ok || dhw == nil {
		t.Fatalf("parseB524DhwFromGrab returned nil")
	}
	if dhw.OperatingMode != "auto" {
		t.Fatalf("operating mode = %q; want auto", dhw.OperatingMode)
	}
	if dhw.CurrentTempC == nil || *dhw.CurrentTempC != 50.0 {
		t.Fatalf("current temp = %v; want 50", dhw.CurrentTempC)
	}
}

func TestVaillantOwnerParseB524GrabLineRejectsUnsupportedOpcode(t *testing.T) {
	t.Parallel()

	line := "f115b52406010001000500 / 0801050000004842 = 1"
	_, _, _, _, ok := parseB524GrabLine(line, 0x15)
	if ok {
		t.Fatalf("expected unsupported opcode line to be rejected")
	}
}

func TestVaillantOwnerVaillantSemanticPollerRegulatorDetectionPresent(t *testing.T) {
	t.Parallel()

	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error: %v", err)
	}

	// 0020028521 = Vaillant calorMATIC 430f VRC 430f (Regulator)
	reg := newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV", SerialNumber: "21-22-09-0010002460-0082-005409-N4"},
		registry.DeviceInfo{Address: 0x60, Manufacturer: "Vaillant", DeviceID: "VRC430", SerialNumber: "21-22-09-0020028521-0082-005409-N4"},
	)

	poller := &vaillantSemanticPoller{
		reg:     reg,
		catalog: catalog,
	}

	got := poller.findRegulatorCapability()
	if got != productids.ControllerPresent {
		t.Fatalf("findRegulatorCapability() = %s; want ControllerPresent", got)
	}
}

func TestVaillantOwnerVaillantSemanticPollerRegulatorDetectionNone(t *testing.T) {
	t.Parallel()

	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error: %v", err)
	}

	// 0010002315 = Vaillant atmoCOMPACT (Boiler) — not a regulator
	reg := newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV", SerialNumber: "21-22-09-0010002460-0082-005409-N4"},
	)

	poller := &vaillantSemanticPoller{
		reg:     reg,
		catalog: catalog,
	}

	got := poller.findRegulatorCapability()
	if got != productids.ControllerNone {
		t.Fatalf("findRegulatorCapability() = %s; want ControllerNone", got)
	}
}

// Pins the classification delta from the ebusreg 2342a8d8 repin: capability is
// read from Catalog.All and the "Thermostat" role now counts as a present
// controller (before: Regulator only, from the metadata-complete ByPartNumber
// subset). Affects the published regulator capability and absence FSM.
func TestVaillantOwnerVaillantSemanticPollerRegulatorDetectionEbusregRepinDelta(t *testing.T) {
	t.Parallel()

	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error: %v", err)
	}
	cases := []struct {
		name string
		part string
		want productids.ControllerCapability
	}{
		{"thermostat role is a present controller", "0020018249", productids.ControllerPresent},
		{"boiler outside ByPartNumber is none, not unknown", "0010002465", productids.ControllerNone},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reg := newTestRegistry(
				registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "X", SerialNumber: "21-22-09-" + tc.part + "-0082-005409-N4"},
			)
			poller := &vaillantSemanticPoller{reg: reg, catalog: catalog}
			if got := poller.findRegulatorCapability(); got != tc.want {
				t.Fatalf("findRegulatorCapability(%s) = %s; want %s", tc.part, got, tc.want)
			}
		})
	}
}

func TestVaillantOwnerVaillantSemanticPollerRegulatorDetectionUnknown(t *testing.T) {
	t.Parallel()

	catalog, err := productids.LoadCatalog()
	if err != nil {
		t.Fatalf("LoadCatalog() error: %v", err)
	}

	// Device with no serial number — cannot extract part number, so unknown.
	reg := newTestRegistry(
		registry.DeviceInfo{Address: 0x15, Manufacturer: "Vaillant", DeviceID: "BASV"},
	)

	poller := &vaillantSemanticPoller{
		reg:     reg,
		catalog: catalog,
	}

	got := poller.findRegulatorCapability()
	if got != productids.ControllerUnknown {
		t.Fatalf("findRegulatorCapability() = %s; want ControllerUnknown", got)
	}
}
