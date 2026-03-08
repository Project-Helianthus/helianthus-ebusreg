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

// B524DiscoveryProfiles defines the maximum exploration bounds for B524
// register discovery per group. Dual-namespace groups (0x08, 0x09, 0x0A)
// have separate profiles for local (OP=0x02) and remote (OP=0x06) opcodes.
//
// Over-estimate policy: where existing profile max exceeds scan-observed max
// (GG=0x0A Remote: profile 0x003F > scan 0x0035, GG=0x0C Remote: profile
// 0x003F > scan 0x002F), existing higher values are kept conservatively as
// the scan may miss optional registers.
var B524DiscoveryProfiles = []B524GroupProfile{
	// Singleton groups (II_max=0x00)
	{Group: 0x00, Opcode: B524OpcodeLocal, InstanceMax: 0x00, RegisterMax: 0x00FF}, // Regulator Parameters
	{Group: 0x01, Opcode: B524OpcodeLocal, InstanceMax: 0x00, RegisterMax: 0x0013}, // DHW

	// Multi-instance local groups
	{Group: 0x02, Opcode: B524OpcodeLocal, InstanceMax: 0x0A, RegisterMax: 0x0025}, // Heating Circuits
	{Group: 0x03, Opcode: B524OpcodeLocal, InstanceMax: 0x0A, RegisterMax: 0x002F}, // Zones

	// Dual-namespace groups: local (OP=0x02) config + remote (OP=0x06) live data
	{Group: 0x08, Opcode: B524OpcodeLocal, InstanceMax: 0x00, RegisterMax: 0x0007},  // Buffer/Solar Cylinder 2 (singleton)
	{Group: 0x08, Opcode: B524OpcodeRemote, InstanceMax: 0x0A, RegisterMax: 0x0004}, // Buffer/Solar Cylinder 2 (remote)
	{Group: 0x09, Opcode: B524OpcodeLocal, InstanceMax: 0x0A, RegisterMax: 0x000F},  // Room Sensors (local config)
	{Group: 0x09, Opcode: B524OpcodeRemote, InstanceMax: 0x0A, RegisterMax: 0x002F}, // Room Sensors (remote live)
	{Group: 0x0A, Opcode: B524OpcodeLocal, InstanceMax: 0x0A, RegisterMax: 0x004D},  // Room State (local config)
	{Group: 0x0A, Opcode: B524OpcodeRemote, InstanceMax: 0x0A, RegisterMax: 0x003F}, // Room State (remote live)

	// Remote-only groups
	{Group: 0x0C, Opcode: B524OpcodeRemote, InstanceMax: 0x0A, RegisterMax: 0x003F}, // Unrecognized
}
