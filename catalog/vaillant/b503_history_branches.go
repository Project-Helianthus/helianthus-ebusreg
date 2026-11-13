package vaillant

type B503HistoryBranchContext struct {
	B503HistoryBranchFamily    string `json:"Ab05An05Family"`
	B503HistoryBranchVariant   string `json:"Ab05An05Variant"`
	B503HistoryBranchOperation uint8  `json:"Ab05An05Operation"`
	B503HistoryBranchDirection string `json:"Ab05An05Direction"`
}

type B503HistoryBranchCell struct {
	B503HistoryBranchRawPresent        bool   `json:"Ab05An05RawPresent"`
	B503HistoryBranchRaw               uint8  `json:"Ab05An05Raw"`
	B503HistoryBranchNormalizedPresent bool   `json:"Ab05An05NormalizedPresent"`
	B503HistoryBranchValue             uint64 `json:"Ab05An05Value"`
}

type B503HistoryBranchFields struct {
	B503HistoryBranchContext          B503HistoryBranchContext `json:"Ab05An05Context"`
	B503HistoryBranchMarkerPresent    bool                     `json:"Ab05An05MarkerPresent"`
	B503HistoryBranchMarker           uint8                    `json:"Ab05An05Marker"`
	B503HistoryBranchRemainderPresent bool                     `json:"Ab05An05RemainderPresent"`
	B503HistoryBranchRemainder        [7]uint8                 `json:"Ab05An05Remainder"`
	B503HistoryBranchErrorPresent     bool                     `json:"Ab05An05ErrorPresent"`
	B503HistoryBranchError            uint64                   `json:"Ab05An05Error"`
	B503HistoryBranchCalendar         [5]B503HistoryBranchCell `json:"Ab05An05Calendar"`
}

type B503HistoryBranchComponent struct {
	B503HistoryBranchCell  B503HistoryBranchCell `json:"Ab05An05Cell"`
	B503HistoryBranchState string                `json:"Ab05An05State"`
}

type B503HistoryBranchView struct {
	B503HistoryBranchFields        B503HistoryBranchFields       `json:"Ab05An05Fields"`
	B503HistoryBranchState         string                        `json:"Ab05An05State"`
	B503HistoryBranchComponents    [5]B503HistoryBranchComponent `json:"Ab05An05Components"`
	B503HistoryBranchComplete      bool                          `json:"Ab05An05Complete"`
	B503HistoryBranchContract      string                        `json:"Ab05An05Contract"`
	B503HistoryBranchMeaning       string                        `json:"Ab05An05Meaning"`
	B503HistoryBranchQualification string                        `json:"Ab05An05Qualification"`
}

func B503HistoryBranchHistory(b503HistoryBranchFields B503HistoryBranchFields) B503HistoryBranchView {
	b503HistoryBranchView := B503HistoryBranchView{
		B503HistoryBranchFields:        b503HistoryBranchFields,
		B503HistoryBranchState:         "unknown-branch",
		B503HistoryBranchContract:      "Proven",
		B503HistoryBranchMeaning:       "Unknown",
		B503HistoryBranchQualification: "Unknown",
	}
	for b503HistoryBranchIndex, b503HistoryBranchCell := range b503HistoryBranchFields.B503HistoryBranchCalendar {
		b503HistoryBranchView.B503HistoryBranchComponents[b503HistoryBranchIndex] = B503HistoryBranchComponent{B503HistoryBranchCell: b503HistoryBranchCell, B503HistoryBranchState: "unknown"}
	}
	if !b503HistoryBranchFields.B503HistoryBranchMarkerPresent || b503HistoryBranchFields.B503HistoryBranchContext.B503HistoryBranchFamily == "" || b503HistoryBranchFields.B503HistoryBranchContext.B503HistoryBranchVariant == "" || b503HistoryBranchFields.B503HistoryBranchContext.B503HistoryBranchDirection == "" || b503HistoryBranchFields.B503HistoryBranchContext.B503HistoryBranchOperation != 1 {
		return b503HistoryBranchView
	}
	switch b503HistoryBranchFields.B503HistoryBranchMarker {
	case 0:
		if !b503HistoryBranchFields.B503HistoryBranchRemainderPresent {
			return b503HistoryBranchView
		}
		for _, b503HistoryBranchByte := range b503HistoryBranchFields.B503HistoryBranchRemainder {
			if b503HistoryBranchByte != 255 {
				return b503HistoryBranchView
			}
		}
		b503HistoryBranchView.B503HistoryBranchState = "no-item"
	case 1, 2:
		b503HistoryBranchComplete := true
		b503HistoryBranchUnavailable := false
		for b503HistoryBranchIndex, b503HistoryBranchCell := range b503HistoryBranchFields.B503HistoryBranchCalendar {
			b503HistoryBranchState := "unknown"
			switch {
			case !b503HistoryBranchCell.B503HistoryBranchRawPresent:
				b503HistoryBranchState = "missing"
			case b503HistoryBranchCell.B503HistoryBranchRaw == 255:
				b503HistoryBranchState = "replacement"
			case !b503HistoryBranchCell.B503HistoryBranchNormalizedPresent:
				b503HistoryBranchState = "missing"
			case b503HistoryBranchFields.B503HistoryBranchMarker == 2:
				b503HistoryBranchState = "out-of-domain"
				if b503HistoryBranchCell.B503HistoryBranchValue >= [5]uint64{0, 0, 1, 1, 0}[b503HistoryBranchIndex] && b503HistoryBranchCell.B503HistoryBranchValue <= [5]uint64{59, 23, 31, 12, 99}[b503HistoryBranchIndex] {
					b503HistoryBranchState = "candidate"
				}
			}
			b503HistoryBranchView.B503HistoryBranchComponents[b503HistoryBranchIndex].B503HistoryBranchState = b503HistoryBranchState
			b503HistoryBranchComplete = b503HistoryBranchComplete && b503HistoryBranchState == "candidate"
			b503HistoryBranchUnavailable = b503HistoryBranchUnavailable || b503HistoryBranchState == "missing"
		}
		if b503HistoryBranchFields.B503HistoryBranchMarker == 1 {
			if !b503HistoryBranchFields.B503HistoryBranchErrorPresent {
				return b503HistoryBranchView
			}
			b503HistoryBranchView.B503HistoryBranchState = "error-without-complete-time"
		} else {
			if b503HistoryBranchUnavailable {
				return b503HistoryBranchView
			}
			b503HistoryBranchView.B503HistoryBranchState = "candidate-components"
			b503HistoryBranchView.B503HistoryBranchComplete = b503HistoryBranchComplete
		}
	default:
		return b503HistoryBranchView
	}
	b503HistoryBranchView.B503HistoryBranchMeaning = "Hypothesis"
	return b503HistoryBranchView
}
