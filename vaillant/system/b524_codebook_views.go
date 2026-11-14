package system

type B524CodebookViewInput struct {
	B524CodebookViewFamily        string `json:"family"`
	B524CodebookViewSoftware      string `json:"software"`
	B524CodebookViewHardware      string `json:"hardware"`
	B524CodebookViewAPIVersion    string `json:"api_version"`
	B524CodebookViewAPIRevision   string `json:"api_revision"`
	B524CodebookViewOP            byte   `json:"op"`
	B524CodebookViewGG            byte   `json:"gg"`
	B524CodebookViewII            *byte  `json:"ii"`
	B524CodebookViewRR            uint16 `json:"rr"`
	B524CodebookViewWidth         string `json:"width"`
	B524CodebookViewRequestedGG   byte   `json:"requested_gg"`
	B524CodebookViewRequestedRR   uint16 `json:"requested_rr"`
	B524CodebookViewReturnedGG    byte   `json:"returned_gg"`
	B524CodebookViewReturnedRR    uint16 `json:"returned_rr"`
	B524CodebookViewResponseState string `json:"response_state"`
	B524CodebookViewRaw           []byte `json:"raw"`
}

type B524CodebookViewResult struct {
	B524CodebookViewContext             B524CodebookViewInput `json:"context"`
	B524CodebookViewInteger             *int                  `json:"integer"`
	B524CodebookViewCandidateLabel      *string               `json:"candidate_label"`
	B524CodebookViewApplicable          bool                  `json:"applicable"`
	B524CodebookViewMalformedWidth      bool                  `json:"malformed_width"`
	B524CodebookViewGrade               string                `json:"grade"`
	B524CodebookViewNativeQualification string                `json:"native_qualification"`
}

func B524CodebookViewInterpret(b524CodebookViewInput B524CodebookViewInput) B524CodebookViewResult {
	b524CodebookViewInput.B524CodebookViewRaw = append([]byte(nil), b524CodebookViewInput.B524CodebookViewRaw...)
	if b524CodebookViewInput.B524CodebookViewII != nil {
		b524CodebookViewII := *b524CodebookViewInput.B524CodebookViewII
		b524CodebookViewInput.B524CodebookViewII = &b524CodebookViewII
	}
	b524CodebookViewResult := B524CodebookViewResult{B524CodebookViewContext: b524CodebookViewInput, B524CodebookViewGrade: "Unknown", B524CodebookViewNativeQualification: "Unknown"}
	b524CodebookViewSize := 0
	switch b524CodebookViewInput.B524CodebookViewWidth {
	case "u8", "signed8":
		b524CodebookViewSize = 1
	case "u16-le":
		b524CodebookViewSize = 2
	default:
		return b524CodebookViewResult
	}
	if len(b524CodebookViewInput.B524CodebookViewRaw) != b524CodebookViewSize {
		b524CodebookViewResult.B524CodebookViewMalformedWidth = true
		return b524CodebookViewResult
	}
	b524CodebookViewValue := int(b524CodebookViewInput.B524CodebookViewRaw[0])
	if b524CodebookViewSize == 2 {
		b524CodebookViewValue |= int(b524CodebookViewInput.B524CodebookViewRaw[1]) << 8
	}
	if b524CodebookViewInput.B524CodebookViewWidth == "signed8" {
		b524CodebookViewValue = int(int8(b524CodebookViewInput.B524CodebookViewRaw[0]))
	}
	b524CodebookViewResult.B524CodebookViewInteger = &b524CodebookViewValue
	if (b524CodebookViewInput.B524CodebookViewFamily != "controller-700" && b524CodebookViewInput.B524CodebookViewFamily != "controller-720") ||
		b524CodebookViewInput.B524CodebookViewOP != 2 || b524CodebookViewInput.B524CodebookViewII == nil ||
		b524CodebookViewInput.B524CodebookViewResponseState != "echo-payload" ||
		b524CodebookViewInput.B524CodebookViewRequestedGG != b524CodebookViewInput.B524CodebookViewGG || b524CodebookViewInput.B524CodebookViewRequestedRR != b524CodebookViewInput.B524CodebookViewRR ||
		b524CodebookViewInput.B524CodebookViewReturnedGG != b524CodebookViewInput.B524CodebookViewGG || b524CodebookViewInput.B524CodebookViewReturnedRR != b524CodebookViewInput.B524CodebookViewRR {
		return b524CodebookViewResult
	}
	var b524CodebookViewLabels map[int]string
	b524CodebookViewFamily700 := b524CodebookViewInput.B524CodebookViewFamily == "controller-700"
	switch {
	case (b524CodebookViewInput.B524CodebookViewGG == 1 && b524CodebookViewInput.B524CodebookViewRR == 3 && *b524CodebookViewInput.B524CodebookViewII == 0) ||
		(b524CodebookViewInput.B524CodebookViewGG == 3 && (b524CodebookViewInput.B524CodebookViewRR == 1 || b524CodebookViewInput.B524CodebookViewRR == 6)):
		if b524CodebookViewInput.B524CodebookViewWidth != "u16-le" {
			return b524CodebookViewResult
		}
		b524CodebookViewLabels = map[int]string{0: "off", 1: "time-controlled", 2: "manual"}
		if b524CodebookViewFamily700 {
			b524CodebookViewLabels = map[int]string{0: "off", 1: "automatic", 2: "day"}
			if b524CodebookViewInput.B524CodebookViewGG == 3 && b524CodebookViewInput.B524CodebookViewRR == 6 {
				b524CodebookViewLabels[3] = "setback"
			}
		}
	case b524CodebookViewInput.B524CodebookViewGG == 2 && b524CodebookViewInput.B524CodebookViewRR == 2:
		if b524CodebookViewInput.B524CodebookViewWidth != "u16-le" {
			return b524CodebookViewResult
		}
		b524CodebookViewLabels = map[int]string{0: "deactivated", 1: "heating", 2: "fixed-flow", 3: "DHW", 4: "return-flow"}
		if b524CodebookViewFamily700 {
			b524CodebookViewLabels[5] = "pool"
		}
	case b524CodebookViewInput.B524CodebookViewGG == 3 && b524CodebookViewInput.B524CodebookViewRR == 0x13:
		if b524CodebookViewInput.B524CodebookViewWidth != "u16-le" {
			return b524CodebookViewResult
		}
		b524CodebookViewLabels = map[int]string{0: "none", 1: "central", 2: "remote-circuit-one", 3: "remote-circuit-two", 4: "remote-circuit-three", 5: "remote-circuit-four"}
		if b524CodebookViewFamily700 {
			b524CodebookViewLabels[6] = "remote-circuit-five"
			b524CodebookViewLabels[7] = "remote-circuit-six"
			b524CodebookViewLabels[8] = "remote-circuit-seven"
			b524CodebookViewLabels[9] = "remote-circuit-eight"
		}
	case b524CodebookViewInput.B524CodebookViewGG == 3 && b524CodebookViewInput.B524CodebookViewRR == 0x1b:
		if b524CodebookViewInput.B524CodebookViewWidth != "u8" && b524CodebookViewInput.B524CodebookViewWidth != "u16-le" {
			return b524CodebookViewResult
		}
		b524CodebookViewLabels = map[int]string{0: "no-special-function", 1: "ventilation-boost", 4: "quick-veto", 7: "holiday", 10: "system-off"}
		if b524CodebookViewFamily700 {
			b524CodebookViewLabels = map[int]string{0: "no-quick-mode", 1: "ventilation-boost", 3: "party", 4: "quick-veto", 5: "one-day-holiday", 6: "one-day-bank-holiday", 7: "holiday", 8: "bank-holiday", 10: "system-off"}
		}
	case b524CodebookViewInput.B524CodebookViewGG == 3 && b524CodebookViewInput.B524CodebookViewRR == 0x1c:
		if b524CodebookViewInput.B524CodebookViewWidth != "signed8" {
			return b524CodebookViewResult
		}
		b524CodebookViewLabels = map[int]string{-1: "unbound", 0: "control-circuit-one", 1: "control-circuit-two", 2: "control-circuit-three"}
		if !b524CodebookViewFamily700 {
			b524CodebookViewLabels[3] = "control-circuit-four"
			b524CodebookViewLabels[4] = "control-circuit-five"
			b524CodebookViewLabels[5] = "control-circuit-six"
			b524CodebookViewLabels[6] = "control-circuit-seven"
			b524CodebookViewLabels[7] = "control-circuit-eight"
			b524CodebookViewLabels[8] = "control-circuit-nine"
		}
	default:
		return b524CodebookViewResult
	}
	b524CodebookViewResult.B524CodebookViewApplicable = true
	b524CodebookViewResult.B524CodebookViewGrade = "Hypothesis"
	if b524CodebookViewLabel, b524CodebookViewOK := b524CodebookViewLabels[b524CodebookViewValue]; b524CodebookViewOK {
		b524CodebookViewResult.B524CodebookViewCandidateLabel = &b524CodebookViewLabel
	}
	return b524CodebookViewResult
}
