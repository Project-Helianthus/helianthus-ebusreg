package vaillant

const (
	Ab11An11Family           = "ab11an11-b511-family"
	Ab11An11OP00Reserved     = "ab11an11-op00-reserved"
	Ab11An11OP00Fan          = "ab11an11-op00-fan"
	Ab11An11OP07Base         = "ab11an11-op07-base"
	Ab11An11OP07ExtendedRoom = "ab11an11-op07-extended-room"
)

type Ab11An11Selection struct {
	Ab11An11Family    string
	Ab11An11Variant   string
	Ab11An11Operation byte
	Ab11An11Direction string
}

type Ab11An11View struct {
	Ab11An11Selection             Ab11An11Selection
	Ab11An11Raw                   []byte
	Ab11An11AvailableExtent       int
	Ab11An11RequiredExtent        int
	Ab11An11CompleteLogicalExtent bool
	Ab11An11Selected              bool
	Ab11An11RawFlag               *byte
	Ab11An11FanRunning            *bool
	Ab11An11RawOverlay            *byte
	Ab11An11RoomSetpointCelsius   *float64
	Ab11An11RawPressure           *byte
	Ab11An11PressurePercent       *byte
	Ab11An11PressureDomain        string
	Ab11An11CandidateGrade        string
	Ab11An11PhysicalQualification string
}

func Ab11An11Interpret(ab11an11Selection Ab11An11Selection, ab11an11LogicalImage []byte) Ab11An11View {
	ab11an11View := Ab11An11View{
		Ab11An11Selection:             ab11an11Selection,
		Ab11An11Raw:                   append([]byte(nil), ab11an11LogicalImage...),
		Ab11An11AvailableExtent:       len(ab11an11LogicalImage),
		Ab11An11PressureDomain:        "Unknown",
		Ab11An11CandidateGrade:        "Unknown",
		Ab11An11PhysicalQualification: "Unknown",
	}
	if ab11an11Selection.Ab11An11Family != Ab11An11Family || ab11an11Selection.Ab11An11Direction != "positive" {
		return ab11an11View
	}
	switch {
	case ab11an11Selection.Ab11An11Operation == 0 && (ab11an11Selection.Ab11An11Variant == Ab11An11OP00Reserved || ab11an11Selection.Ab11An11Variant == Ab11An11OP00Fan):
		ab11an11View.Ab11An11RequiredExtent = 12
		if len(ab11an11LogicalImage) > 8 {
			ab11an11Flag := ab11an11LogicalImage[8]
			ab11an11View.Ab11An11RawFlag = &ab11an11Flag
			if ab11an11Selection.Ab11An11Variant == Ab11An11OP00Fan {
				ab11an11Fan := ab11an11Flag&0x80 != 0
				ab11an11View.Ab11An11FanRunning = &ab11an11Fan
			}
		}
	case ab11an11Selection.Ab11An11Operation == 7 && (ab11an11Selection.Ab11An11Variant == Ab11An11OP07Base || ab11an11Selection.Ab11An11Variant == Ab11An11OP07ExtendedRoom):
		ab11an11View.Ab11An11RequiredExtent = 8
		if len(ab11an11LogicalImage) > 7 {
			ab11an11Overlay := ab11an11LogicalImage[7]
			ab11an11View.Ab11An11RawOverlay = &ab11an11Overlay
		}
		if ab11an11Selection.Ab11An11Variant == Ab11An11OP07ExtendedRoom {
			ab11an11View.Ab11An11RequiredExtent = 9
			ab11an11View.Ab11An11PressureDomain = "unavailable"
			if ab11an11View.Ab11An11RawOverlay != nil {
				ab11an11Room := float64(*ab11an11View.Ab11An11RawOverlay) / 2
				ab11an11View.Ab11An11RoomSetpointCelsius = &ab11an11Room
			}
			if len(ab11an11LogicalImage) > 8 {
				ab11an11Pressure := ab11an11LogicalImage[8]
				ab11an11View.Ab11An11RawPressure = &ab11an11Pressure
				ab11an11View.Ab11An11PressureDomain = "out-of-domain"
				if ab11an11Pressure <= 100 {
					ab11an11Percent := ab11an11Pressure
					ab11an11View.Ab11An11PressurePercent = &ab11an11Percent
					ab11an11View.Ab11An11PressureDomain = "in-domain"
				}
			}
		}
	default:
		return ab11an11View
	}
	ab11an11View.Ab11An11Selected = true
	ab11an11View.Ab11An11CandidateGrade = "Hypothesis"
	ab11an11View.Ab11An11CompleteLogicalExtent = len(ab11an11LogicalImage) >= ab11an11View.Ab11An11RequiredExtent
	return ab11an11View
}
