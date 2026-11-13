package vaillant

import (
	b524NumericCandidateBytes "bytes"
	b524NumericCandidateJSON "encoding/json"
	b524NumericCandidateReflect "reflect"
	b524NumericCandidateTesting "testing"
)

func b524NumericCandidatePointer[b524NumericCandidateType any](b524NumericCandidateValue b524NumericCandidateType) *b524NumericCandidateType {
	return &b524NumericCandidateValue
}

func b524NumericCandidateTestContext() B524NumericCandidateContext {
	return B524NumericCandidateContext{"ab56_family", "ab56_software", "ab56_hardware", "ab56_api_version", "ab56_api_revision"}
}

func b524NumericCandidateTestKey(b524NumericCandidateGroup, b524NumericCandidateInstance uint8, b524NumericCandidateRegister uint16) B524NumericCandidateKey {
	return B524NumericCandidateKey{
		B524NumericCandidateFamily: "ab56_family", B524NumericCandidateSoftware: "ab56_software", B524NumericCandidateHardware: "ab56_hardware",
		B524NumericCandidateAPIVersion: "ab56_api_version", B524NumericCandidateAPIRevision: "ab56_api_revision",
		B524NumericCandidateVariant: "ab56_numeric_candidates", B524NumericCandidateSet: "ab56_b524_twenty_unadmitted",
		B524NumericCandidatePrimary: b524NumericCandidatePointer(uint16(181)), B524NumericCandidateSecondary: b524NumericCandidatePointer(uint8(36)),
		B524NumericCandidateOperation: b524NumericCandidatePointer(uint8(2)), B524NumericCandidateOpcode: b524NumericCandidatePointer(uint8(2)),
		B524NumericCandidateDirection: "none", B524NumericCandidateField: "ab56_numeric_candidate", B524NumericCandidateAlternative: "none",
		B524NumericCandidateGroup: b524NumericCandidatePointer(b524NumericCandidateGroup), B524NumericCandidateInstanceScope: "ab56_candidate_slot",
		B524NumericCandidateInstance: b524NumericCandidatePointer(b524NumericCandidateInstance), B524NumericCandidateRegister: b524NumericCandidatePointer(b524NumericCandidateRegister),
	}
}

func b524NumericCandidateDecode(b524NumericCandidateT *b524NumericCandidateTesting.T, b524NumericCandidateRaw []byte) map[string]b524NumericCandidateJSON.RawMessage {
	b524NumericCandidateT.Helper()
	var b524NumericCandidateObject map[string]b524NumericCandidateJSON.RawMessage
	if b524NumericCandidateError := b524NumericCandidateJSON.Unmarshal(b524NumericCandidateRaw, &b524NumericCandidateObject); b524NumericCandidateError != nil {
		b524NumericCandidateT.Fatal(b524NumericCandidateError)
	}
	return b524NumericCandidateObject
}

func b524NumericCandidateEncode(b524NumericCandidateT *b524NumericCandidateTesting.T, b524NumericCandidateValue any) []byte {
	b524NumericCandidateT.Helper()
	b524NumericCandidateRaw, b524NumericCandidateError := b524NumericCandidateJSON.Marshal(b524NumericCandidateValue)
	if b524NumericCandidateError != nil {
		b524NumericCandidateT.Fatal(b524NumericCandidateError)
	}
	return b524NumericCandidateRaw
}

func b524NumericCandidateEqual(b524NumericCandidateT *b524NumericCandidateTesting.T, b524NumericCandidateActual, b524NumericCandidateExpected []byte) {
	b524NumericCandidateT.Helper()
	var b524NumericCandidateLeft, b524NumericCandidateRight any
	if b524NumericCandidateError := b524NumericCandidateJSON.Unmarshal(b524NumericCandidateActual, &b524NumericCandidateLeft); b524NumericCandidateError != nil {
		b524NumericCandidateT.Fatal(b524NumericCandidateError)
	}
	if b524NumericCandidateError := b524NumericCandidateJSON.Unmarshal(b524NumericCandidateExpected, &b524NumericCandidateRight); b524NumericCandidateError != nil {
		b524NumericCandidateT.Fatal(b524NumericCandidateError)
	}
	if !b524NumericCandidateReflect.DeepEqual(b524NumericCandidateLeft, b524NumericCandidateRight) {
		b524NumericCandidateT.Fatalf("ab56 row differs\ngot: %s\nwant: %s", b524NumericCandidateActual, b524NumericCandidateExpected)
	}
}

func b524NumericCandidateViewForTest(b524NumericCandidateT *b524NumericCandidateTesting.T) B524NumericCandidateView {
	b524NumericCandidateT.Helper()
	b524NumericCandidateView, b524NumericCandidateError := B524NumericCandidateNewView(b524NumericCandidateTestContext())
	if b524NumericCandidateError != nil {
		b524NumericCandidateT.Fatal(b524NumericCandidateError)
	}
	return b524NumericCandidateView
}

func b524NumericCandidateHit(b524NumericCandidateT *b524NumericCandidateTesting.T, b524NumericCandidateView B524NumericCandidateView, b524NumericCandidateKey B524NumericCandidateKey) b524NumericCandidateJSON.RawMessage {
	b524NumericCandidateT.Helper()
	b524NumericCandidateRaw, b524NumericCandidateOK, b524NumericCandidateError := b524NumericCandidateView.B524NumericCandidateLookup(b524NumericCandidateKey)
	if b524NumericCandidateError != nil || !b524NumericCandidateOK || len(b524NumericCandidateRaw) == 0 {
		b524NumericCandidateT.Fatalf("ab56 missing control: %v %v", b524NumericCandidateOK, b524NumericCandidateError)
	}
	return b524NumericCandidateRaw
}

func b524NumericCandidateMiss(b524NumericCandidateT *b524NumericCandidateTesting.T, b524NumericCandidateView B524NumericCandidateView, b524NumericCandidateKey B524NumericCandidateKey) {
	b524NumericCandidateT.Helper()
	b524NumericCandidateRaw, b524NumericCandidateOK, b524NumericCandidateError := b524NumericCandidateView.B524NumericCandidateLookup(b524NumericCandidateKey)
	if b524NumericCandidateError != nil || b524NumericCandidateOK || b524NumericCandidateRaw != nil {
		b524NumericCandidateT.Fatalf("ab56 unexpected match: %s %v %v", b524NumericCandidateRaw, b524NumericCandidateOK, b524NumericCandidateError)
	}
}

func TestAb56FiniteCandidates(b524NumericCandidateT *b524NumericCandidateTesting.T) {
	b524NumericCandidateCases := []struct {
		b524NumericCandidateGroup, b524NumericCandidateInstance                            uint8
		b524NumericCandidateRegister                                                       uint16
		b524NumericCandidateCodec                                                          string
		b524NumericCandidateMinimum, b524NumericCandidateMaximum, b524NumericCandidateStep float64
		b524NumericCandidateRawForm                                                        string
	}{
		{0, 0, 3, "UIN", 0, 12, 1, "00|00|0003|UIN|0|12|1"},
		{0, 0, 128, "EXP", -10, 10, 1, "00|00|0080|EXP|-10.0|10.0|1.0"},
		{1, 0, 1, "UIN", 0, 1, 1, "01|00|0001|UIN|0|1|1"},
		{1, 0, 3, "UIN", 0, 2, 1, "01|00|0003|UIN|0|2|1"},
		{2, 1, 2, "UIN", 0, 4, 1, "02|01|0002|UIN|0|4|1"},
		{2, 1, 4, "EXP", 15, 80, 1, "02|01|0004|EXP|15.0|80.0|1.0"},
		{2, 1, 5, "UCH", 0, 1, 1, "02|01|0005|UCH|0|1|1"},
		{2, 1, 6, "UCH", 0, 1, 1, "02|01|0006|UCH|0|1|1"},
		{2, 2, 2, "UIN", 0, 4, 1, "02|02|0002|UIN|0|4|1"},
		{2, 2, 4, "EXP", 15, 80, 1, "02|02|0004|EXP|15.0|80.0|1.0"},
		{2, 2, 5, "UCH", 0, 1, 1, "02|02|0005|UCH|0|1|1"},
		{2, 2, 6, "UCH", 0, 1, 1, "02|02|0006|UCH|0|1|1"},
		{2, 9, 2, "UIN", 0, 4, 1, "02|09|0002|UIN|0|4|1"},
		{2, 9, 4, "EXP", 15, 80, 1, "02|09|0004|EXP|15.0|80.0|1.0"},
		{2, 9, 5, "UCH", 0, 1, 1, "02|09|0005|UCH|0|1|1"},
		{2, 9, 6, "UCH", 0, 1, 1, "02|09|0006|UCH|0|1|1"},
		{3, 0, 1, "UIN", 0, 2, 1, "03|00|0001|UIN|0|2|1"},
		{3, 0, 2, "EXP", 15, 30, 0.5, "03|00|0002|EXP|15.0|30.0|0.5"},
		{3, 0, 5, "EXP", 5, 30, 1, "03|00|0005|EXP|5.0|30.0|1.0"},
		{3, 0, 6, "UIN", 0, 2, 1, "03|00|0006|UIN|0|2|1"},
	}
	b524NumericCandidateEntries, b524NumericCandidateError := B524NumericCandidateEntries()
	if b524NumericCandidateError != nil {
		b524NumericCandidateT.Fatal(b524NumericCandidateError)
	}
	var b524NumericCandidateRows []b524NumericCandidateJSON.RawMessage
	if b524NumericCandidateError := b524NumericCandidateJSON.Unmarshal(b524NumericCandidateEntries, &b524NumericCandidateRows); b524NumericCandidateError != nil {
		b524NumericCandidateT.Fatal(b524NumericCandidateError)
	}
	if len(b524NumericCandidateRows) != 20 {
		b524NumericCandidateT.Fatalf("ab56 count = %d", len(b524NumericCandidateRows))
	}
	b524NumericCandidateView := b524NumericCandidateViewForTest(b524NumericCandidateT)
	b524NumericCandidateSeen := make(map[string]bool)
	for _, b524NumericCandidateCase := range b524NumericCandidateCases {
		b524NumericCandidateKey := b524NumericCandidateTestKey(b524NumericCandidateCase.b524NumericCandidateGroup, b524NumericCandidateCase.b524NumericCandidateInstance, b524NumericCandidateCase.b524NumericCandidateRegister)
		b524NumericCandidateRaw := b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateKey)
		b524NumericCandidateKey.B524NumericCandidateFamily, b524NumericCandidateKey.B524NumericCandidateSoftware, b524NumericCandidateKey.B524NumericCandidateHardware, b524NumericCandidateKey.B524NumericCandidateAPIVersion, b524NumericCandidateKey.B524NumericCandidateAPIRevision = "", "", "", "", ""
		b524NumericCandidateExpected := map[string]any{
			"key": b524NumericCandidateKey, "codec": b524NumericCandidateCase.b524NumericCandidateCodec, "width_bytes": nil, "width_grade": "Unknown",
			"minimum": b524NumericCandidateCase.b524NumericCandidateMinimum, "maximum": b524NumericCandidateCase.b524NumericCandidateMaximum, "step": b524NumericCandidateCase.b524NumericCandidateStep,
			"grade": "Hypothesis", "physical_applicability": "Unknown", "qualification": "Unknown", "admitted": false,
			"presence": "Unknown", "freshness": "Unknown", "unresolved": []string{"native identity applicability", "physical applicability", "field byte width", "slot presence"},
			"raw_form": b524NumericCandidateCase.b524NumericCandidateRawForm, "raw_fallback": "preserve native bytes; metadata arithmetic grants no execution or edit permission",
		}
		b524NumericCandidateEqual(b524NumericCandidateT, b524NumericCandidateRaw, b524NumericCandidateEncode(b524NumericCandidateT, b524NumericCandidateExpected))
		if b524NumericCandidateSeen[b524NumericCandidateCase.b524NumericCandidateRawForm] {
			b524NumericCandidateT.Fatal("ab56 duplicate expected selector")
		}
		b524NumericCandidateSeen[b524NumericCandidateCase.b524NumericCandidateRawForm] = true
	}
	b524NumericCandidateMiss(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateTestKey(0, 0, 1))
	for _, b524NumericCandidateKey := range []B524NumericCandidateKey{b524NumericCandidateTestKey(2, 0, 2), b524NumericCandidateTestKey(2, 3, 2), b524NumericCandidateTestKey(3, 0, 128)} {
		b524NumericCandidateMiss(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateKey)
	}
}

func TestAb56EveryKeyAxis(b524NumericCandidateT *b524NumericCandidateTesting.T) {
	b524NumericCandidateView := b524NumericCandidateViewForTest(b524NumericCandidateT)
	b524NumericCandidateControl := b524NumericCandidateTestKey(0, 0, 3)
	for b524NumericCandidateIndex := 0; b524NumericCandidateIndex < b524NumericCandidateReflect.TypeOf(b524NumericCandidateControl).NumField(); b524NumericCandidateIndex++ {
		b524NumericCandidateName := b524NumericCandidateReflect.TypeOf(b524NumericCandidateControl).Field(b524NumericCandidateIndex).Name
		b524NumericCandidateT.Run("ab56_"+b524NumericCandidateName, func(b524NumericCandidateT *b524NumericCandidateTesting.T) {
			for _, b524NumericCandidateMode := range []string{"ab56_absent", "ab56_zero", "ab56_other"} {
				b524NumericCandidateChanged := b524NumericCandidateControl
				b524NumericCandidateField := b524NumericCandidateReflect.ValueOf(&b524NumericCandidateChanged).Elem().Field(b524NumericCandidateIndex)
				if b524NumericCandidateField.Kind() == b524NumericCandidateReflect.String {
					if b524NumericCandidateMode == "ab56_other" {
						b524NumericCandidateField.SetString("ab56_other")
					} else {
						b524NumericCandidateField.SetString("")
					}
				} else if b524NumericCandidateMode == "ab56_absent" {
					b524NumericCandidateField.SetZero()
				} else {
					b524NumericCandidatePointerValue := b524NumericCandidateReflect.New(b524NumericCandidateField.Type().Elem())
					if b524NumericCandidateMode == "ab56_other" {
						b524NumericCandidatePointerValue.Elem().SetUint(255)
					}
					b524NumericCandidateField.Set(b524NumericCandidatePointerValue)
				}
				if b524NumericCandidateReflect.DeepEqual(b524NumericCandidateChanged, b524NumericCandidateControl) {
					continue
				}
				b524NumericCandidateMiss(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateChanged)
				b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateControl)
			}
		})
	}
	b524NumericCandidateOP06 := b524NumericCandidateControl
	b524NumericCandidateOP06.B524NumericCandidateOperation, b524NumericCandidateOP06.B524NumericCandidateOpcode = b524NumericCandidatePointer(uint8(6)), b524NumericCandidatePointer(uint8(6))
	b524NumericCandidateMiss(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateOP06)
	b524NumericCandidateMiss(b524NumericCandidateT, B524NumericCandidateView{}, b524NumericCandidateControl)
}

func TestAb56ContextAxes(b524NumericCandidateT *b524NumericCandidateTesting.T) {
	b524NumericCandidateControl := b524NumericCandidateTestContext()
	b524NumericCandidateKey := b524NumericCandidateTestKey(3, 0, 2)
	b524NumericCandidateOriginal := b524NumericCandidateViewForTest(b524NumericCandidateT)
	for b524NumericCandidateIndex := 0; b524NumericCandidateIndex < b524NumericCandidateReflect.TypeOf(b524NumericCandidateControl).NumField(); b524NumericCandidateIndex++ {
		b524NumericCandidateName := b524NumericCandidateReflect.TypeOf(b524NumericCandidateControl).Field(b524NumericCandidateIndex).Name
		b524NumericCandidateT.Run("ab56_"+b524NumericCandidateName, func(b524NumericCandidateT *b524NumericCandidateTesting.T) {
			b524NumericCandidateMissing := b524NumericCandidateControl
			b524NumericCandidateReflect.ValueOf(&b524NumericCandidateMissing).Elem().Field(b524NumericCandidateIndex).SetString("")
			if _, b524NumericCandidateError := B524NumericCandidateNewView(b524NumericCandidateMissing); b524NumericCandidateError == nil {
				b524NumericCandidateT.Fatal("ab56 incomplete context accepted")
			}
			b524NumericCandidateChanged := b524NumericCandidateControl
			b524NumericCandidateReflect.ValueOf(&b524NumericCandidateChanged).Elem().Field(b524NumericCandidateIndex).SetString("ab56_other")
			b524NumericCandidateView, b524NumericCandidateError := B524NumericCandidateNewView(b524NumericCandidateChanged)
			if b524NumericCandidateError != nil {
				b524NumericCandidateT.Fatal(b524NumericCandidateError)
			}
			b524NumericCandidateChangedKey := b524NumericCandidateKey
			b524NumericCandidateReflect.ValueOf(&b524NumericCandidateChangedKey).Elem().FieldByName(b524NumericCandidateName).SetString("ab56_other")
			b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateChangedKey)
			b524NumericCandidateMiss(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateKey)
			b524NumericCandidateMiss(b524NumericCandidateT, b524NumericCandidateOriginal, b524NumericCandidateChangedKey)
			b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateOriginal, b524NumericCandidateKey)
			b524NumericCandidateReflect.ValueOf(&b524NumericCandidateChanged).Elem().Field(b524NumericCandidateIndex).SetString("ab56_after")
			b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateChangedKey)
		})
	}
}

func TestAb56ArithmeticOnly(b524NumericCandidateT *b524NumericCandidateTesting.T) {
	b524NumericCandidateRow := b524NumericCandidateDecode(b524NumericCandidateT, b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateViewForTest(b524NumericCandidateT), b524NumericCandidateTestKey(3, 0, 2)))
	var b524NumericCandidateMinimum, b524NumericCandidateMaximum, b524NumericCandidateStep float64
	for b524NumericCandidateName, b524NumericCandidateDestination := range map[string]*float64{"minimum": &b524NumericCandidateMinimum, "maximum": &b524NumericCandidateMaximum, "step": &b524NumericCandidateStep} {
		if b524NumericCandidateError := b524NumericCandidateJSON.Unmarshal(b524NumericCandidateRow[b524NumericCandidateName], b524NumericCandidateDestination); b524NumericCandidateError != nil {
			b524NumericCandidateT.Fatal(b524NumericCandidateError)
		}
	}
	for _, b524NumericCandidateCase := range []struct {
		b524NumericCandidateValue float64
		b524NumericCandidateWant  bool
	}{{15, true}, {15.5, true}, {30, true}, {15.25, false}, {14.5, false}, {30.5, false}} {
		b524NumericCandidateOffset := (b524NumericCandidateCase.b524NumericCandidateValue - b524NumericCandidateMinimum) / b524NumericCandidateStep
		b524NumericCandidateGot := b524NumericCandidateCase.b524NumericCandidateValue >= b524NumericCandidateMinimum && b524NumericCandidateCase.b524NumericCandidateValue <= b524NumericCandidateMaximum && b524NumericCandidateOffset == float64(int(b524NumericCandidateOffset))
		if b524NumericCandidateGot != b524NumericCandidateCase.b524NumericCandidateWant {
			b524NumericCandidateT.Fatalf("ab56 arithmetic %v", b524NumericCandidateCase.b524NumericCandidateValue)
		}
	}
	if string(b524NumericCandidateRow["admitted"]) != "false" || string(b524NumericCandidateRow["presence"]) != `"Unknown"` || string(b524NumericCandidateRow["freshness"]) != `"Unknown"` {
		b524NumericCandidateT.Fatal("ab56 metadata promoted")
	}
}

func TestAb56StorageOwnership(b524NumericCandidateT *b524NumericCandidateTesting.T) {
	b524NumericCandidateOriginal, b524NumericCandidateError := B524NumericCandidateEntries()
	if b524NumericCandidateError != nil {
		b524NumericCandidateT.Fatal(b524NumericCandidateError)
	}
	b524NumericCandidateSource := append([]byte(nil), b524NumericCandidateOriginal...)
	b524NumericCandidateView, b524NumericCandidateError := b524NumericCandidateNewView(b524NumericCandidateTestContext(), b524NumericCandidateSource)
	if b524NumericCandidateError != nil {
		b524NumericCandidateT.Fatal(b524NumericCandidateError)
	}
	b524NumericCandidateKey := b524NumericCandidateTestKey(3, 0, 2)
	b524NumericCandidateBefore := b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateKey)
	for b524NumericCandidateIndex := range b524NumericCandidateSource {
		b524NumericCandidateSource[b524NumericCandidateIndex] = 'x'
	}
	b524NumericCandidateReturned := b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateKey)
	for b524NumericCandidateIndex := range b524NumericCandidateReturned {
		b524NumericCandidateReturned[b524NumericCandidateIndex] = 'x'
	}
	b524NumericCandidateEqual(b524NumericCandidateT, b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateKey), b524NumericCandidateBefore)
	for b524NumericCandidateIndex := range b524NumericCandidateOriginal {
		b524NumericCandidateOriginal[b524NumericCandidateIndex] = 'x'
	}
	b524NumericCandidateAfter, b524NumericCandidateError := B524NumericCandidateEntries()
	if b524NumericCandidateError != nil || b524NumericCandidateBytes.Equal(b524NumericCandidateOriginal, b524NumericCandidateAfter) {
		b524NumericCandidateT.Fatal("ab56 entries storage shared")
	}
	*b524NumericCandidateKey.B524NumericCandidateRegister = 99
	b524NumericCandidateEqual(b524NumericCandidateT, b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateTestKey(3, 0, 2)), b524NumericCandidateBefore)
}

func TestAb56OptionalAbsentZero(b524NumericCandidateT *b524NumericCandidateTesting.T) {
	b524NumericCandidateControl := b524NumericCandidateTestKey(0, 0, 3)
	b524NumericCandidateBaseRow := b524NumericCandidateDecode(b524NumericCandidateT, b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateViewForTest(b524NumericCandidateT), b524NumericCandidateControl))
	for _, b524NumericCandidateName := range []string{"B524NumericCandidateSecondary", "B524NumericCandidateOperation", "B524NumericCandidateOpcode", "B524NumericCandidateGroup", "B524NumericCandidateInstance", "B524NumericCandidateRegister"} {
		b524NumericCandidateT.Run("ab56_"+b524NumericCandidateName, func(b524NumericCandidateT *b524NumericCandidateTesting.T) {
			for _, b524NumericCandidateAbsent := range []bool{true, false} {
				b524NumericCandidateStored := b524NumericCandidateControl
				b524NumericCandidateStored.B524NumericCandidateFamily, b524NumericCandidateStored.B524NumericCandidateSoftware, b524NumericCandidateStored.B524NumericCandidateHardware, b524NumericCandidateStored.B524NumericCandidateAPIVersion, b524NumericCandidateStored.B524NumericCandidateAPIRevision = "", "", "", "", ""
				b524NumericCandidateField := b524NumericCandidateReflect.ValueOf(&b524NumericCandidateStored).Elem().FieldByName(b524NumericCandidateName)
				b524NumericCandidateOther := b524NumericCandidateControl
				b524NumericCandidateOtherField := b524NumericCandidateReflect.ValueOf(&b524NumericCandidateOther).Elem().FieldByName(b524NumericCandidateName)
				if b524NumericCandidateAbsent {
					b524NumericCandidateField.SetZero()
					b524NumericCandidateOtherField.Set(b524NumericCandidateReflect.New(b524NumericCandidateOtherField.Type().Elem()))
				} else {
					b524NumericCandidateField.Set(b524NumericCandidateReflect.New(b524NumericCandidateField.Type().Elem()))
					b524NumericCandidateOtherField.SetZero()
				}
				b524NumericCandidateRow := b524NumericCandidateDecode(b524NumericCandidateT, b524NumericCandidateEncode(b524NumericCandidateT, b524NumericCandidateBaseRow))
				b524NumericCandidateRow["key"] = b524NumericCandidateEncode(b524NumericCandidateT, b524NumericCandidateStored)
				b524NumericCandidateView, b524NumericCandidateError := b524NumericCandidateNewView(b524NumericCandidateTestContext(), b524NumericCandidateEncode(b524NumericCandidateT, []any{b524NumericCandidateRow}))
				if b524NumericCandidateError != nil {
					b524NumericCandidateT.Fatal(b524NumericCandidateError)
				}
				b524NumericCandidateMatch := b524NumericCandidateControl
				b524NumericCandidateReflect.ValueOf(&b524NumericCandidateMatch).Elem().FieldByName(b524NumericCandidateName).Set(b524NumericCandidateField)
				b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateMatch)
				b524NumericCandidateMiss(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateOther)
			}
		})
	}
}

func TestAb56GradesWidthsAndRaw(b524NumericCandidateT *b524NumericCandidateTesting.T) {
	b524NumericCandidateKey := b524NumericCandidateTestKey(3, 0, 2)
	for _, b524NumericCandidateGradeValue := range []string{"Proven", "Hypothesis", "Unknown"} {
		for _, b524NumericCandidateWidth := range []any{nil, 0, 1, 2} {
			b524NumericCandidateRow := b524NumericCandidateDecode(b524NumericCandidateT, b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateViewForTest(b524NumericCandidateT), b524NumericCandidateKey))
			b524NumericCandidateRow["grade"] = b524NumericCandidateEncode(b524NumericCandidateT, b524NumericCandidateGradeValue)
			b524NumericCandidateRow["width_bytes"] = b524NumericCandidateEncode(b524NumericCandidateT, b524NumericCandidateWidth)
			b524NumericCandidateRow["ab56_extension"] = b524NumericCandidateJSON.RawMessage(`{"ab56_raw":[0,255],"ab56_unresolved":"Unknown"}`)
			b524NumericCandidateRaw := b524NumericCandidateEncode(b524NumericCandidateT, []any{b524NumericCandidateRow})
			b524NumericCandidateView, b524NumericCandidateError := b524NumericCandidateNewView(b524NumericCandidateTestContext(), b524NumericCandidateRaw)
			if b524NumericCandidateError != nil {
				b524NumericCandidateT.Fatal(b524NumericCandidateError)
			}
			b524NumericCandidateGot := b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateKey)
			b524NumericCandidateEqual(b524NumericCandidateT, b524NumericCandidateGot, b524NumericCandidateEncode(b524NumericCandidateT, b524NumericCandidateRow))
			if !b524NumericCandidateBytes.Contains(b524NumericCandidateRaw, b524NumericCandidateGot) {
				b524NumericCandidateT.Fatal("ab56 raw row rewritten")
			}
		}
	}
}

func TestAb56InvalidAndDuplicateData(b524NumericCandidateT *b524NumericCandidateTesting.T) {
	b524NumericCandidateKey := b524NumericCandidateTestKey(3, 0, 2)
	b524NumericCandidateRow := b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateViewForTest(b524NumericCandidateT), b524NumericCandidateKey)
	b524NumericCandidateInvalid := [][]byte{nil, []byte(`null`), []byte(`{}`), []byte(`[null]`), []byte(`[5]`), []byte(`[{}]`), []byte(`[`), b524NumericCandidateEncode(b524NumericCandidateT, []b524NumericCandidateJSON.RawMessage{b524NumericCandidateRow, b524NumericCandidateRow})}
	for _, b524NumericCandidateName := range []string{"codec", "width_bytes", "width_grade", "minimum", "maximum", "step", "grade", "physical_applicability", "qualification", "admitted", "presence", "freshness", "unresolved", "raw_form", "raw_fallback", "key"} {
		b524NumericCandidateObject := b524NumericCandidateDecode(b524NumericCandidateT, b524NumericCandidateRow)
		b524NumericCandidateObject[b524NumericCandidateName] = b524NumericCandidateJSON.RawMessage(`{"ab56_bad":true}`)
		b524NumericCandidateInvalid = append(b524NumericCandidateInvalid, b524NumericCandidateEncode(b524NumericCandidateT, []any{b524NumericCandidateObject}))
		delete(b524NumericCandidateObject, b524NumericCandidateName)
		b524NumericCandidateInvalid = append(b524NumericCandidateInvalid, b524NumericCandidateEncode(b524NumericCandidateT, []any{b524NumericCandidateObject}))
	}
	for b524NumericCandidateName, b524NumericCandidateValue := range map[string]string{"admitted": "true", "presence": `"present"`, "freshness": `"fresh"`, "writable": "true", "step": "0"} {
		b524NumericCandidateObject := b524NumericCandidateDecode(b524NumericCandidateT, b524NumericCandidateRow)
		b524NumericCandidateObject[b524NumericCandidateName] = b524NumericCandidateJSON.RawMessage(b524NumericCandidateValue)
		b524NumericCandidateInvalid = append(b524NumericCandidateInvalid, b524NumericCandidateEncode(b524NumericCandidateT, []any{b524NumericCandidateObject}))
	}
	b524NumericCandidateInvalid = append(b524NumericCandidateInvalid,
		b524NumericCandidateBytes.Replace(b524NumericCandidateEncode(b524NumericCandidateT, []b524NumericCandidateJSON.RawMessage{b524NumericCandidateRow}), []byte(`"grade":"Hypothesis"`), []byte(`"grade":"Unknown","grade":"Hypothesis"`), 1),
		b524NumericCandidateBytes.Replace(b524NumericCandidateEncode(b524NumericCandidateT, []b524NumericCandidateJSON.RawMessage{b524NumericCandidateRow}), []byte(`"primary":181`), []byte(`"primary":181,"primary":181`), 1),
	)
	for _, b524NumericCandidateName := range []string{"family", "software", "hardware", "api_version", "api_revision", "ab56_unknown_axis"} {
		b524NumericCandidateObject := b524NumericCandidateDecode(b524NumericCandidateT, b524NumericCandidateRow)
		b524NumericCandidateObjectKey := b524NumericCandidateDecode(b524NumericCandidateT, b524NumericCandidateObject["key"])
		b524NumericCandidateObjectKey[b524NumericCandidateName] = b524NumericCandidateJSON.RawMessage(`"ab56_forbidden"`)
		b524NumericCandidateObject["key"] = b524NumericCandidateEncode(b524NumericCandidateT, b524NumericCandidateObjectKey)
		b524NumericCandidateInvalid = append(b524NumericCandidateInvalid, b524NumericCandidateEncode(b524NumericCandidateT, []any{b524NumericCandidateObject}))
	}
	b524NumericCandidateExistingView := b524NumericCandidateViewForTest(b524NumericCandidateT)
	b524NumericCandidateSaved := b524NumericCandidateData
	defer func() { b524NumericCandidateData = b524NumericCandidateSaved }()
	for b524NumericCandidateIndex, b524NumericCandidateSource := range b524NumericCandidateInvalid {
		if _, b524NumericCandidateError := b524NumericCandidateNewView(b524NumericCandidateTestContext(), b524NumericCandidateSource); b524NumericCandidateError == nil {
			b524NumericCandidateT.Fatalf("ab56 malformed case %d accepted", b524NumericCandidateIndex)
		}
		b524NumericCandidateEqual(b524NumericCandidateT, b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateExistingView, b524NumericCandidateKey), b524NumericCandidateRow)
		b524NumericCandidateData = string(b524NumericCandidateSource)
		if b524NumericCandidateRaw, b524NumericCandidateError := B524NumericCandidateEntries(); b524NumericCandidateError == nil || b524NumericCandidateRaw != nil {
			b524NumericCandidateT.Fatalf("ab56 entries case %d accepted", b524NumericCandidateIndex)
		}
		if _, b524NumericCandidateError := B524NumericCandidateNewView(b524NumericCandidateTestContext()); b524NumericCandidateError == nil {
			b524NumericCandidateT.Fatalf("ab56 public constructor case %d accepted", b524NumericCandidateIndex)
		}
	}
}

func TestAb56OpaqueIdentityStrings(b524NumericCandidateT *b524NumericCandidateTesting.T) {
	b524NumericCandidateContext := b524NumericCandidateTestContext()
	b524NumericCandidateContext.B524NumericCandidateFamily = "ab56_\xff"
	b524NumericCandidateView, b524NumericCandidateError := B524NumericCandidateNewView(b524NumericCandidateContext)
	if b524NumericCandidateError != nil {
		b524NumericCandidateT.Fatal(b524NumericCandidateError)
	}
	b524NumericCandidateKey := b524NumericCandidateTestKey(3, 0, 2)
	b524NumericCandidateKey.B524NumericCandidateFamily = "ab56_\xff"
	b524NumericCandidateHit(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateKey)
	b524NumericCandidateKey.B524NumericCandidateFamily = "ab56_\ufffd"
	b524NumericCandidateMiss(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateKey)
	b524NumericCandidateKey.B524NumericCandidateFamily = "ab56_\xfe"
	b524NumericCandidateMiss(b524NumericCandidateT, b524NumericCandidateView, b524NumericCandidateKey)
}
