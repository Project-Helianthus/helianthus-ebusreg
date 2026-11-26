package owner

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
)

var vaillantStructuralStartupProbeTargets = []byte{0x08, 0x15, 0x26}

var registryScanDirectedFn = registry.ScanDirected

func snapshotTargetAddressForRouting(snap registry.DeviceEntrySnapshot) byte {
	if addr, ok := snap.AddressByRole(registry.SlotRoleSlave); ok {
		return addr
	}
	return snap.PrimaryDisplayAddress()
}

func snapshotContainsAddress(snap registry.DeviceEntrySnapshot, addr byte) bool {
	for _, a := range snap.Addresses {
		if a == addr {
			return true
		}
	}
	return false
}

func isB524ResponseCoherent(responseData []byte, group byte, addr uint16) bool {
	if len(responseData) < 4 {
		return false
	}
	if len(responseData) >= 5 {
		if responseData[2] == group && (uint16(responseData[3])|uint16(responseData[4])<<8) == addr {
			return true
		}
	}
	return responseData[1] == group && (uint16(responseData[2])|uint16(responseData[3])<<8) == addr
}

func cloneStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	cp := *value
	return &cp
}

func countRegistryDevices(reg *registry.DeviceRegistry) int {
	if reg == nil {
		return 0
	}
	count := 0
	reg.Iterate(func(entry registry.DeviceEntry) bool {
		count++
		return true
	})
	return count
}

func ebusdScanTargetCandidates(config engine.TransportConfig) []engine.TransportConfig {
	if config.Protocol != engine.TransportEbusdTCP {
		return nil
	}
	if !strings.EqualFold(config.Network, "tcp") || config.Address == "" {
		return nil
	}
	return []engine.TransportConfig{config}
}

func dialContext(ctx context.Context, network, address string, timeout time.Duration) (net.Conn, error) {
	dialer := net.Dialer{Timeout: timeout}
	return dialer.DialContext(ctx, network, address)
}

func circuitTemperatureNS(v semreg.Uint64) uint64 {
	n, _ := strconv.ParseUint(string(v), 10, 64)
	return n
}
