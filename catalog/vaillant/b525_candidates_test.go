package vaillant

import (
	b525CandidateBytes "bytes"
	b525CandidateFmt "fmt"
	b525CandidateReflect "reflect"
	b525CandidateTesting "testing"
)

func b525CandidatePointer[b525CandidateT any](b525CandidateValue b525CandidateT) *b525CandidateT {
	return &b525CandidateValue
}

func b525CandidateKey(b525CandidateOp byte, b525CandidateDirection, b525CandidateField string) B525CandidateKey {
	return B525CandidateKey{"Vaillant radio", nil, 0xb5, b525CandidatePointer(byte(0x25)), b525CandidatePointer(b525CandidateOp), b525CandidateDirection, b525CandidateField}
}

func b525CandidateGet(b525CandidateT *b525CandidateTesting.T, b525CandidateOp byte, b525CandidateDirection, b525CandidateField string) B525CandidateEntry {
	b525CandidateT.Helper()
	b525CandidateEntry, b525CandidateOK := B525CandidateLookup(b525CandidateKey(b525CandidateOp, b525CandidateDirection, b525CandidateField))
	if !b525CandidateOK {
		b525CandidateT.Fatalf("ab73 missing %02x/%s/%s", b525CandidateOp, b525CandidateDirection, b525CandidateField)
	}
	return b525CandidateEntry
}

func TestAb73ExactAxes(b525CandidateT *b525CandidateTesting.T) {
	b525CandidateControl := b525CandidateKey(0, "Unknown", "layer")
	if b525CandidateEqualOptional[byte](nil, b525CandidatePointer(byte(0))) || b525CandidateEqualOptional[byte](b525CandidatePointer(byte(0)), nil) || !b525CandidateEqualOptional(b525CandidatePointer(byte(0)), b525CandidatePointer(byte(0))) {
		b525CandidateT.Fatal("ab73 optional absence collapsed into zero")
	}
	for b525CandidateName, b525CandidateMutate := range map[string]func(*B525CandidateKey){
		"ab73_family":           func(b525CandidateK *B525CandidateKey) { b525CandidateK.B525CandidateFamily = "" },
		"ab73_variant":          func(b525CandidateK *B525CandidateKey) { b525CandidateK.B525CandidateVariant = b525CandidatePointer("") },
		"ab73_primary":          func(b525CandidateK *B525CandidateKey) { b525CandidateK.B525CandidatePrimary = 0 },
		"ab73_secondary_absent": func(b525CandidateK *B525CandidateKey) { b525CandidateK.B525CandidateSecondary = nil },
		"ab73_secondary_zero": func(b525CandidateK *B525CandidateKey) {
			b525CandidateK.B525CandidateSecondary = b525CandidatePointer(byte(0))
		},
		"ab73_operation_absent": func(b525CandidateK *B525CandidateKey) { b525CandidateK.B525CandidateOperation = nil },
		"ab73_operation_other": func(b525CandidateK *B525CandidateKey) {
			b525CandidateK.B525CandidateOperation = b525CandidatePointer(byte(1))
		},
		"ab73_direction": func(b525CandidateK *B525CandidateKey) { b525CandidateK.B525CandidateDirection = "request" },
		"ab73_field":     func(b525CandidateK *B525CandidateKey) { b525CandidateK.B525CandidateField = "" },
	} {
		b525CandidateT.Run(b525CandidateName, func(b525CandidateT *b525CandidateTesting.T) {
			if _, b525CandidateOK := B525CandidateLookup(b525CandidateControl); !b525CandidateOK {
				b525CandidateT.Fatal("ab73 invalid control")
			}
			b525CandidateChanged := b525CandidateControl
			b525CandidateMutate(&b525CandidateChanged)
			if _, b525CandidateOK := B525CandidateLookup(b525CandidateChanged); b525CandidateOK {
				b525CandidateT.Fatal("ab73 axis wildcard")
			}
		})
	}
	b525CandidateCatalog := B525Candidates()
	if b525CandidateCatalog.B525CandidateRepresentationGrade != "Proven" || b525CandidateCatalog.B525CandidatePhysicalQualification != "Unknown" {
		b525CandidateT.Fatal("ab73 grade conflation")
	}
	if len(b525CandidateCatalog.B525CandidateEntries) != 127 {
		b525CandidateT.Fatalf("ab73 finite count %d", len(b525CandidateCatalog.B525CandidateEntries))
	}
	b525CandidateSeen := map[string]bool{}
	for _, b525CandidateEntry := range b525CandidateCatalog.B525CandidateEntries {
		b525CandidateIdentity := b525CandidateFmt.Sprintf("%s/%v/%d/%d/%d/%s/%s", b525CandidateEntry.B525CandidateKey.B525CandidateFamily, b525CandidateEntry.B525CandidateKey.B525CandidateVariant, b525CandidateEntry.B525CandidateKey.B525CandidatePrimary, *b525CandidateEntry.B525CandidateKey.B525CandidateSecondary, *b525CandidateEntry.B525CandidateKey.B525CandidateOperation, b525CandidateEntry.B525CandidateKey.B525CandidateDirection, b525CandidateEntry.B525CandidateKey.B525CandidateField)
		if b525CandidateSeen[b525CandidateIdentity] {
			b525CandidateT.Fatal("ab73 duplicate")
		}
		b525CandidateSeen[b525CandidateIdentity] = true
		b525CandidateFound, b525CandidateOK := B525CandidateLookup(b525CandidateEntry.B525CandidateKey)
		if !b525CandidateOK || !b525CandidateReflect.DeepEqual(b525CandidateFound, b525CandidateEntry) || b525CandidateEntry.B525CandidateGrade != "Hypothesis" || b525CandidateEntry.B525CandidateRawPolicy != "preserve_all" {
			b525CandidateT.Fatal("ab73 entry altered")
		}
	}
}

func TestAb73FreshNestedCopies(b525CandidateT *b525CandidateTesting.T) {
	b525CandidateBefore := B525Candidates()
	b525CandidateChanged := B525Candidates()
	b525CandidateChanged.B525CandidateLimitations[0] = "ab73_changed"
	for b525CandidateI := range b525CandidateChanged.B525CandidateEntries {
		b525CandidateE := &b525CandidateChanged.B525CandidateEntries[b525CandidateI]
		b525CandidateE.B525CandidateKey.B525CandidateFamily = "ab73_changed"
		*b525CandidateE.B525CandidateKey.B525CandidateSecondary = 0
		*b525CandidateE.B525CandidateKey.B525CandidateOperation = 255
		for _, b525CandidateP := range []*int{b525CandidateE.B525CandidatePosition, b525CandidateE.B525CandidateWidth, b525CandidateE.B525CandidateMinimum, b525CandidateE.B525CandidateMaximum, b525CandidateE.B525CandidateDefault, b525CandidateE.B525CandidateCount, b525CandidateE.B525CandidateStartIndex} {
			if b525CandidateP != nil {
				*b525CandidateP = 999
			}
		}
		b525CandidateE.B525CandidateMappings[0] = "ab73_changed"
		b525CandidateE.B525CandidateBits[0] = "ab73_changed"
		if len(b525CandidateE.B525CandidateUnknowns) > 0 {
			b525CandidateE.B525CandidateUnknowns[0] = "ab73_changed"
		}
	}
	if !b525CandidateReflect.DeepEqual(b525CandidateBefore, B525Candidates()) {
		b525CandidateT.Fatal("ab73 accessor shared backing")
	}
	b525CandidateLookup := b525CandidateGet(b525CandidateT, 0, "Unknown", "sender-role")
	b525CandidateLookup.B525CandidateMappings[1] = "ab73_changed"
	*b525CandidateLookup.B525CandidateKey.B525CandidateOperation = 255
	if !b525CandidateReflect.DeepEqual(b525CandidateBefore, B525Candidates()) || b525CandidateGet(b525CandidateT, 0, "Unknown", "sender-role").B525CandidateMappings[1] != "remote-eight" {
		b525CandidateT.Fatal("ab73 lookup shared backing")
	}
}

func TestAb73RolesAndHeader(b525CandidateT *b525CandidateTesting.T) {
	b525CandidateCodes := []int{1, 15, 16, 21, 22, 23, 48, 55, 112, 119, 240, 247, 253, 254}
	b525CandidateLabels := []string{"remote-eight", "outdoor-sensor", "primary-control", "base-station", "base-outdoor-sensor", "remote-four", "remote-one", "remote-five", "remote-two", "remote-six", "remote-three", "remote-seven", "diagnosis-tool", "range-extender"}
	for b525CandidateI, b525CandidateName := range []string{"sender-role", "receiver-role"} {
		b525CandidateE := b525CandidateGet(b525CandidateT, 0, "Unknown", b525CandidateName)
		if *b525CandidateE.B525CandidatePosition != 7+b525CandidateI || *b525CandidateE.B525CandidateWidth != 1 || *b525CandidateE.B525CandidateDefault != 1 || len(b525CandidateE.B525CandidateMappings) != 14 {
			b525CandidateT.Fatal("ab73 role extent")
		}
		for b525CandidateJ, b525CandidateCode := range b525CandidateCodes {
			if b525CandidateE.B525CandidateMappings[b525CandidateCode] != b525CandidateLabels[b525CandidateJ] {
				b525CandidateT.Fatal("ab73 role reversal")
			}
		}
		if _, b525CandidateOK := b525CandidateE.B525CandidateMappings[0]; b525CandidateOK {
			b525CandidateT.Fatal("ab73 unmapped 00")
		}
	}
	b525CandidateSender := b525CandidateGet(b525CandidateT, 0, "Unknown", "sender-role")
	b525CandidateReceiver := b525CandidateGet(b525CandidateT, 0, "Unknown", "receiver-role")
	if b525CandidateSender.B525CandidateMappings[1] != "remote-eight" || b525CandidateReceiver.B525CandidateMappings[254] != "range-extender" {
		b525CandidateT.Fatal("ab73 roles collapsed")
	}
	for _, b525CandidateCase := range []struct {
		b525CandidateName                                      string
		b525CandidatePos, b525CandidateWidth, b525CandidateMax int
		b525CandidateCodec                                     string
	}{{"layer", 3, 1, 255, "u8"}, {"radio-domain", 4, 2, 65535, "u16-le"}, {"request-id", 9, 1, 255, "u8"}} {
		b525CandidateE := b525CandidateGet(b525CandidateT, 0, "Unknown", b525CandidateCase.b525CandidateName)
		if *b525CandidateE.B525CandidatePosition != b525CandidateCase.b525CandidatePos || *b525CandidateE.B525CandidateWidth != b525CandidateCase.b525CandidateWidth || *b525CandidateE.B525CandidateMaximum != b525CandidateCase.b525CandidateMax || *b525CandidateE.B525CandidateMinimum != 0 || *b525CandidateE.B525CandidateDefault != 0 || b525CandidateE.B525CandidateCodec != b525CandidateCase.b525CandidateCodec {
			b525CandidateT.Fatal("ab73 OP00 metadata")
		}
	}
	b525CandidateStatus := b525CandidateGet(b525CandidateT, 0, "Unknown", "status")
	for b525CandidateI, b525CandidateName := range []string{"repeat", "group", "encrypted", "frame-type", "outdoor-sensor-pairing", "reserved", "reserved", "reserved"} {
		if b525CandidateStatus.B525CandidateBits[b525CandidateI] != b525CandidateName || *b525CandidateStatus.B525CandidateDefault != 0 {
			b525CandidateT.Fatal("ab73 status bits")
		}
	}
}

func b525CandidateSigned(b525CandidateE B525CandidateEntry, b525CandidateRaw byte) int {
	if b525CandidateE.B525CandidateCodec == "signed8" {
		return int(int8(b525CandidateRaw))
	}
	return int(b525CandidateRaw)
}

func TestAb73RSSIAndFeedback(b525CandidateT *b525CandidateTesting.T) {
	b525CandidateRSSI := b525CandidateGet(b525CandidateT, 3, "Unknown", "RSSI")
	if b525CandidateSigned(b525CandidateRSSI, 255) != -1 || b525CandidateRSSI.B525CandidateUnit != nil || *b525CandidateRSSI.B525CandidateMinimum != -128 || *b525CandidateRSSI.B525CandidateMaximum != 127 || *b525CandidateRSSI.B525CandidateDefault != 0 {
		b525CandidateT.Fatal("ab73 RSSI signed carrier or unit")
	}
	b525CandidateLink := b525CandidateGet(b525CandidateT, 3, "Unknown", "link-quality")
	if len(b525CandidateLink.B525CandidateBits) != 8 || *b525CandidateLink.B525CandidateDefault != 0 {
		b525CandidateT.Fatal("ab73 link bits")
	}
	for b525CandidateI := 0; b525CandidateI < 8; b525CandidateI++ {
		b525CandidateWant := b525CandidateFmt.Sprintf("link-quality-%d", b525CandidateI)
		if b525CandidateI == 7 {
			b525CandidateWant = "CRC-okay"
		}
		if b525CandidateLink.B525CandidateBits[b525CandidateI] != b525CandidateWant {
			b525CandidateT.Fatal("ab73 link bit independence")
		}
	}
	b525CandidateFeedback := b525CandidateGet(b525CandidateT, 4, "positive", "feedback")
	b525CandidateFlags := []string{}
	for b525CandidateI := 0; b525CandidateI < 8; b525CandidateI++ {
		if byte(7)&(1<<b525CandidateI) != 0 {
			b525CandidateFlags = append(b525CandidateFlags, b525CandidateFeedback.B525CandidateBits[b525CandidateI])
		}
	}
	if !b525CandidateReflect.DeepEqual(b525CandidateFlags, []string{"success", "too-short", "too-long"}) || b525CandidateFeedback.B525CandidateRelation != "independent-bits" {
		b525CandidateT.Fatal("ab73 feedback collapsed")
	}
	for _, b525CandidateD := range []string{"request", "positive"} {
		b525CandidateE := b525CandidateGet(b525CandidateT, 4, b525CandidateD, "bytes")
		if *b525CandidateE.B525CandidateCount != 16 || *b525CandidateE.B525CandidateWidth != 16 || *b525CandidateE.B525CandidateDefault != 0 || b525CandidateE.B525CandidateRelation != "Unknown" || b525CandidateE.B525CandidatePosition != nil {
			b525CandidateT.Fatal("ab73 array contract/echo")
		}
		for b525CandidateI := 0; b525CandidateI < 16; b525CandidateI++ {
			b525CandidateRaw := make([]byte, 16)
			b525CandidateRaw[b525CandidateI] = 0xaa
			for b525CandidateJ, b525CandidateB := range b525CandidateRaw {
				b525CandidateWant := 0
				if b525CandidateI == b525CandidateJ {
					b525CandidateWant = 170
				}
				if b525CandidateSigned(b525CandidateE, b525CandidateB) != b525CandidateWant {
					b525CandidateT.Fatal("ab73 independent raw byte")
				}
			}
		}
	}
}

func TestAb73SignedBlocks(b525CandidateT *b525CandidateTesting.T) {
	for b525CandidateOp := byte(16); b525CandidateOp <= 23; b525CandidateOp++ {
		b525CandidateE := b525CandidateGet(b525CandidateT, b525CandidateOp, "request", "values")
		if *b525CandidateE.B525CandidatePosition != 3 || *b525CandidateE.B525CandidateWidth != 32 || *b525CandidateE.B525CandidateCount != 32 || *b525CandidateE.B525CandidateStartIndex != int(b525CandidateOp-16)*32 || b525CandidateE.B525CandidateUnit != nil || len(b525CandidateE.B525CandidateMappings) != 0 || !b525CandidateReflect.DeepEqual(b525CandidateE.B525CandidateUnknowns, []string{"frequency-mapping", "scanning-effects"}) || *b525CandidateE.B525CandidateMinimum != -128 || *b525CandidateE.B525CandidateMaximum != 127 || *b525CandidateE.B525CandidateDefault != 0 {
			b525CandidateT.Fatal("ab73 block alignment/frequency")
		}
		for _, b525CandidateLength := range []int{31, 32, 33} {
			if (b525CandidateLength == *b525CandidateE.B525CandidateCount) != (b525CandidateLength == 32) {
				b525CandidateT.Fatal("ab73 incomplete/extra accepted")
			}
		}
		for b525CandidateI := 0; b525CandidateI < 32; b525CandidateI++ {
			for _, b525CandidateCase := range []struct {
				b525CandidateRaw  byte
				b525CandidateWant int
			}{{128, -128}, {0, 0}, {127, 127}} {
				b525CandidateRaw := make([]byte, 32)
				b525CandidateRaw[b525CandidateI] = b525CandidateCase.b525CandidateRaw
				for b525CandidateJ, b525CandidateB := range b525CandidateRaw {
					b525CandidateWant := 0
					if b525CandidateJ == b525CandidateI {
						b525CandidateWant = b525CandidateCase.b525CandidateWant
					}
					if b525CandidateSigned(b525CandidateE, b525CandidateB) != b525CandidateWant {
						b525CandidateT.Fatal("ab73 signed array")
					}
				}
			}
		}
		for _, b525CandidateRow := range B525Candidates().B525CandidateEntries {
			if *b525CandidateRow.B525CandidateKey.B525CandidateOperation == b525CandidateOp && b525CandidateRow.B525CandidateKey.B525CandidateDirection == "positive" && b525CandidateRow.B525CandidateKey.B525CandidateField != "primary" && b525CandidateRow.B525CandidateKey.B525CandidateField != "secondary" && b525CandidateRow.B525CandidateKey.B525CandidateField != "selector" {
				b525CandidateT.Fatal("ab73 positive is prefix-only")
			}
		}
	}
}

func TestAb73OpaqueAndConfiguration(b525CandidateT *b525CandidateTesting.T) {
	for _, b525CandidateOp := range []byte{1, 2, 24, 25, 26, 27} {
		b525CandidateDirections := []string{"Unknown"}
		if b525CandidateOp >= 24 {
			b525CandidateDirections = []string{"request", "positive"}
		}
		for _, b525CandidateD := range b525CandidateDirections {
			b525CandidateE := b525CandidateGet(b525CandidateT, b525CandidateOp, b525CandidateD, "data")
			if b525CandidateE.B525CandidateCodec != "opaque" || *b525CandidateE.B525CandidateMinimum != 0 || b525CandidateE.B525CandidateMaximum != nil || b525CandidateE.B525CandidateWidth != nil || b525CandidateE.B525CandidateRelation != "Unknown" {
				b525CandidateT.Fatal("ab73 opaque bounds")
			}
			for _, b525CandidateRaw := range [][]byte{{}, {0, 170, 255}} {
				b525CandidateCopy := b525CandidateBytes.Clone(b525CandidateRaw)
				if !b525CandidateBytes.Equal(b525CandidateRaw, b525CandidateCopy) || len(b525CandidateCopy) < *b525CandidateE.B525CandidateMinimum {
					b525CandidateT.Fatal("ab73 opaque representation")
				}
			}
			if b525CandidateOp < 24 {
				if *b525CandidateE.B525CandidatePosition != 3 {
					b525CandidateT.Fatal("ab73 opaque position")
				}
				continue
			}
			if b525CandidateE.B525CandidatePosition != nil {
				b525CandidateT.Fatal("ab73 invented layout")
			}
			b525CandidateConfig := b525CandidateGet(b525CandidateT, b525CandidateOp, b525CandidateD, "configuration")
			if *b525CandidateConfig.B525CandidatePosition != 3 || *b525CandidateConfig.B525CandidateDefault != 0 || b525CandidateConfig.B525CandidateBits[0] != "configuration" || len(b525CandidateConfig.B525CandidateBits) != 1 || b525CandidateConfig.B525CandidateMappings[0] != "read" || len(b525CandidateConfig.B525CandidateMappings) != 1 || b525CandidateConfig.B525CandidateRawPolicy != "preserve_all" {
				b525CandidateT.Fatal("ab73 configuration metadata")
			}
			for _, b525CandidateRaw := range []byte{0, 1, 0xaa, 0xff} {
				b525CandidateKnown := b525CandidateRaw & 1
				b525CandidateUnknown := b525CandidateRaw & 0xfe
				if b525CandidateKnown|b525CandidateUnknown != b525CandidateRaw {
					b525CandidateT.Fatal("ab73 unknown bits lost")
				}
			}
		}
	}
}

func TestAb73PrefixAndUnknowns(b525CandidateT *b525CandidateTesting.T) {
	for _, b525CandidateE := range B525Candidates().B525CandidateEntries {
		for b525CandidateI, b525CandidateName := range []string{"primary", "secondary", "selector"} {
			if b525CandidateE.B525CandidateKey.B525CandidateField != b525CandidateName {
				continue
			}
			b525CandidateValue := []int{181, 37, int(*b525CandidateE.B525CandidateKey.B525CandidateOperation)}[b525CandidateI]
			if *b525CandidateE.B525CandidatePosition != b525CandidateI || *b525CandidateE.B525CandidateWidth != 1 || *b525CandidateE.B525CandidateMinimum != b525CandidateValue || *b525CandidateE.B525CandidateMaximum != b525CandidateValue || b525CandidateE.B525CandidateDefault != nil {
				b525CandidateT.Fatal("ab73 prefix changed")
			}
		}
	}
	for _, b525CandidateC := range []struct {
		b525CandidateOp                byte
		b525CandidateD, b525CandidateF string
		b525CandidateUnknowns          []string
	}{
		{0, "Unknown", "request-id", []string{"lifecycle"}},
		{1, "Unknown", "data", []string{"cipher", "padding", "sender-identity"}},
		{2, "Unknown", "data", []string{"cipher", "padding", "sender-identity"}},
		{3, "Unknown", "RSSI", []string{"physical-unit"}},
		{3, "Unknown", "link-quality", []string{"CRC-algorithm"}},
		{4, "request", "bytes", []string{"active-key", "persistent-key", "echo"}},
		{4, "positive", "bytes", []string{"active-key", "persistent-key", "echo"}},
		{4, "positive", "feedback", []string{"exclusivity"}},
	} {
		if !b525CandidateReflect.DeepEqual(b525CandidateGet(b525CandidateT, b525CandidateC.b525CandidateOp, b525CandidateC.b525CandidateD, b525CandidateC.b525CandidateF).B525CandidateUnknowns, b525CandidateC.b525CandidateUnknowns) {
			b525CandidateT.Fatal("ab73 unknown became claim")
		}
	}
	for b525CandidateOp := byte(24); b525CandidateOp <= 27; b525CandidateOp++ {
		for _, b525CandidateD := range []string{"request", "positive"} {
			if !b525CandidateReflect.DeepEqual(b525CandidateGet(b525CandidateT, b525CandidateOp, b525CandidateD, "configuration").B525CandidateUnknowns, []string{"bits-1..7", "effects"}) || !b525CandidateReflect.DeepEqual(b525CandidateGet(b525CandidateT, b525CandidateOp, b525CandidateD, "data").B525CandidateUnknowns, []string{"layout", "units", "correlation", "effects"}) {
				b525CandidateT.Fatal("ab73 configuration claims")
			}
		}
	}
}
