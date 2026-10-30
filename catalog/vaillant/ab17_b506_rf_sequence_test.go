package vaillant

import (
	ab17Bytes "bytes"
	ab17JSON "encoding/json"
	ab17Fmt "fmt"
	ab17Reflect "reflect"
	ab17Testing "testing"
)

func ab17Decode(ab17T *ab17Testing.T, ab17Raw []byte) map[string]ab17JSON.RawMessage {
	ab17T.Helper()
	var ab17Object map[string]ab17JSON.RawMessage
	if ab17Error := ab17JSON.Unmarshal(ab17Raw, &ab17Object); ab17Error != nil {
		ab17T.Fatal(ab17Error)
	}
	return ab17Object
}

func ab17Rows(ab17T *ab17Testing.T) []ab17JSON.RawMessage {
	ab17T.Helper()
	ab17Data, ab17Error := Ab17Entries()
	if ab17Error != nil {
		ab17T.Fatal(ab17Error)
	}
	var ab17Result []ab17JSON.RawMessage
	if ab17Error := ab17JSON.Unmarshal(ab17Data, &ab17Result); ab17Error != nil {
		ab17T.Fatal(ab17Error)
	}
	return ab17Result
}

func ab17KeyOf(ab17T *ab17Testing.T, ab17Raw []byte) Ab17Key {
	ab17T.Helper()
	var ab17Entry ab17Row
	if ab17Error := ab17JSON.Unmarshal(ab17Raw, &ab17Entry); ab17Error != nil {
		ab17T.Fatal(ab17Error)
	}
	return ab17Entry.Ab17Key
}

func ab17Want(ab17T *ab17Testing.T, ab17Object map[string]ab17JSON.RawMessage, ab17Field, ab17Expected string) {
	ab17T.Helper()
	if string(ab17Object[ab17Field]) != ab17Expected {
		ab17T.Fatalf("%s = %s; want %s", ab17Field, ab17Object[ab17Field], ab17Expected)
	}
}

func TestAb17FiniteMetadata(ab17T *ab17Testing.T) {
	ab17Entries := ab17Rows(ab17T)
	if len(ab17Entries) != 6 {
		ab17T.Fatalf("rows: %d", len(ab17Entries))
	}
	for _, ab17Raw := range ab17Entries {
		ab17Key := ab17KeyOf(ab17T, ab17Raw)
		ab17Lookup, ab17Found, ab17Error := Ab17Lookup(ab17Key)
		if ab17Error != nil || !ab17Found || !ab17Bytes.Equal(ab17Lookup, ab17Raw) {
			ab17T.Fatalf("whole row lookup: %v %v", ab17Found, ab17Error)
		}
		ab17Object := ab17Decode(ab17T, ab17Lookup)
		ab17Want(ab17T, ab17Object, "grade", `"Hypothesis"`)
		ab17Want(ab17T, ab17Object, "qualification", `"Unknown"`)
		ab17Want(ab17T, ab17Object, "ab17_raw_form", `"retain-original-bytes"`)
		for _, ab17Field := range []string{"ab17_extent_bytes", "ab17_minimum_bytes", "ab17_maximum_bytes", "ab17_field_position"} {
			ab17Want(ab17T, ab17Object, ab17Field, "null")
		}
		ab17Unknown := ab17Decode(ab17T, ab17Object["ab17_unknown"])
		for _, ab17Field := range []string{"native_support", "physical_effects", "persistence", "delivery"} {
			ab17Want(ab17T, ab17Unknown, ab17Field, `"Unknown"`)
		}
		if ab17Key.Ab17Family != "vaillant" || ab17Key.Ab17Variant != "B506-RF" || ab17Key.Ab17Primary != 181 || ab17Key.Ab17Secondary == nil || *ab17Key.Ab17Secondary != 6 || ab17Key.Ab17Direction != "none" || ab17Key.Ab17Group != nil || ab17Key.Ab17Register != nil || ab17Key.Ab17InstanceScope != "none" {
			ab17T.Fatal("unexpected finite identity")
		}
		if ab17Key.Ab17Operation == nil {
			ab17T.Fatal("missing operation")
		}
		if ab17Key.Ab17Field != "stop" && *ab17Key.Ab17Operation != 13 {
			ab17T.Fatal("RF metadata operation")
		}
	}
	ab17Count := ab17Decode(ab17T, ab17Decode(ab17T, ab17Entries[0])["ab17_count"])
	ab17Want(ab17T, ab17Count, "ab17_zero_meaning", `"unlimited"`)
	ab17Want(ab17T, ab17Count, "ab17_finite_minimum", "1")
	ab17Want(ab17T, ab17Count, "ab17_finite_maximum", "255")
	ab17Want(ab17T, ab17Count, "ab17_finite_meaning", `"count-value"`)
	for _, ab17Case := range []struct {
		ab17Count    int
		ab17Expected string
	}{{0, "unlimited"}, {1, "finite-1"}, {2, "finite-2"}, {255, "finite-255"}} {
		ab17T.Run(ab17Fmt.Sprintf("ab17-count-%d", ab17Case.ab17Count), func(ab17T *ab17Testing.T) {
			var ab17Minimum, ab17Maximum int
			var ab17Zero, ab17Meaning string
			for _, ab17DecodeCase := range []struct {
				ab17Name   string
				ab17Target any
			}{{"ab17_finite_minimum", &ab17Minimum}, {"ab17_finite_maximum", &ab17Maximum}, {"ab17_zero_meaning", &ab17Zero}, {"ab17_finite_meaning", &ab17Meaning}} {
				if ab17Error := ab17JSON.Unmarshal(ab17Count[ab17DecodeCase.ab17Name], ab17DecodeCase.ab17Target); ab17Error != nil {
					ab17T.Fatal(ab17Error)
				}
			}
			ab17Actual := ab17Zero
			if ab17Case.ab17Count != 0 {
				if ab17Case.ab17Count < ab17Minimum || ab17Case.ab17Count > ab17Maximum || ab17Meaning != "count-value" {
					ab17T.Fatal("finite count unmapped")
				}
				ab17Actual = ab17Fmt.Sprintf("finite-%d", ab17Case.ab17Count)
			}
			if ab17Actual != ab17Case.ab17Expected {
				ab17T.Fatalf("count meaning %q", ab17Actual)
			}
		})
	}
	for ab17Sequence, ab17Expected := range []string{`"ab17-class-a"`, `"ab17-class-a"`, `"ab17-class-a"`, `"ab17-class-b"`} {
		ab17Key := ab17KeyOf(ab17T, ab17Entries[1])
		ab17Key.Ab17Alternative = ab17Fmt.Sprint(ab17Sequence)
		ab17Raw, ab17Found, ab17Error := Ab17Lookup(ab17Key)
		if ab17Error != nil || !ab17Found {
			ab17T.Fatalf("sequence %d: %v", ab17Sequence, ab17Error)
		}
		ab17Object := ab17Decode(ab17T, ab17Raw)
		ab17Want(ab17T, ab17Object, "ab17_extent_class", ab17Expected)
		for _, ab17Payload := range [][]byte{{}, {0x23, 0x45}} {
			ab17T.Run(ab17Fmt.Sprintf("ab17-sequence-%d-payload-length-%d", ab17Sequence, len(ab17Payload)), func(ab17T *ab17Testing.T) {
				ab17Want(ab17T, ab17Object, "qualification", `"Unknown"`)
				ab17Want(ab17T, ab17Object, "ab17_extent_bytes", "null")
			})
		}
	}
	for ab17Sequence := 4; ab17Sequence <= 255; ab17Sequence++ {
		ab17Key := ab17KeyOf(ab17T, ab17Entries[1])
		ab17Key.Ab17Alternative = ab17Fmt.Sprint(ab17Sequence)
		if ab17Raw, ab17Found, ab17Error := Ab17Lookup(ab17Key); ab17Error != nil || ab17Found || ab17Raw != nil {
			ab17T.Fatalf("unmapped sequence %d", ab17Sequence)
		}
	}
	ab17Stop := ab17Decode(ab17T, ab17Entries[5])
	ab17StopKey := ab17KeyOf(ab17T, ab17Entries[5])
	if ab17StopKey.Ab17Field != "stop" || *ab17StopKey.Ab17Operation != 14 {
		ab17T.Fatal("stop identity")
	}
	var ab17Related []int
	if ab17Error := ab17JSON.Unmarshal(ab17Stop["ab17_related_operations"], &ab17Related); ab17Error != nil {
		ab17T.Fatal(ab17Error)
	}
	if !ab17Reflect.DeepEqual(ab17Related, []int{12, 13}) {
		ab17T.Fatal("stop associations")
	}
	for _, ab17Field := range []string{"ab17_activation_order", "ab17_delivery_guarantee", "ab17_completed_stop"} {
		ab17Want(ab17T, ab17Stop, ab17Field, `"Unknown"`)
	}
	for _, ab17Field := range []string{"ab17_retry_count", "ab17_timeout"} {
		ab17Want(ab17T, ab17Stop, ab17Field, "null")
	}
}

func TestAb17IndependentAxes(ab17T *ab17Testing.T) {
	ab17Base := ab17Rows(ab17T)[0]
	ab17BaseKey := ab17KeyOf(ab17T, ab17Base)
	ab17Zero := uint8(0)
	ab17Other := uint8(7)
	ab17Register := uint16(0)
	for _, ab17Case := range []struct {
		ab17Name   string
		ab17Change func(*Ab17Key)
	}{
		{"ab17-family", func(ab17Key *Ab17Key) { ab17Key.Ab17Family = "ab17-other" }},
		{"ab17-variant", func(ab17Key *Ab17Key) { ab17Key.Ab17Variant = "ab17-other" }},
		{"ab17-primary", func(ab17Key *Ab17Key) { ab17Key.Ab17Primary = 182 }},
		{"ab17-secondary", func(ab17Key *Ab17Key) { ab17Key.Ab17Secondary = &ab17Other }},
		{"ab17-operation", func(ab17Key *Ab17Key) { ab17Key.Ab17Operation = &ab17Other }},
		{"ab17-direction", func(ab17Key *Ab17Key) { ab17Key.Ab17Direction = "ab17-other" }},
		{"ab17-field", func(ab17Key *Ab17Key) { ab17Key.Ab17Field = "ab17-other" }},
		{"ab17-alternative", func(ab17Key *Ab17Key) { ab17Key.Ab17Alternative = "ab17-other" }},
		{"ab17-group", func(ab17Key *Ab17Key) { ab17Key.Ab17Group = &ab17Other }},
		{"ab17-instance-scope", func(ab17Key *Ab17Key) { ab17Key.Ab17InstanceScope = "dynamic" }},
		{"ab17-register", func(ab17Key *Ab17Key) { ab17OtherRegister := uint16(1); ab17Key.Ab17Register = &ab17OtherRegister }},
		{"ab17-secondary-absent-zero", func(ab17Key *Ab17Key) { ab17Key.Ab17Secondary = nil }},
		{"ab17-operation-absent-zero", func(ab17Key *Ab17Key) { ab17Key.Ab17Operation = nil }},
		{"ab17-group-absent-zero", func(ab17Key *Ab17Key) { ab17Key.Ab17Group = &ab17Zero }},
		{"ab17-register-absent-zero", func(ab17Key *Ab17Key) { ab17Key.Ab17Register = &ab17Register }},
	} {
		ab17T.Run(ab17Case.ab17Name, func(ab17T *ab17Testing.T) {
			ab17Control := ab17BaseKey
			if ab17Case.ab17Name == "ab17-secondary-absent-zero" {
				ab17Control.Ab17Secondary = &ab17Zero
			}
			if ab17Case.ab17Name == "ab17-operation-absent-zero" {
				ab17Control.Ab17Operation = &ab17Zero
			}
			ab17Changed := ab17Control
			ab17Case.ab17Change(&ab17Changed)
			if ab17Raw, ab17Found, ab17Error := Ab17Lookup(ab17Changed); ab17Error != nil || ab17Found || ab17Raw != nil {
				ab17T.Fatal("public lookup fell back")
			}
			ab17Objects := []map[string]ab17JSON.RawMessage{ab17Decode(ab17T, ab17Base), ab17Decode(ab17T, ab17Base)}
			for ab17Index, ab17Key := range []Ab17Key{ab17Control, ab17Changed} {
				ab17KeyJSON, ab17Error := ab17JSON.Marshal(ab17Key)
				if ab17Error != nil {
					ab17T.Fatal(ab17Error)
				}
				ab17Objects[ab17Index]["key"] = ab17KeyJSON
				ab17Objects[ab17Index]["ab17_id"] = ab17JSON.RawMessage(ab17Fmt.Sprintf(`"ab17-control-%d"`, ab17Index))
			}
			ab17Data, ab17Error := ab17JSON.Marshal(ab17Objects)
			if ab17Error != nil {
				ab17T.Fatal(ab17Error)
			}
			if _, ab17Error := ab17EntriesFrom(string(ab17Data)); ab17Error != nil {
				ab17T.Fatalf("distinct identity rejected: %v", ab17Error)
			}
			for ab17Index, ab17Key := range []Ab17Key{ab17Control, ab17Changed} {
				ab17Raw, ab17Found, ab17Error := ab17LookupFrom(string(ab17Data), ab17Key)
				if ab17Error != nil || !ab17Found {
					ab17T.Fatalf("valid control: %v", ab17Error)
				}
				ab17Want(ab17T, ab17Decode(ab17T, ab17Raw), "ab17_id", ab17Fmt.Sprintf(`"ab17-control-%d"`, ab17Index))
			}
		})
	}
}

func TestAb17Validation(ab17T *ab17Testing.T) {
	ab17Base := ab17Rows(ab17T)[0]
	ab17Key := ab17KeyOf(ab17T, ab17Base)
	for ab17Index, ab17Data := range []string{"{", "null", "{}", "[null]", "[{}]", string(ab17Base), "[] false", "[" + string(ab17Base) + "," + string(ab17Base) + "]"} {
		ab17T.Run(ab17Fmt.Sprintf("ab17-invalid-%d", ab17Index), func(ab17T *ab17Testing.T) {
			if _, ab17Error := ab17EntriesFrom(ab17Data); ab17Error == nil {
				ab17T.Fatal("invalid catalog admitted")
			}
			if ab17Raw, ab17Found, ab17Error := ab17LookupFrom(ab17Data, ab17Key); ab17Error == nil || ab17Found || ab17Raw != nil {
				ab17T.Fatal("lookup did not fail closed")
			}
		})
	}
	for _, ab17Axis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
		ab17T.Run("ab17-missing-"+ab17Axis, func(ab17T *ab17Testing.T) {
			ab17Object := ab17Decode(ab17T, ab17Base)
			ab17Axes := ab17Decode(ab17T, ab17Object["key"])
			delete(ab17Axes, ab17Axis)
			ab17KeyJSON, ab17Error := ab17JSON.Marshal(ab17Axes)
			if ab17Error != nil {
				ab17T.Fatal(ab17Error)
			}
			ab17Object["key"] = ab17KeyJSON
			ab17Data, ab17Error := ab17JSON.Marshal([]any{ab17Object})
			if ab17Error != nil {
				ab17T.Fatal(ab17Error)
			}
			if _, ab17Error := ab17EntriesFrom(string(ab17Data)); ab17Error == nil {
				ab17T.Fatal("missing identity admitted")
			}
		})
	}
}

func TestAb17IndependentBuffers(ab17T *ab17Testing.T) {
	ab17Baseline := ab17Rows(ab17T)[0]
	ab17Key := ab17KeyOf(ab17T, ab17Baseline)
	ab17First, ab17Error := Ab17Entries()
	if ab17Error != nil {
		ab17T.Fatal(ab17Error)
	}
	ab17Second, ab17Error := Ab17Entries()
	if ab17Error != nil {
		ab17T.Fatal(ab17Error)
	}
	ab17Snapshot := append([]byte(nil), ab17Second...)
	for ab17Index := range ab17First {
		ab17First[ab17Index] = 'x'
	}
	if !ab17Bytes.Equal(ab17Second, ab17Snapshot) {
		ab17T.Fatal("entries buffers share memory")
	}
	ab17Third, ab17Error := Ab17Entries()
	if ab17Error != nil || !ab17Bytes.Equal(ab17Third, ab17Snapshot) {
		ab17T.Fatal("entries mutation persisted")
	}
	ab17FirstRow, ab17Found, ab17Error := Ab17Lookup(ab17Key)
	if ab17Error != nil || !ab17Found {
		ab17T.Fatal("lookup")
	}
	ab17SecondRow, ab17Found, ab17Error := Ab17Lookup(ab17Key)
	if ab17Error != nil || !ab17Found {
		ab17T.Fatal("lookup")
	}
	for ab17Index := range ab17FirstRow {
		ab17FirstRow[ab17Index] = 'x'
	}
	if !ab17Bytes.Equal(ab17SecondRow, ab17Baseline) {
		ab17T.Fatal("row buffers share memory")
	}
	ab17ThirdRow, ab17Found, ab17Error := Ab17Lookup(ab17Key)
	if ab17Error != nil || !ab17Found || !ab17Bytes.Equal(ab17ThirdRow, ab17Baseline) {
		ab17T.Fatal("row mutation persisted")
	}
	ab17Fourth, ab17Error := Ab17Entries()
	if ab17Error != nil || !ab17Bytes.Equal(ab17Fourth, ab17Snapshot) {
		ab17T.Fatal("row mutation affected entries")
	}
}
