package system

import ab09an09binary "encoding/binary"

type Ab09An09Selection struct {
	Ab09An09Family    string
	Ab09An09Variant   string
	Ab09An09Operation byte
	Ab09An09Direction string
}

type Ab09An09Fields struct {
	Ab09An09Temperature   []byte
	Ab09An09Flow          []byte
	Ab09An09Volume        []byte
	Ab09An09MinimumVolume []byte
	Ab09An09MaximumVolume []byte
}

type Ab09An09Candidate struct {
	Ab09An09Raw                 []byte
	Ab09An09State               string
	Ab09An09Value               *float64
	Ab09An09Unit                string
	Ab09An09SemanticGrade       string
	Ab09An09NativeQualification string
}

type Ab09An09QuantityView struct {
	Ab09An09Identity      Ab09An09Selection
	Ab09An09Selection     Ab09An09Selection
	Ab09An09Selected      bool
	Ab09An09ContractGrade string
	Ab09An09Temperature   Ab09An09Candidate
	Ab09An09Flow          Ab09An09Candidate
	Ab09An09Volume        Ab09An09Candidate
	Ab09An09MinimumVolume Ab09An09Candidate
	Ab09An09MaximumVolume Ab09An09Candidate
}

func ab09an09Candidate(ab09an09Raw []byte, ab09an09Selected bool, ab09an09Unit string, ab09an09Divisor float64, ab09an09Signed bool) Ab09An09Candidate {
	ab09an09Result := Ab09An09Candidate{
		Ab09An09State:               "ab09an09_absent",
		Ab09An09Unit:                ab09an09Unit,
		Ab09An09SemanticGrade:       "Hypothesis",
		Ab09An09NativeQualification: "Unknown",
	}
	if ab09an09Unit == "" {
		ab09an09Result.Ab09An09SemanticGrade = "Unknown"
	}
	if ab09an09Raw == nil {
		return ab09an09Result
	}
	ab09an09Result.Ab09An09Raw = append([]byte{}, ab09an09Raw...)
	if !ab09an09Selected {
		ab09an09Result.Ab09An09State = "ab09an09_unselected"
		return ab09an09Result
	}
	if ab09an09Unit == "" {
		ab09an09Result.Ab09An09State = "ab09an09_opaque"
		ab09an09Result.Ab09An09SemanticGrade = "Unknown"
		return ab09an09Result
	}
	if len(ab09an09Raw) != 2 {
		ab09an09Result.Ab09An09State = "ab09an09_malformed"
		return ab09an09Result
	}
	ab09an09Word := ab09an09binary.LittleEndian.Uint16(ab09an09Raw)
	ab09an09Value := float64(ab09an09Word)
	if ab09an09Signed {
		ab09an09Value = float64(int16(ab09an09Word))
	}
	ab09an09Value /= ab09an09Divisor
	ab09an09Result.Ab09An09Value = &ab09an09Value
	ab09an09Result.Ab09An09State = "ab09an09_candidate"
	return ab09an09Result
}

func Ab09An09B509QuantityViews(ab09an09Identity, ab09an09Selection Ab09An09Selection, ab09an09Fields Ab09An09Fields) Ab09An09QuantityView {
	ab09an09Selected := ab09an09Identity == ab09an09Selection &&
		ab09an09Selection.Ab09An09Family != "" && ab09an09Selection.Ab09An09Variant != "" &&
		ab09an09Selection.Ab09An09Operation == 0x59 &&
		(ab09an09Selection.Ab09An09Direction == "ab09an09_request" || ab09an09Selection.Ab09An09Direction == "ab09an09_response")
	return Ab09An09QuantityView{
		Ab09An09Identity:      ab09an09Identity,
		Ab09An09Selection:     ab09an09Selection,
		Ab09An09Selected:      ab09an09Selected,
		Ab09An09ContractGrade: "Proven",
		Ab09An09Temperature:   ab09an09Candidate(ab09an09Fields.Ab09An09Temperature, ab09an09Selected, "degrees Celsius", 16, true),
		Ab09An09Flow:          ab09an09Candidate(ab09an09Fields.Ab09An09Flow, ab09an09Selected, "litres/minute", 100, false),
		Ab09An09Volume:        ab09an09Candidate(ab09an09Fields.Ab09An09Volume, ab09an09Selected, "litres", 10, false),
		Ab09An09MinimumVolume: ab09an09Candidate(ab09an09Fields.Ab09An09MinimumVolume, ab09an09Selected, "", 0, false),
		Ab09An09MaximumVolume: ab09an09Candidate(ab09an09Fields.Ab09An09MaximumVolume, ab09an09Selected, "", 0, false),
	}
}
