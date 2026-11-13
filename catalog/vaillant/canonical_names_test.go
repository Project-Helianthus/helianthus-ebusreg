package vaillant

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCanonicalAnnotationBoundary(t *testing.T) {
	raw, err := B524NameAnnotationEntries()
	if err != nil {
		t.Fatal(err)
	}
	var rows []map[string]json.RawMessage
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		for _, forbidden := range []string{"previous_name", "legacy_name", "display_alias"} {
			if _, present := row[forbidden]; present {
				t.Fatalf("documentary field emitted: %s", forbidden)
			}
		}
		var key B524NameAnnotationKey
		if err := json.Unmarshal(row["key"], &key); err != nil {
			t.Fatal(err)
		}
		if *key.B524NameAnnotationOperation == 2 {
			if _, present := row["raw_form"]; present {
				t.Fatal("documentary correspondence emitted as native evidence")
			}
			if _, present := row["native_label"]; present {
				t.Fatal("OP02 acquired OP06 evidence")
			}
		}
	}
	for _, forbidden := range []string{"previous_name", "legacy_name", "display_alias"} {
		row := b524NameAnnotationExpectedRow(t, "06|0A|0007|RoomHumidity")
		row[forbidden] = "synthetic_previous_name"
		if _, _, err := b524NameAnnotationValidate(b524NameAnnotationEncode(t, []any{row})); err == nil {
			t.Fatalf("runtime alias accepted: %s", forbidden)
		}
	}
}

func TestCanonicalAnnotationIdentityCollisions(t *testing.T) {
	base := b524NameAnnotationExpectedRow(t, "06|0A|0007|RoomHumidity")["key"].(B524NameAnnotationKey)
	for axis := 0; axis < reflect.TypeOf(base).NumField(); axis++ {
		name := reflect.TypeOf(base).Field(axis).Name
		t.Run(name, func(t *testing.T) {
			changed := base
			value := reflect.ValueOf(&changed).Elem().Field(axis)
			switch value.Kind() {
			case reflect.String:
				value.SetString("synthetic_other")
			case reflect.Uint8:
				value.SetUint(254)
			case reflect.Ptr:
				value.Set(reflect.New(value.Type().Elem()))
				value.Elem().SetUint(254)
			default:
				t.Fatal("uncovered identity axis")
			}
			first := b524NameAnnotationExpectedRow(t, "06|0A|0007|RoomHumidity")
			second := b524NameAnnotationExpectedRow(t, "06|0A|0007|RoomHumidity")
			second["key"] = changed
			second["native_label"] = "synthetic_distinct_evidence"
			catalog := &b524NameAnnotationCatalog{b524NameAnnotationSource: string(b524NameAnnotationEncode(t, []any{first, second}))}
			index, err := catalog.b524NameAnnotationLoad()
			if err != nil || len(index) != 2 {
				t.Fatalf("identity axis collapsed: %v", err)
			}
			for _, row := range []map[string]any{first, second} {
				b524NameAnnotationEqualJSON(t, index[string(b524NameAnnotationEncode(t, row["key"]))], b524NameAnnotationEncode(t, row))
			}
			second["key"] = base
			if _, _, err := b524NameAnnotationValidate(b524NameAnnotationEncode(t, []any{first, second})); err == nil {
				t.Fatal("duplicate identity admitted")
			}
		})
	}
}

func TestRenamedCatalogFieldJSONKeys(t *testing.T) {
	cases := []struct {
		value any
		keys  map[string]string
	}{
		{B514ActivationVectorIdentity{}, map[string]string{"B514ActivationVectorFamily": "Ab80Family", "B514ActivationVectorVariant": "Ab80Variant", "B514ActivationVectorOperation": "Ab80Operation", "B514ActivationVectorDirection": "Ab80Direction"}},
		{B514ActivationVectorBit{}, map[string]string{"B514ActivationVectorPresent": "Ab80Present", "B514ActivationVectorValue": "Ab80Value"}},
		{B514ActivationVectorScalar{}, map[string]string{"B514ActivationVectorPresent": "Ab80Present", "B514ActivationVectorRaw": "Ab80Raw", "B514ActivationVectorInterpretation": "Ab80Interpretation"}},
		{B525SpectrumBlockContext{}, map[string]string{"B525SpectrumBlockPrimary": "Ab82Primary", "B525SpectrumBlockSecondary": "Ab82Secondary", "B525SpectrumBlockFamily": "Ab82Family", "B525SpectrumBlockVariant": "Ab82Variant", "B525SpectrumBlockOperation": "Ab82Operation", "B525SpectrumBlockDirection": "Ab82Direction", "B525SpectrumBlockSelector": "Ab82Selector"}},
		{B525SpectrumBlockInput{}, map[string]string{"B525SpectrumBlockContext": "Ab82Context", "B525SpectrumBlockValues": "Ab82Values"}},
		{B525SpectrumBlockEntry{}, map[string]string{"B525SpectrumBlockIndex": "Ab82Index", "B525SpectrumBlockRawByte": "Ab82RawByte", "B525SpectrumBlockSigned": "Ab82Signed"}},
		{B525SpectrumBlock{}, map[string]string{"B525SpectrumBlockContext": "Ab82Context", "B525SpectrumBlockState": "Ab82State", "B525SpectrumBlockGrade": "Ab82Grade", "B525SpectrumBlockStart": "Ab82Start", "B525SpectrumBlockRaw": "Ab82Raw", "B525SpectrumBlockEntries": "Ab82Entries", "B525SpectrumBlockMask": "Ab82Mask", "B525SpectrumBlockIncomplete": "Ab82Incomplete", "B525SpectrumBlockExtra": "Ab82Extra", "B525SpectrumBlockRemainder": "Ab82Remainder"}},
		{B525SpectrumJoinBlock{}, map[string]string{"B525SpectrumBlockSelector": "Ab82Selector", "B525SpectrumBlockState": "Ab82State", "B525SpectrumBlockContributed": "Ab82Contributed", "B525SpectrumBlockConflict": "Ab82Conflict"}},
		{B525SpectrumBlockJoined{}, map[string]string{"B525SpectrumBlockState": "Ab82State", "B525SpectrumBlockGrade": "Ab82Grade", "B525SpectrumBlockComplete": "Ab82Complete", "B525SpectrumBlocks": "Ab82Blocks", "B525SpectrumBlockEntries": "Ab82Entries", "B525SpectrumBlockMask": "Ab82Mask", "B525SpectrumBlockMissing": "Ab82Missing", "B525SpectrumBlockConflicts": "Ab82Conflicts"}},
	}
	for _, tc := range cases {
		data, err := json.Marshal(tc.value)
		if err != nil {
			t.Fatal(err)
		}
		var object map[string]json.RawMessage
		if err := json.Unmarshal(data, &object); err != nil {
			t.Fatal(err)
		}
		typ := reflect.TypeOf(tc.value)
		for field, key := range tc.keys {
			if _, ok := object[key]; !ok {
				t.Errorf("%s.%s missing JSON key %s", typ, field, key)
			}
			if _, ok := object[field]; ok {
				t.Errorf("%s exposes renamed key %s", typ, field)
			}
			if member, ok := typ.FieldByName(field); !ok || member.Tag.Get("json") != key {
				t.Errorf("%s.%s does not decode JSON key %s", typ, field, key)
			}
		}
	}
}
