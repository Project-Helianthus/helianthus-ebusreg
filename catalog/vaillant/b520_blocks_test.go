package vaillant

import (
	b520BlockReflect "reflect"
	b520BlockTesting "testing"
)

func b520BlockPtr[B520BlockValue any](b520BlockValue B520BlockValue) *B520BlockValue {
	return &b520BlockValue
}

func b520BlockKey(b520BlockOperation uint16, b520BlockAlternative, b520BlockName string) B520BlockKey {
	b520BlockDirection := "response"
	if b520BlockAlternative == "request" {
		b520BlockDirection = "request"
	}
	return B520BlockKey{B520BlockFamily: "Vaillant", B520BlockVariant: "Unknown", B520BlockPrimary: 181,
		B520BlockSecondary: b520BlockPtr(uint16(32)), B520BlockOperation: b520BlockPtr(b520BlockOperation),
		B520BlockDirection: b520BlockDirection, B520BlockAlternative: b520BlockAlternative, B520BlockField: b520BlockName}
}

func b520BlockGet(b520BlockT *b520BlockTesting.T, b520BlockOperation uint16, b520BlockAlternative, b520BlockName string) B520BlockB520Field {
	b520BlockT.Helper()
	b520BlockField, b520BlockFound, b520BlockError := B520BlockLookup(b520BlockKey(b520BlockOperation, b520BlockAlternative, b520BlockName))
	if b520BlockError != nil || !b520BlockFound {
		b520BlockT.Fatalf("ab92 lookup %d/%s/%s: %v %v", b520BlockOperation, b520BlockAlternative, b520BlockName, b520BlockFound, b520BlockError)
	}
	return b520BlockField
}

func TestAb92ProductionMetadata(b520BlockT *b520BlockTesting.T) {
	b520BlockCatalog, b520BlockError := B520Blocks()
	if b520BlockError != nil {
		b520BlockT.Fatal(b520BlockError)
	}
	if b520BlockCatalog.B520BlockRepresentationGrade != "Proven" || b520BlockCatalog.B520BlockPhysicalQualification != "Unknown" ||
		b520BlockCatalog.B520BlockFamilyCompleteness != "incomplete" || b520BlockCatalog.B520BlockInvocation != "unchanged" ||
		b520BlockCatalog.B520BlockSelection != "explicit-key-only" || b520BlockCatalog.B520BlockNegativeForm != "object-not-physical-NAK" ||
		b520BlockCatalog.B520BlockLogicalPositions != "include-prefix;not-physical-offsets" {
		b520BlockT.Fatal("ab92 admission/representation boundary")
	}
	b520BlockUnknowns := map[string]string{"status_semantics": "Unknown", "block_correlation": "Unknown", "atomicity": "Unknown", "physical_adaptation": "Unknown", "effects": "Unknown", "support": "Unknown", "freshness": "Unknown", "persistence": "Unknown", "execution_permission": "Unknown"}
	if !b520BlockReflect.DeepEqual(b520BlockCatalog.B520BlockUnknowns, b520BlockUnknowns) {
		b520BlockT.Fatal("ab92 unknown claims")
	}
	b520BlockExpected := []B520BlockB520Field{}
	b520BlockAdd := func(b520BlockOperation uint16, b520BlockAlternative, b520BlockName string, b520BlockPosition, b520BlockWidth int, b520BlockMax uint64, b520BlockReserved bool) {
		b520BlockOrder := "not-applicable"
		if b520BlockWidth > 1 {
			b520BlockOrder = "little-endian"
		}
		b520BlockExpected = append(b520BlockExpected, B520BlockB520Field{B520BlockKey: b520BlockKey(b520BlockOperation, b520BlockAlternative, b520BlockName), B520BlockPrefix: []int{181, 32, int(b520BlockOperation)},
			B520BlockLogicalPosition: b520BlockPtr(b520BlockPosition), B520BlockWidthBytes: b520BlockPtr(b520BlockWidth), B520BlockSignedness: "unsigned", B520BlockByteOrder: b520BlockOrder,
			B520BlockConversion: "unchanged", B520BlockDefault: b520BlockPtr(uint64(0)), B520BlockMinimum: b520BlockPtr(uint64(0)), B520BlockMaximum: b520BlockPtr(b520BlockMax),
			B520BlockEnumeration: map[string]string{}, B520BlockReserved: b520BlockReserved, B520BlockGrade: "Hypothesis", B520BlockSemantics: "Unknown", B520BlockRawFallback: "preserve-raw-number"})
	}
	b520BlockData := func(b520BlockOperation uint16, b520BlockAlternative string, b520BlockPosition int, b520BlockMinimum, b520BlockMaximum *int, b520BlockDefault *uint64) {
		b520BlockExpected = append(b520BlockExpected, B520BlockB520Field{B520BlockKey: b520BlockKey(b520BlockOperation, b520BlockAlternative, "data"), B520BlockPrefix: []int{181, 32, int(b520BlockOperation)},
			B520BlockLogicalPosition: b520BlockPtr(b520BlockPosition), B520BlockSignedness: "not-applicable", B520BlockByteOrder: "not-applicable", B520BlockConversion: "unchanged",
			B520BlockDefault: b520BlockDefault, B520BlockMinimumBytes: b520BlockMinimum, B520BlockMaximumBytes: b520BlockMaximum, B520BlockExtent: b520BlockPtr("through-end-of-PDU"),
			B520BlockEnumeration: map[string]string{}, B520BlockGrade: "Hypothesis", B520BlockSemantics: "Unknown", B520BlockRawFallback: "preserve-raw-bytes"})
	}
	for _, b520BlockAlternative := range []string{"request", "positive"} {
		b520BlockAdd(1, b520BlockAlternative, "status", 3, 1, 255, false)
	}
	b520BlockAdd(2, "request", "status", 3, 1, 255, false)
	b520BlockAdd(2, "request", "address", 4, 4, 4294967295, false)
	b520BlockAdd(2, "request", "flags", 8, 1, 255, false)
	b520BlockAdd(2, "request", "size-bits", 9, 1, 255, false)
	b520BlockAdd(2, "request", "block", 10, 1, 254, false)
	b520BlockAdd(2, "request", "start-bit", 11, 1, 255, false)
	b520BlockAdd(2, "positive", "status", 3, 1, 255, false)
	for _, b520BlockAlternative := range []string{"request", "positive"} {
		b520BlockAdd(3, b520BlockAlternative, "status", 3, 1, 255, false)
		b520BlockAdd(3, b520BlockAlternative, "block", 4, 1, 254, false)
	}
	b520BlockData(3, "positive", 5, b520BlockPtr(0), nil, nil)
	b520BlockAdd(4, "request", "status", 3, 1, 255, false)
	b520BlockAdd(4, "request", "block", 4, 1, 255, false)
	b520BlockData(4, "request", 5, nil, nil, nil)
	b520BlockAdd(4, "positive", "status", 3, 1, 255, false)
	b520BlockAdd(5, "request", "status", 3, 1, 255, false)
	for b520BlockIndex, b520BlockName := range []string{"status", "version", "max-block-length", "max-variable-count", "current-variable-count", "max-block-count", "current-block-count", "error-code"} {
		b520BlockAdd(5, "positive", b520BlockName, 3+b520BlockIndex, 1, 255, false)
	}
	for _, b520BlockAlternative := range []string{"request", "positive"} {
		b520BlockMax := uint64(254)
		if b520BlockAlternative == "positive" {
			b520BlockMax = 255
		}
		b520BlockAdd(6, b520BlockAlternative, "status", 3, 1, 255, false)
		b520BlockAdd(6, b520BlockAlternative, "write-block", 4, 1, b520BlockMax, false)
		b520BlockAdd(6, b520BlockAlternative, "read-block", 5, 1, b520BlockMax, false)
		b520BlockData(6, b520BlockAlternative, 6, nil, nil, nil)
	}
	for _, b520BlockAlternative := range []string{"request", "positive"} {
		b520BlockAdd(7, b520BlockAlternative, "reserved-identifier", 3, 4, 4294967295, true)
		b520BlockAdd(7, b520BlockAlternative, "reserved-offset", 7, 2, 65535, true)
		b520BlockData(7, b520BlockAlternative, 9, b520BlockPtr(0), b520BlockPtr(23), b520BlockPtr(uint64(0)))
	}
	for _, b520BlockOperation := range []uint16{2, 3, 4, 6} {
		b520BlockAdd(b520BlockOperation, "negative", "status", 3, 1, 255, false)
		if b520BlockOperation == 6 {
			b520BlockAdd(6, "negative", "write-error", 4, 1, 255, false)
			b520BlockAdd(6, "negative", "read-error", 5, 1, 255, false)
		} else {
			b520BlockAdd(b520BlockOperation, "negative", "error", 4, 1, 255, false)
		}
	}
	if len(b520BlockExpected) != 50 || len(b520BlockCatalog.B520BlockFields) != 50 {
		b520BlockT.Fatal("ab92 exact finite field count")
	}
	for _, b520BlockWant := range b520BlockExpected {
		b520BlockGot, b520BlockFound, b520BlockError := B520BlockLookup(b520BlockWant.B520BlockKey)
		if b520BlockError != nil || !b520BlockFound || !b520BlockReflect.DeepEqual(b520BlockGot, b520BlockWant) {
			b520BlockT.Errorf("ab92 field mismatch %+v: got %+v, found %v, error %v", b520BlockWant.B520BlockKey, b520BlockGot, b520BlockFound, b520BlockError)
		}
	}
}

func b520BlockFits(b520BlockField B520BlockB520Field, b520BlockValue uint64) bool {
	return b520BlockField.B520BlockMinimum != nil && b520BlockField.B520BlockMaximum != nil && b520BlockValue >= *b520BlockField.B520BlockMinimum && b520BlockValue <= *b520BlockField.B520BlockMaximum
}

func b520BlockRead(b520BlockT *b520BlockTesting.T, b520BlockField B520BlockB520Field, b520BlockRaw []byte) uint64 {
	b520BlockT.Helper()
	if b520BlockField.B520BlockLogicalPosition == nil || b520BlockField.B520BlockWidthBytes == nil {
		b520BlockT.Fatal("ab92 no scalar extent")
	}
	b520BlockPosition, b520BlockWidth := *b520BlockField.B520BlockLogicalPosition, *b520BlockField.B520BlockWidthBytes
	if b520BlockPosition < 0 || b520BlockWidth < 1 || b520BlockWidth > 8 || b520BlockPosition+b520BlockWidth > len(b520BlockRaw) {
		b520BlockT.Fatal("ab92 invalid synthetic scalar extent")
	}
	if b520BlockWidth > 1 && b520BlockField.B520BlockByteOrder != "little-endian" {
		b520BlockT.Fatal("ab92 unexpected scalar order")
	}
	b520BlockValue := uint64(0)
	for b520BlockIndex := 0; b520BlockIndex < b520BlockWidth; b520BlockIndex++ {
		b520BlockValue |= uint64(b520BlockRaw[b520BlockPosition+b520BlockIndex]) << uint(8*b520BlockIndex)
	}
	return b520BlockValue
}

func TestAb92SyntheticCases(b520BlockT *b520BlockTesting.T) {
	for _, b520BlockCase := range []struct {
		b520BlockOp                         uint16
		b520BlockAlternative, b520BlockName string
		b520BlockValue                      uint64
		b520BlockWant                       bool
	}{
		{3, "request", "block", 254, true}, {3, "request", "block", 255, false}, {4, "request", "block", 255, true},
		{6, "request", "write-block", 255, false}, {6, "request", "read-block", 255, false}, {6, "positive", "write-block", 255, true}, {6, "positive", "read-block", 255, true},
	} {
		b520BlockField := b520BlockGet(b520BlockT, b520BlockCase.b520BlockOp, b520BlockCase.b520BlockAlternative, b520BlockCase.b520BlockName)
		if b520BlockFits(b520BlockField, b520BlockCase.b520BlockValue) != b520BlockCase.b520BlockWant {
			b520BlockT.Errorf("ab92 domain %+v", b520BlockCase)
		}
		if len(b520BlockField.B520BlockEnumeration) != 0 || b520BlockField.B520BlockReplacement != nil || b520BlockField.B520BlockConversion != "unchanged" {
			b520BlockT.Fatal("ab92 invented sentinel")
		}
	}
	b520BlockAddress := b520BlockGet(b520BlockT, 2, "request", "address")
	if b520BlockRead(b520BlockT, b520BlockAddress, []byte{181, 32, 2, 0, 0x78, 0x56, 0x34, 0x12}) != 0x12345678 || b520BlockAddress.B520BlockSemantics != "Unknown" || b520BlockAddress.B520BlockConversion != "unchanged" {
		b520BlockT.Fatal("ab92 raw address")
	}
	b520BlockRaw := []byte{181, 32, 5, 0, 9, 2, 3, 200, 4, 201, 255}
	for b520BlockIndex, b520BlockName := range []string{"status", "version", "max-block-length", "max-variable-count", "current-variable-count", "max-block-count", "current-block-count", "error-code"} {
		b520BlockField := b520BlockGet(b520BlockT, 5, "positive", b520BlockName)
		b520BlockWant := []uint64{0, 9, 2, 3, 200, 4, 201, 255}[b520BlockIndex]
		if b520BlockRead(b520BlockT, b520BlockField, b520BlockRaw) != b520BlockWant || !b520BlockFits(b520BlockField, b520BlockWant) || b520BlockField.B520BlockConversion != "unchanged" || len(b520BlockField.B520BlockEnumeration) != 0 {
			b520BlockT.Fatal("ab92 independent OP05 raw fields")
		}
	}
	if !b520BlockFits(b520BlockGet(b520BlockT, 4, "request", "block"), 255) {
		b520BlockT.Fatal("ab92 OP05 cannot impose a global maximum")
	}
	for _, b520BlockAlternative := range []string{"request", "positive"} {
		b520BlockField := b520BlockGet(b520BlockT, 7, b520BlockAlternative, "data")
		for _, b520BlockLength := range []int{23, 24} {
			b520BlockRaw := make([]byte, 9+b520BlockLength)
			b520BlockCount := len(b520BlockRaw) - *b520BlockField.B520BlockLogicalPosition
			b520BlockFits := b520BlockCount >= *b520BlockField.B520BlockMinimumBytes && b520BlockCount <= *b520BlockField.B520BlockMaximumBytes
			if b520BlockFits != (b520BlockLength == 23) {
				b520BlockT.Fatal("ab92 OP07 abstract extent")
			}
		}
		for _, b520BlockName := range []string{"reserved-identifier", "reserved-offset"} {
			b520BlockField := b520BlockGet(b520BlockT, 7, b520BlockAlternative, b520BlockName)
			if !b520BlockField.B520BlockReserved || b520BlockField.B520BlockSemantics != "Unknown" || b520BlockField.B520BlockConversion != "unchanged" {
				b520BlockT.Fatal("ab92 reserved is not a pointer")
			}
		}
	}
	b520BlockNegative := []byte{181, 32, 6, 0, 1, 2}
	if b520BlockRead(b520BlockT, b520BlockGet(b520BlockT, 6, "negative", "write-error"), b520BlockNegative) != 1 || b520BlockRead(b520BlockT, b520BlockGet(b520BlockT, 6, "negative", "read-error"), b520BlockNegative) != 2 {
		b520BlockT.Fatal("ab92 independent error bytes")
	}
	b520BlockPositive := b520BlockGet(b520BlockT, 6, "positive", "status")
	b520BlockNegativeField := b520BlockGet(b520BlockT, 6, "negative", "status")
	if !b520BlockReflect.DeepEqual(b520BlockPositive.B520BlockPrefix, b520BlockNegativeField.B520BlockPrefix) || b520BlockSameKey(b520BlockPositive.B520BlockKey, b520BlockNegativeField.B520BlockKey) {
		b520BlockT.Fatal("ab92 prefix does not select branch")
	}
	b520BlockMissing := b520BlockPositive.B520BlockKey
	b520BlockMissing.B520BlockAlternative = ""
	if _, b520BlockFound, _ := B520BlockLookup(b520BlockMissing); b520BlockFound {
		b520BlockT.Fatal("ab92 automatic branch selection")
	}
}

func TestAb92ExactIdentity(b520BlockT *b520BlockTesting.T) {
	b520BlockControl := b520BlockKey(2, "request", "status")
	b520BlockMutations := []struct {
		b520BlockName   string
		b520BlockChange func(*B520BlockKey)
	}{
		{"family", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockFamily = "other" }},
		{"variant", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockVariant = "other" }},
		{"primary", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockPrimary = 182 }},
		{"secondary", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockSecondary = b520BlockPtr(uint16(33)) }},
		{"secondary-absent", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockSecondary = nil }},
		{"secondary-zero", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockSecondary = b520BlockPtr(uint16(0)) }},
		{"operation", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockOperation = b520BlockPtr(uint16(0)) }},
		{"operation-absent", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockOperation = nil }},
		{"direction", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockDirection = "response" }},
		{"alternative", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockAlternative = "negative" }},
		{"field", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockField = "other" }},
		{"opcode-zero", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockOpcode = b520BlockPtr(uint16(0)) }},
		{"group-zero", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockGroup = b520BlockPtr(uint16(0)) }},
		{"instance-zero", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockInstance = b520BlockPtr(uint16(0)) }},
		{"register-zero", func(b520BlockKey *B520BlockKey) { b520BlockKey.B520BlockRegister = b520BlockPtr(uint16(0)) }},
	}
	if _, b520BlockFound, b520BlockError := B520BlockLookup(b520BlockControl); !b520BlockFound || b520BlockError != nil {
		b520BlockT.Fatal("ab92 valid control")
	}
	for _, b520BlockMutation := range b520BlockMutations {
		b520BlockT.Run("ab92-"+b520BlockMutation.b520BlockName, func(b520BlockT *b520BlockTesting.T) {
			b520BlockChanged := b520BlockControl
			b520BlockMutation.b520BlockChange(&b520BlockChanged)
			if _, b520BlockFound, b520BlockError := B520BlockLookup(b520BlockChanged); b520BlockFound || b520BlockError != nil {
				b520BlockT.Fatal("ab92 axis ignored")
			}
		})
	}
	b520BlockB524 := b520BlockControl
	b520BlockB524.B520BlockSecondary = b520BlockPtr(uint16(36))
	b520BlockB524.B520BlockOpcode = b520BlockPtr(uint16(1))
	b520BlockB524.B520BlockGroup = b520BlockPtr(uint16(2))
	b520BlockB524.B520BlockInstance = b520BlockPtr(uint16(3))
	b520BlockB524.B520BlockRegister = b520BlockPtr(uint16(4))
	b520BlockB524Copy := b520BlockB524
	b520BlockB524Copy.B520BlockOpcode = b520BlockPtr(uint16(1))
	if !b520BlockSameKey(b520BlockB524, b520BlockB524Copy) {
		b520BlockT.Fatal("ab92 pointer identity used instead of value")
	}
	for b520BlockAxis := 0; b520BlockAxis < 4; b520BlockAxis++ {
		for _, b520BlockValue := range []*uint16{nil, b520BlockPtr(uint16(0)), b520BlockPtr(uint16(99))} {
			b520BlockChanged := b520BlockB524
			b520BlockAxes := []**uint16{&b520BlockChanged.B520BlockOpcode, &b520BlockChanged.B520BlockGroup, &b520BlockChanged.B520BlockInstance, &b520BlockChanged.B520BlockRegister}
			*b520BlockAxes[b520BlockAxis] = b520BlockValue
			if b520BlockSameKey(b520BlockB524, b520BlockChanged) {
				b520BlockT.Fatalf("ab92 qualified B524 axis %d ignored", b520BlockAxis)
			}
		}
	}
	b520BlockAbsent := b520BlockControl
	b520BlockAbsent.B520BlockOperation = nil
	b520BlockZero := b520BlockAbsent
	b520BlockZero.B520BlockOperation = b520BlockPtr(uint16(0))
	if b520BlockSameKey(b520BlockAbsent, b520BlockZero) {
		b520BlockT.Fatal("ab92 absent operation equals zero")
	}
	b520BlockAbsent = b520BlockControl
	b520BlockAbsent.B520BlockSecondary = nil
	b520BlockZero = b520BlockAbsent
	b520BlockZero.B520BlockSecondary = b520BlockPtr(uint16(0))
	if b520BlockSameKey(b520BlockAbsent, b520BlockZero) {
		b520BlockT.Fatal("ab92 absent secondary equals zero")
	}
}

func TestAb92DeepCopies(b520BlockT *b520BlockTesting.T) {
	b520BlockBefore, b520BlockError := B520Blocks()
	if b520BlockError != nil {
		b520BlockT.Fatal(b520BlockError)
	}
	b520BlockChanged, _ := B520Blocks()
	b520BlockChanged.B520BlockUnknowns["support"] = "Proven"
	for b520BlockIndex := range b520BlockChanged.B520BlockFields {
		b520BlockMutate(&b520BlockChanged.B520BlockFields[b520BlockIndex])
	}
	b520BlockAfter, b520BlockError := B520Blocks()
	if b520BlockError != nil || !b520BlockReflect.DeepEqual(b520BlockBefore, b520BlockAfter) {
		b520BlockT.Fatal("ab92 accessor shares storage")
	}
	for _, b520BlockOriginal := range b520BlockBefore.B520BlockFields {
		b520BlockChanged, b520BlockFound, b520BlockError := B520BlockLookup(b520BlockOriginal.B520BlockKey)
		if !b520BlockFound || b520BlockError != nil {
			b520BlockT.Fatal("ab92 lookup control")
		}
		b520BlockMutate(&b520BlockChanged)
		b520BlockAgain, b520BlockFound, b520BlockError := B520BlockLookup(b520BlockOriginal.B520BlockKey)
		if !b520BlockFound || b520BlockError != nil || !b520BlockReflect.DeepEqual(b520BlockOriginal, b520BlockAgain) {
			b520BlockT.Fatal("ab92 lookup shares storage")
		}
	}
}

func b520BlockMutate(b520BlockField *B520BlockB520Field) {
	b520BlockField.B520BlockPrefix[0] = 0
	*b520BlockField.B520BlockKey.B520BlockSecondary = 0
	*b520BlockField.B520BlockKey.B520BlockOperation = 0
	*b520BlockField.B520BlockLogicalPosition = 99
	if b520BlockField.B520BlockWidthBytes != nil {
		*b520BlockField.B520BlockWidthBytes = 99
	}
	if b520BlockField.B520BlockDefault != nil {
		*b520BlockField.B520BlockDefault = 99
	}
	if b520BlockField.B520BlockMinimum != nil {
		*b520BlockField.B520BlockMinimum = 99
	}
	if b520BlockField.B520BlockMaximum != nil {
		*b520BlockField.B520BlockMaximum = 99
	}
	if b520BlockField.B520BlockMinimumBytes != nil {
		*b520BlockField.B520BlockMinimumBytes = 99
	}
	if b520BlockField.B520BlockMaximumBytes != nil {
		*b520BlockField.B520BlockMaximumBytes = 99
	}
	if b520BlockField.B520BlockExtent != nil {
		*b520BlockField.B520BlockExtent = "changed"
	}
	b520BlockField.B520BlockEnumeration["255"] = "sentinel"
	b520BlockField.B520BlockRawFallback = "changed"
}
