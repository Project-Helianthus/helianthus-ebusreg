package ebusstandard

import _ "embed"

//go:embed testdata/standardfe_metadata.json
var standardFEMetadataFEMetadataJSON []byte

// StandardFEMetadata returns a caller-owned copy of the embedded service FE metadata JSON.
func StandardFEMetadata() []byte {
	return append([]byte(nil), standardFEMetadataFEMetadataJSON...)
}
