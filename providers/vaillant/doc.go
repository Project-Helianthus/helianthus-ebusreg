// Package vaillant assembles the default plane providers for Vaillant-group
// devices: the system, heating, domestic hot water, solar and timer
// providers from the vaillant/... packages, in that order. Each call
// returns new provider values; the package keeps no state and adds no
// matching or decoding of its own.
//
// A device registry built from this list creates every plane whose provider
// matches the device, so the order only fixes the order in which planes and
// projections are listed for a device.
package vaillant
