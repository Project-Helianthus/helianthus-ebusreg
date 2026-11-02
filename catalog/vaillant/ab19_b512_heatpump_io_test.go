package vaillant

import (
	ab19bytes "bytes"
	ab19json "encoding/json"
	ab19reflect "reflect"
	ab19testing "testing"
)

type ab19Row struct {
	Ab19ID          string  `json:"id"`
	Ab19Key         Ab19Key `json:"key"`
	Ab19Position    int     `json:"position"`
	Ab19Width       int     `json:"width_bits"`
	Ab19Encoding    string  `json:"encoding"`
	Ab19Scale       float64 `json:"scale"`
	Ab19Replacement *struct {
		Ab19Raw         uint16 `json:"raw_word"`
		Ab19Selection   string `json:"selection"`
		Ab19Meaning     string `json:"meaning"`
		Ab19Unavailable bool   `json:"unavailable_assertion"`
	} `json:"replacement"`
	Ab19Bits []struct {
		Ab19Bit    uint8  `json:"bit"`
		Ab19Label  string `json:"label"`
		Ab19Status string `json:"status"`
	} `json:"bits"`
	Ab19Values map[string]string `json:"values"`
}

func ab19Pointer[ab19Type any](ab19Value ab19Type) *ab19Type { return &ab19Value }
func ab19Equal(ab19T *ab19testing.T, ab19Got, ab19Want any) {
	ab19T.Helper()
	if !ab19reflect.DeepEqual(ab19Got, ab19Want) {
		ab19T.Fatalf("got %#v; want %#v", ab19Got, ab19Want)
	}
}
func ab19Get(ab19T *ab19testing.T, ab19Op uint8, ab19Direction, ab19Field string) (ab19Row, ab19json.RawMessage) {
	ab19T.Helper()
	ab19Raw, ab19Found, ab19Err := Ab19Lookup(Ab19Key{Ab19Family: "vaillant", Ab19Variant: "B512-heatpump-io", Ab19Primary: 181, Ab19Secondary: ab19Pointer(uint8(18)), Ab19Operation: ab19Pointer(ab19Op), Ab19Direction: ab19Direction, Ab19Field: "ab19_" + ab19Field, Ab19Alternative: "none", Ab19InstanceScope: "none"})
	if ab19Err != nil || !ab19Found {
		ab19T.Fatalf("lookup %d/%s/%s: %v %v", ab19Op, ab19Direction, ab19Field, ab19Found, ab19Err)
	}
	var ab19Entry ab19Row
	if ab19Err := ab19json.Unmarshal(ab19Raw, &ab19Entry); ab19Err != nil {
		ab19T.Fatal(ab19Err)
	}
	return ab19Entry, ab19Raw
}
func ab19Object(ab19T *ab19testing.T, ab19Raw ab19json.RawMessage) map[string]any {
	ab19T.Helper()
	var ab19Fields map[string]any
	if ab19Err := ab19json.Unmarshal(ab19Raw, &ab19Fields); ab19Err != nil {
		ab19T.Fatal(ab19Err)
	}
	return ab19Fields
}
func ab19Marshal(ab19T *ab19testing.T, ab19Value any) string {
	ab19T.Helper()
	ab19Raw, ab19Err := ab19json.Marshal(ab19Value)
	if ab19Err != nil {
		ab19T.Fatal(ab19Err)
	}
	return string(ab19Raw)
}

func TestAb19CatalogContract(ab19T *ab19testing.T) {
	ab19Raw, ab19Err := Ab19Entries()
	if ab19Err != nil {
		ab19T.Fatal(ab19Err)
	}
	var ab19Rows []ab19json.RawMessage
	if ab19Err := ab19json.Unmarshal(ab19Raw, &ab19Rows); ab19Err != nil {
		ab19T.Fatal(ab19Err)
	}
	ab19Equal(ab19T, len(ab19Rows), 34)
	for _, ab19Raw := range ab19Rows {
		var ab19Entry ab19Row
		if ab19Err := ab19json.Unmarshal(ab19Raw, &ab19Entry); ab19Err != nil {
			ab19T.Fatal(ab19Err)
		}
		ab19LookupRaw, ab19Found, ab19Err := Ab19Lookup(ab19Entry.Ab19Key)
		if ab19Err != nil || !ab19Found {
			ab19T.Fatalf("whole-row lookup: %v %v", ab19Found, ab19Err)
		}
		ab19Equal(ab19T, ab19Object(ab19T, ab19LookupRaw), ab19Object(ab19T, ab19Raw))
		ab19Fields := ab19Object(ab19T, ab19Raw)
		ab19Equal(ab19T, ab19Fields["grade"], "Hypothesis")
		ab19Equal(ab19T, ab19Fields["metadata_presence"], "Proven")
		ab19Equal(ab19T, ab19Fields["native_contract"], "https://github.com/Project-Helianthus/helianthus-docs-ebus/blob/main/protocols/ebus-services/ebus-overview.md")
		ab19Equal(ab19T, ab19Fields["unknown"], map[string]any{"native_support": "Unknown", "physical_effects": "Unknown", "persistence": "Unknown", "delivery": "Unknown", "physical_reply_form": "Unknown", "qualification": "Unknown"})
		ab19Equal(ab19T, ab19Fields["raw"], map[string]any{"retained": true, "short": "retain_raw", "extra": "retain_raw", "physical_reply_selection": false})
	}
}

func TestAb19IdentityAxes(ab19T *ab19testing.T) {
	ab19Base, ab19Raw := ab19Get(ab19T, 15, "request", "outputs")
	ab19Changes := []struct {
		ab19Name   string
		ab19Change func(*Ab19Key)
	}{
		{"ab19Family", func(ab19Key *Ab19Key) { ab19Key.Ab19Family = "synthetic-family" }},
		{"ab19Variant", func(ab19Key *Ab19Key) { ab19Key.Ab19Variant = "synthetic-variant" }},
		{"ab19Primary", func(ab19Key *Ab19Key) { ab19Key.Ab19Primary = 1 }},
		{"ab19Secondary", func(ab19Key *Ab19Key) { ab19Key.Ab19Secondary = ab19Pointer(uint8(0)) }},
		{"ab19Operation", func(ab19Key *Ab19Key) { ab19Key.Ab19Operation = ab19Pointer(uint8(0)) }},
		{"ab19Direction", func(ab19Key *Ab19Key) { ab19Key.Ab19Direction = "positive" }},
		{"ab19Field", func(ab19Key *Ab19Key) { ab19Key.Ab19Field = "ab19_unlisted" }},
		{"ab19Alternative", func(ab19Key *Ab19Key) { ab19Key.Ab19Alternative = "ab19_alternative" }},
		{"ab19Group", func(ab19Key *Ab19Key) { ab19Key.Ab19Group = ab19Pointer(uint8(0)) }},
		{"ab19InstanceScope", func(ab19Key *Ab19Key) { ab19Key.Ab19InstanceScope = "dynamic" }},
		{"ab19Register", func(ab19Key *Ab19Key) { ab19Key.Ab19Register = ab19Pointer(uint16(0)) }},
	}
	for _, ab19Case := range ab19Changes {
		ab19T.Run(ab19Case.ab19Name, func(ab19T *ab19testing.T) {
			ab19Key := ab19Base.Ab19Key
			ab19Case.ab19Change(&ab19Key)
			if _, ab19Found, ab19Err := Ab19Lookup(ab19Key); ab19Found || ab19Err != nil {
				ab19T.Fatalf("mismatched axis admitted: %v %v", ab19Found, ab19Err)
			}
			ab19First := ab19Object(ab19T, ab19Raw)
			ab19Second := ab19Object(ab19T, ab19Raw)
			ab19Second["key"] = ab19Key
			ab19Second["id"] = "ab19_control"
			ab19Input := ab19Marshal(ab19T, []any{ab19First, ab19Second})
			for _, ab19Control := range []struct {
				ab19Key Ab19Key
				ab19ID  string
			}{{ab19Base.Ab19Key, ab19Base.Ab19ID}, {ab19Key, "ab19_control"}} {
				ab19Got, ab19Found, ab19Err := ab19Lookup(ab19Input, ab19Control.ab19Key)
				if !ab19Found || ab19Err != nil {
					ab19T.Fatalf("valid control: %v %v", ab19Found, ab19Err)
				}
				ab19Equal(ab19T, ab19Object(ab19T, ab19Got)["id"], ab19Control.ab19ID)
			}
		})
	}
	for _, ab19Axis := range []string{"secondary", "operation", "group", "register"} {
		ab19T.Run("ab19AbsentZero_"+ab19Axis, func(ab19T *ab19testing.T) {
			ab19First := ab19Object(ab19T, ab19Raw)
			ab19Second := ab19Object(ab19T, ab19Raw)
			ab19First["key"].(map[string]any)[ab19Axis] = nil
			ab19Second["key"].(map[string]any)[ab19Axis] = float64(0)
			ab19Second["id"] = "ab19_zero"
			ab19Input := ab19Marshal(ab19T, []any{ab19First, ab19Second})
			for _, ab19Control := range []map[string]any{ab19First, ab19Second} {
				var ab19Key Ab19Key
				if ab19Err := ab19json.Unmarshal([]byte(ab19Marshal(ab19T, ab19Control["key"])), &ab19Key); ab19Err != nil {
					ab19T.Fatal(ab19Err)
				}
				ab19Got, ab19Found, ab19Err := ab19Lookup(ab19Input, ab19Key)
				if !ab19Found || ab19Err != nil {
					ab19T.Fatalf("absent/zero: %v %v", ab19Found, ab19Err)
				}
				ab19Equal(ab19T, ab19Object(ab19T, ab19Got)["id"], ab19Control["id"])
			}
		})
	}
}

func TestAb19ValidationAndIsolation(ab19T *ab19testing.T) {
	ab19Base, ab19Raw := ab19Get(ab19T, 15, "request", "outputs")
	ab19Duplicate := ab19Object(ab19T, ab19Raw)
	ab19Duplicate["id"] = "ab19_duplicate"
	ab19Invalid := []string{"{", "null", "{}", "[]", "[null]", string(ab19Raw), "[" + string(ab19Raw) + "] true", ab19Marshal(ab19T, []any{ab19Object(ab19T, ab19Raw), ab19Duplicate})}
	for _, ab19Axis := range []string{"family", "variant", "primary", "secondary", "operation", "direction", "field", "alternative", "group", "instance_scope", "register"} {
		ab19Missing := ab19Object(ab19T, ab19Raw)
		delete(ab19Missing["key"].(map[string]any), ab19Axis)
		ab19Invalid = append(ab19Invalid, ab19Marshal(ab19T, []any{ab19Missing}))
	}
	for _, ab19Input := range ab19Invalid {
		if _, ab19Err := ab19Entries(ab19Input); ab19Err == nil {
			ab19T.Fatal("invalid catalog admitted by entries")
		}
		if _, ab19Found, ab19Err := ab19Lookup(ab19Input, ab19Base.Ab19Key); ab19Err == nil || ab19Found {
			ab19T.Fatal("invalid catalog admitted by lookup")
		}
	}
	ab19All, ab19Err := Ab19Entries()
	if ab19Err != nil {
		ab19T.Fatal(ab19Err)
	}
	ab19Original := append([]byte(nil), ab19All...)
	for ab19Index := range ab19All {
		ab19All[ab19Index] = 'x'
	}
	for ab19Index := range ab19Raw {
		ab19Raw[ab19Index] = 'x'
	}
	ab19Again, ab19Err := Ab19Entries()
	if ab19Err != nil {
		ab19T.Fatal(ab19Err)
	}
	ab19Equal(ab19T, ab19bytes.Equal(ab19Again, ab19Original), true)
	_, ab19Fresh := ab19Get(ab19T, 15, "request", "outputs")
	ab19LookupOriginal := append([]byte(nil), ab19Fresh...)
	for ab19Index := range ab19Fresh {
		ab19Fresh[ab19Index] = 'x'
	}
	_, ab19Fresh = ab19Get(ab19T, 15, "request", "outputs")
	ab19Equal(ab19T, ab19bytes.Equal(ab19Fresh, ab19LookupOriginal), true)
}

func TestAb19FieldCoordinates(ab19T *ab19testing.T) {
	for _, ab19Case := range []struct {
		ab19Op                   uint8
		ab19Direction, ab19Field string
		ab19Position, ab19Width  int
		ab19Encoding             string
	}{
		{15, "request", "outputs", 3, 8, "bits"}, {15, "request", "valve", 4, 8, "u8"},
		{15, "positive", "temperature_one", 3, 16, "signed-le16"}, {15, "positive", "temperature_two", 6, 16, "signed-le16"},
		{15, "positive", "input", 5, 8, "u8"}, {15, "positive", "pressure", 8, 8, "u8"}, {15, "positive", "backup", 9, 8, "u8"},
		{16, "request", "actuators_one", 3, 8, "bits"}, {16, "request", "actuators_two", 4, 8, "bits"}, {16, "request", "inverse_actuator", 5, 16, "opaque-16"},
		{16, "positive", "status", 3, 8, "bits"}, {16, "positive", "single_phase", 4, 8, "u8"}, {16, "positive", "three_phase", 5, 8, "u8"},
		{16, "positive", "adc_1", 6, 16, "raw-le16"}, {16, "positive", "adc_2", 8, 16, "raw-le16"}, {16, "positive", "adc_3", 10, 16, "raw-le16"},
		{16, "positive", "adc_4", 12, 16, "raw-le16"}, {16, "positive", "adc_5", 14, 16, "raw-le16"}, {16, "positive", "adc_6", 16, 16, "raw-le16"},
		{16, "positive", "board", 18, 8, "u8"},
		{17, "request", "temperature_1", 3, 16, "signed-le16"}, {17, "request", "temperature_2", 5, 16, "signed-le16"}, {17, "request", "temperature_3", 7, 16, "signed-le16"},
		{17, "request", "pressure", 9, 8, "u8"}, {17, "request", "flow", 10, 8, "signed8"}, {17, "request", "pump_status", 11, 8, "raw-u8"}, {17, "request", "compressor", 12, 8, "bits"},
		{17, "positive", "status", 3, 8, "bits"}, {17, "positive", "setpoint", 4, 8, "u8"},
		{18, "request", "emergency", 3, 8, "bits"}, {18, "request", "switch", 4, 8, "raw-u8"}, {18, "request", "co2_1", 5, 16, "raw-le16"}, {18, "request", "co2_2", 7, 16, "raw-le16"},
		{18, "positive", "led", 3, 8, "bits"},
	} {
		ab19T.Run("ab19_"+ab19Case.ab19Direction+"_"+ab19Case.ab19Field+"_"+string(rune('A'+ab19Case.ab19Op)), func(ab19T *ab19testing.T) {
			ab19Row, ab19Raw := ab19Get(ab19T, ab19Case.ab19Op, ab19Case.ab19Direction, ab19Case.ab19Field)
			ab19Equal(ab19T, ab19Row.Ab19Position, ab19Case.ab19Position)
			ab19Equal(ab19T, ab19Row.Ab19Width, ab19Case.ab19Width)
			ab19Equal(ab19T, ab19Row.Ab19Encoding, ab19Case.ab19Encoding)
			if ab19Case.ab19Op != 15 || ab19Case.ab19Field != "pressure" {
				if _, ab19Present := ab19Object(ab19T, ab19Raw)["default"]; ab19Present {
					ab19T.Fatal("synthesized default")
				}
			}
		})
	}
}

func TestAb19Bits(ab19T *ab19testing.T) {
	for _, ab19Case := range []struct {
		ab19Op                   uint8
		ab19Direction, ab19Field string
		ab19Labels               []string
		ab19Unmodeled            bool
	}{
		{15, "request", "outputs", []string{"output one", "output two", "reserved", "reserved", "reserved", "reserved", "reserved", "reserved"}, false},
		{16, "request", "actuators_one", []string{"cooling valve", "passive mixer one", "passive mixer two", "cooling active", "well pump", "backup safety", "backup one", "backup two"}, false},
		{16, "request", "actuators_two", []string{"backup three", "external valve", "external backup", "reserved", "reserved", "reserved", "reserved", "reserved"}, false},
		{16, "positive", "status", []string{"anode channel one", "anode channel two", "smart-grid", "well", "safety limiter", "safety", "reserved", "reserved"}, false},
		{17, "request", "compressor", []string{"compressor-off", "reserved", "reserved", "reserved", "reserved", "reserved", "reserved", "unmodeled"}, true},
		{17, "positive", "status", []string{"compressor-off", "pump-relay", "reserved", "reserved", "reserved", "reserved", "reserved", "reserved"}, false},
		{18, "request", "emergency", []string{"emergency", "reserved", "reserved", "reserved", "reserved", "reserved", "reserved", "unmodeled"}, true},
		{18, "positive", "led", []string{"LED", "reserved", "reserved", "reserved", "reserved", "reserved", "reserved", "unmodeled"}, true},
	} {
		ab19Row, _ := ab19Get(ab19T, ab19Case.ab19Op, ab19Case.ab19Direction, ab19Case.ab19Field)
		ab19Equal(ab19T, len(ab19Row.Ab19Bits), 8)
		for ab19Index, ab19Label := range ab19Case.ab19Labels {
			ab19Bit := ab19Row.Ab19Bits[ab19Index]
			ab19Equal(ab19T, ab19Bit.Ab19Bit, uint8(ab19Index))
			ab19Equal(ab19T, ab19Bit.Ab19Label, ab19Label)
			ab19Status := "candidate"
			if ab19Label == "reserved" || ab19Label == "unmodeled" {
				ab19Status = ab19Label
			}
			ab19Equal(ab19T, ab19Bit.Ab19Status, ab19Status)
			for ab19Toggle := 0; ab19Toggle < 8; ab19Toggle++ {
				ab19Sample := uint8(1 << ab19Toggle)
				ab19Equal(ab19T, ab19Sample&(1<<ab19Bit.Ab19Bit) != 0, ab19Toggle == ab19Index)
			}
		}
		if ab19Case.ab19Unmodeled {
			ab19Equal(ab19T, ab19Row.Ab19Bits[7].Ab19Status, "unmodeled")
		}
	}
}

func ab19Evaluate(ab19Row ab19Row, ab19Raw []byte) (float64, bool, []byte) {
	ab19Retained := append([]byte(nil), ab19Raw...)
	if len(ab19Raw)*8 != ab19Row.Ab19Width {
		return 0, false, ab19Retained
	}
	ab19Word := uint16(ab19Raw[0])
	if len(ab19Raw) == 2 {
		ab19Word |= uint16(ab19Raw[1]) << 8
	}
	if ab19Row.Ab19Replacement != nil && ab19Row.Ab19Replacement.Ab19Selection == "before_scaling" && ab19Word == ab19Row.Ab19Replacement.Ab19Raw {
		return 0, true, ab19Retained
	}
	ab19Value := float64(ab19Word)
	if ab19Row.Ab19Encoding == "signed-le16" {
		ab19Value = float64(int16(ab19Word))
	}
	if ab19Row.Ab19Encoding == "signed8" {
		ab19Value = float64(int8(ab19Word))
	}
	if ab19Row.Ab19Scale != 0 {
		ab19Value *= ab19Row.Ab19Scale
	}
	return ab19Value, false, ab19Retained
}

func TestAb19NumericCandidates(ab19T *ab19testing.T) {
	for _, ab19Field := range []string{"temperature_one", "temperature_two"} {
		ab19Row, _ := ab19Get(ab19T, 15, "positive", ab19Field)
		ab19Equal(ab19T, ab19Row.Ab19Scale, 0.0625)
		if ab19Row.Ab19Replacement == nil {
			ab19T.Fatal("missing replacement")
		}
		ab19Equal(ab19T, ab19Row.Ab19Replacement.Ab19Selection, "before_scaling")
		ab19Equal(ab19T, ab19Row.Ab19Replacement.Ab19Meaning, "replacement")
		ab19Equal(ab19T, ab19Row.Ab19Replacement.Ab19Unavailable, false)
		ab19Value, ab19Replacement, ab19Raw := ab19Evaluate(ab19Row, []byte{0, 128})
		ab19Equal(ab19T, ab19Value, float64(0))
		ab19Equal(ab19T, ab19Replacement, true)
		ab19Equal(ab19T, ab19Raw, []byte{0, 128})
		ab19Value, ab19Replacement, ab19Raw = ab19Evaluate(ab19Row, []byte{240, 255})
		ab19Equal(ab19T, ab19Value, float64(-1))
		ab19Equal(ab19T, ab19Replacement, false)
		ab19Equal(ab19T, ab19Raw, []byte{240, 255})
	}
	ab19Pressure, ab19Raw := ab19Get(ab19T, 15, "positive", "pressure")
	ab19Fields := ab19Object(ab19T, ab19Raw)
	ab19Equal(ab19T, ab19Fields["default"], float64(255))
	ab19Equal(ab19T, ab19Fields["default_conflict"], true)
	ab19Equal(ab19T, ab19Fields["domain"], []any{float64(0), 25.5})
	ab19Equal(ab19T, ab19Fields["raw_domain"], []any{float64(0), float64(255)})
	ab19Value, ab19Replacement, ab19Bytes := ab19Evaluate(ab19Pressure, []byte{255})
	ab19Equal(ab19T, ab19Value, 25.5)
	ab19Equal(ab19T, ab19Replacement, false)
	ab19Equal(ab19T, ab19Bytes, []byte{255})
	for _, ab19Field := range []string{"temperature_1", "temperature_2", "temperature_3"} {
		ab19Row, _ := ab19Get(ab19T, 17, "request", ab19Field)
		ab19Value, ab19Replacement, _ := ab19Evaluate(ab19Row, []byte{240, 255})
		ab19Equal(ab19T, ab19Value, float64(-1))
		ab19Equal(ab19T, ab19Replacement, false)
		if ab19Row.Ab19Replacement != nil {
			ab19T.Fatal("synthesized replacement")
		}
	}
	ab19Flow, _ := ab19Get(ab19T, 17, "request", "flow")
	ab19Value, ab19Replacement, _ = ab19Evaluate(ab19Flow, []byte{255})
	ab19Equal(ab19T, ab19Value, float64(-10))
	ab19Equal(ab19T, ab19Replacement, false)
	for _, ab19Field := range []string{"co2_1", "co2_2"} {
		ab19Row, ab19Raw := ab19Get(ab19T, 18, "request", ab19Field)
		ab19Value, ab19Replacement, ab19Bytes = ab19Evaluate(ab19Row, []byte{255, 255})
		ab19Equal(ab19T, ab19Value, float64(65535))
		ab19Equal(ab19T, ab19Replacement, false)
		ab19Equal(ab19T, ab19Bytes, []byte{255, 255})
		ab19Equal(ab19T, ab19Object(ab19T, ab19Raw)["domain"], []any{float64(0), float64(65535)})
	}
}

func TestAb19OpaqueAndPhaseCandidates(ab19T *ab19testing.T) {
	for _, ab19Field := range []string{"adc_1", "adc_2", "adc_3", "adc_4", "adc_5", "adc_6"} {
		ab19Row, _ := ab19Get(ab19T, 16, "positive", ab19Field)
		for _, ab19Sample := range []struct {
			ab19Bytes []byte
			ab19Want  float64
		}{{[]byte{0, 0}, 0}, {[]byte{1, 0}, 1}, {[]byte{0, 1}, 256}, {[]byte{255, 255}, 65535}} {
			ab19Value, ab19Replacement, ab19Retained := ab19Evaluate(ab19Row, ab19Sample.ab19Bytes)
			ab19Equal(ab19T, ab19Value, ab19Sample.ab19Want)
			ab19Equal(ab19T, ab19Replacement, false)
			ab19Equal(ab19T, ab19Retained, ab19Sample.ab19Bytes)
		}
	}
	ab19Inverse, ab19Raw := ab19Get(ab19T, 16, "request", "inverse_actuator")
	ab19Equal(ab19T, ab19Object(ab19T, ab19Raw)["algorithm"], "Unknown")
	for _, ab19Sample := range [][]byte{{255, 255}, {0, 0}, {170, 85}, {1}, {1, 2, 3}} {
		_, ab19Replacement, ab19Retained := ab19Evaluate(ab19Inverse, ab19Sample)
		ab19Equal(ab19T, ab19Replacement, false)
		ab19Equal(ab19T, ab19Retained, ab19Sample)
	}
	ab19Board, ab19Raw := ab19Get(ab19T, 16, "positive", "board")
	ab19Equal(ab19T, ab19Board.Ab19Values, map[string]string{"0": "opaque", "1": "opaque", "2": "opaque", "3": "opaque", "11": "opaque", "12": "opaque"})
	if _, ab19Present := ab19Board.Ab19Values["13"]; ab19Present {
		ab19T.Fatal("unmapped board code inferred")
	}
	ab19Equal(ab19T, ab19Object(ab19T, ab19Raw)["unmapped"], "retain_raw")
	for _, ab19Sample := range []byte{12, 13} {
		ab19Value, _, ab19Retained := ab19Evaluate(ab19Board, []byte{ab19Sample})
		ab19Equal(ab19T, ab19Value, float64(ab19Sample))
		ab19Equal(ab19T, ab19Retained, []byte{ab19Sample})
	}
	_, ab19Raw = ab19Get(ab19T, 16, "positive", "single_phase")
	ab19Equal(ab19T, ab19Object(ab19T, ab19Raw)["relations"], []any{map[string]any{"raw": float64(50), "comparison": "approximately", "meaning": "powered"}, map[string]any{"raw": float64(0), "comparison": "exact", "meaning": "absent"}})
	ab19Equal(ab19T, ab19Object(ab19T, ab19Raw)["tolerance"], "Unknown")
	ab19Phase, _ := ab19Get(ab19T, 16, "positive", "three_phase")
	ab19Equal(ab19T, ab19Phase.Ab19Values, map[string]string{"100": "powered", "0": "all phases absent", "33": "one phase absent", "66": "two phases absent"})
	for _, ab19Field := range []string{"single_phase", "three_phase"} {
		ab19Row, _ := ab19Get(ab19T, 16, "positive", ab19Field)
		for _, ab19Sample := range []byte{0, 33, 50, 66, 100} {
			ab19Value, _, ab19Retained := ab19Evaluate(ab19Row, []byte{ab19Sample})
			ab19Equal(ab19T, ab19Value, float64(ab19Sample))
			ab19Equal(ab19T, ab19Retained, []byte{ab19Sample})
		}
	}
}

func TestAb19EnumsAndPositiveConflict(ab19T *ab19testing.T) {
	ab19Valve, _ := ab19Get(ab19T, 15, "request", "valve")
	ab19Equal(ab19T, ab19Valve.Ab19Values, map[string]string{"0": "off", "1": "CH", "2": "DHW"})
	ab19Input, _ := ab19Get(ab19T, 15, "positive", "input")
	ab19Equal(ab19T, ab19Input.Ab19Values, map[string]string{"0": "multipurpose", "1": "utility contact"})
	ab19Backup, _ := ab19Get(ab19T, 15, "positive", "backup")
	ab19Equal(ab19T, ab19Backup.Ab19Values, map[string]string{"0": "off", "1": "step one", "2": "step two", "3": "step three"})
	if ab19Backup.Ab19Replacement == nil {
		ab19T.Fatal("backup replacement missing")
	}
	ab19Equal(ab19T, ab19Backup.Ab19Replacement.Ab19Meaning, "replacement")
	ab19Equal(ab19T, ab19Backup.Ab19Replacement.Ab19Unavailable, false)
	_, ab19Replacement, ab19Retained := ab19Evaluate(ab19Backup, []byte{255})
	ab19Equal(ab19T, ab19Replacement, true)
	ab19Equal(ab19T, ab19Retained, []byte{255})
	for _, ab19Case := range []struct {
		ab19Op    uint8
		ab19Field string
	}{{17, "status"}, {17, "setpoint"}, {18, "led"}} {
		ab19Row, ab19Raw := ab19Get(ab19T, ab19Case.ab19Op, "positive", ab19Case.ab19Field)
		ab19Equal(ab19T, *ab19Row.Ab19Key.Ab19Operation, ab19Case.ab19Op)
		ab19Equal(ab19T, ab19Object(ab19T, ab19Raw)["declared_positive"], map[string]any{"primary": float64(181), "secondary": float64(18), "operation": float64(12), "selector_conflict": true, "decoder_selection": false, "acceptance": false})
		for _, ab19Structure := range [][]byte{{}, {0}, {0, 0, 0, 0, 0, 0, 0, 0, 0}} {
			ab19Equal(ab19T, ab19Object(ab19T, ab19Raw)["raw"], map[string]any{"retained": true, "short": "retain_raw", "extra": "retain_raw", "physical_reply_selection": false})
			_, _, ab19Retained := ab19Evaluate(ab19Row, ab19Structure)
			ab19Equal(ab19T, ab19bytes.Equal(ab19Retained, ab19Structure), true)
		}
	}
	_, ab19Raw := ab19Get(ab19T, 17, "positive", "setpoint")
	ab19Equal(ab19T, ab19Object(ab19T, ab19Raw)["domain"], []any{float64(0), float64(100)})
}

func TestAb19UnitsAndScales(ab19T *ab19testing.T) {
	for _, ab19Case := range []struct {
		ab19Op                             uint8
		ab19Direction, ab19Field, ab19Unit string
		ab19Scale                          float64
	}{
		{15, "positive", "temperature_one", "degrees Celsius", 0.0625},
		{15, "positive", "temperature_two", "degrees Celsius", 0.0625},
		{15, "positive", "pressure", "bar", 0.1},
		{16, "positive", "single_phase", "percent", 0},
		{16, "positive", "three_phase", "percent", 0},
		{17, "request", "temperature_1", "degrees Celsius", 0.0625},
		{17, "request", "temperature_2", "degrees Celsius", 0.0625},
		{17, "request", "temperature_3", "degrees Celsius", 0.0625},
		{17, "request", "pressure", "bar", 0.1},
		{17, "request", "flow", "l/h", 10},
		{17, "positive", "setpoint", "percent", 0},
		{18, "request", "co2_1", "ppm", 0},
		{18, "request", "co2_2", "ppm", 0},
	} {
		ab19Row, ab19Raw := ab19Get(ab19T, ab19Case.ab19Op, ab19Case.ab19Direction, ab19Case.ab19Field)
		ab19Equal(ab19T, ab19Row.Ab19Scale, ab19Case.ab19Scale)
		ab19Equal(ab19T, ab19Object(ab19T, ab19Raw)["unit"], ab19Case.ab19Unit)
	}
	ab19Pressure, _ := ab19Get(ab19T, 17, "request", "pressure")
	ab19Value, ab19Replacement, ab19Raw := ab19Evaluate(ab19Pressure, []byte{255})
	ab19Equal(ab19T, ab19Value, 25.5)
	ab19Equal(ab19T, ab19Replacement, false)
	ab19Equal(ab19T, ab19Raw, []byte{255})
}
