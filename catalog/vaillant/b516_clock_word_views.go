package vaillant

// Family and variant that a B516ClockWordViewSelector must name for a clock or word view to be admitted.
const (
	B516ClockWordViewFamily  = "vaillant-clock-family"
	B516ClockWordViewVariant = "b516-candidate"
)

// B516ClockWordViewSelector names the family, variant, operation and direction a view is asked for.
type B516ClockWordViewSelector struct {
	B516ClockWordViewFamily    string `json:"family"`
	B516ClockWordViewVariant   string `json:"variant"`
	B516ClockWordViewOperation byte   `json:"operation"`
	B516ClockWordViewDirection string `json:"direction"`
}

// B516ClockWordViewCell is one raw read cell; Failed marks a read that did not return data.
type B516ClockWordViewCell struct {
	B516ClockWordViewRaw    []byte `json:"raw"`
	B516ClockWordViewFailed bool   `json:"failed"`
}

// B516ClockWordViewField is one decoded cell: State is numeric, missing, failed, replacement or
// malformed-decimal, and Value is set only when State is numeric.
type B516ClockWordViewField struct {
	B516ClockWordViewState string  `json:"state"`
	B516ClockWordViewValue *uint16 `json:"value"`
}

// B516ClockWordViewCandidate reports whether the selector was admitted and the grades of the view;
// support, adaptation, presence, freshness and execution stay Unknown.
type B516ClockWordViewCandidate struct {
	B516ClockWordViewAdmitted           bool   `json:"admitted"`
	B516ClockWordViewContractGrade      string `json:"contract_grade"`
	B516ClockWordViewFieldGrade         string `json:"field_grade"`
	B516ClockWordViewNativeSupport      string `json:"native_support"`
	B516ClockWordViewPhysicalAdaptation string `json:"physical_adaptation"`
	B516ClockWordViewPresence           string `json:"presence"`
	B516ClockWordViewFreshness          string `json:"freshness"`
	B516ClockWordViewExecution          string `json:"execution"`
}

// B516ClockWordViewClockInput holds the raw clock cells and the optional setpoint, mode, special function
// and message bytes of the clock view.
type B516ClockWordViewClockInput struct {
	B516ClockWordViewClock           [5]B516ClockWordViewCell `json:"clock"`
	B516ClockWordViewSetpoint        *byte                    `json:"setpoint"`
	B516ClockWordViewMode            *byte                    `json:"mode"`
	B516ClockWordViewSpecialFunction *byte                    `json:"special_function"`
	B516ClockWordViewMessages        [5]*byte                 `json:"messages"`
}

// B516ClockWordViewClockResult is the clock view result; Clock and SetpointCelsius are nil unless the
// selector was admitted, and Raw is a detached copy of the input.
type B516ClockWordViewClockResult struct {
	B516ClockWordViewCandidate       B516ClockWordViewCandidate  `json:"candidate"`
	B516ClockWordViewRaw             B516ClockWordViewClockInput `json:"raw"`
	B516ClockWordViewClock           *[5]B516ClockWordViewField  `json:"clock"`
	B516ClockWordViewSetpointCelsius *float64                    `json:"setpoint_celsius"`
	B516ClockWordViewTablesGrade     string                      `json:"tables_grade"`
}

// B516ClockWordViewWordResult is the four-word view result; Words is nil unless the selector was admitted,
// and Raw is a detached copy of the input.
type B516ClockWordViewWordResult struct {
	B516ClockWordViewCandidate      B516ClockWordViewCandidate `json:"candidate"`
	B516ClockWordViewRaw            [4]B516ClockWordViewCell   `json:"raw"`
	B516ClockWordViewWords          *[4]B516ClockWordViewField `json:"words"`
	B516ClockWordViewPositionsGrade string                     `json:"positions_grade"`
}

func b516ClockWordViewCandidate(b516ClockWordViewSelector B516ClockWordViewSelector, b516ClockWordViewOperation byte) B516ClockWordViewCandidate {
	return B516ClockWordViewCandidate{
		B516ClockWordViewAdmitted:      b516ClockWordViewSelector.B516ClockWordViewFamily == B516ClockWordViewFamily && b516ClockWordViewSelector.B516ClockWordViewVariant == B516ClockWordViewVariant && b516ClockWordViewSelector.B516ClockWordViewOperation == b516ClockWordViewOperation && b516ClockWordViewSelector.B516ClockWordViewDirection == "positive",
		B516ClockWordViewContractGrade: "Proven", B516ClockWordViewFieldGrade: "Hypothesis", B516ClockWordViewNativeSupport: "Unknown", B516ClockWordViewPhysicalAdaptation: "Unknown", B516ClockWordViewPresence: "Unknown", B516ClockWordViewFreshness: "Unknown", B516ClockWordViewExecution: "Unknown",
	}
}

func b516ClockWordViewCopyCell(b516ClockWordViewCell B516ClockWordViewCell) B516ClockWordViewCell {
	if b516ClockWordViewCell.B516ClockWordViewRaw != nil {
		b516ClockWordViewRaw := make([]byte, len(b516ClockWordViewCell.B516ClockWordViewRaw))
		copy(b516ClockWordViewRaw, b516ClockWordViewCell.B516ClockWordViewRaw)
		b516ClockWordViewCell.B516ClockWordViewRaw = b516ClockWordViewRaw
	}
	return b516ClockWordViewCell
}

func b516ClockWordViewCopyByte(b516ClockWordViewByte *byte) *byte {
	if b516ClockWordViewByte == nil {
		return nil
	}
	b516ClockWordViewValue := *b516ClockWordViewByte
	return &b516ClockWordViewValue
}

func b516ClockWordViewField(b516ClockWordViewCell B516ClockWordViewCell, b516ClockWordViewSize int) B516ClockWordViewField {
	if b516ClockWordViewCell.B516ClockWordViewFailed {
		return B516ClockWordViewField{B516ClockWordViewState: "failed"}
	}
	if len(b516ClockWordViewCell.B516ClockWordViewRaw) == 0 {
		return B516ClockWordViewField{B516ClockWordViewState: "missing"}
	}
	if len(b516ClockWordViewCell.B516ClockWordViewRaw) != b516ClockWordViewSize {
		return B516ClockWordViewField{B516ClockWordViewState: "failed"}
	}
	b516ClockWordViewValue := uint16(b516ClockWordViewCell.B516ClockWordViewRaw[0])
	if b516ClockWordViewSize == 2 {
		b516ClockWordViewValue |= uint16(b516ClockWordViewCell.B516ClockWordViewRaw[1]) << 8
	} else {
		if b516ClockWordViewValue == 255 {
			return B516ClockWordViewField{B516ClockWordViewState: "replacement"}
		}
		if b516ClockWordViewValue>>4 > 9 || b516ClockWordViewValue&15 > 9 {
			return B516ClockWordViewField{B516ClockWordViewState: "malformed-decimal"}
		}
		b516ClockWordViewValue = (b516ClockWordViewValue>>4)*10 + (b516ClockWordViewValue & 15)
	}
	return B516ClockWordViewField{B516ClockWordViewState: "numeric", B516ClockWordViewValue: &b516ClockWordViewValue}
}

// B516ClockWordViewClockView decodes the positive operation 0x08 clock view: each clock cell is one
// packed-decimal byte (0xFF is the replacement value) and the setpoint is in half degrees Celsius. A
// selector that is not admitted returns only the detached raw input and the candidate grades.
func B516ClockWordViewClockView(b516ClockWordViewSelector B516ClockWordViewSelector, b516ClockWordViewInput B516ClockWordViewClockInput) B516ClockWordViewClockResult {
	b516ClockWordViewResult := B516ClockWordViewClockResult{B516ClockWordViewCandidate: b516ClockWordViewCandidate(b516ClockWordViewSelector, 0x08), B516ClockWordViewRaw: b516ClockWordViewInput, B516ClockWordViewTablesGrade: "Unknown"}
	for b516ClockWordViewIndex := range b516ClockWordViewInput.B516ClockWordViewClock {
		b516ClockWordViewResult.B516ClockWordViewRaw.B516ClockWordViewClock[b516ClockWordViewIndex] = b516ClockWordViewCopyCell(b516ClockWordViewInput.B516ClockWordViewClock[b516ClockWordViewIndex])
	}
	b516ClockWordViewResult.B516ClockWordViewRaw.B516ClockWordViewSetpoint = b516ClockWordViewCopyByte(b516ClockWordViewInput.B516ClockWordViewSetpoint)
	b516ClockWordViewResult.B516ClockWordViewRaw.B516ClockWordViewMode = b516ClockWordViewCopyByte(b516ClockWordViewInput.B516ClockWordViewMode)
	b516ClockWordViewResult.B516ClockWordViewRaw.B516ClockWordViewSpecialFunction = b516ClockWordViewCopyByte(b516ClockWordViewInput.B516ClockWordViewSpecialFunction)
	for b516ClockWordViewIndex := range b516ClockWordViewInput.B516ClockWordViewMessages {
		b516ClockWordViewResult.B516ClockWordViewRaw.B516ClockWordViewMessages[b516ClockWordViewIndex] = b516ClockWordViewCopyByte(b516ClockWordViewInput.B516ClockWordViewMessages[b516ClockWordViewIndex])
	}
	if !b516ClockWordViewResult.B516ClockWordViewCandidate.B516ClockWordViewAdmitted {
		return b516ClockWordViewResult
	}
	b516ClockWordViewResult.B516ClockWordViewClock = new([5]B516ClockWordViewField)
	for b516ClockWordViewIndex, b516ClockWordViewCell := range b516ClockWordViewInput.B516ClockWordViewClock {
		b516ClockWordViewResult.B516ClockWordViewClock[b516ClockWordViewIndex] = b516ClockWordViewField(b516ClockWordViewCell, 1)
	}
	if b516ClockWordViewInput.B516ClockWordViewSetpoint != nil {
		b516ClockWordViewValue := float64(*b516ClockWordViewInput.B516ClockWordViewSetpoint) / 2
		b516ClockWordViewResult.B516ClockWordViewSetpointCelsius = &b516ClockWordViewValue
	}
	return b516ClockWordViewResult
}

// B516ClockWordView decodes the positive operation 0x0D view as four little-endian 16-bit words; a cell of
// any other length is failed. A selector that is not admitted returns only the detached raw input.
func B516ClockWordView(b516ClockWordViewSelector B516ClockWordViewSelector, b516ClockWordViewWords [4]B516ClockWordViewCell) B516ClockWordViewWordResult {
	b516ClockWordViewResult := B516ClockWordViewWordResult{B516ClockWordViewCandidate: b516ClockWordViewCandidate(b516ClockWordViewSelector, 0x0d), B516ClockWordViewPositionsGrade: "Unknown"}
	for b516ClockWordViewIndex, b516ClockWordViewCell := range b516ClockWordViewWords {
		b516ClockWordViewResult.B516ClockWordViewRaw[b516ClockWordViewIndex] = b516ClockWordViewCopyCell(b516ClockWordViewCell)
	}
	if !b516ClockWordViewResult.B516ClockWordViewCandidate.B516ClockWordViewAdmitted {
		return b516ClockWordViewResult
	}
	b516ClockWordViewResult.B516ClockWordViewWords = new([4]B516ClockWordViewField)
	for b516ClockWordViewIndex, b516ClockWordViewCell := range b516ClockWordViewWords {
		b516ClockWordViewResult.B516ClockWordViewWords[b516ClockWordViewIndex] = b516ClockWordViewField(b516ClockWordViewCell, 2)
	}
	return b516ClockWordViewResult
}
