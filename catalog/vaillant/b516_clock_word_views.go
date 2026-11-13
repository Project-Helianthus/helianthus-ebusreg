package vaillant

const (
	B516ClockWordViewFamily  = "vaillant-clock-family"
	B516ClockWordViewVariant = "b516-candidate"
)

type B516ClockWordViewSelector struct {
	B516ClockWordViewFamily    string `json:"Ab15An15Family"`
	B516ClockWordViewVariant   string `json:"Ab15An15Variant"`
	B516ClockWordViewOperation byte   `json:"Ab15An15Operation"`
	B516ClockWordViewDirection string `json:"Ab15An15Direction"`
}

type B516ClockWordViewCell struct {
	B516ClockWordViewRaw    []byte `json:"Ab15An15Raw"`
	B516ClockWordViewFailed bool   `json:"Ab15An15Failed"`
}

type B516ClockWordViewField struct {
	B516ClockWordViewState string  `json:"Ab15An15State"`
	B516ClockWordViewValue *uint16 `json:"Ab15An15Value"`
}

type B516ClockWordViewCandidate struct {
	B516ClockWordViewAdmitted           bool   `json:"Ab15An15Admitted"`
	B516ClockWordViewContractGrade      string `json:"Ab15An15ContractGrade"`
	B516ClockWordViewFieldGrade         string `json:"Ab15An15FieldGrade"`
	B516ClockWordViewNativeSupport      string `json:"Ab15An15NativeSupport"`
	B516ClockWordViewPhysicalAdaptation string `json:"Ab15An15PhysicalAdaptation"`
	B516ClockWordViewPresence           string `json:"Ab15An15Presence"`
	B516ClockWordViewFreshness          string `json:"Ab15An15Freshness"`
	B516ClockWordViewExecution          string `json:"Ab15An15Execution"`
}

type B516ClockWordViewClockInput struct {
	B516ClockWordViewClock           [5]B516ClockWordViewCell `json:"Ab15An15Clock"`
	B516ClockWordViewSetpoint        *byte                    `json:"Ab15An15Setpoint"`
	B516ClockWordViewMode            *byte                    `json:"Ab15An15Mode"`
	B516ClockWordViewSpecialFunction *byte                    `json:"Ab15An15SpecialFunction"`
	B516ClockWordViewMessages        [5]*byte                 `json:"Ab15An15Messages"`
}

type B516ClockWordViewClockResult struct {
	B516ClockWordViewCandidate       B516ClockWordViewCandidate  `json:"Ab15An15Candidate"`
	B516ClockWordViewRaw             B516ClockWordViewClockInput `json:"Ab15An15Raw"`
	B516ClockWordViewClock           *[5]B516ClockWordViewField  `json:"Ab15An15Clock"`
	B516ClockWordViewSetpointCelsius *float64                    `json:"Ab15An15SetpointCelsius"`
	B516ClockWordViewTablesGrade     string                      `json:"Ab15An15TablesGrade"`
}

type B516ClockWordViewWordResult struct {
	B516ClockWordViewCandidate      B516ClockWordViewCandidate `json:"Ab15An15Candidate"`
	B516ClockWordViewRaw            [4]B516ClockWordViewCell   `json:"Ab15An15Raw"`
	B516ClockWordViewWords          *[4]B516ClockWordViewField `json:"Ab15An15Words"`
	B516ClockWordViewPositionsGrade string                     `json:"Ab15An15PositionsGrade"`
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
