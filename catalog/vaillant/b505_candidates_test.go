package vaillant

import (
	b505CandidateReflect "reflect"
	b505CandidateTesting "testing"
)

func b505CandidateByte(b505CandidateValue byte) *byte       { return &b505CandidateValue }
func b505CandidateString(b505CandidateValue string) *string { return &b505CandidateValue }

func b505CandidateKey(b505CandidateField string, b505CandidateOperation *byte) B505CandidateKey {
	return B505CandidateKey{B505CandidateFamily: "vaillant", B505CandidateVariant: "B505", B505CandidatePrimary: 0xB5, B505CandidateSecondary: b505CandidateByte(5), B505CandidateOperation: b505CandidateOperation, B505CandidateField: b505CandidateField}
}

func b505CandidateLookup(b505CandidateT *b505CandidateTesting.T, b505CandidateField string, b505CandidateOperation *byte) B505CandidateEntry {
	b505CandidateT.Helper()
	b505CandidateEntry, b505CandidateOK := B505CandidateLookup(b505CandidateKey(b505CandidateField, b505CandidateOperation))
	if !b505CandidateOK {
		b505CandidateT.Fatalf("missing %s", b505CandidateField)
	}
	return b505CandidateEntry
}

func TestAb55CatalogContract(b505CandidateT *b505CandidateTesting.T) {
	b505CandidateCatalog := B505Candidates()
	if b505CandidateCatalog.B505CandidateMetadataPresence != "Proven" || b505CandidateCatalog.B505CandidateQualification != "Unknown" || len(b505CandidateCatalog.B505CandidateEntries) != 8 || len(b505CandidateCatalog.B505CandidateLimitations) != 4 {
		b505CandidateT.Fatal("finite catalog or evidence boundary changed")
	}
	b505CandidateFields := []string{"temperature_1", "temperature_2", "dhw_setpoint", "collector_temperature", "collector_error", "pump_power_display", "remote_control", "priority"}
	b505CandidateOperations := []*byte{b505CandidateByte(0x32), b505CandidateByte(0x32), b505CandidateByte(0x61), b505CandidateByte(0x4D), b505CandidateByte(0x4D), nil, nil, b505CandidateByte(0x4E)}
	b505CandidateWidths := []*int{new(int), new(int), nil, nil, new(int), nil, new(int), new(int)}
	*b505CandidateWidths[0], *b505CandidateWidths[1], *b505CandidateWidths[4], *b505CandidateWidths[6], *b505CandidateWidths[7] = 16, 16, 8, 8, 8
	for b505CandidateIndex, b505CandidateField := range b505CandidateFields {
		b505CandidateEntry := b505CandidateLookup(b505CandidateT, b505CandidateField, b505CandidateOperations[b505CandidateIndex])
		if !b505CandidateReflect.DeepEqual(b505CandidateEntry, b505CandidateCatalog.B505CandidateEntries[b505CandidateIndex]) || !b505CandidateReflect.DeepEqual(b505CandidateEntry.B505CandidateKey, b505CandidateKey(b505CandidateField, b505CandidateOperations[b505CandidateIndex])) {
			b505CandidateT.Fatalf("accessor/lookup identity: %s", b505CandidateField)
		}
		if b505CandidateEntry.B505CandidatePosition != nil || !b505CandidateReflect.DeepEqual(b505CandidateEntry.B505CandidateWidthBits, b505CandidateWidths[b505CandidateIndex]) || b505CandidateEntry.B505CandidateConversion != "Unknown" || b505CandidateEntry.B505CandidateBounds != "Unknown" || b505CandidateEntry.B505CandidateGrade != "Hypothesis" {
			b505CandidateT.Fatalf("unsupported coordinate, width, conversion or grade: %s", b505CandidateField)
		}
		if b505CandidateEntry.B505CandidateUnit != nil || b505CandidateEntry.B505CandidateFreshness != "Unknown" || b505CandidateEntry.B505CandidateExecution != "Unknown" || b505CandidateEntry.B505CandidatePhysicalSupport != "Unknown" || b505CandidateEntry.B505CandidateAcknowledgement != "Unknown" || b505CandidateEntry.B505CandidateAcceptance != "Unknown" || b505CandidateEntry.B505CandidatePersistence != "Unknown" || b505CandidateEntry.B505CandidateDurationEnforcement != "Unknown" {
			b505CandidateT.Fatalf("qualification promoted: %s", b505CandidateField)
		}
		b505CandidateHasDuration := b505CandidateField == "dhw_setpoint" || b505CandidateField == "collector_temperature" || b505CandidateField == "collector_error"
		if b505CandidateHasDuration {
			if b505CandidateEntry.B505CandidateDurationSeconds == nil || *b505CandidateEntry.B505CandidateDurationSeconds != 900 {
				b505CandidateT.Fatalf("duration: %s", b505CandidateField)
			}
		} else if b505CandidateEntry.B505CandidateDurationSeconds != nil {
			b505CandidateT.Fatalf("invented duration: %s", b505CandidateField)
		}
	}
}

func TestAb55ExactKeyAxes(b505CandidateT *b505CandidateTesting.T) {
	b505CandidateControl := b505CandidateKey("temperature_1", b505CandidateByte(0x32))
	b505CandidateMutations := map[string]func(*B505CandidateKey){
		"ab55_family":  func(b505CandidateKey *B505CandidateKey) { b505CandidateKey.B505CandidateFamily = "other-family" },
		"ab55_variant": func(b505CandidateKey *B505CandidateKey) { b505CandidateKey.B505CandidateVariant = "other-variant" },
		"ab55_primary": func(b505CandidateKey *B505CandidateKey) { b505CandidateKey.B505CandidatePrimary = 0 },
		"ab55_secondary": func(b505CandidateKey *B505CandidateKey) {
			b505CandidateKey.B505CandidateSecondary = b505CandidateByte(6)
		},
		"ab55_secondary_absent": func(b505CandidateKey *B505CandidateKey) { b505CandidateKey.B505CandidateSecondary = nil },
		"ab55_secondary_zero": func(b505CandidateKey *B505CandidateKey) {
			b505CandidateKey.B505CandidateSecondary = b505CandidateByte(0)
		},
		"ab55_operation": func(b505CandidateKey *B505CandidateKey) {
			b505CandidateKey.B505CandidateOperation = b505CandidateByte(0x33)
		},
		"ab55_operation_absent": func(b505CandidateKey *B505CandidateKey) { b505CandidateKey.B505CandidateOperation = nil },
		"ab55_operation_zero": func(b505CandidateKey *B505CandidateKey) {
			b505CandidateKey.B505CandidateOperation = b505CandidateByte(0)
		},
		"ab55_direction": func(b505CandidateKey *B505CandidateKey) {
			b505CandidateKey.B505CandidateDirection = b505CandidateString("request")
		},
		"ab55_direction_empty": func(b505CandidateKey *B505CandidateKey) {
			b505CandidateKey.B505CandidateDirection = b505CandidateString("")
		},
		"ab55_field": func(b505CandidateKey *B505CandidateKey) { b505CandidateKey.B505CandidateField = "other-field" },
	}
	for b505CandidateName, b505CandidateMutate := range b505CandidateMutations {
		b505CandidateT.Run(b505CandidateName, func(b505CandidateT *b505CandidateTesting.T) {
			b505CandidateKey := b505CandidateControl
			b505CandidateMutate(&b505CandidateKey)
			if _, b505CandidateOK := B505CandidateLookup(b505CandidateKey); b505CandidateOK {
				b505CandidateT.Fatal("changed axis matched")
			}
			if _, b505CandidateOK := B505CandidateLookup(b505CandidateControl); !b505CandidateOK {
				b505CandidateT.Fatal("valid control missed")
			}
		})
	}
	for _, b505CandidateAxis := range []string{"secondary", "operation", "direction"} {
		b505CandidateLeft, b505CandidateRight := b505CandidateControl, b505CandidateControl
		switch b505CandidateAxis {
		case "secondary":
			b505CandidateLeft.B505CandidateSecondary, b505CandidateRight.B505CandidateSecondary = nil, b505CandidateByte(0)
		case "operation":
			b505CandidateLeft.B505CandidateOperation, b505CandidateRight.B505CandidateOperation = nil, b505CandidateByte(0)
		case "direction":
			b505CandidateLeft.B505CandidateDirection, b505CandidateRight.B505CandidateDirection = nil, b505CandidateString("")
		}
		if b505CandidateKeysEqual(b505CandidateLeft, b505CandidateRight) || b505CandidateKeysEqual(b505CandidateRight, b505CandidateLeft) || !b505CandidateKeysEqual(b505CandidateLeft, b505CandidateLeft) || !b505CandidateKeysEqual(b505CandidateRight, b505CandidateRight) {
			b505CandidateT.Fatalf("absent and zero collapsed: %s", b505CandidateAxis)
		}
	}
	b505CandidateAbsent := b505CandidateKey("remote_control", nil)
	if _, b505CandidateOK := B505CandidateLookup(b505CandidateAbsent); !b505CandidateOK {
		b505CandidateT.Fatal("absent operation control missed")
	}
	b505CandidateAbsent.B505CandidateOperation = b505CandidateByte(0)
	if _, b505CandidateOK := B505CandidateLookup(b505CandidateAbsent); b505CandidateOK {
		b505CandidateT.Fatal("zero operation matched absent")
	}
	b505CandidateSecond := b505CandidateLookup(b505CandidateT, "temperature_2", b505CandidateByte(0x32))
	if b505CandidateSecond.B505CandidateKey.B505CandidateField != "temperature_2" {
		b505CandidateT.Fatal("field identities collapsed")
	}
}

func TestAb55IndependentSentinels(b505CandidateT *b505CandidateTesting.T) {
	b505CandidateFirst := b505CandidateLookup(b505CandidateT, "temperature_1", b505CandidateByte(0x32))
	b505CandidateSecond := b505CandidateLookup(b505CandidateT, "temperature_2", b505CandidateByte(0x32))
	for _, b505CandidateEntry := range []B505CandidateEntry{b505CandidateFirst, b505CandidateSecond} {
		if !b505CandidateReflect.DeepEqual(b505CandidateEntry.B505CandidateBranches, map[string]string{"8000": "preserve-existing", "0000": "circuit-default"}) || b505CandidateEntry.B505CandidateFallback != "numeric-candidate" {
			b505CandidateT.Fatal("sentinel branches changed")
		}
	}
	for _, b505CandidateCase := range []struct {
		B505CandidateRaw  [2]string
		B505CandidateWant [2]string
	}{
		{[2]string{"8000", "0000"}, [2]string{"preserve-existing", "circuit-default"}},
		{[2]string{"0000", "8000"}, [2]string{"circuit-default", "preserve-existing"}},
		{[2]string{"8000", "8000"}, [2]string{"preserve-existing", "preserve-existing"}},
		{[2]string{"0001", "0000"}, [2]string{"numeric-candidate", "circuit-default"}},
	} {
		b505CandidateBefore := b505CandidateCase.B505CandidateRaw
		b505CandidateGot := [2]string{b505CandidateFirst.B505CandidateSelectRaw(b505CandidateCase.B505CandidateRaw[0]), b505CandidateSecond.B505CandidateSelectRaw(b505CandidateCase.B505CandidateRaw[1])}
		if b505CandidateGot != b505CandidateCase.B505CandidateWant || b505CandidateCase.B505CandidateRaw != b505CandidateBefore {
			b505CandidateT.Fatalf("raw %v: got %v", b505CandidateCase.B505CandidateRaw, b505CandidateGot)
		}
	}
	b505CandidateOther := b505CandidateLookup(b505CandidateT, "dhw_setpoint", b505CandidateByte(0x61))
	if b505CandidateOther.B505CandidateSelectRaw("8000") != "Unknown" || len(b505CandidateOther.B505CandidateBranches) != 0 {
		b505CandidateT.Fatal("OP32 sentinel leaked into OP61")
	}
}

func TestAb55OverridesAndPriority(b505CandidateT *b505CandidateTesting.T) {
	b505CandidateSensor := b505CandidateLookup(b505CandidateT, "collector_error", b505CandidateByte(0x4D))
	if !b505CandidateReflect.DeepEqual(b505CandidateSensor.B505CandidateBranches, map[string]string{"55": "short-circuit", "AA": "open"}) || b505CandidateSensor.B505CandidateSelectRaw("55") != "short-circuit" || b505CandidateSensor.B505CandidateSelectRaw("AA") != "open" || b505CandidateSensor.B505CandidateSelectRaw("54") != "Unknown" {
		b505CandidateT.Fatal("sensor mappings changed")
	}
	b505CandidatePump := b505CandidateLookup(b505CandidateT, "pump_power_display", nil)
	if !b505CandidateReflect.DeepEqual(b505CandidatePump.B505CandidateDisplaySteps, []int{0, 25, 50, 75, 100}) {
		b505CandidateT.Fatal("display steps changed")
	}
	b505CandidateRemote := b505CandidateLookup(b505CandidateT, "remote_control", nil)
	if !b505CandidateReflect.DeepEqual(b505CandidateRemote.B505CandidateBranches, map[string]string{"01": "enable-associated-output"}) {
		b505CandidateT.Fatal("remote mappings changed")
	}
	for b505CandidateRaw := 0; b505CandidateRaw < 256; b505CandidateRaw++ {
		if b505CandidateRaw != 1 && b505CandidateRemote.B505CandidateSelectRaw(b505CandidateHex(byte(b505CandidateRaw))) != "Unknown" {
			b505CandidateT.Fatal("unprovided remote interpretation")
		}
	}
	b505CandidatePriority := b505CandidateLookup(b505CandidateT, "priority", b505CandidateByte(0x4E))
	if b505CandidatePriority.B505CandidateSignedness != "unsigned" || !b505CandidateReflect.DeepEqual(b505CandidatePriority.B505CandidatePriority, map[string]string{"0": "lowest", "255": "highest"}) || len(b505CandidatePriority.B505CandidateBranches) != 0 || b505CandidatePriority.B505CandidateFallback != "Unknown" {
		b505CandidateT.Fatal("priority endpoints or temperature independence changed")
	}
}

func b505CandidateHex(b505CandidateRaw byte) string {
	const b505CandidateDigits = "0123456789ABCDEF"
	return string([]byte{b505CandidateDigits[b505CandidateRaw>>4], b505CandidateDigits[b505CandidateRaw&15]})
}

func TestAb55DeepCopy(b505CandidateT *b505CandidateTesting.T) {
	b505CandidateBefore := B505Candidates()
	b505CandidateChanged := B505Candidates()
	b505CandidateChanged.B505CandidateLimitations[0] = "changed"
	for b505CandidateIndex := range b505CandidateChanged.B505CandidateEntries {
		b505CandidateEntry := &b505CandidateChanged.B505CandidateEntries[b505CandidateIndex]
		*b505CandidateEntry.B505CandidateKey.B505CandidateSecondary = 0
		if b505CandidateEntry.B505CandidateKey.B505CandidateOperation != nil {
			*b505CandidateEntry.B505CandidateKey.B505CandidateOperation = 0
		}
		if b505CandidateEntry.B505CandidateWidthBits != nil {
			*b505CandidateEntry.B505CandidateWidthBits = 0
		}
		if b505CandidateEntry.B505CandidateDurationSeconds != nil {
			*b505CandidateEntry.B505CandidateDurationSeconds = 0
		}
		b505CandidateEntry.B505CandidateBranches["8000"] = "changed"
		if b505CandidateEntry.B505CandidatePriority != nil {
			b505CandidateEntry.B505CandidatePriority["0"] = "changed"
		}
		if b505CandidateEntry.B505CandidateDisplaySteps != nil {
			b505CandidateEntry.B505CandidateDisplaySteps[0] = 99
		}
		b505CandidateEntry.B505CandidateKey.B505CandidateField = "changed"
	}
	b505CandidateLookupCopy := b505CandidateLookup(b505CandidateT, "temperature_1", b505CandidateByte(0x32))
	b505CandidateLookupCopy.B505CandidateBranches["8000"] = "changed"
	*b505CandidateLookupCopy.B505CandidateKey.B505CandidateOperation = 0
	*b505CandidateLookupCopy.B505CandidateWidthBits = 0
	if !b505CandidateReflect.DeepEqual(B505Candidates(), b505CandidateBefore) || !b505CandidateReflect.DeepEqual(b505CandidateLookup(b505CandidateT, "temperature_1", b505CandidateByte(0x32)), b505CandidateBefore.B505CandidateEntries[0]) {
		b505CandidateT.Fatal("shared catalog backing storage")
	}
}
