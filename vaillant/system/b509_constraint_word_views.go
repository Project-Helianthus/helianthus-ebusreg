package system

import b509ConstraintWordViewBytes "bytes"

// B509ConstraintWordViewPositiveResponse is the direction value that, together
// with operation 0x33, selects a constraint-word reply for decoding.
const B509ConstraintWordViewPositiveResponse = "positive-response"

// B509ConstraintWordViewInput is one observed B509 constraint exchange: its
// family and variant, operation and direction, the requested and returned
// register addresses and the raw minimum, maximum, default, step and factor
// words.
type B509ConstraintWordViewInput struct {
	B509ConstraintWordViewFamily           string `json:"family"`
	B509ConstraintWordViewVariant          string `json:"variant"`
	B509ConstraintWordViewOperation        byte   `json:"operation"`
	B509ConstraintWordViewDirection        string `json:"direction"`
	B509ConstraintWordViewRequestedAddress []byte `json:"requested_address"`
	B509ConstraintWordViewReturnedAddress  []byte `json:"returned_address"`
	B509ConstraintWordViewMinimum          []byte `json:"minimum"`
	B509ConstraintWordViewMaximum          []byte `json:"maximum"`
	B509ConstraintWordViewDefault          []byte `json:"default"`
	B509ConstraintWordViewStep             []byte `json:"step"`
	B509ConstraintWordViewFactor           []byte `json:"factor"`
}

// B509ConstraintWord is one raw constraint word and its decoding state:
// missing, malformed (not two bytes), unselected, or decoded as a signed
// little-endian value.
type B509ConstraintWord struct {
	B509ConstraintWordViewRaw   []byte `json:"raw"`
	B509ConstraintWordViewState string `json:"state"`
	B509ConstraintWordViewValue int    `json:"value"`
}

// B509ConstraintWordView is the graded interpretation of a constraint
// exchange. Only the software contract is proven and the word
// interpretation is a hypothesis; echo rules, physical qualification,
// writability, freshness and the other aspects stay unknown. A matching
// returned address is reported, not required.
type B509ConstraintWordView struct {
	B509ConstraintWordViewContext                     B509ConstraintWordViewInput `json:"context"`
	B509ConstraintWordViewSelected                    bool                        `json:"selected"`
	B509ConstraintWordViewRawAddressEqual             bool                        `json:"raw_address_equal"`
	B509ConstraintWordViewRequestedEditBoundsAdmitted bool                        `json:"requested_edit_bounds_admitted"`
	B509ConstraintWordViewMinimum                     B509ConstraintWord          `json:"minimum"`
	B509ConstraintWordViewMaximum                     B509ConstraintWord          `json:"maximum"`
	B509ConstraintWordViewDefault                     B509ConstraintWord          `json:"default"`
	B509ConstraintWordViewStep                        B509ConstraintWord          `json:"step"`
	B509ConstraintWordViewFactor                      B509ConstraintWord          `json:"factor"`
	B509ConstraintWordViewSoftwareContract            string                      `json:"software_contract"`
	B509ConstraintWordViewInterpretation              string                      `json:"interpretation"`
	B509ConstraintWordViewNativeEchoRules             string                      `json:"native_echo_rules"`
	B509ConstraintWordViewPhysicalQualification       string                      `json:"physical_qualification"`
	B509ConstraintWordViewPacketAssociation           string                      `json:"packet_association"`
	B509ConstraintWordViewFamilyApplicability         string                      `json:"family_applicability"`
	B509ConstraintWordViewFactorAlgorithm             string                      `json:"factor_algorithm"`
	B509ConstraintWordViewWritability                 string                      `json:"writability"`
	B509ConstraintWordViewPeerEnforcement             string                      `json:"peer_enforcement"`
	B509ConstraintWordViewFreshness                   string                      `json:"freshness"`
	B509ConstraintWordViewExecution                   string                      `json:"execution"`
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

// B509ConstraintWordViewInterpret decodes the constraint words of a positive
// reply to operation 0x33 and leaves any other exchange unselected. It never
// fails and copies every input slice, so the view shares no memory with the
// caller.
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
