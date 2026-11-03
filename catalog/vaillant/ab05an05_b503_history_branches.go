package vaillant

type Ab05An05Context struct {
	Ab05An05Family    string
	Ab05An05Variant   string
	Ab05An05Operation uint8
	Ab05An05Direction string
}

type Ab05An05Cell struct {
	Ab05An05RawPresent        bool
	Ab05An05Raw               uint8
	Ab05An05NormalizedPresent bool
	Ab05An05Value             uint64
}

type Ab05An05Fields struct {
	Ab05An05Context          Ab05An05Context
	Ab05An05MarkerPresent    bool
	Ab05An05Marker           uint8
	Ab05An05RemainderPresent bool
	Ab05An05Remainder        [7]uint8
	Ab05An05ErrorPresent     bool
	Ab05An05Error            uint64
	Ab05An05Calendar         [5]Ab05An05Cell
}

type Ab05An05Component struct {
	Ab05An05Cell  Ab05An05Cell
	Ab05An05State string
}

type Ab05An05View struct {
	Ab05An05Fields        Ab05An05Fields
	Ab05An05State         string
	Ab05An05Components    [5]Ab05An05Component
	Ab05An05Complete      bool
	Ab05An05Contract      string
	Ab05An05Meaning       string
	Ab05An05Qualification string
}

func Ab05An05History(ab05an05Fields Ab05An05Fields) Ab05An05View {
	ab05an05View := Ab05An05View{
		Ab05An05Fields:        ab05an05Fields,
		Ab05An05State:         "unknown-branch",
		Ab05An05Contract:      "Proven",
		Ab05An05Meaning:       "Unknown",
		Ab05An05Qualification: "Unknown",
	}
	for ab05an05Index, ab05an05Cell := range ab05an05Fields.Ab05An05Calendar {
		ab05an05View.Ab05An05Components[ab05an05Index] = Ab05An05Component{Ab05An05Cell: ab05an05Cell, Ab05An05State: "unknown"}
	}
	if !ab05an05Fields.Ab05An05MarkerPresent || ab05an05Fields.Ab05An05Context.Ab05An05Family == "" || ab05an05Fields.Ab05An05Context.Ab05An05Variant == "" || ab05an05Fields.Ab05An05Context.Ab05An05Direction == "" || ab05an05Fields.Ab05An05Context.Ab05An05Operation != 1 {
		return ab05an05View
	}
	switch ab05an05Fields.Ab05An05Marker {
	case 0:
		if !ab05an05Fields.Ab05An05RemainderPresent {
			return ab05an05View
		}
		for _, ab05an05Byte := range ab05an05Fields.Ab05An05Remainder {
			if ab05an05Byte != 255 {
				return ab05an05View
			}
		}
		ab05an05View.Ab05An05State = "no-item"
	case 1, 2:
		ab05an05Complete := true
		ab05an05Unavailable := false
		for ab05an05Index, ab05an05Cell := range ab05an05Fields.Ab05An05Calendar {
			ab05an05State := "unknown"
			switch {
			case !ab05an05Cell.Ab05An05RawPresent:
				ab05an05State = "missing"
			case ab05an05Cell.Ab05An05Raw == 255:
				ab05an05State = "replacement"
			case !ab05an05Cell.Ab05An05NormalizedPresent:
				ab05an05State = "missing"
			case ab05an05Fields.Ab05An05Marker == 2:
				ab05an05State = "out-of-domain"
				if ab05an05Cell.Ab05An05Value >= [5]uint64{0, 0, 1, 1, 0}[ab05an05Index] && ab05an05Cell.Ab05An05Value <= [5]uint64{59, 23, 31, 12, 99}[ab05an05Index] {
					ab05an05State = "candidate"
				}
			}
			ab05an05View.Ab05An05Components[ab05an05Index].Ab05An05State = ab05an05State
			ab05an05Complete = ab05an05Complete && ab05an05State == "candidate"
			ab05an05Unavailable = ab05an05Unavailable || ab05an05State == "missing"
		}
		if ab05an05Fields.Ab05An05Marker == 1 {
			if !ab05an05Fields.Ab05An05ErrorPresent {
				return ab05an05View
			}
			ab05an05View.Ab05An05State = "error-without-complete-time"
		} else {
			if ab05an05Unavailable {
				return ab05an05View
			}
			ab05an05View.Ab05An05State = "candidate-components"
			ab05an05View.Ab05An05Complete = ab05an05Complete
		}
	default:
		return ab05an05View
	}
	ab05an05View.Ab05An05Meaning = "Hypothesis"
	return ab05an05View
}
