package system

import (
	"bytes"
	"testing"
)

func TestAb22EightRawFields(ab22T *testing.T) {
	ab22Payload := []byte{0x01, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67, 0x78}
	ab22Result := Ab22DecodeStatus(Ab22RequestContext{Ab22Family: "ab22-synthetic-family", Ab22Opcode: 0x08}, ab22Payload)
	ab22Want := Ab22StatusFields{
		Ab22BindingZone: 0x01, Ab22SpecialFunctionStatus: 0x12,
		Ab22HeatingOperatingMode: 0x23, Ab22CoolingOperatingMode: 0x34,
		Ab22StatusInfo: 0x45, Ab22FrostProtection: 0x56,
		Ab22HeatingTemperatureRaw: 0x67, Ab22CoolingTemperatureRaw: 0x78,
	}
	if ab22Result.Ab22Fields == nil || *ab22Result.Ab22Fields != ab22Want {
		ab22T.Fatalf("ab22 fields = %+v, want %+v", ab22Result.Ab22Fields, ab22Want)
	}
	if !bytes.Equal(ab22Result.Ab22Payload, ab22Payload) {
		ab22T.Fatalf("ab22 raw bytes = %x, want %x", ab22Result.Ab22Payload, ab22Payload)
	}
}

func TestAb22ContextDimensions(ab22T *testing.T) {
	// Independently vary every context field. Family labels are synthetic and
	// unqualified; retaining them must never introduce family-specific semantics.
	for _, ab22Family := range []string{"ab22-synthetic-family", "ab22-other-family", "unknown", ""} {
		for _, ab22Opcode := range []byte{0x08, 0x03, 0x00, 0xff} {
			ab22Context := Ab22RequestContext{Ab22Family: ab22Family, Ab22Opcode: ab22Opcode}
			ab22Payload := []byte{1, 0x12, 0x23, 0x34, 0x45, 0x56, 0x67, 0x78}
			ab22Result := Ab22DecodeStatus(ab22Context, ab22Payload)
			if ab22Result.Ab22Context != ab22Context {
				ab22T.Fatalf("ab22 context = %+v, want %+v", ab22Result.Ab22Context, ab22Context)
			}
			if (ab22Result.Ab22Fields != nil) != (ab22Opcode == 0x08) {
				ab22T.Fatalf("ab22 family %q opcode %02x fields = %+v", ab22Family, ab22Opcode, ab22Result.Ab22Fields)
			}
			if !bytes.Equal(ab22Result.Ab22Payload, ab22Payload) {
				ab22T.Fatalf("ab22 context altered raw payload: %x", ab22Result.Ab22Payload)
			}
		}
	}
}

func TestAb22IncompleteWidths(ab22T *testing.T) {
	for _, ab22Width := range []int{0, 1, 7, 8, 9} {
		ab22Payload := bytes.Repeat([]byte{0x78}, ab22Width)
		ab22Result := Ab22DecodeStatus(Ab22RequestContext{Ab22Opcode: 0x08}, ab22Payload)
		if (ab22Result.Ab22Fields != nil) != (ab22Width == 8) {
			ab22T.Fatalf("ab22 width %d fields = %+v", ab22Width, ab22Result.Ab22Fields)
		}
		if !bytes.Equal(ab22Result.Ab22Payload, ab22Payload) {
			ab22T.Fatalf("ab22 width %d raw bytes = %x", ab22Width, ab22Result.Ab22Payload)
		}
	}
}

func TestAb22OwnedPayload(ab22T *testing.T) {
	for _, ab22Width := range []int{1, 7, 8, 9} {
		ab22Input := bytes.Repeat([]byte{0x67}, ab22Width)
		ab22Context := Ab22RequestContext{Ab22Family: "unknown", Ab22Opcode: 0x08}
		ab22First := Ab22DecodeStatus(ab22Context, ab22Input)
		ab22Second := Ab22DecodeStatus(ab22Context, ab22Input)
		ab22Input[0] = 0x12
		if ab22First.Ab22Payload[0] != 0x67 || ab22Second.Ab22Payload[0] != 0x67 {
			ab22T.Fatal("ab22 input aliases a result")
		}
		ab22First.Ab22Payload[0] = 0x23
		if ab22Input[0] != 0x12 || ab22Second.Ab22Payload[0] != 0x67 {
			ab22T.Fatal("ab22 results alias input or each other")
		}
		if ab22Width == 8 {
			if ab22First.Ab22Fields.Ab22BindingZone != 0x67 {
				ab22T.Fatal("ab22 fields changed with payload")
			}
			ab22First.Ab22Fields.Ab22BindingZone = 0x34
			if ab22Second.Ab22Fields.Ab22BindingZone != 0x67 {
				ab22T.Fatal("ab22 field storage aliases another result")
			}
		}
	}
}
