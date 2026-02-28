package registry

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

const (
	scanPrimary   = byte(0x07)
	scanSecondary = byte(0x04)
)

const (
	vaillantPrimary         = byte(0xB5)
	vaillantScanIDSecondary = byte(0x09)
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
	registered := make(map[byte]struct{})
	pending := dedupeScanTargets(targets)
	for pass := 0; pass <= scanCollisionMaxPasses && len(pending) > 0; pass++ {
		collisions := make([]byte, 0)
		retries := make([]byte, 0)

		for _, target := range pending {
			if err := ctx.Err(); err != nil {
				return nil, err
			}

			frameType := protocol.FrameTypeForTarget(target)
			if frameType == protocol.FrameTypeUnknown || isInitiatorCapableAddress(target) {
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
				if shouldRetryScanError(err) {
					retries = append(retries, target)
					continue
				}
				if shouldSkipScanError(err) {
					continue
				}
				return nil, fmt.Errorf("scan target %02x: %w", target, err)
			}
			if response == nil {
				err := fmt.Errorf("scan target %02x empty response: %w", target, errors.Join(errScanResponsePayload, ebuserrors.ErrInvalidPayload))
				if shouldRetryScanError(err) {
					retries = append(retries, target)
					continue
				}
				if shouldSkipScanError(err) {
					continue
				}
				return nil, err
			}

			address := response.Source
			if address == 0 {
				address = target
			}

			info, err := parseDeviceInfo(address, response.Data)
			if err != nil {
				if shouldRetryScanError(err) {
					retries = append(retries, target)
					continue
				}
				if shouldSkipScanError(err) {
					continue
				}
				return nil, fmt.Errorf("scan target %02x parse: %w", target, err)
			}

			if info.Manufacturer == "Vaillant" && info.SerialNumber == "" {
				serial, ok, err := readVaillantScanID(ctx, bus, source, address)
				if err != nil {
					return nil, err
				}
				if ok {
					info.SerialNumber = serial
				}
			}
			if info.SerialNumber == "" && info.Manufacturer == "Vaillant" {
				if existing, ok := registry.lookupByIdentity(info); ok && shouldReuseSerial(info, existing) {
					info.SerialNumber = existing.SerialNumber()
				}
			}
			entry := registry.Register(info)
			if address != target {
				alias := info
				alias.Address = target
				entry = registry.Register(alias)
			}
			canonicalAddress := entry.Address()
			if _, ok := registered[canonicalAddress]; !ok {
				registered[canonicalAddress] = struct{}{}
				entries = append(entries, entry)
			}
		}

		if len(collisions) == 0 && len(retries) == 0 {
			break
		}
		if pass == scanCollisionMaxPasses {
			break
		}
		pending = dedupeScanTargets(append(collisions, retries...))

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

// ReadVaillantScanID sends B5.09 identity reads (chunks 0x24..0x27) and
// returns the formatted Vaillant serial number.  It is the exported twin of
// the internal readVaillantScanID so that callers outside this package (e.g.
// the gateway's ebusd-tcp preload enrichment) can invoke it.
func ReadVaillantScanID(ctx context.Context, bus ScanBus, source byte, target byte) (string, bool, error) {
	return readVaillantScanID(ctx, bus, source, target)
}

func readVaillantScanID(ctx context.Context, bus ScanBus, source byte, target byte) (string, bool, error) {
	if bus == nil {
		return "", false, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	raw := make([]byte, 0, 32)
	for qq := byte(0x24); qq <= byte(0x27); qq++ {
		request := protocol.Frame{
			Source:    source,
			Target:    target,
			Primary:   vaillantPrimary,
			Secondary: vaillantScanIDSecondary,
			Data:      []byte{qq},
		}
		response, err := bus.Send(ctx, request)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return "", false, err
			}
			if errors.Is(err, context.DeadlineExceeded) {
				if ctxErr := ctx.Err(); ctxErr != nil {
					return "", false, ctxErr
				}
				return "", false, nil
			}
			return "", false, nil
		}
		if response == nil {
			return "", false, nil
		}
		if len(response.Data) != 9 || response.Data[0] != 0x00 {
			return "", false, nil
		}
		raw = append(raw, response.Data[1:]...)
	}

	trimmed := trimScanIDBytes(raw)
	if len(trimmed) == 0 {
		return "", false, nil
	}

	formatted := formatVaillantSerial(string(trimmed))
	if formatted == "" {
		return "", false, nil
	}
	return formatted, true, nil
}

func trimScanIDBytes(data []byte) []byte {
	end := len(data)
	for end > 0 {
		last := data[end-1]
		if last == 0x00 || last == 0x20 || last == 0xFF {
			end--
			continue
		}
		break
	}
	return data[:end]
}

func formatVaillantSerial(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if len(raw) < 28 {
		return raw
	}
	candidate := raw[:28]
	for i := 0; i < 26; i++ {
		if candidate[i] < '0' || candidate[i] > '9' {
			return raw
		}
	}
	return fmt.Sprintf(
		"%s-%s-%s-%s-%s-%s-%s",
		candidate[0:2],
		candidate[2:4],
		candidate[4:6],
		candidate[6:16],
		candidate[16:20],
		candidate[20:26],
		candidate[26:28],
	)
}

func shouldReuseSerial(info DeviceInfo, existing DeviceEntry) bool {
	if existing == nil || existing.SerialNumber() == "" {
		return false
	}
	if !strings.EqualFold(existing.Manufacturer(), info.Manufacturer) {
		return false
	}
	if info.DeviceID != "" && existing.DeviceID() != info.DeviceID {
		return false
	}
	if info.SoftwareVersion != "" && existing.SoftwareVersion() != info.SoftwareVersion {
		return false
	}
	if info.HardwareVersion != "" && existing.HardwareVersion() != info.HardwareVersion {
		return false
	}
	return true
}

func shouldSkipScanError(err error) bool {
	return errors.Is(err, ebuserrors.ErrNoSuchDevice) ||
		errors.Is(err, ebuserrors.ErrTimeout) ||
		errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, ebuserrors.ErrNACK) ||
		errors.Is(err, ebuserrors.ErrCRCMismatch) ||
		errors.Is(err, errScanResponsePayload)
}

func shouldRetryScanError(err error) bool {
	return errors.Is(err, ebuserrors.ErrTimeout) ||
		errors.Is(err, context.DeadlineExceeded) ||
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

func isInitiatorCapableAddress(address byte) bool {
	return initiatorPartIndex(address&0x0F) > 0 && initiatorPartIndex((address&0xF0)>>4) > 0
}

func initiatorPartIndex(bits byte) byte {
	switch bits {
	case 0x0:
		return 1
	case 0x1:
		return 2
	case 0x3:
		return 3
	case 0x7:
		return 4
	case 0xF:
		return 5
	default:
		return 0
	}
}
