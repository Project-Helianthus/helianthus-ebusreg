package writes

import (
	"encoding/binary"
	"fmt"
	"sync"
)

// RegisterIdentity is one opcode, group, instance and register read from a reply.
type RegisterIdentity struct {
	Opcode   byte
	Group    byte
	Instance byte
	Register uint16
}

// Declared records registers whose correlated reply sets the write attribute.
// Observing a reply does not send a frame.
type Declared struct {
	mu   sync.Mutex
	seen map[RegisterIdentity]struct{}
}

// Observe admits identity when reply is a correlated header and bit 1 of the
// attribute byte is set. A later correlated reply with that bit clear withdraws
// it. A short or uncorrelated reply leaves the record unchanged.
func (d *Declared) Observe(opcode, group, instance byte, register uint16, reply []byte) {
	if d == nil {
		return
	}
	flags, ok := correlatedHeader(group, register, reply)
	if !ok {
		return
	}
	id := RegisterIdentity{Opcode: opcode, Group: group, Instance: instance, Register: register}
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.seen == nil {
		d.seen = map[RegisterIdentity]struct{}{}
	}
	if flags&0x02 != 0 {
		d.seen[id] = struct{}{}
		return
	}
	delete(d.seen, id)
}

// OperationID returns the operation of one admitted identity.
// A pre-port selector with the same family, opcode, group and register keeps
// that operation. Any other identity gets an id of its own four fields.
func (d *Declared) OperationID(id RegisterIdentity) (string, bool) {
	if !d.admitted(id) {
		return "", false
	}
	if operation, ok := catalogOperation("local", id.Opcode, id.Group, id.Register); ok {
		return operation, true
	}
	return fmt.Sprintf("ebusreg.operation.declared.%02x.%02x.%02x.%04x", id.Opcode, id.Group, id.Instance, id.Register), true
}

// OperationForSelector returns the pre-port operation shared by every instance
// of family, opcode, group and register, when one admitted reply exists.
func (d *Declared) OperationForSelector(family string, opcode, group byte, register uint16) (string, bool) {
	if d == nil || !d.hasSelector(opcode, group, register) {
		return "", false
	}
	return catalogOperation(family, opcode, group, register)
}

// Grade is hypothesis for an operation admitted from a device reply.
func (d *Declared) Grade(id RegisterIdentity) (string, bool) {
	if _, ok := d.OperationID(id); !ok {
		return "", false
	}
	return "hypothesis", true
}

func (d *Declared) admitted(id RegisterIdentity) bool {
	if d == nil {
		return false
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.seen[id]
	return ok
}

func (d *Declared) hasSelector(opcode, group byte, register uint16) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	for id := range d.seen {
		if id.Opcode == opcode && id.Group == group && id.Register == register {
			return true
		}
	}
	return false
}

func correlatedHeader(group byte, register uint16, reply []byte) (byte, bool) {
	if len(reply) < 4 {
		return 0, false
	}
	if reply[1] != group || binary.LittleEndian.Uint16(reply[2:4]) != register {
		return 0, false
	}
	return reply[0], true
}

func catalogOperation(family string, opcode, group byte, register uint16) (string, bool) {
	rows, err := cachedInventory()
	if err != nil {
		return "", false
	}
	var found string
	for _, row := range rows {
		numbers, ok := parseSelector(row.Selector)
		if !ok || row.Selector.Family != family || numbers.opcode != opcode || numbers.group != group || numbers.register != register {
			continue
		}
		if found == "" {
			found = row.Operation
			continue
		}
		if found != row.Operation {
			return "", false
		}
	}
	return found, found != ""
}
