package ebus_standard_catalog

import (
	ab81reflect "reflect"
	ab81testing "testing"
)

func ab81Expected() []Ab81Candidate {
	ab81False := false
	return []Ab81Candidate{
		{
			Ab81Key:           Ab81Key{Ab81Family: "testing", Ab81Variant: "partial", Ab81Primary: 0x0F, Ab81Direction: "request", Ab81Field: "declaration"},
			Ab81MetadataGrade: "Proven", Ab81Grade: "Hypothesis", Ab81Body: "Unknown",
			Ab81ACK: "Unknown", Ab81Support: "Unknown", Ab81Effects: "Unknown", Ab81Raw: []byte{0x0F},
			Ab81Positive: &Ab81Positive{Ab81Prefix: "PeBUS", Ab81TypedDataPresent: &ab81False, Ab81Grade: "Hypothesis"},
			Ab81Limits: []string{
				"Logical declaration only; no physical reply length is established.",
				"Absent secondary is distinct from 00 and standard services 0F01/0F02/0F03.",
				"The broader testing-family catalog remains incomplete; physical qualification and runtime adoption remain open before 1.0.",
			},
		},
		{
			Ab81Key:           Ab81Key{Ab81Family: "testing", Ab81Variant: "partial", Ab81Primary: 0x0F, Ab81Direction: "response", Ab81Field: "opaque_reply"},
			Ab81MetadataGrade: "Proven", Ab81Grade: "Unknown", Ab81Body: "Unknown",
			Ab81ACK: "Unknown", Ab81Support: "Unknown", Ab81Effects: "Unknown", Ab81Raw: []byte{0xAA},
			Ab81Limits: []string{"Synthetic opaque reply example only; AA remains raw and Unknown, without physical qualification."},
		},
	}
}

func TestAb81ProductionMetadata(ab81T *ab81testing.T) {
	ab81Want := ab81Expected()
	if ab81Got := Ab81Candidates(); !ab81reflect.DeepEqual(ab81Got, ab81Want) {
		ab81T.Fatalf("bounded production metadata differs: got %+v, want %+v", ab81Got, ab81Want)
	}
	for _, ab81Entry := range ab81Want {
		ab81Got, ab81OK := Ab81Lookup(ab81Entry.Ab81Key)
		if !ab81OK || !ab81reflect.DeepEqual(ab81Got, ab81Entry) {
			ab81T.Fatalf("exact production lookup differs: %+v, found %v", ab81Got, ab81OK)
		}
	}
}

func TestAb81ExactKeyAxes(ab81T *ab81testing.T) {
	ab81Base := ab81Expected()[0].Ab81Key
	ab81Cases := []struct {
		ab81Name   string
		ab81Change func(*Ab81Key)
	}{
		{"ab81_family", func(ab81Key *Ab81Key) { ab81Key.Ab81Family = "other" }},
		{"ab81_variant", func(ab81Key *Ab81Key) { ab81Key.Ab81Variant = "other" }},
		{"ab81_primary", func(ab81Key *Ab81Key) { ab81Key.Ab81Primary = 0x10 }},
		{"ab81_secondary_00", func(ab81Key *Ab81Key) { ab81Value := byte(0); ab81Key.Ab81Secondary = &ab81Value }},
		{"ab81_secondary_01", func(ab81Key *Ab81Key) { ab81Value := byte(1); ab81Key.Ab81Secondary = &ab81Value }},
		{"ab81_secondary_02", func(ab81Key *Ab81Key) { ab81Value := byte(2); ab81Key.Ab81Secondary = &ab81Value }},
		{"ab81_secondary_03", func(ab81Key *Ab81Key) { ab81Value := byte(3); ab81Key.Ab81Secondary = &ab81Value }},
		{"ab81_operation_00", func(ab81Key *Ab81Key) { ab81Value := byte(0); ab81Key.Ab81Operation = &ab81Value }},
		{"ab81_operation_01", func(ab81Key *Ab81Key) { ab81Value := byte(1); ab81Key.Ab81Operation = &ab81Value }},
		{"ab81_direction", func(ab81Key *Ab81Key) { ab81Key.Ab81Direction = "response" }},
		{"ab81_field", func(ab81Key *Ab81Key) { ab81Key.Ab81Field = "opaque_reply" }},
		{"ab81_unknown_family", func(ab81Key *Ab81Key) { ab81Key.Ab81Family = "" }},
		{"ab81_unknown_variant", func(ab81Key *Ab81Key) { ab81Key.Ab81Variant = "" }},
		{"ab81_unknown_direction", func(ab81Key *Ab81Key) { ab81Key.Ab81Direction = "" }},
		{"ab81_unknown_field", func(ab81Key *Ab81Key) { ab81Key.Ab81Field = "" }},
	}
	for _, ab81Case := range ab81Cases {
		ab81T.Run(ab81Case.ab81Name, func(ab81T *ab81testing.T) {
			if _, ab81OK := Ab81Lookup(ab81Base); !ab81OK {
				ab81T.Fatal("valid control missing")
			}
			ab81Changed := ab81Base
			ab81Case.ab81Change(&ab81Changed)
			ab81Got, ab81OK := Ab81Lookup(ab81Changed)
			if ab81OK || !ab81reflect.DeepEqual(ab81Got, Ab81Candidate{}) {
				ab81T.Fatalf("changed axis matched: %+v", ab81Changed)
			}
		})
	}
}

func TestAb81FreshNestedData(ab81T *ab81testing.T) {
	ab81Want := ab81Expected()
	ab81Entries := Ab81Candidates()
	ab81Entries[0].Ab81Raw[0] = 0
	ab81Entries[0].Ab81Limits[0] = "changed"
	*ab81Entries[0].Ab81Positive.Ab81TypedDataPresent = true
	ab81Entries[0].Ab81Positive.Ab81Prefix = "changed"
	ab81Entries[0].Ab81Key.Ab81Family = "changed"
	ab81Entries[1].Ab81Raw[0] = 0
	ab81Entries[1].Ab81Limits[0] = "changed"
	ab81Lookup, ab81OK := Ab81Lookup(ab81Want[0].Ab81Key)
	if !ab81OK {
		ab81T.Fatal("lookup control missing")
	}
	ab81Lookup.Ab81Raw[0] = 0
	ab81Lookup.Ab81Limits[0] = "changed"
	*ab81Lookup.Ab81Positive.Ab81TypedDataPresent = true
	ab81Lookup.Ab81Positive.Ab81Grade = "changed"
	if ab81Got := Ab81Candidates(); !ab81reflect.DeepEqual(ab81Got, ab81Want) {
		ab81T.Fatal("accessor shares nested storage")
	}
	for _, ab81Entry := range ab81Want {
		ab81Got, ab81OK := Ab81Lookup(ab81Entry.Ab81Key)
		if !ab81OK || !ab81reflect.DeepEqual(ab81Got, ab81Entry) {
			ab81T.Fatal("lookup shares nested storage")
		}
	}
}
