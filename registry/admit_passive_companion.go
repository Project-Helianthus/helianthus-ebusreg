package registry

import (
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

// AdmitPassiveCompanionWithCurrentQualifiedIdentityWitness atomically admits
// the canonical companion of source as a passively observed target when source
// currently has a direct, complete qualified-identity witness. The companion
// is derived from the protocol-owned canonical mapping; callers cannot assert
// an arbitrary relationship.
//
// This method deliberately validates the witness and commits the passive slot
// while holding one registry write lock. It does not register identity, create
// an alias, or turn topology evidence into identity authority.
func (r *DeviceRegistry) AdmitPassiveCompanionWithCurrentQualifiedIdentityWitness(source byte, observedAt time.Time) (companion byte, admitted bool) {
	if r == nil || observedAt.IsZero() {
		return 0, false
	}
	companion, ok := protocol.CompanionOfSource(source)
	if !ok {
		return 0, false
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.admitPassiveCompanionWithCurrentQualifiedIdentityWitnessLocked(source, companion, observedAt) {
		return 0, false
	}
	return companion, true
}

// admitPassiveCompanionWithCurrentQualifiedIdentityWitnessLocked performs the
// registry-owned commit for AdmitPassiveCompanionWithCurrentQualifiedIdentityWitness.
// The caller holds r.mu. Keeping validation and the slot update in this helper
// prevents a replacement, retirement, or conflict from interleaving between
// them.
func (r *DeviceRegistry) admitPassiveCompanionWithCurrentQualifiedIdentityWitnessLocked(source, companion byte, observedAt time.Time) bool {
	witness := r.qualifiedWitnesses[source].witness
	if !r.qualifiedIdentityWitnessCurrentLocked(source, witness) {
		return false
	}
	// Address roles are sticky once observed. A companion already established
	// as a source-side role cannot become the canonical target through this narrow
	// admission, so reject before changing any state.
	if slot := r.addressTable[companion]; slot != nil && slot.Role != SlotRoleUnknown && slot.Role != SlotRoleSlave {
		return false
	}

	slot := r.ensureAddressSlotLocked(companion)
	r.markSlotPassiveObservedLocked(slot, SlotRoleSlave, observedAt)
	if slot.Device != nil {
		r.syncEntryFacesLocked(slot.Device)
	}
	r.observationGeneration++
	return true
}
