package owner

import (
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

func observeFirstResponseClassForPassiveEvent(event engine.PassiveClassifiedEvent, outcome engine.DedupOutcomeClass) engine.DedupResponseClass {
	return observeFirstResponseClass(
		event.Request,
		event.FrameType,
		outcome,
		event.Response,
		event.HasResponse,
	)
}

func observeFirstResponseClassForActiveEvent(event protocol.BusEvent) engine.DedupResponseClass {
	return observeFirstResponseClass(
		event.Request,
		event.FrameType,
		engine.DedupOutcomeSuccess,
		event.Response,
		event.HasResponse,
	)
}

func observeFirstResponseClass(request protocol.Frame, frameType protocol.FrameType, outcome engine.DedupOutcomeClass, response protocol.Frame, hasResponse bool) engine.DedupResponseClass {
	if outcome != engine.DedupOutcomeSuccess {
		return engine.DedupResponseErrorOrAmbiguous
	}
	if !hasResponse {
		if frameType == protocol.FrameTypeInitiatorInitiator {
			return engine.DedupResponseACKOnly
		}
		return engine.DedupResponseErrorOrAmbiguous
	}

	family := observeFirstFamilyFromFrame(request)
	intent := observeFirstRequestIntentFromFrame(request)

	switch family {
	case engine.ObserveFirstFamilyB509:
		return observeFirstB509ResponseClass(request, intent, response.Data)
	case engine.ObserveFirstFamilyB524:
		return observeFirstB524ResponseClass(request, intent, response.Data)
	default:
		return observeFirstGenericResponseClass(response.Data)
	}
}

func observeFirstFamilyPolicy(scope engine.ObserveFirstTrafficScope, request protocol.Frame, responseClass engine.DedupResponseClass, observation engine.WatchObservation, flags engine.ObserveFirstFeatureFlagView) engine.ObserveFirstFamilyPolicy {
	family := observeFirstFamilyFromFrame(request)
	intent := observeFirstRequestIntentFromFrame(request)

	policy := engine.ObserveFirstFamilyPolicy{
		Family:            family,
		RequestIntent:     intent,
		ResponseClass:     responseClass,
		CorrelationPolicy: engine.WatchCorrelationPolicyRecordOnly,
		DirectApplyPolicy: engine.ObserveFirstDirectApplyPolicyNever,
	}

	if scope == engine.ObserveFirstTrafficScopePassive && intent == engine.ObserveFirstRequestIntentWrite {
		policy.UsesRuntimeExternalWritePolicy = true
		policy.EffectiveExternalWritePolicy = engine.NormalizeObserveFirstFeatureFlagsFromView(flags).ExternalWritePolicy()
	}

	switch family {
	case engine.ObserveFirstFamilyB516:
		policy.CorrelationPolicy = engine.WatchCorrelationPolicyBroadcastSelector
		if responseClass == engine.DedupResponseValueBearing {
			policy.DirectApplyPolicy = engine.ObserveFirstDirectApplyPolicyEnergyMergeOnly
		}
	case engine.ObserveFirstFamilyB524:
		policy.CorrelationPolicy = engine.WatchCorrelationPolicyRequestResponse
		if observeFirstB524CorrelatedStateRead(request, observation) && responseClass == engine.DedupResponseValueBearing {
			policy.DirectApplyPolicy = engine.ObserveFirstDirectApplyPolicyStateDefault
		}
	case engine.ObserveFirstFamilyB509:
		policy.CorrelationPolicy = engine.WatchCorrelationPolicyRequestResponse
		if observeFirstB509CorrelatedStateRead(request, observation) && responseClass == engine.DedupResponseValueBearing {
			policy.DirectApplyPolicy = engine.ObserveFirstDirectApplyPolicyStateDefault
		}
	case engine.ObserveFirstFamilyB555:
		policy.CorrelationPolicy = engine.WatchCorrelationPolicyRecordInvalidate
	default:
		if scope == engine.ObserveFirstTrafficScopePassive && intent == engine.ObserveFirstRequestIntentBroadcast {
			policy.CorrelationPolicy = engine.WatchCorrelationPolicyBroadcastSelector
		}
	}

	return policy
}

func observeFirstFamilyFromFrame(frame protocol.Frame) engine.ObserveFirstFamily {
	if frame.Primary != 0xB5 {
		return engine.ObserveFirstFamilyOther
	}
	switch frame.Secondary {
	case 0x09:
		return engine.ObserveFirstFamilyB509
	case 0x16:
		return engine.ObserveFirstFamilyB516
	case 0x24:
		return engine.ObserveFirstFamilyB524
	case 0x55:
		return engine.ObserveFirstFamilyB555
	default:
		return engine.ObserveFirstFamilyOther
	}
}

func observeFirstRequestIntentFromFrame(frame protocol.Frame) engine.ObserveFirstRequestIntent {
	switch observeFirstFamilyFromFrame(frame) {
	case engine.ObserveFirstFamilyB509:
		if len(frame.Data) == 0 {
			return engine.ObserveFirstRequestIntentUnknown
		}
		switch frame.Data[0] {
		case 0x0D, 0x29:
			return engine.ObserveFirstRequestIntentRead
		case 0x0E:
			return engine.ObserveFirstRequestIntentWrite
		default:
			return engine.ObserveFirstRequestIntentUnknown
		}
	case engine.ObserveFirstFamilyB516:
		return engine.ObserveFirstRequestIntentBroadcast
	case engine.ObserveFirstFamilyB524:
		if len(frame.Data) == 0 {
			return engine.ObserveFirstRequestIntentUnknown
		}
		switch frame.Data[0] {
		case 0x02, 0x06:
			if len(frame.Data) >= 2 && frame.Data[1] == 0x00 {
				return engine.ObserveFirstRequestIntentRead
			}
			if len(frame.Data) >= 2 {
				return engine.ObserveFirstRequestIntentWrite
			}
			return engine.ObserveFirstRequestIntentUnknown
		case 0x03:
			return engine.ObserveFirstRequestIntentRead
		default:
			return engine.ObserveFirstRequestIntentUnknown
		}
	case engine.ObserveFirstFamilyB555:
		if len(frame.Data) == 0 {
			return engine.ObserveFirstRequestIntentUnknown
		}
		switch frame.Data[0] {
		case 0xA3, 0xA4, 0xA5:
			return engine.ObserveFirstRequestIntentRead
		case 0xA6:
			return engine.ObserveFirstRequestIntentWrite
		default:
			return engine.ObserveFirstRequestIntentUnknown
		}
	default:
		return engine.ObserveFirstRequestIntentUnknown
	}
}

func observeFirstGenericResponseClass(payload []byte) engine.DedupResponseClass {
	if len(payload) == 0 {
		return engine.DedupResponseHeaderOnly
	}
	return engine.DedupResponseValueBearing
}

func observeFirstB509ResponseClass(request protocol.Frame, intent engine.ObserveFirstRequestIntent, payload []byte) engine.DedupResponseClass {
	if len(payload) == 0 {
		return engine.DedupResponseHeaderOnly
	}
	if intent == engine.ObserveFirstRequestIntentWrite {
		return engine.DedupResponseHeaderOnly
	}
	if intent != engine.ObserveFirstRequestIntentRead {
		return observeFirstGenericResponseClass(payload)
	}

	addr, ok := observeFirstExpectedB509Address(request)
	if !ok {
		return observeFirstGenericResponseClass(payload)
	}

	addrHi := byte(addr >> 8)
	addrLo := byte(addr)

	if len(payload) >= 3 && (payload[0] == 0x0D || payload[0] == 0x29) && payload[1] == addrHi && payload[2] == addrLo {
		if len(payload) == 3 {
			return engine.DedupResponseHeaderOnly
		}
		return engine.DedupResponseValueBearing
	}
	if len(payload) >= 2 && payload[0] == addrHi && payload[1] == addrLo {
		if len(payload) == 2 {
			return engine.DedupResponseHeaderOnly
		}
		return engine.DedupResponseValueBearing
	}
	return observeFirstGenericResponseClass(payload)
}

func observeFirstB524ResponseClass(request protocol.Frame, intent engine.ObserveFirstRequestIntent, payload []byte) engine.DedupResponseClass {
	if len(payload) == 0 {
		return engine.DedupResponseHeaderOnly
	}
	if intent != engine.ObserveFirstRequestIntentRead {
		if intent == engine.ObserveFirstRequestIntentWrite {
			return engine.DedupResponseHeaderOnly
		}
		return observeFirstGenericResponseClass(payload)
	}
	if observeFirstIsB524TimerRead(request) {
		return observeFirstGenericResponseClass(payload)
	}

	opcode, group, instance, addr, ok := observeFirstExpectedB524ReadSelector(request)
	if !ok {
		return observeFirstGenericResponseClass(payload)
	}
	if len(payload) == 1 && payload[0] == 0x00 {
		return engine.DedupResponseErrorOrAmbiguous
	}
	if len(payload) < 4 {
		return engine.DedupResponseErrorOrAmbiguous
	}
	if len(payload) >= 5 {
		replyInstance := payload[1]
		replyGroup := payload[2]
		replyAddr := uint16(payload[3]) | uint16(payload[4])<<8
		if replyGroup == group && replyAddr == addr {
			if !observeFirstMatchesB524ReplyInstance(replyInstance, instance) {
				return engine.DedupResponseErrorOrAmbiguous
			}
			if len(payload) == 5 {
				return engine.DedupResponseHeaderOnly
			}
			return engine.DedupResponseValueBearing
		}
	}

	replyGroup := payload[1]
	replyAddr := uint16(payload[2]) | uint16(payload[3])<<8
	if replyGroup != group || replyAddr != addr {
		_ = opcode
		return engine.DedupResponseErrorOrAmbiguous
	}
	if len(payload) == 4 {
		return engine.DedupResponseHeaderOnly
	}
	return engine.DedupResponseValueBearing
}

func observeFirstExpectedB509Address(frame protocol.Frame) (uint16, bool) {
	if len(frame.Data) < 3 {
		return 0, false
	}
	switch frame.Data[0] {
	case 0x0D, 0x29, 0x0E:
		return uint16(frame.Data[1])<<8 | uint16(frame.Data[2]), true
	default:
		return 0, false
	}
}

func observeFirstExpectedB524ReadSelector(frame protocol.Frame) (opcode, group, instance byte, addr uint16, ok bool) {
	if len(frame.Data) < 6 {
		return 0, 0, 0, 0, false
	}
	if frame.Data[0] != 0x02 && frame.Data[0] != 0x06 {
		return 0, 0, 0, 0, false
	}
	if frame.Data[1] != 0x00 {
		return 0, 0, 0, 0, false
	}
	return frame.Data[0], frame.Data[2], frame.Data[3], uint16(frame.Data[4]) | uint16(frame.Data[5])<<8, true
}

func observeFirstB524CorrelatedStateRead(frame protocol.Frame, observation engine.WatchObservation) bool {
	opcode, group, instance, addr, ok := observeFirstExpectedB524ReadSelector(frame)
	if !ok {
		return false
	}
	key := engine.NewB524WatchKey(frame.Target, opcode, group, instance, addr)
	return observeFirstMatchesActiveStateDefaultDescriptor(key, observation)
}

func observeFirstB509CorrelatedStateRead(frame protocol.Frame, observation engine.WatchObservation) bool {
	addr, ok := observeFirstExpectedB509Address(frame)
	if !ok || observeFirstRequestIntentFromFrame(frame) != engine.ObserveFirstRequestIntentRead {
		return false
	}
	key := engine.NewB509WatchKey(frame.Target, addr)
	return observeFirstMatchesActiveStateDefaultDescriptor(key, observation)
}

func observeFirstMatchesActiveStateDefaultDescriptor(key engine.WatchKey, observation engine.WatchObservation) bool {
	if key == nil {
		return false
	}
	if !observation.HasDescriptor || observation.State != engine.WatchObservationStateActive {
		return false
	}
	descriptor := observation.Descriptor
	if descriptor.Key == nil || descriptor.Key.Canonical() != key.Canonical() {
		return false
	}
	if descriptor.SemanticClass != engine.WatchSemanticClassState {
		return false
	}
	if descriptor.CorrelationPolicy != engine.WatchCorrelationPolicyRequestResponse {
		return false
	}
	return descriptor.DirectApplyPolicy == engine.WatchDirectApplyPolicyStateDefault
}

func observeFirstIsB524TimerRead(frame protocol.Frame) bool {
	return len(frame.Data) > 0 && frame.Data[0] == 0x03
}

func observeFirstMatchesB524ReplyInstance(replyInstance, requestedInstance byte) bool {
	if replyInstance == requestedInstance {
		return true
	}
	if requestedInstance < 0xFF && replyInstance == requestedInstance+1 {
		return true
	}
	return false
}

// isScanProbeRaw checks raw request bytes for a 0704 device identity query
// without requiring a full frame parse.  Used to reclassify collision-garbled
// scan probes as scan_collision instead of corrupted_request.
func isScanProbeRaw(raw []byte) bool {
	return len(raw) >= 4 && raw[2] == 0x07 && raw[3] == 0x04
}
