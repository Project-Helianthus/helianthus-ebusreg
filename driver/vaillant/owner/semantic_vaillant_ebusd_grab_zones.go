package owner

import (
	"encoding/hex"
	"errors"
	"math"
	"net"
	"strings"
)

// finishB524GrabZone derives the circuit, mode, preset, action and name of
// one accumulated zone and reports whether the zone is present.
func finishB524GrabZone(instance byte, state *b524GrabZoneState, circuitTypeByInstance map[byte]*uint16) (*vaillantZoneSnapshot, bool) {
	if state == nil || state.snapshot == nil {
		return nil, false
	}
	circuitInstance := resolveAssociatedCircuitInstance(state.roomTemperatureZoneMappingRaw, instance)
	if state.roomTemperatureZoneMappingRaw != nil {
		typed := uint16(circuitInstance)
		state.snapshot.ConfigurationAssociatedCircuitRaw = &typed
	}
	circuitType, hasCircuitType := circuitTypeByInstance[circuitInstance]
	if hasCircuitType && circuitType != nil {
		state.snapshot.ConfigurationCircuitTypeRaw = circuitType
	}

	mode, preset, allowedModes := deriveZoneModeAndPreset(state.configurationOperationMode, state.stateSpecialFunction, circuitType, hasCircuitType)
	hvacAction := deriveZoneHvacAction(state.stateValveStatus, circuitType, hasCircuitType)
	if mode != "" {
		state.snapshot.OperatingMode = mode
	}
	if preset != "" {
		state.snapshot.Preset = preset
	}
	if hvacAction != "" {
		state.snapshot.HvacAction = hvacAction
	}
	if len(allowedModes) > 0 {
		state.snapshot.AllowedModes = allowedModes
	}

	name := composeZoneName(state.nameParts.primary, state.nameParts.prefix, state.nameParts.suffix)
	if strings.TrimSpace(name) != "" {
		state.snapshot.Name = name
	}
	meaningful := strings.TrimSpace(state.snapshot.Name) != "" ||
		state.snapshot.OperatingMode != "" ||
		state.snapshot.Preset != "" ||
		state.snapshot.CurrentTempC != nil ||
		state.snapshot.TargetTempC != nil ||
		state.snapshot.HumidityPct != nil
	if state.hasIndex {
		if !state.indexPresent {
			return nil, false
		}
	} else if !meaningful {
		return nil, false
	}
	state.snapshot.Present = true
	return state.snapshot, true
}

func parseB524DhwFromGrab(lines []string, controller byte) (*vaillantDhwSnapshot, bool) {
	var (
		opModeRaw *uint16
		sfModeRaw *uint16
		current   *float64
		target    *float64
	)

	for _, line := range lines {
		group, instance, addr, payload, ok := parseB524GrabLine(line, controller)
		if !ok || group != localDHW.group || instance != dhwInstance {
			continue
		}
		switch addr {
		case dhw_operation_mode:
			if value, ok := decodeB524Uint16(payload); ok {
				typed := value
				opModeRaw = &typed
			}
		case dhw_special_function:
			if value, ok := decodeB524Uint16(payload); ok {
				typed := value
				sfModeRaw = &typed
			}
		case dhw_current_temp:
			if value, ok := decodeFloat32LE(payload); ok {
				typed := value
				current = &typed
			}
		case dhw_target_temp:
			if value, ok := decodeFloat32LE(payload); ok {
				typed := value
				target = &typed
			}
		}
	}

	operatingMode, preset := deriveDhwModeAndPreset(opModeRaw, sfModeRaw)
	if operatingMode == "" && preset == "" && current == nil && target == nil {
		return nil, false
	}

	dhw := &vaillantDhwSnapshot{
		OperatingMode: operatingMode,
		Preset:        preset,
		CurrentTempC:  current,
		TargetTempC:   target,
	}
	if opModeRaw != nil {
		dhw.ConfigurationDHWOperationMode = formatUintToken(*opModeRaw)
	}
	if sfModeRaw != nil {
		dhw.StateSpecialFunction = formatUintToken(*sfModeRaw)
	}
	return dhw, true
}

func parseB524GrabLine(line string, controller byte) (byte, byte, uint16, []byte, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return 0, 0, 0, nil, false
	}

	slash := strings.Index(line, "/")
	if slash <= 0 {
		return 0, 0, 0, nil, false
	}

	reqText := strings.TrimSpace(line[:slash])
	reqFields := strings.Fields(reqText)
	if len(reqFields) == 0 {
		return 0, 0, 0, nil, false
	}
	reqBytes, ok := decodeHexField(reqFields[0])
	if !ok || len(reqBytes) < 11 {
		return 0, 0, 0, nil, false
	}
	if reqBytes[2] != vaillantExtRegisterPrimary || reqBytes[3] != vaillantExtRegisterSecondary {
		return 0, 0, 0, nil, false
	}
	readOpcode := reqBytes[5]
	if reqBytes[4] != vaillantB524OpcodeRead ||
		!isB524GrabReadOpcode(readOpcode) ||
		reqBytes[6] != vaillantB524OpRead {
		return 0, 0, 0, nil, false
	}

	target := reqBytes[1]
	if controller != 0 && target != controller {
		return 0, 0, 0, nil, false
	}

	group := reqBytes[7]
	instance := reqBytes[8]
	addr := uint16(reqBytes[9]) | uint16(reqBytes[10])<<8

	respPart := strings.TrimSpace(line[slash+1:])
	if respPart == "" {
		return 0, 0, 0, nil, false
	}
	if equal := strings.Index(respPart, "="); equal >= 0 {
		respPart = strings.TrimSpace(respPart[:equal])
	}
	respFields := strings.Fields(respPart)
	if len(respFields) == 0 {
		return 0, 0, 0, nil, false
	}
	if strings.HasPrefix(strings.ToLower(respFields[0]), "err") {
		return 0, 0, 0, nil, false
	}

	respBytes, ok := decodeHexField(respFields[0])
	if !ok || len(respBytes) == 0 {
		return 0, 0, 0, nil, false
	}
	if int(respBytes[0]) == len(respBytes)-1 {
		respBytes = respBytes[1:]
	}

	payload, ok := parseB524ReadPayload(respBytes, readOpcode, group, instance, addr)
	if !ok {
		return 0, 0, 0, nil, false
	}
	return group, instance, addr, payload, true
}

func isB524GrabReadOpcode(opcode byte) bool {
	return opcode == vaillantB524OpcodeLocal || opcode == vaillantB524OpcodeRead
}

func decodeHexField(value string) ([]byte, bool) {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return nil, false
	}
	normalized = strings.ReplaceAll(normalized, " ", "")
	normalized = strings.ReplaceAll(normalized, "\t", "")
	if len(normalized)%2 != 0 {
		return nil, false
	}
	bytes, err := hex.DecodeString(normalized)
	if err != nil {
		return nil, false
	}
	return bytes, true
}

func decodeCString(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	trimmed := raw
	for index, value := range trimmed {
		if value == 0x00 {
			trimmed = trimmed[:index]
			break
		}
	}
	return strings.TrimSpace(string(trimmed))
}

func decodeFloat32LE(raw []byte) (float64, bool) {
	if len(raw) < 4 {
		return 0, false
	}
	bits := uint32(raw[0]) | uint32(raw[1])<<8 | uint32(raw[2])<<16 | uint32(raw[3])<<24
	value := float64(math.Float32frombits(bits))
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, false
	}
	return value, true
}

func isNetTimeout(err error) bool {
	if err == nil {
		return false
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
