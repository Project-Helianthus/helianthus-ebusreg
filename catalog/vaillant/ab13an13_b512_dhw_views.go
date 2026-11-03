package vaillant

const (
	Ab13An13Family  = "vaillant-dhw-family"
	Ab13An13Variant = "b512-op06-candidate"
)

type Ab13An13Scalar struct {
	Ab13An13Raw          byte
	Ab13An13Value        float64
	Ab13An13Unit         string
	Ab13An13UnitConflict bool
	Ab13An13NamedDomain  bool
}

type Ab13An13Request struct {
	Ab13An13TaskActive          *bool
	Ab13An13HeatDemand          *bool
	Ab13An13LoadingActive       *bool
	Ab13An13LoadingPumpFault    *bool
	Ab13An13SensorConditionBits [4]*bool
	Ab13An13Modulation          *Ab13An13Scalar
	Ab13An13Flow                *Ab13An13Scalar
	Ab13An13Temperatures        [3]*Ab13An13Scalar
	Ab13An13InletConditionBits  [2]*bool
}

type Ab13An13Positive struct {
	Ab13An13LoadingDisabled       *bool
	Ab13An13InternalTestActive    *bool
	Ab13An13Ready                 *bool
	Ab13An13FrostProtectionActive *bool
	Ab13An13LegionellaDemand      *bool
	Ab13An13Reset                 *bool
	Ab13An13Setpoint              *Ab13An13Scalar
	Ab13An13ReturnTemperature     *Ab13An13Scalar
	Ab13An13Adjustment            *Ab13An13Scalar
}

type Ab13An13Result struct {
	Ab13An13Raw           []byte
	Ab13An13Unmodelled    map[int]byte
	Ab13An13Admitted      bool
	Ab13An13Complete      bool
	Ab13An13ContractGrade string
	Ab13An13FieldGrade    string
	Ab13An13Qualification string
	Ab13An13Presence      string
	Ab13An13Freshness     string
	Ab13An13Execution     string
	Ab13An13Defaults      map[string]float64
	Ab13An13Request       *Ab13An13Request
	Ab13An13Positive      *Ab13An13Positive
}

func ab13an13Bit(ab13an13Raw byte, ab13an13Index uint) *bool {
	ab13an13Value := ab13an13Raw&(1<<ab13an13Index) != 0
	return &ab13an13Value
}

func ab13an13Scalar(ab13an13Raw byte, ab13an13Divisor float64, ab13an13Unit string) *Ab13An13Scalar {
	return &Ab13An13Scalar{Ab13An13Raw: ab13an13Raw, Ab13An13Value: float64(ab13an13Raw) / ab13an13Divisor, Ab13An13Unit: ab13an13Unit}
}

func Ab13An13Interpret(ab13an13Family, ab13an13Variant, ab13an13Direction string, ab13an13Image []byte) Ab13An13Result {
	ab13an13Result := Ab13An13Result{
		Ab13An13Raw:           append([]byte(nil), ab13an13Image...),
		Ab13An13Unmodelled:    make(map[int]byte),
		Ab13An13ContractGrade: "Proven", Ab13An13FieldGrade: "Hypothesis",
		Ab13An13Qualification: "Unknown", Ab13An13Presence: "Unknown",
		Ab13An13Freshness: "Unknown", Ab13An13Execution: "Unknown",
	}
	for ab13an13Index, ab13an13Raw := range ab13an13Image {
		ab13an13Result.Ab13An13Unmodelled[ab13an13Index] = ab13an13Raw
	}
	if ab13an13Family != Ab13An13Family || ab13an13Variant != Ab13An13Variant ||
		(ab13an13Direction != "request" && ab13an13Direction != "positive") ||
		len(ab13an13Image) < 3 || ab13an13Image[0] != 0xb5 || ab13an13Image[1] != 0x12 || ab13an13Image[2] != 0x06 {
		return ab13an13Result
	}
	ab13an13Result.Ab13An13Admitted = true
	for ab13an13Index := 0; ab13an13Index < 3; ab13an13Index++ {
		delete(ab13an13Result.Ab13An13Unmodelled, ab13an13Index)
	}
	if ab13an13Direction == "request" {
		ab13an13Request := &Ab13An13Request{}
		ab13an13Result.Ab13An13Request = ab13an13Request
		ab13an13Result.Ab13An13Complete = len(ab13an13Image) >= 10
		ab13an13Result.Ab13An13Defaults = map[string]float64{"Ab13An13Flags": 0, "Ab13An13Scalars": 0}
		if len(ab13an13Image) > 3 {
			ab13an13Request.Ab13An13TaskActive = ab13an13Bit(ab13an13Image[3], 0)
			ab13an13Request.Ab13An13HeatDemand = ab13an13Bit(ab13an13Image[3], 1)
			ab13an13Request.Ab13An13LoadingActive = ab13an13Bit(ab13an13Image[3], 2)
			ab13an13Request.Ab13An13LoadingPumpFault = ab13an13Bit(ab13an13Image[3], 3)
			for ab13an13Index := range ab13an13Request.Ab13An13SensorConditionBits {
				ab13an13Request.Ab13An13SensorConditionBits[ab13an13Index] = ab13an13Bit(ab13an13Image[3], uint(ab13an13Index+4))
			}
		}
		if len(ab13an13Image) > 4 {
			ab13an13Request.Ab13An13Modulation = ab13an13Scalar(ab13an13Image[4], 2, "formal Celsius; modulation interpretation")
			ab13an13Request.Ab13An13Modulation.Ab13An13UnitConflict = true
		}
		if len(ab13an13Image) > 5 {
			ab13an13Request.Ab13An13Flow = ab13an13Scalar(ab13an13Image[5], 10, "litres/minute")
		}
		for ab13an13Index := range ab13an13Request.Ab13An13Temperatures {
			if len(ab13an13Image) > 6+ab13an13Index {
				ab13an13Request.Ab13An13Temperatures[ab13an13Index] = ab13an13Scalar(ab13an13Image[6+ab13an13Index], 2, "Celsius")
			}
		}
		if len(ab13an13Image) > 9 {
			for ab13an13Index := range ab13an13Request.Ab13An13InletConditionBits {
				ab13an13Request.Ab13An13InletConditionBits[ab13an13Index] = ab13an13Bit(ab13an13Image[9], uint(ab13an13Index))
			}
		}
		for ab13an13Index := 3; ab13an13Index < len(ab13an13Image) && ab13an13Index < 10; ab13an13Index++ {
			delete(ab13an13Result.Ab13An13Unmodelled, ab13an13Index)
		}
		if len(ab13an13Image) > 9 && ab13an13Image[9]&0xfc != 0 {
			ab13an13Result.Ab13An13Unmodelled[9] = ab13an13Image[9] & 0xfc
		}
	} else {
		ab13an13Positive := &Ab13An13Positive{}
		ab13an13Result.Ab13An13Positive = ab13an13Positive
		ab13an13Result.Ab13An13Complete = len(ab13an13Image) >= 7
		ab13an13Result.Ab13An13Defaults = map[string]float64{"Ab13An13Flags": 0, "Ab13An13Scalars": 0, "Ab13An13Adjustment": -1}
		if len(ab13an13Image) > 3 {
			ab13an13Positive.Ab13An13LoadingDisabled = ab13an13Bit(ab13an13Image[3], 0)
			ab13an13Positive.Ab13An13InternalTestActive = ab13an13Bit(ab13an13Image[3], 1)
			ab13an13Positive.Ab13An13Ready = ab13an13Bit(ab13an13Image[3], 2)
			ab13an13Positive.Ab13An13FrostProtectionActive = ab13an13Bit(ab13an13Image[3], 3)
			ab13an13Positive.Ab13An13LegionellaDemand = ab13an13Bit(ab13an13Image[3], 4)
			ab13an13Positive.Ab13An13Reset = ab13an13Bit(ab13an13Image[3], 5)
		}
		if len(ab13an13Image) > 4 {
			ab13an13Positive.Ab13An13Setpoint = ab13an13Scalar(ab13an13Image[4], 2, "Celsius")
		}
		if len(ab13an13Image) > 5 {
			ab13an13Positive.Ab13An13ReturnTemperature = ab13an13Scalar(ab13an13Image[5], 2, "Celsius")
		}
		if len(ab13an13Image) > 6 {
			ab13an13Value := int8(ab13an13Image[6])
			ab13an13Positive.Ab13An13Adjustment = &Ab13An13Scalar{Ab13An13Raw: ab13an13Image[6], Ab13An13Value: float64(ab13an13Value), Ab13An13NamedDomain: ab13an13Value >= -1 && ab13an13Value <= 2}
		}
		for ab13an13Index := 3; ab13an13Index < len(ab13an13Image) && ab13an13Index < 7; ab13an13Index++ {
			delete(ab13an13Result.Ab13An13Unmodelled, ab13an13Index)
		}
		if len(ab13an13Image) > 3 && ab13an13Image[3]&0xc0 != 0 {
			ab13an13Result.Ab13An13Unmodelled[3] = ab13an13Image[3] & 0xc0
		}
	}
	return ab13an13Result
}
