package vaillant

import (
	b501LegacySequenceBytes "bytes"
	b501LegacySequenceJSON "encoding/json"
	b501LegacySequenceReflect "reflect"
	b501LegacySequenceStrings "strings"
	b501LegacySequenceTesting "testing"
)

func b501LegacySequenceBaseKey() B501LegacySequenceKey {
	b501LegacySequenceSecondary := uint8(1)
	return B501LegacySequenceKey{B501LegacySequenceFamily: "legacy-interface", B501LegacySequenceVariant: "legacy-computer-interface", B501LegacySequencePrimary: 181, B501LegacySequenceSecondary: &b501LegacySequenceSecondary, B501LegacySequenceDirection: "initiator-to-initiator", B501LegacySequenceField: "sequencing", B501LegacySequenceAlternative: "none", B501LegacySequenceInstanceScope: "none"}
}

func b501LegacySequenceDecode(b501LegacySequenceT *b501LegacySequenceTesting.T, b501LegacySequenceData []byte) map[string]any {
	b501LegacySequenceT.Helper()
	var b501LegacySequenceObject map[string]any
	if b501LegacySequenceError := b501LegacySequenceJSON.Unmarshal(b501LegacySequenceData, &b501LegacySequenceObject); b501LegacySequenceError != nil {
		b501LegacySequenceT.Fatal(b501LegacySequenceError)
	}
	return b501LegacySequenceObject
}

func b501LegacySequenceMarshal(b501LegacySequenceT *b501LegacySequenceTesting.T, b501LegacySequenceValue any) []byte {
	b501LegacySequenceT.Helper()
	b501LegacySequenceData, b501LegacySequenceError := b501LegacySequenceJSON.Marshal(b501LegacySequenceValue)
	if b501LegacySequenceError != nil {
		b501LegacySequenceT.Fatal(b501LegacySequenceError)
	}
	return b501LegacySequenceData
}

func b501LegacySequenceGet(b501LegacySequenceT *b501LegacySequenceTesting.T, b501LegacySequenceKey B501LegacySequenceKey) []byte {
	b501LegacySequenceT.Helper()
	b501LegacySequenceRow, b501LegacySequenceFound, b501LegacySequenceError := B501LegacySequenceLookup(b501LegacySequenceKey)
	if b501LegacySequenceError != nil || !b501LegacySequenceFound || len(b501LegacySequenceRow) == 0 {
		b501LegacySequenceT.Fatalf("ab34 lookup: found=%v error=%v", b501LegacySequenceFound, b501LegacySequenceError)
	}
	return b501LegacySequenceRow
}

func TestAb34Contract(b501LegacySequenceT *b501LegacySequenceTesting.T) {
	b501LegacySequenceData, b501LegacySequenceError := B501LegacySequenceEntries()
	if b501LegacySequenceError != nil {
		b501LegacySequenceT.Fatal(b501LegacySequenceError)
	}
	var b501LegacySequenceRows []b501LegacySequenceJSON.RawMessage
	if b501LegacySequenceError = b501LegacySequenceJSON.Unmarshal(b501LegacySequenceData, &b501LegacySequenceRows); b501LegacySequenceError != nil || len(b501LegacySequenceRows) != 3 {
		b501LegacySequenceT.Fatalf("ab34 finite array: %d %v", len(b501LegacySequenceRows), b501LegacySequenceError)
	}
	for b501LegacySequenceIndex, b501LegacySequenceRole := range []string{"request", "positive-status", "return-data"} {
		b501LegacySequenceKey := b501LegacySequenceBaseKey()
		b501LegacySequenceService := "B501"
		if b501LegacySequenceIndex == 1 {
			b501LegacySequenceKey.B501LegacySequenceField, b501LegacySequenceKey.B501LegacySequenceAlternative = "status", "positive"
		}
		if b501LegacySequenceIndex == 2 {
			b501LegacySequenceSecondary := uint8(2)
			b501LegacySequenceKey.B501LegacySequenceSecondary = &b501LegacySequenceSecondary
			b501LegacySequenceKey.B501LegacySequenceField = "data"
			b501LegacySequenceService = "B502"
		}
		b501LegacySequenceRaw := b501LegacySequenceGet(b501LegacySequenceT, b501LegacySequenceKey)
		if !b501LegacySequenceBytes.Equal(b501LegacySequenceRaw, b501LegacySequenceRows[b501LegacySequenceIndex]) {
			b501LegacySequenceT.Fatal("ab34 lookup must return the whole original row")
		}
		b501LegacySequenceRow := b501LegacySequenceDecode(b501LegacySequenceT, b501LegacySequenceRaw)
		b501LegacySequenceWant := map[string]any{
			"ab34_id":                   "ab34_legacy_" + b501LegacySequenceRole,
			"ab34_identity":             b501LegacySequenceDecode(b501LegacySequenceT, b501LegacySequenceMarshal(b501LegacySequenceT, b501LegacySequenceKey)),
			"ab34_representation_grade": "Proven", "ab34_candidate_grade": "Hypothesis",
			"ab34_service": b501LegacySequenceService, "ab34_service_role": b501LegacySequenceRole,
			"ab34_raw_form": map[string]any{"ab34_service": b501LegacySequenceService, "ab34_payload": nil, "ab34_policy": "preserve complete native bytes and unaffected last-known-good fields"},
			"ab34_unknown":  map[string]any{"ab34_synchronization_bytes": "Unknown", "ab34_state_order_beyond_wait": "Unknown", "ab34_timeout_duration": "Unknown", "ab34_segmentation": "Unknown", "ab34_retries": "Unknown", "ab34_return_codec": "Unknown", "ab34_transport_adaptation": "Unknown", "ab34_native_support": "Unknown", "ab34_physical_effects": "Unknown", "ab34_persistence": "Unknown", "ab34_delivery": "Unknown"},
		}
		if b501LegacySequenceIndex == 0 {
			b501LegacySequenceWant["ab34_sequence"] = b501LegacySequenceDecode(b501LegacySequenceT, []byte(`{"ab34_synchronization_required":true,"ab34_synchronization_token":null,"ab34_timeout_cooldown_min_ms":3500,"ab34_cooldown_event":"timeout","ab34_success_cooldown_required":false,"ab34_request_service":"B501","ab34_positive_status_service":"B501","ab34_return_data_service":"B502","ab34_status_is_return_data":false}`))
		} else {
			b501LegacySequenceWant["ab34_context_only"] = true
		}
		if !b501LegacySequenceReflect.DeepEqual(b501LegacySequenceRow, b501LegacySequenceWant) {
			b501LegacySequenceT.Errorf("ab34 %s differs: got %s want %s", b501LegacySequenceRole, b501LegacySequenceRaw, b501LegacySequenceMarshal(b501LegacySequenceT, b501LegacySequenceWant))
		}
	}
}

func TestAb34TimeoutBoundary(b501LegacySequenceT *b501LegacySequenceTesting.T) {
	b501LegacySequenceSequence := b501LegacySequenceDecode(b501LegacySequenceT, b501LegacySequenceGet(b501LegacySequenceT, b501LegacySequenceBaseKey()))["ab34_sequence"].(map[string]any)
	for _, b501LegacySequenceCase := range []struct {
		b501LegacySequenceName    string
		b501LegacySequenceEvent   string
		b501LegacySequenceElapsed int
		b501LegacySequenceWant    bool
	}{
		{"ab34_timeout_3499", "timeout", 3499, false},
		{"ab34_timeout_3500", "timeout", 3500, true},
		{"ab34_timeout_3501", "timeout", 3501, true},
		{"ab34_success_0", "successful-return", 0, true},
		{"ab34_success_3499", "successful-return", 3499, true},
	} {
		b501LegacySequenceT.Run(b501LegacySequenceCase.b501LegacySequenceName, func(b501LegacySequenceT *b501LegacySequenceTesting.T) {
			b501LegacySequenceAllowed := b501LegacySequenceCase.b501LegacySequenceEvent != b501LegacySequenceSequence["ab34_cooldown_event"] || b501LegacySequenceCase.b501LegacySequenceElapsed >= int(b501LegacySequenceSequence["ab34_timeout_cooldown_min_ms"].(float64))
			if b501LegacySequenceAllowed != b501LegacySequenceCase.b501LegacySequenceWant {
				b501LegacySequenceT.Fatal("ab34 timeout-only cooldown boundary differs")
			}
		})
	}
	for _, b501LegacySequenceName := range []string{"ab34_timeout_duration_ms", "ab34_timeout_duration", "ab34_synchronization_bytes"} {
		if _, b501LegacySequencePresent := b501LegacySequenceSequence[b501LegacySequenceName]; b501LegacySequencePresent {
			b501LegacySequenceT.Fatalf("ab34 invented %s", b501LegacySequenceName)
		}
	}
}

func TestAb34EveryIdentityAxis(b501LegacySequenceT *b501LegacySequenceTesting.T) {
	b501LegacySequenceBase := b501LegacySequenceBaseKey()
	b501LegacySequenceOriginal := b501LegacySequenceGet(b501LegacySequenceT, b501LegacySequenceBase)
	for b501LegacySequenceIndex := 0; b501LegacySequenceIndex < b501LegacySequenceReflect.TypeOf(b501LegacySequenceBase).NumField(); b501LegacySequenceIndex++ {
		b501LegacySequenceName := b501LegacySequenceReflect.TypeOf(b501LegacySequenceBase).Field(b501LegacySequenceIndex).Name
		b501LegacySequenceT.Run("ab34_"+b501LegacySequenceName, func(b501LegacySequenceT *b501LegacySequenceTesting.T) {
			b501LegacySequenceChanged := b501LegacySequenceBaseKey()
			b501LegacySequenceField := b501LegacySequenceReflect.ValueOf(&b501LegacySequenceChanged).Elem().Field(b501LegacySequenceIndex)
			switch b501LegacySequenceField.Kind() {
			case b501LegacySequenceReflect.String:
				b501LegacySequenceField.SetString("ab34_synthetic_other")
				if b501LegacySequenceName == "B501LegacySequenceFamily" {
					b501LegacySequenceField.SetString("modern-controller")
				}
			case b501LegacySequenceReflect.Uint8:
				b501LegacySequenceField.SetUint(180)
			case b501LegacySequenceReflect.Pointer:
				b501LegacySequencePointer := b501LegacySequenceReflect.New(b501LegacySequenceField.Type().Elem())
				b501LegacySequencePointer.Elem().SetUint(0)
				b501LegacySequenceField.Set(b501LegacySequencePointer)
			}
			b501LegacySequenceRow, b501LegacySequenceFound, b501LegacySequenceError := B501LegacySequenceLookup(b501LegacySequenceChanged)
			if b501LegacySequenceError != nil || b501LegacySequenceFound || b501LegacySequenceRow != nil {
				b501LegacySequenceT.Fatalf("ab34 changed axis matched: %v %v", b501LegacySequenceFound, b501LegacySequenceError)
			}
			b501LegacySequenceControl := b501LegacySequenceDecode(b501LegacySequenceT, b501LegacySequenceOriginal)
			b501LegacySequenceControl["ab34_identity"] = b501LegacySequenceChanged
			b501LegacySequenceData := b501LegacySequenceMarshal(b501LegacySequenceT, []any{b501LegacySequenceDecode(b501LegacySequenceT, b501LegacySequenceOriginal), b501LegacySequenceControl})
			if _, b501LegacySequenceError = b501LegacySequenceValidate(b501LegacySequenceData); b501LegacySequenceError != nil {
				b501LegacySequenceT.Fatalf("ab34 distinct-axis valid control rejected: %v", b501LegacySequenceError)
			}
			b501LegacySequenceRow, b501LegacySequenceFound, b501LegacySequenceError = b501LegacySequenceLookup(b501LegacySequenceData, b501LegacySequenceChanged)
			if b501LegacySequenceError != nil || !b501LegacySequenceFound || !b501LegacySequenceReflect.DeepEqual(b501LegacySequenceDecode(b501LegacySequenceT, b501LegacySequenceRow), b501LegacySequenceDecode(b501LegacySequenceT, b501LegacySequenceMarshal(b501LegacySequenceT, b501LegacySequenceControl))) {
				b501LegacySequenceT.Fatalf("ab34 distinct-axis lookup control failed: %v", b501LegacySequenceError)
			}
		})
	}
}

func TestAb34OptionalAbsentVersusZero(b501LegacySequenceT *b501LegacySequenceTesting.T) {
	for _, b501LegacySequenceIndex := range []int{3, 4, 8, 10} {
		b501LegacySequenceName := b501LegacySequenceReflect.TypeOf(B501LegacySequenceKey{}).Field(b501LegacySequenceIndex).Name
		b501LegacySequenceT.Run("ab34_"+b501LegacySequenceName, func(b501LegacySequenceT *b501LegacySequenceTesting.T) {
			b501LegacySequenceAbsent, b501LegacySequenceZero := b501LegacySequenceBaseKey(), b501LegacySequenceBaseKey()
			b501LegacySequenceAbsentField := b501LegacySequenceReflect.ValueOf(&b501LegacySequenceAbsent).Elem().Field(b501LegacySequenceIndex)
			b501LegacySequenceAbsentField.Set(b501LegacySequenceReflect.Zero(b501LegacySequenceAbsentField.Type()))
			b501LegacySequenceZeroField := b501LegacySequenceReflect.ValueOf(&b501LegacySequenceZero).Elem().Field(b501LegacySequenceIndex)
			b501LegacySequenceZeroField.Set(b501LegacySequenceReflect.New(b501LegacySequenceZeroField.Type().Elem()))
			b501LegacySequenceRow := b501LegacySequenceDecode(b501LegacySequenceT, b501LegacySequenceGet(b501LegacySequenceT, b501LegacySequenceBaseKey()))
			b501LegacySequenceOther := b501LegacySequenceDecode(b501LegacySequenceT, b501LegacySequenceGet(b501LegacySequenceT, b501LegacySequenceBaseKey()))
			b501LegacySequenceRow["ab34_identity"], b501LegacySequenceOther["ab34_identity"] = b501LegacySequenceAbsent, b501LegacySequenceZero
			b501LegacySequenceOther["ab34_id"] = "ab34_zero_control"
			b501LegacySequenceData := b501LegacySequenceMarshal(b501LegacySequenceT, []any{b501LegacySequenceRow, b501LegacySequenceOther})
			for _, b501LegacySequenceCase := range []struct {
				b501LegacySequenceKey B501LegacySequenceKey
				b501LegacySequenceID  string
			}{{b501LegacySequenceAbsent, "ab34_legacy_request"}, {b501LegacySequenceZero, "ab34_zero_control"}} {
				b501LegacySequenceGot, b501LegacySequenceFound, b501LegacySequenceError := b501LegacySequenceLookup(b501LegacySequenceData, b501LegacySequenceCase.b501LegacySequenceKey)
				if b501LegacySequenceError != nil || !b501LegacySequenceFound || b501LegacySequenceDecode(b501LegacySequenceT, b501LegacySequenceGot)["ab34_id"] != b501LegacySequenceCase.b501LegacySequenceID {
					b501LegacySequenceT.Fatalf("ab34 absent/zero collision: %v", b501LegacySequenceError)
				}
			}
		})
	}
}

func TestAb34ValidationErrors(b501LegacySequenceT *b501LegacySequenceTesting.T) {
	b501LegacySequenceRow := b501LegacySequenceGet(b501LegacySequenceT, b501LegacySequenceBaseKey())
	b501LegacySequenceDuplicate := b501LegacySequenceMarshal(b501LegacySequenceT, []b501LegacySequenceJSON.RawMessage{b501LegacySequenceRow, b501LegacySequenceRow})
	for _, b501LegacySequenceData := range [][]byte{[]byte(`[`), []byte(`{}`), []byte(`null`), []byte(`[null]`), []byte(`[{}]`), b501LegacySequenceDuplicate, []byte(`[{}] trailing`)} {
		if _, b501LegacySequenceError := b501LegacySequenceValidate(b501LegacySequenceData); b501LegacySequenceError == nil {
			b501LegacySequenceT.Errorf("ab34 accepted invalid data %s", b501LegacySequenceData)
		}
		if b501LegacySequenceGot, b501LegacySequenceFound, b501LegacySequenceError := b501LegacySequenceLookup(b501LegacySequenceData, b501LegacySequenceBaseKey()); b501LegacySequenceError == nil || b501LegacySequenceFound || b501LegacySequenceGot != nil {
			b501LegacySequenceT.Errorf("ab34 lookup failed open for %s", b501LegacySequenceData)
		}
	}
	b501LegacySequenceObject := b501LegacySequenceDecode(b501LegacySequenceT, b501LegacySequenceRow)
	for b501LegacySequenceAxis := range b501LegacySequenceObject["ab34_identity"].(map[string]any) {
		b501LegacySequenceT.Run("ab34_missing_"+b501LegacySequenceAxis, func(b501LegacySequenceT *b501LegacySequenceTesting.T) {
			b501LegacySequenceControl := b501LegacySequenceDecode(b501LegacySequenceT, b501LegacySequenceRow)
			delete(b501LegacySequenceControl["ab34_identity"].(map[string]any), b501LegacySequenceAxis)
			if _, b501LegacySequenceError := b501LegacySequenceValidate(b501LegacySequenceMarshal(b501LegacySequenceT, []any{b501LegacySequenceControl})); b501LegacySequenceError == nil {
				b501LegacySequenceT.Fatal("ab34 incomplete identity admitted")
			}
		})
	}
	for _, b501LegacySequenceBad := range []string{"null", "-1", "256", `"B501"`} {
		b501LegacySequenceData := b501LegacySequenceStrings.Replace(string(b501LegacySequenceRow), `"ab34_primary": 181`, `"ab34_primary": `+b501LegacySequenceBad, 1)
		if _, b501LegacySequenceError := b501LegacySequenceValidate([]byte("[" + b501LegacySequenceData + "]")); b501LegacySequenceError == nil {
			b501LegacySequenceT.Errorf("ab34 invalid primary admitted: %s", b501LegacySequenceBad)
		}
	}
}

func TestAb34IndependentBuffers(b501LegacySequenceT *b501LegacySequenceTesting.T) {
	b501LegacySequenceEntries, b501LegacySequenceError := B501LegacySequenceEntries()
	if b501LegacySequenceError != nil {
		b501LegacySequenceT.Fatal(b501LegacySequenceError)
	}
	b501LegacySequenceSnapshot := b501LegacySequenceBytes.Clone(b501LegacySequenceEntries)
	b501LegacySequenceRow := b501LegacySequenceGet(b501LegacySequenceT, b501LegacySequenceBaseKey())
	b501LegacySequenceRowSnapshot := b501LegacySequenceBytes.Clone(b501LegacySequenceRow)
	b501LegacySequenceSecondRow := b501LegacySequenceGet(b501LegacySequenceT, b501LegacySequenceBaseKey())
	for b501LegacySequenceIndex := range b501LegacySequenceRow {
		b501LegacySequenceRow[b501LegacySequenceIndex] ^= 255
	}
	if !b501LegacySequenceBytes.Equal(b501LegacySequenceSecondRow, b501LegacySequenceRowSnapshot) || !b501LegacySequenceBytes.Equal(b501LegacySequenceEntries, b501LegacySequenceSnapshot) || !b501LegacySequenceBytes.Equal(b501LegacySequenceGet(b501LegacySequenceT, b501LegacySequenceBaseKey()), b501LegacySequenceRowSnapshot) {
		b501LegacySequenceT.Fatal("ab34 row memory shared")
	}
	for b501LegacySequenceIndex := range b501LegacySequenceEntries {
		b501LegacySequenceEntries[b501LegacySequenceIndex] ^= 255
	}
	b501LegacySequenceLater, b501LegacySequenceError := B501LegacySequenceEntries()
	if b501LegacySequenceError != nil || !b501LegacySequenceBytes.Equal(b501LegacySequenceLater, b501LegacySequenceSnapshot) || !b501LegacySequenceBytes.Equal(b501LegacySequenceGet(b501LegacySequenceT, b501LegacySequenceBaseKey()), b501LegacySequenceRowSnapshot) {
		b501LegacySequenceT.Fatal("ab34 entries memory shared")
	}
	b501LegacySequenceInput := b501LegacySequenceBytes.Clone(b501LegacySequenceSnapshot)
	b501LegacySequenceOwned, b501LegacySequenceFound, b501LegacySequenceError := b501LegacySequenceLookup(b501LegacySequenceInput, b501LegacySequenceBaseKey())
	if b501LegacySequenceError != nil || !b501LegacySequenceFound {
		b501LegacySequenceT.Fatal("ab34 seam control failed")
	}
	for b501LegacySequenceIndex := range b501LegacySequenceInput {
		b501LegacySequenceInput[b501LegacySequenceIndex] ^= 255
	}
	if !b501LegacySequenceBytes.Equal(b501LegacySequenceOwned, b501LegacySequenceRowSnapshot) {
		b501LegacySequenceT.Fatal("ab34 lookup retains caller input")
	}
}
