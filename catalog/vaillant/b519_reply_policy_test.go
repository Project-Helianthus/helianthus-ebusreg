package vaillant

import (
	b519ReplyPolicyBytes "bytes"
	b519ReplyPolicyJSON "encoding/json"
	b519ReplyPolicyReflect "reflect"
	b519ReplyPolicyTesting "testing"
)

type b519ReplyPolicyEntry struct {
	B519ReplyPolicyID        string             `json:"ab21_id"`
	B519ReplyPolicyKey       B519ReplyPolicyKey `json:"key"`
	B519ReplyPolicyGrade     string             `json:"ab21_grade"`
	B519ReplyPolicyPresence  string             `json:"ab21_metadata_presence"`
	B519ReplyPolicyContract  string             `json:"ab21_native_contract"`
	B519ReplyPolicy          string             `json:"ab21_reply_policy"`
	B519ReplyPolicyRecipient string             `json:"ab21_reply_recipient"`
	B519ReplyPolicyInitiator string             `json:"ab21_initiator_relation"`
	B519ReplyPolicyData      *struct {
		B519ReplyPolicyMinimum        int    `json:"ab21_minimum_length"`
		B519ReplyPolicyCapacity       any    `json:"ab21_peer_capacity"`
		B519ReplyPolicyBlock          bool   `json:"ab21_block_number_in_data"`
		B519ReplyPolicyRepresentation string `json:"ab21_representation"`
	} `json:"ab21_data"`
	B519ReplyPolicyQuery *uint8 `json:"ab21_crc_query_operation"`
	B519ReplyPolicyRaw   struct {
		B519ReplyPolicyService    []uint8 `json:"ab21_service"`
		B519ReplyPolicyPreserve   bool    `json:"ab21_preserve_raw"`
		B519ReplyPolicyUnresolved string  `json:"ab21_unresolved_bytes"`
	} `json:"ab21_raw_form"`
	B519ReplyPolicyUnknown map[string]string `json:"ab21_unknown"`
}

func b519ReplyPolicyPointer(b519ReplyPolicyValue uint8) *uint8 { return &b519ReplyPolicyValue }
func b519ReplyPolicyRegisterPointer(b519ReplyPolicyValue uint16) *uint16 {
	return &b519ReplyPolicyValue
}

func b519ReplyPolicyBaseKey() B519ReplyPolicyKey {
	return B519ReplyPolicyKey{B519ReplyPolicyFamily: "vaillant", B519ReplyPolicyVariant: "Unknown", B519ReplyPolicyPrimary: 181, B519ReplyPolicySecondary: b519ReplyPolicyPointer(25), B519ReplyPolicyOperation: b519ReplyPolicyPointer(3), B519ReplyPolicyDirection: "initiator-to-target", B519ReplyPolicyField: "reply-policy", B519ReplyPolicyAlternative: "mode-policy", B519ReplyPolicyInstanceScope: "none"}
}

func b519ReplyPolicyGet(b519ReplyPolicyT *b519ReplyPolicyTesting.T, b519ReplyPolicyKey B519ReplyPolicyKey) (b519ReplyPolicyJSON.RawMessage, b519ReplyPolicyEntry) {
	b519ReplyPolicyT.Helper()
	b519ReplyPolicyRow, b519ReplyPolicyFound, b519ReplyPolicyError := B519ReplyPolicyLookup(b519ReplyPolicyKey)
	if b519ReplyPolicyError != nil || !b519ReplyPolicyFound {
		b519ReplyPolicyT.Fatalf("lookup: found=%v error=%v", b519ReplyPolicyFound, b519ReplyPolicyError)
	}
	var b519ReplyPolicyResult b519ReplyPolicyEntry
	if b519ReplyPolicyError := b519ReplyPolicyJSON.Unmarshal(b519ReplyPolicyRow, &b519ReplyPolicyResult); b519ReplyPolicyError != nil {
		b519ReplyPolicyT.Fatal(b519ReplyPolicyError)
	}
	if !b519ReplyPolicyReflect.DeepEqual(b519ReplyPolicyResult.B519ReplyPolicyKey, b519ReplyPolicyKey) {
		b519ReplyPolicyT.Fatal("whole row identity mismatch")
	}
	return b519ReplyPolicyRow, b519ReplyPolicyResult
}

func TestAb21Modes(b519ReplyPolicyT *b519ReplyPolicyTesting.T) {
	for _, b519ReplyPolicyCase := range []struct {
		B519ReplyPolicyMode, B519ReplyPolicy, B519ReplyPolicyRecipient, B519ReplyPolicyInitiator string
	}{
		{"initiator-to-target", "target-response", "none", "none"},
		{"initiator-to-initiator", "service-toward-initial-sender", "initial-sender", "initiator-associated-with-receiving-target"},
		{"broadcast", "retained-CRC-for-OP02", "none", "initiator-associated-with-receiving-target"},
	} {
		b519ReplyPolicyT.Run("ab21-"+b519ReplyPolicyCase.B519ReplyPolicyMode, func(b519ReplyPolicyT *b519ReplyPolicyTesting.T) {
			b519ReplyPolicyKey := b519ReplyPolicyBaseKey()
			b519ReplyPolicyKey.B519ReplyPolicyDirection = b519ReplyPolicyCase.B519ReplyPolicyMode
			_, b519ReplyPolicyResult := b519ReplyPolicyGet(b519ReplyPolicyT, b519ReplyPolicyKey)
			if b519ReplyPolicyResult.B519ReplyPolicy != b519ReplyPolicyCase.B519ReplyPolicy || b519ReplyPolicyResult.B519ReplyPolicyRecipient != b519ReplyPolicyCase.B519ReplyPolicyRecipient || b519ReplyPolicyResult.B519ReplyPolicyInitiator != b519ReplyPolicyCase.B519ReplyPolicyInitiator {
				b519ReplyPolicyT.Fatalf("wrong mode relations: %+v", b519ReplyPolicyResult)
			}
			if b519ReplyPolicyCase.B519ReplyPolicyMode == "broadcast" {
				if b519ReplyPolicyResult.B519ReplyPolicyQuery == nil || *b519ReplyPolicyResult.B519ReplyPolicyQuery != 2 {
					b519ReplyPolicyT.Fatal("missing later OP02 relation")
				}
			} else if b519ReplyPolicyResult.B519ReplyPolicyQuery != nil {
				b519ReplyPolicyT.Fatal("unexpected CRC query")
			}
			b519ReplyPolicyParticipants := map[string]string{"initial-sender": "sender-A", "receiving-target": "receiving-target-B", "initiator-associated-with-receiving-target": "associated-initiator-C"}
			if b519ReplyPolicyCase.B519ReplyPolicyMode == "initiator-to-initiator" && b519ReplyPolicyParticipants[b519ReplyPolicyResult.B519ReplyPolicyRecipient] != "sender-A" {
				b519ReplyPolicyT.Fatal("reply must relate to sender-A")
			}
			if b519ReplyPolicyCase.B519ReplyPolicyMode != "initiator-to-target" && b519ReplyPolicyParticipants[b519ReplyPolicyResult.B519ReplyPolicyInitiator] != "associated-initiator-C" {
				b519ReplyPolicyT.Fatal("associated initiator must remain C")
			}
			for _, b519ReplyPolicyData := range [][]byte{{}, {0, 255, 65}} {
				b519ReplyPolicyBlockToken := "ab21-fabricated-block"
				if b519ReplyPolicyResult.B519ReplyPolicyData == nil || b519ReplyPolicyResult.B519ReplyPolicyData.B519ReplyPolicyCapacity != "Unknown" || b519ReplyPolicyResult.B519ReplyPolicyData.B519ReplyPolicyMinimum != 0 || b519ReplyPolicyResult.B519ReplyPolicyData.B519ReplyPolicyBlock || b519ReplyPolicyResult.B519ReplyPolicyData.B519ReplyPolicyRepresentation != "opaque-bytes" {
					b519ReplyPolicyT.Fatal("data range, capacity or block separation changed")
				}
				if len(b519ReplyPolicyData) != 0 && len(b519ReplyPolicyData) != 3 {
					b519ReplyPolicyT.Fatal("invalid synthetic fixture")
				}
				if b519ReplyPolicyBytes.Contains(b519ReplyPolicyData, []byte(b519ReplyPolicyBlockToken)) {
					b519ReplyPolicyT.Fatal("block token in data")
				}
			}
		})
	}
}

func TestAb21AlternativesAndRaw(b519ReplyPolicyT *b519ReplyPolicyTesting.T) {
	b519ReplyPolicyData, b519ReplyPolicyError := B519ReplyPolicyEntries()
	if b519ReplyPolicyError != nil {
		b519ReplyPolicyT.Fatal(b519ReplyPolicyError)
	}
	var b519ReplyPolicyRows []b519ReplyPolicyJSON.RawMessage
	if b519ReplyPolicyError := b519ReplyPolicyJSON.Unmarshal(b519ReplyPolicyData, &b519ReplyPolicyRows); b519ReplyPolicyError != nil {
		b519ReplyPolicyT.Fatal(b519ReplyPolicyError)
	}
	if len(b519ReplyPolicyRows) != 7 {
		b519ReplyPolicyT.Fatalf("finite area: got %d rows", len(b519ReplyPolicyRows))
	}
	b519ReplyPolicyIDs := map[string]bool{}
	for _, b519ReplyPolicyRow := range b519ReplyPolicyRows {
		var b519ReplyPolicyExpected b519ReplyPolicyEntry
		if b519ReplyPolicyError := b519ReplyPolicyJSON.Unmarshal(b519ReplyPolicyRow, &b519ReplyPolicyExpected); b519ReplyPolicyError != nil {
			b519ReplyPolicyT.Fatal(b519ReplyPolicyError)
		}
		b519ReplyPolicyWhole, b519ReplyPolicyActual := b519ReplyPolicyGet(b519ReplyPolicyT, b519ReplyPolicyExpected.B519ReplyPolicyKey)
		if !b519ReplyPolicyBytes.Equal(b519ReplyPolicyWhole, b519ReplyPolicyRow) {
			b519ReplyPolicyT.Fatal("lookup must return the whole original row")
		}
		if b519ReplyPolicyActual.B519ReplyPolicyID == "" || b519ReplyPolicyIDs[b519ReplyPolicyActual.B519ReplyPolicyID] {
			b519ReplyPolicyT.Fatal("missing or duplicate ID")
		}
		b519ReplyPolicyIDs[b519ReplyPolicyActual.B519ReplyPolicyID] = true
		if b519ReplyPolicyActual.B519ReplyPolicyGrade != "Hypothesis" || b519ReplyPolicyActual.B519ReplyPolicyPresence != "Proven" || b519ReplyPolicyActual.B519ReplyPolicyContract != "protocols/vaillant/ebus-vaillant-B519.md" || !b519ReplyPolicyActual.B519ReplyPolicyRaw.B519ReplyPolicyPreserve || b519ReplyPolicyActual.B519ReplyPolicyRaw.B519ReplyPolicyUnresolved != "retained" || !b519ReplyPolicyReflect.DeepEqual(b519ReplyPolicyActual.B519ReplyPolicyRaw.B519ReplyPolicyService, []uint8{181, 25, *b519ReplyPolicyActual.B519ReplyPolicyKey.B519ReplyPolicyOperation}) {
			b519ReplyPolicyT.Fatal("grade or raw retention changed")
		}
		for _, b519ReplyPolicyUnknown := range []string{"ab21_routing", "ab21_actual_acknowledgement", "ab21_storage_effects", "ab21_native_support", "ab21_physical_effects", "ab21_persistence", "ab21_delivery", "ab21_qualification"} {
			if b519ReplyPolicyActual.B519ReplyPolicyUnknown[b519ReplyPolicyUnknown] != "Unknown" {
				b519ReplyPolicyT.Fatalf("%s must remain Unknown", b519ReplyPolicyUnknown)
			}
		}
	}
	for _, b519ReplyPolicyMode := range []string{"initiator-to-initiator", "broadcast"} {
		for _, b519ReplyPolicyAlternative := range []string{"ACK-only", "CRC-positive"} {
			b519ReplyPolicyKey := b519ReplyPolicyBaseKey()
			b519ReplyPolicyKey.B519ReplyPolicyOperation = b519ReplyPolicyPointer(1)
			b519ReplyPolicyKey.B519ReplyPolicyDirection = b519ReplyPolicyMode
			b519ReplyPolicyKey.B519ReplyPolicyAlternative = b519ReplyPolicyAlternative
			_, b519ReplyPolicyActual := b519ReplyPolicyGet(b519ReplyPolicyT, b519ReplyPolicyKey)
			if b519ReplyPolicyActual.B519ReplyPolicy != b519ReplyPolicyAlternative || b519ReplyPolicyActual.B519ReplyPolicyData != nil || b519ReplyPolicyActual.B519ReplyPolicyInitiator != "initiator-associated-with-receiving-target" {
				b519ReplyPolicyT.Fatal("OP01 alternatives must coexist as policies")
			}
		}
	}
}

func TestAb21EveryAxis(b519ReplyPolicyT *b519ReplyPolicyTesting.T) {
	for _, b519ReplyPolicyCase := range []struct {
		B519ReplyPolicyName   string
		B519ReplyPolicyChange func(*B519ReplyPolicyKey)
	}{
		{"ab21-family", func(b519ReplyPolicyKey *B519ReplyPolicyKey) {
			b519ReplyPolicyKey.B519ReplyPolicyFamily = "ab21-other-family"
		}},
		{"ab21-variant", func(b519ReplyPolicyKey *B519ReplyPolicyKey) {
			b519ReplyPolicyKey.B519ReplyPolicyVariant = "ab21-other-variant"
		}},
		{"ab21-primary", func(b519ReplyPolicyKey *B519ReplyPolicyKey) { b519ReplyPolicyKey.B519ReplyPolicyPrimary = 180 }},
		{"ab21-secondary", func(b519ReplyPolicyKey *B519ReplyPolicyKey) {
			b519ReplyPolicyKey.B519ReplyPolicySecondary = b519ReplyPolicyPointer(24)
		}},
		{"ab21-operation", func(b519ReplyPolicyKey *B519ReplyPolicyKey) {
			b519ReplyPolicyKey.B519ReplyPolicyOperation = b519ReplyPolicyPointer(2)
		}},
		{"ab21-direction", func(b519ReplyPolicyKey *B519ReplyPolicyKey) { b519ReplyPolicyKey.B519ReplyPolicyDirection = "none" }},
		{"ab21-field", func(b519ReplyPolicyKey *B519ReplyPolicyKey) {
			b519ReplyPolicyKey.B519ReplyPolicyField = "ab21-other-field"
		}},
		{"ab21-alternative", func(b519ReplyPolicyKey *B519ReplyPolicyKey) { b519ReplyPolicyKey.B519ReplyPolicyAlternative = "none" }},
		{"ab21-group", func(b519ReplyPolicyKey *B519ReplyPolicyKey) {
			b519ReplyPolicyKey.B519ReplyPolicyGroup = b519ReplyPolicyPointer(1)
		}},
		{"ab21-instance-scope", func(b519ReplyPolicyKey *B519ReplyPolicyKey) {
			b519ReplyPolicyKey.B519ReplyPolicyInstanceScope = "dynamic"
		}},
		{"ab21-register", func(b519ReplyPolicyKey *B519ReplyPolicyKey) {
			b519ReplyPolicyKey.B519ReplyPolicyRegister = b519ReplyPolicyRegisterPointer(1)
		}},
	} {
		b519ReplyPolicyT.Run(b519ReplyPolicyCase.B519ReplyPolicyName, func(b519ReplyPolicyT *b519ReplyPolicyTesting.T) {
			b519ReplyPolicyKey := b519ReplyPolicyBaseKey()
			b519ReplyPolicyGet(b519ReplyPolicyT, b519ReplyPolicyKey)
			b519ReplyPolicyCase.B519ReplyPolicyChange(&b519ReplyPolicyKey)
			b519ReplyPolicyRow, b519ReplyPolicyFound, b519ReplyPolicyError := B519ReplyPolicyLookup(b519ReplyPolicyKey)
			if b519ReplyPolicyError != nil || b519ReplyPolicyFound || b519ReplyPolicyRow != nil {
				b519ReplyPolicyT.Fatalf("single axis must miss: found=%v error=%v", b519ReplyPolicyFound, b519ReplyPolicyError)
			}
		})
	}
}

func TestAb21OptionalZeroAndValidation(b519ReplyPolicyT *b519ReplyPolicyTesting.T) {
	b519ReplyPolicyRow, _ := b519ReplyPolicyGet(b519ReplyPolicyT, b519ReplyPolicyBaseKey())
	for _, b519ReplyPolicyAxis := range []string{"secondary", "operation", "group", "register"} {
		b519ReplyPolicyT.Run("ab21-"+b519ReplyPolicyAxis, func(b519ReplyPolicyT *b519ReplyPolicyTesting.T) {
			var b519ReplyPolicyObject map[string]b519ReplyPolicyJSON.RawMessage
			if b519ReplyPolicyError := b519ReplyPolicyJSON.Unmarshal(b519ReplyPolicyRow, &b519ReplyPolicyObject); b519ReplyPolicyError != nil {
				b519ReplyPolicyT.Fatal(b519ReplyPolicyError)
			}
			var b519ReplyPolicyAxes map[string]b519ReplyPolicyJSON.RawMessage
			if b519ReplyPolicyError := b519ReplyPolicyJSON.Unmarshal(b519ReplyPolicyObject["key"], &b519ReplyPolicyAxes); b519ReplyPolicyError != nil {
				b519ReplyPolicyT.Fatal(b519ReplyPolicyError)
			}
			var b519ReplyPolicyControls []b519ReplyPolicyJSON.RawMessage
			var b519ReplyPolicyKeys []B519ReplyPolicyKey
			for _, b519ReplyPolicyValue := range []string{"null", "0"} {
				b519ReplyPolicyAxes[b519ReplyPolicyAxis] = b519ReplyPolicyJSON.RawMessage(b519ReplyPolicyValue)
				b519ReplyPolicyObject["key"], _ = b519ReplyPolicyJSON.Marshal(b519ReplyPolicyAxes)
				b519ReplyPolicyControl, b519ReplyPolicyError := b519ReplyPolicyJSON.Marshal(b519ReplyPolicyObject)
				if b519ReplyPolicyError != nil {
					b519ReplyPolicyT.Fatal(b519ReplyPolicyError)
				}
				var b519ReplyPolicyKey B519ReplyPolicyKey
				if b519ReplyPolicyError := b519ReplyPolicyJSON.Unmarshal(b519ReplyPolicyObject["key"], &b519ReplyPolicyKey); b519ReplyPolicyError != nil {
					b519ReplyPolicyT.Fatal(b519ReplyPolicyError)
				}
				b519ReplyPolicyKeys = append(b519ReplyPolicyKeys, b519ReplyPolicyKey)
				b519ReplyPolicyControls = append(b519ReplyPolicyControls, b519ReplyPolicyControl)
			}
			b519ReplyPolicyCatalog, _ := b519ReplyPolicyJSON.Marshal(b519ReplyPolicyControls)
			if _, b519ReplyPolicyError := b519ReplyPolicyValidate(b519ReplyPolicyCatalog); b519ReplyPolicyError != nil {
				b519ReplyPolicyT.Fatal(b519ReplyPolicyError)
			}
			for b519ReplyPolicyIndex, b519ReplyPolicyKey := range b519ReplyPolicyKeys {
				b519ReplyPolicyGot, b519ReplyPolicyFound, b519ReplyPolicyError := b519ReplyPolicyLookup(b519ReplyPolicyCatalog, b519ReplyPolicyKey)
				if b519ReplyPolicyError != nil || !b519ReplyPolicyFound || !b519ReplyPolicyBytes.Equal(b519ReplyPolicyGot, b519ReplyPolicyControls[b519ReplyPolicyIndex]) {
					b519ReplyPolicyT.Fatal("absent and zero collapsed")
				}
				b519ReplyPolicySingle, _ := b519ReplyPolicyJSON.Marshal([]b519ReplyPolicyJSON.RawMessage{b519ReplyPolicyControls[b519ReplyPolicyIndex]})
				if _, b519ReplyPolicyFound, b519ReplyPolicyError := b519ReplyPolicyLookup(b519ReplyPolicySingle, b519ReplyPolicyKeys[1-b519ReplyPolicyIndex]); b519ReplyPolicyError != nil || b519ReplyPolicyFound {
					b519ReplyPolicyT.Fatal("opposite optional key must miss")
				}
			}
		})
	}
	b519ReplyPolicyDuplicate, _ := b519ReplyPolicyJSON.Marshal([]b519ReplyPolicyJSON.RawMessage{b519ReplyPolicyRow, b519ReplyPolicyRow})
	for _, b519ReplyPolicyInvalid := range [][]byte{[]byte("{"), []byte("null"), []byte("[]"), []byte("[null]"), []byte("[{}]"), []byte("{}"), append(append([]byte{}, b519ReplyPolicyDuplicate...), []byte(" true")...), b519ReplyPolicyDuplicate} {
		if _, b519ReplyPolicyError := b519ReplyPolicyValidate(b519ReplyPolicyInvalid); b519ReplyPolicyError == nil {
			b519ReplyPolicyT.Fatal("invalid embedded data admitted")
		}
		if b519ReplyPolicyGot, b519ReplyPolicyFound, b519ReplyPolicyError := b519ReplyPolicyLookup(b519ReplyPolicyInvalid, b519ReplyPolicyBaseKey()); b519ReplyPolicyError == nil || b519ReplyPolicyFound || b519ReplyPolicyGot != nil {
			b519ReplyPolicyT.Fatal("lookup must fail closed on invalid data")
		}
	}
}

func TestAb21IndependentBuffers(b519ReplyPolicyT *b519ReplyPolicyTesting.T) {
	b519ReplyPolicyEntries, b519ReplyPolicyError := B519ReplyPolicyEntries()
	if b519ReplyPolicyError != nil {
		b519ReplyPolicyT.Fatal(b519ReplyPolicyError)
	}
	b519ReplyPolicySaved := append([]byte(nil), b519ReplyPolicyEntries...)
	b519ReplyPolicyRow, _ := b519ReplyPolicyGet(b519ReplyPolicyT, b519ReplyPolicyBaseKey())
	b519ReplyPolicySavedRow := append([]byte(nil), b519ReplyPolicyRow...)
	b519ReplyPolicyOtherRow, _ := b519ReplyPolicyGet(b519ReplyPolicyT, b519ReplyPolicyBaseKey())
	for b519ReplyPolicyIndex := range b519ReplyPolicyEntries {
		b519ReplyPolicyEntries[b519ReplyPolicyIndex] = 'x'
	}
	for b519ReplyPolicyIndex := range b519ReplyPolicyRow {
		b519ReplyPolicyRow[b519ReplyPolicyIndex] = 'y'
	}
	b519ReplyPolicyNext, b519ReplyPolicyError := B519ReplyPolicyEntries()
	if b519ReplyPolicyError != nil || !b519ReplyPolicyBytes.Equal(b519ReplyPolicyNext, b519ReplyPolicySaved) {
		b519ReplyPolicyT.Fatal("entries share storage")
	}
	b519ReplyPolicyNextRow, _ := b519ReplyPolicyGet(b519ReplyPolicyT, b519ReplyPolicyBaseKey())
	if !b519ReplyPolicyBytes.Equal(b519ReplyPolicyOtherRow, b519ReplyPolicySavedRow) || !b519ReplyPolicyBytes.Equal(b519ReplyPolicyNextRow, b519ReplyPolicySavedRow) {
		b519ReplyPolicyT.Fatal("rows share storage")
	}
	b519ReplyPolicyCatalog := append([]byte(nil), b519ReplyPolicySaved...)
	b519ReplyPolicySeamRow, b519ReplyPolicyFound, b519ReplyPolicyError := b519ReplyPolicyLookup(b519ReplyPolicyCatalog, b519ReplyPolicyBaseKey())
	if b519ReplyPolicyError != nil || !b519ReplyPolicyFound {
		b519ReplyPolicyT.Fatal("seam lookup failed")
	}
	for b519ReplyPolicyIndex := range b519ReplyPolicyCatalog {
		b519ReplyPolicyCatalog[b519ReplyPolicyIndex] = 'z'
	}
	if !b519ReplyPolicyBytes.Equal(b519ReplyPolicySeamRow, b519ReplyPolicySavedRow) {
		b519ReplyPolicyT.Fatal("row aliases validation input")
	}
}
