package vaillant

import (
	ab89reflect "reflect"
	ab89testing "testing"
)

func ab89Int(ab89Value int) *int          { return &ab89Value }
func ab89String(ab89Value string) *string { return &ab89Value }
func ab89Axis(ab89Value uint16) Ab89Axis  { return Ab89Axis{true, ab89Value} }
func ab89Key(ab89Operation uint16, ab89Direction, ab89Name, ab89Variant string) Ab89Key {
	ab89Primary := uint16(0xb5)
	if ab89Direction == "positive" {
		ab89Primary = 0xf5
	}
	return Ab89Key{Ab89Family: "B514", Ab89Variant: ab89Variant, Ab89Primary: ab89Primary, Ab89Secondary: ab89Axis(0x14), Ab89Operation: ab89Axis(ab89Operation), Ab89Direction: ab89Direction, Ab89Field: ab89Name}
}
func ab89Get(ab89T *ab89testing.T, ab89Key Ab89Key) Ab89Field {
	ab89T.Helper()
	ab89Field, ab89Found, ab89Error := Ab89Lookup(ab89Key)
	if ab89Error != nil || !ab89Found {
		ab89T.Fatalf("lookup %+v: found=%v error=%v", ab89Key, ab89Found, ab89Error)
	}
	return ab89Field
}
func ab89Label(ab89Field Ab89Field, ab89Value int) string {
	if ab89Label, ab89Found := ab89Field.Ab89Enumeration[ab89Value]; ab89Found {
		return ab89Label
	}
	for _, ab89Range := range ab89Field.Ab89Ranges {
		if ab89Value >= ab89Range.Ab89Minimum && ab89Value <= ab89Range.Ab89Maximum {
			return ab89Range.Ab89Meaning
		}
	}
	return "unmapped"
}
func ab89Word(ab89T *ab89testing.T, ab89Field Ab89Field, ab89Data []byte) float64 {
	ab89T.Helper()
	if ab89Field.Ab89Position == nil || ab89Field.Ab89Width == nil || ab89Field.Ab89Extent == nil || len(ab89Data) < *ab89Field.Ab89Extent {
		ab89T.Fatal("incomplete abstract object")
	}
	if *ab89Field.Ab89Width != 2 || ab89Field.Ab89ByteOrder != "little-endian" {
		ab89T.Fatal("unexpected word layout")
	}
	ab89Position := *ab89Field.Ab89Position
	ab89Value := float64(uint16(ab89Data[ab89Position]) | uint16(ab89Data[ab89Position+1])<<8)
	if ab89Field.Ab89Conversion != nil {
		ab89Value = ab89Value*ab89Field.Ab89Conversion.Ab89Factor + ab89Field.Ab89Conversion.Ab89Offset
	}
	return ab89Value
}

func TestAb89ProductionMetadata(ab89T *ab89testing.T) {
	ab89Catalog, ab89Error := Ab89MenuFeedback()
	if ab89Error != nil {
		ab89T.Fatal(ab89Error)
	}
	if ab89Catalog.Ab89RepresentationGrade != "Proven" || ab89Catalog.Ab89Qualification != "Unknown" || len(ab89Catalog.Ab89Fields) != 21 || len(ab89Catalog.Ab89Limitations) != 8 {
		ab89T.Fatal("catalog boundary changed")
	}
	ab89Expected := []Ab89Field{}
	ab89Add := func(ab89Op uint16, ab89Direction, ab89Name, ab89Variant string, ab89Extent *int, ab89Enum map[int]string) *Ab89Field {
		ab89Expected = append(ab89Expected, Ab89Field{Ab89Key: ab89Key(ab89Op, ab89Direction, ab89Name, ab89Variant), Ab89Extent: ab89Extent, Ab89Signedness: "Unknown", Ab89ByteOrder: "Unknown", Ab89Enumeration: ab89Enum, Ab89Ranges: []Ab89Range{}, Ab89Reserved: []Ab89Range{}, Ab89Grade: "Hypothesis", Ab89RawFallback: "preserve raw bytes and unresolved alternatives"})
		return &ab89Expected[len(ab89Expected)-1]
	}
	ab89Add(1, "request", "mode", "menu", ab89Int(4), map[int]string{0: "off", 1: "on"}).Ab89Ranges = []Ab89Range{{2, 255, "off"}}
	ab89Add(1, "positive", "answer", "menu", ab89Int(5), map[int]string{0: "normal", 1: "actors-deactivated"})
	ab89Add(1, "positive", "feedback", "menu", ab89Int(5), map[int]string{1: "accepted"})
	for _, ab89Variant := range []string{"generic", "specific-a", "specific-b"} {
		ab89Add(6, "request", "command", ab89Variant, nil, map[int]string{0: "query", 1: "start"})
		if ab89Variant == "generic" {
			ab89Add(6, "positive", "feedback", ab89Variant, nil, map[int]string{1: "in-progress", 255: "unsupported"})
		} else {
			ab89Add(6, "positive", "feedback", ab89Variant, nil, map[int]string{1: "not-successful"}).Ab89Ranges = []Ab89Range{{3, 255, "invalid"}}
		}
	}
	ab89Add(7, "request", "command", "menu", nil, map[int]string{0: "query", 1: "start-liquefied-gas", 2: "start-natural-gas", 3: "abort"}).Ab89Default = ab89Int(0)
	ab89Add(7, "positive", "feedback", "menu", nil, map[int]string{0: "idle", 1: "progress", 2: "passed", 3: "warning", 4: "combustion-error", 5: "gas-type-error", 6: "aborted", 255: "unsupported"}).Ab89Default = ab89Int(0)
	for _, ab89Case := range []struct {
		ab89Direction, ab89Name             string
		ab89Position, ab89Width, ab89Extent int
	}{
		{"request", "command", 3, 1, 6}, {"request", "pressure", 4, 2, 6}, {"positive", "feedback", 3, 1, 13}, {"positive", "pressure", 4, 2, 13}, {"positive", "minimum", 6, 2, 13}, {"positive", "maximum", 8, 2, 13}, {"positive", "step", 10, 1, 13}, {"positive", "flow", 11, 2, 13},
	} {
		ab89Field := ab89Add(8, ab89Case.ab89Direction, ab89Case.ab89Name, "menu", ab89Int(ab89Case.ab89Extent), map[int]string{})
		ab89Field.Ab89Position = ab89Int(ab89Case.ab89Position)
		ab89Field.Ab89Width = ab89Int(ab89Case.ab89Width)
		ab89Field.Ab89Signedness = "unsigned"
		ab89Field.Ab89ByteOrder = "not applicable"
		if ab89Case.ab89Width == 2 {
			ab89Field.Ab89ByteOrder = "little-endian"
		}
		if ab89Case.ab89Name == "step" {
			ab89Field.Ab89Unit = ab89String("mbar")
			ab89Field.Ab89Default = ab89Int(0)
			ab89Field.Ab89Constant = ab89Int(1)
			ab89Field.Ab89Conflict = ab89String("constant 1 mbar conflicts with default zero; unresolved")
		}
		if ab89Case.ab89Name == "flow" {
			ab89Field.Ab89Unit = ab89String("l/min")
			ab89Field.Ab89Conversion = &Ab89Conversion{0.01, 0}
		}
	}
	for _, ab89Case := range []struct {
		ab89Direction string
		ab89Constant  int
	}{{"request", 0xb5}, {"positive", 0xf5}} {
		ab89Field := ab89Add(0, ab89Case.ab89Direction, "prefix", "named-09", nil, map[int]string{})
		ab89Field.Ab89Key.Ab89Operation = Ab89Axis{}
		ab89Field.Ab89Key.Ab89Secondary = Ab89Axis{}
		ab89Field.Ab89Position = ab89Int(0)
		ab89Field.Ab89Width = ab89Int(1)
		ab89Field.Ab89Signedness = "unsigned"
		ab89Field.Ab89ByteOrder = "not applicable"
		ab89Field.Ab89Constant = ab89Int(ab89Case.ab89Constant)
	}
	if !ab89reflect.DeepEqual(ab89Catalog.Ab89Fields, ab89Expected) {
		ab89T.Fatal("production fields differ from independent expected metadata")
	}
	for _, ab89Field := range ab89Expected {
		if !ab89reflect.DeepEqual(ab89Get(ab89T, ab89Field.Ab89Key), ab89Field) {
			ab89T.Fatal("lookup differs")
		}
	}
}

func TestAb89SyntheticCases(ab89T *ab89testing.T) {
	for _, ab89Case := range []struct {
		ab89Op                               uint16
		ab89Direction, ab89Name, ab89Variant string
		ab89Value                            int
		ab89Expected                         string
	}{
		{1, "request", "mode", "menu", 2, "off"}, {1, "positive", "answer", "menu", 1, "actors-deactivated"}, {1, "positive", "feedback", "menu", 1, "accepted"}, {1, "positive", "feedback", "menu", 0, "unmapped"},
		{6, "positive", "feedback", "generic", 1, "in-progress"}, {6, "positive", "feedback", "specific-a", 1, "not-successful"}, {6, "positive", "feedback", "specific-b", 1, "not-successful"}, {6, "positive", "feedback", "generic", 255, "unsupported"}, {6, "positive", "feedback", "specific-a", 255, "invalid"}, {6, "positive", "feedback", "specific-b", 255, "invalid"}, {6, "positive", "feedback", "specific-a", 2, "unmapped"},
		{7, "positive", "feedback", "menu", 6, "aborted"}, {7, "positive", "feedback", "menu", 7, "unmapped"},
	} {
		if ab89Label(ab89Get(ab89T, ab89Key(ab89Case.ab89Op, ab89Case.ab89Direction, ab89Case.ab89Name, ab89Case.ab89Variant)), ab89Case.ab89Value) != ab89Case.ab89Expected {
			ab89T.Fatalf("synthetic case %+v", ab89Case)
		}
	}
	ab89Data := []byte{0, 0, 0, 0, 1, 0, 0, 1, 0x34, 0x12, 0, 0x64, 0}
	for _, ab89Case := range []struct {
		ab89Name     string
		ab89Expected float64
	}{{"pressure", 1}, {"minimum", 256}, {"maximum", 0x1234}, {"flow", 1}} {
		ab89Field := ab89Get(ab89T, ab89Key(8, "positive", ab89Case.ab89Name, "menu"))
		if ab89Word(ab89T, ab89Field, ab89Data) != ab89Case.ab89Expected {
			ab89T.Fatal(ab89Case)
		}
		if len(ab89Data[:12]) >= *ab89Field.Ab89Extent {
			ab89T.Fatal("12-position object admitted")
		}
	}
	ab89Step := ab89Get(ab89T, ab89Key(8, "positive", "step", "menu"))
	if int(ab89Data[*ab89Step.Ab89Position]) != *ab89Step.Ab89Default || *ab89Step.Ab89Constant != 1 || ab89Step.Ab89Conflict == nil {
		ab89T.Fatal("step conflict lost")
	}
	ab89Named := ab89Key(0, "request", "prefix", "named-09")
	ab89Named.Ab89Operation = Ab89Axis{}
	ab89Named.Ab89Secondary = Ab89Axis{}
	ab89Field := ab89Get(ab89T, ab89Named)
	if ab89Field.Ab89Key.Ab89Secondary.Ab89Present || ab89Field.Ab89Key.Ab89Operation.Ab89Present || ab89Field.Ab89Extent != nil {
		ab89T.Fatal("inferred executable named 09")
	}
}

func TestAb89ExactIdentity(ab89T *ab89testing.T) {
	ab89Mutations := []struct {
		ab89Name   string
		ab89Mutate func(*Ab89Key)
	}{
		{"family", func(ab89Key *Ab89Key) { ab89Key.Ab89Family = "other" }},
		{"variant", func(ab89Key *Ab89Key) { ab89Key.Ab89Variant = "other" }},
		{"primary", func(ab89Key *Ab89Key) { ab89Key.Ab89Primary++ }},
		{"direction", func(ab89Key *Ab89Key) { ab89Key.Ab89Direction = "other" }},
		{"field", func(ab89Key *Ab89Key) { ab89Key.Ab89Field = "other" }},
	}
	for _, ab89AxisName := range []string{"Ab89Secondary", "Ab89Operation", "Ab89Opcode", "Ab89Group", "Ab89Instance", "Ab89Register"} {
		ab89Mutations = append(ab89Mutations, struct {
			ab89Name   string
			ab89Mutate func(*Ab89Key)
		}{ab89AxisName + "-presence", func(ab89Key *Ab89Key) {
			ab89AxisValue := ab89reflect.ValueOf(ab89Key).Elem().FieldByName(ab89AxisName).Addr().Interface().(*Ab89Axis)
			ab89AxisValue.Ab89Present = !ab89AxisValue.Ab89Present
		}}, struct {
			ab89Name   string
			ab89Mutate func(*Ab89Key)
		}{ab89AxisName + "-value", func(ab89Key *Ab89Key) {
			ab89AxisValue := ab89reflect.ValueOf(ab89Key).Elem().FieldByName(ab89AxisName).Addr().Interface().(*Ab89Axis)
			ab89AxisValue.Ab89Value ^= 0x8000
		}})
	}
	ab89Catalog, ab89Error := Ab89MenuFeedback()
	if ab89Error != nil {
		ab89T.Fatal(ab89Error)
	}
	for _, ab89Control := range ab89Catalog.Ab89Fields {
		for _, ab89Mutation := range ab89Mutations {
			ab89T.Run(ab89Control.Ab89Key.Ab89Variant+"/"+ab89Control.Ab89Key.Ab89Field+"/"+ab89Mutation.ab89Name, func(ab89T *ab89testing.T) {
				ab89Get(ab89T, ab89Control.Ab89Key)
				ab89Changed := ab89Control.Ab89Key
				ab89Mutation.ab89Mutate(&ab89Changed)
				if _, ab89Found, ab89Error := Ab89Lookup(ab89Changed); ab89Found || ab89Error != nil {
					ab89T.Fatal("nonexact key matched", ab89Changed, ab89Error)
				}
			})
		}
	}
	ab89B524 := ab89Key(0, "request", "synthetic", "synthetic")
	ab89B524.Ab89Family = "B524"
	ab89B524.Ab89Primary = 0xb5
	ab89B524.Ab89Secondary = ab89Axis(0)
	ab89B524.Ab89Opcode = ab89Axis(3)
	ab89B524.Ab89Group = ab89Axis(4)
	ab89B524.Ab89Instance = ab89Axis(5)
	ab89B524.Ab89Register = ab89Axis(6)
	ab89Zero := ab89B524
	ab89Zero.Ab89Opcode = ab89Axis(0)
	ab89Zero.Ab89Group = ab89Axis(0)
	ab89Zero.Ab89Instance = ab89Axis(0)
	ab89Zero.Ab89Register = ab89Axis(0)
	for _, ab89Control := range []Ab89Key{ab89B524, ab89Zero} {
		ab89Synthetic := []Ab89Field{{Ab89Key: ab89Control}}
		if _, ab89Found, _ := ab89LookupFields(ab89Synthetic, ab89Control); !ab89Found {
			ab89T.Fatal("synthetic B524 control missing")
		}
		for _, ab89Mutation := range ab89Mutations {
			ab89Changed := ab89Control
			ab89Mutation.ab89Mutate(&ab89Changed)
			if _, ab89Found, _ := ab89LookupFields(ab89Synthetic, ab89Changed); ab89Found {
				ab89T.Fatal("B524 axis ignored", ab89Mutation.ab89Name)
			}
		}
		if _, ab89Found, _ := Ab89Lookup(ab89Control); ab89Found {
			ab89T.Fatal("synthetic B524 became production declaration")
		}
	}

}

func ab89MutateField(ab89Field *Ab89Field) {
	if ab89Field.Ab89Position != nil {
		*ab89Field.Ab89Position = 99
	}
	if ab89Field.Ab89Width != nil {
		*ab89Field.Ab89Width = 99
	}
	if ab89Field.Ab89Extent != nil {
		*ab89Field.Ab89Extent = 99
	}
	if ab89Field.Ab89Unit != nil {
		*ab89Field.Ab89Unit = "mutated"
	}
	if ab89Field.Ab89Default != nil {
		*ab89Field.Ab89Default = 99
	}
	if ab89Field.Ab89Constant != nil {
		*ab89Field.Ab89Constant = 99
	}
	if ab89Field.Ab89Conflict != nil {
		*ab89Field.Ab89Conflict = "mutated"
	}
	if ab89Field.Ab89Conversion != nil {
		ab89Field.Ab89Conversion.Ab89Factor = 99
		ab89Field.Ab89Conversion.Ab89Offset = 99
	}
	ab89Field.Ab89Enumeration[999] = "mutated"
	if len(ab89Field.Ab89Ranges) > 0 {
		ab89Field.Ab89Ranges[0].Ab89Meaning = "mutated"
	}
	ab89Field.Ab89Ranges = append(ab89Field.Ab89Ranges, Ab89Range{99, 99, "mutated"})
	ab89Field.Ab89Reserved = append(ab89Field.Ab89Reserved, Ab89Range{99, 99, "mutated"})
	ab89Field.Ab89Bounds = append(ab89Field.Ab89Bounds, 99)
	ab89Field.Ab89Key.Ab89Family = "mutated"
}

func TestAb89DetachedCopies(ab89T *ab89testing.T) {
	ab89Original, ab89Error := Ab89MenuFeedback()
	if ab89Error != nil {
		ab89T.Fatal(ab89Error)
	}
	ab89Changed, ab89Error := Ab89MenuFeedback()
	if ab89Error != nil {
		ab89T.Fatal(ab89Error)
	}
	ab89Changed.Ab89Limitations[0] = "mutated"
	for ab89Index := range ab89Changed.Ab89Fields {
		ab89MutateField(&ab89Changed.Ab89Fields[ab89Index])
	}
	ab89Fresh, ab89Error := Ab89MenuFeedback()
	if ab89Error != nil || !ab89reflect.DeepEqual(ab89Fresh, ab89Original) {
		ab89T.Fatal("accessor shares backing storage", ab89Error)
	}
	for _, ab89Field := range ab89Original.Ab89Fields {
		ab89Lookup := ab89Get(ab89T, ab89Field.Ab89Key)
		ab89MutateField(&ab89Lookup)
		if !ab89reflect.DeepEqual(ab89Get(ab89T, ab89Field.Ab89Key), ab89Field) {
			ab89T.Fatal("lookup shares backing storage")
		}
	}
}
