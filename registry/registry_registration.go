package registry

import "time"

func (r *DeviceRegistry) Register(info DeviceInfo) DeviceEntry {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry := r.registerLocked(info)
	state := VerificationStateCandidate
	if canonicalPhysicalIdentity(info).isQualified() {
		state = VerificationStateIdentityConfirmed
	}
	r.observeAddressSlotLocked(info.Address, entry, DiscoverySourceActiveConfirmed, state)
	if state == VerificationStateIdentityConfirmed {
		r.confirmEntryIdentitySlotsLocked(entry)
	}
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
// syncEntryFacesLocked, detachAddressLocked).
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
	incomingHasStableIdentity := physical.isQualified()
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
	// A serial-key match can conceal a contradictory MAC (or vice versa
	// through a retained identity alias). Keep an existing address group
	// unchanged instead of applying ambiguous evidence to either device.
	if existingByAddress != nil && existingByIdentity != nil &&
		existingByIdentity != existingByAddress && !canMergeIdentity(info, existingByIdentity.info) {
		return existingByAddress
	}

	entry := existingByIdentity
	if entry == nil {
		entry = existingByAddress
	}
	if entry == existingByAddress && existingByAddress != nil &&
		(!canMergeIdentity(info, existingByAddress.info) ||
			(!physical.isQualified() && hasConflictingModelSignature(info, existingByAddress.info))) {
		entry = nil
	}
	// A model signature cannot disband an already-evidenced alias group.
	// Keep it together while a later serial/MAC observation establishes
	// whether it belongs to another physical entry.
	if existingByAddress != nil && len(existingByAddress.addresses) > 1 &&
		entry != existingByAddress && !incomingHasStableIdentity &&
		compatibleAliasEnrichment(info, existingByAddress) {
		entry = existingByAddress
		existingByIdentity = nil
	}

	if existingByAddress != nil && existingByAddress != entry {
		if existingByIdentity != nil && isStableIdentityKey(identityKey) &&
			compatibleAliasEnrichment(info, existingByAddress) &&
			compatibleAliasEnrichment(entry.info, existingByAddress) {
			r.mergeRegisteredAliasesLocked(entry, existingByAddress)
		} else {
			r.detachAddressLocked(existingByAddress, info.Address)
		}
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
		// Retain qualified identity aliases when a same-address correction
		// rotates the canonical triple. Partial identity and model signatures
		// never enter this map, so they cannot merge independently observed
		// addresses.
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

// confirmEntryIdentitySlotsLocked promotes the explicitly grouped faces of an
// entry after an active observation supplied a complete qualified identity.
// It never promotes on partial or sentinel identity, so static-seed candidates
// remain candidates until this evidence arrives.
func (r *DeviceRegistry) confirmEntryIdentitySlotsLocked(entry *deviceEntry) {
	if entry == nil {
		return
	}
	for _, address := range entry.addresses {
		slot := r.ensureAddressSlotLocked(address)
		if slot.VerificationState < VerificationStateIdentityConfirmed {
			slot.VerificationState = VerificationStateIdentityConfirmed
		}
	}
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

// mergeRegisteredAliasesLocked joins an already-evidenced address group after
// a stable identity match. Moving just the enriched address would strand its
// companion and discard the address slot's observation history. The caller
// holds r.mu and has rejected conflicting identities and model signatures.
func (r *DeviceRegistry) mergeRegisteredAliasesLocked(dst, src *deviceEntry) {
	r.absorbIdentityLocked(dst, src)
	for _, address := range src.addresses {
		if !containsAddress(dst.addresses, address) {
			dst.addresses = append(dst.addresses, address)
		}
		r.entries[address] = dst
		r.ensureAddressSlotLocked(address).Device = dst
	}
	for _, key := range append(src.identityKeyAliases, src.identityKey) {
		if key == "" || r.identity[key] != src {
			continue
		}
		if isStableIdentityKey(key) {
			r.identity[key] = dst
			if key != dst.identityKey {
				dst.identityKeyAliases = appendUniqueString(dst.identityKeyAliases, key)
			}
		} else {
			delete(r.identity, key)
		}
	}
	r.order = removeEntry(r.order, src)
}

func compatibleAliasEnrichment(info DeviceInfo, previous *deviceEntry) bool {
	manufacturer := normalizeIdentityPart(info.Manufacturer)
	return canMergeIdentity(info, previous.info) &&
		!hasConflictingModelSignature(info, previous.info) &&
		(manufacturer == "" || previous.physical.manufacturer == "" || manufacturer == previous.physical.manufacturer)
}
