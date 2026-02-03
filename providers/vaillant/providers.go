package vaillant

import (
	"github.com/d3vi1/helianthus-ebusreg/registry"
	"github.com/d3vi1/helianthus-ebusreg/vaillant/dhw"
	"github.com/d3vi1/helianthus-ebusreg/vaillant/heating"
	"github.com/d3vi1/helianthus-ebusreg/vaillant/solar"
	"github.com/d3vi1/helianthus-ebusreg/vaillant/system"
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
