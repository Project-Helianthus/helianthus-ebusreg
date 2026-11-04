package b555

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type ab52ExpCell struct {
	Present bool `json:"present"`
	Value   byte `json:"value"`
}

type ab52A4Expect struct {
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

type ab52A5Expect struct {
	RawHex          string      `json:"raw_hex"`
	Status          ab52ExpCell `json:"status"`
	StatusAccepted  bool        `json:"status_accepted"`
	StartHour       ab52ExpCell `json:"start_hour"`
	StartMinute     ab52ExpCell `json:"start_minute"`
	EndHour         ab52ExpCell `json:"end_hour"`
	EndMinute       ab52ExpCell `json:"end_minute"`
	SetpointPresent bool        `json:"setpoint_present"`
	SetpointRaw     uint16      `json:"setpoint_raw"`
	SetpointTenths  float64     `json:"setpoint_tenths"`
	Incomplete      bool        `json:"incomplete"`
	Issues          []string    `json:"issues"`
}

type ab52Fixture struct {
	A4 []struct {
		ID       string       `json:"ab52_id"`
		Declared int          `json:"ab52_declared"`
		RawHex   string       `json:"ab52_raw_hex"`
		Expect   ab52A4Expect `json:"ab52_expect"`
	} `json:"ab52_a4_cases"`
	A5 []struct {
		ID       string       `json:"ab52_id"`
		Declared int          `json:"ab52_declared"`
		RawHex   string       `json:"ab52_raw_hex"`
		Expect   ab52A5Expect `json:"ab52_expect"`
	} `json:"ab52_a5_cases"`
}

func ab52Load(ab52T *testing.T) ab52Fixture {
	ab52T.Helper()
	ab52Data, ab52Err := os.ReadFile("testdata/ab52_b555_extent_views.json")
	if ab52Err != nil {
		ab52T.Fatal(ab52Err)
	}
	var ab52Fix ab52Fixture
	if ab52Err = json.Unmarshal(ab52Data, &ab52Fix); ab52Err != nil {
		ab52T.Fatal(ab52Err)
	}
	return ab52Fix
}

func ab52Unhex(ab52T *testing.T, ab52S string) []byte {
	ab52T.Helper()
	ab52B, ab52Err := hex.DecodeString(ab52S)
	if ab52Err != nil {
		ab52T.Fatal(ab52Err)
	}
	return ab52B
}

func ab52Same(ab52A, ab52B []string) bool {
	return len(ab52A) == len(ab52B) && (len(ab52A) == 0 || reflect.DeepEqual(ab52A, ab52B))
}

func ab52CellOf(ab52C Ab52Byte) ab52ExpCell {
	return ab52ExpCell(ab52C)
}

func TestAb52A4FixtureCases(ab52T *testing.T) {
	ab52Fix := ab52Load(ab52T)
	if len(ab52Fix.A4) < 10 {
		ab52T.Fatalf("ab52 a4 fixture too small: %d", len(ab52Fix.A4))
	}
	for _, ab52Case := range ab52Fix.A4 {
		ab52T.Run(ab52Case.ID, func(ab52T *testing.T) {
			ab52View := Ab52DecodeA4(ab52Case.Declared, ab52Unhex(ab52T, ab52Case.RawHex))
			ab52Want := ab52Case.Expect
			if hex.EncodeToString(ab52View.Ab52Raw) != ab52Want.RawHex {
				ab52T.Errorf("raw %x want %s", ab52View.Ab52Raw, ab52Want.RawHex)
			}
			if ab52View.Ab52Status.Present != ab52Want.StatusPresent || ab52View.Ab52Status.Value != ab52Want.Status ||
				ab52View.Ab52StatusAccepted != ab52Want.StatusAccepted || ab52View.Ab52StatusUnavailable != ab52Want.StatusUnavailable {
				ab52T.Errorf("status %+v accepted=%v unavailable=%v", ab52View.Ab52Status, ab52View.Ab52StatusAccepted, ab52View.Ab52StatusUnavailable)
			}
			if len(ab52View.Ab52DayCounts) != len(ab52Want.DayCounts) {
				ab52T.Fatalf("day counts %v want %v", ab52View.Ab52DayCounts, ab52Want.DayCounts)
			}
			for ab52I, ab52V := range ab52Want.DayCounts {
				if int(ab52View.Ab52DayCounts[ab52I]) != ab52V {
					ab52T.Errorf("day %d = %d want %d", ab52I, ab52View.Ab52DayCounts[ab52I], ab52V)
				}
			}
			if ab52View.Ab52WeekComplete != ab52Want.WeekComplete || ab52View.Ab52WeekValid != ab52Want.WeekValid {
				ab52T.Errorf("week complete=%v valid=%v", ab52View.Ab52WeekComplete, ab52View.Ab52WeekValid)
			}
			if ab52View.Ab52EighthCount.Present != ab52Want.EighthPresent || ab52View.Ab52EighthCount.Value != ab52Want.Eighth {
				ab52T.Errorf("eighth %+v", ab52View.Ab52EighthCount)
			}
			if ab52View.Ab52Incomplete != ab52Want.Incomplete {
				ab52T.Errorf("incomplete=%v", ab52View.Ab52Incomplete)
			}
			if !ab52Same(ab52View.Ab52Issues, ab52Want.Issues) {
				ab52T.Errorf("issues %v want %v", ab52View.Ab52Issues, ab52Want.Issues)
			}
		})
	}
}

func TestAb52A5FixtureCases(ab52T *testing.T) {
	ab52Fix := ab52Load(ab52T)
	if len(ab52Fix.A5) < 8 {
		ab52T.Fatalf("ab52 a5 fixture too small: %d", len(ab52Fix.A5))
	}
	for _, ab52Case := range ab52Fix.A5 {
		ab52T.Run(ab52Case.ID, func(ab52T *testing.T) {
			ab52View := Ab52DecodeA5(ab52Case.Declared, ab52Unhex(ab52T, ab52Case.RawHex))
			ab52Want := ab52Case.Expect
			if hex.EncodeToString(ab52View.Ab52Raw) != ab52Want.RawHex {
				ab52T.Errorf("raw %x want %s", ab52View.Ab52Raw, ab52Want.RawHex)
			}
			ab52Got := []ab52ExpCell{ab52CellOf(ab52View.Ab52Status), ab52CellOf(ab52View.Ab52StartHour), ab52CellOf(ab52View.Ab52StartMinute), ab52CellOf(ab52View.Ab52EndHour), ab52CellOf(ab52View.Ab52EndMinute)}
			ab52Exp := []ab52ExpCell{ab52Want.Status, ab52Want.StartHour, ab52Want.StartMinute, ab52Want.EndHour, ab52Want.EndMinute}
			if !reflect.DeepEqual(ab52Got, ab52Exp) {
				ab52T.Errorf("cells %+v want %+v", ab52Got, ab52Exp)
			}
			if ab52View.Ab52StatusAccepted != ab52Want.StatusAccepted {
				ab52T.Errorf("status accepted=%v", ab52View.Ab52StatusAccepted)
			}
			if ab52View.Ab52SetpointPresent != ab52Want.SetpointPresent || ab52View.Ab52SetpointRaw != ab52Want.SetpointRaw || ab52View.Ab52SetpointTenths != ab52Want.SetpointTenths {
				ab52T.Errorf("setpoint present=%v raw=%d tenths=%v", ab52View.Ab52SetpointPresent, ab52View.Ab52SetpointRaw, ab52View.Ab52SetpointTenths)
			}
			if ab52View.Ab52Incomplete != ab52Want.Incomplete {
				ab52T.Errorf("incomplete=%v", ab52View.Ab52Incomplete)
			}
			if !ab52Same(ab52View.Ab52Issues, ab52Want.Issues) {
				ab52T.Errorf("issues %v want %v", ab52View.Ab52Issues, ab52Want.Issues)
			}
		})
	}
}

func TestAb52A4ExtentMatrix(ab52T *testing.T) {
	ab52Executed := 0
	for ab52Declared := 6; ab52Declared <= 8; ab52Declared++ {
		for ab52Len := 0; ab52Len <= 8; ab52Len++ {
			ab52In := make([]byte, ab52Len)
			for ab52I := range ab52In {
				ab52In[ab52I] = byte(0x40 + ab52I)
			}
			ab52View := Ab52DecodeA4(ab52Declared, ab52In)
			ab52Ext := ab52Len
			if ab52Declared < ab52Ext {
				ab52Ext = ab52Declared
			}
			if !bytes.Equal(ab52View.Ab52Raw, ab52In[:ab52Ext]) {
				ab52T.Errorf("d=%d l=%d raw %x", ab52Declared, ab52Len, ab52View.Ab52Raw)
			}
			ab52Days := ab52Ext - 1
			if ab52Days < 0 {
				ab52Days = 0
			}
			if ab52Days > 7 {
				ab52Days = 7
			}
			if len(ab52View.Ab52DayCounts) != ab52Days {
				ab52T.Errorf("d=%d l=%d days %d want %d", ab52Declared, ab52Len, len(ab52View.Ab52DayCounts), ab52Days)
			}
			ab52Week := ab52Declared >= 8 && ab52Len >= 8
			if ab52View.Ab52WeekComplete != ab52Week || ab52View.Ab52Incomplete != (!ab52Week || ab52Len < ab52Declared) {
				ab52T.Errorf("d=%d l=%d week=%v incomplete=%v", ab52Declared, ab52Len, ab52View.Ab52WeekComplete, ab52View.Ab52Incomplete)
			}
			if ab52View.Ab52EighthCount.Present {
				ab52T.Errorf("d=%d l=%d eighth exposed", ab52Declared, ab52Len)
			}
			if ab52View.Ab52Status.Present != (ab52Ext >= 1) {
				ab52T.Errorf("d=%d l=%d status presence", ab52Declared, ab52Len)
			}
			ab52Executed++
		}
	}
	if ab52Executed != 27 {
		ab52T.Fatalf("executed %d", ab52Executed)
	}
}

func TestAb52A4EighthRequiresBothNine(ab52T *testing.T) {
	ab52In := []byte{0, 1, 2, 3, 4, 5, 6, 7, 0x77, 0x88}
	for ab52Declared := 7; ab52Declared <= 11; ab52Declared++ {
		for ab52Len := 7; ab52Len <= 10; ab52Len++ {
			ab52View := Ab52DecodeA4(ab52Declared, ab52In[:ab52Len])
			ab52Want := ab52Declared >= 9 && ab52Len >= 9
			if ab52View.Ab52EighthCount.Present != ab52Want {
				ab52T.Errorf("d=%d l=%d present=%v", ab52Declared, ab52Len, ab52View.Ab52EighthCount.Present)
			}
			if ab52Want && ab52View.Ab52EighthCount.Value != 0x77 {
				ab52T.Errorf("eighth value %x", ab52View.Ab52EighthCount.Value)
			}
			if len(ab52View.Ab52DayCounts) > 7 {
				ab52T.Errorf("eighth weekday exposed: %d", len(ab52View.Ab52DayCounts))
			}
		}
	}
}

func TestAb52A4DayLimitBoundary(ab52T *testing.T) {
	for ab52Day := 0; ab52Day < 7; ab52Day++ {
		for _, ab52Pair := range []struct {
			value byte
			valid bool
		}{{12, true}, {13, false}, {255, false}} {
			ab52In := make([]byte, 8)
			ab52In[1+ab52Day] = ab52Pair.value
			ab52View := Ab52DecodeA4(8, ab52In)
			if ab52View.Ab52WeekValid != ab52Pair.valid || !ab52View.Ab52WeekComplete {
				ab52T.Errorf("day %d value %d valid=%v", ab52Day, ab52Pair.value, ab52View.Ab52WeekValid)
			}
			if ab52View.Ab52DayCounts[ab52Day] != ab52Pair.value {
				ab52T.Errorf("raw count lost for day %d", ab52Day)
			}
		}
	}
	ab52Eighth := Ab52DecodeA4(9, []byte{0, 1, 1, 1, 1, 1, 1, 1, 200})
	if !ab52Eighth.Ab52WeekValid || ab52Eighth.Ab52EighthCount.Value != 200 {
		ab52T.Errorf("eighth count must not take part in the day limit: %+v", ab52Eighth)
	}
}

func TestAb52A4StatusIsIndependentOfWeekView(ab52T *testing.T) {
	for _, ab52Status := range []byte{0, 3, 1, 2, 4, 255} {
		ab52View := Ab52DecodeA4(8, []byte{ab52Status, 1, 2, 3, 4, 5, 6, 7})
		if ab52View.Ab52StatusAccepted != (ab52Status == 0) || ab52View.Ab52StatusUnavailable != (ab52Status == 3) {
			ab52T.Errorf("status %d classification", ab52Status)
		}
		if !ab52View.Ab52WeekValid || len(ab52View.Ab52DayCounts) != 7 {
			ab52T.Errorf("status %d changed day view", ab52Status)
		}
	}
}

func TestAb52A5ExtentMatrix(ab52T *testing.T) {
	ab52Executed := 0
	for ab52Declared := 6; ab52Declared <= 8; ab52Declared++ {
		for ab52Len := 0; ab52Len <= 8; ab52Len++ {
			ab52In := make([]byte, ab52Len)
			for ab52I := range ab52In {
				ab52In[ab52I] = byte(0x40 + ab52I)
			}
			ab52View := Ab52DecodeA5(ab52Declared, ab52In)
			ab52Ext := ab52Len
			if ab52Declared < ab52Ext {
				ab52Ext = ab52Declared
			}
			if !bytes.Equal(ab52View.Ab52Raw, ab52In[:ab52Ext]) {
				ab52T.Errorf("d=%d l=%d raw %x", ab52Declared, ab52Len, ab52View.Ab52Raw)
			}
			ab52ExpCells := []Ab52Byte{ab52View.Ab52Status, ab52View.Ab52StartHour, ab52View.Ab52StartMinute, ab52View.Ab52EndHour, ab52View.Ab52EndMinute}
			for ab52I, ab52C := range ab52ExpCells {
				if ab52C.Present != (ab52Ext > ab52I) || (ab52C.Present && ab52C.Value != byte(0x40+ab52I)) {
					ab52T.Errorf("d=%d l=%d cell %d %+v", ab52Declared, ab52Len, ab52I, ab52C)
				}
			}
			ab52Set := ab52Declared >= 7 && ab52Len >= 7
			if ab52View.Ab52SetpointPresent != ab52Set {
				ab52T.Errorf("d=%d l=%d setpoint present=%v", ab52Declared, ab52Len, ab52View.Ab52SetpointPresent)
			}
			if ab52Set && ab52View.Ab52SetpointRaw != 0x4645 {
				ab52T.Errorf("setpoint raw %x", ab52View.Ab52SetpointRaw)
			}
			if !ab52Set && (ab52View.Ab52SetpointRaw != 0 || ab52View.Ab52SetpointTenths != 0) {
				ab52T.Errorf("d=%d l=%d setpoint leaked", ab52Declared, ab52Len)
			}
			if ab52View.Ab52Incomplete != (!ab52Set || ab52Len < ab52Declared) {
				ab52T.Errorf("d=%d l=%d incomplete=%v", ab52Declared, ab52Len, ab52View.Ab52Incomplete)
			}
			ab52Executed++
		}
	}
	if ab52Executed != 27 {
		ab52T.Fatalf("executed %d", ab52Executed)
	}
}

func TestAb52A5LiteralPolicyAndSetpointQualification(ab52T *testing.T) {
	ab52View := Ab52DecodeA5(7, []byte{0x03, 0x17, 0x3B, 0x18, 0x3C, 0xFF, 0xFF})
	if ab52View.Ab52Status.Value != 3 || ab52View.Ab52StatusAccepted {
		ab52T.Errorf("status %+v accepted=%v", ab52View.Ab52Status, ab52View.Ab52StatusAccepted)
	}
	if ab52View.Ab52StartHour.Value != 23 || ab52View.Ab52StartMinute.Value != 59 || ab52View.Ab52EndHour.Value != 24 || ab52View.Ab52EndMinute.Value != 60 {
		ab52T.Errorf("times altered: %+v %+v %+v %+v", ab52View.Ab52StartHour, ab52View.Ab52StartMinute, ab52View.Ab52EndHour, ab52View.Ab52EndMinute)
	}
	if !ab52View.Ab52SetpointPresent || ab52View.Ab52SetpointRaw != 65535 || ab52View.Ab52SetpointTenths != 6553.5 {
		ab52T.Errorf("setpoint %+v", ab52View)
	}
	if ab52View.Ab52Incomplete || len(ab52View.Ab52Issues) != 0 {
		ab52T.Errorf("status 3 must not make decoding fail: %+v", ab52View.Ab52Issues)
	}
	if ab52View.Ab52SetpointGrade != Ab52GradeHypothesis || ab52View.Ab52SetpointQualified != Ab52GradeUnknown {
		ab52T.Errorf("grades %q %q", ab52View.Ab52SetpointGrade, ab52View.Ab52SetpointQualified)
	}
	ab52Plain := Ab52DecodeA5(7, []byte{0, 0, 0, 0, 0, 0xFF, 0xFF})
	ab52Other := Ab52DecodeA5(7, []byte{0, 0, 0, 0, 0, 0xFE, 0xFF})
	if ab52Plain.Ab52SetpointRaw != 0xFFFF || ab52Other.Ab52SetpointRaw != 0xFFFE || ab52Plain.Ab52SetpointQualified != Ab52GradeUnknown {
		ab52T.Errorf("sentinel inferred: %+v %+v", ab52Plain, ab52Other)
	}
}

func TestAb52OwnedMemory(ab52T *testing.T) {
	ab52Src4 := []byte{0, 1, 2, 3, 4, 5, 6, 7, 9}
	ab52V4 := Ab52DecodeA4(9, ab52Src4)
	ab52Src4[1], ab52Src4[8] = 0xEE, 0xEE
	if ab52V4.Ab52Raw[1] != 1 || ab52V4.Ab52DayCounts[0] != 1 || ab52V4.Ab52EighthCount.Value != 9 {
		ab52T.Fatalf("a4 aliases source: %+v", ab52V4)
	}
	ab52V4.Ab52Raw[2] = 0xAA
	ab52V4.Ab52DayCounts[3] = 0xBB
	ab52V4.Ab52Issues = append(ab52V4.Ab52Issues, "x")
	ab52Again := Ab52DecodeA4(9, []byte{0, 1, 2, 3, 4, 5, 6, 7, 9})
	if ab52Again.Ab52Raw[2] != 2 || ab52Again.Ab52DayCounts[3] != 4 || len(ab52Again.Ab52Issues) != 0 {
		ab52T.Fatalf("a4 shared state: %+v", ab52Again)
	}
	if &ab52V4.Ab52Raw[0] == &ab52V4.Ab52DayCounts[0] {
		ab52T.Fatal("a4 raw and day counts share memory")
	}
	ab52Src5 := []byte{0, 1, 2, 3, 4, 5, 6}
	ab52V5 := Ab52DecodeA5(7, ab52Src5)
	ab52Src5[5] = 0xEE
	if ab52V5.Ab52Raw[5] != 5 || ab52V5.Ab52SetpointRaw != 0x0605 {
		ab52T.Fatalf("a5 aliases source: %+v", ab52V5)
	}
	ab52V5.Ab52Raw[6] = 0
	if ab52Fresh := Ab52DecodeA5(7, []byte{0, 1, 2, 3, 4, 5, 6}); ab52Fresh.Ab52Raw[6] != 6 {
		ab52T.Fatal("a5 shared state")
	}
}

func TestAb52MalformedExtentsNeverSucceed(ab52T *testing.T) {
	for _, ab52Declared := range []int{-1, -100, 0} {
		ab52V4 := Ab52DecodeA4(ab52Declared, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8})
		ab52V5 := Ab52DecodeA5(ab52Declared, nil)
		if len(ab52V4.Ab52Raw) != 0 || !ab52V4.Ab52Incomplete || ab52V4.Ab52WeekComplete || ab52V4.Ab52WeekValid || ab52V4.Ab52EighthCount.Present || len(ab52V4.Ab52Issues) == 0 {
			ab52T.Errorf("a4 declared=%d: %+v", ab52Declared, ab52V4)
		}
		if len(ab52V5.Ab52Raw) != 0 || !ab52V5.Ab52Incomplete || ab52V5.Ab52SetpointPresent || ab52V5.Ab52Status.Present {
			ab52T.Errorf("a5 declared=%d: %+v", ab52Declared, ab52V5)
		}
	}
	ab52Huge := Ab52DecodeA4(1<<30, []byte{0, 1, 2})
	if len(ab52Huge.Ab52Raw) != 3 || !ab52Huge.Ab52Incomplete {
		ab52T.Errorf("huge declared: %+v", ab52Huge)
	}
}
