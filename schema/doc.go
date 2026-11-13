// Package schema describes eBUS payloads as ordered sequences of named,
// fixed-size fields whose data types come from the transport library, and
// selects the schema that applies to a given device.
//
// Decoding and encoding are strict. Fields are laid out back to back with
// no padding. A field without a type, a field of non-positive size when
// decoding, a payload shorter than the schema, a value missing when
// encoding, or an encoded value of the wrong width fails with an error
// wrapping the transport library's invalid-payload error; nothing partial
// is returned. Bytes beyond the last field are ignored by decoding.
//
// Selection. A selector holds a default schema and an ordered list of
// conditions on the target address and an inclusive hardware-version range;
// the first matching condition wins, otherwise the default applies. The
// hardware version is read as the decimal number formed by all digits in
// the version string. A version without digits, or one whose value exceeds
// 16 bits, matches no condition that sets a version bound.
//
// The JSON loaders are not available in TinyGo builds; there they always
// return the invalid-payload error.
package schema
