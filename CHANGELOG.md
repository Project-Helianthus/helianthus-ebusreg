# Changelog

## Unreleased

- `vaillant/nativepoints` models every named controller register of the
  published B524 register map. The name annotations gain the local
  parameter group (opcode 0x02, group 0x0A) and the opcode 0x06 device
  groups, with the universal names `device_connected`,
  `device_class_address`, `device_error_code` and
  `device_firmware_version` for registers 0x0001 to 0x0004 of every
  group; unnamed rows get no point. Each annotation carries the map's
  `wire` and `decode` columns, and the new
  `catalog/vaillant/data/controller_register_matrix.json` holds the register
  matrix by family (reply width, value type, unit and enumeration for the
  `controller-700` and `controller-720` families). The generator takes a
  controller location's codec, unit and enumeration from the matrix row of
  its family and otherwise from the map; a family the matrix does not list
  takes the map's layout unless the map reads the register on one family
  only, and a location on such a register names that one family even when
  no matrix row lists it. A matrix row that leaves the unit blank keeps the
  map's unit, and a wire type alone does not make a register numeric when
  its decode column states no class. It gives families whose rows differ one location each, fails when
  a scalar interpreter disagrees with the documented layout, and gives the
  opcode 0x06 groups the map's slot interval 1 to 8. A layout the map
  leaves open (`date`, `time`, `bytes`, `unknown`) gets no codec and stays
  native and named; the map's `time` wire with the `version` decode is the
  three-byte major, minor and patch triple. New codecs: `uint8`, `int8`,
  `uint8_or_uint16_le`, `uint32_le`, `day_month_year`,
  `hour_minute_second`, `major_minor_patch` and `text`, and the value
  class `version`; points and locations carry their documented
  enumeration or Boolean labels, from the matrix or from a label list the
  map's Values column states outright, and each location carries the value
  class its own documentation states. `DecodeNative` returns the exact native
  value of any point with a known codec (number, Boolean, enumeration
  code and label, date, time, version or text); an enumeration code
  without a candidate label decodes with an empty label, because the map
  does not show such a code to be invalid. The Boolean range and the
  enumeration reading apply only where the documentation of the location
  states them, so a location that states neither keeps its integer. It
  refuses a wrong width, a
  Boolean or date out of range, a non-UTF-8 text, NaN, infinities and
  codec sentinels with the typed `DecodeError` (`out_of_range`,
  `invalid_text`).
  `Point.ForLocation` narrows a point to one location. The SI export rule,
  its withheld reasons, `Decode`, the boiler and heat-pump dictionaries
  and every existing canonical name are unchanged; 37 more controller
  rows become exportable because the documentation now states their codec
  and unit. Refs #524.

- `scripts/vaillant_race_budget.sh` lists the tests in `catalog/vaillant` and
  runs each one alone with the race detector, count 1 and a 90 second bound.
  The bound names a test that does not finish and exits non-zero. It does
  not pass `-short`, does not skip a test, and does not change the 20 minute
  suite timeout. Refs #522. Related: issue 514, which this change does not
  close.

- Add `vaillant/nativepoints`, the native point catalog of the controller
  register service and the boiler and heat-pump register service. It joins
  the controller name annotations, the register-service symbolic
  annotations, the boiler and heat-pump codec dictionaries, the controller
  discovery profiles and the codecs and units of the `vaillant/system`
  scalar interpreters into one row per device class and canonical name:
  operation and selector or group and register, instance range, value
  class, codec and scale, native unit, controller families, name status
  and `Exportable` with a withheld reason (`name_not_final`,
  `identifier_or_code`, `not_numeric`, `unit_unknown`,
  `pressure_reference_unknown`, `codec_unknown`). The evidence grade is
  Hypothesis on every source row and is not the gate. Heat-pump rows keep
  their real status: refrigerant and water pressures are withheld because
  no pressure reference is documented, superheat, subcooling and
  compressor counts and times are withheld until a machine-readable codec
  exists, and a speed with an unresolved unit is withheld. `Decode` turns a
  retained raw reply of an exportable row into an exact native value and
  an exact SI value (kelvin, pascal, joule, watt, second, cubic metre per
  second, per second or a ratio) and refuses NaN, infinities, codec
  sentinels and short or over-long replies with a typed `DecodeError`.
  `scripts/generate_native_points.py --check` compares the committed table
  with its sources. No existing package, name, tag or embedded file
  changes, and the library performs no bus traffic. Follow-up: the gateway
  native point export reads this catalog to publish exportable values from
  the replies it already retains.

- Rename the embedded data files of `catalog/ebusstandard` and
  `catalog/vaillant` whose names started with a work-item code; no
  production file name starts with a code any more. Each file keeps its
  directory, extension and bytes and takes the name of the Go file that
  embeds it; its `go:embed` line and the `ReadFile` paths that name it
  follow:
  - `catalog/ebusstandard/testdata/ab83_03_counter_bodies.json` -> `standard03_counter_bodies.json`
  - `catalog/ebusstandard/testdata/ab84_05_burner_bodies.json` -> `standard05_burner_bodies.json`
  - `catalog/ebusstandard/testdata/ab77_07_candidates.json` -> `standard07_candidates.json`
  - `catalog/ebusstandard/testdata/ab28_standard07_residual_fields.json` -> `standard07_residual_fields.json`
  - `catalog/ebusstandard/testdata/ab79_09_candidates.json` -> `standard09_candidates.json`
  - `catalog/ebusstandard/testdata/ab81_0f_candidates.json` -> `standard0f_candidates.json`
  - `catalog/ebusstandard/testdata/ab49_fe_metadata.json` -> `standardfe_metadata.json`
  - `catalog/ebusstandard/testdata/ab54_standard_ff_body_states.json` -> `standardff_body_states.json`
  - `catalog/vaillant/data/ab33_appliance_info.json` -> `appliance_info.json`
  - `catalog/vaillant/data/ab34_b501_legacy_sequence.json` -> `b501_legacy_sequence.json`
  - `catalog/vaillant/data/ab86_b503_history_flags.json` -> `b503_history_flags.json`
  - `catalog/vaillant/data/ab96_b503_text_categories.json` -> `b503_text_categories.json`
  - `catalog/vaillant/data/ab53_b504_candidates.json` -> `b504_candidates.json`
  - `catalog/vaillant/data/ab15_b504_circuit_roles.json` -> `b504_circuit_roles.json`
  - `catalog/vaillant/data/ab55_b505_candidates.json` -> `b505_candidates.json`
  - `catalog/vaillant/data/ab16_b505_operating_roles.json` -> `b505_operating_roles.json`
  - `catalog/vaillant/data/ab57_b506_candidates.json` -> `b506_candidates.json`
  - `catalog/vaillant/data/ab17_b506_rf_sequence.json` -> `b506_rf_sequence.json`
  - `catalog/vaillant/data/ab13_b507_energy_candidates.json` -> `b507_energy_candidates.json`
  - `catalog/vaillant/data/ab68_b509_boiler_codec_dictionary.json` -> `b509_boiler_codec_dictionary.json`
  - `catalog/vaillant/data/ab59_b509_candidates.json` -> `b509_candidates.json`
  - `catalog/vaillant/data/ab87_b509_diagnostic_shapes.json` -> `b509_diagnostic_shapes.json`
  - `catalog/vaillant/data/ab70_b509_heatpump_codec_dictionary.json` -> `b509_heatpump_codec_dictionary.json`
  - `catalog/vaillant/data/ab46_b509_identity_shapes.json` -> `b509_identity_shapes.json`
  - `catalog/vaillant/data/ab14_b509_op59_quantities.json` -> `b509_op59_quantities.json`
  - `catalog/vaillant/data/ab40_b509_symbolic_annotations.json` -> `b509_symbolic_annotations.json`
  - `catalog/vaillant/data/ab61_b50e_candidates.json` -> `b50e_candidates.json`
  - `catalog/vaillant/data/ab74_b50e_rule_value_rows.json` -> `b50e_rule_value_rows.json`
  - `catalog/vaillant/data/ab63_b50f_candidates.json` -> `b50f_candidates.json`
  - `catalog/vaillant/data/ab12_b50f_status_spans.json` -> `b50f_status_spans.json`
  - `catalog/vaillant/data/ab65_b510_candidates.json` -> `b510_candidates.json`
  - `catalog/vaillant/data/ab48_b510_operating_defaults.json` -> `b510_operating_defaults.json`
  - `catalog/vaillant/data/ab67_b511_candidates.json` -> `b511_candidates.json`
  - `catalog/vaillant/data/ab50_b511_class_candidates.json` -> `b511_class_candidates.json`
  - `catalog/vaillant/data/ab88_b512_air_fields.json` -> `b512_air_fields.json`
  - `catalog/vaillant/data/ab69_b512_candidates.json` -> `b512_candidates.json`
  - `catalog/vaillant/data/ab97_b512_dhw_stepper.json` -> `b512_dhw_stepper.json`
  - `catalog/vaillant/data/ab18_b512_feedback_fallback.json` -> `b512_feedback_fallback.json`
  - `catalog/vaillant/data/ab19_b512_heatpump_io.json` -> `b512_heatpump_io.json`
  - `catalog/vaillant/data/ab38_b513_symbolic_dictionary.json` -> `b513_symbolic_dictionary.json`
  - `catalog/vaillant/data/ab89_b514_menu_feedback.json` -> `b514_menu_feedback.json`
  - `catalog/vaillant/data/ab20_b516_directional_types.json` -> `b516_directional_types.json`
  - `catalog/vaillant/data/ab51_b516_residual_fields.json` -> `b516_residual_fields.json`
  - `catalog/vaillant/data/ab72_b516_statistics_directions.json` -> `b516_statistics_directions.json`
  - `catalog/vaillant/data/ab21_b519_reply_policy.json` -> `b519_reply_policy.json`
  - `catalog/vaillant/data/ab98_b51a_diagnostic_variants.json` -> `b51a_diagnostic_variants.json`
  - `catalog/vaillant/data/ab90_b51a_mirror_fields.json` -> `b51a_mirror_fields.json`
  - `catalog/vaillant/data/ab91_b51cd_partial_labels.json` -> `b51cd_partial_labels.json`
  - `catalog/vaillant/data/ab92_b520_blocks.json` -> `b520_blocks.json`
  - `catalog/vaillant/data/ab71_b523_candidates.json` -> `b523_candidates.json`
  - `catalog/vaillant/data/ab99_b523_config_mixer.json` -> `b523_config_mixer.json`
  - `catalog/vaillant/data/ab93_b523_sensor_fields.json` -> `b523_sensor_fields.json`
  - `catalog/vaillant/data/ab32_b524_component_selectors.json` -> `b524_component_selectors.json`
  - `catalog/vaillant/data/ab95_b524_energy_roles.json` -> `b524_energy_roles.json`
  - `catalog/vaillant/data/ab30_b524_family_codebooks.json` -> `b524_family_codebooks.json`
  - `catalog/vaillant/data/ab36_b524_name_annotations.json` -> `b524_name_annotations.json`
  - `catalog/vaillant/data/ab56_b524_numeric_candidates.json` -> `b524_numeric_candidates.json`
  - `catalog/vaillant/data/ab85_b524_scalar_shapes.json` -> `b524_scalar_shapes.json`
  - `catalog/vaillant/data/ab73_b525_candidates.json` -> `b525_candidates.json`
  - `catalog/vaillant/data/ab75_b531_candidates.json` -> `b531_candidates.json`
  - `catalog/vaillant/data/ab26_b555_additional_forms.json` -> `b555_additional_forms.json`
  - `catalog/vaillant/data/ab94_b555_numeric_units.json` -> `b555_numeric_units.json`
  - `catalog/vaillant/data/ab47_fan_demand.json` -> `fan_demand.json`
  - `catalog/vaillant/data/ab27_holiday.json` -> `holiday.json`
  - `catalog/vaillant/data/ab35_legacy_timer.json` -> `legacy_timer.json`
  - `catalog/vaillant/data/ab39_opaque_programming.json` -> `opaque_programming.json`
  - `catalog/vaillant/data/ab45_outdoor_module.json` -> `outdoor_module.json`
  - `catalog/vaillant/data/ab41_program_blocks.json` -> `program_blocks.json`
  - `catalog/vaillant/data/ab31_pump_roles.json` -> `pump_roles.json`
  - `catalog/vaillant/data/ab29_room_actions.json` -> `room_actions.json`
  - `catalog/vaillant/data/ab43_simulator_adc.json` -> `simulator_adc.json`
  - `catalog/vaillant/data/ab37_tunnel.json` -> `tunnel.json`
  Test-only testdata fixtures keep their names. Exported names, JSON
  contents and catalog identifiers inside the data are unchanged.
  `scripts/rename_catalog_names.py --slice embedded-data` reproduces
  the change and `--check` confirms it.

- Rename the work-item-coded JSON tags, embedded data keys, string values
  and error texts of `catalog/vaillant` to behavior names. A JSON tag or
  data key keeps the part after the code in snake_case (`ab05_family` ->
  `family`, `Ab13An13LoadingActive` -> `loading_active`, `ab16Role` ->
  `role`), because that part is already unique within its type; the
  procedure values of the B51B local budget and the B511 overlap view
  selections drop the code the same way (`ab23an23_grouped_1024` ->
  `grouped_1024`, `ab11an11-op07-base` -> `op07-base`). An error text
  replaces its code with the behavior of its file in words
  (`ab34 duplicate identity` -> `b501 legacy sequence: duplicate
  identity`). Every key of the embedded data files under
  `catalog/vaillant/data` and every Go tag that reads it change together,
  and `scripts/rename_catalog_strings.py --slice catalog-vaillant` refuses a
  key that would collide inside one type or one JSON object. The strings
  are replaced directly: there is no alias and no dual acceptance, so JSON
  with a coded key no longer decodes into the renamed fields. Go
  identifiers, embedded data file names and the row identifiers inside the
  data are unchanged. The loaders no longer require a coded prefix on a row
  identifier; the non-empty and unique checks stay. The test fixtures and test-local tags that
  decode into the renamed types follow, and the `naming-inventory.json` rows
  follow one-to-one. The consumer step is the gateway re-pin onto this
  revision; release bill of materials reconciliation remains required.

  | Previous string | Current string |
  | --- | --- |
  | `Ab05An05Calendar` | `calendar` |
  | `Ab05An05Cell` | `cell` |
  | `Ab05An05Complete` | `complete` |
  | `Ab05An05Components` | `components` |
  | `Ab05An05Context` | `context` |
  | `Ab05An05Contract` | `contract` |
  | `Ab05An05Direction` | `direction` |
  | `Ab05An05Error` | `error` |
  | `Ab05An05ErrorPresent` | `error_present` |
  | `Ab05An05Family` | `family` |
  | `Ab05An05Fields` | `fields` |
  | `Ab05An05Marker` | `marker` |
  | `Ab05An05MarkerPresent` | `marker_present` |
  | `Ab05An05Meaning` | `meaning` |
  | `Ab05An05NormalizedPresent` | `normalized_present` |
  | `Ab05An05Operation` | `operation` |
  | `Ab05An05Qualification` | `qualification` |
  | `Ab05An05Raw` | `raw` |
  | `Ab05An05RawPresent` | `raw_present` |
  | `Ab05An05Remainder` | `remainder` |
  | `Ab05An05RemainderPresent` | `remainder_present` |
  | `Ab05An05State` | `state` |
  | `Ab05An05Value` | `value` |
  | `Ab05An05Variant` | `variant` |
  | `Ab11An11AvailableExtent` | `available_extent` |
  | `Ab11An11CandidateGrade` | `candidate_grade` |
  | `Ab11An11CompleteLogicalExtent` | `complete_logical_extent` |
  | `Ab11An11Direction` | `direction` |
  | `Ab11An11Family` | `family` |
  | `Ab11An11FanRunning` | `fan_running` |
  | `Ab11An11Operation` | `operation` |
  | `Ab11An11PhysicalQualification` | `physical_qualification` |
  | `Ab11An11PressureDomain` | `pressure_domain` |
  | `Ab11An11PressurePercent` | `pressure_percent` |
  | `Ab11An11Raw` | `raw` |
  | `Ab11An11RawFlag` | `raw_flag` |
  | `Ab11An11RawOverlay` | `raw_overlay` |
  | `Ab11An11RawPressure` | `raw_pressure` |
  | `Ab11An11RequiredExtent` | `required_extent` |
  | `Ab11An11RoomSetpointCelsius` | `room_setpoint_celsius` |
  | `Ab11An11Selected` | `selected` |
  | `Ab11An11Selection` | `selection` |
  | `Ab11An11Variant` | `variant` |
  | `Ab13An13Adjustment` | `adjustment` |
  | `Ab13An13Admitted` | `admitted` |
  | `Ab13An13Complete` | `complete` |
  | `Ab13An13ContractGrade` | `contract_grade` |
  | `Ab13An13Defaults` | `defaults` |
  | `Ab13An13Execution` | `execution` |
  | `Ab13An13FieldGrade` | `field_grade` |
  | `Ab13An13Flags` | `flags` |
  | `Ab13An13Flow` | `flow` |
  | `Ab13An13Freshness` | `freshness` |
  | `Ab13An13FrostProtectionActive` | `frost_protection_active` |
  | `Ab13An13HeatDemand` | `heat_demand` |
  | `Ab13An13InletConditionBits` | `inlet_condition_bits` |
  | `Ab13An13InternalTestActive` | `internal_test_active` |
  | `Ab13An13LegionellaDemand` | `legionella_demand` |
  | `Ab13An13LoadingActive` | `loading_active` |
  | `Ab13An13LoadingDisabled` | `loading_disabled` |
  | `Ab13An13LoadingPumpFault` | `loading_pump_fault` |
  | `Ab13An13Modulation` | `modulation` |
  | `Ab13An13NamedDomain` | `named_domain` |
  | `Ab13An13Positive` | `positive` |
  | `Ab13An13Presence` | `presence` |
  | `Ab13An13Qualification` | `qualification` |
  | `Ab13An13Raw` | `raw` |
  | `Ab13An13Ready` | `ready` |
  | `Ab13An13Request` | `request` |
  | `Ab13An13Reset` | `reset` |
  | `Ab13An13ReturnTemperature` | `return_temperature` |
  | `Ab13An13Scalars` | `scalars` |
  | `Ab13An13SensorConditionBits` | `sensor_condition_bits` |
  | `Ab13An13Setpoint` | `setpoint` |
  | `Ab13An13TaskActive` | `task_active` |
  | `Ab13An13Temperatures` | `temperatures` |
  | `Ab13An13Unit` | `unit` |
  | `Ab13An13UnitConflict` | `unit_conflict` |
  | `Ab13An13Unmodelled` | `unmodelled` |
  | `Ab13An13Value` | `value` |
  | `Ab15An15Admitted` | `admitted` |
  | `Ab15An15Candidate` | `candidate` |
  | `Ab15An15Clock` | `clock` |
  | `Ab15An15ContractGrade` | `contract_grade` |
  | `Ab15An15Direction` | `direction` |
  | `Ab15An15Execution` | `execution` |
  | `Ab15An15Failed` | `failed` |
  | `Ab15An15Family` | `family` |
  | `Ab15An15FieldGrade` | `field_grade` |
  | `Ab15An15Freshness` | `freshness` |
  | `Ab15An15Messages` | `messages` |
  | `Ab15An15Mode` | `mode` |
  | `Ab15An15NativeSupport` | `native_support` |
  | `Ab15An15Operation` | `operation` |
  | `Ab15An15PhysicalAdaptation` | `physical_adaptation` |
  | `Ab15An15PositionsGrade` | `positions_grade` |
  | `Ab15An15Presence` | `presence` |
  | `Ab15An15Raw` | `raw` |
  | `Ab15An15Setpoint` | `setpoint` |
  | `Ab15An15SetpointCelsius` | `setpoint_celsius` |
  | `Ab15An15SpecialFunction` | `special_function` |
  | `Ab15An15State` | `state` |
  | `Ab15An15TablesGrade` | `tables_grade` |
  | `Ab15An15Value` | `value` |
  | `Ab15An15Variant` | `variant` |
  | `Ab15An15Words` | `words` |
  | `Ab16Alternative` | `alternative` |
  | `Ab16Direction` | `direction` |
  | `Ab16Family` | `family` |
  | `Ab16Field` | `field` |
  | `Ab16Group` | `group` |
  | `Ab16InstanceScope` | `instance_scope` |
  | `Ab16Operation` | `operation` |
  | `Ab16Primary` | `primary` |
  | `Ab16Register` | `register` |
  | `Ab16Secondary` | `secondary` |
  | `Ab16Variant` | `variant` |
  | `Ab17An17Admitted` | `admitted` |
  | `Ab17An17AvailableExtent` | `available_extent` |
  | `Ab17An17CompleteTwoWord` | `complete_two_word` |
  | `Ab17An17ContractGrade` | `contract_grade` |
  | `Ab17An17DeclaredCount` | `declared_count` |
  | `Ab17An17InterpretationGrade` | `interpretation_grade` |
  | `Ab17An17Qualification` | `qualification` |
  | `Ab17An17Raw` | `raw` |
  | `Ab17An17Status` | `status` |
  | `Ab17An17StatusAvailable` | `status_available` |
  | `Ab17An17StatusMeaning` | `status_meaning` |
  | `Ab17An17Trailing` | `trailing` |
  | `Ab17An17Word0Available` | `word0_available` |
  | `Ab17An17Word0Bits` | `word0_bits` |
  | `Ab17An17Word0Finite` | `word0_finite` |
  | `Ab17An17Word0NativeSign` | `word0_native_sign` |
  | `Ab17An17Word0NativeUnit` | `word0_native_unit` |
  | `Ab17An17Word0Raw` | `word0_raw` |
  | `Ab17An17Word0Value` | `word0_value` |
  | `Ab17An17Word1Available` | `word1_available` |
  | `Ab17An17Word1Raw` | `word1_raw` |
  | `Ab17An17Word1Scale` | `word1_scale` |
  | `Ab17An17Word1Sign` | `word1_sign` |
  | `Ab17An17Word1Type` | `word1_type` |
  | `Ab19An19Admitted` | `admitted` |
  | `Ab19An19Attempt` | `attempt` |
  | `Ab19An19Available` | `available` |
  | `Ab19An19ContractGrade` | `contract_grade` |
  | `Ab19An19Execution` | `execution` |
  | `Ab19An19Fresh` | `fresh` |
  | `Ab19An19InterpretationGrade` | `interpretation_grade` |
  | `Ab19An19LastKnown` | `last_known` |
  | `Ab19An19Maximum` | `maximum` |
  | `Ab19An19Minimum` | `minimum` |
  | `Ab19An19NativeUnit` | `native_unit` |
  | `Ab19An19Presence` | `presence` |
  | `Ab19An19Prior` | `prior` |
  | `Ab19An19Qualification` | `qualification` |
  | `Ab19An19Raw` | `raw` |
  | `Ab19An19State` | `state` |
  | `Ab19An19Value` | `value` |
  | `Ab23An23Chunks` | `chunks` |
  | `Ab23An23Consumed` | `consumed` |
  | `Ab23An23Procedure` | `procedure` |
  | `Ab23An23Unplanned` | `unplanned` |
  | `Ab80Direction` | `direction` |
  | `Ab80Family` | `family` |
  | `Ab80Interpretation` | `interpretation` |
  | `Ab80Operation` | `operation` |
  | `Ab80Present` | `present` |
  | `Ab80Raw` | `raw` |
  | `Ab80Value` | `value` |
  | `Ab80Variant` | `variant` |
  | `Ab82Blocks` | `blocks` |
  | `Ab82Complete` | `complete` |
  | `Ab82Conflict` | `conflict` |
  | `Ab82Conflicts` | `conflicts` |
  | `Ab82Context` | `context` |
  | `Ab82Contributed` | `contributed` |
  | `Ab82Direction` | `direction` |
  | `Ab82Entries` | `entries` |
  | `Ab82Extra` | `extra` |
  | `Ab82Family` | `family` |
  | `Ab82Grade` | `grade` |
  | `Ab82Incomplete` | `incomplete` |
  | `Ab82Index` | `index` |
  | `Ab82Mask` | `mask` |
  | `Ab82Missing` | `missing` |
  | `Ab82Operation` | `operation` |
  | `Ab82Primary` | `primary` |
  | `Ab82Raw` | `raw` |
  | `Ab82RawByte` | `raw_byte` |
  | `Ab82Remainder` | `remainder` |
  | `Ab82Secondary` | `secondary` |
  | `Ab82Selector` | `selector` |
  | `Ab82Signed` | `signed` |
  | `Ab82Start` | `start` |
  | `Ab82State` | `state` |
  | `Ab82Values` | `values` |
  | `Ab82Variant` | `variant` |
  | `ab01ao01: %s must be a byte array` | `b505 sentinel views: %s must be a byte array` |
  | `ab01ao01: %s must be a string` | `b505 sentinel views: %s must be a string` |
  | `ab01ao01: %s must be an integer` | `b505 sentinel views: %s must be an integer` |
  | `ab01ao01: context must be an object` | `b505 sentinel views: context must be an object` |
  | `ab01ao01: duplicate key %q` | `b505 sentinel views: duplicate key %q` |
  | `ab01ao01: fields must be an object` | `b505 sentinel views: fields must be an object` |
  | `ab01ao01: input must be an object` | `b505 sentinel views: input must be an object` |
  | `ab01ao01: invalid byte` | `b505 sentinel views: invalid byte` |
  | `ab01ao01: invalid object key` | `b505 sentinel views: invalid object key` |
  | `ab01ao01: no arguments are defined for this view` | `b505 sentinel views: no arguments are defined for this view` |
  | `ab01ao01: trailing JSON` | `b505 sentinel views: trailing JSON` |
  | `ab01ao01: unexpected delimiter` | `b505 sentinel views: unexpected delimiter` |
  | `ab01ao01: unknown input key %q` | `b505 sentinel views: unknown input key %q` |
  | `ab05 key: opcode, group, instance and register must be supplied and non-null` | `schema key: opcode, group, instance and register must be supplied and non-null` |
  | `ab05 record %d field %d: extent exceeds body` | `schema record %d field %d: extent exceeds body` |
  | `ab05 record %d field %d: invalid extent` | `schema record %d field %d: invalid extent` |
  | `ab05 record %d field %d: invalid grade` | `schema record %d field %d: invalid grade` |
  | `ab05 record %d field %d: width exceeds body` | `schema record %d field %d: width exceeds body` |
  | `ab05 record %d: duplicate alternative key` | `schema record %d: duplicate alternative key` |
  | `ab05 record %d: invalid grade` | `schema record %d: invalid grade` |
  | `ab05 record %d: invalid key dimension` | `schema record %d: invalid key dimension` |
  | `ab05 record %d: negative body width` | `schema record %d: negative body width` |
  | `ab05_alternative` | `alternative` |
  | `ab05_body_width` | `body_width` |
  | `ab05_codec` | `codec` |
  | `ab05_direction` | `direction` |
  | `ab05_enum` | `enum` |
  | `ab05_family` | `family` |
  | `ab05_fields` | `fields` |
  | `ab05_grade` | `grade` |
  | `ab05_group` | `group` |
  | `ab05_instance` | `instance` |
  | `ab05_key` | `key` |
  | `ab05_name` | `name` |
  | `ab05_opcode` | `opcode` |
  | `ab05_position` | `position` |
  | `ab05_raw_fallback` | `raw_fallback` |
  | `ab05_register` | `register` |
  | `ab05_sentinel` | `sentinel` |
  | `ab05_unit` | `unit` |
  | `ab05_variant` | `variant` |
  | `ab05_width` | `width` |
  | `ab11an11-b511-family` | `b511-family` |
  | `ab11an11-op00-fan` | `op00-fan` |
  | `ab11an11-op00-reserved` | `op00-reserved` |
  | `ab11an11-op07-base` | `op07-base` |
  | `ab11an11-op07-extended-room` | `op07-extended-room` |
  | `ab12: duplicate key` | `b50f status spans: duplicate key` |
  | `ab12: incomplete bit partition` | `b50f status spans: incomplete bit partition` |
  | `ab12: incomplete finite catalog` | `b50f status spans: incomplete finite catalog` |
  | `ab12: inconsistent raw byte` | `b50f status spans: inconsistent raw byte` |
  | `ab12: invalid bit` | `b50f status spans: invalid bit` |
  | `ab12: invalid container/raw partition` | `b50f status spans: invalid container/raw partition` |
  | `ab12: invalid entry %d` | `b50f status spans: invalid entry %d` |
  | `ab12: trailing data` | `b50f status spans: trailing data` |
  | `ab13: duplicate embedded key at %d` | `b507 energy candidates: duplicate embedded key at %d` |
  | `ab13: empty embedded catalog` | `b507 energy candidates: empty embedded catalog` |
  | `ab13: invalid divisor at %d` | `b507 energy candidates: invalid divisor at %d` |
  | `ab13: invalid domain at %d` | `b507 energy candidates: invalid domain at %d` |
  | `ab13: invalid embedded entry %d` | `b507 energy candidates: invalid embedded entry %d` |
  | `ab13: trailing embedded data` | `b507 energy candidates: trailing embedded data` |
  | `ab14: duplicate key` | `b509 op59 quantities: duplicate key` |
  | `ab14: expected five finite entries` | `b509 op59 quantities: expected five finite entries` |
  | `ab14: incomplete quantity codec` | `b509 op59 quantities: incomplete quantity codec` |
  | `ab14: invalid codec %d` | `b509 op59 quantities: invalid codec %d` |
  | `ab14: invalid entry %d` | `b509 op59 quantities: invalid entry %d` |
  | `ab14: raw limit has unprovided metadata` | `b509 op59 quantities: raw limit has unprovided metadata` |
  | `ab14: trailing embedded data` | `b509 op59 quantities: trailing embedded data` |
  | `ab15: duplicate key` | `b504 circuit roles: duplicate key` |
  | `ab15: empty catalog` | `b504 circuit roles: empty catalog` |
  | `ab15: incomplete evidence boundary` | `b504 circuit roles: incomplete evidence boundary` |
  | `ab15: incomplete identity` | `b504 circuit roles: incomplete identity` |
  | `ab15: invalid or duplicate id` | `b504 circuit roles: invalid or duplicate id` |
  | `ab15: missing identity axis %s` | `b504 circuit roles: missing identity axis %s` |
  | `ab15ao15 duplicate or invalid key` | `b512 feedback views: duplicate or invalid key` |
  | `ab15ao15 expected object` | `b512 feedback views: expected object` |
  | `ab15ao15 field must be byte array` | `b512 feedback views: field must be byte array` |
  | `ab15ao15 invalid argument` | `b512 feedback views: invalid argument` |
  | `ab15ao15 invalid byte` | `b512 feedback views: invalid byte` |
  | `ab15ao15 null byte` | `b512 feedback views: null byte` |
  | `ab15ao15 null identity` | `b512 feedback views: null identity` |
  | `ab15ao15 null selector` | `b512 feedback views: null selector` |
  | `ab15ao15 trailing JSON` | `b512 feedback views: trailing JSON` |
  | `ab15ao15 unexpected delimiter` | `b512 feedback views: unexpected delimiter` |
  | `ab15ao15 unknown input key` | `b512 feedback views: unknown input key` |
  | `ab16: duplicate identity` | `b505 operating roles: duplicate identity` |
  | `ab16: expected array` | `b505 operating roles: expected array` |
  | `ab16: incomplete string axes` | `b505 operating roles: incomplete string axes` |
  | `ab16: missing identity` | `b505 operating roles: missing identity` |
  | `ab16Accepted` | `accepted` |
  | `ab16Action` | `action` |
  | `ab16BlockedEligibility` | `blocked_eligibility` |
  | `ab16Boiler` | `boiler` |
  | `ab16BranchOrder` | `branch_order` |
  | `ab16CircuitClass` | `circuit_class` |
  | `ab16DisabledValue` | `disabled_value` |
  | `ab16EnabledEligibility` | `enabled_eligibility` |
  | `ab16EnabledValue` | `enabled_value` |
  | `ab16Encoding` | `encoding` |
  | `ab16ExceptionRole` | `exception_role` |
  | `ab16Intent` | `intent` |
  | `ab16Mode` | `mode` |
  | `ab16Nonzero` | `nonzero` |
  | `ab16NumericBranch` | `numeric_branch` |
  | `ab16OtherModes` | `other_modes` |
  | `ab16Period` | `period` |
  | `ab16Persistent` | `persistent` |
  | `ab16ProgrammingMode` | `programming_mode` |
  | `ab16Quantity` | `quantity` |
  | `ab16Relations` | `relations` |
  | `ab16Replacement` | `replacement` |
  | `ab16ReplacementAction` | `replacement_action` |
  | `ab16ReplacementRaw` | `replacement_raw` |
  | `ab16RequestValue` | `request_value` |
  | `ab16Role` | `role` |
  | `ab16SelfHeating` | `self_heating` |
  | `ab16Tank` | `tank` |
  | `ab16Unknown` | `unknown` |
  | `ab16UserCorrection` | `user_correction` |
  | `ab16Zero` | `zero` |
  | `ab17: duplicate key` | `b506 rf sequence: duplicate key` |
  | `ab17: expected array` | `b506 rf sequence: expected array` |
  | `ab17: incomplete metadata` | `b506 rf sequence: incomplete metadata` |
  | `ab17: missing axis %s` | `b506 rf sequence: missing axis %s` |
  | `ab17: null axis %s` | `b506 rf sequence: null axis %s` |
  | `ab17_activation_order` | `activation_order` |
  | `ab17_completed_stop` | `completed_stop` |
  | `ab17_count` | `count` |
  | `ab17_delivery_guarantee` | `delivery_guarantee` |
  | `ab17_extent_bytes` | `extent_bytes` |
  | `ab17_extent_class` | `extent_class` |
  | `ab17_field_position` | `field_position` |
  | `ab17_finite_maximum` | `finite_maximum` |
  | `ab17_finite_meaning` | `finite_meaning` |
  | `ab17_finite_minimum` | `finite_minimum` |
  | `ab17_id` | `id` |
  | `ab17_maximum_bytes` | `maximum_bytes` |
  | `ab17_minimum_bytes` | `minimum_bytes` |
  | `ab17_raw_form` | `raw_form` |
  | `ab17_related_operations` | `related_operations` |
  | `ab17_retry_count` | `retry_count` |
  | `ab17_timeout` | `timeout` |
  | `ab17_unknown` | `unknown` |
  | `ab17_zero_meaning` | `zero_meaning` |
  | `ab18: duplicate identity` | `b512 feedback fallback: duplicate identity` |
  | `ab18: expected JSON array` | `b512 feedback fallback: expected JSON array` |
  | `ab18: incomplete identity` | `b512 feedback fallback: incomplete identity` |
  | `ab18: missing identity axis %s` | `b512 feedback fallback: missing identity axis %s` |
  | `ab18: unexpected identity axes` | `b512 feedback fallback: unexpected identity axes` |
  | `ab19 catalog: %w` | `b512 heatpump io catalog: %w` |
  | `ab19 catalog: duplicate key` | `b512 heatpump io catalog: duplicate key` |
  | `ab19 catalog: empty array` | `b512 heatpump io catalog: empty array` |
  | `ab19 catalog: incomplete identity` | `b512 heatpump io catalog: incomplete identity` |
  | `ab19 catalog: invalid ID` | `b512 heatpump io catalog: invalid ID` |
  | `ab19 catalog: missing %s` | `b512 heatpump io catalog: missing %s` |
  | `ab20 catalog: %w` | `b516 directional types catalog: %w` |
  | `ab20 catalog: duplicate key` | `b516 directional types catalog: duplicate key` |
  | `ab20 catalog: empty array` | `b516 directional types catalog: empty array` |
  | `ab20 catalog: incomplete identity or duplicate id` | `b516 directional types catalog: incomplete identity or duplicate id` |
  | `ab20 catalog: missing primary` | `b516 directional types catalog: missing primary` |
  | `ab21: duplicate identity at %d` | `b519 reply policy: duplicate identity at %d` |
  | `ab21: empty catalog` | `b519 reply policy: empty catalog` |
  | `ab21: incomplete identity at %d` | `b519 reply policy: incomplete identity at %d` |
  | `ab21: invalid identity at %d` | `b519 reply policy: invalid identity at %d` |
  | `ab21: missing axis %s` | `b519 reply policy: missing axis %s` |
  | `ab21_actual_acknowledgement` | `actual_acknowledgement` |
  | `ab21_block_number_in_data` | `block_number_in_data` |
  | `ab21_crc_query_operation` | `crc_query_operation` |
  | `ab21_data` | `data` |
  | `ab21_delivery` | `delivery` |
  | `ab21_grade` | `grade` |
  | `ab21_id` | `id` |
  | `ab21_initiator_relation` | `initiator_relation` |
  | `ab21_metadata_presence` | `metadata_presence` |
  | `ab21_minimum_length` | `minimum_length` |
  | `ab21_native_contract` | `native_contract` |
  | `ab21_native_support` | `native_support` |
  | `ab21_peer_capacity` | `peer_capacity` |
  | `ab21_persistence` | `persistence` |
  | `ab21_physical_effects` | `physical_effects` |
  | `ab21_preserve_raw` | `preserve_raw` |
  | `ab21_qualification` | `qualification` |
  | `ab21_raw_form` | `raw_form` |
  | `ab21_reply_policy` | `reply_policy` |
  | `ab21_reply_recipient` | `reply_recipient` |
  | `ab21_representation` | `representation` |
  | `ab21_routing` | `routing` |
  | `ab21_service` | `service` |
  | `ab21_storage_effects` | `storage_effects` |
  | `ab21_unknown` | `unknown` |
  | `ab21_unresolved_bytes` | `unresolved_bytes` |
  | `ab23an23: negative progress or length` | `b51b local budget crc: negative progress or length` |
  | `ab23an23: progress overflow` | `b51b local budget crc: progress overflow` |
  | `ab23an23: unknown procedure` | `b51b local budget crc: unknown procedure` |
  | `ab23an23_grouped_1024` | `grouped_1024` |
  | `ab23an23_single_205` | `single_205` |
  | `ab25ao25: argument %q is not a boolean or integer` | `b522 demand views: argument %q is not a boolean or integer` |
  | `ab25ao25: arguments is not an object` | `b522 demand views: arguments is not an object` |
  | `ab25ao25: context missing or not an object` | `b522 demand views: context missing or not an object` |
  | `ab25ao25: context requires integer primary, secondary, operation and string direction, view` | `b522 demand views: context requires integer primary, secondary, operation and string direction, view` |
  | `ab25ao25: duplicate key %q` | `b522 demand views: duplicate key %q` |
  | `ab25ao25: field %q has a non-integer element` | `b522 demand views: field %q has a non-integer element` |
  | `ab25ao25: field %q has an element outside 0..255` | `b522 demand views: field %q has an element outside 0..255` |
  | `ab25ao25: field %q is not an array` | `b522 demand views: field %q is not an array` |
  | `ab25ao25: fields missing or not an object` | `b522 demand views: fields missing or not an object` |
  | `ab25ao25: input is not an object` | `b522 demand views: input is not an object` |
  | `ab25ao25: invalid input: %w` | `b522 demand views: invalid input: %w` |
  | `ab25ao25: nesting too deep` | `b522 demand views: nesting too deep` |
  | `ab25ao25: object key is not a string` | `b522 demand views: object key is not a string` |
  | `ab25ao25: trailing data after input` | `b522 demand views: trailing data after input` |
  | `ab26: duplicate identity at row %d` | `b555 additional forms: duplicate identity at row %d` |
  | `ab26: expected JSON array` | `b555 additional forms: expected JSON array` |
  | `ab26: incomplete identity at row %d` | `b555 additional forms: incomplete identity at row %d` |
  | `ab26: missing identity axis %s` | `b555 additional forms: missing identity axis %s` |
  | `ab30 duplicate identity` | `b524 family codebooks: duplicate identity` |
  | `ab30 entries must be an array` | `b524 family codebooks: entries must be an array` |
  | `ab30 missing identity axis %s` | `b524 family codebooks: missing identity axis %s` |
  | `ab30 string axes require explicit tokens` | `b524 family codebooks: string axes require explicit tokens` |
  | `ab32: duplicate key` | `b524 component selectors: duplicate key` |
  | `ab32: empty identity string` | `b524 component selectors: empty identity string` |
  | `ab32: expected array` | `b524 component selectors: expected array` |
  | `ab32: incomplete identity` | `b524 component selectors: incomplete identity` |
  | `ab32: invalid or duplicate id` | `b524 component selectors: invalid or duplicate id` |
  | `ab32: missing %s` | `b524 component selectors: missing %s` |
  | `ab32: null primary` | `b524 component selectors: null primary` |
  | `ab34 array required` | `b501 legacy sequence: array required` |
  | `ab34 duplicate identity` | `b501 legacy sequence: duplicate identity` |
  | `ab34 incomplete identity` | `b501 legacy sequence: incomplete identity` |
  | `ab34 missing identity axis %s` | `b501 legacy sequence: missing identity axis %s` |
  | `ab34_alternative` | `alternative` |
  | `ab34_candidate_grade` | `candidate_grade` |
  | `ab34_context_only` | `context_only` |
  | `ab34_cooldown_event` | `cooldown_event` |
  | `ab34_delivery` | `delivery` |
  | `ab34_direction` | `direction` |
  | `ab34_family` | `family` |
  | `ab34_field` | `field` |
  | `ab34_group` | `group` |
  | `ab34_id` | `id` |
  | `ab34_identity` | `identity` |
  | `ab34_instance_scope` | `instance_scope` |
  | `ab34_native_support` | `native_support` |
  | `ab34_operation` | `operation` |
  | `ab34_payload` | `payload` |
  | `ab34_persistence` | `persistence` |
  | `ab34_physical_effects` | `physical_effects` |
  | `ab34_policy` | `policy` |
  | `ab34_positive_status_service` | `positive_status_service` |
  | `ab34_primary` | `primary` |
  | `ab34_raw_form` | `raw_form` |
  | `ab34_register` | `register` |
  | `ab34_representation_grade` | `representation_grade` |
  | `ab34_request_service` | `request_service` |
  | `ab34_retries` | `retries` |
  | `ab34_return_codec` | `return_codec` |
  | `ab34_return_data_service` | `return_data_service` |
  | `ab34_secondary` | `secondary` |
  | `ab34_segmentation` | `segmentation` |
  | `ab34_sequence` | `sequence` |
  | `ab34_service` | `service` |
  | `ab34_service_role` | `service_role` |
  | `ab34_state_order_beyond_wait` | `state_order_beyond_wait` |
  | `ab34_status_is_return_data` | `status_is_return_data` |
  | `ab34_success_cooldown_required` | `success_cooldown_required` |
  | `ab34_synchronization_bytes` | `synchronization_bytes` |
  | `ab34_synchronization_required` | `synchronization_required` |
  | `ab34_synchronization_token` | `synchronization_token` |
  | `ab34_timeout_cooldown_min_ms` | `timeout_cooldown_min_ms` |
  | `ab34_timeout_duration` | `timeout_duration` |
  | `ab34_transport_adaptation` | `transport_adaptation` |
  | `ab34_unknown` | `unknown` |
  | `ab34_variant` | `variant` |
  | `ab36: documentary name field is not runtime data` | `b524 name annotations: documentary name field is not runtime data` |
  | `ab36: duplicate key` | `b524 name annotations: duplicate key` |
  | `ab36: expected an array` | `b524 name annotations: expected an array` |
  | `ab36: incomplete key` | `b524 name annotations: incomplete key` |
  | `ab36: missing key` | `b524 name annotations: missing key` |
  | `ab38 dictionary: %w` | `b513 symbolic dictionary: %w` |
  | `ab38 dictionary: duplicate key` | `b513 symbolic dictionary: duplicate key` |
  | `ab38 dictionary: expected nonempty array` | `b513 symbolic dictionary: expected nonempty array` |
  | `ab38 row: %w` | `b513 symbolic dictionary row: %w` |
  | `ab38 row: duplicate member` | `b513 symbolic dictionary row: duplicate member` |
  | `ab38 row: expected object` | `b513 symbolic dictionary row: expected object` |
  | `ab38 row: incomplete dictionary identity` | `b513 symbolic dictionary row: incomplete dictionary identity` |
  | `ab38 row: invalid candidate boundary or raw form` | `b513 symbolic dictionary row: invalid candidate boundary or raw form` |
  | `ab38 row: missing member %s` | `b513 symbolic dictionary row: missing member %s` |
  | `ab38 row: null required member` | `b513 symbolic dictionary row: null required member` |
  | `ab38 row: trailing data` | `b513 symbolic dictionary row: trailing data` |
  | `ab40: duplicate key` | `b509 symbolic annotations: duplicate key` |
  | `ab40: empty string axis; use none when inapplicable` | `b509 symbolic annotations: empty string axis; use none when inapplicable` |
  | `ab40: expected array` | `b509 symbolic annotations: expected array` |
  | `ab40: missing axis %s` | `b509 symbolic annotations: missing axis %s` |
  | `ab40: null optional axis %s` | `b509 symbolic annotations: null optional axis %s` |
  | `ab46: duplicate key at row %d` | `b509 identity shapes: duplicate key at row %d` |
  | `ab46: expected an array` | `b509 identity shapes: expected an array` |
  | `ab46: incomplete key at row %d` | `b509 identity shapes: incomplete key at row %d` |
  | `ab46: invalid key` | `b509 identity shapes: invalid key` |
  | `ab46: row %d must be an object` | `b509 identity shapes: row %d must be an object` |
  | `ab48 catalog: duplicate key` | `b510 operating defaults catalog: duplicate key` |
  | `ab48 catalog: empty key axis; use none for non-applicable strings` | `b510 operating defaults catalog: empty key axis; use none for non-applicable strings` |
  | `ab48 catalog: expected array` | `b510 operating defaults catalog: expected array` |
  | `ab48 catalog: missing key axis %s` | `b510 operating defaults catalog: missing key axis %s` |
  | `ab50: duplicate key` | `b511 class candidates: duplicate key` |
  | `ab50: expected array` | `b511 class candidates: expected array` |
  | `ab50: incomplete key` | `b511 class candidates: incomplete key` |
  | `ab50: missing key` | `b511 class candidates: missing key` |
  | `ab50: missing primary` | `b511 class candidates: missing primary` |
  | `ab51: duplicate key` | `b516 residual fields: duplicate key` |
  | `ab51: empty string axis; use none when not applicable` | `b516 residual fields: empty string axis; use none when not applicable` |
  | `ab51: expected array` | `b516 residual fields: expected array` |
  | `ab51: missing key axis %s` | `b516 residual fields: missing key axis %s` |
  | `ab56: all caller identity strings are required` | `b524 numeric candidates: all caller identity strings are required` |
  | `ab56: duplicate or invalid JSON member` | `b524 numeric candidates: duplicate or invalid JSON member` |
  | `ab56: duplicate tuple key at row %d` | `b524 numeric candidates: duplicate tuple key at row %d` |
  | `ab56: expected a tuple array` | `b524 numeric candidates: expected a tuple array` |
  | `ab56: invalid or promoted metadata at row %d` | `b524 numeric candidates: invalid or promoted metadata at row %d` |
  | `ab56: row %d has an incomplete or identity-qualified tuple key` | `b524 numeric candidates: row %d has an incomplete or identity-qualified tuple key` |
  | `ab56: row %d must be an object` | `b524 numeric candidates: row %d must be an object` |
  | `ab56: unexpected delimiter` | `b524 numeric candidates: unexpected delimiter` |
  | `ab68: duplicate key at row %d` | `b509 boiler codec dictionary: duplicate key at row %d` |
  | `ab68: empty dictionary` | `b509 boiler codec dictionary: empty dictionary` |
  | `ab68: incomplete declaration at row %d` | `b509 boiler codec dictionary: incomplete declaration at row %d` |
  | `ab68: incomplete key at row %d` | `b509 boiler codec dictionary: incomplete key at row %d` |
  | `ab68: invalid %s grade at row %d` | `b509 boiler codec dictionary: invalid %s grade at row %d` |
  | `ab68: missing metadata grade at row %d` | `b509 boiler codec dictionary: missing metadata grade at row %d` |
  | `ab68: selector byte out of range at row %d` | `b509 boiler codec dictionary: selector byte out of range at row %d` |
  | `ab68: trailing JSON: %v` | `b509 boiler codec dictionary: trailing JSON: %v` |
  | `ab70: duplicate key` | `b509 heatpump codec dictionary: duplicate key` |
  | `ab70: empty dictionary` | `b509 heatpump codec dictionary: empty dictionary` |
  | `ab70: incomplete or invalid entry` | `b509 heatpump codec dictionary: incomplete or invalid entry` |
  | `ab70: trailing data` | `b509 heatpump codec dictionary: trailing data` |
  | `ab72 catalog: %w` | `b516 statistics directions catalog: %w` |
  | `ab72 catalog: duplicate key` | `b516 statistics directions catalog: duplicate key` |
  | `ab72 catalog: empty` | `b516 statistics directions catalog: empty` |
  | `ab72 catalog: incomplete key` | `b516 statistics directions catalog: incomplete key` |
  | `ab72 catalog: invalid bit` | `b516 statistics directions catalog: invalid bit` |
  | `ab72 catalog: invalid evidence` | `b516 statistics directions catalog: invalid evidence` |
  | `ab72 catalog: invalid qualification` | `b516 statistics directions catalog: invalid qualification` |
  | `ab72 catalog: invalid width` | `b516 statistics directions catalog: invalid width` |
  | `ab72 catalog: trailing input` | `b516 statistics directions catalog: trailing input` |
  | `ab74 duplicate row key` | `b50e rule value rows: duplicate row key` |
  | `ab74 expected nineteen rows` | `b50e rule value rows: expected nineteen rows` |
  | `ab74 incomplete qualification` | `b50e rule value rows: incomplete qualification` |
  | `ab74 invalid field declaration` | `b50e rule value rows: invalid field declaration` |
  | `ab74 invalid row identity` | `b50e rule value rows: invalid row identity` |
  | `ab74 invalid value base` | `b50e rule value rows: invalid value base` |
  | `ab74 missing explicit metadata field` | `b50e rule value rows: missing explicit metadata field` |
  | `ab74 missing known width` | `b50e rule value rows: missing known width` |
  | `ab74 missing or extra companion` | `b50e rule value rows: missing or extra companion` |
  | `ab74 opaque width must remain absent` | `b50e rule value rows: opaque width must remain absent` |
  | `ab74 trailing JSON: %v` | `b50e rule value rows: trailing JSON: %v` |
  | `ab95 metadata: %w` | `b524 energy roles metadata: %w` |
  | `ab95 metadata: duplicate key at %d` | `b524 energy roles metadata: duplicate key at %d` |
  | `ab95 metadata: expected five targets` | `b524 energy roles metadata: expected five targets` |
  | `ab95 metadata: invalid entry %d` | `b524 energy roles metadata: invalid entry %d` |
  | `ab95 metadata: trailing data` | `b524 energy roles metadata: trailing data` |
  | `ab96: duplicate embedded key` | `b503 text categories: duplicate embedded key` |
  | `ab96: empty embedded catalog` | `b503 text categories: empty embedded catalog` |
  | `ab96: invalid embedded entry %d` | `b503 text categories: invalid embedded entry %d` |
  | `ab96: invalid embedded width or position` | `b503 text categories: invalid embedded width or position` |
  | `ab96: missing embedded primary in entry %d` | `b503 text categories: missing embedded primary in entry %d` |
  | `ab96: trailing embedded data: %v` | `b503 text categories: trailing embedded data: %v` |
  | `ab97 duplicate key` | `b512 dhw stepper: duplicate key` |
  | `ab97 empty catalog` | `b512 dhw stepper: empty catalog` |
  | `ab97 invalid bit` | `b512 dhw stepper: invalid bit` |
  | `ab97 invalid divisor` | `b512 dhw stepper: invalid divisor` |
  | `ab97 invalid entry %d` | `b512 dhw stepper: invalid entry %d` |
  | `ab97 trailing data` | `b512 dhw stepper: trailing data` |
  | `ab98 catalog: %w` | `b51a diagnostic variants catalog: %w` |
  | `ab98 catalog: default outside domain` | `b51a diagnostic variants catalog: default outside domain` |
  | `ab98 catalog: duplicate key` | `b51a diagnostic variants catalog: duplicate key` |
  | `ab98 catalog: empty entries` | `b51a diagnostic variants catalog: empty entries` |
  | `ab98 catalog: incomplete key` | `b51a diagnostic variants catalog: incomplete key` |
  | `ab98 catalog: invalid alternative` | `b51a diagnostic variants catalog: invalid alternative` |
  | `ab98 catalog: invalid bit` | `b51a diagnostic variants catalog: invalid bit` |
  | `ab98 catalog: invalid byte order` | `b51a diagnostic variants catalog: invalid byte order` |
  | `ab98 catalog: invalid condition` | `b51a diagnostic variants catalog: invalid condition` |
  | `ab98 catalog: invalid domain` | `b51a diagnostic variants catalog: invalid domain` |
  | `ab98 catalog: invalid evidence boundary` | `b51a diagnostic variants catalog: invalid evidence boundary` |
  | `ab98 catalog: invalid width` | `b51a diagnostic variants catalog: invalid width` |
  | `ab98 catalog: trailing data` | `b51a diagnostic variants catalog: trailing data` |
  | `ab99: duplicate key` | `b523 config mixer: duplicate key` |
  | `ab99: empty catalog` | `b523 config mixer: empty catalog` |
  | `ab99: invalid entry %d` | `b523 config mixer: invalid entry %d` |
  | `ab99: invalid key coordinates in entry %d` | `b523 config mixer: invalid key coordinates in entry %d` |
  | `ab99: trailing data` | `b523 config mixer: trailing data` |

- Document the projection node identifier as a native-only snapshot
  locator. `NodeID` and `StableNodeID` keep their names and values: the
  value is still the canonical Service path. The paths the native providers
  build are anchored at the bus address and the device segment, so for those
  paths the locator changes when a device is re-addressed and repeats when
  another unit of the same model answers at the same address; a provider
  that builds a path without an address gives no such guarantee in either
  direction. The package documentation, the `NodeID` and `StableNodeID` doc
  comments and the per-device projection comments now say that the value is
  a native-only snapshot locator and is not a resource identity. New tests
  cover re-addressing and a same-model swap.

- Rename the work-item-coded JSON tags of the remaining `vaillant/system`
  and `vaillant/b555` types to behavior names: the calendar, codebook,
  constraint-word, discovery-window, value-codec, status, scalar, unscoped
  hint and B555 extent views. A tag keeps the part after the code in
  snake_case (`Ab52StatusAccepted` -> `status_accepted`,
  `Ab07An07RequestedAddress` -> `requested_address`), because that part is
  already unique within its type. The tags are replaced directly: there is
  no alias and no dual acceptance, so JSON with a coded key no longer
  decodes into the renamed fields and consumers move in the same wave. Go
  identifiers are unchanged. The `TestRenamedFieldJSONKeys` pins and the
  testdata keys that decode into these types follow (the discovery windows,
  the constraint-word inputs and the unscoped hint keys and properties); the
  tags of test-local structs, test function names and testdata file names
  are unchanged. After this change no JSON tag of `vaillant/system` and
  `vaillant/b555` carries a work-item code; the coded strings of
  `catalog/vaillant` are renamed in a later change.
  `scripts/rename_catalog_strings.py --slice system-b555-tags` reproduces
  the change; gateway adoption and release bill of materials reconciliation
  remain required.

  | Previous string | Current string |
  | --- | --- |
  | `Ab01An01APIRevision` | `api_revision` |
  | `Ab01An01APIVersion` | `api_version` |
  | `Ab01An01Alternative` | `alternative` |
  | `Ab01An01Candidate` | `candidate` |
  | `Ab01An01Cell` | `cell` |
  | `Ab01An01Cells` | `cells` |
  | `Ab01An01CombinedTimestamp` | `combined_timestamp` |
  | `Ab01An01CompleteCandidate` | `complete_candidate` |
  | `Ab01An01Components` | `components` |
  | `Ab01An01Context` | `context` |
  | `Ab01An01ContractGrade` | `contract_grade` |
  | `Ab01An01EditValidation` | `edit_validation` |
  | `Ab01An01Family` | `family` |
  | `Ab01An01Grade` | `grade` |
  | `Ab01An01Hardware` | `hardware` |
  | `Ab01An01II` | `ii` |
  | `Ab01An01NativeQualification` | `native_qualification` |
  | `Ab01An01OP` | `op` |
  | `Ab01An01Raw` | `raw` |
  | `Ab01An01RequestedGG` | `requested_gg` |
  | `Ab01An01RequestedRR` | `requested_rr` |
  | `Ab01An01ResponseState` | `response_state` |
  | `Ab01An01ReturnedGG` | `returned_gg` |
  | `Ab01An01ReturnedRR` | `returned_rr` |
  | `Ab01An01Software` | `software` |
  | `Ab01An01State` | `state` |
  | `Ab03An03APIRevision` | `api_revision` |
  | `Ab03An03APIVersion` | `api_version` |
  | `Ab03An03Applicable` | `applicable` |
  | `Ab03An03CandidateLabel` | `candidate_label` |
  | `Ab03An03Context` | `context` |
  | `Ab03An03Family` | `family` |
  | `Ab03An03GG` | `gg` |
  | `Ab03An03Grade` | `grade` |
  | `Ab03An03Hardware` | `hardware` |
  | `Ab03An03II` | `ii` |
  | `Ab03An03Integer` | `integer` |
  | `Ab03An03MalformedWidth` | `malformed_width` |
  | `Ab03An03NativeQualification` | `native_qualification` |
  | `Ab03An03OP` | `op` |
  | `Ab03An03RR` | `rr` |
  | `Ab03An03Raw` | `raw` |
  | `Ab03An03RequestedGG` | `requested_gg` |
  | `Ab03An03RequestedRR` | `requested_rr` |
  | `Ab03An03ResponseState` | `response_state` |
  | `Ab03An03ReturnedGG` | `returned_gg` |
  | `Ab03An03ReturnedRR` | `returned_rr` |
  | `Ab03An03Software` | `software` |
  | `Ab03An03Width` | `width` |
  | `Ab03AvailabilityOnly` | `availability_only` |
  | `Ab03Family` | `family` |
  | `Ab03Group` | `group` |
  | `Ab03Instances` | `instances` |
  | `Ab03Maximum` | `maximum` |
  | `Ab03Minimum` | `minimum` |
  | `Ab03Opcode` | `opcode` |
  | `Ab03PredicateRegister` | `predicate_register` |
  | `Ab03RegisterCeiling` | `register_ceiling` |
  | `Ab03Virtual` | `virtual` |
  | `Ab07An07Context` | `context` |
  | `Ab07An07Default` | `default` |
  | `Ab07An07Direction` | `direction` |
  | `Ab07An07Execution` | `execution` |
  | `Ab07An07Factor` | `factor` |
  | `Ab07An07FactorAlgorithm` | `factor_algorithm` |
  | `Ab07An07Family` | `family` |
  | `Ab07An07FamilyApplicability` | `family_applicability` |
  | `Ab07An07Freshness` | `freshness` |
  | `Ab07An07Interpretation` | `interpretation` |
  | `Ab07An07Maximum` | `maximum` |
  | `Ab07An07Minimum` | `minimum` |
  | `Ab07An07NativeEchoRules` | `native_echo_rules` |
  | `Ab07An07Operation` | `operation` |
  | `Ab07An07PacketAssociation` | `packet_association` |
  | `Ab07An07PeerEnforcement` | `peer_enforcement` |
  | `Ab07An07PhysicalQualification` | `physical_qualification` |
  | `Ab07An07Raw` | `raw` |
  | `Ab07An07RawAddressEqual` | `raw_address_equal` |
  | `Ab07An07RequestedAddress` | `requested_address` |
  | `Ab07An07RequestedEditBoundsAdmitted` | `requested_edit_bounds_admitted` |
  | `Ab07An07ReturnedAddress` | `returned_address` |
  | `Ab07An07Selected` | `selected` |
  | `Ab07An07SoftwareContract` | `software_contract` |
  | `Ab07An07State` | `state` |
  | `Ab07An07Step` | `step` |
  | `Ab07An07Value` | `value` |
  | `Ab07An07Variant` | `variant` |
  | `Ab07An07Writability` | `writability` |
  | `Ab11Bytes` | `bytes` |
  | `Ab11Day` | `day` |
  | `Ab11Grade` | `grade` |
  | `Ab11Hour` | `hour` |
  | `Ab11Interpretation` | `interpretation` |
  | `Ab11Latin1` | `latin1` |
  | `Ab11Minute` | `minute` |
  | `Ab11Month` | `month` |
  | `Ab11Qualified` | `qualified` |
  | `Ab11Raw` | `raw` |
  | `Ab11Second` | `second` |
  | `Ab11Year` | `year` |
  | `Ab22BindingZone` | `binding_zone` |
  | `Ab22Context` | `context` |
  | `Ab22CoolingOperatingMode` | `cooling_operating_mode` |
  | `Ab22CoolingTemperatureRaw` | `cooling_temperature_raw` |
  | `Ab22Family` | `family` |
  | `Ab22Fields` | `fields` |
  | `Ab22FrostProtection` | `frost_protection` |
  | `Ab22HeatingOperatingMode` | `heating_operating_mode` |
  | `Ab22HeatingTemperatureRaw` | `heating_temperature_raw` |
  | `Ab22Opcode` | `opcode` |
  | `Ab22Payload` | `payload` |
  | `Ab22SpecialFunctionStatus` | `special_function_status` |
  | `Ab22StatusInfo` | `status_info` |
  | `Ab52AvailableCount` | `available_count` |
  | `Ab52DayCountGrade` | `day_count_grade` |
  | `Ab52DayCounts` | `day_counts` |
  | `Ab52DeclaredCount` | `declared_count` |
  | `Ab52EighthCount` | `eighth_count` |
  | `Ab52EighthCountGrade` | `eighth_count_grade` |
  | `Ab52EndHour` | `end_hour` |
  | `Ab52EndMinute` | `end_minute` |
  | `Ab52Extent` | `extent` |
  | `Ab52Incomplete` | `incomplete` |
  | `Ab52Issues` | `issues` |
  | `Ab52Raw` | `raw` |
  | `Ab52SetpointGrade` | `setpoint_grade` |
  | `Ab52SetpointPresent` | `setpoint_present` |
  | `Ab52SetpointQualified` | `setpoint_qualified` |
  | `Ab52SetpointRaw` | `setpoint_raw` |
  | `Ab52SetpointTenths` | `setpoint_tenths` |
  | `Ab52StartHour` | `start_hour` |
  | `Ab52StartMinute` | `start_minute` |
  | `Ab52Status` | `status` |
  | `Ab52StatusAccepted` | `status_accepted` |
  | `Ab52StatusGrade` | `status_grade` |
  | `Ab52StatusUnavailable` | `status_unavailable` |
  | `Ab52WeekComplete` | `week_complete` |
  | `Ab52WeekValid` | `week_valid` |
  | `Ab58APIRevision` | `api_revision` |
  | `Ab58APIVersion` | `api_version` |
  | `Ab58Codec` | `codec` |
  | `Ab58CodecGrade` | `codec_grade` |
  | `Ab58Context` | `context` |
  | `Ab58DeclaredUnit` | `declared_unit` |
  | `Ab58Delivery` | `delivery` |
  | `Ab58Effects` | `effects` |
  | `Ab58Execution` | `execution` |
  | `Ab58Family` | `family` |
  | `Ab58Float32` | `float32` |
  | `Ab58Freshness` | `freshness` |
  | `Ab58Grade` | `grade` |
  | `Ab58Group` | `group` |
  | `Ab58Hardware` | `hardware` |
  | `Ab58Input` | `input` |
  | `Ab58Instance` | `instance` |
  | `Ab58Opcode` | `opcode` |
  | `Ab58Persistence` | `persistence` |
  | `Ab58Presence` | `presence` |
  | `Ab58Present` | `present` |
  | `Ab58Raw` | `raw` |
  | `Ab58Register` | `register` |
  | `Ab58ReplyKind` | `reply_kind` |
  | `Ab58ReturnedGroup` | `returned_group` |
  | `Ab58ReturnedRegister` | `returned_register` |
  | `Ab58Software` | `software` |
  | `Ab58State` | `state` |
  | `Ab58Support` | `support` |
  | `Ab58Uint16` | `uint16` |
  | `Ab58Value` | `value` |
  | `Ab60APIRevision` | `api_revision` |
  | `Ab60APIVersion` | `api_version` |
  | `Ab60Codec` | `codec` |
  | `Ab60CodecGrade` | `codec_grade` |
  | `Ab60Context` | `context` |
  | `Ab60DeclaredUnit` | `declared_unit` |
  | `Ab60Delivery` | `delivery` |
  | `Ab60Effects` | `effects` |
  | `Ab60Execution` | `execution` |
  | `Ab60Family` | `family` |
  | `Ab60Float32` | `float32` |
  | `Ab60Freshness` | `freshness` |
  | `Ab60Grade` | `grade` |
  | `Ab60Group` | `group` |
  | `Ab60Hardware` | `hardware` |
  | `Ab60Input` | `input` |
  | `Ab60Instance` | `instance` |
  | `Ab60Opcode` | `opcode` |
  | `Ab60Persistence` | `persistence` |
  | `Ab60Presence` | `presence` |
  | `Ab60Present` | `present` |
  | `Ab60Raw` | `raw` |
  | `Ab60Register` | `register` |
  | `Ab60ReplyKind` | `reply_kind` |
  | `Ab60ReturnedGroup` | `returned_group` |
  | `Ab60ReturnedRegister` | `returned_register` |
  | `Ab60Software` | `software` |
  | `Ab60State` | `state` |
  | `Ab60Support` | `support` |
  | `Ab60Uint16` | `uint16` |
  | `Ab60Value` | `value` |
  | `Ab62APIRevision` | `api_revision` |
  | `Ab62APIVersion` | `api_version` |
  | `Ab62Codec` | `codec` |
  | `Ab62CodecGrade` | `codec_grade` |
  | `Ab62Context` | `context` |
  | `Ab62DeclaredUnit` | `declared_unit` |
  | `Ab62Family` | `family` |
  | `Ab62Grade` | `grade` |
  | `Ab62Group` | `group` |
  | `Ab62Hardware` | `hardware` |
  | `Ab62Input` | `input` |
  | `Ab62Instance` | `instance` |
  | `Ab62Opcode` | `opcode` |
  | `Ab62Present` | `present` |
  | `Ab62Raw` | `raw` |
  | `Ab62Register` | `register` |
  | `Ab62ReplyKind` | `reply_kind` |
  | `Ab62ReturnedGroup` | `returned_group` |
  | `Ab62ReturnedRegister` | `returned_register` |
  | `Ab62Software` | `software` |
  | `Ab62State` | `state` |
  | `Ab62Value` | `value` |
  | `Ab64APIRevision` | `api_revision` |
  | `Ab64APIVersion` | `api_version` |
  | `Ab64Codec` | `codec` |
  | `Ab64CodecGrade` | `codec_grade` |
  | `Ab64Context` | `context` |
  | `Ab64Delivery` | `delivery` |
  | `Ab64Effects` | `effects` |
  | `Ab64Execution` | `execution` |
  | `Ab64Family` | `family` |
  | `Ab64Freshness` | `freshness` |
  | `Ab64Grade` | `grade` |
  | `Ab64Group` | `group` |
  | `Ab64Hardware` | `hardware` |
  | `Ab64Input` | `input` |
  | `Ab64Instance` | `instance` |
  | `Ab64InstanceRule` | `instance_rule` |
  | `Ab64NativeSupport` | `native_support` |
  | `Ab64Opcode` | `opcode` |
  | `Ab64Persistence` | `persistence` |
  | `Ab64Presence` | `presence` |
  | `Ab64Present` | `present` |
  | `Ab64Qualification` | `qualification` |
  | `Ab64Raw` | `raw` |
  | `Ab64Register` | `register` |
  | `Ab64ReplyKind` | `reply_kind` |
  | `Ab64ReturnedGroup` | `returned_group` |
  | `Ab64ReturnedRegister` | `returned_register` |
  | `Ab64Software` | `software` |
  | `Ab64State` | `state` |
  | `Ab64Unit` | `unit` |
  | `Ab64Value` | `value` |
  | `Ab66Applicability` | `applicability` |
  | `Ab66Group` | `group` |
  | `Ab66Hint` | `hint` |
  | `Ab66Key` | `key` |
  | `Ab66Max` | `max` |
  | `Ab66Min` | `min` |
  | `Ab66Opcode` | `opcode` |
  | `Ab66Qualified` | `qualified` |
  | `Ab66Register` | `register` |
  | `Ab66Status` | `status` |
  | `Ab66Step` | `step` |
  | `Ab66Type` | `type` |

- Rename the work-item-coded JSON tags, string values and error texts of
  `catalog/ebusstandard`, `vaillant/system` and `vaillant/b555` to behavior
  names. The strings are replaced directly: there is no alias and no dual
  acceptance, so consumers move in the same wave. A JSON tag or value keeps
  the part after the code in snake_case (`Ab76Family` -> `family`,
  `ab09an09_absent` -> `absent`), because that part is already unique
  within its type; an error text names its behavior in words
  (`ab24: missing family` -> `event plan: missing family`). Inputs follow
  the same table: `B509QuantityViews` selects only the directions `request`
  and `response`, and JSON with a coded key no longer decodes into the
  renamed fields. Testdata fixtures follow where they decode into catalog
  types; the tags of test-local structs, test function names, testdata file
  names and the row identifiers of the embedded catalog data are unchanged.
  The coded JSON tags of the remaining `vaillant/system` and `vaillant/b555`
  types and the coded strings of `catalog/vaillant` are renamed in a later
  change. `scripts/rename_catalog_strings.py --slice catalog-strings`
  reproduces the change; gateway adoption and release bill of materials
  reconciliation remain required.

  | Previous string | Current string |
  | --- | --- |
  | `Ab08Address` | `address` |
  | `Ab08Channel` | `channel` |
  | `Ab08Day` | `day` |
  | `Ab08Family` | `family` |
  | `Ab08Identity` | `identity` |
  | `Ab08Instance` | `instance` |
  | `Ab08Raw` | `raw` |
  | `Ab08Start` | `start` |
  | `Ab08Stop` | `stop` |
  | `Ab08SystemType` | `system_type` |
  | `Ab09An09ContractGrade` | `contract_grade` |
  | `Ab09An09Direction` | `direction` |
  | `Ab09An09Family` | `family` |
  | `Ab09An09Flow` | `flow` |
  | `Ab09An09Identity` | `identity` |
  | `Ab09An09MaximumVolume` | `maximum_volume` |
  | `Ab09An09MinimumVolume` | `minimum_volume` |
  | `Ab09An09NativeQualification` | `native_qualification` |
  | `Ab09An09Operation` | `operation` |
  | `Ab09An09Raw` | `raw` |
  | `Ab09An09Selected` | `selected` |
  | `Ab09An09Selection` | `selection` |
  | `Ab09An09SemanticGrade` | `semantic_grade` |
  | `Ab09An09State` | `state` |
  | `Ab09An09Temperature` | `temperature` |
  | `Ab09An09Unit` | `unit` |
  | `Ab09An09Value` | `value` |
  | `Ab09An09Variant` | `variant` |
  | `Ab09An09Volume` | `volume` |
  | `Ab21An21Availability` | `availability` |
  | `Ab21An21Context` | `context` |
  | `Ab21An21Conversion` | `conversion` |
  | `Ab21An21CountMaterialization` | `count_materialization` |
  | `Ab21An21Execution` | `execution` |
  | `Ab21An21Family` | `family` |
  | `Ab21An21Freshness` | `freshness` |
  | `Ab21An21IntermediateStatus` | `intermediate_status` |
  | `Ab21An21Knowledge` | `knowledge` |
  | `Ab21An21LocalCarrier` | `local_carrier` |
  | `Ab21An21Low16` | `low16` |
  | `Ab21An21OriginalFloat64` | `original_float64` |
  | `Ab21An21OriginalInteger64` | `original_integer64` |
  | `Ab21An21Path` | `path` |
  | `Ab21An21Persistence` | `persistence` |
  | `Ab21An21PhysicalQualification` | `physical_qualification` |
  | `Ab21An21Presence` | `presence` |
  | `Ab21An21RawClassification` | `raw_classification` |
  | `Ab21An21RoundedClassification` | `rounded_classification` |
  | `Ab21An21RoundedFloat64` | `rounded_float64` |
  | `Ab21An21ScaledFloat64` | `scaled_float64` |
  | `Ab21An21SoftwareVariant` | `software_variant` |
  | `Ab21An21Status` | `status` |
  | `Ab21An21TimerKind` | `timer_kind` |
  | `Ab21An21TimerKindEffects` | `timer_kind_effects` |
  | `Ab21An21Units` | `units` |
  | `Ab21An21Unsigned32` | `unsigned32` |
  | `Ab21An21WriteAdmission` | `write_admission` |
  | `Ab24Address` | `address` |
  | `Ab24Code` | `code` |
  | `Ab24Family` | `family` |
  | `Ab24Identity` | `identity` |
  | `Ab24Instance` | `instance` |
  | `Ab24OP09` | `op09` |
  | `Ab24OP0B` | `op0b` |
  | `Ab24Opcode` | `opcode` |
  | `Ab24Payload` | `payload` |
  | `Ab24Raw` | `raw` |
  | `Ab24Success` | `success` |
  | `Ab24SystemType` | `system_type` |
  | `Ab24Values` | `values` |
  | `Ab25APIRevision` | `api_revision` |
  | `Ab25APIVersion` | `api_version` |
  | `Ab25An25CandidateMaximum` | `candidate_maximum` |
  | `Ab25An25Conflict` | `conflict` |
  | `Ab25An25ContractGrade` | `contract_grade` |
  | `Ab25An25Data` | `data` |
  | `Ab25An25Direction` | `direction` |
  | `Ab25An25ExecutionGrade` | `execution_grade` |
  | `Ab25An25Family` | `family` |
  | `Ab25An25High` | `high` |
  | `Ab25An25Key` | `key` |
  | `Ab25An25LE` | `le` |
  | `Ab25An25Low` | `low` |
  | `Ab25An25MeaningGrade` | `meaning_grade` |
  | `Ab25An25PhysicalOffsetsGrade` | `physical_offsets_grade` |
  | `Ab25An25PhysicalSupportGrade` | `physical_support_grade` |
  | `Ab25An25PositivePrefixRelevanceGrade` | `positive_prefix_relevance_grade` |
  | `Ab25An25Profile` | `profile` |
  | `Ab25An25Raw` | `raw` |
  | `Ab25An25Service` | `service` |
  | `Ab25An25State` | `state` |
  | `Ab25An25UnresolvedTail` | `unresolved_tail` |
  | `Ab25An25Value` | `value` |
  | `Ab25Codec` | `codec` |
  | `Ab25Family` | `family` |
  | `Ab25Group` | `group` |
  | `Ab25HardwareIdentity` | `hardware_identity` |
  | `Ab25Key` | `key` |
  | `Ab25Opcode` | `opcode` |
  | `Ab25RawDescription` | `raw_description` |
  | `Ab25Register` | `register` |
  | `Ab25RequestedInstance` | `requested_instance` |
  | `Ab25SoftwareIdentity` | `software_identity` |
  | `Ab25State` | `state` |
  | `Ab25Width` | `width` |
  | `Ab42APIRevision` | `api_revision` |
  | `Ab42APIVersion` | `api_version` |
  | `Ab42Capacity` | `capacity` |
  | `Ab42Context` | `context` |
  | `Ab42EligibilityGrade` | `eligibility_grade` |
  | `Ab42Family` | `family` |
  | `Ab42Hardware` | `hardware` |
  | `Ab42HeatPumpCount` | `heat_pump_count` |
  | `Ab42Identifier` | `identifier` |
  | `Ab42InvalidReason` | `invalid_reason` |
  | `Ab42LayoutGrade` | `layout_grade` |
  | `Ab42Meaning` | `meaning` |
  | `Ab42PresenceGrade` | `presence_grade` |
  | `Ab42QualificationGrade` | `qualification_grade` |
  | `Ab42Raw` | `raw` |
  | `Ab42RepresentationGrade` | `representation_grade` |
  | `Ab42Software` | `software` |
  | `Ab42Valid` | `valid` |
  | `Ab42Value` | `value` |
  | `Ab44APIRevision` | `api_revision` |
  | `Ab44APIVersion` | `api_version` |
  | `Ab44Codec` | `codec` |
  | `Ab44Context` | `context` |
  | `Ab44Family` | `family` |
  | `Ab44GG` | `gg` |
  | `Ab44Grade` | `grade` |
  | `Ab44Hardware` | `hardware` |
  | `Ab44II` | `ii` |
  | `Ab44ImplementationValid` | `implementation_valid` |
  | `Ab44Limits` | `limits` |
  | `Ab44Maximum` | `maximum` |
  | `Ab44MaximumRaw` | `maximum_raw` |
  | `Ab44MetadataOP` | `metadata_op` |
  | `Ab44Minimum` | `minimum` |
  | `Ab44MinimumRaw` | `minimum_raw` |
  | `Ab44NumericValid` | `numeric_valid` |
  | `Ab44OP` | `op` |
  | `Ab44RR` | `rr` |
  | `Ab44Raw` | `raw` |
  | `Ab44Reason` | `reason` |
  | `Ab44Software` | `software` |
  | `Ab44Step` | `step` |
  | `Ab44StepMeaning` | `step_meaning` |
  | `Ab44StepRaw` | `step_raw` |
  | `Ab44Width` | `width` |
  | `Ab76Components` | `components` |
  | `Ab76Default` | `default` |
  | `Ab76Direction` | `direction` |
  | `Ab76Encoding` | `encoding` |
  | `Ab76Family` | `family` |
  | `Ab76Function` | `function` |
  | `Ab76ImplementationGrade` | `implementation_grade` |
  | `Ab76LocalCounter` | `local_counter` |
  | `Ab76LocalHours` | `local_hours` |
  | `Ab76LocalMinute` | `local_minute` |
  | `Ab76Maximum` | `maximum` |
  | `Ab76MeaningGrade` | `meaning_grade` |
  | `Ab76Name` | `name` |
  | `Ab76NativeProperties` | `native_properties` |
  | `Ab76Operation` | `operation` |
  | `Ab76Raw` | `raw` |
  | `Ab76Remainder` | `remainder` |
  | `Ab76Selection` | `selection` |
  | `Ab76Selector` | `selector` |
  | `Ab76State` | `state` |
  | `Ab76Value` | `value` |
  | `Ab76Variant` | `variant` |
  | `Ab76Width` | `width` |
  | `Ab78APIGrade` | `api_grade` |
  | `Ab78Complete` | `complete` |
  | `Ab78Default` | `default` |
  | `Ab78Direction` | `direction` |
  | `Ab78Domain` | `domain` |
  | `Ab78Family` | `family` |
  | `Ab78Fields` | `fields` |
  | `Ab78Interpretation` | `interpretation` |
  | `Ab78Key` | `key` |
  | `Ab78Label` | `label` |
  | `Ab78Maximum` | `maximum` |
  | `Ab78MeaningGrade` | `meaning_grade` |
  | `Ab78Minimum` | `minimum` |
  | `Ab78Name` | `name` |
  | `Ab78NativeGrade` | `native_grade` |
  | `Ab78Operation` | `operation` |
  | `Ab78Primary` | `primary` |
  | `Ab78Raw` | `raw` |
  | `Ab78Remainder` | `remainder` |
  | `Ab78Secondary` | `secondary` |
  | `Ab78Selector` | `selector` |
  | `Ab78State` | `state` |
  | `Ab78Unknowns` | `unknowns` |
  | `Ab78Value` | `value` |
  | `Ab78Variant` | `variant` |
  | `ab07_parameter_configuration` | `parameter_configuration` |
  | `ab07_raw_pair` | `raw_pair` |
  | `ab07_slots` | `slots` |
  | `ab07_start_minutes` | `start_minutes` |
  | `ab07_stop_minutes` | `stop_minutes` |
  | `ab07_unused` | `unused` |
  | `ab08 baseline identity mismatch: %w` | `timer plan baseline identity mismatch: %w` |
  | `ab08 invalid explicit identity: %w` | `timer plan invalid explicit identity: %w` |
  | `ab08 requires seven baseline bytes and three slots: %w` | `timer plan requires seven baseline bytes and three slots: %w` |
  | `ab08 slot %d invalid raw ten-minute pair: %w` | `timer plan slot %d invalid raw ten-minute pair: %w` |
  | `ab08 slot %d outside byte range: %w` | `timer plan slot %d outside byte range: %w` |
  | `ab08 unsupported raw selector or non-zone instance: %w` | `timer plan unsupported raw selector or non-zone instance: %w` |
  | `ab09an09_absent` | `absent` |
  | `ab09an09_candidate` | `candidate` |
  | `ab09an09_malformed` | `malformed` |
  | `ab09an09_opaque` | `opaque` |
  | `ab09an09_request` | `request` |
  | `ab09an09_response` | `response` |
  | `ab09an09_unselected` | `unselected` |
  | `ab21an21-float64` | `float64` |
  | `ab21an21-integer64` | `integer64` |
  | `ab24: exactly seven raw values required` | `event plan: exactly seven raw values required` |
  | `ab24: field outside raw u8 range` | `event plan: field outside raw u8 range` |
  | `ab24: invalid address` | `event plan: invalid address` |
  | `ab24: invalid edit opcode` | `event plan: invalid edit opcode` |
  | `ab24: missing family` | `event plan: missing family` |
  | `ab24: missing, failed or uncorrelated baseline %02x` | `event plan: missing, failed or uncorrelated baseline %02x` |
  | `ab24: unsupported system type` | `event plan: unsupported system type` |
  | `ab24: value outside raw u8 range` | `event plan: value outside raw u8 range` |
  | `ab25 conflicting description records` | `description index: conflicting description records` |
  | `ab28: duplicate key` | `standard07 residual fields: duplicate key` |
  | `ab28: expected JSON array` | `standard07 residual fields: expected JSON array` |
  | `ab28: incomplete identity` | `standard07 residual fields: incomplete identity` |
  | `ab28: missing axis %s` | `standard07 residual fields: missing axis %s` |
  | `ab28: missing key` | `standard07 residual fields: missing key` |
  | `ab42 non-finite binary32` | `b524 system information non-finite binary32` |
  | `ab42 requires exactly four normalized bytes` | `b524 system information requires exactly four normalized bytes` |
  | `ab44_bad_width` | `bad_width` |
  | `ab44_codec_grade_required` | `codec_grade_required` |
  | `ab44_codec_width_mismatch` | `codec_width_mismatch` |
  | `ab44_date_time_raw_only` | `date_time_raw_only` |
  | `ab44_echo_mismatch` | `echo_mismatch` |
  | `ab44_incomplete_context` | `incomplete_context` |
  | `ab44_malformed_length` | `malformed_length` |
  | `ab44_negative_step` | `negative_step` |
  | `ab44_nonfinite_limits` | `nonfinite_limits` |
  | `ab44_numeric_limits` | `numeric_limits` |
  | `ab44_reversed_limits` | `reversed_limits` |
  | `ab44_unsupported_codec` | `unsupported_codec` |
  | `ab44_unsupported_op` | `unsupported_op` |
  | `ab54: duplicate JSON member` | `standardff body states: duplicate JSON member` |
  | `ab54: duplicate key` | `standardff body states: duplicate key` |
  | `ab54: expected JSON array` | `standardff body states: expected JSON array` |
  | `ab54: incomplete key; use none for non-applicable string axes` | `standardff body states: incomplete key; use none for non-applicable string axes` |
  | `ab54: unknown key axis` | `standardff body states: unknown key axis` |
  | `ab81 embedded candidate metadata: ` | `standard0f embedded candidate metadata: ` |
  | `ab83 counter metadata: ` | `standard03 counter body metadata: ` |

- Remove work-item codes from the `vaillant/system` and `vaillant/b555`
  file names and Go identifiers; after this change no Go production file
  name and no exported identifier starts with a work-item code. Each
  production file and its test twin is renamed; testdata fixtures and
  test-only files whose name carries a code keep their names:
  - `vaillant/system/ab01an01_b524_calendar_views.go` -> `b524_calendar_views.go`
  - `vaillant/system/ab03_discovery_catalog.go` -> `discovery_catalog.go`
  - `vaillant/system/ab03an03_b524_codebook_views.go` -> `b524_codebook_views.go`
  - `vaillant/system/ab07an07_b509_constraint_word_views.go` -> `b509_constraint_word_views.go`
  - `vaillant/system/ab08_timer_plan.go` -> `timer_plan.go`
  - `vaillant/system/ab09an09_b509_quantity_views.go` -> `b509_quantity_views.go`
  - `vaillant/system/ab11_value_codec.go` -> `value_codec.go`
  - `vaillant/system/ab22_status.go` -> `status.go`
  - `vaillant/system/ab23_event_decode.go` -> `event_decode.go`
  - `vaillant/system/ab24_event_plan.go` -> `event_plan.go`
  - `vaillant/system/ab25_description_index.go` -> `description_index.go`
  - `vaillant/system/ab42_b524_system_information.go` -> `b524_system_information.go`
  - `vaillant/system/ab44_b524_description_spans.go` -> `b524_description_spans.go`
  - `vaillant/system/ab58_b524_system_scalars.go` -> `b524_system_scalars.go`
  - `vaillant/system/ab60_b524_dhw_ventilation_scalars.go` -> `b524_dhw_ventilation_scalars.go`
  - `vaillant/system/ab62_b524_circuit_scalars.go` -> `b524_circuit_scalars.go`
  - `vaillant/system/ab64_b524_zone_scalars.go` -> `b524_zone_scalars.go`
  - `vaillant/system/ab66_b524_unscoped_hints.go` -> `b524_unscoped_hints.go`
  - `vaillant/b555/ab21an21_b555_numeric_narrowing.go` -> `b555_numeric_narrowing.go`
  - `vaillant/b555/ab52_b555_extent_views.go` -> `b555_extent_views.go`
  In `vaillant/system`, identifiers that start with a code (530 exported,
  643 unexported) take the stem of their destination file;
  `Ab<code>X` becomes `<Stem>X` and `ab<code>X` becomes `<stem>X`:
  - `Ab01An01...` -> `B524CalendarView...`
  - `Ab03...` -> `DiscoveryCatalog...`
  - `Ab03An03...` -> `B524CodebookView...`
  - `Ab07...` -> `B524TimerAccess...`
  - `Ab07An07...` -> `B509ConstraintWordView...`
  - `Ab08...` -> `TimerPlan...`
  - `Ab09An09...` -> `B509QuantityView...`
  - `Ab11...` -> `ValueCodec...`
  - `Ab22...` -> `Status...`
  - `Ab23...` -> `EventDecode...`
  - `Ab24...` -> `EventPlan...`
  - `Ab25...` -> `DescriptionIndex...`
  - `Ab42...` -> `B524SystemInformation...`
  - `Ab44...` -> `B524DescriptionSpan...`
  - `Ab58...` -> `B524SystemScalar...`
  - `Ab60...` -> `B524DHWVentilationScalar...`
  - `Ab62...` -> `B524CircuitScalar...`
  - `Ab64...` -> `B524ZoneScalar...`
  - `Ab66...` -> `B524UnscopedHint...`
  Exported names whose remainder repeats a stem word keep the
  stem's service, cut the stem before the first repeated word and
  drop a service the remainder repeats:
  - `Ab03DiscoveryWindow` -> `DiscoveryWindow`
  - `Ab03DiscoveryWindows` -> `DiscoveryWindows`
  - `Ab03LookupDiscoveryWindow` -> `LookupDiscoveryWindow`
  - `Ab03PresentDiscoveryInstances` -> `PresentDiscoveryInstances`
  - `Ab07An07View` -> `B509ConstraintWordView`
  - `Ab07An07Word` -> `B509ConstraintWord`
  - `Ab07TimerSlot` -> `B524TimerSlot`
  - `Ab08BuildTimerPlan` -> `BuildTimerPlan`
  - `Ab08TimerBaseline` -> `TimerBaseline`
  - `Ab08TimerIdentity` -> `TimerIdentity`
  - `Ab08TimerPlan` -> `TimerPlan`
  - `Ab08TimerSlot` -> `TimerSlot`
  - `Ab09An09B509QuantityViews` -> `B509QuantityViews`
  - `Ab09An09QuantityView` -> `B509QuantityView`
  - `Ab11Codec` -> `ValueCodec`
  - `Ab11DecodeValue` -> `DecodeValue`
  - `Ab22DecodeStatus` -> `DecodeStatus`
  - `Ab22SpecialFunctionStatus` -> `SpecialFunctionStatus`
  - `Ab22StatusInfo` -> `StatusInfo`
  - `Ab23DecodeEvent` -> `DecodeEvent`
  - `Ab23Event` -> `Event`
  - `Ab23EventPair` -> `EventPair`
  - `Ab23PairEvents` -> `PairEvents`
  - `Ab24BuildEventPlan` -> `BuildEventPlan`
  - `Ab24EventBaseline` -> `EventBaseline`
  - `Ab24EventEdit` -> `EventEdit`
  - `Ab24EventIdentity` -> `EventIdentity`
  - `Ab24EventPlan` -> `EventPlan`
  - `Ab25DescriptionIndex` -> `DescriptionIndex`
  - `Ab25DescriptionKey` -> `DescriptionKey`
  - `Ab25DescriptionRecord` -> `DescriptionRecord`
  - `Ab25DescriptionState` -> `DescriptionState`
  - `Ab25NewDescriptionIndex` -> `NewDescriptionIndex`
  - `Ab25RawDescription` -> `RawDescription`
  - `Ab42BuildInformation` -> `B524SystemBuildInformation`
  - `Ab42DecodeInformation` -> `B524SystemDecodeInformation`
  - `Ab42InformationResult` -> `B524SystemInformationResult`
  - `Ab44BuildDescription` -> `B524BuildDescription`
  - `Ab44DecodeDescription` -> `B524DecodeDescription`
  - `Ab44Description` -> `B524Description`
  - `Ab58InterpretB524Scalars` -> `B524SystemInterpretScalars`
  - `Ab58Scalar` -> `B524SystemScalar`
  - `Ab60InterpretB524Scalars` -> `B524DHWVentilationInterpretScalars`
  - `Ab60Scalar` -> `B524DHWVentilationScalar`
  - `Ab62InterpretCircuitScalar` -> `B524InterpretCircuitScalar`
  - `Ab66B524ApplicabilityUnknown` -> `B524UnscopedHintApplicabilityUnknown`
  - `Ab66B524HintAbsent` -> `B524UnscopedHintAbsent`
  - `Ab66B524HintGrade` -> `B524UnscopedHintGrade`
  - `Ab66B524HintProperties` -> `B524UnscopedHintProperties`
  - `Ab66B524HintResult` -> `B524UnscopedHintResult`
  - `Ab66B524HintUnscoped` -> `B524HintUnscoped`
  - `Ab66B524UnscopedHints` -> `B524UnscopedHints`
  - `Ab66Hint` -> `B524UnscopedHint`
  - `Ab66LookupB524UnscopedHint` -> `B524LookupUnscopedHint`
  - `ErrAb25DescriptionConflict` -> `ErrDescriptionConflict`
  The import alias `ab09an09hex` would equal another name and takes
  `b509QuantityViewHexPackage`.
  Only the timer slot projection names of `Ab07` are renamed
  (`b524_timer_access.go` keeps its name):
  - `Ab07RawPair` -> `B524TimerAccessRawPair`
  - `Ab07StartMinutes` -> `B524TimerAccessStartMinutes`
  - `Ab07StopMinutes` -> `B524TimerAccessStopMinutes`
  - `Ab07TimerSlot` -> `B524TimerSlot`
  - `Ab07Unused` -> `B524TimerAccessUnused`
  The locals of `b524_timer_access.go` follow:
  - `ab07Index` -> `b524TimerAccessIndex`
  - `ab07Slots` -> `b524TimerAccessSlots`
  - `ab07Start` -> `b524TimerAccessStart`
  - `ab07StartMinutes` -> `b524TimerAccessStartMinutes`
  - `ab07Stop` -> `b524TimerAccessStop`
  - `ab07StopMinutes` -> `b524TimerAccessStopMinutes`
  In `vaillant/b555`, identifiers that start with a code (79 exported,
  92 unexported) take the stem of their destination file;
  `Ab<code>X` becomes `<Stem>X` and `ab<code>X` becomes `<stem>X`:
  - `Ab21An21...` -> `NumericNarrowing...`
  - `Ab52...` -> `ExtentView...`
  Exported names whose remainder repeats a stem word keep the
  stem's service, cut the stem before the first repeated word and
  drop a service the remainder repeats:
  - `Ab52A4View` -> `ExtentA4View`
  - `Ab52A5View` -> `ExtentA5View`
  - `Ab52Extent` -> `Extent`
  The b555 stems leave out the service the package name carries, so no
  name stutters. Four types keep `<Stem><remainder>` because their short
  form is the default name of a struct field: `StatusStatusFields`,
  `B524SystemInformationInformationContext`,
  `B524UnscopedHintB524HintKey` and `B524UnscopedHintB524HintStatus`;
  `ValueCodecValue` and `StatusStatus` keep it because `Value` and
  `Status` are names the package already uses.
  `router/vaillant_timer_access_test.go` and
  `vaillant/system/ab07_timer_decode_test.go` follow the timer slot names.
  Renamed exported struct fields that had no JSON tag carry an explicit
  tag with the unchanged key, and `TestRenamedFieldJSONKeys` of each
  package pins every key. Test strings that name a renamed field for
  reflection follow the rename. Consumers must update
  references; there are no aliases. JSON tags (including
  `ab07_start_minutes`, `ab07_stop_minutes`, `ab07_raw_pair` and
  `ab07_unused`), string values, error texts, testdata names and test
  function names are unchanged.
  `scripts/rename_catalog_names.py --slice vaillant-system-b555`
  reproduces the change. Gateway adoption and release bill of materials
  reconciliation remain required.

- Remove work-item codes from the remaining `catalog/vaillant` files,
  whose code number is 30 to 99, and from the Go identifiers that
  start with their codes; no `catalog/vaillant` file name carries a
  code any more. Each production file and its test twin is renamed;
  embedded data and testdata fixtures keep their names. The mirror
  field catalog takes the plain name `b51a_mirror_fields.go` next to
  `b51a_mirror_fields_evaluator.go`:
  - `ab30_b524_family_codebooks.go` -> `b524_family_codebooks.go`
  - `ab31_pump_roles.go` -> `pump_roles.go`
  - `ab32_b524_component_selectors.go` -> `b524_component_selectors.go`
  - `ab33_appliance_info.go` -> `appliance_info.go`
  - `ab34_b501_legacy_sequence.go` -> `b501_legacy_sequence.go`
  - `ab35_legacy_timer.go` -> `legacy_timer.go`
  - `ab36_b524_name_annotations.go` -> `b524_name_annotations.go`
  - `ab37_tunnel.go` -> `tunnel.go`
  - `ab38_b513_symbolic_dictionary.go` -> `b513_symbolic_dictionary.go`
  - `ab39_opaque_programming.go` -> `opaque_programming.go`
  - `ab40_b509_symbolic_annotations.go` -> `b509_symbolic_annotations.go`
  - `ab41_program_blocks.go` -> `program_blocks.go`
  - `ab43_simulator_adc.go` -> `simulator_adc.go`
  - `ab45_outdoor_module.go` -> `outdoor_module.go`
  - `ab46_b509_identity_shapes.go` -> `b509_identity_shapes.go`
  - `ab47_fan_demand.go` -> `fan_demand.go`
  - `ab48_b510_operating_defaults.go` -> `b510_operating_defaults.go`
  - `ab50_b511_class_candidates.go` -> `b511_class_candidates.go`
  - `ab51_b516_residual_fields.go` -> `b516_residual_fields.go`
  - `ab53_b504_candidates.go` -> `b504_candidates.go`
  - `ab55_b505_candidates.go` -> `b505_candidates.go`
  - `ab56_b524_numeric_candidates.go` -> `b524_numeric_candidates.go`
  - `ab57_b506_candidates.go` -> `b506_candidates.go`
  - `ab59_b509_candidates.go` -> `b509_candidates.go`
  - `ab61_b50e_candidates.go` -> `b50e_candidates.go`
  - `ab63_b50f_candidates.go` -> `b50f_candidates.go`
  - `ab65_b510_candidates.go` -> `b510_candidates.go`
  - `ab67_b511_candidates.go` -> `b511_candidates.go`
  - `ab68_b509_boiler_codec_dictionary.go` -> `b509_boiler_codec_dictionary.go`
  - `ab69_b512_candidates.go` -> `b512_candidates.go`
  - `ab70_b509_heatpump_codec_dictionary.go` -> `b509_heatpump_codec_dictionary.go`
  - `ab71_b523_candidates.go` -> `b523_candidates.go`
  - `ab72_b516_statistics_directions.go` -> `b516_statistics_directions.go`
  - `ab73_b525_candidates.go` -> `b525_candidates.go`
  - `ab74_b50e_rule_value_rows.go` -> `b50e_rule_value_rows.go`
  - `ab75_b531_candidates.go` -> `b531_candidates.go`
  - `ab80_b514_activation_vectors.go` -> `b514_activation_vectors.go`
  - `ab82_b525_spectrum_blocks.go` -> `b525_spectrum_blocks.go`
  - `ab85_b524_scalar_shapes.go` -> `b524_scalar_shapes.go`
  - `ab86_b503_history_flags.go` -> `b503_history_flags.go`
  - `ab87_b509_diagnostic_shapes.go` -> `b509_diagnostic_shapes.go`
  - `ab88_b512_air_fields.go` -> `b512_air_fields.go`
  - `ab89_b514_menu_feedback.go` -> `b514_menu_feedback.go`
  - `ab90_b51a_mirror_fields.go` -> `b51a_mirror_fields.go`
  - `ab91_b51cd_partial_labels.go` -> `b51cd_partial_labels.go`
  - `ab92_b520_blocks.go` -> `b520_blocks.go`
  - `ab93_b523_sensor_fields.go` -> `b523_sensor_fields.go`
  - `ab94_b555_numeric_units.go` -> `b555_numeric_units.go`
  - `ab95_b524_energy_roles.go` -> `b524_energy_roles.go`
  - `ab96_b503_text_categories.go` -> `b503_text_categories.go`
  - `ab97_b512_dhw_stepper.go` -> `b512_dhw_stepper.go`
  - `ab98_b51a_diagnostic_variants.go` -> `b51a_diagnostic_variants.go`
  - `ab99_b523_config_mixer.go` -> `b523_config_mixer.go`
  Identifiers that start with a code (1542 exported, 3093 unexported)
  take the stem of their destination file; `Ab<code>X` becomes
  `<Stem>X` and `ab<code>X` becomes `<stem>X`:
  - `Ab30...` -> `B524FamilyCodebook...`
  - `Ab31...` -> `PumpRole...`
  - `Ab32...` -> `B524ComponentSelector...`
  - `Ab33...` -> `ApplianceInfo...`
  - `Ab34...` -> `B501LegacySequence...`
  - `Ab35...` -> `LegacyTimer...`
  - `Ab36...` -> `B524NameAnnotation...`
  - `Ab37...` -> `Tunnel...`
  - `Ab38...` -> `B513SymbolicDictionary...`
  - `Ab39...` -> `OpaqueProgramming...`
  - `Ab40...` -> `B509SymbolicAnnotation...`
  - `Ab41...` -> `ProgramBlock...`
  - `Ab43...` -> `SimulatorADC...`
  - `Ab45...` -> `OutdoorModule...`
  - `Ab46...` -> `B509IdentityShape...`
  - `Ab47...` -> `FanDemand...`
  - `Ab48...` -> `B510OperatingDefault...`
  - `Ab50...` -> `B511ClassCandidate...`
  - `Ab51...` -> `B516ResidualField...`
  - `Ab53...` -> `B504Candidate...`
  - `Ab55...` -> `B505Candidate...`
  - `Ab56...` -> `B524NumericCandidate...`
  - `Ab57...` -> `B506Candidate...`
  - `Ab59...` -> `B509Candidate...`
  - `Ab61...` -> `B50ECandidate...`
  - `Ab63...` -> `B50FCandidate...`
  - `Ab65...` -> `B510Candidate...`
  - `Ab67...` -> `B511Candidate...`
  - `Ab68...` -> `B509BoilerCodecDictionary...`
  - `Ab69...` -> `B512Candidate...`
  - `Ab70...` -> `B509HeatpumpCodecDictionary...`
  - `Ab71...` -> `B523Candidate...`
  - `Ab72...` -> `B516StatisticsDirection...`
  - `Ab73...` -> `B525Candidate...`
  - `Ab74...` -> `B50ERuleValueRow...`
  - `Ab75...` -> `B531Candidate...`
  - `Ab80...` -> `B514ActivationVector...`
  - `Ab82...` -> `B525SpectrumBlock...`
  - `Ab85...` -> `B524ScalarShape...`
  - `Ab86...` -> `B503HistoryFlag...`
  - `Ab87...` -> `B509DiagnosticShape...`
  - `Ab88...` -> `B512AirField...`
  - `Ab89...` -> `B514MenuFeedback...`
  - `Ab90...` -> `B51AMirrorField...`
  - `Ab91...` -> `B51CDPartialLabel...`
  - `Ab92...` -> `B520Block...`
  - `Ab93...` -> `B523SensorField...`
  - `Ab94...` -> `B555NumericUnit...`
  - `Ab95...` -> `B524EnergyRole...`
  - `Ab96...` -> `B503TextCategory...`
  - `Ab97...` -> `B512DHWStepper...`
  - `Ab98...` -> `B51ADiagnosticVariant...`
  - `Ab99...` -> `B523ConfigMixer...`
  Exported names whose remainder repeats a stem word keep the stem's
  service, cut the stem before the first repeated word and drop a
  service the remainder repeats. Where two renamed names would meet,
  one keeps `<Stem><remainder>`: `Ab67B511Key` ->
  `B511CandidateB511Key` and `Ab92B520Field` -> `B520BlockB520Field`,
  whose short form is the name of a struct field, the type
  `Ab91PartialLabels` -> `B51CDPartialLabelPartialLabels` next to the
  accessor `B51CDPartialLabels`, and nine struct fields such as
  `Ab88Field` -> `B512AirFieldField`:
  - `Ab30Family` -> `B524Family`
  - `Ab31PumpRoles` -> `PumpRoles`
  - `Ab32Components` -> `B524Components`
  - `Ab33ApplianceInfoCatalog` -> `ApplianceInfoCatalog`
  - `Ab35LegacyTimerCatalog` -> `LegacyTimerCatalog`
  - `Ab37TunnelCatalog` -> `TunnelCatalog`
  - `Ab38DictionaryClass` -> `B513SymbolicDictionaryClass`
  - `Ab38DictionaryNamespace` -> `B513SymbolicDictionaryNamespace`
  - `Ab39OpaqueProgramming` -> `OpaqueProgramming`
  - `Ab41ProgramBlocks` -> `ProgramBlocks`
  - `Ab43SimulatorADC` -> `SimulatorADC`
  - `Ab45OutdoorModuleCatalog` -> `OutdoorModuleCatalog`
  - `Ab47FanDemandCatalog` -> `FanDemandCatalog`
  - `Ab51Field` -> `B516ResidualField`
  - `Ab53B504Candidates` -> `B504Candidates`
  - `Ab53LookupB504Candidate` -> `B504LookupCandidate`
  - `Ab55B505Candidates` -> `B505Candidates`
  - `Ab56CandidateSet` -> `B524NumericCandidateSet`
  - `Ab57B506Candidates` -> `B506Candidates`
  - `Ab57B506Lookup` -> `B506CandidateLookup`
  - `Ab59B509Candidates` -> `B509Candidates`
  - `Ab61AllRulesCandidate` -> `B50EAllRulesCandidate`
  - `Ab61B50ECandidates` -> `B50ECandidates`
  - `Ab61LookupB50ECandidate` -> `B50ELookupCandidate`
  - `Ab63B50FCandidates` -> `B50FCandidates`
  - `Ab63B50FLookup` -> `B50FCandidateLookup`
  - `Ab65B510Candidates` -> `B510Candidates`
  - `Ab67B511Catalog` -> `B511CandidateCatalog`
  - `Ab67B511Declarations` -> `B511CandidateDeclarations`
  - `Ab67B511Entry` -> `B511CandidateEntry`
  - `Ab67B511Lookup` -> `B511CandidateLookup`
  - `Ab68CodecDeclaration` -> `B509BoilerCodecDeclaration`
  - `Ab69B512Candidates` -> `B512Candidates`
  - `Ab69Candidate` -> `B512Candidate`
  - `Ab70CodecDeclaration` -> `B509HeatpumpCodecDeclaration`
  - `Ab71B523Candidates` -> `B523Candidates`
  - `Ab72Direction` -> `B516StatisticsDirection`
  - `Ab73B525Candidates` -> `B525Candidates`
  - `Ab74Value` -> `B50ERuleValue`
  - `Ab75B531Candidates` -> `B531Candidates`
  - `Ab75LookupB531` -> `B531CandidateLookup`
  - `Ab80Vector` -> `B514ActivationVector`
  - `Ab82Block` -> `B525SpectrumBlock`
  - `Ab82BlockLength` -> `B525SpectrumBlockLength`
  - `Ab82Blocks` -> `B525SpectrumBlocks`
  - `Ab82JoinBlock` -> `B525SpectrumJoinBlock`
  - `Ab85ScalarShape` -> `B524ScalarShape`
  - `Ab85ScalarShapes` -> `B524ScalarShapes`
  - `Ab86HistoryCatalog` -> `B503HistoryCatalog`
  - `Ab86HistoryField` -> `B503HistoryField`
  - `Ab86HistoryLookup` -> `B503HistoryLookup`
  - `Ab86HistoryMetadata` -> `B503HistoryMetadata`
  - `Ab87B509DiagnosticShapes` -> `B509DiagnosticShapes`
  - `Ab87DiagnosticCatalog` -> `B509DiagnosticCatalog`
  - `Ab87DiagnosticField` -> `B509DiagnosticField`
  - `Ab87DiagnosticKey` -> `B509DiagnosticKey`
  - `Ab87LookupB509DiagnosticField` -> `B509LookupDiagnosticField`
  - `Ab88AirField` -> `B512AirField`
  - `Ab88AirKey` -> `B512AirKey`
  - `Ab88AirMetadata` -> `B512AirMetadata`
  - `Ab88AirOperation` -> `B512AirOperation`
  - `Ab88B512AirFields` -> `B512AirFields`
  - `Ab88LookupB512AirField` -> `B512LookupAirField`
  - `Ab89MenuFeedback` -> `B514MenuFeedback`
  - `Ab90B51AMirrorFields` -> `B51AMirrorFields`
  - `Ab90LookupB51AMirrorField` -> `B51ALookupMirrorField`
  - `Ab90MirrorAlternative` -> `B51AMirrorAlternative`
  - `Ab90MirrorBit` -> `B51AMirrorBit`
  - `Ab90MirrorField` -> `B51AMirrorField`
  - `Ab90MirrorKey` -> `B51AMirrorKey`
  - `Ab90MirrorMetadata` -> `B51AMirrorMetadata`
  - `Ab91B51CDPartialLabels` -> `B51CDPartialLabels`
  - `Ab91LookupB51CDPartialLabel` -> `B51CDLookupPartialLabel`
  - `Ab91PartialField` -> `B51CDPartialField`
  - `Ab91PartialKey` -> `B51CDPartialKey`
  - `Ab91PartialLabel` -> `B51CDPartialLabel`
  - `Ab92B520Blocks` -> `B520Blocks`
  - `Ab92B520Catalog` -> `B520BlockCatalog`
  - `Ab92LookupB520` -> `B520BlockLookup`
  - `Ab93B523SensorFields` -> `B523SensorFields`
  - `Ab93LookupB523SensorField` -> `B523LookupSensorField`
  - `Ab93SensorBranch` -> `B523SensorBranch`
  - `Ab93SensorCatalog` -> `B523SensorCatalog`
  - `Ab93SensorField` -> `B523SensorField`
  - `Ab93SensorKey` -> `B523SensorKey`
  - `Ab93Sensors` -> `B523Sensors`
  - `Ab94B555NumericUnits` -> `B555NumericUnits`
  - `Ab94LookupB555NumericUnits` -> `B555LookupNumericUnits`
  - `Ab94NumericField` -> `B555NumericField`
  - `Ab94NumericMax` -> `B555NumericMax`
  - `Ab94NumericMetadata` -> `B555NumericMetadata`
  - `Ab94NumericMin` -> `B555NumericMin`
  - `Ab94Unit` -> `B555NumericUnit`
  - `Ab94UnitCode` -> `B555NumericUnitCode`
  - `Ab95CommandRole` -> `B524EnergyCommandRole`
  - `Ab98Variant` -> `B51ADiagnosticVariant`
  Import aliases and test locals that would equal another name take
  a `Package` or `Value` suffix:
  - `ab32bytes` -> `b524ComponentSelectorBytesPackage`
  - `ab33json` -> `applianceInfoJSONPackage`
  - `ab41expected` -> `programBlockExpectedValue`
  - `ab41field` -> `programBlockFieldValue`
  - `ab68json` -> `b509BoilerCodecDictionaryJSONPackage`
  - `ab70json` -> `b509HeatpumpCodecDictionaryJSONPackage`
  - `ab95bytes` -> `b524EnergyRoleBytesPackage`
  - `ab96bytes` -> `b503TextCategoryBytesPackage`
  - `ab97bytes` -> `b512DHWStepperBytesPackage`
  The closure parameter of the symbolic dictionary test takes a short
  name so the file stays within its size baseline:
  - `ab38K` -> `key`
  `canonical_names_test.go` follows the `B524NameAnnotation...` names.
  Embedded struct fields keep no JSON key of their own; renamed
  exported struct fields that had no JSON tag carry an explicit tag
  with the unchanged key, and `TestRenamedCatalogFieldJSONKeys` pins
  every key. Test strings that name a renamed field or function for
  reflection or `go/ast` checks, and the production file name the
  partial label test parses, follow the rename. Consumers must update
  references; there are no aliases. JSON tags, string values, error
  texts, `go:embed` paths, data and fixture names and test function
  names are unchanged.
  `scripts/rename_catalog_names.py --slice vaillant-30-99` reproduces
  the change. Gateway adoption and release bill of materials
  reconciliation remain required.

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
