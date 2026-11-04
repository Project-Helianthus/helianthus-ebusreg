package ebus_standard_catalog

type Ab25An25Key struct {
	Ab25An25Family    string
	Ab25An25Profile   string
	Ab25An25Service   uint16
	Ab25An25Direction string
}

type Ab25An25State struct {
	Ab25An25Key                          Ab25An25Key
	Ab25An25Raw                          []byte
	Ab25An25State                        string
	Ab25An25ContractGrade                string
	Ab25An25MeaningGrade                 string
	Ab25An25ExecutionGrade               string
	Ab25An25PhysicalSupportGrade         string
	Ab25An25PhysicalOffsetsGrade         string
	Ab25An25PositivePrefixRelevanceGrade string
}

type Ab25An25Body struct {
	Ab25An25State          Ab25An25State
	Ab25An25Data           []byte
	Ab25An25UnresolvedTail []byte
}

type Ab25An25Address struct {
	Ab25An25State Ab25An25State
	Ab25An25LE    *uint16
	Ab25An25Low   *uint8
	Ab25An25High  *uint8
}

type Ab25An25Count struct {
	Ab25An25State            Ab25An25State
	Ab25An25Value            uint8
	Ab25An25CandidateMaximum *uint8
	Ab25An25Conflict         bool
}

func ab25an25State(ab25an25Key Ab25An25Key, ab25an25Raw []byte) Ab25An25State {
	return Ab25An25State{
		Ab25An25Key: ab25an25Key, Ab25An25Raw: append([]byte(nil), ab25an25Raw...),
		Ab25An25State: "unknown", Ab25An25ContractGrade: "Proven", Ab25An25MeaningGrade: "Hypothesis",
		Ab25An25ExecutionGrade: "Unknown", Ab25An25PhysicalSupportGrade: "Unknown",
		Ab25An25PhysicalOffsetsGrade: "Unknown", Ab25An25PositivePrefixRelevanceGrade: "Unknown",
	}
}

func ab25an25Matches(ab25an25Selected, ab25an25Key Ab25An25Key) bool {
	return ab25an25Selected == ab25an25Key && ab25an25Selected.Ab25An25Family != "" &&
		ab25an25Selected.Ab25An25Profile != "" && ab25an25Selected.Ab25An25Direction == "positive" &&
		ab25an25Selected.Ab25An25Service >= 0x0900 && ab25an25Selected.Ab25An25Service <= 0x0903
}

func Ab25An25ViewBody(ab25an25Selected, ab25an25Key Ab25An25Key, ab25an25Raw []byte) Ab25An25Body {
	ab25an25View := Ab25An25Body{Ab25An25State: ab25an25State(ab25an25Key, ab25an25Raw)}
	if !ab25an25Matches(ab25an25Selected, ab25an25Key) {
		return ab25an25View
	}
	if ab25an25Key.Ab25An25Service == 0x0902 {
		ab25an25View.Ab25An25State.Ab25An25State = "candidate"
		ab25an25View.Ab25An25Data = append([]byte{}, ab25an25Raw...)
		return ab25an25View
	}
	if len(ab25an25Raw) < 2 || ab25an25Raw[0] != 9 || ab25an25Raw[1] != byte(ab25an25Key.Ab25An25Service) {
		ab25an25View.Ab25An25State.Ab25An25State = "incomplete"
		return ab25an25View
	}
	ab25an25View.Ab25An25State.Ab25An25State = "candidate"
	if ab25an25Key.Ab25An25Service == 0x0901 {
		ab25an25View.Ab25An25UnresolvedTail = append([]byte(nil), ab25an25Raw[2:]...)
		if len(ab25an25Raw) > 2 {
			ab25an25View.Ab25An25State.Ab25An25State = "unresolved"
		}
		return ab25an25View
	}
	ab25an25View.Ab25An25Data = append([]byte{}, ab25an25Raw[2:]...)
	return ab25an25View
}

func Ab25An25ViewAddress(ab25an25Selected, ab25an25Key Ab25An25Key, ab25an25Field []byte) Ab25An25Address {
	ab25an25View := Ab25An25Address{Ab25An25State: ab25an25State(ab25an25Key, ab25an25Field)}
	if !ab25an25Matches(ab25an25Selected, ab25an25Key) {
		return ab25an25View
	}
	if len(ab25an25Field) != 2 {
		ab25an25View.Ab25An25State.Ab25An25State = "incomplete"
		return ab25an25View
	}
	ab25an25View.Ab25An25State.Ab25An25State = "candidate"
	if ab25an25Key.Ab25An25Service <= 0x0901 {
		ab25an25LE := uint16(ab25an25Field[0]) | uint16(ab25an25Field[1])<<8
		ab25an25View.Ab25An25LE = &ab25an25LE
	} else {
		ab25an25Low, ab25an25High := ab25an25Field[0], ab25an25Field[1]
		ab25an25View.Ab25An25Low, ab25an25View.Ab25An25High = &ab25an25Low, &ab25an25High
	}
	return ab25an25View
}

func Ab25An25ViewCount(ab25an25Selected, ab25an25Key Ab25An25Key, ab25an25Count uint8) Ab25An25Count {
	ab25an25View := Ab25An25Count{Ab25An25State: ab25an25State(ab25an25Key, []byte{ab25an25Count}), Ab25An25Value: ab25an25Count}
	if !ab25an25Matches(ab25an25Selected, ab25an25Key) || (ab25an25Key.Ab25An25Service != 0x0900 && ab25an25Key.Ab25An25Service != 0x0902) {
		return ab25an25View
	}
	ab25an25View.Ab25An25State.Ab25An25State = "candidate"
	if ab25an25Key.Ab25An25Service == 0x0902 {
		ab25an25Maximum := uint8(10)
		ab25an25View.Ab25An25CandidateMaximum = &ab25an25Maximum
		ab25an25View.Ab25An25Conflict = ab25an25Count > 10
	}
	return ab25an25View
}
