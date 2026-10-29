package vaillant

import (
	ab33json "encoding/json"
	ab33fmt "fmt"
	ab33reflect "reflect"
	ab33strconv "strconv"
	ab33testing "testing"
)

func ab33Decode(ab33T *ab33testing.T, ab33Bytes []byte) map[string]any {
	ab33T.Helper()
	var ab33Value map[string]any
	if ab33Error := ab33json.Unmarshal(ab33Bytes, &ab33Value); ab33Error != nil {
		ab33T.Fatal(ab33Error)
	}
	return ab33Value
}

func ab33JSON(ab33Text string) any {
	var ab33Value any
	if ab33Error := ab33json.Unmarshal([]byte(ab33Text), &ab33Value); ab33Error != nil {
		panic(ab33Error)
	}
	return ab33Value
}

func ab33ExpectedField(ab33Name string, ab33Position any, ab33Default any, ab33Mappings, ab33Bits, ab33Reserved, ab33Domain string) any {
	return map[string]any{
		"name": ab33Name, "logical_position": ab33Position, "width_bytes": 1,
		"raw_domain": ab33JSON(ab33Domain), "signedness": "unsigned",
		"byte_order": "not-applicable", "conversion": "identity", "unit": nil,
		"default": ab33Default, "mappings": ab33JSON(ab33Mappings),
		"bits": ab33JSON(ab33Bits), "reserved_bits": ab33JSON(ab33Reserved),
		"grade": "Hypothesis", "raw_fallback": "preserve-raw-byte",
	}
}

func ab33ExpectedRow(ab33ID, ab33Direction, ab33Operation, ab33Alternative string, ab33Extent any, ab33Fields ...any) any {
	if ab33Fields == nil {
		ab33Fields = []any{}
	}
	return map[string]any{
		"id": ab33ID, "family": "B513", "variant": nil,
		"direction": ab33Direction, "service": "B513", "operation": ab33Operation,
		"alternative": ab33Alternative, "logical_positions_include_prefix": true,
		"physical_offsets": nil, "extent_bytes": ab33Extent, "grade": "Hypothesis",
		"raw_fallback": "preserve-raw-bytes", "reserved_regions": []any{}, "fields": ab33Fields,
	}
}

func ab33Expected(ab33T *ab33testing.T) map[string]any {
	ab33T.Helper()
	ab33Value := map[string]any{
		"representation_grade": "Proven", "candidate_grade": "Hypothesis", "physical_qualification": "Unknown",
		"coordinate_system": "logical-including-prefix-not-physical-offsets",
		"limitations":       ab33JSON(`["physical-adaptation-Unknown","OP00-complement-Unknown","OP04-address-dictionaries-Unknown","unspecified-positions-and-extents-Unknown","variant-Unknown","no-runtime-registration","no-wire-builder","no-automatic-selection","no-native-codec-promotion","no-executable-capability"]`),
		"entries": []any{
			ab33ExpectedRow("ab33_op05_request", "request", "OP05", "request", nil,
				ab33ExpectedField("request_bits", nil, 8, `{}`, `{"0":{"default":false,"role":"request-data"},"1":{"default":false,"role":"request-data"},"2":{"default":false,"role":"request-data"},"3":{"default":true,"role":"request-data"},"4":{"default":false,"role":"request-data"},"5":{"default":false,"role":"request-data"}}`, `[6,7]`, `[0,255]`)),
			ab33ExpectedRow("ab33_op05_positive", "response", "OP05", "positive", nil,
				ab33ExpectedField("feedback", nil, nil, `{"1":"accepted"}`, `{}`, `[]`, `[1,1]`)),
			ab33ExpectedRow("ab33_op05_negative", "response", "OP05", "negative", nil,
				ab33ExpectedField("feedback", nil, nil, `{"0":"not-accepted-or-not-supported","2..255":"not-accepted-or-not-supported"}`, `{}`, `[]`, `[0,255]`)),
			ab33ExpectedRow("ab33_op06_request", "request", "OP06", "request", nil,
				ab33ExpectedField("light", 3, 0, `{"0":"off","100":"on","255":"replacement"}`, `{}`, `[]`, `[0,255]`),
				ab33ExpectedField("colour", 4, 0, `{}`, `{}`, `[]`, `[0,255]`)),
			ab33ExpectedRow("ab33_op06_positive", "response", "OP06", "positive", nil,
				ab33ExpectedField("feedback", 3, nil, `{"1":"accepted"}`, `{}`, `[]`, `[1,1]`)),
			ab33ExpectedRow("ab33_op06_negative", "response", "OP06", "negative", nil,
				ab33ExpectedField("feedback", 3, nil, `{"0":"rejected","2..255":"rejected"}`, `{}`, `[]`, `[0,255]`)),
			ab33ExpectedRow("ab33_op07_request", "request", "OP07", "request", 0),
			ab33ExpectedRow("ab33_op07_positive", "response", "OP07", "positive", nil,
				ab33ExpectedField("appliance_type", 3, 0, `{"0":"modulating","1":"single-step","2":"two-step"}`, `{}`, `[]`, `[0,255]`),
				ab33ExpectedField("appliance_flags", 4, 0, `{}`, `{"0":{"label":"floor-standing","default":false},"1":{"label":"condensing","default":false},"2":{"label":"fan-assisted","default":false},"3":{"label":"open-vented","default":false},"4":{"label":"high-modulation-range","default":false},"5":{"label":"water-heater","default":false},"6":{"label":"step-modulation","default":false},"7":{"label":"ECS","default":false}}`, `[]`, `[0,255]`),
				ab33ExpectedField("fuel_flags", 5, 0, `{}`, `{"0":{"label":"gas","default":false},"1":{"label":"oil","default":false},"2":{"label":"ABC","default":false}}`, `[3,4,5,6,7]`, `[0,255]`),
				ab33ExpectedField("tank_flags", 6, 0, `{}`, `{"0":{"label":"layer-tank","default":false},"1":{"label":"warm-start","default":false},"2":{"label":"external-tank","default":false},"3":{"label":"combination","default":false}}`, `[4,5,6,7]`, `[0,255]`),
				ab33ExpectedField("feature_flags", 7, 0, `{}`, `{"0":{"label":"solar-compact","default":false},"1":{"label":"basic-solar","default":false},"2":{"label":"AAPS","default":false},"3":{"label":"advanced-power","default":false},"4":{"label":"partial-load","default":false},"5":{"label":"parameter-primary","default":false}}`, `[6,7]`, `[0,255]`),
				ab33ExpectedField("pump_mode", 8, 0, `{"0":false,"1":true}`, `{}`, `[]`, `[0,255]`),
				ab33ExpectedField("opaque_9", 9, 0, `{}`, `{}`, `[]`, `[0,255]`),
				ab33ExpectedField("opaque_10", 10, 0, `{}`, `{}`, `[]`, `[0,255]`),
				ab33ExpectedField("opaque_11", 11, 0, `{}`, `{}`, `[]`, `[0,255]`)),
		},
	}
	ab33Bytes, ab33Error := ab33json.Marshal(ab33Value)
	if ab33Error != nil {
		ab33T.Fatal(ab33Error)
	}
	return ab33Decode(ab33T, ab33Bytes)
}

func ab33Assert(ab33T *ab33testing.T, ab33Got, ab33Want any) {
	ab33T.Helper()
	switch ab33ExpectedValue := ab33Want.(type) {
	case map[string]any:
		ab33ActualValue, ab33OK := ab33Got.(map[string]any)
		if !ab33OK || len(ab33ActualValue) != len(ab33ExpectedValue) {
			ab33T.Fatalf("ab33 object keys: got %v, want %v", ab33Got, ab33Want)
		}
		for ab33Key, ab33Value := range ab33ExpectedValue {
			ab33T.Run("ab33_"+ab33Key, func(ab33T *ab33testing.T) {
				ab33Actual, ab33Present := ab33ActualValue[ab33Key]
				if !ab33Present {
					ab33T.Fatal("ab33 missing explicit property")
				}
				ab33Assert(ab33T, ab33Actual, ab33Value)
			})
		}
	case []any:
		ab33ActualValue, ab33OK := ab33Got.([]any)
		if !ab33OK || len(ab33ActualValue) != len(ab33ExpectedValue) {
			ab33T.Fatalf("ab33 array extent: got %v, want %v", ab33Got, ab33Want)
		}
		for ab33Index, ab33Value := range ab33ExpectedValue {
			ab33T.Run(ab33fmt.Sprintf("ab33_%d", ab33Index), func(ab33T *ab33testing.T) {
				ab33Assert(ab33T, ab33ActualValue[ab33Index], ab33Value)
			})
		}
	default:
		if !ab33reflect.DeepEqual(ab33Got, ab33Want) {
			ab33T.Fatalf("ab33 got %#v, want %#v", ab33Got, ab33Want)
		}
	}
}

func TestAb33Contract(ab33T *ab33testing.T) {
	ab33Assert(ab33T, ab33Decode(ab33T, Ab33ApplianceInfoCatalog()), ab33Expected(ab33T))
}

func TestAb33IndependentBytes(ab33T *ab33testing.T) {
	ab33First := Ab33ApplianceInfoCatalog()
	ab33Control := append([]byte(nil), ab33First...)
	for ab33Index := range ab33First {
		ab33First[ab33Index] = 0
	}
	if !ab33reflect.DeepEqual([]byte(Ab33ApplianceInfoCatalog()), ab33Control) {
		ab33T.Fatal("ab33 accessor shares storage")
	}
}

func TestAb33IdentityAxes(ab33T *ab33testing.T) {
	for ab33Index := 0; ab33Index < 8; ab33Index++ {
		for _, ab33Axis := range []string{"id", "family", "variant", "direction", "service", "operation", "alternative"} {
			ab33T.Run(ab33fmt.Sprintf("ab33_%d_%s", ab33Index, ab33Axis), func(ab33T *ab33testing.T) {
				ab33Catalog := ab33Decode(ab33T, Ab33ApplianceInfoCatalog())
				ab33Rows := ab33Catalog["entries"].([]any)
				ab33Row := ab33Rows[ab33Index].(map[string]any)
				ab33Row[ab33Axis] = "ab33_changed"
				ab33Bytes, ab33Error := ab33json.Marshal(ab33Catalog)
				if ab33Error != nil {
					ab33T.Fatal(ab33Error)
				}
				ab33RoundTrip := ab33Decode(ab33T, ab33Bytes)
				if !ab33reflect.DeepEqual(ab33RoundTrip, ab33Catalog) || ab33reflect.DeepEqual(ab33RoundTrip, ab33Expected(ab33T)) {
					ab33T.Fatal("ab33 identity axis lost or comparison did not detect its mutation")
				}
				if !ab33reflect.DeepEqual(ab33Decode(ab33T, Ab33ApplianceInfoCatalog()), ab33Expected(ab33T)) {
					ab33T.Fatal("ab33 local identity mutation escaped into catalog")
				}
			})
		}
	}
}

func ab33Fields(ab33T *ab33testing.T, ab33Index int) []any {
	ab33T.Helper()
	return ab33Decode(ab33T, Ab33ApplianceInfoCatalog())["entries"].([]any)[ab33Index].(map[string]any)["fields"].([]any)
}

func ab33Label(ab33Field map[string]any, ab33Raw int) any {
	ab33Mappings := ab33Field["mappings"].(map[string]any)
	if ab33Label, ab33Present := ab33Mappings[ab33strconv.Itoa(ab33Raw)]; ab33Present {
		return ab33Label
	}
	if ab33Raw >= 2 && ab33Raw <= 255 {
		return ab33Mappings["2..255"]
	}
	return nil
}

func TestAb33FeedbackAndRequestBits(ab33T *ab33testing.T) {
	ab33Request := ab33Fields(ab33T, 0)[0].(map[string]any)
	ab33Bits := ab33Request["bits"].(map[string]any)
	ab33Packed := 0
	for ab33Bit := 0; ab33Bit < 8; ab33Bit++ {
		ab33Definition, ab33Present := ab33Bits[ab33strconv.Itoa(ab33Bit)]
		if ab33Present != (ab33Bit < 6) {
			ab33T.Fatalf("ab33 bit %d request/reserved boundary", ab33Bit)
		}
		if ab33Present {
			ab33Default := ab33Definition.(map[string]any)["default"].(bool)
			if ab33Default != (ab33Bit == 3) {
				ab33T.Fatalf("ab33 bit %d default", ab33Bit)
			}
			if ab33Default {
				ab33Packed |= 1 << ab33Bit
			}
		}
	}
	if ab33Packed != 0x08 || ab33Request["default"] != float64(ab33Packed) {
		ab33T.Fatal("ab33 request packed default is not 08")
	}
	for _, ab33Case := range []struct {
		ab33Index int
		ab33Want  string
	}{{2, "not-accepted-or-not-supported"}, {5, "rejected"}} {
		for ab33Raw := 0; ab33Raw <= 255; ab33Raw++ {
			ab33Want := any(ab33Case.ab33Want)
			if ab33Raw == 1 {
				ab33Want = nil
			}
			if ab33Label(ab33Fields(ab33T, ab33Case.ab33Index)[0].(map[string]any), ab33Raw) != ab33Want {
				ab33T.Fatalf("ab33 feedback row %d raw %d", ab33Case.ab33Index, ab33Raw)
			}
		}
	}
}

func TestAb33InformationWalkingBitsAndRaw(ab33T *ab33testing.T) {
	ab33FieldsValue := ab33Fields(ab33T, 7)
	for ab33Index, ab33Labels := range [][]string{
		{"floor-standing", "condensing", "fan-assisted", "open-vented", "high-modulation-range", "water-heater", "step-modulation", "ECS"},
		{"gas", "oil", "ABC"},
		{"layer-tank", "warm-start", "external-tank", "combination"},
		{"solar-compact", "basic-solar", "AAPS", "advanced-power", "partial-load", "parameter-primary"},
	} {
		ab33Field := ab33FieldsValue[ab33Index+1].(map[string]any)
		ab33Bits := ab33Field["bits"].(map[string]any)
		for ab33Bit := 0; ab33Bit < 8; ab33Bit++ {
			ab33Raw := 1 << ab33Bit
			ab33Active := []string{}
			for ab33Key, ab33Definition := range ab33Bits {
				ab33Position, ab33Error := ab33strconv.Atoi(ab33Key)
				if ab33Error != nil {
					ab33T.Fatal(ab33Error)
				}
				if ab33Raw&(1<<ab33Position) != 0 {
					ab33Active = append(ab33Active, ab33Definition.(map[string]any)["label"].(string))
				}
			}
			ab33Want := []string{}
			if ab33Bit < len(ab33Labels) {
				ab33Want = append(ab33Want, ab33Labels[ab33Bit])
			}
			if !ab33reflect.DeepEqual(ab33Active, ab33Want) {
				ab33T.Fatalf("ab33 flag byte %d walking bit %d: %v", ab33Index+4, ab33Bit, ab33Active)
			}
		}
	}
	if ab33Label(ab33FieldsValue[0].(map[string]any), 3) != nil {
		ab33T.Fatal("ab33 type 3 acquired a fabricated label")
	}
	for ab33Raw, ab33Want := range []any{false, true, nil} {
		if ab33Label(ab33FieldsValue[5].(map[string]any), ab33Raw) != ab33Want {
			ab33T.Fatalf("ab33 Boolean raw %d", ab33Raw)
		}
	}
	ab33Logical := []byte{0, 0, 0, 3, 0, 0, 0, 0, 2, 7, 8, 9}
	ab33Tail := []byte{}
	for _, ab33Value := range ab33FieldsValue[6:] {
		ab33Field := ab33Value.(map[string]any)
		if len(ab33Field["mappings"].(map[string]any)) != 0 || ab33Field["raw_fallback"] != "preserve-raw-byte" {
			ab33T.Fatal("ab33 opaque field was interpreted")
		}
		ab33Tail = append(ab33Tail, ab33Logical[int(ab33Field["logical_position"].(float64))])
	}
	if !ab33reflect.DeepEqual(ab33Tail, []byte{7, 8, 9}) {
		ab33T.Fatalf("ab33 opaque tail lost: %v", ab33Tail)
	}
}
