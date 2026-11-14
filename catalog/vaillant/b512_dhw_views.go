package vaillant

const (
	B512DHWViewFamily  = "vaillant-dhw-family"
	B512DHWViewVariant = "b512-op06-candidate"
)

type B512DHWViewScalar struct {
	B512DHWViewRaw          byte    `json:"raw"`
	B512DHWViewValue        float64 `json:"value"`
	B512DHWViewUnit         string  `json:"unit"`
	B512DHWViewUnitConflict bool    `json:"unit_conflict"`
	B512DHWViewNamedDomain  bool    `json:"named_domain"`
}

type B512DHWViewRequest struct {
	B512DHWViewTaskActive          *bool                 `json:"task_active"`
	B512DHWViewHeatDemand          *bool                 `json:"heat_demand"`
	B512DHWViewLoadingActive       *bool                 `json:"loading_active"`
	B512DHWViewLoadingPumpFault    *bool                 `json:"loading_pump_fault"`
	B512DHWViewSensorConditionBits [4]*bool              `json:"sensor_condition_bits"`
	B512DHWViewModulation          *B512DHWViewScalar    `json:"modulation"`
	B512DHWViewFlow                *B512DHWViewScalar    `json:"flow"`
	B512DHWViewTemperatures        [3]*B512DHWViewScalar `json:"temperatures"`
	B512DHWViewInletConditionBits  [2]*bool              `json:"inlet_condition_bits"`
}

type B512DHWViewPositive struct {
	B512DHWViewLoadingDisabled       *bool              `json:"loading_disabled"`
	B512DHWViewInternalTestActive    *bool              `json:"internal_test_active"`
	B512DHWViewReady                 *bool              `json:"ready"`
	B512DHWViewFrostProtectionActive *bool              `json:"frost_protection_active"`
	B512DHWViewLegionellaDemand      *bool              `json:"legionella_demand"`
	B512DHWViewReset                 *bool              `json:"reset"`
	B512DHWViewSetpoint              *B512DHWViewScalar `json:"setpoint"`
	B512DHWViewReturnTemperature     *B512DHWViewScalar `json:"return_temperature"`
	B512DHWViewAdjustment            *B512DHWViewScalar `json:"adjustment"`
}

type B512DHWViewResult struct {
	B512DHWViewRaw           []byte               `json:"raw"`
	B512DHWViewUnmodelled    map[int]byte         `json:"unmodelled"`
	B512DHWViewAdmitted      bool                 `json:"admitted"`
	B512DHWViewComplete      bool                 `json:"complete"`
	B512DHWViewContractGrade string               `json:"contract_grade"`
	B512DHWViewFieldGrade    string               `json:"field_grade"`
	B512DHWViewQualification string               `json:"qualification"`
	B512DHWViewPresence      string               `json:"presence"`
	B512DHWViewFreshness     string               `json:"freshness"`
	B512DHWViewExecution     string               `json:"execution"`
	B512DHWViewDefaults      map[string]float64   `json:"defaults"`
	B512DHWViewRequest       *B512DHWViewRequest  `json:"request"`
	B512DHWViewPositive      *B512DHWViewPositive `json:"positive"`
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
		b512DHWViewResult.B512DHWViewDefaults = map[string]float64{"flags": 0, "scalars": 0}
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
		b512DHWViewResult.B512DHWViewDefaults = map[string]float64{"flags": 0, "scalars": 0, "adjustment": -1}
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
