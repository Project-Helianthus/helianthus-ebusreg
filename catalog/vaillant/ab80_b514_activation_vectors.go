package vaillant

const (
	Ab80Family  = "b514-activation"
	Ab80Variant = "vector-candidate"

	Ab80OperationBits     = 3
	Ab80OperationScalars  = 4
	Ab80DirectionPositive = "positive"

	Ab80BitFieldCount    = 10
	Ab80ScalarFieldCount = 12

	Ab80GradeHypothesis = "Hypothesis"
	Ab80GradeUnknown    = "Unknown"

	Ab80StateRawOnly = "raw-only"
	Ab80StateBits    = "bits"
	Ab80StateScalars = "scalars"
)

type Ab80Identity struct {
	Ab80Family    string
	Ab80Variant   string
	Ab80Operation uint8
	Ab80Direction string
}

type Ab80Bit struct {
	Ab80Present bool
	Ab80Value   bool
}

type Ab80Scalar struct {
	Ab80Present        bool
	Ab80Raw            byte
	Ab80Interpretation string
}

type Ab80View struct {
	ab80State     string
	ab80Grade     string
	ab80Raw       []byte
	ab80Mask      uint16
	ab80Complete  bool
	ab80Remainder []byte
	ab80Fields    int
	ab80Scalars   bool
}

func ab80Copy(ab80Source []byte) []byte {
	ab80Out := make([]byte, len(ab80Source))
	copy(ab80Out, ab80Source)
	return ab80Out
}

func Ab80Interpret(ab80Identity Ab80Identity, ab80Vector []byte) Ab80View {
	ab80View := Ab80View{ab80State: Ab80StateRawOnly, ab80Grade: Ab80GradeUnknown, ab80Raw: ab80Copy(ab80Vector), ab80Remainder: []byte{}}
	if ab80Identity.Ab80Family != Ab80Family || ab80Identity.Ab80Variant != Ab80Variant || ab80Identity.Ab80Direction != Ab80DirectionPositive {
		return ab80View
	}
	switch ab80Identity.Ab80Operation {
	case Ab80OperationBits:
		ab80View.ab80State = Ab80StateBits
		ab80View.ab80Fields = Ab80BitFieldCount
	case Ab80OperationScalars:
		ab80View.ab80State = Ab80StateScalars
		ab80View.ab80Fields = Ab80ScalarFieldCount
		ab80View.ab80Scalars = true
	default:
		return ab80View
	}
	ab80View.ab80Grade = Ab80GradeHypothesis
	ab80Present := len(ab80Vector)
	if ab80Present > ab80View.ab80Fields {
		ab80Present = ab80View.ab80Fields
	}
	for ab80Index := 0; ab80Index < ab80Present; ab80Index++ {
		ab80View.ab80Mask |= 1 << uint(ab80Index)
	}
	ab80View.ab80Complete = len(ab80Vector) == ab80View.ab80Fields
	ab80View.ab80Remainder = ab80Copy(ab80Vector[ab80Present:])
	return ab80View
}

func (ab80View Ab80View) Ab80State() string { return ab80View.ab80State }

func (ab80View Ab80View) Ab80Grade() string { return ab80View.ab80Grade }

func (ab80View Ab80View) Ab80RawBytes() []byte { return ab80Copy(ab80View.ab80Raw) }

func (ab80View Ab80View) Ab80PresentMask() uint16 { return ab80View.ab80Mask }

func (ab80View Ab80View) Ab80IsComplete() bool { return ab80View.ab80Complete }

func (ab80View Ab80View) Ab80RemainderBytes() []byte { return ab80Copy(ab80View.ab80Remainder) }

func (ab80View Ab80View) Ab80FieldCount() int { return ab80View.ab80Fields }

func (ab80View Ab80View) Ab80BitAt(ab80Byte int, ab80Bit int) Ab80Bit {
	if ab80View.ab80State != Ab80StateBits || ab80Byte < 0 || ab80Byte >= ab80View.ab80Fields || ab80Bit < 0 || ab80Bit > 7 {
		return Ab80Bit{}
	}
	if ab80View.ab80Mask&(1<<uint(ab80Byte)) == 0 {
		return Ab80Bit{}
	}
	return Ab80Bit{Ab80Present: true, Ab80Value: ab80View.ab80Raw[ab80Byte]&(1<<uint(ab80Bit)) != 0}
}

func (ab80View Ab80View) Ab80ScalarAt(ab80Index int) Ab80Scalar {
	ab80Out := Ab80Scalar{Ab80Interpretation: Ab80GradeUnknown}
	if ab80View.ab80State != Ab80StateScalars || ab80Index < 0 || ab80Index >= ab80View.ab80Fields {
		return ab80Out
	}
	if ab80View.ab80Mask&(1<<uint(ab80Index)) == 0 {
		return ab80Out
	}
	ab80Out.Ab80Present = true
	ab80Out.Ab80Raw = ab80View.ab80Raw[ab80Index]
	return ab80Out
}
