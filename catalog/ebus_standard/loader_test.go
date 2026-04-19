package ebus_standard_catalog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLoadCatalog_DuplicateIdentityKey asserts that a YAML fixture with a
// planted duplicate 14-tuple identity key causes LoadCatalog to fail with
// ErrDuplicateIdentityKey. This is the canonical collision-detection
// guarantee from locked plan §3.
func TestLoadCatalog_DuplicateIdentityKey(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "collision_duplicate.yaml"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	_, err = LoadCatalog(data)
	if err == nil {
		t.Fatalf("LoadCatalog: expected ErrDuplicateIdentityKey, got nil")
	}
	if !errors.Is(err, ErrDuplicateIdentityKey) {
		t.Fatalf("LoadCatalog: expected ErrDuplicateIdentityKey, got %v", err)
	}
}

// TestLoadCatalog_AmbiguousLengthSelector asserts that two entries sharing
// (namespace, PB, SB, selector_decoder) with incompatible length_prefix_mode
// cause LoadCatalog to fail.
func TestLoadCatalog_AmbiguousLengthSelector(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "ambiguous_length_selector.yaml"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	_, err = LoadCatalog(data)
	if err == nil {
		t.Fatalf("LoadCatalog: expected ErrAmbiguousLengthSelector, got nil")
	}
	if !errors.Is(err, ErrAmbiguousLengthSelector) {
		t.Fatalf("LoadCatalog: expected ErrAmbiguousLengthSelector, got %v", err)
	}
}

// TestLoadCatalog_AmbiguousLengthSelector_NoneDecoder asserts that two
// entries sharing the on-wire identity axes with selector_decoder="none"
// but differing only by length_prefix_mode are still rejected as
// ambiguous. Without a selector branch, no decode-time disambiguation is
// possible, so the ambiguity detector must treat "none" as a bundling axis
// rather than skipping it.
func TestLoadCatalog_AmbiguousLengthSelector_NoneDecoder(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "ambiguous_length_selector_none_decoder.yaml"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	_, err = LoadCatalog(data)
	if err == nil {
		t.Fatalf("LoadCatalog: expected ErrAmbiguousLengthSelector, got nil")
	}
	if !errors.Is(err, ErrAmbiguousLengthSelector) {
		t.Fatalf("LoadCatalog: expected ErrAmbiguousLengthSelector, got %v", err)
	}
}

// TestLoadCatalog_MissingPB asserts that a YAML fixture omitting the `pb`
// key in the identity block is rejected with ErrIncompleteIdentityKey. The
// value 0x00 must NOT be accepted as a default; absence of the key is the
// only signal the loader has to distinguish "missing" from "explicit 0x00".
func TestLoadCatalog_MissingPB(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "missing_pb.yaml"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	_, err = LoadCatalog(data)
	if err == nil {
		t.Fatalf("LoadCatalog: expected ErrIncompleteIdentityKey, got nil")
	}
	if !errors.Is(err, ErrIncompleteIdentityKey) {
		t.Fatalf("LoadCatalog: expected ErrIncompleteIdentityKey, got %v", err)
	}
}

// TestLoadCatalog_MissingSB asserts the symmetric rejection when the `sb`
// key is absent.
func TestLoadCatalog_MissingSB(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "missing_sb.yaml"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	_, err = LoadCatalog(data)
	if err == nil {
		t.Fatalf("LoadCatalog: expected ErrIncompleteIdentityKey, got nil")
	}
	if !errors.Is(err, ErrIncompleteIdentityKey) {
		t.Fatalf("LoadCatalog: expected ErrIncompleteIdentityKey, got %v", err)
	}
}

// TestLoadCatalog_ExplicitZeroPBSB asserts that explicitly-set `pb: 0x00`
// and `sb: 0x00` are ACCEPTED — the value zero is legitimate when the YAML
// author wrote it out. Only absence of the key must be rejected.
func TestLoadCatalog_ExplicitZeroPBSB(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "explicit_zero_pb_sb.yaml"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	cat, err := LoadCatalog(data)
	if err != nil {
		t.Fatalf("LoadCatalog: expected success for explicit zero PB/SB, got %v", err)
	}
	if len(cat.Services) != 1 || len(cat.Services[0].Commands) != 1 {
		t.Fatalf("LoadCatalog: expected 1 service with 1 command, got %+v", cat.Services)
	}
	id := cat.Services[0].Commands[0].Identity
	if id.PB == nil || *id.PB != 0x00 {
		t.Fatalf("LoadCatalog: PB: expected explicit 0x00, got %v", id.PB)
	}
	if id.SB == nil || *id.SB != 0x00 {
		t.Fatalf("LoadCatalog: SB: expected explicit 0x00, got %v", id.SB)
	}
}

// TestLoadCatalog_UnknownEnumValue asserts that typos in any of the enum-
// typed identity-key axes are rejected at load time with
// ErrUnknownEnumValue — rather than silently accepted as opaque strings,
// which would break downstream matching without a deterministic error.
//
// One sub-test per axis, each starts from a valid template and replaces
// a single field with a deliberate typo.
func TestLoadCatalog_UnknownEnumValue(t *testing.T) {
	const validTemplate = `
namespace: ebus_standard
version: v1.0-locked
plan_sha256: 9e0a29bb76d99f551904b05749e322aafd3972621858aa6d1acbe49b9ef37305
services:
  - pb: 0x07
    name: System Data
    commands:
      - id: ebus_standard.enum_typo_test
        name: Enum typo probe
        identity:
          namespace: ebus_standard
          pb: 0x07
          sb: 0x04
          selector_path: ""
          telegram_class: {{TELEGRAM_CLASS}}
          direction: {{DIRECTION}}
          request_or_response_role: {{ROLE}}
          broadcast_or_addressed: {{ADDRESSING}}
          answer_policy: {{ANSWER_POLICY}}
          length_prefix_mode: {{LPM}}
          selector_decoder: none
          service_variant: enum_probe
          transport_capability_requirements: [master_slave]
          version: v1.0-locked
        safety_class: read_only_bus_load
`
	type axis struct {
		name        string
		typoedField string
		typoedValue string
	}
	baseline := map[string]string{
		"TELEGRAM_CLASS": "addressed",
		"DIRECTION":      "request",
		"ROLE":           "initiator",
		"ADDRESSING":     "addressed",
		"ANSWER_POLICY":  "answer_required",
		"LPM":            "none",
	}
	// Sanity: the baseline template itself must load successfully.
	t.Run("baseline_loads", func(t *testing.T) {
		yml := applyTemplate(validTemplate, baseline)
		if _, err := LoadCatalog([]byte(yml)); err != nil {
			t.Fatalf("baseline template: expected success, got %v", err)
		}
	})

	cases := []axis{
		{name: "telegram_class", typoedField: "TELEGRAM_CLASS", typoedValue: "addresed_typo"},
		{name: "direction", typoedField: "DIRECTION", typoedValue: "requset"},
		{name: "request_or_response_role", typoedField: "ROLE", typoedValue: "initator"},
		{name: "broadcast_or_addressed", typoedField: "ADDRESSING", typoedValue: "addresed"},
		{name: "answer_policy", typoedField: "ANSWER_POLICY", typoedValue: "answer_req"},
		{name: "length_prefix_mode", typoedField: "LPM", typoedValue: "fixd"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vals := make(map[string]string, len(baseline))
			for k, v := range baseline {
				vals[k] = v
			}
			vals[tc.typoedField] = tc.typoedValue
			yml := applyTemplate(validTemplate, vals)
			_, err := LoadCatalog([]byte(yml))
			if err == nil {
				t.Fatalf("axis %s with typo %q: expected ErrUnknownEnumValue, got nil", tc.name, tc.typoedValue)
			}
			if !errors.Is(err, ErrUnknownEnumValue) {
				t.Fatalf("axis %s with typo %q: expected ErrUnknownEnumValue, got %v", tc.name, tc.typoedValue, err)
			}
			// Error message must name the field and the bad value for
			// deterministic diagnosis.
			if !strings.Contains(err.Error(), tc.name) {
				t.Fatalf("axis %s: error %q does not name the field", tc.name, err.Error())
			}
			if !strings.Contains(err.Error(), tc.typoedValue) {
				t.Fatalf("axis %s: error %q does not include the bad value", tc.name, err.Error())
			}
		})
	}
}

func applyTemplate(tpl string, vals map[string]string) string {
	out := tpl
	for k, v := range vals {
		out = strings.ReplaceAll(out, fmt.Sprintf("{{%s}}", k), v)
	}
	return out
}

// TestLoadCatalog_ServiceWithoutCommands asserts that a YAML fixture whose
// service block deserializes with no commands (typo or omission of the
// `commands:` key) is rejected with ErrServiceMissingCommands. The baseline
// catalog satisfies len(commands) > 0 by construction.
func TestLoadCatalog_ServiceWithoutCommands(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "service_without_commands.yaml"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	_, err = LoadCatalog(data)
	if err == nil {
		t.Fatalf("LoadCatalog: expected ErrServiceMissingCommands, got nil")
	}
	if !errors.Is(err, ErrServiceMissingCommands) {
		t.Fatalf("LoadCatalog: expected ErrServiceMissingCommands, got %v", err)
	}
}

// TestLoadCatalog_EmbeddedBaselineLoads asserts that the embedded baseline
// catalog passes the empty-commands guard — i.e. every declared service has
// at least one command.
func TestLoadCatalog_EmbeddedBaselineLoads(t *testing.T) {
	if _, err := LoadCatalog(EmbeddedYAML()); err != nil {
		t.Fatalf("LoadCatalog(embedded): unexpected error: %v", err)
	}
}

// TestLoadCatalog_TypoNamespace asserts that an identity-key namespace that
// differs from the fixed Namespace constant (e.g. "ebus_standrad") is
// rejected at load time with ErrInvalidNamespace. Without this guard, the
// typo would be silently treated as a distinct identity and bypass the
// duplicate-14-tuple collision detector.
func TestLoadCatalog_TypoNamespace(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "typo_namespace.yaml"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	_, err = LoadCatalog(data)
	if err == nil {
		t.Fatalf("LoadCatalog: expected ErrInvalidNamespace, got nil")
	}
	if !errors.Is(err, ErrInvalidNamespace) {
		t.Fatalf("LoadCatalog: expected ErrInvalidNamespace, got %v", err)
	}
	// Error message must name the field and the bad value for
	// deterministic diagnosis.
	if !strings.Contains(err.Error(), "namespace") {
		t.Fatalf("error %q does not name the field", err.Error())
	}
	if !strings.Contains(err.Error(), "ebus_standrad") {
		t.Fatalf("error %q does not include the bad value", err.Error())
	}
}

// TestEmbeddedCatalog_SHAPinning asserts that the embedded catalog carries
// a ContentSHA256 value and that it matches the SHA of its raw bytes. This
// enforces the plan's "Catalog is SHA-pinned and version-tagged"
// requirement.
func TestEmbeddedCatalog_SHAPinning(t *testing.T) {
	cat := MustEmbeddedCatalog()
	if cat.Version == "" {
		t.Fatalf("embedded catalog: Version empty")
	}
	if cat.Version != CatalogVersion {
		t.Fatalf("embedded catalog: Version=%q, want %q", cat.Version, CatalogVersion)
	}
	if cat.PlanSHA256 != CanonicalPlanSHA256 {
		t.Fatalf("embedded catalog: PlanSHA256=%q, want %q", cat.PlanSHA256, CanonicalPlanSHA256)
	}
	if cat.ContentSHA256 == "" {
		t.Fatalf("embedded catalog: ContentSHA256 empty")
	}
	// The loader must compute ContentSHA256 from the raw YAML bytes, not
	// from any post-parse structure. We re-compute from EmbeddedYAML() and
	// assert equality.
	want := ComputeContentSHA256(EmbeddedYAML())
	if cat.ContentSHA256 != want {
		t.Fatalf("embedded catalog: ContentSHA256=%q, want %q", cat.ContentSHA256, want)
	}
}

// TestEmbeddedCatalog_IdentityKeyCompleteness walks every command in the
// embedded catalog and asserts that its 14-tuple identity key is fully
// populated per canonical §3.
func TestEmbeddedCatalog_IdentityKeyCompleteness(t *testing.T) {
	cat := MustEmbeddedCatalog()
	var failures int
	for _, svc := range cat.Services {
		for _, cmd := range svc.Commands {
			if !cmd.Identity.IsComplete() {
				t.Errorf("command %q: incomplete identity key: %+v", cmd.ID, cmd.Identity)
				failures++
			}
			if cmd.Identity.Namespace != Namespace {
				t.Errorf("command %q: namespace=%q, want %q", cmd.ID, cmd.Identity.Namespace, Namespace)
			}
		}
	}
	if failures > 0 {
		t.Fatalf("%d commands have incomplete identity keys", failures)
	}
}

// TestEmbeddedCatalog_ServiceCoverage asserts that every service required
// by the locked plan's first-delivery baseline is present.
func TestEmbeddedCatalog_ServiceCoverage(t *testing.T) {
	cat := MustEmbeddedCatalog()
	required := []uint8{0x03, 0x05, 0x07, 0x08, 0x09, 0x0F, 0xFE, 0xFF}
	present := make(map[uint8]bool, len(cat.Services))
	for _, svc := range cat.Services {
		present[svc.PB] = true
	}
	for _, pb := range required {
		if !present[pb] {
			t.Errorf("service PB=0x%02X missing from embedded catalog", pb)
		}
	}
}
