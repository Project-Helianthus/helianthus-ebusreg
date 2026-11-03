package b555

const (
	Ab52GradeProven     = "proven"
	Ab52GradeHypothesis = "hypothesis"
	Ab52GradeUnknown    = "unknown"

	Ab52IssueDeclaredNegative    = "declared_negative"
	Ab52IssueAvailableShort      = "available_short"
	Ab52IssueAvailableExceeds    = "available_exceeds_declared"
	Ab52IssueDeclaredBelowWeek   = "declared_below_week"
	Ab52IssueDayCountOverLimit   = "day_count_over_limit"
	Ab52IssueDeclaredBelowA5     = "declared_below_a5_setpoint"
	ab52A4DayCountLimit          = 12
	ab52A4DayCount               = 7
	ab52A4WeekExtent             = 8
	ab52A4EighthExtent           = 9
	ab52A5SetpointExtent         = 7
	ab52A4StatusAcceptedValue    = 0
	ab52A4StatusUnavailableValue = 3
)

type Ab52Byte struct {
	Present bool
	Value   byte
}

type Ab52A4View struct {
	Ab52Raw               []byte
	Ab52DeclaredCount     int
	Ab52AvailableCount    int
	Ab52Extent            int
	Ab52Incomplete        bool
	Ab52Issues            []string
	Ab52Status            Ab52Byte
	Ab52StatusAccepted    bool
	Ab52StatusUnavailable bool
	Ab52StatusGrade       string
	Ab52DayCounts         []byte
	Ab52WeekComplete      bool
	Ab52WeekValid         bool
	Ab52DayCountGrade     string
	Ab52EighthCount       Ab52Byte
	Ab52EighthCountGrade  string
}

type Ab52A5View struct {
	Ab52Raw               []byte
	Ab52DeclaredCount     int
	Ab52AvailableCount    int
	Ab52Extent            int
	Ab52Incomplete        bool
	Ab52Issues            []string
	Ab52Status            Ab52Byte
	Ab52StatusAccepted    bool
	Ab52StartHour         Ab52Byte
	Ab52StartMinute       Ab52Byte
	Ab52EndHour           Ab52Byte
	Ab52EndMinute         Ab52Byte
	Ab52SetpointPresent   bool
	Ab52SetpointRaw       uint16
	Ab52SetpointTenths    float64
	Ab52SetpointGrade     string
	Ab52SetpointQualified string
}

func ab52Extent(declared int, available []byte) (int, []string) {
	issues := []string{}
	if declared < 0 {
		return 0, append(issues, Ab52IssueDeclaredNegative)
	}
	extent := declared
	if len(available) < extent {
		extent = len(available)
		issues = append(issues, Ab52IssueAvailableShort)
	}
	if len(available) > declared {
		issues = append(issues, Ab52IssueAvailableExceeds)
	}
	return extent, issues
}

func ab52Cell(raw []byte, index int) Ab52Byte {
	if index < len(raw) {
		return Ab52Byte{Present: true, Value: raw[index]}
	}
	return Ab52Byte{}
}

func Ab52DecodeA4(declaredCount int, availableRaw []byte) Ab52A4View {
	extent, issues := ab52Extent(declaredCount, availableRaw)
	raw := append([]byte{}, availableRaw[:extent]...)
	view := Ab52A4View{
		Ab52Raw:              raw,
		Ab52DeclaredCount:    declaredCount,
		Ab52AvailableCount:   len(availableRaw),
		Ab52Extent:           extent,
		Ab52StatusGrade:      Ab52GradeHypothesis,
		Ab52DayCountGrade:    Ab52GradeHypothesis,
		Ab52EighthCountGrade: Ab52GradeHypothesis,
		Ab52DayCounts:        []byte{},
	}
	view.Ab52Status = ab52Cell(raw, 0)
	if view.Ab52Status.Present {
		view.Ab52StatusAccepted = view.Ab52Status.Value == ab52A4StatusAcceptedValue
		view.Ab52StatusUnavailable = view.Ab52Status.Value == ab52A4StatusUnavailableValue
	}
	for day := 1; day <= ab52A4DayCount && day < len(raw); day++ {
		view.Ab52DayCounts = append(view.Ab52DayCounts, raw[day])
	}
	if declaredCount >= 0 && declaredCount < ab52A4WeekExtent {
		issues = append(issues, Ab52IssueDeclaredBelowWeek)
	}
	view.Ab52WeekComplete = declaredCount >= ab52A4WeekExtent && len(availableRaw) >= ab52A4WeekExtent
	if view.Ab52WeekComplete {
		view.Ab52WeekValid = true
		for _, count := range view.Ab52DayCounts {
			if count > ab52A4DayCountLimit {
				view.Ab52WeekValid = false
			}
		}
		if !view.Ab52WeekValid {
			issues = append(issues, Ab52IssueDayCountOverLimit)
		}
	}
	if declaredCount >= ab52A4EighthExtent && len(availableRaw) >= ab52A4EighthExtent {
		view.Ab52EighthCount = Ab52Byte{Present: true, Value: availableRaw[8]}
	}
	view.Ab52Incomplete = !view.Ab52WeekComplete || len(availableRaw) < declaredCount
	view.Ab52Issues = issues
	return view
}

func Ab52DecodeA5(declaredCount int, availableRaw []byte) Ab52A5View {
	extent, issues := ab52Extent(declaredCount, availableRaw)
	raw := append([]byte{}, availableRaw[:extent]...)
	view := Ab52A5View{
		Ab52Raw:               raw,
		Ab52DeclaredCount:     declaredCount,
		Ab52AvailableCount:    len(availableRaw),
		Ab52Extent:            extent,
		Ab52SetpointGrade:     Ab52GradeHypothesis,
		Ab52SetpointQualified: Ab52GradeUnknown,
	}
	view.Ab52Status = ab52Cell(raw, 0)
	view.Ab52StatusAccepted = view.Ab52Status.Present && view.Ab52Status.Value == 0
	view.Ab52StartHour = ab52Cell(raw, 1)
	view.Ab52StartMinute = ab52Cell(raw, 2)
	view.Ab52EndHour = ab52Cell(raw, 3)
	view.Ab52EndMinute = ab52Cell(raw, 4)
	if extent >= ab52A5SetpointExtent {
		view.Ab52SetpointPresent = true
		view.Ab52SetpointRaw = uint16(raw[5]) | uint16(raw[6])<<8
		view.Ab52SetpointTenths = float64(view.Ab52SetpointRaw) / 10
	} else if declaredCount >= 0 && declaredCount < ab52A5SetpointExtent {
		issues = append(issues, Ab52IssueDeclaredBelowA5)
	}
	view.Ab52Incomplete = !view.Ab52SetpointPresent || len(availableRaw) < declaredCount
	view.Ab52Issues = issues
	return view
}
