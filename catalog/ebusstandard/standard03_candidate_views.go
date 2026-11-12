package ebusstandard

type Standard03CandidateViewState string

type Standard03CandidateViewEncoding string

const (
	Standard03CandidateViewRawOnly             Standard03CandidateViewState    = "raw-only"
	Standard03CandidateViewIncomplete          Standard03CandidateViewState    = "incomplete"
	Standard03CandidateViewExtra               Standard03CandidateViewState    = "extra"
	Standard03CandidateViewInvalid             Standard03CandidateViewState    = "invalid"
	Standard03CandidateViewComplete            Standard03CandidateViewState    = "complete-candidate"
	Standard03CandidateViewAbsent              Standard03CandidateViewState    = "absent"
	Standard03CandidateViewValid               Standard03CandidateViewState    = "valid"
	Standard03CandidateViewEncodingUnspecified Standard03CandidateViewEncoding = ""
	Standard03CandidateViewDecimalChunks       Standard03CandidateViewEncoding = "unsigned-decimal-chunks"
	Standard03CandidateViewUnsignedU8          Standard03CandidateViewEncoding = "unsigned-u8"
)

type Standard03CandidateViewSelection struct {
	Standard03CandidateViewFamily    string                          `json:"Ab76Family"`
	Standard03CandidateViewVariant   string                          `json:"Ab76Variant"`
	Standard03CandidateViewFunction  *uint8                          `json:"Ab76Function"`
	Standard03CandidateViewOperation *uint8                          `json:"Ab76Operation"`
	Standard03CandidateViewDirection string                          `json:"Ab76Direction"`
	Standard03CandidateViewSelector  string                          `json:"Ab76Selector"`
	Standard03CandidateViewEncoding  Standard03CandidateViewEncoding `json:"Ab76Encoding"`
}

type Standard03CandidateViewComponent struct {
	Standard03CandidateViewName     string                          `json:"Ab76Name"`
	Standard03CandidateViewState    Standard03CandidateViewState    `json:"Ab76State"`
	Standard03CandidateViewRaw      *uint8                          `json:"Ab76Raw"`
	Standard03CandidateViewValue    *uint8                          `json:"Ab76Value"`
	Standard03CandidateViewWidth    int                             `json:"Ab76Width"`
	Standard03CandidateViewMaximum  uint8                           `json:"Ab76Maximum"`
	Standard03CandidateViewDefault  *uint8                          `json:"Ab76Default"`
	Standard03CandidateViewEncoding Standard03CandidateViewEncoding `json:"Ab76Encoding"`
}

type Standard03CandidateView struct {
	Standard03CandidateViewSelection           Standard03CandidateViewSelection   `json:"Ab76Selection"`
	Standard03CandidateViewState               Standard03CandidateViewState       `json:"Ab76State"`
	Standard03CandidateViewRaw                 []byte                             `json:"Ab76Raw"`
	Standard03CandidateViewRemainder           []byte                             `json:"Ab76Remainder"`
	Standard03CandidateViewComponents          []Standard03CandidateViewComponent `json:"Ab76Components"`
	Standard03CandidateViewLocalCounter        *int                               `json:"Ab76LocalCounter"`
	Standard03CandidateViewLocalHours          *int                               `json:"Ab76LocalHours"`
	Standard03CandidateViewLocalMinute         *int                               `json:"Ab76LocalMinute"`
	Standard03CandidateViewImplementationGrade string                             `json:"Ab76ImplementationGrade"`
	Standard03CandidateViewMeaningGrade        string                             `json:"Ab76MeaningGrade"`
	Standard03CandidateViewNativeProperties    map[string]string                  `json:"Ab76NativeProperties"`
}

func standard03CandidateViewCopyByte(standard03CandidateViewValue *uint8) *uint8 {
	if standard03CandidateViewValue == nil {
		return nil
	}
	standard03CandidateViewCopy := *standard03CandidateViewValue
	return &standard03CandidateViewCopy
}

func standard03CandidateViewCopyBody(standard03CandidateViewBody []byte) []byte {
	if standard03CandidateViewBody == nil {
		return nil
	}
	standard03CandidateViewCopy := make([]byte, len(standard03CandidateViewBody))
	copy(standard03CandidateViewCopy, standard03CandidateViewBody)
	return standard03CandidateViewCopy
}

func Standard03CandidateViewInterpret(standard03CandidateViewSelection Standard03CandidateViewSelection, standard03CandidateViewBody []byte) Standard03CandidateView {
	standard03CandidateViewSelection.Standard03CandidateViewFunction = standard03CandidateViewCopyByte(standard03CandidateViewSelection.Standard03CandidateViewFunction)
	standard03CandidateViewSelection.Standard03CandidateViewOperation = standard03CandidateViewCopyByte(standard03CandidateViewSelection.Standard03CandidateViewOperation)
	standard03CandidateViewView := Standard03CandidateView{
		Standard03CandidateViewSelection: standard03CandidateViewSelection, Standard03CandidateViewState: Standard03CandidateViewRawOnly,
		Standard03CandidateViewRaw:                 standard03CandidateViewCopyBody(standard03CandidateViewBody),
		Standard03CandidateViewImplementationGrade: "Proven", Standard03CandidateViewMeaningGrade: "Hypothesis",
		Standard03CandidateViewNativeProperties: map[string]string{
			"physical-support": "Unknown", "presence": "Unknown", "freshness": "Unknown",
			"execution": "Unknown", "effects": "Unknown", "persistence": "Unknown",
			"delivery": "Unknown", "physical-length": "Unknown", "native-units": "Unknown",
			"transport-adaptation": "Unknown",
		},
	}
	if standard03CandidateViewSelection.Standard03CandidateViewFamily != "standard-counter" || standard03CandidateViewSelection.Standard03CandidateViewFunction == nil || *standard03CandidateViewSelection.Standard03CandidateViewFunction != 3 || standard03CandidateViewSelection.Standard03CandidateViewOperation == nil {
		return standard03CandidateViewView
	}
	standard03CandidateViewCandidate := standard03CandidateViewSelection.Standard03CandidateViewSelector == "candidate" && standard03CandidateViewSelection.Standard03CandidateViewDirection == "request"
	standard03CandidateViewStandard := standard03CandidateViewSelection.Standard03CandidateViewSelector == "standard" && standard03CandidateViewSelection.Standard03CandidateViewDirection == "positive"
	if !standard03CandidateViewCandidate && !standard03CandidateViewStandard {
		return standard03CandidateViewView
	}
	var standard03CandidateViewNames []string
	var standard03CandidateViewMaxima []uint8
	switch {
	case standard03CandidateViewSelection.Standard03CandidateViewVariant == "03/04" && *standard03CandidateViewSelection.Standard03CandidateViewOperation == 4 && standard03CandidateViewSelection.Standard03CandidateViewEncoding == Standard03CandidateViewDecimalChunks:
		standard03CandidateViewNames = []string{"c0", "c1", "c2"}
		standard03CandidateViewMaxima = []uint8{99, 99, 99}
	case standard03CandidateViewSelection.Standard03CandidateViewVariant == "03/05" && *standard03CandidateViewSelection.Standard03CandidateViewOperation == 5 && standard03CandidateViewCandidate && standard03CandidateViewSelection.Standard03CandidateViewEncoding == Standard03CandidateViewUnsignedU8:
		standard03CandidateViewNames = []string{"minutes", "hundreds-of-hours", "ten-thousands-of-hours"}
		standard03CandidateViewMaxima = []uint8{59, 59, 59}
	case standard03CandidateViewSelection.Standard03CandidateViewVariant == "03/06" && *standard03CandidateViewSelection.Standard03CandidateViewOperation == 6 && standard03CandidateViewSelection.Standard03CandidateViewEncoding == Standard03CandidateViewUnsignedU8:
		standard03CandidateViewNames = []string{"minutes", "h0", "h1", "h2"}
		standard03CandidateViewMaxima = []uint8{59, 99, 99, 99}
	default:
		return standard03CandidateViewView
	}
	standard03CandidateViewView.Standard03CandidateViewState = Standard03CandidateViewComplete
	if len(standard03CandidateViewBody) < len(standard03CandidateViewNames) {
		standard03CandidateViewView.Standard03CandidateViewState = Standard03CandidateViewIncomplete
	}
	if len(standard03CandidateViewBody) > len(standard03CandidateViewNames) {
		standard03CandidateViewView.Standard03CandidateViewState = Standard03CandidateViewExtra
		standard03CandidateViewView.Standard03CandidateViewRemainder = append([]byte(nil), standard03CandidateViewBody[len(standard03CandidateViewNames):]...)
	}
	for standard03CandidateViewIndex, standard03CandidateViewName := range standard03CandidateViewNames {
		standard03CandidateViewComponent := Standard03CandidateViewComponent{Standard03CandidateViewName: standard03CandidateViewName, Standard03CandidateViewState: Standard03CandidateViewAbsent, Standard03CandidateViewWidth: 1, Standard03CandidateViewMaximum: standard03CandidateViewMaxima[standard03CandidateViewIndex], Standard03CandidateViewEncoding: standard03CandidateViewSelection.Standard03CandidateViewEncoding}
		if standard03CandidateViewSelection.Standard03CandidateViewVariant == "03/04" {
			standard03CandidateViewZero := uint8(0)
			standard03CandidateViewComponent.Standard03CandidateViewDefault = &standard03CandidateViewZero
		}
		if standard03CandidateViewIndex < len(standard03CandidateViewBody) {
			standard03CandidateViewComponent.Standard03CandidateViewRaw = standard03CandidateViewCopyByte(&standard03CandidateViewBody[standard03CandidateViewIndex])
			standard03CandidateViewComponent.Standard03CandidateViewState = Standard03CandidateViewInvalid
			if standard03CandidateViewBody[standard03CandidateViewIndex] <= standard03CandidateViewMaxima[standard03CandidateViewIndex] {
				standard03CandidateViewComponent.Standard03CandidateViewState = Standard03CandidateViewValid
				standard03CandidateViewComponent.Standard03CandidateViewValue = standard03CandidateViewCopyByte(&standard03CandidateViewBody[standard03CandidateViewIndex])
			} else {
				standard03CandidateViewView.Standard03CandidateViewState = Standard03CandidateViewInvalid
			}
		}
		standard03CandidateViewView.Standard03CandidateViewComponents = append(standard03CandidateViewView.Standard03CandidateViewComponents, standard03CandidateViewComponent)
	}
	if standard03CandidateViewView.Standard03CandidateViewState == Standard03CandidateViewComplete {
		switch standard03CandidateViewSelection.Standard03CandidateViewVariant {
		case "03/04":
			standard03CandidateViewCounter := int(standard03CandidateViewBody[0]) + 100*int(standard03CandidateViewBody[1]) + 10000*int(standard03CandidateViewBody[2])
			standard03CandidateViewView.Standard03CandidateViewLocalCounter = &standard03CandidateViewCounter
		case "03/06":
			standard03CandidateViewMinute := int(standard03CandidateViewBody[0])
			standard03CandidateViewHours := int(standard03CandidateViewBody[1]) + 100*int(standard03CandidateViewBody[2]) + 10000*int(standard03CandidateViewBody[3])
			standard03CandidateViewView.Standard03CandidateViewLocalMinute, standard03CandidateViewView.Standard03CandidateViewLocalHours = &standard03CandidateViewMinute, &standard03CandidateViewHours
		}
	}
	return standard03CandidateViewView
}
