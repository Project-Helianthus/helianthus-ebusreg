package registry

import "time"

func (r *DeviceRegistry) Register(info DeviceInfo) DeviceEntry {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry := r.registerLocked(info)
	r.observeAddressSlotLocked(info.Address, entry, DiscoverySourceActiveConfirmed, VerificationStateIdentityConfirmed)
	r.syncEntryFacesLocked(entry)
	r.observationGeneration++
	return entry
}

// registerLocked performs the core identity-merge / planes / projections /
// index work for an incoming DeviceInfo without stamping a discovery
// label on the AddressSlot. It is the primitive shared by the
// active-discovery path (Register, which stamps
// DiscoverySourceActiveConfirmed / VerificationStateIdentityConfirmed)
// and the static-seed path (RegisterStaticSeed, which stamps
// DiscoverySourceStaticSeed / VerificationStateCandidate).
//
// Lock contract: the caller MUST hold r.mu. registerLocked does NOT
// acquire the lock itself. This mirrors the file's existing *Locked
// suffix convention (observeAddressSlotLocked, ensureAddressSlotLocked,
// syncEntryFacesLocked, lookupCompatibleBySignatureLocked,
// detachAddressLocked).
//
// Both callers are responsible for stamping the AddressSlot
// (observeAddressSlotLocked) and refreshing entry.Faces
// (syncEntryFacesLocked) after the call.
func (r *DeviceRegistry) registerLocked(info DeviceInfo) *deviceEntry {
	physical := canonicalPhysicalIdentity(info)
	identityKey := physical.key()
	planes := make([]Plane, 0)
	matched := make([]PlaneProvider, 0, len(r.providers))

	existingByAddress := r.entries[info.Address]
	incomingHasStableIdentity := normalizeIdentityPart(info.SerialNumber) != "" || normalizeIdentityPart(info.MacAddress) != ""
	if !incomingHasStableIdentity && existingByAddress != nil && len(existingByAddress.addresses) > 1 && existingByAddress.identityKey != "" {
		identityKey = existingByAddress.identityKey
	}
	if identityKey == "" && existingByAddress != nil {
		identityKey = existingByAddress.identityKey
	}

	existingByIdentity := (*deviceEntry)(nil)
	if identityKey != "" {
		existingByIdentity = r.identity[identityKey]
	}

	entry := existingByIdentity
	if entry == nil {
		entry = existingByAddress
	}
	if entry == existingByAddress && existingByIdentity == nil && existingByAddress != nil && (!canMergeIdentity(info, existingByAddress.info) || hasConflictingModelSignature(info, existingByAddress.info)) {
		entry = nil
	}
	if entry == nil {
		if fallback, ok := r.lookupCompatibleBySignatureLocked(info); ok {
			entry = fallback
		}
	}

	if existingByAddress != nil && existingByAddress != entry {
		r.detachAddressLocked(existingByAddress, info.Address)
	}

	if entry == nil {
		entry = &deviceEntry{
			primaryAddress: info.Address,
			addresses:      []byte{info.Address},
		}
		r.order = append(r.order, entry)
	} else if !containsAddress(entry.addresses, info.Address) {
		entry.addresses = append(entry.addresses, info.Address)
	}

	storedInfo := info
	if storedInfo.Manufacturer == "" {
		storedInfo.Manufacturer = entry.info.Manufacturer
	}
	preserveExistingDeviceID := existingByIdentity != nil && existingByIdentity != existingByAddress
	if preserveExistingDeviceID && entry.info.DeviceID != "" && storedInfo.DeviceID != entry.info.DeviceID {
		storedInfo.DeviceID = entry.info.DeviceID
	} else if storedInfo.DeviceID == "" {
		storedInfo.DeviceID = entry.info.DeviceID
	}
	if storedInfo.SoftwareVersion == "" {
		storedInfo.SoftwareVersion = entry.info.SoftwareVersion
	}
	if storedInfo.HardwareVersion == "" {
		storedInfo.HardwareVersion = entry.info.HardwareVersion
	}
	if storedInfo.SerialNumber == "" {
		storedInfo.SerialNumber = entry.info.SerialNumber
	}
	if storedInfo.MacAddress == "" {
		storedInfo.MacAddress = entry.info.MacAddress
	}
	storedInfo.Address = entry.primaryAddress
	physical = canonicalPhysicalIdentity(storedInfo)
	identityKey = physical.key()
	if identityKey == "" && entry.identityKey != "" {
		identityKey = entry.identityKey
	}

	for _, provider := range r.providers {
		if provider.Match(storedInfo) {
			matched = append(matched, provider)
			planes = append(planes, provider.CreatePlanes(storedInfo)...)
		}
	}

	projections := make([]Projection, 0)
	for _, provider := range matched {
		projectionProvider, ok := provider.(ProjectionProvider)
		if !ok {
			continue
		}
		projections = append(projections, projectionProvider.CreateProjections(storedInfo, planes)...)
	}

	index, projectionErr := BuildCanonicalIndex(projections)
	if projectionErr != nil {
		projections = nil
	}

	if entry.identityKey != "" && entry.identityKey != identityKey {
		// P0 round-6 (Codex P2 follow-up 2026-05-08): instead of
		// deleting the old primary key, re-point it at this entry
		// and track it as an alias. Without this, a Register call
		// that promotes a serial-derived key to primary while the
		// previous primary was MAC-derived would orphan the MAC
		// lookup path even though the entry still has the MAC in
		// info.MacAddress. Now `lookupByIdentity` by either key
		// continues to resolve to the merged entry, and
		// detachAddressLocked cleans up both via identityKeyAliases.
		//
		// P0 round-7 (Codex P2 follow-up 2026-05-10 on PR #136
		// thread PRRT_kwDORGIkfM6ArzFY): only preserve STABLE keys
		// (sn|... / mac|...) as aliases. Fallback signature keys
		// (`sig|...`) are NOT stable identifiers — multiple units
		// of the same model share the identical fallback signature.
		// If the old primary was sig-derived (entry first seen
		// without serial/MAC, then enriched with a stable serial),
		// preserving the sig key as an alias would silently bypass
		// `lookupCompatibleBySignatureLocked`'s ambiguity-refusal
		// scan when a second device with the same fingerprint
		// exists. A subsequent bare sig-only observation at a new
		// address would resolve directly to this entry via
		// r.identity instead of being routed through the ambiguity
		// check, incorrectly merging into this entry. Drop sig keys
		// rather than aliasing them.
		if isStableIdentityKey(entry.identityKey) {
			r.identity[entry.identityKey] = entry
			entry.identityKeyAliases = appendUniqueString(entry.identityKeyAliases, entry.identityKey)
		} else {
			delete(r.identity, entry.identityKey)
		}
	}
	entry.info = storedInfo
	entry.physical = physical
	entry.identityKey = identityKey
	entry.planes = planes
	entry.projections = projections
	entry.index = index
	entry.indexErr = projectionErr

	if identityKey != "" {
		r.identity[identityKey] = entry
	}
	r.entries[info.Address] = entry

	return entry
}

// RegisterStaticSeed plants identity for an address known from a
// product taxonomy table BEFORE any wire traffic has been observed.
// Mirrors Register's identity-merge behavior but stamps the AddressSlot
// with DiscoverySourceStaticSeed / VerificationStateCandidate so the
// observability surface (`/metrics`, MCP `bus.summary.get`,
// address-table snapshots) correctly shows the slot's provenance as a
// pre-known seed rather than active confirmation.
//
// On a clean cold boot a static-seeded slot subsequently observed
// passively will: NOT advance DiscoverySource (PassiveObserved <
// StaticSeed in the monotonic enum order), WILL advance
// VerificationState from Candidate to Corroborated. An active
// confirmation (e.g. directed scan) DOES advance DiscoverySource to
// ActiveConfirmed (StaticSeed < ActiveConfirmed) AND VerificationState
// to IdentityConfirmed.
//
// Single lock acquisition — composes registerLocked, then the
// shared static-seed stamping primitive, then syncEntryFacesLocked.
func (r *DeviceRegistry) RegisterStaticSeed(info DeviceInfo, role SlotRole, seededAt time.Time) DeviceEntry {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry := r.registerLocked(info)
	slot := r.ensureAddressSlotLocked(info.Address)
	slot.Device = entry
	r.markSlotStaticSeedLocked(slot, role, seededAt)
	r.syncEntryFacesLocked(entry)
	r.observationGeneration++
	return entry
}

// MarkSlotStaticSeed updates an AddressSlot for an address known from
// a product taxonomy seed table, mirroring MarkSlotPassiveObserved
// (lock discipline, monotonic upgrade semantics, idempotence) but
// stamping DiscoverySourceStaticSeed / VerificationStateCandidate.
//
// SCOPE: this API only mutates the AddressSlot. It does NOT attach
// the slot to a device entry; if `slot.Device` is nil at call time
// it stays nil, the address is NOT added to `r.entries`, and the
// device's `addresses` / `Faces` lists are NOT updated. Therefore:
//
//   - To plant a NEW seeded address with identity attached, use
//     RegisterStaticSeed (which composes registerLocked + this
//     stamp). Each face that should appear in `Lookup` /
//     `AddressByRole` needs its own RegisterStaticSeed call; the
//     identity-merge in registerLocked joins them when DeviceInfo
//     identity matches, or they can be aliased post-hoc via
//     AliasAddresses.
//
//   - The use case for MarkSlotStaticSeed in isolation is updating
//     an AddressSlot that was already attached to a device by some
//     prior path (Register / RegisterStaticSeed / AliasAddresses)
//     to upgrade its discovery_source / verification labels — for
//     example, marking a slot newly populated by passive observation
//     as "now also seeded from the static table" so the operator
//     surface reflects that the addresses are pre-known.
//
// Idempotent. Re-calling on a slot already at higher
// DiscoverySource (e.g. ActiveConfirmed) is a no-op for the discovery
// label, though it may still upgrade VerificationState if the
// existing state is below Candidate.
func (r *DeviceRegistry) MarkSlotStaticSeed(address byte, role SlotRole, seededAt time.Time) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	slot := r.ensureAddressSlotLocked(address)
	r.markSlotStaticSeedLocked(slot, role, seededAt)
	if slot.Device != nil {
		r.syncEntryFacesLocked(slot.Device)
	}
	r.observationGeneration++
}

// markSlotStaticSeedLocked is the shared slot-stamping primitive used
// by both MarkSlotStaticSeed and RegisterStaticSeed. Caller MUST hold
// r.mu and is responsible for any subsequent syncEntryFacesLocked
// call. Centralising the stamping rules here prevents drift between
// the two public entry points (Codex P3.5 review NIT FINDING_3).
func (r *DeviceRegistry) markSlotStaticSeedLocked(slot *AddressSlot, role SlotRole, seededAt time.Time) {
	if slot.DiscoverySource < DiscoverySourceStaticSeed {
		slot.DiscoverySource = DiscoverySourceStaticSeed
	}
	if slot.VerificationState < VerificationStateCandidate {
		slot.VerificationState = VerificationStateCandidate
	}
	if role != SlotRoleUnknown && slot.Role == SlotRoleUnknown {
		slot.Role = role
	}
	if slot.FirstObservedAt.IsZero() && !seededAt.IsZero() {
		slot.FirstObservedAt = seededAt
	}
	if !seededAt.IsZero() {
		slot.LastObservedAt = seededAt
	}
}
