package ebus_standard_catalog

import (
	"context"
	"errors"
	"testing"
)

// TestInvoke_SafetyClassGate exercises the invoke-boundary gate for every
// safety class the catalog can carry. Allowed classes must NOT return
// ErrSafetyClassDenied; denied classes MUST return it exactly.
//
// The test builds a synthetic single-command catalog per class so coverage
// does not depend on which classes happen to appear in the real YAML.
func TestInvoke_SafetyClassGate(t *testing.T) {
	cases := []struct {
		class  SafetyClass
		denied bool
	}{
		{SafetyReadOnlySafe, false},
		{SafetyReadOnlyBusLoad, false},
		{SafetyMutating, true},
		{SafetyDestructive, true},
		{SafetyBroadcast, true},
		{SafetyMemoryWrite, true},
	}

	for _, tc := range cases {
		t.Run(string(tc.class), func(t *testing.T) {
			cat := syntheticCatalogForClass(tc.class)
			p, err := NewProvider(cat, true)
			if err != nil {
				t.Fatalf("NewProvider: %v", err)
			}
			_, err = p.Invoke(context.Background(), "ebus_standard.test.method", nil, CallerContextUserFacing)
			if tc.denied {
				if !errors.Is(err, ErrSafetyClassDenied) {
					t.Fatalf("class=%s err=%v, want ErrSafetyClassDenied", tc.class, err)
				}
			} else {
				if errors.Is(err, ErrSafetyClassDenied) {
					t.Fatalf("class=%s err=%v, want NOT ErrSafetyClassDenied", tc.class, err)
				}
			}
		})
	}
}

// TestInvoke_SystemNMRuntimeCallerStillDenied confirms that in M3 the
// system_nm_runtime caller context is carried but does NOT whitelist any
// mutating class. This contract is reserved for M4+ expansion; M3 must not
// leak an earlier-than-planned escalation.
func TestInvoke_SystemNMRuntimeCallerStillDenied(t *testing.T) {
	cat := syntheticCatalogForClass(SafetyMutating)
	p, err := NewProvider(cat, true)
	if err != nil {
		t.Fatalf("NewProvider: %v", err)
	}
	_, err = p.Invoke(context.Background(), "ebus_standard.test.method", nil, CallerContextSystemNMRuntime)
	if !errors.Is(err, ErrSafetyClassDenied) {
		t.Fatalf("err=%v, want ErrSafetyClassDenied for system_nm_runtime in M3", err)
	}
}

// syntheticCatalogForClass builds a one-service one-command catalog whose
// single command carries the supplied safety_class. Its identity key is
// complete and unique so LoadCatalog (indirectly, via the provider) has no
// reason to reject it.
func syntheticCatalogForClass(class SafetyClass) Catalog {
	pb := uint8(0x03)
	sb := uint8(0xF0)
	return Catalog{
		Namespace:  Namespace,
		Version:    CatalogVersion,
		PlanSHA256: CanonicalPlanSHA256,
		Services: []Service{
			{
				PB:   &pb,
				Name: "synthetic",
				Commands: []Command{
					{
						ID:   "ebus_standard.test.method",
						Name: "synthetic",
						Identity: IdentityKey{
							Namespace:                       Namespace,
							PB:                              &pb,
							SB:                              &sb,
							SelectorPath:                    "",
							TelegramClass:                   TelegramClassAddressed,
							Direction:                       DirectionRequest,
							RequestOrResponseRole:           RoleInitiator,
							BroadcastOrAddressed:            AddressedDirect,
							AnswerPolicy:                    AnswerRequired,
							LengthPrefixMode:                LengthPrefixNone,
							SelectorDecoder:                 "none",
							ServiceVariant:                  "synthetic",
							TransportCapabilityRequirements: []string{"master_slave"},
							Version:                         CatalogVersion,
						},
						SafetyClass: class,
					},
				},
			},
		},
	}
}
