package ebusstandard

type Standard09BodyViewKey struct {
	Standard09BodyViewFamily    string `json:"Ab25An25Family"`
	Standard09BodyViewProfile   string `json:"Ab25An25Profile"`
	Standard09BodyViewService   uint16 `json:"Ab25An25Service"`
	Standard09BodyViewDirection string `json:"Ab25An25Direction"`
}

type Standard09BodyViewState struct {
	Standard09BodyViewKey                          Standard09BodyViewKey `json:"Ab25An25Key"`
	Standard09BodyViewRaw                          []byte                `json:"Ab25An25Raw"`
	Standard09BodyViewState                        string                `json:"Ab25An25State"`
	Standard09BodyViewContractGrade                string                `json:"Ab25An25ContractGrade"`
	Standard09BodyViewMeaningGrade                 string                `json:"Ab25An25MeaningGrade"`
	Standard09BodyViewExecutionGrade               string                `json:"Ab25An25ExecutionGrade"`
	Standard09BodyViewPhysicalSupportGrade         string                `json:"Ab25An25PhysicalSupportGrade"`
	Standard09BodyViewPhysicalOffsetsGrade         string                `json:"Ab25An25PhysicalOffsetsGrade"`
	Standard09BodyViewPositivePrefixRelevanceGrade string                `json:"Ab25An25PositivePrefixRelevanceGrade"`
}

type Standard09BodyViewBody struct {
	Standard09BodyViewState          Standard09BodyViewState `json:"Ab25An25State"`
	Standard09BodyViewData           []byte                  `json:"Ab25An25Data"`
	Standard09BodyViewUnresolvedTail []byte                  `json:"Ab25An25UnresolvedTail"`
}

type Standard09BodyViewAddress struct {
	Standard09BodyViewState Standard09BodyViewState `json:"Ab25An25State"`
	Standard09BodyViewLE    *uint16                 `json:"Ab25An25LE"`
	Standard09BodyViewLow   *uint8                  `json:"Ab25An25Low"`
	Standard09BodyViewHigh  *uint8                  `json:"Ab25An25High"`
}

type Standard09BodyViewCount struct {
	Standard09BodyViewState            Standard09BodyViewState `json:"Ab25An25State"`
	Standard09BodyViewValue            uint8                   `json:"Ab25An25Value"`
	Standard09BodyViewCandidateMaximum *uint8                  `json:"Ab25An25CandidateMaximum"`
	Standard09BodyViewConflict         bool                    `json:"Ab25An25Conflict"`
}

func standard09BodyViewState(standard09BodyViewKey Standard09BodyViewKey, standard09BodyViewRaw []byte) Standard09BodyViewState {
	return Standard09BodyViewState{
		Standard09BodyViewKey: standard09BodyViewKey, Standard09BodyViewRaw: append([]byte(nil), standard09BodyViewRaw...),
		Standard09BodyViewState: "unknown", Standard09BodyViewContractGrade: "Proven", Standard09BodyViewMeaningGrade: "Hypothesis",
		Standard09BodyViewExecutionGrade: "Unknown", Standard09BodyViewPhysicalSupportGrade: "Unknown",
		Standard09BodyViewPhysicalOffsetsGrade: "Unknown", Standard09BodyViewPositivePrefixRelevanceGrade: "Unknown",
	}
}

func standard09BodyViewMatches(standard09BodyViewSelected, standard09BodyViewKey Standard09BodyViewKey, standard09BodyViewDirection string) bool {
	return standard09BodyViewSelected == standard09BodyViewKey && standard09BodyViewSelected.Standard09BodyViewFamily != "" &&
		standard09BodyViewSelected.Standard09BodyViewProfile != "" && standard09BodyViewSelected.Standard09BodyViewDirection == standard09BodyViewDirection &&
		standard09BodyViewSelected.Standard09BodyViewService >= 0x0900 && standard09BodyViewSelected.Standard09BodyViewService <= 0x0903
}

// Body views match the positive direction; address and count views match the
// abstract-request direction that carries those request fields.
func Standard09ViewBody(standard09BodyViewSelected, standard09BodyViewKey Standard09BodyViewKey, standard09BodyViewRaw []byte) Standard09BodyViewBody {
	standard09BodyViewView := Standard09BodyViewBody{Standard09BodyViewState: standard09BodyViewState(standard09BodyViewKey, standard09BodyViewRaw)}
	if !standard09BodyViewMatches(standard09BodyViewSelected, standard09BodyViewKey, "positive") {
		return standard09BodyViewView
	}
	if standard09BodyViewKey.Standard09BodyViewService == 0x0902 {
		standard09BodyViewView.Standard09BodyViewState.Standard09BodyViewState = "candidate"
		standard09BodyViewView.Standard09BodyViewData = append([]byte{}, standard09BodyViewRaw...)
		return standard09BodyViewView
	}
	if len(standard09BodyViewRaw) < 2 || standard09BodyViewRaw[0] != 9 || standard09BodyViewRaw[1] != byte(standard09BodyViewKey.Standard09BodyViewService) {
		standard09BodyViewView.Standard09BodyViewState.Standard09BodyViewState = "incomplete"
		return standard09BodyViewView
	}
	standard09BodyViewView.Standard09BodyViewState.Standard09BodyViewState = "candidate"
	if standard09BodyViewKey.Standard09BodyViewService == 0x0901 {
		standard09BodyViewView.Standard09BodyViewUnresolvedTail = append([]byte(nil), standard09BodyViewRaw[2:]...)
		if len(standard09BodyViewRaw) > 2 {
			standard09BodyViewView.Standard09BodyViewState.Standard09BodyViewState = "unresolved"
		}
		return standard09BodyViewView
	}
	standard09BodyViewView.Standard09BodyViewData = append([]byte{}, standard09BodyViewRaw[2:]...)
	return standard09BodyViewView
}

func Standard09ViewAddress(standard09BodyViewSelected, standard09BodyViewKey Standard09BodyViewKey, standard09BodyViewField []byte) Standard09BodyViewAddress {
	standard09BodyViewView := Standard09BodyViewAddress{Standard09BodyViewState: standard09BodyViewState(standard09BodyViewKey, standard09BodyViewField)}
	if !standard09BodyViewMatches(standard09BodyViewSelected, standard09BodyViewKey, "abstract-request") {
		return standard09BodyViewView
	}
	if len(standard09BodyViewField) != 2 {
		standard09BodyViewView.Standard09BodyViewState.Standard09BodyViewState = "incomplete"
		return standard09BodyViewView
	}
	standard09BodyViewView.Standard09BodyViewState.Standard09BodyViewState = "candidate"
	if standard09BodyViewKey.Standard09BodyViewService <= 0x0901 {
		standard09BodyViewLE := uint16(standard09BodyViewField[0]) | uint16(standard09BodyViewField[1])<<8
		standard09BodyViewView.Standard09BodyViewLE = &standard09BodyViewLE
	} else {
		standard09BodyViewLow, standard09BodyViewHigh := standard09BodyViewField[0], standard09BodyViewField[1]
		standard09BodyViewView.Standard09BodyViewLow, standard09BodyViewView.Standard09BodyViewHigh = &standard09BodyViewLow, &standard09BodyViewHigh
	}
	return standard09BodyViewView
}

func Standard09ViewCount(standard09BodyViewSelected, standard09BodyViewKey Standard09BodyViewKey, standard09BodyViewCount uint8) Standard09BodyViewCount {
	standard09BodyViewView := Standard09BodyViewCount{Standard09BodyViewState: standard09BodyViewState(standard09BodyViewKey, []byte{standard09BodyViewCount}), Standard09BodyViewValue: standard09BodyViewCount}
	if !standard09BodyViewMatches(standard09BodyViewSelected, standard09BodyViewKey, "abstract-request") || (standard09BodyViewKey.Standard09BodyViewService != 0x0900 && standard09BodyViewKey.Standard09BodyViewService != 0x0902) {
		return standard09BodyViewView
	}
	standard09BodyViewView.Standard09BodyViewState.Standard09BodyViewState = "candidate"
	if standard09BodyViewKey.Standard09BodyViewService == 0x0902 {
		standard09BodyViewMaximum := uint8(10)
		standard09BodyViewView.Standard09BodyViewCandidateMaximum = &standard09BodyViewMaximum
		standard09BodyViewView.Standard09BodyViewConflict = standard09BodyViewCount > 10
	}
	return standard09BodyViewView
}
