# Changelog

## Unreleased

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
