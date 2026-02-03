package registry

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
	"github.com/d3vi1/helianthus-ebusgo/protocol"
)

const (
	scanPrimary   = byte(0x07)
	scanSecondary = byte(0x04)
)

var errScanResponsePayload = errors.New("scan: invalid response payload")

type ScanBus interface {
	Send(ctx context.Context, frame protocol.Frame) (*protocol.Frame, error)
}

const scanCollisionMaxPasses = 3

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
	found := make(map[byte]struct{})
	pending := dedupeScanTargets(targets)
	for pass := 0; pass <= scanCollisionMaxPasses && len(pending) > 0; pass++ {
		collisions := make([]byte, 0)

		for _, target := range pending {
			if err := ctx.Err(); err != nil {
				return nil, err
			}

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
				if ctx.Err() != nil {
					return nil, ctx.Err()
				}
				// Arbitration collisions are transient and do not imply a missing device.
				// ebusd retries later; mimic that by deferring the address to a later pass.
				if errors.Is(err, ebuserrors.ErrBusCollision) {
					collisions = append(collisions, target)
					continue
				}
				if shouldSkipScanError(err) {
					continue
				}
				return nil, fmt.Errorf("scan target %02x: %w", target, err)
			}
			if response == nil {
				err := fmt.Errorf("scan target %02x empty response: %w", target, errors.Join(errScanResponsePayload, ebuserrors.ErrInvalidPayload))
				if shouldSkipScanError(err) {
					continue
				}
				return nil, err
			}

			address := response.Source
			if address == 0 {
				address = target
			}
			if _, ok := found[address]; ok {
				continue
			}

			info, err := parseDeviceInfo(address, response.Data)
			if err != nil {
				if shouldSkipScanError(err) {
					continue
				}
				return nil, fmt.Errorf("scan target %02x parse: %w", target, err)
			}
			entries = append(entries, registry.Register(info))
			found[address] = struct{}{}
		}

		if len(collisions) == 0 {
			break
		}
		if pass == scanCollisionMaxPasses {
			break
		}
		pending = dedupeScanTargets(collisions)

		timer := time.NewTimer(25 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
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
	if len(payload) < 10 {
		return DeviceInfo{}, fmt.Errorf("short device info payload: %w", errors.Join(errScanResponsePayload, ebuserrors.ErrInvalidPayload))
	}

	manufacturer := fmt.Sprintf("0x%02X", payload[0])
	if payload[0] == 0xB5 {
		manufacturer = "Vaillant"
	}

	return DeviceInfo{
		Address:         address,
		Manufacturer:    manufacturer,
		DeviceID:        strings.Trim(string(payload[1:6]), " \x00"),
		SoftwareVersion: fmt.Sprintf("%02X%02X", payload[6], payload[7]),
		HardwareVersion: fmt.Sprintf("%02X%02X", payload[8], payload[9]),
	}, nil
}

func shouldSkipScanError(err error) bool {
	return errors.Is(err, ebuserrors.ErrNoSuchDevice) ||
		errors.Is(err, ebuserrors.ErrTimeout) ||
		errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, ebuserrors.ErrNACK) ||
		errors.Is(err, ebuserrors.ErrCRCMismatch) ||
		errors.Is(err, errScanResponsePayload)
}

func dedupeScanTargets(targets []byte) []byte {
	if len(targets) == 0 {
		return nil
	}
	seen := make(map[byte]struct{}, len(targets))
	out := make([]byte, 0, len(targets))
	for _, target := range targets {
		if _, ok := seen[target]; ok {
			continue
		}
		seen[target] = struct{}{}
		out = append(out, target)
	}
	return out
}
