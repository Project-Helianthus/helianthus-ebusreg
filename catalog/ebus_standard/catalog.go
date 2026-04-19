package ebus_standard_catalog

// Namespace is the fixed provider namespace for this catalog.
const Namespace = "ebus_standard"

// CanonicalPlanSHA256 is the SHA-256 of the locked canonical plan
// ebus-standard-l7-services-w16-26.locked/00-canonical.md. It pins the
// catalog to a specific plan revision.
const CanonicalPlanSHA256 = "9e0a29bb76d99f551904b05749e322aafd3972621858aa6d1acbe49b9ef37305"

// CatalogVersion is the explicit version tag embedded at compile time. It
// must be bumped whenever the canonical YAML data changes.
const CatalogVersion = "v1.0-locked"

// Parameter describes a single decoded field in a command's payload. The
// body of Parameter is intentionally minimal in M2; full L7-type references
// land in M3 when the generic provider wires decode into helianthus-ebusgo
// primitives.
type Parameter struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Description string `yaml:"description,omitempty"`
}

// Command is a single method entry in the catalog. It carries the full
// identity key plus human-facing metadata, safety class, and parameter
// lists.
type Command struct {
	ID          string      `yaml:"id"`
	Name        string      `yaml:"name"`
	Description string      `yaml:"description,omitempty"`
	Identity    IdentityKey `yaml:"identity"`
	SafetyClass SafetyClass `yaml:"safety_class"`
	Request     []Parameter `yaml:"request,omitempty"`
	Response    []Parameter `yaml:"response,omitempty"`
}

// Service groups commands by primary-byte service code.
//
// PB is pointer-typed so the YAML loader can distinguish an absent key
// (nil → ErrServiceMissingPB) from an explicit zero value (e.g. `pb: 0x00`
// is legitimate and must be accepted). This parallels the treatment of
// IdentityKey.PB/SB: a value-typed uint8 would silently accept schema
// typos because an omitted key deserializes as 0, which would then match
// any command identity whose pb is also 0x00 and defeat the
// service/identity mismatch check.
type Service struct {
	PB          *uint8    `yaml:"pb"`
	Name        string    `yaml:"name"`
	Description string    `yaml:"description,omitempty"`
	Commands    []Command `yaml:"commands"`
}

// PBValue returns the dereferenced PB byte, or 0 if PB is nil. Callers that
// need to distinguish "absent" from "explicit 0x00" must check the pointer
// field directly.
func (s Service) PBValue() uint8 {
	if s.PB == nil {
		return 0
	}
	return *s.PB
}

// Catalog is the root document produced by the YAML loader.
type Catalog struct {
	Namespace  string    `yaml:"namespace"`
	Version    string    `yaml:"version"`
	PlanSHA256 string    `yaml:"plan_sha256"`
	Services   []Service `yaml:"services"`

	// ContentSHA256 is the SHA-256 of the raw YAML bytes. It is populated
	// by the loader, not read from the YAML itself.
	ContentSHA256 string `yaml:"-"`
}
