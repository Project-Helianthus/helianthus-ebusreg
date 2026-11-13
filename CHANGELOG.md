# Changelog

## Unreleased

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
