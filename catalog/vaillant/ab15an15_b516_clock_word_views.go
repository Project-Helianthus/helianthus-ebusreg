package vaillant

const (
	Ab15An15Family  = "vaillant-clock-family"
	Ab15An15Variant = "b516-candidate"
)

type Ab15An15Selector struct {
	Ab15An15Family    string
	Ab15An15Variant   string
	Ab15An15Operation byte
	Ab15An15Direction string
}

type Ab15An15Cell struct {
	Ab15An15Raw    []byte
	Ab15An15Failed bool
}

type Ab15An15Field struct {
	Ab15An15State string
	Ab15An15Value *uint16
}

type Ab15An15Candidate struct {
	Ab15An15Admitted           bool
	Ab15An15ContractGrade      string
	Ab15An15FieldGrade         string
	Ab15An15NativeSupport      string
	Ab15An15PhysicalAdaptation string
	Ab15An15Presence           string
	Ab15An15Freshness          string
	Ab15An15Execution          string
}

type Ab15An15ClockInput struct {
	Ab15An15Clock           [5]Ab15An15Cell
	Ab15An15Setpoint        *byte
	Ab15An15Mode            *byte
	Ab15An15SpecialFunction *byte
	Ab15An15Messages        [5]*byte
}

type Ab15An15ClockResult struct {
	Ab15An15Candidate       Ab15An15Candidate
	Ab15An15Raw             Ab15An15ClockInput
	Ab15An15Clock           *[5]Ab15An15Field
	Ab15An15SetpointCelsius *float64
	Ab15An15TablesGrade     string
}

type Ab15An15WordResult struct {
	Ab15An15Candidate      Ab15An15Candidate
	Ab15An15Raw            [4]Ab15An15Cell
	Ab15An15Words          *[4]Ab15An15Field
	Ab15An15PositionsGrade string
}

func ab15an15Candidate(ab15an15Selector Ab15An15Selector, ab15an15Operation byte) Ab15An15Candidate {
	return Ab15An15Candidate{
		Ab15An15Admitted:      ab15an15Selector.Ab15An15Family == Ab15An15Family && ab15an15Selector.Ab15An15Variant == Ab15An15Variant && ab15an15Selector.Ab15An15Operation == ab15an15Operation && ab15an15Selector.Ab15An15Direction == "positive",
		Ab15An15ContractGrade: "Proven", Ab15An15FieldGrade: "Hypothesis", Ab15An15NativeSupport: "Unknown", Ab15An15PhysicalAdaptation: "Unknown", Ab15An15Presence: "Unknown", Ab15An15Freshness: "Unknown", Ab15An15Execution: "Unknown",
	}
}

func ab15an15CopyCell(ab15an15Cell Ab15An15Cell) Ab15An15Cell {
	if ab15an15Cell.Ab15An15Raw != nil {
		ab15an15Raw := make([]byte, len(ab15an15Cell.Ab15An15Raw))
		copy(ab15an15Raw, ab15an15Cell.Ab15An15Raw)
		ab15an15Cell.Ab15An15Raw = ab15an15Raw
	}
	return ab15an15Cell
}

func ab15an15CopyByte(ab15an15Byte *byte) *byte {
	if ab15an15Byte == nil {
		return nil
	}
	ab15an15Value := *ab15an15Byte
	return &ab15an15Value
}

func ab15an15Field(ab15an15Cell Ab15An15Cell, ab15an15Size int) Ab15An15Field {
	if ab15an15Cell.Ab15An15Failed {
		return Ab15An15Field{Ab15An15State: "failed"}
	}
	if len(ab15an15Cell.Ab15An15Raw) == 0 {
		return Ab15An15Field{Ab15An15State: "missing"}
	}
	if len(ab15an15Cell.Ab15An15Raw) != ab15an15Size {
		return Ab15An15Field{Ab15An15State: "failed"}
	}
	ab15an15Value := uint16(ab15an15Cell.Ab15An15Raw[0])
	if ab15an15Size == 2 {
		ab15an15Value |= uint16(ab15an15Cell.Ab15An15Raw[1]) << 8
	} else {
		if ab15an15Value == 255 {
			return Ab15An15Field{Ab15An15State: "replacement"}
		}
		if ab15an15Value>>4 > 9 || ab15an15Value&15 > 9 {
			return Ab15An15Field{Ab15An15State: "malformed-decimal"}
		}
		ab15an15Value = (ab15an15Value>>4)*10 + (ab15an15Value & 15)
	}
	return Ab15An15Field{Ab15An15State: "numeric", Ab15An15Value: &ab15an15Value}
}

func Ab15An15ClockView(ab15an15Selector Ab15An15Selector, ab15an15Input Ab15An15ClockInput) Ab15An15ClockResult {
	ab15an15Result := Ab15An15ClockResult{Ab15An15Candidate: ab15an15Candidate(ab15an15Selector, 0x08), Ab15An15Raw: ab15an15Input, Ab15An15TablesGrade: "Unknown"}
	for ab15an15Index := range ab15an15Input.Ab15An15Clock {
		ab15an15Result.Ab15An15Raw.Ab15An15Clock[ab15an15Index] = ab15an15CopyCell(ab15an15Input.Ab15An15Clock[ab15an15Index])
	}
	ab15an15Result.Ab15An15Raw.Ab15An15Setpoint = ab15an15CopyByte(ab15an15Input.Ab15An15Setpoint)
	ab15an15Result.Ab15An15Raw.Ab15An15Mode = ab15an15CopyByte(ab15an15Input.Ab15An15Mode)
	ab15an15Result.Ab15An15Raw.Ab15An15SpecialFunction = ab15an15CopyByte(ab15an15Input.Ab15An15SpecialFunction)
	for ab15an15Index := range ab15an15Input.Ab15An15Messages {
		ab15an15Result.Ab15An15Raw.Ab15An15Messages[ab15an15Index] = ab15an15CopyByte(ab15an15Input.Ab15An15Messages[ab15an15Index])
	}
	if !ab15an15Result.Ab15An15Candidate.Ab15An15Admitted {
		return ab15an15Result
	}
	ab15an15Result.Ab15An15Clock = new([5]Ab15An15Field)
	for ab15an15Index, ab15an15Cell := range ab15an15Input.Ab15An15Clock {
		ab15an15Result.Ab15An15Clock[ab15an15Index] = ab15an15Field(ab15an15Cell, 1)
	}
	if ab15an15Input.Ab15An15Setpoint != nil {
		ab15an15Value := float64(*ab15an15Input.Ab15An15Setpoint) / 2
		ab15an15Result.Ab15An15SetpointCelsius = &ab15an15Value
	}
	return ab15an15Result
}

func Ab15An15WordView(ab15an15Selector Ab15An15Selector, ab15an15Words [4]Ab15An15Cell) Ab15An15WordResult {
	ab15an15Result := Ab15An15WordResult{Ab15An15Candidate: ab15an15Candidate(ab15an15Selector, 0x0d), Ab15An15PositionsGrade: "Unknown"}
	for ab15an15Index, ab15an15Cell := range ab15an15Words {
		ab15an15Result.Ab15An15Raw[ab15an15Index] = ab15an15CopyCell(ab15an15Cell)
	}
	if !ab15an15Result.Ab15An15Candidate.Ab15An15Admitted {
		return ab15an15Result
	}
	ab15an15Result.Ab15An15Words = new([4]Ab15An15Field)
	for ab15an15Index, ab15an15Cell := range ab15an15Words {
		ab15an15Result.Ab15An15Words[ab15an15Index] = ab15an15Field(ab15an15Cell, 2)
	}
	return ab15an15Result
}
