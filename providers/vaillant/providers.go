package vaillant

import (
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/b555"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/dhw"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/heating"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/solar"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/system"
)

// System returns a new provider for the Vaillant system plane.
func System() registry.PlaneProvider {
	return system.NewProvider()
}

// Heating returns a new provider for the Vaillant heating plane.
func Heating() registry.PlaneProvider {
	return heating.NewProvider()
}

// DHW returns a new provider for the Vaillant domestic hot water plane.
func DHW() registry.PlaneProvider {
	return dhw.NewProvider()
}

// Solar returns a new provider for the Vaillant solar plane.
func Solar() registry.PlaneProvider {
	return solar.NewProvider()
}

// Timer returns a new provider for the Vaillant B555 weekly timer plane.
func Timer() registry.PlaneProvider {
	return b555.NewProvider()
}

// Default returns a fresh slice holding one new instance of every Vaillant
// plane provider, in the order system, heating, DHW, solar, timer. Each
// provider decides on its own which devices it matches; the order is the
// order in which a registry appends their planes to a matched device.
func Default() []registry.PlaneProvider {
	return []registry.PlaneProvider{
		System(),
		Heating(),
		DHW(),
		Solar(),
		Timer(),
	}
}
