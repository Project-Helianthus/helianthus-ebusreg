package system

import "testing"

var (
	basv0Identity = B524ControllerIdentity{EID: "BASV0", SoftwareRaw: "0217", HardwareRaw: "8803", APIVersion: "1", Revision: "1"}
)

func TestLookupB524BoundDHWBivalenceIsProfileScoped(t *testing.T) {
	t.Parallel()

	got := LookupB524Bound(basv0Identity, B524OpcodeLocal, 0x00, 0x00, 0x0001)
	if got.Status != B524BoundFound {
		t.Fatalf("BASV0 status = %v; want found", got.Status)
	}
	if got.Bound.Min != -20 || got.Bound.Max != 20 || got.Bound.Step != 1 {
		t.Fatalf("BASV0 bound = %+v; want -20..20 step 1", got.Bound)
	}
	if !got.Bound.Qualified || got.Bound.Codec != "EXP" || !got.Bound.AccessFlagsKnown || got.Bound.AccessFlags != 0x03 {
		t.Fatalf("BASV0 bound metadata = %+v; want qualified EXP flags 03", got.Bound)
	}
}

func TestLookupB524BoundExplicitUnknowns(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		id     B524ControllerIdentity
		opcode byte
		group  byte
		inst   byte
		reg    uint16
		want   B524BoundStatus
	}{
		{"other software revision", B524ControllerIdentity{EID: "BASV0", SoftwareRaw: "0218", HardwareRaw: "8803", APIVersion: "1", Revision: "1"}, 0x02, 0x00, 0x00, 0x0001, B524BoundUnknownProfile},
		{"other hardware", B524ControllerIdentity{EID: "BASV0", SoftwareRaw: "0217", HardwareRaw: "8804", APIVersion: "1", Revision: "1"}, 0x02, 0x00, 0x00, 0x0001, B524BoundUnknownProfile},
		{"other api version", B524ControllerIdentity{EID: "BASV0", SoftwareRaw: "0217", HardwareRaw: "8803", APIVersion: "2", Revision: "1"}, 0x02, 0x00, 0x00, 0x0001, B524BoundUnknownProfile},
		{"other revision", B524ControllerIdentity{EID: "BASV0", SoftwareRaw: "0217", HardwareRaw: "8803", APIVersion: "1", Revision: "2"}, 0x02, 0x00, 0x00, 0x0001, B524BoundUnknownProfile},
		{"BASV0 without API version or revision", B524ControllerIdentity{EID: "BASV0", SoftwareRaw: "0217", HardwareRaw: "8803"}, 0x02, 0x00, 0x00, 0x0001, B524BoundUnknownProfile},
		{"BASV0 without API version", B524ControllerIdentity{EID: "BASV0", SoftwareRaw: "0217", HardwareRaw: "8803", Revision: "1"}, 0x02, 0x00, 0x00, 0x0001, B524BoundUnknownProfile},
		{"BASV0 without revision", B524ControllerIdentity{EID: "BASV0", SoftwareRaw: "0217", HardwareRaw: "8803", APIVersion: "1"}, 0x02, 0x00, 0x00, 0x0001, B524BoundUnknownProfile},
		{"BASV2 identity has no seeded profile", B524ControllerIdentity{EID: "BASV2", SoftwareRaw: "0507", HardwareRaw: "1704"}, 0x02, 0x00, 0x00, 0x0001, B524BoundUnknownProfile},
		{"zero identity", B524ControllerIdentity{}, 0x02, 0x00, 0x00, 0x0001, B524BoundUnknownProfile},
		{"remote opcode", basv0Identity, 0x06, 0x00, 0x00, 0x0001, B524BoundUnknownRegister},
		{"byte-swapped register", basv0Identity, 0x02, 0x00, 0x00, 0x0100, B524BoundUnknownRegister},
		{"other instance", basv0Identity, 0x02, 0x00, 0x01, 0x0001, B524BoundUnknownRegister},
		{"other instance 0xFF", basv0Identity, 0x02, 0x00, 0xFF, 0x0001, B524BoundUnknownRegister},
		{"other group", basv0Identity, 0x02, 0x01, 0x00, 0x0001, B524BoundUnknownRegister},
	}
	for _, tc := range cases {
		got := LookupB524Bound(tc.id, tc.opcode, tc.group, tc.inst, tc.reg)
		if got.Status != tc.want {
			t.Errorf("%s: status = %v; want %v", tc.name, got.Status, tc.want)
		}
		if got.Status != B524BoundFound && got.Bound != (B524RegisterBound{}) {
			t.Errorf("%s: unknown result carries bound %+v; no fallback value is allowed", tc.name, got.Bound)
		}
	}
}

func TestLookupB524BoundOtherInstanceCarriesNoBound(t *testing.T) {
	t.Parallel()

	got := LookupB524Bound(basv0Identity, B524OpcodeLocal, 0x00, 0x01, 0x0001)
	if got.Status != B524BoundUnknownRegister {
		t.Fatalf("II01 status = %v; want unknown_register", got.Status)
	}
	if got.Bound != (B524RegisterBound{}) {
		t.Fatalf("II01 bound = %+v; want zero bound", got.Bound)
	}
}
