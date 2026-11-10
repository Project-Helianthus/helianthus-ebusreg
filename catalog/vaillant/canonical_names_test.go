package vaillant

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCanonicalAnnotationBoundary(t *testing.T) {
	raw, err := Ab36Entries()
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
		var key Ab36Key
		if err := json.Unmarshal(row["key"], &key); err != nil {
			t.Fatal(err)
		}
		if *key.Ab36Operation == 2 {
			if _, present := row["raw_form"]; present {
				t.Fatal("documentary correspondence emitted as native evidence")
			}
			if _, present := row["native_label"]; present {
				t.Fatal("OP02 acquired OP06 evidence")
			}
		}
	}
	for _, forbidden := range []string{"previous_name", "legacy_name", "display_alias"} {
		row := ab36ExpectedRow(t, "06|0A|0007|RoomHumidity")
		row[forbidden] = "synthetic_previous_name"
		if _, _, err := ab36Validate(ab36Encode(t, []any{row})); err == nil {
			t.Fatalf("runtime alias accepted: %s", forbidden)
		}
	}
}

func TestCanonicalAnnotationIdentityCollisions(t *testing.T) {
	base := ab36ExpectedRow(t, "06|0A|0007|RoomHumidity")["key"].(Ab36Key)
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
			first := ab36ExpectedRow(t, "06|0A|0007|RoomHumidity")
			second := ab36ExpectedRow(t, "06|0A|0007|RoomHumidity")
			second["key"] = changed
			second["native_label"] = "synthetic_distinct_evidence"
			catalog := &ab36Catalog{ab36Source: string(ab36Encode(t, []any{first, second}))}
			index, err := catalog.ab36Load()
			if err != nil || len(index) != 2 {
				t.Fatalf("identity axis collapsed: %v", err)
			}
			for _, row := range []map[string]any{first, second} {
				ab36EqualJSON(t, index[string(ab36Encode(t, row["key"]))], ab36Encode(t, row))
			}
			second["key"] = base
			if _, _, err := ab36Validate(ab36Encode(t, []any{first, second})); err == nil {
				t.Fatal("duplicate identity admitted")
			}
		})
	}
}
