#!/usr/bin/env python3
"""Reproduce the removal of work-item codes from catalog file and Go names.

The work runs in slices; each slice names one package, its file moves and
the stem every code prefix takes, and is run once against the tree that
still holds its coded files (--slice selects it):

- ebusstandard: every coded file of catalog/ebusstandard.
- vaillant-01-29: the catalog/vaillant files whose code number is 01 to 29.
  Files with higher codes that use these identifiers are edited in place.

Production and test Go files lose their code prefix and keep the behavior part
of their name. Go identifiers that start with a code prefix of the slice are
renamed to a stem named after the destination file. The rewrite is
token-aware: string literals (JSON tags, values, error texts, fixture paths)
and embedded data and testdata file names stay byte-identical, so no
externally visible string changes here. Comments only change where they name
a renamed identifier; go: directives never change. Test functions
(TestAb...) keep their names. Exported struct fields whose name is renamed
and that carry no JSON tag get an explicit tag with the unchanged key, so
serialized output and fixture decoding stay the same.

The script stops before writing anything when a precondition or a collision
check fails, and prints the identifier map with --print-map.
"""

import argparse
import importlib.util
import json
import os
from pathlib import Path
import re
import subprocess
import sys
from types import SimpleNamespace

MODULE = "github.com/Project-Helianthus/helianthus-ebusreg/"

PATH_KEYED = (
    "naming-inventory.json",
    "naming-consumers.json",
    "scripts/behavior-name-baseline.json",
    "scripts/go-lint-baseline.json",
    "scripts/size-baseline.json",
)

CODE = re.compile(r"([aA]b[0-9]{2}(?:[aA][a-z][0-9]{2})?)([A-Za-z0-9_]*)")
IDENT = re.compile(r"[A-Za-z_][A-Za-z0-9_]*")
WORD = re.compile(r"[A-Za-z0-9_]+")
NUMBER = re.compile(r"[0-9][0-9A-Za-z_.]*")
KEYWORDS = frozenset(
    "break case chan const continue default defer else fallthrough for func go "
    "goto if import interface map package range return select struct switch "
    "type var".split())


# Keep the regression coverage in the reproducible rewrite as well as the
# resulting package. Its explicit keys pin the serialization contract.
PROJECTION_JSON_KEYS_TEST = r'''func TestRenamedProjectionJSONKeys(t *testing.T) {
	cases := []struct {
		value any
		keys  map[string]string
	}{
		{Standard09BodyViewKey{}, map[string]string{"Standard09BodyViewFamily": "Ab25An25Family", "Standard09BodyViewProfile": "Ab25An25Profile", "Standard09BodyViewService": "Ab25An25Service", "Standard09BodyViewDirection": "Ab25An25Direction"}},
		{Standard09BodyViewState{}, map[string]string{"Standard09BodyViewKey": "Ab25An25Key", "Standard09BodyViewRaw": "Ab25An25Raw", "Standard09BodyViewState": "Ab25An25State", "Standard09BodyViewContractGrade": "Ab25An25ContractGrade", "Standard09BodyViewMeaningGrade": "Ab25An25MeaningGrade", "Standard09BodyViewExecutionGrade": "Ab25An25ExecutionGrade", "Standard09BodyViewPhysicalSupportGrade": "Ab25An25PhysicalSupportGrade", "Standard09BodyViewPhysicalOffsetsGrade": "Ab25An25PhysicalOffsetsGrade", "Standard09BodyViewPositivePrefixRelevanceGrade": "Ab25An25PositivePrefixRelevanceGrade"}},
		{Standard09BodyViewBody{}, map[string]string{"Standard09BodyViewState": "Ab25An25State", "Standard09BodyViewData": "Ab25An25Data", "Standard09BodyViewUnresolvedTail": "Ab25An25UnresolvedTail"}},
		{Standard09BodyViewAddress{}, map[string]string{"Standard09BodyViewState": "Ab25An25State", "Standard09BodyViewLE": "Ab25An25LE", "Standard09BodyViewLow": "Ab25An25Low", "Standard09BodyViewHigh": "Ab25An25High"}},
		{Standard09BodyViewCount{}, map[string]string{"Standard09BodyViewState": "Ab25An25State", "Standard09BodyViewValue": "Ab25An25Value", "Standard09BodyViewCandidateMaximum": "Ab25An25CandidateMaximum", "Standard09BodyViewConflict": "Ab25An25Conflict"}},
		{Standard05CandidateViewKey{}, map[string]string{"Standard05CandidateViewFamily": "Ab78Family", "Standard05CandidateViewVariant": "Ab78Variant", "Standard05CandidateViewPrimary": "Ab78Primary", "Standard05CandidateViewSecondary": "Ab78Secondary", "Standard05CandidateViewOperation": "Ab78Operation", "Standard05CandidateViewDirection": "Ab78Direction", "Standard05CandidateViewSelector": "Ab78Selector"}},
		{Standard05CandidateViewField{}, map[string]string{"Standard05CandidateViewLabel": "Ab78Label", "Standard05CandidateViewRaw": "Ab78Raw", "Standard05CandidateViewValue": "Ab78Value", "Standard05CandidateViewName": "Ab78Name", "Standard05CandidateViewInterpretation": "Ab78Interpretation", "Standard05CandidateViewDomain": "Ab78Domain", "Standard05CandidateViewDefault": "Ab78Default", "Standard05CandidateViewMinimum": "Ab78Minimum", "Standard05CandidateViewMaximum": "Ab78Maximum"}},
		{Standard05CandidateView{}, map[string]string{"Standard05CandidateViewKey": "Ab78Key", "Standard05CandidateViewRaw": "Ab78Raw", "Standard05CandidateViewRemainder": "Ab78Remainder", "Standard05CandidateViewState": "Ab78State", "Standard05CandidateViewComplete": "Ab78Complete", "Standard05CandidateViewFields": "Ab78Fields", "Standard05CandidateViewAPIGrade": "Ab78APIGrade", "Standard05CandidateViewMeaningGrade": "Ab78MeaningGrade", "Standard05CandidateViewNativeGrade": "Ab78NativeGrade", "Standard05CandidateViewUnknowns": "Ab78Unknowns"}},
		{Standard03CandidateViewSelection{}, map[string]string{"Standard03CandidateViewFamily": "Ab76Family", "Standard03CandidateViewVariant": "Ab76Variant", "Standard03CandidateViewFunction": "Ab76Function", "Standard03CandidateViewOperation": "Ab76Operation", "Standard03CandidateViewDirection": "Ab76Direction", "Standard03CandidateViewSelector": "Ab76Selector", "Standard03CandidateViewEncoding": "Ab76Encoding"}},
		{Standard03CandidateViewComponent{}, map[string]string{"Standard03CandidateViewName": "Ab76Name", "Standard03CandidateViewState": "Ab76State", "Standard03CandidateViewRaw": "Ab76Raw", "Standard03CandidateViewValue": "Ab76Value", "Standard03CandidateViewWidth": "Ab76Width", "Standard03CandidateViewMaximum": "Ab76Maximum", "Standard03CandidateViewDefault": "Ab76Default", "Standard03CandidateViewEncoding": "Ab76Encoding"}},
		{Standard03CandidateView{}, map[string]string{"Standard03CandidateViewSelection": "Ab76Selection", "Standard03CandidateViewState": "Ab76State", "Standard03CandidateViewRaw": "Ab76Raw", "Standard03CandidateViewRemainder": "Ab76Remainder", "Standard03CandidateViewComponents": "Ab76Components", "Standard03CandidateViewLocalCounter": "Ab76LocalCounter", "Standard03CandidateViewLocalHours": "Ab76LocalHours", "Standard03CandidateViewLocalMinute": "Ab76LocalMinute", "Standard03CandidateViewImplementationGrade": "Ab76ImplementationGrade", "Standard03CandidateViewMeaningGrade": "Ab76MeaningGrade", "Standard03CandidateViewNativeProperties": "Ab76NativeProperties"}},
	}
	for _, tc := range cases {
		data, err := json.Marshal(tc.value)
		if err != nil {
			t.Fatal(err)
		}
		var object map[string]json.RawMessage
		if err := json.Unmarshal(data, &object); err != nil {
			t.Fatal(err)
		}
		typ := reflect.TypeOf(tc.value)
		for field, key := range tc.keys {
			if _, ok := object[key]; !ok {
				t.Errorf("%s.%s missing JSON key %s", typ, field, key)
			}
			if _, ok := object[field]; ok {
				t.Errorf("%s exposes renamed key %s", typ, field)
			}
		}
	}
	var selection Standard03CandidateViewSelection
	if err := json.Unmarshal([]byte(`{"Ab76Family":"standard-counter"}`), &selection); err != nil {
		t.Fatal(err)
	}
	if selection.Standard03CandidateViewFamily != "standard-counter" {
		t.Fatal("original family key was not decoded")
	}
}
'''


EBUSSTANDARD = SimpleNamespace(
    name="ebusstandard",
    package_dir="catalog/ebusstandard",
    # Coded file stem -> behavior file stem. Each applies to the .go file and
    # to its _test.go twin. Testdata fixtures keep their names.
    file_moves=(
        ("ab15aq15_standard07_field_views", "standard07_field_views"),
        ("ab17aq17_standardff_extent_view", "standardff_extent_view"),
        ("ab25an25_standard09_body_views", "standard09_body_views"),
        ("ab28_standard07_residual_fields", "standard07_residual_fields"),
        ("ab49_fe_metadata", "standardfe_metadata"),
        ("ab54_standard_ff_body_states", "standardff_body_states"),
        ("ab76_standard03_candidate_views", "standard03_candidate_views"),
        ("ab77_07_candidates", "standard07_candidates"),
        ("ab78_standard05_candidate_views", "standard05_candidate_views"),
        ("ab79_09_candidates", "standard09_candidates"),
        ("ab81_0f_candidates", "standard0f_candidates"),
        ("ab83_03_counter_bodies", "standard03_counter_bodies"),
        ("ab84_05_burner_bodies", "standard05_burner_bodies"),
    ),
    # Code prefix (lowercase) -> UpperCamelCase stem derived from the
    # destination file. Ab<code>X becomes <Stem>X and ab<code>X becomes
    # <stem>X.
    stems={
        "ab15aq15": "Standard07FieldView",
        "ab17aq17": "StandardFFExtentView",
        "ab25an25": "Standard09BodyView",
        "ab28": "Standard07Residual",
        "ab49": "StandardFEMetadata",
        "ab54": "StandardFFBodyState",
        "ab76": "Standard03CandidateView",
        "ab77": "Standard07Candidate",
        "ab78": "Standard05CandidateView",
        "ab79": "Standard09Candidate",
        "ab81": "Standard0fCandidate",
        "ab83": "Standard03CounterBody",
        "ab84": "Standard05BurnerBody",
    },
    # Exported names whose remainder repeats the stem; the stem words that
    # would repeat are dropped.
    exported_overrides={
        "Ab25An25ViewAddress": "Standard09ViewAddress",
        "Ab25An25ViewBody": "Standard09ViewBody",
        "Ab25An25ViewCount": "Standard09ViewCount",
        "Ab49FEMetadata": "StandardFEMetadata",
        "Ab76View": "Standard03CandidateView",
        "Ab77CandidateMetadata": "Standard07CandidateMetadata",
        "Ab77Candidates": "Standard07Candidates",
        "Ab78View": "Standard05CandidateView",
        "Ab79Candidates": "Standard09Candidates",
        "Ab81Candidate": "Standard0fCandidate",
        "Ab81Candidates": "Standard0fCandidates",
        "Ab83Body": "Standard03CounterBody",
        "Ab83BodyState": "Standard03CounterBodyState",
        "Ab83CounterBodies": "Standard03CounterBodies",
        "Ab83LookupBody": "Standard03CounterLookupBody",
        "Ab84BurnerBodies": "Standard05BurnerBodies",
        "Ab84BurnerField": "Standard05BurnerField",
        "Ab84BurnerMetadata": "Standard05BurnerMetadata",
        "Ab84LookupBurnerField": "Standard05LookupBurnerField",
    },
    # Resolved collisions: the encoding/json import alias of the standard 09
    # test would otherwise equal the name of its fixture byte slice
    # (ab25an25JSON).
    collision_overrides={
        "ab25an25json": "standard09BodyViewEncodingJSON",
    },
    # A lowercase remainder (import aliases, short locals) is capitalized;
    # these words take their Go initialism form instead.
    initialisms={"json": "JSON", "id": "ID"},
    # Test case structs decoded from fixtures through untagged field names;
    # the old field name is the fixture key and becomes an explicit tag.
    fixture_bound={
        "ab76_standard03_candidate_views_test.go": ("ab76Cases", {
            "Ab76Label": "Ab76Label", "Ab76Operation": "Ab76Operation",
            "Ab76Body": "Ab76Body", "Ab76Names": "Ab76Names",
            "Ab76Values": "Ab76Values", "Ab76Counter": "Ab76Counter",
            "Ab76Hours": "Ab76Hours", "Ab76Minute": "Ab76Minute"}),
        "ab78_standard05_candidate_views_test.go": ("ab78Cases", {
            "Ab78Label": "ab78Label", "Ab78Body": "ab78Body",
            "Ab78Values": "ab78Values", "Ab78Name": "ab78Name",
            "Ab78Domains": "ab78Domains"}),
    },
    reflect_names={},
    import_unaliases={},
    # Fields are tagged by preserve_projection_json_keys and the explicit
    # fixture_bound table.
    tag_fields=False,
    golden="catalog/ebusstandard/testdata/provider_abi.golden.txt",
    regression=("canonical_names_test.go", lambda tagged, mapping: PROJECTION_JSON_KEYS_TEST),
    internal_dir="catalog/ebusstandard/internal",
)


VAILLANT_01_29 = SimpleNamespace(
    name="vaillant-01-29",
    package_dir="catalog/vaillant",
    # The work-item code is dropped and the rest of the name kept. Two
    # evaluators take an _evaluator suffix because a catalog file holds the
    # plain name: b512_heatpump_io.go (from ab19) and, in the next slice,
    # b51a_mirror_fields.go (from ab90).
    file_moves=(
        ("ab01ao01_b505_sentinel_views", "b505_sentinel_views"),
        ("ab01ap01_b531_capability_state", "b531_capability_state"),
        ("ab01aq01_b504_candidate_fields", "b504_candidate_fields"),
        ("ab03ao03_b506_parameter_views", "b506_parameter_views"),
        ("ab03ap03_b531_ordinary_limit", "b531_ordinary_limit"),
        ("ab03aq03_b505_override_fields", "b505_override_fields"),
        ("ab05_schema", "schema"),
        ("ab05an05_b503_history_branches", "b503_history_branches"),
        ("ab05ao05_b507_quantity_views", "b507_quantity_views"),
        ("ab05ap05_b516_cycle_reducer", "b516_cycle_reducer"),
        ("ab05aq05_b505_conditional_roles", "b505_conditional_roles"),
        ("ab07ao07_b50e_row_spans", "b50e_row_spans"),
        ("ab07ap07_b555_weekly_refresh", "b555_weekly_refresh"),
        ("ab07aq07_b503_control_fields", "b503_control_fields"),
        ("ab09ao09_b50f_field_views", "b50f_field_views"),
        ("ab09ap09_b555_manager_gate", "b555_manager_gate"),
        ("ab09aq09_b515_period_fields", "b515_period_fields"),
        ("ab11an11_b511_overlap_views", "b511_overlap_views"),
        ("ab11ao11_b510_cache_transition", "b510_cache_transition"),
        ("ab11ap11_b512_fan_fields", "b512_fan_fields"),
        ("ab11aq11_b517_local_guard", "b517_local_guard"),
        ("ab12_b50f_status_spans", "b50f_status_spans"),
        ("ab13_b507_energy_candidates", "b507_energy_candidates"),
        ("ab13an13_b512_dhw_views", "b512_dhw_views"),
        ("ab13ao13_b511_failure_transition", "b511_failure_transition"),
        ("ab13ap13_b512_heatpump_io", "b512_heatpump_io_evaluator"),
        ("ab13aq13_b521_field_views", "b521_field_views"),
        ("ab14_b509_op59_quantities", "b509_op59_quantities"),
        ("ab15_b504_circuit_roles", "b504_circuit_roles"),
        ("ab15an15_b516_clock_word_views", "b516_clock_word_views"),
        ("ab15ao15_b512_feedback_views", "b512_feedback_views"),
        ("ab15ap15_b512_emergency_flow", "b512_emergency_flow"),
        ("ab16_b505_operating_roles", "b505_operating_roles"),
        ("ab17_b506_rf_sequence", "b506_rf_sequence"),
        ("ab17an17_b516_two_word_views", "b516_two_word_views"),
        ("ab17ao17_b513_information_views", "b513_information_views"),
        ("ab17ap17_b509_diagnostic_fields", "b509_diagnostic_fields"),
        ("ab18_b512_feedback_fallback", "b512_feedback_fallback"),
        ("ab19_b512_heatpump_io", "b512_heatpump_io"),
        ("ab19an19_b531_limit_numeric_views", "b531_limit_numeric_views"),
        ("ab19ao19_b514_menu_views", "b514_menu_views"),
        ("ab19ap19_b51a_mirror_fields", "b51a_mirror_fields_evaluator"),
        ("ab19aq19_b51a_handshake_fields", "b51a_handshake_fields"),
        ("ab20_b516_directional_types", "b516_directional_types"),
        ("ab21_b519_reply_policy", "b519_reply_policy"),
        ("ab21ao21_b519_opaque_views", "b519_opaque_views"),
        ("ab21ap21_b523_actuator_fields", "b523_actuator_fields"),
        ("ab21aq21_b51a_diagnostic_fields", "b51a_diagnostic_fields"),
        ("ab23an23_b51b_local_budget_crc", "b51b_local_budget_crc"),
        ("ab23ao23_b520_actuator_views", "b520_actuator_views"),
        ("ab23ap23_b516_energy_type_fields", "b516_energy_type_fields"),
        ("ab23aq23_b512_valve_pump_fields", "b512_valve_pump_fields"),
        ("ab25ao25_b522_demand_views", "b522_demand_views"),
        ("ab25ap25_b525_quality_config", "b525_quality_config"),
        ("ab25aq25_b512_power_stepper_fields", "b512_power_stepper_fields"),
        ("ab26_b555_additional_forms", "b555_additional_forms"),
        ("ab27_holiday", "holiday"),
        ("ab27ao27_b523_sensor_views", "b523_sensor_views"),
        ("ab29_room_actions", "room_actions"),
    ),
    # Stem: the destination file name in UpperCamelCase, the service in
    # upper-case hex, initialisms in Go form and the last word singular.
    stems={
        "ab01ao01": "B505SentinelView",
        "ab01ap01": "B531CapabilityState",
        "ab01aq01": "B504CandidateField",
        "ab03ao03": "B506ParameterView",
        "ab03ap03": "B531OrdinaryLimit",
        "ab03aq03": "B505OverrideField",
        "ab05": "Schema",
        "ab05an05": "B503HistoryBranch",
        "ab05ao05": "B507QuantityView",
        "ab05ap05": "B516CycleReducer",
        "ab05aq05": "B505ConditionalRole",
        "ab07ao07": "B50ERowSpan",
        "ab07ap07": "B555WeeklyRefresh",
        "ab07aq07": "B503ControlField",
        "ab09ao09": "B50FFieldView",
        "ab09ap09": "B555ManagerGate",
        "ab09aq09": "B515PeriodField",
        "ab11an11": "B511OverlapView",
        "ab11ao11": "B510CacheTransition",
        "ab11ap11": "B512FanField",
        "ab11aq11": "B517LocalGuard",
        "ab12": "B50FStatusSpan",
        "ab13": "B507EnergyCandidate",
        "ab13an13": "B512DHWView",
        "ab13ao13": "B511FailureTransition",
        "ab13ap13": "B512HeatpumpIOEvaluator",
        "ab13aq13": "B521FieldView",
        "ab14": "B509Op59Quantity",
        "ab15": "B504CircuitRole",
        "ab15an15": "B516ClockWordView",
        "ab15ao15": "B512FeedbackView",
        "ab15ap15": "B512EmergencyFlow",
        "ab16": "B505OperatingRole",
        "ab17": "B506RFSequence",
        "ab17an17": "B516TwoWordView",
        "ab17ao17": "B513InformationView",
        "ab17ap17": "B509DiagnosticField",
        "ab18": "B512FeedbackFallback",
        "ab19": "B512HeatpumpIO",
        "ab19an19": "B531LimitNumericView",
        "ab19ao19": "B514MenuView",
        "ab19ap19": "B51AMirrorFieldEvaluator",
        "ab19aq19": "B51AHandshakeField",
        "ab20": "B516DirectionalType",
        "ab21": "B519ReplyPolicy",
        "ab21ao21": "B519OpaqueView",
        "ab21ap21": "B523ActuatorField",
        "ab21aq21": "B51ADiagnosticField",
        "ab23an23": "B51BLocalBudgetCRC",
        "ab23ao23": "B520ActuatorView",
        "ab23ap23": "B516EnergyTypeField",
        "ab23aq23": "B512ValvePumpField",
        "ab25ao25": "B522DemandView",
        "ab25ap25": "B525QualityConfig",
        "ab25aq25": "B512PowerStepperField",
        "ab26": "B555AdditionalForm",
        "ab27": "Holiday",
        "ab27ao27": "B523SensorView",
        "ab29": "RoomAction",
    },
    exported_overrides={
        "Ab11An11View": "B511OverlapView",
        "Ab13Ap13Evaluate": "B512HeatpumpIOEvaluate",
        "Ab15Role": "B504CircuitRole",
        "Ab15An15WordView": "B516ClockWordView",
        "Ab16Role": "B505OperatingRole",
        "Ab19Ap19Evaluate": "B51AMirrorFieldEvaluate",
        "Ab21Policy": "B519ReplyPolicy",
        "Ab23An23CRC": "B51BLocalBudgetCRC",
        "Ab23An23CRCWithInitial": "B51BLocalBudgetCRCWithInitial",
        "Ab27HolidayCatalog": "HolidayCatalog",
        "Ab29RoomActions": "RoomActions",
    },
    # Resolved collisions: the bytes import alias (ab<code>bytes) would
    # otherwise equal the name of a byte slice of the same file
    # (ab<code>Bytes).
    collision_overrides={
        "ab03ap03bytes": "b531OrdinaryLimitBytesPackage",
        "ab05ao05bytes": "b507QuantityViewBytesPackage",
        "ab07aq07bytes": "b503ControlFieldBytesPackage",
        "ab09ap09bytes": "b555ManagerGateBytesPackage",
        "ab11ao11bytes": "b510CacheTransitionBytesPackage",
        "ab11ap11bytes": "b512FanFieldBytesPackage",
        "ab13ao13bytes": "b511FailureTransitionBytesPackage",
        "ab15ao15bytes": "b512FeedbackViewBytesPackage",
        "ab19bytes": "b512HeatpumpIOBytesPackage",
        "ab19ap19bytes": "b51AMirrorFieldEvaluatorBytesPackage",
    },
    initialisms={"json": "JSON", "id": "ID", "io": "IO", "utf8": "UTF8",
                 "ascii": "ASCII"},
    fixture_bound={},
    # Test string literals that name a renamed struct field for reflection
    # (reflect.Value.FieldByName); they follow the field rename.
    reflect_names={
        "ab05an05_b503_history_branches_test.go": (
            "Ab05An05Family", "Ab05An05Variant", "Ab05An05Operation",
            "Ab05An05Direction"),
        "ab12_b50f_status_spans_test.go": (
            "Ab12Secondary", "Ab12Operation", "Ab12Direction", "Ab12Group",
            "Ab12RequestedInstance", "Ab12Register", "Ab12CommandRole"),
        "ab15_b504_circuit_roles_test.go": (
            "Ab15Family", "Ab15Variant", "Ab15Primary", "Ab15Secondary",
            "Ab15Operation", "Ab15Direction", "Ab15Field", "Ab15Alternative",
            "Ab15Group", "Ab15InstanceScope", "Ab15Register"),
    },
    # File-scoped import aliases that only carried the code; the import loses
    # its alias and the uses take the package name. The stem-derived alias
    # (b511FailureTransitionJSON) would make gofmt split the one-line
    # result closure of B511FailureTransitionInterpret and push the function
    # over the size ratchet (150 lines).
    import_unaliases={
        "ab13ao13json": (("ab13ao13_b511_failure_transition.go",
                          "ab13ao13_b511_failure_transition_test.go"), "encoding/json"),
    },
    # Exported production struct fields are all tagged; test structs are
    # tagged when every renamed field is a key of the package fixtures.
    tag_fields=True,
    golden=None,
    # The holiday test file imports encoding/json, reflect and testing under
    # their own names and stays far below the functions-per-file budget.
    regression=("ab27_holiday_test.go", lambda tagged, mapping: field_json_keys_test(tagged, mapping)),
    internal_dir=None,
)

SLICES = {cfg.name: cfg for cfg in (EBUSSTANDARD, VAILLANT_01_29)}


def lex(source):
    tokens, index, size = [], 0, len(source)
    while index < size:
        char = source[index]
        if source.startswith("//", index):
            end = source.find("\n", index)
            end = size if end < 0 else end
            kind = "comment"
        elif source.startswith("/*", index):
            end = source.index("*/", index) + 2
            kind = "comment"
        elif char in "\"'":
            end = index + 1
            while source[end] != char:
                if source[end] == "\n":
                    raise ValueError("unterminated literal")
                end += 2 if source[end] == "\\" else 1
            end += 1
            kind = "string"
        elif char == "`":
            end = source.index("`", index + 1) + 1
            kind = "string"
        elif char.isascii() and (char.isalpha() or char == "_"):
            end = IDENT.match(source, index).end()
            kind = "ident"
        elif char.isascii() and char.isdigit():
            end = NUMBER.match(source, index).end()
            kind = "number"
        elif not char.isascii():
            raise ValueError("non-ASCII Go source outside literals")
        else:
            end = index + 1
            kind = "other"
        tokens.append((kind, source[index:end]))
        index = end
    if "".join(text for _, text in tokens) != source:
        raise ValueError("lexer did not reproduce the source")
    return tokens


def coded(name, cfg):
    match = CODE.fullmatch(name)
    if match and match.group(1).lower() in cfg.stems:
        return match.group(1).lower(), match.group(2)
    return None


def new_name(name, cfg):
    if name in cfg.import_unaliases:
        return cfg.import_unaliases[name][1].rsplit("/", 1)[1]
    if name in cfg.exported_overrides:
        return cfg.exported_overrides[name]
    if name in cfg.collision_overrides:
        return cfg.collision_overrides[name]
    code, rest = coded(name, cfg)
    if rest[:1].islower():
        rest = cfg.initialisms.get(rest, rest[:1].upper() + rest[1:])
    stem = cfg.stems[code]
    if name[0].islower():
        stem = stem[:1].lower() + stem[1:]
    return stem + rest


def package_files(root, cfg):
    return sorted(p for p in (root / cfg.package_dir).glob("*.go") if p.is_file())


def build_map(root, cfg):
    old, other = set(), set()
    for path in package_files(root, cfg):
        for kind, text in lex(path.read_text()):
            if kind == "ident":
                (old if coded(text, cfg) else other).add(text)
    if cfg.internal_dir:
        for path in sorted((root / cfg.internal_dir).rglob("*.go")):
            if any(kind == "ident" and coded(text, cfg) for kind, text in lex(path.read_text())):
                raise SystemExit(f"coded identifier outside the package scope: {path}")
    for alias, (file_names, import_path) in cfg.import_unaliases.items():
        users = [path.name for path in package_files(root, cfg)
                 if alias in {text for kind, text in lex(path.read_text()) if kind == "ident"}]
        if sorted(users) != sorted(file_names):
            raise SystemExit(f"import alias {alias} is used outside {file_names}")
        for file_name in file_names:
            idents = [text for kind, text in lex((root / cfg.package_dir / file_name).read_text())
                      if kind == "ident"]
            if new_name(alias, cfg) in idents:
                raise SystemExit(f"import alias {alias} cannot take its package name")
        other.discard(new_name(alias, cfg))
    for table in (cfg.exported_overrides, cfg.collision_overrides, cfg.import_unaliases):
        unused = sorted(set(table) - old)
        if unused:
            raise SystemExit(f"override for an absent identifier: {unused}")
    mapping = {name: new_name(name, cfg) for name in sorted(old)}
    seen = {}
    for name, target in mapping.items():
        if target in seen:
            raise SystemExit(f"collision: {seen[target]} and {name} -> {target}")
        seen[target] = name
        if name in cfg.import_unaliases:
            continue
        if target in other or target in KEYWORDS or CODE.fullmatch(target):
            raise SystemExit(f"collision with an existing name: {name} -> {target}")
        if name[0].isupper() != target[0].isupper():
            raise SystemExit(f"export status changes: {name} -> {target}")
    return mapping


def check_consumers(root, mapping, cfg):
    import_path = MODULE + cfg.package_dir
    default_alias = cfg.package_dir.rsplit("/", 1)[1]
    for path in sorted(root.rglob("*.go")):
        rel = path.relative_to(root)
        if rel.parts[0].startswith(".") or rel.parent.as_posix() == cfg.package_dir:
            continue
        tokens = [t for t in lex(path.read_text()) if t[0] != "comment" and t[1].strip()]
        aliases = set()
        for index, (kind, text) in enumerate(tokens):
            if kind == "string" and text.strip('"`') == import_path:
                kind, alias = tokens[index - 1]
                if kind != "ident" or alias == "import":
                    alias = default_alias
                aliases.add(alias)
        for index in range(len(tokens) - 2):
            if (tokens[index][1] in aliases and tokens[index + 1][1] == "."
                    and tokens[index + 2][1] in mapping):
                raise SystemExit(f"consumer {rel} uses {tokens[index + 2][1]}")


def add_fixture_tags(source, variable, fields):
    lines = source.split("\n")
    opener = f"\tvar {variable} []struct {{"
    if lines.count(opener) != 1:
        raise SystemExit(f"fixture struct {variable} not found exactly once")
    index = lines.index(opener) + 1
    tagged = set()
    while lines[index] != "\t}":
        match = re.fullmatch(r"\t\t([A-Za-z0-9_]+)(\s+)(\S+)", lines[index])
        if not match or match.group(1) not in fields:
            raise SystemExit(f"unexpected fixture field line: {lines[index]!r}")
        lines[index] += f' `json:"{fields[match.group(1)]}"`'
        tagged.add(match.group(1))
        index += 1
    if tagged != set(fields):
        raise SystemExit(f"fixture fields differ for {variable}")
    return "\n".join(lines)


def rename_comment(text, mapping):
    if text.startswith("//go:"):
        return text
    return WORD.sub(lambda m: mapping.get(m.group(), m.group()), text)


def preserve_projection_json_keys(source, cfg):
    """Keep implicit JSON field keys before renaming exported struct fields."""
    lines = source.splitlines(keepends=True)
    in_struct = False
    for index, line in enumerate(lines):
        if re.fullmatch(r"type [A-Z]\w* struct \{\n", line):
            in_struct = True
        elif in_struct and line == "}\n":
            in_struct = False
        elif in_struct:
            match = re.fullmatch(r"\t([A-Z]\w*)\s+([^`\n]+)\n", line)
            if match and coded(match.group(1), cfg):
                lines[index] = line.rstrip("\n") + f' `json:"{match.group(1)}"`\n'
    return "".join(lines)


STRUCT_OPEN = re.compile(r".*\bstruct \{")
NAMED_STRUCT = re.compile(r"type ([A-Za-z_]\w*) struct \{")
NESTED_FIELD = re.compile(r"([A-Za-z_]\w*) +[\[\]*.\w]*struct \{")
FIELD_LINE = re.compile(r"([A-Za-z_]\w*(?:, [A-Za-z_]\w*)*)( +)([^`]*?)( +`.*`)?")
EMBEDDED_LINE = re.compile(r"\*?[A-Za-z_][\w.]*( +`.*`)?")


def json_tag(name):
    return f' `json:"{name}"`'


def renamed_field(name, cfg):
    return name[:1].isupper() and coded(name, cfg) is not None


def tag_struct_fields(source, cfg, wanted):
    """Give renamed exported fields an explicit tag with their old name.

    Works on gofmt output: one field declaration per line, struct bodies
    indented one tab deeper than the line that opens them. wanted(names)
    decides per struct from the renamed field names it declares. A
    declaration with several names is split so every field carries its own
    key. Returns the source and the (type, fields) pairs that were tagged.
    """
    lines = source.split("\n")
    edits, tagged, stack = {}, [], []
    for index, line in enumerate(lines):
        body = line.lstrip("\t")
        depth = len(line) - len(body)
        if stack and depth == stack[-1]["depth"] and body.startswith("}"):
            frame = stack.pop()
            if frame["parent"] is not None:
                frame["parent"]["close"] = index
                frame["parent"]["tagged"] = "`" in body
            names = [n for f in frame["fields"] for n in f["names"] if renamed_field(n, cfg)]
            if not names or not wanted(names):
                continue
            done = []
            for field in frame["fields"]:
                need = [n for n in field["names"] if renamed_field(n, cfg)]
                if field["tagged"] or not need:
                    continue
                if field["comment"]:
                    raise SystemExit(f"field with a trailing comment: {field['names']}")
                if field["type"] is None:
                    edits[field["close"]] = lines[field["close"]] + json_tag(field["names"][0])
                elif len(field["names"]) == 1:
                    edits[field["line"]] = lines[field["line"]] + json_tag(field["names"][0])
                else:
                    indent = "\t" * (frame["depth"] + 1)
                    edits[field["line"]] = "\n".join(
                        indent + name + " " + field["type"]
                        + (json_tag(name) if renamed_field(name, cfg) else "")
                        for name in field["names"])
                done.extend(need)
            if done:
                tagged.append((frame["type"], done))
            continue
        if STRUCT_OPEN.fullmatch(body) and not body.startswith("//"):
            frame = {"depth": depth, "fields": [], "parent": None, "type": None}
            named = NAMED_STRUCT.fullmatch(body)
            if named and depth == 0:
                frame["type"] = named.group(1)
            if stack and depth == stack[-1]["depth"] + 1:
                nested = NESTED_FIELD.fullmatch(body)
                if not nested:
                    raise SystemExit(f"unexpected nested struct line: {line!r}")
                field = {"names": [nested.group(1)], "type": None, "line": index,
                         "close": None, "tagged": False, "comment": False}
                stack[-1]["fields"].append(field)
                frame["parent"] = field
            stack.append(frame)
            continue
        if not stack or depth != stack[-1]["depth"] + 1 or not body or body.startswith("//"):
            continue
        code_part, _, comment = body.partition(" //")
        match = FIELD_LINE.fullmatch(code_part)
        if match is None or match.group(3) == "":
            if EMBEDDED_LINE.fullmatch(code_part) is None:
                raise SystemExit(f"unexpected struct field line: {line!r}")
            if renamed_field(code_part.lstrip("*").split(" ")[0], cfg):
                raise SystemExit(f"embedded renamed field: {line!r}")
            continue
        stack[-1]["fields"].append({
            "names": match.group(1).split(", "), "type": match.group(3),
            "line": index, "close": None, "tagged": match.group(4) is not None,
            "comment": bool(comment)})
    if stack:
        raise SystemExit("unterminated struct body")
    out = [edits.get(index, line) for index, line in enumerate(lines)]
    return "\n".join(out), tagged


def json_keys(value, keys):
    if isinstance(value, dict):
        for key, item in value.items():
            keys.add(key.lower())
            json_keys(item, keys)
    elif isinstance(value, list):
        for item in value:
            json_keys(item, keys)
    return keys


def fixture_keys(root, cfg):
    keys = set()
    for folder in ("data", "testdata"):
        for path in sorted((root / cfg.package_dir / folder).glob("*.json")):
            json_keys(json.loads(path.read_text()), keys)
    return keys


def field_json_keys_test(tagged, mapping):
    """Regression test that pins the JSON key of every tagged field."""
    rows = []
    for type_name, fields in tagged:
        pairs = ", ".join(f'"{mapping[name]}": "{name}"' for name in fields)
        rows.append(f"\t\t{{{mapping.get(type_name, type_name)}{{}}, "
                    f"map[string]string{{{pairs}}}}},\n")
    return (
        "func TestRenamedFieldJSONKeys(t *testing.T) {\n"
        "\tcases := []struct {\n"
        "\t\tvalue any\n"
        "\t\tkeys  map[string]string\n"
        "\t}{\n" + "".join(rows) +
        "\t}\n"
        "\tfor _, tc := range cases {\n"
        "\t\tdata, err := json.Marshal(tc.value)\n"
        "\t\tif err != nil {\n"
        "\t\t\tt.Fatal(err)\n"
        "\t\t}\n"
        "\t\tvar object map[string]json.RawMessage\n"
        "\t\tif err := json.Unmarshal(data, &object); err != nil {\n"
        "\t\t\tt.Fatal(err)\n"
        "\t\t}\n"
        "\t\ttyp := reflect.TypeOf(tc.value)\n"
        "\t\tfor field, key := range tc.keys {\n"
        "\t\t\tif _, ok := object[key]; !ok {\n"
        "\t\t\t\tt.Errorf(\"%s.%s missing JSON key %s\", typ, field, key)\n"
        "\t\t\t}\n"
        "\t\t\tif _, ok := object[field]; ok {\n"
        "\t\t\t\tt.Errorf(\"%s exposes renamed key %s\", typ, field)\n"
        "\t\t\t}\n"
        "\t\t\tif member, ok := typ.FieldByName(field); !ok || member.Tag.Get(\"json\") != key {\n"
        "\t\t\t\tt.Errorf(\"%s.%s does not decode JSON key %s\", typ, field, key)\n"
        "\t\t\t}\n"
        "\t\t}\n"
        "\t}\n"
        "}\n")


def drop_import_alias(source, alias, import_path):
    spec = f'\t{alias} "{import_path}"\n'
    if source.count(spec) != 1:
        raise SystemExit(f"import spec of {alias} not found exactly once")
    return source.replace(spec, f'\t"{import_path}"\n')


def rewrite_go(source, mapping, names=()):
    out = []
    for kind, text in lex(source):
        if kind == "ident":
            text = mapping.get(text, text)
        elif kind == "comment":
            text = rename_comment(text, mapping)
        elif kind == "string" and text[:1] == '"' and text[1:-1] in names:
            text = '"' + mapping[text[1:-1]] + '"'
        out.append(text)
    return "".join(out)


def rename_symbol(symbol, mapping):
    return ".".join(mapping.get(part, part) for part in symbol.split("."))


def moved_paths(cfg):
    moves = {}
    for old, new in cfg.file_moves:
        for suffix in (".go", "_test.go"):
            moves[f"{cfg.package_dir}/{old}{suffix}"] = f"{cfg.package_dir}/{new}{suffix}"
    return moves


def load_json(path):
    return json.loads(path.read_text())


def dump_json(path, data, ensure_ascii):
    path.write_text(json.dumps(data, indent=2, ensure_ascii=ensure_ascii) + "\n")


def replace_paths(value, moves):
    if isinstance(value, str):
        return moves.get(value, value)
    if isinstance(value, list):
        return [replace_paths(item, moves) for item in value]
    if isinstance(value, dict):
        return {key: replace_paths(item, moves) for key, item in value.items()}
    return value


def behavior_tokens(root):
    sys.dont_write_bytecode = True
    spec = importlib.util.spec_from_file_location(
        "behavior_names", root / "scripts/check_behavior_names.py")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    vocabulary, _ = module.load_policy(root / "scripts/behavior-name-baseline.json")
    return module, vocabulary


def update_baselines(root, mapping, moves, cfg):
    prefix = cfg.package_dir + "/"
    module, vocabulary = behavior_tokens(root)
    path = root / "scripts/behavior-name-baseline.json"
    data = load_json(path)
    rows = []
    for row in data["baseline"]:
        if row["path"] in moves:
            row = dict(row, path=moves[row["path"]])
            if row["category"] == "filename":
                row["symbol"] = row["path"].rsplit("/", 1)[1]
        if row["path"].startswith(prefix) and row["category"] == "identifier":
            row = dict(row, symbol=mapping.get(row["symbol"], row["symbol"]))
        if module.violations([row], vocabulary):
            rows.append(row)
    data["baseline"] = sorted(rows, key=lambda r: (r["path"], r["symbol"], r["category"]))
    dump_json(path, data, True)

    path = root / "scripts/go-lint-baseline.json"
    data = load_json(path)
    for row in data["entries"]:
        if row["path"].startswith(prefix):
            row["path"] = moves.get(row["path"], row["path"])
            row["symbol"] = rename_symbol(row["symbol"], mapping)
    data["entries"].sort(key=lambda r: (r["path"], r["symbol"], r["rule"]))
    dump_json(path, data, True)

    path = root / "scripts/size-baseline.json"
    data = load_json(path)
    for row in data["debt"]:
        if row["path"].startswith(prefix):
            row["path"] = moves.get(row["path"], row["path"])
            row["symbol"] = rename_symbol(row["symbol"], mapping)
    data["debt"].sort(key=lambda r: (r["path"], r["symbol"], r["metric"]))
    dump_json(path, data, True)


def update_inventories(root, mapping, moves, cfg):
    path = root / "naming-inventory.json"
    data = load_json(path)
    data = replace_paths(data, moves)
    data["source_files"] = sorted(data["source_files"])
    for row in data["rows"]:
        if (row["owning_file"].startswith(cfg.package_dir + "/")
                and row["kind"] in ("go_identifier", "go_field_candidate")):
            row["current_name"] = mapping.get(row["current_name"], row["current_name"])
            if "target" in row:
                row["target"] = mapping.get(row["target"], row["target"])
            for key in ("declaration", "field_candidate"):
                if key in row["identity"]:
                    row["identity"][key] = mapping.get(row["identity"][key], row["identity"][key])
    dump_json(path, data, False)

    path = root / "naming-consumers.json"
    before = load_json(path)
    after = replace_paths(before, moves)
    if after != before:
        dump_json(path, after, True)


def update_golden(root, cfg):
    # Use the snapshot's own serializer so tags added to projections are
    # recorded along with renamed symbols and their signatures.
    subprocess.run(
        ["go", "test", "./" + cfg.package_dir,
         "-run", "^TestProvider_ABISnapshot$", "-count=1", "-update"],
        cwd=root, check=True)


def code_label(code):
    return code[:1].upper() + code[1:4] + code[4:5].upper() + code[5:]


def changelog_entry(mapping, cfg):
    files = "".join(f"  - `{old}.go` -> `{new}.go`\n" for old, new in cfg.file_moves)
    stems = "".join(f"  - `{code_label(code)}...` -> `{stem}...`\n"
                    for code, stem in cfg.stems.items())
    overrides = "".join(f"  - `{old}` -> `{new}`\n"
                        for old, new in sorted(cfg.exported_overrides.items()))
    exported = sum(1 for name in mapping if name[0].isupper())
    counts = (f"  Identifiers that start with a code ({exported} exported, "
              f"{len(mapping) - exported} unexported)\n"
              "  take the stem of their destination file; `Ab<code>X` becomes\n"
              "  `<Stem>X` and `ab<code>X` becomes `<stem>X`:\n")
    if cfg is EBUSSTANDARD:
        return (
            "- Remove work-item codes from `catalog/ebusstandard` file names and Go\n"
            "  identifiers. Each production file and its test twin is renamed; the\n"
            "  testdata fixtures keep their names:\n" + files + counts + stems +
            "  Exported names whose remainder repeats the stem drop the repeated\n"
            "  words:\n" + overrides +
            "  Consumers must update references; there are no aliases. JSON tags,\n"
            "  string values, error texts, fixture names and test function names are\n"
            "  unchanged. `scripts/rename_catalog_names.py` reproduces the change.\n"
            "  Gateway adoption and release bill of materials reconciliation remain\n"
            "  required.\n")
    return (
        "- Remove work-item codes from the `catalog/vaillant` files whose code\n"
        "  number is 01 to 29 and from the Go identifiers that start with their\n"
        "  codes. Each production file and its test twin is renamed; embedded\n"
        "  data and testdata fixtures keep their names. The two evaluators\n"
        "  `b512_heatpump_io_evaluator.go` and `b51a_mirror_fields_evaluator.go`\n"
        "  take an `_evaluator` suffix because a catalog file holds the plain\n"
        "  name:\n" + files +
        counts + stems +
        "  Exported names whose remainder repeats the stem drop the repeated\n"
        "  words:\n" + overrides +
        "  The `bytes` import aliases that would equal a byte slice name take a\n"
        "  `BytesPackage` suffix, and the `encoding/json` imports of\n"
        "  `b511_failure_transition.go` and its test drop their coded alias so\n"
        "  `B511FailureTransitionInterpret` stays within the size ratchet.\n"
        "  Files with codes 30 to 99 that use the renamed `Schema...` names are\n"
        "  edited in place and keep their names until the next slice. Renamed\n"
        "  exported struct fields that had no JSON tag carry an explicit tag\n"
        "  with the unchanged key, so serialized output and embedded data\n"
        "  decoding stay the same; `TestRenamedFieldJSONKeys` pins every key.\n"
        "  Consumers must update references; there are no aliases. JSON tags,\n"
        "  string values, error texts, `go:embed` paths, data and fixture names\n"
        "  and test function names are unchanged.\n"
        "  `scripts/rename_catalog_names.py --slice vaillant-01-29` reproduces\n"
        "  the change. Gateway adoption and release bill of materials\n"
        "  reconciliation remain required.\n")


def append_test(source, text):
    name = text.split("(", 1)[0]
    if name in source:
        raise SystemExit(f"regression test already exists: {name}")
    return source.rstrip() + "\n\n" + text


def run(root, print_map, cfg):
    moves = moved_paths(cfg)
    for old, new in moves.items():
        if not (root / old).is_file():
            raise SystemExit(f"source file missing: {old}")
        if (root / new).exists():
            raise SystemExit(f"target already exists: {new}")
    mapping = build_map(root, cfg)
    check_consumers(root, mapping, cfg)
    if print_map:
        for name, target in mapping.items():
            print(f"{name} {target}")
        return
    keys = fixture_keys(root, cfg) if cfg.tag_fields else set()
    sources, tagged = {}, []
    for path in package_files(root, cfg):
        source = original = path.read_text()
        test = path.name.endswith("_test.go")
        if cfg.tag_fields and not test:
            source, found = tag_struct_fields(source, cfg, lambda names: True)
            if any(type_name is None for type_name, _ in found):
                raise SystemExit(f"renamed field of an anonymous struct in {path.name}")
            tagged.extend(found)
        elif cfg.tag_fields:
            source, _ = tag_struct_fields(
                source, cfg, lambda names: all(n.lower() in keys for n in names))
        elif not test:
            source = preserve_projection_json_keys(source, cfg)
        if path.name in cfg.fixture_bound:
            source = add_fixture_tags(source, *cfg.fixture_bound[path.name])
        for alias, (file_names, import_path) in cfg.import_unaliases.items():
            if path.name in file_names:
                source = drop_import_alias(source, alias, import_path)
        sources[path] = (source, original)
    regression_file, regression = cfg.regression
    for name, names in cfg.reflect_names.items():
        found = {text[1:-1] for kind, text in lex((root / cfg.package_dir / name).read_text())
                 if kind == "string"}
        if not set(names) <= found or not set(names) <= set(mapping):
            raise SystemExit(f"reflection names not found in {name}")
    for path, (source, original) in sources.items():
        source = rewrite_go(source, mapping, cfg.reflect_names.get(path.name, ()))
        if path.name == regression_file:
            source = append_test(source, regression(tagged, mapping))
        sources[path] = (source, original)
    changed = []
    for path, (source, original) in sources.items():
        rel = path.relative_to(root).as_posix()
        if source != original:
            path.write_text(source)
        if source != original or rel in moves:
            changed.append(root / moves.get(rel, rel))
    for old, new in moves.items():
        os.rename(root / old, root / new)
    subprocess.run(["gofmt", "-w", *map(str, changed)], check=True)
    if cfg.golden:
        update_golden(root, cfg)
    update_baselines(root, mapping, moves, cfg)
    update_inventories(root, mapping, moves, cfg)
    changelog = root / "CHANGELOG.md"
    text = changelog.read_text()
    heading = "## Unreleased\n\n"
    if heading not in text:
        raise SystemExit("CHANGELOG.md has no Unreleased section")
    changelog.write_text(text.replace(heading, heading + changelog_entry(mapping, cfg) + "\n", 1))


def test_inventory_targets(tmp_path):
    row = {'owning_file': 'catalog/ebusstandard/ab76_standard03_candidate_views.go', 'kind': 'go_identifier', 'current_name': 'Ab76Interpret', 'target': 'Ab76Interpret', 'status': 'unchanged', 'identity': {'declaration': 'Ab76Interpret'}}
    (tmp_path / 'naming-inventory.json').write_text(json.dumps({'source_files': [row['owning_file']], 'rows': [row]}))
    (tmp_path / 'naming-consumers.json').write_text('{}')
    update_inventories(tmp_path, {'Ab76Interpret': 'Standard03CandidateViewInterpret'}, moved_paths(EBUSSTANDARD), EBUSSTANDARD)
    result = json.loads((tmp_path / 'naming-inventory.json').read_text())['rows'][0]
    assert result['current_name'] == result['target'] == 'Standard03CandidateViewInterpret'


def test_generated_projection_json_keys_test():
    canonical = (Path(__file__).resolve().parent.parent / EBUSSTANDARD.package_dir /
                 "canonical_names_test.go").read_text()
    parent, regression = canonical.split("func TestRenamedProjectionJSONKeys(", 1)
    assert append_test(parent, PROJECTION_JSON_KEYS_TEST) == canonical
    assert regression.count("map[string]string{") == 11


def test_run_generates_projection_json_keys_test(tmp_path, monkeypatch):
    canonical = (Path(__file__).resolve().parent.parent / EBUSSTANDARD.package_dir /
                 "canonical_names_test.go").read_text()
    parent = canonical.split("func TestRenamedProjectionJSONKeys(", 1)[0]
    package = tmp_path / EBUSSTANDARD.package_dir
    package.mkdir(parents=True)
    for old in moved_paths(EBUSSTANDARD):
        (tmp_path / old).write_text("package ebusstandard\n")
    target = package / "canonical_names_test.go"
    target.write_text(parent)
    (tmp_path / "CHANGELOG.md").write_text("## Unreleased\n\n")
    module = sys.modules[__name__]
    monkeypatch.setattr(module, "build_map", lambda root, cfg: {})
    monkeypatch.setattr(module, "add_fixture_tags", lambda source, *args: source)
    monkeypatch.setattr(module, "check_consumers", lambda *args: None)
    for name in ("update_golden", "update_baselines", "update_inventories"):
        monkeypatch.setattr(module, name, lambda *args: None)
    run(tmp_path, False, EBUSSTANDARD)
    assert target.read_text() == canonical


def test_projection_tags():
    source = 'type Ab76Selection struct {\n\tAb76Family string\n\tAb76Value int `json:"value"`\n}\n'
    tagged = preserve_projection_json_keys(source, EBUSSTANDARD)
    assert '`json:"Ab76Family"`' in tagged
    assert '`json:"value"`' in tagged
    assert preserve_projection_json_keys(tagged, EBUSSTANDARD) == tagged


def test_golden_projection_tags(tmp_path):
    import shutil

    package = tmp_path / EBUSSTANDARD.package_dir
    package.mkdir(parents=True)
    (package / "testdata").mkdir()
    (tmp_path / "go.mod").write_text("module snapshotexample\n\ngo 1.23\n")
    snapshot_test = (Path(__file__).resolve().parent.parent / EBUSSTANDARD.package_dir /
                     "provider_abi_snapshot_test.go")
    shutil.copyfile(snapshot_test, package / snapshot_test.name)
    (package / "internal/safetypolicy").mkdir(parents=True)
    source = ('package ebusstandard\n\ntype Ab76Selection struct {\n'
              '\tAb76Family string\n'
              '\tAb76Value int `json:"value"`\n}\n')
    mapping = {"Ab76Selection": "Standard03CandidateViewSelection",
               "Ab76Family": "Standard03CandidateViewFamily",
               "Ab76Value": "Standard03CandidateViewValue"}
    (package / "projection.go").write_text(
        rewrite_go(preserve_projection_json_keys(source, EBUSSTANDARD), mapping))
    golden = tmp_path / EBUSSTANDARD.golden
    golden.write_text('ebusstandard type Ab76Selection struct { Ab76Family string Ab76Value int `json:"value"` }\n')
    update_golden(tmp_path, EBUSSTANDARD)
    assert golden.read_text() == (
        'ebusstandard type Standard03CandidateViewSelection struct { '
        'Standard03CandidateViewFamily string `json:"Ab76Family"` '
        'Standard03CandidateViewValue int `json:"value"` }\n')


def test_struct_field_tags():
    source = (
        "type Ab16Key struct {\n\tAb16Family string\n\tAb16Group *uint8 `json:\"group\"`\n"
        "\tab16hidden int\n}\n\n"
        "func f() {\n\tvar fixture struct {\n\t\tAb19An19Limits []struct {\n"
        "\t\t\tAb19An19Label string\n\t\t\tAb19An19MinimumRaw, Ab19An19MaximumRaw string\n"
        "\t\t}\n\t}\n\tfor range []struct {\n\t\tAb21Name, Ab21Mode string\n\t}{} {\n\t}\n}\n")
    keys = {"ab16family", "ab16group", "ab19an19limits", "ab19an19label", "ab19an19minimumraw",
            "ab19an19maximumraw", "ab21name"}
    result, tagged = tag_struct_fields(
        source, VAILLANT_01_29, lambda names: all(n.lower() in keys for n in names))
    assert '\tAb16Family string `json:"Ab16Family"`\n' in result
    assert '`json:"group"`' in result and "ab16hidden int\n" in result
    assert '\t\t} `json:"Ab19An19Limits"`\n' in result
    assert ('\t\t\tAb19An19MinimumRaw string `json:"Ab19An19MinimumRaw"`\n'
            '\t\t\tAb19An19MaximumRaw string `json:"Ab19An19MaximumRaw"`\n') in result
    assert "\t\tAb21Name, Ab21Mode string\n" in result
    assert tagged == [("Ab16Key", ["Ab16Family"]),
                      (None, ["Ab19An19Label", "Ab19An19MinimumRaw", "Ab19An19MaximumRaw"]),
                      (None, ["Ab19An19Limits"])]
    again, _ = tag_struct_fields(result, VAILLANT_01_29, lambda names: True)
    assert again.count("`json:") == result.count("`json:") + 2


def test_field_json_keys_test():
    text = field_json_keys_test([("Ab16Key", ["Ab16Family", "Ab16Group"])],
                                {"Ab16Key": "B505OperatingRoleKey",
                                 "Ab16Family": "B505OperatingRoleFamily",
                                 "Ab16Group": "B505OperatingRoleGroup"})
    assert ('{B505OperatingRoleKey{}, map[string]string{"B505OperatingRoleFamily": '
            '"Ab16Family", "B505OperatingRoleGroup": "Ab16Group"}},') in text


def test_vaillant_names():
    cfg = VAILLANT_01_29
    assert new_name("Ab05Grade", cfg) == "SchemaGrade"
    assert new_name("ab13reflect", cfg) == "b507EnergyCandidateReflect"
    assert new_name("ab01ao01JSON", cfg) == "b505SentinelViewJSON"
    assert new_name("Ab16Role", cfg) == "B505OperatingRole"
    assert coded("Ab30Key", cfg) is None and coded("Ab21An21Value", cfg) is None
    assert len(moved_paths(cfg)) == 118
    assert set(code for code in cfg.stems) == {old.split("_", 1)[0] for old, _ in cfg.file_moves}


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__,
                                     formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("root", nargs="?", type=Path,
                        default=Path(__file__).resolve().parent.parent)
    parser.add_argument("--slice", required=True, choices=sorted(SLICES),
                        help="the rename slice to reproduce")
    parser.add_argument("--print-map", action="store_true",
                        help="print the identifier map and change nothing")
    arguments = parser.parse_args()
    run(arguments.root.resolve(), arguments.print_map, SLICES[arguments.slice])
    sys.exit(0)
