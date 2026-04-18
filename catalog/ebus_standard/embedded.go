package ebus_standard_catalog

import _ "embed"

//go:embed catalog.yaml
var embeddedYAMLBytes []byte

// EmbeddedYAML returns the raw bytes of the embedded catalog YAML. The
// returned slice is shared; callers must not modify it.
func EmbeddedYAML() []byte {
	return embeddedYAMLBytes
}

// MustEmbeddedCatalog returns the parsed embedded catalog or panics. This
// is safe because the embedded bytes are compile-time constants and any
// load failure is a build-breaking catalog error that CI catches.
func MustEmbeddedCatalog() Catalog {
	cat, err := LoadCatalog(EmbeddedYAML())
	if err != nil {
		panic("ebus_standard: embedded catalog failed to load: " + err.Error())
	}
	return cat
}
