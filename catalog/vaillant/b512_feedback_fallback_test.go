package vaillant

import (
	b512FeedbackFallbackBytes "bytes"
	b512FeedbackFallbackJSON "encoding/json"
	b512FeedbackFallbackReflect "reflect"
	b512FeedbackFallbackTesting "testing"
)

type b512FeedbackFallbackRow struct {
	B512FeedbackFallbackID             string                  `json:"id"`
	B512FeedbackFallbackKey            B512FeedbackFallbackKey `json:"key"`
	B512FeedbackFallbackGrade          string                  `json:"grade"`
	B512FeedbackFallbackPosition       int                     `json:"logical_position"`
	B512FeedbackFallbackType           string                  `json:"type"`
	B512FeedbackFallbackPhysicalOffset string                  `json:"physical_offset"`
	B512FeedbackFallbackRawForm        map[string]string       `json:"raw_form"`
	B512FeedbackFallbackRawFallback    string                  `json:"raw_fallback"`
	B512FeedbackFallbackValues         map[string]string       `json:"values"`
	B512FeedbackFallbackOther          string                  `json:"other_values"`
	B512FeedbackFallbackPolarity       string                  `json:"polarity"`
	B512FeedbackFallbackFeedback       struct {
		B512FeedbackFallbackAbsent        string `json:"actual_signal_absent"`
		B512FeedbackFallbackPresent       string `json:"actual_signal_present"`
		B512FeedbackFallbackEqual         bool   `json:"equal_to_setpoint_establishes_actual"`
		B512FeedbackFallbackDiscriminator string `json:"wire_discriminator"`
		B512FeedbackFallbackPhysical      string `json:"physical_value"`
	} `json:"feedback"`
	B512FeedbackFallbackDefault     bool                                `json:"boolean_default"`
	B512FeedbackFallbackRawDefault  b512FeedbackFallbackJSON.RawMessage `json:"raw_default"`
	B512FeedbackFallbackEncoding    string                              `json:"raw_default_encoding"`
	B512FeedbackFallbackNative      string                              `json:"native_support"`
	B512FeedbackFallbackEffects     string                              `json:"physical_effects"`
	B512FeedbackFallbackPersistence string                              `json:"persistence"`
	B512FeedbackFallbackDelivery    string                              `json:"delivery"`
}

func b512FeedbackFallbackRead(b512FeedbackFallbackT *b512FeedbackFallbackTesting.T) ([]b512FeedbackFallbackJSON.RawMessage, []b512FeedbackFallbackRow) {
	b512FeedbackFallbackT.Helper()
	b512FeedbackFallbackData, b512FeedbackFallbackError := B512FeedbackFallbackEntries()
	if b512FeedbackFallbackError != nil {
		b512FeedbackFallbackT.Fatal(b512FeedbackFallbackError)
	}
	var b512FeedbackFallbackRaw []b512FeedbackFallbackJSON.RawMessage
	var b512FeedbackFallbackRows []b512FeedbackFallbackRow
	if b512FeedbackFallbackError = b512FeedbackFallbackJSON.Unmarshal(b512FeedbackFallbackData, &b512FeedbackFallbackRaw); b512FeedbackFallbackError != nil {
		b512FeedbackFallbackT.Fatal(b512FeedbackFallbackError)
	}
	if b512FeedbackFallbackError = b512FeedbackFallbackJSON.Unmarshal(b512FeedbackFallbackData, &b512FeedbackFallbackRows); b512FeedbackFallbackError != nil {
		b512FeedbackFallbackT.Fatal(b512FeedbackFallbackError)
	}
	return b512FeedbackFallbackRaw, b512FeedbackFallbackRows
}

func TestAb18FiniteMetadata(b512FeedbackFallbackT *b512FeedbackFallbackTesting.T) {
	b512FeedbackFallbackRaw, b512FeedbackFallbackRows := b512FeedbackFallbackRead(b512FeedbackFallbackT)
	if len(b512FeedbackFallbackRows) != 12 {
		b512FeedbackFallbackT.Fatalf("ab18: rows = %d", len(b512FeedbackFallbackRows))
	}
	b512FeedbackFallbackFields := []string{"secondary-heating-pump", "circulation-pump", "solar-pump", "error-relay", "loading-pump", "fume-flap"}
	b512FeedbackFallbackRequest := []byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26}
	b512FeedbackFallbackPositive := []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36}
	for b512FeedbackFallbackSide, b512FeedbackFallbackDirection := range []string{"request", "positive"} {
		for b512FeedbackFallbackIndex, b512FeedbackFallbackField := range b512FeedbackFallbackFields {
			b512FeedbackFallbackT.Run("ab18-"+b512FeedbackFallbackDirection+"-"+b512FeedbackFallbackField, func(b512FeedbackFallbackT *b512FeedbackFallbackTesting.T) {
				b512FeedbackFallbackOffset := b512FeedbackFallbackSide*6 + b512FeedbackFallbackIndex
				b512FeedbackFallbackRow := b512FeedbackFallbackRows[b512FeedbackFallbackOffset]
				b512FeedbackFallbackSecondary, b512FeedbackFallbackOperation := uint8(18), uint8(1)
				b512FeedbackFallbackExpectedKey := B512FeedbackFallbackKey{"B512", "OP01", 181, &b512FeedbackFallbackSecondary, &b512FeedbackFallbackOperation, b512FeedbackFallbackDirection, b512FeedbackFallbackField, "none", nil, "none", nil}
				if !b512FeedbackFallbackReflect.DeepEqual(b512FeedbackFallbackRow.B512FeedbackFallbackKey, b512FeedbackFallbackExpectedKey) || b512FeedbackFallbackRow.B512FeedbackFallbackID != "ab18-"+b512FeedbackFallbackDirection+"-"+b512FeedbackFallbackField {
					b512FeedbackFallbackT.Fatal("ab18: identity")
				}
				b512FeedbackFallbackFoundRow, b512FeedbackFallbackFound, b512FeedbackFallbackError := B512FeedbackFallbackLookup(b512FeedbackFallbackExpectedKey)
				if b512FeedbackFallbackError != nil || !b512FeedbackFallbackFound || !b512FeedbackFallbackBytes.Equal(b512FeedbackFallbackFoundRow, b512FeedbackFallbackRaw[b512FeedbackFallbackOffset]) {
					b512FeedbackFallbackT.Fatal("ab18: whole row lookup", b512FeedbackFallbackError)
				}
				if b512FeedbackFallbackRow.B512FeedbackFallbackPosition != b512FeedbackFallbackIndex+3 || b512FeedbackFallbackRow.B512FeedbackFallbackType != "u8" || b512FeedbackFallbackRow.B512FeedbackFallbackGrade != "Hypothesis" {
					b512FeedbackFallbackT.Fatal("ab18: independent position/type/grade")
				}
				if b512FeedbackFallbackRow.B512FeedbackFallbackPhysicalOffset != "Unknown" || b512FeedbackFallbackRow.B512FeedbackFallbackPolarity != "Unknown" || b512FeedbackFallbackRow.B512FeedbackFallbackNative != "Unknown" || b512FeedbackFallbackRow.B512FeedbackFallbackEffects != "Unknown" || b512FeedbackFallbackRow.B512FeedbackFallbackPersistence != "Unknown" || b512FeedbackFallbackRow.B512FeedbackFallbackDelivery != "Unknown" {
					b512FeedbackFallbackT.Fatal("ab18: qualification")
				}
				if !b512FeedbackFallbackRow.B512FeedbackFallbackDefault || string(b512FeedbackFallbackRow.B512FeedbackFallbackRawDefault) != "null" || b512FeedbackFallbackRow.B512FeedbackFallbackEncoding != "Unknown" {
					b512FeedbackFallbackT.Fatal("ab18: boolean default is not raw 01")
				}
				if b512FeedbackFallbackRow.B512FeedbackFallbackRawForm["type"] != "u8" || b512FeedbackFallbackRow.B512FeedbackFallbackRawForm["retention"] != "preserve each independent raw token" || b512FeedbackFallbackRow.B512FeedbackFallbackRawFallback != "retain raw token and candidate meaning independently" {
					b512FeedbackFallbackT.Fatal("ab18: raw retention")
				}
				if b512FeedbackFallbackRow.B512FeedbackFallbackFeedback.B512FeedbackFallbackAbsent != "setpoint-fallback" || b512FeedbackFallbackRow.B512FeedbackFallbackFeedback.B512FeedbackFallbackPresent != "actual-signal" || b512FeedbackFallbackRow.B512FeedbackFallbackFeedback.B512FeedbackFallbackEqual || b512FeedbackFallbackRow.B512FeedbackFallbackFeedback.B512FeedbackFallbackDiscriminator != "Unknown" || b512FeedbackFallbackRow.B512FeedbackFallbackFeedback.B512FeedbackFallbackPhysical != "Unknown" {
					b512FeedbackFallbackT.Fatal("ab18: feedback provenance")
				}
				b512FeedbackFallbackTokens := append([]byte(nil), b512FeedbackFallbackRequest...)
				if b512FeedbackFallbackSide == 1 {
					b512FeedbackFallbackTokens = append([]byte(nil), b512FeedbackFallbackPositive...)
				}
				b512FeedbackFallbackBefore := append([]byte(nil), b512FeedbackFallbackTokens...)
				b512FeedbackFallbackChanged := append([]byte(nil), b512FeedbackFallbackTokens...)
				b512FeedbackFallbackChanged[b512FeedbackFallbackRow.B512FeedbackFallbackPosition-3] = 0x71 + byte(b512FeedbackFallbackIndex)
				for b512FeedbackFallbackOther := range b512FeedbackFallbackChanged {
					b512FeedbackFallbackExpectedToken := b512FeedbackFallbackBefore[b512FeedbackFallbackOther]
					if b512FeedbackFallbackOther == b512FeedbackFallbackIndex {
						b512FeedbackFallbackExpectedToken = 0x71 + byte(b512FeedbackFallbackIndex)
					}
					if b512FeedbackFallbackChanged[b512FeedbackFallbackOther] != b512FeedbackFallbackExpectedToken {
						b512FeedbackFallbackT.Fatal("ab18: field selection aliases another position")
					}
				}
				for _, b512FeedbackFallbackActualPresent := range []bool{false, true} {
					b512FeedbackFallbackSetpoint := b512FeedbackFallbackBefore[b512FeedbackFallbackIndex]
					b512FeedbackFallbackFeedbackToken := b512FeedbackFallbackSetpoint
					b512FeedbackFallbackSource := b512FeedbackFallbackRow.B512FeedbackFallbackFeedback.B512FeedbackFallbackAbsent
					b512FeedbackFallbackExpectedSource := "setpoint-fallback"
					if b512FeedbackFallbackActualPresent {
						b512FeedbackFallbackSource = b512FeedbackFallbackRow.B512FeedbackFallbackFeedback.B512FeedbackFallbackPresent
						b512FeedbackFallbackExpectedSource = "actual-signal"
					}
					if b512FeedbackFallbackSource != b512FeedbackFallbackExpectedSource || b512FeedbackFallbackFeedbackToken != b512FeedbackFallbackSetpoint || b512FeedbackFallbackRow.B512FeedbackFallbackFeedback.B512FeedbackFallbackEqual {
						b512FeedbackFallbackT.Fatal("ab18: equal tokens preserve source distinction")
					}
				}
				if b512FeedbackFallbackIndex != 3 && b512FeedbackFallbackIndex != 5 {
					b512FeedbackFallbackExpected := "leave-unchanged"
					if b512FeedbackFallbackSide == 1 {
						b512FeedbackFallbackExpected = "function-unavailable"
					}
					if len(b512FeedbackFallbackRow.B512FeedbackFallbackValues) != 1 || b512FeedbackFallbackRow.B512FeedbackFallbackValues["FF"] != b512FeedbackFallbackExpected || b512FeedbackFallbackRow.B512FeedbackFallbackOther != "retain-raw" {
						b512FeedbackFallbackT.Fatal("ab18: direction-specific FF")
					}
					// The metadata selects a synthetic field; request FF retains that field's previous token.
					b512FeedbackFallbackTokens[b512FeedbackFallbackRow.B512FeedbackFallbackPosition-3] = 0xff
					if b512FeedbackFallbackRow.B512FeedbackFallbackValues["FF"] == "leave-unchanged" {
						b512FeedbackFallbackTokens[b512FeedbackFallbackRow.B512FeedbackFallbackPosition-3] = b512FeedbackFallbackBefore[b512FeedbackFallbackIndex]
					}
					for b512FeedbackFallbackOther := range b512FeedbackFallbackTokens {
						b512FeedbackFallbackExpectedToken := b512FeedbackFallbackBefore[b512FeedbackFallbackOther]
						if b512FeedbackFallbackOther == b512FeedbackFallbackIndex && b512FeedbackFallbackSide == 1 {
							b512FeedbackFallbackExpectedToken = 0xff
						}
						if b512FeedbackFallbackTokens[b512FeedbackFallbackOther] != b512FeedbackFallbackExpectedToken {
							b512FeedbackFallbackT.Fatal("ab18: independent raw tokens")
						}
					}
				} else {
					if len(b512FeedbackFallbackRow.B512FeedbackFallbackValues) != 2 || b512FeedbackFallbackRow.B512FeedbackFallbackValues["0F"] != "state-0F" || b512FeedbackFallbackRow.B512FeedbackFallbackValues["F0"] != "state-F0" || b512FeedbackFallbackRow.B512FeedbackFallbackOther != "no-action" {
						b512FeedbackFallbackT.Fatal("ab18: declared states without polarity")
					}
					for _, b512FeedbackFallbackToken := range []string{"00", "01", "FF"} {
						if _, b512FeedbackFallbackAction := b512FeedbackFallbackRow.B512FeedbackFallbackValues[b512FeedbackFallbackToken]; b512FeedbackFallbackAction {
							b512FeedbackFallbackT.Fatal("ab18: unexpected action", b512FeedbackFallbackToken)
						}
					}
				}
			})
		}
	}
}

func TestAb18ExactAxes(b512FeedbackFallbackT *b512FeedbackFallbackTesting.T) {
	b512FeedbackFallbackRaw, b512FeedbackFallbackRows := b512FeedbackFallbackRead(b512FeedbackFallbackT)
	for b512FeedbackFallbackAxis := 0; b512FeedbackFallbackAxis < 11; b512FeedbackFallbackAxis++ {
		b512FeedbackFallbackName := b512FeedbackFallbackReflect.TypeOf(B512FeedbackFallbackKey{}).Field(b512FeedbackFallbackAxis).Name
		b512FeedbackFallbackT.Run(b512FeedbackFallbackName, func(b512FeedbackFallbackT *b512FeedbackFallbackTesting.T) {
			b512FeedbackFallbackKey := b512FeedbackFallbackRows[0].B512FeedbackFallbackKey
			b512FeedbackFallbackValue := b512FeedbackFallbackReflect.ValueOf(&b512FeedbackFallbackKey).Elem().Field(b512FeedbackFallbackAxis)
			switch b512FeedbackFallbackValue.Kind() {
			case b512FeedbackFallbackReflect.String:
				b512FeedbackFallbackValue.SetString(b512FeedbackFallbackValue.String() + "-ab18-control")
			case b512FeedbackFallbackReflect.Uint8:
				b512FeedbackFallbackValue.SetUint(b512FeedbackFallbackValue.Uint() + 1)
			case b512FeedbackFallbackReflect.Ptr:
				b512FeedbackFallbackPointer := b512FeedbackFallbackReflect.New(b512FeedbackFallbackValue.Type().Elem())
				if !b512FeedbackFallbackValue.IsNil() {
					b512FeedbackFallbackPointer.Elem().SetUint(b512FeedbackFallbackValue.Elem().Uint() + 1)
				}
				b512FeedbackFallbackValue.Set(b512FeedbackFallbackPointer)
			default:
				b512FeedbackFallbackT.Fatal("ab18: unhandled axis")
			}
			if _, b512FeedbackFallbackFound, b512FeedbackFallbackError := B512FeedbackFallbackLookup(b512FeedbackFallbackKey); b512FeedbackFallbackError != nil || b512FeedbackFallbackFound {
				b512FeedbackFallbackT.Fatal("ab18: public exact miss", b512FeedbackFallbackError)
			}
			var b512FeedbackFallbackControl map[string]b512FeedbackFallbackJSON.RawMessage
			if b512FeedbackFallbackError := b512FeedbackFallbackJSON.Unmarshal(b512FeedbackFallbackRaw[0], &b512FeedbackFallbackControl); b512FeedbackFallbackError != nil {
				b512FeedbackFallbackT.Fatal(b512FeedbackFallbackError)
			}
			b512FeedbackFallbackControl["key"], _ = b512FeedbackFallbackJSON.Marshal(b512FeedbackFallbackKey)
			b512FeedbackFallbackControlRaw, _ := b512FeedbackFallbackJSON.Marshal(b512FeedbackFallbackControl)
			b512FeedbackFallbackFixture, _ := b512FeedbackFallbackJSON.Marshal([]b512FeedbackFallbackJSON.RawMessage{b512FeedbackFallbackRaw[0], b512FeedbackFallbackControlRaw})
			for _, b512FeedbackFallbackExpected := range []struct {
				B512FeedbackFallbackKey B512FeedbackFallbackKey
				B512FeedbackFallbackRaw b512FeedbackFallbackJSON.RawMessage
			}{{b512FeedbackFallbackRows[0].B512FeedbackFallbackKey, b512FeedbackFallbackRaw[0]}, {b512FeedbackFallbackKey, b512FeedbackFallbackControlRaw}} {
				b512FeedbackFallbackResult, b512FeedbackFallbackFound, b512FeedbackFallbackError := b512FeedbackFallbackLookup(b512FeedbackFallbackFixture, b512FeedbackFallbackExpected.B512FeedbackFallbackKey)
				var b512FeedbackFallbackCompact b512FeedbackFallbackBytes.Buffer
				if b512FeedbackFallbackError := b512FeedbackFallbackJSON.Compact(&b512FeedbackFallbackCompact, b512FeedbackFallbackExpected.B512FeedbackFallbackRaw); b512FeedbackFallbackError != nil {
					b512FeedbackFallbackT.Fatal(b512FeedbackFallbackError)
				}
				if b512FeedbackFallbackError != nil || !b512FeedbackFallbackFound || !b512FeedbackFallbackBytes.Equal(b512FeedbackFallbackResult, b512FeedbackFallbackCompact.Bytes()) {
					b512FeedbackFallbackT.Fatal("ab18: independent valid control", b512FeedbackFallbackError)
				}
			}
		})
	}
}

func TestAb18OptionalPresence(b512FeedbackFallbackT *b512FeedbackFallbackTesting.T) {
	b512FeedbackFallbackRaw, b512FeedbackFallbackRows := b512FeedbackFallbackRead(b512FeedbackFallbackT)
	for _, b512FeedbackFallbackAxis := range []int{3, 4, 8, 10} {
		b512FeedbackFallbackT.Run(b512FeedbackFallbackReflect.TypeOf(B512FeedbackFallbackKey{}).Field(b512FeedbackFallbackAxis).Name, func(b512FeedbackFallbackT *b512FeedbackFallbackTesting.T) {
			b512FeedbackFallbackKeys := []B512FeedbackFallbackKey{b512FeedbackFallbackRows[0].B512FeedbackFallbackKey, b512FeedbackFallbackRows[0].B512FeedbackFallbackKey}
			b512FeedbackFallbackReflect.ValueOf(&b512FeedbackFallbackKeys[0]).Elem().Field(b512FeedbackFallbackAxis).SetZero()
			b512FeedbackFallbackValue := b512FeedbackFallbackReflect.ValueOf(&b512FeedbackFallbackKeys[1]).Elem().Field(b512FeedbackFallbackAxis)
			b512FeedbackFallbackValue.Set(b512FeedbackFallbackReflect.New(b512FeedbackFallbackValue.Type().Elem()))
			var b512FeedbackFallbackControls []b512FeedbackFallbackJSON.RawMessage
			for _, b512FeedbackFallbackKey := range b512FeedbackFallbackKeys {
				var b512FeedbackFallbackControl map[string]b512FeedbackFallbackJSON.RawMessage
				if b512FeedbackFallbackError := b512FeedbackFallbackJSON.Unmarshal(b512FeedbackFallbackRaw[0], &b512FeedbackFallbackControl); b512FeedbackFallbackError != nil {
					b512FeedbackFallbackT.Fatal(b512FeedbackFallbackError)
				}
				b512FeedbackFallbackControl["key"], _ = b512FeedbackFallbackJSON.Marshal(b512FeedbackFallbackKey)
				b512FeedbackFallbackControlRaw, _ := b512FeedbackFallbackJSON.Marshal(b512FeedbackFallbackControl)
				b512FeedbackFallbackControls = append(b512FeedbackFallbackControls, b512FeedbackFallbackControlRaw)
			}
			b512FeedbackFallbackFixture, _ := b512FeedbackFallbackJSON.Marshal(b512FeedbackFallbackControls)
			for b512FeedbackFallbackIndex, b512FeedbackFallbackKey := range b512FeedbackFallbackKeys {
				b512FeedbackFallbackResult, b512FeedbackFallbackFound, b512FeedbackFallbackError := b512FeedbackFallbackLookup(b512FeedbackFallbackFixture, b512FeedbackFallbackKey)
				if b512FeedbackFallbackError != nil || !b512FeedbackFallbackFound || !b512FeedbackFallbackBytes.Equal(b512FeedbackFallbackResult, b512FeedbackFallbackControls[b512FeedbackFallbackIndex]) {
					b512FeedbackFallbackT.Fatal("ab18: absent/zero", b512FeedbackFallbackError)
				}
			}
		})
	}
}

func TestAb18Validation(b512FeedbackFallbackT *b512FeedbackFallbackTesting.T) {
	b512FeedbackFallbackRaw, b512FeedbackFallbackRows := b512FeedbackFallbackRead(b512FeedbackFallbackT)
	b512FeedbackFallbackDuplicate, _ := b512FeedbackFallbackJSON.Marshal([]b512FeedbackFallbackJSON.RawMessage{b512FeedbackFallbackRaw[0], b512FeedbackFallbackRaw[0]})
	for _, b512FeedbackFallbackInput := range [][]byte{[]byte("["), []byte("{}"), []byte("null"), []byte("[null]"), []byte("[{}]"), b512FeedbackFallbackDuplicate, []byte(`[{"key":{"family":"B512"}}]`)} {
		if _, b512FeedbackFallbackError := b512FeedbackFallbackValidate(b512FeedbackFallbackInput); b512FeedbackFallbackError == nil {
			b512FeedbackFallbackT.Fatal("ab18: accepted invalid data", string(b512FeedbackFallbackInput))
		}
		if b512FeedbackFallbackResult, b512FeedbackFallbackFound, b512FeedbackFallbackError := b512FeedbackFallbackLookup(b512FeedbackFallbackInput, b512FeedbackFallbackRows[0].B512FeedbackFallbackKey); b512FeedbackFallbackError == nil || b512FeedbackFallbackFound || b512FeedbackFallbackResult != nil {
			b512FeedbackFallbackT.Fatal("ab18: fail closed")
		}
	}
	// An error anywhere prevents returning an earlier valid matching row.
	b512FeedbackFallbackPoisoned := append(append([]byte{'['}, b512FeedbackFallbackRaw[0]...), []byte(",{}]")...)
	if _, b512FeedbackFallbackFound, b512FeedbackFallbackError := b512FeedbackFallbackLookup(b512FeedbackFallbackPoisoned, b512FeedbackFallbackRows[0].B512FeedbackFallbackKey); b512FeedbackFallbackError == nil || b512FeedbackFallbackFound {
		b512FeedbackFallbackT.Fatal("ab18: whole catalog validation")
	}
}

func TestAb18IndependentBuffers(b512FeedbackFallbackT *b512FeedbackFallbackTesting.T) {
	b512FeedbackFallbackRaw, b512FeedbackFallbackRows := b512FeedbackFallbackRead(b512FeedbackFallbackT)
	b512FeedbackFallbackEntries, b512FeedbackFallbackError := B512FeedbackFallbackEntries()
	if b512FeedbackFallbackError != nil {
		b512FeedbackFallbackT.Fatal(b512FeedbackFallbackError)
	}
	b512FeedbackFallbackExpected := append([]byte(nil), b512FeedbackFallbackEntries...)
	b512FeedbackFallbackOther, b512FeedbackFallbackError := B512FeedbackFallbackEntries()
	if b512FeedbackFallbackError != nil {
		b512FeedbackFallbackT.Fatal(b512FeedbackFallbackError)
	}
	for b512FeedbackFallbackIndex := range b512FeedbackFallbackEntries {
		b512FeedbackFallbackEntries[b512FeedbackFallbackIndex] = 'x'
	}
	b512FeedbackFallbackLater, b512FeedbackFallbackError := B512FeedbackFallbackEntries()
	if b512FeedbackFallbackError != nil || !b512FeedbackFallbackBytes.Equal(b512FeedbackFallbackOther, b512FeedbackFallbackExpected) || !b512FeedbackFallbackBytes.Equal(b512FeedbackFallbackLater, b512FeedbackFallbackExpected) {
		b512FeedbackFallbackT.Fatal("ab18: shared entries")
	}
	for b512FeedbackFallbackIndex, b512FeedbackFallbackRow := range b512FeedbackFallbackRows {
		b512FeedbackFallbackFirst, b512FeedbackFallbackFound, b512FeedbackFallbackError := B512FeedbackFallbackLookup(b512FeedbackFallbackRow.B512FeedbackFallbackKey)
		if b512FeedbackFallbackError != nil || !b512FeedbackFallbackFound {
			b512FeedbackFallbackT.Fatal(b512FeedbackFallbackError)
		}
		b512FeedbackFallbackSecond, b512FeedbackFallbackFound, b512FeedbackFallbackError := B512FeedbackFallbackLookup(b512FeedbackFallbackRow.B512FeedbackFallbackKey)
		if b512FeedbackFallbackError != nil || !b512FeedbackFallbackFound {
			b512FeedbackFallbackT.Fatal(b512FeedbackFallbackError)
		}
		for b512FeedbackFallbackByte := range b512FeedbackFallbackFirst {
			b512FeedbackFallbackFirst[b512FeedbackFallbackByte] = 'x'
		}
		b512FeedbackFallbackThird, b512FeedbackFallbackFound, b512FeedbackFallbackError := B512FeedbackFallbackLookup(b512FeedbackFallbackRow.B512FeedbackFallbackKey)
		if b512FeedbackFallbackError != nil || !b512FeedbackFallbackFound || !b512FeedbackFallbackBytes.Equal(b512FeedbackFallbackSecond, b512FeedbackFallbackRaw[b512FeedbackFallbackIndex]) || !b512FeedbackFallbackBytes.Equal(b512FeedbackFallbackThird, b512FeedbackFallbackRaw[b512FeedbackFallbackIndex]) {
			b512FeedbackFallbackT.Fatal("ab18: shared row")
		}
	}
}
