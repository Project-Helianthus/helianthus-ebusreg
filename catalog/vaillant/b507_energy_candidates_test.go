package vaillant

import (
	b507EnergyCandidateJSON "encoding/json"
	b507EnergyCandidateReflect "reflect"
	b507EnergyCandidateTesting "testing"
)

func b507EnergyCandidatePointer[b507EnergyCandidateT any](b507EnergyCandidateValue b507EnergyCandidateT) *b507EnergyCandidateT {
	return &b507EnergyCandidateValue
}

func b507EnergyCandidateKey(b507EnergyCandidateOperation uint8, b507EnergyCandidateField string) B507EnergyCandidateKey {
	return B507EnergyCandidateKey{B507EnergyCandidateFamily: "vaillant", B507EnergyCandidateVariant: "B507", B507EnergyCandidatePrimary: 0xb5, B507EnergyCandidateSecondary: b507EnergyCandidatePointer(uint8(7)), B507EnergyCandidateOperation: b507EnergyCandidatePointer(b507EnergyCandidateOperation), B507EnergyCandidateField: b507EnergyCandidateField}
}

func b507EnergyCandidateGet(b507EnergyCandidateT *b507EnergyCandidateTesting.T, b507EnergyCandidateOperation uint8, b507EnergyCandidateField string) B507EnergyCandidateEntry {
	b507EnergyCandidateT.Helper()
	b507EnergyCandidateEntry, b507EnergyCandidateFound, b507EnergyCandidateError := B507EnergyCandidateLookup(b507EnergyCandidateKey(b507EnergyCandidateOperation, b507EnergyCandidateField))
	if b507EnergyCandidateError != nil || !b507EnergyCandidateFound {
		b507EnergyCandidateT.Fatalf("lookup: found=%v error=%v", b507EnergyCandidateFound, b507EnergyCandidateError)
	}
	return b507EnergyCandidateEntry
}

func b507EnergyCandidateInstall(b507EnergyCandidateT *b507EnergyCandidateTesting.T, b507EnergyCandidateEntries []B507EnergyCandidateEntry) {
	b507EnergyCandidateT.Helper()
	b507EnergyCandidateSaved := b507EnergyCandidateData
	b507EnergyCandidateT.Cleanup(func() { b507EnergyCandidateData = b507EnergyCandidateSaved })
	b507EnergyCandidateEncoded, b507EnergyCandidateError := b507EnergyCandidateJSON.Marshal(b507EnergyCandidateEntries)
	if b507EnergyCandidateError != nil {
		b507EnergyCandidateT.Fatal(b507EnergyCandidateError)
	}
	b507EnergyCandidateData = b507EnergyCandidateEncoded
}

func TestAb13FiniteMetadata(b507EnergyCandidateT *b507EnergyCandidateTesting.T) {
	b507EnergyCandidateEntries, b507EnergyCandidateError := B507EnergyCandidateEntries()
	if b507EnergyCandidateError != nil || len(b507EnergyCandidateEntries) != 6 {
		b507EnergyCandidateT.Fatalf("entries=%d error=%v", len(b507EnergyCandidateEntries), b507EnergyCandidateError)
	}
	for _, b507EnergyCandidateExpected := range []struct {
		b507EnergyCandidateOperation                       uint8
		b507EnergyCandidateField, b507EnergyCandidateGrade string
	}{
		{6, "capacity", "Hypothesis"}, {6, "cost", "Hypothesis"}, {6, "temperature", "Hypothesis"},
		{8, "preference_word", "Hypothesis"}, {9, "preference_word", "Unknown"}, {11, "preference_word", "Unknown"},
	} {
		b507EnergyCandidateEntry := b507EnergyCandidateGet(b507EnergyCandidateT, b507EnergyCandidateExpected.b507EnergyCandidateOperation, b507EnergyCandidateExpected.b507EnergyCandidateField)
		if b507EnergyCandidateEntry.B507EnergyCandidateGrade != b507EnergyCandidateExpected.b507EnergyCandidateGrade || b507EnergyCandidateEntry.B507EnergyCandidateMetadataPresence != "Proven" || b507EnergyCandidateEntry.B507EnergyCandidatePhysicalQualification != "Unknown" || b507EnergyCandidateEntry.B507EnergyCandidateInvocation != "Unknown" || !b507EnergyCandidateEntry.B507EnergyCandidatePreserveRaw {
			b507EnergyCandidateT.Fatal("grade/qualification/raw boundary")
		}
		if b507EnergyCandidateEntry.B507EnergyCandidatePosition != nil || b507EnergyCandidateEntry.B507EnergyCandidateKey.B507EnergyCandidateDirection != nil || b507EnergyCandidateEntry.B507EnergyCandidateKey.B507EnergyCandidateGroup != nil || b507EnergyCandidateEntry.B507EnergyCandidateKey.B507EnergyCandidateRequestedInstance != nil || b507EnergyCandidateEntry.B507EnergyCandidateKey.B507EnergyCandidateRegister != nil || b507EnergyCandidateEntry.B507EnergyCandidateKey.B507EnergyCandidateCommandRole != nil || b507EnergyCandidateEntry.B507EnergyCandidateReference != "Unknown" {
			b507EnergyCandidateT.Fatal("invented coordinates or reference")
		}
		if b507EnergyCandidateEntry.B507EnergyCandidateNativeContract != "https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/protocols/vaillant/ebus-vaillant-b507.md" {
			b507EnergyCandidateT.Fatal("native contract")
		}
	}
}

func TestAb13Capacity(b507EnergyCandidateT *b507EnergyCandidateTesting.T) {
	b507EnergyCandidateEntry := b507EnergyCandidateGet(b507EnergyCandidateT, 6, "capacity")
	if b507EnergyCandidateEntry.B507EnergyCandidateWidthBits == nil || *b507EnergyCandidateEntry.B507EnergyCandidateWidthBits != 8 || b507EnergyCandidateEntry.B507EnergyCandidateEncoding == nil || *b507EnergyCandidateEntry.B507EnergyCandidateEncoding != "unsigned" || b507EnergyCandidateEntry.B507EnergyCandidateDivisor == nil || *b507EnergyCandidateEntry.B507EnergyCandidateDivisor != 10 || b507EnergyCandidateEntry.B507EnergyCandidateUnit == nil || *b507EnergyCandidateEntry.B507EnergyCandidateUnit != "kW" || !b507EnergyCandidateReflect.DeepEqual(b507EnergyCandidateEntry.B507EnergyCandidateRawDomain, []uint16{0, 250}) {
		b507EnergyCandidateT.Fatal("capacity candidate")
	}
	for _, b507EnergyCandidateCase := range []struct {
		b507EnergyCandidateRaw    uint16
		b507EnergyCandidateInside bool
		b507EnergyCandidateKW     float64
	}{{0, true, 0}, {1, true, 0.1}, {250, true, 25}, {251, false, 0}} {
		b507EnergyCandidateInside := b507EnergyCandidateCase.b507EnergyCandidateRaw >= b507EnergyCandidateEntry.B507EnergyCandidateRawDomain[0] && b507EnergyCandidateCase.b507EnergyCandidateRaw <= b507EnergyCandidateEntry.B507EnergyCandidateRawDomain[1]
		if b507EnergyCandidateInside != b507EnergyCandidateCase.b507EnergyCandidateInside {
			b507EnergyCandidateT.Fatal("out-of-domain must not clip")
		}
		if b507EnergyCandidateInside && float64(b507EnergyCandidateCase.b507EnergyCandidateRaw) / *b507EnergyCandidateEntry.B507EnergyCandidateDivisor != b507EnergyCandidateCase.b507EnergyCandidateKW {
			b507EnergyCandidateT.Fatal("unsigned division candidate")
		}
	}
	if !b507EnergyCandidateReflect.DeepEqual(b507EnergyCandidateEntry.B507EnergyCandidateAlternatives, []string{"Negative mode-related prose remains unresolved; no signed reinterpretation selected."}) {
		b507EnergyCandidateT.Fatal("unresolved alternative lost")
	}
}

func TestAb13OperationSentinels(b507EnergyCandidateT *b507EnergyCandidateTesting.T) {
	for _, b507EnergyCandidateOperation := range []uint8{8, 9, 11} {
		b507EnergyCandidateEntry := b507EnergyCandidateGet(b507EnergyCandidateT, b507EnergyCandidateOperation, "preference_word")
		for b507EnergyCandidateIndex, b507EnergyCandidateRaw := range []uint16{0, 65535, 1} {
			b507EnergyCandidateWant := ""
			if b507EnergyCandidateOperation == 8 {
				b507EnergyCandidateWant = []string{"appliance-unavailable candidate", "maximum-efficiency-preference", ""}[b507EnergyCandidateIndex]
			}
			b507EnergyCandidateGot, b507EnergyCandidateMapped := b507EnergyCandidateEntry.B507EnergyCandidateValues[b507EnergyCandidateRaw]
			if b507EnergyCandidateGot != b507EnergyCandidateWant || b507EnergyCandidateMapped != (b507EnergyCandidateWant != "") {
				b507EnergyCandidateT.Fatal("operation-local label", b507EnergyCandidateOperation, b507EnergyCandidateRaw)
			}
		}
		if b507EnergyCandidateOperation == 8 {
			if b507EnergyCandidateEntry.B507EnergyCandidateWidthBits == nil || *b507EnergyCandidateEntry.B507EnergyCandidateWidthBits != 16 || len(b507EnergyCandidateEntry.B507EnergyCandidateValues) != 2 {
				b507EnergyCandidateT.Fatal("OP08 word")
			}
		} else if b507EnergyCandidateEntry.B507EnergyCandidateWidthBits != nil || b507EnergyCandidateEntry.B507EnergyCandidateValues != nil {
			b507EnergyCandidateT.Fatal("transferred OP08 metadata")
		}
		if b507EnergyCandidateEntry.B507EnergyCandidateDivisor != nil || b507EnergyCandidateEntry.B507EnergyCandidateEncoding != nil || b507EnergyCandidateEntry.B507EnergyCandidateUnit != nil {
			b507EnergyCandidateT.Fatal("invented conversion")
		}
	}
}

func TestAb13CostContexts(b507EnergyCandidateT *b507EnergyCandidateTesting.T) {
	for _, b507EnergyCandidateMode := range []string{"synthetic-mode-a", "synthetic-mode-b"} {
		for _, b507EnergyCandidateReference := range []string{"monetary", "emissions"} {
			b507EnergyCandidateT.Run("ab13_"+b507EnergyCandidateMode+"_"+b507EnergyCandidateReference, func(b507EnergyCandidateT *b507EnergyCandidateTesting.T) {
				b507EnergyCandidateRaw := uint16(173)
				b507EnergyCandidateEntry := b507EnergyCandidateGet(b507EnergyCandidateT, 6, "cost")
				if !b507EnergyCandidateReflect.DeepEqual(b507EnergyCandidateEntry.B507EnergyCandidateDependencies, []string{"operating_mode", "monetary_or_emissions_reference"}) || b507EnergyCandidateEntry.B507EnergyCandidateReference != "Unknown" || b507EnergyCandidateEntry.B507EnergyCandidateQuantityOrigin != "calculated_or_estimated_possible" || b507EnergyCandidateEntry.B507EnergyCandidateDivisor != nil || b507EnergyCandidateEntry.B507EnergyCandidateWidthBits != nil || b507EnergyCandidateEntry.B507EnergyCandidateEncoding != nil || b507EnergyCandidateEntry.B507EnergyCandidateUnit != nil || b507EnergyCandidateEntry.B507EnergyCandidateValues != nil || b507EnergyCandidateEntry.B507EnergyCandidateRawDomain != nil || !b507EnergyCandidateEntry.B507EnergyCandidatePreserveRaw || b507EnergyCandidateRaw != 173 {
					b507EnergyCandidateT.Fatal("cost must remain raw with unresolved reference and converter")
				}
			})
		}
	}
	b507EnergyCandidateTemperature := b507EnergyCandidateGet(b507EnergyCandidateT, 6, "temperature")
	if b507EnergyCandidateTemperature.B507EnergyCandidateQuantityOrigin != "calculated_or_estimated_possible" || b507EnergyCandidateTemperature.B507EnergyCandidateWidthBits != nil || b507EnergyCandidateTemperature.B507EnergyCandidateDivisor != nil || b507EnergyCandidateTemperature.B507EnergyCandidateEncoding != nil || b507EnergyCandidateTemperature.B507EnergyCandidateUnit != nil || b507EnergyCandidateTemperature.B507EnergyCandidateValues != nil || b507EnergyCandidateTemperature.B507EnergyCandidateRawDomain != nil {
		b507EnergyCandidateT.Fatal("temperature is not qualified as direct measurement")
	}
}

func TestAb13ExactKeyAxes(b507EnergyCandidateT *b507EnergyCandidateTesting.T) {
	b507EnergyCandidateEntry := b507EnergyCandidateGet(b507EnergyCandidateT, 6, "capacity")
	b507EnergyCandidateEntry.B507EnergyCandidateKey = B507EnergyCandidateKey{B507EnergyCandidateFamily: "synthetic-family", B507EnergyCandidateVariant: "B524", B507EnergyCandidatePrimary: 181, B507EnergyCandidateSecondary: b507EnergyCandidatePointer(uint8(0)), B507EnergyCandidateOperation: b507EnergyCandidatePointer(uint8(0)), B507EnergyCandidateDirection: b507EnergyCandidatePointer("request"), B507EnergyCandidateField: "synthetic-field", B507EnergyCandidateGroup: b507EnergyCandidatePointer(uint8(0)), B507EnergyCandidateRequestedInstance: b507EnergyCandidatePointer(uint8(0)), B507EnergyCandidateRegister: b507EnergyCandidatePointer(uint16(0)), B507EnergyCandidateCommandRole: b507EnergyCandidatePointer("read")}
	b507EnergyCandidateInstall(b507EnergyCandidateT, []B507EnergyCandidateEntry{b507EnergyCandidateEntry})
	for b507EnergyCandidateName, b507EnergyCandidateChange := range map[string]func(*B507EnergyCandidateKey){
		"ab13_family": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateFamily = "other-family"
		},
		"ab13_variant": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateVariant = "B507"
		},
		"ab13_primary": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidatePrimary = 182
		},
		"ab13_secondary": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateSecondary = b507EnergyCandidatePointer(uint8(1))
		},
		"ab13_secondary_absent": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateSecondary = nil
		},
		"ab13_operation": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateOperation = b507EnergyCandidatePointer(uint8(1))
		},
		"ab13_operation_absent": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateOperation = nil
		},
		"ab13_direction": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateDirection = b507EnergyCandidatePointer("response")
		},
		"ab13_direction_absent": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateDirection = nil
		},
		"ab13_field": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateField = "other-field"
		},
		"ab13_group": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateGroup = b507EnergyCandidatePointer(uint8(1))
		},
		"ab13_group_absent": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateGroup = nil
		},
		"ab13_requested_instance": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateRequestedInstance = b507EnergyCandidatePointer(uint8(1))
		},
		"ab13_requested_instance_absent": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateRequestedInstance = nil
		},
		"ab13_register": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateRegister = b507EnergyCandidatePointer(uint16(1))
		},
		"ab13_register_high_byte": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateRegister = b507EnergyCandidatePointer(uint16(256))
		},
		"ab13_register_absent": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateRegister = nil
		},
		"ab13_command_role": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateCommandRole = b507EnergyCandidatePointer("write")
		},
		"ab13_command_role_absent": func(b507EnergyCandidateKey *B507EnergyCandidateKey) {
			b507EnergyCandidateKey.B507EnergyCandidateCommandRole = nil
		},
	} {
		b507EnergyCandidateT.Run(b507EnergyCandidateName, func(b507EnergyCandidateT *b507EnergyCandidateTesting.T) {
			b507EnergyCandidateControl, b507EnergyCandidateFound, b507EnergyCandidateError := B507EnergyCandidateLookup(b507EnergyCandidateEntry.B507EnergyCandidateKey)
			if b507EnergyCandidateError != nil || !b507EnergyCandidateFound || !b507EnergyCandidateReflect.DeepEqual(b507EnergyCandidateControl, b507EnergyCandidateEntry) {
				b507EnergyCandidateT.Fatal("positive control", b507EnergyCandidateError)
			}
			b507EnergyCandidateChanged := b507EnergyCandidateEntry.B507EnergyCandidateKey
			b507EnergyCandidateChange(&b507EnergyCandidateChanged)
			b507EnergyCandidateGot, b507EnergyCandidateFound, b507EnergyCandidateError := B507EnergyCandidateLookup(b507EnergyCandidateChanged)
			if b507EnergyCandidateError != nil || b507EnergyCandidateFound || !b507EnergyCandidateReflect.DeepEqual(b507EnergyCandidateGot, B507EnergyCandidateEntry{}) {
				b507EnergyCandidateT.Fatal("axis omitted or wildcarded", b507EnergyCandidateError)
			}
		})
	}
}

func TestAb13FreshOwnership(b507EnergyCandidateT *b507EnergyCandidateTesting.T) {
	b507EnergyCandidateBaseline, b507EnergyCandidateError := B507EnergyCandidateEntries()
	if b507EnergyCandidateError != nil {
		b507EnergyCandidateT.Fatal(b507EnergyCandidateError)
	}
	b507EnergyCandidateEntries, b507EnergyCandidateError := B507EnergyCandidateEntries()
	if b507EnergyCandidateError != nil {
		b507EnergyCandidateT.Fatal(b507EnergyCandidateError)
	}
	*b507EnergyCandidateEntries[0].B507EnergyCandidateKey.B507EnergyCandidateSecondary = 0
	*b507EnergyCandidateEntries[0].B507EnergyCandidateKey.B507EnergyCandidateOperation = 0
	*b507EnergyCandidateEntries[0].B507EnergyCandidateWidthBits = 0
	*b507EnergyCandidateEntries[0].B507EnergyCandidateEncoding = "signed8"
	b507EnergyCandidateEntries[0].B507EnergyCandidateRawDomain[1] = 255
	*b507EnergyCandidateEntries[0].B507EnergyCandidateDivisor = 0.1
	*b507EnergyCandidateEntries[0].B507EnergyCandidateUnit = "other"
	b507EnergyCandidateEntries[0].B507EnergyCandidateUnknowns[0] = "other"
	b507EnergyCandidateEntries[0].B507EnergyCandidateAlternatives[0] = "other"
	b507EnergyCandidateEntries[1].B507EnergyCandidateDependencies[0] = "other"
	b507EnergyCandidateEntries[3].B507EnergyCandidateValues[0] = "other"
	b507EnergyCandidateLookup := b507EnergyCandidateGet(b507EnergyCandidateT, 8, "preference_word")
	b507EnergyCandidateLookup.B507EnergyCandidateValues[0] = "other"
	*b507EnergyCandidateLookup.B507EnergyCandidateKey.B507EnergyCandidateOperation = 0
	b507EnergyCandidateLookup.B507EnergyCandidateUnknowns[0] = "other"
	b507EnergyCandidateAgain, b507EnergyCandidateError := B507EnergyCandidateEntries()
	if b507EnergyCandidateError != nil || !b507EnergyCandidateReflect.DeepEqual(b507EnergyCandidateAgain, b507EnergyCandidateBaseline) || !b507EnergyCandidateReflect.DeepEqual(b507EnergyCandidateGet(b507EnergyCandidateT, 8, "preference_word"), b507EnergyCandidateBaseline[3]) {
		b507EnergyCandidateT.Fatal("shared storage", b507EnergyCandidateError)
	}
}

func TestAb13OptionalPointerOwnership(b507EnergyCandidateT *b507EnergyCandidateTesting.T) {
	b507EnergyCandidateEntry := b507EnergyCandidateGet(b507EnergyCandidateT, 6, "capacity")
	b507EnergyCandidateEntry.B507EnergyCandidateKey.B507EnergyCandidateDirection = b507EnergyCandidatePointer("synthetic-direction")
	b507EnergyCandidateEntry.B507EnergyCandidateKey.B507EnergyCandidateGroup = b507EnergyCandidatePointer(uint8(0))
	b507EnergyCandidateEntry.B507EnergyCandidateKey.B507EnergyCandidateRequestedInstance = b507EnergyCandidatePointer(uint8(0))
	b507EnergyCandidateEntry.B507EnergyCandidateKey.B507EnergyCandidateRegister = b507EnergyCandidatePointer(uint16(0))
	b507EnergyCandidateEntry.B507EnergyCandidateKey.B507EnergyCandidateCommandRole = b507EnergyCandidatePointer("synthetic-role")
	b507EnergyCandidateEntry.B507EnergyCandidatePosition = b507EnergyCandidatePointer(uint16(0))
	b507EnergyCandidateInstall(b507EnergyCandidateT, []B507EnergyCandidateEntry{b507EnergyCandidateEntry})
	b507EnergyCandidateEntries, b507EnergyCandidateError := B507EnergyCandidateEntries()
	if b507EnergyCandidateError != nil {
		b507EnergyCandidateT.Fatal(b507EnergyCandidateError)
	}
	b507EnergyCandidateLookup, b507EnergyCandidateFound, b507EnergyCandidateError := B507EnergyCandidateLookup(b507EnergyCandidateEntry.B507EnergyCandidateKey)
	if b507EnergyCandidateError != nil || !b507EnergyCandidateFound {
		b507EnergyCandidateT.Fatal("positive control", b507EnergyCandidateError)
	}
	for _, b507EnergyCandidateChanged := range []*B507EnergyCandidateEntry{&b507EnergyCandidateEntries[0], &b507EnergyCandidateLookup} {
		*b507EnergyCandidateChanged.B507EnergyCandidateKey.B507EnergyCandidateDirection = "changed"
		*b507EnergyCandidateChanged.B507EnergyCandidateKey.B507EnergyCandidateGroup = 1
		*b507EnergyCandidateChanged.B507EnergyCandidateKey.B507EnergyCandidateRequestedInstance = 1
		*b507EnergyCandidateChanged.B507EnergyCandidateKey.B507EnergyCandidateRegister = 1
		*b507EnergyCandidateChanged.B507EnergyCandidateKey.B507EnergyCandidateCommandRole = "changed"
		*b507EnergyCandidateChanged.B507EnergyCandidatePosition = 1
	}
	b507EnergyCandidateAgain, b507EnergyCandidateFound, b507EnergyCandidateError := B507EnergyCandidateLookup(b507EnergyCandidateEntry.B507EnergyCandidateKey)
	if b507EnergyCandidateError != nil || !b507EnergyCandidateFound || !b507EnergyCandidateReflect.DeepEqual(b507EnergyCandidateAgain, b507EnergyCandidateEntry) {
		b507EnergyCandidateT.Fatal("optional pointers share storage", b507EnergyCandidateError)
	}
}

func TestAb13MalformedData(b507EnergyCandidateT *b507EnergyCandidateTesting.T) {
	b507EnergyCandidateSaved := b507EnergyCandidateData
	b507EnergyCandidateT.Cleanup(func() { b507EnergyCandidateData = b507EnergyCandidateSaved })
	for _, b507EnergyCandidateInvalid := range []string{"{", "null", "[]", "[{}]", string(b507EnergyCandidateSaved) + " true", `[{"unexpected":1}]`} {
		b507EnergyCandidateData = []byte(b507EnergyCandidateInvalid)
		b507EnergyCandidateEntries, b507EnergyCandidateError := B507EnergyCandidateEntries()
		if b507EnergyCandidateError == nil || b507EnergyCandidateEntries != nil {
			b507EnergyCandidateT.Fatal("malformed entries succeeded")
		}
		b507EnergyCandidateEntry, b507EnergyCandidateFound, b507EnergyCandidateError := B507EnergyCandidateLookup(b507EnergyCandidateKey(6, "capacity"))
		if b507EnergyCandidateError == nil || b507EnergyCandidateFound || !b507EnergyCandidateReflect.DeepEqual(b507EnergyCandidateEntry, B507EnergyCandidateEntry{}) {
			b507EnergyCandidateT.Fatal("malformed lookup succeeded")
		}
	}
	b507EnergyCandidateData = b507EnergyCandidateSaved
	b507EnergyCandidateEntries, b507EnergyCandidateError := B507EnergyCandidateEntries()
	if b507EnergyCandidateError != nil {
		b507EnergyCandidateT.Fatal(b507EnergyCandidateError)
	}
	b507EnergyCandidateEntries = append(b507EnergyCandidateEntries, b507EnergyCandidateEntries[0])
	b507EnergyCandidateInstall(b507EnergyCandidateT, b507EnergyCandidateEntries)
	if _, b507EnergyCandidateError := B507EnergyCandidateEntries(); b507EnergyCandidateError == nil {
		b507EnergyCandidateT.Fatal("duplicate key accepted")
	}
}
