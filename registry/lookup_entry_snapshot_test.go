package registry

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// P9 — DeviceEntrySnapshot value-typed entry-identity snapshot API.
//
// Pre-P9 the DeviceEntry interface methods (Manufacturer, DeviceID,
// etc.) read d.info.<Field> lock-free. Concurrent Register replaces
// entry.info; readers can observe a torn read of string fields
// (string is ptr+len). LookupEntrySnapshot copies all observable
// identity fields under r.mu.RLock — race-free.

// TestLookupEntrySnapshot_ReturnsCurrentIdentity asserts the snapshot
// reflects the registered entry's current identity fields.
func TestLookupEntrySnapshot_ReturnsCurrentIdentity(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.Register(DeviceInfo{
		Address:         0x10,
		Manufacturer:    "Vaillant",
		DeviceID:        "BASV2",
		SerialNumber:    "SN-X",
		MacAddress:      "AA:BB:CC",
		SoftwareVersion: "1.2",
		HardwareVersion: "3.4",
	})

	snap, ok := reg.LookupEntrySnapshot(0x10)
	if !ok {
		t.Fatalf("LookupEntrySnapshot(0x10) ok=false")
	}
	if snap.Manufacturer != "Vaillant" {
		t.Errorf("Manufacturer = %q; want Vaillant", snap.Manufacturer)
	}
	if snap.DeviceID != "BASV2" {
		t.Errorf("DeviceID = %q; want BASV2", snap.DeviceID)
	}
	if snap.SerialNumber != "SN-X" {
		t.Errorf("SerialNumber = %q; want SN-X", snap.SerialNumber)
	}
	if snap.MacAddress != "AA:BB:CC" {
		t.Errorf("MacAddress = %q; want AA:BB:CC", snap.MacAddress)
	}
	if snap.SoftwareVersion != "1.2" {
		t.Errorf("SoftwareVersion = %q; want 1.2", snap.SoftwareVersion)
	}
	if snap.HardwareVersion != "3.4" {
		t.Errorf("HardwareVersion = %q; want 3.4", snap.HardwareVersion)
	}
	if snap.PrimaryAddress != 0x10 {
		t.Errorf("PrimaryAddress = 0x%02X; want 0x10", snap.PrimaryAddress)
	}
	if len(snap.Addresses) != 1 || snap.Addresses[0] != 0x10 {
		t.Errorf("Addresses = %v; want [0x10]", snap.Addresses)
	}
}

// TestLookupEntrySnapshot_AbsentEntryReturnsZero verifies the false
// return for an address with no entry.
func TestLookupEntrySnapshot_AbsentEntryReturnsZero(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	snap, ok := reg.LookupEntrySnapshot(0x42)
	if ok {
		t.Errorf("ok=true; want false (no entry)")
	}
	if snap.Manufacturer != "" || snap.PrimaryAddress != 0 {
		t.Errorf("snap = %+v; want zero", snap)
	}
}

// TestLookupEntrySnapshot_ReflectsLaterUpdate verifies the snapshot is
// taken at call time — a Register-level identity update visible to a
// later snapshot call.
func TestLookupEntrySnapshot_ReflectsLaterUpdate(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.Register(DeviceInfo{Address: 0x10, Manufacturer: "Initial"})

	pre, _ := reg.LookupEntrySnapshot(0x10)
	if pre.Manufacturer != "Initial" {
		t.Fatalf("pre.Manufacturer = %q; want Initial", pre.Manufacturer)
	}

	reg.Register(DeviceInfo{
		Address:      0x10,
		Manufacturer: "Updated",
		DeviceID:     "DEV",
		SerialNumber: "SN-Y",
	})

	post, _ := reg.LookupEntrySnapshot(0x10)
	if post.Manufacturer != "Updated" {
		t.Errorf("post.Manufacturer = %q; want Updated", post.Manufacturer)
	}
	if post.DeviceID != "DEV" {
		t.Errorf("post.DeviceID = %q; want DEV", post.DeviceID)
	}
}

// TestLookupEntrySnapshot_SliceCopyIsolatedFromRegistry verifies the
// snapshot's Addresses slice is a copy — mutating it must NOT affect
// registry state. Uses AliasAddresses to plant a multi-address
// entry deterministically.
func TestLookupEntrySnapshot_SliceCopyIsolatedFromRegistry(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.Register(DeviceInfo{Address: 0x10, SerialNumber: "SN-X"})
	reg.Register(DeviceInfo{Address: 0x15, SerialNumber: "SN-Y"})
	if err := reg.AliasAddresses(0x10, 0x15); err != nil {
		t.Fatalf("AliasAddresses(0x10, 0x15) error = %v", err)
	}

	snap, _ := reg.LookupEntrySnapshot(0x10)
	if len(snap.Addresses) < 2 {
		t.Fatalf("Addresses = %v; want both 0x10 and 0x15 (alias-merge)", snap.Addresses)
	}

	// Record the original byte at index 0 so we can prove the registry
	// didn't change after mutating the snapshot's slice.
	preMutationAddr0 := snap.Addresses[0]

	// Mutate the snapshot's slice.
	snap.Addresses[0] = 0xFF

	// Re-fetch the snapshot from the registry — it must still have
	// the original addresses.
	snap2, _ := reg.LookupEntrySnapshot(0x10)
	if len(snap2.Addresses) == 0 || snap2.Addresses[0] != preMutationAddr0 {
		t.Errorf("snap2.Addresses[0] = 0x%02X; want 0x%02X (mutation leaked through to registry storage)", snap2.Addresses[0], preMutationAddr0)
	}
	for _, a := range snap2.Addresses {
		if a == 0xFF {
			t.Errorf("snap2.Addresses contains 0xFF; mutation leaked through")
		}
	}
}

// TestIterateSnapshots_VisitsAllEntries asserts the snapshot iteration
// visits every registered entry.
func TestIterateSnapshots_VisitsAllEntries(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.Register(DeviceInfo{Address: 0x10, Manufacturer: "M1"})
	reg.Register(DeviceInfo{Address: 0x26, Manufacturer: "M2"})

	visited := make(map[byte]string)
	reg.IterateSnapshots(func(snap DeviceEntrySnapshot) bool {
		visited[snap.PrimaryAddress] = snap.Manufacturer
		return true
	})

	if len(visited) != 2 {
		t.Errorf("visited %d entries; want 2 (got %v)", len(visited), visited)
	}
	if visited[0x10] != "M1" {
		t.Errorf("visited[0x10] = %q; want M1", visited[0x10])
	}
	if visited[0x26] != "M2" {
		t.Errorf("visited[0x26] = %q; want M2", visited[0x26])
	}
}

// TestIterateSnapshots_StopOnFalse verifies the early-stop semantic.
func TestIterateSnapshots_StopOnFalse(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	reg.Register(DeviceInfo{Address: 0x10})
	reg.Register(DeviceInfo{Address: 0x20})
	reg.Register(DeviceInfo{Address: 0x30})

	count := 0
	reg.IterateSnapshots(func(snap DeviceEntrySnapshot) bool {
		count++
		return false // stop after first
	})
	if count != 1 {
		t.Errorf("callback invocations = %d; want 1 (early stop)", count)
	}
}

// TestDeviceEntrySnapshot_AddressByRole_FacesPath verifies the
// snapshot's AddressByRole resolves via the Faces slice (mirrors the
// live deviceEntry.AddressByRole behaviour).
func TestDeviceEntrySnapshot_AddressByRole_FacesPath(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	now := time.Now()
	reg.RegisterPassiveObserved(DeviceInfo{Address: 0x10}, SlotRoleMaster, now)
	reg.RegisterPassiveObserved(DeviceInfo{Address: 0x15, SerialNumber: "SN-X"}, SlotRoleSlave, now)
	reg.Register(DeviceInfo{Address: 0x10, SerialNumber: "SN-X"})

	snap, ok := reg.LookupEntrySnapshot(0x10)
	if !ok {
		t.Fatalf("LookupEntrySnapshot(0x10) ok=false")
	}

	masterAddr, masterOK := snap.AddressByRole(SlotRoleMaster)
	if !masterOK {
		t.Errorf("AddressByRole(SlotRoleMaster) ok=false; want true (entry has 0x10 master face)")
	}
	if masterAddr != 0x10 {
		t.Errorf("master address = 0x%02X; want 0x10", masterAddr)
	}
}

// TestLookupEntrySnapshot_RaceFreeUnderConcurrentRegisters is the
// central P9 contract: the snapshot's identity fields are stable for
// the caller's lifetime regardless of concurrent Register calls.
//
// Methodology: a writer goroutine repeatedly Register(s) the same
// address with rotating Manufacturer / DeviceID strings. The reader
// snapshots and reads the local copy's fields. Race detector flags
// any unsynced access via the snapshot path.
//
// Pre-P9 a reader doing entry.Manufacturer() on a live *deviceEntry
// pointer would race with the writer's `entry.info = storedInfo` —
// the race detector would trip. The snapshot copy under RLock
// eliminates that surface.
func TestLookupEntrySnapshot_RaceFreeUnderConcurrentRegisters(t *testing.T) {
	t.Parallel()

	reg := NewDeviceRegistry(nil)
	addr := byte(0x10)
	reg.Register(DeviceInfo{Address: addr, Manufacturer: "Initial"})

	stop := make(chan struct{})
	var writes uint64
	var writerWg sync.WaitGroup
	var readerWg sync.WaitGroup

	writerWg.Add(1)
	go func() {
		defer writerWg.Done()
		i := 0
		for {
			select {
			case <-stop:
				return
			default:
			}
			reg.Register(DeviceInfo{
				Address:      addr,
				Manufacturer: "Manufacturer-" + string(rune('A'+(i%26))),
				DeviceID:     "Device-" + string(rune('A'+(i%26))),
				SerialNumber: "Serial-" + string(rune('A'+(i%26))),
			})
			atomic.AddUint64(&writes, 1)
			i++
		}
	}()

	// Start barrier — wait for at least one Register before reads.
	deadline := time.Now().Add(2 * time.Second)
	for atomic.LoadUint64(&writes) == 0 {
		if time.Now().After(deadline) {
			close(stop)
			writerWg.Wait()
			t.Fatal("writer made no Register calls within 2s")
		}
		time.Sleep(time.Microsecond)
	}

	startWrites := atomic.LoadUint64(&writes)

	// 4 readers using LookupEntrySnapshot.
	const numReaders = 4
	for r := 0; r < numReaders; r++ {
		readerWg.Add(1)
		go func() {
			defer readerWg.Done()
			for j := 0; j < 1000; j++ {
				snap, ok := reg.LookupEntrySnapshot(addr)
				if !ok {
					continue
				}
				_ = snap.Manufacturer
				_ = snap.DeviceID
				_ = snap.SerialNumber
				_ = snap.PrimaryAddress
				_ = snap.Addresses
				_ = snap.Faces
			}
		}()
	}

	// Wait for readers FIRST so the writer keeps producing Registers
	// throughout the reader phase.
	readerWg.Wait()
	endWrites := atomic.LoadUint64(&writes)

	close(stop)
	writerWg.Wait()

	if endWrites-startWrites < 1 {
		t.Errorf("writer Register count advanced by 0 during reader phase; want >= 1 (overlap not exercised)")
	}
}
