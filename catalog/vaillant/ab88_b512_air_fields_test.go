package vaillant

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func ab88Get(ab88T *testing.T, ab88Operation uint16, ab88Direction, ab88Name string) Ab88AirField {
	ab88T.Helper()
	ab88Field, ab88OK := Ab88LookupB512AirField(Ab88AirKey{Ab88Family: "vaillant", Ab88Variant: "B512", Ab88Primary: 0xb512, Ab88Operation: Ab88Axis{Ab88Present: true, Ab88Value: ab88Operation}, Ab88Direction: ab88Direction, Ab88Field: ab88Name})
	if !ab88OK {
		ab88T.Fatalf("ab88 missing %d/%s/%s", ab88Operation, ab88Direction, ab88Name)
	}
	return ab88Field
}

func TestAb88ProductionFields(ab88T *testing.T) {
	ab88Expected := []struct {
		ab88Operation           uint16
		ab88Direction, ab88Name string
		ab88Position, ab88Width int
		ab88Signed              bool
		ab88Divisor             float64
		ab88Unit                string
		ab88Default             *float64
	}{}
	ab88Zero, ab88One := 0.0, 1.0
	ab88Add := func(ab88Operation uint16, ab88Direction, ab88Name string, ab88Position, ab88Width int, ab88Signed bool, ab88Divisor float64, ab88Unit string, ab88Default *float64) {
		ab88Expected = append(ab88Expected, struct {
			ab88Operation           uint16
			ab88Direction, ab88Name string
			ab88Position, ab88Width int
			ab88Signed              bool
			ab88Divisor             float64
			ab88Unit                string
			ab88Default             *float64
		}{ab88Operation, ab88Direction, ab88Name, ab88Position, ab88Width, ab88Signed, ab88Divisor, ab88Unit, ab88Default})
	}
	ab88Add(11, "request", "target_fan_1", 3, 2, false, 1, "rpm", &ab88Zero)
	ab88Add(11, "request", "target_fan_2", 5, 2, false, 1, "rpm", &ab88Zero)
	ab88Add(11, "request", "relay", 7, 1, false, 1, "", &ab88Zero)
	ab88Add(11, "positive", "actual_fan_1", 3, 2, false, 1, "rpm", &ab88Zero)
	ab88Add(11, "positive", "actual_fan_2", 5, 2, false, 1, "rpm", &ab88Zero)
	ab88Add(11, "positive", "status", 7, 1, false, 1, "", &ab88Zero)
	ab88Add(11, "positive", "temperature_1", 8, 2, true, 16, "degrees Celsius", &ab88Zero)
	ab88Add(11, "positive", "temperature_2", 10, 2, true, 16, "degrees Celsius", &ab88Zero)
	ab88Add(11, "positive", "temperature_3", 12, 2, true, 16, "degrees Celsius", &ab88Zero)
	for ab88Index, ab88Name := range []string{"pulse", "P", "I", "D"} {
		ab88Add(12, "request", ab88Name, 3+ab88Index, 1, false, 1, "", nil)
	}
	ab88Add(12, "request", "minimum_pwm", 7, 2, true, 256, "percent", &ab88One)
	ab88Add(12, "request", "maximum_pwm", 9, 2, true, 256, "percent", &ab88One)
	ab88Add(12, "request", "fan_count", 11, 1, false, 1, "", nil)
	ab88Add(12, "positive", "feedback", 3, 1, false, 1, "", nil)
	ab88Add(13, "positive", "digital", 3, 1, false, 1, "", &ab88Zero)
	ab88Add(13, "positive", "temperature_1", 4, 2, true, 16, "degrees Celsius", nil)
	ab88Add(13, "positive", "temperature_2", 6, 2, true, 16, "degrees Celsius", nil)
	ab88Add(14, "positive", "temperature", 3, 2, true, 16, "degrees Celsius", &ab88Zero)
	ab88Add(14, "positive", "humidity", 5, 2, false, 16, "percent", &ab88Zero)
	ab88Add(14, "positive", "air_quality", 7, 2, false, 1, "", &ab88Zero)
	ab88Metadata := Ab88B512AirFields()
	if len(ab88Metadata.Ab88Fields) != len(ab88Expected) {
		ab88T.Fatal("ab88 finite field count")
	}
	ab88Seen := map[Ab88AirKey]bool{}
	for _, ab88Want := range ab88Expected {
		ab88T.Run(fmt.Sprintf("ab88_%d_%s_%s", ab88Want.ab88Operation, ab88Want.ab88Direction, ab88Want.ab88Name), func(ab88T *testing.T) {
			ab88Got := ab88Get(ab88T, ab88Want.ab88Operation, ab88Want.ab88Direction, ab88Want.ab88Name)
			if ab88Seen[ab88Got.Ab88Key] {
				ab88T.Fatal("ab88 duplicate")
			}
			ab88Seen[ab88Got.Ab88Key] = true
			ab88Signedness, ab88Order, ab88Conversion := "unsigned", "not_applicable", "identity"
			if ab88Want.ab88Signed {
				ab88Signedness = "signed"
			}
			if ab88Want.ab88Width == 2 {
				ab88Order = "little_endian"
			}
			if ab88Want.ab88Divisor != 1 {
				ab88Conversion = "divide"
			}
			if ab88Got.Ab88LogicalPosition == nil || *ab88Got.Ab88LogicalPosition != ab88Want.ab88Position || ab88Got.Ab88Width == nil || *ab88Got.Ab88Width != ab88Want.ab88Width || ab88Got.Ab88Signedness != ab88Signedness || ab88Got.Ab88ByteOrder != ab88Order || ab88Got.Ab88Conversion != ab88Conversion || ab88Got.Ab88Divisor != ab88Want.ab88Divisor || !reflect.DeepEqual(ab88Got.Ab88Default, ab88Want.ab88Default) {
				ab88T.Fatalf("ab88 layout/conversion/default: %+v", ab88Got)
			}
			if (ab88Got.Ab88Unit == nil) != (ab88Want.ab88Unit == "") || (ab88Got.Ab88Unit != nil && *ab88Got.Ab88Unit != ab88Want.ab88Unit) {
				ab88T.Fatal("ab88 unit")
			}
			if ab88Got.Ab88Grade != "Hypothesis" || ab88Got.Ab88RawFallback != "preserve original field bytes including unlabeled and reserved values" || ab88Got.Ab88Replacement != nil {
				ab88T.Fatal("ab88 evidence or replacement")
			}
			var ab88DefaultLabel *string
			if ab88Want.ab88Name == "feedback" {
				ab88Rejected := "rejected"
				ab88DefaultLabel = &ab88Rejected
			}
			if !reflect.DeepEqual(ab88Got.Ab88DefaultLabel, ab88DefaultLabel) {
				ab88T.Fatal("ab88 default label")
			}
			ab88Bounds := []float64(nil)
			if ab88Want.ab88Name == "minimum_pwm" || ab88Want.ab88Name == "maximum_pwm" {
				ab88Bounds = []float64{-128, 127.99609375}
			}
			if !reflect.DeepEqual(ab88Got.Ab88Bounds, ab88Bounds) {
				ab88T.Fatal("ab88 bounds")
			}
			ab88Enum := map[string]string{}
			var ab88Otherwise *string
			switch ab88Want.ab88Name {
			case "status":
				ab88Enum = map[string]string{"0": "idle", "1": "running", "255": "fault"}
			case "feedback":
				ab88Enum = map[string]string{"1": "accepted"}
				ab88Rejected := "rejected"
				ab88Otherwise = &ab88Rejected
			case "digital":
				ab88Enum = map[string]string{"0": "false", "1": "true"}
			}
			if !reflect.DeepEqual(ab88Got.Ab88Enumeration, ab88Enum) || !reflect.DeepEqual(ab88Got.Ab88Otherwise, ab88Otherwise) {
				ab88T.Fatal("ab88 enumeration")
			}
			ab88Bits := map[string]map[string]string{}
			ab88Reserved := []int{}
			if ab88Want.ab88Name == "relay" {
				ab88Bits = map[string]map[string]string{"0": {"0": "15V", "1": "tray-heater relay"}, "1": {"0": "15V", "1": "tray-heater relay"}}
				ab88Reserved = []int{2, 3, 4, 5, 6, 7}
			}
			if !reflect.DeepEqual(ab88Got.Ab88Bits, ab88Bits) || !reflect.DeepEqual(ab88Got.Ab88ReservedBits, ab88Reserved) {
				ab88T.Fatal("ab88 relay/reserved roles")
			}
		})
	}
	if ab88Metadata.Ab88Representation != "Proven" || ab88Metadata.Ab88PhysicalSupport != "Unknown" || ab88Metadata.Ab88Acknowledgement != "Unknown" || ab88Metadata.Ab88Complete || ab88Metadata.Ab88Coordinates != "logical positions; not wire offsets" || len(ab88Metadata.Ab88Limitations) != 4 {
		ab88T.Fatal("ab88 qualification boundary")
	}
	if !reflect.DeepEqual(ab88Metadata.Ab88Operations, []Ab88AirOperation{{11, true, "Hypothesis"}, {12, true, "Hypothesis"}, {13, false, "Hypothesis"}, {14, false, "Hypothesis"}}) {
		ab88T.Fatal("ab88 no-argument operations")
	}
}

func ab88Evaluate(ab88Field Ab88AirField, ab88Bytes []byte) (float64, bool) {
	if ab88Field.Ab88Width == nil || len(ab88Bytes) != *ab88Field.Ab88Width {
		return 0, false
	}
	ab88Raw := int64(ab88Bytes[0])
	if *ab88Field.Ab88Width == 2 {
		if ab88Field.Ab88ByteOrder != "little_endian" {
			return 0, false
		}
		ab88Raw |= int64(ab88Bytes[1]) << 8
		if ab88Field.Ab88Signedness == "signed" {
			ab88Raw = int64(int16(ab88Raw))
		}
	}
	ab88Value := float64(ab88Raw)
	switch ab88Field.Ab88Conversion {
	case "divide":
		ab88Value /= ab88Field.Ab88Divisor
	case "identity":
	default:
		return 0, false
	}
	return ab88Value, true
}

func TestAb88SyntheticValues(ab88T *testing.T) {
	for _, ab88Case := range []struct {
		ab88Op            uint16
		ab88Dir, ab88Name string
		ab88Bytes         []byte
		ab88Want          float64
	}{
		{11, "positive", "temperature_1", []byte{0x10, 0}, 1}, {11, "positive", "temperature_2", []byte{0xf0, 0xff}, -1},
		{12, "request", "minimum_pwm", []byte{0, 1}, 1}, {12, "request", "minimum_pwm", []byte{0, 0x80}, -128}, {12, "request", "maximum_pwm", []byte{0xff, 0x7f}, 127.99609375},
		{14, "positive", "humidity", []byte{0x40, 6}, 100}, {14, "positive", "humidity", []byte{0x50, 6}, 101}, {14, "positive", "air_quality", []byte{0xff, 0xff}, 65535},
	} {
		ab88Field := ab88Get(ab88T, ab88Case.ab88Op, ab88Case.ab88Dir, ab88Case.ab88Name)
		if ab88Got, ab88OK := ab88Evaluate(ab88Field, ab88Case.ab88Bytes); !ab88OK || ab88Got != ab88Case.ab88Want {
			ab88T.Fatalf("ab88 %s got %v want %v", ab88Case.ab88Name, ab88Got, ab88Case.ab88Want)
		}
		for _, ab88Invalid := range [][]byte{ab88Case.ab88Bytes[:1], append(append([]byte{}, ab88Case.ab88Bytes...), 0)} {
			if _, ab88OK := ab88Evaluate(ab88Field, ab88Invalid); ab88OK {
				ab88T.Fatal("ab88 field extent")
			}
		}
	}
	for _, ab88Dir := range []string{"request", "positive"} {
		ab88Logical := []byte{0, 0, 0, 2, 1, 4, 3, 0, 0x10, 0, 0xf0, 0xff, 0, 0}
		ab88Prefix := "target"
		if ab88Dir == "positive" {
			ab88Prefix = "actual"
		}
		for ab88Index, ab88Want := range []float64{258, 772} {
			ab88Field := ab88Get(ab88T, 11, ab88Dir, fmt.Sprintf("%s_fan_%d", ab88Prefix, ab88Index+1))
			ab88Position := *ab88Field.Ab88LogicalPosition
			if ab88Got, ab88OK := ab88Evaluate(ab88Field, ab88Logical[ab88Position:ab88Position+*ab88Field.Ab88Width]); !ab88OK || ab88Got != ab88Want {
				ab88T.Fatal("ab88 fan identity")
			}
		}
	}
	ab88Relay := ab88Get(ab88T, 11, "request", "relay")
	if ab88Relay.Ab88Bits["2"] != nil || ab88Relay.Ab88ReservedBits[0] != 2 {
		ab88T.Fatal("ab88 bit 2 reserved")
	}
	ab88Feedback := ab88Get(ab88T, 12, "positive", "feedback")
	for ab88Raw, ab88Want := range map[int]string{1: "accepted", 2: "rejected"} {
		ab88Label := ab88Feedback.Ab88Enumeration[strconv.Itoa(ab88Raw)]
		if ab88Label == "" && ab88Feedback.Ab88Otherwise != nil {
			ab88Label = *ab88Feedback.Ab88Otherwise
		}
		if ab88Label != ab88Want {
			ab88T.Fatal("ab88 feedback")
		}
	}
	if ab88Get(ab88T, 13, "positive", "digital").Ab88Enumeration["2"] != "" {
		ab88T.Fatal("ab88 digital 2 unlabeled")
	}
}

func TestAb88ExactIdentity(ab88T *testing.T) {
	ab88Control := ab88Get(ab88T, 12, "request", "minimum_pwm")
	for _, ab88B524 := range []bool{false, true} {
		ab88Field := ab88Control
		if ab88B524 {
			ab88Field.Ab88Key.Ab88Variant = "B524"
			ab88Field.Ab88Key.Ab88Primary = 0xb524
			ab88Field.Ab88Key.Ab88Opcode = Ab88Axis{true, 1}
			ab88Field.Ab88Key.Ab88Group = Ab88Axis{true, 2}
			ab88Field.Ab88Key.Ab88Instance = Ab88Axis{true, 3}
			ab88Field.Ab88Key.Ab88Register = Ab88Axis{true, 4}
		}
		ab88Key := ab88Field.Ab88Key
		for ab88Index := 0; ab88Index < reflect.TypeOf(ab88Key).NumField(); ab88Index++ {
			ab88T.Run(fmt.Sprintf("ab88_%t_%s", ab88B524, reflect.TypeOf(ab88Key).Field(ab88Index).Name), func(ab88T *testing.T) {
				ab88Changed := ab88Key
				ab88Axis := reflect.ValueOf(&ab88Changed).Elem().Field(ab88Index)
				switch ab88Axis.Kind() {
				case reflect.String:
					ab88Axis.SetString(ab88Axis.String() + "_ab88_other")
				case reflect.Uint16:
					ab88Axis.SetUint(ab88Axis.Uint() + 1)
				case reflect.Struct:
					ab88Axis.FieldByName("Ab88Value").SetUint(ab88Axis.FieldByName("Ab88Value").Uint() + 1)
				}
				if _, ab88OK := ab88FindAirField([]Ab88AirField{ab88Field}, ab88Changed); ab88OK {
					ab88T.Fatal("ab88 axis ignored")
				}
				if !ab88B524 {
					if _, ab88OK := Ab88LookupB512AirField(ab88Changed); ab88OK {
						ab88T.Fatal("ab88 public wildcard")
					}
				}
				if ab88Axis.Kind() == reflect.Struct {
					ab88Changed = ab88Key
					ab88Axis = reflect.ValueOf(&ab88Changed).Elem().Field(ab88Index)
					ab88Axis.FieldByName("Ab88Present").SetBool(!ab88Axis.FieldByName("Ab88Present").Bool())
					if _, ab88OK := ab88FindAirField([]Ab88AirField{ab88Field}, ab88Changed); ab88OK {
						ab88T.Fatal("ab88 absent equals present")
					}
					ab88Changed = ab88Key
					ab88Axis = reflect.ValueOf(&ab88Changed).Elem().Field(ab88Index)
					ab88Axis.FieldByName("Ab88Value").SetUint(0)
					ab88Axis.FieldByName("Ab88Present").SetBool(true)
					ab88Zero := ab88Field
					ab88Zero.Ab88Key = ab88Changed
					if _, ab88OK := ab88FindAirField([]Ab88AirField{ab88Zero}, ab88Changed); !ab88OK {
						ab88T.Fatal("ab88 explicit zero control")
					}
					ab88Axis.FieldByName("Ab88Present").SetBool(false)
					if _, ab88OK := ab88FindAirField([]Ab88AirField{ab88Zero}, ab88Changed); ab88OK {
						ab88T.Fatal("ab88 zero presence ignored")
					}
				}
				if _, ab88OK := ab88FindAirField([]Ab88AirField{ab88Field}, ab88Key); !ab88OK {
					ab88T.Fatal("ab88 valid control")
				}
			})
		}
	}
}

func ab88Mutate(ab88Field *Ab88AirField) {
	*ab88Field.Ab88LogicalPosition = 99
	*ab88Field.Ab88Width = 99
	if ab88Field.Ab88Unit != nil {
		*ab88Field.Ab88Unit = "ab88_changed"
	}
	if ab88Field.Ab88Default != nil {
		*ab88Field.Ab88Default = 99
	}
	if ab88Field.Ab88DefaultLabel != nil {
		*ab88Field.Ab88DefaultLabel = "ab88_changed"
	}
	if ab88Field.Ab88Otherwise != nil {
		*ab88Field.Ab88Otherwise = "ab88_changed"
	}
	for ab88Index := range ab88Field.Ab88Bounds {
		ab88Field.Ab88Bounds[ab88Index] = 99
	}
	for ab88Key := range ab88Field.Ab88Enumeration {
		ab88Field.Ab88Enumeration[ab88Key] = "ab88_changed"
	}
	for _, ab88Bits := range ab88Field.Ab88Bits {
		for ab88Key := range ab88Bits {
			ab88Bits[ab88Key] = "ab88_changed"
		}
	}
	for ab88Index := range ab88Field.Ab88ReservedBits {
		ab88Field.Ab88ReservedBits[ab88Index] = 99
	}
}

func TestAb88CopyIsolation(ab88T *testing.T) {
	ab88Before := Ab88B512AirFields()
	ab88Copy := Ab88B512AirFields()
	ab88Copy.Ab88Limitations[0] = "ab88_changed"
	ab88Copy.Ab88Operations[0].Ab88Grade = "ab88_changed"
	for ab88Index := range ab88Copy.Ab88Fields {
		ab88Mutate(&ab88Copy.Ab88Fields[ab88Index])
		ab88Lookup, ab88OK := Ab88LookupB512AirField(ab88Before.Ab88Fields[ab88Index].Ab88Key)
		if !ab88OK {
			ab88T.Fatal("ab88 missing copy")
		}
		ab88Mutate(&ab88Lookup)
		ab88Again, ab88OK := Ab88LookupB512AirField(ab88Before.Ab88Fields[ab88Index].Ab88Key)
		if !ab88OK || !reflect.DeepEqual(ab88Again, ab88Before.Ab88Fields[ab88Index]) {
			ab88T.Fatal("ab88 lookup shares storage")
		}
	}
	if !reflect.DeepEqual(Ab88B512AirFields(), ab88Before) {
		ab88T.Fatal("ab88 accessor shares storage")
	}
}
