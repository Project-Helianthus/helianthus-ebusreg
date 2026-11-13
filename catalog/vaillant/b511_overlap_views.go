package vaillant

const (
	B511OverlapViewFamily           = "ab11an11-b511-family"
	B511OverlapViewOP00Reserved     = "ab11an11-op00-reserved"
	B511OverlapViewOP00Fan          = "ab11an11-op00-fan"
	B511OverlapViewOP07Base         = "ab11an11-op07-base"
	B511OverlapViewOP07ExtendedRoom = "ab11an11-op07-extended-room"
)

type B511OverlapViewSelection struct {
	B511OverlapViewFamily    string `json:"Ab11An11Family"`
	B511OverlapViewVariant   string `json:"Ab11An11Variant"`
	B511OverlapViewOperation *byte  `json:"Ab11An11Operation"`
	B511OverlapViewDirection string `json:"Ab11An11Direction"`
}

type B511OverlapView struct {
	B511OverlapViewSelection             B511OverlapViewSelection `json:"Ab11An11Selection"`
	B511OverlapViewRaw                   []byte                   `json:"Ab11An11Raw"`
	B511OverlapViewAvailableExtent       int                      `json:"Ab11An11AvailableExtent"`
	B511OverlapViewRequiredExtent        int                      `json:"Ab11An11RequiredExtent"`
	B511OverlapViewCompleteLogicalExtent bool                     `json:"Ab11An11CompleteLogicalExtent"`
	B511OverlapViewSelected              bool                     `json:"Ab11An11Selected"`
	B511OverlapViewRawFlag               *byte                    `json:"Ab11An11RawFlag"`
	B511OverlapViewFanRunning            *bool                    `json:"Ab11An11FanRunning"`
	B511OverlapViewRawOverlay            *byte                    `json:"Ab11An11RawOverlay"`
	B511OverlapViewRoomSetpointCelsius   *float64                 `json:"Ab11An11RoomSetpointCelsius"`
	B511OverlapViewRawPressure           *byte                    `json:"Ab11An11RawPressure"`
	B511OverlapViewPressurePercent       *byte                    `json:"Ab11An11PressurePercent"`
	B511OverlapViewPressureDomain        string                   `json:"Ab11An11PressureDomain"`
	B511OverlapViewCandidateGrade        string                   `json:"Ab11An11CandidateGrade"`
	B511OverlapViewPhysicalQualification string                   `json:"Ab11An11PhysicalQualification"`
}

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
