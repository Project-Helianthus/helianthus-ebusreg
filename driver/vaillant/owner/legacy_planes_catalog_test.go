package owner

import (
	"encoding/json"
	"testing"

	catalogvaillant "github.com/Project-Helianthus/helianthus-ebusreg/catalog/vaillant"
)

func assertLegacyRoleMutations(t *testing.T) {
	t.Helper()
	raw, err := catalogvaillant.B504CircuitRoleEntries()
	if err != nil {
		t.Fatal(err)
	}
	var circuitRows []struct {
		Key catalogvaillant.B504CircuitRoleKey `json:"key"`
	}
	if err := json.Unmarshal(raw, &circuitRows); err != nil || len(circuitRows) == 0 {
		t.Fatal(err)
	}
	circuitKey := circuitRows[0].Key
	if found, err := lookupCircuitRole(circuitKey); err != nil || !found {
		t.Fatalf("unmutated circuit role found=%v err=%v", found, err)
	}
	for _, mutation := range []struct {
		field  string
		mutate func(*catalogvaillant.B504CircuitRoleKey)
	}{
		{"family", func(key *catalogvaillant.B504CircuitRoleKey) { key.B504CircuitRoleFamily = "absent-family" }},
		{"variant", func(key *catalogvaillant.B504CircuitRoleKey) { key.B504CircuitRoleVariant = "absent-variant" }},
		{"primary", func(key *catalogvaillant.B504CircuitRoleKey) { key.B504CircuitRolePrimary++ }},
		{"secondary", func(key *catalogvaillant.B504CircuitRoleKey) {
			key.B504CircuitRoleSecondary = mutateOptionalOctet(key.B504CircuitRoleSecondary)
		}},
		{"operation", func(key *catalogvaillant.B504CircuitRoleKey) {
			key.B504CircuitRoleOperation = mutateOptionalOctet(key.B504CircuitRoleOperation)
		}},
		{"direction", func(key *catalogvaillant.B504CircuitRoleKey) { key.B504CircuitRoleDirection = "sideways" }},
		{"field", func(key *catalogvaillant.B504CircuitRoleKey) { key.B504CircuitRoleField = "absent-field" }},
		{"alternative", func(key *catalogvaillant.B504CircuitRoleKey) { key.B504CircuitRoleAlternative = "absent-alternative" }},
		{"group", func(key *catalogvaillant.B504CircuitRoleKey) {
			key.B504CircuitRoleGroup = mutateOptionalOctet(key.B504CircuitRoleGroup)
		}},
		{"instance_scope", func(key *catalogvaillant.B504CircuitRoleKey) { key.B504CircuitRoleInstanceScope = "absent-scope" }},
		{"register", func(key *catalogvaillant.B504CircuitRoleKey) {
			key.B504CircuitRoleRegister = mutateOptionalWide(key.B504CircuitRoleRegister)
		}},
	} {
		mutated := cloneJSON(t, circuitKey)
		mutation.mutate(&mutated)
		if found, err := lookupCircuitRole(mutated); err != nil || found {
			t.Fatalf("%s mutation still matched; dropping %s from equality would keep the original row", mutation.field, mutation.field)
		}
	}
	roleRaw, err := catalogvaillant.B505OperatingRoleEntries()
	if err != nil {
		t.Fatal(err)
	}
	var roleRows []struct {
		Identity catalogvaillant.B505OperatingRoleKey `json:"identity"`
	}
	if err := json.Unmarshal(roleRaw, &roleRows); err != nil || len(roleRows) == 0 {
		t.Fatal(err)
	}
	roleKey := roleRows[0].Identity
	if found, err := lookupOperatingRole(roleKey); err != nil || !found {
		t.Fatalf("unmutated operating role found=%v err=%v", found, err)
	}
	for _, mutation := range []struct {
		field  string
		mutate func(*catalogvaillant.B505OperatingRoleKey)
	}{
		{"family", func(key *catalogvaillant.B505OperatingRoleKey) { key.B505OperatingRoleFamily = "absent-family" }},
		{"variant", func(key *catalogvaillant.B505OperatingRoleKey) { key.B505OperatingRoleVariant = "absent-variant" }},
		{"primary", func(key *catalogvaillant.B505OperatingRoleKey) { key.B505OperatingRolePrimary++ }},
		{"secondary", func(key *catalogvaillant.B505OperatingRoleKey) {
			key.B505OperatingRoleSecondary = mutateOptionalOctet(key.B505OperatingRoleSecondary)
		}},
		{"operation", func(key *catalogvaillant.B505OperatingRoleKey) {
			key.B505OperatingRoleOperation = mutateOptionalOctet(key.B505OperatingRoleOperation)
		}},
		{"direction", func(key *catalogvaillant.B505OperatingRoleKey) { key.B505OperatingRoleDirection = "sideways" }},
		{"field", func(key *catalogvaillant.B505OperatingRoleKey) { key.B505OperatingRoleField = "absent-field" }},
		{"alternative", func(key *catalogvaillant.B505OperatingRoleKey) {
			key.B505OperatingRoleAlternative = "absent-alternative"
		}},
		{"group", func(key *catalogvaillant.B505OperatingRoleKey) {
			key.B505OperatingRoleGroup = mutateOptionalOctet(key.B505OperatingRoleGroup)
		}},
		{"instance_scope", func(key *catalogvaillant.B505OperatingRoleKey) {
			key.B505OperatingRoleInstanceScope = "absent-scope"
		}},
		{"register", func(key *catalogvaillant.B505OperatingRoleKey) {
			key.B505OperatingRoleRegister = mutateOptionalWide(key.B505OperatingRoleRegister)
		}},
	} {
		mutated := cloneJSON(t, roleKey)
		mutation.mutate(&mutated)
		if found, err := lookupOperatingRole(mutated); err != nil || found {
			t.Fatalf("%s mutation still matched; dropping %s from equality would keep the original row", mutation.field, mutation.field)
		}
	}
	sequenceRaw, err := catalogvaillant.B501LegacySequenceEntries()
	if err != nil {
		t.Fatal(err)
	}
	var sequenceRows []struct {
		Identity catalogvaillant.B501LegacySequenceKey `json:"identity"`
	}
	if err := json.Unmarshal(sequenceRaw, &sequenceRows); err != nil || len(sequenceRows) == 0 {
		t.Fatal(err)
	}
	sequenceKey := sequenceRows[0].Identity
	if found, err := lookupLegacySequence(sequenceKey); err != nil || !found {
		t.Fatalf("unmutated sequence found=%v err=%v", found, err)
	}
	for _, mutation := range []struct {
		field  string
		mutate func(*catalogvaillant.B501LegacySequenceKey)
	}{
		{"family", func(key *catalogvaillant.B501LegacySequenceKey) { key.B501LegacySequenceFamily = "absent-family" }},
		{"variant", func(key *catalogvaillant.B501LegacySequenceKey) { key.B501LegacySequenceVariant = "absent-variant" }},
		{"primary", func(key *catalogvaillant.B501LegacySequenceKey) { key.B501LegacySequencePrimary++ }},
		{"secondary", func(key *catalogvaillant.B501LegacySequenceKey) {
			key.B501LegacySequenceSecondary = mutateOptionalOctet(key.B501LegacySequenceSecondary)
		}},
		{"operation", func(key *catalogvaillant.B501LegacySequenceKey) {
			key.B501LegacySequenceOperation = mutateOptionalOctet(key.B501LegacySequenceOperation)
		}},
		{"direction", func(key *catalogvaillant.B501LegacySequenceKey) { key.B501LegacySequenceDirection = "sideways" }},
		{"field", func(key *catalogvaillant.B501LegacySequenceKey) { key.B501LegacySequenceField = "absent-field" }},
		{"alternative", func(key *catalogvaillant.B501LegacySequenceKey) {
			key.B501LegacySequenceAlternative = "absent-alternative"
		}},
		{"group", func(key *catalogvaillant.B501LegacySequenceKey) {
			key.B501LegacySequenceGroup = mutateOptionalOctet(key.B501LegacySequenceGroup)
		}},
		{"instance_scope", func(key *catalogvaillant.B501LegacySequenceKey) {
			key.B501LegacySequenceInstanceScope = "absent-scope"
		}},
		{"register", func(key *catalogvaillant.B501LegacySequenceKey) {
			key.B501LegacySequenceRegister = mutateOptionalWide(key.B501LegacySequenceRegister)
		}},
	} {
		mutated := cloneJSON(t, sequenceKey)
		mutation.mutate(&mutated)
		if found, err := lookupLegacySequence(mutated); err != nil || found {
			t.Fatalf("%s mutation still matched; dropping %s from equality would keep the original row", mutation.field, mutation.field)
		}
	}
}
