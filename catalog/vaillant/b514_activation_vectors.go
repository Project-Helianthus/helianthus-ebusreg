package vaillant

const (
	B514ActivationVectorFamily  = "b514-activation"
	B514ActivationVectorVariant = "vector-candidate"

	B514ActivationVectorOperationBits     = 3
	B514ActivationVectorOperationScalars  = 4
	B514ActivationVectorDirectionPositive = "positive"

	B514ActivationVectorBitFieldCount    = 10
	B514ActivationVectorScalarFieldCount = 12

	B514ActivationVectorGradeHypothesis = "Hypothesis"
	B514ActivationVectorGradeUnknown    = "Unknown"

	B514ActivationVectorStateRawOnly = "raw-only"
	B514ActivationVectorStateBits    = "bits"
	B514ActivationVectorStateScalars = "scalars"
)

type B514ActivationVectorIdentity struct {
	B514ActivationVectorFamily    string `json:"Ab80Family"`
	B514ActivationVectorVariant   string `json:"Ab80Variant"`
	B514ActivationVectorOperation uint8  `json:"Ab80Operation"`
	B514ActivationVectorDirection string `json:"Ab80Direction"`
}

type B514ActivationVectorBit struct {
	B514ActivationVectorPresent bool `json:"Ab80Present"`
	B514ActivationVectorValue   bool `json:"Ab80Value"`
}

type B514ActivationVectorScalar struct {
	B514ActivationVectorPresent        bool   `json:"Ab80Present"`
	B514ActivationVectorRaw            byte   `json:"Ab80Raw"`
	B514ActivationVectorInterpretation string `json:"Ab80Interpretation"`
}

type B514ActivationVectorView struct {
	b514ActivationVectorState     string
	b514ActivationVectorGrade     string
	b514ActivationVectorRaw       []byte
	b514ActivationVectorMask      uint16
	b514ActivationVectorComplete  bool
	b514ActivationVectorRemainder []byte
	b514ActivationVectorFields    int
	b514ActivationVectorScalars   bool
}

func b514ActivationVectorCopy(b514ActivationVectorSource []byte) []byte {
	b514ActivationVectorOut := make([]byte, len(b514ActivationVectorSource))
	copy(b514ActivationVectorOut, b514ActivationVectorSource)
	return b514ActivationVectorOut
}

func B514ActivationVectorInterpret(b514ActivationVectorIdentity B514ActivationVectorIdentity, b514ActivationVectorVector []byte) B514ActivationVectorView {
	b514ActivationVectorView := B514ActivationVectorView{b514ActivationVectorState: B514ActivationVectorStateRawOnly, b514ActivationVectorGrade: B514ActivationVectorGradeUnknown, b514ActivationVectorRaw: b514ActivationVectorCopy(b514ActivationVectorVector), b514ActivationVectorRemainder: []byte{}}
	if b514ActivationVectorIdentity.B514ActivationVectorFamily != B514ActivationVectorFamily || b514ActivationVectorIdentity.B514ActivationVectorVariant != B514ActivationVectorVariant || b514ActivationVectorIdentity.B514ActivationVectorDirection != B514ActivationVectorDirectionPositive {
		return b514ActivationVectorView
	}
	switch b514ActivationVectorIdentity.B514ActivationVectorOperation {
	case B514ActivationVectorOperationBits:
		b514ActivationVectorView.b514ActivationVectorState = B514ActivationVectorStateBits
		b514ActivationVectorView.b514ActivationVectorFields = B514ActivationVectorBitFieldCount
	case B514ActivationVectorOperationScalars:
		b514ActivationVectorView.b514ActivationVectorState = B514ActivationVectorStateScalars
		b514ActivationVectorView.b514ActivationVectorFields = B514ActivationVectorScalarFieldCount
		b514ActivationVectorView.b514ActivationVectorScalars = true
	default:
		return b514ActivationVectorView
	}
	b514ActivationVectorView.b514ActivationVectorGrade = B514ActivationVectorGradeHypothesis
	b514ActivationVectorPresent := len(b514ActivationVectorVector)
	if b514ActivationVectorPresent > b514ActivationVectorView.b514ActivationVectorFields {
		b514ActivationVectorPresent = b514ActivationVectorView.b514ActivationVectorFields
	}
	for b514ActivationVectorIndex := 0; b514ActivationVectorIndex < b514ActivationVectorPresent; b514ActivationVectorIndex++ {
		b514ActivationVectorView.b514ActivationVectorMask |= 1 << uint(b514ActivationVectorIndex)
	}
	b514ActivationVectorView.b514ActivationVectorComplete = len(b514ActivationVectorVector) == b514ActivationVectorView.b514ActivationVectorFields
	b514ActivationVectorView.b514ActivationVectorRemainder = b514ActivationVectorCopy(b514ActivationVectorVector[b514ActivationVectorPresent:])
	return b514ActivationVectorView
}

func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorState() string {
	return b514ActivationVectorView.b514ActivationVectorState
}

func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorGrade() string {
	return b514ActivationVectorView.b514ActivationVectorGrade
}

func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorRawBytes() []byte {
	return b514ActivationVectorCopy(b514ActivationVectorView.b514ActivationVectorRaw)
}

func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorPresentMask() uint16 {
	return b514ActivationVectorView.b514ActivationVectorMask
}

func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorIsComplete() bool {
	return b514ActivationVectorView.b514ActivationVectorComplete
}

func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorRemainderBytes() []byte {
	return b514ActivationVectorCopy(b514ActivationVectorView.b514ActivationVectorRemainder)
}

func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorFieldCount() int {
	return b514ActivationVectorView.b514ActivationVectorFields
}

func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorBitAt(b514ActivationVectorByte int, b514ActivationVectorBit int) B514ActivationVectorBit {
	if b514ActivationVectorView.b514ActivationVectorState != B514ActivationVectorStateBits || b514ActivationVectorByte < 0 || b514ActivationVectorByte >= b514ActivationVectorView.b514ActivationVectorFields || b514ActivationVectorBit < 0 || b514ActivationVectorBit > 7 {
		return B514ActivationVectorBit{}
	}
	if b514ActivationVectorView.b514ActivationVectorMask&(1<<uint(b514ActivationVectorByte)) == 0 {
		return B514ActivationVectorBit{}
	}
	return B514ActivationVectorBit{B514ActivationVectorPresent: true, B514ActivationVectorValue: b514ActivationVectorView.b514ActivationVectorRaw[b514ActivationVectorByte]&(1<<uint(b514ActivationVectorBit)) != 0}
}

func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorScalarAt(b514ActivationVectorIndex int) B514ActivationVectorScalar {
	b514ActivationVectorOut := B514ActivationVectorScalar{B514ActivationVectorInterpretation: B514ActivationVectorGradeUnknown}
	if b514ActivationVectorView.b514ActivationVectorState != B514ActivationVectorStateScalars || b514ActivationVectorIndex < 0 || b514ActivationVectorIndex >= b514ActivationVectorView.b514ActivationVectorFields {
		return b514ActivationVectorOut
	}
	if b514ActivationVectorView.b514ActivationVectorMask&(1<<uint(b514ActivationVectorIndex)) == 0 {
		return b514ActivationVectorOut
	}
	b514ActivationVectorOut.B514ActivationVectorPresent = true
	b514ActivationVectorOut.B514ActivationVectorRaw = b514ActivationVectorView.b514ActivationVectorRaw[b514ActivationVectorIndex]
	return b514ActivationVectorOut
}
