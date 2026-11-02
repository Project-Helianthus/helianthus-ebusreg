package vaillant

import (
	ab72binary "encoding/binary"
	ab72json "encoding/json"
	ab72reflect "reflect"
	ab72strings "strings"
	ab72testing "testing"
)

func ab72Pointer[ab72Type any](ab72Value ab72Type) *ab72Type { return &ab72Value }
func ab72Equal(ab72T *ab72testing.T, ab72Got, ab72Want any) {
	ab72T.Helper()
	if !ab72reflect.DeepEqual(ab72Got, ab72Want) {
		ab72T.Fatalf("got %#v; want %#v", ab72Got, ab72Want)
	}
}
func ab72Get(ab72T *ab72testing.T, ab72Variant string, ab72Operation uint8, ab72Side, ab72Field string, ab72Selector *uint8) Ab72Entry {
	ab72T.Helper()
	ab72Entry, ab72Found, ab72Err := Ab72Lookup(Ab72Key{Ab72Family: "vaillant", Ab72Variant: ab72Variant, Ab72Operation: ab72Pointer(ab72Operation), Ab72Direction: ab72Side, Ab72Field: ab72Field, Ab72Selector: ab72Selector})
	if ab72Err != nil || !ab72Found {
		ab72T.Fatalf("lookup %s/%d/%s/%s: %v %v", ab72Variant, ab72Operation, ab72Side, ab72Field, ab72Found, ab72Err)
	}
	return ab72Entry
}
func TestAb72StatisticsSides(ab72T *ab72testing.T) {
	for _, ab72Case := range []struct {
		ab72Variant, ab72Side string
		ab72Length            uint8
		ab72Value             bool
		ab72Access            *uint8
	}{
		{"ab72-statistics-read", "request", 7, false, ab72Pointer(uint8(2))},
		{"ab72-statistics-read", "positive", 11, true, nil},
		{"ab72-statistics-write", "request", 11, true, ab72Pointer(uint8(2))},
		{"ab72-statistics-write", "positive", 7, false, ab72Pointer(uint8(3))},
	} {
		ab72Equal(ab72T, ab72Get(ab72T, ab72Case.ab72Variant, 16, ab72Case.ab72Side, "layout", nil).Ab72Length, ab72Pointer(ab72Case.ab72Length))
		ab72Flags := ab72Get(ab72T, ab72Case.ab72Variant, 16, ab72Case.ab72Side, "flags", nil)
		ab72Expected := map[string]uint8{"time-base-low": 0, "time-base-high": 1}
		if ab72Case.ab72Access != nil {
			ab72Expected["access"] = *ab72Case.ab72Access
		}
		ab72Equal(ab72T, ab72Flags.Ab72Bits, ab72Expected)
		ab72Equal(ab72T, ab72Flags.Ab72Values["time-base:0"], "total")
		ab72Equal(ab72T, ab72Flags.Ab72Values["time-base:1"], "daily")
		ab72Equal(ab72T, ab72Flags.Ab72Values["time-base:2"], "monthly")
		ab72Equal(ab72T, ab72Flags.Ab72Values["time-base:3"], "yearly")
		ab72Key := ab72Flags.Ab72Key
		ab72Key.Ab72Field = "energy-value"
		ab72Entry, ab72Found, ab72Err := Ab72Lookup(ab72Key)
		ab72Equal(ab72T, ab72Err, nil)
		ab72Equal(ab72T, ab72Found, ab72Case.ab72Value)
		if ab72Found {
			ab72Equal(ab72T, ab72Entry.Ab72Encoding, ab72Pointer("binary32 LE"))
			ab72Equal(ab72T, ab72Entry.Ab72Unit, ab72Pointer("Wh"))
		}
		for _, ab72Index := range []uint8{0, 255} {
			ab72Period := ab72Get(ab72T, ab72Case.ab72Variant, 16, ab72Case.ab72Side, "period-index", ab72Pointer(ab72Index))
			ab72Energy := ab72Get(ab72T, ab72Case.ab72Variant, 16, ab72Case.ab72Side, "energy-index", ab72Pointer(ab72Index))
			ab72PeriodRule, ab72EnergyRule := "Select table row; ignore date.", "Ignore energy type and operation mode."
			if ab72Index == 255 {
				ab72PeriodRule = "Select by date."
				ab72EnergyRule = "Select energy type and operation mode."
			}
			ab72Equal(ab72T, ab72Period.Ab72Rules, []string{"Total ignores period index and date.", ab72PeriodRule})
			ab72Equal(ab72T, ab72Energy.Ab72Rules, []string{ab72EnergyRule})
			ab72Equal(ab72T, ab72Period.Ab72Unknown, []string{"date packing", "offsets", "defaults"})
			ab72Equal(ab72T, ab72Period.Ab72Encoding, (*string)(nil))
		}
		if ab72Case.ab72Side == "positive" {
			ab72Return := ab72Get(ab72T, ab72Case.ab72Variant, 16, "positive", "return-nibble", nil)
			ab72Equal(ab72T, ab72Return.Ab72Values, map[string]string{"0": "value-contained", "1": "NOK"})
			ab72Equal(ab72T, ab72Return.Ab72Unknown, []string{"meanings 2..15", "offsets"})
		}
	}
}
func TestAb72Index(ab72T *ab72testing.T) {
	ab72Equal(ab72T, ab72Get(ab72T, "ab72-statistics-index", 17, "request", "layout", nil).Ab72Length, ab72Pointer(uint8(0)))
	ab72Fields := []string{"Version", "Build", "TotalEnergies", "Days", "DailyEnergies", "Months", "MonthlyEnergies", "Years", "YearlyEnergies"}
	ab72Order := ab72Get(ab72T, "ab72-statistics-index", 17, "positive", "ordered-fields", nil)
	ab72Equal(ab72T, ab72Order.Ab72Order, ab72Fields)
	ab72Equal(ab72T, ab72Order.Ab72Rules, []string{"Counts are returned metadata, not installed counts or enforced bounds."})
	for ab72Index, ab72Field := range ab72Fields {
		ab72Entry := ab72Get(ab72T, "ab72-statistics-index", 17, "positive", ab72Field, nil)
		ab72Equal(ab72T, ab72Entry.Ab72Position, ab72Pointer(uint8(ab72Index)))
		ab72Equal(ab72T, ab72Entry.Ab72Width, ab72Pointer(uint8(8)))
		ab72Equal(ab72T, ab72Entry.Ab72Encoding, ab72Pointer("u8"))
	}
}
func TestAb72MeterClearAlternatives(ab72T *ab72testing.T) {
	for _, ab72Variant := range []string{"ab72-meter-scalar-u8", "ab72-meter-low-six-bits", "ab72-clear-scalar-u8", "ab72-clear-low-six-bits"} {
		ab72Equal(ab72T, ab72Get(ab72T, ab72Variant, 18, "request", "flags", nil).Ab72Bits, map[string]uint8{"total": 0, "daily": 1, "monthly": 2, "yearly": 3})
		ab72Selector := ab72Get(ab72T, ab72Variant, 18, "request", "energy-selector", nil)
		ab72Expected := map[string]uint8{}
		ab72Rules := []string{"Scalar u8 selector."}
		if ab72strings.Contains(ab72Variant, "low-six") {
			ab72Expected = map[string]uint8{"sum-6": 6, "sum-7": 7}
			ab72Rules = []string{"Low six bits selector; independent sum bits 6 and 7."}
		}
		ab72Equal(ab72T, ab72Selector.Ab72Bits, ab72Expected)
		ab72Equal(ab72T, ab72Selector.Ab72Rules, ab72Rules)
		ab72Equal(ab72T, ab72Get(ab72T, ab72Variant, 18, "request", "mode", nil).Ab72Encoding, ab72Pointer("u8"))
		ab72Positive := ab72Get(ab72T, ab72Variant, 18, "positive", "flags", nil)
		ab72Equal(ab72T, ab72Positive.Ab72Bits, map[string]uint8{"not-supported": 0, "energy-unavailable": 1, "mode-unavailable": 2})
		ab72Equal(ab72T, ab72Positive.Ab72Conflicts, []string{"Nested and outer duplicate flags; packing precedence Unknown."})
		ab72Equal(ab72T, ab72Positive.Ab72Unknown, []string{"packing precedence", "offsets", "defaults", "full variant-specific meter/clear tables"})
	}
}
func TestAb72YieldAndConflicts(ab72T *ab72testing.T) {
	for _, ab72Case := range []struct{ ab72Variant, ab72Side string }{{"ab72-yield-read", "positive"}, {"ab72-yield-write", "request"}} {
		ab72Entry := ab72Get(ab72T, ab72Case.ab72Variant, 24, ab72Case.ab72Side, "yield-value", nil)
		ab72Equal(ab72T, ab72Entry.Ab72Encoding, ab72Pointer("LEu32"))
		ab72Equal(ab72T, ab72Entry.Ab72Unit, ab72Pointer("kWh"))
		ab72Equal(ab72T, ab72Entry.Ab72Width, ab72Pointer(uint8(32)))
		ab72Raw := []byte{255, 255, 255, 255}
		ab72Equal(ab72T, ab72binary.LittleEndian.Uint32(ab72Raw), uint32(4294967295))
		ab72Equal(ab72T, ab72Entry.Ab72Rules, []string{"All u32 values retain their raw numeric meaning; no sentinel inferred."})
	}
	for _, ab72Variant := range []string{"ab72-yield-read", "ab72-yield-write"} {
		for _, ab72Side := range []string{"request", "positive"} {
			ab72Equal(ab72T, ab72Get(ab72T, ab72Variant, 24, ab72Side, "yield-type", nil).Ab72Values, map[string]string{"0": "solar", "1": "environment", "2": "energy", "3": "gas-or-oil"})
		}
	}
	ab72Equal(ab72T, ab72Get(ab72T, "ab72-yield-reset", 25, "request", "yield-type", nil).Ab72Values, map[string]string{"0": "solar", "1": "environment", "2": "energy", "3": "gas-or-oil"})
	ab72Feedback := ab72Get(ab72T, "ab72-yield-reset", 25, "positive", "feedback", nil)
	ab72Equal(ab72T, ab72Feedback.Ab72Values, map[string]string{"1": "accepted"})
	ab72Equal(ab72T, ab72Feedback.Ab72Unknown, []string{"other feedback meanings", "effects", "persistence"})
	for _, ab72Field := range []string{"yield-type", "month", "year-index", "mode", "sum"} {
		ab72Entry := ab72Get(ab72T, "ab72-yield-month", 26, "positive", ab72Field, nil)
		ab72Equal(ab72T, ab72Entry.Ab72Position, (*uint8)(nil))
		ab72Equal(ab72T, ab72Entry.Ab72Unit, (*string)(nil))
		ab72Equal(ab72T, ab72Entry.Ab72Conflicts, []string{"Unit-type overlaps second sum byte.", "Month 10/11 conflict."})
		ab72Equal(ab72T, ab72Entry.Ab72Unknown, []string{"offsets", "defaults", "packing precedence", "unit"})
	}
	ab72Equal(ab72T, ab72Get(ab72T, "ab72-yield-month", 26, "positive", "sum", nil).Ab72Encoding, ab72Pointer("LEu32"))
	ab72Equal(ab72T, ab72Get(ab72T, "ab72-yield-month", 26, "positive", "year-index", nil).Ab72Values, map[string]string{"0": "this-year", "1": "last-year"})
	ab72Equal(ab72T, ab72Get(ab72T, "ab72-yield-month", 26, "positive", "mode", nil).Ab72Values, map[string]string{"0": "all", "1": "heating", "2": "DHW"})
}
func TestAb72ExactKeys(ab72T *ab72testing.T) {
	ab72Base := ab72Get(ab72T, "ab72-statistics-read", 16, "request", "period-index", ab72Pointer(uint8(0))).Ab72Key
	for _, ab72Change := range []func(*Ab72Key){
		func(ab72Key *Ab72Key) { ab72Key.Ab72Family = "ab72-other" },
		func(ab72Key *Ab72Key) { ab72Key.Ab72Variant = "ab72-other" },
		func(ab72Key *Ab72Key) { ab72Key.Ab72Operation = ab72Pointer(uint8(0)) },
		func(ab72Key *Ab72Key) { ab72Key.Ab72Operation = nil },
		func(ab72Key *Ab72Key) { ab72Key.Ab72Direction = "ab72-other" },
		func(ab72Key *Ab72Key) { ab72Key.Ab72Selector = nil },
		func(ab72Key *Ab72Key) { ab72Key.Ab72Selector = ab72Pointer(uint8(1)) },
		func(ab72Key *Ab72Key) { ab72Key.Ab72Field = "ab72-other" },
	} {
		ab72Key := ab72Base
		ab72Change(&ab72Key)
		_, ab72Found, ab72Err := Ab72Lookup(ab72Key)
		ab72Equal(ab72T, ab72Err, nil)
		ab72Equal(ab72T, ab72Found, false)
	}
	ab72Entries, ab72Err := Ab72Entries()
	ab72Equal(ab72T, ab72Err, nil)
	ab72Equal(ab72T, len(ab72Entries), 68)
	for _, ab72Entry := range ab72Entries {
		ab72Got, ab72Found, ab72Err := Ab72Lookup(ab72Entry.Ab72Key)
		ab72Equal(ab72T, ab72Err, nil)
		ab72Equal(ab72T, ab72Found, true)
		ab72Equal(ab72T, ab72Got, ab72Entry)
	}
}
func TestAb72ParserAndOwnership(ab72T *ab72testing.T) {
	ab72Entry := ab72Get(ab72T, "ab72-statistics-write", 16, "request", "flags", nil)
	for _, ab72Raw := range []string{"{", "null", "[]", "[{}]", string(ab72Data) + " []", ab72strings.Replace(string(ab72Data), "\"grade\":", "\"extra\": 0, \"grade\":", 1), ab72strings.Replace(string(ab72Data), "Hypothesis", "Proven", 1)} {
		_, ab72Err := ab72Parse([]byte(ab72Raw))
		if ab72Err == nil {
			ab72T.Fatal("invalid input admitted")
		}
	}
	ab72Raw, ab72Err := ab72json.Marshal([]Ab72Entry{ab72Entry, ab72Entry})
	ab72Equal(ab72T, ab72Err, nil)
	_, ab72Err = ab72Parse(ab72Raw)
	if ab72Err == nil {
		ab72T.Fatal("duplicate admitted")
	}
	ab72Raw, ab72Err = ab72json.Marshal([]Ab72Entry{ab72Entry})
	ab72Equal(ab72T, ab72Err, nil)
	ab72Parsed, ab72Err := ab72Parse(ab72Raw)
	ab72Equal(ab72T, ab72Err, nil)
	for ab72Index := range ab72Raw {
		ab72Raw[ab72Index] = 0
	}
	ab72Equal(ab72T, ab72Parsed[0], ab72Entry)
	ab72Entries, ab72Err := Ab72Entries()
	ab72Equal(ab72T, ab72Err, nil)
	for ab72Index := range ab72Entries {
		ab72Owned := &ab72Entries[ab72Index]
		ab72Owned.Ab72Key.Ab72Family = "ab72-mutated"
		*ab72Owned.Ab72Key.Ab72Operation = 0
		if ab72Owned.Ab72Key.Ab72Selector != nil {
			*ab72Owned.Ab72Key.Ab72Selector = 42
		}
		for ab72Key := range ab72Owned.Ab72Bits {
			ab72Owned.Ab72Bits[ab72Key] = 7
		}
		for ab72Key := range ab72Owned.Ab72Values {
			ab72Owned.Ab72Values[ab72Key] = "ab72-mutated"
		}
		for _, ab72Slice := range [][]string{ab72Owned.Ab72Order, ab72Owned.Ab72Rules, ab72Owned.Ab72Conflicts, ab72Owned.Ab72Unknown} {
			for ab72Index := range ab72Slice {
				ab72Slice[ab72Index] = "ab72-mutated"
			}
		}
		if ab72Owned.Ab72Length != nil {
			*ab72Owned.Ab72Length = 255
		}
		if ab72Owned.Ab72Position != nil {
			*ab72Owned.Ab72Position = 255
		}
		if ab72Owned.Ab72Width != nil {
			*ab72Owned.Ab72Width = 1
		}
		if ab72Owned.Ab72Encoding != nil {
			*ab72Owned.Ab72Encoding = "ab72-mutated"
		}
		if ab72Owned.Ab72Unit != nil {
			*ab72Owned.Ab72Unit = "ab72-mutated"
		}
	}
	ab72Fresh, ab72Err := Ab72Entries()
	ab72Equal(ab72T, ab72Err, nil)
	ab72Native, ab72Err := ab72Parse(ab72Data)
	ab72Equal(ab72T, ab72Err, nil)
	ab72Equal(ab72T, ab72Fresh, ab72Native)
	ab72Lookup := ab72Get(ab72T, "ab72-statistics-write", 16, "request", "flags", nil)
	ab72Lookup.Ab72Bits["access"] = 7
	*ab72Lookup.Ab72Key.Ab72Operation = 0
	ab72Equal(ab72T, ab72Get(ab72T, "ab72-statistics-write", 16, "request", "flags", nil), ab72Entry)
	for _, ab72Owned := range ab72Fresh {
		ab72Equal(ab72T, ab72Owned.Ab72Grade, "Hypothesis")
		ab72Equal(ab72T, ab72Owned.Ab72SoftwareContract, "Proven")
		for _, ab72Grade := range []string{ab72Owned.Ab72PhysicalSupport, ab72Owned.Ab72Presence, ab72Owned.Ab72Freshness, ab72Owned.Ab72Execution, ab72Owned.Ab72Effects, ab72Owned.Ab72Persistence, ab72Owned.Ab72Delivery} {
			ab72Equal(ab72T, ab72Grade, "Unknown")
		}
	}
}

func TestAb72FailClosed(ab72T *ab72testing.T) {
	ab72Original := ab72Data
	defer func() { ab72Data = ab72Original }()
	ab72Valid := ab72Get(ab72T, "ab72-statistics-read", 16, "request", "flags", nil)
	ab72Invalid := ab72Valid
	ab72Invalid.Ab72Key.Ab72Field = ""
	ab72Raw, ab72Err := ab72json.Marshal([]Ab72Entry{ab72Valid, ab72Invalid})
	ab72Equal(ab72T, ab72Err, nil)
	ab72Data = ab72Raw
	ab72Entries, ab72Err := Ab72Entries()
	if ab72Err == nil || ab72Entries != nil {
		ab72T.Fatal("partial catalog admitted")
	}
	ab72Entry, ab72Found, ab72Err := Ab72Lookup(ab72Valid.Ab72Key)
	if ab72Err == nil || ab72Found {
		ab72T.Fatal("lookup ignored catalog error")
	}
	ab72Equal(ab72T, ab72Entry, Ab72Entry{})
}
