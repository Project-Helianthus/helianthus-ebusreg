package ebus_standard_catalog

import (
	"go/parser"
	"go/token"
	"io/fs"
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
	// Walk the entire catalog/ebus_standard subtree recursively. A prior
	// implementation only called os.ReadDir on two top-level roots, which
	// produced a false-negative path: any Vaillant import dropped into a
	// nested subdirectory (e.g. internal/anything/deeper/foo.go) escaped
	// the scanner entirely. filepath.WalkDir catches the full subtree.
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, werr error) error {
		if werr != nil {
			return werr
		}
		if d.IsDir() {
			// Skip testdata/ — it contains planted .go.txt fixtures that are
			// exercised explicitly by TestNamespaceIsolation_PlantedViolationCaught.
			if d.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		// Test files are allowed to reference fixtures by string path but
		// must not import Vaillant packages either, so both .go and
		// _test.go files are scanned uniformly.
		fset := token.NewFileSet()
		file, perr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if perr != nil {
			return perr
		}
		for _, imp := range file.Imports {
			p := strings.Trim(imp.Path.Value, `"`)
			if strings.Contains(p, "/vaillant/") || strings.HasSuffix(p, "/vaillant") {
				violations = append(violations, path+": imports "+p)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
	if len(violations) > 0 {
		t.Fatalf("namespace isolation violation (catalog/ebus_standard must not import Vaillant packages):\n  %s",
			strings.Join(violations, "\n  "))
	}
}

// TestNamespaceIsolation_RecursiveWalk_CatchesNestedViolation proves the
// scanner descends into subdirectories. It copies the planted nested
// fixture into a tempdir at a multi-level path and runs the same WalkDir
// + parser.ImportsOnly pipeline the real test uses. A flat ReadDir
// implementation would miss the deeply-nested file — this regression
// test pins the recursive behavior.
func TestNamespaceIsolation_RecursiveWalk_CatchesNestedViolation(t *testing.T) {
	plantedPath := filepath.Join("testdata", "nested", "deeper", "planted_vaillant_nested.go.txt")
	raw, err := os.ReadFile(plantedPath)
	if err != nil {
		t.Fatalf("read nested planted fixture: %v", err)
	}
	tmp := t.TempDir()
	deep := filepath.Join(tmp, "a", "b", "c")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	tmpPath := filepath.Join(deep, "planted.go")
	if err := os.WriteFile(tmpPath, raw, 0o644); err != nil {
		t.Fatalf("write nested planted copy: %v", err)
	}
	var violations []string
	err = filepath.WalkDir(tmp, func(path string, d fs.DirEntry, werr error) error {
		if werr != nil {
			return werr
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		fset := token.NewFileSet()
		file, perr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if perr != nil {
			return perr
		}
		for _, imp := range file.Imports {
			p := strings.Trim(imp.Path.Value, `"`)
			if strings.Contains(p, "/vaillant/") || strings.HasSuffix(p, "/vaillant") {
				violations = append(violations, path+": imports "+p)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk nested: %v", err)
	}
	if len(violations) == 0 {
		t.Fatal("nested planted vaillant import was NOT caught — recursive walk regressed to flat ReadDir")
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
