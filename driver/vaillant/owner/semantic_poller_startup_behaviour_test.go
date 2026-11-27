package owner

import (
	"context"
	"fmt"
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"testing"
	"time"
)

func TestVaillantOwnerDecodeRadioDeviceModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		class *uint8
		want  string
	}{
		{name: "nil class", class: nil, want: ""},
		{name: "vrc720", class: uint8Ptr(0x15), want: "VRC720"},
		{name: "vr71 fm5", class: uint8Ptr(0x26), want: "VR71/FM5"},
		{name: "vr92", class: uint8Ptr(0x35), want: "VR92"},
		{name: "unknown class", class: uint8Ptr(0x44), want: "Unknown (0x44)"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := decodeRadioDeviceModel(test.class); got != test.want {
				t.Fatalf("decodeRadioDeviceModel(%v) = %q; want %q", test.class, got, test.want)
			}
		})
	}
}

func TestVaillantOwnerHasRemoteIdentityEvidence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		classAddress *uint8
		firmware     *string
		hardware     *uint16
		want         bool
	}{
		{name: "vr71 class implies identity", classAddress: uint8Ptr(0x26), want: true},
		{name: "firmware implies identity", firmware: stringPtr("01.00.00"), want: true},
		{name: "hardware implies identity", hardware: uint16Ptr(0x5904), want: true},
		{name: "zero hardware no identity", hardware: uint16Ptr(0), want: false},
		{name: "ffff hardware no identity", hardware: uint16Ptr(0xFFFF), want: false},
		{name: "no evidence", want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := hasRemoteIdentityEvidence(test.classAddress, test.firmware, test.hardware); got != test.want {
				t.Fatalf("hasRemoteIdentityEvidence(...) = %v; want %v", got, test.want)
			}
		})
	}
}

func TestVaillantOwnerMatchesB524ReplyInstance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		reply     byte
		requested byte
		want      bool
	}{
		{name: "exact match", reply: 0x02, requested: 0x02, want: true},
		{name: "one based reply", reply: 0x03, requested: 0x02, want: true},
		{name: "different reply", reply: 0x07, requested: 0x02, want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := matchesB524ReplyInstance(test.reply, test.requested); got != test.want {
				t.Fatalf("matchesB524ReplyInstance(reply=0x%02x, requested=0x%02x) = %v; want %v", test.reply, test.requested, got, test.want)
			}
		})
	}
}

func TestVaillantOwnerSourceFromEbusdGrab(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		protocol engine.TransportProtocol
		ok       bool
		want     semanticSnapshotSource
	}{
		{
			name:     "ebusd-tcp successful grab is live",
			protocol: engine.TransportEbusdTCP,
			ok:       true,
			want:     semanticSnapshotSourceLive,
		},
		{
			name:     "non-ebusd successful grab stays cache",
			protocol: engine.TransportENH,
			ok:       true,
			want:     semanticSnapshotSourceCache,
		},
		{
			name:     "grab failure stays cache",
			protocol: engine.TransportEbusdTCP,
			ok:       false,
			want:     semanticSnapshotSourceCache,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			poller := &vaillantSemanticPoller{
				transportConfig: engine.TransportConfig{Protocol: test.protocol},
			}
			if got := poller.sourceFromEbusdGrab(test.ok); got != test.want {
				t.Fatalf("sourceFromEbusdGrab(protocol=%q, ok=%v) = %v; want %v", test.protocol, test.ok, got, test.want)
			}
		})
	}
}

func TestVaillantOwnerNewVaillantSemanticPollerBoilerTierCadence(t *testing.T) {
	t.Parallel()

	poller := newVaillantSemanticPollerFromDonor(
		Config{},
		&donorGateway{},
		newRetainedSemanticState(),
		nil,
		nil,
	)

	schedules := poller.boilerStatusTierSchedules()
	if len(schedules) != 3 {
		t.Fatalf("len(boilerStatusTierSchedules) = %d; want 3", len(schedules))
	}

	byTier := make(map[boilerStatusTier]boilerStatusTierSchedule, len(schedules))
	for _, schedule := range schedules {
		byTier[schedule.tier] = schedule
	}

	tests := []struct {
		tier     boilerStatusTier
		interval time.Duration
		priority semanticTaskPriority
	}{
		{tier: boilerStatusTierFast, interval: 30 * time.Second, priority: semanticTaskPriorityHigh},
		{tier: boilerStatusTierMedium, interval: 5 * time.Minute, priority: semanticTaskPriorityMedium},
		{tier: boilerStatusTierSlow, interval: 10 * time.Minute, priority: semanticTaskPriorityLow},
	}
	for _, test := range tests {
		schedule, ok := byTier[test.tier]
		if !ok {
			t.Fatalf("missing tier schedule: %v", test.tier)
		}
		if schedule.interval != test.interval {
			t.Fatalf("tier %v interval = %s; want %s", test.tier, schedule.interval, test.interval)
		}
		if schedule.priority != test.priority {
			t.Fatalf("tier %v priority = %v; want %v", test.tier, schedule.priority, test.priority)
		}
	}
}

func TestVaillantOwnerNewVaillantSemanticPollerNormalizesObserveFirstFlagsFromScalarOnlyConfig(t *testing.T) {
	t.Parallel()

	cfg := Config{
		ObserveFirstEnabled:      true,
		PassiveStateDirectApply:  true,
		PassiveConfigDirectApply: true,
		ExternalWritePolicy:      engine.ObserveFirstExternalWritePolicyRecordOnly,
	}
	poller := newVaillantSemanticPollerFromDonor(
		cfg,
		&donorGateway{},
		newRetainedSemanticState(),
		nil,
		nil,
	)

	want := engine.NormalizeObserveFirstFeatureFlags(
		cfg.ObserveFirstEnabled,
		cfg.PassiveStateDirectApply,
		cfg.PassiveConfigDirectApply,
		cfg.ExternalWritePolicy,
	).State()
	got := poller.shadow.FeatureFlags().State()
	if got.ObserveFirstEnabled != want.ObserveFirstEnabled ||
		got.PassiveStateDirectApply != want.PassiveStateDirectApply ||
		got.PassiveConfigDirectApply != want.PassiveConfigDirectApply ||
		got.ExternalWritePolicy != want.ExternalWritePolicy {
		t.Fatalf("shadow feature flags = %+v; want %+v", got, want)
	}
}

func TestVaillantOwnerNewVaillantSemanticPollerAttachesRuntimeShadowCache(t *testing.T) {
	t.Parallel()

	poller := newVaillantSemanticPollerFromDonor(
		Config{},
		&donorGateway{},
		newRetainedSemanticState(),
		nil,
		nil,
	)
	if poller.shadow == nil {
		t.Fatal("poller.shadow = nil; want runtime shadow cache")
	}

	key := engine.NewB509WatchKey(0x08, 0x0200)
	maxAge := poller.prepareSemanticReadWatch(key)
	if maxAge != 30*time.Second {
		t.Fatalf("prepareSemanticReadWatch(default) = %s; want 30s (descriptor policy, not legacy 500ms)", maxAge)
	}

	result := poller.shadow.Write(engine.ShadowWrite{
		Key:        key,
		Source:     engine.ShadowWriteSourcePassive,
		Confidence: engine.ShadowConfidenceHigh,
		Value:      []byte{0x7A},
		ObservedAt: time.Now().Add(-time.Second),
	})
	if !result.Accepted {
		t.Fatalf("shadow write rejected: %s", result.Reason)
	}

	fetchCalls := 0
	value, err := poller.scheduler.GetWatch(context.Background(), key, maxAge, func(context.Context) ([]byte, error) {
		fetchCalls++
		return []byte{0x55}, nil
	})
	if err != nil {
		t.Fatalf("scheduler.GetWatch() error = %v", err)
	}
	if len(value) != 1 || value[0] != 0x7A {
		t.Fatalf("scheduler.GetWatch() value = %v; want shadow value [0x7a]", value)
	}
	if fetchCalls != 0 {
		t.Fatalf("fetch calls = %d; want 0 when scheduler shadow is wired", fetchCalls)
	}
}

func TestVaillantOwnerPrepareSemanticReadWatchUsesObserverDescriptorFreshness(t *testing.T) {
	t.Parallel()

	key := engine.NewB524WatchKey(0x15, vaillantB524OpcodeLocal, localZones.group, 0x00, zoneCurrentTemp)
	observer := staticSemanticReadWatchObserver{
		observation: engine.WatchObservation{
			State:         engine.WatchObservationStateActive,
			HasDescriptor: true,
			Descriptor: engine.WatchDescriptor{
				Key:               key,
				SemanticClass:     engine.WatchSemanticClassState,
				FreshnessProfile:  engine.WatchFreshnessProfileStateSlow,
				DecoderID:         "test.semantic.read.observer",
				CorrelationPolicy: engine.WatchCorrelationPolicyRequestResponse,
				DirectApplyPolicy: engine.WatchDirectApplyPolicyStateDefault,
			},
			Sources: []engine.WatchActivationSource{engine.WatchActivationSourceTooling},
		},
	}

	poller := newVaillantSemanticPollerFromDonor(
		Config{WatchObserver: observer},
		&donorGateway{},
		newRetainedSemanticState(),
		nil,
		nil,
	)
	maxAge := poller.prepareSemanticReadWatch(key)
	if maxAge != 120*time.Second {
		t.Fatalf("prepareSemanticReadWatch(observed state_slow profile) = %s; want 120s (descriptor policy, not legacy 500ms)", maxAge)
	}

	result := poller.shadow.Write(engine.ShadowWrite{
		Key:        key,
		Source:     engine.ShadowWriteSourcePassive,
		Confidence: engine.ShadowConfidenceHigh,
		Value:      []byte{0x01},
		ObservedAt: time.Now(),
	})
	if !result.Accepted {
		t.Fatalf("shadow write rejected after runtime descriptor bootstrap: %s", result.Reason)
	}
}

func TestVaillantOwnerSemanticReadWatchDescriptorWithDecoderIDClassifiesB524Cadence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		key         engine.WatchKey
		wantClass   engine.WatchSemanticClass
		wantProfile engine.WatchFreshnessProfile
		wantPolicy  engine.WatchDirectApplyPolicy
		wantTTL     time.Duration
	}{
		{
			name:        "zone current temp stays fast state",
			key:         engine.NewB524WatchKey(0x15, vaillantB524OpcodeLocal, localZones.group, 0x00, zoneCurrentTemp),
			wantClass:   engine.WatchSemanticClassState,
			wantProfile: engine.WatchFreshnessProfileStateFast,
			wantPolicy:  engine.WatchDirectApplyPolicyStateDefault,
			wantTTL:     30 * time.Second,
		},
		{
			name:        "zone humidity is slow state",
			key:         engine.NewB524WatchKey(0x15, vaillantB524OpcodeLocal, localZones.group, 0x00, zoneCurrentHumidity),
			wantClass:   engine.WatchSemanticClassState,
			wantProfile: engine.WatchFreshnessProfileStateSlow,
			wantPolicy:  engine.WatchDirectApplyPolicyStateDefault,
			wantTTL:     120 * time.Second,
		},
		{
			name:        "zone target temp is config opt in",
			key:         engine.NewB524WatchKey(0x15, vaillantB524OpcodeLocal, localZones.group, 0x00, zoneTargetTemp),
			wantClass:   engine.WatchSemanticClassConfig,
			wantProfile: engine.WatchFreshnessProfileConfig,
			wantPolicy:  engine.WatchDirectApplyPolicyConfigOptIn,
			wantTTL:     5 * time.Minute,
		},
		{
			name:        "zone index is discovery",
			key:         engine.NewB524WatchKey(0x15, vaillantB524OpcodeLocal, localZones.group, 0x00, zoneIndex),
			wantClass:   engine.WatchSemanticClassDiscovery,
			wantProfile: engine.WatchFreshnessProfileDiscovery,
			wantPolicy:  engine.WatchDirectApplyPolicyNever,
			wantTTL:     time.Hour,
		},
		{
			name:        "regulator flow temperature stays fast state",
			key:         engine.NewB524WatchKey(0x15, vaillantB524OpcodeLocal, localRegulator.group, regulatorInstance, systemFlowTemperature),
			wantClass:   engine.WatchSemanticClassState,
			wantProfile: engine.WatchFreshnessProfileStateFast,
			wantPolicy:  engine.WatchDirectApplyPolicyStateDefault,
			wantTTL:     30 * time.Second,
		},
		{
			name:        "regulator energy totals are slow state",
			key:         engine.NewB524WatchKey(0x15, vaillantB524OpcodeLocal, localRegulator.group, regulatorInstance, energyFuelSumHc),
			wantClass:   engine.WatchSemanticClassState,
			wantProfile: engine.WatchFreshnessProfileStateSlow,
			wantPolicy:  engine.WatchDirectApplyPolicyStateDefault,
			wantTTL:     120 * time.Second,
		},
		{
			name:        "regulator installer field is config",
			key:         engine.NewB524WatchKey(0x15, vaillantB524OpcodeLocal, localRegulator.group, regulatorInstance, systemInstallerMenuCode),
			wantClass:   engine.WatchSemanticClassConfig,
			wantProfile: engine.WatchFreshnessProfileConfig,
			wantPolicy:  engine.WatchDirectApplyPolicyConfigOptIn,
			wantTTL:     5 * time.Minute,
		},
		{
			name:        "regulator module configuration is discovery",
			key:         engine.NewB524WatchKey(0x15, vaillantB524OpcodeLocal, localRegulator.group, regulatorInstance, systemModuleConfigurationFunctionModule),
			wantClass:   engine.WatchSemanticClassDiscovery,
			wantProfile: engine.WatchFreshnessProfileDiscovery,
			wantPolicy:  engine.WatchDirectApplyPolicyNever,
			wantTTL:     time.Hour,
		},
		{
			name:        "circuit type is config",
			key:         engine.NewB524WatchKey(0x15, vaillantB524OpcodeLocal, localCircuits.group, 0x00, circuitTypeRegister),
			wantClass:   engine.WatchSemanticClassConfig,
			wantProfile: engine.WatchFreshnessProfileConfig,
			wantPolicy:  engine.WatchDirectApplyPolicyConfigOptIn,
			wantTTL:     5 * time.Minute,
		},
		{
			name:        "remote slot connected is slow state liveness",
			key:         engine.NewB524WatchKey(0x15, vaillantB524OpcodeRead, remoteFunctionalModules.group, 0x01, deviceSlotConnected),
			wantClass:   engine.WatchSemanticClassState,
			wantProfile: engine.WatchFreshnessProfileStateSlow,
			wantPolicy:  engine.WatchDirectApplyPolicyStateDefault,
			wantTTL:     120 * time.Second,
		},
		{
			name:        "remote slot humidity is slow state",
			key:         engine.NewB524WatchKey(0x15, vaillantB524OpcodeRead, remoteFunctionalModules.group, 0x01, deviceSlotRoomHumidity),
			wantClass:   engine.WatchSemanticClassState,
			wantProfile: engine.WatchFreshnessProfileStateSlow,
			wantPolicy:  engine.WatchDirectApplyPolicyStateDefault,
			wantTTL:     120 * time.Second,
		},
		{
			name:        "remote slot hardware identity is discovery",
			key:         engine.NewB524WatchKey(0x15, vaillantB524OpcodeRead, remoteFunctionalModules.group, 0x01, deviceSlotHardwareIdentifier),
			wantClass:   engine.WatchSemanticClassDiscovery,
			wantProfile: engine.WatchFreshnessProfileDiscovery,
			wantPolicy:  engine.WatchDirectApplyPolicyNever,
			wantTTL:     time.Hour,
		},
		{
			name:        "non B524 keeps prior fast state fallback",
			key:         engine.NewB509WatchKey(0x08, 0x0200),
			wantClass:   engine.WatchSemanticClassState,
			wantProfile: engine.WatchFreshnessProfileStateFast,
			wantPolicy:  engine.WatchDirectApplyPolicyStateDefault,
			wantTTL:     30 * time.Second,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			descriptor := semanticReadWatchDescriptorWithDecoderID(test.key, "test.decoder")
			if descriptor.SemanticClass != test.wantClass {
				t.Fatalf("SemanticClass = %q; want %q", descriptor.SemanticClass, test.wantClass)
			}
			if descriptor.FreshnessProfile != test.wantProfile {
				t.Fatalf("FreshnessProfile = %q; want %q", descriptor.FreshnessProfile, test.wantProfile)
			}
			if descriptor.DirectApplyPolicy != test.wantPolicy {
				t.Fatalf("DirectApplyPolicy = %q; want %q", descriptor.DirectApplyPolicy, test.wantPolicy)
			}
			ttl, err := descriptor.EffectiveFreshnessTTL(EngineCollaborator{})
			if err != nil {
				t.Fatalf("EffectiveFreshnessTTL() error = %v", err)
			}
			if ttl != test.wantTTL {
				t.Fatalf("EffectiveFreshnessTTL() = %s; want %s", ttl, test.wantTTL)
			}
		})
	}
}

func TestVaillantOwnerPrepareSemanticReadWatchRuntimeFallbackDescriptorB524DiscoveryBucketsCorrectly(t *testing.T) {
	t.Parallel()

	spy := &watchEfficiencyObserverSpy{}
	poller := &vaillantSemanticPoller{
		watchEfficiency: spy,
		nowFn: func() time.Time {
			return time.Unix(1700000700, 0).UTC()
		},
	}
	key := engine.NewB524WatchKey(0x15, vaillantB524OpcodeLocal, localZones.group, 0x01, zoneIndex)
	runtime := poller.prepareSemanticReadWatchRuntime(key)
	if runtime.hasDescriptor {
		t.Fatal("runtime.hasDescriptor = true; want false when runtime observation descriptor is missing")
	}
	poller.emitWatchReadEfficiency(runtime, runtime.maxAge, engine.SemanticReadExecutionStats{
		ActiveFetchAttempted: true,
		ActiveFetchSucceeded: true,
		ActiveFetchDuration:  time.Second,
	})
	event, ok := spy.latestReadEvent()
	if !ok {
		t.Fatal("watch efficiency read event missing")
	}
	if event.HasDescriptor {
		t.Fatal("event.HasDescriptor = true; want false")
	}
	if event.Key == nil || event.Key.Family() != engine.WatchFamilyB524 {
		t.Fatalf("family = %v; want B524", event.Key)
	}
	if event.Descriptor.FreshnessProfile != engine.WatchFreshnessProfileDiscovery {
		t.Fatalf("freshness = %s; want discovery", event.Descriptor.FreshnessProfile)
	}
}

func TestVaillantOwnerReadB509ValueEmitsWatchEfficiencyShadowHit(t *testing.T) {
	t.Parallel()

	key := engine.NewB509WatchKey(0x08, 0x0200)
	spy := &watchEfficiencyObserverSpy{}
	cfg := observeFirstStateShadowConfig(key)
	cfg.WatchEfficiencyObserver = spy

	poller := newVaillantSemanticPollerFromDonor(
		cfg,
		&donorGateway{},
		newRetainedSemanticState(),
		nil,
		nil,
	)
	poller.sendFrameFn = func(context.Context, protocol.Frame) (*protocol.Frame, error) {
		return nil, fmt.Errorf("unexpected active send on shadow-hit path")
	}

	maxAge := poller.prepareSemanticReadWatch(key)
	if maxAge != 30*time.Second {
		t.Fatalf("prepareSemanticReadWatch() = %s; want 30s", maxAge)
	}
	write := poller.shadow.Write(engine.ShadowWrite{
		Key:        key,
		Source:     engine.ShadowWriteSourcePassive,
		Confidence: engine.ShadowConfidenceHigh,
		Value:      []byte{0x7A},
		ObservedAt: time.Now(),
	})
	if !write.Accepted {
		t.Fatalf("shadow write rejected: %s", write.Reason)
	}

	value, ok := poller.readB509Value(context.Background(), 0x08, 0x0200)
	if !ok {
		t.Fatal("readB509Value() ok = false; want true on shadow hit")
	}
	if len(value) != 1 || value[0] != 0x7A {
		t.Fatalf("readB509Value() value = %v; want [0x7a]", value)
	}

	event, found := spy.latestReadEvent()
	if !found {
		t.Fatal("watch-efficiency read event missing")
	}
	if !event.Stats.ServedFromShadow {
		t.Fatal("event.Stats.ServedFromShadow = false; want true")
	}
	if event.Stats.ActiveFetchAttempted {
		t.Fatal("event.Stats.ActiveFetchAttempted = true; want false on shadow hit")
	}
	if event.Descriptor.FreshnessProfile != engine.WatchFreshnessProfileStateFast {
		t.Fatalf("event descriptor freshness_profile = %q; want state_fast", event.Descriptor.FreshnessProfile)
	}
	if event.Descriptor.Family() != engine.WatchFamilyB509 {
		t.Fatalf("event descriptor family = %q; want B509", event.Descriptor.Family())
	}
}

func TestVaillantOwnerHandleAdjudicatedPassiveEventEmitsWatchEfficiencyDirectApply(t *testing.T) {
	t.Parallel()

	key := engine.NewB509WatchKey(0x08, 0x0200)
	spy := &watchEfficiencyObserverSpy{}
	cfg := observeFirstStateShadowRuntimeConfig(engine.ObserveFirstExternalWritePolicyRecordOnly)
	cfg.WatchEfficiencyObserver = spy

	poller := newVaillantSemanticPollerFromDonor(
		cfg,
		&donorGateway{},
		newRetainedSemanticState(),
		nil,
		nil,
	)
	now := time.Unix(1700000600, 0).UTC()
	poller.nowFn = func() time.Time { return now }

	poller.handleAdjudicatedPassiveEvent(engine.AdjudicatedPassiveEvent{
		Disposition: engine.DedupDispositionUnmatchedThirdParty,
		Fingerprint: engine.PassiveTransactionFingerprint{
			SharedWatchKey: key,
			ObservedAt:     now,
			ResponseClass:  engine.DedupResponseValueBearing,
			FamilyPolicy: engine.ObserveFirstFamilyPolicy{
				RequestIntent:     engine.ObserveFirstRequestIntentRead,
				DirectApplyPolicy: engine.ObserveFirstDirectApplyPolicyStateDefault,
			},
		},
		Event: engine.PassiveClassifiedEvent{
			HasResponse: true,
			Response: protocol.Frame{
				Data: []byte{vaillantB509OpcodeRead, 0x02, 0x00, 0x11},
			},
		},
	})

	event, found := spy.latestDirectApplyEvent()
	if !found {
		t.Fatal("watch-efficiency direct-apply event missing")
	}
	if !event.CandidateEvaluated {
		t.Fatal("event.CandidateEvaluated = false; want true")
	}
	if !event.Accepted {
		t.Fatal("event.Accepted = false; want true")
	}
	if event.Descriptor.FreshnessProfile != engine.WatchFreshnessProfileStateFast {
		t.Fatalf("event descriptor freshness_profile = %q; want state_fast", event.Descriptor.FreshnessProfile)
	}
	if event.Descriptor.Family() != engine.WatchFamilyB509 {
		t.Fatalf("event descriptor family = %q; want B509", event.Descriptor.Family())
	}
}

func TestVaillantOwnerHandleAdjudicatedPassiveEventEmitsWatchEfficiencyCandidateForRejectedDirectApply(t *testing.T) {
	t.Parallel()

	key := engine.NewB509WatchKey(0x08, 0x0200)
	spy := &watchEfficiencyObserverSpy{}
	cfg := observeFirstStateShadowRuntimeConfig(engine.ObserveFirstExternalWritePolicyRecordOnly)
	cfg.WatchEfficiencyObserver = spy

	poller := newVaillantSemanticPollerFromDonor(
		cfg,
		&donorGateway{},
		newRetainedSemanticState(),
		nil,
		nil,
	)
	now := time.Unix(1700000600, 0).UTC()
	poller.nowFn = func() time.Time { return now }
	_ = poller.prepareSemanticReadWatch(key)

	seed := poller.shadow.Write(engine.ShadowWrite{
		Key:        key,
		Source:     engine.ShadowWriteSourcePassive,
		Confidence: engine.ShadowConfidenceHigh,
		Value:      []byte{0x33},
		ObservedAt: now.Add(time.Second),
	})
	if !seed.Accepted {
		t.Fatalf("shadow seed write rejected: %s", seed.Reason)
	}

	poller.handleAdjudicatedPassiveEvent(engine.AdjudicatedPassiveEvent{
		Disposition: engine.DedupDispositionUnmatchedThirdParty,
		Fingerprint: engine.PassiveTransactionFingerprint{
			SharedWatchKey: key,
			ObservedAt:     now,
			ResponseClass:  engine.DedupResponseValueBearing,
			FamilyPolicy: engine.ObserveFirstFamilyPolicy{
				RequestIntent:     engine.ObserveFirstRequestIntentRead,
				DirectApplyPolicy: engine.ObserveFirstDirectApplyPolicyStateDefault,
			},
		},
		Event: engine.PassiveClassifiedEvent{
			HasResponse: true,
			Response: protocol.Frame{
				Data: []byte{vaillantB509OpcodeRead, 0x02, 0x00, 0x11},
			},
		},
	})

	event, found := spy.latestDirectApplyEvent()
	if !found {
		t.Fatal("watch-efficiency direct-apply event missing")
	}
	if !event.CandidateEvaluated {
		t.Fatal("event.CandidateEvaluated = false; want true")
	}
	if event.Accepted {
		t.Fatal("event.Accepted = true; want false for stale rejected write")
	}
}

func TestVaillantOwnerRefreshDiscoveryPrimesBoilerTiersWhenBoilerAppears(t *testing.T) {
	t.Parallel()

	taskScheduler := newSemanticTaskScheduler()
	poller := &vaillantSemanticPoller{
		reg:            newTestRegistry(registry.DeviceInfo{Address: 0x08, Manufacturer: "Vaillant", DeviceID: "BAI00"}),
		tasks:          taskScheduler,
		zones:          make(map[byte]*vaillantZoneSnapshot),
		presence:       make(map[byte]*zonePresenceRecord),
		circuits:       make(map[byte]*vaillantCircuitSnapshot),
		radioDevices:   make(map[radioDeviceKey]*vaillantRadioDeviceSnapshot),
		solarCylinders: make(map[byte]*vaillantCylinderSnapshot),
	}

	poller.refreshDiscovery(context.Background())

	poller.mu.Lock()
	boilerAddress := poller.boilerAddress
	poller.mu.Unlock()
	if boilerAddress != 0x08 {
		t.Fatalf("boilerAddress = 0x%02x; want 0x08", boilerAddress)
	}

	taskScheduler.mu.Lock()
	queueLen := len(taskScheduler.queue)
	taskScheduler.mu.Unlock()
	if queueLen != len(poller.boilerStatusTierSchedules()) {
		t.Fatalf("queued boiler tasks = %d; want %d", queueLen, len(poller.boilerStatusTierSchedules()))
	}

	poller.refreshDiscovery(context.Background())

	taskScheduler.mu.Lock()
	queueLen = len(taskScheduler.queue)
	taskScheduler.mu.Unlock()
	if queueLen != len(poller.boilerStatusTierSchedules()) {
		t.Fatalf("queued boiler tasks after unchanged discovery = %d; want unchanged %d", queueLen, len(poller.boilerStatusTierSchedules()))
	}
}
