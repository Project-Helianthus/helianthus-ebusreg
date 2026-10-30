package vaillant

import (
	ab50Bytes "bytes"
	ab50Binary "encoding/binary"
	ab50JSON "encoding/json"
	ab50Reflect "reflect"
	ab50Testing "testing"
)

func ab50Ptr(ab50Value uint16) *uint16 { return &ab50Value }

func ab50TestKey(ab50Operation uint16, ab50Field string) Ab50Key {
	return Ab50Key{Ab50Family: "Vaillant", Ab50Variant: "class-field-candidate", Ab50Primary: 181, Ab50Secondary: ab50Ptr(17), Ab50Operation: ab50Ptr(ab50Operation), Ab50Direction: "Unknown", Ab50Field: ab50Field, Ab50Alternative: "none", Ab50InstanceScope: "none"}
}

func ab50TestRow(ab50T *ab50Testing.T, ab50Key Ab50Key) map[string]ab50JSON.RawMessage {
	ab50T.Helper()
	ab50Raw, ab50Found, ab50Error := Ab50Lookup(ab50Key)
	if ab50Error != nil || !ab50Found {
		ab50T.Fatalf("lookup: %v %v", ab50Found, ab50Error)
	}
	var ab50Row map[string]ab50JSON.RawMessage
	if ab50Error = ab50JSON.Unmarshal(ab50Raw, &ab50Row); ab50Error != nil {
		ab50T.Fatal(ab50Error)
	}
	return ab50Row
}

func ab50Assert(ab50T *ab50Testing.T, ab50Row map[string]ab50JSON.RawMessage, ab50Name, ab50Expected string) {
	ab50T.Helper()
	var ab50ActualValue, ab50ExpectedValue any
	if ab50JSON.Unmarshal(ab50Row[ab50Name], &ab50ActualValue) != nil || ab50JSON.Unmarshal([]byte(ab50Expected), &ab50ExpectedValue) != nil || !ab50Reflect.DeepEqual(ab50ActualValue, ab50ExpectedValue) {
		ab50T.Fatalf("%s: got %s want %s", ab50Name, ab50Row[ab50Name], ab50Expected)
	}
}

func TestAb50FiniteDeclarations(ab50T *ab50Testing.T) {
	ab50Expected := [][]string{
		{"flags-a", "flags-b", "flags-c", "modulation"},
		{"handshakes", "forced-mode-guards"},
		{"status", "DHW-enabled", "tapping-enabled"},
		{"basic-state", "flags-a", "flags-b", "message"},
		{"status", "temperature-a", "temperature-b"},
		{"setpoint-a", "setpoint-b", "setpoint-c", "setpoint-d", "status-a", "status-b"},
		{"CH-limit", "DHW-limit", "pressure-limit", "cooling-limit", "status"},
		{"power", "energy", "pump-flags", "status-flags"},
		{"volume", "pump", "status"},
		{"yield", "power", "activity", "display", "day-night"},
		{"error", "loading"},
		{"fan", "radiator", "bypass", "error", "maintenance", "assistant"},
		{"operation", "switch-on-flags"},
		{"status", "configuration"},
	}
	ab50EntriesJSON, ab50Error := Ab50Entries()
	if ab50Error != nil {
		ab50T.Fatal(ab50Error)
	}
	var ab50Rows []ab50JSON.RawMessage
	if ab50JSON.Unmarshal(ab50EntriesJSON, &ab50Rows) != nil || len(ab50Rows) != 53 {
		ab50T.Fatal("finite array must contain 53 declarations")
	}
	for ab50Operation, ab50Fields := range ab50Expected {
		for _, ab50Field := range ab50Fields {
			ab50Row := ab50TestRow(ab50T, ab50TestKey(uint16(ab50Operation), ab50Field))
			ab50Assert(ab50T, ab50Row, "grade", `"Hypothesis"`)
			ab50Assert(ab50T, ab50Row, "physical_qualification", `"Unknown"`)
			ab50Assert(ab50T, ab50Row, "offset", "null")
			ab50Assert(ab50T, ab50Row, "default", "null")
			ab50Assert(ab50T, ab50Row, "bit", "null")
			ab50Assert(ab50T, ab50Row, "labels", "{}")
			ab50Assert(ab50T, ab50Row, "raw_form", `"preserve original bytes and unlabelled values without clipping or sentinel inference"`)
			ab50Width, ab50Bounds, ab50Unit, ab50Order := "null", "null", "null", `"Unknown"`
			if ab50Field == "energy" || ab50Field == "yield" {
				ab50Width, ab50Unit, ab50Order = "16", `"kWh"`, `"little-endian"`
			}
			if (ab50Operation == 0 || ab50Operation == 3) && (ab50Field == "flags-a" || ab50Field == "flags-b" || ab50Field == "flags-c") || ab50Operation == 5 && (ab50Field == "status-a" || ab50Field == "status-b") {
				ab50Width = "8"
			}
			if (ab50Operation == 7 || ab50Operation == 9) && ab50Field == "power" || ab50Operation == 8 && ab50Field == "volume" {
				ab50Unit = `"percent"`
			}
			switch {
			case ab50Operation == 2 && ab50Field == "status":
				ab50Bounds = "[0,1]"
			case ab50Operation == 3 && ab50Field == "basic-state":
				ab50Bounds = "[0,6]"
			case ab50Operation == 3 && ab50Field == "message":
				ab50Bounds = "[0,3]"
			case ab50Operation == 7 && ab50Field == "energy":
				ab50Bounds = "[0,9999]"
			case ab50Operation == 9 && ab50Field == "yield":
				ab50Bounds = "[0,65535]"
			case ab50Operation == 12 && ab50Field == "operation":
				ab50Bounds = "[0,3]"
			}
			ab50Assert(ab50T, ab50Row, "width_bits", ab50Width)
			ab50Assert(ab50T, ab50Row, "bounds", ab50Bounds)
			ab50Assert(ab50T, ab50Row, "unit", ab50Unit)
			ab50Assert(ab50T, ab50Row, "byte_order", ab50Order)
			ab50Unresolved := `["absolute position","per-variant defaults","individual labels"]`
			if ab50Width == "null" {
				ab50Unresolved = `["absolute position","per-variant defaults","individual labels","width"]`
			}
			ab50Assert(ab50T, ab50Row, "unresolved", ab50Unresolved)
		}
	}
	ab50Assert(ab50T, ab50TestRow(ab50T, ab50TestKey(6, "status")), "unlabelled_values", "[4,9,16]")
	ab50Assert(ab50T, ab50TestRow(ab50T, ab50TestKey(2, "DHW-enabled")), "note", `"DHW-enabled covers tank-loading/warm-start while tapping remains enabled."`)
}

func TestAb50Applicability(ab50T *ab50Testing.T) {
	for _, ab50Case := range []struct {
		Ab50Operation                   uint16
		Ab50Field, Ab50Variant, Ab50Bit string
	}{{0, "second-modulation", "second-modulation-applicable", "null"}, {1, "warm-start", "warm-start-applicable", "5"}} {
		ab50Key := ab50TestKey(ab50Case.Ab50Operation, ab50Case.Ab50Field)
		ab50Key.Ab50Variant = ab50Case.Ab50Variant
		ab50Row := ab50TestRow(ab50T, ab50Key)
		ab50Assert(ab50T, ab50Row, "bit", ab50Case.Ab50Bit)
		for _, ab50Name := range []string{"offset", "width_bits", "default", "bounds"} {
			ab50Assert(ab50T, ab50Row, ab50Name, "null")
		}
		ab50Assert(ab50T, ab50Row, "grade", `"Hypothesis"`)
		ab50Assert(ab50T, ab50Row, "physical_qualification", `"Unknown"`)
		ab50Assert(ab50T, ab50Row, "labels", `{}`)
		ab50Assert(ab50T, ab50Row, "raw_form", `"preserve original bytes and unlabelled values without clipping or sentinel inference"`)
		for _, ab50Variant := range []string{"class-field-candidate", "warm-start-inapplicable", "none"} {
			ab50Key.Ab50Variant = ab50Variant
			if _, ab50Found, ab50Error := Ab50Lookup(ab50Key); ab50Error != nil || ab50Found {
				ab50T.Fatal("applicability widened", ab50Variant, ab50Error)
			}
		}
	}
}

func TestAb50SyntheticRawDomains(ab50T *ab50Testing.T) {
	for _, ab50Case := range []struct {
		Ab50Operation uint16
		Ab50Field     string
		Ab50Raw       []byte
		Ab50Value     uint16
		Ab50Within    bool
	}{
		{7, "energy", []byte{0x0f, 0x27}, 9999, true}, {7, "energy", []byte{0x10, 0x27}, 10000, false}, {9, "yield", []byte{0xff, 0xff}, 65535, true},
		{3, "message", []byte{3, 0}, 3, true}, {3, "message", []byte{4, 0}, 4, false}, {3, "basic-state", []byte{6, 0}, 6, true}, {3, "basic-state", []byte{7, 0}, 7, false},
	} {
		ab50Row := ab50TestRow(ab50T, ab50TestKey(ab50Case.Ab50Operation, ab50Case.Ab50Field))
		var ab50Bounds []uint16
		if ab50JSON.Unmarshal(ab50Row["bounds"], &ab50Bounds) != nil || len(ab50Bounds) != 2 {
			ab50T.Fatal("bounds")
		}
		ab50Value := ab50Binary.LittleEndian.Uint16(ab50Case.Ab50Raw)
		if ab50Value != ab50Case.Ab50Value || (ab50Value >= ab50Bounds[0] && ab50Value <= ab50Bounds[1]) != ab50Case.Ab50Within {
			ab50T.Fatal("independent raw domain result")
		}
		ab50Assert(ab50T, ab50Row, "raw_form", `"preserve original bytes and unlabelled values without clipping or sentinel inference"`)
	}
}

func TestAb50EveryIdentityAxis(ab50T *ab50Testing.T) {
	ab50Base := ab50TestKey(7, "energy")
	ab50Changes := map[string]func(*Ab50Key){
		"ab50Family":          func(ab50Key *Ab50Key) { ab50Key.Ab50Family = "ab50Other" },
		"ab50Variant":         func(ab50Key *Ab50Key) { ab50Key.Ab50Variant = "ab50Other" },
		"ab50Primary":         func(ab50Key *Ab50Key) { ab50Key.Ab50Primary = 0 },
		"ab50SecondaryAbsent": func(ab50Key *Ab50Key) { ab50Key.Ab50Secondary = nil },
		"ab50SecondaryZero":   func(ab50Key *Ab50Key) { ab50Key.Ab50Secondary = ab50Ptr(0) },
		"ab50OperationAbsent": func(ab50Key *Ab50Key) { ab50Key.Ab50Operation = nil },
		"ab50OperationZero":   func(ab50Key *Ab50Key) { ab50Key.Ab50Operation = ab50Ptr(0) },
		"ab50Direction":       func(ab50Key *Ab50Key) { ab50Key.Ab50Direction = "none" },
		"ab50Field":           func(ab50Key *Ab50Key) { ab50Key.Ab50Field = "yield" },
		"ab50Alternative":     func(ab50Key *Ab50Key) { ab50Key.Ab50Alternative = "ab50Other" },
		"ab50GroupZero":       func(ab50Key *Ab50Key) { ab50Key.Ab50Group = ab50Ptr(0) },
		"ab50GroupOther":      func(ab50Key *Ab50Key) { ab50Key.Ab50Group = ab50Ptr(1) },
		"ab50InstanceScope":   func(ab50Key *Ab50Key) { ab50Key.Ab50InstanceScope = "ab50Other" },
		"ab50RegisterZero":    func(ab50Key *Ab50Key) { ab50Key.Ab50Register = ab50Ptr(0) },
		"ab50RegisterOther":   func(ab50Key *Ab50Key) { ab50Key.Ab50Register = ab50Ptr(1) },
	}
	for ab50Name, ab50Change := range ab50Changes {
		ab50T.Run(ab50Name, func(ab50T *ab50Testing.T) {
			ab50Key := ab50Base
			ab50Change(&ab50Key)
			if _, ab50Found, ab50Error := Ab50Lookup(ab50Key); ab50Found || ab50Error != nil {
				ab50T.Fatal("inexact production match", ab50Error)
			}
			ab50Control := []struct {
				Ab50Key    Ab50Key `json:"key"`
				Ab50Marker string  `json:"ab50_marker"`
			}{{ab50Base, "ab50Base"}, {ab50Key, "ab50Changed"}}
			ab50Source, ab50Error := ab50JSON.Marshal(ab50Control)
			if ab50Error != nil {
				ab50T.Fatal(ab50Error)
			}
			for ab50Index, ab50Wanted := range []Ab50Key{ab50Base, ab50Key} {
				ab50Row, ab50Found, ab50Error := ab50Lookup(ab50Source, ab50Wanted)
				var ab50Actual map[string]ab50JSON.RawMessage
				if ab50Error != nil || !ab50Found || ab50JSON.Unmarshal(ab50Row, &ab50Actual) != nil {
					ab50T.Fatal("valid identity control", ab50Error)
				}
				ab50Expected, _ := ab50JSON.Marshal(ab50Control[ab50Index].Ab50Marker)
				if !ab50Bytes.Equal(ab50Actual["ab50_marker"], ab50Expected) {
					ab50T.Fatal("axis collapsed")
				}
			}
		})
	}
}

func TestAb50ValidationAndStorage(ab50T *ab50Testing.T) {
	ab50Key := ab50TestKey(7, "energy")
	ab50Original, ab50Found, ab50Error := Ab50Lookup(ab50Key)
	if ab50Error != nil || !ab50Found {
		ab50T.Fatal(ab50Error)
	}
	ab50EntriesJSON, ab50Error := Ab50Entries()
	if ab50Error != nil {
		ab50T.Fatal(ab50Error)
	}
	ab50Source := append([]byte(nil), ab50EntriesJSON...)
	ab50Copied, ab50Error := ab50Entries(ab50Source)
	if ab50Error != nil {
		ab50T.Fatal(ab50Error)
	}
	ab50Row, _, ab50Error := ab50Lookup(ab50Source, ab50Key)
	if ab50Error != nil {
		ab50T.Fatal(ab50Error)
	}
	for ab50Index := range ab50Source {
		ab50Source[ab50Index] = 'x'
	}
	if !ab50Bytes.Equal(ab50Copied, ab50EntriesJSON) || !ab50Bytes.Equal(ab50Row, ab50Original) {
		ab50T.Fatal("source aliases returned storage")
	}
	ab50Original[0] = 'x'
	ab50EntriesJSON[0] = 'x'
	ab50Again, _, ab50Error := Ab50Lookup(ab50Key)
	if ab50Error != nil || !ab50Bytes.Equal(ab50Again, ab50Row) {
		ab50T.Fatal("lookup storage leaked")
	}
	ab50AgainEntries, ab50Error := Ab50Entries()
	if ab50Error != nil || !ab50Bytes.Equal(ab50AgainEntries, ab50Copied) {
		ab50T.Fatal("entries storage leaked")
	}
	ab50Duplicate := append([]byte{'['}, ab50Again...)
	ab50Duplicate = append(ab50Duplicate, ',')
	ab50Duplicate = append(ab50Duplicate, ab50Again...)
	ab50Duplicate = append(ab50Duplicate, ']')
	for _, ab50Invalid := range [][]byte{[]byte(`{`), []byte(`null`), []byte(`{}`), []byte(`[null]`), []byte(`[{}]`), []byte(`[{"key":{"family":false}}]`), ab50Duplicate, append(append([]byte{'['}, ab50Again...), []byte(`,{}]`)...)} {
		if ab50Result, ab50Error := ab50Entries(ab50Invalid); ab50Error == nil || ab50Result != nil {
			ab50T.Fatal("invalid entries succeeded")
		}
		if ab50Result, ab50Found, ab50Error := ab50Lookup(ab50Invalid, ab50Key); ab50Error == nil || ab50Found || ab50Result != nil {
			ab50T.Fatal("invalid lookup succeeded")
		}
	}
}

func TestAb50WholeRowMetadataIsolation(ab50T *ab50Testing.T) {
	ab50Key := ab50TestKey(7, "energy")
	ab50Raw, ab50Found, ab50Error := Ab50Lookup(ab50Key)
	if ab50Error != nil || !ab50Found {
		ab50T.Fatal(ab50Error)
	}
	for ab50Name, ab50Value := range map[string]string{
		"grade": `"Unknown"`, "width_bits": `8`,
		"unresolved": `["ab50 unresolved fixture"]`,
		"raw_form":   `{"ab50_raw":[16,39],"ab50_value":10000}`,
	} {
		var ab50Object map[string]ab50JSON.RawMessage
		if ab50JSON.Unmarshal(ab50Raw, &ab50Object) != nil {
			ab50T.Fatal("fixture decode")
		}
		ab50Object[ab50Name] = ab50JSON.RawMessage(ab50Value)
		ab50Changed, ab50Error := ab50JSON.Marshal(ab50Object)
		if ab50Error != nil {
			ab50T.Fatal(ab50Error)
		}
		ab50Source := append(append([]byte{'['}, ab50Changed...), ']')
		ab50Result, ab50Found, ab50Error := ab50Lookup(ab50Source, ab50Key)
		if ab50Error != nil || !ab50Found || !ab50Bytes.Equal(ab50Result, ab50Changed) {
			ab50T.Fatal("whole row was rewritten", ab50Name, ab50Error)
		}
		ab50Later, _, ab50Error := Ab50Lookup(ab50Key)
		if ab50Error != nil || !ab50Bytes.Equal(ab50Later, ab50Raw) {
			ab50T.Fatal("fixture contaminated production")
		}
	}
	for _, ab50Primary := range []string{"null", `"181"`} {
		var ab50Object map[string]ab50JSON.RawMessage
		var ab50Axes map[string]ab50JSON.RawMessage
		if ab50JSON.Unmarshal(ab50Raw, &ab50Object) != nil || ab50JSON.Unmarshal(ab50Object["key"], &ab50Axes) != nil {
			ab50T.Fatal("fixture decode")
		}
		ab50Axes["primary"] = ab50JSON.RawMessage(ab50Primary)
		ab50Object["key"], ab50Error = ab50JSON.Marshal(ab50Axes)
		if ab50Error != nil {
			ab50T.Fatal(ab50Error)
		}
		ab50Changed, ab50Error := ab50JSON.Marshal(ab50Object)
		if ab50Error != nil {
			ab50T.Fatal(ab50Error)
		}
		ab50Source := append(append([]byte{'['}, ab50Changed...), ']')
		if _, ab50Found, ab50Error := ab50Lookup(ab50Source, ab50Key); ab50Error == nil || ab50Found {
			ab50T.Fatal("malformed primary succeeded")
		}
	}
}
