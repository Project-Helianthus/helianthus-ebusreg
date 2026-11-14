#!/usr/bin/env python3
"""Reproduce the rename of the work-item-coded strings of three catalog packages.

The work runs in slices; each slice is run once against the tree that still
holds its coded strings (--slice selects it):

- catalog-strings: the coded JSON tags, string values and error texts of
  catalog/ebusstandard and the first types of vaillant/system and
  vaillant/b555 (TOKENS, TEXTS, FILES below).
- system-b555-tags: the coded JSON tags of the remaining vaillant/system and
  vaillant/b555 types (TAG_TOKENS, TAG_FILES below).

The coded JSON tags, string values and error texts of catalog/ebusstandard,
vaillant/system and vaillant/b555 take behavior names. The tables below are
the whole change; nothing is computed from the source:

- TOKENS: per package, a coded JSON tag or value and its new string. A JSON
  tag or value keeps the part after the code in snake_case; the code is
  dropped because that part is already unique within its type.
- TEXTS: per production file, an error text or reason whose code prefix
  becomes the behavior of the file in words.
- FILES: the files each package table is applied to.

Production files, the provider ABI golden and router/vaillant_timer_access_test.go
take every token. Other test files take every token outside struct tags: the
tags of test-local structs name test vectors, not the catalog, and keep their
keys. A testdata fixture takes only the tokens listed for it, the keys and
values that decode into catalog types; KEEP lists the test-local keys of
unlisted files that equal a token. The rows of naming-inventory.json whose
current name is renamed follow one-to-one.

The system-b555-tags slice applies the same naming rule to struct tags only:

- TAG_TOKENS: per package, a coded JSON tag and its new key, the part after
  the code in snake_case.
- TAG_FILES: the production files whose json struct tags take the new keys;
  after the rewrite none of their tags carries a code.
- PIN_TESTS: test files whose TestRenamedFieldJSONKeys pins the production
  keys; only the quoted keys inside that function follow.
- FIXTURE_PATHS: per testdata fixture, the JSON paths ("*" is any array
  element) whose value decodes into a renamed type. Every key under such a
  path follows, matched without regard to case as encoding/json does; the
  keys of the test-local vector types around it stay.
- TAG_KEEP: the files that hold a token only as a test-local key; they stay.

There are no aliases and no dual acceptance: the old strings are gone. The
script runs once against the tree that still holds the coded strings and
stops before writing when a precondition fails. --check verifies that the
selected slice leaves no coded string in the places it rewrites. --changelog
prints the table of the CHANGELOG entry.
"""

import argparse
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

TOKENS = {
    "catalog/ebusstandard": {
        "Ab25An25CandidateMaximum": "candidate_maximum",
        "Ab25An25Conflict": "conflict",
        "Ab25An25ContractGrade": "contract_grade",
        "Ab25An25Data": "data",
        "Ab25An25Direction": "direction",
        "Ab25An25ExecutionGrade": "execution_grade",
        "Ab25An25Family": "family",
        "Ab25An25High": "high",
        "Ab25An25Key": "key",
        "Ab25An25LE": "le",
        "Ab25An25Low": "low",
        "Ab25An25MeaningGrade": "meaning_grade",
        "Ab25An25PhysicalOffsetsGrade": "physical_offsets_grade",
        "Ab25An25PhysicalSupportGrade": "physical_support_grade",
        "Ab25An25PositivePrefixRelevanceGrade": "positive_prefix_relevance_grade",
        "Ab25An25Profile": "profile",
        "Ab25An25Raw": "raw",
        "Ab25An25Service": "service",
        "Ab25An25State": "state",
        "Ab25An25UnresolvedTail": "unresolved_tail",
        "Ab25An25Value": "value",
        "Ab76Components": "components",
        "Ab76Default": "default",
        "Ab76Direction": "direction",
        "Ab76Encoding": "encoding",
        "Ab76Family": "family",
        "Ab76Function": "function",
        "Ab76ImplementationGrade": "implementation_grade",
        "Ab76LocalCounter": "local_counter",
        "Ab76LocalHours": "local_hours",
        "Ab76LocalMinute": "local_minute",
        "Ab76Maximum": "maximum",
        "Ab76MeaningGrade": "meaning_grade",
        "Ab76Name": "name",
        "Ab76NativeProperties": "native_properties",
        "Ab76Operation": "operation",
        "Ab76Raw": "raw",
        "Ab76Remainder": "remainder",
        "Ab76Selection": "selection",
        "Ab76Selector": "selector",
        "Ab76State": "state",
        "Ab76Value": "value",
        "Ab76Variant": "variant",
        "Ab76Width": "width",
        "Ab78APIGrade": "api_grade",
        "Ab78Complete": "complete",
        "Ab78Default": "default",
        "Ab78Direction": "direction",
        "Ab78Domain": "domain",
        "Ab78Family": "family",
        "Ab78Fields": "fields",
        "Ab78Interpretation": "interpretation",
        "Ab78Key": "key",
        "Ab78Label": "label",
        "Ab78Maximum": "maximum",
        "Ab78MeaningGrade": "meaning_grade",
        "Ab78Minimum": "minimum",
        "Ab78Name": "name",
        "Ab78NativeGrade": "native_grade",
        "Ab78Operation": "operation",
        "Ab78Primary": "primary",
        "Ab78Raw": "raw",
        "Ab78Remainder": "remainder",
        "Ab78Secondary": "secondary",
        "Ab78Selector": "selector",
        "Ab78State": "state",
        "Ab78Unknowns": "unknowns",
        "Ab78Value": "value",
        "Ab78Variant": "variant",
    },
    "vaillant/b555": {
        "Ab21An21Availability": "availability",
        "Ab21An21Context": "context",
        "Ab21An21Conversion": "conversion",
        "Ab21An21CountMaterialization": "count_materialization",
        "Ab21An21Execution": "execution",
        "Ab21An21Family": "family",
        "Ab21An21Freshness": "freshness",
        "Ab21An21IntermediateStatus": "intermediate_status",
        "Ab21An21Knowledge": "knowledge",
        "Ab21An21LocalCarrier": "local_carrier",
        "Ab21An21Low16": "low16",
        "Ab21An21OriginalFloat64": "original_float64",
        "Ab21An21OriginalInteger64": "original_integer64",
        "Ab21An21Path": "path",
        "Ab21An21Persistence": "persistence",
        "Ab21An21PhysicalQualification": "physical_qualification",
        "Ab21An21Presence": "presence",
        "Ab21An21RawClassification": "raw_classification",
        "Ab21An21RoundedClassification": "rounded_classification",
        "Ab21An21RoundedFloat64": "rounded_float64",
        "Ab21An21ScaledFloat64": "scaled_float64",
        "Ab21An21SoftwareVariant": "software_variant",
        "Ab21An21Status": "status",
        "Ab21An21TimerKind": "timer_kind",
        "Ab21An21TimerKindEffects": "timer_kind_effects",
        "Ab21An21Units": "units",
        "Ab21An21Unsigned32": "unsigned32",
        "Ab21An21WriteAdmission": "write_admission",
        "ab21an21-float64": "float64",
        "ab21an21-integer64": "integer64",
    },
    "vaillant/system": {
        "Ab08Address": "address",
        "Ab08Channel": "channel",
        "Ab08Day": "day",
        "Ab08Family": "family",
        "Ab08Identity": "identity",
        "Ab08Instance": "instance",
        "Ab08Raw": "raw",
        "Ab08Start": "start",
        "Ab08Stop": "stop",
        "Ab08SystemType": "system_type",
        "Ab09An09ContractGrade": "contract_grade",
        "Ab09An09Direction": "direction",
        "Ab09An09Family": "family",
        "Ab09An09Flow": "flow",
        "Ab09An09Identity": "identity",
        "Ab09An09MaximumVolume": "maximum_volume",
        "Ab09An09MinimumVolume": "minimum_volume",
        "Ab09An09NativeQualification": "native_qualification",
        "Ab09An09Operation": "operation",
        "Ab09An09Raw": "raw",
        "Ab09An09Selected": "selected",
        "Ab09An09Selection": "selection",
        "Ab09An09SemanticGrade": "semantic_grade",
        "Ab09An09State": "state",
        "Ab09An09Temperature": "temperature",
        "Ab09An09Unit": "unit",
        "Ab09An09Value": "value",
        "Ab09An09Variant": "variant",
        "Ab09An09Volume": "volume",
        "Ab24Address": "address",
        "Ab24Code": "code",
        "Ab24Family": "family",
        "Ab24Identity": "identity",
        "Ab24Instance": "instance",
        "Ab24OP09": "op09",
        "Ab24OP0B": "op0b",
        "Ab24Opcode": "opcode",
        "Ab24Payload": "payload",
        "Ab24Raw": "raw",
        "Ab24Success": "success",
        "Ab24SystemType": "system_type",
        "Ab24Values": "values",
        "Ab25APIRevision": "api_revision",
        "Ab25APIVersion": "api_version",
        "Ab25Codec": "codec",
        "Ab25Family": "family",
        "Ab25Group": "group",
        "Ab25HardwareIdentity": "hardware_identity",
        "Ab25Key": "key",
        "Ab25Opcode": "opcode",
        "Ab25RawDescription": "raw_description",
        "Ab25Register": "register",
        "Ab25RequestedInstance": "requested_instance",
        "Ab25SoftwareIdentity": "software_identity",
        "Ab25State": "state",
        "Ab25Width": "width",
        "Ab42APIRevision": "api_revision",
        "Ab42APIVersion": "api_version",
        "Ab42Capacity": "capacity",
        "Ab42Context": "context",
        "Ab42EligibilityGrade": "eligibility_grade",
        "Ab42Family": "family",
        "Ab42Hardware": "hardware",
        "Ab42HeatPumpCount": "heat_pump_count",
        "Ab42Identifier": "identifier",
        "Ab42InvalidReason": "invalid_reason",
        "Ab42LayoutGrade": "layout_grade",
        "Ab42Meaning": "meaning",
        "Ab42PresenceGrade": "presence_grade",
        "Ab42QualificationGrade": "qualification_grade",
        "Ab42Raw": "raw",
        "Ab42RepresentationGrade": "representation_grade",
        "Ab42Software": "software",
        "Ab42Valid": "valid",
        "Ab42Value": "value",
        "Ab44APIRevision": "api_revision",
        "Ab44APIVersion": "api_version",
        "Ab44Codec": "codec",
        "Ab44Context": "context",
        "Ab44Family": "family",
        "Ab44GG": "gg",
        "Ab44Grade": "grade",
        "Ab44Hardware": "hardware",
        "Ab44II": "ii",
        "Ab44ImplementationValid": "implementation_valid",
        "Ab44Limits": "limits",
        "Ab44Maximum": "maximum",
        "Ab44MaximumRaw": "maximum_raw",
        "Ab44MetadataOP": "metadata_op",
        "Ab44Minimum": "minimum",
        "Ab44MinimumRaw": "minimum_raw",
        "Ab44NumericValid": "numeric_valid",
        "Ab44OP": "op",
        "Ab44RR": "rr",
        "Ab44Raw": "raw",
        "Ab44Reason": "reason",
        "Ab44Software": "software",
        "Ab44Step": "step",
        "Ab44StepMeaning": "step_meaning",
        "Ab44StepRaw": "step_raw",
        "Ab44Width": "width",
        "ab07_parameter_configuration": "parameter_configuration",
        "ab07_raw_pair": "raw_pair",
        "ab07_slots": "slots",
        "ab07_start_minutes": "start_minutes",
        "ab07_stop_minutes": "stop_minutes",
        "ab07_unused": "unused",
        "ab09an09_absent": "absent",
        "ab09an09_candidate": "candidate",
        "ab09an09_malformed": "malformed",
        "ab09an09_opaque": "opaque",
        "ab09an09_request": "request",
        "ab09an09_response": "response",
        "ab09an09_unselected": "unselected",
        "ab44_bad_width": "bad_width",
        "ab44_codec_grade_required": "codec_grade_required",
        "ab44_codec_width_mismatch": "codec_width_mismatch",
        "ab44_date_time_raw_only": "date_time_raw_only",
        "ab44_echo_mismatch": "echo_mismatch",
        "ab44_incomplete_context": "incomplete_context",
        "ab44_malformed_length": "malformed_length",
        "ab44_negative_step": "negative_step",
        "ab44_nonfinite_limits": "nonfinite_limits",
        "ab44_numeric_limits": "numeric_limits",
        "ab44_reversed_limits": "reversed_limits",
        "ab44_unsupported_codec": "unsupported_codec",
        "ab44_unsupported_op": "unsupported_op",
    },
}

TEXTS = {
    "catalog/ebusstandard/standard03_counter_bodies.go": {
        "ab83 counter metadata: ": "standard03 counter body metadata: ",
    },
    "catalog/ebusstandard/standard07_residual_fields.go": {
        "ab28: duplicate key": "standard07 residual fields: duplicate key",
        "ab28: expected JSON array": "standard07 residual fields: expected JSON array",
        "ab28: incomplete identity": "standard07 residual fields: incomplete identity",
        "ab28: missing axis %s": "standard07 residual fields: missing axis %s",
        "ab28: missing key": "standard07 residual fields: missing key",
    },
    "catalog/ebusstandard/standard0f_candidates.go": {
        "ab81 embedded candidate metadata: ": "standard0f embedded candidate metadata: ",
    },
    "catalog/ebusstandard/standardff_body_states.go": {
        "ab54: duplicate JSON member": "standardff body states: duplicate JSON member",
        "ab54: duplicate key": "standardff body states: duplicate key",
        "ab54: expected JSON array": "standardff body states: expected JSON array",
        "ab54: incomplete key; use none for non-applicable string axes": "standardff body states: incomplete key; use none for non-applicable string axes",
        "ab54: unknown key axis": "standardff body states: unknown key axis",
    },
    "vaillant/system/b524_system_information.go": {
        "ab42 non-finite binary32": "b524 system information non-finite binary32",
        "ab42 requires exactly four normalized bytes": "b524 system information requires exactly four normalized bytes",
    },
    "vaillant/system/description_index.go": {
        "ab25 conflicting description records": "description index: conflicting description records",
    },
    "vaillant/system/event_plan.go": {
        "ab24: exactly seven raw values required": "event plan: exactly seven raw values required",
        "ab24: field outside raw u8 range": "event plan: field outside raw u8 range",
        "ab24: invalid address": "event plan: invalid address",
        "ab24: invalid edit opcode": "event plan: invalid edit opcode",
        "ab24: missing family": "event plan: missing family",
        "ab24: missing, failed or uncorrelated baseline %02x": "event plan: missing, failed or uncorrelated baseline %02x",
        "ab24: unsupported system type": "event plan: unsupported system type",
        "ab24: value outside raw u8 range": "event plan: value outside raw u8 range",
    },
    "vaillant/system/timer_plan.go": {
        "ab08 baseline identity mismatch: %w": "timer plan baseline identity mismatch: %w",
        "ab08 invalid explicit identity: %w": "timer plan invalid explicit identity: %w",
        "ab08 requires seven baseline bytes and three slots: %w": "timer plan requires seven baseline bytes and three slots: %w",
        "ab08 slot %d invalid raw ten-minute pair: %w": "timer plan slot %d invalid raw ten-minute pair: %w",
        "ab08 slot %d outside byte range: %w": "timer plan slot %d outside byte range: %w",
        "ab08 unsupported raw selector or non-zone instance: %w": "timer plan unsupported raw selector or non-zone instance: %w",
    },
}

FILES = {
    "catalog/ebusstandard": [
        "catalog/ebusstandard/canonical_names_test.go",
        "catalog/ebusstandard/standard03_candidate_views.go",
        "catalog/ebusstandard/standard05_candidate_views.go",
        "catalog/ebusstandard/standard09_body_views.go",
        "catalog/ebusstandard/testdata/provider_abi.golden.txt",
    ],
    "vaillant/b555": [
        "vaillant/b555/b555_extent_views_test.go",
        "vaillant/b555/b555_numeric_narrowing.go",
        "vaillant/b555/b555_numeric_narrowing_test.go",
        "vaillant/b555/testdata/ab21an21_b555_numeric_narrowing.json",
    ],
    "vaillant/system": [
        "router/vaillant_timer_access_test.go",
        "vaillant/system/ab07_timer_decode_test.go",
        "vaillant/system/b509_quantity_views.go",
        "vaillant/system/b509_quantity_views_test.go",
        "vaillant/system/b524_codebook_views_test.go",
        "vaillant/system/b524_description_spans.go",
        "vaillant/system/b524_description_spans_test.go",
        "vaillant/system/b524_system_information.go",
        "vaillant/system/b524_timer_access.go",
        "vaillant/system/description_index.go",
        "vaillant/system/event_plan.go",
        "vaillant/system/testdata/ab09an09_b509_quantity_views.json",
        "vaillant/system/timer_plan.go",
    ],
}

FIXTURE_TOKENS = {
    "vaillant/b555/testdata/ab21an21_b555_numeric_narrowing.json": [
        "ab21an21-float64",
        "ab21an21-integer64",
    ],
    "vaillant/system/testdata/ab09an09_b509_quantity_views.json": [
        "Ab09An09Direction",
        "Ab09An09Family",
        "Ab09An09Operation",
        "Ab09An09Variant",
        "ab09an09_absent",
        "ab09an09_candidate",
        "ab09an09_malformed",
        "ab09an09_opaque",
        "ab09an09_request",
        "ab09an09_response",
        "ab09an09_unselected",
    ],
}

KEEP = {
    "catalog/ebusstandard/standard03_candidate_views_test.go": ["Ab76Operation"],
    "catalog/ebusstandard/testdata/ab76_standard03_candidate_views.json": ["Ab76Operation"],
    "vaillant/system/b524_system_information_test.go": [
        "Ab42Capacity",
        "Ab42HeatPumpCount",
        "Ab42Identifier",
        "Ab42Meaning",
        "Ab42Raw",
        "Ab42Value",
    ],
}

INVENTORY = "naming-inventory.json"

STRUCT_TAG = re.compile(r'`(?:[a-z]+:"[^"`]*"\s*)+`')
EDGE = "(?<![A-Za-z0-9_]){}(?![A-Za-z0-9_])"
SKIP = {".git", ".tmp"}
HISTORY = {"CHANGELOG.md", "scripts/rename_catalog_names.py", "scripts/rename_catalog_strings.py"}


def pattern(tokens):
    alternatives = "|".join(re.escape(t) for t in sorted(tokens, key=len, reverse=True))
    return re.compile(EDGE.format("(" + alternatives + ")"))


def rewrite(rel, text, table):
    tokens = FIXTURE_TOKENS.get(rel, table)
    rx = pattern(tokens)
    replace = lambda m: table[m.group(1)]
    if not rel.endswith("_test.go") or rel == "router/vaillant_timer_access_test.go":
        return rx.subn(replace, text)
    parts, count, start = [], 0, 0
    for tag in STRUCT_TAG.finditer(text):
        piece, n = rx.subn(replace, text[start:tag.start()])
        parts += [piece, tag.group()]
        count += n
        start = tag.end()
    piece, n = rx.subn(replace, text[start:])
    return "".join(parts + [piece]), count + n


def rewrite_texts(rel, text):
    for old, new in TEXTS.get(rel, {}).items():
        quoted = '"' + old + '"'
        if text.count(quoted) != 1:
            raise SystemExit(f"{rel}: expected one {quoted}")
        text = text.replace(quoted, '"' + new + '"')
    return text


def repository_files(root):
    for path in sorted(root.rglob("*")):
        rel = path.relative_to(root).as_posix()
        if path.is_file() and not SKIP & set(rel.split("/")) and rel not in HISTORY:
            yield rel, path


def check_scope(root):
    listed = {rel for files in FILES.values() for rel in files} | set(TEXTS) | {INVENTORY}
    everything = pattern({t for table in TOKENS.values() for t in table})
    texts = re.compile("|".join(re.escape('"' + old + '"') for table in TEXTS.values() for old in table))
    for rel, path in repository_files(root):
        if rel in listed:
            continue
        try:
            text = path.read_text()
        except UnicodeDecodeError:
            continue
        for rx in (everything, texts):
            for hit in rx.finditer(text):
                if hit.group() not in KEEP.get(rel, ()):
                    raise SystemExit(f"{rel}: unlisted file holds {hit.group()}")


def plan(root):
    check_scope(root)
    changed = {}
    for package, files in FILES.items():
        for rel in files:
            text = (root / rel).read_text()
            new, count = rewrite(rel, text, TOKENS[package])
            if count == 0:
                raise SystemExit(f"{rel}: no coded string found")
            changed[rel] = new
    for rel in TEXTS:
        changed[rel] = rewrite_texts(rel, changed.get(rel) or (root / rel).read_text())
    changed[INVENTORY] = inventory(root)
    return changed


def inventory(root):
    data = json.loads((root / INVENTORY).read_text())
    count = 0
    for row in data["rows"]:
        package = row["owning_file"].rsplit("/", 1)[0]
        new = TOKENS.get(package, {}).get(row["current_name"])
        if new is not None:
            row["current_name"] = new
            count += 1
    if count == 0:
        raise SystemExit(f"{INVENTORY}: no coded row found")
    return json.dumps(data, indent=2, ensure_ascii=False) + "\n"


TAG_TOKENS = {
    "vaillant/b555": {
        "Ab52AvailableCount": "available_count",
        "Ab52DayCountGrade": "day_count_grade",
        "Ab52DayCounts": "day_counts",
        "Ab52DeclaredCount": "declared_count",
        "Ab52EighthCount": "eighth_count",
        "Ab52EighthCountGrade": "eighth_count_grade",
        "Ab52EndHour": "end_hour",
        "Ab52EndMinute": "end_minute",
        "Ab52Extent": "extent",
        "Ab52Incomplete": "incomplete",
        "Ab52Issues": "issues",
        "Ab52Raw": "raw",
        "Ab52SetpointGrade": "setpoint_grade",
        "Ab52SetpointPresent": "setpoint_present",
        "Ab52SetpointQualified": "setpoint_qualified",
        "Ab52SetpointRaw": "setpoint_raw",
        "Ab52SetpointTenths": "setpoint_tenths",
        "Ab52StartHour": "start_hour",
        "Ab52StartMinute": "start_minute",
        "Ab52Status": "status",
        "Ab52StatusAccepted": "status_accepted",
        "Ab52StatusGrade": "status_grade",
        "Ab52StatusUnavailable": "status_unavailable",
        "Ab52WeekComplete": "week_complete",
        "Ab52WeekValid": "week_valid",
    },
    "vaillant/system": {
        "Ab01An01APIRevision": "api_revision",
        "Ab01An01APIVersion": "api_version",
        "Ab01An01Alternative": "alternative",
        "Ab01An01Candidate": "candidate",
        "Ab01An01Cell": "cell",
        "Ab01An01Cells": "cells",
        "Ab01An01CombinedTimestamp": "combined_timestamp",
        "Ab01An01CompleteCandidate": "complete_candidate",
        "Ab01An01Components": "components",
        "Ab01An01Context": "context",
        "Ab01An01ContractGrade": "contract_grade",
        "Ab01An01EditValidation": "edit_validation",
        "Ab01An01Family": "family",
        "Ab01An01Grade": "grade",
        "Ab01An01Hardware": "hardware",
        "Ab01An01II": "ii",
        "Ab01An01NativeQualification": "native_qualification",
        "Ab01An01OP": "op",
        "Ab01An01Raw": "raw",
        "Ab01An01RequestedGG": "requested_gg",
        "Ab01An01RequestedRR": "requested_rr",
        "Ab01An01ResponseState": "response_state",
        "Ab01An01ReturnedGG": "returned_gg",
        "Ab01An01ReturnedRR": "returned_rr",
        "Ab01An01Software": "software",
        "Ab01An01State": "state",
        "Ab03An03APIRevision": "api_revision",
        "Ab03An03APIVersion": "api_version",
        "Ab03An03Applicable": "applicable",
        "Ab03An03CandidateLabel": "candidate_label",
        "Ab03An03Context": "context",
        "Ab03An03Family": "family",
        "Ab03An03GG": "gg",
        "Ab03An03Grade": "grade",
        "Ab03An03Hardware": "hardware",
        "Ab03An03II": "ii",
        "Ab03An03Integer": "integer",
        "Ab03An03MalformedWidth": "malformed_width",
        "Ab03An03NativeQualification": "native_qualification",
        "Ab03An03OP": "op",
        "Ab03An03RR": "rr",
        "Ab03An03Raw": "raw",
        "Ab03An03RequestedGG": "requested_gg",
        "Ab03An03RequestedRR": "requested_rr",
        "Ab03An03ResponseState": "response_state",
        "Ab03An03ReturnedGG": "returned_gg",
        "Ab03An03ReturnedRR": "returned_rr",
        "Ab03An03Software": "software",
        "Ab03An03Width": "width",
        "Ab03AvailabilityOnly": "availability_only",
        "Ab03Family": "family",
        "Ab03Group": "group",
        "Ab03Instances": "instances",
        "Ab03Maximum": "maximum",
        "Ab03Minimum": "minimum",
        "Ab03Opcode": "opcode",
        "Ab03PredicateRegister": "predicate_register",
        "Ab03RegisterCeiling": "register_ceiling",
        "Ab03Virtual": "virtual",
        "Ab07An07Context": "context",
        "Ab07An07Default": "default",
        "Ab07An07Direction": "direction",
        "Ab07An07Execution": "execution",
        "Ab07An07Factor": "factor",
        "Ab07An07FactorAlgorithm": "factor_algorithm",
        "Ab07An07Family": "family",
        "Ab07An07FamilyApplicability": "family_applicability",
        "Ab07An07Freshness": "freshness",
        "Ab07An07Interpretation": "interpretation",
        "Ab07An07Maximum": "maximum",
        "Ab07An07Minimum": "minimum",
        "Ab07An07NativeEchoRules": "native_echo_rules",
        "Ab07An07Operation": "operation",
        "Ab07An07PacketAssociation": "packet_association",
        "Ab07An07PeerEnforcement": "peer_enforcement",
        "Ab07An07PhysicalQualification": "physical_qualification",
        "Ab07An07Raw": "raw",
        "Ab07An07RawAddressEqual": "raw_address_equal",
        "Ab07An07RequestedAddress": "requested_address",
        "Ab07An07RequestedEditBoundsAdmitted": "requested_edit_bounds_admitted",
        "Ab07An07ReturnedAddress": "returned_address",
        "Ab07An07Selected": "selected",
        "Ab07An07SoftwareContract": "software_contract",
        "Ab07An07State": "state",
        "Ab07An07Step": "step",
        "Ab07An07Value": "value",
        "Ab07An07Variant": "variant",
        "Ab07An07Writability": "writability",
        "Ab11Bytes": "bytes",
        "Ab11Day": "day",
        "Ab11Grade": "grade",
        "Ab11Hour": "hour",
        "Ab11Interpretation": "interpretation",
        "Ab11Latin1": "latin1",
        "Ab11Minute": "minute",
        "Ab11Month": "month",
        "Ab11Qualified": "qualified",
        "Ab11Raw": "raw",
        "Ab11Second": "second",
        "Ab11Year": "year",
        "Ab22BindingZone": "binding_zone",
        "Ab22Context": "context",
        "Ab22CoolingOperatingMode": "cooling_operating_mode",
        "Ab22CoolingTemperatureRaw": "cooling_temperature_raw",
        "Ab22Family": "family",
        "Ab22Fields": "fields",
        "Ab22FrostProtection": "frost_protection",
        "Ab22HeatingOperatingMode": "heating_operating_mode",
        "Ab22HeatingTemperatureRaw": "heating_temperature_raw",
        "Ab22Opcode": "opcode",
        "Ab22Payload": "payload",
        "Ab22SpecialFunctionStatus": "special_function_status",
        "Ab22StatusInfo": "status_info",
        "Ab58APIRevision": "api_revision",
        "Ab58APIVersion": "api_version",
        "Ab58Codec": "codec",
        "Ab58CodecGrade": "codec_grade",
        "Ab58Context": "context",
        "Ab58DeclaredUnit": "declared_unit",
        "Ab58Delivery": "delivery",
        "Ab58Effects": "effects",
        "Ab58Execution": "execution",
        "Ab58Family": "family",
        "Ab58Float32": "float32",
        "Ab58Freshness": "freshness",
        "Ab58Grade": "grade",
        "Ab58Group": "group",
        "Ab58Hardware": "hardware",
        "Ab58Input": "input",
        "Ab58Instance": "instance",
        "Ab58Opcode": "opcode",
        "Ab58Persistence": "persistence",
        "Ab58Presence": "presence",
        "Ab58Present": "present",
        "Ab58Raw": "raw",
        "Ab58Register": "register",
        "Ab58ReplyKind": "reply_kind",
        "Ab58ReturnedGroup": "returned_group",
        "Ab58ReturnedRegister": "returned_register",
        "Ab58Software": "software",
        "Ab58State": "state",
        "Ab58Support": "support",
        "Ab58Uint16": "uint16",
        "Ab58Value": "value",
        "Ab60APIRevision": "api_revision",
        "Ab60APIVersion": "api_version",
        "Ab60Codec": "codec",
        "Ab60CodecGrade": "codec_grade",
        "Ab60Context": "context",
        "Ab60DeclaredUnit": "declared_unit",
        "Ab60Delivery": "delivery",
        "Ab60Effects": "effects",
        "Ab60Execution": "execution",
        "Ab60Family": "family",
        "Ab60Float32": "float32",
        "Ab60Freshness": "freshness",
        "Ab60Grade": "grade",
        "Ab60Group": "group",
        "Ab60Hardware": "hardware",
        "Ab60Input": "input",
        "Ab60Instance": "instance",
        "Ab60Opcode": "opcode",
        "Ab60Persistence": "persistence",
        "Ab60Presence": "presence",
        "Ab60Present": "present",
        "Ab60Raw": "raw",
        "Ab60Register": "register",
        "Ab60ReplyKind": "reply_kind",
        "Ab60ReturnedGroup": "returned_group",
        "Ab60ReturnedRegister": "returned_register",
        "Ab60Software": "software",
        "Ab60State": "state",
        "Ab60Support": "support",
        "Ab60Uint16": "uint16",
        "Ab60Value": "value",
        "Ab62APIRevision": "api_revision",
        "Ab62APIVersion": "api_version",
        "Ab62Codec": "codec",
        "Ab62CodecGrade": "codec_grade",
        "Ab62Context": "context",
        "Ab62DeclaredUnit": "declared_unit",
        "Ab62Family": "family",
        "Ab62Grade": "grade",
        "Ab62Group": "group",
        "Ab62Hardware": "hardware",
        "Ab62Input": "input",
        "Ab62Instance": "instance",
        "Ab62Opcode": "opcode",
        "Ab62Present": "present",
        "Ab62Raw": "raw",
        "Ab62Register": "register",
        "Ab62ReplyKind": "reply_kind",
        "Ab62ReturnedGroup": "returned_group",
        "Ab62ReturnedRegister": "returned_register",
        "Ab62Software": "software",
        "Ab62State": "state",
        "Ab62Value": "value",
        "Ab64APIRevision": "api_revision",
        "Ab64APIVersion": "api_version",
        "Ab64Codec": "codec",
        "Ab64CodecGrade": "codec_grade",
        "Ab64Context": "context",
        "Ab64Delivery": "delivery",
        "Ab64Effects": "effects",
        "Ab64Execution": "execution",
        "Ab64Family": "family",
        "Ab64Freshness": "freshness",
        "Ab64Grade": "grade",
        "Ab64Group": "group",
        "Ab64Hardware": "hardware",
        "Ab64Input": "input",
        "Ab64Instance": "instance",
        "Ab64InstanceRule": "instance_rule",
        "Ab64NativeSupport": "native_support",
        "Ab64Opcode": "opcode",
        "Ab64Persistence": "persistence",
        "Ab64Presence": "presence",
        "Ab64Present": "present",
        "Ab64Qualification": "qualification",
        "Ab64Raw": "raw",
        "Ab64Register": "register",
        "Ab64ReplyKind": "reply_kind",
        "Ab64ReturnedGroup": "returned_group",
        "Ab64ReturnedRegister": "returned_register",
        "Ab64Software": "software",
        "Ab64State": "state",
        "Ab64Unit": "unit",
        "Ab64Value": "value",
        "Ab66Applicability": "applicability",
        "Ab66Group": "group",
        "Ab66Hint": "hint",
        "Ab66Key": "key",
        "Ab66Max": "max",
        "Ab66Min": "min",
        "Ab66Opcode": "opcode",
        "Ab66Qualified": "qualified",
        "Ab66Register": "register",
        "Ab66Status": "status",
        "Ab66Step": "step",
        "Ab66Type": "type",
    },
}

TAG_FILES = {
    "vaillant/b555": [
        "vaillant/b555/b555_extent_views.go",
    ],
    "vaillant/system": [
        "vaillant/system/b509_constraint_word_views.go",
        "vaillant/system/b524_calendar_views.go",
        "vaillant/system/b524_circuit_scalars.go",
        "vaillant/system/b524_codebook_views.go",
        "vaillant/system/b524_dhw_ventilation_scalars.go",
        "vaillant/system/b524_system_scalars.go",
        "vaillant/system/b524_unscoped_hints.go",
        "vaillant/system/b524_zone_scalars.go",
        "vaillant/system/discovery_catalog.go",
        "vaillant/system/status.go",
        "vaillant/system/value_codec.go",
    ],
}

PIN_TESTS = {
    "vaillant/b555/b555_extent_views_test.go": "vaillant/b555",
    "vaillant/system/b524_codebook_views_test.go": "vaillant/system",
}

FIXTURE_PATHS = {
    "vaillant/system/testdata/ab03_discovery_windows.json": [
        ("ab03Windows", "*"),
    ],
    "vaillant/system/testdata/ab07an07_b509_constraint_word_views.json": [
        ("*", "Ab07An07Input"),
    ],
    "vaillant/system/testdata/ab66_b524_unscoped_hints.json": [
        ("ab66Keys", "*"),
        ("ab66Controls", "*", "Ab66Key"),
        ("ab66Controls", "*", "Ab66Hint"),
    ],
}

TAG_KEEP = [
    "vaillant/system/b509_constraint_word_views_test.go",
    "vaillant/system/b524_calendar_views_test.go",
    "vaillant/system/b524_circuit_scalars_test.go",
    "vaillant/system/b524_unscoped_hints_test.go",
    "vaillant/system/b524_zone_scalars_test.go",
    "vaillant/system/testdata/ab01an01_b524_calendar_views.json",
    "vaillant/system/testdata/ab03an03_b524_codebook_views.json",
    "vaillant/system/testdata/ab62_b524_circuit_scalars.json",
    "vaillant/system/testdata/ab64_b524_zone_scalars.json",
]

CODED_TAG = re.compile(r'json:"[Aa][Bb][0-9]{2}')
QUOTED = re.compile(r'"([A-Za-z0-9_]+)(?:,[a-z]+)*"')
PIN_FUNCTION = re.compile(r"\nfunc TestRenamedFieldJSONKeys\(.*?\n}\n", re.S)


def tag_package(rel):
    return rel.split("/testdata/")[0]


def rewrite_tags(text, table):
    tags = pattern(table)
    replace = lambda m: tags.sub(lambda t: table[t.group(1)], m.group())
    return STRUCT_TAG.subn(replace, text)


def rewrite_pins(text, table):
    pin = PIN_FUNCTION.search(text)
    if pin is None:
        raise SystemExit("TestRenamedFieldJSONKeys not found")
    quoted = re.compile('"(' + "|".join(re.escape(t) for t in table) + ')"')
    body, count = quoted.subn(lambda m: '"' + table[m.group(1)] + '"', pin.group())
    return text[:pin.start()] + body + text[pin.end():], count


def under(path, roots):
    def step(r, p):
        return isinstance(p, int) if r == "*" else r == p
    return any(len(path) >= len(root) and all(map(step, root, path)) for root in roots)


def rewrite_fixture(rel, text, table):
    lower = {old.lower(): new for old, new in table.items()}
    roots, stack, out, count, i = FIXTURE_PATHS[rel], [], [], 0, 0
    while i < len(text):
        c = text[i]
        if c == '"':
            j = i + 1
            while text[j] != '"':
                j += 2 if text[j] == "\\" else 1
            word, k = text[i + 1:j], j + 1
            while k < len(text) and text[k] in " \t\r\n":
                k += 1
            if k < len(text) and text[k] == ":" and stack and stack[-1][0] == "{":
                stack[-1][1] = word
                if under(tuple(frame[1] for frame in stack[:-1]), roots):
                    if word.lower() not in lower:
                        raise SystemExit(f"{rel}: key {word} under a renamed type is unlisted")
                    word, count = lower[word.lower()], count + 1
            out.append('"' + word + '"')
            i = j + 1
            continue
        if c in "{[":
            stack.append([c, None if c == "{" else 0])
        elif c in "}]":
            stack.pop()
        elif c == "," and stack[-1][0] == "[":
            stack[-1][1] += 1
        out.append(c)
        i += 1
    return "".join(out), count


def tag_scope(root):
    listed = {rel for files in TAG_FILES.values() for rel in files}
    listed |= set(PIN_TESTS) | set(FIXTURE_PATHS) | set(TAG_KEEP)
    lower = {t.lower() for table in TAG_TOKENS.values() for t in table}
    holders = set()
    for rel, path in repository_files(root):
        if not rel.endswith((".go", ".json", ".txt")):
            continue
        if any(m.group(1).lower() in lower for m in QUOTED.finditer(path.read_text())):
            holders.add(rel)
    if holders - listed:
        raise SystemExit(f"unlisted files hold a coded tag: {sorted(holders - listed)}")
    if set(TAG_KEEP) - holders:
        raise SystemExit(f"TAG_KEEP files hold no token: {sorted(set(TAG_KEEP) - holders)}")


def tag_plan(root):
    tag_scope(root)
    changed = {}
    for package, files in TAG_FILES.items():
        for rel in files:
            new, count = rewrite_tags((root / rel).read_text(), TAG_TOKENS[package])
            if count == 0 or CODED_TAG.search(new):
                raise SystemExit(f"{rel}: coded tags not fully covered by the table")
            changed[rel] = new
    for rel, package in PIN_TESTS.items():
        new, count = rewrite_pins((root / rel).read_text(), TAG_TOKENS[package])
        if count == 0:
            raise SystemExit(f"{rel}: no pinned coded key found")
        changed[rel] = new
    for rel in FIXTURE_PATHS:
        new, count = rewrite_fixture(rel, (root / rel).read_text(), TAG_TOKENS[tag_package(rel)])
        if count == 0:
            raise SystemExit(f"{rel}: no coded key under a renamed type")
        changed[rel] = new
    return changed


def tag_check(root):
    remaining = []
    for files in TAG_FILES.values():
        remaining += [rel for rel in files if CODED_TAG.search((root / rel).read_text())]
    for rel, package in PIN_TESTS.items():
        if rewrite_pins((root / rel).read_text(), TAG_TOKENS[package])[1]:
            remaining.append(rel)
    for rel in FIXTURE_PATHS:
        table = {new: new for new in TAG_TOKENS[tag_package(rel)].values()}
        try:
            rewrite_fixture(rel, (root / rel).read_text(), table)
        except SystemExit:
            remaining.append(rel)
    return remaining


def changelog_rows(tables):
    rows = {}
    for table in tables:
        for old, new in table.items():
            if rows.setdefault(old, new) != new:
                raise SystemExit(f"{old}: two new strings")
    lines = ["| Previous string | Current string |", "| --- | --- |"]
    lines += [f"| `{old}` | `{new}` |" for old, new in sorted(rows.items())]
    return "\n".join(lines)


SLICES = {
    "catalog-strings": (plan, lambda: changelog_rows(list(TOKENS.values()) + list(TEXTS.values()))),
    "system-b555-tags": (tag_plan, lambda: changelog_rows(TAG_TOKENS.values())),
}


def main():
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    parser.add_argument("--slice", required=True, choices=sorted(SLICES))
    parser.add_argument("--root", type=Path, default=ROOT)
    parser.add_argument("--dry-run", action="store_true")
    parser.add_argument("--check", action="store_true")
    parser.add_argument("--changelog", action="store_true")
    args = parser.parse_args()
    build, table = SLICES[args.slice]
    if args.changelog:
        print(table())
        return 0
    root = args.root.resolve()
    if args.check:
        if args.slice != "system-b555-tags":
            raise SystemExit("--check covers the system-b555-tags slice")
        remaining = tag_check(root)
        for rel in remaining:
            print(f"{rel}: coded string remains")
        return 1 if remaining else 0
    changed = build(root)
    for rel, text in sorted(changed.items()):
        print(rel)
        if not args.dry_run:
            (root / rel).write_text(text)
    return 0


if __name__ == "__main__":
    sys.exit(main())
