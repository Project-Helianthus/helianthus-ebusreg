package vaillant

import (
	b514MenuFeedbackReflect "reflect"
	b514MenuFeedbackTesting "testing"
)

func b514MenuFeedbackInt(b514MenuFeedbackValue int) *int          { return &b514MenuFeedbackValue }
func b514MenuFeedbackString(b514MenuFeedbackValue string) *string { return &b514MenuFeedbackValue }
func b514MenuFeedbackAxis(b514MenuFeedbackValue uint16) B514MenuFeedbackAxis {
	return B514MenuFeedbackAxis{true, b514MenuFeedbackValue}
}
func b514MenuFeedbackKey(b514MenuFeedbackOperation uint16, b514MenuFeedbackDirection, b514MenuFeedbackName, b514MenuFeedbackVariant string) B514MenuFeedbackKey {
	b514MenuFeedbackPrimary := uint16(0xb5)
	if b514MenuFeedbackDirection == "positive" {
		b514MenuFeedbackPrimary = 0xf5
	}
	return B514MenuFeedbackKey{B514MenuFeedbackFamily: "B514", B514MenuFeedbackVariant: b514MenuFeedbackVariant, B514MenuFeedbackPrimary: b514MenuFeedbackPrimary, B514MenuFeedbackSecondary: b514MenuFeedbackAxis(0x14), B514MenuFeedbackOperation: b514MenuFeedbackAxis(b514MenuFeedbackOperation), B514MenuFeedbackDirection: b514MenuFeedbackDirection, B514MenuFeedbackField: b514MenuFeedbackName}
}
func b514MenuFeedbackGet(b514MenuFeedbackT *b514MenuFeedbackTesting.T, b514MenuFeedbackKey B514MenuFeedbackKey) B514MenuFeedbackField {
	b514MenuFeedbackT.Helper()
	b514MenuFeedbackField, b514MenuFeedbackFound, b514MenuFeedbackError := B514MenuFeedbackLookup(b514MenuFeedbackKey)
	if b514MenuFeedbackError != nil || !b514MenuFeedbackFound {
		b514MenuFeedbackT.Fatalf("lookup %+v: found=%v error=%v", b514MenuFeedbackKey, b514MenuFeedbackFound, b514MenuFeedbackError)
	}
	return b514MenuFeedbackField
}
func b514MenuFeedbackLabel(b514MenuFeedbackField B514MenuFeedbackField, b514MenuFeedbackValue int) string {
	if b514MenuFeedbackLabel, b514MenuFeedbackFound := b514MenuFeedbackField.B514MenuFeedbackEnumeration[b514MenuFeedbackValue]; b514MenuFeedbackFound {
		return b514MenuFeedbackLabel
	}
	for _, b514MenuFeedbackRange := range b514MenuFeedbackField.B514MenuFeedbackRanges {
		if b514MenuFeedbackValue >= b514MenuFeedbackRange.B514MenuFeedbackMinimum && b514MenuFeedbackValue <= b514MenuFeedbackRange.B514MenuFeedbackMaximum {
			return b514MenuFeedbackRange.B514MenuFeedbackMeaning
		}
	}
	return "unmapped"
}
func b514MenuFeedbackWord(b514MenuFeedbackT *b514MenuFeedbackTesting.T, b514MenuFeedbackField B514MenuFeedbackField, b514MenuFeedbackData []byte) float64 {
	b514MenuFeedbackT.Helper()
	if b514MenuFeedbackField.B514MenuFeedbackPosition == nil || b514MenuFeedbackField.B514MenuFeedbackWidth == nil || b514MenuFeedbackField.B514MenuFeedbackExtent == nil || len(b514MenuFeedbackData) < *b514MenuFeedbackField.B514MenuFeedbackExtent {
		b514MenuFeedbackT.Fatal("incomplete abstract object")
	}
	if *b514MenuFeedbackField.B514MenuFeedbackWidth != 2 || b514MenuFeedbackField.B514MenuFeedbackByteOrder != "little-endian" {
		b514MenuFeedbackT.Fatal("unexpected word layout")
	}
	b514MenuFeedbackPosition := *b514MenuFeedbackField.B514MenuFeedbackPosition
	b514MenuFeedbackValue := float64(uint16(b514MenuFeedbackData[b514MenuFeedbackPosition]) | uint16(b514MenuFeedbackData[b514MenuFeedbackPosition+1])<<8)
	if b514MenuFeedbackField.B514MenuFeedbackConversion != nil {
		b514MenuFeedbackValue = b514MenuFeedbackValue*b514MenuFeedbackField.B514MenuFeedbackConversion.B514MenuFeedbackFactor + b514MenuFeedbackField.B514MenuFeedbackConversion.B514MenuFeedbackOffset
	}
	return b514MenuFeedbackValue
}

func TestAb89ProductionMetadata(b514MenuFeedbackT *b514MenuFeedbackTesting.T) {
	b514MenuFeedbackCatalog, b514MenuFeedbackError := B514MenuFeedback()
	if b514MenuFeedbackError != nil {
		b514MenuFeedbackT.Fatal(b514MenuFeedbackError)
	}
	if b514MenuFeedbackCatalog.B514MenuFeedbackRepresentationGrade != "Proven" || b514MenuFeedbackCatalog.B514MenuFeedbackQualification != "Unknown" || len(b514MenuFeedbackCatalog.B514MenuFeedbackFields) != 21 || len(b514MenuFeedbackCatalog.B514MenuFeedbackLimitations) != 8 {
		b514MenuFeedbackT.Fatal("catalog boundary changed")
	}
	b514MenuFeedbackExpected := []B514MenuFeedbackField{}
	b514MenuFeedbackAdd := func(b514MenuFeedbackOp uint16, b514MenuFeedbackDirection, b514MenuFeedbackName, b514MenuFeedbackVariant string, b514MenuFeedbackExtent *int, b514MenuFeedbackEnum map[int]string) *B514MenuFeedbackField {
		b514MenuFeedbackExpected = append(b514MenuFeedbackExpected, B514MenuFeedbackField{B514MenuFeedbackKey: b514MenuFeedbackKey(b514MenuFeedbackOp, b514MenuFeedbackDirection, b514MenuFeedbackName, b514MenuFeedbackVariant), B514MenuFeedbackExtent: b514MenuFeedbackExtent, B514MenuFeedbackSignedness: "Unknown", B514MenuFeedbackByteOrder: "Unknown", B514MenuFeedbackEnumeration: b514MenuFeedbackEnum, B514MenuFeedbackRanges: []B514MenuFeedbackRange{}, B514MenuFeedbackReserved: []B514MenuFeedbackRange{}, B514MenuFeedbackGrade: "Hypothesis", B514MenuFeedbackRawFallback: "preserve raw bytes and unresolved alternatives"})
		return &b514MenuFeedbackExpected[len(b514MenuFeedbackExpected)-1]
	}
	b514MenuFeedbackAdd(1, "request", "mode", "menu", b514MenuFeedbackInt(4), map[int]string{0: "off", 1: "on"}).B514MenuFeedbackRanges = []B514MenuFeedbackRange{{2, 255, "off"}}
	b514MenuFeedbackAdd(1, "positive", "answer", "menu", b514MenuFeedbackInt(5), map[int]string{0: "normal", 1: "actors-deactivated"})
	b514MenuFeedbackAdd(1, "positive", "feedback", "menu", b514MenuFeedbackInt(5), map[int]string{1: "accepted"})
	for _, b514MenuFeedbackVariant := range []string{"generic", "specific-a", "specific-b"} {
		b514MenuFeedbackAdd(6, "request", "command", b514MenuFeedbackVariant, nil, map[int]string{0: "query", 1: "start"})
		if b514MenuFeedbackVariant == "generic" {
			b514MenuFeedbackAdd(6, "positive", "feedback", b514MenuFeedbackVariant, nil, map[int]string{1: "in-progress", 255: "unsupported"})
		} else {
			b514MenuFeedbackAdd(6, "positive", "feedback", b514MenuFeedbackVariant, nil, map[int]string{1: "not-successful"}).B514MenuFeedbackRanges = []B514MenuFeedbackRange{{3, 255, "invalid"}}
		}
	}
	b514MenuFeedbackAdd(7, "request", "command", "menu", nil, map[int]string{0: "query", 1: "start-liquefied-gas", 2: "start-natural-gas", 3: "abort"}).B514MenuFeedbackDefault = b514MenuFeedbackInt(0)
	b514MenuFeedbackAdd(7, "positive", "feedback", "menu", nil, map[int]string{0: "idle", 1: "progress", 2: "passed", 3: "warning", 4: "combustion-error", 5: "gas-type-error", 6: "aborted", 255: "unsupported"}).B514MenuFeedbackDefault = b514MenuFeedbackInt(0)
	for _, b514MenuFeedbackCase := range []struct {
		b514MenuFeedbackDirection, b514MenuFeedbackName                         string
		b514MenuFeedbackPosition, b514MenuFeedbackWidth, b514MenuFeedbackExtent int
	}{
		{"request", "command", 3, 1, 6}, {"request", "pressure", 4, 2, 6}, {"positive", "feedback", 3, 1, 13}, {"positive", "pressure", 4, 2, 13}, {"positive", "minimum", 6, 2, 13}, {"positive", "maximum", 8, 2, 13}, {"positive", "step", 10, 1, 13}, {"positive", "flow", 11, 2, 13},
	} {
		b514MenuFeedbackField := b514MenuFeedbackAdd(8, b514MenuFeedbackCase.b514MenuFeedbackDirection, b514MenuFeedbackCase.b514MenuFeedbackName, "menu", b514MenuFeedbackInt(b514MenuFeedbackCase.b514MenuFeedbackExtent), map[int]string{})
		b514MenuFeedbackField.B514MenuFeedbackPosition = b514MenuFeedbackInt(b514MenuFeedbackCase.b514MenuFeedbackPosition)
		b514MenuFeedbackField.B514MenuFeedbackWidth = b514MenuFeedbackInt(b514MenuFeedbackCase.b514MenuFeedbackWidth)
		b514MenuFeedbackField.B514MenuFeedbackSignedness = "unsigned"
		b514MenuFeedbackField.B514MenuFeedbackByteOrder = "not applicable"
		if b514MenuFeedbackCase.b514MenuFeedbackWidth == 2 {
			b514MenuFeedbackField.B514MenuFeedbackByteOrder = "little-endian"
		}
		if b514MenuFeedbackCase.b514MenuFeedbackName == "step" {
			b514MenuFeedbackField.B514MenuFeedbackUnit = b514MenuFeedbackString("mbar")
			b514MenuFeedbackField.B514MenuFeedbackDefault = b514MenuFeedbackInt(0)
			b514MenuFeedbackField.B514MenuFeedbackConstant = b514MenuFeedbackInt(1)
			b514MenuFeedbackField.B514MenuFeedbackConflict = b514MenuFeedbackString("constant 1 mbar conflicts with default zero; unresolved")
		}
		if b514MenuFeedbackCase.b514MenuFeedbackName == "flow" {
			b514MenuFeedbackField.B514MenuFeedbackUnit = b514MenuFeedbackString("l/min")
			b514MenuFeedbackField.B514MenuFeedbackConversion = &B514MenuFeedbackConversion{0.01, 0}
		}
	}
	for _, b514MenuFeedbackCase := range []struct {
		b514MenuFeedbackDirection string
		b514MenuFeedbackConstant  int
	}{{"request", 0xb5}, {"positive", 0xf5}} {
		b514MenuFeedbackField := b514MenuFeedbackAdd(0, b514MenuFeedbackCase.b514MenuFeedbackDirection, "prefix", "named-09", nil, map[int]string{})
		b514MenuFeedbackField.B514MenuFeedbackKey.B514MenuFeedbackOperation = B514MenuFeedbackAxis{}
		b514MenuFeedbackField.B514MenuFeedbackKey.B514MenuFeedbackSecondary = B514MenuFeedbackAxis{}
		b514MenuFeedbackField.B514MenuFeedbackPosition = b514MenuFeedbackInt(0)
		b514MenuFeedbackField.B514MenuFeedbackWidth = b514MenuFeedbackInt(1)
		b514MenuFeedbackField.B514MenuFeedbackSignedness = "unsigned"
		b514MenuFeedbackField.B514MenuFeedbackByteOrder = "not applicable"
		b514MenuFeedbackField.B514MenuFeedbackConstant = b514MenuFeedbackInt(b514MenuFeedbackCase.b514MenuFeedbackConstant)
	}
	if !b514MenuFeedbackReflect.DeepEqual(b514MenuFeedbackCatalog.B514MenuFeedbackFields, b514MenuFeedbackExpected) {
		b514MenuFeedbackT.Fatal("production fields differ from independent expected metadata")
	}
	for _, b514MenuFeedbackField := range b514MenuFeedbackExpected {
		if !b514MenuFeedbackReflect.DeepEqual(b514MenuFeedbackGet(b514MenuFeedbackT, b514MenuFeedbackField.B514MenuFeedbackKey), b514MenuFeedbackField) {
			b514MenuFeedbackT.Fatal("lookup differs")
		}
	}
}

func TestAb89SyntheticCases(b514MenuFeedbackT *b514MenuFeedbackTesting.T) {
	for _, b514MenuFeedbackCase := range []struct {
		b514MenuFeedbackOp                                                       uint16
		b514MenuFeedbackDirection, b514MenuFeedbackName, b514MenuFeedbackVariant string
		b514MenuFeedbackValue                                                    int
		b514MenuFeedbackExpected                                                 string
	}{
		{1, "request", "mode", "menu", 2, "off"}, {1, "positive", "answer", "menu", 1, "actors-deactivated"}, {1, "positive", "feedback", "menu", 1, "accepted"}, {1, "positive", "feedback", "menu", 0, "unmapped"},
		{6, "positive", "feedback", "generic", 1, "in-progress"}, {6, "positive", "feedback", "specific-a", 1, "not-successful"}, {6, "positive", "feedback", "specific-b", 1, "not-successful"}, {6, "positive", "feedback", "generic", 255, "unsupported"}, {6, "positive", "feedback", "specific-a", 255, "invalid"}, {6, "positive", "feedback", "specific-b", 255, "invalid"}, {6, "positive", "feedback", "specific-a", 2, "unmapped"},
		{7, "positive", "feedback", "menu", 6, "aborted"}, {7, "positive", "feedback", "menu", 7, "unmapped"},
	} {
		if b514MenuFeedbackLabel(b514MenuFeedbackGet(b514MenuFeedbackT, b514MenuFeedbackKey(b514MenuFeedbackCase.b514MenuFeedbackOp, b514MenuFeedbackCase.b514MenuFeedbackDirection, b514MenuFeedbackCase.b514MenuFeedbackName, b514MenuFeedbackCase.b514MenuFeedbackVariant)), b514MenuFeedbackCase.b514MenuFeedbackValue) != b514MenuFeedbackCase.b514MenuFeedbackExpected {
			b514MenuFeedbackT.Fatalf("synthetic case %+v", b514MenuFeedbackCase)
		}
	}
	b514MenuFeedbackData := []byte{0, 0, 0, 0, 1, 0, 0, 1, 0x34, 0x12, 0, 0x64, 0}
	for _, b514MenuFeedbackCase := range []struct {
		b514MenuFeedbackName     string
		b514MenuFeedbackExpected float64
	}{{"pressure", 1}, {"minimum", 256}, {"maximum", 0x1234}, {"flow", 1}} {
		b514MenuFeedbackField := b514MenuFeedbackGet(b514MenuFeedbackT, b514MenuFeedbackKey(8, "positive", b514MenuFeedbackCase.b514MenuFeedbackName, "menu"))
		if b514MenuFeedbackWord(b514MenuFeedbackT, b514MenuFeedbackField, b514MenuFeedbackData) != b514MenuFeedbackCase.b514MenuFeedbackExpected {
			b514MenuFeedbackT.Fatal(b514MenuFeedbackCase)
		}
		if len(b514MenuFeedbackData[:12]) >= *b514MenuFeedbackField.B514MenuFeedbackExtent {
			b514MenuFeedbackT.Fatal("12-position object admitted")
		}
	}
	b514MenuFeedbackStep := b514MenuFeedbackGet(b514MenuFeedbackT, b514MenuFeedbackKey(8, "positive", "step", "menu"))
	if int(b514MenuFeedbackData[*b514MenuFeedbackStep.B514MenuFeedbackPosition]) != *b514MenuFeedbackStep.B514MenuFeedbackDefault || *b514MenuFeedbackStep.B514MenuFeedbackConstant != 1 || b514MenuFeedbackStep.B514MenuFeedbackConflict == nil {
		b514MenuFeedbackT.Fatal("step conflict lost")
	}
	b514MenuFeedbackNamed := b514MenuFeedbackKey(0, "request", "prefix", "named-09")
	b514MenuFeedbackNamed.B514MenuFeedbackOperation = B514MenuFeedbackAxis{}
	b514MenuFeedbackNamed.B514MenuFeedbackSecondary = B514MenuFeedbackAxis{}
	b514MenuFeedbackField := b514MenuFeedbackGet(b514MenuFeedbackT, b514MenuFeedbackNamed)
	if b514MenuFeedbackField.B514MenuFeedbackKey.B514MenuFeedbackSecondary.B514MenuFeedbackPresent || b514MenuFeedbackField.B514MenuFeedbackKey.B514MenuFeedbackOperation.B514MenuFeedbackPresent || b514MenuFeedbackField.B514MenuFeedbackExtent != nil {
		b514MenuFeedbackT.Fatal("inferred executable named 09")
	}
}

func TestAb89ExactIdentity(b514MenuFeedbackT *b514MenuFeedbackTesting.T) {
	b514MenuFeedbackMutations := []struct {
		b514MenuFeedbackName   string
		b514MenuFeedbackMutate func(*B514MenuFeedbackKey)
	}{
		{"family", func(b514MenuFeedbackKey *B514MenuFeedbackKey) { b514MenuFeedbackKey.B514MenuFeedbackFamily = "other" }},
		{"variant", func(b514MenuFeedbackKey *B514MenuFeedbackKey) { b514MenuFeedbackKey.B514MenuFeedbackVariant = "other" }},
		{"primary", func(b514MenuFeedbackKey *B514MenuFeedbackKey) { b514MenuFeedbackKey.B514MenuFeedbackPrimary++ }},
		{"direction", func(b514MenuFeedbackKey *B514MenuFeedbackKey) {
			b514MenuFeedbackKey.B514MenuFeedbackDirection = "other"
		}},
		{"field", func(b514MenuFeedbackKey *B514MenuFeedbackKey) { b514MenuFeedbackKey.B514MenuFeedbackField = "other" }},
	}
	for _, b514MenuFeedbackAxisName := range []string{"B514MenuFeedbackSecondary", "B514MenuFeedbackOperation", "B514MenuFeedbackOpcode", "B514MenuFeedbackGroup", "B514MenuFeedbackInstance", "B514MenuFeedbackRegister"} {
		b514MenuFeedbackMutations = append(b514MenuFeedbackMutations, struct {
			b514MenuFeedbackName   string
			b514MenuFeedbackMutate func(*B514MenuFeedbackKey)
		}{b514MenuFeedbackAxisName + "-presence", func(b514MenuFeedbackKey *B514MenuFeedbackKey) {
			b514MenuFeedbackAxisValue := b514MenuFeedbackReflect.ValueOf(b514MenuFeedbackKey).Elem().FieldByName(b514MenuFeedbackAxisName).Addr().Interface().(*B514MenuFeedbackAxis)
			b514MenuFeedbackAxisValue.B514MenuFeedbackPresent = !b514MenuFeedbackAxisValue.B514MenuFeedbackPresent
		}}, struct {
			b514MenuFeedbackName   string
			b514MenuFeedbackMutate func(*B514MenuFeedbackKey)
		}{b514MenuFeedbackAxisName + "-value", func(b514MenuFeedbackKey *B514MenuFeedbackKey) {
			b514MenuFeedbackAxisValue := b514MenuFeedbackReflect.ValueOf(b514MenuFeedbackKey).Elem().FieldByName(b514MenuFeedbackAxisName).Addr().Interface().(*B514MenuFeedbackAxis)
			b514MenuFeedbackAxisValue.B514MenuFeedbackValue ^= 0x8000
		}})
	}
	b514MenuFeedbackCatalog, b514MenuFeedbackError := B514MenuFeedback()
	if b514MenuFeedbackError != nil {
		b514MenuFeedbackT.Fatal(b514MenuFeedbackError)
	}
	for _, b514MenuFeedbackControl := range b514MenuFeedbackCatalog.B514MenuFeedbackFields {
		for _, b514MenuFeedbackMutation := range b514MenuFeedbackMutations {
			b514MenuFeedbackT.Run(b514MenuFeedbackControl.B514MenuFeedbackKey.B514MenuFeedbackVariant+"/"+b514MenuFeedbackControl.B514MenuFeedbackKey.B514MenuFeedbackField+"/"+b514MenuFeedbackMutation.b514MenuFeedbackName, func(b514MenuFeedbackT *b514MenuFeedbackTesting.T) {
				b514MenuFeedbackGet(b514MenuFeedbackT, b514MenuFeedbackControl.B514MenuFeedbackKey)
				b514MenuFeedbackChanged := b514MenuFeedbackControl.B514MenuFeedbackKey
				b514MenuFeedbackMutation.b514MenuFeedbackMutate(&b514MenuFeedbackChanged)
				if _, b514MenuFeedbackFound, b514MenuFeedbackError := B514MenuFeedbackLookup(b514MenuFeedbackChanged); b514MenuFeedbackFound || b514MenuFeedbackError != nil {
					b514MenuFeedbackT.Fatal("nonexact key matched", b514MenuFeedbackChanged, b514MenuFeedbackError)
				}
			})
		}
	}
	b514MenuFeedbackB524 := b514MenuFeedbackKey(0, "request", "synthetic", "synthetic")
	b514MenuFeedbackB524.B514MenuFeedbackFamily = "B524"
	b514MenuFeedbackB524.B514MenuFeedbackPrimary = 0xb5
	b514MenuFeedbackB524.B514MenuFeedbackSecondary = b514MenuFeedbackAxis(0)
	b514MenuFeedbackB524.B514MenuFeedbackOpcode = b514MenuFeedbackAxis(3)
	b514MenuFeedbackB524.B514MenuFeedbackGroup = b514MenuFeedbackAxis(4)
	b514MenuFeedbackB524.B514MenuFeedbackInstance = b514MenuFeedbackAxis(5)
	b514MenuFeedbackB524.B514MenuFeedbackRegister = b514MenuFeedbackAxis(6)
	b514MenuFeedbackZero := b514MenuFeedbackB524
	b514MenuFeedbackZero.B514MenuFeedbackOpcode = b514MenuFeedbackAxis(0)
	b514MenuFeedbackZero.B514MenuFeedbackGroup = b514MenuFeedbackAxis(0)
	b514MenuFeedbackZero.B514MenuFeedbackInstance = b514MenuFeedbackAxis(0)
	b514MenuFeedbackZero.B514MenuFeedbackRegister = b514MenuFeedbackAxis(0)
	for _, b514MenuFeedbackControl := range []B514MenuFeedbackKey{b514MenuFeedbackB524, b514MenuFeedbackZero} {
		b514MenuFeedbackSynthetic := []B514MenuFeedbackField{{B514MenuFeedbackKey: b514MenuFeedbackControl}}
		if _, b514MenuFeedbackFound, _ := b514MenuFeedbackLookupFields(b514MenuFeedbackSynthetic, b514MenuFeedbackControl); !b514MenuFeedbackFound {
			b514MenuFeedbackT.Fatal("synthetic B524 control missing")
		}
		for _, b514MenuFeedbackMutation := range b514MenuFeedbackMutations {
			b514MenuFeedbackChanged := b514MenuFeedbackControl
			b514MenuFeedbackMutation.b514MenuFeedbackMutate(&b514MenuFeedbackChanged)
			if _, b514MenuFeedbackFound, _ := b514MenuFeedbackLookupFields(b514MenuFeedbackSynthetic, b514MenuFeedbackChanged); b514MenuFeedbackFound {
				b514MenuFeedbackT.Fatal("B524 axis ignored", b514MenuFeedbackMutation.b514MenuFeedbackName)
			}
		}
		if _, b514MenuFeedbackFound, _ := B514MenuFeedbackLookup(b514MenuFeedbackControl); b514MenuFeedbackFound {
			b514MenuFeedbackT.Fatal("synthetic B524 became production declaration")
		}
	}

}

func b514MenuFeedbackMutateField(b514MenuFeedbackField *B514MenuFeedbackField) {
	if b514MenuFeedbackField.B514MenuFeedbackPosition != nil {
		*b514MenuFeedbackField.B514MenuFeedbackPosition = 99
	}
	if b514MenuFeedbackField.B514MenuFeedbackWidth != nil {
		*b514MenuFeedbackField.B514MenuFeedbackWidth = 99
	}
	if b514MenuFeedbackField.B514MenuFeedbackExtent != nil {
		*b514MenuFeedbackField.B514MenuFeedbackExtent = 99
	}
	if b514MenuFeedbackField.B514MenuFeedbackUnit != nil {
		*b514MenuFeedbackField.B514MenuFeedbackUnit = "mutated"
	}
	if b514MenuFeedbackField.B514MenuFeedbackDefault != nil {
		*b514MenuFeedbackField.B514MenuFeedbackDefault = 99
	}
	if b514MenuFeedbackField.B514MenuFeedbackConstant != nil {
		*b514MenuFeedbackField.B514MenuFeedbackConstant = 99
	}
	if b514MenuFeedbackField.B514MenuFeedbackConflict != nil {
		*b514MenuFeedbackField.B514MenuFeedbackConflict = "mutated"
	}
	if b514MenuFeedbackField.B514MenuFeedbackConversion != nil {
		b514MenuFeedbackField.B514MenuFeedbackConversion.B514MenuFeedbackFactor = 99
		b514MenuFeedbackField.B514MenuFeedbackConversion.B514MenuFeedbackOffset = 99
	}
	b514MenuFeedbackField.B514MenuFeedbackEnumeration[999] = "mutated"
	if len(b514MenuFeedbackField.B514MenuFeedbackRanges) > 0 {
		b514MenuFeedbackField.B514MenuFeedbackRanges[0].B514MenuFeedbackMeaning = "mutated"
	}
	b514MenuFeedbackField.B514MenuFeedbackRanges = append(b514MenuFeedbackField.B514MenuFeedbackRanges, B514MenuFeedbackRange{99, 99, "mutated"})
	b514MenuFeedbackField.B514MenuFeedbackReserved = append(b514MenuFeedbackField.B514MenuFeedbackReserved, B514MenuFeedbackRange{99, 99, "mutated"})
	b514MenuFeedbackField.B514MenuFeedbackBounds = append(b514MenuFeedbackField.B514MenuFeedbackBounds, 99)
	b514MenuFeedbackField.B514MenuFeedbackKey.B514MenuFeedbackFamily = "mutated"
}

func TestAb89DetachedCopies(b514MenuFeedbackT *b514MenuFeedbackTesting.T) {
	b514MenuFeedbackOriginal, b514MenuFeedbackError := B514MenuFeedback()
	if b514MenuFeedbackError != nil {
		b514MenuFeedbackT.Fatal(b514MenuFeedbackError)
	}
	b514MenuFeedbackChanged, b514MenuFeedbackError := B514MenuFeedback()
	if b514MenuFeedbackError != nil {
		b514MenuFeedbackT.Fatal(b514MenuFeedbackError)
	}
	b514MenuFeedbackChanged.B514MenuFeedbackLimitations[0] = "mutated"
	for b514MenuFeedbackIndex := range b514MenuFeedbackChanged.B514MenuFeedbackFields {
		b514MenuFeedbackMutateField(&b514MenuFeedbackChanged.B514MenuFeedbackFields[b514MenuFeedbackIndex])
	}
	b514MenuFeedbackFresh, b514MenuFeedbackError := B514MenuFeedback()
	if b514MenuFeedbackError != nil || !b514MenuFeedbackReflect.DeepEqual(b514MenuFeedbackFresh, b514MenuFeedbackOriginal) {
		b514MenuFeedbackT.Fatal("accessor shares backing storage", b514MenuFeedbackError)
	}
	for _, b514MenuFeedbackField := range b514MenuFeedbackOriginal.B514MenuFeedbackFields {
		b514MenuFeedbackLookup := b514MenuFeedbackGet(b514MenuFeedbackT, b514MenuFeedbackField.B514MenuFeedbackKey)
		b514MenuFeedbackMutateField(&b514MenuFeedbackLookup)
		if !b514MenuFeedbackReflect.DeepEqual(b514MenuFeedbackGet(b514MenuFeedbackT, b514MenuFeedbackField.B514MenuFeedbackKey), b514MenuFeedbackField) {
			b514MenuFeedbackT.Fatal("lookup shares backing storage")
		}
	}
}
