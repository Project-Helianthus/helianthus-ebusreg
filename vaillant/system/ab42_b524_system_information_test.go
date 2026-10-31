package system

import (
	ab42bytes "bytes"
	ab42binary "encoding/binary"
	ab42json "encoding/json"
	ab42math "math"
	ab42os "os"
	ab42reflect "reflect"
	ab42testing "testing"
)

func ab42TestPointer[Ab42Type any](ab42Value Ab42Type) *Ab42Type { return &ab42Value }

func TestAb42BuildInformation(ab42T *ab42testing.T) {
	for _, ab42Case := range []struct {
		ab42Identifier uint16
		ab42Expected   []byte
	}{
		{0x1234, []byte{0, 0x34, 0x12}},
		{0, []byte{0, 0, 0}},
		{1, []byte{0, 1, 0}},
		{0xffff, []byte{0, 0xff, 0xff}},
	} {
		ab42First := Ab42BuildInformation(ab42Case.ab42Identifier)
		if !ab42bytes.Equal(ab42First, ab42Case.ab42Expected) {
			ab42T.Fatalf("ab42 request: %x", ab42First)
		}
		ab42First[0] = 99
		if !ab42bytes.Equal(Ab42BuildInformation(ab42Case.ab42Identifier), ab42Case.ab42Expected) {
			ab42T.Fatal("ab42 request aliases prior result")
		}
	}
}

func TestAb42SyntheticInformation(ab42T *ab42testing.T) {
	ab42Data, ab42Err := ab42os.ReadFile("testdata/ab42_b524_system_information.json")
	if ab42Err != nil {
		ab42T.Fatal(ab42Err)
	}
	var ab42Cases []struct {
		Ab42ID            string
		Ab42Identifier    uint16
		Ab42Raw           []byte
		Ab42Value         float32
		Ab42Capacity      bool
		Ab42HeatPumpCount bool
		Ab42Meaning       string
	}
	ab42Decoder := ab42json.NewDecoder(ab42bytes.NewReader(ab42Data))
	ab42Decoder.DisallowUnknownFields()
	if ab42Err = ab42Decoder.Decode(&ab42Cases); ab42Err != nil {
		ab42T.Fatal(ab42Err)
	}
	if len(ab42Cases) != 7 {
		ab42T.Fatalf("ab42 fixture count: %d", len(ab42Cases))
	}
	ab42Seen := map[string]bool{}
	for _, ab42Case := range ab42Cases {
		if ab42Case.Ab42ID == "" || ab42Seen[ab42Case.Ab42ID] {
			ab42T.Fatal("ab42 missing or duplicate fixture")
		}
		ab42Seen[ab42Case.Ab42ID] = true
		ab42T.Run(ab42Case.Ab42ID, func(ab42T *ab42testing.T) {
			ab42Result := Ab42DecodeInformation(Ab42InformationContext{Ab42Identifier: ab42Case.Ab42Identifier}, ab42Case.Ab42Raw)
			if !ab42Result.Ab42Valid || ab42Result.Ab42Value != ab42Case.Ab42Value || ab42Result.Ab42InvalidReason != "" {
				ab42T.Fatalf("ab42 decode: %+v", ab42Result)
			}
			if ab42Result.Ab42Context.Ab42Identifier != ab42Case.Ab42Identifier || !ab42bytes.Equal(ab42Result.Ab42Raw, ab42Case.Ab42Raw) {
				ab42T.Fatal("ab42 evidence lost")
			}
			if ab42Result.Ab42Meaning != ab42Case.Ab42Meaning {
				ab42T.Fatalf("ab42 meaning: %s", ab42Result.Ab42Meaning)
			}
			if (ab42Result.Ab42Capacity != nil) != ab42Case.Ab42Capacity || (ab42Result.Ab42HeatPumpCount != nil) != ab42Case.Ab42HeatPumpCount {
				ab42T.Fatal("ab42 projection state")
			}
			if ab42Result.Ab42Capacity != nil && *ab42Result.Ab42Capacity != ab42Case.Ab42Value {
				ab42T.Fatal("ab42 capacity value")
			}
			if ab42Result.Ab42HeatPumpCount != nil && *ab42Result.Ab42HeatPumpCount != ab42Case.Ab42Value {
				ab42T.Fatal("ab42 count value")
			}
			if ab42Result.Ab42RepresentationGrade != "Proven" || ab42Result.Ab42LayoutGrade != "Hypothesis" || ab42Result.Ab42QualificationGrade != "Unknown" || ab42Result.Ab42PresenceGrade != "Unknown" || ab42Result.Ab42EligibilityGrade != "Unknown" {
				ab42T.Fatal("ab42 representation promoted to physical evidence")
			}
		})
	}
}

func TestAb42IndependentRepresentation(ab42T *ab42testing.T) {
	for _, ab42Case := range []struct {
		ab42Raw      []byte
		ab42Value    float32
		ab42Capacity bool
	}{
		{[]byte{0, 0, 128, 63}, 1, true},
		{[]byte{0, 0, 192, 63}, 1.5, false},
		{[]byte{0, 0, 128, 191}, -1, false},
		{[]byte{0, 0, 0, 0}, 0, true},
		{[]byte{0, 0, 0, 128}, 0, true},
	} {
		for _, ab42Identifier := range []uint16{0, 1, 0x000d, 0xffff, 0x1234} {
			ab42Result := Ab42DecodeInformation(Ab42InformationContext{Ab42Identifier: ab42Identifier}, ab42Case.ab42Raw)
			if !ab42Result.Ab42Valid || ab42Result.Ab42Value != ab42Case.ab42Value {
				ab42T.Fatal("ab42 independent numeric expectation")
			}
			if (ab42Result.Ab42Capacity != nil) != (ab42Case.ab42Capacity && ab42Identifier <= 1) {
				ab42T.Fatal("ab42 capacity is not an integral supported capacity")
			}
			if (ab42Result.Ab42HeatPumpCount != nil) != (ab42Case.ab42Capacity && ab42Identifier == 0x000d) {
				ab42T.Fatal("ab42 count identity conflated")
			}
			ab42Meaning := "Unknown"
			switch ab42Identifier {
			case 0:
				ab42Meaning = "supported circuit capacity"
			case 1:
				ab42Meaning = "supported zone capacity"
			case 0x000d:
				ab42Meaning = "heat-pump count"
			}
			if ab42Result.Ab42Meaning != ab42Meaning || ab42Result.Ab42PresenceGrade != "Unknown" || ab42Result.Ab42EligibilityGrade != "Unknown" {
				ab42T.Fatal("ab42 capacity/count became presence or eligibility")
			}
		}
	}
}

func TestAb42InvalidInformation(ab42T *ab42testing.T) {
	ab42Inputs := [][]byte{nil, {}, {0}, {0, 0}, {0, 0, 128}, {0, 0, 128, 63, 0}, {0, 0, 128, 63, 0, 0, 128, 63}}
	for _, ab42Bits := range []uint32{0x7fc00001, 0x7f800001, 0x7f800000, 0xff800000} {
		ab42Raw := make([]byte, 4)
		ab42binary.LittleEndian.PutUint32(ab42Raw, ab42Bits)
		ab42Inputs = append(ab42Inputs, ab42Raw)
	}
	for _, ab42Raw := range ab42Inputs {
		for _, ab42Identifier := range []uint16{0, 1, 0x000d, 0xffff} {
			ab42Context := Ab42InformationContext{Ab42Identifier: ab42Identifier, Ab42APIVersion: ab42TestPointer(uint16(2)), Ab42APIRevision: ab42TestPointer(uint16(7))}
			ab42Result := Ab42DecodeInformation(ab42Context, ab42Raw)
			if ab42Result.Ab42Valid || ab42Result.Ab42InvalidReason == "" || ab42Result.Ab42Capacity != nil || ab42Result.Ab42HeatPumpCount != nil || ab42Result.Ab42RepresentationGrade != "Unknown" {
				ab42T.Fatalf("ab42 invalid input admitted: %x", ab42Raw)
			}
			if !ab42bytes.Equal(ab42Result.Ab42Raw, ab42Raw) || !ab42reflect.DeepEqual(ab42Result.Ab42Context, ab42Context) {
				ab42T.Fatal("ab42 invalid response discarded independent evidence")
			}
			if len(ab42Raw) > 0 {
				ab42Expected := append([]byte(nil), ab42Raw...)
				ab42Raw[0] ^= 0xff
				if !ab42bytes.Equal(ab42Result.Ab42Raw, ab42Expected) {
					ab42T.Fatal("ab42 invalid raw aliases caller")
				}
				ab42Raw[0] ^= 0xff
				ab42Result.Ab42Raw[0] ^= 0xff
				if !ab42bytes.Equal(ab42Raw, ab42Expected) {
					ab42T.Fatal("ab42 caller raw aliases invalid result")
				}
			}
		}
	}
}

func ab42TestContext() Ab42InformationContext {
	return Ab42InformationContext{
		Ab42Family: ab42TestPointer("ab42-family"), Ab42Software: ab42TestPointer("ab42-software"), Ab42Hardware: ab42TestPointer("ab42-hardware"),
		Ab42APIVersion: ab42TestPointer(uint16(2)), Ab42APIRevision: ab42TestPointer(uint16(7)), Ab42Identifier: 0x1234,
	}
}

func TestAb42ContextAxes(ab42T *ab42testing.T) {
	for _, ab42Axis := range []string{"Ab42Family", "Ab42Software", "Ab42Hardware", "Ab42APIVersion", "Ab42APIRevision", "Ab42Identifier"} {
		ab42T.Run(ab42Axis, func(ab42T *ab42testing.T) {
			for _, ab42Alternative := range []string{"ab42-absent", "ab42-zero", "ab42-alternative"} {
				ab42Context := ab42TestContext()
				ab42Field := ab42reflect.ValueOf(&ab42Context).Elem().FieldByName(ab42Axis)
				if ab42Field.Kind() == ab42reflect.Uint16 {
					if ab42Alternative == "ab42-alternative" {
						ab42Field.SetUint(0xffff)
					} else {
						ab42Field.SetUint(0)
					}
				} else if ab42Alternative == "ab42-absent" {
					ab42Field.SetZero()
				} else {
					ab42Field.Set(ab42reflect.New(ab42Field.Type().Elem()))
					if ab42Alternative == "ab42-alternative" {
						if ab42Field.Elem().Kind() == ab42reflect.String {
							ab42Field.Elem().SetString("ab42-alternative")
						} else {
							ab42Field.Elem().SetUint(19)
						}
					}
				}
				ab42Result := Ab42DecodeInformation(ab42Context, []byte{0, 0, 128, 63})
				if !ab42Result.Ab42Valid || !ab42reflect.DeepEqual(ab42Result.Ab42Context, ab42Context) {
					ab42T.Fatalf("ab42 context axis %s %s lost", ab42Axis, ab42Alternative)
				}
				ab42Invalid := Ab42DecodeInformation(ab42Context, []byte{0})
				if ab42Invalid.Ab42Valid || !ab42reflect.DeepEqual(ab42Invalid.Ab42Context, ab42Context) {
					ab42T.Fatal("ab42 partial failure lost context")
				}
			}
		})
	}
}

func TestAb42Ownership(ab42T *ab42testing.T) {
	ab42Context := ab42TestContext()
	ab42Raw := []byte{0, 0, 128, 63}
	ab42First := Ab42DecodeInformation(ab42Context, ab42Raw)
	ab42Second := Ab42DecodeInformation(ab42Context, ab42Raw)
	for _, ab42Axis := range []string{"Ab42Family", "Ab42Software", "Ab42Hardware", "Ab42APIVersion", "Ab42APIRevision"} {
		ab42Source := ab42reflect.ValueOf(&ab42Context).Elem().FieldByName(ab42Axis).Elem()
		ab42Saved := ab42reflect.ValueOf(ab42First.Ab42Context).FieldByName(ab42Axis).Elem().Interface()
		if ab42Source.Kind() == ab42reflect.String {
			ab42Source.SetString("ab42-mutated")
		} else {
			ab42Source.SetUint(99)
		}
		if !ab42reflect.DeepEqual(ab42reflect.ValueOf(ab42First.Ab42Context).FieldByName(ab42Axis).Elem().Interface(), ab42Saved) {
			ab42T.Fatalf("ab42 caller context aliases %s", ab42Axis)
		}
		ab42Returned := ab42reflect.ValueOf(ab42First.Ab42Context).FieldByName(ab42Axis).Elem()
		if ab42Returned.Kind() == ab42reflect.String {
			ab42Returned.SetString("ab42-result")
		} else {
			ab42Returned.SetUint(100)
		}
		if !ab42reflect.DeepEqual(ab42reflect.ValueOf(ab42Second.Ab42Context).FieldByName(ab42Axis).Elem().Interface(), ab42Saved) {
			ab42T.Fatalf("ab42 returned contexts alias %s", ab42Axis)
		}
		if ab42reflect.DeepEqual(ab42Source.Interface(), ab42Returned.Interface()) {
			ab42T.Fatal("ab42 returned context aliases source")
		}
	}
	ab42Raw[0] = 1
	if !ab42bytes.Equal(ab42First.Ab42Raw, []byte{0, 0, 128, 63}) {
		ab42T.Fatal("ab42 caller raw aliases result")
	}
	ab42First.Ab42Raw[1] = 1
	if !ab42bytes.Equal(ab42Second.Ab42Raw, []byte{0, 0, 128, 63}) || ab42Raw[1] != 0 {
		ab42T.Fatal("ab42 returned raw aliases caller or later result")
	}
	for _, ab42Identifier := range []uint16{0, 0x000d} {
		ab42One := Ab42DecodeInformation(Ab42InformationContext{Ab42Identifier: ab42Identifier}, []byte{0, 0, 128, 63})
		ab42Two := Ab42DecodeInformation(Ab42InformationContext{Ab42Identifier: ab42Identifier}, []byte{0, 0, 128, 63})
		ab42Projection, ab42Other := ab42One.Ab42Capacity, ab42Two.Ab42Capacity
		if ab42Identifier == 0x000d {
			ab42Projection, ab42Other = ab42One.Ab42HeatPumpCount, ab42Two.Ab42HeatPumpCount
		}
		*ab42Projection = 9
		if *ab42Other != 1 || ab42One.Ab42Value != 1 {
			ab42T.Fatal("ab42 projection aliases independent state")
		}
	}
	ab42Later := Ab42DecodeInformation(ab42TestContext(), []byte{0, 0, 128, 63})
	if !ab42reflect.DeepEqual(ab42Later.Ab42Context, ab42TestContext()) || ab42math.Float32bits(ab42Later.Ab42Value) != 0x3f800000 {
		ab42T.Fatal("ab42 mutations changed later decode")
	}
}
