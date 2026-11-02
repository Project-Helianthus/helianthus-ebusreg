package vaillant

import (
	ab65reflect "reflect"
	ab65testing "testing"
)

func ab65Pointer[ab65T any](ab65Value ab65T) *ab65T { return &ab65Value }

func ab65Get(ab65T *ab65testing.T, ab65Operation uint8, ab65Field string) Ab65Entry {
	ab65T.Helper()
	ab65Entry, ab65OK := Ab65Lookup(Ab65Key{"vaillant", "B510", 0xb5, ab65Pointer(uint8(0x10)), ab65Pointer(ab65Operation), ab65Pointer("request"), ab65Field})
	if !ab65OK {
		ab65T.Fatalf("missing OP%02x %s", ab65Operation, ab65Field)
	}
	return ab65Entry
}

func TestAb65FiniteCatalog(ab65T *ab65testing.T) {
	ab65Catalog := Ab65B510Candidates()
	if ab65Catalog.Ab65MetadataPresence != "Proven" || ab65Catalog.Ab65PhysicalQualification != "Unknown" || ab65Catalog.Ab65RoutableInvocation != "Unknown" || len(ab65Catalog.Ab65Entries) != 12 || len(ab65Catalog.Ab65Limitations) != 5 {
		ab65T.Fatal("finite scope or grades changed")
	}
	ab65Fields := map[uint8][]string{0: {"disable"}, 2: {"remote_status", "CH_setpoint_upper_limit", "DHW_setpoint_upper_limit"}, 3: {"climate", "capacity", "cost", "force"}, 4: {"energy"}, 5: {"DHW_setpoint", "adjustment"}, 6: {"status"}}
	for ab65Operation, ab65Names := range ab65Fields {
		for _, ab65Name := range ab65Names {
			ab65Entry := ab65Get(ab65T, ab65Operation, ab65Name)
			if ab65Entry.Ab65Grade != "Hypothesis" || ab65Entry.Ab65Position != nil || ab65Entry.Ab65WidthBits == nil || *ab65Entry.Ab65WidthBits != 8 || ab65Entry.Ab65Interpretation != "Unknown" {
				ab65T.Fatal("invented qualification or coordinates", ab65Name)
			}
		}
	}
	for _, ab65Entry := range ab65Catalog.Ab65Entries {
		ab65Got, ab65OK := Ab65Lookup(ab65Entry.Ab65Key)
		if !ab65OK || !ab65reflect.DeepEqual(ab65Got, ab65Entry) {
			ab65T.Fatal("accessor/lookup mismatch")
		}
	}
}

func TestAb65ExactKeyAxes(ab65T *ab65testing.T) {
	ab65Control := ab65Get(ab65T, 0, "disable").Ab65Key
	ab65Changes := map[string]func(*Ab65Key){
		"ab65_family":           func(ab65Key *Ab65Key) { ab65Key.Ab65Family = "other" },
		"ab65_family_absent":    func(ab65Key *Ab65Key) { ab65Key.Ab65Family = "" },
		"ab65_variant":          func(ab65Key *Ab65Key) { ab65Key.Ab65Variant = "other" },
		"ab65_primary":          func(ab65Key *Ab65Key) { ab65Key.Ab65Primary = 0 },
		"ab65_secondary":        func(ab65Key *Ab65Key) { ab65Key.Ab65Secondary = ab65Pointer(uint8(0)) },
		"ab65_secondary_absent": func(ab65Key *Ab65Key) { ab65Key.Ab65Secondary = nil },
		"ab65_operation":        func(ab65Key *Ab65Key) { ab65Key.Ab65Operation = ab65Pointer(uint8(1)) },
		"ab65_operation_absent": func(ab65Key *Ab65Key) { ab65Key.Ab65Operation = nil },
		"ab65_direction":        func(ab65Key *Ab65Key) { ab65Key.Ab65Direction = ab65Pointer("response") },
		"ab65_direction_absent": func(ab65Key *Ab65Key) { ab65Key.Ab65Direction = nil },
		"ab65_field":            func(ab65Key *Ab65Key) { ab65Key.Ab65Field = "other" },
	}
	for ab65Name, ab65Change := range ab65Changes {
		ab65T.Run(ab65Name, func(ab65T *ab65testing.T) {
			if _, ab65OK := Ab65Lookup(ab65Control); !ab65OK {
				ab65T.Fatal("invalid positive control")
			}
			ab65Key := ab65Control
			ab65Change(&ab65Key)
			if _, ab65OK := Ab65Lookup(ab65Key); ab65OK {
				ab65T.Fatal("axis wildcarded")
			}
		})
	}
}

func TestAb65ScalesAndEnums(ab65T *ab65testing.T) {
	ab65Capacity := ab65Get(ab65T, 3, "capacity")
	if !ab65reflect.DeepEqual(ab65Capacity.Ab65RawDomain, []uint8{0, 250}) || ab65Capacity.Ab65Divisor == nil || *ab65Capacity.Ab65Divisor != 10 || ab65Capacity.Ab65Unit == nil || *ab65Capacity.Ab65Unit != "kW" {
		ab65T.Fatal("capacity domain/scale")
	}
	for ab65Index, ab65Raw := range []uint8{0, 250, 251} {
		ab65Inside := ab65Raw >= ab65Capacity.Ab65RawDomain[0] && ab65Raw <= ab65Capacity.Ab65RawDomain[1]
		if ab65Inside != []bool{true, true, false}[ab65Index] {
			ab65T.Fatal("capacity domain")
		}
		if ab65Inside && float64(ab65Raw) / *ab65Capacity.Ab65Divisor != []float64{0, 25}[ab65Index] {
			ab65T.Fatal("capacity example")
		}
	}
	for _, ab65Entry := range []Ab65Entry{ab65Get(ab65T, 2, "CH_setpoint_upper_limit"), ab65Get(ab65T, 2, "DHW_setpoint_upper_limit"), ab65Get(ab65T, 5, "DHW_setpoint")} {
		if ab65Entry.Ab65Divisor == nil || float64(1) / *ab65Entry.Ab65Divisor != 0.5 || ab65Entry.Ab65Unit == nil || *ab65Entry.Ab65Unit != "degrees Celsius" {
			ab65T.Fatal("temperature 01")
		}
	}
	ab65Expected := map[string]map[uint8]string{
		"climate": {0: "off", 1: "heating", 2: "cooling"},
		"force":   {0: "normal", 1: "maximum", 2: "minimum", 3: "off", 4: "cooling-leak-test", 5: "cooling-charge-test", 6: "capacity-test", 7: "I/O-test"},
		"energy":  {0: "unknown", 1: "gas", 2: "electricity-one", 3: "electricity-two", 4: "fuel", 5: "pellet"},
	}
	for ab65Field, ab65Values := range ab65Expected {
		ab65Operation := uint8(3)
		if ab65Field == "energy" {
			ab65Operation = 4
		}
		ab65Entry := ab65Get(ab65T, ab65Operation, ab65Field)
		if !ab65reflect.DeepEqual(ab65Entry.Ab65Values, ab65Values) {
			ab65T.Fatal("enum labels", ab65Field)
		}
		if _, ab65Mapped := ab65Entry.Ab65Values[8]; ab65Mapped {
			ab65T.Fatal("unmapped value 8")
		}
	}
	ab65Cost := ab65Get(ab65T, 3, "cost")
	if !ab65reflect.DeepEqual(ab65Cost.Ab65RawDomain, []uint8{0, 250}) || ab65Cost.Ab65Unit == nil || *ab65Cost.Ab65Unit != "reciprocal-MWh" || ab65Cost.Ab65Divisor != nil {
		ab65T.Fatal("cost metadata")
	}
	ab65Remote := ab65Get(ab65T, 2, "remote_status")
	if ab65Remote.Ab65Divisor != nil || ab65Remote.Ab65Values != nil || !ab65reflect.DeepEqual(ab65Remote.Ab65RawDomain, []uint8{0, 255}) {
		ab65T.Fatal("raw remote status")
	}
}

func TestAb65AdjustmentAndReserved(ab65T *ab65testing.T) {
	ab65Adjustment := ab65Get(ab65T, 5, "adjustment")
	if len(ab65Adjustment.Ab65Bits) != 8 {
		ab65T.Fatal("missing adjustment bits")
	}
	for ab65Index, ab65Raw := range []uint8{0, 3, 0xfc} {
		ab65Reserved := uint8(0)
		for _, ab65Bit := range ab65Adjustment.Ab65Bits {
			if ab65Bit.Ab65Reserved {
				ab65Reserved |= ab65Raw & (1 << ab65Bit.Ab65Bit)
				continue
			}
			ab65Want := []string{"released", "blocked", "released"}[ab65Index]
			if ab65Bit.Ab65Values[(ab65Raw>>ab65Bit.Ab65Bit)&1] != ab65Want {
				ab65T.Fatal("blocked polarity")
			}
		}
		if ab65Reserved != []uint8{0, 0, 0xfc}[ab65Index] {
			ab65T.Fatal("reserved bits erased")
		}
	}
	for ab65Index, ab65Bit := range ab65Adjustment.Ab65Bits {
		if ab65Bit.Ab65Bit != uint8(ab65Index) || ab65Bit.Ab65Reserved != (ab65Index >= 2) || ab65Bit.Ab65Expected != nil {
			ab65T.Fatal("adjustment bit identity")
		}
	}
	if *ab65Adjustment.Ab65Bits[0].Ab65Meaning != "CH" || *ab65Adjustment.Ab65Bits[1].Ab65Meaning != "DHW" {
		ab65T.Fatal("adjustment field identity")
	}
	ab65Status := ab65Get(ab65T, 6, "status")
	if len(ab65Status.Ab65Bits) != 8 || !ab65reflect.DeepEqual(ab65Status.Ab65Bits[0].Ab65Values, map[uint8]string{0: "off", 1: "on"}) {
		ab65T.Fatal("OP06 bit0")
	}
	for ab65Index, ab65Bit := range ab65Status.Ab65Bits[1:] {
		if ab65Bit.Ab65Bit != uint8(ab65Index+1) || !ab65Bit.Ab65Reserved || ab65Bit.Ab65Expected == nil || *ab65Bit.Ab65Expected != 0 {
			ab65T.Fatal("OP06 reserved zero")
		}
	}
}

func TestAb65DisableModesAndUnresolved(ab65T *ab65testing.T) {
	ab65Disable := ab65Get(ab65T, 0, "disable")
	if len(ab65Disable.Ab65Bits) != 8 || ab65Disable.Ab65Values != nil {
		ab65T.Fatal("invented mode codes")
	}
	for ab65Index, ab65Bit := range ab65Disable.Ab65Bits {
		if ab65Bit.Ab65Bit != uint8(ab65Index) {
			ab65T.Fatal("bit identity")
		}
		if ab65Index < 3 && (!ab65reflect.DeepEqual(ab65Bit.Ab65ActiveModes, []string{"automatic"}) || !ab65reflect.DeepEqual(ab65Bit.Ab65IgnoredModes, []string{"forced-heating", "forced-DHW", "forced-off"})) {
			ab65T.Fatal("disable mode applicability")
		}
	}
	for _, ab65Mode := range []string{"automatic", "forced-heating", "forced-DHW", "forced-off"} {
		ab65Active := uint8(0)
		for _, ab65Bit := range ab65Disable.Ab65Bits[:3] {
			for _, ab65Allowed := range ab65Bit.Ab65ActiveModes {
				if ab65Allowed == ab65Mode {
					ab65Active |= uint8(7) & (1 << ab65Bit.Ab65Bit)
				}
			}
		}
		ab65Want := uint8(0)
		if ab65Mode == "automatic" {
			ab65Want = 7
		}
		if ab65Active != ab65Want {
			ab65T.Fatal("07 mode example", ab65Mode)
		}
	}
	for _, ab65Bit := range ab65Disable.Ab65Bits[3:5] {
		if ab65Bit.Ab65Grade != "Unknown" || ab65Bit.Ab65Meaning == nil || *ab65Bit.Ab65Meaning != "timer assignment unresolved" || ab65Bit.Ab65Values != nil {
			ab65T.Fatal("settled competing timer assignments")
		}
	}
	if ab65Disable.Ab65Bits[5].Ab65Grade != "Unknown" || ab65Disable.Ab65Bits[5].Ab65Meaning != nil {
		ab65T.Fatal("omitted meaning invented")
	}
	if *ab65Disable.Ab65Bits[6].Ab65Meaning != "relates to subsequent flame status" {
		ab65T.Fatal("flame relationship")
	}
	ab65Legionella := ab65Disable.Ab65Bits[7]
	if *ab65Legionella.Ab65Meaning != "starts legionella protection" || ab65Legionella.Ab65PersistsUntilStatus == nil || !*ab65Legionella.Ab65PersistsUntilStatus || ab65Legionella.Ab65CommandsTermination == nil || *ab65Legionella.Ab65CommandsTermination || ab65Legionella.Ab65ActiveModes != nil || ab65Legionella.Ab65IgnoredModes != nil {
		ab65T.Fatal("bit7 independent persistence")
	}
	for _, ab65Raw := range []uint8{7, 0x80, 0x87, 0x18} {
		if (ab65Raw&0x80 != 0) != (ab65Raw == 0x80 || ab65Raw == 0x87) {
			ab65T.Fatal("bit7 conflated")
		}
	}
}

func ab65Mutate(ab65Value ab65reflect.Value) {
	switch ab65Value.Kind() {
	case ab65reflect.Pointer:
		if !ab65Value.IsNil() {
			ab65Mutate(ab65Value.Elem())
		}
	case ab65reflect.Struct:
		for ab65Index := 0; ab65Index < ab65Value.NumField(); ab65Index++ {
			ab65Mutate(ab65Value.Field(ab65Index))
		}
	case ab65reflect.Slice:
		for ab65Index := 0; ab65Index < ab65Value.Len(); ab65Index++ {
			ab65Mutate(ab65Value.Index(ab65Index))
		}
	case ab65reflect.Map:
		for _, ab65Key := range ab65Value.MapKeys() {
			ab65Value.SetMapIndex(ab65Key, ab65reflect.ValueOf("ab65_changed"))
		}
	case ab65reflect.String:
		ab65Value.SetString("ab65_changed")
	case ab65reflect.Bool:
		ab65Value.SetBool(!ab65Value.Bool())
	case ab65reflect.Uint8:
		ab65Value.SetUint(99)
	case ab65reflect.Int:
		ab65Value.SetInt(99)
	case ab65reflect.Float64:
		ab65Value.SetFloat(99)
	}
}

func TestAb65DeepCopy(ab65T *ab65testing.T) {
	ab65Before := Ab65B510Candidates()
	ab65Returned := Ab65B510Candidates()
	ab65Mutate(ab65reflect.ValueOf(&ab65Returned).Elem())
	if !ab65reflect.DeepEqual(Ab65B510Candidates(), ab65Before) {
		ab65T.Fatal("shared nested accessor backing")
	}
	for _, ab65Entry := range ab65Before.Ab65Entries {
		ab65Got, ab65OK := Ab65Lookup(ab65Entry.Ab65Key)
		if !ab65OK {
			ab65T.Fatal("lookup")
		}
		ab65Mutate(ab65reflect.ValueOf(&ab65Got).Elem())
		ab65Fresh, ab65OK := Ab65Lookup(ab65Entry.Ab65Key)
		if !ab65OK || !ab65reflect.DeepEqual(ab65Fresh, ab65Entry) {
			ab65T.Fatal("shared nested lookup backing")
		}
	}
}
