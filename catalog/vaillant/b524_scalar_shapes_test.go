package vaillant

import (
	b524ScalarShapeBytes "bytes"
	b524ScalarShapeContext "context"
	b524ScalarShapeBinary "encoding/binary"
	b524ScalarShapeJSON "encoding/json"
	b524ScalarShapeErrors "errors"
	b524ScalarShapeMath "math"
	b524ScalarShapeReflect "reflect"
	b524ScalarShapeTesting "testing"

	b524ScalarShapeStandard "github.com/Project-Helianthus/helianthus-ebusreg/catalog/ebusstandard"
)

func b524ScalarShapePointer[b524ScalarShapeValue any](b524ScalarShapeInput b524ScalarShapeValue) *b524ScalarShapeValue {
	return &b524ScalarShapeInput
}

func b524ScalarShapeExpectedKey(b524ScalarShapeFamily *string, b524ScalarShapeDirection string, b524ScalarShapeGroup uint8, b524ScalarShapeRegister uint16, b524ScalarShapeField string) B524ScalarShapeKey {
	return B524ScalarShapeKey{B524ScalarShapeFamily: b524ScalarShapeFamily, B524ScalarShapePrimary: b524ScalarShapePointer(uint8(0xB5)), B524ScalarShapeSecondary: b524ScalarShapePointer(uint8(0x24)), B524ScalarShapeDirection: b524ScalarShapePointer(b524ScalarShapeDirection), B524ScalarShapeField: b524ScalarShapePointer(b524ScalarShapeField), B524ScalarShapeOpcode: b524ScalarShapePointer(uint8(2)), B524ScalarShapeGroup: b524ScalarShapePointer(b524ScalarShapeGroup), B524ScalarShapeRegister: b524ScalarShapePointer(b524ScalarShapeRegister)}
}

func b524ScalarShapeExpectedEntries() []B524ScalarShape {
	var b524ScalarShapeEntries []B524ScalarShape
	for _, b524ScalarShapeCase := range []struct {
		b524ScalarShapeFamily                      *string
		b524ScalarShapeDirection                   string
		b524ScalarShapeGroup                       uint8
		b524ScalarShapeRegister                    uint16
		b524ScalarShapeField, b524ScalarShapeCodec string
		b524ScalarShapeWidth                       int
	}{
		{nil, "read", 3, 0x001B, "scalar", "u8", 1},
		{nil, "read", 3, 0x001B, "scalar", "u16", 2},
		{nil, "read", 1, 0x0014, "scalar", "u8", 1},
		{nil, "read", 1, 0x0019, "scalar", "u8", 1},
		{nil, "read", 1, 0x0015, "scalar", "binary32", 4},
		{b524ScalarShapePointer("controller-family"), "write", 3, 0x0006, "mode", "u16", 2},
		{b524ScalarShapePointer("controller-family"), "write", 3, 0x0022, "manual-temperature", "binary32", 4},
		{b524ScalarShapePointer("controller-family"), "write", 3, 0x0007, "comfort-temperature", "binary32", 4},
	} {
		b524ScalarShapeSignedness, b524ScalarShapeOrder := "unsigned", "little-endian"
		if b524ScalarShapeCase.b524ScalarShapeCodec == "binary32" {
			b524ScalarShapeSignedness = "signed"
		}
		if b524ScalarShapeCase.b524ScalarShapeCodec == "u8" {
			b524ScalarShapeOrder = "not applicable"
		}
		b524ScalarShapeEntries = append(b524ScalarShapeEntries, B524ScalarShape{
			B524ScalarShapeKey:   b524ScalarShapeExpectedKey(b524ScalarShapeCase.b524ScalarShapeFamily, b524ScalarShapeCase.b524ScalarShapeDirection, b524ScalarShapeCase.b524ScalarShapeGroup, b524ScalarShapeCase.b524ScalarShapeRegister, b524ScalarShapeCase.b524ScalarShapeField),
			B524ScalarShapeShape: b524ScalarShapeCase.b524ScalarShapeCodec, B524ScalarShapeWidthBytes: b524ScalarShapePointer(b524ScalarShapeCase.b524ScalarShapeWidth),
			B524ScalarShapeSignedness: b524ScalarShapeSignedness, B524ScalarShapeByteOrder: b524ScalarShapeOrder,
			B524ScalarShapeConversion: &B524ScalarShapeConversion{B524ScalarShapeEncoding: b524ScalarShapeCase.b524ScalarShapeCodec, B524ScalarShapeFactor: 1, B524ScalarShapeOffset: 0},
			B524ScalarShapeGrade:      "Hypothesis", B524ScalarShapeRawFallback: "preserve native field bytes on every unsuccessful typed interpretation",
		})
	}
	return b524ScalarShapeEntries
}

func b524ScalarShapeCatalog(b524ScalarShapeT *b524ScalarShapeTesting.T) B524ScalarShapeCatalog {
	b524ScalarShapeT.Helper()
	b524ScalarShapeCatalog, b524ScalarShapeError := B524ScalarShapes()
	if b524ScalarShapeError != nil {
		b524ScalarShapeT.Fatal(b524ScalarShapeError)
	}
	return b524ScalarShapeCatalog
}

func TestAb85ProductionMetadata(b524ScalarShapeT *b524ScalarShapeTesting.T) {
	var b524ScalarShapeNative struct {
		B524ScalarShapeEntries []map[string]b524ScalarShapeJSON.RawMessage `json:"entries"`
	}
	if b524ScalarShapeError := b524ScalarShapeJSON.Unmarshal(b524ScalarShapeScalarShapesJSON, &b524ScalarShapeNative); b524ScalarShapeError != nil {
		b524ScalarShapeT.Fatal(b524ScalarShapeError)
	}
	for _, b524ScalarShapeEntry := range b524ScalarShapeNative.B524ScalarShapeEntries {
		for _, b524ScalarShapeLabel := range []string{"key", "shape", "logical_position", "width_bytes", "signedness", "byte_order", "conversion", "unit", "default", "bounds", "enumeration", "replacement_values", "reserved_values", "grade", "raw_fallback"} {
			if _, b524ScalarShapePresent := b524ScalarShapeEntry[b524ScalarShapeLabel]; !b524ScalarShapePresent {
				b524ScalarShapeT.Fatalf("ab85 missing explicit field %s", b524ScalarShapeLabel)
			}
		}
		var b524ScalarShapeKey map[string]b524ScalarShapeJSON.RawMessage
		if b524ScalarShapeError := b524ScalarShapeJSON.Unmarshal(b524ScalarShapeEntry["key"], &b524ScalarShapeKey); b524ScalarShapeError != nil {
			b524ScalarShapeT.Fatal(b524ScalarShapeError)
		}
		for _, b524ScalarShapeLabel := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "opcode", "group", "instance", "register"} {
			if _, b524ScalarShapePresent := b524ScalarShapeKey[b524ScalarShapeLabel]; !b524ScalarShapePresent {
				b524ScalarShapeT.Fatalf("ab85 missing explicit identity axis %s", b524ScalarShapeLabel)
			}
		}
	}
	b524ScalarShapeGot := b524ScalarShapeCatalog(b524ScalarShapeT)
	b524ScalarShapeWant := B524ScalarShapeCatalog{
		B524ScalarShapeCatalog: "ab85_b524_scalar_shapes", B524ScalarShapeRepresentationGrade: "Proven", B524ScalarShapePhysicalQualification: "Unknown", B524ScalarShapeComplete: false,
		B524ScalarShapeLimitations: []string{
			"Finite supplemental declarations; broader family coverage is incomplete.",
			"Family and instance absence is exact missing qualification, never a wildcard.",
			"Request length and reply length do not choose a semantic alternative or family qualification.",
			"Decode failure retains a prior cached value with uncertain freshness; no cache updater is supplied.",
			"Installed support, omitted instance predicates, complete enumerations, units, persistence and wire adaptation are Unknown.",
			"No builder, executor, provider registration, invocation, admission or runtime cache is supplied.",
			"Logical positions including prefixes are Unknown; widths describe isolated candidate field bytes only.",
		}, B524ScalarShapeEntries: b524ScalarShapeExpectedEntries(),
	}
	if !b524ScalarShapeReflect.DeepEqual(b524ScalarShapeGot, b524ScalarShapeWant) {
		b524ScalarShapeT.Fatalf("ab85 production metadata differs from independent contract: %+v", b524ScalarShapeGot)
	}
}

func TestAb85ExactIdentity(b524ScalarShapeT *b524ScalarShapeTesting.T) {
	for b524ScalarShapeIndex, b524ScalarShapeExpected := range b524ScalarShapeExpectedEntries() {
		b524ScalarShapeControl, b524ScalarShapeError := B524ScalarShapeLookup(b524ScalarShapeExpected.B524ScalarShapeKey)
		if b524ScalarShapeError != nil {
			b524ScalarShapeT.Fatal(b524ScalarShapeError)
		}
		b524ScalarShapeCount := 1
		if b524ScalarShapeIndex < 2 {
			b524ScalarShapeCount = 2
		}
		if len(b524ScalarShapeControl) != b524ScalarShapeCount {
			b524ScalarShapeT.Fatalf("ab85 valid control %d: %v", b524ScalarShapeIndex, b524ScalarShapeControl)
		}
		for b524ScalarShapeAxis := 0; b524ScalarShapeAxis < b524ScalarShapeReflect.TypeOf(B524ScalarShapeKey{}).NumField(); b524ScalarShapeAxis++ {
			for b524ScalarShapeMutation := 0; b524ScalarShapeMutation < 3; b524ScalarShapeMutation++ {
				b524ScalarShapeKey := b524ScalarShapeExpected.B524ScalarShapeKey
				b524ScalarShapeSlot := b524ScalarShapeReflect.ValueOf(&b524ScalarShapeKey).Elem().Field(b524ScalarShapeAxis)
				b524ScalarShapeOld := b524ScalarShapeReflect.ValueOf(b524ScalarShapeExpected.B524ScalarShapeKey).Field(b524ScalarShapeAxis)
				if b524ScalarShapeMutation == 0 {
					b524ScalarShapeSlot.SetZero()
				} else {
					b524ScalarShapeValue := b524ScalarShapeReflect.New(b524ScalarShapeSlot.Type().Elem())
					if b524ScalarShapeValue.Elem().Kind() == b524ScalarShapeReflect.String {
						if b524ScalarShapeMutation == 2 {
							b524ScalarShapeValue.Elem().SetString("ab85_other")
						}
					} else if b524ScalarShapeMutation == 2 {
						b524ScalarShapeNumber := uint64(6)
						if !b524ScalarShapeOld.IsNil() {
							b524ScalarShapeNumber = b524ScalarShapeOld.Elem().Uint() ^ 0x80
						}
						if b524ScalarShapeReflect.TypeOf(b524ScalarShapeKey).Field(b524ScalarShapeAxis).Name == "B524ScalarShapeOpcode" {
							b524ScalarShapeNumber = 6
						}
						if b524ScalarShapeReflect.TypeOf(b524ScalarShapeKey).Field(b524ScalarShapeAxis).Name == "B524ScalarShapeRegister" {
							b524ScalarShapeNumber = b524ScalarShapeOld.Elem().Uint() ^ 0x100
						}
						b524ScalarShapeValue.Elem().SetUint(b524ScalarShapeNumber)
					}
					b524ScalarShapeSlot.Set(b524ScalarShapeValue)
				}
				if b524ScalarShapeReflect.DeepEqual(b524ScalarShapeKey, b524ScalarShapeExpected.B524ScalarShapeKey) {
					continue
				}
				b524ScalarShapeMatches, b524ScalarShapeError := B524ScalarShapeLookup(b524ScalarShapeKey)
				if b524ScalarShapeError != nil || len(b524ScalarShapeMatches) != 0 {
					b524ScalarShapeT.Fatalf("ab85 identity axis %s mutation %d inherited entry: %v %v", b524ScalarShapeReflect.TypeOf(b524ScalarShapeKey).Field(b524ScalarShapeAxis).Name, b524ScalarShapeMutation, b524ScalarShapeMatches, b524ScalarShapeError)
				}
			}
		}
		b524ScalarShapeWrongDirection := b524ScalarShapeExpected.B524ScalarShapeKey
		if *b524ScalarShapeWrongDirection.B524ScalarShapeDirection == "read" {
			b524ScalarShapeWrongDirection.B524ScalarShapeDirection = b524ScalarShapePointer("write")
		} else {
			b524ScalarShapeWrongDirection.B524ScalarShapeDirection = b524ScalarShapePointer("read")
		}
		b524ScalarShapeMatches, b524ScalarShapeError := B524ScalarShapeLookup(b524ScalarShapeWrongDirection)
		if b524ScalarShapeError != nil || len(b524ScalarShapeMatches) != 0 {
			b524ScalarShapeT.Fatal("ab85 read/write leaked", b524ScalarShapeError)
		}
	}
}

func b524ScalarShapeEvaluate(b524ScalarShapeShape B524ScalarShape, b524ScalarShapeRaw []byte) (float64, bool, []byte) {
	b524ScalarShapeEvidence := b524ScalarShapeBytes.Clone(b524ScalarShapeRaw)
	if b524ScalarShapeShape.B524ScalarShapeWidthBytes == nil || len(b524ScalarShapeRaw) != *b524ScalarShapeShape.B524ScalarShapeWidthBytes || b524ScalarShapeShape.B524ScalarShapeConversion == nil {
		return 0, false, b524ScalarShapeEvidence
	}
	var b524ScalarShapeValue float64
	switch b524ScalarShapeShape.B524ScalarShapeConversion.B524ScalarShapeEncoding {
	case "u8":
		if len(b524ScalarShapeRaw) != 1 || b524ScalarShapeShape.B524ScalarShapeSignedness != "unsigned" || b524ScalarShapeShape.B524ScalarShapeByteOrder != "not applicable" {
			return 0, false, b524ScalarShapeEvidence
		}
		b524ScalarShapeValue = float64(b524ScalarShapeRaw[0])
	case "u16":
		if len(b524ScalarShapeRaw) != 2 || b524ScalarShapeShape.B524ScalarShapeSignedness != "unsigned" || b524ScalarShapeShape.B524ScalarShapeByteOrder != "little-endian" {
			return 0, false, b524ScalarShapeEvidence
		}
		b524ScalarShapeValue = float64(b524ScalarShapeBinary.LittleEndian.Uint16(b524ScalarShapeRaw))
	case "binary32":
		if len(b524ScalarShapeRaw) != 4 || b524ScalarShapeShape.B524ScalarShapeSignedness != "signed" || b524ScalarShapeShape.B524ScalarShapeByteOrder != "little-endian" {
			return 0, false, b524ScalarShapeEvidence
		}
		b524ScalarShapeValue = float64(b524ScalarShapeMath.Float32frombits(b524ScalarShapeBinary.LittleEndian.Uint32(b524ScalarShapeRaw)))
	default:
		return 0, false, b524ScalarShapeEvidence
	}
	return b524ScalarShapeValue*b524ScalarShapeShape.B524ScalarShapeConversion.B524ScalarShapeFactor + b524ScalarShapeShape.B524ScalarShapeConversion.B524ScalarShapeOffset, true, b524ScalarShapeEvidence
}

func TestAb85SyntheticFields(b524ScalarShapeT *b524ScalarShapeTesting.T) {
	b524ScalarShapeEntries := b524ScalarShapeCatalog(b524ScalarShapeT).B524ScalarShapeEntries
	for _, b524ScalarShapeCase := range []struct {
		b524ScalarShapeIndex int
		b524ScalarShapeRaw   []byte
		b524ScalarShapeWant  float64
		b524ScalarShapeOK    bool
	}{
		{0, []byte{1}, 1, true}, {1, []byte{1}, 0, false},
		{0, []byte{1, 0}, 0, false}, {1, []byte{1, 0}, 1, true},
		{0, []byte{0, 1}, 0, false}, {1, []byte{0, 1}, 256, true},
		{2, []byte{14}, 14, true}, {3, []byte{19}, 19, true},
		{4, []byte{0, 0, 128, 63}, 1, true},
		{5, []byte{1, 1}, 257, true},
		{6, []byte{0, 0, 128, 63}, 1, true}, {7, []byte{0, 0, 128, 63}, 1, true},
	} {
		b524ScalarShapeValue, b524ScalarShapeOK, b524ScalarShapeRaw := b524ScalarShapeEvaluate(b524ScalarShapeEntries[b524ScalarShapeCase.b524ScalarShapeIndex], b524ScalarShapeCase.b524ScalarShapeRaw)
		if b524ScalarShapeValue != b524ScalarShapeCase.b524ScalarShapeWant || b524ScalarShapeOK != b524ScalarShapeCase.b524ScalarShapeOK || !b524ScalarShapeBytes.Equal(b524ScalarShapeRaw, b524ScalarShapeCase.b524ScalarShapeRaw) {
			b524ScalarShapeT.Fatalf("ab85 synthetic case %+v got %v %v %x", b524ScalarShapeCase, b524ScalarShapeValue, b524ScalarShapeOK, b524ScalarShapeRaw)
		}
	}
	for _, b524ScalarShapeEntry := range b524ScalarShapeEntries {
		for _, b524ScalarShapeLength := range []int{0, *b524ScalarShapeEntry.B524ScalarShapeWidthBytes - 1, *b524ScalarShapeEntry.B524ScalarShapeWidthBytes + 1, 6} {
			if b524ScalarShapeLength == *b524ScalarShapeEntry.B524ScalarShapeWidthBytes {
				continue
			}
			b524ScalarShapeInput := b524ScalarShapeBytes.Repeat([]byte{1}, b524ScalarShapeLength)
			_, b524ScalarShapeOK, b524ScalarShapeRaw := b524ScalarShapeEvaluate(b524ScalarShapeEntry, b524ScalarShapeInput)
			if b524ScalarShapeOK || !b524ScalarShapeBytes.Equal(b524ScalarShapeRaw, b524ScalarShapeInput) {
				b524ScalarShapeT.Fatal("ab85 short/extra bytes gained typed claim or lost raw")
			}
		}
		b524ScalarShapeMatches, b524ScalarShapeError := B524ScalarShapeLookup(b524ScalarShapeEntry.B524ScalarShapeKey)
		if b524ScalarShapeError != nil || len(b524ScalarShapeMatches) == 0 {
			b524ScalarShapeT.Fatal("ab85 exact identity lost", b524ScalarShapeError)
		}
		for _, b524ScalarShapeMatch := range b524ScalarShapeMatches {
			if b524ScalarShapeMatch.B524ScalarShapeGrade != "Hypothesis" {
				b524ScalarShapeT.Fatal("ab85 length promoted candidate grade")
			}
		}
	}
	if b524ScalarShapeReflect.DeepEqual(b524ScalarShapeEntries[2].B524ScalarShapeKey, b524ScalarShapeEntries[3].B524ScalarShapeKey) || b524ScalarShapeReflect.DeepEqual(b524ScalarShapeEntries[6].B524ScalarShapeKey, b524ScalarShapeEntries[7].B524ScalarShapeKey) {
		b524ScalarShapeT.Fatal("ab85 separate field roles merged")
	}
}

func b524ScalarShapeMutateShape(b524ScalarShapeShape *B524ScalarShape) {
	*b524ScalarShapeShape.B524ScalarShapeKey.B524ScalarShapePrimary = 0
	*b524ScalarShapeShape.B524ScalarShapeKey.B524ScalarShapeSecondary = 0
	*b524ScalarShapeShape.B524ScalarShapeKey.B524ScalarShapeDirection = "ab85_mutated"
	*b524ScalarShapeShape.B524ScalarShapeKey.B524ScalarShapeOpcode = 6
	*b524ScalarShapeShape.B524ScalarShapeKey.B524ScalarShapeGroup = 0
	*b524ScalarShapeShape.B524ScalarShapeKey.B524ScalarShapeRegister = 0
	*b524ScalarShapeShape.B524ScalarShapeKey.B524ScalarShapeField = "ab85_mutated"
	if b524ScalarShapeShape.B524ScalarShapeKey.B524ScalarShapeFamily != nil {
		*b524ScalarShapeShape.B524ScalarShapeKey.B524ScalarShapeFamily = "ab85_mutated"
	}
	*b524ScalarShapeShape.B524ScalarShapeWidthBytes = 99
	b524ScalarShapeShape.B524ScalarShapeConversion.B524ScalarShapeOffset = 273.15
	b524ScalarShapeShape.B524ScalarShapeEnumeration = map[string]string{"ab85_mutated": "ab85_mutated"}
	b524ScalarShapeShape.B524ScalarShapeReplacementValues = []string{"ab85_mutated"}
	b524ScalarShapeShape.B524ScalarShapeReservedValues = []string{"ab85_mutated"}
	b524ScalarShapeShape.B524ScalarShapeBounds = &B524ScalarShapeBounds{B524ScalarShapeMinimum: b524ScalarShapePointer(99.0)}
}

func TestAb85DetachedCopies(b524ScalarShapeT *b524ScalarShapeTesting.T) {
	b524ScalarShapeBefore := b524ScalarShapeCatalog(b524ScalarShapeT)
	b524ScalarShapeCopy := b524ScalarShapeCatalog(b524ScalarShapeT)
	b524ScalarShapeCopy.B524ScalarShapeLimitations[0] = "ab85_mutated"
	for b524ScalarShapeIndex := range b524ScalarShapeCopy.B524ScalarShapeEntries {
		b524ScalarShapeMutateShape(&b524ScalarShapeCopy.B524ScalarShapeEntries[b524ScalarShapeIndex])
	}
	if !b524ScalarShapeReflect.DeepEqual(b524ScalarShapeCatalog(b524ScalarShapeT), b524ScalarShapeBefore) {
		b524ScalarShapeT.Fatal("ab85 accessor aliases nested data")
	}
	for _, b524ScalarShapeExpected := range b524ScalarShapeExpectedEntries() {
		b524ScalarShapeFirst, b524ScalarShapeError := B524ScalarShapeLookup(b524ScalarShapeExpected.B524ScalarShapeKey)
		if b524ScalarShapeError != nil {
			b524ScalarShapeT.Fatal(b524ScalarShapeError)
		}
		b524ScalarShapeSecond, b524ScalarShapeError := B524ScalarShapeLookup(b524ScalarShapeExpected.B524ScalarShapeKey)
		if b524ScalarShapeError != nil {
			b524ScalarShapeT.Fatal(b524ScalarShapeError)
		}
		for b524ScalarShapeIndex := range b524ScalarShapeFirst {
			b524ScalarShapeMutateShape(&b524ScalarShapeFirst[b524ScalarShapeIndex])
		}
		b524ScalarShapeAgain, b524ScalarShapeError := B524ScalarShapeLookup(b524ScalarShapeExpected.B524ScalarShapeKey)
		if b524ScalarShapeError != nil || !b524ScalarShapeReflect.DeepEqual(b524ScalarShapeAgain, b524ScalarShapeSecond) {
			b524ScalarShapeT.Fatal("ab85 lookup aliases nested data", b524ScalarShapeError)
		}
	}
}

func TestAb85MetadataDoesNotEnableInvocation(b524ScalarShapeT *b524ScalarShapeTesting.T) {
	b524ScalarShapeProvider, b524ScalarShapeError := b524ScalarShapeStandard.NewProvider(b524ScalarShapeStandard.MustEmbeddedCatalog(), true)
	if b524ScalarShapeError != nil {
		b524ScalarShapeT.Fatal(b524ScalarShapeError)
	}
	for b524ScalarShapePass := 0; b524ScalarShapePass < 2; b524ScalarShapePass++ {
		if b524ScalarShapePass == 1 {
			_ = b524ScalarShapeCatalog(b524ScalarShapeT)
		}
		b524ScalarShapeResult, b524ScalarShapeError := b524ScalarShapeProvider.Invoke(b524ScalarShapeContext.Background(), "ebus_standard.service_data.start_counts", nil, b524ScalarShapeStandard.CallerContextUserFacing)
		if b524ScalarShapeError != nil || b524ScalarShapeResult["executed"] != false {
			b524ScalarShapeT.Fatal("ab85 existing stub changed", b524ScalarShapeResult, b524ScalarShapeError)
		}
		_, b524ScalarShapeError = b524ScalarShapeProvider.Invoke(b524ScalarShapeContext.Background(), "ab85_b524_scalar_shapes", nil, b524ScalarShapeStandard.CallerContextUserFacing)
		if !b524ScalarShapeErrors.Is(b524ScalarShapeError, b524ScalarShapeStandard.ErrUnknownMethod) {
			b524ScalarShapeT.Fatal("ab85 metadata gained method registration", b524ScalarShapeError)
		}
	}
}
