package vaillant

// B503HistoryBranchContext is the family, variant, operation and direction of a B503 history record.
type B503HistoryBranchContext struct {
	B503HistoryBranchFamily    string `json:"family"`
	B503HistoryBranchVariant   string `json:"variant"`
	B503HistoryBranchOperation uint8  `json:"operation"`
	B503HistoryBranchDirection string `json:"direction"`
}

// B503HistoryBranchCell is one calendar cell of a history record with its raw and normalized values
// and their presence.
type B503HistoryBranchCell struct {
	B503HistoryBranchRawPresent        bool   `json:"raw_present"`
	B503HistoryBranchRaw               uint8  `json:"raw"`
	B503HistoryBranchNormalizedPresent bool   `json:"normalized_present"`
	B503HistoryBranchValue             uint64 `json:"value"`
}

// B503HistoryBranchFields are the decoded fields of one history record: the branch marker, the
// remainder, the error code and the calendar cells.
type B503HistoryBranchFields struct {
	B503HistoryBranchContext          B503HistoryBranchContext `json:"context"`
	B503HistoryBranchMarkerPresent    bool                     `json:"marker_present"`
	B503HistoryBranchMarker           uint8                    `json:"marker"`
	B503HistoryBranchRemainderPresent bool                     `json:"remainder_present"`
	B503HistoryBranchRemainder        [7]uint8                 `json:"remainder"`
	B503HistoryBranchErrorPresent     bool                     `json:"error_present"`
	B503HistoryBranchError            uint64                   `json:"error"`
	B503HistoryBranchCalendar         [5]B503HistoryBranchCell `json:"calendar"`
}

// B503HistoryBranchComponent is one calendar cell with its interpretation state.
type B503HistoryBranchComponent struct {
	B503HistoryBranchCell  B503HistoryBranchCell `json:"cell"`
	B503HistoryBranchState string                `json:"state"`
}

// B503HistoryBranchView is the branch interpretation of a history record with its grades.
type B503HistoryBranchView struct {
	B503HistoryBranchFields        B503HistoryBranchFields       `json:"fields"`
	B503HistoryBranchState         string                        `json:"state"`
	B503HistoryBranchComponents    [5]B503HistoryBranchComponent `json:"components"`
	B503HistoryBranchComplete      bool                          `json:"complete"`
	B503HistoryBranchContract      string                        `json:"contract"`
	B503HistoryBranchMeaning       string                        `json:"meaning"`
	B503HistoryBranchQualification string                        `json:"qualification"`
}

// B503HistoryBranchHistory selects the record branch from the marker: 0 is a no-item record when the
// remainder is all 0xFF, 1 an error record without complete time, and 2 a record with candidate
// calendar components. Any other marker or incomplete input leaves the state unknown-branch.
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
