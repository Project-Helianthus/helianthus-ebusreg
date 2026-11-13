// Package dhw provides the domestic hot water plane for Vaillant-group
// devices. The provider matches on the manufacturer name only and adds the
// same plane to every matching device.
//
// The plane is descriptive: it lists methods with their command bytes,
// read-only flag and response schema for the registry, but it does not
// implement the router's plane contract, so it neither builds request
// frames nor decodes replies itself. Its status and parameter reads use
// the B504 operational-data service, its target-temperature write uses
// B505, and the energy-statistics read (B516) is listed only for devices
// whose hardware version is at least 7603.
package dhw
