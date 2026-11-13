package system

import b509QuantityViewBinary "encoding/binary"

type B509QuantityViewSelection struct {
	B509QuantityViewFamily    string `json:"family"`
	B509QuantityViewVariant   string `json:"variant"`
	B509QuantityViewOperation byte   `json:"operation"`
	B509QuantityViewDirection string `json:"direction"`
}

type B509QuantityViewFields struct {
	B509QuantityViewTemperature   []byte `json:"temperature"`
	B509QuantityViewFlow          []byte `json:"flow"`
	B509QuantityViewVolume        []byte `json:"volume"`
	B509QuantityViewMinimumVolume []byte `json:"minimum_volume"`
	B509QuantityViewMaximumVolume []byte `json:"maximum_volume"`
}

type B509QuantityViewCandidate struct {
	B509QuantityViewRaw                 []byte   `json:"raw"`
	B509QuantityViewState               string   `json:"state"`
	B509QuantityViewValue               *float64 `json:"value"`
	B509QuantityViewUnit                string   `json:"unit"`
	B509QuantityViewSemanticGrade       string   `json:"semantic_grade"`
	B509QuantityViewNativeQualification string   `json:"native_qualification"`
}

type B509QuantityView struct {
	B509QuantityViewIdentity      B509QuantityViewSelection `json:"identity"`
	B509QuantityViewSelection     B509QuantityViewSelection `json:"selection"`
	B509QuantityViewSelected      bool                      `json:"selected"`
	B509QuantityViewContractGrade string                    `json:"contract_grade"`
	B509QuantityViewTemperature   B509QuantityViewCandidate `json:"temperature"`
	B509QuantityViewFlow          B509QuantityViewCandidate `json:"flow"`
	B509QuantityViewVolume        B509QuantityViewCandidate `json:"volume"`
	B509QuantityViewMinimumVolume B509QuantityViewCandidate `json:"minimum_volume"`
	B509QuantityViewMaximumVolume B509QuantityViewCandidate `json:"maximum_volume"`
}

func b509QuantityViewCandidate(b509QuantityViewRaw []byte, b509QuantityViewSelected bool, b509QuantityViewUnit string, b509QuantityViewDivisor float64, b509QuantityViewSigned bool) B509QuantityViewCandidate {
	b509QuantityViewResult := B509QuantityViewCandidate{
		B509QuantityViewState:               "absent",
		B509QuantityViewUnit:                b509QuantityViewUnit,
		B509QuantityViewSemanticGrade:       "Hypothesis",
		B509QuantityViewNativeQualification: "Unknown",
	}
	if b509QuantityViewUnit == "" {
		b509QuantityViewResult.B509QuantityViewSemanticGrade = "Unknown"
	}
	if b509QuantityViewRaw == nil {
		return b509QuantityViewResult
	}
	b509QuantityViewResult.B509QuantityViewRaw = append([]byte{}, b509QuantityViewRaw...)
	if !b509QuantityViewSelected {
		b509QuantityViewResult.B509QuantityViewState = "unselected"
		return b509QuantityViewResult
	}
	if b509QuantityViewUnit == "" {
		b509QuantityViewResult.B509QuantityViewState = "opaque"
		b509QuantityViewResult.B509QuantityViewSemanticGrade = "Unknown"
		return b509QuantityViewResult
	}
	if len(b509QuantityViewRaw) != 2 {
		b509QuantityViewResult.B509QuantityViewState = "malformed"
		return b509QuantityViewResult
	}
	b509QuantityViewWord := b509QuantityViewBinary.LittleEndian.Uint16(b509QuantityViewRaw)
	b509QuantityViewValue := float64(b509QuantityViewWord)
	if b509QuantityViewSigned {
		b509QuantityViewValue = float64(int16(b509QuantityViewWord))
	}
	b509QuantityViewValue /= b509QuantityViewDivisor
	b509QuantityViewResult.B509QuantityViewValue = &b509QuantityViewValue
	b509QuantityViewResult.B509QuantityViewState = "candidate"
	return b509QuantityViewResult
}

func B509QuantityViews(b509QuantityViewIdentity, b509QuantityViewSelection B509QuantityViewSelection, b509QuantityViewFields B509QuantityViewFields) B509QuantityView {
	b509QuantityViewSelected := b509QuantityViewIdentity == b509QuantityViewSelection &&
		b509QuantityViewSelection.B509QuantityViewFamily != "" && b509QuantityViewSelection.B509QuantityViewVariant != "" &&
		b509QuantityViewSelection.B509QuantityViewOperation == 0x59 &&
		(b509QuantityViewSelection.B509QuantityViewDirection == "request" || b509QuantityViewSelection.B509QuantityViewDirection == "response")
	return B509QuantityView{
		B509QuantityViewIdentity:      b509QuantityViewIdentity,
		B509QuantityViewSelection:     b509QuantityViewSelection,
		B509QuantityViewSelected:      b509QuantityViewSelected,
		B509QuantityViewContractGrade: "Proven",
		B509QuantityViewTemperature:   b509QuantityViewCandidate(b509QuantityViewFields.B509QuantityViewTemperature, b509QuantityViewSelected, "degrees Celsius", 16, true),
		B509QuantityViewFlow:          b509QuantityViewCandidate(b509QuantityViewFields.B509QuantityViewFlow, b509QuantityViewSelected, "litres per minute", 100, false),
		B509QuantityViewVolume:        b509QuantityViewCandidate(b509QuantityViewFields.B509QuantityViewVolume, b509QuantityViewSelected, "litres", 10, false),
		B509QuantityViewMinimumVolume: b509QuantityViewCandidate(b509QuantityViewFields.B509QuantityViewMinimumVolume, b509QuantityViewSelected, "", 0, false),
		B509QuantityViewMaximumVolume: b509QuantityViewCandidate(b509QuantityViewFields.B509QuantityViewMaximumVolume, b509QuantityViewSelected, "", 0, false),
	}
}
