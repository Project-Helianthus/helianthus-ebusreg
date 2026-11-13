// Package solar provides the solar plane for Vaillant-group devices. The
// provider matches on the manufacturer name only and adds the same plane to
// every matching device.
//
// The plane is descriptive and read-only: it lists methods with their
// command bytes and response schema for the registry, but it does not
// implement the router's plane contract, so it neither builds request
// frames nor decodes replies itself. All of its methods name the B504
// operational-data service without an operation byte, so the command
// bytes alone do not tell the methods apart.
package solar
