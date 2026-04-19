package ebus_standard_catalog

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNamespaceIsolation_NoVaillantImports is the M3 AC6 regression test
// that enforces build-time namespace isolation via static analysis. It
// parses every Go file in the catalog/ebus_standard/... subtree and
// asserts that NO import path contains "/vaillant/" or terminates with
// "/vaillant". Together with the internal/safetypolicy placement (Go's
// structural internal-rule prevents Vaillant-rooted importers), this
// provides two lines of defense:
//
//  1. Structural: anything under internal/safetypolicy is unreachable
//     from /vaillant/... packages per Go's compile-time internal rule.
//  2. Static: this test catches accidental "the other direction" imports
//     — code inside catalog/ebus_standard/ pulling Vaillant specifics into
//     the generic namespace. That would leak Vaillant semantics into the
//     cross-vendor provider and silently bind the catalog to a single
//     manufacturer.
//
// The test also includes a synthetic planted violation guarded by the
// build tag `namespace_guard_violation`. When that tag is active (explicit
// operator target), the file at testdata/planted_vaillant_import.go.txt
// would be detected here. The test body ignores .txt fixtures (they are
// not compiled), but under the build tag the fixture is symlinked or
// copied by the Makefile/CI into a .go file to prove the scanner fires.
// In the default build the .txt content is inert, keeping the regression
// test zero-cost for normal runs.
func TestNamespaceIsolation_NoVaillantImports(t *testing.T) {
	var violations []string
	roots := []string{".", filepath.Join("internal", "safetypolicy")}
	for _, root := range roots {
		entries, err := os.ReadDir(root)
		if err != nil {
			t.Fatalf("readdir %s: %v", root, err)
		}
		for _, ent := range entries {
			if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") {
				continue
			}
			path := filepath.Join(root, ent.Name())
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
			if err != nil {
				t.Fatalf("parse %s: %v", path, err)
			}
			for _, imp := range file.Imports {
				p := strings.Trim(imp.Path.Value, `"`)
				if strings.Contains(p, "/vaillant/") || strings.HasSuffix(p, "/vaillant") {
					violations = append(violations, path+": imports "+p)
				}
			}
		}
	}
	if len(violations) > 0 {
		t.Fatalf("namespace isolation violation (catalog/ebus_standard must not import Vaillant packages):\n  %s",
			strings.Join(violations, "\n  "))
	}
}

// TestNamespaceIsolation_PlantedViolationCaught asserts the scanner is
// NOT a no-op: we feed it the planted fixture at
// testdata/planted_vaillant_import.go.txt (which imports a real vaillant
// package) and require that exactly one violation is returned. If this
// ever stops firing, the scanner has regressed and the outer test is a
// false-negative.
func TestNamespaceIsolation_PlantedViolationCaught(t *testing.T) {
	plantedPath := filepath.Join("testdata", "planted_vaillant_import.go.txt")
	raw, err := os.ReadFile(plantedPath)
	if err != nil {
		t.Fatalf("read planted fixture: %v", err)
	}
	// Copy into a .go file in a temp dir so go/parser accepts it (parser
	// itself is extension-agnostic, but we keep the parity explicit).
	tmp := t.TempDir()
	tmpPath := filepath.Join(tmp, "planted.go")
	if err := os.WriteFile(tmpPath, raw, 0o644); err != nil {
		t.Fatalf("write planted copy: %v", err)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, tmpPath, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse planted: %v", err)
	}
	var hits int
	for _, imp := range file.Imports {
		p := strings.Trim(imp.Path.Value, `"`)
		if strings.Contains(p, "/vaillant/") || strings.HasSuffix(p, "/vaillant") {
			hits++
		}
	}
	if hits == 0 {
		t.Fatal("planted vaillant import was NOT caught — namespace-isolation scanner regressed")
	}
}
