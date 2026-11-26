package input

import "encoding/binary"

const (
	primaryExtendedRegister   byte = 0xb5
	secondaryExtendedRegister byte = 0x24
	primaryBoilerRegister     byte = 0xb5
	secondaryBoilerRegister   byte = 0x09
	primarySchedule           byte = 0xb5
	secondarySchedule         byte = 0x55
	primaryDiagnostic         byte = 0xb5
	secondaryDiagnostic       byte = 0x03
)

const (
	// OpcodeLocal is the B524 opcode of the controller's own register groups.
	OpcodeLocal byte = 0x02
	// OpcodeRemote is the B524 opcode of the remote device slot groups.
	OpcodeRemote byte = 0x06
	// OperationRead is the B524 operation byte of a register read.
	OperationRead byte = 0x00
	// OperationWrite is the B524 operation byte of a register write.
	OperationWrite byte = 0x01
	// BoilerRead is the B509 opcode of a register read.
	BoilerRead byte = 0x0d
	// BoilerWrite is the B509 opcode of a register write.
	BoilerWrite byte = 0x0e
	// InterfaceIdentityAPI is the identity id of the interface API version.
	InterfaceIdentityAPI byte = 0x06
	// InterfaceIdentityRevision is the identity id of the interface revision.
	InterfaceIdentityRevision byte = 0x07
)

const (
	scheduleConfigOpcode byte = 0xa3
	scheduleCountsOpcode byte = 0xa4
	scheduleReadOpcode   byte = 0xa5
	scheduleWriteOpcode  byte = 0xa6
)

// Service is the primary and secondary command byte pair of a request.
type Service struct{ Primary, Secondary byte }

// Selector is the wire geometry of one request: its service and its payload.
type Selector interface {
	Service() Service
	Payload() []byte
}

// B524Selector is a B524 register request: opcode, operation, group, instance and register.
type B524Selector struct {
	Opcode, Operation byte
	Group, Instance   byte
	Register          uint16
}

// Service returns the B524 command pair.
func (s B524Selector) Service() Service {
	return Service{Primary: primaryExtendedRegister, Secondary: secondaryExtendedRegister}
}

// Payload returns the six-byte header with the register in little-endian order.
func (s B524Selector) Payload() []byte {
	return binary.LittleEndian.AppendUint16([]byte{s.Opcode, s.Operation, s.Group, s.Instance}, s.Register)
}

// WithValue returns the header followed by the value bytes of a write.
func (s B524Selector) WithValue(value []byte) []byte {
	return append(s.Payload(), value...)
}

// B509Selector is a B509 register request: opcode and register.
type B509Selector struct {
	Opcode   byte
	Register uint16
}

// Service returns the B509 command pair.
func (s B509Selector) Service() Service {
	return Service{Primary: primaryBoilerRegister, Secondary: secondaryBoilerRegister}
}

// Payload returns the opcode followed by the register in big-endian order.
func (s B509Selector) Payload() []byte {
	return binary.BigEndian.AppendUint16([]byte{s.Opcode}, s.Register)
}

// WithValue returns the header followed by the value bytes of a write.
func (s B509Selector) WithValue(value []byte) []byte {
	return append(s.Payload(), value...)
}

// InterfaceIdentitySelector is the three-byte B524 interface identity request.
type InterfaceIdentitySelector struct{ ID byte }

// Service returns the B524 command pair.
func (s InterfaceIdentitySelector) Service() Service {
	return Service{Primary: primaryExtendedRegister, Secondary: secondaryExtendedRegister}
}

// Payload returns zero, the identity id and zero.
func (s InterfaceIdentitySelector) Payload() []byte { return []byte{0x00, s.ID, 0x00} }

// ScheduleConfigSelector is the B555 program configuration read of one zone and circuit.
type ScheduleConfigSelector struct{ Zone, Circuit byte }

// Service returns the B555 command pair.
func (s ScheduleConfigSelector) Service() Service {
	return Service{Primary: primarySchedule, Secondary: secondarySchedule}
}

// Payload returns the opcode, zone and circuit.
func (s ScheduleConfigSelector) Payload() []byte {
	return []byte{scheduleConfigOpcode, s.Zone, s.Circuit}
}

// ScheduleCountsSelector is the B555 slots-per-weekday read of one zone and circuit.
type ScheduleCountsSelector struct{ Zone, Circuit byte }

// Service returns the B555 command pair.
func (s ScheduleCountsSelector) Service() Service {
	return Service{Primary: primarySchedule, Secondary: secondarySchedule}
}

// Payload returns the opcode, zone and circuit.
func (s ScheduleCountsSelector) Payload() []byte {
	return []byte{scheduleCountsOpcode, s.Zone, s.Circuit}
}

// ScheduleTimerSelector is the B555 read of one timer slot.
type ScheduleTimerSelector struct{ Zone, Circuit, Weekday, Index byte }

// Service returns the B555 command pair.
func (s ScheduleTimerSelector) Service() Service {
	return Service{Primary: primarySchedule, Secondary: secondarySchedule}
}

// Payload returns the opcode, zone, circuit, weekday and slot index.
func (s ScheduleTimerSelector) Payload() []byte {
	return []byte{scheduleReadOpcode, s.Zone, s.Circuit, s.Weekday, s.Index}
}

// ScheduleTimerWriteSelector is the B555 write of one timer slot.
type ScheduleTimerWriteSelector struct {
	Zone, Circuit, Weekday, Index              byte
	Count                                      byte
	StartHour, StartMinute, EndHour, EndMinute byte
	TemperatureRaw                             uint16
}

// Service returns the B555 command pair.
func (s ScheduleTimerWriteSelector) Service() Service {
	return Service{Primary: primarySchedule, Secondary: secondarySchedule}
}

// Payload returns the write opcode, the slot position, the slot count, the times and the raw temperature in little-endian order.
func (s ScheduleTimerWriteSelector) Payload() []byte {
	head := []byte{scheduleWriteOpcode, s.Zone, s.Circuit, s.Weekday, s.Index, s.Count, s.StartHour, s.StartMinute, s.EndHour, s.EndMinute}
	return binary.LittleEndian.AppendUint16(head, s.TemperatureRaw)
}

// DiagnosticSelector is a B503 request body: family, selector and an optional history index.
type DiagnosticSelector struct {
	Family, Selector byte
	Index            *byte
}

// Service returns the B503 command pair.
func (s DiagnosticSelector) Service() Service {
	return Service{Primary: primaryDiagnostic, Secondary: secondaryDiagnostic}
}

// Payload returns the family and selector, followed by the index when one is set.
func (s DiagnosticSelector) Payload() []byte {
	body := []byte{s.Family, s.Selector}
	if s.Index != nil {
		body = append(body, *s.Index)
	}
	return body
}

// LiveMonitorAction is the session-managed purpose of a live-monitor request.
type LiveMonitorAction string

const (
	// LiveMonitorEnable opens the live-monitor session.
	LiveMonitorEnable LiveMonitorAction = "enable"
	// LiveMonitorRead reads the live-monitor main record.
	LiveMonitorRead LiveMonitorAction = "read"
	// LiveMonitorDisable closes the live-monitor session.
	LiveMonitorDisable LiveMonitorAction = "disable"
)

// LiveMonitorSelector is the B503 live-monitor request. All three actions send the same body.
type LiveMonitorSelector struct{ Action LiveMonitorAction }

// Service returns the B503 command pair.
func (s LiveMonitorSelector) Service() Service {
	return Service{Primary: primaryDiagnostic, Secondary: secondaryDiagnostic}
}

// Payload returns the body that every action sends.
func (s LiveMonitorSelector) Payload() []byte { return []byte{0x00, 0x03} }

// Key returns the replay key of the action.
func (s LiveMonitorSelector) Key() string { return "live_monitor_" + string(s.Action) }
