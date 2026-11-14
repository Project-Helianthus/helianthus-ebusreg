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
		{B514ActivationVectorIdentity{}, map[string]string{"B514ActivationVectorFamily": "family", "B514ActivationVectorVariant": "variant", "B514ActivationVectorOperation": "operation", "B514ActivationVectorDirection": "direction"}},
		{B514ActivationVectorBit{}, map[string]string{"B514ActivationVectorPresent": "present", "B514ActivationVectorValue": "value"}},
		{B514ActivationVectorScalar{}, map[string]string{"B514ActivationVectorPresent": "present", "B514ActivationVectorRaw": "raw", "B514ActivationVectorInterpretation": "interpretation"}},
		{B525SpectrumBlockContext{}, map[string]string{"B525SpectrumBlockPrimary": "primary", "B525SpectrumBlockSecondary": "secondary", "B525SpectrumBlockFamily": "family", "B525SpectrumBlockVariant": "variant", "B525SpectrumBlockOperation": "operation", "B525SpectrumBlockDirection": "direction", "B525SpectrumBlockSelector": "selector"}},
		{B525SpectrumBlockInput{}, map[string]string{"B525SpectrumBlockContext": "context", "B525SpectrumBlockValues": "values"}},
		{B525SpectrumBlockEntry{}, map[string]string{"B525SpectrumBlockIndex": "index", "B525SpectrumBlockRawByte": "raw_byte", "B525SpectrumBlockSigned": "signed"}},
		{B525SpectrumBlock{}, map[string]string{"B525SpectrumBlockContext": "context", "B525SpectrumBlockState": "state", "B525SpectrumBlockGrade": "grade", "B525SpectrumBlockStart": "start", "B525SpectrumBlockRaw": "raw", "B525SpectrumBlockEntries": "entries", "B525SpectrumBlockMask": "mask", "B525SpectrumBlockIncomplete": "incomplete", "B525SpectrumBlockExtra": "extra", "B525SpectrumBlockRemainder": "remainder"}},
		{B525SpectrumJoinBlock{}, map[string]string{"B525SpectrumBlockSelector": "selector", "B525SpectrumBlockState": "state", "B525SpectrumBlockContributed": "contributed", "B525SpectrumBlockConflict": "conflict"}},
		{B525SpectrumBlockJoined{}, map[string]string{"B525SpectrumBlockState": "state", "B525SpectrumBlockGrade": "grade", "B525SpectrumBlockComplete": "complete", "B525SpectrumBlocks": "blocks", "B525SpectrumBlockEntries": "entries", "B525SpectrumBlockMask": "mask", "B525SpectrumBlockMissing": "missing", "B525SpectrumBlockConflicts": "conflicts"}},
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
