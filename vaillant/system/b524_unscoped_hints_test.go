package system

import (
	b524UnscopedHintBytes "bytes"
	b524UnscopedHintJSON "encoding/json"
	b524UnscopedHintFmt "fmt"
	b524UnscopedHintIO "io"
	b524UnscopedHintOs "os"
	b524UnscopedHintReflect "reflect"
	b524UnscopedHintTesting "testing"
)

type b524UnscopedHintFixture struct {
	B524UnscopedHintKeys     []B524UnscopedHintB524HintKey `json:"ab66Keys"`
	B524UnscopedHintControls []struct {
		B524UnscopedHintKey B524UnscopedHintB524HintKey `json:"Ab66Key"`
		B524UnscopedHint    B524UnscopedHintProperties  `json:"Ab66Hint"`
	} `json:"ab66Controls"`
}

func b524UnscopedHintLoadFixture(b524UnscopedHintT *b524UnscopedHintTesting.T) b524UnscopedHintFixture {
	b524UnscopedHintT.Helper()
	b524UnscopedHintRaw, b524UnscopedHintErr := b524UnscopedHintOs.ReadFile("testdata/ab66_b524_unscoped_hints.json")
	if b524UnscopedHintErr != nil {
		b524UnscopedHintT.Fatal(b524UnscopedHintErr)
	}
	var b524UnscopedHintData b524UnscopedHintFixture
	b524UnscopedHintDecoder := b524UnscopedHintJSON.NewDecoder(b524UnscopedHintBytes.NewReader(b524UnscopedHintRaw))
	b524UnscopedHintDecoder.DisallowUnknownFields()
	if b524UnscopedHintErr = b524UnscopedHintDecoder.Decode(&b524UnscopedHintData); b524UnscopedHintErr != nil {
		b524UnscopedHintT.Fatal(b524UnscopedHintErr)
	}
	return b524UnscopedHintData
}

func TestAb66EnumerationAndLookup(b524UnscopedHintT *b524UnscopedHintTesting.T) {
	b524UnscopedHintData := b524UnscopedHintLoadFixture(b524UnscopedHintT)
	b524UnscopedHintResults := B524UnscopedHints()
	if len(b524UnscopedHintResults) != 46 || len(b524UnscopedHintData.B524UnscopedHintKeys) != 46 || len(b524StaticConstraintsCatalog) != 46 {
		b524UnscopedHintT.Fatalf("ab66 count: %d/%d/%d", len(b524UnscopedHintResults), len(b524UnscopedHintData.B524UnscopedHintKeys), len(b524StaticConstraintsCatalog))
	}
	b524UnscopedHintSeen := make(map[B524UnscopedHintB524HintKey]bool)
	for b524UnscopedHintI, b524UnscopedHintKey := range b524UnscopedHintData.B524UnscopedHintKeys {
		b524UnscopedHintT.Run(b524UnscopedHintFmt.Sprintf("ab66_key_%d", b524UnscopedHintI), func(b524UnscopedHintT *b524UnscopedHintTesting.T) {
			if b524UnscopedHintSeen[b524UnscopedHintKey] {
				b524UnscopedHintT.Fatal("ab66 duplicate")
			}
			b524UnscopedHintSeen[b524UnscopedHintKey] = true
			b524UnscopedHintExpected := b524StaticConstraintsCatalog[b524ConstraintSelector{Opcode: b524UnscopedHintKey.B524UnscopedHintOpcode, Group: b524UnscopedHintKey.B524UnscopedHintGroup, Register: b524UnscopedHintKey.B524UnscopedHintRegister}]
			b524UnscopedHintWant := B524UnscopedHintResult{
				B524UnscopedHintKey: b524UnscopedHintKey, B524UnscopedHintStatus: B524HintUnscoped, B524UnscopedHintQualified: false,
				B524UnscopedHintApplicability: B524UnscopedHintApplicabilityUnknown,
				B524UnscopedHint:              B524UnscopedHintProperties{B524UnscopedHintType: b524UnscopedHintExpected.Type, B524UnscopedHintMin: b524UnscopedHintExpected.Min, B524UnscopedHintMax: b524UnscopedHintExpected.Max, B524UnscopedHintStep: b524UnscopedHintExpected.Step},
			}
			if b524UnscopedHintResults[b524UnscopedHintI] != b524UnscopedHintWant || B524LookupUnscopedHint(b524UnscopedHintKey) != b524UnscopedHintWant {
				b524UnscopedHintT.Fatalf("ab66 enumeration/lookup mismatch: %+v; want %+v", b524UnscopedHintResults[b524UnscopedHintI], b524UnscopedHintWant)
			}
			b524UnscopedHintKey.B524UnscopedHintOpcode = 6
			b524UnscopedHintAssertAbsent(b524UnscopedHintT, b524UnscopedHintKey)
		})
	}
}

func TestAb66IndependentNumericControls(b524UnscopedHintT *b524UnscopedHintTesting.T) {
	b524UnscopedHintData := b524UnscopedHintLoadFixture(b524UnscopedHintT)
	if len(b524UnscopedHintData.B524UnscopedHintControls) != 3 {
		b524UnscopedHintT.Fatal("ab66 control count")
	}
	for b524UnscopedHintI, b524UnscopedHintControl := range b524UnscopedHintData.B524UnscopedHintControls {
		b524UnscopedHintT.Run(b524UnscopedHintFmt.Sprintf("ab66_numeric_%d", b524UnscopedHintI), func(b524UnscopedHintT *b524UnscopedHintTesting.T) {
			b524UnscopedHintGot := B524LookupUnscopedHint(b524UnscopedHintControl.B524UnscopedHintKey)
			if b524UnscopedHintGot.B524UnscopedHint != b524UnscopedHintControl.B524UnscopedHint {
				b524UnscopedHintT.Fatalf("ab66 properties: %+v want %+v", b524UnscopedHintGot, b524UnscopedHintControl)
			}
		})
	}
}

func b524UnscopedHintAssertAbsent(b524UnscopedHintT *b524UnscopedHintTesting.T, b524UnscopedHintKey B524UnscopedHintB524HintKey) {
	b524UnscopedHintT.Helper()
	b524UnscopedHintWant := B524UnscopedHintResult{B524UnscopedHintKey: b524UnscopedHintKey, B524UnscopedHintStatus: B524UnscopedHintAbsent, B524UnscopedHintApplicability: B524UnscopedHintApplicabilityUnknown}
	if b524UnscopedHintGot := B524LookupUnscopedHint(b524UnscopedHintKey); b524UnscopedHintGot != b524UnscopedHintWant {
		b524UnscopedHintT.Fatalf("ab66 absent: %+v want %+v", b524UnscopedHintGot, b524UnscopedHintWant)
	}
}

func TestAb66IndependentKeyAxes(b524UnscopedHintT *b524UnscopedHintTesting.T) {
	for b524UnscopedHintOpcode := 0; b524UnscopedHintOpcode < 256; b524UnscopedHintOpcode++ {
		if b524UnscopedHintOpcode != 2 {
			b524UnscopedHintAssertAbsent(b524UnscopedHintT, B524UnscopedHintB524HintKey{B524UnscopedHintOpcode: byte(b524UnscopedHintOpcode), B524UnscopedHintGroup: 0, B524UnscopedHintRegister: 2})
		}
	}
	for _, b524UnscopedHintKey := range []B524UnscopedHintB524HintKey{
		{B524UnscopedHintOpcode: 2, B524UnscopedHintGroup: 0, B524UnscopedHintRegister: 1},
		{B524UnscopedHintOpcode: 6, B524UnscopedHintGroup: 0, B524UnscopedHintRegister: 1},
		{B524UnscopedHintOpcode: 2, B524UnscopedHintGroup: 0, B524UnscopedHintRegister: 0},
		{B524UnscopedHintOpcode: 2, B524UnscopedHintGroup: 0, B524UnscopedHintRegister: 0x0200},
		{B524UnscopedHintOpcode: 2, B524UnscopedHintGroup: 0, B524UnscopedHintRegister: 0xffff},
		{B524UnscopedHintOpcode: 2, B524UnscopedHintGroup: 0xff, B524UnscopedHintRegister: 2},
		{},
	} {
		b524UnscopedHintAssertAbsent(b524UnscopedHintT, b524UnscopedHintKey)
	}
	for _, b524UnscopedHintKey := range []B524UnscopedHintB524HintKey{
		{B524UnscopedHintOpcode: 2, B524UnscopedHintGroup: 0, B524UnscopedHintRegister: 2},
		{B524UnscopedHintOpcode: 2, B524UnscopedHintGroup: 3, B524UnscopedHintRegister: 2},
		{B524UnscopedHintOpcode: 2, B524UnscopedHintGroup: 0, B524UnscopedHintRegister: 4},
	} {
		if B524LookupUnscopedHint(b524UnscopedHintKey).B524UnscopedHintStatus != B524HintUnscoped {
			b524UnscopedHintT.Fatal("ab66 valid axis control")
		}
	}
}

func TestAb66OwnedValues(b524UnscopedHintT *b524UnscopedHintTesting.T) {
	b524UnscopedHintBefore := B524UnscopedHints()
	b524UnscopedHintResults := B524UnscopedHints()
	b524UnscopedHintKey := b524UnscopedHintResults[0].B524UnscopedHintKey
	b524UnscopedHintResult := B524LookupUnscopedHint(b524UnscopedHintKey)
	b524UnscopedHintKey.B524UnscopedHintOpcode, b524UnscopedHintKey.B524UnscopedHintGroup, b524UnscopedHintKey.B524UnscopedHintRegister = 6, 255, 65535
	if b524UnscopedHintResult != b524UnscopedHintBefore[0] {
		b524UnscopedHintT.Fatal("ab66 input alias")
	}
	b524UnscopedHintResult.B524UnscopedHintKey = b524UnscopedHintKey
	b524UnscopedHintResult.B524UnscopedHintStatus = B524UnscopedHintAbsent
	b524UnscopedHintResult.B524UnscopedHintQualified = true
	b524UnscopedHintResult.B524UnscopedHintApplicability = "ab66_changed"
	b524UnscopedHintResult.B524UnscopedHint = B524UnscopedHintProperties{B524UnscopedHintType: "ab66_changed", B524UnscopedHintMin: 11, B524UnscopedHintMax: 12, B524UnscopedHintStep: 13}
	for b524UnscopedHintI := range b524UnscopedHintResults {
		b524UnscopedHintResults[b524UnscopedHintI] = b524UnscopedHintResult
	}
	b524UnscopedHintResults = append(b524UnscopedHintResults, b524UnscopedHintResult)
	if len(b524UnscopedHintResults) != len(b524UnscopedHintBefore)+1 || b524UnscopedHintResults[len(b524UnscopedHintBefore)] != b524UnscopedHintResult {
		b524UnscopedHintT.Fatal("ab66 caller append")
	}
	if !b524UnscopedHintReflect.DeepEqual(B524UnscopedHints(), b524UnscopedHintBefore) {
		b524UnscopedHintT.Fatal("ab66 enumeration alias")
	}
	if B524LookupUnscopedHint(b524UnscopedHintBefore[0].B524UnscopedHintKey) != b524UnscopedHintBefore[0] {
		b524UnscopedHintT.Fatal("ab66 lookup alias")
	}
	b524UnscopedHintNative, b524UnscopedHintOK := lookupB524Constraint(2, 0, 2)
	if !b524UnscopedHintOK || b524UnscopedHintNative.Type != "f32_range" || b524UnscopedHintNative.Min != -26 || b524UnscopedHintNative.Max != 10 || b524UnscopedHintNative.Step != 1 {
		b524UnscopedHintT.Fatal("ab66 native table changed")
	}
}

func TestAb66PublicShape(b524UnscopedHintT *b524UnscopedHintTesting.T) {
	for _, b524UnscopedHintCase := range []struct {
		B524UnscopedHintValue  any
		B524UnscopedHintFields []string
	}{
		{B524UnscopedHintB524HintKey{}, []string{"B524UnscopedHintOpcode", "B524UnscopedHintGroup", "B524UnscopedHintRegister"}},
		{B524UnscopedHintProperties{}, []string{"B524UnscopedHintType", "B524UnscopedHintMin", "B524UnscopedHintMax", "B524UnscopedHintStep"}},
		{B524UnscopedHintResult{}, []string{"B524UnscopedHintKey", "B524UnscopedHintStatus", "B524UnscopedHintQualified", "B524UnscopedHintApplicability", "B524UnscopedHint"}},
	} {
		b524UnscopedHintType := b524UnscopedHintReflect.TypeOf(b524UnscopedHintCase.B524UnscopedHintValue)
		if b524UnscopedHintType.NumField() != len(b524UnscopedHintCase.B524UnscopedHintFields) {
			b524UnscopedHintT.Fatal("ab66 unexpected identity/bound field")
		}
		for b524UnscopedHintI, b524UnscopedHintName := range b524UnscopedHintCase.B524UnscopedHintFields {
			if b524UnscopedHintType.Field(b524UnscopedHintI).Name != b524UnscopedHintName {
				b524UnscopedHintT.Fatal("ab66 field contract")
			}
		}
	}
	if B524HintUnscoped != "unscoped_catalog" || B524UnscopedHintAbsent != "absent" || B524UnscopedHintApplicabilityUnknown != "Unknown" {
		b524UnscopedHintT.Fatal("ab66 status/grade contract")
	}
}

func b524UnscopedHintValidateSyntheticKeys(b524UnscopedHintRaw []byte) ([]B524UnscopedHintB524HintKey, error) {
	var b524UnscopedHintWire []struct {
		B524UnscopedHintOpcode   *byte   `json:"Ab66Opcode"`
		B524UnscopedHintGroup    *byte   `json:"Ab66Group"`
		B524UnscopedHintRegister *uint16 `json:"Ab66Register"`
	}
	b524UnscopedHintDecoder := b524UnscopedHintJSON.NewDecoder(b524UnscopedHintBytes.NewReader(b524UnscopedHintRaw))
	b524UnscopedHintDecoder.DisallowUnknownFields()
	if b524UnscopedHintErr := b524UnscopedHintDecoder.Decode(&b524UnscopedHintWire); b524UnscopedHintErr != nil {
		return nil, b524UnscopedHintErr
	}
	var b524UnscopedHintExtra any
	if b524UnscopedHintErr := b524UnscopedHintDecoder.Decode(&b524UnscopedHintExtra); b524UnscopedHintErr != b524UnscopedHintIO.EOF {
		return nil, b524UnscopedHintFmt.Errorf("ab66 trailing input")
	}
	if len(b524UnscopedHintWire) == 0 {
		return nil, b524UnscopedHintFmt.Errorf("ab66 missing keys")
	}
	b524UnscopedHintSeen := make(map[B524UnscopedHintB524HintKey]bool)
	b524UnscopedHintKeys := make([]B524UnscopedHintB524HintKey, 0, len(b524UnscopedHintWire))
	for _, b524UnscopedHintItem := range b524UnscopedHintWire {
		if b524UnscopedHintItem.B524UnscopedHintOpcode == nil || b524UnscopedHintItem.B524UnscopedHintGroup == nil || b524UnscopedHintItem.B524UnscopedHintRegister == nil {
			return nil, b524UnscopedHintFmt.Errorf("ab66 missing field")
		}
		b524UnscopedHintKey := B524UnscopedHintB524HintKey{B524UnscopedHintOpcode: *b524UnscopedHintItem.B524UnscopedHintOpcode, B524UnscopedHintGroup: *b524UnscopedHintItem.B524UnscopedHintGroup, B524UnscopedHintRegister: *b524UnscopedHintItem.B524UnscopedHintRegister}
		if b524UnscopedHintSeen[b524UnscopedHintKey] {
			return nil, b524UnscopedHintFmt.Errorf("ab66 duplicate key")
		}
		b524UnscopedHintSeen[b524UnscopedHintKey] = true
		b524UnscopedHintKeys = append(b524UnscopedHintKeys, b524UnscopedHintKey)
	}
	return b524UnscopedHintKeys, nil
}

func TestAb66SyntheticValidationAndRawOwnership(b524UnscopedHintT *b524UnscopedHintTesting.T) {
	for b524UnscopedHintI, b524UnscopedHintRaw := range []string{
		`[`, `null`, `[]`, `[{}]`,
		`[{"Ab66Opcode":2,"Ab66Register":2}]`,
		`[{"Ab66Group":0,"Ab66Register":2}]`,
		`[{"Ab66Opcode":2,"Ab66Group":0}]`,
		`[{"Ab66Opcode":2,"Ab66Group":null,"Ab66Register":2}]`,
		`[{"Ab66Opcode":2,"Ab66Group":0,"Ab66Register":2,"ab66Family":"ab66_synthetic"}]`,
		`[{"Ab66Opcode":256,"Ab66Group":0,"Ab66Register":2}]`,
		`[{"Ab66Opcode":2,"Ab66Group":0,"Ab66Register":2},{"Ab66Opcode":2,"Ab66Group":0,"Ab66Register":2}]`,
		`[{"Ab66Opcode":2,"Ab66Group":0,"Ab66Register":2}] {}`,
	} {
		b524UnscopedHintT.Run(b524UnscopedHintFmt.Sprintf("ab66_invalid_%d", b524UnscopedHintI), func(b524UnscopedHintT *b524UnscopedHintTesting.T) {
			if _, b524UnscopedHintErr := b524UnscopedHintValidateSyntheticKeys([]byte(b524UnscopedHintRaw)); b524UnscopedHintErr == nil {
				b524UnscopedHintT.Fatal("ab66 accepted invalid fixture")
			}
		})
	}
	b524UnscopedHintRaw := []byte(`[{"Ab66Opcode":0,"Ab66Group":0,"Ab66Register":0},{"Ab66Opcode":2,"Ab66Group":0,"Ab66Register":2}]`)
	b524UnscopedHintSaved := append([]byte(nil), b524UnscopedHintRaw...)
	b524UnscopedHintKeys, b524UnscopedHintErr := b524UnscopedHintValidateSyntheticKeys(b524UnscopedHintRaw)
	if b524UnscopedHintErr != nil || len(b524UnscopedHintKeys) != 2 || b524UnscopedHintKeys[0] != (B524UnscopedHintB524HintKey{}) {
		b524UnscopedHintT.Fatalf("ab66 explicit zero: %v", b524UnscopedHintErr)
	}
	if !b524UnscopedHintBytes.Equal(b524UnscopedHintSaved, b524UnscopedHintRaw) {
		b524UnscopedHintT.Fatal("ab66 raw changed")
	}
	for b524UnscopedHintI := range b524UnscopedHintRaw {
		b524UnscopedHintRaw[b524UnscopedHintI] = ' '
	}
	if b524UnscopedHintKeys[1] != (B524UnscopedHintB524HintKey{B524UnscopedHintOpcode: 2, B524UnscopedHintGroup: 0, B524UnscopedHintRegister: 2}) {
		b524UnscopedHintT.Fatal("ab66 raw alias")
	}
	b524UnscopedHintAssertAbsent(b524UnscopedHintT, b524UnscopedHintKeys[0])
	if B524LookupUnscopedHint(b524UnscopedHintKeys[1]).B524UnscopedHintStatus != B524HintUnscoped {
		b524UnscopedHintT.Fatal("ab66 valid raw control")
	}
}
