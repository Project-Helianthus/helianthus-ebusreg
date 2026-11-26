package owner

import (
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

func observeFirstStateShadowRuntimeConfig(externalWritePolicy engine.ObserveFirstExternalWritePolicy) Config {
	return Config{
		ObserveFirstFlags: engine.NormalizeObserveFirstFeatureFlags(
			true,
			true,
			false,
			externalWritePolicy,
		),
	}
}

func observeFirstStateShadowConfig(key engine.WatchKey) Config {
	descriptor := engine.WatchDescriptor{
		Key:               key,
		SemanticClass:     engine.WatchSemanticClassState,
		FreshnessProfile:  engine.WatchFreshnessProfileStateFast,
		DecoderID:         "test.semantic.passive.shadow",
		CorrelationPolicy: engine.WatchCorrelationPolicyRequestResponse,
		DirectApplyPolicy: engine.WatchDirectApplyPolicyStateDefault,
	}
	return Config{
		ObserveFirstFlags: engine.NormalizeObserveFirstFeatureFlags(
			true,
			true,
			false,
			engine.ObserveFirstExternalWritePolicyRecordOnly,
		),
		WatchObserver: staticSemanticReadWatchObserver{
			observation: engine.WatchObservation{
				State:         engine.WatchObservationStateActive,
				Descriptor:    descriptor,
				HasDescriptor: true,
				Sources:       []engine.WatchActivationSource{engine.WatchActivationSourcePoller},
			},
		},
	}
}

func passiveBroadcastClassifiedEvent(observedAt time.Time) engine.PassiveClassifiedEvent {
	return engine.PassiveClassifiedEvent{
		Kind:       engine.PassiveClassifiedEventBroadcastFrame,
		FrameType:  protocol.FrameTypeBroadcast,
		ObservedAt: observedAt,
	}
}

func passiveB509WriteClassifiedEvent(observedAt time.Time, source, target byte, addr uint16, value []byte) engine.PassiveClassifiedEvent {
	request := protocol.Frame{
		Source:    source,
		Target:    target,
		Primary:   vaillantB509Primary,
		Secondary: vaillantB509Secondary,
		Data:      buildB509WriteSelector(addr, value),
	}
	response := protocol.Frame{
		Source:    target,
		Target:    source,
		Primary:   vaillantB509Primary,
		Secondary: vaillantB509Secondary,
		Data:      []byte{vaillantB509OpcodeWrite, byte(addr >> 8), byte(addr)},
	}
	return engine.PassiveClassifiedEvent{
		Kind:        engine.PassiveClassifiedEventTransaction,
		FrameType:   protocol.FrameTypeInitiatorTarget,
		Request:     request,
		Response:    response,
		HasRequest:  true,
		HasResponse: true,
		ObservedAt:  observedAt,
	}
}

type unknownSemanticWatchKey struct {
	canonical string
}

func waitForAdjudicatedEvent(
	t *testing.T,
	subscription *engine.AdjudicatedPassiveSubscription,
	timeout time.Duration,
	match func(engine.AdjudicatedPassiveEvent) bool,
) engine.AdjudicatedPassiveEvent {
	t.Helper()

	if subscription == nil {
		t.Fatal("waitForAdjudicatedEvent subscription = nil")
	}
	if match == nil {
		t.Fatal("waitForAdjudicatedEvent match = nil")
	}

	deadline := time.NewTimer(timeout)
	defer deadline.Stop()

	for {
		select {
		case <-deadline.C:
			t.Fatalf("timeout waiting for adjudicated event after %s", timeout)
		case event, ok := <-subscription.Events():
			if !ok {
				t.Fatal("subscription closed before matching adjudicated event")
			}
			if match(event) {
				return event
			}
		}
	}
}

func waitForShadowIneligible(
	t *testing.T,
	shadow *engine.ShadowCache,
	key engine.WatchKey,
	timeout time.Duration,
) engine.ShadowLookupResult {
	t.Helper()

	if shadow == nil || key == nil {
		t.Fatal("waitForShadowIneligible requires non-nil shadow and key")
	}

	deadline := time.Now().Add(timeout)
	for {
		lookup := shadow.Lookup(key, 10*time.Second)
		if lookup.Found &&
			(lookup.Entry.State == engine.ShadowEntryStateInvalidated ||
				lookup.Entry.State == engine.ShadowEntryStateTombstone) {
			return lookup
		}
		if time.Now().After(deadline) {
			t.Fatalf("timeout waiting for shadow invalidation on %q", key.Canonical())
		}
		time.Sleep(10 * time.Millisecond)
	}
}

type staticSemanticReadWatchObserver struct {
	observation engine.WatchObservation
}

func (observer staticSemanticReadWatchObserver) Observe(engine.WatchKey) engine.WatchObservation {
	return observer.observation
}
func (key unknownSemanticWatchKey) Canonical() string {
	return key.canonical
}

func (key unknownSemanticWatchKey) String() string {
	return key.canonical
}

func (key unknownSemanticWatchKey) Family() engine.WatchFamily {
	return engine.WatchFamily("UNKNOWN")
}
