package system

import "testing"

var (
	basv0Identity = B524ControllerIdentity{EID: "BASV0", SoftwareRaw: "0217", HardwareRaw: "8803", APIVersion: "1", Revision: "1"}
	basv2Identity = B524ControllerIdentity{EID: "BASV2", SoftwareRaw: "0507", HardwareRaw: "1704"}
)

func TestLookupB524BoundDHWBivalenceIsProfileScoped(t *testing.T) {
	t.Parallel()

	got := LookupB524Bound(basv0Identity, B524OpcodeLocal, 0x00, 0x0001)
	if got.Status != B524BoundFound {
		t.Fatalf("BASV0 status = %v; want found", got.Status)
	}
	if got.Bound.Min != -20 || got.Bound.Max != 20 || got.Bound.Step != 1 {
		t.Fatalf("BASV0 bound = %+v; want -20..20 step 1", got.Bound)
	}
	if !got.Bound.Qualified || got.Bound.Codec != "EXP" || !got.Bound.AccessFlagsKnown || got.Bound.AccessFlags != 0x03 {
		t.Fatalf("BASV0 bound metadata = %+v; want qualified EXP flags 03", got.Bound)
	}

	got = LookupB524Bound(basv2Identity, B524OpcodeLocal, 0x00, 0x0001)
	if got.Status != B524BoundFound || got.Bound.Min != -20 || got.Bound.Max != 50 || got.Bound.Step != 1 {
		t.Fatalf("BASV2 result = %+v; want found -20..50 step 1", got)
	}
	if got.Bound.Qualified {
		t.Fatalf("BASV2 bound must stay unqualified: the docs retain a statement, not a fixture record")
	}
}

func TestLookupB524BoundExplicitUnknowns(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		id     B524ControllerIdentity
		opcode byte
		group  byte
		reg    uint16
		want   B524BoundStatus
	}{
		{"other software revision", B524ControllerIdentity{EID: "BASV0", SoftwareRaw: "0218", HardwareRaw: "8803", APIVersion: "1", Revision: "1"}, 0x02, 0x00, 0x0001, B524BoundUnknownProfile},
		{"other hardware", B524ControllerIdentity{EID: "BASV0", SoftwareRaw: "0217", HardwareRaw: "8804", APIVersion: "1", Revision: "1"}, 0x02, 0x00, 0x0001, B524BoundUnknownProfile},
		{"other api version", B524ControllerIdentity{EID: "BASV0", SoftwareRaw: "0217", HardwareRaw: "8803", APIVersion: "2", Revision: "1"}, 0x02, 0x00, 0x0001, B524BoundUnknownProfile},
		{"other revision", B524ControllerIdentity{EID: "BASV0", SoftwareRaw: "0217", HardwareRaw: "8803", APIVersion: "1", Revision: "2"}, 0x02, 0x00, 0x0001, B524BoundUnknownProfile},
		{"zero identity", B524ControllerIdentity{}, 0x02, 0x00, 0x0001, B524BoundUnknownProfile},
		{"remote opcode", basv0Identity, 0x06, 0x00, 0x0001, B524BoundUnknownRegister},
		{"byte-swapped register", basv0Identity, 0x02, 0x00, 0x0100, B524BoundUnknownRegister},
		{"other group", basv0Identity, 0x02, 0x01, 0x0001, B524BoundUnknownRegister},
	}
	for _, tc := range cases {
		got := LookupB524Bound(tc.id, tc.opcode, tc.group, tc.reg)
		if got.Status != tc.want {
			t.Errorf("%s: status = %v; want %v", tc.name, got.Status, tc.want)
		}
		if got.Status != B524BoundFound && got.Bound != (B524RegisterBound{}) {
			t.Errorf("%s: unknown result carries bound %+v; no fallback value is allowed", tc.name, got.Bound)
		}
	}
}
