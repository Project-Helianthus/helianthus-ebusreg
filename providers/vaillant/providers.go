package vaillant

import (
	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/dhw"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/heating"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/solar"
	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/system"
)

func System() registry.PlaneProvider {
	return system.NewProvider()
}

func Heating() registry.PlaneProvider {
	return heating.NewProvider()
}

func DHW() registry.PlaneProvider {
	return dhw.NewProvider()
}

func Solar() registry.PlaneProvider {
	return solar.NewProvider()
}

func Default() []registry.PlaneProvider {
	return []registry.PlaneProvider{
		System(),
		Heating(),
		DHW(),
		Solar(),
	}
}
