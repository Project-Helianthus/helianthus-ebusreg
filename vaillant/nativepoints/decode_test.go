package nativepoints_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"math/big"
	"testing"

	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/nativepoints"
)

func float32Reply(value float32) []byte {
	return binary.LittleEndian.AppendUint32(nil, math.Float32bits(value))
}

func synthetic(codec nativepoints.Codec, divisor uint16, unit string) nativepoints.Point {
	return nativepoints.Point{
		DeviceClass: nativepoints.DeviceClassBoiler, Name: "synthetic_value", NameStatus: nativepoints.NameFinal,
		ValueClass: nativepoints.ValueNumeric, Codec: codec, ScaleDivisor: divisor, NativeUnit: unit, Exportable: true,
	}
}

func TestDecodeVectors(t *testing.T) {
	outdoor := mustLookup(t, nativepoints.DeviceClassController, "outdoor_temperature")
	preheat := mustLookup(t, nativepoints.DeviceClassController, "preheating_time_maximum")
	humidity := mustLookup(t, nativepoints.DeviceClassController, "zone_room_humidity")
	offset := mustLookup(t, nativepoints.DeviceClassController, "cylinder_loading_offset")
	flow := mustLookup(t, nativepoints.DeviceClassBoiler, "flow_temperature")
	water := mustLookup(t, nativepoints.DeviceClassBoiler, "dhw_water_flow_rate")
	fan := mustLookup(t, nativepoints.DeviceClassBoiler, "fan_speed")
	cases := []struct {
		label  string
		point  nativepoints.Point
		raw    []byte
		native string
		si     string
		unit   string
	}{
		{"float32 shortest decimal", outdoor, float32Reply(21.1), "21.1", "294.25", "K"},
		{"float32 negative", outdoor, float32Reply(-3.7), "-3.7", "269.45", "K"},
		{"float32 tenth", outdoor, float32Reply(0.1), "0.1", "273.25", "K"},
		{"float32 zero", outdoor, float32Reply(0), "0", "273.15", "K"},
		{"float32 percent", humidity, float32Reply(45.5), "45.5", "0.455", "1"},
		{"float32 kelvin difference", offset, float32Reply(5), "5", "5", "K"},
		{"uint16 minutes", preheat, []byte{0x5A, 0x00}, "90", "5400", "s"},
		{"DATA2c positive", flow, []byte{0x58, 0x02}, "37.5", "310.65", "K"},
		{"DATA2c negative", flow, []byte{0xF8, 0xFF}, "-0.5", "272.65", "K"},
		{"DATA2c sixteenth", flow, []byte{0x01, 0x00}, "0.0625", "273.2125", "K"},
		{"DATA2c lowest", flow, []byte{0x01, 0x80}, "-2047.9375", "-1774.7875", "K"},
		{"scaled UIN", water, []byte{0x2C, 0x01}, "3", "0.00005", "m3/s"},
		{"scaled UIN terminating", water, []byte{0x7B, 0x00}, "1.23", "0.0000205", "m3/s"},
		{"scaled UIN rational", water, []byte{0x01, 0x00}, "0.01", "1/6000000", "m3/s"},
		{"UIN revolutions", fan, []byte{0xB8, 0x0B}, "3000", "50", "1/s"},
		{"UIN revolutions rational", fan, []byte{0x01, 0x00}, "1", "1/60", "1/s"},
		{"DATA2b", synthetic(nativepoints.CodecData2b, 1, "°C"), []byte{0x80, 0x14}, "20.5", "293.65", "K"},
		{"bar", synthetic(nativepoints.CodecFloat32LE, 1, "bar"), float32Reply(1.5), "1.5", "150000", "Pa"},
		{"kilowatt hour", synthetic(nativepoints.CodecFloat32LE, 1, "kWh"), float32Reply(1.5), "1.5", "5400000", "J"},
		{"kilowatt", synthetic(nativepoints.CodecFloat32LE, 1, "kW"), float32Reply(2.5), "2.5", "2500", "W"},
		{"watt", synthetic(nativepoints.CodecUIN, 1, "W"), []byte{0x10, 0x00}, "16", "16", "W"},
		{"hour", synthetic(nativepoints.CodecUIN, 1, "h"), []byte{0x02, 0x00}, "2", "7200", "s"},
		{"second", synthetic(nativepoints.CodecUIN, 1, "s"), []byte{0x02, 0x00}, "2", "2", "s"},
		{"count", synthetic(nativepoints.CodecUIN, 1, "count"), []byte{0x07, 0x00}, "7", "7", "1"},
		{"litre per hour", synthetic(nativepoints.CodecUIN, 1, "L/h"), []byte{0x24, 0x00}, "36", "0.00001", "m3/s"},
		{"scale ten", synthetic(nativepoints.CodecUIN, 10, "K"), []byte{0x05, 0x00}, "0.5", "0.5", "K"},
	}
	for _, tc := range cases {
		raw := append([]byte(nil), tc.raw...)
		value, err := nativepoints.Decode(tc.point, raw)
		if err != nil {
			t.Errorf("%s: %v", tc.label, err)
			continue
		}
		if value.Native != tc.native || value.SI != tc.si || value.SIUnit != tc.unit || value.NativeUnit != tc.point.NativeUnit {
			t.Errorf("%s: got %+v; want native %s si %s %s", tc.label, value, tc.native, tc.si, tc.unit)
		}
		if !bytes.Equal(raw, tc.raw) {
			t.Errorf("%s: reply mutated", tc.label)
		}
		for _, text := range []string{value.Native, value.SI} {
			if _, ok := new(big.Rat).SetString(text); !ok {
				t.Errorf("%s: %q is not an exact rational", tc.label, text)
			}
		}
	}
}

func TestDecodeRefusals(t *testing.T) {
	outdoor := mustLookup(t, nativepoints.DeviceClassController, "outdoor_temperature")
	preheat := mustLookup(t, nativepoints.DeviceClassController, "preheating_time_maximum")
	flow := mustLookup(t, nativepoints.DeviceClassBoiler, "flow_temperature")
	water := mustLookup(t, nativepoints.DeviceClassBoiler, "dhw_water_flow_rate")
	pressure := mustLookup(t, nativepoints.DeviceClassBoiler, "water_pressure")
	unknownCodec := synthetic(nativepoints.CodecUnknown, 1, "°C")
	unknownUnit := synthetic(nativepoints.CodecUIN, 1, "")
	foreignUnit := synthetic(nativepoints.CodecUIN, 1, "steps per fortnight")
	zeroScale := synthetic(nativepoints.CodecUIN, 0, "K")
	cases := []struct {
		label   string
		point   nativepoints.Point
		raw     []byte
		failure nativepoints.DecodeFailure
	}{
		{"NaN", outdoor, float32Reply(float32(math.NaN())), nativepoints.DecodeNotANumber},
		{"quiet NaN pattern", outdoor, []byte{0x00, 0x00, 0xC0, 0x7F}, nativepoints.DecodeNotANumber},
		{"positive infinity", outdoor, float32Reply(float32(math.Inf(1))), nativepoints.DecodeInfinite},
		{"negative infinity", outdoor, float32Reply(float32(math.Inf(-1))), nativepoints.DecodeInfinite},
		{"float32 short", outdoor, []byte{0x00, 0x00, 0x80}, nativepoints.DecodeShortReply},
		{"float32 long", outdoor, []byte{0x00, 0x00, 0x80, 0x3F, 0x00}, nativepoints.DecodeLongReply},
		{"float32 empty", outdoor, nil, nativepoints.DecodeShortReply},
		{"uint16 short", preheat, []byte{0x5A}, nativepoints.DecodeShortReply},
		{"uint16 long", preheat, []byte{0x5A, 0x00, 0x00}, nativepoints.DecodeLongReply},
		{"DATA2c sentinel", flow, []byte{0x00, 0x80}, nativepoints.DecodeSentinel},
		{"DATA2c short", flow, []byte{0x58}, nativepoints.DecodeShortReply},
		{"DATA2c long", flow, []byte{0x58, 0x02, 0x00}, nativepoints.DecodeLongReply},
		{"DATA2b sentinel", synthetic(nativepoints.CodecData2b, 1, "°C"), []byte{0x00, 0x80}, nativepoints.DecodeSentinel},
		{"UIN sentinel", water, []byte{0xFF, 0xFF}, nativepoints.DecodeSentinel},
		{"UIN short", water, []byte{0x2C}, nativepoints.DecodeShortReply},
		{"UIN long", water, []byte{0x2C, 0x01, 0x00}, nativepoints.DecodeLongReply},
		{"withheld row", pressure, []byte{0x00, 0x01}, nativepoints.DecodeWithheld},
		{"unknown codec", unknownCodec, []byte{0x00, 0x01}, nativepoints.DecodeCodecUnknown},
		{"unknown unit", unknownUnit, []byte{0x00, 0x01}, nativepoints.DecodeUnitUnknown},
		{"unconvertible unit", foreignUnit, []byte{0x00, 0x01}, nativepoints.DecodeUnitUnknown},
		{"zero scale", zeroScale, []byte{0x00, 0x01}, nativepoints.DecodeCodecUnknown},
	}
	for _, tc := range cases {
		value, err := nativepoints.Decode(tc.point, tc.raw)
		var decodeError *nativepoints.DecodeError
		if !errors.As(err, &decodeError) {
			t.Errorf("%s: error %v (value %+v) is not a DecodeError", tc.label, err, value)
			continue
		}
		if decodeError.Failure != tc.failure || decodeError.Name != tc.point.Name || decodeError.DeviceClass != tc.point.DeviceClass {
			t.Errorf("%s: %+v; want failure %s", tc.label, decodeError, tc.failure)
		}
		if value != (nativepoints.Value{}) {
			t.Errorf("%s: refused decode returned %+v", tc.label, value)
		}
		if decodeError.Error() == "" {
			t.Errorf("%s: empty error text", tc.label)
		}
	}
}
