package ebus_standard_catalog

// EmbeddedYAML returns the raw bytes of the embedded catalog YAML.
//
// RED stub: the real implementation go-embeds catalog.yaml. Until the
// GREEN commit lands, this panics to keep the package compiling while
// making every embedded-catalog test fail loudly.
func EmbeddedYAML() []byte {
	panic("not implemented: embedded ebus_standard catalog YAML (M2 RED stub)")
}

// MustEmbeddedCatalog returns the parsed embedded catalog or panics.
func MustEmbeddedCatalog() Catalog {
	cat, err := LoadCatalog(EmbeddedYAML())
	if err != nil {
		panic("ebus_standard: embedded catalog failed to load: " + err.Error())
	}
	return cat
}

// ComputeContentSHA256 returns the lowercase hex SHA-256 of the given bytes.
// Stub to be implemented alongside the loader.
func ComputeContentSHA256(data []byte) string {
	panic("not implemented: ComputeContentSHA256 (M2 RED stub)")
}
