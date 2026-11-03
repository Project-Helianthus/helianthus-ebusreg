package ebus_standard_catalog

type Ab76State string

type Ab76Encoding string

const (
	Ab76RawOnly             Ab76State    = "raw-only"
	Ab76Incomplete          Ab76State    = "incomplete"
	Ab76Extra               Ab76State    = "extra"
	Ab76Invalid             Ab76State    = "invalid"
	Ab76Complete            Ab76State    = "complete-candidate"
	Ab76Absent              Ab76State    = "absent"
	Ab76Valid               Ab76State    = "valid"
	Ab76EncodingUnspecified Ab76Encoding = ""
	Ab76DecimalChunks       Ab76Encoding = "unsigned-decimal-chunks"
	Ab76UnsignedU8          Ab76Encoding = "unsigned-u8"
)

type Ab76Selection struct {
	Ab76Family    string
	Ab76Variant   string
	Ab76Function  *uint8
	Ab76Operation *uint8
	Ab76Direction string
	Ab76Selector  string
	Ab76Encoding  Ab76Encoding
}

type Ab76Component struct {
	Ab76Name     string
	Ab76State    Ab76State
	Ab76Raw      *uint8
	Ab76Value    *uint8
	Ab76Width    int
	Ab76Maximum  uint8
	Ab76Default  *uint8
	Ab76Encoding Ab76Encoding
}

type Ab76View struct {
	Ab76Selection           Ab76Selection
	Ab76State               Ab76State
	Ab76Raw                 []byte
	Ab76Remainder           []byte
	Ab76Components          []Ab76Component
	Ab76LocalCounter        *int
	Ab76LocalHours          *int
	Ab76LocalMinute         *int
	Ab76ImplementationGrade string
	Ab76MeaningGrade        string
	Ab76NativeProperties    map[string]string
}

func ab76CopyByte(ab76Value *uint8) *uint8 {
	if ab76Value == nil {
		return nil
	}
	ab76Copy := *ab76Value
	return &ab76Copy
}

func ab76CopyBody(ab76Body []byte) []byte {
	if ab76Body == nil {
		return nil
	}
	ab76Copy := make([]byte, len(ab76Body))
	copy(ab76Copy, ab76Body)
	return ab76Copy
}

func Ab76Interpret(ab76Selection Ab76Selection, ab76Body []byte) Ab76View {
	ab76Selection.Ab76Function = ab76CopyByte(ab76Selection.Ab76Function)
	ab76Selection.Ab76Operation = ab76CopyByte(ab76Selection.Ab76Operation)
	ab76View := Ab76View{
		Ab76Selection: ab76Selection, Ab76State: Ab76RawOnly,
		Ab76Raw:                 ab76CopyBody(ab76Body),
		Ab76ImplementationGrade: "Proven", Ab76MeaningGrade: "Hypothesis",
		Ab76NativeProperties: map[string]string{
			"physical-support": "Unknown", "presence": "Unknown", "freshness": "Unknown",
			"execution": "Unknown", "effects": "Unknown", "persistence": "Unknown",
			"delivery": "Unknown", "physical-length": "Unknown", "native-units": "Unknown",
			"transport-adaptation": "Unknown",
		},
	}
	if ab76Selection.Ab76Family != "standard-counter" || ab76Selection.Ab76Function == nil || *ab76Selection.Ab76Function != 3 || ab76Selection.Ab76Operation == nil {
		return ab76View
	}
	ab76Candidate := ab76Selection.Ab76Selector == "candidate" && ab76Selection.Ab76Direction == "request"
	ab76Standard := ab76Selection.Ab76Selector == "standard" && ab76Selection.Ab76Direction == "positive"
	if !ab76Candidate && !ab76Standard {
		return ab76View
	}
	var ab76Names []string
	var ab76Maxima []uint8
	switch {
	case ab76Selection.Ab76Variant == "03/04" && *ab76Selection.Ab76Operation == 4 && ab76Selection.Ab76Encoding == Ab76DecimalChunks:
		ab76Names = []string{"c0", "c1", "c2"}
		ab76Maxima = []uint8{99, 99, 99}
	case ab76Selection.Ab76Variant == "03/05" && *ab76Selection.Ab76Operation == 5 && ab76Candidate && ab76Selection.Ab76Encoding == Ab76UnsignedU8:
		ab76Names = []string{"minutes", "hundreds-of-hours", "ten-thousands-of-hours"}
		ab76Maxima = []uint8{59, 59, 59}
	case ab76Selection.Ab76Variant == "03/06" && *ab76Selection.Ab76Operation == 6 && ab76Selection.Ab76Encoding == Ab76UnsignedU8:
		ab76Names = []string{"minutes", "h0", "h1", "h2"}
		ab76Maxima = []uint8{59, 99, 99, 99}
	default:
		return ab76View
	}
	ab76View.Ab76State = Ab76Complete
	if len(ab76Body) < len(ab76Names) {
		ab76View.Ab76State = Ab76Incomplete
	}
	if len(ab76Body) > len(ab76Names) {
		ab76View.Ab76State = Ab76Extra
		ab76View.Ab76Remainder = append([]byte(nil), ab76Body[len(ab76Names):]...)
	}
	for ab76Index, ab76Name := range ab76Names {
		ab76Component := Ab76Component{Ab76Name: ab76Name, Ab76State: Ab76Absent, Ab76Width: 1, Ab76Maximum: ab76Maxima[ab76Index], Ab76Encoding: ab76Selection.Ab76Encoding}
		if ab76Selection.Ab76Variant == "03/04" {
			ab76Zero := uint8(0)
			ab76Component.Ab76Default = &ab76Zero
		}
		if ab76Index < len(ab76Body) {
			ab76Component.Ab76Raw = ab76CopyByte(&ab76Body[ab76Index])
			ab76Component.Ab76State = Ab76Invalid
			if ab76Body[ab76Index] <= ab76Maxima[ab76Index] {
				ab76Component.Ab76State = Ab76Valid
				ab76Component.Ab76Value = ab76CopyByte(&ab76Body[ab76Index])
			} else {
				ab76View.Ab76State = Ab76Invalid
			}
		}
		ab76View.Ab76Components = append(ab76View.Ab76Components, ab76Component)
	}
	if ab76View.Ab76State == Ab76Complete {
		switch ab76Selection.Ab76Variant {
		case "03/04":
			ab76Counter := int(ab76Body[0]) + 100*int(ab76Body[1]) + 10000*int(ab76Body[2])
			ab76View.Ab76LocalCounter = &ab76Counter
		case "03/06":
			ab76Minute := int(ab76Body[0])
			ab76Hours := int(ab76Body[1]) + 100*int(ab76Body[2]) + 10000*int(ab76Body[3])
			ab76View.Ab76LocalMinute, ab76View.Ab76LocalHours = &ab76Minute, &ab76Hours
		}
	}
	return ab76View
}
