package vaillant

import (
	b511ClassCandidateBytes "bytes"
	b511ClassCandidateBinary "encoding/binary"
	b511ClassCandidateJSON "encoding/json"
	b511ClassCandidateReflect "reflect"
	b511ClassCandidateTesting "testing"
)

func b511ClassCandidatePtr(b511ClassCandidateValue uint16) *uint16 { return &b511ClassCandidateValue }

func b511ClassCandidateTestKey(b511ClassCandidateOperation uint16, b511ClassCandidateField string) B511ClassCandidateKey {
	return B511ClassCandidateKey{B511ClassCandidateFamily: "Vaillant", B511ClassCandidateVariant: "class-field-candidate", B511ClassCandidatePrimary: 181, B511ClassCandidateSecondary: b511ClassCandidatePtr(17), B511ClassCandidateOperation: b511ClassCandidatePtr(b511ClassCandidateOperation), B511ClassCandidateDirection: "Unknown", B511ClassCandidateField: b511ClassCandidateField, B511ClassCandidateAlternative: "none", B511ClassCandidateInstanceScope: "none"}
}

func b511ClassCandidateTestRow(b511ClassCandidateT *b511ClassCandidateTesting.T, b511ClassCandidateKey B511ClassCandidateKey) map[string]b511ClassCandidateJSON.RawMessage {
	b511ClassCandidateT.Helper()
	b511ClassCandidateRaw, b511ClassCandidateFound, b511ClassCandidateError := B511ClassCandidateLookup(b511ClassCandidateKey)
	if b511ClassCandidateError != nil || !b511ClassCandidateFound {
		b511ClassCandidateT.Fatalf("lookup: %v %v", b511ClassCandidateFound, b511ClassCandidateError)
	}
	var b511ClassCandidateRow map[string]b511ClassCandidateJSON.RawMessage
	if b511ClassCandidateError = b511ClassCandidateJSON.Unmarshal(b511ClassCandidateRaw, &b511ClassCandidateRow); b511ClassCandidateError != nil {
		b511ClassCandidateT.Fatal(b511ClassCandidateError)
	}
	return b511ClassCandidateRow
}

func b511ClassCandidateAssert(b511ClassCandidateT *b511ClassCandidateTesting.T, b511ClassCandidateRow map[string]b511ClassCandidateJSON.RawMessage, b511ClassCandidateName, b511ClassCandidateExpected string) {
	b511ClassCandidateT.Helper()
	var b511ClassCandidateActualValue, b511ClassCandidateExpectedValue any
	if b511ClassCandidateJSON.Unmarshal(b511ClassCandidateRow[b511ClassCandidateName], &b511ClassCandidateActualValue) != nil || b511ClassCandidateJSON.Unmarshal([]byte(b511ClassCandidateExpected), &b511ClassCandidateExpectedValue) != nil || !b511ClassCandidateReflect.DeepEqual(b511ClassCandidateActualValue, b511ClassCandidateExpectedValue) {
		b511ClassCandidateT.Fatalf("%s: got %s want %s", b511ClassCandidateName, b511ClassCandidateRow[b511ClassCandidateName], b511ClassCandidateExpected)
	}
}

func TestAb50FiniteDeclarations(b511ClassCandidateT *b511ClassCandidateTesting.T) {
	b511ClassCandidateExpected := [][]string{
		{"flags-a", "flags-b", "flags-c", "modulation"},
		{"handshakes", "forced-mode-guards"},
		{"status", "DHW-enabled", "tapping-enabled"},
		{"basic-state", "flags-a", "flags-b", "message"},
		{"status", "temperature-a", "temperature-b"},
		{"setpoint-a", "setpoint-b", "setpoint-c", "setpoint-d", "status-a", "status-b"},
		{"CH-limit", "DHW-limit", "pressure-limit", "cooling-limit", "status"},
		{"power", "energy", "pump-flags", "status-flags"},
		{"volume", "pump", "status"},
		{"yield", "power", "activity", "display", "day-night"},
		{"error", "loading"},
		{"fan", "radiator", "bypass", "error", "maintenance", "assistant"},
		{"operation", "switch-on-flags"},
		{"status", "configuration"},
	}
	b511ClassCandidateEntriesJSON, b511ClassCandidateError := B511ClassCandidateEntries()
	if b511ClassCandidateError != nil {
		b511ClassCandidateT.Fatal(b511ClassCandidateError)
	}
	var b511ClassCandidateRows []b511ClassCandidateJSON.RawMessage
	if b511ClassCandidateJSON.Unmarshal(b511ClassCandidateEntriesJSON, &b511ClassCandidateRows) != nil || len(b511ClassCandidateRows) != 53 {
		b511ClassCandidateT.Fatal("finite array must contain 53 declarations")
	}
	for b511ClassCandidateOperation, b511ClassCandidateFields := range b511ClassCandidateExpected {
		for _, b511ClassCandidateField := range b511ClassCandidateFields {
			b511ClassCandidateRow := b511ClassCandidateTestRow(b511ClassCandidateT, b511ClassCandidateTestKey(uint16(b511ClassCandidateOperation), b511ClassCandidateField))
			b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "grade", `"Hypothesis"`)
			b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "physical_qualification", `"Unknown"`)
			b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "offset", "null")
			b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "default", "null")
			b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "bit", "null")
			b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "labels", "{}")
			b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "raw_form", `"preserve original bytes and unlabelled values without clipping or sentinel inference"`)
			b511ClassCandidateWidth, b511ClassCandidateBounds, b511ClassCandidateUnit, b511ClassCandidateOrder := "null", "null", "null", `"Unknown"`
			if b511ClassCandidateField == "energy" || b511ClassCandidateField == "yield" {
				b511ClassCandidateWidth, b511ClassCandidateUnit, b511ClassCandidateOrder = "16", `"kWh"`, `"little-endian"`
			}
			if (b511ClassCandidateOperation == 0 || b511ClassCandidateOperation == 3) && (b511ClassCandidateField == "flags-a" || b511ClassCandidateField == "flags-b" || b511ClassCandidateField == "flags-c") || b511ClassCandidateOperation == 5 && (b511ClassCandidateField == "status-a" || b511ClassCandidateField == "status-b") {
				b511ClassCandidateWidth = "8"
			}
			if (b511ClassCandidateOperation == 7 || b511ClassCandidateOperation == 9) && b511ClassCandidateField == "power" || b511ClassCandidateOperation == 8 && b511ClassCandidateField == "volume" {
				b511ClassCandidateUnit = `"percent"`
			}
			switch {
			case b511ClassCandidateOperation == 2 && b511ClassCandidateField == "status":
				b511ClassCandidateBounds = "[0,1]"
			case b511ClassCandidateOperation == 3 && b511ClassCandidateField == "basic-state":
				b511ClassCandidateBounds = "[0,6]"
			case b511ClassCandidateOperation == 3 && b511ClassCandidateField == "message":
				b511ClassCandidateBounds = "[0,3]"
			case b511ClassCandidateOperation == 7 && b511ClassCandidateField == "energy":
				b511ClassCandidateBounds = "[0,9999]"
			case b511ClassCandidateOperation == 9 && b511ClassCandidateField == "yield":
				b511ClassCandidateBounds = "[0,65535]"
			case b511ClassCandidateOperation == 12 && b511ClassCandidateField == "operation":
				b511ClassCandidateBounds = "[0,3]"
			}
			b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "width_bits", b511ClassCandidateWidth)
			b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "bounds", b511ClassCandidateBounds)
			b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "unit", b511ClassCandidateUnit)
			b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "byte_order", b511ClassCandidateOrder)
			b511ClassCandidateUnresolved := `["absolute position","per-variant defaults","individual labels"]`
			if b511ClassCandidateWidth == "null" {
				b511ClassCandidateUnresolved = `["absolute position","per-variant defaults","individual labels","width"]`
			}
			b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "unresolved", b511ClassCandidateUnresolved)
		}
	}
	b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateTestRow(b511ClassCandidateT, b511ClassCandidateTestKey(6, "status")), "unlabelled_values", "[4,9,16]")
	b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateTestRow(b511ClassCandidateT, b511ClassCandidateTestKey(2, "DHW-enabled")), "note", `"DHW-enabled covers tank-loading/warm-start while tapping remains enabled."`)
}

func TestAb50Applicability(b511ClassCandidateT *b511ClassCandidateTesting.T) {
	for _, b511ClassCandidateCase := range []struct {
		B511ClassCandidateOperation                                               uint16
		B511ClassCandidateField, B511ClassCandidateVariant, B511ClassCandidateBit string
	}{{0, "second-modulation", "second-modulation-applicable", "null"}, {1, "warm-start", "warm-start-applicable", "5"}} {
		b511ClassCandidateKey := b511ClassCandidateTestKey(b511ClassCandidateCase.B511ClassCandidateOperation, b511ClassCandidateCase.B511ClassCandidateField)
		b511ClassCandidateKey.B511ClassCandidateVariant = b511ClassCandidateCase.B511ClassCandidateVariant
		b511ClassCandidateRow := b511ClassCandidateTestRow(b511ClassCandidateT, b511ClassCandidateKey)
		b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "bit", b511ClassCandidateCase.B511ClassCandidateBit)
		for _, b511ClassCandidateName := range []string{"offset", "width_bits", "default", "bounds"} {
			b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, b511ClassCandidateName, "null")
		}
		b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "grade", `"Hypothesis"`)
		b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "physical_qualification", `"Unknown"`)
		b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "labels", `{}`)
		b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "raw_form", `"preserve original bytes and unlabelled values without clipping or sentinel inference"`)
		for _, b511ClassCandidateVariant := range []string{"class-field-candidate", "warm-start-inapplicable", "none"} {
			b511ClassCandidateKey.B511ClassCandidateVariant = b511ClassCandidateVariant
			if _, b511ClassCandidateFound, b511ClassCandidateError := B511ClassCandidateLookup(b511ClassCandidateKey); b511ClassCandidateError != nil || b511ClassCandidateFound {
				b511ClassCandidateT.Fatal("applicability widened", b511ClassCandidateVariant, b511ClassCandidateError)
			}
		}
	}
}

func TestAb50SyntheticRawDomains(b511ClassCandidateT *b511ClassCandidateTesting.T) {
	for _, b511ClassCandidateCase := range []struct {
		B511ClassCandidateOperation uint16
		B511ClassCandidateField     string
		B511ClassCandidateRaw       []byte
		B511ClassCandidateValue     uint16
		B511ClassCandidateWithin    bool
	}{
		{7, "energy", []byte{0x0f, 0x27}, 9999, true}, {7, "energy", []byte{0x10, 0x27}, 10000, false}, {9, "yield", []byte{0xff, 0xff}, 65535, true},
		{3, "message", []byte{3, 0}, 3, true}, {3, "message", []byte{4, 0}, 4, false}, {3, "basic-state", []byte{6, 0}, 6, true}, {3, "basic-state", []byte{7, 0}, 7, false},
	} {
		b511ClassCandidateRow := b511ClassCandidateTestRow(b511ClassCandidateT, b511ClassCandidateTestKey(b511ClassCandidateCase.B511ClassCandidateOperation, b511ClassCandidateCase.B511ClassCandidateField))
		var b511ClassCandidateBounds []uint16
		if b511ClassCandidateJSON.Unmarshal(b511ClassCandidateRow["bounds"], &b511ClassCandidateBounds) != nil || len(b511ClassCandidateBounds) != 2 {
			b511ClassCandidateT.Fatal("bounds")
		}
		b511ClassCandidateValue := b511ClassCandidateBinary.LittleEndian.Uint16(b511ClassCandidateCase.B511ClassCandidateRaw)
		if b511ClassCandidateValue != b511ClassCandidateCase.B511ClassCandidateValue || (b511ClassCandidateValue >= b511ClassCandidateBounds[0] && b511ClassCandidateValue <= b511ClassCandidateBounds[1]) != b511ClassCandidateCase.B511ClassCandidateWithin {
			b511ClassCandidateT.Fatal("independent raw domain result")
		}
		b511ClassCandidateAssert(b511ClassCandidateT, b511ClassCandidateRow, "raw_form", `"preserve original bytes and unlabelled values without clipping or sentinel inference"`)
	}
}

func TestAb50EveryIdentityAxis(b511ClassCandidateT *b511ClassCandidateTesting.T) {
	b511ClassCandidateBase := b511ClassCandidateTestKey(7, "energy")
	b511ClassCandidateChanges := map[string]func(*B511ClassCandidateKey){
		"ab50Family": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidateFamily = "ab50Other"
		},
		"ab50Variant": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidateVariant = "ab50Other"
		},
		"ab50Primary": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidatePrimary = 0
		},
		"ab50SecondaryAbsent": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidateSecondary = nil
		},
		"ab50SecondaryZero": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidateSecondary = b511ClassCandidatePtr(0)
		},
		"ab50OperationAbsent": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidateOperation = nil
		},
		"ab50OperationZero": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidateOperation = b511ClassCandidatePtr(0)
		},
		"ab50Direction": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidateDirection = "none"
		},
		"ab50Field": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidateField = "yield"
		},
		"ab50Alternative": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidateAlternative = "ab50Other"
		},
		"ab50GroupZero": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidateGroup = b511ClassCandidatePtr(0)
		},
		"ab50GroupOther": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidateGroup = b511ClassCandidatePtr(1)
		},
		"ab50InstanceScope": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidateInstanceScope = "ab50Other"
		},
		"ab50RegisterZero": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidateRegister = b511ClassCandidatePtr(0)
		},
		"ab50RegisterOther": func(b511ClassCandidateKey *B511ClassCandidateKey) {
			b511ClassCandidateKey.B511ClassCandidateRegister = b511ClassCandidatePtr(1)
		},
	}
	for b511ClassCandidateName, b511ClassCandidateChange := range b511ClassCandidateChanges {
		b511ClassCandidateT.Run(b511ClassCandidateName, func(b511ClassCandidateT *b511ClassCandidateTesting.T) {
			b511ClassCandidateKey := b511ClassCandidateBase
			b511ClassCandidateChange(&b511ClassCandidateKey)
			if _, b511ClassCandidateFound, b511ClassCandidateError := B511ClassCandidateLookup(b511ClassCandidateKey); b511ClassCandidateFound || b511ClassCandidateError != nil {
				b511ClassCandidateT.Fatal("inexact production match", b511ClassCandidateError)
			}
			b511ClassCandidateControl := []struct {
				B511ClassCandidateKey    B511ClassCandidateKey `json:"key"`
				B511ClassCandidateMarker string                `json:"ab50_marker"`
			}{{b511ClassCandidateBase, "ab50Base"}, {b511ClassCandidateKey, "ab50Changed"}}
			b511ClassCandidateSource, b511ClassCandidateError := b511ClassCandidateJSON.Marshal(b511ClassCandidateControl)
			if b511ClassCandidateError != nil {
				b511ClassCandidateT.Fatal(b511ClassCandidateError)
			}
			for b511ClassCandidateIndex, b511ClassCandidateWanted := range []B511ClassCandidateKey{b511ClassCandidateBase, b511ClassCandidateKey} {
				b511ClassCandidateRow, b511ClassCandidateFound, b511ClassCandidateError := b511ClassCandidateLookup(b511ClassCandidateSource, b511ClassCandidateWanted)
				var b511ClassCandidateActual map[string]b511ClassCandidateJSON.RawMessage
				if b511ClassCandidateError != nil || !b511ClassCandidateFound || b511ClassCandidateJSON.Unmarshal(b511ClassCandidateRow, &b511ClassCandidateActual) != nil {
					b511ClassCandidateT.Fatal("valid identity control", b511ClassCandidateError)
				}
				b511ClassCandidateExpected, _ := b511ClassCandidateJSON.Marshal(b511ClassCandidateControl[b511ClassCandidateIndex].B511ClassCandidateMarker)
				if !b511ClassCandidateBytes.Equal(b511ClassCandidateActual["ab50_marker"], b511ClassCandidateExpected) {
					b511ClassCandidateT.Fatal("axis collapsed")
				}
			}
		})
	}
}

func TestAb50ValidationAndStorage(b511ClassCandidateT *b511ClassCandidateTesting.T) {
	b511ClassCandidateKey := b511ClassCandidateTestKey(7, "energy")
	b511ClassCandidateOriginal, b511ClassCandidateFound, b511ClassCandidateError := B511ClassCandidateLookup(b511ClassCandidateKey)
	if b511ClassCandidateError != nil || !b511ClassCandidateFound {
		b511ClassCandidateT.Fatal(b511ClassCandidateError)
	}
	b511ClassCandidateEntriesJSON, b511ClassCandidateError := B511ClassCandidateEntries()
	if b511ClassCandidateError != nil {
		b511ClassCandidateT.Fatal(b511ClassCandidateError)
	}
	b511ClassCandidateSource := append([]byte(nil), b511ClassCandidateEntriesJSON...)
	b511ClassCandidateCopied, b511ClassCandidateError := b511ClassCandidateEntries(b511ClassCandidateSource)
	if b511ClassCandidateError != nil {
		b511ClassCandidateT.Fatal(b511ClassCandidateError)
	}
	b511ClassCandidateRow, _, b511ClassCandidateError := b511ClassCandidateLookup(b511ClassCandidateSource, b511ClassCandidateKey)
	if b511ClassCandidateError != nil {
		b511ClassCandidateT.Fatal(b511ClassCandidateError)
	}
	for b511ClassCandidateIndex := range b511ClassCandidateSource {
		b511ClassCandidateSource[b511ClassCandidateIndex] = 'x'
	}
	if !b511ClassCandidateBytes.Equal(b511ClassCandidateCopied, b511ClassCandidateEntriesJSON) || !b511ClassCandidateBytes.Equal(b511ClassCandidateRow, b511ClassCandidateOriginal) {
		b511ClassCandidateT.Fatal("source aliases returned storage")
	}
	b511ClassCandidateOriginal[0] = 'x'
	b511ClassCandidateEntriesJSON[0] = 'x'
	b511ClassCandidateAgain, _, b511ClassCandidateError := B511ClassCandidateLookup(b511ClassCandidateKey)
	if b511ClassCandidateError != nil || !b511ClassCandidateBytes.Equal(b511ClassCandidateAgain, b511ClassCandidateRow) {
		b511ClassCandidateT.Fatal("lookup storage leaked")
	}
	b511ClassCandidateAgainEntries, b511ClassCandidateError := B511ClassCandidateEntries()
	if b511ClassCandidateError != nil || !b511ClassCandidateBytes.Equal(b511ClassCandidateAgainEntries, b511ClassCandidateCopied) {
		b511ClassCandidateT.Fatal("entries storage leaked")
	}
	b511ClassCandidateDuplicate := append([]byte{'['}, b511ClassCandidateAgain...)
	b511ClassCandidateDuplicate = append(b511ClassCandidateDuplicate, ',')
	b511ClassCandidateDuplicate = append(b511ClassCandidateDuplicate, b511ClassCandidateAgain...)
	b511ClassCandidateDuplicate = append(b511ClassCandidateDuplicate, ']')
	for _, b511ClassCandidateInvalid := range [][]byte{[]byte(`{`), []byte(`null`), []byte(`{}`), []byte(`[null]`), []byte(`[{}]`), []byte(`[{"key":{"family":false}}]`), b511ClassCandidateDuplicate, append(append([]byte{'['}, b511ClassCandidateAgain...), []byte(`,{}]`)...)} {
		if b511ClassCandidateResult, b511ClassCandidateError := b511ClassCandidateEntries(b511ClassCandidateInvalid); b511ClassCandidateError == nil || b511ClassCandidateResult != nil {
			b511ClassCandidateT.Fatal("invalid entries succeeded")
		}
		if b511ClassCandidateResult, b511ClassCandidateFound, b511ClassCandidateError := b511ClassCandidateLookup(b511ClassCandidateInvalid, b511ClassCandidateKey); b511ClassCandidateError == nil || b511ClassCandidateFound || b511ClassCandidateResult != nil {
			b511ClassCandidateT.Fatal("invalid lookup succeeded")
		}
	}
}

func TestAb50WholeRowMetadataIsolation(b511ClassCandidateT *b511ClassCandidateTesting.T) {
	b511ClassCandidateKey := b511ClassCandidateTestKey(7, "energy")
	b511ClassCandidateRaw, b511ClassCandidateFound, b511ClassCandidateError := B511ClassCandidateLookup(b511ClassCandidateKey)
	if b511ClassCandidateError != nil || !b511ClassCandidateFound {
		b511ClassCandidateT.Fatal(b511ClassCandidateError)
	}
	for b511ClassCandidateName, b511ClassCandidateValue := range map[string]string{
		"grade": `"Unknown"`, "width_bits": `8`,
		"unresolved": `["ab50 unresolved fixture"]`,
		"raw_form":   `{"ab50_raw":[16,39],"ab50_value":10000}`,
	} {
		var b511ClassCandidateObject map[string]b511ClassCandidateJSON.RawMessage
		if b511ClassCandidateJSON.Unmarshal(b511ClassCandidateRaw, &b511ClassCandidateObject) != nil {
			b511ClassCandidateT.Fatal("fixture decode")
		}
		b511ClassCandidateObject[b511ClassCandidateName] = b511ClassCandidateJSON.RawMessage(b511ClassCandidateValue)
		b511ClassCandidateChanged, b511ClassCandidateError := b511ClassCandidateJSON.Marshal(b511ClassCandidateObject)
		if b511ClassCandidateError != nil {
			b511ClassCandidateT.Fatal(b511ClassCandidateError)
		}
		b511ClassCandidateSource := append(append([]byte{'['}, b511ClassCandidateChanged...), ']')
		b511ClassCandidateResult, b511ClassCandidateFound, b511ClassCandidateError := b511ClassCandidateLookup(b511ClassCandidateSource, b511ClassCandidateKey)
		if b511ClassCandidateError != nil || !b511ClassCandidateFound || !b511ClassCandidateBytes.Equal(b511ClassCandidateResult, b511ClassCandidateChanged) {
			b511ClassCandidateT.Fatal("whole row was rewritten", b511ClassCandidateName, b511ClassCandidateError)
		}
		b511ClassCandidateLater, _, b511ClassCandidateError := B511ClassCandidateLookup(b511ClassCandidateKey)
		if b511ClassCandidateError != nil || !b511ClassCandidateBytes.Equal(b511ClassCandidateLater, b511ClassCandidateRaw) {
			b511ClassCandidateT.Fatal("fixture contaminated production")
		}
	}
	for _, b511ClassCandidatePrimary := range []string{"null", `"181"`} {
		var b511ClassCandidateObject map[string]b511ClassCandidateJSON.RawMessage
		var b511ClassCandidateAxes map[string]b511ClassCandidateJSON.RawMessage
		if b511ClassCandidateJSON.Unmarshal(b511ClassCandidateRaw, &b511ClassCandidateObject) != nil || b511ClassCandidateJSON.Unmarshal(b511ClassCandidateObject["key"], &b511ClassCandidateAxes) != nil {
			b511ClassCandidateT.Fatal("fixture decode")
		}
		b511ClassCandidateAxes["primary"] = b511ClassCandidateJSON.RawMessage(b511ClassCandidatePrimary)
		b511ClassCandidateObject["key"], b511ClassCandidateError = b511ClassCandidateJSON.Marshal(b511ClassCandidateAxes)
		if b511ClassCandidateError != nil {
			b511ClassCandidateT.Fatal(b511ClassCandidateError)
		}
		b511ClassCandidateChanged, b511ClassCandidateError := b511ClassCandidateJSON.Marshal(b511ClassCandidateObject)
		if b511ClassCandidateError != nil {
			b511ClassCandidateT.Fatal(b511ClassCandidateError)
		}
		b511ClassCandidateSource := append(append([]byte{'['}, b511ClassCandidateChanged...), ']')
		if _, b511ClassCandidateFound, b511ClassCandidateError := b511ClassCandidateLookup(b511ClassCandidateSource, b511ClassCandidateKey); b511ClassCandidateError == nil || b511ClassCandidateFound {
			b511ClassCandidateT.Fatal("malformed primary succeeded")
		}
	}
}
