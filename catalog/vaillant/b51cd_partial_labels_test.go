package vaillant

import (
	b51CDPartialLabelContext "context"
	b51CDPartialLabelErrors "errors"
	b51CDPartialLabelStandard "github.com/Project-Helianthus/helianthus-ebusreg/catalog/ebusstandard"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

func b51CDPartialLabelPointer[b51CDPartialLabelType any](b51CDPartialLabelValue b51CDPartialLabelType) *b51CDPartialLabelType {
	return &b51CDPartialLabelValue
}

func b51CDPartialLabelCatalog(b51CDPartialLabelT *testing.T) B51CDPartialLabelPartialLabels {
	b51CDPartialLabelT.Helper()
	b51CDPartialLabelCatalog, b51CDPartialLabelError := B51CDPartialLabels()
	if b51CDPartialLabelError != nil {
		b51CDPartialLabelT.Fatal(b51CDPartialLabelError)
	}
	return b51CDPartialLabelCatalog
}

func TestAb91CanonicalResponseDirection(b51CDPartialLabelT *testing.T) {
	b51CDPartialLabelReplies := 0
	for _, b51CDPartialLabelEntry := range b51CDPartialLabelCatalog(b51CDPartialLabelT).B51CDPartialLabelEntries {
		if b51CDPartialLabelEntry.B51CDPartialLabelKey.B51CDPartialLabelDirection != string(b51CDPartialLabelStandard.DirectionRequest) {
			continue
		}
		b51CDPartialLabelKey := b51CDPartialLabelEntry.B51CDPartialLabelKey
		b51CDPartialLabelKey.B51CDPartialLabelDirection = string(b51CDPartialLabelStandard.DirectionResponse)
		b51CDPartialLabelReply, b51CDPartialLabelFound, b51CDPartialLabelError := B51CDLookupPartialLabel(b51CDPartialLabelKey)
		if b51CDPartialLabelError != nil || !b51CDPartialLabelFound {
			b51CDPartialLabelT.Fatalf("ab91 canonical response missing for %s: %v", b51CDPartialLabelKey.B51CDPartialLabelLabel, b51CDPartialLabelError)
		}
		if !reflect.DeepEqual(b51CDPartialLabelReply.B51CDPartialLabelKey, b51CDPartialLabelKey) || b51CDPartialLabelReply.B51CDPartialLabelField.B51CDPartialLabelConstant == nil || *b51CDPartialLabelReply.B51CDPartialLabelField.B51CDPartialLabelConstant != 0xF5 {
			b51CDPartialLabelT.Fatalf("ab91 canonical response lost positive constant for %s", b51CDPartialLabelKey.B51CDPartialLabelLabel)
		}
		b51CDPartialLabelKey.B51CDPartialLabelDirection = "positive"
		if _, b51CDPartialLabelFound, b51CDPartialLabelError := B51CDLookupPartialLabel(b51CDPartialLabelKey); b51CDPartialLabelError != nil || b51CDPartialLabelFound {
			b51CDPartialLabelT.Fatal("ab91 positive classification matched a direction")
		}
		b51CDPartialLabelReplies++
	}
	if b51CDPartialLabelReplies != 23 {
		b51CDPartialLabelT.Fatalf("ab91 got %d canonical replies, want 23", b51CDPartialLabelReplies)
	}
}

func TestAb91ExactPartialLabels(b51CDPartialLabelT *testing.T) {
	b51CDPartialLabelCatalog := b51CDPartialLabelCatalog(b51CDPartialLabelT)
	b51CDPartialLabelWantCatalog := B51CDPartialLabelPartialLabels{
		B51CDPartialLabelRepresentation: "Proven", B51CDPartialLabelPhysicalQualification: "Unknown", B51CDPartialLabelInstalledSupport: "Unknown",
		B51CDPartialLabelEquivalence: "Unknown", B51CDPartialLabelBodyLayout: "Unknown", B51CDPartialLabelCodec: "Unknown", B51CDPartialLabelGating: "Unknown",
		B51CDPartialLabelConfirmation: "Unknown", B51CDPartialLabelEffect: "Unknown", B51CDPartialLabelPhysicalFraming: "Unknown", B51CDPartialLabelControls: "opaque", B51CDPartialLabelCompleteness: "incomplete",
		B51CDPartialLabelLimitations: []string{
			"Lexical suffixes identify labels only; they encode no secondary, operation or body byte.",
			"Logical byte and bit positions are not physical framing offsets.",
			"Partial forms are not empty encodable packets and establish no enum, constructor or decoder.",
			"No invocation, emitter, registration, runtime cache or transport is provided.",
			"Synthetic metadata representation establishes no physical support, freshness, persistence or execution permission.",
		},
	}
	b51CDPartialLabelLabels := map[string][]string{
		"B51C": {"00", "01", "02", "03", "10", "13", "14", "19", "1B", "1C", "1D"},
		"B51D": {"00", "01", "02", "03", "10", "11", "12", "13", "14", "15", "16", "17"},
	}
	if len(b51CDPartialLabelCatalog.B51CDPartialLabelEntries) != 46 {
		b51CDPartialLabelT.Fatalf("ab91 got %d rows, want 46", len(b51CDPartialLabelCatalog.B51CDPartialLabelEntries))
	}
	b51CDPartialLabelSeen := make(map[string]bool)
	for b51CDPartialLabelVariant, b51CDPartialLabelSuffixes := range b51CDPartialLabelLabels {
		for _, b51CDPartialLabelSuffix := range b51CDPartialLabelSuffixes {
			for b51CDPartialLabelDirection, b51CDPartialLabelConstant := range map[string]uint8{"request": 0xB5, "response": 0xF5} {
				b51CDPartialLabelKey := B51CDPartialKey{B51CDPartialLabelFamily: "vaillant", B51CDPartialLabelVariant: b51CDPartialLabelVariant, B51CDPartialLabelLabel: b51CDPartialLabelVariant + "_" + b51CDPartialLabelSuffix, B51CDPartialLabelPrimary: 0xB5, B51CDPartialLabelDirection: b51CDPartialLabelDirection, B51CDPartialLabelField: "coded_constant"}
				b51CDPartialLabelT.Run("ab91_"+b51CDPartialLabelKey.B51CDPartialLabelLabel+"_"+b51CDPartialLabelDirection, func(b51CDPartialLabelT *testing.T) {
					b51CDPartialLabelRow, b51CDPartialLabelFound, b51CDPartialLabelError := B51CDLookupPartialLabel(b51CDPartialLabelKey)
					if b51CDPartialLabelError != nil || !b51CDPartialLabelFound {
						b51CDPartialLabelT.Fatalf("ab91 lookup: %v %v", b51CDPartialLabelFound, b51CDPartialLabelError)
					}
					b51CDPartialLabelWant := B51CDPartialLabel{
						B51CDPartialLabelKey: b51CDPartialLabelKey, B51CDPartialLabelGrade: "Hypothesis", B51CDPartialLabelRawFallback: "preserve full original carrier",
						B51CDPartialLabelField: B51CDPartialField{B51CDPartialLabelLogicalByte: b51CDPartialLabelPointer(0), B51CDPartialLabelLogicalBit: b51CDPartialLabelPointer(0), B51CDPartialLabelWidthBits: b51CDPartialLabelPointer(8), B51CDPartialLabelConstant: b51CDPartialLabelPointer(b51CDPartialLabelConstant), B51CDPartialLabelSignedness: "Unknown", B51CDPartialLabelByteOrder: "Unknown", B51CDPartialLabelConversion: "Unknown", B51CDPartialLabelGrade: "Hypothesis", B51CDPartialLabelRawFallback: "preserve original bytes without interpreting the body"},
					}
					if !reflect.DeepEqual(b51CDPartialLabelRow, b51CDPartialLabelWant) {
						b51CDPartialLabelT.Fatalf("ab91 metadata mismatch: %#v", b51CDPartialLabelRow)
					}
					b51CDPartialLabelMatches := 0
					for _, b51CDPartialLabelEntry := range b51CDPartialLabelCatalog.B51CDPartialLabelEntries {
						if reflect.DeepEqual(b51CDPartialLabelEntry, b51CDPartialLabelWant) {
							b51CDPartialLabelMatches++
						}
					}
					if b51CDPartialLabelMatches != 1 {
						b51CDPartialLabelT.Fatalf("ab91 want one independent row, got %d", b51CDPartialLabelMatches)
					}
					for _, b51CDPartialLabelCase := range []struct {
						b51CDPartialLabelRaw  []byte
						b51CDPartialLabelWant bool
					}{
						{[]byte{b51CDPartialLabelConstant, 0x73, 0xA2}, true}, {[]byte{b51CDPartialLabelConstant ^ 1, 0x73}, false}, {nil, false},
					} {
						b51CDPartialLabelOriginal := append([]byte(nil), b51CDPartialLabelCase.b51CDPartialLabelRaw...)
						if b51CDPartialLabelConstantMatches(b51CDPartialLabelRow.B51CDPartialLabelField, b51CDPartialLabelCase.b51CDPartialLabelRaw) != b51CDPartialLabelCase.b51CDPartialLabelWant {
							b51CDPartialLabelT.Fatal("ab91 synthetic partial constant check")
						}
						if !reflect.DeepEqual(b51CDPartialLabelOriginal, b51CDPartialLabelCase.b51CDPartialLabelRaw) {
							b51CDPartialLabelT.Fatal("ab91 opaque raw body changed")
						}
					}
				})
			}
			b51CDPartialLabelSeen[b51CDPartialLabelVariant+"_"+b51CDPartialLabelSuffix] = true
		}
	}
	if len(b51CDPartialLabelSeen) != 23 {
		b51CDPartialLabelT.Fatal("ab91 label identities merged")
	}
	b51CDPartialLabelCatalog.B51CDPartialLabelEntries = nil
	if !reflect.DeepEqual(b51CDPartialLabelCatalog, b51CDPartialLabelWantCatalog) {
		b51CDPartialLabelT.Fatal("ab91 qualification or unknown-body contract changed")
	}
}

func b51CDPartialLabelConstantMatches(b51CDPartialLabelField B51CDPartialField, b51CDPartialLabelRaw []byte) bool {
	if b51CDPartialLabelField.B51CDPartialLabelLogicalByte == nil || b51CDPartialLabelField.B51CDPartialLabelLogicalBit == nil || b51CDPartialLabelField.B51CDPartialLabelWidthBits == nil || b51CDPartialLabelField.B51CDPartialLabelConstant == nil {
		return false
	}
	b51CDPartialLabelPosition := *b51CDPartialLabelField.B51CDPartialLabelLogicalByte
	if b51CDPartialLabelPosition < 0 || b51CDPartialLabelPosition >= len(b51CDPartialLabelRaw) {
		return false
	}
	b51CDPartialLabelMask := uint16(1)<<uint(*b51CDPartialLabelField.B51CDPartialLabelWidthBits) - 1
	return uint16(b51CDPartialLabelRaw[b51CDPartialLabelPosition])>>uint(*b51CDPartialLabelField.B51CDPartialLabelLogicalBit)&b51CDPartialLabelMask == uint16(*b51CDPartialLabelField.B51CDPartialLabelConstant)
}

func TestAb91ExactIdentityAxes(b51CDPartialLabelT *testing.T) {
	b51CDPartialLabelBase := b51CDPartialLabelCatalog(b51CDPartialLabelT).B51CDPartialLabelEntries[0]
	b51CDPartialLabelQualified := b51CDPartialLabelBase
	b51CDPartialLabelQualified.B51CDPartialLabelKey = B51CDPartialKey{B51CDPartialLabelFamily: "vaillant", B51CDPartialLabelVariant: "B524", B51CDPartialLabelLabel: "ab91_synthetic", B51CDPartialLabelPrimary: 0xB5, B51CDPartialLabelSecondary: b51CDPartialLabelPointer(uint16(0x24)), B51CDPartialLabelOperation: b51CDPartialLabelPointer(uint16(1)), B51CDPartialLabelDirection: "request", B51CDPartialLabelField: "ab91_synthetic", B51CDPartialLabelOpcode: b51CDPartialLabelPointer(uint16(2)), B51CDPartialLabelGroup: b51CDPartialLabelPointer(uint16(3)), B51CDPartialLabelInstance: b51CDPartialLabelPointer(uint16(4)), B51CDPartialLabelRegister: b51CDPartialLabelPointer(uint16(5))}
	for _, b51CDPartialLabelControl := range []B51CDPartialLabel{b51CDPartialLabelBase, b51CDPartialLabelQualified} {
		for _, b51CDPartialLabelAxis := range []string{"B51CDPartialLabelFamily", "B51CDPartialLabelVariant", "B51CDPartialLabelLabel", "B51CDPartialLabelPrimary", "B51CDPartialLabelSecondary", "B51CDPartialLabelOperation", "B51CDPartialLabelDirection", "B51CDPartialLabelField", "B51CDPartialLabelOpcode", "B51CDPartialLabelGroup", "B51CDPartialLabelInstance", "B51CDPartialLabelRegister"} {
			b51CDPartialLabelT.Run("ab91_"+b51CDPartialLabelControl.B51CDPartialLabelKey.B51CDPartialLabelVariant+"_"+b51CDPartialLabelAxis, func(b51CDPartialLabelT *testing.T) {
				b51CDPartialLabelKey := b51CDPartialLabelControl.B51CDPartialLabelKey
				b51CDPartialLabelMember := reflect.ValueOf(&b51CDPartialLabelKey).Elem().FieldByName(b51CDPartialLabelAxis)
				switch b51CDPartialLabelMember.Kind() {
				case reflect.String:
					b51CDPartialLabelMember.SetString("ab91_distinct")
				case reflect.Uint8:
					b51CDPartialLabelMember.SetUint(0)
				case reflect.Pointer:
					b51CDPartialLabelMember.Set(reflect.ValueOf(b51CDPartialLabelPointer(uint16(0))))
				default:
					b51CDPartialLabelT.Fatal("ab91 untested axis")
				}
				if _, b51CDPartialLabelFound := b51CDPartialLabelFindPartialLabel([]B51CDPartialLabel{b51CDPartialLabelControl}, b51CDPartialLabelKey); b51CDPartialLabelFound {
					b51CDPartialLabelT.Fatal("ab91 distinct axis wildcard matched")
				}
				if b51CDPartialLabelControl.B51CDPartialLabelKey.B51CDPartialLabelVariant != "B524" {
					if _, b51CDPartialLabelFound, b51CDPartialLabelError := B51CDLookupPartialLabel(b51CDPartialLabelKey); b51CDPartialLabelError != nil || b51CDPartialLabelFound {
						b51CDPartialLabelT.Fatal("ab91 production lookup ignored axis")
					}
				}
				b51CDPartialLabelDistinct := b51CDPartialLabelControl
				b51CDPartialLabelDistinct.B51CDPartialLabelKey = b51CDPartialLabelKey
				b51CDPartialLabelDistinct.B51CDPartialLabelRawFallback = "ab91_distinct"
				b51CDPartialLabelRows := []B51CDPartialLabel{b51CDPartialLabelControl, b51CDPartialLabelDistinct}
				if b51CDPartialLabelGot, b51CDPartialLabelFound := b51CDPartialLabelFindPartialLabel(b51CDPartialLabelRows, b51CDPartialLabelKey); !b51CDPartialLabelFound || !reflect.DeepEqual(b51CDPartialLabelGot, b51CDPartialLabelDistinct) {
					b51CDPartialLabelT.Fatal("ab91 qualified positive control failed")
				}
				if b51CDPartialLabelGot, b51CDPartialLabelFound := b51CDPartialLabelFindPartialLabel(b51CDPartialLabelRows, b51CDPartialLabelControl.B51CDPartialLabelKey); !b51CDPartialLabelFound || !reflect.DeepEqual(b51CDPartialLabelGot, b51CDPartialLabelControl) {
					b51CDPartialLabelT.Fatal("ab91 original positive control failed")
				}
			})
		}
	}
	for _, b51CDPartialLabelAxis := range []string{"B51CDPartialLabelSecondary", "B51CDPartialLabelOperation", "B51CDPartialLabelOpcode", "B51CDPartialLabelGroup", "B51CDPartialLabelInstance", "B51CDPartialLabelRegister"} {
		b51CDPartialLabelT.Run("ab91_absent_zero_nonzero_"+b51CDPartialLabelAxis, func(b51CDPartialLabelT *testing.T) {
			for _, b51CDPartialLabelCoordinate := range []*uint16{nil, b51CDPartialLabelPointer(uint16(0)), b51CDPartialLabelPointer(uint16(1))} {
				b51CDPartialLabelRow := b51CDPartialLabelQualified
				reflect.ValueOf(&b51CDPartialLabelRow.B51CDPartialLabelKey).Elem().FieldByName(b51CDPartialLabelAxis).Set(reflect.ValueOf(b51CDPartialLabelCoordinate))
				for _, b51CDPartialLabelQuery := range []*uint16{nil, b51CDPartialLabelPointer(uint16(0)), b51CDPartialLabelPointer(uint16(1))} {
					b51CDPartialLabelKey := b51CDPartialLabelRow.B51CDPartialLabelKey
					reflect.ValueOf(&b51CDPartialLabelKey).Elem().FieldByName(b51CDPartialLabelAxis).Set(reflect.ValueOf(b51CDPartialLabelQuery))
					_, b51CDPartialLabelFound := b51CDPartialLabelFindPartialLabel([]B51CDPartialLabel{b51CDPartialLabelRow}, b51CDPartialLabelKey)
					if b51CDPartialLabelFound != reflect.DeepEqual(b51CDPartialLabelCoordinate, b51CDPartialLabelQuery) {
						b51CDPartialLabelT.Fatal("ab91 absent and numeric coordinates conflated")
					}
				}
			}
		})
	}
}

func b51CDPartialLabelMutateNested(b51CDPartialLabelRow *B51CDPartialLabel) {
	*b51CDPartialLabelRow.B51CDPartialLabelField.B51CDPartialLabelLogicalByte = 7
	*b51CDPartialLabelRow.B51CDPartialLabelField.B51CDPartialLabelLogicalBit = 6
	*b51CDPartialLabelRow.B51CDPartialLabelField.B51CDPartialLabelWidthBits = 3
	*b51CDPartialLabelRow.B51CDPartialLabelField.B51CDPartialLabelConstant = 0
	b51CDPartialLabelRow.B51CDPartialLabelKey.B51CDPartialLabelLabel = "ab91_mutated"
	b51CDPartialLabelRow.B51CDPartialLabelKey.B51CDPartialLabelSecondary = b51CDPartialLabelPointer(uint16(0))
	b51CDPartialLabelRow.B51CDPartialLabelKey.B51CDPartialLabelOperation = b51CDPartialLabelPointer(uint16(0))
	b51CDPartialLabelRow.B51CDPartialLabelField.B51CDPartialLabelEnumeration = map[string]string{"ab91_mutated": "ab91_mutated"}
	b51CDPartialLabelRow.B51CDPartialLabelField.B51CDPartialLabelReserved = []int64{1}
	b51CDPartialLabelRow.B51CDPartialLabelField.B51CDPartialLabelReplacement = []int64{2}
	b51CDPartialLabelRow.B51CDPartialLabelField.B51CDPartialLabelBounds = []int64{0, 1}
	b51CDPartialLabelRow.B51CDPartialLabelBodyWidth = b51CDPartialLabelPointer(0)
}

func TestAb91CopyIsolation(b51CDPartialLabelT *testing.T) {
	b51CDPartialLabelBefore := b51CDPartialLabelCatalog(b51CDPartialLabelT)
	b51CDPartialLabelCatalog := b51CDPartialLabelCatalog(b51CDPartialLabelT)
	b51CDPartialLabelKey := b51CDPartialLabelBefore.B51CDPartialLabelEntries[0].B51CDPartialLabelKey
	b51CDPartialLabelMutateNested(&b51CDPartialLabelCatalog.B51CDPartialLabelEntries[0])
	b51CDPartialLabelCatalog.B51CDPartialLabelLimitations[0] = "ab91_mutated"
	b51CDPartialLabelAfter, b51CDPartialLabelError := B51CDPartialLabels()
	if b51CDPartialLabelError != nil || !reflect.DeepEqual(b51CDPartialLabelAfter, b51CDPartialLabelBefore) {
		b51CDPartialLabelT.Fatal("ab91 accessor shared backing")
	}
	b51CDPartialLabelRow, b51CDPartialLabelFound, b51CDPartialLabelError := B51CDLookupPartialLabel(b51CDPartialLabelKey)
	if b51CDPartialLabelError != nil || !b51CDPartialLabelFound {
		b51CDPartialLabelT.Fatal("ab91 lookup missing")
	}
	b51CDPartialLabelMutateNested(&b51CDPartialLabelRow)
	b51CDPartialLabelAgain, b51CDPartialLabelFound, b51CDPartialLabelError := B51CDLookupPartialLabel(b51CDPartialLabelKey)
	if b51CDPartialLabelError != nil || !b51CDPartialLabelFound || !reflect.DeepEqual(b51CDPartialLabelAgain, b51CDPartialLabelBefore.B51CDPartialLabelEntries[0]) {
		b51CDPartialLabelT.Fatal("ab91 lookup shared backing")
	}
	b51CDPartialLabelAfter, b51CDPartialLabelError = B51CDPartialLabels()
	if b51CDPartialLabelError != nil || !reflect.DeepEqual(b51CDPartialLabelAfter, b51CDPartialLabelBefore) {
		b51CDPartialLabelT.Fatal("ab91 lookup mutated accessor")
	}
}

func TestAb91MetadataOnlySurface(b51CDPartialLabelT *testing.T) {
	b51CDPartialLabelFile, b51CDPartialLabelError := parser.ParseFile(token.NewFileSet(), "b51cd_partial_labels.go", nil, 0)
	if b51CDPartialLabelError != nil {
		b51CDPartialLabelT.Fatal(b51CDPartialLabelError)
	}
	b51CDPartialLabelFunctions := map[string]bool{"B51CDPartialLabels": false, "B51CDLookupPartialLabel": false, "b51CDPartialLabelFindPartialLabel": false}
	for _, b51CDPartialLabelDeclaration := range b51CDPartialLabelFile.Decls {
		if b51CDPartialLabelFunction, b51CDPartialLabelOK := b51CDPartialLabelDeclaration.(*ast.FuncDecl); b51CDPartialLabelOK {
			if _, b51CDPartialLabelAllowed := b51CDPartialLabelFunctions[b51CDPartialLabelFunction.Name.Name]; !b51CDPartialLabelAllowed || b51CDPartialLabelFunction.Recv != nil {
				b51CDPartialLabelT.Fatal("ab91 unexpected callable or emitter")
			}
			b51CDPartialLabelFunctions[b51CDPartialLabelFunction.Name.Name] = true
		}
	}
	for b51CDPartialLabelName, b51CDPartialLabelSeen := range b51CDPartialLabelFunctions {
		if !b51CDPartialLabelSeen {
			b51CDPartialLabelT.Fatal(b51CDPartialLabelName)
		}
	}
	b51CDPartialLabelAllowedImports := map[string]bool{"\"embed\"": true, "\"encoding/json\"": true, "\"reflect\"": true}
	for _, b51CDPartialLabelImport := range b51CDPartialLabelFile.Imports {
		if !b51CDPartialLabelAllowedImports[b51CDPartialLabelImport.Path.Value] {
			b51CDPartialLabelT.Fatal("ab91 unexpected invocation dependency")
		}
	}
	for _, b51CDPartialLabelDeclaration := range b51CDPartialLabelFile.Decls {
		b51CDPartialLabelGeneral, b51CDPartialLabelOK := b51CDPartialLabelDeclaration.(*ast.GenDecl)
		if !b51CDPartialLabelOK || b51CDPartialLabelGeneral.Tok != token.VAR {
			continue
		}
		for _, b51CDPartialLabelSpec := range b51CDPartialLabelGeneral.Specs {
			b51CDPartialLabelValue := b51CDPartialLabelSpec.(*ast.ValueSpec)
			if len(b51CDPartialLabelValue.Names) != 1 || b51CDPartialLabelValue.Names[0].Name != "b51CDPartialLabelPartialLabelsJSON" || len(b51CDPartialLabelValue.Values) != 0 {
				b51CDPartialLabelT.Fatal("ab91 unexpected registration or runtime cache")
			}
		}
	}
}

func TestAb91ExistingStubRemainsNonExecuting(b51CDPartialLabelT *testing.T) {
	b51CDPartialLabelProvider, b51CDPartialLabelError := b51CDPartialLabelStandard.NewProvider(b51CDPartialLabelStandard.Catalog{
		Namespace: "ebus_standard", Version: "ab91_synthetic",
		Services: []b51CDPartialLabelStandard.Service{{PB: b51CDPartialLabelPointer(uint8(0xB5)), Commands: []b51CDPartialLabelStandard.Command{{ID: "ab91_stub", SafetyClass: b51CDPartialLabelStandard.SafetyReadOnlySafe}}}},
	}, true)
	if b51CDPartialLabelError != nil {
		b51CDPartialLabelT.Fatal(b51CDPartialLabelError)
	}
	for b51CDPartialLabelPhase := 0; b51CDPartialLabelPhase < 2; b51CDPartialLabelPhase++ {
		if b51CDPartialLabelPhase == 1 {
			b51CDPartialLabelCatalog := b51CDPartialLabelCatalog(b51CDPartialLabelT)
			for _, b51CDPartialLabelRow := range b51CDPartialLabelCatalog.B51CDPartialLabelEntries {
				if _, b51CDPartialLabelFound, b51CDPartialLabelError := B51CDLookupPartialLabel(b51CDPartialLabelRow.B51CDPartialLabelKey); b51CDPartialLabelError != nil || !b51CDPartialLabelFound {
					b51CDPartialLabelT.Fatal("ab91 lookup failed")
				}
				if _, b51CDPartialLabelError := b51CDPartialLabelProvider.Invoke(b51CDPartialLabelContext.Background(), b51CDPartialLabelRow.B51CDPartialLabelKey.B51CDPartialLabelLabel, nil, b51CDPartialLabelStandard.CallerContextUserFacing); !b51CDPartialLabelErrors.Is(b51CDPartialLabelError, b51CDPartialLabelStandard.ErrUnknownMethod) {
					b51CDPartialLabelT.Fatal("ab91 metadata label became callable")
				}
			}
		}
		b51CDPartialLabelResult, b51CDPartialLabelError := b51CDPartialLabelProvider.Invoke(b51CDPartialLabelContext.Background(), "ab91_stub", nil, b51CDPartialLabelStandard.CallerContextUserFacing)
		if b51CDPartialLabelError != nil || b51CDPartialLabelResult["executed"] != false {
			b51CDPartialLabelT.Fatal("ab91 existing synthetic stub became executable")
		}
	}
}
