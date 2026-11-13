package vaillant

import (
	b512DHWStepperJSON "encoding/json"
	b512DHWStepperReflect "reflect"
	b512DHWStepperStrconv "strconv"
	b512DHWStepperSync "sync"
	b512DHWStepperTesting "testing"
)

func b512DHWStepperByte(b512DHWStepperValue uint8) *uint8 { return &b512DHWStepperValue }
func b512DHWStepperKey(b512DHWStepperOperation *uint8, b512DHWStepperDirection, b512DHWStepperField string) B512DHWStepperKey {
	return B512DHWStepperKey{B512DHWStepperFamily: "vaillant", B512DHWStepperVariant: "B512", B512DHWStepperPrimary: 181, B512DHWStepperSecondary: b512DHWStepperByte(18), B512DHWStepperOperation: b512DHWStepperOperation, B512DHWStepperDirection: b512DHWStepperDirection, B512DHWStepperField: b512DHWStepperField}
}
func b512DHWStepperRequire(b512DHWStepperT *b512DHWStepperTesting.T, b512DHWStepperOperation uint8, b512DHWStepperDirection, b512DHWStepperField string) B512DHWStepperEntry {
	b512DHWStepperT.Helper()
	b512DHWStepperEntry, b512DHWStepperFound, b512DHWStepperErr := B512DHWStepperLookup(b512DHWStepperKey(b512DHWStepperByte(b512DHWStepperOperation), b512DHWStepperDirection, b512DHWStepperField))
	if b512DHWStepperErr != nil || !b512DHWStepperFound {
		b512DHWStepperT.Fatalf("lookup %d %s %s: %v %v", b512DHWStepperOperation, b512DHWStepperDirection, b512DHWStepperField, b512DHWStepperFound, b512DHWStepperErr)
	}
	return b512DHWStepperEntry
}
func b512DHWStepperScalar(b512DHWStepperT *b512DHWStepperTesting.T, b512DHWStepperEntry B512DHWStepperEntry, b512DHWStepperRaw []byte) float64 {
	b512DHWStepperT.Helper()
	if b512DHWStepperEntry.B512DHWStepperWidthBits == nil {
		b512DHWStepperT.Fatal("missing width")
	}
	b512DHWStepperValue := int(b512DHWStepperRaw[0])
	if *b512DHWStepperEntry.B512DHWStepperWidthBits == 16 {
		if b512DHWStepperEntry.B512DHWStepperByteOrder == nil || *b512DHWStepperEntry.B512DHWStepperByteOrder != "little-endian" {
			b512DHWStepperT.Fatal("byte order")
		}
		b512DHWStepperValue += int(b512DHWStepperRaw[1]) << 8
		if b512DHWStepperEntry.B512DHWStepperSigned != nil && *b512DHWStepperEntry.B512DHWStepperSigned {
			b512DHWStepperValue = int(int16(b512DHWStepperValue))
		}
	} else if b512DHWStepperEntry.B512DHWStepperSigned != nil && *b512DHWStepperEntry.B512DHWStepperSigned {
		b512DHWStepperValue = int(int8(b512DHWStepperValue))
	}
	b512DHWStepperResult := float64(b512DHWStepperValue)
	if b512DHWStepperEntry.B512DHWStepperDivisor != nil {
		b512DHWStepperResult /= *b512DHWStepperEntry.B512DHWStepperDivisor
	}
	return b512DHWStepperResult
}
func TestAb97Scalars(b512DHWStepperT *b512DHWStepperTesting.T) {
	b512DHWStepperRaw := b512DHWStepperRequire(b512DHWStepperT, 3, "positive", "temperature_raw")
	if b512DHWStepperScalar(b512DHWStepperT, b512DHWStepperRaw, []byte{0x34, 0x12}) != 4660 || b512DHWStepperRaw.B512DHWStepperDivisor != nil || b512DHWStepperRaw.B512DHWStepperUnit != nil || *b512DHWStepperRaw.B512DHWStepperPosition != 4 {
		b512DHWStepperT.Fatal("OP03 raw must remain unscaled")
	}
	if b512DHWStepperRequire(b512DHWStepperT, 3, "request", "selector").B512DHWStepperValues["0"] != "valve-temperature-input" {
		b512DHWStepperT.Fatal("selector")
	}
	b512DHWStepperStatus := b512DHWStepperRequire(b512DHWStepperT, 3, "positive", "status")
	if !b512DHWStepperReflect.DeepEqual(b512DHWStepperStatus.B512DHWStepperValues, map[string]string{"0": "okay", "1": "disconnected", "2": "short-circuit"}) || *b512DHWStepperStatus.B512DHWStepperPosition != 3 {
		b512DHWStepperT.Fatal("status")
	}
	for _, b512DHWStepperCase := range []struct {
		b512DHWStepperOp                             uint8
		b512DHWStepperDirection, b512DHWStepperField string
		b512DHWStepperPosition                       uint8
		b512DHWStepperRaw                            []byte
		b512DHWStepperWant                           float64
		b512DHWStepperUnit                           string
	}{
		{6, "request", "flow", 5, []byte{1}, 0.1, "l/min"},
		{6, "request", "temperature_6", 6, []byte{1}, 0.5, "Celsius"},
		{6, "request", "temperature_7", 7, []byte{1}, 0.5, "Celsius"},
		{6, "request", "temperature_8", 8, []byte{1}, 0.5, "Celsius"},
		{6, "positive", "dhw_setpoint", 4, []byte{1}, 0.5, "Celsius"},
		{6, "positive", "return", 5, []byte{1}, 0.5, "Celsius"},
		{8, "request", "inlet", 3, []byte{0xf0, 0xff}, -1, "Celsius"},
		{8, "positive", "dhw_setpoint", 3, []byte{0xf0, 0xff}, -1, "Celsius"},
	} {
		b512DHWStepperEntry := b512DHWStepperRequire(b512DHWStepperT, b512DHWStepperCase.b512DHWStepperOp, b512DHWStepperCase.b512DHWStepperDirection, b512DHWStepperCase.b512DHWStepperField)
		if b512DHWStepperScalar(b512DHWStepperT, b512DHWStepperEntry, b512DHWStepperCase.b512DHWStepperRaw) != b512DHWStepperCase.b512DHWStepperWant || *b512DHWStepperEntry.B512DHWStepperPosition != b512DHWStepperCase.b512DHWStepperPosition || b512DHWStepperEntry.B512DHWStepperUnit == nil || *b512DHWStepperEntry.B512DHWStepperUnit != b512DHWStepperCase.b512DHWStepperUnit || *b512DHWStepperEntry.B512DHWStepperDefault != 0 {
			b512DHWStepperT.Fatalf("scalar %s", b512DHWStepperCase.b512DHWStepperField)
		}
	}
	b512DHWStepperAdjustment := b512DHWStepperRequire(b512DHWStepperT, 6, "positive", "adjustment")
	if b512DHWStepperScalar(b512DHWStepperT, b512DHWStepperAdjustment, []byte{255}) != -1 || b512DHWStepperScalar(b512DHWStepperT, b512DHWStepperAdjustment, []byte{2}) != 2 || *b512DHWStepperAdjustment.B512DHWStepperDefault != -1 || *b512DHWStepperAdjustment.B512DHWStepperPosition != 6 || !b512DHWStepperReflect.DeepEqual(b512DHWStepperAdjustment.B512DHWStepperValues, map[string]string{"-1": "-1", "0": "0", "1": "1", "2": "2"}) {
		b512DHWStepperT.Fatal("signed adjustment")
	}
	b512DHWStepperModulation := b512DHWStepperRequire(b512DHWStepperT, 6, "request", "modulation")
	if b512DHWStepperScalar(b512DHWStepperT, b512DHWStepperModulation, []byte{1}) != 0.5 || *b512DHWStepperModulation.B512DHWStepperPosition != 4 || b512DHWStepperModulation.B512DHWStepperUnit != nil || *b512DHWStepperModulation.B512DHWStepperFormalUnit != "Celsius" || b512DHWStepperModulation.B512DHWStepperUnknown[0] != "modulation unit conflicts with formal Celsius unit" {
		b512DHWStepperT.Fatal("unit conflict lost")
	}
	b512DHWStepperStatus = b512DHWStepperRequire(b512DHWStepperT, 8, "request", "status")
	if *b512DHWStepperStatus.B512DHWStepperPosition != 5 || *b512DHWStepperStatus.B512DHWStepperDefault != 0 {
		b512DHWStepperT.Fatal("OP08 status position/default")
	}
	for b512DHWStepperCode, b512DHWStepperWant := range map[uint8]string{0: "valid", 0x55: "short", 0xaa: "interruption-or-disconnected", 0x54: ""} {
		if b512DHWStepperStatus.B512DHWStepperValues[b512DHWStepperStrconv.Itoa(int(b512DHWStepperCode))] != b512DHWStepperWant {
			b512DHWStepperT.Fatal("OP08 status mapping")
		}
	}
}

func TestAb97Flags(b512DHWStepperT *b512DHWStepperTesting.T) {
	for _, b512DHWStepperCase := range []struct {
		b512DHWStepperOp        uint8
		b512DHWStepperDirection string
		b512DHWStepperPosition  uint8
		b512DHWStepperFields    []string
		b512DHWStepperNamed     int
	}{
		{6, "request", 3, []string{"task-active", "heat-demand", "loading-active", "loading-pump-fault", "condition_4", "condition_5", "condition_6", "condition_7"}, 4},
		{6, "request", 9, []string{"inlet-condition_0", "inlet-condition_1", "reserved_2", "reserved_3", "reserved_4", "reserved_5", "reserved_6", "reserved_7"}, 0},
		{6, "positive", 3, []string{"loading-disabled", "internal-test-active", "ready", "frost-protection-active", "legionella-demand", "reset", "reserved_6", "reserved_7"}, 6},
		{8, "request", 6, []string{"stepper-release", "reserved_1", "reserved_2", "reserved_3", "reserved_4", "reserved_5", "reserved_6", "reserved_7"}, 1},
		{8, "positive", 5, []string{"stepper-release", "reserved_1", "reserved_2", "reserved_3", "reserved_4", "reserved_5", "reserved_6", "reserved_7"}, 1},
		{9, "request", 3, []string{"stepper-release", "reserved_1", "reserved_2", "reserved_3", "reserved_4", "reserved_5", "reserved_6", "reserved_7"}, 1},
		{9, "positive", 3, []string{"stepper-release", "reserved_1", "reserved_2", "reserved_3", "reserved_4", "reserved_5", "reserved_6", "reserved_7"}, 1},
	} {
		for b512DHWStepperBit, b512DHWStepperField := range b512DHWStepperCase.b512DHWStepperFields {
			b512DHWStepperEntry := b512DHWStepperRequire(b512DHWStepperT, b512DHWStepperCase.b512DHWStepperOp, b512DHWStepperCase.b512DHWStepperDirection, b512DHWStepperField)
			if *b512DHWStepperEntry.B512DHWStepperPosition != b512DHWStepperCase.b512DHWStepperPosition || *b512DHWStepperEntry.B512DHWStepperWidthBits != 1 || b512DHWStepperEntry.B512DHWStepperBit == nil {
				b512DHWStepperT.Fatal("flag coordinate")
			}
			for b512DHWStepperToggle := 0; b512DHWStepperToggle < 8; b512DHWStepperToggle++ {
				if ((uint8(1)<<b512DHWStepperToggle)&(uint8(1)<<*b512DHWStepperEntry.B512DHWStepperBit) != 0) != (b512DHWStepperToggle == b512DHWStepperBit) {
					b512DHWStepperT.Fatalf("flag swap: %s", b512DHWStepperField)
				}
			}
			if b512DHWStepperBit < b512DHWStepperCase.b512DHWStepperNamed {
				if b512DHWStepperEntry.B512DHWStepperLabel == nil || *b512DHWStepperEntry.B512DHWStepperLabel != b512DHWStepperField {
					b512DHWStepperT.Fatal("label")
				}
			} else if b512DHWStepperEntry.B512DHWStepperLabel != nil {
				b512DHWStepperT.Fatal("invented condition label")
			}
			b512DHWStepperReserved := len(b512DHWStepperField) >= 8 && b512DHWStepperField[:8] == "reserved"
			if b512DHWStepperEntry.B512DHWStepperReserved != b512DHWStepperReserved {
				b512DHWStepperT.Fatal("reserved bit lost")
			}
			if b512DHWStepperReserved {
				if b512DHWStepperEntry.B512DHWStepperDefault == nil || *b512DHWStepperEntry.B512DHWStepperDefault != 0 {
					b512DHWStepperT.Fatal("reserved default")
				}
			} else if b512DHWStepperEntry.B512DHWStepperBooleanDefault == nil || *b512DHWStepperEntry.B512DHWStepperBooleanDefault {
				b512DHWStepperT.Fatal("Boolean default")
			}
			if b512DHWStepperField == "stepper-release" && !b512DHWStepperReflect.DeepEqual(b512DHWStepperEntry.B512DHWStepperUnknown, []string{"release does not prove movement", "physical adaptation", "support", "actual movement"}) {
				b512DHWStepperT.Fatal("release implies movement")
			}
		}
	}
}

func TestAb97ExactKeys(b512DHWStepperT *b512DHWStepperTesting.T) {
	b512DHWStepperBase := b512DHWStepperKey(b512DHWStepperByte(3), "positive", "temperature_raw")
	for b512DHWStepperName, b512DHWStepperChange := range map[string]func(*B512DHWStepperKey){
		"family":  func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperFamily = "other" },
		"variant": func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperVariant = "other" },
		"primary": func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperPrimary = 0 },
		"secondary": func(b512DHWStepperKey *B512DHWStepperKey) {
			b512DHWStepperKey.B512DHWStepperSecondary = b512DHWStepperByte(0)
		},
		"secondary_absent": func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperSecondary = nil },
		"operation": func(b512DHWStepperKey *B512DHWStepperKey) {
			b512DHWStepperKey.B512DHWStepperOperation = b512DHWStepperByte(0)
		},
		"operation_absent": func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperOperation = nil },
		"direction":        func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperDirection = "request" },
		"field":            func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperField = "other" },
		"group": func(b512DHWStepperKey *B512DHWStepperKey) {
			b512DHWStepperKey.B512DHWStepperGroup = b512DHWStepperByte(0)
		},
		"requested_instance": func(b512DHWStepperKey *B512DHWStepperKey) {
			b512DHWStepperKey.B512DHWStepperRequestedInstance = b512DHWStepperByte(0)
		},
		"register": func(b512DHWStepperKey *B512DHWStepperKey) {
			b512DHWStepperValue := uint16(0)
			b512DHWStepperKey.B512DHWStepperRegister = &b512DHWStepperValue
		},
		"command_role": func(b512DHWStepperKey *B512DHWStepperKey) {
			b512DHWStepperValue := "read"
			b512DHWStepperKey.B512DHWStepperCommandRole = &b512DHWStepperValue
		},
	} {
		b512DHWStepperT.Run("ab97_"+b512DHWStepperName, func(b512DHWStepperT *b512DHWStepperTesting.T) {
			b512DHWStepperChanged := b512DHWStepperBase
			b512DHWStepperChange(&b512DHWStepperChanged)
			b512DHWStepperEntry, b512DHWStepperFound, b512DHWStepperErr := B512DHWStepperLookup(b512DHWStepperChanged)
			if b512DHWStepperErr != nil || b512DHWStepperFound || !b512DHWStepperReflect.DeepEqual(b512DHWStepperEntry, B512DHWStepperEntry{}) {
				b512DHWStepperT.Fatalf("non-exact lookup: %v %v", b512DHWStepperFound, b512DHWStepperErr)
			}
		})
	}
	for _, b512DHWStepperDirection := range []string{"request", "positive"} {
		b512DHWStepperEntry := b512DHWStepperRequire(b512DHWStepperT, 10, b512DHWStepperDirection, "declaration")
		if b512DHWStepperEntry.B512DHWStepperPosition != nil || b512DHWStepperEntry.B512DHWStepperWidthBits != nil || b512DHWStepperEntry.B512DHWStepperDivisor != nil {
			b512DHWStepperT.Fatal("invented OP0A body")
		}
		b512DHWStepperKey := b512DHWStepperKey(nil, b512DHWStepperDirection, "declaration")
		b512DHWStepperKey.B512DHWStepperVariant = "B51214"
		b512DHWStepperKey.B512DHWStepperSecondary = b512DHWStepperByte(245)
		b512DHWStepperEntry, b512DHWStepperFound, b512DHWStepperErr := B512DHWStepperLookup(b512DHWStepperKey)
		if b512DHWStepperErr != nil || !b512DHWStepperFound || b512DHWStepperEntry.B512DHWStepperKey.B512DHWStepperOperation != nil || b512DHWStepperEntry.B512DHWStepperWidthBits != nil {
			b512DHWStepperT.Fatal("partial B51214 declaration")
		}
		for _, b512DHWStepperOp := range []uint8{0, 10, 0x14} {
			b512DHWStepperKey.B512DHWStepperOperation = b512DHWStepperByte(b512DHWStepperOp)
			if _, b512DHWStepperFound, b512DHWStepperErr := B512DHWStepperLookup(b512DHWStepperKey); b512DHWStepperFound || b512DHWStepperErr != nil {
				b512DHWStepperT.Fatal("name encoded as operation")
			}
		}
	}
}

func b512DHWStepperMutate(b512DHWStepperValue b512DHWStepperReflect.Value) {
	if !b512DHWStepperValue.IsValid() {
		return
	}
	switch b512DHWStepperValue.Kind() {
	case b512DHWStepperReflect.Pointer:
		if !b512DHWStepperValue.IsNil() {
			b512DHWStepperMutate(b512DHWStepperValue.Elem())
		}
	case b512DHWStepperReflect.Struct:
		for b512DHWStepperIndex := 0; b512DHWStepperIndex < b512DHWStepperValue.NumField(); b512DHWStepperIndex++ {
			b512DHWStepperMutate(b512DHWStepperValue.Field(b512DHWStepperIndex))
		}
	case b512DHWStepperReflect.Slice:
		for b512DHWStepperIndex := 0; b512DHWStepperIndex < b512DHWStepperValue.Len(); b512DHWStepperIndex++ {
			b512DHWStepperMutate(b512DHWStepperValue.Index(b512DHWStepperIndex))
		}
	case b512DHWStepperReflect.Map:
		if !b512DHWStepperValue.IsNil() {
			b512DHWStepperValue.SetMapIndex(b512DHWStepperReflect.ValueOf("ab97_changed"), b512DHWStepperReflect.ValueOf("ab97_changed"))
		}
	case b512DHWStepperReflect.String:
		b512DHWStepperValue.SetString("ab97_changed")
	case b512DHWStepperReflect.Uint8, b512DHWStepperReflect.Uint16:
		b512DHWStepperValue.SetUint(0)
	case b512DHWStepperReflect.Int:
		b512DHWStepperValue.SetInt(99)
	case b512DHWStepperReflect.Float64:
		b512DHWStepperValue.SetFloat(99)
	case b512DHWStepperReflect.Bool:
		b512DHWStepperValue.SetBool(!b512DHWStepperValue.Bool())
	}
}
func TestAb97IndependentStorage(b512DHWStepperT *b512DHWStepperTesting.T) {
	b512DHWStepperBefore, b512DHWStepperErr := B512DHWStepperEntries()
	if b512DHWStepperErr != nil {
		b512DHWStepperT.Fatal(b512DHWStepperErr)
	}
	if len(b512DHWStepperBefore) != 74 {
		b512DHWStepperT.Fatalf("finite entries: %d", len(b512DHWStepperBefore))
	}
	b512DHWStepperChanged, b512DHWStepperErr := B512DHWStepperEntries()
	if b512DHWStepperErr != nil {
		b512DHWStepperT.Fatal(b512DHWStepperErr)
	}
	b512DHWStepperMutate(b512DHWStepperReflect.ValueOf(&b512DHWStepperChanged).Elem())
	for _, b512DHWStepperExpected := range b512DHWStepperBefore {
		if b512DHWStepperExpected.B512DHWStepperPosition != nil && (b512DHWStepperExpected.B512DHWStepperCoordinateBasis == nil || *b512DHWStepperExpected.B512DHWStepperCoordinateBasis != "logical; includes three-position B5/12/operation prefix") {
			b512DHWStepperT.Fatal("logical coordinate basis lost")
		}
		if b512DHWStepperExpected.B512DHWStepperPosition == nil && b512DHWStepperExpected.B512DHWStepperCoordinateBasis != nil {
			b512DHWStepperT.Fatal("invented coordinate basis")
		}
		if b512DHWStepperExpected.B512DHWStepperGrade != "Hypothesis" || b512DHWStepperExpected.B512DHWStepperQualification != "Unknown" || b512DHWStepperExpected.B512DHWStepperRawPolicy != "preserve native bytes, reserved bits and independent last-known-good fields" {
			b512DHWStepperT.Fatal("grade or raw boundary")
		}
		b512DHWStepperFound, b512DHWStepperOK, b512DHWStepperErr := B512DHWStepperLookup(b512DHWStepperExpected.B512DHWStepperKey)
		if b512DHWStepperErr != nil || !b512DHWStepperOK || !b512DHWStepperReflect.DeepEqual(b512DHWStepperFound, b512DHWStepperExpected) {
			b512DHWStepperT.Fatal("accessor aliases storage")
		}
		b512DHWStepperMutate(b512DHWStepperReflect.ValueOf(&b512DHWStepperFound).Elem())
		b512DHWStepperAgain, b512DHWStepperOK, b512DHWStepperErr := B512DHWStepperLookup(b512DHWStepperExpected.B512DHWStepperKey)
		if b512DHWStepperErr != nil || !b512DHWStepperOK || !b512DHWStepperReflect.DeepEqual(b512DHWStepperAgain, b512DHWStepperExpected) {
			b512DHWStepperT.Fatal("lookup aliases storage")
		}
	}
	b512DHWStepperAfter, b512DHWStepperErr := B512DHWStepperEntries()
	if b512DHWStepperErr != nil || !b512DHWStepperReflect.DeepEqual(b512DHWStepperBefore, b512DHWStepperAfter) {
		b512DHWStepperT.Fatal("nested mutation escaped")
	}
}
func TestAb97MalformedData(b512DHWStepperT *b512DHWStepperTesting.T) {
	b512DHWStepperOriginal := b512DHWStepperData
	defer func() {
		b512DHWStepperData = b512DHWStepperOriginal
		b512DHWStepperCatalog = b512DHWStepperSync.OnceValues(b512DHWStepperLoadEntries)
	}()
	b512DHWStepperRows, b512DHWStepperErr := B512DHWStepperEntries()
	if b512DHWStepperErr != nil {
		b512DHWStepperT.Fatal(b512DHWStepperErr)
	}
	b512DHWStepperDuplicate, b512DHWStepperErr := b512DHWStepperJSON.Marshal([]B512DHWStepperEntry{b512DHWStepperRows[0], b512DHWStepperRows[0]})
	if b512DHWStepperErr != nil {
		b512DHWStepperT.Fatal(b512DHWStepperErr)
	}
	b512DHWStepperCases := [][]byte{[]byte("{"), []byte("null"), []byte("[]"), []byte("[{}]"), append(append([]byte{}, b512DHWStepperOriginal...), []byte(" null")...), b512DHWStepperDuplicate}
	for _, b512DHWStepperChange := range []func(*B512DHWStepperEntry){
		func(b512DHWStepperEntry *B512DHWStepperEntry) { b512DHWStepperEntry.B512DHWStepperGrade = "Proven" },
		func(b512DHWStepperEntry *B512DHWStepperEntry) {
			b512DHWStepperEntry.B512DHWStepperQualification = "qualified"
		},
		func(b512DHWStepperEntry *B512DHWStepperEntry) {
			b512DHWStepperValue := float64(0)
			b512DHWStepperEntry.B512DHWStepperDivisor = &b512DHWStepperValue
		},
		func(b512DHWStepperEntry *B512DHWStepperEntry) {
			b512DHWStepperEntry.B512DHWStepperBit = b512DHWStepperByte(8)
		},
	} {
		b512DHWStepperEntry := b512DHWStepperRows[0]
		b512DHWStepperChange(&b512DHWStepperEntry)
		b512DHWStepperBytes, b512DHWStepperErr := b512DHWStepperJSON.Marshal([]B512DHWStepperEntry{b512DHWStepperEntry})
		if b512DHWStepperErr != nil {
			b512DHWStepperT.Fatal(b512DHWStepperErr)
		}
		b512DHWStepperCases = append(b512DHWStepperCases, b512DHWStepperBytes)
	}
	for b512DHWStepperIndex, b512DHWStepperInvalid := range b512DHWStepperCases {
		b512DHWStepperData = b512DHWStepperInvalid
		b512DHWStepperCatalog = b512DHWStepperSync.OnceValues(b512DHWStepperLoadEntries)
		if b512DHWStepperEntries, b512DHWStepperErr := B512DHWStepperEntries(); b512DHWStepperErr == nil || b512DHWStepperEntries != nil {
			b512DHWStepperT.Fatalf("invalid catalog %d accepted", b512DHWStepperIndex)
		}
		if b512DHWStepperEntry, b512DHWStepperFound, b512DHWStepperErr := B512DHWStepperLookup(b512DHWStepperRows[0].B512DHWStepperKey); b512DHWStepperErr == nil || b512DHWStepperFound || !b512DHWStepperReflect.DeepEqual(b512DHWStepperEntry, B512DHWStepperEntry{}) {
			b512DHWStepperT.Fatal("invalid lookup succeeded")
		}
	}
}

func TestAb97B524Identity(b512DHWStepperT *b512DHWStepperTesting.T) {
	b512DHWStepperOriginal := b512DHWStepperData
	defer func() {
		b512DHWStepperData = b512DHWStepperOriginal
		b512DHWStepperCatalog = b512DHWStepperSync.OnceValues(b512DHWStepperLoadEntries)
	}()
	b512DHWStepperRows, b512DHWStepperErr := B512DHWStepperEntries()
	if b512DHWStepperErr != nil {
		b512DHWStepperT.Fatal(b512DHWStepperErr)
	}
	b512DHWStepperRegister := uint16(1)
	b512DHWStepperRole := "synthetic-read"
	b512DHWStepperBase := B512DHWStepperKey{B512DHWStepperFamily: "synthetic-family", B512DHWStepperVariant: "synthetic-B524", B512DHWStepperPrimary: 181, B512DHWStepperSecondary: b512DHWStepperByte(36), B512DHWStepperOperation: b512DHWStepperByte(1), B512DHWStepperDirection: "request", B512DHWStepperField: "synthetic-field", B512DHWStepperGroup: b512DHWStepperByte(1), B512DHWStepperRequestedInstance: b512DHWStepperByte(1), B512DHWStepperRegister: &b512DHWStepperRegister, B512DHWStepperCommandRole: &b512DHWStepperRole}
	b512DHWStepperRow := b512DHWStepperRows[0]
	b512DHWStepperRow.B512DHWStepperKey = b512DHWStepperBase
	b512DHWStepperData, b512DHWStepperErr = b512DHWStepperJSON.Marshal([]B512DHWStepperEntry{b512DHWStepperRow})
	b512DHWStepperCatalog = b512DHWStepperSync.OnceValues(b512DHWStepperLoadEntries)
	if b512DHWStepperErr != nil {
		b512DHWStepperT.Fatal(b512DHWStepperErr)
	}
	if _, b512DHWStepperFound, b512DHWStepperErr := B512DHWStepperLookup(b512DHWStepperBase); b512DHWStepperErr != nil || !b512DHWStepperFound {
		b512DHWStepperT.Fatal("synthetic base missing")
	}
	for b512DHWStepperName, b512DHWStepperChange := range map[string]func(*B512DHWStepperKey){
		"family":  func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperFamily = "other" },
		"variant": func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperVariant = "other" },
		"primary": func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperPrimary = 0 },
		"secondary": func(b512DHWStepperKey *B512DHWStepperKey) {
			b512DHWStepperKey.B512DHWStepperSecondary = b512DHWStepperByte(0)
		},
		"operation": func(b512DHWStepperKey *B512DHWStepperKey) {
			b512DHWStepperKey.B512DHWStepperOperation = b512DHWStepperByte(0)
		},
		"direction": func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperDirection = "positive" },
		"field":     func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperField = "other" },
		"group": func(b512DHWStepperKey *B512DHWStepperKey) {
			b512DHWStepperKey.B512DHWStepperGroup = b512DHWStepperByte(0)
		},
		"group_absent": func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperGroup = nil },
		"requested_instance": func(b512DHWStepperKey *B512DHWStepperKey) {
			b512DHWStepperKey.B512DHWStepperRequestedInstance = b512DHWStepperByte(0)
		},
		"requested_instance_absent": func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperRequestedInstance = nil },
		"register": func(b512DHWStepperKey *B512DHWStepperKey) {
			b512DHWStepperValue := uint16(0)
			b512DHWStepperKey.B512DHWStepperRegister = &b512DHWStepperValue
		},
		"register_absent": func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperRegister = nil },
		"command_role": func(b512DHWStepperKey *B512DHWStepperKey) {
			b512DHWStepperValue := "synthetic-write"
			b512DHWStepperKey.B512DHWStepperCommandRole = &b512DHWStepperValue
		},
		"command_role_absent": func(b512DHWStepperKey *B512DHWStepperKey) { b512DHWStepperKey.B512DHWStepperCommandRole = nil },
	} {
		b512DHWStepperT.Run("ab97_"+b512DHWStepperName, func(b512DHWStepperT *b512DHWStepperTesting.T) {
			b512DHWStepperChanged := b512DHWStepperBase
			b512DHWStepperChange(&b512DHWStepperChanged)
			if _, b512DHWStepperFound, b512DHWStepperErr := B512DHWStepperLookup(b512DHWStepperChanged); b512DHWStepperErr != nil || b512DHWStepperFound {
				b512DHWStepperT.Fatal("synthetic B524 identity conflated")
			}
		})
	}
	b512DHWStepperChanged, b512DHWStepperFound, b512DHWStepperErr := B512DHWStepperLookup(b512DHWStepperBase)
	if b512DHWStepperErr != nil || !b512DHWStepperFound {
		b512DHWStepperT.Fatal("synthetic lookup")
	}
	*b512DHWStepperChanged.B512DHWStepperKey.B512DHWStepperRegister = 0
	*b512DHWStepperChanged.B512DHWStepperKey.B512DHWStepperGroup = 0
	*b512DHWStepperChanged.B512DHWStepperKey.B512DHWStepperRequestedInstance = 0
	*b512DHWStepperChanged.B512DHWStepperKey.B512DHWStepperCommandRole = "changed"
	b512DHWStepperAgain, b512DHWStepperFound, b512DHWStepperErr := B512DHWStepperLookup(b512DHWStepperBase)
	if b512DHWStepperErr != nil || !b512DHWStepperFound || !b512DHWStepperReflect.DeepEqual(b512DHWStepperAgain, b512DHWStepperRow) {
		b512DHWStepperT.Fatal("synthetic B524 axes alias storage")
	}
}

func TestAb97ConcurrentLookup(b512DHWStepperT *b512DHWStepperTesting.T) {
	// Exercise concurrent initialization as well as independent returned storage.
	b512DHWStepperCatalog = b512DHWStepperSync.OnceValues(b512DHWStepperLoadEntries)
	for b512DHWStepperIndex := 0; b512DHWStepperIndex < 16; b512DHWStepperIndex++ {
		b512DHWStepperT.Run(b512DHWStepperStrconv.Itoa(b512DHWStepperIndex), func(b512DHWStepperT *b512DHWStepperTesting.T) {
			b512DHWStepperT.Parallel()
			for b512DHWStepperRead := 0; b512DHWStepperRead < 10; b512DHWStepperRead++ {
				b512DHWStepperEntry := b512DHWStepperRequire(b512DHWStepperT, 3, "positive", "status")
				if b512DHWStepperEntry.B512DHWStepperValues["0"] != "okay" {
					b512DHWStepperT.Fatal("concurrent mutation escaped")
				}
				b512DHWStepperMutate(b512DHWStepperReflect.ValueOf(&b512DHWStepperEntry).Elem())
			}
		})
	}
}

func BenchmarkAb97Lookup(b512DHWStepperB *b512DHWStepperTesting.B) {
	b512DHWStepperKey := b512DHWStepperKey(b512DHWStepperByte(3), "positive", "status")
	if _, b512DHWStepperFound, b512DHWStepperErr := B512DHWStepperLookup(b512DHWStepperKey); b512DHWStepperErr != nil || !b512DHWStepperFound {
		b512DHWStepperB.Fatalf("lookup: %v %v", b512DHWStepperFound, b512DHWStepperErr)
	}
	b512DHWStepperB.ReportAllocs()
	b512DHWStepperB.ResetTimer()
	for b512DHWStepperIndex := 0; b512DHWStepperIndex < b512DHWStepperB.N; b512DHWStepperIndex++ {
		if _, b512DHWStepperFound, b512DHWStepperErr := B512DHWStepperLookup(b512DHWStepperKey); b512DHWStepperErr != nil || !b512DHWStepperFound {
			b512DHWStepperB.Fatalf("lookup: %v %v", b512DHWStepperFound, b512DHWStepperErr)
		}
	}
}
