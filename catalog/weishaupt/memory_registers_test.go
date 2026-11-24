package weishaupt

import (
	"bytes"
	"errors"
	"testing"
)

// The page's example IDs round-trip through the read request and the write layout.
func TestKromschroederObservedIDsBuildTheDocumentedRequests(t *testing.T) {
	id, ok := ObservedMemoryID(0x0004)
	if !ok || !bytes.Equal(id[:], []byte{0x58, 0x01, 0x04}) {
		t.Fatalf("observed ID for 0x0004 = % X", id)
	}
	if got := EncodeMemoryRead(id); !bytes.Equal(got, []byte{0x58, 0x01, 0x04}) {
		t.Fatalf("5000 request = % X", got)
	}
	// 5001 write: the ID then a one-byte and a two-byte value.
	got, err := EncodeMemoryWrite(id, []byte{0x2A})
	if err != nil || !bytes.Equal(got, []byte{0x58, 0x01, 0x04, 0x2A}) {
		t.Fatalf("5001 one-byte write = % X, err %v", got, err)
	}
	got, err = EncodeMemoryWrite(id, []byte{0x34, 0x12})
	if err != nil || !bytes.Equal(got, []byte{0x58, 0x01, 0x04, 0x34, 0x12}) {
		t.Fatalf("5001 two-byte write = % X, err %v", got, err)
	}
	if _, err := EncodeMemoryWrite(id, nil); !errors.Is(err, ErrValueSize) {
		t.Fatalf("empty write err = %v; want ErrValueSize", err)
	}
	if _, err := EncodeMemoryWrite(id, []byte{1, 2, 3}); !errors.Is(err, ErrValueSize) {
		t.Fatalf("three-byte write err = %v; want ErrValueSize", err)
	}
	if _, ok := ObservedMemoryID(0x0006); ok {
		t.Fatal("an address the page shows no ID for returned an ID")
	}
}

func TestKromschroederReplyKeepsRawAndChecksWidth(t *testing.T) {
	f, err := DecodeMemoryReply([]byte{0x7F})
	if err != nil || f.Name != "value_raw" || !bytes.Equal(f.Value.([]byte), []byte{0x7F}) {
		t.Fatalf("one-byte reply = %#v, err %v", f, err)
	}
	if _, err := DecodeMemoryReply([]byte{1, 2, 3}); !errors.Is(err, ErrWidth) {
		t.Fatalf("three-byte reply err = %v; want ErrWidth", err)
	}
}

func TestKromschroederSelectedRegisters(t *testing.T) {
	reg, ok := MemoryRegisterByAddress(0x0068)
	if !ok || reg.Name != "TVIST" || reg.Description != "Flow temperature instant" {
		t.Fatalf("0x0068 = %#v", reg)
	}
	if len(memoryRegisterTable) != 33 {
		t.Fatalf("selected registers = %d; the page lists 33", len(memoryRegisterTable))
	}
}
