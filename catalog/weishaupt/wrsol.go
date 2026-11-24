package weishaupt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
)

// Primary and secondary bytes of the documented WRSol messages.
const (
	PBSystem         byte = 0x07
	SBDateTime       byte = 0x00
	SBIdentification byte = 0x04
	PBMemoryServer   byte = 0x09
	SBReadRAM        byte = 0x00

	// ManufacturerWRSol is the manufacturer byte the WRSol identification
	// response states (S3 = 10h, Weishaupt).
	ManufacturerWRSol byte = 0x10
)

// Errors of the Weishaupt catalogue.
var (
	ErrWidth     = errors.New("weishaupt: payload width does not match the documented layout")
	ErrRange     = errors.New("weishaupt: value outside the documented range")
	ErrManuf     = errors.New("weishaupt: manufacturer byte does not match the page")
	ErrValueSize = errors.New("weishaupt: value does not fit the documented type")
	ErrBCD       = errors.New("weishaupt: invalid BCD digit")
)

// Field is one named decoded value. Unit is empty when the page states none.
type Field struct {
	Name  string
	Value any
	Unit  string
}

// Identification is the decoded 07 04 reply of a WRSol. Versions are the BCD
// values the page states; DeviceID is the ASCII device ID the page lists.
type Identification struct {
	Manufacturer    byte
	DeviceID        string
	SoftwareVersion int
	SoftwareRev     int
	HardwareVersion int
	HardwareRev     int
}

// DecodeIdentification decodes the ten-byte 07 04 reply: the manufacturer byte
// (0x10), the five-byte ASCII device ID and four BCD version fields.
func DecodeIdentification(payload []byte) (Identification, error) {
	if len(payload) != 10 {
		return Identification{}, fmt.Errorf("%w: identification wants 10 bytes, got %d", ErrWidth, len(payload))
	}
	if payload[0] != ManufacturerWRSol {
		return Identification{}, fmt.Errorf("%w: 0x%02X", ErrManuf, payload[0])
	}
	var out Identification
	out.Manufacturer = payload[0]
	out.DeviceID = strings.Trim(string(payload[1:6]), " \x00")
	var err error
	if out.SoftwareVersion, err = decodeBCD(payload[6]); err != nil {
		return Identification{}, err
	}
	if out.SoftwareRev, err = decodeBCD(payload[7]); err != nil {
		return Identification{}, err
	}
	if out.HardwareVersion, err = decodeBCD(payload[8]); err != nil {
		return Identification{}, err
	}
	if out.HardwareRev, err = decodeBCD(payload[9]); err != nil {
		return Identification{}, err
	}
	return out, nil
}

// DateTime is the calendar value of the 07 00 message. Weekday is 1..7 and
// Year is 0..99, as the page states.
type DateTime struct {
	Seconds, Minutes, Hours, Day, Month, Weekday, Year int
}

// EncodeDateTime builds the nine-byte 07 00 payload: the outdoor temperature as
// DATA2b (signed 16-bit little-endian, 1/256 degrees C) followed by the BCD
// time fields. The page's range is -50.0 to 50.0 degrees C; each BCD field is
// checked against its documented range.
func EncodeDateTime(outdoorC float64, t DateTime) ([]byte, error) {
	if outdoorC < -50.0 || outdoorC > 50.0 {
		return nil, fmt.Errorf("%w: outdoor temperature %.2f", ErrRange, outdoorC)
	}
	fields := []struct {
		name     string
		value    int
		min, max int
	}{
		{"seconds", t.Seconds, 0, 59},
		{"minutes", t.Minutes, 0, 59},
		{"hours", t.Hours, 0, 23},
		{"day", t.Day, 1, 31},
		{"month", t.Month, 1, 12},
		{"weekday", t.Weekday, 1, 7},
		{"year", t.Year, 0, 99},
	}
	bcd := make([]byte, len(fields))
	for i, f := range fields {
		if f.value < f.min || f.value > f.max {
			return nil, fmt.Errorf("%w: %s %d", ErrRange, f.name, f.value)
		}
		bcd[i] = encodeBCD(f.value)
	}
	raw := math.Round(outdoorC * 256)
	out := make([]byte, 9)
	binary.LittleEndian.PutUint16(out[0:2], uint16(int16(raw)))
	copy(out[2:], bcd)
	return out, nil
}

// EncodeRAMRead builds the 09 00 request: the start address little-endian and
// the byte count 02 that the page states.
func EncodeRAMRead(address uint16) []byte {
	return []byte{byte(address), byte(address >> 8), 0x02}
}

// RAMRegister is one row of the WRSol RAM address map. Format is the wire type
// the page states ("Integer" signed, "Word" unsigned, or a non-numeric layout
// kept raw); Scale is "X1" or "X10".
type RAMRegister struct {
	Address     uint16
	Name        string
	Format      string
	Scale       string
	Unit        string
	Description string
}

// RAMByAddress returns the RAM map row for an address.
func RAMByAddress(address uint16) (RAMRegister, bool) {
	for _, reg := range ramTable {
		if reg.Address == address {
			return reg, true
		}
	}
	return RAMRegister{}, false
}

// DecodeRAMRead decodes the two-byte 09 00 reply for a RAM row. Integer and Word
// rows are scaled by X10 (divide by 10) or X1; other layouts stay raw.
func DecodeRAMRead(reg RAMRegister, payload []byte) (Field, error) {
	if len(payload) != 2 {
		return Field{}, fmt.Errorf("%w: RAM read wants 2 bytes, got %d", ErrWidth, len(payload))
	}
	field := Field{Name: reg.Name, Unit: reg.Unit}
	raw := binary.LittleEndian.Uint16(payload)
	divisor := 1.0
	if reg.Scale == "X10" {
		divisor = 10
	}
	switch reg.Format {
	case "Integer":
		field.Value = float64(int16(raw)) / divisor
	case "Word":
		field.Value = float64(raw) / divisor
	default:
		field.Value = append([]byte(nil), payload...)
		field.Unit = ""
	}
	return field, nil
}

func decodeBCD(b byte) (int, error) {
	if b>>4 > 9 || b&0x0F > 9 {
		return 0, fmt.Errorf("%w: 0x%02X", ErrBCD, b)
	}
	return int(b>>4)*10 + int(b&0x0F), nil
}

func encodeBCD(v int) byte {
	return byte(v/10)<<4 | byte(v%10)
}

// ramTable is the WRSol RAM address map transcribed from the Weishaupt WRSol page.
var ramTable = []RAMRegister{
	{Address: 0xF484, Name: "TSO", Format: "Integer", Scale: "X10", Unit: "C", Description: "Tank Top Temperature (Speicher Oben Temperatur)"},
	{Address: 0xF488, Name: "TSU", Format: "Integer", Scale: "X10", Unit: "C", Description: "Tank Bottom Temperature (Speicher Unten Temperatur)"},
	{Address: 0xF48A, Name: "TWT / PWT", Format: "Integer", Scale: "X10", Unit: "C", Description: "DHW Plate Heat Exchanger Temperature (Brauchwasser Plattentauscher Temperatur)"},
	{Address: 0xF48C, Name: "TZW", Format: "Integer", Scale: "X10", Unit: "C", Description: "Circulation Return Temperature (Zirkulations Rucklauftemperatur)"},
	{Address: 0xF492, Name: "Drehzahl_BW_P1", Format: "Integer", Scale: "X10", Unit: "%", Description: "Collector Tank Pump Speed HV 50 (Drehzahl Kollektor Speicher Pumpe HV 50)"},
	{Address: 0xF494, Name: "Drehzahl_BW_P2", Format: "Integer", Scale: "X10", Unit: "%", Description: "Collector Buffer Pump Speed HV 50 (Drehzahl Kollektor Puffer Pumpe HV 50)"},
	{Address: 0xF496, Name: "Drehzahl_BW_WT", Format: "Integer", Scale: "X10", Unit: "%", Description: "DHW Plate Heat Exchanger Pump Speed (Drehzahl Brauchwasser Plattentauscher Pumpe)"},
	{Address: 0xF4C8, Name: "TSB", Format: "Integer", Scale: "X10", Unit: "C", Description: "Swimming Pool Temperature (Schwimmbad Temperatur)"},
	{Address: 0xF4F2, Name: "TPO", Format: "Integer", Scale: "X10", Unit: "C", Description: "Buffer Top Temperature (Puffer Oben Temperatur)"},
	{Address: 0xF4F4, Name: "TPU", Format: "Integer", Scale: "X10", Unit: "C", Description: "Buffer Bottom Temperature (Puffer Unten Temperatur)"},
	{Address: 0xF4F6, Name: "TWT", Format: "Integer", Scale: "X10", Unit: "C", Description: "Buffer Plate Heat Exchanger Temperature (Puffer Plattentauscher Temperatur)"},
	{Address: 0xF4F8, Name: "THR", Format: "Integer", Scale: "X10", Unit: "C", Description: "Heating Return Temperature (Heizungs Rucklauftemperatur)"},
	{Address: 0xF4FC, Name: "Drehzahl_Puf_WT", Format: "Integer", Scale: "X10", Unit: "%", Description: "Buffer Plate Heat Exchanger Pump Speed (Drehzahl Puffer Plattentauscher Pumpe)"},
	{Address: 0xF504, Name: "TWW", Format: "Integer", Scale: "X10", Unit: "C", Description: "Hot Water Temperature HV 14 (Warmwasser Temperatur HV 14)"},
	{Address: 0xF506, Name: "TKW", Format: "Integer", Scale: "X10", Unit: "C", Description: "Cold Water Temperature HV 14 (Kaltwasser Temperatur HV 14)"},
	{Address: 0xF508, Name: "Drehzahl_WW", Format: "Integer", Scale: "X10", Unit: "%", Description: "DHW Charging Pump Speed HV 14 (Drehzahl Warmwasser Ladepumpe HV 14)"},
	{Address: 0xF51C, Name: "TKO", Format: "Integer", Scale: "X10", Unit: "C", Description: "Collector 1 Temperature (Kollektor 1 Temperatur)"},
	{Address: 0xF530, Name: "Drehzahl_K1", Format: "Integer", Scale: "X10", Unit: "%", Description: "Collector 1 Pump Speed (Drehzahl Kollektor 1 Pumpe)"},
	{Address: 0xF532, Name: "MittlDrehz_1", Format: "Integer", Scale: "X10", Unit: "%", Description: "Average Collector 1 Pump Speed (Mittlere Drehzahl Kollektor 1 Pumpe)"},
	{Address: 0xF536, Name: "SolarLeistung1", Format: "Word", Scale: "X1", Unit: "kW", Description: "Solar Power Collector 1 (Solare Leistung Kollektor 1)"},
	{Address: 0xF53A, Name: "Gewinn1", Format: "Word", Scale: "X1", Unit: "kWh", Description: "Solar Gain Collector 1 (Solargewinn Kollektor 1)"},
	{Address: 0xF53C, Name: "GewinnTotLow1", Format: "Word", Scale: "X1", Unit: "kWh", Description: "Total Solar Gain Low Collector 1, 0--9999 kWh (Solargewinn total Kollektor 1)"},
	{Address: 0xF53E, Name: "GewinnTot1", Format: "Word", Scale: "X10", Unit: "MWh", Description: "Total Solar Gain High Collector 1, 0.1--999.9 MWh (Solargewinn total Kollektor 1)"},
	{Address: 0xF540, Name: "Bstd_Kol1", Format: "Integer", Scale: "X1", Unit: "h", Description: "Collector 1 Pump Operating Hours (Betriebsstunden Kollektor Pumpe 1)"},
	{Address: 0xF544, Name: "TK2", Format: "Integer", Scale: "X10", Unit: "C", Description: "Collector 2 Temperature (Kollektor 2 Temperatur)"},
	{Address: 0xF558, Name: "Drehzahl_K2", Format: "Integer", Scale: "X10", Unit: "%", Description: "Collector 2 Pump Speed (Drehzahl Kollektor 2 Pumpe)"},
	{Address: 0xF55A, Name: "MittlDrehz_2", Format: "Integer", Scale: "X10", Unit: "%", Description: "Average Collector 2 Pump Speed (Mittlere Drehzahl Kollektor 2 Pumpe)"},
	{Address: 0xF55E, Name: "SolarLeistung2", Format: "Word", Scale: "X1", Unit: "kW", Description: "Solar Power Collector 2 (Solare Leistung Kollektor 2)"},
	{Address: 0xF562, Name: "Gewinn2", Format: "Word", Scale: "X1", Unit: "kWh", Description: "Solar Gain Collector 2 (Solargewinn Kollektor 2)"},
	{Address: 0xF564, Name: "GewinnTotLow2", Format: "Word", Scale: "X1", Unit: "kWh", Description: "Total Solar Gain Low Collector 2, 0--9999 kWh (Solargewinn total Kollektor 2)"},
	{Address: 0xF566, Name: "GewinnTot2", Format: "Word", Scale: "X10", Unit: "MWh", Description: "Total Solar Gain High Collector 2, 0.1--999.9 MWh (Solargewinn total Kollektor 2)"},
	{Address: 0xF568, Name: "Bstd_Kol2", Format: "Integer", Scale: "X1", Unit: "h", Description: "Collector 2 Pump Operating Hours (Betriebsstunden Kollektor Pumpe 2)"},
	{Address: 0xF56C, Name: "TKV", Format: "Integer", Scale: "X10", Unit: "C", Description: "Collector Flow Temperature (Kollektor Vorlauftemperatur)"},
	{Address: 0xF56E, Name: "TKR", Format: "Integer", Scale: "X10", Unit: "C", Description: "Collector Return Temperature (Kollektor Rucklauftemperatur)"},
	{Address: 0xF570, Name: "TBY", Format: "Integer", Scale: "X10", Unit: "C", Description: "Bypass Temperature (Beipass Temperatur)"},
	{Address: 0xF572, Name: "STR", Format: "Integer", Scale: "X10", Unit: "C", Description: "Radiation Temperature (Strahlungs Temperatur)"},
	{Address: 0xF58A, Name: "TKOL_hoechst", Format: "Integer", Scale: "X10", Unit: "C", Description: "Collector Maximum Temperature (Kollektor Hochsttemperatur)"},
	{Address: 0xF58C, Name: "VolumStrom", Format: "Integer", Scale: "X10", Unit: "l/h", Description: "Volume Flow (Volumenstrom)"},
	{Address: 0xF4CA, Name: "TFK", Format: "Integer", Scale: "X10", Unit: "C", Description: "Solid Fuel Boiler Flow Temperature (Feststoffkessel Vorlauftemperatur)"},
	{Address: 0xF4CC, Name: "TFKR", Format: "Integer", Scale: "X10", Unit: "C", Description: "Solid Fuel Boiler Return Temperature (Feststoffkessel Rucklauftemperatur)"},
	{Address: 0xF5B4, Name: "WE_sperre", Format: "Integer", Scale: "X10", Unit: "C", Description: "Heat Generator Lockout (Warmeerzeuger Sperre)"},
	{Address: 0xFC00, Name: "ActualTime", Format: "T, hh, mm, ss", Scale: "--", Unit: "", Description: "Current Time and Weekday (Aktuelle Uhrzeit und Wochentag)"},
	{Address: 0xF430, Name: "ErrorBits", Format: "Word", Scale: "--", Unit: "", Description: "Error Register, Bit 0--15 (Fehlerregister)"},
}
