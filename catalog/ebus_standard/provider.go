package ebus_standard_catalog

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Project-Helianthus/helianthus-ebusreg/catalog/ebus_standard/internal/safetypolicy"
)

// Provider is the generic ebus_standard L7 provider. Method dispatch is
// driven by the loaded catalog; there is no per-device hard-coding in this
// package.
//
// Namespace isolation (M3 AC6): this provider and its internal helpers MUST
// NOT import any Vaillant-specific package under github.com/Project-
// Helianthus/helianthus-ebusreg/vaillant/.... Compile-time enforcement is
// provided by placing the safety gate in catalog/ebus_standard/internal/
// safetypolicy (Go's internal/ visibility rule excludes packages outside
// catalog/ebus_standard/). A second line of defense is the static-import
// regression test in provider_namespace_isolation_test.go which parses
// every .go file in this package and fails if a Vaillant import appears.
//
// Disable switch (M3 AC4): the env variable EBUS_STANDARD_PROVIDER_ENABLED
// is read once at provider construction via NewProviderFromEnv. When the
// value is literally "0" or "false" (case-insensitive), the provider is
// constructed in the disabled state and every call to Invoke /
// Identification returns ErrProviderDisabled. Default is enabled.
type Provider struct {
	catalog   Catalog
	methodIdx map[string]Command
	enabled   bool
}

// IdentificationDescriptor carries the canonical manufacturer / device-id /
// software / hardware values produced by the 0x07 0x04 service, wrapped in
// provenance metadata (never mutates DeviceInfo).
type IdentificationDescriptor struct {
	Manufacturer    string
	DeviceID        string
	SoftwareVersion string
	HardwareVersion string
	// Provenance fields per architecture doc 07-identity-provenance.md
	// (source=ebus_standard.identification, catalog version, decode
	// validity).
	Source         string
	CatalogVersion string
	Valid          bool
}

// DeviceInfo is the existing registry identity tuple used for provenance
// comparison. The shape is intentionally local — the provider does not
// import the registry package to keep catalog/ebus_standard self-contained
// (see AC6 namespace isolation). Callers that already have a
// registry.DeviceInfo translate into this struct at the boundary.
type DeviceInfo struct {
	Manufacturer    string
	DeviceID        string
	SoftwareVersion string
	HardwareVersion string
}

// CallerContext tags an Invoke call with the intent of the caller. In M3
// it only affects future whitelist expansion of mutating / destructive
// classes; currently every mutating class is denied regardless of context.
type CallerContext int

// CallerContext values.
const (
	// CallerContextUserFacing is the default. It denies every safety class
	// outside read_only_safe / read_only_bus_load.
	CallerContextUserFacing CallerContext = iota
	// CallerContextSystemNMRuntime is reserved for a future whitelist
	// expansion (network-management runtime). In M3 it is carried but NOT
	// interpreted — callers passing it receive the same default-deny as
	// user_facing.
	CallerContextSystemNMRuntime
)

// toInternal converts a public CallerContext into the safetypolicy.Caller
// tag. The enum values are mirrored deliberately so the translation is a
// pure literal switch.
func (c CallerContext) toInternal() safetypolicy.Caller {
	switch c {
	case CallerContextSystemNMRuntime:
		return safetypolicy.CallerSystemNMRuntime
	default:
		return safetypolicy.CallerUserFacing
	}
}

// Sentinel errors returned by the provider.
var (
	// ErrProviderDisabled is returned by every provider entrypoint when
	// EBUS_STANDARD_PROVIDER_ENABLED resolves to a disabled value at
	// construction time.
	ErrProviderDisabled = errors.New("ebus_standard: provider disabled via EBUS_STANDARD_PROVIDER_ENABLED")

	// ErrSafetyClassDenied is returned at the Invoke boundary when the
	// catalog's safety_class for the requested method is not permitted for
	// the caller context.
	ErrSafetyClassDenied = errors.New("ebus_standard: safety_class denied at invoke boundary")

	// ErrUnknownMethod is returned when Invoke receives a method ID that
	// does not exist in the loaded catalog.
	ErrUnknownMethod = errors.New("ebus_standard: unknown method id")

	// ErrDuplicateMethodID is returned by NewProvider when two catalog
	// commands share the same ID. This is distinct from
	// ErrDuplicateIdentityKey (loader-level, identity-fingerprint based): a
	// catalog may pass loader validation yet still contain two commands
	// with identical string IDs (e.g. YAML typo or merge mistake). Without
	// this check the second entry would silently overwrite the first in
	// the method index and Invoke would dispatch the wrong command.
	ErrDuplicateMethodID = errors.New("ebus_standard: duplicate catalog method id")

	// ErrEmptyMethodID is returned by NewProvider when a catalog command
	// has an empty string ID. Without this check the bad entry would be
	// indexed under "" and every normal method-name lookup would still
	// succeed, while Invoke("") would silently dispatch the malformed
	// command. LoadCatalog does not currently enforce non-empty IDs, so
	// the provider converts this silent-correctness failure mode into a
	// loud configuration error at construction time.
	ErrEmptyMethodID = errors.New("ebus_standard: empty catalog method id")
)

// DisableEnvVar is the exact env-var name read at construction time.
const DisableEnvVar = "EBUS_STANDARD_PROVIDER_ENABLED"

// NewProvider constructs a provider from a loaded catalog, with the
// enabled flag passed explicitly. Use NewProviderFromEnv to read the env
// var. The method index is built eagerly so Invoke lookups are O(1).
//
// Returns ErrDuplicateMethodID wrapped with the colliding command names
// when two catalog commands share the same ID. This is distinct from the
// loader-level ErrDuplicateIdentityKey (which fingerprints on the identity
// tuple). Two commands with identical IDs but different identity keys
// would pass the loader yet silently overwrite each other in the method
// index; this check converts that silent-correctness failure mode into a
// loud configuration error at construction time.
func NewProvider(cat Catalog, enabled bool) (*Provider, error) {
	idx := make(map[string]Command, len(cat.Services)*4)
	for _, svc := range cat.Services {
		for _, cmd := range svc.Commands {
			if cmd.ID == "" {
				return nil, fmt.Errorf("%w: service=%q command=%q",
					ErrEmptyMethodID, svc.Name, cmd.Name)
			}
			if prev, dup := idx[cmd.ID]; dup {
				return nil, fmt.Errorf("%w: id=%q first=%q second=%q",
					ErrDuplicateMethodID, cmd.ID, prev.Name, cmd.Name)
			}
			idx[cmd.ID] = cmd
		}
	}
	return &Provider{
		catalog:   cat,
		methodIdx: idx,
		enabled:   enabled,
	}, nil
}

// NewProviderFromEnv constructs a provider reading the disable switch from
// the process environment at this point in time. The env value is captured
// by copy; subsequent mutations to the environment do not affect the
// provider. An unset or empty variable defaults to enabled (blast-radius
// principle: the generic provider participates in the platform by default;
// operators explicitly opt out by setting "0" or "false").
//
// Propagates ErrDuplicateMethodID from NewProvider when the catalog has
// colliding command IDs.
func NewProviderFromEnv(cat Catalog) (*Provider, error) {
	return NewProvider(cat, envEnabled(os.Getenv(DisableEnvVar)))
}

// envEnabled applies the disable-switch contract. "0" / "false" (any case)
// => disabled. Empty / anything else => enabled.
func envEnabled(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "0", "false":
		return false
	default:
		return true
	}
}

// IsEnabled reports whether the provider is currently enabled.
func (p *Provider) IsEnabled() bool { return p != nil && p.enabled }

// Invoke dispatches a catalog method by ID after consulting the
// invoke-boundary safety gate. Error precedence:
//
//  1. ErrProviderDisabled — provider constructed in disabled state.
//  2. ErrUnknownMethod — method ID not present in the catalog.
//  3. ErrSafetyClassDenied — safety_class not permitted for callerCtx.
//  4. (future) decode / transport errors from the command body.
//
// M3 scope does not yet execute the wire request — the method signature is
// generated from the catalog and dispatch returns a stub result map for
// gated (allowed) methods. Decode into L7 primitives lands when the M1
// helianthus-ebusgo/protocol/ebus_standard/types package is published and
// the go.mod is bumped.
func (p *Provider) Invoke(ctx context.Context, methodID string, params map[string]any, caller CallerContext) (map[string]any, error) {
	if !p.IsEnabled() {
		return nil, ErrProviderDisabled
	}
	cmd, ok := p.methodIdx[methodID]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownMethod, methodID)
	}
	if !safetypolicy.Allow(safetypolicy.Class(cmd.SafetyClass), caller.toInternal()) {
		return nil, fmt.Errorf("%w: method=%q class=%s", ErrSafetyClassDenied, methodID, cmd.SafetyClass)
	}
	// M3: returning a minimal descriptor that proves the gate passed.
	// Downstream milestones wire the actual decode/encode once the L7
	// types package is available.
	_ = ctx
	_ = params
	return map[string]any{
		"method_id":       cmd.ID,
		"safety_class":    string(cmd.SafetyClass),
		"catalog_version": p.catalog.Version,
		"executed":        false, // M3 stub — handler body not yet wired
	}, nil
}

// Identification returns the canonical 0x07 0x04 identification descriptor
// wrapped in provenance metadata. It never mutates DeviceInfo; see
// CompareIdentity for the provenance merge helper.
//
// M3 scope: returns an empty descriptor with provenance labels populated.
// The wire decode lands when the gateway transport layer passes an
// observed Identification frame through to the provider in a later
// milestone.
func (p *Provider) Identification(ctx context.Context) (IdentificationDescriptor, error) {
	_ = ctx
	if !p.IsEnabled() {
		return IdentificationDescriptor{}, ErrProviderDisabled
	}
	return IdentificationDescriptor{
		Source:         string(SourceEBUSStandardIdent),
		CatalogVersion: p.catalog.Version,
		Valid:          false, // no wire frame yet in M3
	}, nil
}

// IdentitySource tags a single evidence contribution for the provenance
// comparison helper.
type IdentitySource string

// Evidence source labels per docs/architecture/ebus_standard/07-identity-
// provenance.md.
const (
	SourceDeviceInfo         IdentitySource = "device_info"
	SourceEBUSStandardIdent  IdentitySource = "ebus_standard.identification"
	SourceOperatorSeed       IdentitySource = "operator_seed"
	SourcePassiveObservation IdentitySource = "passive_observation"
)

// FieldProvenance is the per-field merge result produced by
// CompareIdentity. When Agreement is false, both Preferred and Sources
// remain populated.
type FieldProvenance struct {
	Preferred string
	Sources   []IdentitySourceValue
	Agreement bool
}

// IdentitySourceValue is a single labeled contribution.
type IdentitySourceValue struct {
	Source IdentitySource
	Value  string
	Valid  bool
}

// ProvenanceRecord is the full per-field provenance for a device.
type ProvenanceRecord struct {
	Manufacturer    FieldProvenance
	DeviceID        FieldProvenance
	SoftwareVersion FieldProvenance
	HardwareVersion FieldProvenance
}

// CompareIdentity merges an existing DeviceInfo with an ebus_standard
// Identification descriptor. The existing DeviceInfo is never mutated; on
// disagreement BOTH values are retained with source labels. Deterministic
// precedence (see 07-identity-provenance.md §"Deterministic Precedence"):
// device_info > ebus_standard.identification > operator_seed > passive >
// "unknown". This function implements steps 1-2 (the two sources it has
// evidence for); operator_seed and passive sources are stacked by the
// caller via future helpers.
func CompareIdentity(existing DeviceInfo, desc IdentificationDescriptor) ProvenanceRecord {
	mk := func(existingVal, descVal string) FieldProvenance {
		// Both sources are considered valid when they are non-empty. A
		// malformed descriptor value is the caller's responsibility to tag
		// as invalid before passing it in (via the Valid flag on the
		// descriptor as a whole — in M3 that flag applies uniformly to the
		// four fields; future milestones may introduce per-field validity).
		var srcs []IdentitySourceValue
		if existingVal != "" {
			srcs = append(srcs, IdentitySourceValue{
				Source: SourceDeviceInfo, Value: existingVal, Valid: true,
			})
		}
		if descVal != "" {
			srcs = append(srcs, IdentitySourceValue{
				Source: SourceEBUSStandardIdent, Value: descVal, Valid: desc.Valid,
			})
		}
		// Preferred value: deterministic precedence => device_info wins.
		// When existingVal is empty, the descriptor value is only promoted
		// to Preferred if the descriptor is marked valid. An invalid
		// descriptor (desc.Valid == false) must NOT surface malformed or
		// failed-decode identity as the canonical preferred value; leave
		// Preferred empty so the caller falls through to the next source
		// in the deterministic precedence chain (operator_seed > passive >
		// "unknown"). The invalid value is still retained in Sources with
		// Valid=false so the evidence trail is preserved.
		preferred := existingVal
		if preferred == "" && desc.Valid {
			preferred = descVal
		}
		// Agreement requires BOTH sources present AND equal. If only one
		// side provides evidence, treat it as agreement (there is nothing
		// to disagree with).
		agreement := existingVal == "" || descVal == "" || existingVal == descVal
		return FieldProvenance{
			Preferred: preferred,
			Sources:   srcs,
			Agreement: agreement,
		}
	}
	return ProvenanceRecord{
		Manufacturer:    mk(existing.Manufacturer, desc.Manufacturer),
		DeviceID:        mk(existing.DeviceID, desc.DeviceID),
		SoftwareVersion: mk(existing.SoftwareVersion, desc.SoftwareVersion),
		HardwareVersion: mk(existing.HardwareVersion, desc.HardwareVersion),
	}
}
