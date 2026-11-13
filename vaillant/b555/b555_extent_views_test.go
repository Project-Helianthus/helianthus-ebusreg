package b555

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type extentViewExpCell struct {
	Present bool `json:"present"`
	Value   byte `json:"value"`
}

type extentViewA4Expect struct {
	RawHex            string   `json:"raw_hex"`
	StatusPresent     bool     `json:"status_present"`
	Status            byte     `json:"status"`
	StatusAccepted    bool     `json:"status_accepted"`
	StatusUnavailable bool     `json:"status_unavailable"`
	DayCounts         []int    `json:"day_counts"`
	WeekComplete      bool     `json:"week_complete"`
	WeekValid         bool     `json:"week_valid"`
	EighthPresent     bool     `json:"eighth_present"`
	Eighth            byte     `json:"eighth"`
	Incomplete        bool     `json:"incomplete"`
	Issues            []string `json:"issues"`
}

type extentViewA5Expect struct {
	RawHex          string            `json:"raw_hex"`
	Status          extentViewExpCell `json:"status"`
	StatusAccepted  bool              `json:"status_accepted"`
	StartHour       extentViewExpCell `json:"start_hour"`
	StartMinute     extentViewExpCell `json:"start_minute"`
	EndHour         extentViewExpCell `json:"end_hour"`
	EndMinute       extentViewExpCell `json:"end_minute"`
	SetpointPresent bool              `json:"setpoint_present"`
	SetpointRaw     uint16            `json:"setpoint_raw"`
	SetpointTenths  float64           `json:"setpoint_tenths"`
	Incomplete      bool              `json:"incomplete"`
	Issues          []string          `json:"issues"`
}

type extentViewFixture struct {
	A4 []struct {
		ID       string             `json:"ab52_id"`
		Declared int                `json:"ab52_declared"`
		RawHex   string             `json:"ab52_raw_hex"`
		Expect   extentViewA4Expect `json:"ab52_expect"`
	} `json:"ab52_a4_cases"`
	A5 []struct {
		ID       string             `json:"ab52_id"`
		Declared int                `json:"ab52_declared"`
		RawHex   string             `json:"ab52_raw_hex"`
		Expect   extentViewA5Expect `json:"ab52_expect"`
	} `json:"ab52_a5_cases"`
}

func extentViewLoad(extentViewT *testing.T) extentViewFixture {
	extentViewT.Helper()
	extentViewData, extentViewErr := os.ReadFile("testdata/ab52_b555_extent_views.json")
	if extentViewErr != nil {
		extentViewT.Fatal(extentViewErr)
	}
	var extentViewFix extentViewFixture
	if extentViewErr = json.Unmarshal(extentViewData, &extentViewFix); extentViewErr != nil {
		extentViewT.Fatal(extentViewErr)
	}
	return extentViewFix
}

func extentViewUnhex(extentViewT *testing.T, extentViewS string) []byte {
	extentViewT.Helper()
	extentViewB, extentViewErr := hex.DecodeString(extentViewS)
	if extentViewErr != nil {
		extentViewT.Fatal(extentViewErr)
	}
	return extentViewB
}

func extentViewSame(extentViewA, extentViewB []string) bool {
	return len(extentViewA) == len(extentViewB) && (len(extentViewA) == 0 || reflect.DeepEqual(extentViewA, extentViewB))
}

func extentViewCellOf(extentViewC ExtentViewByte) extentViewExpCell {
	return extentViewExpCell(extentViewC)
}

func TestAb52A4FixtureCases(extentViewT *testing.T) {
	extentViewFix := extentViewLoad(extentViewT)
	if len(extentViewFix.A4) < 10 {
		extentViewT.Fatalf("ab52 a4 fixture too small: %d", len(extentViewFix.A4))
	}
	for _, extentViewCase := range extentViewFix.A4 {
		extentViewT.Run(extentViewCase.ID, func(extentViewT *testing.T) {
			extentViewView := ExtentViewDecodeA4(extentViewCase.Declared, extentViewUnhex(extentViewT, extentViewCase.RawHex))
			extentViewWant := extentViewCase.Expect
			if hex.EncodeToString(extentViewView.ExtentViewRaw) != extentViewWant.RawHex {
				extentViewT.Errorf("raw %x want %s", extentViewView.ExtentViewRaw, extentViewWant.RawHex)
			}
			if extentViewView.ExtentViewStatus.Present != extentViewWant.StatusPresent || extentViewView.ExtentViewStatus.Value != extentViewWant.Status ||
				extentViewView.ExtentViewStatusAccepted != extentViewWant.StatusAccepted || extentViewView.ExtentViewStatusUnavailable != extentViewWant.StatusUnavailable {
				extentViewT.Errorf("status %+v accepted=%v unavailable=%v", extentViewView.ExtentViewStatus, extentViewView.ExtentViewStatusAccepted, extentViewView.ExtentViewStatusUnavailable)
			}
			if len(extentViewView.ExtentViewDayCounts) != len(extentViewWant.DayCounts) {
				extentViewT.Fatalf("day counts %v want %v", extentViewView.ExtentViewDayCounts, extentViewWant.DayCounts)
			}
			for extentViewI, extentViewV := range extentViewWant.DayCounts {
				if int(extentViewView.ExtentViewDayCounts[extentViewI]) != extentViewV {
					extentViewT.Errorf("day %d = %d want %d", extentViewI, extentViewView.ExtentViewDayCounts[extentViewI], extentViewV)
				}
			}
			if extentViewView.ExtentViewWeekComplete != extentViewWant.WeekComplete || extentViewView.ExtentViewWeekValid != extentViewWant.WeekValid {
				extentViewT.Errorf("week complete=%v valid=%v", extentViewView.ExtentViewWeekComplete, extentViewView.ExtentViewWeekValid)
			}
			if extentViewView.ExtentViewEighthCount.Present != extentViewWant.EighthPresent || extentViewView.ExtentViewEighthCount.Value != extentViewWant.Eighth {
				extentViewT.Errorf("eighth %+v", extentViewView.ExtentViewEighthCount)
			}
			if extentViewView.ExtentViewIncomplete != extentViewWant.Incomplete {
				extentViewT.Errorf("incomplete=%v", extentViewView.ExtentViewIncomplete)
			}
			if !extentViewSame(extentViewView.ExtentViewIssues, extentViewWant.Issues) {
				extentViewT.Errorf("issues %v want %v", extentViewView.ExtentViewIssues, extentViewWant.Issues)
			}
		})
	}
}

func TestAb52A5FixtureCases(extentViewT *testing.T) {
	extentViewFix := extentViewLoad(extentViewT)
	if len(extentViewFix.A5) < 8 {
		extentViewT.Fatalf("ab52 a5 fixture too small: %d", len(extentViewFix.A5))
	}
	for _, extentViewCase := range extentViewFix.A5 {
		extentViewT.Run(extentViewCase.ID, func(extentViewT *testing.T) {
			extentViewView := ExtentViewDecodeA5(extentViewCase.Declared, extentViewUnhex(extentViewT, extentViewCase.RawHex))
			extentViewWant := extentViewCase.Expect
			if hex.EncodeToString(extentViewView.ExtentViewRaw) != extentViewWant.RawHex {
				extentViewT.Errorf("raw %x want %s", extentViewView.ExtentViewRaw, extentViewWant.RawHex)
			}
			extentViewGot := []extentViewExpCell{extentViewCellOf(extentViewView.ExtentViewStatus), extentViewCellOf(extentViewView.ExtentViewStartHour), extentViewCellOf(extentViewView.ExtentViewStartMinute), extentViewCellOf(extentViewView.ExtentViewEndHour), extentViewCellOf(extentViewView.ExtentViewEndMinute)}
			extentViewExp := []extentViewExpCell{extentViewWant.Status, extentViewWant.StartHour, extentViewWant.StartMinute, extentViewWant.EndHour, extentViewWant.EndMinute}
			if !reflect.DeepEqual(extentViewGot, extentViewExp) {
				extentViewT.Errorf("cells %+v want %+v", extentViewGot, extentViewExp)
			}
			if extentViewView.ExtentViewStatusAccepted != extentViewWant.StatusAccepted {
				extentViewT.Errorf("status accepted=%v", extentViewView.ExtentViewStatusAccepted)
			}
			if extentViewView.ExtentViewSetpointPresent != extentViewWant.SetpointPresent || extentViewView.ExtentViewSetpointRaw != extentViewWant.SetpointRaw || extentViewView.ExtentViewSetpointTenths != extentViewWant.SetpointTenths {
				extentViewT.Errorf("setpoint present=%v raw=%d tenths=%v", extentViewView.ExtentViewSetpointPresent, extentViewView.ExtentViewSetpointRaw, extentViewView.ExtentViewSetpointTenths)
			}
			if extentViewView.ExtentViewIncomplete != extentViewWant.Incomplete {
				extentViewT.Errorf("incomplete=%v", extentViewView.ExtentViewIncomplete)
			}
			if !extentViewSame(extentViewView.ExtentViewIssues, extentViewWant.Issues) {
				extentViewT.Errorf("issues %v want %v", extentViewView.ExtentViewIssues, extentViewWant.Issues)
			}
		})
	}
}

func TestAb52A4ExtentMatrix(extentViewT *testing.T) {
	extentViewExecuted := 0
	for extentViewDeclared := 6; extentViewDeclared <= 8; extentViewDeclared++ {
		for extentViewLen := 0; extentViewLen <= 8; extentViewLen++ {
			extentViewIn := make([]byte, extentViewLen)
			for extentViewI := range extentViewIn {
				extentViewIn[extentViewI] = byte(0x40 + extentViewI)
			}
			extentViewView := ExtentViewDecodeA4(extentViewDeclared, extentViewIn)
			extentViewExt := extentViewLen
			if extentViewDeclared < extentViewExt {
				extentViewExt = extentViewDeclared
			}
			if !bytes.Equal(extentViewView.ExtentViewRaw, extentViewIn[:extentViewExt]) {
				extentViewT.Errorf("d=%d l=%d raw %x", extentViewDeclared, extentViewLen, extentViewView.ExtentViewRaw)
			}
			extentViewDays := extentViewExt - 1
			if extentViewDays < 0 {
				extentViewDays = 0
			}
			if extentViewDays > 7 {
				extentViewDays = 7
			}
			if len(extentViewView.ExtentViewDayCounts) != extentViewDays {
				extentViewT.Errorf("d=%d l=%d days %d want %d", extentViewDeclared, extentViewLen, len(extentViewView.ExtentViewDayCounts), extentViewDays)
			}
			extentViewWeek := extentViewDeclared >= 8 && extentViewLen >= 8
			if extentViewView.ExtentViewWeekComplete != extentViewWeek || extentViewView.ExtentViewIncomplete != (!extentViewWeek || extentViewLen < extentViewDeclared) {
				extentViewT.Errorf("d=%d l=%d week=%v incomplete=%v", extentViewDeclared, extentViewLen, extentViewView.ExtentViewWeekComplete, extentViewView.ExtentViewIncomplete)
			}
			if extentViewView.ExtentViewEighthCount.Present {
				extentViewT.Errorf("d=%d l=%d eighth exposed", extentViewDeclared, extentViewLen)
			}
			if extentViewView.ExtentViewStatus.Present != (extentViewExt >= 1) {
				extentViewT.Errorf("d=%d l=%d status presence", extentViewDeclared, extentViewLen)
			}
			extentViewExecuted++
		}
	}
	if extentViewExecuted != 27 {
		extentViewT.Fatalf("executed %d", extentViewExecuted)
	}
}

func TestAb52A4EighthRequiresBothNine(extentViewT *testing.T) {
	extentViewIn := []byte{0, 1, 2, 3, 4, 5, 6, 7, 0x77, 0x88}
	for extentViewDeclared := 7; extentViewDeclared <= 11; extentViewDeclared++ {
		for extentViewLen := 7; extentViewLen <= 10; extentViewLen++ {
			extentViewView := ExtentViewDecodeA4(extentViewDeclared, extentViewIn[:extentViewLen])
			extentViewWant := extentViewDeclared >= 9 && extentViewLen >= 9
			if extentViewView.ExtentViewEighthCount.Present != extentViewWant {
				extentViewT.Errorf("d=%d l=%d present=%v", extentViewDeclared, extentViewLen, extentViewView.ExtentViewEighthCount.Present)
			}
			if extentViewWant && extentViewView.ExtentViewEighthCount.Value != 0x77 {
				extentViewT.Errorf("eighth value %x", extentViewView.ExtentViewEighthCount.Value)
			}
			if len(extentViewView.ExtentViewDayCounts) > 7 {
				extentViewT.Errorf("eighth weekday exposed: %d", len(extentViewView.ExtentViewDayCounts))
			}
		}
	}
}

func TestAb52A4DayLimitBoundary(extentViewT *testing.T) {
	for extentViewDay := 0; extentViewDay < 7; extentViewDay++ {
		for _, extentViewPair := range []struct {
			value byte
			valid bool
		}{{12, true}, {13, false}, {255, false}} {
			extentViewIn := make([]byte, 8)
			extentViewIn[1+extentViewDay] = extentViewPair.value
			extentViewView := ExtentViewDecodeA4(8, extentViewIn)
			if extentViewView.ExtentViewWeekValid != extentViewPair.valid || !extentViewView.ExtentViewWeekComplete {
				extentViewT.Errorf("day %d value %d valid=%v", extentViewDay, extentViewPair.value, extentViewView.ExtentViewWeekValid)
			}
			if extentViewView.ExtentViewDayCounts[extentViewDay] != extentViewPair.value {
				extentViewT.Errorf("raw count lost for day %d", extentViewDay)
			}
		}
	}
	extentViewEighth := ExtentViewDecodeA4(9, []byte{0, 1, 1, 1, 1, 1, 1, 1, 200})
	if !extentViewEighth.ExtentViewWeekValid || extentViewEighth.ExtentViewEighthCount.Value != 200 {
		extentViewT.Errorf("eighth count must not take part in the day limit: %+v", extentViewEighth)
	}
}

func TestAb52A4StatusIsIndependentOfWeekView(extentViewT *testing.T) {
	for _, extentViewStatus := range []byte{0, 3, 1, 2, 4, 255} {
		extentViewView := ExtentViewDecodeA4(8, []byte{extentViewStatus, 1, 2, 3, 4, 5, 6, 7})
		if extentViewView.ExtentViewStatusAccepted != (extentViewStatus == 0) || extentViewView.ExtentViewStatusUnavailable != (extentViewStatus == 3) {
			extentViewT.Errorf("status %d classification", extentViewStatus)
		}
		if !extentViewView.ExtentViewWeekValid || len(extentViewView.ExtentViewDayCounts) != 7 {
			extentViewT.Errorf("status %d changed day view", extentViewStatus)
		}
	}
}

func TestAb52A5ExtentMatrix(extentViewT *testing.T) {
	extentViewExecuted := 0
	for extentViewDeclared := 6; extentViewDeclared <= 8; extentViewDeclared++ {
		for extentViewLen := 0; extentViewLen <= 8; extentViewLen++ {
			extentViewIn := make([]byte, extentViewLen)
			for extentViewI := range extentViewIn {
				extentViewIn[extentViewI] = byte(0x40 + extentViewI)
			}
			extentViewView := ExtentViewDecodeA5(extentViewDeclared, extentViewIn)
			extentViewExt := extentViewLen
			if extentViewDeclared < extentViewExt {
				extentViewExt = extentViewDeclared
			}
			if !bytes.Equal(extentViewView.ExtentViewRaw, extentViewIn[:extentViewExt]) {
				extentViewT.Errorf("d=%d l=%d raw %x", extentViewDeclared, extentViewLen, extentViewView.ExtentViewRaw)
			}
			extentViewExpCells := []ExtentViewByte{extentViewView.ExtentViewStatus, extentViewView.ExtentViewStartHour, extentViewView.ExtentViewStartMinute, extentViewView.ExtentViewEndHour, extentViewView.ExtentViewEndMinute}
			for extentViewI, extentViewC := range extentViewExpCells {
				if extentViewC.Present != (extentViewExt > extentViewI) || (extentViewC.Present && extentViewC.Value != byte(0x40+extentViewI)) {
					extentViewT.Errorf("d=%d l=%d cell %d %+v", extentViewDeclared, extentViewLen, extentViewI, extentViewC)
				}
			}
			extentViewSet := extentViewDeclared >= 7 && extentViewLen >= 7
			if extentViewView.ExtentViewSetpointPresent != extentViewSet {
				extentViewT.Errorf("d=%d l=%d setpoint present=%v", extentViewDeclared, extentViewLen, extentViewView.ExtentViewSetpointPresent)
			}
			if extentViewSet && extentViewView.ExtentViewSetpointRaw != 0x4645 {
				extentViewT.Errorf("setpoint raw %x", extentViewView.ExtentViewSetpointRaw)
			}
			if !extentViewSet && (extentViewView.ExtentViewSetpointRaw != 0 || extentViewView.ExtentViewSetpointTenths != 0) {
				extentViewT.Errorf("d=%d l=%d setpoint leaked", extentViewDeclared, extentViewLen)
			}
			if extentViewView.ExtentViewIncomplete != (!extentViewSet || extentViewLen < extentViewDeclared) {
				extentViewT.Errorf("d=%d l=%d incomplete=%v", extentViewDeclared, extentViewLen, extentViewView.ExtentViewIncomplete)
			}
			extentViewExecuted++
		}
	}
	if extentViewExecuted != 27 {
		extentViewT.Fatalf("executed %d", extentViewExecuted)
	}
}

func TestAb52A5LiteralPolicyAndSetpointQualification(extentViewT *testing.T) {
	extentViewView := ExtentViewDecodeA5(7, []byte{0x03, 0x17, 0x3B, 0x18, 0x3C, 0xFF, 0xFF})
	if extentViewView.ExtentViewStatus.Value != 3 || extentViewView.ExtentViewStatusAccepted {
		extentViewT.Errorf("status %+v accepted=%v", extentViewView.ExtentViewStatus, extentViewView.ExtentViewStatusAccepted)
	}
	if extentViewView.ExtentViewStartHour.Value != 23 || extentViewView.ExtentViewStartMinute.Value != 59 || extentViewView.ExtentViewEndHour.Value != 24 || extentViewView.ExtentViewEndMinute.Value != 60 {
		extentViewT.Errorf("times altered: %+v %+v %+v %+v", extentViewView.ExtentViewStartHour, extentViewView.ExtentViewStartMinute, extentViewView.ExtentViewEndHour, extentViewView.ExtentViewEndMinute)
	}
	if !extentViewView.ExtentViewSetpointPresent || extentViewView.ExtentViewSetpointRaw != 65535 || extentViewView.ExtentViewSetpointTenths != 6553.5 {
		extentViewT.Errorf("setpoint %+v", extentViewView)
	}
	if extentViewView.ExtentViewIncomplete || len(extentViewView.ExtentViewIssues) != 0 {
		extentViewT.Errorf("status 3 must not make decoding fail: %+v", extentViewView.ExtentViewIssues)
	}
	if extentViewView.ExtentViewSetpointGrade != ExtentViewGradeHypothesis || extentViewView.ExtentViewSetpointQualified != ExtentViewGradeUnknown {
		extentViewT.Errorf("grades %q %q", extentViewView.ExtentViewSetpointGrade, extentViewView.ExtentViewSetpointQualified)
	}
	extentViewPlain := ExtentViewDecodeA5(7, []byte{0, 0, 0, 0, 0, 0xFF, 0xFF})
	extentViewOther := ExtentViewDecodeA5(7, []byte{0, 0, 0, 0, 0, 0xFE, 0xFF})
	if extentViewPlain.ExtentViewSetpointRaw != 0xFFFF || extentViewOther.ExtentViewSetpointRaw != 0xFFFE || extentViewPlain.ExtentViewSetpointQualified != ExtentViewGradeUnknown {
		extentViewT.Errorf("sentinel inferred: %+v %+v", extentViewPlain, extentViewOther)
	}
}

func TestAb52OwnedMemory(extentViewT *testing.T) {
	extentViewSrc4 := []byte{0, 1, 2, 3, 4, 5, 6, 7, 9}
	extentViewV4 := ExtentViewDecodeA4(9, extentViewSrc4)
	extentViewSrc4[1], extentViewSrc4[8] = 0xEE, 0xEE
	if extentViewV4.ExtentViewRaw[1] != 1 || extentViewV4.ExtentViewDayCounts[0] != 1 || extentViewV4.ExtentViewEighthCount.Value != 9 {
		extentViewT.Fatalf("a4 aliases source: %+v", extentViewV4)
	}
	extentViewV4.ExtentViewRaw[2] = 0xAA
	extentViewV4.ExtentViewDayCounts[3] = 0xBB
	extentViewV4.ExtentViewIssues = append(extentViewV4.ExtentViewIssues, "x")
	extentViewAgain := ExtentViewDecodeA4(9, []byte{0, 1, 2, 3, 4, 5, 6, 7, 9})
	if extentViewAgain.ExtentViewRaw[2] != 2 || extentViewAgain.ExtentViewDayCounts[3] != 4 || len(extentViewAgain.ExtentViewIssues) != 0 {
		extentViewT.Fatalf("a4 shared state: %+v", extentViewAgain)
	}
	if &extentViewV4.ExtentViewRaw[0] == &extentViewV4.ExtentViewDayCounts[0] {
		extentViewT.Fatal("a4 raw and day counts share memory")
	}
	extentViewSrc5 := []byte{0, 1, 2, 3, 4, 5, 6}
	extentViewV5 := ExtentViewDecodeA5(7, extentViewSrc5)
	extentViewSrc5[5] = 0xEE
	if extentViewV5.ExtentViewRaw[5] != 5 || extentViewV5.ExtentViewSetpointRaw != 0x0605 {
		extentViewT.Fatalf("a5 aliases source: %+v", extentViewV5)
	}
	extentViewV5.ExtentViewRaw[6] = 0
	if extentViewFresh := ExtentViewDecodeA5(7, []byte{0, 1, 2, 3, 4, 5, 6}); extentViewFresh.ExtentViewRaw[6] != 6 {
		extentViewT.Fatal("a5 shared state")
	}
}

func TestAb52MalformedExtentsNeverSucceed(extentViewT *testing.T) {
	for _, extentViewDeclared := range []int{-1, -100, 0} {
		extentViewV4 := ExtentViewDecodeA4(extentViewDeclared, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8})
		extentViewV5 := ExtentViewDecodeA5(extentViewDeclared, nil)
		if len(extentViewV4.ExtentViewRaw) != 0 || !extentViewV4.ExtentViewIncomplete || extentViewV4.ExtentViewWeekComplete || extentViewV4.ExtentViewWeekValid || extentViewV4.ExtentViewEighthCount.Present || len(extentViewV4.ExtentViewIssues) == 0 {
			extentViewT.Errorf("a4 declared=%d: %+v", extentViewDeclared, extentViewV4)
		}
		if len(extentViewV5.ExtentViewRaw) != 0 || !extentViewV5.ExtentViewIncomplete || extentViewV5.ExtentViewSetpointPresent || extentViewV5.ExtentViewStatus.Present {
			extentViewT.Errorf("a5 declared=%d: %+v", extentViewDeclared, extentViewV5)
		}
	}
	extentViewHuge := ExtentViewDecodeA4(1<<30, []byte{0, 1, 2})
	if len(extentViewHuge.ExtentViewRaw) != 3 || !extentViewHuge.ExtentViewIncomplete {
		extentViewT.Errorf("huge declared: %+v", extentViewHuge)
	}
}

func TestRenamedFieldJSONKeys(t *testing.T) {
	cases := []struct {
		value any
		keys  map[string]string
	}{
		{NumericNarrowingContext{}, map[string]string{"NumericNarrowingTimerKind": "timer_kind", "NumericNarrowingFamily": "family", "NumericNarrowingSoftwareVariant": "software_variant"}},
		{NumericNarrowingKnowledge{}, map[string]string{"NumericNarrowingConversion": "conversion", "NumericNarrowingLocalCarrier": "local_carrier", "NumericNarrowingUnits": "units", "NumericNarrowingTimerKindEffects": "timer_kind_effects", "NumericNarrowingPersistence": "persistence", "NumericNarrowingAvailability": "availability", "NumericNarrowingCountMaterialization": "count_materialization", "NumericNarrowingWriteAdmission": "write_admission", "NumericNarrowingPhysicalQualification": "physical_qualification", "NumericNarrowingPresence": "presence", "NumericNarrowingFreshness": "freshness", "NumericNarrowingExecution": "execution"}},
		{NumericNarrowingResult{}, map[string]string{"NumericNarrowingContext": "context", "NumericNarrowingPath": "path", "NumericNarrowingOriginalInteger64": "original_integer64", "NumericNarrowingOriginalFloat64": "original_float64", "NumericNarrowingRoundedFloat64": "rounded_float64", "NumericNarrowingScaledFloat64": "scaled_float64", "NumericNarrowingRawClassification": "raw_classification", "NumericNarrowingRoundedClassification": "rounded_classification", "NumericNarrowingIntermediateStatus": "intermediate_status", "NumericNarrowingUnsigned32": "unsigned32", "NumericNarrowingLow16": "low16", "NumericNarrowingStatus": "status", "NumericNarrowingKnowledge": "knowledge"}},
		{ExtentA4View{}, map[string]string{"ExtentViewRaw": "Ab52Raw", "ExtentViewDeclaredCount": "Ab52DeclaredCount", "ExtentViewAvailableCount": "Ab52AvailableCount", "Extent": "Ab52Extent", "ExtentViewIncomplete": "Ab52Incomplete", "ExtentViewIssues": "Ab52Issues", "ExtentViewStatus": "Ab52Status", "ExtentViewStatusAccepted": "Ab52StatusAccepted", "ExtentViewStatusUnavailable": "Ab52StatusUnavailable", "ExtentViewStatusGrade": "Ab52StatusGrade", "ExtentViewDayCounts": "Ab52DayCounts", "ExtentViewWeekComplete": "Ab52WeekComplete", "ExtentViewWeekValid": "Ab52WeekValid", "ExtentViewDayCountGrade": "Ab52DayCountGrade", "ExtentViewEighthCount": "Ab52EighthCount", "ExtentViewEighthCountGrade": "Ab52EighthCountGrade"}},
		{ExtentA5View{}, map[string]string{"ExtentViewRaw": "Ab52Raw", "ExtentViewDeclaredCount": "Ab52DeclaredCount", "ExtentViewAvailableCount": "Ab52AvailableCount", "Extent": "Ab52Extent", "ExtentViewIncomplete": "Ab52Incomplete", "ExtentViewIssues": "Ab52Issues", "ExtentViewStatus": "Ab52Status", "ExtentViewStatusAccepted": "Ab52StatusAccepted", "ExtentViewStartHour": "Ab52StartHour", "ExtentViewStartMinute": "Ab52StartMinute", "ExtentViewEndHour": "Ab52EndHour", "ExtentViewEndMinute": "Ab52EndMinute", "ExtentViewSetpointPresent": "Ab52SetpointPresent", "ExtentViewSetpointRaw": "Ab52SetpointRaw", "ExtentViewSetpointTenths": "Ab52SetpointTenths", "ExtentViewSetpointGrade": "Ab52SetpointGrade", "ExtentViewSetpointQualified": "Ab52SetpointQualified"}},
	}
	for _, tc := range cases {
		data, err := json.Marshal(tc.value)
		if err != nil {
			t.Fatal(err)
		}
		var object map[string]json.RawMessage
		if err := json.Unmarshal(data, &object); err != nil {
			t.Fatal(err)
		}
		typ := reflect.TypeOf(tc.value)
		for field, key := range tc.keys {
			if _, ok := object[key]; !ok {
				t.Errorf("%s.%s missing JSON key %s", typ, field, key)
			}
			if _, ok := object[field]; ok {
				t.Errorf("%s exposes renamed key %s", typ, field)
			}
			if member, ok := typ.FieldByName(field); !ok || member.Tag.Get("json") != key {
				t.Errorf("%s.%s does not decode JSON key %s", typ, field, key)
			}
		}
	}
}
