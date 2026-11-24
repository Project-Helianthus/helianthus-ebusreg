package wolf

import (
	"bytes"
	"errors"
	"math"
	"testing"
)

func fieldValue(t *testing.T, m Message, name string) Field {
	t.Helper()
	for _, f := range m.Fields {
		if f.Name == name {
			return f
		}
	}
	t.Fatalf("message %s has no field %s", m.Name, name)
	return Field{}
}

func TestWolfBurnerStatusAndModulationDecode(t *testing.T) {
	m, err := DecodeFrame(PBBurner, SBBurnerStatus, []byte{0x01, 0x05, 0xFF, 0x3C})
	if err != nil || !m.Known || m.Name != "burner_status" {
		t.Fatalf("burner status = %#v, err %v", m, err)
	}
	if fieldValue(t, m, "modulation_raw").Value != byte(0xFF) {
		t.Fatalf("dummy modulation byte not kept raw")
	}

	m, err = DecodeFrame(PBBurner, SBModulationActual, []byte{40, 100})
	if err != nil {
		t.Fatal(err)
	}
	if f := fieldValue(t, m, "modulation_actual"); f.Value != byte(40) || f.Unit != "%" {
		t.Fatalf("modulation actual = %#v", f)
	}
	if _, err := DecodeFrame(PBBurner, SBModulationActual, []byte{101, 100}); !errors.Is(err, ErrRange) {
		t.Fatalf("modulation above 100 err = %v; want ErrRange", err)
	}
}

func TestWolfHeatRequestNamesTheTenStates(t *testing.T) {
	cases := map[byte]string{0x00: "aus", 0x01: "keine", 0x44: "Reglerstop_stufig", 0x55: "Brauchwasser", 0x66: "Brauchwasser_Reglerstop", 0xAA: "Heizen", 0xBB: "Brauchwasser_Heizen", 0xCC: "Emissionskontrolle", 0xDD: "TUeV", 0xEE: "Reglerstop"}
	for code, name := range cases {
		m, err := DecodeFrame(PBBurner, SBHeatRequestStatus, []byte{code, 60, 0x10})
		if err != nil {
			t.Fatal(err)
		}
		if got := fieldValue(t, m, "heat_request_state").Value; got != name {
			t.Fatalf("state 0x%02X = %v; want %s", code, got, name)
		}
	}
	m, _ := DecodeFrame(PBBurner, SBHeatRequestStatus, []byte{0x11, 60})
	if got := fieldValue(t, m, "heat_request_state").Value; got != byte(0x11) {
		t.Fatalf("undocumented state = %v; want raw 0x11", got)
	}
}

func TestWolfBoilerSetpointIsTemp2(t *testing.T) {
	// 0x0280 is 640 sixteenths: 40.0 degrees C, the page's example.
	m, err := DecodeFrame(PBBoiler, SBBoilerSetpoint, []byte{0x80, 0x02, 0x00, 0x01})
	if err != nil {
		t.Fatal(err)
	}
	if f := fieldValue(t, m, "boiler_setpoint"); math.Abs(f.Value.(float64)-40.0) > 1e-9 || f.Unit != "degC" {
		t.Fatalf("boiler setpoint = %#v", f)
	}
}

func TestWolfSolarFramesDecode(t *testing.T) {
	m, err := DecodeFrame(PBWolf, SBSolarPump, []byte{solarPumpOn, 0x00, 0x20, 0x30})
	if err != nil || fieldValue(t, m, "pump_state").Value != "on" {
		t.Fatalf("pump on = %#v, err %v", m, err)
	}
	m, _ = DecodeFrame(PBWolf, SBSolarPump, []byte{0x01, 0x00})
	if fieldValue(t, m, "pump_state").Value != byte(0x01) {
		t.Fatalf("undocumented pump byte not kept raw")
	}
	// 0x0180 is 384/256 = 1.5 W.
	m, err = DecodeFrame(PBWolf, SBSolarPower, []byte{0x80, 0x01, 0x07, 0x09})
	if err != nil {
		t.Fatal(err)
	}
	if f := fieldValue(t, m, "solar_power"); math.Abs(f.Value.(float64)-1.5) > 1e-9 || f.Unit != "W" {
		t.Fatalf("solar power = %#v", f)
	}
	if _, err := DecodeFrame(PBWolf, SBSolarPower, []byte{0x80}); !errors.Is(err, ErrWidth) {
		t.Fatalf("short solar power err = %v; want ErrWidth", err)
	}
	m, _ = DecodeFrame(PBWolf, SBSolarStatus, []byte{0x05, 0x00, 0x01, 0x02})
	if !m.Known || fieldValue(t, m, "solar_status_raw").Value != byte(0x05) {
		t.Fatalf("solar status = %#v", m)
	}
}

// An undocumented PB/SB pair is kept raw and not decoded.
func TestWolfUnknownFrameStaysRaw(t *testing.T) {
	m, err := DecodeFrame(0x09, 0x01, []byte{0xAA, 0xBB})
	if err != nil || m.Known || !bytes.Equal(m.Raw, []byte{0xAA, 0xBB}) {
		t.Fatalf("unknown frame = %#v, err %v", m, err)
	}
}

// Service-register reads: the table's ID, scaled UIN values.
func TestWolfServiceRegisterReadDecodes(t *testing.T) {
	reg, ok := ServiceRegisterByID([3]byte{0x84, 0x22, 0x00})
	if !ok || reg.Name != "hg01 hysterese_vorlauf" || reg.Scale != 10 {
		t.Fatalf("hysteresis register = %#v", reg)
	}
	if got := EncodeServiceRead(reg); !bytes.Equal(got, []byte{0x84, 0x22, 0x00}) {
		t.Fatalf("read request = % X", got)
	}
	// 0x0064 = 100 with scale 10: 10.0 K.
	f, err := DecodeServiceRead(reg, []byte{0x64, 0x00})
	if err != nil || f.Value.(float64) != 10.0 {
		t.Fatalf("read value = %#v, err %v", f, err)
	}
	if _, err := DecodeServiceRead(reg, []byte{0x64}); !errors.Is(err, ErrWidth) {
		t.Fatalf("short read err = %v; want ErrWidth", err)
	}
}

// Service-register write: the documented layout is id_lo, id_hi, the UIN value
// little-endian, then the four HEX bytes the caller supplies.
func TestWolfServiceRegisterWriteLayout(t *testing.T) {
	reg, _ := ServiceRegisterByID([3]byte{0x84, 0x22, 0x00})
	got, err := EncodeServiceWrite(reg, 0x0064, [4]byte{0x01, 0x02, 0x03, 0x04})
	if err != nil || !bytes.Equal(got, []byte{0x84, 0x22, 0x64, 0x00, 0x01, 0x02, 0x03, 0x04}) {
		t.Fatalf("write = % X, err %v", got, err)
	}
	ro, ok := ServiceRegisterByID([3]byte{0xCD, 0x59, 0x01})
	if !ok || ro.Writable {
		t.Fatalf("read-only register = %#v", ro)
	}
	if _, err := EncodeServiceWrite(ro, 1, [4]byte{}); !errors.Is(err, ErrReadOnly) {
		t.Fatalf("read-only write err = %v; want ErrReadOnly", err)
	}
}

func TestWolfServiceTableCounts(t *testing.T) {
	writable := 0
	for _, reg := range serviceRegisterTable {
		if reg.Writable {
			writable++
		}
	}
	if len(serviceRegisterTable) != 26 || writable != 11 {
		t.Fatalf("table rows = %d, writable = %d; the page's table has 26 rows and 11 R+W", len(serviceRegisterTable), writable)
	}
}

// Each documented write builds the documented request layout.
func TestWolfWriteEncodersFollowDocumentedLayouts(t *testing.T) {
	if got := EncodeBurnerStatus(0x05, 0xFF, []byte{0x3C}); !bytes.Equal(got, []byte{0x01, 0x05, 0xFF, 0x3C}) {
		t.Fatalf("burner status = % X", got)
	}
	if got, err := EncodeModulation(40, 100); err != nil || !bytes.Equal(got, []byte{40, 100}) {
		t.Fatalf("modulation = % X, err %v", got, err)
	}
	if _, err := EncodeModulation(101, 0); !errors.Is(err, ErrRange) {
		t.Fatalf("modulation range err = %v", err)
	}
	if got := EncodeHeatRequest(0xAA, 60, []byte{0x10}); !bytes.Equal(got, []byte{0xAA, 60, 0x10}) {
		t.Fatalf("heat request = % X", got)
	}
	if got, err := EncodeBoilerSetpoint(40.0, 0x02, 0x03); err != nil || !bytes.Equal(got, []byte{0x80, 0x02, 0x02, 0x03}) {
		t.Fatalf("boiler setpoint = % X, err %v", got, err)
	}
	if got := EncodeSolarStatus(0x05, 0x00, []byte{0x01}); !bytes.Equal(got, []byte{0x05, 0x00, 0x01}) {
		t.Fatalf("solar status = % X", got)
	}
	if got := EncodeSolarPump(true, 0x00, nil); !bytes.Equal(got, []byte{0xBD, 0x00}) {
		t.Fatalf("solar pump on = % X", got)
	}
	if got := EncodeSolarPump(false, 0x00, nil); !bytes.Equal(got, []byte{0xBC, 0x00}) {
		t.Fatalf("solar pump off = % X", got)
	}
	if got, err := EncodeSolarPower(1.5, 0x07, 0x09); err != nil || !bytes.Equal(got, []byte{0x80, 0x01, 0x07, 0x09}) {
		t.Fatalf("solar power = % X, err %v", got, err)
	}
}

// The encoders read back through the matching decoders.
func TestWolfWriteEncodersRoundTripThroughDecoders(t *testing.T) {
	setpoint, err := EncodeBoilerSetpoint(-5.5, 0x01, 0x02)
	if err != nil {
		t.Fatal(err)
	}
	m, err := DecodeFrame(PBBoiler, SBBoilerSetpoint, setpoint)
	if err != nil {
		t.Fatal(err)
	}
	if got := fieldValue(t, m, "boiler_setpoint").Value.(float64); math.Abs(got-(-5.5)) > 1e-9 {
		t.Fatalf("round trip = %v", got)
	}
	power, err := EncodeSolarPower(-12.25, 0x00, 0x00)
	if err != nil {
		t.Fatal(err)
	}
	m, err = DecodeFrame(PBWolf, SBSolarPower, power)
	if err != nil {
		t.Fatal(err)
	}
	if got := fieldValue(t, m, "solar_power").Value.(float64); math.Abs(got-(-12.25)) > 1e-9 {
		t.Fatalf("solar round trip = %v", got)
	}
}
