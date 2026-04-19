package ebus_standard_catalog

import (
	"bytes"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// updateABIGolden mirrors the `go test -update` convention used elsewhere
// in the Helianthus toolchain. When true, the test overwrites the golden
// fixture instead of asserting equality. The flag complements UPDATE=1 via
// the OS environment; either mechanism works.
var updateABIGolden = flag.Bool("update", false, "update provider ABI golden fixture")

// TestProvider_ABISnapshot is the CONTRACT_ABI_SNAPSHOT_TEST (M3 AC5).
// It parses every non-test .go file in this package (and the internal/
// safetypolicy helper) and serializes the exported-symbol shape (type
// declarations, method receivers + signatures, and top-level function
// signatures). On a diff vs the golden fixture the test fails with a
// directive to update the fixture AND add a PR-body rationale.
func TestProvider_ABISnapshot(t *testing.T) {
	got := mustComputeABISnapshot(t)

	goldenPath := filepath.Join("testdata", "provider_abi.golden.txt")
	if *updateABIGolden || os.Getenv("UPDATE") == "1" {
		if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		t.Logf("golden fixture updated at %s", goldenPath)
		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden %s: %v (run `go test ./catalog/ebus_standard -update` or UPDATE=1 go test to create it, then add a PR-body rationale for the ABI change)", goldenPath, err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("provider ABI drift detected.\n"+
			"Directive: run `go test ./catalog/ebus_standard -update` (or UPDATE=1 go test) to refresh\n"+
			"  testdata/provider_abi.golden.txt, THEN add a PR-body rationale explaining the ABI\n"+
			"  change. Unapproved ABI drift is a merge blocker per canonical plan M3 §ADD-05.\n\n"+
			"diff (got vs want):\n%s", unifiedDiff(string(want), string(got)))
	}
}

// mustComputeABISnapshot walks the source tree rooted at the ebus_standard
// catalog package and emits a deterministic textual representation of every
// exported API symbol.
func mustComputeABISnapshot(t *testing.T) []byte {
	t.Helper()
	var lines []string
	// Package roots to cover. Relative paths are resolved against the
	// package directory (where `go test` runs).
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
			if strings.HasSuffix(ent.Name(), "_test.go") {
				continue
			}
			path := filepath.Join(root, ent.Name())
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				t.Fatalf("parse %s: %v", path, err)
			}
			pkgLabel := file.Name.Name
			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.GenDecl:
					for _, spec := range d.Specs {
						switch s := spec.(type) {
						case *ast.TypeSpec:
							if !s.Name.IsExported() {
								continue
							}
							lines = append(lines, formatTypeDecl(pkgLabel, s, fset))
						case *ast.ValueSpec:
							// Exported constants and vars (errors, sentinels).
							for _, n := range s.Names {
								if !n.IsExported() {
									continue
								}
								kind := "var"
								if d.Tok == token.CONST {
									kind = "const"
								}
								lines = append(lines, pkgLabel+" "+kind+" "+n.Name)
							}
						}
					}
				case *ast.FuncDecl:
					if !d.Name.IsExported() {
						continue
					}
					lines = append(lines, formatFuncDecl(pkgLabel, d, fset))
				}
			}
		}
	}
	sort.Strings(lines)
	var buf bytes.Buffer
	for _, ln := range lines {
		buf.WriteString(ln)
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

func formatTypeDecl(pkg string, s *ast.TypeSpec, fset *token.FileSet) string {
	var buf bytes.Buffer
	buf.WriteString(pkg)
	buf.WriteString(" type ")
	buf.WriteString(s.Name.Name)
	buf.WriteString(" = ")
	_ = printNode(&buf, fset, s.Type)
	return compactWhitespace(buf.String())
}

func formatFuncDecl(pkg string, d *ast.FuncDecl, fset *token.FileSet) string {
	var buf bytes.Buffer
	buf.WriteString(pkg)
	buf.WriteString(" func ")
	if d.Recv != nil && len(d.Recv.List) == 1 {
		buf.WriteString("(")
		_ = printNode(&buf, fset, d.Recv.List[0].Type)
		buf.WriteString(") ")
	}
	buf.WriteString(d.Name.Name)
	_ = printNode(&buf, fset, d.Type)
	return compactWhitespace(buf.String())
}

func printNode(buf *bytes.Buffer, fset *token.FileSet, node ast.Node) error {
	// Use a minimal printer to avoid pulling go/format (which would pin us to
	// a specific gofmt revision). The ast package's default String via
	// position-aware extraction suffices for deterministic output across Go
	// toolchain revisions.
	start := fset.Position(node.Pos()).Offset
	end := fset.Position(node.End()).Offset
	file := fset.File(node.Pos())
	if file == nil {
		return nil
	}
	src, err := os.ReadFile(file.Name())
	if err != nil {
		return err
	}
	if start < 0 || end > len(src) {
		return nil
	}
	buf.Write(src[start:end])
	return nil
}

func compactWhitespace(s string) string {
	var b strings.Builder
	prevSpace := false
	for _, r := range s {
		switch r {
		case ' ', '\t', '\n', '\r':
			if !prevSpace {
				b.WriteRune(' ')
				prevSpace = true
			}
		default:
			b.WriteRune(r)
			prevSpace = false
		}
	}
	return strings.TrimSpace(b.String())
}

// unifiedDiff produces a minimal line-diff. Pulling in a full diff library
// would be overkill for a fixture compare; this prints both sides line-by-
// line so the operator can eyeball the drift.
func unifiedDiff(a, b string) string {
	aLines := strings.Split(a, "\n")
	bLines := strings.Split(b, "\n")
	var buf bytes.Buffer
	buf.WriteString("--- want\n+++ got\n")
	max := len(aLines)
	if len(bLines) > max {
		max = len(bLines)
	}
	for i := 0; i < max; i++ {
		var av, bv string
		if i < len(aLines) {
			av = aLines[i]
		}
		if i < len(bLines) {
			bv = bLines[i]
		}
		if av != bv {
			if av != "" {
				buf.WriteString("- ")
				buf.WriteString(av)
				buf.WriteByte('\n')
			}
			if bv != "" {
				buf.WriteString("+ ")
				buf.WriteString(bv)
				buf.WriteByte('\n')
			}
		}
	}
	return buf.String()
}
