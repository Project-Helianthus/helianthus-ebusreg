//go:build !tinygo
// +build !tinygo

package ebus_standard_catalog

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// loadCatalogImpl parses a YAML document, validates it, and populates
// ContentSHA256 from the raw bytes. It returns a wrapped sentinel error on
// collision or ambiguity.
func loadCatalogImpl(data []byte) (Catalog, error) {
	var cat Catalog
	if err := yaml.Unmarshal(data, &cat); err != nil {
		return Catalog{}, fmt.Errorf("ebus_standard: yaml unmarshal: %w", err)
	}

	// Default namespace from constant when not explicitly set. Every
	// identity key must still carry it, but the document-level namespace
	// is informational.
	if cat.Namespace == "" {
		cat.Namespace = Namespace
	}

	// Reject catalogs with no services at all. A YAML document that omits
	// the `services:` key (or supplies an empty list) would otherwise load
	// successfully as an effectively empty catalog, silently dropping every
	// method definition at runtime and bypassing the per-service typo/
	// omission guards (ErrServiceMissingPB, ErrServiceMissingCommands,
	// identity/namespace checks, duplicate detection). Fail fast here so a
	// malformed document cannot masquerade as a valid (but empty) catalog.
	if len(cat.Services) == 0 {
		return Catalog{}, fmt.Errorf("%w: document has no services: entries",
			ErrCatalogMissingServices)
	}

	// Validate identity-key completeness and check safety_class values.
	for si := range cat.Services {
		svc := &cat.Services[si]
		// Reject services without an explicit `pb:` key. A value-typed
		// uint8 would silently deserialize omission as 0x00 and defeat
		// the service/identity mismatch check below when an identity.pb
		// also happens to be 0x00. Presence is checked before the empty-
		// commands guard so that diagnostics include the pb axis even
		// when multiple fields are malformed.
		if svc.PB == nil {
			return Catalog{}, fmt.Errorf("%w: service %q",
				ErrServiceMissingPB, svc.Name)
		}
		// Reject services with no commands. An empty commands list is
		// almost always a YAML typo (wrong key name, missing block) and
		// must fail loudly rather than silently accept a service with no
		// method definitions.
		if len(svc.Commands) == 0 {
			return Catalog{}, fmt.Errorf("%w: service %q (pb=0x%02X)",
				ErrServiceMissingCommands, svc.Name, svc.PBValue())
		}
		for ci := range svc.Commands {
			cmd := &svc.Commands[ci]
			if !cmd.Identity.IsComplete() {
				return Catalog{}, fmt.Errorf("%w: command %q", ErrIncompleteIdentityKey, cmd.ID)
			}
			if err := validateNamespace(cmd.ID, cmd.Identity); err != nil {
				return Catalog{}, err
			}
			// Guard against a typo in the service header that would
			// otherwise silently group commands under the wrong service
			// code. IsComplete() only asserts PB is non-nil, and the
			// duplicate detector fingerprints the identity pb (not the
			// service pb), so a mismatch is invisible without this check.
			if cmd.Identity.PBValue() != svc.PBValue() {
				return Catalog{}, fmt.Errorf(
					"%w: service %q pb=0x%02X, command %q identity.pb=0x%02X",
					ErrServicePBMismatch, svc.Name, svc.PBValue(), cmd.ID, cmd.Identity.PBValue())
			}
			if !isKnownSafetyClass(cmd.SafetyClass) {
				return Catalog{}, fmt.Errorf("%w: command %q safety_class=%q", ErrUnknownSafetyClass, cmd.ID, cmd.SafetyClass)
			}
			if err := validateIdentityEnums(cmd.ID, cmd.Identity); err != nil {
				return Catalog{}, err
			}
		}
	}

	// Duplicate 14-tuple detection.
	if err := detectDuplicateIdentityKeys(cat); err != nil {
		return Catalog{}, err
	}

	// Ambiguous length-selector detection.
	if err := detectAmbiguousLengthSelectors(cat); err != nil {
		return Catalog{}, err
	}

	cat.ContentSHA256 = ComputeContentSHA256(data)
	return cat, nil
}

// validateIdentityEnums checks every enum-typed identity axis against the
// constants declared in identity.go. Any value that is non-empty but not
// one of the enumerated constants is rejected with ErrUnknownEnumValue so
// typos fail deterministically at load time rather than silently breaking
// downstream matching.
//
// Empty values are NOT handled here — they are caught earlier by
// IdentityKey.IsComplete(), which returns ErrIncompleteIdentityKey.
func validateIdentityEnums(cmdID string, k IdentityKey) error {
	switch k.TelegramClass {
	case TelegramClassAddressed, TelegramClassBroadcast,
		TelegramClassInitiatorInitiator, TelegramClassControllerBroadcast:
	default:
		return fmt.Errorf("%w: command %q field=telegram_class value=%q",
			ErrUnknownEnumValue, cmdID, k.TelegramClass)
	}
	switch k.Direction {
	case DirectionRequest, DirectionResponse:
	default:
		return fmt.Errorf("%w: command %q field=direction value=%q",
			ErrUnknownEnumValue, cmdID, k.Direction)
	}
	switch k.RequestOrResponseRole {
	case RoleInitiator, RoleResponder, RoleOriginator:
	default:
		return fmt.Errorf("%w: command %q field=request_or_response_role value=%q",
			ErrUnknownEnumValue, cmdID, k.RequestOrResponseRole)
	}
	switch k.BroadcastOrAddressed {
	case AddressedDirect, AddressedBroadcast:
	default:
		return fmt.Errorf("%w: command %q field=broadcast_or_addressed value=%q",
			ErrUnknownEnumValue, cmdID, k.BroadcastOrAddressed)
	}
	switch k.AnswerPolicy {
	case AnswerRequired, AnswerNone:
	default:
		return fmt.Errorf("%w: command %q field=answer_policy value=%q",
			ErrUnknownEnumValue, cmdID, k.AnswerPolicy)
	}
	switch k.LengthPrefixMode {
	case LengthPrefixNone, LengthPrefixFixed, LengthPrefixByte, LengthPrefixTypedPayload:
	default:
		return fmt.Errorf("%w: command %q field=length_prefix_mode value=%q",
			ErrUnknownEnumValue, cmdID, k.LengthPrefixMode)
	}
	return nil
}

// validateNamespace enforces that the identity-key namespace matches the
// package's fixed Namespace constant exactly. IdentityKey.IsComplete() only
// asserts non-emptiness, so a typo like "ebus_standrad" would otherwise slip
// through and bypass duplicate detection (fingerprints include Namespace).
func validateNamespace(cmdID string, k IdentityKey) error {
	if k.Namespace != Namespace {
		return fmt.Errorf("%w: command %q field=namespace value=%q",
			ErrInvalidNamespace, cmdID, k.Namespace)
	}
	return nil
}

func isKnownSafetyClass(s SafetyClass) bool {
	switch s {
	case SafetyReadOnlySafe, SafetyReadOnlyBusLoad, SafetyMutating,
		SafetyDestructive, SafetyBroadcast, SafetyMemoryWrite:
		return true
	}
	return false
}

// identityKeyFingerprint returns a deterministic string covering every
// field of the 14-tuple. Equal fingerprints mean duplicate identity keys.
//
// Slice/map fields MUST be serialized with an injective encoding: naive
// %v formatting of []string renders both []string{"a b"} and
// []string{"a","b"} as "[a b]", which would collapse two distinct
// identities into one fingerprint and cause false-positive
// ErrDuplicateIdentityKey (or, worse, hide real collisions). JSON
// encoding is canonical, escapes embedded separators, and round-trips
// losslessly for any []string, so it is used for every slice/map field
// in the tuple (currently just TransportCapabilityRequirements).
func identityKeyFingerprint(k IdentityKey) string {
	tcrJSON, err := json.Marshal(k.TransportCapabilityRequirements)
	if err != nil {
		// json.Marshal of []string cannot fail in practice; fall back to
		// a clearly non-injective sentinel so a future failure is
		// visible rather than silently masked.
		tcrJSON = []byte(fmt.Sprintf("<json-error:%v>", err))
	}
	return fmt.Sprintf(
		"ns=%s|pb=%02X|sb=%02X|sel=%s|tc=%s|dir=%s|rr=%s|ba=%s|ap=%s|lpm=%s|sd=%s|sv=%s|tcr=%s|ver=%s",
		k.Namespace, k.PBValue(), k.SBValue(), k.SelectorPath, k.TelegramClass, k.Direction,
		k.RequestOrResponseRole, k.BroadcastOrAddressed, k.AnswerPolicy,
		k.LengthPrefixMode, k.SelectorDecoder, k.ServiceVariant,
		string(tcrJSON), k.Version,
	)
}

func detectDuplicateIdentityKeys(cat Catalog) error {
	seen := make(map[string]string)
	for _, svc := range cat.Services {
		for _, cmd := range svc.Commands {
			fp := identityKeyFingerprint(cmd.Identity)
			if prev, ok := seen[fp]; ok {
				return fmt.Errorf("%w: %q collides with %q (fingerprint=%s)",
					ErrDuplicateIdentityKey, cmd.ID, prev, fp)
			}
			seen[fp] = cmd.ID
		}
	}
	return nil
}

// detectAmbiguousLengthSelectors flags entries that share
// (namespace, PB, SB, selector_decoder, selector_path, direction, role) but
// carry incompatible length_prefix_mode values. Such entries cannot be
// disambiguated at decode time because the only differing axis is the
// length-prefix rule itself, which the decoder applies BEFORE it knows
// which branch to pick.
//
// When selector_decoder is "none" (or empty), there is no selector branch
// at all — so two entries sharing the remaining on-wire identity axes but
// differing only by length_prefix_mode are STILL ambiguous (there is no
// branch to disambiguate on). These cases are bundled under a canonical
// "none" decoder key rather than skipped.
func detectAmbiguousLengthSelectors(cat Catalog) error {
	type bucket struct {
		cmdID string
		lpm   LengthPrefixMode
	}
	buckets := make(map[string][]bucket)
	for _, svc := range cat.Services {
		for _, cmd := range svc.Commands {
			k := cmd.Identity
			// Canonicalise empty selector_decoder to "none" so that the
			// two values are treated as the same bundling axis and do not
			// fragment otherwise-ambiguous buckets.
			decoder := k.SelectorDecoder
			if decoder == "" {
				decoder = "none"
			}
			// When there is no selector branch (decoder == "none"),
			// selector_path is not a real disambiguator: with no decoder
			// to interpret it, the value is inert at decode time. Two
			// entries with the same on-wire identity and different
			// length_prefix_mode values must still collide even if they
			// carry different selector_path strings, otherwise a YAML
			// typo or cosmetic difference bypasses the ambiguity check.
			// Drop selector_path from the bundling key in this case.
			selectorPath := k.SelectorPath
			if decoder == "none" {
				selectorPath = ""
			}
			key := fmt.Sprintf("%s|%02X|%02X|%s|%s|%s|%s",
				k.Namespace, k.PBValue(), k.SBValue(), decoder,
				selectorPath, k.Direction, k.RequestOrResponseRole)
			buckets[key] = append(buckets[key], bucket{cmd.ID, k.LengthPrefixMode})
		}
	}
	for key, entries := range buckets {
		if len(entries) < 2 {
			continue
		}
		// Ambiguous if any two entries in the same bucket have different
		// length_prefix_mode values.
		first := entries[0].lpm
		for _, e := range entries[1:] {
			if e.lpm != first {
				return fmt.Errorf("%w: bucket %s: %q(lpm=%s) vs %q(lpm=%s)",
					ErrAmbiguousLengthSelector, key,
					entries[0].cmdID, first, e.cmdID, e.lpm)
			}
		}
	}
	return nil
}

// ComputeContentSHA256 returns the lowercase hex SHA-256 of the given bytes.
func ComputeContentSHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
