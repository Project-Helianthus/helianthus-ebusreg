// Package brink catalogues the Brink 40xx heat-recovery-ventilation protocol
// family as documented on the public page: process byte 0x40, a one-byte
// register ID space, the 4022 sensor and status reads, the 4050 self-describing
// parameter read, the 4080 parameter write, and the documented write
// sub-protocols 0x40FF, 0x4091, 0x40A1 and 0x40CB. Each table entry is the
// page's statement; nothing here is observed on a bus.
//
// The page identifies the family by primary byte 0x40 and states no
// manufacturer byte. The default address 0x3C can be changed, so neither value
// is a family key on the registry's device identity. The package decides no
// support and no qualification: its grades are the page's, and every
// reply that the page does not define is kept raw.
package brink
