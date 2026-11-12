#!/usr/bin/env python3
"""Reproduce the removal of work-item codes from catalog/ebusstandard names.

Production and test Go files lose their code prefix and keep the behavior part
of their name. Go identifiers that start with a code prefix are renamed to a
stem named after the standard service of the destination file. The rewrite is
token-aware: string literals (JSON tags, values, error texts, fixture paths)
and testdata file names stay byte-identical, so no externally visible string
changes here. Comments only change where they name a renamed identifier.
Test functions (TestAb...) keep their names. Two test case structs bind
fixture keys through untagged field names; their fields get explicit tags
carrying the unchanged fixture keys.

Run once against a tree that still holds the coded files. The script stops
before writing anything when a precondition or a collision check fails, and
prints the identifier map with --print-map.
"""

import argparse
import importlib.util
import json
import os
from pathlib import Path
import re
import subprocess
import sys

PACKAGE_DIR = "catalog/ebusstandard"
IMPORT_PATH = "github.com/Project-Helianthus/helianthus-ebusreg/" + PACKAGE_DIR
GOLDEN = PACKAGE_DIR + "/testdata/provider_abi.golden.txt"

# Coded file stem -> behavior file stem. Each applies to the .go file and to
# its _test.go twin. Testdata fixtures keep their names.
FILE_MOVES = (
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
)

# Code prefix (lowercase) -> UpperCamelCase stem derived from the destination
# file. Ab<code>X becomes <Stem>X and ab<code>X becomes <stem>X.
STEMS = {
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
}

# Exported names whose remainder repeats the stem; the stem words that would
# repeat are dropped.
EXPORTED_OVERRIDES = {
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
}

# Resolved collisions: the encoding/json import alias of the standard 09 test
# would otherwise equal the name of its fixture byte slice (ab25an25JSON).
COLLISION_OVERRIDES = {
    "ab25an25json": "standard09BodyViewEncodingJSON",
}

# A lowercase remainder (import aliases, short locals) is capitalized; these
# words take their Go initialism form instead.
INITIALISMS = {"json": "JSON", "id": "ID"}

# Test case structs decoded from fixtures through untagged field names; the
# old field name is the fixture key and becomes an explicit tag.
FIXTURE_BOUND = {
    "ab76_standard03_candidate_views_test.go": ("ab76Cases", {
        "Ab76Label": "Ab76Label", "Ab76Operation": "Ab76Operation",
        "Ab76Body": "Ab76Body", "Ab76Names": "Ab76Names",
        "Ab76Values": "Ab76Values", "Ab76Counter": "Ab76Counter",
        "Ab76Hours": "Ab76Hours", "Ab76Minute": "Ab76Minute"}),
    "ab78_standard05_candidate_views_test.go": ("ab78Cases", {
        "Ab78Label": "ab78Label", "Ab78Body": "ab78Body",
        "Ab78Values": "ab78Values", "Ab78Name": "ab78Name",
        "Ab78Domains": "ab78Domains"}),
}

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


def coded(name):
    match = CODE.fullmatch(name)
    if match and match.group(1).lower() in STEMS:
        return match.group(1).lower(), match.group(2)
    return None


def new_name(name):
    if name in EXPORTED_OVERRIDES:
        return EXPORTED_OVERRIDES[name]
    if name in COLLISION_OVERRIDES:
        return COLLISION_OVERRIDES[name]
    code, rest = coded(name)
    if rest[:1].islower():
        rest = INITIALISMS.get(rest, rest[:1].upper() + rest[1:])
    stem = STEMS[code]
    if name[0].islower():
        stem = stem[:1].lower() + stem[1:]
    return stem + rest


def package_files(root):
    return sorted(p for p in (root / PACKAGE_DIR).glob("*.go") if p.is_file())


def build_map(root):
    old, other = set(), set()
    for path in package_files(root):
        for kind, text in lex(path.read_text()):
            if kind == "ident":
                (old if coded(text) else other).add(text)
    for path in sorted((root / PACKAGE_DIR / "internal").rglob("*.go")):
        if any(kind == "ident" and coded(text) for kind, text in lex(path.read_text())):
            raise SystemExit(f"coded identifier outside the package scope: {path}")
    for table in (EXPORTED_OVERRIDES, COLLISION_OVERRIDES):
        unused = sorted(set(table) - old)
        if unused:
            raise SystemExit(f"override for an absent identifier: {unused}")
    mapping = {name: new_name(name) for name in sorted(old)}
    seen = {}
    for name, target in mapping.items():
        if target in seen:
            raise SystemExit(f"collision: {seen[target]} and {name} -> {target}")
        seen[target] = name
        if target in other or target in KEYWORDS or coded(target):
            raise SystemExit(f"collision with an existing name: {name} -> {target}")
        if name[0].isupper() != target[0].isupper():
            raise SystemExit(f"export status changes: {name} -> {target}")
    return mapping


def check_consumers(root, mapping):
    for path in sorted(root.rglob("*.go")):
        rel = path.relative_to(root)
        if rel.parts[0].startswith(".") or rel.parent.as_posix() == PACKAGE_DIR:
            continue
        tokens = [t for t in lex(path.read_text()) if t[0] != "comment" and t[1].strip()]
        aliases = set()
        for index, (kind, text) in enumerate(tokens):
            if kind == "string" and text.strip('"`') == IMPORT_PATH:
                kind, alias = tokens[index - 1]
                if kind != "ident" or alias == "import":
                    alias = "ebusstandard"
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


def preserve_projection_json_keys(source):
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
            if match and coded(match.group(1)):
                lines[index] = line.rstrip("\n") + f' `json:"{match.group(1)}"`\n'
    return "".join(lines)


def rewrite_go(source, mapping):
    out = []
    for kind, text in lex(source):
        if kind == "ident":
            text = mapping.get(text, text)
        elif kind == "comment":
            text = rename_comment(text, mapping)
        out.append(text)
    return "".join(out)


def rename_symbol(symbol, mapping):
    return ".".join(mapping.get(part, part) for part in symbol.split("."))


def moved_paths():
    moves = {}
    for old, new in FILE_MOVES:
        for suffix in (".go", "_test.go"):
            moves[f"{PACKAGE_DIR}/{old}{suffix}"] = f"{PACKAGE_DIR}/{new}{suffix}"
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


def update_baselines(root, mapping, moves):
    module, vocabulary = behavior_tokens(root)
    path = root / "scripts/behavior-name-baseline.json"
    data = load_json(path)
    rows = []
    for row in data["baseline"]:
        if row["path"] in moves:
            row = dict(row, path=moves[row["path"]])
            if row["category"] == "filename":
                row["symbol"] = row["path"].rsplit("/", 1)[1]
        if row["path"].startswith(PACKAGE_DIR + "/") and row["category"] == "identifier":
            row = dict(row, symbol=mapping.get(row["symbol"], row["symbol"]))
        if module.violations([row], vocabulary):
            rows.append(row)
    data["baseline"] = sorted(rows, key=lambda r: (r["path"], r["symbol"], r["category"]))
    dump_json(path, data, True)

    path = root / "scripts/go-lint-baseline.json"
    data = load_json(path)
    for row in data["entries"]:
        if row["path"].startswith(PACKAGE_DIR + "/"):
            row["path"] = moves.get(row["path"], row["path"])
            row["symbol"] = rename_symbol(row["symbol"], mapping)
    data["entries"].sort(key=lambda r: (r["path"], r["symbol"], r["rule"]))
    dump_json(path, data, True)

    path = root / "scripts/size-baseline.json"
    data = load_json(path)
    for row in data["debt"]:
        if row["path"].startswith(PACKAGE_DIR + "/"):
            row["path"] = moves.get(row["path"], row["path"])
            row["symbol"] = rename_symbol(row["symbol"], mapping)
    data["debt"].sort(key=lambda r: (r["path"], r["symbol"], r["metric"]))
    dump_json(path, data, True)


def update_inventories(root, mapping, moves):
    path = root / "naming-inventory.json"
    data = load_json(path)
    data = replace_paths(data, moves)
    data["source_files"] = sorted(data["source_files"])
    for row in data["rows"]:
        if (row["owning_file"].startswith(PACKAGE_DIR + "/")
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


def update_golden(root, mapping):
    # Use the snapshot's own serializer so tags added to projections are
    # recorded along with renamed symbols and their signatures.
    subprocess.run(
        ["go", "test", "./" + PACKAGE_DIR,
         "-run", "^TestProvider_ABISnapshot$", "-count=1", "-update"],
        cwd=root, check=True)


def changelog_entry(mapping):
    files = "".join(f"  - `{old}.go` -> `{new}.go`\n" for old, new in FILE_MOVES)
    stems = "".join(f"  - `{code[:1].upper() + code[1:4] + code[4:5].upper() + code[5:]}...`"
                    f" -> `{stem}...`\n" for code, stem in STEMS.items())
    overrides = "".join(f"  - `{old}` -> `{new}`\n"
                        for old, new in sorted(EXPORTED_OVERRIDES.items()))
    exported = sum(1 for name in mapping if name[0].isupper())
    return (
        "- Remove work-item codes from `catalog/ebusstandard` file names and Go\n"
        "  identifiers. Each production file and its test twin is renamed; the\n"
        "  testdata fixtures keep their names:\n" + files +
        f"  Identifiers that start with a code ({exported} exported, "
        f"{len(mapping) - exported} unexported)\n"
        "  take the stem of their destination file; `Ab<code>X` becomes\n"
        "  `<Stem>X` and `ab<code>X` becomes `<stem>X`:\n" + stems +
        "  Exported names whose remainder repeats the stem drop the repeated\n"
        "  words:\n" + overrides +
        "  Consumers must update references; there are no aliases. JSON tags,\n"
        "  string values, error texts, fixture names and test function names are\n"
        "  unchanged. `scripts/rename_catalog_names.py` reproduces the change.\n"
        "  Gateway adoption and release bill of materials reconciliation remain\n"
        "  required.\n")


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


def add_projection_json_keys_test(source):
    if "func TestRenamedProjectionJSONKeys(" in source:
        raise SystemExit("projection JSON-key regression test already exists")
    return source.rstrip() + "\n\n" + PROJECTION_JSON_KEYS_TEST


def run(root, print_map):
    moves = moved_paths()
    for old, new in moves.items():
        if not (root / old).is_file():
            raise SystemExit(f"source file missing: {old}")
        if (root / new).exists():
            raise SystemExit(f"target already exists: {new}")
    mapping = build_map(root)
    check_consumers(root, mapping)
    if print_map:
        for name, target in mapping.items():
            print(f"{name} {target}")
        return
    sources = {}
    for path in package_files(root):
        source = path.read_text()
        if not path.name.endswith("_test.go"):
            source = preserve_projection_json_keys(source)
        if path.name in FIXTURE_BOUND:
            source = add_fixture_tags(source, *FIXTURE_BOUND[path.name])
        source = rewrite_go(source, mapping)
        if path.name == "canonical_names_test.go":
            source = add_projection_json_keys_test(source)
        sources[path] = source
    for path, source in sources.items():
        path.write_text(source)
    for old, new in moves.items():
        os.rename(root / old, root / new)
    changed = [str(root / moves.get(p.relative_to(root).as_posix(), p.relative_to(root).as_posix()))
               for p in sources]
    subprocess.run(["gofmt", "-w", *changed], check=True)
    update_golden(root, mapping)
    update_baselines(root, mapping, moves)
    update_inventories(root, mapping, moves)
    changelog = root / "CHANGELOG.md"
    text = changelog.read_text()
    heading = "## Unreleased\n\n"
    if heading not in text:
        raise SystemExit("CHANGELOG.md has no Unreleased section")
    changelog.write_text(text.replace(heading, heading + changelog_entry(mapping) + "\n", 1))


def test_inventory_targets(tmp_path):
    row = {'owning_file': 'catalog/ebusstandard/ab76_standard03_candidate_views.go', 'kind': 'go_identifier', 'current_name': 'Ab76Interpret', 'target': 'Ab76Interpret', 'status': 'unchanged', 'identity': {'declaration': 'Ab76Interpret'}}
    (tmp_path / 'naming-inventory.json').write_text(json.dumps({'source_files': [row['owning_file']], 'rows': [row]}))
    (tmp_path / 'naming-consumers.json').write_text('{}')
    update_inventories(tmp_path, {'Ab76Interpret': 'Standard03CandidateViewInterpret'}, moved_paths())
    result = json.loads((tmp_path / 'naming-inventory.json').read_text())['rows'][0]
    assert result['current_name'] == result['target'] == 'Standard03CandidateViewInterpret'

def test_generated_projection_json_keys_test():
    canonical = (Path(__file__).resolve().parent.parent / PACKAGE_DIR /
                 "canonical_names_test.go").read_text()
    parent, regression = canonical.split("func TestRenamedProjectionJSONKeys(", 1)
    assert add_projection_json_keys_test(parent) == canonical
    assert regression.count("map[string]string{") == 11


def test_run_generates_projection_json_keys_test(tmp_path, monkeypatch):
    canonical = (Path(__file__).resolve().parent.parent / PACKAGE_DIR /
                 "canonical_names_test.go").read_text()
    parent = canonical.split("func TestRenamedProjectionJSONKeys(", 1)[0]
    package = tmp_path / PACKAGE_DIR
    package.mkdir(parents=True)
    for old in moved_paths():
        (tmp_path / old).write_text("package ebusstandard\n")
    target = package / "canonical_names_test.go"
    target.write_text(parent)
    (tmp_path / "CHANGELOG.md").write_text("## Unreleased\n\n")
    module = sys.modules[__name__]
    monkeypatch.setattr(module, "build_map", lambda root: {})
    monkeypatch.setattr(module, "add_fixture_tags", lambda source, *args: source)
    monkeypatch.setattr(module, "check_consumers", lambda *args: None)
    for name in ("update_golden", "update_baselines", "update_inventories"):
        monkeypatch.setattr(module, name, lambda *args: None)
    run(tmp_path, False)
    assert target.read_text() == canonical


def test_projection_tags():
    source = 'type Ab76Selection struct {\n\tAb76Family string\n\tAb76Value int `json:"value"`\n}\n'
    tagged = preserve_projection_json_keys(source)
    assert '`json:"Ab76Family"`' in tagged
    assert '`json:"value"`' in tagged
    assert preserve_projection_json_keys(tagged) == tagged


def test_golden_projection_tags(tmp_path):
    import shutil

    package = tmp_path / PACKAGE_DIR
    package.mkdir(parents=True)
    (package / "testdata").mkdir()
    (tmp_path / "go.mod").write_text("module snapshotexample\n\ngo 1.23\n")
    snapshot_test = Path(__file__).resolve().parent.parent / PACKAGE_DIR / "provider_abi_snapshot_test.go"
    shutil.copyfile(snapshot_test, package / snapshot_test.name)
    (package / "internal/safetypolicy").mkdir(parents=True)
    source = ('package ebusstandard\n\ntype Ab76Selection struct {\n'
              '\tAb76Family string\n'
              '\tAb76Value int `json:"value"`\n}\n')
    mapping = {"Ab76Selection": "Standard03CandidateViewSelection",
               "Ab76Family": "Standard03CandidateViewFamily",
               "Ab76Value": "Standard03CandidateViewValue"}
    (package / "projection.go").write_text(
        rewrite_go(preserve_projection_json_keys(source), mapping))
    golden = tmp_path / GOLDEN
    golden.write_text('ebusstandard type Ab76Selection struct { Ab76Family string Ab76Value int `json:"value"` }\n')
    update_golden(tmp_path, mapping)
    assert golden.read_text() == (
        'ebusstandard type Standard03CandidateViewSelection struct { '
        'Standard03CandidateViewFamily string `json:"Ab76Family"` '
        'Standard03CandidateViewValue int `json:"value"` }\n')


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("root", nargs="?", type=Path,
                        default=Path(__file__).resolve().parent.parent)
    parser.add_argument("--print-map", action="store_true",
                        help="print the identifier map and change nothing")
    arguments = parser.parse_args()
    run(arguments.root.resolve(), arguments.print_map)
    sys.exit(0)
