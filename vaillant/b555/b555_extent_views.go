package b555

const (
	ExtentViewGradeProven     = "proven"
	ExtentViewGradeHypothesis = "hypothesis"
	ExtentViewGradeUnknown    = "unknown"

	ExtentViewIssueDeclaredNegative    = "declared_negative"
	ExtentViewIssueAvailableShort      = "available_short"
	ExtentViewIssueAvailableExceeds    = "available_exceeds_declared"
	ExtentViewIssueDeclaredBelowWeek   = "declared_below_week"
	ExtentViewIssueDayCountOverLimit   = "day_count_over_limit"
	ExtentViewIssueDeclaredBelowA5     = "declared_below_a5_setpoint"
	extentViewA4DayCountLimit          = 12
	extentViewA4DayCount               = 7
	extentViewA4WeekExtent             = 8
	extentViewA4EighthExtent           = 9
	extentViewA5SetpointExtent         = 7
	extentViewA4StatusAcceptedValue    = 0
	extentViewA4StatusUnavailableValue = 3
)

type ExtentViewByte struct {
	Present bool
	Value   byte
}

type ExtentA4View struct {
	ExtentViewRaw               []byte         `json:"Ab52Raw"`
	ExtentViewDeclaredCount     int            `json:"Ab52DeclaredCount"`
	ExtentViewAvailableCount    int            `json:"Ab52AvailableCount"`
	Extent                      int            `json:"Ab52Extent"`
	ExtentViewIncomplete        bool           `json:"Ab52Incomplete"`
	ExtentViewIssues            []string       `json:"Ab52Issues"`
	ExtentViewStatus            ExtentViewByte `json:"Ab52Status"`
	ExtentViewStatusAccepted    bool           `json:"Ab52StatusAccepted"`
	ExtentViewStatusUnavailable bool           `json:"Ab52StatusUnavailable"`
	ExtentViewStatusGrade       string         `json:"Ab52StatusGrade"`
	ExtentViewDayCounts         []byte         `json:"Ab52DayCounts"`
	ExtentViewWeekComplete      bool           `json:"Ab52WeekComplete"`
	ExtentViewWeekValid         bool           `json:"Ab52WeekValid"`
	ExtentViewDayCountGrade     string         `json:"Ab52DayCountGrade"`
	ExtentViewEighthCount       ExtentViewByte `json:"Ab52EighthCount"`
	ExtentViewEighthCountGrade  string         `json:"Ab52EighthCountGrade"`
}

type ExtentA5View struct {
	ExtentViewRaw               []byte         `json:"Ab52Raw"`
	ExtentViewDeclaredCount     int            `json:"Ab52DeclaredCount"`
	ExtentViewAvailableCount    int            `json:"Ab52AvailableCount"`
	Extent                      int            `json:"Ab52Extent"`
	ExtentViewIncomplete        bool           `json:"Ab52Incomplete"`
	ExtentViewIssues            []string       `json:"Ab52Issues"`
	ExtentViewStatus            ExtentViewByte `json:"Ab52Status"`
	ExtentViewStatusAccepted    bool           `json:"Ab52StatusAccepted"`
	ExtentViewStartHour         ExtentViewByte `json:"Ab52StartHour"`
	ExtentViewStartMinute       ExtentViewByte `json:"Ab52StartMinute"`
	ExtentViewEndHour           ExtentViewByte `json:"Ab52EndHour"`
	ExtentViewEndMinute         ExtentViewByte `json:"Ab52EndMinute"`
	ExtentViewSetpointPresent   bool           `json:"Ab52SetpointPresent"`
	ExtentViewSetpointRaw       uint16         `json:"Ab52SetpointRaw"`
	ExtentViewSetpointTenths    float64        `json:"Ab52SetpointTenths"`
	ExtentViewSetpointGrade     string         `json:"Ab52SetpointGrade"`
	ExtentViewSetpointQualified string         `json:"Ab52SetpointQualified"`
}

func extentViewExtent(declared int, available []byte) (int, []string) {
	issues := []string{}
	if declared < 0 {
		return 0, append(issues, ExtentViewIssueDeclaredNegative)
	}
	extent := declared
	if len(available) < extent {
		extent = len(available)
		issues = append(issues, ExtentViewIssueAvailableShort)
	}
	if len(available) > declared {
		issues = append(issues, ExtentViewIssueAvailableExceeds)
	}
	return extent, issues
}

func extentViewCell(raw []byte, index int) ExtentViewByte {
	if index < len(raw) {
		return ExtentViewByte{Present: true, Value: raw[index]}
	}
	return ExtentViewByte{}
}

func ExtentViewDecodeA4(declaredCount int, availableRaw []byte) ExtentA4View {
	extent, issues := extentViewExtent(declaredCount, availableRaw)
	raw := append([]byte{}, availableRaw[:extent]...)
	view := ExtentA4View{
		ExtentViewRaw:              raw,
		ExtentViewDeclaredCount:    declaredCount,
		ExtentViewAvailableCount:   len(availableRaw),
		Extent:                     extent,
		ExtentViewStatusGrade:      ExtentViewGradeHypothesis,
		ExtentViewDayCountGrade:    ExtentViewGradeHypothesis,
		ExtentViewEighthCountGrade: ExtentViewGradeHypothesis,
		ExtentViewDayCounts:        []byte{},
	}
	view.ExtentViewStatus = extentViewCell(raw, 0)
	if view.ExtentViewStatus.Present {
		view.ExtentViewStatusAccepted = view.ExtentViewStatus.Value == extentViewA4StatusAcceptedValue
		view.ExtentViewStatusUnavailable = view.ExtentViewStatus.Value == extentViewA4StatusUnavailableValue
	}
	for day := 1; day <= extentViewA4DayCount && day < len(raw); day++ {
		view.ExtentViewDayCounts = append(view.ExtentViewDayCounts, raw[day])
	}
	if declaredCount >= 0 && declaredCount < extentViewA4WeekExtent {
		issues = append(issues, ExtentViewIssueDeclaredBelowWeek)
	}
	view.ExtentViewWeekComplete = declaredCount >= extentViewA4WeekExtent && len(availableRaw) >= extentViewA4WeekExtent
	if view.ExtentViewWeekComplete {
		view.ExtentViewWeekValid = true
		for _, count := range view.ExtentViewDayCounts {
			if count > extentViewA4DayCountLimit {
				view.ExtentViewWeekValid = false
			}
		}
		if !view.ExtentViewWeekValid {
			issues = append(issues, ExtentViewIssueDayCountOverLimit)
		}
	}
	if declaredCount >= extentViewA4EighthExtent && len(availableRaw) >= extentViewA4EighthExtent {
		view.ExtentViewEighthCount = ExtentViewByte{Present: true, Value: availableRaw[8]}
	}
	view.ExtentViewIncomplete = !view.ExtentViewWeekComplete || len(availableRaw) < declaredCount
	view.ExtentViewIssues = issues
	return view
}

func ExtentViewDecodeA5(declaredCount int, availableRaw []byte) ExtentA5View {
	extent, issues := extentViewExtent(declaredCount, availableRaw)
	raw := append([]byte{}, availableRaw[:extent]...)
	view := ExtentA5View{
		ExtentViewRaw:               raw,
		ExtentViewDeclaredCount:     declaredCount,
		ExtentViewAvailableCount:    len(availableRaw),
		Extent:                      extent,
		ExtentViewSetpointGrade:     ExtentViewGradeHypothesis,
		ExtentViewSetpointQualified: ExtentViewGradeUnknown,
	}
	view.ExtentViewStatus = extentViewCell(raw, 0)
	view.ExtentViewStatusAccepted = view.ExtentViewStatus.Present && view.ExtentViewStatus.Value == 0
	view.ExtentViewStartHour = extentViewCell(raw, 1)
	view.ExtentViewStartMinute = extentViewCell(raw, 2)
	view.ExtentViewEndHour = extentViewCell(raw, 3)
	view.ExtentViewEndMinute = extentViewCell(raw, 4)
	if extent >= extentViewA5SetpointExtent {
		view.ExtentViewSetpointPresent = true
		view.ExtentViewSetpointRaw = uint16(raw[5]) | uint16(raw[6])<<8
		view.ExtentViewSetpointTenths = float64(view.ExtentViewSetpointRaw) / 10
	} else if declaredCount >= 0 && declaredCount < extentViewA5SetpointExtent {
		issues = append(issues, ExtentViewIssueDeclaredBelowA5)
	}
	view.ExtentViewIncomplete = !view.ExtentViewSetpointPresent || len(availableRaw) < declaredCount
	view.ExtentViewIssues = issues
	return view
}
