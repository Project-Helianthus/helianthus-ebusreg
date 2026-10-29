package vaillant

import (
	ab29bytes "bytes"
	ab29json "encoding/json"
	ab29fmt "fmt"
	ab29reflect "reflect"
	ab29testing "testing"
)

func ab29Decode(ab29T *ab29testing.T, ab29Data []byte) map[string]any {
	ab29T.Helper()
	var ab29Value map[string]any
	if ab29Error := ab29json.Unmarshal(ab29Data, &ab29Value); ab29Error != nil {
		ab29T.Fatal(ab29Error)
	}
	return ab29Value
}

func ab29Rows(ab29T *ab29testing.T) []any {
	ab29T.Helper()
	return ab29Decode(ab29T, Ab29RoomActions())["rows"].([]any)
}

func ab29Equal(ab29T *ab29testing.T, ab29Actual, ab29Expected any) {
	ab29T.Helper()
	if !ab29reflect.DeepEqual(ab29Actual, ab29Expected) {
		ab29T.Fatalf("got %#v; want %#v", ab29Actual, ab29Expected)
	}
}

func TestAb29Copy(ab29T *ab29testing.T) {
	ab29First, ab29Control := Ab29RoomActions(), Ab29RoomActions()
	ab29First[0] = '!'
	if !ab29bytes.Equal(Ab29RoomActions(), ab29Control) {
		ab29T.Fatal("accessor shares bytes")
	}
	ab29Control[len(ab29Control)-1] = '!'
	if !ab29json.Valid(Ab29RoomActions()) {
		ab29T.Fatal("second copy aliases catalog")
	}
}

func TestAb29Rows(ab29T *ab29testing.T) {
	ab29Doc := ab29Decode(ab29T, Ab29RoomActions())
	for ab29Key, ab29Expected := range map[string]any{
		"catalog": "ab29_room_actions", "representation_grade": "Proven", "candidate_grade": "Hypothesis", "physical_qualification": "Unknown",
		"coordinate_system": "logical positions include prefixes; not physical offsets",
	} {
		ab29T.Run("ab29_"+ab29Key, func(ab29T *ab29testing.T) { ab29Equal(ab29T, ab29Doc[ab29Key], ab29Expected) })
	}
	ab29RowsValue := ab29Rows(ab29T)
	ab29Equal(ab29T, len(ab29RowsValue), 18)
	ab29Cases := []struct {
		ab29Family, ab29Operation, ab29Variant, ab29Direction string
		ab29Extent                                            any
		ab29Names                                             []any
		ab29Limitations                                       []any
	}{
		{"B50A", "OP00", "room-action-temperature", "read", float64(3), []any{"room_action", "temperature"}, []any{"Unavailable sentinel is Unknown"}},
		{"B50C", "OP00", "keycode-digits", "read", float64(4), []any{"digit_1", "digit_2", "digit_3", "digit_4"}, []any{"Four independent numeric digits; not ASCII, BCD or u32"}},
		{"B50C", "OP01", "signed-value", "read", float64(2), []any{"value"}, []any{"Unit and sentinel are Unknown"}},
		{"B50C", "OP02", "circuit-names", "read", nil, []any{"circuit_name"}, []any{"Independent per-circuit values and failures; retain successful values on partial failure", "Bounded aggregate 1..15 is not installed cardinality", "Echoed circuit, string trimming, B524 namespace equivalence and wildcard semantics of zero are Unknown"}},
		{"B50D", "OP00", "ui-room-action", "UI-write", float64(1), []any{"room_action"}, []any{"UI-direction write is distinct from remote-control writes and completed UI effects"}},
		{"B50D", "OP01", "ui-holiday-status", "UI-write", float64(1), []any{"holiday_status"}, []any{"UI-direction write is distinct from remote-control writes and completed UI effects"}},
	}
	ab29IDs := []string{"ab29_b50a_op00", "ab29_b50c_op00", "ab29_b50c_op01", "ab29_b50c_op02", "ab29_b50d_op00", "ab29_b50d_op01"}
	for ab29Index, ab29Case := range ab29Cases {
		for ab29AlternativeIndex, ab29Alternative := range []string{"request", "positive", "negative"} {
			ab29Row := ab29RowsValue[ab29Index*3+ab29AlternativeIndex].(map[string]any)
			ab29Expected := map[string]any{
				"id": ab29IDs[ab29Index] + "_" + ab29Alternative, "family": ab29Case.ab29Family, "variant": ab29Case.ab29Variant, "direction": ab29Case.ab29Direction,
				"service": ab29Case.ab29Family, "operation": ab29Case.ab29Operation, "alternative": ab29Alternative, "grade": "Hypothesis", "physical_qualification": "Unknown",
				"logical_positions": nil, "extent_bytes": nil, "raw_domain": nil, "signedness": "Unknown", "byte_order": "Unknown", "conversion": nil, "unit": nil, "default": nil,
				"finite_mappings": map[string]any{}, "reserved_regions": []any{}, "relative_layout": []any{},
				"limitations":  append(append([]any{}, ab29Case.ab29Limitations...), "Absolute logical positions including any prefix are Unknown", "Unknown alternative layouts do not assert a qualified response format"),
				"raw_fallback": "preserve raw bytes and unresolved alternatives",
			}
			ab29Known := ab29Alternative == "positive" && ab29Case.ab29Direction == "read" || ab29Alternative == "request" && ab29Case.ab29Direction == "UI-write"
			if ab29Known {
				ab29Expected["extent_bytes"], ab29Expected["raw_domain"], ab29Expected["signedness"], ab29Expected["byte_order"], ab29Expected["relative_layout"] = ab29Case.ab29Extent, "field domains", "per field", "per field", ab29Case.ab29Names
			} else if ab29Alternative == "request" && ab29Case.ab29Direction == "read" && ab29Case.ab29Operation != "OP02" {
				ab29Expected["extent_bytes"], ab29Expected["raw_domain"] = float64(0), "no arguments"
			}
			ab29T.Run(ab29IDs[ab29Index]+"_"+ab29Alternative, func(ab29T *ab29testing.T) {
				for ab29Key, ab29Value := range ab29Expected {
					ab29T.Run("ab29_"+ab29Key, func(ab29T *ab29testing.T) {
						if _, ab29Present := ab29Row[ab29Key]; !ab29Present {
							ab29T.Fatal("missing key")
						}
						ab29Equal(ab29T, ab29Row[ab29Key], ab29Value)
					})
				}
				ab29Equal(ab29T, len(ab29Row), len(ab29Expected)+1)
				if !ab29Known {
					ab29Equal(ab29T, ab29Row["fields"], []any{})
				}
			})
		}
	}
}

func TestAb29Fields(ab29T *ab29testing.T) {
	ab29Expected := []struct {
		ab29Row    int
		ab29Fields string
	}{
		{1, `[
 {"name":"room_action","width_bytes":1,"raw_domain":[0,255],"signedness":"unsigned","byte_order":"not applicable","conversion":{"factor":1,"offset":0},"unit":null,"default":null,"finite_mappings":{"0":"no-action","1":"switch-on","2":"thermostat"},"reserved_regions":[{"range":[3,255],"meaning":"unlabeled; retain raw"}],"sentinel":"Unknown"},
 {"name":"temperature","width_bytes":2,"raw_domain":[-32768,32767],"signedness":"signed","byte_order":"little-endian","conversion":{"factor":0.0625,"offset":0,"candidate_range":[-2048,2047.9375]},"unit":"degrees C","default":null,"finite_mappings":{},"reserved_regions":[],"sentinel":"Unknown; neither 8000 nor FFFF is a qualified unavailable marker"}]`},
		{4, `[
 {"name":"digit_1","width_bytes":1,"raw_domain":[0,255],"signedness":"unsigned","byte_order":"not applicable","conversion":{"factor":1,"offset":0},"unit":null,"default":0,"finite_mappings":{},"reserved_regions":[{"range":[10,255],"meaning":"outside candidate digit domain; retain raw"}],"sentinel":"Unknown"},
 {"name":"digit_2","width_bytes":1,"raw_domain":[0,255],"signedness":"unsigned","byte_order":"not applicable","conversion":{"factor":1,"offset":0},"unit":null,"default":0,"finite_mappings":{},"reserved_regions":[{"range":[10,255],"meaning":"outside candidate digit domain; retain raw"}],"sentinel":"Unknown"},
 {"name":"digit_3","width_bytes":1,"raw_domain":[0,255],"signedness":"unsigned","byte_order":"not applicable","conversion":{"factor":1,"offset":0},"unit":null,"default":0,"finite_mappings":{},"reserved_regions":[{"range":[10,255],"meaning":"outside candidate digit domain; retain raw"}],"sentinel":"Unknown"},
 {"name":"digit_4","width_bytes":1,"raw_domain":[0,255],"signedness":"unsigned","byte_order":"not applicable","conversion":{"factor":1,"offset":0},"unit":null,"default":0,"finite_mappings":{},"reserved_regions":[{"range":[10,255],"meaning":"outside candidate digit domain; retain raw"}],"sentinel":"Unknown"}]`},
		{7, `[{"name":"value","width_bytes":2,"raw_domain":[-32768,32767],"signedness":"signed","byte_order":"little-endian","conversion":{"factor":1,"offset":0},"unit":null,"default":null,"finite_mappings":{},"reserved_regions":[],"sentinel":"Unknown"}]`},
		{10, `[{"name":"circuit_name","width_bytes":null,"raw_domain":null,"signedness":"Unknown","byte_order":"Unknown","conversion":null,"unit":null,"default":null,"finite_mappings":{},"reserved_regions":[],"sentinel":"Unknown"}]`},
		{12, `[{"name":"room_action","width_bytes":1,"raw_domain":[0,255],"signedness":"unsigned","byte_order":"not applicable","conversion":{"factor":1,"offset":0},"unit":null,"default":0,"finite_mappings":{"0":"no-action","1":"switch-on","2":"thermostat"},"reserved_regions":[{"range":[3,255],"meaning":"unlabeled; retain raw"}],"sentinel":"Unknown"}]`},
		{15, `[{"name":"holiday_status","width_bytes":1,"raw_domain":[0,255],"signedness":"unsigned","byte_order":"not applicable","conversion":{"factor":1,"offset":0},"unit":null,"default":0,"finite_mappings":{"0":"0","1":"1"},"reserved_regions":[{"range":[2,255],"meaning":"unlabeled; retain raw"}],"sentinel":"Unknown"}]`},
	}
	ab29RowsValue := ab29Rows(ab29T)
	for _, ab29Case := range ab29Expected {
		var ab29Fields []any
		if ab29Error := ab29json.Unmarshal([]byte(ab29Case.ab29Fields), &ab29Fields); ab29Error != nil {
			ab29T.Fatal(ab29Error)
		}
		ab29Actual := ab29RowsValue[ab29Case.ab29Row].(map[string]any)["fields"].([]any)
		ab29Equal(ab29T, len(ab29Actual), len(ab29Fields))
		for ab29Index, ab29Field := range ab29Fields {
			ab29Want := ab29Field.(map[string]any)
			ab29Want["grade"], ab29Want["raw_fallback"], ab29Want["logical_positions"] = "Hypothesis", "preserve raw bytes", nil
			ab29Got := ab29Actual[ab29Index].(map[string]any)
			ab29Equal(ab29T, len(ab29Got), len(ab29Want))
			for ab29Key, ab29Value := range ab29Want {
				ab29T.Run(ab29fmt.Sprintf("ab29_row_%d_field_%d_%s", ab29Case.ab29Row, ab29Index, ab29Key), func(ab29T *ab29testing.T) {
					if _, ab29Present := ab29Got[ab29Key]; !ab29Present {
						ab29T.Fatal("missing key")
					}
					ab29Equal(ab29T, ab29Got[ab29Key], ab29Value)
				})
			}
		}
	}
}

func ab29Candidate(ab29Field map[string]any, ab29Raw []byte) float64 {
	ab29Carrier := int(ab29Raw[0])
	if ab29Field["width_bytes"] == float64(2) {
		ab29Carrier += int(ab29Raw[1]) << 8
	}
	if ab29Field["signedness"] == "signed" && ab29Carrier >= 32768 {
		ab29Carrier -= 65536
	}
	ab29Conversion := ab29Field["conversion"].(map[string]any)
	return float64(ab29Carrier)*ab29Conversion["factor"].(float64) + ab29Conversion["offset"].(float64)
}

func TestAb29Synthetic(ab29T *ab29testing.T) {
	ab29RowsValue := ab29Rows(ab29T)
	ab29Room := ab29RowsValue[1].(map[string]any)["fields"].([]any)
	for _, ab29Case := range []struct {
		ab29Raw   []byte
		ab29Label string
		ab29Value float64
	}{
		{[]byte{2, 128, 255}, "thermostat", -8},
		{[]byte{3, 0, 0}, "", 0},
		{[]byte{0, 0, 128}, "no-action", -2048},
		{[]byte{0, 255, 255}, "no-action", -0.0625},
	} {
		ab29T.Run(ab29fmt.Sprintf("ab29_%v", ab29Case.ab29Raw), func(ab29T *ab29testing.T) {
			ab29Mapping := ab29Room[0].(map[string]any)["finite_mappings"].(map[string]any)
			ab29Label, _ := ab29Mapping[ab29fmt.Sprint(ab29Case.ab29Raw[0])].(string)
			ab29Equal(ab29T, ab29Label, ab29Case.ab29Label)
			ab29Equal(ab29T, ab29Candidate(ab29Room[1].(map[string]any), ab29Case.ab29Raw[1:]), ab29Case.ab29Value)
		})
	}
	ab29Digits := ab29RowsValue[4].(map[string]any)["fields"].([]any)
	ab29Equal(ab29T, len(ab29Digits), 4)
	for ab29Index, ab29Raw := range []byte{1, 2, 3, 4} {
		ab29Equal(ab29T, ab29Candidate(ab29Digits[ab29Index].(map[string]any), []byte{ab29Raw}), float64(ab29Index+1))
	}
	ab29Numeric := ab29RowsValue[7].(map[string]any)["fields"].([]any)[0].(map[string]any)
	ab29Equal(ab29T, ab29Candidate(ab29Numeric, []byte{254, 255}), float64(-2))
	ab29Equal(ab29T, ab29Numeric["unit"], nil)
	for _, ab29Index := range []int{12, 15} {
		ab29Mapping := ab29RowsValue[ab29Index].(map[string]any)["fields"].([]any)[0].(map[string]any)["finite_mappings"].(map[string]any)
		if _, ab29Present := ab29Mapping["255"]; ab29Present {
			ab29T.Fatal("unknown byte gained a label")
		}
	}
	// This catalog declares independent circuit outcomes; it supplies no aggregate decoder.
	ab29CircuitPolicy := ab29RowsValue[10].(map[string]any)["limitations"].([]any)
	ab29Equal(ab29T, ab29CircuitPolicy[0], "Independent per-circuit values and failures; retain successful values on partial failure")
}
