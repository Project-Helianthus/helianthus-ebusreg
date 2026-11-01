package vaillant

import (
	ab51Bytes "bytes"
	ab51Binary "encoding/binary"
	ab51JSON "encoding/json"
	ab51Reflect "reflect"
	ab51Strings "strings"
	ab51Testing "testing"
)

func ab51Ptr(ab51Value uint16) *uint16 { return &ab51Value }

func ab51Key(ab51Operation uint16, ab51Field string) Ab51Key {
	return Ab51Key{Ab51Family: "Vaillant", Ab51Variant: "Unknown", Ab51Primary: 181, Ab51Secondary: ab51Ptr(22), Ab51Operation: ab51Ptr(ab51Operation), Ab51Direction: "Unknown", Ab51Field: ab51Field, Ab51Alternative: "none", Ab51InstanceScope: "none"}
}

func ab51Object(ab51T *ab51Testing.T, ab51Raw []byte) map[string]any {
	ab51T.Helper()
	var ab51Value map[string]any
	if ab51Error := ab51JSON.Unmarshal(ab51Raw, &ab51Value); ab51Error != nil {
		ab51T.Fatal(ab51Error)
	}
	return ab51Value
}

func ab51Get(ab51T *ab51Testing.T, ab51Operation uint16, ab51Field string) map[string]any {
	ab51T.Helper()
	ab51Raw, ab51Found, ab51Error := Ab51Lookup(ab51Key(ab51Operation, ab51Field))
	if ab51Error != nil || !ab51Found {
		ab51T.Fatalf("lookup %s: %t %v", ab51Field, ab51Found, ab51Error)
	}
	return ab51Object(ab51T, ab51Raw)
}

func ab51Equal(ab51T *ab51Testing.T, ab51Actual, ab51Expected any) {
	ab51T.Helper()
	ab51Want, ab51Error := ab51JSON.Marshal(ab51Expected)
	if ab51Error != nil {
		ab51T.Fatal(ab51Error)
	}
	ab51Got, ab51Error := ab51JSON.Marshal(ab51Actual)
	if ab51Error != nil || !ab51Bytes.Equal(ab51Got, ab51Want) {
		ab51T.Fatalf("got %s want %s (%v)", ab51Got, ab51Want, ab51Error)
	}
}

func TestAb51FiniteMetadata(ab51T *ab51Testing.T) {
	ab51Source, ab51Error := Ab51Entries()
	if ab51Error != nil {
		ab51T.Fatal(ab51Error)
	}
	var ab51Rows []ab51JSON.RawMessage
	if ab51Error = ab51JSON.Unmarshal(ab51Source, &ab51Rows); ab51Error != nil || len(ab51Rows) != 20 {
		ab51T.Fatalf("finite array: %d %v", len(ab51Rows), ab51Error)
	}
	ab51Sequence := []string{"setpoint", "mode", "special_function", "message_1", "message_2", "message_3", "message_4", "message_5", "minute", "hour", "day", "month", "year"}
	for _, ab51Row := range ab51Rows {
		ab51Value := ab51Object(ab51T, ab51Row)
		var ab51Identity struct {
			Ab51Key Ab51Key `json:"key"`
		}
		if ab51Error = ab51JSON.Unmarshal(ab51Row, &ab51Identity); ab51Error != nil {
			ab51T.Fatal(ab51Error)
		}
		ab51LookupRow, ab51Found, ab51LookupError := Ab51Lookup(ab51Identity.Ab51Key)
		if !ab51Found || ab51LookupError != nil || !ab51Bytes.Equal(ab51LookupRow, ab51Row) {
			ab51T.Fatal("whole exact row lost", ab51LookupError)
		}
		for _, ab51Name := range []string{"physical_qualification", "presence", "freshness", "execution", "position_grade", "variant_association_grade"} {
			ab51Equal(ab51T, ab51Value[ab51Name], "Unknown")
		}
		ab51Equal(ab51T, ab51Value["grade"], "Hypothesis")
		ab51Equal(ab51T, ab51Value["metadata_presence"], "Proven")
		ab51Equal(ab51T, ab51Value["position"], nil)
		if len(ab51Value["unresolved"].([]any)) == 0 || len(ab51Value["excluded"].([]any)) != 3 {
			ab51T.Fatal("unresolved/exclusion boundary lost")
		}
		ab51Equal(ab51T, ab51Value["excluded"], []string{"Passive energy, OP14 aggregation, time-positive alternatives, energy-type flags, meter/clear enums and physical identity are separate.", "OP03 decimal software-version and identity-bearing structures have no supplied fields and are not instantiated.", "Voltage conversion is variant-specific and Unknown here."})
		if !ab51Strings.Contains(ab51Value["native_contract"].(string), "no native promotion") {
			ab51T.Fatal("native boundary lost")
		}
		ab51Width, ab51WidthGrade := any(nil), "Unknown"
		if ab51Identity.Ab51Key.Ab51Field == "setpoint" || ab51Strings.HasPrefix(ab51Identity.Ab51Key.Ab51Field, "message_") {
			ab51Width, ab51WidthGrade = float64(8), "Hypothesis"
		}
		if *ab51Identity.Ab51Key.Ab51Operation == 13 {
			ab51Width, ab51WidthGrade = float64(16), "Hypothesis"
		}
		ab51Equal(ab51T, ab51Value["width_bits"], ab51Width)
		ab51Equal(ab51T, ab51Value["width_grade"], ab51WidthGrade)
		ab51Unit := any(nil)
		if ab51Identity.Ab51Key.Ab51Field == "setpoint" {
			ab51Unit = "degrees Celsius"
		}
		ab51Equal(ab51T, ab51Value["unit"], ab51Unit)
	}
	for _, ab51Field := range ab51Sequence {
		ab51Value := ab51Get(ab51T, 8, ab51Field)
		ab51Equal(ab51T, ab51Value["sequence"], ab51Sequence)
		if ab51Strings.HasPrefix(ab51Field, "message_") {
			ab51Equal(ab51T, ab51Value["independent"], true)
		}
	}
	ab51Equal(ab51T, ab51Get(ab51T, 8, "setpoint")["raw_form"], map[string]any{"encoding": "unsigned", "divide_by": 2, "preserve": "original bytes"})
	for _, ab51Field := range []string{"mode", "special_function", "message_1", "message_2", "message_3", "message_4", "message_5"} {
		ab51Equal(ab51T, ab51Get(ab51T, 8, ab51Field)["raw_form"], map[string]any{"encoding": "opaque", "preserve": "original bytes"})
	}
	for _, ab51Field := range []string{"page", "data"} {
		ab51Value := ab51Get(ab51T, 14, ab51Field)
		ab51Equal(ab51T, ab51Value["maximum_length"], nil)
		ab51Equal(ab51T, ab51Value["raw_form"], map[string]any{"encoding": "opaque", "preserve": "original bytes"})
	}
	ab51Value := ab51Get(ab51T, 14, "data")
	ab51Equal(ab51T, ab51Value["variable"], true)
	ab51Equal(ab51T, ab51Value["empty_metadata"], true)
	ab51Equal(ab51T, ab51Value["peer_acceptance"], "Unknown")
	ab51Value = ab51Get(ab51T, 2, "unresolved_time_associations")
	ab51Equal(ab51T, ab51Value["association_grade"], "Unknown")
	ab51Equal(ab51T, ab51Value["conflicts"], map[string]any{
		"minute": "minute metadata on a second field conflicts with unprovided field associations",
		"day":    map[string]any{"default": 0, "calendar_domain": []int{1, 31}, "status": "out-of-domain calendar conflict", "association_grade": "Unknown"},
		"month":  map[string]any{"default": 0, "calendar_domain": []int{1, 12}, "status": "out-of-domain calendar conflict", "association_grade": "Unknown"},
	})
}

func ab51Time(ab51T *ab51Testing.T, ab51Metadata map[string]any, ab51Raw byte) (int, string) {
	ab51T.Helper()
	ab51Form := ab51Metadata["raw_form"].(map[string]any)
	ab51Equal(ab51T, ab51Form, map[string]any{"encoding": "packed-decimal", "selection_order": []string{"raw FF replacement", "decimal validation", "decimal conversion"}, "replacement_raw": 255, "invalid": "retain raw", "preserve": "original bytes"})
	if int(ab51Raw) == int(ab51Form["replacement_raw"].(float64)) {
		return 0, "replacement"
	}
	if ab51Raw>>4 > 9 || ab51Raw&15 > 9 {
		return int(ab51Raw), "raw"
	}
	return int(ab51Raw>>4)*10 + int(ab51Raw&15), "value"
}

func TestAb51SyntheticTime(ab51T *ab51Testing.T) {
	ab51Order := []string{"minute", "hour", "day", "month", "year"}
	ab51Vector := []byte{0x12, 0x23, 0x31, 0x12, 0x99}
	ab51Want := []int{12, 23, 31, 12, 99}
	for ab51Index, ab51Field := range ab51Order {
		ab51Metadata := ab51Get(ab51T, 8, ab51Field)
		ab51Equal(ab51T, ab51Metadata["time_order"], ab51Order)
		ab51Equal(ab51T, ab51Metadata["century"], nil)
		ab51Equal(ab51T, ab51Metadata["timezone"], nil)
		for _, ab51Case := range []struct {
			Ab51Raw   byte
			Ab51Want  int
			Ab51State string
		}{{0x12, 12, "value"}, {0xff, 0, "replacement"}, {0xfa, 250, "raw"}, {ab51Vector[ab51Index], ab51Want[ab51Index], "value"}} {
			ab51Got, ab51State := ab51Time(ab51T, ab51Metadata, ab51Case.Ab51Raw)
			if ab51Got != ab51Case.Ab51Want || ab51State != ab51Case.Ab51State {
				ab51T.Fatal("time rule changed", ab51Field, ab51Got, ab51State)
			}
		}
	}
	ab51Equal(ab51T, ab51Vector, []byte{0x12, 0x23, 0x31, 0x12, 0x99})
}

func TestAb51SyntheticWords(ab51T *ab51Testing.T) {
	ab51Names := []string{"word_1", "word_2", "word_3", "word_4"}
	ab51Raw := []byte{0, 0, 1, 0, 0, 1, 255, 255}
	ab51Want := []uint16{0, 1, 256, 65535}
	for ab51Index, ab51Name := range ab51Names {
		ab51Metadata := ab51Get(ab51T, 13, ab51Name)
		ab51Equal(ab51T, ab51Metadata["raw_form"], map[string]any{"encoding": "unsigned", "byte_order": "little-endian", "preserve": "original bytes"})
		ab51Equal(ab51T, ab51Metadata["independent"], true)
		ab51Equal(ab51T, ab51Metadata["label"], nil)
		ab51Equal(ab51T, ab51Metadata["sentinel"], nil)
		if ab51Binary.LittleEndian.Uint16(ab51Raw[2*ab51Index:]) != ab51Want[ab51Index] {
			ab51T.Fatal("independent opaque words")
		}
	}
	ab51Equal(ab51T, ab51Raw, []byte{0, 0, 1, 0, 0, 1, 255, 255})
}

func ab51Encode(ab51T *ab51Testing.T, ab51Value any) []byte {
	ab51T.Helper()
	ab51Raw, ab51Error := ab51JSON.Marshal(ab51Value)
	if ab51Error != nil {
		ab51T.Fatal(ab51Error)
	}
	return ab51Raw
}

func TestAb51EveryIdentityAxis(ab51T *ab51Testing.T) {
	ab51Control := ab51Key(8, "minute")
	for ab51Index := 0; ab51Index < ab51Reflect.TypeOf(ab51Control).NumField(); ab51Index++ {
		ab51Axis := ab51Reflect.TypeOf(ab51Control).Field(ab51Index)
		ab51T.Run("ab51_"+ab51Axis.Name, func(ab51T *ab51Testing.T) {
			ab51Candidates := []Ab51Key{}
			if ab51Axis.Type.Kind() == ab51Reflect.Pointer {
				for _, ab51Option := range []*uint16{nil, ab51Ptr(0), ab51Ptr(123)} {
					ab51Candidate := ab51Control
					ab51Reflect.ValueOf(&ab51Candidate).Elem().Field(ab51Index).Set(ab51Reflect.ValueOf(ab51Option))
					ab51Candidates = append(ab51Candidates, ab51Candidate)
				}
			} else {
				ab51Candidate := ab51Control
				ab51Value := ab51Reflect.ValueOf(&ab51Candidate).Elem().Field(ab51Index)
				if ab51Value.Kind() == ab51Reflect.String {
					ab51Value.SetString("ab51-other")
				} else {
					ab51Value.SetUint(0)
				}
				ab51Candidates = append(ab51Candidates, ab51Candidate)
			}
			for _, ab51Candidate := range ab51Candidates {
				_, ab51Found, ab51Error := Ab51Lookup(ab51Candidate)
				if ab51Error != nil || ab51Found != ab51Reflect.DeepEqual(ab51Control, ab51Candidate) {
					ab51T.Fatal("production axis ignored", ab51Axis.Name, ab51Error)
				}
				ab51Source := ab51Encode(ab51T, []any{map[string]any{"key": ab51Candidate, "ab51_fixture": "ab51-axis-control"}})
				for _, ab51Query := range append([]Ab51Key{ab51Control}, ab51Candidates...) {
					ab51Raw, ab51Found, ab51Error := ab51Lookup(ab51Source, ab51Query)
					if ab51Error != nil || ab51Found != ab51Reflect.DeepEqual(ab51Candidate, ab51Query) {
						ab51T.Fatal("seam axis ignored", ab51Axis.Name, ab51Error)
					}
					if ab51Found {
						ab51Equal(ab51T, ab51Object(ab51T, ab51Raw)["ab51_fixture"], "ab51-axis-control")
					}
				}
			}
		})
	}
}

func TestAb51ValidationAndCopies(ab51T *ab51Testing.T) {
	ab51Original, ab51Error := Ab51Entries()
	if ab51Error != nil {
		ab51T.Fatal(ab51Error)
	}
	ab51KeyValue := ab51Key(8, "minute")
	ab51OriginalRow, ab51Found, ab51Error := Ab51Lookup(ab51KeyValue)
	if !ab51Found || ab51Error != nil {
		ab51T.Fatal(ab51Error)
	}
	ab51Copy, ab51Error := Ab51Entries()
	if ab51Error != nil {
		ab51T.Fatal(ab51Error)
	}
	for ab51Index := range ab51Copy {
		ab51Copy[ab51Index] = 'x'
	}
	ab51Fresh, ab51Error := Ab51Entries()
	if ab51Error != nil || !ab51Bytes.Equal(ab51Original, ab51Fresh) {
		ab51T.Fatal("entries alias", ab51Error)
	}
	ab51Source := append([]byte(nil), ab51Original...)
	ab51Returned, ab51Found, ab51Error := ab51Lookup(ab51Source, ab51KeyValue)
	if !ab51Found || ab51Error != nil {
		ab51T.Fatal(ab51Error)
	}
	for ab51Index := range ab51Source {
		ab51Source[ab51Index] = 'x'
	}
	if !ab51Bytes.Equal(ab51Returned, ab51OriginalRow) {
		ab51T.Fatal("source alias")
	}
	for ab51Index := range ab51Returned {
		ab51Returned[ab51Index] = 'x'
	}
	ab51Returned, ab51Found, ab51Error = Ab51Lookup(ab51KeyValue)
	if !ab51Found || ab51Error != nil || !ab51Bytes.Equal(ab51Returned, ab51OriginalRow) {
		ab51T.Fatal("lookup alias", ab51Error)
	}
	ab51Returned[0] = 'x'
	ab51Returned, ab51Found, ab51Error = Ab51Lookup(ab51KeyValue)
	if !ab51Found || ab51Error != nil || !ab51Bytes.Equal(ab51Returned, ab51OriginalRow) {
		ab51T.Fatal("public lookup alias", ab51Error)
	}
	ab51Duplicate := ab51Encode(ab51T, []ab51JSON.RawMessage{ab51OriginalRow, ab51OriginalRow})
	for _, ab51Bad := range [][]byte{nil, []byte("{"), []byte("null"), []byte("{}"), []byte("[null]"), []byte("[{}]"), []byte("[1]"), []byte("[] []"), ab51Duplicate} {
		if _, _, ab51Error = ab51Validate(ab51Bad); ab51Error == nil {
			ab51T.Fatal("invalid source accepted", string(ab51Bad))
		}
		if ab51Raw, ab51Found, ab51Error := ab51Lookup(ab51Bad, Ab51Key{}); ab51Error == nil || ab51Found || ab51Raw != nil {
			ab51T.Fatal("invalid source became absence/success")
		}
	}
	for _, ab51Name := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
		ab51ObjectValue := ab51Object(ab51T, ab51OriginalRow)
		ab51Members := ab51ObjectValue["key"].(map[string]any)
		ab51Members[ab51Name] = []any{"ab51-invalid-type"}
		if _, _, ab51Error = ab51Validate(ab51Encode(ab51T, []any{ab51ObjectValue})); ab51Error == nil {
			ab51T.Fatal("invalid axis accepted", ab51Name)
		}
	}
	ab51ObjectValue := ab51Object(ab51T, ab51OriginalRow)
	ab51ObjectValue["key"].(map[string]any)["ab51_unknown_axis"] = 0
	if _, _, ab51Error = ab51Validate(ab51Encode(ab51T, []any{ab51ObjectValue})); ab51Error == nil {
		ab51T.Fatal("unknown namespace axis accepted")
	}
	ab51Missing := ab51KeyValue
	ab51Missing.Ab51Field = "ab51-absent"
	if ab51Raw, ab51Found, ab51Error := Ab51Lookup(ab51Missing); ab51Error != nil || ab51Found || ab51Raw != nil {
		ab51T.Fatal("absence contract")
	}
}

func TestAb51WidthAndGradeRemainRowMetadata(ab51T *ab51Testing.T) {
	ab51Raw, ab51Found, ab51Error := Ab51Lookup(ab51Key(13, "word_1"))
	if !ab51Found || ab51Error != nil {
		ab51T.Fatal(ab51Error)
	}
	for _, ab51Axis := range []string{"width_bits", "grade"} {
		ab51Value := ab51Object(ab51T, ab51Raw)
		if ab51Axis == "width_bits" {
			ab51Value[ab51Axis] = 8
		} else {
			ab51Value[ab51Axis] = "Unknown"
		}
		ab51Source := ab51Encode(ab51T, []any{ab51Value})
		ab51Got, ab51Found, ab51Error := ab51Lookup(ab51Source, ab51Key(13, "word_1"))
		if !ab51Found || ab51Error != nil {
			ab51T.Fatal(ab51Error)
		}
		ab51Equal(ab51T, ab51Object(ab51T, ab51Got), ab51Value)
		if _, _, ab51Error = ab51Validate(ab51Encode(ab51T, []any{ab51Object(ab51T, ab51Raw), ab51Value})); ab51Error == nil {
			ab51T.Fatal("metadata difference hid duplicate identity", ab51Axis)
		}
	}
}

func TestAb51OptionalOmissionAndRequiredAxes(ab51T *ab51Testing.T) {
	ab51Raw, ab51Found, ab51Error := Ab51Lookup(ab51Key(8, "minute"))
	if !ab51Found || ab51Error != nil {
		ab51T.Fatal(ab51Error)
	}
	for _, ab51Axis := range []string{"secondary", "operation", "group", "register"} {
		ab51Value := ab51Object(ab51T, ab51Raw)
		delete(ab51Value["key"].(map[string]any), ab51Axis)
		var ab51Identity struct {
			Ab51Key Ab51Key `json:"key"`
		}
		if ab51Error = ab51JSON.Unmarshal(ab51Encode(ab51T, ab51Value), &ab51Identity); ab51Error != nil {
			ab51T.Fatal(ab51Error)
		}
		ab51Source := ab51Encode(ab51T, []any{ab51Value})
		if _, ab51Found, ab51Error = ab51Lookup(ab51Source, ab51Identity.Ab51Key); !ab51Found || ab51Error != nil {
			ab51T.Fatal("omitted optional axis lost", ab51Axis, ab51Error)
		}
		ab51Null := ab51Object(ab51T, ab51Encode(ab51T, ab51Value))
		ab51Null["key"].(map[string]any)[ab51Axis] = nil
		if _, _, ab51Error = ab51Validate(ab51Encode(ab51T, []any{ab51Value, ab51Null})); ab51Error == nil {
			ab51T.Fatal("null/omitted duplicate hidden", ab51Axis)
		}
		ab51Zero := ab51Object(ab51T, ab51Encode(ab51T, ab51Value))
		ab51Zero["key"].(map[string]any)[ab51Axis] = 0
		if _, _, ab51Error = ab51Validate(ab51Encode(ab51T, []any{ab51Value, ab51Zero})); ab51Error != nil {
			ab51T.Fatal("zero collapsed into absence", ab51Axis, ab51Error)
		}
	}
	for _, ab51Axis := range []string{"family", "variant", "primary", "direction", "field", "alternative", "instance_scope"} {
		for _, ab51Mode := range []string{"ab51-omitted", "ab51-null", "ab51-empty"} {
			if ab51Axis == "primary" && ab51Mode == "ab51-empty" {
				continue
			}
			ab51Value := ab51Object(ab51T, ab51Raw)
			ab51Members := ab51Value["key"].(map[string]any)
			switch ab51Mode {
			case "ab51-omitted":
				delete(ab51Members, ab51Axis)
			case "ab51-null":
				ab51Members[ab51Axis] = nil
			case "ab51-empty":
				ab51Members[ab51Axis] = ""
			}
			if _, _, ab51Error = ab51Validate(ab51Encode(ab51T, []any{ab51Value})); ab51Error == nil {
				ab51T.Fatal("invalid required axis accepted", ab51Axis, ab51Mode)
			}
		}
	}
}
