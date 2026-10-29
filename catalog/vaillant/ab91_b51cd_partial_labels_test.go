package vaillant

import (
	ab91Context "context"
	ab91Errors "errors"
	ab91Standard "github.com/Project-Helianthus/helianthus-ebusreg/catalog/ebus_standard"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

func ab91Pointer[ab91Type any](ab91Value ab91Type) *ab91Type { return &ab91Value }

func ab91Catalog(ab91T *testing.T) Ab91PartialLabels {
	ab91T.Helper()
	ab91Catalog, ab91Error := Ab91B51CDPartialLabels()
	if ab91Error != nil {
		ab91T.Fatal(ab91Error)
	}
	return ab91Catalog
}

func TestAb91ExactPartialLabels(ab91T *testing.T) {
	ab91Catalog := ab91Catalog(ab91T)
	ab91WantCatalog := Ab91PartialLabels{
		Ab91Representation: "Proven", Ab91PhysicalQualification: "Unknown", Ab91InstalledSupport: "Unknown",
		Ab91Equivalence: "Unknown", Ab91BodyLayout: "Unknown", Ab91Codec: "Unknown", Ab91Gating: "Unknown",
		Ab91Confirmation: "Unknown", Ab91Effect: "Unknown", Ab91PhysicalFraming: "Unknown", Ab91Controls: "opaque", Ab91Completeness: "incomplete",
		Ab91Limitations: []string{
			"Lexical suffixes identify labels only; they encode no secondary, operation or body byte.",
			"Logical byte and bit positions are not physical framing offsets.",
			"Partial forms are not empty encodable packets and establish no enum, constructor or decoder.",
			"No invocation, emitter, registration, runtime cache or transport is provided.",
			"Synthetic metadata representation establishes no physical support, freshness, persistence or execution permission.",
		},
	}
	ab91Labels := map[string][]string{
		"B51C": {"00", "01", "02", "03", "10", "13", "14", "19", "1B", "1C", "1D"},
		"B51D": {"00", "01", "02", "03", "10", "11", "12", "13", "14", "15", "16", "17"},
	}
	if len(ab91Catalog.Ab91Entries) != 46 {
		ab91T.Fatalf("ab91 got %d rows, want 46", len(ab91Catalog.Ab91Entries))
	}
	ab91Seen := make(map[string]bool)
	for ab91Variant, ab91Suffixes := range ab91Labels {
		for _, ab91Suffix := range ab91Suffixes {
			for ab91Direction, ab91Constant := range map[string]uint8{"request": 0xB5, "positive": 0xF5} {
				ab91Key := Ab91PartialKey{Ab91Family: "vaillant", Ab91Variant: ab91Variant, Ab91Label: ab91Variant + "_" + ab91Suffix, Ab91Primary: 0xB5, Ab91Direction: ab91Direction, Ab91Field: "coded_constant"}
				ab91T.Run("ab91_"+ab91Key.Ab91Label+"_"+ab91Direction, func(ab91T *testing.T) {
					ab91Row, ab91Found, ab91Error := Ab91LookupB51CDPartialLabel(ab91Key)
					if ab91Error != nil || !ab91Found {
						ab91T.Fatalf("ab91 lookup: %v %v", ab91Found, ab91Error)
					}
					ab91Want := Ab91PartialLabel{
						Ab91Key: ab91Key, Ab91Grade: "Hypothesis", Ab91RawFallback: "preserve full original carrier",
						Ab91Field: Ab91PartialField{Ab91LogicalByte: ab91Pointer(0), Ab91LogicalBit: ab91Pointer(0), Ab91WidthBits: ab91Pointer(8), Ab91Constant: ab91Pointer(ab91Constant), Ab91Signedness: "Unknown", Ab91ByteOrder: "Unknown", Ab91Conversion: "Unknown", Ab91Grade: "Hypothesis", Ab91RawFallback: "preserve original bytes without interpreting the body"},
					}
					if !reflect.DeepEqual(ab91Row, ab91Want) {
						ab91T.Fatalf("ab91 metadata mismatch: %#v", ab91Row)
					}
					ab91Matches := 0
					for _, ab91Entry := range ab91Catalog.Ab91Entries {
						if reflect.DeepEqual(ab91Entry, ab91Want) {
							ab91Matches++
						}
					}
					if ab91Matches != 1 {
						ab91T.Fatalf("ab91 want one independent row, got %d", ab91Matches)
					}
					for _, ab91Case := range []struct {
						ab91Raw  []byte
						ab91Want bool
					}{
						{[]byte{ab91Constant, 0x73, 0xA2}, true}, {[]byte{ab91Constant ^ 1, 0x73}, false}, {nil, false},
					} {
						ab91Original := append([]byte(nil), ab91Case.ab91Raw...)
						if ab91ConstantMatches(ab91Row.Ab91Field, ab91Case.ab91Raw) != ab91Case.ab91Want {
							ab91T.Fatal("ab91 synthetic partial constant check")
						}
						if !reflect.DeepEqual(ab91Original, ab91Case.ab91Raw) {
							ab91T.Fatal("ab91 opaque raw body changed")
						}
					}
				})
			}
			ab91Seen[ab91Variant+"_"+ab91Suffix] = true
		}
	}
	if len(ab91Seen) != 23 {
		ab91T.Fatal("ab91 label identities merged")
	}
	ab91Catalog.Ab91Entries = nil
	if !reflect.DeepEqual(ab91Catalog, ab91WantCatalog) {
		ab91T.Fatal("ab91 qualification or unknown-body contract changed")
	}
}

func ab91ConstantMatches(ab91Field Ab91PartialField, ab91Raw []byte) bool {
	if ab91Field.Ab91LogicalByte == nil || ab91Field.Ab91LogicalBit == nil || ab91Field.Ab91WidthBits == nil || ab91Field.Ab91Constant == nil {
		return false
	}
	ab91Position := *ab91Field.Ab91LogicalByte
	if ab91Position < 0 || ab91Position >= len(ab91Raw) {
		return false
	}
	ab91Mask := uint16(1)<<uint(*ab91Field.Ab91WidthBits) - 1
	return uint16(ab91Raw[ab91Position])>>uint(*ab91Field.Ab91LogicalBit)&ab91Mask == uint16(*ab91Field.Ab91Constant)
}

func TestAb91ExactIdentityAxes(ab91T *testing.T) {
	ab91Base := ab91Catalog(ab91T).Ab91Entries[0]
	ab91Qualified := ab91Base
	ab91Qualified.Ab91Key = Ab91PartialKey{Ab91Family: "vaillant", Ab91Variant: "B524", Ab91Label: "ab91_synthetic", Ab91Primary: 0xB5, Ab91Secondary: ab91Pointer(uint16(0x24)), Ab91Operation: ab91Pointer(uint16(1)), Ab91Direction: "request", Ab91Field: "ab91_synthetic", Ab91Opcode: ab91Pointer(uint16(2)), Ab91Group: ab91Pointer(uint16(3)), Ab91Instance: ab91Pointer(uint16(4)), Ab91Register: ab91Pointer(uint16(5))}
	for _, ab91Control := range []Ab91PartialLabel{ab91Base, ab91Qualified} {
		for _, ab91Axis := range []string{"Ab91Family", "Ab91Variant", "Ab91Label", "Ab91Primary", "Ab91Secondary", "Ab91Operation", "Ab91Direction", "Ab91Field", "Ab91Opcode", "Ab91Group", "Ab91Instance", "Ab91Register"} {
			ab91T.Run("ab91_"+ab91Control.Ab91Key.Ab91Variant+"_"+ab91Axis, func(ab91T *testing.T) {
				ab91Key := ab91Control.Ab91Key
				ab91Member := reflect.ValueOf(&ab91Key).Elem().FieldByName(ab91Axis)
				switch ab91Member.Kind() {
				case reflect.String:
					ab91Member.SetString("ab91_distinct")
				case reflect.Uint8:
					ab91Member.SetUint(0)
				case reflect.Pointer:
					ab91Member.Set(reflect.ValueOf(ab91Pointer(uint16(0))))
				default:
					ab91T.Fatal("ab91 untested axis")
				}
				if _, ab91Found := ab91FindPartialLabel([]Ab91PartialLabel{ab91Control}, ab91Key); ab91Found {
					ab91T.Fatal("ab91 distinct axis wildcard matched")
				}
				if ab91Control.Ab91Key.Ab91Variant != "B524" {
					if _, ab91Found, ab91Error := Ab91LookupB51CDPartialLabel(ab91Key); ab91Error != nil || ab91Found {
						ab91T.Fatal("ab91 production lookup ignored axis")
					}
				}
				ab91Distinct := ab91Control
				ab91Distinct.Ab91Key = ab91Key
				ab91Distinct.Ab91RawFallback = "ab91_distinct"
				ab91Rows := []Ab91PartialLabel{ab91Control, ab91Distinct}
				if ab91Got, ab91Found := ab91FindPartialLabel(ab91Rows, ab91Key); !ab91Found || !reflect.DeepEqual(ab91Got, ab91Distinct) {
					ab91T.Fatal("ab91 qualified positive control failed")
				}
				if ab91Got, ab91Found := ab91FindPartialLabel(ab91Rows, ab91Control.Ab91Key); !ab91Found || !reflect.DeepEqual(ab91Got, ab91Control) {
					ab91T.Fatal("ab91 original positive control failed")
				}
			})
		}
	}
	for _, ab91Axis := range []string{"Ab91Secondary", "Ab91Operation", "Ab91Opcode", "Ab91Group", "Ab91Instance", "Ab91Register"} {
		ab91T.Run("ab91_absent_zero_nonzero_"+ab91Axis, func(ab91T *testing.T) {
			for _, ab91Coordinate := range []*uint16{nil, ab91Pointer(uint16(0)), ab91Pointer(uint16(1))} {
				ab91Row := ab91Qualified
				reflect.ValueOf(&ab91Row.Ab91Key).Elem().FieldByName(ab91Axis).Set(reflect.ValueOf(ab91Coordinate))
				for _, ab91Query := range []*uint16{nil, ab91Pointer(uint16(0)), ab91Pointer(uint16(1))} {
					ab91Key := ab91Row.Ab91Key
					reflect.ValueOf(&ab91Key).Elem().FieldByName(ab91Axis).Set(reflect.ValueOf(ab91Query))
					_, ab91Found := ab91FindPartialLabel([]Ab91PartialLabel{ab91Row}, ab91Key)
					if ab91Found != reflect.DeepEqual(ab91Coordinate, ab91Query) {
						ab91T.Fatal("ab91 absent and numeric coordinates conflated")
					}
				}
			}
		})
	}
}

func ab91MutateNested(ab91Row *Ab91PartialLabel) {
	*ab91Row.Ab91Field.Ab91LogicalByte = 7
	*ab91Row.Ab91Field.Ab91LogicalBit = 6
	*ab91Row.Ab91Field.Ab91WidthBits = 3
	*ab91Row.Ab91Field.Ab91Constant = 0
	ab91Row.Ab91Key.Ab91Label = "ab91_mutated"
	ab91Row.Ab91Key.Ab91Secondary = ab91Pointer(uint16(0))
	ab91Row.Ab91Key.Ab91Operation = ab91Pointer(uint16(0))
	ab91Row.Ab91Field.Ab91Enumeration = map[string]string{"ab91_mutated": "ab91_mutated"}
	ab91Row.Ab91Field.Ab91Reserved = []int64{1}
	ab91Row.Ab91Field.Ab91Replacement = []int64{2}
	ab91Row.Ab91Field.Ab91Bounds = []int64{0, 1}
	ab91Row.Ab91BodyWidth = ab91Pointer(0)
}

func TestAb91CopyIsolation(ab91T *testing.T) {
	ab91Before := ab91Catalog(ab91T)
	ab91Catalog := ab91Catalog(ab91T)
	ab91Key := ab91Before.Ab91Entries[0].Ab91Key
	ab91MutateNested(&ab91Catalog.Ab91Entries[0])
	ab91Catalog.Ab91Limitations[0] = "ab91_mutated"
	ab91After, ab91Error := Ab91B51CDPartialLabels()
	if ab91Error != nil || !reflect.DeepEqual(ab91After, ab91Before) {
		ab91T.Fatal("ab91 accessor shared backing")
	}
	ab91Row, ab91Found, ab91Error := Ab91LookupB51CDPartialLabel(ab91Key)
	if ab91Error != nil || !ab91Found {
		ab91T.Fatal("ab91 lookup missing")
	}
	ab91MutateNested(&ab91Row)
	ab91Again, ab91Found, ab91Error := Ab91LookupB51CDPartialLabel(ab91Key)
	if ab91Error != nil || !ab91Found || !reflect.DeepEqual(ab91Again, ab91Before.Ab91Entries[0]) {
		ab91T.Fatal("ab91 lookup shared backing")
	}
	ab91After, ab91Error = Ab91B51CDPartialLabels()
	if ab91Error != nil || !reflect.DeepEqual(ab91After, ab91Before) {
		ab91T.Fatal("ab91 lookup mutated accessor")
	}
}

func TestAb91MetadataOnlySurface(ab91T *testing.T) {
	ab91File, ab91Error := parser.ParseFile(token.NewFileSet(), "ab91_b51cd_partial_labels.go", nil, 0)
	if ab91Error != nil {
		ab91T.Fatal(ab91Error)
	}
	ab91Functions := map[string]bool{"Ab91B51CDPartialLabels": false, "Ab91LookupB51CDPartialLabel": false, "ab91FindPartialLabel": false}
	for _, ab91Declaration := range ab91File.Decls {
		if ab91Function, ab91OK := ab91Declaration.(*ast.FuncDecl); ab91OK {
			if _, ab91Allowed := ab91Functions[ab91Function.Name.Name]; !ab91Allowed || ab91Function.Recv != nil {
				ab91T.Fatal("ab91 unexpected callable or emitter")
			}
			ab91Functions[ab91Function.Name.Name] = true
		}
	}
	for ab91Name, ab91Seen := range ab91Functions {
		if !ab91Seen {
			ab91T.Fatal(ab91Name)
		}
	}
	ab91AllowedImports := map[string]bool{"\"embed\"": true, "\"encoding/json\"": true, "\"reflect\"": true}
	for _, ab91Import := range ab91File.Imports {
		if !ab91AllowedImports[ab91Import.Path.Value] {
			ab91T.Fatal("ab91 unexpected invocation dependency")
		}
	}
	for _, ab91Declaration := range ab91File.Decls {
		ab91General, ab91OK := ab91Declaration.(*ast.GenDecl)
		if !ab91OK || ab91General.Tok != token.VAR {
			continue
		}
		for _, ab91Spec := range ab91General.Specs {
			ab91Value := ab91Spec.(*ast.ValueSpec)
			if len(ab91Value.Names) != 1 || ab91Value.Names[0].Name != "ab91PartialLabelsJSON" || len(ab91Value.Values) != 0 {
				ab91T.Fatal("ab91 unexpected registration or runtime cache")
			}
		}
	}
}

func TestAb91ExistingStubRemainsNonExecuting(ab91T *testing.T) {
	ab91Provider, ab91Error := ab91Standard.NewProvider(ab91Standard.Catalog{
		Namespace: "ebus_standard", Version: "ab91_synthetic",
		Services: []ab91Standard.Service{{PB: ab91Pointer(uint8(0xB5)), Commands: []ab91Standard.Command{{ID: "ab91_stub", SafetyClass: ab91Standard.SafetyReadOnlySafe}}}},
	}, true)
	if ab91Error != nil {
		ab91T.Fatal(ab91Error)
	}
	for ab91Phase := 0; ab91Phase < 2; ab91Phase++ {
		if ab91Phase == 1 {
			ab91Catalog := ab91Catalog(ab91T)
			for _, ab91Row := range ab91Catalog.Ab91Entries {
				if _, ab91Found, ab91Error := Ab91LookupB51CDPartialLabel(ab91Row.Ab91Key); ab91Error != nil || !ab91Found {
					ab91T.Fatal("ab91 lookup failed")
				}
				if _, ab91Error := ab91Provider.Invoke(ab91Context.Background(), ab91Row.Ab91Key.Ab91Label, nil, ab91Standard.CallerContextUserFacing); !ab91Errors.Is(ab91Error, ab91Standard.ErrUnknownMethod) {
					ab91T.Fatal("ab91 metadata label became callable")
				}
			}
		}
		ab91Result, ab91Error := ab91Provider.Invoke(ab91Context.Background(), "ab91_stub", nil, ab91Standard.CallerContextUserFacing)
		if ab91Error != nil || ab91Result["executed"] != false {
			ab91T.Fatal("ab91 existing synthetic stub became executable")
		}
	}
}
