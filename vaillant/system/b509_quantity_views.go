package system

import b509QuantityViewBinary "encoding/binary"

type B509QuantityViewSelection struct {
	B509QuantityViewFamily    string `json:"Ab09An09Family"`
	B509QuantityViewVariant   string `json:"Ab09An09Variant"`
	B509QuantityViewOperation byte   `json:"Ab09An09Operation"`
	B509QuantityViewDirection string `json:"Ab09An09Direction"`
}

type B509QuantityViewFields struct {
	B509QuantityViewTemperature   []byte `json:"Ab09An09Temperature"`
	B509QuantityViewFlow          []byte `json:"Ab09An09Flow"`
	B509QuantityViewVolume        []byte `json:"Ab09An09Volume"`
	B509QuantityViewMinimumVolume []byte `json:"Ab09An09MinimumVolume"`
	B509QuantityViewMaximumVolume []byte `json:"Ab09An09MaximumVolume"`
}

type B509QuantityViewCandidate struct {
	B509QuantityViewRaw                 []byte   `json:"Ab09An09Raw"`
	B509QuantityViewState               string   `json:"Ab09An09State"`
	B509QuantityViewValue               *float64 `json:"Ab09An09Value"`
	B509QuantityViewUnit                string   `json:"Ab09An09Unit"`
	B509QuantityViewSemanticGrade       string   `json:"Ab09An09SemanticGrade"`
	B509QuantityViewNativeQualification string   `json:"Ab09An09NativeQualification"`
}

type B509QuantityView struct {
	B509QuantityViewIdentity      B509QuantityViewSelection `json:"Ab09An09Identity"`
	B509QuantityViewSelection     B509QuantityViewSelection `json:"Ab09An09Selection"`
	B509QuantityViewSelected      bool                      `json:"Ab09An09Selected"`
	B509QuantityViewContractGrade string                    `json:"Ab09An09ContractGrade"`
	B509QuantityViewTemperature   B509QuantityViewCandidate `json:"Ab09An09Temperature"`
	B509QuantityViewFlow          B509QuantityViewCandidate `json:"Ab09An09Flow"`
	B509QuantityViewVolume        B509QuantityViewCandidate `json:"Ab09An09Volume"`
	B509QuantityViewMinimumVolume B509QuantityViewCandidate `json:"Ab09An09MinimumVolume"`
	B509QuantityViewMaximumVolume B509QuantityViewCandidate `json:"Ab09An09MaximumVolume"`
}

func b509QuantityViewCandidate(b509QuantityViewRaw []byte, b509QuantityViewSelected bool, b509QuantityViewUnit string, b509QuantityViewDivisor float64, b509QuantityViewSigned bool) B509QuantityViewCandidate {
	b509QuantityViewResult := B509QuantityViewCandidate{
		B509QuantityViewState:               "ab09an09_absent",
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
		b509QuantityViewResult.B509QuantityViewState = "ab09an09_unselected"
		return b509QuantityViewResult
	}
	if b509QuantityViewUnit == "" {
		b509QuantityViewResult.B509QuantityViewState = "ab09an09_opaque"
		b509QuantityViewResult.B509QuantityViewSemanticGrade = "Unknown"
		return b509QuantityViewResult
	}
	if len(b509QuantityViewRaw) != 2 {
		b509QuantityViewResult.B509QuantityViewState = "ab09an09_malformed"
		return b509QuantityViewResult
	}
	b509QuantityViewWord := b509QuantityViewBinary.LittleEndian.Uint16(b509QuantityViewRaw)
	b509QuantityViewValue := float64(b509QuantityViewWord)
	if b509QuantityViewSigned {
		b509QuantityViewValue = float64(int16(b509QuantityViewWord))
	}
	b509QuantityViewValue /= b509QuantityViewDivisor
	b509QuantityViewResult.B509QuantityViewValue = &b509QuantityViewValue
	b509QuantityViewResult.B509QuantityViewState = "ab09an09_candidate"
	return b509QuantityViewResult
}

func B509QuantityViews(b509QuantityViewIdentity, b509QuantityViewSelection B509QuantityViewSelection, b509QuantityViewFields B509QuantityViewFields) B509QuantityView {
	b509QuantityViewSelected := b509QuantityViewIdentity == b509QuantityViewSelection &&
		b509QuantityViewSelection.B509QuantityViewFamily != "" && b509QuantityViewSelection.B509QuantityViewVariant != "" &&
		b509QuantityViewSelection.B509QuantityViewOperation == 0x59 &&
		(b509QuantityViewSelection.B509QuantityViewDirection == "ab09an09_request" || b509QuantityViewSelection.B509QuantityViewDirection == "ab09an09_response")
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
