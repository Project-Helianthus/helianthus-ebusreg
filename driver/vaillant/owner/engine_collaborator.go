package owner

import (
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

type EngineCollaborator struct{}

func (EngineCollaborator) ResponseClassForPassive(event engine.PassiveClassifiedEvent, outcome engine.DedupOutcomeClass) engine.DedupResponseClass {
	return observeFirstResponseClassForPassiveEvent(event, outcome)
}

func (EngineCollaborator) ResponseClassForActive(event protocol.BusEvent) engine.DedupResponseClass {
	return observeFirstResponseClassForActiveEvent(event)
}

func (EngineCollaborator) FamilyPolicy(scope engine.ObserveFirstTrafficScope, request protocol.Frame, responseClass engine.DedupResponseClass, observation engine.WatchObservation, flags engine.ObserveFirstFeatureFlagView) engine.ObserveFirstFamilyPolicy {
	return observeFirstFamilyPolicy(scope, request, responseClass, observation, flags)
}

func (EngineCollaborator) PassiveWatchKeyFromEvent(event engine.PassiveClassifiedEvent) (engine.WatchKey, bool) {
	return PassiveWatchKeyFromEvent(event)
}

func (EngineCollaborator) PassiveWatchKeyFromFrame(frame protocol.Frame) (engine.WatchKey, bool) {
	return PassiveWatchKeyFromFrame(frame)
}

func (EngineCollaborator) RequestIntentFromFrame(frame protocol.Frame) engine.ObserveFirstRequestIntent {
	return observeFirstRequestIntentFromFrame(frame)
}

func (EngineCollaborator) IsScanProbeRaw(raw []byte) bool {
	return isScanProbeRaw(raw)
}

// IsScanTimeout reports the donor pending-probe test: a parsed request whose
// primary and secondary bytes are 0x07 and 0x04.
func (EngineCollaborator) IsScanTimeout(request protocol.Frame) bool {
	return request.Primary == 0x07 && request.Secondary == 0x04
}

func (EngineCollaborator) CloneWatchKey(key engine.WatchKey) (engine.WatchKey, error) {
	return cloneWatchKey(key)
}

func (EngineCollaborator) ValidateB555WatchKey(key engine.WatchKey) error {
	return validateB555WatchKey(key)
}
