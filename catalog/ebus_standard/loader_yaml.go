package ebus_standard_catalog

import (
	"crypto/sha256"
	"encoding/hex"
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

	// Validate identity-key completeness and check safety_class values.
	for si := range cat.Services {
		svc := &cat.Services[si]
		for ci := range svc.Commands {
			cmd := &svc.Commands[ci]
			if !cmd.Identity.IsComplete() {
				return Catalog{}, fmt.Errorf("%w: command %q", ErrIncompleteIdentityKey, cmd.ID)
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
func identityKeyFingerprint(k IdentityKey) string {
	return fmt.Sprintf(
		"ns=%s|pb=%02X|sb=%02X|sel=%s|tc=%s|dir=%s|rr=%s|ba=%s|ap=%s|lpm=%s|sd=%s|sv=%s|tcr=%v|ver=%s",
		k.Namespace, k.PBValue(), k.SBValue(), k.SelectorPath, k.TelegramClass, k.Direction,
		k.RequestOrResponseRole, k.BroadcastOrAddressed, k.AnswerPolicy,
		k.LengthPrefixMode, k.SelectorDecoder, k.ServiceVariant,
		k.TransportCapabilityRequirements, k.Version,
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
func detectAmbiguousLengthSelectors(cat Catalog) error {
	type bucket struct {
		cmdID string
		lpm   LengthPrefixMode
	}
	buckets := make(map[string][]bucket)
	for _, svc := range cat.Services {
		for _, cmd := range svc.Commands {
			k := cmd.Identity
			if k.SelectorDecoder == "" || k.SelectorDecoder == "none" {
				continue
			}
			key := fmt.Sprintf("%s|%02X|%02X|%s|%s|%s|%s",
				k.Namespace, k.PBValue(), k.SBValue(), k.SelectorDecoder,
				k.SelectorPath, k.Direction, k.RequestOrResponseRole)
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
