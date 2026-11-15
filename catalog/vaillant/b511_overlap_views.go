package vaillant

// The family and the variants of the B511 overlap views.
const (
	B511OverlapViewFamily           = "b511-family"
	B511OverlapViewOP00Reserved     = "op00-reserved"
	B511OverlapViewOP00Fan          = "op00-fan"
	B511OverlapViewOP07Base         = "op07-base"
	B511OverlapViewOP07ExtendedRoom = "op07-extended-room"
)

// B511OverlapViewSelection selects a B511 overlap view by family, variant, operation and direction.
type B511OverlapViewSelection struct {
	B511OverlapViewFamily    string `json:"family"`
	B511OverlapViewVariant   string `json:"variant"`
	B511OverlapViewOperation *byte  `json:"operation"`
	B511OverlapViewDirection string `json:"direction"`
}

// B511OverlapView is the interpreted overlap of a B511 logical image: the raw flag, fan, room setpoint
// in degrees Celsius and pressure in percent, each set only when the selected variant and the
// available extent provide it.
type B511OverlapView struct {
	B511OverlapViewSelection             B511OverlapViewSelection `json:"selection"`
	B511OverlapViewRaw                   []byte                   `json:"raw"`
	B511OverlapViewAvailableExtent       int                      `json:"available_extent"`
	B511OverlapViewRequiredExtent        int                      `json:"required_extent"`
	B511OverlapViewCompleteLogicalExtent bool                     `json:"complete_logical_extent"`
	B511OverlapViewSelected              bool                     `json:"selected"`
	B511OverlapViewRawFlag               *byte                    `json:"raw_flag"`
	B511OverlapViewFanRunning            *bool                    `json:"fan_running"`
	B511OverlapViewRawOverlay            *byte                    `json:"raw_overlay"`
	B511OverlapViewRoomSetpointCelsius   *float64                 `json:"room_setpoint_celsius"`
	B511OverlapViewRawPressure           *byte                    `json:"raw_pressure"`
	B511OverlapViewPressurePercent       *byte                    `json:"pressure_percent"`
	B511OverlapViewPressureDomain        string                   `json:"pressure_domain"`
	B511OverlapViewCandidateGrade        string                   `json:"candidate_grade"`
	B511OverlapViewPhysicalQualification string                   `json:"physical_qualification"`
}

// B511OverlapViewInterpret interprets the logical image for the selection and returns a detached view;
// an unsupported selection or a short image leaves the values unset.
func B511OverlapViewInterpret(b511OverlapViewSelection B511OverlapViewSelection, b511OverlapViewLogicalImage []byte) B511OverlapView {
	if b511OverlapViewSelection.B511OverlapViewOperation != nil {
		b511OverlapViewOp := *b511OverlapViewSelection.B511OverlapViewOperation
		b511OverlapViewSelection.B511OverlapViewOperation = &b511OverlapViewOp
	}
	b511OverlapViewView := B511OverlapView{
		B511OverlapViewSelection:             b511OverlapViewSelection,
		B511OverlapViewRaw:                   append([]byte(nil), b511OverlapViewLogicalImage...),
		B511OverlapViewAvailableExtent:       len(b511OverlapViewLogicalImage),
		B511OverlapViewPressureDomain:        "Unknown",
		B511OverlapViewCandidateGrade:        "Unknown",
		B511OverlapViewPhysicalQualification: "Unknown",
	}
	if b511OverlapViewSelection.B511OverlapViewFamily != B511OverlapViewFamily || b511OverlapViewSelection.B511OverlapViewDirection != "positive" {
		return b511OverlapViewView
	}
	if b511OverlapViewSelection.B511OverlapViewOperation == nil {
		return b511OverlapViewView
	}
	switch {
	case *b511OverlapViewSelection.B511OverlapViewOperation == 0 && (b511OverlapViewSelection.B511OverlapViewVariant == B511OverlapViewOP00Reserved || b511OverlapViewSelection.B511OverlapViewVariant == B511OverlapViewOP00Fan):
		b511OverlapViewView.B511OverlapViewRequiredExtent = 12
		if len(b511OverlapViewLogicalImage) > 8 {
			b511OverlapViewFlag := b511OverlapViewLogicalImage[8]
			b511OverlapViewView.B511OverlapViewRawFlag = &b511OverlapViewFlag
			if b511OverlapViewSelection.B511OverlapViewVariant == B511OverlapViewOP00Fan {
				b511OverlapViewFan := b511OverlapViewFlag&0x80 != 0
				b511OverlapViewView.B511OverlapViewFanRunning = &b511OverlapViewFan
			}
		}
	case *b511OverlapViewSelection.B511OverlapViewOperation == 7 && (b511OverlapViewSelection.B511OverlapViewVariant == B511OverlapViewOP07Base || b511OverlapViewSelection.B511OverlapViewVariant == B511OverlapViewOP07ExtendedRoom):
		b511OverlapViewView.B511OverlapViewRequiredExtent = 8
		if len(b511OverlapViewLogicalImage) > 7 {
			b511OverlapViewOverlay := b511OverlapViewLogicalImage[7]
			b511OverlapViewView.B511OverlapViewRawOverlay = &b511OverlapViewOverlay
		}
		if b511OverlapViewSelection.B511OverlapViewVariant == B511OverlapViewOP07ExtendedRoom {
			b511OverlapViewView.B511OverlapViewRequiredExtent = 9
			b511OverlapViewView.B511OverlapViewPressureDomain = "unavailable"
			if b511OverlapViewView.B511OverlapViewRawOverlay != nil {
				b511OverlapViewRoom := float64(*b511OverlapViewView.B511OverlapViewRawOverlay) / 2
				b511OverlapViewView.B511OverlapViewRoomSetpointCelsius = &b511OverlapViewRoom
			}
			if len(b511OverlapViewLogicalImage) > 8 {
				b511OverlapViewPressure := b511OverlapViewLogicalImage[8]
				b511OverlapViewView.B511OverlapViewRawPressure = &b511OverlapViewPressure
				b511OverlapViewView.B511OverlapViewPressureDomain = "out-of-domain"
				if b511OverlapViewPressure <= 100 {
					b511OverlapViewPercent := b511OverlapViewPressure
					b511OverlapViewView.B511OverlapViewPressurePercent = &b511OverlapViewPercent
					b511OverlapViewView.B511OverlapViewPressureDomain = "in-domain"
				}
			}
		}
	default:
		return b511OverlapViewView
	}
	b511OverlapViewView.B511OverlapViewSelected = true
	b511OverlapViewView.B511OverlapViewCandidateGrade = "Hypothesis"
	b511OverlapViewView.B511OverlapViewCompleteLogicalExtent = len(b511OverlapViewLogicalImage) >= b511OverlapViewView.B511OverlapViewRequiredExtent
	return b511OverlapViewView
}
