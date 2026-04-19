//go:build tinygo
// +build tinygo

package ebus_standard_catalog

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// ErrTinyGoUnsupported is returned by LoadCatalog on TinyGo builds. The
// full loader depends on gopkg.in/yaml.v3, which is not compatible with
// TinyGo's reflection subset. TinyGo consumers are expected to use
// pre-validated catalogs (e.g. embedded JSON or a generated Go literal)
// produced on the host side; parsing raw YAML at runtime is not supported.
var ErrTinyGoUnsupported = errors.New(
	"ebus_standard: YAML catalog loading is not supported on TinyGo builds")

// loadCatalogImpl is the TinyGo-build stub. It satisfies the dispatch
// declared in loader.go without pulling in gopkg.in/yaml.v3.
func loadCatalogImpl(_ []byte) (Catalog, error) {
	return Catalog{}, ErrTinyGoUnsupported
}

// ComputeContentSHA256 returns the lowercase hex SHA-256 of the given
// bytes. Provided on TinyGo builds so consumers that receive catalogs by
// another path (embedded literal, pre-serialised snapshot) can still
// compute / verify content hashes.
func ComputeContentSHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
