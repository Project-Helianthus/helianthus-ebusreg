package system

import (
	ab44Binary "encoding/binary"
	ab44JSON "encoding/json"
	ab44Math "math"
	ab44OS "os"
	ab44Reflect "reflect"
	ab44Testing "testing"
)

func ab44Ptr[ab44T any](ab44Value ab44T) *ab44T { return &ab44Value }

func ab44TestContext() Ab44Context {
	return Ab44Context{
		Ab44Family: ab44Ptr("ab44-family-a"), Ab44Software: ab44Ptr(uint16(1)), Ab44Hardware: ab44Ptr(uint16(2)),
		Ab44APIVersion: ab44Ptr(byte(1)), Ab44APIRevision: ab44Ptr(byte(0)),
		Ab44OP: ab44Ptr(byte(1)), Ab44GG: ab44Ptr(byte(2)), Ab44II: ab44Ptr(byte(1)), Ab44RR: ab44Ptr(uint16(2)),
	}
}

func ab44TestEXP(ab44Minimum, ab44Maximum, ab44Step float32) []byte {
	ab44Raw := []byte{2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ab44Binary.LittleEndian.PutUint32(ab44Raw[3:7], ab44Math.Float32bits(ab44Minimum))
	ab44Binary.LittleEndian.PutUint32(ab44Raw[7:11], ab44Math.Float32bits(ab44Maximum))
	ab44Binary.LittleEndian.PutUint32(ab44Raw[11:15], ab44Math.Float32bits(ab44Step))
	return ab44Raw
}

func TestAb44Fixtures(ab44T *ab44Testing.T) {
	ab44Data, ab44Err := ab44OS.ReadFile("testdata/ab44_b524_description_spans.json")
	if ab44Err != nil {
		ab44T.Fatal(ab44Err)
	}
	var ab44Fixtures []struct {
		Ab44ID      string
		Ab44Codec   Ab44Codec
		Ab44Width   int
		Ab44RR      uint16
		Ab44Raw     []byte
		Ab44Minimum float64
		Ab44Maximum float64
		Ab44Step    float64
		Ab44Numeric bool
	}
	if ab44Err = ab44JSON.Unmarshal(ab44Data, &ab44Fixtures); ab44Err != nil {
		ab44T.Fatal(ab44Err)
	}
	if len(ab44Fixtures) != 6 {
		ab44T.Fatal("ab44 fixture count", len(ab44Fixtures))
	}
	ab44Seen := map[string]bool{}
	for _, ab44Fixture := range ab44Fixtures {
		if ab44Seen[ab44Fixture.Ab44ID] {
			ab44T.Fatal("ab44 duplicate fixture")
		}
		ab44Seen[ab44Fixture.Ab44ID] = true
		ab44T.Run(ab44Fixture.Ab44ID, func(ab44T *ab44Testing.T) {
			ab44Context := ab44TestContext()
			ab44Context.Ab44RR = ab44Ptr(ab44Fixture.Ab44RR)
			for _, ab44Grade := range []Ab44Grade{Ab44Hypothesis, Ab44Proven} {
				ab44Result := Ab44DecodeDescription(ab44Context, ab44Fixture.Ab44Codec, ab44Fixture.Ab44Width, ab44Grade, ab44Fixture.Ab44Raw)
				if !ab44Result.Ab44ImplementationValid || ab44Result.Ab44NumericValid != ab44Fixture.Ab44Numeric || ab44Result.Ab44Grade != ab44Grade || !ab44Reflect.DeepEqual(ab44Result.Ab44Raw, ab44Fixture.Ab44Raw) || !ab44Reflect.DeepEqual(ab44Result.Ab44Context, ab44Context) || ab44Result.Ab44MetadataOP == nil || *ab44Result.Ab44MetadataOP != 2 {
					ab44T.Fatalf("ab44 unexpected result: %+v", ab44Result)
				}
				if !ab44Reflect.DeepEqual(ab44Result.Ab44MinimumRaw, ab44Fixture.Ab44Raw[3:3+ab44Fixture.Ab44Width]) || !ab44Reflect.DeepEqual(ab44Result.Ab44MaximumRaw, ab44Fixture.Ab44Raw[3+ab44Fixture.Ab44Width:3+2*ab44Fixture.Ab44Width]) || !ab44Reflect.DeepEqual(ab44Result.Ab44StepRaw, ab44Fixture.Ab44Raw[3+2*ab44Fixture.Ab44Width:]) {
					ab44T.Fatal("ab44 spans")
				}
				if ab44Fixture.Ab44Numeric {
					ab44Expected := Ab44Limits{ab44Fixture.Ab44Minimum, ab44Fixture.Ab44Maximum, ab44Fixture.Ab44Step}
					if ab44Result.Ab44Limits == nil || *ab44Result.Ab44Limits != ab44Expected {
						ab44T.Fatalf("ab44 limits: %+v", ab44Result.Ab44Limits)
					}
				} else if ab44Result.Ab44Limits != nil || ab44Result.Ab44StepMeaning != Ab44Unknown || ab44Result.Ab44Reason != "ab44_date_time_raw_only" {
					ab44T.Fatal("ab44 date/time numeric permission")
				}
			}
		})
	}
}

func TestAb44Build(ab44T *ab44Testing.T) {
	ab44Context := ab44TestContext()
	ab44Request, ab44Err := Ab44BuildDescription(ab44Context)
	if ab44Err != nil || !ab44Reflect.DeepEqual(ab44Request, []byte{1, 2, 1, 2, 0}) {
		ab44T.Fatal("ab44 request", ab44Request, ab44Err)
	}
	*ab44Context.Ab44OP = 7
	*ab44Context.Ab44RR = 0x1234
	ab44Request, ab44Err = Ab44BuildDescription(ab44Context)
	if ab44Err != nil || !ab44Reflect.DeepEqual(ab44Request, []byte{7, 2, 1, 0x34, 0x12}) {
		ab44T.Fatal("ab44 remote high register", ab44Request, ab44Err)
	}
	for ab44OP := 0; ab44OP < 256; ab44OP++ {
		*ab44Context.Ab44OP = byte(ab44OP)
		ab44Request, ab44Err = Ab44BuildDescription(ab44Context)
		if ab44OP == 1 || ab44OP == 7 {
			continue
		}
		if ab44Err == nil || ab44Request != nil {
			ab44T.Fatal("ab44 unsupported op", ab44OP)
		}
	}
}

func TestAb44ContextAxes(ab44T *ab44Testing.T) {
	ab44Axes := []string{"ab44Family", "ab44Software", "ab44Hardware", "ab44APIVersion", "ab44APIRevision", "ab44OP", "ab44GG", "ab44II", "ab44RR"}
	for ab44Index, ab44Axis := range ab44Axes {
		ab44T.Run(ab44Axis, func(ab44T *ab44Testing.T) {
			for ab44Variant := 0; ab44Variant < 3; ab44Variant++ {
				ab44Context := ab44TestContext()
				ab44Field := ab44Reflect.ValueOf(&ab44Context).Elem().Field(ab44Index)
				if ab44Variant == 0 {
					ab44Field.SetZero()
				} else if ab44Index == 0 {
					ab44Alternative := "ab44-family-b"
					if ab44Variant == 1 {
						ab44Alternative = ""
					}
					ab44Field.Set(ab44Reflect.ValueOf(ab44Ptr(ab44Alternative)))
				} else {
					ab44Alternative := uint64(9)
					if ab44Variant == 1 {
						ab44Alternative = 0
					}
					if ab44Index == 5 && ab44Variant == 2 {
						ab44Alternative = 7
					}
					ab44Field.Elem().SetUint(ab44Alternative)
				}
				ab44Raw := ab44TestEXP(-20, 20, 1)
				if ab44Context.Ab44GG != nil {
					ab44Raw[0] = *ab44Context.Ab44GG
				}
				if ab44Context.Ab44RR != nil {
					ab44Binary.LittleEndian.PutUint16(ab44Raw[1:3], *ab44Context.Ab44RR)
				}
				ab44ExpectedValid := ab44Variant != 0 && (ab44Index != 0 || ab44Variant != 1) && (ab44Index != 5 || ab44Variant != 1)
				ab44Result := Ab44DecodeDescription(ab44Context, Ab44EXP, 4, Ab44Hypothesis, ab44Raw)
				if !ab44Reflect.DeepEqual(ab44Result.Ab44Context, ab44Context) || ab44Result.Ab44NumericValid != ab44ExpectedValid {
					ab44T.Fatalf("ab44 context variant %d: %+v", ab44Variant, ab44Result)
				}
				ab44Request, ab44Err := Ab44BuildDescription(ab44Context)
				if (ab44Err == nil) != ab44ExpectedValid {
					ab44T.Fatal("ab44 build context", ab44Err)
				}
				if ab44ExpectedValid {
					if !ab44Reflect.DeepEqual(ab44Request, []byte{*ab44Context.Ab44OP, *ab44Context.Ab44GG, *ab44Context.Ab44II, byte(*ab44Context.Ab44RR), byte(*ab44Context.Ab44RR >> 8)}) {
						ab44T.Fatal("ab44 build identity")
					}
					ab44ExpectedOP := byte(2)
					if *ab44Context.Ab44OP == 7 {
						ab44ExpectedOP = 6
					}
					if ab44Result.Ab44MetadataOP == nil || *ab44Result.Ab44MetadataOP != ab44ExpectedOP {
						ab44T.Fatal("ab44 direction")
					}
				} else if ab44Result.Ab44Limits != nil || ab44Result.Ab44ImplementationValid || ab44Result.Ab44Reason == "" {
					ab44T.Fatal("ab44 incomplete context admitted")
				}
			}
		})
	}
}

func TestAb44RejectRawOnly(ab44T *ab44Testing.T) {
	ab44Raw := []byte{2, 2, 0, 0, 0, 4, 0, 1, 0}
	ab44Cases := []struct {
		Ab44ID     string
		Ab44Codec  Ab44Codec
		Ab44Width  int
		Ab44Grade  Ab44Grade
		Ab44Reason string
	}{
		{"ab44_absent_codec", "", 2, Ab44Proven, "ab44_unsupported_codec"},
		{"ab44_unknown_codec", "ab44_future", 2, Ab44Proven, "ab44_unsupported_codec"},
		{"ab44_wrong_codec", Ab44EXP, 2, Ab44Proven, "ab44_codec_width_mismatch"},
		{"ab44_uch_width", Ab44UCH, 2, Ab44Proven, "ab44_codec_width_mismatch"},
		{"ab44_wrong_width", Ab44UIN, 1, Ab44Proven, "ab44_malformed_length"},
		{"ab44_absent_width", Ab44UIN, 0, Ab44Proven, "ab44_bad_width"},
		{"ab44_negative_width", Ab44UIN, -1, Ab44Proven, "ab44_bad_width"},
		{"ab44_large_width", Ab44UIN, int(^uint(0) >> 1), Ab44Proven, "ab44_bad_width"},
		{"ab44_absent_grade", Ab44UIN, 2, "", "ab44_codec_grade_required"},
		{"ab44_unknown_grade", Ab44UIN, 2, Ab44Unknown, "ab44_codec_grade_required"},
		{"ab44_other_grade", Ab44UIN, 2, "ab44_other", "ab44_codec_grade_required"},
	}
	for _, ab44Case := range ab44Cases {
		ab44T.Run(ab44Case.Ab44ID, func(ab44T *ab44Testing.T) {
			ab44Result := Ab44DecodeDescription(ab44TestContext(), ab44Case.Ab44Codec, ab44Case.Ab44Width, ab44Case.Ab44Grade, ab44Raw)
			if ab44Result.Ab44ImplementationValid || ab44Result.Ab44NumericValid || ab44Result.Ab44Limits != nil || ab44Result.Ab44Reason != ab44Case.Ab44Reason || ab44Result.Ab44Codec != ab44Case.Ab44Codec || ab44Result.Ab44Width != ab44Case.Ab44Width || ab44Result.Ab44Grade != ab44Case.Ab44Grade || !ab44Reflect.DeepEqual(ab44Result.Ab44Raw, ab44Raw) {
				ab44T.Fatalf("ab44 rejection: %+v", ab44Result)
			}
			if ab44Case.Ab44Width == 2 && (!ab44Reflect.DeepEqual(ab44Result.Ab44MinimumRaw, []byte{0, 0}) || !ab44Reflect.DeepEqual(ab44Result.Ab44MaximumRaw, []byte{4, 0}) || !ab44Reflect.DeepEqual(ab44Result.Ab44StepRaw, []byte{1, 0})) {
				ab44T.Fatal("ab44 rejected interpretation lost raw spans")
			}
		})
	}
	for _, ab44Codec := range []Ab44Codec{Ab44UCH, Ab44UIN, Ab44EXP, Ab44Date, Ab44Time} {
		for ab44Width := 1; ab44Width <= 4; ab44Width++ {
			ab44ExpectedWidth := map[Ab44Codec]int{Ab44UCH: 1, Ab44UIN: 2, Ab44EXP: 4, Ab44Date: 3, Ab44Time: 3}[ab44Codec]
			if ab44Width == ab44ExpectedWidth {
				continue
			}
			ab44Wrong := make([]byte, 3+3*ab44Width)
			copy(ab44Wrong, []byte{2, 2, 0})
			ab44Result := Ab44DecodeDescription(ab44TestContext(), ab44Codec, ab44Width, Ab44Proven, ab44Wrong)
			if ab44Result.Ab44Reason != "ab44_codec_width_mismatch" || ab44Result.Ab44ImplementationValid {
				ab44T.Fatal("ab44 codec width inferred", ab44Codec, ab44Width)
			}
		}
	}
}

func TestAb44MalformedAndEcho(ab44T *ab44Testing.T) {
	ab44Raw := ab44TestEXP(-20, 20, 1)
	for ab44Index := 0; ab44Index < 3; ab44Index++ {
		ab44Changed := append([]byte(nil), ab44Raw...)
		ab44Changed[ab44Index] ^= 1
		ab44Result := Ab44DecodeDescription(ab44TestContext(), Ab44EXP, 4, Ab44Proven, ab44Changed)
		if ab44Result.Ab44Reason != "ab44_echo_mismatch" || ab44Result.Ab44NumericValid || ab44Result.Ab44Limits != nil || !ab44Reflect.DeepEqual(ab44Result.Ab44Raw, ab44Changed) || !ab44Reflect.DeepEqual(ab44Result.Ab44StepRaw, ab44Raw[11:]) {
			ab44T.Fatal("ab44 echo", ab44Index, ab44Result)
		}
	}
	for ab44Length := 0; ab44Length <= 30; ab44Length++ {
		if ab44Length == 15 {
			continue
		}
		ab44Malformed := make([]byte, ab44Length)
		copy(ab44Malformed, ab44Raw)
		ab44Result := Ab44DecodeDescription(ab44TestContext(), Ab44EXP, 4, Ab44Proven, ab44Malformed)
		if ab44Result.Ab44Reason != "ab44_malformed_length" || ab44Result.Ab44ImplementationValid || ab44Result.Ab44Limits != nil || len(ab44Result.Ab44Raw) != ab44Length {
			ab44T.Fatal("ab44 malformed", ab44Length)
		}
	}
	ab44Duplicate := append(append([]byte(nil), ab44Raw...), ab44Raw...)
	ab44Legacy := []byte{1, 2, 2, 0}
	for _, ab44Malformed := range [][]byte{nil, ab44Duplicate, ab44Legacy} {
		ab44Result := Ab44DecodeDescription(ab44TestContext(), Ab44EXP, 4, Ab44Proven, ab44Malformed)
		if ab44Result.Ab44ImplementationValid || ab44Result.Ab44Reason != "ab44_malformed_length" {
			ab44T.Fatal("ab44 duplicate/legacy")
		}
	}
	for _, ab44OP := range []byte{0, 2, 6, 8} {
		ab44Context := ab44TestContext()
		*ab44Context.Ab44OP = ab44OP
		ab44Result := Ab44DecodeDescription(ab44Context, Ab44EXP, 4, Ab44Proven, ab44Raw)
		if ab44Result.Ab44Reason != "ab44_unsupported_op" || ab44Result.Ab44ImplementationValid {
			ab44T.Fatal("ab44 reply direction", ab44OP)
		}
	}
}

func TestAb44NumericSpans(ab44T *ab44Testing.T) {
	for ab44Index, ab44Expected := range []Ab44Limits{{-19, 20, 1}, {-20, 21, 1}, {-20, 20, 2}, {-20, 20, 0}} {
		ab44Raw := ab44TestEXP(float32(ab44Expected.Ab44Minimum), float32(ab44Expected.Ab44Maximum), float32(ab44Expected.Ab44Step))
		ab44Result := Ab44DecodeDescription(ab44TestContext(), Ab44EXP, 4, Ab44Hypothesis, ab44Raw)
		if !ab44Result.Ab44NumericValid || ab44Result.Ab44Limits == nil || *ab44Result.Ab44Limits != ab44Expected {
			ab44T.Fatal("ab44 independent span", ab44Index)
		}
	}
	ab44Cases := []struct {
		Ab44Raw    []byte
		Ab44Reason string
	}{
		{ab44TestEXP(-20, 20, -1), "ab44_negative_step"},
		{ab44TestEXP(21, 20, 1), "ab44_reversed_limits"},
	}
	for ab44Index := 0; ab44Index < 3; ab44Index++ {
		for _, ab44Nonfinite := range []float32{float32(ab44Math.NaN()), float32(ab44Math.Inf(1)), float32(ab44Math.Inf(-1))} {
			ab44Raw := ab44TestEXP(-20, 20, 1)
			ab44Binary.LittleEndian.PutUint32(ab44Raw[3+4*ab44Index:], ab44Math.Float32bits(ab44Nonfinite))
			ab44Cases = append(ab44Cases, struct {
				Ab44Raw    []byte
				Ab44Reason string
			}{ab44Raw, "ab44_nonfinite_limits"})
		}
	}
	for _, ab44Case := range ab44Cases {
		ab44Result := Ab44DecodeDescription(ab44TestContext(), Ab44EXP, 4, Ab44Proven, ab44Case.Ab44Raw)
		if ab44Result.Ab44NumericValid || ab44Result.Ab44ImplementationValid || ab44Result.Ab44Limits != nil || ab44Result.Ab44Grade != Ab44Proven || ab44Result.Ab44Reason != ab44Case.Ab44Reason || !ab44Reflect.DeepEqual(ab44Result.Ab44MinimumRaw, ab44Case.Ab44Raw[3:7]) || !ab44Reflect.DeepEqual(ab44Result.Ab44MaximumRaw, ab44Case.Ab44Raw[7:11]) || !ab44Reflect.DeepEqual(ab44Result.Ab44StepRaw, ab44Case.Ab44Raw[11:15]) {
			ab44T.Fatalf("ab44 invalid numeric: %+v", ab44Result)
		}
	}
	for _, ab44Codec := range []Ab44Codec{Ab44Date, Ab44Time} {
		for ab44Index := 3; ab44Index < 12; ab44Index++ {
			ab44Raw := []byte{2, 2, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
			ab44Raw[ab44Index] = 255
			ab44Result := Ab44DecodeDescription(ab44TestContext(), ab44Codec, 3, Ab44Proven, ab44Raw)
			if !ab44Result.Ab44ImplementationValid || ab44Result.Ab44NumericValid || ab44Result.Ab44Limits != nil || ab44Result.Ab44StepMeaning != Ab44Unknown || !ab44Reflect.DeepEqual(ab44Result.Ab44StepRaw, ab44Raw[9:]) {
				ab44T.Fatal("ab44 date raw mutation")
			}
		}
	}
}

func TestAb44Ownership(ab44T *ab44Testing.T) {
	ab44Context := ab44TestContext()
	ab44ExpectedContext := ab44TestContext()
	ab44Raw := ab44TestEXP(-20, 20, 1)
	ab44ExpectedRaw := ab44TestEXP(-20, 20, 1)
	ab44Result := Ab44DecodeDescription(ab44Context, Ab44EXP, 4, Ab44Hypothesis, ab44Raw)
	for ab44Index := range ab44Raw {
		ab44Raw[ab44Index] = 255
	}
	for ab44Index := 0; ab44Index < 9; ab44Index++ {
		ab44Reflect.ValueOf(ab44Context).Field(ab44Index).Elem().SetZero()
	}
	if !ab44Reflect.DeepEqual(ab44Result.Ab44Context, ab44ExpectedContext) || !ab44Reflect.DeepEqual(ab44Result.Ab44Raw, ab44ExpectedRaw) || !ab44Reflect.DeepEqual(ab44Result.Ab44MinimumRaw, ab44ExpectedRaw[3:7]) {
		ab44T.Fatal("ab44 input alias")
	}
	ab44Result.Ab44Raw[3] = 255
	if !ab44Reflect.DeepEqual(ab44Result.Ab44MinimumRaw, ab44ExpectedRaw[3:7]) {
		ab44T.Fatal("ab44 raw aliases spans")
	}
	ab44Result.Ab44MinimumRaw[0] = 254
	ab44Result.Ab44MaximumRaw[0] = 253
	ab44Result.Ab44StepRaw[0] = 252
	*ab44Result.Ab44MetadataOP = 99
	ab44Result.Ab44Limits.Ab44Minimum = 99
	for ab44Index := 0; ab44Index < 9; ab44Index++ {
		ab44Reflect.ValueOf(ab44Result.Ab44Context).Field(ab44Index).Elem().SetZero()
	}
	ab44Next := Ab44DecodeDescription(ab44ExpectedContext, Ab44EXP, 4, Ab44Hypothesis, ab44ExpectedRaw)
	if !ab44Reflect.DeepEqual(ab44Next.Ab44Context, ab44ExpectedContext) || !ab44Reflect.DeepEqual(ab44Next.Ab44Raw, ab44ExpectedRaw) || !ab44Reflect.DeepEqual(ab44Next.Ab44MinimumRaw, ab44ExpectedRaw[3:7]) || !ab44Reflect.DeepEqual(ab44Next.Ab44MaximumRaw, ab44ExpectedRaw[7:11]) || !ab44Reflect.DeepEqual(ab44Next.Ab44StepRaw, ab44ExpectedRaw[11:]) || *ab44Next.Ab44Limits != (Ab44Limits{-20, 20, 1}) || *ab44Next.Ab44MetadataOP != 2 {
		ab44T.Fatal("ab44 later result alias")
	}
	ab44Request, ab44Err := Ab44BuildDescription(ab44ExpectedContext)
	if ab44Err != nil {
		ab44T.Fatal(ab44Err)
	}
	ab44Request[0] = 99
	ab44NextRequest, ab44Err := Ab44BuildDescription(ab44ExpectedContext)
	if ab44Err != nil || !ab44Reflect.DeepEqual(ab44NextRequest, []byte{1, 2, 1, 2, 0}) {
		ab44T.Fatal("ab44 request alias")
	}
}

func TestAb44ZeroAndRawOwnership(ab44T *ab44Testing.T) {
	for _, ab44Case := range []struct {
		Ab44Codec Ab44Codec
		Ab44Width int
		Ab44Raw   []byte
	}{
		{Ab44UCH, 1, []byte{2, 2, 0, 0, 0, 0}},
		{Ab44UIN, 2, []byte{2, 2, 0, 0, 0, 0, 0, 0, 0}},
		{Ab44EXP, 4, []byte{2, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
	} {
		ab44Result := Ab44DecodeDescription(ab44TestContext(), ab44Case.Ab44Codec, ab44Case.Ab44Width, Ab44Hypothesis, ab44Case.Ab44Raw)
		if !ab44Result.Ab44NumericValid || ab44Result.Ab44Limits == nil || *ab44Result.Ab44Limits != (Ab44Limits{0, 0, 0}) {
			ab44T.Fatal("ab44 literal zero/equal bounds", ab44Case.Ab44Codec)
		}
	}
	for _, ab44Grade := range []Ab44Grade{"", Ab44Unknown} {
		ab44Context := ab44TestContext()
		ab44Raw := []byte{2, 2, 0, 0, 0, 4, 0, 1, 0}
		ab44Result := Ab44DecodeDescription(ab44Context, Ab44UIN, 2, ab44Grade, ab44Raw)
		ab44Raw[3] = 99
		*ab44Context.Ab44II = 99
		if ab44Result.Ab44Raw[3] != 0 || *ab44Result.Ab44Context.Ab44II != 1 || !ab44Reflect.DeepEqual(ab44Result.Ab44MinimumRaw, []byte{0, 0}) {
			ab44T.Fatal("ab44 raw-only input alias")
		}
		ab44Result.Ab44MinimumRaw[0] = 88
		if ab44Result.Ab44Raw[3] != 0 || ab44Result.Ab44MaximumRaw[0] != 4 || ab44Result.Ab44StepRaw[0] != 1 {
			ab44T.Fatal("ab44 raw-only span alias")
		}
	}
}
