package vaillant

const (
	B531LimitNumericViewFamily    = "vaillant-b531-family"
	B531LimitNumericViewVariant   = "b531-op02-candidate"
	B531LimitNumericViewDirection = "positive"
)

type B531LimitNumericViewField struct {
	B531LimitNumericViewRaw       []byte `json:"raw"`
	B531LimitNumericViewState     string `json:"state"`
	B531LimitNumericViewAvailable bool   `json:"available"`
	B531LimitNumericViewValue     uint32 `json:"value"`
	B531LimitNumericViewFresh     bool   `json:"fresh"`
}

type B531LimitNumericViewPair struct {
	B531LimitNumericViewMinimum B531LimitNumericViewField `json:"minimum"`
	B531LimitNumericViewMaximum B531LimitNumericViewField `json:"maximum"`
	B531LimitNumericViewFresh   bool                      `json:"fresh"`
}

type B531LimitNumericViewResult struct {
	B531LimitNumericViewAdmitted            bool                     `json:"admitted"`
	B531LimitNumericViewAttempt             B531LimitNumericViewPair `json:"attempt"`
	B531LimitNumericViewPrior               B531LimitNumericViewPair `json:"prior"`
	B531LimitNumericViewLastKnown           B531LimitNumericViewPair `json:"last_known"`
	B531LimitNumericViewContractGrade       string                   `json:"contract_grade"`
	B531LimitNumericViewInterpretationGrade string                   `json:"interpretation_grade"`
	B531LimitNumericViewNativeUnit          string                   `json:"native_unit"`
	B531LimitNumericViewQualification       string                   `json:"qualification"`
	B531LimitNumericViewPresence            string                   `json:"presence"`
	B531LimitNumericViewExecution           string                   `json:"execution"`
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
