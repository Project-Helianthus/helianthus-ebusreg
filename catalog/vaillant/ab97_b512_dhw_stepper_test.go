package vaillant

import (
	ab97json "encoding/json"
	ab97reflect "reflect"
	ab97strconv "strconv"
	ab97testing "testing"
)

func ab97Byte(ab97Value uint8) *uint8 { return &ab97Value }
func ab97Key(ab97Operation *uint8, ab97Direction, ab97Field string) Ab97Key {
	return Ab97Key{Ab97Family: "vaillant", Ab97Variant: "B512", Ab97Primary: 181, Ab97Secondary: ab97Byte(18), Ab97Operation: ab97Operation, Ab97Direction: ab97Direction, Ab97Field: ab97Field}
}
func ab97Require(ab97T *ab97testing.T, ab97Operation uint8, ab97Direction, ab97Field string) Ab97Entry {
	ab97T.Helper()
	ab97Entry, ab97Found, ab97Err := Ab97Lookup(ab97Key(ab97Byte(ab97Operation), ab97Direction, ab97Field))
	if ab97Err != nil || !ab97Found {
		ab97T.Fatalf("lookup %d %s %s: %v %v", ab97Operation, ab97Direction, ab97Field, ab97Found, ab97Err)
	}
	return ab97Entry
}
func ab97Scalar(ab97T *ab97testing.T, ab97Entry Ab97Entry, ab97Raw []byte) float64 {
	ab97T.Helper()
	if ab97Entry.Ab97WidthBits == nil {
		ab97T.Fatal("missing width")
	}
	ab97Value := int(ab97Raw[0])
	if *ab97Entry.Ab97WidthBits == 16 {
		if ab97Entry.Ab97ByteOrder == nil || *ab97Entry.Ab97ByteOrder != "little-endian" {
			ab97T.Fatal("byte order")
		}
		ab97Value += int(ab97Raw[1]) << 8
		if ab97Entry.Ab97Signed != nil && *ab97Entry.Ab97Signed {
			ab97Value = int(int16(ab97Value))
		}
	} else if ab97Entry.Ab97Signed != nil && *ab97Entry.Ab97Signed {
		ab97Value = int(int8(ab97Value))
	}
	ab97Result := float64(ab97Value)
	if ab97Entry.Ab97Divisor != nil {
		ab97Result /= *ab97Entry.Ab97Divisor
	}
	return ab97Result
}
func TestAb97Scalars(ab97T *ab97testing.T) {
	ab97Raw := ab97Require(ab97T, 3, "positive", "temperature_raw")
	if ab97Scalar(ab97T, ab97Raw, []byte{0x34, 0x12}) != 4660 || ab97Raw.Ab97Divisor != nil || ab97Raw.Ab97Unit != nil || *ab97Raw.Ab97Position != 4 {
		ab97T.Fatal("OP03 raw must remain unscaled")
	}
	if ab97Require(ab97T, 3, "request", "selector").Ab97Values["0"] != "valve-temperature-input" {
		ab97T.Fatal("selector")
	}
	ab97Status := ab97Require(ab97T, 3, "positive", "status")
	if !ab97reflect.DeepEqual(ab97Status.Ab97Values, map[string]string{"0": "okay", "1": "disconnected", "2": "short-circuit"}) || *ab97Status.Ab97Position != 3 {
		ab97T.Fatal("status")
	}
	for _, ab97Case := range []struct {
		ab97Op                   uint8
		ab97Direction, ab97Field string
		ab97Position             uint8
		ab97Raw                  []byte
		ab97Want                 float64
		ab97Unit                 string
	}{
		{6, "request", "flow", 5, []byte{1}, 0.1, "l/min"},
		{6, "request", "temperature_6", 6, []byte{1}, 0.5, "Celsius"},
		{6, "request", "temperature_7", 7, []byte{1}, 0.5, "Celsius"},
		{6, "request", "temperature_8", 8, []byte{1}, 0.5, "Celsius"},
		{6, "positive", "dhw_setpoint", 4, []byte{1}, 0.5, "Celsius"},
		{6, "positive", "return", 5, []byte{1}, 0.5, "Celsius"},
		{8, "request", "inlet", 3, []byte{0xf0, 0xff}, -1, "Celsius"},
		{8, "positive", "dhw_setpoint", 3, []byte{0xf0, 0xff}, -1, "Celsius"},
	} {
		ab97Entry := ab97Require(ab97T, ab97Case.ab97Op, ab97Case.ab97Direction, ab97Case.ab97Field)
		if ab97Scalar(ab97T, ab97Entry, ab97Case.ab97Raw) != ab97Case.ab97Want || *ab97Entry.Ab97Position != ab97Case.ab97Position || ab97Entry.Ab97Unit == nil || *ab97Entry.Ab97Unit != ab97Case.ab97Unit || *ab97Entry.Ab97Default != 0 {
			ab97T.Fatalf("scalar %s", ab97Case.ab97Field)
		}
	}
	ab97Adjustment := ab97Require(ab97T, 6, "positive", "adjustment")
	if ab97Scalar(ab97T, ab97Adjustment, []byte{255}) != -1 || ab97Scalar(ab97T, ab97Adjustment, []byte{2}) != 2 || *ab97Adjustment.Ab97Default != -1 || *ab97Adjustment.Ab97Position != 6 || !ab97reflect.DeepEqual(ab97Adjustment.Ab97Values, map[string]string{"-1": "-1", "0": "0", "1": "1", "2": "2"}) {
		ab97T.Fatal("signed adjustment")
	}
	ab97Modulation := ab97Require(ab97T, 6, "request", "modulation")
	if ab97Scalar(ab97T, ab97Modulation, []byte{1}) != 0.5 || *ab97Modulation.Ab97Position != 4 || ab97Modulation.Ab97Unit != nil || *ab97Modulation.Ab97FormalUnit != "Celsius" || ab97Modulation.Ab97Unknown[0] != "modulation unit conflicts with formal Celsius unit" {
		ab97T.Fatal("unit conflict lost")
	}
	ab97Status = ab97Require(ab97T, 8, "request", "status")
	if *ab97Status.Ab97Position != 5 || *ab97Status.Ab97Default != 0 {
		ab97T.Fatal("OP08 status position/default")
	}
	for ab97Code, ab97Want := range map[uint8]string{0: "valid", 0x55: "short", 0xaa: "interruption-or-disconnected", 0x54: ""} {
		if ab97Status.Ab97Values[ab97strconv.Itoa(int(ab97Code))] != ab97Want {
			ab97T.Fatal("OP08 status mapping")
		}
	}
}

func TestAb97Flags(ab97T *ab97testing.T) {
	for _, ab97Case := range []struct {
		ab97Op        uint8
		ab97Direction string
		ab97Position  uint8
		ab97Fields    []string
		ab97Named     int
	}{
		{6, "request", 3, []string{"task-active", "heat-demand", "loading-active", "loading-pump-fault", "condition_4", "condition_5", "condition_6", "condition_7"}, 4},
		{6, "request", 9, []string{"inlet-condition_0", "inlet-condition_1", "reserved_2", "reserved_3", "reserved_4", "reserved_5", "reserved_6", "reserved_7"}, 0},
		{6, "positive", 3, []string{"loading-disabled", "internal-test-active", "ready", "frost-protection-active", "legionella-demand", "reset", "reserved_6", "reserved_7"}, 6},
		{8, "request", 6, []string{"stepper-release", "reserved_1", "reserved_2", "reserved_3", "reserved_4", "reserved_5", "reserved_6", "reserved_7"}, 1},
		{8, "positive", 5, []string{"stepper-release", "reserved_1", "reserved_2", "reserved_3", "reserved_4", "reserved_5", "reserved_6", "reserved_7"}, 1},
		{9, "request", 3, []string{"stepper-release", "reserved_1", "reserved_2", "reserved_3", "reserved_4", "reserved_5", "reserved_6", "reserved_7"}, 1},
		{9, "positive", 3, []string{"stepper-release", "reserved_1", "reserved_2", "reserved_3", "reserved_4", "reserved_5", "reserved_6", "reserved_7"}, 1},
	} {
		for ab97Bit, ab97Field := range ab97Case.ab97Fields {
			ab97Entry := ab97Require(ab97T, ab97Case.ab97Op, ab97Case.ab97Direction, ab97Field)
			if *ab97Entry.Ab97Position != ab97Case.ab97Position || *ab97Entry.Ab97WidthBits != 1 || ab97Entry.Ab97Bit == nil {
				ab97T.Fatal("flag coordinate")
			}
			for ab97Toggle := 0; ab97Toggle < 8; ab97Toggle++ {
				if ((uint8(1)<<ab97Toggle)&(uint8(1)<<*ab97Entry.Ab97Bit) != 0) != (ab97Toggle == ab97Bit) {
					ab97T.Fatalf("flag swap: %s", ab97Field)
				}
			}
			if ab97Bit < ab97Case.ab97Named {
				if ab97Entry.Ab97Label == nil || *ab97Entry.Ab97Label != ab97Field {
					ab97T.Fatal("label")
				}
			} else if ab97Entry.Ab97Label != nil {
				ab97T.Fatal("invented condition label")
			}
			ab97Reserved := len(ab97Field) >= 8 && ab97Field[:8] == "reserved"
			if ab97Entry.Ab97Reserved != ab97Reserved {
				ab97T.Fatal("reserved bit lost")
			}
			if ab97Reserved {
				if ab97Entry.Ab97Default == nil || *ab97Entry.Ab97Default != 0 {
					ab97T.Fatal("reserved default")
				}
			} else if ab97Entry.Ab97BooleanDefault == nil || *ab97Entry.Ab97BooleanDefault {
				ab97T.Fatal("Boolean default")
			}
			if ab97Field == "stepper-release" && !ab97reflect.DeepEqual(ab97Entry.Ab97Unknown, []string{"release does not prove movement", "physical adaptation", "support", "actual movement"}) {
				ab97T.Fatal("release implies movement")
			}
		}
	}
}

func TestAb97ExactKeys(ab97T *ab97testing.T) {
	ab97Base := ab97Key(ab97Byte(3), "positive", "temperature_raw")
	for ab97Name, ab97Change := range map[string]func(*Ab97Key){
		"family":             func(ab97Key *Ab97Key) { ab97Key.Ab97Family = "other" },
		"variant":            func(ab97Key *Ab97Key) { ab97Key.Ab97Variant = "other" },
		"primary":            func(ab97Key *Ab97Key) { ab97Key.Ab97Primary = 0 },
		"secondary":          func(ab97Key *Ab97Key) { ab97Key.Ab97Secondary = ab97Byte(0) },
		"secondary_absent":   func(ab97Key *Ab97Key) { ab97Key.Ab97Secondary = nil },
		"operation":          func(ab97Key *Ab97Key) { ab97Key.Ab97Operation = ab97Byte(0) },
		"operation_absent":   func(ab97Key *Ab97Key) { ab97Key.Ab97Operation = nil },
		"direction":          func(ab97Key *Ab97Key) { ab97Key.Ab97Direction = "request" },
		"field":              func(ab97Key *Ab97Key) { ab97Key.Ab97Field = "other" },
		"group":              func(ab97Key *Ab97Key) { ab97Key.Ab97Group = ab97Byte(0) },
		"requested_instance": func(ab97Key *Ab97Key) { ab97Key.Ab97RequestedInstance = ab97Byte(0) },
		"register":           func(ab97Key *Ab97Key) { ab97Value := uint16(0); ab97Key.Ab97Register = &ab97Value },
		"command_role":       func(ab97Key *Ab97Key) { ab97Value := "read"; ab97Key.Ab97CommandRole = &ab97Value },
	} {
		ab97T.Run("ab97_"+ab97Name, func(ab97T *ab97testing.T) {
			ab97Changed := ab97Base
			ab97Change(&ab97Changed)
			ab97Entry, ab97Found, ab97Err := Ab97Lookup(ab97Changed)
			if ab97Err != nil || ab97Found || !ab97reflect.DeepEqual(ab97Entry, Ab97Entry{}) {
				ab97T.Fatalf("non-exact lookup: %v %v", ab97Found, ab97Err)
			}
		})
	}
	for _, ab97Direction := range []string{"request", "positive"} {
		ab97Entry := ab97Require(ab97T, 10, ab97Direction, "declaration")
		if ab97Entry.Ab97Position != nil || ab97Entry.Ab97WidthBits != nil || ab97Entry.Ab97Divisor != nil {
			ab97T.Fatal("invented OP0A body")
		}
		ab97Key := ab97Key(nil, ab97Direction, "declaration")
		ab97Key.Ab97Variant = "B51214"
		ab97Key.Ab97Secondary = ab97Byte(245)
		ab97Entry, ab97Found, ab97Err := Ab97Lookup(ab97Key)
		if ab97Err != nil || !ab97Found || ab97Entry.Ab97Key.Ab97Operation != nil || ab97Entry.Ab97WidthBits != nil {
			ab97T.Fatal("partial B51214 declaration")
		}
		for _, ab97Op := range []uint8{0, 10, 0x14} {
			ab97Key.Ab97Operation = ab97Byte(ab97Op)
			if _, ab97Found, ab97Err := Ab97Lookup(ab97Key); ab97Found || ab97Err != nil {
				ab97T.Fatal("name encoded as operation")
			}
		}
	}
}

func ab97Mutate(ab97Value ab97reflect.Value) {
	if !ab97Value.IsValid() {
		return
	}
	switch ab97Value.Kind() {
	case ab97reflect.Pointer:
		if !ab97Value.IsNil() {
			ab97Mutate(ab97Value.Elem())
		}
	case ab97reflect.Struct:
		for ab97Index := 0; ab97Index < ab97Value.NumField(); ab97Index++ {
			ab97Mutate(ab97Value.Field(ab97Index))
		}
	case ab97reflect.Slice:
		for ab97Index := 0; ab97Index < ab97Value.Len(); ab97Index++ {
			ab97Mutate(ab97Value.Index(ab97Index))
		}
	case ab97reflect.Map:
		if !ab97Value.IsNil() {
			ab97Value.SetMapIndex(ab97reflect.ValueOf("ab97_changed"), ab97reflect.ValueOf("ab97_changed"))
		}
	case ab97reflect.String:
		ab97Value.SetString("ab97_changed")
	case ab97reflect.Uint8, ab97reflect.Uint16:
		ab97Value.SetUint(0)
	case ab97reflect.Int:
		ab97Value.SetInt(99)
	case ab97reflect.Float64:
		ab97Value.SetFloat(99)
	case ab97reflect.Bool:
		ab97Value.SetBool(!ab97Value.Bool())
	}
}
func TestAb97IndependentStorage(ab97T *ab97testing.T) {
	ab97Before, ab97Err := Ab97Entries()
	if ab97Err != nil {
		ab97T.Fatal(ab97Err)
	}
	if len(ab97Before) != 74 {
		ab97T.Fatalf("finite entries: %d", len(ab97Before))
	}
	ab97Changed, ab97Err := Ab97Entries()
	if ab97Err != nil {
		ab97T.Fatal(ab97Err)
	}
	ab97Mutate(ab97reflect.ValueOf(&ab97Changed).Elem())
	for _, ab97Expected := range ab97Before {
		if ab97Expected.Ab97Position != nil && (ab97Expected.Ab97CoordinateBasis == nil || *ab97Expected.Ab97CoordinateBasis != "logical; includes three-position B5/12/operation prefix") {
			ab97T.Fatal("logical coordinate basis lost")
		}
		if ab97Expected.Ab97Position == nil && ab97Expected.Ab97CoordinateBasis != nil {
			ab97T.Fatal("invented coordinate basis")
		}
		if ab97Expected.Ab97Grade != "Hypothesis" || ab97Expected.Ab97Qualification != "Unknown" || ab97Expected.Ab97RawPolicy != "preserve native bytes, reserved bits and independent last-known-good fields" {
			ab97T.Fatal("grade or raw boundary")
		}
		ab97Found, ab97OK, ab97Err := Ab97Lookup(ab97Expected.Ab97Key)
		if ab97Err != nil || !ab97OK || !ab97reflect.DeepEqual(ab97Found, ab97Expected) {
			ab97T.Fatal("accessor aliases storage")
		}
		ab97Mutate(ab97reflect.ValueOf(&ab97Found).Elem())
		ab97Again, ab97OK, ab97Err := Ab97Lookup(ab97Expected.Ab97Key)
		if ab97Err != nil || !ab97OK || !ab97reflect.DeepEqual(ab97Again, ab97Expected) {
			ab97T.Fatal("lookup aliases storage")
		}
	}
	ab97After, ab97Err := Ab97Entries()
	if ab97Err != nil || !ab97reflect.DeepEqual(ab97Before, ab97After) {
		ab97T.Fatal("nested mutation escaped")
	}
}
func TestAb97MalformedData(ab97T *ab97testing.T) {
	ab97Original := ab97Data
	defer func() { ab97Data = ab97Original }()
	ab97Rows, ab97Err := Ab97Entries()
	if ab97Err != nil {
		ab97T.Fatal(ab97Err)
	}
	ab97Duplicate, ab97Err := ab97json.Marshal([]Ab97Entry{ab97Rows[0], ab97Rows[0]})
	if ab97Err != nil {
		ab97T.Fatal(ab97Err)
	}
	ab97Cases := [][]byte{[]byte("{"), []byte("null"), []byte("[]"), []byte("[{}]"), append(append([]byte{}, ab97Original...), []byte(" null")...), ab97Duplicate}
	for _, ab97Change := range []func(*Ab97Entry){
		func(ab97Entry *Ab97Entry) { ab97Entry.Ab97Grade = "Proven" },
		func(ab97Entry *Ab97Entry) { ab97Entry.Ab97Qualification = "qualified" },
		func(ab97Entry *Ab97Entry) { ab97Value := float64(0); ab97Entry.Ab97Divisor = &ab97Value },
		func(ab97Entry *Ab97Entry) { ab97Entry.Ab97Bit = ab97Byte(8) },
	} {
		ab97Entry := ab97Rows[0]
		ab97Change(&ab97Entry)
		ab97Bytes, ab97Err := ab97json.Marshal([]Ab97Entry{ab97Entry})
		if ab97Err != nil {
			ab97T.Fatal(ab97Err)
		}
		ab97Cases = append(ab97Cases, ab97Bytes)
	}
	for ab97Index, ab97Invalid := range ab97Cases {
		ab97Data = ab97Invalid
		if ab97Entries, ab97Err := Ab97Entries(); ab97Err == nil || ab97Entries != nil {
			ab97T.Fatalf("invalid catalog %d accepted", ab97Index)
		}
		if ab97Entry, ab97Found, ab97Err := Ab97Lookup(ab97Rows[0].Ab97Key); ab97Err == nil || ab97Found || !ab97reflect.DeepEqual(ab97Entry, Ab97Entry{}) {
			ab97T.Fatal("invalid lookup succeeded")
		}
	}
}

func TestAb97B524Identity(ab97T *ab97testing.T) {
	ab97Original := ab97Data
	defer func() { ab97Data = ab97Original }()
	ab97Rows, ab97Err := Ab97Entries()
	if ab97Err != nil {
		ab97T.Fatal(ab97Err)
	}
	ab97Register := uint16(1)
	ab97Role := "synthetic-read"
	ab97Base := Ab97Key{Ab97Family: "synthetic-family", Ab97Variant: "synthetic-B524", Ab97Primary: 181, Ab97Secondary: ab97Byte(36), Ab97Operation: ab97Byte(1), Ab97Direction: "request", Ab97Field: "synthetic-field", Ab97Group: ab97Byte(1), Ab97RequestedInstance: ab97Byte(1), Ab97Register: &ab97Register, Ab97CommandRole: &ab97Role}
	ab97Row := ab97Rows[0]
	ab97Row.Ab97Key = ab97Base
	ab97Data, ab97Err = ab97json.Marshal([]Ab97Entry{ab97Row})
	if ab97Err != nil {
		ab97T.Fatal(ab97Err)
	}
	if _, ab97Found, ab97Err := Ab97Lookup(ab97Base); ab97Err != nil || !ab97Found {
		ab97T.Fatal("synthetic base missing")
	}
	for ab97Name, ab97Change := range map[string]func(*Ab97Key){
		"family":                    func(ab97Key *Ab97Key) { ab97Key.Ab97Family = "other" },
		"variant":                   func(ab97Key *Ab97Key) { ab97Key.Ab97Variant = "other" },
		"primary":                   func(ab97Key *Ab97Key) { ab97Key.Ab97Primary = 0 },
		"secondary":                 func(ab97Key *Ab97Key) { ab97Key.Ab97Secondary = ab97Byte(0) },
		"operation":                 func(ab97Key *Ab97Key) { ab97Key.Ab97Operation = ab97Byte(0) },
		"direction":                 func(ab97Key *Ab97Key) { ab97Key.Ab97Direction = "positive" },
		"field":                     func(ab97Key *Ab97Key) { ab97Key.Ab97Field = "other" },
		"group":                     func(ab97Key *Ab97Key) { ab97Key.Ab97Group = ab97Byte(0) },
		"group_absent":              func(ab97Key *Ab97Key) { ab97Key.Ab97Group = nil },
		"requested_instance":        func(ab97Key *Ab97Key) { ab97Key.Ab97RequestedInstance = ab97Byte(0) },
		"requested_instance_absent": func(ab97Key *Ab97Key) { ab97Key.Ab97RequestedInstance = nil },
		"register":                  func(ab97Key *Ab97Key) { ab97Value := uint16(0); ab97Key.Ab97Register = &ab97Value },
		"register_absent":           func(ab97Key *Ab97Key) { ab97Key.Ab97Register = nil },
		"command_role":              func(ab97Key *Ab97Key) { ab97Value := "synthetic-write"; ab97Key.Ab97CommandRole = &ab97Value },
		"command_role_absent":       func(ab97Key *Ab97Key) { ab97Key.Ab97CommandRole = nil },
	} {
		ab97T.Run("ab97_"+ab97Name, func(ab97T *ab97testing.T) {
			ab97Changed := ab97Base
			ab97Change(&ab97Changed)
			if _, ab97Found, ab97Err := Ab97Lookup(ab97Changed); ab97Err != nil || ab97Found {
				ab97T.Fatal("synthetic B524 identity conflated")
			}
		})
	}
	ab97Changed, ab97Found, ab97Err := Ab97Lookup(ab97Base)
	if ab97Err != nil || !ab97Found {
		ab97T.Fatal("synthetic lookup")
	}
	*ab97Changed.Ab97Key.Ab97Register = 0
	*ab97Changed.Ab97Key.Ab97Group = 0
	*ab97Changed.Ab97Key.Ab97RequestedInstance = 0
	*ab97Changed.Ab97Key.Ab97CommandRole = "changed"
	ab97Again, ab97Found, ab97Err := Ab97Lookup(ab97Base)
	if ab97Err != nil || !ab97Found || !ab97reflect.DeepEqual(ab97Again, ab97Row) {
		ab97T.Fatal("synthetic B524 axes alias storage")
	}
}
