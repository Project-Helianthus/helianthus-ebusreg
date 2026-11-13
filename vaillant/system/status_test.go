package system

import (
	"bytes"
	"testing"
)

func TestAb22EightRawFields(statusT *testing.T) {
	statusPayload := []byte{0x01, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67, 0x78}
	statusResult := DecodeStatus(StatusRequestContext{StatusFamily: "ab22-synthetic-family", StatusOpcode: 0x08}, statusPayload)
	statusWant := StatusStatusFields{
		StatusBindingZone: 0x01, SpecialFunctionStatus: 0x12,
		StatusHeatingOperatingMode: 0x23, StatusCoolingOperatingMode: 0x34,
		StatusInfo: 0x45, StatusFrostProtection: 0x56,
		StatusHeatingTemperatureRaw: 0x67, StatusCoolingTemperatureRaw: 0x78,
	}
	if statusResult.StatusFields == nil || *statusResult.StatusFields != statusWant {
		statusT.Fatalf("ab22 fields = %+v, want %+v", statusResult.StatusFields, statusWant)
	}
	if !bytes.Equal(statusResult.StatusPayload, statusPayload) {
		statusT.Fatalf("ab22 raw bytes = %x, want %x", statusResult.StatusPayload, statusPayload)
	}
}

func TestAb22ContextDimensions(statusT *testing.T) {
	// Independently vary every context field. Family labels are synthetic and
	// unqualified; retaining them must never introduce family-specific semantics.
	for _, statusFamily := range []string{"ab22-synthetic-family", "ab22-other-family", "unknown", ""} {
		for _, statusOpcode := range []byte{0x08, 0x03, 0x00, 0xff} {
			statusContext := StatusRequestContext{StatusFamily: statusFamily, StatusOpcode: statusOpcode}
			statusPayload := []byte{1, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67, 0x78}
			statusResult := DecodeStatus(statusContext, statusPayload)
			if statusResult.StatusContext != statusContext {
				statusT.Fatalf("ab22 context = %+v, want %+v", statusResult.StatusContext, statusContext)
			}
			if (statusResult.StatusFields != nil) != (statusOpcode == 0x08) {
				statusT.Fatalf("ab22 family %q opcode %02x fields = %+v", statusFamily, statusOpcode, statusResult.StatusFields)
			}
			if !bytes.Equal(statusResult.StatusPayload, statusPayload) {
				statusT.Fatalf("ab22 context altered raw payload: %x", statusResult.StatusPayload)
			}
		}
	}
}

func TestAb22IncompleteWidths(statusT *testing.T) {
	for _, statusWidth := range []int{0, 1, 7, 8, 9} {
		statusPayload := bytes.Repeat([]byte{0x78}, statusWidth)
		statusResult := DecodeStatus(StatusRequestContext{StatusOpcode: 0x08}, statusPayload)
		if (statusResult.StatusFields != nil) != (statusWidth == 8) {
			statusT.Fatalf("ab22 width %d fields = %+v", statusWidth, statusResult.StatusFields)
		}
		if !bytes.Equal(statusResult.StatusPayload, statusPayload) {
			statusT.Fatalf("ab22 width %d raw bytes = %x", statusWidth, statusResult.StatusPayload)
		}
	}
}

func TestAb22OwnedPayload(statusT *testing.T) {
	for _, statusWidth := range []int{1, 7, 8, 9} {
		statusInput := bytes.Repeat([]byte{0x67}, statusWidth)
		statusContext := StatusRequestContext{StatusFamily: "unknown", StatusOpcode: 0x08}
		statusFirst := DecodeStatus(statusContext, statusInput)
		statusSecond := DecodeStatus(statusContext, statusInput)
		statusInput[0] = 0x12
		if statusFirst.StatusPayload[0] != 0x67 || statusSecond.StatusPayload[0] != 0x67 {
			statusT.Fatal("ab22 input aliases a result")
		}
		statusFirst.StatusPayload[0] = 0x23
		if statusInput[0] != 0x12 || statusSecond.StatusPayload[0] != 0x67 {
			statusT.Fatal("ab22 results alias input or each other")
		}
		if statusWidth == 8 {
			if statusFirst.StatusFields.StatusBindingZone != 0x67 {
				statusT.Fatal("ab22 fields changed with payload")
			}
			statusFirst.StatusFields.StatusBindingZone = 0x34
			if statusSecond.StatusFields.StatusBindingZone != 0x67 {
				statusT.Fatal("ab22 field storage aliases another result")
			}
		}
	}
}
