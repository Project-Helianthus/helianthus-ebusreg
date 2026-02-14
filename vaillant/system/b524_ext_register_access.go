package system

import (
	"fmt"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
	"github.com/d3vi1/helianthus-ebusgo/types"
)

const (
	methodGetExtRegister = "get_ext_register"
	methodSetExtRegister = "set_ext_register"
)

const (
	extRegisterOpLocal  = byte(0x02)
	extRegisterOpRemote = byte(0x06)
	extRegisterOpRead   = byte(0x00)
	extRegisterOpWrite  = byte(0x01)
)

type extRegisterReadTemplate struct {
	primary   byte
	secondary byte
}

func (template extRegisterReadTemplate) Primary() byte {
	return template.primary
}

func (template extRegisterReadTemplate) Secondary() byte {
	return template.secondary
}

func (template extRegisterReadTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("ext register read template missing params: %w", ebuserrors.ErrInvalidPayload)
	}

	opcode, err := extRegisterOpcode(params)
	if err != nil {
		return nil, fmt.Errorf("ext register read template opcode: %w", err)
	}
	group, ok := uint8Param(params, "group")
	if !ok {
		return nil, fmt.Errorf("ext register read template group: %w", ebuserrors.ErrInvalidPayload)
	}
	instance, ok := uint8Param(params, "instance")
	if !ok {
		return nil, fmt.Errorf("ext register read template instance: %w", ebuserrors.ErrInvalidPayload)
	}
	addr, ok, err := registerAddrParam(params, "addr")
	if err != nil {
		return nil, fmt.Errorf("ext register read template addr: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("ext register read template addr: %w", ebuserrors.ErrInvalidPayload)
	}

	return []byte{opcode, extRegisterOpRead, group, instance, byte(addr), byte(addr >> 8)}, nil
}

type extRegisterWriteTemplate struct {
	primary   byte
	secondary byte
}

func (template extRegisterWriteTemplate) Primary() byte {
	return template.primary
}

func (template extRegisterWriteTemplate) Secondary() byte {
	return template.secondary
}

func (template extRegisterWriteTemplate) Build(params map[string]any) ([]byte, error) {
	if params == nil {
		return nil, fmt.Errorf("ext register write template missing params: %w", ebuserrors.ErrInvalidPayload)
	}

	opcode, err := extRegisterOpcode(params)
	if err != nil {
		return nil, fmt.Errorf("ext register write template opcode: %w", err)
	}
	group, ok := uint8Param(params, "group")
	if !ok {
		return nil, fmt.Errorf("ext register write template group: %w", ebuserrors.ErrInvalidPayload)
	}
	instance, ok := uint8Param(params, "instance")
	if !ok {
		return nil, fmt.Errorf("ext register write template instance: %w", ebuserrors.ErrInvalidPayload)
	}
	addr, ok, err := registerAddrParam(params, "addr")
	if err != nil {
		return nil, fmt.Errorf("ext register write template addr: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("ext register write template addr: %w", ebuserrors.ErrInvalidPayload)
	}

	data, hasData, err := bytesParam(params, "data")
	if err != nil {
		return nil, fmt.Errorf("ext register write template data: %w", err)
	}
	if !hasData {
		data, _, err = bytesParam(params, "payload")
		if err != nil {
			return nil, fmt.Errorf("ext register write template payload: %w", err)
		}
	}

	if len(data) > 0xF9 {
		return nil, fmt.Errorf("ext register write template data too long: %w", ebuserrors.ErrInvalidPayload)
	}

	payload := make([]byte, 0, 6+len(data))
	payload = append(payload, opcode, extRegisterOpWrite, group, instance, byte(addr), byte(addr>>8))
	payload = append(payload, data...)
	return payload, nil
}

func decodeExtRegisterResponse(cmd byte, opcode byte, group, instance byte, addr uint16, payload []byte) map[string]types.Value {
	values := map[string]types.Value{
		"optype":   {Value: cmd, Valid: true},
		"opcode":   {Value: opcode, Valid: true},
		"group":    {Value: group, Valid: true},
		"instance": {Value: instance, Valid: true},
		"addr":     {Value: addr, Valid: true},
		"addr_hex": {Value: fmt.Sprintf("%04X", addr), Valid: true},
		"payload":  {Value: append([]byte(nil), payload...), Valid: true},
	}

	constraint, hasConstraint := lookupB524Constraint(group, addr)

	if len(payload) >= 4 {
		values["prefix"] = types.Value{Value: append([]byte(nil), payload[:4]...), Valid: true}
		replyGroup := payload[1]
		replyAddr := uint16(payload[3])<<8 | uint16(payload[2])
		values["reply_group"] = types.Value{Value: replyGroup, Valid: true}
		values["reply_addr"] = types.Value{Value: replyAddr, Valid: true}
		values["reply_addr_hex"] = types.Value{Value: fmt.Sprintf("%04X", replyAddr), Valid: true}
		if resolved, ok := lookupB524Constraint(replyGroup, replyAddr); ok {
			constraint = resolved
			hasConstraint = true
		}
	} else {
		values["prefix"] = types.Value{Valid: false}
	}

	if len(payload) >= 1 {
		values["reply_kind"] = types.Value{Value: payload[0], Valid: true}
	} else {
		values["reply_kind"] = types.Value{Valid: false}
	}

	if len(payload) == 1 && payload[0] == 0x00 {
		values["value"] = types.Value{Valid: false}
		addB524ConstraintValues(values, constraint, hasConstraint)
		return values
	}
	if len(payload) <= 4 {
		values["value"] = types.Value{Valid: false}
		addB524ConstraintValues(values, constraint, hasConstraint)
		return values
	}

	values["value"] = types.Value{Value: append([]byte(nil), payload[4:]...), Valid: true}
	addB524ConstraintValues(values, constraint, hasConstraint)
	return values
}

func addB524ConstraintValues(values map[string]types.Value, constraint b524Constraint, hasConstraint bool) {
	if !hasConstraint {
		return
	}
	values["constraints"] = types.Value{Value: constraint.mapValue(), Valid: true}
	values["constraint_record"] = types.Value{Value: values["addr"].Value, Valid: true}
	values["constraint_record_hex"] = types.Value{Value: values["addr_hex"].Value, Valid: true}
	values["constraint_type"] = types.Value{Value: constraint.Type, Valid: true}
	values["constraint_min"] = types.Value{Value: constraint.Min, Valid: true}
	values["constraint_max"] = types.Value{Value: constraint.Max, Valid: true}
	values["constraint_step"] = types.Value{Value: constraint.Step, Valid: true}
}

func extRegisterOpcode(params map[string]any) (byte, error) {
	if params == nil {
		return extRegisterOpLocal, nil
	}
	if opcode, ok := uint8Param(params, "opcode"); ok {
		if opcode != extRegisterOpLocal && opcode != extRegisterOpRemote {
			return 0, ebuserrors.ErrInvalidPayload
		}
		return opcode, nil
	}
	if opcode, ok := uint8Param(params, "op"); ok {
		if opcode != extRegisterOpLocal && opcode != extRegisterOpRemote {
			return 0, ebuserrors.ErrInvalidPayload
		}
		return opcode, nil
	}
	return extRegisterOpLocal, nil
}
