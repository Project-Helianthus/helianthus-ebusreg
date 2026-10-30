package system

type Ab03An03Input struct {
	Ab03An03Family        string
	Ab03An03Software      string
	Ab03An03Hardware      string
	Ab03An03APIVersion    string
	Ab03An03APIRevision   string
	Ab03An03OP            byte
	Ab03An03GG            byte
	Ab03An03II            *byte
	Ab03An03RR            uint16
	Ab03An03Width         string
	Ab03An03RequestedGG   byte
	Ab03An03RequestedRR   uint16
	Ab03An03ReturnedGG    byte
	Ab03An03ReturnedRR    uint16
	Ab03An03ResponseState string
	Ab03An03Raw           []byte
}

type Ab03An03Result struct {
	Ab03An03Context             Ab03An03Input
	Ab03An03Integer             *int
	Ab03An03CandidateLabel      *string
	Ab03An03Applicable          bool
	Ab03An03MalformedWidth      bool
	Ab03An03Grade               string
	Ab03An03NativeQualification string
}

func Ab03An03Interpret(ab03an03Input Ab03An03Input) Ab03An03Result {
	ab03an03Input.Ab03An03Raw = append([]byte(nil), ab03an03Input.Ab03An03Raw...)
	if ab03an03Input.Ab03An03II != nil {
		ab03an03II := *ab03an03Input.Ab03An03II
		ab03an03Input.Ab03An03II = &ab03an03II
	}
	ab03an03Result := Ab03An03Result{Ab03An03Context: ab03an03Input, Ab03An03Grade: "Unknown", Ab03An03NativeQualification: "Unknown"}
	ab03an03Size := 0
	switch ab03an03Input.Ab03An03Width {
	case "u8", "signed8":
		ab03an03Size = 1
	case "u16-le":
		ab03an03Size = 2
	default:
		return ab03an03Result
	}
	if len(ab03an03Input.Ab03An03Raw) != ab03an03Size {
		ab03an03Result.Ab03An03MalformedWidth = true
		return ab03an03Result
	}
	ab03an03Value := int(ab03an03Input.Ab03An03Raw[0])
	if ab03an03Size == 2 {
		ab03an03Value |= int(ab03an03Input.Ab03An03Raw[1]) << 8
	}
	if ab03an03Input.Ab03An03Width == "signed8" {
		ab03an03Value = int(int8(ab03an03Input.Ab03An03Raw[0]))
	}
	ab03an03Result.Ab03An03Integer = &ab03an03Value
	if (ab03an03Input.Ab03An03Family != "controller-700" && ab03an03Input.Ab03An03Family != "controller-720") ||
		ab03an03Input.Ab03An03OP != 2 || ab03an03Input.Ab03An03II == nil ||
		ab03an03Input.Ab03An03ResponseState != "echo-payload" ||
		ab03an03Input.Ab03An03RequestedGG != ab03an03Input.Ab03An03GG || ab03an03Input.Ab03An03RequestedRR != ab03an03Input.Ab03An03RR ||
		ab03an03Input.Ab03An03ReturnedGG != ab03an03Input.Ab03An03GG || ab03an03Input.Ab03An03ReturnedRR != ab03an03Input.Ab03An03RR {
		return ab03an03Result
	}
	var ab03an03Labels map[int]string
	ab03an03Family700 := ab03an03Input.Ab03An03Family == "controller-700"
	switch {
	case (ab03an03Input.Ab03An03GG == 1 && ab03an03Input.Ab03An03RR == 3 && *ab03an03Input.Ab03An03II == 0) ||
		(ab03an03Input.Ab03An03GG == 3 && (ab03an03Input.Ab03An03RR == 1 || ab03an03Input.Ab03An03RR == 6)):
		if ab03an03Input.Ab03An03Width != "u16-le" {
			return ab03an03Result
		}
		ab03an03Labels = map[int]string{0: "off", 1: "time-controlled", 2: "manual"}
		if ab03an03Family700 {
			ab03an03Labels = map[int]string{0: "off", 1: "automatic", 2: "day"}
			if ab03an03Input.Ab03An03GG == 3 && ab03an03Input.Ab03An03RR == 6 {
				ab03an03Labels[3] = "setback"
			}
		}
	case ab03an03Input.Ab03An03GG == 2 && ab03an03Input.Ab03An03RR == 2:
		if ab03an03Input.Ab03An03Width != "u16-le" {
			return ab03an03Result
		}
		ab03an03Labels = map[int]string{0: "deactivated", 1: "heating", 2: "fixed-flow", 3: "DHW", 4: "return-flow"}
		if ab03an03Family700 {
			ab03an03Labels[5] = "pool"
		}
	case ab03an03Input.Ab03An03GG == 3 && ab03an03Input.Ab03An03RR == 0x13:
		if ab03an03Input.Ab03An03Width != "u16-le" {
			return ab03an03Result
		}
		ab03an03Labels = map[int]string{0: "none", 1: "central", 2: "remote-circuit-one", 3: "remote-circuit-two", 4: "remote-circuit-three", 5: "remote-circuit-four"}
		if ab03an03Family700 {
			ab03an03Labels[6] = "remote-circuit-five"
			ab03an03Labels[7] = "remote-circuit-six"
			ab03an03Labels[8] = "remote-circuit-seven"
			ab03an03Labels[9] = "remote-circuit-eight"
		}
	case ab03an03Input.Ab03An03GG == 3 && ab03an03Input.Ab03An03RR == 0x1b:
		if ab03an03Input.Ab03An03Width != "u8" && ab03an03Input.Ab03An03Width != "u16-le" {
			return ab03an03Result
		}
		ab03an03Labels = map[int]string{0: "no-special-function", 1: "ventilation-boost", 4: "quick-veto", 7: "holiday", 10: "system-off"}
		if ab03an03Family700 {
			ab03an03Labels = map[int]string{0: "no-quick-mode", 1: "ventilation-boost", 3: "party", 4: "quick-veto", 5: "one-day-holiday", 6: "one-day-bank-holiday", 7: "holiday", 8: "bank-holiday", 10: "system-off"}
		}
	case ab03an03Input.Ab03An03GG == 3 && ab03an03Input.Ab03An03RR == 0x1c:
		if ab03an03Input.Ab03An03Width != "signed8" {
			return ab03an03Result
		}
		ab03an03Labels = map[int]string{-1: "unbound", 0: "control-circuit-one", 1: "control-circuit-two", 2: "control-circuit-three"}
		if !ab03an03Family700 {
			ab03an03Labels[3] = "control-circuit-four"
			ab03an03Labels[4] = "control-circuit-five"
			ab03an03Labels[5] = "control-circuit-six"
			ab03an03Labels[6] = "control-circuit-seven"
			ab03an03Labels[7] = "control-circuit-eight"
			ab03an03Labels[8] = "control-circuit-nine"
		}
	default:
		return ab03an03Result
	}
	ab03an03Result.Ab03An03Applicable = true
	ab03an03Result.Ab03An03Grade = "Hypothesis"
	if ab03an03Label, ab03an03OK := ab03an03Labels[ab03an03Value]; ab03an03OK {
		ab03an03Result.Ab03An03CandidateLabel = &ab03an03Label
	}
	return ab03an03Result
}
