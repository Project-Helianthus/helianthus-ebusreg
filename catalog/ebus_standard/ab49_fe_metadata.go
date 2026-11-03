package ebus_standard_catalog

import _ "embed"

//go:embed testdata/ab49_fe_metadata.json
var ab49FEMetadataJSON []byte

func Ab49FEMetadata() []byte {
	return append([]byte(nil), ab49FEMetadataJSON...)
}
