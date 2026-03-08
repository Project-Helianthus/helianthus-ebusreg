package b555

import (
	"github.com/Project-Helianthus/helianthus-ebusgo/types"
)

// decodeConfigResponse decodes an A3 CONFIG_READ response (9 bytes).
//
// Response format:
//
//	[status] [max_slots] [time_res] [min_dur] [has_temp] [temp_slots] [min_temp] [max_temp] [pad]
func decodeConfigResponse(zone, hc byte, payload []byte) map[string]types.Value {
	values := map[string]types.Value{
		"opcode":  {Value: opcodeConfigRead, Valid: true},
		"zone":    {Value: zone, Valid: true},
		"hc":      {Value: hc, Valid: true},
		"payload": {Value: append([]byte(nil), payload...), Valid: true},
	}

	if len(payload) < 9 {
		values["status"] = types.Value{Valid: false}
		return values
	}

	status := payload[0]
	values["status"] = types.Value{Value: status, Valid: true}
	values["max_slots"] = types.Value{Value: payload[1], Valid: true}
	values["time_resolution"] = types.Value{Value: payload[2], Valid: true}
	values["min_duration"] = types.Value{Value: payload[3], Valid: true}
	values["has_temperature"] = types.Value{Value: payload[4] == 0x01, Valid: true}
	values["temp_slots"] = types.Value{Value: payload[5], Valid: true}

	if payload[6] == 0xFF {
		values["min_temp_c"] = types.Value{Valid: false}
	} else {
		values["min_temp_c"] = types.Value{Value: float64(payload[6]), Valid: true}
	}

	if payload[7] == 0xFF {
		values["max_temp_c"] = types.Value{Valid: false}
	} else {
		values["max_temp_c"] = types.Value{Value: float64(payload[7]), Valid: true}
	}

	// status 0x03 = unavailable
	values["available"] = types.Value{Value: status == 0x00, Valid: true}

	return values
}

// decodeSlotsResponse decodes an A4 SLOTS_READ response (9 bytes).
//
// Response format:
//
//	[status] [Mon] [Tue] [Wed] [Thu] [Fri] [Sat] [Sun] [pad]
func decodeSlotsResponse(zone, hc byte, payload []byte) map[string]types.Value {
	values := map[string]types.Value{
		"opcode":  {Value: opcodeSlotsRead, Valid: true},
		"zone":    {Value: zone, Valid: true},
		"hc":      {Value: hc, Valid: true},
		"payload": {Value: append([]byte(nil), payload...), Valid: true},
	}

	if len(payload) < 9 {
		values["status"] = types.Value{Valid: false}
		return values
	}

	status := payload[0]
	values["status"] = types.Value{Value: status, Valid: true}
	values["available"] = types.Value{Value: status == 0x00, Valid: true}

	// Slot counts per weekday: Mon=0x00..Sun=0x06
	days := make([]int, 7)
	for i := 0; i < 7; i++ {
		days[i] = int(payload[1+i])
	}
	values["slots_per_day"] = types.Value{Value: days, Valid: true}

	return values
}

// decodeTimerReadResponse decodes an A5 TIMER_READ response (7 bytes).
//
// Response format:
//
//	[status] [Sh] [Sm] [Eh] [Em] [Tlo] [Thi]
func decodeTimerReadResponse(zone, hc, weekday, slot byte, payload []byte) map[string]types.Value {
	values := map[string]types.Value{
		"opcode":  {Value: opcodeTimerRead, Valid: true},
		"zone":    {Value: zone, Valid: true},
		"hc":      {Value: hc, Valid: true},
		"weekday": {Value: weekday, Valid: true},
		"slot":    {Value: slot, Valid: true},
		"payload": {Value: append([]byte(nil), payload...), Valid: true},
	}

	if len(payload) < 7 {
		values["status"] = types.Value{Valid: false}
		return values
	}

	status := payload[0]
	values["status"] = types.Value{Value: status, Valid: true}
	values["available"] = types.Value{Value: status == 0x00, Valid: true}
	values["start_hour"] = types.Value{Value: payload[1], Valid: true}
	values["start_minute"] = types.Value{Value: payload[2], Valid: true}
	values["end_hour"] = types.Value{Value: payload[3], Valid: true}
	values["end_minute"] = types.Value{Value: payload[4], Valid: true}

	tempRaw := uint16(payload[5]) | uint16(payload[6])<<8
	if tempRaw == 0xFFFF {
		values["temperature_raw"] = types.Value{Value: uint16(0xFFFF), Valid: true}
		values["temperature_c"] = types.Value{Valid: false}
	} else {
		values["temperature_raw"] = types.Value{Value: tempRaw, Valid: true}
		values["temperature_c"] = types.Value{Value: float64(tempRaw) / 10.0, Valid: true}
	}

	return values
}

// decodeTimerWriteResponse decodes an A6 TIMER_WRITE response (1 byte).
//
// Response format:
//
//	[error_code]
//
// Error codes: 0x00=ACK, 0x01=parameter out of range, 0x03=timer unavailable,
// 0x06=validation failure.
func decodeTimerWriteResponse(zone, hc, weekday, slot byte, payload []byte) map[string]types.Value {
	values := map[string]types.Value{
		"opcode":  {Value: opcodeTimerWrite, Valid: true},
		"zone":    {Value: zone, Valid: true},
		"hc":      {Value: hc, Valid: true},
		"weekday": {Value: weekday, Valid: true},
		"slot":    {Value: slot, Valid: true},
		"payload": {Value: append([]byte(nil), payload...), Valid: true},
	}

	if len(payload) == 0 {
		values["error_code"] = types.Value{Valid: false}
		values["accepted"] = types.Value{Value: false, Valid: true}
		return values
	}

	errorCode := payload[0]
	values["error_code"] = types.Value{Value: errorCode, Valid: true}
	values["accepted"] = types.Value{Value: errorCode == 0x00, Valid: true}

	return values
}
