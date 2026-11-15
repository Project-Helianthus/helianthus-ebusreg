package vaillant

// The selection values, field counts, grades and states of the B514 activation vector views.
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

// B514ActivationVectorIdentity selects a B514 activation vector view by family, variant, operation and
// direction.
type B514ActivationVectorIdentity struct {
	B514ActivationVectorFamily    string `json:"family"`
	B514ActivationVectorVariant   string `json:"variant"`
	B514ActivationVectorOperation uint8  `json:"operation"`
	B514ActivationVectorDirection string `json:"direction"`
}

// B514ActivationVectorBit is one activation bit; Value is meaningful only when Present is true.
type B514ActivationVectorBit struct {
	B514ActivationVectorPresent bool `json:"present"`
	B514ActivationVectorValue   bool `json:"value"`
}

// B514ActivationVectorScalar is one activation scalar; Raw is meaningful only when Present is true.
type B514ActivationVectorScalar struct {
	B514ActivationVectorPresent        bool   `json:"present"`
	B514ActivationVectorRaw            byte   `json:"raw"`
	B514ActivationVectorInterpretation string `json:"interpretation"`
}

// B514ActivationVectorView is the immutable interpretation of an activation vector; its accessors
// return copies.
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

// B514ActivationVectorInterpret interprets the vector for the identity: operation 3 yields bit fields
// and operation 4 scalar fields. Any other identity yields a raw-only view.
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

// B514ActivationVectorState returns the view state: raw-only, bits or scalars.
func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorState() string {
	return b514ActivationVectorView.b514ActivationVectorState
}

// B514ActivationVectorGrade returns the grade of the interpretation.
func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorGrade() string {
	return b514ActivationVectorView.b514ActivationVectorGrade
}

// B514ActivationVectorRawBytes returns a copy of the raw vector.
func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorRawBytes() []byte {
	return b514ActivationVectorCopy(b514ActivationVectorView.b514ActivationVectorRaw)
}

// B514ActivationVectorPresentMask returns the mask whose bit i is set when field i is present.
func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorPresentMask() uint16 {
	return b514ActivationVectorView.b514ActivationVectorMask
}

// B514ActivationVectorIsComplete reports whether the vector holds exactly the field count of the
// selected operation.
func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorIsComplete() bool {
	return b514ActivationVectorView.b514ActivationVectorComplete
}

// B514ActivationVectorRemainderBytes returns a copy of the bytes after the last field.
func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorRemainderBytes() []byte {
	return b514ActivationVectorCopy(b514ActivationVectorView.b514ActivationVectorRemainder)
}

// B514ActivationVectorFieldCount returns the number of fields of the selected operation, or zero for a
// raw-only view.
func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorFieldCount() int {
	return b514ActivationVectorView.b514ActivationVectorFields
}

// B514ActivationVectorBitAt returns the given bit of the given field byte; it is not present for a
// scalar or raw-only view, an index out of range or an absent field.
func (b514ActivationVectorView B514ActivationVectorView) B514ActivationVectorBitAt(b514ActivationVectorByte int, b514ActivationVectorBit int) B514ActivationVectorBit {
	if b514ActivationVectorView.b514ActivationVectorState != B514ActivationVectorStateBits || b514ActivationVectorByte < 0 || b514ActivationVectorByte >= b514ActivationVectorView.b514ActivationVectorFields || b514ActivationVectorBit < 0 || b514ActivationVectorBit > 7 {
		return B514ActivationVectorBit{}
	}
	if b514ActivationVectorView.b514ActivationVectorMask&(1<<uint(b514ActivationVectorByte)) == 0 {
		return B514ActivationVectorBit{}
	}
	return B514ActivationVectorBit{B514ActivationVectorPresent: true, B514ActivationVectorValue: b514ActivationVectorView.b514ActivationVectorRaw[b514ActivationVectorByte]&(1<<uint(b514ActivationVectorBit)) != 0}
}

// B514ActivationVectorScalarAt returns the scalar field at the given index; it is not present for a
// bit or raw-only view, an index out of range or an absent field.
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
