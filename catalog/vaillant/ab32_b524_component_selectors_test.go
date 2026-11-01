package vaillant

import (
	ab32bytes "bytes"
	ab32json "encoding/json"
	ab32reflect "reflect"
	ab32testing "testing"
)

func ab32Pointer[ab32T ~uint8 | ~uint16](ab32Value ab32T) *ab32T { return &ab32Value }

func ab32BaseKey() Ab32Key {
	return Ab32Key{"controller-700", "none", 181, ab32Pointer(uint8(36)), ab32Pointer(uint8(2)), "read", "clock-date", "none", ab32Pointer(uint8(0)), "00", ab32Pointer(uint16(52))}
}

func ab32Get(ab32T *ab32testing.T, ab32Key Ab32Key) map[string]ab32json.RawMessage {
	ab32T.Helper()
	ab32Raw, ab32Found, ab32Error := Ab32Lookup(ab32Key)
	if ab32Error != nil || !ab32Found {
		ab32T.Fatalf("lookup: %v %v", ab32Found, ab32Error)
	}
	var ab32Row map[string]ab32json.RawMessage
	if ab32Error := ab32json.Unmarshal(ab32Raw, &ab32Row); ab32Error != nil {
		ab32T.Fatal(ab32Error)
	}
	return ab32Row
}

func TestAb32FiniteSelectors(ab32T *ab32testing.T) {
	ab32Entries, ab32Error := Ab32Entries()
	if ab32Error != nil {
		ab32T.Fatal(ab32Error)
	}
	var ab32Rows []struct {
		Ab32Key Ab32Key `json:"key"`
	}
	if ab32Error := ab32json.Unmarshal(ab32Entries, &ab32Rows); ab32Error != nil {
		ab32T.Fatal(ab32Error)
	}
	if len(ab32Rows) != 15 {
		ab32T.Fatalf("rows: %d", len(ab32Rows))
	}
	var ab32WholeRows []ab32json.RawMessage
	if ab32Error := ab32json.Unmarshal(ab32Entries, &ab32WholeRows); ab32Error != nil {
		ab32T.Fatal(ab32Error)
	}
	for ab32Index, ab32ExpectedRow := range ab32WholeRows {
		ab32ActualRow, ab32Found, ab32Error := Ab32Lookup(ab32Rows[ab32Index].Ab32Key)
		if ab32Error != nil || !ab32Found || !ab32bytes.Equal(ab32ActualRow, ab32ExpectedRow) {
			ab32T.Fatal("whole row was not retained")
		}
	}
	ab32Seen := make(map[string]bool)
	for _, ab32Family := range []string{"controller-700", "controller-720"} {
		ab32Specs := []struct {
			Ab32Field    string
			Ab32Group    uint8
			Ab32Scope    string
			Ab32Register uint16
		}{
			{"clock-date", 0, "00", 0x34}, {"clock-time", 0, "00", 0x35},
			{"system-holiday-start", 0, "00", 0x71}, {"system-holiday-end", 0, "00", 0x72},
			{"zone-holiday-start", 3, "dynamic", 3}, {"zone-holiday-end", 3, "dynamic", 4},
		}
		if ab32Family == "controller-720" {
			ab32Specs[2].Ab32Register = 0xda
			ab32Specs[3].Ab32Register = 0xdb
			ab32Specs = append(ab32Specs, struct {
				Ab32Field    string
				Ab32Group    uint8
				Ab32Scope    string
				Ab32Register uint16
			}{"zone-time-001e", 3, "dynamic", 0x1e}, struct {
				Ab32Field    string
				Ab32Group    uint8
				Ab32Scope    string
				Ab32Register uint16
			}{"zone-time-0020", 3, "dynamic", 0x20}, struct {
				Ab32Field    string
				Ab32Group    uint8
				Ab32Scope    string
				Ab32Register uint16
			}{"zone-time-0021", 3, "dynamic", 0x21})
		}
		for _, ab32Spec := range ab32Specs {
			ab32Key := ab32BaseKey()
			ab32Key.Ab32Family = ab32Family
			ab32Key.Ab32Field = ab32Spec.Ab32Field
			ab32Key.Ab32Group = ab32Pointer(ab32Spec.Ab32Group)
			ab32Key.Ab32InstanceScope = ab32Spec.Ab32Scope
			ab32Key.Ab32Register = ab32Pointer(ab32Spec.Ab32Register)
			ab32Row := ab32Get(ab32T, ab32Key)
			ab32ID := "ab32-" + ab32Family + "-" + ab32Spec.Ab32Field
			if string(ab32Row["id"]) != `"`+ab32ID+`"` || ab32Seen[ab32ID] {
				ab32T.Fatal("identity mismatch")
			}
			ab32Seen[ab32ID] = true
			var ab32Actual Ab32Key
			if ab32Error := ab32json.Unmarshal(ab32Row["key"], &ab32Actual); ab32Error != nil || !ab32reflect.DeepEqual(ab32Actual, ab32Key) {
				ab32T.Fatal("whole key lost")
			}
		}
	}
}

func TestAb32EachLookupAxis(ab32T *ab32testing.T) {
	ab32Changes := []struct {
		Ab32Name   string
		Ab32Change func(*Ab32Key)
		Ab32Found  bool
	}{
		{"ab32Family", func(ab32Key *Ab32Key) { ab32Key.Ab32Family = "controller-720" }, true},
		{"ab32Variant", func(ab32Key *Ab32Key) { ab32Key.Ab32Variant = "other" }, false},
		{"ab32Primary", func(ab32Key *Ab32Key) { ab32Key.Ab32Primary = 180 }, false},
		{"ab32Secondary", func(ab32Key *Ab32Key) { ab32Key.Ab32Secondary = ab32Pointer(uint8(35)) }, false},
		{"ab32Operation", func(ab32Key *Ab32Key) { ab32Key.Ab32Operation = ab32Pointer(uint8(6)) }, false},
		{"ab32Direction", func(ab32Key *Ab32Key) { ab32Key.Ab32Direction = "write" }, false},
		{"ab32Field", func(ab32Key *Ab32Key) { ab32Key.Ab32Field = "zone-name-suffix" }, false},
		{"ab32Alternative", func(ab32Key *Ab32Key) { ab32Key.Ab32Alternative = "other" }, false},
		{"ab32Group", func(ab32Key *Ab32Key) { ab32Key.Ab32Group = ab32Pointer(uint8(3)) }, false},
		{"ab32InstanceScope", func(ab32Key *Ab32Key) { ab32Key.Ab32InstanceScope = "dynamic" }, false},
		{"ab32Register", func(ab32Key *Ab32Key) { ab32Key.Ab32Register = ab32Pointer(uint16(53)) }, false},
	}
	for _, ab32Case := range ab32Changes {
		ab32T.Run(ab32Case.Ab32Name, func(ab32T *ab32testing.T) {
			ab32Key := ab32BaseKey()
			ab32Case.Ab32Change(&ab32Key)
			ab32Raw, ab32Found, ab32Error := Ab32Lookup(ab32Key)
			if ab32Error != nil || ab32Found != ab32Case.Ab32Found {
				ab32T.Fatalf("exact axis: %v %v", ab32Found, ab32Error)
			}
			if ab32Found && !ab32bytes.Contains(ab32Raw, []byte(`"ab32-controller-720-clock-date"`)) {
				ab32T.Fatal("family fallback")
			}
		})
	}
}

func TestAb32AbsentAndZero(ab32T *ab32testing.T) {
	for _, ab32Axis := range []string{"secondary", "operation", "group", "register"} {
		ab32T.Run("ab32"+ab32Axis, func(ab32T *ab32testing.T) {
			ab32Key := ab32BaseKey()
			ab32Set := func(ab32Key *Ab32Key, ab32Present bool) {
				switch ab32Axis {
				case "secondary":
					ab32Key.Ab32Secondary = nil
					if ab32Present {
						ab32Key.Ab32Secondary = ab32Pointer(uint8(0))
					}
				case "operation":
					ab32Key.Ab32Operation = nil
					if ab32Present {
						ab32Key.Ab32Operation = ab32Pointer(uint8(0))
					}
				case "group":
					ab32Key.Ab32Group = nil
					if ab32Present {
						ab32Key.Ab32Group = ab32Pointer(uint8(0))
					}
				case "register":
					ab32Key.Ab32Register = nil
					if ab32Present {
						ab32Key.Ab32Register = ab32Pointer(uint16(0))
					}
				}
			}
			ab32Set(&ab32Key, false)
			ab32Other := ab32BaseKey()
			ab32Set(&ab32Other, true)
			ab32Data, ab32Error := ab32json.Marshal([]any{map[string]any{"id": "ab32-absent", "key": ab32Key}, map[string]any{"id": "ab32-zero", "key": ab32Other}})
			if ab32Error != nil {
				ab32T.Fatal(ab32Error)
			}
			for ab32Index, ab32Query := range []Ab32Key{ab32Key, ab32Other} {
				ab32Raw, ab32Found, ab32Error := ab32LookupData(ab32Data, ab32Query)
				ab32Want := []string{`"ab32-absent"`, `"ab32-zero"`}[ab32Index]
				if ab32Error != nil || !ab32Found || !ab32bytes.Contains(ab32Raw, []byte(ab32Want)) {
					ab32T.Fatalf("absent/zero: %s %v", ab32Raw, ab32Error)
				}
			}
		})
	}
}

func TestAb32ValidationAndCopies(ab32T *ab32testing.T) {
	ab32Key := ab32BaseKey()
	ab32Data, ab32Error := Ab32Entries()
	if ab32Error != nil {
		ab32T.Fatal(ab32Error)
	}
	ab32Original := append([]byte(nil), ab32Data...)
	ab32Row, ab32Found, ab32Error := Ab32Lookup(ab32Key)
	if ab32Error != nil || !ab32Found {
		ab32T.Fatal("lookup")
	}
	ab32OriginalRow := append([]byte(nil), ab32Row...)
	for ab32Index := range ab32Data {
		ab32Data[ab32Index] = 'x'
	}
	for ab32Index := range ab32Row {
		ab32Row[ab32Index] = 'x'
	}
	ab32Again, ab32Error := Ab32Entries()
	if ab32Error != nil || !ab32bytes.Equal(ab32Again, ab32Original) {
		ab32T.Fatal("entries alias")
	}
	ab32Again, ab32Found, ab32Error = Ab32Lookup(ab32Key)
	if ab32Error != nil || !ab32Found || !ab32bytes.Equal(ab32Again, ab32OriginalRow) {
		ab32T.Fatal("row alias")
	}
	ab32Duplicate, _ := ab32json.Marshal([]any{map[string]any{"id": "ab32-first", "key": ab32Key}, map[string]any{"id": "ab32-second", "key": ab32Key}})
	ab32Invalid := [][]byte{[]byte(`{`), []byte(`null`), []byte(`{}`), []byte(`[null]`), []byte(`[{"id":"ab32-x","key":{}}]`), ab32Duplicate}
	for ab32Index, ab32Bad := range ab32Invalid {
		if _, _, ab32Error := ab32Validate(ab32Bad); ab32Error == nil {
			ab32T.Fatalf("invalid accepted %d", ab32Index)
		}
		if _, ab32Found, ab32Error := ab32LookupData(ab32Bad, ab32Key); ab32Error == nil || ab32Found {
			ab32T.Fatalf("lookup invalid accepted %d", ab32Index)
		}
	}
}

func TestAb32ComponentMetadata(ab32T *ab32testing.T) {
	ab32Data, ab32Error := Ab32Entries()
	if ab32Error != nil {
		ab32T.Fatal(ab32Error)
	}
	var ab32Rows []struct {
		Ab32Key        Ab32Key `json:"key"`
		Ab32Grade      string  `json:"grade"`
		Ab32Kind       string  `json:"kind"`
		Ab32Components []struct {
			Ab32Name     string `json:"name"`
			Ab32Position int    `json:"position"`
			Ab32Encoding string `json:"encoding"`
			Ab32Scale    int    `json:"scale"`
			Ab32Offset   int    `json:"offset"`
			Ab32Minimum  *int   `json:"named_minimum"`
			Ab32Maximum  *int   `json:"named_maximum"`
			Ab32Raw      string `json:"raw_form"`
		} `json:"components"`
		Ab32Relations map[string]any    `json:"relations"`
		Ab32Unknown   map[string]string `json:"unknown"`
		Ab32Raw       string            `json:"raw_form"`
	}
	if ab32Error := ab32json.Unmarshal(ab32Data, &ab32Rows); ab32Error != nil {
		ab32T.Fatal(ab32Error)
	}
	for _, ab32Row := range ab32Rows {
		ab32Lookup := ab32Get(ab32T, ab32Row.Ab32Key)
		if string(ab32Lookup["grade"]) != `"Hypothesis"` || ab32Row.Ab32Grade != "Hypothesis" {
			ab32T.Fatal("grade")
		}
		ab32WantRelations := map[string]any{"selector_context": "whole exact key", "paired_selector": nil, "component_order": "declared order", "failure_scope": "individual component", "retain_unaffected": true, "retain_raw": true, "missing_component": "explicit missing cell", "complete_timestamp": false, "freshness_inferred": false, "constraint_admission": false}
		if !ab32reflect.DeepEqual(ab32Row.Ab32Relations, ab32WantRelations) {
			ab32T.Fatal("component retention or unrelated pairing")
		}
		ab32UnknownNames := []string{"native_support", "physical_effects", "persistence", "delivery", "century", "timezone", "joint_date_validity", "zone_name_selectors", "string_encoding", "string_width", "string_termination"}
		if len(ab32Row.Ab32Unknown) != len(ab32UnknownNames) {
			ab32T.Fatal("unknown coverage")
		}
		for _, ab32Name := range ab32UnknownNames {
			if ab32Row.Ab32Unknown[ab32Name] != "Unknown" {
				ab32T.Fatal("lost Unknown", ab32Name)
			}
		}
		if ab32Row.Ab32Raw != "ordered component bytes with independent missing or failed cells" || len(ab32Row.Ab32Components) != 3 {
			ab32T.Fatal("raw or count")
		}
		ab32Names := []string{"day", "month", "year"}
		ab32Bytes := []byte{0x1f, 0x0c, 0x63}
		ab32Expected := []int{31, 12, 99}
		if ab32Row.Ab32Kind == "time" {
			ab32Names = []string{"hour", "minute", "second"}
			ab32Bytes = []byte{0x17, 0x3b, 0x3a}
			ab32Expected = []int{23, 59, 58}
		} else if ab32Row.Ab32Kind != "date" {
			ab32T.Fatal("kind")
		}
		for ab32Index, ab32Component := range ab32Row.Ab32Components {
			if ab32Component.Ab32Name != ab32Names[ab32Index] || ab32Component.Ab32Position != ab32Index || ab32Component.Ab32Encoding != "u8" || ab32Component.Ab32Scale != 1 || ab32Component.Ab32Offset != 0 || ab32Component.Ab32Raw != "individual byte" {
				ab32T.Fatal("component representation")
			}
			if int(ab32Bytes[ab32Index])*ab32Component.Ab32Scale+ab32Component.Ab32Offset != ab32Expected[ab32Index] || int(byte(0x12))*ab32Component.Ab32Scale+ab32Component.Ab32Offset != 18 {
				ab32T.Fatal("plain byte examples")
			}
			if ab32Row.Ab32Kind == "date" {
				if ab32Component.Ab32Minimum == nil || ab32Component.Ab32Maximum == nil || *ab32Component.Ab32Minimum != []int{1, 1, 0}[ab32Index] || *ab32Component.Ab32Maximum != []int{31, 12, 99}[ab32Index] {
					ab32T.Fatal("named domains")
				}
			}
		}
	}
}

func TestAb32NoSelectorFallback(ab32T *ab32testing.T) {
	for _, ab32Case := range []struct {
		Ab32Family   string
		Ab32Field    string
		Ab32Register uint16
		Ab32Group    uint8
		Ab32Scope    string
	}{
		{"controller-700", "system-holiday-start", 0xda, 0, "00"}, {"controller-720", "system-holiday-start", 0x71, 0, "00"},
		{"controller-700", "system-holiday-end", 0xdb, 0, "00"}, {"controller-720", "system-holiday-end", 0x72, 0, "00"},
		{"controller-700", "zone-time-001e", 0x1e, 3, "dynamic"}, {"controller-700", "zone-time-0020", 0x20, 3, "dynamic"}, {"controller-700", "zone-time-0021", 0x21, 3, "dynamic"},
		{"controller-720", "zone-name-suffix", 3, 3, "dynamic"},
	} {
		ab32Key := ab32BaseKey()
		ab32Key.Ab32Family = ab32Case.Ab32Family
		ab32Key.Ab32Field = ab32Case.Ab32Field
		ab32Key.Ab32Register = ab32Pointer(ab32Case.Ab32Register)
		ab32Key.Ab32Group = ab32Pointer(ab32Case.Ab32Group)
		ab32Key.Ab32InstanceScope = ab32Case.Ab32Scope
		if _, ab32Found, ab32Error := Ab32Lookup(ab32Key); ab32Found || ab32Error != nil {
			ab32T.Fatal("forbidden selector fallback")
		}
	}
}

func TestAb32PartialCells(ab32T *ab32testing.T) {
	ab32Row := ab32Get(ab32T, ab32BaseKey())
	var ab32Components []struct {
		Ab32Name     string `json:"name"`
		Ab32Position int    `json:"position"`
	}
	if ab32Error := ab32json.Unmarshal(ab32Row["components"], &ab32Components); ab32Error != nil {
		ab32T.Fatal(ab32Error)
	}
	ab32Cells := []*uint8{ab32Pointer(uint8(31)), nil, ab32Pointer(uint8(99))}
	ab32Expected := map[string]*uint8{"day": ab32Pointer(uint8(31)), "month": nil, "year": ab32Pointer(uint8(99))}
	for _, ab32Component := range ab32Components {
		ab32Want, ab32Present := ab32Expected[ab32Component.Ab32Name]
		if !ab32Present || ab32Component.Ab32Position < 0 || ab32Component.Ab32Position >= len(ab32Cells) || !ab32reflect.DeepEqual(ab32Cells[ab32Component.Ab32Position], ab32Want) {
			ab32T.Fatal("partial raw component mapping")
		}
	}
	var ab32Relations map[string]any
	if ab32Error := ab32json.Unmarshal(ab32Row["relations"], &ab32Relations); ab32Error != nil {
		ab32T.Fatal(ab32Error)
	}
	if ab32Relations["retain_unaffected"] != true || ab32Relations["retain_raw"] != true || ab32Relations["failure_scope"] != "individual component" || ab32Relations["missing_component"] != "explicit missing cell" || ab32Relations["complete_timestamp"] != false || ab32Relations["freshness_inferred"] != false {
		ab32T.Fatal("partial failure policy")
	}
	if _, ab32Present := ab32Relations["fresh"]; ab32Present {
		ab32T.Fatal("invented freshness")
	}
}
