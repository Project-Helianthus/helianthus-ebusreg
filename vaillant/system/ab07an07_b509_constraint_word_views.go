package system

import ab07an07bytes "bytes"

const Ab07An07PositiveResponse = "positive-response"

type Ab07An07Input struct {
	Ab07An07Family           string
	Ab07An07Variant          string
	Ab07An07Operation        byte
	Ab07An07Direction        string
	Ab07An07RequestedAddress []byte
	Ab07An07ReturnedAddress  []byte
	Ab07An07Minimum          []byte
	Ab07An07Maximum          []byte
	Ab07An07Default          []byte
	Ab07An07Step             []byte
	Ab07An07Factor           []byte
}

type Ab07An07Word struct {
	Ab07An07Raw   []byte
	Ab07An07State string
	Ab07An07Value int
}

type Ab07An07View struct {
	Ab07An07Context                     Ab07An07Input
	Ab07An07Selected                    bool
	Ab07An07RawAddressEqual             bool
	Ab07An07RequestedEditBoundsAdmitted bool
	Ab07An07Minimum                     Ab07An07Word
	Ab07An07Maximum                     Ab07An07Word
	Ab07An07Default                     Ab07An07Word
	Ab07An07Step                        Ab07An07Word
	Ab07An07Factor                      Ab07An07Word
	Ab07An07SoftwareContract            string
	Ab07An07Interpretation              string
	Ab07An07NativeEchoRules             string
	Ab07An07PhysicalQualification       string
	Ab07An07PacketAssociation           string
	Ab07An07FamilyApplicability         string
	Ab07An07FactorAlgorithm             string
	Ab07An07Writability                 string
	Ab07An07PeerEnforcement             string
	Ab07An07Freshness                   string
	Ab07An07Execution                   string
}

func ab07an07Clone(ab07an07Raw []byte) []byte {
	if ab07an07Raw == nil {
		return nil
	}
	return append([]byte{}, ab07an07Raw...)
}

func ab07an07Word(ab07an07Raw []byte, ab07an07Selected bool) Ab07An07Word {
	ab07an07Result := Ab07An07Word{Ab07An07Raw: ab07an07Clone(ab07an07Raw), Ab07An07State: "missing"}
	if ab07an07Raw == nil {
		return ab07an07Result
	}
	ab07an07Result.Ab07An07State = "malformed"
	if len(ab07an07Raw) != 2 {
		return ab07an07Result
	}
	ab07an07Result.Ab07An07State = "unselected"
	if !ab07an07Selected {
		return ab07an07Result
	}
	ab07an07Result.Ab07An07State = "decoded"
	ab07an07Result.Ab07An07Value = int(int16(uint16(ab07an07Raw[0]) | uint16(ab07an07Raw[1])<<8))
	return ab07an07Result
}

func Ab07An07Interpret(ab07an07Input Ab07An07Input) Ab07An07View {
	ab07an07Selected := ab07an07Input.Ab07An07Operation == 0x33 && ab07an07Input.Ab07An07Direction == Ab07An07PositiveResponse
	ab07an07Input.Ab07An07RequestedAddress = ab07an07Clone(ab07an07Input.Ab07An07RequestedAddress)
	ab07an07Input.Ab07An07ReturnedAddress = ab07an07Clone(ab07an07Input.Ab07An07ReturnedAddress)
	ab07an07Input.Ab07An07Minimum = ab07an07Clone(ab07an07Input.Ab07An07Minimum)
	ab07an07Input.Ab07An07Maximum = ab07an07Clone(ab07an07Input.Ab07An07Maximum)
	ab07an07Input.Ab07An07Default = ab07an07Clone(ab07an07Input.Ab07An07Default)
	ab07an07Input.Ab07An07Step = ab07an07Clone(ab07an07Input.Ab07An07Step)
	ab07an07Input.Ab07An07Factor = ab07an07Clone(ab07an07Input.Ab07An07Factor)
	return Ab07An07View{
		Ab07An07Context:          ab07an07Input,
		Ab07An07Selected:         ab07an07Selected,
		Ab07An07RawAddressEqual:  ab07an07bytes.Equal(ab07an07Input.Ab07An07RequestedAddress, ab07an07Input.Ab07An07ReturnedAddress),
		Ab07An07Minimum:          ab07an07Word(ab07an07Input.Ab07An07Minimum, ab07an07Selected),
		Ab07An07Maximum:          ab07an07Word(ab07an07Input.Ab07An07Maximum, ab07an07Selected),
		Ab07An07Default:          ab07an07Word(ab07an07Input.Ab07An07Default, ab07an07Selected),
		Ab07An07Step:             ab07an07Word(ab07an07Input.Ab07An07Step, ab07an07Selected),
		Ab07An07Factor:           ab07an07Word(ab07an07Input.Ab07An07Factor, ab07an07Selected),
		Ab07An07SoftwareContract: "Proven", Ab07An07Interpretation: "Hypothesis",
		Ab07An07NativeEchoRules: "Unknown", Ab07An07PhysicalQualification: "Unknown",
		Ab07An07PacketAssociation: "Unknown", Ab07An07FamilyApplicability: "Unknown",
		Ab07An07FactorAlgorithm: "Unknown", Ab07An07Writability: "Unknown",
		Ab07An07PeerEnforcement: "Unknown", Ab07An07Freshness: "Unknown", Ab07An07Execution: "Unknown",
	}
}
