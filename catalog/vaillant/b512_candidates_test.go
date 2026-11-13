package vaillant

import (
	b512CandidateReflect "reflect"
	b512CandidateStrconv "strconv"
	b512CandidateTesting "testing"
)

func b512CandidateByte(b512CandidateValue byte) *byte { return &b512CandidateValue }

func b512CandidateGet(b512CandidateT *b512CandidateTesting.T, b512CandidateOperation byte, b512CandidateVariant, b512CandidateDirection, b512CandidateField string) B512Candidate {
	b512CandidateT.Helper()
	b512CandidateEntry, b512CandidateOK := B512CandidateLookup(B512CandidateKey{"B512", b512CandidateVariant, 181, b512CandidateByte(18), b512CandidateByte(b512CandidateOperation), b512CandidateDirection, b512CandidateField})
	if !b512CandidateOK {
		b512CandidateT.Fatalf("missing OP%02d %s", b512CandidateOperation, b512CandidateField)
	}
	return b512CandidateEntry
}

func b512CandidateMeaning(b512CandidateEntry B512Candidate, b512CandidateRaw int) string {
	if b512CandidateValue, b512CandidateOK := b512CandidateEntry.B512CandidateMappings[b512CandidateStrconv.Itoa(b512CandidateRaw)]; b512CandidateOK {
		return b512CandidateValue
	}
	for _, b512CandidateRange := range b512CandidateEntry.B512CandidateRanges {
		if b512CandidateRaw >= b512CandidateRange.B512CandidateMin && b512CandidateRaw <= b512CandidateRange.B512CandidateMax {
			return b512CandidateRange.B512CandidateMeaning
		}
	}
	return "Unknown"
}

func TestAb69Metadata(b512CandidateT *b512CandidateTesting.T) {
	b512CandidateCatalog := B512Candidates()
	if b512CandidateCatalog.B512CandidateRepresentationGrade != "Proven" || b512CandidateCatalog.B512CandidatePhysicalQualification != "Unknown" || b512CandidateCatalog.B512CandidateInstalledSupport != "Unknown" || b512CandidateCatalog.B512CandidateResponseAdaptation != "Unknown" || b512CandidateCatalog.B512CandidateActualSelection != "Unknown" || b512CandidateCatalog.B512CandidateScope == "" || len(b512CandidateCatalog.B512CandidateEntries) != 15 {
		b512CandidateT.Fatal("catalog boundary changed")
	}
	for _, b512CandidateEntry := range b512CandidateCatalog.B512CandidateEntries {
		b512CandidateFound, b512CandidateOK := B512CandidateLookup(b512CandidateEntry.B512CandidateKey)
		if !b512CandidateOK || !b512CandidateReflect.DeepEqual(b512CandidateFound, b512CandidateEntry) || b512CandidateEntry.B512CandidateGrade != "Hypothesis" {
			b512CandidateT.Fatal("entry/lookup/grade mismatch")
		}
		if *b512CandidateEntry.B512CandidateKey.B512CandidateOperation != 4 && (b512CandidateEntry.B512CandidateTimeoutSeconds != nil || b512CandidateEntry.B512CandidateTimeoutRestart != nil) {
			b512CandidateT.Fatal("timeout leaked")
		}
		if b512CandidateEntry.B512CandidateKey.B512CandidateField != "virtual_error" && b512CandidateMeaning(b512CandidateEntry, 65535) == "no-error" {
			b512CandidateT.Fatal("virtual error sentinel leaked across fields")
		}
		if *b512CandidateEntry.B512CandidateKey.B512CandidateOperation != 5 && b512CandidateEntry.B512CandidateExternalMode != nil {
			b512CandidateT.Fatal("external mode leaked")
		}
	}
}

func TestAb69ExactKey(b512CandidateT *b512CandidateTesting.T) {
	b512CandidateControl := B512CandidateKey{"B512", "circulation", 181, b512CandidateByte(18), b512CandidateByte(0), "request", "setpoint"}
	b512CandidateCases := []struct {
		b512CandidateName   string
		b512CandidateMutate func(*B512CandidateKey)
	}{
		{"family", func(b512CandidateKey *B512CandidateKey) { b512CandidateKey.B512CandidateFamily = "Unknown" }},
		{"variant", func(b512CandidateKey *B512CandidateKey) { b512CandidateKey.B512CandidateVariant = "Unknown" }},
		{"primary", func(b512CandidateKey *B512CandidateKey) { b512CandidateKey.B512CandidatePrimary = 0 }},
		{"secondary-absent", func(b512CandidateKey *B512CandidateKey) { b512CandidateKey.B512CandidateSecondary = nil }},
		{"secondary-zero", func(b512CandidateKey *B512CandidateKey) {
			b512CandidateKey.B512CandidateSecondary = b512CandidateByte(0)
		}},
		{"operation-absent", func(b512CandidateKey *B512CandidateKey) { b512CandidateKey.B512CandidateOperation = nil }},
		{"operation-other", func(b512CandidateKey *B512CandidateKey) {
			b512CandidateKey.B512CandidateOperation = b512CandidateByte(2)
		}},
		{"direction", func(b512CandidateKey *B512CandidateKey) { b512CandidateKey.B512CandidateDirection = "positive" }},
		{"field", func(b512CandidateKey *B512CandidateKey) { b512CandidateKey.B512CandidateField = "Unknown" }},
	}
	for _, b512CandidateCase := range b512CandidateCases {
		b512CandidateT.Run("ab69-"+b512CandidateCase.b512CandidateName, func(b512CandidateT *b512CandidateTesting.T) {
			if _, b512CandidateOK := B512CandidateLookup(b512CandidateControl); !b512CandidateOK {
				b512CandidateT.Fatal("valid control failed")
			}
			b512CandidateChanged := b512CandidateControl
			b512CandidateCase.b512CandidateMutate(&b512CandidateChanged)
			if _, b512CandidateOK := B512CandidateLookup(b512CandidateChanged); b512CandidateOK {
				b512CandidateT.Fatal("axis ignored or wildcarded")
			}
		})
	}
	if b512CandidateOptionalEqual(nil, b512CandidateByte(0)) || b512CandidateOptionalEqual(b512CandidateByte(0), nil) || !b512CandidateOptionalEqual(nil, nil) {
		b512CandidateT.Fatal("absent conflated with zero")
	}
}

func TestAb69PumpValveCases(b512CandidateT *b512CandidateTesting.T) {
	b512CandidateCases := []struct {
		b512CandidateOp                                                  byte
		b512CandidateVariant, b512CandidateDirection, b512CandidateField string
		b512CandidateRaw                                                 int
		b512CandidateWant                                                string
	}{
		{0, "circulation", "request", "setpoint", 0, "off"},
		{0, "circulation", "request", "setpoint", 100, "on/speed"},
		{0, "circulation", "request", "setpoint", 101, "unresolved"},
		{0, "circulation", "request", "setpoint", 255, "unresolved"},
		{2, "valve", "request", "valve", 0, "CH"},
		{2, "valve", "request", "valve", 50, "midpoint"},
		{2, "valve", "request", "valve", 100, "DHW"},
		{2, "valve", "request", "valve", 254, "standby"},
		{2, "valve", "request", "valve", 255, "unchanged"},
		{2, "valve", "positive", "contact", 0, "open"},
		{2, "valve", "positive", "contact", 1, "closed"},
		{2, "valve", "positive", "virtual_error", 65535, "no-error"},
		{2, "valve", "positive", "tank_temperature", 65535, "Unknown"},
		{4, "heating-pump", "request", "setpoint", 0, "off regardless of internal demand"},
		{4, "heating-pump", "request", "setpoint", 100, "on/speed"},
		{4, "heating-pump", "request", "setpoint", 255, "replacement/no-remote-control"},
	}
	for _, b512CandidateCase := range b512CandidateCases {
		b512CandidateEntry := b512CandidateGet(b512CandidateT, b512CandidateCase.b512CandidateOp, b512CandidateCase.b512CandidateVariant, b512CandidateCase.b512CandidateDirection, b512CandidateCase.b512CandidateField)
		if b512CandidateMeaning(b512CandidateEntry, b512CandidateCase.b512CandidateRaw) != b512CandidateCase.b512CandidateWant {
			b512CandidateT.Fatalf("%+v", b512CandidateCase)
		}
	}
	b512CandidateCirculation := b512CandidateGet(b512CandidateT, 0, "circulation", "request", "setpoint")
	if !b512CandidateReflect.DeepEqual(b512CandidateCirculation.B512CandidateRawDomain, []int{0, 255}) || *b512CandidateCirculation.B512CandidateDefault != 0 || *b512CandidateCirculation.B512CandidateUnit != "percent" {
		b512CandidateT.Fatal("circulation domain/default/unit")
	}
	b512CandidateAck := b512CandidateGet(b512CandidateT, 0, "circulation", "positive", "acknowledgement")
	if b512CandidateAck.B512CandidatePosition != nil || b512CandidateAck.B512CandidateWidth != nil || !b512CandidateReflect.DeepEqual(b512CandidateAck.B512CandidateNotes, []string{"prefix-only", "physical acknowledgement: Unknown"}) {
		b512CandidateT.Fatal("ack invented")
	}
	b512CandidateFeedback := b512CandidateGet(b512CandidateT, 2, "valve", "positive", "valve_feedback")
	if !b512CandidateReflect.DeepEqual(b512CandidateFeedback.B512CandidateAlternatives, []string{"actual", "setpoint fallback"}) || !b512CandidateReflect.DeepEqual(b512CandidateFeedback.B512CandidateNotes, []string{"actual-versus-fallback selection: Unknown"}) {
		b512CandidateT.Fatal("fallback promoted")
	}
	b512CandidateTank := b512CandidateGet(b512CandidateT, 2, "valve", "positive", "tank_temperature")
	b512CandidateBytes := []byte{0xf0, 0xff}
	b512CandidateRaw := int(b512CandidateBytes[0]) | int(b512CandidateBytes[1])<<8
	if *b512CandidateTank.B512CandidateSigned && b512CandidateRaw >= 32768 {
		b512CandidateRaw -= 65536
	}
	if float64(b512CandidateRaw) / *b512CandidateTank.B512CandidateDivisor != -1 || *b512CandidateTank.B512CandidateByteOrder != "little-endian" || !b512CandidateReflect.DeepEqual(b512CandidateTank.B512CandidateRawDomain, []int{-32768, 32767}) {
		b512CandidateT.Fatal("signed temperature conversion")
	}
	b512CandidateError := b512CandidateGet(b512CandidateT, 2, "valve", "positive", "virtual_error")
	if *b512CandidateError.B512CandidateSigned || *b512CandidateError.B512CandidateByteOrder != "little-endian" || !b512CandidateReflect.DeepEqual(b512CandidateError.B512CandidateNotes, []string{"supplements fault-service information"}) {
		b512CandidateT.Fatal("virtual error metadata")
	}
	b512CandidateHeating := b512CandidateGet(b512CandidateT, 4, "heating-pump", "request", "setpoint")
	if *b512CandidateHeating.B512CandidateTimeoutSeconds != 900 || *b512CandidateHeating.B512CandidateTimeoutRestart != "every received service" {
		b512CandidateT.Fatal("timeout")
	}
	b512CandidateAccepted := b512CandidateGet(b512CandidateT, 4, "heating-pump", "positive", "feedback")
	if b512CandidateAccepted.B512CandidatePosition != nil {
		b512CandidateT.Fatal("unknown feedback position invented")
	}
	for b512CandidateRaw := 0; b512CandidateRaw <= 255; b512CandidateRaw++ {
		b512CandidateWant := "rejected"
		if b512CandidateRaw == 1 {
			b512CandidateWant = "accepted"
		}
		if b512CandidateMeaning(b512CandidateAccepted, b512CandidateRaw) != b512CandidateWant {
			b512CandidateT.Fatal("acceptance byte")
		}
	}
	for b512CandidateChannel := 1; b512CandidateChannel <= 3; b512CandidateChannel++ {
		b512CandidateRelay := b512CandidateGet(b512CandidateT, 5, "relay", "request", "relay_"+b512CandidateStrconv.Itoa(b512CandidateChannel))
		if *b512CandidateRelay.B512CandidatePosition != b512CandidateChannel+2 || *b512CandidateRelay.B512CandidateWidth != 1 || *b512CandidateRelay.B512CandidateExternalMode != 8 || b512CandidateRelay.B512CandidateTimeoutSeconds != nil || b512CandidateRelay.B512CandidateTimeoutRestart != nil || !b512CandidateReflect.DeepEqual(b512CandidateRelay.B512CandidateRawDomain, []int{0, 255}) {
			b512CandidateT.Fatal("relay independent metadata")
		}
		for _, b512CandidateCase := range []struct {
			b512CandidateRaw  int
			b512CandidateWant string
		}{{0, "off"}, {1, "on-or-speed"}, {100, "on-or-speed"}, {255, "unchanged"}} {
			if b512CandidateMeaning(b512CandidateRelay, b512CandidateCase.b512CandidateRaw) != b512CandidateCase.b512CandidateWant {
				b512CandidateT.Fatal("relay meaning")
			}
		}
	}
}

func TestAb69CoordinatesPower(b512CandidateT *b512CandidateTesting.T) {
	for _, b512CandidateCase := range []struct {
		b512CandidateOp                                                  byte
		b512CandidateVariant, b512CandidateDirection, b512CandidateField string
		b512CandidatePosition, b512CandidateWidth                        int
	}{
		{0, "circulation", "request", "setpoint", 3, 1}, {2, "valve", "request", "valve", 3, 1}, {2, "valve", "positive", "valve_feedback", 3, 1}, {2, "valve", "positive", "contact", 4, 1}, {2, "valve", "positive", "tank_temperature", 5, 2}, {2, "valve", "positive", "virtual_error", 7, 2}, {4, "heating-pump", "request", "setpoint", 3, 1}, {7, "power", "positive", "maximum_power", 3, 1}, {7, "power", "positive", "minimum_power", 4, 1},
	} {
		b512CandidateEntry := b512CandidateGet(b512CandidateT, b512CandidateCase.b512CandidateOp, b512CandidateCase.b512CandidateVariant, b512CandidateCase.b512CandidateDirection, b512CandidateCase.b512CandidateField)
		if *b512CandidateEntry.B512CandidatePosition != b512CandidateCase.b512CandidatePosition || *b512CandidateEntry.B512CandidateWidth != b512CandidateCase.b512CandidateWidth {
			b512CandidateT.Fatal("coordinate/width")
		}
	}
	b512CandidateArgs := b512CandidateGet(b512CandidateT, 7, "power", "request", "arguments")
	if b512CandidateArgs.B512CandidatePosition != nil || b512CandidateArgs.B512CandidateWidth != nil || !b512CandidateReflect.DeepEqual(b512CandidateArgs.B512CandidateNotes, []string{"no arguments"}) {
		b512CandidateT.Fatal("arguments invented")
	}
	b512CandidateMaximum := b512CandidateGet(b512CandidateT, 7, "power", "positive", "maximum_power")
	if *b512CandidateMaximum.B512CandidateDefault != 0 || *b512CandidateMaximum.B512CandidateUnit != "kW" || *b512CandidateMaximum.B512CandidateSigned || !b512CandidateReflect.DeepEqual(b512CandidateMaximum.B512CandidateRawDomain, []int{0, 255}) {
		b512CandidateT.Fatal("maximum power")
	}
	b512CandidateMinimum := b512CandidateGet(b512CandidateT, 7, "power", "positive", "minimum_power")
	if *b512CandidateMinimum.B512CandidateDefault != 0 || *b512CandidateMinimum.B512CandidateSigned || *b512CandidateMinimum.B512CandidateUnit != "percent" || !b512CandidateReflect.DeepEqual(b512CandidateMinimum.B512CandidateRawDomain, []int{0, 250}) || 250 / *b512CandidateMinimum.B512CandidateDivisor != 125 || b512CandidateMeaning(b512CandidateMinimum, 251) != "outside declared domain; retain raw" {
		b512CandidateT.Fatal("minimum power clipped or domain expanded")
	}
}

func TestAb69DeepCopy(b512CandidateT *b512CandidateTesting.T) {
	b512CandidateBefore := B512Candidates()
	b512CandidateChanged := B512Candidates()
	b512CandidateChanged.B512CandidateEntries[0].B512CandidateMappings["0"] = "changed"
	b512CandidateChanged.B512CandidateEntries[0].B512CandidateRanges[0].B512CandidateMin = -1
	b512CandidateChanged.B512CandidateEntries[0].B512CandidateRawDomain[1] = 0
	*b512CandidateChanged.B512CandidateEntries[0].B512CandidatePosition = 99
	*b512CandidateChanged.B512CandidateEntries[0].B512CandidateKey.B512CandidateOperation = 99
	b512CandidateChanged.B512CandidateEntries[1].B512CandidateNotes[0] = "changed"
	b512CandidateChanged.B512CandidateEntries[3].B512CandidateAlternatives[0] = "changed"
	if !b512CandidateReflect.DeepEqual(b512CandidateBefore, B512Candidates()) {
		b512CandidateT.Fatal("shared accessor storage")
	}
	b512CandidateLookup := b512CandidateGet(b512CandidateT, 0, "circulation", "request", "setpoint")
	b512CandidateLookup.B512CandidateMappings["0"] = "changed"
	*b512CandidateLookup.B512CandidateWidth = 99
	if !b512CandidateReflect.DeepEqual(b512CandidateBefore.B512CandidateEntries[0], b512CandidateGet(b512CandidateT, 0, "circulation", "request", "setpoint")) {
		b512CandidateT.Fatal("shared lookup storage")
	}
}
