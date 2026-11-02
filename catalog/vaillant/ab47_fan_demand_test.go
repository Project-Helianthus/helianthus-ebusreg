package vaillant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func ab47Catalog(ab47T *testing.T) map[string]any {
	ab47T.Helper()
	var ab47Value map[string]any
	if ab47Error := json.Unmarshal(Ab47FanDemandCatalog(), &ab47Value); ab47Error != nil {
		ab47T.Fatal(ab47Error)
	}
	return ab47Value
}

func ab47Field(ab47Name any, ab47Width int, ab47Max int, ab47Default any, ab47Mappings map[string]any, ab47Conversion string, ab47Unit any, ab47Bits any, ab47Reserved []any) map[string]any {
	if ab47Mappings == nil {
		ab47Mappings = map[string]any{}
	}
	if ab47Reserved == nil {
		ab47Reserved = []any{}
	}
	ab47Order := "not_applicable"
	if ab47Width > 8 {
		ab47Order = "Unknown"
	}
	return map[string]any{"name": ab47Name, "logical_position": nil, "width_bits": ab47Width, "raw_domain": []int{0, ab47Max}, "signedness": "unsigned", "byte_order": ab47Order, "conversion": ab47Conversion, "unit": ab47Unit, "default": ab47Default, "mappings": ab47Mappings, "bits": ab47Bits, "reserved": ab47Reserved, "grade": "Hypothesis", "raw_fallback": "preserve original bytes including reserved bits"}
}

func ab47ExpectedRows() []any {
	ab47Rows := []any{}
	ab47Add := func(ab47ID string, ab47Direction any, ab47Alternative string, ab47Extent any, ab47Fields []any) {
		ab47Rows = append(ab47Rows, map[string]any{"id": ab47ID, "family": "vaillant", "variant": "B522", "direction": ab47Direction, "service": "B522", "operation": nil, "alternative": ab47Alternative, "logical_extent_bytes": ab47Extent, "logical_coordinates": "include prefix positions; not physical offsets", "fields": ab47Fields, "grade": "Hypothesis", "raw_fallback": "preserve full original carrier", "physical_mapping": "Unknown", "installed_support": "Unknown", "execution": "Unknown", "semantic_equivalence": "Unknown"})
	}
	ab47Add("ab47_current_mode", nil, "current-mode", nil, []any{ab47Field("current_mode", 8, 255, 0, map[string]any{
		"0": []any{"no heating demand to source/circuits"},
		"1": []any{"no source demand but demand from a sink", nil},
		"2": []any{"source heating demand"},
		"3": []any{"no source cooling demand but demand from a sink"},
		"4": []any{"source cooling demand"},
	}, "identity", nil, nil, nil)})
	ab47Add("ab47_request", "request", "request", 7, []any{
		ab47Field("fan", 8, 255, 0, nil, "identity", nil, nil, nil),
		ab47Field("automatic", 1, 1, 0, map[string]any{"0": []any{"disabled"}, "1": []any{"enabled"}}, "identity", nil, []int{0, 0}, []any{map[string]any{"bits": []int{1, 7}, "expected": 0}}),
		ab47Field("room_temperature", 8, 255, 0, nil, "raw * 0.5", "degrees C", nil, nil),
		ab47Field("mode", 8, 255, 0, nil, "identity", nil, nil, nil),
	})
	ab47Positive := []any{}
	for ab47Bit := 0; ab47Bit < 3; ab47Bit++ {
		ab47Positive = append(ab47Positive, ab47Field(nil, 1, 1, false, map[string]any{"0": []any{false}, "1": []any{true}}, "identity", nil, []int{ab47Bit, ab47Bit}, nil))
	}
	ab47Add("ab47_positive", "response", "positive", 4, ab47Positive)
	ab47Rows[2].(map[string]any)["reserved"] = []any{map[string]any{"bits": []int{3, 7}, "expected": 0}}
	ab47Add("ab47_negative", "response", "negative", nil, []any{})
	ab47Unset := ab47Field(nil, 16, 65535, nil, map[string]any{"FF FF": []any{"candidate unset"}}, "Unknown", nil, nil, nil)
	ab47Unset["width_bits"] = nil
	ab47Unset["raw_domain"] = nil
	ab47Unset["signedness"] = "Unknown"
	ab47Add("ab47_unset", nil, "FF FF unset", nil, []any{ab47Unset})
	return ab47Rows
}

func ab47Normalized(ab47T *testing.T, ab47Value any) any {
	ab47T.Helper()
	ab47Bytes, ab47Error := json.Marshal(ab47Value)
	if ab47Error != nil {
		ab47T.Fatal(ab47Error)
	}
	var ab47Result any
	if ab47Error = json.Unmarshal(ab47Bytes, &ab47Result); ab47Error != nil {
		ab47T.Fatal(ab47Error)
	}
	return ab47Result
}

func TestAb47ExactCatalog(ab47T *testing.T) {
	ab47Got := ab47Catalog(ab47T)
	ab47Want := ab47Normalized(ab47T, ab47ExpectedRows()).([]any)
	ab47Rows, ab47OK := ab47Got["entries"].([]any)
	if !ab47OK || len(ab47Rows) != 5 {
		ab47T.Fatal("ab47 expected exactly five separate rows")
	}
	for ab47Index, ab47Row := range ab47Want {
		for ab47Key, ab47Value := range ab47Row.(map[string]any) {
			ab47T.Run(fmt.Sprintf("ab47_row_%d/%s", ab47Index, ab47Key), func(ab47T *testing.T) {
				if !reflect.DeepEqual(ab47Rows[ab47Index].(map[string]any)[ab47Key], ab47Value) {
					ab47T.Fatalf("ab47 %s: got %#v want %#v", ab47Key, ab47Rows[ab47Index].(map[string]any)[ab47Key], ab47Value)
				}
			})
		}
		if !reflect.DeepEqual(ab47Rows[ab47Index], ab47Row) {
			ab47T.Errorf("ab47 row %d differs from finite contract", ab47Index)
		}
	}
	if ab47Got["software_representation"] != "Proven" || ab47Got["candidate_semantics"] != "Hypothesis" || ab47Got["physical_qualification"] != "Unknown" {
		ab47T.Fatal("ab47 grades conflated")
	}
	ab47Limitations := []string{
		"Partial field catalog; absent names and coordinates remain unknown.",
		"Raw current modes 5..255 and all mode values remain unlabeled.",
		"Current-mode raw 1 has two unresolved labels; the competing label is Unknown.",
		"Heating priority during simultaneous cooling and bypass behavior are candidate semantics only.",
		"FF is not a sentinel for room temperature, mode or fan.",
		"Separate FF FF unset interpretation does not promote individual FF carriers.",
		"No runtime registration, wire builder, automatic selection, native codec promotion or executable capability.",
	}
	if !reflect.DeepEqual(ab47Got["limitations"], ab47Normalized(ab47T, ab47Limitations)) {
		ab47T.Fatal("ab47 limitations changed")
	}
}

func TestAb47IdentityAxes(ab47T *testing.T) {
	ab47Rows := ab47Catalog(ab47T)["entries"].([]any)
	for ab47Index, ab47Value := range ab47Rows {
		for _, ab47Axis := range []string{"family", "variant", "direction", "service", "operation", "alternative"} {
			ab47T.Run(fmt.Sprintf("ab47_row_%d/%s", ab47Index, ab47Axis), func(ab47T *testing.T) {
				ab47Original := ab47Value.(map[string]any)
				ab47Changed := ab47Normalized(ab47T, ab47Original).(map[string]any)
				ab47Changed[ab47Axis] = "ab47_distinct"
				ab47RoundTrip := ab47Normalized(ab47T, []any{ab47Original, ab47Changed}).([]any)
				if !reflect.DeepEqual(ab47RoundTrip[0], ab47Original) || reflect.DeepEqual(ab47RoundTrip[0], ab47RoundTrip[1]) || ab47RoundTrip[1].(map[string]any)[ab47Axis] != "ab47_distinct" {
					ab47T.Fatal("ab47 identity axis lost")
				}
				ab47Changed[ab47Axis] = ab47Original[ab47Axis]
				if !reflect.DeepEqual(ab47Changed, ab47Original) {
					ab47T.Fatal("ab47 invalid control")
				}
			})
		}
	}
}

func TestAb47RawCarriers(ab47T *testing.T) {
	ab47Rows := ab47Catalog(ab47T)["entries"].([]any)
	ab47Fields := ab47Rows[1].(map[string]any)["fields"].([]any)
	ab47Room := ab47Fields[2].(map[string]any)
	if ab47Room["conversion"] != "raw * 0.5" {
		ab47T.Fatal("ab47 temperature conversion changed")
	}
	for ab47Index, ab47Raw := range []float64{0, 1, 255} {
		if ab47Raw*0.5 != []float64{0, 0.5, 127.5}[ab47Index] {
			ab47T.Fatal("ab47 temperature examples")
		}
	}
	ab47Fan := ab47Fields[0].(map[string]any)
	if ab47Fan["raw_domain"].([]any)[1] != float64(255) || ab47Fan["conversion"] != "identity" || len(ab47Fan["mappings"].(map[string]any)) != 0 {
		ab47T.Fatal("ab47 fan255 lost")
	}
	ab47Modes := ab47Rows[0].(map[string]any)["fields"].([]any)[0].(map[string]any)["mappings"].(map[string]any)
	if len(ab47Modes["1"].([]any)) != 2 || ab47Modes["1"].([]any)[1] != nil || ab47Modes["5"] != nil {
		ab47T.Fatal("ab47 unresolved mode alternatives lost")
	}
	ab47Automatic := ab47Fields[1].(map[string]any)
	ab47Positive := ab47Rows[2].(map[string]any)
	if !reflect.DeepEqual(ab47Automatic["bits"], []any{float64(0), float64(0)}) || !reflect.DeepEqual(ab47Automatic["reserved"], ab47Normalized(ab47T, []any{map[string]any{"bits": []int{1, 7}, "expected": 0}})) || len(ab47Positive["fields"].([]any)) != 3 || !reflect.DeepEqual(ab47Positive["reserved"], ab47Normalized(ab47T, []any{map[string]any{"bits": []int{3, 7}, "expected": 0}})) {
		ab47T.Fatal("ab47 bit roles changed")
	}
	for _, ab47Raw := range []byte{1, 128} {
		ab47Enabled := ab47Raw&1 != 0
		ab47Discrepancy := ab47Raw&254 != 0
		if ab47Enabled != (ab47Raw == 1) || ab47Discrepancy != (ab47Raw == 128) {
			ab47T.Fatal("ab47 automatic examples")
		}
	}
	if byte(8)&7 != 0 || byte(8)&248 != 8 {
		ab47T.Fatal("ab47 positive08 must remain reserved")
	}
}

func TestAb47CopyIsolation(ab47T *testing.T) {
	ab47Before := bytes.Clone(Ab47FanDemandCatalog())
	ab47Returned := Ab47FanDemandCatalog()
	if len(ab47Returned) == 0 {
		ab47T.Fatal("ab47 empty catalog")
	}
	for ab47Index := range ab47Returned {
		ab47Returned[ab47Index] ^= 255
	}
	if !bytes.Equal(Ab47FanDemandCatalog(), ab47Before) {
		ab47T.Fatal("ab47 returned bytes alias catalog")
	}
}
