package vaillant

import (
	ab15bytes "bytes"
	ab15json "encoding/json"
	ab15reflect "reflect"
	ab15strings "strings"
	ab15testing "testing"
)

func ab15Pointer[ab15T any](ab15Value ab15T) *ab15T { return &ab15Value }

func ab15Equal(ab15T *ab15testing.T, ab15Got, ab15Want any) {
	ab15T.Helper()
	if !ab15reflect.DeepEqual(ab15Got, ab15Want) {
		ab15T.Fatalf("ab15 got %#v; want %#v", ab15Got, ab15Want)
	}
}

func ab15Key(ab15Operation uint8, ab15Field, ab15Alternative string) Ab15Key {
	return Ab15Key{Ab15Family: "vaillant", Ab15Variant: "B504", Ab15Primary: 181, Ab15Secondary: ab15Pointer(uint8(4)), Ab15Operation: ab15Pointer(ab15Operation), Ab15Direction: "none", Ab15Field: ab15Field, Ab15Alternative: ab15Alternative, Ab15InstanceScope: "dynamic"}
}

func ab15Get(ab15T *ab15testing.T, ab15Key Ab15Key) ab15Row {
	ab15T.Helper()
	ab15Raw, ab15Found, ab15Err := Ab15Lookup(ab15Key)
	if ab15Err != nil || !ab15Found {
		ab15T.Fatalf("ab15 lookup %v: %v, %v", ab15Key, ab15Found, ab15Err)
	}
	var ab15Entry ab15Row
	if ab15Err := ab15json.Unmarshal(ab15Raw, &ab15Entry); ab15Err != nil {
		ab15T.Fatal(ab15Err)
	}
	ab15Equal(ab15T, ab15Entry.Ab15Key, ab15Key)
	return ab15Entry
}

func TestAb15FiniteRoles(ab15T *ab15testing.T) {
	ab15Cases := []struct {
		ab15Operation                        uint8
		ab15Field, ab15Alternative, ab15Role string
		ab15Applies                          []string
		ab15Relations                        []string
	}{
		{11, "reload-delay", "none", "solar-storage-charging-delay", []string{"solar-storage"}, []string{"reload delay concerns solar storage charging"}},
		{11, "parallel-tank-filling", "none", "global-mixer-circuit-filling", []string{"all-mixer-circuits"}, []string{"parallel-tank filling applies globally to mixer circuits"}},
		{45, "room-setpoint", "none", "inverter-indoor-room-setpoint", []string{"inverter-indoor-circuit"}, []string{"room temperature is the regulation input"}},
		{45, "room-temperature", "none", "regulation-input", []string{"inverter-indoor-circuit"}, []string{"regulation input for room setpoint"}},
		{13, "system-cooling-support", "none", "system-cooling-support", []string{"system"}, []string{"request bit 0 expresses cooling support", "system generation occupies a different byte"}},
		{13, "system-generation", "none", "system-generation", []string{"system"}, []string{"different byte from cooling-support request bit 0"}},
		{13, "dhw-request", "dhw", "dhw-request", []string{"dhw-circuit"}, []string{"DHW interpretation applies only to a DHW circuit"}},
		{13, "circuit-temperature", "burner", "common-sensor", []string{"burner-circuit"}, []string{"circuit role determines temperature meaning"}},
		{13, "circuit-temperature", "mixer", "flow", []string{"mixer-circuit"}, []string{"circuit role determines temperature meaning"}},
		{13, "circuit-temperature", "dhw", "tank", []string{"dhw-circuit"}, []string{"circuit role determines temperature meaning"}},
		{13, "desired-temperature", "burner", "desired-common-sensor", []string{"burner-circuit"}, []string{"circuit role determines desired-temperature meaning"}},
		{13, "desired-temperature", "mixer", "desired-flow", []string{"mixer-circuit"}, []string{"circuit role determines desired-temperature meaning"}},
		{13, "desired-temperature", "dhw", "desired-tank", []string{"dhw-circuit"}, []string{"circuit role determines desired-temperature meaning"}},
		{24, "desired-flow", "unsuppressed", "unsuppressed-desired-flow", []string{"circuit"}, []string{"independent of suppressed desired flow during tank charging"}},
		{24, "desired-flow", "suppressed", "suppressed-desired-flow", []string{"circuit"}, []string{"zero during tank charging does not replace unsuppressed demand"}},
		{1, "drying-day", "none", "drying-program-indication", []string{"circuit"}, []string{"nonzero drying day identifies a drying program", "zero has no nonzero-program indication"}},
		{51, "averaged-outside-temperature", "none", "correction-included-average", []string{"system"}, []string{"averaged outside temperature includes configured correction"}},
		{33, "maximum-solar-temperature", "none", "maximum-solar-temperature", []string{"all-solar-circuits"}, []string{"maximum solar temperature spans solar circuits"}},
		{33, "tank-state", "none", "temperature-to-maximum-percentage", []string{"solar-circuit"}, []string{"distinct calculated quantity from collector-to-tank switching ratio"}},
		{33, "solar-yield", "none", "collector-to-tank-switching-ratio", []string{"solar-circuit"}, []string{"distinct calculated quantity from temperature-to-maximum percentage"}},
	}
	ab15Data, ab15Err := Ab15Entries()
	if ab15Err != nil {
		ab15T.Fatal(ab15Err)
	}
	var ab15Rows []ab15json.RawMessage
	if ab15Err = ab15json.Unmarshal(ab15Data, &ab15Rows); ab15Err != nil {
		ab15T.Fatal(ab15Err)
	}
	ab15Equal(ab15T, len(ab15Rows), len(ab15Cases))
	for _, ab15Case := range ab15Cases {
		ab15T.Run("ab15-"+ab15Case.ab15Field+"-"+ab15Case.ab15Alternative, func(ab15T *ab15testing.T) {
			ab15Entry := ab15Get(ab15T, ab15Key(ab15Case.ab15Operation, ab15Case.ab15Field, ab15Case.ab15Alternative))
			ab15Equal(ab15T, ab15Entry.Ab15ID, "ab15-"+ab15Case.ab15Field+"-"+ab15Case.ab15Alternative)
			ab15Equal(ab15T, ab15Entry.Ab15Role, ab15Case.ab15Role)
			ab15Equal(ab15T, ab15Entry.Ab15AppliesTo, ab15Case.ab15Applies)
			ab15Equal(ab15T, ab15Entry.Ab15Relations, ab15Case.ab15Relations)
			ab15Equal(ab15T, ab15Entry.Ab15Grade, "Hypothesis")
			ab15Equal(ab15T, ab15Entry.Ab15MetadataPresence, "Proven")
			ab15Equal(ab15T, ab15Entry.Ab15PhysicalQualification, "Unknown")
			ab15Equal(ab15T, ab15Entry.Ab15NativeContract, "protocols/ebus-services/ebus-overview.md")
			ab15Unknown := []string{"byte positions", "widths", "wire direction", "generation values", "suppression algorithms", "native support", "physical effects", "persistence", "delivery"}
			if ab15Case.ab15Field == "averaged-outside-temperature" {
				ab15Unknown = append(ab15Unknown, "averaging and correction arithmetic")
			}
			if ab15Case.ab15Field == "tank-state" || ab15Case.ab15Field == "solar-yield" {
				ab15Unknown = append(ab15Unknown, "exact formula")
			}
			ab15Equal(ab15T, ab15Entry.Ab15Unknown, ab15Unknown)
			var ab15Raw map[string]any
			if ab15Err := ab15json.Unmarshal(ab15Entry.Ab15Raw, &ab15Raw); ab15Err != nil {
				ab15T.Fatal(ab15Err)
			}
			ab15Equal(ab15T, ab15Raw, map[string]any{"form": "opaque", "preserve": true, "interpretation": "Unknown"})
			ab15Whole, ab15Found, ab15Err := Ab15Lookup(ab15Entry.Ab15Key)
			if ab15Err != nil || !ab15Found {
				ab15T.Fatal(ab15Err)
			}
			ab15Matches := 0
			for _, ab15Original := range ab15Rows {
				var ab15OriginalEntry ab15Row
				if ab15Err := ab15json.Unmarshal(ab15Original, &ab15OriginalEntry); ab15Err != nil {
					ab15T.Fatal(ab15Err)
				}
				if ab15reflect.DeepEqual(ab15OriginalEntry.Ab15Key, ab15Entry.Ab15Key) {
					ab15Equal(ab15T, []byte(ab15Whole), []byte(ab15Original))
					ab15Matches++
				}
			}
			ab15Equal(ab15T, ab15Matches, 1)
			ab15Condition := "none"
			var ab15Bit *uint8
			if ab15Case.ab15Field == "drying-day" {
				ab15Condition = "nonzero"
			}
			if ab15Case.ab15Field == "system-cooling-support" {
				ab15Bit = ab15Pointer(uint8(0))
			}
			ab15Equal(ab15T, ab15Entry.Ab15Condition, ab15Condition)
			ab15Equal(ab15T, ab15Entry.Ab15Bit, ab15Bit)
		})
	}
}

func TestAb15SyntheticCells(ab15T *ab15testing.T) {
	for _, ab15Case := range []struct{ ab15Context, ab15Role string }{{"burner", "common-sensor"}, {"mixer", "flow"}, {"dhw", "tank"}} {
		ab15Entry := ab15Get(ab15T, ab15Key(13, "circuit-temperature", ab15Case.ab15Context))
		for _, ab15Token := range []int{40, 41} {
			ab15Cell := map[string]int{ab15Entry.Ab15Role: ab15Token}
			ab15Equal(ab15T, ab15Cell, map[string]int{ab15Case.ab15Role: ab15Token})
		}
	}
	ab15DHW := ab15Get(ab15T, ab15Key(13, "dhw-request", "dhw"))
	ab15Equal(ab15T, ab15DHW.Ab15AppliesTo, []string{"dhw-circuit"})
	ab15Raw, ab15Found, ab15Err := Ab15Lookup(ab15Key(13, "dhw-request", "mixer"))
	if ab15Err != nil || ab15Found || ab15Raw != nil {
		ab15T.Fatal("ab15 mixer acquired DHW interpretation")
	}
	for _, ab15Case := range []struct{ ab15Unsuppressed, ab15Suppressed int }{{45, 0}, {46, 0}, {45, 1}} {
		ab15Demand := ab15Get(ab15T, ab15Key(24, "desired-flow", "unsuppressed"))
		ab15Suppressed := ab15Get(ab15T, ab15Key(24, "desired-flow", "suppressed"))
		ab15Cells := map[string]int{ab15Demand.Ab15Role: ab15Case.ab15Unsuppressed, ab15Suppressed.Ab15Role: ab15Case.ab15Suppressed}
		ab15Equal(ab15T, ab15Cells, map[string]int{"unsuppressed-desired-flow": ab15Case.ab15Unsuppressed, "suppressed-desired-flow": ab15Case.ab15Suppressed})
	}
	ab15Drying := ab15Get(ab15T, ab15Key(1, "drying-day", "none"))
	for _, ab15Case := range []struct {
		ab15Day       int
		ab15Indicated bool
	}{{0, false}, {1, true}, {2, true}} {
		ab15Equal(ab15T, ab15Drying.Ab15Condition == "nonzero" && ab15Case.ab15Day != 0, ab15Case.ab15Indicated)
	}
	for _, ab15Case := range []struct{ ab15Percentage, ab15Ratio int }{{50, 50}, {51, 50}, {50, 51}} {
		ab15Tank := ab15Get(ab15T, ab15Key(33, "tank-state", "none"))
		ab15Yield := ab15Get(ab15T, ab15Key(33, "solar-yield", "none"))
		ab15Cells := map[string]int{ab15Tank.Ab15Role: ab15Case.ab15Percentage, ab15Yield.Ab15Role: ab15Case.ab15Ratio}
		ab15Equal(ab15T, ab15Cells, map[string]int{"temperature-to-maximum-percentage": ab15Case.ab15Percentage, "collector-to-tank-switching-ratio": ab15Case.ab15Ratio})
	}
}

func TestAb15EveryKeyAxis(ab15T *ab15testing.T) {
	ab15Base := ab15Key(11, "reload-delay", "none")
	ab15Controls := []struct {
		ab15Name  string
		ab15Value any
	}{
		{"Ab15Family", "other-family"}, {"Ab15Variant", "other-variant"}, {"Ab15Primary", uint8(180)},
		{"Ab15Secondary", ab15Pointer(uint8(5))}, {"Ab15Operation", ab15Pointer(uint8(12))},
		{"Ab15Direction", "request"}, {"Ab15Field", "other-field"}, {"Ab15Alternative", "other-role"},
		{"Ab15Group", ab15Pointer(uint8(1))}, {"Ab15InstanceScope", "00"}, {"Ab15Register", ab15Pointer(uint16(1))},
		{"Ab15Secondary", (*uint8)(nil)}, {"Ab15Secondary", ab15Pointer(uint8(0))},
		{"Ab15Operation", (*uint8)(nil)}, {"Ab15Operation", ab15Pointer(uint8(0))},
		{"Ab15Group", ab15Pointer(uint8(0))}, {"Ab15Register", ab15Pointer(uint16(0))},
	}
	ab15Original, ab15Found, ab15Err := Ab15Lookup(ab15Base)
	if ab15Err != nil || !ab15Found {
		ab15T.Fatal("ab15 positive control missing", ab15Err)
	}
	for _, ab15Control := range ab15Controls {
		ab15T.Run("ab15-"+ab15Control.ab15Name, func(ab15T *ab15testing.T) {
			ab15Changed := ab15Base
			ab15reflect.ValueOf(&ab15Changed).Elem().FieldByName(ab15Control.ab15Name).Set(ab15reflect.ValueOf(ab15Control.ab15Value))
			ab15Raw, ab15Found, ab15Err := Ab15Lookup(ab15Changed)
			if ab15Err != nil || ab15Found || ab15Raw != nil {
				ab15T.Fatal("ab15 inexact key admitted", ab15Err)
			}
			ab15Again, ab15Found, ab15Err := Ab15Lookup(ab15Base)
			if ab15Err != nil || !ab15Found {
				ab15T.Fatal("ab15 positive control lost", ab15Err)
			}
			ab15Equal(ab15T, ab15Again, ab15Original)
			ab15One := ab15Get(ab15T, ab15Base)
			ab15Two := ab15One
			ab15Two.Ab15ID = "ab15-synthetic-axis-control"
			ab15Two.Ab15Key = ab15Changed
			ab15Fixture, ab15Err := ab15json.Marshal([]ab15Row{ab15One, ab15Two})
			if ab15Err != nil {
				ab15T.Fatal(ab15Err)
			}
			if _, ab15Err = ab15Entries(ab15Fixture); ab15Err != nil {
				ab15T.Fatal("ab15 distinct identities rejected", ab15Err)
			}
			for _, ab15Key := range []Ab15Key{ab15Base, ab15Changed} {
				ab15Raw, ab15Found, ab15Err := ab15Lookup(ab15Fixture, ab15Key)
				if ab15Err != nil || !ab15Found {
					ab15T.Fatal("ab15 seam control missing", ab15Err)
				}
				var ab15Entry ab15Row
				if ab15Err = ab15json.Unmarshal(ab15Raw, &ab15Entry); ab15Err != nil {
					ab15T.Fatal(ab15Err)
				}
				ab15Equal(ab15T, ab15Entry.Ab15Key, ab15Key)
			}
		})
	}
}

func TestAb15AbsentVersusZero(ab15T *ab15testing.T) {
	for _, ab15Axis := range []string{"Ab15Secondary", "Ab15Operation", "Ab15Group", "Ab15Register"} {
		ab15T.Run("ab15-"+ab15Axis, func(ab15T *ab15testing.T) {
			ab15Absent := ab15Get(ab15T, ab15Key(11, "reload-delay", "none"))
			ab15Field := ab15reflect.ValueOf(&ab15Absent.Ab15Key).Elem().FieldByName(ab15Axis)
			ab15Field.Set(ab15reflect.Zero(ab15Field.Type()))
			ab15Zero := ab15Absent
			ab15Zero.Ab15ID = "ab15-synthetic-zero"
			ab15Field = ab15reflect.ValueOf(&ab15Zero.Ab15Key).Elem().FieldByName(ab15Axis)
			ab15Field.Set(ab15reflect.New(ab15Field.Type().Elem()))
			ab15Fixture, ab15Err := ab15json.Marshal([]ab15Row{ab15Absent, ab15Zero})
			if ab15Err != nil {
				ab15T.Fatal(ab15Err)
			}
			for _, ab15Entry := range []ab15Row{ab15Absent, ab15Zero} {
				ab15Raw, ab15Found, ab15Err := ab15Lookup(ab15Fixture, ab15Entry.Ab15Key)
				if ab15Err != nil || !ab15Found {
					ab15T.Fatal("ab15 absent/zero control", ab15Err)
				}
				var ab15Got ab15Row
				if ab15Err = ab15json.Unmarshal(ab15Raw, &ab15Got); ab15Err != nil {
					ab15T.Fatal(ab15Err)
				}
				ab15Expected, ab15Err := ab15json.Marshal(ab15Entry)
				if ab15Err != nil {
					ab15T.Fatal(ab15Err)
				}
				ab15Equal(ab15T, []byte(ab15Raw), ab15Expected)
				ab15Equal(ab15T, ab15Got.Ab15Key, ab15Entry.Ab15Key)
			}
		})
	}
}

func TestAb15ValidationErrors(ab15T *ab15testing.T) {
	ab15Entry := ab15Get(ab15T, ab15Key(11, "reload-delay", "none"))
	ab15Duplicate := ab15Entry
	ab15Duplicate.Ab15ID = "ab15-duplicate-key"
	ab15Fixture, ab15Err := ab15json.Marshal([]ab15Row{ab15Entry, ab15Duplicate})
	if ab15Err != nil {
		ab15T.Fatal(ab15Err)
	}
	ab15Inputs := []string{"{", "null", "[]", "{}", "[null]", "[{}]", string(ab15Fixture), string(ab15Fixture) + " []"}
	ab15Valid, ab15Err := ab15json.Marshal([]ab15Row{ab15Entry})
	if ab15Err != nil {
		ab15T.Fatal(ab15Err)
	}
	ab15Inputs = append(ab15Inputs, ab15strings.Replace(string(ab15Valid), `"Hypothesis"`, `"Proven"`, 1), ab15strings.Replace(string(ab15Valid), `"physical_qualification":"Unknown"`, `"physical_qualification":"Proven"`, 1), ab15strings.Replace(string(ab15Valid), `"raw":{"form":"opaque","preserve":true,"interpretation":"Unknown"}`, `"raw":null`, 1))
	for _, ab15Axis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
		var ab15Object []map[string]ab15json.RawMessage
		if ab15Err := ab15json.Unmarshal(ab15Valid, &ab15Object); ab15Err != nil {
			ab15T.Fatal(ab15Err)
		}
		var ab15Identity map[string]ab15json.RawMessage
		if ab15Err := ab15json.Unmarshal(ab15Object[0]["key"], &ab15Identity); ab15Err != nil {
			ab15T.Fatal(ab15Err)
		}
		delete(ab15Identity, ab15Axis)
		ab15Object[0]["key"], ab15Err = ab15json.Marshal(ab15Identity)
		if ab15Err != nil {
			ab15T.Fatal(ab15Err)
		}
		ab15Missing, ab15Err := ab15json.Marshal(ab15Object)
		if ab15Err != nil {
			ab15T.Fatal(ab15Err)
		}
		ab15Inputs = append(ab15Inputs, string(ab15Missing))
	}
	for ab15Index, ab15Input := range ab15Inputs {
		ab15Raw, ab15Err := ab15Entries([]byte(ab15Input))
		if ab15Err == nil || ab15Raw != nil {
			ab15T.Fatalf("ab15 invalid entries accepted: %d", ab15Index)
		}
		ab15Raw, ab15Found, ab15Err := ab15Lookup([]byte(ab15Input), ab15Entry.Ab15Key)
		if ab15Err == nil || ab15Raw != nil || ab15Found {
			ab15T.Fatalf("ab15 invalid lookup accepted: %d", ab15Index)
		}
	}
}

func TestAb15IndependentBuffers(ab15T *ab15testing.T) {
	ab15Before, ab15Err := Ab15Entries()
	if ab15Err != nil {
		ab15T.Fatal(ab15Err)
	}
	ab15Other, ab15Err := Ab15Entries()
	if ab15Err != nil {
		ab15T.Fatal(ab15Err)
	}
	for ab15Index := range ab15Other {
		ab15Other[ab15Index] = 'x'
	}
	ab15After, ab15Err := Ab15Entries()
	if ab15Err != nil {
		ab15T.Fatal(ab15Err)
	}
	ab15Equal(ab15T, ab15After, ab15Before)
	ab15Key := ab15Key(11, "reload-delay", "none")
	ab15One, ab15Found, ab15Err := Ab15Lookup(ab15Key)
	if ab15Err != nil || !ab15Found {
		ab15T.Fatal(ab15Err)
	}
	ab15Two, ab15Found, ab15Err := Ab15Lookup(ab15Key)
	if ab15Err != nil || !ab15Found {
		ab15T.Fatal(ab15Err)
	}
	for ab15Index := range ab15One {
		ab15One[ab15Index] = 'x'
	}
	ab15Three, ab15Found, ab15Err := Ab15Lookup(ab15Key)
	if ab15Err != nil || !ab15Found {
		ab15T.Fatal(ab15Err)
	}
	ab15Equal(ab15T, ab15Three, ab15Two)
	ab15Equal(ab15T, ab15bytes.Contains(ab15Before, ab15Three), true)
	ab15Source := append([]byte(nil), ab15Before...)
	ab15Copy, ab15Err := ab15Entries(ab15Source)
	if ab15Err != nil {
		ab15T.Fatal(ab15Err)
	}
	for ab15Index := range ab15Copy {
		ab15Copy[ab15Index] = 'x'
	}
	ab15Equal(ab15T, ab15Source, []byte(ab15Before))
	ab15Copy, ab15Found, ab15Err = ab15Lookup(ab15Source, ab15Key)
	if ab15Err != nil || !ab15Found {
		ab15T.Fatal(ab15Err)
	}
	for ab15Index := range ab15Copy {
		ab15Copy[ab15Index] = 'x'
	}
	ab15Equal(ab15T, ab15Source, []byte(ab15Before))
}
