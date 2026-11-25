package ebus

import (
	"encoding/binary"
	"errors"
	"math"
	"strconv"

	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
)

type interfaceIdentityField byte

const (
	interfaceAPI      interfaceIdentityField = 0x06
	interfaceRevision interfaceIdentityField = 0x07
)

type interfaceIdentityScope struct {
	Source, Target       byte
	Observation, Session uint64
}

type interfaceIdentityRead struct {
	Field interfaceIdentityField
	Scope interfaceIdentityScope
}

func (s interfaceIdentityRead) request() protocol.Frame {
	return protocol.Frame{Source: s.Scope.Source, Target: s.Scope.Target, Primary: 0xb5, Secondary: 0x24, Data: []byte{0, byte(s.Field), 0}}
}

type interfaceIdentityValue struct {
	Read      interfaceIdentityRead
	Present   bool
	Raw       []byte
	Qualifier string
}

func (s interfaceIdentityRead) qualify(reply *protocol.Frame) interfaceIdentityValue {
	r := interfaceIdentityValue{Read: s, Present: reply != nil}
	if reply == nil {
		return r
	}
	r.Raw = reply.Data
	if len(r.Raw) != 4 {
		return r
	}
	x := float64(math.Float32frombits(binary.LittleEndian.Uint32(r.Raw)))
	if math.IsNaN(x) || x < 0 || x > 1<<24 || math.Trunc(x) != x {
		return r
	}
	r.Qualifier = strconv.FormatUint(uint64(x), 10)
	return r
}

func matchingInterfaceIdentity(api, revision interfaceIdentityValue) bool {
	return api.Read.Field == interfaceAPI && revision.Read.Field == interfaceRevision &&
		api.Read.Scope == revision.Read.Scope && api.Qualifier != "" && revision.Qualifier != ""
}

type localRegisterWrite struct {
	Opcode, Group, Instance byte
	Register                uint16
	Encoded                 []byte
}

func (s localRegisterWrite) request(source, target byte) protocol.Frame {
	data := binary.LittleEndian.AppendUint16([]byte{s.Opcode, 1, s.Group, s.Instance}, s.Register)
	return protocol.Frame{Source: source, Target: target, Primary: 0xb5, Secondary: 0x24, Data: append(data, s.Encoded...)}
}

func localRegisterAcknowledgement(reply *protocol.Frame) error {
	if reply == nil {
		return errors.New("b524 write returned nil response")
	}
	return nil
}
