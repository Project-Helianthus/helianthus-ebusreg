package ebus_standard_catalog

import (
	ab78json "encoding/json"
	ab78os "os"
	ab78reflect "reflect"
	ab78testing "testing"
)

func ab78Byte(ab78Value uint8) *uint8 { return &ab78Value }

func ab78Control() Ab78Key {
	return Ab78Key{Ab78Family: "burner", Ab78Variant: "0501", Ab78Primary: ab78Byte(5), Ab78Secondary: ab78Byte(1), Ab78Direction: "abstract-request", Ab78Selector: "burner-five-field"}
}

func TestAb78Examples(ab78T *ab78testing.T) {
	ab78Data, ab78Err := ab78os.ReadFile("testdata/ab78_standard05_candidate_views.json")
	if ab78Err != nil {
		ab78T.Fatal(ab78Err)
	}
	var ab78Cases []struct {
		Ab78Label   string
		Ab78Body    []int
		Ab78Values  []*int
		Ab78Name    string
		Ab78Domains []Ab78State
	}
	if ab78Err := ab78json.Unmarshal(ab78Data, &ab78Cases); ab78Err != nil {
		ab78T.Fatal(ab78Err)
	}
	if len(ab78Cases) != 7 {
		ab78T.Fatal("ab78 fixture inventory")
	}
	for _, ab78Case := range ab78Cases {
		ab78T.Run(ab78Case.Ab78Label, func(ab78T *ab78testing.T) {
			ab78Body := make([]byte, len(ab78Case.Ab78Body))
			for ab78Index, ab78Value := range ab78Case.Ab78Body {
				ab78Body[ab78Index] = byte(ab78Value)
			}
			ab78View := Ab78Interpret(ab78Control(), ab78Body)
			if !ab78View.Ab78Complete || ab78View.Ab78State != "candidate-body" || len(ab78View.Ab78Remainder) != 0 || !ab78reflect.DeepEqual(ab78View.Ab78Raw, ab78Body) {
				ab78T.Fatal("ab78 body boundary")
			}
			for ab78Index, ab78Field := range ab78View.Ab78Fields {
				ab78WantState := Ab78State("interpreted")
				if ab78Case.Ab78Values[ab78Index] == nil {
					ab78WantState = "unmapped"
				}
				if ab78Field.Ab78Raw == nil || *ab78Field.Ab78Raw != ab78Body[ab78Index] || !ab78reflect.DeepEqual(ab78Field.Ab78Value, ab78Case.Ab78Values[ab78Index]) || ab78Field.Ab78Domain != ab78Case.Ab78Domains[ab78Index] || ab78Field.Ab78Interpretation != ab78WantState {
					ab78T.Fatalf("ab78 field %d: %+v", ab78Index, ab78Field)
				}
			}
			if ab78View.Ab78Fields[0].Ab78Name != ab78Case.Ab78Name {
				ab78T.Fatal("ab78 heat name")
			}
		})
	}
}

func TestAb78SelectionAxes(ab78T *ab78testing.T) {
	ab78Cases := []struct {
		Ab78Label  string
		Ab78Mutate func(*Ab78Key)
	}{
		{"ab78-family", func(ab78Key *Ab78Key) { ab78Key.Ab78Family = "unknown" }},
		{"ab78-variant-unknown", func(ab78Key *Ab78Key) { ab78Key.Ab78Variant = "unknown" }},
		{"ab78-variant", func(ab78Key *Ab78Key) { ab78Key.Ab78Variant = "0504" }},
		{"ab78-primary", func(ab78Key *Ab78Key) { ab78Key.Ab78Primary = ab78Byte(4) }},
		{"ab78-primary-absent", func(ab78Key *Ab78Key) { ab78Key.Ab78Primary = nil }},
		{"ab78-primary-zero", func(ab78Key *Ab78Key) { ab78Key.Ab78Primary = ab78Byte(0) }},
		{"ab78-secondary", func(ab78Key *Ab78Key) { ab78Key.Ab78Secondary = ab78Byte(4) }},
		{"ab78-secondary-absent", func(ab78Key *Ab78Key) { ab78Key.Ab78Secondary = nil }},
		{"ab78-secondary-zero", func(ab78Key *Ab78Key) { ab78Key.Ab78Secondary = ab78Byte(0) }},
		{"ab78-operation-zero", func(ab78Key *Ab78Key) { ab78Key.Ab78Operation = ab78Byte(0) }},
		{"ab78-operation", func(ab78Key *Ab78Key) { ab78Key.Ab78Operation = ab78Byte(1) }},
		{"ab78-positive", func(ab78Key *Ab78Key) { ab78Key.Ab78Direction = "abstract-positive" }},
		{"ab78-direction", func(ab78Key *Ab78Key) { ab78Key.Ab78Direction = "unknown" }},
		{"ab78-selector", func(ab78Key *Ab78Key) { ab78Key.Ab78Selector = "unknown" }},
		{"ab78-selector-absent", func(ab78Key *Ab78Key) { ab78Key.Ab78Selector = "" }},
	}
	for _, ab78Case := range ab78Cases {
		ab78T.Run(ab78Case.Ab78Label, func(ab78T *ab78testing.T) {
			ab78Key := ab78Control()
			if Ab78Interpret(ab78Key, []byte{0, 0, 0, 0, 0}).Ab78State != "candidate-body" {
				ab78T.Fatal("ab78 valid control")
			}
			ab78Case.Ab78Mutate(&ab78Key)
			for _, ab78Body := range [][]byte{nil, {5, 1}, {170, 90, 35, 255, 255}, {5, 1, 170, 90, 35, 255, 255}} {
				ab78View := Ab78Interpret(ab78Key, ab78Body)
				if ab78View.Ab78State != "raw-only" || ab78View.Ab78Complete || ab78View.Ab78Fields != [5]Ab78Field{} || !ab78reflect.DeepEqual(ab78View.Ab78Raw, ab78Body) || !ab78reflect.DeepEqual(ab78View.Ab78Remainder, ab78Body) {
					ab78T.Fatal("ab78 changed axis interpreted")
				}
			}
		})
	}
	ab78Key := ab78Control()
	ab78Key.Ab78Variant = "0504"
	ab78Key.Ab78Secondary = ab78Byte(4)
	ab78View := Ab78Interpret(ab78Key, []byte{0, 0, 0, 0, 0})
	if ab78View.Ab78State != "raw-only" || ab78View.Ab78Key.Ab78Variant != "0504" || *ab78View.Ab78Key.Ab78Secondary != 4 {
		ab78T.Fatal("ab78 separate identity lost")
	}
}

func TestAb78Lengths(ab78T *ab78testing.T) {
	ab78Body := []byte{170, 90, 35, 255, 255, 17, 34}
	for ab78Length := 0; ab78Length <= len(ab78Body); ab78Length++ {
		ab78View := Ab78Interpret(ab78Control(), ab78Body[:ab78Length])
		if ab78View.Ab78Complete != (ab78Length == 5) {
			ab78T.Fatalf("ab78 completeness %d", ab78Length)
		}
		for ab78Index, ab78Field := range ab78View.Ab78Fields {
			if ab78Index >= ab78Length {
				if ab78Field.Ab78Raw != nil || ab78Field.Ab78Value != nil || ab78Field.Ab78Name != "" || ab78Field.Ab78Interpretation != "absent" || ab78Field.Ab78Domain != "unknown" {
					ab78T.Fatalf("ab78 filled absent %d", ab78Index)
				}
			} else if ab78Field.Ab78Raw == nil || *ab78Field.Ab78Raw != ab78Body[ab78Index] || ab78Field.Ab78Value == nil || *ab78Field.Ab78Value != []int{170, 90, 35, -1, 255}[ab78Index] {
				ab78T.Fatal("ab78 partial value lost")
			}
		}
		if ab78Length > 5 && !ab78reflect.DeepEqual(ab78View.Ab78Remainder, ab78Body[5:ab78Length]) {
			ab78T.Fatal("ab78 remainder lost")
		}
	}
	ab78Zero := Ab78Interpret(ab78Control(), []byte{0, 0, 0, 0, 0})
	for _, ab78Field := range ab78Zero.Ab78Fields {
		if ab78Field.Ab78Raw == nil || ab78Field.Ab78Value == nil || *ab78Field.Ab78Value != 0 || ab78Field.Ab78Interpretation != "interpreted" {
			ab78T.Fatal("ab78 explicit zero absent")
		}
	}
}

func TestAb78FieldAxes(ab78T *ab78testing.T) {
	for ab78Index := 0; ab78Index < 5; ab78Index++ {
		for ab78Value := 0; ab78Value <= 255; ab78Value++ {
			ab78Body := []byte{0, 0, 0, 0, 0}
			ab78Body[ab78Index] = byte(ab78Value)
			ab78View := Ab78Interpret(ab78Control(), ab78Body)
			for ab78Other, ab78Field := range ab78View.Ab78Fields {
				if ab78Other != ab78Index && (ab78Field.Ab78Value == nil || *ab78Field.Ab78Value != 0) {
					ab78T.Fatal("ab78 field interference")
				}
			}
			ab78Field := ab78View.Ab78Fields[ab78Index]
			ab78Want := ab78Value
			if ab78Index == 3 && ab78Value >= 128 {
				ab78Want = ab78Value - 256
			}
			if ab78Index == 0 {
				ab78Names := map[int]string{0: "shutdown", 85: "DHW", 170: "heating", 204: "emission-check", 238: "controller-stop"}
				ab78Name, ab78Known := ab78Names[ab78Value]
				if !ab78Known {
					if ab78Field.Ab78Interpretation != "unmapped" || ab78Field.Ab78Value != nil || ab78Field.Ab78Name != "" || ab78Field.Ab78Domain != "unknown" {
						ab78T.Fatal("ab78 unmapped heat named")
					}
					continue
				}
				if ab78Field.Ab78Name != ab78Name {
					ab78T.Fatal("ab78 heat mapping")
				}
			}
			ab78Domain := Ab78State("within-candidate-domain")
			if (ab78Index == 1 || ab78Index == 2) && ab78Value > 90 {
				ab78Domain = "out-of-candidate-domain"
			}
			if ab78Field.Ab78Raw == nil || int(*ab78Field.Ab78Raw) != ab78Value || ab78Field.Ab78Value == nil || *ab78Field.Ab78Value != ab78Want || ab78Field.Ab78Domain != ab78Domain || ab78Field.Ab78Interpretation != "interpreted" {
				ab78T.Fatalf("ab78 axis %d value %d", ab78Index, ab78Value)
			}
		}
	}
}

func TestAb78Metadata(ab78T *ab78testing.T) {
	ab78View := Ab78Interpret(ab78Control(), nil)
	if ab78View.Ab78APIGrade != "Proven" || ab78View.Ab78MeaningGrade != "Hypothesis" || ab78View.Ab78NativeGrade != "Unknown" || !ab78reflect.DeepEqual(ab78View.Ab78Unknowns, []string{"physical support", "presence", "freshness", "physical length", "delivery", "effects", "persistence", "execution", "ACK", "command acceptance", "empty wire response"}) {
		ab78T.Fatal("ab78 grades or unknowns")
	}
	for ab78Index, ab78Field := range ab78View.Ab78Fields {
		if ab78Field.Ab78Label != []string{"heat-state", "boiler-target", "DHW-target", "outside-temperature", "setting"}[ab78Index] || ab78Field.Ab78Default != 0 || ab78Field.Ab78Minimum != []int{0, 0, 0, -128, 0}[ab78Index] || ab78Field.Ab78Maximum != []int{255, 90, 90, 127, 255}[ab78Index] {
			ab78T.Fatal("ab78 metadata")
		}
	}
}

func TestAb78OwnedMemory(ab78T *ab78testing.T) {
	ab78Key := ab78Control()
	ab78Body := []byte{170, 90, 35, 255, 255, 17}
	ab78Baseline := Ab78Interpret(ab78Control(), append([]byte(nil), ab78Body...))
	ab78View := Ab78Interpret(ab78Key, ab78Body)
	*ab78Key.Ab78Primary = 0
	*ab78Key.Ab78Secondary = 0
	ab78Body[0] = 0
	ab78Body[5] = 0
	if !ab78reflect.DeepEqual(ab78View, ab78Baseline) {
		ab78T.Fatal("ab78 input alias")
	}
	ab78View.Ab78Raw[0] = 1
	ab78View.Ab78Raw[5] = 2
	if *ab78View.Ab78Fields[0].Ab78Raw != 170 || ab78View.Ab78Remainder[0] != 17 {
		ab78T.Fatal("ab78 internal alias")
	}
	for ab78Index := range ab78View.Ab78Fields {
		*ab78View.Ab78Fields[ab78Index].Ab78Raw = 0
		*ab78View.Ab78Fields[ab78Index].Ab78Value = 0
	}
	ab78View.Ab78Unknowns[0] = "changed"
	*ab78View.Ab78Key.Ab78Primary = 0
	*ab78View.Ab78Key.Ab78Secondary = 0
	if !ab78reflect.DeepEqual(Ab78Interpret(ab78Control(), []byte{170, 90, 35, 255, 255, 17}), ab78Baseline) {
		ab78T.Fatal("ab78 return alias")
	}
	ab78Key = ab78Control()
	ab78Key.Ab78Operation = ab78Byte(0)
	ab78Raw := Ab78Interpret(ab78Key, ab78Body)
	*ab78Key.Ab78Operation = 1
	ab78Raw.Ab78Raw[0] = 9
	if *ab78Raw.Ab78Key.Ab78Operation != 0 || ab78Raw.Ab78Remainder[0] != ab78Body[0] {
		ab78T.Fatal("ab78 raw-only alias")
	}
}
