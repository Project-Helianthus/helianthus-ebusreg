package vaillant

import (
	ab21bytes "bytes"
	ab21json "encoding/json"
	ab21reflect "reflect"
	ab21testing "testing"
)

type ab21Entry struct {
	Ab21ID        string  `json:"ab21_id"`
	Ab21Key       Ab21Key `json:"key"`
	Ab21Grade     string  `json:"ab21_grade"`
	Ab21Presence  string  `json:"ab21_metadata_presence"`
	Ab21Contract  string  `json:"ab21_native_contract"`
	Ab21Policy    string  `json:"ab21_reply_policy"`
	Ab21Recipient string  `json:"ab21_reply_recipient"`
	Ab21Initiator string  `json:"ab21_initiator_relation"`
	Ab21Data      *struct {
		Ab21Minimum        int    `json:"ab21_minimum_length"`
		Ab21Capacity       any    `json:"ab21_peer_capacity"`
		Ab21Block          bool   `json:"ab21_block_number_in_data"`
		Ab21Representation string `json:"ab21_representation"`
	} `json:"ab21_data"`
	Ab21Query *uint8 `json:"ab21_crc_query_operation"`
	Ab21Raw   struct {
		Ab21Service    []uint8 `json:"ab21_service"`
		Ab21Preserve   bool    `json:"ab21_preserve_raw"`
		Ab21Unresolved string  `json:"ab21_unresolved_bytes"`
	} `json:"ab21_raw_form"`
	Ab21Unknown map[string]string `json:"ab21_unknown"`
}

func ab21Pointer(ab21Value uint8) *uint8           { return &ab21Value }
func ab21RegisterPointer(ab21Value uint16) *uint16 { return &ab21Value }

func ab21BaseKey() Ab21Key {
	return Ab21Key{Ab21Family: "vaillant", Ab21Variant: "Unknown", Ab21Primary: 181, Ab21Secondary: ab21Pointer(25), Ab21Operation: ab21Pointer(3), Ab21Direction: "initiator-to-target", Ab21Field: "reply-policy", Ab21Alternative: "mode-policy", Ab21InstanceScope: "none"}
}

func ab21Get(ab21T *ab21testing.T, ab21Key Ab21Key) (ab21json.RawMessage, ab21Entry) {
	ab21T.Helper()
	ab21Row, ab21Found, ab21Error := Ab21Lookup(ab21Key)
	if ab21Error != nil || !ab21Found {
		ab21T.Fatalf("lookup: found=%v error=%v", ab21Found, ab21Error)
	}
	var ab21Result ab21Entry
	if ab21Error := ab21json.Unmarshal(ab21Row, &ab21Result); ab21Error != nil {
		ab21T.Fatal(ab21Error)
	}
	if !ab21reflect.DeepEqual(ab21Result.Ab21Key, ab21Key) {
		ab21T.Fatal("whole row identity mismatch")
	}
	return ab21Row, ab21Result
}

func TestAb21Modes(ab21T *ab21testing.T) {
	for _, ab21Case := range []struct {
		Ab21Mode, Ab21Policy, Ab21Recipient, Ab21Initiator string
	}{
		{"initiator-to-target", "target-response", "none", "none"},
		{"initiator-to-initiator", "service-toward-initial-sender", "initial-sender", "initiator-associated-with-receiving-target"},
		{"broadcast", "retained-CRC-for-OP02", "none", "initiator-associated-with-receiving-target"},
	} {
		ab21T.Run("ab21-"+ab21Case.Ab21Mode, func(ab21T *ab21testing.T) {
			ab21Key := ab21BaseKey()
			ab21Key.Ab21Direction = ab21Case.Ab21Mode
			_, ab21Result := ab21Get(ab21T, ab21Key)
			if ab21Result.Ab21Policy != ab21Case.Ab21Policy || ab21Result.Ab21Recipient != ab21Case.Ab21Recipient || ab21Result.Ab21Initiator != ab21Case.Ab21Initiator {
				ab21T.Fatalf("wrong mode relations: %+v", ab21Result)
			}
			if ab21Case.Ab21Mode == "broadcast" {
				if ab21Result.Ab21Query == nil || *ab21Result.Ab21Query != 2 {
					ab21T.Fatal("missing later OP02 relation")
				}
			} else if ab21Result.Ab21Query != nil {
				ab21T.Fatal("unexpected CRC query")
			}
			ab21Participants := map[string]string{"initial-sender": "sender-A", "receiving-target": "receiving-target-B", "initiator-associated-with-receiving-target": "associated-initiator-C"}
			if ab21Case.Ab21Mode == "initiator-to-initiator" && ab21Participants[ab21Result.Ab21Recipient] != "sender-A" {
				ab21T.Fatal("reply must relate to sender-A")
			}
			if ab21Case.Ab21Mode != "initiator-to-target" && ab21Participants[ab21Result.Ab21Initiator] != "associated-initiator-C" {
				ab21T.Fatal("associated initiator must remain C")
			}
			for _, ab21Data := range [][]byte{{}, {0, 255, 65}} {
				ab21BlockToken := "ab21-fabricated-block"
				if ab21Result.Ab21Data == nil || ab21Result.Ab21Data.Ab21Capacity != "Unknown" || ab21Result.Ab21Data.Ab21Minimum != 0 || ab21Result.Ab21Data.Ab21Block || ab21Result.Ab21Data.Ab21Representation != "opaque-bytes" {
					ab21T.Fatal("data range, capacity or block separation changed")
				}
				if len(ab21Data) != 0 && len(ab21Data) != 3 {
					ab21T.Fatal("invalid synthetic fixture")
				}
				if ab21bytes.Contains(ab21Data, []byte(ab21BlockToken)) {
					ab21T.Fatal("block token in data")
				}
			}
		})
	}
}

func TestAb21AlternativesAndRaw(ab21T *ab21testing.T) {
	ab21Data, ab21Error := Ab21Entries()
	if ab21Error != nil {
		ab21T.Fatal(ab21Error)
	}
	var ab21Rows []ab21json.RawMessage
	if ab21Error := ab21json.Unmarshal(ab21Data, &ab21Rows); ab21Error != nil {
		ab21T.Fatal(ab21Error)
	}
	if len(ab21Rows) != 7 {
		ab21T.Fatalf("finite area: got %d rows", len(ab21Rows))
	}
	ab21IDs := map[string]bool{}
	for _, ab21Row := range ab21Rows {
		var ab21Expected ab21Entry
		if ab21Error := ab21json.Unmarshal(ab21Row, &ab21Expected); ab21Error != nil {
			ab21T.Fatal(ab21Error)
		}
		ab21Whole, ab21Actual := ab21Get(ab21T, ab21Expected.Ab21Key)
		if !ab21bytes.Equal(ab21Whole, ab21Row) {
			ab21T.Fatal("lookup must return the whole original row")
		}
		if ab21Actual.Ab21ID == "" || ab21IDs[ab21Actual.Ab21ID] {
			ab21T.Fatal("missing or duplicate ID")
		}
		ab21IDs[ab21Actual.Ab21ID] = true
		if ab21Actual.Ab21Grade != "Hypothesis" || ab21Actual.Ab21Presence != "Proven" || ab21Actual.Ab21Contract != "protocols/vaillant/ebus-vaillant-B519.md" || !ab21Actual.Ab21Raw.Ab21Preserve || ab21Actual.Ab21Raw.Ab21Unresolved != "retained" || !ab21reflect.DeepEqual(ab21Actual.Ab21Raw.Ab21Service, []uint8{181, 25, *ab21Actual.Ab21Key.Ab21Operation}) {
			ab21T.Fatal("grade or raw retention changed")
		}
		for _, ab21Unknown := range []string{"ab21_routing", "ab21_actual_acknowledgement", "ab21_storage_effects", "ab21_native_support", "ab21_physical_effects", "ab21_persistence", "ab21_delivery", "ab21_qualification"} {
			if ab21Actual.Ab21Unknown[ab21Unknown] != "Unknown" {
				ab21T.Fatalf("%s must remain Unknown", ab21Unknown)
			}
		}
	}
	for _, ab21Mode := range []string{"initiator-to-initiator", "broadcast"} {
		for _, ab21Alternative := range []string{"ACK-only", "CRC-positive"} {
			ab21Key := ab21BaseKey()
			ab21Key.Ab21Operation = ab21Pointer(1)
			ab21Key.Ab21Direction = ab21Mode
			ab21Key.Ab21Alternative = ab21Alternative
			_, ab21Actual := ab21Get(ab21T, ab21Key)
			if ab21Actual.Ab21Policy != ab21Alternative || ab21Actual.Ab21Data != nil || ab21Actual.Ab21Initiator != "initiator-associated-with-receiving-target" {
				ab21T.Fatal("OP01 alternatives must coexist as policies")
			}
		}
	}
}

func TestAb21EveryAxis(ab21T *ab21testing.T) {
	for _, ab21Case := range []struct {
		Ab21Name   string
		Ab21Change func(*Ab21Key)
	}{
		{"ab21-family", func(ab21Key *Ab21Key) { ab21Key.Ab21Family = "ab21-other-family" }},
		{"ab21-variant", func(ab21Key *Ab21Key) { ab21Key.Ab21Variant = "ab21-other-variant" }},
		{"ab21-primary", func(ab21Key *Ab21Key) { ab21Key.Ab21Primary = 180 }},
		{"ab21-secondary", func(ab21Key *Ab21Key) { ab21Key.Ab21Secondary = ab21Pointer(24) }},
		{"ab21-operation", func(ab21Key *Ab21Key) { ab21Key.Ab21Operation = ab21Pointer(2) }},
		{"ab21-direction", func(ab21Key *Ab21Key) { ab21Key.Ab21Direction = "none" }},
		{"ab21-field", func(ab21Key *Ab21Key) { ab21Key.Ab21Field = "ab21-other-field" }},
		{"ab21-alternative", func(ab21Key *Ab21Key) { ab21Key.Ab21Alternative = "none" }},
		{"ab21-group", func(ab21Key *Ab21Key) { ab21Key.Ab21Group = ab21Pointer(1) }},
		{"ab21-instance-scope", func(ab21Key *Ab21Key) { ab21Key.Ab21InstanceScope = "dynamic" }},
		{"ab21-register", func(ab21Key *Ab21Key) { ab21Key.Ab21Register = ab21RegisterPointer(1) }},
	} {
		ab21T.Run(ab21Case.Ab21Name, func(ab21T *ab21testing.T) {
			ab21Key := ab21BaseKey()
			ab21Get(ab21T, ab21Key)
			ab21Case.Ab21Change(&ab21Key)
			ab21Row, ab21Found, ab21Error := Ab21Lookup(ab21Key)
			if ab21Error != nil || ab21Found || ab21Row != nil {
				ab21T.Fatalf("single axis must miss: found=%v error=%v", ab21Found, ab21Error)
			}
		})
	}
}

func TestAb21OptionalZeroAndValidation(ab21T *ab21testing.T) {
	ab21Row, _ := ab21Get(ab21T, ab21BaseKey())
	for _, ab21Axis := range []string{"secondary", "operation", "group", "register"} {
		ab21T.Run("ab21-"+ab21Axis, func(ab21T *ab21testing.T) {
			var ab21Object map[string]ab21json.RawMessage
			if ab21Error := ab21json.Unmarshal(ab21Row, &ab21Object); ab21Error != nil {
				ab21T.Fatal(ab21Error)
			}
			var ab21Axes map[string]ab21json.RawMessage
			if ab21Error := ab21json.Unmarshal(ab21Object["key"], &ab21Axes); ab21Error != nil {
				ab21T.Fatal(ab21Error)
			}
			var ab21Controls []ab21json.RawMessage
			var ab21Keys []Ab21Key
			for _, ab21Value := range []string{"null", "0"} {
				ab21Axes[ab21Axis] = ab21json.RawMessage(ab21Value)
				ab21Object["key"], _ = ab21json.Marshal(ab21Axes)
				ab21Control, ab21Error := ab21json.Marshal(ab21Object)
				if ab21Error != nil {
					ab21T.Fatal(ab21Error)
				}
				var ab21Key Ab21Key
				if ab21Error := ab21json.Unmarshal(ab21Object["key"], &ab21Key); ab21Error != nil {
					ab21T.Fatal(ab21Error)
				}
				ab21Keys = append(ab21Keys, ab21Key)
				ab21Controls = append(ab21Controls, ab21Control)
			}
			ab21Catalog, _ := ab21json.Marshal(ab21Controls)
			if _, ab21Error := ab21Validate(ab21Catalog); ab21Error != nil {
				ab21T.Fatal(ab21Error)
			}
			for ab21Index, ab21Key := range ab21Keys {
				ab21Got, ab21Found, ab21Error := ab21Lookup(ab21Catalog, ab21Key)
				if ab21Error != nil || !ab21Found || !ab21bytes.Equal(ab21Got, ab21Controls[ab21Index]) {
					ab21T.Fatal("absent and zero collapsed")
				}
				ab21Single, _ := ab21json.Marshal([]ab21json.RawMessage{ab21Controls[ab21Index]})
				if _, ab21Found, ab21Error := ab21Lookup(ab21Single, ab21Keys[1-ab21Index]); ab21Error != nil || ab21Found {
					ab21T.Fatal("opposite optional key must miss")
				}
			}
		})
	}
	ab21Duplicate, _ := ab21json.Marshal([]ab21json.RawMessage{ab21Row, ab21Row})
	for _, ab21Invalid := range [][]byte{[]byte("{"), []byte("null"), []byte("[]"), []byte("[null]"), []byte("[{}]"), []byte("{}"), append(append([]byte{}, ab21Duplicate...), []byte(" true")...), ab21Duplicate} {
		if _, ab21Error := ab21Validate(ab21Invalid); ab21Error == nil {
			ab21T.Fatal("invalid embedded data admitted")
		}
		if ab21Got, ab21Found, ab21Error := ab21Lookup(ab21Invalid, ab21BaseKey()); ab21Error == nil || ab21Found || ab21Got != nil {
			ab21T.Fatal("lookup must fail closed on invalid data")
		}
	}
}

func TestAb21IndependentBuffers(ab21T *ab21testing.T) {
	ab21Entries, ab21Error := Ab21Entries()
	if ab21Error != nil {
		ab21T.Fatal(ab21Error)
	}
	ab21Saved := append([]byte(nil), ab21Entries...)
	ab21Row, _ := ab21Get(ab21T, ab21BaseKey())
	ab21SavedRow := append([]byte(nil), ab21Row...)
	ab21OtherRow, _ := ab21Get(ab21T, ab21BaseKey())
	for ab21Index := range ab21Entries {
		ab21Entries[ab21Index] = 'x'
	}
	for ab21Index := range ab21Row {
		ab21Row[ab21Index] = 'y'
	}
	ab21Next, ab21Error := Ab21Entries()
	if ab21Error != nil || !ab21bytes.Equal(ab21Next, ab21Saved) {
		ab21T.Fatal("entries share storage")
	}
	ab21NextRow, _ := ab21Get(ab21T, ab21BaseKey())
	if !ab21bytes.Equal(ab21OtherRow, ab21SavedRow) || !ab21bytes.Equal(ab21NextRow, ab21SavedRow) {
		ab21T.Fatal("rows share storage")
	}
	ab21Catalog := append([]byte(nil), ab21Saved...)
	ab21SeamRow, ab21Found, ab21Error := ab21Lookup(ab21Catalog, ab21BaseKey())
	if ab21Error != nil || !ab21Found {
		ab21T.Fatal("seam lookup failed")
	}
	for ab21Index := range ab21Catalog {
		ab21Catalog[ab21Index] = 'z'
	}
	if !ab21bytes.Equal(ab21SeamRow, ab21SavedRow) {
		ab21T.Fatal("row aliases validation input")
	}
}
