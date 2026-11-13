# Changelog

## Unreleased

- Add package documentation to the 13 packages that had none: the module
  root, `catalog/vaillant`, `internal/match`, `providers/vaillant`,
  `registry`, `router`, `schema`, `vaillant/b555`, `vaillant/dhw`,
  `vaillant/heating`, `vaillant/productids`, `vaillant/solar` and
  `vaillant/system`. Each package gets a new `doc.go` that holds only the
  package comment (responsibility, ownership boundary and non-obvious
  contracts) and the package clause; no other Go file changes. The 13
  `revive/package-comments` rows leave `scripts/go-lint-baseline.json`.
- Remove work-item codes from `catalog/ebusstandard` file names and Go
  identifiers. Each production file and its test twin is renamed; the
  testdata fixtures keep their names:
  - `ab15aq15_standard07_field_views.go` -> `standard07_field_views.go`
  - `ab17aq17_standardff_extent_view.go` -> `standardff_extent_view.go`
  - `ab25an25_standard09_body_views.go` -> `standard09_body_views.go`
  - `ab28_standard07_residual_fields.go` -> `standard07_residual_fields.go`
  - `ab49_fe_metadata.go` -> `standardfe_metadata.go`
  - `ab54_standard_ff_body_states.go` -> `standardff_body_states.go`
  - `ab76_standard03_candidate_views.go` -> `standard03_candidate_views.go`
  - `ab77_07_candidates.go` -> `standard07_candidates.go`
  - `ab78_standard05_candidate_views.go` -> `standard05_candidate_views.go`
  - `ab79_09_candidates.go` -> `standard09_candidates.go`
  - `ab81_0f_candidates.go` -> `standard0f_candidates.go`
  - `ab83_03_counter_bodies.go` -> `standard03_counter_bodies.go`
  - `ab84_05_burner_bodies.go` -> `standard05_burner_bodies.go`
  Identifiers that start with a code (349 exported, 622 unexported)
  take the stem of their destination file; `Ab<code>X` becomes
  `<Stem>X` and `ab<code>X` becomes `<stem>X`:
  - `Ab15Aq15...` -> `Standard07FieldView...`
  - `Ab17Aq17...` -> `StandardFFExtentView...`
  - `Ab25An25...` -> `Standard09BodyView...`
  - `Ab28...` -> `Standard07Residual...`
  - `Ab49...` -> `StandardFEMetadata...`
  - `Ab54...` -> `StandardFFBodyState...`
  - `Ab76...` -> `Standard03CandidateView...`
  - `Ab77...` -> `Standard07Candidate...`
  - `Ab78...` -> `Standard05CandidateView...`
  - `Ab79...` -> `Standard09Candidate...`
  - `Ab81...` -> `Standard0fCandidate...`
  - `Ab83...` -> `Standard03CounterBody...`
  - `Ab84...` -> `Standard05BurnerBody...`
  Exported names whose remainder repeats the stem drop the repeated
  words:
  - `Ab25An25ViewAddress` -> `Standard09ViewAddress`
  - `Ab25An25ViewBody` -> `Standard09ViewBody`
  - `Ab25An25ViewCount` -> `Standard09ViewCount`
  - `Ab49FEMetadata` -> `StandardFEMetadata`
  - `Ab76View` -> `Standard03CandidateView`
  - `Ab77CandidateMetadata` -> `Standard07CandidateMetadata`
  - `Ab77Candidates` -> `Standard07Candidates`
  - `Ab78View` -> `Standard05CandidateView`
  - `Ab79Candidates` -> `Standard09Candidates`
  - `Ab81Candidate` -> `Standard0fCandidate`
  - `Ab81Candidates` -> `Standard0fCandidates`
  - `Ab83Body` -> `Standard03CounterBody`
  - `Ab83BodyState` -> `Standard03CounterBodyState`
  - `Ab83CounterBodies` -> `Standard03CounterBodies`
  - `Ab83LookupBody` -> `Standard03CounterLookupBody`
  - `Ab84BurnerBodies` -> `Standard05BurnerBodies`
  - `Ab84BurnerField` -> `Standard05BurnerField`
  - `Ab84BurnerMetadata` -> `Standard05BurnerMetadata`
  - `Ab84LookupBurnerField` -> `Standard05LookupBurnerField`
  Consumers must update references; there are no aliases. JSON tags,
  string values, error texts, fixture names and test function names are
  unchanged. `scripts/rename_catalog_names.py` reproduces the change.
  Gateway adoption and release bill of materials reconciliation remain
  required.

- Move `catalog/ebus_standard` (package `ebus_standard_catalog`) to
  `catalog/ebusstandard` (package `ebusstandard`), together with its
  `internal/safetypolicy` package. Consumers must update imports; the old
  path has no forwarding package. The `ebus_standard` service namespace,
  method identifiers, catalog data and error texts are unchanged. Gateway
  adoption and release bill of materials reconciliation remain required.

- Remove the unused duplicate `pv/` package and its tests. Protocol-independent
  photovoltaic semantics belong to `helianthus-semreg` (`packs/pv`); native
  eBUS registry ownership remains here. Gateway adoption and release bill of
  materials reconciliation remain required.

- Standard 0700 and 0702 candidate metadata and field views emit
  `outdoor_temperature` for the documented outdoor quantity. Previous keys do
  not resolve as candidates. `duration`
  stays provisional with its existing zero semantics.
- B524 annotation APIs emit the existing public OP02 canonical names and their
  open questions without documentary legacy-name or display-alias fields.
  All 142 provisional OP02 names remain unchanged. The two OP06 native labels
  remain distinct evidence and do not become canonical names or runtime aliases.
  Documentary correspondence is retained in the non-runtime naming inventory.
- Vaillant public documentation references adopt lowercase page filenames.
  Consumer adoption and renewed release evidence are required before release.
