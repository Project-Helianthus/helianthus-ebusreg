package system

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	ebuserrors "github.com/d3vi1/helianthus-ebusgo/errors"
)

func TestToByteNumericCoercion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
		want  byte
		ok    bool
	}{
		{name: "float64 integer", value: float64(17), want: 17, ok: true},
		{name: "float64 fractional", value: float64(17.5), ok: false},
		{name: "json number integer", value: json.Number("18"), want: 18, ok: true},
		{name: "json number float integer", value: json.Number("19.0"), want: 19, ok: true},
		{name: "json number float fractional", value: json.Number("19.5"), ok: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, ok := toByte(test.value)
			if ok != test.ok {
				t.Fatalf("ok=%t, want %t", ok, test.ok)
			}
			if ok && got != test.want {
				t.Fatalf("value=%d, want %d", got, test.want)
			}
		})
	}
}

func TestRegisterAddrParamNumericCoercion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   any
		want    uint16
		wantErr bool
	}{
		{name: "float64 integer", value: float64(0x030D), want: 0x030D},
		{name: "json number integer", value: json.Number("781"), want: 781},
		{name: "json number float integer", value: json.Number("781.0"), want: 781},
		{name: "json number fractional", value: json.Number("781.5"), wantErr: true},
		{name: "float64 out of range", value: float64(70000), wantErr: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, ok, err := registerAddrParam(map[string]any{"addr": test.value}, "addr")
			if !ok {
				t.Fatalf("expected present value")
			}
			if test.wantErr {
				if !errors.Is(err, ebuserrors.ErrInvalidPayload) {
					t.Fatalf("expected ErrInvalidPayload, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != test.want {
				t.Fatalf("value=0x%04X, want 0x%04X", got, test.want)
			}
		})
	}
}

func TestRegisterWriteTemplateBuildNumericGraphQLShapes(t *testing.T) {
	t.Parallel()

	template := registerWriteTemplate{primary: 0xB5, secondary: 0x09}
	payload, err := template.Build(map[string]any{
		"addr": float64(0x030D),
		"data": []any{float64(1), json.Number("2"), json.Number("3.0")},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []byte{registerCmdWrite, 0x03, 0x0D, 0x01, 0x02, 0x03}
	if !reflect.DeepEqual(payload, want) {
		t.Fatalf("payload=% X, want % X", payload, want)
	}
}

func TestRegisterWriteTemplateBuildRejectsFractionalData(t *testing.T) {
	t.Parallel()

	template := registerWriteTemplate{primary: 0xB5, secondary: 0x09}
	_, err := template.Build(map[string]any{
		"addr": 0x030D,
		"data": []any{float64(1.5)},
	})
	if !errors.Is(err, ebuserrors.ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}
}
