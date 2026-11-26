package input

import "fmt"

const (
	boilerReadNilResponseText    = "b509 read returned nil response"
	boilerReadFailedText         = "b509 read failed"
	boilerReadFailedFormat       = "b509 read failed: target=0x%02x addr=0x%04x"
	extendedReadNilResponseText  = "b524 read returned nil response"
	extendedReadFailedText       = "b524 read failed"
	extendedReadFailedFormat     = "b524 read failed: opcode=0x%02x group=0x%02x instance=0x%02x addr=0x%04x"
	extendedInstanceMismatchText = "b524 read mismatch: want opcode=0x%02x group=0x%02x instance=0x%02x addr=0x%04x; got reply-instance=0x%02x (group=0x%02x addr=0x%04x)"
	extendedHeaderMismatchText   = "b524 read mismatch: want opcode=0x%02x group=0x%02x instance=0x%02x addr=0x%04x; got kind=0x%02x group=0x%02x addr=0x%04x len=%d"
)

// ReadOutcome is the native result of one read: a value or no value.
type ReadOutcome string

const (
	// OutcomeValue is a read that returned a payload.
	OutcomeValue ReadOutcome = "value"
	// OutcomeNoValue is a read that returned no payload.
	OutcomeNoValue ReadOutcome = "no_value"
)

// Diagnostics are the per-read outcomes the poller produces. SendError is the
// error the send returned at the acquisition boundary. MismatchText is the
// line the poller logged for a mismatched reply and stays empty when no line
// was emitted.
type Diagnostics struct {
	Outcome      ReadOutcome
	SendError    string
	MismatchText string
}

// FailureSurface is where a failure text of the poller can be observed.
type FailureSurface string

const (
	// SurfaceLog is a text the poller writes to its log.
	SurfaceLog FailureSurface = "log"
	// SurfaceSchedulerCallback is an error value inside the scheduler callback that the read functions do not return.
	SurfaceSchedulerCallback FailureSurface = "scheduler_callback"
)

// FailureAnchor ties one poller text to the donor file, function and line that produce it.
type FailureAnchor struct {
	Text     string
	File     string
	Function string
	Line     int
	Surface  FailureSurface
}

// BoilerReadNilResponse returns the text for a B509 read that got no frame.
func BoilerReadNilResponse() string { return boilerReadNilResponseText }

// BoilerReadFailed returns the text for a B509 reply that does not parse.
func BoilerReadFailed(target byte, addr uint16) string {
	return fmt.Sprintf(boilerReadFailedFormat, target, addr)
}

// ExtendedReadNilResponse returns the text for a B524 read that got no frame.
func ExtendedReadNilResponse() string { return extendedReadNilResponseText }

// ExtendedReadFailed returns the text for a B524 reply that does not parse.
func ExtendedReadFailed(opcode, group, instance byte, addr uint16) string {
	return fmt.Sprintf(extendedReadFailedFormat, opcode, group, instance, addr)
}

// ExtendedInstanceMismatch returns the log line for a reply that names the wanted group and register with another instance.
func ExtendedInstanceMismatch(opcode, group, instance byte, addr uint16, replyInstance, replyGroup byte, replyAddr uint16) string {
	return fmt.Sprintf(extendedInstanceMismatchText, opcode, group, instance, addr, replyInstance, replyGroup, replyAddr)
}

// ExtendedHeaderMismatch returns the log line for a reply that names another group or register.
func ExtendedHeaderMismatch(opcode, group, instance byte, addr uint16, kind, replyGroup byte, replyAddr uint16, length int) string {
	return fmt.Sprintf(extendedHeaderMismatchText, opcode, group, instance, addr, kind, replyGroup, replyAddr, length)
}

// FailureAnchors lists the poller texts of a read with the place each is produced.
var FailureAnchors = []FailureAnchor{
	{Text: boilerReadNilResponseText, File: "cmd/gateway/semantic_vaillant_b509_access.go", Function: "readB509ValueWithMaxAge", Line: 293, Surface: SurfaceSchedulerCallback},
	{Text: boilerReadFailedFormat, File: "cmd/gateway/semantic_vaillant_b509_access.go", Function: "readB509ValueWithMaxAge", Line: 297, Surface: SurfaceSchedulerCallback},
	{Text: boilerReadFailedText, File: "cmd/gateway/semantic_vaillant_b509_access.go", Function: "readB509ValueWithMaxAge", Line: 305, Surface: SurfaceSchedulerCallback},
	{Text: extendedReadNilResponseText, File: "cmd/gateway/semantic_vaillant_b524_access.go", Function: "readB524ValueWithMaxAge", Line: 260, Surface: SurfaceSchedulerCallback},
	{Text: extendedReadFailedFormat, File: "cmd/gateway/semantic_vaillant_b524_access.go", Function: "readB524ValueWithMaxAge", Line: 267, Surface: SurfaceSchedulerCallback},
	{Text: extendedReadFailedText, File: "cmd/gateway/semantic_vaillant_b524_access.go", Function: "readB524ValueWithMaxAge", Line: 280, Surface: SurfaceSchedulerCallback},
	{Text: extendedInstanceMismatchText, File: "cmd/gateway/semantic_vaillant_b524_access.go", Function: "parseB524ReadPayload", Line: 342, Surface: SurfaceLog},
	{Text: extendedHeaderMismatchText, File: "cmd/gateway/semantic_vaillant_b524_access.go", Function: "parseB524ReadPayload", Line: 364, Surface: SurfaceLog},
}
