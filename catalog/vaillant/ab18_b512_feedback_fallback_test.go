package vaillant

import (
	ab18bytes "bytes"
	ab18json "encoding/json"
	ab18reflect "reflect"
	ab18testing "testing"
)

type ab18Row struct {
	Ab18ID             string            `json:"id"`
	Ab18Key            Ab18Key           `json:"key"`
	Ab18Grade          string            `json:"grade"`
	Ab18Position       int               `json:"logical_position"`
	Ab18Type           string            `json:"type"`
	Ab18PhysicalOffset string            `json:"physical_offset"`
	Ab18RawForm        map[string]string `json:"raw_form"`
	Ab18RawFallback    string            `json:"raw_fallback"`
	Ab18Values         map[string]string `json:"values"`
	Ab18Other          string            `json:"other_values"`
	Ab18Polarity       string            `json:"polarity"`
	Ab18Feedback       struct {
		Ab18Absent        string `json:"actual_signal_absent"`
		Ab18Present       string `json:"actual_signal_present"`
		Ab18Equal         bool   `json:"equal_to_setpoint_establishes_actual"`
		Ab18Discriminator string `json:"wire_discriminator"`
		Ab18Physical      string `json:"physical_value"`
	} `json:"feedback"`
	Ab18Default     bool                `json:"boolean_default"`
	Ab18RawDefault  ab18json.RawMessage `json:"raw_default"`
	Ab18Encoding    string              `json:"raw_default_encoding"`
	Ab18Native      string              `json:"native_support"`
	Ab18Effects     string              `json:"physical_effects"`
	Ab18Persistence string              `json:"persistence"`
	Ab18Delivery    string              `json:"delivery"`
}

func ab18Read(ab18T *ab18testing.T) ([]ab18json.RawMessage, []ab18Row) {
	ab18T.Helper()
	ab18Data, ab18Error := Ab18Entries()
	if ab18Error != nil {
		ab18T.Fatal(ab18Error)
	}
	var ab18Raw []ab18json.RawMessage
	var ab18Rows []ab18Row
	if ab18Error = ab18json.Unmarshal(ab18Data, &ab18Raw); ab18Error != nil {
		ab18T.Fatal(ab18Error)
	}
	if ab18Error = ab18json.Unmarshal(ab18Data, &ab18Rows); ab18Error != nil {
		ab18T.Fatal(ab18Error)
	}
	return ab18Raw, ab18Rows
}

func TestAb18FiniteMetadata(ab18T *ab18testing.T) {
	ab18Raw, ab18Rows := ab18Read(ab18T)
	if len(ab18Rows) != 12 {
		ab18T.Fatalf("ab18: rows = %d", len(ab18Rows))
	}
	ab18Fields := []string{"secondary-heating-pump", "circulation-pump", "solar-pump", "error-relay", "loading-pump", "fume-flap"}
	ab18Request := []byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26}
	ab18Positive := []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36}
	for ab18Side, ab18Direction := range []string{"request", "positive"} {
		for ab18Index, ab18Field := range ab18Fields {
			ab18T.Run("ab18-"+ab18Direction+"-"+ab18Field, func(ab18T *ab18testing.T) {
				ab18Offset := ab18Side*6 + ab18Index
				ab18Row := ab18Rows[ab18Offset]
				ab18Secondary, ab18Operation := uint8(18), uint8(1)
				ab18ExpectedKey := Ab18Key{"B512", "OP01", 181, &ab18Secondary, &ab18Operation, ab18Direction, ab18Field, "none", nil, "none", nil}
				if !ab18reflect.DeepEqual(ab18Row.Ab18Key, ab18ExpectedKey) || ab18Row.Ab18ID != "ab18-"+ab18Direction+"-"+ab18Field {
					ab18T.Fatal("ab18: identity")
				}
				ab18FoundRow, ab18Found, ab18Error := Ab18Lookup(ab18ExpectedKey)
				if ab18Error != nil || !ab18Found || !ab18bytes.Equal(ab18FoundRow, ab18Raw[ab18Offset]) {
					ab18T.Fatal("ab18: whole row lookup", ab18Error)
				}
				if ab18Row.Ab18Position != ab18Index+3 || ab18Row.Ab18Type != "u8" || ab18Row.Ab18Grade != "Hypothesis" {
					ab18T.Fatal("ab18: independent position/type/grade")
				}
				if ab18Row.Ab18PhysicalOffset != "Unknown" || ab18Row.Ab18Polarity != "Unknown" || ab18Row.Ab18Native != "Unknown" || ab18Row.Ab18Effects != "Unknown" || ab18Row.Ab18Persistence != "Unknown" || ab18Row.Ab18Delivery != "Unknown" {
					ab18T.Fatal("ab18: qualification")
				}
				if !ab18Row.Ab18Default || string(ab18Row.Ab18RawDefault) != "null" || ab18Row.Ab18Encoding != "Unknown" {
					ab18T.Fatal("ab18: boolean default is not raw 01")
				}
				if ab18Row.Ab18RawForm["type"] != "u8" || ab18Row.Ab18RawForm["retention"] != "preserve each independent raw token" || ab18Row.Ab18RawFallback != "retain raw token and candidate meaning independently" {
					ab18T.Fatal("ab18: raw retention")
				}
				if ab18Row.Ab18Feedback.Ab18Absent != "setpoint-fallback" || ab18Row.Ab18Feedback.Ab18Present != "actual-signal" || ab18Row.Ab18Feedback.Ab18Equal || ab18Row.Ab18Feedback.Ab18Discriminator != "Unknown" || ab18Row.Ab18Feedback.Ab18Physical != "Unknown" {
					ab18T.Fatal("ab18: feedback provenance")
				}
				ab18Tokens := append([]byte(nil), ab18Request...)
				if ab18Side == 1 {
					ab18Tokens = append([]byte(nil), ab18Positive...)
				}
				ab18Before := append([]byte(nil), ab18Tokens...)
				ab18Changed := append([]byte(nil), ab18Tokens...)
				ab18Changed[ab18Row.Ab18Position-3] = 0x71 + byte(ab18Index)
				for ab18Other := range ab18Changed {
					ab18ExpectedToken := ab18Before[ab18Other]
					if ab18Other == ab18Index {
						ab18ExpectedToken = 0x71 + byte(ab18Index)
					}
					if ab18Changed[ab18Other] != ab18ExpectedToken {
						ab18T.Fatal("ab18: field selection aliases another position")
					}
				}
				for _, ab18ActualPresent := range []bool{false, true} {
					ab18Setpoint := ab18Before[ab18Index]
					ab18FeedbackToken := ab18Setpoint
					ab18Source := ab18Row.Ab18Feedback.Ab18Absent
					ab18ExpectedSource := "setpoint-fallback"
					if ab18ActualPresent {
						ab18Source = ab18Row.Ab18Feedback.Ab18Present
						ab18ExpectedSource = "actual-signal"
					}
					if ab18Source != ab18ExpectedSource || ab18FeedbackToken != ab18Setpoint || ab18Row.Ab18Feedback.Ab18Equal {
						ab18T.Fatal("ab18: equal tokens preserve source distinction")
					}
				}
				if ab18Index != 3 && ab18Index != 5 {
					ab18Expected := "leave-unchanged"
					if ab18Side == 1 {
						ab18Expected = "function-unavailable"
					}
					if len(ab18Row.Ab18Values) != 1 || ab18Row.Ab18Values["FF"] != ab18Expected || ab18Row.Ab18Other != "retain-raw" {
						ab18T.Fatal("ab18: direction-specific FF")
					}
					// The metadata selects a synthetic field; request FF retains that field's previous token.
					ab18Tokens[ab18Row.Ab18Position-3] = 0xff
					if ab18Row.Ab18Values["FF"] == "leave-unchanged" {
						ab18Tokens[ab18Row.Ab18Position-3] = ab18Before[ab18Index]
					}
					for ab18Other := range ab18Tokens {
						ab18ExpectedToken := ab18Before[ab18Other]
						if ab18Other == ab18Index && ab18Side == 1 {
							ab18ExpectedToken = 0xff
						}
						if ab18Tokens[ab18Other] != ab18ExpectedToken {
							ab18T.Fatal("ab18: independent raw tokens")
						}
					}
				} else {
					if len(ab18Row.Ab18Values) != 2 || ab18Row.Ab18Values["0F"] != "state-0F" || ab18Row.Ab18Values["F0"] != "state-F0" || ab18Row.Ab18Other != "no-action" {
						ab18T.Fatal("ab18: declared states without polarity")
					}
					for _, ab18Token := range []string{"00", "01", "FF"} {
						if _, ab18Action := ab18Row.Ab18Values[ab18Token]; ab18Action {
							ab18T.Fatal("ab18: unexpected action", ab18Token)
						}
					}
				}
			})
		}
	}
}

func TestAb18ExactAxes(ab18T *ab18testing.T) {
	ab18Raw, ab18Rows := ab18Read(ab18T)
	for ab18Axis := 0; ab18Axis < 11; ab18Axis++ {
		ab18Name := ab18reflect.TypeOf(Ab18Key{}).Field(ab18Axis).Name
		ab18T.Run(ab18Name, func(ab18T *ab18testing.T) {
			ab18Key := ab18Rows[0].Ab18Key
			ab18Value := ab18reflect.ValueOf(&ab18Key).Elem().Field(ab18Axis)
			switch ab18Value.Kind() {
			case ab18reflect.String:
				ab18Value.SetString(ab18Value.String() + "-ab18-control")
			case ab18reflect.Uint8:
				ab18Value.SetUint(ab18Value.Uint() + 1)
			case ab18reflect.Ptr:
				ab18Pointer := ab18reflect.New(ab18Value.Type().Elem())
				if !ab18Value.IsNil() {
					ab18Pointer.Elem().SetUint(ab18Value.Elem().Uint() + 1)
				}
				ab18Value.Set(ab18Pointer)
			default:
				ab18T.Fatal("ab18: unhandled axis")
			}
			if _, ab18Found, ab18Error := Ab18Lookup(ab18Key); ab18Error != nil || ab18Found {
				ab18T.Fatal("ab18: public exact miss", ab18Error)
			}
			var ab18Control map[string]ab18json.RawMessage
			if ab18Error := ab18json.Unmarshal(ab18Raw[0], &ab18Control); ab18Error != nil {
				ab18T.Fatal(ab18Error)
			}
			ab18Control["key"], _ = ab18json.Marshal(ab18Key)
			ab18ControlRaw, _ := ab18json.Marshal(ab18Control)
			ab18Fixture, _ := ab18json.Marshal([]ab18json.RawMessage{ab18Raw[0], ab18ControlRaw})
			for _, ab18Expected := range []struct {
				Ab18Key Ab18Key
				Ab18Raw ab18json.RawMessage
			}{{ab18Rows[0].Ab18Key, ab18Raw[0]}, {ab18Key, ab18ControlRaw}} {
				ab18Result, ab18Found, ab18Error := ab18Lookup(ab18Fixture, ab18Expected.Ab18Key)
				var ab18Compact ab18bytes.Buffer
				if ab18Error := ab18json.Compact(&ab18Compact, ab18Expected.Ab18Raw); ab18Error != nil {
					ab18T.Fatal(ab18Error)
				}
				if ab18Error != nil || !ab18Found || !ab18bytes.Equal(ab18Result, ab18Compact.Bytes()) {
					ab18T.Fatal("ab18: independent valid control", ab18Error)
				}
			}
		})
	}
}

func TestAb18OptionalPresence(ab18T *ab18testing.T) {
	ab18Raw, ab18Rows := ab18Read(ab18T)
	for _, ab18Axis := range []int{3, 4, 8, 10} {
		ab18T.Run(ab18reflect.TypeOf(Ab18Key{}).Field(ab18Axis).Name, func(ab18T *ab18testing.T) {
			ab18Keys := []Ab18Key{ab18Rows[0].Ab18Key, ab18Rows[0].Ab18Key}
			ab18reflect.ValueOf(&ab18Keys[0]).Elem().Field(ab18Axis).SetZero()
			ab18Value := ab18reflect.ValueOf(&ab18Keys[1]).Elem().Field(ab18Axis)
			ab18Value.Set(ab18reflect.New(ab18Value.Type().Elem()))
			var ab18Controls []ab18json.RawMessage
			for _, ab18Key := range ab18Keys {
				var ab18Control map[string]ab18json.RawMessage
				if ab18Error := ab18json.Unmarshal(ab18Raw[0], &ab18Control); ab18Error != nil {
					ab18T.Fatal(ab18Error)
				}
				ab18Control["key"], _ = ab18json.Marshal(ab18Key)
				ab18ControlRaw, _ := ab18json.Marshal(ab18Control)
				ab18Controls = append(ab18Controls, ab18ControlRaw)
			}
			ab18Fixture, _ := ab18json.Marshal(ab18Controls)
			for ab18Index, ab18Key := range ab18Keys {
				ab18Result, ab18Found, ab18Error := ab18Lookup(ab18Fixture, ab18Key)
				if ab18Error != nil || !ab18Found || !ab18bytes.Equal(ab18Result, ab18Controls[ab18Index]) {
					ab18T.Fatal("ab18: absent/zero", ab18Error)
				}
			}
		})
	}
}

func TestAb18Validation(ab18T *ab18testing.T) {
	ab18Raw, ab18Rows := ab18Read(ab18T)
	ab18Duplicate, _ := ab18json.Marshal([]ab18json.RawMessage{ab18Raw[0], ab18Raw[0]})
	for _, ab18Input := range [][]byte{[]byte("["), []byte("{}"), []byte("null"), []byte("[null]"), []byte("[{}]"), ab18Duplicate, []byte(`[{"key":{"family":"B512"}}]`)} {
		if _, ab18Error := ab18Validate(ab18Input); ab18Error == nil {
			ab18T.Fatal("ab18: accepted invalid data", string(ab18Input))
		}
		if ab18Result, ab18Found, ab18Error := ab18Lookup(ab18Input, ab18Rows[0].Ab18Key); ab18Error == nil || ab18Found || ab18Result != nil {
			ab18T.Fatal("ab18: fail closed")
		}
	}
	// An error anywhere prevents returning an earlier valid matching row.
	ab18Poisoned := append(append([]byte{'['}, ab18Raw[0]...), []byte(",{}]")...)
	if _, ab18Found, ab18Error := ab18Lookup(ab18Poisoned, ab18Rows[0].Ab18Key); ab18Error == nil || ab18Found {
		ab18T.Fatal("ab18: whole catalog validation")
	}
}

func TestAb18IndependentBuffers(ab18T *ab18testing.T) {
	ab18Raw, ab18Rows := ab18Read(ab18T)
	ab18Entries, ab18Error := Ab18Entries()
	if ab18Error != nil {
		ab18T.Fatal(ab18Error)
	}
	ab18Expected := append([]byte(nil), ab18Entries...)
	ab18Other, ab18Error := Ab18Entries()
	if ab18Error != nil {
		ab18T.Fatal(ab18Error)
	}
	for ab18Index := range ab18Entries {
		ab18Entries[ab18Index] = 'x'
	}
	ab18Later, ab18Error := Ab18Entries()
	if ab18Error != nil || !ab18bytes.Equal(ab18Other, ab18Expected) || !ab18bytes.Equal(ab18Later, ab18Expected) {
		ab18T.Fatal("ab18: shared entries")
	}
	for ab18Index, ab18Row := range ab18Rows {
		ab18First, ab18Found, ab18Error := Ab18Lookup(ab18Row.Ab18Key)
		if ab18Error != nil || !ab18Found {
			ab18T.Fatal(ab18Error)
		}
		ab18Second, ab18Found, ab18Error := Ab18Lookup(ab18Row.Ab18Key)
		if ab18Error != nil || !ab18Found {
			ab18T.Fatal(ab18Error)
		}
		for ab18Byte := range ab18First {
			ab18First[ab18Byte] = 'x'
		}
		ab18Third, ab18Found, ab18Error := Ab18Lookup(ab18Row.Ab18Key)
		if ab18Error != nil || !ab18Found || !ab18bytes.Equal(ab18Second, ab18Raw[ab18Index]) || !ab18bytes.Equal(ab18Third, ab18Raw[ab18Index]) {
			ab18T.Fatal("ab18: shared row")
		}
	}
}
