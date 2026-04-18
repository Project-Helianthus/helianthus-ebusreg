package ebus_standard_catalog

// loadCatalogImpl is the real implementation path. M2 RED leaves this as a
// panicking stub; the GREEN implementation replaces the body.
func loadCatalogImpl(data []byte) (Catalog, error) {
	panic("not implemented: ebus_standard catalog loader (M2 RED stub)")
}
