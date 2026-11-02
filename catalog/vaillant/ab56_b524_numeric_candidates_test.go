package vaillant

import (
	ab56Bytes "bytes"
	ab56JSON "encoding/json"
	ab56Reflect "reflect"
	ab56Testing "testing"
)

func ab56Pointer[ab56Type any](ab56Value ab56Type) *ab56Type { return &ab56Value }

func ab56TestContext() Ab56Context {
	return Ab56Context{"ab56_family", "ab56_software", "ab56_hardware", "ab56_api_version", "ab56_api_revision"}
}

func ab56TestKey(ab56Group, ab56Instance uint8, ab56Register uint16) Ab56Key {
	return Ab56Key{
		Ab56Family: "ab56_family", Ab56Software: "ab56_software", Ab56Hardware: "ab56_hardware",
		Ab56APIVersion: "ab56_api_version", Ab56APIRevision: "ab56_api_revision",
		Ab56Variant: "ab56_numeric_candidates", Ab56CandidateSet: "ab56_b524_twenty_unadmitted",
		Ab56Primary: ab56Pointer(uint16(181)), Ab56Secondary: ab56Pointer(uint8(36)),
		Ab56Operation: ab56Pointer(uint8(2)), Ab56Opcode: ab56Pointer(uint8(2)),
		Ab56Direction: "none", Ab56Field: "ab56_numeric_candidate", Ab56Alternative: "none",
		Ab56Group: ab56Pointer(ab56Group), Ab56InstanceScope: "ab56_candidate_slot",
		Ab56Instance: ab56Pointer(ab56Instance), Ab56Register: ab56Pointer(ab56Register),
	}
}

func ab56Decode(ab56T *ab56Testing.T, ab56Raw []byte) map[string]ab56JSON.RawMessage {
	ab56T.Helper()
	var ab56Object map[string]ab56JSON.RawMessage
	if ab56Error := ab56JSON.Unmarshal(ab56Raw, &ab56Object); ab56Error != nil {
		ab56T.Fatal(ab56Error)
	}
	return ab56Object
}

func ab56Encode(ab56T *ab56Testing.T, ab56Value any) []byte {
	ab56T.Helper()
	ab56Raw, ab56Error := ab56JSON.Marshal(ab56Value)
	if ab56Error != nil {
		ab56T.Fatal(ab56Error)
	}
	return ab56Raw
}

func ab56Equal(ab56T *ab56Testing.T, ab56Actual, ab56Expected []byte) {
	ab56T.Helper()
	var ab56Left, ab56Right any
	if ab56Error := ab56JSON.Unmarshal(ab56Actual, &ab56Left); ab56Error != nil {
		ab56T.Fatal(ab56Error)
	}
	if ab56Error := ab56JSON.Unmarshal(ab56Expected, &ab56Right); ab56Error != nil {
		ab56T.Fatal(ab56Error)
	}
	if !ab56Reflect.DeepEqual(ab56Left, ab56Right) {
		ab56T.Fatalf("ab56 row differs\ngot: %s\nwant: %s", ab56Actual, ab56Expected)
	}
}

func ab56ViewForTest(ab56T *ab56Testing.T) Ab56View {
	ab56T.Helper()
	ab56View, ab56Error := Ab56NewView(ab56TestContext())
	if ab56Error != nil {
		ab56T.Fatal(ab56Error)
	}
	return ab56View
}

func ab56Hit(ab56T *ab56Testing.T, ab56View Ab56View, ab56Key Ab56Key) ab56JSON.RawMessage {
	ab56T.Helper()
	ab56Raw, ab56OK, ab56Error := ab56View.Ab56Lookup(ab56Key)
	if ab56Error != nil || !ab56OK || len(ab56Raw) == 0 {
		ab56T.Fatalf("ab56 missing control: %v %v", ab56OK, ab56Error)
	}
	return ab56Raw
}

func ab56Miss(ab56T *ab56Testing.T, ab56View Ab56View, ab56Key Ab56Key) {
	ab56T.Helper()
	ab56Raw, ab56OK, ab56Error := ab56View.Ab56Lookup(ab56Key)
	if ab56Error != nil || ab56OK || ab56Raw != nil {
		ab56T.Fatalf("ab56 unexpected match: %s %v %v", ab56Raw, ab56OK, ab56Error)
	}
}

func TestAb56FiniteCandidates(ab56T *ab56Testing.T) {
	ab56Cases := []struct {
		ab56Group, ab56Instance            uint8
		ab56Register                       uint16
		ab56Codec                          string
		ab56Minimum, ab56Maximum, ab56Step float64
		ab56RawForm                        string
	}{
		{0, 0, 3, "UIN", 0, 12, 1, "00|00|0003|UIN|0|12|1"},
		{0, 0, 128, "EXP", -10, 10, 1, "00|00|0080|EXP|-10.0|10.0|1.0"},
		{1, 0, 1, "UIN", 0, 1, 1, "01|00|0001|UIN|0|1|1"},
		{1, 0, 3, "UIN", 0, 2, 1, "01|00|0003|UIN|0|2|1"},
		{2, 1, 2, "UIN", 0, 4, 1, "02|01|0002|UIN|0|4|1"},
		{2, 1, 4, "EXP", 15, 80, 1, "02|01|0004|EXP|15.0|80.0|1.0"},
		{2, 1, 5, "UCH", 0, 1, 1, "02|01|0005|UCH|0|1|1"},
		{2, 1, 6, "UCH", 0, 1, 1, "02|01|0006|UCH|0|1|1"},
		{2, 2, 2, "UIN", 0, 4, 1, "02|02|0002|UIN|0|4|1"},
		{2, 2, 4, "EXP", 15, 80, 1, "02|02|0004|EXP|15.0|80.0|1.0"},
		{2, 2, 5, "UCH", 0, 1, 1, "02|02|0005|UCH|0|1|1"},
		{2, 2, 6, "UCH", 0, 1, 1, "02|02|0006|UCH|0|1|1"},
		{2, 9, 2, "UIN", 0, 4, 1, "02|09|0002|UIN|0|4|1"},
		{2, 9, 4, "EXP", 15, 80, 1, "02|09|0004|EXP|15.0|80.0|1.0"},
		{2, 9, 5, "UCH", 0, 1, 1, "02|09|0005|UCH|0|1|1"},
		{2, 9, 6, "UCH", 0, 1, 1, "02|09|0006|UCH|0|1|1"},
		{3, 0, 1, "UIN", 0, 2, 1, "03|00|0001|UIN|0|2|1"},
		{3, 0, 2, "EXP", 15, 30, 0.5, "03|00|0002|EXP|15.0|30.0|0.5"},
		{3, 0, 5, "EXP", 5, 30, 1, "03|00|0005|EXP|5.0|30.0|1.0"},
		{3, 0, 6, "UIN", 0, 2, 1, "03|00|0006|UIN|0|2|1"},
	}
	ab56Entries, ab56Error := Ab56Entries()
	if ab56Error != nil {
		ab56T.Fatal(ab56Error)
	}
	var ab56Rows []ab56JSON.RawMessage
	if ab56Error := ab56JSON.Unmarshal(ab56Entries, &ab56Rows); ab56Error != nil {
		ab56T.Fatal(ab56Error)
	}
	if len(ab56Rows) != 20 {
		ab56T.Fatalf("ab56 count = %d", len(ab56Rows))
	}
	ab56View := ab56ViewForTest(ab56T)
	ab56Seen := make(map[string]bool)
	for _, ab56Case := range ab56Cases {
		ab56Key := ab56TestKey(ab56Case.ab56Group, ab56Case.ab56Instance, ab56Case.ab56Register)
		ab56Raw := ab56Hit(ab56T, ab56View, ab56Key)
		ab56Key.Ab56Family, ab56Key.Ab56Software, ab56Key.Ab56Hardware, ab56Key.Ab56APIVersion, ab56Key.Ab56APIRevision = "", "", "", "", ""
		ab56Expected := map[string]any{
			"key": ab56Key, "codec": ab56Case.ab56Codec, "width_bytes": nil, "width_grade": "Unknown",
			"minimum": ab56Case.ab56Minimum, "maximum": ab56Case.ab56Maximum, "step": ab56Case.ab56Step,
			"grade": "Hypothesis", "physical_applicability": "Unknown", "qualification": "Unknown", "admitted": false,
			"presence": "Unknown", "freshness": "Unknown", "unresolved": []string{"native identity applicability", "physical applicability", "field byte width", "slot presence"},
			"raw_form": ab56Case.ab56RawForm, "raw_fallback": "preserve native bytes; metadata arithmetic grants no execution or edit permission",
		}
		ab56Equal(ab56T, ab56Raw, ab56Encode(ab56T, ab56Expected))
		if ab56Seen[ab56Case.ab56RawForm] {
			ab56T.Fatal("ab56 duplicate expected selector")
		}
		ab56Seen[ab56Case.ab56RawForm] = true
	}
	ab56Miss(ab56T, ab56View, ab56TestKey(0, 0, 1))
	for _, ab56Key := range []Ab56Key{ab56TestKey(2, 0, 2), ab56TestKey(2, 3, 2), ab56TestKey(3, 0, 128)} {
		ab56Miss(ab56T, ab56View, ab56Key)
	}
}

func TestAb56EveryKeyAxis(ab56T *ab56Testing.T) {
	ab56View := ab56ViewForTest(ab56T)
	ab56Control := ab56TestKey(0, 0, 3)
	for ab56Index := 0; ab56Index < ab56Reflect.TypeOf(ab56Control).NumField(); ab56Index++ {
		ab56Name := ab56Reflect.TypeOf(ab56Control).Field(ab56Index).Name
		ab56T.Run("ab56_"+ab56Name, func(ab56T *ab56Testing.T) {
			for _, ab56Mode := range []string{"ab56_absent", "ab56_zero", "ab56_other"} {
				ab56Changed := ab56Control
				ab56Field := ab56Reflect.ValueOf(&ab56Changed).Elem().Field(ab56Index)
				if ab56Field.Kind() == ab56Reflect.String {
					if ab56Mode == "ab56_other" {
						ab56Field.SetString("ab56_other")
					} else {
						ab56Field.SetString("")
					}
				} else if ab56Mode == "ab56_absent" {
					ab56Field.SetZero()
				} else {
					ab56PointerValue := ab56Reflect.New(ab56Field.Type().Elem())
					if ab56Mode == "ab56_other" {
						ab56PointerValue.Elem().SetUint(255)
					}
					ab56Field.Set(ab56PointerValue)
				}
				if ab56Reflect.DeepEqual(ab56Changed, ab56Control) {
					continue
				}
				ab56Miss(ab56T, ab56View, ab56Changed)
				ab56Hit(ab56T, ab56View, ab56Control)
			}
		})
	}
	ab56OP06 := ab56Control
	ab56OP06.Ab56Operation, ab56OP06.Ab56Opcode = ab56Pointer(uint8(6)), ab56Pointer(uint8(6))
	ab56Miss(ab56T, ab56View, ab56OP06)
	ab56Miss(ab56T, Ab56View{}, ab56Control)
}

func TestAb56ContextAxes(ab56T *ab56Testing.T) {
	ab56Control := ab56TestContext()
	ab56Key := ab56TestKey(3, 0, 2)
	ab56Original := ab56ViewForTest(ab56T)
	for ab56Index := 0; ab56Index < ab56Reflect.TypeOf(ab56Control).NumField(); ab56Index++ {
		ab56Name := ab56Reflect.TypeOf(ab56Control).Field(ab56Index).Name
		ab56T.Run("ab56_"+ab56Name, func(ab56T *ab56Testing.T) {
			ab56Missing := ab56Control
			ab56Reflect.ValueOf(&ab56Missing).Elem().Field(ab56Index).SetString("")
			if _, ab56Error := Ab56NewView(ab56Missing); ab56Error == nil {
				ab56T.Fatal("ab56 incomplete context accepted")
			}
			ab56Changed := ab56Control
			ab56Reflect.ValueOf(&ab56Changed).Elem().Field(ab56Index).SetString("ab56_other")
			ab56View, ab56Error := Ab56NewView(ab56Changed)
			if ab56Error != nil {
				ab56T.Fatal(ab56Error)
			}
			ab56ChangedKey := ab56Key
			ab56Reflect.ValueOf(&ab56ChangedKey).Elem().FieldByName(ab56Name).SetString("ab56_other")
			ab56Hit(ab56T, ab56View, ab56ChangedKey)
			ab56Miss(ab56T, ab56View, ab56Key)
			ab56Miss(ab56T, ab56Original, ab56ChangedKey)
			ab56Hit(ab56T, ab56Original, ab56Key)
			ab56Reflect.ValueOf(&ab56Changed).Elem().Field(ab56Index).SetString("ab56_after")
			ab56Hit(ab56T, ab56View, ab56ChangedKey)
		})
	}
}

func TestAb56ArithmeticOnly(ab56T *ab56Testing.T) {
	ab56Row := ab56Decode(ab56T, ab56Hit(ab56T, ab56ViewForTest(ab56T), ab56TestKey(3, 0, 2)))
	var ab56Minimum, ab56Maximum, ab56Step float64
	for ab56Name, ab56Destination := range map[string]*float64{"minimum": &ab56Minimum, "maximum": &ab56Maximum, "step": &ab56Step} {
		if ab56Error := ab56JSON.Unmarshal(ab56Row[ab56Name], ab56Destination); ab56Error != nil {
			ab56T.Fatal(ab56Error)
		}
	}
	for _, ab56Case := range []struct {
		ab56Value float64
		ab56Want  bool
	}{{15, true}, {15.5, true}, {30, true}, {15.25, false}, {14.5, false}, {30.5, false}} {
		ab56Offset := (ab56Case.ab56Value - ab56Minimum) / ab56Step
		ab56Got := ab56Case.ab56Value >= ab56Minimum && ab56Case.ab56Value <= ab56Maximum && ab56Offset == float64(int(ab56Offset))
		if ab56Got != ab56Case.ab56Want {
			ab56T.Fatalf("ab56 arithmetic %v", ab56Case.ab56Value)
		}
	}
	if string(ab56Row["admitted"]) != "false" || string(ab56Row["presence"]) != `"Unknown"` || string(ab56Row["freshness"]) != `"Unknown"` {
		ab56T.Fatal("ab56 metadata promoted")
	}
}

func TestAb56StorageOwnership(ab56T *ab56Testing.T) {
	ab56Original, ab56Error := Ab56Entries()
	if ab56Error != nil {
		ab56T.Fatal(ab56Error)
	}
	ab56Source := append([]byte(nil), ab56Original...)
	ab56View, ab56Error := ab56NewView(ab56TestContext(), ab56Source)
	if ab56Error != nil {
		ab56T.Fatal(ab56Error)
	}
	ab56Key := ab56TestKey(3, 0, 2)
	ab56Before := ab56Hit(ab56T, ab56View, ab56Key)
	for ab56Index := range ab56Source {
		ab56Source[ab56Index] = 'x'
	}
	ab56Returned := ab56Hit(ab56T, ab56View, ab56Key)
	for ab56Index := range ab56Returned {
		ab56Returned[ab56Index] = 'x'
	}
	ab56Equal(ab56T, ab56Hit(ab56T, ab56View, ab56Key), ab56Before)
	for ab56Index := range ab56Original {
		ab56Original[ab56Index] = 'x'
	}
	ab56After, ab56Error := Ab56Entries()
	if ab56Error != nil || ab56Bytes.Equal(ab56Original, ab56After) {
		ab56T.Fatal("ab56 entries storage shared")
	}
	*ab56Key.Ab56Register = 99
	ab56Equal(ab56T, ab56Hit(ab56T, ab56View, ab56TestKey(3, 0, 2)), ab56Before)
}

func TestAb56OptionalAbsentZero(ab56T *ab56Testing.T) {
	ab56Control := ab56TestKey(0, 0, 3)
	ab56BaseRow := ab56Decode(ab56T, ab56Hit(ab56T, ab56ViewForTest(ab56T), ab56Control))
	for _, ab56Name := range []string{"Ab56Secondary", "Ab56Operation", "Ab56Opcode", "Ab56Group", "Ab56Instance", "Ab56Register"} {
		ab56T.Run("ab56_"+ab56Name, func(ab56T *ab56Testing.T) {
			for _, ab56Absent := range []bool{true, false} {
				ab56Stored := ab56Control
				ab56Stored.Ab56Family, ab56Stored.Ab56Software, ab56Stored.Ab56Hardware, ab56Stored.Ab56APIVersion, ab56Stored.Ab56APIRevision = "", "", "", "", ""
				ab56Field := ab56Reflect.ValueOf(&ab56Stored).Elem().FieldByName(ab56Name)
				ab56Other := ab56Control
				ab56OtherField := ab56Reflect.ValueOf(&ab56Other).Elem().FieldByName(ab56Name)
				if ab56Absent {
					ab56Field.SetZero()
					ab56OtherField.Set(ab56Reflect.New(ab56OtherField.Type().Elem()))
				} else {
					ab56Field.Set(ab56Reflect.New(ab56Field.Type().Elem()))
					ab56OtherField.SetZero()
				}
				ab56Row := ab56Decode(ab56T, ab56Encode(ab56T, ab56BaseRow))
				ab56Row["key"] = ab56Encode(ab56T, ab56Stored)
				ab56View, ab56Error := ab56NewView(ab56TestContext(), ab56Encode(ab56T, []any{ab56Row}))
				if ab56Error != nil {
					ab56T.Fatal(ab56Error)
				}
				ab56Match := ab56Control
				ab56Reflect.ValueOf(&ab56Match).Elem().FieldByName(ab56Name).Set(ab56Field)
				ab56Hit(ab56T, ab56View, ab56Match)
				ab56Miss(ab56T, ab56View, ab56Other)
			}
		})
	}
}

func TestAb56GradesWidthsAndRaw(ab56T *ab56Testing.T) {
	ab56Key := ab56TestKey(3, 0, 2)
	for _, ab56GradeValue := range []string{"Proven", "Hypothesis", "Unknown"} {
		for _, ab56Width := range []any{nil, 0, 1, 2} {
			ab56Row := ab56Decode(ab56T, ab56Hit(ab56T, ab56ViewForTest(ab56T), ab56Key))
			ab56Row["grade"] = ab56Encode(ab56T, ab56GradeValue)
			ab56Row["width_bytes"] = ab56Encode(ab56T, ab56Width)
			ab56Row["ab56_extension"] = ab56JSON.RawMessage(`{"ab56_raw":[0,255],"ab56_unresolved":"Unknown"}`)
			ab56Raw := ab56Encode(ab56T, []any{ab56Row})
			ab56View, ab56Error := ab56NewView(ab56TestContext(), ab56Raw)
			if ab56Error != nil {
				ab56T.Fatal(ab56Error)
			}
			ab56Got := ab56Hit(ab56T, ab56View, ab56Key)
			ab56Equal(ab56T, ab56Got, ab56Encode(ab56T, ab56Row))
			if !ab56Bytes.Contains(ab56Raw, ab56Got) {
				ab56T.Fatal("ab56 raw row rewritten")
			}
		}
	}
}

func TestAb56InvalidAndDuplicateData(ab56T *ab56Testing.T) {
	ab56Key := ab56TestKey(3, 0, 2)
	ab56Row := ab56Hit(ab56T, ab56ViewForTest(ab56T), ab56Key)
	ab56Invalid := [][]byte{nil, []byte(`null`), []byte(`{}`), []byte(`[null]`), []byte(`[5]`), []byte(`[{}]`), []byte(`[`), ab56Encode(ab56T, []ab56JSON.RawMessage{ab56Row, ab56Row})}
	for _, ab56Name := range []string{"codec", "width_bytes", "width_grade", "minimum", "maximum", "step", "grade", "physical_applicability", "qualification", "admitted", "presence", "freshness", "unresolved", "raw_form", "raw_fallback", "key"} {
		ab56Object := ab56Decode(ab56T, ab56Row)
		ab56Object[ab56Name] = ab56JSON.RawMessage(`{"ab56_bad":true}`)
		ab56Invalid = append(ab56Invalid, ab56Encode(ab56T, []any{ab56Object}))
		delete(ab56Object, ab56Name)
		ab56Invalid = append(ab56Invalid, ab56Encode(ab56T, []any{ab56Object}))
	}
	for ab56Name, ab56Value := range map[string]string{"admitted": "true", "presence": `"present"`, "freshness": `"fresh"`, "writable": "true", "step": "0"} {
		ab56Object := ab56Decode(ab56T, ab56Row)
		ab56Object[ab56Name] = ab56JSON.RawMessage(ab56Value)
		ab56Invalid = append(ab56Invalid, ab56Encode(ab56T, []any{ab56Object}))
	}
	ab56Invalid = append(ab56Invalid,
		ab56Bytes.Replace(ab56Encode(ab56T, []ab56JSON.RawMessage{ab56Row}), []byte(`"grade":"Hypothesis"`), []byte(`"grade":"Unknown","grade":"Hypothesis"`), 1),
		ab56Bytes.Replace(ab56Encode(ab56T, []ab56JSON.RawMessage{ab56Row}), []byte(`"primary":181`), []byte(`"primary":181,"primary":181`), 1),
	)
	for _, ab56Name := range []string{"family", "software", "hardware", "api_version", "api_revision", "ab56_unknown_axis"} {
		ab56Object := ab56Decode(ab56T, ab56Row)
		ab56ObjectKey := ab56Decode(ab56T, ab56Object["key"])
		ab56ObjectKey[ab56Name] = ab56JSON.RawMessage(`"ab56_forbidden"`)
		ab56Object["key"] = ab56Encode(ab56T, ab56ObjectKey)
		ab56Invalid = append(ab56Invalid, ab56Encode(ab56T, []any{ab56Object}))
	}
	ab56ExistingView := ab56ViewForTest(ab56T)
	ab56Saved := ab56Data
	defer func() { ab56Data = ab56Saved }()
	for ab56Index, ab56Source := range ab56Invalid {
		if _, ab56Error := ab56NewView(ab56TestContext(), ab56Source); ab56Error == nil {
			ab56T.Fatalf("ab56 malformed case %d accepted", ab56Index)
		}
		ab56Equal(ab56T, ab56Hit(ab56T, ab56ExistingView, ab56Key), ab56Row)
		ab56Data = string(ab56Source)
		if ab56Raw, ab56Error := Ab56Entries(); ab56Error == nil || ab56Raw != nil {
			ab56T.Fatalf("ab56 entries case %d accepted", ab56Index)
		}
		if _, ab56Error := Ab56NewView(ab56TestContext()); ab56Error == nil {
			ab56T.Fatalf("ab56 public constructor case %d accepted", ab56Index)
		}
	}
}

func TestAb56OpaqueIdentityStrings(ab56T *ab56Testing.T) {
	ab56Context := ab56TestContext()
	ab56Context.Ab56Family = "ab56_\xff"
	ab56View, ab56Error := Ab56NewView(ab56Context)
	if ab56Error != nil {
		ab56T.Fatal(ab56Error)
	}
	ab56Key := ab56TestKey(3, 0, 2)
	ab56Key.Ab56Family = "ab56_\xff"
	ab56Hit(ab56T, ab56View, ab56Key)
	ab56Key.Ab56Family = "ab56_\ufffd"
	ab56Miss(ab56T, ab56View, ab56Key)
	ab56Key.Ab56Family = "ab56_\xfe"
	ab56Miss(ab56T, ab56View, ab56Key)
}
