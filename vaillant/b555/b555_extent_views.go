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
	ExtentViewRaw               []byte         `json:"raw"`
	ExtentViewDeclaredCount     int            `json:"declared_count"`
	ExtentViewAvailableCount    int            `json:"available_count"`
	Extent                      int            `json:"extent"`
	ExtentViewIncomplete        bool           `json:"incomplete"`
	ExtentViewIssues            []string       `json:"issues"`
	ExtentViewStatus            ExtentViewByte `json:"status"`
	ExtentViewStatusAccepted    bool           `json:"status_accepted"`
	ExtentViewStatusUnavailable bool           `json:"status_unavailable"`
	ExtentViewStatusGrade       string         `json:"status_grade"`
	ExtentViewDayCounts         []byte         `json:"day_counts"`
	ExtentViewWeekComplete      bool           `json:"week_complete"`
	ExtentViewWeekValid         bool           `json:"week_valid"`
	ExtentViewDayCountGrade     string         `json:"day_count_grade"`
	ExtentViewEighthCount       ExtentViewByte `json:"eighth_count"`
	ExtentViewEighthCountGrade  string         `json:"eighth_count_grade"`
}

type ExtentA5View struct {
	ExtentViewRaw               []byte         `json:"raw"`
	ExtentViewDeclaredCount     int            `json:"declared_count"`
	ExtentViewAvailableCount    int            `json:"available_count"`
	Extent                      int            `json:"extent"`
	ExtentViewIncomplete        bool           `json:"incomplete"`
	ExtentViewIssues            []string       `json:"issues"`
	ExtentViewStatus            ExtentViewByte `json:"status"`
	ExtentViewStatusAccepted    bool           `json:"status_accepted"`
	ExtentViewStartHour         ExtentViewByte `json:"start_hour"`
	ExtentViewStartMinute       ExtentViewByte `json:"start_minute"`
	ExtentViewEndHour           ExtentViewByte `json:"end_hour"`
	ExtentViewEndMinute         ExtentViewByte `json:"end_minute"`
	ExtentViewSetpointPresent   bool           `json:"setpoint_present"`
	ExtentViewSetpointRaw       uint16         `json:"setpoint_raw"`
	ExtentViewSetpointTenths    float64        `json:"setpoint_tenths"`
	ExtentViewSetpointGrade     string         `json:"setpoint_grade"`
	ExtentViewSetpointQualified string         `json:"setpoint_qualified"`
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
