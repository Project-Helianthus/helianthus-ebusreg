package vaillant

import (
	b524ComponentSelectorBytesPackage "bytes"
	b524ComponentSelectorJSON "encoding/json"
	b524ComponentSelectorReflect "reflect"
	b524ComponentSelectorTesting "testing"
)

func b524ComponentSelectorPointer[b524ComponentSelectorT ~uint8 | ~uint16](b524ComponentSelectorValue b524ComponentSelectorT) *b524ComponentSelectorT {
	return &b524ComponentSelectorValue
}

func b524ComponentSelectorBaseKey() B524ComponentSelectorKey {
	return B524ComponentSelectorKey{"controller-700", "none", 181, b524ComponentSelectorPointer(uint8(36)), b524ComponentSelectorPointer(uint8(2)), "read", "clock-date", "none", b524ComponentSelectorPointer(uint8(0)), "00", b524ComponentSelectorPointer(uint16(52))}
}

func b524ComponentSelectorGet(b524ComponentSelectorT *b524ComponentSelectorTesting.T, b524ComponentSelectorKey B524ComponentSelectorKey) map[string]b524ComponentSelectorJSON.RawMessage {
	b524ComponentSelectorT.Helper()
	b524ComponentSelectorRaw, b524ComponentSelectorFound, b524ComponentSelectorError := B524ComponentSelectorLookup(b524ComponentSelectorKey)
	if b524ComponentSelectorError != nil || !b524ComponentSelectorFound {
		b524ComponentSelectorT.Fatalf("lookup: %v %v", b524ComponentSelectorFound, b524ComponentSelectorError)
	}
	var b524ComponentSelectorRow map[string]b524ComponentSelectorJSON.RawMessage
	if b524ComponentSelectorError := b524ComponentSelectorJSON.Unmarshal(b524ComponentSelectorRaw, &b524ComponentSelectorRow); b524ComponentSelectorError != nil {
		b524ComponentSelectorT.Fatal(b524ComponentSelectorError)
	}
	return b524ComponentSelectorRow
}

func TestAb32FiniteSelectors(b524ComponentSelectorT *b524ComponentSelectorTesting.T) {
	b524ComponentSelectorEntries, b524ComponentSelectorError := B524ComponentSelectorEntries()
	if b524ComponentSelectorError != nil {
		b524ComponentSelectorT.Fatal(b524ComponentSelectorError)
	}
	var b524ComponentSelectorRows []struct {
		B524ComponentSelectorKey B524ComponentSelectorKey `json:"key"`
	}
	if b524ComponentSelectorError := b524ComponentSelectorJSON.Unmarshal(b524ComponentSelectorEntries, &b524ComponentSelectorRows); b524ComponentSelectorError != nil {
		b524ComponentSelectorT.Fatal(b524ComponentSelectorError)
	}
	if len(b524ComponentSelectorRows) != 15 {
		b524ComponentSelectorT.Fatalf("rows: %d", len(b524ComponentSelectorRows))
	}
	var b524ComponentSelectorWholeRows []b524ComponentSelectorJSON.RawMessage
	if b524ComponentSelectorError := b524ComponentSelectorJSON.Unmarshal(b524ComponentSelectorEntries, &b524ComponentSelectorWholeRows); b524ComponentSelectorError != nil {
		b524ComponentSelectorT.Fatal(b524ComponentSelectorError)
	}
	for b524ComponentSelectorIndex, b524ComponentSelectorExpectedRow := range b524ComponentSelectorWholeRows {
		b524ComponentSelectorActualRow, b524ComponentSelectorFound, b524ComponentSelectorError := B524ComponentSelectorLookup(b524ComponentSelectorRows[b524ComponentSelectorIndex].B524ComponentSelectorKey)
		if b524ComponentSelectorError != nil || !b524ComponentSelectorFound || !b524ComponentSelectorBytesPackage.Equal(b524ComponentSelectorActualRow, b524ComponentSelectorExpectedRow) {
			b524ComponentSelectorT.Fatal("whole row was not retained")
		}
	}
	b524ComponentSelectorSeen := make(map[string]bool)
	for _, b524ComponentSelectorFamily := range []string{"controller-700", "controller-720"} {
		b524ComponentSelectorSpecs := []struct {
			B524ComponentSelectorField    string
			B524ComponentSelectorGroup    uint8
			B524ComponentSelectorScope    string
			B524ComponentSelectorRegister uint16
		}{
			{"clock-date", 0, "00", 0x34}, {"clock-time", 0, "00", 0x35},
			{"system-holiday-start", 0, "00", 0x71}, {"system-holiday-end", 0, "00", 0x72},
			{"zone-holiday-start", 3, "dynamic", 3}, {"zone-holiday-end", 3, "dynamic", 4},
		}
		if b524ComponentSelectorFamily == "controller-720" {
			b524ComponentSelectorSpecs[2].B524ComponentSelectorRegister = 0xda
			b524ComponentSelectorSpecs[3].B524ComponentSelectorRegister = 0xdb
			b524ComponentSelectorSpecs = append(b524ComponentSelectorSpecs, struct {
				B524ComponentSelectorField    string
				B524ComponentSelectorGroup    uint8
				B524ComponentSelectorScope    string
				B524ComponentSelectorRegister uint16
			}{"zone-time-001e", 3, "dynamic", 0x1e}, struct {
				B524ComponentSelectorField    string
				B524ComponentSelectorGroup    uint8
				B524ComponentSelectorScope    string
				B524ComponentSelectorRegister uint16
			}{"zone-time-0020", 3, "dynamic", 0x20}, struct {
				B524ComponentSelectorField    string
				B524ComponentSelectorGroup    uint8
				B524ComponentSelectorScope    string
				B524ComponentSelectorRegister uint16
			}{"zone-time-0021", 3, "dynamic", 0x21})
		}
		for _, b524ComponentSelectorSpec := range b524ComponentSelectorSpecs {
			b524ComponentSelectorKey := b524ComponentSelectorBaseKey()
			b524ComponentSelectorKey.B524ComponentSelectorFamily = b524ComponentSelectorFamily
			b524ComponentSelectorKey.B524ComponentSelectorField = b524ComponentSelectorSpec.B524ComponentSelectorField
			b524ComponentSelectorKey.B524ComponentSelectorGroup = b524ComponentSelectorPointer(b524ComponentSelectorSpec.B524ComponentSelectorGroup)
			b524ComponentSelectorKey.B524ComponentSelectorInstanceScope = b524ComponentSelectorSpec.B524ComponentSelectorScope
			b524ComponentSelectorKey.B524ComponentSelectorRegister = b524ComponentSelectorPointer(b524ComponentSelectorSpec.B524ComponentSelectorRegister)
			b524ComponentSelectorRow := b524ComponentSelectorGet(b524ComponentSelectorT, b524ComponentSelectorKey)
			b524ComponentSelectorID := "ab32-" + b524ComponentSelectorFamily + "-" + b524ComponentSelectorSpec.B524ComponentSelectorField
			if string(b524ComponentSelectorRow["id"]) != `"`+b524ComponentSelectorID+`"` || b524ComponentSelectorSeen[b524ComponentSelectorID] {
				b524ComponentSelectorT.Fatal("identity mismatch")
			}
			b524ComponentSelectorSeen[b524ComponentSelectorID] = true
			var b524ComponentSelectorActual B524ComponentSelectorKey
			if b524ComponentSelectorError := b524ComponentSelectorJSON.Unmarshal(b524ComponentSelectorRow["key"], &b524ComponentSelectorActual); b524ComponentSelectorError != nil || !b524ComponentSelectorReflect.DeepEqual(b524ComponentSelectorActual, b524ComponentSelectorKey) {
				b524ComponentSelectorT.Fatal("whole key lost")
			}
		}
	}
}

func TestAb32EachLookupAxis(b524ComponentSelectorT *b524ComponentSelectorTesting.T) {
	b524ComponentSelectorChanges := []struct {
		B524ComponentSelectorName   string
		B524ComponentSelectorChange func(*B524ComponentSelectorKey)
		B524ComponentSelectorFound  bool
	}{
		{"ab32Family", func(b524ComponentSelectorKey *B524ComponentSelectorKey) {
			b524ComponentSelectorKey.B524ComponentSelectorFamily = "controller-720"
		}, true},
		{"ab32Variant", func(b524ComponentSelectorKey *B524ComponentSelectorKey) {
			b524ComponentSelectorKey.B524ComponentSelectorVariant = "other"
		}, false},
		{"ab32Primary", func(b524ComponentSelectorKey *B524ComponentSelectorKey) {
			b524ComponentSelectorKey.B524ComponentSelectorPrimary = 180
		}, false},
		{"ab32Secondary", func(b524ComponentSelectorKey *B524ComponentSelectorKey) {
			b524ComponentSelectorKey.B524ComponentSelectorSecondary = b524ComponentSelectorPointer(uint8(35))
		}, false},
		{"ab32Operation", func(b524ComponentSelectorKey *B524ComponentSelectorKey) {
			b524ComponentSelectorKey.B524ComponentSelectorOperation = b524ComponentSelectorPointer(uint8(6))
		}, false},
		{"ab32Direction", func(b524ComponentSelectorKey *B524ComponentSelectorKey) {
			b524ComponentSelectorKey.B524ComponentSelectorDirection = "write"
		}, false},
		{"ab32Field", func(b524ComponentSelectorKey *B524ComponentSelectorKey) {
			b524ComponentSelectorKey.B524ComponentSelectorField = "zone-name-suffix"
		}, false},
		{"ab32Alternative", func(b524ComponentSelectorKey *B524ComponentSelectorKey) {
			b524ComponentSelectorKey.B524ComponentSelectorAlternative = "other"
		}, false},
		{"ab32Group", func(b524ComponentSelectorKey *B524ComponentSelectorKey) {
			b524ComponentSelectorKey.B524ComponentSelectorGroup = b524ComponentSelectorPointer(uint8(3))
		}, false},
		{"ab32InstanceScope", func(b524ComponentSelectorKey *B524ComponentSelectorKey) {
			b524ComponentSelectorKey.B524ComponentSelectorInstanceScope = "dynamic"
		}, false},
		{"ab32Register", func(b524ComponentSelectorKey *B524ComponentSelectorKey) {
			b524ComponentSelectorKey.B524ComponentSelectorRegister = b524ComponentSelectorPointer(uint16(53))
		}, false},
	}
	for _, b524ComponentSelectorCase := range b524ComponentSelectorChanges {
		b524ComponentSelectorT.Run(b524ComponentSelectorCase.B524ComponentSelectorName, func(b524ComponentSelectorT *b524ComponentSelectorTesting.T) {
			b524ComponentSelectorKey := b524ComponentSelectorBaseKey()
			b524ComponentSelectorCase.B524ComponentSelectorChange(&b524ComponentSelectorKey)
			b524ComponentSelectorRaw, b524ComponentSelectorFound, b524ComponentSelectorError := B524ComponentSelectorLookup(b524ComponentSelectorKey)
			if b524ComponentSelectorError != nil || b524ComponentSelectorFound != b524ComponentSelectorCase.B524ComponentSelectorFound {
				b524ComponentSelectorT.Fatalf("exact axis: %v %v", b524ComponentSelectorFound, b524ComponentSelectorError)
			}
			if b524ComponentSelectorFound && !b524ComponentSelectorBytesPackage.Contains(b524ComponentSelectorRaw, []byte(`"ab32-controller-720-clock-date"`)) {
				b524ComponentSelectorT.Fatal("family fallback")
			}
		})
	}
}

func TestAb32AbsentAndZero(b524ComponentSelectorT *b524ComponentSelectorTesting.T) {
	for _, b524ComponentSelectorAxis := range []string{"secondary", "operation", "group", "register"} {
		b524ComponentSelectorT.Run("ab32"+b524ComponentSelectorAxis, func(b524ComponentSelectorT *b524ComponentSelectorTesting.T) {
			b524ComponentSelectorKey := b524ComponentSelectorBaseKey()
			b524ComponentSelectorSet := func(b524ComponentSelectorKey *B524ComponentSelectorKey, b524ComponentSelectorPresent bool) {
				switch b524ComponentSelectorAxis {
				case "secondary":
					b524ComponentSelectorKey.B524ComponentSelectorSecondary = nil
					if b524ComponentSelectorPresent {
						b524ComponentSelectorKey.B524ComponentSelectorSecondary = b524ComponentSelectorPointer(uint8(0))
					}
				case "operation":
					b524ComponentSelectorKey.B524ComponentSelectorOperation = nil
					if b524ComponentSelectorPresent {
						b524ComponentSelectorKey.B524ComponentSelectorOperation = b524ComponentSelectorPointer(uint8(0))
					}
				case "group":
					b524ComponentSelectorKey.B524ComponentSelectorGroup = nil
					if b524ComponentSelectorPresent {
						b524ComponentSelectorKey.B524ComponentSelectorGroup = b524ComponentSelectorPointer(uint8(0))
					}
				case "register":
					b524ComponentSelectorKey.B524ComponentSelectorRegister = nil
					if b524ComponentSelectorPresent {
						b524ComponentSelectorKey.B524ComponentSelectorRegister = b524ComponentSelectorPointer(uint16(0))
					}
				}
			}
			b524ComponentSelectorSet(&b524ComponentSelectorKey, false)
			b524ComponentSelectorOther := b524ComponentSelectorBaseKey()
			b524ComponentSelectorSet(&b524ComponentSelectorOther, true)
			b524ComponentSelectorData, b524ComponentSelectorError := b524ComponentSelectorJSON.Marshal([]any{map[string]any{"id": "ab32-absent", "key": b524ComponentSelectorKey}, map[string]any{"id": "ab32-zero", "key": b524ComponentSelectorOther}})
			if b524ComponentSelectorError != nil {
				b524ComponentSelectorT.Fatal(b524ComponentSelectorError)
			}
			for b524ComponentSelectorIndex, b524ComponentSelectorQuery := range []B524ComponentSelectorKey{b524ComponentSelectorKey, b524ComponentSelectorOther} {
				b524ComponentSelectorRaw, b524ComponentSelectorFound, b524ComponentSelectorError := b524ComponentSelectorLookupData(b524ComponentSelectorData, b524ComponentSelectorQuery)
				b524ComponentSelectorWant := []string{`"ab32-absent"`, `"ab32-zero"`}[b524ComponentSelectorIndex]
				if b524ComponentSelectorError != nil || !b524ComponentSelectorFound || !b524ComponentSelectorBytesPackage.Contains(b524ComponentSelectorRaw, []byte(b524ComponentSelectorWant)) {
					b524ComponentSelectorT.Fatalf("absent/zero: %s %v", b524ComponentSelectorRaw, b524ComponentSelectorError)
				}
			}
		})
	}
}

func TestAb32ValidationAndCopies(b524ComponentSelectorT *b524ComponentSelectorTesting.T) {
	b524ComponentSelectorKey := b524ComponentSelectorBaseKey()
	b524ComponentSelectorData, b524ComponentSelectorError := B524ComponentSelectorEntries()
	if b524ComponentSelectorError != nil {
		b524ComponentSelectorT.Fatal(b524ComponentSelectorError)
	}
	b524ComponentSelectorOriginal := append([]byte(nil), b524ComponentSelectorData...)
	b524ComponentSelectorRow, b524ComponentSelectorFound, b524ComponentSelectorError := B524ComponentSelectorLookup(b524ComponentSelectorKey)
	if b524ComponentSelectorError != nil || !b524ComponentSelectorFound {
		b524ComponentSelectorT.Fatal("lookup")
	}
	b524ComponentSelectorOriginalRow := append([]byte(nil), b524ComponentSelectorRow...)
	for b524ComponentSelectorIndex := range b524ComponentSelectorData {
		b524ComponentSelectorData[b524ComponentSelectorIndex] = 'x'
	}
	for b524ComponentSelectorIndex := range b524ComponentSelectorRow {
		b524ComponentSelectorRow[b524ComponentSelectorIndex] = 'x'
	}
	b524ComponentSelectorAgain, b524ComponentSelectorError := B524ComponentSelectorEntries()
	if b524ComponentSelectorError != nil || !b524ComponentSelectorBytesPackage.Equal(b524ComponentSelectorAgain, b524ComponentSelectorOriginal) {
		b524ComponentSelectorT.Fatal("entries alias")
	}
	b524ComponentSelectorAgain, b524ComponentSelectorFound, b524ComponentSelectorError = B524ComponentSelectorLookup(b524ComponentSelectorKey)
	if b524ComponentSelectorError != nil || !b524ComponentSelectorFound || !b524ComponentSelectorBytesPackage.Equal(b524ComponentSelectorAgain, b524ComponentSelectorOriginalRow) {
		b524ComponentSelectorT.Fatal("row alias")
	}
	b524ComponentSelectorDuplicate, _ := b524ComponentSelectorJSON.Marshal([]any{map[string]any{"id": "ab32-first", "key": b524ComponentSelectorKey}, map[string]any{"id": "ab32-second", "key": b524ComponentSelectorKey}})
	b524ComponentSelectorInvalid := [][]byte{[]byte(`{`), []byte(`null`), []byte(`{}`), []byte(`[null]`), []byte(`[{"id":"ab32-x","key":{}}]`), b524ComponentSelectorDuplicate}
	for b524ComponentSelectorIndex, b524ComponentSelectorBad := range b524ComponentSelectorInvalid {
		if _, _, b524ComponentSelectorError := b524ComponentSelectorValidate(b524ComponentSelectorBad); b524ComponentSelectorError == nil {
			b524ComponentSelectorT.Fatalf("invalid accepted %d", b524ComponentSelectorIndex)
		}
		if _, b524ComponentSelectorFound, b524ComponentSelectorError := b524ComponentSelectorLookupData(b524ComponentSelectorBad, b524ComponentSelectorKey); b524ComponentSelectorError == nil || b524ComponentSelectorFound {
			b524ComponentSelectorT.Fatalf("lookup invalid accepted %d", b524ComponentSelectorIndex)
		}
	}
}

func TestAb32ComponentMetadata(b524ComponentSelectorT *b524ComponentSelectorTesting.T) {
	b524ComponentSelectorData, b524ComponentSelectorError := B524ComponentSelectorEntries()
	if b524ComponentSelectorError != nil {
		b524ComponentSelectorT.Fatal(b524ComponentSelectorError)
	}
	var b524ComponentSelectorRows []struct {
		B524ComponentSelectorKey   B524ComponentSelectorKey `json:"key"`
		B524ComponentSelectorGrade string                   `json:"grade"`
		B524ComponentSelectorKind  string                   `json:"kind"`
		B524Components             []struct {
			B524ComponentSelectorName     string `json:"name"`
			B524ComponentSelectorPosition int    `json:"position"`
			B524ComponentSelectorEncoding string `json:"encoding"`
			B524ComponentSelectorScale    int    `json:"scale"`
			B524ComponentSelectorOffset   int    `json:"offset"`
			B524ComponentSelectorMinimum  *int   `json:"named_minimum"`
			B524ComponentSelectorMaximum  *int   `json:"named_maximum"`
			B524ComponentSelectorRaw      string `json:"raw_form"`
		} `json:"components"`
		B524ComponentSelectorRelations map[string]any    `json:"relations"`
		B524ComponentSelectorUnknown   map[string]string `json:"unknown"`
		B524ComponentSelectorRaw       string            `json:"raw_form"`
	}
	if b524ComponentSelectorError := b524ComponentSelectorJSON.Unmarshal(b524ComponentSelectorData, &b524ComponentSelectorRows); b524ComponentSelectorError != nil {
		b524ComponentSelectorT.Fatal(b524ComponentSelectorError)
	}
	for _, b524ComponentSelectorRow := range b524ComponentSelectorRows {
		b524ComponentSelectorLookup := b524ComponentSelectorGet(b524ComponentSelectorT, b524ComponentSelectorRow.B524ComponentSelectorKey)
		if string(b524ComponentSelectorLookup["grade"]) != `"Hypothesis"` || b524ComponentSelectorRow.B524ComponentSelectorGrade != "Hypothesis" {
			b524ComponentSelectorT.Fatal("grade")
		}
		b524ComponentSelectorWantRelations := map[string]any{"selector_context": "whole exact key", "paired_selector": nil, "component_order": "declared order", "failure_scope": "individual component", "retain_unaffected": true, "retain_raw": true, "missing_component": "explicit missing cell", "complete_timestamp": false, "freshness_inferred": false, "constraint_admission": false}
		if !b524ComponentSelectorReflect.DeepEqual(b524ComponentSelectorRow.B524ComponentSelectorRelations, b524ComponentSelectorWantRelations) {
			b524ComponentSelectorT.Fatal("component retention or unrelated pairing")
		}
		b524ComponentSelectorUnknownNames := []string{"native_support", "physical_effects", "persistence", "delivery", "century", "timezone", "joint_date_validity", "zone_name_selectors", "string_encoding", "string_width", "string_termination"}
		if len(b524ComponentSelectorRow.B524ComponentSelectorUnknown) != len(b524ComponentSelectorUnknownNames) {
			b524ComponentSelectorT.Fatal("unknown coverage")
		}
		for _, b524ComponentSelectorName := range b524ComponentSelectorUnknownNames {
			if b524ComponentSelectorRow.B524ComponentSelectorUnknown[b524ComponentSelectorName] != "Unknown" {
				b524ComponentSelectorT.Fatal("lost Unknown", b524ComponentSelectorName)
			}
		}
		if b524ComponentSelectorRow.B524ComponentSelectorRaw != "ordered component bytes with independent missing or failed cells" || len(b524ComponentSelectorRow.B524Components) != 3 {
			b524ComponentSelectorT.Fatal("raw or count")
		}
		b524ComponentSelectorNames := []string{"day", "month", "year"}
		b524ComponentSelectorBytes := []byte{0x1f, 0x0c, 0x63}
		b524ComponentSelectorExpected := []int{31, 12, 99}
		if b524ComponentSelectorRow.B524ComponentSelectorKind == "time" {
			b524ComponentSelectorNames = []string{"hour", "minute", "second"}
			b524ComponentSelectorBytes = []byte{0x17, 0x3b, 0x3a}
			b524ComponentSelectorExpected = []int{23, 59, 58}
		} else if b524ComponentSelectorRow.B524ComponentSelectorKind != "date" {
			b524ComponentSelectorT.Fatal("kind")
		}
		for b524ComponentSelectorIndex, b524ComponentSelectorComponent := range b524ComponentSelectorRow.B524Components {
			if b524ComponentSelectorComponent.B524ComponentSelectorName != b524ComponentSelectorNames[b524ComponentSelectorIndex] || b524ComponentSelectorComponent.B524ComponentSelectorPosition != b524ComponentSelectorIndex || b524ComponentSelectorComponent.B524ComponentSelectorEncoding != "u8" || b524ComponentSelectorComponent.B524ComponentSelectorScale != 1 || b524ComponentSelectorComponent.B524ComponentSelectorOffset != 0 || b524ComponentSelectorComponent.B524ComponentSelectorRaw != "individual byte" {
				b524ComponentSelectorT.Fatal("component representation")
			}
			if int(b524ComponentSelectorBytes[b524ComponentSelectorIndex])*b524ComponentSelectorComponent.B524ComponentSelectorScale+b524ComponentSelectorComponent.B524ComponentSelectorOffset != b524ComponentSelectorExpected[b524ComponentSelectorIndex] || int(byte(0x12))*b524ComponentSelectorComponent.B524ComponentSelectorScale+b524ComponentSelectorComponent.B524ComponentSelectorOffset != 18 {
				b524ComponentSelectorT.Fatal("plain byte examples")
			}
			if b524ComponentSelectorRow.B524ComponentSelectorKind == "date" {
				if b524ComponentSelectorComponent.B524ComponentSelectorMinimum == nil || b524ComponentSelectorComponent.B524ComponentSelectorMaximum == nil || *b524ComponentSelectorComponent.B524ComponentSelectorMinimum != []int{1, 1, 0}[b524ComponentSelectorIndex] || *b524ComponentSelectorComponent.B524ComponentSelectorMaximum != []int{31, 12, 99}[b524ComponentSelectorIndex] {
					b524ComponentSelectorT.Fatal("named domains")
				}
			}
		}
	}
}

func TestAb32NoSelectorFallback(b524ComponentSelectorT *b524ComponentSelectorTesting.T) {
	for _, b524ComponentSelectorCase := range []struct {
		B524ComponentSelectorFamily   string
		B524ComponentSelectorField    string
		B524ComponentSelectorRegister uint16
		B524ComponentSelectorGroup    uint8
		B524ComponentSelectorScope    string
	}{
		{"controller-700", "system-holiday-start", 0xda, 0, "00"}, {"controller-720", "system-holiday-start", 0x71, 0, "00"},
		{"controller-700", "system-holiday-end", 0xdb, 0, "00"}, {"controller-720", "system-holiday-end", 0x72, 0, "00"},
		{"controller-700", "zone-time-001e", 0x1e, 3, "dynamic"}, {"controller-700", "zone-time-0020", 0x20, 3, "dynamic"}, {"controller-700", "zone-time-0021", 0x21, 3, "dynamic"},
		{"controller-720", "zone-name-suffix", 3, 3, "dynamic"},
	} {
		b524ComponentSelectorKey := b524ComponentSelectorBaseKey()
		b524ComponentSelectorKey.B524ComponentSelectorFamily = b524ComponentSelectorCase.B524ComponentSelectorFamily
		b524ComponentSelectorKey.B524ComponentSelectorField = b524ComponentSelectorCase.B524ComponentSelectorField
		b524ComponentSelectorKey.B524ComponentSelectorRegister = b524ComponentSelectorPointer(b524ComponentSelectorCase.B524ComponentSelectorRegister)
		b524ComponentSelectorKey.B524ComponentSelectorGroup = b524ComponentSelectorPointer(b524ComponentSelectorCase.B524ComponentSelectorGroup)
		b524ComponentSelectorKey.B524ComponentSelectorInstanceScope = b524ComponentSelectorCase.B524ComponentSelectorScope
		if _, b524ComponentSelectorFound, b524ComponentSelectorError := B524ComponentSelectorLookup(b524ComponentSelectorKey); b524ComponentSelectorFound || b524ComponentSelectorError != nil {
			b524ComponentSelectorT.Fatal("forbidden selector fallback")
		}
	}
}

func TestAb32PartialCells(b524ComponentSelectorT *b524ComponentSelectorTesting.T) {
	b524ComponentSelectorRow := b524ComponentSelectorGet(b524ComponentSelectorT, b524ComponentSelectorBaseKey())
	var b524ComponentSelectorComponents []struct {
		B524ComponentSelectorName     string `json:"name"`
		B524ComponentSelectorPosition int    `json:"position"`
	}
	if b524ComponentSelectorError := b524ComponentSelectorJSON.Unmarshal(b524ComponentSelectorRow["components"], &b524ComponentSelectorComponents); b524ComponentSelectorError != nil {
		b524ComponentSelectorT.Fatal(b524ComponentSelectorError)
	}
	b524ComponentSelectorCells := []*uint8{b524ComponentSelectorPointer(uint8(31)), nil, b524ComponentSelectorPointer(uint8(99))}
	b524ComponentSelectorExpected := map[string]*uint8{"day": b524ComponentSelectorPointer(uint8(31)), "month": nil, "year": b524ComponentSelectorPointer(uint8(99))}
	for _, b524ComponentSelectorComponent := range b524ComponentSelectorComponents {
		b524ComponentSelectorWant, b524ComponentSelectorPresent := b524ComponentSelectorExpected[b524ComponentSelectorComponent.B524ComponentSelectorName]
		if !b524ComponentSelectorPresent || b524ComponentSelectorComponent.B524ComponentSelectorPosition < 0 || b524ComponentSelectorComponent.B524ComponentSelectorPosition >= len(b524ComponentSelectorCells) || !b524ComponentSelectorReflect.DeepEqual(b524ComponentSelectorCells[b524ComponentSelectorComponent.B524ComponentSelectorPosition], b524ComponentSelectorWant) {
			b524ComponentSelectorT.Fatal("partial raw component mapping")
		}
	}
	var b524ComponentSelectorRelations map[string]any
	if b524ComponentSelectorError := b524ComponentSelectorJSON.Unmarshal(b524ComponentSelectorRow["relations"], &b524ComponentSelectorRelations); b524ComponentSelectorError != nil {
		b524ComponentSelectorT.Fatal(b524ComponentSelectorError)
	}
	if b524ComponentSelectorRelations["retain_unaffected"] != true || b524ComponentSelectorRelations["retain_raw"] != true || b524ComponentSelectorRelations["failure_scope"] != "individual component" || b524ComponentSelectorRelations["missing_component"] != "explicit missing cell" || b524ComponentSelectorRelations["complete_timestamp"] != false || b524ComponentSelectorRelations["freshness_inferred"] != false {
		b524ComponentSelectorT.Fatal("partial failure policy")
	}
	if _, b524ComponentSelectorPresent := b524ComponentSelectorRelations["fresh"]; b524ComponentSelectorPresent {
		b524ComponentSelectorT.Fatal("invented freshness")
	}
}
