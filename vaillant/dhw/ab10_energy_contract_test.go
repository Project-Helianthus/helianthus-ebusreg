package dhw

import (
	ab10bytes "bytes"
	ab10testing "testing"

	"github.com/Project-Helianthus/helianthus-ebusreg/registry"
)

func TestAb10ActiveThreeByteSelector(ab10t *ab10testing.T) {
	for _, ab10case := range []struct {
		ab10name                          string
		ab10period, ab10source, ab10usage byte
	}{
		{"ab10_control", 1, 4, 3},
		{"ab10_period_only", 2, 4, 3},
		{"ab10_source_only", 1, 9, 3},
		{"ab10_usage_only", 1, 4, 4},
	} {
		ab10t.Run(ab10case.ab10name, func(ab10t *ab10testing.T) {
			ab10got, ab10err := (energyTemplate{primary: 0xB5, secondary: 0x16}).Build(map[string]any{
				"period": ab10case.ab10period,
				"source": ab10case.ab10source,
				"usage":  ab10case.ab10usage,
			})
			ab10want := []byte{ab10case.ab10period, ab10case.ab10source, ab10case.ab10usage}
			if ab10err != nil || !ab10bytes.Equal(ab10got, ab10want) {
				ab10t.Fatalf("active Build = %x, %v; want exactly %x", ab10got, ab10err, ab10want)
			}
		})
	}
}

func TestAb10DescriptorInclusionThreshold(ab10t *ab10testing.T) {
	for _, ab10case := range []struct {
		ab10hardware string
		ab10want     int
	}{
		{"7602", 0},
		{"7603", 1},
		{"7604", 1},
	} {
		ab10t.Run("ab10_"+ab10case.ab10hardware, func(ab10t *ab10testing.T) {
			ab10planes := NewProvider().CreatePlanes(registry.DeviceInfo{Manufacturer: "Vaillant", HardwareVersion: ab10case.ab10hardware})
			if len(ab10planes) != 1 {
				ab10t.Fatalf("plane count = %d; want 1", len(ab10planes))
			}
			ab10count := 0
			for _, ab10method := range ab10planes[0].Methods() {
				if ab10method.Name() == methodEnergyStats {
					ab10count++
				}
			}
			if ab10count != ab10case.ab10want {
				ab10t.Fatalf("energy descriptor count = %d; want %d", ab10count, ab10case.ab10want)
			}
		})
	}
}
