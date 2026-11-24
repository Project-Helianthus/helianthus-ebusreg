// Package wolf catalogues the Wolf eBUS messages as the public Wolf page
// documents them: the burner-controller broadcasts (PB 05 and PB 08), the
// solar extensions (PB 50 SB 14, 17 and 18) and the service-parameter access
// (PB 50 SB 22 read and SB 23 write). Identity comes from the manufacturer byte
// 0x19 that the page states; PB 50 is a command byte, not a manufacturer.
//
// Every documented message the page marks as a write or command has an
// encoder; the plane methods in providers/wolf call them by hand, and nothing
// is sent unless a method is called. The package decides no support and no
// qualification. A value the page does not define is kept raw and unnamed in
// its unit.
package wolf
