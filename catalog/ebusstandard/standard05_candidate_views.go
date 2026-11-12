package ebusstandard

type Standard05CandidateViewState string

type Standard05CandidateViewKey struct {
	Standard05CandidateViewFamily    string `json:"Ab78Family"`
	Standard05CandidateViewVariant   string `json:"Ab78Variant"`
	Standard05CandidateViewPrimary   *uint8 `json:"Ab78Primary"`
	Standard05CandidateViewSecondary *uint8 `json:"Ab78Secondary"`
	Standard05CandidateViewOperation *uint8 `json:"Ab78Operation"`
	Standard05CandidateViewDirection string `json:"Ab78Direction"`
	Standard05CandidateViewSelector  string `json:"Ab78Selector"`
}

type Standard05CandidateViewField struct {
	Standard05CandidateViewLabel          string                       `json:"Ab78Label"`
	Standard05CandidateViewRaw            *uint8                       `json:"Ab78Raw"`
	Standard05CandidateViewValue          *int                         `json:"Ab78Value"`
	Standard05CandidateViewName           string                       `json:"Ab78Name"`
	Standard05CandidateViewInterpretation Standard05CandidateViewState `json:"Ab78Interpretation"`
	Standard05CandidateViewDomain         Standard05CandidateViewState `json:"Ab78Domain"`
	Standard05CandidateViewDefault        int                          `json:"Ab78Default"`
	Standard05CandidateViewMinimum        int                          `json:"Ab78Minimum"`
	Standard05CandidateViewMaximum        int                          `json:"Ab78Maximum"`
}

type Standard05CandidateView struct {
	Standard05CandidateViewKey          Standard05CandidateViewKey      `json:"Ab78Key"`
	Standard05CandidateViewRaw          []byte                          `json:"Ab78Raw"`
	Standard05CandidateViewRemainder    []byte                          `json:"Ab78Remainder"`
	Standard05CandidateViewState        Standard05CandidateViewState    `json:"Ab78State"`
	Standard05CandidateViewComplete     bool                            `json:"Ab78Complete"`
	Standard05CandidateViewFields       [5]Standard05CandidateViewField `json:"Ab78Fields"`
	Standard05CandidateViewAPIGrade     string                          `json:"Ab78APIGrade"`
	Standard05CandidateViewMeaningGrade string                          `json:"Ab78MeaningGrade"`
	Standard05CandidateViewNativeGrade  string                          `json:"Ab78NativeGrade"`
	Standard05CandidateViewUnknowns     []string                        `json:"Ab78Unknowns"`
}

func standard05CandidateViewCopyByte(standard05CandidateViewByte *uint8) *uint8 {
	if standard05CandidateViewByte == nil {
		return nil
	}
	standard05CandidateViewValue := *standard05CandidateViewByte
	return &standard05CandidateViewValue
}

func Standard05CandidateViewInterpret(standard05CandidateViewKey Standard05CandidateViewKey, standard05CandidateViewBody []byte) Standard05CandidateView {
	standard05CandidateViewKey.Standard05CandidateViewPrimary = standard05CandidateViewCopyByte(standard05CandidateViewKey.Standard05CandidateViewPrimary)
	standard05CandidateViewKey.Standard05CandidateViewSecondary = standard05CandidateViewCopyByte(standard05CandidateViewKey.Standard05CandidateViewSecondary)
	standard05CandidateViewKey.Standard05CandidateViewOperation = standard05CandidateViewCopyByte(standard05CandidateViewKey.Standard05CandidateViewOperation)
	standard05CandidateViewView := Standard05CandidateView{
		Standard05CandidateViewKey:       standard05CandidateViewKey,
		Standard05CandidateViewRaw:       append([]byte(nil), standard05CandidateViewBody...),
		Standard05CandidateViewRemainder: append([]byte(nil), standard05CandidateViewBody...),
		Standard05CandidateViewState:     "raw-only",
		Standard05CandidateViewAPIGrade:  "Proven", Standard05CandidateViewMeaningGrade: "Hypothesis", Standard05CandidateViewNativeGrade: "Unknown",
		Standard05CandidateViewUnknowns: []string{"physical support", "presence", "freshness", "physical length", "delivery", "effects", "persistence", "execution", "ACK", "command acceptance", "empty wire response"},
	}
	if standard05CandidateViewKey.Standard05CandidateViewFamily != "burner" || standard05CandidateViewKey.Standard05CandidateViewVariant != "0501" ||
		standard05CandidateViewKey.Standard05CandidateViewPrimary == nil || *standard05CandidateViewKey.Standard05CandidateViewPrimary != 5 ||
		standard05CandidateViewKey.Standard05CandidateViewSecondary == nil || *standard05CandidateViewKey.Standard05CandidateViewSecondary != 1 ||
		standard05CandidateViewKey.Standard05CandidateViewOperation != nil || standard05CandidateViewKey.Standard05CandidateViewDirection != "abstract-request" ||
		standard05CandidateViewKey.Standard05CandidateViewSelector != "burner-five-field" {
		return standard05CandidateViewView
	}
	standard05CandidateViewView.Standard05CandidateViewState = "candidate-body"
	standard05CandidateViewView.Standard05CandidateViewComplete = len(standard05CandidateViewBody) == 5
	standard05CandidateViewView.Standard05CandidateViewRemainder = nil
	if len(standard05CandidateViewBody) > 5 {
		standard05CandidateViewView.Standard05CandidateViewRemainder = append([]byte(nil), standard05CandidateViewBody[5:]...)
	}
	standard05CandidateViewView.Standard05CandidateViewFields = [5]Standard05CandidateViewField{
		{Standard05CandidateViewLabel: "heat-state", Standard05CandidateViewMaximum: 255},
		{Standard05CandidateViewLabel: "boiler-target", Standard05CandidateViewMaximum: 90},
		{Standard05CandidateViewLabel: "DHW-target", Standard05CandidateViewMaximum: 90},
		{Standard05CandidateViewLabel: "outside-temperature", Standard05CandidateViewMinimum: -128, Standard05CandidateViewMaximum: 127},
		{Standard05CandidateViewLabel: "setting", Standard05CandidateViewMaximum: 255},
	}
	for standard05CandidateViewIndex := range standard05CandidateViewView.Standard05CandidateViewFields {
		standard05CandidateViewField := &standard05CandidateViewView.Standard05CandidateViewFields[standard05CandidateViewIndex]
		standard05CandidateViewField.Standard05CandidateViewInterpretation = "absent"
		standard05CandidateViewField.Standard05CandidateViewDomain = "unknown"
		if standard05CandidateViewIndex >= len(standard05CandidateViewBody) {
			continue
		}
		standard05CandidateViewRaw := standard05CandidateViewBody[standard05CandidateViewIndex]
		standard05CandidateViewField.Standard05CandidateViewRaw = &standard05CandidateViewRaw
		standard05CandidateViewValue := int(standard05CandidateViewRaw)
		if standard05CandidateViewIndex == 3 {
			standard05CandidateViewValue = int(int8(standard05CandidateViewRaw))
		}
		standard05CandidateViewField.Standard05CandidateViewInterpretation = "interpreted"
		standard05CandidateViewField.Standard05CandidateViewDomain = "within-candidate-domain"
		if standard05CandidateViewIndex == 0 {
			switch standard05CandidateViewRaw {
			case 0x00:
				standard05CandidateViewField.Standard05CandidateViewName = "shutdown"
			case 0x55:
				standard05CandidateViewField.Standard05CandidateViewName = "DHW"
			case 0xaa:
				standard05CandidateViewField.Standard05CandidateViewName = "heating"
			case 0xcc:
				standard05CandidateViewField.Standard05CandidateViewName = "emission-check"
			case 0xee:
				standard05CandidateViewField.Standard05CandidateViewName = "controller-stop"
			default:
				standard05CandidateViewField.Standard05CandidateViewInterpretation = "unmapped"
				standard05CandidateViewField.Standard05CandidateViewDomain = "unknown"
				continue
			}
		}
		standard05CandidateViewField.Standard05CandidateViewValue = &standard05CandidateViewValue
		if standard05CandidateViewValue < standard05CandidateViewField.Standard05CandidateViewMinimum || standard05CandidateViewValue > standard05CandidateViewField.Standard05CandidateViewMaximum {
			standard05CandidateViewField.Standard05CandidateViewDomain = "out-of-candidate-domain"
		}
	}
	return standard05CandidateViewView
}
