package ebus_standard_catalog

import (
	"bytes"
	"flag"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
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
	// Use the default build context so files behind build tags that are
	// inactive in the current build (e.g. loader_tinygo.go under `tinygo`
	// while we're running the non-tinygo default) are excluded. Without
	// this filter, mutually-exclusive tagged files were parsed together
	// and produced duplicate symbols in the golden fixture (observed:
	// ComputeContentSHA256 appearing twice from embedded.go + its tinygo
	// counterpart). go/build.Context.MatchFile applies the same logic the
	// compiler uses to decide which files belong in the package build.
	bctx := build.Default
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
			match, merr := bctx.MatchFile(root, ent.Name())
			if merr != nil {
				t.Fatalf("build.MatchFile %s/%s: %v", root, ent.Name(), merr)
			}
			if !match {
				continue
			}
			path := filepath.Join(root, ent.Name())
			fset := token.NewFileSet()
			// Parse WITHOUT parser.ParseComments: comments are irrelevant
			// to the exported-API shape. Excluding them at parse time
			// keeps the downstream AST serialization free of comment
			// nodes and prevents comment-only edits from causing ABI
			// golden drift (pre-fix behavior: raw-src slicing captured
			// any comments that happened to fall inside a node's Pos..End
			// range, causing false-positive CI failures).
			file, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				t.Fatalf("parse %s: %v", path, err)
			}
			// Detach any parsed comment groups defensively.
			file.Comments = nil
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

// TestABISnapshot_CommentsInsensitive is the regression gate for the
// go/printer-based serialization: it feeds the snapshot pipeline two
// variants of the same exported type — one bare, one decorated with doc
// comments and inline comments inside a struct field — and asserts the
// formatter produces byte-identical output for both. A regression that
// re-introduces raw-source slicing would make the commented variant
// capture comment bytes and diverge.
func TestABISnapshot_CommentsInsensitive(t *testing.T) {
	bare := `package sample

type ExportedShape struct {
	Field1 int
	Field2 string
}

func ExportedFunc(a int) error { return nil }
`
	commented := `package sample

// ExportedShape is the documented shape.
//
// Multiple paragraphs of doc comments MUST NOT leak into the ABI snapshot.
type ExportedShape struct {
	// Field1 is the first field.
	Field1 int // trailing comment on Field1
	// Field2 is the second field.
	Field2 string
}

// ExportedFunc has a lengthy doc comment that describes
// every edge case in excruciating detail.
func ExportedFunc(a int) error { return nil }
`
	render := func(src string) string {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "sample.go", src, 0)
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		file.Comments = nil
		var lines []string
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					if s, ok := spec.(*ast.TypeSpec); ok && s.Name.IsExported() {
						lines = append(lines, formatTypeDecl(file.Name.Name, s, fset))
					}
				}
			case *ast.FuncDecl:
				if d.Name.IsExported() {
					lines = append(lines, formatFuncDecl(file.Name.Name, d, fset))
				}
			}
		}
		sort.Strings(lines)
		return strings.Join(lines, "\n")
	}
	if got, want := render(commented), render(bare); got != want {
		t.Fatalf("ABI snapshot leaked comment bytes.\n  bare=%q\n  comm=%q", want, got)
	}
}

// TestABISnapshot_PreservesTypeDeclForm asserts that formatTypeDecl
// distinguishes defined types from aliases and renders type parameters
// for generics. These are ABI-significant shape differences: a
// defined→alias flip changes method-set reachability; a generic type
// param add/remove changes the exported surface. The snapshot MUST
// catch both.
func TestABISnapshot_PreservesTypeDeclForm(t *testing.T) {
	render := func(src string) []string {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "sample.go", src, 0)
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		file.Comments = nil
		var lines []string
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range gd.Specs {
				if s, ok := spec.(*ast.TypeSpec); ok && s.Name.IsExported() {
					lines = append(lines, formatTypeDecl(file.Name.Name, s, fset))
				}
			}
		}
		sort.Strings(lines)
		return lines
	}

	// Defined type: no `=`.
	defined := render(`package sample
type Defined string
`)
	if len(defined) != 1 || !strings.Contains(defined[0], "type Defined string") {
		t.Fatalf("defined type render wrong: %q", defined)
	}
	if strings.Contains(defined[0], "Defined = string") {
		t.Fatalf("defined type rendered as alias: %q", defined[0])
	}

	// Alias type: `=` present.
	alias := render(`package sample
type Alias = string
`)
	if len(alias) != 1 || !strings.Contains(alias[0], "type Alias = string") {
		t.Fatalf("alias type render wrong: %q", alias)
	}

	// Defined vs alias MUST produce distinct snapshot lines. This is the
	// core regression: the pre-fix implementation hardcoded `" = "` and
	// collapsed the two forms.
	definedLine := render(`package sample
type T string
`)
	aliasLine := render(`package sample
type T = string
`)
	if definedLine[0] == aliasLine[0] {
		t.Fatalf("defined and alias rendered identically: %q", definedLine[0])
	}

	// Generic type: type params present.
	generic := render(`package sample
type Generic[T any] struct{ Val T }
`)
	if len(generic) != 1 || !strings.Contains(generic[0], "Generic[T any]") {
		t.Fatalf("generic type params missing: %q", generic)
	}

	// Adding/removing a type parameter MUST produce a snapshot diff.
	oneParam := render(`package sample
type G[T any] struct{ V T }
`)
	twoParams := render(`package sample
type G[T any, U comparable] struct{ V T }
`)
	if oneParam[0] == twoParams[0] {
		t.Fatalf("type-param count change not detected: %q", oneParam[0])
	}
}

func formatTypeDecl(pkg string, s *ast.TypeSpec, fset *token.FileSet) string {
	var buf bytes.Buffer
	buf.WriteString(pkg)
	buf.WriteString(" type ")
	buf.WriteString(s.Name.Name)
	// Emit type parameters for generic type declarations (Go 1.18+):
	// `type Generic[T any, U comparable] struct{...}`. TypeParams is nil
	// for non-generic types, so guard the render.
	if s.TypeParams != nil && len(s.TypeParams.List) > 0 {
		buf.WriteString("[")
		for i, field := range s.TypeParams.List {
			if i > 0 {
				buf.WriteString(", ")
			}
			for j, name := range field.Names {
				if j > 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(name.Name)
			}
			if len(field.Names) > 0 {
				buf.WriteString(" ")
			}
			_ = printNode(&buf, fset, field.Type)
		}
		buf.WriteString("]")
	}
	// Distinguish a type alias (`type T = U`, TypeSpec.Assign != NoPos)
	// from a defined type (`type T U`). The two forms are ABI-distinct:
	// flipping one to the other changes method-set reachability and
	// conversion rules, so the snapshot MUST preserve the separator.
	if s.Assign != token.NoPos {
		buf.WriteString(" = ")
	} else {
		buf.WriteString(" ")
	}
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
	// Serialize the AST node canonically via go/printer with mode=0 (no
	// raw-source/comment passthrough). This replaces a prior
	// implementation that sliced raw source between node.Pos() and
	// node.End(): that approach silently captured inline comments sitting
	// inside struct literals and type declarations, which caused the
	// golden fixture to drift on comment-only edits and produced
	// false-positive ABI-gate failures in CI. Serializing the node itself
	// is comments-insensitive by construction (comments live on
	// *ast.File.Comments, not on the type/func node, and are additionally
	// scrubbed by the caller via file.Comments = nil defensively).
	cfg := printer.Config{Mode: 0, Tabwidth: 8}
	return cfg.Fprint(buf, fset, node)
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
