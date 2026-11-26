package owner

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

type semanticReadWatchRuntime struct {
	key           engine.WatchKey
	descriptor    engine.WatchDescriptor
	hasDescriptor bool
	maxAge        time.Duration
}

func (p *vaillantSemanticPoller) prepareSemanticReadWatch(key engine.WatchKey) time.Duration {
	return p.prepareSemanticReadWatchRuntime(key).maxAge
}

func (p *vaillantSemanticPoller) prepareSemanticReadWatchRuntime(key engine.WatchKey) semanticReadWatchRuntime {
	return p.resolveSemanticReadWatchRuntime(
		key,
		semanticReadWatchDescriptorDefault,
		[]engine.WatchActivationSource{engine.WatchActivationSourcePoller},
		semanticReadActivationSources,
	)
}

func (p *vaillantSemanticPoller) bootstrapPassiveSharedWatchKey(key engine.WatchKey) semanticReadWatchRuntime {
	return p.resolveSemanticReadWatchRuntime(
		key,
		semanticReadWatchDescriptorPassiveFallback,
		[]engine.WatchActivationSource{engine.WatchActivationSourceTooling},
		semanticReadPassiveActivationSources,
	)
}

func (p *vaillantSemanticPoller) resolveSemanticReadWatchRuntime(
	key engine.WatchKey,
	fallback func(engine.WatchKey) engine.WatchDescriptor,
	defaultSources []engine.WatchActivationSource,
	sourceResolver func([]engine.WatchActivationSource) []engine.WatchActivationSource,
) semanticReadWatchRuntime {
	runtime := semanticReadWatchRuntime{
		key: key,
	}
	if key == nil {
		return runtime
	}

	descriptor := fallback(key)
	hasDescriptor := false
	sources := defaultSources
	if sourceResolver == nil {
		sourceResolver = semanticReadActivationSources
	}
	if p != nil && p.watchObserver != nil {
		observation := p.watchObserver.Observe(key)
		if observation.HasDescriptor {
			descriptor = semanticReadWatchDescriptorFromObservation(key, observation.Descriptor)
			hasDescriptor = true
		}
		sources = sourceResolver(observation.Sources)
	}

	maxAge, err := descriptor.EffectiveFreshnessTTL(EngineCollaborator{})
	if err != nil {
		descriptor = fallback(key)
		maxAge, _ = descriptor.EffectiveFreshnessTTL(EngineCollaborator{})
		hasDescriptor = false
	}
	if maxAge < 0 {
		maxAge = 0
	}

	if p != nil && p.shadow != nil {
		if err := p.shadow.BootstrapRuntimeDescriptor(descriptor, sources...); err != nil {
			log.Printf("semantic_read_shadow_bootstrap_failed key=%q err=%v", key.Canonical(), err)
		}
	}

	runtime.descriptor = descriptor
	runtime.hasDescriptor = hasDescriptor
	runtime.maxAge = maxAge
	return runtime
}

func (p *vaillantSemanticPoller) emitWatchReadEfficiency(runtime semanticReadWatchRuntime, maxAge time.Duration, stats engine.SemanticReadExecutionStats) {
	if p == nil || p.watchEfficiency == nil || runtime.key == nil {
		return
	}
	observedAt := p.now()
	p.watchEfficiency.ObserveWatchRead(engine.WatchEfficiencyReadEvent{
		Key:           runtime.key,
		Descriptor:    runtime.descriptor,
		HasDescriptor: runtime.hasDescriptor,
		MaxAge:        maxAge,
		Stats:         stats,
		ObservedAt:    observedAt,
	})
}

func (p *vaillantSemanticPoller) emitWatchDirectApplyEfficiency(runtime semanticReadWatchRuntime, observedAt time.Time, candidateEvaluated bool, accepted bool) {
	if p == nil || p.watchEfficiency == nil || runtime.key == nil {
		return
	}
	if observedAt.IsZero() {
		observedAt = p.now()
	}
	p.watchEfficiency.ObserveWatchDirectApply(engine.WatchEfficiencyDirectApplyEvent{
		Key:                runtime.key,
		Descriptor:         runtime.descriptor,
		HasDescriptor:      runtime.hasDescriptor,
		ObservedAt:         observedAt,
		CandidateEvaluated: candidateEvaluated,
		Accepted:           accepted,
	})
}

func semanticReadWatchDescriptorDefault(key engine.WatchKey) engine.WatchDescriptor {
	return semanticReadWatchDescriptorWithDecoderID(key, "semantic.read.poller")
}

func semanticReadWatchDescriptorPassiveFallback(key engine.WatchKey) engine.WatchDescriptor {
	return semanticReadWatchDescriptorWithDecoderID(key, "semantic.read.passive_fallback")
}

func semanticReadWatchDescriptorWithDecoderID(key engine.WatchKey, decoderID string) engine.WatchDescriptor {
	semanticClass, freshnessProfile, directApplyPolicy := semanticReadWatchDescriptorPolicy(key)
	return engine.WatchDescriptor{
		Key:               key,
		SemanticClass:     semanticClass,
		FreshnessProfile:  freshnessProfile,
		DecoderID:         decoderID,
		CorrelationPolicy: engine.WatchCorrelationPolicyRequestResponse,
		DirectApplyPolicy: directApplyPolicy,
	}
}

func semanticReadWatchDescriptorPolicy(key engine.WatchKey) (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	switch typed := key.(type) {
	case engine.B524WatchKey:
		return semanticReadB524WatchDescriptorPolicy(typed)
	case *engine.B524WatchKey:
		if typed != nil {
			return semanticReadB524WatchDescriptorPolicy(*typed)
		}
	}
	return semanticReadStateFastDescriptorPolicy()
}

func semanticReadStateFastDescriptorPolicy() (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	return engine.WatchSemanticClassState, engine.WatchFreshnessProfileStateFast, engine.WatchDirectApplyPolicyStateDefault
}

func semanticReadStateSlowDescriptorPolicy() (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	return engine.WatchSemanticClassState, engine.WatchFreshnessProfileStateSlow, engine.WatchDirectApplyPolicyStateDefault
}

func semanticReadConfigDescriptorPolicy() (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	return engine.WatchSemanticClassConfig, engine.WatchFreshnessProfileConfig, engine.WatchDirectApplyPolicyConfigOptIn
}

func semanticReadDiscoveryDescriptorPolicy() (engine.WatchSemanticClass, engine.WatchFreshnessProfile, engine.WatchDirectApplyPolicy) {
	return engine.WatchSemanticClassDiscovery, engine.WatchFreshnessProfileDiscovery, engine.WatchDirectApplyPolicyNever
}

func semanticReadWatchDescriptorFromObservation(key engine.WatchKey, observed engine.WatchDescriptor) engine.WatchDescriptor {
	normalized := observed
	normalized.Key = key
	if normalized.SemanticClass == "" {
		normalized.SemanticClass = engine.WatchSemanticClassState
	}
	if normalized.FreshnessProfile == "" {
		normalized.FreshnessProfile = engine.WatchFreshnessProfileStateFast
	}
	if normalized.DecoderID == "" {
		normalized.DecoderID = "semantic.read.poller"
	}
	if normalized.CorrelationPolicy == "" {
		normalized.CorrelationPolicy = engine.WatchCorrelationPolicyRequestResponse
	}
	if normalized.DirectApplyPolicy == "" {
		normalized.DirectApplyPolicy = engine.WatchDirectApplyPolicyStateDefault
	}
	return normalized
}

func semanticReadActivationSources(sources []engine.WatchActivationSource) []engine.WatchActivationSource {
	if len(sources) == 0 {
		return []engine.WatchActivationSource{engine.WatchActivationSourcePoller}
	}

	out := make([]engine.WatchActivationSource, 0, len(sources)+1)
	seen := make(map[engine.WatchActivationSource]struct{}, len(sources)+1)
	for _, source := range sources {
		if source == "" {
			continue
		}
		if _, ok := seen[source]; ok {
			continue
		}
		seen[source] = struct{}{}
		out = append(out, source)
	}
	if _, ok := seen[engine.WatchActivationSourcePoller]; !ok {
		out = append(out, engine.WatchActivationSourcePoller)
	}
	return out
}

func semanticReadPassiveActivationSources(sources []engine.WatchActivationSource) []engine.WatchActivationSource {
	if len(sources) == 0 {
		return []engine.WatchActivationSource{engine.WatchActivationSourceTooling}
	}

	seen := make(map[engine.WatchActivationSource]struct{}, len(sources))
	out := make([]engine.WatchActivationSource, 0, len(sources))
	for _, source := range sources {
		if source == "" {
			continue
		}
		if _, ok := seen[source]; ok {
			continue
		}
		seen[source] = struct{}{}
		out = append(out, source)
	}
	if len(out) == 0 {
		return []engine.WatchActivationSource{engine.WatchActivationSourceTooling}
	}
	return out
}

func (p *vaillantSemanticPoller) readB509Value(ctx context.Context, target byte, addr uint16) ([]byte, bool) {
	return p.readB509ValueWithMaxAge(ctx, target, addr, -1)
}

func (p *vaillantSemanticPoller) readB509ValueLive(ctx context.Context, target byte, addr uint16) ([]byte, bool) {
	return p.readB509ValueWithMaxAge(ctx, target, addr, 0)
}

func (p *vaillantSemanticPoller) readB509ValueWithMaxAge(ctx context.Context, target byte, addr uint16, maxAgeOverride time.Duration) ([]byte, bool) {
	if p == nil || (p.bus == nil && p.sendFrameFn == nil) || target == 0 {
		return nil, false
	}
	if ctx == nil {
		ctx = context.Background()
	}

	p.mu.Lock()
	source := p.source
	timeout := p.requestTimeout
	p.mu.Unlock()
	if timeout <= 0 {
		timeout = 2 * time.Second
	}

	watchKey := engine.NewB509WatchKey(target, addr)
	watchRuntime := p.prepareSemanticReadWatchRuntime(watchKey)
	maxAge := watchRuntime.maxAge
	if maxAgeOverride >= 0 {
		maxAge = maxAgeOverride
	}
	value, stats, err := p.scheduler.GetWatchWithStats(ctx, watchKey, maxAge, func(ctx context.Context) ([]byte, error) {
		var lastErr error
		for attempt := 0; attempt < 3; attempt++ {
			p.readMu.Lock()
			reqCtx, cancel := context.WithTimeout(ctx, timeout)
			request := protocol.Frame{
				FrameType: protocol.FrameTypeInitiatorTarget,
				Source:    source,
				Target:    target,
				Primary:   vaillantB509Primary,
				Secondary: vaillantB509Secondary,
				Data:      buildB509ReadSelector(addr),
			}
			response, err := p.sendSemanticFrame(reqCtx, request)
			cancel()
			p.readMu.Unlock()
			if addr == boiler_b509_flow_temperature {
				p.retainBoilerThermalRead(request, response, err)
			}

			if err != nil {
				lastErr = err
			} else if response == nil {
				lastErr = fmt.Errorf("b509 read returned nil response")
			} else if payload, ok := parseB509ReadPayload(response.Data, addr); ok {
				return payload, nil
			} else {
				lastErr = fmt.Errorf("b509 read failed: target=0x%02x addr=0x%04x", target, addr)
			}

			if attempt < 2 && !p.pause(ctx, 75*time.Millisecond) {
				return nil, ctx.Err()
			}
		}
		if lastErr == nil {
			lastErr = fmt.Errorf("b509 read failed")
		}
		return nil, lastErr
	})
	p.emitWatchReadEfficiency(watchRuntime, maxAge, stats)
	if stats.ActiveFetchSucceeded {
		p.pause(ctx, 25*time.Millisecond)
	}
	if err != nil {
		return nil, false
	}
	return value, true
}

func (p *vaillantSemanticPoller) writeB509Value(ctx context.Context, target byte, addr uint16, value []byte) error {
	if p == nil || (p.bus == nil && p.sendFrameFn == nil) {
		return fmt.Errorf("b509 write unavailable")
	}
	if target == 0 {
		return fmt.Errorf("b509 write target is zero")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	p.mu.Lock()
	source := p.source
	timeout := p.requestTimeout
	p.mu.Unlock()
	if timeout <= 0 {
		timeout = 2 * time.Second
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		p.readMu.Lock()
		reqCtx, cancel := context.WithTimeout(ctx, timeout)
		request := protocol.Frame{
			FrameType: protocol.FrameTypeInitiatorTarget,
			Source:    source,
			Target:    target,
			Primary:   vaillantB509Primary,
			Secondary: vaillantB509Secondary,
			Data:      buildB509WriteSelector(addr, value),
		}
		response, err := p.sendSemanticFrame(reqCtx, request)
		cancel()
		p.readMu.Unlock()

		if err != nil {
			lastErr = err
		} else if response == nil {
			lastErr = fmt.Errorf("b509 write returned nil response")
		} else if parseB509WriteAck(response.Data, addr) {
			return nil
		} else {
			lastErr = fmt.Errorf("b509 write ack invalid: target=0x%02x addr=0x%04x", target, addr)
		}

		if attempt < 2 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(75 * time.Millisecond):
			}
		}
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("b509 write failed")
	}
	return lastErr
}

func (p *vaillantSemanticPoller) sendSemanticFrame(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error) {
	if p == nil {
		return nil, fmt.Errorf("semantic poller unavailable")
	}
	if p.sendFrameFn != nil {
		return p.sendFrameFn(ctx, frame)
	}
	if p.bus == nil {
		return nil, fmt.Errorf("semantic bus unavailable")
	}
	return p.bus.Send(ctx, frame)
}

func (p *vaillantSemanticPoller) readB509Data2b(ctx context.Context, target byte, addr uint16) *float64 {
	raw, ok := p.readB509Value(ctx, target, addr)
	if !ok {
		return nil
	}
	value, ok := decodeData2b(raw)
	if !ok {
		return nil
	}
	return &value
}

func (p *vaillantSemanticPoller) readB509Data2c(ctx context.Context, target byte, addr uint16) *float64 {
	raw, ok := p.readB509Value(ctx, target, addr)
	if !ok {
		return nil
	}
	value, ok := decodeData2c(raw)
	if !ok {
		return nil
	}
	return &value
}

func (p *vaillantSemanticPoller) readB509OnOff(ctx context.Context, target byte, addr uint16) *bool {
	raw, ok := p.readB509Value(ctx, target, addr)
	if !ok {
		return nil
	}
	value, ok := decodeOnOff(raw)
	if !ok {
		return nil
	}
	return &value
}

func (p *vaillantSemanticPoller) readB509UINInt(ctx context.Context, target byte, addr uint16) *int {
	raw, ok := p.readB509Value(ctx, target, addr)
	if !ok {
		return nil
	}
	value, ok := decodeUIN(raw)
	if !ok {
		return nil
	}
	out := int(value)
	return &out
}

func (p *vaillantSemanticPoller) readB509SINScaled(ctx context.Context, target byte, addr uint16, divisor float64) *float64 {
	raw, ok := p.readB509Value(ctx, target, addr)
	if !ok {
		return nil
	}
	value, ok := decodeSIN(raw)
	if !ok {
		return nil
	}
	out := float64(value)
	if divisor != 0 {
		out /= divisor
	}
	return &out
}

func (p *vaillantSemanticPoller) readB509UnsignedHundredths(ctx context.Context, target byte, addr uint16) *float64 {
	raw, ok := p.readB509Value(ctx, target, addr)
	if !ok {
		return nil
	}
	value, ok := decodeUnsignedHundredths(raw)
	if !ok {
		return nil
	}
	return &value
}

func (p *vaillantSemanticPoller) readB509Percent0(ctx context.Context, target byte, addr uint16) *float64 {
	raw, ok := p.readB509Value(ctx, target, addr)
	if !ok {
		return nil
	}
	value, ok := decodePercent0(raw)
	if !ok {
		return nil
	}
	out := float64(value)
	return &out
}

func (p *vaillantSemanticPoller) readB509Hoursum2(ctx context.Context, target byte, addr uint16) *float64 {
	raw, ok := p.readB509Value(ctx, target, addr)
	if !ok {
		return nil
	}
	value, ok := decodeHoursum2(raw)
	if !ok {
		return nil
	}
	return &value
}
