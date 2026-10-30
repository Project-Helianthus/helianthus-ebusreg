package ebus_standard_catalog

import (
	ab28bytes "bytes"
	ab28json "encoding/json"
	ab28fmt "fmt"
	ab28reflect "reflect"
	ab28strings "strings"
	ab28testing "testing"
)

func ab28TestKey(ab28Secondary uint8, ab28Direction, ab28Field string) Ab28Key {
	return Ab28Key{Ab28Family: "standard", Ab28Variant: "standard07", Ab28Primary: 7, Ab28Secondary: &ab28Secondary, Ab28Direction: ab28Direction, Ab28Field: ab28Field, Ab28Alternative: "none", Ab28InstanceScope: "none"}
}

func ab28TestRow(ab28T *ab28testing.T, ab28Key Ab28Key) map[string]ab28json.RawMessage {
	ab28T.Helper()
	ab28Row, ab28Found, ab28Err := Ab28Lookup(ab28Key)
	if ab28Err != nil || !ab28Found {
		ab28T.Fatalf("lookup %+v: found=%v error=%v", ab28Key, ab28Found, ab28Err)
	}
	var ab28Object map[string]ab28json.RawMessage
	if ab28Err := ab28json.Unmarshal(ab28Row, &ab28Object); ab28Err != nil {
		ab28T.Fatal(ab28Err)
	}
	return ab28Object
}

func ab28Assert(ab28T *ab28testing.T, ab28Object map[string]ab28json.RawMessage, ab28Field, ab28Expected string) {
	ab28T.Helper()
	var ab28ActualValue, ab28ExpectedValue any
	if ab28Err := ab28json.Unmarshal(ab28Object[ab28Field], &ab28ActualValue); ab28Err != nil {
		ab28T.Fatalf("%s: %v", ab28Field, ab28Err)
	}
	if ab28Err := ab28json.Unmarshal([]byte(ab28Expected), &ab28ExpectedValue); ab28Err != nil {
		ab28T.Fatal(ab28Err)
	}
	if !ab28reflect.DeepEqual(ab28ActualValue, ab28ExpectedValue) {
		ab28T.Fatalf("%s: got %s, want %s", ab28Field, ab28Object[ab28Field], ab28Expected)
	}
}

func TestAb28FiniteInventory(ab28T *ab28testing.T) {
	ab28Data, ab28Err := Ab28Entries()
	if ab28Err != nil {
		ab28T.Fatal(ab28Err)
	}
	var ab28Rows []struct {
		Ab28ID  string  `json:"id"`
		Ab28Key Ab28Key `json:"key"`
	}
	if ab28Err := ab28json.Unmarshal(ab28Data, &ab28Rows); ab28Err != nil {
		ab28T.Fatal(ab28Err)
	}
	if len(ab28Rows) != 37 {
		ab28T.Fatalf("finite inventory: %d", len(ab28Rows))
	}
	ab28Seen := make(map[string]bool)
	for _, ab28Row := range ab28Rows {
		if !ab28strings.HasPrefix(ab28Row.Ab28ID, "ab28-") || ab28Seen[ab28Row.Ab28ID] {
			ab28T.Fatalf("invalid or repeated ID: %s", ab28Row.Ab28ID)
		}
		ab28Seen[ab28Row.Ab28ID] = true
		if ab28Row.Ab28Key.Ab28Secondary == nil || !ab28bytes.Contains([]byte{0, 2, 3, 5, 254, 255}, []byte{*ab28Row.Ab28Key.Ab28Secondary}) {
			ab28T.Fatal("invented or out-of-slice service")
		}
		ab28Object := ab28TestRow(ab28T, ab28Row.Ab28Key)
		ab28Assert(ab28T, ab28Object, "grade", `"Hypothesis"`)
		ab28Assert(ab28T, ab28Object, "unknown", `{"grade":"Unknown","native_support":"Unknown","physical_effects":"Unknown","persistence":"Unknown","delivery":"Unknown","qualification":"Unknown"}`)
		ab28Assert(ab28T, ab28Object, "id", ab28fmt.Sprintf("%q", ab28Row.Ab28ID))
		var ab28Raw map[string]ab28json.RawMessage
		if ab28Err := ab28json.Unmarshal(ab28Object["raw"], &ab28Raw); ab28Err != nil {
			ab28T.Fatal(ab28Err)
		}
		ab28Assert(ab28T, ab28Raw, "retained", `true`)
		ab28Assert(ab28T, ab28Raw, "coordinates", `"after PB/SB"`)
	}
	for _, ab28Secondary := range []uint8{1, 4} {
		if _, ab28Found, ab28Err := Ab28Lookup(ab28TestKey(ab28Secondary, "request", "temperature")); ab28Err != nil || ab28Found {
			ab28T.Fatal("out-of-slice body admitted")
		}
	}
}

func TestAb28ResidualFields(ab28T *ab28testing.T) {
	ab28Weekday := ab28TestRow(ab28T, ab28TestKey(0, "request", "weekday"))
	ab28Assert(ab28T, ab28Weekday, "encoding", `"raw u8"`)
	ab28Assert(ab28T, ab28Weekday, "offset", `7`)
	ab28Assert(ab28T, ab28Weekday, "width", `1`)
	ab28Assert(ab28T, ab28Weekday, "labels", `{"0":"Monday","1":"Tuesday","2":"Wednesday","3":"Thursday","4":"Friday","5":"Saturday","6":"Sunday"}`)
	ab28Assert(ab28T, ab28Weekday, "unlabelled_minimum", `7`)
	ab28Assert(ab28T, ab28Weekday, "unlabelled_maximum", `255`)
	ab28Assert(ab28T, ab28Weekday, "examples", `[{"raw":[6],"label":"Sunday"},{"raw":[7],"label":null},{"raw":[255],"label":null}]`)
	for _, ab28Case := range []struct {
		Ab28Field                            string
		Ab28Offset, Ab28Minimum, Ab28Maximum int
	}{{"second", 2, 0, 59}, {"minute", 3, 0, 59}, {"hour", 4, 0, 23}, {"day", 5, 1, 31}, {"month", 6, 1, 12}, {"year", 8, 0, 99}} {
		ab28Object := ab28TestRow(ab28T, ab28TestKey(0, "request", ab28Case.Ab28Field))
		ab28Assert(ab28T, ab28Object, "encoding", `"packed-decimal"`)
		ab28Assert(ab28T, ab28Object, "width", `1`)
		for ab28Field, ab28Value := range map[string]int{"offset": ab28Case.Ab28Offset, "minimum": ab28Case.Ab28Minimum, "maximum": ab28Case.Ab28Maximum} {
			ab28Assert(ab28T, ab28Object, ab28Field, ab28fmt.Sprint(ab28Value))
		}
		for _, ab28Field := range []string{"century", "timezone", "calendar_validation", "complete_native_decimal_conversion"} {
			ab28Assert(ab28T, ab28Object, ab28Field, `"Unknown"`)
		}
		ab28Assert(ab28T, ab28Object, "executable_schema", `false`)
	}
	for _, ab28Secondary := range []uint8{0, 2} {
		ab28Field := "temperature"
		if ab28Secondary == 0 {
			ab28Field = "outside_temperature"
		}
		ab28Object := ab28TestRow(ab28T, ab28TestKey(ab28Secondary, "request", ab28Field))
		for ab28Field, ab28Expected := range map[string]string{"offset": "0", "width": "2", "encoding": `"signed LE16 /256"`, "unit": `"degrees Celsius"`, "raw_minimum": "-12800", "raw_maximum": "12800", "minimum": "-50", "maximum": "50", "sentinel_8000": `"Unknown"`, "examples": `[{"raw":[0,50],"value":50},{"raw":[0,206],"value":-50}]`} {
			ab28Assert(ab28T, ab28Object, ab28Field, ab28Expected)
		}
	}
	ab28Duration := ab28TestRow(ab28T, ab28TestKey(2, "request", "duration"))
	for ab28Field, ab28Expected := range map[string]string{"offset": "2", "width": "1", "encoding": `"raw u8"`, "unit": `"minutes"`, "default": "0", "zero_permanent_interpretation": `"Unknown"`, "examples": `[{"raw":[0],"value":0}]`} {
		ab28Assert(ab28T, ab28Duration, ab28Field, ab28Expected)
	}
	for _, ab28Case := range []struct {
		Ab28Secondary uint8
		Ab28Literal   int
	}{{0, 6}, {2, 1}} {
		ab28Object := ab28TestRow(ab28T, ab28TestKey(ab28Case.Ab28Secondary, "none", "unassociated_default"))
		ab28Assert(ab28T, ab28Object, "literal", ab28fmt.Sprint(ab28Case.Ab28Literal))
		ab28Assert(ab28T, ab28Object, "field_association", `"Unknown"`)
	}
	for _, ab28Key := range []Ab28Key{ab28TestKey(0, "request", "weekday"), ab28TestKey(0, "request", "outside_temperature"), ab28TestKey(2, "request", "temperature")} {
		if _, ab28OK := ab28TestRow(ab28T, ab28Key)["default"]; ab28OK {
			ab28T.Fatal("unassociated default attached to field")
		}
	}
}

func TestAb28BlocksAndPairs(ab28T *ab28testing.T) {
	for _, ab28Case := range []struct {
		Ab28Secondary         uint8
		Ab28Field             string
		Ab28Offset, Ab28Width int
	}{{3, "block", 0, 1}, {5, "secondary_block", 0, 2}, {5, "primary_block", 1, 2}} {
		ab28Object := ab28TestRow(ab28T, ab28TestKey(ab28Case.Ab28Secondary, "request", ab28Case.Ab28Field))
		for ab28Field, ab28Expected := range map[string]string{"offset": ab28fmt.Sprint(ab28Case.Ab28Offset), "width": "1", "encoding": `"raw u8"`, "request_body_width": ab28fmt.Sprint(ab28Case.Ab28Width)} {
			ab28Assert(ab28T, ab28Object, ab28Field, ab28Expected)
		}
	}
	for _, ab28Secondary := range []uint8{3, 5} {
		for _, ab28Case := range []struct {
			Ab28Field    string
			Ab28Offset   int
			Ab28Examples string
		}{{"version", 0, `[{"raw":[18],"value":12}]`}, {"revision", 1, `[{"raw":[52],"value":34}]`}} {
			ab28Object := ab28TestRow(ab28T, ab28TestKey(ab28Secondary, "positive", ab28Case.Ab28Field))
			for ab28Field, ab28Expected := range map[string]string{"offset": ab28fmt.Sprint(ab28Case.Ab28Offset), "width": "1", "encoding": `"packed-decimal"`, "independent": "true", "combined_api_number": "false", "positive_body_width": "10", "version_control_overlap_0704": `"Unknown"`, "complete_native_decimal_conversion": `"Unknown"`, "executable_schema": "false", "examples": ab28Case.Ab28Examples} {
				ab28Assert(ab28T, ab28Object, ab28Field, ab28Expected)
			}
		}
		for ab28Index := 0; ab28Index < 8; ab28Index++ {
			ab28Object := ab28TestRow(ab28T, ab28TestKey(ab28Secondary, "positive", ab28fmt.Sprintf("bitmap_%d", ab28Index)))
			for ab28Field, ab28Expected := range map[string]string{"offset": ab28fmt.Sprint(ab28Index + 2), "width": "1", "encoding": `"raw bitmap u8"`, "independent": "true", "bit_selection_algorithm": `"Unknown"`, "admitted_methods": "[]", "examples": `[{"raw":[128],"value":128,"selected_service":"Unknown"}]`} {
				ab28Assert(ab28T, ab28Object, ab28Field, ab28Expected)
			}
		}
	}
	for _, ab28Case := range []struct {
		Ab28Secondary uint8
		Ab28Field     string
	}{{254, "inquiry"}, {255, "sign_of_life"}} {
		ab28Object := ab28TestRow(ab28T, ab28TestKey(ab28Case.Ab28Secondary, "request", ab28Case.Ab28Field))
		for ab28Field, ab28Expected := range map[string]string{"argument_body_width": "0", "emission": `"Unknown"`, "physical_response": `"Unknown"`, "direct_inquiry_reply": `"Unknown"`, "raw": ab28fmt.Sprintf(`{"retained":true,"coordinates":"after PB/SB","pair":[7,%d],"body":[]}`, ab28Case.Ab28Secondary)} {
			ab28Assert(ab28T, ab28Object, ab28Field, ab28Expected)
		}
	}
}

func TestAb28IndependentSelectors(ab28T *ab28testing.T) {
	for _, ab28Case := range []struct {
		Ab28Secondary           uint8
		Ab28Field, Ab28Examples string
	}{
		{3, "block", `[{"raw":[1],"value":1},{"raw":[2],"value":2}]`},
		{5, "secondary_block", `[{"raw":[1,2],"value":1},{"raw":[3,2],"value":3}]`},
		{5, "primary_block", `[{"raw":[1,2],"value":2},{"raw":[1,3],"value":3}]`},
	} {
		ab28Assert(ab28T, ab28TestRow(ab28T, ab28TestKey(ab28Case.Ab28Secondary, "request", ab28Case.Ab28Field)), "examples", ab28Case.Ab28Examples)
	}
}

func TestAb28ExactAxes(ab28T *ab28testing.T) {
	ab28Zero8 := uint8(0)
	ab28Zero16 := uint16(0)
	ab28One8 := uint8(1)
	ab28Changes := []struct {
		Ab28Name   string
		Ab28Change func(*Ab28Key)
	}{
		{"family", func(ab28Key *Ab28Key) { ab28Key.Ab28Family = "synthetic" }},
		{"variant", func(ab28Key *Ab28Key) { ab28Key.Ab28Variant = "synthetic" }},
		{"primary", func(ab28Key *Ab28Key) { ab28Key.Ab28Primary = 8 }},
		{"secondary", func(ab28Key *Ab28Key) { ab28Key.Ab28Secondary = &ab28One8 }},
		{"operation", func(ab28Key *Ab28Key) { ab28Key.Ab28Operation = &ab28Zero8 }},
		{"direction", func(ab28Key *Ab28Key) { ab28Key.Ab28Direction = "positive" }},
		{"field", func(ab28Key *Ab28Key) { ab28Key.Ab28Field = "synthetic" }},
		{"alternative", func(ab28Key *Ab28Key) { ab28Key.Ab28Alternative = "synthetic" }},
		{"group", func(ab28Key *Ab28Key) { ab28Key.Ab28Group = &ab28Zero8 }},
		{"instance_scope", func(ab28Key *Ab28Key) { ab28Key.Ab28InstanceScope = "dynamic" }},
		{"register", func(ab28Key *Ab28Key) { ab28Key.Ab28Register = &ab28Zero16 }},
		{"secondary_absent_zero", func(ab28Key *Ab28Key) { ab28Key.Ab28Secondary = nil }},
	}
	for _, ab28Case := range ab28Changes {
		ab28T.Run("ab28_"+ab28Case.Ab28Name, func(ab28T *ab28testing.T) {
			ab28Base := ab28TestKey(0, "request", "outside_temperature")
			ab28Changed := ab28Base
			ab28Case.Ab28Change(&ab28Changed)
			ab28Original, ab28Found, ab28Err := Ab28Lookup(ab28Base)
			if ab28Err != nil || !ab28Found {
				ab28T.Fatalf("base control: %v %v", ab28Found, ab28Err)
			}
			if ab28Row, ab28Found, ab28Err := Ab28Lookup(ab28Changed); ab28Err != nil || ab28Found || ab28Row != nil {
				ab28T.Fatalf("non-exact key matched: %s %v %v", ab28Row, ab28Found, ab28Err)
			}
			var ab28Compact ab28bytes.Buffer
			if ab28Err := ab28json.Compact(&ab28Compact, ab28Original); ab28Err != nil {
				ab28T.Fatal(ab28Err)
			}
			ab28Original = ab28Compact.Bytes()
			var ab28Object map[string]ab28json.RawMessage
			if ab28Err := ab28json.Unmarshal(ab28Original, &ab28Object); ab28Err != nil {
				ab28T.Fatal(ab28Err)
			}
			ab28ChangedJSON, ab28Err := ab28json.Marshal(ab28Changed)
			if ab28Err != nil {
				ab28T.Fatal(ab28Err)
			}
			ab28Object["key"] = ab28ChangedJSON
			ab28Object["id"] = ab28json.RawMessage(`"ab28-synthetic-control"`)
			ab28Other, ab28Err := ab28json.Marshal(ab28Object)
			if ab28Err != nil {
				ab28T.Fatal(ab28Err)
			}
			ab28Controls, ab28Err := ab28json.Marshal([]ab28json.RawMessage{ab28Original, ab28Other})
			if ab28Err != nil {
				ab28T.Fatal(ab28Err)
			}
			for _, ab28Control := range []struct {
				Ab28Key      Ab28Key
				Ab28Expected ab28json.RawMessage
			}{{ab28Base, ab28Original}, {ab28Changed, ab28Other}} {
				ab28Row, ab28Found, ab28Err := ab28Lookup(ab28Controls, ab28Control.Ab28Key)
				if ab28Err != nil || !ab28Found || !ab28bytes.Equal(ab28Row, ab28Control.Ab28Expected) {
					ab28T.Fatalf("distinct valid control: %s %v %v", ab28Row, ab28Found, ab28Err)
				}
			}
		})
	}
	// ab28: independent nonzero optional controls prevent presence-only matching.
	for _, ab28Axis := range []string{"secondary", "operation", "group", "register"} {
		ab28T.Run("ab28_"+ab28Axis+"_value", func(ab28T *ab28testing.T) {
			ab28Base := ab28TestKey(0, "request", "outside_temperature")
			ab28Base.Ab28Operation = &ab28Zero8
			ab28Base.Ab28Group = &ab28Zero8
			ab28Base.Ab28Register = &ab28Zero16
			ab28Changed := ab28Base
			ab28One16 := uint16(1)
			switch ab28Axis {
			case "secondary":
				ab28Changed.Ab28Secondary = &ab28One8
			case "operation":
				ab28Changed.Ab28Operation = &ab28One8
			case "group":
				ab28Changed.Ab28Group = &ab28One8
			case "register":
				ab28Changed.Ab28Register = &ab28One16
			}
			ab28Controls, ab28Err := ab28json.Marshal([]map[string]any{{"id": "ab28-zero", "key": ab28Base}, {"id": "ab28-one", "key": ab28Changed}})
			if ab28Err != nil {
				ab28T.Fatal(ab28Err)
			}
			for _, ab28Control := range []struct {
				Ab28Key Ab28Key
				Ab28ID  string
			}{{ab28Base, "ab28-zero"}, {ab28Changed, "ab28-one"}} {
				ab28Row, ab28Found, ab28Err := ab28Lookup(ab28Controls, ab28Control.Ab28Key)
				if ab28Err != nil || !ab28Found {
					ab28T.Fatalf("optional value: %v %v", ab28Found, ab28Err)
				}
				var ab28Object map[string]ab28json.RawMessage
				if ab28Err := ab28json.Unmarshal(ab28Row, &ab28Object); ab28Err != nil {
					ab28T.Fatal(ab28Err)
				}
				ab28Assert(ab28T, ab28Object, "id", ab28fmt.Sprintf("%q", ab28Control.Ab28ID))
			}
		})
	}
}

func TestAb28ValidationAndCopies(ab28T *ab28testing.T) {
	ab28Key := ab28TestKey(0, "request", "outside_temperature")
	ab28Row, ab28Found, ab28Err := Ab28Lookup(ab28Key)
	if ab28Err != nil || !ab28Found {
		ab28T.Fatal(ab28Err)
	}
	ab28Duplicate, ab28Err := ab28json.Marshal([]ab28json.RawMessage{ab28Row, ab28Row})
	if ab28Err != nil {
		ab28T.Fatal(ab28Err)
	}
	for _, ab28Bad := range [][]byte{[]byte(`{`), []byte(`{}`), []byte(`null`), []byte(`[null]`), []byte(`[{}]`), []byte(`[{"key":{"primary":256}}]`), ab28Duplicate} {
		if _, ab28Err := ab28Validate(ab28Bad); ab28Err == nil {
			ab28T.Fatalf("invalid data admitted: %s", ab28Bad)
		}
		if ab28Result, ab28Found, ab28Err := ab28Lookup(ab28Bad, ab28Key); ab28Err == nil || ab28Found || ab28Result != nil {
			ab28T.Fatal("invalid catalog did not fail closed")
		}
	}
	var ab28Object map[string]ab28json.RawMessage
	if ab28Err := ab28json.Unmarshal(ab28Row, &ab28Object); ab28Err != nil {
		ab28T.Fatal(ab28Err)
	}
	var ab28Axes map[string]ab28json.RawMessage
	if ab28Err := ab28json.Unmarshal(ab28Object["key"], &ab28Axes); ab28Err != nil {
		ab28T.Fatal(ab28Err)
	}
	for ab28Axis, ab28Value := range ab28Axes {
		delete(ab28Axes, ab28Axis)
		ab28KeyJSON, ab28Err := ab28json.Marshal(ab28Axes)
		if ab28Err != nil {
			ab28T.Fatal(ab28Err)
		}
		ab28Object["key"] = ab28KeyJSON
		ab28Bad, ab28Err := ab28json.Marshal([]map[string]ab28json.RawMessage{ab28Object})
		if ab28Err != nil {
			ab28T.Fatal(ab28Err)
		}
		if _, ab28Err := ab28Validate(ab28Bad); ab28Err == nil {
			ab28T.Fatalf("missing %s admitted", ab28Axis)
		}
		ab28Axes[ab28Axis] = ab28Value
	}
	ab28Entries, ab28Err := Ab28Entries()
	if ab28Err != nil {
		ab28T.Fatal(ab28Err)
	}
	ab28EntriesAgain, ab28Err := Ab28Entries()
	if ab28Err != nil {
		ab28T.Fatal(ab28Err)
	}
	ab28RowAgain, ab28Found, ab28Err := Ab28Lookup(ab28Key)
	if ab28Err != nil || !ab28Found {
		ab28T.Fatal(ab28Err)
	}
	ab28ExpectedRow := append([]byte(nil), ab28Row...)
	ab28ExpectedEntries := append([]byte(nil), ab28Entries...)
	for ab28Index := range ab28Row {
		ab28Row[ab28Index] = '!'
	}
	for ab28Index := range ab28Entries {
		ab28Entries[ab28Index] = '!'
	}
	if !ab28bytes.Equal(ab28RowAgain, ab28ExpectedRow) || !ab28bytes.Equal(ab28EntriesAgain, ab28ExpectedEntries) {
		ab28T.Fatal("simultaneous results share memory")
	}
	ab28Later, ab28Found, ab28Err := Ab28Lookup(ab28Key)
	if ab28Err != nil || !ab28Found || !ab28bytes.Equal(ab28Later, ab28ExpectedRow) {
		ab28T.Fatal("row mutation affected later lookup")
	}
	ab28LaterEntries, ab28Err := Ab28Entries()
	if ab28Err != nil || !ab28bytes.Equal(ab28LaterEntries, ab28ExpectedEntries) {
		ab28T.Fatal("buffer mutation affected later entries")
	}
	ab28Source := append([]byte{'['}, ab28ExpectedRow...)
	ab28Source = append(ab28Source, ']')
	ab28Isolated, ab28Found, ab28Err := ab28Lookup(ab28Source, ab28Key)
	if ab28Err != nil || !ab28Found {
		ab28T.Fatal(ab28Err)
	}
	for ab28Index := range ab28Source {
		ab28Source[ab28Index] = '!'
	}
	if !ab28bytes.Equal(ab28Isolated, ab28ExpectedRow) {
		ab28T.Fatal("lookup aliases input")
	}
}
