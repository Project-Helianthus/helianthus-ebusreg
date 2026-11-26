package owner

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

func TestVaillantOwnerHandleAdjudicatedPassiveEventUnmatchedValueBearingReadWritesShadow(t *testing.T) {
	t.Parallel()

	key := engine.NewB509WatchKey(0x08, 0x0200)
	poller := newVaillantSemanticAcquisition(
		observeFirstStateShadowRuntimeConfig(engine.ObserveFirstExternalWritePolicyRecordOnly),
		Substrate{},
		nil,
	)

	observedAt := time.Now()
	event := engine.AdjudicatedPassiveEvent{
		Event: engine.PassiveClassifiedEvent{
			HasRequest:  true,
			HasResponse: true,
			Request: protocol.Frame{
				Source:    0x10,
				Target:    0x08,
				Primary:   vaillantB509Primary,
				Secondary: vaillantB509Secondary,
				Data:      buildB509ReadSelector(0x0200),
			},
			Response: protocol.Frame{
				Source:    0x08,
				Target:    0x10,
				Primary:   vaillantB509Primary,
				Secondary: vaillantB509Secondary,
				Data:      []byte{vaillantB509OpcodeRead, 0x02, 0x00, 0xAA, 0x55},
			},
		},
		Fingerprint: engine.PassiveTransactionFingerprint{
			SharedWatchKey: key,
			ObservedAt:     observedAt,
			ResponseClass:  engine.DedupResponseValueBearing,
			FamilyPolicy: engine.ObserveFirstFamilyPolicy{
				RequestIntent:     engine.ObserveFirstRequestIntentRead,
				ResponseClass:     engine.DedupResponseValueBearing,
				DirectApplyPolicy: engine.ObserveFirstDirectApplyPolicyStateDefault,
			},
		},
		Disposition: engine.DedupDispositionUnmatchedThirdParty,
	}

	poller.handleAdjudicatedPassiveEvent(event)

	lookup := poller.shadow.Lookup(key, 10*time.Second)
	if !lookup.Found {
		t.Fatal("shadow lookup found = false; want passive value-bearing read in shadow")
	}
	if !lookup.Eligible {
		t.Fatalf("shadow lookup eligible = false; want true for freshly written passive read (state=%s)", lookup.Entry.State)
	}
	if !slices.Equal(lookup.Entry.Value, []byte{0xAA, 0x55}) {
		t.Fatalf("shadow value = %x; want %x", lookup.Entry.Value, []byte{0xAA, 0x55})
	}
	if lookup.Entry.ObservedAt.IsZero() || !lookup.Entry.ObservedAt.Equal(observedAt) {
		t.Fatalf("shadow observed_at = %s; want %s", lookup.Entry.ObservedAt.UTC().Format(time.RFC3339), observedAt.UTC().Format(time.RFC3339))
	}
	if lookup.Entry.Pinned {
		t.Fatal("shadow entry pinned = true; want passive fallback bootstrap to remain unpinned")
	}
	if summary := poller.shadow.Summary(); summary.StaticPinnedFootprint != 0 {
		t.Fatalf("Summary().StaticPinnedFootprint = %d; want 0 for passive fallback bootstrap", summary.StaticPinnedFootprint)
	}
}

func TestVaillantOwnerHandleAdjudicatedPassiveEventUnmatchedValueBearingB524ReadWritesShadow(t *testing.T) {
	t.Parallel()

	key := engine.NewB524WatchKey(0x15, 0x02, 0x03, 0x01, 0x001C)
	poller := newVaillantSemanticAcquisition(
		observeFirstStateShadowRuntimeConfig(engine.ObserveFirstExternalWritePolicyRecordOnly),
		Substrate{},
		nil,
	)

	observedAt := time.Now()
	event := engine.AdjudicatedPassiveEvent{
		Event: engine.PassiveClassifiedEvent{
			HasRequest:  true,
			HasResponse: true,
			Request: protocol.Frame{
				Source:    0x10,
				Target:    0x15,
				Primary:   vaillantExtRegisterPrimary,
				Secondary: vaillantExtRegisterSecondary,
				Data:      buildB524ReadSelector(0x02, 0x03, 0x01, 0x001C),
			},
			Response: protocol.Frame{
				Source:    0x15,
				Target:    0x10,
				Primary:   vaillantExtRegisterPrimary,
				Secondary: vaillantExtRegisterSecondary,
				Data:      []byte{0x02, 0x01, 0x03, 0x1C, 0x00, 0x42, 0x01},
			},
		},
		Fingerprint: engine.PassiveTransactionFingerprint{
			SharedWatchKey: key,
			ObservedAt:     observedAt,
			ResponseClass:  engine.DedupResponseValueBearing,
			FamilyPolicy: engine.ObserveFirstFamilyPolicy{
				RequestIntent:     engine.ObserveFirstRequestIntentRead,
				ResponseClass:     engine.DedupResponseValueBearing,
				DirectApplyPolicy: engine.ObserveFirstDirectApplyPolicyStateDefault,
			},
		},
		Disposition: engine.DedupDispositionUnmatchedThirdParty,
	}

	poller.handleAdjudicatedPassiveEvent(event)

	lookup := poller.shadow.Lookup(key, 10*time.Second)
	if !lookup.Found {
		t.Fatal("shadow lookup found = false; want passive B524 value-bearing read in shadow")
	}
	if !lookup.Eligible {
		t.Fatalf("shadow lookup eligible = false; want true for freshly written passive B524 read (state=%s)", lookup.Entry.State)
	}
	if !slices.Equal(lookup.Entry.Value, []byte{0x42, 0x01}) {
		t.Fatalf("shadow value = %x; want %x", lookup.Entry.Value, []byte{0x42, 0x01})
	}
}

func TestVaillantOwnerPassiveShadowLaneEnabledEnergyMergeOnlyPolicyDisabled(t *testing.T) {
	t.Parallel()

	flags := engine.NormalizeObserveFirstFeatureFlags(
		true,
		true,
		true,
		engine.ObserveFirstExternalWritePolicyRecordOnly,
	)
	policy := engine.ObserveFirstFamilyPolicy{
		RequestIntent:     engine.ObserveFirstRequestIntentRead,
		DirectApplyPolicy: engine.ObserveFirstDirectApplyPolicyEnergyMergeOnly,
	}

	if passiveShadowLaneEnabled(flags, policy) {
		t.Fatal("passiveShadowLaneEnabled() = true; want false for energy_merge_only carve-out")
	}
}

func TestVaillantOwnerHandleAdjudicatedPassiveEventDirectApplyAdmissionHonorsNormalizedPolicyFlags(t *testing.T) {
	t.Parallel()

	stateKey := engine.NewB509WatchKey(0x08, 0x0200)
	configKey := engine.NewB524WatchKey(0x15, 0x02, 0x03, 0x01, 0x001C)
	tests := []struct {
		name   string
		flags  engine.ObserveFirstFeatureFlags
		key    engine.WatchKey
		policy engine.ObserveFirstDirectApplyPolicy
		want   bool
	}{
		{
			name:   "state default is admitted",
			flags:  engine.NormalizeObserveFirstFeatureFlags(true, true, false, engine.ObserveFirstExternalWritePolicyRecordOnly),
			key:    stateKey,
			policy: engine.ObserveFirstDirectApplyPolicyStateDefault,
			want:   true,
		},
		{
			name:   "config opt in remains denied without runtime admission",
			flags:  engine.NormalizeObserveFirstFeatureFlags(true, true, true, engine.ObserveFirstExternalWritePolicyRecordAndInvalidate),
			key:    configKey,
			policy: engine.ObserveFirstDirectApplyPolicyConfigOptIn,
			want:   false,
		},
		{
			name:   "global disable is denied",
			flags:  engine.NormalizeObserveFirstFeatureFlags(false, true, true, engine.ObserveFirstExternalWritePolicyRecordAndInvalidate),
			key:    stateKey,
			policy: engine.ObserveFirstDirectApplyPolicyStateDefault,
			want:   false,
		},
		{
			name:   "state disabled normalizes config direct apply off",
			flags:  engine.NormalizeObserveFirstFeatureFlags(true, false, true, engine.ObserveFirstExternalWritePolicyRecordAndInvalidate),
			key:    configKey,
			policy: engine.ObserveFirstDirectApplyPolicyConfigOptIn,
			want:   false,
		},
		{
			name:   "never is denied",
			flags:  engine.NormalizeObserveFirstFeatureFlags(true, true, true, engine.ObserveFirstExternalWritePolicyRecordAndInvalidate),
			key:    stateKey,
			policy: engine.ObserveFirstDirectApplyPolicyNever,
			want:   false,
		},
		{
			name:   "energy merge only is denied",
			flags:  engine.NormalizeObserveFirstFeatureFlags(true, true, true, engine.ObserveFirstExternalWritePolicyRecordAndInvalidate),
			key:    stateKey,
			policy: engine.ObserveFirstDirectApplyPolicyEnergyMergeOnly,
			want:   false,
		},
		{
			name:   "unknown policy is denied",
			flags:  engine.NormalizeObserveFirstFeatureFlags(true, true, true, engine.ObserveFirstExternalWritePolicyRecordAndInvalidate),
			key:    stateKey,
			policy: engine.ObserveFirstDirectApplyPolicy("unknown"),
			want:   false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			poller := newVaillantSemanticAcquisition(
				Config{ObserveFirstFlags: test.flags},
				Substrate{},
				nil,
			)
			event := passiveDirectApplyTestEvent(test.key, test.policy)
			poller.handleAdjudicatedPassiveEvent(event)

			lookup := poller.shadow.Lookup(test.key, 10*time.Second)
			if lookup.Found != test.want {
				t.Fatalf("shadow found = %t; want %t for policy=%q flags=%+v", lookup.Found, test.want, test.policy, test.flags)
			}
		})
	}
}

func passiveDirectApplyTestEvent(key engine.WatchKey, policy engine.ObserveFirstDirectApplyPolicy) engine.AdjudicatedPassiveEvent {
	request := protocol.Frame{Source: 0x10, Target: 0x08, Primary: vaillantB509Primary, Secondary: vaillantB509Secondary, Data: buildB509ReadSelector(0x0200)}
	response := protocol.Frame{Source: 0x08, Target: 0x10, Primary: vaillantB509Primary, Secondary: vaillantB509Secondary, Data: []byte{vaillantB509OpcodeRead, 0x02, 0x00, 0xAA}}
	if typed, ok := key.(engine.B524WatchKey); ok {
		request = protocol.Frame{Source: 0x10, Target: typed.Target, Primary: vaillantExtRegisterPrimary, Secondary: vaillantExtRegisterSecondary, Data: buildB524ReadSelector(typed.Opcode, typed.Group, typed.Instance, typed.RegisterAddress)}
		response = protocol.Frame{Source: typed.Target, Target: 0x10, Primary: vaillantExtRegisterPrimary, Secondary: vaillantExtRegisterSecondary, Data: []byte{typed.Opcode, typed.Instance, typed.Group, byte(typed.RegisterAddress), byte(typed.RegisterAddress >> 8), 0xAA}}
	}
	return engine.AdjudicatedPassiveEvent{
		Event: engine.PassiveClassifiedEvent{
			HasRequest:  true,
			HasResponse: true,
			Request:     request,
			Response:    response,
		},
		Fingerprint: engine.PassiveTransactionFingerprint{
			SharedWatchKey: key,
			ObservedAt:     time.Unix(400, 0),
			ResponseClass:  engine.DedupResponseValueBearing,
			FamilyPolicy: engine.ObserveFirstFamilyPolicy{
				RequestIntent:     engine.ObserveFirstRequestIntentRead,
				ResponseClass:     engine.DedupResponseValueBearing,
				DirectApplyPolicy: policy,
			},
		},
		Disposition: engine.DedupDispositionUnmatchedThirdParty,
	}
}

func TestVaillantOwnerHandleAdjudicatedPassiveEventUnmatchedExternalWriteInvalidatesShadow(t *testing.T) {
	t.Parallel()

	key := engine.NewB509WatchKey(0x08, 0x0200)
	poller := newVaillantSemanticAcquisition(
		observeFirstStateShadowRuntimeConfig(engine.ObserveFirstExternalWritePolicyInvalidateOnly),
		Substrate{},
		nil,
	)

	_ = poller.prepareSemanticReadWatch(key)
	seedAt := time.Unix(100, 0)
	seed := poller.shadow.Write(engine.ShadowWrite{
		Key:        key,
		Source:     engine.ShadowWriteSourcePassive,
		Confidence: engine.ShadowConfidenceHigh,
		Value:      []byte{0x11, 0x22},
		ObservedAt: seedAt,
	})
	if !seed.Accepted {
		t.Fatalf("seed shadow write rejected: %s", seed.Reason)
	}

	invalidatedAt := seedAt.Add(5 * time.Second)
	event := engine.AdjudicatedPassiveEvent{
		Event: engine.PassiveClassifiedEvent{
			HasRequest: true,
			Request: protocol.Frame{
				Source:    0x10,
				Target:    0x08,
				Primary:   vaillantB509Primary,
				Secondary: vaillantB509Secondary,
				Data:      buildB509WriteSelector(0x0200, []byte{0x33, 0x44}),
			},
		},
		Fingerprint: engine.PassiveTransactionFingerprint{
			SharedWatchKey: key,
			ObservedAt:     invalidatedAt,
			FamilyPolicy: engine.ObserveFirstFamilyPolicy{
				RequestIntent:                  engine.ObserveFirstRequestIntentWrite,
				UsesRuntimeExternalWritePolicy: true,
				EffectiveExternalWritePolicy:   engine.ObserveFirstExternalWritePolicyInvalidateOnly,
			},
		},
		Disposition: engine.DedupDispositionUnmatchedThirdParty,
	}

	poller.handleAdjudicatedPassiveEvent(event)

	lookup := poller.shadow.Lookup(key, 10*time.Second)
	if !lookup.Found {
		t.Fatal("shadow lookup found = false; want invalidated entry preserved")
	}
	if lookup.Eligible {
		t.Fatalf("shadow lookup eligible = true; want false after external-write invalidation (state=%s)", lookup.Entry.State)
	}
	if lookup.Entry.State != engine.ShadowEntryStateInvalidated && lookup.Entry.State != engine.ShadowEntryStateTombstone {
		t.Fatalf("shadow state = %s; want invalidated or tombstone", lookup.Entry.State)
	}
}

func TestVaillantOwnerHandleAdjudicatedPassiveEventIgnoresNonApplicableEvents(t *testing.T) {
	t.Parallel()

	baseEvent := func() engine.AdjudicatedPassiveEvent {
		return engine.AdjudicatedPassiveEvent{
			Event: engine.PassiveClassifiedEvent{
				HasRequest:  true,
				HasResponse: true,
				Request: protocol.Frame{
					Source:    0x10,
					Target:    0x08,
					Primary:   vaillantB509Primary,
					Secondary: vaillantB509Secondary,
					Data:      buildB509ReadSelector(0x0200),
				},
				Response: protocol.Frame{
					Source:    0x08,
					Target:    0x10,
					Primary:   vaillantB509Primary,
					Secondary: vaillantB509Secondary,
					Data:      []byte{vaillantB509OpcodeRead, 0x02, 0x00, 0xAA},
				},
			},
			Fingerprint: engine.PassiveTransactionFingerprint{
				SharedWatchKey: engine.NewB509WatchKey(0x08, 0x0200),
				ObservedAt:     time.Unix(300, 0),
				ResponseClass:  engine.DedupResponseValueBearing,
				FamilyPolicy: engine.ObserveFirstFamilyPolicy{
					RequestIntent:     engine.ObserveFirstRequestIntentRead,
					ResponseClass:     engine.DedupResponseValueBearing,
					DirectApplyPolicy: engine.ObserveFirstDirectApplyPolicyStateDefault,
				},
			},
			Disposition: engine.DedupDispositionUnmatchedThirdParty,
		}
	}

	tests := []struct {
		name   string
		mutate func(*engine.AdjudicatedPassiveEvent)
	}{
		{
			name: "matched active duplicate",
			mutate: func(event *engine.AdjudicatedPassiveEvent) {
				event.MatchedActiveDuplicate = true
			},
		},
		{
			name: "observability only",
			mutate: func(event *engine.AdjudicatedPassiveEvent) {
				event.ObservabilityOnly = true
			},
		},
		{
			name: "local participant inbound",
			mutate: func(event *engine.AdjudicatedPassiveEvent) {
				event.LocalParticipantInbound = true
			},
		},
		{
			name: "header only response class",
			mutate: func(event *engine.AdjudicatedPassiveEvent) {
				event.Fingerprint.ResponseClass = engine.DedupResponseHeaderOnly
				event.Fingerprint.FamilyPolicy.ResponseClass = engine.DedupResponseHeaderOnly
			},
		},
		{
			name: "ack only response class",
			mutate: func(event *engine.AdjudicatedPassiveEvent) {
				event.Fingerprint.ResponseClass = engine.DedupResponseACKOnly
				event.Fingerprint.FamilyPolicy.ResponseClass = engine.DedupResponseACKOnly
			},
		},
		{
			name: "no response bytes",
			mutate: func(event *engine.AdjudicatedPassiveEvent) {
				event.Event.HasResponse = false
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			poller := newVaillantSemanticAcquisition(
				Config{},
				Substrate{},
				nil,
			)
			event := baseEvent()
			test.mutate(&event)

			poller.handleAdjudicatedPassiveEvent(event)

			key := engine.NewB509WatchKey(0x08, 0x0200)
			lookup := poller.shadow.Lookup(key, 10*time.Second)
			if lookup.Found {
				t.Fatalf("shadow lookup found = true; want ignored event to skip shadow mutation (state=%s value=%x)", lookup.Entry.State, lookup.Entry.Value)
			}
		})
	}
}

func TestVaillantOwnerHandleAdjudicatedPassiveEventReadSkipsShadowWhenObserveFirstMasterDisabled(t *testing.T) {
	t.Parallel()

	key := engine.NewB509WatchKey(0x08, 0x0200)
	cfg := observeFirstStateShadowConfig(key)
	cfg.ObserveFirstFlags = engine.DefaultObserveFirstFeatureFlags()
	poller := newVaillantSemanticAcquisition(
		cfg,
		Substrate{},
		nil,
	)

	event := engine.AdjudicatedPassiveEvent{
		Event: engine.PassiveClassifiedEvent{
			HasRequest:  true,
			HasResponse: true,
			Request: protocol.Frame{
				Source:    0x10,
				Target:    0x08,
				Primary:   vaillantB509Primary,
				Secondary: vaillantB509Secondary,
				Data:      buildB509ReadSelector(0x0200),
			},
			Response: protocol.Frame{
				Source:    0x08,
				Target:    0x10,
				Primary:   vaillantB509Primary,
				Secondary: vaillantB509Secondary,
				Data:      []byte{vaillantB509OpcodeRead, 0x02, 0x00, 0xAA},
			},
		},
		Fingerprint: engine.PassiveTransactionFingerprint{
			SharedWatchKey: key,
			ObservedAt:     time.Unix(310, 0),
			ResponseClass:  engine.DedupResponseValueBearing,
			FamilyPolicy: engine.ObserveFirstFamilyPolicy{
				RequestIntent:     engine.ObserveFirstRequestIntentRead,
				ResponseClass:     engine.DedupResponseValueBearing,
				DirectApplyPolicy: engine.ObserveFirstDirectApplyPolicyStateDefault,
			},
		},
		Disposition: engine.DedupDispositionUnmatchedThirdParty,
	}

	poller.handleAdjudicatedPassiveEvent(event)

	if lookup := poller.shadow.Lookup(key, 10*time.Second); lookup.Found {
		t.Fatalf("shadow lookup found = true; want observe-first global gate to block passive shadow write (state=%s value=%x)", lookup.Entry.State, lookup.Entry.Value)
	}
}

func TestVaillantOwnerHandleAdjudicatedPassiveEventReadWithoutSharedWatchKeyDoesNotPolluteShadow(t *testing.T) {
	t.Parallel()

	key := engine.NewB509WatchKey(0x08, 0x0200)
	poller := newVaillantSemanticAcquisition(
		observeFirstStateShadowRuntimeConfig(engine.ObserveFirstExternalWritePolicyRecordOnly),
		Substrate{},
		nil,
	)

	event := engine.AdjudicatedPassiveEvent{
		Event: engine.PassiveClassifiedEvent{
			HasRequest:  true,
			HasResponse: true,
			Request: protocol.Frame{
				Source:    0x10,
				Target:    0x08,
				Primary:   vaillantB509Primary,
				Secondary: vaillantB509Secondary,
				Data:      buildB509ReadSelector(0x0200),
			},
			Response: protocol.Frame{
				Source:    0x08,
				Target:    0x10,
				Primary:   vaillantB509Primary,
				Secondary: vaillantB509Secondary,
				Data:      []byte{vaillantB509OpcodeRead, 0x02, 0x00, 0xAB},
			},
		},
		Fingerprint: engine.PassiveTransactionFingerprint{
			ObservedAt:    time.Unix(320, 0),
			ResponseClass: engine.DedupResponseValueBearing,
			FamilyPolicy: engine.ObserveFirstFamilyPolicy{
				RequestIntent:     engine.ObserveFirstRequestIntentRead,
				ResponseClass:     engine.DedupResponseValueBearing,
				DirectApplyPolicy: engine.ObserveFirstDirectApplyPolicyStateDefault,
			},
		},
		Disposition: engine.DedupDispositionUnmatchedThirdParty,
	}

	poller.handleAdjudicatedPassiveEvent(event)

	if lookup := poller.shadow.Lookup(key, 10*time.Second); lookup.Found {
		t.Fatalf("shadow lookup found = true; want shared-watchkey-absent traffic to be ignored (state=%s value=%x)", lookup.Entry.State, lookup.Entry.Value)
	}
}

func TestVaillantOwnerHandleAdjudicatedPassiveEventReadUnknownWatchKeyDoesNotPolluteShadow(t *testing.T) {
	t.Parallel()

	key := engine.NewB509WatchKey(0x08, 0x0200)
	cfg := observeFirstStateShadowRuntimeConfig(engine.ObserveFirstExternalWritePolicyRecordOnly)
	poller := newVaillantSemanticAcquisition(
		cfg,
		Substrate{},
		nil,
	)

	event := engine.AdjudicatedPassiveEvent{
		Event: engine.PassiveClassifiedEvent{
			HasRequest:  true,
			HasResponse: true,
			Request: protocol.Frame{
				Source:    0x10,
				Target:    0x08,
				Primary:   vaillantB509Primary,
				Secondary: vaillantB509Secondary,
				Data:      buildB509ReadSelector(0x0200),
			},
			Response: protocol.Frame{
				Source:    0x08,
				Target:    0x10,
				Primary:   vaillantB509Primary,
				Secondary: vaillantB509Secondary,
				Data:      []byte{vaillantB509OpcodeRead, 0x02, 0x00, 0xAC},
			},
		},
		Fingerprint: engine.PassiveTransactionFingerprint{
			SharedWatchKey: unknownSemanticWatchKey{canonical: "unknown:08:0200"},
			ObservedAt:     time.Unix(330, 0),
			ResponseClass:  engine.DedupResponseValueBearing,
			FamilyPolicy: engine.ObserveFirstFamilyPolicy{
				RequestIntent:     engine.ObserveFirstRequestIntentRead,
				ResponseClass:     engine.DedupResponseValueBearing,
				DirectApplyPolicy: engine.ObserveFirstDirectApplyPolicyStateDefault,
			},
		},
		Disposition: engine.DedupDispositionUnmatchedThirdParty,
	}

	poller.handleAdjudicatedPassiveEvent(event)

	if lookup := poller.shadow.Lookup(key, 10*time.Second); lookup.Found {
		t.Fatalf("shadow lookup found = true; want unknown shared watch key to be ignored (state=%s value=%x)", lookup.Entry.State, lookup.Entry.Value)
	}
}

func TestVaillantOwnerAttachPassiveShadowProducerDedupSuppressedWriteStillInvalidatesShadow(t *testing.T) {
	policies := []engine.ObserveFirstExternalWritePolicy{
		engine.ObserveFirstExternalWritePolicyInvalidateOnly,
		engine.ObserveFirstExternalWritePolicyRecordAndInvalidate,
	}

	for _, policy := range policies {
		policy := policy
		t.Run(string(policy), func(t *testing.T) {
			key := engine.NewB509WatchKey(0x08, 0x0200)
			cfg := observeFirstStateShadowRuntimeConfig(policy)
			dedupConfig := engine.Config{
				ObserveFirstFlags:                  cfg.ObserveFirstFlags,
				PassiveDedupRecoveryHysteresis:     time.Nanosecond,
				PassiveDedupRecoveryEventThreshold: 2,
				Collaborator:                       EngineCollaborator{},
			}

			deduplicator, err := engine.NewActivePassiveDeduplicator(dedupConfig)
			if err != nil {
				t.Fatalf("NewActivePassiveDeduplicator() error = %v", err)
			}
			t.Cleanup(func() {
				_ = deduplicator.Close()
			})

			poller := newVaillantSemanticAcquisition(
				cfg,
				Substrate{},
				nil,
			)

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)
			if err := poller.AttachPassiveShadowProducer(ctx, deduplicator); err != nil {
				t.Fatalf("AttachPassiveShadowProducer() error = %v", err)
			}

			witness, err := deduplicator.Subscribe("shadow-regression-witness", engine.DedupSubscriberCritical, 128)
			if err != nil {
				t.Fatalf("Subscribe(witness) error = %v", err)
			}
			t.Cleanup(witness.Close)

			_ = poller.prepareSemanticReadWatch(key)
			seedAt := time.Now()
			seed := poller.shadow.Write(engine.ShadowWrite{
				Key:        key,
				Source:     engine.ShadowWriteSourcePassive,
				Confidence: engine.ShadowConfidenceHigh,
				Value:      []byte{0x11, 0x22},
				ObservedAt: seedAt,
			})
			if !seed.Accepted {
				t.Fatalf("seed shadow write rejected: %s", seed.Reason)
			}

			base := seedAt.Add(10 * time.Millisecond)
			deduplicator.OnPassiveClassifiedEvent(passiveBroadcastClassifiedEvent(base))
			deduplicator.OnPassiveClassifiedEvent(passiveB509WriteClassifiedEvent(base.Add(2*time.Nanosecond), 0x10, 0x08, 0x0200, []byte{0x33, 0x44}))
			deduplicator.OnPassiveClassifiedEvent(passiveBroadcastClassifiedEvent(base.Add(deduplicator.Budgets().PendingGraceTimeout + 4*time.Nanosecond)))

			adjudicated := waitForAdjudicatedEvent(t, witness, 2*time.Second, func(event engine.AdjudicatedPassiveEvent) bool {
				return event.Disposition == engine.DedupDispositionUnmatchedThirdParty &&
					event.FamilyPolicy.RequestIntent == engine.ObserveFirstRequestIntentWrite
			})
			if !adjudicated.SuppressShadow {
				t.Fatalf("SuppressShadow = false; want true for policy %q", policy)
			}
			if adjudicated.Fingerprint.SharedWatchKey == nil {
				t.Fatal("SharedWatchKey = nil; want dedup-emitted shared key for passive write invalidation")
			}
			if got := adjudicated.Fingerprint.SharedWatchKey.Canonical(); got != key.Canonical() {
				t.Fatalf("SharedWatchKey.Canonical() = %q; want %q", got, key.Canonical())
			}

			lookup := waitForShadowIneligible(t, poller.shadow, key, 2*time.Second)
			if lookup.Entry.State != engine.ShadowEntryStateInvalidated && lookup.Entry.State != engine.ShadowEntryStateTombstone {
				t.Fatalf("shadow state = %s; want invalidated or tombstone", lookup.Entry.State)
			}
			if lookup.Entry.InvalidationReason != engine.ShadowInvalidationReasonExternalWrite {
				t.Fatalf("invalidation reason = %s; want %s", lookup.Entry.InvalidationReason, engine.ShadowInvalidationReasonExternalWrite)
			}
			if lookup.Entry.InvalidationSource != engine.ShadowInvalidationSourcePassive {
				t.Fatalf("invalidation source = %s; want %s", lookup.Entry.InvalidationSource, engine.ShadowInvalidationSourcePassive)
			}
		})
	}
}

func TestVaillantOwnerAttachPassiveShadowProducerResubscribesAfterCriticalOverflow(t *testing.T) {
	originalPriority := passiveShadowSubscriberPriority
	originalBuffer := passiveShadowSubscriberBuffer
	originalRetryDelay := passiveShadowRetryDelay
	passiveShadowSubscriberPriority = engine.DedupSubscriberCritical
	passiveShadowSubscriberBuffer = 1
	passiveShadowRetryDelay = time.Millisecond
	t.Cleanup(func() {
		passiveShadowSubscriberPriority = originalPriority
		passiveShadowSubscriberBuffer = originalBuffer
		passiveShadowRetryDelay = originalRetryDelay
	})

	key := engine.NewB509WatchKey(0x08, 0x0200)
	cfg := observeFirstStateShadowRuntimeConfig(engine.ObserveFirstExternalWritePolicyInvalidateOnly)
	dedupConfig := engine.Config{
		ObserveFirstFlags:                  cfg.ObserveFirstFlags,
		PassiveDedupRecoveryHysteresis:     time.Nanosecond,
		PassiveDedupRecoveryEventThreshold: 2,
		Collaborator:                       EngineCollaborator{},
	}

	deduplicator, err := engine.NewActivePassiveDeduplicator(dedupConfig)
	if err != nil {
		t.Fatalf("NewActivePassiveDeduplicator() error = %v", err)
	}
	t.Cleanup(func() {
		_ = deduplicator.Close()
	})

	poller := newVaillantSemanticAcquisition(
		cfg,
		Substrate{},
		nil,
	)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if err := poller.AttachPassiveShadowProducer(ctx, deduplicator); err != nil {
		t.Fatalf("AttachPassiveShadowProducer() error = %v", err)
	}

	witness, err := deduplicator.Subscribe("shadow-overflow-witness", engine.DedupSubscriberCritical, 4096)
	if err != nil {
		t.Fatalf("Subscribe(witness) error = %v", err)
	}
	t.Cleanup(witness.Close)

	_ = poller.prepareSemanticReadWatch(key)
	seedAt := time.Now()
	seed := poller.shadow.Write(engine.ShadowWrite{
		Key:        key,
		Source:     engine.ShadowWriteSourcePassive,
		Confidence: engine.ShadowConfidenceHigh,
		Value:      []byte{0x01},
		ObservedAt: seedAt,
	})
	if !seed.Accepted {
		t.Fatalf("seed shadow write rejected: %s", seed.Reason)
	}

	base := seedAt.Add(10 * time.Millisecond)
	deduplicator.OnPassiveClassifiedEvent(passiveBroadcastClassifiedEvent(base))
	for index := 0; index < 256; index++ {
		deduplicator.OnPassiveClassifiedEvent(passiveB509WriteClassifiedEvent(base.Add(2*time.Nanosecond), 0x10, 0x08, 0x0200, []byte{byte(index)}))
	}
	deduplicator.OnPassiveClassifiedEvent(passiveBroadcastClassifiedEvent(base.Add(deduplicator.Budgets().PendingGraceTimeout + 4*time.Nanosecond)))

	_ = waitForAdjudicatedEvent(t, witness, 2*time.Second, func(event engine.AdjudicatedPassiveEvent) bool {
		return event.Disposition == engine.DedupDispositionDiscontinuity &&
			event.Event.DiscontinuityReason == engine.PassiveDiscontinuityCriticalSubscriberFault
	})
	time.Sleep(50 * time.Millisecond)

	reseedAt := base.Add(time.Second)
	startGeneration := poller.shadow.SnapshotEligibility(key).Generation
	reseed := poller.shadow.Write(engine.ShadowWrite{
		Key:             key,
		Source:          engine.ShadowWriteSourceActiveConfirmed,
		Confidence:      engine.ShadowConfidenceHigh,
		Value:           []byte{0x7A},
		ObservedAt:      reseedAt,
		StartGeneration: startGeneration,
	})
	if !reseed.Accepted {
		t.Fatalf("reseed shadow write rejected: %s", reseed.Reason)
	}
	if lookup := poller.shadow.Lookup(key, 10*time.Second); !lookup.Found || lookup.Entry.State != engine.ShadowEntryStatePresent {
		t.Fatalf("shadow reseed lookup = %+v; want present seeded entry before resubscribe verification", lookup)
	}

	wave := reseedAt.Add(10 * time.Millisecond)
	deduplicator.OnPassiveClassifiedEvent(passiveBroadcastClassifiedEvent(wave))
	deduplicator.OnPassiveClassifiedEvent(passiveB509WriteClassifiedEvent(wave.Add(2*time.Nanosecond), 0x11, 0x08, 0x0200, []byte{0xAA}))
	deduplicator.OnPassiveClassifiedEvent(passiveBroadcastClassifiedEvent(wave.Add(deduplicator.Budgets().PendingGraceTimeout + 4*time.Nanosecond)))

	lookup := waitForShadowIneligible(t, poller.shadow, key, 2*time.Second)
	if lookup.Entry.InvalidationSource != engine.ShadowInvalidationSourcePassive {
		t.Fatalf("invalidation source = %s; want %s", lookup.Entry.InvalidationSource, engine.ShadowInvalidationSourcePassive)
	}
}
