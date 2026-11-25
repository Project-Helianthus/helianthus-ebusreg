package ebus

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/Project-Helianthus/helianthus-ebusgo/engine"
	ebuserrors "github.com/Project-Helianthus/helianthus-ebusgo/errors"
	"github.com/Project-Helianthus/helianthus-ebusgo/protocol"
	"github.com/Project-Helianthus/helianthus-ebusgo/transport"
)

const fixtureSource, fixtureTarget = byte(0x10), byte(0x08)

type simulatedAdapter struct {
	*transport.Loopback
	inbound []byte
	echoes  int
	written []byte
	// inFlight runs once, when the first reply byte is delivered: the send has begun and has not completed.
	inFlight func()
}

func (a *simulatedAdapter) Write(data []byte) (int, error) {
	a.written = append(a.written, data...)
	a.echoes += len(data)
	return a.Loopback.Write(data)
}

func (a *simulatedAdapter) ReadByte() (byte, error) {
	if a.echoes > 0 {
		a.echoes--
		return a.Loopback.ReadByte()
	}
	if len(a.inbound) == 0 {
		return 0, ebuserrors.ErrTimeout
	}
	if f := a.inFlight; f != nil {
		a.inFlight = nil
		f()
	}
	b := a.inbound[0]
	a.inbound = a.inbound[1:]
	return b, nil
}

func fixtureJSON[T any](t *testing.T, path string) T {
	t.Helper()
	var v T
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return v
}

func octets(t *testing.T, text string) []byte {
	t.Helper()
	b, err := hex.DecodeString(text)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func fixtureReply(t *testing.T, frame protocol.Frame, wantRequest, raw []byte, failure string) ([]byte, error) {
	t.Helper()
	reply, err := simulateExchange(t, frame, fixtureSource, fixtureTarget, wantRequest, raw, failure)
	if reply == nil {
		return nil, err
	}
	return reply.Data, err
}

func simulateExchange(t *testing.T, frame protocol.Frame, source, target byte, wantRequest, raw []byte, failure string) (*protocol.Frame, error) {
	t.Helper()
	return simulateExchangeInFlight(t, frame, source, target, wantRequest, raw, failure, nil)
}

func simulateExchangeInFlight(t *testing.T, frame protocol.Frame, source, target byte, wantRequest, raw []byte, failure string, inFlight func()) (*protocol.Frame, error) {
	t.Helper()
	adapter := &simulatedAdapter{Loopback: transport.NewLoopback(), inFlight: inFlight}
	var wantErr error
	switch failure {
	case "nack":
		adapter.inbound = []byte{protocol.SymbolNack, protocol.SymbolNack}
		wantErr = ebuserrors.ErrNACK
	case "timeout":
		wantErr = ebuserrors.ErrTimeout
	case "":
		reply, err := protocol.EncodeSlaveResponse(raw)
		if err != nil {
			t.Fatal(err)
		}
		adapter.inbound = append([]byte{protocol.SymbolAck}, reply...)
	default:
		t.Fatalf("unknown synthetic failure %q", failure)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	bus := protocol.NewBus(adapter, protocol.BusConfig{}, 1)
	bus.Run(ctx)
	var response *protocol.Frame
	got, err := engine.NewSemanticReadScheduler().Get(ctx, "fixture", 0, func(ctx context.Context) ([]byte, error) {
		var err error
		response, err = bus.Send(ctx, frame)
		if response == nil {
			return nil, err
		}
		return response.Data, err
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("adapter: got error %v, want %v", err, wantErr)
	}
	if !bytes.Equal(got, raw) {
		t.Fatalf("adapter delivered %x, want %x", got, raw)
	}
	if (response != nil) != (failure == "") {
		t.Fatalf("adapter frame presence was %t for %q", response != nil, failure)
	}
	assertRequest(t, adapter.written, source, target, wantRequest, failure)
	return response, err
}

func assertRequest(t *testing.T, written []byte, source, target byte, request []byte, failure string) {
	t.Helper()
	body := []byte{source, target, request[0], request[1], byte(len(request) - 2)}
	body = append(body, request[2:]...)
	wire := protocol.EscapeBytes(append(body, protocol.CRC(body)))
	switch failure {
	case "":
		wire = append(wire, protocol.SymbolAck)
	case "nack":
		wire = append(wire, wire...)
	}
	if !bytes.HasPrefix(written, wire) {
		t.Fatalf("writes %x did not start with %x", written, wire)
	}
	if rest := written[len(wire):]; len(rest) != 0 && !bytes.Equal(rest, []byte{protocol.SymbolSyn}) {
		t.Fatalf("unexpected extra writes %x", rest)
	}
}
