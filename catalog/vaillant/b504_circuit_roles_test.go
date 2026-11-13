package vaillant

import (
	b504CircuitRoleBytes "bytes"
	b504CircuitRoleJSON "encoding/json"
	b504CircuitRoleReflect "reflect"
	b504CircuitRoleStrings "strings"
	b504CircuitRoleTesting "testing"
)

func b504CircuitRolePointer[b504CircuitRoleT any](b504CircuitRoleValue b504CircuitRoleT) *b504CircuitRoleT {
	return &b504CircuitRoleValue
}

func b504CircuitRoleEqual(b504CircuitRoleT *b504CircuitRoleTesting.T, b504CircuitRoleGot, b504CircuitRoleWant any) {
	b504CircuitRoleT.Helper()
	if !b504CircuitRoleReflect.DeepEqual(b504CircuitRoleGot, b504CircuitRoleWant) {
		b504CircuitRoleT.Fatalf("ab15 got %#v; want %#v", b504CircuitRoleGot, b504CircuitRoleWant)
	}
}

func b504CircuitRoleKey(b504CircuitRoleOperation uint8, b504CircuitRoleField, b504CircuitRoleAlternative string) B504CircuitRoleKey {
	return B504CircuitRoleKey{B504CircuitRoleFamily: "vaillant", B504CircuitRoleVariant: "B504", B504CircuitRolePrimary: 181, B504CircuitRoleSecondary: b504CircuitRolePointer(uint8(4)), B504CircuitRoleOperation: b504CircuitRolePointer(b504CircuitRoleOperation), B504CircuitRoleDirection: "none", B504CircuitRoleField: b504CircuitRoleField, B504CircuitRoleAlternative: b504CircuitRoleAlternative, B504CircuitRoleInstanceScope: "dynamic"}
}

func b504CircuitRoleGet(b504CircuitRoleT *b504CircuitRoleTesting.T, b504CircuitRoleKey B504CircuitRoleKey) b504CircuitRoleRow {
	b504CircuitRoleT.Helper()
	b504CircuitRoleRaw, b504CircuitRoleFound, b504CircuitRoleErr := B504CircuitRoleLookup(b504CircuitRoleKey)
	if b504CircuitRoleErr != nil || !b504CircuitRoleFound {
		b504CircuitRoleT.Fatalf("ab15 lookup %v: %v, %v", b504CircuitRoleKey, b504CircuitRoleFound, b504CircuitRoleErr)
	}
	var b504CircuitRoleEntry b504CircuitRoleRow
	if b504CircuitRoleErr := b504CircuitRoleJSON.Unmarshal(b504CircuitRoleRaw, &b504CircuitRoleEntry); b504CircuitRoleErr != nil {
		b504CircuitRoleT.Fatal(b504CircuitRoleErr)
	}
	b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleEntry.B504CircuitRoleKey, b504CircuitRoleKey)
	return b504CircuitRoleEntry
}

func TestAb15FiniteRoles(b504CircuitRoleT *b504CircuitRoleTesting.T) {
	b504CircuitRoleCases := []struct {
		b504CircuitRoleOperation                                              uint8
		b504CircuitRoleField, b504CircuitRoleAlternative, b504CircuitRoleRole string
		b504CircuitRoleApplies                                                []string
		b504CircuitRoleRelations                                              []string
	}{
		{11, "reload-delay", "none", "solar-storage-charging-delay", []string{"solar-storage"}, []string{"reload delay concerns solar storage charging"}},
		{11, "parallel-tank-filling", "none", "global-mixer-circuit-filling", []string{"all-mixer-circuits"}, []string{"parallel-tank filling applies globally to mixer circuits"}},
		{45, "room-setpoint", "none", "inverter-indoor-room-setpoint", []string{"inverter-indoor-circuit"}, []string{"room temperature is the regulation input"}},
		{45, "room-temperature", "none", "regulation-input", []string{"inverter-indoor-circuit"}, []string{"regulation input for room setpoint"}},
		{13, "system-cooling-support", "none", "system-cooling-support", []string{"system"}, []string{"request bit 0 expresses cooling support", "system generation occupies a different byte"}},
		{13, "system-generation", "none", "system-generation", []string{"system"}, []string{"different byte from cooling-support request bit 0"}},
		{13, "dhw-request", "dhw", "dhw-request", []string{"dhw-circuit"}, []string{"DHW interpretation applies only to a DHW circuit"}},
		{13, "circuit-temperature", "burner", "common-sensor", []string{"burner-circuit"}, []string{"circuit role determines temperature meaning"}},
		{13, "circuit-temperature", "mixer", "flow", []string{"mixer-circuit"}, []string{"circuit role determines temperature meaning"}},
		{13, "circuit-temperature", "dhw", "tank", []string{"dhw-circuit"}, []string{"circuit role determines temperature meaning"}},
		{13, "desired-temperature", "burner", "desired-common-sensor", []string{"burner-circuit"}, []string{"circuit role determines desired-temperature meaning"}},
		{13, "desired-temperature", "mixer", "desired-flow", []string{"mixer-circuit"}, []string{"circuit role determines desired-temperature meaning"}},
		{13, "desired-temperature", "dhw", "desired-tank", []string{"dhw-circuit"}, []string{"circuit role determines desired-temperature meaning"}},
		{24, "desired-flow", "unsuppressed", "unsuppressed-desired-flow", []string{"circuit"}, []string{"independent of suppressed desired flow during tank charging"}},
		{24, "desired-flow", "suppressed", "suppressed-desired-flow", []string{"circuit"}, []string{"zero during tank charging does not replace unsuppressed demand"}},
		{1, "drying-day", "none", "drying-program-indication", []string{"circuit"}, []string{"nonzero drying day identifies a drying program", "zero has no nonzero-program indication"}},
		{51, "averaged-outside-temperature", "none", "correction-included-average", []string{"system"}, []string{"averaged outside temperature includes configured correction"}},
		{33, "maximum-solar-temperature", "none", "maximum-solar-temperature", []string{"all-solar-circuits"}, []string{"maximum solar temperature spans solar circuits"}},
		{33, "tank-state", "none", "temperature-to-maximum-percentage", []string{"solar-circuit"}, []string{"distinct calculated quantity from collector-to-tank switching ratio"}},
		{33, "solar-yield", "none", "collector-to-tank-switching-ratio", []string{"solar-circuit"}, []string{"distinct calculated quantity from temperature-to-maximum percentage"}},
	}
	b504CircuitRoleData, b504CircuitRoleErr := B504CircuitRoleEntries()
	if b504CircuitRoleErr != nil {
		b504CircuitRoleT.Fatal(b504CircuitRoleErr)
	}
	var b504CircuitRoleRows []b504CircuitRoleJSON.RawMessage
	if b504CircuitRoleErr = b504CircuitRoleJSON.Unmarshal(b504CircuitRoleData, &b504CircuitRoleRows); b504CircuitRoleErr != nil {
		b504CircuitRoleT.Fatal(b504CircuitRoleErr)
	}
	b504CircuitRoleEqual(b504CircuitRoleT, len(b504CircuitRoleRows), len(b504CircuitRoleCases))
	for _, b504CircuitRoleCase := range b504CircuitRoleCases {
		b504CircuitRoleT.Run("ab15-"+b504CircuitRoleCase.b504CircuitRoleField+"-"+b504CircuitRoleCase.b504CircuitRoleAlternative, func(b504CircuitRoleT *b504CircuitRoleTesting.T) {
			b504CircuitRoleEntry := b504CircuitRoleGet(b504CircuitRoleT, b504CircuitRoleKey(b504CircuitRoleCase.b504CircuitRoleOperation, b504CircuitRoleCase.b504CircuitRoleField, b504CircuitRoleCase.b504CircuitRoleAlternative))
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleEntry.B504CircuitRoleID, "ab15-"+b504CircuitRoleCase.b504CircuitRoleField+"-"+b504CircuitRoleCase.b504CircuitRoleAlternative)
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleEntry.B504CircuitRole, b504CircuitRoleCase.b504CircuitRoleRole)
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleEntry.B504CircuitRoleAppliesTo, b504CircuitRoleCase.b504CircuitRoleApplies)
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleEntry.B504CircuitRoleRelations, b504CircuitRoleCase.b504CircuitRoleRelations)
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleEntry.B504CircuitRoleGrade, "Hypothesis")
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleEntry.B504CircuitRoleMetadataPresence, "Proven")
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleEntry.B504CircuitRolePhysicalQualification, "Unknown")
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleEntry.B504CircuitRoleNativeContract, "protocols/ebus-services/ebus-overview.md")
			b504CircuitRoleUnknown := []string{"byte positions", "widths", "wire direction", "generation values", "suppression algorithms", "native support", "physical effects", "persistence", "delivery"}
			if b504CircuitRoleCase.b504CircuitRoleField == "averaged-outside-temperature" {
				b504CircuitRoleUnknown = append(b504CircuitRoleUnknown, "averaging and correction arithmetic")
			}
			if b504CircuitRoleCase.b504CircuitRoleField == "tank-state" || b504CircuitRoleCase.b504CircuitRoleField == "solar-yield" {
				b504CircuitRoleUnknown = append(b504CircuitRoleUnknown, "exact formula")
			}
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleEntry.B504CircuitRoleUnknown, b504CircuitRoleUnknown)
			var b504CircuitRoleRaw map[string]any
			if b504CircuitRoleErr := b504CircuitRoleJSON.Unmarshal(b504CircuitRoleEntry.B504CircuitRoleRaw, &b504CircuitRoleRaw); b504CircuitRoleErr != nil {
				b504CircuitRoleT.Fatal(b504CircuitRoleErr)
			}
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleRaw, map[string]any{"form": "opaque", "preserve": true, "interpretation": "Unknown"})
			b504CircuitRoleWhole, b504CircuitRoleFound, b504CircuitRoleErr := B504CircuitRoleLookup(b504CircuitRoleEntry.B504CircuitRoleKey)
			if b504CircuitRoleErr != nil || !b504CircuitRoleFound {
				b504CircuitRoleT.Fatal(b504CircuitRoleErr)
			}
			b504CircuitRoleMatches := 0
			for _, b504CircuitRoleOriginal := range b504CircuitRoleRows {
				var b504CircuitRoleOriginalEntry b504CircuitRoleRow
				if b504CircuitRoleErr := b504CircuitRoleJSON.Unmarshal(b504CircuitRoleOriginal, &b504CircuitRoleOriginalEntry); b504CircuitRoleErr != nil {
					b504CircuitRoleT.Fatal(b504CircuitRoleErr)
				}
				if b504CircuitRoleReflect.DeepEqual(b504CircuitRoleOriginalEntry.B504CircuitRoleKey, b504CircuitRoleEntry.B504CircuitRoleKey) {
					b504CircuitRoleEqual(b504CircuitRoleT, []byte(b504CircuitRoleWhole), []byte(b504CircuitRoleOriginal))
					b504CircuitRoleMatches++
				}
			}
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleMatches, 1)
			b504CircuitRoleCondition := "none"
			var b504CircuitRoleBit *uint8
			if b504CircuitRoleCase.b504CircuitRoleField == "drying-day" {
				b504CircuitRoleCondition = "nonzero"
			}
			if b504CircuitRoleCase.b504CircuitRoleField == "system-cooling-support" {
				b504CircuitRoleBit = b504CircuitRolePointer(uint8(0))
			}
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleEntry.B504CircuitRoleCondition, b504CircuitRoleCondition)
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleEntry.B504CircuitRoleBit, b504CircuitRoleBit)
		})
	}
}

func TestAb15SyntheticCells(b504CircuitRoleT *b504CircuitRoleTesting.T) {
	for _, b504CircuitRoleCase := range []struct{ b504CircuitRoleContext, b504CircuitRoleRole string }{{"burner", "common-sensor"}, {"mixer", "flow"}, {"dhw", "tank"}} {
		b504CircuitRoleEntry := b504CircuitRoleGet(b504CircuitRoleT, b504CircuitRoleKey(13, "circuit-temperature", b504CircuitRoleCase.b504CircuitRoleContext))
		for _, b504CircuitRoleToken := range []int{40, 41} {
			b504CircuitRoleCell := map[string]int{b504CircuitRoleEntry.B504CircuitRole: b504CircuitRoleToken}
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleCell, map[string]int{b504CircuitRoleCase.b504CircuitRoleRole: b504CircuitRoleToken})
		}
	}
	b504CircuitRoleDHW := b504CircuitRoleGet(b504CircuitRoleT, b504CircuitRoleKey(13, "dhw-request", "dhw"))
	b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleDHW.B504CircuitRoleAppliesTo, []string{"dhw-circuit"})
	b504CircuitRoleRaw, b504CircuitRoleFound, b504CircuitRoleErr := B504CircuitRoleLookup(b504CircuitRoleKey(13, "dhw-request", "mixer"))
	if b504CircuitRoleErr != nil || b504CircuitRoleFound || b504CircuitRoleRaw != nil {
		b504CircuitRoleT.Fatal("ab15 mixer acquired DHW interpretation")
	}
	for _, b504CircuitRoleCase := range []struct{ b504CircuitRoleUnsuppressed, b504CircuitRoleSuppressed int }{{45, 0}, {46, 0}, {45, 1}} {
		b504CircuitRoleDemand := b504CircuitRoleGet(b504CircuitRoleT, b504CircuitRoleKey(24, "desired-flow", "unsuppressed"))
		b504CircuitRoleSuppressed := b504CircuitRoleGet(b504CircuitRoleT, b504CircuitRoleKey(24, "desired-flow", "suppressed"))
		b504CircuitRoleCells := map[string]int{b504CircuitRoleDemand.B504CircuitRole: b504CircuitRoleCase.b504CircuitRoleUnsuppressed, b504CircuitRoleSuppressed.B504CircuitRole: b504CircuitRoleCase.b504CircuitRoleSuppressed}
		b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleCells, map[string]int{"unsuppressed-desired-flow": b504CircuitRoleCase.b504CircuitRoleUnsuppressed, "suppressed-desired-flow": b504CircuitRoleCase.b504CircuitRoleSuppressed})
	}
	b504CircuitRoleDrying := b504CircuitRoleGet(b504CircuitRoleT, b504CircuitRoleKey(1, "drying-day", "none"))
	for _, b504CircuitRoleCase := range []struct {
		b504CircuitRoleDay       int
		b504CircuitRoleIndicated bool
	}{{0, false}, {1, true}, {2, true}} {
		b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleDrying.B504CircuitRoleCondition == "nonzero" && b504CircuitRoleCase.b504CircuitRoleDay != 0, b504CircuitRoleCase.b504CircuitRoleIndicated)
	}
	for _, b504CircuitRoleCase := range []struct{ b504CircuitRolePercentage, b504CircuitRoleRatio int }{{50, 50}, {51, 50}, {50, 51}} {
		b504CircuitRoleTank := b504CircuitRoleGet(b504CircuitRoleT, b504CircuitRoleKey(33, "tank-state", "none"))
		b504CircuitRoleYield := b504CircuitRoleGet(b504CircuitRoleT, b504CircuitRoleKey(33, "solar-yield", "none"))
		b504CircuitRoleCells := map[string]int{b504CircuitRoleTank.B504CircuitRole: b504CircuitRoleCase.b504CircuitRolePercentage, b504CircuitRoleYield.B504CircuitRole: b504CircuitRoleCase.b504CircuitRoleRatio}
		b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleCells, map[string]int{"temperature-to-maximum-percentage": b504CircuitRoleCase.b504CircuitRolePercentage, "collector-to-tank-switching-ratio": b504CircuitRoleCase.b504CircuitRoleRatio})
	}
}

func TestAb15EveryKeyAxis(b504CircuitRoleT *b504CircuitRoleTesting.T) {
	b504CircuitRoleBase := b504CircuitRoleKey(11, "reload-delay", "none")
	b504CircuitRoleControls := []struct {
		b504CircuitRoleName  string
		b504CircuitRoleValue any
	}{
		{"B504CircuitRoleFamily", "other-family"}, {"B504CircuitRoleVariant", "other-variant"}, {"B504CircuitRolePrimary", uint8(180)},
		{"B504CircuitRoleSecondary", b504CircuitRolePointer(uint8(5))}, {"B504CircuitRoleOperation", b504CircuitRolePointer(uint8(12))},
		{"B504CircuitRoleDirection", "request"}, {"B504CircuitRoleField", "other-field"}, {"B504CircuitRoleAlternative", "other-role"},
		{"B504CircuitRoleGroup", b504CircuitRolePointer(uint8(1))}, {"B504CircuitRoleInstanceScope", "00"}, {"B504CircuitRoleRegister", b504CircuitRolePointer(uint16(1))},
		{"B504CircuitRoleSecondary", (*uint8)(nil)}, {"B504CircuitRoleSecondary", b504CircuitRolePointer(uint8(0))},
		{"B504CircuitRoleOperation", (*uint8)(nil)}, {"B504CircuitRoleOperation", b504CircuitRolePointer(uint8(0))},
		{"B504CircuitRoleGroup", b504CircuitRolePointer(uint8(0))}, {"B504CircuitRoleRegister", b504CircuitRolePointer(uint16(0))},
	}
	b504CircuitRoleOriginal, b504CircuitRoleFound, b504CircuitRoleErr := B504CircuitRoleLookup(b504CircuitRoleBase)
	if b504CircuitRoleErr != nil || !b504CircuitRoleFound {
		b504CircuitRoleT.Fatal("ab15 positive control missing", b504CircuitRoleErr)
	}
	for _, b504CircuitRoleControl := range b504CircuitRoleControls {
		b504CircuitRoleT.Run("ab15-"+b504CircuitRoleControl.b504CircuitRoleName, func(b504CircuitRoleT *b504CircuitRoleTesting.T) {
			b504CircuitRoleChanged := b504CircuitRoleBase
			b504CircuitRoleReflect.ValueOf(&b504CircuitRoleChanged).Elem().FieldByName(b504CircuitRoleControl.b504CircuitRoleName).Set(b504CircuitRoleReflect.ValueOf(b504CircuitRoleControl.b504CircuitRoleValue))
			b504CircuitRoleRaw, b504CircuitRoleFound, b504CircuitRoleErr := B504CircuitRoleLookup(b504CircuitRoleChanged)
			if b504CircuitRoleErr != nil || b504CircuitRoleFound || b504CircuitRoleRaw != nil {
				b504CircuitRoleT.Fatal("ab15 inexact key admitted", b504CircuitRoleErr)
			}
			b504CircuitRoleAgain, b504CircuitRoleFound, b504CircuitRoleErr := B504CircuitRoleLookup(b504CircuitRoleBase)
			if b504CircuitRoleErr != nil || !b504CircuitRoleFound {
				b504CircuitRoleT.Fatal("ab15 positive control lost", b504CircuitRoleErr)
			}
			b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleAgain, b504CircuitRoleOriginal)
			b504CircuitRoleOne := b504CircuitRoleGet(b504CircuitRoleT, b504CircuitRoleBase)
			b504CircuitRoleTwo := b504CircuitRoleOne
			b504CircuitRoleTwo.B504CircuitRoleID = "ab15-synthetic-axis-control"
			b504CircuitRoleTwo.B504CircuitRoleKey = b504CircuitRoleChanged
			b504CircuitRoleFixture, b504CircuitRoleErr := b504CircuitRoleJSON.Marshal([]b504CircuitRoleRow{b504CircuitRoleOne, b504CircuitRoleTwo})
			if b504CircuitRoleErr != nil {
				b504CircuitRoleT.Fatal(b504CircuitRoleErr)
			}
			if _, b504CircuitRoleErr = b504CircuitRoleEntries(b504CircuitRoleFixture); b504CircuitRoleErr != nil {
				b504CircuitRoleT.Fatal("ab15 distinct identities rejected", b504CircuitRoleErr)
			}
			for _, b504CircuitRoleKey := range []B504CircuitRoleKey{b504CircuitRoleBase, b504CircuitRoleChanged} {
				b504CircuitRoleRaw, b504CircuitRoleFound, b504CircuitRoleErr := b504CircuitRoleLookup(b504CircuitRoleFixture, b504CircuitRoleKey)
				if b504CircuitRoleErr != nil || !b504CircuitRoleFound {
					b504CircuitRoleT.Fatal("ab15 seam control missing", b504CircuitRoleErr)
				}
				var b504CircuitRoleEntry b504CircuitRoleRow
				if b504CircuitRoleErr = b504CircuitRoleJSON.Unmarshal(b504CircuitRoleRaw, &b504CircuitRoleEntry); b504CircuitRoleErr != nil {
					b504CircuitRoleT.Fatal(b504CircuitRoleErr)
				}
				b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleEntry.B504CircuitRoleKey, b504CircuitRoleKey)
			}
		})
	}
}

func TestAb15AbsentVersusZero(b504CircuitRoleT *b504CircuitRoleTesting.T) {
	for _, b504CircuitRoleAxis := range []string{"B504CircuitRoleSecondary", "B504CircuitRoleOperation", "B504CircuitRoleGroup", "B504CircuitRoleRegister"} {
		b504CircuitRoleT.Run("ab15-"+b504CircuitRoleAxis, func(b504CircuitRoleT *b504CircuitRoleTesting.T) {
			b504CircuitRoleAbsent := b504CircuitRoleGet(b504CircuitRoleT, b504CircuitRoleKey(11, "reload-delay", "none"))
			b504CircuitRoleField := b504CircuitRoleReflect.ValueOf(&b504CircuitRoleAbsent.B504CircuitRoleKey).Elem().FieldByName(b504CircuitRoleAxis)
			b504CircuitRoleField.Set(b504CircuitRoleReflect.Zero(b504CircuitRoleField.Type()))
			b504CircuitRoleZero := b504CircuitRoleAbsent
			b504CircuitRoleZero.B504CircuitRoleID = "ab15-synthetic-zero"
			b504CircuitRoleField = b504CircuitRoleReflect.ValueOf(&b504CircuitRoleZero.B504CircuitRoleKey).Elem().FieldByName(b504CircuitRoleAxis)
			b504CircuitRoleField.Set(b504CircuitRoleReflect.New(b504CircuitRoleField.Type().Elem()))
			b504CircuitRoleFixture, b504CircuitRoleErr := b504CircuitRoleJSON.Marshal([]b504CircuitRoleRow{b504CircuitRoleAbsent, b504CircuitRoleZero})
			if b504CircuitRoleErr != nil {
				b504CircuitRoleT.Fatal(b504CircuitRoleErr)
			}
			for _, b504CircuitRoleEntry := range []b504CircuitRoleRow{b504CircuitRoleAbsent, b504CircuitRoleZero} {
				b504CircuitRoleRaw, b504CircuitRoleFound, b504CircuitRoleErr := b504CircuitRoleLookup(b504CircuitRoleFixture, b504CircuitRoleEntry.B504CircuitRoleKey)
				if b504CircuitRoleErr != nil || !b504CircuitRoleFound {
					b504CircuitRoleT.Fatal("ab15 absent/zero control", b504CircuitRoleErr)
				}
				var b504CircuitRoleGot b504CircuitRoleRow
				if b504CircuitRoleErr = b504CircuitRoleJSON.Unmarshal(b504CircuitRoleRaw, &b504CircuitRoleGot); b504CircuitRoleErr != nil {
					b504CircuitRoleT.Fatal(b504CircuitRoleErr)
				}
				b504CircuitRoleExpected, b504CircuitRoleErr := b504CircuitRoleJSON.Marshal(b504CircuitRoleEntry)
				if b504CircuitRoleErr != nil {
					b504CircuitRoleT.Fatal(b504CircuitRoleErr)
				}
				b504CircuitRoleEqual(b504CircuitRoleT, []byte(b504CircuitRoleRaw), b504CircuitRoleExpected)
				b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleGot.B504CircuitRoleKey, b504CircuitRoleEntry.B504CircuitRoleKey)
			}
		})
	}
}

func TestAb15ValidationErrors(b504CircuitRoleT *b504CircuitRoleTesting.T) {
	b504CircuitRoleEntry := b504CircuitRoleGet(b504CircuitRoleT, b504CircuitRoleKey(11, "reload-delay", "none"))
	b504CircuitRoleDuplicate := b504CircuitRoleEntry
	b504CircuitRoleDuplicate.B504CircuitRoleID = "ab15-duplicate-key"
	b504CircuitRoleFixture, b504CircuitRoleErr := b504CircuitRoleJSON.Marshal([]b504CircuitRoleRow{b504CircuitRoleEntry, b504CircuitRoleDuplicate})
	if b504CircuitRoleErr != nil {
		b504CircuitRoleT.Fatal(b504CircuitRoleErr)
	}
	b504CircuitRoleInputs := []string{"{", "null", "[]", "{}", "[null]", "[{}]", string(b504CircuitRoleFixture), string(b504CircuitRoleFixture) + " []"}
	b504CircuitRoleValid, b504CircuitRoleErr := b504CircuitRoleJSON.Marshal([]b504CircuitRoleRow{b504CircuitRoleEntry})
	if b504CircuitRoleErr != nil {
		b504CircuitRoleT.Fatal(b504CircuitRoleErr)
	}
	b504CircuitRoleInputs = append(b504CircuitRoleInputs, b504CircuitRoleStrings.Replace(string(b504CircuitRoleValid), `"Hypothesis"`, `"Proven"`, 1), b504CircuitRoleStrings.Replace(string(b504CircuitRoleValid), `"physical_qualification":"Unknown"`, `"physical_qualification":"Proven"`, 1), b504CircuitRoleStrings.Replace(string(b504CircuitRoleValid), `"raw":{"form":"opaque","preserve":true,"interpretation":"Unknown"}`, `"raw":null`, 1))
	for _, b504CircuitRoleAxis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
		var b504CircuitRoleObject []map[string]b504CircuitRoleJSON.RawMessage
		if b504CircuitRoleErr := b504CircuitRoleJSON.Unmarshal(b504CircuitRoleValid, &b504CircuitRoleObject); b504CircuitRoleErr != nil {
			b504CircuitRoleT.Fatal(b504CircuitRoleErr)
		}
		var b504CircuitRoleIdentity map[string]b504CircuitRoleJSON.RawMessage
		if b504CircuitRoleErr := b504CircuitRoleJSON.Unmarshal(b504CircuitRoleObject[0]["key"], &b504CircuitRoleIdentity); b504CircuitRoleErr != nil {
			b504CircuitRoleT.Fatal(b504CircuitRoleErr)
		}
		delete(b504CircuitRoleIdentity, b504CircuitRoleAxis)
		b504CircuitRoleObject[0]["key"], b504CircuitRoleErr = b504CircuitRoleJSON.Marshal(b504CircuitRoleIdentity)
		if b504CircuitRoleErr != nil {
			b504CircuitRoleT.Fatal(b504CircuitRoleErr)
		}
		b504CircuitRoleMissing, b504CircuitRoleErr := b504CircuitRoleJSON.Marshal(b504CircuitRoleObject)
		if b504CircuitRoleErr != nil {
			b504CircuitRoleT.Fatal(b504CircuitRoleErr)
		}
		b504CircuitRoleInputs = append(b504CircuitRoleInputs, string(b504CircuitRoleMissing))
	}
	for b504CircuitRoleIndex, b504CircuitRoleInput := range b504CircuitRoleInputs {
		b504CircuitRoleRaw, b504CircuitRoleErr := b504CircuitRoleEntries([]byte(b504CircuitRoleInput))
		if b504CircuitRoleErr == nil || b504CircuitRoleRaw != nil {
			b504CircuitRoleT.Fatalf("ab15 invalid entries accepted: %d", b504CircuitRoleIndex)
		}
		b504CircuitRoleRaw, b504CircuitRoleFound, b504CircuitRoleErr := b504CircuitRoleLookup([]byte(b504CircuitRoleInput), b504CircuitRoleEntry.B504CircuitRoleKey)
		if b504CircuitRoleErr == nil || b504CircuitRoleRaw != nil || b504CircuitRoleFound {
			b504CircuitRoleT.Fatalf("ab15 invalid lookup accepted: %d", b504CircuitRoleIndex)
		}
	}
}

func TestAb15IndependentBuffers(b504CircuitRoleT *b504CircuitRoleTesting.T) {
	b504CircuitRoleBefore, b504CircuitRoleErr := B504CircuitRoleEntries()
	if b504CircuitRoleErr != nil {
		b504CircuitRoleT.Fatal(b504CircuitRoleErr)
	}
	b504CircuitRoleOther, b504CircuitRoleErr := B504CircuitRoleEntries()
	if b504CircuitRoleErr != nil {
		b504CircuitRoleT.Fatal(b504CircuitRoleErr)
	}
	for b504CircuitRoleIndex := range b504CircuitRoleOther {
		b504CircuitRoleOther[b504CircuitRoleIndex] = 'x'
	}
	b504CircuitRoleAfter, b504CircuitRoleErr := B504CircuitRoleEntries()
	if b504CircuitRoleErr != nil {
		b504CircuitRoleT.Fatal(b504CircuitRoleErr)
	}
	b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleAfter, b504CircuitRoleBefore)
	b504CircuitRoleKey := b504CircuitRoleKey(11, "reload-delay", "none")
	b504CircuitRoleOne, b504CircuitRoleFound, b504CircuitRoleErr := B504CircuitRoleLookup(b504CircuitRoleKey)
	if b504CircuitRoleErr != nil || !b504CircuitRoleFound {
		b504CircuitRoleT.Fatal(b504CircuitRoleErr)
	}
	b504CircuitRoleTwo, b504CircuitRoleFound, b504CircuitRoleErr := B504CircuitRoleLookup(b504CircuitRoleKey)
	if b504CircuitRoleErr != nil || !b504CircuitRoleFound {
		b504CircuitRoleT.Fatal(b504CircuitRoleErr)
	}
	for b504CircuitRoleIndex := range b504CircuitRoleOne {
		b504CircuitRoleOne[b504CircuitRoleIndex] = 'x'
	}
	b504CircuitRoleThree, b504CircuitRoleFound, b504CircuitRoleErr := B504CircuitRoleLookup(b504CircuitRoleKey)
	if b504CircuitRoleErr != nil || !b504CircuitRoleFound {
		b504CircuitRoleT.Fatal(b504CircuitRoleErr)
	}
	b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleThree, b504CircuitRoleTwo)
	b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleBytes.Contains(b504CircuitRoleBefore, b504CircuitRoleThree), true)
	b504CircuitRoleSource := append([]byte(nil), b504CircuitRoleBefore...)
	b504CircuitRoleCopy, b504CircuitRoleErr := b504CircuitRoleEntries(b504CircuitRoleSource)
	if b504CircuitRoleErr != nil {
		b504CircuitRoleT.Fatal(b504CircuitRoleErr)
	}
	for b504CircuitRoleIndex := range b504CircuitRoleCopy {
		b504CircuitRoleCopy[b504CircuitRoleIndex] = 'x'
	}
	b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleSource, []byte(b504CircuitRoleBefore))
	b504CircuitRoleCopy, b504CircuitRoleFound, b504CircuitRoleErr = b504CircuitRoleLookup(b504CircuitRoleSource, b504CircuitRoleKey)
	if b504CircuitRoleErr != nil || !b504CircuitRoleFound {
		b504CircuitRoleT.Fatal(b504CircuitRoleErr)
	}
	for b504CircuitRoleIndex := range b504CircuitRoleCopy {
		b504CircuitRoleCopy[b504CircuitRoleIndex] = 'x'
	}
	b504CircuitRoleEqual(b504CircuitRoleT, b504CircuitRoleSource, []byte(b504CircuitRoleBefore))
}
