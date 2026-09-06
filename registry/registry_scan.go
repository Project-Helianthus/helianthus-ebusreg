package registry

import "time"

func (r *DeviceRegistry) ensureAddressSlotLocked(address byte) *AddressSlot {
	slot := r.addressTable[address]
	if slot == nil {
		slot = &AddressSlot{Addr: address}
		r.addressTable[address] = slot
	}
	return slot
}

// MarkSlotPassiveObserved updates an AddressSlot for an address that was
// passively observed by the gateway (e.g. by AddressTableInserter on
// positive ACK following a complete request). Writes Role / Discovery
// Source / VerificationState / FirstObservedAt / LastObservedAt under the
// registry write lock so concurrent readers via LookupSlot / Lookup do
// not see torn state.
//
// This API replaces direct *AddressSlot field mutation by the gateway
// inserter, which was racy with other readers (Codex P2 follow-up from
// PR #565). Idempotent: re-marking the same slot only advances
// VerificationState monotonically. DiscoverySource records the first
// non-unknown native admission path and is retained thereafter.
//
// SCOPE: this API only mutates the AddressSlot. It does NOT attach
// the slot to a device entry. To plant a NEW passively-observed
// address with identity attached AND label it correctly in a single
// critical section, use RegisterPassiveObserved (which composes
// registerLocked + this primitive). Calling Register followed by
// MarkSlotPassiveObserved is not a substitute for the atomic admission API:
// a preceding Register records an active origin, which subsequent passive
// observation correctly retains.
func (r *DeviceRegistry) MarkSlotPassiveObserved(address byte, role SlotRole, observedAt time.Time) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	slot := r.ensureAddressSlotLocked(address)
	r.markSlotPassiveObservedLocked(slot, role, observedAt)
	// Phase C M-C6a: refresh entry.Faces so AddressByRole sees the
	// updated SlotRole. Without this sync, MarkSlotPassiveObserved
	// would leave Faces stale and AddressByRole(SlotRoleSlave) on a
	// just-passively-observed slot would return (0, false).
	if slot.Device != nil {
		r.syncEntryFacesLocked(slot.Device)
	}
	r.observationGeneration++
}

// markSlotPassiveObservedLocked is the shared slot-stamping primitive
// used by both MarkSlotPassiveObserved and RegisterPassiveObserved.
// Caller MUST hold r.mu and is responsible for any subsequent
// syncEntryFacesLocked call. Centralising the stamping rules here
// prevents drift between the two public entry points (mirrors the
// markSlotStaticSeedLocked design from P3.5).
func (r *DeviceRegistry) markSlotPassiveObservedLocked(slot *AddressSlot, role SlotRole, observedAt time.Time) {
	recordDiscoverySource(slot, DiscoverySourcePassiveObserved)
	if slot.VerificationState < VerificationStateCorroborated {
		slot.VerificationState = VerificationStateCorroborated
	}
	if role != SlotRoleUnknown && slot.Role == SlotRoleUnknown {
		slot.Role = role
	}
	if slot.FirstObservedAt.IsZero() && !observedAt.IsZero() {
		slot.FirstObservedAt = observedAt
	}
	if !observedAt.IsZero() {
		slot.LastObservedAt = observedAt
	}
}

// RegisterPassiveObserved plants identity for an address newly observed
// on the wire by the gateway's passive inserter. Mirrors Register's
// identity-merge behaviour but stamps the AddressSlot with
// DiscoverySourcePassiveObserved / VerificationStateCorroborated so
// the observability surface (`/metrics`, MCP `bus.summary.get`,
// address-table snapshots) correctly shows the slot's provenance as
// passive observation rather than active confirmation.
//
// P8 fix: previously the gateway inserter called Register (which
// stamps ActiveConfirmed/IdentityConfirmed) followed by
// MarkSlotPassiveObserved. That ordering recorded active discovery as the
// original source, so passively observed slots were misreported as active.
// RegisterPassiveObserved performs the identity-merge AND the
// passive-label stamping atomically under a single lock acquisition,
// avoiding the misorder.
//
// Later observations retain PassiveObserved as this face's original native
// discovery source. They may advance VerificationState independently, so a
// passively observed face can become IdentityConfirmed after active evidence.
//
// Single lock acquisition — composes registerLocked, then the shared
// passive-observation primitive, then syncEntryFacesLocked.
func (r *DeviceRegistry) RegisterPassiveObserved(info DeviceInfo, role SlotRole, observedAt time.Time) DeviceEntry {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry, retainedConflict := r.registerLocked(info)
	if retainedConflict {
		return entry
	}
	slot := r.ensureAddressSlotLocked(info.Address)
	slot.Device = entry
	r.markSlotPassiveObservedLocked(slot, role, observedAt)
	r.syncEntryFacesLocked(entry)
	r.observationGeneration++
	r.reconcileQualifiedIdentityWitnessesLocked()
	return entry
}

func (r *DeviceRegistry) observeAddressSlotLocked(address byte, entry *deviceEntry, source DiscoverySource, state VerificationState) {
	now := time.Now()
	slot := r.ensureAddressSlotLocked(address)
	slot.Device = entry
	recordDiscoverySource(slot, source)
	if slot.VerificationState < state {
		slot.VerificationState = state
	}
	if slot.FirstObservedAt.IsZero() {
		slot.FirstObservedAt = now
	}
	slot.LastObservedAt = now
}

// recordDiscoverySource preserves the first non-unknown native admission
// source for an address face. Verification confidence is tracked separately
// by VerificationState and may advance after the original source is recorded.
// Caller must hold r.mu.
func recordDiscoverySource(slot *AddressSlot, source DiscoverySource) {
	if slot != nil && slot.DiscoverySource == DiscoverySourceUnknown && source != DiscoverySourceUnknown {
		slot.DiscoverySource = source
	}
}

func (r *DeviceRegistry) syncEntryFacesLocked(entry *deviceEntry) {
	if entry == nil {
		return
	}
	faces := make([]BusFace, 0, len(entry.addresses))
	for _, address := range entry.addresses {
		slot := r.ensureAddressSlotLocked(address)
		if slot.Device == nil {
			slot.Device = entry
		}
		faces = append(faces, BusFace{
			Addr:              address,
			Role:              slot.Role,
			DiscoverySource:   slot.DiscoverySource,
			VerificationState: slot.VerificationState,
		})
	}
	entry.Faces = faces
}

type deviceEntry struct {
	primaryAddress byte
	addresses      []byte
	physical       physicalIdentity
	identityKey    string
	// identityKeyAliases records legacy additional r.identity bindings so the
	// shared cleanup path can retire them. New entries retain only their current
	// qualified triple: explicit address aliases and LKG are not independent
	// identity authority.
	identityKeyAliases []string
	info               DeviceInfo
	planes             []Plane
	projections        []Projection
	index              CanonicalIndex
	indexErr           error
	Faces              []BusFace
}

// PrimaryDisplayAddress returns a representative address for log/UI
// display. Returns the canonical primary if set, otherwise the
// originally registered info.Address. Use this for log lines,
// MCP/GraphQL device.address fields, UI labels — anywhere the value
// is shown to humans rather than written to the wire. For wire
// routing, use AddressByRole(SlotRole) which is class-aware.
//
// Phase C M-C6c: replaces deviceEntry.Address(), whose name conflated
// display and routing semantics for aliased canonical pairs.
