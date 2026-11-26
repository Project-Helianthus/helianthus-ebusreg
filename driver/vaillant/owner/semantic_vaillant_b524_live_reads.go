package owner

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

func (p *vaillantSemanticPoller) readB524ValueLive(ctx context.Context, opcode, group, instance byte, addr uint16) ([]byte, bool) {
	if p == nil || (p.bus == nil && p.sendFrameFn == nil) {
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

	for attempt := 0; attempt < 3; attempt++ {
		p.readMu.Lock()
		reqCtx, cancel := context.WithTimeout(ctx, timeout)
		request := protocol.Frame{
			FrameType: protocol.FrameTypeInitiatorTarget,
			Source:    source,
			Target:    target,
			Primary:   vaillantExtRegisterPrimary,
			Secondary: vaillantExtRegisterSecondary,
			Data:      buildB524ReadSelector(opcode, group, instance, addr),
		}
		response, err := p.sendSemanticFrame(reqCtx, request)
		cancel()
		p.readMu.Unlock()

		if err == nil && response != nil {
			payload, ok := parseB524ReadPayload(response.Data, opcode, group, instance, addr)
			if ok {
				return payload, true
			}
		}

		if attempt < 2 {
			select {
			case <-ctx.Done():
				return nil, false
			case <-time.After(75 * time.Millisecond):
			}
		}
	}
	return nil, false
}

// Reader kinds of the fixed-width B524 readers, used as keys of
// semantic_b524_width_mismatch_total.
const (
	b524ReaderKindUint8         = "u8"
	b524ReaderKindUint16        = "uint16"
	b524ReaderKindFloat32       = "float32le"
	b524ReaderKindUint32        = "uint32le"
	b524ReaderKindDateThreeByte = "date_hda3"

	// b524WidthMismatchLogInterval bounds the width mismatch log to one
	// line per selector in this interval.
	b524WidthMismatchLogInterval = 10 * time.Minute
)

// b524WidthMismatchSelector is the selector the width mismatch line names.
type b524WidthMismatchSelector struct {
	opcode byte
	group  byte
	addr   uint16
}

// b524ReplyWidthMatches reports whether a reply of got bytes fits a reader
// that accepts exactly one of the widths in accepted. A mismatch is counted
// under kind and logged at most once per selector per log interval; the
// caller then returns no value. want is the first accepted width.
func (p *vaillantSemanticPoller) b524ReplyWidthMatches(kind string, opcode, group byte, addr uint16, got int, accepted ...int) bool {
	for _, width := range accepted {
		if got == width {
			return true
		}
	}
	semanticB524WidthMismatchTotal.Add(kind, 1)
	if p == nil {
		return false
	}
	selector := b524WidthMismatchSelector{opcode: opcode, group: group, addr: addr}
	now := p.now()
	p.b524WidthMismatchMu.Lock()
	last, seen := p.b524WidthMismatchLoggedAt[selector]
	logLine := !seen || now.Sub(last) >= b524WidthMismatchLogInterval
	if logLine {
		if p.b524WidthMismatchLoggedAt == nil {
			p.b524WidthMismatchLoggedAt = make(map[b524WidthMismatchSelector]time.Time)
		}
		p.b524WidthMismatchLoggedAt[selector] = now
	}
	p.b524WidthMismatchMu.Unlock()
	if logLine {
		want := 0
		if len(accepted) > 0 {
			want = accepted[0]
		}
		log.Printf("semantic: b524 reply width mismatch opcode=0x%02X group=0x%02X register=0x%04X want=%d got=%d", opcode, group, addr, want, got)
	}
	return false
}

func (p *vaillantSemanticPoller) readB524Float32LE(ctx context.Context, opcode, group, instance byte, addr uint16) (float64, bool) {
	// The circuit RR0020 selector is raw evidence in the native catalog.
	// Reject it before acquisition or cache lookup, regardless of payload
	// width; this rejection is not a width mismatch.
	if opcode == localCircuits.opcode && group == localCircuits.group && addr == circuit_calc_flow_temp {
		return 0, false
	}
	raw, ok := p.readB524Value(ctx, opcode, group, instance, addr)
	if !ok || !p.b524ReplyWidthMatches(b524ReaderKindFloat32, opcode, group, addr, len(raw), 4) {
		return 0, false
	}
	bits := binary.LittleEndian.Uint32(raw[:4])
	value := float64(math.Float32frombits(bits))
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, false
	}
	return value, true
}

func (p *vaillantSemanticPoller) readB524Uint8(ctx context.Context, opcode, group, instance byte, addr uint16) *uint8 {
	raw, ok := p.readB524Value(ctx, opcode, group, instance, addr)
	if !ok || !p.b524ReplyWidthMatches(b524ReaderKindUint8, opcode, group, addr, len(raw), 1) {
		return nil
	}
	if raw[0] == 0xFF {
		return nil
	}
	v := raw[0]
	return &v
}

func (p *vaillantSemanticPoller) readB524Uint16Pointer(ctx context.Context, opcode, group, instance byte, addr uint16) *uint16 {
	v, ok := p.readB524Uint16(ctx, opcode, group, instance, addr)
	if !ok || v == nil {
		return nil
	}
	if *v == 0xFFFF {
		return nil
	}
	return v
}

func (p *vaillantSemanticPoller) readB524Float32(ctx context.Context, opcode, group, instance byte, addr uint16) *float64 {
	v, ok := p.readB524Float32LE(ctx, opcode, group, instance, addr)
	if !ok {
		return nil
	}
	return &v
}

func (p *vaillantSemanticPoller) readB524Bool(ctx context.Context, opcode, group, instance byte, addr uint16) *bool {
	v := p.readB524Uint8(ctx, opcode, group, instance, addr)
	if v == nil {
		return nil
	}
	if *v > 1 {
		return nil
	}
	value := *v == 1
	return &value
}

func (p *vaillantSemanticPoller) readB524Firmware(ctx context.Context, opcode, group, instance byte, addr uint16) *string {
	raw, ok := p.readB524Value(ctx, opcode, group, instance, addr)
	if !ok {
		return nil
	}
	return decodeB524FirmwareVersion(raw)
}

func (p *vaillantSemanticPoller) readB524NullTerminatedString(ctx context.Context, opcode, group, instance byte, addr uint16) (string, bool) {
	raw, ok := p.readB524Value(ctx, opcode, group, instance, addr)
	if !ok || len(raw) == 0 {
		return "", false
	}
	trimmed := raw
	for i, b := range trimmed {
		if b == 0x00 {
			trimmed = trimmed[:i]
			break
		}
	}
	return string(trimmed), true
}

func (p *vaillantSemanticPoller) readB524NullTerminatedStringSanitized(ctx context.Context, opcode, group, instance byte, addr uint16) (string, bool) {
	value, ok := p.readB524NullTerminatedString(ctx, opcode, group, instance, addr)
	if !ok {
		return "", false
	}
	sanitized := make([]byte, len(value))
	for i := 0; i < len(value); i++ {
		if value[i] > 0x7F {
			sanitized[i] = '?'
		} else {
			sanitized[i] = value[i]
		}
	}
	return string(sanitized), true
}

func (p *vaillantSemanticPoller) readB524DateThreeByte(ctx context.Context, opcode, group, instance byte, addr uint16) (string, bool) {
	raw, ok := p.readB524Value(ctx, opcode, group, instance, addr)
	if !ok || !p.b524ReplyWidthMatches(b524ReaderKindDateThreeByte, opcode, group, addr, len(raw), 3) {
		return "", false
	}
	day, month, yearOffset := int(raw[0]), int(raw[1]), int(raw[2])
	if day < 1 || day > 31 || month < 1 || month > 12 || yearOffset > 99 {
		return "", false
	}
	year := 2000 + yearOffset
	iso := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	return iso, true
}

func (p *vaillantSemanticPoller) readB524ZoneNamePart(ctx context.Context, instance byte, addr uint16) (string, bool) {
	raw, ok := p.readB524NullTerminatedString(ctx, localZones.opcode, localZones.group, instance, addr)
	if !ok {
		return "", false
	}
	return strings.TrimSpace(raw), true
}

func (p *vaillantSemanticPoller) readB524Uint16(ctx context.Context, opcode, group, instance byte, addr uint16) (*uint16, bool) {
	raw, ok := p.readB524Value(ctx, opcode, group, instance, addr)
	// The codec takes the two-byte form and the one-byte form.
	if !ok || !p.b524ReplyWidthMatches(b524ReaderKindUint16, opcode, group, addr, len(raw), 2, 1) {
		return nil, false
	}
	parsed, ok := decodeB524Uint16(raw)
	if !ok {
		return nil, false
	}
	return &parsed, true
}

func (p *vaillantSemanticPoller) readB524Uint32LE(ctx context.Context, opcode, group, instance byte, addr uint16) (uint32, bool) {
	raw, ok := p.readB524Value(ctx, opcode, group, instance, addr)
	if !ok || !p.b524ReplyWidthMatches(b524ReaderKindUint32, opcode, group, addr, len(raw), 4) {
		return 0, false
	}
	v := binary.LittleEndian.Uint32(raw[:4])
	if v == 0xFFFFFFFF {
		return 0, false
	}
	return v, true
}
