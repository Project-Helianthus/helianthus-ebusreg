package ebusstandard

import (
	standard0fCandidateReflect "reflect"
	standard0fCandidateTesting "testing"
)

func standard0fCandidateExpected() []Standard0fCandidate {
	standard0fCandidateFalse := false
	return []Standard0fCandidate{
		{
			Standard0fCandidateKey:           Standard0fCandidateKey{Standard0fCandidateFamily: "testing", Standard0fCandidateVariant: "partial", Standard0fCandidatePrimary: 0x0F, Standard0fCandidateDirection: "request", Standard0fCandidateField: "declaration"},
			Standard0fCandidateMetadataGrade: "Proven", Standard0fCandidateGrade: "Hypothesis", Standard0fCandidateBody: "Unknown",
			Standard0fCandidateACK: "Unknown", Standard0fCandidateSupport: "Unknown", Standard0fCandidateEffects: "Unknown", Standard0fCandidateRaw: []byte{0x0F},
			Standard0fCandidatePositive: &Standard0fCandidatePositive{Standard0fCandidatePrefix: "PeBUS", Standard0fCandidateTypedDataPresent: &standard0fCandidateFalse, Standard0fCandidateGrade: "Hypothesis"},
			Standard0fCandidateLimits: []string{
				"Logical declaration only; no physical reply length is established.",
				"Absent secondary is distinct from 00 and standard services 0F01/0F02/0F03.",
				"The broader testing-family catalog remains incomplete; physical qualification and runtime adoption remain open before 1.0.",
			},
		},
		{
			Standard0fCandidateKey:           Standard0fCandidateKey{Standard0fCandidateFamily: "testing", Standard0fCandidateVariant: "partial", Standard0fCandidatePrimary: 0x0F, Standard0fCandidateDirection: "response", Standard0fCandidateField: "opaque_reply"},
			Standard0fCandidateMetadataGrade: "Proven", Standard0fCandidateGrade: "Unknown", Standard0fCandidateBody: "Unknown",
			Standard0fCandidateACK: "Unknown", Standard0fCandidateSupport: "Unknown", Standard0fCandidateEffects: "Unknown", Standard0fCandidateRaw: []byte{0xAA},
			Standard0fCandidateLimits: []string{"Synthetic opaque reply example only; AA remains raw and Unknown, without physical qualification."},
		},
	}
}

func TestAb81ProductionMetadata(standard0fCandidateT *standard0fCandidateTesting.T) {
	standard0fCandidateWant := standard0fCandidateExpected()
	if standard0fCandidateGot := Standard0fCandidates(); !standard0fCandidateReflect.DeepEqual(standard0fCandidateGot, standard0fCandidateWant) {
		standard0fCandidateT.Fatalf("bounded production metadata differs: got %+v, want %+v", standard0fCandidateGot, standard0fCandidateWant)
	}
	for _, standard0fCandidateEntry := range standard0fCandidateWant {
		standard0fCandidateGot, standard0fCandidateOK := Standard0fCandidateLookup(standard0fCandidateEntry.Standard0fCandidateKey)
		if !standard0fCandidateOK || !standard0fCandidateReflect.DeepEqual(standard0fCandidateGot, standard0fCandidateEntry) {
			standard0fCandidateT.Fatalf("exact production lookup differs: %+v, found %v", standard0fCandidateGot, standard0fCandidateOK)
		}
	}
}

func TestAb81ExactKeyAxes(standard0fCandidateT *standard0fCandidateTesting.T) {
	standard0fCandidateBase := standard0fCandidateExpected()[0].Standard0fCandidateKey
	standard0fCandidateCases := []struct {
		standard0fCandidateName   string
		standard0fCandidateChange func(*Standard0fCandidateKey)
	}{
		{"ab81_family", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateKey.Standard0fCandidateFamily = "other"
		}},
		{"ab81_variant", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateKey.Standard0fCandidateVariant = "other"
		}},
		{"ab81_primary", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateKey.Standard0fCandidatePrimary = 0x10
		}},
		{"ab81_secondary_00", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateValue := byte(0)
			standard0fCandidateKey.Standard0fCandidateSecondary = &standard0fCandidateValue
		}},
		{"ab81_secondary_01", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateValue := byte(1)
			standard0fCandidateKey.Standard0fCandidateSecondary = &standard0fCandidateValue
		}},
		{"ab81_secondary_02", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateValue := byte(2)
			standard0fCandidateKey.Standard0fCandidateSecondary = &standard0fCandidateValue
		}},
		{"ab81_secondary_03", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateValue := byte(3)
			standard0fCandidateKey.Standard0fCandidateSecondary = &standard0fCandidateValue
		}},
		{"ab81_operation_00", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateValue := byte(0)
			standard0fCandidateKey.Standard0fCandidateOperation = &standard0fCandidateValue
		}},
		{"ab81_operation_01", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateValue := byte(1)
			standard0fCandidateKey.Standard0fCandidateOperation = &standard0fCandidateValue
		}},
		{"ab81_direction", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateKey.Standard0fCandidateDirection = "response"
		}},
		{"ab81_field", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateKey.Standard0fCandidateField = "opaque_reply"
		}},
		{"ab81_unknown_family", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateKey.Standard0fCandidateFamily = ""
		}},
		{"ab81_unknown_variant", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateKey.Standard0fCandidateVariant = ""
		}},
		{"ab81_unknown_direction", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateKey.Standard0fCandidateDirection = ""
		}},
		{"ab81_unknown_field", func(standard0fCandidateKey *Standard0fCandidateKey) {
			standard0fCandidateKey.Standard0fCandidateField = ""
		}},
	}
	for _, standard0fCandidateCase := range standard0fCandidateCases {
		standard0fCandidateT.Run(standard0fCandidateCase.standard0fCandidateName, func(standard0fCandidateT *standard0fCandidateTesting.T) {
			if _, standard0fCandidateOK := Standard0fCandidateLookup(standard0fCandidateBase); !standard0fCandidateOK {
				standard0fCandidateT.Fatal("valid control missing")
			}
			standard0fCandidateChanged := standard0fCandidateBase
			standard0fCandidateCase.standard0fCandidateChange(&standard0fCandidateChanged)
			standard0fCandidateGot, standard0fCandidateOK := Standard0fCandidateLookup(standard0fCandidateChanged)
			if standard0fCandidateOK || !standard0fCandidateReflect.DeepEqual(standard0fCandidateGot, Standard0fCandidate{}) {
				standard0fCandidateT.Fatalf("changed axis matched: %+v", standard0fCandidateChanged)
			}
		})
	}
}

func TestAb81FreshNestedData(standard0fCandidateT *standard0fCandidateTesting.T) {
	standard0fCandidateWant := standard0fCandidateExpected()
	standard0fCandidateEntries := Standard0fCandidates()
	standard0fCandidateEntries[0].Standard0fCandidateRaw[0] = 0
	standard0fCandidateEntries[0].Standard0fCandidateLimits[0] = "changed"
	*standard0fCandidateEntries[0].Standard0fCandidatePositive.Standard0fCandidateTypedDataPresent = true
	standard0fCandidateEntries[0].Standard0fCandidatePositive.Standard0fCandidatePrefix = "changed"
	standard0fCandidateEntries[0].Standard0fCandidateKey.Standard0fCandidateFamily = "changed"
	standard0fCandidateEntries[1].Standard0fCandidateRaw[0] = 0
	standard0fCandidateEntries[1].Standard0fCandidateLimits[0] = "changed"
	standard0fCandidateLookup, standard0fCandidateOK := Standard0fCandidateLookup(standard0fCandidateWant[0].Standard0fCandidateKey)
	if !standard0fCandidateOK {
		standard0fCandidateT.Fatal("lookup control missing")
	}
	standard0fCandidateLookup.Standard0fCandidateRaw[0] = 0
	standard0fCandidateLookup.Standard0fCandidateLimits[0] = "changed"
	*standard0fCandidateLookup.Standard0fCandidatePositive.Standard0fCandidateTypedDataPresent = true
	standard0fCandidateLookup.Standard0fCandidatePositive.Standard0fCandidateGrade = "changed"
	if standard0fCandidateGot := Standard0fCandidates(); !standard0fCandidateReflect.DeepEqual(standard0fCandidateGot, standard0fCandidateWant) {
		standard0fCandidateT.Fatal("accessor shares nested storage")
	}
	for _, standard0fCandidateEntry := range standard0fCandidateWant {
		standard0fCandidateGot, standard0fCandidateOK := Standard0fCandidateLookup(standard0fCandidateEntry.Standard0fCandidateKey)
		if !standard0fCandidateOK || !standard0fCandidateReflect.DeepEqual(standard0fCandidateGot, standard0fCandidateEntry) {
			standard0fCandidateT.Fatal("lookup shares nested storage")
		}
	}
}
