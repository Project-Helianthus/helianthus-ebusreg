package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type b555AdditionalFormField struct {
	B555AdditionalFormName     string `json:"name"`
	B555AdditionalFormByte     int    `json:"byte"`
	B555AdditionalFormBit      int    `json:"bit"`
	B555AdditionalFormWidth    int    `json:"width_bits"`
	B555AdditionalFormChannel  int    `json:"channel"`
	B555AdditionalFormWeekday  int    `json:"weekday_index"`
	B555AdditionalFormMeaning  string `json:"meaning"`
	B555AdditionalFormEncoding string `json:"encoding"`
	B555AdditionalFormDomain   []int  `json:"domain"`
	B555AdditionalFormUnit     string `json:"unit"`
	B555AdditionalFormSentinel string `json:"sentinel"`
}

type b555AdditionalFormRow struct {
	B555AdditionalFormID        string                    `json:"id"`
	B555AdditionalFormKey       B555AdditionalFormKey     `json:"key"`
	B555AdditionalFormGrade     string                    `json:"grade"`
	B555AdditionalFormOrigin    string                    `json:"coordinate_origin"`
	B555AdditionalFormPhysical  string                    `json:"physical_offset"`
	B555AdditionalFormFields    []b555AdditionalFormField `json:"fields"`
	B555AdditionalFormRelations struct {
		B555AdditionalFormRequired  int    `json:"required_accessible_initialized_bytes"`
		B555AdditionalFormDeclared  bool   `json:"declared_extent_is_received_extent"`
		B555AdditionalFormLocal     bool   `json:"local_count_is_received_extent"`
		B555AdditionalFormCapacity  bool   `json:"backing_capacity_is_received_extent"`
		B555AdditionalFormStatus    bool   `json:"status_is_extent"`
		B555AdditionalFormThreshold bool   `json:"numeric_threshold_is_extent"`
		B555AdditionalFormExtra     string `json:"extra_bytes"`
		B555AdditionalFormEmpty     bool   `json:"empty_arguments_are_physical_empty_packet"`
	} `json:"relations"`
	B555AdditionalFormRaw           string            `json:"raw_form"`
	B555AdditionalFormUnknown       map[string]string `json:"unknown"`
	B555AdditionalFormQualification string            `json:"qualification"`
	B555AdditionalFormExecutable    bool              `json:"executable"`
}

func b555AdditionalFormRows(b555AdditionalFormT *testing.T) ([]json.RawMessage, []b555AdditionalFormRow) {
	b555AdditionalFormT.Helper()
	b555AdditionalFormData, b555AdditionalFormError := B555AdditionalFormEntries()
	if b555AdditionalFormError != nil {
		b555AdditionalFormT.Fatal(b555AdditionalFormError)
	}
	var b555AdditionalFormRaw []json.RawMessage
	var b555AdditionalFormDecoded []b555AdditionalFormRow
	if b555AdditionalFormError = json.Unmarshal(b555AdditionalFormData, &b555AdditionalFormRaw); b555AdditionalFormError != nil {
		b555AdditionalFormT.Fatal(b555AdditionalFormError)
	}
	if b555AdditionalFormError = json.Unmarshal(b555AdditionalFormData, &b555AdditionalFormDecoded); b555AdditionalFormError != nil {
		b555AdditionalFormT.Fatal(b555AdditionalFormError)
	}
	return b555AdditionalFormRaw, b555AdditionalFormDecoded
}

func b555AdditionalFormReply(b555AdditionalFormT *testing.T, b555AdditionalFormOperation uint8) b555AdditionalFormRow {
	b555AdditionalFormT.Helper()
	_, b555AdditionalFormDecoded := b555AdditionalFormRows(b555AdditionalFormT)
	for _, b555AdditionalFormCandidate := range b555AdditionalFormDecoded {
		if *b555AdditionalFormCandidate.B555AdditionalFormKey.B555AdditionalFormOperation == b555AdditionalFormOperation && b555AdditionalFormCandidate.B555AdditionalFormKey.B555AdditionalFormDirection == "reply" {
			b555AdditionalFormRaw, b555AdditionalFormFound, b555AdditionalFormError := B555AdditionalFormLookup(b555AdditionalFormCandidate.B555AdditionalFormKey)
			if b555AdditionalFormError != nil || !b555AdditionalFormFound {
				b555AdditionalFormT.Fatalf("lookup: %v %v", b555AdditionalFormFound, b555AdditionalFormError)
			}
			var b555AdditionalFormResult b555AdditionalFormRow
			if b555AdditionalFormError = json.Unmarshal(b555AdditionalFormRaw, &b555AdditionalFormResult); b555AdditionalFormError != nil {
				b555AdditionalFormT.Fatal(b555AdditionalFormError)
			}
			return b555AdditionalFormResult
		}
	}
	b555AdditionalFormT.Fatalf("missing operation %x", b555AdditionalFormOperation)
	return b555AdditionalFormRow{}
}

func TestAb26Catalog(b555AdditionalFormT *testing.T) {
	b555AdditionalFormRaw, b555AdditionalFormDecoded := b555AdditionalFormRows(b555AdditionalFormT)
	if len(b555AdditionalFormDecoded) != 14 {
		b555AdditionalFormT.Fatalf("rows: %d", len(b555AdditionalFormDecoded))
	}
	b555AdditionalFormIDs := map[string]bool{}
	for b555AdditionalFormIndex, b555AdditionalFormCandidate := range b555AdditionalFormDecoded {
		if !strings.HasPrefix(b555AdditionalFormCandidate.B555AdditionalFormID, "ab26-") || b555AdditionalFormIDs[b555AdditionalFormCandidate.B555AdditionalFormID] {
			b555AdditionalFormT.Fatal("invalid ID")
		}
		b555AdditionalFormIDs[b555AdditionalFormCandidate.B555AdditionalFormID] = true
		b555AdditionalFormResult, b555AdditionalFormFound, b555AdditionalFormError := B555AdditionalFormLookup(b555AdditionalFormCandidate.B555AdditionalFormKey)
		if b555AdditionalFormError != nil || !b555AdditionalFormFound || !bytes.Equal(b555AdditionalFormResult, b555AdditionalFormRaw[b555AdditionalFormIndex]) {
			b555AdditionalFormT.Fatal("whole-row lookup mismatch")
		}
		if b555AdditionalFormCandidate.B555AdditionalFormGrade != "Hypothesis" || b555AdditionalFormCandidate.B555AdditionalFormQualification != "Unknown" || b555AdditionalFormCandidate.B555AdditionalFormPhysical != "Unknown" || b555AdditionalFormCandidate.B555AdditionalFormOrigin != "immediately_after_command" || b555AdditionalFormCandidate.B555AdditionalFormExecutable {
			b555AdditionalFormT.Fatal("grade or qualification")
		}
		if b555AdditionalFormCandidate.B555AdditionalFormRaw != "retain_all_original_bytes_including_reserved_and_extra_bytes" {
			b555AdditionalFormT.Fatal("raw fallback")
		}
		for _, b555AdditionalFormName := range []string{"native_support", "physical_adaptation", "physical_effects", "persistence", "delivery", "other_timer_variants", "unit", "sentinel"} {
			if b555AdditionalFormCandidate.B555AdditionalFormUnknown[b555AdditionalFormName] != "Unknown" {
				b555AdditionalFormT.Fatal(b555AdditionalFormName)
			}
		}
		b555AdditionalFormRelations := b555AdditionalFormCandidate.B555AdditionalFormRelations
		if b555AdditionalFormRelations.B555AdditionalFormDeclared || b555AdditionalFormRelations.B555AdditionalFormLocal || b555AdditionalFormRelations.B555AdditionalFormCapacity || b555AdditionalFormRelations.B555AdditionalFormStatus || b555AdditionalFormRelations.B555AdditionalFormThreshold || b555AdditionalFormRelations.B555AdditionalFormEmpty || b555AdditionalFormRelations.B555AdditionalFormExtra != "retain_raw" {
			b555AdditionalFormT.Fatal("independent extents/raw retention")
		}
		b555AdditionalFormKey := b555AdditionalFormCandidate.B555AdditionalFormKey
		if b555AdditionalFormKey.B555AdditionalFormFamily != "vaillant" || b555AdditionalFormKey.B555AdditionalFormVariant != "b555" || b555AdditionalFormKey.B555AdditionalFormPrimary != 0xb5 || b555AdditionalFormKey.B555AdditionalFormSecondary == nil || *b555AdditionalFormKey.B555AdditionalFormSecondary != 0x55 || b555AdditionalFormKey.B555AdditionalFormGroup != nil || b555AdditionalFormKey.B555AdditionalFormRegister != nil || b555AdditionalFormKey.B555AdditionalFormInstanceScope != "none" {
			b555AdditionalFormT.Fatal("identity")
		}
		if b555AdditionalFormKey.B555AdditionalFormDirection == "request" && (len(b555AdditionalFormCandidate.B555AdditionalFormFields) != 0 || b555AdditionalFormRelations.B555AdditionalFormRequired != 0) {
			b555AdditionalFormT.Fatal("declared empty arguments")
		}
	}
}

func TestAb26ExactAxes(b555AdditionalFormT *testing.T) {
	b555AdditionalFormBase := b555AdditionalFormReply(b555AdditionalFormT, 0xa0).B555AdditionalFormKey
	b555AdditionalFormZero := uint8(0)
	b555AdditionalFormZeroRegister := uint16(0)
	b555AdditionalFormChanges := []struct {
		B555AdditionalFormName   string
		B555AdditionalFormChange func(*B555AdditionalFormKey)
	}{
		{"ab26Family", func(b555AdditionalFormKey *B555AdditionalFormKey) {
			b555AdditionalFormKey.B555AdditionalFormFamily = "other"
		}},
		{"ab26Variant", func(b555AdditionalFormKey *B555AdditionalFormKey) {
			b555AdditionalFormKey.B555AdditionalFormVariant = "other"
		}},
		{"ab26Primary", func(b555AdditionalFormKey *B555AdditionalFormKey) {
			b555AdditionalFormKey.B555AdditionalFormPrimary = 0
		}},
		{"ab26SecondaryAbsent", func(b555AdditionalFormKey *B555AdditionalFormKey) {
			b555AdditionalFormKey.B555AdditionalFormSecondary = nil
		}},
		{"ab26SecondaryZero", func(b555AdditionalFormKey *B555AdditionalFormKey) {
			b555AdditionalFormKey.B555AdditionalFormSecondary = &b555AdditionalFormZero
		}},
		{"ab26OperationAbsent", func(b555AdditionalFormKey *B555AdditionalFormKey) {
			b555AdditionalFormKey.B555AdditionalFormOperation = nil
		}},
		{"ab26OperationZero", func(b555AdditionalFormKey *B555AdditionalFormKey) {
			b555AdditionalFormKey.B555AdditionalFormOperation = &b555AdditionalFormZero
		}},
		{"ab26Direction", func(b555AdditionalFormKey *B555AdditionalFormKey) {
			b555AdditionalFormKey.B555AdditionalFormDirection = "none"
		}},
		{"ab26Field", func(b555AdditionalFormKey *B555AdditionalFormKey) {
			b555AdditionalFormKey.B555AdditionalFormField = "other"
		}},
		{"ab26Alternative", func(b555AdditionalFormKey *B555AdditionalFormKey) {
			b555AdditionalFormKey.B555AdditionalFormAlternative = "other"
		}},
		{"ab26Group", func(b555AdditionalFormKey *B555AdditionalFormKey) {
			b555AdditionalFormKey.B555AdditionalFormGroup = &b555AdditionalFormZero
		}},
		{"ab26InstanceScope", func(b555AdditionalFormKey *B555AdditionalFormKey) {
			b555AdditionalFormKey.B555AdditionalFormInstanceScope = "dynamic"
		}},
		{"ab26Register", func(b555AdditionalFormKey *B555AdditionalFormKey) {
			b555AdditionalFormKey.B555AdditionalFormRegister = &b555AdditionalFormZeroRegister
		}},
	}
	for _, b555AdditionalFormCase := range b555AdditionalFormChanges {
		b555AdditionalFormT.Run(b555AdditionalFormCase.B555AdditionalFormName, func(b555AdditionalFormT *testing.T) {
			b555AdditionalFormKey := b555AdditionalFormBase
			b555AdditionalFormCase.B555AdditionalFormChange(&b555AdditionalFormKey)
			b555AdditionalFormResult, b555AdditionalFormFound, b555AdditionalFormError := B555AdditionalFormLookup(b555AdditionalFormKey)
			if b555AdditionalFormError != nil || b555AdditionalFormFound || b555AdditionalFormResult != nil {
				b555AdditionalFormT.Fatalf("inexact key admitted: %v %v", b555AdditionalFormFound, b555AdditionalFormError)
			}
			b555AdditionalFormData, b555AdditionalFormError := json.Marshal([]struct {
				B555AdditionalFormKey B555AdditionalFormKey `json:"key"`
			}{{b555AdditionalFormBase}, {b555AdditionalFormKey}})
			if b555AdditionalFormError != nil {
				b555AdditionalFormT.Fatal(b555AdditionalFormError)
			}
			if _, _, b555AdditionalFormError = b555AdditionalFormValidate(b555AdditionalFormData); b555AdditionalFormError != nil {
				b555AdditionalFormT.Fatalf("distinct keys rejected: %v", b555AdditionalFormError)
			}
		})
	}
	// Optional absent and explicit zero are distinct even where production rows use a nonzero value.
	for _, b555AdditionalFormAxis := range []string{"secondary", "operation", "group", "register"} {
		b555AdditionalFormT.Run("ab26Optional"+b555AdditionalFormAxis, func(b555AdditionalFormT *testing.T) {
			b555AdditionalFormFirst := b555AdditionalFormBase
			b555AdditionalFormSecond := b555AdditionalFormBase
			switch b555AdditionalFormAxis {
			case "secondary":
				b555AdditionalFormFirst.B555AdditionalFormSecondary = nil
				b555AdditionalFormSecond.B555AdditionalFormSecondary = &b555AdditionalFormZero
			case "operation":
				b555AdditionalFormFirst.B555AdditionalFormOperation = nil
				b555AdditionalFormSecond.B555AdditionalFormOperation = &b555AdditionalFormZero
			case "group":
				b555AdditionalFormFirst.B555AdditionalFormGroup = nil
				b555AdditionalFormSecond.B555AdditionalFormGroup = &b555AdditionalFormZero
			case "register":
				b555AdditionalFormFirst.B555AdditionalFormRegister = nil
				b555AdditionalFormSecond.B555AdditionalFormRegister = &b555AdditionalFormZeroRegister
			}
			b555AdditionalFormData, b555AdditionalFormError := json.Marshal([]struct {
				B555AdditionalFormKey B555AdditionalFormKey `json:"key"`
			}{{b555AdditionalFormFirst}, {b555AdditionalFormSecond}})
			if b555AdditionalFormError != nil {
				b555AdditionalFormT.Fatal(b555AdditionalFormError)
			}
			if _, _, b555AdditionalFormError = b555AdditionalFormValidate(b555AdditionalFormData); b555AdditionalFormError != nil {
				b555AdditionalFormT.Fatal(b555AdditionalFormError)
			}
		})
	}
}

func TestAb26ValidationAndCopies(b555AdditionalFormT *testing.T) {
	b555AdditionalFormRaw, b555AdditionalFormDecoded := b555AdditionalFormRows(b555AdditionalFormT)
	for _, b555AdditionalFormInvalid := range []string{"{", "null", "{}", "[null]", "[{}]", "[1]", "[{\"key\":{\"primary\":256}}]", "[" + string(b555AdditionalFormRaw[0]) + "," + string(b555AdditionalFormRaw[0]) + "]"} {
		if _, _, b555AdditionalFormError := b555AdditionalFormValidate([]byte(b555AdditionalFormInvalid)); b555AdditionalFormError == nil {
			b555AdditionalFormT.Fatal("invalid data admitted")
		}
	}
	var b555AdditionalFormObject map[string]json.RawMessage
	if b555AdditionalFormError := json.Unmarshal(b555AdditionalFormRaw[0], &b555AdditionalFormObject); b555AdditionalFormError != nil {
		b555AdditionalFormT.Fatal(b555AdditionalFormError)
	}
	var b555AdditionalFormIdentity map[string]json.RawMessage
	if b555AdditionalFormError := json.Unmarshal(b555AdditionalFormObject["key"], &b555AdditionalFormIdentity); b555AdditionalFormError != nil {
		b555AdditionalFormT.Fatal(b555AdditionalFormError)
	}
	for b555AdditionalFormAxis, b555AdditionalFormValue := range b555AdditionalFormIdentity {
		delete(b555AdditionalFormIdentity, b555AdditionalFormAxis)
		b555AdditionalFormKeyData, _ := json.Marshal(b555AdditionalFormIdentity)
		b555AdditionalFormObject["key"] = b555AdditionalFormKeyData
		b555AdditionalFormData, _ := json.Marshal([]map[string]json.RawMessage{b555AdditionalFormObject})
		if _, _, b555AdditionalFormError := b555AdditionalFormValidate(b555AdditionalFormData); b555AdditionalFormError == nil {
			b555AdditionalFormT.Fatalf("missing %s admitted", b555AdditionalFormAxis)
		}
		b555AdditionalFormIdentity[b555AdditionalFormAxis] = b555AdditionalFormValue
	}
	b555AdditionalFormEntries, b555AdditionalFormError := B555AdditionalFormEntries()
	if b555AdditionalFormError != nil {
		b555AdditionalFormT.Fatal(b555AdditionalFormError)
	}
	b555AdditionalFormExpected := append([]byte(nil), b555AdditionalFormEntries...)
	b555AdditionalFormOther, _ := B555AdditionalFormEntries()
	for b555AdditionalFormIndex := range b555AdditionalFormEntries {
		b555AdditionalFormEntries[b555AdditionalFormIndex] = 'x'
	}
	b555AdditionalFormLater, b555AdditionalFormError := B555AdditionalFormEntries()
	if b555AdditionalFormError != nil || !bytes.Equal(b555AdditionalFormLater, b555AdditionalFormExpected) || !bytes.Equal(b555AdditionalFormOther, b555AdditionalFormExpected) {
		b555AdditionalFormT.Fatal("entries share memory")
	}
	for b555AdditionalFormIndex, b555AdditionalFormCandidate := range b555AdditionalFormDecoded {
		b555AdditionalFormResult, _, _ := B555AdditionalFormLookup(b555AdditionalFormCandidate.B555AdditionalFormKey)
		b555AdditionalFormOther, _, _ := B555AdditionalFormLookup(b555AdditionalFormCandidate.B555AdditionalFormKey)
		for b555AdditionalFormByte := range b555AdditionalFormResult {
			b555AdditionalFormResult[b555AdditionalFormByte] = 'x'
		}
		b555AdditionalFormLater, b555AdditionalFormFound, b555AdditionalFormError := B555AdditionalFormLookup(b555AdditionalFormCandidate.B555AdditionalFormKey)
		if b555AdditionalFormError != nil || !b555AdditionalFormFound || !bytes.Equal(b555AdditionalFormLater, b555AdditionalFormRaw[b555AdditionalFormIndex]) || !bytes.Equal(b555AdditionalFormOther, b555AdditionalFormLater) {
			b555AdditionalFormT.Fatal("row shares memory")
		}
	}
}

func b555AdditionalFormValue(b555AdditionalFormField b555AdditionalFormField, b555AdditionalFormRaw []byte) uint64 {
	var b555AdditionalFormResult uint64
	for b555AdditionalFormBit := 0; b555AdditionalFormBit < b555AdditionalFormField.B555AdditionalFormWidth; b555AdditionalFormBit++ {
		b555AdditionalFormPosition := b555AdditionalFormField.B555AdditionalFormByte*8 + b555AdditionalFormField.B555AdditionalFormBit + b555AdditionalFormBit
		b555AdditionalFormResult |= uint64((b555AdditionalFormRaw[b555AdditionalFormPosition/8]>>uint(b555AdditionalFormPosition%8))&1) << uint(b555AdditionalFormBit)
	}
	return b555AdditionalFormResult
}

func TestAb26Masks(b555AdditionalFormT *testing.T) {
	for _, b555AdditionalFormOperation := range []uint8{0xa0, 0xa7} {
		b555AdditionalFormRow := b555AdditionalFormReply(b555AdditionalFormT, b555AdditionalFormOperation)
		b555AdditionalFormLabel := "availability"
		if b555AdditionalFormOperation == 0xa7 {
			b555AdditionalFormLabel = "visibility"
		}
		if b555AdditionalFormRow.B555AdditionalFormKey.B555AdditionalFormField != b555AdditionalFormLabel || b555AdditionalFormRow.B555AdditionalFormRelations.B555AdditionalFormRequired != 8 || len(b555AdditionalFormRow.B555AdditionalFormFields) != 55 {
			b555AdditionalFormT.Fatal("mask meaning or extent")
		}
		b555AdditionalFormOccupied := [64]bool{}
		for b555AdditionalFormIndex, b555AdditionalFormField := range b555AdditionalFormRow.B555AdditionalFormFields {
			b555AdditionalFormExpectedByte, b555AdditionalFormExpectedBit, b555AdditionalFormExpectedWidth := 0, b555AdditionalFormIndex, 1
			if b555AdditionalFormIndex < 5 {
				if b555AdditionalFormField.B555AdditionalFormName != []string{"DHW", "circulation", "noise_reduction", "ventilation", "high_tariff"}[b555AdditionalFormIndex] {
					b555AdditionalFormT.Fatal("global bit label")
				}
			} else if b555AdditionalFormIndex == 5 {
				b555AdditionalFormExpectedBit = 5
				b555AdditionalFormExpectedWidth = 3
			} else if b555AdditionalFormIndex == 6 {
				b555AdditionalFormExpectedByte = 1
				b555AdditionalFormExpectedBit = 0
				b555AdditionalFormExpectedWidth = 8
			} else {
				b555AdditionalFormChannel := (b555AdditionalFormIndex-7)/4 + 1
				b555AdditionalFormRole := (b555AdditionalFormIndex - 7) % 4
				// Independent explicit byte/nibble controls, including channel 12 at byte 7.
				b555AdditionalFormExpectedByte = []int{2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7}[b555AdditionalFormChannel-1]
				b555AdditionalFormExpectedBit = []int{0, 4, 0, 4, 0, 4, 0, 4, 0, 4, 0, 4}[b555AdditionalFormChannel-1] + b555AdditionalFormRole
				if b555AdditionalFormField.B555AdditionalFormChannel != b555AdditionalFormChannel || b555AdditionalFormField.B555AdditionalFormName != []string{"zone_heating", "zone_cooling", "external_DHW", "external_fixed_flow"}[b555AdditionalFormRole] {
					b555AdditionalFormT.Fatal("channel role")
				}
			}
			if b555AdditionalFormField.B555AdditionalFormByte != b555AdditionalFormExpectedByte || b555AdditionalFormField.B555AdditionalFormBit != b555AdditionalFormExpectedBit || b555AdditionalFormField.B555AdditionalFormWidth != b555AdditionalFormExpectedWidth {
				b555AdditionalFormT.Fatalf("coordinate %d: %+v", b555AdditionalFormIndex, b555AdditionalFormField)
			}
			if b555AdditionalFormIndex == 5 || b555AdditionalFormIndex == 6 {
				if b555AdditionalFormField.B555AdditionalFormName != "reserved" || b555AdditionalFormField.B555AdditionalFormEncoding != "raw" {
					b555AdditionalFormT.Fatal("reserved retention")
				}
			} else if b555AdditionalFormField.B555AdditionalFormMeaning != b555AdditionalFormLabel {
				b555AdditionalFormT.Fatal("availability/visibility alias")
			}
			for b555AdditionalFormBit := 0; b555AdditionalFormBit < b555AdditionalFormExpectedWidth; b555AdditionalFormBit++ {
				b555AdditionalFormPosition := b555AdditionalFormExpectedByte*8 + b555AdditionalFormExpectedBit + b555AdditionalFormBit
				if b555AdditionalFormOccupied[b555AdditionalFormPosition] {
					b555AdditionalFormT.Fatal("overlap")
				}
				b555AdditionalFormOccupied[b555AdditionalFormPosition] = true
				b555AdditionalFormRaw := make([]byte, 8)
				b555AdditionalFormRaw[b555AdditionalFormPosition/8] = 1 << uint(b555AdditionalFormPosition%8)
				for b555AdditionalFormOtherIndex, b555AdditionalFormOtherField := range b555AdditionalFormRow.B555AdditionalFormFields {
					b555AdditionalFormExpected := uint64(0)
					if b555AdditionalFormOtherIndex == b555AdditionalFormIndex {
						b555AdditionalFormExpected = 1 << uint(b555AdditionalFormBit)
					}
					if b555AdditionalFormValue(b555AdditionalFormOtherField, b555AdditionalFormRaw) != b555AdditionalFormExpected {
						b555AdditionalFormT.Fatal("one-bit isolation or retention")
					}
				}
			}
		}
		for _, b555AdditionalFormPresent := range b555AdditionalFormOccupied {
			if !b555AdditionalFormPresent {
				b555AdditionalFormT.Fatal("lost bit")
			}
		}
	}
}

func TestAb26ReservedReplies(b555AdditionalFormT *testing.T) {
	for _, b555AdditionalFormOperation := range []uint8{0xa1, 0xa2, 0xa8, 0xa9} {
		b555AdditionalFormRow := b555AdditionalFormReply(b555AdditionalFormT, b555AdditionalFormOperation)
		if b555AdditionalFormRow.B555AdditionalFormKey.B555AdditionalFormField != "reserved" || b555AdditionalFormRow.B555AdditionalFormExecutable || b555AdditionalFormRow.B555AdditionalFormRelations.B555AdditionalFormRequired != 6 || len(b555AdditionalFormRow.B555AdditionalFormFields) != 1 {
			b555AdditionalFormT.Fatal("reserved settings promoted")
		}
		b555AdditionalFormField := b555AdditionalFormRow.B555AdditionalFormFields[0]
		if b555AdditionalFormField.B555AdditionalFormName != "reserved" || b555AdditionalFormField.B555AdditionalFormByte != 0 || b555AdditionalFormField.B555AdditionalFormBit != 0 || b555AdditionalFormField.B555AdditionalFormWidth != 48 || b555AdditionalFormField.B555AdditionalFormEncoding != "raw" || b555AdditionalFormValue(b555AdditionalFormField, []byte{1, 2, 3, 4, 5, 6}) != 0x060504030201 {
			b555AdditionalFormT.Fatal("reserved bytes")
		}
	}
}

func TestAb26PositiveForms(b555AdditionalFormT *testing.T) {
	for _, b555AdditionalFormCase := range []struct {
		B555AdditionalFormOperation uint8
		B555AdditionalFormRequired  int
	}{{0xa4, 9}, {0xa5, 7}} {
		b555AdditionalFormRow := b555AdditionalFormReply(b555AdditionalFormT, b555AdditionalFormCase.B555AdditionalFormOperation)
		if b555AdditionalFormRow.B555AdditionalFormRelations.B555AdditionalFormRequired != b555AdditionalFormCase.B555AdditionalFormRequired || b555AdditionalFormRow.B555AdditionalFormKey.B555AdditionalFormAlternative != "positive" {
			b555AdditionalFormT.Fatal("positive extent")
		}
		for _, b555AdditionalFormAccessible := range []int{b555AdditionalFormCase.B555AdditionalFormRequired - 1, b555AdditionalFormCase.B555AdditionalFormRequired, b555AdditionalFormCase.B555AdditionalFormRequired + 1} {
			for _, b555AdditionalFormInitialized := range []int{0, b555AdditionalFormCase.B555AdditionalFormRequired - 1, 16} {
				for _, b555AdditionalFormDeclared := range []int{0, 16} {
					for _, b555AdditionalFormLocal := range []int{0, 16} {
						b555AdditionalFormT.Run(fmt.Sprintf("ab26%x-accessible%d-initialized%d-declared%d-local%d", b555AdditionalFormCase.B555AdditionalFormOperation, b555AdditionalFormAccessible, b555AdditionalFormInitialized, b555AdditionalFormDeclared, b555AdditionalFormLocal), func(b555AdditionalFormT *testing.T) {
							b555AdditionalFormStorage := make([]byte, b555AdditionalFormAccessible, 16)
							b555AdditionalFormBefore := append([]byte(nil), b555AdditionalFormStorage...)
							b555AdditionalFormComplete := len(b555AdditionalFormStorage) >= b555AdditionalFormRow.B555AdditionalFormRelations.B555AdditionalFormRequired && b555AdditionalFormInitialized >= b555AdditionalFormRow.B555AdditionalFormRelations.B555AdditionalFormRequired
							b555AdditionalFormExpected := b555AdditionalFormAccessible >= b555AdditionalFormCase.B555AdditionalFormRequired && b555AdditionalFormInitialized >= b555AdditionalFormCase.B555AdditionalFormRequired
							if b555AdditionalFormComplete != b555AdditionalFormExpected || !bytes.Equal(b555AdditionalFormStorage, b555AdditionalFormBefore) {
								b555AdditionalFormT.Fatal("extent independence/raw retention")
							}
						})
					}
				}
			}
		}
		if b555AdditionalFormCase.B555AdditionalFormOperation == 0xa4 {
			if len(b555AdditionalFormRow.B555AdditionalFormFields) != 9 {
				b555AdditionalFormT.Fatal("weekday fields")
			}
			for b555AdditionalFormIndex, b555AdditionalFormField := range b555AdditionalFormRow.B555AdditionalFormFields {
				b555AdditionalFormName := "weekday_count"
				if b555AdditionalFormIndex == 0 {
					b555AdditionalFormName = "status"
				}
				if b555AdditionalFormIndex == 8 {
					b555AdditionalFormName = "24_hour_count"
				}
				if b555AdditionalFormField.B555AdditionalFormByte != b555AdditionalFormIndex || b555AdditionalFormField.B555AdditionalFormWidth != 8 || b555AdditionalFormField.B555AdditionalFormName != b555AdditionalFormName {
					b555AdditionalFormT.Fatal("weekday coordinate")
				}
				if b555AdditionalFormIndex > 0 && b555AdditionalFormIndex < 8 && b555AdditionalFormField.B555AdditionalFormWeekday != b555AdditionalFormIndex {
					b555AdditionalFormT.Fatal("weekday index")
				}
			}
		} else {
			if len(b555AdditionalFormRow.B555AdditionalFormFields) != 6 {
				b555AdditionalFormT.Fatal("interval fields")
			}
			for b555AdditionalFormIndex, b555AdditionalFormField := range b555AdditionalFormRow.B555AdditionalFormFields {
				if b555AdditionalFormField.B555AdditionalFormByte != b555AdditionalFormIndex || b555AdditionalFormField.B555AdditionalFormName != []string{"status", "start_hour", "start_minute", "end_hour", "end_minute", "setpoint"}[b555AdditionalFormIndex] {
					b555AdditionalFormT.Fatal("interval coordinate")
				}
				if b555AdditionalFormIndex < 5 && b555AdditionalFormField.B555AdditionalFormWidth != 8 {
					b555AdditionalFormT.Fatal("interval width")
				}
				if b555AdditionalFormIndex > 0 && b555AdditionalFormIndex < 5 {
					b555AdditionalFormMax := 59
					if b555AdditionalFormIndex%2 == 1 {
						b555AdditionalFormMax = 24
					}
					if !reflect.DeepEqual(b555AdditionalFormField.B555AdditionalFormDomain, []int{0, b555AdditionalFormMax}) || b555AdditionalFormField.B555AdditionalFormEncoding != "unsigned" {
						b555AdditionalFormT.Fatal("time domain")
					}
				}
			}
			b555AdditionalFormField := b555AdditionalFormRow.B555AdditionalFormFields[5]
			if b555AdditionalFormField.B555AdditionalFormWidth != 16 || b555AdditionalFormField.B555AdditionalFormEncoding != "unsigned_little_endian" || b555AdditionalFormField.B555AdditionalFormUnit != "Unknown" || b555AdditionalFormField.B555AdditionalFormSentinel != "Unknown" {
				b555AdditionalFormT.Fatal("setpoint representation")
			}
			for _, b555AdditionalFormCase := range []struct {
				B555AdditionalFormRaw      []byte
				B555AdditionalFormExpected uint64
			}{{[]byte{0, 24, 59, 0, 0, 255, 255, 42}, 65535}, {[]byte{255, 0, 0, 24, 59, 0x34, 0x12, 42}, 4660}} {
				b555AdditionalFormBefore := append([]byte(nil), b555AdditionalFormCase.B555AdditionalFormRaw...)
				if b555AdditionalFormValue(b555AdditionalFormField, b555AdditionalFormCase.B555AdditionalFormRaw) != b555AdditionalFormCase.B555AdditionalFormExpected || !bytes.Equal(b555AdditionalFormBefore, b555AdditionalFormCase.B555AdditionalFormRaw) {
					b555AdditionalFormT.Fatal("raw numeric or extra byte changed")
				}
			}
		}
	}
}
