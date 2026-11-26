package owner

import (
	"fmt"
	"sync"

	"github.com/Project-Helianthus/helianthus-ebusreg/driver/vaillant/output"
)

type recordingAdapterSink struct {
	mu   sync.Mutex
	last *output.AdapterHardwareInfo
}

func (s *recordingAdapterSink) publishAdapterInfo(info vaillantAdapterInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last = adapterHardwareInfoFromNative(info)
}

func (s *recordingAdapterSink) AdapterHardwareInfo() *output.AdapterHardwareInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.last
}

func adapterHardwareInfoFromNative(native vaillantAdapterInfo) *output.AdapterHardwareInfo {
	info := &output.AdapterHardwareInfo{}

	if native.identity != nil {
		v := native.identity
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

	info.HardwareID = native.hardwareID
	info.HardwareConfig = native.hardwareConfig
	info.TemperatureC = native.temperatureC
	info.SupplyVoltageMV = native.supplyMV
	info.BusVoltageMaxDV = native.busMaxDV
	info.BusVoltageMinDV = native.busMinDV
	info.ResetCause = native.resetCause
	info.ResetCauseCode = native.resetCode
	info.RestartCount = native.restartCount
	info.WiFiRSSIDBm = native.wifiRSSI

	if !native.lastIdentity.IsZero() {
		t := native.lastIdentity
		info.LastIdentityQuery = &t
	}
	if !native.lastTelemetry.IsZero() {
		t := native.lastTelemetry
		info.LastTelemetryQuery = &t
	}
	return info
}
