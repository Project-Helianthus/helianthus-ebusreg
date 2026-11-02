package ebus_standard_catalog

import (
	ab54bytes "bytes"
	ab54json "encoding/json"
	ab54fmt "fmt"
	ab54reflect "reflect"
	ab54strings "strings"
	ab54testing "testing"
)

func ab54Byte(ab54Value uint8) *uint8   { return &ab54Value }
func ab54Word(ab54Value uint16) *uint16 { return &ab54Value }

func ab54TestKey(ab54Secondary uint8) Ab54Key {
	ab54Direction := "none"
	if ab54Secondary >= 3 {
		ab54Direction = "abstract-request"
	}
	return Ab54Key{Ab54Family: "standard-service", Ab54Variant: ab54fmt.Sprintf("FF%02d", ab54Secondary), Ab54Primary: 255,
		Ab54Secondary: ab54Byte(ab54Secondary), Ab54Direction: ab54Direction, Ab54Field: "body-completeness",
		Ab54Alternative: "none", Ab54InstanceScope: "none"}
}

func ab54Object(ab54T *ab54testing.T, ab54Data []byte) map[string]ab54json.RawMessage {
	ab54T.Helper()
	var ab54Result map[string]ab54json.RawMessage
	if ab54Err := ab54json.Unmarshal(ab54Data, &ab54Result); ab54Err != nil {
		ab54T.Fatal(ab54Err)
	}
	return ab54Result
}

func ab54Expect(ab54T *ab54testing.T, ab54Row map[string]ab54json.RawMessage, ab54Field string, ab54Want string) {
	ab54T.Helper()
	if !ab54bytes.Equal(ab54bytes.TrimSpace(ab54Row[ab54Field]), []byte(ab54Want)) {
		ab54T.Fatalf("ab54 %s = %s, want %s", ab54Field, ab54Row[ab54Field], ab54Want)
	}
}

func ab54Get(ab54T *ab54testing.T, ab54Key Ab54Key) ab54json.RawMessage {
	ab54T.Helper()
	ab54Row, ab54Found, ab54Err := Ab54Lookup(ab54Key)
	if ab54Err != nil || !ab54Found {
		ab54T.Fatalf("ab54 missing row: found=%v error=%v", ab54Found, ab54Err)
	}
	return ab54Row
}

func TestAb54Inventory(ab54T *ab54testing.T) {
	ab54Data, ab54Err := Ab54Entries()
	if ab54Err != nil {
		ab54T.Fatal(ab54Err)
	}
	var ab54Rows []ab54json.RawMessage
	if ab54Err := ab54json.Unmarshal(ab54Data, &ab54Rows); ab54Err != nil {
		ab54T.Fatal(ab54Err)
	}
	if len(ab54Rows) != 7 {
		ab54T.Fatalf("ab54 identities = %d, want 7", len(ab54Rows))
	}
	ab54Queries := []string{"none", "none", "none", "net-status-query", "monitored-participants-query", "failed-nodes-query", "required-services-query"}
	ab54NN := []string{"Unknown", "0", "0", "Unknown", "1..10", "Unknown", "1+2n for complete service pairs"}
	for ab54SB := uint8(0); ab54SB < 7; ab54SB++ {
		ab54Raw := ab54Get(ab54T, ab54TestKey(ab54SB))
		if !ab54bytes.Equal(ab54Raw, ab54Rows[ab54SB]) {
			ab54T.Fatal("ab54 lookup did not return whole exact row")
		}
		ab54Row := ab54Object(ab54T, ab54Raw)
		ab54Expect(ab54T, ab54Row, "id", ab54fmt.Sprintf("\"ab54-ff%02d\"", ab54SB))
		ab54Expect(ab54T, ab54Row, "grade", `"Proven"`)
		ab54Expect(ab54T, ab54Row, "invocable", "false")
		ab54Expect(ab54T, ab54Row, "executable", "false")
		ab54Expect(ab54T, ab54Row, "semantic_validity", `"Unknown"`)
		ab54Expect(ab54T, ab54Row, "physical_qualification", `"Unknown"`)
		ab54Expect(ab54T, ab54Row, "query_identity", `"`+ab54Queries[ab54SB]+`"`)
		ab54Reply := `"Unknown"`
		if ab54SB >= 3 {
			ab54Reply = `"same-service"`
		}
		ab54Expect(ab54T, ab54Row, "reply_service", ab54Reply)
		ab54Form := ab54Object(ab54T, ab54Row["raw_form"])
		ab54Expect(ab54T, ab54Form, "PB", `"FF"`)
		ab54Expect(ab54T, ab54Form, "SB", ab54fmt.Sprintf("\"%02d\"", ab54SB))
		ab54Expect(ab54T, ab54Form, "NN", `"`+ab54NN[ab54SB]+`"`)
		ab54Expect(ab54T, ab54Form, "body", `"Unknown"`)
		ab54Unknowns := ab54Object(ab54T, ab54Row["unknowns"])
		for _, ab54Axis := range []string{"native_emission", "forwarding_eligibility", "target_identity", "correlation", "lifecycle_effects", "installed_support", "physical_ack", "exact_body_fields", "pair_meaning", "presence", "freshness"} {
			ab54Expect(ab54T, ab54Unknowns, ab54Axis, `"Unknown"`)
		}
		if ab54SB != 4 {
			ab54Expect(ab54T, ab54Row, "minimum_overhead_width", "null")
		}
		if ab54SB != 6 {
			ab54Expect(ab54T, ab54Row, "service_pair_width", "null")
		}
		if ab54SB != 4 && ab54SB != 6 {
			ab54Expect(ab54T, ab54Row, "cases", "[]")
		}
	}
	ab54Reset := ab54Object(ab54T, ab54Get(ab54T, ab54TestKey(1)))
	ab54Expect(ab54T, ab54Reset, "policy_boundary", `"NN0; reset of dynamic target configuration requires explicit emission policy; no emitter added"`)
	ab54Availability := ab54Object(ab54T, ab54Get(ab54T, ab54TestKey(2)))
	ab54Expect(ab54T, ab54Availability, "policy_boundary", `"NN0 independent of responder availability; no responder inference"`)
}

func TestAb54BodyCompleteness(ab54T *ab54testing.T) {
	for _, ab54SB := range []uint8{4, 6} {
		ab54Row := ab54Object(ab54T, ab54Get(ab54T, ab54TestKey(ab54SB)))
		var ab54Cases []map[string]ab54json.RawMessage
		if ab54Err := ab54json.Unmarshal(ab54Row["cases"], &ab54Cases); ab54Err != nil {
			ab54T.Fatal(ab54Err)
		}
		ab54Count := 2
		if ab54SB == 6 {
			ab54Count = 5
		}
		if len(ab54Cases) != ab54Count {
			ab54T.Fatal("ab54 case inventory changed")
		}
		for ab54Index, ab54Case := range ab54Cases {
			ab54NN := ab54Index + 1
			ab54Expect(ab54T, ab54Case, "id", ab54fmt.Sprintf("\"ab54-ff%02d-length%d\"", ab54SB, ab54NN))
			ab54Expect(ab54T, ab54Case, "nn", ab54fmt.Sprint(ab54NN))
			ab54Expect(ab54T, ab54Case, "grade", `"Proven"`)
			ab54Expect(ab54T, ab54Case, "semantic_validity", `"Unknown"`)
			var ab54Raw []int
			if ab54Err := ab54json.Unmarshal(ab54Case["raw"], &ab54Raw); ab54Err != nil {
				ab54T.Fatal(ab54Err)
			}
			if !ab54reflect.DeepEqual(ab54Raw, make([]int, ab54NN)) {
				ab54T.Fatal("ab54 synthetic raw body changed")
			}
			if ab54SB == 4 {
				ab54Expect(ab54T, ab54Row, "minimum_overhead_width", "2")
				ab54Expect(ab54T, ab54Case, "minimum_overhead_met", []string{"false", "true"}[ab54Index])
				ab54Expect(ab54T, ab54Case, "body_complete", []string{"false", "null"}[ab54Index])
				ab54Expect(ab54T, ab54Case, "state", []string{`"truncated/raw"`, `"minimum-overhead-only"`}[ab54Index])
			} else {
				ab54Expect(ab54T, ab54Row, "service_pair_width", "2")
				ab54Expect(ab54T, ab54Case, "complete_pairs", []string{"0", "0", "1", "1", "2"}[ab54Index])
				ab54Expect(ab54T, ab54Case, "pair_shape_complete", []string{"true", "false", "true", "false", "true"}[ab54Index])
				ab54Expect(ab54T, ab54Case, "body_complete", []string{"null", "false", "null", "false", "null"}[ab54Index])
				ab54State, ab54Remainder := `"arithmetic-only"`, []int{}
				if ab54NN == 2 || ab54NN == 4 {
					ab54State, ab54Remainder = `"incomplete/raw-remainder"`, []int{0}
				}
				ab54Expect(ab54T, ab54Case, "state", ab54State)
				var ab54GotRemainder []int
				if ab54Err := ab54json.Unmarshal(ab54Case["raw_remainder"], &ab54GotRemainder); ab54Err != nil {
					ab54T.Fatal(ab54Err)
				}
				if !ab54reflect.DeepEqual(ab54GotRemainder, ab54Remainder) {
					ab54T.Fatal("ab54 lost trailing raw byte")
				}
			}
		}
	}
}

func TestAb54ExactAxes(ab54T *ab54testing.T) {
	ab54Control := ab54TestKey(4)
	ab54Raw := ab54Get(ab54T, ab54Control)
	ab54Mutations := []struct {
		Ab54Name   string
		Ab54Change func(*Ab54Key)
	}{
		{"ab54-family", func(ab54K *Ab54Key) { ab54K.Ab54Family = "other-family" }},
		{"ab54-variant", func(ab54K *Ab54Key) { ab54K.Ab54Variant = "FF03" }},
		{"ab54-primary", func(ab54K *Ab54Key) { ab54K.Ab54Primary = 254 }},
		{"ab54-primary-zero", func(ab54K *Ab54Key) { ab54K.Ab54Primary = 0 }},
		{"ab54-secondary", func(ab54K *Ab54Key) { ab54K.Ab54Secondary = ab54Byte(3) }},
		{"ab54-secondary-absent", func(ab54K *Ab54Key) { ab54K.Ab54Secondary = nil }},
		{"ab54-secondary-zero", func(ab54K *Ab54Key) { ab54K.Ab54Secondary = ab54Byte(0) }},
		{"ab54-operation-zero", func(ab54K *Ab54Key) { ab54K.Ab54Operation = ab54Byte(0) }},
		{"ab54-operation-value", func(ab54K *Ab54Key) { ab54K.Ab54Operation = ab54Byte(1) }},
		{"ab54-direction", func(ab54K *Ab54Key) { ab54K.Ab54Direction = "abstract-positive" }},
		{"ab54-field", func(ab54K *Ab54Key) { ab54K.Ab54Field = "other-field" }},
		{"ab54-alternative", func(ab54K *Ab54Key) { ab54K.Ab54Alternative = "other-alternative" }},
		{"ab54-group-zero", func(ab54K *Ab54Key) { ab54K.Ab54Group = ab54Word(0) }},
		{"ab54-group-value", func(ab54K *Ab54Key) { ab54K.Ab54Group = ab54Word(1) }},
		{"ab54-instance-scope", func(ab54K *Ab54Key) { ab54K.Ab54InstanceScope = "other-scope" }},
		{"ab54-register-zero", func(ab54K *Ab54Key) { ab54K.Ab54Register = ab54Word(0) }},
		{"ab54-register-value", func(ab54K *Ab54Key) { ab54K.Ab54Register = ab54Word(1) }},
	}
	for _, ab54Mutation := range ab54Mutations {
		ab54T.Run(ab54Mutation.Ab54Name, func(ab54T *ab54testing.T) {
			ab54Key := ab54Control
			ab54Mutation.Ab54Change(&ab54Key)
			if ab54Got, ab54Found, ab54Err := Ab54Lookup(ab54Key); ab54Err != nil || ab54Found || ab54Got != nil {
				ab54T.Fatal("ab54 changed axis selected production row")
			}
			ab54Other := ab54Object(ab54T, ab54Raw)
			ab54Other["key"], _ = ab54json.Marshal(ab54Key)
			ab54Other["id"] = ab54json.RawMessage(`"ab54-synthetic-axis-control"`)
			ab54OtherRaw, _ := ab54json.Marshal(ab54Other)
			ab54Source := append(append(append([]byte{'['}, ab54Raw...), ','), ab54OtherRaw...)
			ab54Source = append(ab54Source, ']')
			for _, ab54Expected := range []struct {
				Ab54Key Ab54Key
				Ab54Raw []byte
			}{{ab54Control, ab54Raw}, {ab54Key, ab54OtherRaw}} {
				ab54Got, ab54Found, ab54Err := ab54Lookup(ab54Source, ab54Expected.Ab54Key)
				if ab54Err != nil || !ab54Found || !ab54bytes.Equal(ab54Got, ab54Expected.Ab54Raw) {
					ab54T.Fatalf("ab54 valid independent control: %v", ab54Err)
				}
			}
		})
	}
	ab54Zero := ab54TestKey(0)
	ab54Get(ab54T, ab54Zero)
	ab54Zero.Ab54Secondary = nil
	if _, ab54Found, ab54Err := Ab54Lookup(ab54Zero); ab54Err != nil || ab54Found {
		ab54T.Fatal("ab54 absent SB matched zero")
	}
}

func TestAb54InvalidAndDuplicate(ab54T *ab54testing.T) {
	ab54Raw := ab54Get(ab54T, ab54TestKey(4))
	ab54Duplicate := "[" + string(ab54Raw) + "," + string(ab54Raw) + "]"
	for ab54Index, ab54Source := range []string{"", "null", "{}", "[", "[null]", "[{}]", "[1]", "[] false", ab54Duplicate,
		"[" + ab54strings.Replace(string(ab54Raw), `"primary": 255`, `"primary": 255, "primary": 254`, 1) + "]",
		"[" + ab54strings.Replace(string(ab54Raw), `"primary": 255`, `"primary": 255, "Primary": 254`, 1) + "]",
		"[" + ab54strings.Replace(string(ab54Raw), `"key": {`, `"key": {}, "key": {`, 1) + "]",
		"[" + ab54strings.Replace(string(ab54Raw), `"primary": 255,`, ``, 1) + "]",
		"[" + ab54strings.Replace(string(ab54Raw), `"primary": 255`, `"primary": null`, 1) + "]",
		"[" + ab54strings.Replace(string(ab54Raw), `"alternative": "none"`, `"alternative": ""`, 1) + "]",
		"[" + ab54strings.Replace(string(ab54Raw), `"primary": 255`, `"primary": 256`, 1) + "]",
	} {
		ab54T.Run(ab54fmt.Sprintf("ab54-invalid-%d", ab54Index), func(ab54T *ab54testing.T) {
			if ab54Entries, ab54Index, ab54Err := ab54Validate([]byte(ab54Source)); ab54Err == nil || ab54Entries != nil || ab54Index != nil {
				ab54T.Fatal("ab54 invalid source succeeded")
			}
			if ab54Got, ab54Found, ab54Err := ab54Lookup([]byte(ab54Source), ab54TestKey(4)); ab54Err == nil || ab54Found || ab54Got != nil {
				ab54T.Fatal("ab54 invalid lookup succeeded")
			}
		})
	}
	ab54AbsentKey := ab54TestKey(4)
	ab54AbsentKey.Ab54Family = "absent"
	if _, ab54Found, ab54Err := ab54Lookup([]byte(ab54Duplicate), ab54AbsentKey); ab54Err == nil || ab54Found {
		ab54T.Fatal("ab54 invalid catalog was hidden by a miss")
	}
}

func TestAb54IndependentStorage(ab54T *ab54testing.T) {
	ab54Before := ab54Get(ab54T, ab54TestKey(4))
	ab54Entries, ab54Err := Ab54Entries()
	if ab54Err != nil {
		ab54T.Fatal(ab54Err)
	}
	ab54Expected := append([]byte(nil), ab54Entries...)
	ab54Entries[0] = '!'
	ab54Returned := ab54Get(ab54T, ab54TestKey(4))
	ab54Returned[0] = '!'
	ab54Again, ab54Err := Ab54Entries()
	if ab54Err != nil || !ab54bytes.Equal(ab54Again, ab54Expected) || !ab54bytes.Equal(ab54Before, ab54Get(ab54T, ab54TestKey(4))) {
		ab54T.Fatal("ab54 returned mutation changed later API results")
	}
	ab54Source := append([]byte(nil), ab54Expected...)
	ab54Detached, ab54Index, ab54Err := ab54Validate(ab54Source)
	if ab54Err != nil {
		ab54T.Fatal(ab54Err)
	}
	for ab54Offset := range ab54Source {
		ab54Source[ab54Offset] = '!'
	}
	if !ab54bytes.Equal(ab54Detached, ab54Expected) {
		ab54T.Fatal("ab54 source aliases entries")
	}
	ab54Identity, _ := ab54json.Marshal(ab54TestKey(4))
	if !ab54bytes.Equal(ab54Index[string(ab54Identity)], ab54Before) {
		ab54T.Fatal("ab54 source aliases row")
	}
	ab54Index[string(ab54Identity)][0] = '!'
	if !ab54bytes.Equal(ab54Detached, ab54Expected) {
		ab54T.Fatal("ab54 entries alias indexed row")
	}
}

func TestAb54MetadataPreservation(ab54T *ab54testing.T) {
	ab54Raw := ab54Get(ab54T, ab54TestKey(4))
	for _, ab54Grade := range []string{"Proven", "Hypothesis", "Unknown"} {
		for _, ab54Width := range []string{"null", "0", "1", "3"} {
			ab54Row := ab54Object(ab54T, ab54Raw)
			ab54Row["grade"] = ab54json.RawMessage(`"` + ab54Grade + `"`)
			ab54Row["minimum_overhead_width"] = ab54json.RawMessage(ab54Width)
			ab54Row["ab54_unresolved"] = ab54json.RawMessage(`{"ab54_raw":[0,255],"ab54_grade":"Unknown"}`)
			ab54Changed, _ := ab54json.Marshal(ab54Row)
			ab54Source := append(append([]byte{'['}, ab54Changed...), ']')
			ab54Got, ab54Found, ab54Err := ab54Lookup(ab54Source, ab54TestKey(4))
			if ab54Err != nil || !ab54Found || !ab54bytes.Equal(ab54Got, ab54Changed) {
				ab54T.Fatal("ab54 width, grade or unresolved raw metadata lost")
			}
		}
	}
}

func TestAb54OptionalValues(ab54T *ab54testing.T) {
	for _, ab54Axis := range []string{"ab54-secondary", "ab54-operation", "ab54-group", "ab54-register"} {
		ab54T.Run(ab54Axis, func(ab54T *ab54testing.T) {
			ab54Keys := []Ab54Key{ab54TestKey(4), ab54TestKey(4)}
			ab54Rows := make([]ab54json.RawMessage, 2)
			for ab54Index := range ab54Keys {
				switch ab54Axis {
				case "ab54-secondary":
					ab54Keys[ab54Index].Ab54Secondary = ab54Byte(uint8(ab54Index))
				case "ab54-operation":
					ab54Keys[ab54Index].Ab54Operation = ab54Byte(uint8(ab54Index))
				case "ab54-group":
					ab54Keys[ab54Index].Ab54Group = ab54Word(uint16(ab54Index))
				case "ab54-register":
					ab54Keys[ab54Index].Ab54Register = ab54Word(uint16(ab54Index))
				}
				ab54Row := ab54Object(ab54T, ab54Get(ab54T, ab54TestKey(4)))
				ab54Row["key"], _ = ab54json.Marshal(ab54Keys[ab54Index])
				ab54Row["id"] = ab54json.RawMessage(ab54fmt.Sprintf("\"ab54-synthetic-optional-%d\"", ab54Index))
				ab54Rows[ab54Index], _ = ab54json.Marshal(ab54Row)
			}
			ab54Source, _ := ab54json.Marshal(ab54Rows)
			for ab54Index, ab54Key := range ab54Keys {
				ab54Got, ab54Found, ab54Err := ab54Lookup(ab54Source, ab54Key)
				var ab54GotValue, ab54WantValue any
				if ab54Err != nil || !ab54Found {
					ab54T.Fatalf("ab54 explicit zero and nonzero collided: %v", ab54Err)
				}
				if ab54Err := ab54json.Unmarshal(ab54Got, &ab54GotValue); ab54Err != nil {
					ab54T.Fatal(ab54Err)
				}
				if ab54Err := ab54json.Unmarshal(ab54Rows[ab54Index], &ab54WantValue); ab54Err != nil {
					ab54T.Fatal(ab54Err)
				}
				if !ab54reflect.DeepEqual(ab54GotValue, ab54WantValue) {
					ab54T.Fatal("ab54 wrong optional value selected")
				}
			}
		})
	}
}
