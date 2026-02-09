package system

type B524GroupProfile struct {
	Group       byte
	Opcode      byte
	InstanceMax byte
	RegisterMax uint16
}

const (
	B524OpcodeLocal  = byte(0x02)
	B524OpcodeRemote = byte(0x06)
)

var B524DiscoveryProfiles = []B524GroupProfile{
	{Group: 0x02, Opcode: B524OpcodeLocal, InstanceMax: 0x0A, RegisterMax: 0x0021},
	{Group: 0x03, Opcode: B524OpcodeLocal, InstanceMax: 0x0A, RegisterMax: 0x002F},
	{Group: 0x09, Opcode: B524OpcodeRemote, InstanceMax: 0x0A, RegisterMax: 0x002F},
	{Group: 0x0A, Opcode: B524OpcodeRemote, InstanceMax: 0x0A, RegisterMax: 0x003F},
	{Group: 0x0C, Opcode: B524OpcodeRemote, InstanceMax: 0x0A, RegisterMax: 0x003F},
}
