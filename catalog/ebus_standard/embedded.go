//go:build !tinygo
// +build !tinygo

package ebus_standard_catalog

import _ "embed"

//go:embed catalog.yaml
var embeddedYAMLBytes []byte

// EmbeddedYAML returns a defensive copy of the raw bytes of the embedded
// catalog YAML. A copy is returned (rather than the package-global slice)
// so that mutations by external callers cannot corrupt subsequent catalog
// loads or introduce data races across goroutines. Modifications to the
// returned slice do not affect future calls.
func EmbeddedYAML() []byte {
	return append([]byte(nil), embeddedYAMLBytes...)
}

// MustEmbeddedCatalog returns the parsed embedded catalog or panics. This
// is safe because the embedded bytes are compile-time constants and any
// load failure is a build-breaking catalog error that CI catches.
//
// This symbol is defined only on non-TinyGo builds because the YAML loader
// is unavailable under TinyGo (see loader_tinygo.go). TinyGo consumers
// must construct a Catalog via a pre-validated path (embedded literal or
// host-generated snapshot) rather than calling this helper; a compile-time
// "undefined: MustEmbeddedCatalog" guides them to the correct API instead
// of deterministically panicking at runtime.
func MustEmbeddedCatalog() Catalog {
	cat, err := LoadCatalog(EmbeddedYAML())
	if err != nil {
		panic("ebus_standard: embedded catalog failed to load: " + err.Error())
	}
	return cat
}
