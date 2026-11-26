package owner

import (
	"context"
	"fmt"
	"log"
	"time"

	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

func buildB524ReadSelector(opcode, group, instance byte, addr uint16) []byte {
	return []byte{
		opcode,
		vaillantB524OpRead,
		group,
		instance,
		byte(addr),
		byte(addr >> 8),
	}
}

func buildB524WriteSelector(opcode, group, instance byte, addr uint16, data []byte) []byte {
	selector := []byte{
		opcode,
		vaillantB524OpWrite,
		group,
		instance,
		byte(addr),
		byte(addr >> 8),
	}
	return append(selector, data...)
}

func parseB524ReadPayload(payload []byte, opcode, group, instance byte, addr uint16) ([]byte, bool) {
	if len(payload) == 0 {
		return nil, false
	}

	if len(payload) == 1 && payload[0] == 0x00 {
		return nil, false
	}
	if len(payload) < 4 {
		return nil, false
	}

	replyKind := payload[0]

	if len(payload) >= 5 {
		replyInstance := payload[1]
		replyGroup := payload[2]
		replyAddr := uint16(payload[3]) | uint16(payload[4])<<8
		if replyGroup == group && replyAddr == addr {
			if !matchesB524ReplyInstance(replyInstance, instance) {
				log.Printf(
					"b524 read mismatch: want opcode=0x%02x group=0x%02x instance=0x%02x addr=0x%04x; got reply-instance=0x%02x (group=0x%02x addr=0x%04x)",
					opcode,
					group,
					instance,
					addr,
					replyInstance,
					replyGroup,
					replyAddr,
				)
				return nil, false
			}
			if len(payload) == 5 {
				return nil, false
			}
			return payload[5:], true
		}
	}

	replyGroup := payload[1]
	replyAddr := uint16(payload[2]) | uint16(payload[3])<<8
	if replyGroup != group || replyAddr != addr {
		log.Printf(
			"b524 read mismatch: want opcode=0x%02x group=0x%02x instance=0x%02x addr=0x%04x; got kind=0x%02x group=0x%02x addr=0x%04x len=%d",
			opcode,
			group,
			instance,
			addr,
			replyKind,
			replyGroup,
			replyAddr,
			len(payload),
		)
		return nil, false
	}
	if len(payload) == 4 {
		return nil, false
	}
	return payload[4:], true
}

// isB524ProbeCoherent checks if a B524 read response proves the target
// understands B524 register reads. Unlike parseB524ReadPayload, this accepts
// header-only responses (no value bytes) as coherent — the device returned a
// valid B524 structure which proves capability.
func isB524ProbeCoherent(payload []byte, group byte, addr uint16) bool {
	return isB524ResponseCoherent(payload, group, addr)
}

// probeB524Register sends a single B524 read request to target and checks
// if the response is coherent. Retries up to 3 times with 200ms backoff
// to handle adapter-direct bus contention.
func (p *vaillantSemanticPoller) probeB524Register(ctx context.Context, target, opcode, group, instance byte, addr uint16) bool {
	if p == nil || p.bus == nil {
		return false
	}

	p.mu.Lock()
	source := p.source
	timeout := p.requestTimeout
	p.mu.Unlock()

	// B524 probes need at least 5s per attempt in adapter-direct mode:
	// arbitration(~200ms) + frame(~220ms) + adapter-timeout(~550ms) +
	// collision-retry(~200ms) + second-attempt(~970ms) ≈ 2.1s minimum.
	// The config default SemanticRequestTimeout=2s is insufficient.
	const minB524ProbeTimeout = 5 * time.Second
	if timeout < minB524ProbeTimeout {
		timeout = minB524ProbeTimeout
	}

	// Retry probes up to 3 times — adapter-direct mode has bus contention
	// that can cause individual probes to timeout on busy buses.
	for attempt := 0; attempt < 3; attempt++ {
		// Back off between retries to avoid colliding with the same
		// bus contention pattern (50ms arbitration + waitForSyn).
		if attempt > 0 && !p.pause(ctx, 200*time.Millisecond) {
			return false
		}

		data := buildB524ReadSelector(opcode, group, instance, addr)
		frame := protocol.Frame{
			FrameType: protocol.FrameTypeInitiatorTarget,
			Source:    source,
			Target:    target,
			Primary:   vaillantExtRegisterPrimary,
			Secondary: vaillantExtRegisterSecondary,
			Data:      data,
		}

		p.readMu.Lock()
		reqCtx, cancel := context.WithTimeout(ctx, timeout)
		response, err := p.sendSemanticFrame(reqCtx, frame)
		cancel()
		p.readMu.Unlock()

		if err != nil {
			if ctx.Err() != nil {
				return false // parent context cancelled
			}
			// Don't retry definitive rejections (NACK, no device) —
			// these are permanent, not transient bus contention.
			if ebuserrors.IsDefinitive(err) {
				return false
			}
			continue // retry on transient errors
		}
		if response == nil {
			continue
		}
		return isB524ProbeCoherent(response.Data, group, addr)
	}
	return false
}

func matchesB524ReplyInstance(replyInstance, requestedInstance byte) bool {
	if replyInstance == requestedInstance {
		return true
	}
	if requestedInstance < 0xFF && replyInstance == requestedInstance+1 {
		return true
	}
	return false
}

func (p *vaillantSemanticPoller) writeB524Value(ctx context.Context, opcode, group, instance byte, addr uint16, data []byte) error {
	if p == nil || (p.bus == nil && p.sendFrameFn == nil) {
		return fmt.Errorf("b524 write unavailable")
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
		return fmt.Errorf("b524 write target is zero")
	}
	if timeout <= 0 {
		timeout = 2 * time.Second
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		p.readMu.Lock()
		reqCtx, cancel := context.WithTimeout(ctx, timeout)
		request := protocol.Frame{
			FrameType: protocol.FrameTypeInitiatorTarget,
			Source:    source,
			Target:    target,
			Primary:   vaillantExtRegisterPrimary,
			Secondary: vaillantExtRegisterSecondary,
			Data:      buildB524WriteSelector(opcode, group, instance, addr, data),
		}
		response, err := p.sendSemanticFrame(reqCtx, request)
		cancel()
		p.readMu.Unlock()

		if err != nil {
			lastErr = err
		} else if response != nil {
			// B524 write responses are typically short ACKs (< 4 bytes payload).
			return nil
		} else {
			lastErr = fmt.Errorf("b524 write returned nil response")
		}

		if attempt < 2 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(75 * time.Millisecond):
			}
		}
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("b524 write failed")
	}
	return lastErr
}
