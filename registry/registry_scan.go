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
// DiscoverySource / VerificationState monotonically (the slot retains
// the higher of the existing and new value, matching
// observeAddressSlotLocked's monotonic semantics).
//
// SCOPE: this API only mutates the AddressSlot. It does NOT attach
// the slot to a device entry. To plant a NEW passively-observed
// address with identity attached AND label it correctly in a single
// critical section, use RegisterPassiveObserved (which composes
// registerLocked + this primitive). Calling Register followed by
// MarkSlotPassiveObserved produces a label-misorder because Register
// stamps DiscoverySourceActiveConfirmed and the monotonic guard then
// refuses to downgrade — that was the P8 bug fixed in
// RegisterPassiveObserved.
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
	if slot.DiscoverySource < DiscoverySourcePassiveObserved {
		slot.DiscoverySource = DiscoverySourcePassiveObserved
	}
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
// MarkSlotPassiveObserved. The monotonic ladder
// (PassiveObserved < ActiveConfirmed) made the second call a no-op,
// so passively-observed slots were misreported as `active_confirmed`.
// RegisterPassiveObserved performs the identity-merge AND the
// passive-label stamping atomically under a single lock acquisition,
// avoiding the misorder.
//
// Subsequent label progression (after RegisterPassiveObserved):
//   - An active confirmation (e.g. directed scan) DOES advance the
//     DiscoverySource to ActiveConfirmed (PassiveObserved <
//     ActiveConfirmed) AND VerificationState to IdentityConfirmed.
//   - A static-seed mark on a passively-observed slot DOES advance
//     DiscoverySource to StaticSeed (PassiveObserved < StaticSeed) —
//     pre-known taxonomy outranks wire-only inference.
//
// Single lock acquisition — composes registerLocked, then the shared
// passive-observation primitive, then syncEntryFacesLocked.
func (r *DeviceRegistry) RegisterPassiveObserved(info DeviceInfo, role SlotRole, observedAt time.Time) DeviceEntry {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry := r.registerLocked(info)
	slot := r.ensureAddressSlotLocked(info.Address)
	slot.Device = entry
	r.markSlotPassiveObservedLocked(slot, role, observedAt)
	r.syncEntryFacesLocked(entry)
	r.observationGeneration++
	return entry
}

func (r *DeviceRegistry) observeAddressSlotLocked(address byte, entry *deviceEntry, source DiscoverySource, state VerificationState) {
	now := time.Now()
	slot := r.ensureAddressSlotLocked(address)
	slot.Device = entry
	if slot.DiscoverySource < source {
		slot.DiscoverySource = source
	}
	if slot.VerificationState < state {
		slot.VerificationState = state
	}
	if slot.FirstObservedAt.IsZero() {
		slot.FirstObservedAt = now
	}
	slot.LastObservedAt = now
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
	// identityKeyAliases tracks ADDITIONAL r.identity keys that
	// resolve to this entry beyond its own identityKey. Populated by
	// AliasAddresses when canonical and removed-secondary had distinct
	// identity keys (e.g. canonical=MAC-keyed, secondary=serial-keyed)
	// and the secondary's key is re-pointed at canonical instead of
	// being deleted. detachAddressLocked iterates this slice to clean
	// up r.identity bindings when the merged entry is removed,
	// preventing orphan keys from resolving to a removed *deviceEntry.
	// (Codex P2 round-4 finding 2026-05-08 on PR #136.)
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
