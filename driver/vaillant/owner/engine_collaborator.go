package owner

import (
	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

// EngineCollaborator is the engine collaborator type of the Vaillant owner.
type EngineCollaborator struct{}

// ResponseClassForPassive implements the response class for passive step of the Vaillant owner.
func (EngineCollaborator) ResponseClassForPassive(event engine.PassiveClassifiedEvent, outcome engine.DedupOutcomeClass) engine.DedupResponseClass {
	return observeFirstResponseClassForPassiveEvent(event, outcome)
}

// ResponseClassForActive implements the response class for active step of the Vaillant owner.
func (EngineCollaborator) ResponseClassForActive(event protocol.BusEvent) engine.DedupResponseClass {
	return observeFirstResponseClassForActiveEvent(event)
}

// FamilyPolicy implements the family policy step of the Vaillant owner.
func (EngineCollaborator) FamilyPolicy(scope engine.ObserveFirstTrafficScope, request protocol.Frame, responseClass engine.DedupResponseClass, observation engine.WatchObservation, flags engine.ObserveFirstFeatureFlagView) engine.ObserveFirstFamilyPolicy {
	return observeFirstFamilyPolicy(scope, request, responseClass, observation, flags)
}

// PassiveWatchKeyFromEvent implements the passive watch key from event step of the Vaillant owner.
func (EngineCollaborator) PassiveWatchKeyFromEvent(event engine.PassiveClassifiedEvent) (engine.WatchKey, bool) {
	return PassiveWatchKeyFromEvent(event)
}

// PassiveWatchKeyFromFrame implements the passive watch key from frame step of the Vaillant owner.
func (EngineCollaborator) PassiveWatchKeyFromFrame(frame protocol.Frame) (engine.WatchKey, bool) {
	return PassiveWatchKeyFromFrame(frame)
}

// RequestIntentFromFrame implements the request intent from frame step of the Vaillant owner.
func (EngineCollaborator) RequestIntentFromFrame(frame protocol.Frame) engine.ObserveFirstRequestIntent {
	return observeFirstRequestIntentFromFrame(frame)
}

// IsScanProbeRaw implements the is scan probe raw step of the Vaillant owner.
func (EngineCollaborator) IsScanProbeRaw(raw []byte) bool {
	return isScanProbeRaw(raw)
}

// IsScanTimeout reports the donor pending-probe test: a parsed request whose
// primary and secondary bytes are 0x07 and 0x04.
func (EngineCollaborator) IsScanTimeout(request protocol.Frame) bool {
	return request.Primary == 0x07 && request.Secondary == 0x04
}

// CloneWatchKey implements the clone watch key step of the Vaillant owner.
func (EngineCollaborator) CloneWatchKey(key engine.WatchKey) (engine.WatchKey, error) {
	return cloneWatchKey(key)
}

// ValidateB555WatchKey implements the validate b555 watch key step of the Vaillant owner.
func (EngineCollaborator) ValidateB555WatchKey(key engine.WatchKey) error {
	return validateB555WatchKey(key)
}
