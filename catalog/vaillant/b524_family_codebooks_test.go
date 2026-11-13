package vaillant

import (
	b524FamilyCodebookBytes "bytes"
	b524FamilyCodebookJSON "encoding/json"
	b524FamilyCodebookFmt "fmt"
	b524FamilyCodebookReflect "reflect"
	b524FamilyCodebookStrconv "strconv"
	b524FamilyCodebookStrings "strings"
	b524FamilyCodebookTesting "testing"
)

type b524FamilyCodebookMeaning struct {
	B524FamilyCodebookValue int    `json:"value"`
	B524FamilyCodebookLabel string `json:"label"`
}

type b524FamilyCodebookRow struct {
	B524FamilyCodebookID             string                               `json:"id"`
	B524FamilyCodebookKey            B524FamilyCodebookKey                `json:"key"`
	B524FamilyCodebookGrade          string                               `json:"grade"`
	B524FamilyCodebookRepresentation string                               `json:"representation_grade"`
	B524FamilyCodebookEnumeration    map[string]b524FamilyCodebookMeaning `json:"enumeration"`
	B524FamilyCodebookRawForm        struct {
		B524FamilyCodebookEncoding   string `json:"encoding"`
		B524FamilyCodebookWidth      int    `json:"width_bytes"`
		B524FamilyCodebookSignedness string `json:"signedness"`
		B524FamilyCodebookPreserved  bool   `json:"preserved"`
	} `json:"raw_form"`
	B524FamilyCodebookFallback  string            `json:"raw_fallback"`
	B524FamilyCodebookRelations map[string]string `json:"relations"`
	B524FamilyCodebookUnknown   map[string]string `json:"unknown"`
}

func b524FamilyCodebookPointer[b524FamilyCodebookType any](b524FamilyCodebookValue b524FamilyCodebookType) *b524FamilyCodebookType {
	return &b524FamilyCodebookValue
}

func b524FamilyCodebookKey(b524FamilyCodebookFamily string, b524FamilyCodebookGroup uint8, b524FamilyCodebookRegister uint16, b524FamilyCodebookScope, b524FamilyCodebookAlternative string) B524FamilyCodebookKey {
	return B524FamilyCodebookKey{B524Family: b524FamilyCodebookFamily, B524FamilyCodebookVariant: "none", B524FamilyCodebookPrimary: 0xB5, B524FamilyCodebookSecondary: b524FamilyCodebookPointer(uint8(0x24)), B524FamilyCodebookOperation: b524FamilyCodebookPointer(uint8(2)), B524FamilyCodebookDirection: "read", B524FamilyCodebookField: "value", B524FamilyCodebookAlternative: b524FamilyCodebookAlternative, B524FamilyCodebookGroup: b524FamilyCodebookPointer(b524FamilyCodebookGroup), B524FamilyCodebookInstanceScope: b524FamilyCodebookScope, B524FamilyCodebookRegister: b524FamilyCodebookPointer(b524FamilyCodebookRegister)}
}

func b524FamilyCodebookGet(b524FamilyCodebookT *b524FamilyCodebookTesting.T, b524FamilyCodebookKey B524FamilyCodebookKey) b524FamilyCodebookRow {
	b524FamilyCodebookT.Helper()
	b524FamilyCodebookRaw, b524FamilyCodebookFound, b524FamilyCodebookError := B524FamilyCodebookLookup(b524FamilyCodebookKey)
	if b524FamilyCodebookError != nil || !b524FamilyCodebookFound {
		b524FamilyCodebookT.Fatalf("ab30 valid control: %v %v", b524FamilyCodebookFound, b524FamilyCodebookError)
	}
	var b524FamilyCodebookResult b524FamilyCodebookRow
	if b524FamilyCodebookError := b524FamilyCodebookJSON.Unmarshal(b524FamilyCodebookRaw, &b524FamilyCodebookResult); b524FamilyCodebookError != nil {
		b524FamilyCodebookT.Fatal(b524FamilyCodebookError)
	}
	if !b524FamilyCodebookReflect.DeepEqual(b524FamilyCodebookResult.B524FamilyCodebookKey, b524FamilyCodebookKey) {
		b524FamilyCodebookT.Fatal("ab30 returned a different whole row")
	}
	return b524FamilyCodebookResult
}

func TestAb30Codebooks(b524FamilyCodebookT *b524FamilyCodebookTesting.T) {
	b524FamilyCodebookRaw, b524FamilyCodebookError := B524FamilyCodebookEntries()
	if b524FamilyCodebookError != nil {
		b524FamilyCodebookT.Fatal(b524FamilyCodebookError)
	}
	var b524FamilyCodebookRows []b524FamilyCodebookRow
	if b524FamilyCodebookError := b524FamilyCodebookJSON.Unmarshal(b524FamilyCodebookRaw, &b524FamilyCodebookRows); b524FamilyCodebookError != nil {
		b524FamilyCodebookT.Fatal(b524FamilyCodebookError)
	}
	if len(b524FamilyCodebookRows) != 16 {
		b524FamilyCodebookT.Fatalf("ab30 finite row count: %d", len(b524FamilyCodebookRows))
	}
	b524FamilyCodebookSeen := map[string]bool{}
	for _, b524FamilyCodebookFamily := range []string{"controller-700", "controller-720"} {
		for _, b524FamilyCodebookCase := range []struct {
			b524FamilyCodebookGroup                                                                                          uint8
			b524FamilyCodebookRegister                                                                                       uint16
			b524FamilyCodebookScope, b524FamilyCodebookAlternative, b524FamilyCodebookLabels700, b524FamilyCodebookLabels720 string
		}{
			{1, 0x0003, "00", "u16", "0:off,1:automatic,2:day", "0:off,1:time-controlled,2:manual"},
			{2, 0x0002, "dynamic", "u16", "0:deactivated,1:heating,2:fixed-flow,3:DHW,4:return-flow,5:pool", "0:deactivated,1:heating,2:fixed-flow,3:DHW,4:return-flow"},
			{3, 0x0001, "dynamic", "u16", "0:off,1:automatic,2:day", "0:off,1:time-controlled,2:manual"},
			{3, 0x0006, "dynamic", "u16", "0:off,1:automatic,2:day,3:setback", "0:off,1:time-controlled,2:manual"},
			{3, 0x0013, "dynamic", "u16", "0:none,1:central,2:remote-circuit-one,3:remote-circuit-two,4:remote-circuit-three,5:remote-circuit-four,6:remote-circuit-five,7:remote-circuit-six,8:remote-circuit-seven,9:remote-circuit-eight", "0:none,1:central,2:remote-circuit-one,3:remote-circuit-two,4:remote-circuit-three,5:remote-circuit-four"},
			{3, 0x001B, "dynamic", "u8", "0:no-quick-mode,1:ventilation-boost,3:party,4:quick-veto,5:one-day-holiday,6:one-day-bank-holiday,7:holiday,8:bank-holiday,10:system-off", "0:no-special-function,1:ventilation-boost,4:quick-veto,7:holiday,10:system-off"},
			{3, 0x001B, "dynamic", "u16", "0:no-quick-mode,1:ventilation-boost,3:party,4:quick-veto,5:one-day-holiday,6:one-day-bank-holiday,7:holiday,8:bank-holiday,10:system-off", "0:no-special-function,1:ventilation-boost,4:quick-veto,7:holiday,10:system-off"},
			{3, 0x001C, "dynamic", "signed8", "255:unbound,0:control-circuit-one,1:control-circuit-two,2:control-circuit-three", "255:unbound,0:control-circuit-one,1:control-circuit-two,2:control-circuit-three,3:control-circuit-four,4:control-circuit-five,5:control-circuit-six,6:control-circuit-seven,7:control-circuit-eight,8:control-circuit-nine"},
		} {
			b524FamilyCodebookT.Run(b524FamilyCodebookFmt.Sprintf("ab30_%s_%d_%d_%s", b524FamilyCodebookFamily, b524FamilyCodebookCase.b524FamilyCodebookGroup, b524FamilyCodebookCase.b524FamilyCodebookRegister, b524FamilyCodebookCase.b524FamilyCodebookAlternative), func(b524FamilyCodebookT *b524FamilyCodebookTesting.T) {
				b524FamilyCodebookGot := b524FamilyCodebookGet(b524FamilyCodebookT, b524FamilyCodebookKey(b524FamilyCodebookFamily, b524FamilyCodebookCase.b524FamilyCodebookGroup, b524FamilyCodebookCase.b524FamilyCodebookRegister, b524FamilyCodebookCase.b524FamilyCodebookScope, b524FamilyCodebookCase.b524FamilyCodebookAlternative))
				if !b524FamilyCodebookStrings.HasPrefix(b524FamilyCodebookGot.B524FamilyCodebookID, "ab30_") || b524FamilyCodebookSeen[b524FamilyCodebookGot.B524FamilyCodebookID] {
					b524FamilyCodebookT.Fatal("ab30 row IDs must be unique and prefixed")
				}
				b524FamilyCodebookSeen[b524FamilyCodebookGot.B524FamilyCodebookID] = true
				b524FamilyCodebookLabels := b524FamilyCodebookCase.b524FamilyCodebookLabels700
				if b524FamilyCodebookFamily == "controller-720" {
					b524FamilyCodebookLabels = b524FamilyCodebookCase.b524FamilyCodebookLabels720
				}
				b524FamilyCodebookExpected := map[string]b524FamilyCodebookMeaning{}
				for _, b524FamilyCodebookPair := range b524FamilyCodebookStrings.Split(b524FamilyCodebookLabels, ",") {
					b524FamilyCodebookParts := b524FamilyCodebookStrings.SplitN(b524FamilyCodebookPair, ":", 2)
					b524FamilyCodebookValue, b524FamilyCodebookError := b524FamilyCodebookStrconv.Atoi(b524FamilyCodebookParts[0])
					if b524FamilyCodebookError != nil {
						b524FamilyCodebookT.Fatal(b524FamilyCodebookError)
					}
					if b524FamilyCodebookCase.b524FamilyCodebookAlternative == "signed8" && b524FamilyCodebookValue == 255 {
						b524FamilyCodebookValue = -1
					}
					b524FamilyCodebookExpected[b524FamilyCodebookParts[0]] = b524FamilyCodebookMeaning{b524FamilyCodebookValue, b524FamilyCodebookParts[1]}
				}
				if !b524FamilyCodebookReflect.DeepEqual(b524FamilyCodebookGot.B524FamilyCodebookEnumeration, b524FamilyCodebookExpected) {
					b524FamilyCodebookT.Fatalf("ab30 independent exact codebook: got %v want %v", b524FamilyCodebookGot.B524FamilyCodebookEnumeration, b524FamilyCodebookExpected)
				}
				if b524FamilyCodebookGot.B524FamilyCodebookGrade != "Hypothesis" || b524FamilyCodebookGot.B524FamilyCodebookRepresentation != "Proven" {
					b524FamilyCodebookT.Fatal("ab30 grades")
				}
				b524FamilyCodebookWidth, b524FamilyCodebookSignedness := 1, "unsigned"
				if b524FamilyCodebookCase.b524FamilyCodebookAlternative == "u16" {
					b524FamilyCodebookWidth = 2
				}
				if b524FamilyCodebookCase.b524FamilyCodebookAlternative == "signed8" {
					b524FamilyCodebookSignedness = "signed"
				}
				if b524FamilyCodebookGot.B524FamilyCodebookRawForm.B524FamilyCodebookEncoding != b524FamilyCodebookCase.b524FamilyCodebookAlternative || b524FamilyCodebookGot.B524FamilyCodebookRawForm.B524FamilyCodebookWidth != b524FamilyCodebookWidth || b524FamilyCodebookGot.B524FamilyCodebookRawForm.B524FamilyCodebookSignedness != b524FamilyCodebookSignedness || !b524FamilyCodebookGot.B524FamilyCodebookRawForm.B524FamilyCodebookPreserved {
					b524FamilyCodebookT.Fatal("ab30 raw form")
				}
				if b524FamilyCodebookGot.B524FamilyCodebookFallback != "Preserve raw bytes; unlisted codes and unknown families have Unknown meaning." {
					b524FamilyCodebookT.Fatal("ab30 raw fallback")
				}
				if !b524FamilyCodebookReflect.DeepEqual(b524FamilyCodebookGot.B524FamilyCodebookRelations, map[string]string{"identity_grade": "Proven", "candidate_opcode": "OP02", "instance_scope_role": "metadata-only", "width_selection": "explicit alternative; reply length does not select", "family_selection": "explicit; no auto-detection"}) {
					b524FamilyCodebookT.Fatal("ab30 relations must not infer installed cardinality or width")
				}
				b524FamilyCodebookUnknown := map[string]string{}
				for _, b524FamilyCodebookName := range []string{"native_support", "physical_effects", "persistence", "delivery", "installed_cardinality", "present_slots", "freshness", "unlisted_codes"} {
					b524FamilyCodebookUnknown[b524FamilyCodebookName] = "Unknown"
				}
				if !b524FamilyCodebookReflect.DeepEqual(b524FamilyCodebookGot.B524FamilyCodebookUnknown, b524FamilyCodebookUnknown) {
					b524FamilyCodebookT.Fatal("ab30 qualification remains Unknown")
				}
			})
		}
	}
	for _, b524FamilyCodebookRow := range b524FamilyCodebookRows {
		if !b524FamilyCodebookSeen[b524FamilyCodebookRow.B524FamilyCodebookID] {
			b524FamilyCodebookT.Fatal("ab30 unexpected production row")
		}
	}
}

func TestAb30EveryIdentityAxis(b524FamilyCodebookT *b524FamilyCodebookTesting.T) {
	b524FamilyCodebookControl := b524FamilyCodebookKey("controller-700", 3, 0x001B, "dynamic", "u8")
	b524FamilyCodebookGet(b524FamilyCodebookT, b524FamilyCodebookControl)
	for b524FamilyCodebookAxis := 0; b524FamilyCodebookAxis < b524FamilyCodebookReflect.TypeOf(b524FamilyCodebookControl).NumField(); b524FamilyCodebookAxis++ {
		b524FamilyCodebookName := b524FamilyCodebookReflect.TypeOf(b524FamilyCodebookControl).Field(b524FamilyCodebookAxis).Name
		b524FamilyCodebookT.Run("ab30_"+b524FamilyCodebookName, func(b524FamilyCodebookT *b524FamilyCodebookTesting.T) {
			b524FamilyCodebookChanged := b524FamilyCodebookControl
			b524FamilyCodebookSlot := b524FamilyCodebookReflect.ValueOf(&b524FamilyCodebookChanged).Elem().Field(b524FamilyCodebookAxis)
			switch b524FamilyCodebookSlot.Kind() {
			case b524FamilyCodebookReflect.String:
				b524FamilyCodebookSlot.SetString("ab30_other")
			case b524FamilyCodebookReflect.Uint8:
				b524FamilyCodebookSlot.SetUint(0)
			case b524FamilyCodebookReflect.Ptr:
				b524FamilyCodebookValue := b524FamilyCodebookReflect.New(b524FamilyCodebookSlot.Type().Elem())
				b524FamilyCodebookValue.Elem().SetUint(6)
				b524FamilyCodebookSlot.Set(b524FamilyCodebookValue)
			default:
				b524FamilyCodebookT.Fatal("ab30 unhandled key axis")
			}
			b524FamilyCodebookRaw, b524FamilyCodebookFound, b524FamilyCodebookError := B524FamilyCodebookLookup(b524FamilyCodebookChanged)
			if b524FamilyCodebookError != nil || b524FamilyCodebookFound || b524FamilyCodebookRaw != nil {
				b524FamilyCodebookT.Fatalf("ab30 changed axis inherited labels: %s %v %v", b524FamilyCodebookName, b524FamilyCodebookFound, b524FamilyCodebookError)
			}
			b524FamilyCodebookFirst := b524FamilyCodebookGet(b524FamilyCodebookT, b524FamilyCodebookControl)
			b524FamilyCodebookSecond := b524FamilyCodebookFirst
			b524FamilyCodebookSecond.B524FamilyCodebookKey = b524FamilyCodebookChanged
			b524FamilyCodebookSecond.B524FamilyCodebookID = "ab30_decoy"
			b524FamilyCodebookData, b524FamilyCodebookError := b524FamilyCodebookJSON.Marshal([]b524FamilyCodebookRow{b524FamilyCodebookFirst, b524FamilyCodebookSecond})
			if b524FamilyCodebookError != nil {
				b524FamilyCodebookT.Fatal(b524FamilyCodebookError)
			}
			b524FamilyCodebookRows, _, b524FamilyCodebookError := b524FamilyCodebookValidate(b524FamilyCodebookData)
			if b524FamilyCodebookError != nil || len(b524FamilyCodebookRows) != 2 {
				b524FamilyCodebookT.Fatalf("ab30 valid distinct pair: %v", b524FamilyCodebookError)
			}
			b524FamilyCodebookRaw, b524FamilyCodebookFound, b524FamilyCodebookError = b524FamilyCodebookLookup(b524FamilyCodebookData, b524FamilyCodebookChanged)
			var b524FamilyCodebookGot b524FamilyCodebookRow
			if b524FamilyCodebookError != nil || !b524FamilyCodebookFound {
				b524FamilyCodebookT.Fatalf("ab30 decoy control: %v", b524FamilyCodebookError)
			}
			if b524FamilyCodebookError := b524FamilyCodebookJSON.Unmarshal(b524FamilyCodebookRaw, &b524FamilyCodebookGot); b524FamilyCodebookError != nil || b524FamilyCodebookGot.B524FamilyCodebookID != "ab30_decoy" {
				b524FamilyCodebookT.Fatalf("ab30 whole decoy row: %v", b524FamilyCodebookError)
			}
		})
	}
}

func TestAb30OptionalAbsentAndZero(b524FamilyCodebookT *b524FamilyCodebookTesting.T) {
	b524FamilyCodebookControl := b524FamilyCodebookGet(b524FamilyCodebookT, b524FamilyCodebookKey("controller-700", 3, 27, "dynamic", "u8"))
	for _, b524FamilyCodebookName := range []string{"B524FamilyCodebookSecondary", "B524FamilyCodebookOperation", "B524FamilyCodebookGroup", "B524FamilyCodebookRegister"} {
		b524FamilyCodebookT.Run("ab30_"+b524FamilyCodebookName, func(b524FamilyCodebookT *b524FamilyCodebookTesting.T) {
			b524FamilyCodebookAbsent, b524FamilyCodebookZero := b524FamilyCodebookControl, b524FamilyCodebookControl
			b524FamilyCodebookAbsent.B524FamilyCodebookID = "ab30_absent"
			b524FamilyCodebookZero.B524FamilyCodebookID = "ab30_zero"
			b524FamilyCodebookReflect.ValueOf(&b524FamilyCodebookAbsent.B524FamilyCodebookKey).Elem().FieldByName(b524FamilyCodebookName).SetZero()
			b524FamilyCodebookSlot := b524FamilyCodebookReflect.ValueOf(&b524FamilyCodebookZero.B524FamilyCodebookKey).Elem().FieldByName(b524FamilyCodebookName)
			b524FamilyCodebookSlot.Set(b524FamilyCodebookReflect.New(b524FamilyCodebookSlot.Type().Elem()))
			b524FamilyCodebookData, b524FamilyCodebookError := b524FamilyCodebookJSON.Marshal([]b524FamilyCodebookRow{b524FamilyCodebookAbsent, b524FamilyCodebookZero})
			if b524FamilyCodebookError != nil {
				b524FamilyCodebookT.Fatal(b524FamilyCodebookError)
			}
			for _, b524FamilyCodebookExpected := range []b524FamilyCodebookRow{b524FamilyCodebookAbsent, b524FamilyCodebookZero} {
				b524FamilyCodebookRaw, b524FamilyCodebookFound, b524FamilyCodebookError := b524FamilyCodebookLookup(b524FamilyCodebookData, b524FamilyCodebookExpected.B524FamilyCodebookKey)
				var b524FamilyCodebookGot b524FamilyCodebookRow
				if b524FamilyCodebookError != nil || !b524FamilyCodebookFound {
					b524FamilyCodebookT.Fatalf("ab30 absent/zero control: %v", b524FamilyCodebookError)
				}
				if b524FamilyCodebookError := b524FamilyCodebookJSON.Unmarshal(b524FamilyCodebookRaw, &b524FamilyCodebookGot); b524FamilyCodebookError != nil || b524FamilyCodebookGot.B524FamilyCodebookID != b524FamilyCodebookExpected.B524FamilyCodebookID {
					b524FamilyCodebookT.Fatalf("ab30 absent and zero collapsed: %v", b524FamilyCodebookError)
				}
				b524FamilyCodebookRaw, b524FamilyCodebookFound, b524FamilyCodebookError = B524FamilyCodebookLookup(b524FamilyCodebookExpected.B524FamilyCodebookKey)
				if b524FamilyCodebookError != nil || b524FamilyCodebookFound || b524FamilyCodebookRaw != nil {
					b524FamilyCodebookT.Fatal("ab30 missing optional identity inherited production labels")
				}
			}
		})
	}
}

func TestAb30NamespaceAndRawFallback(b524FamilyCodebookT *b524FamilyCodebookTesting.T) {
	for _, b524FamilyCodebookFamily := range []string{"controller-700", "controller-720"} {
		b524FamilyCodebookSigned := b524FamilyCodebookGet(b524FamilyCodebookT, b524FamilyCodebookKey(b524FamilyCodebookFamily, 3, 28, "dynamic", "signed8"))
		if b524FamilyCodebookSigned.B524FamilyCodebookEnumeration["255"] != (b524FamilyCodebookMeaning{-1, "unbound"}) {
			b524FamilyCodebookT.Fatal("ab30 FF candidate signed8")
		}
		b524FamilyCodebookUnsigned := b524FamilyCodebookGet(b524FamilyCodebookT, b524FamilyCodebookKey(b524FamilyCodebookFamily, 3, 6, "dynamic", "u16"))
		for _, b524FamilyCodebookRaw := range []string{"255", "65535", "-1"} {
			if _, b524FamilyCodebookFound := b524FamilyCodebookUnsigned.B524FamilyCodebookEnumeration[b524FamilyCodebookRaw]; b524FamilyCodebookFound {
				b524FamilyCodebookT.Fatal("ab30 unrelated u16 sign extension")
			}
		}
		for _, b524FamilyCodebookAlternative := range []string{"u8", "u16"} {
			b524FamilyCodebookGot := b524FamilyCodebookGet(b524FamilyCodebookT, b524FamilyCodebookKey(b524FamilyCodebookFamily, 3, 27, "dynamic", b524FamilyCodebookAlternative))
			if _, b524FamilyCodebookFound := b524FamilyCodebookGot.B524FamilyCodebookEnumeration["2"]; b524FamilyCodebookFound {
				b524FamilyCodebookT.Fatal("ab30 unlisted raw code labelled")
			}
		}
		for _, b524FamilyCodebookChanged := range []B524FamilyCodebookKey{
			b524FamilyCodebookKey("ab30_unknown-family", 2, 2, "dynamic", "u16"),
			b524FamilyCodebookKey(b524FamilyCodebookFamily, 3, 28, "dynamic", "u16"),
			b524FamilyCodebookKey(b524FamilyCodebookFamily, 3, 27, "00", "u8"),
			b524FamilyCodebookKey(b524FamilyCodebookFamily, 1, 3, "dynamic", "u16"),
			b524FamilyCodebookKey(b524FamilyCodebookFamily, 3, 27, "dynamic", "none"),
		} {
			b524FamilyCodebookRaw, b524FamilyCodebookFound, b524FamilyCodebookError := B524FamilyCodebookLookup(b524FamilyCodebookChanged)
			if b524FamilyCodebookError != nil || b524FamilyCodebookFound || b524FamilyCodebookRaw != nil {
				b524FamilyCodebookT.Fatal("ab30 cross-namespace fallback")
			}
		}
		b524FamilyCodebookOpcode := b524FamilyCodebookKey(b524FamilyCodebookFamily, 3, 27, "dynamic", "u8")
		b524FamilyCodebookOpcode.B524FamilyCodebookOperation = b524FamilyCodebookPointer(uint8(6))
		b524FamilyCodebookRaw, b524FamilyCodebookFound, b524FamilyCodebookError := B524FamilyCodebookLookup(b524FamilyCodebookOpcode)
		if b524FamilyCodebookError != nil || b524FamilyCodebookFound || b524FamilyCodebookRaw != nil {
			b524FamilyCodebookT.Fatal("ab30 OP06 inherited OP02")
		}
	}
}

func TestAb30InvalidEmbeddedData(b524FamilyCodebookT *b524FamilyCodebookTesting.T) {
	b524FamilyCodebookKey := b524FamilyCodebookKey("controller-700", 3, 27, "dynamic", "u8")
	b524FamilyCodebookControlRow := b524FamilyCodebookGet(b524FamilyCodebookT, b524FamilyCodebookKey)
	b524FamilyCodebookDuplicate, b524FamilyCodebookError := b524FamilyCodebookJSON.Marshal([]b524FamilyCodebookRow{b524FamilyCodebookControlRow, b524FamilyCodebookControlRow})
	if b524FamilyCodebookError != nil {
		b524FamilyCodebookT.Fatal(b524FamilyCodebookError)
	}
	b524FamilyCodebookCases := [][]byte{[]byte("["), []byte("null"), []byte("{}"), []byte("[null]"), []byte("[{}]"), b524FamilyCodebookDuplicate}
	b524FamilyCodebookValid, b524FamilyCodebookError := b524FamilyCodebookJSON.Marshal([]b524FamilyCodebookRow{b524FamilyCodebookControlRow})
	if b524FamilyCodebookError != nil {
		b524FamilyCodebookT.Fatal(b524FamilyCodebookError)
	}
	var b524FamilyCodebookObjects []map[string]b524FamilyCodebookJSON.RawMessage
	if b524FamilyCodebookError := b524FamilyCodebookJSON.Unmarshal(b524FamilyCodebookValid, &b524FamilyCodebookObjects); b524FamilyCodebookError != nil {
		b524FamilyCodebookT.Fatal(b524FamilyCodebookError)
	}
	var b524FamilyCodebookIdentity map[string]b524FamilyCodebookJSON.RawMessage
	if b524FamilyCodebookError := b524FamilyCodebookJSON.Unmarshal(b524FamilyCodebookObjects[0]["key"], &b524FamilyCodebookIdentity); b524FamilyCodebookError != nil {
		b524FamilyCodebookT.Fatal(b524FamilyCodebookError)
	}
	for b524FamilyCodebookAxis := range b524FamilyCodebookIdentity {
		for _, b524FamilyCodebookReplacement := range []b524FamilyCodebookJSON.RawMessage{nil, []byte("true")} {
			b524FamilyCodebookCopy := map[string]b524FamilyCodebookJSON.RawMessage{}
			for b524FamilyCodebookName, b524FamilyCodebookValue := range b524FamilyCodebookIdentity {
				b524FamilyCodebookCopy[b524FamilyCodebookName] = b524FamilyCodebookValue
			}
			if b524FamilyCodebookReplacement == nil {
				delete(b524FamilyCodebookCopy, b524FamilyCodebookAxis)
			} else {
				b524FamilyCodebookCopy[b524FamilyCodebookAxis] = b524FamilyCodebookReplacement
			}
			b524FamilyCodebookKeyData, b524FamilyCodebookError := b524FamilyCodebookJSON.Marshal(b524FamilyCodebookCopy)
			if b524FamilyCodebookError != nil {
				b524FamilyCodebookT.Fatal(b524FamilyCodebookError)
			}
			b524FamilyCodebookInvalidData, b524FamilyCodebookError := b524FamilyCodebookJSON.Marshal([]map[string]b524FamilyCodebookJSON.RawMessage{{"key": b524FamilyCodebookKeyData}})
			if b524FamilyCodebookError != nil {
				b524FamilyCodebookT.Fatal(b524FamilyCodebookError)
			}
			b524FamilyCodebookCases = append(b524FamilyCodebookCases, b524FamilyCodebookInvalidData)
		}
	}
	for b524FamilyCodebookIndex, b524FamilyCodebookData := range b524FamilyCodebookCases {
		b524FamilyCodebookT.Run(b524FamilyCodebookFmt.Sprintf("ab30_invalid_%d", b524FamilyCodebookIndex), func(b524FamilyCodebookT *b524FamilyCodebookTesting.T) {
			if _, _, b524FamilyCodebookError := b524FamilyCodebookValidate(b524FamilyCodebookData); b524FamilyCodebookError == nil {
				b524FamilyCodebookT.Fatal("ab30 invalid data admitted")
			}
			b524FamilyCodebookRaw, b524FamilyCodebookFound, b524FamilyCodebookError := b524FamilyCodebookLookup(b524FamilyCodebookData, b524FamilyCodebookKey)
			if b524FamilyCodebookError == nil || b524FamilyCodebookFound || b524FamilyCodebookRaw != nil {
				b524FamilyCodebookT.Fatal("ab30 invalid lookup must fail closed")
			}
		})
	}
}

func TestAb30IndependentBuffers(b524FamilyCodebookT *b524FamilyCodebookTesting.T) {
	b524FamilyCodebookKey := b524FamilyCodebookKey("controller-700", 3, 27, "dynamic", "u8")
	b524FamilyCodebookEntries, b524FamilyCodebookError := B524FamilyCodebookEntries()
	if b524FamilyCodebookError != nil {
		b524FamilyCodebookT.Fatal(b524FamilyCodebookError)
	}
	b524FamilyCodebookExpected := append([]byte(nil), b524FamilyCodebookEntries...)
	b524FamilyCodebookFirst, b524FamilyCodebookFound, b524FamilyCodebookError := B524FamilyCodebookLookup(b524FamilyCodebookKey)
	if b524FamilyCodebookError != nil || !b524FamilyCodebookFound {
		b524FamilyCodebookT.Fatal("ab30 lookup control")
	}
	b524FamilyCodebookRowExpected := append([]byte(nil), b524FamilyCodebookFirst...)
	b524FamilyCodebookSecond, b524FamilyCodebookFound, b524FamilyCodebookError := B524FamilyCodebookLookup(b524FamilyCodebookKey)
	if b524FamilyCodebookError != nil || !b524FamilyCodebookFound {
		b524FamilyCodebookT.Fatal("ab30 second lookup control")
	}
	for b524FamilyCodebookIndex := range b524FamilyCodebookEntries {
		b524FamilyCodebookEntries[b524FamilyCodebookIndex] = 'x'
	}
	for b524FamilyCodebookIndex := range b524FamilyCodebookFirst {
		b524FamilyCodebookFirst[b524FamilyCodebookIndex] = 'y'
	}
	b524FamilyCodebookLater, b524FamilyCodebookError := B524FamilyCodebookEntries()
	if b524FamilyCodebookError != nil || !b524FamilyCodebookBytes.Equal(b524FamilyCodebookLater, b524FamilyCodebookExpected) {
		b524FamilyCodebookT.Fatal("ab30 entries share storage")
	}
	b524FamilyCodebookRowLater, b524FamilyCodebookFound, b524FamilyCodebookError := B524FamilyCodebookLookup(b524FamilyCodebookKey)
	if b524FamilyCodebookError != nil || !b524FamilyCodebookFound || !b524FamilyCodebookBytes.Equal(b524FamilyCodebookRowLater, b524FamilyCodebookRowExpected) || !b524FamilyCodebookBytes.Equal(b524FamilyCodebookSecond, b524FamilyCodebookRowExpected) {
		b524FamilyCodebookT.Fatal("ab30 row buffers share storage")
	}
	b524FamilyCodebookRows, _, b524FamilyCodebookError := b524FamilyCodebookValidate(b524FamilyCodebookExpected)
	if b524FamilyCodebookError != nil {
		b524FamilyCodebookT.Fatal(b524FamilyCodebookError)
	}
	for _, b524FamilyCodebookRow := range b524FamilyCodebookRows {
		for b524FamilyCodebookIndex := range b524FamilyCodebookRow {
			b524FamilyCodebookRow[b524FamilyCodebookIndex] = 'z'
		}
	}
	if !b524FamilyCodebookBytes.Equal(b524FamilyCodebookLater, b524FamilyCodebookExpected) {
		b524FamilyCodebookT.Fatal("ab30 decoded rows alias input")
	}
}
