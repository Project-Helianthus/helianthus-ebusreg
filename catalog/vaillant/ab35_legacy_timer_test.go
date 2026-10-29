package vaillant

import (
	ab35bytes "bytes"
	ab35json "encoding/json"
	ab35fmt "fmt"
	ab35reflect "reflect"
	ab35testing "testing"
)

func ab35Decode(ab35T *ab35testing.T, ab35Data []byte) map[string]any {
	ab35T.Helper()
	var ab35Value map[string]any
	if ab35Error := ab35json.Unmarshal(ab35Data, &ab35Value); ab35Error != nil {
		ab35T.Fatal(ab35Error)
	}
	return ab35Value
}

func ab35Expected(ab35T *ab35testing.T) map[string]any {
	ab35T.Helper()
	ab35Want := ab35Decode(ab35T, []byte(`{
 "id":"ab35_legacy_timer","representation_grade":"Proven","candidate_grade":"Hypothesis","physical_qualification":"Unknown",
 "opt_in":true,"runtime_registration":false,"wire_builder":false,"automatic_selection":false,"native_codec_promotion":false,"write_admission":false,
 "coordinates":"logical coordinates include prefix positions; they are not physical offsets",
 "limitations":{"physical_forms":"Unknown","competing_two_position_gap":"Unknown","mode_labels":"Unknown","enforcement":"Unknown","deletion_reset":"Unknown","atomicity":"Unknown","overlap":"Unknown","midnight_crossing":"Unknown","timezone_DST":"Unknown","physical_schedule_effects":"Unknown"},
 "ordering_candidate":{"grade":"Hypothesis","written_start_equal_end":"disabled interval","interval_order":["first","second","third"],"out_of_order_qualification":"Unknown"}
 }`))
	ab35Entries := []any{}
	for ab35Index, ab35Identity := range []struct {
		ab35Operation, ab35Alternative, ab35Direction string
		ab35Extent                                    any
	}{
		{"read", "request", "request", float64(4)}, {"read", "positive", "response", float64(11)}, {"read", "negative", "response", nil},
		{"write", "request", "request", float64(11)}, {"write", "positive", "response", float64(3)}, {"write", "negative", "response", nil},
	} {
		ab35Row := ab35Decode(ab35T, []byte(`{"family":"Vaillant legacy timer","variant":"periods","service":"B515","logical_positions":null,"physical_extent":null,"width":null,"raw_domain":null,"signedness":"Unknown","byte_order":"Unknown","conversion":null,"unit":null,"default":null,"finite_mappings":null,"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve complete native bytes and independent useful fields"}`))
		ab35Row["id"] = "ab35_" + ab35Identity.ab35Operation + "_" + ab35Identity.ab35Alternative
		ab35Row["operation"], ab35Row["alternative"], ab35Row["direction"], ab35Row["logical_extent"] = ab35Identity.ab35Operation, ab35Identity.ab35Alternative, ab35Identity.ab35Direction, ab35Identity.ab35Extent
		ab35Fields := map[string]any{}
		if ab35Index != 2 && ab35Index != 4 && ab35Index != 5 {
			for _, ab35Name := range []string{"weekday", "channel", "mode"} {
				ab35Field := ab35Decode(ab35T, []byte(`{"logical_positions":null,"width":null,"raw_domain":[0,255],"signedness":"Unknown","byte_order":"Unknown","conversion":null,"unit":null,"default":null,"finite_mappings":null,"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve raw carrier"}`))
				if ab35Name == "weekday" {
					ab35Field["default"] = "Monday"
				}
				if ab35Name == "channel" {
					ab35Field["default"] = float64(0)
					if ab35Index == 0 || ab35Index == 3 {
						ab35Field["default"] = map[string]any{"circulation": float64(2), "DHW": float64(1), "mixer_one": float64(0), "mixer_two": float64(0)}
					}
				}
				ab35Fields[ab35Name] = ab35Field
			}
		}
		if ab35Index == 1 || ab35Index == 3 {
			ab35Period := ab35Decode(ab35T, []byte(`{"logical_positions":null,"width":null,"raw_domain":null,"signedness":"Unknown","byte_order":"Unknown","conversion":"raw * 10 minutes for raw 0..144 only","unit":"minute","default":null,"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve raw label; 145 and 255 have Unknown roles","undeclared_time_conversion":"09:05","rounding":"Unknown"}`))
			ab35Mappings := map[string]any{}
			for ab35Minutes := 0; ab35Minutes <= 1440; ab35Minutes += 10 {
				ab35Mappings[ab35fmt.Sprint(ab35Minutes/10)] = ab35fmt.Sprintf("%02d:%02d", ab35Minutes/60, ab35Minutes%60)
			}
			ab35Period["finite_mappings"] = ab35Mappings
			ab35Fields["period"] = ab35Period
		}
		ab35Row["fields"] = ab35Fields
		ab35Entries = append(ab35Entries, ab35Row)
	}
	for ab35Index := 1; ab35Index <= 5; ab35Index++ {
		ab35Row := ab35Decode(ab35T, []byte(`{"family":"Vaillant legacy timer","direction":null,"service":"B515","operation":null,"logical_positions":null,"logical_extent":null,"physical_extent":null,"width":null,"raw_domain":[0,3],"signedness":"Unknown","byte_order":"Unknown","conversion":null,"unit":null,"default":null,"finite_mappings":null,"reserved_regions":null,"grade":"Hypothesis","raw_fallback":"preserve raw mode and unresolved variant","mapping_injective":false,"inverse_default":null,"fields":{}}`))
		ab35Row["id"] = ab35fmt.Sprintf("ab35_mode_candidate_%d", ab35Index)
		ab35Row["variant"] = ab35fmt.Sprintf("mode_candidate_%d", ab35Index)
		ab35Row["alternative"] = ab35Row["variant"]
		ab35Entries = append(ab35Entries, ab35Row)
	}
	ab35Want["entries"] = ab35Entries
	return ab35Want
}

func ab35Assert(ab35T *ab35testing.T, ab35Path string, ab35Got, ab35Want any) {
	ab35T.Helper()
	switch ab35Value := ab35Want.(type) {
	case map[string]any:
		ab35Actual, ab35OK := ab35Got.(map[string]any)
		if !ab35OK || len(ab35Actual) != len(ab35Value) {
			ab35T.Fatalf("%s: object keys differ", ab35Path)
		}
		for ab35Key, ab35Child := range ab35Value {
			ab35ActualChild, ab35Present := ab35Actual[ab35Key]
			if !ab35Present {
				ab35T.Fatalf("%s/%s: missing explicit property", ab35Path, ab35Key)
			}
			ab35Assert(ab35T, ab35Path+"/"+ab35Key, ab35ActualChild, ab35Child)
		}
	case []any:
		ab35Actual, ab35OK := ab35Got.([]any)
		if !ab35OK || len(ab35Actual) != len(ab35Value) {
			ab35T.Fatalf("%s: alternative/domain count differs", ab35Path)
		}
		for ab35Index, ab35Child := range ab35Value {
			ab35Assert(ab35T, ab35fmt.Sprintf("%s/%d", ab35Path, ab35Index), ab35Actual[ab35Index], ab35Child)
		}
	default:
		if !ab35reflect.DeepEqual(ab35Got, ab35Want) {
			ab35T.Errorf("%s: got %v, want %v", ab35Path, ab35Got, ab35Want)
		}
	}
}

func TestAb35CatalogContract(ab35T *ab35testing.T) {
	ab35Assert(ab35T, "ab35", ab35Decode(ab35T, Ab35LegacyTimerCatalog()), ab35Expected(ab35T))
}

func TestAb35CopyOwnership(ab35T *ab35testing.T) {
	ab35First, ab35Second := Ab35LegacyTimerCatalog(), Ab35LegacyTimerCatalog()
	if len(ab35First) == 0 || !ab35bytes.Equal(ab35First, ab35Second) {
		ab35T.Fatal("ab35 inconsistent catalog")
	}
	ab35Snapshot := ab35bytes.Clone(ab35Second)
	for ab35Index := range ab35First {
		ab35First[ab35Index] ^= 255
	}
	if !ab35bytes.Equal(ab35Snapshot, ab35Second) || !ab35bytes.Equal(ab35Snapshot, Ab35LegacyTimerCatalog()) {
		ab35T.Fatal("ab35 caller mutation changed catalog")
	}
}

func TestAb35PeriodVocabulary(ab35T *ab35testing.T) {
	ab35Rows := ab35Decode(ab35T, Ab35LegacyTimerCatalog())["entries"].([]any)
	for _, ab35Index := range []int{1, 3} {
		ab35Period := ab35Rows[ab35Index].(map[string]any)["fields"].(map[string]any)["period"].(map[string]any)
		ab35Mappings := ab35Period["finite_mappings"].(map[string]any)
		if len(ab35Mappings) != 145 {
			ab35T.Fatal("ab35 requires 145 labels")
		}
		for ab35Raw, ab35Want := range map[string]string{"0": "00:00", "54": "09:00", "144": "24:00"} {
			if ab35Mappings[ab35Raw] != ab35Want {
				ab35T.Fatalf("ab35 %s: %v", ab35Raw, ab35Mappings[ab35Raw])
			}
		}
		for _, ab35Raw := range []string{"145", "255"} {
			if _, ab35Present := ab35Mappings[ab35Raw]; ab35Present {
				ab35T.Fatal("ab35 unknown raw role assigned")
			}
		}
		for _, ab35Label := range ab35Mappings {
			if ab35Label == "09:05" {
				ab35T.Fatal("ab35 undeclared rounding")
			}
		}
		if ab35Period["rounding"] != "Unknown" {
			ab35T.Fatal("ab35 rounding promoted")
		}
	}
}

func TestAb35IdentityAxes(ab35T *ab35testing.T) {
	ab35Rows := ab35Decode(ab35T, Ab35LegacyTimerCatalog())["entries"].([]any)
	ab35Want := ab35Expected(ab35T)["entries"].([]any)
	for ab35Index, ab35Entry := range ab35Rows {
		for _, ab35Axis := range []string{"id", "family", "variant", "direction", "service", "operation", "alternative"} {
			ab35T.Run(ab35fmt.Sprintf("ab35_%d_%s", ab35Index, ab35Axis), func(ab35T *ab35testing.T) {
				ab35Row := ab35Entry.(map[string]any)
				if !ab35reflect.DeepEqual(ab35Row[ab35Axis], ab35Want[ab35Index].(map[string]any)[ab35Axis]) {
					ab35T.Fatal("ab35 identity changed")
				}
				ab35Pair := []any{ab35Row, ab35Row}
				ab35Wire, ab35Error := ab35json.Marshal(ab35Pair)
				if ab35Error != nil {
					ab35T.Fatal(ab35Error)
				}
				var ab35RoundTrip []map[string]any
				if ab35Error = ab35json.Unmarshal(ab35Wire, &ab35RoundTrip); ab35Error != nil {
					ab35T.Fatal(ab35Error)
				}
				ab35RoundTrip[1][ab35Axis] = "ab35_distinct"
				ab35Wire, ab35Error = ab35json.Marshal(ab35RoundTrip)
				if ab35Error != nil {
					ab35T.Fatal(ab35Error)
				}
				if ab35Error = ab35json.Unmarshal(ab35Wire, &ab35RoundTrip); ab35Error != nil {
					ab35T.Fatal(ab35Error)
				}
				if ab35reflect.DeepEqual(ab35RoundTrip[0], ab35RoundTrip[1]) || ab35RoundTrip[1][ab35Axis] != "ab35_distinct" {
					ab35T.Fatal("ab35 identity axis collapsed")
				}
				delete(ab35RoundTrip[1], ab35Axis)
				if ab35reflect.DeepEqual(ab35RoundTrip[0], ab35RoundTrip[1]) {
					ab35T.Fatal("ab35 absent axis equated with explicit null/value")
				}
			})
		}
	}
}

func ab35Invert(ab35Mapping map[int]string) (map[string]int, bool) {
	ab35Inverse := map[string]int{}
	for ab35Raw, ab35Label := range ab35Mapping {
		if _, ab35Present := ab35Inverse[ab35Label]; ab35Present {
			return nil, false
		}
		ab35Inverse[ab35Label] = ab35Raw
	}
	return ab35Inverse, true
}

func TestAb35ModeAmbiguity(ab35T *ab35testing.T) {
	ab35Control := map[int]string{0: "ab35_a", 1: "ab35_b", 2: "ab35_c", 3: "ab35_d"}
	ab35Inverse, ab35OK := ab35Invert(ab35Control)
	if !ab35OK || len(ab35Inverse) != 4 {
		ab35T.Fatal("ab35 injective control rejected")
	}
	for ab35Raw, ab35Label := range ab35Control {
		if ab35Inverse[ab35Label] != ab35Raw {
			ab35T.Fatal("ab35 control inverse changed")
		}
	}
	ab35Inverse, ab35OK = ab35Invert(map[int]string{0: "ab35_same", 1: "ab35_same", 2: "ab35_same", 3: "ab35_same"})
	if ab35OK || ab35Inverse != nil {
		ab35T.Fatal("ab35 ambiguous inversion accepted")
	}
	ab35Rows := ab35Decode(ab35T, Ab35LegacyTimerCatalog())["entries"].([]any)
	for _, ab35Entry := range ab35Rows[6:] {
		ab35Row := ab35Entry.(map[string]any)
		if ab35Row["mapping_injective"] != false || ab35Row["finite_mappings"] != nil || ab35Row["inverse_default"] != nil {
			ab35T.Fatal("ab35 candidate mode table invented or invertible")
		}
	}
}

func TestAb35OrderingCandidate(ab35T *ab35testing.T) {
	ab35Order := ab35Decode(ab35T, Ab35LegacyTimerCatalog())["ordering_candidate"].(map[string]any)
	if ab35Order["written_start_equal_end"] != "disabled interval" || ab35Order["out_of_order_qualification"] != "Unknown" {
		ab35T.Fatal("ab35 ordering promoted")
	}
	for _, ab35Example := range []struct {
		ab35Name           string
		ab35Start, ab35End int
		ab35Want           string
	}{
		{"ab35_equal", 54, 54, "disabled interval"},
		{"ab35_unequal", 54, 60, "Unknown"},
	} {
		ab35T.Run(ab35Example.ab35Name, func(ab35T *ab35testing.T) {
			ab35Claim := "Unknown"
			if ab35Example.ab35Start == ab35Example.ab35End {
				ab35Claim = ab35Order["written_start_equal_end"].(string)
			}
			if ab35Claim != ab35Example.ab35Want {
				ab35T.Fatal("ab35 equality candidate changed")
			}
		})
	}
	for ab35Index, ab35Intervals := range [][3][2]int{
		{{90, 96}, {54, 60}, {108, 114}},
		{{54, 60}, {108, 114}, {90, 96}},
	} {
		ab35T.Run(ab35fmt.Sprintf("ab35_out_of_order_%d", ab35Index), func(ab35T *ab35testing.T) {
			if ab35Intervals[0][0] < ab35Intervals[1][0] && ab35Intervals[1][0] < ab35Intervals[2][0] {
				ab35T.Fatal("ab35 invalid out-of-order fixture")
			}
			if ab35Order["out_of_order_qualification"] != "Unknown" {
				ab35T.Fatal("ab35 synthetic ordering qualified")
			}
		})
	}
	ab35Assert(ab35T, "ab35_order", ab35Order["interval_order"], []any{"first", "second", "third"})
}
