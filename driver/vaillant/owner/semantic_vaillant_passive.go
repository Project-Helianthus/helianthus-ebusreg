package owner

import (
	"context"
	"log"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
)

func (p *vaillantSemanticPoller) AttachPassiveShadowProducer(ctx context.Context, deduplicator *engine.ActivePassiveDeduplicator) error {
	if p == nil || deduplicator == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	subscription, err := deduplicator.Subscribe(
		"semantic-vaillant-shadow",
		passiveShadowSubscriberPriority,
		passiveShadowSubscriberBuffer,
	)
	if err != nil {
		return err
	}

	go func(subscription *engine.AdjudicatedPassiveSubscription) {
		for {
			closedUnexpectedly := false
		consume:
			for {
				select {
				case <-ctx.Done():
					subscription.Close()
					return
				case event, ok := <-subscription.Events():
					if !ok {
						closedUnexpectedly = ctx.Err() == nil
						break consume
					}
					p.handleAdjudicatedPassiveEvent(event)
				}
			}
			subscription.Close()
			if !closedUnexpectedly {
				return
			}
			for {
				if !waitForPassiveShadowRetry(ctx) {
					return
				}
				next, err := deduplicator.Subscribe(
					"semantic-vaillant-shadow",
					passiveShadowSubscriberPriority,
					passiveShadowSubscriberBuffer,
				)
				if err != nil {
					continue
				}
				subscription = next
				break
			}
		}
	}(subscription)

	return nil
}

func waitForPassiveShadowRetry(ctx context.Context) bool {
	delay := passiveShadowRetryDelay
	if delay <= 0 {
		delay = 100 * time.Millisecond
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func (p *vaillantSemanticPoller) handleAdjudicatedPassiveEvent(event engine.AdjudicatedPassiveEvent) {
	if p == nil || p.shadow == nil {
		return
	}
	if event.Disposition != engine.DedupDispositionUnmatchedThirdParty ||
		event.ObservabilityOnly ||
		event.LocalParticipantInbound ||
		event.MatchedActiveDuplicate {
		return
	}

	familyPolicy := event.Fingerprint.FamilyPolicy
	if !passiveShadowLaneEnabled(p.shadow.FeatureFlags(), familyPolicy) {
		return
	}

	key, ok := clonePassiveAdjudicatedWatchKey(event)
	if !ok || key == nil {
		return
	}
	runtime := p.bootstrapPassiveSharedWatchKey(key)

	switch familyPolicy.RequestIntent {
	case engine.ObserveFirstRequestIntentWrite:
		if !shouldInvalidatePassiveExternalWrite(familyPolicy) {
			return
		}
		invalidatedAt := event.Fingerprint.ObservedAt
		if invalidatedAt.IsZero() {
			invalidatedAt = p.now()
		}
		p.shadow.Invalidate(engine.ShadowInvalidation{
			Key:           key,
			Reason:        engine.ShadowInvalidationReasonExternalWrite,
			Source:        engine.ShadowInvalidationSourcePassive,
			InvalidatedAt: invalidatedAt,
		})
		return
	case engine.ObserveFirstRequestIntentRead:
	default:
		return
	}
	if event.SuppressShadow {
		return
	}

	if event.Fingerprint.ResponseClass != engine.DedupResponseValueBearing {
		return
	}
	value, ok := parsePassiveShadowPayload(event, key)
	if !ok || len(value) == 0 {
		return
	}

	observedAt := event.Fingerprint.ObservedAt
	if observedAt.IsZero() {
		observedAt = p.now()
	}
	result := p.shadow.Write(engine.ShadowWrite{
		Key:        key,
		Source:     engine.ShadowWriteSourcePassive,
		Confidence: engine.ShadowConfidenceHigh,
		Value:      value,
		ObservedAt: observedAt,
	})
	p.emitWatchDirectApplyEfficiency(runtime, observedAt, true, result.Accepted)
	if !result.Accepted {
		log.Printf("semantic_passive_shadow_write_rejected key=%q reason=%s", key.Canonical(), result.Reason)
		return
	}
}

func passiveShadowLaneEnabled(flags engine.ObserveFirstFeatureFlags, policy engine.ObserveFirstFamilyPolicy) bool {
	if !flags.ObserveFirstEnabled() {
		return false
	}

	switch policy.RequestIntent {
	case engine.ObserveFirstRequestIntentRead:
		if policy.DirectApplyPolicy == engine.ObserveFirstDirectApplyPolicyStateDefault {
			return flags.PassiveStateDirectApply()
		}
		return false
	case engine.ObserveFirstRequestIntentWrite:
		if !policy.UsesRuntimeExternalWritePolicy {
			return false
		}
		return policy.EffectiveExternalWritePolicy != engine.ObserveFirstExternalWritePolicyRecordOnly
	default:
		return false
	}
}

func shouldInvalidatePassiveExternalWrite(policy engine.ObserveFirstFamilyPolicy) bool {
	if !policy.UsesRuntimeExternalWritePolicy {
		return false
	}
	switch policy.EffectiveExternalWritePolicy {
	case engine.ObserveFirstExternalWritePolicyInvalidateOnly,
		engine.ObserveFirstExternalWritePolicyRecordAndInvalidate:
		return true
	default:
		return false
	}
}

func clonePassiveAdjudicatedWatchKey(event engine.AdjudicatedPassiveEvent) (engine.WatchKey, bool) {
	return cloneSemanticWatchKey(event.Fingerprint.SharedWatchKey)
}

func cloneSemanticWatchKey(key engine.WatchKey) (engine.WatchKey, bool) {
	switch typed := key.(type) {
	case engine.B509WatchKey:
		cloned := engine.NewB509WatchKey(typed.Target, typed.RegisterAddress)
		return cloned, true
	case *engine.B509WatchKey:
		if typed == nil {
			return nil, false
		}
		cloned := engine.NewB509WatchKey(typed.Target, typed.RegisterAddress)
		return cloned, true
	case engine.B524WatchKey:
		cloned := engine.NewB524WatchKey(typed.Target, typed.Opcode, typed.Group, typed.Instance, typed.RegisterAddress)
		return cloned, true
	case *engine.B524WatchKey:
		if typed == nil {
			return nil, false
		}
		cloned := engine.NewB524WatchKey(typed.Target, typed.Opcode, typed.Group, typed.Instance, typed.RegisterAddress)
		return cloned, true
	default:
		return nil, false
	}
}

func parsePassiveShadowPayload(event engine.AdjudicatedPassiveEvent, key engine.WatchKey) ([]byte, bool) {
	if !event.Event.HasResponse {
		return nil, false
	}
	switch typed := key.(type) {
	case engine.B509WatchKey:
		return parseB509ReadPayload(event.Event.Response.Data, typed.RegisterAddress)
	case *engine.B509WatchKey:
		if typed == nil {
			return nil, false
		}
		return parseB509ReadPayload(event.Event.Response.Data, typed.RegisterAddress)
	case engine.B524WatchKey:
		return parseB524ReadPayload(event.Event.Response.Data, typed.Opcode, typed.Group, typed.Instance, typed.RegisterAddress)
	case *engine.B524WatchKey:
		if typed == nil {
			return nil, false
		}
		return parseB524ReadPayload(event.Event.Response.Data, typed.Opcode, typed.Group, typed.Instance, typed.RegisterAddress)
	default:
		return nil, false
	}
}
