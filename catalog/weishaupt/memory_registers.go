package weishaupt

import "fmt"

// Primary and secondary bytes of the Kromschroeder memory access.
const (
	PBKromschroeder byte = 0x50
	SBMemoryRead    byte = 0x00
	SBMemoryWrite   byte = 0x01
)

// MemoryID is the three-byte obfuscated ID of a memory register: B0, B1 (the
// high nibble of the page), and B2 (the address low byte). The page says the
// algorithm for B0 is not recoverable from the public material, so an ID is
// given by the caller, never derived from an address.
type MemoryID [3]byte

// MemoryRegister is one selected RAM register of the page. Its width and unit
// are not stated on the page; a reply is therefore decoded raw.
type MemoryRegister struct {
	Address     uint16
	Name        string
	Description string
}

// MemoryRegisterByAddress returns the selected register at an address.
func MemoryRegisterByAddress(address uint16) (MemoryRegister, bool) {
	for _, reg := range memoryRegisterTable {
		if reg.Address == address {
			return reg, true
		}
	}
	return MemoryRegister{}, false
}

// observedMemoryIDs lists the ID bytes the page shows as examples. They are not a
// rule for other addresses.
var observedMemoryIDs = map[uint16]MemoryID{
	0x0004: {0x58, 0x01, 0x04},
	0x0005: {0x59, 0x01, 0x05},
	0x0020: {0x7C, 0x01, 0x20},
	0x0080: {0xDC, 0x01, 0x80},
	0x00D0: {0x68, 0x02, 0xD0},
}

// ObservedMemoryID returns the ID the page shows for an address as an example.
func ObservedMemoryID(address uint16) (MemoryID, bool) {
	id, ok := observedMemoryIDs[address]
	return id, ok
}

// EncodeMemoryRead builds the 5000 request: the three ID bytes.
func EncodeMemoryRead(id MemoryID) []byte {
	return append([]byte(nil), id[:]...)
}

// DecodeMemoryReply keeps the 5001 read reply raw. The page gives a width of one
// or two bytes; any other length is an error.
func DecodeMemoryReply(payload []byte) (Field, error) {
	if len(payload) < 1 || len(payload) > 2 {
		return Field{}, fmt.Errorf("%w: memory reply wants 1 or 2 bytes, got %d", ErrWidth, len(payload))
	}
	return Field{Name: "value_raw", Value: append([]byte(nil), payload...)}, nil
}

// EncodeMemoryWrite builds the 5001 write request: the three ID bytes followed
// by a value of one or two bytes, as the page states.
func EncodeMemoryWrite(id MemoryID, value []byte) ([]byte, error) {
	if len(value) < 1 || len(value) > 2 {
		return nil, fmt.Errorf("%w: memory write value wants 1 or 2 bytes, got %d", ErrValueSize, len(value))
	}
	out := make([]byte, 0, 3+len(value))
	out = append(out, id[:]...)
	return append(out, value...), nil
}

// memoryRegisterTable is the selected Kromschroeder RAM register set transcribed from the Kromschroeder page.
var memoryRegisterTable = []MemoryRegister{
	{Address: 0x0004, Name: "Divi", Description: "Divider value"},
	{Address: 0x0005, Name: "Aufruf_von", Description: "State machine caller"},
	{Address: 0x0006, Name: "Bezeichnung_INr", Description: "Designation / internal number"},
	{Address: 0x0007, Name: "Drehzahl", Description: "Fan/blower RPM"},
	{Address: 0x0023, Name: "ERR_MERK", Description: "Error marker"},
	{Address: 0x0026, Name: "EBUS_BITS", Description: "eBUS status bits"},
	{Address: 0x0055, Name: "ICL20SEC", Description: "20-second interval clock"},
	{Address: 0x0056, Name: "ICL1HRS", Description: "1-hour interval clock"},
	{Address: 0x005C, Name: "TVISTAD", Description: "Flow temperature actual (damped)"},
	{Address: 0x005D, Name: "TWWISTAD", Description: "DHW temperature actual (damped)"},
	{Address: 0x005E, Name: "TAUSSEN", Description: "Outside temperature"},
	{Address: 0x005F, Name: "TABGAS", Description: "Exhaust gas temperature"},
	{Address: 0x0068, Name: "TVIST", Description: "Flow temperature instant"},
	{Address: 0x006A, Name: "TWWIST", Description: "DHW temperature instant"},
	{Address: 0x006E, Name: "TVSOLL", Description: "Flow temperature setpoint"},
	{Address: 0x006F, Name: "TWWSOLL_HS", Description: "DHW setpoint (heat exchanger side)"},
	{Address: 0x0060, Name: "IONIST", Description: "Ionisation current actual"},
	{Address: 0x0061, Name: "PWMGPV", Description: "PWM pump output"},
	{Address: 0x0070, Name: "NGSOLL", Description: "Fan speed setpoint"},
	{Address: 0x0071, Name: "NGIST", Description: "Fan speed actual"},
	{Address: 0x0072, Name: "HZZUSTBM", Description: "Heating circuit state bitmap"},
	{Address: 0x0074, Name: "WWZUSTBM", Description: "DHW circuit state bitmap"},
	{Address: 0x0081, Name: "HZ_HYS", Description: "Heating hysteresis"},
	{Address: 0x0082, Name: "TBRMINOFF", Description: "Minimum off-time"},
	{Address: 0x0083, Name: "HZ_KP", Description: "Heating PID proportional gain"},
	{Address: 0x0084, Name: "HZ_KTND", Description: "Heating PID derivative"},
	{Address: 0x0085, Name: "NG_KESSEL_MIN", Description: "Min fan speed"},
	{Address: 0x0087, Name: "NG_HZ_MAX", Description: "Max fan speed (heating)"},
	{Address: 0x008B, Name: "WW_HYS", Description: "DHW hysteresis"},
	{Address: 0x008D, Name: "NG_WW_MAX", Description: "Max fan speed (DHW)"},
	{Address: 0x0094, Name: "NG_START_PWM", Description: "Start PWM value"},
	{Address: 0x009B, Name: "HK_S", Description: "Heating curve slope"},
	{Address: 0x009D, Name: "NG_KESSEL_MAX", Description: "Max fan speed (global)"},
}
