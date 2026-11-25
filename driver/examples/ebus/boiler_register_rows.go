package ebus

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

type boilerRegisterRead struct {
	Target   byte
	Register uint16
}

type boilerRegisterWrite struct {
	Target   byte
	Register uint16
	Encoded  []byte
}

func (s boilerRegisterRead) request(source byte) protocol.Frame {
	return protocol.Frame{Source: source, Target: s.Target, Primary: 0xb5, Secondary: 0x09, Data: binary.BigEndian.AppendUint16([]byte{0x0d}, s.Register)}
}

func (s boilerRegisterWrite) request(source byte) protocol.Frame {
	data := binary.BigEndian.AppendUint16([]byte{0x0e}, s.Register)
	return protocol.Frame{Source: source, Target: s.Target, Primary: 0xb5, Secondary: 0x09, Data: append(data, s.Encoded...)}
}

type boilerRegisterValue struct {
	Raw, Tail []byte
	PumpOn    *bool
}

func (s boilerRegisterRead) pumpReply(raw []byte) boilerRegisterValue {
	r := boilerRegisterValue{Raw: raw}
	hi, lo := byte(s.Register>>8), byte(s.Register)
	tail := raw
	switch {
	case len(raw) >= 3 && raw[0] == 0x0d && raw[1] == hi && raw[2] == lo:
		tail = raw[3:]
	case len(raw) >= 2 && raw[0] == hi && raw[1] == lo:
		tail = raw[2:]
	}
	if len(tail) != 0 {
		r.Tail = tail
		on := tail[0] != 0xf0
		r.PumpOn = &on
	}
	return r
}

func (s boilerRegisterWrite) acknowledgement(reply *protocol.Frame) error {
	if reply == nil {
		return errors.New("b509 write returned nil response")
	}
	if !s.accepts(reply.Data) {
		return fmt.Errorf("b509 write ack invalid: target=0x%02x addr=0x%04x", s.Target, s.Register)
	}
	return nil
}

func (s boilerRegisterWrite) accepts(raw []byte) bool {
	hi, lo := byte(s.Register>>8), byte(s.Register)
	switch {
	case len(raw) >= 3 && raw[0] == 0x0e && raw[1] == hi && raw[2] == lo:
		raw = raw[3:]
	case len(raw) >= 2 && raw[0] == hi && raw[1] == lo:
		raw = raw[2:]
	}
	return len(raw) == 0 || raw[0] != 0
}
