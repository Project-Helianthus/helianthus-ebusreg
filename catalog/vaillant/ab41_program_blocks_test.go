package vaillant

import (
	ab41bytes "bytes"
	ab41binary "encoding/binary"
	ab41json "encoding/json"
	ab41fmt "fmt"
	ab41reflect "reflect"
	ab41testing "testing"
)

func ab41Read(ab41t *ab41testing.T) map[string]any {
	ab41t.Helper()
	var ab41catalog map[string]any
	if ab41err := ab41json.Unmarshal(Ab41ProgramBlocks(), &ab41catalog); ab41err != nil {
		ab41t.Fatal(ab41err)
	}
	return ab41catalog
}

func ab41Field(ab41name string, ab41position, ab41width, ab41maximum, ab41default any, ab41unit any, ab41order, ab41conversion string, ab41mapping map[string]any) map[string]any {
	ab41signed := "unsigned"
	if ab41width == nil {
		ab41signed = "Unknown"
	}
	return map[string]any{"name": ab41name, "position": ab41position, "width": ab41width, "raw_domain": map[string]any{"minimum": 0, "maximum": ab41maximum}, "signedness": ab41signed, "byte_order": ab41order, "conversion": ab41conversion, "unit": ab41unit, "default": ab41default, "mappings": ab41mapping, "reserved_regions": nil, "grade": "Hypothesis", "raw_fallback": "preserve raw bytes"}
}

func ab41Expected() map[string]any {
	ab41rows := []any{}
	for ab41op := 6; ab41op <= 14; ab41op++ {
		for _, ab41alternative := range []string{"request", "positive", "negative"} {
			ab41fields := []any{}
			var ab41extent any
			var ab41expression any
			var ab41prefix any
			ab41direction := "response"
			if ab41alternative != "negative" {
				ab41prefix = []any{map[string]any{"position": 0, "value": 181}, map[string]any{"position": 1, "value": 27}, map[string]any{"position": 2, "value": ab41op}}
			}
			if ab41alternative == "request" {
				ab41direction = "request"
				ab41extent = 3
				switch ab41op {
				case 8:
					ab41extent = 11
					ab41fields = append(ab41fields, ab41Field("block", 3, 4, 4294967295, 0, nil, "little-endian", "identity", map[string]any{}), ab41Field("segment", 7, 4, 4294967295, 0, nil, "little-endian", "identity", map[string]any{}))
				case 9:
					ab41extent = nil
					ab41expression = "5+k"
					ab41data := ab41Field("data", 5, nil, nil, nil, nil, "Unknown", "opaque", map[string]any{})
					ab41data["length"] = map[string]any{"minimum": 0, "maximum": nil}
					ab41fields = append(ab41fields, ab41Field("offset", 3, 2, 65535, 0, nil, "little-endian", "identity", map[string]any{}), ab41data)
				case 13:
					ab41extent = 4
					ab41fields = append(ab41fields, ab41Field("timeout", 3, 1, 255, 15, "minutes", "not-applicable", "identity", map[string]any{}))
				case 14:
					ab41extent = 4
					ab41fields = append(ab41fields, ab41Field("direction", 3, 1, 255, 240, nil, "not-applicable", "identity", map[string]any{"15": "boot-loader", "240": "application"}))
				}
			}
			if ab41alternative == "positive" {
				ab41extent = 4
				switch ab41op {
				case 6:
					ab41extent = 6
					for ab41position := 3; ab41position <= 5; ab41position++ {
						ab41fields = append(ab41fields, ab41Field(ab41fmt.Sprintf("raw_%d", ab41position), ab41position, 1, 255, 0, nil, "not-applicable", "identity", map[string]any{}))
					}
				case 8:
					ab41fields = append(ab41fields, ab41Field("feedback", 3, 1, 255, nil, nil, "not-applicable", "identity", map[string]any{"0": "rejected", "51": "block-unavailable", "85": "accepted", "170": "protected"}))
				case 11:
					ab41extent = 5
					ab41fields = append(ab41fields, ab41Field("CRC", 3, 2, 65535, nil, nil, "little-endian", "identity", map[string]any{}))
				case 12:
					ab41fields = append(ab41fields, ab41Field("state", 3, 1, 255, nil, nil, "not-applicable", "identity", map[string]any{"0": "hardware-init", "1": "wait-programming-request", "2": "check-application-CRC", "3": "preparation", "4": "active", "5": "application-running"}))
				default:
					ab41mapping := map[string]any{}
					for ab41value := 0; ab41value < 256; ab41value++ {
						ab41label := "rejected"
						if ab41value == 1 {
							ab41label = "accepted"
						}
						ab41mapping[ab41fmt.Sprint(ab41value)] = ab41label
					}
					ab41fields = append(ab41fields, ab41Field("feedback", 3, 1, 255, nil, nil, "not-applicable", "identity", ab41mapping))
				}
			}
			ab41rows = append(ab41rows, map[string]any{"id": ab41fmt.Sprintf("ab41_op%02x_%s", ab41op, ab41alternative), "family": "Vaillant programming", "variant": "B51B partial", "direction": ab41direction, "service": "B51B", "operation": ab41op, "alternative": ab41alternative, "logical_prefix": ab41prefix, "extent": ab41extent, "extent_expression": ab41expression, "fields": ab41fields, "reserved_regions": nil, "grade": "Hypothesis", "raw_fallback": "preserve complete native bytes", "mutative_declaration": ab41alternative == "request" && (ab41op == 7 || ab41op == 8 || ab41op == 9 || ab41op == 10 || ab41op == 13 || ab41op == 14)})
		}
	}
	ab41expected := map[string]any{"representation_grade": "Proven", "candidate_grade": "Hypothesis", "qualification": "Unknown", "coordinates": "logical including prefix; not physical offsets", "limitations": []string{"negative response positions and extents Unknown", "reserved regions Unknown", "full forwarding, segmentation and local CRC variants Unknown", "physical addresses, response adaptation and geometry Unknown", "peer CRC algorithm, coverage and persistence Unknown", "empty OP09 data does not establish peer acceptance", "commit feedback does not prove persistence", "FFFFFFFF block is not a sentinel", "OP08 position 2 is operation, not block", "OP0D zero has no special meaning", "no runtime registration, wire builder, automatic selection or write admission"}, "entries": ab41rows}
	ab41encoded, _ := ab41json.Marshal(ab41expected)
	var ab41normalized map[string]any
	_ = ab41json.Unmarshal(ab41encoded, &ab41normalized)
	return ab41normalized
}

func ab41Assert(ab41t *ab41testing.T, ab41path string, ab41got, ab41want any) {
	ab41t.Helper()
	if ab41map, ab41ok := ab41want.(map[string]any); ab41ok {
		ab41actual, ab41ok := ab41got.(map[string]any)
		if !ab41ok || len(ab41actual) != len(ab41map) {
			ab41t.Fatalf("%s object differs", ab41path)
		}
		for ab41key, ab41value := range ab41map {
			ab41t.Run(ab41key, func(ab41t *ab41testing.T) { ab41Assert(ab41t, ab41path+"/"+ab41key, ab41actual[ab41key], ab41value) })
		}
		return
	}
	if ab41slice, ab41ok := ab41want.([]any); ab41ok {
		ab41actual, ab41ok := ab41got.([]any)
		if !ab41ok || len(ab41actual) != len(ab41slice) {
			ab41t.Fatalf("%s array differs", ab41path)
		}
		for ab41index, ab41value := range ab41slice {
			ab41t.Run(ab41fmt.Sprint(ab41index), func(ab41t *ab41testing.T) { ab41Assert(ab41t, ab41path, ab41actual[ab41index], ab41value) })
		}
		return
	}
	if !ab41reflect.DeepEqual(ab41got, ab41want) {
		ab41t.Fatalf("%s got %#v want %#v", ab41path, ab41got, ab41want)
	}
}

func TestAb41CatalogContract(ab41t *ab41testing.T) {
	ab41Assert(ab41t, "ab41", ab41Read(ab41t), ab41Expected())
}

func TestAb41IndependentBytes(ab41t *ab41testing.T) {
	ab41first := Ab41ProgramBlocks()
	ab41control := append([]byte(nil), ab41first...)
	for ab41index := range ab41first {
		ab41first[ab41index] ^= 255
	}
	if !ab41bytes.Equal(Ab41ProgramBlocks(), ab41control) {
		ab41t.Fatal("ab41 accessor aliases embedded storage")
	}
}

func TestAb41IdentityAxes(ab41t *ab41testing.T) {
	for _, ab41axis := range []string{"family", "variant", "direction", "service", "operation", "alternative"} {
		ab41t.Run(ab41axis, func(ab41t *ab41testing.T) {
			ab41control := ab41Read(ab41t)["entries"].([]any)[0].(map[string]any)
			ab41encoded, _ := ab41json.Marshal(ab41control)
			var ab41changed map[string]any
			_ = ab41json.Unmarshal(ab41encoded, &ab41changed)
			if ab41axis == "operation" {
				ab41changed[ab41axis] = float64(7)
			} else {
				ab41changed[ab41axis] = "ab41_distinct"
			}
			ab41encoded, _ = ab41json.Marshal([]any{ab41control, ab41changed})
			var ab41roundtrip []map[string]any
			_ = ab41json.Unmarshal(ab41encoded, &ab41roundtrip)
			if len(ab41roundtrip) != 2 || ab41reflect.DeepEqual(ab41roundtrip[0], ab41roundtrip[1]) || !ab41reflect.DeepEqual(ab41roundtrip[0], ab41control) || !ab41reflect.DeepEqual(ab41roundtrip[1], ab41changed) {
				ab41t.Fatal("ab41 identity axis lost")
			}
		})
	}
}

func TestAb41SyntheticInterpretation(ab41t *ab41testing.T) {
	ab41rows := ab41Read(ab41t)["entries"].([]any)
	ab41block := []byte{181, 27, 8, 4, 3, 2, 1, 8, 7, 6, 5}
	ab41fields := ab41rows[6].(map[string]any)["fields"].([]any)
	for ab41index, ab41want := range []uint32{0x01020304, 0x05060708} {
		ab41field := ab41fields[ab41index].(map[string]any)
		ab41position := int(ab41field["position"].(float64))
		if ab41field["byte_order"] != "little-endian" || ab41field["width"] != float64(4) || ab41binary.LittleEndian.Uint32(ab41block[ab41position:]) != ab41want {
			ab41t.Fatal("ab41 block/segment or endian differs")
		}
	}
	ab41data := []byte{181, 27, 9, 0x34, 0x12, 0, 255}
	ab41offset := ab41rows[9].(map[string]any)["fields"].([]any)[0].(map[string]any)
	if ab41binary.LittleEndian.Uint16(ab41data[int(ab41offset["position"].(float64)):]) != 0x1234 || !ab41bytes.Equal(ab41data[5:], []byte{0, 255}) {
		ab41t.Fatal("ab41 opaque offset vector differs")
	}
	ab41timeout := ab41rows[21].(map[string]any)["fields"].([]any)[0].(map[string]any)
	if ab41timeout["unit"] != "minutes" || ab41timeout["default"] != float64(byte(15)) || ab41timeout["conversion"] != "identity" {
		ab41t.Fatal("ab41 minutes changed")
	}
	ab41state := ab41rows[19].(map[string]any)["fields"].([]any)[0].(map[string]any)["mappings"].(map[string]any)
	if _, ab41found := ab41state["6"]; ab41found {
		ab41t.Fatal("ab41 state 6 labeled")
	}
	ab41special := ab41rows[7].(map[string]any)["fields"].([]any)[0].(map[string]any)["mappings"].(map[string]any)
	ab41generic := ab41rows[4].(map[string]any)["fields"].([]any)[0].(map[string]any)["mappings"].(map[string]any)
	if ab41special["85"] != "accepted" || ab41special["170"] != "protected" || ab41generic["85"] != "rejected" || ab41generic["170"] != "rejected" {
		ab41t.Fatal("ab41 feedback tables conflated")
	}
	if ab41binary.LittleEndian.Uint16([]byte{0x34, 0x12}) != 0x1234 {
		ab41t.Fatal("ab41 CRC numeric control")
	}
}
