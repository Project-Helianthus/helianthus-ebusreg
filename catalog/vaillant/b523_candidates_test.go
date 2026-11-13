package vaillant

import (
	b523CandidateFmt "fmt"
	b523CandidateReflect "reflect"
	b523CandidateTesting "testing"
)

func b523CandidatePtr[b523CandidateT any](b523CandidateValue b523CandidateT) *b523CandidateT {
	return &b523CandidateValue
}

func b523CandidateKey(b523CandidateOperation byte, b523CandidateField string) B523CandidateKey {
	b523CandidateDirection := "request"
	if b523CandidateField == "feedback" {
		b523CandidateDirection = "response"
	}
	return B523CandidateKey{B523CandidateFamily: "vaillant", B523CandidatePrimary: 0xb5, B523CandidateSecondary: b523CandidatePtr(byte(0x23)), B523CandidateOperation: b523CandidatePtr(b523CandidateOperation), B523CandidateDirection: b523CandidateDirection, B523CandidateField: b523CandidateField}
}

func b523CandidateLookup(b523CandidateT *b523CandidateTesting.T, b523CandidateKeyValue B523CandidateKey) B523CandidateEntry {
	b523CandidateT.Helper()
	b523CandidateEntry, b523CandidateFound, b523CandidateError := B523CandidateLookup(b523CandidateKeyValue)
	if b523CandidateError != nil || !b523CandidateFound {
		b523CandidateT.Fatalf("lookup: found=%v error=%v key=%+v", b523CandidateFound, b523CandidateError, b523CandidateKeyValue)
	}
	return b523CandidateEntry
}

func b523CandidateRule(b523CandidateMin, b523CandidateMax byte, b523CandidateResult string, b523CandidateDivisor *float64, b523CandidateDuration *string) B523CandidateRule {
	return B523CandidateRule{B523CandidateRawMin: b523CandidateMin, B523CandidateRawMax: b523CandidateMax, B523CandidateResult: b523CandidateResult, B523CandidateDivideBy: b523CandidateDivisor, B523CandidateDurationUnit: b523CandidateDuration, B523CandidateGrade: "Hypothesis"}
}

func b523CandidateExpected(b523CandidateOperation byte, b523CandidateField string, b523CandidatePosition int, b523CandidateRequest *int, b523CandidateAlternatives []B523CandidateAlternative, b523CandidateDefault *string) B523CandidateEntry {
	return B523CandidateEntry{B523CandidateKey: b523CandidateKey(b523CandidateOperation, b523CandidateField), B523CandidateLogicalPosition: b523CandidatePtr(b523CandidatePosition), B523CandidateWidthBits: b523CandidatePtr(8), B523CandidateRequestExtent: b523CandidateRequest, B523CandidatePositiveExtent: b523CandidatePtr(4), B523CandidateDefaultRaw: b523CandidatePtr(byte(0)), B523CandidateDefaultResult: b523CandidateDefault, B523CandidateGrade: "Hypothesis", B523CandidateRawDomain: []byte{0, 255}, B523CandidateAlternatives: b523CandidateAlternatives}
}

func b523CandidateEvaluate(b523CandidateAlternative B523CandidateAlternative, b523CandidateRaw byte) (string, *float64, *int) {
	for _, b523CandidateRuleValue := range b523CandidateAlternative.B523CandidateRules {
		if b523CandidateRaw < b523CandidateRuleValue.B523CandidateRawMin || b523CandidateRaw > b523CandidateRuleValue.B523CandidateRawMax {
			continue
		}
		var b523CandidateNumeric *float64
		var b523CandidateDuration *int
		if b523CandidateRuleValue.B523CandidateDivideBy != nil {
			b523CandidateNumeric = b523CandidatePtr(float64(b523CandidateRaw) / *b523CandidateRuleValue.B523CandidateDivideBy)
		}
		if b523CandidateRuleValue.B523CandidateDurationUnit != nil {
			b523CandidateDuration = b523CandidatePtr(int(b523CandidateRaw))
		}
		return b523CandidateRuleValue.B523CandidateResult, b523CandidateNumeric, b523CandidateDuration
	}
	return "outside_domain", nil, nil
}

func TestAb71FiniteCatalog(b523CandidateT *b523CandidateTesting.T) {
	b523CandidateCatalog, b523CandidateError := B523Candidates()
	if b523CandidateError != nil {
		b523CandidateT.Fatal(b523CandidateError)
	}
	if b523CandidateCatalog.B523CandidateRepresentationGrade != "Proven" || b523CandidateCatalog.B523CandidateNativeSupport != "Unknown" || b523CandidateCatalog.B523CandidateRuntimeEffects != "Unknown" || b523CandidateCatalog.B523CandidateOtherModuleLayouts != "Unknown" {
		b523CandidateT.Fatal("claim boundaries")
	}
	b523CandidateLimitations := []string{
		"Finite candidate metadata; broader family catalog remains incomplete.",
		"Logical positions include the three-position prefix; abstract extents are not physical reply lengths.",
		"Acceptance is not physical ACK, persistence or readback.",
		"No routing, timed relay controller, physical qualification or runtime admission.",
		"Native qualification, freshness and execution authority are not established by candidate metadata.",
	}
	if !b523CandidateReflect.DeepEqual(b523CandidateCatalog.B523CandidateLimitations, b523CandidateLimitations) {
		b523CandidateT.Fatal("limitations")
	}
	if len(b523CandidateCatalog.B523CandidateEntries) != 26 {
		b523CandidateT.Fatalf("entries: %d", len(b523CandidateCatalog.B523CandidateEntries))
	}
	b523CandidateSeen := map[string]bool{}
	for _, b523CandidateEntry := range b523CandidateCatalog.B523CandidateEntries {
		b523CandidateID := b523CandidateFmt.Sprintf("%d/%s", *b523CandidateEntry.B523CandidateKey.B523CandidateOperation, b523CandidateEntry.B523CandidateKey.B523CandidateField)
		if b523CandidateSeen[b523CandidateID] {
			b523CandidateT.Fatal("duplicate", b523CandidateID)
		}
		b523CandidateSeen[b523CandidateID] = true
		if !b523CandidateReflect.DeepEqual(b523CandidateEntry, b523CandidateLookup(b523CandidateT, b523CandidateEntry.B523CandidateKey)) {
			b523CandidateT.Fatal("accessor/lookup disagreement")
		}
	}
}

func TestAb71RelayPositions(b523CandidateT *b523CandidateTesting.T) {
	for _, b523CandidateOperation := range []byte{1, 5} {
		b523CandidateCount, b523CandidateExtent := 6, 10
		if b523CandidateOperation == 5 {
			b523CandidateCount, b523CandidateExtent = 12, 17
		}
		for b523CandidateIndex := 0; b523CandidateIndex < b523CandidateCount; b523CandidateIndex++ {
			b523CandidateT.Run(b523CandidateFmt.Sprintf("ab71_op%d_relay%d", b523CandidateOperation, b523CandidateIndex+1), func(b523CandidateT *b523CandidateTesting.T) {
				b523CandidateField := b523CandidateFmt.Sprintf("relay_%d", b523CandidateIndex+1)
				b523CandidateEntry := b523CandidateLookup(b523CandidateT, b523CandidateKey(b523CandidateOperation, b523CandidateField))
				b523CandidateRules := []B523CandidateRule{b523CandidateRule(0, 0, "off", nil, nil), b523CandidateRule(1, 253, "on_then_off", nil, b523CandidatePtr("seconds")), b523CandidateRule(254, 254, "on", nil, nil), b523CandidateRule(255, 255, "unchanged", nil, nil)}
				b523CandidateWant := b523CandidateExpected(b523CandidateOperation, b523CandidateField, 3+b523CandidateIndex, b523CandidatePtr(b523CandidateExtent), []B523CandidateAlternative{{B523CandidateName: "relay_behavior", B523CandidateGrade: "Hypothesis", B523CandidateRules: b523CandidateRules}}, nil)
				if !b523CandidateReflect.DeepEqual(b523CandidateEntry, b523CandidateWant) {
					b523CandidateT.Fatalf("relay metadata: %+v", b523CandidateEntry)
				}
				for b523CandidateIndexValue, b523CandidateRaw := range []byte{0, 1, 253, 254, 255} {
					b523CandidateResult, b523CandidateNumeric, b523CandidateDuration := b523CandidateEvaluate(b523CandidateEntry.B523CandidateAlternatives[0], b523CandidateRaw)
					b523CandidateResults := []string{"off", "on_then_off", "on_then_off", "on", "unchanged"}
					b523CandidateDurations := []*int{nil, b523CandidatePtr(1), b523CandidatePtr(253), nil, nil}
					if b523CandidateResult != b523CandidateResults[b523CandidateIndexValue] || b523CandidateNumeric != nil || !b523CandidateReflect.DeepEqual(b523CandidateDuration, b523CandidateDurations[b523CandidateIndexValue]) {
						b523CandidateT.Fatalf("raw %d: %s %v", b523CandidateRaw, b523CandidateResult, b523CandidateDuration)
					}
				}
			})
		}
	}
}

func TestAb71PWMConflict(b523CandidateT *b523CandidateTesting.T) {
	for _, b523CandidateOperation := range []byte{1, 5} {
		b523CandidatePositions, b523CandidateExtent := []int{9}, 10
		if b523CandidateOperation == 5 {
			b523CandidatePositions, b523CandidateExtent = []int{15, 16}, 17
		}
		for b523CandidateIndex, b523CandidatePosition := range b523CandidatePositions {
			b523CandidateT.Run(b523CandidateFmt.Sprintf("ab71_op%d_pwm%d", b523CandidateOperation, b523CandidateIndex+1), func(b523CandidateT *b523CandidateTesting.T) {
				b523CandidateField := b523CandidateFmt.Sprintf("pwm_%d", b523CandidateIndex+1)
				b523CandidateEntry := b523CandidateLookup(b523CandidateT, b523CandidateKey(b523CandidateOperation, b523CandidateField))
				b523CandidateWant := b523CandidateExpected(b523CandidateOperation, b523CandidateField, b523CandidatePosition, b523CandidatePtr(b523CandidateExtent), []B523CandidateAlternative{
					{B523CandidateName: "numeric_codec", B523CandidateGrade: "Hypothesis", B523CandidateRules: []B523CandidateRule{b523CandidateRule(0, 250, "percent", b523CandidatePtr(2.0), nil)}},
					{B523CandidateName: "behavior", B523CandidateGrade: "Hypothesis", B523CandidateRules: []B523CandidateRule{b523CandidateRule(0, 200, "percent", b523CandidatePtr(2.0), nil), b523CandidateRule(201, 255, "unchanged", nil, nil)}},
				}, nil)
				if !b523CandidateReflect.DeepEqual(b523CandidateEntry, b523CandidateWant) {
					b523CandidateT.Fatal("PWM metadata")
				}
				for b523CandidateIndexValue, b523CandidateRaw := range []byte{200, 201, 250, 255} {
					b523CandidateResult, b523CandidateNumeric, b523CandidateDuration := b523CandidateEvaluate(b523CandidateEntry.B523CandidateAlternatives[0], b523CandidateRaw)
					b523CandidateNumbers := []*float64{b523CandidatePtr(100.0), b523CandidatePtr(100.5), b523CandidatePtr(125.0), nil}
					b523CandidateExpectedResult := "percent"
					if b523CandidateRaw == 255 {
						b523CandidateExpectedResult = "outside_domain"
					}
					if b523CandidateResult != b523CandidateExpectedResult || !b523CandidateReflect.DeepEqual(b523CandidateNumeric, b523CandidateNumbers[b523CandidateIndexValue]) || b523CandidateDuration != nil {
						b523CandidateT.Fatalf("numeric raw %d", b523CandidateRaw)
					}
					b523CandidateResult, b523CandidateNumeric, b523CandidateDuration = b523CandidateEvaluate(b523CandidateEntry.B523CandidateAlternatives[1], b523CandidateRaw)
					b523CandidateExpectedResult = "unchanged"
					var b523CandidateNumber *float64
					if b523CandidateRaw == 200 {
						b523CandidateExpectedResult = "percent"
						b523CandidateNumber = b523CandidatePtr(100.0)
					}
					if b523CandidateResult != b523CandidateExpectedResult || !b523CandidateReflect.DeepEqual(b523CandidateNumeric, b523CandidateNumber) || b523CandidateDuration != nil {
						b523CandidateT.Fatalf("behavior raw %d", b523CandidateRaw)
					}
				}
			})
		}
	}
}

func TestAb71Feedback(b523CandidateT *b523CandidateTesting.T) {
	for _, b523CandidateOperation := range []byte{0, 1, 2, 4, 5} {
		b523CandidateT.Run(b523CandidateFmt.Sprintf("ab71_op%d", b523CandidateOperation), func(b523CandidateT *b523CandidateTesting.T) {
			b523CandidateEntry := b523CandidateLookup(b523CandidateT, b523CandidateKey(b523CandidateOperation, "feedback"))
			b523CandidateRejected := "rejected"
			if b523CandidateOperation == 4 || b523CandidateOperation == 5 {
				b523CandidateRejected = "rejected-or-unsupported"
			}
			var b523CandidateExtent *int
			if b523CandidateOperation == 1 {
				b523CandidateExtent = b523CandidatePtr(10)
			}
			if b523CandidateOperation == 5 {
				b523CandidateExtent = b523CandidatePtr(17)
			}
			b523CandidateWant := b523CandidateExpected(b523CandidateOperation, "feedback", 3, b523CandidateExtent, []B523CandidateAlternative{{B523CandidateName: "feedback", B523CandidateGrade: "Hypothesis", B523CandidateRules: []B523CandidateRule{b523CandidateRule(0, 0, b523CandidateRejected, nil, nil), b523CandidateRule(1, 1, "accepted", nil, nil), b523CandidateRule(2, 255, b523CandidateRejected, nil, nil)}}}, b523CandidatePtr(b523CandidateRejected))
			if b523CandidateExtent == nil {
				b523CandidateWant.B523CandidatePositiveExtent = nil
			}
			if !b523CandidateReflect.DeepEqual(b523CandidateEntry, b523CandidateWant) {
				b523CandidateT.Fatal("feedback metadata")
			}
			for _, b523CandidateRaw := range []byte{0, 1, 2, 255} {
				b523CandidateWantResult := b523CandidateRejected
				if b523CandidateRaw == 1 {
					b523CandidateWantResult = "accepted"
				}
				b523CandidateResult, b523CandidateNumeric, b523CandidateDuration := b523CandidateEvaluate(b523CandidateEntry.B523CandidateAlternatives[0], b523CandidateRaw)
				if b523CandidateResult != b523CandidateWantResult || b523CandidateNumeric != nil || b523CandidateDuration != nil {
					b523CandidateT.Fatalf("feedback %d", b523CandidateRaw)
				}
			}
		})
	}
}

func TestAb71ExactKeyAxes(b523CandidateT *b523CandidateTesting.T) {
	b523CandidateMutations := map[string]func(*B523CandidateKey){
		"ab71_family": func(b523CandidateKeyValue *B523CandidateKey) {
			b523CandidateKeyValue.B523CandidateFamily = "other_family"
		},
		"ab71_variant": func(b523CandidateKeyValue *B523CandidateKey) {
			b523CandidateKeyValue.B523CandidateVariant = b523CandidatePtr("other_variant")
		},
		"ab71_variant_empty": func(b523CandidateKeyValue *B523CandidateKey) {
			b523CandidateKeyValue.B523CandidateVariant = b523CandidatePtr("")
		},
		"ab71_primary": func(b523CandidateKeyValue *B523CandidateKey) { b523CandidateKeyValue.B523CandidatePrimary = 0 },
		"ab71_secondary": func(b523CandidateKeyValue *B523CandidateKey) {
			b523CandidateKeyValue.B523CandidateSecondary = b523CandidatePtr(byte(0x24))
		},
		"ab71_secondary_absent": func(b523CandidateKeyValue *B523CandidateKey) { b523CandidateKeyValue.B523CandidateSecondary = nil },
		"ab71_secondary_zero": func(b523CandidateKeyValue *B523CandidateKey) {
			b523CandidateKeyValue.B523CandidateSecondary = b523CandidatePtr(byte(0))
		},
		"ab71_operation": func(b523CandidateKeyValue *B523CandidateKey) {
			b523CandidateKeyValue.B523CandidateOperation = b523CandidatePtr(byte(127))
		},
		"ab71_operation_absent": func(b523CandidateKeyValue *B523CandidateKey) { b523CandidateKeyValue.B523CandidateOperation = nil },
		"ab71_operation_zero": func(b523CandidateKeyValue *B523CandidateKey) {
			b523CandidateKeyValue.B523CandidateOperation = b523CandidatePtr(byte(0))
		},
		"ab71_direction": func(b523CandidateKeyValue *B523CandidateKey) {
			b523CandidateKeyValue.B523CandidateDirection = "response"
		},
		"ab71_field": func(b523CandidateKeyValue *B523CandidateKey) {
			b523CandidateKeyValue.B523CandidateField = "unknown_field"
		},
	}
	for b523CandidateName, b523CandidateMutation := range b523CandidateMutations {
		b523CandidateT.Run(b523CandidateName, func(b523CandidateT *b523CandidateTesting.T) {
			b523CandidateControl := b523CandidateKey(1, "relay_1")
			b523CandidateLookup(b523CandidateT, b523CandidateControl)
			b523CandidateChanged := b523CandidateKey(1, "relay_1")
			b523CandidateMutation(&b523CandidateChanged)
			if _, b523CandidateFound, b523CandidateError := B523CandidateLookup(b523CandidateChanged); b523CandidateFound || b523CandidateError != nil {
				b523CandidateT.Fatalf("single-axis mutation found=%v error=%v", b523CandidateFound, b523CandidateError)
			}
			b523CandidateLookup(b523CandidateT, b523CandidateControl)
		})
	}
	b523CandidateZero := b523CandidateKey(0, "feedback")
	b523CandidateLookup(b523CandidateT, b523CandidateZero)
	b523CandidateZero.B523CandidateOperation = nil
	if _, b523CandidateFound, b523CandidateError := B523CandidateLookup(b523CandidateZero); b523CandidateFound || b523CandidateError != nil {
		b523CandidateT.Fatal("absent operation matched zero")
	}
	b523CandidateAllocated := b523CandidateKey(1, "relay_1")
	if !b523CandidateReflect.DeepEqual(b523CandidateLookup(b523CandidateT, b523CandidateAllocated), b523CandidateLookup(b523CandidateT, b523CandidateKey(1, "relay_1"))) {
		b523CandidateT.Fatal("pointer identity used as key")
	}
}

func TestAb71DeepCopy(b523CandidateT *b523CandidateTesting.T) {
	b523CandidateBaseline, b523CandidateError := B523Candidates()
	if b523CandidateError != nil {
		b523CandidateT.Fatal(b523CandidateError)
	}
	b523CandidateChanged, b523CandidateError := B523Candidates()
	if b523CandidateError != nil {
		b523CandidateT.Fatal(b523CandidateError)
	}
	b523CandidateChanged.B523CandidateLimitations[0] = "changed"
	for b523CandidateIndex := range b523CandidateChanged.B523CandidateEntries {
		b523CandidateEntry := &b523CandidateChanged.B523CandidateEntries[b523CandidateIndex]
		*b523CandidateEntry.B523CandidateKey.B523CandidateSecondary = 0
		*b523CandidateEntry.B523CandidateKey.B523CandidateOperation = 127
		*b523CandidateEntry.B523CandidateLogicalPosition = 99
		*b523CandidateEntry.B523CandidateWidthBits = 99
		if b523CandidateEntry.B523CandidateRequestExtent != nil {
			*b523CandidateEntry.B523CandidateRequestExtent = 99
		}
		if b523CandidateEntry.B523CandidatePositiveExtent != nil {
			*b523CandidateEntry.B523CandidatePositiveExtent = 99
		}
		*b523CandidateEntry.B523CandidateDefaultRaw = 99
		if b523CandidateEntry.B523CandidateDefaultResult != nil {
			*b523CandidateEntry.B523CandidateDefaultResult = "changed"
		}
		b523CandidateEntry.B523CandidateRawDomain[0] = 99
		for b523CandidateAlternativeIndex := range b523CandidateEntry.B523CandidateAlternatives {
			b523CandidateAlternative := &b523CandidateEntry.B523CandidateAlternatives[b523CandidateAlternativeIndex]
			b523CandidateAlternative.B523CandidateName = "changed"
			for b523CandidateRuleIndex := range b523CandidateAlternative.B523CandidateRules {
				b523CandidateRuleValue := &b523CandidateAlternative.B523CandidateRules[b523CandidateRuleIndex]
				b523CandidateRuleValue.B523CandidateResult = "changed"
				if b523CandidateRuleValue.B523CandidateDivideBy != nil {
					*b523CandidateRuleValue.B523CandidateDivideBy = 99
				}
				if b523CandidateRuleValue.B523CandidateDurationUnit != nil {
					*b523CandidateRuleValue.B523CandidateDurationUnit = "changed"
				}
			}
		}
	}
	b523CandidateFresh, b523CandidateError := B523Candidates()
	if b523CandidateError != nil || !b523CandidateReflect.DeepEqual(b523CandidateFresh, b523CandidateBaseline) {
		b523CandidateT.Fatal("catalog storage shared")
	}
	for _, b523CandidateExpectedEntry := range b523CandidateBaseline.B523CandidateEntries {
		b523CandidateReturned := b523CandidateLookup(b523CandidateT, b523CandidateExpectedEntry.B523CandidateKey)
		*b523CandidateReturned.B523CandidateLogicalPosition = 99
		b523CandidateReturned.B523CandidateAlternatives[0].B523CandidateRules[0].B523CandidateResult = "changed"
		if !b523CandidateReflect.DeepEqual(b523CandidateLookup(b523CandidateT, b523CandidateExpectedEntry.B523CandidateKey), b523CandidateExpectedEntry) {
			b523CandidateT.Fatal("lookup storage shared")
		}
	}
}
