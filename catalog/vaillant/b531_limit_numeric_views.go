package vaillant

const (
	B531LimitNumericViewFamily    = "vaillant-b531-family"
	B531LimitNumericViewVariant   = "b531-op02-candidate"
	B531LimitNumericViewDirection = "positive"
)

type B531LimitNumericViewField struct {
	B531LimitNumericViewRaw       []byte `json:"Ab19An19Raw"`
	B531LimitNumericViewState     string `json:"Ab19An19State"`
	B531LimitNumericViewAvailable bool   `json:"Ab19An19Available"`
	B531LimitNumericViewValue     uint32 `json:"Ab19An19Value"`
	B531LimitNumericViewFresh     bool   `json:"Ab19An19Fresh"`
}

type B531LimitNumericViewPair struct {
	B531LimitNumericViewMinimum B531LimitNumericViewField `json:"Ab19An19Minimum"`
	B531LimitNumericViewMaximum B531LimitNumericViewField `json:"Ab19An19Maximum"`
	B531LimitNumericViewFresh   bool                      `json:"Ab19An19Fresh"`
}

type B531LimitNumericViewResult struct {
	B531LimitNumericViewAdmitted            bool                     `json:"Ab19An19Admitted"`
	B531LimitNumericViewAttempt             B531LimitNumericViewPair `json:"Ab19An19Attempt"`
	B531LimitNumericViewPrior               B531LimitNumericViewPair `json:"Ab19An19Prior"`
	B531LimitNumericViewLastKnown           B531LimitNumericViewPair `json:"Ab19An19LastKnown"`
	B531LimitNumericViewContractGrade       string                   `json:"Ab19An19ContractGrade"`
	B531LimitNumericViewInterpretationGrade string                   `json:"Ab19An19InterpretationGrade"`
	B531LimitNumericViewNativeUnit          string                   `json:"Ab19An19NativeUnit"`
	B531LimitNumericViewQualification       string                   `json:"Ab19An19Qualification"`
	B531LimitNumericViewPresence            string                   `json:"Ab19An19Presence"`
	B531LimitNumericViewExecution           string                   `json:"Ab19An19Execution"`
}

func b531LimitNumericViewClone(b531LimitNumericViewField B531LimitNumericViewField) B531LimitNumericViewField {
	b531LimitNumericViewField.B531LimitNumericViewRaw = append([]byte(nil), b531LimitNumericViewField.B531LimitNumericViewRaw...)
	return b531LimitNumericViewField
}

func b531LimitNumericViewDecode(b531LimitNumericViewRaw []byte, b531LimitNumericViewReadOK, b531LimitNumericViewAdmitted bool) B531LimitNumericViewField {
	b531LimitNumericViewField := B531LimitNumericViewField{B531LimitNumericViewRaw: append([]byte(nil), b531LimitNumericViewRaw...), B531LimitNumericViewState: "not-admitted"}
	if !b531LimitNumericViewAdmitted {
		return b531LimitNumericViewField
	}
	switch {
	case !b531LimitNumericViewReadOK:
		b531LimitNumericViewField.B531LimitNumericViewState = "read-failed"
	case len(b531LimitNumericViewRaw) == 0:
		b531LimitNumericViewField.B531LimitNumericViewState = "missing"
	case len(b531LimitNumericViewRaw) != 2:
		b531LimitNumericViewField.B531LimitNumericViewState = "malformed"
	default:
		b531LimitNumericViewField.B531LimitNumericViewState = "success"
		b531LimitNumericViewField.B531LimitNumericViewAvailable = true
		b531LimitNumericViewField.B531LimitNumericViewFresh = true
		b531LimitNumericViewField.B531LimitNumericViewValue = (uint32(b531LimitNumericViewRaw[0]) | uint32(b531LimitNumericViewRaw[1])<<8) * 10
	}
	return b531LimitNumericViewField
}

func b531LimitNumericViewRetain(b531LimitNumericViewAttempt, b531LimitNumericViewPrior B531LimitNumericViewField) B531LimitNumericViewField {
	if b531LimitNumericViewAttempt.B531LimitNumericViewAvailable {
		return b531LimitNumericViewClone(b531LimitNumericViewAttempt)
	}
	b531LimitNumericViewPrior = b531LimitNumericViewClone(b531LimitNumericViewPrior)
	b531LimitNumericViewPrior.B531LimitNumericViewFresh = false
	return b531LimitNumericViewPrior
}

func B531LimitNumericViewInterpret(b531LimitNumericViewFamily, b531LimitNumericViewVariant string, b531LimitNumericViewOperation byte, b531LimitNumericViewDirection string, b531LimitNumericViewMinimumRaw, b531LimitNumericViewMaximumRaw []byte, b531LimitNumericViewMinimumReadOK, b531LimitNumericViewMaximumReadOK bool, b531LimitNumericViewPrior B531LimitNumericViewPair) B531LimitNumericViewResult {
	b531LimitNumericViewResult := B531LimitNumericViewResult{
		B531LimitNumericViewAdmitted:      b531LimitNumericViewFamily == B531LimitNumericViewFamily && b531LimitNumericViewVariant == B531LimitNumericViewVariant && b531LimitNumericViewOperation == 0x02 && b531LimitNumericViewDirection == B531LimitNumericViewDirection,
		B531LimitNumericViewPrior:         B531LimitNumericViewPair{B531LimitNumericViewMinimum: b531LimitNumericViewClone(b531LimitNumericViewPrior.B531LimitNumericViewMinimum), B531LimitNumericViewMaximum: b531LimitNumericViewClone(b531LimitNumericViewPrior.B531LimitNumericViewMaximum), B531LimitNumericViewFresh: b531LimitNumericViewPrior.B531LimitNumericViewFresh},
		B531LimitNumericViewContractGrade: "Proven", B531LimitNumericViewInterpretationGrade: "Hypothesis", B531LimitNumericViewNativeUnit: "Unknown", B531LimitNumericViewQualification: "Unknown", B531LimitNumericViewPresence: "Unknown", B531LimitNumericViewExecution: "Unknown",
	}
	b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewMinimum = b531LimitNumericViewDecode(b531LimitNumericViewMinimumRaw, b531LimitNumericViewMinimumReadOK, b531LimitNumericViewResult.B531LimitNumericViewAdmitted)
	b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewMaximum = b531LimitNumericViewDecode(b531LimitNumericViewMaximumRaw, b531LimitNumericViewMaximumReadOK, b531LimitNumericViewResult.B531LimitNumericViewAdmitted)
	b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewFresh = b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewMinimum.B531LimitNumericViewFresh && b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewMaximum.B531LimitNumericViewFresh
	b531LimitNumericViewResult.B531LimitNumericViewLastKnown = B531LimitNumericViewPair{
		B531LimitNumericViewMinimum: b531LimitNumericViewRetain(b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewMinimum, b531LimitNumericViewPrior.B531LimitNumericViewMinimum),
		B531LimitNumericViewMaximum: b531LimitNumericViewRetain(b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewMaximum, b531LimitNumericViewPrior.B531LimitNumericViewMaximum),
		B531LimitNumericViewFresh:   b531LimitNumericViewResult.B531LimitNumericViewAttempt.B531LimitNumericViewFresh,
	}
	return b531LimitNumericViewResult
}

func B531LimitNumericViewGetter(b531LimitNumericViewRaw uint32) uint16 {
	return uint16(b531LimitNumericViewRaw / 10)
}
