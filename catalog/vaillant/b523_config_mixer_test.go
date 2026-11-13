package vaillant

import (
	b523ConfigMixerJSON "encoding/json"
	b523ConfigMixerReflect "reflect"
	b523ConfigMixerTesting "testing"
)

func b523ConfigMixerPointer[b523ConfigMixerT any](b523ConfigMixerValue b523ConfigMixerT) *b523ConfigMixerT {
	return &b523ConfigMixerValue
}
func b523ConfigMixerKey(b523ConfigMixerVariant string, b523ConfigMixerOperation uint8, b523ConfigMixerDirection, b523ConfigMixerField string) B523ConfigMixerKey {
	return B523ConfigMixerKey{B523ConfigMixerFamily: "vaillant", B523ConfigMixerVariant: b523ConfigMixerVariant, B523ConfigMixerPrimary: 0xb5, B523ConfigMixerSecondary: b523ConfigMixerPointer(uint8(0x23)), B523ConfigMixerOperation: b523ConfigMixerPointer(b523ConfigMixerOperation), B523ConfigMixerDirection: b523ConfigMixerDirection, B523ConfigMixerField: b523ConfigMixerField}
}
func b523ConfigMixerGet(b523ConfigMixerT *b523ConfigMixerTesting.T, b523ConfigMixerKey B523ConfigMixerKey) B523ConfigMixerEntry {
	b523ConfigMixerT.Helper()
	b523ConfigMixerEntry, b523ConfigMixerOK, b523ConfigMixerError := B523ConfigMixerLookup(b523ConfigMixerKey)
	if b523ConfigMixerError != nil || !b523ConfigMixerOK {
		b523ConfigMixerT.Fatalf("lookup: %v, %v", b523ConfigMixerOK, b523ConfigMixerError)
	}
	return b523ConfigMixerEntry
}
func TestAb99Finite(b523ConfigMixerT *b523ConfigMixerTesting.T) {
	b523ConfigMixerEntries, b523ConfigMixerError := B523ConfigMixerEntries()
	if b523ConfigMixerError != nil || len(b523ConfigMixerEntries) != 34 {
		b523ConfigMixerT.Fatalf("finite catalog: %d %v", len(b523ConfigMixerEntries), b523ConfigMixerError)
	}
	for _, b523ConfigMixerEntry := range b523ConfigMixerEntries {
		if b523ConfigMixerEntry.B523ConfigMixerGrade != "Hypothesis" || b523ConfigMixerEntry.B523ConfigMixerSupport != "Unknown" || b523ConfigMixerEntry.B523ConfigMixerPhysicalAdaptation != "Unknown" || b523ConfigMixerEntry.B523ConfigMixerExecutedMovement != "Unknown" || b523ConfigMixerEntry.B523ConfigMixerSelectedInterpretation != "Unknown" || b523ConfigMixerEntry.B523ConfigMixerUnknownMetadata != "Unknown" || !b523ConfigMixerEntry.B523ConfigMixerPreserveRaw {
			b523ConfigMixerT.Fatal("grade/admission")
		}
		if !b523ConfigMixerReflect.DeepEqual(b523ConfigMixerGet(b523ConfigMixerT, b523ConfigMixerEntry.B523ConfigMixerKey), b523ConfigMixerEntry) {
			b523ConfigMixerT.Fatal("lookup/accessor disagreement")
		}
		if b523ConfigMixerEntry.B523ConfigMixerKey.B523ConfigMixerGroup != nil || b523ConfigMixerEntry.B523ConfigMixerKey.B523ConfigMixerRequestedInstance != nil || b523ConfigMixerEntry.B523ConfigMixerKey.B523ConfigMixerRegister != nil || b523ConfigMixerEntry.B523ConfigMixerKey.B523ConfigMixerCommandRole != nil {
			b523ConfigMixerT.Fatal("invented B524 axes")
		}
	}
}
func TestAb99ExactAxes(b523ConfigMixerT *b523ConfigMixerTesting.T) {
	b523ConfigMixerControl := b523ConfigMixerKey("configuration", 0, "request", "S1")
	b523ConfigMixerChanges := map[string]func(*B523ConfigMixerKey){
		"ab99_family":  func(b523ConfigMixerKey *B523ConfigMixerKey) { b523ConfigMixerKey.B523ConfigMixerFamily = "other" },
		"ab99_variant": func(b523ConfigMixerKey *B523ConfigMixerKey) { b523ConfigMixerKey.B523ConfigMixerVariant = "other" },
		"ab99_primary": func(b523ConfigMixerKey *B523ConfigMixerKey) { b523ConfigMixerKey.B523ConfigMixerPrimary = 0 },
		"ab99_secondary": func(b523ConfigMixerKey *B523ConfigMixerKey) {
			b523ConfigMixerKey.B523ConfigMixerSecondary = b523ConfigMixerPointer(uint8(0))
		},
		"ab99_secondary_absent": func(b523ConfigMixerKey *B523ConfigMixerKey) { b523ConfigMixerKey.B523ConfigMixerSecondary = nil },
		"ab99_operation": func(b523ConfigMixerKey *B523ConfigMixerKey) {
			b523ConfigMixerKey.B523ConfigMixerOperation = b523ConfigMixerPointer(uint8(1))
		},
		"ab99_operation_absent": func(b523ConfigMixerKey *B523ConfigMixerKey) { b523ConfigMixerKey.B523ConfigMixerOperation = nil },
		"ab99_direction":        func(b523ConfigMixerKey *B523ConfigMixerKey) { b523ConfigMixerKey.B523ConfigMixerDirection = "positive" },
		"ab99_field":            func(b523ConfigMixerKey *B523ConfigMixerKey) { b523ConfigMixerKey.B523ConfigMixerField = "other" },
		"ab99_group_zero": func(b523ConfigMixerKey *B523ConfigMixerKey) {
			b523ConfigMixerKey.B523ConfigMixerGroup = b523ConfigMixerPointer(uint8(0))
		},
		"ab99_instance_zero": func(b523ConfigMixerKey *B523ConfigMixerKey) {
			b523ConfigMixerKey.B523ConfigMixerRequestedInstance = b523ConfigMixerPointer(uint8(0))
		},
		"ab99_register_zero": func(b523ConfigMixerKey *B523ConfigMixerKey) {
			b523ConfigMixerKey.B523ConfigMixerRegister = b523ConfigMixerPointer(uint16(0))
		},
		"ab99_role": func(b523ConfigMixerKey *B523ConfigMixerKey) {
			b523ConfigMixerKey.B523ConfigMixerCommandRole = b523ConfigMixerPointer("read")
		},
	}
	for b523ConfigMixerName, b523ConfigMixerChange := range b523ConfigMixerChanges {
		b523ConfigMixerT.Run(b523ConfigMixerName, func(b523ConfigMixerT *b523ConfigMixerTesting.T) {
			b523ConfigMixerGet(b523ConfigMixerT, b523ConfigMixerControl)
			b523ConfigMixerChanged := b523ConfigMixerControl
			b523ConfigMixerChange(&b523ConfigMixerChanged)
			b523ConfigMixerEntry, b523ConfigMixerOK, b523ConfigMixerError := B523ConfigMixerLookup(b523ConfigMixerChanged)
			if b523ConfigMixerError != nil || b523ConfigMixerOK || !b523ConfigMixerReflect.DeepEqual(b523ConfigMixerEntry, B523ConfigMixerEntry{}) {
				b523ConfigMixerT.Fatal("key axis ignored")
			}
		})
	}
}
func TestAb99OptionalAxes(b523ConfigMixerT *b523ConfigMixerTesting.T) {
	b523ConfigMixerSaved := b523ConfigMixerData
	b523ConfigMixerT.Cleanup(func() { b523ConfigMixerData = b523ConfigMixerSaved })
	b523ConfigMixerEntries, b523ConfigMixerError := B523ConfigMixerEntries()
	if b523ConfigMixerError != nil {
		b523ConfigMixerT.Fatal(b523ConfigMixerError)
	}
	b523ConfigMixerControl := b523ConfigMixerEntries[0]
	b523ConfigMixerControl.B523ConfigMixerKey.B523ConfigMixerGroup = b523ConfigMixerPointer(uint8(1))
	b523ConfigMixerControl.B523ConfigMixerKey.B523ConfigMixerRequestedInstance = b523ConfigMixerPointer(uint8(2))
	b523ConfigMixerControl.B523ConfigMixerKey.B523ConfigMixerRegister = b523ConfigMixerPointer(uint16(3))
	b523ConfigMixerControl.B523ConfigMixerKey.B523ConfigMixerCommandRole = b523ConfigMixerPointer("read")
	b523ConfigMixerData, b523ConfigMixerError = b523ConfigMixerJSON.Marshal([]B523ConfigMixerEntry{b523ConfigMixerControl})
	if b523ConfigMixerError != nil {
		b523ConfigMixerT.Fatal(b523ConfigMixerError)
	}
	b523ConfigMixerChanges := map[string]func(*B523ConfigMixerKey){
		"ab99_group": func(b523ConfigMixerKey *B523ConfigMixerKey) {
			b523ConfigMixerKey.B523ConfigMixerGroup = b523ConfigMixerPointer(uint8(0))
		},
		"ab99_group_absent": func(b523ConfigMixerKey *B523ConfigMixerKey) { b523ConfigMixerKey.B523ConfigMixerGroup = nil },
		"ab99_instance": func(b523ConfigMixerKey *B523ConfigMixerKey) {
			b523ConfigMixerKey.B523ConfigMixerRequestedInstance = b523ConfigMixerPointer(uint8(0))
		},
		"ab99_instance_absent": func(b523ConfigMixerKey *B523ConfigMixerKey) {
			b523ConfigMixerKey.B523ConfigMixerRequestedInstance = nil
		},
		"ab99_register": func(b523ConfigMixerKey *B523ConfigMixerKey) {
			b523ConfigMixerKey.B523ConfigMixerRegister = b523ConfigMixerPointer(uint16(0))
		},
		"ab99_register_absent": func(b523ConfigMixerKey *B523ConfigMixerKey) { b523ConfigMixerKey.B523ConfigMixerRegister = nil },
		"ab99_role": func(b523ConfigMixerKey *B523ConfigMixerKey) {
			b523ConfigMixerKey.B523ConfigMixerCommandRole = b523ConfigMixerPointer("write")
		},
		"ab99_role_absent": func(b523ConfigMixerKey *B523ConfigMixerKey) { b523ConfigMixerKey.B523ConfigMixerCommandRole = nil },
	}
	for b523ConfigMixerName, b523ConfigMixerChange := range b523ConfigMixerChanges {
		b523ConfigMixerT.Run(b523ConfigMixerName, func(b523ConfigMixerT *b523ConfigMixerTesting.T) {
			b523ConfigMixerGet(b523ConfigMixerT, b523ConfigMixerControl.B523ConfigMixerKey)
			b523ConfigMixerKey := b523ConfigMixerControl.B523ConfigMixerKey
			b523ConfigMixerChange(&b523ConfigMixerKey)
			if _, b523ConfigMixerOK, b523ConfigMixerError := B523ConfigMixerLookup(b523ConfigMixerKey); b523ConfigMixerOK || b523ConfigMixerError != nil {
				b523ConfigMixerT.Fatal("B524 axis ignored")
			}
		})
	}
}
func TestAb99OP00(b523ConfigMixerT *b523ConfigMixerTesting.T) {
	for b523ConfigMixerIndex, b523ConfigMixerField := range []string{"relay_groups_134", "relay_groups_256", "S1", "S2", "S3", "S4", "S5", "S6"} {
		b523ConfigMixerEntry := b523ConfigMixerGet(b523ConfigMixerT, b523ConfigMixerKey("configuration", 0, "request", b523ConfigMixerField))
		b523ConfigMixerWant := map[uint8]string{0: "kindA", 1: "kindB", 2: "digital"}
		if b523ConfigMixerIndex < 2 {
			b523ConfigMixerWant = map[uint8]string{0: "output", 1: "mixer"}
		}
		if *b523ConfigMixerEntry.B523ConfigMixerPosition != uint8(b523ConfigMixerIndex+3) || *b523ConfigMixerEntry.B523ConfigMixerDefault != 0 || *b523ConfigMixerEntry.B523ConfigMixerRequestExtent != 11 || *b523ConfigMixerEntry.B523ConfigMixerPositiveExtent != 4 || b523ConfigMixerEntry.B523ConfigMixerWidthBits != nil || !b523ConfigMixerReflect.DeepEqual(b523ConfigMixerEntry.B523ConfigMixerValues, b523ConfigMixerWant) {
			b523ConfigMixerT.Fatal("OP00 mapping", b523ConfigMixerField)
		}
		if b523ConfigMixerIndex >= 2 && b523ConfigMixerEntry.B523ConfigMixerOtherValues != "undefined-or-ignored" {
			b523ConfigMixerT.Fatal("undefined enum")
		}
	}
}
func TestAb99OP04(b523ConfigMixerT *b523ConfigMixerTesting.T) {
	for b523ConfigMixerIndex, b523ConfigMixerField := range []string{"R1/R7/R8", "R2/R9/R10", "R3/R11/R12"} {
		b523ConfigMixerEntry := b523ConfigMixerGet(b523ConfigMixerT, b523ConfigMixerKey("configuration", 4, "request", b523ConfigMixerField))
		if *b523ConfigMixerEntry.B523ConfigMixerPosition != uint8(b523ConfigMixerIndex+3) || !b523ConfigMixerReflect.DeepEqual(b523ConfigMixerEntry.B523ConfigMixerValues, map[uint8]string{0: "output", 1: "mixer"}) || *b523ConfigMixerEntry.B523ConfigMixerDefault != 0 {
			b523ConfigMixerT.Fatal("OP04 relays")
		}
	}
	for b523ConfigMixerIndex, b523ConfigMixerField := range []string{"S1", "S2", "S3", "S4", "S5", "S6", "S7", "S8", "S9", "S10", "S11", "S12", "S13"} {
		b523ConfigMixerEntry := b523ConfigMixerGet(b523ConfigMixerT, b523ConfigMixerKey("configuration", 4, "request", b523ConfigMixerField))
		if *b523ConfigMixerEntry.B523ConfigMixerPosition != []uint8{6, 6, 7, 7, 8, 8, 9, 9, 10, 10, 11, 12, 12}[b523ConfigMixerIndex] || *b523ConfigMixerEntry.B523ConfigMixerRequestExtent != 13 || *b523ConfigMixerEntry.B523ConfigMixerPositiveExtent != 4 || b523ConfigMixerEntry.B523ConfigMixerWidthBits != nil || b523ConfigMixerEntry.B523ConfigMixerOneHot == nil || *b523ConfigMixerEntry.B523ConfigMixerOneHot || !b523ConfigMixerReflect.DeepEqual(b523ConfigMixerEntry.B523ConfigMixerAlternatives, []string{"independent sensor flags", "competing nibble-code interpretation unresolved", "repeated sensor associations unresolved"}) {
			b523ConfigMixerT.Fatal("OP04 alternatives", b523ConfigMixerField)
		}
		b523ConfigMixerOffset := []uint8{0, 4, 0, 4, 0, 4, 0, 4, 0, 4, 0, 0, 4}[b523ConfigMixerIndex]
		b523ConfigMixerCount := 4
		if b523ConfigMixerIndex == 10 {
			b523ConfigMixerCount = 8
		}
		if len(b523ConfigMixerEntry.B523ConfigMixerBits) != b523ConfigMixerCount {
			b523ConfigMixerT.Fatal("bit count")
		}
		for b523ConfigMixerBitIndex, b523ConfigMixerBit := range b523ConfigMixerEntry.B523ConfigMixerBits {
			b523ConfigMixerLabel := ""
			if b523ConfigMixerBitIndex < 3 {
				b523ConfigMixerLabel = []string{"kindA", "kindB", "digital"}[b523ConfigMixerBitIndex]
			} else if (b523ConfigMixerIndex == 8 || b523ConfigMixerIndex == 9) && b523ConfigMixerBitIndex == 3 {
				b523ConfigMixerLabel = "PWM"
			}
			if b523ConfigMixerBit.B523ConfigMixerIndex != b523ConfigMixerOffset+uint8(b523ConfigMixerBitIndex) || b523ConfigMixerBit.B523ConfigMixerReserved != (b523ConfigMixerLabel == "") || (b523ConfigMixerLabel == "" && b523ConfigMixerBit.B523ConfigMixerLabel != nil) || (b523ConfigMixerLabel != "" && (b523ConfigMixerBit.B523ConfigMixerLabel == nil || *b523ConfigMixerBit.B523ConfigMixerLabel != b523ConfigMixerLabel)) {
				b523ConfigMixerT.Fatal("sensor bit mapping", b523ConfigMixerField, b523ConfigMixerBitIndex)
			}
			b523ConfigMixerDefault := uint8(0)
			if b523ConfigMixerBitIndex == 0 {
				b523ConfigMixerDefault = 1
			}
			if b523ConfigMixerBit.B523ConfigMixerDefault != b523ConfigMixerDefault {
				b523ConfigMixerT.Fatal("sensor defaults")
			}
			for b523ConfigMixerRaw := 0; b523ConfigMixerRaw < 256; b523ConfigMixerRaw++ {
				if ((uint8(b523ConfigMixerRaw) >> b523ConfigMixerBit.B523ConfigMixerIndex) & 1) != ((uint8(b523ConfigMixerRaw) >> uint8(int(b523ConfigMixerOffset)+b523ConfigMixerBitIndex)) & 1) {
					b523ConfigMixerT.Fatal("independent bit extraction")
				}
			}
		}
		if uint8(3<<b523ConfigMixerOffset)&(1<<b523ConfigMixerEntry.B523ConfigMixerBits[0].B523ConfigMixerIndex) == 0 || uint8(3<<b523ConfigMixerOffset)&(1<<b523ConfigMixerEntry.B523ConfigMixerBits[1].B523ConfigMixerIndex) == 0 {
			b523ConfigMixerT.Fatal("two flags preserved")
		}
	}
}
func TestAb99Mixer(b523ConfigMixerT *b523ConfigMixerTesting.T) {
	for _, b523ConfigMixerVariant := range []string{"mixer-smaller", "mixer-larger"} {
		b523ConfigMixerCircuit := b523ConfigMixerGet(b523ConfigMixerT, b523ConfigMixerKey(b523ConfigMixerVariant, 2, "request", "circuit"))
		b523ConfigMixerWant := map[uint8]string{0: "first group", 1: "second group"}
		if b523ConfigMixerVariant == "mixer-larger" {
			b523ConfigMixerWant[2] = "third group"
		}
		if *b523ConfigMixerCircuit.B523ConfigMixerPosition != 3 || !b523ConfigMixerReflect.DeepEqual(b523ConfigMixerCircuit.B523ConfigMixerValues, b523ConfigMixerWant) || b523ConfigMixerCircuit.B523ConfigMixerOtherValues != "ignored" {
			b523ConfigMixerT.Fatal("circuit eligibility")
		}
		b523ConfigMixerMode := b523ConfigMixerGet(b523ConfigMixerT, b523ConfigMixerKey(b523ConfigMixerVariant, 2, "request", "mode"))
		if *b523ConfigMixerMode.B523ConfigMixerPosition != 4 || !b523ConfigMixerReflect.DeepEqual(b523ConfigMixerMode.B523ConfigMixerValues, map[uint8]string{0: "stop/pump-off/close", 1: "heating-regulation", 2: "cooling-regulation", 3: "pump-off/open", 4: "pump-on/close", 5: "pump-on/open"}) || b523ConfigMixerMode.B523ConfigMixerOtherValues != "ignored" {
			b523ConfigMixerT.Fatal("modes")
		}
		b523ConfigMixerSetpoint := b523ConfigMixerGet(b523ConfigMixerT, b523ConfigMixerKey(b523ConfigMixerVariant, 2, "request", "setpoint"))
		if *b523ConfigMixerSetpoint.B523ConfigMixerPosition != 5 || *b523ConfigMixerSetpoint.B523ConfigMixerWidthBits != 8 || *b523ConfigMixerSetpoint.B523ConfigMixerDefault != 0 || *b523ConfigMixerSetpoint.B523ConfigMixerUnit != "Celsius" || !b523ConfigMixerReflect.DeepEqual(b523ConfigMixerSetpoint.B523ConfigMixerRelevantModes, []uint8{1, 2}) || !b523ConfigMixerReflect.DeepEqual(b523ConfigMixerSetpoint.B523ConfigMixerIrrelevantModes, []uint8{0, 3, 4, 5}) || len(b523ConfigMixerSetpoint.B523ConfigMixerCandidates) != 2 {
			b523ConfigMixerT.Fatal("setpoint relevance")
		}
		b523ConfigMixerSigned := b523ConfigMixerSetpoint.B523ConfigMixerCandidates[0]
		b523ConfigMixerOther := b523ConfigMixerSetpoint.B523ConfigMixerCandidates[1]
		b523ConfigMixerRaw := uint8(0xbe)
		if b523ConfigMixerSigned.B523ConfigMixerSigned == nil || !*b523ConfigMixerSigned.B523ConfigMixerSigned || *b523ConfigMixerSigned.B523ConfigMixerScale != 0.5 || *b523ConfigMixerSigned.B523ConfigMixerMinimum != -64 || *b523ConfigMixerSigned.B523ConfigMixerMaximum != 63.5 || float64(int8(b523ConfigMixerRaw))**b523ConfigMixerSigned.B523ConfigMixerScale != -33 || *b523ConfigMixerOther.B523ConfigMixerRaw != 190 || *b523ConfigMixerOther.B523ConfigMixerValue != 95 || *b523ConfigMixerOther.B523ConfigMixerClamp != 95 || b523ConfigMixerOther.B523ConfigMixerScale != nil || b523ConfigMixerOther.B523ConfigMixerSigned != nil {
			b523ConfigMixerT.Fatal("competing BE interpretations")
		}
		b523ConfigMixerFeedback := b523ConfigMixerGet(b523ConfigMixerT, b523ConfigMixerKey(b523ConfigMixerVariant, 2, "positive", "feedback"))
		if *b523ConfigMixerFeedback.B523ConfigMixerPosition != 3 || b523ConfigMixerFeedback.B523ConfigMixerWidthBits != nil || b523ConfigMixerFeedback.B523ConfigMixerValues != nil || b523ConfigMixerFeedback.B523ConfigMixerCandidates != nil {
			b523ConfigMixerT.Fatal("raw feedback")
		}
		b523ConfigMixerMovement := b523ConfigMixerGet(b523ConfigMixerT, b523ConfigMixerKey(b523ConfigMixerVariant, 2, "positive", "movement"))
		if b523ConfigMixerMovement.B523ConfigMixerSigned == nil || !*b523ConfigMixerMovement.B523ConfigMixerSigned || *b523ConfigMixerMovement.B523ConfigMixerPosition != 4 || *b523ConfigMixerMovement.B523ConfigMixerWidthBits != 8 || *b523ConfigMixerMovement.B523ConfigMixerUnit != "percent" || b523ConfigMixerMovement.B523ConfigMixerOtherValues != "unresolved" || !b523ConfigMixerReflect.DeepEqual(b523ConfigMixerMovement.B523ConfigMixerIntervals, []B523ConfigMixerInterval{{-100, -1, "closing"}, {0, 0, "stopped"}, {1, 100, "opening"}}) {
			b523ConfigMixerT.Fatal("movement ranges")
		}
		for b523ConfigMixerIndex, b523ConfigMixerRaw := range []uint8{0x9c, 0x64, 0x9b, 0x65, 0} {
			b523ConfigMixerValue := int(int8(b523ConfigMixerRaw))
			b523ConfigMixerLabel := "unresolved"
			for _, b523ConfigMixerInterval := range b523ConfigMixerMovement.B523ConfigMixerIntervals {
				if b523ConfigMixerValue >= b523ConfigMixerInterval.B523ConfigMixerMinimum && b523ConfigMixerValue <= b523ConfigMixerInterval.B523ConfigMixerMaximum {
					b523ConfigMixerLabel = b523ConfigMixerInterval.B523ConfigMixerLabel
				}
			}
			if b523ConfigMixerLabel != []string{"closing", "opening", "unresolved", "unresolved", "stopped"}[b523ConfigMixerIndex] {
				b523ConfigMixerT.Fatal("movement polarity/bounds")
			}
		}
	}
}
func b523ConfigMixerMutate(b523ConfigMixerValue b523ConfigMixerReflect.Value) {
	switch b523ConfigMixerValue.Kind() {
	case b523ConfigMixerReflect.Pointer:
		if !b523ConfigMixerValue.IsNil() {
			b523ConfigMixerMutate(b523ConfigMixerValue.Elem())
		}
	case b523ConfigMixerReflect.Struct:
		for b523ConfigMixerIndex := 0; b523ConfigMixerIndex < b523ConfigMixerValue.NumField(); b523ConfigMixerIndex++ {
			b523ConfigMixerMutate(b523ConfigMixerValue.Field(b523ConfigMixerIndex))
		}
	case b523ConfigMixerReflect.Slice:
		for b523ConfigMixerIndex := 0; b523ConfigMixerIndex < b523ConfigMixerValue.Len(); b523ConfigMixerIndex++ {
			b523ConfigMixerMutate(b523ConfigMixerValue.Index(b523ConfigMixerIndex))
		}
	case b523ConfigMixerReflect.Map:
		for _, b523ConfigMixerKey := range b523ConfigMixerValue.MapKeys() {
			b523ConfigMixerValue.SetMapIndex(b523ConfigMixerKey, b523ConfigMixerReflect.ValueOf("ab99_changed"))
		}
	case b523ConfigMixerReflect.String:
		b523ConfigMixerValue.SetString("ab99_changed")
	case b523ConfigMixerReflect.Bool:
		b523ConfigMixerValue.SetBool(!b523ConfigMixerValue.Bool())
	case b523ConfigMixerReflect.Uint8, b523ConfigMixerReflect.Uint16:
		b523ConfigMixerValue.SetUint(99)
	case b523ConfigMixerReflect.Int:
		b523ConfigMixerValue.SetInt(99)
	case b523ConfigMixerReflect.Float64:
		b523ConfigMixerValue.SetFloat(99)
	}
}
func TestAb99Copies(b523ConfigMixerT *b523ConfigMixerTesting.T) {
	b523ConfigMixerBefore, b523ConfigMixerError := B523ConfigMixerEntries()
	if b523ConfigMixerError != nil {
		b523ConfigMixerT.Fatal(b523ConfigMixerError)
	}
	b523ConfigMixerReturned, b523ConfigMixerError := B523ConfigMixerEntries()
	if b523ConfigMixerError != nil {
		b523ConfigMixerT.Fatal(b523ConfigMixerError)
	}
	b523ConfigMixerMutate(b523ConfigMixerReflect.ValueOf(&b523ConfigMixerReturned).Elem())
	b523ConfigMixerAfter, b523ConfigMixerError := B523ConfigMixerEntries()
	if b523ConfigMixerError != nil || !b523ConfigMixerReflect.DeepEqual(b523ConfigMixerBefore, b523ConfigMixerAfter) {
		b523ConfigMixerT.Fatal("accessor storage shared")
	}
	for _, b523ConfigMixerEntry := range b523ConfigMixerBefore {
		b523ConfigMixerGot := b523ConfigMixerGet(b523ConfigMixerT, b523ConfigMixerEntry.B523ConfigMixerKey)
		b523ConfigMixerMutate(b523ConfigMixerReflect.ValueOf(&b523ConfigMixerGot).Elem())
		if !b523ConfigMixerReflect.DeepEqual(b523ConfigMixerGet(b523ConfigMixerT, b523ConfigMixerEntry.B523ConfigMixerKey), b523ConfigMixerEntry) {
			b523ConfigMixerT.Fatal("lookup storage shared")
		}
	}
}
func TestAb99Malformed(b523ConfigMixerT *b523ConfigMixerTesting.T) {
	b523ConfigMixerSaved := b523ConfigMixerData
	b523ConfigMixerT.Cleanup(func() { b523ConfigMixerData = b523ConfigMixerSaved })
	b523ConfigMixerEntries, b523ConfigMixerError := B523ConfigMixerEntries()
	if b523ConfigMixerError != nil {
		b523ConfigMixerT.Fatal(b523ConfigMixerError)
	}
	b523ConfigMixerDuplicate, b523ConfigMixerError := b523ConfigMixerJSON.Marshal([]B523ConfigMixerEntry{b523ConfigMixerEntries[0], b523ConfigMixerEntries[0]})
	if b523ConfigMixerError != nil {
		b523ConfigMixerT.Fatal(b523ConfigMixerError)
	}
	for _, b523ConfigMixerInvalid := range [][]byte{[]byte("{"), []byte("null"), []byte("[]"), []byte("[{}]"), append(append([]byte{}, b523ConfigMixerSaved...), []byte(" true")...), b523ConfigMixerDuplicate, []byte(`[{"ab99_unknown":1}]`)} {
		b523ConfigMixerData = b523ConfigMixerInvalid
		if b523ConfigMixerGot, b523ConfigMixerError := B523ConfigMixerEntries(); b523ConfigMixerError == nil || b523ConfigMixerGot != nil {
			b523ConfigMixerT.Fatal("invalid catalog exposed")
		}
		if b523ConfigMixerGot, b523ConfigMixerOK, b523ConfigMixerError := B523ConfigMixerLookup(b523ConfigMixerKey("missing", 0, "request", "missing")); b523ConfigMixerError == nil || b523ConfigMixerOK || !b523ConfigMixerReflect.DeepEqual(b523ConfigMixerGot, B523ConfigMixerEntry{}) {
			b523ConfigMixerT.Fatal("invalid catalog hidden as miss")
		}
	}
}

func TestAb99AbsentZero(b523ConfigMixerT *b523ConfigMixerTesting.T) {
	b523ConfigMixerSaved := b523ConfigMixerData
	b523ConfigMixerT.Cleanup(func() { b523ConfigMixerData = b523ConfigMixerSaved })
	b523ConfigMixerEntries, b523ConfigMixerError := B523ConfigMixerEntries()
	if b523ConfigMixerError != nil {
		b523ConfigMixerT.Fatal(b523ConfigMixerError)
	}
	b523ConfigMixerControl := b523ConfigMixerEntries[0]
	b523ConfigMixerControl.B523ConfigMixerKey.B523ConfigMixerGroup = b523ConfigMixerPointer(uint8(0))
	b523ConfigMixerControl.B523ConfigMixerKey.B523ConfigMixerRequestedInstance = b523ConfigMixerPointer(uint8(0))
	b523ConfigMixerControl.B523ConfigMixerKey.B523ConfigMixerRegister = b523ConfigMixerPointer(uint16(0))
	b523ConfigMixerControl.B523ConfigMixerKey.B523ConfigMixerCommandRole = b523ConfigMixerPointer("")
	b523ConfigMixerData, b523ConfigMixerError = b523ConfigMixerJSON.Marshal([]B523ConfigMixerEntry{b523ConfigMixerControl})
	if b523ConfigMixerError != nil {
		b523ConfigMixerT.Fatal(b523ConfigMixerError)
	}
	b523ConfigMixerChanges := map[string]func(*B523ConfigMixerKey){
		"ab99_operation": func(b523ConfigMixerKey *B523ConfigMixerKey) { b523ConfigMixerKey.B523ConfigMixerOperation = nil },
		"ab99_group":     func(b523ConfigMixerKey *B523ConfigMixerKey) { b523ConfigMixerKey.B523ConfigMixerGroup = nil },
		"ab99_instance": func(b523ConfigMixerKey *B523ConfigMixerKey) {
			b523ConfigMixerKey.B523ConfigMixerRequestedInstance = nil
		},
		"ab99_register": func(b523ConfigMixerKey *B523ConfigMixerKey) { b523ConfigMixerKey.B523ConfigMixerRegister = nil },
		"ab99_role":     func(b523ConfigMixerKey *B523ConfigMixerKey) { b523ConfigMixerKey.B523ConfigMixerCommandRole = nil },
	}
	for b523ConfigMixerName, b523ConfigMixerChange := range b523ConfigMixerChanges {
		b523ConfigMixerT.Run(b523ConfigMixerName, func(b523ConfigMixerT *b523ConfigMixerTesting.T) {
			b523ConfigMixerGet(b523ConfigMixerT, b523ConfigMixerControl.B523ConfigMixerKey)
			b523ConfigMixerKey := b523ConfigMixerControl.B523ConfigMixerKey
			b523ConfigMixerChange(&b523ConfigMixerKey)
			if _, b523ConfigMixerOK, b523ConfigMixerError := B523ConfigMixerLookup(b523ConfigMixerKey); b523ConfigMixerOK || b523ConfigMixerError != nil {
				b523ConfigMixerT.Fatal("absent axis conflated with explicit zero")
			}
		})
	}
}

func TestAb99RequiredCoordinates(b523ConfigMixerT *b523ConfigMixerTesting.T) {
	b523ConfigMixerSaved := b523ConfigMixerData
	b523ConfigMixerT.Cleanup(func() { b523ConfigMixerData = b523ConfigMixerSaved })
	for _, b523ConfigMixerAxis := range []string{"primary", "secondary", "operation"} {
		for _, b523ConfigMixerForm := range []string{"missing", "null", "wrong"} {
			if b523ConfigMixerAxis == "operation" && b523ConfigMixerForm == "wrong" {
				continue
			}
			b523ConfigMixerT.Run(b523ConfigMixerAxis+"/"+b523ConfigMixerForm, func(b523ConfigMixerT *b523ConfigMixerTesting.T) {
				var b523ConfigMixerEntries []map[string]any
				if b523ConfigMixerError := b523ConfigMixerJSON.Unmarshal(b523ConfigMixerSaved, &b523ConfigMixerEntries); b523ConfigMixerError != nil {
					b523ConfigMixerT.Fatal(b523ConfigMixerError)
				}
				b523ConfigMixerKey := b523ConfigMixerEntries[0]["key"].(map[string]any)
				switch b523ConfigMixerForm {
				case "missing":
					delete(b523ConfigMixerKey, b523ConfigMixerAxis)
				case "null":
					b523ConfigMixerKey[b523ConfigMixerAxis] = nil
				case "wrong":
					b523ConfigMixerKey[b523ConfigMixerAxis] = 0
				}
				var b523ConfigMixerError error
				b523ConfigMixerData, b523ConfigMixerError = b523ConfigMixerJSON.Marshal(b523ConfigMixerEntries)
				if b523ConfigMixerError != nil {
					b523ConfigMixerT.Fatal(b523ConfigMixerError)
				}
				if b523ConfigMixerGot, b523ConfigMixerError := B523ConfigMixerEntries(); b523ConfigMixerError == nil || b523ConfigMixerGot != nil {
					b523ConfigMixerT.Fatal("invalid coordinates exposed")
				}
				if b523ConfigMixerGot, b523ConfigMixerOK, b523ConfigMixerError := B523ConfigMixerLookup(B523ConfigMixerKey{}); b523ConfigMixerError == nil || b523ConfigMixerOK || !b523ConfigMixerReflect.DeepEqual(b523ConfigMixerGot, B523ConfigMixerEntry{}) {
					b523ConfigMixerT.Fatal("invalid coordinates hidden as miss")
				}
			})
		}
	}
}
