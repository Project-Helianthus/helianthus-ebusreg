# Changelog

## Unreleased

- Remove work-item codes from the `catalog/vaillant` files whose code
  number is 01 to 29 and from the Go identifiers that start with their
  codes. Each production file and its test twin is renamed; embedded
  data and testdata fixtures keep their names. The two evaluators
  `b512_heatpump_io_evaluator.go` and `b51a_mirror_fields_evaluator.go`
  take an `_evaluator` suffix because a catalog file holds the plain
  name:
  - `ab01ao01_b505_sentinel_views.go` -> `b505_sentinel_views.go`
  - `ab01ap01_b531_capability_state.go` -> `b531_capability_state.go`
  - `ab01aq01_b504_candidate_fields.go` -> `b504_candidate_fields.go`
  - `ab03ao03_b506_parameter_views.go` -> `b506_parameter_views.go`
  - `ab03ap03_b531_ordinary_limit.go` -> `b531_ordinary_limit.go`
  - `ab03aq03_b505_override_fields.go` -> `b505_override_fields.go`
  - `ab05_schema.go` -> `schema.go`
  - `ab05an05_b503_history_branches.go` -> `b503_history_branches.go`
  - `ab05ao05_b507_quantity_views.go` -> `b507_quantity_views.go`
  - `ab05ap05_b516_cycle_reducer.go` -> `b516_cycle_reducer.go`
  - `ab05aq05_b505_conditional_roles.go` -> `b505_conditional_roles.go`
  - `ab07ao07_b50e_row_spans.go` -> `b50e_row_spans.go`
  - `ab07ap07_b555_weekly_refresh.go` -> `b555_weekly_refresh.go`
  - `ab07aq07_b503_control_fields.go` -> `b503_control_fields.go`
  - `ab09ao09_b50f_field_views.go` -> `b50f_field_views.go`
  - `ab09ap09_b555_manager_gate.go` -> `b555_manager_gate.go`
  - `ab09aq09_b515_period_fields.go` -> `b515_period_fields.go`
  - `ab11an11_b511_overlap_views.go` -> `b511_overlap_views.go`
  - `ab11ao11_b510_cache_transition.go` -> `b510_cache_transition.go`
  - `ab11ap11_b512_fan_fields.go` -> `b512_fan_fields.go`
  - `ab11aq11_b517_local_guard.go` -> `b517_local_guard.go`
  - `ab12_b50f_status_spans.go` -> `b50f_status_spans.go`
  - `ab13_b507_energy_candidates.go` -> `b507_energy_candidates.go`
  - `ab13an13_b512_dhw_views.go` -> `b512_dhw_views.go`
  - `ab13ao13_b511_failure_transition.go` -> `b511_failure_transition.go`
  - `ab13ap13_b512_heatpump_io.go` -> `b512_heatpump_io_evaluator.go`
  - `ab13aq13_b521_field_views.go` -> `b521_field_views.go`
  - `ab14_b509_op59_quantities.go` -> `b509_op59_quantities.go`
  - `ab15_b504_circuit_roles.go` -> `b504_circuit_roles.go`
  - `ab15an15_b516_clock_word_views.go` -> `b516_clock_word_views.go`
  - `ab15ao15_b512_feedback_views.go` -> `b512_feedback_views.go`
  - `ab15ap15_b512_emergency_flow.go` -> `b512_emergency_flow.go`
  - `ab16_b505_operating_roles.go` -> `b505_operating_roles.go`
  - `ab17_b506_rf_sequence.go` -> `b506_rf_sequence.go`
  - `ab17an17_b516_two_word_views.go` -> `b516_two_word_views.go`
  - `ab17ao17_b513_information_views.go` -> `b513_information_views.go`
  - `ab17ap17_b509_diagnostic_fields.go` -> `b509_diagnostic_fields.go`
  - `ab18_b512_feedback_fallback.go` -> `b512_feedback_fallback.go`
  - `ab19_b512_heatpump_io.go` -> `b512_heatpump_io.go`
  - `ab19an19_b531_limit_numeric_views.go` -> `b531_limit_numeric_views.go`
  - `ab19ao19_b514_menu_views.go` -> `b514_menu_views.go`
  - `ab19ap19_b51a_mirror_fields.go` -> `b51a_mirror_fields_evaluator.go`
  - `ab19aq19_b51a_handshake_fields.go` -> `b51a_handshake_fields.go`
  - `ab20_b516_directional_types.go` -> `b516_directional_types.go`
  - `ab21_b519_reply_policy.go` -> `b519_reply_policy.go`
  - `ab21ao21_b519_opaque_views.go` -> `b519_opaque_views.go`
  - `ab21ap21_b523_actuator_fields.go` -> `b523_actuator_fields.go`
  - `ab21aq21_b51a_diagnostic_fields.go` -> `b51a_diagnostic_fields.go`
  - `ab23an23_b51b_local_budget_crc.go` -> `b51b_local_budget_crc.go`
  - `ab23ao23_b520_actuator_views.go` -> `b520_actuator_views.go`
  - `ab23ap23_b516_energy_type_fields.go` -> `b516_energy_type_fields.go`
  - `ab23aq23_b512_valve_pump_fields.go` -> `b512_valve_pump_fields.go`
  - `ab25ao25_b522_demand_views.go` -> `b522_demand_views.go`
  - `ab25ap25_b525_quality_config.go` -> `b525_quality_config.go`
  - `ab25aq25_b512_power_stepper_fields.go` -> `b512_power_stepper_fields.go`
  - `ab26_b555_additional_forms.go` -> `b555_additional_forms.go`
  - `ab27_holiday.go` -> `holiday.go`
  - `ab27ao27_b523_sensor_views.go` -> `b523_sensor_views.go`
  - `ab29_room_actions.go` -> `room_actions.go`
  Identifiers that start with a code (709 exported, 3863 unexported)
  take the stem of their destination file; `Ab<code>X` becomes
  `<Stem>X` and `ab<code>X` becomes `<stem>X`:
  - `Ab01Ao01...` -> `B505SentinelView...`
  - `Ab01Ap01...` -> `B531CapabilityState...`
  - `Ab01Aq01...` -> `B504CandidateField...`
  - `Ab03Ao03...` -> `B506ParameterView...`
  - `Ab03Ap03...` -> `B531OrdinaryLimit...`
  - `Ab03Aq03...` -> `B505OverrideField...`
  - `Ab05...` -> `Schema...`
  - `Ab05An05...` -> `B503HistoryBranch...`
  - `Ab05Ao05...` -> `B507QuantityView...`
  - `Ab05Ap05...` -> `B516CycleReducer...`
  - `Ab05Aq05...` -> `B505ConditionalRole...`
  - `Ab07Ao07...` -> `B50ERowSpan...`
  - `Ab07Ap07...` -> `B555WeeklyRefresh...`
  - `Ab07Aq07...` -> `B503ControlField...`
  - `Ab09Ao09...` -> `B50FFieldView...`
  - `Ab09Ap09...` -> `B555ManagerGate...`
  - `Ab09Aq09...` -> `B515PeriodField...`
  - `Ab11An11...` -> `B511OverlapView...`
  - `Ab11Ao11...` -> `B510CacheTransition...`
  - `Ab11Ap11...` -> `B512FanField...`
  - `Ab11Aq11...` -> `B517LocalGuard...`
  - `Ab12...` -> `B50FStatusSpan...`
  - `Ab13...` -> `B507EnergyCandidate...`
  - `Ab13An13...` -> `B512DHWView...`
  - `Ab13Ao13...` -> `B511FailureTransition...`
  - `Ab13Ap13...` -> `B512HeatpumpIOEvaluator...`
  - `Ab13Aq13...` -> `B521FieldView...`
  - `Ab14...` -> `B509Op59Quantity...`
  - `Ab15...` -> `B504CircuitRole...`
  - `Ab15An15...` -> `B516ClockWordView...`
  - `Ab15Ao15...` -> `B512FeedbackView...`
  - `Ab15Ap15...` -> `B512EmergencyFlow...`
  - `Ab16...` -> `B505OperatingRole...`
  - `Ab17...` -> `B506RFSequence...`
  - `Ab17An17...` -> `B516TwoWordView...`
  - `Ab17Ao17...` -> `B513InformationView...`
  - `Ab17Ap17...` -> `B509DiagnosticField...`
  - `Ab18...` -> `B512FeedbackFallback...`
  - `Ab19...` -> `B512HeatpumpIO...`
  - `Ab19An19...` -> `B531LimitNumericView...`
  - `Ab19Ao19...` -> `B514MenuView...`
  - `Ab19Ap19...` -> `B51AMirrorFieldEvaluator...`
  - `Ab19Aq19...` -> `B51AHandshakeField...`
  - `Ab20...` -> `B516DirectionalType...`
  - `Ab21...` -> `B519ReplyPolicy...`
  - `Ab21Ao21...` -> `B519OpaqueView...`
  - `Ab21Ap21...` -> `B523ActuatorField...`
  - `Ab21Aq21...` -> `B51ADiagnosticField...`
  - `Ab23An23...` -> `B51BLocalBudgetCRC...`
  - `Ab23Ao23...` -> `B520ActuatorView...`
  - `Ab23Ap23...` -> `B516EnergyTypeField...`
  - `Ab23Aq23...` -> `B512ValvePumpField...`
  - `Ab25Ao25...` -> `B522DemandView...`
  - `Ab25Ap25...` -> `B525QualityConfig...`
  - `Ab25Aq25...` -> `B512PowerStepperField...`
  - `Ab26...` -> `B555AdditionalForm...`
  - `Ab27...` -> `Holiday...`
  - `Ab27Ao27...` -> `B523SensorView...`
  - `Ab29...` -> `RoomAction...`
  Exported names whose remainder repeats the stem drop the repeated
  words:
  - `Ab11An11View` -> `B511OverlapView`
  - `Ab13Ap13Evaluate` -> `B512HeatpumpIOEvaluate`
  - `Ab15An15WordView` -> `B516ClockWordView`
  - `Ab15Role` -> `B504CircuitRole`
  - `Ab16Role` -> `B505OperatingRole`
  - `Ab19Ap19Evaluate` -> `B51AMirrorFieldEvaluate`
  - `Ab21Policy` -> `B519ReplyPolicy`
  - `Ab23An23CRC` -> `B51BLocalBudgetCRC`
  - `Ab23An23CRCWithInitial` -> `B51BLocalBudgetCRCWithInitial`
  - `Ab27HolidayCatalog` -> `HolidayCatalog`
  - `Ab29RoomActions` -> `RoomActions`
  The `bytes` import aliases that would equal a byte slice name take a
  `BytesPackage` suffix, and the `encoding/json` imports of
  `b511_failure_transition.go` and its test drop their coded alias so
  `B511FailureTransitionInterpret` stays within the size ratchet.
  Files with codes 30 to 99 that use the renamed `Schema...` names are
  edited in place and keep their names until the next slice. Renamed
  exported struct fields that had no JSON tag carry an explicit tag
  with the unchanged key, so serialized output and embedded data
  decoding stay the same; `TestRenamedFieldJSONKeys` pins every key.
  Consumers must update references; there are no aliases. JSON tags,
  string values, error texts, `go:embed` paths, data and fixture names
  and test function names are unchanged.
  `scripts/rename_catalog_names.py --slice vaillant-01-29` reproduces
  the change. Gateway adoption and release bill of materials
  reconciliation remain required.

- Add package documentation to the 13 packages that had none: the module
  root, `catalog/vaillant`, `internal/match`, `providers/vaillant`,
  `registry`, `router`, `schema`, `vaillant/b555`, `vaillant/dhw`,
  `vaillant/heating`, `vaillant/productids`, `vaillant/solar` and
  `vaillant/system`. Twelve packages get a new `doc.go` that holds only the
  package comment (responsibility, ownership boundary and non-obvious
  contracts) and the package clause; `vaillant/b555` carries its package
  comment on the package clause of the existing `vaillant/b555/b555.go`
  (comment lines only); no other Go file changes. The 13
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
