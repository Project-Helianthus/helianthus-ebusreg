package system

type B524CalendarViewCell struct {
	B524CalendarViewRaw   byte   `json:"raw"`
	B524CalendarViewState string `json:"state"`
}

type B524CalendarViewInput struct {
	B524CalendarViewFamily        string                  `json:"family"`
	B524CalendarViewSoftware      string                  `json:"software"`
	B524CalendarViewHardware      string                  `json:"hardware"`
	B524CalendarViewAPIVersion    string                  `json:"api_version"`
	B524CalendarViewAPIRevision   string                  `json:"api_revision"`
	B524CalendarViewOP            byte                    `json:"op"`
	B524CalendarViewRequestedGG   byte                    `json:"requested_gg"`
	B524CalendarViewII            *byte                   `json:"ii"`
	B524CalendarViewRequestedRR   uint16                  `json:"requested_rr"`
	B524CalendarViewAlternative   string                  `json:"alternative"`
	B524CalendarViewReturnedGG    byte                    `json:"returned_gg"`
	B524CalendarViewReturnedRR    uint16                  `json:"returned_rr"`
	B524CalendarViewResponseState string                  `json:"response_state"`
	B524CalendarViewCells         [3]B524CalendarViewCell `json:"cells"`
	B524CalendarViewRaw           []byte                  `json:"raw"`
}

type B524CalendarViewComponent struct {
	B524CalendarViewCell      B524CalendarViewCell `json:"cell"`
	B524CalendarViewState     string               `json:"state"`
	B524CalendarViewCandidate *byte                `json:"candidate"`
}

type B524CalendarViewResult struct {
	B524CalendarViewContext             B524CalendarViewInput        `json:"context"`
	B524CalendarViewComponents          [3]B524CalendarViewComponent `json:"components"`
	B524CalendarViewCompleteCandidate   bool                         `json:"complete_candidate"`
	B524CalendarViewGrade               string                       `json:"grade"`
	B524CalendarViewContractGrade       string                       `json:"contract_grade"`
	B524CalendarViewNativeQualification string                       `json:"native_qualification"`
	B524CalendarViewCombinedTimestamp   string                       `json:"combined_timestamp"`
	B524CalendarViewEditValidation      string                       `json:"edit_validation"`
}

func B524CalendarViewInterpret(b524CalendarViewInput B524CalendarViewInput) B524CalendarViewResult {
	b524CalendarViewInput.B524CalendarViewRaw = append([]byte(nil), b524CalendarViewInput.B524CalendarViewRaw...)
	if b524CalendarViewInput.B524CalendarViewII != nil {
		b524CalendarViewII := *b524CalendarViewInput.B524CalendarViewII
		b524CalendarViewInput.B524CalendarViewII = &b524CalendarViewII
	}
	b524CalendarViewResult := B524CalendarViewResult{B524CalendarViewContext: b524CalendarViewInput, B524CalendarViewGrade: "Unknown", B524CalendarViewContractGrade: "Proven", B524CalendarViewNativeQualification: "Unknown", B524CalendarViewCombinedTimestamp: "Unknown", B524CalendarViewEditValidation: "Unknown"}
	for b524CalendarViewIndex, b524CalendarViewCell := range b524CalendarViewInput.B524CalendarViewCells {
		b524CalendarViewResult.B524CalendarViewComponents[b524CalendarViewIndex] = B524CalendarViewComponent{B524CalendarViewCell: b524CalendarViewCell, B524CalendarViewState: "Unknown"}
	}
	if !b524CalendarViewMatches(b524CalendarViewInput) {
		return b524CalendarViewResult
	}
	b524CalendarViewResult.B524CalendarViewGrade = "Hypothesis"
	b524CalendarViewResult.B524CalendarViewCompleteCandidate = true
	for b524CalendarViewIndex, b524CalendarViewCell := range b524CalendarViewInput.B524CalendarViewCells {
		b524CalendarViewComponent := &b524CalendarViewResult.B524CalendarViewComponents[b524CalendarViewIndex]
		if b524CalendarViewCell.B524CalendarViewState != "present" {
			if b524CalendarViewCell.B524CalendarViewState == "missing" || b524CalendarViewCell.B524CalendarViewState == "failed" {
				b524CalendarViewComponent.B524CalendarViewState = b524CalendarViewCell.B524CalendarViewState
			}
			b524CalendarViewResult.B524CalendarViewCompleteCandidate = false
			continue
		}
		if b524CalendarViewInput.B524CalendarViewAlternative == "date" && ((b524CalendarViewIndex == 0 && (b524CalendarViewCell.B524CalendarViewRaw < 1 || b524CalendarViewCell.B524CalendarViewRaw > 31)) || (b524CalendarViewIndex == 1 && (b524CalendarViewCell.B524CalendarViewRaw < 1 || b524CalendarViewCell.B524CalendarViewRaw > 12)) || (b524CalendarViewIndex == 2 && b524CalendarViewCell.B524CalendarViewRaw > 99)) {
			b524CalendarViewComponent.B524CalendarViewState = "out-of-domain"
			b524CalendarViewResult.B524CalendarViewCompleteCandidate = false
			continue
		}
		b524CalendarViewValue := b524CalendarViewCell.B524CalendarViewRaw
		b524CalendarViewComponent.B524CalendarViewCandidate = &b524CalendarViewValue
		b524CalendarViewComponent.B524CalendarViewState = "candidate"
	}
	return b524CalendarViewResult
}

func b524CalendarViewMatches(b524CalendarViewInput B524CalendarViewInput) bool {
	if (b524CalendarViewInput.B524CalendarViewFamily != "controller-700" && b524CalendarViewInput.B524CalendarViewFamily != "controller-720") || b524CalendarViewInput.B524CalendarViewSoftware == "" || b524CalendarViewInput.B524CalendarViewHardware == "" || b524CalendarViewInput.B524CalendarViewAPIVersion == "" || b524CalendarViewInput.B524CalendarViewAPIRevision == "" || b524CalendarViewInput.B524CalendarViewOP != 2 || b524CalendarViewInput.B524CalendarViewII == nil || b524CalendarViewInput.B524CalendarViewResponseState != "echo-payload" || b524CalendarViewInput.B524CalendarViewRequestedGG != b524CalendarViewInput.B524CalendarViewReturnedGG || b524CalendarViewInput.B524CalendarViewRequestedRR != b524CalendarViewInput.B524CalendarViewReturnedRR {
		return false
	}
	b524CalendarViewAlternative := ""
	switch b524CalendarViewInput.B524CalendarViewRequestedGG {
	case 0:
		if *b524CalendarViewInput.B524CalendarViewII != 0 {
			return false
		}
		switch b524CalendarViewInput.B524CalendarViewRequestedRR {
		case 0x34:
			b524CalendarViewAlternative = "date"
		case 0x35:
			b524CalendarViewAlternative = "time"
		case 0x71, 0x72:
			if b524CalendarViewInput.B524CalendarViewFamily == "controller-700" {
				b524CalendarViewAlternative = "date"
			}
		case 0xda, 0xdb:
			if b524CalendarViewInput.B524CalendarViewFamily == "controller-720" {
				b524CalendarViewAlternative = "date"
			}
		}
	case 1:
		if *b524CalendarViewInput.B524CalendarViewII != 0 {
			return false
		}
		switch b524CalendarViewInput.B524CalendarViewRequestedRR {
		case 9, 10:
			b524CalendarViewAlternative = "date"
		case 0x10, 0x11:
			if b524CalendarViewInput.B524CalendarViewFamily == "controller-720" {
				b524CalendarViewAlternative = "time"
			}
		}
	case 3:
		switch b524CalendarViewInput.B524CalendarViewRequestedRR {
		case 3, 4:
			b524CalendarViewAlternative = "date"
		case 0x24:
			if b524CalendarViewInput.B524CalendarViewFamily == "controller-720" {
				b524CalendarViewAlternative = "date"
			}
		case 0x1e, 0x20, 0x21:
			if b524CalendarViewInput.B524CalendarViewFamily == "controller-720" {
				b524CalendarViewAlternative = "time"
			}
		}
	case 9:
		if *b524CalendarViewInput.B524CalendarViewII != 0 || b524CalendarViewInput.B524CalendarViewFamily != "controller-720" {
			return false
		}
		switch b524CalendarViewInput.B524CalendarViewRequestedRR {
		case 7, 9:
			b524CalendarViewAlternative = "date"
		case 8, 10:
			b524CalendarViewAlternative = "time"
		}
	}
	return b524CalendarViewAlternative != "" && b524CalendarViewInput.B524CalendarViewAlternative == b524CalendarViewAlternative
}
