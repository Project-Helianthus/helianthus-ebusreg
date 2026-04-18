package ebus_standard_catalog

import (
	"errors"
	"os"
	"path/filepath"
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
