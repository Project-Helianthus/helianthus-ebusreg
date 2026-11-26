package owner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

func (p *vaillantSemanticPoller) readB509Hoursum2Int(ctx context.Context, target byte, addr uint16) *int {
	value := p.readB509Hoursum2(ctx, target, addr)
	if value == nil {
		return nil
	}
	v := int(*value)
	return &v
}

func (p *vaillantSemanticPoller) readB509PhoneBCD(ctx context.Context, target byte, addr uint16) *string {
	raw, ok := p.readB509Value(ctx, target, addr)
	if !ok || len(raw) == 0 {
		return nil
	}
	s := decodeBCDPhone(raw)
	if s == "" {
		return nil
	}
	return &s
}

func decodeBCDPhone(b []byte) string {
	var buf strings.Builder
	for _, v := range b {
		hi, lo := v>>4, v&0x0F
		if hi <= 9 {
			buf.WriteByte('0' + hi)
		} else {
			break
		}
		if lo <= 9 {
			buf.WriteByte('0' + lo)
		} else {
			break
		}
	}
	return buf.String()
}

func (p *vaillantSemanticPoller) readB509UCHInt(ctx context.Context, target byte, addr uint16) *int {
	raw, ok := p.readB509Value(ctx, target, addr)
	if !ok {
		return nil
	}
	value, ok := decodeUCH(raw)
	if !ok {
		return nil
	}
	out := int(value)
	return &out
}

func (p *vaillantSemanticPoller) readB509UCHFloat(ctx context.Context, target byte, addr uint16) *float64 {
	raw, ok := p.readB509Value(ctx, target, addr)
	if !ok {
		return nil
	}
	value, ok := decodeUCH(raw)
	if !ok {
		return nil
	}
	out := float64(value)
	return &out
}

func buildB509ReadSelector(addr uint16) []byte {
	return []byte{
		vaillantB509OpcodeRead,
		byte(addr >> 8),
		byte(addr),
	}
}

func buildB509WriteSelector(addr uint16, value []byte) []byte {
	data := make([]byte, 0, 3+len(value))
	data = append(data, vaillantB509OpcodeWrite, byte(addr>>8), byte(addr))
	data = append(data, value...)
	return data
}

func parseB509ReadPayload(payload []byte, addr uint16) ([]byte, bool) {
	if len(payload) == 0 {
		return nil, false
	}

	addrHi := byte(addr >> 8)
	addrLo := byte(addr)

	if len(payload) >= 3 && payload[0] == vaillantB509OpcodeRead && payload[1] == addrHi && payload[2] == addrLo {
		if len(payload) == 3 {
			return nil, false
		}
		return payload[3:], true
	}
	if len(payload) >= 2 && payload[0] == addrHi && payload[1] == addrLo {
		if len(payload) == 2 {
			return nil, false
		}
		return payload[2:], true
	}
	return payload, true
}

func parseB509WriteAck(payload []byte, addr uint16) bool {
	if len(payload) == 0 {
		return true
	}
	if len(payload) == 1 {
		return payload[0] != 0x00
	}

	addrHi := byte(addr >> 8)
	addrLo := byte(addr)

	if len(payload) >= 3 && payload[0] == vaillantB509OpcodeWrite && payload[1] == addrHi && payload[2] == addrLo {
		if len(payload) == 3 {
			return true
		}
		return payload[3] != 0x00
	}
	if len(payload) >= 2 && payload[0] == addrHi && payload[1] == addrLo {
		if len(payload) == 2 {
			return true
		}
		return payload[2] != 0x00
	}
	return payload[0] != 0x00
}

var (
	// Reserved for GW-2 wiring; keep implementations linked and lint-clean.
	_ = (*vaillantSemanticPoller).readB509Value
	_ = (*vaillantSemanticPoller).writeB509Value
)

// readB555Frame sends a B555 frame and returns the responder payload.
func (p *vaillantSemanticPoller) readB555Frame(ctx context.Context, data []byte) ([]byte, bool) {
	if p == nil || p.bus == nil {
		return nil, false
	}
	if ctx == nil {
		ctx = context.Background()
	}

	p.mu.Lock()
	source := p.source
	target := p.controller
	timeout := p.requestTimeout
	p.mu.Unlock()

	if target == 0 {
		return nil, false
	}
	if timeout <= 0 {
		timeout = 2 * time.Second
	}

	p.readMu.Lock()
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	request := protocol.Frame{
		FrameType: protocol.FrameTypeInitiatorTarget,
		Source:    source,
		Target:    target,
		Primary:   vaillantB555Primary,
		Secondary: vaillantB555Secondary,
		Data:      data,
	}
	response, err := p.bus.Send(reqCtx, request)
	cancel()
	p.readMu.Unlock()

	if err != nil || response == nil {
		return nil, false
	}
	return response.Data, true
}

func (p *vaillantSemanticPoller) writeB555Frame(ctx context.Context, data []byte) ([]byte, bool) {
	if p == nil || p.bus == nil {
		return nil, false
	}
	if ctx == nil {
		ctx = context.Background()
	}

	p.mu.Lock()
	source := p.source
	target := p.controller
	timeout := p.requestTimeout
	p.mu.Unlock()

	if target == 0 {
		return nil, false
	}
	if timeout <= 0 {
		timeout = 2 * time.Second
	}

	p.readMu.Lock()
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	request := protocol.Frame{
		FrameType: protocol.FrameTypeInitiatorTarget,
		Source:    source,
		Target:    target,
		Primary:   vaillantB555Primary,
		Secondary: vaillantB555Secondary,
		Data:      data,
	}
	response, err := p.bus.Send(reqCtx, request)
	cancel()
	p.readMu.Unlock()

	if err != nil || response == nil {
		return nil, false
	}
	return response.Data, true
}

func (p *vaillantSemanticPoller) writeB555TimerSlot(ctx context.Context, zone, hc, weekday, slotIdx, slotCount, startHour, startMinute, endHour, endMinute byte, temperatureRaw uint16) (byte, error) {
	data := []byte{
		b555OpcodeTimerWrite,
		zone, hc, weekday,
		slotIdx, slotCount,
		startHour, startMinute,
		endHour, endMinute,
		byte(temperatureRaw & 0xFF), byte(temperatureRaw >> 8),
	}
	resp, ok := p.writeB555Frame(ctx, data)
	if !ok {
		return 0xFF, fmt.Errorf("b555 timer write: no response")
	}
	if len(resp) < 1 {
		return 0xFF, fmt.Errorf("b555 timer write: empty response")
	}
	return resp[0], nil
}
