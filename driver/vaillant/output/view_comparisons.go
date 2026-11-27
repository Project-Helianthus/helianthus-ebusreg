package output

import (
	"fmt"
	"slices"
)

// AdapterHardwareInfoFromNative maps the native values to the corresponding output record.
func AdapterHardwareInfoFromNative(native NativeAdapterInfo) *AdapterHardwareInfo {
	info := &AdapterHardwareInfo{}

	if native.Identity != nil {
		v := native.Identity
		info.FirmwareVersion = fmt.Sprintf("0x%02x", v.Version)
		info.Features = v.Features
		info.Jumpers = v.Jumpers
		info.IsWiFi = v.IsWiFi
		info.IsEthernet = v.IsEthernet
		info.VersionResponseLen = v.VersionResponseLen()
		info.InfoSupported = v.SupportsInfo

		if v.HasChecksum {
			info.FirmwareChecksum = fmt.Sprintf("0x%04x", v.Checksum)
		}
		if v.HasBootloader {
			info.BootloaderVersion = fmt.Sprintf("0x%02x", v.BootloaderVersion)
			info.BootloaderChecksum = fmt.Sprintf("0x%04x", v.BootloaderChecksum)
		}

		var flags []string
		if v.Jumpers&0x01 != 0 {
			flags = append(flags, "enhanced")
		}
		if v.IsHighSpeed {
			flags = append(flags, "high_speed")
		}
		if v.IsEthernet {
			flags = append(flags, "ethernet")
		}
		if v.IsWiFi {
			flags = append(flags, "wifi")
		}
		if v.IsV31 {
			flags = append(flags, "v3.1")
		}
		if v.Jumpers&0x20 != 0 {
			flags = append(flags, "soft_config")
		}
		info.JumperFlags = flags
	}

	info.HardwareID = native.HardwareID
	info.HardwareConfig = native.HardwareConfig
	info.TemperatureC = native.TemperatureC
	info.SupplyVoltageMV = native.SupplyMV
	info.BusVoltageMaxDV = native.BusMaxDV
	info.BusVoltageMinDV = native.BusMinDV
	info.ResetCause = native.ResetCause
	info.ResetCauseCode = native.ResetCode
	info.RestartCount = native.RestartCount
	info.WiFiRSSIDBm = native.WifiRSSI

	if !native.LastIdentity.IsZero() {
		t := native.LastIdentity
		info.LastIdentityQuery = &t
	}
	if !native.LastTelemetry.IsZero() {
		t := native.LastTelemetry
		info.LastTelemetryQuery = &t
	}
	return info
}

func enumerationsEqual(a, b []CatalogueEnumeration) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Group != b[i].Group || a[i].Register != b[i].Register || a[i].Alternative != b[i].Alternative || a[i].Label != b[i].Label || a[i].Level != b[i].Level || !slices.Equal(a[i].Raw, b[i].Raw) || !intPtrEquals(a[i].Code, b[i].Code) {
			return false
		}
	}
	return true
}

func configBoundEquals(a, b *ConfigBound) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

// CircuitsEqual compares published circuits values.
func CircuitsEqual(a, b []CircuitStatus) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		left := a[i]
		right := b[i]
		if left.Index != right.Index || left.CircuitType != right.CircuitType || left.HasMixer != right.HasMixer {
			return false
		}
		if !boolPtrEquals(left.State.PumpActive, right.State.PumpActive) {
			return false
		}
		if !floatPtrEquals(left.State.MixerPositionPct, right.State.MixerPositionPct) {
			return false
		}
		if !floatPtrEquals(left.State.FlowTemperatureC, right.State.FlowTemperatureC) {
			return false
		}
		if !floatPtrEquals(left.State.FlowSetpointC, right.State.FlowSetpointC) {
			return false
		}
		if !floatPtrEquals(left.State.CalcFlowTempC, right.State.CalcFlowTempC) {
			return false
		}
		if left.State.CircuitState != right.State.CircuitState {
			return false
		}
		if !floatPtrEquals(left.State.Humidity, right.State.Humidity) {
			return false
		}
		if !floatPtrEquals(left.State.DewPoint, right.State.DewPoint) {
			return false
		}
		if !floatPtrEquals(left.State.PumpHours, right.State.PumpHours) {
			return false
		}
		if !intPtrEquals(left.State.PumpStarts, right.State.PumpStarts) {
			return false
		}
		if !floatPtrEquals(left.Config.HeatingCurve, right.Config.HeatingCurve) {
			return false
		}
		if !floatPtrEquals(left.Config.FlowTempMaxC, right.Config.FlowTempMaxC) {
			return false
		}
		if !floatPtrEquals(left.Config.FlowTempMinC, right.Config.FlowTempMinC) {
			return false
		}
		if !floatPtrEquals(left.Config.SummerLimitC, right.Config.SummerLimitC) {
			return false
		}
		if !intPtrEquals(left.Config.SpecialFunctionMode, right.Config.SpecialFunctionMode) {
			return false
		}
		if left.Config.RoomTempControl != right.Config.RoomTempControl {
			return false
		}
		if !boolPtrEquals(left.Config.CoolingEnabled, right.Config.CoolingEnabled) {
			return false
		}
		if !enumerationsEqual(left.Config.enumerations, right.Config.enumerations) {
			return false
		}
		if left.ManagingDevice.Role != right.ManagingDevice.Role {
			return false
		}
		if !stringPtrEquals(left.ManagingDevice.DeviceID, right.ManagingDevice.DeviceID) {
			return false
		}
		if !intPtrEquals(left.ManagingDevice.Address, right.ManagingDevice.Address) {
			return false
		}
	}
	return true
}

// DeriveCircuitManagingDevice maps scheme and mode evidence to the published managing device.
func DeriveCircuitManagingDevice(system *NativeSystemSnapshot, moduleMode NativeFunctionModuleMode) ManagingDevice {
	return circuitManagingDevice(CircuitsManagedByFunctionModule(system, moduleMode))
}

func circuitManagingDevice(managedByFunctionModule bool) ManagingDevice {
	if managedByFunctionModule {
		deviceID := functionModuleManagingDeviceID
		address := functionModuleManagingDeviceAddress
		return ManagingDevice{
			Role:     ManagingDeviceRoleFunctionModule,
			DeviceID: &deviceID,
			Address:  &address,
		}
	}
	return ManagingDevice{Role: ManagingDeviceRoleUnknown}
}

// DhwEquals compares published dhw values.
func DhwEquals(a, b *DhwStatus) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Config.OperatingMode != b.Config.OperatingMode || a.Config.Preset != b.Config.Preset {
		return false
	}
	if a.State.SpecialFunction != b.State.SpecialFunction {
		return false
	}
	if !enumerationsEqual(a.State.enumerations, b.State.enumerations) {
		return false
	}
	if !floatPtrEquals(a.State.CurrentTempC, b.State.CurrentTempC) {
		return false
	}
	if !floatPtrEquals(a.Config.TargetTempC, b.Config.TargetTempC) {
		return false
	}
	if !floatPtrEquals(a.State.HeatingDemandPct, b.State.HeatingDemandPct) {
		return false
	}
	if a.Config.HolidayStartDate != b.Config.HolidayStartDate || a.Config.HolidayEndDate != b.Config.HolidayEndDate {
		return false
	}
	return true
}

// RadioDevicesEqual compares published radio devices values.
func RadioDevicesEqual(a, b []RadioDevice) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		left := a[i]
		right := b[i]
		if left.Group != right.Group || left.Instance != right.Instance || left.SlotMode != right.SlotMode {
			return false
		}
		if !boolPtrEquals(left.DeviceConnected, right.DeviceConnected) {
			return false
		}
		if !intPtrEquals(left.DeviceClassAddress, right.DeviceClassAddress) {
			return false
		}
		if left.DeviceModel != right.DeviceModel {
			return false
		}
		if !stringPtrEquals(left.FirmwareVersion, right.FirmwareVersion) {
			return false
		}
		if !intPtrEquals(left.HardwareIdentifier, right.HardwareIdentifier) {
			return false
		}
		if !intPtrEquals(left.RemoteControlAddress, right.RemoteControlAddress) {
			return false
		}
		if !boolPtrEquals(left.DevicePaired, right.DevicePaired) {
			return false
		}
		if !intPtrEquals(left.ReceptionStrength, right.ReceptionStrength) {
			return false
		}
		if !intPtrEquals(left.ZoneAssignment, right.ZoneAssignment) {
			return false
		}
		if !floatPtrEquals(left.RoomTemperatureC, right.RoomTemperatureC) {
			return false
		}
		if !floatPtrEquals(left.RoomHumidityPct, right.RoomHumidityPct) {
			return false
		}
	}
	return true
}

// ZoneEquals compares published zone values.
func ZoneEquals(a, b Zone) bool {
	if a.ID != b.ID || a.Name != b.Name {
		return false
	}

	if a.State.HvacAction != b.State.HvacAction || a.State.SpecialFunction != b.State.SpecialFunction {
		return false
	}
	if !enumerationsEqual(a.State.enumerations, b.State.enumerations) {
		return false
	}
	if !floatPtrEquals(a.State.CurrentTempC, b.State.CurrentTempC) {
		return false
	}
	if !floatPtrEquals(a.State.CurrentHumidityPct, b.State.CurrentHumidityPct) {
		return false
	}
	if !floatPtrEquals(a.State.HeatingDemandPct, b.State.HeatingDemandPct) {
		return false
	}
	if !floatPtrEquals(a.State.ValvePositionPct, b.State.ValvePositionPct) {
		return false
	}

	if a.Config.OperatingMode != b.Config.OperatingMode || a.Config.Preset != b.Config.Preset {
		return false
	}
	if !floatPtrEquals(a.Config.TargetTempC, b.Config.TargetTempC) {
		return false
	}
	if !slices.Equal(a.Config.AllowedModes, b.Config.AllowedModes) {
		return false
	}
	if a.Config.CircuitType != b.Config.CircuitType {
		return false
	}
	if !intPtrEquals(a.Config.AssociatedCircuit, b.Config.AssociatedCircuit) {
		return false
	}
	if !intPtrEquals(a.Config.RoomTemperatureZoneMapping, b.Config.RoomTemperatureZoneMapping) {
		return false
	}
	if a.Config.QuickVeto != b.Config.QuickVeto || a.Config.QuickVetoExpiry != b.Config.QuickVetoExpiry {
		return false
	}
	if !floatPtrEquals(a.Config.QuickVetoSetpointC, b.Config.QuickVetoSetpointC) {
		return false
	}
	if !floatPtrEquals(a.Config.QuickVetoDurationH, b.Config.QuickVetoDurationH) {
		return false
	}
	if a.Config.HolidayStartDate != b.Config.HolidayStartDate || a.Config.HolidayEndDate != b.Config.HolidayEndDate {
		return false
	}
	if !floatPtrEquals(a.Config.HolidaySetpointC, b.Config.HolidaySetpointC) {
		return false
	}
	if a.Config.HolidayStartTime != b.Config.HolidayStartTime || a.Config.HolidayEndTime != b.Config.HolidayEndTime {
		return false
	}
	return true
}
