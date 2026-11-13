package system

import b509ConstraintWordViewBytes "bytes"

const B509ConstraintWordViewPositiveResponse = "positive-response"

type B509ConstraintWordViewInput struct {
	B509ConstraintWordViewFamily           string `json:"Ab07An07Family"`
	B509ConstraintWordViewVariant          string `json:"Ab07An07Variant"`
	B509ConstraintWordViewOperation        byte   `json:"Ab07An07Operation"`
	B509ConstraintWordViewDirection        string `json:"Ab07An07Direction"`
	B509ConstraintWordViewRequestedAddress []byte `json:"Ab07An07RequestedAddress"`
	B509ConstraintWordViewReturnedAddress  []byte `json:"Ab07An07ReturnedAddress"`
	B509ConstraintWordViewMinimum          []byte `json:"Ab07An07Minimum"`
	B509ConstraintWordViewMaximum          []byte `json:"Ab07An07Maximum"`
	B509ConstraintWordViewDefault          []byte `json:"Ab07An07Default"`
	B509ConstraintWordViewStep             []byte `json:"Ab07An07Step"`
	B509ConstraintWordViewFactor           []byte `json:"Ab07An07Factor"`
}

type B509ConstraintWord struct {
	B509ConstraintWordViewRaw   []byte `json:"Ab07An07Raw"`
	B509ConstraintWordViewState string `json:"Ab07An07State"`
	B509ConstraintWordViewValue int    `json:"Ab07An07Value"`
}

type B509ConstraintWordView struct {
	B509ConstraintWordViewContext                     B509ConstraintWordViewInput `json:"Ab07An07Context"`
	B509ConstraintWordViewSelected                    bool                        `json:"Ab07An07Selected"`
	B509ConstraintWordViewRawAddressEqual             bool                        `json:"Ab07An07RawAddressEqual"`
	B509ConstraintWordViewRequestedEditBoundsAdmitted bool                        `json:"Ab07An07RequestedEditBoundsAdmitted"`
	B509ConstraintWordViewMinimum                     B509ConstraintWord          `json:"Ab07An07Minimum"`
	B509ConstraintWordViewMaximum                     B509ConstraintWord          `json:"Ab07An07Maximum"`
	B509ConstraintWordViewDefault                     B509ConstraintWord          `json:"Ab07An07Default"`
	B509ConstraintWordViewStep                        B509ConstraintWord          `json:"Ab07An07Step"`
	B509ConstraintWordViewFactor                      B509ConstraintWord          `json:"Ab07An07Factor"`
	B509ConstraintWordViewSoftwareContract            string                      `json:"Ab07An07SoftwareContract"`
	B509ConstraintWordViewInterpretation              string                      `json:"Ab07An07Interpretation"`
	B509ConstraintWordViewNativeEchoRules             string                      `json:"Ab07An07NativeEchoRules"`
	B509ConstraintWordViewPhysicalQualification       string                      `json:"Ab07An07PhysicalQualification"`
	B509ConstraintWordViewPacketAssociation           string                      `json:"Ab07An07PacketAssociation"`
	B509ConstraintWordViewFamilyApplicability         string                      `json:"Ab07An07FamilyApplicability"`
	B509ConstraintWordViewFactorAlgorithm             string                      `json:"Ab07An07FactorAlgorithm"`
	B509ConstraintWordViewWritability                 string                      `json:"Ab07An07Writability"`
	B509ConstraintWordViewPeerEnforcement             string                      `json:"Ab07An07PeerEnforcement"`
	B509ConstraintWordViewFreshness                   string                      `json:"Ab07An07Freshness"`
	B509ConstraintWordViewExecution                   string                      `json:"Ab07An07Execution"`
}

func b509ConstraintWordViewClone(b509ConstraintWordViewRaw []byte) []byte {
	if b509ConstraintWordViewRaw == nil {
		return nil
	}
	return append([]byte{}, b509ConstraintWordViewRaw...)
}

func b509ConstraintWordViewWord(b509ConstraintWordViewRaw []byte, b509ConstraintWordViewSelected bool) B509ConstraintWord {
	b509ConstraintWordViewResult := B509ConstraintWord{B509ConstraintWordViewRaw: b509ConstraintWordViewClone(b509ConstraintWordViewRaw), B509ConstraintWordViewState: "missing"}
	if b509ConstraintWordViewRaw == nil {
		return b509ConstraintWordViewResult
	}
	b509ConstraintWordViewResult.B509ConstraintWordViewState = "malformed"
	if len(b509ConstraintWordViewRaw) != 2 {
		return b509ConstraintWordViewResult
	}
	b509ConstraintWordViewResult.B509ConstraintWordViewState = "unselected"
	if !b509ConstraintWordViewSelected {
		return b509ConstraintWordViewResult
	}
	b509ConstraintWordViewResult.B509ConstraintWordViewState = "decoded"
	b509ConstraintWordViewResult.B509ConstraintWordViewValue = int(int16(uint16(b509ConstraintWordViewRaw[0]) | uint16(b509ConstraintWordViewRaw[1])<<8))
	return b509ConstraintWordViewResult
}

func B509ConstraintWordViewInterpret(b509ConstraintWordViewInput B509ConstraintWordViewInput) B509ConstraintWordView {
	b509ConstraintWordViewSelected := b509ConstraintWordViewInput.B509ConstraintWordViewOperation == 0x33 && b509ConstraintWordViewInput.B509ConstraintWordViewDirection == B509ConstraintWordViewPositiveResponse
	b509ConstraintWordViewInput.B509ConstraintWordViewRequestedAddress = b509ConstraintWordViewClone(b509ConstraintWordViewInput.B509ConstraintWordViewRequestedAddress)
	b509ConstraintWordViewInput.B509ConstraintWordViewReturnedAddress = b509ConstraintWordViewClone(b509ConstraintWordViewInput.B509ConstraintWordViewReturnedAddress)
	b509ConstraintWordViewInput.B509ConstraintWordViewMinimum = b509ConstraintWordViewClone(b509ConstraintWordViewInput.B509ConstraintWordViewMinimum)
	b509ConstraintWordViewInput.B509ConstraintWordViewMaximum = b509ConstraintWordViewClone(b509ConstraintWordViewInput.B509ConstraintWordViewMaximum)
	b509ConstraintWordViewInput.B509ConstraintWordViewDefault = b509ConstraintWordViewClone(b509ConstraintWordViewInput.B509ConstraintWordViewDefault)
	b509ConstraintWordViewInput.B509ConstraintWordViewStep = b509ConstraintWordViewClone(b509ConstraintWordViewInput.B509ConstraintWordViewStep)
	b509ConstraintWordViewInput.B509ConstraintWordViewFactor = b509ConstraintWordViewClone(b509ConstraintWordViewInput.B509ConstraintWordViewFactor)
	return B509ConstraintWordView{
		B509ConstraintWordViewContext:          b509ConstraintWordViewInput,
		B509ConstraintWordViewSelected:         b509ConstraintWordViewSelected,
		B509ConstraintWordViewRawAddressEqual:  b509ConstraintWordViewBytes.Equal(b509ConstraintWordViewInput.B509ConstraintWordViewRequestedAddress, b509ConstraintWordViewInput.B509ConstraintWordViewReturnedAddress),
		B509ConstraintWordViewMinimum:          b509ConstraintWordViewWord(b509ConstraintWordViewInput.B509ConstraintWordViewMinimum, b509ConstraintWordViewSelected),
		B509ConstraintWordViewMaximum:          b509ConstraintWordViewWord(b509ConstraintWordViewInput.B509ConstraintWordViewMaximum, b509ConstraintWordViewSelected),
		B509ConstraintWordViewDefault:          b509ConstraintWordViewWord(b509ConstraintWordViewInput.B509ConstraintWordViewDefault, b509ConstraintWordViewSelected),
		B509ConstraintWordViewStep:             b509ConstraintWordViewWord(b509ConstraintWordViewInput.B509ConstraintWordViewStep, b509ConstraintWordViewSelected),
		B509ConstraintWordViewFactor:           b509ConstraintWordViewWord(b509ConstraintWordViewInput.B509ConstraintWordViewFactor, b509ConstraintWordViewSelected),
		B509ConstraintWordViewSoftwareContract: "Proven", B509ConstraintWordViewInterpretation: "Hypothesis",
		B509ConstraintWordViewNativeEchoRules: "Unknown", B509ConstraintWordViewPhysicalQualification: "Unknown",
		B509ConstraintWordViewPacketAssociation: "Unknown", B509ConstraintWordViewFamilyApplicability: "Unknown",
		B509ConstraintWordViewFactorAlgorithm: "Unknown", B509ConstraintWordViewWritability: "Unknown",
		B509ConstraintWordViewPeerEnforcement: "Unknown", B509ConstraintWordViewFreshness: "Unknown", B509ConstraintWordViewExecution: "Unknown",
	}
}
