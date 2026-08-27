package registry

import "strings"

func (r *DeviceRegistry) lookupByIdentity(info DeviceInfo) (DeviceEntry, bool) {
	identity := canonicalPhysicalIdentity(info).key()
	if identity == "" {
		return r.lookupBySignature(info)
	}

	r.mu.RLock()
	entry, ok := r.identity[identity]
	r.mu.RUnlock()
	if !ok {
		return r.lookupBySignature(info)
	}
	return entry, true
}

func (r *DeviceRegistry) lookupBySignature(info DeviceInfo) (DeviceEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.lookupCompatibleBySignatureLocked(info)
	if !ok {
		return nil, false
	}
	return entry, true
}

func (r *DeviceRegistry) lookupCompatibleBySignatureLocked(info DeviceInfo) (*deviceEntry, bool) {
	signature := canonicalPhysicalIdentity(info).withFallbackModelSignature()
	if signature == "" {
		return nil, false
	}
	var match *deviceEntry
	for _, candidate := range r.order {
		if candidate == nil {
			continue
		}
		if candidate.physical.withFallbackModelSignature() != signature {
			continue
		}
		if !canMergeIdentity(info, candidate.info) {
			continue
		}
		if match != nil && match != candidate {
			return nil, false
		}
		match = candidate
	}
	if match == nil {
		return nil, false
	}
	return match, true
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
				// Secondary is fully removed: this is the BASV2 /
				// NETX3 scenario where the identity-bearing entry
				// becomes orphaned. Promote secondary's identity
				// fields onto canonical (only when canonical's
				// fields are empty) and transfer the r.identity
				// binding. (Codex P2 round-2 finding 2026-05-08 on
				// PR #136: absorb must NOT fire when secondary
				// survives, otherwise canonical and the surviving
				// secondary expose duplicate identity.)
				r.absorbIdentityLocked(canonical, secondary)
				if secondary.identityKey != "" {
					switch canonical.identityKey {
					case "":
						// Canonical had no key; adopt secondary's
						// as canonical's primary key.
						canonical.identityKey = secondary.identityKey
						r.identity[canonical.identityKey] = canonical
					case secondary.identityKey:
						// Same key already bound; ensure the row
						// points to canonical (defensive).
						r.identity[canonical.identityKey] = canonical
					default:
						// Distinct keys (e.g. canonical has a MAC-
						// derived key, secondary has a serial-derived
						// key). Re-point secondary's key at canonical
						// so future lookupByIdentity calls on EITHER
						// key resolve to the merged entry. Do NOT
						// delete — that would orphan the lookup path.
						// (Codex P2 round-3 finding 2026-05-08 on
						// PR #136: previously the delete here meant
						// SerialNumber() was visible on the merged
						// entry but lookupByIdentity-by-serial could
						// not find it.)
						//
						// Track the alias key on canonical so
						// detachAddressLocked can clean it up if the
						// merged entry is later removed. Without
						// this, the orphan key would resolve to a
						// removed *deviceEntry until r.identity gets
						// rebuilt. (Codex P2 round-4 finding.)
						//
						// P0 round-7 (Codex P2 round-7 finding
						// 2026-05-10 on PR #136 thread
						// PRRT_kwDORGIkfM6ArzFY): only preserve
						// STABLE keys (sn|... / mac|...) as aliases.
						// `sig|...` fallback keys are NOT
						// per-device — preserving one would let a
						// later bare sig-only observation bypass
						// `lookupCompatibleBySignatureLocked`'s
						// ambiguity-refusal scan and silently merge
						// into canonical.
						if isStableIdentityKey(secondary.identityKey) {
							r.identity[secondary.identityKey] = canonical
							canonical.identityKeyAliases = appendUniqueString(canonical.identityKeyAliases, secondary.identityKey)
						} else {
							delete(r.identity, secondary.identityKey)
						}
					}
				}
				r.order = removeEntry(r.order, secondary)
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
				r.absorbIdentityLocked(canonical, secondary)
				if secondary.identityKey != "" {
					switch canonical.identityKey {
					case "":
						// Canonical had no key; adopt secondary's
						// as canonical's primary key.
						canonical.identityKey = secondary.identityKey
						r.identity[canonical.identityKey] = canonical
					case secondary.identityKey:
						// Same key already bound; ensure the row
						// points to canonical (defensive).
						r.identity[canonical.identityKey] = canonical
					default:
						// Distinct keys (e.g. canonical has a MAC-
						// derived key, secondary has a serial-derived
						// key). Re-point secondary's key at canonical
						// so future lookupByIdentity calls on EITHER
						// key resolve to the merged entry. Do NOT
						// delete — that would orphan the lookup path.
						// (Codex P2 round-3 finding 2026-05-08 on
						// PR #136: previously the delete here meant
						// SerialNumber() was visible on the merged
						// entry but lookupByIdentity-by-serial could
						// not find it.)
						//
						// Track the alias key on canonical so
						// detachAddressLocked can clean it up if the
						// merged entry is later removed. Without
						// this, the orphan key would resolve to a
						// removed *deviceEntry until r.identity gets
						// rebuilt. (Codex P2 round-4 finding.)
						//
						// P0 round-7 (Codex P2 round-7 finding
						// 2026-05-10 on PR #136 thread
						// PRRT_kwDORGIkfM6ArzFY): only preserve
						// STABLE keys (sn|... / mac|...) as aliases.
						// `sig|...` fallback keys are NOT
						// per-device — preserving one would let a
						// later bare sig-only observation bypass
						// `lookupCompatibleBySignatureLocked`'s
						// ambiguity-refusal scan and silently merge
						// into canonical.
						if isStableIdentityKey(secondary.identityKey) {
							r.identity[secondary.identityKey] = canonical
							canonical.identityKeyAliases = appendUniqueString(canonical.identityKeyAliases, secondary.identityKey)
						} else {
							delete(r.identity, secondary.identityKey)
						}
					}
				}
				r.order = removeEntry(r.order, secondary)
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

	r.observationGeneration++
	return nil
}

// appendUniqueString returns dst with s appended if not already
// present. Used to track identity-key aliases on a deviceEntry
// without duplication.
func appendUniqueString(dst []string, s string) []string {
	for _, existing := range dst {
		if existing == s {
			return dst
		}
	}
	return append(dst, s)
}

// isStableIdentityKey reports whether key is a stable, per-device
// identity key (serial- or MAC-derived) versus a fallback model
// signature key (`sig|...`) shared by every unit of the same model.
//
// Only stable keys are safe to preserve in r.identity as identity
// aliases when the entry's primary identityKey is rotated (e.g. on
// late-enrichment from sig-only → serial). Preserving a `sig|...` key
// would silently bypass `lookupCompatibleBySignatureLocked`'s
// ambiguity-refusal scan when a second device with the same
// fingerprint exists in the registry: a subsequent bare sig-only
// observation at a new address would resolve directly to the first
// entry via r.identity instead of being routed through the ambiguity
// check. (Codex P2 round-7 finding 2026-05-08 on PR #136 thread
// PRRT_kwDORGIkfM6ArzFY.)
//
// Mirrors the prefix taxonomy in physicalIdentity.key() /
// withFallbackModelSignature() in registry/identity.go.
func isStableIdentityKey(key string) bool {
	return strings.HasPrefix(key, "sn|") || strings.HasPrefix(key, "mac|")
}

// absorbIdentityLocked copies non-empty identity-bearing fields and
// derived state from src onto dst when dst's corresponding fields are
// empty. Re-keys r.identity[dst.identityKey] = dst when a new
// identityKey is adopted from src. Holds r.mu (caller's responsibility).
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
//   - identityKey (re-keyed in r.identity)
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
		// P0 round-7 follow-up (Codex P2 on PR #143): when dst.physical
		// is non-zero we don't replace it, but we just absorbed
		// info.DeviceID / SoftwareVersion / HardwareVersion above from
		// src. dst.physical is stale (was computed pre-absorb). A
		// future bare sig-only Register goes through
		// lookupCompatibleBySignatureLocked which compares
		// candidate.physical.withFallbackModelSignature() — if dst.physical
		// still lacks the absorbed sig fields, the candidate scan
		// won't match and a duplicate device is created. Recompute
		// dst.physical from the freshly-merged dst.info so the model
		// signature reflects everything we've absorbed.
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

	entry.addresses = removeAddress(entry.addresses, address)
	if len(entry.addresses) == 0 {
		if entry.identityKey != "" {
			delete(r.identity, entry.identityKey)
		}
		// P0 round-4 (Codex P2 follow-up 2026-05-08): also drop any
		// identityKeyAliases that AliasAddresses re-pointed at this
		// entry. Without this cleanup, orphan keys remain in
		// r.identity resolving to a removed *deviceEntry and a
		// later Register({key}) would attach to an entry no longer
		// in r.order.
		for _, alias := range entry.identityKeyAliases {
			if alias == "" {
				continue
			}
			if r.identity[alias] == entry {
				delete(r.identity, alias)
			}
		}
		entry.identityKeyAliases = nil
		r.order = removeEntry(r.order, entry)
		return
	}
	if entry.primaryAddress == address {
		entry.primaryAddress = entry.addresses[0]
		entry.info.Address = entry.primaryAddress
	}
	r.syncEntryFacesLocked(entry)
}
