package system

type Ab01An01Cell struct {
	Ab01An01Raw   byte
	Ab01An01State string
}

type Ab01An01Input struct {
	Ab01An01Family        string
	Ab01An01Software      string
	Ab01An01Hardware      string
	Ab01An01APIVersion    string
	Ab01An01APIRevision   string
	Ab01An01OP            byte
	Ab01An01RequestedGG   byte
	Ab01An01II            *byte
	Ab01An01RequestedRR   uint16
	Ab01An01Alternative   string
	Ab01An01ReturnedGG    byte
	Ab01An01ReturnedRR    uint16
	Ab01An01ResponseState string
	Ab01An01Cells         [3]Ab01An01Cell
	Ab01An01Raw           []byte
}

type Ab01An01Component struct {
	Ab01An01Cell      Ab01An01Cell
	Ab01An01State     string
	Ab01An01Candidate *byte
}

type Ab01An01Result struct {
	Ab01An01Context             Ab01An01Input
	Ab01An01Components          [3]Ab01An01Component
	Ab01An01CompleteCandidate   bool
	Ab01An01Grade               string
	Ab01An01ContractGrade       string
	Ab01An01NativeQualification string
	Ab01An01CombinedTimestamp   string
	Ab01An01EditValidation      string
}

func Ab01An01Interpret(ab01an01Input Ab01An01Input) Ab01An01Result {
	ab01an01Input.Ab01An01Raw = append([]byte(nil), ab01an01Input.Ab01An01Raw...)
	if ab01an01Input.Ab01An01II != nil {
		ab01an01II := *ab01an01Input.Ab01An01II
		ab01an01Input.Ab01An01II = &ab01an01II
	}
	ab01an01Result := Ab01An01Result{Ab01An01Context: ab01an01Input, Ab01An01Grade: "Unknown", Ab01An01ContractGrade: "Proven", Ab01An01NativeQualification: "Unknown", Ab01An01CombinedTimestamp: "Unknown", Ab01An01EditValidation: "Unknown"}
	for ab01an01Index, ab01an01Cell := range ab01an01Input.Ab01An01Cells {
		ab01an01Result.Ab01An01Components[ab01an01Index] = Ab01An01Component{Ab01An01Cell: ab01an01Cell, Ab01An01State: "Unknown"}
	}
	if !ab01an01Matches(ab01an01Input) {
		return ab01an01Result
	}
	ab01an01Result.Ab01An01Grade = "Hypothesis"
	ab01an01Result.Ab01An01CompleteCandidate = true
	for ab01an01Index, ab01an01Cell := range ab01an01Input.Ab01An01Cells {
		ab01an01Component := &ab01an01Result.Ab01An01Components[ab01an01Index]
		if ab01an01Cell.Ab01An01State != "present" {
			if ab01an01Cell.Ab01An01State == "missing" || ab01an01Cell.Ab01An01State == "failed" {
				ab01an01Component.Ab01An01State = ab01an01Cell.Ab01An01State
			}
			ab01an01Result.Ab01An01CompleteCandidate = false
			continue
		}
		if ab01an01Input.Ab01An01Alternative == "date" && ((ab01an01Index == 0 && (ab01an01Cell.Ab01An01Raw < 1 || ab01an01Cell.Ab01An01Raw > 31)) || (ab01an01Index == 1 && (ab01an01Cell.Ab01An01Raw < 1 || ab01an01Cell.Ab01An01Raw > 12)) || (ab01an01Index == 2 && ab01an01Cell.Ab01An01Raw > 99)) {
			ab01an01Component.Ab01An01State = "out-of-domain"
			ab01an01Result.Ab01An01CompleteCandidate = false
			continue
		}
		ab01an01Value := ab01an01Cell.Ab01An01Raw
		ab01an01Component.Ab01An01Candidate = &ab01an01Value
		ab01an01Component.Ab01An01State = "candidate"
	}
	return ab01an01Result
}

func ab01an01Matches(ab01an01Input Ab01An01Input) bool {
	if (ab01an01Input.Ab01An01Family != "controller-700" && ab01an01Input.Ab01An01Family != "controller-720") || ab01an01Input.Ab01An01Software == "" || ab01an01Input.Ab01An01Hardware == "" || ab01an01Input.Ab01An01APIVersion == "" || ab01an01Input.Ab01An01APIRevision == "" || ab01an01Input.Ab01An01OP != 2 || ab01an01Input.Ab01An01II == nil || ab01an01Input.Ab01An01ResponseState != "echo-payload" || ab01an01Input.Ab01An01RequestedGG != ab01an01Input.Ab01An01ReturnedGG || ab01an01Input.Ab01An01RequestedRR != ab01an01Input.Ab01An01ReturnedRR {
		return false
	}
	ab01an01Alternative := ""
	switch ab01an01Input.Ab01An01RequestedGG {
	case 0:
		if *ab01an01Input.Ab01An01II != 0 {
			return false
		}
		switch ab01an01Input.Ab01An01RequestedRR {
		case 0x34:
			ab01an01Alternative = "date"
		case 0x35:
			ab01an01Alternative = "time"
		case 0x71, 0x72:
			if ab01an01Input.Ab01An01Family == "controller-700" {
				ab01an01Alternative = "date"
			}
		case 0xda, 0xdb:
			if ab01an01Input.Ab01An01Family == "controller-720" {
				ab01an01Alternative = "date"
			}
		}
	case 1:
		if *ab01an01Input.Ab01An01II != 0 {
			return false
		}
		switch ab01an01Input.Ab01An01RequestedRR {
		case 9, 10:
			ab01an01Alternative = "date"
		case 0x10, 0x11:
			if ab01an01Input.Ab01An01Family == "controller-720" {
				ab01an01Alternative = "time"
			}
		}
	case 3:
		switch ab01an01Input.Ab01An01RequestedRR {
		case 3, 4:
			ab01an01Alternative = "date"
		case 0x24:
			if ab01an01Input.Ab01An01Family == "controller-720" {
				ab01an01Alternative = "date"
			}
		case 0x1e, 0x20, 0x21:
			if ab01an01Input.Ab01An01Family == "controller-720" {
				ab01an01Alternative = "time"
			}
		}
	case 9:
		if *ab01an01Input.Ab01An01II != 0 || ab01an01Input.Ab01An01Family != "controller-720" {
			return false
		}
		switch ab01an01Input.Ab01An01RequestedRR {
		case 7, 9:
			ab01an01Alternative = "date"
		case 8, 10:
			ab01an01Alternative = "time"
		}
	}
	return ab01an01Alternative != "" && ab01an01Input.Ab01An01Alternative == ab01an01Alternative
}
