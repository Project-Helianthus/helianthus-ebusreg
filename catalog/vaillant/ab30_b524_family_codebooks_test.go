package vaillant

import (
	ab30bytes "bytes"
	ab30json "encoding/json"
	ab30fmt "fmt"
	ab30reflect "reflect"
	ab30strconv "strconv"
	ab30strings "strings"
	ab30testing "testing"
)

type ab30Meaning struct {
	Ab30Value int    `json:"value"`
	Ab30Label string `json:"label"`
}

type ab30Row struct {
	Ab30ID             string                 `json:"id"`
	Ab30Key            Ab30Key                `json:"key"`
	Ab30Grade          string                 `json:"grade"`
	Ab30Representation string                 `json:"representation_grade"`
	Ab30Enumeration    map[string]ab30Meaning `json:"enumeration"`
	Ab30RawForm        struct {
		Ab30Encoding   string `json:"encoding"`
		Ab30Width      int    `json:"width_bytes"`
		Ab30Signedness string `json:"signedness"`
		Ab30Preserved  bool   `json:"preserved"`
	} `json:"raw_form"`
	Ab30Fallback  string            `json:"raw_fallback"`
	Ab30Relations map[string]string `json:"relations"`
	Ab30Unknown   map[string]string `json:"unknown"`
}

func ab30Pointer[ab30Type any](ab30Value ab30Type) *ab30Type { return &ab30Value }

func ab30Key(ab30Family string, ab30Group uint8, ab30Register uint16, ab30Scope, ab30Alternative string) Ab30Key {
	return Ab30Key{Ab30Family: ab30Family, Ab30Variant: "none", Ab30Primary: 0xB5, Ab30Secondary: ab30Pointer(uint8(0x24)), Ab30Operation: ab30Pointer(uint8(2)), Ab30Direction: "read", Ab30Field: "value", Ab30Alternative: ab30Alternative, Ab30Group: ab30Pointer(ab30Group), Ab30InstanceScope: ab30Scope, Ab30Register: ab30Pointer(ab30Register)}
}

func ab30Get(ab30T *ab30testing.T, ab30Key Ab30Key) ab30Row {
	ab30T.Helper()
	ab30Raw, ab30Found, ab30Error := Ab30Lookup(ab30Key)
	if ab30Error != nil || !ab30Found {
		ab30T.Fatalf("ab30 valid control: %v %v", ab30Found, ab30Error)
	}
	var ab30Result ab30Row
	if ab30Error := ab30json.Unmarshal(ab30Raw, &ab30Result); ab30Error != nil {
		ab30T.Fatal(ab30Error)
	}
	if !ab30reflect.DeepEqual(ab30Result.Ab30Key, ab30Key) {
		ab30T.Fatal("ab30 returned a different whole row")
	}
	return ab30Result
}

func TestAb30Codebooks(ab30T *ab30testing.T) {
	ab30Raw, ab30Error := Ab30Entries()
	if ab30Error != nil {
		ab30T.Fatal(ab30Error)
	}
	var ab30Rows []ab30Row
	if ab30Error := ab30json.Unmarshal(ab30Raw, &ab30Rows); ab30Error != nil {
		ab30T.Fatal(ab30Error)
	}
	if len(ab30Rows) != 16 {
		ab30T.Fatalf("ab30 finite row count: %d", len(ab30Rows))
	}
	ab30Seen := map[string]bool{}
	for _, ab30Family := range []string{"controller-700", "controller-720"} {
		for _, ab30Case := range []struct {
			ab30Group                                                uint8
			ab30Register                                             uint16
			ab30Scope, ab30Alternative, ab30Labels700, ab30Labels720 string
		}{
			{1, 0x0003, "00", "u16", "0:off,1:automatic,2:day", "0:off,1:time-controlled,2:manual"},
			{2, 0x0002, "dynamic", "u16", "0:deactivated,1:heating,2:fixed-flow,3:DHW,4:return-flow,5:pool", "0:deactivated,1:heating,2:fixed-flow,3:DHW,4:return-flow"},
			{3, 0x0001, "dynamic", "u16", "0:off,1:automatic,2:day", "0:off,1:time-controlled,2:manual"},
			{3, 0x0006, "dynamic", "u16", "0:off,1:automatic,2:day,3:setback", "0:off,1:time-controlled,2:manual"},
			{3, 0x0013, "dynamic", "u16", "0:none,1:central,2:remote-circuit-one,3:remote-circuit-two,4:remote-circuit-three,5:remote-circuit-four,6:remote-circuit-five,7:remote-circuit-six,8:remote-circuit-seven,9:remote-circuit-eight", "0:none,1:central,2:remote-circuit-one,3:remote-circuit-two,4:remote-circuit-three,5:remote-circuit-four"},
			{3, 0x001B, "dynamic", "u8", "0:no-quick-mode,1:ventilation-boost,3:party,4:quick-veto,5:one-day-holiday,6:one-day-bank-holiday,7:holiday,8:bank-holiday,10:system-off", "0:no-special-function,1:ventilation-boost,4:quick-veto,7:holiday,10:system-off"},
			{3, 0x001B, "dynamic", "u16", "0:no-quick-mode,1:ventilation-boost,3:party,4:quick-veto,5:one-day-holiday,6:one-day-bank-holiday,7:holiday,8:bank-holiday,10:system-off", "0:no-special-function,1:ventilation-boost,4:quick-veto,7:holiday,10:system-off"},
			{3, 0x001C, "dynamic", "signed8", "255:unbound,0:control-circuit-one,1:control-circuit-two,2:control-circuit-three", "255:unbound,0:control-circuit-one,1:control-circuit-two,2:control-circuit-three,3:control-circuit-four,4:control-circuit-five,5:control-circuit-six,6:control-circuit-seven,7:control-circuit-eight,8:control-circuit-nine"},
		} {
			ab30T.Run(ab30fmt.Sprintf("ab30_%s_%d_%d_%s", ab30Family, ab30Case.ab30Group, ab30Case.ab30Register, ab30Case.ab30Alternative), func(ab30T *ab30testing.T) {
				ab30Got := ab30Get(ab30T, ab30Key(ab30Family, ab30Case.ab30Group, ab30Case.ab30Register, ab30Case.ab30Scope, ab30Case.ab30Alternative))
				if !ab30strings.HasPrefix(ab30Got.Ab30ID, "ab30_") || ab30Seen[ab30Got.Ab30ID] {
					ab30T.Fatal("ab30 row IDs must be unique and prefixed")
				}
				ab30Seen[ab30Got.Ab30ID] = true
				ab30Labels := ab30Case.ab30Labels700
				if ab30Family == "controller-720" {
					ab30Labels = ab30Case.ab30Labels720
				}
				ab30Expected := map[string]ab30Meaning{}
				for _, ab30Pair := range ab30strings.Split(ab30Labels, ",") {
					ab30Parts := ab30strings.SplitN(ab30Pair, ":", 2)
					ab30Value, ab30Error := ab30strconv.Atoi(ab30Parts[0])
					if ab30Error != nil {
						ab30T.Fatal(ab30Error)
					}
					if ab30Case.ab30Alternative == "signed8" && ab30Value == 255 {
						ab30Value = -1
					}
					ab30Expected[ab30Parts[0]] = ab30Meaning{ab30Value, ab30Parts[1]}
				}
				if !ab30reflect.DeepEqual(ab30Got.Ab30Enumeration, ab30Expected) {
					ab30T.Fatalf("ab30 independent exact codebook: got %v want %v", ab30Got.Ab30Enumeration, ab30Expected)
				}
				if ab30Got.Ab30Grade != "Hypothesis" || ab30Got.Ab30Representation != "Proven" {
					ab30T.Fatal("ab30 grades")
				}
				ab30Width, ab30Signedness := 1, "unsigned"
				if ab30Case.ab30Alternative == "u16" {
					ab30Width = 2
				}
				if ab30Case.ab30Alternative == "signed8" {
					ab30Signedness = "signed"
				}
				if ab30Got.Ab30RawForm.Ab30Encoding != ab30Case.ab30Alternative || ab30Got.Ab30RawForm.Ab30Width != ab30Width || ab30Got.Ab30RawForm.Ab30Signedness != ab30Signedness || !ab30Got.Ab30RawForm.Ab30Preserved {
					ab30T.Fatal("ab30 raw form")
				}
				if ab30Got.Ab30Fallback != "Preserve raw bytes; unlisted codes and unknown families have Unknown meaning." {
					ab30T.Fatal("ab30 raw fallback")
				}
				if !ab30reflect.DeepEqual(ab30Got.Ab30Relations, map[string]string{"identity_grade": "Proven", "candidate_opcode": "OP02", "instance_scope_role": "metadata-only", "width_selection": "explicit alternative; reply length does not select", "family_selection": "explicit; no auto-detection"}) {
					ab30T.Fatal("ab30 relations must not infer installed cardinality or width")
				}
				ab30Unknown := map[string]string{}
				for _, ab30Name := range []string{"native_support", "physical_effects", "persistence", "delivery", "installed_cardinality", "present_slots", "freshness", "unlisted_codes"} {
					ab30Unknown[ab30Name] = "Unknown"
				}
				if !ab30reflect.DeepEqual(ab30Got.Ab30Unknown, ab30Unknown) {
					ab30T.Fatal("ab30 qualification remains Unknown")
				}
			})
		}
	}
	for _, ab30Row := range ab30Rows {
		if !ab30Seen[ab30Row.Ab30ID] {
			ab30T.Fatal("ab30 unexpected production row")
		}
	}
}

func TestAb30EveryIdentityAxis(ab30T *ab30testing.T) {
	ab30Control := ab30Key("controller-700", 3, 0x001B, "dynamic", "u8")
	ab30Get(ab30T, ab30Control)
	for ab30Axis := 0; ab30Axis < ab30reflect.TypeOf(ab30Control).NumField(); ab30Axis++ {
		ab30Name := ab30reflect.TypeOf(ab30Control).Field(ab30Axis).Name
		ab30T.Run("ab30_"+ab30Name, func(ab30T *ab30testing.T) {
			ab30Changed := ab30Control
			ab30Slot := ab30reflect.ValueOf(&ab30Changed).Elem().Field(ab30Axis)
			switch ab30Slot.Kind() {
			case ab30reflect.String:
				ab30Slot.SetString("ab30_other")
			case ab30reflect.Uint8:
				ab30Slot.SetUint(0)
			case ab30reflect.Ptr:
				ab30Value := ab30reflect.New(ab30Slot.Type().Elem())
				ab30Value.Elem().SetUint(6)
				ab30Slot.Set(ab30Value)
			default:
				ab30T.Fatal("ab30 unhandled key axis")
			}
			ab30Raw, ab30Found, ab30Error := Ab30Lookup(ab30Changed)
			if ab30Error != nil || ab30Found || ab30Raw != nil {
				ab30T.Fatalf("ab30 changed axis inherited labels: %s %v %v", ab30Name, ab30Found, ab30Error)
			}
			ab30First := ab30Get(ab30T, ab30Control)
			ab30Second := ab30First
			ab30Second.Ab30Key = ab30Changed
			ab30Second.Ab30ID = "ab30_decoy"
			ab30Data, ab30Error := ab30json.Marshal([]ab30Row{ab30First, ab30Second})
			if ab30Error != nil {
				ab30T.Fatal(ab30Error)
			}
			ab30Rows, _, ab30Error := ab30Validate(ab30Data)
			if ab30Error != nil || len(ab30Rows) != 2 {
				ab30T.Fatalf("ab30 valid distinct pair: %v", ab30Error)
			}
			ab30Raw, ab30Found, ab30Error = ab30Lookup(ab30Data, ab30Changed)
			var ab30Got ab30Row
			if ab30Error != nil || !ab30Found {
				ab30T.Fatalf("ab30 decoy control: %v", ab30Error)
			}
			if ab30Error := ab30json.Unmarshal(ab30Raw, &ab30Got); ab30Error != nil || ab30Got.Ab30ID != "ab30_decoy" {
				ab30T.Fatalf("ab30 whole decoy row: %v", ab30Error)
			}
		})
	}
}

func TestAb30OptionalAbsentAndZero(ab30T *ab30testing.T) {
	ab30Control := ab30Get(ab30T, ab30Key("controller-700", 3, 27, "dynamic", "u8"))
	for _, ab30Name := range []string{"Ab30Secondary", "Ab30Operation", "Ab30Group", "Ab30Register"} {
		ab30T.Run("ab30_"+ab30Name, func(ab30T *ab30testing.T) {
			ab30Absent, ab30Zero := ab30Control, ab30Control
			ab30Absent.Ab30ID = "ab30_absent"
			ab30Zero.Ab30ID = "ab30_zero"
			ab30reflect.ValueOf(&ab30Absent.Ab30Key).Elem().FieldByName(ab30Name).SetZero()
			ab30Slot := ab30reflect.ValueOf(&ab30Zero.Ab30Key).Elem().FieldByName(ab30Name)
			ab30Slot.Set(ab30reflect.New(ab30Slot.Type().Elem()))
			ab30Data, ab30Error := ab30json.Marshal([]ab30Row{ab30Absent, ab30Zero})
			if ab30Error != nil {
				ab30T.Fatal(ab30Error)
			}
			for _, ab30Expected := range []ab30Row{ab30Absent, ab30Zero} {
				ab30Raw, ab30Found, ab30Error := ab30Lookup(ab30Data, ab30Expected.Ab30Key)
				var ab30Got ab30Row
				if ab30Error != nil || !ab30Found {
					ab30T.Fatalf("ab30 absent/zero control: %v", ab30Error)
				}
				if ab30Error := ab30json.Unmarshal(ab30Raw, &ab30Got); ab30Error != nil || ab30Got.Ab30ID != ab30Expected.Ab30ID {
					ab30T.Fatalf("ab30 absent and zero collapsed: %v", ab30Error)
				}
				ab30Raw, ab30Found, ab30Error = Ab30Lookup(ab30Expected.Ab30Key)
				if ab30Error != nil || ab30Found || ab30Raw != nil {
					ab30T.Fatal("ab30 missing optional identity inherited production labels")
				}
			}
		})
	}
}

func TestAb30NamespaceAndRawFallback(ab30T *ab30testing.T) {
	for _, ab30Family := range []string{"controller-700", "controller-720"} {
		ab30Signed := ab30Get(ab30T, ab30Key(ab30Family, 3, 28, "dynamic", "signed8"))
		if ab30Signed.Ab30Enumeration["255"] != (ab30Meaning{-1, "unbound"}) {
			ab30T.Fatal("ab30 FF candidate signed8")
		}
		ab30Unsigned := ab30Get(ab30T, ab30Key(ab30Family, 3, 6, "dynamic", "u16"))
		for _, ab30Raw := range []string{"255", "65535", "-1"} {
			if _, ab30Found := ab30Unsigned.Ab30Enumeration[ab30Raw]; ab30Found {
				ab30T.Fatal("ab30 unrelated u16 sign extension")
			}
		}
		for _, ab30Alternative := range []string{"u8", "u16"} {
			ab30Got := ab30Get(ab30T, ab30Key(ab30Family, 3, 27, "dynamic", ab30Alternative))
			if _, ab30Found := ab30Got.Ab30Enumeration["2"]; ab30Found {
				ab30T.Fatal("ab30 unlisted raw code labelled")
			}
		}
		for _, ab30Changed := range []Ab30Key{
			ab30Key("ab30_unknown-family", 2, 2, "dynamic", "u16"),
			ab30Key(ab30Family, 3, 28, "dynamic", "u16"),
			ab30Key(ab30Family, 3, 27, "00", "u8"),
			ab30Key(ab30Family, 1, 3, "dynamic", "u16"),
			ab30Key(ab30Family, 3, 27, "dynamic", "none"),
		} {
			ab30Raw, ab30Found, ab30Error := Ab30Lookup(ab30Changed)
			if ab30Error != nil || ab30Found || ab30Raw != nil {
				ab30T.Fatal("ab30 cross-namespace fallback")
			}
		}
		ab30Opcode := ab30Key(ab30Family, 3, 27, "dynamic", "u8")
		ab30Opcode.Ab30Operation = ab30Pointer(uint8(6))
		ab30Raw, ab30Found, ab30Error := Ab30Lookup(ab30Opcode)
		if ab30Error != nil || ab30Found || ab30Raw != nil {
			ab30T.Fatal("ab30 OP06 inherited OP02")
		}
	}
}

func TestAb30InvalidEmbeddedData(ab30T *ab30testing.T) {
	ab30Key := ab30Key("controller-700", 3, 27, "dynamic", "u8")
	ab30ControlRow := ab30Get(ab30T, ab30Key)
	ab30Duplicate, ab30Error := ab30json.Marshal([]ab30Row{ab30ControlRow, ab30ControlRow})
	if ab30Error != nil {
		ab30T.Fatal(ab30Error)
	}
	ab30Cases := [][]byte{[]byte("["), []byte("null"), []byte("{}"), []byte("[null]"), []byte("[{}]"), ab30Duplicate}
	ab30Valid, ab30Error := ab30json.Marshal([]ab30Row{ab30ControlRow})
	if ab30Error != nil {
		ab30T.Fatal(ab30Error)
	}
	var ab30Objects []map[string]ab30json.RawMessage
	if ab30Error := ab30json.Unmarshal(ab30Valid, &ab30Objects); ab30Error != nil {
		ab30T.Fatal(ab30Error)
	}
	var ab30Identity map[string]ab30json.RawMessage
	if ab30Error := ab30json.Unmarshal(ab30Objects[0]["key"], &ab30Identity); ab30Error != nil {
		ab30T.Fatal(ab30Error)
	}
	for ab30Axis := range ab30Identity {
		for _, ab30Replacement := range []ab30json.RawMessage{nil, []byte("true")} {
			ab30Copy := map[string]ab30json.RawMessage{}
			for ab30Name, ab30Value := range ab30Identity {
				ab30Copy[ab30Name] = ab30Value
			}
			if ab30Replacement == nil {
				delete(ab30Copy, ab30Axis)
			} else {
				ab30Copy[ab30Axis] = ab30Replacement
			}
			ab30KeyData, ab30Error := ab30json.Marshal(ab30Copy)
			if ab30Error != nil {
				ab30T.Fatal(ab30Error)
			}
			ab30InvalidData, ab30Error := ab30json.Marshal([]map[string]ab30json.RawMessage{{"key": ab30KeyData}})
			if ab30Error != nil {
				ab30T.Fatal(ab30Error)
			}
			ab30Cases = append(ab30Cases, ab30InvalidData)
		}
	}
	for ab30Index, ab30Data := range ab30Cases {
		ab30T.Run(ab30fmt.Sprintf("ab30_invalid_%d", ab30Index), func(ab30T *ab30testing.T) {
			if _, _, ab30Error := ab30Validate(ab30Data); ab30Error == nil {
				ab30T.Fatal("ab30 invalid data admitted")
			}
			ab30Raw, ab30Found, ab30Error := ab30Lookup(ab30Data, ab30Key)
			if ab30Error == nil || ab30Found || ab30Raw != nil {
				ab30T.Fatal("ab30 invalid lookup must fail closed")
			}
		})
	}
}

func TestAb30IndependentBuffers(ab30T *ab30testing.T) {
	ab30Key := ab30Key("controller-700", 3, 27, "dynamic", "u8")
	ab30Entries, ab30Error := Ab30Entries()
	if ab30Error != nil {
		ab30T.Fatal(ab30Error)
	}
	ab30Expected := append([]byte(nil), ab30Entries...)
	ab30First, ab30Found, ab30Error := Ab30Lookup(ab30Key)
	if ab30Error != nil || !ab30Found {
		ab30T.Fatal("ab30 lookup control")
	}
	ab30RowExpected := append([]byte(nil), ab30First...)
	ab30Second, ab30Found, ab30Error := Ab30Lookup(ab30Key)
	if ab30Error != nil || !ab30Found {
		ab30T.Fatal("ab30 second lookup control")
	}
	for ab30Index := range ab30Entries {
		ab30Entries[ab30Index] = 'x'
	}
	for ab30Index := range ab30First {
		ab30First[ab30Index] = 'y'
	}
	ab30Later, ab30Error := Ab30Entries()
	if ab30Error != nil || !ab30bytes.Equal(ab30Later, ab30Expected) {
		ab30T.Fatal("ab30 entries share storage")
	}
	ab30RowLater, ab30Found, ab30Error := Ab30Lookup(ab30Key)
	if ab30Error != nil || !ab30Found || !ab30bytes.Equal(ab30RowLater, ab30RowExpected) || !ab30bytes.Equal(ab30Second, ab30RowExpected) {
		ab30T.Fatal("ab30 row buffers share storage")
	}
	ab30Rows, _, ab30Error := ab30Validate(ab30Expected)
	if ab30Error != nil {
		ab30T.Fatal(ab30Error)
	}
	for _, ab30Row := range ab30Rows {
		for ab30Index := range ab30Row {
			ab30Row[ab30Index] = 'z'
		}
	}
	if !ab30bytes.Equal(ab30Later, ab30Expected) {
		ab30T.Fatal("ab30 decoded rows alias input")
	}
}
