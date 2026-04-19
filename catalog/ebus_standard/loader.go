package ebus_standard_catalog

import (
	"errors"
)

// Sentinel errors returned by the catalog loader.
var (
	// ErrDuplicateIdentityKey is returned when two catalog entries share
	// the full 14-tuple identity key.
	ErrDuplicateIdentityKey = errors.New("ebus_standard: duplicate catalog identity key")

	// ErrAmbiguousLengthSelector is returned when two entries share
	// (namespace, PB, SB, selector_decoder) but have incompatible
	// length-selector semantics that cannot be disambiguated at decode
	// time.
	ErrAmbiguousLengthSelector = errors.New("ebus_standard: ambiguous length-selector branch")

	// ErrIncompleteIdentityKey is returned when a catalog entry fails to
	// populate every required field of the 14-tuple identity key.
	ErrIncompleteIdentityKey = errors.New("ebus_standard: catalog entry has incomplete identity key")

	// ErrUnknownSafetyClass is returned when an entry's safety_class is
	// not one of the enumerated values.
	ErrUnknownSafetyClass = errors.New("ebus_standard: unknown safety_class")

	// ErrUnknownEnumValue is returned when an identity-key enum field
	// (telegram_class, direction, request_or_response_role,
	// broadcast_or_addressed, answer_policy, length_prefix_mode) carries
	// a value that is not one of the constants defined in identity.go.
	ErrUnknownEnumValue = errors.New("ebus_standard: unknown enum value")

	// ErrServiceMissingCommands is returned when a service entry deserializes
	// with an empty commands list. A service with no commands is always a
	// typo or omission in the YAML (wrong key name, empty list) and must
	// fail loudly rather than silently load an empty service.
	ErrServiceMissingCommands = errors.New("ebus_standard: service has no commands")

	// ErrInvalidNamespace is returned when an identity-key namespace does
	// not exactly match the package's fixed Namespace constant. A typo
	// (e.g. "ebus_standrad") is otherwise silently treated as a distinct
	// identity and bypasses duplicate detection.
	ErrInvalidNamespace = errors.New("ebus_standard: identity-key namespace must be ebus_standard")
)

// LoadCatalog parses and validates a YAML catalog document. The returned
// catalog has ContentSHA256 populated from the raw bytes.
//
// Implementation lives in loader_yaml.go so that this file can be consumed
// by TinyGo builds without the yaml.v3 dependency. M2 RED scaffolding keeps
// a panicking stub here only until loader_yaml.go is in place.
func LoadCatalog(data []byte) (Catalog, error) {
	return loadCatalogImpl(data)
}
