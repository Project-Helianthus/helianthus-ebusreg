package ebus

import (
	"encoding/binary"
	"errors"
	"maps"
	"math/big"
	"strconv"
	"strings"

	"github.com/Project-Helianthus/helianthus-ebusreg/vaillant/nativepoints"
	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	"github.com/Project-Helianthus/helianthus-semreg/semreg/v1/roweval"
)

// B524Selector holds the wire fields that identify a B5/24 register request.
type B524Selector struct {
	Primary, Secondary byte
	Opcode, Operation  byte
	Group, Instance    byte
	Register           uint16
}

func (s B524Selector) payload() []byte {
	return []byte{s.Opcode, s.Operation, s.Group, s.Instance, byte(s.Register), byte(s.Register >> 8)}
}

type registerReadPurpose string

const (
	scalarRead     registerReadPurpose = "scalar"
	discoveryProbe registerReadPurpose = "discovery"
)

type registerRead struct {
	Selector     B524Selector
	Purpose      registerReadPurpose
	Raw, Value   []byte
	InstanceEcho *byte
	Coherent     bool
	Reason       string
}

func correlateRegisterRead(s B524Selector, purpose registerReadPurpose, raw []byte, err error) registerRead {
	r := registerRead{Selector: s, Purpose: purpose, Raw: append([]byte(nil), raw...)}
	switch {
	case err != nil:
		r.Reason = "acquisition_error"
	case len(raw) == 0:
		r.Reason = "empty"
	case len(raw) == 1:
		r.Reason = "status_only"
	case len(raw) < 4:
		r.Reason = "short_header"
	default:
		r.correlate()
	}
	return r
}

func (r *registerRead) correlate() {
	raw, s := r.Raw, r.Selector
	var header int
	switch {
	case len(raw) >= 5 && raw[2] == s.Group && binary.LittleEndian.Uint16(raw[3:5]) == s.Register:
		r.InstanceEcho = &raw[1]
		if r.Purpose != discoveryProbe && raw[1] != s.Instance && (s.Instance == 0xff || raw[1] != s.Instance+1) {
			r.Reason = "unmatched_instance"
			return
		}
		header = 5
	case raw[1] == s.Group && binary.LittleEndian.Uint16(raw[2:4]) == s.Register:
		header = 4
	default:
		r.Reason = "unmatched_echo"
		return
	}
	r.Coherent = true
	if r.Purpose == discoveryProbe {
		return
	}
	if len(raw) == header {
		r.Reason = "header_only"
		return
	}
	r.Value = raw[header:]
}

const dhwFeature = "hwc_enabled"

func circuitTypeLabel(code uint16) string {
	switch code {
	case 0:
		return "off"
	case 1:
		return "on"
	default:
		return ""
	}
}

func point(name string) nativepoints.Point {
	p, _ := nativepoints.Lookup(nativepoints.DeviceClassController, name)
	return p
}

func pointSelector(p nativepoints.Point) B524Selector {
	loc := p.Locations[0]
	return B524Selector{Primary: p.Primary, Secondary: p.Secondary, Opcode: loc.Opcode, Group: loc.Group, Register: loc.Register}
}

func dhwGateSelector() B524Selector { return pointSelector(point("dhw_circuit_type")) }

func dhwFeatures(g registerRead) map[string]bool {
	if g.Selector != dhwGateSelector() || g.Reason != "" || !g.Coherent || g.Purpose != scalarRead {
		return nil
	}
	uin := nativepoints.Point{Codec: nativepoints.CodecUIN, ScaleDivisor: 1}
	decoded, err := nativepoints.DecodeNative(uin, g.Value)
	if err != nil {
		return nil
	}
	return map[string]bool{dhwFeature: decoded.Native != "0"}
}

var thermalPack = semreg.PackRef{ID: "helianthus.pack.thermal", Version: "1.0.0"}

func thermalRef(id semreg.DefinitionID) semreg.DefinitionRef {
	return semreg.DefinitionRef{Pack: thermalPack, ID: id, Version: "1.0.0"}
}

func dhwTable(policy semreg.PolicyID) roweval.Table[B524Selector] {
	setpoint, measured := point("dhw_temperature_setpoint"), point("dhw_temperature")
	return roweval.Table[B524Selector]{
		Native: []roweval.NativeRow[B524Selector]{dhwNative(setpoint), dhwNative(measured)},
		Bindings: []roweval.BindingRow{
			dhwBinding(setpoint.Name, "thermal.setpoint.dhw_temperature_setpoint", policy, map[string]string{
				"direction": "controller_to_initiator", "presentation": "controller",
			}),
			dhwBinding(setpoint.Name, "thermal.setpoint.dhw_temperature_setpoint", policy, map[string]string{
				"direction": "controller_to_initiator", "presentation": "controller_mirror",
			}),
			dhwBinding(measured.Name, "thermal.measurement.dhw_temperature", policy, map[string]string{
				"direction": "target_to_initiator",
			}),
		},
		Hooks: map[string]roweval.RowHook{"dhw_celsius_to_kelvin": dhwKelvin},
	}
}

func dhwNative(p nativepoints.Point) roweval.NativeRow[B524Selector] {
	loc := p.Locations[0]
	return roweval.NativeRow[B524Selector]{
		ID: p.Name, Unit: loc.NativeUnit, At: pointSelector(p), When: roweval.When{Requires: []string{dhwFeature}},
		Wire:   roweval.Layout{Width: 32, Codec: string(loc.Codec), BitOrder: "lsb0", ByteOrder: "little_endian", WordOrder: "native"},
		Decode: roweval.Convert{Numeric: &roweval.Affine{Num: 1, Den: int64(loc.ScaleDivisor)}},
	}
}

func dhwBinding(nativeID string, fact semreg.DefinitionID, policy semreg.PolicyID, equal map[string]string) roweval.BindingRow {
	return roweval.BindingRow{
		NativeID: nativeID, SubjectRole: "resource", Service: thermalRef("thermal.service.dhw"),
		Capability: thermalRef("thermal.capability.read.dhw"), When: roweval.When{Equal: equal},
		Read: &roweval.ReadRow{Fact: thermalRef(fact), Unit: "unit.kelvin", Policy: policy,
			Dimensions: map[semreg.DefinitionID]string{"thermal.dimension.dhw": "resource"},
			ToRegistry: roweval.Convert{Numeric: &roweval.Affine{Num: 1, Den: 1, Add: kelvinOffset}, Hook: "dhw_celsius_to_kelvin"}},
	}
}

var kelvinOffset = semreg.Decimal{Coefficient: "27315", Exponent10: -2}

func dhwKelvin(in roweval.DetachedInputs, _ roweval.RowID, _ roweval.Direction, kelvin roweval.TypedValue) (roweval.DerivedValue, roweval.Reason) {
	if kelvin.Quantity == nil {
		return roweval.DerivedValue{}, roweval.ReasonConversion
	}
	k, ok := new(big.Rat).SetString(decimalOperand(kelvin.Quantity.Number))
	offset, offsetOK := new(big.Rat).SetString(decimalOperand(kelvinOffset))
	if !ok || !offsetOK {
		return roweval.DerivedValue{}, roweval.ReasonConversion
	}
	celsius := k.Sub(k, offset)
	q := in.Qualifiers().Values
	for _, b := range []struct {
		key string
		cmp int
	}{{"native_min", -1}, {"native_max", 1}} {
		if text, present := q[b.key]; present {
			bound, valid := new(big.Rat).SetString(text)
			if !valid {
				return roweval.DerivedValue{}, roweval.ReasonConversion
			}
			if celsius.Cmp(bound) == b.cmp {
				return roweval.DerivedValue{}, roweval.ReasonOutsideDomain
			}
		}
	}
	return roweval.DerivedValue{Value: kelvin}, ""
}

func decimalOperand(d semreg.Decimal) string {
	return d.Coefficient + "e" + strconv.FormatInt(int64(d.Exponent10), 10)
}

func nativeDecimal(s string) (semreg.Decimal, error) {
	integer, fraction, _ := strings.Cut(s, ".")
	negative := strings.HasPrefix(integer, "-")
	c := strings.TrimLeft(strings.TrimPrefix(integer, "-")+fraction, "0")
	if c == "" {
		return semreg.Decimal{Coefficient: "0"}, nil
	}
	coefficient := strings.TrimRight(c, "0")
	if negative {
		coefficient = "-" + coefficient
		c = "-" + c
	}
	d := semreg.Decimal{Coefficient: coefficient, Exponent10: int32(len(c) - len(coefficient) - len(fraction))}
	return d, d.Validate()
}

func dhwDescriptors(p nativepoints.Point) roweval.Descriptors[B524Selector] {
	loc := p.Locations[0]
	p.Codec, p.NativeUnit, p.ScaleDivisor = loc.Codec, loc.NativeUnit, 1
	identity := func(x uint64, _ uint16) (uint64, error) { return x, nil }
	order := roweval.Order{Decode: identity, Encode: identity}
	return roweval.Descriptors[B524Selector]{
		Codecs: map[string]roweval.Codec{string(loc.Codec): {
			Decode: func(raw uint64, _ uint16) (semreg.Decimal, error) {
				v, err := nativepoints.DecodeNative(p, binary.LittleEndian.AppendUint32(nil, uint32(raw)))
				if err != nil {
					return semreg.Decimal{}, err
				}
				return nativeDecimal(v.Native)
			},
			Encode: func(semreg.Decimal, uint16) (uint64, error) { return 0, errors.New("no DHW write binding") },
		}},
		BitOrders: map[string]roweval.Order{"lsb0": order}, ByteOrders: map[string]roweval.Order{"little_endian": order},
		WordOrders: map[string]roweval.Order{"native": order},
	}
}

func dhwPoint(s B524Selector) nativepoints.Point {
	for _, name := range []string{"dhw_temperature_setpoint", "dhw_temperature"} {
		p := point(name)
		if pointSelector(p) == s {
			return p
		}
	}
	return nativepoints.Point{}
}

type dhwResult struct {
	Read, Gate  registerRead
	Evaluation  roweval.Evaluation
	DecodeError error
}

func evaluateDHW(table roweval.Table[B524Selector], a roweval.Acquired[B524Selector], read, gate registerRead) (dhwResult, error) {
	out := dhwResult{Read: read, Gate: gate}
	if read.Purpose != scalarRead || !read.Coherent || len(read.Value) == 0 {
		out.Evaluation = roweval.Evaluation{Reason: roweval.ReasonNoMatch, Detail: read.Reason}
		return out, nil
	}
	a.Selector, a.Payload = read.Selector, read.Value
	a.Qualifiers.Features = dhwFeatures(gate)
	captureDHWError(&out, &a)
	var err error
	out.Evaluation, err = roweval.Evaluate(table, a)
	if out.Evaluation.Reason == roweval.ReasonShortPayload {
		_, out.DecodeError = nativepoints.DecodeNative(dhwPoint(read.Selector), read.Value)
	}
	return out, err
}

func captureDHWError(out *dhwResult, a *roweval.Acquired[B524Selector]) {
	p := dhwPoint(a.Selector)
	id := string(p.Codec)
	c := a.Descriptors.Codecs[id]
	if c.Decode == nil {
		return
	}
	decode := c.Decode
	c.Decode = func(bits uint64, width uint16) (semreg.Decimal, error) {
		var value semreg.Decimal
		var err error
		if len(a.Payload) != int(width)/8 {
			_, err = nativepoints.DecodeNative(p, a.Payload)
		} else {
			value, err = decode(bits, width)
		}
		out.DecodeError = err
		return value, err
	}
	a.Descriptors.Codecs = maps.Clone(a.Descriptors.Codecs)
	a.Descriptors.Codecs[id] = c
}
