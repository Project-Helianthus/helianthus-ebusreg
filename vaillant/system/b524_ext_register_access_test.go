package system

import "testing"

func TestDecodeExtRegisterResponseMarksStaticConstraintsUnscoped(t *testing.T) {
	t.Parallel()

	mapped := decodeExtRegisterResponse(0x02, B524OpcodeLocal, 0x01, 0x01, 0x0004, []byte{0x01, 0x01, 0x04, 0x00, 0x2A})
	if _, ok := mapped["constraints"]; !ok {
		t.Fatal("constraints absent for a catalog entry; want present")
	}
	if got := mapped["constraint_qualified"]; !got.Valid || got.Value != false {
		t.Fatalf("constraint_qualified = %+v; want valid false", got)
	}
	if got := mapped["constraint_scope"]; !got.Valid || got.Value != "unscoped_catalog" {
		t.Fatalf("constraint_scope = %+v; want valid unscoped_catalog", got)
	}
	// Existing keys stay.
	for _, key := range []string{"constraint_record", "constraint_record_hex", "constraint_type", "constraint_min", "constraint_max", "constraint_step"} {
		if got := mapped[key]; !got.Valid {
			t.Errorf("%s = %+v; want valid (additive change must keep existing keys)", key, got)
		}
	}

	unmapped := decodeExtRegisterResponse(0x02, B524OpcodeLocal, 0x00, 0x00, 0x5C00, []byte{0x00})
	for _, key := range []string{"constraints", "constraint_qualified", "constraint_scope"} {
		if got, ok := unmapped[key]; ok {
			t.Errorf("%s = %+v; want absent without a catalog entry", key, got)
		}
	}
}
