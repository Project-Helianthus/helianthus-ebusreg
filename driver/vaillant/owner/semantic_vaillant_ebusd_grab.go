package owner

import (
	"bufio"
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
)

const (
	ebusdGrabCommand               = "grab result all\n"
	ebusdGrabFollowupReadWindow    = 2 * time.Second
	ebusdGrabMinimumInitialTimeout = 2 * time.Second
)

func (p *vaillantSemanticPoller) refreshFromEbusdGrabWithInstances(ctx context.Context) (map[byte]bool, bool) {
	if p == nil {
		return nil, false
	}
	if !isEbusdTCPTransport(p.transportConfig) {
		return nil, false
	}

	p.mu.Lock()
	controller := p.controller
	p.mu.Unlock()

	zones, ok := readB524ZonesFromEbusdGrab(ctx, p.transportConfig, controller)
	if !ok || len(zones) == 0 {
		return nil, false
	}
	hydrated := make(map[byte]bool, len(zones))
	for instance := range zones {
		hydrated[instance] = true
	}

	p.mu.Lock()
	if p.controller == 0 && controller != 0 {
		p.controller = controller
	}
	for instance, incoming := range zones {
		p.markZonePresentLocked(instance)
		entry := p.zones[instance]
		if entry == nil {
			continue
		}
		mergeZoneSnapshotFields(entry, incoming, semanticSnapshotSourceLive, zoneGrabFieldSet)
	}
	p.mu.Unlock()
	return hydrated, true
}

func isEbusdTCPTransport(cfg engine.TransportConfig) bool {
	if strings.TrimSpace(strings.ToLower(cfg.Network)) != "tcp" {
		return false
	}
	return strings.TrimSpace(strings.ToLower(string(cfg.Protocol))) == "ebusd-tcp"
}

func readB524ZonesFromEbusdGrab(ctx context.Context, cfg engine.TransportConfig, controller byte) (map[byte]*vaillantZoneSnapshot, bool) {
	candidates := ebusdScanTargetCandidates(cfg)
	for _, candidate := range candidates {
		lines, err := ebusdCommandLines(ctx, candidate, ebusdGrabCommand)
		if err != nil || len(lines) == 0 {
			continue
		}
		zones := parseB524ZonesFromGrab(lines, controller)
		if len(zones) > 0 {
			return zones, true
		}
	}
	return nil, false
}

func (p *vaillantSemanticPoller) refreshDHWFromEbusdGrab(ctx context.Context) bool {
	if p == nil {
		return false
	}
	if !isEbusdTCPTransport(p.transportConfig) {
		return false
	}
	p.mu.Lock()
	controller := p.controller
	p.mu.Unlock()

	candidates := ebusdScanTargetCandidates(p.transportConfig)
	for _, candidate := range candidates {
		lines, err := ebusdCommandLines(ctx, candidate, ebusdGrabCommand)
		if err != nil || len(lines) == 0 {
			continue
		}
		dhw, ok := parseB524DhwFromGrab(lines, controller)
		if !ok || dhw == nil {
			continue
		}
		p.mu.Lock()
		if p.dhw == nil {
			p.dhw = &vaillantDhwSnapshot{}
		}
		mergeDhwSnapshotFields(p.dhw, dhw, semanticSnapshotSourceLive, dhwFieldSet)
		p.markDHWUpdatedNowLocked()
		p.mu.Unlock()
		return true
	}
	return false
}

func ebusdCommandLines(ctx context.Context, cfg engine.TransportConfig, command string) ([]string, error) {
	command = strings.TrimSpace(command)
	if command == "" || cfg.Address == "" || !isEbusdTCPTransport(cfg) {
		return nil, nil
	}

	dial := cfg.Dial
	if dial == nil {
		dial = dialContext
	}

	dialCtx := ctx
	cancel := func() {}
	if dialCtx == nil {
		dialCtx = context.Background()
	}
	if cfg.DialTimeout > 0 {
		dialCtx, cancel = context.WithTimeout(dialCtx, cfg.DialTimeout)
	}
	defer cancel()

	conn, err := dial(dialCtx, cfg.Network, cfg.Address, cfg.DialTimeout)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()

	reader := bufio.NewReader(conn)
	if _, err := io.WriteString(conn, command+"\n"); err != nil {
		return nil, err
	}

	initialTimeout := cfg.ReadTimeout
	if initialTimeout < ebusdGrabMinimumInitialTimeout {
		initialTimeout = ebusdGrabMinimumInitialTimeout
	}
	if initialTimeout > 0 {
		_ = conn.SetReadDeadline(time.Now().Add(initialTimeout))
	}

	lines := make([]string, 0, 256)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if isNetTimeout(err) || errors.Is(err, io.EOF) {
				break
			}
			return lines, err
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		lines = append(lines, trimmed)
		if initialTimeout > 0 {
			_ = conn.SetReadDeadline(time.Now().Add(ebusdGrabFollowupReadWindow))
		}
	}

	return lines, nil
}

// b524GrabZoneNameParts holds the three name registers of one zone seen in a
// grab dump; they are composed after every line has been applied.
type b524GrabZoneNameParts struct {
	primary string
	prefix  string
	suffix  string
}

// b524GrabZoneState accumulates the registers of one zone instance from a
// grab dump before the zone is derived and admitted.
type b524GrabZoneState struct {
	snapshot                      *vaillantZoneSnapshot
	nameParts                     b524GrabZoneNameParts
	hasIndex                      bool
	indexPresent                  bool
	configurationOperationMode    *uint16
	stateSpecialFunction          *uint16
	stateValveStatus              *uint16
	roomTemperatureZoneMappingRaw *uint16
}

func parseB524ZonesFromGrab(lines []string, controller byte) map[byte]*vaillantZoneSnapshot {
	states := make(map[byte]*b524GrabZoneState)
	zones := make(map[byte]*vaillantZoneSnapshot)
	circuitTypeByInstance := make(map[byte]*uint16)

	for _, line := range lines {
		group, instance, addr, payload, ok := parseB524GrabLine(line, controller)
		if !ok {
			continue
		}
		if group == localCircuits.group && addr == circuitTypeRegister {
			if value, ok := decodeB524Uint16(payload); ok {
				typed := value
				circuitTypeByInstance[instance] = &typed
			}
			continue
		}
		if group != localZones.group {
			continue
		}
		state := states[instance]
		if state == nil {
			state = &b524GrabZoneState{
				snapshot: &vaillantZoneSnapshot{
					Instance: instance,
				},
			}
			states[instance] = state
		}
		applyB524GrabZoneRegister(state, addr, payload)
	}

	for instance, state := range states {
		if snapshot, ok := finishB524GrabZone(instance, state, circuitTypeByInstance); ok {
			zones[instance] = snapshot
		}
	}

	return zones
}

// applyB524GrabZoneRegister applies one zone register payload of a grab dump
// to the accumulated zone state.
func applyB524GrabZoneRegister(state *b524GrabZoneState, addr uint16, payload []byte) {
	snapshot := state.snapshot
	switch addr {
	case zoneIndex:
		state.hasIndex = true
		if len(payload) > 0 && payload[0] != 0xFF {
			state.indexPresent = true
		}
	case zoneName:
		parts := state.nameParts
		parts.primary = decodeCString(payload)
		state.nameParts = parts
	case zoneNamePrefix:
		parts := state.nameParts
		parts.prefix = decodeCString(payload)
		state.nameParts = parts
	case zoneNameSuffix:
		parts := state.nameParts
		parts.suffix = decodeCString(payload)
		state.nameParts = parts
	case zoneHeatingOpMode:
		if value, ok := decodeB524Uint16(payload); ok {
			typed := value
			state.configurationOperationMode = &typed
			snapshot.ConfigurationHeatingOperationMode = formatUintToken(value)
		}
	case zoneSpecialFunction:
		if value, ok := decodeB524Uint16(payload); ok {
			typed := value
			state.stateSpecialFunction = &typed
			snapshot.StateSpecialFunction = formatUintToken(value)
		}
	case zoneValveStatus:
		if value, ok := decodeB524Uint16(payload); ok {
			typed := value
			state.stateValveStatus = &typed
			snapshot.StateValveStatusRaw = &typed
		}
	case zoneRoomTemperatureZoneMappingRawRegister:
		if value, ok := decodeB524Uint16(payload); ok {
			typed := value
			state.roomTemperatureZoneMappingRaw = &typed
			snapshot.ConfigurationRoomTemperatureZoneMappingRaw = &typed
		}
	default:
		applyB524GrabZoneTemperatureRegister(snapshot, addr, payload)
	}
}

// applyB524GrabZoneTemperatureRegister applies the temperature and humidity
// registers of one zone from a grab dump.
func applyB524GrabZoneTemperatureRegister(snapshot *vaillantZoneSnapshot, addr uint16, payload []byte) {
	switch addr {
	case zoneCurrentTemp:
		if value, ok := decodeFloat32LE(payload); ok {
			v := value
			snapshot.CurrentTempC = &v
		}
	case zoneTargetTemp:
		if value, ok := decodeFloat32LE(payload); ok {
			v := value
			snapshot.TargetTempC = &v
		}
	case zoneFallbackManualTemp:
		if snapshot.TargetTempC == nil {
			if value, ok := decodeFloat32LE(payload); ok {
				v := value
				snapshot.TargetTempC = &v
			}
		}
	case zoneCurrentHumidity:
		if value, ok := decodeFloat32LE(payload); ok {
			v := value
			snapshot.HumidityPct = &v
		}
	}
}
