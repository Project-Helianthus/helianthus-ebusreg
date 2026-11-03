package ebus_standard_catalog

import (
	ab84context "context"
	ab84errors "errors"
	ab84reflect "reflect"
	ab84testing "testing"
)

func ab84Int(ab84Value int) *int { return &ab84Value }

func ab84Key(ab84Name, ab84Direction string) Ab84Key {
	return Ab84Key{Ab84Family: "burner", Ab84Variant: "05_01_candidate", Ab84Primary: Ab84Axis{true, 5}, Ab84Secondary: Ab84Axis{true, 1}, Ab84Direction: ab84Direction, Ab84Field: ab84Name}
}

func ab84ExpectedField(ab84Name, ab84Direction string, ab84Position int, ab84Signedness, ab84Unit string) Ab84BurnerField {
	return Ab84BurnerField{Ab84Key: ab84Key(ab84Name, ab84Direction), Ab84LogicalPosition: ab84Int(ab84Position), Ab84WidthBits: ab84Int(8), Ab84Signedness: ab84Signedness, Ab84ByteOrder: "not_applicable", Ab84Conversion: "identity", Ab84Unit: ab84Unit, Ab84Replacement: []int{}, Ab84Grade: "Hypothesis", Ab84RawFallback: "retain_raw_bytes"}
}

func ab84ExpectedFields() []Ab84BurnerField {
	ab84Fields := []Ab84BurnerField{
		ab84ExpectedField("primary", "request", 0, "unsigned", "none"),
		ab84ExpectedField("secondary", "request", 1, "unsigned", "none"),
		ab84ExpectedField("heat_state", "request", 2, "unsigned", "none"),
		ab84ExpectedField("boiler_target", "request", 3, "unsigned", "degrees_Celsius"),
		ab84ExpectedField("dhw_target", "request", 4, "unsigned", "degrees_Celsius"),
		ab84ExpectedField("outside_temperature", "request", 5, "signed", "degrees_Celsius"),
		ab84ExpectedField("setting_percentage", "request", 6, "unsigned", "percent"),
		ab84ExpectedField("primary", "positive_declaration", 0, "unsigned", "none"),
		ab84ExpectedField("secondary", "positive_declaration", 1, "unsigned", "none"),
	}
	ab84Fields[0].Ab84Constant = ab84Int(5)
	ab84Fields[1].Ab84Constant = ab84Int(1)
	ab84Fields[7].Ab84Constant = ab84Int(5)
	ab84Fields[8].Ab84Constant = ab84Int(1)
	for ab84Index := 2; ab84Index < 7; ab84Index++ {
		ab84Fields[ab84Index].Ab84Default = ab84Int(0)
	}
	ab84Fields[2].Ab84Enumeration = map[int]string{0: "shutdown", 85: "DHW", 170: "heating", 204: "emission_check", 238: "controller_stop"}
	ab84Fields[3].Ab84Bounds = &Ab84Bounds{0, 90}
	ab84Fields[4].Ab84Bounds = &Ab84Bounds{0, 90}
	ab84Fields[5].Ab84Bounds = &Ab84Bounds{-128, 127}
	ab84Fields[6].Ab84Bounds = &Ab84Bounds{0, 255}
	for _, ab84Name := range []string{"physical_positive_body", "acknowledgement", "outside_3f_replacement", "control_permission", "scheduling", "other_05_candidate_bodies", "dd_check_semantics"} {
		ab84Fields = append(ab84Fields, Ab84BurnerField{Ab84Key: ab84Key(ab84Name, "unprovided"), Ab84Signedness: "Unknown", Ab84ByteOrder: "Unknown", Ab84Conversion: "Unknown", Ab84Unit: "Unknown", Ab84Grade: "Unknown", Ab84RawFallback: "retain_raw_bytes"})
	}
	return ab84Fields
}

func TestAb84ProductionMetadata(ab84T *ab84testing.T) {
	ab84Expected := Ab84BurnerBodies{Ab84RepresentationGrade: "Proven", Ab84Coordinates: "logical byte coordinates only; no physical offsets or positive-body layout", Ab84NativeSupport: "Unknown", Ab84Qualification: "Unknown", Ab84Freshness: "Unknown", Ab84Persistence: "Unknown", Ab84ExecutionPermission: "Unknown", Ab84Fields: ab84ExpectedFields()}
	if ab84Actual := Ab84BurnerMetadata(); !ab84reflect.DeepEqual(ab84Actual, ab84Expected) {
		ab84T.Fatalf("production metadata differs from independent expectation:\n%+v", ab84Actual)
	}
	for _, ab84ExpectedField := range ab84Expected.Ab84Fields {
		ab84Actual, ab84OK := Ab84LookupBurnerField(ab84ExpectedField.Ab84Key)
		if !ab84OK || !ab84reflect.DeepEqual(ab84Actual, ab84ExpectedField) {
			ab84T.Fatalf("lookup %s: %+v, %v", ab84ExpectedField.Ab84Key.Ab84Field, ab84Actual, ab84OK)
		}
	}
}

func ab84Evaluate(ab84Field Ab84BurnerField, ab84Logical []byte) (int, string, bool) {
	if ab84Field.Ab84LogicalPosition == nil || ab84Field.Ab84WidthBits == nil || *ab84Field.Ab84WidthBits != 8 || ab84Field.Ab84Conversion != "identity" || ab84Field.Ab84ByteOrder != "not_applicable" {
		return 0, "", false
	}
	ab84Position := *ab84Field.Ab84LogicalPosition
	if ab84Position < 0 || ab84Position >= len(ab84Logical) {
		return 0, "", false
	}
	ab84Value := int(ab84Logical[ab84Position])
	for _, ab84Replacement := range ab84Field.Ab84Replacement {
		if ab84Value == ab84Replacement {
			return ab84Value, "", false
		}
	}
	if ab84Field.Ab84Signedness == "signed" && ab84Value >= 128 {
		ab84Value -= 256
	} else if ab84Field.Ab84Signedness != "signed" && ab84Field.Ab84Signedness != "unsigned" {
		return 0, "", false
	}
	ab84Valid := ab84Field.Ab84Bounds == nil || (ab84Value >= ab84Field.Ab84Bounds.Ab84Min && ab84Value <= ab84Field.Ab84Bounds.Ab84Max)
	return ab84Value, ab84Field.Ab84Enumeration[ab84Value], ab84Valid
}

func TestAb84SyntheticCandidate(ab84T *ab84testing.T) {
	ab84Logical := []byte{5, 1, 0xAA, 90, 0, 0xFF, 255}
	for _, ab84Case := range []struct {
		Ab84Name  string
		Ab84Value int
		Ab84Label string
	}{
		{"heat_state", 170, "heating"}, {"boiler_target", 90, ""}, {"dhw_target", 0, ""}, {"outside_temperature", -1, ""}, {"setting_percentage", 255, ""},
	} {
		ab84Field, ab84OK := Ab84LookupBurnerField(ab84Key(ab84Case.Ab84Name, "request"))
		if !ab84OK || ab84Field.Ab84Grade != "Hypothesis" {
			ab84T.Fatal("missing hypothesis field", ab84Case.Ab84Name)
		}
		ab84Value, ab84Label, ab84Valid := ab84Evaluate(ab84Field, ab84Logical)
		if !ab84Valid || ab84Value != ab84Case.Ab84Value || ab84Label != ab84Case.Ab84Label {
			ab84T.Fatalf("%s: %d %q %v", ab84Case.Ab84Name, ab84Value, ab84Label, ab84Valid)
		}
	}
	for _, ab84Name := range []string{"boiler_target", "dhw_target"} {
		ab84Field, _ := Ab84LookupBurnerField(ab84Key(ab84Name, "request"))
		ab84Logical[*ab84Field.Ab84LogicalPosition] = 91
		if _, _, ab84Valid := ab84Evaluate(ab84Field, ab84Logical); ab84Valid {
			ab84T.Fatal("target 91 accepted", ab84Name)
		}
	}
	ab84Outside, _ := Ab84LookupBurnerField(ab84Key("outside_temperature", "request"))
	ab84Logical[*ab84Outside.Ab84LogicalPosition] = 0x3F
	if ab84Value, _, ab84Valid := ab84Evaluate(ab84Outside, ab84Logical); !ab84Valid || ab84Value != 63 {
		ab84T.Fatal("3F is numeric raw 63", ab84Value, ab84Valid)
	}
	ab84Heat, _ := Ab84LookupBurnerField(ab84Key("heat_state", "request"))
	if _, ab84OK := ab84Heat.Ab84Enumeration[0xDD]; ab84OK {
		ab84T.Fatal("DD semantics enabled")
	}
	for ab84Raw, ab84Label := range map[byte]string{0: "shutdown", 0x55: "DHW", 0xAA: "heating", 0xCC: "emission_check", 0xEE: "controller_stop"} {
		ab84Logical[*ab84Heat.Ab84LogicalPosition] = ab84Raw
		if _, ab84Actual, ab84Valid := ab84Evaluate(ab84Heat, ab84Logical); !ab84Valid || ab84Actual != ab84Label {
			ab84T.Fatal("heat alternative", ab84Raw, ab84Actual)
		}
	}
}

func ab84PositiveComplete(ab84Fields []Ab84BurnerField, ab84Logical []byte) bool {
	ab84Count := 0
	for _, ab84Field := range ab84Fields {
		if ab84Field.Ab84Key.Ab84Direction != "positive_declaration" {
			continue
		}
		ab84Count++
		if ab84Field.Ab84Constant == nil {
			return false
		}
		ab84Value, _, ab84Valid := ab84Evaluate(ab84Field, ab84Logical)
		if !ab84Valid || ab84Value != *ab84Field.Ab84Constant {
			return false
		}
	}
	return ab84Count == 2 && len(ab84Logical) == 2
}

func TestAb84PositiveDeclaration(ab84T *ab84testing.T) {
	ab84Fields := Ab84BurnerMetadata().Ab84Fields
	for _, ab84Case := range []struct {
		Ab84Bytes    []byte
		Ab84Complete bool
	}{
		{[]byte{5, 1}, true}, {[]byte{5}, false}, {[]byte{1}, false}, {[]byte{5, 4}, false}, {[]byte{5, 1, 0}, false}, {nil, false},
	} {
		if ab84Actual := ab84PositiveComplete(ab84Fields, ab84Case.Ab84Bytes); ab84Actual != ab84Case.Ab84Complete {
			ab84T.Fatalf("logical constants %v: %v", ab84Case.Ab84Bytes, ab84Actual)
		}
	}
}

func ab84KeyMutations(ab84Control Ab84Key) []Ab84Key {
	ab84Keys := []Ab84Key{}
	for ab84Index := 0; ab84Index < 11; ab84Index++ {
		ab84Changed := ab84Control
		switch ab84Index {
		case 0:
			ab84Changed.Ab84Family += "_other"
		case 1:
			ab84Changed.Ab84Variant += "_other"
		case 2:
			ab84Changed.Ab84Primary.Ab84Value++
		case 3:
			ab84Changed.Ab84Secondary.Ab84Value++
		case 4:
			ab84Changed.Ab84Operation.Ab84Value++
		case 5:
			ab84Changed.Ab84Direction += "_other"
		case 6:
			ab84Changed.Ab84Field += "_other"
		case 7:
			ab84Changed.Ab84Opcode.Ab84Value++
		case 8:
			ab84Changed.Ab84Group.Ab84Value++
		case 9:
			ab84Changed.Ab84Instance.Ab84Value++
		case 10:
			ab84Changed.Ab84Register.Ab84Value++
		}
		ab84Keys = append(ab84Keys, ab84Changed)
	}
	for ab84Index := 0; ab84Index < 7; ab84Index++ {
		ab84Changed := ab84Control
		ab84Axes := []*Ab84Axis{&ab84Changed.Ab84Primary, &ab84Changed.Ab84Secondary, &ab84Changed.Ab84Operation, &ab84Changed.Ab84Opcode, &ab84Changed.Ab84Group, &ab84Changed.Ab84Instance, &ab84Changed.Ab84Register}
		ab84Axes[ab84Index].Ab84Present = !ab84Axes[ab84Index].Ab84Present
		ab84Keys = append(ab84Keys, ab84Changed)
	}
	return ab84Keys
}

func TestAb84ExactKeys(ab84T *ab84testing.T) {
	ab84Control := ab84Key("heat_state", "request")
	if _, ab84OK := Ab84LookupBurnerField(ab84Control); !ab84OK {
		ab84T.Fatal("valid control missing")
	}
	for ab84Index, ab84Changed := range ab84KeyMutations(ab84Control) {
		if _, ab84OK := Ab84LookupBurnerField(ab84Changed); ab84OK {
			ab84T.Fatal("identity axis ignored", ab84Index)
		}
	}
	ab84Standard := ab84Control
	ab84Standard.Ab84Secondary.Ab84Value = 4
	ab84Standard.Ab84Direction = "response"
	if _, ab84OK := Ab84LookupBurnerField(ab84Standard); ab84OK {
		ab84T.Fatal("standard 05/04 aliases candidate")
	}
	// Synthetic keys exercise the same lookup without declaring B524 semantics.
	ab84B524 := Ab84Key{Ab84Family: "synthetic_family", Ab84Variant: "synthetic_variant", Ab84Direction: "response", Ab84Field: "synthetic_field", Ab84Primary: Ab84Axis{true, 0xB5}, Ab84Secondary: Ab84Axis{true, 0x24}, Ab84Operation: Ab84Axis{true, 0}, Ab84Opcode: Ab84Axis{true, 0}, Ab84Group: Ab84Axis{true, 0}, Ab84Instance: Ab84Axis{true, 0}, Ab84Register: Ab84Axis{true, 0}}
	for _, ab84Control := range []Ab84Key{ab84B524, {Ab84Family: "synthetic_family", Ab84Variant: "synthetic_variant", Ab84Direction: "response", Ab84Field: "synthetic_field"}} {
		ab84Fields := []Ab84BurnerField{{Ab84Key: ab84Control, Ab84Grade: "Unknown"}}
		if _, ab84OK := ab84LookupBurnerField(ab84Fields, ab84Control); !ab84OK {
			ab84T.Fatal("synthetic control missing")
		}
		for ab84Index, ab84Changed := range ab84KeyMutations(ab84Control) {
			if _, ab84OK := ab84LookupBurnerField(ab84Fields, ab84Changed); ab84OK {
				ab84T.Fatal("synthetic axis ignored", ab84Index)
			}
		}
	}
}

func ab84Mutate(ab84Field *Ab84BurnerField) {
	ab84Field.Ab84Key.Ab84Family = "mutated"
	if ab84Field.Ab84LogicalPosition != nil {
		*ab84Field.Ab84LogicalPosition = 999
	}
	if ab84Field.Ab84WidthBits != nil {
		*ab84Field.Ab84WidthBits = 999
	}
	if ab84Field.Ab84Default != nil {
		*ab84Field.Ab84Default = 999
	}
	if ab84Field.Ab84Constant != nil {
		*ab84Field.Ab84Constant = 999
	}
	if ab84Field.Ab84Bounds != nil {
		ab84Field.Ab84Bounds.Ab84Min = 999
		ab84Field.Ab84Bounds.Ab84Max = 999
	}
	if ab84Field.Ab84Enumeration != nil {
		ab84Field.Ab84Enumeration[0] = "mutated"
	}
	ab84Field.Ab84Replacement = append(ab84Field.Ab84Replacement, 999)
	ab84Field.Ab84Reserved = append(ab84Field.Ab84Reserved, 999)
}

func TestAb84DetachedCopies(ab84T *ab84testing.T) {
	ab84Bodies := Ab84BurnerMetadata()
	for ab84Index := range ab84Bodies.Ab84Fields {
		ab84Mutate(&ab84Bodies.Ab84Fields[ab84Index])
	}
	if !ab84reflect.DeepEqual(Ab84BurnerMetadata().Ab84Fields, ab84ExpectedFields()) {
		ab84T.Fatal("accessor storage shared")
	}
	for _, ab84Expected := range ab84ExpectedFields() {
		ab84Field, ab84OK := Ab84LookupBurnerField(ab84Expected.Ab84Key)
		if !ab84OK {
			ab84T.Fatal("lookup missing")
		}
		ab84Mutate(&ab84Field)
		ab84Again, ab84OK := Ab84LookupBurnerField(ab84Expected.Ab84Key)
		if !ab84OK || !ab84reflect.DeepEqual(ab84Again, ab84Expected) {
			ab84T.Fatal("lookup storage shared", ab84Expected.Ab84Key)
		}
	}
}

func TestAb84ProviderUnchanged(ab84T *ab84testing.T) {
	ab84Provider, ab84Error := NewProvider(MustEmbeddedCatalog(), true)
	if ab84Error != nil {
		ab84T.Fatal(ab84Error)
	}
	_ = Ab84BurnerMetadata()
	for _, ab84Caller := range []CallerContext{CallerContextUserFacing, CallerContextSystemNMRuntime} {
		if _, ab84Error := ab84Provider.Invoke(ab84context.Background(), "ebus_standard.burner.op_data_c2b", nil, ab84Caller); !ab84errors.Is(ab84Error, ErrSafetyClassDenied) {
			ab84T.Fatal("burner became callable", ab84Error)
		}
	}
	ab84Result, ab84Error := ab84Provider.Invoke(ab84context.Background(), "ebus_standard.system_data.identification_addressed", nil, CallerContextUserFacing)
	if ab84Error != nil {
		ab84T.Fatal(ab84Error)
	}
	if ab84Executed, ab84OK := ab84Result["executed"].(bool); !ab84OK || ab84Executed {
		ab84T.Fatal("existing stub executed", ab84Result)
	}
	ab84Identity, ab84Error := ab84Provider.Identification(ab84context.Background())
	if ab84Error != nil || ab84Identity.Valid || ab84Identity.Manufacturer != "" || ab84Identity.DeviceID != "" {
		ab84T.Fatal("identification gained evidence", ab84Identity, ab84Error)
	}
}
