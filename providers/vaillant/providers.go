package vaillant

import (
	"github.com/d3vi1/helianthus-ebusreg/registry"
	"github.com/d3vi1/helianthus-ebusreg/vaillant/system"
	"github.com/d3vi1/helianthus-ebusreg/vendor/vaillant/dhw"
	"github.com/d3vi1/helianthus-ebusreg/vendor/vaillant/heating"
	"github.com/d3vi1/helianthus-ebusreg/vendor/vaillant/solar"
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
