package system

import (
	ab66bytes "bytes"
	ab66json "encoding/json"
	ab66fmt "fmt"
	ab66io "io"
	ab66os "os"
	ab66reflect "reflect"
	ab66testing "testing"
)

type ab66Fixture struct {
	Ab66Keys     []Ab66B524HintKey `json:"ab66Keys"`
	Ab66Controls []struct {
		Ab66Key  Ab66B524HintKey
		Ab66Hint Ab66B524HintProperties
	} `json:"ab66Controls"`
}

func ab66LoadFixture(ab66T *ab66testing.T) ab66Fixture {
	ab66T.Helper()
	ab66Raw, ab66Err := ab66os.ReadFile("testdata/ab66_b524_unscoped_hints.json")
	if ab66Err != nil {
		ab66T.Fatal(ab66Err)
	}
	var ab66Data ab66Fixture
	ab66Decoder := ab66json.NewDecoder(ab66bytes.NewReader(ab66Raw))
	ab66Decoder.DisallowUnknownFields()
	if ab66Err = ab66Decoder.Decode(&ab66Data); ab66Err != nil {
		ab66T.Fatal(ab66Err)
	}
	return ab66Data
}

func TestAb66EnumerationAndLookup(ab66T *ab66testing.T) {
	ab66Data := ab66LoadFixture(ab66T)
	ab66Results := Ab66B524UnscopedHints()
	if len(ab66Results) != 46 || len(ab66Data.Ab66Keys) != 46 || len(b524StaticConstraintsCatalog) != 46 {
		ab66T.Fatalf("ab66 count: %d/%d/%d", len(ab66Results), len(ab66Data.Ab66Keys), len(b524StaticConstraintsCatalog))
	}
	ab66Seen := make(map[Ab66B524HintKey]bool)
	for ab66I, ab66Key := range ab66Data.Ab66Keys {
		ab66T.Run(ab66fmt.Sprintf("ab66_key_%d", ab66I), func(ab66T *ab66testing.T) {
			if ab66Seen[ab66Key] {
				ab66T.Fatal("ab66 duplicate")
			}
			ab66Seen[ab66Key] = true
			ab66Expected := b524StaticConstraintsCatalog[b524ConstraintSelector{Opcode: ab66Key.Ab66Opcode, Group: ab66Key.Ab66Group, Register: ab66Key.Ab66Register}]
			ab66Want := Ab66B524HintResult{
				Ab66Key: ab66Key, Ab66Status: Ab66B524HintUnscoped, Ab66Qualified: false,
				Ab66Applicability: Ab66B524ApplicabilityUnknown,
				Ab66Hint:          Ab66B524HintProperties{Ab66Type: ab66Expected.Type, Ab66Min: ab66Expected.Min, Ab66Max: ab66Expected.Max, Ab66Step: ab66Expected.Step},
			}
			if ab66Results[ab66I] != ab66Want || Ab66LookupB524UnscopedHint(ab66Key) != ab66Want {
				ab66T.Fatalf("ab66 enumeration/lookup mismatch: %+v; want %+v", ab66Results[ab66I], ab66Want)
			}
			ab66Key.Ab66Opcode = 6
			ab66AssertAbsent(ab66T, ab66Key)
		})
	}
}

func TestAb66IndependentNumericControls(ab66T *ab66testing.T) {
	ab66Data := ab66LoadFixture(ab66T)
	if len(ab66Data.Ab66Controls) != 3 {
		ab66T.Fatal("ab66 control count")
	}
	for ab66I, ab66Control := range ab66Data.Ab66Controls {
		ab66T.Run(ab66fmt.Sprintf("ab66_numeric_%d", ab66I), func(ab66T *ab66testing.T) {
			ab66Got := Ab66LookupB524UnscopedHint(ab66Control.Ab66Key)
			if ab66Got.Ab66Hint != ab66Control.Ab66Hint {
				ab66T.Fatalf("ab66 properties: %+v want %+v", ab66Got, ab66Control)
			}
		})
	}
}

func ab66AssertAbsent(ab66T *ab66testing.T, ab66Key Ab66B524HintKey) {
	ab66T.Helper()
	ab66Want := Ab66B524HintResult{Ab66Key: ab66Key, Ab66Status: Ab66B524HintAbsent, Ab66Applicability: Ab66B524ApplicabilityUnknown}
	if ab66Got := Ab66LookupB524UnscopedHint(ab66Key); ab66Got != ab66Want {
		ab66T.Fatalf("ab66 absent: %+v want %+v", ab66Got, ab66Want)
	}
}

func TestAb66IndependentKeyAxes(ab66T *ab66testing.T) {
	for ab66Opcode := 0; ab66Opcode < 256; ab66Opcode++ {
		if ab66Opcode != 2 {
			ab66AssertAbsent(ab66T, Ab66B524HintKey{Ab66Opcode: byte(ab66Opcode), Ab66Group: 0, Ab66Register: 2})
		}
	}
	for _, ab66Key := range []Ab66B524HintKey{
		{Ab66Opcode: 2, Ab66Group: 0, Ab66Register: 1},
		{Ab66Opcode: 6, Ab66Group: 0, Ab66Register: 1},
		{Ab66Opcode: 2, Ab66Group: 0, Ab66Register: 0},
		{Ab66Opcode: 2, Ab66Group: 0, Ab66Register: 0x0200},
		{Ab66Opcode: 2, Ab66Group: 0, Ab66Register: 0xffff},
		{Ab66Opcode: 2, Ab66Group: 0xff, Ab66Register: 2},
		{},
	} {
		ab66AssertAbsent(ab66T, ab66Key)
	}
	for _, ab66Key := range []Ab66B524HintKey{
		{Ab66Opcode: 2, Ab66Group: 0, Ab66Register: 2},
		{Ab66Opcode: 2, Ab66Group: 3, Ab66Register: 2},
		{Ab66Opcode: 2, Ab66Group: 0, Ab66Register: 4},
	} {
		if Ab66LookupB524UnscopedHint(ab66Key).Ab66Status != Ab66B524HintUnscoped {
			ab66T.Fatal("ab66 valid axis control")
		}
	}
}

func TestAb66OwnedValues(ab66T *ab66testing.T) {
	ab66Before := Ab66B524UnscopedHints()
	ab66Results := Ab66B524UnscopedHints()
	ab66Key := ab66Results[0].Ab66Key
	ab66Result := Ab66LookupB524UnscopedHint(ab66Key)
	ab66Key.Ab66Opcode, ab66Key.Ab66Group, ab66Key.Ab66Register = 6, 255, 65535
	if ab66Result != ab66Before[0] {
		ab66T.Fatal("ab66 input alias")
	}
	ab66Result.Ab66Key = ab66Key
	ab66Result.Ab66Status = Ab66B524HintAbsent
	ab66Result.Ab66Qualified = true
	ab66Result.Ab66Applicability = "ab66_changed"
	ab66Result.Ab66Hint = Ab66B524HintProperties{Ab66Type: "ab66_changed", Ab66Min: 11, Ab66Max: 12, Ab66Step: 13}
	for ab66I := range ab66Results {
		ab66Results[ab66I] = ab66Result
	}
	ab66Results = append(ab66Results, ab66Result)
	if len(ab66Results) != len(ab66Before)+1 || ab66Results[len(ab66Before)] != ab66Result {
		ab66T.Fatal("ab66 caller append")
	}
	if !ab66reflect.DeepEqual(Ab66B524UnscopedHints(), ab66Before) {
		ab66T.Fatal("ab66 enumeration alias")
	}
	if Ab66LookupB524UnscopedHint(ab66Before[0].Ab66Key) != ab66Before[0] {
		ab66T.Fatal("ab66 lookup alias")
	}
	ab66Native, ab66OK := lookupB524Constraint(2, 0, 2)
	if !ab66OK || ab66Native.Type != "f32_range" || ab66Native.Min != -26 || ab66Native.Max != 10 || ab66Native.Step != 1 {
		ab66T.Fatal("ab66 native table changed")
	}
}

func TestAb66PublicShape(ab66T *ab66testing.T) {
	for _, ab66Case := range []struct {
		Ab66Value  any
		Ab66Fields []string
	}{
		{Ab66B524HintKey{}, []string{"Ab66Opcode", "Ab66Group", "Ab66Register"}},
		{Ab66B524HintProperties{}, []string{"Ab66Type", "Ab66Min", "Ab66Max", "Ab66Step"}},
		{Ab66B524HintResult{}, []string{"Ab66Key", "Ab66Status", "Ab66Qualified", "Ab66Applicability", "Ab66Hint"}},
	} {
		ab66Type := ab66reflect.TypeOf(ab66Case.Ab66Value)
		if ab66Type.NumField() != len(ab66Case.Ab66Fields) {
			ab66T.Fatal("ab66 unexpected identity/bound field")
		}
		for ab66I, ab66Name := range ab66Case.Ab66Fields {
			if ab66Type.Field(ab66I).Name != ab66Name {
				ab66T.Fatal("ab66 field contract")
			}
		}
	}
	if Ab66B524HintUnscoped != "unscoped_catalog" || Ab66B524HintAbsent != "absent" || Ab66B524ApplicabilityUnknown != "Unknown" {
		ab66T.Fatal("ab66 status/grade contract")
	}
}

func ab66ValidateSyntheticKeys(ab66Raw []byte) ([]Ab66B524HintKey, error) {
	var ab66Wire []struct {
		Ab66Opcode   *byte
		Ab66Group    *byte
		Ab66Register *uint16
	}
	ab66Decoder := ab66json.NewDecoder(ab66bytes.NewReader(ab66Raw))
	ab66Decoder.DisallowUnknownFields()
	if ab66Err := ab66Decoder.Decode(&ab66Wire); ab66Err != nil {
		return nil, ab66Err
	}
	var ab66Extra any
	if ab66Err := ab66Decoder.Decode(&ab66Extra); ab66Err != ab66io.EOF {
		return nil, ab66fmt.Errorf("ab66 trailing input")
	}
	if len(ab66Wire) == 0 {
		return nil, ab66fmt.Errorf("ab66 missing keys")
	}
	ab66Seen := make(map[Ab66B524HintKey]bool)
	ab66Keys := make([]Ab66B524HintKey, 0, len(ab66Wire))
	for _, ab66Item := range ab66Wire {
		if ab66Item.Ab66Opcode == nil || ab66Item.Ab66Group == nil || ab66Item.Ab66Register == nil {
			return nil, ab66fmt.Errorf("ab66 missing field")
		}
		ab66Key := Ab66B524HintKey{Ab66Opcode: *ab66Item.Ab66Opcode, Ab66Group: *ab66Item.Ab66Group, Ab66Register: *ab66Item.Ab66Register}
		if ab66Seen[ab66Key] {
			return nil, ab66fmt.Errorf("ab66 duplicate key")
		}
		ab66Seen[ab66Key] = true
		ab66Keys = append(ab66Keys, ab66Key)
	}
	return ab66Keys, nil
}

func TestAb66SyntheticValidationAndRawOwnership(ab66T *ab66testing.T) {
	for ab66I, ab66Raw := range []string{
		`[`, `null`, `[]`, `[{}]`,
		`[{"Ab66Opcode":2,"Ab66Register":2}]`,
		`[{"Ab66Group":0,"Ab66Register":2}]`,
		`[{"Ab66Opcode":2,"Ab66Group":0}]`,
		`[{"Ab66Opcode":2,"Ab66Group":null,"Ab66Register":2}]`,
		`[{"Ab66Opcode":2,"Ab66Group":0,"Ab66Register":2,"ab66Family":"ab66_synthetic"}]`,
		`[{"Ab66Opcode":256,"Ab66Group":0,"Ab66Register":2}]`,
		`[{"Ab66Opcode":2,"Ab66Group":0,"Ab66Register":2},{"Ab66Opcode":2,"Ab66Group":0,"Ab66Register":2}]`,
		`[{"Ab66Opcode":2,"Ab66Group":0,"Ab66Register":2}] {}`,
	} {
		ab66T.Run(ab66fmt.Sprintf("ab66_invalid_%d", ab66I), func(ab66T *ab66testing.T) {
			if _, ab66Err := ab66ValidateSyntheticKeys([]byte(ab66Raw)); ab66Err == nil {
				ab66T.Fatal("ab66 accepted invalid fixture")
			}
		})
	}
	ab66Raw := []byte(`[{"Ab66Opcode":0,"Ab66Group":0,"Ab66Register":0},{"Ab66Opcode":2,"Ab66Group":0,"Ab66Register":2}]`)
	ab66Saved := append([]byte(nil), ab66Raw...)
	ab66Keys, ab66Err := ab66ValidateSyntheticKeys(ab66Raw)
	if ab66Err != nil || len(ab66Keys) != 2 || ab66Keys[0] != (Ab66B524HintKey{}) {
		ab66T.Fatalf("ab66 explicit zero: %v", ab66Err)
	}
	if !ab66bytes.Equal(ab66Saved, ab66Raw) {
		ab66T.Fatal("ab66 raw changed")
	}
	for ab66I := range ab66Raw {
		ab66Raw[ab66I] = ' '
	}
	if ab66Keys[1] != (Ab66B524HintKey{Ab66Opcode: 2, Ab66Group: 0, Ab66Register: 2}) {
		ab66T.Fatal("ab66 raw alias")
	}
	ab66AssertAbsent(ab66T, ab66Keys[0])
	if Ab66LookupB524UnscopedHint(ab66Keys[1]).Ab66Status != Ab66B524HintUnscoped {
		ab66T.Fatal("ab66 valid raw control")
	}
}
