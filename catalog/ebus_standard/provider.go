package ebus_standard_catalog

import (
	"context"
	"errors"
)

// Provider is the generic ebus_standard L7 provider. Methods are generated
// at construction time from the loaded catalog; no per-device hard-coding
// exists in this package.
//
// Namespace isolation (M3 AC6): this provider and its internal helpers MUST
// NOT import any Vaillant-specific package under github.com/Project-
// Helianthus/helianthus-ebusreg/vaillant/.... Compile-time enforcement is
// provided by placing the safety gate in catalog/ebus_standard/internal/
// safetypolicy (Go's internal/ visibility rule excludes packages outside
// catalog/ebus_standard/). A second line of defense is the static-import
// regression test in provider_namespace_isolation_test.go which parses every
// Go file in this package and fails if a Vaillant import appears.
//
// Disable switch (M3 AC4): the env variable EBUS_STANDARD_PROVIDER_ENABLED
// is read once at provider construction via NewProviderFromEnv. When the
// value is literally "0" or "false" (case-insensitive), the provider is
// constructed in the disabled state and every call to Invoke /
// Identification returns ErrProviderDisabled. Default is enabled.
type Provider struct {
	// scaffolding only — populated by RED stubs so the package compiles.
}

// IdentificationDescriptor carries the canonical manufacturer / device-id /
// software / hardware values produced by the 0x07 0x04 service, wrapped in
// provenance metadata (never mutates DeviceInfo).
type IdentificationDescriptor struct {
	Manufacturer    string
	DeviceID        string
	SoftwareVersion string
	HardwareVersion string
	// Provenance fields are populated per architecture doc 07-identity-
	// provenance.md (source=ebus_standard.identification, catalog version,
	// decode validity, timestamp).
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
)

// DisableEnvVar is the exact env-var name read at construction time.
const DisableEnvVar = "EBUS_STANDARD_PROVIDER_ENABLED"

// NewProvider constructs a provider from a loaded catalog, with the enabled
// flag passed explicitly. Use NewProviderFromEnv to read the env var.
//
// RED stub: returns a non-functional provider.
func NewProvider(cat Catalog, enabled bool) *Provider {
	_ = cat
	_ = enabled
	return &Provider{}
}

// NewProviderFromEnv constructs a provider reading the disable switch from
// the process environment at this point in time. The env value is captured
// by copy; subsequent mutations to the environment do not affect the
// provider.
//
// RED stub: not implemented.
func NewProviderFromEnv(cat Catalog) *Provider {
	_ = cat
	return &Provider{}
}

// IsEnabled reports whether the provider is currently enabled.
//
// RED stub.
func (p *Provider) IsEnabled() bool { return false }

// Invoke dispatches a catalog method by ID after consulting the safety
// class policy. See ErrSafetyClassDenied for the deny rules.
//
// RED stub: always returns "not implemented".
func (p *Provider) Invoke(ctx context.Context, methodID string, params map[string]any, caller CallerContext) (map[string]any, error) {
	_ = ctx
	_ = methodID
	_ = params
	_ = caller
	return nil, errors.New("ebus_standard: Invoke not implemented (RED)")
}

// Identification returns the canonical 0x07 0x04 identification descriptor
// wrapped in provenance metadata. It never mutates DeviceInfo.
//
// RED stub: not implemented.
func (p *Provider) Identification(ctx context.Context) (IdentificationDescriptor, error) {
	_ = ctx
	return IdentificationDescriptor{}, errors.New("ebus_standard: Identification not implemented (RED)")
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
// disagreement both values are retained with source labels. Deterministic
// precedence (see 07-identity-provenance.md §"Deterministic Precedence"):
// device_info > ebus_standard.identification > operator_seed > passive >
// "unknown". This function implements steps 1-2 of that order.
//
// RED stub.
func CompareIdentity(existing DeviceInfo, desc IdentificationDescriptor) ProvenanceRecord {
	_ = existing
	_ = desc
	return ProvenanceRecord{}
}
