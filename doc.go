// Package ebusreg is the module root of the Helianthus eBUS registry and
// projection layer. It declares no API of its own; importing it provides
// nothing. The functionality lives in the subpackages:
//
//   - registry: device discovery, address faces, identity admission and the
//     typed planes and projections attached to each device.
//   - schema: sequential payload schemas and hardware-version-dependent
//     schema selection.
//   - router: method invocation and broadcast dispatch over an eBUS
//     transport supplied by the caller.
//   - providers/vaillant and vaillant/...: the default plane providers for
//     Vaillant-group devices.
//   - catalog/...: embedded, qualification-graded catalogues of eBUS service
//     fields and views.
//
// The module sits between the transport library below it, which owns
// framing, arbitration and bus access, and the gateway above it, which owns
// runtime wiring and every API surface. Neither transport nor consumer
// policy belongs here.
package ebusreg
