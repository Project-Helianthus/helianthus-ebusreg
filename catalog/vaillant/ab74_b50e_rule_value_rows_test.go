package vaillant

import (
	ab74bytes "bytes"
	ab74json "encoding/json"
	ab74fmt "fmt"
	ab74reflect "reflect"
	ab74testing "testing"
)

func ab74Pointer[ab74Type any](ab74Input ab74Type) *ab74Type { return &ab74Input }

func ab74Key(ab74Selector byte) Ab74Key {
	return Ab74Key{ab74Pointer("vaillant"), ab74Pointer("B50E"), ab74Pointer(byte(0xB5)), ab74Pointer(byte(0x0E)), ab74Pointer(byte(2)), ab74Pointer("positive"), ab74Pointer(ab74Selector), ab74Pointer("rule_value_row")}
}

func ab74ExpectedField(ab74Offset int, ab74Width *int, ab74Encoding string) Ab74Field {
	return Ab74Field{Ab74Offset: ab74Pointer(ab74Offset), Ab74Width: ab74Width, Ab74Encoding: ab74Encoding, Ab74ByteOrder: "Unknown", Ab74Enumeration: map[string]string{}}
}

func ab74Expected() []Ab74Entry {
	var ab74Rows []Ab74Entry
	for ab74Index, ab74Case := range []struct {
		ab74Encoding                                       string
		ab74Width, ab74CompanionOffset, ab74CompanionWidth *int
		ab74CompanionEncoding                              string
	}{
		{"opaque", nil, nil, nil, ""},
		{"boolean_u8", ab74Pointer(1), nil, nil, ""},
		{"u8", ab74Pointer(1), nil, nil, ""},
		{"i8", ab74Pointer(1), nil, nil, ""},
		{"u16", ab74Pointer(2), nil, nil, ""},
		{"i16", ab74Pointer(2), nil, nil, ""},
		{"u32", ab74Pointer(4), nil, nil, ""},
		{"i32", ab74Pointer(4), nil, nil, ""},
		{"secure_boolean_u8", ab74Pointer(1), nil, nil, ""},
		{"u8", ab74Pointer(1), ab74Pointer(1), ab74Pointer(1), "raw_inverse"},
		{"i8", ab74Pointer(1), ab74Pointer(1), ab74Pointer(1), "raw_inverse"},
		{"u16", ab74Pointer(2), ab74Pointer(2), ab74Pointer(2), "raw_inverse"},
		{"i16", ab74Pointer(2), ab74Pointer(2), ab74Pointer(2), "raw_inverse"},
		{"u32", ab74Pointer(4), ab74Pointer(4), ab74Pointer(4), "raw_inverse"},
		{"i32", ab74Pointer(4), ab74Pointer(4), ab74Pointer(4), "raw_inverse"},
		{"opaque", nil, nil, nil, ""},
		{"opaque", nil, nil, nil, ""},
		{"i16", ab74Pointer(2), ab74Pointer(2), ab74Pointer(1), "status_u8"},
		{"i16", ab74Pointer(2), ab74Pointer(2), ab74Pointer(2), "raw_inverse"},
	} {
		ab74Row := Ab74Entry{Ab74Key: ab74Key(byte(ab74Index)), Ab74Value: ab74ExpectedField(0, ab74Case.ab74Width, ab74Case.ab74Encoding), Ab74Base: "row-relative", Ab74Grade: "Hypothesis", Ab74MetadataPresence: "Proven", Ab74PhysicalSupport: "Unknown", Ab74Effects: "Unknown", Ab74Persistence: "Unknown", Ab74Delivery: "Unknown", Ab74DefaultQualification: "Unknown", Ab74RawFallback: "Preserve isolated row bytes and independently useful fields after partial failure; retain prior good values with uncertain freshness. Inverse bytes are opaque companions, without complement validation. No numeric decoding without byte-order qualification. Opaque end-of-data alternatives have unknown width, never zero width. OP03 buffer bytes are outside this table."}
		if ab74Case.ab74CompanionOffset != nil {
			ab74Companion := ab74ExpectedField(*ab74Case.ab74CompanionOffset, ab74Case.ab74CompanionWidth, ab74Case.ab74CompanionEncoding)
			ab74Row.Ab74Companion = &ab74Companion
		}
		if ab74Index == 8 {
			ab74Row.Ab74Value.Ab74Enumeration = map[string]string{"0F": "true", "F0": "false"}
		}
		if ab74Index >= 17 {
			ab74Row.Ab74Value.Ab74Scale = ab74Pointer(1.0 / 16)
			ab74Row.Ab74Value.Ab74Unit = ab74Pointer("Celsius")
			ab74Row.Ab74Value.Ab74CarrierMinimum = ab74Pointer(-32768)
			ab74Row.Ab74Value.Ab74CarrierMaximum = ab74Pointer(32767)
			ab74Row.Ab74Value.Ab74ScaledMinimum = ab74Pointer(-2048.0)
			ab74Row.Ab74Value.Ab74ScaledMaximum = ab74Pointer(2047.9375)
		}
		if ab74Index == 17 {
			ab74Row.Ab74Companion.Ab74Enumeration = map[string]string{"00": "okay", "55": "short-circuit", "AA": "interrupted"}
		}
		ab74Rows = append(ab74Rows, ab74Row)
	}
	return ab74Rows
}

func TestAb74ExactMetadata(ab74T *ab74testing.T) {
	ab74Rows, ab74Error := Ab74Entries()
	if ab74Error != nil {
		ab74T.Fatal(ab74Error)
	}
	if !ab74reflect.DeepEqual(ab74Rows, ab74Expected()) {
		ab74T.Fatal("ab74 independent nineteen-row contract differs")
	}
	for _, ab74Row := range ab74Expected() {
		ab74Found, ab74OK, ab74Error := Ab74Lookup(ab74Row.Ab74Key)
		if ab74Error != nil || !ab74OK || !ab74reflect.DeepEqual(ab74Found, ab74Row) {
			ab74T.Fatal("ab74 exact control", ab74Error)
		}
	}
}

func TestAb74IdentityAxes(ab74T *ab74testing.T) {
	for ab74Selector := byte(0); ab74Selector <= 18; ab74Selector++ {
		for ab74Axis := 0; ab74Axis < ab74reflect.TypeOf(Ab74Key{}).NumField(); ab74Axis++ {
			for ab74Mutation := 0; ab74Mutation < 3; ab74Mutation++ {
				ab74Input := ab74Key(ab74Selector)
				ab74Slot := ab74reflect.ValueOf(&ab74Input).Elem().Field(ab74Axis)
				if ab74Mutation == 0 {
					ab74Slot.SetZero()
				} else {
					ab74Value := ab74reflect.New(ab74Slot.Type().Elem())
					if ab74Value.Elem().Kind() == ab74reflect.String {
						if ab74Mutation == 2 {
							ab74Value.Elem().SetString("ab74_other")
						}
					} else if ab74Mutation == 2 {
						ab74Value.Elem().SetUint(19)
						if ab74Axis == 4 {
							ab74Value.Elem().SetUint(3)
						}
					}
					ab74Slot.Set(ab74Value)
				}
				ab74Found, ab74OK, ab74Error := Ab74Lookup(ab74Input)
				ab74Want := ab74Axis == 6 && ab74Mutation == 1
				if ab74Error != nil || ab74OK != ab74Want || (ab74Want && !ab74reflect.DeepEqual(ab74Found, ab74Expected()[0])) {
					ab74T.Fatalf("ab74 selector %d axis %d mutation %d: %v %v", ab74Selector, ab74Axis, ab74Mutation, ab74OK, ab74Error)
				}
			}
		}
	}
}

func ab74Marshal(ab74T *ab74testing.T, ab74Input any) []byte {
	ab74T.Helper()
	ab74Raw, ab74Error := ab74json.Marshal(ab74Input)
	if ab74Error != nil {
		ab74T.Fatal(ab74Error)
	}
	return ab74Raw
}

func TestAb74RejectInvalidCatalog(ab74T *ab74testing.T) {
	for _, ab74Raw := range [][]byte{nil, []byte("{"), []byte("null"), []byte("[]"), append(ab74bytes.Clone(ab74RowsJSON), []byte(" {}")...)} {
		if ab74Rows, ab74Error := ab74Decode(ab74Raw); ab74Error == nil || ab74Rows != nil {
			ab74T.Fatal("ab74 malformed catalog admitted")
		}
	}
	for _, ab74Mutation := range []string{"ab74_partial", "ab74_extra", "ab74_duplicate", "ab74_unknown", "ab74_zero_width", "ab74_order", "ab74_default", "ab74_base"} {
		ab74Rows := ab74Expected()
		switch ab74Mutation {
		case "ab74_partial":
			ab74Rows = ab74Rows[:18]
		case "ab74_extra":
			ab74Rows = append(ab74Rows, ab74Rows[0])
		case "ab74_duplicate":
			ab74Rows[1].Ab74Key = ab74Rows[0].Ab74Key
		case "ab74_unknown":
			ab74Rows[0].Ab74Key.Ab74Selector = ab74Pointer(byte(19))
		case "ab74_zero_width":
			ab74Rows[0].Ab74Value.Ab74Width = ab74Pointer(0)
		case "ab74_order":
			ab74Rows[17].Ab74Value.Ab74ByteOrder = "little-endian"
		case "ab74_default":
			ab74Rows[1].Ab74Value.Ab74Default = ab74Pointer(0.0)
		case "ab74_base":
			ab74Rows[1].Ab74Value.Ab74Offset = ab74Pointer(1)
		}
		if ab74Got, ab74Error := ab74Decode(ab74Marshal(ab74T, ab74Rows)); ab74Error == nil || ab74Got != nil {
			ab74T.Fatal(ab74Mutation, "admitted")
		}
	}
	var ab74Rows []map[string]any
	if ab74Error := ab74json.Unmarshal(ab74RowsJSON, &ab74Rows); ab74Error != nil {
		ab74T.Fatal(ab74Error)
	}
	for _, ab74Object := range []map[string]any{ab74Rows[0], ab74Rows[0]["key"].(map[string]any), ab74Rows[0]["value"].(map[string]any)} {
		for ab74Name, ab74Value := range ab74Object {
			delete(ab74Object, ab74Name)
			if _, ab74Error := ab74Decode(ab74Marshal(ab74T, ab74Rows)); ab74Error == nil {
				ab74T.Fatalf("ab74 missing %s admitted", ab74Name)
			}
			ab74Object[ab74Name] = ab74Value
		}
		ab74Object["ab74_unknown"] = true
		if _, ab74Error := ab74Decode(ab74Marshal(ab74T, ab74Rows)); ab74Error == nil {
			ab74T.Fatal("ab74 unknown property admitted")
		}
		delete(ab74Object, "ab74_unknown")
	}
	if _, ab74Error := ab74Decode(ab74RowsJSON); ab74Error != nil {
		ab74T.Fatal("ab74 valid control", ab74Error)
	}
}

func ab74Observe(ab74Row Ab74Entry, ab74Raw []byte) ([]byte, []byte, []byte) {
	ab74Evidence := ab74bytes.Clone(ab74Raw)
	ab74Extract := func(ab74Field *Ab74Field) []byte {
		if ab74Field == nil || ab74Field.Ab74Offset == nil || ab74Field.Ab74Width == nil || ab74Row.Ab74Base != "row-relative" {
			return nil
		}
		ab74Start, ab74Length := *ab74Field.Ab74Offset, *ab74Field.Ab74Width
		if ab74Start < 0 || ab74Length <= 0 || ab74Start > len(ab74Raw) || ab74Length > len(ab74Raw)-ab74Start {
			return nil
		}
		return ab74bytes.Clone(ab74Raw[ab74Start : ab74Start+ab74Length])
	}
	return ab74Extract(&ab74Row.Ab74Value), ab74Extract(ab74Row.Ab74Companion), ab74Evidence
}

func TestAb74RawCompanions(ab74T *ab74testing.T) {
	for _, ab74Case := range []struct {
		ab74Selector                       byte
		ab74ValueWidth, ab74CompanionWidth int
	}{
		{9, 1, 1}, {10, 1, 1}, {11, 2, 2}, {12, 2, 2}, {13, 4, 4}, {14, 4, 4}, {18, 2, 2},
	} {
		ab74Row, ab74OK, ab74Error := Ab74Lookup(ab74Key(ab74Case.ab74Selector))
		if ab74Error != nil || !ab74OK {
			ab74T.Fatal(ab74Error)
		}
		for ab74Index := 0; ab74Index < ab74Case.ab74CompanionWidth; ab74Index++ {
			ab74Input := ab74bytes.Repeat([]byte{0x23}, ab74Case.ab74ValueWidth+ab74Case.ab74CompanionWidth)
			ab74Input[ab74Case.ab74ValueWidth+ab74Index] = byte(0x40 + ab74Index)
			ab74Value, ab74Companion, ab74Evidence := ab74Observe(ab74Row, ab74Input)
			if !ab74bytes.Equal(ab74Value, ab74Input[:ab74Case.ab74ValueWidth]) || !ab74bytes.Equal(ab74Companion, ab74Input[ab74Case.ab74ValueWidth:]) || !ab74bytes.Equal(ab74Evidence, ab74Input) {
				ab74T.Fatalf("ab74 row %d companion byte %d lost", ab74Case.ab74Selector, ab74Index)
			}
			ab74Value[0], ab74Companion[0], ab74Evidence[0] = 0, 0, 0
			if ab74Input[0] != 0x23 || ab74Input[ab74Case.ab74ValueWidth] == 0 {
				ab74T.Fatal("ab74 returned bytes alias input")
			}
		}
		for ab74Length := 0; ab74Length < ab74Case.ab74ValueWidth+ab74Case.ab74CompanionWidth; ab74Length++ {
			ab74Input := ab74bytes.Repeat([]byte{0x23}, ab74Length)
			ab74Value, ab74Companion, ab74Evidence := ab74Observe(ab74Row, ab74Input)
			if (ab74Value != nil) != (ab74Length >= ab74Case.ab74ValueWidth) || ab74Companion != nil || !ab74bytes.Equal(ab74Evidence, ab74Input) {
				ab74T.Fatal("ab74 partial row lost independent field or raw evidence")
			}
		}
		ab74Input := ab74bytes.Repeat([]byte{0x23}, ab74Case.ab74ValueWidth+ab74Case.ab74CompanionWidth+1)
		_, _, ab74Evidence := ab74Observe(ab74Row, ab74Input)
		if !ab74bytes.Equal(ab74Evidence, ab74Input) {
			ab74T.Fatal("ab74 extra raw byte discarded")
		}
	}
	for _, ab74Selector := range []byte{0, 15, 16} {
		ab74Row, _, ab74Error := Ab74Lookup(ab74Key(ab74Selector))
		if ab74Error != nil {
			ab74T.Fatal(ab74Error)
		}
		ab74Value, ab74Companion, ab74Evidence := ab74Observe(ab74Row, []byte{0, 0x55, 0xAB})
		if ab74Value != nil || ab74Companion != nil || !ab74bytes.Equal(ab74Evidence, []byte{0, 0x55, 0xAB}) {
			ab74T.Fatal("ab74 opaque row acquired width")
		}
	}
}

func TestAb74FiniteEnumerations(ab74T *ab74testing.T) {
	ab74Boolean, _, ab74Error := Ab74Lookup(ab74Key(8))
	if ab74Error != nil {
		ab74T.Fatal(ab74Error)
	}
	ab74Status, _, ab74Error := Ab74Lookup(ab74Key(17))
	if ab74Error != nil {
		ab74T.Fatal(ab74Error)
	}
	for ab74Byte := 0; ab74Byte < 256; ab74Byte++ {
		ab74Label := ab74fmt.Sprintf("%02X", ab74Byte)
		ab74WantBoolean, ab74WantStatus := "", ""
		switch ab74Byte {
		case 0x0F:
			ab74WantBoolean = "true"
		case 0xF0:
			ab74WantBoolean = "false"
		}
		switch ab74Byte {
		case 0x00:
			ab74WantStatus = "okay"
		case 0x55:
			ab74WantStatus = "short-circuit"
		case 0xAA:
			ab74WantStatus = "interrupted"
		}
		ab74GotBoolean, ab74BooleanOK := ab74Boolean.Ab74Value.Ab74Enumeration[ab74Label]
		ab74GotStatus, ab74StatusOK := ab74Status.Ab74Companion.Ab74Enumeration[ab74Label]
		if ab74GotBoolean != ab74WantBoolean || ab74BooleanOK != (ab74WantBoolean != "") || ab74GotStatus != ab74WantStatus || ab74StatusOK != (ab74WantStatus != "") {
			ab74T.Fatalf("ab74 byte %s gained or lost meaning", ab74Label)
		}
	}
	ab74Ordinary, _, ab74Error := Ab74Lookup(ab74Key(1))
	if ab74Error != nil || len(ab74Ordinary.Ab74Value.Ab74Enumeration) != 0 {
		ab74T.Fatal("ab74 incomplete Boolean gained universal meaning")
	}
	ab74Value, ab74Companion, ab74Evidence := ab74Observe(ab74Status, []byte{0x80, 0x01, 0xFF})
	if !ab74bytes.Equal(ab74Value, []byte{0x80, 0x01}) || !ab74bytes.Equal(ab74Companion, []byte{0xFF}) || !ab74bytes.Equal(ab74Evidence, []byte{0x80, 0x01, 0xFF}) || ab74Status.Ab74Value.Ab74ByteOrder != "Unknown" {
		ab74T.Fatal("ab74 unknown status or byte order lost raw candidate")
	}
}

func ab74MutateMemory(ab74Input ab74reflect.Value) {
	switch ab74Input.Kind() {
	case ab74reflect.Pointer:
		if !ab74Input.IsNil() {
			ab74MutateMemory(ab74Input.Elem())
		}
	case ab74reflect.Struct:
		for ab74Index := 0; ab74Index < ab74Input.NumField(); ab74Index++ {
			ab74MutateMemory(ab74Input.Field(ab74Index))
		}
	case ab74reflect.Map:
		if !ab74Input.IsNil() {
			ab74Input.SetMapIndex(ab74reflect.ValueOf("ab74_mutated"), ab74reflect.ValueOf("ab74_mutated"))
		}
	case ab74reflect.String:
		ab74Input.SetString("ab74_mutated")
	case ab74reflect.Int:
		ab74Input.SetInt(99)
	case ab74reflect.Uint8:
		ab74Input.SetUint(99)
	case ab74reflect.Float64:
		ab74Input.SetFloat(99)
	}
}

func TestAb74DetachedMemory(ab74T *ab74testing.T) {
	ab74Before, ab74Error := Ab74Entries()
	if ab74Error != nil {
		ab74T.Fatal(ab74Error)
	}
	ab74Copy, ab74Error := Ab74Entries()
	if ab74Error != nil {
		ab74T.Fatal(ab74Error)
	}
	for ab74Index := range ab74Copy {
		ab74MutateMemory(ab74reflect.ValueOf(&ab74Copy[ab74Index]).Elem())
	}
	ab74After, ab74Error := Ab74Entries()
	if ab74Error != nil || !ab74reflect.DeepEqual(ab74Before, ab74After) {
		ab74T.Fatal("ab74 entries share nested memory", ab74Error)
	}
	for ab74Index, ab74Expected := range ab74Before {
		ab74KeyCopy := ab74Key(byte(ab74Index))
		ab74First, ab74OK, ab74Error := Ab74Lookup(ab74KeyCopy)
		if ab74Error != nil || !ab74OK {
			ab74T.Fatal(ab74Error)
		}
		ab74MutateMemory(ab74reflect.ValueOf(&ab74First).Elem())
		ab74Again, ab74OK, ab74Error := Ab74Lookup(ab74KeyCopy)
		if ab74Error != nil || !ab74OK || !ab74reflect.DeepEqual(ab74Expected, ab74Again) {
			ab74T.Fatal("ab74 lookup or caller identity shares nested memory", ab74Error)
		}
	}
	ab74Input := ab74bytes.Clone(ab74RowsJSON)
	ab74Decoded, ab74Error := ab74Decode(ab74Input)
	if ab74Error != nil {
		ab74T.Fatal(ab74Error)
	}
	clear(ab74Input)
	if !ab74reflect.DeepEqual(ab74Decoded, ab74Before) {
		ab74T.Fatal("ab74 decoded catalog aliases raw JSON")
	}
}
