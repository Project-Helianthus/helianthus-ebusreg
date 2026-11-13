package vaillant

const (
	B512DHWViewFamily  = "vaillant-dhw-family"
	B512DHWViewVariant = "b512-op06-candidate"
)

type B512DHWViewScalar struct {
	B512DHWViewRaw          byte    `json:"Ab13An13Raw"`
	B512DHWViewValue        float64 `json:"Ab13An13Value"`
	B512DHWViewUnit         string  `json:"Ab13An13Unit"`
	B512DHWViewUnitConflict bool    `json:"Ab13An13UnitConflict"`
	B512DHWViewNamedDomain  bool    `json:"Ab13An13NamedDomain"`
}

type B512DHWViewRequest struct {
	B512DHWViewTaskActive          *bool                 `json:"Ab13An13TaskActive"`
	B512DHWViewHeatDemand          *bool                 `json:"Ab13An13HeatDemand"`
	B512DHWViewLoadingActive       *bool                 `json:"Ab13An13LoadingActive"`
	B512DHWViewLoadingPumpFault    *bool                 `json:"Ab13An13LoadingPumpFault"`
	B512DHWViewSensorConditionBits [4]*bool              `json:"Ab13An13SensorConditionBits"`
	B512DHWViewModulation          *B512DHWViewScalar    `json:"Ab13An13Modulation"`
	B512DHWViewFlow                *B512DHWViewScalar    `json:"Ab13An13Flow"`
	B512DHWViewTemperatures        [3]*B512DHWViewScalar `json:"Ab13An13Temperatures"`
	B512DHWViewInletConditionBits  [2]*bool              `json:"Ab13An13InletConditionBits"`
}

type B512DHWViewPositive struct {
	B512DHWViewLoadingDisabled       *bool              `json:"Ab13An13LoadingDisabled"`
	B512DHWViewInternalTestActive    *bool              `json:"Ab13An13InternalTestActive"`
	B512DHWViewReady                 *bool              `json:"Ab13An13Ready"`
	B512DHWViewFrostProtectionActive *bool              `json:"Ab13An13FrostProtectionActive"`
	B512DHWViewLegionellaDemand      *bool              `json:"Ab13An13LegionellaDemand"`
	B512DHWViewReset                 *bool              `json:"Ab13An13Reset"`
	B512DHWViewSetpoint              *B512DHWViewScalar `json:"Ab13An13Setpoint"`
	B512DHWViewReturnTemperature     *B512DHWViewScalar `json:"Ab13An13ReturnTemperature"`
	B512DHWViewAdjustment            *B512DHWViewScalar `json:"Ab13An13Adjustment"`
}

type B512DHWViewResult struct {
	B512DHWViewRaw           []byte               `json:"Ab13An13Raw"`
	B512DHWViewUnmodelled    map[int]byte         `json:"Ab13An13Unmodelled"`
	B512DHWViewAdmitted      bool                 `json:"Ab13An13Admitted"`
	B512DHWViewComplete      bool                 `json:"Ab13An13Complete"`
	B512DHWViewContractGrade string               `json:"Ab13An13ContractGrade"`
	B512DHWViewFieldGrade    string               `json:"Ab13An13FieldGrade"`
	B512DHWViewQualification string               `json:"Ab13An13Qualification"`
	B512DHWViewPresence      string               `json:"Ab13An13Presence"`
	B512DHWViewFreshness     string               `json:"Ab13An13Freshness"`
	B512DHWViewExecution     string               `json:"Ab13An13Execution"`
	B512DHWViewDefaults      map[string]float64   `json:"Ab13An13Defaults"`
	B512DHWViewRequest       *B512DHWViewRequest  `json:"Ab13An13Request"`
	B512DHWViewPositive      *B512DHWViewPositive `json:"Ab13An13Positive"`
}

func b512DHWViewBit(b512DHWViewRaw byte, b512DHWViewIndex uint) *bool {
	b512DHWViewValue := b512DHWViewRaw&(1<<b512DHWViewIndex) != 0
	return &b512DHWViewValue
}

func b512DHWViewScalar(b512DHWViewRaw byte, b512DHWViewDivisor float64, b512DHWViewUnit string) *B512DHWViewScalar {
	return &B512DHWViewScalar{B512DHWViewRaw: b512DHWViewRaw, B512DHWViewValue: float64(b512DHWViewRaw) / b512DHWViewDivisor, B512DHWViewUnit: b512DHWViewUnit}
}

func B512DHWViewInterpret(b512DHWViewFamily, b512DHWViewVariant, b512DHWViewDirection string, b512DHWViewImage []byte) B512DHWViewResult {
	b512DHWViewResult := B512DHWViewResult{
		B512DHWViewRaw:           append([]byte(nil), b512DHWViewImage...),
		B512DHWViewUnmodelled:    make(map[int]byte),
		B512DHWViewContractGrade: "Proven", B512DHWViewFieldGrade: "Hypothesis",
		B512DHWViewQualification: "Unknown", B512DHWViewPresence: "Unknown",
		B512DHWViewFreshness: "Unknown", B512DHWViewExecution: "Unknown",
	}
	for b512DHWViewIndex, b512DHWViewRaw := range b512DHWViewImage {
		b512DHWViewResult.B512DHWViewUnmodelled[b512DHWViewIndex] = b512DHWViewRaw
	}
	if b512DHWViewFamily != B512DHWViewFamily || b512DHWViewVariant != B512DHWViewVariant ||
		(b512DHWViewDirection != "request" && b512DHWViewDirection != "positive") ||
		len(b512DHWViewImage) < 3 || b512DHWViewImage[0] != 0xb5 || b512DHWViewImage[1] != 0x12 || b512DHWViewImage[2] != 0x06 {
		return b512DHWViewResult
	}
	b512DHWViewResult.B512DHWViewAdmitted = true
	for b512DHWViewIndex := 0; b512DHWViewIndex < 3; b512DHWViewIndex++ {
		delete(b512DHWViewResult.B512DHWViewUnmodelled, b512DHWViewIndex)
	}
	if b512DHWViewDirection == "request" {
		b512DHWViewRequest := &B512DHWViewRequest{}
		b512DHWViewResult.B512DHWViewRequest = b512DHWViewRequest
		b512DHWViewResult.B512DHWViewComplete = len(b512DHWViewImage) >= 10
		b512DHWViewResult.B512DHWViewDefaults = map[string]float64{"Ab13An13Flags": 0, "Ab13An13Scalars": 0}
		if len(b512DHWViewImage) > 3 {
			b512DHWViewRequest.B512DHWViewTaskActive = b512DHWViewBit(b512DHWViewImage[3], 0)
			b512DHWViewRequest.B512DHWViewHeatDemand = b512DHWViewBit(b512DHWViewImage[3], 1)
			b512DHWViewRequest.B512DHWViewLoadingActive = b512DHWViewBit(b512DHWViewImage[3], 2)
			b512DHWViewRequest.B512DHWViewLoadingPumpFault = b512DHWViewBit(b512DHWViewImage[3], 3)
			for b512DHWViewIndex := range b512DHWViewRequest.B512DHWViewSensorConditionBits {
				b512DHWViewRequest.B512DHWViewSensorConditionBits[b512DHWViewIndex] = b512DHWViewBit(b512DHWViewImage[3], uint(b512DHWViewIndex+4))
			}
		}
		if len(b512DHWViewImage) > 4 {
			b512DHWViewRequest.B512DHWViewModulation = b512DHWViewScalar(b512DHWViewImage[4], 2, "formal Celsius; modulation interpretation")
			b512DHWViewRequest.B512DHWViewModulation.B512DHWViewUnitConflict = true
		}
		if len(b512DHWViewImage) > 5 {
			b512DHWViewRequest.B512DHWViewFlow = b512DHWViewScalar(b512DHWViewImage[5], 10, "litres/minute")
		}
		for b512DHWViewIndex := range b512DHWViewRequest.B512DHWViewTemperatures {
			if len(b512DHWViewImage) > 6+b512DHWViewIndex {
				b512DHWViewRequest.B512DHWViewTemperatures[b512DHWViewIndex] = b512DHWViewScalar(b512DHWViewImage[6+b512DHWViewIndex], 2, "Celsius")
			}
		}
		if len(b512DHWViewImage) > 9 {
			for b512DHWViewIndex := range b512DHWViewRequest.B512DHWViewInletConditionBits {
				b512DHWViewRequest.B512DHWViewInletConditionBits[b512DHWViewIndex] = b512DHWViewBit(b512DHWViewImage[9], uint(b512DHWViewIndex))
			}
		}
		for b512DHWViewIndex := 3; b512DHWViewIndex < len(b512DHWViewImage) && b512DHWViewIndex < 10; b512DHWViewIndex++ {
			delete(b512DHWViewResult.B512DHWViewUnmodelled, b512DHWViewIndex)
		}
		if len(b512DHWViewImage) > 9 && b512DHWViewImage[9]&0xfc != 0 {
			b512DHWViewResult.B512DHWViewUnmodelled[9] = b512DHWViewImage[9] & 0xfc
		}
	} else {
		b512DHWViewPositive := &B512DHWViewPositive{}
		b512DHWViewResult.B512DHWViewPositive = b512DHWViewPositive
		b512DHWViewResult.B512DHWViewComplete = len(b512DHWViewImage) >= 7
		b512DHWViewResult.B512DHWViewDefaults = map[string]float64{"Ab13An13Flags": 0, "Ab13An13Scalars": 0, "Ab13An13Adjustment": -1}
		if len(b512DHWViewImage) > 3 {
			b512DHWViewPositive.B512DHWViewLoadingDisabled = b512DHWViewBit(b512DHWViewImage[3], 0)
			b512DHWViewPositive.B512DHWViewInternalTestActive = b512DHWViewBit(b512DHWViewImage[3], 1)
			b512DHWViewPositive.B512DHWViewReady = b512DHWViewBit(b512DHWViewImage[3], 2)
			b512DHWViewPositive.B512DHWViewFrostProtectionActive = b512DHWViewBit(b512DHWViewImage[3], 3)
			b512DHWViewPositive.B512DHWViewLegionellaDemand = b512DHWViewBit(b512DHWViewImage[3], 4)
			b512DHWViewPositive.B512DHWViewReset = b512DHWViewBit(b512DHWViewImage[3], 5)
		}
		if len(b512DHWViewImage) > 4 {
			b512DHWViewPositive.B512DHWViewSetpoint = b512DHWViewScalar(b512DHWViewImage[4], 2, "Celsius")
		}
		if len(b512DHWViewImage) > 5 {
			b512DHWViewPositive.B512DHWViewReturnTemperature = b512DHWViewScalar(b512DHWViewImage[5], 2, "Celsius")
		}
		if len(b512DHWViewImage) > 6 {
			b512DHWViewValue := int8(b512DHWViewImage[6])
			b512DHWViewPositive.B512DHWViewAdjustment = &B512DHWViewScalar{B512DHWViewRaw: b512DHWViewImage[6], B512DHWViewValue: float64(b512DHWViewValue), B512DHWViewNamedDomain: b512DHWViewValue >= -1 && b512DHWViewValue <= 2}
		}
		for b512DHWViewIndex := 3; b512DHWViewIndex < len(b512DHWViewImage) && b512DHWViewIndex < 7; b512DHWViewIndex++ {
			delete(b512DHWViewResult.B512DHWViewUnmodelled, b512DHWViewIndex)
		}
		if len(b512DHWViewImage) > 3 && b512DHWViewImage[3]&0xc0 != 0 {
			b512DHWViewResult.B512DHWViewUnmodelled[3] = b512DHWViewImage[3] & 0xc0
		}
	}
	return b512DHWViewResult
}
