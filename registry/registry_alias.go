package registry

import "strings"

func (r *DeviceRegistry) lookupByIdentity(info DeviceInfo) (DeviceEntry, bool) {
	identity := canonicalPhysicalIdentity(info).key()
	if identity == "" {
		return nil, false
	}

	r.mu.RLock()
	entry, ok := r.identity[identity]
	ok = ok && isCurrentQualifiedIdentityKey(entry, identity)
	r.mu.RUnlock()
	return entry, ok
}

// currentIdentityEntryLocked returns an index candidate only when its key is
// the entry's current complete normalized triple. Qualified triples are exact
// current identity claims, not historical evidence. Caller holds r.mu.
func (r *DeviceRegistry) currentIdentityEntryLocked(identity string) *deviceEntry {
	entry := r.identity[identity]
	if !isCurrentQualifiedIdentityKey(entry, identity) {
		if entry != nil && r.identity[identity] == entry {
			delete(r.identity, identity)
		}
		return nil
	}
	return entry
}

func isCurrentQualifiedIdentityKey(entry *deviceEntry, identity string) bool {
	return entry != nil && isStableIdentityKey(identity) &&
		entry.identityKey == identity && entry.physical.key() == identity
}

// retainCurrentIdentityBindingLocked removes every historical identity
// binding for entry. Explicit address aliases stay in r.entries and topology,
// but a distinct old triple cannot select this entry across addresses.
func (r *DeviceRegistry) retainCurrentIdentityBindingLocked(entry *deviceEntry, identity string) {
	if entry == nil {
		return
	}
	for key, indexed := range r.identity {
		if indexed == entry && key != identity {
			delete(r.identity, key)
		}
	}
	entry.identityKeyAliases = nil
}

func (r *DeviceRegistry) removeIdentityBindingsLocked(entry *deviceEntry) {
	r.retainCurrentIdentityBindingLocked(entry, "")
}

func canMergeIdentity(incoming DeviceInfo, existing DeviceInfo) bool {
	normalizedIncomingSerial := normalizeIdentityPart(incoming.SerialNumber)
	normalizedExistingSerial := normalizeIdentityPart(existing.SerialNumber)
	if normalizedIncomingSerial != "" && normalizedExistingSerial != "" && normalizedIncomingSerial != normalizedExistingSerial {
		return false
	}
	normalizedIncomingMAC := normalizeIdentityPart(incoming.MacAddress)
	normalizedExistingMAC := normalizeIdentityPart(existing.MacAddress)
	if normalizedIncomingMAC != "" && normalizedExistingMAC != "" && normalizedIncomingMAC != normalizedExistingMAC {
		return false
	}
	return true
}

func hasConflictingModelSignature(incoming DeviceInfo, existing DeviceInfo) bool {
	incomingDeviceID := normalizeIdentityPart(incoming.DeviceID)
	existingDeviceID := normalizeIdentityPart(existing.DeviceID)
	incomingSoftware := normalizeIdentityPart(incoming.SoftwareVersion)
	existingSoftware := normalizeIdentityPart(existing.SoftwareVersion)
	incomingHardware := normalizeIdentityPart(incoming.HardwareVersion)
	existingHardware := normalizeIdentityPart(existing.HardwareVersion)

	if incomingDeviceID == "" || existingDeviceID == "" || incomingSoftware == "" || existingSoftware == "" || incomingHardware == "" || existingHardware == "" {
		return false
	}

	return incomingDeviceID != existingDeviceID || incomingSoftware != existingSoftware || incomingHardware != existingHardware
}

// Lookup returns the canonical *DeviceEntry for the given address, or
// (nil, false) if no device occupies that slot. Preserved (signature
// unchanged) so existing callers continue to compile.
func (r *DeviceRegistry) AliasAddresses(a, b byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	slotA := r.ensureAddressSlotLocked(a)
	slotB := r.ensureAddressSlotLocked(b)

	switch {
	case slotA.Device != nil:
		canonical := slotA.Device
		if secondary := slotB.Device; secondary != nil && secondary != canonical {
			secondary.addresses = removeAddress(secondary.addresses, b)
			if len(secondary.addresses) == 0 {
				r.mergeRemovedAliasEntryLocked(canonical, secondary)
			} else {
				// Secondary survives at remaining addresses (e.g.
				// 0x15 + 0x16 with same identity, then alias 0x10
				// to 0x15 → secondary still owns 0x16 with the
				// original identity). DO NOT absorb identity onto
				// canonical and DO NOT transfer identityKey:
				// secondary keeps its r.identity row + manufacturer
				// fields, canonical gets identity via Register's
				// identity-merge path on a future write.
				if secondary.primaryAddress == b {
					secondary.primaryAddress = secondary.addresses[0]
					secondary.info.Address = secondary.primaryAddress
				}
				r.syncEntryFacesLocked(secondary)
			}
		}
		slotB.Device = canonical
		if !containsAddress(canonical.addresses, b) {
			canonical.addresses = append(canonical.addresses, b)
		}
		r.entries[b] = canonical
		r.syncEntryFacesLocked(canonical)
	case slotB.Device != nil:
		canonical := slotB.Device
		if secondary := slotA.Device; secondary != nil && secondary != canonical {
			// Symmetric case. See slotA branch.
			secondary.addresses = removeAddress(secondary.addresses, a)
			if len(secondary.addresses) == 0 {
				r.mergeRemovedAliasEntryLocked(canonical, secondary)
			} else {
				if secondary.primaryAddress == a {
					secondary.primaryAddress = secondary.addresses[0]
					secondary.info.Address = secondary.primaryAddress
				}
				r.syncEntryFacesLocked(secondary)
			}
		}
		slotA.Device = canonical
		if !containsAddress(canonical.addresses, a) {
			canonical.addresses = append(canonical.addresses, a)
		}
		r.entries[a] = canonical
		r.syncEntryFacesLocked(canonical)
	}

	if r.entries[a] != nil && r.entries[a] == r.entries[b] {
		r.recordTopologyAliasLocked(a, b)
	}
	r.observationGeneration++
	return nil
}

// mergeRemovedAliasEntryLocked retains useful fields from an explicitly
// aliased, now-removed entry without retaining that entry's distinct triple as
// a lookup or cross-address merge authority.
func (r *DeviceRegistry) mergeRemovedAliasEntryLocked(canonical, secondary *deviceEntry) {
	r.absorbIdentityLocked(canonical, secondary)
	r.removeIdentityBindingsLocked(secondary)

	canonical.physical = canonicalPhysicalIdentity(canonical.info)
	canonical.identityKey = canonical.physical.key()
	r.retainCurrentIdentityBindingLocked(canonical, canonical.identityKey)
	if canonical.identityKey != "" {
		r.identity[canonical.identityKey] = canonical
	}
	r.order = removeEntry(r.order, secondary)
}

// recordTopologyAliasLocked remembers explicit source-target or canonical-
// companion evidence without conflating it with qualified identity joins.
// AliasAddresses has already applied the relation; this method only records
// it for later current-session confirmation propagation. Caller holds r.mu.
func (r *DeviceRegistry) recordTopologyAliasLocked(a, b byte) {
	if a == b {
		return
	}
	if r.topology[a] == nil {
		r.topology[a] = make(map[byte]struct{})
	}
	if r.topology[b] == nil {
		r.topology[b] = make(map[byte]struct{})
	}
	r.topology[a][b] = struct{}{}
	r.topology[b][a] = struct{}{}
}

// forgetTopologyAddressLocked invalidates every explicit topology relation
// involving address. An identity conflict can detach one face from an alias
// group and later rejoin it through the qualified-identity index alone; the
// old relation is not current-session evidence for that new membership.
//
// Caller holds r.mu.
func (r *DeviceRegistry) forgetTopologyAddressLocked(address byte) {
	for companion := range r.topology[address] {
		delete(r.topology[companion], address)
		if len(r.topology[companion]) == 0 {
			delete(r.topology, companion)
		}
	}
	delete(r.topology, address)
}

// isStableIdentityKey reports whether key is a stable, per-device
// identity key: a complete, normalized manufacturer/device/serial triple.
//
// Partial identity and model evidence never enter r.identity, so they cannot
// select or merge independently observed addresses.
func isStableIdentityKey(key string) bool {
	return strings.HasPrefix(key, stableIdentityKeyPrefix)
}

// absorbIdentityLocked copies non-empty identity-bearing fields and
// derived state from src onto dst when dst's corresponding fields are
// empty. Its caller rebuilds the current r.identity binding after deciding
// whether the source entry survives. Holds r.mu (caller's responsibility).
//
// Phase post-C P0: introduced to fix AliasAddresses identity loss
// (see live observation on 2026-05-08: BASV2 0x10↔0x15 + NETX3
// 0xF1↔0xF6 pairs aliased correctly but with manufacturer="" because
// the identity-bearing target-face entry was the secondary in the
// merge).
//
// Fields absorbed (only when dst's value is empty/zero AND src's value
// is non-empty):
//   - info.Manufacturer, info.DeviceID, info.SerialNumber, info.MacAddress
//   - info.SoftwareVersion, info.HardwareVersion
//   - physical (only if dst's physicalIdentity is zero)
//   - planes, projections, index, indexErr (only when dst has none)
//
// This function does NOT touch addresses or primaryAddress — those are
// owned by the caller (AliasAddresses) since the merge's address-graph
// semantics are independent of identity.
func (r *DeviceRegistry) absorbIdentityLocked(dst, src *deviceEntry) {
	if dst == nil || src == nil || dst == src {
		return
	}
	if dst.info.Manufacturer == "" && src.info.Manufacturer != "" {
		dst.info.Manufacturer = src.info.Manufacturer
	}
	if dst.info.DeviceID == "" && src.info.DeviceID != "" {
		dst.info.DeviceID = src.info.DeviceID
	}
	if dst.info.SerialNumber == "" && src.info.SerialNumber != "" {
		dst.info.SerialNumber = src.info.SerialNumber
	}
	if dst.info.MacAddress == "" && src.info.MacAddress != "" {
		dst.info.MacAddress = src.info.MacAddress
	}
	if dst.info.SoftwareVersion == "" && src.info.SoftwareVersion != "" {
		dst.info.SoftwareVersion = src.info.SoftwareVersion
	}
	if dst.info.HardwareVersion == "" && src.info.HardwareVersion != "" {
		dst.info.HardwareVersion = src.info.HardwareVersion
	}
	// Adopt physicalIdentity + identityKey if dst has neither and src has both.
	emptyPhysical := physicalIdentity{}
	if dst.physical == emptyPhysical && src.physical != emptyPhysical {
		dst.physical = src.physical
	} else if dst.physical != emptyPhysical {
		// We just absorbed fields into dst.info, so its cached normalized
		// identity must reflect the same-address topology merge.
		dst.physical = canonicalPhysicalIdentity(dst.info)
	}
	// NOTE: identityKey transfer is intentionally NOT done here.
	// absorbIdentityLocked runs BEFORE AliasAddresses removes
	// `b`/`a` from src.addresses, so we don't yet know whether src
	// will survive (multiple addresses → survives at remaining face)
	// or be fully removed (only the aliased address → removed). The
	// caller (AliasAddresses) handles identityKey transfer in the
	// "secondary fully removed" branch only — see the comment at
	// the call site. This avoids the live-validation P2 finding
	// where surviving multi-address secondaries lost their identity
	// map binding because absorb prematurely moved identityKey.
	//
	// Adopt planes / projections / index if dst has none. These
	// derive from info; absorbing them avoids re-running providers.
	if len(dst.planes) == 0 && len(src.planes) > 0 {
		dst.planes = src.planes
	}
	if len(dst.projections) == 0 && len(src.projections) > 0 {
		dst.projections = src.projections
	}
	if dst.index.canonicalByID == nil && src.index.canonicalByID != nil {
		dst.index = src.index
	}
	if dst.indexErr == nil && src.indexErr != nil {
		dst.indexErr = src.indexErr
	}
}

func (r *DeviceRegistry) Iterate(fn func(DeviceEntry) bool) {
	r.mu.RLock()
	order := make([]*deviceEntry, len(r.order))
	copy(order, r.order)
	r.mu.RUnlock()

	for _, entry := range order {
		if entry == nil {
			continue
		}
		if !fn(entry) {
			return
		}
	}
}

func (r *DeviceRegistry) detachAddressLocked(entry *deviceEntry, address byte) {
	if entry == nil {
		return
	}
	delete(r.entries, address)
	if slot := r.addressTable[address]; slot != nil && slot.Device == entry {
		r.addressTable[address] = nil
	}
	if !containsAddress(entry.addresses, address) {
		return
	}

	// Leaving an address group invalidates any topology evidence attached to
	// that face. A later qualified-identity rejoin has no topology semantics;
	// only a newly observed AliasAddresses relation may restore propagation.
	r.forgetTopologyAddressLocked(address)
	entry.addresses = removeAddress(entry.addresses, address)
	if len(entry.addresses) == 0 {
		r.removeIdentityBindingsLocked(entry)
		r.order = removeEntry(r.order, entry)
		return
	}
	if entry.primaryAddress == address {
		entry.primaryAddress = entry.addresses[0]
		entry.info.Address = entry.primaryAddress
	}
	r.syncEntryFacesLocked(entry)
}
