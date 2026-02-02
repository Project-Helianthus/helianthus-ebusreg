package registry

import (
	"context"
	"errors"
	"fmt"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
	"github.com/d3vi1/helianthus-ebusgo/protocol"
)

const (
	scanPrimary   = byte(0x07)
	scanSecondary = byte(0x04)
)

type ScanBus interface {
	Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error)
}

// Scan performs a 07 04 identification scan over the provided targets.
func Scan(ctx context.Context, bus ScanBus, registry *DeviceRegistry, source byte, targets []byte) ([]DeviceEntry, error) {
	if bus == nil {
		return nil, fmt.Errorf("scan missing bus: %w", ebuserrors.ErrInvalidPayload)
	}
	if registry == nil {
		return nil, fmt.Errorf("scan missing registry: %w", ebuserrors.ErrInvalidPayload)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if len(targets) == 0 {
		targets = DefaultScanTargets()
	}

	entries := make([]DeviceEntry, 0)
	for _, target := range targets {
		frameType := protocol.FrameTypeForTarget(target)
		if frameType == protocol.FrameTypeMasterMaster || frameType == protocol.FrameTypeUnknown {
			continue
		}
		request := protocol.Frame{
			Source:    source,
			Target:    target,
			Primary:   scanPrimary,
			Secondary: scanSecondary,
		}
		response, err := bus.Send(ctx, request)
		if err != nil {
			if errors.Is(err, ebuserrors.ErrNoSuchDevice) {
				continue
			}
			return nil, fmt.Errorf("scan target %02x: %w", target, err)
		}
		if response == nil {
			return nil, fmt.Errorf("scan target %02x empty response: %w", target, ebuserrors.ErrInvalidPayload)
		}

		address := response.Source
		if address == 0 {
			address = target
		}
		info, err := parseDeviceInfo(address, response.Data)
		if err != nil {
			return nil, fmt.Errorf("scan target %02x parse: %w", target, err)
		}
		entries = append(entries, registry.Register(info))
	}

	return entries, nil
}

// DefaultScanTargets returns the default address range for scanning.
func DefaultScanTargets() []byte {
	targets := make([]byte, 0, 0xFD)
	for address := 0x01; address < 0xFE; address++ {
		targets = append(targets, byte(address))
	}
	return targets
}

func parseDeviceInfo(address byte, payload []byte) (DeviceInfo, error) {
	if len(payload) < 8 {
		return DeviceInfo{}, fmt.Errorf("short device info payload: %w", ebuserrors.ErrInvalidPayload)
	}
	return DeviceInfo{
		Address:         address,
		Manufacturer:    fmt.Sprintf("%02X%02X", payload[0], payload[1]),
		DeviceID:        fmt.Sprintf("%02X%02X", payload[2], payload[3]),
		SoftwareVersion: fmt.Sprintf("%02X%02X", payload[4], payload[5]),
		HardwareVersion: fmt.Sprintf("%02X%02X", payload[6], payload[7]),
	}, nil
}
