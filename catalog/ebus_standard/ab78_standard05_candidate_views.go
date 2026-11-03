package ebus_standard_catalog

type Ab78State string

type Ab78Key struct {
	Ab78Family    string
	Ab78Variant   string
	Ab78Primary   *uint8
	Ab78Secondary *uint8
	Ab78Operation *uint8
	Ab78Direction string
	Ab78Selector  string
}

type Ab78Field struct {
	Ab78Label          string
	Ab78Raw            *uint8
	Ab78Value          *int
	Ab78Name           string
	Ab78Interpretation Ab78State
	Ab78Domain         Ab78State
	Ab78Default        int
	Ab78Minimum        int
	Ab78Maximum        int
}

type Ab78View struct {
	Ab78Key          Ab78Key
	Ab78Raw          []byte
	Ab78Remainder    []byte
	Ab78State        Ab78State
	Ab78Complete     bool
	Ab78Fields       [5]Ab78Field
	Ab78APIGrade     string
	Ab78MeaningGrade string
	Ab78NativeGrade  string
	Ab78Unknowns     []string
}

func ab78CopyByte(ab78Byte *uint8) *uint8 {
	if ab78Byte == nil {
		return nil
	}
	ab78Value := *ab78Byte
	return &ab78Value
}

func Ab78Interpret(ab78Key Ab78Key, ab78Body []byte) Ab78View {
	ab78Key.Ab78Primary = ab78CopyByte(ab78Key.Ab78Primary)
	ab78Key.Ab78Secondary = ab78CopyByte(ab78Key.Ab78Secondary)
	ab78Key.Ab78Operation = ab78CopyByte(ab78Key.Ab78Operation)
	ab78View := Ab78View{
		Ab78Key:       ab78Key,
		Ab78Raw:       append([]byte(nil), ab78Body...),
		Ab78Remainder: append([]byte(nil), ab78Body...),
		Ab78State:     "raw-only",
		Ab78APIGrade:  "Proven", Ab78MeaningGrade: "Hypothesis", Ab78NativeGrade: "Unknown",
		Ab78Unknowns: []string{"physical support", "presence", "freshness", "physical length", "delivery", "effects", "persistence", "execution", "ACK", "command acceptance", "empty wire response"},
	}
	if ab78Key.Ab78Family != "burner" || ab78Key.Ab78Variant != "0501" ||
		ab78Key.Ab78Primary == nil || *ab78Key.Ab78Primary != 5 ||
		ab78Key.Ab78Secondary == nil || *ab78Key.Ab78Secondary != 1 ||
		ab78Key.Ab78Operation != nil || ab78Key.Ab78Direction != "abstract-request" ||
		ab78Key.Ab78Selector != "burner-five-field" {
		return ab78View
	}
	ab78View.Ab78State = "candidate-body"
	ab78View.Ab78Complete = len(ab78Body) == 5
	ab78View.Ab78Remainder = nil
	if len(ab78Body) > 5 {
		ab78View.Ab78Remainder = append([]byte(nil), ab78Body[5:]...)
	}
	ab78View.Ab78Fields = [5]Ab78Field{
		{Ab78Label: "heat-state", Ab78Maximum: 255},
		{Ab78Label: "boiler-target", Ab78Maximum: 90},
		{Ab78Label: "DHW-target", Ab78Maximum: 90},
		{Ab78Label: "outside-temperature", Ab78Minimum: -128, Ab78Maximum: 127},
		{Ab78Label: "setting", Ab78Maximum: 255},
	}
	for ab78Index := range ab78View.Ab78Fields {
		ab78Field := &ab78View.Ab78Fields[ab78Index]
		ab78Field.Ab78Interpretation = "absent"
		ab78Field.Ab78Domain = "unknown"
		if ab78Index >= len(ab78Body) {
			continue
		}
		ab78Raw := ab78Body[ab78Index]
		ab78Field.Ab78Raw = &ab78Raw
		ab78Value := int(ab78Raw)
		if ab78Index == 3 {
			ab78Value = int(int8(ab78Raw))
		}
		ab78Field.Ab78Interpretation = "interpreted"
		ab78Field.Ab78Domain = "within-candidate-domain"
		if ab78Index == 0 {
			switch ab78Raw {
			case 0x00:
				ab78Field.Ab78Name = "shutdown"
			case 0x55:
				ab78Field.Ab78Name = "DHW"
			case 0xaa:
				ab78Field.Ab78Name = "heating"
			case 0xcc:
				ab78Field.Ab78Name = "emission-check"
			case 0xee:
				ab78Field.Ab78Name = "controller-stop"
			default:
				ab78Field.Ab78Interpretation = "unmapped"
				ab78Field.Ab78Domain = "unknown"
				continue
			}
		}
		ab78Field.Ab78Value = &ab78Value
		if ab78Value < ab78Field.Ab78Minimum || ab78Value > ab78Field.Ab78Maximum {
			ab78Field.Ab78Domain = "out-of-candidate-domain"
		}
	}
	return ab78View
}
