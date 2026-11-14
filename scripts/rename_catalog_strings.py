#!/usr/bin/env python3
"""Reproduce the rename of the work-item-coded strings of three catalog packages.

The work runs in slices; each slice is run once against the tree that still
holds its coded strings (--slice selects it):

- catalog-strings: the coded JSON tags, string values and error texts of
  catalog/ebusstandard and the first types of vaillant/system and
  vaillant/b555 (TOKENS, TEXTS, FILES below).
- system-b555-tags: the coded JSON tags of the remaining vaillant/system and
  vaillant/b555 types (TAG_TOKENS, TAG_FILES below).
- catalog-vaillant: the coded JSON tags, embedded data keys, string values
  and error texts of catalog/vaillant (VAILLANT_TOKENS, VAILLANT_VALUES,
  VAILLANT_TEXTS below).

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

The catalog-vaillant slice applies the same naming rule to catalog/vaillant:

- VAILLANT_TOKENS: a coded JSON tag, data key or value and its new string,
  the part after the code in snake_case.
- VAILLANT_VALUES: the tokens that are string values; a testdata fixture
  takes them only as whole JSON values.
- VAILLANT_TEXTS: per production file, an error text whose code becomes the
  behavior of the file in words.
- VAILLANT_KEPT: the row-identifier prefixes the loaders check; they stay
  with the row identifiers of the data.

Every string literal of the Go files of catalog/vaillant, tests included,
takes every token, so a test-local tag and the fixture it reads change
together. Every key of catalog/vaillant/data and catalog/vaillant/testdata
takes the token it equals without regard to case, as encoding/json matches
it; row identifiers and other values stay. The script refuses a key that
would collide inside one struct type or one JSON object.

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


VAILLANT_TOKENS = {
    "Ab05An05Calendar": "calendar",
    "Ab05An05Cell": "cell",
    "Ab05An05Complete": "complete",
    "Ab05An05Components": "components",
    "Ab05An05Context": "context",
    "Ab05An05Contract": "contract",
    "Ab05An05Direction": "direction",
    "Ab05An05Error": "error",
    "Ab05An05ErrorPresent": "error_present",
    "Ab05An05Family": "family",
    "Ab05An05Fields": "fields",
    "Ab05An05Marker": "marker",
    "Ab05An05MarkerPresent": "marker_present",
    "Ab05An05Meaning": "meaning",
    "Ab05An05NormalizedPresent": "normalized_present",
    "Ab05An05Operation": "operation",
    "Ab05An05Qualification": "qualification",
    "Ab05An05Raw": "raw",
    "Ab05An05RawPresent": "raw_present",
    "Ab05An05Remainder": "remainder",
    "Ab05An05RemainderPresent": "remainder_present",
    "Ab05An05State": "state",
    "Ab05An05Value": "value",
    "Ab05An05Variant": "variant",
    "Ab11An11AvailableExtent": "available_extent",
    "Ab11An11CandidateGrade": "candidate_grade",
    "Ab11An11CompleteLogicalExtent": "complete_logical_extent",
    "Ab11An11Direction": "direction",
    "Ab11An11Family": "family",
    "Ab11An11FanRunning": "fan_running",
    "Ab11An11Operation": "operation",
    "Ab11An11PhysicalQualification": "physical_qualification",
    "Ab11An11PressureDomain": "pressure_domain",
    "Ab11An11PressurePercent": "pressure_percent",
    "Ab11An11Raw": "raw",
    "Ab11An11RawFlag": "raw_flag",
    "Ab11An11RawOverlay": "raw_overlay",
    "Ab11An11RawPressure": "raw_pressure",
    "Ab11An11RequiredExtent": "required_extent",
    "Ab11An11RoomSetpointCelsius": "room_setpoint_celsius",
    "Ab11An11Selected": "selected",
    "Ab11An11Selection": "selection",
    "Ab11An11Variant": "variant",
    "Ab13An13Adjustment": "adjustment",
    "Ab13An13Admitted": "admitted",
    "Ab13An13Complete": "complete",
    "Ab13An13ContractGrade": "contract_grade",
    "Ab13An13Defaults": "defaults",
    "Ab13An13Execution": "execution",
    "Ab13An13FieldGrade": "field_grade",
    "Ab13An13Flags": "flags",
    "Ab13An13Flow": "flow",
    "Ab13An13Freshness": "freshness",
    "Ab13An13FrostProtectionActive": "frost_protection_active",
    "Ab13An13HeatDemand": "heat_demand",
    "Ab13An13InletConditionBits": "inlet_condition_bits",
    "Ab13An13InternalTestActive": "internal_test_active",
    "Ab13An13LegionellaDemand": "legionella_demand",
    "Ab13An13LoadingActive": "loading_active",
    "Ab13An13LoadingDisabled": "loading_disabled",
    "Ab13An13LoadingPumpFault": "loading_pump_fault",
    "Ab13An13Modulation": "modulation",
    "Ab13An13NamedDomain": "named_domain",
    "Ab13An13Positive": "positive",
    "Ab13An13Presence": "presence",
    "Ab13An13Qualification": "qualification",
    "Ab13An13Raw": "raw",
    "Ab13An13Ready": "ready",
    "Ab13An13Request": "request",
    "Ab13An13Reset": "reset",
    "Ab13An13ReturnTemperature": "return_temperature",
    "Ab13An13Scalars": "scalars",
    "Ab13An13SensorConditionBits": "sensor_condition_bits",
    "Ab13An13Setpoint": "setpoint",
    "Ab13An13TaskActive": "task_active",
    "Ab13An13Temperatures": "temperatures",
    "Ab13An13Unit": "unit",
    "Ab13An13UnitConflict": "unit_conflict",
    "Ab13An13Unmodelled": "unmodelled",
    "Ab13An13Value": "value",
    "Ab15An15Admitted": "admitted",
    "Ab15An15Candidate": "candidate",
    "Ab15An15Clock": "clock",
    "Ab15An15ContractGrade": "contract_grade",
    "Ab15An15Direction": "direction",
    "Ab15An15Execution": "execution",
    "Ab15An15Failed": "failed",
    "Ab15An15Family": "family",
    "Ab15An15FieldGrade": "field_grade",
    "Ab15An15Freshness": "freshness",
    "Ab15An15Messages": "messages",
    "Ab15An15Mode": "mode",
    "Ab15An15NativeSupport": "native_support",
    "Ab15An15Operation": "operation",
    "Ab15An15PhysicalAdaptation": "physical_adaptation",
    "Ab15An15PositionsGrade": "positions_grade",
    "Ab15An15Presence": "presence",
    "Ab15An15Raw": "raw",
    "Ab15An15Setpoint": "setpoint",
    "Ab15An15SetpointCelsius": "setpoint_celsius",
    "Ab15An15SpecialFunction": "special_function",
    "Ab15An15State": "state",
    "Ab15An15TablesGrade": "tables_grade",
    "Ab15An15Value": "value",
    "Ab15An15Variant": "variant",
    "Ab15An15Words": "words",
    "Ab16Alternative": "alternative",
    "Ab16Direction": "direction",
    "Ab16Family": "family",
    "Ab16Field": "field",
    "Ab16Group": "group",
    "Ab16InstanceScope": "instance_scope",
    "Ab16Operation": "operation",
    "Ab16Primary": "primary",
    "Ab16Register": "register",
    "Ab16Secondary": "secondary",
    "Ab16Variant": "variant",
    "Ab17An17Admitted": "admitted",
    "Ab17An17AvailableExtent": "available_extent",
    "Ab17An17CompleteTwoWord": "complete_two_word",
    "Ab17An17ContractGrade": "contract_grade",
    "Ab17An17DeclaredCount": "declared_count",
    "Ab17An17InterpretationGrade": "interpretation_grade",
    "Ab17An17Qualification": "qualification",
    "Ab17An17Raw": "raw",
    "Ab17An17Status": "status",
    "Ab17An17StatusAvailable": "status_available",
    "Ab17An17StatusMeaning": "status_meaning",
    "Ab17An17Trailing": "trailing",
    "Ab17An17Word0Available": "word0_available",
    "Ab17An17Word0Bits": "word0_bits",
    "Ab17An17Word0Finite": "word0_finite",
    "Ab17An17Word0NativeSign": "word0_native_sign",
    "Ab17An17Word0NativeUnit": "word0_native_unit",
    "Ab17An17Word0Raw": "word0_raw",
    "Ab17An17Word0Value": "word0_value",
    "Ab17An17Word1Available": "word1_available",
    "Ab17An17Word1Raw": "word1_raw",
    "Ab17An17Word1Scale": "word1_scale",
    "Ab17An17Word1Sign": "word1_sign",
    "Ab17An17Word1Type": "word1_type",
    "Ab19An19Admitted": "admitted",
    "Ab19An19Attempt": "attempt",
    "Ab19An19Available": "available",
    "Ab19An19ContractGrade": "contract_grade",
    "Ab19An19Execution": "execution",
    "Ab19An19Fresh": "fresh",
    "Ab19An19InterpretationGrade": "interpretation_grade",
    "Ab19An19LastKnown": "last_known",
    "Ab19An19Maximum": "maximum",
    "Ab19An19Minimum": "minimum",
    "Ab19An19NativeUnit": "native_unit",
    "Ab19An19Presence": "presence",
    "Ab19An19Prior": "prior",
    "Ab19An19Qualification": "qualification",
    "Ab19An19Raw": "raw",
    "Ab19An19State": "state",
    "Ab19An19Value": "value",
    "Ab23An23Chunks": "chunks",
    "Ab23An23Consumed": "consumed",
    "Ab23An23Procedure": "procedure",
    "Ab23An23Unplanned": "unplanned",
    "Ab80Direction": "direction",
    "Ab80Family": "family",
    "Ab80Interpretation": "interpretation",
    "Ab80Operation": "operation",
    "Ab80Present": "present",
    "Ab80Raw": "raw",
    "Ab80Value": "value",
    "Ab80Variant": "variant",
    "Ab82Blocks": "blocks",
    "Ab82Complete": "complete",
    "Ab82Conflict": "conflict",
    "Ab82Conflicts": "conflicts",
    "Ab82Context": "context",
    "Ab82Contributed": "contributed",
    "Ab82Direction": "direction",
    "Ab82Entries": "entries",
    "Ab82Extra": "extra",
    "Ab82Family": "family",
    "Ab82Grade": "grade",
    "Ab82Incomplete": "incomplete",
    "Ab82Index": "index",
    "Ab82Mask": "mask",
    "Ab82Missing": "missing",
    "Ab82Operation": "operation",
    "Ab82Primary": "primary",
    "Ab82Raw": "raw",
    "Ab82RawByte": "raw_byte",
    "Ab82Remainder": "remainder",
    "Ab82Secondary": "secondary",
    "Ab82Selector": "selector",
    "Ab82Signed": "signed",
    "Ab82Start": "start",
    "Ab82State": "state",
    "Ab82Values": "values",
    "Ab82Variant": "variant",
    "ab05_alternative": "alternative",
    "ab05_body_width": "body_width",
    "ab05_codec": "codec",
    "ab05_direction": "direction",
    "ab05_enum": "enum",
    "ab05_family": "family",
    "ab05_fields": "fields",
    "ab05_grade": "grade",
    "ab05_group": "group",
    "ab05_instance": "instance",
    "ab05_key": "key",
    "ab05_name": "name",
    "ab05_opcode": "opcode",
    "ab05_position": "position",
    "ab05_raw_fallback": "raw_fallback",
    "ab05_register": "register",
    "ab05_sentinel": "sentinel",
    "ab05_unit": "unit",
    "ab05_variant": "variant",
    "ab05_width": "width",
    "ab11an11-b511-family": "b511-family",
    "ab11an11-op00-fan": "op00-fan",
    "ab11an11-op00-reserved": "op00-reserved",
    "ab11an11-op07-base": "op07-base",
    "ab11an11-op07-extended-room": "op07-extended-room",
    "ab16Accepted": "accepted",
    "ab16Action": "action",
    "ab16BlockedEligibility": "blocked_eligibility",
    "ab16Boiler": "boiler",
    "ab16BranchOrder": "branch_order",
    "ab16CircuitClass": "circuit_class",
    "ab16DisabledValue": "disabled_value",
    "ab16EnabledEligibility": "enabled_eligibility",
    "ab16EnabledValue": "enabled_value",
    "ab16Encoding": "encoding",
    "ab16ExceptionRole": "exception_role",
    "ab16Intent": "intent",
    "ab16Mode": "mode",
    "ab16Nonzero": "nonzero",
    "ab16NumericBranch": "numeric_branch",
    "ab16OtherModes": "other_modes",
    "ab16Period": "period",
    "ab16Persistent": "persistent",
    "ab16ProgrammingMode": "programming_mode",
    "ab16Quantity": "quantity",
    "ab16Relations": "relations",
    "ab16Replacement": "replacement",
    "ab16ReplacementAction": "replacement_action",
    "ab16ReplacementRaw": "replacement_raw",
    "ab16RequestValue": "request_value",
    "ab16Role": "role",
    "ab16SelfHeating": "self_heating",
    "ab16Tank": "tank",
    "ab16Unknown": "unknown",
    "ab16UserCorrection": "user_correction",
    "ab16Zero": "zero",
    "ab17_activation_order": "activation_order",
    "ab17_completed_stop": "completed_stop",
    "ab17_count": "count",
    "ab17_delivery_guarantee": "delivery_guarantee",
    "ab17_extent_bytes": "extent_bytes",
    "ab17_extent_class": "extent_class",
    "ab17_field_position": "field_position",
    "ab17_finite_maximum": "finite_maximum",
    "ab17_finite_meaning": "finite_meaning",
    "ab17_finite_minimum": "finite_minimum",
    "ab17_id": "id",
    "ab17_maximum_bytes": "maximum_bytes",
    "ab17_minimum_bytes": "minimum_bytes",
    "ab17_raw_form": "raw_form",
    "ab17_related_operations": "related_operations",
    "ab17_retry_count": "retry_count",
    "ab17_timeout": "timeout",
    "ab17_unknown": "unknown",
    "ab17_zero_meaning": "zero_meaning",
    "ab21_actual_acknowledgement": "actual_acknowledgement",
    "ab21_block_number_in_data": "block_number_in_data",
    "ab21_crc_query_operation": "crc_query_operation",
    "ab21_data": "data",
    "ab21_delivery": "delivery",
    "ab21_grade": "grade",
    "ab21_id": "id",
    "ab21_initiator_relation": "initiator_relation",
    "ab21_metadata_presence": "metadata_presence",
    "ab21_minimum_length": "minimum_length",
    "ab21_native_contract": "native_contract",
    "ab21_native_support": "native_support",
    "ab21_peer_capacity": "peer_capacity",
    "ab21_persistence": "persistence",
    "ab21_physical_effects": "physical_effects",
    "ab21_preserve_raw": "preserve_raw",
    "ab21_qualification": "qualification",
    "ab21_raw_form": "raw_form",
    "ab21_reply_policy": "reply_policy",
    "ab21_reply_recipient": "reply_recipient",
    "ab21_representation": "representation",
    "ab21_routing": "routing",
    "ab21_service": "service",
    "ab21_storage_effects": "storage_effects",
    "ab21_unknown": "unknown",
    "ab21_unresolved_bytes": "unresolved_bytes",
    "ab23an23_grouped_1024": "grouped_1024",
    "ab23an23_single_205": "single_205",
    "ab34_alternative": "alternative",
    "ab34_candidate_grade": "candidate_grade",
    "ab34_context_only": "context_only",
    "ab34_cooldown_event": "cooldown_event",
    "ab34_delivery": "delivery",
    "ab34_direction": "direction",
    "ab34_family": "family",
    "ab34_field": "field",
    "ab34_group": "group",
    "ab34_id": "id",
    "ab34_identity": "identity",
    "ab34_instance_scope": "instance_scope",
    "ab34_native_support": "native_support",
    "ab34_operation": "operation",
    "ab34_payload": "payload",
    "ab34_persistence": "persistence",
    "ab34_physical_effects": "physical_effects",
    "ab34_policy": "policy",
    "ab34_positive_status_service": "positive_status_service",
    "ab34_primary": "primary",
    "ab34_raw_form": "raw_form",
    "ab34_register": "register",
    "ab34_representation_grade": "representation_grade",
    "ab34_request_service": "request_service",
    "ab34_retries": "retries",
    "ab34_return_codec": "return_codec",
    "ab34_return_data_service": "return_data_service",
    "ab34_secondary": "secondary",
    "ab34_segmentation": "segmentation",
    "ab34_sequence": "sequence",
    "ab34_service": "service",
    "ab34_service_role": "service_role",
    "ab34_state_order_beyond_wait": "state_order_beyond_wait",
    "ab34_status_is_return_data": "status_is_return_data",
    "ab34_success_cooldown_required": "success_cooldown_required",
    "ab34_synchronization_bytes": "synchronization_bytes",
    "ab34_synchronization_required": "synchronization_required",
    "ab34_synchronization_token": "synchronization_token",
    "ab34_timeout_cooldown_min_ms": "timeout_cooldown_min_ms",
    "ab34_timeout_duration": "timeout_duration",
    "ab34_transport_adaptation": "transport_adaptation",
    "ab34_unknown": "unknown",
    "ab34_variant": "variant",
}

VAILLANT_VALUES = [
    "ab11an11-b511-family",
    "ab11an11-op00-fan",
    "ab11an11-op00-reserved",
    "ab11an11-op07-base",
    "ab11an11-op07-extended-room",
    "ab23an23_grouped_1024",
    "ab23an23_single_205",
]

VAILLANT_TEXTS = {
    "catalog/vaillant/b501_legacy_sequence.go": {
        "ab34 array required": "b501 legacy sequence: array required",
        "ab34 duplicate identity": "b501 legacy sequence: duplicate identity",
        "ab34 incomplete identity": "b501 legacy sequence: incomplete identity",
        "ab34 missing identity axis %s": "b501 legacy sequence: missing identity axis %s",
    },
    "catalog/vaillant/b503_text_categories.go": {
        "ab96: duplicate embedded key": "b503 text categories: duplicate embedded key",
        "ab96: empty embedded catalog": "b503 text categories: empty embedded catalog",
        "ab96: invalid embedded entry %d": "b503 text categories: invalid embedded entry %d",
        "ab96: invalid embedded width or position": "b503 text categories: invalid embedded width or position",
        "ab96: missing embedded primary in entry %d": "b503 text categories: missing embedded primary in entry %d",
        "ab96: trailing embedded data: %v": "b503 text categories: trailing embedded data: %v",
    },
    "catalog/vaillant/b504_circuit_roles.go": {
        "ab15: duplicate key": "b504 circuit roles: duplicate key",
        "ab15: empty catalog": "b504 circuit roles: empty catalog",
        "ab15: incomplete evidence boundary": "b504 circuit roles: incomplete evidence boundary",
        "ab15: incomplete identity": "b504 circuit roles: incomplete identity",
        "ab15: invalid or duplicate id": "b504 circuit roles: invalid or duplicate id",
        "ab15: missing identity axis %s": "b504 circuit roles: missing identity axis %s",
    },
    "catalog/vaillant/b505_operating_roles.go": {
        "ab16: duplicate identity": "b505 operating roles: duplicate identity",
        "ab16: expected array": "b505 operating roles: expected array",
        "ab16: incomplete string axes": "b505 operating roles: incomplete string axes",
        "ab16: missing identity": "b505 operating roles: missing identity",
    },
    "catalog/vaillant/b505_sentinel_views.go": {
        "ab01ao01: %s must be a byte array": "b505 sentinel views: %s must be a byte array",
        "ab01ao01: %s must be a string": "b505 sentinel views: %s must be a string",
        "ab01ao01: %s must be an integer": "b505 sentinel views: %s must be an integer",
        "ab01ao01: context must be an object": "b505 sentinel views: context must be an object",
        "ab01ao01: duplicate key %q": "b505 sentinel views: duplicate key %q",
        "ab01ao01: fields must be an object": "b505 sentinel views: fields must be an object",
        "ab01ao01: input must be an object": "b505 sentinel views: input must be an object",
        "ab01ao01: invalid byte": "b505 sentinel views: invalid byte",
        "ab01ao01: invalid object key": "b505 sentinel views: invalid object key",
        "ab01ao01: no arguments are defined for this view": "b505 sentinel views: no arguments are defined for this view",
        "ab01ao01: trailing JSON": "b505 sentinel views: trailing JSON",
        "ab01ao01: unexpected delimiter": "b505 sentinel views: unexpected delimiter",
        "ab01ao01: unknown input key %q": "b505 sentinel views: unknown input key %q",
    },
    "catalog/vaillant/b506_rf_sequence.go": {
        "ab17: duplicate key": "b506 rf sequence: duplicate key",
        "ab17: expected array": "b506 rf sequence: expected array",
        "ab17: incomplete metadata": "b506 rf sequence: incomplete metadata",
        "ab17: missing axis %s": "b506 rf sequence: missing axis %s",
        "ab17: null axis %s": "b506 rf sequence: null axis %s",
    },
    "catalog/vaillant/b507_energy_candidates.go": {
        "ab13: duplicate embedded key at %d": "b507 energy candidates: duplicate embedded key at %d",
        "ab13: empty embedded catalog": "b507 energy candidates: empty embedded catalog",
        "ab13: invalid divisor at %d": "b507 energy candidates: invalid divisor at %d",
        "ab13: invalid domain at %d": "b507 energy candidates: invalid domain at %d",
        "ab13: invalid embedded entry %d": "b507 energy candidates: invalid embedded entry %d",
        "ab13: trailing embedded data": "b507 energy candidates: trailing embedded data",
    },
    "catalog/vaillant/b509_boiler_codec_dictionary.go": {
        "ab68: duplicate key at row %d": "b509 boiler codec dictionary: duplicate key at row %d",
        "ab68: empty dictionary": "b509 boiler codec dictionary: empty dictionary",
        "ab68: incomplete declaration at row %d": "b509 boiler codec dictionary: incomplete declaration at row %d",
        "ab68: incomplete key at row %d": "b509 boiler codec dictionary: incomplete key at row %d",
        "ab68: invalid %s grade at row %d": "b509 boiler codec dictionary: invalid %s grade at row %d",
        "ab68: missing metadata grade at row %d": "b509 boiler codec dictionary: missing metadata grade at row %d",
        "ab68: selector byte out of range at row %d": "b509 boiler codec dictionary: selector byte out of range at row %d",
        "ab68: trailing JSON: %v": "b509 boiler codec dictionary: trailing JSON: %v",
    },
    "catalog/vaillant/b509_heatpump_codec_dictionary.go": {
        "ab70: duplicate key": "b509 heatpump codec dictionary: duplicate key",
        "ab70: empty dictionary": "b509 heatpump codec dictionary: empty dictionary",
        "ab70: incomplete or invalid entry": "b509 heatpump codec dictionary: incomplete or invalid entry",
        "ab70: trailing data": "b509 heatpump codec dictionary: trailing data",
    },
    "catalog/vaillant/b509_identity_shapes.go": {
        "ab46: duplicate key at row %d": "b509 identity shapes: duplicate key at row %d",
        "ab46: expected an array": "b509 identity shapes: expected an array",
        "ab46: incomplete key at row %d": "b509 identity shapes: incomplete key at row %d",
        "ab46: invalid key": "b509 identity shapes: invalid key",
        "ab46: row %d must be an object": "b509 identity shapes: row %d must be an object",
    },
    "catalog/vaillant/b509_op59_quantities.go": {
        "ab14: duplicate key": "b509 op59 quantities: duplicate key",
        "ab14: expected five finite entries": "b509 op59 quantities: expected five finite entries",
        "ab14: incomplete quantity codec": "b509 op59 quantities: incomplete quantity codec",
        "ab14: invalid codec %d": "b509 op59 quantities: invalid codec %d",
        "ab14: invalid entry %d": "b509 op59 quantities: invalid entry %d",
        "ab14: raw limit has unprovided metadata": "b509 op59 quantities: raw limit has unprovided metadata",
        "ab14: trailing embedded data": "b509 op59 quantities: trailing embedded data",
    },
    "catalog/vaillant/b509_symbolic_annotations.go": {
        "ab40: duplicate key": "b509 symbolic annotations: duplicate key",
        "ab40: empty string axis; use none when inapplicable": "b509 symbolic annotations: empty string axis; use none when inapplicable",
        "ab40: expected array": "b509 symbolic annotations: expected array",
        "ab40: missing axis %s": "b509 symbolic annotations: missing axis %s",
        "ab40: null optional axis %s": "b509 symbolic annotations: null optional axis %s",
    },
    "catalog/vaillant/b50e_rule_value_rows.go": {
        "ab74 duplicate row key": "b50e rule value rows: duplicate row key",
        "ab74 expected nineteen rows": "b50e rule value rows: expected nineteen rows",
        "ab74 incomplete qualification": "b50e rule value rows: incomplete qualification",
        "ab74 invalid field declaration": "b50e rule value rows: invalid field declaration",
        "ab74 invalid row identity": "b50e rule value rows: invalid row identity",
        "ab74 invalid value base": "b50e rule value rows: invalid value base",
        "ab74 missing explicit metadata field": "b50e rule value rows: missing explicit metadata field",
        "ab74 missing known width": "b50e rule value rows: missing known width",
        "ab74 missing or extra companion": "b50e rule value rows: missing or extra companion",
        "ab74 opaque width must remain absent": "b50e rule value rows: opaque width must remain absent",
        "ab74 trailing JSON: %v": "b50e rule value rows: trailing JSON: %v",
    },
    "catalog/vaillant/b50f_status_spans.go": {
        "ab12: duplicate key": "b50f status spans: duplicate key",
        "ab12: incomplete bit partition": "b50f status spans: incomplete bit partition",
        "ab12: incomplete finite catalog": "b50f status spans: incomplete finite catalog",
        "ab12: inconsistent raw byte": "b50f status spans: inconsistent raw byte",
        "ab12: invalid bit": "b50f status spans: invalid bit",
        "ab12: invalid container/raw partition": "b50f status spans: invalid container/raw partition",
        "ab12: invalid entry %d": "b50f status spans: invalid entry %d",
        "ab12: trailing data": "b50f status spans: trailing data",
    },
    "catalog/vaillant/b510_operating_defaults.go": {
        "ab48 catalog: duplicate key": "b510 operating defaults catalog: duplicate key",
        "ab48 catalog: empty key axis; use none for non-applicable strings": "b510 operating defaults catalog: empty key axis; use none for non-applicable strings",
        "ab48 catalog: expected array": "b510 operating defaults catalog: expected array",
        "ab48 catalog: missing key axis %s": "b510 operating defaults catalog: missing key axis %s",
    },
    "catalog/vaillant/b511_class_candidates.go": {
        "ab50: duplicate key": "b511 class candidates: duplicate key",
        "ab50: expected array": "b511 class candidates: expected array",
        "ab50: incomplete key": "b511 class candidates: incomplete key",
        "ab50: missing key": "b511 class candidates: missing key",
        "ab50: missing primary": "b511 class candidates: missing primary",
    },
    "catalog/vaillant/b512_dhw_stepper.go": {
        "ab97 duplicate key": "b512 dhw stepper: duplicate key",
        "ab97 empty catalog": "b512 dhw stepper: empty catalog",
        "ab97 invalid bit": "b512 dhw stepper: invalid bit",
        "ab97 invalid divisor": "b512 dhw stepper: invalid divisor",
        "ab97 invalid entry %d": "b512 dhw stepper: invalid entry %d",
        "ab97 trailing data": "b512 dhw stepper: trailing data",
    },
    "catalog/vaillant/b512_feedback_fallback.go": {
        "ab18: duplicate identity": "b512 feedback fallback: duplicate identity",
        "ab18: expected JSON array": "b512 feedback fallback: expected JSON array",
        "ab18: incomplete identity": "b512 feedback fallback: incomplete identity",
        "ab18: missing identity axis %s": "b512 feedback fallback: missing identity axis %s",
        "ab18: unexpected identity axes": "b512 feedback fallback: unexpected identity axes",
    },
    "catalog/vaillant/b512_feedback_views.go": {
        "ab15ao15 duplicate or invalid key": "b512 feedback views: duplicate or invalid key",
        "ab15ao15 expected object": "b512 feedback views: expected object",
        "ab15ao15 field must be byte array": "b512 feedback views: field must be byte array",
        "ab15ao15 invalid argument": "b512 feedback views: invalid argument",
        "ab15ao15 invalid byte": "b512 feedback views: invalid byte",
        "ab15ao15 null byte": "b512 feedback views: null byte",
        "ab15ao15 null identity": "b512 feedback views: null identity",
        "ab15ao15 null selector": "b512 feedback views: null selector",
        "ab15ao15 trailing JSON": "b512 feedback views: trailing JSON",
        "ab15ao15 unexpected delimiter": "b512 feedback views: unexpected delimiter",
        "ab15ao15 unknown input key": "b512 feedback views: unknown input key",
    },
    "catalog/vaillant/b512_heatpump_io.go": {
        "ab19 catalog: %w": "b512 heatpump io catalog: %w",
        "ab19 catalog: duplicate key": "b512 heatpump io catalog: duplicate key",
        "ab19 catalog: empty array": "b512 heatpump io catalog: empty array",
        "ab19 catalog: incomplete identity": "b512 heatpump io catalog: incomplete identity",
        "ab19 catalog: invalid ID": "b512 heatpump io catalog: invalid ID",
        "ab19 catalog: missing %s": "b512 heatpump io catalog: missing %s",
    },
    "catalog/vaillant/b513_symbolic_dictionary.go": {
        "ab38 dictionary: %w": "b513 symbolic dictionary: %w",
        "ab38 dictionary: duplicate key": "b513 symbolic dictionary: duplicate key",
        "ab38 dictionary: expected nonempty array": "b513 symbolic dictionary: expected nonempty array",
        "ab38 row: %w": "b513 symbolic dictionary row: %w",
        "ab38 row: duplicate member": "b513 symbolic dictionary row: duplicate member",
        "ab38 row: expected object": "b513 symbolic dictionary row: expected object",
        "ab38 row: incomplete dictionary identity": "b513 symbolic dictionary row: incomplete dictionary identity",
        "ab38 row: invalid candidate boundary or raw form": "b513 symbolic dictionary row: invalid candidate boundary or raw form",
        "ab38 row: missing member %s": "b513 symbolic dictionary row: missing member %s",
        "ab38 row: null required member": "b513 symbolic dictionary row: null required member",
        "ab38 row: trailing data": "b513 symbolic dictionary row: trailing data",
    },
    "catalog/vaillant/b516_directional_types.go": {
        "ab20 catalog: %w": "b516 directional types catalog: %w",
        "ab20 catalog: duplicate key": "b516 directional types catalog: duplicate key",
        "ab20 catalog: empty array": "b516 directional types catalog: empty array",
        "ab20 catalog: incomplete identity or duplicate id": "b516 directional types catalog: incomplete identity or duplicate id",
        "ab20 catalog: missing primary": "b516 directional types catalog: missing primary",
    },
    "catalog/vaillant/b516_residual_fields.go": {
        "ab51: duplicate key": "b516 residual fields: duplicate key",
        "ab51: empty string axis; use none when not applicable": "b516 residual fields: empty string axis; use none when not applicable",
        "ab51: expected array": "b516 residual fields: expected array",
        "ab51: missing key axis %s": "b516 residual fields: missing key axis %s",
    },
    "catalog/vaillant/b516_statistics_directions.go": {
        "ab72 catalog: %w": "b516 statistics directions catalog: %w",
        "ab72 catalog: duplicate key": "b516 statistics directions catalog: duplicate key",
        "ab72 catalog: empty": "b516 statistics directions catalog: empty",
        "ab72 catalog: incomplete key": "b516 statistics directions catalog: incomplete key",
        "ab72 catalog: invalid bit": "b516 statistics directions catalog: invalid bit",
        "ab72 catalog: invalid evidence": "b516 statistics directions catalog: invalid evidence",
        "ab72 catalog: invalid qualification": "b516 statistics directions catalog: invalid qualification",
        "ab72 catalog: invalid width": "b516 statistics directions catalog: invalid width",
        "ab72 catalog: trailing input": "b516 statistics directions catalog: trailing input",
    },
    "catalog/vaillant/b519_reply_policy.go": {
        "ab21: duplicate identity at %d": "b519 reply policy: duplicate identity at %d",
        "ab21: empty catalog": "b519 reply policy: empty catalog",
        "ab21: incomplete identity at %d": "b519 reply policy: incomplete identity at %d",
        "ab21: invalid identity at %d": "b519 reply policy: invalid identity at %d",
        "ab21: missing axis %s": "b519 reply policy: missing axis %s",
    },
    "catalog/vaillant/b51a_diagnostic_variants.go": {
        "ab98 catalog: %w": "b51a diagnostic variants catalog: %w",
        "ab98 catalog: default outside domain": "b51a diagnostic variants catalog: default outside domain",
        "ab98 catalog: duplicate key": "b51a diagnostic variants catalog: duplicate key",
        "ab98 catalog: empty entries": "b51a diagnostic variants catalog: empty entries",
        "ab98 catalog: incomplete key": "b51a diagnostic variants catalog: incomplete key",
        "ab98 catalog: invalid alternative": "b51a diagnostic variants catalog: invalid alternative",
        "ab98 catalog: invalid bit": "b51a diagnostic variants catalog: invalid bit",
        "ab98 catalog: invalid byte order": "b51a diagnostic variants catalog: invalid byte order",
        "ab98 catalog: invalid condition": "b51a diagnostic variants catalog: invalid condition",
        "ab98 catalog: invalid domain": "b51a diagnostic variants catalog: invalid domain",
        "ab98 catalog: invalid evidence boundary": "b51a diagnostic variants catalog: invalid evidence boundary",
        "ab98 catalog: invalid width": "b51a diagnostic variants catalog: invalid width",
        "ab98 catalog: trailing data": "b51a diagnostic variants catalog: trailing data",
    },
    "catalog/vaillant/b51b_local_budget_crc.go": {
        "ab23an23: negative progress or length": "b51b local budget crc: negative progress or length",
        "ab23an23: progress overflow": "b51b local budget crc: progress overflow",
        "ab23an23: unknown procedure": "b51b local budget crc: unknown procedure",
    },
    "catalog/vaillant/b522_demand_views.go": {
        "ab25ao25: argument %q is not a boolean or integer": "b522 demand views: argument %q is not a boolean or integer",
        "ab25ao25: arguments is not an object": "b522 demand views: arguments is not an object",
        "ab25ao25: context missing or not an object": "b522 demand views: context missing or not an object",
        "ab25ao25: context requires integer primary, secondary, operation and string direction, view": "b522 demand views: context requires integer primary, secondary, operation and string direction, view",
        "ab25ao25: duplicate key %q": "b522 demand views: duplicate key %q",
        "ab25ao25: field %q has a non-integer element": "b522 demand views: field %q has a non-integer element",
        "ab25ao25: field %q has an element outside 0..255": "b522 demand views: field %q has an element outside 0..255",
        "ab25ao25: field %q is not an array": "b522 demand views: field %q is not an array",
        "ab25ao25: fields missing or not an object": "b522 demand views: fields missing or not an object",
        "ab25ao25: input is not an object": "b522 demand views: input is not an object",
        "ab25ao25: invalid input: %w": "b522 demand views: invalid input: %w",
        "ab25ao25: nesting too deep": "b522 demand views: nesting too deep",
        "ab25ao25: object key is not a string": "b522 demand views: object key is not a string",
        "ab25ao25: trailing data after input": "b522 demand views: trailing data after input",
    },
    "catalog/vaillant/b523_config_mixer.go": {
        "ab99: duplicate key": "b523 config mixer: duplicate key",
        "ab99: empty catalog": "b523 config mixer: empty catalog",
        "ab99: invalid entry %d": "b523 config mixer: invalid entry %d",
        "ab99: invalid key coordinates in entry %d": "b523 config mixer: invalid key coordinates in entry %d",
        "ab99: trailing data": "b523 config mixer: trailing data",
    },
    "catalog/vaillant/b524_component_selectors.go": {
        "ab32: duplicate key": "b524 component selectors: duplicate key",
        "ab32: empty identity string": "b524 component selectors: empty identity string",
        "ab32: expected array": "b524 component selectors: expected array",
        "ab32: incomplete identity": "b524 component selectors: incomplete identity",
        "ab32: invalid or duplicate id": "b524 component selectors: invalid or duplicate id",
        "ab32: missing %s": "b524 component selectors: missing %s",
        "ab32: null primary": "b524 component selectors: null primary",
    },
    "catalog/vaillant/b524_energy_roles.go": {
        "ab95 metadata: %w": "b524 energy roles metadata: %w",
        "ab95 metadata: duplicate key at %d": "b524 energy roles metadata: duplicate key at %d",
        "ab95 metadata: expected five targets": "b524 energy roles metadata: expected five targets",
        "ab95 metadata: invalid entry %d": "b524 energy roles metadata: invalid entry %d",
        "ab95 metadata: trailing data": "b524 energy roles metadata: trailing data",
    },
    "catalog/vaillant/b524_family_codebooks.go": {
        "ab30 duplicate identity": "b524 family codebooks: duplicate identity",
        "ab30 entries must be an array": "b524 family codebooks: entries must be an array",
        "ab30 missing identity axis %s": "b524 family codebooks: missing identity axis %s",
        "ab30 string axes require explicit tokens": "b524 family codebooks: string axes require explicit tokens",
    },
    "catalog/vaillant/b524_name_annotations.go": {
        "ab36: documentary name field is not runtime data": "b524 name annotations: documentary name field is not runtime data",
        "ab36: duplicate key": "b524 name annotations: duplicate key",
        "ab36: expected an array": "b524 name annotations: expected an array",
        "ab36: incomplete key": "b524 name annotations: incomplete key",
        "ab36: missing key": "b524 name annotations: missing key",
    },
    "catalog/vaillant/b524_numeric_candidates.go": {
        "ab56: all caller identity strings are required": "b524 numeric candidates: all caller identity strings are required",
        "ab56: duplicate or invalid JSON member": "b524 numeric candidates: duplicate or invalid JSON member",
        "ab56: duplicate tuple key at row %d": "b524 numeric candidates: duplicate tuple key at row %d",
        "ab56: expected a tuple array": "b524 numeric candidates: expected a tuple array",
        "ab56: invalid or promoted metadata at row %d": "b524 numeric candidates: invalid or promoted metadata at row %d",
        "ab56: row %d has an incomplete or identity-qualified tuple key": "b524 numeric candidates: row %d has an incomplete or identity-qualified tuple key",
        "ab56: row %d must be an object": "b524 numeric candidates: row %d must be an object",
        "ab56: unexpected delimiter": "b524 numeric candidates: unexpected delimiter",
    },
    "catalog/vaillant/b555_additional_forms.go": {
        "ab26: duplicate identity at row %d": "b555 additional forms: duplicate identity at row %d",
        "ab26: expected JSON array": "b555 additional forms: expected JSON array",
        "ab26: incomplete identity at row %d": "b555 additional forms: incomplete identity at row %d",
        "ab26: missing identity axis %s": "b555 additional forms: missing identity axis %s",
    },
    "catalog/vaillant/schema.go": {
        "ab05 key: opcode, group, instance and register must be supplied and non-null": "schema key: opcode, group, instance and register must be supplied and non-null",
        "ab05 record %d field %d: extent exceeds body": "schema record %d field %d: extent exceeds body",
        "ab05 record %d field %d: invalid extent": "schema record %d field %d: invalid extent",
        "ab05 record %d field %d: invalid grade": "schema record %d field %d: invalid grade",
        "ab05 record %d field %d: width exceeds body": "schema record %d field %d: width exceeds body",
        "ab05 record %d: duplicate alternative key": "schema record %d: duplicate alternative key",
        "ab05 record %d: invalid grade": "schema record %d: invalid grade",
        "ab05 record %d: invalid key dimension": "schema record %d: invalid key dimension",
        "ab05 record %d: negative body width": "schema record %d: negative body width",
    },
}

VAILLANT_KEPT = ["ab15-", "ab17-", "ab19_", "ab20-", "ab32-"]

VAILLANT_ROOT = "catalog/vaillant"
GO_STRING = re.compile(r'//[^\n]*|/\*.*?\*/|\'(?:\\.|[^\'\\])*\'|"(?:\\.|[^"\\\n])*"|`[^`]*`', re.S)
JSON_STRING = re.compile(r'"((?:\\.|[^"\\])*)"(\s*:)?')
STRUCT_OPEN = re.compile(r"\bstruct\s*\{")
JSON_TAG = re.compile(r'json:"([^",]*)')
CODED_LITERAL = re.compile(r'"((?:ab|aa|Ab|Aa)[0-9]{2}[^"]*)"')


def vaillant_go_files(root, tests):
    for path in sorted((root / VAILLANT_ROOT).glob("*.go")):
        if path.name.endswith("_test.go") == tests:
            yield path.relative_to(root).as_posix(), path


def rewrite_go_strings(text, rx, table):
    count = 0

    def literal(m):
        nonlocal count
        if m.group().startswith("/"):
            return m.group()
        new, n = rx.subn(lambda t: table[t.group(1)], m.group())
        count += n
        return new
    return GO_STRING.sub(literal, text), count


def struct_collisions(rel, text):
    for opening in STRUCT_OPEN.finditer(text):
        depth, i, names = 1, opening.end(), []
        line_start = i
        while depth:
            c = text[i]
            if c == "`":
                end = text.index("`", i + 1)
                if depth == 1:
                    names += JSON_TAG.findall(text[i:end])
                i = end
            elif c == "{":
                depth += 1
            elif c == "}":
                depth -= 1
            i += 1
        folded = [n.lower() for n in names if n and n != "-"]
        for name in set(folded):
            if folded.count(name) > 1:
                raise SystemExit(f"{rel}: json key {name} would collide inside one type")


def rewrite_json(rel, text, keys, values):
    count = 0
    folded = {old.lower(): new for old, new in keys.items()}

    def string(m):
        nonlocal count
        word, colon = m.group(1), m.group(2)
        new = folded.get(word.lower()) if colon else values.get(word)
        if new is not None:
            count += 1
            return '"' + new + '"' + (colon or "")
        return m.group()
    new = JSON_STRING.sub(string, text)

    def unique(pairs):
        folded = [k.lower() for k, _ in pairs]
        for key in set(folded):
            if folded.count(key) > 1:
                raise SystemExit(f"{rel}: key {key} would collide inside one object")
        return dict(pairs)
    json.loads(new, object_pairs_hook=unique)
    return new, count


def vaillant_scope(root):
    production = {rel for rel, _ in vaillant_go_files(root, False)}
    listed = set(VAILLANT_TEXTS)
    for rel, path in vaillant_go_files(root, False):
        hits = set(CODED_LITERAL.findall(path.read_text())) - set(VAILLANT_KEPT)
        texts = VAILLANT_TEXTS.get(rel, {})
        unlisted = {h for h in hits if h not in texts and h.split(",")[0] not in VAILLANT_TOKENS}
        if unlisted:
            raise SystemExit(f"{rel}: unlisted coded string {sorted(unlisted)}")
        if hits:
            listed.add(rel)
    if listed - production:
        raise SystemExit(f"listed files missing: {sorted(listed - production)}")
    return listed


def vaillant_plan(root):
    listed = vaillant_scope(root)
    rx = pattern(VAILLANT_TOKENS)
    values = {old: VAILLANT_TOKENS[old] for old in VAILLANT_VALUES}
    changed = {}
    for tests in (False, True):
        for rel, path in vaillant_go_files(root, tests):
            text = path.read_text()
            new, count = rewrite_go_strings(text, rx, VAILLANT_TOKENS)
            for old, replacement in VAILLANT_TEXTS.get(rel, {}).items():
                quoted = '"' + old + '"'
                found = new.count(quoted)
                if found == 0:
                    raise SystemExit(f"{rel}: expected {quoted}")
                new, count = new.replace(quoted, '"' + replacement + '"'), count + found
            if rel in listed and count == 0:
                raise SystemExit(f"{rel}: no coded string found")
            if count:
                struct_collisions(rel, new)
                changed[rel] = new
    for directory, table in (("data", {}), ("testdata", values)):
        for path in sorted((root / VAILLANT_ROOT / directory).glob("*.json")):
            rel = path.relative_to(root).as_posix()
            new, count = rewrite_json(rel, path.read_text(), VAILLANT_TOKENS, table)
            if count:
                changed[rel] = new
    changed[INVENTORY] = vaillant_inventory(root)
    return changed


def vaillant_inventory(root):
    data = json.loads((root / INVENTORY).read_text())
    count = 0
    for row in data["rows"]:
        if row["owning_file"].startswith(VAILLANT_ROOT + "/") and row["current_name"] in VAILLANT_TOKENS:
            row["current_name"] = VAILLANT_TOKENS[row["current_name"]]
            count += 1
    if count == 0:
        raise SystemExit(f"{INVENTORY}: no coded row found")
    return json.dumps(data, indent=2, ensure_ascii=False) + "\n"


def vaillant_check(root):
    remaining = []
    for rel, path in vaillant_go_files(root, False):
        hits = set(CODED_LITERAL.findall(path.read_text())) - set(VAILLANT_KEPT)
        remaining += [f"{rel}: {hit}" for hit in sorted(hits)]
    for path in sorted((root / VAILLANT_ROOT / "data").glob("*.json")):
        for m in JSON_STRING.finditer(path.read_text()):
            if m.group(2) and re.match(r"(?:ab|aa|Ab|Aa)[0-9]{2}", m.group(1)):
                remaining.append(f"{path.relative_to(root).as_posix()}: key {m.group(1)}")
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
    "catalog-vaillant": (vaillant_plan, lambda: changelog_rows([VAILLANT_TOKENS] + list(VAILLANT_TEXTS.values()))),
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
        checks = {"system-b555-tags": tag_check, "catalog-vaillant": vaillant_check}
        if args.slice not in checks:
            raise SystemExit("--check covers the system-b555-tags and catalog-vaillant slices")
        remaining = checks[args.slice](root)
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
